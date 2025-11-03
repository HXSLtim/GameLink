package user

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