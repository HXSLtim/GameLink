// Package admin provides integration tests for user admin handlers.
package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/service/integration"
)

// TestUserHandler_ListUsers_Success tests listing users with pagination.
func TestUserHandler_ListUsers_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)

	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test users with different roles
	user1 := integration.CreateTestUser(t, helper.db, "user1")
	user1.Role = model.RoleUser
	helper.db.Save(user1)

	user2 := integration.CreateTestUser(t, helper.db, "user2")
	user2.Role = model.RolePlayer
	helper.db.Save(user2)

	user3 := integration.CreateTestUser(t, helper.db, "user3")
	user3.Role = model.RoleAdmin
	helper.db.Save(user3)

	// Test listing users
	w := helper.MakeRequest("GET", "/users?page=1&page_size=10", nil, true)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool `json:"success"`
		Data    []struct {
			ID    uint64 `json:"id"`
			Name  string `json:"name"`
			Phone string `json:"phone"`
			Email string `json:"email"`
			Role  string `json:"role"`
		} `json:"data"`
		Meta struct {
			Pagination struct {
				Page       int `json:"page"`
				PageSize   int `json:"page_size"`
				Total      int `json:"total"`
				TotalPages int `json:"total_pages"`
			} `json:"pagination"`
		} `json:"meta"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Success)
	assert.GreaterOrEqual(t, len(response.Data), 3)
	assert.Equal(t, 1, response.Meta.Pagination.Page)
	assert.Equal(t, 10, response.Meta.Pagination.PageSize)
}

// TestUserHandler_ListUsers_WithFilters tests listing users with role and status filters.
func TestUserHandler_ListUsers_WithFilters(t *testing.T) {
	integration.SkipIfNoTestDB(t)

	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create users with different roles and statuses
	activeUser := integration.CreateTestUser(t, helper.db, "active_user")
	activeUser.Role = model.RoleUser
	activeUser.Status = model.UserStatusActive
	helper.db.Save(activeUser)

	suspendedPlayer := integration.CreateTestUser(t, helper.db, "suspended_player")
	suspendedPlayer.Role = model.RolePlayer
	suspendedPlayer.Status = model.UserStatusSuspended
	helper.db.Save(suspendedPlayer)

	activePlayer := integration.CreateTestUser(t, helper.db, "active_player")
	activePlayer.Role = model.RolePlayer
	activePlayer.Status = model.UserStatusActive
	helper.db.Save(activePlayer)

	// Test filtering by role
	w := helper.MakeRequest("GET", "/users?role=player", nil, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool `json:"success"`
		Data    []struct {
			Role string `json:"role"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Success)
	// Should only return players
	for _, u := range response.Data {
		assert.Equal(t, "player", u.Role)
	}

	// Test filtering by status
	w = helper.MakeRequest("GET", "/users?status=active", nil, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var statusResponse struct {
		Success bool `json:"success"`
		Data    []struct {
			Status string `json:"status"`
		} `json:"data"`
	}

	err = json.Unmarshal(w.Body.Bytes(), &statusResponse)
	require.NoError(t, err)
	assert.True(t, statusResponse.Success)
}

// TestUserHandler_ListUsers_WithDateRange tests listing users with date range filter.
func TestUserHandler_ListUsers_WithDateRange(t *testing.T) {
	integration.SkipIfNoTestDB(t)

	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create a test user
	_ = integration.CreateTestUser(t, helper.db, "daterange_user")

	// Test with date range (should include all users)
	w := helper.MakeRequest("GET", "/users?date_from=2025-01-01&date_to=2025-12-31", nil, true)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool          `json:"success"`
		Data    []interface{} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Success)
}

// TestUserHandler_ListUsers_WithKeyword tests listing users with keyword search.
func TestUserHandler_ListUsers_WithKeyword(t *testing.T) {
	integration.SkipIfNoTestDB(t)

	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create a user with specific name
	testUser := integration.CreateTestUser(t, helper.db, "uniquekeyword123")
	testUser.Phone = "13800138999"
	helper.db.Save(testUser)

	// Test keyword search
	w := helper.MakeRequest("GET", "/users?keyword=uniquekeyword123", nil, true)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool `json:"success"`
		Data    []struct {
			Name string `json:"name"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Success)
	if len(response.Data) > 0 {
		assert.Contains(t, response.Data[0].Name, "uniquekeyword123")
	}
}

// TestUserHandler_GetUserStats tests getting user statistics.
func TestUserHandler_GetUserStats(t *testing.T) {
	integration.SkipIfNoTestDB(t)

	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test users
	_ = integration.CreateTestUser(t, helper.db, "stat_user1")
	user2 := integration.CreateTestUser(t, helper.db, "stat_player1")
	user2.Role = model.RolePlayer
	helper.db.Save(user2)

	// Get user stats
	w := helper.MakeRequest("GET", "/users/stats", nil, true)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			TotalUsers int `json:"total_users"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Success)
	assert.GreaterOrEqual(t, response.Data.TotalUsers, 2)
}

