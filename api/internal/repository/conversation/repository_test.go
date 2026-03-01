package conversation

import (
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
)

func setupConversationRepoDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.ChatGroup{}, &model.ChatGroupMember{}, &model.ChatMessage{}))

	return db
}

func createConversationFixture(t *testing.T, repo *Repository, userID uint64, agentID uint64) *model.ChatGroup {
	t.Helper()

	group := &model.ChatGroup{
		GroupName:  "客服会话",
		GroupType:  model.ChatGroupTypePrivate,
		CreatedBy:  userID,
		MaxMembers: 2,
		IsActive:   true,
	}

	members := []*model.ChatGroupMember{
		{
			GroupID:  group.ID,
			UserID:   userID,
			Role:     model.ChatMemberRoleOwner,
			Nickname: "用户",
			JoinedAt: time.Now(),
			IsActive: true,
		},
		{
			GroupID:  group.ID,
			UserID:   agentID,
			Role:     model.ChatMemberRoleMember,
			Nickname: "客服",
			JoinedAt: time.Now(),
			IsActive: true,
		},
	}

	require.NoError(t, repo.Create(context.Background(), group, members))
	return group
}

func TestConversationRepository_CRUDAndMessageFlow(t *testing.T) {
	db := setupConversationRepoDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	group := createConversationFixture(t, repo, 1001, 2001)

	loaded, err := repo.Get(ctx, group.ID)
	require.NoError(t, err)
	assert.Equal(t, group.ID, loaded.ID)
	assert.Len(t, loaded.Members, 2)

	groups, total, err := repo.ListByUser(ctx, 1001, 1, 20)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, groups, 1)

	msg := &model.ChatMessage{
		GroupID:     group.ID,
		SenderID:    1001,
		Content:     "你好，客服",
		MessageType: model.ChatMessageTypeText,
		AuditStatus: model.ChatMessageAuditApproved,
	}
	require.NoError(t, repo.CreateMessage(ctx, msg))

	messages, msgTotal, err := repo.ListMessages(ctx, group.ID, 1, 20, nil)
	require.NoError(t, err)
	assert.Equal(t, int64(1), msgTotal)
	assert.Len(t, messages, 1)
	assert.Equal(t, "你好，客服", messages[0].Content)

	loaded.GroupName = "已更新会话"
	require.NoError(t, repo.Update(ctx, loaded))
	updated, err := repo.Get(ctx, group.ID)
	require.NoError(t, err)
	assert.Equal(t, "已更新会话", updated.GroupName)

	require.NoError(t, repo.Delete(ctx, group.ID))
	closed, err := repo.Get(ctx, group.ID)
	require.NoError(t, err)
	assert.False(t, closed.IsActive)
}

func TestConversationRepository_FindAndCountByAgent(t *testing.T) {
	db := setupConversationRepoDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	groupA := createConversationFixture(t, repo, 1001, 2001)
	_ = createConversationFixture(t, repo, 1002, 2001)
	_ = createConversationFixture(t, repo, 1003, 2002)

	found, err := repo.FindActivePrivateByParticipants(ctx, 1001, []uint64{2001, 2002})
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, groupA.ID, found.ID)

	counts, err := repo.CountActiveByAgentIDs(ctx, []uint64{2001, 2002, 2003})
	require.NoError(t, err)
	assert.Equal(t, int64(2), counts[2001])
	assert.Equal(t, int64(1), counts[2002])
	assert.Equal(t, int64(0), counts[2003])
}
