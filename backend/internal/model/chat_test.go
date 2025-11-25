package model_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gamelink/internal/model"
)

func TestChatGroupModel(t *testing.T) {
	now := time.Now()
	deactivatedAt := now.Add(24 * time.Hour)
	orderID := uint64(100)

	chatGroup := &model.ChatGroup{
		Base: model.Base{
			ID: 1,
		},
		GroupName:      "游戏讨论组",
		GroupType:      model.ChatGroupTypePublic,
		RelatedOrderID: &orderID,
		CreatedBy:      200,
		MaxMembers:     50,
		IsActive:       true,
		AutoDestroy:    false,
		DeactivatedAt:  &deactivatedAt,
		AvatarURL:      "https://example.com/avatar.jpg",
		Description:    "这是一个关于游戏的讨论组",
		Settings:       `{"allowFile": true, "maxFileSize": 10485760}`,
	}

	assert.Equal(t, uint64(1), chatGroup.ID)
	assert.Equal(t, "游戏讨论组", chatGroup.GroupName)
	assert.Equal(t, model.ChatGroupTypePublic, chatGroup.GroupType)
	assert.Equal(t, &orderID, chatGroup.RelatedOrderID)
	assert.Equal(t, uint64(200), chatGroup.CreatedBy)
	assert.Equal(t, 50, chatGroup.MaxMembers)
	assert.True(t, chatGroup.IsActive)
	assert.False(t, chatGroup.AutoDestroy)
	assert.Equal(t, &deactivatedAt, chatGroup.DeactivatedAt)
	assert.Equal(t, "https://example.com/avatar.jpg", chatGroup.AvatarURL)
	assert.Equal(t, "这是一个关于游戏的讨论组", chatGroup.Description)
	assert.Equal(t, `{"allowFile": true, "maxFileSize": 10485760}`, chatGroup.Settings)
}

func TestChatGroupJSONSerialization(t *testing.T) {
	now := time.Now()
	orderID := uint64(100)

	chatGroup := &model.ChatGroup{
		Base: model.Base{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		GroupName:      "订单聊天",
		GroupType:      model.ChatGroupTypeOrder,
		RelatedOrderID: &orderID,
		CreatedBy:      200,
		MaxMembers:     10,
		IsActive:       true,
		AutoDestroy:    true,
		AvatarURL:      "https://example.com/order_avatar.png",
		Description:    "订单相关聊天",
		Settings:       `{"allowFile": false}`,
	}

	// 序列化
	data, err := json.Marshal(chatGroup)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "订单聊天")
	assert.Contains(t, string(data), "order")

	// 反序列化
	var decoded model.ChatGroup
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, chatGroup.ID, decoded.ID)
	assert.Equal(t, chatGroup.GroupName, decoded.GroupName)
	assert.Equal(t, chatGroup.GroupType, decoded.GroupType)
	assert.Equal(t, *chatGroup.RelatedOrderID, *decoded.RelatedOrderID)
	assert.Equal(t, chatGroup.CreatedBy, decoded.CreatedBy)
}

func TestChatGroupTableName(t *testing.T) {
	chatGroup := model.ChatGroup{}
	tableName := chatGroup.TableName()
	assert.Equal(t, "chat_groups", tableName)
}

func TestChatGroupConstants(t *testing.T) {
	// 测试聊天组类型常量
	assert.Equal(t, model.ChatGroupType("public"), model.ChatGroupTypePublic)
	assert.Equal(t, model.ChatGroupType("order"), model.ChatGroupTypeOrder)

	// 测试消息类型常量
	assert.Equal(t, model.ChatMessageType("text"), model.ChatMessageTypeText)
	assert.Equal(t, model.ChatMessageType("image"), model.ChatMessageTypeImage)
	assert.Equal(t, model.ChatMessageType("file"), model.ChatMessageTypeFile)
	assert.Equal(t, model.ChatMessageType("system"), model.ChatMessageTypeSystem)

	// 测试审核状态常量
	assert.Equal(t, model.ChatMessageAuditStatus("pending"), model.ChatMessageAuditPending)
	assert.Equal(t, model.ChatMessageAuditStatus("approved"), model.ChatMessageAuditApproved)
	assert.Equal(t, model.ChatMessageAuditStatus("rejected"), model.ChatMessageAuditRejected)
}

