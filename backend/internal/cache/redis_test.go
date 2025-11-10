package cache

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"

	"gamelink/internal/config"
)

func TestNewRedisOperations(t *testing.T) {
	mr := miniredis.RunT(t)
	defer mr.Close()

	cache, err := NewRedis(config.RedisConfig{
		Addr: mr.Addr(),
	})
	if err != nil {
		t.Fatalf("NewRedis failed: %v", err)
	}
	ctx := context.Background()

	if err := cache.Set(ctx, "foo", "bar", time.Minute); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	got, ok, err := cache.Get(ctx, "foo")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !ok || got != "bar" {
		t.Fatalf("expected hit with bar, got ok=%v value=%q", ok, got)
	}

	if err := cache.Delete(ctx, "foo"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, ok, err = cache.Get(ctx, "foo")
	if err != nil {
		t.Fatalf("Get after delete failed: %v", err)
	}
	if ok {
		t.Fatalf("expected cache miss after delete")
	}

	if err := cache.Close(ctx); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}

func TestNewRedisPingFailure(t *testing.T) {
	_, err := NewRedis(config.RedisConfig{
		Addr: "127.0.0.1:0",
	})
	if err == nil {
		t.Fatalf("expected error when redis is unreachable")
	}
}
