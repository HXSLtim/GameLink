package router

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"gamelink/internal/handler"
	adminhandler "gamelink/internal/handler/admin"
	"gamelink/internal/handler/middleware"
	notificationhandler "gamelink/internal/handler/notification"
	"gamelink/internal/model"
	adminrepo "gamelink/internal/repository/admin"
	commissionrepo "gamelink/internal/repository/commission"
	orderrepo "gamelink/internal/repository/implementations"
	rankingrepo "gamelink/internal/repository/ranking"
	serviceitemrepo "gamelink/internal/repository/serviceitem"
	settlementcompanyrepo "gamelink/internal/repository/settlementcompany"
	statsrepo "gamelink/internal/repository/stats"
	userrepo "gamelink/internal/repository/user"
	withdrawrepo "gamelink/internal/repository/withdraw"
	adminservice "gamelink/internal/service/admin"
	authservice "gamelink/internal/service/auth"
	settlementcompanysvc "gamelink/internal/service/settlementcompany"
	"gamelink/internal/ws"
	"gamelink/pkg/auth"
	"gamelink/pkg/cache"
	"gamelink/pkg/config"
	"gamelink/pkg/lifecycle"
)

// Router 包含所有路由配置和依赖
type Router struct {
	engine         *gin.Engine
	cfg            config.AppConfig
	orm            *gorm.DB
	sqlDB          *sql.DB
	cacheClient    cache.Cache
	adminSvc       *adminservice.AdminService
	jwtMgr         *auth.JWTManager
	authSvc        *authservice.AuthService
	authMiddleware gin.HandlerFunc
	permMiddleware *middleware.PermissionMiddleware
	services       *appServices
	lifecycle      *lifecycle.Registry
}

// NewRouter 创建新的路由器实例
func NewRouter(cfg config.AppConfig, orm *gorm.DB, sqlDB *sql.DB, cacheClient cache.Cache, adminSvc *adminservice.AdminService, lifecycleRegistry *lifecycle.Registry) *Router {
	return &Router{
		cfg:         cfg,
		orm:         orm,
		sqlDB:       sqlDB,
		cacheClient: cacheClient,
		adminSvc:    adminSvc,
		lifecycle:   lifecycleRegistry,
	}
}

// Setup 初始化路由器和中间件
func (r *Router) Setup() *gin.Engine {
	gin.SetMode(resolveGinMode())

	// 创建路由引擎
	r.engine = gin.New()

	// 初始化认证服务
	r.setupAuth()

	// 初始化业务服务（必须在注册中间件之前，因为MetricsMiddleware需要monitor service）
	r.setupServices()

	// 注册全局中间件（按顺序执行）
	r.engine.Use(middleware.RequestID())
	r.engine.Use(middleware.SlogLogger())                                   // 结构化访问日志
	r.engine.Use(middleware.MetricsMiddleware(r.services.realtimeSvc))      // HTTP 指标，传入monitor service
	r.engine.Use(middleware.RateLimit(middleware.DefaultRateLimitConfig())) // 限流中间件
	r.engine.Use(middleware.Signature(r.cfg.Signature))                     // HMAC-SHA256签名验证
	r.engine.Use(middleware.Crypto(r.cfg.Crypto))                           // 请求解密
	r.engine.Use(middleware.ErrorMap())                                     // 统一错误映射
	r.engine.Use(middleware.Recovery())                                     // 统一JSON恢复中间件
	r.engine.Use(middleware.CORS())                                         // CORS中间件

	// 注册所有路由
	r.registerRoutes()

	return r.engine
}

// setupAuth 初始化认证相关服务
func (r *Router) setupAuth() {
	jwtSecret := strings.TrimSpace(r.cfg.Auth.JWTSecret)
	tokenTTL := time.Duration(r.cfg.Auth.TokenTTLHours) * time.Hour
	if tokenTTL <= 0 {
		tokenTTL = 24 * time.Hour
	}

	r.jwtMgr = auth.NewJWTManager(jwtSecret, tokenTTL)
	r.authSvc = authservice.NewAuthService(userrepo.NewUserRepository(r.orm), r.jwtMgr)
	r.authMiddleware = middleware.JWTAuth(jwtSecret)
}