func TestChatGroupZeroValues(t *testing.T) {
	chatGroup := &model.ChatGroup{
		GroupName:      "",
		GroupType:      "",
		RelatedOrderID: nil,
		CreatedBy:      0,
		MaxMembers:     0,
		IsActive:       false,
		AutoDestroy:    false,
		DeactivatedAt:  nil,
		AvatarURL:      "",
		Description:    "",
		Settings:       "",
	}

	assert.Equal(t, "", chatGroup.GroupName)
	assert.Equal(t, model.ChatGroupType(""), chatGroup.GroupType)
	assert.Nil(t, chatGroup.RelatedOrderID)
	assert.Equal(t, uint64(0), chatGroup.CreatedBy)
	assert.Equal(t, 0, chatGroup.MaxMembers)
	assert.False(t, chatGroup.IsActive)
	assert.False(t, chatGroup.AutoDestroy)
	assert.Nil(t, chatGroup.DeactivatedAt)
	assert.Equal(t, "", chatGroup.AvatarURL)
	assert.Equal(t, "", chatGroup.Description)
	assert.Equal(t, "", chatGroup.Settings)
}

func TestChatGroupEdgeCases(t *testing.T) {
	// 测试长文本
	longName := "这是一个非常长的聊天组名称，用于测试字符串长度的边界情况，可能包含很多描述性信息"
	longDescription := "这是一个非常长的聊天组描述，可以包含很多详细信息，比如聊天组的具体用途、适用场景、规则说明等等。这种长文本测试可以确保我们的模型能够处理各种长度的输入。"
	longAvatarURL := "https://example.com/very/long/avatar/url/path/that/might/be/used/for/testing/purposes/and/should/handle/long/paths/without/issues.jpg"
	longSettings := `{"allowFile": true, "maxFileSize": 10485760, "allowImage": true, "maxImageSize": 5242880, "allowVideo": false, "allowAudio": true, "maxMembers": 1000, "requireApproval": false, "autoModerate": true, "keywords": ["游戏", "娱乐", "聊天", "交友"]}`

	chatGroup1 := &model.ChatGroup{
		GroupName:   longName,
		Description: longDescription,
		AvatarURL:   longAvatarURL,
		Settings:    longSettings,
	}

	assert.Equal(t, longName, chatGroup1.GroupName)
	assert.Equal(t, longDescription, chatGroup1.Description)
	assert.Equal(t, longAvatarURL, chatGroup1.AvatarURL)
	assert.Equal(t, longSettings, chatGroup1.Settings)

	// 测试特殊字符
	chatGroup2 := &model.ChatGroup{
		GroupName:   "聊天组@#$%^&*()_+-=[]{}|\\",
		Description: "描述包含特殊字符：<>{}[]|\\\"quotes\" and 'apostrophes' and @#$%^&*()_+-=[]{}|;':\",./<>?😊🚀",
		AvatarURL:   "https://example.com/avatar@special#chars.png",
		Settings:    `{"specialChars": "@#$%^&*()_+-=[]{}|\\\"quotes\""}`,
	}
	assert.Equal(t, "聊天组@#$%^&*()_+-=[]{}|\\", chatGroup2.GroupName)
	assert.Contains(t, chatGroup2.Description, "特殊字符")
	assert.Equal(t, "https://example.com/avatar@special#chars.png", chatGroup2.AvatarURL)

	// 测试大数值
	chatGroup3 := &model.ChatGroup{
		CreatedBy:  ^uint64(0), // 最大uint64值
		MaxMembers: ^int(0) >> 1, // 最大int值
	}
	assert.Equal(t, ^uint64(0), chatGroup3.CreatedBy)
	assert.Equal(t, ^int(0)>>1, chatGroup3.MaxMembers)

	// 测试零值
	chatGroup4 := &model.ChatGroup{
		MaxMembers: 0,
	}
	assert.Equal(t, 0, chatGroup4.MaxMembers)
}

