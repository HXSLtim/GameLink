package user

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
	orderimpl "gamelink/internal/repository/implementations"
	userrepo "gamelink/internal/repository/user"
	withdrawrepo "gamelink/internal/repository/withdraw"
	"gamelink/pkg/testutil"
)

func setupEarningsTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.User{},
		&model.Player{},
		&model.Order{},
		&model.Withdraw{},
		&model.Game{},
		&model.Wallet{},
	)
	return db
}

func createEarningsTestData(t *testing.T, db *gorm.DB) (*model.User, *model.Player) {
	t.Helper()

	// 创建用户
	user := &model.User{
		Phone:        "13800000001",
		Email:        "player@test.com",
		Name:         "Player User",
		Role:         model.RolePlayer,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(user).Error)

	// 创建游戏
	game := &model.Game{Key: "lol", Name: "英雄联盟", Category: "moba"}
	require.NoError(t, db.Create(game).Error)

	// 创建陪玩师
	player := &model.Player{
		UserID:             user.ID,
		Nickname:           "Pro Player",
		MainGameID:         game.ID,
		HourlyRateCents:    5000,
		VerificationStatus: model.VerificationVerified,
	}
	require.NoError(t, db.Create(player).Error)

	return user, player
}

func createEarningsService(db *gorm.DB) *EarningsService {
	playerRepo := userrepo.NewPlayerRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	return NewEarningsService(playerRepo, orderRepo, withdrawRepo)
}

func TestEarningsService_GetEarningsSummary(t *testing.T) {
	db := setupEarningsTestDB(t)
	defer testutil.CleanDB(t, db)

	user, player := createEarningsTestData(t, db)
	svc := createEarningsService(db)

	// 创建已完成的订单
	now := time.Now()
	order := &model.Order{
		UserID:          1,
		ItemID:          1,
		Title:           "测试订单",
		Status:          model.OrderStatusCompleted,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		CompletedAt:     &now,
	}
	order.SetPlayerID(player.ID)
	require.NoError(t, db.Create(order).Error)

	t.Run("获取收益概览成功", func(t *testing.T) {
		summary, err := svc.GetEarningsSummary(context.Background(), user.ID)
		require.NoError(t, err)
		assert.NotNil(t, summary)
		assert.GreaterOrEqual(t, summary.TotalEarnings, int64(0))
	})

	t.Run("非陪玩师获取收益应失败", func(t *testing.T) {
		// 创建普通用户
		normalUser := &model.User{
			Phone:        "13800000002",
			Email:        "normal@test.com",
			Name:         "Normal User",
			Role:         model.RoleUser,
			Status:       model.UserStatusActive,
			PasswordHash: "hashed",
		}
		require.NoError(t, db.Create(normalUser).Error)

		_, err := svc.GetEarningsSummary(context.Background(), normalUser.ID)
		assert.Error(t, err)
	})
}

func TestEarningsService_GetEarningsTrend(t *testing.T) {
	db := setupEarningsTestDB(t)
	defer testutil.CleanDB(t, db)

	user, _ := createEarningsTestData(t, db)
	svc := createEarningsService(db)

	t.Run("获取7天收益趋势", func(t *testing.T) {
		trend, err := svc.GetEarningsTrend(context.Background(), user.ID, 7)
		require.NoError(t, err)
		assert.NotNil(t, trend)
		assert.Len(t, trend.Trend, 7)
	})

	t.Run("获取30天收益趋势", func(t *testing.T) {
		trend, err := svc.GetEarningsTrend(context.Background(), user.ID, 30)
		require.NoError(t, err)
		assert.Len(t, trend.Trend, 30)
	})

	t.Run("天数小于7自动调整为7", func(t *testing.T) {
		trend, err := svc.GetEarningsTrend(context.Background(), user.ID, 3)
		require.NoError(t, err)
		assert.Len(t, trend.Trend, 7)
	})

	t.Run("天数大于90自动调整为90", func(t *testing.T) {
		trend, err := svc.GetEarningsTrend(context.Background(), user.ID, 100)
		require.NoError(t, err)
		assert.Len(t, trend.Trend, 90)
	})

	t.Run("非陪玩师获取趋势应失败", func(t *testing.T) {
		_, err := svc.GetEarningsTrend(context.Background(), 99999, 7)
		assert.Error(t, err)
	})
}

