package config

import (
	"log"
	"strings"
)

const defaultDBType = "sqlite"

var dsnSamples = map[string]string{
	"sqlite":    "file:./var/dev.db?mode=rwc&_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)",
	"postgres":  "postgres://user:password@localhost:5432/gamelink?sslmode=disable",
	"mysql":     "user:password@tcp(localhost:3306)/gamelink?parseTime=true&charset=utf8mb4",
	"sqlserver": "sqlserver://user:password@localhost:1433?database=gamelink",
}

// SampleDSN 返回常见数据库类型的 DSN 示例。
func SampleDSN(dbType string) string {
	return dsnSamples[dbType]
}

func normalizeDBType(input string) string {
	value := strings.ToLower(strings.TrimSpace(input))
	if _, ok := dsnSamples[value]; ok {
		return value
	}
	log.Printf("未知 DB_TYPE %q，使用默认值 %s", input, defaultDBType)
	return defaultDBType
}
