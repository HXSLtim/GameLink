/**
 * Sentry Error Monitoring Integration
 *
 * This module provides production error monitoring using Sentry.
 * It is only initialized when VITE_SENTRY_ENABLED is set to 'true'.
 *
 * @module utils/monitoring/sentry
 */

import * as Sentry from '@sentry/react';
import type { Breadcrumb, ErrorEvent, EventHint } from '@sentry/react';

type ErrorContext = Record<string, unknown>;

/**
 * Sentry configuration options
 */
export interface SentryConfig {
  /** Sentry DSN (Data Source Name) */
  dsn: string;
  /** Environment name (e.g., production, staging) */
  environment?: string;
  /** Sample rate for error capture (0-1) */
  sampleRate?: number;
  /** Sample rate for performance tracing (0-1) */
  tracesSampleRate?: number;
  /** Release version */
  release?: string;
  /** Application name */
  appName?: string;
  /** Before send callback for error filtering */
  beforeSend?: (event: ErrorEvent, hint: EventHint) => ErrorEvent | null;
  /** Before breadcrumb callback */
  beforeBreadcrumb?: (breadcrumb: Breadcrumb) => Breadcrumb | null;
}

/**
 * Sentry initialization status
 */
let isInitialized = false;

/**
 * Initialize Sentry with configuration
 * Only initializes once, subsequent calls are ignored
 *
 * @param config - Sentry configuration
 */
export function initSentry(config: SentryConfig): void {
  if (isInitialized) {
    console.warn('[Sentry] Already initialized, skipping duplicate initialization');
    return;
  }

  // Only initialize in production mode
  if (import.meta.env.DEV) {
    console.log('[Sentry] Skipping initialization in development mode');
    return;
  }

  // Validate DSN
  if (!config.dsn || !config.dsn.startsWith('https://')) {
    console.error('[Sentry] Invalid DSN provided, skipping initialization');
    return;
  }

  try {
    Sentry.init({
      dsn: config.dsn,
      environment: config.environment ?? import.meta.env.MODE,
      sampleRate: config.sampleRate ?? 1.0,
      tracesSampleRate: config.tracesSampleRate ?? 0.1,
      release: config.release ?? import.meta.env.VITE_APP_VERSION,
      integrations: [Sentry.browserTracingIntegration()],
      beforeSend: config.beforeSend,
      beforeBreadcrumb: config.beforeBreadcrumb,

      // Filter out sensitive data
      normalizeDepth: 5,

      // Ignore specific errors
      ignoreErrors: [
        // Random plugins/extensions
        'top.GLOBALS',
        // Original Facebook error
        'fb_xd_fragment',
        // Facebook bugged error
        'Script error.',
        // Safari specific error
        'QuotaExceededError',
        // Network errors that happen frequently
        'Non-Error promise rejection captured',
        // ResizeObserver errors (often non-critical)
        'ResizeObserver loop limit exceeded',
      ],

      // Ignore specific URLs from being traced
      tracePropagationTargets: [
        // Only trace requests to our API
        /^https:\/\/api\.gamelink\.com/,
        /^https:\/\/.*\.gamelink\.com/,
        // Include localhost for development
        'localhost',
        '127.0.0.1',
      ],
    });

    isInitialized = true;
    console.log('[Sentry] Initialized successfully');
  } catch (error) {
    console.error('[Sentry] Failed to initialize:', error);
  }
}

/**
 * Capture an exception and send it to Sentry
 *
 * @param error - The error to capture
 * @param context - Additional context information
 */
export function captureException(error: Error, context?: ErrorContext): void {
  if (!isInitialized) {
    console.warn('[Sentry] Not initialized, error not sent to Sentry');
    console.error('Error:', error);
    return;
  }

  Sentry.withScope((scope) => {
    if (context) {
      Object.entries(context).forEach(([key, value]) => {
        scope.setContext(key, value as Record<string, unknown>);
      });
    }
    Sentry.captureException(error);
  });
}

