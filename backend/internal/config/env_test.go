package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("APP_ENV", "unit-test")
	cfg := Load()
	if cfg.Port == "" {
		t.Fatal("expected default port to be set")
	}
	if cfg.Database.Type == "" {
		t.Fatal("expected database type to be set")
	}
}

func TestValidateProductionRequiresDSN(t *testing.T) {
	cfg := AppConfig{}
	if err := Validate("production", cfg); err == nil {
		t.Fatal("expected validation error when DSN missing")
	}
	cfg.Database.DSN = "postgres://example"
	if err := Validate("production", cfg); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}
