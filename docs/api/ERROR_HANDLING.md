# Error Response Handling Guide

## Overview

GameLink uses a unified error response system centered around the `resp` package. All handlers MUST use this system for consistent API responses.

## Core Package: `internal/handler/resp`

### Success Response Functions

| Function | Description | HTTP Status |
|----------|-------------|-------------|
| `OK[T](c, data)` | Success with "OK" message | 200 |
| `Success[T](c, message, data)` | Success with custom message | 200 |
| `Created[T](c, data)` | Resource created | 201 |
| `Updated[T](c, data)` | Resource updated | 200 |
| `Deleted(c)` | Resource deleted | 200 |
| `List[T](c, data, pagination)` | Paginated list | 200 |

### Error Response Functions

| Function | Description | Usage |
|----------|-------------|------|
| `Error(c, err)` | Auto-detects error type | **Primary error handler** |
| `ErrorMsg(c, status, msg)` | Custom error response | Direct control needed |
| `BadRequest(c, msg)` | 400 error | Invalid input |
| `Unauthorized(c, msg)` | 401 error | Not authenticated |
| `Forbidden(c, msg)` | 403 error | No permission |
| `NotFound(c, msg)` | 404 error | Resource not found |
| `InternalError(c, msg)` | 500 error | Server error |

## Error Types Supported

The `resp.Error()` function automatically handles:

1. **`*apierr.APIError`** - Structured API errors with details
   ```go
   resp.Error(c, apierr.BadRequest("invalid input"))
   resp.Error(c, apierr.NotFound("user not found"))
   ```

2. **`repository.ErrNotFound`** - Database not found errors
   ```go
   resp.Error(c, repository.ErrNotFound)
   ```

3. **Validation errors** - Auto-detected validation failures
   ```go
   resp.Error(c, err) // Detects validation errors
   ```

4. **Generic errors** - Fallback to 500
   ```go
   resp.Error(c, errors.New("something failed"))
   ```

## Handler-Level Wrappers

### Admin Handlers (`internal/handler/admin/helpers.go`)

```go
// Success wrappers
respondSuccess(c, data)
respondSuccessWithMsg(c, message, data)
respondCreated(c, data)
respondUpdated(c, data)
respondDeleted(c)
respondList(c, data, pagination)

// Error wrappers
respondError(c, err) // Handles adminservice.ErrNotFound
respondAPIError(c, err) // Alias for respondError

// Parameter parsing with auto-response
ParseIDAndRespond(c, "id")
ValidateAndRespond(c, &request)
parsePagination(c)
```

### User/Player Handlers

```go
// Use local wrappers or resp package directly
respondAPIError(c, err)
```

## Response Format

### Success Response

```json
{
  "success": true,
  "code": 200,
  "message": "OK",
  "data": { ... },
  "traceId": "uuid-here"
}
```

### List Response (with pagination)

```json
{
  "success": true,
  "code": 200,
  "message": "OK",
  "data": [ ... ],
  "pagination": {
    "page": 1,
    "pageSize": 20,
    "total": 100,
    "totalPages": 5,
    "hasNext": true,
    "hasPrev": false
  },
  "traceId": "uuid-here"
}
```

### Error Response

```json
{
  "success": false,
  "code": 400,
  "message": "invalid input",
  "traceId": "uuid-here",
  "meta": {
    "details": "field 'email' is required",
    "field": "email",
    "timestamp": 1234567890
  }
}
```

## Best Practices

### 1. Always Use `resp.Error()` for Errors

❌ **Bad:**
```go
c.JSON(500, gin.H{"error": "failed"})
```

✅ **Good:**
```go
resp.Error(c, apierr.InternalError("failed"))
```

### 2. Use `apierr` Package for Structured Errors

❌ **Bad:**
```go
resp.ErrorMsg(c, 400, "invalid input")
```

✅ **Good:**
```go
resp.Error(c, apierr.BadRequest("invalid input"))
```

### 3. Chain Error Details

```go
resp.Error(c, apierr.BadRequest("validation failed")
    .WithDetails("email is required")
    .WithField("email"))
```

### 4. Handle Service Errors Correctly

```go
result, err := service.DoSomething(ctx, input)
if err != nil {
    respondError(c, err)  // Handles all error types
    return
}
respondSuccess(c, result)
```

### 5. Return Early on Errors

```go
id, ok := ParseIDAndRespond(c, "id")
if !ok {
    return  // Error already sent
}
// Continue with valid id...
```

## Error Message Constants

Use predefined error constants from `apierr`:

```go
const (
    ErrInvalidJSONPayload   = "invalid JSON payload"
    ErrInvalidID            = "invalid ID"
    ErrInvalidPage          = "invalid page"
    ErrInvalidPageSize      = "invalid page size"
    ErrUserNotFound         = "user not found"
    ErrOrderNotFound        = "order not found"
    ErrPaymentNotFound      = "payment not found"
    ErrPlayerNotFound       = "player not found"
    ErrGameNotFound         = "game not found"
)
```

## Trace ID Support

All responses automatically include `traceId` when set in context:

```go
// In middleware
c.Set("request_id", uuid.New().String())

// In handler
resp.OK(c, data) // Automatically includes traceId
```

## Testing Error Responses

```go
func TestHandler_Error(t *testing.T) {
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)

    // Trigger error
    Handler(c)

    // Assert response
    assert.Equal(t, 404, w.Code)
    var resp model.SuccessResponse
    json.Unmarshal(w.Body.Bytes(), &resp)
    assert.False(t, resp.Success)
    assert.Equal(t, "resource not found", resp.Message)
}
```

## Migration Guide

If you find legacy code using direct `c.JSON` or `c.String` for errors:

### Before
```go
if err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
```

### After
```go
if err != nil {
    resp.Error(c, apierr.BadRequest(err.Error()))
    return
}
```
