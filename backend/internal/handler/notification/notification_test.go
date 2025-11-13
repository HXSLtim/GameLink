package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	notificationservice "gamelink/internal/service/notification"
)

// ---- Mock Repository ----

type mockNotificationRepoForHandler struct {
	notifications map[uint64]*model.NotificationEvent
	nextID        uint64
}

func newMockNotificationRepoForHandler() *mockNotificationRepoForHandler {
	return &mockNotificationRepoForHandler{
		notifications: make(map[uint64]*model.NotificationEvent),
		nextID:        1,
	}
}

func (m *mockNotificationRepoForHandler) Create(ctx context.Context, notification *model.NotificationEvent) error {
	notification.ID = m.nextID
	m.nextID++
	m.notifications[notification.ID] = notification
	return nil
}

func (m *mockNotificationRepoForHandler) Get(ctx context.Context, id uint64) (*model.NotificationEvent, error) {
	if notif, ok := m.notifications[id]; ok {
		return notif, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockNotificationRepoForHandler) List(ctx context.Context, opts repository.NotificationListOptions) ([]model.NotificationEvent, int64, error) {
	var result []model.NotificationEvent
	for _, notif := range m.notifications {
		if notif.UserID == opts.UserID {
			result = append(result, *notif)
		}
	}
	return result, int64(len(result)), nil
}

func (m *mockNotificationRepoForHandler) Update(ctx context.Context, notification *model.NotificationEvent) error {
	if _, ok := m.notifications[notification.ID]; ok {
		m.notifications[notification.ID] = notification
		return nil
	}
	return repository.ErrNotFound
}

func (m *mockNotificationRepoForHandler) Delete(ctx context.Context, id uint64) error {
	delete(m.notifications, id)
	return nil
}

func (m *mockNotificationRepoForHandler) CountUnread(ctx context.Context, userID uint64) (int64, error) {
	count := 0
	for _, notif := range m.notifications {
		if notif.UserID == userID && notif.ReadAt == nil {
			count++
		}
	}
	return int64(count), nil
}

func (m *mockNotificationRepoForHandler) MarkRead(ctx context.Context, userID uint64, notificationIDs []uint64) error {
	for _, id := range notificationIDs {
		if notif, ok := m.notifications[id]; ok && notif.UserID == userID {
			now := time.Now()
			notif.ReadAt = &now
		}
	}
	return nil
}

func (m *mockNotificationRepoForHandler) ListByUser(ctx context.Context, opts repository.NotificationListOptions) ([]model.NotificationEvent, int64, error) {
	var result []model.NotificationEvent
	for _, notif := range m.notifications {
		if notif.UserID == opts.UserID {
			result = append(result, *notif)
		}
	}
	return result, int64(len(result)), nil
}

// ---- Setup Helper ----

func setupNotificationService() *notificationservice.Service {
	repo := newMockNotificationRepoForHandler()
	return notificationservice.NewService(repo)
}

// ---- Tests for listNotificationsHandler ----

func TestListNotificationsHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := setupNotificationService()

	router := gin.New()
	router.GET("/notifications", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		listNotificationsHandler(c, svc)
	})

	req := httptest.NewRequest(http.MethodGet, "/notifications?page=1&pageSize=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[*notificationservice.ListResponse]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestListNotificationsHandler_WithUnreadFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := setupNotificationService()

	router := gin.New()
	router.GET("/notifications", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		listNotificationsHandler(c, svc)
	})

	req := httptest.NewRequest(http.MethodGet, "/notifications?page=1&pageSize=20&unread=true", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[*notificationservice.ListResponse]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestListNotificationsHandler_WithPriorityFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := setupNotificationService()

	router := gin.New()
	router.GET("/notifications", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		listNotificationsHandler(c, svc)
	})

	req := httptest.NewRequest(http.MethodGet, "/notifications?page=1&pageSize=20&priority=high&priority=medium", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[*notificationservice.ListResponse]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestListNotificationsHandler_InvalidQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := setupNotificationService()

	router := gin.New()
	router.GET("/notifications", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		listNotificationsHandler(c, svc)
	})

	// Invalid page parameter - should return 400
	req := httptest.NewRequest(http.MethodGet, "/notifications?page=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return 400 for invalid query parameter
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListNotificationsHandler_MissingUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := setupNotificationService()

	router := gin.New()
	router.GET("/notifications", func(c *gin.Context) {
		// Don't set user_id
		listNotificationsHandler(c, svc)
	})

	req := httptest.NewRequest(http.MethodGet, "/notifications", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[*notificationservice.ListResponse]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

// ---- Tests for markNotificationsReadHandler ----

func TestMarkNotificationsReadHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := setupNotificationService()

	router := gin.New()
	router.POST("/notifications/read", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		markNotificationsReadHandler(c, svc)
	})

	body := map[string]interface{}{
		"ids": []uint64{1, 2, 3},
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/notifications/read", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "已标记为已读", resp.Message)
}

func TestMarkNotificationsReadHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := setupNotificationService()

	router := gin.New()
	router.POST("/notifications/read", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		markNotificationsReadHandler(c, svc)
	})

	req := httptest.NewRequest(http.MethodPost, "/notifications/read", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMarkNotificationsReadHandler_EmptyIDs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := setupNotificationService()

	router := gin.New()
	router.POST("/notifications/read", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		markNotificationsReadHandler(c, svc)
	})

	body := map[string]interface{}{
		"ids": []uint64{},
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/notifications/read", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestMarkNotificationsReadHandler_MissingUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := setupNotificationService()

	router := gin.New()
	router.POST("/notifications/read", func(c *gin.Context) {
		// Don't set user_id
		markNotificationsReadHandler(c, svc)
	})

	body := map[string]interface{}{
		"ids": []uint64{1, 2},
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/notifications/read", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ---- Tests for unreadCountHandler ----

func TestUnreadCountHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := setupNotificationService()

	router := gin.New()
	router.GET("/notifications/unread-count", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		unreadCountHandler(c, svc)
	})

	req := httptest.NewRequest(http.MethodGet, "/notifications/unread-count", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[map[string]int64]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Data)
	assert.Contains(t, resp.Data, "unread")
}

