package sensitiveword

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository/sensitiveword"
	"gamelink/pkg/testutil"
)

func TestSensitiveWordService_AddSensitiveWord(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.SensitiveWord{})
	repo := sensitiveword.NewSensitiveWordRepository(db)
	service := NewSensitiveWordService(repo)
	ctx := context.Background()

	t.Run("add valid sensitive word", func(t *testing.T) {
		req := AddSensitiveWordRequest{
			Word:     "测试敏感词",
			Category: model.SensitiveWordCategoryOther,
			Severity: model.SensitiveWordSeverityLow,
		}

		dto, err := service.AddSensitiveWord(ctx, req)
		require.NoError(t, err)
		assert.NotZero(t, dto.ID)
		assert.Equal(t, "测试敏感词", dto.Word)
	})

	t.Run("add duplicate sensitive word", func(t *testing.T) {
		req := AddSensitiveWordRequest{
			Word:     "重复词",
			Category: model.SensitiveWordCategoryOther,
			Severity: model.SensitiveWordSeverityLow,
		}

		_, err := service.AddSensitiveWord(ctx, req)
		require.NoError(t, err)

		_, err = service.AddSensitiveWord(ctx, req)
		assert.ErrorIs(t, err, ErrDuplicate)
	})

	t.Run("add with invalid category", func(t *testing.T) {
		req := AddSensitiveWordRequest{
			Word:     "测试",
			Category: model.SensitiveWordCategory("invalid"),
			Severity: model.SensitiveWordSeverityLow,
		}

		_, err := service.AddSensitiveWord(ctx, req)
		assert.ErrorIs(t, err, ErrValidation)
	})

	t.Run("add with empty word", func(t *testing.T) {
		req := AddSensitiveWordRequest{
			Word:     "   ",
			Category: model.SensitiveWordCategoryOther,
			Severity: model.SensitiveWordSeverityLow,
		}

		_, err := service.AddSensitiveWord(ctx, req)
		assert.ErrorIs(t, err, ErrValidation)
	})
}

func TestSensitiveWordService_UpdateSensitiveWord(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.SensitiveWord{})
	repo := sensitiveword.NewSensitiveWordRepository(db)
	service := NewSensitiveWordService(repo)
	ctx := context.Background()

	t.Run("update existing sensitive word", func(t *testing.T) {
		// 创建敏感词
		addReq := AddSensitiveWordRequest{
			Word:     "原始词",
			Category: model.SensitiveWordCategoryOther,
			Severity: model.SensitiveWordSeverityLow,
		}
		dto, err := service.AddSensitiveWord(ctx, addReq)
		require.NoError(t, err)

		// 更新敏感词
		updateReq := UpdateSensitiveWordRequest{
			Word:     "更新词",
			Category: model.SensitiveWordCategoryPolitical,
			Severity: model.SensitiveWordSeverityHigh,
		}
		err = service.UpdateSensitiveWord(ctx, dto.ID, updateReq)
		require.NoError(t, err)

		// 验证更新
		word, err := repo.Get(ctx, dto.ID)
		require.NoError(t, err)
		assert.Equal(t, "更新词", word.Word)
		assert.Equal(t, model.SensitiveWordCategoryPolitical, word.Category)
		assert.Equal(t, model.SensitiveWordSeverityHigh, word.Severity)
	})

	t.Run("update to duplicate word", func(t *testing.T) {
		// 创建两个敏感词
		_, err := service.AddSensitiveWord(ctx, AddSensitiveWordRequest{
			Word:     "词1",
			Category: model.SensitiveWordCategoryOther,
			Severity: model.SensitiveWordSeverityLow,
		})
		require.NoError(t, err)

		dto2, err := service.AddSensitiveWord(ctx, AddSensitiveWordRequest{
			Word:     "词2",
			Category: model.SensitiveWordCategoryOther,
			Severity: model.SensitiveWordSeverityLow,
		})
		require.NoError(t, err)

		// 尝试将词2更新为词1
		err = service.UpdateSensitiveWord(ctx, dto2.ID, UpdateSensitiveWordRequest{
			Word:     "词1",
			Category: model.SensitiveWordCategoryOther,
			Severity: model.SensitiveWordSeverityLow,
		})
		assert.ErrorIs(t, err, ErrDuplicate)
	})

	t.Run("update non-existent word", func(t *testing.T) {
		err := service.UpdateSensitiveWord(ctx, 99999, UpdateSensitiveWordRequest{
			Word:     "不存在",
			Category: model.SensitiveWordCategoryOther,
			Severity: model.SensitiveWordSeverityLow,
		})
		assert.ErrorIs(t, err, ErrNotFound)
	})
}

