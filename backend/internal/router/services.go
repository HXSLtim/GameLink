package router

import (
	"gamelink/internal/cache"
	chatrepo "gamelink/internal/repository/chat"
	commissionrepo "gamelink/internal/repository/commission"
	feedrepo "gamelink/internal/repository/feed"
	gamerepo "gamelink/internal/repository/game"
	orderrepo "gamelink/internal/repository/implementations"
	notificationrepo "gamelink/internal/repository/notification"
	paymentrepo "gamelink/internal/repository/payment"
	playerrepo "gamelink/internal/repository/player"
	playertagrepo "gamelink/internal/repository/player_tag"
	reviewrepo "gamelink/internal/repository/review"
	reviewreplyrepo "gamelink/internal/repository/reviewreply"
	serviceitemrepo "gamelink/internal/repository/serviceitem"
	userrepo "gamelink/internal/repository/user"
	withdrawrepo "gamelink/internal/repository/withdraw"
	"gamelink/internal/scheduler"
	chatservice "gamelink/internal/service/chat"
	commissionservice "gamelink/internal/service/commission"
	earningsservice "gamelink/internal/service/earnings"
	feedservice "gamelink/internal/service/feed"
	giftservice "gamelink/internal/service/gift"
	itemservice "gamelink/internal/service/item"
	notificationservice "gamelink/internal/service/notification"
	orderservice "gamelink/internal/service/order"
	paymentservice "gamelink/internal/service/payment"
	playerservice "gamelink/internal/service/player"
	reviewservice "gamelink/internal/service/review"

	"gorm.io/gorm"
)

// appServices 包含所有领域服务实例和调度器句柄，供路由注册使用。
type appServices struct {
	commissionSvc       *commissionservice.CommissionService
	serviceItemSvc      *itemservice.ServiceItemService
	giftSvc             *giftservice.GiftService
	orderSvc            *orderservice.OrderService
	paymentSvc          *paymentservice.PaymentService
	playerSvc           *playerservice.PlayerService
	reviewSvc           *reviewservice.ReviewService
	earningsSvc         *earningsservice.EarningsService
	chatSvc             *chatservice.ChatService
	feedSvc             *feedservice.Service
	notificationSvc     *notificationservice.Service
	settlementScheduler *scheduler.SettlementScheduler
	chatRetention       *scheduler.ChatRetentionScheduler
}

// initServices 初始化领域服务和调度任务（但不启动调度器）。
func initServices(orm *gorm.DB, cacheClient cache.Cache) *appServices {
	// 仓库实例（仅在此函数内部复用）
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
	reviewReplyRepo := reviewreplyrepo.NewReviewReplyRepository(orm)
	playerTagRepo := playertagrepo.NewPlayerTagRepository(orm)
	withdrawRepo := withdrawrepo.NewWithdrawRepository(orm)
	commissionRepo := commissionrepo.NewCommissionRepository(orm)
	serviceItemRepo := serviceitemrepo.NewServiceItemRepository(orm)
	feedRepo := feedrepo.NewFeedRepository(orm)
	notificationRepo := notificationrepo.NewNotificationRepository(orm)

	// 领域服务
	commissionSvc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)
	serviceItemSvc := itemservice.NewServiceItemService(serviceItemRepo, gameRepo, playerRepo)
	giftSvc := giftservice.NewGiftService(serviceItemRepo, orderRepo, playerRepo, commissionRepo)
	orderSvc := orderservice.NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)
	// 注入聊天群仓库用于订单聊天自动销毁
	orderSvc.SetChatGroupRepository(chatGroupRepo)
	paymentSvc := paymentservice.NewPaymentService(paymentRepo, orderRepo)
	playerSvc := playerservice.NewPlayerService(playerRepo, userRepo, gameRepo, orderRepo, reviewRepo, playerTagRepo, cacheClient)
	reviewSvc := reviewservice.NewReviewService(reviewRepo, orderRepo, playerRepo, userRepo, reviewReplyRepo)
	earningsSvc := earningsservice.NewEarningsService(playerRepo, orderRepo, withdrawRepo)
	chatSvc := chatservice.NewChatService(chatGroupRepo, chatMemberRepo, chatMessageRepo, chatReportRepo, cacheClient)
	feedSvc := feedservice.NewService(feedRepo, nil)
	notificationSvc := notificationservice.NewService(notificationRepo)

	// 调度器（先构造，调用方负责 Start/Stop）
	settlementScheduler := scheduler.NewSettlementScheduler(commissionSvc)
	chatRetention := scheduler.NewChatRetentionScheduler(chatGroupRepo, chatMessageRepo, 30)

	return &appServices{
		commissionSvc:       commissionSvc,
		serviceItemSvc:      serviceItemSvc,
		giftSvc:             giftSvc,
		orderSvc:            orderSvc,
		paymentSvc:          paymentSvc,
		playerSvc:           playerSvc,
		reviewSvc:           reviewSvc,
		earningsSvc:         earningsSvc,
		chatSvc:             chatSvc,
		feedSvc:             feedSvc,
		notificationSvc:     notificationSvc,
		settlementScheduler: settlementScheduler,
		chatRetention:       chatRetention,
	}
}
