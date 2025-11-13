package user

import (
	"gamelink/internal/model"
	"github.com/gin-gonic/gin"
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
