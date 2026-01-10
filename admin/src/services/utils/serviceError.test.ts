/**
 * Property-Based Tests for Service Error Utilities
 *
 * **Feature: admin-phase3-improvements, Property 1: Service Error Format Consistency**
 * **Validates: Requirements 1.3**
 *
 * Tests that service errors maintain consistent format across all error scenarios.
 */

import { describe, it, expect } from 'vitest';
import * as fc from 'fast-check';
import {
  ServiceException,
  ServiceErrorCodes,
  type ServiceError,
} from './serviceError';

/**
 * Validates that a ServiceError object has the required structure
 */
function isValidServiceError(error: ServiceError): boolean {
  // Must have code (non-empty string)
  if (typeof error.code !== 'string' || error.code.length === 0) {
    return false;
  }

  // Must have message (non-empty string)
  if (typeof error.message !== 'string' || error.message.length === 0) {
    return false;
  }

  // Details must be undefined or an object
  if (error.details !== undefined && typeof error.details !== 'object') {
    return false;
  }

  // originalError must be undefined or an Error instance
  if (error.originalError !== undefined && !(error.originalError instanceof Error)) {
    return false;
  }

  return true;
}

describe('ServiceException - Property Tests', () => {
  /**
   * **Feature: admin-phase3-improvements, Property 1: Service Error Format Consistency**
   * **Validates: Requirements 1.3**
   *
   * For any service method that throws an error, the error object SHALL contain
   * code, message, and optionally details fields in a consistent format.
   */
  it('Property 1: toError() always returns valid ServiceError format', () => {
    // Arbitrary for error codes
    const errorCodeArb = fc.constantFrom(...Object.values(ServiceErrorCodes));

    // Arbitrary for non-empty messages
    const messageArb = fc.string({ minLength: 1, maxLength: 500 });

    // Arbitrary for optional details
    const detailsArb = fc.option(
      fc.dictionary(fc.string({ minLength: 1, maxLength: 50 }), fc.jsonValue()),
      { nil: undefined }
    );

    fc.assert(
      fc.property(errorCodeArb, messageArb, detailsArb, (code, message, details) => {
        const exception = new ServiceException(
          code,
          message,
          details as Record<string, unknown> | undefined
        );
        const error = exception.toError();

        // Verify the error has valid format
        expect(isValidServiceError(error)).toBe(true);

        // Verify specific fields match
        expect(error.code).toBe(code);
        expect(error.message).toBe(message);
        expect(error.details).toEqual(details);
      }),
      { numRuns: 100 }
    );
  });

  it('Property 1: fromError() always produces valid ServiceError format', () => {
    // Arbitrary for various error types
    const errorArb = fc.oneof(
      // String errors
      fc.string({ minLength: 1 }).map((s) => s),
      // Error objects
      fc.string({ minLength: 1 }).map((msg) => new Error(msg)),
      // ServiceException
      fc.tuple(
        fc.constantFrom(...Object.values(ServiceErrorCodes)),
        fc.string({ minLength: 1 })
      ).map(([code, msg]) => new ServiceException(code, msg)),
      // Null/undefined
      fc.constant(null),
      fc.constant(undefined)
    );

    const contextArb = fc.option(fc.string({ minLength: 1, maxLength: 100 }), { nil: undefined });

    fc.assert(
      fc.property(errorArb, contextArb, (error, context) => {
        const exception = ServiceException.fromError(error, context);
        const serviceError = exception.toError();

        // Must always produce valid format
        expect(isValidServiceError(serviceError)).toBe(true);

        // Code must be a valid error code
        expect(typeof serviceError.code).toBe('string');
        expect(serviceError.code.length).toBeGreaterThan(0);

        // Message must be non-empty
        expect(typeof serviceError.message).toBe('string');
        expect(serviceError.message.length).toBeGreaterThan(0);
      }),
      { numRuns: 100 }
    );
  });

  it('Property 1: validation() helper produces valid ServiceError format', () => {
    const messageArb = fc.string({ minLength: 1, maxLength: 200 });
    const fieldErrorsArb = fc.array(
      fc.record({
        field: fc.string({ minLength: 1, maxLength: 50 }),
        message: fc.string({ minLength: 1, maxLength: 200 }),
      }),
      { minLength: 1, maxLength: 10 }
    );

    fc.assert(
      fc.property(messageArb, fieldErrorsArb, (message, errors) => {
        const exception = ServiceException.validation(message, errors);
        const serviceError = exception.toError();

        // Must have valid format
        expect(isValidServiceError(serviceError)).toBe(true);

        // Must have VALIDATION_ERROR code
        expect(serviceError.code).toBe(ServiceErrorCodes.VALIDATION_ERROR);

        // Must include errors in details
        expect(serviceError.details).toBeDefined();
        expect(serviceError.details?.errors).toEqual(errors);
      }),
      { numRuns: 100 }
    );
  });

  it('Property 1: notFound() helper produces valid ServiceError format', () => {
    const resourceArb = fc.string({ minLength: 1, maxLength: 50 });
    const idArb = fc.option(
      fc.oneof(fc.integer({ min: 1 }), fc.string({ minLength: 1, maxLength: 50 })),
      { nil: undefined }
    );

    fc.assert(
      fc.property(resourceArb, idArb, (resource, id) => {
        const exception = ServiceException.notFound(resource, id);
        const serviceError = exception.toError();

        // Must have valid format
        expect(isValidServiceError(serviceError)).toBe(true);

        // Must have NOT_FOUND code
        expect(serviceError.code).toBe(ServiceErrorCodes.NOT_FOUND);

        // Must include resource in details
        expect(serviceError.details).toBeDefined();
        expect(serviceError.details?.resource).toBe(resource);
      }),
      { numRuns: 100 }
    );
  });
});

describe('ServiceException - Unit Tests', () => {
  it('should create exception with all properties', () => {
    const originalError = new Error('Original');
    const exception = new ServiceException(
      ServiceErrorCodes.VALIDATION_ERROR,
      'Test message',
      { field: 'email' },
      originalError
    );

    expect(exception.code).toBe(ServiceErrorCodes.VALIDATION_ERROR);
    expect(exception.message).toBe('Test message');
    expect(exception.details).toEqual({ field: 'email' });
    expect(exception.originalError).toBe(originalError);
    expect(exception.name).toBe('ServiceException');
  });

  it('should be instanceof Error', () => {
    const exception = new ServiceException(ServiceErrorCodes.UNKNOWN_ERROR, 'Test');
    expect(exception).toBeInstanceOf(Error);
    expect(exception).toBeInstanceOf(ServiceException);
  });

  it('fromError should preserve ServiceException', () => {
    const original = new ServiceException(ServiceErrorCodes.NOT_FOUND, 'Not found');
    const result = ServiceException.fromError(original);
    expect(result).toBe(original);
  });

  it('fromError should wrap regular Error', () => {
    const original = new Error('Regular error');
    const result = ServiceException.fromError(original, 'Context');

    expect(result.code).toBe(ServiceErrorCodes.UNKNOWN_ERROR);
    expect(result.message).toBe('Context: Regular error');
    expect(result.originalError).toBe(original);
  });
});
