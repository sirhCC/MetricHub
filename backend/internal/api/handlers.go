package api

import (
	"net/http"
	"time"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/sirhCC/MetricHub/pkg/metrics"
	"go.uber.org/zap"
)

// Health check handlers
func (r *Router) healthCheck(c *gin.Context) {
	response := gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"version":   "0.1.0",
		"checks": gin.H{
			"database": "healthy",
			"redis":    "healthy",
		},
	}

	r.logger.Info("Health check requested", zap.String("status", "healthy"))
	c.JSON(http.StatusOK, response)
}

func (r *Router) databaseHealth(c *gin.Context) {
	// TODO: Implement actual database health check
	response := gin.H{
		"status":      "connected",
		"latency":     "12ms",
		"connections": 5,
		"timestamp":   time.Now().UTC(),
	}

	c.JSON(http.StatusOK, response)
}

func (r *Router) redisHealth(c *gin.Context) {
	// TODO: Implement actual Redis health check
	response := gin.H{
		"status":       "connected",
		"memory_usage": "45MB",
		"hit_rate":     "94.2%",
		"timestamp":    time.Now().UTC(),
	}

	c.JSON(http.StatusOK, response)
}

// DORA Metrics handlers
func (r *Router) getDoraMetrics(c *gin.Context) {
	tr := r.parseTimeRange(c)
	result, err := r.calculator.CalculateAll(r.deployments, r.incidents, tr)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"calculation failed"}); return }
	classification := r.calculator.ClassifyPerformance(result)
	overall := r.calculator.GetOverallPerformance(classification)
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"deployment_frequency": result.DeploymentFrequency,
			"lead_time": result.LeadTime.String(),
			"mttr": result.MTTR.String(),
			"change_failure_rate": result.ChangeFailureRate,
			"classification": classification,
			"overall_performance": overall,
		},
		"metadata": gin.H{
			"time_range": gin.H{"start": tr.Start, "end": tr.End},
			"last_updated": time.Now().UTC(),
			"data_quality": result.DataQuality,
			"deployments_count": len(r.deployments),
			"incidents_count": len(r.incidents),
		},
	})
}

func (r *Router) getDeploymentFrequency(c *gin.Context) {
	tr := r.parseTimeRange(c)
	freq := r.calculator.CalculateDeploymentFrequency(r.deployments, tr)
	c.JSON(http.StatusOK, gin.H{"value": freq, "unit":"per_day", "time_range": tr})
}

func (r *Router) getLeadTime(c *gin.Context) {
	lt := r.calculator.CalculateLeadTime(r.deployments)
	c.JSON(http.StatusOK, gin.H{"value": lt.String(), "unit":"duration"})
}

func (r *Router) getMTTR(c *gin.Context) {
	mttr := r.calculator.CalculateMTTR(r.incidents)
	c.JSON(http.StatusOK, gin.H{"value": mttr.String(), "unit":"duration"})
}

func (r *Router) getChangeFailureRate(c *gin.Context) {
	cfr := r.calculator.CalculateChangeFailureRate(r.deployments, r.incidents)
	c.JSON(http.StatusOK, gin.H{"value": cfr, "unit":"ratio"})
}

// Plugin handlers
func (r *Router) listPlugins(c *gin.Context) {
	// TODO: Get actual plugin list from database
	response := gin.H{
		"plugins": []gin.H{
			{
				"id":        "github",
				"name":      "GitHub Integration",
				"version":   "1.2.0",
				"status":    "active",
				"last_sync": time.Now().Add(-2 * time.Minute).UTC(),
			},
			{
				"id":        "jenkins",
				"name":      "Jenkins CI/CD",
				"version":   "1.1.5",
				"status":    "inactive",
				"last_sync": time.Now().Add(-1 * time.Hour).UTC(),
			},
			{
				"id":        "datadog",
				"name":      "Datadog Monitoring",
				"version":   "1.0.3",
				"status":    "error",
				"last_sync": time.Now().Add(-5 * time.Hour).UTC(),
			},
		},
	}

	r.logger.Info("Plugin list requested")
	c.JSON(http.StatusOK, response)
}

func (r *Router) pluginHealth(c *gin.Context) {
	pluginName := c.Param("name")

	// TODO: Get actual plugin health from database
	response := gin.H{
		"plugin":     pluginName,
		"status":     "healthy",
		"last_check": time.Now().UTC(),
		"metrics": gin.H{
			"requests_per_minute": 25,
			"error_rate":         0.02,
			"response_time":      120,
		},
	}

	r.logger.Info("Plugin health requested", zap.String("plugin", pluginName))
	c.JSON(http.StatusOK, response)
}

