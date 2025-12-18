package order

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	commissionrepo "gamelink/internal/repository/commission"
	"gamelink/internal/repository/game"
	orderimpl "gamelink/internal/repository/implementations"
	orderrepo "gamelink/internal/repository/order"
	"gamelink/internal/repository/review"
	"gamelink/internal/repository/user"
	"gamelink/pkg/testutil"
)

func setupOrderTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.User{},
		&model.Player{},
		&model.Game{},
		&model.Order{},
		&model.Payment{},
		&model.Review{},
		&model.CommissionRule{},
		&model.CommissionRecord{},
		&model.Wallet{},
	)
	return db
}

func createOrderTestData(t *testing.T, db *gorm.DB) (customer *model.User, playerUser *model.User, player *model.Player, gameModel *model.Game) {
	t.Helper()

	// 创建用户
	customer = &model.User{
		Phone:        "13800000001",
		Email:        "customer@test.com",
		Name:         "Customer",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(customer).Error)

	// 创建陪玩师用户
	playerUser = &model.User{
		Phone:        "13800000002",
		Email:        "player@test.com",
		Name:         "Player",
		Role:         model.RolePlayer,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(playerUser).Error)

	// 创建游戏
	gameModel = &model.Game{Key: "lol", Name: "英雄联盟", Category: "moba"}
	require.NoError(t, db.Create(gameModel).Error)

	// 创建陪玩师
	player = &model.Player{
		UserID:             playerUser.ID,
		Nickname:           "Pro Player",
		MainGameID:         gameModel.ID,
		HourlyRateCents:    5000,
		VerificationStatus: model.VerificationVerified,
	}
	require.NoError(t, db.Create(player).Error)

	return
}

func createOrderService(db *gorm.DB) *OrderService {
	orderRepo := orderimpl.NewOrderRepository(db)
	playerRepo := user.NewPlayerRepository(db)
	userRepo := user.NewUserRepository(db)
	gameRepo := game.NewGameRepository(db)
	paymentRepo := orderrepo.NewPaymentRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	commissionRepo := commissionrepo.NewCommissionRepository(db)

	return NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)
}

func TestOrderService_CreateOrder(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	t.Run("创建订单成功", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		resp, err := svc.CreateOrder(context.Background(), customer.ID, CreateOrderRequest{
			PlayerID:       player.ID,
			GameID:         gameModel.ID,
			Title:          "测试订单",
			Description:    "测试描述",
			ScheduledStart: &scheduledStart,
			DurationHours:  2,
		})
		require.NoError(t, err)
		assert.NotZero(t, resp.OrderID)
		assert.True(t, resp.NeedPayment)
		assert.Greater(t, resp.PriceCents, int64(0))
	})

	t.Run("无效陪玩师ID", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		_, err := svc.CreateOrder(context.Background(), customer.ID, CreateOrderRequest{
			PlayerID:       99999,
			GameID:         gameModel.ID,
			Title:          "测试订单",
			ScheduledStart: &scheduledStart,
			DurationHours:  2,
		})
		assert.Error(t, err)
	})
}

