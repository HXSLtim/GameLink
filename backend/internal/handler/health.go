package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterHealth 注册健康检查路由。
func RegisterHealth(router gin.IRoutes) {
	router.GET("/healthz", Health)
}

// HealthStatus 健康检查响应
type HealthStatus struct {
	Status string `json:"status" example:"ok"`
}

// Health 返回服务运行状态。
// @Summary      健康检查
// @Description  返回服务运行状态，用于负载均衡器和监控系统
// @Tags         System
// @Produce      json
// @Success      200  {object}  HealthStatus
// @Router       /healthz [get]
func Health(c *gin.Context) {
	RespondSuccess(c, "OK", gin.H{"status": "ok"})
}
