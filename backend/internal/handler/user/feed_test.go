package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gamelink/internal/model"
	feedrepo "gamelink/internal/repository/feed"
	feedservice "gamelink/internal/service/feed"
)

func setupFeedTest(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&model.Feed{}, &model.FeedImage{}, &model.FeedReport{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	repo := feedrepo.NewFeedRepository(db)
	svc := feedservice.NewService(repo, feedservice.NewDefaultModerationEngine())

	engine := gin.New()
	engine.Use(func(c *gin.Context) {
		c.Set("user_id", uint64(1))
		c.Set("request_id", "trace-123")
		c.Next()
	})
	auth := func(c *gin.Context) { c.Next() }
	RegisterFeedRoutes(engine.Group("/user"), svc, auth)
	return engine, db
}

func TestCreateFeed_Success(t *testing.T) {
	router, _ := setupFeedTest(t)
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
	router, _ := setupFeedTest(t)
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
	router, _ := setupFeedTest(t)
	body := `{"content":"这是违规内容","images":[]}`
	req := httptest.NewRequest(http.MethodPost, "/user/feeds", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w.Code)
	}
}
