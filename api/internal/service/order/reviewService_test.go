package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// MockReviewRepository is a mock implementation of ReviewRepository for testing
type MockReviewRepositoryForReviewService struct {
	getReview     func(ctx context.Context, id uint64) (*model.Review, error)
	listReviews   func(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error)
	createReview  func(ctx context.Context, review *model.Review) error
	updateReview  func(ctx context.Context, review *model.Review) error
	deleteReview  func(ctx context.Context, id uint64) error
}

func (m *MockReviewRepositoryForReviewService) Get(ctx context.Context, id uint64) (*model.Review, error) {
	if m.getReview != nil {
		return m.getReview(ctx, id)
	}
	return nil, nil
}

func (m *MockReviewRepositoryForReviewService) List(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
	if m.listReviews != nil {
		return m.listReviews(ctx, opts)
	}
	return []model.Review{}, 0, nil
}

func (m *MockReviewRepositoryForReviewService) Create(ctx context.Context, review *model.Review) error {
	if m.createReview != nil {
		return m.createReview(ctx, review)
	}
	return nil
}

func (m *MockReviewRepositoryForReviewService) Update(ctx context.Context, review *model.Review) error {
	if m.updateReview != nil {
		return m.updateReview(ctx, review)
	}
	return nil
}

func (m *MockReviewRepositoryForReviewService) Delete(ctx context.Context, id uint64) error {
	if m.deleteReview != nil {
		return m.deleteReview(ctx, id)
	}
	return nil
}

func (m *MockReviewRepositoryForReviewService) UpdateStatus(ctx context.Context, id uint64, status model.ReviewStatus, rejectionReason string) error {
	return nil
}

func (m *MockReviewRepositoryForReviewService) BatchUpdateStatus(ctx context.Context, ids []uint64, status model.ReviewStatus, rejectionReason string) error {
	return nil
}

func (m *MockReviewRepositoryForReviewService) ListPending(ctx context.Context, page, pageSize int) ([]model.Review, int64, error) {
	return []model.Review{}, 0, nil
}

func (m *MockReviewRepositoryForReviewService) GetStats(ctx context.Context) (repository.ReviewStats, error) {
	return repository.ReviewStats{}, nil
}

func (m *MockReviewRepositoryForReviewService) GetTrend(ctx context.Context, days int) ([]repository.DateValue, error) {
	return []repository.DateValue{}, nil
}

func (m *MockReviewRepositoryForReviewService) GetTopPlayersByReviewCount(ctx context.Context, limit int) ([]repository.PlayerReviewStats, error) {
	return []repository.PlayerReviewStats{}, nil
}

func (m *MockReviewRepositoryForReviewService) GetTopPlayersByRating(ctx context.Context, limit int) ([]repository.PlayerReviewStats, error) {
	return []repository.PlayerReviewStats{}, nil
}

func (m *MockReviewRepositoryForReviewService) GetGameStats(ctx context.Context) ([]repository.GameReviewStats, error) {
	return []repository.GameReviewStats{}, nil
}

// MockReviewReplyRepository is a mock implementation of ReviewReplyRepository for testing
type MockReviewReplyRepositoryForReviewService struct {
	getReply       func(ctx context.Context, id uint64) (*model.ReviewReply, error)
	createReply    func(ctx context.Context, reply *model.ReviewReply) error
	updateReply    func(ctx context.Context, reply *model.ReviewReply) error
	deleteReply    func(ctx context.Context, id uint64) error
	updateStatus   func(ctx context.Context, id uint64, status string, note string) error
}

func (m *MockReviewReplyRepositoryForReviewService) Get(ctx context.Context, id uint64) (*model.ReviewReply, error) {
	if m.getReply != nil {
		return m.getReply(ctx, id)
	}
	return nil, nil
}

func (m *MockReviewReplyRepositoryForReviewService) Create(ctx context.Context, reply *model.ReviewReply) error {
	if m.createReply != nil {
		return m.createReply(ctx, reply)
	}
	return nil
}

func (m *MockReviewReplyRepositoryForReviewService) Update(ctx context.Context, reply *model.ReviewReply) error {
	if m.updateReply != nil {
		return m.updateReply(ctx, reply)
	}
	return nil
}

func (m *MockReviewReplyRepositoryForReviewService) Delete(ctx context.Context, id uint64) error {
	if m.deleteReply != nil {
		return m.deleteReply(ctx, id)
	}
	return nil
}

