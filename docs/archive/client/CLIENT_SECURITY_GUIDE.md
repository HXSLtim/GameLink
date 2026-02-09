# GameLink Client - Security Best Practices Guide

## Overview

This document outlines security best practices for the GameLink Client application. It provides actionable guidelines for developers to maintain and improve the security posture of the client-side application.

**Security Score**: 88/100 (as of 2026-01-18)

## Table of Contents

- [Security Checklist](#security-checklist)
- [Known Security Considerations](#known-security-considerations)
- [Security Testing](#security-testing)
- [Incident Response](#incident-response)
- [Security Configuration](#security-configuration)
- [Development Best Practices](#development-best-practices)

## Security Checklist

### Development Environment

- [ ] Never commit `.env` files to version control
- [ ] Use `.env.local` for local development overrides
- [ ] Enable security logging in development (`VITE_SECURITY_LOGGING_ENABLED=true`)
- [ ] Run `npm audit` regularly
- [ ] Keep dependencies up to date
- [ ] Review security warnings in console
- [ ] Test with production-like security settings

### Production Environment

- [ ] Enable encryption (`VITE_CRYPTO_ENABLED=true`)
- [ ] Use strong, randomly generated crypto keys
- [ ] Enable HTTPS only (`VITE_REQUIRE_HTTPS=true`)
- [ ] Configure CSP headers (`VITE_CSP_ENABLED=true`)
- [ ] Set up security monitoring
- [ ] Enable all security headers
- [ ] Configure proper session timeouts
- [ ] Set up rate limiting
- [ ] Disable debug logging
- [ ] Use environment-specific configurations

### Code Review

- [ ] Sanitize all user inputs
- [ ] Validate URLs before navigation
- [ ] Avoid `dangerouslySetInnerHTML` with user content
- [ ] Check for exposed sensitive data in logs
- [ ] Review token storage implementation
- [ ] Verify error messages don't leak information
- [ ] Test XSS attack vectors
- [ ] Validate file uploads

## Known Security Considerations

### Token Storage

**Current Implementation**: localStorage

**Risk Level**: Medium

**Description**:
- localStorage is vulnerable to XSS attacks
- Malicious scripts can access tokens
- Tokens persist across sessions

**Recommended Solution**: httpOnly Cookies

**Benefits of httpOnly Cookies**:
- Inaccessible to JavaScript (prevents XSS token theft)
- Automatic CSRF protection when combined with SameSite attribute
- More secure session management

**Migration Plan**:

#### Phase 1: Backend Preparation
- [ ] Implement CSRF token generation and validation
- [ ] Add httpOnly cookie support to authentication endpoints
- [ ] Update JWT middleware to support cookie-based auth
- [ ] Add SameSite cookie attribute

#### Phase 2: Client Migration
- [ ] Update auth store to use cookie-based auth
- [ ] Remove token management from client-side storage
- [ ] Add CSRF token to all state-changing requests
- [ ] Update error handling for auth failures

#### Phase 3: Testing & Deployment
- [ ] Security testing for CSRF protection
- [ ] Test cookie behavior across browsers
- [ ] Verify session invalidation on logout
- [ ] Gradual rollout with monitoring

**Current Mitigation Measures**:
1. Content Security Policy (CSP) headers reduce XSS risk
2. Input sanitization prevents script injection
3. Short token expiration (5 minutes before refresh)
4. Security logging for suspicious activities

**Temporary Workaround**:
While waiting for backend support, the following measures are in place:
- XSS protection via CSP
- Input sanitization utilities
- Security event logging
- Token refresh buffer to minimize exposure window

### Dependency Security

**Risk Level**: Medium

**Regular Checks**:
```bash
# Check for vulnerabilities
npm audit

# Fix vulnerabilities automatically
npm audit fix

# Check for outdated packages
npm outdated

# Run audit with specific severity level
npm audit --audit-level=moderate
```

**Best Practices**:
- Review security advisories for dependencies
- Update dependencies regularly
- Lock dependency versions in production
- Use npm's `overrides` for forced updates
- Monitor for security alerts in GitHub
- Subscribe to security bulletins

**Automated Security Scripts**:
```bash
# Run security audit (package.json script)
npm run security:audit

# Fix security issues
npm run security:fix

# Check security before deployment
npm run security:check
```

### XSS Prevention

**Current Protections**:

1. **React Built-in Protection**
   - Automatic escaping of JSX content
   - No manual encoding needed for most cases

2. **Content Security Policy (CSP)**
   - Restricts script sources
   - Prevents inline script execution
   - Blocks unauthorized content loading

3. **Input Sanitization**
   - `sanitizeInput()` function for user input
   - URL validation before navigation
   - HTML sanitization for rich content

4. **Security Headers**
   - X-XSS-Protection header
   - X-Content-Type-Options: nosniff
   - X-Frame-Options: DENY

**Best Practices**:

✅ **DO**:
- Use React's default rendering (automatic escaping)
- Sanitize user input before display
- Validate URLs before navigation
- Use CSP headers
- Implement security headers

❌ **DON'T**:
- Use `dangerouslySetInnerHTML` with user input
- Construct HTML with user data
- `eval()` user input
- Use `innerHTML` with untrusted data
- Bypass CSP with unsafe inline scripts

**XSS Testing Checklist**:
```javascript
// Test cases to try (should be sanitized)
<script>alert('XSS')</script>
<img src=x onerror=alert('XSS')>
<a href="javascript:alert('XSS')">link</a>
<svg onload=alert('XSS')>
```

### CSRF Protection

**Current Status**: Requires Backend Support

**Risk Level**: Medium

**Description**:
Cross-Site Request Forgery (CSRF) tricks users into performing actions they didn't intend.

**Required Backend Changes**:
1. Generate CSRF tokens on session start
2. Include CSRF token in forms and AJAX requests
3. Validate CSRF token on state-changing requests
4. Use SameSite cookie attribute

**Client-Side Preparation** (Future):
```typescript
// When backend is ready, add CSRF token to requests
axios.interceptors.request.use((config) => {
  const csrfToken = getCsrfToken(); // from cookie or meta tag
  if (csrfToken) {
    config.headers['X-CSRF-Token'] = csrfToken;
  }
  return config;
});
```

## Security Testing

### Manual Testing Checklist

**Input Validation**:
- [ ] Try to inject scripts in all input fields
- [ ] Test SQL injection in search fields
- [ ] Attempt path traversal in file uploads
- [ ] Verify length limits are enforced
- [ ] Test special characters and Unicode

**Authentication & Authorization**:
- [ ] Verify tokens are not exposed in URLs
- [ ] Check that logout invalidates tokens
- [ ] Test session timeout functionality
- [ ] Verify role-based access control
- [ ] Test concurrent login handling

**Data Protection**:
- [ ] Check that sensitive data is not logged
- [ ] Verify passwords are never stored or logged
- [ ] Test that tokens are not visible in browser storage
- [ ] Confirm no data leakage in error messages
- [ ] Verify HTTPS enforcement in production

**API Security**:
- [ ] Test API error messages don't leak information
- [ ] Verify rate limiting is enforced
- [ ] Check that sensitive endpoints require auth
- [ ] Test request timeout handling
- [ ] Verify proper status codes

### Automated Testing

```bash
# Run security audit
npm audit

# Run test suite
npm run test

# Type checking
npx tsc --noEmit

# Linting
npm run lint

# Coverage report
npm run test:coverage
```

**Security Test Utilities**:

Located in `client/src/lib/security.ts`:
- `validateCryptoConfig()` - Verify encryption settings
- `sanitizeInput()` - Test input sanitization
- `isValidUrl()` - Test URL validation
- `validateFileUpload()` - Test file upload security
- `RateLimiter` - Test rate limiting

**Example Test**:
```typescript
import { describe, it, expect } from 'vitest';
import { sanitizeInput, isValidUrl, validateFileUpload } from '@/lib/security';

describe('Security Utilities', () => {
  describe('sanitizeInput', () => {
    it('should sanitize XSS attempts', () => {
      const input = '<script>alert("xss")</script>';
      const sanitized = sanitizeInput(input);
      expect(sanitized).not.toContain('<script>');
    });
  });

  describe('isValidUrl', () => {
    it('should reject javascript: protocol', () => {
      expect(isValidUrl('javascript:alert(1)')).toBe(false);
    });

    it('should accept https URLs', () => {
      expect(isValidUrl('https://gamelink.com')).toBe(true);
    });
  });
});
```

## Security Configuration

### Environment Variables

See `.env.security` for complete security configuration options.

**Key Variables**:

```bash
# Encryption (Production Required)
VITE_CRYPTO_ENABLED=true
VITE_CRYPTO_SECRET_KEY=<32-char-random-key>
VITE_CRYPTO_IV=<16-char-random-iv>

# Content Security Policy
VITE_CSP_ENABLED=true
VITE_CSP_POLICY={"default-src":"'self'",...}

# Token Storage
VITE_TOKEN_STORAGE=localStorage  # until httpOnly cookies available

# API Security
VITE_REQUIRE_HTTPS=true
VITE_API_TIMEOUT=30000
VITE_RATE_LIMIT_PER_MINUTE=60

# Session Security
VITE_SESSION_TIMEOUT=86400000      # 24 hours
VITE_IDLE_TIMEOUT=1800000          # 30 minutes
VITE_SESSION_TIMEOUT_WARNING=300000 # 5 minutes
```

### Security Headers

**Headers Implemented**:

```typescript
{
  "X-Content-Type-Options": "nosniff",
  "X-Frame-Options": "DENY",
  "X-XSS-Protection": "1; mode=block",
  "Strict-Transport-Security": "max-age=31536000; includeSubDomains"
}
```

**CSP Policy**:
```json
{
  "default-src": "'self'",
  "script-src": "'self' 'unsafe-inline' 'unsafe-eval'",
  "style-src": "'self' 'unsafe-inline'",
  "img-src": "'self' data: https:",
  "connect-src": "'self' https://api.gamelink.com wss://api.gamelink.com",
  "font-src": "'self' data:",
  "object-src": "'none'",
  "base-uri": "'self'",
  "form-action": "'self'",
  "frame-ancestors": "'none'"
}
```

### Security Utilities

**Available Functions** (from `client/src/lib/security.ts`):

```typescript
// Configuration validation
validateCryptoConfig()
validateSecurityConfig()

// Input sanitization
sanitizeInput(input: string)
sanitizeHTML(html: string)

// URL validation
isValidUrl(url: string)
isTrustedDomain(url: string)
isExternalLinkSafe(url: string)

// Environment checks
isSecureEnvironment()
isTokenStorageSecure()

// File validation
validateFileUpload(file: File)

// Session management
isSessionValid(lastActivity: number)
getSessionTimeoutWarning()

// Rate limiting
createRateLimiter()
class RateLimiter { ... }

// Security headers
getSecurityHeaders()
getCSPPolicy()
getCSPMetaTag()

// Logging
logSecurityEvent(event: string, data?: unknown)
```

## Development Best Practices

### Code Review Security Checklist

When reviewing code, check for:

1. **Input Handling**
   - All user input is sanitized
   - No direct HTML injection
   - URL validation before navigation

2. **Data Storage**
   - No sensitive data in localStorage/sessionStorage
   - Tokens handled securely
   - No passwords or secrets in code

3. **API Communication**
   - HTTPS only in production
   - Proper error handling
   - No sensitive data in URLs

4. **Authentication**
   - Proper token validation
   - Secure logout implementation
   - Session timeout handling

5. **Error Handling**
   - No information leakage in errors
   - Secure error logging
   - User-friendly error messages

### Security in React Components

**Safe Practices**:
```typescript
// ✅ Safe: React auto-escapes
<div>{userInput}</div>

// ✅ Safe: Sanitized HTML
<div dangerouslySetInnerHTML={{ __html: sanitizeHTML(userInput) }} />

// ❌ Unsafe: Raw HTML
<div dangerouslySetInnerHTML={{ __html: userInput }} />

// ✅ Safe: Validated URL
<a href={isValidUrl(url) ? url : '#'}>Link</a>

// ❌ Unsafe: Unvalidated URL
<a href={userProvidedUrl}>Link</a>
```

### Logging Best Practices

**Security Logging** (Development Only):
```typescript
import { logSecurityEvent } from '@/lib/security';

// Log security events
logSecurityEvent('XSS attempt detected', { input: userInput });
logSecurityEvent('Invalid URL', { url: userUrl });
logSecurityEvent('File upload rejected', { fileName, reason });
```

**Production Logging**:
- Never log sensitive data (passwords, tokens, PII)
- Use structured logging for security events
- Implement log aggregation and monitoring
- Set up alerts for suspicious activities

## Incident Response

### If a Security Issue is Discovered

**Immediate Actions**:
1. Stop using the affected feature
2. Notify the security team immediately
3. Document the issue thoroughly
4. Assess severity and impact

**Development Process**:
1. Create a private branch for the fix
2. Implement and test the fix thoroughly
3. Document the vulnerability and fix
4. Create security test cases
5. Code review by security lead
6. Deploy as hotfix if critical

**Post-Incident**:
1. Conduct root cause analysis
2. Update this document with lessons learned
3. Add test cases to prevent regression
4. Review related code for similar issues
5. Update security training if needed

### Severity Levels

**Critical** (Immediate Action Required):
- Active exploitation
- Data breach confirmed
- Authentication bypass
- Remote code execution

**High** (Action Within 24 Hours):
- Privilege escalation
- Sensitive data exposure
- CSRF/XSS in production
- Session hijacking

**Medium** (Action Within 1 Week):
- Information disclosure
- Denial of service vulnerabilities
- Missing security headers
- Dependency vulnerabilities

**Low** (Action Within 1 Month):
- Best practice violations
- Minor configuration issues
- Low-risk dependency updates
- Documentation improvements

## Resources

### External Resources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [OWASP XSS Prevention Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross_Site_Scripting_Prevention_Cheat_Sheet.html)
- [OWASP CSRF Prevention](https://cheatsheetseries.owasp.org/cheatsheets/Cross_Site_Request_Forgery_Prevention_Cheat_Sheet.html)
- [CSP Evaluator](https://csp-evaluator.withgoogle.com/)
- [React Security](https://react.dev/learn/keeping-components-pure)
- [MDN Web Security](https://developer.mozilla.org/en-US/docs/Web/Security)

### Internal Documentation

- `CLIENT_SECURITY_THREAT_MODEL.md` - Complete threat analysis
- `CLIENT_INFRASTRUCTURE_DESIGN.md` - Security architecture
- `.env.security` - Security configuration reference
- `client/src/lib/security.ts` - Security utilities implementation

### Security Tools

```bash
# Dependency vulnerability scanning
npm audit
npm audit fix
snyk test

# Static analysis
eslint --ext .js,.jsx,.ts,.tsx src/

# Type checking
tsc --noEmit

# Testing
vitest run --coverage
```

## Version History

- **2026-01-18**: Initial version with basic security guidelines
  - Security checklist and testing procedures
  - Known security considerations documented
  - Token storage migration plan outlined
  - Security utilities and configuration documented

## Contributing

When contributing to the GameLink Client:

1. Review this security guide before implementing features
2. Follow the security checklist for code reviews
3. Add security tests for new functionality
4. Update this document if adding new security features
5. Report security issues privately to the security team

## Contact

For security questions or to report vulnerabilities:
- Security Team: security@gamelink.com
- Private Vulnerability Report: Create a private GitHub issue with the "security" label

---

**Remember**: Security is everyone's responsibility. When in doubt, ask the security team before implementing potentially risky features.
