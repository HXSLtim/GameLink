package db

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	sqlite "github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"gamelink/internal/config"
)

func TestOpenUnsupportedType(t *testing.T) {
	cfg := config.AppConfig{
		Database: config.DatabaseConfig{
			Type: "oracle",
			DSN:  "oracle://localhost",
		},
	}
	_, err := Open(cfg)
	if err == nil || !strings.Contains(err.Error(), "暂不支持的数据库类型") {
		t.Fatalf("expected unsupported type error, got %v", err)
	}
}

func TestEnsureSQLiteDirCreatesDirectory(t *testing.T) {
	temp := t.TempDir()
	target := filepath.Join(temp, "nested", "db.sqlite")
	dsn := "file:" + target + "?cache=shared"

	if err := ensureSQLiteDir(dsn); err != nil {
		t.Fatalf("ensureSQLiteDir failed: %v", err)
	}

	if _, err := os.Stat(filepath.Dir(target)); err != nil {
		t.Fatalf("expected directory to exist: %v", err)
	}
}

func TestEnsureSQLiteDirNoopForNonFileDSN(t *testing.T) {
	if err := ensureSQLiteDir("sqlite::memory:?cache=shared"); err != nil {
		t.Fatalf("expected no error for in-memory DSN, got %v", err)
	}
}

func TestConfigureConnectionSetsLimits(t *testing.T) {
	gdb, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite memory failed: %v", err)
	}

	if err := configureConnection(gdb, 2); err != nil {
		t.Fatalf("configureConnection failed: %v", err)
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		t.Fatalf("sql DB: %v", err)
	}
	if stats := sqlDB.Stats(); stats.MaxOpenConnections != 2 {
		t.Fatalf("expected MaxOpenConnections=2, got %d", stats.MaxOpenConnections)
	}
}
