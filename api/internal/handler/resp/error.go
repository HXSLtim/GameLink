package resp

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service"
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

// MapServiceError maps service layer errors to standard API error responses
func MapServiceError(err error) (int, *model.SuccessResponse) {
	return MapServiceErrorWithPath(err, "")
}

// MapServiceErrorWithPath maps service layer errors to standard API error responses with path-based messages
func MapServiceErrorWithPath(err error, path string) (int, *model.SuccessResponse) {
	if err == nil {
		return http.StatusOK, &model.SuccessResponse{
			Success: true,
			Code:    http.StatusOK,
			Message: "success",
		}
	}

	// Check if it's an APIError type
	if apiErr, ok := err.(*apierr.APIError); ok {
		meta := make(map[string]any)
		if apiErr.Details != "" {
			meta["details"] = apiErr.Details
		}
		if apiErr.Field != "" {
			meta["field"] = apiErr.Field
		}
		if apiErr.Timestamp != 0 {
			meta["timestamp"] = apiErr.Timestamp
		}
		if apiErr.Extensions != nil {
			for k, v := range apiErr.Extensions {
				meta[k] = v
			}
		}

		resp := &model.SuccessResponse{
			Success: false,
			Code:    apiErr.Code,
			Message: apiErr.Message,
		}

		if len(meta) > 0 {
			resp.Meta = meta
		}

		return apiErr.Code, resp
	}

	// Check service layer errors
	switch {
	case errors.Is(err, service.ErrInvalidCredentials):
		return http.StatusUnauthorized, &model.SuccessResponse{
			Success: false,
			Code:    http.StatusUnauthorized,
			Message: "invalid credentials",
		}
	case errors.Is(err, service.ErrUserDisabled):
		return http.StatusForbidden, &model.SuccessResponse{
			Success: false,
			Code:    http.StatusForbidden,
			Message: "user account is disabled",
		}
	case errors.Is(err, service.ErrOrderInvalidTransition):
		return http.StatusConflict, &model.SuccessResponse{
			Success: false,
			Code:    http.StatusConflict,
			Message: apierr.ErrOrderInvalidTransition,
		}
	case errors.Is(err, service.ErrUserNotFound):
		return http.StatusNotFound, &model.SuccessResponse{
			Success: false,
			Code:    http.StatusNotFound,
			Message: apierr.ErrUserNotFound,
		}
	case errors.Is(err, service.ErrNotFound) || errors.Is(err, repository.ErrNotFound):
		message := "resource not found"
		if path != "" {
			message = getDomainNotFoundMessage(path)
		}
		return http.StatusNotFound, &model.SuccessResponse{
			Success: false,
			Code:    http.StatusNotFound,
			Message: message,
		}
	default:
		return http.StatusInternalServerError, &model.SuccessResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
}

// getDomainNotFoundMessage returns a stable message for 404 based on route path
func getDomainNotFoundMessage(path string) string {
	p := strings.ToLower(path)
	switch {
	case strings.Contains(p, "/users"):
		return apierr.ErrUserNotFound
	case strings.Contains(p, "/orders"):
		return apierr.ErrOrderNotFound
	case strings.Contains(p, "/payments"):
		return apierr.ErrPaymentNotFound
	case strings.Contains(p, "/players"):
		return apierr.ErrPlayerNotFound
	case strings.Contains(p, "/games"):
		return apierr.ErrGameNotFound
	default:
		return "resource not found"
	}
}
