package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/review"
)

// ============================================================================
// Review CRUD Tests
// ============================================================================

func TestReviewRepository_Create(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := review.NewReviewRepository(db)
	ctx := context.Background()

	// Create user and player
	user := CreateUniqueTestUser(t, db, "review_user")
	playerUser := CreateUniqueTestUser(t, db, "review_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "review_game")

	// Create order with details
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, 10000)

	reviewObj := &model.Review{
		OrderID:  order.ID,
		UserID:   user.ID,
		PlayerID: player.ID,
		Score:    5,
		Content:  "Great service!",
		Status:   model.ReviewStatusPending,
	}

	err := repo.Create(ctx, reviewObj)
	require.NoError(t, err)
	assert.NotZero(t, reviewObj.ID)
}

func TestReviewRepository_Get(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := review.NewReviewRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "get_review_user")
	playerUser := CreateUniqueTestUser(t, db, "get_review_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "get_review_game")
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, 10000)

	reviewObj := &model.Review{
		OrderID:  order.ID,
		UserID:   user.ID,
		PlayerID: player.ID,
		Score:    4,
		Content:  "Good service",
		Status:   model.ReviewStatusPending,
	}
	require.NoError(t, repo.Create(ctx, reviewObj))

	got, err := repo.Get(ctx, reviewObj.ID)
	require.NoError(t, err)
	assert.Equal(t, model.Rating(4), got.Score)
	assert.Equal(t, "Good service", got.Content)
}

func TestReviewRepository_Update(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := review.NewReviewRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "update_review_user")
	playerUser := CreateUniqueTestUser(t, db, "update_review_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "update_review_game")
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, 10000)

	reviewObj := &model.Review{
		OrderID:  order.ID,
		UserID:   user.ID,
		PlayerID: player.ID,
		Score:    3,
		Content:  "Average",
		Status:   model.ReviewStatusPending,
	}
	require.NoError(t, repo.Create(ctx, reviewObj))

	// Update
	reviewObj.Score = 5
	reviewObj.Content = "Actually great!"
	err := repo.Update(ctx, reviewObj)
	require.NoError(t, err)

	got, err := repo.Get(ctx, reviewObj.ID)
	require.NoError(t, err)
	assert.Equal(t, model.Rating(5), got.Score)
	assert.Equal(t, "Actually great!", got.Content)
}

func TestReviewRepository_Delete(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := review.NewReviewRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "delete_review_user")
	playerUser := CreateUniqueTestUser(t, db, "delete_review_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "delete_review_game")
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, 10000)

	reviewObj := &model.Review{
		OrderID:  order.ID,
		UserID:   user.ID,
		PlayerID: player.ID,
		Score:    2,
		Content:  "Not good",
		Status:   model.ReviewStatusPending,
	}
	require.NoError(t, repo.Create(ctx, reviewObj))

	err := repo.Delete(ctx, reviewObj.ID)
	require.NoError(t, err)

	_, err = repo.Get(ctx, reviewObj.ID)
	assert.Error(t, err)
}

func TestReviewRepository_List(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := review.NewReviewRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "list_review_user")
	playerUser := CreateUniqueTestUser(t, db, "list_review_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "list_review_game")

	// Create multiple reviews
	for i := 0; i < 3; i++ {
		order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, 10000)
		reviewObj := &model.Review{
			OrderID:  order.ID,
			UserID:   user.ID,
			PlayerID: player.ID,
			Score:    5,
			Content:  "Great!",
			Status:   model.ReviewStatusApproved,
		}
		require.NoError(t, repo.Create(ctx, reviewObj))
	}

	reviews, total, err := repo.List(ctx, repository.ReviewListOptions{
		PlayerID: &player.ID,
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, reviews, 3)
}

