package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func openPostgres(cfg ConfigProvider, metrics MetricsProvider) (*gorm.DB, error) {
	gormDB, err := gorm.Open(postgres.Open(cfg.GetDatabaseDSN()), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Warn),
		PrepareStmt:            true, // 预编译语句缓存，提升性能
		SkipDefaultTransaction: true, // 跳过默认事务，提升单条查询性能
	})
	if err != nil {
		return nil, fmt.Errorf("打开 postgres 失败: %w", err)
	}

	// PostgreSQL 连接池优化：使用配置中的连接数
	// 生产环境建议：100 最大连接，50 空闲连接
	maxConns := cfg.GetMaxOpenConns()
	if maxConns <= 0 {
		maxConns = 100 // 默认值
	}
	if err := configureConnection(gormDB, maxConns); err != nil {
		return nil, err
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
