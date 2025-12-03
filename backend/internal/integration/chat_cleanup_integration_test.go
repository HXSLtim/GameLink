package integration

import (
	"testing"
	"time"

	"gamelink/internal/model"
	chatrepo "gamelink/internal/repository/chat"
	chatservice "gamelink/internal/service/chat"
	"gamelink/pkg/testutil"
)

// 聊天组清理：删除已停用且超过截止时间的订单群
func TestChatGroupCleanupInactiveOrderGroups(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateChatModels(t, db)

	repo := chatrepo.NewChatGroupRepository(db)
	svc := chatservice.NewCleanupService(repo)

	// seed groups
	old := time.Now().Add(-48 * time.Hour)
	active := &model.ChatGroup{
		GroupName:      "active order",
		GroupType:      model.ChatGroupTypeOrder,
		IsActive:       true,
		RelatedOrderID: u64ptr(1),
	}
	inactiveOld := &model.ChatGroup{
		GroupName:      "old inactive",
		GroupType:      model.ChatGroupTypeOrder,
		IsActive:       false,
		RelatedOrderID: u64ptr(2),
		DeactivatedAt:  &old,
	}
	inactiveRecent := &model.ChatGroup{
		GroupName:      "recent inactive",
		GroupType:      model.ChatGroupTypeOrder,
		IsActive:       false,
		RelatedOrderID: u64ptr(3),
		DeactivatedAt:  timePtr(time.Now().Add(-2 * time.Hour)),
	}
	_ = repo.Create(ctx(), active)
	_ = repo.Create(ctx(), inactiveOld)
	_ = repo.Create(ctx(), inactiveRecent)

	// 强制更新停用时间，确保满足查询条件
	if err := db.Model(&model.ChatGroup{}).Where("group_name = ?", "old inactive").
		Updates(map[string]any{"is_active": false, "deactivated_at": old}).Error; err != nil {
		t.Fatalf("force update old inactive: %v", err)
	}

	var all []model.ChatGroup
	if err := db.Find(&all).Error; err != nil {
		t.Fatalf("find all groups: %v", err)
	}
	foundOld := false
	for _, g := range all {
		if g.GroupName == "old inactive" {
			foundOld = true
			if g.DeactivatedAt == nil {
				t.Fatalf("expected deactivated_at set for old inactive group")
			}
		}
	}
	if !foundOld {
		t.Fatalf("seed old inactive group not found")
	}

	// cleanup with cutoff 24h ago should delete only old inactive
	eligible, err := repo.ListDeactivatedBefore(ctx(), time.Now().Add(-24*time.Hour), 100)
	if err != nil {
		t.Fatalf("list deactivated before cleanup: %v", err)
	}
	if len(eligible) == 0 {
		t.Fatalf("expected eligible groups >0, got 0")
	}
	deleted, err := svc.CleanupInactiveOrderGroups(ctx(), time.Now().Add(-24*time.Hour), 100)
	if err != nil {
		t.Fatalf("cleanup: %v", err)
	}
	if deleted != 1 {
		t.Fatalf("expected 1 deleted, got %d", deleted)
	}

	// verify remaining groups count = 2
	var remaining int64
	if err := db.Model(&model.ChatGroup{}).Count(&remaining).Error; err != nil {
		t.Fatalf("count remaining: %v", err)
	}
	if remaining != 2 {
		t.Fatalf("expected 2 groups remaining, got %d", remaining)
	}
}
