package lifecycle

import "context"

// Service describes a component that has an explicit start/stop lifecycle.
type Service interface {
	Name() string
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// Hook is a lightweight Service implementation backed by simple callbacks.
type Hook struct {
	name  string
	start func(context.Context) error
	stop  func(context.Context) error
}

// NewHook wraps the provided callbacks into a Service implementation.
func NewHook(name string, start func(context.Context) error, stop func(context.Context) error) *Hook {
	return &Hook{
		name:  name,
		start: start,
		stop:  stop,
	}
}

// Name implements Service.
func (h *Hook) Name() string { return h.name }

// Start implements Service.
func (h *Hook) Start(ctx context.Context) error {
	if h.start == nil {
		return nil
	}
	return h.start(ctx)
}

// Stop implements Service.
func (h *Hook) Stop(ctx context.Context) error {
	if h.stop == nil {
		return nil
	}
	return h.stop(ctx)
}
