package middleware

import (
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
        switch err {
        case service.ErrValidation:
            c.JSON(http.StatusBadRequest, model.APIResponse[any]{Success: false, Code: http.StatusBadRequest, Message: "validation failed"})
            return
        case service.ErrNotFound, repository.ErrNotFound:
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
