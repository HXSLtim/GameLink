package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterHealth 注册健康检查路由。
func RegisterHealth(router gin.IRoutes) {
	router.GET("/healthz", Health)
}

// Health 返回服务运行状态。
func Health(c *gin.Context) {
	RespondSuccess(c, "OK", gin.H{"status": "ok"})
}
