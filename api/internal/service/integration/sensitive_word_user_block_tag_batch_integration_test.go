// Package integration provides batch operation integration tests for SensitiveWord, UserBlock, and UserTag modules.
package integration

import (
	"context"
	"testing"

	"gamelink/internal/model"
	userblockrepo "gamelink/internal/repository/userblock"
	userrepository "gamelink/internal/repository/user"
	sensitivewordrepo "gamelink/internal/repository/sensitiveword"
	"gamelink/internal/service/sensitiveword"
	userblockservice "gamelink/internal/service/userblock"
	userservice "gamelink/internal/service/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// ============================================================================
// Test Helper Functions
// ============================================================================

// CreateTestSensitiveWordWithStatus creates a test sensitive word with specified status.
func CreateTestSensitiveWordWithStatus(t *testing.T, db *gorm.DB, word string, category model.SensitiveWordCategory, severity model.SensitiveWordSeverity, isActive bool) *model.SensitiveWord {
	t.Helper()
	sw := &model.SensitiveWord{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Word:        word,
		Category:    category,
		MatchType:   model.SensitiveWordMatchTypeExact,
		Severity:    severity,
		Replacement: "***",
		IsActive:    isActive,
	}
	if err := db.Create(sw).Error; err != nil {
		t.Fatalf("Failed to create test sensitive word: %v", err)
	}
	return sw
}

// CreateTestUserBlock creates a test user block record.
func CreateTestUserBlock(t *testing.T, db *gorm.DB, blockerID, blockedID uint64, blockerType, blockedType model.BlockUserType, reason string) *model.UserBlock {
	t.Helper()
	block := &model.UserBlock{
		Base: model.Base{
			ExtJSON: "{}",
		},
		BlockerID:   blockerID,
		BlockerType: blockerType,
		BlockedID:   blockedID,
		BlockedType: blockedType,
		Reason:      reason,
		Status:      model.BlockStatusActive,
		BlockedAt:   db.NowFunc(),
	}
	if err := db.Create(block).Error; err != nil {
		t.Fatalf("Failed to create test user block: %v", err)
	}
	return block
}

// CreateTestUserTag creates a test user tag.
func CreateTestUserTag(t *testing.T, db *gorm.DB, name, color, description string) *model.UserTag {
	t.Helper()
	tag := &model.UserTag{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Name:        name,
		Color:       color,
		Description: description,
		CreatedBy:   1,
	}
	if err := db.Create(tag).Error; err != nil {
		t.Fatalf("Failed to create test user tag: %v", err)
	}
	return tag
}

// AssignTagToUser assigns a tag to a user via the join table.
func AssignTagToUser(t *testing.T, db *gorm.DB, userID, tagID uint64) {
	t.Helper()
	relation := &model.UserTagRelation{
		Base: model.Base{
			ExtJSON: "{}",
		},
		UserID: userID,
		TagID:  tagID,
	}
	if err := db.Create(relation).Error; err != nil {
		t.Fatalf("Failed to assign tag to user: %v", err)
	}
}

// GetSensitiveWordCount returns the count of sensitive words in database.
func GetSensitiveWordCount(t *testing.T, db *gorm.DB) int64 {
	t.Helper()
	var count int64
	if err := db.Model(&model.SensitiveWord{}).Count(&count).Error; err != nil {
		t.Fatalf("Failed to count sensitive words: %v", err)
	}
	return count
}

// SensitiveWordExists checks if a sensitive word exists by ID.
func SensitiveWordExists(t *testing.T, db *gorm.DB, id uint64) bool {
	t.Helper()
	var count int64
	if err := db.Model(&model.SensitiveWord{}).Where("id = ?", id).Count(&count).Error; err != nil {
		t.Fatalf("Failed to check sensitive word existence: %v", err)
	}
	return count > 0
}

// GetSensitiveWordByID retrieves a sensitive word by ID.
func GetSensitiveWordByID(t *testing.T, db *gorm.DB, id uint64) *model.SensitiveWord {
	t.Helper()
	var sw model.SensitiveWord
	if err := db.First(&sw, id).Error; err != nil {
		return nil
	}
	return &sw
}

// GetUserBlockCount returns the count of user blocks in database.
func GetUserBlockCount(t *testing.T, db *gorm.DB) int64 {
	t.Helper()
	var count int64
	if err := db.Model(&model.UserBlock{}).Count(&count).Error; err != nil {
		t.Fatalf("Failed to count user blocks: %v", err)
	}
	return count
}

// UserBlockExists checks if a user block exists by ID.
func UserBlockExists(t *testing.T, db *gorm.DB, id uint64) bool {
	t.Helper()
	var count int64
	if err := db.Model(&model.UserBlock{}).Where("id = ?", id).Count(&count).Error; err != nil {
		t.Fatalf("Failed to check user block existence: %v", err)
	}
	return count > 0
}

// GetUserBlockByID retrieves a user block by ID.
func GetUserBlockByID(t *testing.T, db *gorm.DB, id uint64) *model.UserBlock {
	t.Helper()
	var block model.UserBlock
	if err := db.First(&block, id).Error; err != nil {
		return nil
	}
	return &block
}

// GetUserTagCount returns the count of user tags in database.
func GetUserTagCount(t *testing.T, db *gorm.DB) int64 {
	t.Helper()
	var count int64
	if err := db.Model(&model.UserTag{}).Count(&count).Error; err != nil {
		t.Fatalf("Failed to count user tags: %v", err)
	}
	return count
}

// UserTagExists checks if a user tag exists by ID.
func UserTagExists(t *testing.T, db *gorm.DB, id uint64) bool {
	t.Helper()
	var count int64
	if err := db.Model(&model.UserTag{}).Where("id = ?", id).Count(&count).Error; err != nil {
		t.Fatalf("Failed to check user tag existence: %v", err)
	}
	return count > 0
}

// GetUserTagByID retrieves a user tag by ID.
func GetUserTagByID(t *testing.T, db *gorm.DB, id uint64) *model.UserTag {
	t.Helper()
	var tag model.UserTag
	if err := db.First(&tag, id).Error; err != nil {
		return nil
	}
	return &tag
}

// GetUserTagRelationCount returns the count of user tag relations.
func GetUserTagRelationCount(t *testing.T, db *gorm.DB) int64 {
	t.Helper()
	var count int64
	if err := db.Model(&model.UserTagRelation{}).Count(&count).Error; err != nil {
		t.Fatalf("Failed to count user tag relations: %v", err)
	}
	return count
}

// UserHasTag checks if a user has a specific tag.
func UserHasTag(t *testing.T, db *gorm.DB, userID, tagID uint64) bool {
	t.Helper()
	var count int64
	if err := db.Model(&model.UserTagRelation{}).
		Where("user_id = ? AND tag_id = ?", userID, tagID).
		Count(&count).Error; err != nil {
		t.Fatalf("Failed to check user tag relation: %v", err)
	}
	return count > 0
}

// ============================================================================
// SensitiveWord Batch Operations Tests
// ============================================================================

func TestSensitiveWordService_BatchAddSensitiveWords_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	swRepo := sensitivewordrepo.NewSensitiveWordRepository(db)
	swService := sensitiveword.NewSensitiveWordService(swRepo)

	// Execute batch add
	req := sensitiveword.BatchAddSensitiveWordsRequest{
		Words:     []string{"test1", "test2", "test3"},
		Category:  model.SensitiveWordCategoryAbuse,
		Severity:  model.SensitiveWordSeverityHigh,
		MatchType: model.SensitiveWordMatchTypeExact,
	}
	result, err := swService.BatchAddSensitiveWords(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Len(t, result.SuccessItems, 3)
	assert.Empty(t, result.FailedItems)

	// Verify database state
	finalCount := GetSensitiveWordCount(t, db)
	assert.Equal(t, int64(3), finalCount)
}

func TestSensitiveWordService_BatchAddSensitiveWords_WithDuplicates(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	swRepo := sensitivewordrepo.NewSensitiveWordRepository(db)
	swService := sensitiveword.NewSensitiveWordService(swRepo)

	// Create an existing sensitive word
	existing := CreateTestSensitiveWordWithStatus(t, db, "existing", model.SensitiveWordCategoryAbuse, model.SensitiveWordSeverityHigh, true)

	// Execute batch add with duplicates
	req := sensitiveword.BatchAddSensitiveWordsRequest{
		Words:     []string{"new1", "existing", "new2"},
		Category:  model.SensitiveWordCategoryAbuse,
		Severity:  model.SensitiveWordSeverityHigh,
		MatchType: model.SensitiveWordMatchTypeExact,
	}
	result, err := swService.BatchAddSensitiveWords(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 2, result.SuccessCount) // new1, new2
	assert.Equal(t, 1, result.FailedCount)  // existing (duplicate)
	assert.Len(t, result.SuccessItems, 2)
	assert.Len(t, result.FailedItems, 1)

	// Verify failed item
	failedItem := result.FailedItems[0]
	assert.Equal(t, uint64(0), failedItem.ID)
	assert.Contains(t, failedItem.Message, "already exists")

	// Verify database state
	finalCount := GetSensitiveWordCount(t, db)
	assert.Equal(t, int64(3), finalCount) // existing + new1 + new2
	assert.True(t, SensitiveWordExists(t, db, existing.ID))
}

func TestSensitiveWordService_BatchAddSensitiveWords_WithEmptyWords(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	swRepo := sensitivewordrepo.NewSensitiveWordRepository(db)
	swService := sensitiveword.NewSensitiveWordService(swRepo)

	// Execute batch add with empty words
	req := sensitiveword.BatchAddSensitiveWordsRequest{
		Words:     []string{"valid1", "", "  ", "valid2"},
		Category:  model.SensitiveWordCategoryAbuse,
		Severity:  model.SensitiveWordSeverityHigh,
		MatchType: model.SensitiveWordMatchTypeExact,
	}
	result, err := swService.BatchAddSensitiveWords(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 4, result.TotalCount)
	assert.Equal(t, 2, result.SuccessCount) // valid1, valid2
	assert.Equal(t, 2, result.FailedCount)  // empty strings
	assert.Len(t, result.SuccessItems, 2)
	assert.Len(t, result.FailedItems, 2)

	// Verify failed items
	for _, item := range result.FailedItems {
		assert.Contains(t, item.Message, "cannot be empty")
	}

	// Verify database state
	finalCount := GetSensitiveWordCount(t, db)
	assert.Equal(t, int64(2), finalCount)
}

func TestSensitiveWordService_BatchAddSensitiveWords_ExceedsLimit(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	swRepo := sensitivewordrepo.NewSensitiveWordRepository(db)
	swService := sensitiveword.NewSensitiveWordService(swRepo)

	// Create 101 words (exceeds limit)
	words := make([]string, 101)
	for i := 0; i < 101; i++ {
		words[i] = "test" + string(rune('0'+i))
	}

	req := sensitiveword.BatchAddSensitiveWordsRequest{
		Words:     words,
		Category:  model.SensitiveWordCategoryAbuse,
		Severity:  model.SensitiveWordSeverityHigh,
		MatchType: model.SensitiveWordMatchTypeExact,
	}
	result, err := swService.BatchAddSensitiveWords(ctx, req)

	// Should return validation error
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestSensitiveWordService_BatchDeleteSensitiveWords_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	swRepo := sensitivewordrepo.NewSensitiveWordRepository(db)
	swService := sensitiveword.NewSensitiveWordService(swRepo)

	// Create test sensitive words
	sw1 := CreateTestSensitiveWord(t, db, "word1", model.SensitiveWordCategoryAbuse)
	sw2 := CreateTestSensitiveWord(t, db, "word2", model.SensitiveWordCategoryAbuse)
	sw3 := CreateTestSensitiveWord(t, db, "word3", model.SensitiveWordCategoryAbuse)

	// Execute batch delete
	result, err := swService.BatchDeleteSensitiveWords(ctx, []uint64{sw1.ID, sw2.ID, sw3.ID})

	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Len(t, result.SuccessItems, 3)
	assert.Contains(t, result.SuccessItems, sw1.ID)
	assert.Contains(t, result.SuccessItems, sw2.ID)
	assert.Contains(t, result.SuccessItems, sw3.ID)
	assert.Empty(t, result.FailedItems)

	// Verify database state
	finalCount := GetSensitiveWordCount(t, db)
	assert.Equal(t, int64(0), finalCount)
	assert.False(t, SensitiveWordExists(t, db, sw1.ID))
	assert.False(t, SensitiveWordExists(t, db, sw2.ID))
	assert.False(t, SensitiveWordExists(t, db, sw3.ID))
}

func TestSensitiveWordService_BatchDeleteSensitiveWords_PartialNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	swRepo := sensitivewordrepo.NewSensitiveWordRepository(db)
	swService := sensitiveword.NewSensitiveWordService(swRepo)

	// Create one sensitive word
	sw1 := CreateTestSensitiveWord(t, db, "word1", model.SensitiveWordCategoryAbuse)
	nonExistentID := uint64(999999)

	// Execute batch delete with mixed valid and invalid IDs
	result, err := swService.BatchDeleteSensitiveWords(ctx, []uint64{sw1.ID, nonExistentID, nonExistentID + 1})

	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 1, result.SuccessCount) // sw1
	assert.Equal(t, 2, result.FailedCount)  // non-existent
	assert.Contains(t, result.SuccessItems, sw1.ID)
	assert.Len(t, result.FailedItems, 2)

	// Verify failed items
	for _, item := range result.FailedItems {
		assert.Contains(t, item.Message, "not found")
	}

	// Verify only valid word was deleted
	assert.False(t, SensitiveWordExists(t, db, sw1.ID))
}

