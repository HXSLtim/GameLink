package main

// GameLink API
//
// @title           GameLink API
// @version         0.3.0
// @description     GameLink 平台 API，包含健康检查、认证与管理端能力。
// @BasePath        /api/v1
// @schemes         http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"gamelink/internal/auth"
	"gamelink/internal/cache"
	"gamelink/internal/config"
	"gamelink/internal/db"
	"gamelink/internal/handler"
	adminhandler "gamelink/internal/handler/admin"
	"gamelink/internal/handler/middleware"
	playerhandler "gamelink/internal/handler/player"
	userhandler "gamelink/internal/handler/user"
	"gamelink/internal/logging"
	"gamelink/internal/model"
	commissionrepo "gamelink/internal/repository/commission"
	"gamelink/internal/repository/common"
	gamerepo "gamelink/internal/repository/game"
	chatrepo "gamelink/internal/repository/chat"
	orderrepo "gamelink/internal/repository/order"
	paymentrepo "gamelink/internal/repository/payment"
	permissionrepo "gamelink/internal/repository/permission"
	playerrepo "gamelink/internal/repository/player"
	playertagrepo "gamelink/internal/repository/player_tag"
	rankingrepo "gamelink/internal/repository/ranking"
	reviewrepo "gamelink/internal/repository/review"
	rolerepo "gamelink/internal/repository/role"
	serviceitemrepo "gamelink/internal/repository/serviceitem"
	statsrepo "gamelink/internal/repository/stats"
	userrepo "gamelink/internal/repository/user"
	withdrawrepo "gamelink/internal/repository/withdraw"
	"gamelink/internal/scheduler"
	adminservice "gamelink/internal/service/admin"
	authservice "gamelink/internal/service/auth"
	chatservice "gamelink/internal/service/chat"
	commissionservice "gamelink/internal/service/commission"
	earningsservice "gamelink/internal/service/earnings"
	giftservice "gamelink/internal/service/gift"
	itemservice "gamelink/internal/service/item"
	orderservice "gamelink/internal/service/order"
	paymentservice "gamelink/internal/service/payment"
	permissionservice "gamelink/internal/service/permission"
	playerservice "gamelink/internal/service/player"
	reviewservice "gamelink/internal/service/review"
	roleservice "gamelink/internal/service/role"
	statsservice "gamelink/internal/service/stats"
)

