package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
	commissionrepo "gamelink/internal/repository/commission"
	gamerepo "gamelink/internal/repository/game"
	orderrepo "gamelink/internal/repository/implementations"
	paymentrepo "gamelink/internal/repository/payment"
	reviewrepo "gamelink/internal/repository/review"
	playerrepo "gamelink/internal/repository/user"
	orderservice "gamelink/internal/service/order"
	"gamelink/pkg/testutil"
)

// setupPlayerOrderTestDB 设置陪玩师订单测试数据库
func setupPlayerOrderTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.User{},
		&model.Game{},
		&model.Player{},
		&model.Order{},
		&model.ServiceItem{},
		&model.Payment{},
		&model.Wallet{},
		&model.Review{},
		&model.CommissionRule{},
		&model.CommissionRecord{},
	)
	return db
}

// setupOrderService 创建订单服务
func setupOrderService(db *gorm.DB) *orderservice.OrderService {
	orderRepo := orderrepo.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	userRepo := playerrepo.NewUserRepository(db)
	gameRepo := gamerepo.NewGameRepository(db)
	paymentRepo := paymentrepo.NewPaymentRepository(db)
	reviewRepo := reviewrepo.NewReviewRepository(db)
	commissionRepo := commissionrepo.NewCommissionRepository(db)

	return orderservice.NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)
}

// TestPlayerAcceptOrder 测试陪玩师接单
func TestPlayerAcceptOrder(t *testing.T) {
	db := setupPlayerOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	// 创建测试数据
	customer := &model.User{
		Phone:        "13900000001",
		Email:        "customer@test.com",
		Name:         "Customer",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(customer).Error)

	playerUser := &model.User{
		Phone:        "13900000002",
		Email:        "player@test.com",
		Name:         "Player",
		Role:         model.RolePlayer,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(playerUser).Error)

	game := &model.Game{Key: "lol", Name: "英雄联盟", Category: "moba"}
	require.NoError(t, db.Create(game).Error)

	serviceItem := &model.ServiceItem{
		ItemCode:       "escort-lol",
		Name:           "LOL陪玩",
		Category:       "escort",
		BasePriceCents: 5000,
		IsActive:       true,
		Tags:           "[]",
	}
	require.NoError(t, db.Create(serviceItem).Error)

	player := &model.Player{
		UserID:             playerUser.ID,
		Nickname:           "测试陪玩",
		MainGameID:         game.ID,
		HourlyRateCents:    5000,
		VerificationStatus: model.VerificationVerified,
	}
	require.NoError(t, db.Create(player).Error)

	// 创建待接单订单（已支付，等待陪玩师接单）
	order := &model.Order{
		OrderNo:         model.GenerateEscortOrderNo(),
		UserID:          customer.ID,
		ItemID:          serviceItem.ID,
		GameID:          &game.ID,
		Status:          model.OrderStatusConfirmed, // 已支付状态才能被接单
		TotalPriceCents: 5000,
		UnitPriceCents:  5000,
		Currency:        model.CurrencyCNY,
		Quantity:        1,
		OrderConfig:     "{}",
	}
	require.NoError(t, db.Create(order).Error)

	orderSvc := setupOrderService(db)

	t.Run("陪玩师成功接单", func(t *testing.T) {
		err := orderSvc.AcceptOrder(context.Background(), playerUser.ID, order.ID)
		require.NoError(t, err)

		// 验证订单状态变为进行中
		var updatedOrder model.Order
		require.NoError(t, db.First(&updatedOrder, order.ID).Error)
		assert.Equal(t, model.OrderStatusInProgress, updatedOrder.Status)
		assert.NotNil(t, updatedOrder.PlayerID)
		assert.Equal(t, player.ID, *updatedOrder.PlayerID)
	})

	t.Run("重复接单应失败", func(t *testing.T) {
		err := orderSvc.AcceptOrder(context.Background(), playerUser.ID, order.ID)
		assert.Error(t, err)
	})

	t.Run("非陪玩师不能接单", func(t *testing.T) {
		// 创建新订单
		order2 := &model.Order{
			OrderNo:         model.GenerateEscortOrderNo(),
			UserID:          customer.ID,
			ItemID:          serviceItem.ID,
			GameID:          &game.ID,
			Status:          model.OrderStatusConfirmed, // 已支付状态
			TotalPriceCents: 5000,
			UnitPriceCents:  5000,
			Currency:        model.CurrencyCNY,
			Quantity:        1,
			OrderConfig:     "{}",
		}
		require.NoError(t, db.Create(order2).Error)

		// 普通用户尝试接单
		err := orderSvc.AcceptOrder(context.Background(), customer.ID, order2.ID)
		assert.Error(t, err)
	})
}

