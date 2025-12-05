package model_test

import (
	"testing"

	"gamelink/internal/model"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestReviewMigration(t *testing.T) {
	// Create an in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Run AutoMigrate
	err = db.AutoMigrate(&model.Review{})
	assert.NoError(t, err)

	// Verify the table was created
	assert.True(t, db.Migrator().HasTable(&model.Review{}))

	// Verify all columns exist
	assert.True(t, db.Migrator().HasColumn(&model.Review{}, "id"))
	assert.True(t, db.Migrator().HasColumn(&model.Review{}, "order_id"))
	assert.True(t, db.Migrator().HasColumn(&model.Review{}, "user_id"))
	assert.True(t, db.Migrator().HasColumn(&model.Review{}, "player_id"))
	assert.True(t, db.Migrator().HasColumn(&model.Review{}, "score"))
	assert.True(t, db.Migrator().HasColumn(&model.Review{}, "content"))
	assert.True(t, db.Migrator().HasColumn(&model.Review{}, "status"))
	assert.True(t, db.Migrator().HasColumn(&model.Review{}, "is_reported"))
	assert.True(t, db.Migrator().HasColumn(&model.Review{}, "images"))
	assert.True(t, db.Migrator().HasColumn(&model.Review{}, "rejection_reason"))
	assert.True(t, db.Migrator().HasColumn(&model.Review{}, "created_at"))
	assert.True(t, db.Migrator().HasColumn(&model.Review{}, "updated_at"))
}

func TestReviewDatabaseOperations(t *testing.T) {
	// Create an in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Run AutoMigrate
	err = db.AutoMigrate(&model.Review{})
	assert.NoError(t, err)

	// Test Create
	review := &model.Review{
		OrderID:         100,
		UserID:          200,
		PlayerID:        300,
		Score:           5,
		Content:         "Great service!",
		Status:          model.ReviewStatusPending,
		IsReported:      false,
		Images:          model.StringArray{"https://example.com/image1.jpg", "https://example.com/image2.jpg"},
		RejectionReason: "",
	}

	err = db.Create(review).Error
	assert.NoError(t, err)
	assert.NotZero(t, review.ID)

	// Test Read
	var retrieved model.Review
	err = db.First(&retrieved, review.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, review.OrderID, retrieved.OrderID)
	assert.Equal(t, review.UserID, retrieved.UserID)
	assert.Equal(t, review.PlayerID, retrieved.PlayerID)
	assert.Equal(t, review.Score, retrieved.Score)
	assert.Equal(t, review.Content, retrieved.Content)
	assert.Equal(t, review.Status, retrieved.Status)
	assert.Equal(t, review.IsReported, retrieved.IsReported)
	assert.Len(t, retrieved.Images, 2)
	assert.Equal(t, "https://example.com/image1.jpg", retrieved.Images[0])
	assert.Equal(t, "https://example.com/image2.jpg", retrieved.Images[1])

	// Test Update
	err = db.Model(&retrieved).Updates(map[string]interface{}{
		"status":           model.ReviewStatusApproved,
		"is_reported":      true,
		"rejection_reason": "Test reason",
	}).Error
	assert.NoError(t, err)

	// Verify Update
	var updated model.Review
	err = db.First(&updated, review.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, model.ReviewStatusApproved, updated.Status)
	assert.True(t, updated.IsReported)
	assert.Equal(t, "Test reason", updated.RejectionReason)

	// Test Delete
	err = db.Delete(&updated).Error
	assert.NoError(t, err)

	// Verify Delete (soft delete)
	var deleted model.Review
	err = db.First(&deleted, review.ID).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestReviewStatusIndex(t *testing.T) {
	// Create an in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Run AutoMigrate
	err = db.AutoMigrate(&model.Review{})
	assert.NoError(t, err)

	// Create multiple reviews with different statuses
	reviews := []model.Review{
		{OrderID: 1, UserID: 1, PlayerID: 1, Score: 5, Status: model.ReviewStatusPending},
		{OrderID: 2, UserID: 2, PlayerID: 2, Score: 4, Status: model.ReviewStatusApproved},
		{OrderID: 3, UserID: 3, PlayerID: 3, Score: 3, Status: model.ReviewStatusRejected},
		{OrderID: 4, UserID: 4, PlayerID: 4, Score: 5, Status: model.ReviewStatusPending},
		{OrderID: 5, UserID: 5, PlayerID: 5, Score: 2, Status: model.ReviewStatusDeleted},
	}

	for i := range reviews {
		err = db.Create(&reviews[i]).Error
		assert.NoError(t, err)
	}

	// Query by status (should use index)
	var pendingReviews []model.Review
	err = db.Where("status = ?", model.ReviewStatusPending).Find(&pendingReviews).Error
	assert.NoError(t, err)
	assert.Len(t, pendingReviews, 2)

	// Query by is_reported (should use index)
	var reportedReviews []model.Review
	err = db.Where("is_reported = ?", false).Find(&reportedReviews).Error
	assert.NoError(t, err)
	assert.Len(t, reportedReviews, 5)
}

func TestReviewImagesJSONStorage(t *testing.T) {
	// Create an in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Run AutoMigrate
	err = db.AutoMigrate(&model.Review{})
	assert.NoError(t, err)

	// Test with multiple images
	review := &model.Review{
		OrderID:  100,
		UserID:   200,
		PlayerID: 300,
		Score:    5,
		Content:  "Great service with photos!",
		Status:   model.ReviewStatusApproved,
		Images: model.StringArray{
			"https://example.com/image1.jpg",
			"https://example.com/image2.jpg",
			"https://example.com/image3.jpg",
		},
	}

	err = db.Create(review).Error
	assert.NoError(t, err)

	// Retrieve and verify
	var retrieved model.Review
	err = db.First(&retrieved, review.ID).Error
	assert.NoError(t, err)
	assert.Len(t, retrieved.Images, 3)
	assert.Equal(t, "https://example.com/image1.jpg", retrieved.Images[0])
	assert.Equal(t, "https://example.com/image2.jpg", retrieved.Images[1])
	assert.Equal(t, "https://example.com/image3.jpg", retrieved.Images[2])

	// Test with empty images
	review2 := &model.Review{
		OrderID:  101,
		UserID:   201,
		PlayerID: 301,
		Score:    4,
		Content:  "Good service without photos",
		Status:   model.ReviewStatusApproved,
		Images:   model.StringArray{},
	}

	err = db.Create(review2).Error
	assert.NoError(t, err)

	var retrieved2 model.Review
	err = db.First(&retrieved2, review2.ID).Error
	assert.NoError(t, err)
	assert.Len(t, retrieved2.Images, 0)
}