func TestChatGroupWithOrderID(t *testing.T) {
	orderID := uint64(500)

	chatGroup := &model.ChatGroup{
		GroupType:      model.ChatGroupTypeOrder,
		RelatedOrderID: &orderID,
	}

	assert.Equal(t, model.ChatGroupTypeOrder, chatGroup.GroupType)
	assert.NotNil(t, chatGroup.RelatedOrderID)
	assert.Equal(t, uint64(500), *chatGroup.RelatedOrderID)
}

func TestChatGroupWithoutOrderID(t *testing.T) {
	chatGroup := &model.ChatGroup{
		GroupType:      model.ChatGroupTypePublic,
		RelatedOrderID: nil,
	}

	assert.Equal(t, model.ChatGroupTypePublic, chatGroup.GroupType)
	assert.Nil(t, chatGroup.RelatedOrderID)
}

func TestChatGroupMembers(t *testing.T) {
	member1 := model.ChatGroupMember{
		Base:     model.Base{ID: 1},
		GroupID:  100,
		UserID:   200,
		Role:     "admin",
		Nickname: "管理员",
	}

	member2 := model.ChatGroupMember{
		Base:     model.Base{ID: 2},
		GroupID:  100,
		UserID:   300,
		Role:     "member",
		Nickname: "普通成员",
	}

	chatGroup := &model.ChatGroup{
		GroupName: "测试群组",
		Members:   []model.ChatGroupMember{member1, member2},
	}

	assert.Len(t, chatGroup.Members, 2)
	assert.Equal(t, "管理员", chatGroup.Members[0].Nickname)
	assert.Equal(t, "普通成员", chatGroup.Members[1].Nickname)
}

func TestChatGroupJSONFields(t *testing.T) {
	orderID := uint64(100)

	chatGroup := &model.ChatGroup{
		Base: model.Base{
			ID: 1,
		},
		GroupName:      "测试群组",
		GroupType:      model.ChatGroupTypePublic,
		RelatedOrderID: &orderID,
		CreatedBy:      200,
		MaxMembers:     100,
		IsActive:       true,
		AutoDestroy:    false,
		AvatarURL:      "https://example.com/avatar.jpg",
		Description:    "这是一个测试群组",
		Settings:       `{"test": "value"}`,
	}

	data, err := json.Marshal(chatGroup)
	assert.NoError(t, err)

	// 验证JSON结构
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// 检查必需的字段
	assert.Contains(t, result, "id")
	assert.Contains(t, result, "groupName")
	assert.Contains(t, result, "groupType")
	assert.Contains(t, result, "relatedOrderId")
	assert.Contains(t, result, "createdBy")
	assert.Contains(t, result, "maxMembers")
	assert.Contains(t, result, "isActive")
	assert.Contains(t, result, "autoDestroy")
	assert.Contains(t, result, "avatarUrl")
	assert.Contains(t, result, "description")
	assert.Contains(t, result, "settings")

	// 验证值
	assert.Equal(t, float64(1), result["id"])
	assert.Equal(t, "测试群组", result["groupName"])
	assert.Equal(t, "public", result["groupType"])
	assert.Equal(t, float64(100), result["relatedOrderId"])
	assert.Equal(t, float64(200), result["createdBy"])
	assert.Equal(t, float64(100), result["maxMembers"])
	assert.Equal(t, true, result["isActive"])
	assert.Equal(t, false, result["autoDestroy"])
	assert.Equal(t, "https://example.com/avatar.jpg", result["avatarUrl"])
	assert.Equal(t, "这是一个测试群组", result["description"])
}

