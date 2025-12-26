// Package testutil provides utilities for testing.
package testutil

import (
	"testing"

	sqlite "github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// NewMemoryDB creates an in-memory SQLite database for testing.
func NewMemoryDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to create memory db: %v", err)
	}

	return db
}

// MigrateTables auto migrates table structures for testing.
func MigrateTables(t *testing.T, db *gorm.DB, models ...interface{}) {
	t.Helper()

	if err := db.AutoMigrate(models...); err != nil {
		t.Fatalf("failed to migrate tables: %v", err)
	}
}

// CleanDB cleans up the database after testing.
func CleanDB(t *testing.T, db *gorm.DB) {
	t.Helper()

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get db instance: %v", err)
	}

	if err := sqlDB.Close(); err != nil {
		t.Fatalf("failed to close db: %v", err)
	}
}
