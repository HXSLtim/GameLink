package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// AppConfig 汇总服务运行所需的核心配置。
type AppConfig struct {
	Port          string
	EnableSwagger bool
	Database      DatabaseConfig
	Cache         CacheConfig
}

// DatabaseConfig 描述数据库驱动与连接信息。
type DatabaseConfig struct {
	Type string `yaml:"type"`
	DSN  string `yaml:"dsn"`
}

// CacheConfig 描述缓存配置。
type CacheConfig struct {
	Type  string      `yaml:"type"`
	Redis RedisConfig `yaml:"redis"`
}

// RedisConfig 包含 redis 连接信息。
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type fileConfig struct {
	Server struct {
		Port          string `yaml:"port"`
		EnableSwagger *bool  `yaml:"enable_swagger"`
	} `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Cache    CacheConfig    `yaml:"cache"`
}

// Load 读取配置文件及环境变量，生成最终配置。
func Load() AppConfig {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	cfg := AppConfig{
		Port:          "8080",
		EnableSwagger: true,
		Database: DatabaseConfig{
			Type: defaultDBType,
			DSN:  "",
		},
		Cache: CacheConfig{
			Type: "memory",
			Redis: RedisConfig{
				Addr:     "127.0.0.1:6379",
				Password: "",
				DB:       0,
			},
		},
	}

	loadFromFile(env, &cfg)
	overrideFromEnv(&cfg)

	if cfg.Database.DSN == "" {
		if env != "production" {
			if sample := SampleDSN(cfg.Database.Type); sample != "" {
				cfg.Database.DSN = sample
				log.Printf("DB_DSN 未配置，使用 %s 示例 DSN：%s", cfg.Database.Type, sample)
			}
		} else {
			log.Printf("DB_DSN 未配置，生产环境将保持为空并由外部注入")
		}
	}

	return cfg
}

func loadFromFile(env string, cfg *AppConfig) {
	configPath := filepath.Join("configs", fmt.Sprintf("config.%s.yaml", env))
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Printf("读取配置文件 %s 失败：%v，使用默认值", configPath, err)
		return
	}

	var fc fileConfig
	if err := yaml.Unmarshal(data, &fc); err != nil {
		log.Printf("解析配置文件 %s 失败：%v，使用默认值", configPath, err)
		return
	}

	if fc.Server.Port != "" {
		cfg.Port = fc.Server.Port
	}
	if fc.Server.EnableSwagger != nil {
		cfg.EnableSwagger = *fc.Server.EnableSwagger
	}
	if fc.Database.Type != "" {
		cfg.Database.Type = normalizeDBType(fc.Database.Type)
	}
	if fc.Database.DSN != "" {
		cfg.Database.DSN = fc.Database.DSN
	}
	if fc.Cache.Type != "" {
		cfg.Cache.Type = strings.ToLower(fc.Cache.Type)
	}
	if fc.Cache.Redis.Addr != "" {
		cfg.Cache.Redis.Addr = fc.Cache.Redis.Addr
	}
	if fc.Cache.Redis.Password != "" {
		cfg.Cache.Redis.Password = fc.Cache.Redis.Password
	}
	if fc.Cache.Redis.DB != 0 {
		cfg.Cache.Redis.DB = fc.Cache.Redis.DB
	}
}

func overrideFromEnv(cfg *AppConfig) {
	if port := os.Getenv("SERVICE_PORT"); port != "" {
		cfg.Port = port
	}

	if v := os.Getenv("ENABLE_SWAGGER"); v != "" {
		if enabled, err := strconv.ParseBool(v); err != nil {
			log.Printf("ENABLE_SWAGGER=%q 无法解析，保持原值 %v", v, cfg.EnableSwagger)
		} else {
			cfg.EnableSwagger = enabled
		}
	}

	if dbType := os.Getenv("DB_TYPE"); dbType != "" {
		cfg.Database.Type = normalizeDBType(dbType)
	}

	if dsn := os.Getenv("DB_DSN"); dsn != "" {
		cfg.Database.DSN = dsn
	}

	if cacheType := os.Getenv("CACHE_TYPE"); cacheType != "" {
		cfg.Cache.Type = strings.ToLower(cacheType)
	}
	if redisAddr := os.Getenv("REDIS_ADDR"); redisAddr != "" {
		cfg.Cache.Redis.Addr = redisAddr
	}
	if redisPassword := os.Getenv("REDIS_PASSWORD"); redisPassword != "" {
		cfg.Cache.Redis.Password = redisPassword
	}
	if redisDB := os.Getenv("REDIS_DB"); redisDB != "" {
		if db, err := strconv.Atoi(redisDB); err != nil {
			log.Printf("REDIS_DB=%q 无法解析，保持原值 %d", redisDB, cfg.Cache.Redis.DB)
		} else {
			cfg.Cache.Redis.DB = db
		}
	}
}