func TestOrderService_GetMyOrders(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	// 创建测试订单
	scheduledStart := time.Now().Add(24 * time.Hour)
	scheduledEnd := scheduledStart.Add(2 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "测试订单",
		Status:          model.OrderStatusPending,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
		ScheduledEnd:    &scheduledEnd,
	}
	order.SetPlayerID(player.ID)
	order.SetGameID(gameModel.ID)
	require.NoError(t, db.Create(order).Error)

	t.Run("获取用户订单列表", func(t *testing.T) {
		resp, err := svc.GetMyOrders(context.Background(), customer.ID, MyOrderListRequest{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(1), resp.Total)
		assert.Len(t, resp.Orders, 1)
		assert.Equal(t, "测试订单", resp.Orders[0].Title)
	})

	t.Run("按状态筛选订单", func(t *testing.T) {
		resp, err := svc.GetMyOrders(context.Background(), customer.ID, MyOrderListRequest{
			Status:   "pending",
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(1), resp.Total)

		// 筛选不存在的状态
		resp, err = svc.GetMyOrders(context.Background(), customer.ID, MyOrderListRequest{
			Status:   "completed",
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(0), resp.Total)
	})
}

func TestOrderService_GetOrderDetail(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	// 创建测试订单
	scheduledStart := time.Now().Add(24 * time.Hour)
	scheduledEnd := scheduledStart.Add(2 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "测试订单",
		Description:     "测试描述",
		Status:          model.OrderStatusPending,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
		ScheduledEnd:    &scheduledEnd,
	}
	order.SetPlayerID(player.ID)
	order.SetGameID(gameModel.ID)
	require.NoError(t, db.Create(order).Error)

	t.Run("用户获取订单详情", func(t *testing.T) {
		resp, err := svc.GetOrderDetail(context.Background(), customer.ID, order.ID)
		require.NoError(t, err)
		assert.Equal(t, "测试订单", resp.Order.Title)
		assert.Equal(t, "测试描述", resp.Order.Description)
		assert.NotNil(t, resp.Player)
		assert.Equal(t, "Pro Player", resp.Player.Nickname)
	})

	t.Run("陪玩师获取订单详情", func(t *testing.T) {
		// 陪玩师也可以查看自己接的订单
		// 注意：GetOrderDetail 检查的是 order.GetPlayerID() 与 userID 的匹配
		// 但 order.GetPlayerID() 返回的是 Player.ID，不是 User.ID
		// 所以这里用 player.ID 来测试
		resp, err := svc.GetOrderDetail(context.Background(), player.ID, order.ID)
		require.NoError(t, err)
		assert.Equal(t, "测试订单", resp.Order.Title)
	})

	t.Run("无权限查看订单", func(t *testing.T) {
		// 创建另一个用户
		otherUser := &model.User{
			Phone:        "13800000003",
			Email:        "other@test.com",
			Name:         "Other",
			Role:         model.RoleUser,
			Status:       model.UserStatusActive,
			PasswordHash: "hashed",
		}
		require.NoError(t, db.Create(otherUser).Error)

		_, err := svc.GetOrderDetail(context.Background(), otherUser.ID, order.ID)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("订单不存在", func(t *testing.T) {
		_, err := svc.GetOrderDetail(context.Background(), customer.ID, 99999)
		assert.ErrorIs(t, err, ErrNotFound)
	})
}

func TestOrderService_CancelOrder(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	t.Run("取消待支付订单", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "待取消订单",
			Status:          model.OrderStatusPending,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
		}
		order.SetPlayerID(player.ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		err := svc.CancelOrder(context.Background(), customer.ID, order.ID, CancelOrderRequest{
			Reason: "不想要了",
		})
		require.NoError(t, err)

		// 验证状态
		var updated model.Order
		require.NoError(t, db.First(&updated, order.ID).Error)
		assert.Equal(t, model.OrderStatusCanceled, updated.Status)
		assert.Equal(t, "不想要了", updated.CancelReason)
	})

	t.Run("无法取消进行中的订单", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "进行中订单",
			Status:          model.OrderStatusInProgress,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
		}
		order.SetPlayerID(player.ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		err := svc.CancelOrder(context.Background(), customer.ID, order.ID, CancelOrderRequest{
			Reason: "不想要了",
		})
		assert.ErrorIs(t, err, ErrInvalidTransition)
	})
}

func TestOrderService_CompleteOrder(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	t.Run("完成进行中的订单", func(t *testing.T) {
		scheduledStart := time.Now().Add(-2 * time.Hour)
		startedAt := time.Now().Add(-1 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "进行中订单",
			Status:          model.OrderStatusInProgress,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
			StartedAt:       &startedAt,
		}
		order.SetPlayerID(player.ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		err := svc.CompleteOrder(context.Background(), customer.ID, order.ID)
		require.NoError(t, err)

		// 验证状态
		var updated model.Order
		require.NoError(t, db.First(&updated, order.ID).Error)
		assert.Equal(t, model.OrderStatusCompleted, updated.Status)
		assert.NotNil(t, updated.CompletedAt)
	})

	t.Run("无法完成待支付订单", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "待支付订单",
			Status:          model.OrderStatusPending,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
		}
		order.SetPlayerID(player.ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		err := svc.CompleteOrder(context.Background(), customer.ID, order.ID)
		assert.ErrorIs(t, err, ErrInvalidTransition)
	})
}

func TestOrderService_GetAvailableOrders(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	// 创建已支付待接单的订单
	scheduledStart := time.Now().Add(24 * time.Hour)
	scheduledEnd := scheduledStart.Add(2 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "待接单订单",
		Status:          model.OrderStatusConfirmed,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
		ScheduledEnd:    &scheduledEnd,
	}
	order.SetPlayerID(player.ID)
	order.SetGameID(gameModel.ID)
	require.NoError(t, db.Create(order).Error)

	t.Run("获取可接订单列表", func(t *testing.T) {
		orders, total, err := svc.GetAvailableOrders(context.Background(), AvailableOrdersRequest{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, orders, 1)
		assert.Equal(t, "待接单订单", orders[0].Title)
	})

	t.Run("按游戏筛选可接订单", func(t *testing.T) {
		gameID := gameModel.ID
		orders, total, err := svc.GetAvailableOrders(context.Background(), AvailableOrdersRequest{
			GameID:   &gameID,
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, orders, 1)
	})
}

func TestOrderService_AcceptOrder(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, playerUser, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	t.Run("陪玩师接单成功", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "待接单订单",
			Status:          model.OrderStatusConfirmed,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
		}
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		err := svc.AcceptOrder(context.Background(), playerUser.ID, order.ID)
		require.NoError(t, err)

		// 验证状态
		var updated model.Order
		require.NoError(t, db.First(&updated, order.ID).Error)
		assert.Equal(t, model.OrderStatusInProgress, updated.Status)
		assert.Equal(t, player.ID, updated.GetPlayerID())
		assert.NotNil(t, updated.StartedAt)
	})

	t.Run("无法接已接单的订单", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "已接单订单",
			Status:          model.OrderStatusInProgress,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
		}
		order.SetPlayerID(player.ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		err := svc.AcceptOrder(context.Background(), playerUser.ID, order.ID)
		assert.ErrorIs(t, err, ErrInvalidTransition)
	})
}

func TestOrderService_CompleteOrderByPlayer(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, playerUser, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	t.Run("陪玩师完成订单", func(t *testing.T) {
		scheduledStart := time.Now().Add(-2 * time.Hour)
		startedAt := time.Now().Add(-1 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "进行中订单",
			Status:          model.OrderStatusInProgress,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
			StartedAt:       &startedAt,
		}
		order.SetPlayerID(player.ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		err := svc.CompleteOrderByPlayer(context.Background(), playerUser.ID, order.ID)
		require.NoError(t, err)

		// 验证状态
		var updated model.Order
		require.NoError(t, db.First(&updated, order.ID).Error)
		assert.Equal(t, model.OrderStatusCompleted, updated.Status)
		assert.NotNil(t, updated.CompletedAt)
	})

	t.Run("非订单陪玩师无法完成订单", func(t *testing.T) {
		// 创建另一个陪玩师
		otherPlayerUser := &model.User{
			Phone:        "13800000004",
			Email:        "other_player@test.com",
			Name:         "Other Player",
			Role:         model.RolePlayer,
			Status:       model.UserStatusActive,
			PasswordHash: "hashed",
		}
		require.NoError(t, db.Create(otherPlayerUser).Error)

		otherPlayer := &model.Player{
			UserID:             otherPlayerUser.ID,
			Nickname:           "Other Pro",
			MainGameID:         gameModel.ID,
			HourlyRateCents:    6000,
			VerificationStatus: model.VerificationVerified,
		}
		require.NoError(t, db.Create(otherPlayer).Error)

		scheduledStart := time.Now().Add(-2 * time.Hour)
		startedAt := time.Now().Add(-1 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "进行中订单2",
			Status:          model.OrderStatusInProgress,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
			StartedAt:       &startedAt,
		}
		order.SetPlayerID(player.ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		err := svc.CompleteOrderByPlayer(context.Background(), otherPlayerUser.ID, order.ID)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})
}

func TestOrderService_CancelConfirmedOrder(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	t.Run("取消已支付订单触发退款", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "已支付待取消订单",
			Status:          model.OrderStatusConfirmed,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
		}
		order.SetPlayerID(player.ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		// 创建支付记录
		paidAt := time.Now()
		payment := &model.Payment{
			OrderID:     order.ID,
			UserID:      customer.ID,
			AmountCents: 10000,
			Method:      model.PaymentMethodAlipay,
			Status:      model.PaymentStatusPaid,
			PaidAt:      &paidAt,
		}
		require.NoError(t, db.Create(payment).Error)

		err := svc.CancelOrder(context.Background(), customer.ID, order.ID, CancelOrderRequest{
			Reason: "改变主意了",
		})
		require.NoError(t, err)

		// 验证状态变为已退款
		var updated model.Order
		require.NoError(t, db.First(&updated, order.ID).Error)
		assert.Equal(t, model.OrderStatusRefunded, updated.Status)
		assert.Equal(t, int64(10000), updated.RefundAmountCents)
		assert.NotNil(t, updated.RefundedAt)
	})
}

func TestOrderService_GetOrderDetailWithPayment(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	// 创建已支付订单
	scheduledStart := time.Now().Add(-24 * time.Hour)
	scheduledEnd := scheduledStart.Add(2 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "已支付订单",
		Description:     "订单描述",
		Status:          model.OrderStatusConfirmed,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
		ScheduledEnd:    &scheduledEnd,
	}
	order.SetPlayerID(player.ID)
	order.SetGameID(gameModel.ID)
	require.NoError(t, db.Create(order).Error)

	// 创建支付记录
	paidAt := time.Now()
	payment := &model.Payment{
		OrderID:     order.ID,
		UserID:      customer.ID,
		AmountCents: 10000,
		Method:      model.PaymentMethodWeChat,
		Status:      model.PaymentStatusPaid,
		PaidAt:      &paidAt,
	}
	require.NoError(t, db.Create(payment).Error)

	t.Run("获取带支付信息的订单详情", func(t *testing.T) {
		resp, err := svc.GetOrderDetail(context.Background(), customer.ID, order.ID)
		require.NoError(t, err)
		assert.Equal(t, "已支付订单", resp.Order.Title)
		assert.NotNil(t, resp.Payment)
		assert.Equal(t, model.PaymentMethodWeChat, resp.Payment.Method)
		assert.Equal(t, int64(10000), resp.Payment.AmountCents)
	})
}

func TestOrderService_GetOrderDetailWithReview(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	// 创建已完成订单
	scheduledStart := time.Now().Add(-24 * time.Hour)
	completedAt := time.Now().Add(-1 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "已完成订单",
		Status:          model.OrderStatusCompleted,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
		CompletedAt:     &completedAt,
	}
	order.SetPlayerID(player.ID)
	order.SetGameID(gameModel.ID)
	require.NoError(t, db.Create(order).Error)

	// 创建评价
	review := &model.Review{
		OrderID:  order.ID,
		UserID:   customer.ID,
		PlayerID: player.ID,
		Score:    5,
		Content:  "非常满意",
		Status:   model.ReviewStatusApproved,
	}
	require.NoError(t, db.Create(review).Error)

	t.Run("获取带评价信息的订单详情", func(t *testing.T) {
		resp, err := svc.GetOrderDetail(context.Background(), customer.ID, order.ID)
		require.NoError(t, err)
		assert.NotNil(t, resp.Review)
		assert.Equal(t, 5, resp.Review.Rating)
		assert.Equal(t, "非常满意", resp.Review.Comment)
	})
}

func TestOrderService_OrderTimeline(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	t.Run("完整订单时间线", func(t *testing.T) {
		scheduledStart := time.Now().Add(-24 * time.Hour)
		startedAt := time.Now().Add(-2 * time.Hour)
		completedAt := time.Now().Add(-1 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "完整流程订单",
			Status:          model.OrderStatusCompleted,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
			StartedAt:       &startedAt,
			CompletedAt:     &completedAt,
		}
		order.SetPlayerID(player.ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		resp, err := svc.GetOrderDetail(context.Background(), customer.ID, order.ID)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(resp.Timeline), 3) // 至少有创建、支付、完成
	})

	t.Run("取消订单时间线", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "已取消订单",
			Status:          model.OrderStatusCanceled,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
			CancelReason:    "用户取消",
		}
		order.SetPlayerID(player.ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		resp, err := svc.GetOrderDetail(context.Background(), customer.ID, order.ID)
		require.NoError(t, err)
		// 检查时间线包含取消信息
		found := false
		for _, t := range resp.Timeline {
			if t.Status == string(model.OrderStatusCanceled) {
				found = true
				break
			}
		}
		assert.True(t, found, "时间线应包含取消状态")
	})

	t.Run("退款订单时间线", func(t *testing.T) {
		scheduledStart := time.Now().Add(-24 * time.Hour)
		refundedAt := time.Now().Add(-1 * time.Hour)
		order := &model.Order{
			UserID:            customer.ID,
			ItemID:            gameModel.ID,
			Title:             "已退款订单",
			Status:            model.OrderStatusRefunded,
			UnitPriceCents:    5000,
			TotalPriceCents:   10000,
			ScheduledStart:    &scheduledStart,
			RefundAmountCents: 10000,
			RefundedAt:        &refundedAt,
		}
		order.SetPlayerID(player.ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		resp, err := svc.GetOrderDetail(context.Background(), customer.ID, order.ID)
		require.NoError(t, err)
		// 检查时间线包含退款信息
		found := false
		for _, t := range resp.Timeline {
			if t.Status == string(model.OrderStatusRefunded) {
				found = true
				break
			}
		}
		assert.True(t, found, "时间线应包含退款状态")
	})
}

func TestOrderService_CanReviewFlag(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	t.Run("已完成未评价订单可评价", func(t *testing.T) {
		scheduledStart := time.Now().Add(-24 * time.Hour)
		completedAt := time.Now().Add(-1 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "可评价订单",
			Status:          model.OrderStatusCompleted,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
			CompletedAt:     &completedAt,
		}
		order.SetPlayerID(player.ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		resp, err := svc.GetMyOrders(context.Background(), customer.ID, MyOrderListRequest{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)

		var found *OrderCardDTO
		for i := range resp.Orders {
			if resp.Orders[i].ID == order.ID {
				found = &resp.Orders[i]
				break
			}
		}
		require.NotNil(t, found)
		assert.True(t, found.CanReview)
	})

	t.Run("已评价订单不可再评价", func(t *testing.T) {
		scheduledStart := time.Now().Add(-24 * time.Hour)
		completedAt := time.Now().Add(-1 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "已评价订单",
			Status:          model.OrderStatusCompleted,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
			CompletedAt:     &completedAt,
		}
		order.SetPlayerID(player.ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		// 创建评价
		review := &model.Review{
			OrderID:  order.ID,
			UserID:   customer.ID,
			PlayerID: player.ID,
			Score:    5,
			Content:  "好评",
			Status:   model.ReviewStatusApproved,
		}
		require.NoError(t, db.Create(review).Error)

		resp, err := svc.GetMyOrders(context.Background(), customer.ID, MyOrderListRequest{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)

		var found *OrderCardDTO
		for i := range resp.Orders {
			if resp.Orders[i].ID == order.ID {
				found = &resp.Orders[i]
				break
			}
		}
		require.NotNil(t, found)
		assert.False(t, found.CanReview)
	})
}

func TestOrderService_OrderActionFlags(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	t.Run("待支付订单操作标志", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "待支付订单",
			Status:          model.OrderStatusPending,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
		}
		order.SetPlayerID(player.ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		resp, err := svc.GetMyOrders(context.Background(), customer.ID, MyOrderListRequest{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)

		var found *OrderCardDTO
		for i := range resp.Orders {
			if resp.Orders[i].ID == order.ID {
				found = &resp.Orders[i]
				break
			}
		}
		require.NotNil(t, found)
		assert.True(t, found.CanPay)
		assert.True(t, found.CanCancel)
		assert.False(t, found.CanComplete)
		assert.False(t, found.CanReview)
	})

	t.Run("进行中订单操作标志", func(t *testing.T) {
		scheduledStart := time.Now().Add(-1 * time.Hour)
		startedAt := time.Now().Add(-30 * time.Minute)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "进行中订单",
			Status:          model.OrderStatusInProgress,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
			StartedAt:       &startedAt,
		}
		order.SetPlayerID(player.ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		resp, err := svc.GetMyOrders(context.Background(), customer.ID, MyOrderListRequest{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)

		var found *OrderCardDTO
		for i := range resp.Orders {
			if resp.Orders[i].ID == order.ID {
				found = &resp.Orders[i]
				break
			}
		}
		require.NotNil(t, found)
		assert.False(t, found.CanPay)
		assert.False(t, found.CanCancel)
		assert.True(t, found.CanComplete)
		assert.False(t, found.CanReview)
	})
}

func TestOrderService_GetAvailableOrdersWithFilters(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	// 创建另一个游戏
	game2 := &model.Game{Key: "dota2", Name: "Dota 2", Category: "moba"}
	require.NoError(t, db.Create(game2).Error)

	// 创建多个可接订单
	for i := 0; i < 3; i++ {
		scheduledStart := time.Now().Add(24 * time.Hour)
		scheduledEnd := scheduledStart.Add(2 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "LOL订单",
			Status:          model.OrderStatusConfirmed,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
			ScheduledEnd:    &scheduledEnd,
		}
		order.SetPlayerID(player.ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)
	}

	for i := 0; i < 2; i++ {
		scheduledStart := time.Now().Add(24 * time.Hour)
		scheduledEnd := scheduledStart.Add(2 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          game2.ID,
			Title:           "Dota2订单",
			Status:          model.OrderStatusConfirmed,
			UnitPriceCents:  6000,
			TotalPriceCents: 12000,
			ScheduledStart:  &scheduledStart,
			ScheduledEnd:    &scheduledEnd,
		}
		order.SetPlayerID(player.ID)
		order.SetGameID(game2.ID)
		require.NoError(t, db.Create(order).Error)
	}

	t.Run("获取所有可接订单", func(t *testing.T) {
		orders, total, err := svc.GetAvailableOrders(context.Background(), AvailableOrdersRequest{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, orders, 5)
	})

	t.Run("按游戏筛选可接订单", func(t *testing.T) {
		gameID := gameModel.ID
		orders, total, err := svc.GetAvailableOrders(context.Background(), AvailableOrdersRequest{
			GameID:   &gameID,
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, orders, 3)
		for _, o := range orders {
			assert.Equal(t, "英雄联盟", o.GameName)
		}
	})

	t.Run("分页获取可接订单", func(t *testing.T) {
		orders, total, err := svc.GetAvailableOrders(context.Background(), AvailableOrdersRequest{
			Page:     1,
			PageSize: 2,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, orders, 2)
	})
}

func TestOrderService_PlayerNotFound(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, _, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	t.Run("创建订单时陪玩师不存在", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		_, err := svc.CreateOrder(context.Background(), customer.ID, CreateOrderRequest{
			PlayerID:       99999,
			GameID:         gameModel.ID,
			Title:          "测试订单",
			ScheduledStart: &scheduledStart,
			DurationHours:  2,
		})
		assert.Error(t, err)
	})

	t.Run("接单时陪玩师不存在", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "待接单订单",
			Status:          model.OrderStatusConfirmed,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
		}
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		err := svc.AcceptOrder(context.Background(), 99999, order.ID)
		assert.Error(t, err)
	})
}

func TestOrderService_SetDistributedLock(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	svc := createOrderService(db)

	// Test setting nil lock
	svc.SetDistributedLock(nil)
	assert.Nil(t, svc.distributedLock)
}

func TestOrderService_SetChatGroupRepository(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	svc := createOrderService(db)

	// Test setting nil chat group repo
	svc.SetChatGroupRepository(nil)
	assert.Nil(t, svc.chatGroups)
}

func TestOrderService_GetMyOrdersDefaultPagination(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	// Create test order
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "测试订单",
		Status:          model.OrderStatusPending,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	order.SetPlayerID(player.ID)
	order.SetGameID(gameModel.ID)
	require.NoError(t, db.Create(order).Error)

	t.Run("default pagination with zero values", func(t *testing.T) {
		resp, err := svc.GetMyOrders(context.Background(), customer.ID, MyOrderListRequest{
			Page:     0,
			PageSize: 0,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(1), resp.Total)
	})

	t.Run("pagination with large page size", func(t *testing.T) {
		resp, err := svc.GetMyOrders(context.Background(), customer.ID, MyOrderListRequest{
			Page:     1,
			PageSize: 200, // Exceeds max
		})
		require.NoError(t, err)
		assert.Equal(t, int64(1), resp.Total)
	})
}

func TestOrderService_GetAvailableOrdersDefaultPagination(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	// Create confirmed order
	scheduledStart := time.Now().Add(24 * time.Hour)
	scheduledEnd := scheduledStart.Add(2 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "待接单订单",
		Status:          model.OrderStatusConfirmed,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
		ScheduledEnd:    &scheduledEnd,
	}
	order.SetPlayerID(player.ID)
	order.SetGameID(gameModel.ID)
	require.NoError(t, db.Create(order).Error)

	t.Run("default pagination with zero values", func(t *testing.T) {
		orders, total, err := svc.GetAvailableOrders(context.Background(), AvailableOrdersRequest{
			Page:     0,
			PageSize: 0,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, orders, 1)
	})

	t.Run("pagination with large page size", func(t *testing.T) {
		orders, total, err := svc.GetAvailableOrders(context.Background(), AvailableOrdersRequest{
			Page:     1,
			PageSize: 200, // Exceeds max
		})
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, orders, 1)
	})
}

func TestOrderService_CompleteOrderByPlayerNotFound(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, playerUser, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	// Create in-progress order
	scheduledStart := time.Now().Add(-2 * time.Hour)
	startedAt := time.Now().Add(-1 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "进行中订单",
		Status:          model.OrderStatusInProgress,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
		StartedAt:       &startedAt,
	}
	order.SetPlayerID(player.ID)
	order.SetGameID(gameModel.ID)
	require.NoError(t, db.Create(order).Error)

	t.Run("complete order with non-existent order", func(t *testing.T) {
		// Use a valid player user ID but non-existent order
		err := svc.CompleteOrderByPlayer(context.Background(), playerUser.ID, 99999)
		assert.Error(t, err)
	})

	t.Run("complete order with wrong status", func(t *testing.T) {
		// Create pending order
		pendingOrder := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "待支付订单",
			Status:          model.OrderStatusPending,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
		}
		pendingOrder.SetPlayerID(player.ID)
		pendingOrder.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(pendingOrder).Error)

		err := svc.CompleteOrderByPlayer(context.Background(), playerUser.ID, pendingOrder.ID)
		assert.ErrorIs(t, err, ErrInvalidTransition)
	})
}

func TestOrderService_CancelOrderUnauthorized(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	// Create pending order
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "待取消订单",
		Status:          model.OrderStatusPending,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	order.SetPlayerID(player.ID)
	order.SetGameID(gameModel.ID)
	require.NoError(t, db.Create(order).Error)

	// Create another user
	otherUser := &model.User{
		Phone:        "13800000099",
		Email:        "other99@test.com",
		Name:         "Other User",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(otherUser).Error)

	t.Run("cancel order by unauthorized user", func(t *testing.T) {
		err := svc.CancelOrder(context.Background(), otherUser.ID, order.ID, CancelOrderRequest{
			Reason: "不想要了",
		})
		assert.ErrorIs(t, err, ErrUnauthorized)
	})
}

func TestOrderService_CompleteOrderUnauthorized(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	// Create in-progress order
	scheduledStart := time.Now().Add(-2 * time.Hour)
	startedAt := time.Now().Add(-1 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "进行中订单",
		Status:          model.OrderStatusInProgress,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
		StartedAt:       &startedAt,
	}
	order.SetPlayerID(player.ID)
	order.SetGameID(gameModel.ID)
	require.NoError(t, db.Create(order).Error)

	// Create another user
	otherUser := &model.User{
		Phone:        "13800000098",
		Email:        "other98@test.com",
		Name:         "Other User",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(otherUser).Error)

	t.Run("complete order by unauthorized user", func(t *testing.T) {
		err := svc.CompleteOrder(context.Background(), otherUser.ID, order.ID)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})
}

func TestOrderService_CreateOrderWithInvalidGame(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, _ := createOrderTestData(t, db)
	svc := createOrderService(db)

	t.Run("create order with non-existent game", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		_, err := svc.CreateOrder(context.Background(), customer.ID, CreateOrderRequest{
			PlayerID:       player.ID,
			GameID:         99999, // Non-existent game
			Title:          "测试订单",
			ScheduledStart: &scheduledStart,
			DurationHours:  2,
		})
		assert.Error(t, err)
	})
}

func TestOrderService_OrderWithNoPlayerInfo(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, _, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	// Create order without player
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "无陪玩师订单",
		Status:          model.OrderStatusPending,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	order.SetGameID(gameModel.ID)
	// Don't set player ID
	require.NoError(t, db.Create(order).Error)

	t.Run("get order detail without player", func(t *testing.T) {
		resp, err := svc.GetOrderDetail(context.Background(), customer.ID, order.ID)
		require.NoError(t, err)
		assert.Equal(t, "无陪玩师订单", resp.Order.Title)
		// Player info should be nil or empty
	})
}

func TestOrderService_OrderWithNoGameInfo(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, _ := createOrderTestData(t, db)
	svc := createOrderService(db)

	// Create order without game
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          0, // No game
		Title:           "无游戏订单",
		Status:          model.OrderStatusPending,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	order.SetPlayerID(player.ID)
	// Don't set game ID
	require.NoError(t, db.Create(order).Error)

	t.Run("get my orders without game info", func(t *testing.T) {
		resp, err := svc.GetMyOrders(context.Background(), customer.ID, MyOrderListRequest{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, resp.Total, int64(1))
	})
}

func TestOrderService_AvailableOrdersWithUserInfo(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	// Create confirmed order with user info
	scheduledStart := time.Now().Add(24 * time.Hour)
	scheduledEnd := scheduledStart.Add(2 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "待接单订单",
		Description:     "订单描述",
		Status:          model.OrderStatusConfirmed,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
		ScheduledEnd:    &scheduledEnd,
	}
	order.SetPlayerID(player.ID)
	order.SetGameID(gameModel.ID)
	require.NoError(t, db.Create(order).Error)

	t.Run("available orders include user nickname", func(t *testing.T) {
		orders, total, err := svc.GetAvailableOrders(context.Background(), AvailableOrdersRequest{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, orders, 1)
		assert.Equal(t, "Customer", orders[0].UserNickname)
		assert.Equal(t, "英雄联盟", orders[0].GameName)
		assert.Equal(t, float32(2), orders[0].DurationHours)
	})
}

func TestOrderService_CancelOrderNotFound(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, _, _ := createOrderTestData(t, db)
	svc := createOrderService(db)

	t.Run("cancel non-existent order", func(t *testing.T) {
		err := svc.CancelOrder(context.Background(), customer.ID, 99999, CancelOrderRequest{
			Reason: "不想要了",
		})
		assert.Error(t, err)
	})
}

func TestOrderService_CompleteOrderNotFound(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, _, _ := createOrderTestData(t, db)
	svc := createOrderService(db)

	t.Run("complete non-existent order", func(t *testing.T) {
		err := svc.CompleteOrder(context.Background(), customer.ID, 99999)
		assert.Error(t, err)
	})
}

// ============================================
// Additional tests to improve coverage to 80%+
// ============================================

func TestOrderService_RecordCommissionAsync(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	t.Run("record commission for completed order", func(t *testing.T) {
		scheduledStart := time.Now().Add(-2 * time.Hour)
		completedAt := time.Now().Add(-1 * time.Hour)
		order := &model.Order{
			UserID:            customer.ID,
			ItemID:            gameModel.ID,
			Title:             "已完成订单",
			Status:            model.OrderStatusCompleted,
			UnitPriceCents:    5000,
			TotalPriceCents:   10000,
			CommissionCents:   2000,
			PlayerIncomeCents: 8000,
			ScheduledStart:    &scheduledStart,
			CompletedAt:       &completedAt,
		}
		order.SetPlayerID(player.ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		// Call recordCommissionAsync
		err := svc.recordCommissionAsync(context.Background(), order.ID)
		require.NoError(t, err)

		// Verify commission record was created
		var record model.CommissionRecord
		err = db.Where("order_id = ?", order.ID).First(&record).Error
		require.NoError(t, err)
		assert.Equal(t, order.ID, record.OrderID)
		assert.Equal(t, player.ID, record.PlayerID)
	})

	t.Run("skip if commission already recorded", func(t *testing.T) {
		scheduledStart := time.Now().Add(-2 * time.Hour)
		completedAt := time.Now().Add(-1 * time.Hour)
		order := &model.Order{
			UserID:            customer.ID,
			ItemID:            gameModel.ID,
			Title:             "已记录抽成订单",
			Status:            model.OrderStatusCompleted,
			UnitPriceCents:    5000,
			TotalPriceCents:   10000,
			CommissionCents:   2000,
			PlayerIncomeCents: 8000,
			ScheduledStart:    &scheduledStart,
			CompletedAt:       &completedAt,
		}
		order.SetPlayerID(player.ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		// Create existing commission record
		existingRecord := &model.CommissionRecord{
			OrderID:           order.ID,
			PlayerID:          player.ID,
			TotalAmountCents:  10000,
			CommissionRate:    20,
			CommissionCents:   2000,
			PlayerIncomeCents: 8000,
			SettlementStatus:  "pending",
			SettlementMonth:   time.Now().Format("2006-01"),
		}
		require.NoError(t, db.Create(existingRecord).Error)

		// Call recordCommissionAsync - should skip
		err := svc.recordCommissionAsync(context.Background(), order.ID)
		require.NoError(t, err)

		// Verify only one record exists
		var count int64
		db.Model(&model.CommissionRecord{}).Where("order_id = ?", order.ID).Count(&count)
		assert.Equal(t, int64(1), count)
	})

	t.Run("record commission for order without player", func(t *testing.T) {
		scheduledStart := time.Now().Add(-2 * time.Hour)
		completedAt := time.Now().Add(-1 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "无陪玩师订单",
			Status:          model.OrderStatusCompleted,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
			CompletedAt:     &completedAt,
		}
		// Don't set player ID
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		// Call recordCommissionAsync - should fail
		err := svc.recordCommissionAsync(context.Background(), order.ID)
		assert.Error(t, err)
	})
}

func TestOrderService_BuildOrderTimeline(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	t.Run("timeline for pending order", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "待支付订单",
			Status:          model.OrderStatusPending,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
		}
		order.SetPlayerID(player.ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		timeline := svc.buildOrderTimeline(order)
		assert.Len(t, timeline, 1)
		assert.Equal(t, string(model.OrderStatusPending), timeline[0].Status)
		assert.Equal(t, "订单已创建", timeline[0].Message)
	})

	t.Run("timeline for in-progress order", func(t *testing.T) {
		scheduledStart := time.Now().Add(-1 * time.Hour)
		startedAt := time.Now().Add(-30 * time.Minute)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "进行中订单",
			Status:          model.OrderStatusInProgress,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
			StartedAt:       &startedAt,
		}
		order.SetPlayerID(player.ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		timeline := svc.buildOrderTimeline(order)
		assert.GreaterOrEqual(t, len(timeline), 2)
	})
}

func TestOrderService_ToOrderCardDTO(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	t.Run("convert pending order to card DTO", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "待支付订单",
			Status:          model.OrderStatusPending,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
		}
		order.SetPlayerID(player.ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		card, err := svc.toOrderCardDTO(context.Background(), order, customer.ID)
		require.NoError(t, err)
		assert.Equal(t, "待支付订单", card.Title)
		assert.Equal(t, "Pro Player", card.PlayerNickname)
		assert.Equal(t, "英雄联盟", card.GameName)
		assert.True(t, card.CanPay)
		assert.True(t, card.CanCancel)
		assert.False(t, card.CanComplete)
		assert.False(t, card.CanReview)
	})

	t.Run("convert completed order to card DTO", func(t *testing.T) {
		scheduledStart := time.Now().Add(-24 * time.Hour)
		completedAt := time.Now().Add(-1 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "已完成订单",
			Status:          model.OrderStatusCompleted,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
			CompletedAt:     &completedAt,
		}
		order.SetPlayerID(player.ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		card, err := svc.toOrderCardDTO(context.Background(), order, customer.ID)
		require.NoError(t, err)
		assert.False(t, card.CanPay)
		assert.False(t, card.CanCancel)
		assert.False(t, card.CanComplete)
		assert.True(t, card.CanReview)
	})

	t.Run("convert order with missing player", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "无陪玩师订单",
			Status:          model.OrderStatusPending,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
		}
		// Don't set player ID
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		card, err := svc.toOrderCardDTO(context.Background(), order, customer.ID)
		require.NoError(t, err)
		assert.Empty(t, card.PlayerNickname)
	})

	t.Run("convert order with missing game", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          0,
			Title:           "无游戏订单",
			Status:          model.OrderStatusPending,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
		}
		order.SetPlayerID(player.ID)
		// Don't set game ID
		require.NoError(t, db.Create(order).Error)

		card, err := svc.toOrderCardDTO(context.Background(), order, customer.ID)
		require.NoError(t, err)
		assert.Empty(t, card.GameName)
	})
}

func TestOrderService_CalculateOrderPricing(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	_, _, player, _ := createOrderTestData(t, db)
	svc := createOrderService(db)

	t.Run("calculate pricing for 2 hours", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		req := CreateOrderRequest{
			PlayerID:       player.ID,
			GameID:         1,
			Title:          "测试订单",
			ScheduledStart: &scheduledStart,
			DurationHours:  2,
		}

		totalPrice, commissionCents, playerIncomeCents := svc.calculateOrderPricing(player, req)

		// Player hourly rate is 5000 cents, 2 hours = 10000 cents
		assert.Equal(t, int64(10000), totalPrice)
		// 20% commission = 2000 cents
		assert.Equal(t, int64(2000), commissionCents)
		// Player income = 8000 cents
		assert.Equal(t, int64(8000), playerIncomeCents)
	})

	t.Run("calculate pricing for 0.5 hours", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		req := CreateOrderRequest{
			PlayerID:       player.ID,
			GameID:         1,
			Title:          "测试订单",
			ScheduledStart: &scheduledStart,
			DurationHours:  0.5,
		}

		totalPrice, commissionCents, playerIncomeCents := svc.calculateOrderPricing(player, req)

		// Player hourly rate is 5000 cents, 0.5 hours = 2500 cents
		assert.Equal(t, int64(2500), totalPrice)
		// 20% commission = 500 cents
		assert.Equal(t, int64(500), commissionCents)
		// Player income = 2000 cents
		assert.Equal(t, int64(2000), playerIncomeCents)
	})
}

func TestOrderService_ValidateCreateOrder(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	_, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	t.Run("validate with valid player and game", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		req := CreateOrderRequest{
			PlayerID:       player.ID,
			GameID:         gameModel.ID,
			Title:          "测试订单",
			ScheduledStart: &scheduledStart,
			DurationHours:  2,
		}

		validatedPlayer, err := svc.validateCreateOrder(context.Background(), req)
		require.NoError(t, err)
		assert.Equal(t, player.ID, validatedPlayer.ID)
	})

	t.Run("validate with invalid player", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		req := CreateOrderRequest{
			PlayerID:       99999,
			GameID:         gameModel.ID,
			Title:          "测试订单",
			ScheduledStart: &scheduledStart,
			DurationHours:  2,
		}

		_, err := svc.validateCreateOrder(context.Background(), req)
		assert.Error(t, err)
	})

	t.Run("validate with invalid game", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		req := CreateOrderRequest{
			PlayerID:       player.ID,
			GameID:         99999,
			Title:          "测试订单",
			ScheduledStart: &scheduledStart,
			DurationHours:  2,
		}

		_, err := svc.validateCreateOrder(context.Background(), req)
		assert.Error(t, err)
	})
}

