package main

// GameLink API
//
// @title           GameLink API
// @version         0.3.0
// @description     GameLink 平台 API，包含健康检查、认证与管理端能力
// @BasePath        /api/v1
// @schemes         http https
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

	"gamelink/internal/config"
	"gamelink/internal/container"
	"gamelink/internal/lifecycle"
)

func main() {
	app, err := container.NewApplication()
	if err != nil {
		log.Fatalf("failed to bootstrap application: %v", err)
	}

	logCryptoStatus(app.Config)

	if err := app.Lifecycle.Start(context.Background()); err != nil {
		log.Fatalf("failed to start services: %v", err)
	}

	startServer(app.Engine, app.Config.Port, app.Lifecycle)
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
