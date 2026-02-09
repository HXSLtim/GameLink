# Performance Benchmarking Guide for GameLink Backend

This guide provides comprehensive information about performance benchmarking for the GameLink Go backend, including setup, execution, result interpretation, and SLA definitions.

## Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Quick Start](#quick-start)
4. [Benchmark Framework](#benchmark-framework)
5. [Benchmark Suites](#benchmark-suites)
6. [Load Testing](#load-testing)
7. [Performance SLAs](#performance-slas)
8. [Interpreting Results](#interpreting-results)
9. [Continuous Benchmarking](#continuous-benchmarking)
10. [Troubleshooting](#troubleshooting)

---

## Overview

The GameLink benchmarking framework provides:

- **Micro-benchmarks**: Fine-grained performance testing of individual functions and methods
- **HTTP benchmarks**: End-to-end API endpoint performance testing
- **Load testing**: Stress testing with simulated traffic using Vegeta
- **Database benchmarks**: Raw database operation performance
- **Regression detection**: Automated performance regression tracking

### Key Features

- **Baseline metrics**: Established performance baselines for critical operations
- **Automated comparison**: Built-in tools to compare results across runs
- **Profiling integration**: CPU and memory profiling support
- **CI/CD ready**: Can be integrated into continuous integration pipelines

---

## Prerequisites

### Required Tools

```bash
# Install benchmark and load testing tools
make bench-tools

# This installs:
# - benchstat: For comparing benchmark results
# - vegeta: For HTTP load testing
```

### Database Setup

Benchmarks require a dedicated PostgreSQL database:

```bash
# Option 1: Use environment variables (recommended)
export BENCH_DB_HOST=localhost
export BENCH_DB_PORT=5432
export BENCH_DB_USER=gamelink
export BENCH_DB_PASSWORD=gamelink
export BENCH_DB_NAME=gamelink_bench

# Option 2: Use Docker
docker run -d \
  --name gamelink-bench-db \
  -e POSTGRES_USER=gamelink \
  -e POSTGRES_PASSWORD=gamelink \
  -e POSTGRES_DB=gamelink_bench \
  -p 5432:5432 \
  postgres:16
```

---

## Quick Start

### 1. Run All Benchmarks

```bash
# Run all benchmarks with profiling
make bench-all

# Results will be saved to test/results/bench/
```

### 2. Run Specific Suite

```bash
# Order service benchmarks
make bench-order

# Authentication benchmarks
make bench-auth

# Database benchmarks
make bench-db
```

### 3. Run Load Tests

```bash
# Start the API server first
go run cmd/main.go

# In another terminal, run load tests
make load-test-auth
```

---

## Benchmark Framework

### Framework Components

Located in `api/internal/benchmark/framework/`:

- **framework.go**: Core benchmark infrastructure
  - `BenchmarkSuite`: Container for services and repositories
  - `BenchmarkConfig`: Configuration management
  - Helper functions for data creation and cleanup

### Key Features

```go
// Create a benchmark suite
suite, orderService := setupOrderBenchmarkSuite(b)

// Create test data helpers
user := suite.CreateBenchmarkUser(b, "13800138000")
player := suite.CreateBenchmarkPlayer(b, user.ID, "Test Player")

// Automatic cleanup
suite.CleanBenchmarkData(b)
suite.ResetBenchmarkSequence(b)
```

### Configuration

```go
type BenchmarkConfig struct {
    DBHost          string  // Database host
    DBPort          string  // Database port
    DBUser          string  // Database user
    DBPassword      string  // Database password
    DBName          string  // Database name
    EnableProfiling bool    // Enable CPU/memory profiling
    CPUProfile      string  // CPU profile output path
    MemProfile      string  // Memory profile output path
}
```

---

## Benchmark Suites

### 1. Order Service Benchmarks

**Location**: `api/internal/service/order/order_bench_test.go`

#### Tests

| Benchmark Name | Description | Key Metrics |
|---------------|-------------|-------------|
| `BenchmarkOrderCreation_Simple` | Basic order creation | ns/op, B/op, allocs/op |
| `BenchmarkOrderListing` | List orders with pagination | Throughput, latency |
| `BenchmarkOrderGetByID` | Retrieve single order | Cache hit ratio |
| `BenchmarkOrderStatusUpdate` | Update order status | Transaction time |
| `BenchmarkOrderCancellation` | Cancel order | Rollback performance |
| `BenchmarkOrderComplexQuery` | Complex queries with joins | Query optimization |
| `BenchmarkOrderConcurrentCreation` | Concurrent order creation | Lock contention |

#### Running

```bash
make bench-order

# Or with Go directly
go test -v -bench=. -benchtime=10s -run=^$ ./internal/service/order/... -benchmem
```

#### Expected Baseline Metrics

| Operation | Target P50 | Target P95 | Target P99 |
|-----------|-----------|-----------|-----------|
| Create Order | < 50ms | < 100ms | < 200ms |
| List Orders (20 items) | < 30ms | < 60ms | < 100ms |
| Get Order by ID | < 10ms | < 20ms | < 40ms |
| Update Status | < 20ms | < 40ms | < 80ms |
| Cancel Order | < 50ms | < 100ms | < 200ms |

---

### 2. Authentication Benchmarks

**Location**: `api/internal/service/auth/auth_bench_test.go`

#### Tests

| Benchmark Name | Description |
|---------------|-------------|
| `BenchmarkLogin` | User login with password verification |
| `BenchmarkTokenGeneration` | JWT token generation |
| `BenchmarkTokenVerification` | JWT token verification |
| `BenchmarkMe` | /auth/me endpoint (token verify + user fetch) |
| `BenchmarkRegister` | User registration |
| `BenchmarkConcurrentLogin` | Concurrent login attempts |

#### Expected Baseline Metrics

| Operation | Target P50 | Target P95 | Target P99 |
|-----------|-----------|-----------|-----------|
| Login | < 30ms | < 60ms | < 100ms |
| Token Generation | < 1ms | < 2ms | < 5ms |
| Token Verification | < 1ms | < 2ms | < 5ms |
| Register | < 50ms | < 100ms | < 200ms |

---

### 3. Payment Service Benchmarks

**Location**: `api/internal/service/payment/payment_bench_test.go`

#### Tests

| Benchmark Name | Description |
|---------------|-------------|
| `BenchmarkCreatePayment` | Create payment record |
| `BenchmarkPaymentStatusUpdate` | Update payment status |
| `BenchmarkGetPaymentByID` | Retrieve payment by ID |
| `BenchmarkGetPaymentsByOrder` | List payments for an order |
| `BenchmarkPaymentRefund` | Process refund |
| `BenchmarkPaymentStatistics` | Calculate payment statistics |

#### Expected Baseline Metrics

| Operation | Target P50 | Target P95 | Target P99 |
|-----------|-----------|-----------|-----------|
| Create Payment | < 20ms | < 40ms | < 80ms |
| Update Status | < 15ms | < 30ms | < 60ms |
| Get Payment by ID | < 10ms | < 20ms | < 40ms |
| Refund | < 100ms | < 200ms | < 400ms |

---

### 4. Database Benchmarks

**Location**: `api/internal/repository/benchmarks/database_bench_test.go`

#### Tests

| Benchmark Name | Description |
|---------------|-------------|
| `BenchmarkDBInsert_Single` | Single row insert |
| `BenchmarkDBInsert_Batch` | Batch insert (100 rows) |
| `BenchmarkDBSelectByPrimaryKey` | Primary key lookup |
| `BenchmarkDBSelectByIndex` | Indexed column lookup |
| `BenchmarkDBUpdate` | Update operation |
| `BenchmarkDBJoin_Simple` | Simple 2-table join |
| `BenchmarkDBJoin_Complex` | Complex 4-table join |
| `BenchmarkDBTransaction` | Transaction with multiple operations |

#### Expected Baseline Metrics

| Operation | Target P50 | Target P95 | Target P99 |
|-----------|-----------|-----------|-----------|
| Single Insert | < 5ms | < 10ms | < 20ms |
| Batch Insert (100) | < 50ms | < 100ms | < 200ms |
| Primary Key Lookup | < 1ms | < 2ms | < 5ms |
| Indexed Lookup | < 2ms | < 5ms | < 10ms |
| Simple Join | < 10ms | < 20ms | < 40ms |

---

### 5. HTTP Endpoint Benchmarks

**Location**: `api/internal/handler/benchmarks/http_bench_test.go`

#### Tests

| Benchmark Name | Description |
|---------------|-------------|
| `BenchmarkHTTP_Login` | POST /api/v1/auth/login |
| `BenchmarkHTTP_Register` | POST /api/v1/auth/register |
| `BenchmarkHTTP_Me` | GET /api/v1/auth/me (with auth) |
| `BenchmarkHTTP_ConcurrentRequests` | Concurrent HTTP requests |
| `BenchmarkHTTP_Middleware` | Middleware overhead |
| `BenchmarkHTTP_FullRequestResponse` | Full HTTP cycle |

#### Expected Baseline Metrics

| Operation | Target P50 | Target P95 | Target P99 |
|-----------|-----------|-----------|-----------|
| POST /auth/login | < 40ms | < 80ms | < 150ms |
| POST /auth/register | < 60ms | < 120ms | < 250ms |
| GET /auth/me | < 20ms | < 40ms | < 80ms |

---

## Load Testing

### Vegeta Load Testing

**Location**: `api/tests/load/vegeta/`

#### Setup

1. Start the API server:
```bash
go run cmd/main.go
```

2. Configure targets in `vegeta/auth_targets.txt`:
```txt
POST http://localhost:8080/api/v1/auth/login
Content-Type: application/json

{"username":"13800138000","password":"password123"}
```

3. Run load test:
```bash
make load-test-auth
```

#### Scenarios

| Scenario | Rate | Duration | Purpose |
|----------|------|----------|---------|
| Auth endpoints | 100 req/s | 30s | Baseline auth performance |
| Order endpoints | 50 req/s | 60s | Order creation and listing |
| Stress test | 500 req/s | 5m | Find breaking point |

#### Custom Load Test

```bash
# Custom scenario
vegeta attack \
  -targets=tests/load/vegeta/auth_targets.txt \
  -rate=200 \
  -duration=60s \
  -output=test/results/load/custom.bin

# Generate report
vegeta report -type=html -input=test/results/load/custom.bin -output=report.html
```

---

## Performance SLAs

### Service Level Agreements

Based on benchmark testing, the following SLAs are established:

#### API Endpoints

| Endpoint | P50 Latency | P95 Latency | P99 Latency | Throughput |
|----------|-------------|-------------|-------------|------------|
| POST /auth/login | < 40ms | < 80ms | < 150ms | 500 req/s |
| POST /auth/register | < 60ms | < 120ms | < 250ms | 200 req/s |
| GET /auth/me | < 20ms | < 40ms | < 80ms | 1000 req/s |
| POST /orders | < 50ms | < 100ms | < 200ms | 300 req/s |
| GET /orders | < 30ms | < 60ms | < 100ms | 500 req/s |
| PUT /orders/{id}/cancel | < 50ms | < 100ms | < 200ms | 200 req/s |

#### Service Methods

| Method | P50 Latency | P95 Latency | P99 Latency |
|--------|-------------|-------------|-------------|
| OrderService.CreateOrder | < 50ms | < 100ms | < 200ms |
| OrderService.GetUserOrders | < 30ms | < 60ms | < 100ms |
| AuthService.Login | < 30ms | < 60ms | < 100ms |
| AuthService.Register | < 50ms | < 100ms | < 200ms |
| PaymentService.CreatePayment | < 20ms | < 40ms | < 80ms |

#### Database Operations

| Operation | P50 Latency | P95 Latency | P99 Latency |
|-----------|-------------|-------------|-------------|
| Single Insert | < 5ms | < 10ms | < 20ms |
| Batch Insert (100) | < 50ms | < 100ms | < 200ms |
| Primary Key Lookup | < 1ms | < 2ms | < 5ms |
| Indexed Lookup | < 2ms | < 5ms | < 10ms |
| Simple Join | < 10ms | < 20ms | < 40ms |

### Alert Thresholds

Set up monitoring alerts when:

- **P95 latency** exceeds 2x baseline
- **P99 latency** exceeds 3x baseline
- **Error rate** exceeds 1%
- **Throughput** drops below 80% of baseline

---

## Interpreting Results

### Benchmark Output Format

```
BenchmarkOrderCreation_Simple-16     100000      105234 ns/op    5123 B/op    45 allocs/op
```

| Field | Description |
|-------|-------------|
| `-16` | GOMAXPROCS value (parallelism) |
| `100000` | Number of iterations executed |
| `105234 ns/op` | Average time per operation (105 microseconds) |
| `5123 B/op` | Memory allocated per operation (5 KB) |
| `45 allocs/op` | Number of memory allocations per operation |

### Using benchstat for Comparison

```bash
# Save baseline results
make bench-order > baseline.txt

# Make code changes

# Run benchmarks again
make bench-order > current.txt

# Compare results
make bench-compare OLD=baseline.txt NEW=current.txt

# Or use benchstat directly
benchstat baseline.txt current.txt
```

#### Interpreting benchstat Output

```
name          old time/op  new time/op  delta
OrderCreate   105ms ± 2%   110ms ± 3%  +4.81%  (p=0.016 n=10)
OrderList      30ms ± 1%    25ms ± 1%  -16.67%  (p=0.000 n=10)
```

- **+4.81%**: Performance regression (slower)
- **-16.67%**: Performance improvement (faster)
- **p value**: Statistical significance (< 0.05 is significant)

### Profiling Results

#### CPU Profiling

```bash
# Run benchmarks with CPU profiling
make bench-all

# Analyze CPU profile
go tool pprof -http=:8080 test/results/bench/cpu.prof

# Open browser to http://localhost:8080
```

Look for:
- Hot paths (functions consuming most CPU)
- Unexpected function calls
- Inefficient algorithms

#### Memory Profiling

```bash
# Analyze memory profile
go tool pprof -http=:8080 test/results/bench/mem.prof

# Or view top allocations
go tool pprof -top test/results/bench/mem.prof
```

Look for:
- Excessive allocations
- Memory leaks ( allocations that don't get freed)
- Large object allocations

---

## Continuous Benchmarking

### CI/CD Integration

Add to your CI pipeline (`.github/workflows/ci.yml`):

```yaml
name: Benchmark

on:
  push:
    branches: [main, dev]
  pull_request:
    branches: [main]

jobs:
  benchmark:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.25'

      - name: Run benchmarks
        run: |
          make bench-tools
          make bench-all

      - name: Store benchmark result
        uses: benchmark-action/github-action-benchmark@v1
        with:
          tool: 'go'
          output-file-path: test/results/bench/results.txt
          github-token: ${{ secrets.GITHUB_TOKEN }}
          auto-push: true
```

### Baseline Management

1. **Establish baselines** for each critical operation
2. **Store baselines** in version control
3. **Compare every PR** against baselines
4. **Block regressions** that exceed thresholds

```bash
# In CI script
if [ -f "baseline.txt" ]; then
  benchstat baseline.txt current.txt > comparison.txt
  if grep -q "~.*%.*p=0.000" comparison.txt; then
    echo "Significant performance regression detected!"
    exit 1
  fi
fi
```

---

## Troubleshooting

### Common Issues

#### 1. Database Connection Errors

**Problem**: `Failed to connect to benchmark database`

**Solution**:
```bash
# Check database is running
docker ps | grep gamelink-bench-db

# Verify environment variables
echo $BENCH_DB_HOST
echo $BENCH_DB_NAME

# Test connection
psql -h localhost -U gamelink -d gamelimk_bench
```

#### 2. Unstable Benchmark Results

**Problem**: Results vary significantly between runs

**Solution**:
- Increase `benchtime` to get more samples: `-benchtime=30s`
- Run multiple times and average: `for i in {1..5}; do make bench-order >> results.txt; done`
- Ensure consistent system load (close other applications)
- Use dedicated benchmarking database (not shared)

#### 3. Memory Issues During Benchmarks

**Problem**: `out of memory` errors

**Solution**:
```bash
# Reduce concurrent operations
export GOMAXPROCS=2

# Run smaller batches
go test -bench=. -benchtime=5s

# Increase swap if needed
sudo swapon /swapfile
```

#### 4. Vegeta Load Test Failures

**Problem**: All requests fail with connection errors

**Solution**:
```bash
# Verify API server is running
curl http://localhost:8080/api/v1/healthz

# Check firewall settings
sudo ufw allow 8080

# Reduce rate to avoid overwhelming server
make load-test-auth RATE=10 DURATION=10s
```

---

## Best Practices

### Running Benchmarks

1. **Consistent Environment**
   - Use same hardware for comparison
   - Close unnecessary applications
   - Run multiple times for accuracy

2. **Isolate Benchmarks**
   - Use dedicated benchmark database
   - Avoid network calls to external services
   - Mock external dependencies

3. **Profile Before Optimizing**
   - Use CPU profiling to find bottlenecks
   - Use memory profiling to find leaks
   - Focus on high-impact optimizations

4. **Track Results Over Time**
   - Store historical benchmark results
   - Plot trends to detect gradual degradation
   - Correlate performance with code changes

### Optimization Checklist

- [ ] Identified bottleneck through profiling
- [ ] Established baseline metric
- [ ] Implemented optimization
- [ ] Re-ran benchmarks with same conditions
- [ ] Verified improvement is statistically significant
- [ ] Updated baseline documentation
- [ ] Added regression test

---

## Additional Resources

### Documentation

- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)
- [Vegeta Documentation](https://github.com/tsenart/vegeta)
- [pprof Profiling](https://go.dev/blog/pprof)

### Tools

- **benchstat**: Compare benchmark results
- **vegeta**: HTTP load testing
- **pprof**: CPU and memory profiling
- **go test -bench**: Built-in benchmarking

### Related Project Files

- `api/internal/benchmark/framework/framework.go` - Benchmark framework
- `api/Makefile` - Benchmark make targets
- `api/tests/load/` - Load testing scripts

---

## Performance Benchmarking Checklist

When running performance benchmarks for GameLink:

- [ ] Database is set up and accessible
- [ ] Benchmark tools are installed (`make bench-tools`)
- [ ] API server is running (for HTTP benchmarks)
- [ ] System load is minimal
- [ ] Baseline metrics are documented
- [ ] Benchmarks run for sufficient duration (≥10s)
- [ ] Results are saved for comparison
- [ ] Significant changes are investigated
- [ ] Documentation is updated
- [ ] CI/CD is configured for regression detection

---

For questions or issues related to performance benchmarking, please refer to the project documentation or open an issue on GitHub.
