package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Recovery returns a custom JSON recovery middleware that outputs
// a unified API envelope on panic.
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": "internal server error",
		})
	})
}
