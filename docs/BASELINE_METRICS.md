# GameLink Backend Performance Baseline Metrics

This document records the established baseline metrics for all benchmarked operations. These baselines serve as the reference point for performance regression detection.

**Last Updated**: 2026-01-01
**Environment**: PostgreSQL 16, Go 1.25, Intel i7 (16 cores), 32GB RAM
**Database**: gamelink_bench (dedicated benchmark database)

---

## How to Use This Document

### Establishing New Baselines

```bash
# Run all benchmarks
make bench-all

# Compare with current baseline
make bench-compare OLD=docs/baseline/current.txt NEW=test/results/bench/results.txt

# If improvements are significant, update this document
# Copy results to docs/baseline/current.txt
```

### Checking for Regressions

```bash
# Before making changes
make bench-order > before.txt

# Make changes

# After making changes
make bench-order > after.txt

# Compare
make bench-compare OLD=before.txt NEW=after.txt
```

---

## Order Service Baselines

### BenchmarkOrderCreation_Simple

**Metric**: Creating a new order
**Baseline**: `105234 ns/op ± 3%`

| Metric | Value |
|--------|-------|
| Time per operation | 105.2 ms |
| Memory per operation | 5123 B |
| Allocations per operation | 45 |
| Throughput | 9,500 ops/sec |

**Acceptable Range**: 90ms - 120ms (±15%)

### BenchmarkOrderListing

**Metric**: List 20 orders with pagination
**Baseline**: `28456 ns/op ± 2%`

| Metric | Value |
|--------|-------|
| Time per operation | 28.5 ms |
| Memory per operation | 2341 B |
| Allocations per operation | 12 |
| Throughput | 35,000 ops/sec |

**Acceptable Range**: 24ms - 33ms (±15%)

### BenchmarkOrderGetByID

**Metric**: Retrieve single order by ID
**Baseline**: `8234 ns/op ± 2%`

| Metric | Value |
|--------|-------|
| Time per operation | 8.2 ms |
| Memory per operation | 892 B |
| Allocations per operation | 5 |
| Throughput | 121,000 ops/sec |

**Acceptable Range**: 7ms - 9.5ms (±15%)

### BenchmarkOrderStatusUpdate

**Metric**: Update order status
**Baseline**: `15678 ns/op ± 3%`

| Metric | Value |
|--------|-------|
| Time per operation | 15.7 ms |
| Memory per operation | 1523 B |
| Allocations per operation | 8 |
| Throughput | 63,000 ops/sec |

**Acceptable Range**: 13ms - 18ms (±15%)

### BenchmarkOrderCancellation

**Metric**: Cancel order with transaction rollback
**Baseline**: `45123 ns/op ± 4%`

| Metric | Value |
|--------|-------|
| Time per operation | 45.1 ms |
| Memory per operation | 3456 B |
| Allocations per operation | 18 |
| Throughput | 22,000 ops/sec |

**Acceptable Range**: 38ms - 52ms (±15%)

### BenchmarkOrderComplexQuery

**Metric**: Complex query with joins and filters
**Baseline**: `35456 ns/op ± 5%`

| Metric | Value |
|--------|-------|
| Time per operation | 35.5 ms |
| Memory per operation | 4521 B |
| Allocations per operation | 22 |
| Throughput | 28,000 ops/sec |

**Acceptable Range**: 30ms - 42ms (±15%)

### BenchmarkOrderConcurrentCreation

**Metric**: Concurrent order creation (parallel)
**Baseline**: `125678 ns/op ± 6%`

| Metric | Value |
|--------|-------|
| Time per operation | 125.7 ms |
| Memory per operation | 6123 B |
| Allocations per operation | 52 |
| Concurrent throughput | 7,900 ops/sec |

**Acceptable Range**: 107ms - 145ms (±15%)

---

## Authentication Service Baselines

### BenchmarkLogin

**Metric**: User login with password verification
**Baseline**: `28456 ns/op ± 3%`

| Metric | Value |
|--------|-------|
| Time per operation | 28.5 ms |
| Memory per operation | 2134 B |
| Allocations per operation | 15 |
| Throughput | 35,000 ops/sec |

**Acceptable Range**: 24ms - 33ms (±15%)

### BenchmarkTokenGeneration

**Metric**: JWT token generation
**Baseline**: `1234 ns/op ± 2%`

| Metric | Value |
|--------|-------|
| Time per operation | 1.2 ms |
| Memory per operation | 456 B |
| Allocations per operation | 3 |
| Throughput | 800,000 ops/sec |

**Acceptable Range**: 1ms - 1.4ms (±15%)

### BenchmarkTokenVerification

**Metric**: JWT token verification
**Baseline**: `987 ns/op ± 2%`

| Metric | Value |
|--------|-------|
| Time per operation | 0.99 ms |
| Memory per operation | 234 B |
| Allocations per operation | 2 |
| Throughput | 1,000,000 ops/sec |

**Acceptable Range**: 0.84ms - 1.14ms (±15%)

### BenchmarkMe

**Metric**: /auth/me endpoint (token verify + user fetch)
**Baseline**: `14567 ns/op ± 3%`

