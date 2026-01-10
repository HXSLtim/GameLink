// Package admin provides comprehensive tests for UserTag handlers.
package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/handler/testutil"
	"gamelink/internal/model"
	userrepository "gamelink/internal/repository/user"
	"gamelink/internal/service/integration"
	userservice "gamelink/internal/service/user"
	"gamelink/pkg/cache"
)

// ============================================================================
// Test Context & Setup
// ============================================================================

// UserTagTestContext provides test context for user tag handler tests.
type UserTagTestContext struct {
	Router     *gin.Engine
	DB         *gorm.DB
	Service    *userservice.UserTagService
	AdminUser  *model.User
	AdminToken string
}

// SetupUserTagTest initializes test environment for user tag handler tests.
func SetupUserTagTest(t *testing.T) *UserTagTestContext {
	t.Helper()

	// Setup test database
	db := testutil.SetupTestDB(t)

	// Create repositories
	userRepo := userrepository.NewUserRepository(db)
	tagRepo := userrepository.NewUserTagRepository(db)
	c := cache.NewMemory()

	// Create service
	svc := userservice.NewUserTagService(tagRepo, userRepo, c)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create super admin user
	adminUser := testutil.CreateSuperAdmin(t, db)
	adminToken := testutil.GenerateTestToken(adminUser.ID)

	return &UserTagTestContext{
		Router:     router,
		DB:         db,
		Service:    svc,
		AdminUser:  adminUser,
		AdminToken: adminToken,
	}
}

// RegisterTagRoutes registers user tag routes for testing.
func (ctx *UserTagTestContext) RegisterTagRoutes() {
	RegisterTagRoutes(ctx.Router.Group("/api/v1/admin"), ctx.Service)
}

// Helper to make authenticated request
func (ctx *UserTagTestContext) makeRequest(method, path string, body interface{}) *httptest.ResponseRecorder {
	var bodyReader *bytes.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		require.NoError(nil, err, "Failed to marshal request body")
		bodyReader = bytes.NewReader(jsonData)
	} else {
		bodyReader = bytes.NewReader([]byte{})
	}

	req, err := http.NewRequest(method, path, bodyReader)
	require.NoError(nil, err, "Failed to create request")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ctx.AdminToken))

	w := httptest.NewRecorder()
	ctx.Router.ServeHTTP(w, req)

	return w
}

// Helper to create a test tag via service
func (ctx *UserTagTestContext) createTestTag(t *testing.T, name, color, description string) *model.UserTag {
	t.Helper()
	return integration.CreateTestUserTag(t, ctx.DB, name, color, description)
}

// ============================================================================
// CreateTagHandler Tests
// ============================================================================

func TestUserTagHandler_CreateTag_Success(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	payload := map[string]interface{}{
		"name":        "VIP Customer",
		"color":       "#FF5733",
		"description": "High value customer",
	}

	w := ctx.makeRequest("POST", "/api/v1/admin/user-tags", payload)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["id"])
	assert.Equal(t, "VIP Customer", data["name"])
	assert.Equal(t, "#FF5733", data["color"])
	assert.Equal(t, "High value customer", data["description"])
}

func TestUserTagHandler_CreateTag_DuplicateName(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	// Create first tag
	ctx.createTestTag(t, "VIP Customer", "#FF5733", "First tag")

	// Try to create duplicate
	payload := map[string]interface{}{
		"name":        "VIP Customer",
		"color":       "#00FF00",
		"description": "Duplicate tag",
	}

	w := ctx.makeRequest("POST", "/api/v1/admin/user-tags", payload)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response["success"].(bool))
	assert.Contains(t, response["message"].(string), "标签名称已存在")
}

