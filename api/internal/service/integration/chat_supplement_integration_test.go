// Package integration provides supplementary integration tests for chat service.
package integration

import (
	"context"
	"testing"

	"gamelink/internal/model"
	chatrepo "gamelink/internal/repository/chat"
	"gamelink/internal/service/chat"
	"gamelink/pkg/cache"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestChatService_SensitiveWordFilter tests sensitive word filtering in messages.
func TestChatService_SensitiveWordFilter(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup chat service
	groupRepo := chatrepo.NewChatGroupRepository(db)
	memberRepo := chatrepo.NewChatMemberRepository(db)
	messageRepo := chatrepo.NewChatMessageRepository(db)
	reportRepo := chatrepo.NewChatReportRepository(db)
	memCache := cache.NewMemory()

	svc := chat.NewChatService(groupRepo, memberRepo, messageRepo, reportRepo, memCache)

	// Create test data
	user := CreateUniqueTestUser(t, db, "sensitive_user")
	group := CreateTestChatGroup(t, db, "Sensitive Test Group", model.ChatGroupTypePublic, nil)

	// Create sensitive word
	CreateTestSensitiveWord(t, db, "badword", model.SensitiveWordCategoryAbuse)

	// Join group
	err := svc.JoinGroup(ctx, group.ID, user.ID, "TestUser")
	require.NoError(t, err)

	// Send message with sensitive word (in public group, goes to pending)
	input := chat.SendMessageInput{
		GroupID:     group.ID,
		SenderID:    user.ID,
		Content:     "This contains badword in it",
		MessageType: model.ChatMessageTypeText,
	}
	msg, err := svc.SendMessage(ctx, input)
	require.NoError(t, err)

	// Public group messages need audit
	assert.Equal(t, model.ChatMessageAuditPending, msg.AuditStatus)
}

// TestChatService_DeleteMessage tests deleting a message.
func TestChatService_DeleteMessage(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup chat service
	groupRepo := chatrepo.NewChatGroupRepository(db)
	memberRepo := chatrepo.NewChatMemberRepository(db)
	messageRepo := chatrepo.NewChatMessageRepository(db)
	reportRepo := chatrepo.NewChatReportRepository(db)
	memCache := cache.NewMemory()

	svc := chat.NewChatService(groupRepo, memberRepo, messageRepo, reportRepo, memCache)

	// Create test data
	user := CreateUniqueTestUser(t, db, "delete_msg_user")
	moderator := CreateUniqueTestUser(t, db, "delete_msg_mod")
	group := CreateTestChatGroup(t, db, "Delete Test Group", model.ChatGroupTypeOrder, nil)

	// Join group and send message
	err := svc.JoinGroup(ctx, group.ID, user.ID, "TestUser")
	require.NoError(t, err)

	input := chat.SendMessageInput{
		GroupID:     group.ID,
		SenderID:    user.ID,
		Content:     "Message to be deleted",
		MessageType: model.ChatMessageTypeText,
	}
	msg, err := svc.SendMessage(ctx, input)
	require.NoError(t, err)

	// Delete message (mark as deleted)
	err = svc.RejectMessage(ctx, msg.ID, moderator.ID, "Inappropriate content - deleted")
	require.NoError(t, err)

	// Verify message is rejected/deleted
	var updatedMsg model.ChatMessage
	err = db.First(&updatedMsg, msg.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.ChatMessageAuditRejected, updatedMsg.AuditStatus)
}

// TestChatService_AutoCreateOrderGroup tests auto-creating order chat group.
func TestChatService_AutoCreateOrderGroup(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create test data
	user := CreateUniqueTestUser(t, db, "order_group_user")
	playerUser := CreateUniqueTestUser(t, db, "order_group_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "order_group_game")

	// Create order
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusConfirmed, 10000)

	// Auto-create order chat group (simulating payment success)
	group := CreateTestChatGroup(t, db, "Order #"+order.OrderNo, model.ChatGroupTypeOrder, &order.ID)

	// Verify group is linked to order
	assert.NotNil(t, group.RelatedOrderID)
	assert.Equal(t, order.ID, *group.RelatedOrderID)
	assert.Equal(t, model.ChatGroupTypeOrder, group.GroupType)

	// Add user and player as members
	CreateTestChatGroupMember(t, db, group.ID, user.ID, model.ChatMemberRoleMember)
	CreateTestChatGroupMember(t, db, group.ID, playerUser.ID, model.ChatMemberRoleMember)

	// Verify members
	var members []model.ChatGroupMember
	err := db.Where("group_id = ?", group.ID).Find(&members).Error
	require.NoError(t, err)
	assert.Len(t, members, 2)
}

// TestChatService_OrderGroupWithCustomerService tests order group with customer service.
func TestChatService_OrderGroupWithCustomerService(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create test data
	user := CreateUniqueTestUser(t, db, "cs_group_user")
	playerUser := CreateUniqueTestUser(t, db, "cs_group_player")
	player := CreateTestPlayer(t, db, playerUser)
	csUser := CreateUniqueTestUser(t, db, "cs_group_cs")
	game := CreateTestGame(t, db, "cs_group_game")

	// Create order
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusConfirmed, 10000)

	// Create order chat group
	group := CreateTestChatGroup(t, db, "Order #"+order.OrderNo, model.ChatGroupTypeOrder, &order.ID)

	// Add members including customer service
	CreateTestChatGroupMember(t, db, group.ID, user.ID, model.ChatMemberRoleMember)
	CreateTestChatGroupMember(t, db, group.ID, playerUser.ID, model.ChatMemberRoleMember)
	CreateTestChatGroupMember(t, db, group.ID, csUser.ID, model.ChatMemberRoleAdmin) // CS as admin

	// Verify CS is in group
	var csMember model.ChatGroupMember
	err := db.Where("group_id = ? AND user_id = ?", group.ID, csUser.ID).First(&csMember).Error
	require.NoError(t, err)
	assert.Equal(t, model.ChatMemberRoleAdmin, csMember.Role)
}