func TestReviewRepository_ListPending(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := review.NewReviewRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "pending_review_user")
	playerUser := CreateUniqueTestUser(t, db, "pending_review_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "pending_review_game")

	// Create pending reviews
	for i := 0; i < 2; i++ {
		order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, 10000)
		reviewObj := &model.Review{
			OrderID:  order.ID,
			UserID:   user.ID,
			PlayerID: player.ID,
			Score:    4,
			Content:  "Pending review",
			Status:   model.ReviewStatusPending,
		}
		require.NoError(t, repo.Create(ctx, reviewObj))
	}

	reviews, total, err := repo.ListPending(ctx, 1, 10)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(2))
	for _, r := range reviews {
		assert.Equal(t, model.ReviewStatusPending, r.Status)
	}
}

func TestReviewRepository_UpdateStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := review.NewReviewRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "status_review_user")
	playerUser := CreateUniqueTestUser(t, db, "status_review_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "status_review_game")
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, 10000)

	reviewObj := &model.Review{
		OrderID:  order.ID,
		UserID:   user.ID,
		PlayerID: player.ID,
		Score:    5,
		Content:  "Awaiting approval",
		Status:   model.ReviewStatusPending,
	}
	require.NoError(t, repo.Create(ctx, reviewObj))

	// Approve
	err := repo.UpdateStatus(ctx, reviewObj.ID, model.ReviewStatusApproved, "")
	require.NoError(t, err)

	got, err := repo.Get(ctx, reviewObj.ID)
	require.NoError(t, err)
	assert.Equal(t, model.ReviewStatusApproved, got.Status)
}

func TestReviewRepository_UpdateStatus_Rejected(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := review.NewReviewRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "reject_review_user")
	playerUser := CreateUniqueTestUser(t, db, "reject_review_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "reject_review_game")
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, 10000)

	reviewObj := &model.Review{
		OrderID:  order.ID,
		UserID:   user.ID,
		PlayerID: player.ID,
		Score:    1,
		Content:  "Bad words here",
		Status:   model.ReviewStatusPending,
	}
	require.NoError(t, repo.Create(ctx, reviewObj))

	// Reject with reason
	err := repo.UpdateStatus(ctx, reviewObj.ID, model.ReviewStatusRejected, "Contains inappropriate content")
	require.NoError(t, err)

	got, err := repo.Get(ctx, reviewObj.ID)
	require.NoError(t, err)
	assert.Equal(t, model.ReviewStatusRejected, got.Status)
	assert.Equal(t, "Contains inappropriate content", got.RejectionReason)
}

func TestReviewRepository_BatchUpdateStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := review.NewReviewRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "batch_review_user")
	playerUser := CreateUniqueTestUser(t, db, "batch_review_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "batch_review_game")

	var ids []uint64
	for i := 0; i < 3; i++ {
		order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, 10000)
		reviewObj := &model.Review{
			OrderID:  order.ID,
			UserID:   user.ID,
			PlayerID: player.ID,
			Score:    5,
			Content:  "Batch review",
			Status:   model.ReviewStatusPending,
		}
		require.NoError(t, repo.Create(ctx, reviewObj))
		ids = append(ids, reviewObj.ID)
	}

	// Batch approve
	err := repo.BatchUpdateStatus(ctx, ids, model.ReviewStatusApproved, "")
	require.NoError(t, err)

	// Verify
	for _, id := range ids {
		got, err := repo.Get(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, model.ReviewStatusApproved, got.Status)
	}
}

func TestReviewRepository_GetStats(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := review.NewReviewRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "stats_review_user")
	playerUser := CreateUniqueTestUser(t, db, "stats_review_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "stats_review_game")

	// Create approved reviews with different scores
	scores := []model.Rating{5, 5, 4, 4, 3}
	for _, score := range scores {
		order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, 10000)
		reviewObj := &model.Review{
			OrderID:  order.ID,
			UserID:   user.ID,
			PlayerID: player.ID,
			Score:    score,
			Content:  "Stats review",
			Status:   model.ReviewStatusApproved,
		}
		require.NoError(t, repo.Create(ctx, reviewObj))
	}

	stats, err := repo.GetStats(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, stats.TotalReviews, int64(5))
	assert.Greater(t, stats.AverageRating, float64(0))
	assert.NotNil(t, stats.RatingDistribution)
}