func TestSensitiveWordService_BatchUpdateSensitiveWordStatus_EnableSuccess(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	swRepo := sensitivewordrepo.NewSensitiveWordRepository(db)
	swService := sensitiveword.NewSensitiveWordService(swRepo)

	// Create sensitive words with various active states
	sw1 := CreateTestSensitiveWordWithStatus(t, db, "word1", model.SensitiveWordCategoryAbuse, model.SensitiveWordSeverityHigh, false)
	sw2 := CreateTestSensitiveWordWithStatus(t, db, "word2", model.SensitiveWordCategoryAbuse, model.SensitiveWordSeverityHigh, false)
	sw3 := CreateTestSensitiveWordWithStatus(t, db, "word3", model.SensitiveWordCategoryAbuse, model.SensitiveWordSeverityHigh, true)

	// Execute batch enable
	req := sensitiveword.BatchUpdateSensitiveWordStatusRequest{
		IDs:      []uint64{sw1.ID, sw2.ID, sw3.ID},
		IsActive: true,
	}
	result, err := swService.BatchUpdateSensitiveWordStatus(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Len(t, result.SuccessItems, 3)

	// Verify database state - all should be active
	updatedSW1 := GetSensitiveWordByID(t, db, sw1.ID)
	updatedSW2 := GetSensitiveWordByID(t, db, sw2.ID)
	updatedSW3 := GetSensitiveWordByID(t, db, sw3.ID)

	assert.NotNil(t, updatedSW1)
	assert.NotNil(t, updatedSW2)
	assert.NotNil(t, updatedSW3)

	assert.True(t, updatedSW1.IsActive, "word1 should be active")
	assert.True(t, updatedSW2.IsActive, "word2 should be active")
	assert.True(t, updatedSW3.IsActive, "word3 should be active")
}

func TestSensitiveWordService_BatchUpdateSensitiveWordStatus_DisableSuccess(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	swRepo := sensitivewordrepo.NewSensitiveWordRepository(db)
	swService := sensitiveword.NewSensitiveWordService(swRepo)

	// Create sensitive words
	sw1 := CreateTestSensitiveWordWithStatus(t, db, "word1", model.SensitiveWordCategoryAbuse, model.SensitiveWordSeverityHigh, true)
	sw2 := CreateTestSensitiveWordWithStatus(t, db, "word2", model.SensitiveWordCategoryAbuse, model.SensitiveWordSeverityHigh, true)
	sw3 := CreateTestSensitiveWordWithStatus(t, db, "word3", model.SensitiveWordCategoryAbuse, model.SensitiveWordSeverityHigh, false)

	// Execute batch disable
	req := sensitiveword.BatchUpdateSensitiveWordStatusRequest{
		IDs:      []uint64{sw1.ID, sw2.ID, sw3.ID},
		IsActive: false,
	}
	result, err := swService.BatchUpdateSensitiveWordStatus(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)

	// Verify database state - all should be inactive
	updatedSW1 := GetSensitiveWordByID(t, db, sw1.ID)
	updatedSW2 := GetSensitiveWordByID(t, db, sw2.ID)
	updatedSW3 := GetSensitiveWordByID(t, db, sw3.ID)

	assert.False(t, updatedSW1.IsActive, "word1 should be inactive")
	assert.False(t, updatedSW2.IsActive, "word2 should be inactive")
	assert.False(t, updatedSW3.IsActive, "word3 should be inactive")
}

func TestSensitiveWordService_BatchUpdateSensitiveWordStatus_PartialNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	swRepo := sensitivewordrepo.NewSensitiveWordRepository(db)
	swService := sensitiveword.NewSensitiveWordService(swRepo)

	// Create one sensitive word
	sw1 := CreateTestSensitiveWord(t, db, "word1", model.SensitiveWordCategoryAbuse)
	nonExistentID := uint64(999999)

	// Execute batch update with mixed valid and invalid IDs
	req := sensitiveword.BatchUpdateSensitiveWordStatusRequest{
		IDs:      []uint64{sw1.ID, nonExistentID},
		IsActive: true,
	}
	result, err := swService.BatchUpdateSensitiveWordStatus(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 2, result.TotalCount)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Contains(t, result.SuccessItems, sw1.ID)

	// Verify failed item
	assert.Len(t, result.FailedItems, 1)
	assert.Equal(t, nonExistentID, result.FailedItems[0].ID)
	assert.Contains(t, result.FailedItems[0].Message, "not found")

	// Verify valid word was updated
	updatedSW := GetSensitiveWordByID(t, db, sw1.ID)
	assert.True(t, updatedSW.IsActive)
}

