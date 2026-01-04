# Comprehensive Payment & Withdraw Test Suite Documentation

## Overview

This document describes the comprehensive test suite created for the Payment and Withdraw modules in the GameLink project. These tests focus on **financial accuracy**, **security**, and **edge cases** for critical financial operations.

## Test Files Created

### 1. Payment Handler Tests
**File**: `api/internal/handler/user/payment_handler_test.go`

**Purpose**: Tests for payment HTTP handlers with focus on webhook security.

**Key Test Categories**:
- **Webhook Signature Verification**: Valid/Invalid signature validation
- **Replay Attack Prevention**: Duplicate transaction ID detection
- **Tampered Payload Detection**: Modified payload rejection
- **Idempotency**: Multiple callback handling
- **Financial Calculations**: Wallet/Combined payment accuracy
- **State Machine**: Valid/Invalid payment state transitions
- **Security**: Amount validation, overpayment prevention
- **Edge Cases**: Zero amounts, maximum amounts, one cent precision

**Test Count**: 25+ test functions

### 2. Payment Integration Tests
**File**: `api/internal/service/payment/payment_integration_test.go`

**Purpose**: End-to-end integration tests with real database.

**Key Test Scenarios**:
- **E2E Payment Flows**:
  - Third-party payment (WeChat/Alipay)
  - Wallet payment flow
  - Combined payment flow
  - Refund flow
  - Combined payment refund

- **Error Handling**:
  - Double payment attempts
  - Insufficient wallet balance
  - Invalid order status
  - Unauthorized users

- **Callback Handling**:
  - Valid signature processing
  - Amount mismatch detection
  - Provider mismatch detection
  - Idempotency verification

- **List and Filter**:
  - Filter by status
  - Filter by method
  - Filter by date range
  - Pagination

**Test Count**: 20+ test functions

### 3. Financial Calculation Tests
**File**: `api/internal/service/payment/financial_calculations_test.go`

**Purpose**: Precision tests for financial calculations.

**Test Categories**:

#### Basic Calculations
- Single deduction accuracy
- Combined payment ratio accuracy
- Full refund accuracy
- Partial refund accuracy

#### Edge Cases
- Minimum amount (1 cent / 0.01 yuan)
- Maximum amount (int64 max)
- Zero balance
- Rounding precision

#### Precision Tests
- Cent accuracy (1-100 cents)
- Large amount precision
- Various rounding scenarios

#### Commission Calculations
- Platform commission (15%, 20%, 25%)
- Tiered rates
- Edge case rounding (0.99 yuan)

#### Refund Calculations
- Partial refunds (50%, 10%, 90%)
- Combined payment refunds
- Wallet restoration verification

#### Multiple Transactions
- Sequential deductions
- Refund cycles (pay → refund → pay → refund)

**Test Count**: 30+ test functions

### 4. Security Tests
**File**: `api/internal/service/payment/security_test.go`

**Purpose**: Security-focused tests for payment operations.

**Test Categories**:

#### Replay Attack Prevention
- Duplicate transaction ID detection
- Different amount same transaction ID
- Timestamp replay attacks
- Multiple payments same order

#### Webhook Signature Verification
- Valid signature acceptance
- Invalid signature rejection
- Tampered payload detection
- Different secret key rejection

#### Idempotency
- Create payment same order
- Callback multiple times
- Refund multiple times
- Cancel multiple times

#### Race Conditions
- Concurrent payment creation (10 goroutines)
- Concurrent refunds (5 goroutines)
- Distributed lock verification

#### Authorization
- Cross-user payment attempts
- Cross-user refund attempts

#### Input Validation
- Negative amounts
- Zero amounts
- Invalid payment methods

**Test Count**: 25+ test functions

### 5. Withdraw Integration Tests
**File**: `api/internal/handler/admin/withdraw_integration_test.go`

**Purpose**: End-to-end withdraw tests with financial accuracy.

**Test Categories**:

#### Financial Accuracy
- Balance validation
- Tax calculation (5% test rate)
- Partial refund handling
- Multiple withdrawals

#### Idempotency
- Approve twice
- Complete twice
- Batch operation idempotency

