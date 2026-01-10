package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/content"
	"gamelink/internal/repository/implementations"
	"gamelink/internal/repository/order"
	"gamelink/internal/repository/player"
	"gamelink/internal/repository/user"
	orderservice "gamelink/internal/service/order"

	"gorm.io/gorm"
)

// ============================================================================
// Review Service Integration Tests
// ============================================================================

// setupReviewService creates a new ReviewService with all required dependencies
func setupReviewService(db *gorm.DB) *orderservice.ReviewService {
	// The order repository functions as OrderReader
	return orderservice.NewReviewService(
		order.NewReviewRepository(db),
		implementations.NewOrderRepository(db), // OrderReader interface
		player.NewPlayerRepository(db),
		user.NewUserRepository(db),
		order.NewReviewReplyRepository(db),
		content.NewNotificationRepository(db),
	)
}

// ============================================================================
// CreateReview Tests
// ============================================================================

func TestReviewService_CreateReview_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	service := setupReviewService(db)
	ctx := context.Background()

	// Setup test data
	testUser := CreateUniqueTestUser(t, db, "review_user")
	playerUser := CreateUniqueTestUser(t, db, "review_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "review_game")
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, game, model.OrderStatusCompleted, 10000)

	// Test creating a review
	req := orderservice.CreateReviewRequest{
		OrderID:   testOrder.ID,
		Rating:    5,
		Comment:   "Excellent service! Very professional.",
		Anonymous: false,
	}

	resp, err := service.CreateReview(ctx, testUser.ID, req)
	require.NoError(t, err)
	assert.NotZero(t, resp.ReviewID)

	// Verify review was created
	reviewRepo := order.NewReviewRepository(db)
	createdReview, err := reviewRepo.Get(ctx, resp.ReviewID)
	require.NoError(t, err)
	assert.Equal(t, testOrder.ID, createdReview.OrderID)
	assert.Equal(t, testUser.ID, createdReview.UserID)
	assert.Equal(t, testPlayer.ID, createdReview.PlayerID)
	assert.Equal(t, model.Rating(5), createdReview.Score)
	assert.Equal(t, "Excellent service! Very professional.", createdReview.Content)
	assert.Equal(t, model.ReviewStatusApproved, createdReview.Status) // Auto-approved
}

func TestReviewService_CreateReview_OrderNotCompleted(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	service := setupReviewService(db)
	ctx := context.Background()

	testUser := CreateUniqueTestUser(t, db, "pending_user")
	playerUser := CreateUniqueTestUser(t, db, "pending_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "pending_game")
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, game, model.OrderStatusPending, 10000)

	req := orderservice.CreateReviewRequest{
		OrderID: testOrder.ID,
		Rating:  4,
		Comment: "Good service",
	}

	_, err := service.CreateReview(ctx, testUser.ID, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "订单未完成")
}

func TestReviewService_CreateReview_UnauthorizedUser(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	service := setupReviewService(db)
	ctx := context.Background()

	testUser := CreateUniqueTestUser(t, db, "order_owner")
	otherUser := CreateUniqueTestUser(t, db, "other_user")
	playerUser := CreateUniqueTestUser(t, db, "unauth_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "unauth_game")
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, game, model.OrderStatusCompleted, 10000)

	req := orderservice.CreateReviewRequest{
		OrderID: testOrder.ID,
		Rating:  5,
		Comment: "Trying to review someone else's order",
	}

	// Different user trying to review
	_, err := service.CreateReview(ctx, otherUser.ID, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "无权操作")
}

