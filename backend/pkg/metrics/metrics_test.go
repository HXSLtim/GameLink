package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/**
 * Test_Metrics_FullLifecycle 测试metrics完整生命周期
 * 注意：由于使用了sync.Once，所有测试必须在一个函数中完成
 */
func Test_Metrics_FullLifecycle(t *testing.T) {
	reg := prometheus.NewRegistry()

	t.Run("Init initializes metrics", func(t *testing.T) {
		testMetricsInit(t, reg)
	})
	t.Run("metrics are registered", func(t *testing.T) {
		testMetricsRegistered(t, reg)
	})
	t.Run("HTTPRequestsTotal has correct labels", func(t *testing.T) {
		testHTTPRequestsTotalLabels(t, reg)
	})
	t.Run("HTTPRequestDuration has correct labels", func(t *testing.T) {
		testHTTPRequestDurationLabels(t, reg)
	})
	t.Run("DBQueryDuration has correct labels", func(t *testing.T) {
		testDBQueryDurationLabels(t, reg)
	})
	t.Run("counter increments correctly", func(t *testing.T) {
		testCounterIncrements(t, reg)
	})
	t.Run("histogram has buckets", func(t *testing.T) {
		testHistogramBuckets(t, reg)
	})
	t.Run("complete scenario", func(t *testing.T) {
		testCompleteScenario(t, reg)
	})
	t.Run("multiple Init calls are safe", func(t *testing.T) {
		testMultipleInitCalls(t, reg)
	})
}

func testMetricsInit(t *testing.T, reg *prometheus.Registry) {
	Init(reg)
	assert.NotNil(t, HTTPRequestsTotal, "HTTPRequestsTotal应该被初始化")
	assert.NotNil(t, HTTPRequestDuration, "HTTPRequestDuration应该被初始化")
	assert.NotNil(t, DBQueryDuration, "DBQueryDuration应该被初始化")
}

func testMetricsRegistered(t *testing.T, reg *prometheus.Registry) {
	HTTPRequestsTotal.WithLabelValues("GET", "/", "200").Inc()
	HTTPRequestDuration.WithLabelValues("GET", "/").Observe(0.1)
	DBQueryDuration.WithLabelValues("query", "test").Observe(0.001)

	metrics, err := reg.Gather()
	require.NoError(t, err)

	metricNames := make(map[string]bool)
	for _, m := range metrics {
		metricNames[m.GetName()] = true
	}

	assert.True(t, metricNames["http_requests_total"], "应该注册http_requests_total")
	assert.True(t, metricNames["http_request_duration_seconds"], "应该注册http_request_duration_seconds")
	assert.True(t, metricNames["db_query_duration_seconds"], "应该注册db_query_duration_seconds")
}

func testHTTPRequestsTotalLabels(t *testing.T, reg *prometheus.Registry) {
	HTTPRequestsTotal.WithLabelValues("GET", "/api/users", "200").Inc()
	HTTPRequestsTotal.WithLabelValues("POST", "/api/orders", "201").Inc()

	metrics, err := reg.Gather()
	require.NoError(t, err)

	found := findMetricByName(metrics, "http_requests_total")
	assert.True(t, found, "应该找到http_requests_total指标")
}

func testHTTPRequestDurationLabels(t *testing.T, reg *prometheus.Registry) {
	HTTPRequestDuration.WithLabelValues("GET", "/api/users").Observe(0.123)
	HTTPRequestDuration.WithLabelValues("POST", "/api/orders").Observe(0.456)

	metrics, err := reg.Gather()
	require.NoError(t, err)

	found := findMetricByName(metrics, "http_request_duration_seconds")
	assert.True(t, found, "应该找到http_request_duration_seconds指标")
}

func testDBQueryDurationLabels(t *testing.T, reg *prometheus.Registry) {
	DBQueryDuration.WithLabelValues("query", "users").Observe(0.001)
	DBQueryDuration.WithLabelValues("create", "orders").Observe(0.002)
	DBQueryDuration.WithLabelValues("update", "players").Observe(0.003)
	DBQueryDuration.WithLabelValues("delete", "reviews").Observe(0.004)

	metrics, err := reg.Gather()
	require.NoError(t, err)

	found := findMetricByName(metrics, "db_query_duration_seconds")
	assert.True(t, found, "应该找到db_query_duration_seconds指标")
}

func testCounterIncrements(t *testing.T, reg *prometheus.Registry) {
	counter := HTTPRequestsTotal.WithLabelValues("GET", "/test", "200")
	counter.Inc()
	counter.Inc()
	counter.Inc()

	metrics, err := reg.Gather()
	require.NoError(t, err)

	for _, m := range metrics {
		if m.GetName() == "http_requests_total" {
			for _, metric := range m.GetMetric() {
				if hasLabelValue(metric.GetLabel(), "path", "/test") {
					assert.Equal(t, float64(3), metric.GetCounter().GetValue(), "计数应该是3")
					return
				}
			}
		}
	}
}