func TestUserTagHandler_CreateTag_ValidationError(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	tests := []struct {
		name        string
		payload     map[string]interface{}
		expectedMsg string
	}{
		{
			name: "missing name",
			payload: map[string]interface{}{
				"color":       "#FF5733",
				"description": "Test",
			},
			expectedMsg: "参数验证失败",
		},
		{
			name: "missing color",
			payload: map[string]interface{}{
				"name":        "Test Tag",
				"description": "Test",
			},
			expectedMsg: "参数验证失败",
		},
		{
			name: "invalid color format",
			payload: map[string]interface{}{
				"name":        "Test Tag",
				"color":       "invalid",
				"description": "Test",
			},
			expectedMsg: "颜色格式不正确",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := ctx.makeRequest("POST", "/api/v1/admin/user-tags", tt.payload)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.False(t, response["success"].(bool))
		})
	}
}

// ============================================================================
// ListTagsHandler Tests
// ============================================================================

func TestUserTagHandler_ListTags_Success(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	// Create test tags
	ctx.createTestTag(t, "VIP", "#FF5733", "VIP Customer")
	ctx.createTestTag(t, "Blacklist", "#000000", "Banned users")
	ctx.createTestTag(t, "New User", "#00FF00", "New registrations")

	w := ctx.makeRequest("GET", "/api/v1/admin/user-tags", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].([]interface{})
	assert.Len(t, data, 3)
}

func TestUserTagHandler_ListTags_Empty(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	w := ctx.makeRequest("GET", "/api/v1/admin/user-tags", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].([]interface{})
	assert.Len(t, data, 0)
}

// ============================================================================
// GetTagHandler Tests
// ============================================================================

func TestUserTagHandler_GetTag_Success(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	tag := ctx.createTestTag(t, "VIP", "#FF5733", "VIP Customer")

	w := ctx.makeRequest("GET", fmt.Sprintf("/api/v1/admin/user-tags/%d", tag.ID), nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, tag.ID, uint64(data["id"].(float64)))
	assert.Equal(t, "VIP", data["name"])
	assert.Equal(t, "#FF5733", data["color"])
}

func TestUserTagHandler_GetTag_NotFound(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	w := ctx.makeRequest("GET", "/api/v1/admin/user-tags/999999", nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response["success"].(bool))
}

func TestUserTagHandler_GetTag_InvalidID(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	w := ctx.makeRequest("GET", "/api/v1/admin/user-tags/invalid", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response["success"].(bool))
	assert.Contains(t, response["message"].(string), "无效的标签ID")
}

// ============================================================================
// UpdateTagHandler Tests
// ============================================================================

func TestUserTagHandler_UpdateTag_Success(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	tag := ctx.createTestTag(t, "VIP", "#FF5733", "VIP Customer")

	payload := map[string]interface{}{
		"name":        "VIP Gold",
		"color":       "#FFD700",
		"description": "Gold tier VIP customers",
	}

	w := ctx.makeRequest("PUT", fmt.Sprintf("/api/v1/admin/user-tags/%d", tag.ID), payload)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "VIP Gold", data["name"])
	assert.Equal(t, "#FFD700", data["color"])
	assert.Equal(t, "Gold tier VIP customers", data["description"])

	// Verify in database
	var updatedTag model.UserTag
	err = ctx.DB.First(&updatedTag, tag.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "VIP Gold", updatedTag.Name)
}

func TestUserTagHandler_UpdateTag_NotFound(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	payload := map[string]interface{}{
		"name":  "Updated",
		"color": "#FF5733",
	}

	w := ctx.makeRequest("PUT", "/api/v1/admin/user-tags/999999", payload)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response["success"].(bool))
}

func TestUserTagHandler_UpdateTag_InvalidID(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	payload := map[string]interface{}{
		"name": "Test",
	}

	w := ctx.makeRequest("PUT", "/api/v1/admin/user-tags/invalid", payload)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["message"].(string), "无效的标签ID")
}

// ============================================================================
// DeleteTagHandler Tests
// ============================================================================

