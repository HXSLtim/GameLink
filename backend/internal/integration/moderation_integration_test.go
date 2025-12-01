package integration

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"gamelink/internal/auth"
	"gamelink/internal/cache"
	"gamelink/internal/handler"
	adminhandler "gamelink/internal/handler/admin"
	userhandler "gamelink/internal/handler/user"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	chatsrepo "gamelink/internal/repository/chat"
	"gamelink/internal/repository/game"
	orderimpl "gamelink/internal/repository/implementations"
	"gamelink/internal/repository/payment"
	"gamelink/internal/repository/player"
	"gamelink/internal/repository/role"
	"gamelink/internal/repository/serviceitem"
	"gamelink/internal/repository/user"
	adminservice "gamelink/internal/service/admin"
	authservice "gamelink/internal/service/auth"
	chatservice "gamelink/internal/service/chat"
	"gamelink/internal/testutil"
)

// 封禁：登录成功 -> 管理端封禁 -> 登录被拒
func TestUserBanAfterAdminUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateAuthModels(t, db)

	userRepo := user.NewUserRepository(db)
	roleRepo := role.NewRoleRepository(db)
	gameRepo := game.NewGameRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("Passw0rd"), bcrypt.DefaultCost)
	u := &model.User{Name: "BanUser", Email: "ban@example.com", Phone: "19912345678", PasswordHash: string(hashed), Role: model.RoleUser, Status: model.UserStatusActive}
	if err := userRepo.Create(context.Background(), u); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	jwtMgr := auth.NewJWTManager("secret", time.Hour)
	authSvc := authservice.NewAuthService(userRepo, jwtMgr)
	serviceItemRepo := serviceitem.NewServiceItemRepository(db)
	adminSvc := adminservice.NewAdminService(gameRepo, userRepo, playerRepo, orderRepo, paymentRepo, roleRepo, serviceItemRepo, cache.NewMemory())

	authRouter := gin.New()
	api := authRouter.Group("/api/v1")
	handler.RegisterAuthRoutes(api, authSvc)
	adminGroup := api.Group("/admin")
	userHandler := adminhandler.NewUserHandler(adminSvc)
	adminGroup.PUT("/users/:id/status", userHandler.UpdateUserStatus)

	loginPayload := map[string]string{"username": u.Email, "password": "Passw0rd"}
	loginResp := doJSON(authRouter, http.MethodPost, "/api/v1/auth/login", loginPayload, "")
	if loginResp.Code != http.StatusOK {
		t.Fatalf("login status=%d body=%s", loginResp.Code, loginResp.Body.String())
	}

	// ban user
	banResp := doJSON(authRouter, http.MethodPut, "/api/v1/admin/users/"+uintToStr(u.ID)+"/status", map[string]string{"status": "banned"}, "")
	if banResp.Code != http.StatusOK {
		t.Fatalf("ban status=%d body=%s", banResp.Code, banResp.Body.String())
	}

	// login should now be forbidden
	loginResp2 := doJSON(authRouter, http.MethodPost, "/api/v1/auth/login", loginPayload, "")
	if loginResp2.Code != http.StatusForbidden {
		t.Fatalf("expected forbidden after ban, got %d body=%s", loginResp2.Code, loginResp2.Body.String())
	}
}

// 举报聊天消息：创建举报记录
func TestChatMessageReport(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateReportModels(t, db)

	// seed user
	userRepo := user.NewUserRepository(db)
	u := &model.User{Name: "Reporter", Email: "r@example.com", Phone: "18800000001", PasswordHash: "x", Role: model.RoleUser, Status: model.UserStatusActive}
	_ = userRepo.Create(context.Background(), u)

	reportRepo := chatsrepo.NewChatReportRepository(db)
	chatSvc := chatservice.NewChatService(nil, nil, nil, reportRepo, cache.NewMemory())

	router := gin.New()
	api := router.Group("/api/v1")
	auth := fakeAuthMiddleware(u.ID)
	userGroup := api.Group("/user")
	userhandler.RegisterChatRoutes(userGroup, chatSvc, auth)

	payload := map[string]string{"reason": "spam", "evidence": "screenshot"}
	resp := doJSON(router, http.MethodPost, "/api/v1/user/chat/messages/1/report", payload, "")
	if resp.Code != http.StatusOK {
		t.Fatalf("report status=%d body=%s", resp.Code, resp.Body.String())
	}

	// verify record exists
	reports, total, err := reportRepo.List(context.Background(), repository.ChatReportListOptions{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("list reports: %v", err)
	}
	if total == 0 || len(reports) == 0 || reports[0].Reason == "" {
		t.Fatalf("expected report record, got %+v", reports)
	}
}

func migrateAuthModels(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(&model.User{}, &model.RoleModel{}, &model.Permission{}, &model.RolePermission{}, &model.UserRole{}); err != nil {
		t.Fatalf("migrate auth models: %v", err)
	}
}

func migrateReportModels(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(&model.User{}, &model.ChatReport{}); err != nil {
		t.Fatalf("migrate report models: %v", err)
	}
}
