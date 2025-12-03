package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func openPostgres(cfg ConfigProvider, metrics MetricsProvider) (*gorm.DB, error) {
	gormDB, err := gorm.Open(postgres.Open(cfg.GetDatabaseDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("打开 postgres 失败: %w", err)
	}

	if err := configureConnection(gormDB, 25); err != nil {
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
