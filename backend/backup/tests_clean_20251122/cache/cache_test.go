package cache

import (
	"testing"

	"gamelink/internal/config"
)

func TestNew(t *testing.T) {
	t.Run("Default to memory", func(t *testing.T) {
		c, err := New(config.CacheConfig{})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if c == nil {
			t.Fatal("Expected non-nil cache")
		}
		defer c.Close(nil)
	})

	t.Run("Memory type", func(t *testing.T) {
		c, err := New(config.CacheConfig{Type: "memory"})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if c == nil {
			t.Fatal("Expected non-nil cache")
		}
		defer c.Close(nil)
	})

	t.Run("Unknown type defaults to memory", func(t *testing.T) {
		c, err := New(config.CacheConfig{Type: "unknown"})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if c == nil {
			t.Fatal("Expected non-nil cache")
		}
		defer c.Close(nil)
	})

	t.Run("Redis type with invalid config", func(t *testing.T) {
		// This should fail because Redis is not running
		c, err := New(config.CacheConfig{
			Type: "redis",
			Redis: config.RedisConfig{
				Addr: "localhost:9999", // Non-existent Redis
			},
		})
		if err == nil {
			t.Error("Expected error for invalid Redis config")
			if c != nil {
				c.Close(nil)
			}
		}
	})
}
