package withdraw

import (
	"context"
	"testing"

	"gamelink/internal/model"
	"gamelink/internal/repository"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(
		&model.Withdraw{},
		&model.Order{},
		&model.Player{},
	)
	require.NoError(t, err)

	return db
}

func TestWithdrawRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewWithdrawRepository(db)
	ctx := context.Background()

	withdraw := &model.Withdraw{
		PlayerID:    1,
		UserID:      1,
		AmountCents: 10000,
		Method:      model.WithdrawMethodAlipay,
		AccountInfo: "test@example.com",
		Status:      model.WithdrawStatusPending,
	}

	err := repo.Create(ctx, withdraw)

	assert.NoError(t, err)
	assert.NotZero(t, withdraw.ID)
}

func TestWithdrawRepository_Get(t *testing.T) {
	db := setupTestDB(t)
	repo := NewWithdrawRepository(db)
	ctx := context.Background()

	withdraw := &model.Withdraw{
		PlayerID:    1,
		AmountCents: 10000,
	}
	err := repo.Create(ctx, withdraw)
	require.NoError(t, err)

	// 测试获取存在的记录
	result, err := repo.Get(ctx, withdraw.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(10000), result.AmountCents)

	// 测试获取不存在的记录
	_, err = repo.Get(ctx, 999)
	assert.Error(t, err)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestWithdrawRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewWithdrawRepository(db)
	ctx := context.Background()

	withdraw := &model.Withdraw{
		PlayerID:    1,
		AmountCents: 10000,
		Status:      model.WithdrawStatusPending,
	}
	err := repo.Create(ctx, withdraw)
	require.NoError(t, err)

	// 更新状态
	withdraw.Status = model.WithdrawStatusCompleted

	err = repo.Update(ctx, withdraw)
	assert.NoError(t, err)

	// 验证更新成功
	updated, err := repo.Get(ctx, withdraw.ID)
	assert.NoError(t, err)
	assert.Equal(t, model.WithdrawStatusCompleted, updated.Status)
}

func TestWithdrawRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewWithdrawRepository(db)
	ctx := context.Background()

	// 创建测试数据
	withdraws := []model.Withdraw{
		{PlayerID: 1, UserID: 1, AmountCents: 10000, Status: model.WithdrawStatusPending},
		{PlayerID: 1, UserID: 1, AmountCents: 20000, Status: model.WithdrawStatusCompleted},
		{PlayerID: 2, UserID: 2, AmountCents: 30000, Status: model.WithdrawStatusPending},
	}
	for _, w := range withdraws {
		err := repo.Create(ctx, &w)
		require.NoError(t, err)
	}

	// 测试按PlayerID过滤
	playerID := uint64(1)
	opts := WithdrawListOptions{
		PlayerID: &playerID,
		Page:     1,
		PageSize: 10,
	}
	results, total, err := repo.List(ctx, opts)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, results, 2)

	// 测试按Status过滤
	status := model.WithdrawStatusPending
	opts2 := WithdrawListOptions{
		Status:   &status,
		Page:     1,
		PageSize: 10,
	}
	results2, total2, err := repo.List(ctx, opts2)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total2)
	assert.Len(t, results2, 2)
}

func TestWithdrawRepository_GetPlayerBalance(t *testing.T) {
	db := setupTestDB(t)
	repo := NewWithdrawRepository(db)
	ctx := context.Background()

	// 创建测试player
	playerID := uint64(1)
	player := model.Player{Base: model.Base{ID: playerID}, UserID: 1}
	db.Create(&player)

	// 创建完成的订单
	orders := []model.Order{
		{PlayerID: &playerID, Status: model.OrderStatusCompleted, TotalPriceCents: 10000},
		{PlayerID: &playerID, Status: model.OrderStatusCompleted, TotalPriceCents: 20000},
		{PlayerID: &playerID, Status: model.OrderStatusInProgress, TotalPriceCents: 5000},
	}
	for _, o := range orders {
		db.Create(&o)
	}

	// 创建提现记录
	withdraws := []model.Withdraw{
		{PlayerID: playerID, UserID: 1, AmountCents: 5000, Status: model.WithdrawStatusCompleted},
		{PlayerID: playerID, UserID: 1, AmountCents: 3000, Status: model.WithdrawStatusPending},
	}
	for _, w := range withdraws {
		err := repo.Create(ctx, &w)
		require.NoError(t, err)
	}

	// 获取余额
	balance, err := repo.GetPlayerBalance(ctx, playerID)
	assert.NoError(t, err)
	assert.NotNil(t, balance)
	assert.Equal(t, int64(30000), balance.TotalEarnings)    // 10000 + 20000
	assert.Equal(t, int64(5000), balance.WithdrawTotal)     // 已完成提现
	assert.Equal(t, int64(3000), balance.PendingWithdraw)   // 待处理提现
	assert.Equal(t, int64(5000), balance.PendingBalance)    // 进行中的订单
	assert.Equal(t, int64(17000), balance.AvailableBalance) // 30000 - 5000 - 3000 - 5000
}

