# Error Monitoring with Sentry

This document describes the error monitoring integration with Sentry for the GameLink Admin Panel.

## Overview

The admin panel uses Sentry for production error monitoring. Errors are automatically captured and reported to Sentry when they occur in the application.

## Configuration

### Environment Variables

Configure Sentry in your `.env` file:

```bash
# Enable Sentry error monitoring
VITE_SENTRY_ENABLED=true

# Sentry DSN (Data Source Name) - Get from Sentry dashboard
VITE_SENTRY_DSN=https://examplePublicKey@o0.ingest.sentry.io/0

# Environment name (e.g., production, staging, development)
VITE_SENTRY_ENVIRONMENT=production

# Error sample rate (0.0 to 1.0, recommended: 1.0 for production)
VITE_SENTRY_SAMPLE_RATE=1.0

# Performance monitoring sample rate (0.0 to 1.0, recommended: 0.1 for production)
VITE_SENTRY_TRACES_SAMPLE_RATE=0.1
```

### Getting a Sentry DSN

1. Create a Sentry account at https://sentry.io
2. Create a new project
3. Select "React" as the platform
4. Copy the DSN from the project settings

## Features

### Automatic Error Capture

Errors are automatically captured from:

- **React Error Boundary**: Catches React component errors
- **Uncaught Exceptions**: Global JavaScript errors
- **Unhandled Promise Rejections**: Promise rejections without catch handlers
- **Network Errors**: Failed API requests (when using the monitoring utilities)

### Manual Error Reporting

Use the monitoring utilities to manually report errors:

```typescript
import { captureException, captureMessage, addBreadcrumb } from '@/utils/monitoring';

// Capture an exception
try {
  // Some operation
} catch (error) {
  captureException(error as Error, {
    context: 'Additional context',
  });
}

// Capture a message
captureMessage('Something important happened', 'warning', {
  customData: 'value',
});

// Add a breadcrumb for context
addBreadcrumb('user-action', 'User clicked button', {
  buttonId: 'submit-btn',
});
```

### User Context

Set user context to track which users experience errors:

```typescript
import { setUser, clearUser } from '@/utils/monitoring';

// Set user when they log in
setUser(userId, userEmail, username);

// Clear user when they log out
clearUser();
```

### Tags and Extra Data

Add custom tags and data to help filter errors:

```typescript
import { setTag, setExtra } from '@/utils/monitoring';

// Set a tag for filtering
setTag('page', 'dashboard');
setTag('feature', 'user-management');

// Set extra data
setExtra('customData', { some: 'data' });
```

## ErrorBoundary Integration

The `ErrorBoundary` component automatically sends errors to Sentry in production:

```tsx
import { ErrorBoundary } from '@/components/ErrorBoundary';

<ErrorBoundary>
  <App />
</ErrorBoundary>
```

## Development vs Production

- **Development**: Errors are logged to the console only
- **Production**: Errors are sent to Sentry (if enabled)

Sentry is only initialized when:
- `VITE_SENTRY_ENABLED=true`
- Not in development mode
- Valid DSN is provided

## Source Maps

Source maps are automatically uploaded to Sentry during production builds when:
- `SENTRY_AUTH_TOKEN` environment variable is set
- Running in production mode

## Ignored Errors

The following errors are ignored by default to reduce noise:

- Random plugins/extensions errors
- Facebook bugged errors
- Safari QuotaExceededError
- Non-Error promise rejections
- ResizeObserver loop limit exceeded

## Performance Monitoring

Performance monitoring is available with configurable sample rates. This helps track:

- Page load times
- API request durations
- Component render times

Configure with `VITE_SENTRY_TRACES_SAMPLE_RATE`.

## Security Notes

1. Never commit actual `.env` file with credentials
2. DSN is exposed in browser, but this is intentional for client-side error reporting
3. Sensitive data is automatically filtered from error reports
4. Use appropriate sample rates to manage quota

## Troubleshooting

### Errors not appearing in Sentry

1. Check `VITE_SENTRY_ENABLED=true` in `.env`
2. Verify DSN is correct
3. Check browser console for initialization errors
4. Ensure you're in production mode

### Source maps not working

1. Set `SENTRY_AUTH_TOKEN` environment variable
2. Build with `pnpm build`
3. Check Sentry project settings for uploaded artifacts

## Additional Resources

- [Sentry Documentation](https://docs.sentry.io/)
- [Sentry React Integration](https://docs.sentry.io/platforms/javascript/guides/react/)
- [Error Monitoring Best Practices](https://docs.sentry.io/product/alerts/)
