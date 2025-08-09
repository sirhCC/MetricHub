package api

import (
	"net/http"
	"strconv"
	"time"
	"sync"
	"github.com/google/uuid"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirhCC/MetricHub/internal/storage"
	"github.com/sirhCC/MetricHub/pkg/metrics"
	"go.uber.org/zap"
)

// Router holds the dependencies for API handlers
type Router struct {
	logger *zap.Logger
	db     *storage.Database
	redis  *storage.Redis
	mu     sync.RWMutex
	// In-memory development stores
	deployments []metrics.Deployment
	incidents   []metrics.Incident
	calculator  *metrics.DORACalculator
}

// NewRouter creates a new API router with all dependencies
func NewRouter(logger *zap.Logger, db *storage.Database, redis *storage.Redis) *gin.Engine {
	r := &Router{
		logger: logger,
		db:     db,
		redis:  redis,
		calculator: metrics.NewDORACalculator(),
	}

	// Create Gin router
	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(r.requestIDMiddleware())
	router.Use(r.loggingMiddleware())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:5174", "http://localhost:5175"},
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

		// DORA metrics endpoints
		metricsGroup := api.Group("/metrics")
		{
			metricsGroup.GET("/dora", r.getDoraMetrics)
			metricsGroup.GET("/dora/deployment-frequency", r.getDeploymentFrequency)
			metricsGroup.GET("/dora/lead-time", r.getLeadTime)
			metricsGroup.GET("/dora/mttr", r.getMTTR)
			metricsGroup.GET("/dora/change-failure-rate", r.getChangeFailureRate)
		}

		// Plugin endpoints (placeholder for now)
		plugins := api.Group("/plugins")
		{
			plugins.GET("", r.listPlugins)
			plugins.GET("/:name/health", r.pluginHealth)
		}

		// Webhook endpoints (placeholder for now)
		api.POST("/webhook/:plugin", r.handleWebhook)

		// Ingestion (in-memory dev only) + listing
		api.POST("/deployments", r.createDeployment)
		api.GET("/deployments", r.listDeployments)
		api.POST("/incidents", r.createIncident)
		api.GET("/incidents", r.listIncidents)
		api.POST("/incidents/:id/resolve", r.resolveIncident)
		api.GET("/state", r.listState)
	}

	// Serve static files and handle SPA routing
	router.Static("/static", "./web/static")
	router.NoRoute(func(c *gin.Context) {
		// For SPA routing, serve index.html for any non-API routes
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] != "/api" {
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
		requestID := param.Request.Header.Get("X-Request-ID")
		if requestID == "" { // fallback to context value if present
			if v := param.Request.Context().Value("request_id"); v != nil {
				if rid, ok := v.(string); ok {
					requestID = rid
				}
			}
		}
		if requestID == "" { // final safeguard
			requestID = "unknown"
		}
		r.logger.Info("HTTP request",
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.Int("status", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("client_ip", param.ClientIP),
			zap.String("user_agent", param.Request.UserAgent()),
			zap.String("request_id", requestID),
		)
		return ""
	})
}

// requestIDMiddleware injects a unique request ID if not provided
func (r *Router) requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader("X-Request-ID")
		if rid == "" {
			rid = uuid.NewString()
			c.Request.Header.Set("X-Request-ID", rid)
		}
		c.Set("request_id", rid)
		c.Writer.Header().Set("X-Request-ID", rid)
		c.Next()
	}
}

// parseTimeRange parses an optional ?days=N parameter (default 30, max 365)
func (r *Router) parseTimeRange(c *gin.Context) metrics.TimeRange {
	days := 30
	if v := c.Query("days"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 365 { days = n }
	}
	end := time.Now()
	start := end.AddDate(0,0,-days+1)
	return metrics.TimeRange{Start: start, End: end}
}
