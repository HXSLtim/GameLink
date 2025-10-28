package logging

import "context"

type ctxKey string

const (
	keyRequestID ctxKey = "request_id"
)

// WithRequestID stores request id into context.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyRequestID, id)
}

// RequestIDFromContext returns the request id if present.
func RequestIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(keyRequestID)
	if v == nil {
		return "", false
	}
	if s, ok := v.(string); ok {
		return s, true
	}
	return "", false
}
