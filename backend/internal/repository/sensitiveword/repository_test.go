package sensitiveword

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/testutil"
)

func TestSensitiveWordRepository_Create(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.SensitiveWord{})
	repo := NewSensitiveWordRepository(db)
	ctx := context.Background()

	t.Run("create new sensitive word", func(t *testing.T) {
		word := &model.SensitiveWord{
			Word:     "测试敏感词",
			Category: model.SensitiveWordCategoryOther,
			Severity: model.SensitiveWordSeverityLow,
		}

		err := repo.Create(ctx, word)
		require.NoError(t, err)
		assert.NotZero(t, word.ID)
	})

	t.Run("create duplicate sensitive word", func(t *testing.T) {
		word1 := &model.SensitiveWord{
			Word:     "重复敏感词",
			Category: model.SensitiveWordCategoryOther,
			Severity: model.SensitiveWordSeverityLow,
		}
		err := repo.Create(ctx, word1)
		require.NoError(t, err)

		word2 := &model.SensitiveWord{
			Word:     "重复敏感词",
			Category: model.SensitiveWordCategoryPolitical,
			Severity: model.SensitiveWordSeverityHigh,
		}
		err = repo.Create(ctx, word2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "已存在")
	})
}

func TestSensitiveWordRepository_Get(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.SensitiveWord{})
	repo := NewSensitiveWordRepository(db)
	ctx := context.Background()

	t.Run("get existing sensitive word", func(t *testing.T) {
		word := &model.SensitiveWord{
			Word:     "测试词",
			Category: model.SensitiveWordCategoryOther,
			Severity: model.SensitiveWordSeverityMedium,
		}
		err := repo.Create(ctx, word)
		require.NoError(t, err)

		retrieved, err := repo.Get(ctx, word.ID)
		require.NoError(t, err)
		assert.Equal(t, word.Word, retrieved.Word)
		assert.Equal(t, word.Category, retrieved.Category)
		assert.Equal(t, word.Severity, retrieved.Severity)
	})

	t.Run("get non-existent sensitive word", func(t *testing.T) {
		_, err := repo.Get(ctx, 99999)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})
}

