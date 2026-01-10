# Prometheus Metrics Quick Start

## Installation & Setup

The metrics system is already integrated into GameLink. No additional installation required.

## Quick Start

### 1. Start the Application

```bash
cd api

# Development mode (metrics endpoint open)
go run cmd/main.go

# Production mode (metrics endpoint secured)
APP_ENV=production METRICS_ALLOWED_CIDRS="10.0.0.0/8" go run cmd/main.go
```

### 2. Access Metrics

```bash
# From localhost (always works)
curl http://localhost:8080/metrics

# View specific metric
curl http://localhost:8080/metrics | grep http_requests_total
```

### 3. Setup Prometheus

Create `prometheus.yml`:

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'gamelink'
    static_configs:
      - targets: ['localhost:8080']
```

Run Prometheus:

```bash
docker run -d \
  -p 9090:9090 \
  -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
  prom/prometheus
```

Access at http://localhost:9090

### 4. Setup Grafana

```bash
docker run -d \
  -p 3000:3000 \
  -e "GF_SECURITY_ADMIN_PASSWORD=admin" \
  grafana/grafana
```

Access at http://localhost:3000 (admin/admin)

Add Prometheus datasource: http://localhost:9090

## Key Metrics

### HTTP Performance

```promql
# Request rate (requests/sec)
rate(http_requests_total[5m])

# P95 latency
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Error rate (%)
sum(rate(http_requests_total{status=~"5.."}[5m])) /
sum(rate(http_requests_total[5m])) * 100
```

### Slow Requests

```promql
# Slow requests per second
rate(http_slow_requests_total[5m])

# Total slow requests
http_slow_requests_total
```

### Database Performance

```promql
# Average query duration
rate(db_query_duration_seconds_sum[5m]) /
rate(db_query_duration_seconds_count[5m])

# Slow queries (>100ms)
rate(db_query_duration_seconds_bucket{le="0.1"}[5m])
```

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `APP_ENV` | Environment (dev/production) | `production` |
| `METRICS_ALLOWED_CIDRS` | Allowed networks for /metrics | `10.0.0.0/8,192.168.0.0/16` |

## Testing

```bash
# Run unit tests
go test ./pkg/metrics/... -v
go test ./internal/handler/middleware/... -v -run Prometheus

# Generate load
ab -n 1000 -c 10 http://localhost:8080/api/v1/health

# Check metrics
curl http://localhost:8080/metrics | grep http_requests
```

## Troubleshooting

**403 Forbidden on /metrics?**
```bash
# Check you're accessing from allowed IP
curl http://localhost:8080/metrics  # Works from localhost

# Set allowed networks if needed
export METRICS_ALLOWED_CIDRS="10.0.0.0/8"
```

**No metrics appearing?**
- Check application logs for initialization
- Verify Prometheus is scraping correct port
- Check firewall rules

## Documentation

- Full documentation: `docs/PROMETHEUS_METRICS.md`
- Implementation details: `docs/METRICS_IMPLEMENTATION.md`
