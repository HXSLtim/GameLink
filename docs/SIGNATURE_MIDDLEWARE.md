# HMAC-SHA256 Request Signature Middleware

## Overview

The signature middleware provides HMAC-SHA256 request signature validation to ensure request integrity and authenticity. This middleware validates that incoming HTTP requests have been signed with a shared secret key, preventing tampering and unauthorized access.

## Architecture

### Signature Algorithm

The signature is calculated using HMAC-SHA256 with the following format:

```
SIGNATURE = HMAC-SHA256(METHOD:PATH:BODY, secret_key)
```

**Components:**
- **METHOD**: HTTP method in uppercase (e.g., "POST", "PUT", "DELETE")
- **PATH**: Request URL path (e.g., "/api/v1/users")
- **BODY**: Raw request body as string
- **secret_key**: Shared secret key configured on server and client

**Example:**
```
METHOD = "POST"
PATH = "/api/v1/users"
BODY = '{"name":"test","value":123}'
MESSAGE = "POST:/api/v1/users:{"name":"test","value":123}"
SIGNATURE = HMAC-SHA256(MESSAGE, "your-secret-key")
```

### Security Features

1. **Request Integrity**: Ensures request body hasn't been tampered with
2. **Authentication**: Verifies request comes from a trusted source with the secret key
3. **Timing Attack Prevention**: Uses constant-time comparison (`hmac.Equal`)
4. **Configurable Methods**: Only validates configured HTTP methods
5. **Path Exclusions**: Allows excluding specific paths from signature validation

## Configuration

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `SIGNATURE_ENABLED` | Enable/disable signature validation | No | `false` |
| `SIGNATURE_SECRET_KEY` | HMAC secret key (32+ bytes recommended) | Yes (if enabled) | - |
| `SIGNATURE_HEADER_NAME` | Header name for signature | No | `X-Signature` |
| `SIGNATURE_METHODS` | HTTP methods to validate (comma-separated) | No | `POST,PUT,PATCH,DELETE` |
| `SIGNATURE_EXCLUDE_PATHS` | Paths to exclude (comma-separated) | No | `/api/v1/health,/api/v1/ping` |

### Configuration Files

#### Development (`api/configs/config.development.yaml`)

```yaml
signature:
  enabled: false
  secret_key: ""
  header_name: "X-Signature"
  methods:
    - POST
    - PUT
    - PATCH
    - DELETE
  exclude_paths:
    - "/api/v1/health"
    - "/api/v1/ping"
    - "/api/v1/auth/refresh"
```

#### Production (`api/configs/config.production.yaml`)

```yaml
signature:
  enabled: true
  secret_key: ""  # Must be set via environment variable
  header_name: "X-Signature"
  methods:
    - POST
    - PUT
    - PATCH
    - DELETE
  exclude_paths:
    - "/api/v1/health"
    - "/api/v1/ping"
    - "/api/v1/auth/refresh"
```

### Secret Key Generation

Generate a secure 32-byte secret key:

```bash
# Using OpenSSL
openssl rand -base64 32

# Using Python
python -c "import secrets; print(secrets.token_urlsafe(32))"

# Using Node.js
node -e "console.log(require('crypto').randomBytes(32).toString('base64'))"
```

## Usage

### Client-Side Implementation

#### JavaScript/TypeScript (Node.js)

```typescript
import crypto from 'crypto';

function calculateSignature(method: string, path: string, body: string, secretKey: string): string {
  const message = `${method.toUpperCase()}:${path}:${body}`;
  const hmac = crypto.createHmac('sha256', secretKey);
  hmac.update(message);
  return hmac.digest('hex');
}

// Example usage
const method = 'POST';
const path = '/api/v1/users';
const body = JSON.stringify({ name: 'test', value: 123 });
const secretKey = process.env.SIGNATURE_SECRET_KEY;

const signature = calculateSignature(method, path, body, secretKey);

// Make request with signature
fetch('http://localhost:8080/api/v1/users', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'X-Signature': signature,
  },
  body: body,
});
```

#### JavaScript/TypeScript (Browser)

```typescript
import { hmacSha256 } from '@/utils/crypto';

async function calculateSignature(method: string, path: string, body: string, secretKey: string): Promise<string> {
  const message = `${method.toUpperCase()}:${path}:${body}`;
  return await hmacSha256(message, secretKey);
}

// Example usage
const method = 'POST';
const path = '/api/v1/users';
const body = JSON.stringify({ name: 'test', value: 123 });
const secretKey = import.meta.env.VITE_SIGNATURE_SECRET_KEY; // Frontend env var

const signature = await calculateSignature(method, path, body, secretKey);

// Make request with signature
fetch('/api/v1/users', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'X-Signature': signature,
  },
  body: body,
});
```

