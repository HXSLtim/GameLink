# Timeout Control Patterns in GameLink

This document describes the timeout control patterns used throughout the GameLink backend services to ensure reliable external service calls and prevent resource exhaustion.

## Overview

All external service calls MUST have appropriate timeout controls to prevent:
- Hanging requests that consume goroutines indefinitely
- Cascading failures when external services are slow
- Resource exhaustion under high load
- Poor user experience from unresponsive operations

## Timeout Duration Guidelines

| Operation Type | Timeout | Rationale |
|----------------|---------|-----------|
| **Database operations** | 5-10s | Simple queries should complete quickly; complex aggregations may need up to 10s |
| **Payment gateway calls** | 60s | Banking systems can be slow; refunds especially require extra time |
| **HTTP client calls** | 10-30s | Individual API requests should be fast; payment APIs get 30s |
| **WebSocket operations** | 10s | Real-time communication should be fast |
| **File uploads (OSS)** | 60s | Large files need time to upload |
| **File deletion (OSS)** | 10s | Delete operations should be fast |
| **SMS/API gateway calls** | 10s | External API calls should respond quickly |
| **Analytics queries** | 10s | Multiple database operations but should be optimized |

## Implementation Patterns

### Pattern 1: Context Timeout Wrapper (Recommended)

Use `context.WithTimeout()` for service-level operations:

```go
func (s *MyService) DoSomething(ctx context.Context, req Request) (*Response, error) {
    // Timeout choice: 10s for this operation
    // Reason: [explain why this timeout is appropriate]
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    // All downstream calls use the ctx with timeout
    result, err := s.repo.Get(ctx, id)
    if err != nil {
        return nil, err
    }

    return result, nil
}
```

**Benefits:**
- Context is propagated to all downstream calls (database, HTTP, etc.)
- Single timeout controls the entire operation
- Deferred cancel ensures cleanup

**When to use:**
- Service methods that coordinate multiple operations
- Methods that call repositories or external services
- Any public API method

### Pattern 2: HTTP Client Timeout

For individual HTTP requests, set `http.Client` timeout:

```go
// Timeout choice: 10s for external API calls
// Reason: Individual HTTP requests should complete quickly
client := &http.Client{Timeout: 10 * time.Second}
resp, err := client.Do(req)
```

**Benefits:**
- Prevents hanging HTTP connections
- Works even when request lacks context
- Explicit per-request timeout

**When to use:**
- Direct HTTP calls to external APIs
- When you don't control the caller's context
- In conjunction with Pattern 1 for defense in depth

### Pattern 3: Context-Aware HTTP Requests

Use `http.NewRequestWithContext()` to link HTTP request timeout to context:

```go
func (s *MyService) CallExternalAPI(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return err
    }

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    // ...
}
```

**Benefits:**
- HTTP request respects context cancellation
- Timeout is enforced at both context and HTTP client level
- Caller can cancel the operation

**When to use:**
- All HTTP requests in production code
- When you have a context available
- For better integration with Go's context propagation

## Files Modified

### External API Calls

1. **api/internal/service/sms/tencent.go**
   - Added: 10s timeout for SMS API calls
   - Line: 126-128
   - Pattern: HTTP client timeout + context propagation

2. **api/internal/service/payment/alipayProvider.go**
   - Added: 60s timeout wrapper for `Refund()` and `CreateOrder()`
   - Added: 30s HTTP client timeout for `doRequest()`
   - Lines: 63-66, 115-118, 218-220
   - Pattern: Context timeout wrapper + HTTP client timeout

3. **api/internal/service/payment/wechatProvider.go**
   - Added: 60s timeout wrapper for `Refund()` and `CreateOrder()`
   - Added: 30s HTTP client timeout for `doRefundRequest()`
   - Lines: 40-43, 198-201, 169-171
   - Pattern: Context timeout wrapper + HTTP client timeout

4. **api/internal/service/oss/service.go**
   - Added: 60s timeout for file uploads (can be large)
   - Added: 10s timeout for file deletions (should be fast)
   - Lines: 140-142, 172-173
   - Pattern: HTTP client timeout with documentation

### Database Operations

5. **api/internal/service/analytics/analytics.go**
   - Added: 10s timeout for all public analytics methods
   - Methods: `GetActiveUsers()`, `GetRetention()`, `GetPaymentAnalytics()`, `GetConversionFunnel()`
   - Lines: 90-93, 176-179, 269-272, 371-373
   - Pattern: Context timeout wrapper
   - Note: Service already uses `db.WithContext(ctx)` for DB operations

### Already Compliant

6. **api/internal/service/audit/audit.go**
   - Already has: 10s timeout for batch database writes
   - Line: 628
   - Pattern: Context timeout in background goroutine

