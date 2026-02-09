# Client Deployment Guide

## Prerequisites

- Node.js 18+ or 20+
- npm or yarn
- Valid environment configuration

## Pre-Build Validation

Before building for production, run all checks:

```bash
# Validate environment variables
npm run validate:env

# Run pre-build checks (TypeScript, ESLint, tests)
npm run prebuild

# Run complete CI pipeline locally
npm run ci
```

## Building for Production

```bash
# Standard build
npm run build

# Build with validation
npm run prebuild && npm run build
```

## Deployment Checklist

- [ ] All tests passing
- [ ] TypeScript compilation successful
- [ ] ESLint check passed
- [ ] Security audit passed
- [ ] Environment variables validated
- [ ] Production build successful
- [ ] Build artifacts tested
- [ ] Health checks passing

## CI/CD Pipeline

The client has automated CI/CD via GitHub Actions:
- Runs on every push to main/dev
- Executes TypeScript, ESLint, tests
- Generates coverage reports
- Uploads build artifacts
- Runs security scans

### Workflow Triggers

The `.github/workflows/client-ci.yml` workflow runs on:
- Push to `main` or `dev` branches (changes in `client/` directory)
- Pull requests targeting `main` or `dev` branches

### Jobs

1. **Client Quality Checks** (Node.js 18.x and 20.x)
   - TypeScript type check
   - ESLint validation
   - Unit tests with coverage
   - Security audit
   - Production build
   - Coverage upload to Codecov
   - Build artifact upload

2. **Integration Tests**
   - Runs against mock backend service
   - Tests API connectivity

3. **Security Scan**
   - Trivy vulnerability scanner
   - SARIF results upload to GitHub Security

### Quality Gates

- **TypeScript errors**: Build fails
- **ESLint errors**: Build fails
- **Test failures**: Build fails
- **Security vulnerabilities**: Displayed as warnings (build continues)

## Environment Variables

### Required Variables

```bash
VITE_API_BASE_URL=<API endpoint URL>
```

### Optional Variables

```bash
# Encryption
VITE_CRYPTO_ENABLED=true
VITE_CRYPTO_SECRET_KEY=<32-character secret>
VITE_CRYPTO_IV=<16-character IV>

# Security
VITE_CSP_ENABLED=true
VITE_TOKEN_STORAGE=localStorage|sessionStorage
```

### Environment Files

- `.env` - Local development (not committed)
- `.env.development` - Development environment
- `.env.production` - Production environment
- `.env.example` - Example configuration (committed)
- `.env.security` - Security-specific configuration

## Health Checks

The client includes a health check system in `src/lib/health.ts`:

```typescript
import { performHealthChecks, getHealthReport } from '@/lib/health';

// Run health checks
const health = await performHealthChecks();
console.log(health.status); // 'healthy' | 'degraded' | 'unhealthy'

// Get formatted report
const report = await getHealthReport();
console.log(report);
```

### Health Check Components

1. **API Connectivity**: Tests backend `/api/v1/healthz` endpoint
2. **LocalStorage**: Verifies browser storage availability
3. **Crypto Configuration**: Validates encryption settings

## Monitoring

After deployment:

### Performance Metrics
- Monitor page load times
- Track API response times
- Check bundle size impact

### Error Tracking
- Review browser console errors
- Monitor API failure rates
- Check crash reports

### Security Monitoring
- Review security scan results
- Monitor for new vulnerabilities
- Check dependency updates

## Troubleshooting

### Build Fails

**Problem**: TypeScript compilation errors

```bash
# Solution: Run type check to see all errors
npx tsc --noEmit
```

**Problem**: ESLint errors

```bash
# Solution: Run lint to see specific issues
npm run lint

# Fix auto-fixable issues
npm run lint -- --fix
```

**Problem**: Test failures

```bash
# Solution: Run tests with detailed output
npm run test:run -- --reporter=verbose
```

### Environment Issues

**Problem**: Missing environment variables

```bash
# Solution: Validate environment
npm run validate:env

# Create missing file from example
cp .env.example .env.production
# Edit .env.production with actual values
```

**Problem**: Incorrect crypto configuration

```bash
# Solution: Validate crypto key lengths
# VITE_CRYPTO_SECRET_KEY must be exactly 32 characters
# VITE_CRYPTO_IV must be exactly 16 characters
```

### Runtime Issues

**Problem**: API connectivity failures

```bash
# Solution: Check VITE_API_BASE_URL
# Ensure backend is accessible
# Test endpoint: curl https://your-api.com/api/v1/healthz
```

**Problem**: LocalStorage unavailable

```bash
# Solution: Check browser settings
# Ensure cookies/storage enabled
# Consider using sessionStorage instead
```

## Rollback Procedure

If issues occur after deployment:

1. **Identify the problematic commit**
   ```bash
   git log --oneline -10
   ```

2. **Revert to previous stable version**
   ```bash
   git revert <commit-hash>
   # OR
   git checkout <previous-tag>
   ```

3. **Run full test suite**
   ```bash
   npm run ci
   ```

4. **Redeploy**
   ```bash
   npm run build
   # Deploy dist/ directory
   ```

## Performance Optimization

### Build Size Analysis

```bash
# Build with bundle analysis
npm run build:analyze

# View report
# Open dist/stats.html in browser
```

### Code Splitting

The client uses dynamic imports for code splitting:
- Route-based splitting (React Router)
- Component-based splitting (lazy loading)
- Vendor splitting (Vite default)

### Caching Strategy

- Static assets: Cached by filename hashes
- API responses: Configurable via Cache-Control headers
- LocalStorage: User preferences and auth tokens

## Security Best Practices

### Dependency Management

```bash
# Regularly audit dependencies
npm audit

# Fix vulnerabilities
npm audit fix

# Check for outdated packages
npm outdated
```

### Environment Variable Security

- Never commit `.env` files
- Use `.env.example` as template
- Rotate crypto keys periodically
- Use different keys for each environment

### Content Security Policy

```bash
# Enable CSP in production
VITE_CSP_ENABLED=true
```

## Version History

| Date | Version | Changes |
|------|---------|---------|
| 2026-01-18 | 1.0.0 | Initial deployment guide |

## Additional Resources

- [Vite Documentation](https://vitejs.dev/)
- [React Documentation](https://react.dev/)
- [TypeScript Documentation](https://www.typescriptlang.org/)
- [Testing Library](https://testing-library.com/)
- [GameLink Project Docs](../.kiro/steering/00-INDEX.md)

## Support

For deployment issues:
1. Check this guide's troubleshooting section
2. Review CI/CD logs in GitHub Actions
3. Consult team documentation in `.kiro/steering/`
4. Contact DevOps team
