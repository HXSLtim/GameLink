/**
 * Service Layer Error Utilities
 * Provides standardized error handling for the service layer
 *
 * @module services/utils/serviceError
 */

/**
 * Service error codes for consistent error identification
 */
export const ServiceErrorCodes = {
  // General
  UNKNOWN_ERROR: 'UNKNOWN_ERROR',
  VALIDATION_ERROR: 'VALIDATION_ERROR',
  NOT_FOUND: 'NOT_FOUND',
  UNAUTHORIZED: 'UNAUTHORIZED',
  FORBIDDEN: 'FORBIDDEN',
  NETWORK_ERROR: 'NETWORK_ERROR',
  TIMEOUT_ERROR: 'TIMEOUT_ERROR',

  // User
  USER_EMAIL_EXISTS: 'USER_EMAIL_EXISTS',
  USER_PHONE_EXISTS: 'USER_PHONE_EXISTS',
  USER_INVALID_STATUS: 'USER_INVALID_STATUS',
  USER_INVALID_EMAIL: 'USER_INVALID_EMAIL',
  USER_INVALID_PHONE: 'USER_INVALID_PHONE',
  USER_WEAK_PASSWORD: 'USER_WEAK_PASSWORD',

  // Order
  ORDER_CANNOT_CANCEL: 'ORDER_CANNOT_CANCEL',
  ORDER_CANNOT_REFUND: 'ORDER_CANNOT_REFUND',
  ORDER_INVALID_REFUND_AMOUNT: 'ORDER_INVALID_REFUND_AMOUNT',
  ORDER_NOT_FOUND: 'ORDER_NOT_FOUND',

  // Player
  PLAYER_INVALID_VERIFICATION: 'PLAYER_INVALID_VERIFICATION',
  PLAYER_ALREADY_EXISTS: 'PLAYER_ALREADY_EXISTS',
  PLAYER_NOT_FOUND: 'PLAYER_NOT_FOUND',

  // Import
  IMPORT_INVALID_FILE: 'IMPORT_INVALID_FILE',
  IMPORT_FILE_TOO_LARGE: 'IMPORT_FILE_TOO_LARGE',
  IMPORT_INVALID_STRUCTURE: 'IMPORT_INVALID_STRUCTURE',
  IMPORT_VALIDATION_FAILED: 'IMPORT_VALIDATION_FAILED',
  IMPORT_DUPLICATE_FOUND: 'IMPORT_DUPLICATE_FOUND',

  // Batch Operations
  BATCH_PARTIAL_FAILURE: 'BATCH_PARTIAL_FAILURE',
  BATCH_OPERATION_IN_PROGRESS: 'BATCH_OPERATION_IN_PROGRESS',
} as const;

export type ServiceErrorCode = (typeof ServiceErrorCodes)[keyof typeof ServiceErrorCodes];

/**
 * Service error interface for consistent error structure
 */
export interface ServiceError {
  /** Error code for programmatic handling */
  code: string;
  /** Human-readable error message */
  message: string;
  /** Optional additional details about the error */
  details?: Record<string, unknown>;
  /** Original error if this wraps another error */
  originalError?: Error;
}

/**
 * ServiceException class for throwing typed service errors
 *
 * @example
 * ```typescript
 * throw new ServiceException(
 *   ServiceErrorCodes.VALIDATION_ERROR,
 *   'Invalid email format',
 *   { field: 'email', value: 'invalid' }
 * );
 * ```
 */
export class ServiceException extends Error {
  public readonly code: string;
  public readonly details?: Record<string, unknown>;
  public readonly originalError?: Error;

  constructor(
    code: string,
    message: string,
    details?: Record<string, unknown>,
    originalError?: Error
  ) {
    super(message);
    this.name = 'ServiceException';
    this.code = code;
    this.details = details;
    this.originalError = originalError;

    // Maintains proper stack trace for where error was thrown (V8 engines)
    const ErrorWithCapture = Error as typeof Error & {
      captureStackTrace?: (target: object, constructor: unknown) => void;
    };
    if (typeof ErrorWithCapture.captureStackTrace === 'function') {
      ErrorWithCapture.captureStackTrace(this, ServiceException);
    }
  }

  /**
   * Convert to a plain ServiceError object
   */
  toError(): ServiceError {
    return {
      code: this.code,
      message: this.message,
      details: this.details,
      originalError: this.originalError,
    };
  }

  /**
   * Create a ServiceException from an unknown error
   */
  static fromError(error: unknown, context?: string): ServiceException {
    if (error instanceof ServiceException) {
      return error;
    }

    const message = error instanceof Error ? error.message : String(error);
    const contextMessage = context ? `${context}: ${message}` : message;

    return new ServiceException(
      ServiceErrorCodes.UNKNOWN_ERROR,
      contextMessage,
      undefined,
      error instanceof Error ? error : undefined
    );
  }

  /**
   * Create a validation error with field details
   */
  static validation(
    message: string,
    errors: Array<{ field: string; message: string }>
  ): ServiceException {
    return new ServiceException(ServiceErrorCodes.VALIDATION_ERROR, message, { errors });
  }

  /**
   * Create a not found error
   */
  static notFound(resource: string, id?: number | string): ServiceException {
    const message = id ? `${resource} with ID ${id} not found` : `${resource} not found`;
    return new ServiceException(ServiceErrorCodes.NOT_FOUND, message, { resource, id });
  }
}
