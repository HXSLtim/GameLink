package handler

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

// MapServiceError maps service layer errors to standard API error responses
func MapServiceError(err error) (int, *model.APIResponse[any]) {
	return MapServiceErrorWithPath(err, "")
}

// MapServiceErrorWithPath maps service layer errors to standard API error responses with path-based messages
func MapServiceErrorWithPath(err error, path string) (int, *model.APIResponse[any]) {
	if err == nil {
		return http.StatusOK, &model.APIResponse[any]{
			Success: true,
			Code:    http.StatusOK,
			Message: "success",
		}
	}

	// Check if it's an APIError type
	if apiErr, ok := err.(*apierr.APIError); ok {
		meta := make(map[string]interface{})
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

		resp := &model.APIResponse[any]{
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
		return http.StatusUnauthorized, &model.APIResponse[any]{
			Success: false,
			Code:    http.StatusUnauthorized,
			Message: "invalid credentials",
		}
	case errors.Is(err, service.ErrUserDisabled):
		return http.StatusForbidden, &model.APIResponse[any]{
			Success: false,
			Code:    http.StatusForbidden,
			Message: "user account is disabled",
		}
	case errors.Is(err, service.ErrOrderInvalidTransition):
		return http.StatusConflict, &model.APIResponse[any]{
			Success: false,
			Code:    http.StatusConflict,
			Message: apierr.ErrOrderInvalidTransition,
		}
	case errors.Is(err, service.ErrUserNotFound):
		return http.StatusNotFound, &model.APIResponse[any]{
			Success: false,
			Code:    http.StatusNotFound,
			Message: apierr.ErrUserNotFound,
		}
	case errors.Is(err, service.ErrNotFound) || errors.Is(err, repository.ErrNotFound):
		message := "resource not found"
		if path != "" {
			message = getDomainNotFoundMessage(path)
		}
		return http.StatusNotFound, &model.APIResponse[any]{
			Success: false,
			Code:    http.StatusNotFound,
			Message: message,
		}
	default:
		return http.StatusInternalServerError, &model.APIResponse[any]{
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

// RespondWithServiceError responds with a standard error format
func RespondWithServiceError(c *gin.Context, err error) {
	RespondWithServiceErrorAndPath(c, err, "")
}

// RespondWithServiceErrorAndPath responds with a standard error format with path-based messages
func RespondWithServiceErrorAndPath(c *gin.Context, err error, path string) {
	statusCode, response := MapServiceErrorWithPath(err, path)

	if requestID := GetRequestID(c); requestID != "" {
		response.TraceID = requestID
	}

	c.JSON(statusCode, response)
}
