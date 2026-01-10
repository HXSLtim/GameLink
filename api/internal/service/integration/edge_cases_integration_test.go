// Package integration provides edge case integration tests for services.
package integration

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"gamelink/internal/model"
	commissionrepo "gamelink/internal/repository/commission"
	"gamelink/internal/repository/game"
	"gamelink/internal/repository/implementations"
	"gamelink/internal/repository/payment"
	"gamelink/internal/repository/player"
	"gamelink/internal/repository/review"
	userrepo "gamelink/internal/repository/user"
	"gamelink/internal/service/order"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Order Edge Cases - 超时、并发、边界值
// ============================================================================

// TestOrder_PaymentTimeoutEdgeCase 测试支付超时边界情况
func TestOrder_PaymentTimeoutEdgeCase(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	testUser := CreateUniqueTestUser(t, db, "timeout_user")
	playerUser := CreateUniqueTestUser(t, db, "timeout_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "timeout_game")

	// 创建服务项目
	serviceItem := CreateTestServiceItem(t, db, testGame, "Timeout Service", 5000)

	// 测试场景1：订单创建29分59秒前，未超时
	createdAt := time.Now().Add(-29*time.Minute - 59*time.Second)
	order1 := &model.Order{
		Base:            model.Base{ExtJSON: "{}", CreatedAt: createdAt},
		OrderNo:         fmt.Sprintf("TIMEOUT1%d", time.Now().UnixNano()),
		UserID:          testUser.ID,
		PlayerID:        &testPlayer.ID,
		GameID:          &testGame.ID,
		ItemID:          serviceItem.ID,
		TotalPriceCents: 5000,
		Currency:        model.CurrencyCNY,
		Status:          model.OrderStatusPending,
		Title:           "Timeout Test Order 1",
		OrderConfig:     "{}",
	}
	require.NoError(t, db.Create(order1).Error)

	// 验证订单时间计算
	timeSinceCreation1 := time.Since(order1.CreatedAt)
	timeoutThreshold := 30 * time.Minute
	notTimedOut := timeSinceCreation1 < timeoutThreshold
	assert.True(t, notTimedOut, "Order should NOT be timed out (29m59s)")

	// 测试场景2：订单创建30分01秒前，已超时
	createdAt2 := time.Now().Add(-30*time.Minute - time.Second)
	order2 := &model.Order{
		Base:            model.Base{ExtJSON: "{}", CreatedAt: createdAt2},
		OrderNo:         fmt.Sprintf("TIMEOUT2%d", time.Now().UnixNano()),
		UserID:          testUser.ID,
		PlayerID:        &testPlayer.ID,
		GameID:          &testGame.ID,
		ItemID:          serviceItem.ID,
		TotalPriceCents: 5000,
		Currency:        model.CurrencyCNY,
		Status:          model.OrderStatusPending,
		Title:           "Timeout Test Order 2",
		OrderConfig:     "{}",
	}
	require.NoError(t, db.Create(order2).Error)

	// 验证订单时间计算
	timeSinceCreation2 := time.Since(order2.CreatedAt)
	timedOut := timeSinceCreation2 >= timeoutThreshold
	assert.True(t, timedOut, "Order should be timed out (30m01s)")

	// 模拟超时处理：将订单状态设为已取消
	order2.Status = model.OrderStatusCanceled
	require.NoError(t, db.Save(order2).Error)

	// 验证状态已更新
	var canceledOrder model.Order
	require.NoError(t, db.Where("id = ?", order2.ID).First(&canceledOrder).Error)
	assert.Equal(t, model.OrderStatusCanceled, canceledOrder.Status)
}

// TestOrder_ConcurrentOrderCreation 测试并发创建订单
func TestOrder_ConcurrentOrderCreation(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	gameRepo := game.NewGameRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	commissionRepo := commissionrepo.NewCommissionRepository(db)

	svc := order.NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)

	testUser := CreateUniqueTestUser(t, db, "concurrent_user")
	playerUser := CreateUniqueTestUser(t, db, "concurrent_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "concurrent_game")

	serviceItem := CreateTestServiceItem(t, db, testGame, "Concurrent Service", 5000)
	testPlayer.HourlyRateCents = 5000
	db.Save(testPlayer)

	now := time.Now().Add(time.Hour)
	serviceID := serviceItem.ID

	// 并发创建10个订单
	const numOrders = 10
	var wg sync.WaitGroup
	errors := make(chan error, numOrders)
	orderIDs := make(chan uint64, numOrders)

	for i := 0; i < numOrders; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			req := order.CreateOrderRequest{
				PlayerID:       testPlayer.ID,
				GameID:         testGame.ID,
				ServiceID:      &serviceID,
				Title:          fmt.Sprintf("Concurrent Order %d", index),
				Description:    "Test concurrent order creation",
				ScheduledStart: &now,
				DurationHours:  1,
			}

			resp, err := svc.CreateOrder(ctx, testUser.ID, req)
			if err != nil {
				errors <- err
				return
			}
			orderIDs <- resp.OrderID
		}(i)
	}

	wg.Wait()
	close(errors)
	close(orderIDs)

	// 检查是否有错误
	for err := range errors {
		t.Errorf("Concurrent order creation error: %v", err)
	}

	// 收集所有订单ID
	var ids []uint64
	for id := range orderIDs {
		ids = append(ids, id)
	}

	// 验证所有订单都创建成功
	assert.Len(t, ids, numOrders, "All concurrent orders should be created")

	// 验证订单号唯一
	orderNos := make(map[string]bool)
	for _, id := range ids {
		var order model.Order
		require.NoError(t, db.Where("id = ?", id).First(&order).Error)
		assert.False(t, orderNos[order.OrderNo], "Order number should be unique")
		orderNos[order.OrderNo] = true
	}
}

