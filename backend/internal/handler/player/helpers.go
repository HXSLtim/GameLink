package player

import (
    "github.com/gin-gonic/gin"
    "gamelink/internal/model"
)

// 本包内通用的响应封装
func respondJSON[T any](c *gin.Context, status int, payload model.APIResponse[T]) {
    c.JSON(status, payload)
}

func respondError(c *gin.Context, status int, msg string) {
    respondJSON(c, status, model.APIResponse[any]{
        Success: false,
        Code:    status,
        Message: msg,
    })
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