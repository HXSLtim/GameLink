package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/testutil"
)

func TestUserTagRepository(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.UserTag{}, &model.UserTagRelation{}, &model.User{})
	repo := NewUserTagRepository(db)
	ctx := context.Background()

	// seed users
	u1 := &model.User{Name: "U1", Email: "u1@example.com", Phone: "10000000001", Role: model.RoleUser, Status: model.UserStatusActive}
	u2 := &model.User{Name: "U2", Email: "u2@example.com", Phone: "10000000002", Role: model.RoleUser, Status: model.UserStatusActive}
	assert.NoError(t, db.WithContext(ctx).Create(u1).Error)
	assert.NoError(t, db.WithContext(ctx).Create(u2).Error)

	t.Run("CreateTag", func(t *testing.T) {
		tag := &model.UserTag{
			Name:        "VIP",
			Color:       "#FF6B6B",
			Description: "VIP用户",
			CreatedBy:   u1.ID,
		}
		err := repo.CreateTag(ctx, tag)
		assert.NoError(t, err)
		assert.NotZero(t, tag.ID)
	})

	// load tag id
	var vip model.UserTag
	assert.NoError(t, db.WithContext(ctx).Where("name = ?", "VIP").First(&vip).Error)

	t.Run("GetTag", func(t *testing.T) {
		got, err := repo.GetTag(ctx, vip.ID)
		assert.NoError(t, err)
		assert.Equal(t, vip.ID, got.ID)
		_, err = repo.GetTag(ctx, 99999)
		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("AddTagToUser idempotent", func(t *testing.T) {
		err := repo.AddTagToUser(ctx, u1.ID, vip.ID)
		assert.NoError(t, err)
		// add again should not duplicate
		err = repo.AddTagToUser(ctx, u1.ID, vip.ID)
		assert.NoError(t, err)
		var count int64
		assert.NoError(t, db.Model(&model.UserTagRelation{}).Where("user_id = ? AND tag_id = ?", u1.ID, vip.ID).Count(&count).Error)
		assert.Equal(t, int64(1), count)
	})

	t.Run("GetUserTags", func(t *testing.T) {
		tags, err := repo.GetUserTags(ctx, u1.ID)
		assert.NoError(t, err)
		assert.Len(t, tags, 1)
		assert.Equal(t, "VIP", tags[0].Name)
	})

	t.Run("GetUsersByTag pagination", func(t *testing.T) {
		// add for u2
		assert.NoError(t, repo.AddTagToUser(ctx, u2.ID, vip.ID))
		users, total, err := repo.GetUsersByTag(ctx, vip.ID, 1, 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, users, 1)
		// second page
		users2, total2, err := repo.GetUsersByTag(ctx, vip.ID, 2, 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total2)
		assert.Len(t, users2, 1)
	})

	t.Run("RemoveTagFromUser not found", func(t *testing.T) {
		err := repo.RemoveTagFromUser(ctx, 12345, vip.ID)
		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("UpdateTag and DeleteTag", func(t *testing.T) {
		vip.Color = "#000000"
		err := repo.UpdateTag(ctx, &vip)
		assert.NoError(t, err)
		// delete
		err = repo.DeleteTag(ctx, vip.ID)
		assert.NoError(t, err)
		// delete again should not found
		err = repo.DeleteTag(ctx, vip.ID)
		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})
}
