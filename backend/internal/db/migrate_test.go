package db

import (
	"testing"

	"gamelink/internal/config"
)

// Test that runDataFixups converts legacy 'cancelled' to 'canceled'.
func TestRunDataFixups_OrderStatusSpelling(t *testing.T) {
	cfg := config.AppConfig{
		Database: config.DatabaseConfig{Type: "sqlite", DSN: "file::memory:?cache=shared"},
	}
	db1, err := Open(cfg)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	// insert legacy value
	if err := db1.Exec("INSERT INTO orders (id, user_id, game_id, title, status, price_cents) VALUES (1,1,1,'t','cancelled',0)").Error; err != nil {
		t.Fatalf("insert legacy: %v", err)
	}
	// open again to trigger fixups on same shared memory database
	db2, err := Open(cfg)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	var got string
	if err := db2.Raw("SELECT status FROM orders WHERE id=1").Scan(&got).Error; err != nil {
		t.Fatalf("select: %v", err)
	}
	if got != "canceled" {
		t.Fatalf("expected 'canceled', got %q", got)
	}
}