// ============================================================================
// UserBlock Batch Operations Tests
// ============================================================================

func TestUserBlockService_BatchUnblock_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test users
	user1 := CreateTestUser(t, db, "user1")
	user2 := CreateTestUser(t, db, "user2")
	user3 := CreateTestUser(t, db, "user3")
	user4 := CreateTestUser(t, db, "user4")
	admin := CreateTestUser(t, db, "admin")

	// Create repository and service
	blockRepo := userblockrepo.NewUserBlockRepository(db)
	userRepo := userrepository.NewUserRepository(db)
	blockService := userblockservice.NewUserBlockService(blockRepo, userRepo)

	// Create test user blocks
	block1 := CreateTestUserBlock(t, db, user1.ID, user4.ID, model.BlockUserTypeUser, model.BlockUserTypeUser, "reason1")
	block2 := CreateTestUserBlock(t, db, user2.ID, user4.ID, model.BlockUserTypeUser, model.BlockUserTypeUser, "reason2")
	block3 := CreateTestUserBlock(t, db, user3.ID, user4.ID, model.BlockUserTypeUser, model.BlockUserTypeUser, "reason3")

	// Execute batch unblock
	result, err := blockService.BatchUnblock(ctx, []uint64{block1.ID, block2.ID, block3.ID}, admin.ID, "admin unblock")

	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Empty(t, result.FailedIDs)

	// Verify database state - blocks should be canceled
	updatedBlock1 := GetUserBlockByID(t, db, block1.ID)
	updatedBlock2 := GetUserBlockByID(t, db, block2.ID)
	updatedBlock3 := GetUserBlockByID(t, db, block3.ID)

	assert.NotNil(t, updatedBlock1)
	assert.NotNil(t, updatedBlock2)
	assert.NotNil(t, updatedBlock3)

	assert.Equal(t, model.BlockStatusAdminCanceled, updatedBlock1.Status)
	assert.Equal(t, model.BlockStatusAdminCanceled, updatedBlock2.Status)
	assert.Equal(t, model.BlockStatusAdminCanceled, updatedBlock3.Status)

	assert.NotNil(t, updatedBlock1.CanceledBy)
	assert.NotNil(t, updatedBlock2.CanceledBy)
	assert.NotNil(t, updatedBlock3.CanceledBy)
	assert.Equal(t, admin.ID, *updatedBlock1.CanceledBy)
	assert.Equal(t, admin.ID, *updatedBlock2.CanceledBy)
	assert.Equal(t, admin.ID, *updatedBlock3.CanceledBy)
}