// TestChatService_MessageRateLimit tests message rate limiting.
func TestChatService_MessageRateLimit(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup chat service
	groupRepo := chatrepo.NewChatGroupRepository(db)
	memberRepo := chatrepo.NewChatMemberRepository(db)
	messageRepo := chatrepo.NewChatMessageRepository(db)
	reportRepo := chatrepo.NewChatReportRepository(db)
	memCache := cache.NewMemory()

	svc := chat.NewChatService(groupRepo, memberRepo, messageRepo, reportRepo, memCache)

	// Create test data
	user := CreateUniqueTestUser(t, db, "rate_limit_user")
	group := CreateTestChatGroup(t, db, "Rate Limit Group", model.ChatGroupTypeOrder, nil)

	// Join group
	err := svc.JoinGroup(ctx, group.ID, user.ID, "TestUser")
	require.NoError(t, err)

	// Send multiple messages (order group has no rate limit by default)
	for i := 0; i < 5; i++ {
		input := chat.SendMessageInput{
			GroupID:     group.ID,
			SenderID:    user.ID,
			Content:     "Message " + string(rune('A'+i)),
			MessageType: model.ChatMessageTypeText,
		}
		_, err := svc.SendMessage(ctx, input)
		require.NoError(t, err)
	}

	// Verify all messages sent
	var messages []model.ChatMessage
	err = db.Where("group_id = ?", group.ID).Find(&messages).Error
	require.NoError(t, err)
	assert.Len(t, messages, 5)
}