func TestSensitiveWordService_DeleteSensitiveWord(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.SensitiveWord{})
	repo := sensitiveword.NewSensitiveWordRepository(db)
	service := NewSensitiveWordService(repo)
	ctx := context.Background()

	t.Run("delete existing word", func(t *testing.T) {
		dto, err := service.AddSensitiveWord(ctx, AddSensitiveWordRequest{
			Word:     "待删除",
			Category: model.SensitiveWordCategoryOther,
			Severity: model.SensitiveWordSeverityLow,
		})
		require.NoError(t, err)

		err = service.DeleteSensitiveWord(ctx, dto.ID)
		require.NoError(t, err)

		_, err = repo.Get(ctx, dto.ID)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("delete non-existent word", func(t *testing.T) {
		err := service.DeleteSensitiveWord(ctx, 99999)
		assert.ErrorIs(t, err, ErrNotFound)
	})
}

func TestSensitiveWordService_ListSensitiveWords(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.SensitiveWord{})
	repo := sensitiveword.NewSensitiveWordRepository(db)
	service := NewSensitiveWordService(repo)
	ctx := context.Background()

	// 创建测试数据
	words := []AddSensitiveWordRequest{
		{Word: "政治词", Category: model.SensitiveWordCategoryPolitical, Severity: model.SensitiveWordSeverityHigh},
		{Word: "色情词", Category: model.SensitiveWordCategoryPornographic, Severity: model.SensitiveWordSeverityHigh},
		{Word: "暴力词", Category: model.SensitiveWordCategoryViolent, Severity: model.SensitiveWordSeverityMedium},
	}
	for _, w := range words {
		_, err := service.AddSensitiveWord(ctx, w)
		require.NoError(t, err)
	}

	t.Run("list all words", func(t *testing.T) {
		resp, err := service.ListSensitiveWords(ctx, ListSensitiveWordsRequest{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(3), resp.Total)
		assert.Len(t, resp.Words, 3)
	})

	t.Run("list with keyword", func(t *testing.T) {
		resp, err := service.ListSensitiveWords(ctx, ListSensitiveWordsRequest{
			Page:     1,
			PageSize: 10,
			Keyword:  "政治",
		})
		require.NoError(t, err)
		assert.Equal(t, int64(1), resp.Total)
		assert.Len(t, resp.Words, 1)
	})

	t.Run("list with category filter", func(t *testing.T) {
		category := model.SensitiveWordCategoryPolitical
		resp, err := service.ListSensitiveWords(ctx, ListSensitiveWordsRequest{
			Page:     1,
			PageSize: 10,
			Category: &category,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(1), resp.Total)
	})
}

func TestSensitiveWordService_DetectSensitiveWords(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.SensitiveWord{})
	repo := sensitiveword.NewSensitiveWordRepository(db)
	service := NewSensitiveWordService(repo)
	ctx := context.Background()

	// 创建测试敏感词
	_, err := service.AddSensitiveWord(ctx, AddSensitiveWordRequest{
		Word:     "敏感",
		Category: model.SensitiveWordCategoryOther,
		Severity: model.SensitiveWordSeverityHigh,
	})
	require.NoError(t, err)

	_, err = service.AddSensitiveWord(ctx, AddSensitiveWordRequest{
		Word:     "违规",
		Category: model.SensitiveWordCategoryOther,
		Severity: model.SensitiveWordSeverityMedium,
	})
	require.NoError(t, err)

	t.Run("detect sensitive words", func(t *testing.T) {
		resp, err := service.DetectSensitiveWords(ctx, DetectSensitiveWordsRequest{
			Content: "这是一段包含敏感词汇的文本",
		})
		require.NoError(t, err)
		assert.True(t, resp.HasSensitiveWords)
		assert.NotEmpty(t, resp.DetectedWords)
		assert.Contains(t, resp.HighlightedText, "***")
	})

	t.Run("detect multiple sensitive words", func(t *testing.T) {
		resp, err := service.DetectSensitiveWords(ctx, DetectSensitiveWordsRequest{
			Content: "这段文本包含敏感和违规内容",
		})
		require.NoError(t, err)
		assert.True(t, resp.HasSensitiveWords)
		assert.Len(t, resp.DetectedWords, 2)
	})

	t.Run("no sensitive words", func(t *testing.T) {
		resp, err := service.DetectSensitiveWords(ctx, DetectSensitiveWordsRequest{
			Content: "这是一段正常的文本",
		})
		require.NoError(t, err)
		assert.False(t, resp.HasSensitiveWords)
		assert.Empty(t, resp.DetectedWords)
		assert.Equal(t, "这是一段正常的文本", resp.HighlightedText)
	})

	t.Run("empty content", func(t *testing.T) {
		resp, err := service.DetectSensitiveWords(ctx, DetectSensitiveWordsRequest{
			Content: "",
		})
		require.NoError(t, err)
		assert.False(t, resp.HasSensitiveWords)
		assert.Empty(t, resp.DetectedWords)
	})

	t.Run("case insensitive detection", func(t *testing.T) {
		resp, err := service.DetectSensitiveWords(ctx, DetectSensitiveWordsRequest{
			Content: "这段文本包含敏感内容",
		})
		require.NoError(t, err)
		assert.True(t, resp.HasSensitiveWords)
	})
}

func TestSensitiveWordService_Cache(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.SensitiveWord{})
	repo := sensitiveword.NewSensitiveWordRepository(db)
	service := NewSensitiveWordService(repo)
	ctx := context.Background()

	// 添加敏感词
	_, err := service.AddSensitiveWord(ctx, AddSensitiveWordRequest{
		Word:     "缓存测试",
		Category: model.SensitiveWordCategoryOther,
		Severity: model.SensitiveWordSeverityLow,
	})
	require.NoError(t, err)

	// 第一次检测（会加载到缓存）
	resp1, err := service.DetectSensitiveWords(ctx, DetectSensitiveWordsRequest{
		Content: "这是缓存测试",
	})
	require.NoError(t, err)
	assert.True(t, resp1.HasSensitiveWords)

	// 第二次检测（应该使用缓存）
	resp2, err := service.DetectSensitiveWords(ctx, DetectSensitiveWordsRequest{
		Content: "这是缓存测试",
	})
	require.NoError(t, err)
	assert.True(t, resp2.HasSensitiveWords)

	// 添加新敏感词（会清除缓存）
	_, err = service.AddSensitiveWord(ctx, AddSensitiveWordRequest{
		Word:     "新词",
		Category: model.SensitiveWordCategoryOther,
		Severity: model.SensitiveWordSeverityLow,
	})
	require.NoError(t, err)

	// 检测新词（应该能检测到）
	resp3, err := service.DetectSensitiveWords(ctx, DetectSensitiveWordsRequest{
		Content: "这是新词测试",
	})
	require.NoError(t, err)
	assert.True(t, resp3.HasSensitiveWords)
}
