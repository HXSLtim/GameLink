# Environment Variables Setup Guide

This document explains how to properly configure environment variables for the GameLink project to ensure security and avoid hardcoded secrets.

## Quick Start

### Backend (Go API)

1. **Copy the example environment file:**
   ```bash
   cp .env.example .env
   ```

2. **Edit `.env` with your actual values:**
   ```bash
   # Generate secure keys
   openssl rand -base64 32  # For CRYPTO_SECRET_KEY (32 bytes)
   openssl rand -base64 16  # For CRYPTO_IV (16 bytes)
   openssl rand -base64 32  # For JWT_SECRET_KEY
   openssl rand -base64 24  # For SUPER_ADMIN_PASSWORD
   ```

3. **Update the `.env` file with the generated values:**
   ```bash
   CRYPTO_ENABLED=true
   CRYPTO_SECRET_KEY=<your-32-byte-secret-key>
   CRYPTO_IV=<your-16-byte-iv>
   JWT_SECRET_KEY=<your-jwt-secret>
   SUPER_ADMIN_PASSWORD=<your-admin-password>
   DB_DSN=host=localhost port=5432 user=gamelink password=<db-password> dbname=gamelink sslmode=disable
   REDIS_PASSWORD=<redis-password>
   ```

4. **Start the application:**
   ```bash
   cd api
   go run cmd/main.go
   ```

### Frontend (Admin Panel)

1. **Copy the example environment file:**
   ```bash
   cd admin
   cp .env.example .env
   ```

2. **Edit `.env` with your configuration:**
   ```bash
   VITE_API_BASE_URL=http://localhost:8080
   VITE_CRYPTO_ENABLED=false  # Set to true if backend has encryption enabled
   ```

3. **Start the development server:**
   ```bash
   npm run dev
   ```

## Security Requirements

### Production Environment

In production mode, the following validations are enforced:

1. **Encryption MUST be enabled:**
   - `CRYPTO_ENABLED=true` is required
   - `CACHE_TYPE` must be `redis` (not `memory`)

2. **Required secrets:**
   - `DB_DSN` - Database connection string
   - `JWT_SECRET_KEY` - At least 32 characters recommended
   - `SUPER_ADMIN_EMAIL` - Valid email address
   - `SUPER_ADMIN_PASSWORD` - At least 8 characters with uppercase, lowercase, number, and special character
   - `CRYPTO_SECRET_KEY` - Exactly 32 bytes for AES-256-CBC
   - `CRYPTO_IV` - Exactly 16 bytes

3. **Forbidden values:**
   - Hardcoded default keys are rejected
   - Empty secrets are rejected when encryption is enabled

### Development Environment

In development mode:
- Encryption can be disabled (`CRYPTO_ENABLED=false`)
- Weaker passwords are allowed
- Some validations are relaxed

## Environment Variables Reference

### Application

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_ENV` | `development` | Environment: development, production, staging |
| `SERVICE_PORT` | `8080` | Server port |
| `ENABLE_SWAGGER` | `false` | Enable Swagger documentation |

### Database

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_TYPE` | `postgres` | Database type: postgres, sqlite |
| `DB_DSN` | - | Database connection string |
| `DB_MAX_CONNS` | `25` | Maximum database connections |
| `DB_MAX_IDLE` | `5` | Maximum idle connections |

### Cache

| Variable | Default | Description |
|----------|---------|-------------|
| `CACHE_TYPE` | `memory` | Cache type: redis, memory |
| `REDIS_ADDR` | `127.0.0.1:6379` | Redis address |
| `REDIS_PASSWORD` | - | Redis password |
| `REDIS_DB` | `0` | Redis database number |

### Security - Encryption