func TestUserBlockService_BatchUnblock_PartialNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test users
	user1 := CreateTestUser(t, db, "user1")
	user2 := CreateTestUser(t, db, "user2")
	admin := CreateTestUser(t, db, "admin")

	// Create repository and service
	blockRepo := userblockrepo.NewUserBlockRepository(db)
	userRepo := userrepository.NewUserRepository(db)
	blockService := userblockservice.NewUserBlockService(blockRepo, userRepo)

	// Create one user block
	block1 := CreateTestUserBlock(t, db, user1.ID, user2.ID, model.BlockUserTypeUser, model.BlockUserTypeUser, "reason1")
	nonExistentID := uint64(999999)

	// Execute batch unblock with mixed valid and invalid IDs
	result, err := blockService.BatchUnblock(ctx, []uint64{block1.ID, nonExistentID}, admin.ID, "admin unblock")

	require.NoError(t, err)
	assert.Equal(t, 2, result.TotalCount)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Contains(t, result.FailedIDs, nonExistentID)

	// Verify valid block was canceled
	updatedBlock := GetUserBlockByID(t, db, block1.ID)
	assert.Equal(t, model.BlockStatusAdminCanceled, updatedBlock.Status)
}

func TestUserBlockService_BatchUnblock_AlreadyCanceled(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test users
	user1 := CreateTestUser(t, db, "user1")
	user2 := CreateTestUser(t, db, "user2")
	user3 := CreateTestUser(t, db, "user3")
	admin := CreateTestUser(t, db, "admin")

	// Create repository and service
	blockRepo := userblockrepo.NewUserBlockRepository(db)
	userRepo := userrepository.NewUserRepository(db)
	blockService := userblockservice.NewUserBlockService(blockRepo, userRepo)

	// Create test user blocks
	block1 := CreateTestUserBlock(t, db, user1.ID, user2.ID, model.BlockUserTypeUser, model.BlockUserTypeUser, "reason1")

	// Manually cancel block2
	block2 := CreateTestUserBlock(t, db, user2.ID, user3.ID, model.BlockUserTypeUser, model.BlockUserTypeUser, "reason2")
	db.Model(block2).Update("status", model.BlockStatusCanceled)

	// Execute batch unblock
	result, err := blockService.BatchUnblock(ctx, []uint64{block1.ID, block2.ID}, admin.ID, "admin unblock")

	require.NoError(t, err)
	assert.Equal(t, 2, result.TotalCount)
	assert.Equal(t, 1, result.SuccessCount) // block1
	assert.Equal(t, 1, result.FailedCount)  // block2 (already canceled)
	assert.Contains(t, result.FailedIDs, block2.ID)

	// Verify block1 was canceled
	updatedBlock1 := GetUserBlockByID(t, db, block1.ID)
	assert.Equal(t, model.BlockStatusAdminCanceled, updatedBlock1.Status)
}

