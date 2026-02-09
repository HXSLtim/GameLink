# Prometheus Metrics Monitoring

This document describes the Prometheus metrics collection system for the GameLink application.

## Overview

GameLink uses Prometheus for comprehensive application monitoring, including HTTP request tracking, slow request detection, database query performance, and business metrics.

## Architecture

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│   HTTP Request  │────▶│  Prometheus      │────▶│  Prometheus     │
│                 │     │  Middleware      │     │  Server         │
└─────────────────┘     └──────────────────┘     └─────────────────┘
                               │
                               ▼
                        ┌──────────────────┐
                        │  Metrics         │
                        │  Collector       │
                        └──────────────────┘
                               │
                    ┌──────────┼──────────┐
                    ▼          ▼          ▼
              ┌─────────┐ ┌────────┐ ┌────────────┐
              │ HTTP    │ │  DB    │ │  Business  │
              │ Metrics │ │ Metrics│ │  Metrics   │
              └─────────┘ └────────┘ └────────────┘
```

## Metrics Collected

### HTTP Metrics

| Metric Name | Type | Labels | Description |
|------------|------|--------|-------------|
| `http_requests_total` | Counter | method, path, status | Total number of HTTP requests |
| `http_request_duration_seconds` | Histogram | method, path | HTTP request latency in seconds |
| `http_slow_requests_total` | Counter | method, path | Requests taking >1 second |

**Buckets for duration histogram**:
```
0.005s, 0.01s, 0.025s, 0.05s, 0.1s, 0.25s, 0.5s, 1.0s, 2.5s, 5.0s, 10.0s
```

### Database Metrics

| Metric Name | Type | Labels | Description |
|------------|------|--------|-------------|
| `db_query_duration_seconds` | Histogram | operation, table | Database query latency in seconds |

**Operations tracked**: `query`, `create`, `update`, `delete`

**Buckets for query duration**:
```
0.001s, 0.005s, 0.01s, 0.025s, 0.05s, 0.1s, 0.25s, 0.5s, 1.0s, 2.5s
```

### Business Metrics

| Metric Name | Type | Labels | Description |
|------------|------|--------|-------------|
| `gamelink_orders_created_total` | Counter | status, game_type | Total orders created |
| `gamelink_orders_completed_total` | Counter | game_type, player_tier | Total orders completed |
| `gamelink_orders_canceled_total` | Counter | reason, canceled_by | Total orders canceled |
| `gamelink_order_duration_hours` | Histogram | game_type | Order duration in hours |
| `gamelink_payments_succeeded_total` | Counter | method, currency | Successful payments |
| `gamelink_payments_failed_total` | Counter | method, failure_reason | Failed payments |
| `gamelink_users_active` | Gauge | role, status | Currently active users |
| `gamelink_commission_total_cents` | Counter | commission_type | Total commission earned |

## Configuration

### Environment Variables

| Variable | Description | Example | Default |
|----------|-------------|---------|---------|
| `METRICS_ENABLED` | Enable/disable metrics endpoint | `true` | `true` |
| `METRICS_ALLOWED_CIDRS` | Allowed CIDR blocks for metrics access | `10.0.0.0/8,192.168.0.0/16` | `localhost only` |
| `APP_ENV` | Application environment | `production` | `development` |

### Security

The `/metrics` endpoint is secured by default:

- **Production**: Only localhost (127.0.0.0/8) and configured CIDRs can access
- **Development/Staging**: Open access (authentication disabled)

Configure allowed networks via environment variable:

```bash
# Allow private network access
export METRICS_ALLOWED_CIDRS="10.0.0.0/8,192.168.0.0/16,172.16.0.0/12"

# Allow specific VPN subnet
export METRICS_ALLOWED_CIDRS="10.8.0.0/16"
```

## Usage

### Starting the Application

Metrics are automatically initialized on startup:

```bash
# Development (metrics endpoint open)
go run cmd/main.go

# Production (metrics endpoint secured)
APP_ENV=production METRICS_ALLOWED_CIDRS="10.0.0.0/8" go run cmd/main.go
```

### Accessing Metrics

```bash
# Direct access (development)
curl http://localhost:8080/metrics

# With authentication (production)
curl http://localhost:8080/metrics
# Returns 403 Forbidden if not from allowed IP
```

### Example Metrics Output

```
# HTTP requests
http_requests_total{method="GET",path="/api/v1/users",status="200"} 1523
http_requests_total{method="POST",path="/api/v1/orders",status="201"} 342

# Request duration
http_request_duration_seconds_bucket{method="GET",path="/api/v1/users",le="0.1"} 1500
http_request_duration_seconds_bucket{method="GET",path="/api/v1/users",le="0.5"} 1520
http_request_duration_seconds_sum{method="GET",path="/api/v1/users"} 89.5
http_request_duration_seconds_count{method="GET",path="/api/v1/users"} 1523

# Slow requests
http_slow_requests_total{method="POST",path="/api/v1/orders"} 15

