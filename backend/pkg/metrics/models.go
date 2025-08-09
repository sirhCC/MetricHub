package metrics

import (
	"time"
)

// DORA Metrics Domain Models
// Following Clean Architecture principles with rich domain models

// DeploymentStatus represents the status of a deployment
type DeploymentStatus string

const (
	DeploymentStatusPending   DeploymentStatus = "pending"
	DeploymentStatusRunning   DeploymentStatus = "running"
	DeploymentStatusSuccess   DeploymentStatus = "success"
	DeploymentStatusFailed    DeploymentStatus = "failed"
	DeploymentStatusCancelled DeploymentStatus = "cancelled"
)

// IncidentSeverity represents the severity of an incident
type IncidentSeverity string

const (
	SeverityLow      IncidentSeverity = "low"
	SeverityMedium   IncidentSeverity = "medium"
	SeverityHigh     IncidentSeverity = "high"
	SeverityCritical IncidentSeverity = "critical"
)

// Deployment represents a deployment event in the system
type Deployment struct {
	ID          string            `json:"id"`
	Service     string            `json:"service"`
	Environment string            `json:"environment"`
	Version     string            `json:"version"`
	Status      DeploymentStatus  `json:"status"`
	StartTime   time.Time         `json:"start_time"`
	EndTime     *time.Time        `json:"end_time,omitempty"`
	CommitSHA   string            `json:"commit_sha"`
	CommitTime  time.Time         `json:"commit_time"`
	Author      string            `json:"author"`
	Repository  string            `json:"repository"`
	Branch      string            `json:"branch"`
	BuildURL    string            `json:"build_url,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// LeadTime calculates the lead time for this deployment
func (d *Deployment) LeadTime() time.Duration {
	if d.EndTime == nil {
		return 0
	}
	return d.EndTime.Sub(d.CommitTime)
}

// Duration calculates the deployment duration
func (d *Deployment) Duration() time.Duration {
	if d.EndTime == nil {
		return 0
	}
	return d.EndTime.Sub(d.StartTime)
}

// IsSuccessful returns true if the deployment was successful
func (d *Deployment) IsSuccessful() bool {
	return d.Status == DeploymentStatusSuccess
}

// IsFailed returns true if the deployment failed
func (d *Deployment) IsFailed() bool {
	return d.Status == DeploymentStatusFailed
}

// Incident represents an incident/outage in the system
type Incident struct {
	ID           string            `json:"id"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	Service      string            `json:"service"`
	Environment  string            `json:"environment"`
	Severity     IncidentSeverity  `json:"severity"`
	StartTime    time.Time         `json:"start_time"`
	ResolvedTime *time.Time        `json:"resolved_time,omitempty"`
	RootCause    string            `json:"root_cause,omitempty"`
	Assignee     string            `json:"assignee,omitempty"`
	Tags         map[string]string `json:"tags,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// MTTR calculates the mean time to recovery for this incident
func (i *Incident) MTTR() time.Duration {
	if i.ResolvedTime == nil {
		return 0
	}
	return i.ResolvedTime.Sub(i.StartTime)
}

// IsResolved returns true if the incident is resolved
func (i *Incident) IsResolved() bool {
	return i.ResolvedTime != nil
}

// TimeRange represents a time range for metrics calculation
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// Duration returns the duration of the time range
func (tr *TimeRange) Duration() time.Duration {
	return tr.End.Sub(tr.Start)
}

// Days returns the number of days in the time range
func (tr *TimeRange) Days() float64 {
	return tr.Duration().Hours() / 24
}

// Contains checks if a time is within the range
func (tr *TimeRange) Contains(t time.Time) bool {
	return t.After(tr.Start) && t.Before(tr.End)
}

// DORAMetrics represents the calculated DORA metrics
type DORAMetrics struct {
	DeploymentFrequency float64       `json:"deployment_frequency"` // Deployments per day
	LeadTime            time.Duration `json:"lead_time"`            // Average lead time
	MTTR                time.Duration `json:"mttr"`                 // Mean time to recovery
	ChangeFailureRate   float64       `json:"change_failure_rate"`  // Percentage as decimal (0.15 = 15%)
	TimeRange           TimeRange     `json:"time_range"`
	CalculatedAt        time.Time     `json:"calculated_at"`
	DataQuality         string        `json:"data_quality"` // high, medium, low
}

// MetricsFilter represents filters for metrics queries
type MetricsFilter struct {
	TimeRange    TimeRange         `json:"time_range"`
	Services     []string          `json:"services,omitempty"`
	Environments []string          `json:"environments,omitempty"`
	Tags         map[string]string `json:"tags,omitempty"`
}

// Predefined time ranges for convenience
func Last7Days() TimeRange {
	now := time.Now()
	return TimeRange{
		Start: now.AddDate(0, 0, -7),
		End:   now,
	}
}

func Last30Days() TimeRange {
	now := time.Now()
	return TimeRange{
		Start: now.AddDate(0, 0, -30),
		End:   now,
	}
}

func Last90Days() TimeRange {
	now := time.Now()
	return TimeRange{
		Start: now.AddDate(0, 0, -90),
		End:   now,
	}
}

func ThisMonth() TimeRange {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return TimeRange{
		Start: start,
		End:   now,
	}
}

func LastMonth() TimeRange {
	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)
	start := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, lastMonth.Location())
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return TimeRange{
		Start: start,
		End:   end,
	}
}