func TestUserTagHandler_DeleteTag_Success(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	tag := ctx.createTestTag(t, "VIP", "#FF5733", "VIP Customer")

	w := ctx.makeRequest("DELETE", fmt.Sprintf("/api/v1/admin/user-tags/%d", tag.ID), nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	// Verify deletion (should be soft deleted)
	var deletedTag model.UserTag
	err = ctx.DB.Unscoped().First(&deletedTag, tag.ID).Error
	require.NoError(t, err)
	assert.NotZero(t, deletedTag.DeletedAt)
}

func TestUserTagHandler_DeleteTag_NotFound(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	w := ctx.makeRequest("DELETE", "/api/v1/admin/user-tags/999999", nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response["success"].(bool))
}

func TestUserTagHandler_DeleteTag_InvalidID(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	w := ctx.makeRequest("DELETE", "/api/v1/admin/user-tags/invalid", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["message"].(string), "无效的标签ID")
}

// ============================================================================
// AddUserTagHandler Tests
// ============================================================================

func TestUserTagHandler_AddUserTag_Success(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	testUser := integration.CreateTestUser(t, ctx.DB, "testuser")
	tag := ctx.createTestTag(t, "VIP", "#FF5733", "VIP Customer")

	payload := map[string]interface{}{
		"tagId": tag.ID,
	}

	w := ctx.makeRequest("POST", fmt.Sprintf("/api/v1/admin/users/%d/tags", testUser.ID), payload)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	// Verify assignment in database
	var count int64
	ctx.DB.Table("user_tag_relations").Where("user_id = ? AND tag_id = ?", testUser.ID, tag.ID).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestUserTagHandler_AddUserTag_UserNotFound(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	tag := ctx.createTestTag(t, "VIP", "#FF5733", "VIP Customer")

	payload := map[string]interface{}{
		"tagId": tag.ID,
	}

	w := ctx.makeRequest("POST", "/api/v1/admin/users/999999/tags", payload)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response["success"].(bool))
}

func TestUserTagHandler_AddUserTag_TagNotFound(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	testUser := integration.CreateTestUser(t, ctx.DB, "testuser")

	payload := map[string]interface{}{
		"tagId": uint64(999999),
	}

	w := ctx.makeRequest("POST", fmt.Sprintf("/api/v1/admin/users/%d/tags", testUser.ID), payload)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response["success"].(bool))
}

func TestUserTagHandler_AddUserTag_Duplicate(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	testUser := integration.CreateTestUser(t, ctx.DB, "testuser")
	tag := ctx.createTestTag(t, "VIP", "#FF5733", "VIP Customer")

	// Add tag first time
	payload := map[string]interface{}{
		"tagId": tag.ID,
	}
	ctx.makeRequest("POST", fmt.Sprintf("/api/v1/admin/users/%d/tags", testUser.ID), payload)

	// Try to add again
	w := ctx.makeRequest("POST", fmt.Sprintf("/api/v1/admin/users/%d/tags", testUser.ID), payload)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response["success"].(bool))
	assert.Contains(t, response["message"].(string), "用户已拥有此标签")
}

// ============================================================================
// BatchSetUserTagsHandler Tests
// ============================================================================

func TestUserTagHandler_BatchSetUserTags_Success(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	testUser := integration.CreateTestUser(t, ctx.DB, "testuser")
	tag1 := ctx.createTestTag(t, "VIP", "#FF5733", "VIP")
	tag2 := ctx.createTestTag(t, "Premium", "#FFD700", "Premium")

	payload := map[string]interface{}{
		"tagIds": []uint64{tag1.ID, tag2.ID},
	}

	w := ctx.makeRequest("PUT", fmt.Sprintf("/api/v1/admin/users/%d/tags", testUser.ID), payload)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	// Verify assignments
	var count int64
	ctx.DB.Table("user_tag_relations").Where("user_id = ?", testUser.ID).Count(&count)
	assert.Equal(t, int64(2), count)
}

