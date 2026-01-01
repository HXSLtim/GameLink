# E2E Testing with Playwright

Comprehensive end-to-end testing for the GameLink Admin Panel using Playwright.

## Overview

This E2E test suite covers critical admin workflows including:
- Authentication and authorization
- User management (CRUD operations)
- Order management (view, cancel, refund)
- Payment management (view, process refunds)
- Player management (CRUD, verification)

## Test Structure

```
tests/e2e/
├── fixtures/
│   └── test-data.fixture.ts    # Test data generation
├── helpers/
│   └── api-helpers.ts           # API utilities for setup/teardown
├── pages/
│   ├── LoginPage.ts             # Login page object
│   ├── UserManagementPage.ts    # User management page object
│   ├── OrderManagementPage.ts   # Order management page object
│   ├── PaymentManagementPage.ts # Payment management page object
│   └── PlayerManagementPage.ts  # Player management page object
├── auth.spec.ts                 # Authentication tests
├── user-management.spec.ts      # User management tests
├── order-management.spec.ts     # Order management tests
├── payment-management.spec.ts   # Payment management tests
└── player-management.spec.ts    # Player management tests
```

## Prerequisites

1. **Backend API** running on `http://localhost:8080`
2. **Admin user** credentials (default: `admin` / `admin123`)
3. **PostgreSQL** database with test data

## Environment Variables

Create a `.env` file in the `admin/` directory:

```bash
# Admin credentials for E2E tests
TEST_ADMIN_USERNAME=admin
TEST_ADMIN_PASSWORD=admin123
TEST_ADMIN_EMAIL=admin@gamelink.com

# API URL (default: http://localhost:8080/api/v1)
API_URL=http://localhost:8080/api/v1

# Admin panel URL (default: http://localhost:5173)
BASE_URL=http://localhost:5173
```

## Installation

```bash
cd admin
npm install
npx playwright install chromium
```

## Running Tests

### Run all E2E tests (headless)
```bash
npm run test:e2e
```

### Run tests with UI mode
```bash
npm run test:e2e:ui
```

### Run tests in headed mode (show browser)
```bash
npm run test:e2e:headed
```

### Debug tests
```bash
npm run test:e2e:debug
```

### Run specific test file
```bash
npx playwright test auth.spec.ts
```

### Run tests matching a pattern
```bash
npx playwright test --grep "Authentication"
```

### View test report
```bash
npm run test:e2e:report
```

## Test Coverage

### Authentication Tests (`auth.spec.ts`)
- ✅ Login with valid credentials
- ✅ Login with invalid credentials
- ✅ Form validation
- ✅ Session persistence
- ✅ Logout functionality
- ✅ Token validation
- ✅ UI/UX interactions
- ✅ Accessibility

**Total**: 15+ test cases

### User Management Tests (`user-management.spec.ts`)
- ✅ List users with pagination
- ✅ Create new user
- ✅ Edit user information
- ✅ Delete user
- ✅ Search and filter users
- ✅ View user details
- ✅ Batch operations
- ✅ Export functionality
- ✅ Form validation

**Total**: 20+ test cases

### Order Management Tests (`order-management.spec.ts`)
- ✅ List orders with pagination
- ✅ View order details
- ✅ Cancel orders
- ✅ Refund orders
- ✅ Search and filter orders
- ✅ Batch operations
- ✅ Order status tracking
- ✅ Export functionality

**Total**: 18+ test cases

### Payment Management Tests (`payment-management.spec.ts`)
- ✅ List payment records
- ✅ View payment details
- ✅ Process refunds
- ✅ Search and filter payments
- ✅ Payment status tracking
- ✅ Export functionality
- ✅ Amount display validation

**Total**: 15+ test cases

### Player Management Tests (`player-management.spec.ts`)
- ✅ List players with pagination
- ✅ View player details
- ✅ Create new player
- ✅ Edit player information
- ✅ Delete player
- ✅ Player verification (approve/reject)
- ✅ Search and filter players
- ✅ Batch operations
- ✅ Export functionality

**Total**: 22+ test cases

**Total E2E Tests**: 90+ test cases

## Page Object Model

The tests use the Page Object Model pattern for maintainability:

```typescript
// Example: Using a page object
const loginPage = new LoginPage(page);
await loginPage.goto();
await loginPage.login('username', 'password');
await loginPage.verifyLoginSuccess();
```

