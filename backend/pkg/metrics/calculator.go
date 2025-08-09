package metrics

import (
	"time"
)

// DORACalculator implements the DORA metrics calculation logic
// Following the official DORA research methodology
type DORACalculator struct{}

// NewDORACalculator creates a new DORA metrics calculator
func NewDORACalculator() *DORACalculator {
	return &DORACalculator{}
}

// CalculateAll calculates all DORA metrics for the given data
func (c *DORACalculator) CalculateAll(deployments []Deployment, incidents []Incident, timeRange TimeRange) (*DORAMetrics, error) {
	// Filter data to time range
	filteredDeployments := c.filterDeployments(deployments, timeRange)
	filteredIncidents := c.filterIncidents(incidents, timeRange)

	// Calculate individual metrics
	deploymentFreq := c.CalculateDeploymentFrequency(filteredDeployments, timeRange)
	leadTime := c.CalculateLeadTime(filteredDeployments)
	mttr := c.CalculateMTTR(filteredIncidents)
	changeFailureRate := c.CalculateChangeFailureRate(filteredDeployments, filteredIncidents)

	// Determine data quality
	dataQuality := c.assessDataQuality(filteredDeployments, filteredIncidents, timeRange)

	return &DORAMetrics{
		DeploymentFrequency: deploymentFreq,
		LeadTime:           leadTime,
		MTTR:               mttr,
		ChangeFailureRate:  changeFailureRate,
		TimeRange:          timeRange,
		CalculatedAt:       time.Now(),
		DataQuality:        dataQuality,
	}, nil
}

// CalculateDeploymentFrequency calculates how often deployments occur
// Returns deployments per day
func (c *DORACalculator) CalculateDeploymentFrequency(deployments []Deployment, timeRange TimeRange) float64 {
	successfulDeployments := 0
	for _, deployment := range deployments {
		if deployment.IsSuccessful() {
			successfulDeployments++
		}
	}

	days := timeRange.Days()
	if days == 0 {
		return 0
	}

	return float64(successfulDeployments) / days
}

// CalculateLeadTime calculates the average time from commit to production
func (c *DORACalculator) CalculateLeadTime(deployments []Deployment) time.Duration {
	if len(deployments) == 0 {
		return 0
	}

	var totalLeadTime time.Duration
	validDeployments := 0

	for _, deployment := range deployments {
		if deployment.IsSuccessful() && deployment.EndTime != nil {
			leadTime := deployment.LeadTime()
			if leadTime > 0 {
				totalLeadTime += leadTime
				validDeployments++
			}
		}
	}

	if validDeployments == 0 {
		return 0
	}

	return totalLeadTime / time.Duration(validDeployments)
}

// CalculateMTTR calculates the mean time to recovery from incidents
func (c *DORACalculator) CalculateMTTR(incidents []Incident) time.Duration {
	if len(incidents) == 0 {
		return 0
	}

	var totalRecoveryTime time.Duration
	resolvedIncidents := 0

	for _, incident := range incidents {
		if incident.IsResolved() {
			mttr := incident.MTTR()
			if mttr > 0 {
				totalRecoveryTime += mttr
				resolvedIncidents++
			}
		}
	}

	if resolvedIncidents == 0 {
		return 0
	}

	return totalRecoveryTime / time.Duration(resolvedIncidents)
}

// CalculateChangeFailureRate calculates the percentage of deployments that cause incidents
func (c *DORACalculator) CalculateChangeFailureRate(deployments []Deployment, incidents []Incident) float64 {
	if len(deployments) == 0 {
		return 0
	}

	totalDeployments := len(deployments)
	failedDeployments := 0

	// Count direct deployment failures
	for _, deployment := range deployments {
		if deployment.IsFailed() {
			failedDeployments++
		}
	}

	// Count deployments that caused incidents (simplified approach)
	// In a real implementation, you'd correlate incidents with deployments more precisely
	for _, incident := range incidents {
		// Find deployments that happened shortly before the incident
		for _, deployment := range deployments {
			if deployment.IsSuccessful() && 
			   deployment.EndTime != nil && 
			   incident.StartTime.Sub(*deployment.EndTime) > 0 && 
			   incident.StartTime.Sub(*deployment.EndTime) < 2*time.Hour {
				failedDeployments++
				break // Only count one deployment per incident
			}
		}
	}

	return float64(failedDeployments) / float64(totalDeployments)
}

