package integration

import (
	"encoding/json"
	"net/http"
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
	userTagHandler := adminhandler.NewUserTagHandler(tagService)
	router := gin.New()
	api := router.Group("/api/v1/admin")
	api.Use(fakeAuthMiddleware(adminUser.ID))

	// 标签管理路由
	api.POST("/tags", userTagHandler.CreateTag)
	api.GET("/tags", userTagHandler.ListTags)
	api.PUT("/tags/:id", userTagHandler.UpdateTag)
	api.DELETE("/tags/:id", userTagHandler.DeleteTag)
	api.POST("/users/:userId/tags/:tagId", userTagHandler.AddTagToUser)
	api.DELETE("/users/:userId/tags/:tagId", userTagHandler.RemoveTagFromUser)
	api.PUT("/users/:userId/tags", userTagHandler.BatchSetUserTags)
	api.GET("/users/:userId/tags", userTagHandler.GetUserTags)

	// 1. 创建标签
	tag1Payload := map[string]interface{}{
		"name":        "VIP",
		"color":       "#FF6B6B",
		"description": "VIP用户",
	}
	resp := doJSON(router, http.MethodPost, "/api/v1/admin/tags", tag1Payload, "")
	assert.Equal(t, http.StatusOK, resp.Code, "创建标签应该成功")
	var tag1Resp apiResp[adminhandler.TagResponse]
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
	resp = doJSON(router, http.MethodPost, "/api/v1/admin/tags", tag2Payload, "")
	assert.Equal(t, http.StatusOK, resp.Code, "创建第二个标签应该成功")
	var tag2Resp apiResp[adminhandler.TagResponse]
	json.Unmarshal(resp.Body.Bytes(), &tag2Resp)
	activeTagID := tag2Resp.Data.ID

	// 2. 获取标签列表
	resp = doJSON(router, http.MethodGet, "/api/v1/admin/tags", nil, "")
	assert.Equal(t, http.StatusOK, resp.Code)
	var tagsResp apiResp[[]adminhandler.TagResponse]
	json.Unmarshal(resp.Body.Bytes(), &tagsResp)
	assert.Len(t, tagsResp.Data, 2, "应该返回两个标签")

	// 3. 更新标签
	updatePayload := map[string]interface{}{
		"name":        "VIP用户",
		"color":       "#FF0000",
		"description": "高级VIP用户",
	}
	resp = doJSON(router, http.MethodPut, "/api/v1/admin/tags/"+strconv.FormatUint(vipTagID, 10), updatePayload, "")
	assert.Equal(t, http.StatusOK, resp.Code, "更新标签应该成功")
	var updatedTag apiResp[adminhandler.TagResponse]
	json.Unmarshal(resp.Body.Bytes(), &updatedTag)
	assert.Equal(t, "VIP用户", updatedTag.Data.Name)
	assert.Equal(t, "#FF0000", updatedTag.Data.Color)

	// 4. 给用户添加标签
	resp = doJSON(router, http.MethodPost, "/api/v1/admin/users/"+strconv.FormatUint(testUser.ID, 10)+"/tags/"+strconv.FormatUint(vipTagID, 10), nil, "")
	assert.Equal(t, http.StatusOK, resp.Code, "添加用户标签应该成功")

	// 5. 给同一用户添加第二个标签
	resp = doJSON(router, http.MethodPost, "/api/v1/admin/users/"+strconv.FormatUint(testUser.ID, 10)+"/tags/"+strconv.FormatUint(activeTagID, 10), nil, "")
	assert.Equal(t, http.StatusOK, resp.Code)

	// 6. 获取用户的所有标签
	resp = doJSON(router, http.MethodGet, "/api/v1/admin/users/"+strconv.FormatUint(testUser.ID, 10)+"/tags", nil, "")
	assert.Equal(t, http.StatusOK, resp.Code)
	var userTagsResp apiResp[[]adminhandler.TagResponse]
	json.Unmarshal(resp.Body.Bytes(), &userTagsResp)
	assert.Len(t, userTagsResp.Data, 2, "用户应该有两个标签")

	// 7. 批量设置用户标签（覆盖式）
	batchTagsPayload := map[string]interface{}{
		"tagIds": []uint64{activeTagID}, // 只保留活跃标签
	}
	resp = doJSON(router, http.MethodPut, "/api/v1/admin/users/"+strconv.FormatUint(testUser.ID, 10)+"/tags", batchTagsPayload, "")
	assert.Equal(t, http.StatusOK, resp.Code, "批量设置标签应该成功")

	// 8. 验证标签已被覆盖
	resp = doJSON(router, http.MethodGet, "/api/v1/admin/users/"+strconv.FormatUint(testUser.ID, 10)+"/tags", nil, "")
	json.Unmarshal(resp.Body.Bytes(), &userTagsResp)
	assert.Len(t, userTagsResp.Data, 1, "用户应该只剩一个标签")
	assert.Equal(t, "活跃", userTagsResp.Data[0].Name)

	// 9. 移除用户标签
	resp = doJSON(router, http.MethodDelete, "/api/v1/admin/users/"+strconv.FormatUint(testUser.ID, 10)+"/tags/"+strconv.FormatUint(activeTagID, 10), nil, "")
	assert.Equal(t, http.StatusOK, resp.Code, "移除标签应该成功")

	// 10. 验证标签已移除
	resp = doJSON(router, http.MethodGet, "/api/v1/admin/users/"+strconv.FormatUint(testUser.ID, 10)+"/tags", nil, "")
	json.Unmarshal(resp.Body.Bytes(), &userTagsResp)
	assert.Len(t, userTagsResp.Data, 0, "用户应该没有标签")

	// 11. 删除标签
	resp = doJSON(router, http.MethodDelete, "/api/v1/admin/tags/"+strconv.FormatUint(vipTagID, 10), nil, "")
	assert.Equal(t, http.StatusOK, resp.Code, "删除标签应该成功")

	// 12. 验证标签已被删除
	resp = doJSON(router, http.MethodGet, "/api/v1/admin/tags", nil, "")
	json.Unmarshal(resp.Body.Bytes(), &tagsResp)
	assert.Len(t, tagsResp.Data, 1, "应该只剩一个标签")
}

