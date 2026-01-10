# GameLink Performance Benchmarks - Quick Start

This guide helps you quickly get started with performance benchmarking for the GameLink backend.

## Prerequisites

```bash
# Install benchmark tools
make bench-tools
```

## Setup Benchmark Database

```bash
# Option 1: Use Docker (recommended)
docker run -d \
  --name gamelink-bench-db \
  -e POSTGRES_USER=gamelink \
  -e POSTGRES_PASSWORD=gamelink \
  -e POSTGRES_DB=gamelink_bench \
  -p 5432:5432 \
  postgres:16

# Option 2: Use environment variables
export BENCH_DB_HOST=localhost
export BENCH_DB_PORT=5432
export BENCH_DB_USER=gamelink
export BENCH_DB_PASSWORD=gamelink
export BENCH_DB_NAME=gamelink_bench
```

## Run Benchmarks

### All Benchmarks with Profiling

```bash
make bench-all
```

Results saved to: `test/results/bench/`

### Specific Suite

```bash
# Order service
make bench-order

# Authentication
make bench-auth

# Payment processing
make bench-payment

# Database operations
make bench-db

# HTTP endpoints
make bench-http
```

## Run Load Tests

### Start API Server First

```bash
# Terminal 1
go run cmd/main.go
```

### Run Load Tests

```bash
# Terminal 2
make load-test-auth
make load-test-order
```

## Compare Results

```bash
# Save baseline
make bench-order > baseline.txt

# Make code changes

# Run again
make bench-order > current.txt

# Compare
make bench-compare OLD=baseline.txt NEW=current.txt
```

## View Profiles

```bash
# After running make bench-all
go tool pprof -http=:8080 test/results/bench/cpu.prof
# Open browser to http://localhost:8080
```

## CI/CD Integration

Add to your workflow:

```yaml
- name: Run benchmarks
  run: |
    make bench-tools
    make bench-all

- name: Check regressions
  run: |
    make bench-compare OLD=baseline.txt NEW=current.txt
```

## Documentation

For detailed information, see:
- [Performance Benchmarking Guide](../../docs/PERFORMANCE_BENCHMARKING.md)
- [Baseline Metrics](../../docs/BASELINE_METRICS.md)

## Common Issues

### Database connection failed
```bash
# Check database is running
docker ps | grep gamelink-bench-db

# Restart if needed
docker start gamelink-bench-db
```

### Unstable results
```bash
# Run for longer duration
go test -bench=. -benchtime=30s ./internal/service/order/...
```

### Load test connection errors
```bash
# Ensure API server is running
curl http://localhost:8080/api/v1/healthz
```

## Performance SLAs

| Operation | P50 | P95 | P99 |
|-----------|-----|-----|-----|
| Create Order | < 50ms | < 100ms | < 200ms |
| Login | < 30ms | < 60ms | < 100ms |
| Create Payment | < 20ms | < 40ms | < 80ms |

See [BASELINE_METRICS.md](../../docs/BASELINE_METRICS.md) for complete SLAs.

## Next Steps

1. Establish baselines for your environment
2. Set up continuous benchmarking in CI/CD
3. Configure alerts for performance regressions
4. Regularly review and update baselines
