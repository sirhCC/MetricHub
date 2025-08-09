package api

import (
	"net/http"
	"sort"
	"time"

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
	var deps []metrics.Deployment
	var incs []metrics.Incident
	if r.deploymentRepo != nil && r.incidentRepo != nil { // database path
		var err error
		deps, err = r.deploymentRepo.ListRange(c.Request.Context(), tr.Start, tr.End)
		if err != nil { respondError(c, ErrInternal, "failed to load deployments", nil); return }
		incs, err = r.incidentRepo.ListRange(c.Request.Context(), tr.Start, tr.End)
		if err != nil { respondError(c, ErrInternal, "failed to load incidents", nil); return }
	} else { // in-memory fallback
		r.mu.RLock()
		deps = append([]metrics.Deployment(nil), r.deployments...)
		incs = append([]metrics.Incident(nil), r.incidents...)
		r.mu.RUnlock()
	}
	result, err := r.calculator.CalculateAll(deps, incs, tr)
	if err != nil { respondError(c, ErrInternal, "calculation failed", nil); return }
	classification := r.calculator.ClassifyPerformance(result)
	overall := r.calculator.GetOverallPerformance(classification)
	respondOK(c, gin.H{
		"deployment_frequency": result.DeploymentFrequency,
		"lead_time":            result.LeadTime.String(),
		"mttr":                 result.MTTR.String(),
		"change_failure_rate":  result.ChangeFailureRate,
		"classification":       classification,
		"overall_performance":  overall,
		"time_range":           gin.H{"start": tr.Start, "end": tr.End},
		"last_updated":         time.Now().UTC(),
		"data_quality":         result.DataQuality,
		"deployments_count":    len(deps),
		"incidents_count":      len(incs),
	})
}

func (r *Router) getDeploymentFrequency(c *gin.Context) {
	tr := r.parseTimeRange(c)
	var deps []metrics.Deployment
	if r.deploymentRepo != nil {
		var err error
		deps, err = r.deploymentRepo.ListRange(c.Request.Context(), tr.Start, tr.End)
		if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load deployments"}); return }
	} else {
		r.mu.RLock(); deps = append([]metrics.Deployment(nil), r.deployments...); r.mu.RUnlock()
	}
	freq := r.calculator.CalculateDeploymentFrequency(deps, tr)
	c.JSON(http.StatusOK, gin.H{"value": freq, "unit": "per_day", "time_range": tr})
}

func (r *Router) getLeadTime(c *gin.Context) {
	var deps []metrics.Deployment
	if r.deploymentRepo != nil {
		tr := r.parseTimeRange(c)
		var err error
		deps, err = r.deploymentRepo.ListRange(c.Request.Context(), tr.Start, tr.End)
		if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load deployments"}); return }
	} else { r.mu.RLock(); deps = append([]metrics.Deployment(nil), r.deployments...); r.mu.RUnlock() }
	lt := r.calculator.CalculateLeadTime(deps)
	c.JSON(http.StatusOK, gin.H{"value": lt.String(), "unit": "duration"})
}

func (r *Router) getMTTR(c *gin.Context) {
	var incs []metrics.Incident
	if r.incidentRepo != nil {
		tr := r.parseTimeRange(c)
		var err error
		incs, err = r.incidentRepo.ListRange(c.Request.Context(), tr.Start, tr.End)
		if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load incidents"}); return }
	} else { r.mu.RLock(); incs = append([]metrics.Incident(nil), r.incidents...); r.mu.RUnlock() }
	mttr := r.calculator.CalculateMTTR(incs)
	c.JSON(http.StatusOK, gin.H{"value": mttr.String(), "unit": "duration"})
}

