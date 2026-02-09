package notification

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/handler/testutil"
	"gamelink/internal/model"
	notificationrepo "gamelink/internal/repository/notification"
	contentservice "gamelink/internal/service/content"
)

type userNotificationTestContext struct {
	Router *gin.Engine
	DB     *gorm.DB
	User   *model.User
	Svc    *contentservice.NotificationService
}

func setupUserNotificationTest(t *testing.T) *userNotificationTestContext {
	t.Helper()

	db := testutil.SetupTestDB(t)
	router := testutil.SetupGinTest(t)
	user := testutil.CreateAdminUser(t, db, model.RoleUser)

	repo := notificationrepo.NewNotificationRepository(db)
	svc := contentservice.NewNotificationService(repo)

	return &userNotificationTestContext{
		Router: router,
		DB:     db,
		User:   user,
		Svc:    svc,
	}
}

func (ctx *userNotificationTestContext) registerRoutes() {
	authMiddleware := func(c *gin.Context) {
		c.Set("user_id", ctx.User.ID)
		c.Next()
	}
	RegisterUserRoutes(ctx.Router, ctx.Svc, authMiddleware)
}

func TestUserNotifications_ListAndRead(t *testing.T) {
	ctx := setupUserNotificationTest(t)
	ctx.registerRoutes()

	readAt := time.Now().Add(-time.Hour)
	events := []model.NotificationEvent{
		{
			UserID:   ctx.User.ID,
			Title:    "Order Update",
			Message:  "Order confirmed",
			Priority: model.NotificationPriorityNormal,
			Channel:  "web",
		},
		{
			UserID:   ctx.User.ID,
			Title:    "System Notice",
			Message:  "Welcome",
			Priority: model.NotificationPriorityLow,
			Channel:  "web",
			ReadAt:   &readAt,
		},
	}
	for i := range events {
		require.NoError(t, ctx.DB.Create(&events[i]).Error)
	}

	listResp := testutil.MakeRequest(t, ctx.Router, http.MethodGet, "/user/notifications", nil)
	testutil.AssertSuccess(t, listResp)

	var listBody model.APIResponse[contentservice.NotificationListResponse]
	require.NoError(t, json.Unmarshal(listResp.Body.Bytes(), &listBody))
	assert.Equal(t, int64(2), listBody.Data.Total)
	assert.Equal(t, int64(1), listBody.Data.UnreadCount)

	unreadResp := testutil.MakeRequest(t, ctx.Router, http.MethodGet, "/user/notifications/unread-count", nil)
	testutil.AssertSuccess(t, unreadResp)

	var unreadBody model.APIResponse[map[string]int64]
	require.NoError(t, json.Unmarshal(unreadResp.Body.Bytes(), &unreadBody))
	assert.Equal(t, int64(1), unreadBody.Data["unread"])

	readPath := testutil.BuildPath("/user/notifications/:id/read", map[string]string{
		"id": testutil.Uint64ToStr(events[0].ID),
	})
	readResp := testutil.MakeRequest(t, ctx.Router, http.MethodPost, readPath, nil)
	testutil.AssertSuccess(t, readResp)

	readAllResp := testutil.MakeRequest(t, ctx.Router, http.MethodPost, "/user/notifications/read-all", nil)
	testutil.AssertSuccess(t, readAllResp)
}
