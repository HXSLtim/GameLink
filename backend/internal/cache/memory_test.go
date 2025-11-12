package cache

import (
    "context"
    "testing"
    "time"
)

func TestMemoryCache_SetGetDeleteAndExpiry(t *testing.T) {
    c := NewMemory()
    ctx := context.Background()
    if _, ok, _ := c.Get(ctx, "k"); ok { t.Fatalf("expected miss") }
    if err := c.Set(ctx, "k", "v", 10*time.Millisecond); err != nil { t.Fatalf("set: %v", err) }
    v, ok, _ := c.Get(ctx, "k")
    if !ok || v != "v" { t.Fatalf("get: %v", v) }
    time.Sleep(20 * time.Millisecond)
    if _, ok, _ := c.Get(ctx, "k"); ok { t.Fatalf("expected expired miss") }
    if err := c.Set(ctx, "k2", "v2", 0); err != nil { t.Fatalf("set2: %v", err) }
    if err := c.Delete(ctx, "k2"); err != nil { t.Fatalf("del: %v", err) }
    if _, ok, _ := c.Get(ctx, "k2"); ok { t.Fatalf("expected deleted miss") }
    if err := c.Close(ctx); err != nil { t.Fatalf("close: %v", err) }
}

