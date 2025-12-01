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

	"gamelink/internal/auth"
	"gamelink/internal/cache"
	"gamelink/internal/config"
	"gamelink/internal/handler"
	adminhandler "gamelink/internal/handler/admin"
	"gamelink/internal/handler/middleware"
	notificationhandler "gamelink/internal/handler/notification"
	"gamelink/internal/lifecycle"
	"gamelink/internal/model"
	commissionrepo "gamelink/internal/repository/commission"
	orderrepo "gamelink/internal/repository/implementations"
	menurepo "gamelink/internal/repository/menu"
	permissionrepo "gamelink/internal/repository/permission"
	playerrepo "gamelink/internal/repository/player"
	rankingrepo "gamelink/internal/repository/ranking"
	serviceitemrepo "gamelink/internal/repository/serviceitem"
	rolerepo "gamelink/internal/repository/role"
	statsrepo "gamelink/internal/repository/stats"
	userrepo "gamelink/internal/repository/user"
	withdrawrepo "gamelink/internal/repository/withdraw"
	adminservice "gamelink/internal/service/admin"
	authservice "gamelink/internal/service/auth"
	menuservice "gamelink/internal/service/menu"
	permissionservice "gamelink/internal/service/permission"
	roleservice "gamelink/internal/service/role"
	statsservice "gamelink/internal/service/stats"
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

	// 注册全局中间件（按顺序执行）
	r.engine.Use(middleware.RequestID())
	r.engine.Use(middleware.SlogLogger())                                   // 结构化访问日志
	r.engine.Use(middleware.MetricsMiddleware())                            // HTTP 指标
	r.engine.Use(middleware.RateLimit(middleware.DefaultRateLimitConfig())) // 限流中间件
	r.engine.Use(middleware.Crypto(r.cfg.Crypto))                           // 请求解密
	r.engine.Use(middleware.ErrorMap())                                     // 统一错误映射
	r.engine.Use(middleware.Recovery())                                     // 统一JSON恢复中间件
	r.engine.Use(middleware.CORS())                                         // CORS中间件

	// 初始化认证服务
	r.setupAuth()

	// 初始化业务服务
	r.setupServices()

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
		// Serve embedded OpenAPI v3 at /swagger and /swagger.json
		handler.RegisterSwagger(r.engine)
		// Serve gin-swagger UI backed by /swagger.json for compatibility
		r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger.json")))
	} else {
		log.Println("swagger endpoint disabled by configuration")
	}
}

// registerAdminRoutes 注册管理端路由
func (r *Router) registerAdminRoutes(api *gin.RouterGroup) {
	// RBAC - 权限服务
	permRepo := permissionrepo.NewPermissionRepository(r.orm)
	permService := permissionservice.NewPermissionService(permRepo, r.cacheClient)
	roleSvc := roleservice.NewRoleService(rolerepo.NewRoleRepository(r.orm), r.cacheClient)
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
	menuSvc := menuservice.NewService(menurepo.NewMenuRepository(r.orm))
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
	playerRepo := playerrepo.NewPlayerRepository(r.orm)
	orderRepo := orderrepo.NewOrderRepository(r.orm)
	commissionRepo := commissionrepo.NewCommissionRepository(r.orm)
	serviceItemRepo := serviceitemrepo.NewServiceItemRepository(r.orm)
	adminhandler.RegisterDashboardRoutes(rbacGroup, userRepo, playerRepo, orderRepo, withdrawRepo, serviceItemRepo, commissionRepo)

	// Ranking Commission routes
	rankingCommissionRepo := rankingrepo.NewRankingCommissionRepository(r.orm)
	adminhandler.RegisterRankingCommissionRoutes(rbacGroup, rankingCommissionRepo)
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
