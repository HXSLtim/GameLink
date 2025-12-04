package main

// GameLink API
//
// @title           GameLink API
// @version         0.3.0
// @description     GameLink 平台 API，包含健康检查、认证与管理端能力
// @BasePath        /api/v1
// @schemes         http https
//
// @Tags Auth,User,Player,Admin,Admin/Commission,Admin/Games,Admin/Roles,Admin/Orders,Admin/Players,Admin/Payments,Admin/Reviews,Admin/Stats,Admin/Withdraws,Player - Profile,Player - Orders,Player - Earnings,Player - Commission,Player - Gifts,Player - Reviews,User - Orders,User - Payments,User - Wallet,User - Players,User - Reviews,User - Gifts,User - Chat,User - Feeds,System,Notifications
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

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

	// Initialize metrics
	metrics.Init(app.PrometheusRegistry)
	metrics.InitBusinessMetrics(app.PrometheusRegistry)

	// Expose metrics endpoint
	app.Engine.GET("/metrics", gin.WrapH(promhttp.HandlerFor(app.PrometheusRegistry, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	})))

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
