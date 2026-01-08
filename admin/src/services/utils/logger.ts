/**
 * Service Logger Module
 *
 * Provides logging infrastructure for service layer operations with:
 * - Parameter sanitization for sensitive data
 * - Optional external error tracking integration
 * - Structured logging with context
 *
 * @module services/utils/logger
 */

/**
 * Sensitive keys that should be redacted in logs
 */
const SENSITIVE_KEYS = [
  'password',
  'token',
  'secret',
  'key',
  'authorization',
  'credential',
  'apikey',
  'api_key',
  'accesstoken',
  'access_token',
  'refreshtoken',
  'refresh_token',
  'privatekey',
  'private_key',
];

/**
 * Log level constants
 */
export const LogLevel = {
  DEBUG: 0,
  INFO: 1,
  WARN: 2,
  ERROR: 3,
} as const;

export type LogLevel = (typeof LogLevel)[keyof typeof LogLevel];

/**
 * Log entry structure for structured logging
 */
export interface LogEntry {
  level: LogLevel;
  message: string;
  serviceName: string;
  timestamp: Date;
  context?: Record<string, unknown>;
  error?: Error;
  duration?: number;
}

/**
 * External error tracker interface for integration with services like Sentry
 */
export interface ExternalErrorTracker {
  captureException(error: Error, context?: Record<string, unknown>): void;
  captureMessage(message: string, level: 'info' | 'warning' | 'error'): void;
}

/**
 * Logger configuration options
 */
export interface LoggerConfig {
  /** Minimum log level to output */
  minLevel?: LogLevel;
  /** Whether to include timestamps in console output */
  includeTimestamp?: boolean;
  /** External error tracker for production error reporting */
  errorTracker?: ExternalErrorTracker;
  /** Custom log handler for testing or custom output */
  customHandler?: (entry: LogEntry) => void;
}

/**
 * Service Logger Interface
 *
 * Provides structured logging methods for service operations.
 * All methods support optional context for additional metadata.
 */
export interface ServiceLogger {
  /**
   * Log debug message (development only)
   * @param message - Log message
   * @param context - Optional context data (will be sanitized)
   */
  debug(message: string, context?: Record<string, unknown>): void;

  /**
   * Log info message
   * @param message - Log message
   * @param context - Optional context data (will be sanitized)
   */
  info(message: string, context?: Record<string, unknown>): void;

  /**
   * Log warning message
   * @param message - Log message
   * @param context - Optional context data (will be sanitized)
   */
  warn(message: string, context?: Record<string, unknown>): void;

  /**
   * Log error message with optional error object
   * @param message - Log message
   * @param error - Optional error object
   * @param context - Optional context data (will be sanitized)
   */
  error(message: string, error?: Error, context?: Record<string, unknown>): void;

  /**
   * Get the service name this logger is associated with
   */
  getServiceName(): string;

  /**
   * Get all log entries (for testing/debugging)
   */
  getLogEntries(): LogEntry[];

  /**
   * Clear all stored log entries
   */
  clearLogEntries(): void;
}

/**
 * Sanitize parameters by redacting sensitive values
 *
 * @param params - Parameters to sanitize
 * @returns Sanitized parameters with sensitive values replaced by '[REDACTED]'
 */
export function sanitizeParams(params: Record<string, unknown>): Record<string, unknown> {
  const sanitized: Record<string, unknown> = {};

  for (const [key, value] of Object.entries(params)) {
    const lowerKey = key.toLowerCase();

    // Check if key contains any sensitive keyword
    if (SENSITIVE_KEYS.some((sk) => lowerKey.includes(sk))) {
      sanitized[key] = '[REDACTED]';
    } else if (typeof value === 'object' && value !== null) {
      if (Array.isArray(value)) {
        // Sanitize array elements if they are objects
        sanitized[key] = value.map((item) =>
          typeof item === 'object' && item !== null
            ? sanitizeParams(item as Record<string, unknown>)
            : item
        );
      } else {
        // Recursively sanitize nested objects
        sanitized[key] = sanitizeParams(value as Record<string, unknown>);
      }
    } else {
      sanitized[key] = value;
    }
  }

  return sanitized;
}

/**
 * Check if a key is sensitive and should be redacted
 *
 * @param key - The key to check
 * @returns true if the key should be redacted
 */
export function isSensitiveKey(key: string): boolean {
  const lowerKey = key.toLowerCase();
  return SENSITIVE_KEYS.some((sk) => lowerKey.includes(sk));
}

/**
 * Default Service Logger Implementation
 *
 * Provides console-based logging with:
 * - Automatic parameter sanitization
 * - Development-only debug logging
 * - Optional external error tracking integration
 * - Log entry storage for testing
 */
