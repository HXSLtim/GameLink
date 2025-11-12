package logging

import (
	"context"
	"testing"
)

func TestRequestIDContextHelpers(t *testing.T) {
	base := context.Background()
	ctx := WithRequestID(base, "req-123")

	got, ok := RequestIDFromContext(ctx)
	if !ok || got != "req-123" {
		t.Fatalf("RequestIDFromContext()=%q,%v want req-123,true", got, ok)
	}

	_, ok = RequestIDFromContext(base)
	if ok {
		t.Fatal("expected missing request id to return ok=false")
	}

	// ensure type mismatch gracefully handled
	ctx = context.WithValue(base, keyRequestID, 123)
	if _, ok := RequestIDFromContext(ctx); ok {
		t.Fatal("expected invalid type to return ok=false")
	}
}

func TestActorUserIDContextHelpers(t *testing.T) {
	base := context.Background()
	ctx := WithActorUserID(base, 42)

	got, ok := ActorUserIDFromContext(ctx)
	if !ok || got != 42 {
		t.Fatalf("ActorUserIDFromContext()=%d,%v want 42,true", got, ok)
	}

	_, ok = ActorUserIDFromContext(base)
	if ok {
		t.Fatal("expected missing actor user id to return ok=false")
	}

	ctx = context.WithValue(base, keyActorUserID, "not-uint64")
	if _, ok := ActorUserIDFromContext(ctx); ok {
		t.Fatal("expected invalid type to return ok=false")
	}
}
