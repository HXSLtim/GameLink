package db

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	// SQLite GORM driver implemented in pure Go (no CGO required)
	// github.com/glebarez/sqlite is based on modernc.org/sqlite
	sqlite "github.com/glebarez/sqlite"
)

func openSQLite(cfg ConfigProvider, metrics MetricsProvider) (*gorm.DB, error) {
	dsn := cfg.GetDatabaseDSN()
	if err := ensureSQLiteDir(dsn); err != nil {
		return nil, err
	}

	gormDB, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("打开 sqlite 失败: %w", err)
	}

	if err := configureConnection(gormDB, cfg.GetMaxOpenConns()); err != nil {
		return nil, err
	}

	if err := gormDB.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return nil, fmt.Errorf("启用 sqlite 外键失败: %w", err)
	}

	if err := autoMigrate(gormDB); err != nil {
		return nil, err
	}

	if err := runDataFixups(gormDB); err != nil {
		return nil, err
	}

	if err := ensureIndexes(gormDB); err != nil {
		return nil, err
	}

	if cfg.IsSeedEnabled() {
		if err := applySeeds(gormDB); err != nil {
			return nil, err
		}
	}

	if metrics != nil {
		_ = metrics.InstrumentGorm(gormDB)
	}

	return gormDB, nil
}

func configureConnection(db *gorm.DB, maxOpen int) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxOpen)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	return nil
}

func ensureSQLiteDir(dsn string) error {
	if !strings.HasPrefix(dsn, "file:") {
		return nil
	}

	path := strings.TrimPrefix(dsn, "file:")
	if idx := strings.Index(path, "?"); idx >= 0 {
		path = path[:idx]
	}

	dir := filepath.Dir(path)
	if dir == "." {
		return nil
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("创建 sqlite 目录失败: %w", err)
	}
	return nil
}
