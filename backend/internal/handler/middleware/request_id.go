package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"

	"gamelink/internal/logging"
)

// RequestID ensures every request carries an X-Request-ID, setting it on the
// response and injecting into context for downstream logging.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader("X-Request-ID")
		if rid == "" {
			rid = randomID()
		}
		c.Header("X-Request-ID", rid)
		c.Set("request_id", rid)
		c.Request = c.Request.WithContext(logging.WithRequestID(c.Request.Context(), rid))
		c.Next()
	}
}

func randomID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return ""
	}
	return hex.EncodeToString(b[:])
}
