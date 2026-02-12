// Package integration provides integration tests for chat service.
package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	chatrepo "gamelink/internal/repository/chat"
	"gamelink/internal/service/chat"
	"gamelink/pkg/cache"
)

// setupChatService creates a ChatService with real database repositories.
func setupChatService(t *testing.T) (*chat.ChatService, context.Context) {
	t.Helper()
	db := SetupTestDB(t)
	ctx := context.Background()

	groupRepo := chatrepo.NewChatGroupRepository(db)
	memberRepo := chatrepo.NewChatMemberRepository(db)
	messageRepo := chatrepo.NewChatMessageRepository(db)
	reportRepo := chatrepo.NewChatReportRepository(db)
	memCache := cache.NewMemory()

	svc := chat.NewChatService(groupRepo, memberRepo, messageRepo, reportRepo, nil, memCache)
	return svc, ctx
}

// TestChatService_JoinGroup tests joining a chat group.
func TestChatService_JoinGroup(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc, ctx := setupChatService(t)

	// Create test user and group
	user := CreateUniqueTestUser(t, db, "chat_user")
	group := CreateTestChatGroup(t, db, "Test Group", model.ChatGroupTypePublic, nil)

	// Test joining group
	err := svc.JoinGroup(ctx, group.ID, user.ID, "TestNickname")
	require.NoError(t, err)

	// Verify membership
	member, err := svc.EnsureMembership(ctx, group.ID, user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, member.UserID)
	assert.Equal(t, "TestNickname", member.Nickname)
	assert.True(t, member.IsActive)
}

// TestChatService_JoinGroup_InactiveGroup tests joining an inactive group.
func TestChatService_JoinGroup_InactiveGroup(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc, ctx := setupChatService(t)

	user := CreateUniqueTestUser(t, db, "chat_user")
	group := CreateTestChatGroup(t, db, "Inactive Group", model.ChatGroupTypePublic, nil)

	// Deactivate the group
	db.Model(&model.ChatGroup{}).Where("id = ?", group.ID).Update("is_active", false)

	// Try to join - should fail
	err := svc.JoinGroup(ctx, group.ID, user.ID, "TestNickname")
	require.Error(t, err)
}

// TestChatService_LeaveGroup tests leaving a chat group.
func TestChatService_LeaveGroup(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc, ctx := setupChatService(t)

	user := CreateUniqueTestUser(t, db, "chat_user")
	group := CreateTestChatGroup(t, db, "Test Group", model.ChatGroupTypePublic, nil)

	// Join first
	err := svc.JoinGroup(ctx, group.ID, user.ID, "TestNickname")
	require.NoError(t, err)

	// Leave group
	err = svc.LeaveGroup(ctx, group.ID, user.ID)
	require.NoError(t, err)

	// Verify membership is inactive
	_, err = svc.EnsureMembership(ctx, group.ID, user.ID)
	require.Error(t, err) // Should fail because member is inactive
}

// TestChatService_LeaveGroup_NotMember tests leaving when not a member.
func TestChatService_LeaveGroup_NotMember(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc, ctx := setupChatService(t)

	user := CreateUniqueTestUser(t, db, "chat_user")
	group := CreateTestChatGroup(t, db, "Test Group", model.ChatGroupTypePublic, nil)

	// Try to leave without joining
	err := svc.LeaveGroup(ctx, group.ID, user.ID)
	require.Error(t, err)
}

// TestChatService_SendMessage tests sending a message.
func TestChatService_SendMessage(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc, ctx := setupChatService(t)

	user := CreateUniqueTestUser(t, db, "chat_user")
	group := CreateTestChatGroup(t, db, "Order Group", model.ChatGroupTypeOrder, nil)

	// Join group first
	err := svc.JoinGroup(ctx, group.ID, user.ID, "TestUser")
	require.NoError(t, err)

	// Send message
	input := chat.SendMessageInput{
		GroupID:     group.ID,
		SenderID:    user.ID,
		Content:     "Hello, World!",
		MessageType: model.ChatMessageTypeText,
	}
	msg, err := svc.SendMessage(ctx, input)
	require.NoError(t, err)
	assert.NotNil(t, msg)
	assert.Equal(t, "Hello, World!", msg.Content)
	assert.Equal(t, model.ChatMessageTypeText, msg.MessageType)
	// Order group messages are auto-approved
	assert.Equal(t, model.ChatMessageAuditApproved, msg.AuditStatus)
}

// TestChatService_SendMessage_PublicGroup tests sending message in public group (pending audit).
func TestChatService_SendMessage_PublicGroup(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc, ctx := setupChatService(t)

	user := CreateUniqueTestUser(t, db, "chat_user")
	group := CreateTestChatGroup(t, db, "Public Group", model.ChatGroupTypePublic, nil)

	// Join group
	err := svc.JoinGroup(ctx, group.ID, user.ID, "TestUser")
	require.NoError(t, err)

	// Send message
	input := chat.SendMessageInput{
		GroupID:     group.ID,
		SenderID:    user.ID,
		Content:     "Public message",
		MessageType: model.ChatMessageTypeText,
	}
	msg, err := svc.SendMessage(ctx, input)
	require.NoError(t, err)
	// Public group messages need audit
	assert.Equal(t, model.ChatMessageAuditPending, msg.AuditStatus)
}