// TestAdminUserBatchOperations 测试用户批量操作
func TestAdminUserBatchOperations(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateUserModels(t, db)

	// 创建管理员用户
	userRepo := userrepo.NewUserRepository(db)
	adminUser := &model.User{
		Name:         "Admin",
		Email:        "admin@example.com",
		Phone:        "18800000001",
		PasswordHash: "hashed",
		Role:         model.RoleAdmin,
		Status:       model.UserStatusActive,
		Points:       1000,
	}
	assert.NoError(t, userRepo.Create(context.Background(), adminUser))

	// 创建测试用户
	users := make([]*model.User, 5)
	for i := 0; i < 5; i++ {
		users[i] = &model.User{
			Name:         "User" + string(rune('A'+i)),
			Email:        "user" + string(rune('A'+i)) + "@example.com",
			Phone:        "1880000000" + string(rune('1'+i)),
			PasswordHash: "hashed",
			Role:         model.RoleUser,
			Status:       model.UserStatusActive,
			Points:       100,
		}
		assert.NoError(t, userRepo.Create(context.Background(), users[i]))
	}

	// 初始化批量操作服务
	tagRepo := userrepo.NewUserTagRepository(db)
	batchService := user.NewBatchOperationService(db, userRepo, tagRepo)

	// 创建路由
	batchHandler := adminhandler.NewBatchOperationHandler(batchService)
	router := gin.New()
	api := router.Group("/api/v1/admin")
	api.Use(fakeAuthMiddleware(adminUser.ID))

	// 批量操作路由
	api.POST("/users/batch/role", batchHandler.BatchUpdateRole)
	api.POST("/users/batch/status", batchHandler.BatchUpdateStatus)
	api.POST("/users/batch/delete", batchHandler.BatchDeleteUsers)
	api.POST("/users/batch/points", batchHandler.BatchAddPoints)
	api.POST("/users/batch/notification", batchHandler.BatchSendNotification)

	userIDs := []uint64{users[0].ID, users[1].ID, users[2].ID}

	// 1. 批量更新用户角色
	rolePayload := map[string]interface{}{
		"userIds": userIDs,
		"role":    "player",
	}
	resp := doJSON(router, http.MethodPost, "/api/v1/admin/users/batch/role", rolePayload, "")
	assert.Equal(t, http.StatusOK, resp.Code, "批量更新角色应该成功")
	var roleResp apiResp[adminhandler.BatchResponse]
	json.Unmarshal(resp.Body.Bytes(), &roleResp)
	assert.Equal(t, 3, roleResp.Data.SuccessCount)
	assert.Equal(t, 0, roleResp.Data.FailedCount)

	// 验证角色已更新
	updatedUser, _ := userRepo.Get(context.Background(), users[0].ID)
	assert.Equal(t, model.RolePlayer, updatedUser.Role)

	// 2. 批量更新用户状态
	statusPayload := map[string]interface{}{
		"userIds": userIDs,
		"status":  "suspended",
	}
	resp = doJSON(router, http.MethodPost, "/api/v1/admin/users/batch/status", statusPayload, "")
	assert.Equal(t, http.StatusOK, resp.Code, "批量更新状态应该成功")

	// 验证状态已更新
	updatedUser, _ = userRepo.Get(context.Background(), users[0].ID)
	assert.Equal(t, model.UserStatusSuspended, updatedUser.Status)

	// 3. 批量增加积分
	pointsPayload := map[string]interface{}{
		"userIds": userIDs,
		"points":  50,
		"reason":  "测试奖励",
	}
	resp = doJSON(router, http.MethodPost, "/api/v1/admin/users/batch/points", pointsPayload, "")
	assert.Equal(t, http.StatusOK, resp.Code, "批量增加积分应该成功")

	// 验证积分已增加
	updatedUser, _ = userRepo.Get(context.Background(), users[0].ID)
	assert.Equal(t, int64(150), updatedUser.Points) // 原100 + 50

	// 4. 批量发送通知
	notificationPayload := map[string]interface{}{
		"userIds": userIDs,
		"title":   "测试通知",
		"content": "这是一条测试通知",
	}
	resp = doJSON(router, http.MethodPost, "/api/v1/admin/users/batch/notification", notificationPayload, "")
	assert.Equal(t, http.StatusOK, resp.Code, "批量发送通知应该成功")

	// 5. 批量删除用户（只删除两个）
	deletePayload := map[string]interface{}{
		"userIds": []uint64{users[3].ID, users[4].ID},
	}
	resp = doJSON(router, http.MethodPost, "/api/v1/admin/users/batch/delete", deletePayload, "")
	assert.Equal(t, http.StatusOK, resp.Code, "批量删除应该成功")

	// 验证用户已删除
	_, err := userRepo.Get(context.Background(), users[3].ID)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

// TestAdminUserTagWithCache 测试标签缓存功能
func TestAdminUserTagWithCache(t *testing.T) {
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

	// 初始化缓存
	memCache := cache.NewMemory()

	// 初始化服务
	tagRepo := userrepo.NewUserTagRepository(db)
	tagService := user.NewUserTagService(tagRepo, userRepo, memCache)

	// 创建路由
	userTagHandler := adminhandler.NewUserTagHandler(tagService)
	router := gin.New()
	api := router.Group("/api/v1/admin")
	api.Use(fakeAuthMiddleware(adminUser.ID))

	api.POST("/tags", userTagHandler.CreateTag)
	api.GET("/tags", userTagHandler.ListTags)

	// 1. 创建标签
	tagPayload := map[string]interface{}{
		"name":        "VIP",
		"color":       "#FF6B6B",
		"description": "VIP用户",
	}
	resp := doJSON(router, http.MethodPost, "/api/v1/admin/tags", tagPayload, "")
	assert.Equal(t, http.StatusCreated, resp.Code)
	var tagResp apiResp[adminhandler.Tag]
	json.Unmarshal(resp.Body.Bytes(), &tagResp)
	tagID := tagResp.Data.ID

	// 2. 第一次获取标签列表（缓存未命中，从数据库读取）
	resp = doJSON(router, http.MethodGet, "/api/v1/admin/tags", nil, "")
	assert.Equal(t, http.StatusOK, resp.Code)

	// 3. 验证缓存已写入
	cacheKey := "user_tags:list"
	cached, ok, _ := memCache.Get(context.Background(), cacheKey)
	assert.True(t, ok, "缓存应该存在")
	assert.NotEmpty(t, cached, "缓存数据不应该为空")

	// 4. 创建新标签，应该清除缓存
	tag2Payload := map[string]interface{}{
		"name":        "活跃",
		"color":       "#4CAF50",
		"description": "活跃用户",
	}
	resp = doJSON(router, http.MethodPost, "/api/v1/admin/tags", tag2Payload, "")
	assert.Equal(t, http.StatusCreated, resp.Code)

	// 5. 验证缓存已被清除
	_, ok, _ = memCache.Get(context.Background(), cacheKey)
	assert.False(t, ok, "创建标签后缓存应该被清除")

	// 6. 再次获取标签列表，应该重新缓存
	resp = doJSON(router, http.MethodGet, "/api/v1/admin/tags", nil, "")
	assert.Equal(t, http.StatusOK, resp.Code)
	_, ok, _ = memCache.Get(context.Background(), cacheKey)
	assert.True(t, ok, "重新获取后缓存应该存在")

	// 7. 更新标签，应该清除缓存
	updatePayload := map[string]interface{}{
		"name":        "VIP更新",
		"color":       "#FF0000",
		"description": "VIP用户更新",
	}
	resp = doJSON(router, http.MethodPut, "/api/v1/admin/tags/"+uint2str(tagID), updatePayload, "")
	assert.Equal(t, http.StatusOK, resp.Code)
	_, ok, _ = memCache.Get(context.Background(), cacheKey)
	assert.False(t, ok, "更新标签后缓存应该被清除")

	// 8. 删除标签，应该清除缓存
	resp = doJSON(router, http.MethodDelete, "/api/v1/admin/tags/"+uint2str(tagID), nil, "")
	assert.Equal(t, http.StatusOK, resp.Code)
	_, ok, _ = memCache.Get(context.Background(), cacheKey)
	assert.False(t, ok, "删除标签后缓存应该被清除")
}

// TestAdminUserBatchOperationsErrorCases 测试批量操作错误场景
func TestAdminUserBatchOperationsErrorCases(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateUserModels(t, db)

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

	// 初始化批量操作服务
	tagRepo := userrepo.NewUserTagRepository(db)
	batchService := user.NewBatchOperationService(db, userRepo, tagRepo)

	batchHandler := adminhandler.NewBatchOperationHandler(batchService)
	router := gin.New()
	api := router.Group("/api/v1/admin")
	api.Use(fakeAuthMiddleware(adminUser.ID))

	api.POST("/users/batch/role", batchHandler.BatchUpdateRole)

	// 测试1: 空的用户ID列表
	payload := map[string]interface{}{
		"userIds": []uint64{},
		"role":    "player",
	}
	resp := doJSON(router, http.MethodPost, "/api/v1/admin/users/batch/role", payload, "")
	assert.Equal(t, http.StatusBadRequest, resp.Code, "空用户列表应该返回400")

	// 测试2: 用户ID超过1000个限制
	largeUserIDs := make([]uint64, 1001)
	for i := 0; i < 1001; i++ {
		largeUserIDs[i] = uint64(i + 1)
	}
	payload = map[string]interface{}{
		"userIds": largeUserIDs,
		"role":    "player",
	}
	resp = doJSON(router, http.MethodPost, "/api/v1/admin/users/batch/role", payload, "")
	assert.Equal(t, http.StatusBadRequest, resp.Code, "用户过多应该返回400")

	// 测试3: 不存在的用户ID
	payload = map[string]interface{}{
		"userIds": []uint64{99999, 99998},
		"role":    "player",
	}
	resp = doJSON(router, http.MethodPost, "/api/v1/admin/users/batch/role", payload, "")
	assert.Equal(t, http.StatusOK, resp.Code)
	var batchResp apiResp[adminhandler.BatchResponse]
	json.Unmarshal(resp.Body.Bytes(), &batchResp)
	assert.Equal(t, 0, batchResp.Data.SuccessCount)
	assert.Equal(t, 2, batchResp.Data.FailedCount)
}

// 辅助函数
func migrateUserModels(t *testing.T, db *gorm.DB) {