/**
 * Capture a message and send it to Sentry
 *
 * @param message - The message to capture
 * @param level - Severity level (fatal, error, warning, log, info, debug)
 * @param context - Additional context information
 */
export function captureMessage(
  message: string,
  level: 'fatal' | 'error' | 'warning' | 'log' | 'info' | 'debug' = 'info',
  context?: ErrorContext
): void {
  if (!isInitialized) {
    console.warn('[Sentry] Not initialized, message not sent to Sentry');
    console.log(`[${level.toUpperCase()}]`, message);
    return;
  }

  Sentry.withScope((scope) => {
    if (context) {
      Object.entries(context).forEach(([key, value]) => {
        scope.setContext(key, value as Record<string, unknown>);
      });
    }
    scope.setLevel(level);
    Sentry.captureMessage(message);
  });
}

/**
 * Add a breadcrumb to the current scope
 * Breadcrumbs are used to trace what happened before an error
 *
 * @param category - Category of the breadcrumb
 * @param message - Breadcrumb message
 * @param data - Additional data
 */
export function addBreadcrumb(
  category: string,
  message: string,
  data?: Record<string, unknown>
): void {
  if (!isInitialized) {
    return;
  }

  Sentry.addBreadcrumb({
    category,
    message,
    data,
    level: 'info',
  });
}

/**
 * Set a user for the current scope
 * Useful for tracking which user experienced an error
 *
 * @param id - User ID
 * @param email - User email (optional)
 * @param username - Username (optional)
 */
export function setUser(id: string, email?: string, username?: string): void {
  if (!isInitialized) {
    return;
  }

  Sentry.setUser({
    id,
    email,
    username,
  });
}

/**
 * Clear the currently set user
 */
export function clearUser(): void {
  if (!isInitialized) {
    return;
  }

  Sentry.setUser(null);
}

/**
 * Set a tag for the current scope
 * Tags can be used to filter and aggregate errors
 *
 * @param key - Tag key
 * @param value - Tag value
 */
export function setTag(key: string, value: string): void {
  if (!isInitialized) {
    return;
  }

  Sentry.setTag(key, value);
}

/**
 * Set extra context data for the current scope
 *
 * @param key - Context key
 * @param value - Context value
 */
export function setExtra(key: string, value: unknown): void {
  if (!isInitialized) {
    return;
  }

  Sentry.setExtra(key, value);
}

/**
 * Check if Sentry is initialized
 */
export function isSentryInitialized(): boolean {
  return isInitialized;
}

/**
 * Initialize Sentry from environment variables
 * Automatically reads configuration from import.meta.env
 */
export function initSentryFromEnv(): void {
  const enabled = import.meta.env.VITE_SENTRY_ENABLED === 'true';

  if (!enabled) {
    console.log('[Sentry] Disabled via environment variable');
    return;
  }

  const dsn = import.meta.env.VITE_SENTRY_DSN;

  if (!dsn) {
    console.warn('[Sentry] VITE_SENTRY_DSN not set, skipping initialization');
    return;
  }

  initSentry({
    dsn,
    environment: import.meta.env.VITE_SENTRY_ENVIRONMENT ?? 'production',
    sampleRate: parseFloat(import.meta.env.VITE_SENTRY_SAMPLE_RATE ?? '1.0'),
    tracesSampleRate: parseFloat(import.meta.env.VITE_SENTRY_TRACES_SAMPLE_RATE ?? '0.1'),
    release: import.meta.env.VITE_APP_VERSION,
    appName: import.meta.env.VITE_APP_TITLE,
  });
}

/**
 * Default export
 */
export default {
  initSentry,
  initSentryFromEnv,
  captureException,
  captureMessage,
  addBreadcrumb,
  setUser,
  clearUser,
  setTag,
  setExtra,
  isInitialized: isSentryInitialized,
};