// setupServices 初始化业务服务
func (r *Router) setupServices() {
	services := initServices(r.orm, r.cacheClient)
	r.services = services

	if r.lifecycle == nil {
		services.settlementScheduler.Start()
		services.chatRetention.Start()
		// Start monitor services
		go services.wsHub.Run()
		services.realtimeSvc.Start(context.Background())
		return
	}

	r.lifecycle.RegisterHook("scheduler:settlement", func(context.Context) error {
		services.settlementScheduler.Start()
		return nil
	}, func(context.Context) error {
		services.settlementScheduler.Stop()
		return nil
	})

	r.lifecycle.RegisterHook("scheduler:chat-retention", func(context.Context) error {
		services.chatRetention.Start()
		return nil
	}, func(context.Context) error {
		services.chatRetention.Stop()
		return nil
	})

	// WebSocket Hub 生命周期管理
	r.lifecycle.RegisterHook("websocket:hub", func(ctx context.Context) error {
		go services.wsHub.Run()
		return nil
	}, func(ctx context.Context) error {
		// Hub 会在所有连接关闭后自动停止
		return nil
	})

	// 实时监控服务生命周期管理
	r.lifecycle.RegisterHook("monitor:realtime", func(ctx context.Context) error {
		services.realtimeSvc.Start(ctx)
		return nil
	}, func(ctx context.Context) error {
		services.realtimeSvc.Stop()
		return nil
	})
}

// registerRoutes 注册所有路由
func (r *Router) registerRoutes() {
	// 根路由和健康检查
	handler.RegisterRoot(r.engine)
	handler.RegisterHealth(r.engine)

	// 版本化 API 分组
	api := r.engine.Group("/api/v1")
	handler.RegisterRoot(api)
	handler.RegisterHealth(api)

	// 认证路由
	handler.RegisterAuthRoutes(api, r.authSvc)

	// 用户端路由
	registerUserRoutes(api, r.authMiddleware, r.services)

	// 陪玩端路由
	registerPlayerRoutes(api, r.authMiddleware, r.services)

	// Swagger 路由
	r.registerSwaggerRoutes()

	// 管理端路由
	r.registerAdminRoutes(api)
}

// registerSwaggerRoutes 注册 Swagger 相关路由
func (r *Router) registerSwaggerRoutes() {
	if r.cfg.EnableSwagger {
		log.Println("swagger endpoint enabled at /swagger")
		// Serve swagger.json file
		r.engine.StaticFile("/swagger.json", "./docs/swagger.json")
		// Serve gin-swagger UI backed by /swagger.json for compatibility
		r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger.json")))
	} else {
		log.Println("swagger endpoint disabled by configuration")
	}
}