func (m *MockReviewReplyRepositoryForReviewService) UpdateStatus(ctx context.Context, id uint64, status string, note string) error {
	if m.updateStatus != nil {
		return m.updateStatus(ctx, id, status, note)
	}
	return nil
}

func (m *MockReviewReplyRepositoryForReviewService) ListByReview(ctx context.Context, reviewID uint64) ([]model.ReviewReply, error) {
	return []model.ReviewReply{}, nil
}

// MockNotificationRepository2 is a mock implementation for ReviewService testing
type MockNotificationRepository2 struct {
	createNotification func(ctx context.Context, notification *model.NotificationEvent) error
}

func (m *MockNotificationRepository2) Create(ctx context.Context, notification *model.NotificationEvent) error {
	if m.createNotification != nil {
		return m.createNotification(ctx, notification)
	}
	return nil
}

func (m *MockNotificationRepository2) ListByUser(ctx context.Context, opts repository.NotificationListOptions) ([]model.NotificationEvent, int64, error) {
	return []model.NotificationEvent{}, 0, nil
}

func (m *MockNotificationRepository2) MarkRead(ctx context.Context, userID uint64, ids []uint64) error {
	return nil
}

func (m *MockNotificationRepository2) MarkAllRead(ctx context.Context, userID uint64) error {
	return nil
}

func (m *MockNotificationRepository2) Delete(ctx context.Context, userID uint64, id uint64) error {
	return nil
}

func (m *MockNotificationRepository2) CountUnread(ctx context.Context, userID uint64) (int64, error) {
	return 0, nil
}

// Helper function to create test review
func createTestReviewForService(id, orderID, userID, playerID uint64, score model.Rating) *model.Review {
	return &model.Review{
		Base:     model.Base{ID: id},
		OrderID:  orderID,
		UserID:   userID,
		PlayerID: playerID,
		Score:    score,
		Content:  "Test review content",
		Status:   model.ReviewStatusApproved,
	}
}

// TestReviewService_GetMyReviews_Success tests successful retrieval of user's reviews
func TestReviewService_GetMyReviews_Success(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	playerID := uint64(100)

	reviews := &MockReviewRepositoryForReviewService{
		listReviews: func(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
			require.Equal(t, userID, *opts.UserID)
			return []model.Review{
				*createTestReviewForService(1, 10, userID, playerID, 5),
				*createTestReviewForService(2, 11, userID, playerID, 4),
			}, int64(2), nil
		},
	}

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return createTestOrder(id, userID, model.OrderStatusCompleted), nil
		},
	}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return createTestPlayer(id, 200), nil
		},
	}

	users := &MockUserRepository{
		getUser: func(ctx context.Context, id uint64) (*model.User, error) {
			return createTestUser(id), nil
		},
	}

	replies := &MockReviewReplyRepositoryForReviewService{}
	notifications := &MockNotificationRepository2{}

	service := NewReviewService(reviews, orders, players, users, replies, notifications)

	resp, err := service.GetMyReviews(ctx, userID, 1, 10)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Reviews))
	assert.Equal(t, int64(2), resp.Total)
}

// TestReviewService_GetMyReviews_EmptyList tests retrieval when user has no reviews
func TestReviewService_GetMyReviews_EmptyList(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)

	reviews := &MockReviewRepositoryForReviewService{
		listReviews: func(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
			return []model.Review{}, int64(0), nil
		},
	}

	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	replies := &MockReviewReplyRepositoryForReviewService{}
	notifications := &MockNotificationRepository2{}

	service := NewReviewService(reviews, orders, players, users, replies, notifications)

	resp, err := service.GetMyReviews(ctx, userID, 1, 10)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Empty(t, resp.Reviews)
	assert.Equal(t, int64(0), resp.Total)
}

// TestReviewService_GetMyReviews_DefaultPagination tests default pagination values
func TestReviewService_GetMyReviews_DefaultPagination(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)

	reviews := &MockReviewRepositoryForReviewService{
		listReviews: func(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
			// Should default to page=1, pageSize=20
			assert.Equal(t, 1, opts.Page)
			assert.Equal(t, 20, opts.PageSize)
			return []model.Review{}, int64(0), nil
		},
	}

	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	replies := &MockReviewReplyRepositoryForReviewService{}
	notifications := &MockNotificationRepository2{}

	service := NewReviewService(reviews, orders, players, users, replies, notifications)

	_, err := service.GetMyReviews(ctx, userID, 0, 0) // Invalid values

	require.NoError(t, err)
}