func TestReviewService_CreateReview_DuplicateReview(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	service := setupReviewService(db)
	ctx := context.Background()

	testUser := CreateUniqueTestUser(t, db, "duplicate_user")
	playerUser := CreateUniqueTestUser(t, db, "duplicate_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "duplicate_game")
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, game, model.OrderStatusCompleted, 10000)

	// Create first review
	firstReq := orderservice.CreateReviewRequest{
		OrderID: testOrder.ID,
		Rating:  5,
		Comment: "First review",
	}
	_, err := service.CreateReview(ctx, testUser.ID, firstReq)
	require.NoError(t, err)

	// Try to create second review for same order
	secondReq := orderservice.CreateReviewRequest{
		OrderID: testOrder.ID,
		Rating:  4,
		Comment: "Second review",
	}
	_, err = service.CreateReview(ctx, testUser.ID, secondReq)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "已评价")
}

func TestReviewService_CreateReview_InvalidRating(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	service := setupReviewService(db)
	ctx := context.Background()

	testUser := CreateUniqueTestUser(t, db, "invalid_rating_user")
	playerUser := CreateUniqueTestUser(t, db, "invalid_rating_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "invalid_rating_game")
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, game, model.OrderStatusCompleted, 10000)

	// Test rating below minimum (1)
	req := orderservice.CreateReviewRequest{
		OrderID: testOrder.ID,
		Rating:  0, // Invalid
		Comment: "Test",
	}
	_, err := service.CreateReview(ctx, testUser.ID, req)
	assert.Error(t, err)

	// Test rating above maximum (5)
	req.Rating = 6
	_, err = service.CreateReview(ctx, testUser.ID, req)
	assert.Error(t, err)
}

func TestReviewService_CreateReview_PlayerRatingUpdated(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	service := setupReviewService(db)
	ctx := context.Background()

	playerUser := CreateUniqueTestUser(t, db, "rated_player_user")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "rating_game")

	// Create multiple reviews for the same player
	for i := 0; i < 5; i++ {
		testUser := CreateUniqueTestUser(t, db, "user_"+string(rune('a'+i)))
		testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, game, model.OrderStatusCompleted, 10000)

		req := orderservice.CreateReviewRequest{
			OrderID: testOrder.ID,
			Rating:  5,
			Comment: "Great service",
		}
		_, err := service.CreateReview(ctx, testUser.ID, req)
		require.NoError(t, err)
	}

	// Verify player rating was updated
	playerRepo := player.NewPlayerRepository(db)
	updatedPlayer, err := playerRepo.Get(ctx, testPlayer.ID)
	require.NoError(t, err)
	assert.Equal(t, uint32(5), updatedPlayer.RatingCount)
	assert.InDelta(t, 5.0, float64(updatedPlayer.RatingAverage), 0.01)
}

// ============================================================================
// GetMyReviews Tests
// ============================================================================

func TestReviewService_GetMyReviews_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	service := setupReviewService(db)
	ctx := context.Background()

	testUser := CreateUniqueTestUser(t, db, "myreviews_user")
	playerUser := CreateUniqueTestUser(t, db, "myreviews_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "myreviews_game")

	// Create multiple reviews
	for i := 0; i < 3; i++ {
		testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, game, model.OrderStatusCompleted, 10000)
		req := orderservice.CreateReviewRequest{
			OrderID: testOrder.ID,
			Rating:  5 - i,
			Comment: "Review content",
		}
		_, err := service.CreateReview(ctx, testUser.ID, req)
		require.NoError(t, err)
	}

	// Get my reviews
	resp, err := service.GetMyReviews(ctx, testUser.ID, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(3), resp.Total)
	assert.Len(t, resp.Reviews, 3)
}

func TestReviewService_GetMyReviews_EmptyList(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	service := setupReviewService(db)
	ctx := context.Background()

	testUser := CreateUniqueTestUser(t, db, "noreviews_user")

	resp, err := service.GetMyReviews(ctx, testUser.ID, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(0), resp.Total)
	assert.Len(t, resp.Reviews, 0)
}