// registerAdminRoutes 注册管理端路由
func (r *Router) registerAdminRoutes(api *gin.RouterGroup) {
	// RBAC - 权限服务
	permRepo := adminrepo.NewPermissionRepository(r.orm)
	permService := adminservice.NewPermissionService(permRepo, r.cacheClient)
	roleSvc := adminservice.NewRoleService(adminrepo.NewRoleRepository(r.orm), r.cacheClient)
	r.permMiddleware = middleware.NewPermissionMiddleware(r.jwtMgr, permService, roleSvc)

	// Notification center routes
	notificationhandler.RegisterRoutes(api, r.services.notificationSvc, r.authMiddleware)

	// Register admin routes under versioned prefix: /api/v1/admin
	rbacGroup := api.Group("/admin")

	// CSRF 保护配置
	// 根据环境变量决定是否启用 CSRF
	csrfEnabled := os.Getenv("CSRF_ENABLED") == "true"
	if csrfEnabled {
		// 生产环境使用更严格的配置
		isProduction := os.Getenv("APP_ENV") == "production"
		csrfConfig := middleware.DefaultCSRFConfig
		csrfConfig.CookieSecure = isProduction
		if isProduction {
			csrfConfig.CookieSameSite = http.SameSiteStrictMode
		}
		rbacGroup.Use(middleware.CSRF(csrfConfig))
	}

	// Stats routes (需要先创建statsSvc，因为RegisterRoutes需要它)
	statsSvc := adminservice.NewStatsService(statsrepo.NewStatsRepository(r.orm))

	adminhandler.RegisterRoutes(rbacGroup, r.adminSvc, statsSvc, r.permMiddleware, adminhandler.WithSensitiveWordService(r.services.sensitiveWordSvc))
	adminhandler.RegisterStatsRoutes(rbacGroup, statsSvc, r.permMiddleware)
	adminhandler.RegisterReviewStatsRoutes(rbacGroup, r.services.reviewStatsSvc, r.permMiddleware)
	adminhandler.RegisterReviewSettingsRoutes(rbacGroup, r.services.reviewSettingsSvc, r.permMiddleware)

	// 创建菜单服务用于批量同步
	menuSvc := adminservice.NewMenuService(adminrepo.NewMenuRepository(r.orm))

	// 注册同步专用路由（不受限流限制，用于前端初始化）
	adminhandler.RegisterSyncRoutesWithServices(rbacGroup, roleSvc, permService, menuSvc, r.orm, r.permMiddleware)

	// System info routes
	adminhandler.RegisterSystemRoutes(api, r.cfg, r.sqlDB, r.cacheClient, r.orm, r.permMiddleware)

	// 注册角色和权限管理路由
	r.registerRBACRoutes(rbacGroup, roleSvc, permService)

	// Admin 端业务路由
	r.registerAdminBusinessRoutes(rbacGroup)

	// 结算公司管理路由
	settlementCompanyRepo := settlementcompanyrepo.NewSettlementCompanyRepository(r.orm)
	settlementCompanySvc := settlementcompanysvc.NewSettlementCompanyService(settlementCompanyRepo, nil)
	adminhandler.RegisterSettlementCompanyRoutes(rbacGroup, settlementCompanySvc, r.permMiddleware)

	// 内容管理路由
	r.registerContentRoutes(rbacGroup)

	// 同步 API 路由到权限表
	r.syncAPIPermissions(permService, roleSvc)
}

