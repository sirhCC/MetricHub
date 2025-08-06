package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version,omitempty"`
	Checks    map[string]string `json:"checks,omitempty"`
}

// healthCheck returns the overall system health
func (r *Router) healthCheck(c *gin.Context) {
	checks := make(map[string]string)

	// Check database
	if err := r.db.Health(); err != nil {
		checks["database"] = "unhealthy: " + err.Error()
	} else {
		checks["database"] = "healthy"
	}

	// Check Redis
	if err := r.redis.Health(); err != nil {
		checks["redis"] = "unhealthy: " + err.Error()
	} else {
		checks["redis"] = "healthy"
	}

	// Determine overall status
	status := "healthy"
	for _, check := range checks {
		if check != "healthy" {
			status = "unhealthy"
			break
		}
	}

	response := HealthResponse{
		Status:    status,
		Timestamp: time.Now(),
		Version:   "0.1.0", // TODO: Get from build info
		Checks:    checks,
	}

	if status == "healthy" {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusServiceUnavailable, response)
	}
}

// databaseHealth checks database connectivity
func (r *Router) databaseHealth(c *gin.Context) {
	if err := r.db.Health(); err != nil {
		c.JSON(http.StatusServiceUnavailable, HealthResponse{
			Status:    "unhealthy",
			Timestamp: time.Now(),
			Checks:    map[string]string{"database": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Checks:    map[string]string{"database": "healthy"},
	})
}

// redisHealth checks Redis connectivity
func (r *Router) redisHealth(c *gin.Context) {
	if err := r.redis.Health(); err != nil {
		c.JSON(http.StatusServiceUnavailable, HealthResponse{
			Status:    "unhealthy",
			Timestamp: time.Now(),
			Checks:    map[string]string{"redis": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Checks:    map[string]string{"redis": "healthy"},
	})
}
