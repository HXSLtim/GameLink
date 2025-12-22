// Package resp provides unified response helpers for all handlers.
// 统一响应包，所有 handler 包应使用此包进行响应。
package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
)

// JSON writes a generic JSON response with trace ID support.
func JSON[T any](c *gin.Context, status int, payload model.APIResponse[T]) {
	if payload.TraceID == "" {
		if rid, ok := c.Get("request_id"); ok {
			if ridStr, ok := rid.(string); ok {
				payload.TraceID = ridStr
			}
		}
	}
	c.JSON(status, payload)
}

// Success sends a successful response with message and optional data.
func Success[T any](c *gin.Context, message string, data T) {
	JSON(c, http.StatusOK, model.APIResponse[T]{
		Success: true,
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

// OK sends a successful response with "OK" message.
func OK[T any](c *gin.Context, data T) {
	Success(c, "OK", data)
}

// Created sends a 201 Created response.
func Created[T any](c *gin.Context, data T) {
	JSON(c, http.StatusCreated, model.APIResponse[T]{
		Success: true,
		Code:    http.StatusCreated,
		Message: "created",
		Data:    data,
	})
}

// Updated sends a successful update response.
func Updated[T any](c *gin.Context, data T) {
	JSON(c, http.StatusOK, model.APIResponse[T]{
		Success: true,
		Code:    http.StatusOK,
		Message: "updated",
		Data:    data,
	})
}

// Deleted sends a successful delete response.
func Deleted(c *gin.Context) {
	JSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "deleted",
	})
}

// List sends a paginated list response.
func List[T any](c *gin.Context, data []T, pagination *model.Pagination) {
	if data == nil {
		data = make([]T, 0)
	}
	JSON(c, http.StatusOK, model.APIResponse[[]T]{
		Success:    true,
		Code:       http.StatusOK,
		Message:    "OK",
		Data:       data,
		Pagination: pagination,
	})
}
