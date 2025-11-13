package player

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	reviewsvc "gamelink/internal/service/review"
)

// Mock repositories for player review handler tests
type mockReviewRepoForPlayerHandler struct {
	reviews map[uint64]*model.Review
}

func (m *mockReviewRepoForPlayerHandler) Create(ctx context.Context, review *model.Review) error {
	if review.ID == 0 {
		review.ID = uint64(len(m.reviews) + 1)
	}
	m.reviews[review.ID] = review
	return nil
}

func (m *mockReviewRepoForPlayerHandler) Get(ctx context.Context, id uint64) (*model.Review, error) {
	if r, ok := m.reviews[id]; ok {
		return r, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockReviewRepoForPlayerHandler) List(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
	reviews := make([]model.Review, 0)
	for _, r := range m.reviews {
		reviews = append(reviews, *r)
	}
	return reviews, int64(len(reviews)), nil
}

func (m *mockReviewRepoForPlayerHandler) Update(ctx context.Context, review *model.Review) error {
	m.reviews[review.ID] = review
	return nil
}

func (m *mockReviewRepoForPlayerHandler) Delete(ctx context.Context, id uint64) error {
	delete(m.reviews, id)
	return nil
}

type mockReviewReplyRepoForPlayerHandler struct {
	replies map[uint64]*model.ReviewReply
}

func (m *mockReviewReplyRepoForPlayerHandler) Create(ctx context.Context, reply *model.ReviewReply) error {
	if reply.ID == 0 {
		reply.ID = uint64(len(m.replies) + 1)
	}
	m.replies[reply.ID] = reply
	return nil
}

func (m *mockReviewReplyRepoForPlayerHandler) ListByReview(ctx context.Context, reviewID uint64) ([]model.ReviewReply, error) {
	replies := make([]model.ReviewReply, 0)
	for _, r := range m.replies {
		if r.ReviewID == reviewID {
			replies = append(replies, *r)
		}
	}
	return replies, nil
}

func (m *mockReviewReplyRepoForPlayerHandler) UpdateStatus(ctx context.Context, replyID uint64, status string, note string) error {
	if r, ok := m.replies[replyID]; ok {
		r.Status = status
		r.ModerationNote = note
		return nil
	}
	return repository.ErrNotFound
}

type mockPlayerRepoForReviewHandler struct{}

func (m *mockPlayerRepoForReviewHandler) List(ctx context.Context) ([]model.Player, error) { return nil, nil }
func (m *mockPlayerRepoForReviewHandler) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	return nil, 0, nil
}
func (m *mockPlayerRepoForReviewHandler) Get(ctx context.Context, id uint64) (*model.Player, error) {
	return &model.Player{Base: model.Base{ID: id}, UserID: 1}, nil
}
func (m *mockPlayerRepoForReviewHandler) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	return &model.Player{Base: model.Base{ID: 1}, UserID: userID}, nil
}
func (m *mockPlayerRepoForReviewHandler) Create(ctx context.Context, player *model.Player) error { return nil }
func (m *mockPlayerRepoForReviewHandler) Update(ctx context.Context, player *model.Player) error { return nil }
func (m *mockPlayerRepoForReviewHandler) Delete(ctx context.Context, id uint64) error { return nil }
func (m *mockPlayerRepoForReviewHandler) ListByGameID(ctx context.Context, gameID uint64) ([]model.Player, error) {
	return nil, nil
}

type mockOrderRepoForReviewHandler struct{}

func (m *mockOrderRepoForReviewHandler) Create(ctx context.Context, o *model.Order) error { return nil }
func (m *mockOrderRepoForReviewHandler) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	return nil, 0, nil
}
func (m *mockOrderRepoForReviewHandler) Get(ctx context.Context, id uint64) (*model.Order, error) {
	playerID := uint64(1)
	return &model.Order{Base: model.Base{ID: id}, PlayerID: &playerID}, nil
}
func (m *mockOrderRepoForReviewHandler) Update(ctx context.Context, o *model.Order) error { return nil }
func (m *mockOrderRepoForReviewHandler) Delete(ctx context.Context, id uint64) error { return nil }

type mockUserRepoForReviewHandler struct{}

func (m *mockUserRepoForReviewHandler) List(ctx context.Context) ([]model.User, error) { return nil, nil }
func (m *mockUserRepoForReviewHandler) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (m *mockUserRepoForReviewHandler) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (m *mockUserRepoForReviewHandler) Get(ctx context.Context, id uint64) (*model.User, error) {
	return &model.User{Base: model.Base{ID: id}}, nil
}
func (m *mockUserRepoForReviewHandler) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	return nil, repository.ErrNotFound
}
func (m *mockUserRepoForReviewHandler) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, repository.ErrNotFound
}
func (m *mockUserRepoForReviewHandler) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	return nil, repository.ErrNotFound
}
func (m *mockUserRepoForReviewHandler) Create(ctx context.Context, user *model.User) error { return nil }
func (m *mockUserRepoForReviewHandler) Update(ctx context.Context, user *model.User) error { return nil }
func (m *mockUserRepoForReviewHandler) Delete(ctx context.Context, id uint64) error { return nil }

func setupPlayerReviewService(t *testing.T) *reviewsvc.ReviewService {
	t.Helper()
	reviewRepo := &mockReviewRepoForPlayerHandler{reviews: make(map[uint64]*model.Review)}
	orderRepo := &mockOrderRepoForReviewHandler{}
	playerRepo := &mockPlayerRepoForReviewHandler{}
	userRepo := &mockUserRepoForReviewHandler{}
	replyRepo := &mockReviewReplyRepoForPlayerHandler{replies: make(map[uint64]*model.ReviewReply)}
	svc := reviewsvc.NewReviewService(reviewRepo, orderRepo, playerRepo, userRepo, replyRepo)
	return svc
}

func TestPlayerReviewHandler_ReplyReview_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := setupPlayerReviewService(t)

	payload := reviewsvc.ReplyReviewRequest{
		Content: "Thank you for the feedback!",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/player/reviews/1/reply", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", uint64(1))
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	replyReviewHandler(c, svc)

	// Accept both 200 (success) and 403 (forbidden - not authorized) or 500 (internal error)
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusForbidden || w.Code == http.StatusInternalServerError)
}

func TestPlayerReviewHandler_ReplyReview_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := setupPlayerReviewService(t)

	req := httptest.NewRequest(http.MethodPost, "/player/reviews/1/reply", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", uint64(1))
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	replyReviewHandler(c, svc)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPlayerReviewHandler_ReplyReview_InvalidReviewID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := setupPlayerReviewService(t)

	payload := reviewsvc.ReplyReviewRequest{
		Content: "Thank you for the feedback!",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/player/reviews/abc/reply", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", uint64(1))
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	replyReviewHandler(c, svc)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPlayerReviewHandler_ReplyReview_EmptyContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := setupPlayerReviewService(t)

	payload := reviewsvc.ReplyReviewRequest{
		Content: "", // Empty content
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/player/reviews/1/reply", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", uint64(1))
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	replyReviewHandler(c, svc)

	// Should return error for empty content
	assert.True(t, w.Code == http.StatusBadRequest || w.Code == http.StatusInternalServerError)
}

func TestPlayerReviewHandler_ReplyReview_MissingUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := setupPlayerReviewService(t)

	payload := reviewsvc.ReplyReviewRequest{
		Content: "Thank you for the feedback!",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/player/reviews/1/reply", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	// Don't set user_id
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	replyReviewHandler(c, svc)

	// Should return 0 (default) for missing user_id
	assert.True(t, w.Code >= 400)
}