func TestReviewService_GetMyReviews_Pagination(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	service := setupReviewService(db)
	ctx := context.Background()

	testUser := CreateUniqueTestUser(t, db, "pagination_user")
	playerUser := CreateUniqueTestUser(t, db, "pagination_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "pagination_game")

	// Create 5 reviews
	for i := 0; i < 5; i++ {
		testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, game, model.OrderStatusCompleted, 10000)
		req := orderservice.CreateReviewRequest{
			OrderID: testOrder.ID,
			Rating:  5,
			Comment: "Review content",
		}
		_, err := service.CreateReview(ctx, testUser.ID, req)
		require.NoError(t, err)
	}

	// Get first page
	page1, err := service.GetMyReviews(ctx, testUser.ID, 1, 2)
	require.NoError(t, err)
	assert.Equal(t, int64(5), page1.Total)
	assert.Len(t, page1.Reviews, 2)

	// Get second page
	page2, err := service.GetMyReviews(ctx, testUser.ID, 2, 2)
	require.NoError(t, err)
	assert.Len(t, page2.Reviews, 2)
}

// ============================================================================
// GetPlayerReviews Tests
// ============================================================================

func TestReviewService_GetPlayerReviews_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	service := setupReviewService(db)
	ctx := context.Background()

	playerUser := CreateUniqueTestUser(t, db, "playerreviews_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "playerreviews_game")

	// Create multiple reviews from different users
	for i := 0; i < 3; i++ {
		testUser := CreateUniqueTestUser(t, db, "reviewer_"+string(rune('a'+i)))
		testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, game, model.OrderStatusCompleted, 10000)
		req := orderservice.CreateReviewRequest{
			OrderID: testOrder.ID,
			Rating:  5,
			Comment: "Great service",
		}
		_, err := service.CreateReview(ctx, testUser.ID, req)
		require.NoError(t, err)
	}

	// Get player reviews
	reviews, total, err := service.GetPlayerReviews(ctx, testPlayer.ID, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, reviews, 3)

	// Verify all reviews belong to the player
	for _, r := range reviews {
		assert.NotEmpty(t, r.UserNickname)
		assert.NotEmpty(t, r.CreatedAt)
	}
}

func TestReviewService_GetPlayerReviews_NoReviews(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	service := setupReviewService(db)
	ctx := context.Background()

	playerUser := CreateUniqueTestUser(t, db, "new_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	reviews, total, err := service.GetPlayerReviews(ctx, testPlayer.ID, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Len(t, reviews, 0)
}

// ============================================================================
// ReplyReview Tests
// ============================================================================

func TestReviewService_ReplyReview_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	service := setupReviewService(db)
	ctx := context.Background()

	testUser := CreateUniqueTestUser(t, db, "reply_user")
	playerUser := CreateUniqueTestUser(t, db, "reply_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "reply_game")
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, game, model.OrderStatusCompleted, 10000)

	// Create a review first
	req := orderservice.CreateReviewRequest{
		OrderID: testOrder.ID,
		Rating:  5,
		Comment: "Great service!",
	}
	resp, err := service.CreateReview(ctx, testUser.ID, req)
	require.NoError(t, err)

	// Player replies to review
	replyReq := orderservice.ReplyReviewRequest{
		Content: "Thank you for your feedback!",
	}
	replyResp, err := service.ReplyReview(ctx, playerUser.ID, resp.ReviewID, replyReq)
	require.NoError(t, err)
	assert.NotZero(t, replyResp.ReplyID)
	assert.NotEmpty(t, replyResp.Status)

	// Verify reply was created
	replyRepo := order.NewReviewReplyRepository(db)
	reply, err := replyRepo.Get(ctx, replyResp.ReplyID)
	require.NoError(t, err)
	assert.Equal(t, resp.ReviewID, reply.ReviewID)
	assert.Equal(t, playerUser.ID, reply.AuthorID)
	assert.Equal(t, "Thank you for your feedback!", reply.Content)
}

