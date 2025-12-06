package sensitiveword

import (
	"context"
	"strings"
	"testing"

	"gamelink/internal/model"
	"gamelink/internal/repository/sensitiveword"
	"gamelink/pkg/testutil"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: review-management-module, Property 7: 敏感词检测准确性**
// **Validates: Requirements 2.4, 5.5**
//
// For any content containing sensitive words, the system must detect the sensitive words
// and mark the review as pending status
func TestSensitiveWordDetectionAccuracy(t *testing.T) {
	properties := gopter.NewProperties(nil)

	// Property 1: Any content containing a sensitive word must be detected
	properties.Property("content with sensitive words must be detected", prop.ForAll(
		func(sensitiveWord string, prefix string, suffix string) bool {
			// Setup test database
			db := testutil.NewMemoryDB(t)
			testutil.MigrateTables(t, db, &model.SensitiveWord{})

			// Create repository and service
			repo := sensitiveword.NewSensitiveWordRepository(db)
			service := NewSensitiveWordService(repo)

			ctx := context.Background()

			// Add the sensitive word to the database
			word := &model.SensitiveWord{
				Word:     sensitiveWord,
				Category: model.SensitiveWordCategoryOther,
				Severity: model.SensitiveWordSeverityMedium,
			}
			if err := repo.Create(ctx, word); err != nil {
				// If word already exists (duplicate), skip this test case
				if strings.Contains(err.Error(), "已存在") {
					return true
				}
				t.Logf("Failed to create sensitive word: %v", err)
				return false
			}

			// Create content containing the sensitive word
			content := prefix + sensitiveWord + suffix

			// Detect sensitive words
			result, err := service.DetectSensitiveWords(ctx, DetectSensitiveWordsRequest{
				Content: content,
			})
			if err != nil {
				t.Logf("Detection failed: %v", err)
				return false
			}

			// Verify that sensitive words were detected
			if !result.HasSensitiveWords {
				t.Logf("Failed to detect sensitive word '%s' in content '%s'", sensitiveWord, content)
				return false
			}

			// Verify that the detected word matches
			found := false
			for _, detected := range result.DetectedWords {
				if strings.EqualFold(detected.Word, sensitiveWord) {
					found = true
					break
				}
			}

			if !found {
				t.Logf("Sensitive word '%s' not found in detected words", sensitiveWord)
				return false
			}

			return true
		},
		genNonEmptySensitiveWord(),
		gen.AlphaString(),
		gen.AlphaString(),
	))

	// Property 2: Content without sensitive words should not be flagged
	properties.Property("content without sensitive words should not be flagged", prop.ForAll(
		func(content string) bool {
			// Setup test database
			db := testutil.NewMemoryDB(t)
			testutil.MigrateTables(t, db, &model.SensitiveWord{})

			// Create repository and service
			repo := sensitiveword.NewSensitiveWordRepository(db)
			service := NewSensitiveWordService(repo)

			ctx := context.Background()

			// Add some sensitive words that are NOT in the content
			sensitiveWords := []string{"badword1", "badword2", "badword3"}
			for _, sw := range sensitiveWords {
				word := &model.SensitiveWord{
					Word:     sw,
					Category: model.SensitiveWordCategoryOther,
					Severity: model.SensitiveWordSeverityLow,
				}
				_ = repo.Create(ctx, word) // Ignore errors for duplicates
			}

			// Ensure content doesn't contain any of the sensitive words
			contentLower := strings.ToLower(content)
			containsSensitiveWord := false
			for _, sw := range sensitiveWords {
				if strings.Contains(contentLower, strings.ToLower(sw)) {
					containsSensitiveWord = true
					break
				}
			}

			// If content contains sensitive words, skip this test case
			if containsSensitiveWord {
				return true
			}

			// Detect sensitive words
			result, err := service.DetectSensitiveWords(ctx, DetectSensitiveWordsRequest{
				Content: content,
			})
			if err != nil {
				t.Logf("Detection failed: %v", err)
				return false
			}

			// Verify that no sensitive words were detected
			if result.HasSensitiveWords {
				t.Logf("False positive: detected sensitive words in clean content '%s': %v", content, result.DetectedWords)
				return false
			}

			return true
		},
		gen.AlphaString().SuchThat(func(s string) bool {
			// Generate strings that don't contain our test sensitive words
			lower := strings.ToLower(s)
			return !strings.Contains(lower, "badword1") &&
				!strings.Contains(lower, "badword2") &&
				!strings.Contains(lower, "badword3")
		}),
	))

	// Property 3: Detection should be case-insensitive
	properties.Property("detection should be case-insensitive", prop.ForAll(
		func(sensitiveWord string, caseVariation int) bool {
			// Setup test database
			db := testutil.NewMemoryDB(t)
			testutil.MigrateTables(t, db, &model.SensitiveWord{})

			// Create repository and service
			repo := sensitiveword.NewSensitiveWordRepository(db)
			service := NewSensitiveWordService(repo)

			ctx := context.Background()

			// Skip empty sensitive words
			if strings.TrimSpace(sensitiveWord) == "" {
				return true
			}

			// Add the sensitive word in lowercase
			word := &model.SensitiveWord{
				Word:     strings.ToLower(sensitiveWord),
				Category: model.SensitiveWordCategoryOther,
				Severity: model.SensitiveWordSeverityHigh,
			}
			if err := repo.Create(ctx, word); err != nil {
				if strings.Contains(err.Error(), "已存在") {
					return true
				}
				t.Logf("Failed to create sensitive word: %v", err)
				return false
			}

			// Create content with different case variations
			// Use absolute value to handle negative numbers properly
			variation := caseVariation
			if variation < 0 {
				variation = -variation
			}

			var content string
			switch variation % 3 {
			case 0:
				content = strings.ToUpper(sensitiveWord)
			case 1:
				content = strings.ToLower(sensitiveWord)
			case 2:
				// Mixed case
				if len(sensitiveWord) > 1 {
					content = strings.ToUpper(string(sensitiveWord[0])) + strings.ToLower(sensitiveWord[1:])
				} else if len(sensitiveWord) == 1 {
					// For single character, just use uppercase
					content = strings.ToUpper(sensitiveWord)
				} else {
					content = sensitiveWord
				}
			default:
				// Fallback (shouldn't happen)
				content = sensitiveWord
			}

			// Skip if content is empty (shouldn't happen, but safety check)
			if content == "" {
				t.Logf("Generated empty content for sensitive word '%s', skipping", sensitiveWord)
				return true
			}

			// Detect sensitive words
			result, err := service.DetectSensitiveWords(ctx, DetectSensitiveWordsRequest{
				Content: content,
			})
			if err != nil {
				t.Logf("Detection failed: %v", err)
				return false
			}

			// Verify that sensitive words were detected regardless of case
			if !result.HasSensitiveWords {
				t.Logf("Failed to detect case variation '%s' of sensitive word '%s'", content, sensitiveWord)
				return false
			}

			return true
		},
		genNonEmptySensitiveWord(),
		gen.Int(),
	))

	// Property 4: Multiple occurrences of the same sensitive word should all be detected
	// Note: The implementation groups occurrences by word, storing positions in a single DetectedWord
	properties.Property("multiple occurrences should all be detected", prop.ForAll(
		func(sensitiveWord string, count uint8) bool {
			if count == 0 {
				count = 1
			}
			if count > 10 {
				count = 10 // Limit to reasonable number
			}

			// Setup test database
			db := testutil.NewMemoryDB(t)
			testutil.MigrateTables(t, db, &model.SensitiveWord{})

			// Create repository and service
			repo := sensitiveword.NewSensitiveWordRepository(db)
			service := NewSensitiveWordService(repo)

			ctx := context.Background()

			// Add the sensitive word
			word := &model.SensitiveWord{
				Word:     sensitiveWord,
				Category: model.SensitiveWordCategoryViolent,
				Severity: model.SensitiveWordSeverityHigh,
			}
			if err := repo.Create(ctx, word); err != nil {
				if strings.Contains(err.Error(), "已存在") {
					return true
				}
				t.Logf("Failed to create sensitive word: %v", err)
				return false
			}

			// Create content with multiple occurrences
			parts := make([]string, count)
			for i := range parts {
				parts[i] = sensitiveWord
			}
			content := strings.Join(parts, " ")

			// Detect sensitive words
			result, err := service.DetectSensitiveWords(ctx, DetectSensitiveWordsRequest{
				Content: content,
			})
			if err != nil {
				t.Logf("Detection failed: %v", err)
				return false
			}

			// Verify that all occurrences were detected
			if !result.HasSensitiveWords {
				t.Logf("Failed to detect sensitive word in content with %d occurrences", count)
				return false
			}

			// Count detected positions (implementation groups by word, positions in Positions slice)
			positionCount := 0
			for _, detected := range result.DetectedWords {
				if strings.EqualFold(detected.Word, sensitiveWord) {
					positionCount += len(detected.Positions)
				}
			}

			if positionCount != int(count) {
				t.Logf("Expected %d positions, but detected %d", count, positionCount)
				return false
			}

			return true
		},
		genNonEmptySensitiveWord(),
		gen.UInt8(),
	))

	// Property 5: Empty content should not trigger detection
	properties.Property("empty content should not trigger detection", prop.ForAll(
		func() bool {
			// Setup test database
			db := testutil.NewMemoryDB(t)
			testutil.MigrateTables(t, db, &model.SensitiveWord{})

			// Create repository and service
			repo := sensitiveword.NewSensitiveWordRepository(db)
			service := NewSensitiveWordService(repo)

			ctx := context.Background()

			// Add some sensitive words
			word := &model.SensitiveWord{
				Word:     "testword",
				Category: model.SensitiveWordCategoryOther,
				Severity: model.SensitiveWordSeverityLow,
			}
			_ = repo.Create(ctx, word)

			// Test with empty content
			result, err := service.DetectSensitiveWords(ctx, DetectSensitiveWordsRequest{
				Content: "",
			})
			if err != nil {
				t.Logf("Detection failed: %v", err)
				return false
			}

			// Verify no detection for empty content
			if result.HasSensitiveWords {
				t.Logf("False positive: detected sensitive words in empty content")
				return false
			}

			return true
		},
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// genNonEmptySensitiveWord generates non-empty strings suitable for sensitive words
func genNonEmptySensitiveWord() gopter.Gen {
	return gen.AlphaString().SuchThat(func(s string) bool {
		trimmed := strings.TrimSpace(s)
		return len(trimmed) > 0 && len(trimmed) <= 50
	}).Map(func(s string) string {
		return strings.TrimSpace(s)
	})
}