#### Edge Cases
- Zero amount
- Maximum amount
- Insufficient balance
- One cent withdrawal

#### State Machine
- Valid transitions (pending→approved, pending→rejected, approved→completed)
- Invalid transitions (completed→approved)

#### Batch Operations
- Mixed success/failure
- Maximum batch size (100 withdrawals)
- Exceeds maximum (101 withdrawals)

#### Routing
- Company assignment
- No company assignment (error case)

#### Statistics
- Routing stats aggregation
- Company withdrawal stats

**Test Count**: 30+ test functions

### 6. Withdraw Service Tests (Existing)
**File**: `api/internal/service/withdraw/service_test.go`

**Status**: Already exists with comprehensive mock-based tests.

**Coverage**:
- WithdrawRoutingService tests
- WithdrawRoutingStatsService tests
- Batch operations tests
- Model method tests

## Test Constants

### Financial Values
```go
const (
    OneCent       = 1     // 0.01 yuan
    SmallAmount   = 100   // 1 yuan
    MediumAmount  = 5000  // 50 yuan
    StandardAmount= 10000 // 100 yuan
    LargeAmount   = 50000 // 500 yuan
)
```

### Commission Rates
```go
const (
    Commission15Percent = 0.15
    Commission20Percent = 0.20 // Standard rate
    Commission25Percent = 0.25
)
```

### Tax Rates
```go
const (
    TaxRate5Percent = 0.05 // Test tax rate for withdrawals
)
```

## Running the Tests

### Run All Payment Tests
```bash
cd api
go test ./internal/handler/user/... -v -run TestPayment
go test ./internal/service/payment/... -v -run TestPayment
```

### Run All Withdraw Tests
```bash
cd api
go test ./internal/handler/admin/... -v -run TestWithdraw
go test ./internal/service/withdraw/... -v
```

### Run Financial Calculation Tests
```bash
cd api
go test ./internal/service/payment/... -v -run TestFinancialCalculations
```

### Run Security Tests
```bash
cd api
go test ./internal/service/payment/... -v -run TestSecurity
```

### Run with Coverage
```bash
cd api
go test ./internal/service/payment/... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Run with Race Detector
```bash
cd api
go test ./internal/service/payment/... -race -v
```

## Coverage Goals

| Module | Target Coverage | Current Status |
|--------|----------------|----------------|
| Payment Handler | 90%+ | ✅ Tests created |
| Payment Service | 90%+ | ✅ Tests created |
| Withdraw Handler | 90%+ | ✅ Tests created |
| Withdraw Service | 90%+ | ✅ Tests created |

## Test Infrastructure

### Test Database Setup
All integration tests use `testutil.SetupTestDB(t)` which:
- Initializes PostgreSQL test database
- Auto-migrates all models
- Cleans tables before each test
- Uses mutex to prevent concurrent TRUNCATE deadlocks

### Test Helpers

#### Create Test User
```go
testUser := testutil.CreateAdminUser(t, db, model.RoleUser)
```

#### Create Test Player
```go
testPlayer := testutil.CreateTestPlayer(t, db, testUser.ID)
```

#### Create Test Wallet
```go
wallet := ctx.CreateTestWallet(t, userID, balanceCents)
```

#### Create Test Order
```go
order := ctx.CreateTestOrder(t, status, priceCents)
```

#### Create Test Payment
```go
payment := testutil.CreateTestPayment(t, db, order.ID, userID, status)
```

### Mock Objects

#### Mock Payment Gateway
```go
type MockPaymentGateway struct {
    mockCreatePayment func(order *model.Order) (string, error)
    mockRefund func(paymentId string, amount int64) error
}
```

## Security Test Scenarios

### 1. Replay Attack Tests

#### Scenario: Duplicate Transaction ID
```go
// First webhook with transaction ID "TXN_001"
callbackData1 := {trade_no: "TXN_001", ...}
service.HandlePaymentCallback(ctx, "wechat", callbackData1) // Success

// Second webhook with SAME transaction ID
callbackData2 := {trade_no: "TXN_001", ...}
service.HandlePaymentCallback(ctx, "wechat", callbackData2) // Should be rejected or idempotent
```

#### Scenario: Tampered Amount
```go
// Original signature for amount 10000
signature = generateSignature(payload, secret)