func (r *Router) getChangeFailureRate(c *gin.Context) {
	tr := r.parseTimeRange(c)
	var deps []metrics.Deployment
	var incs []metrics.Incident
	if r.deploymentRepo != nil && r.incidentRepo != nil {
		var err error
		deps, err = r.deploymentRepo.ListRange(c.Request.Context(), tr.Start, tr.End)
		if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load deployments"}); return }
		incs, err = r.incidentRepo.ListRange(c.Request.Context(), tr.Start, tr.End)
		if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load incidents"}); return }
	} else { r.mu.RLock(); deps = append([]metrics.Deployment(nil), r.deployments...); incs = append([]metrics.Incident(nil), r.incidents...); r.mu.RUnlock() }
	cfr := r.calculator.CalculateChangeFailureRate(deps, incs)
	c.JSON(http.StatusOK, gin.H{"value": cfr, "unit": "ratio"})
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
			"error_rate":          0.02,
			"response_time":       120,
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
		date := now.AddDate(0, 0, -(29 - i))
		var value interface{}

		switch metricType {
		case "deployment_frequency":
			value = 1 + (i % 5) // 1-5 deployments per day
		case "lead_time":
			value = 2 + (i % 12) // 2-14 hours
		case "mttr":
			value = 1 + (i % 4) // 1-5 hours
		case "change_failure_rate":
			value = 5 + (i % 10) // 5-15%
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
	ID          string     `json:"id"`
	StartedAt   *time.Time `json:"started_at"`
	EndedAt     *time.Time `json:"ended_at"`
	Status      string     `json:"status"`
	CommitSHA   string     `json:"commit_sha"`
	Service     string     `json:"service"`
	Environment string     `json:"environment"`
}

func (r *Router) createDeployment(c *gin.Context) {
	var req deploymentRequest
	if err := c.ShouldBindJSON(&req); err != nil { respondError(c, ErrValidation, "invalid json", nil); return }
	start := time.Now(); if req.StartedAt != nil { start = *req.StartedAt }
	dep := metrics.Deployment{
		ID:          req.ID,
		Service:     req.Service,
		Environment: req.Environment,
		StartTime:   start,
		EndTime:     req.EndedAt,
		Status:      metrics.DeploymentStatus(req.Status),
		CommitSHA:   req.CommitSHA,
		CommitTime:  start,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if r.deploymentRepo != nil {
		if err := r.deploymentRepo.Create(c.Request.Context(), &dep); err != nil { respondError(c, ErrInternal, "failed to persist deployment", nil); return }
	} else { r.mu.Lock(); r.deployments = append(r.deployments, dep); r.mu.Unlock() }
	c.JSON(http.StatusCreated, gin.H{"data": gin.H{"deployment": dep}, "trace_id": requestIDFromContext(c)})
}

type incidentRequest struct {
	ID          string     `json:"id"`
	StartedAt   *time.Time `json:"started_at"`
	ResolvedAt  *time.Time `json:"resolved_at"`
	Severity    string     `json:"severity"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Service     string     `json:"service"`
	Environment string     `json:"environment"`
}

func (r *Router) createIncident(c *gin.Context) {
	var req incidentRequest
	if err := c.ShouldBindJSON(&req); err != nil { respondError(c, ErrValidation, "invalid json", nil); return }
	start := time.Now(); if req.StartedAt != nil { start = *req.StartedAt }
	inc := metrics.Incident{
		ID:           req.ID,
		Title:        req.Title,
		Description:  req.Description,
		Service:      req.Service,
		Environment:  req.Environment,
		Severity:     metrics.IncidentSeverity(req.Severity),
		StartTime:    start,
		ResolvedTime: req.ResolvedAt,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if r.incidentRepo != nil {
		if err := r.incidentRepo.Create(c.Request.Context(), &inc); err != nil { respondError(c, ErrInternal, "failed to persist incident", nil); return }
	} else { r.mu.Lock(); r.incidents = append(r.incidents, inc); r.mu.Unlock() }
	c.JSON(http.StatusCreated, gin.H{"data": gin.H{"incident": inc}, "trace_id": requestIDFromContext(c)})
}

func (r *Router) resolveIncident(c *gin.Context) {
	id := c.Param("id")
	now := time.Now()
	if r.incidentRepo != nil {
		if err := r.incidentRepo.Resolve(c.Request.Context(), id, now); err != nil { respondError(c, ErrNotFound, "incident not found", nil); return }
		respondOK(c, gin.H{"resolved_at": now})
		return
	}
	for i := range r.incidents {
		if r.incidents[i].ID == id {
			if r.incidents[i].ResolvedTime != nil { respondError(c, ErrConflict, "incident already resolved", nil); return }
			r.mu.Lock(); r.incidents[i].ResolvedTime = &now; r.incidents[i].UpdatedAt = time.Now(); r.mu.Unlock()
			respondOK(c, gin.H{"resolved_at": now}); return
		}
	}
	respondError(c, ErrNotFound, "incident not found", nil)
}

func (r *Router) listState(c *gin.Context) {
	tr := r.parseTimeRange(c)
	if r.deploymentRepo != nil && r.incidentRepo != nil {
		deps, err := r.deploymentRepo.ListRange(c.Request.Context(), tr.Start, tr.End); if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load deployments"}); return }
		incs, err := r.incidentRepo.ListRange(c.Request.Context(), tr.Start, tr.End); if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load incidents"}); return }
		sort.Slice(deps, func(i, j int) bool { return deps[i].StartTime.Before(deps[j].StartTime) })
		sort.Slice(incs, func(i, j int) bool { return incs[i].StartTime.Before(incs[j].StartTime) })
		c.JSON(http.StatusOK, gin.H{"deployments": deps, "incidents": incs}); return
	}
	r.mu.RLock(); deps := append([]metrics.Deployment(nil), r.deployments...); incs := append([]metrics.Incident(nil), r.incidents...); r.mu.RUnlock()
	sort.Slice(deps, func(i, j int) bool { return deps[i].StartTime.Before(deps[j].StartTime) })
	sort.Slice(incs, func(i, j int) bool { return incs[i].StartTime.Before(incs[j].StartTime) })
	c.JSON(http.StatusOK, gin.H{"deployments": deps, "incidents": incs})
}

// listDeployments returns deployments only
func (r *Router) listDeployments(c *gin.Context) {
	tr := r.parseTimeRange(c)
	if r.deploymentRepo != nil { deps, err := r.deploymentRepo.ListRange(c.Request.Context(), tr.Start, tr.End); if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load deployments"}); return }; sort.Slice(deps, func(i, j int) bool { return deps[i].StartTime.Before(deps[j].StartTime) }); c.JSON(http.StatusOK, gin.H{"deployments": deps, "count": len(deps)}); return }
	r.mu.RLock(); deps := append([]metrics.Deployment(nil), r.deployments...); r.mu.RUnlock(); sort.Slice(deps, func(i, j int) bool { return deps[i].StartTime.Before(deps[j].StartTime) }); c.JSON(http.StatusOK, gin.H{"deployments": deps, "count": len(deps)})
}

// listIncidents returns incidents only
func (r *Router) listIncidents(c *gin.Context) {
	tr := r.parseTimeRange(c)
	if r.incidentRepo != nil { incs, err := r.incidentRepo.ListRange(c.Request.Context(), tr.Start, tr.End); if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load incidents"}); return }; sort.Slice(incs, func(i, j int) bool { return incs[i].StartTime.Before(incs[j].StartTime) }); c.JSON(http.StatusOK, gin.H{"incidents": incs, "count": len(incs)}); return }
	r.mu.RLock(); incs := append([]metrics.Incident(nil), r.incidents...); r.mu.RUnlock(); sort.Slice(incs, func(i, j int) bool { return incs[i].StartTime.Before(incs[j].StartTime) }); c.JSON(http.StatusOK, gin.H{"incidents": incs, "count": len(incs)})
}
