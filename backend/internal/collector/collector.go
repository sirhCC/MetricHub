package collector

// Package collector will coordinate metric collection cycles across plugins.
// Phase 1 placeholder: defines basic interface and noop implementation.

import "context"

// CycleCollector defines the minimal interface for a metric collection cycle.
type CycleCollector interface {
	Collect(ctx context.Context) error
}

// NoopCollector is a stub used until real implementations arrive.
type NoopCollector struct{}

func (n *NoopCollector) Collect(ctx context.Context) error { return nil }
