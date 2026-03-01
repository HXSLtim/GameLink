package config

import (
	"bufio"
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
	Signature     SignatureConfig
	Seed          SeedConfig
	SuperAdmin    SuperAdminConfig
	AdminAuth     AdminAuthConfig
	ExternalAPI   ExternalAPIConfig
}

// DatabaseConfig 描述数据库驱动与连接信息。
type DatabaseConfig struct {
	Type       string   `yaml:"type"`
	DSN        string   `yaml:"dsn"`
	ReaderDSNs []string `yaml:"reader_dsns"` // 从库 DSN 列表（读写分离）
	MaxConns   int      `yaml:"max_conns"`   // 最大连接数
	MaxIdle    int      `yaml:"max_idle"`    // 最大空闲连接数
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

// SignatureConfig 描述HMAC-SHA256请求签名验证配置。
type SignatureConfig struct {
	Enabled      bool     `yaml:"enabled"`
	SecretKey    string   `yaml:"secret_key"`
	HeaderName   string   `yaml:"header_name"`
	Methods      []string `yaml:"methods"`
	ExcludePaths []string `yaml:"exclude_paths"`
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

// ExternalAPIConfig 描述外部API凭证配置。
type ExternalAPIConfig struct {
	// 微信支付
	WeChatPay WeChatPayConfig `yaml:"wechat_pay"`
	// 支付宝
	Alipay AlipayConfig `yaml:"alipay"`
	// 短信服务
	SMS SMSConfig `yaml:"sms"`
	// OSS对象存储
	OSS OSSConfig `yaml:"oss"`
}

// WeChatPayConfig 微信支付配置
type WeChatPayConfig struct {
	AppID       string `yaml:"app_id"`        // 应用ID
	MchID       string `yaml:"mch_id"`        // 商户号
	APIKey      string `yaml:"api_key"`       // API密钥
	APICertPath string `yaml:"api_cert_path"` // API证书路径
	NotifyURL   string `yaml:"notify_url"`    // 支付结果通知URL
	Enabled     bool   `yaml:"enabled"`       // 是否启用
}

// AlipayConfig 支付宝配置
type AlipayConfig struct {
	AppID          string `yaml:"app_id"`           // 应用ID
	PrivateKeyPath string `yaml:"private_key_path"` // 应用私钥路径
	PublicKeyPath  string `yaml:"public_key_path"`  // 支付宝公钥路径
	NotifyURL      string `yaml:"notify_url"`       // 支付结果通知URL
	Enabled        bool   `yaml:"enabled"`          // 是否启用
}

// SMSConfig 短信服务配置
type SMSConfig struct {
	Provider  string `yaml:"provider"`   //服务商: aliyun, tencent
	AccessKey string `yaml:"access_key"` // AccessKey ID
	SecretKey string `yaml:"secret_key"` // AccessKey Secret
	SignName  string `yaml:"sign_name"`  // 短信签名
	Enabled   bool   `yaml:"enabled"`    // 是否启用
}

// OSSConfig 对象存储配置
type OSSConfig struct {
	Provider  string `yaml:"provider"`   // 服务商: aliyun, qcloud, minio
	Endpoint  string `yaml:"endpoint"`   // 访问域名
	AccessKey string `yaml:"access_key"` // AccessKey ID
	SecretKey string `yaml:"secret_key"` // AccessKey Secret
	Bucket    string `yaml:"bucket"`     // 存储桶名称
	Region    string `yaml:"region"`     // 区域
	Enabled   bool   `yaml:"enabled"`    // 是否启用
}

type cryptoFileConfig struct {
	Enabled      *bool    `yaml:"enabled"`
	SecretKey    string   `yaml:"secret_key"`
	IV           string   `yaml:"iv"`
	Methods      []string `yaml:"methods"`
	ExcludePaths []string `yaml:"exclude_paths"`
	UseSignature *bool    `yaml:"use_signature"`
}

type signatureFileConfig struct {
	Enabled      *bool    `yaml:"enabled"`
	SecretKey    string   `yaml:"secret_key"`
	HeaderName   string   `yaml:"header_name"`
	Methods      []string `yaml:"methods"`
	ExcludePaths []string `yaml:"exclude_paths"`
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

type externalAPIFileConfig struct {
	WeChatPay weChatPayFileConfig `yaml:"wechat_pay"`
	Alipay    alipayFileConfig    `yaml:"alipay"`
	SMS       smsFileConfig       `yaml:"sms"`
	OSS       ossFileConfig       `yaml:"oss"`
}

type weChatPayFileConfig struct {
	AppID       *string `yaml:"app_id"`
	MchID       *string `yaml:"mch_id"`
	APIKey      *string `yaml:"api_key"`
	APICertPath *string `yaml:"api_cert_path"`
	NotifyURL   *string `yaml:"notify_url"`
	Enabled     *bool   `yaml:"enabled"`
}

type alipayFileConfig struct {
	AppID          *string `yaml:"app_id"`
	PrivateKeyPath *string `yaml:"private_key_path"`
	PublicKeyPath  *string `yaml:"public_key_path"`
	NotifyURL      *string `yaml:"notify_url"`
	Enabled        *bool   `yaml:"enabled"`
}

type smsFileConfig struct {
	Provider  *string `yaml:"provider"`
	AccessKey *string `yaml:"access_key"`
	SecretKey *string `yaml:"secret_key"`
	SignName  *string `yaml:"sign_name"`
	Enabled   *bool   `yaml:"enabled"`
}

type ossFileConfig struct {
	Provider  *string `yaml:"provider"`
	Endpoint  *string `yaml:"endpoint"`
	AccessKey *string `yaml:"access_key"`
	SecretKey *string `yaml:"secret_key"`
	Bucket    *string `yaml:"bucket"`
	Region    *string `yaml:"region"`
	Enabled   *bool   `yaml:"enabled"`
}

type fileConfig struct {
	Server struct {
		Port          string `yaml:"port"`
		EnableSwagger *bool  `yaml:"enable_swagger"`
	} `yaml:"server"`
	Database    DatabaseConfig        `yaml:"database"`
	Cache       CacheConfig           `yaml:"cache"`
	Crypto      cryptoFileConfig      `yaml:"crypto"`
	Signature   signatureFileConfig   `yaml:"signature"`
	Auth        authFileConfig        `yaml:"auth"`
	Seed        SeedConfig            `yaml:"seed"`
	SuperAdmin  superAdminFileConfig  `yaml:"super_admin"`
	AdminAuth   adminAuthFileConfig   `yaml:"admin_auth"`
	ExternalAPI externalAPIFileConfig `yaml:"external_apis"`
}

// Load 读取配置文件及环境变量，生成最终配置。
func Load() AppConfig {
	loadDotEnvFallback()

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	cfg := AppConfig{
		Port:          DefaultPort,
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
		Signature: SignatureConfig{
			Enabled:      false,
			SecretKey:    "", // 必须显式配置，无默认值
			HeaderName:   "X-Signature",
			Methods:      []string{"POST", "PUT", "PATCH", "DELETE"},
			ExcludePaths: []string{"/api/v1/health", "/api/v1/ping", "/api/v1/auth/refresh"},
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
		ExternalAPI: ExternalAPIConfig{
			WeChatPay: WeChatPayConfig{Enabled: false},
			Alipay:    AlipayConfig{Enabled: false},
			SMS:       SMSConfig{Enabled: false},
			OSS:       OSSConfig{Enabled: false},
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

// loadDotEnvFallback loads key-value pairs from local .env files when
// process-level environment variables are absent.
//
// Priority rule:
// 1) Existing process env (highest, never overridden)
// 2) .env in current working directory
// 3) ../.env (useful when running from api/ and .env lives at repo root)
func loadDotEnvFallback() {
	candidates := []string{".env", filepath.Join("..", ".env")}
	for _, path := range candidates {
		if _, err := os.Stat(path); err != nil {
			continue
		}

		file, err := os.Open(path)
		if err != nil {
			log.Printf("读取 .env 文件失败 (%s): %v", path, err)
			continue
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			if strings.HasPrefix(line, "export ") {
				line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
			}

			key, value, ok := strings.Cut(line, "=")
			if !ok {
				continue
			}

			key = strings.TrimSpace(key)
			value = strings.TrimSpace(value)
			if key == "" {
				continue
			}

			if len(value) >= 2 {
				if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
					(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
					value = value[1 : len(value)-1]
				}
			}

			if _, exists := os.LookupEnv(key); exists {
				continue
			}
			_ = os.Setenv(key, value)
		}

		if err := scanner.Err(); err != nil {
			log.Printf("解析 .env 文件失败 (%s): %v", path, err)
		}

		_ = file.Close()
	}
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
	applySignatureConfig(&fc, cfg)
	applyAuthConfig(&fc, cfg)
	applySuperAdminConfig(&fc, cfg)
	applyExternalAPIConfig(&fc, cfg)
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

func applySignatureConfig(fc *fileConfig, cfg *AppConfig) {
	if fc.Signature.SecretKey != "" {
		cfg.Signature.SecretKey = fc.Signature.SecretKey
	}
	if fc.Signature.HeaderName != "" {
		cfg.Signature.HeaderName = fc.Signature.HeaderName
	}
	if len(fc.Signature.Methods) > 0 {
		cfg.Signature.Methods = normalizeHTTPMethods(fc.Signature.Methods)
	}
	if len(fc.Signature.ExcludePaths) > 0 {
		cfg.Signature.ExcludePaths = normalizePaths(fc.Signature.ExcludePaths)
	}
	if fc.Signature.Enabled != nil {
		cfg.Signature.Enabled = *fc.Signature.Enabled
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

func applyExternalAPIConfig(fc *fileConfig, cfg *AppConfig) {
	// WeChat Pay
	if fc.ExternalAPI.WeChatPay.AppID != nil {
		cfg.ExternalAPI.WeChatPay.AppID = *fc.ExternalAPI.WeChatPay.AppID
	}
	if fc.ExternalAPI.WeChatPay.MchID != nil {
		cfg.ExternalAPI.WeChatPay.MchID = *fc.ExternalAPI.WeChatPay.MchID
	}
	if fc.ExternalAPI.WeChatPay.APIKey != nil {
		cfg.ExternalAPI.WeChatPay.APIKey = *fc.ExternalAPI.WeChatPay.APIKey
	}
	if fc.ExternalAPI.WeChatPay.APICertPath != nil {
		cfg.ExternalAPI.WeChatPay.APICertPath = *fc.ExternalAPI.WeChatPay.APICertPath
	}
	if fc.ExternalAPI.WeChatPay.NotifyURL != nil {
		cfg.ExternalAPI.WeChatPay.NotifyURL = *fc.ExternalAPI.WeChatPay.NotifyURL
	}
	if fc.ExternalAPI.WeChatPay.Enabled != nil {
		cfg.ExternalAPI.WeChatPay.Enabled = *fc.ExternalAPI.WeChatPay.Enabled
	}

	// Alipay
	if fc.ExternalAPI.Alipay.AppID != nil {
		cfg.ExternalAPI.Alipay.AppID = *fc.ExternalAPI.Alipay.AppID
	}
	if fc.ExternalAPI.Alipay.PrivateKeyPath != nil {
		cfg.ExternalAPI.Alipay.PrivateKeyPath = *fc.ExternalAPI.Alipay.PrivateKeyPath
	}
	if fc.ExternalAPI.Alipay.PublicKeyPath != nil {
		cfg.ExternalAPI.Alipay.PublicKeyPath = *fc.ExternalAPI.Alipay.PublicKeyPath
	}
	if fc.ExternalAPI.Alipay.NotifyURL != nil {
		cfg.ExternalAPI.Alipay.NotifyURL = *fc.ExternalAPI.Alipay.NotifyURL
	}
	if fc.ExternalAPI.Alipay.Enabled != nil {
		cfg.ExternalAPI.Alipay.Enabled = *fc.ExternalAPI.Alipay.Enabled
	}

	// SMS
	if fc.ExternalAPI.SMS.Provider != nil {
		cfg.ExternalAPI.SMS.Provider = *fc.ExternalAPI.SMS.Provider
	}
	if fc.ExternalAPI.SMS.AccessKey != nil {
		cfg.ExternalAPI.SMS.AccessKey = *fc.ExternalAPI.SMS.AccessKey
	}
	if fc.ExternalAPI.SMS.SecretKey != nil {
		cfg.ExternalAPI.SMS.SecretKey = *fc.ExternalAPI.SMS.SecretKey
	}
	if fc.ExternalAPI.SMS.SignName != nil {
		cfg.ExternalAPI.SMS.SignName = *fc.ExternalAPI.SMS.SignName
	}
	if fc.ExternalAPI.SMS.Enabled != nil {
		cfg.ExternalAPI.SMS.Enabled = *fc.ExternalAPI.SMS.Enabled
	}

	// OSS
	if fc.ExternalAPI.OSS.Provider != nil {
		cfg.ExternalAPI.OSS.Provider = *fc.ExternalAPI.OSS.Provider
	}
	if fc.ExternalAPI.OSS.Endpoint != nil {
		cfg.ExternalAPI.OSS.Endpoint = *fc.ExternalAPI.OSS.Endpoint
	}
	if fc.ExternalAPI.OSS.AccessKey != nil {
		cfg.ExternalAPI.OSS.AccessKey = *fc.ExternalAPI.OSS.AccessKey
	}
	if fc.ExternalAPI.OSS.SecretKey != nil {
		cfg.ExternalAPI.OSS.SecretKey = *fc.ExternalAPI.OSS.SecretKey
	}
	if fc.ExternalAPI.OSS.Bucket != nil {
		cfg.ExternalAPI.OSS.Bucket = *fc.ExternalAPI.OSS.Bucket
	}
	if fc.ExternalAPI.OSS.Region != nil {
		cfg.ExternalAPI.OSS.Region = *fc.ExternalAPI.OSS.Region
	}
	if fc.ExternalAPI.OSS.Enabled != nil {
		cfg.ExternalAPI.OSS.Enabled = *fc.ExternalAPI.OSS.Enabled
	}
}

func overrideFromEnv(cfg *AppConfig) {
	overrideServerFromEnv(cfg)
	overrideDatabaseFromEnv(cfg)
	overrideCacheFromEnv(cfg)
	overrideCryptoFromEnv(cfg)
	overrideSignatureFromEnv(cfg)
	overrideAuthFromEnv(cfg)
	overrideSuperAdminFromEnv(cfg)
	overrideExternalAPIFromEnv(cfg)
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
	// 读写分离：从库 DSN 列表（逗号分隔）
	if readerDSNs := os.Getenv("DB_READER_DSNS"); readerDSNs != "" {
		cfg.Database.ReaderDSNs = strings.Split(readerDSNs, ",")
	}
	// 连接池配置
	if maxConns := os.Getenv("DB_MAX_CONNS"); maxConns != "" {
		if n, err := strconv.Atoi(maxConns); err == nil && n > 0 {
			cfg.Database.MaxConns = n
		}
	}
	if maxIdle := os.Getenv("DB_MAX_IDLE"); maxIdle != "" {
		if n, err := strconv.Atoi(maxIdle); err == nil && n > 0 {
			cfg.Database.MaxIdle = n
		}
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

func overrideSignatureFromEnv(cfg *AppConfig) {
	if v := os.Getenv("SIGNATURE_ENABLED"); v != "" {
		if enabled, err := strconv.ParseBool(v); err != nil {
			log.Printf("SIGNATURE_ENABLED=%q 无法解析，保持原值 %v", v, cfg.Signature.Enabled)
		} else {
			cfg.Signature.Enabled = enabled
		}
	}
	if secret := os.Getenv("SIGNATURE_SECRET_KEY"); secret != "" {
		cfg.Signature.SecretKey = secret
	}
	if header := os.Getenv("SIGNATURE_HEADER_NAME"); header != "" {
		cfg.Signature.HeaderName = header
	}
	if methods := os.Getenv("SIGNATURE_METHODS"); methods != "" {
		cfg.Signature.Methods = normalizeHTTPMethods(strings.Split(methods, ","))
	}
	if excludes := os.Getenv("SIGNATURE_EXCLUDE_PATHS"); excludes != "" {
		cfg.Signature.ExcludePaths = normalizePaths(strings.Split(excludes, ","))
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

func overrideExternalAPIFromEnv(cfg *AppConfig) {
	// WeChat Pay
	if appID := os.Getenv("WECHAT_PAY_APP_ID"); appID != "" {
		cfg.ExternalAPI.WeChatPay.AppID = appID
	}
	if mchID := os.Getenv("WECHAT_PAY_MCH_ID"); mchID != "" {
		cfg.ExternalAPI.WeChatPay.MchID = mchID
	}
	if apiKey := os.Getenv("WECHAT_PAY_API_KEY"); apiKey != "" {
		cfg.ExternalAPI.WeChatPay.APIKey = apiKey
	}
	if certPath := os.Getenv("WECHAT_PAY_API_CERT_PATH"); certPath != "" {
		cfg.ExternalAPI.WeChatPay.APICertPath = certPath
	}
	if notifyURL := os.Getenv("WECHAT_PAY_NOTIFY_URL"); notifyURL != "" {
		cfg.ExternalAPI.WeChatPay.NotifyURL = notifyURL
	}
	if enabled := os.Getenv("WECHAT_PAY_ENABLED"); enabled != "" {
		if e, err := strconv.ParseBool(enabled); err == nil {
			cfg.ExternalAPI.WeChatPay.Enabled = e
		}
	}

	// Alipay
	if appID := os.Getenv("ALIPAY_APP_ID"); appID != "" {
		cfg.ExternalAPI.Alipay.AppID = appID
	}
	if privateKeyPath := os.Getenv("ALIPAY_PRIVATE_KEY_PATH"); privateKeyPath != "" {
		cfg.ExternalAPI.Alipay.PrivateKeyPath = privateKeyPath
	}
	if publicKeyPath := os.Getenv("ALIPAY_PUBLIC_KEY_PATH"); publicKeyPath != "" {
		cfg.ExternalAPI.Alipay.PublicKeyPath = publicKeyPath
	}
	if notifyURL := os.Getenv("ALIPAY_NOTIFY_URL"); notifyURL != "" {
		cfg.ExternalAPI.Alipay.NotifyURL = notifyURL
	}
	if enabled := os.Getenv("ALIPAY_ENABLED"); enabled != "" {
		if e, err := strconv.ParseBool(enabled); err == nil {
			cfg.ExternalAPI.Alipay.Enabled = e
		}
	}

	// SMS
	if provider := os.Getenv("SMS_PROVIDER"); provider != "" {
		cfg.ExternalAPI.SMS.Provider = provider
	}
	if accessKey := os.Getenv("SMS_ACCESS_KEY"); accessKey != "" {
		cfg.ExternalAPI.SMS.AccessKey = accessKey
	}
	if secretKey := os.Getenv("SMS_SECRET_KEY"); secretKey != "" {
		cfg.ExternalAPI.SMS.SecretKey = secretKey
	}
	if signName := os.Getenv("SMS_SIGN_NAME"); signName != "" {
		cfg.ExternalAPI.SMS.SignName = signName
	}
	if enabled := os.Getenv("SMS_ENABLED"); enabled != "" {
		if e, err := strconv.ParseBool(enabled); err == nil {
			cfg.ExternalAPI.SMS.Enabled = e
		}
	}

	// OSS
	if provider := os.Getenv("OSS_PROVIDER"); provider != "" {
		cfg.ExternalAPI.OSS.Provider = provider
	}
	if endpoint := os.Getenv("OSS_ENDPOINT"); endpoint != "" {
		cfg.ExternalAPI.OSS.Endpoint = endpoint
	}
	if accessKey := os.Getenv("OSS_ACCESS_KEY"); accessKey != "" {
		cfg.ExternalAPI.OSS.AccessKey = accessKey
	}
	if secretKey := os.Getenv("OSS_SECRET_KEY"); secretKey != "" {
		cfg.ExternalAPI.OSS.SecretKey = secretKey
	}
	if bucket := os.Getenv("OSS_BUCKET"); bucket != "" {
		cfg.ExternalAPI.OSS.Bucket = bucket
	}
	if region := os.Getenv("OSS_REGION"); region != "" {
		cfg.ExternalAPI.OSS.Region = region
	}
	if enabled := os.Getenv("OSS_ENABLED"); enabled != "" {
		if e, err := strconv.ParseBool(enabled); err == nil {
			cfg.ExternalAPI.OSS.Enabled = e
		}
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
