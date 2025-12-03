package router

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/handler"
)

// registerHealthRoutes 注册健康检查路由
func registerHealthRoutes(router *gin.Engine) {
	handler.RegisterHealth(router)
}
