package ranking

import (
    "context"
    "testing"

    "github.com/glebarez/sqlite"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"

    "gamelink/internal/model"
)

func setupCommissionDB(t *testing.T) *gorm.DB {
    t.Helper()
    db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
    require.NoError(t, err)
    require.NoError(t, db.AutoMigrate(&model.RankingCommissionConfig{}))
    return db
}

func TestRankingCommissionRepository_CreateAndGet(t *testing.T) {
    db := setupCommissionDB(t)
    repo := NewRankingCommissionRepository(db)
    ctx := context.Background()

    cfg := &model.RankingCommissionConfig{Name: "月度收入", RankingType: model.RankingTypeIncome, Period: "monthly", Month: "2025-01", RulesJSON: "[]", IsActive: true}
    require.NoError(t, repo.CreateConfig(ctx, cfg))
    require.NotZero(t, cfg.ID)

    got, err := repo.GetConfig(ctx, cfg.ID)
    require.NoError(t, err)
    assert.Equal(t, cfg.Name, got.Name)

    _, err = repo.GetConfig(ctx, 999)
    require.ErrorIs(t, err, ErrNotFound)
}

func TestRankingCommissionRepository_ListFiltersAndPagination(t *testing.T) {
    db := setupCommissionDB(t)
    repo := NewRankingCommissionRepository(db)
    ctx := context.Background()

    for i := 1; i <= 25; i++ {
        rt := model.RankingTypeIncome
        if i%2 == 0 { rt = model.RankingTypeOrderCount }
        cfg := &model.RankingCommissionConfig{Name: "cfg", RankingType: rt, Period: "monthly", Month: "2025-01", RulesJSON: "[]", IsActive: i%3 != 0}
        require.NoError(t, repo.CreateConfig(ctx, cfg))
    }

    page := 2
    pageSize := 10
    list, total, err := repo.ListConfigs(ctx, RankingCommissionConfigListOptions{Page: page, PageSize: pageSize})
    require.NoError(t, err)
    assert.Equal(t, int64(25), total)
    require.Len(t, list, 10)

    rt := model.RankingTypeOrderCount
    active := true
    filtered, total2, err := repo.ListConfigs(ctx, RankingCommissionConfigListOptions{RankingType: &rt, Month: strPtr("2025-01"), IsActive: &active, Page: 1, PageSize: 100})
    require.NoError(t, err)
    assert.GreaterOrEqual(t, total2, int64(1))
    for _, c := range filtered { assert.Equal(t, rt, c.RankingType); assert.Equal(t, "2025-01", c.Month); assert.True(t, c.IsActive) }
}

func TestRankingCommissionRepository_UpdateAndDelete(t *testing.T) {
    db := setupCommissionDB(t)
    repo := NewRankingCommissionRepository(db)
    ctx := context.Background()

    cfg := &model.RankingCommissionConfig{Name: "old", RankingType: model.RankingTypeIncome, Period: "monthly", Month: "2025-02", RulesJSON: "[]", IsActive: true}
    require.NoError(t, repo.CreateConfig(ctx, cfg))

    newName := "new"
    cfg.Name = newName
    cfg.IsActive = false
    require.NoError(t, repo.UpdateConfig(ctx, cfg))

    got, err := repo.GetConfig(ctx, cfg.ID)
    require.NoError(t, err)
    assert.Equal(t, newName, got.Name)
    assert.False(t, got.IsActive)

    require.NoError(t, repo.DeleteConfig(ctx, cfg.ID))
    _, err = repo.GetConfig(ctx, cfg.ID)
    require.ErrorIs(t, err, ErrNotFound)
}

func TestRankingCommissionRepository_GetActiveConfigForMonth(t *testing.T) {
    db := setupCommissionDB(t)
    repo := NewRankingCommissionRepository(db)
    ctx := context.Background()

    rt := model.RankingTypeIncome
    cfg1 := &model.RankingCommissionConfig{Name: "inactive", RankingType: rt, Period: "monthly", Month: "2025-03", RulesJSON: "[]", IsActive: false}
    cfg2 := &model.RankingCommissionConfig{Name: "active", RankingType: rt, Period: "monthly", Month: "2025-03", RulesJSON: "[]", IsActive: true}
    require.NoError(t, repo.CreateConfig(ctx, cfg1))
    require.NoError(t, repo.CreateConfig(ctx, cfg2))

    got, err := repo.GetActiveConfigForMonth(ctx, rt, "2025-03")
    require.NoError(t, err)
    assert.Equal(t, "active", got.Name)

    _, err = repo.GetActiveConfigForMonth(ctx, model.RankingTypeOrderCount, "2025-03")
    require.ErrorIs(t, err, ErrNotFound)
}