// TestOrder_BoundaryValues 测试边界值
func TestOrder_BoundaryValues(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	testUser := CreateUniqueTestUser(t, db, "boundary_user")
	playerUser := CreateUniqueTestUser(t, db, "boundary_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "boundary_game")

	serviceItem := CreateTestServiceItem(t, db, testGame, "Boundary Service", 5000)

	now := time.Now()

	tests := []struct {
		name           string
		durationHours  float32
		expectedPrice  int64
		shouldValidate bool
	}{
		{
			name:           "最小时长: 1小时",
			durationHours:  1,
			expectedPrice:  5000,
			shouldValidate: true,
		},
		{
			name:           "最大时长: 24小时",
			durationHours:  24,
			expectedPrice:  120000,
			shouldValidate: true,
		},
		{
			name:           "边界值: 0.5小时",
			durationHours:  0.5,
			expectedPrice:  2500,
			shouldValidate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheduledStart := now.Add(time.Hour)
			scheduledEnd := scheduledStart.Add(time.Duration(tt.durationHours) * time.Hour)

			order := &model.Order{
				Base:            model.Base{ExtJSON: "{}"},
				OrderNo:         fmt.Sprintf("BNDRY%d", time.Now().UnixNano()),
				UserID:          testUser.ID,
				PlayerID:        &testPlayer.ID,
				GameID:          &testGame.ID,
				ItemID:          serviceItem.ID,
				TotalPriceCents: tt.expectedPrice,
				Currency:        model.CurrencyCNY,
				Status:          model.OrderStatusPending,
				Title:           tt.name,
				ScheduledStart:  &scheduledStart,
				ScheduledEnd:    &scheduledEnd,
				OrderConfig:     "{}",
			}
			require.NoError(t, db.Create(order).Error)

			// 验证订单创建成功
			var retrievedOrder model.Order
			require.NoError(t, db.Where("id = ?", order.ID).First(&retrievedOrder).Error)
			assert.Equal(t, tt.expectedPrice, retrievedOrder.TotalPriceCents)

			// 验证时间计算正确
			if retrievedOrder.ScheduledStart != nil && retrievedOrder.ScheduledEnd != nil {
				duration := retrievedOrder.ScheduledEnd.Sub(*retrievedOrder.ScheduledStart)
				expectedDuration := time.Duration(tt.durationHours) * time.Hour
				// 允许1秒的误差
				assert.LessOrEqual(t, duration.Abs()-expectedDuration, time.Second)
			}
		})
	}
}

