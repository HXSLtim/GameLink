package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/pkg/apierr"
)

// RespondAPIError sends an error JSON response using unified error mapping
func RespondAPIError(c *gin.Context, err error) {
	RespondWithServiceError(c, err)
}

// RespondSuccess sends a successful JSON response
func RespondSuccess(c *gin.Context, message string, data interface{}) {
	response := model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: message,
	}

	if data != nil {
		if _, isEmpty := data.(struct{}); !isEmpty {
			response.Data = data
		}
	}

	if traceID := GetRequestID(c); traceID != "" {
		response.TraceID = traceID
	}

	c.JSON(http.StatusOK, response)
}

// RespondBadRequest sends a 400 Bad Request response
func RespondBadRequest(c *gin.Context, message string) {
	RespondAPIError(c, apierr.BadRequest(message))
}

// GetRequestID gets the request ID from the context
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

// ParseIDParam parses and validates an ID parameter from the URL
func ParseIDParam(c *gin.Context, paramName string) (uint64, error) {
	idStr := c.Param(paramName)
	if idStr == "" {
		return 0, apierr.BadRequest(fmt.Sprintf("缺少%s参数", paramName))
	}

	var id uint64
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		return 0, apierr.BadRequest(fmt.Sprintf("无效的%s格式", paramName))
	}

	if id == 0 {
		return 0, apierr.BadRequest(fmt.Sprintf("%s不能为0", paramName))
	}

	return id, nil
}