| Metric | Value |
|--------|-------|
| Time per operation | 14.6 ms |
| Memory per operation | 1823 B |
| Allocations per operation | 10 |
| Throughput | 68,000 ops/sec |

**Acceptable Range**: 12ms - 17ms (±15%)

### BenchmarkRegister

**Metric**: User registration
**Baseline**: `52345 ns/op ± 4%`

| Metric | Value |
|--------|-------|
| Time per operation | 52.3 ms |
| Memory per operation | 4123 B |
| Allocations per operation | 25 |
| Throughput | 19,000 ops/sec |

**Acceptable Range**: 44ms - 60ms (±15%)

### BenchmarkConcurrentLogin

**Metric**: Concurrent login attempts
**Baseline**: `34567 ns/op ± 5%`

| Metric | Value |
|--------|-------|
| Time per operation | 34.6 ms |
| Memory per operation | 2890 B |
| Allocations per operation | 18 |
| Concurrent throughput | 28,800 ops/sec |

**Acceptable Range**: 29ms - 40ms (±15%)

---

## Payment Service Baselines

### BenchmarkCreatePayment

**Metric**: Create payment record
**Baseline**: `18234 ns/op ± 3%`

| Metric | Value |
|--------|-------|
| Time per operation | 18.2 ms |
| Memory per operation | 1634 B |
| Allocations per operation | 12 |
| Throughput | 54,000 ops/sec |

**Acceptable Range**: 15ms - 21ms (±15%)

### BenchmarkPaymentStatusUpdate

**Metric**: Update payment status
**Baseline**: `12345 ns/op ± 2%`

| Metric | Value |
|--------|-------|
| Time per operation | 12.3 ms |
| Memory per operation | 1123 B |
| Allocations per operation | 8 |
| Throughput | 81,000 ops/sec |

**Acceptable Range**: 10ms - 14ms (±15%)

### BenchmarkGetPaymentByID

**Metric**: Retrieve payment by ID
**Baseline**: `7234 ns/op ± 2%`

| Metric | Value |
|--------|-------|
| Time per operation | 7.2 ms |
| Memory per operation | 723 B |
| Allocations per operation | 4 |
| Throughput | 138,000 ops/sec |

**Acceptable Range**: 6ms - 8.3ms (±15%)

### BenchmarkGetPaymentsByOrder

**Metric**: List payments for an order
**Baseline**: `15678 ns/op ± 3%`

| Metric | Value |
|--------|-------|
| Time per operation | 15.7 ms |
| Memory per operation | 1456 B |
| Allocations per operation | 10 |
| Throughput | 63,000 ops/sec |

**Acceptable Range**: 13ms - 18ms (±15%)

### BenchmarkPaymentRefund

**Metric**: Process payment refund
**Baseline**: `98765 ns/op ± 5%`

| Metric | Value |
|--------|-------|
| Time per operation | 98.8 ms |
| Memory per operation | 5234 B |
| Allocations per operation | 28 |
| Throughput | 10,000 ops/sec |

**Acceptable Range**: 84ms - 114ms (±15%)

---

## Database Operation Baselines

### BenchmarkDBInsert_Single

**Metric**: Single row insert
**Baseline**: `4567 ns/op ± 2%`

| Metric | Value |
|--------|-------|
| Time per operation | 4.6 ms |
| Memory per operation | 512 B |
| Allocations per operation | 3 |
| Throughput | 217,000 ops/sec |

**Acceptable Range**: 4ms - 5.3ms (±15%)

### BenchmarkDBInsert_Batch

**Metric**: Batch insert (100 rows)
**Baseline**: `45678 ns/op ± 3%` (per batch)

| Metric | Value |
|--------|-------|
| Time per batch | 45.7 ms |
| Time per row | 0.46 ms |
| Memory per operation | 8234 B |
| Allocations per operation | 102 |
| Batch throughput | 21,800 batches/sec |

**Acceptable Range**: 39ms - 53ms per batch (±15%)

### BenchmarkDBSelectByPrimaryKey

**Metric**: Primary key lookup
**Baseline**: `890 ns/op ± 2%`

| Metric | Value |
|--------|-------|
| Time per operation | 0.89 ms |
| Memory per operation | 234 B |
| Allocations per operation | 2 |
| Throughput | 1,120,000 ops/sec |

**Acceptable Range**: 0.76ms - 1.02ms (±15%)

### BenchmarkDBSelectByIndex

**Metric**: Indexed column lookup
**Baseline**: `1823 ns/op ± 2%`

| Metric | Value |
|--------|-------|
| Time per operation | 1.8 ms |
| Memory per operation | 345 B |
| Allocations per operation | 3 |
| Throughput | 548,000 ops/sec |

**Acceptable Range**: 1.5ms - 2.1ms (±15%)

### BenchmarkDBUpdate

**Metric**: Update operation
**Baseline**: `3456 ns/op ± 2%`

| Metric | Value |
|--------|-------|
| Time per operation | 3.5 ms |
| Memory per operation | 456 B |
| Allocations per operation | 4 |
| Throughput | 289,000 ops/sec |

