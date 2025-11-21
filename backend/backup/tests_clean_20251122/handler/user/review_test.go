package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service/review"
)

// ---- Fake ReviewRepository for user_review tests ----

type mockReviewRepoForUserReview struct {
	reviews map[uint64]*model.Review
}

func newMockReviewRepoForUserReview() *mockReviewRepoForUserReview {
	return &mockReviewRepoForUserReview{
		reviews: map[uint64]*model.Review{
			1: {Base: model.Base{ID: 1}, OrderID: 10, UserID: 100, PlayerID: 1, Score: 5, Content: "Great!"},
			2: {Base: model.Base{ID: 2}, OrderID: 11, UserID: 100, PlayerID: 2, Score: 4, Content: "Good"},
		},
	}
}

func (m *mockReviewRepoForUserReview) List(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
	var result []model.Review
	for _, r := range m.reviews {
		// Filter by user if specified
		if opts.UserID != nil && *opts.UserID != r.UserID {
			continue
		}
		// Filter by order if specified
		if opts.OrderID != nil && *opts.OrderID != r.OrderID {
			continue
		}
		result = append(result, *r)
	}
	return result, int64(len(result)), nil
}

func (m *mockReviewRepoForUserReview) Get(ctx context.Context, id uint64) (*model.Review, error) {
	if r, ok := m.reviews[id]; ok {
		return r, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockReviewRepoForUserReview) GetByOrderID(ctx context.Context, orderID uint64) (*model.Review, error) {
	for _, r := range m.reviews {
		if r.OrderID == orderID {
			return r, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *mockReviewRepoForUserReview) Create(ctx context.Context, review *model.Review) error {
	review.ID = uint64(len(m.reviews) + 1)
	m.reviews[review.ID] = review
	return nil
}

func (m *mockReviewRepoForUserReview) Update(ctx context.Context, review *model.Review) error {
	m.reviews[review.ID] = review
	return nil
}

func (m *mockReviewRepoForUserReview) Delete(ctx context.Context, id uint64) error {
	delete(m.reviews, id)
	return nil
}

// ---- Fake ReviewReplyRepository for user_review tests ----

type mockReviewReplyRepoForUserReview struct{}

func (m *mockReviewReplyRepoForUserReview) List(ctx context.Context, opts interface{}) ([]model.ReviewReply, int64, error) {
	return []model.ReviewReply{}, 0, nil
}

func (m *mockReviewReplyRepoForUserReview) Get(ctx context.Context, id uint64) (*model.ReviewReply, error) {
	return nil, repository.ErrNotFound
}

func (m *mockReviewReplyRepoForUserReview) GetByReviewID(ctx context.Context, reviewID uint64) (*model.ReviewReply, error) {
	return nil, repository.ErrNotFound
}

func (m *mockReviewReplyRepoForUserReview) Create(ctx context.Context, reply *model.ReviewReply) error {
	return nil
}

func (m *mockReviewReplyRepoForUserReview) Update(ctx context.Context, reply *model.ReviewReply) error {
	return nil
}

func (m *mockReviewReplyRepoForUserReview) Delete(ctx context.Context, id uint64) error {
	return nil
}

func (m *mockReviewReplyRepoForUserReview) ListByReview(ctx context.Context, reviewID uint64) ([]model.ReviewReply, error) {
	return []model.ReviewReply{}, nil
}

func (m *mockReviewReplyRepoForUserReview) UpdateStatus(ctx context.Context, id uint64, status string, reason string) error {
	return nil
}

// ---- Tests for user_review.go ----
// Note: Fake repositories are shared from user_order_test.go

func TestCreateReviewHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	reviewRepo := newMockReviewRepoForUserReview()
	orderRepo := newFakeOrderRepository()
	// Create a completed order for user 100 with a different order ID (not 10 or 11 which already have reviews)
	playerID := uint64(1)
	order := &model.Order{
		UserID:   100,
		PlayerID: &playerID,
		Status:   model.OrderStatusCompleted,
		ItemID:   1,
		OrderNo:  "REVIEW-TEST-CREATE",
	}
	orderRepo.Create(context.Background(), order)
	// order.ID will be set by Create method to 1 (first order in empty repo)

	reviewSvc := review.NewReviewService(reviewRepo, orderRepo, &fakePlayerRepository{}, &fakeUserRepository{}, &mockReviewReplyRepoForUserReview{})

	router := gin.New()
	router.POST("/user/reviews", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		createReviewHandler(c, reviewSvc)
	})

	reqBody := review.CreateReviewRequest{
		OrderID: order.ID, // Use the actual ID assigned by Create
		Rating:  5,
		Comment: "Excellent service!",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/user/reviews", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp model.APIResponse[review.CreateReviewResponse]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Fatalf("Expected success=true, got %v", resp.Success)
	}
}

func TestCreateReviewHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	reviewRepo := newMockReviewRepoForUserReview()
	reviewSvc := review.NewReviewService(reviewRepo, newFakeOrderRepository(), &fakePlayerRepository{}, &fakeUserRepository{}, &mockReviewReplyRepoForUserReview{})

	router := gin.New()
	router.POST("/user/reviews", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		createReviewHandler(c, reviewSvc)
	})

	req := httptest.NewRequest(http.MethodPost, "/user/reviews", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}

func TestCreateReviewHandler_AlreadyReviewed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	reviewRepo := newMockReviewRepoForUserReview()
	orderRepo := newFakeOrderRepository()
	// Create a completed order that matches the review mock data (OrderID 10)
	playerID := uint64(1)
	order := &model.Order{
		Base:     model.Base{ID: 10},
		UserID:   100,
		PlayerID: &playerID,
		Status:   model.OrderStatusCompleted,
		ItemID:   1,
		OrderNo:  "REVIEW-TEST-DUPLICATE",
	}
	orderRepo.orders[10] = order // Manually add to match mock review data

	reviewSvc := review.NewReviewService(reviewRepo, orderRepo, &fakePlayerRepository{}, &fakeUserRepository{}, &mockReviewReplyRepoForUserReview{})

	router := gin.New()
	router.POST("/user/reviews", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		createReviewHandler(c, reviewSvc)
	})

	// Order 10 already has a review (review ID 1 in mock data has OrderID 10)
	reqBody := review.CreateReviewRequest{
		OrderID: 10,
		Rating:  5,
		Comment: "Duplicate review",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/user/reviews", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}

func TestGetMyReviewsHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	reviewRepo := newMockReviewRepoForUserReview()
	reviewSvc := review.NewReviewService(reviewRepo, newFakeOrderRepository(), &fakePlayerRepository{}, &fakeUserRepository{}, &mockReviewReplyRepoForUserReview{})

	router := gin.New()
	router.GET("/user/reviews/my", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		getMyReviewsHandler(c, reviewSvc)
	})

	req := httptest.NewRequest(http.MethodGet, "/user/reviews/my?page=1&pageSize=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp model.APIResponse[review.MyReviewListResponse]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Fatalf("Expected success=true, got %v", resp.Success)
	}
}

func TestGetMyReviewsHandler_WithPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)

	reviewRepo := newMockReviewRepoForUserReview()
	reviewSvc := review.NewReviewService(reviewRepo, newFakeOrderRepository(), &fakePlayerRepository{}, &fakeUserRepository{}, &mockReviewReplyRepoForUserReview{})

	router := gin.New()
	router.GET("/user/reviews/my", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		getMyReviewsHandler(c, reviewSvc)
	})

	req := httptest.NewRequest(http.MethodGet, "/user/reviews/my?page=2&pageSize=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
}