// TestChatService_SendMessage_NotMember tests sending message when not a member.
func TestChatService_SendMessage_NotMember(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc, ctx := setupChatService(t)

	user := CreateUniqueTestUser(t, db, "chat_user")
	group := CreateTestChatGroup(t, db, "Test Group", model.ChatGroupTypeOrder, nil)

	// Try to send without joining
	input := chat.SendMessageInput{
		GroupID:     group.ID,
		SenderID:    user.ID,
		Content:     "Hello",
		MessageType: model.ChatMessageTypeText,
	}
	_, err := svc.SendMessage(ctx, input)
	require.Error(t, err)
}

// TestChatService_SendMessage_EmptyContent tests sending empty message.
func TestChatService_SendMessage_EmptyContent(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc, ctx := setupChatService(t)

	user := CreateUniqueTestUser(t, db, "chat_user")
	group := CreateTestChatGroup(t, db, "Test Group", model.ChatGroupTypeOrder, nil)

	// Join group
	err := svc.JoinGroup(ctx, group.ID, user.ID, "TestUser")
	require.NoError(t, err)

	// Try to send empty message
	input := chat.SendMessageInput{
		GroupID:     group.ID,
		SenderID:    user.ID,
		Content:     "",
		MessageType: model.ChatMessageTypeText,
	}
	_, err = svc.SendMessage(ctx, input)
	require.Error(t, err)
}

// TestChatService_ListMessages tests listing messages.
func TestChatService_ListMessages(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc, ctx := setupChatService(t)

	user := CreateUniqueTestUser(t, db, "chat_user")
	group := CreateTestChatGroup(t, db, "Order Group", model.ChatGroupTypeOrder, nil)

	// Join and send messages
	err := svc.JoinGroup(ctx, group.ID, user.ID, "TestUser")
	require.NoError(t, err)

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

	// List messages
	opts := chat.ListMessagesOptions{
		Page:     1,
		PageSize: 10,
	}
	messages, total, err := svc.ListMessages(ctx, user.ID, group.ID, opts)
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, messages, 5)
}

// TestChatService_ListMessages_NotMember tests listing messages when not a member.
func TestChatService_ListMessages_NotMember(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc, ctx := setupChatService(t)

	user := CreateUniqueTestUser(t, db, "chat_user")
	group := CreateTestChatGroup(t, db, "Test Group", model.ChatGroupTypeOrder, nil)

	// Try to list without being a member
	opts := chat.ListMessagesOptions{Page: 1, PageSize: 10}
	_, _, err := svc.ListMessages(ctx, user.ID, group.ID, opts)
	require.Error(t, err)
}

// TestChatService_ListUserGroups tests listing user's groups.
func TestChatService_ListUserGroups(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc, ctx := setupChatService(t)

	user := CreateUniqueTestUser(t, db, "chat_user")

	// Create and join multiple groups
	for i := 0; i < 3; i++ {
		group := CreateTestChatGroup(t, db, "Group "+string(rune('A'+i)), model.ChatGroupTypePublic, nil)
		err := svc.JoinGroup(ctx, group.ID, user.ID, "TestUser")
		require.NoError(t, err)
	}

	// List user's groups
	groups, total, err := svc.ListUserGroups(ctx, user.ID, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, groups, 3)
}

// TestChatService_MarkRead tests marking messages as read.
func TestChatService_MarkRead(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc, ctx := setupChatService(t)

	user := CreateUniqueTestUser(t, db, "chat_user")
	group := CreateTestChatGroup(t, db, "Order Group", model.ChatGroupTypeOrder, nil)

	// Join and send a message
	err := svc.JoinGroup(ctx, group.ID, user.ID, "TestUser")
	require.NoError(t, err)

	input := chat.SendMessageInput{
		GroupID:     group.ID,
		SenderID:    user.ID,
		Content:     "Test message",
		MessageType: model.ChatMessageTypeText,
	}
	msg, err := svc.SendMessage(ctx, input)
	require.NoError(t, err)

	// Mark as read
	err = svc.MarkRead(ctx, group.ID, user.ID, msg.ID)
	require.NoError(t, err)

	// Verify by checking membership
	member, err := svc.EnsureMembership(ctx, group.ID, user.ID)
	require.NoError(t, err)
	assert.NotNil(t, member.LastReadMessageID)
	assert.Equal(t, msg.ID, *member.LastReadMessageID)
}

