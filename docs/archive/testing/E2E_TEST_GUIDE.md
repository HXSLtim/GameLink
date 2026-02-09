# End-to-End Integration Tests Guide

## Overview

This document describes the comprehensive end-to-end (E2E) integration tests for the GameLink platform. These tests validate complete user journeys across the entire system, from frontend to backend.

## Test Architecture

### Test Files

| File | Purpose | Key Scenarios |
|------|---------|---------------|
| `e2e_helper.go` | Shared utilities and HTTP test client | Helper functions, assertions, test data setup |
| `e2e_order_test.go` | Order lifecycle (MOST CRITICAL) | Solo, team, gift orders, cancellation, timeout, T+7 settlement |
| `e2e_payment_test.go` | Payment and withdrawal flow | Recharge, payment methods, refund, withdrawal, limits |
| `e2e_dispute_test.go` | Dispute resolution workflow | Initiation, assignment, resolution, SLA, rollback, batch ops |
| `e2e_vip_test.go` | User registration and VIP | Referral codes, VIP purchase, benefits, experience tracking, coupons |
| `e2e_content_test.go` | Content moderation workflow | Sensitive words, approval/rejection, reports, batch moderation |

## Running E2E Tests

### Prerequisites

1. **Start Test Database**
   ```bash
   docker-compose -f docker-compose.test.yml up -d
   ```

2. **Set Environment Variables**
   ```bash
   export TEST_DB_HOST=localhost
   export TEST_DB_PORT=5432
   export TEST_DB_USER=gamelink
   export TEST_DB_PASSWORD=gamelink
   export TEST_DB_NAME=gamelink_test
   ```

### Running Tests

```bash
# Run all E2E tests
cd api
go test ./internal/service/integration/... -run TestE2E -v

# Run specific E2E test
go test ./internal/service/integration/... -run TestE2E_OrderLifecycle_Solo -v

# Run with coverage
go test ./internal/service/integration/... -run TestE2E -cover -coverprofile=e2e_coverage.out
```

## Test Scenarios

### 1. Order Lifecycle (`e2e_order_test.go`)

#### TestE2E_OrderLifecycle_Solo ⭐ MOST CRITICAL
**Purpose**: Validates complete solo order flow from creation to settlement

**Steps**:
1. Customer creates solo order
2. Customer makes payment (wallet)
3. Player accepts order
4. Player completes service
5. System calculates commission (20%)
6. Player income frozen (T+7)
7. User leaves review

**Validations**:
- Order status transitions: pending → confirmed → in_progress → completed
- Payment status: pending → paid
- Wallet balance deducted correctly
- Commission calculated: 20% platform, 80% player
- Player frozen balance = 80% of order amount
- Review created (pending approval)

#### TestE2E_OrderLifecycle_Team
**Purpose**: Validates team order flow with multiple players

**Steps**:
1. Customer creates team order (5 players needed)
2. Order items created for each slot
3. Players assigned to slots
4. Verify team composition

#### TestE2E_OrderLifecycle_Gift
**Purpose**: Validates gift order flow

**Steps**:
1. Customer creates gift order for another user
2. Order type = gift
3. Receiver ID set
4. Payment processed immediately

#### TestE2E_OrderCancellation
**Purpose**: Validates order cancellation and refund

**Steps**:
1. Customer creates and pays for order
2. Customer cancels order
3. Payment refunded automatically
4. Wallet balance restored

#### TestE2E_OrderTimeout
**Purpose**: Validates order timeout handling (24h)

**Steps**:
1. Create order with timestamp 25 hours ago
2. Verify order can be cancelled due to timeout

#### TestE2E_OrderWithVIP
**Purpose**: Validates VIP discount application

**Steps**:
1. Create VIP user with gold level
2. User creates order
3. Verify 5% discount applied (95% of original price)

#### TestE2E_OrderWithCoupon
**Purpose**: Validates coupon discount application