func TestEarningsService_GetWithdrawHistory(t *testing.T) {
	db := setupEarningsTestDB(t)
	defer testutil.CleanDB(t, db)

	user, player := createEarningsTestData(t, db)
	svc := createEarningsService(db)

	// 创建提现记录
	withdraw := &model.Withdraw{
		PlayerID:    player.ID,
		UserID:      user.ID,
		AmountCents: 10000,
		Method:      model.WithdrawMethodAlipay,
		AccountInfo: "test@alipay.com",
		Status:      model.WithdrawStatusPending,
	}
	require.NoError(t, db.Create(withdraw).Error)

	t.Run("获取提现记录成功", func(t *testing.T) {
		history, err := svc.GetWithdrawHistory(context.Background(), user.ID, 1, 10)
		require.NoError(t, err)
		assert.NotNil(t, history)
		assert.Equal(t, int64(1), history.Total)
		assert.Len(t, history.Records, 1)
	})

	t.Run("分页参数默认值", func(t *testing.T) {
		history, err := svc.GetWithdrawHistory(context.Background(), user.ID, 0, 0)
		require.NoError(t, err)
		assert.NotNil(t, history)
	})

	t.Run("非陪玩师获取提现记录应失败", func(t *testing.T) {
		_, err := svc.GetWithdrawHistory(context.Background(), 99999, 1, 10)
		assert.Error(t, err)
	})
}