func TestReviewService_ReplyReview_UnauthorizedPlayer(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	service := setupReviewService(db)
	ctx := context.Background()

	testUser := CreateUniqueTestUser(t, db, "auth_review_user")
	playerUser := CreateUniqueTestUser(t, db, "auth_player")
	otherPlayerUser := CreateUniqueTestUser(t, db, "other_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	_ = CreateTestPlayer(t, db, otherPlayerUser) // Create other player but don't use
	game := CreateTestGame(t, db, "auth_game")
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, game, model.OrderStatusCompleted, 10000)

	// Create a review
	req := orderservice.CreateReviewRequest{
		OrderID: testOrder.ID,
		Rating:  5,
		Comment: "Great service!",
	}
	resp, err := service.CreateReview(ctx, testUser.ID, req)
	require.NoError(t, err)

	// Different player tries to reply
	replyReq := orderservice.ReplyReviewRequest{
		Content: "This is not my review to reply to",
	}
	_, err = service.ReplyReview(ctx, otherPlayerUser.ID, resp.ReviewID, replyReq)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "无权操作")
}

func TestReviewService_ReplyReview_EmptyContent(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	service := setupReviewService(db)
	ctx := context.Background()

	testUser := CreateUniqueTestUser(t, db, "empty_reply_user")
	playerUser := CreateUniqueTestUser(t, db, "empty_reply_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "empty_reply_game")
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, game, model.OrderStatusCompleted, 10000)

	// Create a review
	req := orderservice.CreateReviewRequest{
		OrderID: testOrder.ID,
		Rating:  5,
		Comment: "Great service!",
	}
	resp, err := service.CreateReview(ctx, testUser.ID, req)
	require.NoError(t, err)

	// Try to reply with empty content
	replyReq := orderservice.ReplyReviewRequest{
		Content: "   ", // Only whitespace
	}
	_, err = service.ReplyReview(ctx, playerUser.ID, resp.ReviewID, replyReq)
	assert.Error(t, err)
}

func TestReviewService_ReplyReview_TooLongContent(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	service := setupReviewService(db)
	ctx := context.Background()

	testUser := CreateUniqueTestUser(t, db, "long_reply_user")
	playerUser := CreateUniqueTestUser(t, db, "long_reply_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "long_reply_game")
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, game, model.OrderStatusCompleted, 10000)

	// Create a review
	req := orderservice.CreateReviewRequest{
		OrderID: testOrder.ID,
		Rating:  5,
		Comment: "Great service!",
	}
	resp, err := service.CreateReview(ctx, testUser.ID, req)
	require.NoError(t, err)

	// Try to reply with content exceeding 500 characters
	longContent := ""
	for i := 0; i < 501; i++ {
		longContent += "a"
	}
	replyReq := orderservice.ReplyReviewRequest{
		Content: longContent,
	}
	_, err = service.ReplyReview(ctx, playerUser.ID, resp.ReviewID, replyReq)
	assert.Error(t, err)
}

// ============================================================================
// UpdateReply Tests
// ============================================================================

func TestReviewService_UpdateReply_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	service := setupReviewService(db)
	ctx := context.Background()

	testUser := CreateUniqueTestUser(t, db, "update_reply_user")
	playerUser := CreateUniqueTestUser(t, db, "update_reply_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "update_reply_game")
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, game, model.OrderStatusCompleted, 10000)

	// Create review and reply
	req := orderservice.CreateReviewRequest{
		OrderID: testOrder.ID,
		Rating:  5,
		Comment: "Great service!",
	}
	resp, err := service.CreateReview(ctx, testUser.ID, req)
	require.NoError(t, err)

	replyReq := orderservice.ReplyReviewRequest{
		Content: "Initial reply",
	}
	replyResp, err := service.ReplyReview(ctx, playerUser.ID, resp.ReviewID, replyReq)
	require.NoError(t, err)

	// Update reply
	updateReq := orderservice.UpdateReplyRequest{
		Content: "Updated reply content",
	}
	updateResp, err := service.UpdateReply(ctx, playerUser.ID, replyResp.ReplyID, updateReq)
	require.NoError(t, err)
	assert.Equal(t, replyResp.ReplyID, updateResp.ReplyID)

	// Verify update
	replyRepo := order.NewReviewReplyRepository(db)
	updatedReply, err := replyRepo.Get(ctx, replyResp.ReplyID)
	require.NoError(t, err)
	assert.Equal(t, "Updated reply content", updatedReply.Content)
}

