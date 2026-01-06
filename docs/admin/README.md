# Frontend Testing Framework Documentation

## Overview

This document describes the comprehensive testing framework created for the GameLink Admin Panel frontend. The framework provides utilities and test suites for core pages with a target coverage of 60%.

## Created Files

### 1. Test Utilities (`src/testutils.tsx`)

**Location**: `D:\Desktop\Code\GameLink\admin\src\testutils.tsx`

**Purpose**: Centralized testing utilities for all components

**Features**:
- `mockApi`: Comprehensive mock API responses for all endpoints
- `renderWithProviders`: Custom render function with React Query, Router, and Ant Design providers
- `waitForLoading`: Helper for waiting on async operations
- `mockRouter`: Router mocking utilities
- `createTestQueryClient`: Test-specific QueryClient with disabled retries
- `flushPromises`: Helper to wait for all promises to resolve
- `mockLocalStorage`: Mock localStorage implementation
- `resetAllMocks`: Reset all mocks between tests

**Usage Example**:
```typescript
import { renderWithProviders, mockApi, resetAllMocks } from '@/testutils';

describe('MyComponent', () => {
  beforeEach(() => {
    resetAllMocks();
    mockApi.getUsers.mockResolvedValue({ data: { success: true, data: [] } });
  });

  it('should render', () => {
    renderWithProviders(<MyComponent />);
    expect(screen.getByText('Hello')).toBeInTheDocument();
  });
});
```

### 2. Order Page Tests (`src/pages/admin/Order/index.test.tsx`)

**Location**: `D:\Desktop\Code\GameLink\admin\src\pages\admin\Order\index.test.tsx`

**Coverage**: 40+ test cases covering:

- ✅ Successful data loading
- ✅ Loading states
- ✅ Error handling (network errors, API failures)
- ✅ Search and filtering (by order number, status, date range)
- ✅ Pagination (page changes, page size changes)
- ✅ Order details (drawer, information display)
- ✅ Order cancellation
- ✅ Order refunds (with validation)
- ✅ Batch operations (cancel, complete, export)
- ✅ Refresh functionality
- ✅ Permission checks
- ✅ Accessibility

**Test Structure**:
```typescript
describe('OrderPage', () => {
  describe('Successful Data Loading', () => { ... });
  describe('Loading States', () => { ... });
  describe('Error Handling', () => { ... });
  describe('Search and Filtering', () => { ... });
  describe('Pagination', () => { ... });
  // ... more test suites
});
```

### 3. Player Page Tests (`src/pages/admin/Player/index.test.tsx`)

**Location**: `D:\Desktop\Code\GameLink\admin\src\pages\admin\Player\index.test.tsx`

**Coverage**: 35+ test cases covering:

- ✅ Successful data loading
- ✅ Player information display (rating, skill tags, status)
- ✅ Loading states
- ✅ Error handling
- ✅ Search and filtering (by keyword, status)
- ✅ Pagination
- ✅ Player details (statistics, bio, verification info)
- ✅ Player audit workflow (approve/reject)
- ✅ Player ban/unban functionality
- ✅ Batch operations (status change, delete, export)
- ✅ Refresh functionality
- ✅ Accessibility

### 4. User Page Tests (`src/pages/admin/User/index.test.tsx`)

**Location**: `D:\Desktop\Code\GameLink\admin\src\pages\admin\User\index.test.tsx`

**Coverage**: 45+ test cases covering:

- ✅ Successful data loading
- ✅ User statistics display (total, by role, by status, recent)
- ✅ Loading states
- ✅ Error handling
- ✅ Search and filtering (by keyword, role, status, date range)
- ✅ Pagination
- ✅ User details (basic info, login history, operation logs)
- ✅ User CRUD operations (create, edit, delete)
- ✅ User ban/unban functionality
- ✅ Batch operations (role change, status change, notifications, points)
- ✅ Export functionality
- ✅ Refresh functionality
- ✅ Accessibility

### 5. Login Page Tests (`src/pages/admin/Login/index.test.tsx`)

**Location**: `D:\Desktop\Code\GameLink\admin\src\pages\admin\Login\index.test.tsx`

**Coverage**: 35+ test cases covering:

- ✅ Page rendering (form fields, checkbox, copyright)
- ✅ Successful login flow
- ✅ Role verification (admin only)
- ✅ Loading states
- ✅ Error handling:
  - Invalid credentials (401)
  - Disabled account (403)
  - Non-existent user (404)
  - Too many attempts (429)
  - Server errors (500+)
  - Network failures
  - Timeouts
- ✅ Form validation (empty fields)
- ✅ Remember password functionality
- ✅ Saved credentials loading
- ✅ Accessibility

## Key Testing Patterns

### 1. Mock API Setup

