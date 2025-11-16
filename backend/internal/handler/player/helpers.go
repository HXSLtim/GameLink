package player

import (
	"strconv"

	"github.com/gin-gonic/gin"

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

// parseUintParam 从路径参数中解析无符号整型 ID，调用方负责根据返回的 error 决定如何写入错误响应。
func parseUintParam(c *gin.Context, name string) (uint64, error) {
	value := c.Param(name)
	return strconv.ParseUint(value, 10, 64)
}

// 从上下文获取用户ID（由JWT中间件写入）
func getUserIDFromContext(c *gin.Context) uint64 {
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