func TestUserTagHandler_BatchSetUserTags_ReplaceExisting(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	testUser := integration.CreateTestUser(t, ctx.DB, "testuser")
	tag1 := ctx.createTestTag(t, "VIP", "#FF5733", "VIP")
	tag2 := ctx.createTestTag(t, "Premium", "#FFD700", "Premium")
	tag3 := ctx.createTestTag(t, "Gold", "#FFFF00", "Gold")

	// Set initial tags
	ctx.DB.Exec("INSERT INTO user_tag_relations (user_id, tag_id) VALUES (?, ?)", testUser.ID, tag1.ID)

	// Batch set should replace
	payload := map[string]interface{}{
		"tagIds": []uint64{tag2.ID, tag3.ID},
	}

	w := ctx.makeRequest("PUT", fmt.Sprintf("/api/v1/admin/users/%d/tags", testUser.ID), payload)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify only new tags exist
	var count int64
	ctx.DB.Table("user_tag_relations").Where("user_id = ? AND tag_id = ?", testUser.ID, tag1.ID).Count(&count)
	assert.Equal(t, int64(0), count)

	ctx.DB.Table("user_tag_relations").Where("user_id = ? AND tag_id IN (?, ?)", testUser.ID, tag2.ID, tag3.ID).Count(&count)
	assert.Equal(t, int64(2), count)
}

// ============================================================================
// GetUserTagsHandler Tests
// ============================================================================

func TestUserTagHandler_GetUserTags_Success(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	testUser := integration.CreateTestUser(t, ctx.DB, "testuser")
	tag1 := ctx.createTestTag(t, "VIP", "#FF5733", "VIP")
	tag2 := ctx.createTestTag(t, "Premium", "#FFD700", "Premium")

	// Assign tags
	ctx.DB.Exec("INSERT INTO user_tag_relations (user_id, tag_id) VALUES (?, ?)", testUser.ID, tag1.ID)
	ctx.DB.Exec("INSERT INTO user_tag_relations (user_id, tag_id) VALUES (?, ?)", testUser.ID, tag2.ID)

	w := ctx.makeRequest("GET", fmt.Sprintf("/api/v1/admin/users/%d/tags", testUser.ID), nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].([]interface{})
	assert.Len(t, data, 2)
}

func TestUserTagHandler_GetUserTags_Empty(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	testUser := integration.CreateTestUser(t, ctx.DB, "testuser")

	w := ctx.makeRequest("GET", fmt.Sprintf("/api/v1/admin/users/%d/tags", testUser.ID), nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].([]interface{})
	assert.Len(t, data, 0)
}

func TestUserTagHandler_GetUserTags_UserNotFound(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	w := ctx.makeRequest("GET", "/api/v1/admin/users/999999/tags", nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response["success"].(bool))
}

// ============================================================================
// RemoveUserTagHandler Tests
// ============================================================================

