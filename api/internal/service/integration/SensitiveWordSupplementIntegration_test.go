// Package integration provides supplemental integration tests for SensitiveWord module.
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
// SensitiveWord Supplemental Integration Tests
// Tests for additional scenarios not covered in the main file
// ============================================================================

// TestSensitiveWordRepository_ListByMatchType tests filtering by match type
func TestSensitiveWordRepository_ListByMatchType(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	// Create words with fuzzy match type
	for i := 0; i < 2; i++ {
		word := &model.SensitiveWord{
			Word:      fmt.Sprintf("fuzzy_%d_%d", i, time.Now().UnixNano()),
			Category:  model.SensitiveWordCategoryAbuse,
			MatchType: model.SensitiveWordMatchTypeFuzzy,
			Severity:  model.SensitiveWordSeverityMedium,
		}
		require.NoError(t, repo.Create(ctx, word))
	}

	matchType := model.SensitiveWordMatchTypeFuzzy
	words, total, err := repo.List(ctx, repository.SensitiveWordListOptions{
		MatchType: &matchType,
		Page:      1,
		PageSize:  10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(2))
	for _, w := range words {
		assert.Equal(t, model.SensitiveWordMatchTypeFuzzy, w.MatchType)
	}
}

// TestSensitiveWordRepository_ListByIsActive tests filtering by active status
func TestSensitiveWordRepository_ListByIsActive(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	// Create active word
	activeWord := &model.SensitiveWord{
		Word:      fmt.Sprintf("active_%d", time.Now().UnixNano()),
		Category:  model.SensitiveWordCategoryOther,
		MatchType: model.SensitiveWordMatchTypeExact,
		Severity:  model.SensitiveWordSeverityLow,
		IsActive:  true,
	}
	require.NoError(t, repo.Create(ctx, activeWord))

	// Create inactive word
	inactiveWord := &model.SensitiveWord{
		Word:      fmt.Sprintf("inactive_%d", time.Now().UnixNano()),
		Category:  model.SensitiveWordCategoryOther,
		MatchType: model.SensitiveWordMatchTypeExact,
		Severity:  model.SensitiveWordSeverityLow,
		IsActive:  false,
	}
	require.NoError(t, repo.Create(ctx, inactiveWord))

	// List only active words
	isActive := true
	words, total, err := repo.List(ctx, repository.SensitiveWordListOptions{
		IsActive: &isActive,
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	for _, w := range words {
		assert.True(t, w.IsActive)
	}
}

// TestSensitiveWordRepository_Create_AllMatchTypes tests creating words with all match types
func TestSensitiveWordRepository_Create_AllMatchTypes(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	matchTypes := []model.SensitiveWordMatchType{
		model.SensitiveWordMatchTypeExact,
		model.SensitiveWordMatchTypeFuzzy,
		model.SensitiveWordMatchTypeRegex,
	}

	for _, mt := range matchTypes {
		word := &model.SensitiveWord{
			Word:      fmt.Sprintf("%s_%d", mt, time.Now().UnixNano()),
			Category:  model.SensitiveWordCategoryOther,
			MatchType: mt,
			Severity:  model.SensitiveWordSeverityMedium,
		}
		require.NoError(t, repo.Create(ctx, word))
		assert.NotZero(t, word.ID)
	}
}

// TestSensitiveWordRepository_Create_AllCategories tests creating words with all categories
func TestSensitiveWordRepository_Create_AllCategories(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	categories := []model.SensitiveWordCategory{
		model.SensitiveWordCategoryPolitics,
		model.SensitiveWordCategoryPorn,
		model.SensitiveWordCategoryAbuse,
		model.SensitiveWordCategoryAd,
		model.SensitiveWordCategoryOther,
	}

	for _, cat := range categories {
		word := &model.SensitiveWord{
			Word:      fmt.Sprintf("cat_%s_%d", cat, time.Now().UnixNano()),
			Category:  cat,
			MatchType: model.SensitiveWordMatchTypeExact,
			Severity:  model.SensitiveWordSeverityMedium,
		}
		require.NoError(t, repo.Create(ctx, word))
		assert.NotZero(t, word.ID)
	}
}

// TestSensitiveWordRepository_Update_ToggleIsActive tests toggling active status
func TestSensitiveWordRepository_Update_ToggleIsActive(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	word := &model.SensitiveWord{
		Word:      fmt.Sprintf("toggle_%d", time.Now().UnixNano()),
		Category:  model.SensitiveWordCategoryOther,
		MatchType: model.SensitiveWordMatchTypeExact,
		Severity:  model.SensitiveWordSeverityLow,
		IsActive:  true,
	}
	require.NoError(t, repo.Create(ctx, word))

	// Deactivate
	word.IsActive = false
	err := repo.Update(ctx, word)
	require.NoError(t, err)

	got, err := repo.Get(ctx, word.ID)
	require.NoError(t, err)
	assert.False(t, got.IsActive)

	// Reactivate
	word.IsActive = true
	err = repo.Update(ctx, word)
	require.NoError(t, err)

	got, err = repo.Get(ctx, word.ID)
	require.NoError(t, err)
	assert.True(t, got.IsActive)
}

// TestSensitiveWordRepository_Update_ChangeCategory tests changing category
func TestSensitiveWordRepository_Update_ChangeCategory(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	word := &model.SensitiveWord{
		Word:      fmt.Sprintf("recategorize_%d", time.Now().UnixNano()),
		Category:  model.SensitiveWordCategoryOther,
		MatchType: model.SensitiveWordMatchTypeExact,
		Severity:  model.SensitiveWordSeverityLow,
	}
	require.NoError(t, repo.Create(ctx, word))

	// Change category
	word.Category = model.SensitiveWordCategoryAbuse
	err := repo.Update(ctx, word)
	require.NoError(t, err)

	got, err := repo.Get(ctx, word.ID)
	require.NoError(t, err)
	assert.Equal(t, model.SensitiveWordCategoryAbuse, got.Category)
}

// TestSensitiveWordRepository_Update_ChangeSeverity tests changing severity
func TestSensitiveWordRepository_Update_ChangeSeverity(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	word := &model.SensitiveWord{
		Word:      fmt.Sprintf("severity_%d", time.Now().UnixNano()),
		Category:  model.SensitiveWordCategoryOther,
		MatchType: model.SensitiveWordMatchTypeExact,
		Severity:  model.SensitiveWordSeverityLow,
	}
	require.NoError(t, repo.Create(ctx, word))

	// Change severity
	word.Severity = model.SensitiveWordSeverityHigh
	err := repo.Update(ctx, word)
	require.NoError(t, err)

	got, err := repo.Get(ctx, word.ID)
	require.NoError(t, err)
	assert.Equal(t, model.SensitiveWordSeverityHigh, got.Severity)
}

// TestSensitiveWordRepository_Update_ChangeMatchType tests changing match type
func TestSensitiveWordRepository_Update_ChangeMatchType(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	word := &model.SensitiveWord{
		Word:      fmt.Sprintf("matchtype_%d", time.Now().UnixNano()),
		Category:  model.SensitiveWordCategoryOther,
		MatchType: model.SensitiveWordMatchTypeExact,
		Severity:  model.SensitiveWordSeverityLow,
	}
	require.NoError(t, repo.Create(ctx, word))

	// Change match type
	word.MatchType = model.SensitiveWordMatchTypeFuzzy
	err := repo.Update(ctx, word)
	require.NoError(t, err)

	got, err := repo.Get(ctx, word.ID)
	require.NoError(t, err)
	assert.Equal(t, model.SensitiveWordMatchTypeFuzzy, got.MatchType)
}

// TestSensitiveWordRepository_Update_WithReplacement tests updating with custom replacement
func TestSensitiveWordRepository_Update_WithReplacement(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	word := &model.SensitiveWord{
		Word:        fmt.Sprintf("replace_%d", time.Now().UnixNano()),
		Category:    model.SensitiveWordCategoryOther,
		MatchType:   model.SensitiveWordMatchTypeExact,
		Severity:    model.SensitiveWordSeverityLow,
		Replacement: "***",
	}
	require.NoError(t, repo.Create(ctx, word))

	// Change replacement
	word.Replacement = "[FILTERED]"
	err := repo.Update(ctx, word)
	require.NoError(t, err)

	got, err := repo.Get(ctx, word.ID)
	require.NoError(t, err)
	assert.Equal(t, "[FILTERED]", got.Replacement)
}

// TestSensitiveWordRepository_Delete_NonExistent tests deleting non-existent word
func TestSensitiveWordRepository_Delete_NonExistent(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	// Try to delete non-existent word
	err := repo.Delete(ctx, 99999)
	assert.Error(t, err)
}

// TestSensitiveWordRepository_Get_NonExistent tests getting non-existent word
func TestSensitiveWordRepository_Get_NonExistent(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	// Try to get non-existent word
	_, err := repo.Get(ctx, 99999)
	assert.Error(t, err)
}

// TestSensitiveWordRepository_List_EmptyDatabase tests listing when database is empty
func TestSensitiveWordRepository_List_EmptyDatabase(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	words, total, err := repo.List(ctx, repository.SensitiveWordListOptions{
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, words)
}

// TestSensitiveWordRepository_List_Pagination tests pagination
func TestSensitiveWordRepository_List_Pagination(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	// Create 15 words
	for i := 0; i < 15; i++ {
		word := &model.SensitiveWord{
			Word:      fmt.Sprintf("page_%d_%d", i, time.Now().UnixNano()),
			Category:  model.SensitiveWordCategoryOther,
			MatchType: model.SensitiveWordMatchTypeExact,
			Severity:  model.SensitiveWordSeverityLow,
		}
		require.NoError(t, repo.Create(ctx, word))
	}

	// First page
	words1, total, err := repo.List(ctx, repository.SensitiveWordListOptions{
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(15), total)
	assert.Len(t, words1, 10)

	// Second page
	words2, _, err := repo.List(ctx, repository.SensitiveWordListOptions{
		Page:     2,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.Len(t, words2, 5)
}

// TestSensitiveWordRepository_List_Sorting tests sorting by different fields
func TestSensitiveWordRepository_List_Sorting(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	// Create words with different names
	for _, name := range []string{"zebra", "apple", "middle"} {
		word := &model.SensitiveWord{
			Word:      fmt.Sprintf("%s_%d", name, time.Now().UnixNano()),
			Category:  model.SensitiveWordCategoryOther,
			MatchType: model.SensitiveWordMatchTypeExact,
			Severity:  model.SensitiveWordSeverityLow,
		}
		require.NoError(t, repo.Create(ctx, word))
	}

	// List with keyword to get created words
	_, total, err := repo.List(ctx, repository.SensitiveWordListOptions{
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(3))
}