func TestUserBlockService_BatchDelete_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test users
	user1 := CreateTestUser(t, db, "user1")
	user2 := CreateTestUser(t, db, "user2")
	user3 := CreateTestUser(t, db, "user3")
	user4 := CreateTestUser(t, db, "user4")

	// Create repository and service
	blockRepo := userblockrepo.NewUserBlockRepository(db)
	userRepo := userrepository.NewUserRepository(db)
	blockService := userblockservice.NewUserBlockService(blockRepo, userRepo)

	// Create test user blocks
	block1 := CreateTestUserBlock(t, db, user1.ID, user4.ID, model.BlockUserTypeUser, model.BlockUserTypeUser, "reason1")
	block2 := CreateTestUserBlock(t, db, user2.ID, user4.ID, model.BlockUserTypeUser, model.BlockUserTypeUser, "reason2")
	block3 := CreateTestUserBlock(t, db, user3.ID, user4.ID, model.BlockUserTypeUser, model.BlockUserTypeUser, "reason3")

	// Execute batch delete
	result, err := blockService.BatchDelete(ctx, []uint64{block1.ID, block2.ID, block3.ID})

	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Empty(t, result.FailedIDs)

	// Verify database state
	finalCount := GetUserBlockCount(t, db)
	assert.Equal(t, int64(0), finalCount)
	assert.False(t, UserBlockExists(t, db, block1.ID))
	assert.False(t, UserBlockExists(t, db, block2.ID))
	assert.False(t, UserBlockExists(t, db, block3.ID))
}

func TestUserBlockService_BatchDelete_PartialNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test users
	user1 := CreateTestUser(t, db, "user1")
	user2 := CreateTestUser(t, db, "user2")

	// Create repository and service
	blockRepo := userblockrepo.NewUserBlockRepository(db)
	userRepo := userrepository.NewUserRepository(db)
	blockService := userblockservice.NewUserBlockService(blockRepo, userRepo)

	// Create one user block
	block1 := CreateTestUserBlock(t, db, user1.ID, user2.ID, model.BlockUserTypeUser, model.BlockUserTypeUser, "reason1")
	nonExistentID := uint64(999999)

	// Execute batch delete with mixed valid and invalid IDs
	result, err := blockService.BatchDelete(ctx, []uint64{block1.ID, nonExistentID, nonExistentID + 1})

	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 2, result.FailedCount)
	assert.Len(t, result.FailedIDs, 2)

	// Verify only valid block was deleted
	assert.False(t, UserBlockExists(t, db, block1.ID))
}

// ============================================================================
// UserTag Batch Operations Tests
// ============================================================================

func TestUserTagService_BatchDeleteTags_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	tagRepo := userrepository.NewUserTagRepository(db)
	userRepo := userrepository.NewUserRepository(db)
	tagService := userservice.NewUserTagService(tagRepo, userRepo, nil)

	// Create test tags
	tag1 := CreateTestUserTag(t, db, "VIP", "#FF6B6B", "VIP users")
	tag2 := CreateTestUserTag(t, db, "Regular", "#4ECDC4", "Regular users")
	tag3 := CreateTestUserTag(t, db, "Premium", "#45B7D1", "Premium users")

	// Verify initial count
	initialCount := GetUserTagCount(t, db)
	assert.Equal(t, int64(3), initialCount)

	// Execute batch delete
	result, err := tagService.BatchDeleteTags(ctx, []uint64{tag1.ID, tag2.ID, tag3.ID})

	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Len(t, result.SuccessItems, 3)
	assert.Contains(t, result.SuccessItems, tag1.ID)
	assert.Contains(t, result.SuccessItems, tag2.ID)
	assert.Contains(t, result.SuccessItems, tag3.ID)
	assert.Empty(t, result.FailedItems)

	// Verify database state
	finalCount := GetUserTagCount(t, db)
	assert.Equal(t, int64(0), finalCount)
	assert.False(t, UserTagExists(t, db, tag1.ID))
	assert.False(t, UserTagExists(t, db, tag2.ID))
	assert.False(t, UserTagExists(t, db, tag3.ID))
}

