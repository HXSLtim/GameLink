package router

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"

	"gamelink/pkg/config"
)

func TestRegisterPublicRoutes_RegistersPaymentNotifyRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := gin.New()
	api := engine.Group("/api/v1")

	registerPublicRoutes(api, nil, nil, config.AppConfig{}, nil)

	requireRoute(t, engine, http.MethodPost, "/api/v1/public/payments/wechat/notify")
	requireRoute(t, engine, http.MethodPost, "/api/v1/public/payments/alipay/notify")
}

func requireRoute(t *testing.T, engine *gin.Engine, method, path string) {
	t.Helper()

	for _, route := range engine.Routes() {
		if route.Method == method && route.Path == path {
			return
		}
	}

	t.Fatalf("route not registered: %s %s", method, path)
}