func TestChatGroupMemberModel(t *testing.T) {
	now := time.Now()
	lastReadAt := now.Add(-5 * time.Minute)
	lastReadMessageID := uint64(150)

	member := &model.ChatGroupMember{
		Base: model.Base{
			ID: 1,
		},
		GroupID:           100,
		UserID:            200,
		Role:              "admin",
		Nickname:          "群管理员",
		JoinedAt:          now,
		LastReadAt:        &lastReadAt,
		LastReadMessageID: &lastReadMessageID,
		IsMuted:           false,
		IsActive:          true,
	}

	assert.Equal(t, uint64(1), member.ID)
	assert.Equal(t, uint64(100), member.GroupID)
	assert.Equal(t, uint64(200), member.UserID)
	assert.Equal(t, "admin", member.Role)
	assert.Equal(t, "群管理员", member.Nickname)
	assert.Equal(t, now, member.JoinedAt)
	assert.Equal(t, &lastReadAt, member.LastReadAt)
	assert.Equal(t, &lastReadMessageID, member.LastReadMessageID)
	assert.False(t, member.IsMuted)
	assert.True(t, member.IsActive)
}

func TestChatGroupMemberJSONSerialization(t *testing.T) {
	member := &model.ChatGroupMember{
		GroupID:  100,
		UserID:   200,
		Role:     "member",
		Nickname: "普通成员",
		IsMuted:  false,
		IsActive: true,
	}

	// 序列化
	data, err := json.Marshal(member)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "普通成员")
	assert.Contains(t, string(data), "member")

	// 反序列化
	var decoded model.ChatGroupMember
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, member.ID, decoded.ID)
	assert.Equal(t, member.GroupID, decoded.GroupID)
	assert.Equal(t, member.UserID, decoded.UserID)
	assert.Equal(t, member.Role, decoded.Role)
	assert.Equal(t, member.Nickname, decoded.Nickname)
}

func TestChatGroupMemberZeroValues(t *testing.T) {
	member := &model.ChatGroupMember{
		GroupID:           0,
		UserID:            0,
		Role:              "",
		Nickname:          "",
		LastReadAt:        nil,
		LastReadMessageID: nil,
		IsMuted:           false,
		IsActive:          false,
	}

	assert.Equal(t, uint64(0), member.GroupID)
	assert.Equal(t, uint64(0), member.UserID)
	assert.Equal(t, "", member.Role)
	assert.Equal(t, "", member.Nickname)
	assert.Nil(t, member.LastReadAt)
	assert.Nil(t, member.LastReadMessageID)
	assert.False(t, member.IsMuted)
	assert.False(t, member.IsActive)
}

func TestChatGroupMemberEdgeCases(t *testing.T) {
	// 测试长昵称
	longNickname := "这是一个非常长的群成员昵称，用于测试字符串长度的边界情况，可能包含很多描述性信息"

	member1 := &model.ChatGroupMember{
		Nickname: longNickname,
	}
	assert.Equal(t, longNickname, member1.Nickname)

	// 测试特殊字符
	member2 := &model.ChatGroupMember{
		Role:     "special@role#123",
		Nickname: "昵称@#$%^&*()_+-=[]{}|\\",
	}
	assert.Equal(t, "special@role#123", member2.Role)
	assert.Equal(t, "昵称@#$%^&*()_+-=[]{}|\\", member2.Nickname)

	// 测试大数值
	member3 := &model.ChatGroupMember{
		GroupID:           ^uint64(0),
		UserID:            ^uint64(0),
		LastReadMessageID: &[]uint64{^uint64(0)}[0],
	}
	assert.Equal(t, ^uint64(0), member3.GroupID)
	assert.Equal(t, ^uint64(0), member3.UserID)
	assert.Equal(t, ^uint64(0), *member3.LastReadMessageID)
}