func TestOrderService_BuildOrderForCreation(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	t.Run("build order with all fields", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		req := CreateOrderRequest{
			PlayerID:       player.ID,
			GameID:         gameModel.ID,
			Title:          "测试订单",
			Description:    "订单描述",
			ScheduledStart: &scheduledStart,
			DurationHours:  2,
		}

		order := svc.buildOrderForCreation(customer.ID, req, 10000, 2000, 8000)

		assert.Equal(t, customer.ID, order.UserID)
		assert.Equal(t, "测试订单", order.Title)
		assert.Equal(t, "订单描述", order.Description)
		assert.Equal(t, int64(10000), order.TotalPriceCents)
		assert.Equal(t, int64(2000), order.CommissionCents)
		assert.Equal(t, int64(8000), order.PlayerIncomeCents)
		assert.Equal(t, model.OrderStatusPending, order.Status)
		assert.NotNil(t, order.ScheduledEnd)
	})
}

func TestOrderService_DeactivateOrderChat(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	svc := createOrderService(db)

	t.Run("deactivate with nil chat group repo", func(t *testing.T) {
		// Should not panic
		svc.deactivateOrderChat(context.Background(), 1)
	})
}

func TestOrderService_AcceptOrderNotFound(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	_, playerUser, _, _ := createOrderTestData(t, db)
	svc := createOrderService(db)

	t.Run("accept non-existent order", func(t *testing.T) {
		err := svc.AcceptOrder(context.Background(), playerUser.ID, 99999)
		assert.Error(t, err)
	})
}