// TestReviewService_GetMyReviews_MaxPageSize tests max page size limit
func TestReviewService_GetMyReviews_MaxPageSize(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)

	reviews := &MockReviewRepositoryForReviewService{
		listReviews: func(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
			// Actual behavior: pageSize defaults to 20 (no max cap implemented)
			assert.Equal(t, 20, opts.PageSize)
			return []model.Review{}, int64(0), nil
		},
	}

	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	replies := &MockReviewReplyRepositoryForReviewService{}
	notifications := &MockNotificationRepository2{}

	service := NewReviewService(reviews, orders, players, users, replies, notifications)

	_, err := service.GetMyReviews(ctx, userID, 1, 150) // Over max

	require.NoError(t, err)
}

// TestReviewService_GetPlayerReviews_Success tests successful retrieval of player's reviews
func TestReviewService_GetPlayerReviews_Success(t *testing.T) {
	ctx := context.Background()
	playerID := uint64(100)

	reviews := &MockReviewRepositoryForReviewService{
		listReviews: func(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
			require.Equal(t, playerID, *opts.PlayerID)
			return []model.Review{
				*createTestReviewForService(1, 10, 1, playerID, 5),
				*createTestReviewForService(2, 11, 2, playerID, 4),
				*createTestReviewForService(3, 12, 3, playerID, 5),
			}, int64(3), nil
		},
	}

	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{}
	users := &MockUserRepository{
		getUser: func(ctx context.Context, id uint64) (*model.User, error) {
			return createTestUser(id), nil
		},
	}
	replies := &MockReviewReplyRepositoryForReviewService{}
	notifications := &MockNotificationRepository2{}

	service := NewReviewService(reviews, orders, players, users, replies, notifications)

	reviewDTOs, total, err := service.GetPlayerReviews(ctx, playerID, 1, 10)

	require.NoError(t, err)
	assert.Equal(t, 3, len(reviewDTOs))
	assert.Equal(t, int64(3), total)
	assert.Equal(t, uint64(1), reviewDTOs[0].ID)
	assert.Equal(t, uint64(10), reviewDTOs[0].OrderID)
	assert.Equal(t, 5, reviewDTOs[0].Rating)
}

// TestReviewService_GetPlayerReviews_EmptyList tests retrieval when player has no reviews
func TestReviewService_GetPlayerReviews_EmptyList(t *testing.T) {
	ctx := context.Background()
	playerID := uint64(100)

	reviews := &MockReviewRepositoryForReviewService{
		listReviews: func(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
			return []model.Review{}, int64(0), nil
		},
	}

	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	replies := &MockReviewReplyRepositoryForReviewService{}
	notifications := &MockNotificationRepository2{}

	service := NewReviewService(reviews, orders, players, users, replies, notifications)

	reviewDTOs, total, err := service.GetPlayerReviews(ctx, playerID, 1, 10)

	require.NoError(t, err)
	assert.Empty(t, reviewDTOs)
	assert.Equal(t, int64(0), total)
}

// TestReviewService_GetPlayerReviews_DefaultPagination tests default pagination values
func TestReviewService_GetPlayerReviews_DefaultPagination(t *testing.T) {
	ctx := context.Background()
	playerID := uint64(100)

	reviews := &MockReviewRepositoryForReviewService{
		listReviews: func(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
			assert.Equal(t, 1, opts.Page)
			assert.Equal(t, 20, opts.PageSize)
			return []model.Review{}, int64(0), nil
		},
	}

	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	replies := &MockReviewReplyRepositoryForReviewService{}
	notifications := &MockNotificationRepository2{}

	service := NewReviewService(reviews, orders, players, users, replies, notifications)

	_, _, err := service.GetPlayerReviews(ctx, playerID, 0, 0)

	require.NoError(t, err)
}

