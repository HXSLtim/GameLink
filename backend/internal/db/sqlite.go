package db

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gamelink/internal/config"
	"gamelink/internal/metrics"
	
	// 使用纯 Go 实现的 SQLite GORM 驱动（无需 CGO）
	// github.com/glebarez/sqlite 基于 modernc.org/sqlite
	sqlite "github.com/glebarez/sqlite"
)

func openSQLite(cfg config.AppConfig) (*gorm.DB, error) {
	dsn := cfg.Database.DSN
	if err := ensureSQLiteDir(dsn); err != nil {
		return nil, err
	}

	gormDB, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("打开 sqlite 失败: %w", err)
	}

	if err := configureConnection(gormDB, 1); err != nil {
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

	if cfg.Seed.Enabled {
		if err := applySeeds(gormDB); err != nil {
			return nil, err
		}
	}

	_ = metrics.InstrumentGorm(gormDB)

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