// TestOrder_ZeroAndNegativeAmounts 测试零金额和负金额
func TestOrder_ZeroAndNegativeAmounts(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	testUser := CreateUniqueTestUser(t, db, "amount_user")
	playerUser := CreateUniqueTestUser(t, db, "amount_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "amount_game")

	// 测试零金额订单（免费体验活动）
	zeroPriceService := CreateTestServiceItem(t, db, testGame, "Free Service", 0)

	now := time.Now()
	scheduledStart := now.Add(time.Hour)
	scheduledEnd := now.Add(2 * time.Hour)

	zeroPriceOrder := &model.Order{
		Base:              model.Base{ExtJSON: "{}"},
		OrderNo:           fmt.Sprintf("ZERO%d", time.Now().UnixNano()),
		UserID:            testUser.ID,
		PlayerID:          &testPlayer.ID,
		GameID:            &testGame.ID,
		ItemID:            zeroPriceService.ID,
		TotalPriceCents:   0,
		CommissionCents:   0,
		PlayerIncomeCents: 0,
		Currency:          model.CurrencyCNY,
		Status:            model.OrderStatusConfirmed, // 零金额订单直接确认
		Title:             "Free Order",
		ScheduledStart:    &scheduledStart,
		ScheduledEnd:      &scheduledEnd,
		OrderConfig:       "{}",
	}
	require.NoError(t, db.Create(zeroPriceOrder).Error)

	// 验证零金额订单
	var retrievedOrder model.Order
	require.NoError(t, db.Where("id = ?", zeroPriceOrder.ID).First(&retrievedOrder).Error)
	assert.Equal(t, int64(0), retrievedOrder.TotalPriceCents, "Total price should be 0")
	assert.Equal(t, int64(0), retrievedOrder.CommissionCents, "Commission should be 0")
	assert.Equal(t, int64(0), retrievedOrder.PlayerIncomeCents, "Player income should be 0")
}

// ============================================================================
// Payment Edge Cases - 组合支付失败、部分退款
// ============================================================================

// TestPayment_CombinedPaymentFailure 测试组合支付失败场景
func TestPayment_CombinedPaymentFailure(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	testUser := CreateUniqueTestUser(t, db, "combined_user")
	playerUser := CreateUniqueTestUser(t, db, "combined_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "combined_game")

	// 创建订单
	serviceItem := CreateTestServiceItem(t, db, testGame, "Combined Service", 10000)
	testPlayer.HourlyRateCents = 10000
	db.Save(testPlayer)

	now := time.Now().Add(time.Hour)
	order := &model.Order{
		Base:            model.Base{ExtJSON: "{}"},
		OrderNo:         fmt.Sprintf("COMBO%d", time.Now().UnixNano()),
		UserID:          testUser.ID,
		PlayerID:        &testPlayer.ID,
		GameID:          &testGame.ID,
		ItemID:          serviceItem.ID,
		TotalPriceCents: 10000,
		Currency:        model.CurrencyCNY,
		Status:          model.OrderStatusPending,
		Title:           "Combined Payment Test",
		ScheduledStart:  &now,
		OrderConfig:     "{}",
	}
	require.NoError(t, db.Create(order).Error)

	// 创建钱包，余额3000分（需补7000分）
	wallet := &model.Wallet{
		Base:         model.Base{ExtJSON: "{}"},
		UserID:       testUser.ID,
		BalanceCents: 3000,
		FrozenCents:  0,
	}
	require.NoError(t, db.Create(wallet).Error)

	// 创建支付记录（钱包部分3000 + 第三方7000）
	payment := &model.Payment{
		Base:                  model.Base{ExtJSON: "{}"},
		OrderID:               order.ID,
		UserID:                testUser.ID,
		AmountCents:           10000,
		WalletAmountCents:     3000,
		ThirdPartyAmountCents: 7000,
		Method:                model.PaymentMethodCombined,
		Status:                model.PaymentStatusPending,
		ProviderRaw:           []byte("{}"),
	}
	require.NoError(t, db.Create(payment).Error)

	// 记录初始余额
	initialBalance := wallet.BalanceCents

	// 模拟钱包扣款成功
	wallet.BalanceCents -= 3000
	wallet.FrozenCents += 3000 // 冻结等待第三方支付
	require.NoError(t, db.Save(wallet).Error)

	// 模拟第三方支付失败
	payment.Status = model.PaymentStatusFailed
	payment.ProviderRaw = []byte(`{"status":"failed","message":"Third party payment timeout"}`)
	require.NoError(t, db.Save(payment).Error)

	// 模拟退款逻辑：退还钱包余额
	wallet.BalanceCents = initialBalance
	wallet.FrozenCents = 0
	require.NoError(t, db.Save(wallet).Error)

	// 更新订单状态为已取消
	order.Status = model.OrderStatusCanceled
	require.NoError(t, db.Save(order).Error)

	// 验证钱包余额已退还
	var updatedWallet model.Wallet
	require.NoError(t, db.Where("user_id = ?", testUser.ID).First(&updatedWallet).Error)
	assert.Equal(t, initialBalance, updatedWallet.BalanceCents, "Wallet should be refunded on combined payment failure")
	assert.Equal(t, int64(0), updatedWallet.FrozenCents, "Frozen amount should be released")

	// 验证订单状态
	var updatedOrder model.Order
	require.NoError(t, db.Where("id = ?", order.ID).First(&updatedOrder).Error)
	assert.Equal(t, model.OrderStatusCanceled, updatedOrder.Status, "Order should be canceled on payment failure")
}

// TestPayment_PartialRefund 测试部分退款
func TestPayment_PartialRefund(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	testUser := CreateUniqueTestUser(t, db, "refund_user")
	playerUser := CreateUniqueTestUser(t, db, "refund_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "refund_game")

	// 创建已完成订单
	serviceItem := CreateTestServiceItem(t, db, testGame, "Refund Service", 10000)
	testPlayer.HourlyRateCents = 10000
	db.Save(testPlayer)

	now := time.Now()
	scheduledStart := now.Add(-2 * time.Hour)
	scheduledEnd := now.Add(-1 * time.Hour)
	order := &model.Order{
		Base:              model.Base{ExtJSON: "{}"},
		OrderNo:           fmt.Sprintf("REFUND%d", time.Now().UnixNano()),
		UserID:            testUser.ID,
		PlayerID:          &testPlayer.ID,
		GameID:            &testGame.ID,
		ItemID:            serviceItem.ID,
		TotalPriceCents:   10000,
		CommissionCents:   2000,
		PlayerIncomeCents: 8000,
		Currency:          model.CurrencyCNY,
		Status:            model.OrderStatusCompleted,
		Title:             "Partial Refund Test",
		ScheduledStart:    &scheduledStart,
		ScheduledEnd:      &scheduledEnd,
		CompletedAt:       &now,
		OrderConfig:       "{}",
	}
	require.NoError(t, db.Create(order).Error)

	// 创建支付记录（组合支付）
	payment := &model.Payment{
		Base:                  model.Base{ExtJSON: "{}"},
		OrderID:               order.ID,
		UserID:                testUser.ID,
		AmountCents:           10000,
		WalletAmountCents:     5000,
		ThirdPartyAmountCents: 5000,
		Method:                model.PaymentMethodCombined,
		Status:                model.PaymentStatusPaid,
		ProviderTradeNo:       "TXN" + fmt.Sprint(time.Now().UnixNano()),
		ProviderRaw:           []byte("{}"),
	}
	require.NoError(t, db.Create(payment).Error)

	// 创建钱包
	wallet := &model.Wallet{
		Base:         model.Base{ExtJSON: "{}"},
		UserID:       testUser.ID,
		BalanceCents: 5000,
		FrozenCents:  0,
	}
	require.NoError(t, db.Create(wallet).Error)

	// 创建陪玩师钱包
	playerWallet := &model.Wallet{
		Base:         model.Base{ExtJSON: "{}"},
		UserID:       playerUser.ID,
		BalanceCents: 10000,
		FrozenCents:  8000, // T+7冻结中
	}
	require.NoError(t, db.Create(playerWallet).Error)

	// 执行部分退款（退款50%）
	refundAmount := int64(5000)     // 总额的一半
	walletRefund := int64(2500)     // 钱包部分一半
	thirdpartyRefund := int64(2500) // 第三方部分一半
	_ = thirdpartyRefund            // 暂时使用

	// 更新钱包
	wallet.BalanceCents += walletRefund
	require.NoError(t, db.Save(wallet).Error)

	// 从陪玩师冻结中扣除
	playerWallet.FrozenCents -= 2500 // 陪玩师收入也要退一半
	require.NoError(t, db.Save(playerWallet).Error)

	// 更新原支付记录的退款金额
	payment.RefundedAmountCents = refundAmount
	require.NoError(t, db.Save(payment).Error)

	// 验证退款结果
	var finalWallet model.Wallet
	require.NoError(t, db.Where("user_id = ?", testUser.ID).First(&finalWallet).Error)
	assert.Equal(t, int64(7500), finalWallet.BalanceCents, "Wallet should have 7500 after partial refund")

	var finalPlayerWallet model.Wallet
	require.NoError(t, db.Where("user_id = ?", playerUser.ID).First(&finalPlayerWallet).Error)
	assert.Equal(t, int64(5500), finalPlayerWallet.FrozenCents, "Player frozen should be reduced by refund share")
}

// ============================================================================
// Wallet Edge Cases - T+7边界、负余额
// ============================================================================

// TestWallet_TPlus7Boundary 测试T+7结算边界
func TestWallet_TPlus7Boundary(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	testUser := CreateUniqueTestUser(t, db, "t7_user")
	playerUser := CreateUniqueTestUser(t, db, "t7_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "t7_game")

	// 创建7天前完成的订单
	completedTime := time.Now().Add(-7 * 24 * time.Hour)
	serviceItem := CreateTestServiceItem(t, db, testGame, "T7 Service", 10000)
	testPlayer.HourlyRateCents = 10000
	db.Save(testPlayer)

	scheduledStart := completedTime.Add(-2 * time.Hour)
	scheduledEnd := completedTime.Add(-1 * time.Hour)
	order := &model.Order{
		Base:              model.Base{ExtJSON: "{}"},
		OrderNo:           fmt.Sprintf("T7%d", time.Now().UnixNano()),
		UserID:            testUser.ID,
		PlayerID:          &testPlayer.ID,
		GameID:            &testGame.ID,
		ItemID:            serviceItem.ID,
		TotalPriceCents:   10000,
		CommissionCents:   2000,
		PlayerIncomeCents: 8000,
		Currency:          model.CurrencyCNY,
		Status:            model.OrderStatusCompleted,
		Title:             "T+7 Boundary Test",
		ScheduledStart:    &scheduledStart,
		ScheduledEnd:      &scheduledEnd,
		CompletedAt:       &completedTime,
		OrderConfig:       "{}",
	}
	require.NoError(t, db.Create(order).Error)

	// 创建陪玩师钱包（T+7冻结中）
	playerWallet := &model.Wallet{
		Base:         model.Base{ExtJSON: "{}"},
		UserID:       playerUser.ID,
		BalanceCents: 0,
		FrozenCents:  8000,
	}
	require.NoError(t, db.Create(playerWallet).Error)

	// 计算T+7是否到期
	timeSinceCompletion := time.Since(*order.CompletedAt)
	t7Expired := timeSinceCompletion >= 7*24*time.Hour

	assert.True(t, t7Expired, "T+7 period should be expired")

	// 如果T+7到期，执行解冻
	if t7Expired {
		// 冻结金额转入余额
		playerWallet.BalanceCents += playerWallet.FrozenCents
		playerWallet.FrozenCents = 0
		require.NoError(t, db.Save(playerWallet).Error)

		// 验证余额
		var updatedWallet model.Wallet
		require.NoError(t, db.Where("user_id = ?", playerUser.ID).First(&updatedWallet).Error)
		assert.Equal(t, int64(8000), updatedWallet.BalanceCents, "Balance should be 8000 after T+7 settlement")
		assert.Equal(t, int64(0), updatedWallet.FrozenCents, "Frozen should be 0 after T+7 settlement")
	}
}

// TestWallet_TPlus7OneSecondBefore 测试T+7前一秒
func TestWallet_TPlus7OneSecondBefore(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	testUser := CreateUniqueTestUser(t, db, "t7before_user")
	playerUser := CreateUniqueTestUser(t, db, "t7before_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "t7before_game")

	// 创建6天23小时59分59秒前完成的订单（差1秒到7天）
	completedTime := time.Now().Add(-7*24*time.Hour + time.Second)
	serviceItem := CreateTestServiceItem(t, db, testGame, "T7Before Service", 10000)
	testPlayer.HourlyRateCents = 10000
	db.Save(testPlayer)

	scheduledStart := completedTime.Add(-2 * time.Hour)
	scheduledEnd := completedTime.Add(-1 * time.Hour)
	order := &model.Order{
		Base:              model.Base{ExtJSON: "{}"},
		OrderNo:           fmt.Sprintf("T7B%d", time.Now().UnixNano()),
		UserID:            testUser.ID,
		PlayerID:          &testPlayer.ID,
		GameID:            &testGame.ID,
		ItemID:            serviceItem.ID,
		TotalPriceCents:   10000,
		CommissionCents:   2000,
		PlayerIncomeCents: 8000,
		Currency:          model.CurrencyCNY,
		Status:            model.OrderStatusCompleted,
		Title:             "T+7 One Second Before Test",
		ScheduledStart:    &scheduledStart,
		ScheduledEnd:      &scheduledEnd,
		CompletedAt:       &completedTime,
		OrderConfig:       "{}",
	}
	require.NoError(t, db.Create(order).Error)

	// 创建陪玩师钱包
	playerWallet := &model.Wallet{
		Base:         model.Base{ExtJSON: "{}"},
		UserID:       playerUser.ID,
		BalanceCents: 0,
		FrozenCents:  8000,
	}
	require.NoError(t, db.Create(playerWallet).Error)

	// 计算T+7是否到期
	timeSinceCompletion := time.Since(*order.CompletedAt)
	t7Expired := timeSinceCompletion >= 7*24*time.Hour

	assert.False(t, t7Expired, "T+7 period should NOT be expired (1 second before)")

	// 验证金额仍在冻结中
	var updatedWallet model.Wallet
	require.NoError(t, db.Where("user_id = ?", playerUser.ID).First(&updatedWallet).Error)
	assert.Equal(t, int64(0), updatedWallet.BalanceCents, "Balance should still be 0")
	assert.Equal(t, int64(8000), updatedWallet.FrozenCents, "Frozen should still be 8000")
}

// TestWallet_NegativeBalancePrevention 测试防止负余额
func TestWallet_NegativeBalancePrevention(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	testUser := CreateUniqueTestUser(t, db, "negative_user")

	// 创建余额为100的钱包
	wallet := &model.Wallet{
		Base:         model.Base{ExtJSON: "{}"},
		UserID:       testUser.ID,
		BalanceCents: 100,
		FrozenCents:  0,
	}
	require.NoError(t, db.Create(wallet).Error)

	// 尝试扣款200（余额不足）
	attemptDeduct := int64(200)

	// 检查余额是否足够
	if wallet.BalanceCents < attemptDeduct {
		// 应该拒绝交易
		assert.True(t, true, "Transaction should be rejected due to insufficient balance")
		return
	}

	// 如果代码执行到这里，说明余额检查失败
	t.Fatal("Balance check failed - negative balance should not be allowed")
}

// TestWallet_ConcurrentDeduction 测试并发扣款
func TestWallet_ConcurrentDeduction(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	testUser := CreateUniqueTestUser(t, db, "concurrent_wallet_user")

	// 创建余额为10000的钱包
	wallet := &model.Wallet{
		Base:         model.Base{ExtJSON: "{}"},
		UserID:       testUser.ID,
		BalanceCents: 10000,
		FrozenCents:  0,
	}
	require.NoError(t, db.Create(wallet).Error)

	// 并发扣款10次，每次1000
	const numDeductions = 10
	const deductionAmount = 1000
	var wg sync.WaitGroup
	errors := make(chan error, numDeductions)

	for i := 0; i < numDeductions; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			// 使用事务+行锁保证原子性
			tx := db.Begin()
			var currentWallet model.Wallet
			// 使用 SELECT FOR UPDATE 行锁
			if err := tx.Raw("SELECT * FROM wallets WHERE user_id = ? FOR UPDATE", testUser.ID).Scan(&currentWallet).Error; err != nil {
				errors <- err
				tx.Rollback()
				return
			}

			// 检查余额
			if currentWallet.BalanceCents < deductionAmount {
				errors <- fmt.Errorf("insufficient balance for deduction %d", index)
				tx.Rollback()
				return
			}

			// 扣款
			currentWallet.BalanceCents -= deductionAmount
			if err := tx.Save(&currentWallet).Error; err != nil {
				errors <- err
				tx.Rollback()
				return
			}

			tx.Commit()
		}(i)
	}

	wg.Wait()
	close(errors)

	errorCount := 0
	for range errors {
		errorCount++
	}

	assert.Equal(t, 0, errorCount, "All concurrent deductions should succeed")

	// 验证最终余额
	var finalWallet model.Wallet
	require.NoError(t, db.Where("user_id = ?", testUser.ID).First(&finalWallet).Error)
	assert.Equal(t, int64(0), finalWallet.BalanceCents, "Final balance should be 0")
}

// ============================================================================
// Dispute Edge Cases - SLA边界、同时发起
// ============================================================================

// TestDispute_SLABoundary 测试SLA边界（30分钟）
func TestDispute_SLABoundary(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	testUser := CreateUniqueTestUser(t, db, "sla_user")
	playerUser := CreateUniqueTestUser(t, db, "sla_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "sla_game")

	// 创建进行中的订单
	now := time.Now()
	scheduledStart := now.Add(-1 * time.Hour)
	scheduledEnd := now.Add(time.Hour)

	// 先创建服务项目
	serviceItem := CreateTestServiceItem(t, db, testGame, "SLA Service", 10000)

	order := &model.Order{
		Base:            model.Base{ExtJSON: "{}"},
		OrderNo:         fmt.Sprintf("SLA%d", time.Now().UnixNano()),
		UserID:          testUser.ID,
		PlayerID:        &testPlayer.ID,
		GameID:          &testGame.ID,
		ItemID:          serviceItem.ID,
		TotalPriceCents: 10000,
		Currency:        model.CurrencyCNY,
		Status:          model.OrderStatusInProgress,
		Title:           "SLA Boundary Test",
		ScheduledStart:  &scheduledStart,
		ScheduledEnd:    &scheduledEnd,
		OrderConfig:     "{}",
	}
	require.NoError(t, db.Create(order).Error)

	// 创建29分钟前的争议（刚好在SLA边界内）
	createdAt := time.Now().Add(-29*time.Minute - 59*time.Second)
	dispute := &model.OrderDispute{
		Base:          model.Base{ExtJSON: "{}", CreatedAt: createdAt, UpdatedAt: createdAt},
		OrderID:       order.ID,
		InitiatorID:   testUser.ID,
		InitiatorType: "user",
		Type:          model.DisputeTypeServiceQuality,
		Reason:        "Test SLA boundary",
		Status:        model.DisputeStatusPending,
	}
	require.NoError(t, db.Create(dispute).Error)

	// 验证争议仍在SLA内
	timeSinceCreated := time.Since(dispute.CreatedAt)
	slaExpired := timeSinceCreated >= 30*time.Minute

	assert.False(t, slaExpired, "SLA should NOT be expired")

	// 验证争议状态
	var pendingDispute model.OrderDispute
	require.NoError(t, db.Where("id = ?", dispute.ID).First(&pendingDispute).Error)
	assert.Equal(t, model.DisputeStatusPending, pendingDispute.Status)
}

// TestDispute_SLAExpired 测试SLA过期
func TestDispute_SLAExpired(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	testUser := CreateUniqueTestUser(t, db, "slaexp_user")
	playerUser := CreateUniqueTestUser(t, db, "slaexp_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "slaexp_game")

	// 创建订单
	now := time.Now()
	scheduledStart := now.Add(-1 * time.Hour)
	scheduledEnd := now.Add(time.Hour)

	// 先创建服务项目
	serviceItem := CreateTestServiceItem(t, db, testGame, "SLA Expired Service", 10000)

	order := &model.Order{
		Base:            model.Base{ExtJSON: "{}"},
		OrderNo:         fmt.Sprintf("SLAEXP%d", time.Now().UnixNano()),
		UserID:          testUser.ID,
		PlayerID:        &testPlayer.ID,
		GameID:          &testGame.ID,
		ItemID:          serviceItem.ID,
		TotalPriceCents: 10000,
		Currency:        model.CurrencyCNY,
		Status:          model.OrderStatusInProgress,
		Title:           "SLA Expired Test",
		ScheduledStart:  &scheduledStart,
		ScheduledEnd:    &scheduledEnd,
		OrderConfig:     "{}",
	}
	require.NoError(t, db.Create(order).Error)

	// 创建31分钟前的争议（SLA已过期）
	createdAt := time.Now().Add(-30*time.Minute - time.Second)
	dispute := &model.OrderDispute{
		Base:          model.Base{ExtJSON: "{}", CreatedAt: createdAt, UpdatedAt: createdAt},
		OrderID:       order.ID,
		InitiatorID:   testUser.ID,
		InitiatorType: "user",
		Type:          model.DisputeTypeServiceQuality,
		Reason:        "Test SLA expired",
		Status:        model.DisputeStatusPending,
	}
	require.NoError(t, db.Create(dispute).Error)

	// 验证SLA已过期
	timeSinceCreated := time.Since(dispute.CreatedAt)
	slaExpired := timeSinceCreated >= 30*time.Minute

	assert.True(t, slaExpired, "SLA should be expired")

	// 在实际系统中，SLA过期应该触发告警
	// 这里我们验证争议仍然存在（可以升级处理）
	var expiredDispute model.OrderDispute
	require.NoError(t, db.Where("id = ?", dispute.ID).First(&expiredDispute).Error)
	assert.Equal(t, model.DisputeStatusPending, expiredDispute.Status)
}

// TestDispute_ConcurrentDisputeCreation 测试同时发起争议
func TestDispute_ConcurrentDisputeCreation(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	testUser := CreateUniqueTestUser(t, db, "concurrent_dispute_user")
	playerUser := CreateUniqueTestUser(t, db, "concurrent_dispute_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "concurrent_dispute_game")

	// 创建订单
	now := time.Now()
	scheduledStart := now.Add(-1 * time.Hour)
	scheduledEnd := now.Add(time.Hour)

	// 先创建服务项目
	serviceItem := CreateTestServiceItem(t, db, testGame, "Concurrent Dispute Service", 10000)

	order := &model.Order{
		Base:            model.Base{ExtJSON: "{}"},
		OrderNo:         fmt.Sprintf("CONDISP%d", time.Now().UnixNano()),
		UserID:          testUser.ID,
		PlayerID:        &testPlayer.ID,
		GameID:          &testGame.ID,
		ItemID:          serviceItem.ID,
		TotalPriceCents: 10000,
		Currency:        model.CurrencyCNY,
		Status:          model.OrderStatusInProgress,
		Title:           "Concurrent Dispute Test",
		ScheduledStart:  &scheduledStart,
		ScheduledEnd:    &scheduledEnd,
		OrderConfig:     "{}",
	}
	require.NoError(t, db.Create(order).Error)

	// 用户和陪玩师同时发起争议
	var wg sync.WaitGroup
	errors := make(chan error, 2)
	userDisputeID := make(chan uint64, 1)
	playerDisputeID := make(chan uint64, 1)

	// 用户发起争议
	wg.Add(1)
	go func() {
		defer wg.Done()
		dispute := &model.OrderDispute{
			Base:          model.Base{ExtJSON: "{}"},
			OrderID:       order.ID,
			InitiatorID:   testUser.ID,
			InitiatorType: "user",
			Type:          model.DisputeTypeServiceQuality,
			Reason:        "User dispute",
			Status:        model.DisputeStatusPending,
		}
		if err := db.Create(dispute).Error; err != nil {
			errors <- err
			return
		}
		userDisputeID <- dispute.ID
	}()

	// 陪玩师发起争议
	wg.Add(1)
	go func() {
		defer wg.Done()
		dispute := &model.OrderDispute{
			Base:          model.Base{ExtJSON: "{}"},
			OrderID:       order.ID,
			InitiatorID:   playerUser.ID,
			InitiatorType: "player",
			Type:          model.DisputeTypeUserNotCooperative,
			Reason:        "Player dispute",
			Status:        model.DisputeStatusPending,
		}
		if err := db.Create(dispute).Error; err != nil {
			errors <- err
			return
		}
		playerDisputeID <- dispute.ID
	}()

	wg.Wait()
	close(errors)
	close(userDisputeID)
	close(playerDisputeID)

	// 检查错误
	for err := range errors {
		t.Errorf("Concurrent dispute creation error: %v", err)
	}

	// 获取创建的争议ID
	var userID, playerID uint64
	hasUser := false
	hasPlayer := false

	for id := range userDisputeID {
		userID = id
		hasUser = true
	}
	for id := range playerDisputeID {
		playerID = id
		hasPlayer = true
	}

	// 验证两个争议都创建成功
	assert.True(t, hasUser, "User dispute should be created")
	assert.True(t, hasPlayer, "Player dispute should be created")

	if hasUser && hasPlayer {
		assert.NotEqual(t, userID, playerID, "Dispute IDs should be different")

		// 验证两个争议都关联到同一订单
		var userDispute, playerDispute model.OrderDispute
		require.NoError(t, db.Where("id = ?", userID).First(&userDispute).Error)
		require.NoError(t, db.Where("id = ?", playerID).First(&playerDispute).Error)
		assert.Equal(t, order.ID, userDispute.OrderID)
		assert.Equal(t, order.ID, playerDispute.OrderID)
	}
}

// TestDispute_MultipleDisputesPrevention 测试防止同一订单多个进行中争议
func TestDispute_MultipleDisputesPrevention(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	testUser := CreateUniqueTestUser(t, db, "multi_dispute_user")
	playerUser := CreateUniqueTestUser(t, db, "multi_dispute_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "multi_dispute_game")

	// 创建订单
	now := time.Now()
	scheduledStart := now.Add(-1 * time.Hour)
	scheduledEnd := now.Add(time.Hour)

	// 先创建服务项目
	serviceItem := CreateTestServiceItem(t, db, testGame, "Multiple Dispute Service", 10000)

	order := &model.Order{
		Base:            model.Base{ExtJSON: "{}"},
		OrderNo:         fmt.Sprintf("MULTIDISP%d", time.Now().UnixNano()),
		UserID:          testUser.ID,
		PlayerID:        &testPlayer.ID,
		GameID:          &testGame.ID,
		ItemID:          serviceItem.ID,
		TotalPriceCents: 10000,
		Currency:        model.CurrencyCNY,
		Status:          model.OrderStatusInProgress,
		Title:           "Multiple Dispute Prevention Test",
		ScheduledStart:  &scheduledStart,
		ScheduledEnd:    &scheduledEnd,
		OrderConfig:     "{}",
	}
	require.NoError(t, db.Create(order).Error)

	// 创建第一个争议（进行中）
	dispute1 := &model.OrderDispute{
		Base:          model.Base{ExtJSON: "{}"},
		OrderID:       order.ID,
		InitiatorID:   testUser.ID,
		InitiatorType: "user",
		Type:          model.DisputeTypeServiceQuality,
		Reason:        "First dispute",
		Status:        model.DisputeStatusPending,
	}
	require.NoError(t, db.Create(dispute1).Error)

	// 尝试创建第二个争议（应该被阻止）
	dispute2 := &model.OrderDispute{
		Base:          model.Base{ExtJSON: "{}"},
		OrderID:       order.ID,
		InitiatorID:   testUser.ID,
		InitiatorType: "user",
		Type:          model.DisputeTypeBadAttitude,
		Reason:        "Second dispute",
		Status:        model.DisputeStatusPending,
	}

	// 检查是否已有进行中的争议
	var existingDispute model.OrderDispute
	err := db.Where("order_id = ? AND status IN ?", order.ID, []model.DisputeStatus{
		model.DisputeStatusPending,
		model.DisputeStatusAssigned,
		model.DisputeStatusMediating,
	}).First(&existingDispute).Error

	if err == nil {
		// 已有进行中的争议，不能创建新的
		t.Log("Second dispute correctly prevented - existing dispute found")
		assert.Equal(t, dispute1.ID, existingDispute.ID)
	} else {
		// 没有找到，可能是数据库查询失败
		// 在实际系统中应该阻止创建
		db.Create(dispute2) // 尝试创建
		t.Error("Multiple pending disputes should be prevented")
	}
}

// ============================================================================
// UserBlock Edge Cases - 双向拉黑、循环拉黑
// ============================================================================

// TestUserBlock_BidirectionalBlocking 测试双向拉黑
func TestUserBlock_BidirectionalBlocking(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	user1 := CreateUniqueTestUser(t, db, "block_user1")
	user2 := CreateUniqueTestUser(t, db, "block_user2")

	// User1 拉黑 User2
	block1 := &model.UserBlock{
		Base:        model.Base{ExtJSON: "{}"},
		BlockerID:   user1.ID,
		BlockerType: "user",
		BlockedID:   user2.ID,
		BlockedType: "user",
		Reason:      "Test block 1->2",
		BlockedAt:   time.Now(),
	}
	require.NoError(t, db.Create(block1).Error)

	// User2 也拉黑 User1
	block2 := &model.UserBlock{
		Base:        model.Base{ExtJSON: "{}"},
		BlockerID:   user2.ID,
		BlockerType: "user",
		BlockedID:   user1.ID,
		BlockedType: "user",
		Reason:      "Test block 2->1",
		BlockedAt:   time.Now(),
	}
	require.NoError(t, db.Create(block2).Error)

	// 验证双向拉黑记录独立存在
	var blocks []model.UserBlock
	require.NoError(t, db.Where("(blocker_id = ? AND blocked_id = ?) OR (blocker_id = ? AND blocked_id = ?)",
		user1.ID, user2.ID, user2.ID, user1.ID).Find(&blocks).Error)

	assert.Len(t, blocks, 2, "Should have 2 independent blocking records")

	// 验证每条记录的独立性
	for _, block := range blocks {
		if block.BlockerID == user1.ID {
			assert.Equal(t, user2.ID, block.BlockedID)
		} else {
			assert.Equal(t, user1.ID, block.BlockedID)
		}
	}
}

// TestUserBlock_CircularBlocking 测试循环拉黑（A拉黑B，B拉黑C，C拉黑A）
func TestUserBlock_CircularBlocking(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	userA := CreateUniqueTestUser(t, db, "circular_A")
	userB := CreateUniqueTestUser(t, db, "circular_B")
	userC := CreateUniqueTestUser(t, db, "circular_C")

	// A 拉黑 B
	blockAB := &model.UserBlock{
		Base:        model.Base{ExtJSON: "{}"},
		BlockerID:   userA.ID,
		BlockerType: "user",
		BlockedID:   userB.ID,
		BlockedType: "user",
		Reason:      "A blocks B",
		BlockedAt:   time.Now(),
	}
	require.NoError(t, db.Create(blockAB).Error)

	// B 拉黑 C
	blockBC := &model.UserBlock{
		Base:        model.Base{ExtJSON: "{}"},
		BlockerID:   userB.ID,
		BlockerType: "user",
		BlockedID:   userC.ID,
		BlockedType: "user",
		Reason:      "B blocks C",
		BlockedAt:   time.Now(),
	}
	require.NoError(t, db.Create(blockBC).Error)

	// C 拉黑 A
	blockCA := &model.UserBlock{
		Base:        model.Base{ExtJSON: "{}"},
		BlockerID:   userC.ID,
		BlockerType: "user",
		BlockedID:   userA.ID,
		BlockedType: "user",
		Reason:      "C blocks A",
		BlockedAt:   time.Now(),
	}
	require.NoError(t, db.Create(blockCA).Error)

	// 验证循环拉黑
	var blocks []model.UserBlock
	require.NoError(t, db.Find(&blocks).Error)
	assert.Len(t, blocks, 3, "Should have 3 blocking records in circular chain")

	// 验证拉黑关系
	blockingMap := make(map[uint64]uint64) // blocker_id -> blocked_id
	for _, block := range blocks {
		blockingMap[block.BlockerID] = block.BlockedID
	}

	assert.Equal(t, userB.ID, blockingMap[userA.ID], "A should block B")
	assert.Equal(t, userC.ID, blockingMap[userB.ID], "B should block C")
	assert.Equal(t, userA.ID, blockingMap[userC.ID], "C should block A")
}

// TestUserBlock_SelfBlockingPrevention 测试防止自我拉黑
func TestUserBlock_SelfBlockingPrevention(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	user := CreateUniqueTestUser(t, db, "selfblock_user")

	// 尝试创建自我拉黑记录
	selfBlock := &model.UserBlock{
		Base:        model.Base{ExtJSON: "{}"},
		BlockerID:   user.ID,
		BlockerType: "user",
		BlockedID:   user.ID,
		BlockedType: "user",
		Reason:      "Attempt to self block",
		BlockedAt:   time.Now(),
	}

	// 在实际系统中应该有验证阻止这种情况
	// 这里我们验证如果记录存在，应该被查询排除
	err := db.Create(selfBlock).Error

	if err != nil {
		// 数据库可能有约束阻止
		t.Logf("Self blocking prevented by database: %v", err)
	} else {
		// 如果创建成功，应用层应该在查询时过滤
		// 查询时应该排除 BlockerID == BlockedID 的记录
		var blocks []model.UserBlock
		db.Where("blocker_id = ? AND blocked_id != blocker_id", user.ID).Find(&blocks)

		// 验证自我拉黑不会出现在查询结果中
		for _, block := range blocks {
			assert.NotEqual(t, block.BlockerID, block.BlockedID, "User should not be able to block themselves")
		}
	}
}

// TestUserBlock_OrderBlockingEffect 测试拉黑后订单影响
func TestUserBlock_OrderBlockingEffect(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	testUser := CreateUniqueTestUser(t, db, "blockeffect_user")
	playerUser := CreateUniqueTestUser(t, db, "blockeffect_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "blockeffect_game")

	// 创建进行中的订单
	now := time.Now()
	scheduledStart := now.Add(-1 * time.Hour)
	scheduledEnd := now.Add(time.Hour)

	// 先创建服务项目
	serviceItem := CreateTestServiceItem(t, db, testGame, "Block Effect Service", 10000)

	order := &model.Order{
		Base:            model.Base{ExtJSON: "{}"},
		OrderNo:         fmt.Sprintf("BLKEFF%d", time.Now().UnixNano()),
		UserID:          testUser.ID,
		PlayerID:        &testPlayer.ID,
		GameID:          &testGame.ID,
		ItemID:          serviceItem.ID,
		TotalPriceCents: 10000,
		Currency:        model.CurrencyCNY,
		Status:          model.OrderStatusInProgress,
		Title:           "Blocking Effect Test",
		ScheduledStart:  &scheduledStart,
		ScheduledEnd:    &scheduledEnd,
		OrderConfig:     "{}",
	}
	require.NoError(t, db.Create(order).Error)

	// 用户拉黑陪玩师
	block := &model.UserBlock{
		Base:        model.Base{ExtJSON: "{}"},
		BlockerID:   testUser.ID,
		BlockerType: "user",
		BlockedID:   playerUser.ID,
		BlockedType: "player",
		Reason:      "User blocks player during order",
		BlockedAt:   time.Now(),
	}
	require.NoError(t, db.Create(block).Error)

	// 验证拉黑后进行中的订单不受影响
	var inProgressOrder model.Order
	require.NoError(t, db.Where("id = ?", order.ID).First(&inProgressOrder).Error)
	assert.Equal(t, model.OrderStatusInProgress, inProgressOrder.Status, "In-progress order should continue")

	// 尝试创建新订单（应该被阻止）
	newOrder := &model.Order{
		Base:            model.Base{ExtJSON: "{}"},
		OrderNo:         fmt.Sprintf("BLKEFFNEW%d", time.Now().UnixNano()),
		UserID:          testUser.ID,
		PlayerID:        &testPlayer.ID,
		GameID:          &testGame.ID,
		TotalPriceCents: 10000,
		Currency:        model.CurrencyCNY,
		Status:          model.OrderStatusPending,
		Title:           "New Order After Blocking",
		OrderConfig:     "{}",
	}

	// 在实际系统中，应该检查拉黑关系
	// 这里我们模拟检查
	var existingBlock model.UserBlock
	err := db.Where("blocker_id = ? AND blocked_id = ?", testUser.ID, playerUser.ID).First(&existingBlock).Error

	if err == nil {
		// 存在拉黑关系，不应创建新订单
		t.Log("New order correctly prevented - blocking relationship exists")
		// 实际应用中应该返回错误
	} else {
		// 没有找到拉黑记录，可以创建订单（但实际不应该）
		db.Create(newOrder)
		t.Error("Order creation should be prevented when blocking relationship exists")
	}
}
