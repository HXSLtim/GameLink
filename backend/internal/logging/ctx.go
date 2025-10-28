package logging

import "context"

type ctxKey string

const (
    keyRequestID ctxKey = "request_id"
    keyActorUserID ctxKey = "actor_user_id"
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

// WithActorUserID stores current actor (user id) into context.
func WithActorUserID(ctx context.Context, id uint64) context.Context {
    return context.WithValue(ctx, keyActorUserID, id)
}

// ActorUserIDFromContext returns the user id if present.
func ActorUserIDFromContext(ctx context.Context) (uint64, bool) {
    v := ctx.Value(keyActorUserID)
    if v == nil { return 0, false }
    if id, ok := v.(uint64); ok { return id, true }
    return 0, false
}
