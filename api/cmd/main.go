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
	"gamelink/pkg/lifecycle"
	"gamelink/pkg/metrics"
)

func main() {
	bootStart := time.Now()

	t := time.Now()
	app, err := container.NewApplication()
	if err != nil {
		log.Fatalf("failed to bootstrap application: %v", err)
	}
	log.Printf("[startup] container init: %v", time.Since(t))

	logCryptoStatus(app.Config)

	// Initialize metrics collector (includes HTTP and DB metrics)
	// 注意：InitBusinessMetrics 已在 Wire 生成的 initializeApplication 中调用，无需重复
	t = time.Now()
	metrics.NewCollector(app.PrometheusRegistry)
	log.Printf("[startup] metrics init: %v", time.Since(t))

	// Configure metrics authentication based on environment
	metricsAuthConfig := middleware.DefaultMetricsAuthConfig()

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "production" {
		if allowedCIDRs := os.Getenv("METRICS_ALLOWED_CIDRS"); allowedCIDRs != "" {
			metricsAuthConfig.AllowedCIDRs = parseCIDRs(allowedCIDRs)
		}
		metricsAuthConfig.Enabled = true
		log.Println("metrics endpoint: secured with IP whitelist")
	} else {
		metricsAuthConfig.Enabled = false
		log.Println("metrics endpoint: open (development mode)")
	}

	app.Engine.GET("/metrics",
		middleware.MetricsAuth(metricsAuthConfig),
		gin.WrapH(promhttp.HandlerFor(app.PrometheusRegistry, promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		})),
	)

	t = time.Now()
	if err := app.Lifecycle.Start(context.Background()); err != nil {
		log.Fatalf("failed to start services: %v", err)
	}
	log.Printf("[startup] lifecycle start: %v", time.Since(t))

	// 注意：用户管理表迁移已合并到主 autoMigrate 流程中，
	// 无需再单独打开第二个数据库连接。

	log.Printf("[startup] total boot time: %v", time.Since(bootStart))
	startServer(app.Engine, app.Config.Port, app.Lifecycle)
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

	// HTTP server shutdown: drain in-flight requests
	httpCtx, httpCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer httpCancel()

	if err := server.Shutdown(httpCtx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}

	// Lifecycle services shutdown: schedulers, ws hub, etc.
	svcCtx, svcCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer svcCancel()

	if err := lifecycle.Stop(svcCtx); err != nil {
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