#### cURL

```bash
#!/bin/bash

METHOD="POST"
PATH="/api/v1/users"
BODY='{"name":"test","value":123}'
SECRET_KEY="your-secret-key"

# Calculate signature
MESSAGE="${METHOD}:${PATH}:${BODY}"
SIGNATURE=$(echo -n "$MESSAGE" | openssl dgst -sha256 -hmac "$SECRET_KEY" | awk '{print $2}')

# Make request
curl -X POST "http://localhost:8080${PATH}" \
  -H "Content-Type: application/json" \
  -H "X-Signature: $SIGNATURE" \
  -d "$BODY"
```

#### Python

```python
import hmac
import hashlib
import requests
import json

def calculate_signature(method: str, path: str, body: str, secret_key: str) -> str:
    message = f"{method.upper()}:{path}:{body}"
    return hmac.new(
        secret_key.encode(),
        message.encode(),
        hashlib.sha256
    ).hexdigest()

# Example usage
method = 'POST'
path = '/api/v1/users'
body = json.dumps({'name': 'test', 'value': 123})
secret_key = 'your-secret-key'

signature = calculate_signature(method, path, body, secret_key)

# Make request
response = requests.post(
    f'http://localhost:8080{path}',
    headers={
        'Content-Type': 'application/json',
        'X-Signature': signature,
    },
    data=body,
)
```

### Go (Backend)

```go
import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "strings"
)

func calculateSignature(method, path string, body []byte, secretKey string) string {
    message := strings.ToUpper(method) + ":" + path + ":" + string(body)
    h := hmac.New(sha256.New, []byte(secretKey))
    h.Write([]byte(message))
    return hex.EncodeToString(h.Sum(nil))
}

// Example usage
signature := calculateSignature("POST", "/api/v1/users", bodyBytes, secretKey)
req.Header.Set("X-Signature", signature)
```

## Middleware Behavior

### Request Validation Flow

1. **Configuration Check**: If `enabled: false`, skip validation
2. **Method Check**: Skip if request method not in `methods` list
3. **Path Exclusion**: Skip if path matches any exclude pattern
4. **Signature Extraction**: Extract signature from configured header
5. **Body Reading**: Read request body (non-destructive)
6. **Signature Calculation**: Calculate expected signature
7. **Validation**: Compare using constant-time comparison
8. **Result**: Store validation result in context

### Exclusion Patterns

The middleware supports three exclusion patterns:

1. **Exact Match**: `/api/v1/health` matches only `/api/v1/health`
2. **Wildcard Prefix**: `/api/v1/public/*` matches `/api/v1/public/` and all subpaths
3. **Catch-All**: `*` matches any path

### Error Responses

| Error | Status Code | Message |
|-------|-------------|---------|
| Missing signature header | 401 | `missing signature header` |
| Invalid signature | 401 | `invalid signature` |
| Unable to read body | 400 | `unable to read request body` |

## Testing

### Unit Tests

Run the signature middleware tests:

```bash
cd api
go test ./internal/handler/middleware/... -run TestSignature -v
```

### Test Coverage

The test suite includes:

- Disabled middleware bypass
- Missing secret key handling
- Valid signature acceptance
- Invalid signature rejection
- Missing signature header handling
- Path exclusion patterns
- Method filtering
- Custom header names
- Empty body handling
- Signature determinism
- Context value storage

### Integration Testing

Test signature validation with a real server:

```bash
# Start server with signature enabled
SIGNATURE_ENABLED=true \
SIGNATURE_SECRET_KEY="test-secret-key" \
go run cmd/main.go

# Test valid signature
BODY='{"test":"data"}'
MESSAGE="POST:/api/v1/users:${BODY}"
SIGNATURE=$(echo -n "$MESSAGE" | openssl dgst -sha256 -hmac "test-secret-key" | awk '{print $2}')

curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "X-Signature: $SIGNATURE" \
  -d "$BODY"

# Test invalid signature (should fail)
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "X-Signature: invalid" \
  -d "$BODY"
```

## Deployment

### Production Checklist

- [ ] Set `SIGNATURE_ENABLED=true` in production environment
- [ ] Generate and set `SIGNATURE_SECRET_KEY` via environment variable
- [ ] Ensure all clients have access to the secret key
- [ ] Configure appropriate exclude paths for public endpoints
- [ ] Update all API clients to include signature header
- [ ] Monitor signature validation failures in logs
- [ ] Implement key rotation strategy

### Environment Variables (Production)

