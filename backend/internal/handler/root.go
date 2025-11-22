package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoot 注册根路径路由。
func RegisterRoot(router gin.IRoutes) {
	router.GET("/", rootIndex)
}

func rootIndex(c *gin.Context) {
	RespondSuccess(c, "GameLink API", gin.H{
		"endpoints": []string{
			"/swagger",
			"/swagger.json",
			"/healthz",
			"/api/v1/",
			"/api/v1/healthz",
			"/api/v1/admin/games",
			"/api/v1/admin/users",
			"/api/v1/admin/players",
			"/api/v1/admin/orders",
			"/api/v1/admin/payments",
		},
	})
}
