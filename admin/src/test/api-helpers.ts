/**
 * API Test Helpers
 *
 * Shared utilities for API testing including mock handlers,
 * response builders, and test utilities
 */

import { AxiosResponse } from 'axios';

/**
 * Standard API response structure
 */
export interface ApiSuccessResponse<T> {
  success: true;
  code: number;
  message: string;
  data: T;
}

export interface ApiErrorResponse {
  success: false;
  code: number;
  message: string;
  data?: unknown;
}

/**
 * Build a successful API response
 */
export function buildSuccessResponse<T>(data: T, message = 'Success'): ApiSuccessResponse<T> {
  return {
    success: true,
    code: 200,
    message,
    data,
  };
}

/**
 * Build an Axios response with API data
 */
export function buildAxiosResponse<T>(data: T, status = 200): AxiosResponse<ApiSuccessResponse<T>> {
  return {
    data: buildSuccessResponse(data),
    status,
    statusText: 'OK',
    headers: {},
    config: {} as any,
  };
}

/**
 * Build an error API response
 */
export function buildErrorResponse(code: number, message: string): ApiErrorResponse {
  return {
    success: false,
    code,
    message,
  };
}

/**
 * Mock authentication tokens for testing
 */
export const mockTokens = {
  valid: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.valid',
  expired: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.expired',
  refreshed: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.refreshed',
};

/**
 * Mock user data
 */
export const mockUsers = {
  admin: {
    id: 1,
    username: 'admin',
    email: 'admin@gamelink.com',
    role: 'admin',
  },
  player: {
    id: 2,
    username: 'player1',
    email: 'player@test.com',
    role: 'player',
  },
  user: {
    id: 3,
    username: 'user1',
    email: 'user@test.com',
    role: 'user',
  },
};

/**
 * Mock login responses
 */
export const mockLoginResponses = {
  success: {
    token: mockTokens.valid,
    user: mockUsers.admin,
  },
  invalidCredentials: {
    success: false,
    code: 401,
    message: 'Invalid username or password',
  },
};

/**
 * Wait for async operations to complete
 */
export function waitForAsync(timeout = 0): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, timeout));
}

/**
 * Clear all storage (localStorage and sessionStorage)
 */
export function clearAllStorage(): void {
  localStorage.clear();
  sessionStorage.clear();
}

/**
 * Set a mock token in both storage locations
 */
export function setMockToken(token = mockTokens.valid): void {
  localStorage.setItem('token', token);
  sessionStorage.setItem('auth_token', token);
}

/**
 * Create a mock rejection error
 */
export function createAxiosError(
  status: number,
  message: string,
  code?: string
): any {
  const error: any = new Error(message);
  error.response = {
    status,
    data: {
      success: false,
      code: status,
      message,
    },
  };
  error.code = code || 'ERR_BAD_REQUEST';
  return error;
}

/**
 * Mock window.location for redirect testing
 */
export function mockWindowLocation(): { href: string } {
  const location = { href: '' };
  Object.defineProperty(window, 'location', {
    writable: true,
    value: location,
  });
  return location;
}

/**
 * Test helper to verify interceptor behavior
 */
export function expectAuthHeader(token: string): Record<string, string> {
  return {
    Authorization: `Bearer ${token}`,
  };
}

/**
 * Common test scenarios for API endpoints
 */
export const apiTestScenarios = {
  unauthorized: {
    status: 401,
    message: 'Unauthorized',
  },
  forbidden: {
    status: 403,
    message: 'Forbidden',
  },
  notFound: {
    status: 404,
    message: 'Resource not found',
  },
  serverError: {
    status: 500,
    message: 'Internal server error',
  },
  validationError: {
    status: 400,
    message: 'Validation error',
  },
};

/**
 * Mock pagination response
 */
export function buildPaginatedResponse<T>(
  items: T[],
  page = 1,
  pageSize = 10,
  total?: number
): ApiSuccessResponse<{
  items: T[];
  totalCount: number;
  page: number;
  pageSize: number;
}> {
  return buildSuccessResponse({
    items,
    totalCount: total ?? items.length,
    page,
    pageSize,
  });
}

/**
 * Generate test data with unique IDs
 */
export function generateTestId(): string {
  return `test_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
}

/**
 * Verify if an object is a valid API response
 */
export function isValidApiResponse(response: unknown): response is ApiSuccessResponse<unknown> {
  return (
    typeof response === 'object' &&
    response !== null &&
    'success' in response &&
    'code' in response &&
    'data' in response
  );
}

/**
 * Test timeout - use for testing timeout scenarios
 */
export const TEST_TIMEOUT = 10000;

/**
 * Retry helper for flaky tests
 */
export async function retry<T>(
  fn: () => Promise<T>,
  maxRetries = 3,
  delay = 100
): Promise<T> {
  let lastError: Error | undefined;

  for (let i = 0; i < maxRetries; i++) {
    try {
      return await fn();
    } catch (error) {
      lastError = error as Error;
      if (i < maxRetries - 1) {
        await new Promise(resolve => setTimeout(resolve, delay));
      }
    }
  }

  throw lastError;
}
