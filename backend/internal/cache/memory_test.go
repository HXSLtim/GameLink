package cache

import (
	"context"
	"testing"
	"time"
)

func TestMemoryCacheSetGet(t *testing.T) {
	ctx := context.Background()
	c := NewMemory()
	defer func() { _ = c.Close(ctx) }()

	if err := c.Set(ctx, "key", "value", 0); err != nil {
		t.Fatalf("Set error: %v", err)
	}
	v, ok, err := c.Get(ctx, "key")
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if !ok || v != "value" {
		t.Fatalf("unexpected value: ok=%v v=%s", ok, v)
	}

	if err := c.Set(ctx, "expire", "soon", 10*time.Millisecond); err != nil {
		t.Fatalf("Set with ttl error: %v", err)
	}
	time.Sleep(20 * time.Millisecond)
	if _, ok, _ := c.Get(ctx, "expire"); ok {
		t.Fatal("value should have expired")
	}
}
