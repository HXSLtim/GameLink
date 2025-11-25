package handler

import (
	_ "embed"

	"github.com/gin-gonic/gin"
)

//go:generate swag init -g cmd/main.go -o internal/handler/swagger --parseDependency --parseInternal

//go:embed swagger/index.html
var swaggerHTML []byte

//go:embed swagger/openapi.json
var swaggerJSON []byte

// RegisterSwagger 注册 Swagger UI 与文档路由。
// /swagger 返回内嵌 UI，/swagger.json 返回由 swag 自动生成的 OpenAPI。
func RegisterSwagger(router gin.IRoutes) {
	router.GET("/swagger", SwaggerUI)
	router.GET("/swagger.json", SwaggerSpec)
}

// SwaggerUI 提供内嵌 Swagger UI 页面。
func SwaggerUI(c *gin.Context) {
	c.Data(200, "text/html; charset=utf-8", swaggerHTML)
}

// SwaggerSpec 返回当前服务的 OpenAPI 描述。
func SwaggerSpec(c *gin.Context) {
	c.Data(200, "application/json", swaggerJSON)
}