// Webhook handler
func (r *Router) handleWebhook(c *gin.Context) {
	pluginName := c.Param("plugin")

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		r.logger.Error("Failed to parse webhook payload", 
			zap.String("plugin", pluginName), 
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	// TODO: Process webhook payload and store data in database
	response := gin.H{
		"message":      "Webhook received",
		"plugin":       pluginName,
		"processed_at": time.Now().UTC(),
		"processed":    true,
	}

	r.logger.Info("Webhook processed", 
		zap.String("plugin", pluginName),
		zap.Any("payload", payload))

	c.JSON(http.StatusOK, response)
}

// Helper function to generate mock historical data
func generateMockHistoricalData(metricType string) []gin.H {
	data := make([]gin.H, 30)
	now := time.Now()

	for i := 0; i < 30; i++ {
		date := now.AddDate(0, 0, -(29-i))
		var value interface{}

		switch metricType {
		case "deployment_frequency":
			value = 1 + (i%5) // 1-5 deployments per day
		case "lead_time":
			value = 2 + (i%12) // 2-14 hours
		case "mttr":
			value = 1 + (i%4) // 1-5 hours
		case "change_failure_rate":
			value = 5 + (i%10) // 5-15%
		default:
			value = i
		}

		data[i] = gin.H{
			"date":  date.Format("2006-01-02"),
			"value": value,
		}
	}

	return data
}

// Ingestion endpoints (development in-memory only)
type deploymentRequest struct {
	ID        string     `json:"id"`
	StartedAt *time.Time `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"`
	Status    string     `json:"status"`
	CommitSHA string     `json:"commit_sha"`
	Service   string     `json:"service"`
	Environment string   `json:"environment"`
}

func (r *Router) createDeployment(c *gin.Context) {
	var req deploymentRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid json"}); return }
	start := time.Now()
	if req.StartedAt != nil { start = *req.StartedAt }
	status := metrics.DeploymentStatus(req.Status)
	dep := metrics.Deployment{
		ID: req.ID,
		Service: req.Service,
		Environment: req.Environment,
		StartTime: start,
		EndTime: req.EndedAt,
		Status: status,
		CommitSHA: req.CommitSHA,
		CommitTime: start, // placeholder
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	r.deployments = append(r.deployments, dep)
	c.JSON(http.StatusCreated, gin.H{"deployment": dep})
}

type incidentRequest struct {
	ID           string     `json:"id"`
	StartedAt    *time.Time `json:"started_at"`
	ResolvedAt   *time.Time `json:"resolved_at"`
	Severity     string     `json:"severity"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Service      string     `json:"service"`
	Environment  string     `json:"environment"`
}

func (r *Router) createIncident(c *gin.Context) {
	var req incidentRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid json"}); return }
	start := time.Now()
	if req.StartedAt != nil { start = *req.StartedAt }
	severity := metrics.IncidentSeverity(req.Severity)
	inc := metrics.Incident{
		ID: req.ID,
		Title: req.Title,
		Description: req.Description,
		Service: req.Service,
		Environment: req.Environment,
		Severity: severity,
		StartTime: start,
		ResolvedTime: req.ResolvedAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	r.incidents = append(r.incidents, inc)
	c.JSON(http.StatusCreated, gin.H{"incident": inc})
}

func (r *Router) resolveIncident(c *gin.Context) {
	id := c.Param("id")
	now := time.Now()
	for i := range r.incidents {
		if r.incidents[i].ID == id {
			if r.incidents[i].ResolvedTime != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"already resolved"}); return }
			r.incidents[i].ResolvedTime = &now
			r.incidents[i].UpdatedAt = time.Now()
			c.JSON(http.StatusOK, gin.H{"resolved_at": now}); return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error":"not found"})
}

func (r *Router) listState(c *gin.Context) {
	// stable ordering
	sort.Slice(r.deployments, func(i,j int) bool { return r.deployments[i].StartTime.Before(r.deployments[j].StartTime) })
	sort.Slice(r.incidents, func(i,j int) bool { return r.incidents[i].StartTime.Before(r.incidents[j].StartTime) })
	c.JSON(http.StatusOK, gin.H{"deployments": r.deployments, "incidents": r.incidents})
}
