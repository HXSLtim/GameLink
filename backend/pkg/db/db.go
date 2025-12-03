package db

import (
	"fmt"

	"gorm.io/gorm"
)

// DatabaseConfigProvider defines the database configuration interface
type DatabaseConfigProvider interface {
	GetDatabaseType() string
	ConfigProvider
}

// Open 根据配置创建数据库连接。
func Open(cfg DatabaseConfigProvider, metrics MetricsProvider) (*gorm.DB, error) {
	switch cfg.GetDatabaseType() {
	case "sqlite":
		return openSQLite(cfg, metrics)
	case "postgres":
		return openPostgres(cfg, metrics)
	default:
		return nil, fmt.Errorf("暂不支持的数据库类型: %s", cfg.GetDatabaseType())
	}
}
