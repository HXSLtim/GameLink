package db

import (
	"testing"

	"gamelink/internal/config"
	"gamelink/internal/model"
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
	var adminID uint64
	if err := db1.Raw("SELECT id FROM users WHERE role = ? LIMIT 1", model.RoleAdmin).Scan(&adminID).Error; err != nil {
		t.Fatalf("lookup admin: %v", err)
	}
	if adminID == 0 {
		t.Fatalf("expected super admin to be present")
	}
	if err := db1.Exec("INSERT INTO games (id, key, name) VALUES (1,'g','game')").Error; err != nil {
		t.Fatalf("seed game: %v", err)
	}
	// Create a service item first (needed for foreign key constraint)
	if err := db1.Exec("INSERT INTO service_items (id, item_code, name, sub_category, base_price_cents, service_hours) VALUES (1, 'TEST', 'Test Item', 'solo', 10000, 1)").Error; err != nil {
		t.Fatalf("seed service_item: %v", err)
	}
	// insert legacy value
	if err := db1.Exec("INSERT INTO orders (id, order_no, user_id, item_id, game_id, title, status, unit_price_cents, total_price_cents) VALUES (1, 'TEST001', ?, 1, 1, 't', 'cancelled', 0, 0)", adminID).Error; err != nil {
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
