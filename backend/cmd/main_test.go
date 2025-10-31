package main

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestResolveGinMode(t *testing.T) {
	t.Setenv("GIN_MODE", "")
	t.Setenv("APP_ENV", "")
	if got := resolveGinMode(); got != gin.DebugMode {
		t.Fatalf("expected debug mode, got %s", got)
	}

	t.Setenv("APP_ENV", "production")
	if got := resolveGinMode(); got != gin.ReleaseMode {
		t.Fatalf("expected release mode, got %s", got)
	}

	t.Setenv("GIN_MODE", "test")
	if got := resolveGinMode(); got != "test" {
		t.Fatalf("expected env override to 'test', got %s", got)
	}
}
