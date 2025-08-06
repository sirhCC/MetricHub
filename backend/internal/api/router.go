package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirhCC/MetricHub/internal/storage"
	"go.uber.org/zap"
)

// Router holds the dependencies for API handlers
type Router struct {
	logger *zap.Logger
	db     *storage.Database
	redis  *storage.Redis
}

// NewRouter creates a new API router with all dependencies
func NewRouter(logger *zap.Logger, db *storage.Database, redis *storage.Redis) *gin.Engine {
	r := &Router{
		logger: logger,
		db:     db,
		redis:  redis,
	}

	// Create Gin router
	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(r.loggingMiddleware())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API routes
	api := router.Group("/api/v1")
	{
		// Health check endpoints
		api.GET("/health", r.healthCheck)
		api.GET("/health/database", r.databaseHealth)
		api.GET("/health/redis", r.redisHealth)

		// DORA metrics endpoints (placeholder for now)
		metrics := api.Group("/metrics")
		{
			metrics.GET("/dora", r.getDoraMetrics)
			metrics.GET("/dora/deployment-frequency", r.getDeploymentFrequency)
			metrics.GET("/dora/lead-time", r.getLeadTime)
			metrics.GET("/dora/mttr", r.getMTTR)
			metrics.GET("/dora/change-failure-rate", r.getChangeFailureRate)
		}

		// Plugin endpoints (placeholder for now)
		plugins := api.Group("/plugins")
		{
			plugins.GET("", r.listPlugins)
			plugins.GET("/:name/health", r.pluginHealth)
		}

		// Webhook endpoints (placeholder for now)
		api.POST("/webhook/:plugin", r.handleWebhook)
	}

	// Serve static files and handle SPA routing
	router.Static("/static", "./web/static")
	router.NoRoute(func(c *gin.Context) {
		// For SPA routing, serve index.html for any non-API routes
		if c.Request.URL.Path[:4] != "/api" {
			c.File("./web/index.html")
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "endpoint not found"})
		}
	})

	return router
}

// loggingMiddleware adds request logging
func (r *Router) loggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		r.logger.Info("HTTP request",
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.Int("status", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("client_ip", param.ClientIP),
			zap.String("user_agent", param.Request.UserAgent()),
		)
		return ""
	})
}