func TestUserTagHandler_RemoveUserTag_Success(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	testUser := integration.CreateTestUser(t, ctx.DB, "testuser")
	tag := ctx.createTestTag(t, "VIP", "#FF5733", "VIP")

	// Assign tag
	ctx.DB.Exec("INSERT INTO user_tag_relations (user_id, tag_id) VALUES (?, ?)", testUser.ID, tag.ID)

	w := ctx.makeRequest("DELETE", fmt.Sprintf("/api/v1/admin/users/%d/tags/%d", testUser.ID, tag.ID), nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	// Verify removal
	var count int64
	ctx.DB.Table("user_tag_relations").Where("user_id = ? AND tag_id = ?", testUser.ID, tag.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestUserTagHandler_RemoveUserTag_InvalidUserID(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	w := ctx.makeRequest("DELETE", "/api/v1/admin/users/invalid/tags/1", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["message"].(string), "无效的用户ID")
}

func TestUserTagHandler_RemoveUserTag_InvalidTagID(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	testUser := integration.CreateTestUser(t, ctx.DB, "testuser")

	w := ctx.makeRequest("DELETE", fmt.Sprintf("/api/v1/admin/users/%d/tags/invalid", testUser.ID), nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["message"].(string), "无效的标签ID")
}

// ============================================================================
// GetUsersByTagHandler Tests
// ============================================================================

func TestUserTagHandler_GetUsersByTag_Success(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	tag := ctx.createTestTag(t, "VIP", "#FF5733", "VIP")
	user1 := integration.CreateTestUser(t, ctx.DB, "user1")
	user2 := integration.CreateTestUser(t, ctx.DB, "user2")

	// Assign tags
	ctx.DB.Exec("INSERT INTO user_tag_relations (user_id, tag_id) VALUES (?, ?)", user1.ID, tag.ID)
	ctx.DB.Exec("INSERT INTO user_tag_relations (user_id, tag_id) VALUES (?, ?)", user2.ID, tag.ID)

	w := ctx.makeRequest("GET", fmt.Sprintf("/api/v1/admin/user-tags/%d/users", tag.ID), nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].([]interface{})
	assert.Len(t, data, 2)

	// Check pagination
	pagination := response["pagination"].(map[string]interface{})
	assert.Equal(t, float64(1), pagination["page"])
	assert.Equal(t, float64(10), pagination["page_size"])
	assert.Equal(t, float64(2), pagination["total"])
}

func TestUserTagHandler_GetUsersByTag_Pagination(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	tag := ctx.createTestTag(t, "VIP", "#FF5733", "VIP")

	// Create 15 users
	for i := 0; i < 15; i++ {
		user := integration.CreateTestUser(t, ctx.DB, fmt.Sprintf("user%d", i))
		ctx.DB.Exec("INSERT INTO user_tag_relations (user_id, tag_id) VALUES (?, ?)", user.ID, tag.ID)
	}

	w := ctx.makeRequest("GET", fmt.Sprintf("/api/v1/admin/user-tags/%d/users?page=2&page_size=5", tag.ID), nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].([]interface{})
	assert.Len(t, data, 5)

	pagination := response["pagination"].(map[string]interface{})
	assert.Equal(t, float64(2), pagination["page"])
	assert.Equal(t, float64(5), pagination["page_size"])
	assert.Equal(t, float64(15), pagination["total"])
}

func TestUserTagHandler_GetUsersByTag_TagNotFound(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	w := ctx.makeRequest("GET", "/api/v1/admin/user-tags/999999/users", nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response["success"].(bool))
}

// ============================================================================
// BatchDeleteTagsHandler Tests
// ============================================================================

func TestUserTagHandler_BatchDeleteTags_Success(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	tag1 := ctx.createTestTag(t, "VIP", "#FF5733", "VIP")
	tag2 := ctx.createTestTag(t, "Premium", "#FFD700", "Premium")
	tag3 := ctx.createTestTag(t, "Gold", "#FFFF00", "Gold")

	payload := map[string]interface{}{
		"ids": []uint64{tag1.ID, tag2.ID, tag3.ID},
	}

	w := ctx.makeRequest("DELETE", "/api/v1/admin/user-tags/batch", payload)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(3), data["successCount"])
	assert.Equal(t, float64(0), data["failedCount"])
}

func TestUserTagHandler_BatchDeleteTags_PartialFailure(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	tag1 := ctx.createTestTag(t, "VIP", "#FF5733", "VIP")

	payload := map[string]interface{}{
		"ids": []uint64{tag1.ID, 999999},
	}

	w := ctx.makeRequest("DELETE", "/api/v1/admin/user-tags/batch", payload)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["successCount"])
	assert.Equal(t, float64(1), data["failedCount"])
}

func TestUserTagHandler_BatchDeleteTags_EmptyList(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	payload := map[string]interface{}{
		"ids": []uint64{},
	}

	w := ctx.makeRequest("DELETE", "/api/v1/admin/user-tags/batch", payload)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response["success"].(bool))
}

// ============================================================================
// BatchDeleteUserTags (POST) Tests
// ============================================================================

func TestUserTagHandler_BatchDeleteUserTags_Success(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	tag1 := ctx.createTestTag(t, "VIP", "#FF5733", "VIP")
	tag2 := ctx.createTestTag(t, "Premium", "#FFD700", "Premium")

	payload := map[string]interface{}{
		"tag_ids": []uint64{tag1.ID, tag2.ID},
	}

	w := ctx.makeRequest("POST", "/api/v1/admin/user-tags/batch/delete", payload)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["successCount"])
	assert.Equal(t, float64(0), data["failedCount"])
}