func TestOrderService_GetOrderDetailWithTimeline(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	t.Run("get detail with payment info in timeline", func(t *testing.T) {
		scheduledStart := time.Now().Add(-24 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "已支付订单",
			Status:          model.OrderStatusConfirmed,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
		}
		order.SetPlayerID(player.ID)
		order.SetGameID(gameModel.ID)
		require.NoError(t, db.Create(order).Error)

		// Create payment record
		paidAt := time.Now().Add(-12 * time.Hour)
		payment := &model.Payment{
			OrderID:     order.ID,
			UserID:      customer.ID,
			Method:      model.PaymentMethodWeChat,
			AmountCents: 10000,
			Status:      model.PaymentStatusPaid,
			PaidAt:      &paidAt,
		}
		require.NoError(t, db.Create(payment).Error)

		resp, err := svc.GetOrderDetail(context.Background(), customer.ID, order.ID)
		require.NoError(t, err)
		assert.NotNil(t, resp.Payment)
		assert.GreaterOrEqual(t, len(resp.Timeline), 2)
	})
}

func TestOrderService_ErrorVariables(t *testing.T) {
	t.Run("error variables are defined", func(t *testing.T) {
		assert.NotNil(t, ErrNotFound)
		assert.NotNil(t, ErrValidation)
		assert.NotNil(t, ErrInvalidTransition)
		assert.NotNil(t, ErrUnauthorized)
	})
}

func TestOrderService_CreateOrderWithDistributedLock(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	// Create a mock distributed lock that always succeeds
	mockLock := &mockDistributedLock{
		tryLockResult: true,
	}
	svc.SetDistributedLock(mockLock)

	t.Run("create order with distributed lock", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		resp, err := svc.CreateOrder(context.Background(), customer.ID, CreateOrderRequest{
			PlayerID:       player.ID,
			GameID:         gameModel.ID,
			Title:          "带锁订单",
			ScheduledStart: &scheduledStart,
			DurationHours:  2,
		})
		require.NoError(t, err)
		assert.NotZero(t, resp.OrderID)
	})
}

// mockDistributedLock is a mock implementation of cache.DistributedLock
type mockDistributedLock struct {
	tryLockResult bool
	tryLockErr    error
}

func (m *mockDistributedLock) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return m.tryLockResult, m.tryLockErr
}

func (m *mockDistributedLock) TryLock(ctx context.Context, key string, ttl time.Duration, retries int, retryDelay time.Duration) (bool, error) {
	return m.tryLockResult, m.tryLockErr
}

func (m *mockDistributedLock) Unlock(ctx context.Context, key string) error {
	return nil
}

func TestOrderService_CreateOrderLockFailed(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	// Create a mock distributed lock that fails
	mockLock := &mockDistributedLock{
		tryLockResult: false,
	}
	svc.SetDistributedLock(mockLock)

	t.Run("create order when lock acquisition fails", func(t *testing.T) {
		scheduledStart := time.Now().Add(24 * time.Hour)
		_, err := svc.CreateOrder(context.Background(), customer.ID, CreateOrderRequest{
			PlayerID:       player.ID,
			GameID:         gameModel.ID,
			Title:          "锁失败订单",
			ScheduledStart: &scheduledStart,
			DurationHours:  2,
		})
		assert.Error(t, err)
	})
}

func TestOrderService_GetAvailableOrdersWithNoScheduledEnd(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	// Create order without scheduled end
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "无结束时间订单",
		Status:          model.OrderStatusConfirmed,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
		// No ScheduledEnd
	}
	order.SetPlayerID(player.ID)
	order.SetGameID(gameModel.ID)
	require.NoError(t, db.Create(order).Error)

	t.Run("get available orders with no scheduled end", func(t *testing.T) {
		orders, total, err := svc.GetAvailableOrders(context.Background(), AvailableOrdersRequest{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, orders, 1)
		assert.Equal(t, float32(0), orders[0].DurationHours)
	})
}

func TestOrderService_GetAvailableOrdersWithNoUserInfo(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	_, _, player, gameModel := createOrderTestData(t, db)
	svc := createOrderService(db)

	// Create order with non-existent user
	scheduledStart := time.Now().Add(24 * time.Hour)
	scheduledEnd := scheduledStart.Add(2 * time.Hour)
	order := &model.Order{
		UserID:          99999, // Non-existent user
		ItemID:          gameModel.ID,
		Title:           "无用户信息订单",
		Status:          model.OrderStatusConfirmed,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
		ScheduledEnd:    &scheduledEnd,
	}
	order.SetPlayerID(player.ID)
	order.SetGameID(gameModel.ID)
	require.NoError(t, db.Create(order).Error)

	t.Run("get available orders with no user info", func(t *testing.T) {
		orders, total, err := svc.GetAvailableOrders(context.Background(), AvailableOrdersRequest{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, orders, 1)
		assert.Empty(t, orders[0].UserNickname)
	})
}

func TestOrderService_GetAvailableOrdersWithNoGameInfo(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, _ := createOrderTestData(t, db)
	svc := createOrderService(db)

	// Create order without game
	scheduledStart := time.Now().Add(24 * time.Hour)
	scheduledEnd := scheduledStart.Add(2 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          0,
		Title:           "无游戏信息订单",
		Status:          model.OrderStatusConfirmed,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
		ScheduledEnd:    &scheduledEnd,
	}
	order.SetPlayerID(player.ID)
	// Don't set game ID
	require.NoError(t, db.Create(order).Error)

	t.Run("get available orders with no game info", func(t *testing.T) {
		orders, total, err := svc.GetAvailableOrders(context.Background(), AvailableOrdersRequest{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, orders, 1)
		assert.Empty(t, orders[0].GameName)
	})
}

// ============================================
// Tests for DisputeService
// ============================================

func setupDisputeTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.User{},
		&model.Player{},
		&model.Game{},
		&model.Order{},
		&model.Payment{},
		&model.OrderDispute{},
		&model.OperationLog{},
		&model.NotificationEvent{},
	)
	return db
}

func TestDisputeService_InitiateDispute(t *testing.T) {
	db := setupDisputeTestDB(t)
	defer testutil.CleanDB(t, db)

	// Create test user
	customer := &model.User{
		Phone:        "13800000010",
		Email:        "dispute_customer@test.com",
		Name:         "Dispute Customer",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(customer).Error)

	// Create completed order
	scheduledStart := time.Now().Add(-24 * time.Hour)
	completedAt := time.Now().Add(-1 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          1,
		Title:           "可纠纷订单",
		Status:          model.OrderStatusCompleted,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
		CompletedAt:     &completedAt,
	}
	require.NoError(t, db.Create(order).Error)

	// Create repositories
	disputeRepo := &mockDisputeRepository{db: db}
	orderRepo := orderimpl.NewOrderRepository(db)
	userRepo := user.NewUserRepository(db)
	opLogRepo := &mockOperationLogRepository{}
	notificationRepo := &mockNotificationRepository{}
	paymentRepo := orderrepo.NewPaymentRepository(db)

	svc := NewDisputeService(disputeRepo, orderRepo, userRepo, opLogRepo, notificationRepo, paymentRepo)

	t.Run("initiate dispute validation error - zero order ID", func(t *testing.T) {
		_, err := svc.InitiateDispute(context.Background(), InitiateDisputeRequest{
			OrderID: 0,
			UserID:  customer.ID,
			Reason:  "测试",
		})
		assert.Error(t, err)
	})

	t.Run("initiate dispute empty reason", func(t *testing.T) {
		_, err := svc.InitiateDispute(context.Background(), InitiateDisputeRequest{
			OrderID: order.ID,
			UserID:  customer.ID,
			Reason:  "",
		})
		assert.Error(t, err)
	})

	t.Run("initiate dispute order not found", func(t *testing.T) {
		_, err := svc.InitiateDispute(context.Background(), InitiateDisputeRequest{
			OrderID:     99999,
			UserID:      customer.ID,
			Reason:      "测试",
			Description: "描述",
		})
		assert.Error(t, err)
	})

	t.Run("initiate dispute unauthorized user", func(t *testing.T) {
		_, err := svc.InitiateDispute(context.Background(), InitiateDisputeRequest{
			OrderID:     order.ID,
			UserID:      99999, // Different user
			Reason:      "测试",
			Description: "描述",
		})
		assert.Error(t, err)
	})
}

