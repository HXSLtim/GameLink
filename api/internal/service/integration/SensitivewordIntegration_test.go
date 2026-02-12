package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/sensitiveword"
)

// ============================================================================
// SensitiveWord CRUD Tests
// ============================================================================

func TestSensitiveWordRepository_Create(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	word := &model.SensitiveWord{
		Word:     fmt.Sprintf("badword_%d", time.Now().UnixNano()),
		Category: model.SensitiveWordCategoryAbuse,
		Severity: model.SensitiveWordSeverityMedium,
	}

	err := repo.Create(ctx, word)
	require.NoError(t, err)
	assert.NotZero(t, word.ID)
}

func TestSensitiveWordRepository_Create_Duplicate(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	wordText := fmt.Sprintf("duplicate_%d", time.Now().UnixNano())
	word1 := &model.SensitiveWord{
		Word:     wordText,
		Category: model.SensitiveWordCategoryAbuse,
		Severity: model.SensitiveWordSeverityLow,
	}
	require.NoError(t, repo.Create(ctx, word1))

	// Try to create duplicate
	word2 := &model.SensitiveWord{
		Word:     wordText,
		Category: model.SensitiveWordCategoryAbuse,
		Severity: model.SensitiveWordSeverityLow,
	}
	err := repo.Create(ctx, word2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "已存在")
}

func TestSensitiveWordRepository_Get(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	word := &model.SensitiveWord{
		Word:     fmt.Sprintf("getword_%d", time.Now().UnixNano()),
		Category: model.SensitiveWordCategoryPorn,
		Severity: model.SensitiveWordSeverityHigh,
	}
	require.NoError(t, repo.Create(ctx, word))

	got, err := repo.Get(ctx, word.ID)
	require.NoError(t, err)
	assert.Equal(t, word.Word, got.Word)
	assert.Equal(t, model.SensitiveWordCategoryPorn, got.Category)
	assert.Equal(t, model.SensitiveWordSeverityHigh, got.Severity)
}

func TestSensitiveWordRepository_Update(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	word := &model.SensitiveWord{
		Word:     fmt.Sprintf("updateword_%d", time.Now().UnixNano()),
		Category: model.SensitiveWordCategoryAd,
		Severity: model.SensitiveWordSeverityLow,
	}
	require.NoError(t, repo.Create(ctx, word))

	// Update
	word.Category = model.SensitiveWordCategoryAbuse
	word.Severity = model.SensitiveWordSeverityMedium
	err := repo.Update(ctx, word)
	require.NoError(t, err)

	got, err := repo.Get(ctx, word.ID)
	require.NoError(t, err)
	assert.Equal(t, model.SensitiveWordCategoryAbuse, got.Category)
	assert.Equal(t, model.SensitiveWordSeverityMedium, got.Severity)
}

func TestSensitiveWordRepository_Delete(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	word := &model.SensitiveWord{
		Word:     fmt.Sprintf("deleteword_%d", time.Now().UnixNano()),
		Category: model.SensitiveWordCategoryOther,
		Severity: model.SensitiveWordSeverityLow,
	}
	require.NoError(t, repo.Create(ctx, word))

	err := repo.Delete(ctx, word.ID)
	require.NoError(t, err)

	_, err = repo.Get(ctx, word.ID)
	assert.Error(t, err)
}

func TestSensitiveWordRepository_List(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	// Create multiple words
	for i := 0; i < 3; i++ {
		word := &model.SensitiveWord{
			Word:     fmt.Sprintf("listword_%d_%d", i, time.Now().UnixNano()),
			Category: model.SensitiveWordCategoryAbuse,
			Severity: model.SensitiveWordSeverityMedium,
		}
		require.NoError(t, repo.Create(ctx, word))
	}

	words, total, err := repo.List(ctx, repository.SensitiveWordListOptions{
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(3))
	assert.GreaterOrEqual(t, len(words), 3)
}

func TestSensitiveWordRepository_ListByCategory(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	// Create words with specific category
	for i := 0; i < 2; i++ {
		word := &model.SensitiveWord{
			Word:     fmt.Sprintf("politics_%d_%d", i, time.Now().UnixNano()),
			Category: model.SensitiveWordCategoryPolitics,
			Severity: model.SensitiveWordSeverityHigh,
		}
		require.NoError(t, repo.Create(ctx, word))
	}

	category := model.SensitiveWordCategoryPolitics
	words, total, err := repo.List(ctx, repository.SensitiveWordListOptions{
		Category: &category,
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(2))
	for _, w := range words {
		assert.Equal(t, model.SensitiveWordCategoryPolitics, w.Category)
	}
}

func TestSensitiveWordRepository_ListBySeverity(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	// Create words with high severity
	for i := 0; i < 2; i++ {
		word := &model.SensitiveWord{
			Word:     fmt.Sprintf("highsev_%d_%d", i, time.Now().UnixNano()),
			Category: model.SensitiveWordCategoryPorn,
			Severity: model.SensitiveWordSeverityHigh,
		}
		require.NoError(t, repo.Create(ctx, word))
	}

	severity := model.SensitiveWordSeverityHigh
	words, total, err := repo.List(ctx, repository.SensitiveWordListOptions{
		Severity: &severity,
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(2))
	for _, w := range words {
		assert.Equal(t, model.SensitiveWordSeverityHigh, w.Severity)
	}
}

func TestSensitiveWordRepository_ListByKeyword(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	uniquePrefix := fmt.Sprintf("unique_%d", time.Now().UnixNano())
	word := &model.SensitiveWord{
		Word:     uniquePrefix + "_searchable",
		Category: model.SensitiveWordCategoryOther,
		Severity: model.SensitiveWordSeverityLow,
	}
	require.NoError(t, repo.Create(ctx, word))

	words, total, err := repo.List(ctx, repository.SensitiveWordListOptions{
		Keyword:  uniquePrefix,
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	assert.GreaterOrEqual(t, len(words), 1)
}

func TestSensitiveWordRepository_GetAll(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	// Create some words
	for i := 0; i < 3; i++ {
		word := &model.SensitiveWord{
			Word:     fmt.Sprintf("allword_%d_%d", i, time.Now().UnixNano()),
			Category: model.SensitiveWordCategoryAbuse,
			Severity: model.SensitiveWordSeverityMedium,
		}
		require.NoError(t, repo.Create(ctx, word))
	}

	words, err := repo.GetAll(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(words), 3)
}
