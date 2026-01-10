package middleware

import (
	"log/slog"
	"strconv"
	"time"

	monitorservice "gamelink/internal/service/monitor"
	"gamelink/pkg/metrics"

	"github.com/gin-gonic/gin"
)

const (
	// SlowRequestThreshold defines the duration (in seconds) above which a request is considered slow
	SlowRequestThreshold = 1.0
)

// PrometheusMiddleware returns a gin middleware that records HTTP metrics using Prometheus
// It tracks:
// - Total request count by method, path, and status
// - Request duration histogram by method and path
// - Slow request counter for requests taking longer than SlowRequestThreshold
// - Logs warnings for slow requests
// Note: The metrics collector must be initialized before this middleware is used
func PrometheusMiddleware(monitorSvc *monitorservice.RealtimeService) gin.HandlerFunc {
	collector := metrics.GetCollector()
	if collector == nil {
		// Collector not initialized - return a no-op middleware to avoid panics
		// In production, ensure metrics.NewCollector is called before setting up routes
		return func(c *gin.Context) {
			if monitorSvc != nil {
				monitorSvc.IncrementRequestCount()
			}
			c.Next()
		}
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		// Track active connections
		collector.IncrementActiveConnections("http")
		defer collector.DecrementActiveConnections("http")

		// Record request count for calculating requests per second (legacy support)
		if monitorSvc != nil {
			monitorSvc.IncrementRequestCount()
		}

		// Process request
		c.Next()

		// Calculate metrics after request completes
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method

		// Record HTTP metrics using the collector
		collector.RecordHTTPRequest(method, path, status, duration)

		// Log warning for slow requests
		if duration > SlowRequestThreshold {
			logSlowRequest(c, path, duration)
		}

		// Record error metrics for 4xx and 5xx status codes
		if c.Writer.Status() >= 400 && metrics.BusinessMetrics != nil {
			recordErrorMetrics(c, path, method)
		}
	}
}

// logSlowRequest logs a warning for slow requests
func logSlowRequest(c *gin.Context, path string, duration float64) {
	attrs := []slog.Attr{
		slog.Float64("duration_seconds", duration),
		slog.String("method", c.Request.Method),
		slog.String("path", path),
		slog.Int("status", c.Writer.Status()),
		slog.String("ip", c.ClientIP()),
	}

	// Add request ID if available
	if rid, exists := c.Get("request_id"); exists {
		if ridStr, ok := rid.(string); ok {
			attrs = append(attrs, slog.String("request_id", ridStr))
		}
	}

	// Add user ID if available
	if uid, exists := c.Get("user_id"); exists {
		attrs = append(attrs, slog.Any("user_id", uid))
	}

	slog.LogAttrs(c.Request.Context(), slog.LevelWarn, "slow_request_detected", attrs...)
}

// recordErrorMetrics records error metrics for 4xx and 5xx responses
func recordErrorMetrics(c *gin.Context, path, method string) {
	errorType := "client_error"
	if c.Writer.Status() >= 500 {
		errorType = "server_error"
	}

	metrics.BusinessMetrics.ErrorsTotal.WithLabelValues(
		errorType,
		path,
		method,
	).Inc()
}