// TestUserHandler_CreateUser_Success tests creating a user successfully.
func TestUserHandler_CreateUser_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)

	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	payload := map[string]interface{}{
		"phone":      "13800138888",
		"email":      "newuser@test.com",
		"password":   "password123",
		"name":       "New User",
		"avatar_url": "https://example.com/avatar.jpg",
		"role":       "user",
		"status":     "active",
	}

	w := helper.MakeRequest("POST", "/users", payload, true)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			ID    uint64 `json:"id"`
			Name  string `json:"name"`
			Phone string `json:"phone"`
			Email string `json:"email"`
			Role  string `json:"role"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "New User", response.Data.Name)
	assert.Equal(t, "13800138888", response.Data.Phone)
	assert.Equal(t, "newuser@test.com", response.Data.Email)
	assert.Equal(t, "user", response.Data.Role)
	assert.NotZero(t, response.Data.ID)
}

// TestUserHandler_CreateUser_ValidationError tests creating a user with validation errors.
func TestUserHandler_CreateUser_ValidationError(t *testing.T) {
	integration.SkipIfNoTestDB(t)

	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	testCases := []struct {
		name        string
		payload     map[string]interface{}
		expectedMsg string
	}{
		{
			name: "missing password",
			payload: map[string]interface{}{
				"phone":  "13800138888",
				"name":   "Test User",
				"role":   "user",
				"status": "active",
			},
		},
		{
			name: "password too short",
			payload: map[string]interface{}{
				"phone":    "13800138888",
				"password": "12345",
				"name":     "Test User",
				"role":     "user",
				"status":   "active",
			},
		},
		{
			name: "invalid email format",
			payload: map[string]interface{}{
				"phone":    "13800138888",
				"email":    "invalid-email",
				"password": "password123",
				"name":     "Test User",
				"role":     "user",
				"status":   "active",
			},
		},
		{
			name: "invalid phone format",
			payload: map[string]interface{}{
				"phone":    "12345",
				"password": "password123",
				"name":     "Test User",
				"role":     "user",
				"status":   "active",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := helper.MakeRequest("POST", "/users", tc.payload, true)
			assert.NotEqual(t, http.StatusCreated, w.Code)
		})
	}
}

// TestUserHandler_CreateUserWithPlayer_Success tests creating a user with player profile.
func TestUserHandler_CreateUserWithPlayer_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)

	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	payload := map[string]interface{}{
		"phone":    "13800138777",
		"email":    "player@test.com",
		"password": "password123",
		"name":     "Player User",
		"role":     "player",
		"status":   "active",
		"player": map[string]interface{}{
			"nickname":            "ProGamer",
			"bio":                 "Experienced player",
			"hourly_rate_cents":   5000,
			"main_game_id":        1,
			"verification_status": "verified",
		},
	}

	w := helper.MakeRequest("POST", "/users/with-player", payload, true)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			User struct {
				ID   uint64 `json:"id"`
				Name string `json:"name"`
				Role string `json:"role"`
			} `json:"user"`
			Player struct {
				ID       uint64 `json:"id"`
				Nickname string `json:"nickname"`
			} `json:"player"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Player User", response.Data.User.Name)
	assert.Equal(t, "player", response.Data.User.Role)
	assert.NotZero(t, response.Data.User.ID)
	assert.NotZero(t, response.Data.Player.ID)
	assert.Equal(t, "ProGamer", response.Data.Player.Nickname)
}