// TestChatService_ApproveMessage tests approving a message.
func TestChatService_ApproveMessage(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc, ctx := setupChatService(t)

	user := CreateUniqueTestUser(t, db, "chat_user")
	moderator := CreateUniqueTestUser(t, db, "moderator")
	group := CreateTestChatGroup(t, db, "Public Group", model.ChatGroupTypePublic, nil)

	// Join and send message (will be pending in public group)
	err := svc.JoinGroup(ctx, group.ID, user.ID, "TestUser")
	require.NoError(t, err)

	input := chat.SendMessageInput{
		GroupID:     group.ID,
		SenderID:    user.ID,
		Content:     "Pending message",
		MessageType: model.ChatMessageTypeText,
	}
	msg, err := svc.SendMessage(ctx, input)
	require.NoError(t, err)
	assert.Equal(t, model.ChatMessageAuditPending, msg.AuditStatus)

	// Approve message
	err = svc.ApproveMessage(ctx, msg.ID, moderator.ID)
	require.NoError(t, err)

	// Verify status changed
	var updated model.ChatMessage
	db.First(&updated, msg.ID)
	assert.Equal(t, model.ChatMessageAuditApproved, updated.AuditStatus)
}

// TestChatService_RejectMessage tests rejecting a message.
func TestChatService_RejectMessage(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc, ctx := setupChatService(t)

	user := CreateUniqueTestUser(t, db, "chat_user")
	moderator := CreateUniqueTestUser(t, db, "moderator")
	group := CreateTestChatGroup(t, db, "Public Group", model.ChatGroupTypePublic, nil)

	// Join and send message
	err := svc.JoinGroup(ctx, group.ID, user.ID, "TestUser")
	require.NoError(t, err)

	input := chat.SendMessageInput{
		GroupID:     group.ID,
		SenderID:    user.ID,
		Content:     "Bad message",
		MessageType: model.ChatMessageTypeText,
	}
	msg, err := svc.SendMessage(ctx, input)
	require.NoError(t, err)

	// Reject message
	err = svc.RejectMessage(ctx, msg.ID, moderator.ID, "Inappropriate content")
	require.NoError(t, err)

	// Verify status changed
	var updated model.ChatMessage
	db.First(&updated, msg.ID)
	assert.Equal(t, model.ChatMessageAuditRejected, updated.AuditStatus)
}

// TestChatService_ReportMessage tests reporting a message.
func TestChatService_ReportMessage(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc, ctx := setupChatService(t)

	sender := CreateUniqueTestUser(t, db, "sender")
	reporter := CreateUniqueTestUser(t, db, "reporter")
	group := CreateTestChatGroup(t, db, "Order Group", model.ChatGroupTypeOrder, nil)

	// Both join group
	err := svc.JoinGroup(ctx, group.ID, sender.ID, "Sender")
	require.NoError(t, err)
	err = svc.JoinGroup(ctx, group.ID, reporter.ID, "Reporter")
	require.NoError(t, err)

	// Sender sends message
	input := chat.SendMessageInput{
		GroupID:     group.ID,
		SenderID:    sender.ID,
		Content:     "Offensive message",
		MessageType: model.ChatMessageTypeText,
	}
	msg, err := svc.SendMessage(ctx, input)
	require.NoError(t, err)

	// Reporter reports the message
	err = svc.ReportMessage(ctx, reporter.ID, msg.ID, "harassment", "Screenshot evidence")
	require.NoError(t, err)

	// Verify report was created
	var report model.ChatReport
	err = db.Where("message_id = ?", msg.ID).First(&report).Error
	require.NoError(t, err)
	assert.Equal(t, reporter.ID, report.ReporterID)
	assert.Equal(t, "harassment", report.Reason)
	assert.Equal(t, "pending", report.Status)
}

// TestChatService_OrderGroupWithOrder tests creating order-related chat group.
func TestChatService_OrderGroupWithOrder(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc, ctx := setupChatService(t)

	user := CreateUniqueTestUser(t, db, "order_user")
	playerUser := CreateUniqueTestUser(t, db, "player_user")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "chat_test_game")
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusConfirmed, 10000)

	// Create order chat group
	group := CreateTestChatGroup(t, db, "Order Chat", model.ChatGroupTypeOrder, &order.ID)

	// Both user and player join
	err := svc.JoinGroup(ctx, group.ID, user.ID, "Customer")
	require.NoError(t, err)
	err = svc.JoinGroup(ctx, group.ID, playerUser.ID, "Player")
	require.NoError(t, err)

	// Verify both are members
	_, err = svc.EnsureMembership(ctx, group.ID, user.ID)
	require.NoError(t, err)
	_, err = svc.EnsureMembership(ctx, group.ID, playerUser.ID)
	require.NoError(t, err)

	// Verify group is linked to order
	assert.NotNil(t, group.RelatedOrderID)
	assert.Equal(t, order.ID, *group.RelatedOrderID)
}