**Steps**:
1. Issue coupon to user
2. User creates order with coupon
3. Verify coupon discount applied
4. Verify coupon marked as used

#### TestE2E_TPlus7Settlement
**Purpose**: Validates T+7 settlement for player income

**Steps**:
1. Create order completed 8 days ago
2. Create commission record (pending)
3. Create player wallet with frozen balance
4. Process settlement
5. Verify frozen → available balance transfer

### 2. Payment and Withdrawal (`e2e_payment_test.go`)

#### TestE2E_PaymentWithdrawal_FullFlow ⭐ CRITICAL
**Purpose**: Validates complete payment and withdrawal cycle

**Steps**:
1. Customer recharges wallet (50,000 cents)
2. Customer places 2 orders (10,000 + 15,000 cents)
3. Players complete orders
4. Income frozen in player wallets
5. T+7 settlement processed (8 days later)
6. Player requests withdrawal (5,000 cents)
7. Admin approves withdrawal
8. Payment gateway transfers money
9. Withdrawal marked complete

**Validations**:
- Customer wallet: 50,000 → 25,000 → 0
- Player 1 frozen: 0 → 8,000 → 0 (settled)
- Player 1 available: 0 → 8,000 → 3,000 (after withdrawal)
- Withdrawal status: pending → approved → completed
- Daily/monthly limits enforced
- Idempotent processing (no double-pay)

#### TestE2E_Payment_CombinedPayment
**Purpose**: Validates wallet + third-party combined payment

**Steps**:
1. Customer has 5,000 cents in wallet
2. Order costs 10,000 cents
3. Use 5,000 from wallet + 5,000 from WeChat
4. Verify split payment recorded

#### TestE2E_Payment_Refund
**Purpose**: Validates payment refund on order cancellation

**Steps**:
1. Customer pays for order (wallet)
2. Customer cancels order
3. Payment refunded automatically
4. Wallet balance restored

#### TestE2E_Payment_ThirdParty
**Purpose**: Validates third-party payment (WeChat/Alipay)

**Steps**:
1. Customer creates WeChat payment
2. Payment URL generated
3. Simulate payment webhook
4. Order status updated to confirmed

#### TestE2E_Withdraw_Rejection
**Purpose**: Validates withdrawal rejection flow

**Steps**:
1. Player requests withdrawal
2. Admin rejects (wrong account info)
3. Wallet unfrozen
4. Withdrawal marked rejected

#### TestE2E_Withdraw_Idempotency
**Purpose**: Validates no double-pay on withdrawal

**Steps**:
1. Player requests withdrawal
2. Admin approves
3. Try to approve again (should fail)
4. Verify only one transaction

#### TestE2E_Payment_MultipleMethods
**Purpose**: Validates various payment methods

**Methods Tested**:
- Wallet payment
- WeChat Pay
- Alipay
- Combined payment

#### TestE2E_Payment_DisputeRefund
**Purpose**: Validates refund during dispute resolution

**Steps**:
1. Customer pays for order
2. Dispute initiated and resolved with refund
3. Payment refunded
4. Order marked refunded
5. Commission cancelled

### 3. Dispute Resolution (`e2e_dispute_test.go`)

#### TestE2E_DisputeResolution_FullFlow ⭐ CRITICAL
**Purpose**: Validates complete dispute resolution workflow

**Steps**:
1. Complete order successfully
2. User initiates dispute within 7-day window
3. System sets 30-minute SLA deadline
4. Admin assigns dispute to customer service
5. CS investigates (creates chat snapshot)
6. Independent CS reviews (dual-CS mechanism)
7. CS resolves with refund
8. Payment refunded
9. Commission cancelled
10. Notifications sent

**Validations**:
- 7-day filing window enforced
- 30-minute SLA tracked
- Dual-CS mechanism (original + independent)
- Dispute status: pending → assigned → mediating → resolved
- Order status: disputed → refunded
- Payment status: paid → refunded
- Commission status: pending → cancelled