// TestPlayerCompleteOrder 测试陪玩师完成订单
func TestPlayerCompleteOrder(t *testing.T) {
	db := setupPlayerOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	// 创建测试数据
	customer := &model.User{
		Phone:        "13900000011",
		Email:        "customer@test.com",
		Name:         "Customer",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(customer).Error)

	playerUser := &model.User{
		Phone:        "13900000012",
		Email:        "player@test.com",
		Name:         "Player",
		Role:         model.RolePlayer,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(playerUser).Error)

	game := &model.Game{Key: "lol", Name: "英雄联盟", Category: "moba"}
	require.NoError(t, db.Create(game).Error)

	serviceItem := &model.ServiceItem{
		ItemCode:       "escort-lol",
		Name:           "LOL陪玩",
		Category:       "escort",
		BasePriceCents: 5000,
		IsActive:       true,
		Tags:           "[]",
	}
	require.NoError(t, db.Create(serviceItem).Error)

	player := &model.Player{
		UserID:             playerUser.ID,
		Nickname:           "测试陪玩",
		MainGameID:         game.ID,
		HourlyRateCents:    5000,
		VerificationStatus: model.VerificationVerified,
	}
	require.NoError(t, db.Create(player).Error)

	// 创建进行中的订单
	now := time.Now()
	order := &model.Order{
		OrderNo:         model.GenerateEscortOrderNo(),
		UserID:          customer.ID,
		PlayerID:        &player.ID,
		ItemID:          serviceItem.ID,
		GameID:          &game.ID,
		Status:          model.OrderStatusInProgress,
		TotalPriceCents: 5000,
		UnitPriceCents:  5000,
		Currency:        model.CurrencyCNY,
		Quantity:        1,
		StartedAt:       &now,
		OrderConfig:     "{}",
	}
	require.NoError(t, db.Create(order).Error)

	orderSvc := setupOrderService(db)

	t.Run("陪玩师成功完成订单", func(t *testing.T) {
		err := orderSvc.CompleteOrderByPlayer(context.Background(), playerUser.ID, order.ID)
		require.NoError(t, err)

		// 验证订单状态
		var updatedOrder model.Order
		require.NoError(t, db.First(&updatedOrder, order.ID).Error)
		assert.Equal(t, model.OrderStatusCompleted, updatedOrder.Status)
		assert.NotNil(t, updatedOrder.CompletedAt)
	})

	t.Run("非订单陪玩师不能完成订单", func(t *testing.T) {
		// 创建另一个陪玩师
		otherPlayerUser := &model.User{
			Phone:        "13900000013",
			Email:        "other@test.com",
			Name:         "Other Player",
			Role:         model.RolePlayer,
			Status:       model.UserStatusActive,
			PasswordHash: "hashed",
		}
		require.NoError(t, db.Create(otherPlayerUser).Error)

		otherPlayer := &model.Player{
			UserID:             otherPlayerUser.ID,
			Nickname:           "其他陪玩",
			MainGameID:         game.ID,
			VerificationStatus: model.VerificationVerified,
		}
		require.NoError(t, db.Create(otherPlayer).Error)

		// 创建新订单
		order2 := &model.Order{
			OrderNo:         model.GenerateEscortOrderNo(),
			UserID:          customer.ID,
			PlayerID:        &player.ID, // 分配给原陪玩师
			ItemID:          serviceItem.ID,
			GameID:          &game.ID,
			Status:          model.OrderStatusInProgress,
			TotalPriceCents: 5000,
			UnitPriceCents:  5000,
			Currency:        model.CurrencyCNY,
			Quantity:        1,
			StartedAt:       &now,
			OrderConfig:     "{}",
		}
		require.NoError(t, db.Create(order2).Error)

		// 其他陪玩师尝试完成订单
		err := orderSvc.CompleteOrderByPlayer(context.Background(), otherPlayerUser.ID, order2.ID)
		assert.Error(t, err)
	})
}

