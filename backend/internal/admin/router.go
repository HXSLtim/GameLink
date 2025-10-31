package admin

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/config"
	mw "gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	"gamelink/internal/service"
)

// RegisterRoutes 注册后台管理相关路由。
// 使用细粒度权限控制（method+path 级别）。
func RegisterRoutes(router gin.IRouter, svc *service.AdminService, pm *mw.PermissionMiddleware) {
	gameHandler := NewGameHandler(svc)
	userHandler := NewUserHandler(svc)
	playerHandler := NewPlayerHandler(svc)
	orderHandler := NewOrderHandler(svc)
	paymentHandler := NewPaymentHandler(svc)
	reviewHandler := NewReviewHandler(svc)

	group := router.Group("/admin")
	// 所有管理接口均需要认证 + 速率限制
	cfg := config.Load()
	if os.Getenv("APP_ENV") == "production" {
		group.Use(pm.RequireAuth(), mw.RateLimitAdmin())
	} else {
		// 使用配置中的 admin_auth.mode
		switch strings.ToLower(cfg.AdminAuth.Mode) {
		case "jwt":
			group.Use(pm.RequireAuth(), mw.RateLimitAdmin())
		default:
			// 开发模式：保留旧的 AdminAuth（Bearer Token）
			group.Use(mw.AdminAuth(), mw.RateLimitAdmin())
		}
	}
	{
		// 游戏管理 - 使用细粒度权限
		group.GET("/games", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/games"), gameHandler.ListGames)
		group.POST("/games", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/games"), gameHandler.CreateGame)
		group.GET("/games/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/games/:id"), gameHandler.GetGame)
		group.PUT("/games/:id", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/games/:id"), gameHandler.UpdateGame)
		group.DELETE("/games/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/games/:id"), gameHandler.DeleteGame)
		group.GET("/games/:id/logs", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/games/:id/logs"), gameHandler.ListGameLogs)

		// 用户管理 - 使用细粒度权限
		group.GET("/users", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/users"), userHandler.ListUsers)
		group.POST("/users", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/users"), userHandler.CreateUser)
		group.POST("/users/with-player", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/users/with-player"), userHandler.CreateUserWithPlayer)
		group.GET("/users/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/users/:id"), userHandler.GetUser)
		group.PUT("/users/:id", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/users/:id"), userHandler.UpdateUser)
		group.DELETE("/users/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/users/:id"), userHandler.DeleteUser)
		group.PUT("/users/:id/status", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/users/:id/status"), userHandler.UpdateUserStatus)
		group.PUT("/users/:id/role", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/users/:id/role"), userHandler.UpdateUserRole)
		group.GET("/users/:id/orders", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/users/:id/orders"), userHandler.ListUserOrders)
		group.GET("/users/:id/logs", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/users/:id/logs"), userHandler.ListUserLogs)

		// 陪玩师管理 - 使用细粒度权限
		group.GET("/players", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/players"), playerHandler.ListPlayers)
		group.POST("/players", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/players"), playerHandler.CreatePlayer)
		group.GET("/players/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/players/:id"), playerHandler.GetPlayer)
		group.PUT("/players/:id", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/players/:id"), playerHandler.UpdatePlayer)
		group.DELETE("/players/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/players/:id"), playerHandler.DeletePlayer)
		group.PUT("/players/:id/verification", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/players/:id/verification"), playerHandler.UpdatePlayerVerification)
		group.PUT("/players/:id/games", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/players/:id/games"), playerHandler.UpdatePlayerGames)
		group.PUT("/players/:id/skill-tags", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/players/:id/skill-tags"), playerHandler.UpdatePlayerSkillTags)
		group.GET("/players/:id/logs", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/players/:id/logs"), playerHandler.ListPlayerLogs)

		// 订单管理 - 使用细粒度权限
		group.GET("/orders", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/orders"), orderHandler.ListOrders)
		group.POST("/orders", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/orders"), orderHandler.CreateOrder)
		group.GET("/orders/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/orders/:id"), orderHandler.GetOrder)
		group.PUT("/orders/:id", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/orders/:id"), orderHandler.UpdateOrder)
		group.DELETE("/orders/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/orders/:id"), orderHandler.DeleteOrder)
		group.POST("/orders/:id/review", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/orders/:id/review"), orderHandler.ReviewOrder)
		group.POST("/orders/:id/cancel", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/orders/:id/cancel"), orderHandler.CancelOrder)
		group.POST("/orders/:id/assign", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/orders/:id/assign"), orderHandler.AssignOrder)
		group.POST("/orders/:id/confirm", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/orders/:id/confirm"), orderHandler.ConfirmOrder)
		group.POST("/orders/:id/start", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/orders/:id/start"), orderHandler.StartOrder)
		group.POST("/orders/:id/complete", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/orders/:id/complete"), orderHandler.CompleteOrder)
		group.POST("/orders/:id/refund", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/orders/:id/refund"), orderHandler.RefundOrder)
		group.GET("/orders/:id/logs", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/orders/:id/logs"), orderHandler.ListOrderLogs)
		group.GET("/orders/:id/timeline", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/orders/:id/timeline"), orderHandler.GetOrderTimeline)
		group.GET("/orders/:id/payments", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/orders/:id/payments"), orderHandler.ListOrderPayments)
		group.GET("/orders/:id/refunds", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/orders/:id/refunds"), orderHandler.ListOrderRefunds)
		group.GET("/orders/:id/reviews", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/orders/:id/reviews"), orderHandler.ListOrderReviews)

		// 支付管理 - 使用细粒度权限
		group.GET("/payments", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/payments"), paymentHandler.ListPayments)
		group.POST("/payments", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/payments"), paymentHandler.CreatePayment)
		group.GET("/payments/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/payments/:id"), paymentHandler.GetPayment)
		group.PUT("/payments/:id", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/payments/:id"), paymentHandler.UpdatePayment)
		group.DELETE("/payments/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/payments/:id"), paymentHandler.DeletePayment)
		group.POST("/payments/:id/refund", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/payments/:id/refund"), paymentHandler.RefundPayment)
		group.POST("/payments/:id/capture", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/payments/:id/capture"), paymentHandler.CapturePayment)
		group.GET("/payments/:id/logs", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/payments/:id/logs"), paymentHandler.ListPaymentLogs)

		// 评价管理 - 使用细粒度权限
		group.GET("/reviews", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/reviews"), reviewHandler.ListReviews)
		group.POST("/reviews", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/reviews"), reviewHandler.CreateReview)
		group.GET("/reviews/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/reviews/:id"), reviewHandler.GetReview)
		group.PUT("/reviews/:id", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/reviews/:id"), reviewHandler.UpdateReview)
		group.DELETE("/reviews/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/reviews/:id"), reviewHandler.DeleteReview)
		group.GET("/players/:id/reviews", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/players/:id/reviews"), reviewHandler.ListPlayerReviews)
		group.GET("/reviews/:id/logs", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/reviews/:id/logs"), reviewHandler.ListReviewLogs)
	}
}

// RegisterStatsRoutes 注册统计相关路由。
// 使用细粒度权限控制（method+path 级别）。
func RegisterStatsRoutes(router gin.IRouter, stats *service.StatsService, pm *mw.PermissionMiddleware) {
	h := NewStatsHandler(stats)
	group := router.Group("/admin")
	// 统计接口均需要认证 + 速率限制
	cfg := config.Load()
	if os.Getenv("APP_ENV") == "production" {
		group.Use(pm.RequireAuth(), mw.RateLimitAdmin())
	} else {
		// 使用配置中的 admin_auth.mode
		switch strings.ToLower(cfg.AdminAuth.Mode) {
		case "jwt":
			group.Use(pm.RequireAuth(), mw.RateLimitAdmin())
		default:
			group.Use(mw.AdminAuth(), mw.RateLimitAdmin())
		}
	}
	// 统计接口 - 使用细粒度权限
	group.GET("/stats/dashboard", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/stats/dashboard"), h.Dashboard)
	group.GET("/stats/revenue-trend", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/stats/revenue-trend"), h.RevenueTrend)
	group.GET("/stats/user-growth", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/stats/user-growth"), h.UserGrowth)
	group.GET("/stats/orders", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/stats/orders"), h.OrdersSummary)
	group.GET("/stats/top-players", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/stats/top-players"), h.TopPlayers)
	group.GET("/stats/audit/overview", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/stats/audit/overview"), h.AuditOverview)
	group.GET("/stats/audit/trend", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/stats/audit/trend"), h.AuditTrend)
}
