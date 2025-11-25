package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	adminhandler "gamelink/internal/handler/admin"
	"gamelink/internal/model"
	"gamelink/internal/repository/game"
	"gamelink/internal/repository/player"
	"gamelink/internal/repository/role"
	"gamelink/internal/repository/user"
	adminservice "gamelink/internal/service/admin"
	"gamelink/internal/testutil"
	"gamelink/internal/cache"
)

// 管理端用户管理：创建用户 -> 查询 -> 更新状态/角色 -> 删除
func TestAdminUserManagement(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateUserModels(t, db)

	// repos & service
	userRepo := user.NewUserRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	gameRepo := game.NewGameRepository(db)
	roleRepo := role.NewRoleRepository(db)
	memCache := cache.NewMemory()
	adminSvc := adminservice.NewAdminService(gameRepo, userRepo, playerRepo, nil, nil, roleRepo, memCache)
	userHandler := adminhandler.NewUserHandler(adminSvc)

	router := gin.New()
	adminGroup := router.Group("/api/v1/admin")
	adminGroup.POST("/users", userHandler.CreateUser)
	adminGroup.GET("/users/:id", userHandler.GetUser)
	adminGroup.PUT("/users/:id/status", userHandler.UpdateUserStatus)
	adminGroup.PUT("/users/:id/role", userHandler.UpdateUserRole)
	adminGroup.DELETE("/users/:id", userHandler.DeleteUser)

	// 创建用户
	createPayload := map[string]interface{}{
		"name":     "Alice",
		"email":    "alice@example.com",
		"password": "Passw0rd",
		"role":     "user",
		"status":   "active",
	}
	createResp := doJSON(router, http.MethodPost, "/api/v1/admin/users", createPayload, "")
	if createResp.Code != http.StatusCreated {
		t.Fatalf("create user status=%d body=%s", createResp.Code, createResp.Body.String())
	}
	var createParsed apiResp[model.User]
	if err := json.Unmarshal(createResp.Body.Bytes(), &createParsed); err != nil {
		t.Fatalf("parse create resp: %v", err)
	}
	userID := createParsed.Data.ID
	if userID == 0 {
		t.Fatalf("expected user id > 0, got %d", userID)
	}

	// 查询用户
	getResp := doJSON(router, http.MethodGet, "/api/v1/admin/users/"+uintToStr(userID), nil, "")
	if getResp.Code != http.StatusOK {
		t.Fatalf("get user status=%d body=%s", getResp.Code, getResp.Body.String())
	}

	// 更新状态
	statusResp := doJSON(router, http.MethodPut, "/api/v1/admin/users/"+uintToStr(userID)+"/status", map[string]string{
		"status": "banned",
	}, "")
	if statusResp.Code != http.StatusOK {
		t.Fatalf("update status code=%d body=%s", statusResp.Code, statusResp.Body.String())
	}

	// 更新角色
	roleResp := doJSON(router, http.MethodPut, "/api/v1/admin/users/"+uintToStr(userID)+"/role", map[string]string{
		"role": "admin",
	}, "")
	if roleResp.Code != http.StatusOK {
		t.Fatalf("update role code=%d body=%s", roleResp.Code, roleResp.Body.String())
	}

	// 删除用户
	delResp := doJSON(router, http.MethodDelete, "/api/v1/admin/users/"+uintToStr(userID), nil, "")
	if delResp.Code != http.StatusOK {
		t.Fatalf("delete user code=%d body=%s", delResp.Code, delResp.Body.String())
	}
}

func migrateUserModels(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(
		&model.User{},
		&model.RoleModel{},
		&model.Permission{},
	); err != nil {
		t.Fatalf("migrate user models: %v", err)
	}
}
