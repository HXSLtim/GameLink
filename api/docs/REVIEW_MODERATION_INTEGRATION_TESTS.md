# Review Moderation Integration Tests

## Overview

This document describes the integration tests implemented for the review moderation functionality in the GameLink platform.

## Test File

`backend/internal/integration/review_moderation_integration_test.go`

## Test Coverage

### 1. TestReviewModerationFlow
**Requirements: 2.1, 2.2, 2.3**

Tests the complete review moderation workflow:
- Creates a pending review
- Lists pending reviews and verifies the review appears in the list
- Approves a review and verifies status changes to `approved`
- Creates another pending review
- Rejects a review with a reason and verifies status changes to `rejected`
- Verifies rejection reason is stored correctly
- Tests that already-approved reviews cannot be re-approved
- Tests that already-rejected reviews cannot be re-rejected

**Key Assertions:**
- Pending reviews appear in the pending list
- Approve operation updates status to `approved`
- Reject operation updates status to `rejected` and stores reason
- Status transitions are enforced (no re-approval/re-rejection)

### 2. TestSensitiveWordAutoMarking
**Requirements: 2.4**

Tests automatic sensitive word detection and marking:
- Creates sensitive words in the database ("垃圾", "骗子", "差评")
- Creates a review containing sensitive words
- Verifies the review appears in the pending list
- Creates a review without sensitive words
- Verifies both reviews are in the pending list
- Approves the clean review
- Rejects the review with sensitive words
- Verifies final statuses are correct

**Key Assertions:**
- Reviews with sensitive words are marked for moderation
- Reviews without sensitive words can be approved
- Sensitive word content is preserved in the review
- Moderators can reject reviews containing sensitive words

### 3. TestBatchModeration
**Requirements: 2.5**

Tests batch approval and rejection operations:
- Creates 5 pending reviews
- Batch approves the first 3 reviews
- Verifies the first 3 reviews are approved
- Verifies the last 2 reviews remain pending
- Batch rejects the last 2 reviews with a reason
- Verifies the last 2 reviews are rejected with the correct reason
- Tests atomicity: attempts to batch approve a mix of pending and non-pending reviews
- Verifies atomicity: all reviews maintain their original status when batch operation fails

**Key Assertions:**
- Batch approve successfully updates multiple reviews
- Batch reject successfully updates multiple reviews with reason
- Batch operations are atomic (all-or-nothing)
- Non-pending reviews cannot be included in batch operations

### 4. TestBatchModerationEmptyList

Tests edge case handling:
- Attempts to batch approve an empty list of reviews
- Attempts to batch reject an empty list of reviews
- Verifies both operations fail appropriately

**Key Assertions:**
- Empty batch operations are rejected
- System handles edge cases gracefully

### 5. TestRejectReviewWithoutReason

Tests validation requirements:
- Creates a pending review
- Attempts to reject the review without providing a reason
- Verifies the rejection fails
- Verifies the review status remains `pending`

**Key Assertions:**
- Rejection reason is required
- Reviews cannot be rejected without a reason
- Failed rejections don't change review status

## Test Infrastructure

### Helper Functions

#### `migrateReviewModerationModels(t *testing.T, db *gorm.DB)`
Migrates all necessary database models for review moderation tests.

#### `seedReviewModerationData(t *testing.T, db *gorm.DB) reviewModerationSeed`
Seeds the database with test data:
- Admin user
- Normal user
- Player user and player profile
- Game
- Service item
- Completed order

Returns a `reviewModerationSeed` struct with all created entity IDs.

#### `setupReviewModerationRouter(t *testing.T, db *gorm.DB, adminUserID uint64) (*gin.Engine, *adminservice.AdminService)`
Sets up the test router with:
- All necessary repositories
- AdminService with transaction manager
- Review moderation endpoints
- Fake authentication middleware

## Running the Tests

### Option 1: Using the test script (recommended)
```bash
cd backend
bash test_review_moderation.sh
```

### Option 2: Temporarily moving conflicting test files
```bash
cd backend/internal/integration
mv feed_integration_test.go feed_integration_test.go.bak
mv moderation_integration_test.go moderation_integration_test.go.bak
mv review_integration_test.go review_integration_test.go.bak
mv wallet_integration_test.go wallet_integration_test.go.bak
mv payment_refund_wallet_integration_test.go payment_refund_wallet_integration_test.go.bak

cd ../..
go test -v -run "TestReviewModeration|TestSensitiveWord|TestBatchModeration|TestRejectReview" ./internal/integration

# Restore files
cd internal/integration
mv feed_integration_test.go.bak feed_integration_test.go
mv moderation_integration_test.go.bak moderation_integration_test.go
mv review_integration_test.go.bak review_integration_test.go
mv wallet_integration_test.go.bak wallet_integration_test.go
mv payment_refund_wallet_integration_test.go.bak payment_refund_wallet_integration_test.go
```

## Test Results

All tests pass successfully:
```
=== RUN   TestReviewModerationFlow
--- PASS: TestReviewModerationFlow (0.03s)
=== RUN   TestSensitiveWordAutoMarking
--- PASS: TestSensitiveWordAutoMarking (0.01s)
=== RUN   TestBatchModeration
--- PASS: TestBatchModeration (0.01s)
=== RUN   TestBatchModerationEmptyList
--- PASS: TestBatchModerationEmptyList (0.01s)
=== RUN   TestRejectReviewWithoutReason
--- PASS: TestRejectReviewWithoutReason (0.01s)
PASS
ok      gamelink/internal/integration   0.678s
```

## Implementation Notes

1. **Database Setup**: Tests use an in-memory SQLite database for fast, isolated testing
2. **Transaction Management**: AdminService is configured with a transaction manager to ensure data consistency
3. **Authentication**: Tests use a fake authentication middleware to simulate admin user sessions
4. **Data Seeding**: Each test creates its own test data to ensure isolation
5. **Cleanup**: The testutil package handles database cleanup after each test

## Requirements Coverage

- ✅ **Requirement 2.1**: Complete review moderation flow tested
- ✅ **Requirement 2.2**: Review approval tested
- ✅ **Requirement 2.3**: Review rejection with reason tested
- ✅ **Requirement 2.4**: Sensitive word detection and marking tested
- ✅ **Requirement 2.5**: Batch approval and rejection tested, including atomicity

## Future Enhancements

Potential areas for additional test coverage:
1. Test operation log creation for moderation actions
2. Test notification sending to reviewers after moderation
3. Test concurrent moderation operations
4. Test moderation with different admin permission levels
5. Test sensitive word highlighting in the UI response