func TestReviewService_UpdateReply_Unauthorized(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	service := setupReviewService(db)
	ctx := context.Background()

	testUser := CreateUniqueTestUser(t, db, "unauth_update_user")
	playerUser := CreateUniqueTestUser(t, db, "unauth_update_player")
	otherPlayerUser := CreateUniqueTestUser(t, db, "other_update_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	_ = CreateTestPlayer(t, db, otherPlayerUser) // Create other player but don't use
	game := CreateTestGame(t, db, "unauth_update_game")
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, game, model.OrderStatusCompleted, 10000)

	// Create review and reply
	req := orderservice.CreateReviewRequest{
		OrderID: testOrder.ID,
		Rating:  5,
		Comment: "Great service!",
	}
	resp, err := service.CreateReview(ctx, testUser.ID, req)
	require.NoError(t, err)

	replyReq := orderservice.ReplyReviewRequest{
		Content: "Original reply",
	}
	replyResp, err := service.ReplyReview(ctx, playerUser.ID, resp.ReviewID, replyReq)
	require.NoError(t, err)

	// Different player tries to update
	updateReq := orderservice.UpdateReplyRequest{
		Content: "Trying to update someone else's reply",
	}
	_, err = service.UpdateReply(ctx, otherPlayerUser.ID, replyResp.ReplyID, updateReq)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "无权操作")
}

// ============================================================================
// DeleteReply Tests
// ============================================================================

func TestReviewService_DeleteReply_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	service := setupReviewService(db)
	ctx := context.Background()

	testUser := CreateUniqueTestUser(t, db, "delete_reply_user")
	playerUser := CreateUniqueTestUser(t, db, "delete_reply_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "delete_reply_game")
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, game, model.OrderStatusCompleted, 10000)

	// Create review and reply
	req := orderservice.CreateReviewRequest{
		OrderID: testOrder.ID,
		Rating:  5,
		Comment: "Great service!",
	}
	resp, err := service.CreateReview(ctx, testUser.ID, req)
	require.NoError(t, err)

	replyReq := orderservice.ReplyReviewRequest{
		Content: "Reply to be deleted",
	}
	replyResp, err := service.ReplyReview(ctx, playerUser.ID, resp.ReviewID, replyReq)
	require.NoError(t, err)

	// Delete reply
	err = service.DeleteReply(ctx, playerUser.ID, replyResp.ReplyID)
	require.NoError(t, err)

	// Verify deletion
	replyRepo := order.NewReviewReplyRepository(db)
	_, err = replyRepo.Get(ctx, replyResp.ReplyID)
	assert.Error(t, err)
}

func TestReviewService_DeleteReply_Unauthorized(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	service := setupReviewService(db)
	ctx := context.Background()

	testUser := CreateUniqueTestUser(t, db, "unauth_delete_user")
	playerUser := CreateUniqueTestUser(t, db, "unauth_delete_player")
	otherPlayerUser := CreateUniqueTestUser(t, db, "other_delete_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	_ = CreateTestPlayer(t, db, otherPlayerUser) // Create other player but don't use
	game := CreateTestGame(t, db, "unauth_delete_game")
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, game, model.OrderStatusCompleted, 10000)

	// Create review and reply
	req := orderservice.CreateReviewRequest{
		OrderID: testOrder.ID,
		Rating:  5,
		Comment: "Great service!",
	}
	resp, err := service.CreateReview(ctx, testUser.ID, req)
	require.NoError(t, err)

	replyReq := orderservice.ReplyReviewRequest{
		Content: "Reply to be deleted",
	}
	replyResp, err := service.ReplyReview(ctx, playerUser.ID, resp.ReviewID, replyReq)
	require.NoError(t, err)

	// Different player tries to delete
	err = service.DeleteReply(ctx, otherPlayerUser.ID, replyResp.ReplyID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "无权操作")
}