func TestUserTagHandler_BatchDeleteUserTags_TooManyTags(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	// Create 101 tag IDs
	tagIDs := make([]uint64, 101)
	for i := 0; i < 101; i++ {
		tagIDs[i] = uint64(i + 1)
	}

	payload := map[string]interface{}{
		"tag_ids": tagIDs,
	}

	w := ctx.makeRequest("POST", "/api/v1/admin/user-tags/batch/delete", payload)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response["success"].(bool))
	assert.Contains(t, response["message"].(string), "maximum 100 tags")
}

// ============================================================================
// BatchAssignTagsToUsers Tests
// ============================================================================

func TestUserTagHandler_BatchAssignTagsToUsers_Success(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	user1 := integration.CreateTestUser(t, ctx.DB, "user1")
	user2 := integration.CreateTestUser(t, ctx.DB, "user2")
	tag1 := ctx.createTestTag(t, "VIP", "#FF5733", "VIP")
	tag2 := ctx.createTestTag(t, "Premium", "#FFD700", "Premium")

	payload := map[string]interface{}{
		"user_ids": []uint64{user1.ID, user2.ID},
		"tag_ids":  []uint64{tag1.ID, tag2.ID},
	}

	w := ctx.makeRequest("POST", "/api/v1/admin/user-tags/batch/assign", payload)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	// Total count should be 4 (2 users * 2 tags)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(4), data["totalCount"])
	assert.Equal(t, float64(4), data["successCount"])
	assert.Equal(t, float64(0), data["failedCount"])

	// Verify in database
	var count int64
	ctx.DB.Table("user_tag_relations").Count(&count)
	assert.Equal(t, int64(4), count)
}

func TestUserTagHandler_BatchAssignTagsToUsers_TooManyUsers(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	// Create 101 user IDs
	userIDs := make([]uint64, 101)
	for i := 0; i < 101; i++ {
		userIDs[i] = uint64(i + 1)
	}

	payload := map[string]interface{}{
		"user_ids": userIDs,
		"tag_ids":  []uint64{1, 2},
	}

	w := ctx.makeRequest("POST", "/api/v1/admin/user-tags/batch/assign", payload)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response["success"].(bool))
	assert.Contains(t, response["message"].(string), "maximum 100 users")
}

func TestUserTagHandler_BatchAssignTagsToUsers_EmptyUserIDs(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	payload := map[string]interface{}{
		"user_ids": []uint64{},
		"tag_ids":  []uint64{1, 2},
	}

	w := ctx.makeRequest("POST", "/api/v1/admin/user-tags/batch/assign", payload)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response["success"].(bool))
	assert.Contains(t, response["message"].(string), "user_ids is required")
}

// ============================================================================
// BatchRemoveTagsFromUsers Tests
// ============================================================================

func TestUserTagHandler_BatchRemoveTagsFromUsers_Success(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	user1 := integration.CreateTestUser(t, ctx.DB, "user1")
	user2 := integration.CreateTestUser(t, ctx.DB, "user2")
	tag1 := ctx.createTestTag(t, "VIP", "#FF5733", "VIP")
	tag2 := ctx.createTestTag(t, "Premium", "#FFD700", "Premium")

	// Assign tags first
	ctx.DB.Exec("INSERT INTO user_tag_relations (user_id, tag_id) VALUES (?, ?)", user1.ID, tag1.ID)
	ctx.DB.Exec("INSERT INTO user_tag_relations (user_id, tag_id) VALUES (?, ?)", user1.ID, tag2.ID)
	ctx.DB.Exec("INSERT INTO user_tag_relations (user_id, tag_id) VALUES (?, ?)", user2.ID, tag1.ID)
	ctx.DB.Exec("INSERT INTO user_tag_relations (user_id, tag_id) VALUES (?, ?)", user2.ID, tag2.ID)

	payload := map[string]interface{}{
		"user_ids": []uint64{user1.ID, user2.ID},
		"tag_ids":  []uint64{tag1.ID, tag2.ID},
	}

	w := ctx.makeRequest("POST", "/api/v1/admin/user-tags/batch/remove", payload)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(4), data["successCount"])
	assert.Equal(t, float64(0), data["failedCount"])

	// Verify removal
	var count int64
	ctx.DB.Table("user_tag_relations").Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestUserTagHandler_BatchRemoveTagsFromUsers_PartialSuccess(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	user1 := integration.CreateTestUser(t, ctx.DB, "user1")
	tag1 := ctx.createTestTag(t, "VIP", "#FF5733", "VIP")

	// Assign one tag
	ctx.DB.Exec("INSERT INTO user_tag_relations (user_id, tag_id) VALUES (?, ?)", user1.ID, tag1.ID)

	payload := map[string]interface{}{
		"user_ids": []uint64{user1.ID},
		"tag_ids":  []uint64{tag1.ID, 999999},
	}

	w := ctx.makeRequest("POST", "/api/v1/admin/user-tags/batch/remove", payload)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	// One removal should succeed, one should fail or be handled gracefully
	assert.True(t, data["successCount"].(float64) >= 1)
}

