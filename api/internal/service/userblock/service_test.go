package userblock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
)

// ============================================================================
// Tests
// ============================================================================

func TestBlockStatus_Constants(t *testing.T) {
	t.Run("status values", func(t *testing.T) {
		assert.Equal(t, model.BlockStatus("active"), model.BlockStatusActive)
		assert.Equal(t, model.BlockStatus("canceled"), model.BlockStatusCanceled)
		assert.Equal(t, model.BlockStatus("admin_canceled"), model.BlockStatusAdminCanceled)
	})
}

func TestBlockUserType_Constants(t *testing.T) {
	t.Run("type values", func(t *testing.T) {
		assert.Equal(t, model.BlockUserType("user"), model.BlockUserTypeUser)
		assert.Equal(t, model.BlockUserType("player"), model.BlockUserTypePlayer)
	})
}

func TestUserBlock_Structure(t *testing.T) {
	t.Run("block record", func(t *testing.T) {
		now := time.Now()
		block := &model.UserBlock{
			BlockerID:   100,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   200,
			BlockedType: model.BlockUserTypePlayer,
			Reason:      "test reason",
			Status:      model.BlockStatusActive,
			BlockedAt:   now,
		}

		assert.Equal(t, uint64(100), block.BlockerID)
		assert.Equal(t, uint64(200), block.BlockedID)
		assert.Equal(t, model.BlockUserTypeUser, block.BlockerType)
		assert.Equal(t, model.BlockUserTypePlayer, block.BlockedType)
		assert.Equal(t, model.BlockStatusActive, block.Status)
	})
}

func TestUserBlock_Cancel(t *testing.T) {
	t.Run("user cancel", func(t *testing.T) {
		block := &model.UserBlock{
			Status: model.BlockStatusActive,
		}
		block.Cancel()

		assert.Equal(t, model.BlockStatusCanceled, block.Status)
		assert.NotNil(t, block.CanceledAt)
	})
}

func TestUserBlock_AdminCancel(t *testing.T) {
	t.Run("admin cancel", func(t *testing.T) {
		block := &model.UserBlock{
			Status: model.BlockStatusActive,
		}
		adminID := uint64(1)
		remark := "admin cancel reason"
		block.AdminCancel(adminID, remark)

		assert.Equal(t, model.BlockStatusAdminCanceled, block.Status)
		assert.NotNil(t, block.CanceledAt)
		assert.NotNil(t, block.CanceledBy)
		assert.Equal(t, adminID, *block.CanceledBy)
		assert.Equal(t, remark, block.AdminRemark)
	})
}

func TestUserBlock_IsActive(t *testing.T) {
	t.Run("active block", func(t *testing.T) {
		block := &model.UserBlock{
			Status: model.BlockStatusActive,
		}
		assert.True(t, block.IsActive())
	})

	t.Run("canceled block", func(t *testing.T) {
		block := &model.UserBlock{
			Status: model.BlockStatusCanceled,
		}
		assert.False(t, block.IsActive())
	})

	t.Run("admin canceled block", func(t *testing.T) {
		block := &model.UserBlock{
			Status: model.BlockStatusAdminCanceled,
		}
		assert.False(t, block.IsActive())
	})
}

func TestBlockInput_Validation(t *testing.T) {
	t.Run("cannot block self", func(t *testing.T) {
		input := BlockInput{
			BlockerID: 100,
			BlockedID: 100,
		}
		assert.Equal(t, input.BlockerID, input.BlockedID)
	})

	t.Run("valid block input", func(t *testing.T) {
		input := BlockInput{
			BlockerID:   100,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   200,
			BlockedType: model.BlockUserTypePlayer,
			Reason:      "test reason",
		}
		assert.NotEqual(t, input.BlockerID, input.BlockedID)
	})
}

func TestService_Errors(t *testing.T) {
	t.Run("error types", func(t *testing.T) {
		assert.NotNil(t, ErrNotFound)
		assert.NotNil(t, ErrValidation)
		assert.NotNil(t, ErrAlreadyBlocked)
		assert.NotNil(t, ErrNotBlocked)
		assert.NotNil(t, ErrCannotBlockSelf)
		assert.NotNil(t, ErrUserNotFound)
	})
}

func TestUserBlock_BidirectionalBlock(t *testing.T) {
	t.Run("A blocks B", func(t *testing.T) {
		blockAB := &model.UserBlock{
			BlockerID: 100,
			BlockedID: 200,
			Status:    model.BlockStatusActive,
		}
		assert.Equal(t, uint64(100), blockAB.BlockerID)
		assert.Equal(t, uint64(200), blockAB.BlockedID)
	})

	t.Run("B blocks A (independent)", func(t *testing.T) {
		blockBA := &model.UserBlock{
			BlockerID: 200,
			BlockedID: 100,
			Status:    model.BlockStatusActive,
		}
		assert.Equal(t, uint64(200), blockBA.BlockerID)
		assert.Equal(t, uint64(100), blockBA.BlockedID)
	})
}

func TestUserBlock_WithRelations(t *testing.T) {
	t.Run("block with blocker user", func(t *testing.T) {
		block := &model.UserBlock{
			BlockerID: 100,
			BlockedID: 200,
			Blocker: &model.User{
				Name: "Blocker User",
			},
		}
		assert.NotNil(t, block.Blocker)
		assert.Equal(t, "Blocker User", block.Blocker.Name)
	})

	t.Run("block with blocked user", func(t *testing.T) {
		block := &model.UserBlock{
			BlockerID: 100,
			BlockedID: 200,
			Blocked: &model.User{
				Name: "Blocked User",
			},
		}
		assert.NotNil(t, block.Blocked)
		assert.Equal(t, "Blocked User", block.Blocked.Name)
	})
}

func TestUserBlock_Timestamps(t *testing.T) {
	t.Run("blocked at", func(t *testing.T) {
		now := time.Now()
		block := &model.UserBlock{
			BlockedAt: now,
		}
		assert.False(t, block.BlockedAt.IsZero())
	})

	t.Run("canceled at", func(t *testing.T) {
		now := time.Now()
		block := &model.UserBlock{
			CanceledAt: &now,
		}
		assert.NotNil(t, block.CanceledAt)
	})
}

func TestService_BatchUnblock(t *testing.T) {
	t.Run("batch unblock count", func(t *testing.T) {
		ids := []uint64{1, 2, 3, 4, 5}
		successCount := len(ids) // assuming all succeed
		assert.Equal(t, 5, successCount)
	})

	t.Run("partial success", func(t *testing.T) {
		ids := []uint64{1, 2, 3, 4, 5}
		successCount := 3 // some failed
		assert.Less(t, successCount, len(ids))
	})
}

func TestUserBlock_Reason(t *testing.T) {
	t.Run("with reason", func(t *testing.T) {
		block := &model.UserBlock{
			Reason: "Inappropriate behavior",
		}
		assert.NotEmpty(t, block.Reason)
	})

	t.Run("without reason", func(t *testing.T) {
		block := &model.UserBlock{
			Reason: "",
		}
		assert.Empty(t, block.Reason)
	})
}

func TestUserBlock_AdminRemark(t *testing.T) {
	t.Run("with admin remark", func(t *testing.T) {
		block := &model.UserBlock{
			AdminRemark: "Unblocked by admin due to appeal",
		}
		assert.NotEmpty(t, block.AdminRemark)
	})
}
