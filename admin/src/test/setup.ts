import '@testing-library/jest-dom';
import { vi } from 'vitest';
import { ReactElement } from 'react';

// Fix jsdom CSS variable parsing issue
// Patch CSSStyleDeclaration.prototype to handle CSS custom properties in border shorthand
const originalSetProperty = CSSStyleDeclaration.prototype.setProperty;
CSSStyleDeclaration.prototype.setProperty = function (property: string, value: string | null, priority?: string) {
  // Convert value to string and handle CSS custom properties
  let stringValue = String(value ?? '');

  if (stringValue.includes('var(')) {
    // Handle border shorthand with CSS variables (Ant Design issue)
    if (property === 'border' && stringValue.includes('--ant')) {
      stringValue = '1px solid #d9d9d9';
    }
    // Handle individual border properties with CSS variables
    else if (property.includes('border') && stringValue.includes('--ant')) {
      if (property.includes('width')) stringValue = '1px';
      else if (property.includes('style')) stringValue = 'solid';
      else if (property.includes('color')) stringValue = '#d9d9d9';
      else stringValue = stringValue.replace(/var\(--[^)]+\)/g, 'initial');
    }
    // Replace other CSS variables with initial values
    else {
      stringValue = stringValue.replace(/var\(--[^)]+\)/g, 'initial');
    }
  }

  return originalSetProperty.call(this, property, stringValue, priority);
};

// Mock window.matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: (query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: () => {},
    removeListener: () => {},
    addEventListener: () => {},
    removeEventListener: () => {},
    dispatchEvent: () => false,
  }),
});

// Mock ResizeObserver
window.ResizeObserver = class ResizeObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
} as unknown as typeof ResizeObserver;

// Mock IntersectionObserver
window.IntersectionObserver = class IntersectionObserver {
  root = null;
  rootMargin = '';
  thresholds = [];
  observe() {}
  unobserve() {}
  disconnect() {}
  takeRecords() { return []; }
} as unknown as typeof IntersectionObserver;

// Mock scrollTo
window.scrollTo = () => {};

// Mock localStorage
const localStorageMock = (() => {
  let store: Record<string, string> = {};
  return {
    getItem: (key: string) => store[key] || null,
    setItem: (key: string, value: string) => { store[key] = value.toString(); },
    removeItem: (key: string) => { delete store[key]; },
    clear: () => { store = {}; },
    length: 0,
    key: () => null,
  };
})();
Object.defineProperty(window, 'localStorage', { value: localStorageMock });

// Mock window.location with pathname
Object.defineProperty(window, 'location', {
  writable: true,
  value: {
    pathname: '/admin/dashboard',
    href: 'http://localhost:5173/admin/dashboard',
    origin: 'http://localhost:5173',
    protocol: 'http:',
    host: 'localhost:5173',
    hostname: 'localhost',
    port: '5173',
  },
});

// Mock sessionStorage
const sessionStorageMock = (() => {
  let store: Record<string, string> = {};
  return {
    getItem: (key: string) => store[key] || null,
    setItem: (key: string, value: string) => { store[key] = value.toString(); },
    removeItem: (key: string) => { delete store[key]; },
    clear: () => { store = {}; },
    length: 0,
    key: () => null,
  };
})();
Object.defineProperty(window, 'sessionStorage', { value: sessionStorageMock });

// Mock antd App.useApp() hook
vi.mock('antd', async () => {
  const actual = await vi.importActual('antd');
  return {
    ...actual as object,
    App: {
      useApp: () => ({
        message: {
          success: vi.fn(),
          error: vi.fn(),
          warning: vi.fn(),
          info: vi.fn(),
          loading: vi.fn(),
          open: vi.fn(),
        },
        modal: {
          confirm: vi.fn(),
          info: vi.fn(),
          success: vi.fn(),
          error: vi.fn(),
          warning: vi.fn(),
        },
        notification: {
          success: vi.fn(),
          error: vi.fn(),
          info: vi.fn(),
          warning: vi.fn(),
          open: vi.fn(),
        },
      }),
    },
  };
});

// Mock AdminContext to provide super admin permissions for tests
// This ensures PermissionGuard components allow all actions in tests
vi.mock('@/context/AdminContext', () => ({
  default: {
    Provider: ({ children }: { children: ReactElement }) => children,
  },
}));

vi.mock('@/context/useAdmin', () => ({
  useAdmin: () => ({
    rawMenus: [],
    menus: [],
    permissions: ['*'], // Super admin - all permissions
    loading: false,
    refreshMenus: vi.fn().mockResolvedValue(undefined),
    hasPermission: vi.fn(() => true), // Always return true
    hasAllPermissions: vi.fn(() => true),
    hasAnyPermission: vi.fn(() => true),
    isSuperAdmin: true, // Super admin mode
    permissionVersion: 0,
    notifyPermissionChange: vi.fn(),
  }),
}));
