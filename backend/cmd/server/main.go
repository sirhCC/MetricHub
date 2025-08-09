package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirhCC/MetricHub/internal/api"
	"github.com/sirhCC/MetricHub/internal/config"
	"github.com/sirhCC/MetricHub/internal/storage"
	"go.uber.org/zap"
)

// kvToZap converts a variadic list of key,value pairs into zap fields.
func kvToZap(kv ...interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		k, _ := kv[i].(string)
		v := kv[i+1]
		fields = append(fields, zap.Any(k, v))
	}
	return fields
}

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize storage (optional for development)
	var db *storage.Database
	var redis *storage.Redis

	if cfg.DatabaseURL != "" {
		var err error
		db, err = storage.NewDatabase(cfg.DatabaseURL)
		if err != nil {
			logger.Warn("Failed to initialize database (continuing without it)", zap.Error(err))
		} else {
			defer db.Close()
			logger.Info("Database connection established")
			// Run migrations if enabled
			if cfg.AutoMigrate {
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				if err := storage.ApplyMigrations(ctx, db.GetDB(), "./migrations", func(msg string, kv ...interface{}) { logger.Info(msg, kvToZap(kv...)...) }); err != nil {
					logger.Error("migration failed", zap.Error(err))
				} else {
					logger.Info("database migrations applied")
				}
				cancel()
			}
		}
	} else {
		logger.Info("Running without database connection (development mode)")
	}

	if cfg.RedisURL != "" {
		var err error
		redis, err = storage.NewRedis(cfg.RedisURL)
		if err != nil {
			logger.Warn("Failed to initialize Redis (continuing without it)", zap.Error(err))
		} else {
			defer redis.Close()
			logger.Info("Redis connection established")
		}
	} else {
		logger.Info("Running without Redis connection (development mode)")
	}

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize API router
	router := api.NewRouter(logger, db, redis)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      router,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting MetricHub server",
			zap.String("address", srv.Addr),
			zap.String("environment", cfg.Environment))

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	logger.Info("Shutting down server...", zap.String("signal", sig.String()))

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