### Benefits:
- **Reusable**: Page objects can be used across multiple tests
- **Maintainable**: Changes to UI only require updating page objects
- **Readable**: Test code reads like business requirements

## Test Data Management

### Fixtures
Test data is generated using fixtures with timestamps to ensure uniqueness:

```typescript
test('should create user', async ({ testData }) => {
  await userManagementPage.createUser({
    name: testData.testUser.name,
    email: testData.testUser.email,
    phone: testData.testUser.phone,
    role: 'user',
  });
});
```

### Cleanup
Test data is automatically cleaned up after each test:

```typescript
test.afterEach(async () => {
  if (testUserId) {
    await deleteTestUser(adminToken, testUserId);
  }
});
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: E2E Tests

on: [push, pull_request]

jobs:
  e2e:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '20'

      - name: Install dependencies
        run: |
          cd admin
          npm ci

      - name: Install Playwright
        run: |
          cd admin
          npx playwright install --with-deps

      - name: Start backend
        run: |
          cd api
          go run cmd/main.go &
          sleep 10

      - name: Run E2E tests
        run: |
          cd admin
          npm run test:e2e

      - name: Upload test report
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: playwright-report
          path: admin/playwright-report/
```

## Best Practices

### 1. Test Isolation
Each test should be independent and clean up after itself:

```typescript
test.beforeEach(async ({ page }) => {
  // Setup: Login and navigate
  await loginPage.goto();
  await loginPage.login(username, password);
});

test.afterEach(async ({ page }) => {
  // Cleanup: Delete test data
  await deleteTestUser(testUserId);
});
```

### 2. Wait Strategies
Use Playwright's auto-waiting features:

```typescript
// ✅ Good - Auto-wait for element
await page.click('button');

// ❌ Bad - Manual sleep
await page.waitForTimeout(1000);
await page.click('button');
```

### 3. Selectors
Use accessible selectors:

```typescript
// ✅ Good - Role-based
await page.getByRole('button', { name: 'Submit' }).click();

// ✅ Good - Text-based
await page.getByText('Welcome').click();

// ❌ Bad - CSS selectors
await page.$('.btn-primary').click();
```

### 4. Assertions
Use specific assertions:

```typescript
// ✅ Good - Specific assertion
await expect(page).toHaveURL('/dashboard');
await expect(message).toHaveText('Success');

// ❌ Bad - Generic
expect(page.url()).toContain('dashboard');
```

## Debugging

### Playwright Inspector
```bash
npm run test:e2e:debug
```

### Screenshots on Failure
Screenshots are automatically captured on test failures and saved to:
- `test-results/`

### Traces
Traces are captured for failed tests:
```bash
npx playwright show-trace test-results/[test-name]/trace.zip
```

## Troubleshooting

### Tests fail with "Browser not found"
```bash
npm run test:e2e:install
```

### Tests timeout
- Increase timeout in `playwright.config.ts`
- Check if backend is running
- Verify network connectivity

### Flaky tests
- Use specific wait conditions
- Avoid hard-coded timeouts
- Check for race conditions

## Maintenance

### Update locators
When UI changes, update page objects:
```typescript
// pages/UserManagementPage.ts
private readonly createButton = this.page.getByRole('button', { name: /新增|create/i });
```

### Add new tests
Follow the existing pattern:
1. Create page object methods
2. Write test cases
3. Add test data fixture if needed
4. Update documentation

## Performance

### Run tests in parallel
Tests run in parallel by default. To disable:
```typescript
// playwright.config.ts
export default defineConfig({
  workers: 1, // Disable parallel execution
});
```

### Test execution time
- Full suite: ~5-10 minutes (depends on test data)
- Single file: ~30 seconds - 2 minutes

## Resources

- [Playwright Documentation](https://playwright.dev/)
- [Page Object Model](https://playwright.dev/docs/pom)
- [Best Practices](https://playwright.dev/docs/best-practices)
- [API Reference](https://playwright.dev/docs/api/class-playwright)

## Contributing

When adding new E2E tests:

1. **Page Object First**: Add methods to page objects before writing tests
2. **Test Data**: Use fixtures for test data generation
3. **Cleanup**: Ensure test data is cleaned up in `afterEach`
4. **Documentation**: Update this README with new test coverage
5. **Review**: Run tests locally before pushing

## License

See project root LICENSE file.
