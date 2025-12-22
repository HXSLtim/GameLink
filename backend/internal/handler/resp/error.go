package resp

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
)

// Error sends an error response. Handles apierr.APIError and common error types.
func Error(c *gin.Context, err error) {
	// Handle apierr.APIError
	if apiErr, ok := err.(*apierr.APIError); ok {
		JSON(c, apiErr.Code, model.APIResponse[any]{
			Success: false,
			Code:    apiErr.Code,
			Message: apiErr.Message,
			TraceID: apiErr.RequestID,
		})
		return
	}

	// Handle repository.ErrNotFound
	if errors.Is(err, repository.ErrNotFound) {
		JSON(c, http.StatusNotFound, model.APIResponse[any]{
			Success: false,
			Code:    http.StatusNotFound,
			Message: "resource not found",
		})
		return
	}

	// Handle validation errors
	if apierr.IsValidationError(err) {
		JSON(c, http.StatusBadRequest, model.APIResponse[any]{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	// Default: internal server error
	JSON(c, http.StatusInternalServerError, model.APIResponse[any]{
		Success: false,
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	})
}

// ErrorMsg sends an error response with status code and message.
func ErrorMsg(c *gin.Context, status int, msg string) {
	JSON(c, status, model.APIResponse[any]{
		Success: false,
		Code:    status,
		Message: msg,
	})
}

// BadRequest sends a 400 Bad Request response.
func BadRequest(c *gin.Context, message string) {
	Error(c, apierr.BadRequest(message))
}

// Unauthorized sends a 401 Unauthorized response.
func Unauthorized(c *gin.Context, message string) {
	Error(c, apierr.Unauthorized(message))
}

// Forbidden sends a 403 Forbidden response.
func Forbidden(c *gin.Context, message string) {
	Error(c, apierr.Forbidden(message))
}

// NotFound sends a 404 Not Found response.
func NotFound(c *gin.Context, message string) {
	Error(c, apierr.NotFound(message))
}

// InternalError sends a 500 Internal Server Error response.
func InternalError(c *gin.Context, message string) {
	Error(c, apierr.InternalError(message))
}
