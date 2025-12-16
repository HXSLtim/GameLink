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
