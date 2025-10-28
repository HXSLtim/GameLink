package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"gamelink/internal/metrics"
)

// MetricsMiddleware records HTTP metrics.
func MetricsMiddleware() gin.HandlerFunc {
	// ensure metrics are initialized
	metrics.Init(prometheus.DefaultRegisterer)
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method
		status := strconv.Itoa(c.Writer.Status())
		metrics.HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
		metrics.HTTPRequestDuration.WithLabelValues(method, path).Observe(time.Since(start).Seconds())
	}
}

// MetricsHandler exposes /metrics endpoint using promhttp.DefaultGatherer.
func MetricsHandler() gin.HandlerFunc {
	metrics.Init(prometheus.DefaultRegisterer)
	h := promhttp.Handler()
	return func(c *gin.Context) { h.ServeHTTP(c.Writer, c.Request) }
}