// Helper methods

func (c *DORACalculator) filterDeployments(deployments []Deployment, timeRange TimeRange) []Deployment {
	var filtered []Deployment
	for _, deployment := range deployments {
		if timeRange.Contains(deployment.StartTime) {
			filtered = append(filtered, deployment)
		}
	}
	return filtered
}

func (c *DORACalculator) filterIncidents(incidents []Incident, timeRange TimeRange) []Incident {
	var filtered []Incident
	for _, incident := range incidents {
		if timeRange.Contains(incident.StartTime) {
			filtered = append(filtered, incident)
		}
	}
	return filtered
}

func (c *DORACalculator) assessDataQuality(deployments []Deployment, incidents []Incident, timeRange TimeRange) string {
	deploymentCount := len(deployments)
	incidentCount := len(incidents)
	days := timeRange.Days()

	// Simple heuristic for data quality assessment
	if days < 7 {
		return "low" // Not enough time range
	}

	deploymentsPerDay := float64(deploymentCount) / days

	// High quality: Regular deployments and incident data
	if deploymentsPerDay >= 0.5 && incidentCount > 0 {
		return "high"
	}

	// Medium quality: Some data but not comprehensive
	if deploymentCount > 0 || incidentCount > 0 {
		return "medium"
	}

	// Low quality: Insufficient data
	return "low"
}

// Performance classification based on DORA research
type PerformanceLevel string

const (
	Elite      PerformanceLevel = "elite"
	High       PerformanceLevel = "high"
	Medium     PerformanceLevel = "medium"
	Low        PerformanceLevel = "low"
)

// ClassifyPerformance classifies the team's performance based on DORA metrics
func (c *DORACalculator) ClassifyPerformance(metrics *DORAMetrics) map[string]PerformanceLevel {
	classification := make(map[string]PerformanceLevel)

	// Deployment Frequency classification
	switch {
	case metrics.DeploymentFrequency >= 1: // Multiple times per day
		classification["deployment_frequency"] = Elite
	case metrics.DeploymentFrequency >= 0.14: // Once per week
		classification["deployment_frequency"] = High
	case metrics.DeploymentFrequency >= 0.033: // Once per month
		classification["deployment_frequency"] = Medium
	default:
		classification["deployment_frequency"] = Low
	}

	// Lead Time classification
	leadTimeHours := metrics.LeadTime.Hours()
	switch {
	case leadTimeHours <= 24: // Less than 1 day
		classification["lead_time"] = Elite
	case leadTimeHours <= 168: // Less than 1 week
		classification["lead_time"] = High
	case leadTimeHours <= 720: // Less than 1 month
		classification["lead_time"] = Medium
	default:
		classification["lead_time"] = Low
	}

	// MTTR classification
	mttrHours := metrics.MTTR.Hours()
	switch {
	case mttrHours <= 1: // Less than 1 hour
		classification["mttr"] = Elite
	case mttrHours <= 24: // Less than 1 day
		classification["mttr"] = High
	case mttrHours <= 168: // Less than 1 week
		classification["mttr"] = Medium
	default:
		classification["mttr"] = Low
	}

	// Change Failure Rate classification
	switch {
	case metrics.ChangeFailureRate <= 0.15: // 0-15%
		classification["change_failure_rate"] = Elite
	case metrics.ChangeFailureRate <= 0.20: // 16-20%
		classification["change_failure_rate"] = High
	case metrics.ChangeFailureRate <= 0.30: // 21-30%
		classification["change_failure_rate"] = Medium
	default:
		classification["change_failure_rate"] = Low
	}

	return classification
}

// GetOverallPerformance determines the overall performance level
func (c *DORACalculator) GetOverallPerformance(classification map[string]PerformanceLevel) PerformanceLevel {
	levels := make(map[PerformanceLevel]int)
	
	for _, level := range classification {
		levels[level]++
	}

	// Majority rule with preference for higher performance
	if levels[Elite] >= 3 {
		return Elite
	}
	if levels[Elite] + levels[High] >= 3 {
		return High
	}
	if levels[Low] >= 3 {
		return Low
	}
	return Medium
}
