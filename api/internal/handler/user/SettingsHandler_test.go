package user

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/handler/testutil"
	"gamelink/internal/model"
	notificationrepo "gamelink/internal/repository/notification"
	userrepo "gamelink/internal/repository/user"
	notificationservice "gamelink/internal/service/notification"
	userservice "gamelink/internal/service/user"
)

type userSettingsTestContext struct {
	Router   *gin.Engine
	DB       *gorm.DB
	User     *model.User
	Settings *userservice.SettingsService
	Notify   *notificationservice.SettingsService
}

func setupUserSettingsTest(t *testing.T) *userSettingsTestContext {
	t.Helper()

	db := testutil.SetupTestDB(t)
	router := testutil.SetupGinTest(t)
	testUser := testutil.CreateAdminUser(t, db, model.RoleUser)

	settingsRepo := userrepo.NewUserSettingsRepository(db)
	notificationSettingsRepo := notificationrepo.NewNotificationSettingRepository(db)

	settingsSvc := userservice.NewSettingsService(settingsRepo)
	notificationSvc := notificationservice.NewSettingsService(notificationSettingsRepo)

	return &userSettingsTestContext{
		Router:   router,
		DB:       db,
		User:     testUser,
		Settings: settingsSvc,
		Notify:   notificationSvc,
	}
}

func (ctx *userSettingsTestContext) registerRoutes() {
	authMiddleware := func(c *gin.Context) {
		c.Set("user_id", ctx.User.ID)
		c.Next()
	}
	userGroup := ctx.Router.Group("/user")
	RegisterSettingsRoutes(userGroup, ctx.Settings, ctx.Notify, authMiddleware)
}

func TestUserSettings_GetAndUpdate(t *testing.T) {
	ctx := setupUserSettingsTest(t)
	ctx.registerRoutes()

	getResp := testutil.MakeRequest(t, ctx.Router, http.MethodGet, "/user/settings", nil)
	testutil.AssertSuccess(t, getResp)

	var getBody map[string]interface{}
	require.NoError(t, json.Unmarshal(getResp.Body.Bytes(), &getBody))
	data := getBody["data"].(map[string]interface{})
	assert.Equal(t, "auto", data["theme"])
	assert.Equal(t, "zh-CN", data["language"])

	payload := map[string]interface{}{
		"theme":    "dark",
		"language": "en-US",
		"notifications": map[string]bool{
			"orderReminder": false,
			"newMessage":    true,
			"systemNotice":  false,
		},
		"privacy": map[string]bool{
			"showOnlineStatus":     false,
			"allowStrangerMessage": true,
		},
	}
	updateResp := testutil.MakeRequest(t, ctx.Router, http.MethodPut, "/user/settings", payload)
	testutil.AssertSuccess(t, updateResp)

	var updateBody map[string]interface{}
	require.NoError(t, json.Unmarshal(updateResp.Body.Bytes(), &updateBody))
	updateData := updateBody["data"].(map[string]interface{})
	assert.Equal(t, "dark", updateData["theme"])
	assert.Equal(t, "en-US", updateData["language"])
}

func TestUserNotificationSettings_GetAndUpdate(t *testing.T) {
	ctx := setupUserSettingsTest(t)
	ctx.registerRoutes()

	getResp := testutil.MakeRequest(t, ctx.Router, http.MethodGet, "/user/notification-settings", nil)
	testutil.AssertSuccess(t, getResp)

	payload := map[string]bool{
		"orderReminder": false,
		"newMessage":    false,
		"systemNotice":  true,
	}
	updateResp := testutil.MakeRequest(t, ctx.Router, http.MethodPut, "/user/notification-settings", payload)
	testutil.AssertSuccess(t, updateResp)

	var updateBody map[string]interface{}
	require.NoError(t, json.Unmarshal(updateResp.Body.Bytes(), &updateBody))
	updateData := updateBody["data"].(map[string]interface{})
	assert.Equal(t, false, updateData["orderReminder"])
	assert.Equal(t, false, updateData["newMessage"])
	assert.Equal(t, true, updateData["systemNotice"])
}
