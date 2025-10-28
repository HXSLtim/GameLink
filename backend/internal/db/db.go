package db

import (
	"fmt"

	"gorm.io/gorm"

	"gamelink/internal/config"
)

// Open 根据配置创建数据库连接。
func Open(cfg config.AppConfig) (*gorm.DB, error) {
	switch cfg.Database.Type {
	case "sqlite":
		return openSQLite(cfg)
	case "postgres":
		return openPostgres(cfg)
	default:
		return nil, fmt.Errorf("暂不支持的数据库类型: %s", cfg.Database.Type)
	}
}