func TestChatGroupMemberRoles(t *testing.T) {
	// 测试不同角色
	roles := []string{"admin", "moderator", "member", "guest", "owner"}

	for _, role := range roles {
		member := &model.ChatGroupMember{
			Role: role,
		}
		assert.Equal(t, role, member.Role)
	}
}

func TestChatMessageModel(t *testing.T) {
	now := time.Now()
	replyToID := uint64(50)
	moderatedBy := uint64(100)
	moderatedAt := now.Add(5 * time.Minute)

	message := &model.ChatMessage{
		Base: model.Base{
			ID: 1,
		},
		GroupID:       100,
		SenderID:      200,
		Content:       "这是一条测试消息",
		MessageType:   model.ChatMessageTypeText,
		ReplyToID:     &replyToID,
		ImageURL:      "https://example.com/image.jpg",
		Metadata:      `{"fontSize": 14, "color": "#000000"}`,
		IsDeleted:     false,
		AuditStatus:   model.ChatMessageAuditApproved,
		ModeratedBy:   &moderatedBy,
		ModeratedAt:   &moderatedAt,
		RejectReason:  "",
	}

	assert.Equal(t, uint64(1), message.ID)
	assert.Equal(t, uint64(100), message.GroupID)
	assert.Equal(t, uint64(200), message.SenderID)
	assert.Equal(t, "这是一条测试消息", message.Content)
	assert.Equal(t, model.ChatMessageTypeText, message.MessageType)
	assert.Equal(t, &replyToID, message.ReplyToID)
	assert.Equal(t, "https://example.com/image.jpg", message.ImageURL)
	assert.Equal(t, `{"fontSize": 14, "color": "#000000"}`, message.Metadata)
	assert.False(t, message.IsDeleted)
	assert.Equal(t, model.ChatMessageAuditApproved, message.AuditStatus)
	assert.Equal(t, &moderatedBy, message.ModeratedBy)
	assert.Equal(t, &moderatedAt, message.ModeratedAt)
	assert.Equal(t, "", message.RejectReason)
}

func TestChatMessageJSONSerialization(t *testing.T) {
	now := time.Now()

	message := &model.ChatMessage{
		Base: model.Base{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		GroupID:     100,
		SenderID:    200,
		Content:     "这是一条测试消息",
		MessageType: model.ChatMessageTypeText,
		IsDeleted:   false,
		AuditStatus: model.ChatMessageAuditPending,
	}

	// 序列化
	data, err := json.Marshal(message)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "这是一条测试消息")
	assert.Contains(t, string(data), "text")

	// 反序列化
	var decoded model.ChatMessage
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, message.ID, decoded.ID)
	assert.Equal(t, message.GroupID, decoded.GroupID)
	assert.Equal(t, message.SenderID, decoded.SenderID)
	assert.Equal(t, message.Content, decoded.Content)
	assert.Equal(t, message.MessageType, decoded.MessageType)
}

func TestChatMessageTableName(t *testing.T) {
	message := model.ChatMessage{}
	tableName := message.TableName()
	assert.Equal(t, "chat_messages", tableName)
}

func TestChatMessageZeroValues(t *testing.T) {
	message := &model.ChatMessage{
		GroupID:       0,
		SenderID:      0,
		Content:       "",
		MessageType:   "",
		ReplyToID:     nil,
		ImageURL:      "",
		Metadata:      "",
		IsDeleted:     false,
		AuditStatus:   "",
		ModeratedBy:   nil,
		ModeratedAt:   nil,
		RejectReason:  "",
	}

	assert.Equal(t, uint64(0), message.GroupID)
	assert.Equal(t, uint64(0), message.SenderID)
	assert.Equal(t, "", message.Content)
	assert.Equal(t, model.ChatMessageType(""), message.MessageType)
	assert.Nil(t, message.ReplyToID)
	assert.Equal(t, "", message.ImageURL)
	assert.Equal(t, "", message.Metadata)
	assert.False(t, message.IsDeleted)
	assert.Equal(t, model.ChatMessageAuditStatus(""), message.AuditStatus)
	assert.Nil(t, message.ModeratedBy)
	assert.Nil(t, message.ModeratedAt)
	assert.Equal(t, "", message.RejectReason)
}