func main() {
	cfg := config.Load()
	if err := config.Validate(os.Getenv("APP_ENV"), cfg); err != nil {
		log.Fatalf("配置校验失败: %v", err)
	}
	if cfg.Crypto.Enabled {
		log.Printf("crypto middleware enabled, methods=%v exclude=%v use_signature=%v", cfg.Crypto.Methods, cfg.Crypto.ExcludePaths, cfg.Crypto.UseSignature)
	} else {
		log.Println("crypto middleware disabled")
	}

	orm, err := db.Open(cfg)
	if err != nil {
		log.Fatalf("打开数据库失败: %v", err)
	}
	sqlDB, err := orm.DB()
	if err != nil {
		log.Fatalf("获取底层连接失败: %v", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("close db error: %v", err)
		}
	}()

	cacheClient, err := cache.New(cfg.Cache)
	if err != nil {
		log.Fatalf("初始化缓存失败: %v", err)
	}
	defer func() {
		if err := cacheClient.Close(context.Background()); err != nil {
			log.Printf("close cache error: %v", err)
		}
	}()

	// RBAC - 初始化 RoleRepository（需要在 AdminService 之前）
	roleRepo := rolerepo.NewRoleRepository(orm)

	adminSvc := adminservice.NewAdminService(
		gamerepo.NewGameRepository(orm),
		userrepo.NewUserRepository(orm),
		playerrepo.NewPlayerRepository(orm),
		orderrepo.NewOrderRepository(orm),
		paymentrepo.NewPaymentRepository(orm),
		roleRepo,
		cacheClient,
	)

	// Inject transaction manager for composite operations
	uow := common.NewUnitOfWork(orm)
	adminSvc.SetTxManager(uow)

	gin.SetMode(resolveGinMode())

	// 初始化结构化日志（slog），从 LOG_LEVEL 读取级别
	_ = logging.Init(os.Getenv("LOG_LEVEL"))

	router := gin.New()

	// 注册全局中间件（按顺序执行）
	router.Use(middleware.RequestID())
	router.Use(middleware.SlogLogger())        // 结构化访问日志
	router.Use(middleware.MetricsMiddleware()) // HTTP 指标
	router.Use(middleware.Crypto(cfg.Crypto))  // 请求解密（与前端 AES 中间件对齐）
	router.Use(middleware.ErrorMap())          // 统一错误映射（ErrValidation/ErrNotFound）
	router.Use(middleware.Recovery())          // 统一JSON恢复中间件
	router.Use(middleware.CORS())              // CORS中间件（跨域处理）

	// Register root and health on both base and versioned API for compatibility
	handler.RegisterRoot(router)
	handler.RegisterHealth(router)

	// Versioned API group
	api := router.Group("/api/v1")
	handler.RegisterRoot(api)
	handler.RegisterHealth(api)

	// Metrics endpoint
	router.GET("/metrics", middleware.MetricsHandler())

	// Auth service and routes
	jwtSecret := strings.TrimSpace(cfg.Auth.JWTSecret)
	if jwtSecret == "" {
		if os.Getenv("APP_ENV") == "production" {
			log.Fatal("JWT secret must be provided via configs.auth.jwt_secret or JWT_SECRET_KEY")
		}
		jwtSecret = config.DefaultDevJWTSecret
	}
	tokenTTL := time.Duration(cfg.Auth.TokenTTLHours) * time.Hour
	if tokenTTL <= 0 {
		tokenTTL = 24 * time.Hour
	}
	jwtMgr := auth.NewJWTManager(jwtSecret, tokenTTL)
	authSvc := authservice.NewAuthService(userrepo.NewUserRepository(orm), jwtMgr)
	handler.RegisterAuthRoutes(api, authSvc)

	// Initialize repositories (reuse where possible)
	userRepo := userrepo.NewUserRepository(orm)
	playerRepo := playerrepo.NewPlayerRepository(orm)
	gameRepo := gamerepo.NewGameRepository(orm)
	orderRepo := orderrepo.NewOrderRepository(orm)
	chatGroupRepo := chatrepo.NewChatGroupRepository(orm)
	chatMemberRepo := chatrepo.NewChatMemberRepository(orm)
	chatMessageRepo := chatrepo.NewChatMessageRepository(orm)
	chatReportRepo := chatrepo.NewChatReportRepository(orm)
	paymentRepo := paymentrepo.NewPaymentRepository(orm)
	reviewRepo := reviewrepo.NewReviewRepository(orm)
	playerTagRepo := playertagrepo.NewPlayerTagRepository(orm)
	withdrawRepo := withdrawrepo.NewWithdrawRepository(orm)
	commissionRepo := commissionrepo.NewCommissionRepository(orm)
	serviceItemRepo := serviceitemrepo.NewServiceItemRepository(orm)
	rankingCommissionRepo := rankingrepo.NewRankingCommissionRepository(orm)

	// Initialize user-side services
	commissionSvc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)
	serviceItemSvc := itemservice.NewServiceItemService(serviceItemRepo, gameRepo, playerRepo)
	giftSvc := giftservice.NewGiftService(serviceItemRepo, orderRepo, playerRepo, commissionRepo)
	orderSvc := orderservice.NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)
	// Inject chat group repo for order chat auto-destroy
	orderSvc.SetChatGroupRepository(chatGroupRepo)
	paymentSvc := paymentservice.NewPaymentService(paymentRepo, orderRepo)
	playerSvc := playerservice.NewPlayerService(playerRepo, userRepo, gameRepo, orderRepo, reviewRepo, playerTagRepo, cacheClient)
	reviewSvc := reviewservice.NewReviewService(reviewRepo, orderRepo, playerRepo, userRepo)
	earningsSvc := earningsservice.NewEarningsService(playerRepo, orderRepo, withdrawRepo)
	chatSvc := chatservice.NewChatService(chatGroupRepo, chatMemberRepo, chatMessageRepo, chatReportRepo, cacheClient)

	// Initialize settlement scheduler
	settlementScheduler := scheduler.NewSettlementScheduler(commissionSvc)
	settlementScheduler.Start()
	defer settlementScheduler.Stop()

	// Initialize chat retention scheduler (30 days retention)
	chatRetention := scheduler.NewChatRetentionScheduler(chatGroupRepo, chatMessageRepo, 30)
	chatRetention.Start()
	defer chatRetention.Stop()

	// Register user-side routes (require authentication)
	authMiddleware := middleware.JWTAuth()
	userGroup := api.Group("/user")
	userGroup.Use(authMiddleware)
	{
		userhandler.RegisterOrderRoutes(userGroup, orderSvc, authMiddleware)
		userhandler.RegisterPaymentRoutes(userGroup, paymentSvc, authMiddleware)
		userhandler.RegisterPlayerRoutes(userGroup, playerSvc, authMiddleware)
		userhandler.RegisterReviewRoutes(userGroup, reviewSvc, authMiddleware)
		userhandler.RegisterGiftRoutes(userGroup, giftSvc, serviceItemSvc, authMiddleware)
		userhandler.RegisterChatRoutes(userGroup, chatSvc, authMiddleware)
	}

	// Register player-side routes (require authentication)
	playerGroup := api.Group("/player")
	playerGroup.Use(authMiddleware)
	{
		playerhandler.RegisterProfileRoutes(playerGroup, playerSvc, authMiddleware)
		playerhandler.RegisterOrderRoutes(playerGroup, orderSvc, authMiddleware)
		playerhandler.RegisterEarningsRoutes(playerGroup, earningsSvc, authMiddleware)
		playerhandler.RegisterCommissionRoutes(playerGroup, commissionSvc, authMiddleware)
		playerhandler.RegisterGiftRoutes(playerGroup, giftSvc, authMiddleware)
	}

	if cfg.EnableSwagger {
		log.Println("swagger endpoint enabled at /swagger")
		// Serve embedded OpenAPI v3 at /swagger and /swagger.json
		handler.RegisterSwagger(router)
		// Serve gin-swagger UI backed by /swagger.json for compatibility
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger.json")))
	} else {
		log.Println("swagger endpoint disabled by configuration")
	}

	// RBAC - 权限服务（RoleRepository 已在 AdminService 创建时初始化）
	permRepo := permissionrepo.NewPermissionRepository(orm)
	permService := permissionservice.NewPermissionService(permRepo, cacheClient)
	roleSvc := roleservice.NewRoleService(roleRepo, cacheClient)

	// 权限中间件
	permMiddleware := middleware.NewPermissionMiddleware(jwtMgr, permService, roleSvc)

	// Register admin routes under versioned prefix: /api/v1/admin（使用新的权限中间件）
	adminhandler.RegisterRoutes(api, adminSvc, permMiddleware)

	// Stats routes（使用新的权限中间件）
	statsSvc := statsservice.NewStatsService(statsrepo.NewStatsRepository(orm))
	adminhandler.RegisterStatsRoutes(api, statsSvc, permMiddleware)

	// System info routes（使用新的权限中间件）
	adminhandler.RegisterSystemRoutes(api, cfg, sqlDB, cacheClient, permMiddleware)

	// 注册角色和权限管理路由（使用细粒度权限控制）
	roleHandler := adminhandler.NewRoleHandler(roleSvc)
	permHandler := adminhandler.NewPermissionHandler(permService)

	// RBAC routes
	rbacGroup := api.Group("/admin")
	rbacGroup.Use(permMiddleware.RequireAuth()) // 所有 RBAC 接口需要认证
	{
		// 角色管理 - 使用细粒度权限
		rbacGroup.GET("/roles", permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/roles"), roleHandler.ListRoles)
		rbacGroup.GET("/roles/:id", permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/roles/:id"), roleHandler.GetRole)
		rbacGroup.POST("/roles", permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/roles"), roleHandler.CreateRole)
		rbacGroup.PUT("/roles/:id", permMiddleware.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/roles/:id"), roleHandler.UpdateRole)
		rbacGroup.DELETE("/roles/:id", permMiddleware.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/roles/:id"), roleHandler.DeleteRole)
		rbacGroup.PUT("/roles/:id/permissions", permMiddleware.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/roles/:id/permissions"), roleHandler.AssignPermissions)
		rbacGroup.POST("/roles/assign-user", permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/roles/assign-user"), roleHandler.AssignRolesToUser)
		rbacGroup.GET("/users/:id/roles", permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/users/:id/roles"), roleHandler.GetUserRoles)

		// 权限管理 - 使用细粒度权限
		rbacGroup.GET("/permissions", permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/permissions"), permHandler.ListPermissions)
		rbacGroup.GET("/permissions/groups", permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/permissions/groups"), permHandler.GetPermissionGroups)
		rbacGroup.GET("/permissions/:id", permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/permissions/:id"), permHandler.GetPermission)
		rbacGroup.POST("/permissions", permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/permissions"), permHandler.CreatePermission)
		rbacGroup.PUT("/permissions/:id", permMiddleware.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/permissions/:id"), permHandler.UpdatePermission)
		rbacGroup.DELETE("/permissions/:id", permMiddleware.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/permissions/:id"), permHandler.DeletePermission)
		rbacGroup.GET("/roles/:id/permissions", permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/roles/:id/permissions"), permHandler.GetRolePermissions)
		rbacGroup.GET("/users/:id/permissions", permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/users/:id/permissions"), permHandler.GetUserPermissions)
	}

	// Commission management routes (admin)
	adminhandler.RegisterCommissionRoutes(rbacGroup, commissionSvc, settlementScheduler)

	// Service Item management routes (admin) - 统一管理护航服务和礼物
	adminhandler.RegisterServiceItemRoutes(rbacGroup, serviceItemSvc)

	// Withdraw management routes (admin) - 提现审核管理
	adminhandler.RegisterWithdrawRoutes(rbacGroup, withdrawRepo)

	// Dashboard routes (admin) - 数据统计和Dashboard
	adminhandler.RegisterDashboardRoutes(rbacGroup, userRepo, playerRepo, orderRepo, withdrawRepo, serviceItemRepo, commissionRepo)

	// Stats routes (admin) - 统计分析
	adminhandler.RegisterStatsAnalysisRoutes(rbacGroup, orderRepo, commissionRepo, serviceItemRepo)

	// Ranking Commission routes (admin) - 排名抽成配置
	adminhandler.RegisterRankingCommissionRoutes(rbacGroup, rankingCommissionRepo)

	// 同步 API 路由到权限表（开发环境自动同步）
	if os.Getenv("APP_ENV") != "production" || os.Getenv("SYNC_API_PERMISSIONS") == "true" {
		log.Println("同步 API 权限到数据库...")
		syncConfig := middleware.APISyncConfig{
			GroupFilter: "/api/v1/admin",
			SkipPaths: []string{
				"/api/v1/health",
				"/api/v1/metrics",
				"/api/v1/swagger",
			},
			DryRun: false,
		}
		if err := middleware.SyncAPIPermissions(router, permService, syncConfig); err != nil {
			log.Printf("同步权限失败: %v", err)
		}

		// 权限同步后，为默认角色分配权限
		log.Println("为默认角色分配权限...")
		if err := assignDefaultRolePermissions(context.Background(), roleSvc, permService); err != nil {
			log.Printf("分配默认权限失败: %v", err)
		}
	}

	addr := fmt.Sprintf(":%s", cfg.Port)

	server := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("user-service listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Println("server stopped gracefully")
}

// assignDefaultRolePermissions 为默认角色（admin 和 super_admin）分配所有管理权限。
func assignDefaultRolePermissions(ctx context.Context, roleSvc *roleservice.RoleService, permService *permissionservice.PermissionService) error {
	// 获取所有权限
	allPermissions, err := permService.ListPermissions(ctx)
	if err != nil {
		return fmt.Errorf("failed to list permissions: %w", err)
	}

	if len(allPermissions) == 0 {
		log.Println("没有权限需要分配，跳过")
		return nil
	}

	// 提取所有权限 ID
	permissionIDs := make([]uint64, 0, len(allPermissions))
	for _, perm := range allPermissions {
		permissionIDs = append(permissionIDs, perm.ID)
	}

	// 为 admin 和 super_admin 角色分配所有权限
	roleSlugs := []string{
		string(model.RoleSlugSuperAdmin),
		string(model.RoleSlugAdmin),
	}

	for _, roleSlug := range roleSlugs {
		role, err := roleSvc.GetRoleBySlug(ctx, roleSlug)
		if err != nil {
			log.Printf("警告：未找到角色 %s，跳过: %v", roleSlug, err)
			continue
		}

		// 分配权限（替换现有权限）
		if err := roleSvc.AssignPermissionsToRole(ctx, role.ID, permissionIDs); err != nil {
			log.Printf("警告：为角色 %s 分配权限失败: %v", roleSlug, err)
			continue
		}

		log.Printf("已为角色 %s (id=%d) 分配 %d 个权限", roleSlug, role.ID, len(permissionIDs))
	}

	return nil
}

func resolveGinMode() string {
	if mode := os.Getenv("GIN_MODE"); mode != "" {
		return mode
	}
	if env := os.Getenv("APP_ENV"); env == "production" {
		return gin.ReleaseMode
	}
	return gin.DebugMode
}
