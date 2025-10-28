package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	once sync.Once

	// HTTPRequestsTotal counts requests by method, path and status code.
	HTTPRequestsTotal *prometheus.CounterVec

	// HTTPRequestDuration measures request duration seconds by method and path.
	HTTPRequestDuration *prometheus.HistogramVec

	// DBQueryDuration measures gorm operation duration seconds by op (query/create/update/delete) and table.
	DBQueryDuration *prometheus.HistogramVec
)

// Init registers metrics. Safe to call multiple times.
func Init(reg prometheus.Registerer) {
	once.Do(func() {
		HTTPRequestsTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		)
		HTTPRequestDuration = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		)
		DBQueryDuration = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "db_query_duration_seconds",
				Help:    "GORM operation duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"op", "table"},
		)
		reg.MustRegister(HTTPRequestsTotal, HTTPRequestDuration, DBQueryDuration)
	})
}
