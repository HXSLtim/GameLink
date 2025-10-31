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

func TestMemoryCacheDelete(t *testing.T) {
	ctx := context.Background()
	c := NewMemory()
	defer func() { _ = c.Close(ctx) }()

	// Set a value
	if err := c.Set(ctx, "key", "value", 0); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	// Verify it exists
	v, ok, err := c.Get(ctx, "key")
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if !ok || v != "value" {
		t.Fatalf("unexpected value: ok=%v v=%s", ok, v)
	}

	// Delete it
	if err := c.Delete(ctx, "key"); err != nil {
		t.Fatalf("Delete error: %v", err)
	}

	// Verify it's gone
	_, ok, err = c.Get(ctx, "key")
	if err != nil {
		t.Fatalf("Get error after delete: %v", err)
	}
	if ok {
		t.Fatal("key should be deleted")
	}
}

func TestMemoryCacheGetNonExistent(t *testing.T) {
	ctx := context.Background()
	c := NewMemory()
	defer func() { _ = c.Close(ctx) }()

	v, ok, err := c.Get(ctx, "non_existent_key")
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if ok {
		t.Fatalf("Expected key not found, but got value: %s", v)
	}
}

func TestMemoryCacheJanitor(t *testing.T) {
	ctx := context.Background()
	c := NewMemory()
	defer func() { _ = c.Close(ctx) }()

	// Set multiple keys with short TTL
	for i := 0; i < 5; i++ {
		key := string(rune('a' + i))
		if err := c.Set(ctx, key, "value", 50*time.Millisecond); err != nil {
			t.Fatalf("Set error: %v", err)
		}
	}

	// Wait for janitor to clean up
	time.Sleep(200 * time.Millisecond)

	// Verify all keys are expired
	for i := 0; i < 5; i++ {
		key := string(rune('a' + i))
		_, ok, _ := c.Get(ctx, key)
		if ok {
			t.Fatalf("key %s should have been cleaned up by janitor", key)
		}
	}
}

func TestMemoryCacheClose(t *testing.T) {
	ctx := context.Background()
	c := NewMemory()

	// Set some values
	c.Set(ctx, "key1", "value1", 0)
	c.Set(ctx, "key2", "value2", 0)

	// Close the cache
	if err := c.Close(ctx); err != nil {
		t.Fatalf("Close error: %v", err)
	}

	// Verify values are cleared
	_, ok, _ := c.Get(ctx, "key1")
	if ok {
		t.Fatal("Cache should be cleared after close")
	}
}