// ============================================================================
// BatchAssignTagsHandler (Legacy) Tests
// ============================================================================

func TestUserTagHandler_BatchAssignTags_Success(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	user1 := integration.CreateTestUser(t, ctx.DB, "user1")
	user2 := integration.CreateTestUser(t, ctx.DB, "user2")
	tag1 := ctx.createTestTag(t, "VIP", "#FF5733", "VIP")
	tag2 := ctx.createTestTag(t, "Premium", "#FFD700", "Premium")

	payload := map[string]interface{}{
		"userTagPairs": []map[string]uint64{
			{"userId": user1.ID, "tagId": tag1.ID},
			{"userId": user1.ID, "tagId": tag2.ID},
			{"userId": user2.ID, "tagId": tag1.ID},
		},
	}

	w := ctx.makeRequest("POST", "/api/v1/admin/users/user-tags/batch/assign", payload)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(3), data["successCount"])
	assert.Equal(t, float64(0), data["failedCount"])
}

func TestUserTagHandler_BatchAssignTags_EmptyPairs(t *testing.T) {
	ctx := SetupUserTagTest(t)
	ctx.RegisterTagRoutes()

	payload := map[string]interface{}{
		"userTagPairs": []map[string]uint64{},
	}

	w := ctx.makeRequest("POST", "/api/v1/admin/users/user-tags/batch/assign", payload)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response["success"].(bool))
}

// ============================================================================
// Pagination Helper Tests
// ============================================================================

func TestGetPaginationParams_DefaultValues(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/api/v1/admin/user-tags", nil)

	page, pageSize := getPaginationParams(c)

	assert.Equal(t, 1, page)
	assert.Equal(t, 10, pageSize)
}

func TestGetPaginationParams_CustomValues(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/api/v1/admin/user-tags?page=3&page_size=25", nil)

	page, pageSize := getPaginationParams(c)

	assert.Equal(t, 3, page)
	assert.Equal(t, 25, pageSize)
}

func TestGetPaginationParams_InvalidValues(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/api/v1/admin/user-tags?page=invalid&page_size=invalid", nil)

	page, pageSize := getPaginationParams(c)

	// Should fall back to defaults
	assert.Equal(t, 1, page)
	assert.Equal(t, 10, pageSize)
}

// ============================================================================
// Response Helper Tests
// ============================================================================

func TestToTagResponse(t *testing.T) {
	tag := &model.UserTag{
		Base: model.Base{
			ID: 123,
		},
		Name:        "VIP",
		Color:       "#FF5733",
		Description: "VIP Customer",
	}

	response := toTagResponse(tag)

	assert.Equal(t, uint64(123), response.ID)
	assert.Equal(t, "VIP", response.Name)
	assert.Equal(t, "#FF5733", response.Color)
	assert.Equal(t, "VIP Customer", response.Description)
	assert.NotEmpty(t, response.CreatedAt)
	assert.NotEmpty(t, response.UpdatedAt)
}