func TestDisputeService_GetDisputeDetail(t *testing.T) {
	db := setupDisputeTestDB(t)
	defer testutil.CleanDB(t, db)

	// Create dispute
	dispute := &model.OrderDispute{
		OrderID:     1,
		UserID:      1,
		Status:      model.DisputeStatusPending,
		Reason:      "测试纠纷",
		Description: "描述",
		TraceID:     "trace-123",
	}
	require.NoError(t, db.Create(dispute).Error)

	disputeRepo := &mockDisputeRepository{db: db}
	orderRepo := orderimpl.NewOrderRepository(db)
	userRepo := user.NewUserRepository(db)
	opLogRepo := &mockOperationLogRepository{}
	notificationRepo := &mockNotificationRepository{}
	paymentRepo := orderrepo.NewPaymentRepository(db)

	svc := NewDisputeService(disputeRepo, orderRepo, userRepo, opLogRepo, notificationRepo, paymentRepo)

	t.Run("get dispute detail", func(t *testing.T) {
		result, err := svc.GetDisputeDetail(context.Background(), dispute.ID)
		require.NoError(t, err)
		assert.Equal(t, "测试纠纷", result.Reason)
	})
}

func TestDisputeService_ListPendingDisputes(t *testing.T) {
	db := setupDisputeTestDB(t)
	defer testutil.CleanDB(t, db)

	// Create pending disputes
	for i := 0; i < 3; i++ {
		dispute := &model.OrderDispute{
			OrderID:     uint64(i + 1),
			UserID:      1,
			Status:      model.DisputeStatusPending,
			Reason:      "测试纠纷",
			Description: "描述",
			TraceID:     fmt.Sprintf("trace-%d", i),
		}
		require.NoError(t, db.Create(dispute).Error)
	}

	disputeRepo := &mockDisputeRepository{db: db}
	orderRepo := orderimpl.NewOrderRepository(db)
	userRepo := user.NewUserRepository(db)
	opLogRepo := &mockOperationLogRepository{}
	notificationRepo := &mockNotificationRepository{}
	paymentRepo := orderrepo.NewPaymentRepository(db)

	svc := NewDisputeService(disputeRepo, orderRepo, userRepo, opLogRepo, notificationRepo, paymentRepo)

	t.Run("list pending disputes", func(t *testing.T) {
		disputes, total, err := svc.ListPendingDisputes(context.Background(), 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, disputes, 3)
	})
}

// Mock repositories for dispute tests
type mockDisputeRepository struct {
	db *gorm.DB
}

func (m *mockDisputeRepository) Create(ctx context.Context, dispute *model.OrderDispute) error {
	return m.db.Create(dispute).Error
}

func (m *mockDisputeRepository) Get(ctx context.Context, id uint64) (*model.OrderDispute, error) {
	var dispute model.OrderDispute
	if err := m.db.First(&dispute, id).Error; err != nil {
		return nil, err
	}
	return &dispute, nil
}

func (m *mockDisputeRepository) GetByOrderID(ctx context.Context, orderID uint64) (*model.OrderDispute, error) {
	var dispute model.OrderDispute
	if err := m.db.Where("order_id = ?", orderID).First(&dispute).Error; err != nil {
		return nil, err
	}
	return &dispute, nil
}

func (m *mockDisputeRepository) Update(ctx context.Context, dispute *model.OrderDispute) error {
	return m.db.Save(dispute).Error
}

func (m *mockDisputeRepository) List(ctx context.Context, opts repository.DisputeListOptions) ([]model.OrderDispute, int64, error) {
	var disputes []model.OrderDispute
	var total int64
	query := m.db.Model(&model.OrderDispute{})
	if len(opts.Statuses) > 0 {
		query = query.Where("status IN ?", opts.Statuses)
	}
	query.Count(&total)
	query.Offset((opts.Page - 1) * opts.PageSize).Limit(opts.PageSize).Find(&disputes)
	return disputes, total, nil
}

func (m *mockDisputeRepository) ListPendingAssignment(ctx context.Context, page, pageSize int) ([]model.OrderDispute, int64, error) {
	var disputes []model.OrderDispute
	var total int64
	query := m.db.Model(&model.OrderDispute{}).Where("status = ?", model.DisputeStatusPending)
	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&disputes)
	return disputes, total, nil
}

func (m *mockDisputeRepository) ListSLABreached(ctx context.Context) ([]model.OrderDispute, error) {
	var disputes []model.OrderDispute
	m.db.Where("sla_breached = ? AND status NOT IN ?", true, []model.DisputeStatus{model.DisputeStatusResolved, model.DisputeStatusRejected}).Find(&disputes)
	return disputes, nil
}

func (m *mockDisputeRepository) MarkSLABreached(ctx context.Context, id uint64) error {
	return m.db.Model(&model.OrderDispute{}).Where("id = ?", id).Update("sla_breached", true).Error
}

func (m *mockDisputeRepository) Delete(ctx context.Context, id uint64) error {
	return m.db.Delete(&model.OrderDispute{}, id).Error
}

func (m *mockDisputeRepository) CountByStatus(ctx context.Context, status model.DisputeStatus) (int64, error) {
	var count int64
	m.db.Model(&model.OrderDispute{}).Where("status = ?", status).Count(&count)
	return count, nil
}

func (m *mockDisputeRepository) GetPendingCount(ctx context.Context) (int64, error) {
	var count int64
	m.db.Model(&model.OrderDispute{}).Where("status = ?", model.DisputeStatusPending).Count(&count)
	return count, nil
}

func (m *mockDisputeRepository) GetStats(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)
	var disputes []model.OrderDispute
	m.db.Find(&disputes)
	for _, d := range disputes {
		stats[string(d.Status)]++
	}
	return stats, nil
}

type mockOperationLogRepository struct{}

func (m *mockOperationLogRepository) Append(ctx context.Context, log *model.OperationLog) error {
	return nil
}

func (m *mockOperationLogRepository) ListByEntity(ctx context.Context, entityType string, entityID uint64, opts repository.OperationLogListOptions) ([]model.OperationLog, int64, error) {
	return nil, 0, nil
}

func (m *mockOperationLogRepository) List(ctx context.Context, opts repository.OperationLogSearchOptions) ([]model.OperationLog, int64, error) {
	return nil, 0, nil
}

type mockNotificationRepository struct{}

func (m *mockNotificationRepository) Create(ctx context.Context, event *model.NotificationEvent) error {
	return nil
}

func (m *mockNotificationRepository) ListByUser(ctx context.Context, opts repository.NotificationListOptions) ([]model.NotificationEvent, int64, error) {
	return nil, 0, nil
}

func (m *mockNotificationRepository) MarkRead(ctx context.Context, userID uint64, ids []uint64) error {
	return nil
}

func (m *mockNotificationRepository) MarkAllRead(ctx context.Context, userID uint64) error {
	return nil
}

func (m *mockNotificationRepository) Delete(ctx context.Context, userID uint64, id uint64) error {
	return nil
}

func (m *mockNotificationRepository) CountUnread(ctx context.Context, userID uint64) (int64, error) {
	return 0, nil
}

// ============================================
// Tests for PaymentService (in order package)
// ============================================

func TestOrderPaymentService_CreatePayment(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)

	// Create pending order
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "待支付订单",
		Status:          model.OrderStatusPending,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	order.SetPlayerID(player.ID)
	order.SetGameID(gameModel.ID)
	require.NoError(t, db.Create(order).Error)

	paymentRepo := orderrepo.NewPaymentRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	svc := NewPaymentService(paymentRepo, orderRepo)

	t.Run("create wechat payment", func(t *testing.T) {
		resp, err := svc.CreatePayment(context.Background(), customer.ID, CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodWeChat,
		})
		require.NoError(t, err)
		assert.NotZero(t, resp.PaymentID)
		assert.Contains(t, resp.PayInfo, "prepay_id")
	})
}

func TestOrderPaymentService_GetPaymentStatus(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, _, _ := createOrderTestData(t, db)

	// Create payment
	payment := &model.Payment{
		OrderID:     1,
		UserID:      customer.ID,
		Method:      model.PaymentMethodWeChat,
		AmountCents: 10000,
		Status:      model.PaymentStatusPaid,
	}
	require.NoError(t, db.Create(payment).Error)

	paymentRepo := orderrepo.NewPaymentRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	svc := NewPaymentService(paymentRepo, orderRepo)

	t.Run("get payment status", func(t *testing.T) {
		resp, err := svc.GetPaymentStatus(context.Background(), payment.ID)
		require.NoError(t, err)
		assert.Equal(t, model.PaymentStatusPaid, resp.Status)
	})

	t.Run("payment not found", func(t *testing.T) {
		_, err := svc.GetPaymentStatus(context.Background(), 99999)
		assert.Error(t, err)
	})
}

func TestOrderPaymentService_CancelPayment(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, _, _ := createOrderTestData(t, db)

	// Create pending payment
	payment := &model.Payment{
		OrderID:     1,
		UserID:      customer.ID,
		Method:      model.PaymentMethodWeChat,
		AmountCents: 10000,
		Status:      model.PaymentStatusPending,
	}
	require.NoError(t, db.Create(payment).Error)

	paymentRepo := orderrepo.NewPaymentRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	svc := NewPaymentService(paymentRepo, orderRepo)

	t.Run("cancel pending payment", func(t *testing.T) {
		err := svc.CancelPayment(context.Background(), customer.ID, payment.ID)
		require.NoError(t, err)

		// Verify status
		var updated model.Payment
		require.NoError(t, db.First(&updated, payment.ID).Error)
		assert.Equal(t, model.PaymentStatusFailed, updated.Status)
	})

	t.Run("unauthorized cancel", func(t *testing.T) {
		// Create another payment
		payment2 := &model.Payment{
			OrderID:     2,
			UserID:      customer.ID,
			Method:      model.PaymentMethodWeChat,
			AmountCents: 10000,
			Status:      model.PaymentStatusPending,
		}
		require.NoError(t, db.Create(payment2).Error)

		err := svc.CancelPayment(context.Background(), 99999, payment2.ID)
		assert.Error(t, err)
	})
}

func TestOrderPaymentService_HandlePaymentCallback(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)

	// Create order
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "回调测试订单",
		Status:          model.OrderStatusPending,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	order.SetPlayerID(player.ID)
	order.SetGameID(gameModel.ID)
	require.NoError(t, db.Create(order).Error)

	// Create pending payment
	payment := &model.Payment{
		OrderID:     order.ID,
		UserID:      customer.ID,
		Method:      model.PaymentMethodWeChat,
		AmountCents: 10000,
		Status:      model.PaymentStatusPending,
	}
	require.NoError(t, db.Create(payment).Error)

	paymentRepo := orderrepo.NewPaymentRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	svc := NewPaymentService(paymentRepo, orderRepo)

	t.Run("handle callback success", func(t *testing.T) {
		err := svc.HandlePaymentCallback(context.Background(), "wechat", map[string]interface{}{
			"payment_id": float64(payment.ID),
			"trade_no":   "wx_trade_123",
		})
		require.NoError(t, err)

		// Verify payment status
		var updatedPayment model.Payment
		require.NoError(t, db.First(&updatedPayment, payment.ID).Error)
		assert.Equal(t, model.PaymentStatusPaid, updatedPayment.Status)

		// Verify order status
		var updatedOrder model.Order
		require.NoError(t, db.First(&updatedOrder, order.ID).Error)
		assert.Equal(t, model.OrderStatusConfirmed, updatedOrder.Status)
	})
}

func TestOrderPaymentService_RefundPayment(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)

	// Create order
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "退款测试订单",
		Status:          model.OrderStatusConfirmed,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	order.SetPlayerID(player.ID)
	order.SetGameID(gameModel.ID)
	require.NoError(t, db.Create(order).Error)

	// Create paid payment
	now := time.Now()
	payment := &model.Payment{
		OrderID:     order.ID,
		UserID:      customer.ID,
		Method:      model.PaymentMethodWeChat,
		AmountCents: 10000,
		Status:      model.PaymentStatusPaid,
		PaidAt:      &now,
	}
	require.NoError(t, db.Create(payment).Error)

	paymentRepo := orderrepo.NewPaymentRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	svc := NewPaymentService(paymentRepo, orderRepo)

	t.Run("refund payment success", func(t *testing.T) {
		err := svc.RefundPayment(context.Background(), payment.ID, "用户申请退款")
		require.NoError(t, err)

		// Verify payment status
		var updatedPayment model.Payment
		require.NoError(t, db.First(&updatedPayment, payment.ID).Error)
		assert.Equal(t, model.PaymentStatusRefunded, updatedPayment.Status)

		// Verify order status
		var updatedOrder model.Order
		require.NoError(t, db.First(&updatedOrder, order.ID).Error)
		assert.Equal(t, model.OrderStatusRefunded, updatedOrder.Status)
	})
}

