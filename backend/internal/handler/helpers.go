package handler

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
)

// writeJSON writes a JSON response to the client
func writeJSON(c *gin.Context, status int, payload any) {
	c.JSON(status, payload)
}

// writeJSONError writes a JSON error response to the client
func writeJSONError(c *gin.Context, status int, message string) {
	writeJSON(c, status, model.APIResponse[any]{
		Success: false,
		Code:    status,
		Message: message,
	})
}

// writeJSONSuccess writes a JSON success response with data to the client
func writeJSONSuccess(c *gin.Context, status int, data any) {
	writeJSON(c, status, model.APIResponse[any]{
		Success: true,
		Code:    status,
		Data:    data,
	})
}
