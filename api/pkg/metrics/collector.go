package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// Collector manages all Prometheus metrics for the GameLink application
type Collector struct {
	// HTTP metrics
	HTTPRequestsTotal     *prometheus.CounterVec
	HTTPRequestDuration   *prometheus.HistogramVec
	HTTPSlowRequestsTotal *prometheus.CounterVec

	// Database metrics
	DBQueryDuration *prometheus.HistogramVec

	// System metrics
	ActiveConnections *prometheus.GaugeVec
}

var (
	collectorOnce sync.Once
	collector     *Collector
)

// NewCollector creates and initializes the metrics collector
// Safe to call multiple times - will return the same instance
func NewCollector(reg prometheus.Registerer) *Collector {
	if reg == nil {
		reg = prometheus.DefaultRegisterer
	}

	collectorOnce.Do(func() {
		collector = &Collector{
			// HTTP metrics
			HTTPRequestsTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: "http_requests_total",
					Help: "Total number of HTTP requests by method, path, and status",
				},
				[]string{"method", "path", "status"},
			),

			HTTPRequestDuration: prometheus.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "http_request_duration_seconds",
					Help:    "HTTP request latency in seconds by method and path",
					Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0},
				},
				[]string{"method", "path"},
			),

			HTTPSlowRequestsTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: "http_slow_requests_total",
					Help: "Total number of slow HTTP requests (>1s) by method and path",
				},
				[]string{"method", "path"},
			),

			// Database metrics
			DBQueryDuration: prometheus.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "db_query_duration_seconds",
					Help:    "Database query latency in seconds by operation and table",
					Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5},
				},
				[]string{"operation", "table"},
			),

			// System metrics
			ActiveConnections: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "active_connections",
					Help: "Current number of active connections by type",
				},
				[]string{"type"},
			),
		}

		// Register all metrics
		registerMetrics(reg, []prometheus.Collector{
			collector.HTTPRequestsTotal,
			collector.HTTPRequestDuration,
			collector.HTTPSlowRequestsTotal,
			collector.DBQueryDuration,
			collector.ActiveConnections,
		})
	})

	return collector
}

// GetCollector returns the initialized collector instance
func GetCollector() *Collector {
	return collector
}

// registerMetrics safely registers Prometheus collectors
func registerMetrics(reg prometheus.Registerer, collectors []prometheus.Collector) {
	for _, c := range collectors {
		if c == nil {
			continue
		}
		if err := reg.Register(c); err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				panic(err)
			}
		}
	}
}

// RecordHTTPRequest records an HTTP request with method, path, status, and duration
func (c *Collector) RecordHTTPRequest(method, path, status string, duration float64) {
	if c == nil {
		return
	}

	// Increment total request counter
	if c.HTTPRequestsTotal != nil {
		c.HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
	}

	// Record request duration
	if c.HTTPRequestDuration != nil {
		c.HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)
	}

	// Track slow requests (>1 second)
	if duration > 1.0 && c.HTTPSlowRequestsTotal != nil {
		c.HTTPSlowRequestsTotal.WithLabelValues(method, path).Inc()
	}
}

// RecordDBQuery records a database query with operation, table, and duration
func (c *Collector) RecordDBQuery(operation, table string, duration float64) {
	if c == nil || c.DBQueryDuration == nil {
		return
	}
	c.DBQueryDuration.WithLabelValues(operation, table).Observe(duration)
}

// SetActiveConnections sets the current number of active connections
func (c *Collector) SetActiveConnections(connType string, value float64) {
	if c == nil || c.ActiveConnections == nil {
		return
	}
	c.ActiveConnections.WithLabelValues(connType).Set(value)
}

// IncrementActiveConnections increments the active connection count
func (c *Collector) IncrementActiveConnections(connType string) {
	if c == nil || c.ActiveConnections == nil {
		return
	}
	c.ActiveConnections.WithLabelValues(connType).Inc()
}

// DecrementActiveConnections decrements the active connection count
func (c *Collector) DecrementActiveConnections(connType string) {
	if c == nil || c.ActiveConnections == nil {
		return
	}
	c.ActiveConnections.WithLabelValues(connType).Dec()
}
