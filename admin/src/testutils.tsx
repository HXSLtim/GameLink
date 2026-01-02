/**
 * Testing Utilities for GameLink Admin Panel
 *
 * Provides common testing helpers for React components including:
 * - Provider wrappers for tests
 * - Router mocking
 * - Loading state helpers
 * - Test helpers and utilities
 */

import React, { ReactElement } from 'react';
import { render, RenderOptions } from '@testing-library/react';
import { BrowserRouter, MemoryRouter, Route, Routes } from 'react-router-dom';
import { ConfigProvider, theme } from 'antd';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

/**
 * Mock user data for testing
 */
export const mockUser = {
  id: 1,
  name: 'Test User',
  email: 'test@example.com',
  phone: '13800138000',
  role: 'admin' as const,
  status: 'active' as const,
  createdAt: '2024-01-01T00:00:00Z',
  avatarUrl: 'https://example.com/avatar.jpg',
};

/**
 * NOTE: API mocks should be defined in individual test files
 * to avoid hoisting issues. Use this pattern:
 *
 * const mockApi = {
 *   getUsers: vi.fn(),
 *   updateUser: vi.fn(),
 * };
 *
 * vi.mock('@/api/admin', () => ({
 *   adminApi: mockApi,
 * }));
 */

/**
 * Mock router navigate function
 */
export const mockNavigate = vi.fn();

/**
 * Create a test query client
 */
export function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        gcTime: 0,
      },
      mutations: {
        retry: false,
      },
    },
    logger: {
      log: console.log,
      warn: console.warn,
      error: () => {}, // Suppress error logs in tests
    },
  });
}

/**
 * Custom render function with providers
 */
interface CustomRenderOptions extends Omit<RenderOptions, 'wrapper'> {
  queryClient?: QueryClient;
  route?: string;
  router?: 'memory' | 'browser';
}

export function renderWithProviders(
  ui: ReactElement,
  {
    queryClient = createTestQueryClient(),
    route = '/',
    router = 'memory',
    ...renderOptions
  }: CustomRenderOptions = {}
) {
  const Wrapper = ({ children }: { children: React.ReactNode }) => {
    const routerElement =
      router === 'memory' ? (
        <MemoryRouter initialEntries={[route]}>
          <Routes>
            <Route path="*" element={children} />
          </Routes>
        </MemoryRouter>
      ) : (
        <BrowserRouter>{children}</BrowserRouter>
      );

    return (
      <ConfigProvider
        theme={{
          algorithm: theme.defaultAlgorithm,
          token: {
            colorPrimary: '#1890ff',
            colorSuccess: '#52c41a',
            colorWarning: '#faad14',
            colorError: '#ff4d4f',
          },
        }}
      >
        <QueryClientProvider client={queryClient}>
          {routerElement}
        </QueryClientProvider>
      </ConfigProvider>
    );
  };

  return {
    ...render(ui, { wrapper: Wrapper, ...renderOptions }),
  };
}

/**
 * Mock router hook
 */
export function mockRouter() {
  return {
    navigate: mockNavigate,
    location: { pathname: '/', search: '', hash: '', state: null, key: 'default' },
    params: {},
  };
}

/**
 * Wait for loading to complete (useful for async operations)
 */
export async function waitForLoading(
  callback: () => void,
  timeout = 1000
): Promise<void> {
  return new Promise((resolve) => {
    const timer = setTimeout(() => {
      callback();
      resolve();
    }, timeout);
    return () => clearTimeout(timer);
  });
}

/**
 * Create a mock event
 */
export function createMockEvent(
  type: string,
  target: Record<string, unknown> = {}
): Event {
  return {
    type,
    target: {
      value: '',
      checked: false,
      ...target,
    },
    preventDefault: vi.fn(),
    stopPropagation: vi.fn(),
  } as unknown as Event;
}

/**
 * Mock form submission
 */
export function mockFormSubmit(values: Record<string, unknown>) {
  return {
    preventDefault: vi.fn(),
    stopPropagation: vi.fn(),
    target: {
      ...values,
    },
  };
}

/**
 * Helper to mock window.location
 */
export function mockWindowLocation(url: string) {
  delete (window as unknown as { location?: Partial<Location> }).location;
  window.location = {
    href: url,
    origin: 'http://localhost:5173',
    pathname: url.replace(window.location.origin, ''),
    search: '',
    hash: '',
    assign: vi.fn(),
    reload: vi.fn(),
    replace: vi.fn(),
  } as unknown as Location;
}

/**
 * Reset all mocks
 */
export function resetAllMocks() {
  vi.clearAllMocks();
  mockNavigate.mockReset();
}

/**
 * Cleanup after tests
 */
export function cleanup() {
  resetAllMocks();
}

/**
 * Wait for async updates to complete
 */
export async function flushPromises(): Promise<void> {
  return new Promise((resolve) => {
    setTimeout(resolve, 0);
  });
}

/**
 * Mock permissions
 */
export const mockPermissions = {
  USER_PERMISSIONS: {
    LIST: 'user:list',
    CREATE: 'user:create',
    UPDATE: 'user:update',
    DELETE: 'user:delete',
    STATUS: 'user:status',
  },
  ORDER_PERMISSIONS: {
    LIST: 'order:list',
    CREATE: 'order:create',
    UPDATE: 'order:update',
    DELETE: 'order:delete',
    CANCEL: 'order:cancel',
    REFUND: 'order:refund',
  },
  PLAYER_PERMISSIONS: {
    LIST: 'player:list',
    AUDIT: 'player:audit',
    UPDATE: 'player:update',
    DELETE: 'player:delete',
  },
};

/**
 * Mock localStorage
 */
export function mockLocalStorage() {
  const store: Record<string, string> = {
    token: 'test-token',
    user_role: 'admin',
    user_info: JSON.stringify(mockUser),
  };

  return {
    getItem: (key: string) => store[key] || null,
    setItem: (key: string, value: string) => {
      store[key] = value;
    },
    removeItem: (key: string) => {
      delete store[key];
    },
    clear: () => {
      Object.keys(store).forEach((key) => delete store[key]);
    },
    length: Object.keys(store).length,
    key: (index: number) => Object.keys(store)[index] || null,
  };
}

/**
 * Test helper to verify API was called with correct params
 */
export function expectApiCalled(
  fn: ReturnType<typeof vi.fn>,
  ...args: unknown[]
) {
  expect(fn).toHaveBeenCalledWith(...args);
}

/**
 * Test helper to verify number of API calls
 */
export function expectApiCallCount(fn: ReturnType<typeof vi.fn>, count: number) {
  expect(fn).toHaveBeenCalledTimes(count);
}

// Re-export testing library utilities
export * from '@testing-library/react';
export { default as userEvent } from '@testing-library/user-event';
