//go:build integration
package chat

import (
    "context"
    "testing"
    "time"

    "github.com/glebarez/sqlite"
    "gorm.io/gorm"

    "gamelink/internal/model"
    "gamelink/internal/repository"
)

func setupMessageRepoTestDB(t *testing.T) *gorm.DB {
    t.Helper()

    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("open db: %v", err)
    }

    if err := db.AutoMigrate(&model.ChatMessage{}); err != nil {
        t.Fatalf("migrate: %v", err)
    }

    return db
}

func TestChatMessageRepository_CreateAndGet(t *testing.T) {
    db := setupMessageRepoTestDB(t)
    repo := NewChatMessageRepository(db)

    msg := &model.ChatMessage{
        GroupID:     1,
        SenderID:    2,
        Content:     "hello",
        MessageType: model.ChatMessageTypeText,
        AuditStatus: model.ChatMessageAuditPending,
    }
    if err := repo.Create(context.Background(), msg); err != nil {
        t.Fatalf("create message: %v", err)
    }
    if msg.ID == 0 {
        t.Fatal("expected ID to be assigned")
    }

    got, err := repo.Get(context.Background(), msg.ID)
    if err != nil {
        t.Fatalf("get message: %v", err)
    }
    if got.Content != msg.Content {
        t.Fatalf("expected content %q, got %q", msg.Content, got.Content)
    }
}

func TestChatMessageRepository_CreateBatch(t *testing.T) {
    db := setupMessageRepoTestDB(t)
    repo := NewChatMessageRepository(db)

    if err := repo.CreateBatch(context.Background(), nil); err != nil {
        t.Fatalf("create batch empty: %v", err)
    }

    messages := []*model.ChatMessage{
        {
            GroupID:     10,
            SenderID:    1,
            Content:     "first",
            MessageType: model.ChatMessageTypeText,
            AuditStatus: model.ChatMessageAuditApproved,
        },
        {
            GroupID:     10,
            SenderID:    2,
            Content:     "second",
            MessageType: model.ChatMessageTypeImage,
            AuditStatus: model.ChatMessageAuditRejected,
        },
    }

    if err := repo.CreateBatch(context.Background(), messages); err != nil {
        t.Fatalf("create batch: %v", err)
    }

    opts := repository.ChatMessageListOptions{
        GroupID:       10,
        Page:          1,
        PageSize:      10,
        MessageType:   func() *model.ChatMessageType { v := model.ChatMessageTypeText; return &v }(),
        AuditStatuses: []model.ChatMessageAuditStatus{model.ChatMessageAuditApproved, model.ChatMessageAuditRejected},
        BeforeID:      &messages[1].ID,
        AfterID:       func() *uint64 { v := messages[0].ID - 1; return &v }(),
    }
    now := time.Now()
    opts.DateFrom = &now
    opts.DateTo = &now
    // update timestamps to now to satisfy date filters
    if err := db.Model(&model.ChatMessage{}).Where("group_id = ?", 10).Update("created_at", now).Error; err != nil {
        t.Fatalf("update timestamps: %v", err)
    }

    list, total, err := repo.ListByGroup(context.Background(), opts)
    if err != nil {
        t.Fatalf("list by group: %v", err)
    }
    if total != 2 {
        t.Fatalf("expected total 2, got %d", total)
    }
    if len(list) != 1 || list[0].Content != "first" {
        t.Fatalf("expected filtered message 'first', got %#v", list)
    }
}

func TestChatMessageRepository_MarkDeleted(t *testing.T) {
    db := setupMessageRepoTestDB(t)
    repo := NewChatMessageRepository(db)

    msg := &model.ChatMessage{GroupID: 1, SenderID: 1, Content: "to delete"}
    if err := repo.Create(context.Background(), msg); err != nil {
        t.Fatalf("create message: %v", err)
    }

    if err := repo.MarkDeleted(context.Background(), msg.ID, 999); err != nil {
        t.Fatalf("mark deleted: %v", err)
    }

    var stored model.ChatMessage
    if err := db.First(&stored, msg.ID).Error; err != nil {
        t.Fatalf("load message: %v", err)
    }
    if !stored.IsDeleted {
        t.Fatal("expected message to be marked deleted")
    }
}

func TestChatMessageRepository_ListForModeration(t *testing.T) {
    db := setupMessageRepoTestDB(t)
    repo := NewChatMessageRepository(db)

    msgs := []*model.ChatMessage{
        {GroupID: 1, SenderID: 1, Content: "pending", AuditStatus: model.ChatMessageAuditPending},
        {GroupID: 1, SenderID: 2, Content: "approved", AuditStatus: model.ChatMessageAuditApproved},
        {GroupID: 2, SenderID: 1, Content: "rejected", AuditStatus: model.ChatMessageAuditRejected},
    }
    if err := repo.CreateBatch(context.Background(), msgs); err != nil {
        t.Fatalf("seed messages: %v", err)
    }

    opts := repository.ChatMessageModerationListOptions{
        Page:        1,
        PageSize:    10,
        GroupID:     func() *uint64 { v := uint64(1); return &v }(),
        SenderID:    func() *uint64 { v := uint64(1); return &v }(),
        AuditStatus: func() *model.ChatMessageAuditStatus { v := model.ChatMessageAuditPending; return &v }(),
    }

    list, total, err := repo.ListForModeration(context.Background(), opts)
    if err != nil {
        t.Fatalf("list moderation: %v", err)
    }
    if total != 1 || len(list) != 1 {
        t.Fatalf("expected 1 pending message, got total=%d len=%d", total, len(list))
    }
    if list[0].Content != "pending" {
        t.Fatalf("expected 'pending', got %q", list[0].Content)
    }
}

func TestChatMessageRepository_UpdateAuditStatus(t *testing.T) {
    db := setupMessageRepoTestDB(t)
    repo := NewChatMessageRepository(db)

    msg := &model.ChatMessage{GroupID: 1, SenderID: 1, Content: "needs review", AuditStatus: model.ChatMessageAuditPending}
    if err := repo.Create(context.Background(), msg); err != nil {
        t.Fatalf("create message: %v", err)
    }

    if err := repo.UpdateAuditStatus(context.Background(), msg.ID, model.ChatMessageAuditRejected, nil, "spam"); err != nil {
        t.Fatalf("update audit status nil moderator: %v", err)
    }
    var updated model.ChatMessage
    if err := db.First(&updated, msg.ID).Error; err != nil {
        t.Fatalf("load message: %v", err)
    }
    if updated.AuditStatus != model.ChatMessageAuditRejected || updated.RejectReason != "spam" {
        t.Fatalf("unexpected audit status after rejection: %+v", updated)
    }

    moderator := uint64(42)
    if err := repo.UpdateAuditStatus(context.Background(), msg.ID, model.ChatMessageAuditApproved, &moderator, ""); err != nil {
        t.Fatalf("update audit status moderator: %v", err)
    }
    if err := db.First(&updated, msg.ID).Error; err != nil {
        t.Fatalf("reload message: %v", err)
    }
    if updated.AuditStatus != model.ChatMessageAuditApproved {
        t.Fatalf("expected approved status, got %s", updated.AuditStatus)
    }
    if updated.ModeratedBy == nil || *updated.ModeratedBy != moderator {
        t.Fatalf("expected moderatedBy=%d, got %+v", moderator, updated.ModeratedBy)
    }
    if updated.ModeratedAt == nil {
        t.Fatal("expected moderatedAt to be set")
    }
}
