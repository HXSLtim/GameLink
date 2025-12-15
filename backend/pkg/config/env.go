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

const (
	// 默认Token有效期（小时）
	defaultTokenTTL = 24
)

// deprecatedDefaultJWTSecret 用于检测是否在使用已废弃的默认密钥
const deprecatedDefaultJWTSecret = "gamelink-default-secret-key-change-in-development"

// AppConfig 汇总服务运行所需的核心配置。
type AppConfig struct {
	Port          string
	EnableSwagger bool
	Database      DatabaseConfig
	Cache         CacheConfig
	Crypto        CryptoConfig
	Auth          AuthConfig
	Seed          SeedConfig
	SuperAdmin    SuperAdminConfig
	AdminAuth     AdminAuthConfig
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

// CryptoConfig 描述前后端对称加密配置。
type CryptoConfig struct {
	Enabled      bool     `yaml:"enabled"`
	SecretKey    string   `yaml:"secret_key"`
	IV           string   `yaml:"iv"`
	Methods      []string `yaml:"methods"`
	ExcludePaths []string `yaml:"exclude_paths"`
	UseSignature bool     `yaml:"use_signature"`
}

// AuthConfig 描述鉴权配置。
type AuthConfig struct {
	JWTSecret     string `yaml:"jwt_secret"`
	TokenTTLHours int    `yaml:"token_ttl_hours"`
}

// SeedConfig 控制是否注入演示数据。
type SeedConfig struct {
	Enabled bool `yaml:"enabled"`
}

// SuperAdminConfig 描述超级管理员初始化配置。
type SuperAdminConfig struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Phone    string `yaml:"phone"`
}

// AdminAuthConfig 描述管理员认证模式配置。
type AdminAuthConfig struct {
	Mode string `yaml:"mode"`
}

type cryptoFileConfig struct {
	Enabled      *bool    `yaml:"enabled"`
	SecretKey    string   `yaml:"secret_key"`
	IV           string   `yaml:"iv"`
	Methods      []string `yaml:"methods"`
	ExcludePaths []string `yaml:"exclude_paths"`
	UseSignature *bool    `yaml:"use_signature"`
}

type authFileConfig struct {
	JWTSecret     string `yaml:"jwt_secret"`
	TokenTTLHours *int   `yaml:"token_ttl_hours"`
}

type superAdminFileConfig struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Phone    string `yaml:"phone"`
}

type adminAuthFileConfig struct {
	Mode string `yaml:"mode"`
}