```bash
export SIGNATURE_ENABLED=true
export SIGNATURE_SECRET_KEY="<your-32-byte-secret-key>"
export SIGNATURE_METHODS="POST,PUT,PATCH,DELETE"
export SIGNATURE_EXCLUDE_PATHS="/api/v1/health,/api/v1/ping"
```

### Docker Compose

```yaml
services:
  api:
    environment:
      - SIGNATURE_ENABLED=true
      - SIGNATURE_SECRET_KEY=${SIGNATURE_SECRET_KEY}
      - SIGNATURE_METHODS=POST,PUT,PATCH,DELETE
      - SIGNATURE_EXCLUDE_PATHS=/api/v1/health,/api/v1/ping
```

## Security Considerations

### Best Practices

1. **Secret Key Management**:
   - Use 32+ byte secret keys for HMAC-SHA256
   - Store secret keys in secure vaults (e.g., HashiCorp Vault, AWS Secrets Manager)
   - Rotate secret keys periodically (recommended: every 90 days)
   - Never commit secret keys to version control

2. **Transport Security**:
   - Always use HTTPS in production
   - Signature validation is additional security, not a replacement for TLS

3. **Error Handling**:
   - Don't reveal detailed signature calculation errors
   - Log signature validation failures for security monitoring
   - Implement rate limiting for failed signature attempts

4. **Key Distribution**:
   - Only share secret keys with trusted clients
   - Use different keys for different environments (dev/staging/prod)
   - Implement secure key distribution mechanisms

### Common Pitfalls

1. **Inconsistent Body Serialization**:
   - Ensure client and server use identical body serialization
   - Be careful with whitespace, encoding, and field ordering

2. **Time-Sensitive Data**:
   - Signatures don't include timestamps by default
   - Consider adding timestamp fields for replay attack prevention

3. **Query Parameters**:
   - Query parameters are NOT included in signature calculation
   - Add query params to body if they need signature coverage

## Troubleshooting

### Signature Validation Failures

**Symptom**: `401 Unauthorized` with "invalid signature"

**Possible Causes**:
1. Body serialization mismatch
2. Incorrect HTTP method casing
3. Path mismatch (trailing slashes)
4. Secret key mismatch
5. Whitespace in body

**Debug Steps**:
```bash
# Client-side calculation
METHOD="POST"
PATH="/api/v1/users"
BODY='{"name":"test"}'
MESSAGE="${METHOD}:${PATH}:${BODY}"
echo "Message: $MESSAGE"
echo "Signature: $(echo -n "$MESSAGE" | openssl dgst -sha256 -hmac "your-key")"

# Server-side calculation (check logs)
# Enable debug logging to see expected vs received signature
```

### Missing Signature Header

**Symptom**: `401 Unauthorized` with "missing signature header"

**Solution**:
- Ensure client sends `X-Signature` header (or custom header name)
- Check for typos in header name
- Verify middleware configuration for custom header names

### Performance Impact

The signature middleware adds minimal overhead:
- HMAC-SHA256 calculation: ~1µs per request
- Body reading: O(n) where n is body size
- Constant-time comparison: O(k) where k is signature length (64 hex chars)

For high-throughput scenarios, consider:
- Caching signatures for idempotent requests
- Using hardware acceleration for crypto operations
- Skipping signature validation for trusted internal services

## API Reference

### Middleware Function

```go
func Signature(cfg config.SignatureConfig) gin.HandlerFunc
```

### Helper Functions

```go
// Calculate signature manually
func calculateSignature(method, path string, body []byte, secretKey string) string

// Check if signature was validated in handlers
func IsSignatureValid(c *gin.Context) bool

// Get signature validation result
func GetSignatureValidation(c *gin.Context) (valid bool, signature string)
```

### Configuration Structure

```go
type SignatureConfig struct {
    Enabled      bool     // Enable/disable signature validation
    SecretKey    string   // HMAC secret key
    HeaderName   string   // Header name for signature (default: "X-Signature")
    Methods      []string // HTTP methods to validate
    ExcludePaths []string // Paths to exclude from validation
}
```

## Related Documentation

- [Crypto Middleware](./CRYPTO_MIDDLEWARE.md) - AES-256-CBC encryption
- [JWT Authentication](./AUTH.md) - Token-based authentication
- [Rate Limiting](./RATE_LIMITING.md) - Request throttling
- [Security Best Practices](./SECURITY.md) - Overall security guidelines

## Changelog

### Version 1.0.0 (2025-01-03)

- Initial implementation of HMAC-SHA256 signature middleware
- Support for configurable HTTP methods
- Path exclusion patterns (exact, wildcard, catch-all)
- Custom header names
- Comprehensive unit tests
- Production-ready configuration

## License

Copyright (c) 2025 GameLink. All rights reserved.
