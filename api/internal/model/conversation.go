package model

import "time"

// ConversationStatus represents customer service conversation status.
type ConversationStatus string

const (
	ConversationStatusActive ConversationStatus = "active"
	ConversationStatusClosed ConversationStatus = "closed"
)

// Conversation describes a customer-service conversation summary.
// It is used by customer-service REST responses.
type Conversation struct {
	ID            uint64             `json:"id"`
	GroupID       uint64             `json:"groupId"`
	UserID        uint64             `json:"userId"`
	AgentID       uint64             `json:"agentId"`
	AgentName     string             `json:"agentName"`
	AgentAvatar   string             `json:"agentAvatar,omitempty"`
	Status        ConversationStatus `json:"status"`
	IsAgentOnline bool               `json:"isAgentOnline"`
	LastMessage   string             `json:"lastMessage,omitempty"`
	LastMessageAt *time.Time         `json:"lastMessageAt,omitempty"`
	CreatedAt     time.Time          `json:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt"`
}

// ConversationMessage represents a chat message in customer-service APIs.
type ConversationMessage struct {
	ID          uint64          `json:"id"`
	GroupID     uint64          `json:"groupId"`
	SenderID    uint64          `json:"senderId"`
	Content     string          `json:"content"`
	MessageType ChatMessageType `json:"messageType"`
	IsMe        bool            `json:"isMe"`
	CreatedAt   string          `json:"createdAt"`
}

// ConversationListResponse is returned by GET /user/customer-service/conversations.
type ConversationListResponse struct {
	Conversations []Conversation        `json:"conversations"`
	Messages      []ConversationMessage `json:"messages"`
	Total         int64                 `json:"total"`
	Page          int                   `json:"page"`
	PageSize      int                   `json:"pageSize"`
	HasMore       bool                  `json:"hasMore"`
}

// ConversationCreateRequest is used by POST /user/customer-service/conversations.
type ConversationCreateRequest struct {
	Content string `json:"content" binding:"required"`
}