7. **api/internal/service/chat/service.go**
   - All methods use context propagation
   - Repository calls use context
   - No external HTTP calls (WebSocket handled separately)

## Testing Timeout Behavior

### Unit Test Example

```go
func TestServiceTimeout(t *testing.T) {
    // Create a context that's already canceled
    ctx, cancel := context.WithCancel(context.Background())
    cancel() // Cancel immediately

    // Service should respect the canceled context
    _, err := service.DoSomething(ctx, req)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "context canceled")
}
```

### Integration Test Example

```go
func TestServiceExternalCallTimeout(t *testing.T) {
    // Use a slow HTTP server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(5 * time.Second) // Simulate slow response
        w.WriteHeader(http.StatusOK)
    }))
    defer server.Close()

    // Service should timeout before server responds
    ctx := context.Background()
    _, err := service.CallExternalAPI(ctx, server.URL)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "deadline exceeded")
}
```

## Common Mistakes to Avoid

### ❌ Don't: Use context.Background() in service methods

```go
func (s *MyService) DoSomething(ctx context.Context) error {
    // WRONG: Ignores caller's context
    innerCtx := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    return s.repo.Create(innerCtx, data)
}
```

### ✅ Do: Derive from caller's context

```go
func (s *MyService) DoSomething(ctx context.Context) error {
    // CORRECT: Respects caller's context and deadlines
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    return s.repo.Create(ctx, data)
}
```

### ❌ Don't: Forget deferred cancel

```go
func (s *MyService) DoSomething(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    // WRONG: Missing defer cancel() - context leaks!

    return s.repo.Create(ctx, data)
}
```

### ✅ Do: Always defer cancel()

```go
func (s *MyService) DoSomething(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel() // CORRECT: Ensures context cleanup

    return s.repo.Create(ctx, data)
}
```

### ❌ Don't: Set excessively long timeouts

```go
// WRONG: 5 minutes is too long for a simple API call
ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
```

### ✅ Do: Set appropriate timeouts

```go
// CORRECT: 10s is reasonable for most API calls
// Timeout choice: 10s for external API calls
// Reason: Individual HTTP requests should complete quickly
ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
```

## Monitoring and Alerting

### Metrics to Track

1. **Timeout Rate**: Percentage of requests that timeout vs. total requests
2. **Timeout Duration**: P50, P95, P99 latencies
3. **External Service Health**: Error rates from external APIs

### Example Prometheus Metrics

```go
var (
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "service_request_duration_milliseconds",
            Help:    "Request duration in milliseconds",
            Buckets: prometheus.ExponentialBuckets(10, 2, 10),
        },
        []string{"service", "method", "status"},
    )

    requestTimeouts = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "service_request_timeouts_total",
            Help: "Total number of request timeouts",
        },
        []string{"service", "method"},
    )
)
```

### Alerting Rules

```yaml
# Alert when timeout rate exceeds 5% for 5 minutes
- alert: HighTimeoutRate
  expr: rate(service_request_timeouts_total[5m]) / rate(service_request_duration_milliseconds_count[5m]) > 0.05
  for: 5m
  annotations:
    summary: "High timeout rate detected for {{ $labels.service }}/{{ $labels.method }}"
```

## Checklist for New Code

When adding new service methods or external calls:

- [ ] Does this method call external services (HTTP, database, message queue)?
- [ ] Have I added `context.WithTimeout()` at the service method level?
- [ ] Is the timeout duration appropriate for the operation type?
- [ ] Have I added a comment explaining the timeout choice?
- [ ] Is `defer cancel()` called immediately after `context.WithTimeout()`?
- [ ] Do HTTP clients have explicit timeouts set?
- [ ] Are contexts propagated to all downstream calls (repositories, HTTP requests)?
- [ ] Have I written tests for timeout behavior?
- [ ] Is the timeout documented in this file or inline comments?

## References

- [Go Context Documentation](https://pkg.go.dev/context)
- [Go net/http Client Timeout](https://pkg.go.dev/net/http#Client)
- [Google Cloud Timeouts Best Practices](https://cloud.google.com/architecture/timeout-best-practices)
- [GameLink Testing Standards](../.kiro/steering/05-testing-standard.md)

## Changelog

| Date | Change | Files |
|------|--------|-------|
| 2026-01-03 | Initial timeout controls added | oss/service.go, sms/tencent.go, payment/*, analytics/analytics.go |
| 2026-01-03 | Documentation created | docs/TIMEOUT_PATTERNS.md |

---

**Maintainer**: Backend Team
**Last Updated**: 2026-01-03
**Status**: Active