# Database queries
db_query_duration_seconds_bucket{operation="query",table="users",le="0.01"} 5000
db_query_duration_seconds_sum{operation="query",table="users"} 25.3

# Business metrics
gamelink_orders_created_total{status="created",game_type="王者荣耀"} 450
gamelink_payments_succeeded_total{method="alipay",currency="CNY"} 320
```

## Slow Request Detection

Requests taking longer than **1 second** are automatically:

1. **Counted** in the `http_slow_requests_total` metric
2. **Logged** with `slog.LevelWarn` including:
   - Request duration
   - Method and path
   - Status code
   - Client IP
   - Request ID (if available)
   - User ID (if available)

Example slow request log:

```json
{
  "level": "WARN",
  "msg": "slow_request_detected",
  "duration_seconds": 1.523,
  "method": "POST",
  "path": "/api/v1/orders",
  "status": 201,
  "ip": "192.168.1.100",
  "request_id": "abc123",
  "user_id": 456
}
```

## Prometheus Configuration

### prometheus.yml

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'gamelink-api'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 10s
```

### Docker Compose (Prometheus + Grafana)

```yaml
version: '3.8'

services:
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus-data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana-data:/var/lib/grafana

volumes:
  prometheus-data:
  grafana-data:
```

## Grafana Dashboards

### Recommended Queries

**Request Rate** (requests per second):
```promql
rate(http_requests_total[5m])
```

**Error Rate** (percentage):
```promql(
  sum(rate(http_requests_total{status=~"5.."}[5m])) /
  sum(rate(http_requests_total[5m])) * 100
)
```

**P95 Latency**:
```promql
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))
```

**Slow Requests** (requests per second):
```promql
rate(http_slow_requests_total[5m])
```

**Database Query Performance**:
```promql
rate(db_query_duration_seconds_sum[5m]) / rate(db_query_duration_seconds_count[5m])
```

## Performance Considerations

### Metrics Overhead

- **Per-request overhead**: ~50-100 microseconds
- **Memory footprint**: ~1-2 MB for 10k metrics
- **No impact** when metrics endpoint is not being scraped

### Best Practices

1. **Scrape interval**: Use 10-15 seconds for production
2. **Retention**: Keep 15-30 days of metrics data
3. **Labels**: Avoid high-cardinality labels (e.g., user IDs)
4. **Aggregation**: Aggregate in Prometheus, not in application

## Troubleshooting

### Metrics Endpoint Returns 403

**Problem**: Access to `/metrics` is forbidden

**Solution**:
```bash
# Check if you're accessing from allowed IP
curl http://localhost:8080/metrics  # Should work from localhost

# Configure allowed CIDRs for remote access
export METRICS_ALLOWED_CIDRS="10.0.0.0/8,192.168.0.0/16"
```

### No Metrics Appearing

**Problem**: Metrics are not being recorded

**Solutions**:
1. Check if collector is initialized: `metrics.NewCollector(registry)`
2. Verify middleware is registered: `router.Use(middleware.PrometheusMiddleware())`
3. Check logs for initialization errors

### High Memory Usage

**Problem**: Metrics consuming too much memory

**Solutions**:
1. Reduce histogram bucket counts
2. Avoid high-cardinality labels
3. Increase scrape interval to reduce metric churn
4. Use metric relabeling in Prometheus

## Testing

### Unit Tests

```bash
# Run all metrics tests
go test ./pkg/metrics/... -v
go test ./internal/handler/middleware/*prometheus* -v

# Run with coverage
go test ./pkg/metrics/... -cover -coverprofile=coverage.out
```

### Integration Tests

```bash
# Start test server
go run cmd/main.go

# Scrape metrics
curl http://localhost:8080/metrics | grep http_requests_total

# Generate load
ab -n 1000 -c 10 http://localhost:8080/api/v1/health

# Verify metrics
curl http://localhost:8080/metrics | grep http_requests_total
```

## Alerting Rules

### Example Prometheus Alerts

```yaml
groups:
  - name: gamelink_alerts
    interval: 30s
    rules:
      # High error rate
      - alert: HighErrorRate
        expr: |
          sum(rate(http_requests_total{status=~"5.."}[5m])) /
          sum(rate(http_requests_total[5m])) > 0.05
        for: 5m
        annotations:
          summary: "High error rate detected"

      # Slow requests
      - alert: SlowRequests
        expr: |
          rate(http_slow_requests_total[5m]) > 10
        for: 5m
        annotations:
          summary: "Too many slow requests"

      # High latency
      - alert: HighLatency
        expr: |
          histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        annotations:
          summary: "P95 latency exceeds 1 second"
```

## References

- [Prometheus Go Client Documentation](https://prometheus.io/docs/guides/go-application/)
- [Prometheus Best Practices](https://prometheus.io/docs/practices/naming/)
- [Histogram Quantiles](https://prometheus.io/docs/practices/histograms/)