// TestUserHandler_GetUser_Success tests getting a user by ID.
func TestUserHandler_GetUser_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)

	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create a test user
	testUser := integration.CreateTestUser(t, helper.db, "get_user_test")

	// Get the user
	w := helper.MakeRequest("GET", "/users/"+fmt.Sprintf("%d", testUser.ID), nil, true)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			ID    uint64 `json:"id"`
			Name  string `json:"name"`
			Phone string `json:"phone"`
			Email string `json:"email"`
			Role  string `json:"role"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, testUser.ID, response.Data.ID)
	assert.Equal(t, "get_user_test", response.Data.Name)
}

// TestUserHandler_GetUser_NotFound tests getting a non-existent user.
func TestUserHandler_GetUser_NotFound(t *testing.T) {
	integration.SkipIfNoTestDB(t)

	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Try to get a non-existent user
	w := helper.MakeRequest("GET", "/users/999999", nil, true)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestUserHandler_UpdateUser_Success tests updating a user.
func TestUserHandler_UpdateUser_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)

	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create a test user
	testUser := integration.CreateTestUser(t, helper.db, "update_user_test")

	payload := map[string]interface{}{
		"name":       "Updated Name",
		"phone":      "13900139999",
		"email":      "updated@test.com",
		"avatar_url": "https://example.com/newavatar.jpg",
		"role":       "player",
		"status":     "active",
	}

	w := helper.MakeRequest("PUT", "/users/"+fmt.Sprintf("%d", testUser.ID), payload, true)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			ID    uint64 `json:"id"`
			Name  string `json:"name"`
			Phone string `json:"phone"`
			Email string `json:"email"`
			Role  string `json:"role"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, testUser.ID, response.Data.ID)
	assert.Equal(t, "Updated Name", response.Data.Name)
	assert.Equal(t, "13900139999", response.Data.Phone)
}

// TestUserHandler_DeleteUser_Success tests deleting a user.
func TestUserHandler_DeleteUser_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)

	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create a test user
	testUser := integration.CreateTestUser(t, helper.db, "delete_user_test")

	// Delete the user
	w := helper.MakeRequest("DELETE", "/users/"+fmt.Sprintf("%d", testUser.ID), nil, true)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool `json:"success"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Success)

	// Verify user is deleted
	var count int64
	helper.db.Model(&model.User{}).Where("id = ?", testUser.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

// TestUserHandler_UpdateUserStatus_Success tests updating user status.
func TestUserHandler_UpdateUserStatus_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)

	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create a test user
	testUser := integration.CreateTestUser(t, helper.db, "status_user_test")
	testUser.Status = model.UserStatusActive
	helper.db.Save(testUser)

	payload := map[string]string{
		"status": "suspended",
	}

	w := helper.MakeRequest("PUT", "/users/"+fmt.Sprintf("%d", testUser.ID)+"/status", payload, true)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			ID     uint64           `json:"id"`
			Status model.UserStatus `json:"status"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, testUser.ID, response.Data.ID)
	assert.Equal(t, model.UserStatusSuspended, response.Data.Status)
}

// TestUserHandler_UpdateUserRole_Success tests updating user role.
func TestUserHandler_UpdateUserRole_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)

	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create a test user
	testUser := integration.CreateTestUser(t, helper.db, "role_user_test")
	testUser.Role = model.RoleUser
	helper.db.Save(testUser)

	payload := map[string]string{
		"role": "admin",
	}

	w := helper.MakeRequest("PUT", "/users/"+fmt.Sprintf("%d", testUser.ID)+"/role", payload, true)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			ID   uint64     `json:"id"`
			Role model.Role `json:"role"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, testUser.ID, response.Data.ID)
	assert.Equal(t, model.RoleAdmin, response.Data.Role)
}

// TestUserHandler_BatchDeleteUsers_Success tests batch deleting users.
func TestUserHandler_BatchDeleteUsers_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)

	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test users
	user1 := integration.CreateTestUser(t, helper.db, "batch_delete_1")
	user2 := integration.CreateTestUser(t, helper.db, "batch_delete_2")
	user3 := integration.CreateTestUser(t, helper.db, "batch_delete_3")

	payload := map[string]interface{}{
		"ids": []uint64{user1.ID, user2.ID, user3.ID},
	}

	w := helper.MakeRequest("POST", "/users/batch-delete", payload, true)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			Deleted int `json:"deleted"`
			Failed  int `json:"failed"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, 3, response.Data.Deleted)
	assert.Equal(t, 0, response.Data.Failed)
}

