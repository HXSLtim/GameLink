package router

import (
	"github.com/gin-gonic/gin"

	publichandler "gamelink/internal/handler/public"
	userhandler "gamelink/internal/handler/user"
	"gamelink/internal/repository"
	gamerepo "gamelink/internal/repository/game"
	playerrepo "gamelink/internal/repository/player"
	serviceitemrepo "gamelink/internal/repository/serviceitem"
	userrepo "gamelink/internal/repository/user"
	authservice "gamelink/internal/service/auth"
	"gorm.io/gorm"
)

// registerPublicRoutes 注册公共 API 路由（无需认证）
func registerPublicRoutes(api *gin.RouterGroup, orm *gorm.DB) {
	publicGroup := api.Group("/public")

	// 初始化仓库
	userRepo := userrepo.NewUserRepository(orm)
	playerRepo := playerrepo.NewPlayerRepository(orm)
	gameRepo := gamerepo.NewGameRepository(orm)
	serviceItemRepo := serviceitemrepo.NewServiceItemRepository(orm)

	// 初始化服务
	wechatSvc := authservice.NewWeChatAuthService(userRepo, playerRepo)

	// 注册认证路由
	authHandler := publichandler.NewAuthHandler(wechatSvc)
	authHandler.RegisterRoutes(publicGroup)

	// 注册陪玩师列表路由
	playerHandler := publichandler.NewPlayerHandler(playerRepo, userRepo)
	playerHandler.RegisterRoutes(publicGroup)

	// 注册游戏列表路由
	gameHandler := publichandler.NewGameHandler(gameRepo)
	gameHandler.RegisterRoutes(publicGroup)

	// 注册服务项目列表路由
	serviceItemHandler := publichandler.NewServiceItemHandler(serviceItemRepo)
	serviceItemHandler.RegisterRoutes(publicGroup)

	// 注册搜索路由
	searchHandler := publichandler.NewSearchHandler(playerRepo, gameRepo, userRepo)
	publichandler.RegisterSearchRoutes(publicGroup, searchHandler)
}

// registerRoleSwitchRoutes 注册角色切换路由（需要认证）
func registerRoleSwitchRoutes(userGroup *gin.RouterGroup, userRepo repository.UserRepository, playerRepo repository.PlayerRepository) {
	roleSvc := authservice.NewRoleService(userRepo, playerRepo)
	roleHandler := userhandler.NewRoleHandler(roleSvc)
	roleHandler.RegisterRoutes(userGroup)
}
