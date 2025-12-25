package lifecycle

import (
	"context"
	"sync"
)

// Registry collects lifecycle services that should be started and stopped with
// the application.
type Registry struct {
	mu       sync.Mutex
	services []Service
}

// NewRegistry creates an empty registry.
func NewRegistry() *Registry {
	return &Registry{services: make([]Service, 0)}
}

// Register adds a Service to the registry.
func (r *Registry) Register(s Service) {
	if s == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.services = append(r.services, s)
}

// RegisterHook is a helper that registers start/stop callbacks as a Service.
func (r *Registry) RegisterHook(name string, start func(context.Context) error, stop func(context.Context) error) {
	r.Register(NewHook(name, start, stop))
}

// Services returns a copy of the registered services in registration order.
func (r *Registry) Services() []Service {
	r.mu.Lock()
	defer r.mu.Unlock()
	dst := make([]Service, len(r.services))
	copy(dst, r.services)
	return dst
}