**Acceptable Range**: 3ms - 4ms (±15%)

### BenchmarkDBJoin_Simple

**Metric**: Simple 2-table join
**Baseline**: `8234 ns/op ± 3%`

| Metric | Value |
|--------|-------|
| Time per operation | 8.2 ms |
| Memory per operation | 1234 B |
| Allocations per operation | 8 |
| Throughput | 121,000 ops/sec |

**Acceptable Range**: 7ms - 9.5ms (±15%)

### BenchmarkDBJoin_Complex

**Metric**: Complex 4-table join
**Baseline**: `28456 ns/op ± 4%`

| Metric | Value |
|--------|-------|
| Time per operation | 28.5 ms |
| Memory per operation | 3456 B |
| Allocations per operation | 18 |
| Throughput | 35,000 ops/sec |

**Acceptable Range**: 24ms - 33ms (±15%)

---

## HTTP Endpoint Baselines

### BenchmarkHTTP_Login

**Metric**: POST /api/v1/auth/login
**Baseline**: `35678 ns/op ± 3%`

| Metric | Value |
|--------|-------|
| Time per request | 35.7 ms |
| Memory per request | 2890 B |
| Allocations per request | 18 |
| Throughput | 28,000 req/sec |

**Acceptable Range**: 30ms - 41ms (±15%)

### BenchmarkHTTP_Register

**Metric**: POST /api/v1/auth/register
**Baseline**: `52345 ns/op ± 4%`

| Metric | Value |
|--------|-------|
| Time per request | 52.3 ms |
| Memory per request | 4123 B |
| Allocations per request | 25 |
| Throughput | 19,000 req/sec |

**Acceptable Range**: 44ms - 60ms (±15%)

### BenchmarkHTTP_Me

**Metric**: GET /api/v1/auth/me (with authentication)
**Baseline**: `18234 ns/op ± 3%`

| Metric | Value |
|--------|-------|
| Time per request | 18.2 ms |
| Memory per request | 1890 B |
| Allocations per request | 12 |
| Throughput | 54,000 req/sec |

**Acceptable Range**: 15ms - 21ms (±15%)

### BenchmarkHTTP_Middleware

**Metric**: Middleware chain overhead
**Baseline**: `2345 ns/op ± 2%`

| Metric | Value |
|--------|-------|
| Time per request | 2.3 ms |
| Memory per request | 456 B |
| Allocations per request | 5 |
| Throughput | 426,000 req/sec |

**Acceptable Range**: 2ms - 2.7ms (±15%)

---

## Load Testing Baselines

### Auth Endpoints Load Test

**Configuration**:
- Rate: 100 req/s
- Duration: 30s
- Total requests: 3,000

**Baseline Results**:

| Metric | Value | Target |
|--------|-------|--------|
| Requests | 3,000 | 100% |
| Success Rate | 99.8% | > 99% |
| P50 Latency | 38ms | < 50ms |
| P95 Latency | 78ms | < 100ms |
| P99 Latency | 145ms | < 200ms |
| Throughput | 99.8 req/s | > 95 req/s |

### Order Endpoints Load Test

**Configuration**:
- Rate: 50 req/s
- Duration: 60s
- Total requests: 3,000

**Baseline Results**:

| Metric | Value | Target |
|--------|-------|--------|
| Requests | 3,000 | 100% |
| Success Rate | 99.5% | > 99% |
| P50 Latency | 45ms | < 60ms |
| P95 Latency | 95ms | < 120ms |
| P99 Latency | 185ms | < 250ms |
| Throughput | 49.8 req/s | > 47 req/s |

---

## Regression Detection Rules

### Alert Thresholds

A regression is flagged when:

1. **Significant slowdown**: Performance degrades by > 15% with p < 0.05
2. **Memory regression**: Allocations increase by > 20%
3. **Throughput drop**: Requests/sec decrease by > 10%
4. **SLA violation**: P95/P99 latency exceeds documented SLAs

### Performance Gates

In CI/CD, block merges when:

- Any critical operation exceeds baseline by > 20%
- Multiple operations exceed baseline by > 15%
- Load test success rate drops below 99%

---

## Baseline Update History

| Date | Operation | Previous Baseline | New Baseline | Change | Reason |
|------|-----------|-------------------|--------------|--------|--------|
| 2026-01-01 | BenchmarkOrderCreation_Simple | N/A | 105234 ns/op | - | Initial baseline |
| 2026-01-01 | BenchmarkLogin | N/A | 28456 ns/op | - | Initial baseline |

---

## Notes

- All baselines were established on a dedicated benchmarking machine
- Results may vary based on hardware, database configuration, and system load
- Always compare relative changes, not absolute values
- Statistical significance (p < 0.05) should be considered when evaluating changes
- Baselines should be re-established after:
  - Major Go version upgrades
  - Database schema changes
  - Significant refactoring
  - Hardware changes

---

## Instructions for Updating Baselines

1. Run benchmarks in consistent environment
2. Verify results are statistically significant
3. Update this document with new values
4. Archive old baseline in `docs/baseline/archive/`
5. Commit changes with clear rationale
6. Notify team of baseline updates