export class DefaultServiceLogger implements ServiceLogger {
  private serviceName: string;
  private isDev: boolean;
  private config: Required<Omit<LoggerConfig, 'errorTracker' | 'customHandler'>> & {
    errorTracker?: ExternalErrorTracker;
    customHandler?: (entry: LogEntry) => void;
  };
  private logEntries: LogEntry[] = [];
  private readonly MAX_LOG_ENTRIES = 1000;

  constructor(serviceName: string, config: LoggerConfig = {}) {
    this.serviceName = serviceName;
    this.isDev = typeof import.meta !== 'undefined' ? (import.meta.env?.DEV ?? false) : false;
    this.config = {
      minLevel: config.minLevel ?? (this.isDev ? LogLevel.DEBUG : LogLevel.INFO),
      includeTimestamp: config.includeTimestamp ?? true,
      errorTracker: config.errorTracker,
      customHandler: config.customHandler,
    };
  }

  debug(message: string, context?: Record<string, unknown>): void {
    if (this.config.minLevel > LogLevel.DEBUG) return;

    const entry = this.createLogEntry(LogLevel.DEBUG, message, context);
    this.storeEntry(entry);

    if (this.config.customHandler) {
      this.config.customHandler(entry);
      return;
    }

    if (this.isDev) {
      const sanitizedContext = context ? sanitizeParams(context) : undefined;
      console.debug(this.formatMessage(message), sanitizedContext ?? '');
    }
  }

  info(message: string, context?: Record<string, unknown>): void {
    if (this.config.minLevel > LogLevel.INFO) return;

    const entry = this.createLogEntry(LogLevel.INFO, message, context);
    this.storeEntry(entry);

    if (this.config.customHandler) {
      this.config.customHandler(entry);
      return;
    }

    const sanitizedContext = context ? sanitizeParams(context) : undefined;
    console.info(this.formatMessage(message), sanitizedContext ?? '');
  }

  warn(message: string, context?: Record<string, unknown>): void {
    if (this.config.minLevel > LogLevel.WARN) return;

    const entry = this.createLogEntry(LogLevel.WARN, message, context);
    this.storeEntry(entry);

    if (this.config.customHandler) {
      this.config.customHandler(entry);
      return;
    }

    const sanitizedContext = context ? sanitizeParams(context) : undefined;
    console.warn(this.formatMessage(message), sanitizedContext ?? '');

    // Send warnings to external tracker if configured
    if (this.config.errorTracker) {
      this.config.errorTracker.captureMessage(
        `[${this.serviceName}] ${message}`,
        'warning'
      );
    }
  }

  error(message: string, error?: Error, context?: Record<string, unknown>): void {
    const entry = this.createLogEntry(LogLevel.ERROR, message, context, error);
    this.storeEntry(entry);

    if (this.config.customHandler) {
      this.config.customHandler(entry);
      return;
    }

    const sanitizedContext = context ? sanitizeParams(context) : undefined;
    console.error(this.formatMessage(message), { error, ...sanitizedContext });

    // Send errors to external tracker if configured
    if (this.config.errorTracker && error) {
      this.config.errorTracker.captureException(error, {
        serviceName: this.serviceName,
        message,
        ...sanitizedContext,
      });
    }
  }

  getServiceName(): string {
    return this.serviceName;
  }

  getLogEntries(): LogEntry[] {
    return [...this.logEntries];
  }

  clearLogEntries(): void {
    this.logEntries = [];
  }

  private formatMessage(message: string): string {
    const prefix = `[${this.serviceName}]`;
    if (this.config.includeTimestamp) {
      const timestamp = new Date().toISOString();
      return `${timestamp} ${prefix} ${message}`;
    }
    return `${prefix} ${message}`;
  }

  private createLogEntry(
    level: LogLevel,
    message: string,
    context?: Record<string, unknown>,
    error?: Error
  ): LogEntry {
    return {
      level,
      message,
      serviceName: this.serviceName,
      timestamp: new Date(),
      context: context ? sanitizeParams(context) : undefined,
      error,
    };
  }

  private storeEntry(entry: LogEntry): void {
    this.logEntries.push(entry);
    // Keep only last MAX_LOG_ENTRIES entries
    if (this.logEntries.length > this.MAX_LOG_ENTRIES) {
      this.logEntries = this.logEntries.slice(-this.MAX_LOG_ENTRIES);
    }
  }
}

/**
 * Create a logger instance for a service
 *
 * @param serviceName - Name of the service
 * @param config - Optional logger configuration
 * @returns ServiceLogger instance
 */
export function createServiceLogger(
  serviceName: string,
  config?: LoggerConfig
): ServiceLogger {
  return new DefaultServiceLogger(serviceName, config);
}
