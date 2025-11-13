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

func newDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("db: %v", err)
	}
	if err := db.AutoMigrate(&model.ChatGroup{}, &model.ChatGroupMember{}, &model.ChatMessage{}, &model.ChatReport{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func TestChatGroupRepository_CRUD_AndQueries(t *testing.T) {
	db := newDB(t)
	repo := NewChatGroupRepository(db)
	ctx := context.Background()

	g := &model.ChatGroup{GroupName: "G", GroupType: model.ChatGroupTypeOrder, CreatedBy: 1, IsActive: true}
	if err := repo.Create(ctx, g); err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := repo.Get(ctx, g.ID)
	if err != nil || got.ID != g.ID {
		t.Fatalf("get: %v", err)
	}

	// Update and deactivate
	got.Description = "desc"
	if err := repo.Update(ctx, got); err != nil {
		t.Fatalf("update: %v", err)
	}
	if err := repo.Deactivate(ctx, got.ID); err != nil {
		t.Fatalf("deactivate: %v", err)
	}

	// List deactivated before now
	list, err := repo.ListDeactivatedBefore(ctx, time.Now().Add(1*time.Second), 10)
	if err != nil || len(list) == 0 {
		t.Fatalf("list deactivated: %v len=%d", err, len(list))
	}

	// Delete by IDs
	if err := repo.DeleteByIDs(ctx, []uint64{got.ID}); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestChatMemberRepository_CRUD(t *testing.T) {
	db := newDB(t)
	gr := NewChatGroupRepository(db)
	mr := NewChatMemberRepository(db)
	ctx := context.Background()
	g := &model.ChatGroup{GroupName: "G", GroupType: model.ChatGroupTypePublic, CreatedBy: 1}
	_ = gr.Create(ctx, g)

	m := &model.ChatGroupMember{GroupID: g.ID, UserID: 2, Nickname: "n", JoinedAt: time.Now()}
	if err := mr.Add(ctx, m); err != nil {
		t.Fatalf("add: %v", err)
	}
	m.Nickname = "n2"
	if err := mr.Update(ctx, m); err != nil {
		t.Fatalf("update: %v", err)
	}
	got, err := mr.Get(ctx, g.ID, 2)
	if err != nil || got.Nickname != "n2" {
		t.Fatalf("get: %v", err)
	}
	if err := mr.Remove(ctx, g.ID, 2); err != nil {
		t.Fatalf("remove: %v", err)
	}
}

func TestChatMessageRepository_CRUD_AndFilters(t *testing.T) {
	db := newDB(t)
	gr := NewChatGroupRepository(db)
	repo := NewChatMessageRepository(db)
	ctx := context.Background()
	g := &model.ChatGroup{GroupName: "G", GroupType: model.ChatGroupTypePublic, CreatedBy: 1}
	_ = gr.Create(ctx, g)

	m1 := &model.ChatMessage{GroupID: g.ID, SenderID: 1, Content: "a", AuditStatus: model.ChatMessageAuditApproved}
	m2 := &model.ChatMessage{GroupID: g.ID, SenderID: 2, Content: "b", AuditStatus: model.ChatMessageAuditRejected}
	_ = repo.Create(ctx, m1)
	_ = repo.Create(ctx, m2)

	// ListByGroup with audit statuses
	opts := repository.ChatMessageListOptions{Page: 1, PageSize: 10, GroupID: g.ID, AuditStatuses: []model.ChatMessageAuditStatus{model.ChatMessageAuditApproved}}
	msgs, total, err := repo.ListByGroup(ctx, opts)
	if err != nil || total == 0 || len(msgs) == 0 {
		t.Fatalf("list: %v", err)
	}

	// Get and mark deleted
	got, err := repo.Get(ctx, m1.ID)
	if err != nil || got.ID != m1.ID {
		t.Fatalf("get: %v", err)
	}
	if err := repo.MarkDeleted(ctx, m1.ID, 1); err != nil {
		t.Fatalf("mark: %v", err)
	}

	// Moderation list and update audit status
	mopts := repository.ChatMessageModerationListOptions{Page: 1, PageSize: 10, GroupID: &g.ID}
	q, total2, err := repo.ListForModeration(ctx, mopts)
	if err != nil || total2 < 0 {
		t.Fatalf("mod list: %v", err)
	}
	_ = q
	if err := repo.UpdateAuditStatus(ctx, m2.ID, model.ChatMessageAuditApproved, nil, ""); err != nil {
		t.Fatalf("update audit: %v", err)
	}

	// Delete by group IDs
	if err := repo.DeleteByGroupIDs(ctx, []uint64{g.ID}); err != nil {
		t.Fatalf("delete by groups: %v", err)
	}
}

func TestChatReportRepository_CRUD_List(t *testing.T) {
	db := newDB(t)
	rr := NewChatReportRepository(db)
	ctx := context.Background()
	rep := &model.ChatReport{MessageID: 1, ReporterID: 1, Reason: "spam", Status: "pending"}
	if err := rr.Create(ctx, rep); err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := rr.Get(ctx, rep.ID)
	if err != nil || got.ID != rep.ID {
		t.Fatalf("get: %v", err)
	}
	got.Status = "approved"
	if err := rr.Update(ctx, got); err != nil {
		t.Fatalf("update: %v", err)
	}
	ropts := repository.ChatReportListOptions{Page: 1, PageSize: 10, Status: "approved"}
	list, total, err := rr.List(ctx, ropts)
	if err != nil || total == 0 || len(list) == 0 {
		t.Fatalf("list: %v", err)
	}
}

func TestChatListByUserAndMembers_FiltersPagination(t *testing.T) {
	db := newDB(t)
	groupRepo := NewChatGroupRepository(db)
	memberRepo := NewChatMemberRepository(db)
	ctx := context.Background()

	orderType := model.ChatGroupTypeOrder
	related1 := uint64(1001)
	related2 := uint64(1002)

	// 准备三个群组，其中两个属于同一用户，包含不同的活跃状态和关键字，便于测试过滤逻辑
	g1 := &model.ChatGroup{GroupName: "order-foo", Description: "foo desc", GroupType: orderType, CreatedBy: 1, IsActive: true, RelatedOrderID: &related1}
	g2 := &model.ChatGroup{GroupName: "order-bar", Description: "bar desc", GroupType: orderType, CreatedBy: 1, IsActive: false, RelatedOrderID: &related2}
	g3 := &model.ChatGroup{GroupName: "other", Description: "other", GroupType: model.ChatGroupTypePublic, CreatedBy: 2, IsActive: true}
	if err := groupRepo.Create(ctx, g1); err != nil {
		t.Fatalf("create g1: %v", err)
	}
	if err := groupRepo.Create(ctx, g2); err != nil {
		t.Fatalf("create g2: %v", err)
	}
	if err := groupRepo.Create(ctx, g3); err != nil {
		t.Fatalf("create g3: %v", err)
	}
	// 由于 GORM 对零值字段会应用 default，因此需要显式更新数据库中的 is_active=false
	if err := db.Model(&model.ChatGroup{}).Where("id = ?", g2.ID).
		Updates(map[string]any{"is_active": false}).Error; err != nil {
		t.Fatalf("deactivate g2: %v", err)
	}

	// user 10 属于 g1、g2，user 20 只在 g3 中
	members := []*model.ChatGroupMember{
		{GroupID: g1.ID, UserID: 10, Nickname: "alpha", Role: "member"},
		{GroupID: g2.ID, UserID: 10, Nickname: "beta", Role: "admin"},
		{GroupID: g3.ID, UserID: 20, Nickname: "gamma", Role: "member"},
	}
	for _, m := range members {
		if err := memberRepo.Add(ctx, m); err != nil {
			t.Fatalf("add member: %v", err)
		}
	}

	// 1) 基础：只返回 user 10 且 active 的群组
	opts := repository.ChatGroupListOptions{Page: 1, PageSize: 10}
	list, total, err := groupRepo.ListByUser(ctx, 10, opts)
	if err != nil {
		t.Fatalf("ListByUser basic: %v", err)
	}
	if total != 1 || len(list) != 1 || list[0].ID != g1.ID {
		t.Fatalf("unexpected basic list, total=%d len=%d first=%+v", total, len(list), list)
	}

	// 2) IncludeInactive=true 时应包含 g2
	opts.IncludeInactive = true
	list, total, err = groupRepo.ListByUser(ctx, 10, opts)
	if err != nil {
		t.Fatalf("ListByUser include inactive: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected 2 groups for user 10 when including inactive, got %d", total)
	}

	// 3) 关键字和 GroupType 过滤，只匹配包含 foo 的订单群
	opts.Keyword = "foo"
	opts.GroupType = &orderType
	list, total, err = groupRepo.ListByUser(ctx, 10, opts)
	if err != nil {
		t.Fatalf("ListByUser keyword+type: %v", err)
	}
	if total != 1 || len(list) != 1 || list[0].ID != g1.ID {
		t.Fatalf("expected only g1 for keyword filter, got total=%d len=%d first=%+v", total, len(list), list)
	}

	// 4) RelatedOrderID 过滤，应只返回 g2
	opts.Keyword = ""
	opts.IncludeInactive = true
	opts.RelatedOrderID = &related2
	list, total, err = groupRepo.ListByUser(ctx, 10, opts)
	if err != nil {
		t.Fatalf("ListByUser related order: %v", err)
	}
	if total != 1 || len(list) != 1 || list[0].ID != g2.ID {
		t.Fatalf("expected only g2 for related order filter, got total=%d len=%d first=%+v", total, len(list), list)
	}

	// 5) ListMembers：默认分页（page/pageSize<=0）应返回所有成员
	mOpts := repository.ChatGroupMemberListOptions{}
	memList, memTotal, err := groupRepo.ListMembers(ctx, g1.ID, mOpts)
	if err != nil {
		t.Fatalf("ListMembers basic: %v", err)
	}
	if memTotal != 1 || len(memList) != 1 {
		t.Fatalf("expected 1 member in g1, got total=%d len=%d", memTotal, len(memList))
	}

	// 6) ListMembers：按角色和关键字过滤
	mOpts = repository.ChatGroupMemberListOptions{
		Page:     1,
		PageSize: 1,
		Role:     "admin",
		Keyword:  "beta",
	}
	memList, memTotal, err = groupRepo.ListMembers(ctx, g2.ID, mOpts)
	if err != nil {
		t.Fatalf("ListMembers filtered: %v", err)
	}
	if memTotal != 1 || len(memList) != 1 || memList[0].UserID != 10 {
		t.Fatalf("expected single admin member userID=10, got total=%d len=%d first=%+v", memTotal, len(memList), memList)
	}
}
