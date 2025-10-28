package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// SlogLogger is a structured HTTP access logger based on slog.
func SlogLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		dur := time.Since(start)

		status := c.Writer.Status()
		level := slog.LevelInfo
		switch {
		case status >= 500:
			level = slog.LevelError
		case status >= 400:
			level = slog.LevelWarn
		}
		attrs := []slog.Attr{
			slog.Int("status", status),
			slog.String("method", c.Request.Method),
			slog.String("path", c.FullPath()),
			slog.String("ip", c.ClientIP()),
			slog.String("duration", dur.String()),
		}
		if rid, ok := c.Get("request_id"); ok {
			if s, ok2 := rid.(string); ok2 {
				attrs = append(attrs, slog.String("request_id", s))
			}
		}
		if uid, ok := c.Get("user_id"); ok {
			attrs = append(attrs, slog.Any("user_id", uid))
		}
		slog.LogAttrs(c.Request.Context(), level, "http_request", attrs...)
	}
}
