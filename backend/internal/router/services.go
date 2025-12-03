package router

import (
	"gamelink/pkg/cache"
	chatrepo "gamelink/internal/repository/chat"
	commissionrepo "gamelink/internal/repository/commission"
	feedrepo "gamelink/internal/repository/content"
	gamerepo "gamelink/internal/repository/game"
	orderrepo "gamelink/internal/repository/implementations"
	notificationrepo "gamelink/internal/repository/content"
	paymentrepo "gamelink/internal/repository/order"
	userrepo "gamelink/internal/repository/user"
	reviewrepo "gamelink/internal/repository/order"
	reviewreplyrepo "gamelink/internal/repository/order"
	serviceitemrepo "gamelink/internal/repository/serviceitem"
	adminrepo "gamelink/internal/repository/admin"
	withdrawrepo "gamelink/internal/repository/withdraw"
	"gamelink/pkg/scheduler"
	chatservice "gamelink/internal/service/chat"
	commissionservice "gamelink/internal/service/commission"
	contentservice "gamelink/internal/service/content"
	giftservice "gamelink/internal/service/gift"
	itemservice "gamelink/internal/service/item"
	orderservice "gamelink/internal/service/order"
	serviceplayer "gamelink/internal/service/player"
	userservice "gamelink/internal/service/user"
	walletservice "gamelink/internal/service/wallet"

	"gorm.io/gorm"
)

// appServices 包含所有领域服务实例和调度器句柄，供路由注册使用。
type appServices struct {
	commissionSvc       *commissionservice.CommissionService
	serviceItemSvc      *itemservice.ServiceItemService
	giftSvc             *giftservice.GiftService
	orderSvc            *orderservice.OrderService
	paymentSvc          *orderservice.PaymentService
	playerSvc           *serviceplayer.PlayerService
	reviewSvc           *orderservice.ReviewService
	disputeSvc          *orderservice.DisputeService
	earningsSvc         *userservice.EarningsService
	chatSvc             *chatservice.ChatService
	feedSvc             *contentservice.FeedService
	notificationSvc     *contentservice.NotificationService
	walletSvc           *walletservice.WalletService
	settlementScheduler *scheduler.SettlementScheduler
	chatRetention       *scheduler.ChatRetentionScheduler
}

// initServices 初始化领域服务和调度任务（但不启动调度器）。
func initServices(orm *gorm.DB, cacheClient cache.Cache) *appServices {
	// 仓库实例（仅在此函数内部复用）
	userRepo := userrepo.NewUserRepository(orm)
	playerRepo := userrepo.NewPlayerRepository(orm)
	gameRepo := gamerepo.NewGameRepository(orm)
	orderRepo := orderrepo.NewOrderRepository(orm)
	chatGroupRepo := chatrepo.NewChatGroupRepository(orm)
	chatMemberRepo := chatrepo.NewChatMemberRepository(orm)
	chatMessageRepo := chatrepo.NewChatMessageRepository(orm)
	chatReportRepo := chatrepo.NewChatReportRepository(orm)
	paymentRepo := paymentrepo.NewPaymentRepository(orm)
	reviewRepo := reviewrepo.NewReviewRepository(orm)
	reviewReplyRepo := reviewreplyrepo.NewReviewReplyRepository(orm)
	playerTagRepo := userrepo.NewPlayerTagRepository(orm)
	withdrawRepo := withdrawrepo.NewWithdrawRepository(orm)
	commissionRepo := commissionrepo.NewCommissionRepository(orm)
	serviceItemRepo := serviceitemrepo.NewServiceItemRepository(orm)
	feedRepo := feedrepo.NewFeedRepository(orm)
	notificationRepo := notificationrepo.NewNotificationRepository(orm)
	walletRepo := userrepo.NewWalletRepository(orm)

	// 领域服务
	commissionSvc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)
	serviceItemSvc := itemservice.NewServiceItemService(serviceItemRepo, gameRepo, playerRepo)
	giftSvc := giftservice.NewGiftService(serviceItemRepo, orderRepo, playerRepo, commissionRepo)
	orderSvc := orderservice.NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)
	// 注入聊天群仓库用于订单聊天自动销毁
	orderSvc.SetChatGroupRepository(chatGroupRepo)
	paymentSvc := orderservice.NewPaymentService(paymentRepo, orderRepo)
	playerSvc := serviceplayer.NewPlayerService(playerRepo, userRepo, gameRepo, orderRepo, reviewRepo, playerTagRepo, cacheClient)
	reviewSvc := orderservice.NewReviewService(reviewRepo, orderRepo, playerRepo, userRepo, reviewReplyRepo)
	disputeRepo := reviewrepo.NewDisputeRepository(orm)
	operationLogRepo := adminrepo.NewOperationLogRepository(orm)
	disputeSvc := orderservice.NewDisputeService(disputeRepo, orderRepo, userRepo, operationLogRepo, notificationRepo, paymentRepo)
	earningsSvc := userservice.NewEarningsService(playerRepo, orderRepo, withdrawRepo)
	chatSvc := chatservice.NewChatService(chatGroupRepo, chatMemberRepo, chatMessageRepo, chatReportRepo, cacheClient)
	feedSvc := contentservice.NewFeedService(feedRepo, nil)
	notificationSvc := contentservice.NewNotificationService(notificationRepo)
	walletSvc := walletservice.NewWalletService(walletRepo, paymentRepo, orderRepo)

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
		disputeSvc:          disputeSvc,
		earningsSvc:         earningsSvc,
		chatSvc:             chatSvc,
		feedSvc:             feedSvc,
		notificationSvc:     notificationSvc,
		walletSvc:           walletSvc,
		settlementScheduler: settlementScheduler,
		chatRetention:       chatRetention,
	}
}