func TestEarningsService_RequestWithdraw(t *testing.T) {
	db := setupEarningsTestDB(t)
	defer testutil.CleanDB(t, db)

	user, player := createEarningsTestData(t, db)
	svc := createEarningsService(db)

	// 创建已完成的订单以产生收益
	now := time.Now()
	for i := 0; i < 10; i++ {
		order := &model.Order{
			UserID:          1,
			ItemID:          1,
			Title:           "测试订单",
			Status:          model.OrderStatusCompleted,
			UnitPriceCents:  50000,
			TotalPriceCents: 100000, // 1000元
			CompletedAt:     &now,
		}
		order.SetPlayerID(player.ID)
		require.NoError(t, db.Create(order).Error)
	}

	t.Run("提现金额过小应失败", func(t *testing.T) {
		_, err := svc.RequestWithdraw(context.Background(), user.ID, WithdrawRequest{
			AmountCents: 5000, // 50元，小于100元
			Method:      "alipay",
			AccountInfo: "test@alipay.com",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "100元")
	})

	t.Run("非陪玩师申请提现应失败", func(t *testing.T) {
		_, err := svc.RequestWithdraw(context.Background(), 99999, WithdrawRequest{
			AmountCents: 10000,
			Method:      "alipay",
			AccountInfo: "test@alipay.com",
		})
		assert.Error(t, err)
	})
}

func TestEarningsService_FormatCents(t *testing.T) {
	tests := []struct {
		cents    int64
		expected string
	}{
		{10000, "100.00元"},
		{12345, "123.45元"},
		{100, "1.00元"},
		{0, "0.00元"},
		{1, "0.01元"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatCents(tt.cents)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEarningsService_StrPtr(t *testing.T) {
	t.Run("字符串指针辅助函数", func(t *testing.T) {
		s := "test"
		ptr := strPtr(s)
		assert.NotNil(t, ptr)
		assert.Equal(t, s, *ptr)
	})
}

func TestEarningsService_Errors(t *testing.T) {
	t.Run("错误变量定义正确", func(t *testing.T) {
		assert.NotNil(t, ErrNotFound)
		assert.NotNil(t, ErrValidation)
		assert.NotNil(t, ErrInsufficientBalance)
		assert.NotNil(t, ErrUnauthorized)
		assert.NotNil(t, ErrDailyLimitExceeded)
		assert.NotNil(t, ErrMonthlyLimitExceeded)
		assert.NotNil(t, ErrPendingWithdrawExists)
	})
}

func TestWithdrawStatus(t *testing.T) {
	t.Run("提现状态常量", func(t *testing.T) {
		assert.Equal(t, WithdrawStatus("pending"), WithdrawPending)
		assert.Equal(t, WithdrawStatus("approved"), WithdrawApproved)
		assert.Equal(t, WithdrawStatus("rejected"), WithdrawRejected)
		assert.Equal(t, WithdrawStatus("completed"), WithdrawCompleted)
	})
}

// ============================================
// Additional tests to improve coverage to 80%+
// ============================================

func TestEarningsService_RequestWithdraw_Success(t *testing.T) {
	db := setupEarningsTestDB(t)
	defer testutil.CleanDB(t, db)

	user, player := createEarningsTestData(t, db)
	svc := createEarningsService(db)

	// 创建大量已完成订单以产生足够收益
	now := time.Now()
	for i := 0; i < 20; i++ {
		order := &model.Order{
			UserID:          1,
			ItemID:          1,
			Title:           "测试订单",
			Status:          model.OrderStatusCompleted,
			UnitPriceCents:  100000,
			TotalPriceCents: 200000, // 2000元
			CompletedAt:     &now,
		}
		order.SetPlayerID(player.ID)
		require.NoError(t, db.Create(order).Error)
	}

	t.Run("提现成功", func(t *testing.T) {
		resp, err := svc.RequestWithdraw(context.Background(), user.ID, WithdrawRequest{
			AmountCents: 10000, // 100元
			Method:      "alipay",
			AccountInfo: "test@alipay.com",
		})
		require.NoError(t, err)
		assert.NotZero(t, resp.WithdrawID)
		assert.Equal(t, "pending", resp.Status)
	})
}

func TestEarningsService_RequestWithdraw_AmountTooLarge(t *testing.T) {
	db := setupEarningsTestDB(t)
	defer testutil.CleanDB(t, db)

	user, _ := createEarningsTestData(t, db)
	svc := createEarningsService(db)

	t.Run("提现金额超过单笔限额", func(t *testing.T) {
		_, err := svc.RequestWithdraw(context.Background(), user.ID, WithdrawRequest{
			AmountCents: 1000000001, // 超过1000万元
			Method:      "alipay",
			AccountInfo: "test@alipay.com",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "1000万元")
	})
}

func TestEarningsService_RequestWithdraw_PendingExists(t *testing.T) {
	db := setupEarningsTestDB(t)
	defer testutil.CleanDB(t, db)

	user, player := createEarningsTestData(t, db)
	svc := createEarningsService(db)

	// 创建已完成订单
	now := time.Now()
	for i := 0; i < 10; i++ {
		order := &model.Order{
			UserID:          1,
			ItemID:          1,
			Title:           "测试订单",
			Status:          model.OrderStatusCompleted,
			UnitPriceCents:  100000,
			TotalPriceCents: 200000,
			CompletedAt:     &now,
		}
		order.SetPlayerID(player.ID)
		require.NoError(t, db.Create(order).Error)
	}

	// 创建一个待处理的提现记录
	withdraw := &model.Withdraw{
		PlayerID:    player.ID,
		UserID:      user.ID,
		AmountCents: 10000,
		Method:      model.WithdrawMethodAlipay,
		AccountInfo: "test@alipay.com",
		Status:      model.WithdrawStatusPending,
	}
	require.NoError(t, db.Create(withdraw).Error)

	t.Run("存在待处理提现时应失败", func(t *testing.T) {
		_, err := svc.RequestWithdraw(context.Background(), user.ID, WithdrawRequest{
			AmountCents: 10000,
			Method:      "alipay",
			AccountInfo: "test2@alipay.com",
		})
		assert.ErrorIs(t, err, ErrPendingWithdrawExists)
	})
}

func TestEarningsService_RequestWithdraw_DailyLimitExceeded(t *testing.T) {
	db := setupEarningsTestDB(t)
	defer testutil.CleanDB(t, db)

	user, player := createEarningsTestData(t, db)
	svc := createEarningsService(db)

	// 创建大量已完成订单
	now := time.Now()
	for i := 0; i < 50; i++ {
		order := &model.Order{
			UserID:          1,
			ItemID:          1,
			Title:           "测试订单",
			Status:          model.OrderStatusCompleted,
			UnitPriceCents:  500000,
			TotalPriceCents: 1000000, // 10000元
			CompletedAt:     &now,
		}
		order.SetPlayerID(player.ID)
		require.NoError(t, db.Create(order).Error)
	}

	// 创建今日已完成的提现记录（接近每日限额）
	withdraw := &model.Withdraw{
		PlayerID:    player.ID,
		UserID:      user.ID,
		AmountCents: 4900000, // 49000元，接近每日限额5万
		Method:      model.WithdrawMethodAlipay,
		AccountInfo: "test@alipay.com",
		Status:      model.WithdrawStatusCompleted,
	}
	require.NoError(t, db.Create(withdraw).Error)

	t.Run("超过每日提现限额", func(t *testing.T) {
		_, err := svc.RequestWithdraw(context.Background(), user.ID, WithdrawRequest{
			AmountCents: 200000, // 2000元，会超过每日限额
			Method:      "alipay",
			AccountInfo: "test2@alipay.com",
		})
		assert.Error(t, err)
		// 应该提示剩余额度
		assert.Contains(t, err.Error(), "每日")
	})
}

func TestEarningsService_RequestWithdraw_MonthlyLimitExceeded(t *testing.T) {
	db := setupEarningsTestDB(t)
	defer testutil.CleanDB(t, db)

	user, player := createEarningsTestData(t, db)
	svc := createEarningsService(db)

	// 创建大量已完成订单
	now := time.Now()
	for i := 0; i < 100; i++ {
		order := &model.Order{
			UserID:          1,
			ItemID:          1,
			Title:           "测试订单",
			Status:          model.OrderStatusCompleted,
			UnitPriceCents:  1000000,
			TotalPriceCents: 2000000, // 20000元
			CompletedAt:     &now,
		}
		order.SetPlayerID(player.ID)
		require.NoError(t, db.Create(order).Error)
	}

	// 创建本月已完成的提现记录（接近每月限额20万）
	// 使用过去几天的日期，避免触发每日限额（5万）
	// 每天最多5万，所以分4天，每天4.9万 = 19.6万
	for i := 1; i <= 4; i++ {
		pastDate := now.AddDate(0, 0, -i)
		withdraw := &model.Withdraw{
			PlayerID:    player.ID,
			UserID:      user.ID,
			AmountCents: 4900000, // 49000元
			Method:      model.WithdrawMethodAlipay,
			AccountInfo: "test@alipay.com",
			Status:      model.WithdrawStatusCompleted,
		}
		require.NoError(t, db.Create(withdraw).Error)
		// 手动更新创建时间到过去
		db.Model(withdraw).Update("created_at", pastDate)
	}

	t.Run("超过每月提现限额", func(t *testing.T) {
		_, err := svc.RequestWithdraw(context.Background(), user.ID, WithdrawRequest{
			AmountCents: 500000, // 5000元，会超过每月限额（19.6万+0.5万=20.1万>20万）
			Method:      "alipay",
			AccountInfo: "test2@alipay.com",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "每月")
	})
}

func TestEarningsService_RequestWithdraw_InsufficientBalance(t *testing.T) {
	db := setupEarningsTestDB(t)
	defer testutil.CleanDB(t, db)

	user, player := createEarningsTestData(t, db)
	svc := createEarningsService(db)

	// 只创建少量订单
	now := time.Now()
	order := &model.Order{
		UserID:          1,
		ItemID:          1,
		Title:           "测试订单",
		Status:          model.OrderStatusCompleted,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000, // 100元
		CompletedAt:     &now,
	}
	order.SetPlayerID(player.ID)
	require.NoError(t, db.Create(order).Error)

	t.Run("余额不足", func(t *testing.T) {
		_, err := svc.RequestWithdraw(context.Background(), user.ID, WithdrawRequest{
			AmountCents: 1000000, // 10000元，超过可提现余额
			Method:      "alipay",
			AccountInfo: "test@alipay.com",
		})
		assert.ErrorIs(t, err, ErrInsufficientBalance)
	})
}

func TestEarningsService_HasPendingWithdraw(t *testing.T) {
	db := setupEarningsTestDB(t)
	defer testutil.CleanDB(t, db)

	user, player := createEarningsTestData(t, db)
	svc := createEarningsService(db)

	t.Run("无待处理提现", func(t *testing.T) {
		hasPending, err := svc.hasPendingWithdraw(context.Background(), player.ID)
		require.NoError(t, err)
		assert.False(t, hasPending)
	})

	t.Run("有待处理提现", func(t *testing.T) {
		withdraw := &model.Withdraw{
			PlayerID:    player.ID,
			UserID:      user.ID,
			AmountCents: 10000,
			Method:      model.WithdrawMethodAlipay,
			AccountInfo: "test@alipay.com",
			Status:      model.WithdrawStatusPending,
		}
		require.NoError(t, db.Create(withdraw).Error)

		hasPending, err := svc.hasPendingWithdraw(context.Background(), player.ID)
		require.NoError(t, err)
		assert.True(t, hasPending)
	})
}

func TestEarningsService_CalculateDailyWithdrawTotal(t *testing.T) {
	db := setupEarningsTestDB(t)
	defer testutil.CleanDB(t, db)

	user, player := createEarningsTestData(t, db)
	svc := createEarningsService(db)

	t.Run("无提现记录", func(t *testing.T) {
		total, err := svc.calculateDailyWithdrawTotal(context.Background(), player.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
	})

	t.Run("有今日提现记录", func(t *testing.T) {
		// 创建今日提现记录
		withdraw := &model.Withdraw{
			PlayerID:    player.ID,
			UserID:      user.ID,
			AmountCents: 10000,
			Method:      model.WithdrawMethodAlipay,
			AccountInfo: "test@alipay.com",
			Status:      model.WithdrawStatusCompleted,
		}
		require.NoError(t, db.Create(withdraw).Error)

		total, err := svc.calculateDailyWithdrawTotal(context.Background(), player.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(10000), total)
	})

	t.Run("排除已拒绝的提现", func(t *testing.T) {
		// 创建被拒绝的提现记录
		withdraw := &model.Withdraw{
			PlayerID:    player.ID,
			UserID:      user.ID,
			AmountCents: 50000,
			Method:      model.WithdrawMethodAlipay,
			AccountInfo: "test@alipay.com",
			Status:      model.WithdrawStatusRejected,
		}
		require.NoError(t, db.Create(withdraw).Error)

		total, err := svc.calculateDailyWithdrawTotal(context.Background(), player.ID)
		require.NoError(t, err)
		// 被拒绝的不计入
		assert.Equal(t, int64(10000), total)
	})
}

func TestEarningsService_CalculateMonthlyWithdrawTotal(t *testing.T) {
	db := setupEarningsTestDB(t)
	defer testutil.CleanDB(t, db)

	user, player := createEarningsTestData(t, db)
	svc := createEarningsService(db)

	t.Run("无提现记录", func(t *testing.T) {
		total, err := svc.calculateMonthlyWithdrawTotal(context.Background(), player.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
	})

	t.Run("有本月提现记录", func(t *testing.T) {
		// 创建本月提现记录
		withdraw := &model.Withdraw{
			PlayerID:    player.ID,
			UserID:      user.ID,
			AmountCents: 100000,
			Method:      model.WithdrawMethodAlipay,
			AccountInfo: "test@alipay.com",
			Status:      model.WithdrawStatusCompleted,
		}
		require.NoError(t, db.Create(withdraw).Error)

		total, err := svc.calculateMonthlyWithdrawTotal(context.Background(), player.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(100000), total)
	})

	t.Run("排除已拒绝的提现", func(t *testing.T) {
		// 创建被拒绝的提现记录
		withdraw := &model.Withdraw{
			PlayerID:    player.ID,
			UserID:      user.ID,
			AmountCents: 200000,
			Method:      model.WithdrawMethodAlipay,
			AccountInfo: "test@alipay.com",
			Status:      model.WithdrawStatusRejected,
		}
		require.NoError(t, db.Create(withdraw).Error)

		total, err := svc.calculateMonthlyWithdrawTotal(context.Background(), player.ID)
		require.NoError(t, err)
		// 被拒绝的不计入
		assert.Equal(t, int64(100000), total)
	})
}

func TestEarningsService_GetEarningsSummary_WithWithdrawRepo(t *testing.T) {
	db := setupEarningsTestDB(t)
	defer testutil.CleanDB(t, db)

	user, player := createEarningsTestData(t, db)
	svc := createEarningsService(db)

	// 创建已完成订单
	now := time.Now()
	order := &model.Order{
		UserID:          1,
		ItemID:          1,
		Title:           "测试订单",
		Status:          model.OrderStatusCompleted,
		UnitPriceCents:  50000,
		TotalPriceCents: 100000,
		CompletedAt:     &now,
	}
	order.SetPlayerID(player.ID)
	require.NoError(t, db.Create(order).Error)

	// 创建提现记录
	withdraw := &model.Withdraw{
		PlayerID:    player.ID,
		UserID:      user.ID,
		AmountCents: 20000,
		Method:      model.WithdrawMethodAlipay,
		AccountInfo: "test@alipay.com",
		Status:      model.WithdrawStatusCompleted,
	}
	require.NoError(t, db.Create(withdraw).Error)

	t.Run("获取收益概览包含提现信息", func(t *testing.T) {
		summary, err := svc.GetEarningsSummary(context.Background(), user.ID)
		require.NoError(t, err)
		assert.NotNil(t, summary)
		assert.GreaterOrEqual(t, summary.TotalEarnings, int64(0))
	})
}

func TestEarningsService_CountOrders(t *testing.T) {
	db := setupEarningsTestDB(t)
	defer testutil.CleanDB(t, db)

	_, player := createEarningsTestData(t, db)
	svc := createEarningsService(db)

	// 创建已完成订单
	now := time.Now()
	for i := 0; i < 3; i++ {
		order := &model.Order{
			UserID:          1,
			ItemID:          1,
			Title:           "测试订单",
			Status:          model.OrderStatusCompleted,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			CompletedAt:     &now,
		}
		order.SetPlayerID(player.ID)
		require.NoError(t, db.Create(order).Error)
	}

	t.Run("统计订单数", func(t *testing.T) {
		playerIDPtr := &player.ID
		count, err := svc.countOrders(context.Background(), playerIDPtr, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})

	t.Run("按时间范围统计", func(t *testing.T) {
		playerIDPtr := &player.ID
		todayStart := time.Now().Truncate(24 * time.Hour)
		todayEnd := todayStart.Add(24 * time.Hour)
		count, err := svc.countOrders(context.Background(), playerIDPtr, &todayStart, &todayEnd)
		require.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})
}