// registerRBACRoutes 注册 RBAC 相关路由
func (r *Router) registerRBACRoutes(rbacGroup *gin.RouterGroup, roleSvc *adminservice.RoleService, permService *adminservice.PermissionService) {
	roleHandler := adminhandler.NewRoleHandler(roleSvc)
	// roleBatchHandler := adminhandler.NewRoleBatchHandler(roleSvc) // Role batch handler to be implemented when needed
	permHandler := adminhandler.NewPermissionHandlerWithRoleService(permService, roleSvc)
	menuSvc := adminservice.NewMenuService(adminrepo.NewMenuRepository(r.orm))
	menuHandler := adminhandler.NewMenuHandlerWithRoleService(menuSvc, permService, roleSvc)

	// 注意：同步专用路由 /sync/roles 和 /sync/permissions 已在 adminhandler.RegisterRoutes 中注册
	// 这里只注册常规 RBAC 路由（需要认证和权限检查）
	rbacGroup.Use(r.permMiddleware.RequireAuth())
	{
		// 当前用户权限和菜单 API（无需额外权限检查，只需认证）
		// GET /api/admin/me/permissions - 获取当前用户权限码列表，超级管理员返回 ['*']
		// GET /api/admin/me/menus - 获取当前用户可访问的菜单
		rbacGroup.GET("/me/permissions", permHandler.GetCurrentUserPermissions)
		rbacGroup.GET("/me/menus", menuHandler.ListMyMenus)

		// 角色管理
		rbacGroup.GET("/roles", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/roles"), roleHandler.ListRoles)
		rbacGroup.GET("/roles/:id", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/roles/:id"), roleHandler.GetRole)
		rbacGroup.POST("/roles", r.permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/roles"), roleHandler.CreateRole)
		rbacGroup.PUT("/roles/:id", r.permMiddleware.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/roles/:id"), roleHandler.UpdateRole)
		rbacGroup.DELETE("/roles/:id", r.permMiddleware.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/roles/:id"), roleHandler.DeleteRole)
		// 角色权限分配 API
		rbacGroup.GET("/roles/:id/permissions", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/roles/:id/permissions"), roleHandler.GetRolePermissionIDs)
		rbacGroup.PUT("/roles/:id/permissions/batch", r.permMiddleware.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/roles/:id/permissions/batch"), roleHandler.AssignPermissions)
		rbacGroup.POST("/roles/:id/permissions/:pid", r.permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/roles/:id/permissions/:pid"), roleHandler.AddPermissionToRole)
		rbacGroup.DELETE("/roles/:id/permissions/:pid", r.permMiddleware.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/roles/:id/permissions/:pid"), roleHandler.RemovePermissionFromRole)
		// 用户角色分配 API
		rbacGroup.POST("/roles/assign-user", r.permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/roles/assign-user"), roleHandler.AssignRolesToUser)
		rbacGroup.GET("/users/:id/roles", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/users/:id/roles"), roleHandler.GetUserRoles)
		rbacGroup.PUT("/users/:id/roles", r.permMiddleware.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/users/:id/roles"), roleHandler.UpdateUserRoles)
		rbacGroup.PUT("/users/roles/batch", r.permMiddleware.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/users/roles/batch"), roleHandler.BatchAssignRolesToUsers)
		// Batch role operations (simplified API) - to be implemented when RoleBatchHandler is created
		// For now, using legacy batch endpoints below
		// rbacGroup.POST("/roles/batch/delete", roleBatchHandler.BatchDeleteRoles)
		// rbacGroup.POST("/roles/batch/assign-permissions", roleBatchHandler.BatchAssignPermissions)
		// 批量角色操作 API（旧接口，保留兼容性）
		rbacGroup.DELETE("/roles/batch", r.permMiddleware.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/roles/batch"), roleHandler.BatchDeleteRoles)
		rbacGroup.PUT("/roles/batch/permissions", r.permMiddleware.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/roles/batch/permissions"), roleHandler.BatchAssignPermissionsToRoles)

		// 权限管理
		rbacGroup.GET("/permissions/me", permHandler.GetCurrentUserPermissions) // 保留旧路径兼容性
		rbacGroup.GET("/permissions", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/permissions"), permHandler.ListPermissions)
		rbacGroup.GET("/permissions/groups", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/permissions/groups"), permHandler.GetPermissionGroups)
		rbacGroup.GET("/permissions/tree", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/permissions/tree"), permHandler.GetPermissionTree)
		rbacGroup.GET("/permissions/tree/grouped", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/permissions/tree/grouped"), permHandler.GetPermissionTreeByGroup)
		rbacGroup.GET("/permissions/:id", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/permissions/:id"), permHandler.GetPermission)
		rbacGroup.POST("/permissions", r.permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/permissions"), permHandler.CreatePermission)
		rbacGroup.PUT("/permissions/:id", r.permMiddleware.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/permissions/:id"), permHandler.UpdatePermission)
		rbacGroup.PATCH("/permissions/:id", r.permMiddleware.RequirePermission(model.HTTPMethodPATCH, "/api/v1/admin/permissions/:id"), permHandler.PatchPermission)
		rbacGroup.DELETE("/permissions/:id", r.permMiddleware.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/permissions/:id"), permHandler.DeletePermission)
		// Note: /roles/:id/permissions is already registered above with roleHandler.GetRolePermissionIDs
		rbacGroup.GET("/users/:id/permissions", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/users/:id/permissions"), permHandler.GetUserPermissions)
		// 批量权限操作 API
		rbacGroup.POST("/permissions/batch/delete", r.permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/permissions/batch/delete"), permHandler.BatchDeletePermissions)
		rbacGroup.DELETE("/permissions/batch", r.permMiddleware.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/permissions/batch"), permHandler.BatchDelete)

		// 菜单管理（动态路由配置）
		rbacGroup.GET("/menus/me", menuHandler.ListMyMenus) // 保留旧路径兼容性
		rbacGroup.GET("/menus", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/menus"), menuHandler.List)
		rbacGroup.POST("/menus", r.permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/menus"), menuHandler.Create)
		rbacGroup.GET("/menus/:id", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/menus/:id"), menuHandler.Get)
		rbacGroup.PUT("/menus/:id", r.permMiddleware.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/menus/:id"), menuHandler.Update)
		rbacGroup.DELETE("/menus/:id", r.permMiddleware.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/menus/:id"), menuHandler.Delete)
		// Batch menu operations (new format using POST) - to be implemented in MenuHandler
		// For now, using legacy batch endpoints below
		// rbacGroup.POST("/menus/batch/delete", menuHandler.BatchDeleteMenus)
		// rbacGroup.POST("/menus/batch/status", menuHandler.BatchUpdateMenuStatus)
		// rbacGroup.POST("/menus/batch/order", menuHandler.BatchUpdateMenuOrder)
		// 旧批量菜单操作 API (保持向后兼容)
		rbacGroup.DELETE("/menus/batch", r.permMiddleware.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/menus/batch"), menuHandler.BatchDelete)
		rbacGroup.PUT("/menus/batch/status", r.permMiddleware.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/menus/batch/status"), menuHandler.BatchUpdateStatus)
		rbacGroup.PUT("/menus/batch/sort", r.permMiddleware.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/menus/batch/sort"), menuHandler.BatchUpdateSort)
	}
}

// registerAdminBusinessRoutes 注册管理端业务路由
func (r *Router) registerAdminBusinessRoutes(rbacGroup *gin.RouterGroup) {
	// Dispute management routes
	disputeHandler := adminhandler.NewDisputeHandler(r.services.disputeSvc)
	disputeGroup := rbacGroup.Group("/disputes")
	disputeGroup.Use(r.permMiddleware.RequireAuth())
	{
		disputeGroup.GET("", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/disputes"), disputeHandler.ListDisputes)
		disputeGroup.GET("/stats", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/disputes/stats"), disputeHandler.GetDisputeStats)
		disputeGroup.GET("/pending", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/disputes/pending"), disputeHandler.ListPendingDisputes)
		disputeGroup.GET("/:id", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/disputes/:id"), disputeHandler.GetDisputeDetail)
		disputeGroup.POST("/:id/assign", r.permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/disputes/:id/assign"), disputeHandler.AssignDispute)
		disputeGroup.POST("/:id/rollback", r.permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/disputes/:id/rollback"), disputeHandler.RollbackAssignment)
		disputeGroup.POST("/:id/resolve", r.permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/disputes/:id/resolve"), disputeHandler.ResolveDispute)
		// Batch operations
		disputeGroup.POST("/batch/assign", r.permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/disputes/batch/assign"), disputeHandler.BatchAssignDisputes)
		disputeGroup.PUT("/batch/status", r.permMiddleware.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/disputes/batch/status"), disputeHandler.BatchUpdateDisputesStatus)
		disputeGroup.POST("/batch/close", r.permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/disputes/batch/close"), disputeHandler.BatchCloseDisputes)
	}

	// Commission management routes
	adminhandler.RegisterCommissionRoutes(rbacGroup, r.services.commissionSvc, r.services.settlementScheduler)

	// Service Item management routes
	adminhandler.RegisterServiceItemRoutes(rbacGroup, r.services.serviceItemSvc)

	// Withdraw management routes
	withdrawRepo := withdrawrepo.NewWithdrawRepository(r.orm)
	adminhandler.RegisterWithdrawRoutes(rbacGroup, withdrawRepo, r.services.withdrawRoutingSvc)

	// Dashboard routes
	userRepo := userrepo.NewUserRepository(r.orm)
	playerRepo := userrepo.NewPlayerRepository(r.orm)
	orderRepo := orderrepo.NewOrderRepository(r.orm)
	commissionRepo := commissionrepo.NewCommissionRepository(r.orm)
	serviceItemRepo := serviceitemrepo.NewServiceItemRepository(r.orm)
	adminhandler.RegisterDashboardRoutes(rbacGroup, userRepo, playerRepo, orderRepo, withdrawRepo, serviceItemRepo, commissionRepo)

	// Ranking Commission routes
	rankingCommissionRepo := rankingrepo.NewRankingCommissionRepository(r.orm)
	adminhandler.RegisterRankingCommissionRoutes(rbacGroup, rankingCommissionRepo)

	// User Tag routes (用户标签管理)
	adminhandler.RegisterTagRoutes(rbacGroup, r.services.tagSvc)

	// User Batch Operation routes (用户批量操作)
	adminhandler.RegisterBatchRoutes(rbacGroup, r.services.batchSvc)

	// Monitor routes (实时监控)
	r.registerMonitorRoutes(rbacGroup)

	// Routing Rule routes (分流规则管理)
	adminhandler.RegisterRoutingRuleRoutes(rbacGroup, r.services.routingRuleSvc, r.permMiddleware)

	// Analytics routes (运营分析)
	r.registerAnalyticsRoutes(rbacGroup)

	// KPI routes (KPI 仪表板)
	r.registerKPIRoutes(rbacGroup)

	// Statistics routes (统计指标管理)
	r.registerStatisticsRoutes(rbacGroup)

	// Player rank routes (陪玩师等级/认证管理)
	r.registerPlayerRankRoutes(rbacGroup)

	// Order timeout routes (订单超时管理)
	r.registerOrderTimeoutRoutes(rbacGroup)

	// User block routes (用户拉黑管理)
	r.registerUserBlockRoutes(rbacGroup)

	// VIP routes (VIP会员管理)
	r.registerVipRoutes(rbacGroup)

	// Coupon routes (优惠券管理)
	r.registerCouponRoutes(rbacGroup)

	// Recharge routes (充值管理)
	r.registerRechargeRoutes(rbacGroup)

	// Activity routes (活动管理)
	r.registerActivityRoutes(rbacGroup)

	// Team routes (团队管理)
	r.registerTeamRoutes(rbacGroup)

	// Referral routes (推荐管理)
	r.registerReferralRoutes(rbacGroup)

	// Game Category routes (游戏分类管理)
	r.registerGameCategoryRoutes(rbacGroup)

	// Order Group routes (主订单/订单拆分管理)
	adminhandler.RegisterOrderGroupRoutes(rbacGroup, r.services.orderSvc, r.services.orderGroupRepo)
}

// registerMonitorRoutes 注册监控相关路由
func (r *Router) registerMonitorRoutes(rbacGroup *gin.RouterGroup) {
	monitorHandler := adminhandler.NewMonitorHandler(r.services.realtimeSvc, r.services.alertRepo)

	monitorGroup := rbacGroup.Group("/monitor")
	monitorGroup.Use(r.permMiddleware.RequireAuth())
	{
		// 实时监控 API
		monitorGroup.GET("/system-status", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/monitor/system-status"), monitorHandler.GetSystemStatus)
		monitorGroup.GET("/online-users", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/monitor/online-users"), monitorHandler.GetOnlineUsers)
		monitorGroup.GET("/order-queue", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/monitor/order-queue"), monitorHandler.GetOrderQueue)
		monitorGroup.GET("/alerts", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/monitor/alerts"), monitorHandler.GetAlerts)
		monitorGroup.PUT("/alerts/:id/read", r.permMiddleware.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/monitor/alerts/:id/read"), monitorHandler.MarkAlertRead)
		monitorGroup.PUT("/alerts/batch-read", r.permMiddleware.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/monitor/alerts/batch-read"), monitorHandler.BatchMarkAlertsRead)
	}

	// WebSocket 路由 (使用 ws.Handler 的认证机制)
	wsHandler := ws.NewHandler(r.services.wsHub)
	rbacGroup.GET("/ws/monitor", r.permMiddleware.RequireAuth(), wsHandler.ServeWS)
}

// registerAnalyticsRoutes 注册运营分析路由
func (r *Router) registerAnalyticsRoutes(rbacGroup *gin.RouterGroup) {
	analyticsHandler := adminhandler.NewAnalyticsHandler(r.services.analyticsSvc)

	analyticsGroup := rbacGroup.Group("/analytics")
	analyticsGroup.Use(r.permMiddleware.RequireAuth())
	{
		analyticsGroup.GET("/active-users", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/analytics/active-users"), analyticsHandler.GetActiveUsers)
		analyticsGroup.GET("/retention", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/analytics/retention"), analyticsHandler.GetRetention)
		analyticsGroup.GET("/payment", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/analytics/payment"), analyticsHandler.GetPaymentAnalytics)
		analyticsGroup.GET("/conversion", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/analytics/conversion"), analyticsHandler.GetConversionFunnel)
	}
}

// registerKPIRoutes 注册 KPI 仪表板路由
func (r *Router) registerKPIRoutes(rbacGroup *gin.RouterGroup) {
	kpiHandler := adminhandler.NewKPIHandler(r.services.kpiSvc)

	kpiGroup := rbacGroup.Group("/kpi")
	kpiGroup.Use(r.permMiddleware.RequireAuth())
	{
		// KPI 概览和趋势
		kpiGroup.GET("/overview", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/kpi/overview"), kpiHandler.GetOverview)
		kpiGroup.GET("/trend", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/kpi/trend"), kpiHandler.GetTrend)

		// KPI 目标管理
		kpiGroup.GET("/targets", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/kpi/targets"), kpiHandler.GetTargets)
		kpiGroup.POST("/targets", r.permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/kpi/targets"), kpiHandler.CreateTarget)
		kpiGroup.PUT("/targets/:id", r.permMiddleware.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/kpi/targets/:id"), kpiHandler.UpdateTarget)
		kpiGroup.DELETE("/targets/:id", r.permMiddleware.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/kpi/targets/:id"), kpiHandler.DeleteTarget)
	}
}

// syncAPIPermissions 同步 API 权限
func (r *Router) syncAPIPermissions(permService *adminservice.PermissionService, roleSvc *adminservice.RoleService) {
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
		if err := middleware.SyncAPIPermissions(r.engine, permService, syncConfig); err != nil {
			log.Printf("同步权限失败: %v", err)
		}

		// 权限同步后，为默认角色分配权限
		log.Println("为默认角色分配权限...")
		if err := AssignDefaultRolePermissions(context.Background(), roleSvc, permService); err != nil {
			log.Printf("分配默认权限失败: %v", err)
		}
	}
}

// registerStatisticsRoutes 注册统计指标路由
func (r *Router) registerStatisticsRoutes(rbacGroup *gin.RouterGroup) {
	statisticsHandler := adminhandler.NewStatisticsHandler(r.services.statisticsSvc, r.services.statisticsEvaluator)

	statsGroup := rbacGroup.Group("/statistics")
	statsGroup.Use(r.permMiddleware.RequireAuth())
	{
		// 刷新统计
		statsGroup.POST("/user/:id/refresh", r.permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/statistics/user/:id/refresh"), statisticsHandler.RefreshUserStatistics)
		statsGroup.POST("/player/:id/refresh", r.permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/statistics/player/:id/refresh"), statisticsHandler.RefreshPlayerStatistics)
		statsGroup.POST("/refresh-all", r.permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/statistics/refresh-all"), statisticsHandler.RefreshAllStatistics)

		// 标签评估
		statsGroup.GET("/user/:id/evaluate-tags", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/statistics/user/:id/evaluate-tags"), statisticsHandler.EvaluateUserTags)
		statsGroup.POST("/user/:id/sync-tags", r.permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/statistics/user/:id/sync-tags"), statisticsHandler.SyncUserTags)
	}
}

// registerContentRoutes 注册内容管理路由
func (r *Router) registerContentRoutes(rbacGroup *gin.RouterGroup) {
	contentHandler := adminhandler.NewContentHandler(
		r.services.adminFeedSvc,
		r.services.chatModerationSvc,
		r.services.feedReportSvc,
		r.services.contentStatsSvc,
	)
	categoryHandler := adminhandler.NewContentCategoryHandler(r.services.contentCategorySvc)

	adminhandler.RegisterContentRoutes(rbacGroup, contentHandler, categoryHandler, r.permMiddleware)

	// Sensitive word routes (敏感词管理)
	sensitiveWordHandler := adminhandler.NewSensitiveWordHandler(r.services.sensitiveWordSvc)
	adminhandler.RegisterSensitiveWordRoutes(rbacGroup, sensitiveWordHandler, r.permMiddleware)
}

// registerPlayerRankRoutes 注册陪玩师等级/认证管理路由
func (r *Router) registerPlayerRankRoutes(rbacGroup *gin.RouterGroup) {
	adminhandler.RegisterGameRankRoutes(rbacGroup, r.services.gameRankSvc, r.permMiddleware)
	adminhandler.RegisterPlayerRankRoutes(rbacGroup, r.services.playerRankSvc, r.permMiddleware)
	adminhandler.RegisterPlayerCertificationRoutes(rbacGroup, r.services.playerCertificationSvc, r.permMiddleware)
}

// registerOrderTimeoutRoutes 注册订单超时管理路由
func (r *Router) registerOrderTimeoutRoutes(rbacGroup *gin.RouterGroup) {
	adminhandler.RegisterOrderTimeoutRoutes(rbacGroup, r.services.orderTimeoutSvc, r.permMiddleware)
}

// registerUserBlockRoutes 注册用户拉黑管理路由
func (r *Router) registerUserBlockRoutes(rbacGroup *gin.RouterGroup) {
	adminhandler.RegisterUserBlockRoutes(rbacGroup, r.services.userBlockSvc, r.permMiddleware)
}

// registerVipRoutes 注册VIP管理路由
func (r *Router) registerVipRoutes(rbacGroup *gin.RouterGroup) {
	adminhandler.RegisterVipRoutes(rbacGroup, r.services.vipSvc, r.permMiddleware)
}

// registerCouponRoutes 注册优惠券管理路由
func (r *Router) registerCouponRoutes(rbacGroup *gin.RouterGroup) {
	adminhandler.RegisterCouponRoutes(rbacGroup, r.services.couponSvc, r.permMiddleware)
}

// registerRechargeRoutes 注册充值管理路由
func (r *Router) registerRechargeRoutes(rbacGroup *gin.RouterGroup) {
	adminhandler.RegisterRechargeRoutes(rbacGroup, r.services.rechargeSvc, r.permMiddleware)
}

// registerActivityRoutes 注册活动管理路由
func (r *Router) registerActivityRoutes(rbacGroup *gin.RouterGroup) {
	adminhandler.RegisterActivityRoutes(rbacGroup, r.services.activitySvc, r.permMiddleware)
}

// registerTeamRoutes 注册团队管理路由
func (r *Router) registerTeamRoutes(rbacGroup *gin.RouterGroup) {
	adminhandler.RegisterTeamRoutes(rbacGroup, r.services.teamSvc, r.permMiddleware)
}

// registerReferralRoutes 注册推荐管理路由
func (r *Router) registerReferralRoutes(rbacGroup *gin.RouterGroup) {
	adminhandler.RegisterReferralRoutes(rbacGroup, r.services.referralSvc, r.permMiddleware)
}

// registerGameCategoryRoutes 注册游戏分类管理路由
func (r *Router) registerGameCategoryRoutes(rbacGroup *gin.RouterGroup) {
	categoryHandler := adminhandler.NewGameCategoryHandler(r.adminSvc)

	categoryGroup := rbacGroup.Group("/game-categories")
	categoryGroup.Use(r.permMiddleware.RequireAuth())
	{
		// 游戏分类 CRUD
		categoryGroup.GET("", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/game-categories"), categoryHandler.ListCategories)
		categoryGroup.POST("", r.permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/game-categories"), categoryHandler.CreateCategory)
		categoryGroup.GET("/:id", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/game-categories/:id"), categoryHandler.GetCategory)
		categoryGroup.PUT("/:id", r.permMiddleware.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/game-categories/:id"), categoryHandler.UpdateCategory)
		categoryGroup.DELETE("/:id", r.permMiddleware.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/game-categories/:id"), categoryHandler.DeleteCategory)

		// 批量操作
		categoryGroup.POST("/batch/status", r.permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/game-categories/batch/status"), categoryHandler.BatchUpdateStatus)
		categoryGroup.POST("/batch/delete", r.permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/game-categories/batch/delete"), categoryHandler.BatchDeleteCategories)
	}
}

// resolveGinMode 解析 Gin 运行模式
func resolveGinMode() string {
	if mode := os.Getenv("GIN_MODE"); mode != "" {
		return mode
	}
	if env := os.Getenv("APP_ENV"); env == "production" {
		return gin.ReleaseMode
	}
	return gin.DebugMode
}
