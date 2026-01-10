package player

import (
	"gamelink/internal/handler/resp"
	"gamelink/internal/model"

	"github.com/gin-gonic/gin"
)

// respondJSON writes a JSON response with trace ID support.
// Kept for backward compatibility with existing handlers.
func respondJSON[T any](c *gin.Context, status int, payload model.APIResponse[T]) {
	resp.JSON(c, status, payload)
}

// respondError sends an error response with status code and message.
func respondError(c *gin.Context, status int, msg string) {
	resp.ErrorMsg(c, status, msg)
}

// respondSuccess sends a successful response with message and optional data.
func respondSuccess[T any](c *gin.Context, message string, data T) {
	resp.Success(c, message, data)
}

// respondAPIError sends an API error response.
func respondAPIError(c *gin.Context, err error) {
	resp.Error(c, err)
}

// parseUintParam parses a uint64 path parameter.
func parseUintParam(c *gin.Context, name string) (uint64, error) {
	return resp.ParseUintParam(c, name)
}

// getUserIDFromContext extracts user_id from context.
func getUserIDFromContext(c *gin.Context) uint64 {
	return resp.GetUserID(c)
}
