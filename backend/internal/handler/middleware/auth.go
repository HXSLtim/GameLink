package middleware

import (
    "net/http"
    "os"
    "strconv"

    "github.com/gin-gonic/gin"

    "gamelink/internal/logging"
)

// AdminAuth enforces a simple bearer token for admin endpoints.
// Behavior:
// - If APP_ENV=production and ADMIN_TOKEN is empty -> reject all with 503 to avoid unsafe exposure.
// - If ADMIN_TOKEN is set -> require Authorization: Bearer <ADMIN_TOKEN>.
// - Otherwise (development) -> pass-through.
func AdminAuth() gin.HandlerFunc {
	env := os.Getenv("APP_ENV")
	token := os.Getenv("ADMIN_TOKEN")

	// If production and no token configured, block access.
	if env == "production" && token == "" {
		return func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"success": false,
				"code":    http.StatusServiceUnavailable,
				"message": "admin auth not configured",
			})
		}
	}

    // If no token configured (development), allow requests and optional actor injection
    if token == "" {
        return func(c *gin.Context) {
            // In non-production, allow an optional X-Admin-User-ID to mark actor for audit.
            if os.Getenv("APP_ENV") != "production" {
                if v := c.GetHeader("X-Admin-User-ID"); v != "" {
                    if id, err := strconv.ParseUint(v, 10, 64); err == nil && id != 0 {
                        c.Request = c.Request.WithContext(logging.WithActorUserID(c.Request.Context(), id))
                    }
                }
            }
            c.Next()
        }
    }

	// Enforce bearer token
	prefix := "Bearer "
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth != prefix+token {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"message": "unauthorized",
			})
			return
		}
		c.Next()
	}
}