type fileConfig struct {
	Server struct {
		Port          string `yaml:"port"`
		EnableSwagger *bool  `yaml:"enable_swagger"`
	} `yaml:"server"`
	Database   DatabaseConfig       `yaml:"database"`
	Cache      CacheConfig          `yaml:"cache"`
	Crypto     cryptoFileConfig     `yaml:"crypto"`
	Auth       authFileConfig       `yaml:"auth"`
	Seed       SeedConfig           `yaml:"seed"`
	SuperAdmin superAdminFileConfig `yaml:"super_admin"`
	AdminAuth  adminAuthFileConfig  `yaml:"admin_auth"`
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
		Crypto: CryptoConfig{
			Enabled:      false,
			SecretKey:    "", // 必须显式配置，无默认值
			IV:           "", // 必须显式配置，无默认值
			Methods:      []string{"POST", "PUT", "PATCH"},
			ExcludePaths: []string{"/api/v1/health", "/api/v1/ping", "/api/v1/auth/refresh"},
			UseSignature: true,
		},
		Auth: AuthConfig{
			JWTSecret:     "",
			TokenTTLHours: defaultTokenTTL,
		},
		Seed: SeedConfig{Enabled: false},
		SuperAdmin: SuperAdminConfig{
			Email:    "", // 必须显式配置，无默认值
			Password: "", // 必须显式配置，无默认值
			Name:     "Super Admin",
			Phone:    "",
		},
		AdminAuth: AdminAuthConfig{
			Mode: "admin", // 默认使用 AdminAuth，生产环境建议使用 jwt
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

	if cfg.Auth.TokenTTLHours <= 0 {
		cfg.Auth.TokenTTLHours = defaultTokenTTL
	}

	// 检查是否在使用已废弃的默认 JWT secret
	if strings.TrimSpace(cfg.Auth.JWTSecret) == deprecatedDefaultJWTSecret {
		log.Printf("⚠️  警告：正在使用已废弃的默认 JWT secret '%s'", deprecatedDefaultJWTSecret)
		log.Printf("   请立即通过配置文件或 JWT_SECRET_KEY 环境变量设置安全的密钥！")
		log.Printf("   建议密钥长度至少 32 个字符，包含大小写字母、数字和特殊符号")
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

	applyServerConfig(&fc, cfg)
	applyDatabaseConfig(&fc, cfg)
	applyCacheConfig(&fc, cfg)
	applyCryptoConfig(&fc, cfg)
	applyAuthConfig(&fc, cfg)
	applySuperAdminConfig(&fc, cfg)
}

func applyServerConfig(fc *fileConfig, cfg *AppConfig) {
	if fc.Server.Port != "" {
		cfg.Port = fc.Server.Port
	}
	if fc.Server.EnableSwagger != nil {
		cfg.EnableSwagger = *fc.Server.EnableSwagger
	}
}

func applyDatabaseConfig(fc *fileConfig, cfg *AppConfig) {
	if fc.Database.Type != "" {
		cfg.Database.Type = normalizeDBType(fc.Database.Type)
	}
	if fc.Database.DSN != "" {
		cfg.Database.DSN = fc.Database.DSN
	}
}

func applyCacheConfig(fc *fileConfig, cfg *AppConfig) {
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

func applyCryptoConfig(fc *fileConfig, cfg *AppConfig) {
	if fc.Crypto.SecretKey != "" {
		cfg.Crypto.SecretKey = fc.Crypto.SecretKey
	}
	if fc.Crypto.IV != "" {
		cfg.Crypto.IV = fc.Crypto.IV
	}
	if len(fc.Crypto.Methods) > 0 {
		cfg.Crypto.Methods = normalizeHTTPMethods(fc.Crypto.Methods)
	}
	if len(fc.Crypto.ExcludePaths) > 0 {
		cfg.Crypto.ExcludePaths = normalizePaths(fc.Crypto.ExcludePaths)
	}
	if fc.Crypto.Enabled != nil {
		cfg.Crypto.Enabled = *fc.Crypto.Enabled
	}
	if fc.Crypto.UseSignature != nil {
		cfg.Crypto.UseSignature = *fc.Crypto.UseSignature
	}
}

func applyAuthConfig(fc *fileConfig, cfg *AppConfig) {
	if fc.Auth.JWTSecret != "" {
		cfg.Auth.JWTSecret = fc.Auth.JWTSecret
	}
	if fc.Auth.TokenTTLHours != nil {
		cfg.Auth.TokenTTLHours = *fc.Auth.TokenTTLHours
	}
	if fc.Seed.Enabled {
		cfg.Seed.Enabled = fc.Seed.Enabled
	}
	if fc.AdminAuth.Mode != "" {
		cfg.AdminAuth.Mode = fc.AdminAuth.Mode
	}
}

func applySuperAdminConfig(fc *fileConfig, cfg *AppConfig) {
	if fc.SuperAdmin.Email != "" {
		cfg.SuperAdmin.Email = fc.SuperAdmin.Email
	}
	if fc.SuperAdmin.Password != "" {
		cfg.SuperAdmin.Password = fc.SuperAdmin.Password
	}
	if fc.SuperAdmin.Name != "" {
		cfg.SuperAdmin.Name = fc.SuperAdmin.Name
	}
	if fc.SuperAdmin.Phone != "" {
		cfg.SuperAdmin.Phone = fc.SuperAdmin.Phone
	}
}

func overrideFromEnv(cfg *AppConfig) {
	overrideServerFromEnv(cfg)
	overrideDatabaseFromEnv(cfg)
	overrideCacheFromEnv(cfg)
	overrideCryptoFromEnv(cfg)
	overrideAuthFromEnv(cfg)
	overrideSuperAdminFromEnv(cfg)
}

func overrideServerFromEnv(cfg *AppConfig) {
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
}

func overrideDatabaseFromEnv(cfg *AppConfig) {
	if dbType := os.Getenv("DB_TYPE"); dbType != "" {
		cfg.Database.Type = normalizeDBType(dbType)
	}
	if dsn := os.Getenv("DB_DSN"); dsn != "" {
		cfg.Database.DSN = dsn
	}
}

func overrideCacheFromEnv(cfg *AppConfig) {
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

func overrideCryptoFromEnv(cfg *AppConfig) {
	if v := os.Getenv("CRYPTO_ENABLED"); v != "" {
		if enabled, err := strconv.ParseBool(v); err != nil {
			log.Printf("CRYPTO_ENABLED=%q 无法解析，保持原值 %v", v, cfg.Crypto.Enabled)
		} else {
			cfg.Crypto.Enabled = enabled
		}
	}
	if secret := os.Getenv("CRYPTO_SECRET_KEY"); secret != "" {
		cfg.Crypto.SecretKey = secret
	}
	if iv := os.Getenv("CRYPTO_IV"); iv != "" {
		cfg.Crypto.IV = iv
	}
	if methods := os.Getenv("CRYPTO_METHODS"); methods != "" {
		cfg.Crypto.Methods = normalizeHTTPMethods(strings.Split(methods, ","))
	}
	if excludes := os.Getenv("CRYPTO_EXCLUDE_PATHS"); excludes != "" {
		cfg.Crypto.ExcludePaths = normalizePaths(strings.Split(excludes, ","))
	}
	if useSignature := os.Getenv("CRYPTO_USE_SIGNATURE"); useSignature != "" {
		if enabled, err := strconv.ParseBool(useSignature); err != nil {
			log.Printf("CRYPTO_USE_SIGNATURE=%q 无法解析，保持原值 %v", useSignature, cfg.Crypto.UseSignature)
		} else {
			cfg.Crypto.UseSignature = enabled
		}
	}
}

func overrideAuthFromEnv(cfg *AppConfig) {
	if jwtSecret := os.Getenv("JWT_SECRET_KEY"); jwtSecret != "" {
		cfg.Auth.JWTSecret = jwtSecret
	}
	if ttl := os.Getenv("JWT_TOKEN_TTL_HOURS"); ttl != "" {
		if hours, err := strconv.Atoi(ttl); err != nil {
			log.Printf("JWT_TOKEN_TTL_HOURS=%q 无法解析，保持原值 %d", ttl, cfg.Auth.TokenTTLHours)
		} else {
			cfg.Auth.TokenTTLHours = hours
		}
	}
	if seed := os.Getenv("SEED_ENABLED"); seed != "" {
		if enabled, err := strconv.ParseBool(seed); err != nil {
			log.Printf("SEED_ENABLED=%q 无法解析，保持原值 %v", seed, cfg.Seed.Enabled)
		} else {
			cfg.Seed.Enabled = enabled
		}
	}
	if mode := os.Getenv("ADMIN_AUTH_MODE"); mode != "" {
		cfg.AdminAuth.Mode = strings.ToLower(mode)
	}
}

func overrideSuperAdminFromEnv(cfg *AppConfig) {
	if email := os.Getenv("SUPER_ADMIN_EMAIL"); email != "" {
		cfg.SuperAdmin.Email = email
	}
	if password := os.Getenv("SUPER_ADMIN_PASSWORD"); password != "" {
		cfg.SuperAdmin.Password = password
	}
	if name := os.Getenv("SUPER_ADMIN_NAME"); name != "" {
		cfg.SuperAdmin.Name = name
	}
	if phone := os.Getenv("SUPER_ADMIN_PHONE"); phone != "" {
		cfg.SuperAdmin.Phone = phone
	}
}

func normalizeHTTPMethods(methods []string) []string {
	var normalized []string
	for _, method := range methods {
		trimmed := strings.TrimSpace(strings.ToUpper(method))
		if trimmed == "" {
			continue
		}
		normalized = append(normalized, trimmed)
	}
	if len(normalized) == 0 {
		return []string{"POST", "PUT", "PATCH"}
	}
	return normalized
}

func normalizePaths(paths []string) []string {
	var normalized []string
	for _, path := range paths {
		trimmed := strings.TrimSpace(path)
		if trimmed == "" {
			continue
		}
		normalized = append(normalized, trimmed)
	}
	return normalized
}