func TestSensitiveWordRepository_List(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.SensitiveWord{})
	repo := NewSensitiveWordRepository(db)
	ctx := context.Background()

	// 创建测试数据
	words := []model.SensitiveWord{
		{Word: "政治词汇", Category: model.SensitiveWordCategoryPolitical, Severity: model.SensitiveWordSeverityHigh},
		{Word: "色情词汇", Category: model.SensitiveWordCategoryPornographic, Severity: model.SensitiveWordSeverityHigh},
		{Word: "暴力词汇", Category: model.SensitiveWordCategoryViolent, Severity: model.SensitiveWordSeverityMedium},
		{Word: "广告词汇", Category: model.SensitiveWordCategoryAdvertising, Severity: model.SensitiveWordSeverityLow},
	}
	for i := range words {
		err := repo.Create(ctx, &words[i])
		require.NoError(t, err)
	}

	t.Run("list all sensitive words", func(t *testing.T) {
		result, total, err := repo.List(ctx, repository.SensitiveWordListOptions{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(4), total)
		assert.Len(t, result, 4)
	})

	t.Run("list with keyword filter", func(t *testing.T) {
		result, total, err := repo.List(ctx, repository.SensitiveWordListOptions{
			Page:     1,
			PageSize: 10,
			Keyword:  "政治",
		})
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, result, 1)
		assert.Equal(t, "政治词汇", result[0].Word)
	})

	t.Run("list with category filter", func(t *testing.T) {
		category := model.SensitiveWordCategoryPolitical
		result, total, err := repo.List(ctx, repository.SensitiveWordListOptions{
			Page:     1,
			PageSize: 10,
			Category: &category,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, result, 1)
		assert.Equal(t, model.SensitiveWordCategoryPolitical, result[0].Category)
	})

	t.Run("list with severity filter", func(t *testing.T) {
		severity := model.SensitiveWordSeverityHigh
		result, total, err := repo.List(ctx, repository.SensitiveWordListOptions{
			Page:     1,
			PageSize: 10,
			Severity: &severity,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, result, 2)
	})

	t.Run("list with pagination", func(t *testing.T) {
		result, total, err := repo.List(ctx, repository.SensitiveWordListOptions{
			Page:     1,
			PageSize: 2,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(4), total)
		assert.Len(t, result, 2)

		result, total, err = repo.List(ctx, repository.SensitiveWordListOptions{
			Page:     2,
			PageSize: 2,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(4), total)
		assert.Len(t, result, 2)
	})
}

func TestSensitiveWordRepository_Update(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.SensitiveWord{})
	repo := NewSensitiveWordRepository(db)
	ctx := context.Background()

	t.Run("update existing sensitive word", func(t *testing.T) {
		word := &model.SensitiveWord{
			Word:     "原始词",
			Category: model.SensitiveWordCategoryOther,
			Severity: model.SensitiveWordSeverityLow,
		}
		err := repo.Create(ctx, word)
		require.NoError(t, err)

		word.Word = "更新词"
		word.Category = model.SensitiveWordCategoryPolitical
		word.Severity = model.SensitiveWordSeverityHigh
		err = repo.Update(ctx, word)
		require.NoError(t, err)

		retrieved, err := repo.Get(ctx, word.ID)
		require.NoError(t, err)
		assert.Equal(t, "更新词", retrieved.Word)
		assert.Equal(t, model.SensitiveWordCategoryPolitical, retrieved.Category)
		assert.Equal(t, model.SensitiveWordSeverityHigh, retrieved.Severity)
	})

	t.Run("update to duplicate word", func(t *testing.T) {
		word1 := &model.SensitiveWord{
			Word:     "词汇1",
			Category: model.SensitiveWordCategoryOther,
			Severity: model.SensitiveWordSeverityLow,
		}
		err := repo.Create(ctx, word1)
		require.NoError(t, err)

		word2 := &model.SensitiveWord{
			Word:     "词汇2",
			Category: model.SensitiveWordCategoryOther,
			Severity: model.SensitiveWordSeverityLow,
		}
		err = repo.Create(ctx, word2)
		require.NoError(t, err)

		word2.Word = "词汇1"
		err = repo.Update(ctx, word2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "已存在")
	})

	t.Run("update non-existent sensitive word", func(t *testing.T) {
		word := &model.SensitiveWord{
			Base:     model.Base{ID: 99999},
			Word:     "不存在",
			Category: model.SensitiveWordCategoryOther,
			Severity: model.SensitiveWordSeverityLow,
		}
		err := repo.Update(ctx, word)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})
}

func TestSensitiveWordRepository_Delete(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.SensitiveWord{})
	repo := NewSensitiveWordRepository(db)
	ctx := context.Background()

	t.Run("delete existing sensitive word", func(t *testing.T) {
		word := &model.SensitiveWord{
			Word:     "待删除词",
			Category: model.SensitiveWordCategoryOther,
			Severity: model.SensitiveWordSeverityLow,
		}
		err := repo.Create(ctx, word)
		require.NoError(t, err)

		err = repo.Delete(ctx, word.ID)
		require.NoError(t, err)

		_, err = repo.Get(ctx, word.ID)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})

	t.Run("delete non-existent sensitive word", func(t *testing.T) {
		err := repo.Delete(ctx, 99999)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})
}

func TestSensitiveWordRepository_GetAll(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.SensitiveWord{})
	repo := NewSensitiveWordRepository(db)
	ctx := context.Background()

	t.Run("get all sensitive words", func(t *testing.T) {
		words := []model.SensitiveWord{
			{Word: "词1", Category: model.SensitiveWordCategoryOther, Severity: model.SensitiveWordSeverityLow},
			{Word: "词2", Category: model.SensitiveWordCategoryOther, Severity: model.SensitiveWordSeverityMedium},
			{Word: "词3", Category: model.SensitiveWordCategoryOther, Severity: model.SensitiveWordSeverityHigh},
		}
		for i := range words {
			err := repo.Create(ctx, &words[i])
			require.NoError(t, err)
		}

		result, err := repo.GetAll(ctx)
		require.NoError(t, err)
		assert.Len(t, result, 3)
	})

	t.Run("get all when empty", func(t *testing.T) {
		db2 := testutil.NewMemoryDB(t)
		testutil.MigrateTables(t, db2, &model.SensitiveWord{})
		repo2 := NewSensitiveWordRepository(db2)

		result, err := repo2.GetAll(ctx)
		require.NoError(t, err)
		assert.Empty(t, result)
	})
}
