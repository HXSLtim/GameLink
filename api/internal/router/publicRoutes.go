package router

import (
	"github.com/gin-gonic/gin"

	publichandler "gamelink/internal/handler/public"
	userhandler "gamelink/internal/handler/user"
	"gamelink/internal/repository"
	chatrepo "gamelink/internal/repository/chat"
	gamerepo "gamelink/internal/repository/game"
	orderrepo "gamelink/internal/repository/order"
	playerrepo "gamelink/internal/repository/player"
	referralrepo "gamelink/internal/repository/referral"
	serviceitemrepo "gamelink/internal/repository/serviceitem"
	userrepo "gamelink/internal/repository/user"
	authservice "gamelink/internal/service/auth"
	"gamelink/internal/service/external"
	referralservice "gamelink/internal/service/referral"
	"gamelink/internal/service/sms"
	"gamelink/internal/service/verification"
	"gamelink/pkg/cache"
	"gamelink/pkg/config"

	"gorm.io/gorm"
)

// registerPublicRoutes 注册公共 API 路由（无需认证）
func registerPublicRoutes(api *gin.RouterGroup, orm *gorm.DB, cacheClient cache.Cache) {
	publicGroup := api.Group("/public")

	// 初始化仓库
	userRepo := userrepo.NewUserRepository(orm)
	playerRepo := playerrepo.NewPlayerRepository(orm)
	gameRepo := gamerepo.NewGameRepository(orm)
	serviceItemRepo := serviceitemrepo.NewServiceItemRepository(orm)
	chatGroupRepo := chatrepo.NewChatGroupRepository(orm)
	reviewRepo := orderrepo.NewReviewRepository(orm)
	referralRepo := referralrepo.NewReferralRepository(orm)

	// 初始化服务
	wechatSvc := authservice.NewWeChatAuthService(userRepo, playerRepo)
	wechatSvc.SetReferralService(referralservice.NewReferralService(referralRepo))
	wechatSvc.SetReferralTrigger(referralservice.NewTriggerService(orm))

	// 注册认证路由
	authHandler := publichandler.NewAuthHandler(wechatSvc)
	authHandler.RegisterRoutes(publicGroup)

	// 兼容路径：/api/v1/auth/wechat/*
	authGroup := api.Group("/auth")
	authGroup.POST("/wechat/login", authHandler.WeChatLogin)
	authGroup.POST("/wechat/refresh", authHandler.RefreshToken)

	// 注册陪玩师列表路由
	playerHandler := publichandler.NewPlayerHandler(playerRepo, userRepo)
	playerHandler.RegisterRoutes(publicGroup)

	// 注册陪玩师评价路由
	reviewHandler := publichandler.NewPublicReviewHandler(reviewRepo, userRepo)
	reviewHandler.RegisterRoutes(publicGroup)

	// 注册游戏列表路由
	gameHandler := publichandler.NewGameHandler(gameRepo)
	gameHandler.RegisterRoutes(publicGroup)

	// 注册服务项目列表路由
	serviceItemHandler := publichandler.NewServiceItemHandler(serviceItemRepo)
	serviceItemHandler.RegisterRoutes(publicGroup)

	// 注册公共频道路由
	publicChatHandler := publichandler.NewPublicChatHandler(chatGroupRepo)
	publicChatHandler.RegisterRoutes(publicGroup)

	// 注册搜索路由
	searchHandler := publichandler.NewSearchHandler(playerRepo, gameRepo, userRepo)
	publichandler.RegisterSearchRoutes(publicGroup, searchHandler)

	// 注册验证码路由
	appCfg := config.Load()
	externalCfg := external.NewConfig(appCfg.ExternalAPI)
	smsSvc := sms.NewService(externalCfg)
	verificationSvc := verification.NewService(cacheClient, smsSvc)
	verificationHandler := publichandler.NewVerificationHandler(verificationSvc)
	verificationHandler.RegisterRoutes(publicGroup)

	// 注册 Banner 路由
	bannerHandler := publichandler.NewBannerHandler(orm)
	bannerHandler.RegisterRoutes(publicGroup)
}

// registerRoleSwitchRoutes 注册角色切换路由（需要认证）
func registerRoleSwitchRoutes(userGroup *gin.RouterGroup, userRepo repository.UserRepository, playerRepo repository.PlayerRepository) {
	roleSvc := authservice.NewRoleService(userRepo, playerRepo)
	roleHandler := userhandler.NewRoleHandler(roleSvc)
	roleHandler.RegisterRoutes(userGroup)
}
