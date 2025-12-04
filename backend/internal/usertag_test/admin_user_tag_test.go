package usertag_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"gamelink/internal/handler/admin"
	"gamelink/internal/model"
	userrepo "gamelink/internal/repository/user"
	"gamelink/internal/service/user"
	"gamelink/pkg/cache"
	"gamelink/pkg/testutil"
)

func migrateUserModels(t *testing.T, db *gorm.DB) {
	models := []interface{}{
		&model.User{},
		&model.UserLoginHistory{},
		&model.UserBehavior{},
		&model.UserTag{},
		&model.UserTagRelation{},
	}
	for _, m := range models {
		assert.NoError(t, db.AutoMigrate(m))
	}
}

func doJSON(router *gin.Engine, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func fakeAuthMiddleware(userID uint64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userId", uint64ToStr(userID))
		c.Next()
	}
}

func uint64ToStr(v uint64) string {
	return strconv.FormatUint(v, 10)
}

type apiResp[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// TestAdminUserTagManagement 测试用户标签管理全流程
func TestAdminUserTagManagement(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateUserModels(t, db)

	// 创建测试用户
	userRepo := userrepo.NewUserRepository(db)
	adminUser := &model.User{
		Name:         "Admin",
		Email:        "admin@example.com",
		Phone:        "18800000001",
		PasswordHash: "hashed",
		Role:         model.RoleAdmin,
		Status:       model.UserStatusActive,
	}
	assert.NoError(t, userRepo.Create(context.Background(), adminUser))

	// 创建测试用户
	testUser := &model.User{
		Name:         "TestUser",
		Email:        "user@example.com",
		Phone:        "18800000002",
		PasswordHash: "hashed",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
	}
	assert.NoError(t, userRepo.Create(context.Background(), testUser))

	// 初始化服务
	tagRepo := userrepo.NewUserTagRepository(db)
	tagService := user.NewUserTagService(tagRepo, userRepo, cache.NewMemory())

	// 创建路由
	router := gin.New()
	api := router.Group("/api/v1/admin")
	api.Use(fakeAuthMiddleware(adminUser.ID))

	// 使用RegisterTagRoutes注册所有路由
	admin.RegisterTagRoutes(api, tagService)

	// 1. 创建标签
	tag1Payload := map[string]interface{}{
		"name":        "VIP",
		"color":       "#FF6B6B",
		"description": "VIP用户",
	}
	resp := doJSON(router, http.MethodPost, "/api/v1/admin/user-tags", tag1Payload, "")
	assert.Equal(t, http.StatusOK, resp.Code, "创建标签应该成功")
	var tag1Resp apiResp[admin.TagResponse]
	json.Unmarshal(resp.Body.Bytes(), &tag1Resp)
	assert.Equal(t, "VIP", tag1Resp.Data.Name)
	assert.Equal(t, "#FF6B6B", tag1Resp.Data.Color)
	vipTagID := tag1Resp.Data.ID

	// 创建第二个标签
	tag2Payload := map[string]interface{}{
		"name":        "活跃",
		"color":       "#4CAF50",
		"description": "活跃用户",
	}
	resp = doJSON(router, http.MethodPost, "/api/v1/admin/user-tags", tag2Payload, "")
	assert.Equal(t, http.StatusOK, resp.Code, "创建第二个标签应该成功")
	var tag2Resp apiResp[admin.TagResponse]
	json.Unmarshal(resp.Body.Bytes(), &tag2Resp)

	// 2. 获取标签列表
	resp = doJSON(router, http.MethodGet, "/api/v1/admin/user-tags", nil, "")
	assert.Equal(t, http.StatusOK, resp.Code)
	var tagsResp apiResp[[]admin.TagResponse]
	json.Unmarshal(resp.Body.Bytes(), &tagsResp)
	assert.Len(t, tagsResp.Data, 2, "应该返回两个标签")

	// 3. 更新标签
	updatePayload := map[string]interface{}{
		"name":        "VIP用户",
		"color":       "#FF0000",
		"description": "高级VIP用户",
	}
	resp = doJSON(router, http.MethodPut, "/api/v1/admin/user-tags/"+uint64ToStr(vipTagID), updatePayload, "")
	assert.Equal(t, http.StatusOK, resp.Code, "更新标签应该成功")
	var updatedTag apiResp[admin.TagResponse]
	json.Unmarshal(resp.Body.Bytes(), &updatedTag)
	assert.Equal(t, "VIP用户", updatedTag.Data.Name)
	assert.Equal(t, "#FF0000", updatedTag.Data.Color)

	// 4. 给用户添加标签
	resp = doJSON(router, http.MethodPost, "/api/v1/admin/users/"+uint64ToStr(testUser.ID)+"/tags", map[string]interface{}{"tagId": vipTagID}, "")
	assert.Equal(t, http.StatusOK, resp.Code, "添加用户标签应该成功")

	// 5. 获取用户的所有标签
	resp = doJSON(router, http.MethodGet, "/api/v1/admin/users/"+uint64ToStr(testUser.ID)+"/tags", nil, "")
	assert.Equal(t, http.StatusOK, resp.Code)
	var userTagsResp apiResp[[]admin.TagResponse]
	json.Unmarshal(resp.Body.Bytes(), &userTagsResp)
	assert.Len(t, userTagsResp.Data, 1, "用户应该有一个标签")

	// 6. 移除用户标签
	resp = doJSON(router, http.MethodDelete, "/api/v1/admin/users/"+uint64ToStr(testUser.ID)+"/tags/"+uint64ToStr(vipTagID), nil, "")
	assert.Equal(t, http.StatusOK, resp.Code, "移除标签应该成功")

	// 7. 删除标签
	resp = doJSON(router, http.MethodDelete, "/api/v1/admin/user-tags/"+uint64ToStr(vipTagID), nil, "")
	assert.Equal(t, http.StatusOK, resp.Code, "删除标签应该成功")

	// 8. 验证标签已被删除
	resp = doJSON(router, http.MethodGet, "/api/v1/admin/user-tags", nil, "")
	json.Unmarshal(resp.Body.Bytes(), &tagsResp)
	assert.Len(t, tagsResp.Data, 1, "应该只剩一个标签")
}
