package admin

import (
    "os"

    "github.com/gin-gonic/gin"

    mw "gamelink/internal/handler/middleware"
    "gamelink/internal/service"
)

// RegisterRoutes 注册后台管理相关路由。
func RegisterRoutes(router gin.IRouter, svc *service.AdminService) {
	gameHandler := NewGameHandler(svc)
	userHandler := NewUserHandler(svc)
	playerHandler := NewPlayerHandler(svc)
    orderHandler := NewOrderHandler(svc)
    paymentHandler := NewPaymentHandler(svc)
    reviewHandler := NewReviewHandler(svc)

	group := router.Group("/admin")
	// 生产环境强制 JWT + RBAC；开发环境可由 ADMIN_AUTH_MODE 控制
	if os.Getenv("APP_ENV") == "production" {
		group.Use(mw.JWTAuth(), mw.RequireRole("admin"), mw.RateLimitAdmin())
	} else {
		switch os.Getenv("ADMIN_AUTH_MODE") {
		case "jwt", "JWT":
			group.Use(mw.JWTAuth(), mw.RequireRole("admin"), mw.RateLimitAdmin())
		default:
			group.Use(mw.AdminAuth(), mw.RateLimitAdmin())
		}
	}
	{
		group.GET("/games", gameHandler.ListGames)
		group.POST("/games", gameHandler.CreateGame)
		group.GET("/games/:id", gameHandler.GetGame)
		group.PUT("/games/:id", gameHandler.UpdateGame)
		group.DELETE("/games/:id", gameHandler.DeleteGame)
        group.GET("/games/:id/logs", gameHandler.ListGameLogs)

		group.GET("/users", userHandler.ListUsers)
		group.POST("/users", userHandler.CreateUser)
		group.POST("/users/with-player", userHandler.CreateUserWithPlayer)
		group.GET("/users/:id", userHandler.GetUser)
		group.PUT("/users/:id", userHandler.UpdateUser)
		group.DELETE("/users/:id", userHandler.DeleteUser)
		group.PUT("/users/:id/status", userHandler.UpdateUserStatus)
		group.PUT("/users/:id/role", userHandler.UpdateUserRole)
		group.GET("/users/:id/orders", userHandler.ListUserOrders)
        group.GET("/users/:id/logs", userHandler.ListUserLogs)

		group.GET("/players", playerHandler.ListPlayers)
		group.POST("/players", playerHandler.CreatePlayer)
		group.GET("/players/:id", playerHandler.GetPlayer)
		group.PUT("/players/:id", playerHandler.UpdatePlayer)
		group.DELETE("/players/:id", playerHandler.DeletePlayer)
        group.PUT("/players/:id/verification", playerHandler.UpdatePlayerVerification)
        group.PUT("/players/:id/games", playerHandler.UpdatePlayerGames)
        group.PUT("/players/:id/skill-tags", playerHandler.UpdatePlayerSkillTags)
        group.GET("/players/:id/logs", playerHandler.ListPlayerLogs)

        group.GET("/orders", orderHandler.ListOrders)
        group.POST("/orders", orderHandler.CreateOrder)
        group.GET("/orders/:id", orderHandler.GetOrder)
        group.PUT("/orders/:id", orderHandler.UpdateOrder)
        group.DELETE("/orders/:id", orderHandler.DeleteOrder)
        group.POST("/orders/:id/review", orderHandler.ReviewOrder)
        group.POST("/orders/:id/cancel", orderHandler.CancelOrder)
        group.POST("/orders/:id/assign", orderHandler.AssignOrder)
        group.GET("/orders/:id/logs", orderHandler.ListOrderLogs)

        group.GET("/payments", paymentHandler.ListPayments)
        group.POST("/payments", paymentHandler.CreatePayment)
        group.GET("/payments/:id", paymentHandler.GetPayment)
        group.PUT("/payments/:id", paymentHandler.UpdatePayment)
        group.DELETE("/payments/:id", paymentHandler.DeletePayment)
        group.POST("/payments/:id/refund", paymentHandler.RefundPayment)
        group.POST("/payments/:id/capture", paymentHandler.CapturePayment)
        group.GET("/payments/:id/logs", paymentHandler.ListPaymentLogs)

        // Reviews
        group.GET("/reviews", reviewHandler.ListReviews)
        group.POST("/reviews", reviewHandler.CreateReview)
        group.GET("/reviews/:id", reviewHandler.GetReview)
        group.PUT("/reviews/:id", reviewHandler.UpdateReview)
        group.DELETE("/reviews/:id", reviewHandler.DeleteReview)
        group.GET("/players/:id/reviews", reviewHandler.ListPlayerReviews)
        group.GET("/reviews/:id/logs", reviewHandler.ListReviewLogs)
	}
}

// RegisterStatsRoutes 注册统计相关路由。
func RegisterStatsRoutes(router gin.IRouter, stats *service.StatsService) {
    h := NewStatsHandler(stats)
    group := router.Group("/admin")
    // 统计接口与 admin 采用同样的鉴权策略（简化处理，开发环境回退 token 验证）
    if os.Getenv("APP_ENV") == "production" {
        group.Use(mw.JWTAuth(), mw.RequireRole("admin"), mw.RateLimitAdmin())
    } else {
        switch os.Getenv("ADMIN_AUTH_MODE") {
        case "jwt", "JWT":
            group.Use(mw.JWTAuth(), mw.RequireRole("admin"), mw.RateLimitAdmin())
        default:
            group.Use(mw.AdminAuth(), mw.RateLimitAdmin())
        }
    }
    group.GET("/stats/dashboard", h.Dashboard)
    group.GET("/stats/revenue-trend", h.RevenueTrend)
    group.GET("/stats/user-growth", h.UserGrowth)
    group.GET("/stats/orders", h.OrdersSummary)
    group.GET("/stats/top-players", h.TopPlayers)
    group.GET("/stats/audit/overview", h.AuditOverview)
    group.GET("/stats/audit/trend", h.AuditTrend)
}