// TestPlayerGetAvailableOrders 测试获取可接订单列表
func TestPlayerGetAvailableOrders(t *testing.T) {
	db := setupPlayerOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	// 创建测试数据
	customer := &model.User{
		Phone:        "13900000021",
		Email:        "customer@test.com",
		Name:         "Customer",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(customer).Error)

	game := &model.Game{Key: "lol", Name: "英雄联盟", Category: "moba"}
	require.NoError(t, db.Create(game).Error)

	serviceItem := &model.ServiceItem{
		ItemCode:       "escort-lol",
		Name:           "LOL陪玩",
		Category:       "escort",
		BasePriceCents: 5000,
		IsActive:       true,
		Tags:           "[]",
	}
	require.NoError(t, db.Create(serviceItem).Error)

	// 创建多个订单
	// GetAvailableOrders 查询 confirmed 状态的订单（已支付待接单）
	orders := []model.OrderStatus{
		model.OrderStatusConfirmed,  // 可接（已支付待接单）
		model.OrderStatusConfirmed,  // 可接（已支付待接单）
		model.OrderStatusPending,    // 未支付
		model.OrderStatusInProgress, // 进行中
		model.OrderStatusCompleted,  // 已完成
	}

	for i, status := range orders {
		order := &model.Order{
			OrderNo:         model.GenerateEscortOrderNo(),
			UserID:          customer.ID,
			ItemID:          serviceItem.ID,
			GameID:          &game.ID,
			Status:          status,
			TotalPriceCents: int64(5000 + i*1000),
			UnitPriceCents:  int64(5000 + i*1000),
			Currency:        model.CurrencyCNY,
			Quantity:        1,
			OrderConfig:     "{}",
		}
		require.NoError(t, db.Create(order).Error)
	}

	orderSvc := setupOrderService(db)

	t.Run("获取可接订单列表", func(t *testing.T) {
		orderList, total, err := orderSvc.GetAvailableOrders(context.Background(), orderservice.AvailableOrdersRequest{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		// 只有2个confirmed状态的订单（已支付待接单）
		assert.Equal(t, int64(2), total)
		assert.Equal(t, 2, len(orderList))
	})

	t.Run("按游戏筛选可接订单", func(t *testing.T) {
		orderList, total, err := orderSvc.GetAvailableOrders(context.Background(), orderservice.AvailableOrdersRequest{
			GameID:   &game.ID,
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Equal(t, 2, len(orderList))
	})
}

// TestPlayerGetMyOrders 测试获取用户订单列表
// 注意：GetMyOrders 是用户端方法，按 UserID 查询用户创建的订单
func TestPlayerGetMyOrders(t *testing.T) {
	db := setupPlayerOrderTestDB(t)
	defer testutil.CleanDB(t, db)

	// 创建测试数据
	customer := &model.User{
		Phone:        "13900000031",
		Email:        "customer@test.com",
		Name:         "Customer",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(customer).Error)

	playerUser := &model.User{
		Phone:        "13900000032",
		Email:        "player@test.com",
		Name:         "Player",
		Role:         model.RolePlayer,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(playerUser).Error)

	game := &model.Game{Key: "lol", Name: "英雄联盟", Category: "moba"}
	require.NoError(t, db.Create(game).Error)

	serviceItem := &model.ServiceItem{
		ItemCode:       "escort-lol",
		Name:           "LOL陪玩",
		Category:       "escort",
		BasePriceCents: 5000,
		IsActive:       true,
		Tags:           "[]",
	}
	require.NoError(t, db.Create(serviceItem).Error)

	player := &model.Player{
		UserID:             playerUser.ID,
		Nickname:           "测试陪玩",
		MainGameID:         game.ID,
		HourlyRateCents:    5000,
		VerificationStatus: model.VerificationVerified,
	}
	require.NoError(t, db.Create(player).Error)

	// 创建用户的订单（按 UserID 查询）
	statuses := []model.OrderStatus{
		model.OrderStatusConfirmed,
		model.OrderStatusInProgress,
		model.OrderStatusCompleted,
	}

	for _, status := range statuses {
		order := &model.Order{
			OrderNo:         model.GenerateEscortOrderNo(),
			UserID:          customer.ID, // 订单属于 customer
			PlayerID:        &player.ID,
			ItemID:          serviceItem.ID,
			GameID:          &game.ID,
			Status:          status,
			TotalPriceCents: 5000,
			UnitPriceCents:  5000,
			Currency:        model.CurrencyCNY,
			Quantity:        1,
			OrderConfig:     "{}",
		}
		require.NoError(t, db.Create(order).Error)
	}

	orderSvc := setupOrderService(db)

	t.Run("获取用户所有订单", func(t *testing.T) {
		// 使用 customer.ID 查询，因为订单是 customer 创建的
		resp, err := orderSvc.GetMyOrders(context.Background(), customer.ID, orderservice.MyOrderListRequest{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, 3, len(resp.Orders))
	})

	t.Run("按状态筛选订单", func(t *testing.T) {
		resp, err := orderSvc.GetMyOrders(context.Background(), customer.ID, orderservice.MyOrderListRequest{
			Status:   "completed",
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(resp.Orders))
		assert.Equal(t, model.OrderStatusCompleted, resp.Orders[0].Status)
	})

	t.Run("陪玩师查询自己创建的订单应为空", func(t *testing.T) {
		// playerUser 没有创建任何订单，所以应该返回空
		resp, err := orderSvc.GetMyOrders(context.Background(), playerUser.ID, orderservice.MyOrderListRequest{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, 0, len(resp.Orders))
	})
}