func TestUserTagService_BatchDeleteTags_PartialNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	tagRepo := userrepository.NewUserTagRepository(db)
	userRepo := userrepository.NewUserRepository(db)
	tagService := userservice.NewUserTagService(tagRepo, userRepo, nil)

	// Create one tag
	tag1 := CreateTestUserTag(t, db, "VIP", "#FF6B6B", "VIP users")
	nonExistentID := uint64(999999)

	// Execute batch delete with mixed valid and invalid IDs
	result, err := tagService.BatchDeleteTags(ctx, []uint64{tag1.ID, nonExistentID, nonExistentID + 1})

	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 2, result.FailedCount)
	assert.Contains(t, result.SuccessItems, tag1.ID)
	assert.Len(t, result.FailedItems, 2)

	// Verify failed items
	for _, item := range result.FailedItems {
		assert.Equal(t, nonExistentID, item.ID) // or nonExistentID + 1
		assert.Contains(t, item.Message, "不存在")
	}

	// Verify only valid tag was deleted
	assert.False(t, UserTagExists(t, db, tag1.ID))
}

func TestUserTagService_BatchAssignTagsToUsers_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test users
	user1 := CreateTestUser(t, db, "user1")
	user2 := CreateTestUser(t, db, "user2")

	// Create repository and service
	tagRepo := userrepository.NewUserTagRepository(db)
	userRepo := userrepository.NewUserRepository(db)
	tagService := userservice.NewUserTagService(tagRepo, userRepo, nil)

	// Create test tags
	tag1 := CreateTestUserTag(t, db, "VIP", "#FF6B6B", "VIP users")
	tag2 := CreateTestUserTag(t, db, "Premium", "#45B7D1", "Premium users")

	// Execute batch assign (2 users x 2 tags = 4 assignments)
	result, err := tagService.BatchAssignTagsToUsers(ctx, []uint64{user1.ID, user2.ID}, []uint64{tag1.ID, tag2.ID})

	require.NoError(t, err)
	assert.Equal(t, 4, result.TotalCount)
	assert.Equal(t, 4, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Len(t, result.SuccessItems, 2) // 2 unique users
	assert.Contains(t, result.SuccessItems, user1.ID)
	assert.Contains(t, result.SuccessItems, user2.ID)
	assert.Empty(t, result.FailedItems)

	// Verify database state - each user should have both tags
	assert.True(t, UserHasTag(t, db, user1.ID, tag1.ID))
	assert.True(t, UserHasTag(t, db, user1.ID, tag2.ID))
	assert.True(t, UserHasTag(t, db, user2.ID, tag1.ID))
	assert.True(t, UserHasTag(t, db, user2.ID, tag2.ID))

	// Verify total relations count
	relationCount := GetUserTagRelationCount(t, db)
	assert.Equal(t, int64(4), relationCount)
}

func TestUserTagService_BatchAssignTagsToUsers_UserNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test user
	user1 := CreateTestUser(t, db, "user1")
	nonExistentUserID := uint64(999999)

	// Create repository and service
	tagRepo := userrepository.NewUserTagRepository(db)
	userRepo := userrepository.NewUserRepository(db)
	tagService := userservice.NewUserTagService(tagRepo, userRepo, nil)

	// Create test tag
	tag1 := CreateTestUserTag(t, db, "VIP", "#FF6B6B", "VIP users")

	// Execute batch assign with non-existent user
	result, err := tagService.BatchAssignTagsToUsers(ctx, []uint64{user1.ID, nonExistentUserID}, []uint64{tag1.ID})

	require.NoError(t, err)
	assert.Equal(t, 2, result.TotalCount) // 2 users x 1 tag = 2 assignments
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Contains(t, result.SuccessItems, user1.ID)
	assert.Len(t, result.FailedItems, 1)

	// Verify failed item
	failedItem := result.FailedItems[0]
	assert.Contains(t, failedItem.Message, "用户不存在")

	// Verify only valid user got the tag
	assert.True(t, UserHasTag(t, db, user1.ID, tag1.ID))
}

func TestUserTagService_BatchAssignTagsToUsers_TagNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test user
	user1 := CreateTestUser(t, db, "user1")
	nonExistentTagID := uint64(999999)

	// Create repository and service
	tagRepo := userrepository.NewUserTagRepository(db)
	userRepo := userrepository.NewUserRepository(db)
	tagService := userservice.NewUserTagService(tagRepo, userRepo, nil)

	// Create test tag
	tag1 := CreateTestUserTag(t, db, "VIP", "#FF6B6B", "VIP users")

	// Execute batch assign with non-existent tag
	result, err := tagService.BatchAssignTagsToUsers(ctx, []uint64{user1.ID}, []uint64{tag1.ID, nonExistentTagID})

	require.NoError(t, err)
	assert.Equal(t, 2, result.TotalCount) // 1 user x 2 tags = 2 assignments
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Len(t, result.FailedItems, 1)

	// Verify failed item contains tag info
	failedItem := result.FailedItems[0]
	assert.Contains(t, failedItem.Message, "标签不存在")
	assert.Contains(t, failedItem.Message, "999999") // tag ID in message
}

