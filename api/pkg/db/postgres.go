package db

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func openPostgres(cfg ConfigProvider, metrics MetricsProvider) (*gorm.DB, error) {
	totalStart := time.Now()

	// ── 阶段 1: 建立连接 ──
	t := time.Now()
	gormDB, err := gorm.Open(postgres.Open(cfg.GetDatabaseDSN()), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Warn),
		PrepareStmt:            true, // 预编译语句缓存，提升性能
		SkipDefaultTransaction: true, // 跳过默认事务，提升单条查询性能
	})
	if err != nil {
		return nil, fmt.Errorf("打开 postgres 失败: %w", err)
	}
	log.Printf("[startup] db connect: %v", time.Since(t))

	// ── 阶段 2: 连接池 ──
	t = time.Now()
	maxConns := cfg.GetMaxOpenConns()
	if maxConns <= 0 {
		maxConns = 100
	}
	if err := configureConnection(gormDB, maxConns); err != nil {
		return nil, err
	}
	log.Printf("[startup] db pool (max=%d): %v", maxConns, time.Since(t))

	// ── 阶段 3: AutoMigrate ──
	t = time.Now()
	if err := autoMigrate(gormDB); err != nil {
		return nil, err
	}
	log.Printf("[startup] auto-migrate: %v", time.Since(t))

	// ── 阶段 4: Data fixups ──
	t = time.Now()
	if err := runDataFixups(gormDB); err != nil {
		return nil, err
	}
	log.Printf("[startup] data fixups: %v", time.Since(t))

	// ── 阶段 5: Indexes ──
	t = time.Now()
	if err := ensureIndexes(gormDB); err != nil {
		return nil, err
	}
	log.Printf("[startup] indexes: %v", time.Since(t))

	// ── 阶段 6: Seed data ──
	if cfg.IsSeedEnabled() {
		t = time.Now()
		if err := applySeeds(gormDB); err != nil {
			return nil, err
		}
		log.Printf("[startup] seed data: %v", time.Since(t))
	} else {
		log.Printf("[startup] seed data: skipped (disabled)")
	}

	// ── 阶段 7: Metrics ──
	if metrics != nil {
		_ = metrics.InstrumentGorm(gormDB)
	}

	log.Printf("[startup] db total: %v", time.Since(totalStart))
	return gormDB, nil
}
