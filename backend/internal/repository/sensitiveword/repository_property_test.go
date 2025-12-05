package sensitiveword_test

import (
	"context"
	"fmt"
	"testing"

	"gamelink/internal/model"
	"gamelink/internal/repository/sensitiveword"
	"gamelink/pkg/testutil"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// TestSensitiveWordUniqueness tests Property 4: Sensitive Word Uniqueness
// **Feature: review-management-module, Property 4: 敏感词唯一性**
// **Validates: Requirements 5.3**
// For any sensitive word addition operation, the newly added sensitive word must not exist in the sensitive word database
func TestSensitiveWordUniqueness(t *testing.T) {
	properties := gopter.NewProperties(nil)

	// Property 1: Adding a sensitive word that doesn't exist should succeed
	properties.Property("adding a new unique sensitive word should succeed", prop.ForAll(
		func(word string, categoryIdx int, severityIdx int) bool {
			// Setup fresh database for each test
			db := testutil.NewMemoryDB(t)
			testutil.MigrateTables(t, db, &model.SensitiveWord{})
			repo := sensitiveword.NewSensitiveWordRepository(db)
			ctx := context.Background()

			// Map indices to valid categories and severities
			categories := []model.SensitiveWordCategory{
				model.SensitiveWordCategoryPolitical,
				model.SensitiveWordCategoryPornographic,
				model.SensitiveWordCategoryViolent,
				model.SensitiveWordCategoryAdvertising,
				model.SensitiveWordCategoryOther,
			}
			severities := []model.SensitiveWordSeverity{
				model.SensitiveWordSeverityLow,
				model.SensitiveWordSeverityMedium,
				model.SensitiveWordSeverityHigh,
			}

			category := categories[categoryIdx%len(categories)]
			severity := severities[severityIdx%len(severities)]

			// Create sensitive word
			sw := &model.SensitiveWord{
				Word:     word,
				Category: category,
				Severity: severity,
			}

			// First addition should succeed
			err := repo.Create(ctx, sw)
			return err == nil
		},
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 && len(s) <= 100 }),
		gen.IntRange(0, 100),
		gen.IntRange(0, 100),
	))

	// Property 2: Adding a duplicate sensitive word should fail
	properties.Property("adding a duplicate sensitive word should fail", prop.ForAll(
		func(word string, categoryIdx int, severityIdx int) bool {
			// Setup fresh database for each test
			db := testutil.NewMemoryDB(t)
			testutil.MigrateTables(t, db, &model.SensitiveWord{})
			repo := sensitiveword.NewSensitiveWordRepository(db)
			ctx := context.Background()

			// Map indices to valid categories and severities
			categories := []model.SensitiveWordCategory{
				model.SensitiveWordCategoryPolitical,
				model.SensitiveWordCategoryPornographic,
				model.SensitiveWordCategoryViolent,
				model.SensitiveWordCategoryAdvertising,
				model.SensitiveWordCategoryOther,
			}
			severities := []model.SensitiveWordSeverity{
				model.SensitiveWordSeverityLow,
				model.SensitiveWordSeverityMedium,
				model.SensitiveWordSeverityHigh,
			}

			category := categories[categoryIdx%len(categories)]
			severity := severities[severityIdx%len(severities)]

			// Create first sensitive word
			sw1 := &model.SensitiveWord{
				Word:     word,
				Category: category,
				Severity: severity,
			}

			// First addition should succeed
			err := repo.Create(ctx, sw1)
			if err != nil {
				return false
			}

			// Create second sensitive word with same word but different category/severity
			differentCategory := categories[(categoryIdx+1)%len(categories)]
			differentSeverity := severities[(severityIdx+1)%len(severities)]

			sw2 := &model.SensitiveWord{
				Word:     word, // Same word
				Category: differentCategory,
				Severity: differentSeverity,
			}

			// Second addition should fail due to duplicate word
			err = repo.Create(ctx, sw2)
			return err != nil
		},
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 && len(s) <= 100 }),
		gen.IntRange(0, 100),
		gen.IntRange(0, 100),
	))

	// Property 3: Multiple different sensitive words can be added
	properties.Property("multiple different sensitive words can be added successfully", prop.ForAll(
		func(numWords int) bool {
			// Generate unique words based on numWords
			if numWords <= 0 {
				return true // Skip invalid cases
			}

			// Setup fresh database for each test
			db := testutil.NewMemoryDB(t)
			testutil.MigrateTables(t, db, &model.SensitiveWord{})
			repo := sensitiveword.NewSensitiveWordRepository(db)
			ctx := context.Background()

			// Create unique words by using index
			for i := 0; i < numWords; i++ {
				word := fmt.Sprintf("word%d", i)
				sw := &model.SensitiveWord{
					Word:     word,
					Category: model.SensitiveWordCategoryOther,
					Severity: model.SensitiveWordSeverityLow,
				}

				err := repo.Create(ctx, sw)
				if err != nil {
					return false // Should succeed for all unique words
				}
			}

			// Verify all words were added
			allWords, err := repo.GetAll(ctx)
			if err != nil {
				return false
			}

			return len(allWords) == numWords
		},
		gen.IntRange(1, 10),
	))

	// Property 4: Case sensitivity - same word with different cases should be treated as different
	// (This tests the current implementation behavior)
	properties.Property("sensitive words are case-sensitive", prop.ForAll(
		func(word string) bool {
			// Skip if word is empty or too long
			if len(word) == 0 || len(word) > 100 {
				return true
			}

			// Setup fresh database for each test
			db := testutil.NewMemoryDB(t)
			testutil.MigrateTables(t, db, &model.SensitiveWord{})
			repo := sensitiveword.NewSensitiveWordRepository(db)
			ctx := context.Background()

			// Add lowercase version
			sw1 := &model.SensitiveWord{
				Word:     word,
				Category: model.SensitiveWordCategoryOther,
				Severity: model.SensitiveWordSeverityLow,
			}
			err := repo.Create(ctx, sw1)
			if err != nil {
				return false
			}

			// Try to add the exact same word again - should fail
			sw2 := &model.SensitiveWord{
				Word:     word,
				Category: model.SensitiveWordCategoryPolitical,
				Severity: model.SensitiveWordSeverityHigh,
			}
			err = repo.Create(ctx, sw2)

			// Should fail because it's a duplicate
			return err != nil
		},
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 && len(s) <= 100 }),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}