| Variable | Default | Description |
|----------|---------|-------------|
| `CRYPTO_ENABLED` | `false` | Enable request/response encryption |
| `CRYPTO_SECRET_KEY` | - | 32-byte AES-256-CBC key |
| `CRYPTO_IV` | - | 16-byte initialization vector |
| `CRYPTO_METHODS` | `POST,PUT,PATCH` | HTTP methods to encrypt |
| `CRYPTO_EXCLUDE_PATHS` | `/api/v1/health,...` | Paths to exclude from encryption |
| `CRYPTO_USE_SIGNATURE` | `true` | Enable SHA-256 signature |

### Security - Authentication

| Variable | Default | Description |
|----------|---------|-------------|
| `JWT_SECRET_KEY` | - | JWT signing secret (min 16 chars) |
| `JWT_TOKEN_TTL_HOURS` | `24` | Token time-to-live in hours |

### Super Admin

| Variable | Default | Description |
|----------|---------|-------------|
| `SUPER_ADMIN_EMAIL` | `admin@gamelink.com` | Super admin email |
| `SUPER_ADMIN_PASSWORD` | - | Super admin password |
| `SUPER_ADMIN_NAME` | `Super Admin` | Super admin display name |
| `SUPER_ADMIN_PHONE` | - | Super admin phone number |

## Generating Secure Keys

### Using OpenSSL

```bash
# Crypto secret key (32 bytes, base64 encoded = 44 chars)
openssl rand -base64 32

# Crypto IV (16 bytes, base64 encoded = 24 chars)
openssl rand -base64 16

# JWT secret key (32+ bytes)
openssl rand -base64 32

# Admin password (24 bytes, base64 encoded)
openssl rand -base64 24
```

### Using Python

```python
import secrets
import base64

# 32-byte crypto key
print(base64.b64encode(secrets.token_bytes(32)).decode())

# 16-byte IV
print(base64.b64encode(secrets.token_bytes(16)).decode())

# 32-byte JWT secret
print(secrets.token_urlsafe(32))

# 24-byte password
print(base64.b64encode(secrets.token_bytes(24)).decode())
```

## Validation Errors

Common validation errors and how to fix them:

### "CRYPTO_SECRET_KEY is required when encryption is enabled"
**Fix:** Set `CRYPTO_SECRET_KEY` environment variable with a 32-byte value.

### "CRYPTO_SECRET_KEY must be exactly 32 bytes"
**Fix:** Generate a new key with `openssl rand -base64 32` and use it directly (not base64 decoded).

### "CRYPTO_IV is using a hardcoded default value"
**Fix:** Generate a new IV with `openssl rand -base64 16` and set it in `CRYPTO_IV`.

### "SUPER_ADMIN_PASSWORD must be at least 8 characters in production"
**Fix:** Use a stronger password with at least 8 characters including uppercase, lowercase, number, and special character.

### "CRYPTO_ENABLED must be true in production"
**Fix:** Set `CRYPTO_ENABLED=true` in production environment.

## Best Practices

1. **Never commit `.env` files** to version control
2. **Use different credentials** for different environments
3. **Rotate keys regularly** (at least every 90 days)
4. **Use strong, unique passwords** for all services
5. **Store secrets securely** in production using:
   - Kubernetes Secrets
   - AWS Secrets Manager
   - Azure Key Vault
   - HashiCorp Vault
6. **Limit `.env` file permissions** (chmod 600)
7. **Document key rotation procedures**

## Troubleshooting

### Application won't start with "configuration validation failed"

Check the error message for specific validation failures:
- Are all required secrets set?
- Are crypto keys the correct length?
- Is the password strong enough for production?

### Encryption errors in frontend

If you see "Encryption configuration error" in the browser console:
- Check that `VITE_CRYPTO_ENABLED` matches backend `CRYPTO_ENABLED`
- Verify crypto keys match between frontend and backend
- Ensure frontend has access to environment variables

### Database connection errors

- Verify `DB_DSN` is correctly formatted
- Check database is running and accessible
- Ensure password doesn't contain URL special characters (like `%`)

## Additional Resources

- [Deployment Guide](../README.md#deployment)
- [Security Best Practices](../CLAUDE.md#security)
- [Production Deployment](../scripts/README.md)
