# GameLink Client - Security Quick Reference

> Quick guide for implementing security features in GameLink Client

## 🚨 Quick Start

### 1. Import Security Utilities

```typescript
import {
    sanitizeInput,
    isValidUrl,
    validateFileUpload,
    createRateLimiter,
    logSecurityEvent
} from '@/lib/security';
```

### 2. Common Security Tasks

#### User Input Sanitization
```typescript
// ❌ Unsafe
<div>{userInput}</div>

// ✅ Safe (React auto-escapes, but sanitize anyway)
<div>{sanitizeInput(userInput)}</div>
```

#### URL Validation
```typescript
if (isValidUrl(userUrl) && isExternalLinkSafe(userUrl)) {
    window.open(userUrl, '_blank');
}
```

#### File Upload Validation
```typescript
const validation = validateFileUpload(file);
if (!validation.valid) {
    toast.error(validation.error);
    return;
}
```

#### Rate Limiting
```typescript
const limiter = createRateLimiter();
if (!limiter.canMakeRequest()) {
    toast.error('Too many requests');
    return;
}
```

## 📋 Environment Variables

### Development (.env.local)
```bash
VITE_CRYPTO_ENABLED=false
VITE_SECURITY_LOGGING_ENABLED=true
VITE_CSP_ENABLED=true
```

### Production (.env.production)
```bash
VITE_CRYPTO_ENABLED=true
VITE_CRYPTO_SECRET_KEY=<32-char-random-key>
VITE_CRYPTO_IV=<16-char-random-iv>
VITE_SECURITY_LOGGING_ENABLED=false
VITE_CSP_ENABLED=true
VITE_REQUIRE_HTTPS=true
```

## 🛡️ Security Checklist

### Before Committing Code
- [ ] All user input is sanitized
- [ ] URLs are validated before use
- [ ] File uploads are validated
- [ ] No sensitive data in logs
- [ ] Security tests pass
- [ ] No new vulnerabilities introduced

### Before Deployment
- [ ] `npm audit` shows 0 vulnerabilities
- [ ] Encryption enabled in production
- [ ] CSP headers configured
- [ ] HTTPS enforced
- [ ] Security logging disabled
- [ ] Session timeouts configured

## 🧪 Testing

```bash
# Run security tests
npm run test -- security

# Check vulnerabilities
npm run security:audit

# Full test suite
npm run test:coverage
```

## ⚠️ Common Mistakes

### ❌ Don't Do This

```typescript
// Don't use dangerouslySetInnerHTML with user input
<div dangerouslySetInnerHTML={{ __html: userInput }} />

// Don't trust URLs without validation
<a href={userUrl}>Link</a>

// Don't store sensitive data
localStorage.setItem('password', password);

// Don't ignore file validation
uploadFile(file); // Always validate first!
```

### ✅ Do This Instead

```typescript
// Sanitize HTML first
<div dangerouslySetInnerHTML={{ __html: sanitizeHTML(userInput) }} />

// Validate URLs
<a href={isValidUrl(userUrl) ? userUrl : '#'}>Link</a>

// Never store sensitive data
// Use secure auth store instead

// Always validate files
const v = validateFileUpload(file);
if (!v.valid) return showError(v.error);
```

## 📚 Key Security Functions

### Input Validation
- `sanitizeInput(input)` - Sanitize user input
- `sanitizeHTML(html)` - Sanitize HTML content
- `isValidUrl(url)` - Validate URL protocol
- `isTrustedDomain(url)` - Check trusted domains

### File Security
- `validateFileUpload(file)` - Validate file upload

### Session & Auth
- `isSessionValid(lastActivity)` - Check session timeout
- `getTokenStorageMethod()` - Get storage method
- `isTokenStorageSecure()` - Check if secure

### Rate Limiting
- `createRateLimiter()` - Create rate limiter instance
- `rateLimiter.canMakeRequest()` - Check if allowed
- `rateLimiter.getRemainingRequests()` - Get remaining count

### Environment
- `isSecureEnvironment()` - Check HTTPS
- `getSecurityHeaders()` - Get security headers
- `getCSPPolicy()` - Get CSP policy

### Logging
- `logSecurityEvent(event, data)` - Log security events (dev only)

## 🔒 Security Headers

The following security headers are configured:

```typescript
{
  "X-Content-Type-Options": "nosniff",
  "X-Frame-Options": "DENY",
  "X-XSS-Protection": "1; mode=block",
  "Strict-Transport-Security": "max-age=31536000; includeSubDomains"
}
```

## 🚨 Incident Response

If you discover a security issue:

1. **Stop** - Don't use the affected feature
2. **Report** - Email security@gamelink.com immediately
3. **Document** - Write down what you found
4. **Fix** - Create a private branch with the fix
5. **Test** - Thoroughly test the fix
6. **Deploy** - Hotfix if critical

## 📖 Full Documentation

- **Complete Guide**: `docs/CLIENT_SECURITY_GUIDE.md`
- **Implementation**: `docs/SECURITY_IMPLEMENTATION_SUMMARY.md`
- **Configuration**: `.env.security`
- **Utilities**: `client/src/lib/security.ts`

## 🔗 Quick Links

- Run Security Audit: `npm run security:audit`
- Run Tests: `npm run test -- security`
- View Security Score: 88/100
- Known Vulnerabilities: 0

---

**Last Updated**: 2026-01-18
**Security Score**: 88/100
**Status**: ✅ All Security Measures Implemented