// ============================================================================
// Edge Cases and Business Rules Tests
// ============================================================================

func TestReviewService_AverageRatingCalculation(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	service := setupReviewService(db)
	ctx := context.Background()

	playerUser := CreateUniqueTestUser(t, db, "avg_player_user")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "avg_game")

	// Create reviews with different ratings: 5, 4, 3, 2, 1 (average = 3)
	ratings := []int{5, 4, 3, 2, 1}
	for i, rating := range ratings {
		testUser := CreateUniqueTestUser(t, db, "avg_user_"+string(rune('a'+i)))
		testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, game, model.OrderStatusCompleted, 10000)
		req := orderservice.CreateReviewRequest{
			OrderID: testOrder.ID,
			Rating:  rating,
			Comment: "Review",
		}
		_, err := service.CreateReview(ctx, testUser.ID, req)
		require.NoError(t, err)
	}

	// Verify player average rating
	playerRepo := player.NewPlayerRepository(db)
	updatedPlayer, err := playerRepo.Get(ctx, testPlayer.ID)
	require.NoError(t, err)
	assert.Equal(t, uint32(5), updatedPlayer.RatingCount)
	assert.InDelta(t, 3.0, float64(updatedPlayer.RatingAverage), 0.01)
}

func TestReviewService_MultipleReviewsFromDifferentUsers(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	service := setupReviewService(db)
	ctx := context.Background()

	playerUser := CreateUniqueTestUser(t, db, "multi_review_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "multi_review_game")

	// Create multiple users reviewing the same player
	for i := 0; i < 5; i++ {
		user := CreateUniqueTestUser(t, db, "reviewer_"+string(rune('a'+i)))
		order := CreateTestOrderWithDetails(t, db, user, testPlayer, game, model.OrderStatusCompleted, 10000)

		req := orderservice.CreateReviewRequest{
			OrderID: order.ID,
			Rating:  4 + i%2, // Alternate between 4 and 5
			Comment: "Review content",
		}
		_, err := service.CreateReview(ctx, user.ID, req)
		require.NoError(t, err)
	}

	// Get player reviews
	reviews, total, err := service.GetPlayerReviews(ctx, testPlayer.ID, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, reviews, 5)
}

func TestReviewService_PlayerRatingAverageEdgeCases(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	service := setupReviewService(db)
	ctx := context.Background()

	// Test case 1: Only 1-star reviews
	playerUser1 := CreateUniqueTestUser(t, db, "low_rated_player")
	testPlayer1 := CreateTestPlayer(t, db, playerUser1)
	game1 := CreateTestGame(t, db, "low_rated_game")

	for i := 0; i < 3; i++ {
		user := CreateUniqueTestUser(t, db, "low_user_"+string(rune('a'+i)))
		order := CreateTestOrderWithDetails(t, db, user, testPlayer1, game1, model.OrderStatusCompleted, 10000)
		req := orderservice.CreateReviewRequest{
			OrderID: order.ID,
			Rating:  1,
			Comment: "Poor service",
		}
		_, err := service.CreateReview(ctx, user.ID, req)
		require.NoError(t, err)
	}

	playerRepo := player.NewPlayerRepository(db)
	updatedPlayer1, err := playerRepo.Get(ctx, testPlayer1.ID)
	require.NoError(t, err)
	assert.InDelta(t, 1.0, float64(updatedPlayer1.RatingAverage), 0.01)

	// Test case 2: Only 5-star reviews
	playerUser2 := CreateUniqueTestUser(t, db, "high_rated_player")
	testPlayer2 := CreateTestPlayer(t, db, playerUser2)
	game2 := CreateTestGame(t, db, "high_rated_game")

	for i := 0; i < 3; i++ {
		user := CreateUniqueTestUser(t, db, "high_user_"+string(rune('a'+i)))
		order := CreateTestOrderWithDetails(t, db, user, testPlayer2, game2, model.OrderStatusCompleted, 10000)
		req := orderservice.CreateReviewRequest{
			OrderID: order.ID,
			Rating:  5,
			Comment: "Excellent service",
		}
		_, err := service.CreateReview(ctx, user.ID, req)
		require.NoError(t, err)
	}

	updatedPlayer2, err := playerRepo.Get(ctx, testPlayer2.ID)
	require.NoError(t, err)
	assert.InDelta(t, 5.0, float64(updatedPlayer2.RatingAverage), 0.01)
}

func TestReviewService_ReviewDoesNotExist(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	service := setupReviewService(db)
	ctx := context.Background()

	playerUser := CreateUniqueTestUser(t, db, "nonexist_player")
	_ = CreateTestPlayer(t, db, playerUser)

	// Try to reply to non-existent review
	replyReq := orderservice.ReplyReviewRequest{
		Content: "Reply to non-existent review",
	}
	_, err := service.ReplyReview(ctx, playerUser.ID, 999999, replyReq)
	assert.Error(t, err)
}

// ============================================================================
// Multi-Player Order Tests
// ============================================================================

func TestReviewService_MultiPlayerOrder_AllPlayersReviewed(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	service := setupReviewService(db)
	ctx := context.Background()

	testUser := CreateUniqueTestUser(t, db, "multiplayer_user")
	playerUser1 := CreateUniqueTestUser(t, db, "multiplayer_player1")
	playerUser2 := CreateUniqueTestUser(t, db, "multiplayer_player2")
	testPlayer1 := CreateTestPlayer(t, db, playerUser1)
	testPlayer2 := CreateTestPlayer(t, db, playerUser2)
	game := CreateTestGame(t, db, "multiplayer_game")

	// Create multi-player order (team order)
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer1, game, model.OrderStatusCompleted, 20000)

	// Create order items for each player
	orderItem1 := CreateTestOrderItem(t, db, testOrder, 1, 10000, model.OrderItemStatusCompleted)
	orderItem2 := CreateTestOrderItem(t, db, testOrder, 2, 10000, model.OrderItemStatusCompleted)

	// Add players to order
	CreateTestOrderPlayer(t, db, testOrder, orderItem1, testPlayer1, model.OrderPlayerStatusCompleted)
	CreateTestOrderPlayer(t, db, testOrder, orderItem2, testPlayer2, model.OrderPlayerStatusCompleted)

	// Create review for player 1
	req1 := orderservice.CreateReviewRequest{
		OrderID: testOrder.ID,
		Rating:  5,
		Comment: "Player 1 was great",
	}
	resp1, err := service.CreateReview(ctx, testUser.ID, req1)
	require.NoError(t, err)
	assert.NotZero(t, resp1.ReviewID)

	// Create review for player 2
	req2 := orderservice.CreateReviewRequest{
		OrderID: testOrder.ID,
		Rating:  4,
		Comment: "Player 2 was good",
	}
	resp2, err := service.CreateReview(ctx, testUser.ID, req2)
	require.NoError(t, err)
	assert.NotZero(t, resp2.ReviewID)

	// Verify both reviews were created
	reviewRepo := order.NewReviewRepository(db)
	reviews, total, err := reviewRepo.List(ctx, repository.ReviewListOptions{
		OrderID:  &testOrder.ID,
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(2))
	assert.GreaterOrEqual(t, len(reviews), 2)
}