func TestChatMessageEdgeCases(t *testing.T) {
	// 测试长内容
	longContent := "这是一条非常长的聊天消息，可能包含很多详细信息，用户会在这里发送他们的想法、问题、建议等等。这种长消息测试可以确保我们的模型能够处理各种长度的文本输入，包括很长的消息内容。"

	message1 := &model.ChatMessage{
		Content: longContent,
	}
	assert.Equal(t, longContent, message1.Content)

	// 测试特殊字符
	message2 := &model.ChatMessage{
		Content:      "消息@#$%^&*()_+-=[]{}|\\",
		ImageURL:     "https://example.com/image@special#chars.png",
		Metadata:     `{"special": "@#$%^&*()_+-=[]{}|\\\"quotes\""}`,
		RejectReason: "拒绝原因@#$%^&*()_+-=[]{}|;':\",./<>?",
	}
	assert.Equal(t, "消息@#$%^&*()_+-=[]{}|\\", message2.Content)
	assert.Equal(t, "https://example.com/image@special#chars.png", message2.ImageURL)
	assert.Contains(t, message2.Metadata, "special")
	assert.Equal(t, "拒绝原因@#$%^&*()_+-=[]{}|;':\",./<>?", message2.RejectReason)

	// 测试大数值
	message3 := &model.ChatMessage{
		GroupID:       ^uint64(0),
		SenderID:      ^uint64(0),
		ReplyToID:     &[]uint64{^uint64(0)}[0],
		ModeratedBy:   &[]uint64{^uint64(0)}[0],
	}
	assert.Equal(t, ^uint64(0), message3.GroupID)
	assert.Equal(t, ^uint64(0), message3.SenderID)
	assert.Equal(t, ^uint64(0), *message3.ReplyToID)
	assert.Equal(t, ^uint64(0), *message3.ModeratedBy)
}

func TestChatMessageTypes(t *testing.T) {
	// 测试所有消息类型
	messageTypes := []model.ChatMessageType{
		model.ChatMessageTypeText,
		model.ChatMessageTypeImage,
		model.ChatMessageTypeFile,
		model.ChatMessageTypeSystem,
	}

	for _, msgType := range messageTypes {
		message := &model.ChatMessage{
			MessageType: msgType,
		}
		assert.Equal(t, msgType, message.MessageType)
	}
}

func TestChatMessageAuditStatuses(t *testing.T) {
	// 测试所有审核状态
	auditStatuses := []model.ChatMessageAuditStatus{
		model.ChatMessageAuditPending,
		model.ChatMessageAuditApproved,
		model.ChatMessageAuditRejected,
	}

	for _, status := range auditStatuses {
		message := &model.ChatMessage{
			AuditStatus: status,
		}
		assert.Equal(t, status, message.AuditStatus)
	}
}

func TestChatMessageWithReply(t *testing.T) {
	replyToID := uint64(75)

	message := &model.ChatMessage{
		ReplyToID: &replyToID,
	}

	assert.NotNil(t, message.ReplyToID)
	assert.Equal(t, uint64(75), *message.ReplyToID)
}

func TestChatMessageWithoutReply(t *testing.T) {
	message := &model.ChatMessage{
		ReplyToID: nil,
	}

	assert.Nil(t, message.ReplyToID)
}

func TestChatMessageRelations(t *testing.T) {
	message := &model.ChatMessage{
		Group: model.ChatGroup{
			GroupName: "测试群组",
		},
	}

	assert.Equal(t, "测试群组", message.Group.GroupName)
}

func TestChatGroupMemberRelations(t *testing.T) {
	member := &model.ChatGroupMember{
		Group: model.ChatGroup{
			GroupName: "测试群组",
		},
	}

	assert.Equal(t, "测试群组", member.Group.GroupName)
}