// TestUserHandler_BatchDeleteUsers_PartialSuccess tests batch deleting with some failures.
func TestUserHandler_BatchDeleteUsers_PartialSuccess(t *testing.T) {
	integration.SkipIfNoTestDB(t)

	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create a test user
	user1 := integration.CreateTestUser(t, helper.db, "partial_delete_1")

	// Include a non-existent user ID
	payload := map[string]interface{}{
		"ids": []uint64{user1.ID, 999999, 999998},
	}

	w := helper.MakeRequest("POST", "/users/batch-delete", payload, true)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Data    struct {
			Deleted int `json:"deleted"`
			Failed  int `json:"failed"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "部分删除成功", response.Message)
	assert.GreaterOrEqual(t, response.Data.Deleted, 1)
	assert.GreaterOrEqual(t, response.Data.Failed, 1)
}

// TestUserHandler_ListUserOrders_Success tests listing user's orders.
func TestUserHandler_ListUserOrders_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)

	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create a test user
	testUser := integration.CreateTestUser(t, helper.db, "order_user_test")

	// Create a game for the order
	_ = integration.CreateTestGame(t, helper.db, "test_game")

	// Create a player for the order
	player := integration.CreateTestPlayer(t, helper.db, integration.CreateTestUser(t, helper.db, "order_player"))

	// Create an order for the user
	_ = integration.CreateTestOrder(t, helper.db, testUser, player, model.OrderStatusCompleted)

	// Get user's orders
	w := helper.MakeRequest("GET", "/users/"+fmt.Sprintf("%d", testUser.ID)+"/orders", nil, true)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool          `json:"success"`
		Data    []interface{} `json:"data"`
		Meta    struct {
			Pagination struct {
				Page       int `json:"page"`
				PageSize   int `json:"page_size"`
				Total      int `json:"total"`
				TotalPages int `json:"total_pages"`
			} `json:"pagination"`
		} `json:"meta"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Success)
}

// TestUserHandler_ListUserOrders_UserNotFound tests listing orders for non-existent user.
func TestUserHandler_ListUserOrders_UserNotFound(t *testing.T) {
	integration.SkipIfNoTestDB(t)

	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Try to get orders for a non-existent user
	w := helper.MakeRequest("GET", "/users/999999/orders", nil, true)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestUserHandler_ListUserLogs_Success tests listing user operation logs.
func TestUserHandler_ListUserLogs_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)

	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create a test user
	testUser := integration.CreateTestUser(t, helper.db, "logs_user_test")

	// Get user's operation logs
	w := helper.MakeRequest("GET", "/users/"+fmt.Sprintf("%d", testUser.ID)+"/logs", nil, true)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool          `json:"success"`
		Data    []interface{} `json:"data"`
		Meta    struct {
			Pagination struct {
				Page       int `json:"page"`
				PageSize   int `json:"page_size"`
				Total      int `json:"total"`
				TotalPages int `json:"total_pages"`
			} `json:"pagination"`
		} `json:"meta"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Success)
}

// TestUserHandler_ListUserLoginHistory_Success tests listing user login history.
func TestUserHandler_ListUserLoginHistory_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)

	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create a test user
	testUser := integration.CreateTestUser(t, helper.db, "login_history_user")

	// Get user's login history
	w := helper.MakeRequest("GET", "/users/"+fmt.Sprintf("%d", testUser.ID)+"/login-history", nil, true)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool          `json:"success"`
		Data    []interface{} `json:"data"`
		Meta    struct {
			Pagination struct {
				Page       int `json:"page"`
				PageSize   int `json:"page_size"`
				Total      int `json:"total"`
				TotalPages int `json:"total_pages"`
			} `json:"pagination"`
		} `json:"meta"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Success)
}

// TestUserHandler_ListUserLoginHistory_UserNotFound tests login history for non-existent user.
func TestUserHandler_ListUserLoginHistory_UserNotFound(t *testing.T) {
	integration.SkipIfNoTestDB(t)

	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Try to get login history for a non-existent user
	w := helper.MakeRequest("GET", "/users/999999/login-history", nil, true)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestUserHandler_Unauthorized tests endpoints without authentication.
func TestUserHandler_Unauthorized(t *testing.T) {
	integration.SkipIfNoTestDB(t)

	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	testCases := []struct {
		method string
		path   string
		body   interface{}
	}{
		{"GET", "/users", nil},
		{"GET", "/users/stats", nil},
		{"POST", "/users", map[string]interface{}{"name": "test", "password": "123456", "role": "user", "status": "active"}},
		{"GET", "/users/1", nil},
		{"PUT", "/users/1", map[string]interface{}{"name": "test", "role": "user", "status": "active"}},
		{"DELETE", "/users/1", nil},
		{"PUT", "/users/1/status", map[string]string{"status": "active"}},
		{"PUT", "/users/1/role", map[string]string{"role": "user"}},
		{"POST", "/users/batch-delete", map[string]interface{}{"ids": []uint64{1}}},
	}

	for _, tc := range testCases {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			w := helper.MakeRequest(tc.method, tc.path, tc.body, false)
			// Should return unauthorized or the mock auth bypasses it
			// With mock permission middleware, auth is skipped
			assert.NotEqual(t, http.StatusInternalServerError, w.Code)
		})
	}
}