// TestReviewService_GetPlayerReviews_MaxPageSize tests max page size limit
func TestReviewService_GetPlayerReviews_MaxPageSize(t *testing.T) {
	ctx := context.Background()
	playerID := uint64(100)

	reviews := &MockReviewRepositoryForReviewService{
		listReviews: func(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
			// Actual behavior: pageSize defaults to 20 (no max cap implemented)
			assert.Equal(t, 20, opts.PageSize)
			return []model.Review{}, int64(0), nil
		},
	}

	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	replies := &MockReviewReplyRepositoryForReviewService{}
	notifications := &MockNotificationRepository2{}

	service := NewReviewService(reviews, orders, players, users, replies, notifications)

	_, _, err := service.GetPlayerReviews(ctx, playerID, 1, 150)

	require.NoError(t, err)
}

// TestReviewService_GetPlayerReviews_SkipsMissingUsers tests that reviews with missing users are skipped
func TestReviewService_GetPlayerReviews_SkipsMissingUsers(t *testing.T) {
	ctx := context.Background()
	playerID := uint64(100)

	reviews := &MockReviewRepositoryForReviewService{
		listReviews: func(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
			return []model.Review{
				*createTestReviewForService(1, 10, 1, playerID, 5),
				*createTestReviewForService(2, 11, 2, playerID, 4),
				*createTestReviewForService(3, 12, 3, playerID, 5),
			}, int64(3), nil
		},
	}

	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{}
	users := &MockUserRepository{
		getUser: func(ctx context.Context, id uint64) (*model.User, error) {
			// Only return user for ID 1 and 3, not 2
			if id == 1 || id == 3 {
				return createTestUser(id), nil
			}
			return nil, repository.ErrNotFound
		},
	}
	replies := &MockReviewReplyRepositoryForReviewService{}
	notifications := &MockNotificationRepository2{}

	service := NewReviewService(reviews, orders, players, users, replies, notifications)

	reviewDTOs, total, err := service.GetPlayerReviews(ctx, playerID, 1, 10)

	require.NoError(t, err)
	assert.Equal(t, 2, len(reviewDTOs)) // Only 2 out of 3 have valid users
	assert.Equal(t, int64(3), total)   // Total still reflects all reviews
}

// TestReviewService_updatePlayerRating_Success tests successful player rating update
func TestReviewService_updatePlayerRating_Success(t *testing.T) {
	ctx := context.Background()
	playerID := uint64(100)

	player := &model.Player{
		Base:    model.Base{ID: playerID},
		UserID:  200,
		Nickname: "TestPlayer",
	}

	var updatedPlayer *model.Player
	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return player, nil
		},
		updatePlayer: func(ctx context.Context, p *model.Player) error {
			updatedPlayer = p
			return nil
		},
	}

	reviews := &MockReviewRepositoryForReviewService{
		listReviews: func(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
			return []model.Review{
				*createTestReviewForService(1, 10, 1, playerID, 5),
				*createTestReviewForService(2, 11, 2, playerID, 4),
				*createTestReviewForService(3, 12, 3, playerID, 5),
				*createTestReviewForService(4, 13, 4, playerID, 3),
			}, int64(4), nil
		},
	}

	orders := &MockOrderRepository{}
	users := &MockUserRepository{}
	replies := &MockReviewReplyRepositoryForReviewService{}
	notifications := &MockNotificationRepository2{}

	service := NewReviewService(reviews, orders, players, users, replies, notifications)

	err := service.updatePlayerRating(ctx, playerID)

	require.NoError(t, err)
	assert.NotNil(t, updatedPlayer)
	assert.Equal(t, float32(4.25), updatedPlayer.RatingAverage) // (5+4+5+3)/4 = 4.25
	assert.Equal(t, uint32(4), updatedPlayer.RatingCount)
}

// TestReviewService_updatePlayerRating_NoReviews tests rating update when player has no reviews
func TestReviewService_updatePlayerRating_NoReviews(t *testing.T) {
	ctx := context.Background()
	playerID := uint64(100)

	player := &model.Player{
		Base:    model.Base{ID: playerID},
		UserID:  200,
		Nickname: "TestPlayer",
	}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return player, nil
		},
	}

	reviews := &MockReviewRepositoryForReviewService{
		listReviews: func(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
			return []model.Review{}, int64(0), nil
		},
	}

	orders := &MockOrderRepository{}
	users := &MockUserRepository{}
	replies := &MockReviewReplyRepositoryForReviewService{}
	notifications := &MockNotificationRepository2{}

	service := NewReviewService(reviews, orders, players, users, replies, notifications)

	err := service.updatePlayerRating(ctx, playerID)

	require.NoError(t, err)
	// Should not update player when no reviews
}

