import { Message } from '@arco-design/web-react';

/**
 * Error severity levels
 */
export enum ErrorSeverity {
  INFO = 'info',
  WARNING = 'warning',
  ERROR = 'error',
  CRITICAL = 'critical',
}

/**
 * Application error class
 */
export class AppError extends Error {
  public readonly code: string;
  public readonly severity: ErrorSeverity;
  public readonly timestamp: Date;
  public readonly context?: Record<string, unknown>;

  constructor(
    message: string,
    code: string = 'UNKNOWN_ERROR',
    severity: ErrorSeverity = ErrorSeverity.ERROR,
    context?: Record<string, unknown>,
  ) {
    super(message);
    this.name = 'AppError';
    this.code = code;
    this.severity = severity;
    this.timestamp = new Date();
    this.context = context;

    // Maintains proper stack trace for where our error was thrown (only available on V8)
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, AppError);
    }
  }
}

/**
 * Error logger interface
 */
interface ErrorLogger {
  log: (error: Error | AppError, context?: Record<string, unknown>) => void;
}

/**
 * Console error logger
 */
class ConsoleErrorLogger implements ErrorLogger {
  log(error: Error | AppError, context?: Record<string, unknown>): void {
    const isAppError = error instanceof AppError;

    console.group(`🔴 Error: ${error.message}`);
    console.error('Type:', error.name);
    console.error('Message:', error.message);

    if (isAppError) {
      console.error('Code:', error.code);
      console.error('Severity:', error.severity);
      console.error('Timestamp:', error.timestamp.toISOString());
      if (error.context) {
        console.error('Context:', error.context);
      }
    }

    if (context) {
      console.error('Additional Context:', context);
    }

    if (error.stack) {
      console.error('Stack:', error.stack);
    }

    console.groupEnd();
  }
}

/**
 * Error handler class
 */
class ErrorHandler {
  private logger: ErrorLogger;
  private isDevelopment: boolean;

  constructor() {
    this.logger = new ConsoleErrorLogger();
    this.isDevelopment = import.meta.env.DEV;
  }

  /**
   * Set custom error logger
   */
  setLogger(logger: ErrorLogger): void {
    this.logger = logger;
  }

  /**
   * Handle error with user notification
   */
  handle(error: Error | AppError, showToUser: boolean = true): void {
    // Log error
    this.logger.log(error);

    // Show user notification
    if (showToUser) {
      this.notifyUser(error);
    }

    // Report to external service in production
    if (!this.isDevelopment) {
      this.reportToService(error);
    }
  }

  /**
   * Handle async errors
   */
  async handleAsync<T>(
    promise: Promise<T>,
    errorMessage?: string,
  ): Promise<[T | null, Error | null]> {
    try {
      const data = await promise;
      return [data, null];
    } catch (error) {
      const err = this.normalizeError(error, errorMessage);
      this.handle(err);
      return [null, err];
    }
  }

  /**
   * Normalize error to Error instance
   */
  private normalizeError(error: unknown, defaultMessage?: string): Error {
    if (error instanceof AppError) {
      return error;
    }

    if (error instanceof Error) {
      return error;
    }

    if (typeof error === 'string') {
      return new Error(error);
    }

    return new Error(defaultMessage || 'An unknown error occurred');
  }

  /**
   * Notify user about error
   */
  private notifyUser(error: Error | AppError): void {
    let message = error.message;
    let duration = 3000;

    if (error instanceof AppError) {
      switch (error.severity) {
        case ErrorSeverity.INFO:
          Message.info(message);
          return;
        case ErrorSeverity.WARNING:
          Message.warning(message);
          duration = 4000;
          break;
        case ErrorSeverity.CRITICAL:
          duration = 5000;
          message = `严重错误: ${message}`;
          break;
      }
    }

    Message.error({
      content: message,
      duration,
    });
  }

  /**
   * Report error to external service
   */
  private reportToService(error: Error | AppError): void {
    // TODO: Implement error reporting service (e.g., Sentry)
    // In production, you would send errors to monitoring service
    console.log('📊 Error reported to monitoring service:', error.message);
  }
}

// Singleton instance
export const errorHandler = new ErrorHandler();

/**
 * Helper function to handle API errors
 */
export function handleApiError(error: unknown, operation: string): void {
  const message = error instanceof Error ? error.message : '操作失败';
  const appError = new AppError(`${operation}失败: ${message}`, 'API_ERROR', ErrorSeverity.ERROR, {
    operation,
    originalError: error,
  });
  errorHandler.handle(appError);
}

/**
 * Helper function to create user-friendly error messages
 */
export function getUserFriendlyMessage(error: unknown): string {
  if (error instanceof AppError) {
    return error.message;
  }

  if (error instanceof Error) {
    return error.message;
  }

  return '操作失败，请重试';
}

