package main

// GameLink API
//
// @title           GameLink API
// @version         0.3.0
// @description     GameLink 游戏陪玩服务平台 API
// @description     提供用户认证、订单管理、支付结算、陪玩师管理等功能
// @termsOfService  https://gamelink.com/terms/
// @contact.name    GameLink Support
// @contact.email   support@gamelink.com
// @license.name    MIT
// @license.url     https://opensource.org/licenses/MIT
// @host            localhost:8080
// @BasePath        /api/v1
// @schemes         http https
//
// @tag.name        Auth
// @tag.description 用户认证相关接口
// @tag.name        User - Orders
// @tag.description 用户订单管理
// @tag.name        User - Payments
// @tag.description 用户支付相关
// @tag.name        User - Wallet
// @tag.description 用户钱包管理
// @tag.name        User - Players
// @tag.description 浏览陪玩师
// @tag.name        User - Reviews
// @tag.description 用户评价管理
// @tag.name        User - Gifts
// @tag.description 礼物功能
// @tag.name        User - Chat
// @tag.description 聊天功能
// @tag.name        Player - Profile
// @tag.description 陪玩师个人资料
// @tag.name        Player - Orders
// @tag.description 陪玩师订单管理
// @tag.name        Player - Earnings
// @tag.description 陪玩师收益管理
// @tag.name        Player - Commission
// @tag.description 陪玩师抽成管理
// @tag.name        Admin - Dashboard
// @tag.description 管理后台仪表盘
// @tag.name        Admin - Stats
// @tag.description 管理后台统计
// @tag.name        Admin/Users
// @tag.description 用户管理
// @tag.name        Admin/Players
// @tag.description 陪玩师管理
// @tag.name        Admin/Orders
// @tag.description 订单管理
// @tag.name        Admin/Monitor
// @tag.description 系统监控
// @tag.name        Notifications
// @tag.description 通知中心
// @tag.name        System
// @tag.description 系统接口
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT Bearer Token 认证，格式: Bearer {token}

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"gamelink/internal/handler/middleware"
	"gamelink/pkg/config"
	"gamelink/pkg/container"
	"gamelink/pkg/db"
	"gamelink/pkg/lifecycle"
	"gamelink/pkg/metrics"

	"gorm.io/gorm"
)

type cfgAdapter struct{ cfg config.AppConfig }

func (a *cfgAdapter) GetDatabaseDSN() string  { return a.cfg.Database.DSN }
func (a *cfgAdapter) IsSeedEnabled() bool     { return false }
func (a *cfgAdapter) GetMaxOpenConns() int    { return 1 }
func (a *cfgAdapter) GetDatabaseType() string { return a.cfg.Database.Type }

type noopMetrics struct{}

func (n *noopMetrics) InstrumentGorm(db interface{}) error { return nil }

func main() {
	app, err := container.NewApplication()
	if err != nil {
		log.Fatalf("failed to bootstrap application: %v", err)
	}

	logCryptoStatus(app.Config)

	// Initialize metrics collector (includes HTTP and DB metrics)
	metrics.NewCollector(app.PrometheusRegistry)

	// Initialize business metrics
	metrics.InitBusinessMetrics(app.PrometheusRegistry)

	// Configure metrics authentication based on environment
	metricsAuthConfig := middleware.DefaultMetricsAuthConfig()

	// In production, require admin authentication or IP whitelist
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "production" {
		// Check if METRICS_ALLOWED_CIDRS is set for private network access
		if allowedCIDRs := os.Getenv("METRICS_ALLOWED_CIDRS"); allowedCIDRs != "" {
			metricsAuthConfig.AllowedCIDRs = parseCIDRs(allowedCIDRs)
		}
		metricsAuthConfig.Enabled = true
		log.Println("metrics endpoint: secured with IP whitelist")
	} else {
		// In development/staging, allow all access
		metricsAuthConfig.Enabled = false
		log.Println("metrics endpoint: open (development mode)")
	}

	// Expose metrics endpoint with authentication
	app.Engine.GET("/metrics",
		middleware.MetricsAuth(metricsAuthConfig),
		gin.WrapH(promhttp.HandlerFor(app.PrometheusRegistry, promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		})),
	)

	if err := app.Lifecycle.Start(context.Background()); err != nil {
		log.Fatalf("failed to start services: %v", err)
	}

	// 迁移用户管理相关表结构
	if err := migrateUserManagement(app.Config); err != nil {
		log.Fatalf("failed to migrate user management tables: %v", err)
	}

	startServer(app.Engine, app.Config.Port, app.Lifecycle)
}

// migrateUserManagement 打开数据库并执行用户管理模块的迁移
func migrateUserManagement(cfg config.AppConfig) error {
	orm, err := db.Open(&cfgAdapter{cfg: cfg}, &noopMetrics{})
	if err != nil {
		return err
	}
	if err := db.MigrateUserManagement(orm); err != nil {
		return err
	}
	return closeSQL(orm)
}

// closeSQL 关闭底层 SQL 连接
func closeSQL(orm *gorm.DB) error {
	sqlDB, err := orm.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func startServer(router *gin.Engine, port string, lifecycle *lifecycle.Manager) {
	addr := fmt.Sprintf(":%s", port)

	server := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second, // 读取请求超时（缩短以快速释放连接）
		WriteTimeout:      30 * time.Second, // 写入响应超时
		IdleTimeout:       60 * time.Second, // Keep-Alive 空闲超时（缩短以提高连接复用效率）
		MaxHeaderBytes:    1 << 20,          // 1MB 最大请求头
	}

	go func() {
		log.Printf("api listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	if err := lifecycle.Stop(ctx); err != nil {
		log.Printf("service shutdown encountered errors: %v", err)
	}

	log.Println("Server stopped gracefully")
}

func logCryptoStatus(cfg config.AppConfig) {
	if cfg.Crypto.Enabled {
		log.Printf("crypto middleware enabled, methods=%v exclude=%v use_signature=%v", cfg.Crypto.Methods, cfg.Crypto.ExcludePaths, cfg.Crypto.UseSignature)
		return
	}
	log.Println("crypto middleware disabled")
}

// parseCIDRs parses a comma-separated list of CIDR blocks
func parseCIDRs(cidrsStr string) []string {
	if cidrsStr == "" {
		return nil
	}

	cidrs := make([]string, 0)
	for _, cidr := range strings.Split(cidrsStr, ",") {
		trimmed := strings.TrimSpace(cidr)
		if trimmed != "" {
			cidrs = append(cidrs, trimmed)
		}
	}
	return cidrs
}
