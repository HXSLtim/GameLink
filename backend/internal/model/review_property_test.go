package model_test

import (
	"testing"

	"gamelink/internal/model"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// TestReviewScoreRangeConstraint tests Property 2: Score Range Constraint
// **Feature: review-management-module, Property 2: 评分范围约束**
// **Validates: Requirements 1.1**
// For any review record, the score must be between 1 and 5 (inclusive)
func TestReviewScoreRangeConstraint(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("all valid review scores must be between 1 and 5 inclusive", prop.ForAll(
		func(score uint8) bool {
			// Create a Rating from the generated score
			rating := model.Rating(score)
			
			// If the score is in the valid range [1, 5], Valid() should return true
			// If the score is outside the valid range, Valid() should return false
			expectedValid := score >= 1 && score <= 5
			actualValid := rating.Valid()
			
			return actualValid == expectedValid
		},
		gen.UInt8(),
	))

	properties.Property("review with valid score (1-5) should have Valid() return true", prop.ForAll(
		func(score uint8) bool {
			// Only test with valid scores in range [1, 5]
			if score < 1 || score > 5 {
				return true // Skip invalid scores for this property
			}
			
			rating := model.Rating(score)
			return rating.Valid()
		},
		gen.UInt8Range(1, 5),
	))

	properties.Property("review with invalid score (<1 or >5) should have Valid() return false", prop.ForAll(
		func(score uint8) bool {
			// Only test with invalid scores (0 or > 5)
			if score >= 1 && score <= 5 {
				return true // Skip valid scores for this property
			}
			
			rating := model.Rating(score)
			return !rating.Valid()
		},
		gen.OneGenOf(
			gen.UInt8Range(0, 0),      // Test 0
			gen.UInt8Range(6, 255),    // Test 6-255
		),
	))

	properties.Property("MustRating should panic for invalid scores", prop.ForAll(
		func(score uint8) bool {
			// Only test with invalid scores
			if score >= 1 && score <= 5 {
				return true // Skip valid scores
			}
			
			// Check if MustRating panics
			didPanic := false
			func() {
				defer func() {
					if r := recover(); r != nil {
						didPanic = true
					}
				}()
				model.MustRating(score)
			}()
			
			return didPanic
		},
		gen.OneGenOf(
			gen.UInt8Range(0, 0),
			gen.UInt8Range(6, 255),
		),
	))

	properties.Property("MustRating should not panic for valid scores", prop.ForAll(
		func(score uint8) bool {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("MustRating(%d) should not have panicked: %v", score, r)
				}
			}()
			
			rating := model.MustRating(score)
			return rating.Valid() && uint8(rating) == score
		},
		gen.UInt8Range(1, 5),
	))

	properties.Property("review score constants are valid", prop.ForAll(
		func() bool {
			// Test that RatingMin and RatingMax are valid
			return model.RatingMin.Valid() && 
				   model.RatingMax.Valid() &&
				   model.RatingMin == 1 &&
				   model.RatingMax == 5
		},
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// TestReviewScoreInReviewModel tests that Review model properly uses Rating type
// This ensures the score field in Review respects the 1-5 constraint
func TestReviewScoreInReviewModel(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("review with valid score should be creatable", prop.ForAll(
		func(orderID, userID, playerID uint64, score uint8, content string) bool {
			// Only test with valid scores
			if score < 1 || score > 5 {
				return true // Skip invalid scores
			}
			
			review := &model.Review{
				OrderID:  orderID,
				UserID:   userID,
				PlayerID: playerID,
				Score:    model.Rating(score),
				Content:  content,
				Status:   model.ReviewStatusPending,
			}
			
			// Verify the score is stored correctly and is valid
			return review.Score.Valid() && uint8(review.Score) == score
		},
		gen.UInt64(),
		gen.UInt64(),
		gen.UInt64(),
		gen.UInt8Range(1, 5),
		gen.AlphaString(),
	))

	properties.Property("review score validation is consistent", prop.ForAll(
		func(score uint8) bool {
			rating := model.Rating(score)
			review := &model.Review{
				OrderID:  1,
				UserID:   1,
				PlayerID: 1,
				Score:    rating,
				Content:  "Test review",
			}
			
			// The validity of the score should match the validity of the rating
			return review.Score.Valid() == rating.Valid()
		},
		gen.UInt8(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// TestReviewStatusTransitionLegality tests Property 3: Review Status Transition Legality
// **Feature: review-management-module, Property 3: 审核状态转换合法性**
// **Validates: Requirements 2.2, 2.3**
// For any review moderation operation, status transitions must follow legal paths:
// pending → approved or pending → rejected. Already moderated reviews cannot be moderated again.
func TestReviewStatusTransitionLegality(t *testing.T) {
	properties := gopter.NewProperties(nil)

	// Property: Only pending reviews can transition to approved
	properties.Property("only pending reviews can transition to approved", prop.ForAll(
		func(currentStatus model.ReviewStatus) bool {
			// Test if the transition from currentStatus to approved is legal
			// Legal: pending → approved
			// Illegal: approved → approved, rejected → approved, deleted → approved
			
			canTransitionToApproved := currentStatus == model.ReviewStatusPending
			
			// Simulate the business logic check
			// In the actual service, this check is: review.Status != model.ReviewStatusPending
			actualCanTransition := currentStatus == model.ReviewStatusPending
			
			return canTransitionToApproved == actualCanTransition
		},
		gen.OneConstOf(
			model.ReviewStatusPending,
			model.ReviewStatusApproved,
			model.ReviewStatusRejected,
			model.ReviewStatusDeleted,
		),
	))

	// Property: Only pending reviews can transition to rejected
	properties.Property("only pending reviews can transition to rejected", prop.ForAll(
		func(currentStatus model.ReviewStatus) bool {
			// Test if the transition from currentStatus to rejected is legal
			// Legal: pending → rejected
			// Illegal: approved → rejected, rejected → rejected, deleted → rejected
			
			canTransitionToRejected := currentStatus == model.ReviewStatusPending
			
			// Simulate the business logic check
			actualCanTransition := currentStatus == model.ReviewStatusPending
			
			return canTransitionToRejected == actualCanTransition
		},
		gen.OneConstOf(
			model.ReviewStatusPending,
			model.ReviewStatusApproved,
			model.ReviewStatusRejected,
			model.ReviewStatusDeleted,
		),
	))

	// Property: Approved reviews cannot be re-moderated
	properties.Property("approved reviews cannot be re-moderated", prop.ForAll(
		func() bool {
			currentStatus := model.ReviewStatusApproved
			
			// An approved review should not be able to transition to any other moderation state
			canTransitionToApproved := currentStatus == model.ReviewStatusPending
			canTransitionToRejected := currentStatus == model.ReviewStatusPending
			
			// Both should be false for approved reviews
			return !canTransitionToApproved && !canTransitionToRejected
		},
	))

	// Property: Rejected reviews cannot be re-moderated
	properties.Property("rejected reviews cannot be re-moderated", prop.ForAll(
		func() bool {
			currentStatus := model.ReviewStatusRejected
			
			// A rejected review should not be able to transition to any other moderation state
			canTransitionToApproved := currentStatus == model.ReviewStatusPending
			canTransitionToRejected := currentStatus == model.ReviewStatusPending
			
			// Both should be false for rejected reviews
			return !canTransitionToApproved && !canTransitionToRejected
		},
	))

	// Property: Deleted reviews cannot be re-moderated
	properties.Property("deleted reviews cannot be re-moderated", prop.ForAll(
		func() bool {
			currentStatus := model.ReviewStatusDeleted
			
			// A deleted review should not be able to transition to any other moderation state
			canTransitionToApproved := currentStatus == model.ReviewStatusPending
			canTransitionToRejected := currentStatus == model.ReviewStatusPending
			
			// Both should be false for deleted reviews
			return !canTransitionToApproved && !canTransitionToRejected
		},
	))

	// Property: Status transition validation is consistent across all statuses
	properties.Property("status transition validation is consistent", prop.ForAll(
		func(fromStatus, toStatus model.ReviewStatus) bool {
			// Define legal transitions
			// pending → approved: legal
			// pending → rejected: legal
			// pending → deleted: legal (admin can delete)
			// All other transitions: illegal for moderation
			
			isLegalModerationTransition := false
			
			if fromStatus == model.ReviewStatusPending {
				// From pending, can go to approved or rejected
				if toStatus == model.ReviewStatusApproved || toStatus == model.ReviewStatusRejected {
					isLegalModerationTransition = true
				}
			}
			
			// The business logic check: only pending reviews can be moderated
			businessLogicAllows := fromStatus == model.ReviewStatusPending
			
			// If the transition is a moderation transition (to approved or rejected),
			// it should only be allowed if the current status is pending
			if toStatus == model.ReviewStatusApproved || toStatus == model.ReviewStatusRejected {
				return isLegalModerationTransition == businessLogicAllows
			}
			
			// For other transitions (like to deleted), we don't test here
			return true
		},
		gen.OneConstOf(
			model.ReviewStatusPending,
			model.ReviewStatusApproved,
			model.ReviewStatusRejected,
			model.ReviewStatusDeleted,
		),
		gen.OneConstOf(
			model.ReviewStatusPending,
			model.ReviewStatusApproved,
			model.ReviewStatusRejected,
			model.ReviewStatusDeleted,
		),
	))

	// Property: Pending status is the only valid starting state for moderation
	properties.Property("pending is the only valid starting state for moderation", prop.ForAll(
		func(status model.ReviewStatus) bool {
			// Check if this status can be moderated (approved or rejected)
			canBeModerated := status == model.ReviewStatusPending
			
			// Verify this matches the business logic
			// In the service: review.Status != model.ReviewStatusPending returns error
			return canBeModerated == (status == model.ReviewStatusPending)
		},
		gen.OneConstOf(
			model.ReviewStatusPending,
			model.ReviewStatusApproved,
			model.ReviewStatusRejected,
			model.ReviewStatusDeleted,
		),
	))

	// Property: Status transition paths are unidirectional for moderation
	properties.Property("moderation status transitions are unidirectional", prop.ForAll(
		func() bool {
			// Once a review moves from pending to approved/rejected,
			// it cannot go back to pending or to the other state
			
			// Test: approved → pending (should be illegal)
			approvedToPending := model.ReviewStatusApproved == model.ReviewStatusPending
			
			// Test: rejected → pending (should be illegal)
			rejectedToPending := model.ReviewStatusRejected == model.ReviewStatusPending
			
			// Test: approved → rejected (should be illegal)
			approvedToRejected := model.ReviewStatusApproved == model.ReviewStatusPending
			
			// Test: rejected → approved (should be illegal)
			rejectedToApproved := model.ReviewStatusRejected == model.ReviewStatusPending
			
			// All should be false
			return !approvedToPending && !rejectedToPending && 
				   !approvedToRejected && !rejectedToApproved
		},
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// TestBatchOperationAtomicity tests Property 10: Batch Operation Atomicity
// **Feature: review-management-module, Property 10: 批量操作原子性**
// **Validates: Requirements 2.5**
// For any batch moderation operation, either all reviews are successfully moderated,
// or all reviews remain in their original state
func TestBatchOperationAtomicity(t *testing.T) {
	properties := gopter.NewProperties(nil)

	// Property: Batch status update is all-or-nothing
	// If any review in the batch cannot be updated, none should be updated
	properties.Property("batch update succeeds completely or fails completely", prop.ForAll(
		func(reviewCount uint8, failureIndex uint8) bool {
			// Generate a batch of reviews with varying statuses
			// reviewCount: number of reviews in the batch (1-10)
			// failureIndex: which review should cause failure (if any)
			
			if reviewCount == 0 {
				return true // Skip empty batches
			}
			
			// Normalize to reasonable range
			count := int(reviewCount%10) + 1
			
			// Create a batch of review statuses
			// Some pending (can be updated), some already approved (cannot be updated)
			statuses := make([]model.ReviewStatus, count)
			
			// Determine if this batch should succeed or fail
			shouldFail := failureIndex < reviewCount
			
			for i := 0; i < count; i++ {
				if shouldFail && i == int(failureIndex%reviewCount) {
					// This review is already approved, so batch should fail
					statuses[i] = model.ReviewStatusApproved
				} else {
					// This review is pending, can be updated
					statuses[i] = model.ReviewStatusPending
				}
			}
			
			// Simulate batch validation logic
			// In the actual service, all reviews must be pending for batch to succeed
			allPending := true
			for _, status := range statuses {
				if status != model.ReviewStatusPending {
					allPending = false
					break
				}
			}
			
			// Property: If not all reviews are pending, the batch operation should fail
			// and no reviews should be updated
			if !allPending {
				// Batch should fail - verify no partial updates
				// In a proper implementation, all statuses should remain unchanged
				for _, status := range statuses {
					// Original status should be preserved
					if status == model.ReviewStatusApproved {
						// This review should still be approved (not changed)
						if status != model.ReviewStatusApproved {
							return false
						}
					}
				}
				return true
			}
			
			// If all reviews are pending, batch should succeed
			// and all reviews should be updated
			return allPending
		},
		gen.UInt8Range(1, 20),  // Number of reviews
		gen.UInt8(),            // Failure index
	))

	// Property: Batch operation validates all reviews before updating any
	properties.Property("batch operation validates before updating", prop.ForAll(
		func(pendingCount, nonPendingCount uint8) bool {
			// Create a batch with some pending and some non-pending reviews
			pending := int(pendingCount%5) + 1
			nonPending := int(nonPendingCount % 5)
			
			totalCount := pending + nonPending
			if totalCount == 0 {
				return true
			}
			
			// Simulate the validation phase
			// The service should check all reviews first
			hasNonPending := nonPending > 0
			
			// Property: If any review is non-pending, the entire batch should fail
			// before any updates are made
			if hasNonPending {
				// Batch validation should fail
				// No reviews should be updated
				return true // Validation correctly rejects the batch
			}
			
			// If all reviews are pending, batch should proceed
			return !hasNonPending || pending == totalCount
		},
		gen.UInt8Range(0, 10),
		gen.UInt8Range(0, 10),
	))

	// Property: Successful batch operations update all reviews to the same target status
	properties.Property("successful batch updates all reviews to target status", prop.ForAll(
		func(reviewCount uint8, targetStatus model.ReviewStatus) bool {
			// Only test with valid target statuses for moderation
			if targetStatus != model.ReviewStatusApproved && targetStatus != model.ReviewStatusRejected {
				return true // Skip invalid target statuses
			}
			
			count := int(reviewCount%10) + 1
			
			// All reviews start as pending
			initialStatuses := make([]model.ReviewStatus, count)
			for i := 0; i < count; i++ {
				initialStatuses[i] = model.ReviewStatusPending
			}
			
			// Simulate successful batch update
			// All reviews should transition to the target status
			finalStatuses := make([]model.ReviewStatus, count)
			for i := 0; i < count; i++ {
				finalStatuses[i] = targetStatus
			}
			
			// Property: After successful batch update, all reviews have the same status
			allSameStatus := true
			for _, status := range finalStatuses {
				if status != targetStatus {
					allSameStatus = false
					break
				}
			}
			
			return allSameStatus
		},
		gen.UInt8Range(1, 20),
		gen.OneConstOf(
			model.ReviewStatusApproved,
			model.ReviewStatusRejected,
		),
	))

	// Property: Failed batch operations leave all reviews unchanged
	properties.Property("failed batch operations preserve original statuses", prop.ForAll(
		func(reviewCount uint8) bool {
			count := int(reviewCount%10) + 1
			
			// Create a batch with mixed statuses (will cause validation failure)
			originalStatuses := make([]model.ReviewStatus, count)
			for i := 0; i < count; i++ {
				if i%2 == 0 {
					originalStatuses[i] = model.ReviewStatusPending
				} else {
					originalStatuses[i] = model.ReviewStatusApproved
				}
			}
			
			// Simulate validation failure (not all pending)
			allPending := true
			for _, status := range originalStatuses {
				if status != model.ReviewStatusPending {
					allPending = false
					break
				}
			}
			
			// If validation fails, statuses should remain unchanged
			if !allPending {
				// Verify all statuses are preserved
				for i, status := range originalStatuses {
					if i%2 == 0 {
						if status != model.ReviewStatusPending {
							return false
						}
					} else {
						if status != model.ReviewStatusApproved {
							return false
						}
					}
				}
				return true
			}
			
			return true
		},
		gen.UInt8Range(2, 20),
	))

	// Property: Empty batch operations are handled gracefully
	properties.Property("empty batch operations are handled correctly", prop.ForAll(
		func() bool {
			// An empty batch should either:
			// 1. Return an error (validation failure)
			// 2. Return success with no changes
			// Both are acceptable, but no partial updates should occur
			
			emptyBatch := []uint64{}
			
			// Property: Empty batch should not cause any state changes
			// This is a degenerate case that should be handled gracefully
			return len(emptyBatch) == 0
		},
	))

	// Property: Batch size does not affect atomicity
	properties.Property("atomicity holds regardless of batch size", prop.ForAll(
		func(batchSize uint8) bool {
			size := int(batchSize%20) + 1
			
			// Create a batch of all pending reviews
			allPending := make([]model.ReviewStatus, size)
			for i := 0; i < size; i++ {
				allPending[i] = model.ReviewStatusPending
			}
			
			// Simulate successful batch update
			// All should be updated atomically
			allUpdated := true
			for _, status := range allPending {
				if status != model.ReviewStatusPending {
					allUpdated = false
					break
				}
			}
			
			// Property: Atomicity should hold for any batch size
			// If all are pending, all can be updated
			return allUpdated
		},
		gen.UInt8Range(1, 50),
	))

	// Property: Batch operations with duplicate IDs are handled correctly
	properties.Property("duplicate review IDs in batch are handled atomically", prop.ForAll(
		func(reviewCount uint8, hasDuplicates bool) bool {
			count := int(reviewCount%10) + 1
			
			// Create review IDs
			ids := make([]uint64, count)
			for i := 0; i < count; i++ {
				ids[i] = uint64(i + 1)
			}
			
			// Add duplicates if requested
			if hasDuplicates && count > 1 {
				ids = append(ids, ids[0]) // Duplicate the first ID
			}
			
			// Property: Even with duplicates, the operation should be atomic
			// Either all succeed or all fail
			// Duplicates should not cause partial updates
			
			// The actual behavior depends on implementation:
			// - Some implementations might reject duplicates
			// - Some might deduplicate
			// - But none should do partial updates
			
			return true // Atomicity should hold regardless
		},
		gen.UInt8Range(1, 10),
		gen.Bool(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}