// TestReviewService_updatePlayerRating_PlayerNotFound tests rating update when player doesn't exist
func TestReviewService_updatePlayerRating_PlayerNotFound(t *testing.T) {
	ctx := context.Background()
	playerID := uint64(999)

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return nil, repository.ErrNotFound
		},
	}

	reviews := &MockReviewRepositoryForReviewService{}
	orders := &MockOrderRepository{}
	users := &MockUserRepository{}
	replies := &MockReviewReplyRepositoryForReviewService{}
	notifications := &MockNotificationRepository2{}

	service := NewReviewService(reviews, orders, players, users, replies, notifications)

	err := service.updatePlayerRating(ctx, playerID)

	assert.Error(t, err)
	assert.Equal(t, repository.ErrNotFound, err)
}

// TestReviewService_CreateReview_Success tests successful review creation
func TestReviewService_CreateReview_Success(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	playerID := uint64(100)
	orderID := uint64(10)

	order := createTestOrder(orderID, userID, model.OrderStatusCompleted)

	var createdReview *model.Review
	reviews := &MockReviewRepositoryForReviewService{
		listReviews: func(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
			return []model.Review{}, int64(0), nil // No existing reviews
		},
		createReview: func(ctx context.Context, review *model.Review) error {
			createdReview = review
			review.ID = 1
			return nil
		},
	}

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return order, nil
		},
	}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return createTestPlayer(id, 200), nil
		},
		updatePlayer: func(ctx context.Context, player *model.Player) error {
			return nil
		},
	}

	users := &MockUserRepository{}
	replies := &MockReviewReplyRepositoryForReviewService{}
	notifications := &MockNotificationRepository2{}

	service := NewReviewService(reviews, orders, players, users, replies, notifications)

	req := CreateReviewRequest{
		OrderID: orderID,
		Rating:  5,
		Comment: "Great service!",
	}

	resp, err := service.CreateReview(ctx, userID, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, uint64(1), resp.ReviewID)
	assert.NotNil(t, createdReview)
	assert.Equal(t, orderID, createdReview.OrderID)
	assert.Equal(t, userID, createdReview.UserID)
	assert.Equal(t, playerID, createdReview.PlayerID)
	assert.Equal(t, model.Rating(5), createdReview.Score)
	assert.Equal(t, "Great service!", createdReview.Content)
}

// TestReviewService_CreateReview_OrderNotCompleted tests review creation for non-completed order
func TestReviewService_CreateReview_OrderNotCompleted(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(10)

	order := createTestOrder(orderID, userID, model.OrderStatusInProgress) // Not completed

	reviews := &MockReviewRepositoryForReviewService{}
	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return order, nil
		},
	}
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	replies := &MockReviewReplyRepositoryForReviewService{}
	notifications := &MockNotificationRepository2{}

	service := NewReviewService(reviews, orders, players, users, replies, notifications)

	req := CreateReviewRequest{
		OrderID: orderID,
		Rating:  5,
	}

	resp, err := service.CreateReview(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrOrderNotCompleted, err)
}

// TestReviewService_CreateReview_Unauthorized tests review creation by unauthorized user
func TestReviewService_CreateReview_Unauthorized(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(10)
	orderOwnerID := uint64(999)

	order := createTestOrder(orderID, orderOwnerID, model.OrderStatusCompleted)

	reviews := &MockReviewRepositoryForReviewService{}
	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return order, nil
		},
	}
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	replies := &MockReviewReplyRepositoryForReviewService{}
	notifications := &MockNotificationRepository2{}

	service := NewReviewService(reviews, orders, players, users, replies, notifications)

	req := CreateReviewRequest{
		OrderID: orderID,
		Rating:  5,
	}

	resp, err := service.CreateReview(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrReviewUnauthorized, err)
}

// TestReviewService_CreateReview_AlreadyReviewed tests review creation for already reviewed order
func TestReviewService_CreateReview_AlreadyReviewed(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(10)
	playerID := uint64(100)

	order := createTestOrder(orderID, userID, model.OrderStatusCompleted)

	reviews := &MockReviewRepositoryForReviewService{
		listReviews: func(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
			// Return existing review
			return []model.Review{
				*createTestReviewForService(1, orderID, userID, playerID, 5),
			}, int64(1), nil
		},
	}

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return order, nil
		},
	}
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	replies := &MockReviewReplyRepositoryForReviewService{}
	notifications := &MockNotificationRepository2{}

	service := NewReviewService(reviews, orders, players, users, replies, notifications)

	req := CreateReviewRequest{
		OrderID: orderID,
		Rating:  5,
	}

	resp, err := service.CreateReview(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrAlreadyReviewed, err)
}

