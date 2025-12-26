package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoot 注册根路径路由。
func RegisterRoot(router gin.IRoutes) {
	router.GET("/", rootIndex)
}

func rootIndex(c *gin.Context) {
	// 当访问根路径时，重定向到 Swagger UI
	c.Redirect(302, "/swagger/index.html")
}