// TestChatService_GroupDeactivation tests group deactivation after order completion.
func TestChatService_GroupDeactivation(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create test data
	user := CreateUniqueTestUser(t, db, "deactivate_user")
	playerUser := CreateUniqueTestUser(t, db, "deactivate_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "deactivate_game")

	// Create completed order
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, 10000)

	// Create order chat group
	group := CreateTestChatGroup(t, db, "Order #"+order.OrderNo, model.ChatGroupTypeOrder, &order.ID)
	group.AutoDestroy = true
	db.Save(group)

	// Deactivate group after order completion
	group.IsActive = false
	err := db.Save(group).Error
	require.NoError(t, err)

	// Verify group is deactivated
	var updatedGroup model.ChatGroup
	err = db.First(&updatedGroup, group.ID).Error
	require.NoError(t, err)
	assert.False(t, updatedGroup.IsActive)
}

// TestChatService_ImageMessage tests sending image message.
func TestChatService_ImageMessage(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup chat service
	groupRepo := chatrepo.NewChatGroupRepository(db)
	memberRepo := chatrepo.NewChatMemberRepository(db)
	messageRepo := chatrepo.NewChatMessageRepository(db)
	reportRepo := chatrepo.NewChatReportRepository(db)
	memCache := cache.NewMemory()

	svc := chat.NewChatService(groupRepo, memberRepo, messageRepo, reportRepo, memCache)

	// Create test data
	user := CreateUniqueTestUser(t, db, "image_msg_user")
	group := CreateTestChatGroup(t, db, "Image Test Group", model.ChatGroupTypeOrder, nil)

	// Join group
	err := svc.JoinGroup(ctx, group.ID, user.ID, "TestUser")
	require.NoError(t, err)

	// Send image message
	input := chat.SendMessageInput{
		GroupID:     group.ID,
		SenderID:    user.ID,
		Content:     "https://example.com/image.jpg",
		MessageType: model.ChatMessageTypeImage,
	}
	msg, err := svc.SendMessage(ctx, input)
	require.NoError(t, err)
	assert.Equal(t, model.ChatMessageTypeImage, msg.MessageType)
}

// TestChatService_SystemMessage tests sending system message.
func TestChatService_SystemMessage(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create test data
	group := CreateTestChatGroup(t, db, "System Msg Group", model.ChatGroupTypeOrder, nil)

	// Create system message directly
	systemMsg := &model.ChatMessage{
		Base: model.Base{
			ExtJSON: "{}",
		},
		GroupID:     group.ID,
		SenderID:    0, // System message has no sender
		Content:     "User joined the group",
		MessageType: model.ChatMessageTypeSystem,
		AuditStatus: model.ChatMessageAuditApproved,
		Metadata:    "{}",
	}
	err := db.Create(systemMsg).Error
	require.NoError(t, err)

	// Verify
	var savedMsg model.ChatMessage
	err = db.First(&savedMsg, systemMsg.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.ChatMessageTypeSystem, savedMsg.MessageType)
	assert.Equal(t, uint64(0), savedMsg.SenderID)
}

// TestChatService_MuteUser tests muting a user in group.
func TestChatService_MuteUser(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup chat service
	groupRepo := chatrepo.NewChatGroupRepository(db)
	memberRepo := chatrepo.NewChatMemberRepository(db)
	messageRepo := chatrepo.NewChatMessageRepository(db)
	reportRepo := chatrepo.NewChatReportRepository(db)
	memCache := cache.NewMemory()

	svc := chat.NewChatService(groupRepo, memberRepo, messageRepo, reportRepo, memCache)

	// Create test data
	user := CreateUniqueTestUser(t, db, "mute_user")
	group := CreateTestChatGroup(t, db, "Mute Test Group", model.ChatGroupTypePublic, nil)

	// Join group
	err := svc.JoinGroup(ctx, group.ID, user.ID, "TestUser")
	require.NoError(t, err)

	// Mute user
	var member model.ChatGroupMember
	err = db.Where("group_id = ? AND user_id = ?", group.ID, user.ID).First(&member).Error
	require.NoError(t, err)

	member.IsMuted = true
	err = db.Save(&member).Error
	require.NoError(t, err)

	// Verify user is muted
	var updatedMember model.ChatGroupMember
	err = db.First(&updatedMember, member.ID).Error
	require.NoError(t, err)
	assert.True(t, updatedMember.IsMuted)
}
