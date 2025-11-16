package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	apierr "gamelink/internal/handler"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service"
)

// ErrorMap inspects gin errors and maps known errors to standard envelope responses.
// Only applies when handler hasn't already written a response.
func ErrorMap() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Writer.Written() {
			return
		}
		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors[0].Err
		msg := err.Error()

		switch {
		case errors.Is(err, service.ErrValidation) ||
			msg == apierr.ErrInvalidJSONPayload ||
			msg == apierr.ErrMissingRequiredFields ||
			msg == apierr.ErrMissingFieldsOrShortPassword:
			c.JSON(http.StatusBadRequest, model.APIResponse[any]{Success: false, Code: http.StatusBadRequest, Message: "validation failed"})
			return
		case errors.Is(err, service.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, model.APIResponse[any]{Success: false, Code: http.StatusUnauthorized, Message: "invalid credentials"})
			return
		case errors.Is(err, service.ErrUserDisabled):
			c.JSON(http.StatusForbidden, model.APIResponse[any]{Success: false, Code: http.StatusForbidden, Message: "user account is disabled"})
			return
		case errors.Is(err, service.ErrOrderInvalidTransition) || msg == apierr.ErrOrderInvalidTransition:
			c.JSON(http.StatusConflict, model.APIResponse[any]{Success: false, Code: http.StatusConflict, Message: apierr.ErrOrderInvalidTransition})
			return
		case errors.Is(err, service.ErrUserNotFound) || msg == apierr.ErrUserNotFound:
			c.JSON(http.StatusNotFound, model.APIResponse[any]{Success: false, Code: http.StatusNotFound, Message: apierr.ErrUserNotFound})
			return
		case msg == apierr.ErrInvalidID ||
			msg == apierr.ErrInvalidUserID ||
			msg == apierr.ErrInvalidOrderID ||
			msg == apierr.ErrInvalidPlayerID ||
			msg == apierr.ErrInvalidGameID:
			c.JSON(http.StatusBadRequest, model.APIResponse[any]{Success: false, Code: http.StatusBadRequest, Message: "invalid id format"})
			return
		case msg == apierr.ErrInvalidPage || msg == apierr.ErrInvalidPageSize:
			c.JSON(http.StatusBadRequest, model.APIResponse[any]{Success: false, Code: http.StatusBadRequest, Message: "invalid pagination parameters"})
			return
		case msg == apierr.ErrInvalidEmailFormat:
			c.JSON(http.StatusBadRequest, model.APIResponse[any]{Success: false, Code: http.StatusBadRequest, Message: "invalid email format"})
			return
		case msg == apierr.ErrInvalidPhoneFormat:
			c.JSON(http.StatusBadRequest, model.APIResponse[any]{Success: false, Code: http.StatusBadRequest, Message: "invalid phone format"})
			return
		case errors.Is(err, service.ErrNotFound) || errors.Is(err, repository.ErrNotFound):
			// Return domain-specific not found messages based on route path
			path := c.FullPath()
			if path == "" {
				path = c.Request.URL.Path
			}
			msg := domainNotFoundMessage(path)
			c.JSON(http.StatusNotFound, model.APIResponse[any]{Success: false, Code: http.StatusNotFound, Message: msg})
			return
		default:
			c.JSON(http.StatusInternalServerError, model.APIResponse[any]{Success: false, Code: http.StatusInternalServerError, Message: "internal server error"})
			return
		}
	}
}

// domainNotFoundMessage returns a stable message for 404 based on route path.
func domainNotFoundMessage(path string) string {
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
		return "not found"
	}
}