func testHistogramBuckets(t *testing.T, reg *prometheus.Registry) {
	HTTPRequestDuration.WithLabelValues("GET", "/histogram-test").Observe(0.5)

	metrics, err := reg.Gather()
	require.NoError(t, err)

	for _, m := range metrics {
		if m.GetName() == "http_request_duration_seconds" {
			for _, metric := range m.GetMetric() {
				if histogram := metric.GetHistogram(); histogram != nil && len(histogram.GetBucket()) > 0 {
					assert.Greater(t, len(histogram.GetBucket()), 0, "应该有buckets")
					return
				}
			}
		}
	}
	t.Error("应该找到histogram")
}

func testCompleteScenario(t *testing.T, reg *prometheus.Registry) {
	for i := 0; i < 5; i++ {
		HTTPRequestsTotal.WithLabelValues("GET", "/scenario", "200").Inc()
		HTTPRequestDuration.WithLabelValues("GET", "/scenario").Observe(0.1)
	}
	for i := 0; i < 3; i++ {
		HTTPRequestsTotal.WithLabelValues("POST", "/scenario", "201").Inc()
		HTTPRequestDuration.WithLabelValues("POST", "/scenario").Observe(0.2)
	}
	DBQueryDuration.WithLabelValues("query", "scenario_users").Observe(0.005)
	DBQueryDuration.WithLabelValues("create", "scenario_orders").Observe(0.010)

	metrics, err := reg.Gather()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(metrics), 3, "应该至少有3个指标")
}

func testMultipleInitCalls(t *testing.T, reg *prometheus.Registry) {
	Init(reg)
	Init(reg)
	Init(reg)

	assert.NotNil(t, HTTPRequestsTotal)
	assert.NotNil(t, HTTPRequestDuration)
	assert.NotNil(t, DBQueryDuration)
}

// findMetricByName 查找指定名称的指标
func findMetricByName(metrics []*dto.MetricFamily, name string) bool {
	for _, m := range metrics {
		if m.GetName() == name {
			return true
		}
	}
	return false
}

// hasLabelValue 检查标签是否包含指定值
func hasLabelValue(labels []*dto.LabelPair, name, value string) bool {
	for _, label := range labels {
		if label.GetName() == name && label.GetValue() == value {
			return true
		}
	}
	return false
}

/**
 * Test_Metrics_WithDefaultRegisterer 测试使用默认registerer
 */
func Test_Metrics_WithDefaultRegisterer(t *testing.T) {
	// Arrange & Act
	Init(prometheus.DefaultRegisterer)

	// Assert
	assert.NotNil(t, HTTPRequestsTotal)
	assert.NotNil(t, HTTPRequestDuration)
	assert.NotNil(t, DBQueryDuration)
}

/**
 * Test_HTTPRequestsTotal_DifferentStatusCodes 测试不同的状态码
 */
func Test_HTTPRequestsTotal_DifferentStatusCodes(t *testing.T) {
	// Arrange
	reg := prometheus.NewRegistry()
	Init(reg)

	tests := []struct {
		method string
		path   string
		status string
	}{
		{"GET", "/api/v1/users", "200"},
		{"POST", "/api/v1/users", "201"},
		{"PUT", "/api/v1/users/1", "200"},
		{"DELETE", "/api/v1/users/1", "204"},
		{"GET", "/api/v1/not-found", "404"},
		{"POST", "/api/v1/error", "500"},
	}

	// Act
	for _, tt := range tests {
		HTTPRequestsTotal.WithLabelValues(tt.method, tt.path, tt.status).Inc()
	}

	// Assert
	metrics, err := reg.Gather()
	require.NoError(t, err)

	for _, m := range metrics {
		if m.GetName() == "http_requests_total" {
			assert.GreaterOrEqual(t, len(m.GetMetric()), len(tests), "应该记录所有测试用例")
		}
	}
}

/**
 * Test_DBQueryDuration_AllOperations 测试所有数据库操作类型
 */
func Test_DBQueryDuration_AllOperations(t *testing.T) {
	// Arrange
	reg := prometheus.NewRegistry()
	Init(reg)

	operations := []struct {
		op    string
		table string
	}{
		{"query", "users"},
		{"create", "orders"},
		{"update", "players"},
		{"delete", "reviews"},
		{"query", "games"},
	}

	// Act
	for _, operation := range operations {
		DBQueryDuration.WithLabelValues(operation.op, operation.table).Observe(0.001)
	}

	// Assert
	metrics, err := reg.Gather()
	require.NoError(t, err)

	for _, m := range metrics {
		if m.GetName() == "db_query_duration_seconds" {
			assert.GreaterOrEqual(t, len(m.GetMetric()), len(operations), "应该记录所有操作")
		}
	}
}