#### TestE2E_Dispute_Rejection
**Purpose**: Validates dispute rejection (user was wrong)

**Steps**:
1. User initiates dispute
2. CS investigates
3. CS rejects dispute
4. Order status restored to completed
5. Commission not affected

#### TestE2E_Dispute_PartialRefund
**Purpose**: Validates partial refund (25%)

**Steps**:
1. User initiates dispute (service duration issue)
2. CS resolves with partial refund
3. 25% refunded, 75% kept by player

#### TestE2E_Dispute_PlayerInitiated
**Purpose**: Validates dispute initiated by player

**Steps**:
1. Player initiates dispute (user not cooperative)
2. CS investigates
3. CS rejects player's dispute

#### TestE2E_Dispute_7DayWindow
**Purpose**: Validates 7-day filing window

**Steps**:
1. Try to file dispute for 8-day-old order (fail)
2. File dispute for 5-day-old order (success)

#### TestE2E_Dispute_SLABreach
**Purpose**: Validates SLA breach handling

**Steps**:
1. Create dispute
2. Set SLA deadline to past
3. Run SLA breach check
4. Dispute marked as SLA breached

#### TestE2E_Dispute_Rollback
**Purpose**: Validates dispute assignment rollback

**Steps**:
1. Assign dispute to CS1
2. CS1 rolls back (cannot handle)
3. Re-assign to CS2

#### TestE2E_Dispute_BatchOperations
**Purpose**: Validates batch dispute operations

**Steps**:
1. Create 5 disputes
2. Batch assign all to CS
3. Batch update 3 to mediating
4. Batch close all with rejection

### 4. User Registration and VIP (`e2e_vip_test.go`)

#### TestE2E_UserRegistration_WithReferral
**Purpose**: Validates user registration with referral code

**Steps**:
1. Create referrer (existing user)
2. Generate referral code for referrer
3. New user registers with referral code
4. Complete referral after registration
5. Verify rewards issued

**Validations**:
- Referral code generated
- Referral record created
- Referrer reward: 1,000 points
- Referee reward: 500 points
- Code usage count updated

#### TestE2E_VIP_Purchase
**Purpose**: Validates VIP membership purchase and benefits

**Steps**:
1. Create VIP levels (bronze, silver, gold)
2. User purchases VIP (add experience)
3. User upgrades to higher level
4. Apply VIP discount to orders
5. Handle VIP expiration
6. Renew VIP membership

**Validations**:
- VIP level progression
- Experience points tracking
- Discount application (2%, 5%, 10%)
- Expiration handling

#### TestE2E_VIP_Benefits
**Purpose**: Validates VIP benefit application

**Benefits Tested**:
- Order discount (2%, 5%, 10%)
- Free orders (0, 1, 3)
- Priority support (false/true)
- Exclusive events (false/true)
- Commission discount (0%, 5%)

#### TestE2E_VIP_ExperienceTracking
**Purpose**: Validates VIP experience point tracking

**Steps**:
1. Create VIP levels (100, 500, 1000 exp)
2. User starts with 0 exp
3. Add experience: 50 → still level1
4. Add experience: 100 → total 150 → still level1
5. Add experience: 400 → total 550 → level2
6. Add experience: 600 → total 1150 → level3

#### TestE2E_Coupon_Usage
**Purpose**: Validates coupon usage with orders

**Steps**:
1. Create coupon template
2. Issue coupon to user
3. User uses coupon on order
4. Verify discount applied
5. Try to reuse coupon (fail)

### 5. Content Moderation (`e2e_content_test.go`)

#### TestE2E_ContentModeration_FullFlow ⭐ CRITICAL
**Purpose**: Validates complete content moderation workflow

**Steps**:
1. Configure sensitive words (垃圾, 骗子, 微信)
2. User creates feed with sensitive content
3. System auto-detects sensitive words
4. Feed marked as pending review
5. Admin reviews and approves
6. Feed becomes visible
7. User receives notification

