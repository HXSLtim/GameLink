package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/apierr"
	"gamelink/internal/model"
)

// RespondJSON sends a JSON response with the given status code and data
func RespondJSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}

// RespondSuccess sends a successful JSON response
func RespondSuccess(c *gin.Context, message string, data interface{}) {
	response := model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: message,
	}

	// Set data if provided (handle nil and empty struct cases)
	if data != nil {
		// Skip adding data for empty structs (common for DELETE/PUT operations)
		if _, isEmpty := data.(struct{}); !isEmpty {
			response.Data = data
		}
	}

	// Add trace ID if available
	if traceID := GetRequestID(c); traceID != "" {
		response.TraceID = traceID
	}

	c.JSON(http.StatusOK, response)
}

// RespondCreated sends a created JSON response
func RespondCreated(c *gin.Context, data interface{}) {
	response := model.APIResponse[any]{
		Success: true,
		Code:    http.StatusCreated,
		Message: "created",
	}

	if data != nil {
		// Skip adding data for empty structs
		if _, isEmpty := data.(struct{}); !isEmpty {
			response.Data = data
		}
	}

	// Add trace ID if available
	if traceID := GetRequestID(c); traceID != "" {
		response.TraceID = traceID
	}

	c.JSON(http.StatusCreated, response)
}

// RespondAPIError sends an error JSON response
func RespondAPIError(c *gin.Context, err error) {
	// 使用统一的错误映射
	RespondWithServiceError(c, err)
}

// RespondError sends an error JSON response with the given status code and message
func RespondError(c *gin.Context, statusCode int, message string) {
	RespondAPIError(c, apierr.New(statusCode, message))
}

// RespondBadRequest sends a 400 Bad Request response
func RespondBadRequest(c *gin.Context, message string) {
	RespondAPIError(c, apierr.BadRequest(message))
}

// RespondUnauthorized sends a 401 Unauthorized response
func RespondUnauthorized(c *gin.Context, message string) {
	RespondAPIError(c, apierr.Unauthorized(message))
}

// RespondForbidden sends a 403 Forbidden response
func RespondForbidden(c *gin.Context, message string) {
	RespondAPIError(c, apierr.Forbidden(message))
}

// RespondNotFound sends a 404 Not Found response
func RespondNotFound(c *gin.Context, message string) {
	RespondAPIError(c, apierr.NotFound(message))
}

// RespondInternalError sends a 500 Internal Server Error response
func RespondInternalError(c *gin.Context, message string) {
	RespondAPIError(c, apierr.InternalError(message))
}

// RespondValidationError sends a validation error response
func RespondValidationError(c *gin.Context, field, message string) {
	err := apierr.NewValidationError(field, message)
	// Marshal the full ValidationError structure, not just the embedded APIError
	// This ensures the field-specific details (field, value, tag) are included
	var marshaled map[string]interface{}
	if body, marshalErr := json.Marshal(err); marshalErr == nil {
		if unmarshalErr := json.Unmarshal(body, &marshaled); unmarshalErr == nil {
			// Remove empty fields to keep response clean
			if err.Value == "" {
				delete(marshaled, "value")
			}
			if err.Tag == "" {
				delete(marshaled, "tag")
			}
			c.JSON(err.Code, marshaled)
			return
		}
	}
	// Fallback to basic error response if marshaling fails
	RespondAPIError(c, &err.APIError)
}

// BindAndValidate binds the request body and validates it
func BindAndValidate(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		RespondBadRequest(c, "无效的请求格式: "+err.Error())
		c.Abort()
		return err
	}
	return nil
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

// AddRequestID adds a request ID to the context and response header
func AddRequestID(c *gin.Context) {
	requestID := GetRequestID(c)
	if requestID == "" {
		// Generate a new request ID if not present
		requestID = generateRequestID()
		c.Set("request_id", requestID)
	}
	c.Header("X-Request-ID", requestID)
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	// Simple implementation, can be enhanced with UUID
	return fmt.Sprintf("req-%d", time.Now().UnixNano())
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

// ParseQueryInt parses an integer query parameter
func ParseQueryInt(c *gin.Context, paramName string, defaultValue int) (int, error) {
	valueStr := c.Query(paramName)
	if valueStr == "" {
		return defaultValue, nil
	}

	var value int
	if _, err := fmt.Sscanf(valueStr, "%d", &value); err != nil {
		return 0, apierr.BadRequest(fmt.Sprintf("无效的%s格式", paramName))
	}

	return value, nil
}

// ParseQueryUint64 parses a uint64 query parameter
func ParseQueryUint64(c *gin.Context, paramName string) (uint64, error) {
	valueStr := c.Query(paramName)
	if valueStr == "" {
		return 0, nil
	}

	var value uint64
	if _, err := fmt.Sscanf(valueStr, "%d", &value); err != nil {
		return 0, apierr.BadRequest(fmt.Sprintf("无效的%s格式", paramName))
	}

	return value, nil
}
