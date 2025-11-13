package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	feedservice "gamelink/internal/service/feed"
)

// mockFeedRepository for testing
type mockFeedRepository struct {
	feeds map[uint64]*model.Feed
}

func newMockFeedRepository() *mockFeedRepository {
	return &mockFeedRepository{feeds: make(map[uint64]*model.Feed)}
}

func (m *mockFeedRepository) Create(ctx context.Context, feed *model.Feed) error {
	if feed.ID == 0 {
		feed.ID = uint64(len(m.feeds) + 1)
	}
	m.feeds[feed.ID] = feed
	return nil
}

func (m *mockFeedRepository) Get(ctx context.Context, id uint64) (*model.Feed, error) {
	if f, ok := m.feeds[id]; ok {
		return f, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockFeedRepository) List(ctx context.Context, opts repository.FeedListOptions) ([]model.Feed, error) {
	feeds := make([]model.Feed, 0)
	for _, f := range m.feeds {
		feeds = append(feeds, *f)
	}
	return feeds, nil
}

func (m *mockFeedRepository) UpdateModeration(ctx context.Context, feedID uint64, status model.FeedModerationStatus, note string, manual bool) error {
	if f, ok := m.feeds[feedID]; ok {
		f.ModerationStatus = status
		f.ModerationNote = note
		return nil
	}
	return repository.ErrNotFound
}

func (m *mockFeedRepository) CreateReport(ctx context.Context, report *model.FeedReport) error {
	return nil
}

func setupFeedTest(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	
	repo := newMockFeedRepository()
	svc := feedservice.NewService(repo, feedservice.NewDefaultModerationEngine())

	engine := gin.New()
	engine.Use(func(c *gin.Context) {
		c.Set("user_id", uint64(1))
		c.Set("request_id", "trace-123")
		c.Next()
	})
	auth := func(c *gin.Context) { c.Next() }
	RegisterFeedRoutes(engine.Group("/user"), svc, auth)
	return engine
}

func TestCreateFeed_Success(t *testing.T) {
	router := setupFeedTest(t)
	body := `{"content":"今天天气真好","visibility":"public","images":[{"url":"https://img/1.png","sizeBytes":1024}]}`
	req := httptest.NewRequest(http.MethodPost, "/user/feeds", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Success bool `json:"success"`
		Data    struct {
			ModerationStatus string `json:"moderationStatus"`
			Content          string `json:"content"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !resp.Success {
		t.Fatalf("expected success")
	}
	if resp.Data.ModerationStatus == "" {
		t.Fatalf("expected moderation status populated")
	}
}

func TestCreateFeed_TooManyImages(t *testing.T) {
	router := setupFeedTest(t)
	images := make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		images = append(images, fmt.Sprintf(`{"url":"https://img/%d.png","sizeBytes":1024}`, i))
	}
	body := `{"content":"hello","images":[` + strings.Join(images, ",") + `]}`
	req := httptest.NewRequest(http.MethodPost, "/user/feeds", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w.Code)
	}
}

func TestCreateFeed_SensitiveContent(t *testing.T) {
	router := setupFeedTest(t)
	body := `{"content":"这是违规内容","images":[]}`
	req := httptest.NewRequest(http.MethodPost, "/user/feeds", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w.Code)
	}
}
