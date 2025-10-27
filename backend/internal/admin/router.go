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

	group := router.Group("/admin")
	// Protect admin endpoints: prefer JWT + role guard when configured, else fallback to AdminAuth
	switch os.Getenv("ADMIN_AUTH_MODE") {
	case "jwt", "JWT":
		group.Use(mw.JWTAuth(), mw.RequireRole("admin"), mw.RateLimitAdmin())
	default:
		group.Use(mw.AdminAuth(), mw.RateLimitAdmin())
	}
	{
		group.GET("/games", gameHandler.ListGames)
		group.POST("/games", gameHandler.CreateGame)
		group.GET("/games/:id", gameHandler.GetGame)
		group.PUT("/games/:id", gameHandler.UpdateGame)
		group.DELETE("/games/:id", gameHandler.DeleteGame)

		group.GET("/users", userHandler.ListUsers)
		group.POST("/users", userHandler.CreateUser)
		group.GET("/users/:id", userHandler.GetUser)
		group.PUT("/users/:id", userHandler.UpdateUser)
		group.DELETE("/users/:id", userHandler.DeleteUser)

		group.GET("/players", playerHandler.ListPlayers)
		group.POST("/players", playerHandler.CreatePlayer)
		group.GET("/players/:id", playerHandler.GetPlayer)
		group.PUT("/players/:id", playerHandler.UpdatePlayer)
		group.DELETE("/players/:id", playerHandler.DeletePlayer)

		group.GET("/orders", orderHandler.ListOrders)
		group.GET("/orders/:id", orderHandler.GetOrder)
		group.PUT("/orders/:id", orderHandler.UpdateOrder)
		group.DELETE("/orders/:id", orderHandler.DeleteOrder)

		group.GET("/payments", paymentHandler.ListPayments)
		group.GET("/payments/:id", paymentHandler.GetPayment)
		group.PUT("/payments/:id", paymentHandler.UpdatePayment)
		group.DELETE("/payments/:id", paymentHandler.DeletePayment)
	}
}