func TestUserTagService_BatchAssignTagsToUsers_DuplicateAssignment(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test user
	user1 := CreateTestUser(t, db, "user1")

	// Create repository and service
	tagRepo := userrepository.NewUserTagRepository(db)
	userRepo := userrepository.NewUserRepository(db)
	tagService := userservice.NewUserTagService(tagRepo, userRepo, nil)

	// Create test tag
	tag1 := CreateTestUserTag(t, db, "VIP", "#FF6B6B", "VIP users")

	// Pre-assign tag to user
	AssignTagToUser(t, db, user1.ID, tag1.ID)

	// Execute batch assign - should handle duplicate gracefully
	result, err := tagService.BatchAssignTagsToUsers(ctx, []uint64{user1.ID}, []uint64{tag1.ID})

	// The service's AddTagToUser should fail for duplicate
	require.NoError(t, err)
	assert.Equal(t, 1, result.TotalCount)
	assert.Equal(t, 0, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Len(t, result.FailedItems, 1)

	// Verify failed item
	failedItem := result.FailedItems[0]
	assert.Contains(t, failedItem.Message, "已拥有此标签")
}

func TestUserTagService_BatchRemoveTagsFromUsers_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test users
	user1 := CreateTestUser(t, db, "user1")
	user2 := CreateTestUser(t, db, "user2")

	// Create repository and service
	tagRepo := userrepository.NewUserTagRepository(db)
	userRepo := userrepository.NewUserRepository(db)
	tagService := userservice.NewUserTagService(tagRepo, userRepo, nil)

	// Create test tags
	tag1 := CreateTestUserTag(t, db, "VIP", "#FF6B6B", "VIP users")
	tag2 := CreateTestUserTag(t, db, "Premium", "#45B7D1", "Premium users")

	// Assign tags to users
	AssignTagToUser(t, db, user1.ID, tag1.ID)
	AssignTagToUser(t, db, user1.ID, tag2.ID)
	AssignTagToUser(t, db, user2.ID, tag1.ID)
	AssignTagToUser(t, db, user2.ID, tag2.ID)

	// Verify initial state
	assert.Equal(t, int64(4), GetUserTagRelationCount(t, db))

	// Execute batch remove (2 users x 2 tags = 4 removals)
	result, err := tagService.BatchRemoveTagsFromUsers(ctx, []uint64{user1.ID, user2.ID}, []uint64{tag1.ID, tag2.ID})

	require.NoError(t, err)
	assert.Equal(t, 4, result.TotalCount)
	assert.Equal(t, 4, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Len(t, result.SuccessItems, 2) // 2 unique users
	assert.Contains(t, result.SuccessItems, user1.ID)
	assert.Contains(t, result.SuccessItems, user2.ID)
	assert.Empty(t, result.FailedItems)

	// Verify database state - all relations removed
	assert.False(t, UserHasTag(t, db, user1.ID, tag1.ID))
	assert.False(t, UserHasTag(t, db, user1.ID, tag2.ID))
	assert.False(t, UserHasTag(t, db, user2.ID, tag1.ID))
	assert.False(t, UserHasTag(t, db, user2.ID, tag2.ID))

	relationCount := GetUserTagRelationCount(t, db)
	assert.Equal(t, int64(0), relationCount)
}

func TestUserTagService_BatchRemoveTagsFromUsers_UserNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test user
	user1 := CreateTestUser(t, db, "user1")
	nonExistentUserID := uint64(999999)

	// Create repository and service
	tagRepo := userrepository.NewUserTagRepository(db)
	userRepo := userrepository.NewUserRepository(db)
	tagService := userservice.NewUserTagService(tagRepo, userRepo, nil)

	// Create test tag
	tag1 := CreateTestUserTag(t, db, "VIP", "#FF6B6B", "VIP users")

	// Assign tag to user
	AssignTagToUser(t, db, user1.ID, tag1.ID)

	// Execute batch remove with non-existent user
	result, err := tagService.BatchRemoveTagsFromUsers(ctx, []uint64{user1.ID, nonExistentUserID}, []uint64{tag1.ID})

	require.NoError(t, err)
	assert.Equal(t, 2, result.TotalCount)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Len(t, result.FailedItems, 1)

	// Verify failed item
	failedItem := result.FailedItems[0]
	assert.Contains(t, failedItem.Message, "用户不存在")
}

func TestUserTagService_BatchRemoveTagsFromUsers_TagNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test user
	user1 := CreateTestUser(t, db, "user1")
	nonExistentTagID := uint64(999999)

	// Create repository and service
	tagRepo := userrepository.NewUserTagRepository(db)
	userRepo := userrepository.NewUserRepository(db)
	tagService := userservice.NewUserTagService(tagRepo, userRepo, nil)

	// Create test tag
	tag1 := CreateTestUserTag(t, db, "VIP", "#FF6B6B", "VIP users")

	// Assign tag to user
	AssignTagToUser(t, db, user1.ID, tag1.ID)

	// Execute batch remove with non-existent tag
	result, err := tagService.BatchRemoveTagsFromUsers(ctx, []uint64{user1.ID}, []uint64{tag1.ID, nonExistentTagID})

	require.NoError(t, err)
	assert.Equal(t, 2, result.TotalCount)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Len(t, result.FailedItems, 1)

	// Verify valid tag was removed
	assert.False(t, UserHasTag(t, db, user1.ID, tag1.ID))
}