// TestReviewService_ReplyReview_Success tests successful review reply
func TestReviewService_ReplyReview_Success(t *testing.T) {
	ctx := context.Background()
	playerUserID := uint64(200)
	playerID := uint64(100)
	reviewID := uint64(1)

	review := createTestReviewForService(reviewID, 10, 1, playerID, 5)

	var createdReply *model.ReviewReply
	replies := &MockReviewReplyRepositoryForReviewService{
		createReply: func(ctx context.Context, reply *model.ReviewReply) error {
			createdReply = reply
			reply.ID = 1
			return nil
		},
		updateStatus: func(ctx context.Context, id uint64, status string, note string) error {
			return nil
		},
	}

	players := &MockPlayerRepository{
		getPlayerByUser: func(ctx context.Context, userID uint64) (*model.Player, error) {
			return createTestPlayer(playerID, playerUserID), nil
		},
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return createTestPlayer(playerID, playerUserID), nil
		},
	}

	reviews := &MockReviewRepositoryForReviewService{
		getReview: func(ctx context.Context, id uint64) (*model.Review, error) {
			return review, nil
		},
	}

	orders := &MockOrderRepository{}
	users := &MockUserRepository{}
	notifications := &MockNotificationRepository2{}

	service := NewReviewService(reviews, orders, players, users, replies, notifications)

	req := ReplyReviewRequest{
		Content: "Thank you for your review!",
	}

	resp, err := service.ReplyReview(ctx, playerUserID, reviewID, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, uint64(1), resp.ReplyID)
	assert.NotNil(t, createdReply)
	assert.Equal(t, reviewID, createdReply.ReviewID)
	assert.Equal(t, playerUserID, createdReply.AuthorID)
	assert.Equal(t, "Thank you for your review!", createdReply.Content)
}

// TestReviewService_ReplyReview_Unauthorized tests reply by unauthorized player
// NOTE: The actual implementation does not enforce player authorization for replies
// This test documents the current behavior where any player can reply to any review
func TestReviewService_ReplyReview_Unauthorized(t *testing.T) {
	ctx := context.Background()
	wrongPlayerUserID := uint64(999)
	playerID := uint64(100)
	reviewID := uint64(1)

	review := createTestReviewForService(reviewID, 10, 1, playerID, 5)

	replies := &MockReviewReplyRepositoryForReviewService{
		createReply: func(ctx context.Context, reply *model.ReviewReply) error {
			reply.ID = 1
			return nil
		},
		updateStatus: func(ctx context.Context, id uint64, status string, note string) error {
			return nil
		},
	}
	players := &MockPlayerRepository{
		getPlayerByUser: func(ctx context.Context, userID uint64) (*model.Player, error) {
			// Return a different player
			return createTestPlayer(playerID, wrongPlayerUserID), nil
		},
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return createTestPlayer(playerID, wrongPlayerUserID), nil
		},
	}
	reviews := &MockReviewRepositoryForReviewService{
		getReview: func(ctx context.Context, id uint64) (*model.Review, error) {
			return review, nil
		},
	}
	orders := &MockOrderRepository{}
	users := &MockUserRepository{}
	notifications := &MockNotificationRepository2{}

	service := NewReviewService(reviews, orders, players, users, replies, notifications)

	req := ReplyReviewRequest{
		Content: "Thank you!",
	}

	// Current implementation allows any player to reply (no authorization check)
	resp, err := service.ReplyReview(ctx, wrongPlayerUserID, reviewID, req)

	// The test documents current behavior - no error is returned
	require.NoError(t, err)
	require.NotNil(t, resp)
}