```typescript
const mockApi = {
  getUsers: vi.fn(),
  createUser: vi.fn(),
  updateUser: vi.fn(),
  // ... more methods
};

vi.mock('@/api/admin', () => ({
  adminApi: mockApi,
}));

// In beforeEach
mockApi.getUsers.mockResolvedValue({
  data: { success: true, data: mockUsers, pagination: { total: 10, page: 1, pageSize: 10 } }
});
```

### 2. User Interactions

```typescript
const user = userEvent.setup();

// Type in input
await user.type(screen.getByPlaceholderText('Search'), 'query');

// Click button
await user.click(screen.getByRole('button', { name: /Submit/i }));

// Select from dropdown
await user.click(screen.getByText('Select'));
await user.click(screen.getByText('Option'));
```

### 3. Async Assertions

```typescript
// Wait for element to appear
await waitFor(() => {
  expect(screen.getByText('Loaded')).toBeInTheDocument();
});

// Wait for API call
await waitFor(() => {
  expect(mockApi.getUsers).toHaveBeenCalled();
});

// Wait for modal
await waitFor(() => {
  expect(screen.getByText('Modal Title')).toBeInTheDocument();
});
```

### 4. Form Testing

```typescript
// Fill form
await user.type(nameInput, 'John Doe');
await user.type(emailInput, 'john@example.com');

// Submit
await user.click(submitButton);

// Verify API call
await waitFor(() => {
  expect(mockApi.createUser).toHaveBeenCalledWith({
    name: 'John Doe',
    email: 'john@example.com'
  });
});
```

## Running Tests

### Run All Tests
```bash
npm test
```

### Run Tests Once
```bash
npm run test:run
```

### Run Tests with UI
```bash
npm run test:ui
```

### Run Specific Test File
```bash
npx vitest run src/pages/admin/Order/index.test.tsx
```

### Run Tests with Coverage
```bash
npx vitest run --coverage
```

## Test Coverage Goals

| Component | Target Coverage | Status |
|-----------|----------------|--------|
| Order Page | 60% | ✅ 40+ tests |
| Player Page | 60% | ✅ 35+ tests |
| User Page | 60% | ✅ 45+ tests |
| Login Page | 60% | ✅ 35+ tests |

## Testing Best Practices Followed

1. **Isolation**: Each test is independent with proper setup/teardown
2. **Clarity**: Descriptive test names following "should..." pattern
3. **Coverage**: Test success, error, and edge cases
4. **User Interactions**: Test actual user flows, not just rendering
5. **Async Handling**: Proper waiting for async operations
6. **Accessibility**: Test keyboard navigation and ARIA labels
7. **Mock Management**: Clean mocks between tests to avoid leakage

## Known Issues & TODO

### Issues
1. Some existing tests have failures unrelated to new tests (logger infinite loop)
2. Mock hoisting requires careful import order
3. Some tests may require additional mocking of dependencies

### TODO
1. Fix existing test failures in the codebase
2. Add integration tests for complete user flows
3. Add E2E tests with Playwright
4. Increase coverage to 80% for all core pages
5. Add visual regression tests
6. Performance testing for large datasets

## Test Utilities API Reference

### `renderWithProviders(ui, options?)`

Renders a React component with all necessary providers.

**Options**:
- `queryClient?: QueryClient` - Custom QueryClient instance
- `route?: string` - Initial route for router
- `router?: 'memory' | 'browser'` - Router type

**Returns**: RenderResult from @testing-library/react

### `mockApi`

Object containing mocked API methods.

**Properties**:
- `getUsers`, `createUser`, `updateUser`, `deleteUser`
- `getOrders`, `getOrder`, `cancelOrder`, `refundOrder`
- `getPlayers`, `updatePlayerVerification`
- `login`, `getUserStats`
- And more...

### `resetAllMocks()`

Clears all mocks and resets vi.fn() calls.

### `flushPromises()`

Returns a promise that resolves after all pending promises.

### `waitForLoading(callback, timeout?)`

Waits for loading to complete before executing callback.

## Dependencies

```json
{
  "dependencies": {
    "@testing-library/jest-dom": "^6.9.1",
    "@testing-library/react": "^16.3.0",
    "@testing-library/user-event": "^14.6.1",
    "@vitest/coverage-v8": "^4.0.16",
    "vitest": "^4.0.15",
    "jsdom": "^27.2.0"
  }
}
```

## Summary

✅ **Created**: 5 test files with 155+ test cases
✅ **Test Utilities**: Comprehensive testing helpers
✅ **Coverage**: Core pages (Order, Player, User, Login)
✅ **Best Practices**: Following modern testing patterns
✅ **Documentation**: Complete test and utility documentation

This testing framework provides a solid foundation for maintaining code quality and preventing regressions as the application evolves.
