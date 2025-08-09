package metrics

import (
	"testing"
	"time"
)

// helper to create pointer time
func ptrTime(t time.Time) *time.Time { return &t }

func TestCalculateDeploymentFrequency(t *testing.T) {
	calc := NewDORACalculator()
	tr := TimeRange{Start: time.Now().Add(-24 * time.Hour), End: time.Now()}
	end := time.Now()
	deployments := []Deployment{
		{Status: DeploymentStatusSuccess, StartTime: time.Now().Add(-23 * time.Hour), EndTime: &end},
		{Status: DeploymentStatusFailed, StartTime: time.Now().Add(-20 * time.Hour), EndTime: &end},
	}
	freq := calc.CalculateDeploymentFrequency(deployments, tr)
	if freq <= 0 || freq > 2 { // simple sanity check: 2 deployments in 1 day
		t.Fatalf("unexpected deployment frequency: %v", freq)
	}
}

func TestCalculateLeadTime(t *testing.T) {
	calc := NewDORACalculator()
	now := time.Now()
	successfulEnd := now
	deployments := []Deployment{
		{Status: DeploymentStatusSuccess, CommitTime: now.Add(-2 * time.Hour), EndTime: &successfulEnd},
		{Status: DeploymentStatusFailed, CommitTime: now.Add(-3 * time.Hour), EndTime: &successfulEnd},
	}
	lt := calc.CalculateLeadTime(deployments)
	if lt <= 0 || lt > 3*time.Hour {
		t.Fatalf("unexpected lead time: %v", lt)
	}
}

func TestCalculateMTTR(t *testing.T) {
	calc := NewDORACalculator()
	now := time.Now()
	incidents := []Incident{
		{StartTime: now.Add(-3 * time.Hour), ResolvedTime: ptrTime(now.Add(-2 * time.Hour))},
		{StartTime: now.Add(-6 * time.Hour), ResolvedTime: ptrTime(now.Add(-5 * time.Hour))},
	}
	mttr := calc.CalculateMTTR(incidents)
	if mttr <= 0 || mttr > 2*time.Hour {
		t.Fatalf("unexpected mttr: %v", mttr)
	}
}

func TestCalculateChangeFailureRate(t *testing.T) {
	calc := NewDORACalculator()
	now := time.Now()
	end1 := now.Add(-4 * time.Hour)
	end2 := now.Add(-3 * time.Hour)
	deployments := []Deployment{
		{Status: DeploymentStatusSuccess, EndTime: &end1, StartTime: end1.Add(-10 * time.Minute)},
		{Status: DeploymentStatusSuccess, EndTime: &end2, StartTime: end2.Add(-10 * time.Minute)},
	}
	incidents := []Incident{{StartTime: end2.Add(30 * time.Minute), ResolvedTime: ptrTime(end2.Add(90 * time.Minute))}}
	rate := calc.CalculateChangeFailureRate(deployments, incidents)
	if rate <= 0 || rate > 1 { // expect at least one failure correlation
		t.Fatalf("unexpected change failure rate: %v", rate)
	}
}