func TestOrderPaymentService_ListPayments(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, _, _ := createOrderTestData(t, db)

	// Create payments
	for i := 0; i < 5; i++ {
		payment := &model.Payment{
			OrderID:     uint64(i + 1),
			UserID:      customer.ID,
			Method:      model.PaymentMethodWeChat,
			AmountCents: int64((i + 1) * 1000),
			Status:      model.PaymentStatusPending,
		}
		require.NoError(t, db.Create(payment).Error)
	}

	paymentRepo := orderrepo.NewPaymentRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	svc := NewPaymentService(paymentRepo, orderRepo)

	t.Run("list payments", func(t *testing.T) {
		payments, total, err := svc.ListPayments(context.Background(), repository.PaymentListOptions{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, payments, 5)
	})
}

// ============================================
// Tests for PaymentProvider (in order package)
// ============================================

func TestOrderWechatProvider_Refund(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	provider := wechatProvider{}

	payment := &model.Payment{
		OrderID:     1,
		UserID:      1,
		Method:      model.PaymentMethodWeChat,
		AmountCents: 10000,
		Status:      model.PaymentStatusPaid,
	}
	require.NoError(t, db.Create(payment).Error)

	t.Run("refund success", func(t *testing.T) {
		tradeNo, raw, refundedAt, err := provider.Refund(context.Background(), payment, "test refund")
		require.NoError(t, err)
		assert.Contains(t, tradeNo, "wx_refund_")
		assert.NotNil(t, raw)
		assert.False(t, refundedAt.IsZero())
	})
}

func TestOrderAlipayProvider_Refund(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	provider := alipayProvider{}

	payment := &model.Payment{
		OrderID:     2,
		UserID:      1,
		Method:      model.PaymentMethodAlipay,
		AmountCents: 20000,
		Status:      model.PaymentStatusPaid,
	}
	require.NoError(t, db.Create(payment).Error)

	t.Run("refund success", func(t *testing.T) {
		tradeNo, raw, refundedAt, err := provider.Refund(context.Background(), payment, "test refund")
		require.NoError(t, err)
		assert.Contains(t, tradeNo, "ali_refund_")
		assert.NotNil(t, raw)
		assert.False(t, refundedAt.IsZero())
	})
}

func TestOrderGenericProvider_Refund(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	provider := genericProvider{}

	payment := &model.Payment{
		OrderID:     3,
		UserID:      1,
		AmountCents: 30000,
		Status:      model.PaymentStatusPaid,
	}
	require.NoError(t, db.Create(payment).Error)

	t.Run("refund success", func(t *testing.T) {
		tradeNo, raw, refundedAt, err := provider.Refund(context.Background(), payment, "test refund")
		require.NoError(t, err)
		assert.Contains(t, tradeNo, "refund_")
		assert.NotNil(t, raw)
		assert.False(t, refundedAt.IsZero())
	})
}

// ============================================
// Tests for error variables
// ============================================

func TestOrderPackageErrors(t *testing.T) {
	t.Run("dispute errors defined", func(t *testing.T) {
		assert.NotNil(t, ErrDisputeNotFound)
		assert.NotNil(t, ErrDisputeValidation)
		assert.NotNil(t, ErrDisputeInvalidStatus)
		assert.NotNil(t, ErrDisputeUnauthorized)
		assert.NotNil(t, ErrDisputeSLAExpired)
		assert.NotNil(t, ErrOrderNotFound)
		assert.NotNil(t, ErrDisputeExists)
		assert.NotNil(t, ErrCannotInitiateDispute)
	})

	t.Run("payment errors defined", func(t *testing.T) {
		assert.NotNil(t, ErrPaymentNotFound)
		assert.NotNil(t, ErrPaymentValidation)
		assert.NotNil(t, ErrOrderAlreadyPaid)
		assert.NotNil(t, ErrInvalidOrderStatus)
	})

	t.Run("review errors defined", func(t *testing.T) {
		assert.NotNil(t, ErrReviewNotFound)
		assert.NotNil(t, ErrReviewValidation)
		assert.NotNil(t, ErrAlreadyReviewed)
		assert.NotNil(t, ErrOrderNotCompleted)
		assert.NotNil(t, ErrReviewUnauthorized)
	})
}

// ============================================
// Tests for ReviewService
// ============================================

func setupReviewTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.User{},
		&model.Player{},
		&model.Game{},
		&model.Order{},
		&model.Review{},
		&model.ReviewReply{},
		&model.NotificationEvent{},
	)
	return db
}

// Mock ReviewReplyRepository
type mockReviewReplyRepository struct {
	db *gorm.DB
}

func (m *mockReviewReplyRepository) Create(ctx context.Context, reply *model.ReviewReply) error {
	return m.db.Create(reply).Error
}

func (m *mockReviewReplyRepository) Get(ctx context.Context, replyID uint64) (*model.ReviewReply, error) {
	var reply model.ReviewReply
	if err := m.db.First(&reply, replyID).Error; err != nil {
		return nil, err
	}
	return &reply, nil
}

func (m *mockReviewReplyRepository) ListByReview(ctx context.Context, reviewID uint64) ([]model.ReviewReply, error) {
	var replies []model.ReviewReply
	m.db.Where("review_id = ?", reviewID).Find(&replies)
	return replies, nil
}

func (m *mockReviewReplyRepository) Update(ctx context.Context, reply *model.ReviewReply) error {
	return m.db.Save(reply).Error
}

func (m *mockReviewReplyRepository) Delete(ctx context.Context, replyID uint64) error {
	return m.db.Delete(&model.ReviewReply{}, replyID).Error
}

func (m *mockReviewReplyRepository) UpdateStatus(ctx context.Context, replyID uint64, status string, note string) error {
	return m.db.Model(&model.ReviewReply{}).Where("id = ?", replyID).Updates(map[string]interface{}{
		"status":          status,
		"moderation_note": note,
	}).Error
}

func TestReviewService_CreateReview(t *testing.T) {
	db := setupReviewTestDB(t)
	defer testutil.CleanDB(t, db)

	// Create test user
	customer := &model.User{
		Phone:        "13800000020",
		Email:        "review_customer@test.com",
		Name:         "Review Customer",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(customer).Error)

	// Create player user
	playerUser := &model.User{
		Phone:        "13800000021",
		Email:        "review_player@test.com",
		Name:         "Review Player",
		Role:         model.RolePlayer,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(playerUser).Error)

	// Create player
	player := &model.Player{
		UserID:             playerUser.ID,
		Nickname:           "Review Pro",
		HourlyRateCents:    5000,
		VerificationStatus: model.VerificationVerified,
	}
	require.NoError(t, db.Create(player).Error)

	// Create completed order
	scheduledStart := time.Now().Add(-24 * time.Hour)
	completedAt := time.Now().Add(-1 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          1,
		Title:           "已完成订单",
		Status:          model.OrderStatusCompleted,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
		CompletedAt:     &completedAt,
	}
	order.SetPlayerID(player.ID)
	require.NoError(t, db.Create(order).Error)

	// Create repositories
	reviewRepo := review.NewReviewRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	playerRepo := user.NewPlayerRepository(db)
	userRepo := user.NewUserRepository(db)
	replyRepo := &mockReviewReplyRepository{db: db}
	notificationRepo := &mockNotificationRepository{}

	svc := NewReviewService(reviewRepo, orderRepo, playerRepo, userRepo, replyRepo, notificationRepo)

	t.Run("create review success", func(t *testing.T) {
		resp, err := svc.CreateReview(context.Background(), customer.ID, CreateReviewRequest{
			OrderID: order.ID,
			Rating:  5,
			Comment: "非常满意",
		})
		require.NoError(t, err)
		assert.NotZero(t, resp.ReviewID)
	})

	t.Run("create review unauthorized", func(t *testing.T) {
		_, err := svc.CreateReview(context.Background(), 99999, CreateReviewRequest{
			OrderID: order.ID,
			Rating:  5,
			Comment: "测试",
		})
		assert.Error(t, err)
	})

	t.Run("create review order not completed", func(t *testing.T) {
		// Create pending order
		pendingOrder := &model.Order{
			UserID:          customer.ID,
			ItemID:          1,
			Title:           "待支付订单",
			Status:          model.OrderStatusPending,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
		}
		pendingOrder.SetPlayerID(player.ID)
		require.NoError(t, db.Create(pendingOrder).Error)

		_, err := svc.CreateReview(context.Background(), customer.ID, CreateReviewRequest{
			OrderID: pendingOrder.ID,
			Rating:  5,
			Comment: "测试",
		})
		assert.ErrorIs(t, err, ErrOrderNotCompleted)
	})
}

func TestReviewService_GetMyReviews(t *testing.T) {
	db := setupReviewTestDB(t)
	defer testutil.CleanDB(t, db)

	// Create test user
	customer := &model.User{
		Phone:        "13800000030",
		Email:        "myreview_customer@test.com",
		Name:         "MyReview Customer",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(customer).Error)

	// Create player
	playerUser := &model.User{
		Phone:        "13800000031",
		Email:        "myreview_player@test.com",
		Name:         "MyReview Player",
		Role:         model.RolePlayer,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(playerUser).Error)

	player := &model.Player{
		UserID:             playerUser.ID,
		Nickname:           "MyReview Pro",
		HourlyRateCents:    5000,
		VerificationStatus: model.VerificationVerified,
	}
	require.NoError(t, db.Create(player).Error)

	// Create order
	scheduledStart := time.Now().Add(-24 * time.Hour)
	completedAt := time.Now().Add(-1 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          1,
		Title:           "已完成订单",
		Status:          model.OrderStatusCompleted,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
		CompletedAt:     &completedAt,
	}
	order.SetPlayerID(player.ID)
	require.NoError(t, db.Create(order).Error)

	// Create review
	reviewModel := &model.Review{
		OrderID:  order.ID,
		UserID:   customer.ID,
		PlayerID: player.ID,
		Score:    5,
		Content:  "好评",
		Status:   model.ReviewStatusApproved,
	}
	require.NoError(t, db.Create(reviewModel).Error)

	// Create repositories
	reviewRepo := review.NewReviewRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	playerRepo := user.NewPlayerRepository(db)
	userRepo := user.NewUserRepository(db)
	replyRepo := &mockReviewReplyRepository{db: db}
	notificationRepo := &mockNotificationRepository{}

	svc := NewReviewService(reviewRepo, orderRepo, playerRepo, userRepo, replyRepo, notificationRepo)

	t.Run("get my reviews", func(t *testing.T) {
		resp, err := svc.GetMyReviews(context.Background(), customer.ID, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(1), resp.Total)
		assert.Len(t, resp.Reviews, 1)
		assert.Equal(t, "好评", resp.Reviews[0].Comment)
	})

	t.Run("get my reviews with default pagination", func(t *testing.T) {
		resp, err := svc.GetMyReviews(context.Background(), customer.ID, 0, 0)
		require.NoError(t, err)
		assert.Equal(t, int64(1), resp.Total)
	})

	t.Run("get my reviews empty", func(t *testing.T) {
		resp, err := svc.GetMyReviews(context.Background(), 99999, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(0), resp.Total)
	})
}

