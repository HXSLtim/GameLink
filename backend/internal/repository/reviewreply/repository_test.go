package reviewreply

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

func setupReviewReplyTest(t *testing.T) repository.ReviewReplyRepository {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&model.ReviewReply{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return NewReviewReplyRepository(db)
}

func TestReviewReplyRepository_Create(t *testing.T) {
	repo := setupReviewReplyTest(t)
	ctx := context.Background()

	reply := &model.ReviewReply{
		ReviewID: 1,
		AuthorID: 100,
		Content:  "Thank you for your feedback",
		Status:   "pending",
	}

	err := repo.Create(ctx, reply)
	assert.NoError(t, err)
	assert.NotZero(t, reply.ID)
}

func TestReviewReplyRepository_ListByReview(t *testing.T) {
	repo := setupReviewReplyTest(t)
	ctx := context.Background()

	// Create replies for different reviews
	for i := 0; i < 5; i++ {
		reply := &model.ReviewReply{
			ReviewID: uint64(i%2 + 1), // Alternate between review 1 and 2
			AuthorID: 100,
			Content:  "Reply " + string(rune(i)),
			Status:   "pending",
		}
		err := repo.Create(ctx, reply)
		assert.NoError(t, err)
	}

	// List replies for review 1
	replies, err := repo.ListByReview(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, replies, 3)

	// List replies for review 2
	replies, err = repo.ListByReview(ctx, 2)
	assert.NoError(t, err)
	assert.Len(t, replies, 2)

	// List replies for non-existent review
	replies, err = repo.ListByReview(ctx, 99999)
	assert.NoError(t, err)
	assert.Len(t, replies, 0)
}

func TestReviewReplyRepository_UpdateStatus(t *testing.T) {
	repo := setupReviewReplyTest(t)
	ctx := context.Background()

	// Create a reply
	reply := &model.ReviewReply{
		ReviewID: 1,
		AuthorID: 100,
		Content:  "Thank you for your feedback",
		Status:   "pending",
	}
	err := repo.Create(ctx, reply)
	assert.NoError(t, err)

	// Update status
	err = repo.UpdateStatus(ctx, reply.ID, "approved", "Looks good")
	assert.NoError(t, err)

	// Verify update
	replies, err := repo.ListByReview(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, replies, 1)
	assert.Equal(t, "approved", replies[0].Status)
	assert.Equal(t, "Looks good", replies[0].ModerationNote)
	assert.NotNil(t, replies[0].ModeratedAt)
}

func TestReviewReplyRepository_UpdateStatus_NotFound(t *testing.T) {
	repo := setupReviewReplyTest(t)
	ctx := context.Background()

	// Try to update non-existent reply
	err := repo.UpdateStatus(ctx, 99999, "approved", "Note")
	assert.Equal(t, repository.ErrNotFound, err)
}

func TestReviewReplyRepository_ListByReviewOrdering(t *testing.T) {
	repo := setupReviewReplyTest(t)
	ctx := context.Background()

	// Create replies in random order
	for i := 0; i < 3; i++ {
		reply := &model.ReviewReply{
			ReviewID: 1,
			AuthorID: uint64(100 + i),
			Content:  "Reply " + string(rune(i)),
			Status:   "pending",
		}
		err := repo.Create(ctx, reply)
		assert.NoError(t, err)
	}

	// List should be ordered by created_at ASC
	replies, err := repo.ListByReview(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, replies, 3)

	// Verify ordering (earlier created should come first)
	for i := 0; i < len(replies)-1; i++ {
		assert.True(t, replies[i].CreatedAt.Before(replies[i+1].CreatedAt) || replies[i].CreatedAt.Equal(replies[i+1].CreatedAt))
	}
}

func TestReviewReplyRepository_CreateMultiple(t *testing.T) {
	repo := setupReviewReplyTest(t)
	ctx := context.Background()

	// Create multiple replies
	for i := 0; i < 10; i++ {
		reply := &model.ReviewReply{
			ReviewID: 1,
			AuthorID: 100,
			Content:  "Reply " + string(rune(i)),
			Status:   "pending",
		}
		err := repo.Create(ctx, reply)
		assert.NoError(t, err)
		assert.NotZero(t, reply.ID)
	}

	// Verify all were created
	replies, err := repo.ListByReview(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, replies, 10)
}

func TestReviewReplyRepository_UpdateStatus_Rejected(t *testing.T) {
	repo := setupReviewReplyTest(t)
	ctx := context.Background()

	// Create a reply
	reply := &model.ReviewReply{
		ReviewID: 1,
		AuthorID: 100,
		Content:  "Inappropriate reply",
		Status:   "pending",
	}
	err := repo.Create(ctx, reply)
	assert.NoError(t, err)

	// Reject the reply
	err = repo.UpdateStatus(ctx, reply.ID, "rejected", "Violates community guidelines")
	assert.NoError(t, err)

	// Verify rejection
	replies, err := repo.ListByReview(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, replies, 1)
	assert.Equal(t, "rejected", replies[0].Status)
	assert.Equal(t, "Violates community guidelines", replies[0].ModerationNote)
}
