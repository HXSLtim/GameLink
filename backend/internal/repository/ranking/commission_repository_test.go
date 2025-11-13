package ranking

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
)

// 覆盖 RankingCommissionRepository 的完整生命周期与过滤逻辑

func TestRankingCommissionRepository_CreateAndGetConfig(t *testing.T) {
	db := setupRankingDB(t)
	repo := NewRankingCommissionRepository(db)
	ctx := context.Background()

	cfg := &model.RankingCommissionConfig{
		Name:        "收入榜 2025-01",
		RankingType: model.RankingTypeIncome,
		Period:      "monthly",
		Month:       "2025-01",
		RulesJSON:   `[{"rankStart":1,"rankEnd":3,"commissionRate":10}]`,
		Description: "前 3 名收入抽成配置",
		IsActive:    true,
	}

	require.NoError(t, repo.CreateConfig(ctx, cfg))
	require.NotZero(t, cfg.ID)

	got, err := repo.GetConfig(ctx, cfg.ID)
	require.NoError(t, err)
	assert.Equal(t, cfg.Name, got.Name)
	assert.Equal(t, cfg.RankingType, got.RankingType)
	assert.Equal(t, cfg.Month, got.Month)
	assert.Equal(t, cfg.RulesJSON, got.RulesJSON)

	// 未找到配置时应返回 ErrNotFound
	_, err = repo.GetConfig(ctx, cfg.ID+999)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestRankingCommissionRepository_GetActiveConfigForMonth(t *testing.T) {
	db := setupRankingDB(t)
	repo := NewRankingCommissionRepository(db)
	ctx := context.Background()

	rt := model.RankingTypeIncome

	active := &model.RankingCommissionConfig{
		Name:        "收入榜 2025-02",
		RankingType: rt,
		Period:      "monthly",
		Month:       "2025-02",
		RulesJSON:   "[]",
		IsActive:    true,
	}
	inactive := &model.RankingCommissionConfig{
		Name:        "收入榜 2025-02 旧版",
		RankingType: rt,
		Period:      "monthly",
		Month:       "2025-02",
		RulesJSON:   "[]",
		IsActive:    false,
	}

	require.NoError(t, repo.CreateConfig(ctx, inactive))
	// 显式将数据库中的 is_active 设置为 false，避免受默认值影响
	require.NoError(t, db.Model(&model.RankingCommissionConfig{}).
		Where("id = ?", inactive.ID).
		Update("is_active", false).Error)
	require.NoError(t, repo.CreateConfig(ctx, active))

	got, err := repo.GetActiveConfigForMonth(ctx, rt, "2025-02")
	require.NoError(t, err)
	assert.Equal(t, active.ID, got.ID)

	// 不存在激活配置时返回 ErrNotFound
	_, err = repo.GetActiveConfigForMonth(ctx, rt, "2099-01")
	require.ErrorIs(t, err, ErrNotFound)
}

func TestRankingCommissionRepository_ListConfigs_FilterAndPagination(t *testing.T) {
	db := setupRankingDB(t)
	repo := NewRankingCommissionRepository(db)
	ctx := context.Background()

	rtIncome := model.RankingTypeIncome
	rtOrder := model.RankingTypeOrderCount

	monthA := "2025-01"
	monthB := "2025-02"

	// 插入多条配置，用于测试过滤与分页
	var inserted []*model.RankingCommissionConfig
	for i := 0; i < 15; i++ {
		cfg := &model.RankingCommissionConfig{
			Name:        "配置 A",
			RankingType: rtIncome,
			Period:      "monthly",
			Month:       monthA,
			RulesJSON:   "[]",
			IsActive:    true,
		}
		if i%2 == 0 {
			cfg.RankingType = rtOrder
		}
		if i%3 == 0 {
			cfg.Month = monthB
		}
		require.NoError(t, repo.CreateConfig(ctx, cfg))
		inserted = append(inserted, cfg)
	}

	t.Run("按类型+月份+激活状态过滤", func(t *testing.T) {
		opts := RankingCommissionConfigListOptions{
			RankingType: &rtIncome,
			Month:       &monthA,
		}
		list, total, err := repo.ListConfigs(ctx, opts)
		require.NoError(t, err)
		assert.Greater(t, total, int64(0))
		for _, c := range list {
			assert.Equal(t, rtIncome, c.RankingType)
			assert.Equal(t, monthA, c.Month)
		}
	})

	t.Run("默认分页参数", func(t *testing.T) {
		opts := RankingCommissionConfigListOptions{}
		list, total, err := repo.ListConfigs(ctx, opts)
		require.NoError(t, err)
		assert.Equal(t, total, int64(len(list)))
		// 默认 PageSize=20，因此数据量较小时应一次性返回
		assert.LessOrEqual(t, len(list), 20)
	})

	t.Run("自定义分页", func(t *testing.T) {
		opts := RankingCommissionConfigListOptions{Page: 2, PageSize: 5}
		list, total, err := repo.ListConfigs(ctx, opts)
		require.NoError(t, err)
		assert.Equal(t, int64(len(inserted)), total)
		assert.Len(t, list, 5)
	})
}

func TestRankingCommissionRepository_UpdateAndDeleteConfig(t *testing.T) {
	db := setupRankingDB(t)
	repo := NewRankingCommissionRepository(db)
	ctx := context.Background()

	cfg := &model.RankingCommissionConfig{
		Name:        "更新前配置",
		RankingType: model.RankingTypeQuality,
		Period:      "monthly",
		Month:       "2025-03",
		RulesJSON:   "[]",
		IsActive:    true,
	}
	require.NoError(t, repo.CreateConfig(ctx, cfg))

	// 更新配置内容
	cfg.Name = "更新后配置"
	cfg.RulesJSON = `[{"rankStart":1,"rankEnd":5,"commissionRate":15}]`
	cfg.IsActive = false
	require.NoError(t, repo.UpdateConfig(ctx, cfg))

	got, err := repo.GetConfig(ctx, cfg.ID)
	require.NoError(t, err)
	assert.Equal(t, "更新后配置", got.Name)
	assert.Equal(t, cfg.RulesJSON, got.RulesJSON)
	assert.False(t, got.IsActive)

	// 删除配置
	require.NoError(t, repo.DeleteConfig(ctx, cfg.ID))
	_, err = repo.GetConfig(ctx, cfg.ID)
	require.ErrorIs(t, err, ErrNotFound)
}