**Validations**:
- Sensitive words detected correctly
- Feed auto-moderated
- Admin approval workflow
- Notifications sent

#### TestE2E_ContentModeration_Rejection
**Purpose**: Validates feed rejection

**Steps**:
1. User creates feed with critical sensitive word
2. Admin rejects feed
3. User notified of rejection

#### TestE2E_ContentModeration_BatchOperations
**Purpose**: Validates batch content moderation

**Steps**:
1. Create 5 pending feeds
2. Batch approve 3 feeds
3. Batch reject 2 feeds

#### TestE2E_ContentReport
**Purpose**: Validates user reporting inappropriate content

**Steps**:
1. User creates feed (approved)
2. Another user reports feed
3. Admin investigates
4. Admin rejects feed
5. Report marked as processed

#### TestE2E_SensitiveWordManagement
**Purpose**: Validates sensitive word management

**Steps**:
1. Create sensitive words
2. List sensitive words
3. Check content for sensitive words
4. Filter content (replace words)
5. Update sensitive word
6. Deactivate sensitive word

#### TestE2E_ContentCategories
**Purpose**: Validates content category management

**Steps**:
1. Create categories (游戏攻略, 陪玩分享, 日常闲聊)
2. Create feed with category
3. Update category
4. Deactivate category

## Key Validations

### Database Assertions

The E2E tests use helper functions to validate database state:

```go
// Assert order status
AssertOrderStatus(t, db, orderID, "completed")

// Assert payment status
AssertPaymentStatus(t, db, paymentID, "paid")

// Assert wallet balance
AssertWalletBalance(t, db, userID, expectedBalance, expectedFrozen)

// Assert commission record exists
AssertCommissionRecord(t, db, orderID)

// Assert dispute status
AssertDisputeStatus(t, db, disputeID, "resolved")
```

### Time-Based Operations

```go
// Simulate time travel for testing T+7
TimeTravel(t, db, orderID, daysAgo int)

// Wait for async operations
WaitForCondition(t, condition func() bool, timeout, msg)
```

## Coverage Goals

| Scenario | Coverage Target | Status |
|----------|----------------|--------|
| Order Lifecycle (Solo) | 100% | ✅ |
| Payment Flow | 100% | ✅ |
| Withdrawal Flow | 100% | ✅ |
| Dispute Resolution | 100% | ✅ |
| VIP Benefits | 100% | ✅ |
| Content Moderation | 100% | ✅ |

## Success Criteria

All E2E tests must:
1. ✅ Pass without manual intervention
2. ✅ Clean up test data after each test
3. ✅ Validate database state after each step
4. ✅ Test both success and failure scenarios
5. ✅ Cover all critical business rules

## Troubleshooting

### Tests Skipped

If tests are skipped:
```bash
# Check database is running
docker-compose -f docker-compose.test.yml ps

# Check environment variables
env | grep TEST_DB
```

### Timeout Errors

Increase timeout:
```bash
go test ./internal/service/integration/... -run TestE2E -timeout 30m
```

### Database Lock Errors

Run tests sequentially:
```bash
go test ./internal/service/integration/... -run TestE2E -p 1
```

## Best Practices

1. **Always clean up test data** - Each test should leave the database in a clean state
2. **Use unique test data** - Use timestamps or UUIDs to avoid collisions
3. **Validate at each step** - Don't wait until the end to check results
4. **Test edge cases** - Not just happy path
5. **Use descriptive test names** - TestE2E_Feature_Scenario

## CI/CD Integration

The E2E tests run in CI/CD pipeline:
- **ci.yml**: Runs on every PR (change detection)
- **nightly.yml**: Runs full E2E suite nightly
- **deploy.yml**: Runs smoke tests before deployment

## Related Documentation

- [Integration Test Plan](./INTEGRATION_TEST_PLAN.md)
- [Testing Standards](../.kiro/steering/05-testing-standard.md)
- [Project Management](../.kiro/steering/06-project-management.md)
