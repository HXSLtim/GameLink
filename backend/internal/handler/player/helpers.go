package player

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/apierr"
	"gamelink/internal/handler"
	"gamelink/internal/model"
)

// 本包内通用的响应封装
func respondJSON[T any](c *gin.Context, status int, payload model.APIResponse[T]) {
	if payload.TraceID == "" {
		if rid, ok := c.Get("request_id"); ok {
			if ridStr, ok := rid.(string); ok {
				payload.TraceID = ridStr
			}
		}
	}
	c.JSON(status, payload)
}

func respondError(c *gin.Context, status int, msg string) {
	respondJSON(c, status, model.APIResponse[any]{
		Success: false,
		Code:    status,
		Message: msg,
	})
}

// respondSuccess sends a successful response with message and optional data
func respondSuccess[T any](c *gin.Context, message string, data T) {
	handler.RespondSuccess(c, message, data)
}

// respondAPIError sends an API error response
func respondAPIError(c *gin.Context, err *apierr.APIError) {
	handler.RespondAPIError(c, err)
}

// parseUintParam 从路径参数中解析无符号整型 ID，调用方负责根据返回的 error 决定如何写入错误响应。
func parseUintParam(c *gin.Context, name string) (uint64, error) {
	value := c.Param(name)
	return strconv.ParseUint(value, 10, 64)
}

// getUserIDFromContext 从上下文获取用户ID
func getUserIDFromContext(c *gin.Context) uint64 {
	// 从 JWT 中间件设置的上下文中获取用户ID
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	userID, ok := userIDVal.(uint64)
	if !ok {
		return 0
	}
	return userID
}