func TestUnreadCountHandler_MissingUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := setupNotificationService()

	router := gin.New()
	router.GET("/notifications/unread-count", func(c *gin.Context) {
		// Don't set user_id
		unreadCountHandler(c, svc)
	})

	req := httptest.NewRequest(http.MethodGet, "/notifications/unread-count", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[map[string]int64]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestUnreadCountHandler_WithData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newMockNotificationRepoForHandler()
	// Add some notifications
	repo.Create(context.Background(), &model.NotificationEvent{
		UserID:    100,
		Title:     "Test 1",
		Message:   "Message 1",
		Priority:  model.NotificationPriorityHigh,
		ReadAt:    nil,
	})
	repo.Create(context.Background(), &model.NotificationEvent{
		UserID:    100,
		Title:     "Test 2",
		Message:   "Message 2",
		Priority:  model.NotificationPriorityNormal,
		ReadAt:    nil,
	})

	svc := notificationservice.NewService(repo)

	router := gin.New()
	router.GET("/notifications/unread-count", func(c *gin.Context) {
		c.Set("user_id", uint64(100))
		unreadCountHandler(c, svc)
	})

	req := httptest.NewRequest(http.MethodGet, "/notifications/unread-count", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[map[string]int64]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, int64(2), resp.Data["unread"])
}

// ---- Tests for helper functions ----

func TestGetUserIDFromContext_Valid(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("user_id", uint64(123))

	userID := getUserIDFromContext(c)
	assert.Equal(t, uint64(123), userID)
}

func TestGetUserIDFromContext_Missing(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	userID := getUserIDFromContext(c)
	assert.Equal(t, uint64(0), userID)
}

func TestGetUserIDFromContext_InvalidType(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("user_id", "not-a-uint64")

	userID := getUserIDFromContext(c)
	assert.Equal(t, uint64(0), userID)
}

func TestRespondJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("request_id", "trace-123")

	payload := model.APIResponse[string]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    "test",
	}

	respondJSON(c, http.StatusOK, payload)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var resp model.APIResponse[string]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "trace-123", resp.TraceID)
}

func TestRespondError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	respondError(c, http.StatusBadRequest, "test error")

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp model.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "test error", resp.Message)
}