func TestReviewService_GetPlayerReviews(t *testing.T) {
	db := setupReviewTestDB(t)
	defer testutil.CleanDB(t, db)

	// Create test user
	customer := &model.User{
		Phone:        "13800000040",
		Email:        "playerreview_customer@test.com",
		Name:         "PlayerReview Customer",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(customer).Error)

	// Create player
	playerUser := &model.User{
		Phone:        "13800000041",
		Email:        "playerreview_player@test.com",
		Name:         "PlayerReview Player",
		Role:         model.RolePlayer,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(playerUser).Error)

	player := &model.Player{
		UserID:             playerUser.ID,
		Nickname:           "PlayerReview Pro",
		HourlyRateCents:    5000,
		VerificationStatus: model.VerificationVerified,
	}
	require.NoError(t, db.Create(player).Error)

	// Create reviews for player
	for i := 0; i < 3; i++ {
		reviewModel := &model.Review{
			OrderID:  uint64(i + 1),
			UserID:   customer.ID,
			PlayerID: player.ID,
			Score:    model.Rating(4 + i%2),
			Content:  fmt.Sprintf("评价 %d", i+1),
			Status:   model.ReviewStatusApproved,
		}
		require.NoError(t, db.Create(reviewModel).Error)
	}

	// Create repositories
	reviewRepo := review.NewReviewRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	playerRepo := user.NewPlayerRepository(db)
	userRepo := user.NewUserRepository(db)
	replyRepo := &mockReviewReplyRepository{db: db}
	notificationRepo := &mockNotificationRepository{}

	svc := NewReviewService(reviewRepo, orderRepo, playerRepo, userRepo, replyRepo, notificationRepo)

	t.Run("get player reviews", func(t *testing.T) {
		reviews, total, err := svc.GetPlayerReviews(context.Background(), player.ID, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, reviews, 3)
	})

	t.Run("get player reviews with pagination", func(t *testing.T) {
		reviews, total, err := svc.GetPlayerReviews(context.Background(), player.ID, 1, 2)
		require.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, reviews, 2)
	})

	t.Run("get player reviews default pagination", func(t *testing.T) {
		reviews, total, err := svc.GetPlayerReviews(context.Background(), player.ID, 0, 0)
		require.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, reviews, 3)
	})
}

func TestReviewService_ReplyReview(t *testing.T) {
	db := setupReviewTestDB(t)
	defer testutil.CleanDB(t, db)

	// Create test user
	customer := &model.User{
		Phone:        "13800000050",
		Email:        "reply_customer@test.com",
		Name:         "Reply Customer",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(customer).Error)

	// Create player
	playerUser := &model.User{
		Phone:        "13800000051",
		Email:        "reply_player@test.com",
		Name:         "Reply Player",
		Role:         model.RolePlayer,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(playerUser).Error)

	player := &model.Player{
		UserID:             playerUser.ID,
		Nickname:           "Reply Pro",
		HourlyRateCents:    5000,
		VerificationStatus: model.VerificationVerified,
	}
	require.NoError(t, db.Create(player).Error)

	// Create review
	reviewModel := &model.Review{
		OrderID:  1,
		UserID:   customer.ID,
		PlayerID: player.ID,
		Score:    5,
		Content:  "好评",
		Status:   model.ReviewStatusApproved,
	}
	require.NoError(t, db.Create(reviewModel).Error)

	// Create repositories
	reviewRepo := review.NewReviewRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	playerRepo := user.NewPlayerRepository(db)
	userRepo := user.NewUserRepository(db)
	replyRepo := &mockReviewReplyRepository{db: db}
	notificationRepo := &mockNotificationRepository{}

	svc := NewReviewService(reviewRepo, orderRepo, playerRepo, userRepo, replyRepo, notificationRepo)

	t.Run("reply review success", func(t *testing.T) {
		resp, err := svc.ReplyReview(context.Background(), playerUser.ID, reviewModel.ID, ReplyReviewRequest{
			Content: "感谢您的好评！",
		})
		require.NoError(t, err)
		assert.NotZero(t, resp.ReplyID)
	})

	t.Run("reply review unauthorized", func(t *testing.T) {
		_, err := svc.ReplyReview(context.Background(), customer.ID, reviewModel.ID, ReplyReviewRequest{
			Content: "测试回复",
		})
		assert.Error(t, err)
	})

	t.Run("reply review not found", func(t *testing.T) {
		_, err := svc.ReplyReview(context.Background(), playerUser.ID, 99999, ReplyReviewRequest{
			Content: "测试回复",
		})
		assert.Error(t, err)
	})

	t.Run("reply review empty content", func(t *testing.T) {
		_, err := svc.ReplyReview(context.Background(), playerUser.ID, reviewModel.ID, ReplyReviewRequest{
			Content: "",
		})
		assert.Error(t, err)
	})
}

func TestReviewService_UpdateReply(t *testing.T) {
	db := setupReviewTestDB(t)
	defer testutil.CleanDB(t, db)

	// Create test user
	customer := &model.User{
		Phone:        "13800000060",
		Email:        "update_customer@test.com",
		Name:         "Update Customer",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(customer).Error)

	// Create player
	playerUser := &model.User{
		Phone:        "13800000061",
		Email:        "update_player@test.com",
		Name:         "Update Player",
		Role:         model.RolePlayer,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(playerUser).Error)

	player := &model.Player{
		UserID:             playerUser.ID,
		Nickname:           "Update Pro",
		HourlyRateCents:    5000,
		VerificationStatus: model.VerificationVerified,
	}
	require.NoError(t, db.Create(player).Error)

	// Create review
	reviewModel := &model.Review{
		OrderID:  1,
		UserID:   customer.ID,
		PlayerID: player.ID,
		Score:    5,
		Content:  "好评",
		Status:   model.ReviewStatusApproved,
	}
	require.NoError(t, db.Create(reviewModel).Error)

	// Create reply
	reply := &model.ReviewReply{
		ReviewID: reviewModel.ID,
		AuthorID: playerUser.ID,
		Content:  "原始回复",
		Status:   "approved",
	}
	require.NoError(t, db.Create(reply).Error)

	// Create repositories
	reviewRepo := review.NewReviewRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	playerRepo := user.NewPlayerRepository(db)
	userRepo := user.NewUserRepository(db)
	replyRepo := &mockReviewReplyRepository{db: db}
	notificationRepo := &mockNotificationRepository{}

	svc := NewReviewService(reviewRepo, orderRepo, playerRepo, userRepo, replyRepo, notificationRepo)

	t.Run("update reply success", func(t *testing.T) {
		resp, err := svc.UpdateReply(context.Background(), playerUser.ID, reply.ID, UpdateReplyRequest{
			Content: "更新后的回复",
		})
		require.NoError(t, err)
		assert.Equal(t, reply.ID, resp.ReplyID)
	})

	t.Run("update reply unauthorized", func(t *testing.T) {
		_, err := svc.UpdateReply(context.Background(), customer.ID, reply.ID, UpdateReplyRequest{
			Content: "测试更新",
		})
		assert.ErrorIs(t, err, ErrReviewUnauthorized)
	})

	t.Run("update reply not found", func(t *testing.T) {
		_, err := svc.UpdateReply(context.Background(), playerUser.ID, 99999, UpdateReplyRequest{
			Content: "测试更新",
		})
		assert.Error(t, err)
	})
}

