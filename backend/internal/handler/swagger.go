package handler

import (
	_ "embed"

	"github.com/gin-gonic/gin"
)

//go:embed swagger/index.html
var swaggerHTML []byte

//go:embed swagger/openapi.json
var swaggerJSON []byte

// RegisterSwagger 注册 Swagger UI 与文档路由。
func RegisterSwagger(router gin.IRoutes) {
	router.GET("/swagger", SwaggerUI)
	router.GET("/swagger.json", SwaggerSpec)
}

// SwaggerUI 提供内嵌的 Swagger UI 页面。
func SwaggerUI(c *gin.Context) {
	c.Data(200, "text/html; charset=utf-8", swaggerHTML)
}

// SwaggerSpec 返回当前服务的 OpenAPI 描述。
func SwaggerSpec(c *gin.Context) {
	c.Data(200, "application/json", swaggerJSON)
}
