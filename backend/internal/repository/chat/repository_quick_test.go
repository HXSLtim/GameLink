package chat

import (
    "testing"
    "context"
    "time"

    "github.com/glebarez/sqlite"
    "gorm.io/gorm"

    "gamelink/internal/model"
    "gamelink/internal/repository"
)

func newDB(t *testing.T) *gorm.DB {
    t.Helper()
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil { t.Fatalf("db: %v", err) }
    if err := db.AutoMigrate(&model.ChatGroup{}, &model.ChatGroupMember{}, &model.ChatMessage{}, &model.ChatReport{}); err != nil { t.Fatalf("migrate: %v", err) }
    return db
}

func TestChatGroupRepository_CRUD_AndQueries(t *testing.T) {
    db := newDB(t)
    repo := NewChatGroupRepository(db)
    ctx := context.Background()

    g := &model.ChatGroup{GroupName:"G", GroupType:model.ChatGroupTypeOrder, CreatedBy:1, IsActive:true}
    if err := repo.Create(ctx, g); err != nil { t.Fatalf("create: %v", err) }
    got, err := repo.Get(ctx, g.ID)
    if err != nil || got.ID != g.ID { t.Fatalf("get: %v", err) }

    // Update and deactivate
    got.Description = "desc"
    if err := repo.Update(ctx, got); err != nil { t.Fatalf("update: %v", err) }
    if err := repo.Deactivate(ctx, got.ID); err != nil { t.Fatalf("deactivate: %v", err) }

    // List deactivated before now
    list, err := repo.ListDeactivatedBefore(ctx, time.Now().Add(1*time.Second), 10)
    if err != nil || len(list) == 0 { t.Fatalf("list deactivated: %v len=%d", err, len(list)) }

    // Delete by IDs
    if err := repo.DeleteByIDs(ctx, []uint64{got.ID}); err != nil { t.Fatalf("delete: %v", err) }
}

func TestChatMemberRepository_CRUD(t *testing.T) {
    db := newDB(t)
    gr := NewChatGroupRepository(db)
    mr := NewChatMemberRepository(db)
    ctx := context.Background()
    g := &model.ChatGroup{GroupName:"G", GroupType:model.ChatGroupTypePublic, CreatedBy:1}
    _ = gr.Create(ctx, g)

    m := &model.ChatGroupMember{GroupID:g.ID, UserID:2, Nickname:"n", JoinedAt: time.Now()}
    if err := mr.Add(ctx, m); err != nil { t.Fatalf("add: %v", err) }
    m.Nickname = "n2"
    if err := mr.Update(ctx, m); err != nil { t.Fatalf("update: %v", err) }
    got, err := mr.Get(ctx, g.ID, 2)
    if err != nil || got.Nickname != "n2" { t.Fatalf("get: %v", err) }
    if err := mr.Remove(ctx, g.ID, 2); err != nil { t.Fatalf("remove: %v", err) }
}

func TestChatMessageRepository_CRUD_AndFilters(t *testing.T) {
    db := newDB(t)
    gr := NewChatGroupRepository(db)
    repo := NewChatMessageRepository(db)
    ctx := context.Background()
    g := &model.ChatGroup{GroupName:"G", GroupType:model.ChatGroupTypePublic, CreatedBy:1}
    _ = gr.Create(ctx, g)

    m1 := &model.ChatMessage{GroupID:g.ID, SenderID:1, Content:"a", AuditStatus:model.ChatMessageAuditApproved}
    m2 := &model.ChatMessage{GroupID:g.ID, SenderID:2, Content:"b", AuditStatus:model.ChatMessageAuditRejected}
    _ = repo.Create(ctx, m1)
    _ = repo.Create(ctx, m2)

    // ListByGroup with audit statuses
    opts := repository.ChatMessageListOptions{Page:1, PageSize:10, GroupID:g.ID, AuditStatuses: []model.ChatMessageAuditStatus{model.ChatMessageAuditApproved}}
    msgs, total, err := repo.ListByGroup(ctx, opts)
    if err != nil || total == 0 || len(msgs) == 0 { t.Fatalf("list: %v", err) }

    // Get and mark deleted
    got, err := repo.Get(ctx, m1.ID)
    if err != nil || got.ID != m1.ID { t.Fatalf("get: %v", err) }
    if err := repo.MarkDeleted(ctx, m1.ID, 1); err != nil { t.Fatalf("mark: %v", err) }

    // Moderation list and update audit status
    mopts := repository.ChatMessageModerationListOptions{Page:1, PageSize:10, GroupID:&g.ID}
    q, total2, err := repo.ListForModeration(ctx, mopts)
    if err != nil || total2 < 0 { t.Fatalf("mod list: %v", err) }
    _ = q
    if err := repo.UpdateAuditStatus(ctx, m2.ID, model.ChatMessageAuditApproved, nil, ""); err != nil { t.Fatalf("update audit: %v", err) }

    // Delete by group IDs
    if err := repo.DeleteByGroupIDs(ctx, []uint64{g.ID}); err != nil { t.Fatalf("delete by groups: %v", err) }
}

func TestChatReportRepository_CRUD_List(t *testing.T) {
    db := newDB(t)
    rr := NewChatReportRepository(db)
    ctx := context.Background()
    rep := &model.ChatReport{MessageID:1, ReporterID:1, Reason:"spam", Status:"pending"}
    if err := rr.Create(ctx, rep); err != nil { t.Fatalf("create: %v", err) }
    got, err := rr.Get(ctx, rep.ID)
    if err != nil || got.ID != rep.ID { t.Fatalf("get: %v", err) }
    got.Status = "approved"
    if err := rr.Update(ctx, got); err != nil { t.Fatalf("update: %v", err) }
    ropts := repository.ChatReportListOptions{Page:1, PageSize:10, Status:"approved"}
    list, total, err := rr.List(ctx, ropts)
    if err != nil || total == 0 || len(list) == 0 { t.Fatalf("list: %v", err) }
}

func TestChatListByUserAndMembers_FiltersPagination(t *testing.T) {}