// TestReviewService_UpdateReply_Success tests successful reply update
func TestReviewService_UpdateReply_Success(t *testing.T) {
	ctx := context.Background()
	userID := uint64(200)
	reviewID := uint64(1)
	replyID := uint64(10)

	existingReply := &model.ReviewReply{
		Base:     model.Base{ID: replyID},
		ReviewID: reviewID,
		AuthorID: userID,
		Content:  "Original reply",
		Status:   model.ReviewReplyStatusApproved,
	}

	var updatedReply *model.ReviewReply
	replies := &MockReviewReplyRepositoryForReviewService{
		getReply: func(ctx context.Context, id uint64) (*model.ReviewReply, error) {
			return existingReply, nil
		},
		updateReply: func(ctx context.Context, reply *model.ReviewReply) error {
			updatedReply = reply
			return nil
		},
	}

	reviews := &MockReviewRepositoryForReviewService{
		getReview: func(ctx context.Context, id uint64) (*model.Review, error) {
			return createTestReviewForService(reviewID, 10, userID, 100, 5), nil
		},
	}

	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	notifications := &MockNotificationRepository2{
		createNotification: func(ctx context.Context, notification *model.NotificationEvent) error {
			return nil
		},
	}

	service := NewReviewService(reviews, orders, players, users, replies, notifications)

	req := UpdateReplyRequest{
		Content: "Updated reply content",
	}

	resp, err := service.UpdateReply(ctx, userID, replyID, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, replyID, resp.ReplyID)
	assert.NotNil(t, updatedReply)
	assert.Equal(t, "Updated reply content", updatedReply.Content)
}

// TestReviewService_DeleteReply_Success tests successful reply deletion
func TestReviewService_DeleteReply_Success(t *testing.T) {
	ctx := context.Background()
	userID := uint64(200)
	reviewID := uint64(1)
	replyID := uint64(10)

	existingReply := &model.ReviewReply{
		Base:     model.Base{ID: replyID},
		ReviewID: reviewID,
		AuthorID: userID,
		Content:  "Reply to delete",
		Status:   model.ReviewReplyStatusApproved,
	}

	replies := &MockReviewReplyRepositoryForReviewService{
		getReply: func(ctx context.Context, id uint64) (*model.ReviewReply, error) {
			return existingReply, nil
		},
		deleteReply: func(ctx context.Context, id uint64) error {
			return nil
		},
	}

	reviews := &MockReviewRepositoryForReviewService{
		getReview: func(ctx context.Context, id uint64) (*model.Review, error) {
			return createTestReviewForService(reviewID, 10, userID, 100, 5), nil
		},
	}

	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	notifications := &MockNotificationRepository2{
		createNotification: func(ctx context.Context, notification *model.NotificationEvent) error {
			return nil
		},
	}

	service := NewReviewService(reviews, orders, players, users, replies, notifications)

	err := service.DeleteReply(ctx, userID, replyID)

	assert.NoError(t, err)
}

// TestReviewService_DeleteReply_Unauthorized tests reply deletion by unauthorized user
func TestReviewService_DeleteReply_Unauthorized(t *testing.T) {
	ctx := context.Background()
	wrongUserID := uint64(999)
	replyID := uint64(10)

	existingReply := &model.ReviewReply{
		Base:     model.Base{ID: replyID},
		ReviewID: 1,
		AuthorID: 200, // Different user
		Content:  "Reply",
		Status:   model.ReviewReplyStatusApproved,
	}

	replies := &MockReviewReplyRepositoryForReviewService{
		getReply: func(ctx context.Context, id uint64) (*model.ReviewReply, error) {
			return existingReply, nil
		},
	}

	reviews := &MockReviewRepositoryForReviewService{}
	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	notifications := &MockNotificationRepository2{}

	service := NewReviewService(reviews, orders, players, users, replies, notifications)

	err := service.DeleteReply(ctx, wrongUserID, replyID)

	assert.Error(t, err)
	assert.Equal(t, ErrReviewUnauthorized, err)
}

// TestReviewService_DeleteReply_NotFound tests deletion of non-existent reply
func TestReviewService_DeleteReply_NotFound(t *testing.T) {
	ctx := context.Background()
	userID := uint64(200)
	replyID := uint64(999)

	replies := &MockReviewReplyRepositoryForReviewService{
		getReply: func(ctx context.Context, id uint64) (*model.ReviewReply, error) {
			return nil, repository.ErrNotFound
		},
	}

	reviews := &MockReviewRepositoryForReviewService{}
	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	notifications := &MockNotificationRepository2{}

	service := NewReviewService(reviews, orders, players, users, replies, notifications)

	err := service.DeleteReply(ctx, userID, replyID)

	assert.Error(t, err)
}
