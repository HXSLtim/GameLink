/**
 * Unified Logger Utility
 *
 * Provides consistent logging across the application with environment-aware behavior.
 * In production, logs are suppressed to avoid exposing sensitive information.
 * In development, logs include stack traces for debugging.
 *
 * @example
 * ```ts
 * import { logger } from '@/utils/logger';
 *
 * logger.info('User logged in', { userId: 123 });
 * logger.warn('API rate limit approaching', { remaining: 5 });
 * logger.error('Failed to fetch data', error);
 * ```
 */

type LogLevel = 'info' | 'warn' | 'error' | 'debug';

interface LogContext {
  [key: string]: unknown;
}

// Allow LogContext to be any object type for flexibility
type LogContextInput = LogContext | Record<string, unknown> | unknown[] | unknown;

interface LogEntry {
  level: LogLevel;
  message: string;
  context?: LogContext;
  timestamp: string;
  stack?: string;
}

/**
 * Logger class with environment-aware logging
 */
class Logger {
  private isDevelopment: boolean;
  private isProduction: boolean;

  constructor() {
    this.isDevelopment = import.meta.env.MODE === 'development';
    this.isProduction = import.meta.env.MODE === 'production';
  }

  /**
   * Format log entry with timestamp and context
   */
  private formatEntry(level: LogLevel, message: string, context?: LogContextInput): LogEntry {
    return {
      level,
      message,
      context: context as LogContext | undefined,
      timestamp: new Date().toISOString(),
    };
  }

  /**
   * Get stack trace for error logging (development only)
   */
  private getStackTrace(): string | undefined {
    if (!this.isDevelopment) {
      return undefined;
    }

    const stack = new Error().stack;
    // Remove the first few lines (getStackTrace, formatEntry, log method)
    if (stack) {
      const lines = stack.split('\n').slice(4);
      return lines.join('\n');
    }
    return undefined;
  }

  /**
   * Core logging method
   */
  private log(level: LogLevel, message: string, context?: LogContextInput, error?: Error): void {
    // In production, suppress all console output
    if (this.isProduction) {
      return;
    }

    const entry = this.formatEntry(level, message, context);

    // Add stack trace for errors in development
    if (error && this.isDevelopment) {
      entry.stack = this.getStackTrace();
    }

    // Console output with appropriate styling
    const timestamp = entry.timestamp;
    const prefix = `[${timestamp}] [${level.toUpperCase()}]`;

    switch (level) {
      case 'info':
        console.log(prefix, message, context || '');
        break;
      case 'warn':
        console.warn(prefix, message, context || '');
        break;
      case 'error':
        console.error(prefix, message, context || '', error || '');
        if (entry.stack) {
          console.error('Stack trace:', entry.stack);
        }
        break;
      case 'debug':
        if (this.isDevelopment) {
          console.log(prefix, message, context || '');
        }
        break;
    }
  }

  /**
   * Log informational message
   * Use for general application flow and successful operations
   */
  public info(message: string, context?: LogContextInput): void {
    this.log('info', message, context);
  }

  /**
   * Log warning message
   * Use for potentially problematic situations that don't prevent execution
   */
  public warn(message: string, context?: LogContextInput): void {
    this.log('warn', message, context);
  }

  /**
   * Log error message
   * Use for errors and exceptions that affect functionality
   *
   * Overloads for different call patterns:
   * - logger.error(message)
   * - logger.error(message, error)
   * - logger.error(message, context)
   * - logger.error(message, error, context)
   */
  public error(message: string, errorOrContext?: Error | unknown | LogContextInput, context?: LogContextInput): void {
    // Determine if second param is error or context
    let errorObj: Error | undefined;
    let contextObj: LogContextInput | undefined;

    if (errorOrContext) {
      if (errorOrContext instanceof Error || (errorOrContext as Error).message) {
        // It's an Error object
        errorObj = errorOrContext as Error;
        contextObj = context;
      } else {
        // It's a context object
        contextObj = errorOrContext as LogContextInput;
      }
    }

    this.log('error', message, contextObj, errorObj);
  }

  /**
   * Log debug message
   * Only outputs in development mode
   */
  public debug(message: string, context?: LogContextInput): void {
    this.log('debug', message, context);
  }

  /**
   * Log API request
   * Use for tracking API calls
   */
  public api(method: string, url: string, context?: LogContextInput): void {
    this.debug(`API ${method}`, { url, ...context });
  }

  /**
   * Log API response
   * Use for tracking API responses
   */
  public apiResponse(method: string, url: string, status: number, context?: LogContextInput): void {
    this.debug(`API Response ${method}`, { url, status, ...context });
  }

  /**
   * Log user action
   * Use for tracking important user interactions
   */
  public userAction(action: string, context?: LogContextInput): void {
    this.info(`User action: ${action}`, context);
  }

  /**
   * Log component lifecycle
   * Use for debugging component mounts/unmounts
   */
  public lifecycle(component: string, phase: 'mount' | 'unmount' | 'update', context?: LogContextInput): void {
    this.debug(`Component ${phase}`, { component, ...context });
  }
}

// Export singleton instance
export const logger = new Logger();

// Export types for context
export type { LogContext, LogContextInput };
