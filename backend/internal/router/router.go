package router

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"gamelink/pkg/auth"
	"gamelink/pkg/cache"
	"gamelink/pkg/config"
	"gamelink/internal/handler"
	adminhandler "gamelink/internal/handler/admin"
	"gamelink/internal/handler/middleware"
	notificationhandler "gamelink/internal/handler/notification"
	"gamelink/internal/ws"
	"gamelink/pkg/lifecycle"
	"gamelink/internal/model"
	commissionrepo "gamelink/internal/repository/commission"
	orderrepo "gamelink/internal/repository/implementations"
	adminrepo "gamelink/internal/repository/admin"
	userrepo "gamelink/internal/repository/user"
	rankingrepo "gamelink/internal/repository/ranking"
	serviceitemrepo "gamelink/internal/repository/serviceitem"
	statsrepo "gamelink/internal/repository/stats"
	withdrawrepo "gamelink/internal/repository/withdraw"
	adminservice "gamelink/internal/service/admin"
	authservice "gamelink/internal/service/auth"
	menuservice "gamelink/internal/service/admin"
	permissionservice "gamelink/internal/service/admin"
	roleservice "gamelink/internal/service/admin"
	statsservice "gamelink/internal/service/admin"
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
	r.engine.Use(middleware.SlogLogger())                                                       // 结构化访问日志
	r.engine.Use(middleware.MetricsMiddleware(r.services.realtimeSvc))                          // HTTP 指标，传入monitor service
	r.engine.Use(middleware.RateLimit(middleware.DefaultRateLimitConfig()))                     // 限流中间件
	r.engine.Use(middleware.Crypto(r.cfg.Crypto))                                               // 请求解密
	r.engine.Use(middleware.ErrorMap())                                                         // 统一错误映射
	r.engine.Use(middleware.Recovery())                                                         // 统一JSON恢复中间件
	r.engine.Use(middleware.CORS())                                                             // CORS中间件

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
	r.authMiddleware = middleware.JWTAuth()
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
	permService := permissionservice.NewPermissionService(permRepo, r.cacheClient)
	roleSvc := roleservice.NewRoleService(adminrepo.NewRoleRepository(r.orm), r.cacheClient)
	r.permMiddleware = middleware.NewPermissionMiddleware(r.jwtMgr, permService, roleSvc)

	// Notification center routes
	notificationhandler.RegisterRoutes(api, r.services.notificationSvc, r.authMiddleware)

	// Register admin routes under versioned prefix: /api/v1/admin
	rbacGroup := api.Group("/admin")
	adminhandler.RegisterRoutes(rbacGroup, r.adminSvc, r.permMiddleware)

	// Stats routes
	statsSvc := statsservice.NewStatsService(statsrepo.NewStatsRepository(r.orm))
	adminhandler.RegisterStatsRoutes(rbacGroup, statsSvc, r.permMiddleware)

	// System info routes
	adminhandler.RegisterSystemRoutes(api, r.cfg, r.sqlDB, r.cacheClient, r.permMiddleware)

	// 注册角色和权限管理路由
	r.registerRBACRoutes(rbacGroup, roleSvc, permService)

	// Admin 端业务路由
	r.registerAdminBusinessRoutes(rbacGroup)

	// 同步 API 路由到权限表
	r.syncAPIPermissions(permService, roleSvc)
}

// registerRBACRoutes 注册 RBAC 相关路由
func (r *Router) registerRBACRoutes(rbacGroup *gin.RouterGroup, roleSvc *roleservice.RoleService, permService *permissionservice.PermissionService) {
	roleHandler := adminhandler.NewRoleHandler(roleSvc)
	permHandler := adminhandler.NewPermissionHandler(permService)
	menuSvc := menuservice.NewMenuService(adminrepo.NewMenuRepository(r.orm))
	menuHandler := adminhandler.NewMenuHandler(menuSvc)
	rbacGroup.Use(r.permMiddleware.RequireAuth()) // 所有 RBAC 接口需要认证
	{
		// 角色管理
		rbacGroup.GET("/roles", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/roles"), roleHandler.ListRoles)
		rbacGroup.GET("/roles/:id", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/roles/:id"), roleHandler.GetRole)
		rbacGroup.POST("/roles", r.permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/roles"), roleHandler.CreateRole)
		rbacGroup.PUT("/roles/:id", r.permMiddleware.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/roles/:id"), roleHandler.UpdateRole)
		rbacGroup.DELETE("/roles/:id", r.permMiddleware.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/roles/:id"), roleHandler.DeleteRole)
		rbacGroup.PUT("/roles/:id/permissions", r.permMiddleware.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/roles/:id/permissions"), roleHandler.AssignPermissions)
		rbacGroup.POST("/roles/assign-user", r.permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/roles/assign-user"), roleHandler.AssignRolesToUser)
		rbacGroup.GET("/users/:id/roles", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/users/:id/roles"), roleHandler.GetUserRoles)

		// 权限管理
		rbacGroup.GET("/permissions/me", permHandler.GetCurrentUserPermissions)
		rbacGroup.GET("/permissions", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/permissions"), permHandler.ListPermissions)
		rbacGroup.GET("/permissions/groups", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/permissions/groups"), permHandler.GetPermissionGroups)
		rbacGroup.GET("/permissions/:id", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/permissions/:id"), permHandler.GetPermission)
		rbacGroup.POST("/permissions", r.permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/permissions"), permHandler.CreatePermission)
		rbacGroup.PUT("/permissions/:id", r.permMiddleware.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/permissions/:id"), permHandler.UpdatePermission)
		rbacGroup.DELETE("/permissions/:id", r.permMiddleware.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/permissions/:id"), permHandler.DeletePermission)
		rbacGroup.GET("/roles/:id/permissions", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/roles/:id/permissions"), permHandler.GetRolePermissions)
		rbacGroup.GET("/users/:id/permissions", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/users/:id/permissions"), permHandler.GetUserPermissions)

		// 菜单管理（动态路由配置）
		rbacGroup.GET("/menus", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/menus"), menuHandler.List)
		rbacGroup.POST("/menus", r.permMiddleware.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/menus"), menuHandler.Create)
		rbacGroup.GET("/menus/:id", r.permMiddleware.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/menus/:id"), menuHandler.Get)
		rbacGroup.PUT("/menus/:id", r.permMiddleware.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/menus/:id"), menuHandler.Update)
		rbacGroup.DELETE("/menus/:id", r.permMiddleware.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/menus/:id"), menuHandler.Delete)
	}
}

// registerAdminBusinessRoutes 注册管理端业务路由
func (r *Router) registerAdminBusinessRoutes(rbacGroup *gin.RouterGroup) {
	// Commission management routes
	adminhandler.RegisterCommissionRoutes(rbacGroup, r.services.commissionSvc, r.services.settlementScheduler)

	// Service Item management routes
	adminhandler.RegisterServiceItemRoutes(rbacGroup, r.services.serviceItemSvc)

	// Withdraw management routes
	withdrawRepo := withdrawrepo.NewWithdrawRepository(r.orm)
	adminhandler.RegisterWithdrawRoutes(rbacGroup, withdrawRepo)

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

	// Analytics routes (运营分析)
	r.registerAnalyticsRoutes(rbacGroup)

	// KPI routes (KPI 仪表板)
	r.registerKPIRoutes(rbacGroup)
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
func (r *Router) syncAPIPermissions(permService *permissionservice.PermissionService, roleSvc *roleservice.RoleService) {
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