// Attacker modifies payload to amount 5000 but keeps original signature
tamperedPayload = {amount: 5000, ...}

// Signature verification should fail
verifySignature(tamperedPayload, signature) // Should reject
```

### 2. Idempotency Tests

#### Scenario: Multiple Refunds
```go
// First refund
service.RefundPayment(ctx, paymentID, "first refund") // Success

// Second refund (idempotent or error)
service.RefundPayment(ctx, paymentID, "second refund")
// Expected: Either succeeds idempotently OR fails with "already refunded"
// Either way, balance should only be credited ONCE
```

### 3. Race Condition Tests

#### Scenario: Concurrent Payment Creation
```go
// Launch 10 goroutines trying to create payment for same order
for i := 0; i < 10; i++ {
    go func() {
        service.CreatePayment(ctx, userID, req)
    }()
}

// Expected: Only 1 succeeds, 9 fail
// Verify: Only 1 payment record in database
```

## Financial Calculation Test Examples

### Wallet Payment Accuracy
```go
Initial Balance: 10000
Order Amount:    5000
Expected Final:  5000

// Verify
assert.Equal(t, int64(5000), resp.WalletDeducted)
assert.Equal(t, int64(5000), wallet.BalanceCents)
```

### Combined Payment Calculation
```go
Order Total:       20000
Wallet Balance:    15000
Requested Wallet:  8000
Expected 3rd Party: 12000

// Verify calculation
assert.Equal(t, int64(8000), resp.WalletDeducted)
assert.Equal(t, int64(12000), resp.ThirdPartyAmount)
assert.Equal(t, int64(20000), resp.WalletDeducted + resp.ThirdPartyAmount)
```

### Refund Calculation
```go
Paid Amount:   10000
Wallet After:  0
Refund Amount: 10000
Expected Final: 10000

// Verify
refundErr := service.RefundPayment(ctx, paymentID, "test refund")
walletAfterRefund := getWallet(userID)
assert.Equal(t, int64(10000), walletAfterRefund.BalanceCents)
```

## Edge Cases Tested

### Minimum Amount
- 1 cent (0.01 yuan) payments
- 1 cent withdrawals
- 1 cent refunds

### Maximum Amount
- int64 maximum value
- Large balance handling (no overflow)

### Zero Amount
- Zero amount rejection
- Zero balance handling

### Rounding
- 0.99 yuan commission calculation
- Odd amount splits in combined payments

## Known Limitations

1. **Distributed Locking**: Some tests mock distributed locks instead of using real Redis
2. **External APIs**: Payment gateways (WeChat/Alipay) are mocked
3. **Bank Transfers**: Actual bank transfers are not tested (would require test bank accounts)

## Continuous Integration

### GitHub Actions Workflow
Add to `.github/workflows/test.yml`:

```yaml
name: Financial Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_DB: gamelink_test
          POSTGRES_USER: gamelink
          POSTGRES_PASSWORD: gamelink
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run Payment Tests
        run: |
          cd api
          go test ./internal/service/payment/... -v -race -coverprofile=coverage.out
      - name: Run Withdraw Tests
        run: |
          cd api
          go test ./internal/service/withdraw/... -v -race
```

## Maintenance

### When to Update Tests

1. **New Payment Method**: Add tests in `payment_handler_test.go`
2. **New Withdraw Status**: Add state machine tests
3. **Commission Rate Change**: Update calculation tests
4. **New Security Feature**: Add security tests

### Test Checklist for New Features

- [ ] Create payment/withdraw test
- [ ] Financial accuracy test
- [ ] Idempotency test
- [ ] Security test (if applicable)
- [ ] Edge case test
- [ ] Integration test

## Contact

For questions about these tests, please refer to:
- `CLAUDE.md` - Project conventions
- `.kiro/steering/05-testing-standard.md` - Testing standards
- `docs/INTEGRATION_TEST_PLAN.md` - Integration test guidelines

---

**Total Test Coverage**: 130+ test functions across 6 test files

**Primary Focus**: Financial accuracy, security, and edge cases

**Target**: 90%+ coverage for Payment and Withdraw modules
