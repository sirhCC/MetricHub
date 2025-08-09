package plugins

// Package plugins will host the plugin manager and plugin interface definitions.
// Phase 1 placeholder.

import "context"

// Plugin defines the minimal contract all metric source plugins must satisfy.
type Plugin interface {
	Name() string
	Description() string
	Initialize(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

// Manager orchestrates registered plugins.
type Manager struct{ plugins []Plugin }

// Register adds a plugin to the manager.
func (m *Manager) Register(p Plugin) { m.plugins = append(m.plugins, p) }

// InitializeAll runs Initialize on all registered plugins.
func (m *Manager) InitializeAll(ctx context.Context) error {
	for _, p := range m.plugins {
		if err := p.Initialize(ctx); err != nil {
			return err
		}
	}
	return nil
}

// ShutdownAll runs Shutdown on all registered plugins (best-effort).
func (m *Manager) ShutdownAll(ctx context.Context) error {
	var firstErr error
	for _, p := range m.plugins {
		if err := p.Shutdown(ctx); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
