package lifecycle

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
)

// Manager coordinates the start/stop lifecycle for the registered services.
type Manager struct {
	registry *Registry

	mu      sync.Mutex
	started []Service
}

// NewManager constructs a Manager bound to the provided registry.
func NewManager(reg *Registry) *Manager {
	return &Manager{registry: reg}
}

// Start iterates through registered services in registration order.
// If any service fails to start the already started ones are stopped in reverse
// order and the error is returned.
func (m *Manager) Start(ctx context.Context) error {
	services := m.registry.Services()
	started := make([]Service, 0, len(services))

	for _, svc := range services {
		if err := svc.Start(ctx); err != nil {
			m.stopReverse(ctx, started)
			return fmt.Errorf("start %s: %w", svc.Name(), err)
		}
		slog.Info("service started", "service", svc.Name())
		started = append(started, svc)
	}

	m.mu.Lock()
	m.started = started
	m.mu.Unlock()

	return nil
}

// Stop stops the started services in reverse order. All stop errors are joined.
func (m *Manager) Stop(ctx context.Context) error {
	m.mu.Lock()
	started := m.started
	m.started = nil
	m.mu.Unlock()

	return m.stopReverse(ctx, started)
}

func (m *Manager) stopReverse(ctx context.Context, services []Service) error {
	var errs []error
	for i := len(services) - 1; i >= 0; i-- {
		svc := services[i]
		if err := svc.Stop(ctx); err != nil {
			errs = append(errs, fmt.Errorf("stop %s: %w", svc.Name(), err))
			slog.Error("service stop failed", "service", svc.Name(), "error", err)
		} else {
			slog.Info("service stopped", "service", svc.Name())
		}
	}
	return errors.Join(errs...)
}