func TestUserTagService_BatchRemoveTagsFromUsers_TagNotAssigned(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test users
	user1 := CreateTestUser(t, db, "user1")
	user2 := CreateTestUser(t, db, "user2")

	// Create repository and service
	tagRepo := userrepository.NewUserTagRepository(db)
	userRepo := userrepository.NewUserRepository(db)
	tagService := userservice.NewUserTagService(tagRepo, userRepo, nil)

	// Create test tags
	tag1 := CreateTestUserTag(t, db, "VIP", "#FF6B6B", "VIP users")
	tag2 := CreateTestUserTag(t, db, "Premium", "#45B7D1", "Premium users")

	// Assign only tag1 to user1
	AssignTagToUser(t, db, user1.ID, tag1.ID)

	// Execute batch remove - user1 doesn't have tag2, user2 doesn't have any
	result, err := tagService.BatchRemoveTagsFromUsers(ctx, []uint64{user1.ID, user2.ID}, []uint64{tag1.ID, tag2.ID})

	require.NoError(t, err)
	assert.Equal(t, 4, result.TotalCount)
	assert.Equal(t, 4, result.SuccessCount) // All succeed (idempotent for non-assigned tags)
	assert.Equal(t, 0, result.FailedCount)
	assert.Len(t, result.SuccessItems, 2)

	// Verify tag1 was removed from user1
	assert.False(t, UserHasTag(t, db, user1.ID, tag1.ID))
}

// ============================================================================
// Combined Complex Scenario Tests
// ============================================================================

func TestSensitiveWordUserBlockTag_ComplexBatchScenario(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repositories and services
	swRepo := sensitivewordrepo.NewSensitiveWordRepository(db)
	swService := sensitiveword.NewSensitiveWordService(swRepo)

	blockRepo := userblockrepo.NewUserBlockRepository(db)
	userRepo := userrepository.NewUserRepository(db)
	blockService := userblockservice.NewUserBlockService(blockRepo, userRepo)

	tagRepo := userrepository.NewUserTagRepository(db)
	tagService := userservice.NewUserTagService(tagRepo, userRepo, nil)

	// Create test users
	user1 := CreateTestUser(t, db, "user1")
	user2 := CreateTestUser(t, db, "user2")
	user3 := CreateTestUser(t, db, "user3")
	admin := CreateTestUser(t, db, "admin")

	// Test 1: Batch add sensitive words
	swReq := sensitiveword.BatchAddSensitiveWordsRequest{
		Words:     []string{"badword1", "badword2", "badword3"},
		Category:  model.SensitiveWordCategoryAbuse,
		Severity:  model.SensitiveWordSeverityHigh,
		MatchType: model.SensitiveWordMatchTypeExact,
	}
	swResult, err := swService.BatchAddSensitiveWords(ctx, swReq)
	require.NoError(t, err)
	assert.Equal(t, 3, swResult.SuccessCount)

	// Test 2: Batch create user blocks
	blockInput := []userblockservice.BlockInputItemForBatch{
		{
			BlockerID:   user1.ID,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   user2.ID,
			BlockedType: model.BlockUserTypeUser,
			Reason:      "Test block 1",
		},
		{
			BlockerID:   user2.ID,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   user3.ID,
			BlockedType: model.BlockUserTypeUser,
			Reason:      "Test block 2",
		},
	}
	blockResult, err := blockService.BatchBlock(ctx, blockInput)
	require.NoError(t, err)
	assert.Equal(t, 2, blockResult.SuccessCount)

	// Test 3: Create and assign tags
	tag1 := CreateTestUserTag(t, db, "Problem User", "#FF6B6B", "Users with issues")
	tag2 := CreateTestUserTag(t, db, "Blocked", "#4ECDC4", "Blocked users")

	// Assign tags to users
	assignResult, err := tagService.BatchAssignTagsToUsers(ctx, []uint64{user1.ID, user2.ID}, []uint64{tag1.ID, tag2.ID})
	require.NoError(t, err)
	assert.Equal(t, 4, assignResult.SuccessCount) // 2 users x 2 tags

	// Test 4: Batch update sensitive word status
	// Get the first sensitive word ID
	var sw model.SensitiveWord
	db.First(&sw)
	swStatusReq := sensitiveword.BatchUpdateSensitiveWordStatusRequest{
		IDs:      []uint64{sw.ID},
		IsActive: false,
	}
	swStatusResult, err := swService.BatchUpdateSensitiveWordStatus(ctx, swStatusReq)
	require.NoError(t, err)
	assert.Equal(t, 1, swStatusResult.SuccessCount)

	// Test 5: Batch unblock users
	var blocks []model.UserBlock
	db.Find(&blocks)
	blockIDs := make([]uint64, len(blocks))
	for i, b := range blocks {
		blockIDs[i] = b.ID
	}
	unblockResult, err := blockService.BatchUnblock(ctx, blockIDs, admin.ID, "Batch unblock")
	require.NoError(t, err)
	assert.Equal(t, len(blocks), unblockResult.SuccessCount)

	// Test 6: Batch remove tags
	removeResult, err := tagService.BatchRemoveTagsFromUsers(ctx, []uint64{user1.ID, user2.ID}, []uint64{tag1.ID, tag2.ID})
	require.NoError(t, err)
	assert.Equal(t, 4, removeResult.SuccessCount) // 2 users x 2 tags

	// Test 7: Batch delete sensitive words
	var sws []model.SensitiveWord
	db.Find(&sws)
	swIDs := make([]uint64, len(sws))
	for i, w := range sws {
		swIDs[i] = w.ID
	}
	deleteResult, err := swService.BatchDeleteSensitiveWords(ctx, swIDs)
	require.NoError(t, err)
	assert.Equal(t, len(sws), deleteResult.SuccessCount)

	// Final verification
	assert.Equal(t, int64(0), GetSensitiveWordCount(t, db))
	assert.Equal(t, int64(2), GetUserTagCount(t, db)) // Tags still exist, only relations deleted
	assert.Equal(t, int64(0), GetUserTagRelationCount(t, db))
}
