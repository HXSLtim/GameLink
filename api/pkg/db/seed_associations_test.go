package db

import (
	"testing"

	"gamelink/pkg/testutil"
)

func TestSeedAssociations(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	t.Cleanup(func() { testutil.CleanDB(t, db) })

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get db instance: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)

	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		t.Fatalf("failed to enable sqlite foreign keys: %v", err)
	}

	if err := autoMigrate(db); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	if err := applySeeds(db); err != nil {
		t.Fatalf("seed failed: %v", err)
	}
}