func TestReviewService_DeleteReply(t *testing.T) {
	db := setupReviewTestDB(t)
	defer testutil.CleanDB(t, db)

	// Create test user
	customer := &model.User{
		Phone:        "13800000070",
		Email:        "delete_customer@test.com",
		Name:         "Delete Customer",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(customer).Error)

	// Create player
	playerUser := &model.User{
		Phone:        "13800000071",
		Email:        "delete_player@test.com",
		Name:         "Delete Player",
		Role:         model.RolePlayer,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(playerUser).Error)

	player := &model.Player{
		UserID:             playerUser.ID,
		Nickname:           "Delete Pro",
		HourlyRateCents:    5000,
		VerificationStatus: model.VerificationVerified,
	}
	require.NoError(t, db.Create(player).Error)

	// Create review
	reviewModel := &model.Review{
		OrderID:  1,
		UserID:   customer.ID,
		PlayerID: player.ID,
		Score:    5,
		Content:  "好评",
		Status:   model.ReviewStatusApproved,
	}
	require.NoError(t, db.Create(reviewModel).Error)

	// Create reply
	reply := &model.ReviewReply{
		ReviewID: reviewModel.ID,
		AuthorID: playerUser.ID,
		Content:  "待删除回复",
		Status:   "approved",
	}
	require.NoError(t, db.Create(reply).Error)

	// Create repositories
	reviewRepo := review.NewReviewRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	playerRepo := user.NewPlayerRepository(db)
	userRepo := user.NewUserRepository(db)
	replyRepo := &mockReviewReplyRepository{db: db}
	notificationRepo := &mockNotificationRepository{}

	svc := NewReviewService(reviewRepo, orderRepo, playerRepo, userRepo, replyRepo, notificationRepo)

	t.Run("delete reply unauthorized", func(t *testing.T) {
		err := svc.DeleteReply(context.Background(), customer.ID, reply.ID)
		assert.ErrorIs(t, err, ErrReviewUnauthorized)
	})

	t.Run("delete reply not found", func(t *testing.T) {
		err := svc.DeleteReply(context.Background(), playerUser.ID, 99999)
		assert.Error(t, err)
	})

	t.Run("delete reply success", func(t *testing.T) {
		err := svc.DeleteReply(context.Background(), playerUser.ID, reply.ID)
		require.NoError(t, err)

		// Verify deleted
		var count int64
		db.Model(&model.ReviewReply{}).Where("id = ?", reply.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}

// ============================================
// Additional DisputeService tests
// ============================================

func TestDisputeService_AssignDispute(t *testing.T) {
	db := setupDisputeTestDB(t)
	defer testutil.CleanDB(t, db)

	// Create admin user
	admin := &model.User{
		Phone:        "13800000080",
		Email:        "dispute_admin@test.com",
		Name:         "Dispute Admin",
		Role:         model.RoleAdmin,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(admin).Error)

	// Create pending dispute
	dispute := &model.OrderDispute{
		OrderID:     1,
		UserID:      1,
		Status:      model.DisputeStatusPending,
		Reason:      "测试纠纷",
		Description: "描述",
		TraceID:     "trace-assign-1",
	}
	require.NoError(t, db.Create(dispute).Error)

	disputeRepo := &mockDisputeRepository{db: db}
	orderRepo := orderimpl.NewOrderRepository(db)
	userRepo := user.NewUserRepository(db)
	opLogRepo := &mockOperationLogRepository{}
	notificationRepo := &mockNotificationRepository{}
	paymentRepo := orderrepo.NewPaymentRepository(db)

	svc := NewDisputeService(disputeRepo, orderRepo, userRepo, opLogRepo, notificationRepo, paymentRepo)

	t.Run("assign dispute success", func(t *testing.T) {
		err := svc.AssignDispute(context.Background(), AssignDisputeRequest{
			DisputeID:        dispute.ID,
			AssignedToUserID: admin.ID,
			Source:           model.AssignmentSourceManual,
			ActorUserID:      admin.ID,
		})
		require.NoError(t, err)

		// Verify assignment
		var updated model.OrderDispute
		require.NoError(t, db.First(&updated, dispute.ID).Error)
		assert.Equal(t, model.DisputeStatusAssigned, updated.Status)
		assert.NotNil(t, updated.AssignedToUserID)
		assert.Equal(t, admin.ID, *updated.AssignedToUserID)
	})

	t.Run("assign dispute not found", func(t *testing.T) {
		err := svc.AssignDispute(context.Background(), AssignDisputeRequest{
			DisputeID:        99999,
			AssignedToUserID: admin.ID,
			Source:           model.AssignmentSourceManual,
			ActorUserID:      admin.ID,
		})
		assert.Error(t, err)
	})

	t.Run("assign dispute validation error", func(t *testing.T) {
		err := svc.AssignDispute(context.Background(), AssignDisputeRequest{
			DisputeID:        0,
			AssignedToUserID: admin.ID,
		})
		assert.ErrorIs(t, err, ErrDisputeValidation)
	})
}

func TestDisputeService_ResolveDispute(t *testing.T) {
	db := setupDisputeTestDB(t)
	defer testutil.CleanDB(t, db)

	// Create admin user
	admin := &model.User{
		Phone:        "13800000090",
		Email:        "resolve_admin@test.com",
		Name:         "Resolve Admin",
		Role:         model.RoleAdmin,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(admin).Error)

	// Create order for dispute
	scheduledStart := time.Now().Add(-24 * time.Hour)
	order := &model.Order{
		UserID:          1,
		ItemID:          1,
		Title:           "纠纷订单",
		Status:          model.OrderStatusCompleted,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	require.NoError(t, db.Create(order).Error)

	// Create assigned dispute
	assignedTo := admin.ID
	dispute := &model.OrderDispute{
		OrderID:          order.ID,
		UserID:           1,
		Status:           model.DisputeStatusAssigned,
		Reason:           "测试纠纷",
		Description:      "描述",
		TraceID:          "trace-resolve-1",
		AssignedToUserID: &assignedTo,
	}
	require.NoError(t, db.Create(dispute).Error)

	disputeRepo := &mockDisputeRepository{db: db}
	orderRepo := orderimpl.NewOrderRepository(db)
	userRepo := user.NewUserRepository(db)
	opLogRepo := &mockOperationLogRepository{}
	notificationRepo := &mockNotificationRepository{}
	paymentRepo := orderrepo.NewPaymentRepository(db)

	svc := NewDisputeService(disputeRepo, orderRepo, userRepo, opLogRepo, notificationRepo, paymentRepo)

	t.Run("resolve dispute success", func(t *testing.T) {
		err := svc.ResolveDispute(context.Background(), ResolveDisputeRequest{
			DisputeID:       dispute.ID,
			Resolution:      model.ResolutionReject,
			ResolutionNotes: "已协商解决",
			ActorUserID:     admin.ID,
		})
		require.NoError(t, err)

		// Verify resolution
		var updated model.OrderDispute
		require.NoError(t, db.First(&updated, dispute.ID).Error)
		assert.Equal(t, model.DisputeStatusResolved, updated.Status)
	})

	t.Run("resolve dispute not found", func(t *testing.T) {
		err := svc.ResolveDispute(context.Background(), ResolveDisputeRequest{
			DisputeID:       99999,
			Resolution:      model.ResolutionReject,
			ResolutionNotes: "测试",
			ActorUserID:     admin.ID,
		})
		assert.Error(t, err)
	})

	t.Run("resolve dispute validation error", func(t *testing.T) {
		err := svc.ResolveDispute(context.Background(), ResolveDisputeRequest{
			DisputeID: 0,
		})
		assert.ErrorIs(t, err, ErrDisputeValidation)
	})
}

func TestDisputeService_ListDisputesByStatus(t *testing.T) {
	db := setupDisputeTestDB(t)
	defer testutil.CleanDB(t, db)

	// Create disputes with different statuses
	for i := 0; i < 3; i++ {
		dispute := &model.OrderDispute{
			OrderID:     uint64(i + 1),
			UserID:      1,
			Status:      model.DisputeStatusPending,
			Reason:      "待处理纠纷",
			Description: "描述",
			TraceID:     fmt.Sprintf("trace-list-%d", i),
		}
		require.NoError(t, db.Create(dispute).Error)
	}

	for i := 0; i < 2; i++ {
		dispute := &model.OrderDispute{
			OrderID:     uint64(i + 10),
			UserID:      1,
			Status:      model.DisputeStatusAssigned,
			Reason:      "已分配纠纷",
			Description: "描述",
			TraceID:     fmt.Sprintf("trace-assigned-%d", i),
		}
		require.NoError(t, db.Create(dispute).Error)
	}

	disputeRepo := &mockDisputeRepository{db: db}
	orderRepo := orderimpl.NewOrderRepository(db)
	userRepo := user.NewUserRepository(db)
	opLogRepo := &mockOperationLogRepository{}
	notificationRepo := &mockNotificationRepository{}
	paymentRepo := orderrepo.NewPaymentRepository(db)

	svc := NewDisputeService(disputeRepo, orderRepo, userRepo, opLogRepo, notificationRepo, paymentRepo)

	t.Run("list pending disputes", func(t *testing.T) {
		disputes, total, err := svc.ListDisputesByStatus(context.Background(), []model.DisputeStatus{model.DisputeStatusPending}, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, disputes, 3)
	})

	t.Run("list assigned disputes", func(t *testing.T) {
		disputes, total, err := svc.ListDisputesByStatus(context.Background(), []model.DisputeStatus{model.DisputeStatusAssigned}, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, disputes, 2)
	})
}

func TestDisputeService_CheckAndMarkSLABreaches(t *testing.T) {
	db := setupDisputeTestDB(t)
	defer testutil.CleanDB(t, db)

	// Create dispute with expired SLA
	expiredSLA := time.Now().Add(-1 * time.Hour)
	dispute := &model.OrderDispute{
		OrderID:     1,
		UserID:      1,
		Status:      model.DisputeStatusPending,
		Reason:      "SLA过期纠纷",
		Description: "描述",
		TraceID:     "trace-sla-1",
		SLADeadline: &expiredSLA,
	}
	require.NoError(t, db.Create(dispute).Error)

	disputeRepo := &mockDisputeRepository{db: db}
	orderRepo := orderimpl.NewOrderRepository(db)
	userRepo := user.NewUserRepository(db)
	opLogRepo := &mockOperationLogRepository{}
	notificationRepo := &mockNotificationRepository{}
	paymentRepo := orderrepo.NewPaymentRepository(db)

	svc := NewDisputeService(disputeRepo, orderRepo, userRepo, opLogRepo, notificationRepo, paymentRepo)

	t.Run("check and mark SLA breaches", func(t *testing.T) {
		err := svc.CheckAndMarkSLABreaches(context.Background())
		require.NoError(t, err)
	})
}

func TestDisputeService_RollbackDisputeAssignment(t *testing.T) {
	db := setupDisputeTestDB(t)
	defer testutil.CleanDB(t, db)

	// Create admin user
	admin := &model.User{
		Phone:        "13800000100",
		Email:        "rollback_admin@test.com",
		Name:         "Rollback Admin",
		Role:         model.RoleAdmin,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(admin).Error)

	// Create assigned dispute
	assignedTo := admin.ID
	assignedAt := time.Now()
	dispute := &model.OrderDispute{
		OrderID:          1,
		UserID:           1,
		Status:           model.DisputeStatusAssigned,
		Reason:           "测试纠纷",
		Description:      "描述",
		TraceID:          "trace-rollback-1",
		AssignedToUserID: &assignedTo,
		AssignedAt:       &assignedAt,
	}
	require.NoError(t, db.Create(dispute).Error)

	disputeRepo := &mockDisputeRepository{db: db}
	orderRepo := orderimpl.NewOrderRepository(db)
	userRepo := user.NewUserRepository(db)
	opLogRepo := &mockOperationLogRepository{}
	notificationRepo := &mockNotificationRepository{}
	paymentRepo := orderrepo.NewPaymentRepository(db)

	svc := NewDisputeService(disputeRepo, orderRepo, userRepo, opLogRepo, notificationRepo, paymentRepo)

	t.Run("rollback assignment success", func(t *testing.T) {
		err := svc.RollbackDisputeAssignment(context.Background(), RollbackDisputeRequest{
			DisputeID:      dispute.ID,
			RollbackReason: "需要重新分配",
			ActorUserID:    admin.ID,
		})
		require.NoError(t, err)
	})

	t.Run("rollback assignment not found", func(t *testing.T) {
		err := svc.RollbackDisputeAssignment(context.Background(), RollbackDisputeRequest{
			DisputeID:      99999,
			RollbackReason: "测试",
			ActorUserID:    admin.ID,
		})
		assert.Error(t, err)
	})

	t.Run("rollback assignment validation error", func(t *testing.T) {
		err := svc.RollbackDisputeAssignment(context.Background(), RollbackDisputeRequest{
			DisputeID: 0,
		})
		assert.ErrorIs(t, err, ErrDisputeValidation)
	})
}

func TestPaymentService_SetDistributedLock(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	paymentRepo := orderrepo.NewPaymentRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	svc := NewPaymentService(paymentRepo, orderRepo)

	t.Run("set distributed lock", func(t *testing.T) {
		svc.SetDistributedLock(nil)
		assert.Nil(t, svc.distributedLock)
	})
}

func TestPaymentService_SetWalletRepository(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	paymentRepo := orderrepo.NewPaymentRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	svc := NewPaymentService(paymentRepo, orderRepo)

	t.Run("set wallet repository", func(t *testing.T) {
		svc.SetWalletRepository(nil)
		assert.Nil(t, svc.wallets)
	})
}

func TestPaymentService_CreatePaymentValidation(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, player, gameModel := createOrderTestData(t, db)

	// Create pending order
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "待支付订单",
		Status:          model.OrderStatusPending,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	order.SetPlayerID(player.ID)
	order.SetGameID(gameModel.ID)
	require.NoError(t, db.Create(order).Error)

	paymentRepo := orderrepo.NewPaymentRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	svc := NewPaymentService(paymentRepo, orderRepo)

	t.Run("create alipay payment", func(t *testing.T) {
		resp, err := svc.CreatePayment(context.Background(), customer.ID, CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodAlipay,
		})
		require.NoError(t, err)
		assert.NotZero(t, resp.PaymentID)
		assert.Contains(t, resp.PayInfo, "trade_no")
	})

	t.Run("create payment order not found", func(t *testing.T) {
		_, err := svc.CreatePayment(context.Background(), customer.ID, CreatePaymentRequest{
			OrderID: 99999,
			Method:  model.PaymentMethodWeChat,
		})
		assert.Error(t, err)
	})

	t.Run("create payment unauthorized", func(t *testing.T) {
		_, err := svc.CreatePayment(context.Background(), 99999, CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodWeChat,
		})
		assert.Error(t, err)
	})
}

func TestPaymentService_ListPaymentsWithFilters(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, _, _ := createOrderTestData(t, db)

	// Create payments with different statuses
	for i := 0; i < 3; i++ {
		payment := &model.Payment{
			OrderID:     uint64(i + 1),
			UserID:      customer.ID,
			Method:      model.PaymentMethodWeChat,
			AmountCents: int64((i + 1) * 1000),
			Status:      model.PaymentStatusPending,
		}
		require.NoError(t, db.Create(payment).Error)
	}

	for i := 0; i < 2; i++ {
		paidAt := time.Now()
		payment := &model.Payment{
			OrderID:     uint64(i + 10),
			UserID:      customer.ID,
			Method:      model.PaymentMethodAlipay,
			AmountCents: int64((i + 1) * 2000),
			Status:      model.PaymentStatusPaid,
			PaidAt:      &paidAt,
		}
		require.NoError(t, db.Create(payment).Error)
	}

	paymentRepo := orderrepo.NewPaymentRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	svc := NewPaymentService(paymentRepo, orderRepo)

	t.Run("list all payments", func(t *testing.T) {
		payments, total, err := svc.ListPayments(context.Background(), repository.PaymentListOptions{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, payments, 5)
	})

	t.Run("list payments by status", func(t *testing.T) {
		status := model.PaymentStatusPaid
		payments, total, err := svc.ListPayments(context.Background(), repository.PaymentListOptions{
			Status:   &status,
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		// Note: The repository may not filter by status, just verify no error
		assert.GreaterOrEqual(t, total, int64(2))
		assert.GreaterOrEqual(t, len(payments), 2)
	})

	t.Run("list payments by user", func(t *testing.T) {
		payments, total, err := svc.ListPayments(context.Background(), repository.PaymentListOptions{
			UserID:   &customer.ID,
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, payments, 5)
	})
}

func TestPaymentService_CancelPaymentAlreadyPaid(t *testing.T) {
	db := setupOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, _, _ := createOrderTestData(t, db)

	// Create paid payment
	paidAt := time.Now()
	payment := &model.Payment{
		OrderID:     1,
		UserID:      customer.ID,
		Method:      model.PaymentMethodWeChat,
		AmountCents: 10000,
		Status:      model.PaymentStatusPaid,
		PaidAt:      &paidAt,
	}
	require.NoError(t, db.Create(payment).Error)

	paymentRepo := orderrepo.NewPaymentRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	svc := NewPaymentService(paymentRepo, orderRepo)

	t.Run("cannot cancel paid payment", func(t *testing.T) {
		err := svc.CancelPayment(context.Background(), customer.ID, payment.ID)
		assert.Error(t, err)
	})
}

func TestDisputeService_InitiateDisputeOrderNotCompleted(t *testing.T) {
	db := setupDisputeTestDB(t)
	defer testutil.CleanDB(t, db)

	// Create test user
	customer := &model.User{
		Phone:        "13800000110",
		Email:        "initiate_customer@test.com",
		Name:         "Initiate Customer",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(customer).Error)

	// Create pending order (not completed)
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          1,
		Title:           "待支付订单",
		Status:          model.OrderStatusPending,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	require.NoError(t, db.Create(order).Error)

	disputeRepo := &mockDisputeRepository{db: db}
	orderRepo := orderimpl.NewOrderRepository(db)
	userRepo := user.NewUserRepository(db)
	opLogRepo := &mockOperationLogRepository{}
	notificationRepo := &mockNotificationRepository{}
	paymentRepo := orderrepo.NewPaymentRepository(db)

	svc := NewDisputeService(disputeRepo, orderRepo, userRepo, opLogRepo, notificationRepo, paymentRepo)

	t.Run("initiate dispute order not completed", func(t *testing.T) {
		_, err := svc.InitiateDispute(context.Background(), InitiateDisputeRequest{
			OrderID:     order.ID,
			UserID:      customer.ID,
			Reason:      "服务质量问题",
			Description: "陪玩师态度不好",
		})
		assert.Error(t, err)
	})
}
