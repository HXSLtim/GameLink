package chat

import (
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"gamelink/internal/model"
)

func newMemDB(t *testing.T) *gorm.DB {
	 t.Helper()
	 db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	 if err != nil { t.Fatalf("open db: %v", err) }
	 if err := db.AutoMigrate(&model.ChatGroup{}); err != nil { t.Fatalf("migrate: %v", err) }
	 return db
}

func TestChatGroupRepository_Deactivate_SetsTimestamp(t *testing.T) {
	 db := newMemDB(t)
	 repo := NewChatGroupRepository(db)
	 g := &model.ChatGroup{ GroupName: "order-1", GroupType: model.ChatGroupTypeOrder, IsActive: true }
	 if err := db.Create(g).Error; err != nil { t.Fatalf("seed: %v", err) }

	 if err := repo.Deactivate(context.Background(), g.ID); err != nil { t.Fatalf("deactivate: %v", err) }
	 var got model.ChatGroup
	 if err := db.First(&got, g.ID).Error; err != nil { t.Fatalf("load: %v", err) }
	 if got.IsActive { t.Fatal("expected inactive") }
	 if got.DeactivatedAt == nil { t.Fatal("expected deactivated_at set") }
}

func TestChatGroupRepository_ListDeactivatedBefore_OrderOnly(t *testing.T) {
	 db := newMemDB(t)
	 repo := NewChatGroupRepository(db)
	 now := time.Now()
	 older := now.AddDate(0,0,-40)
	 newer := now.AddDate(0,0,-20)

	 // order older than 30d
	 g1 := &model.ChatGroup{ GroupName: "o-older", GroupType: model.ChatGroupTypeOrder, IsActive: false, DeactivatedAt: &older }
	 // order newer than 30d
	 g2 := &model.ChatGroup{ GroupName: "o-newer", GroupType: model.ChatGroupTypeOrder, IsActive: false, DeactivatedAt: &newer }
	 // public older than 30d (should be ignored)
	 g3 := &model.ChatGroup{ GroupName: "p-older", GroupType: model.ChatGroupTypePublic, IsActive: false, DeactivatedAt: &older }
	 for _, g := range []*model.ChatGroup{g1,g2,g3} {
		 if err := db.Create(g).Error; err != nil { t.Fatalf("seed: %v", err) }
	 }

	 cutoff := now.AddDate(0,0,-30)
	 list, err := repo.ListDeactivatedBefore(context.Background(), cutoff, 10)
	 if err != nil { t.Fatalf("list: %v", err) }
	 if len(list) != 1 || list[0].GroupName != "o-older" { t.Fatalf("unexpected list: %+v", list) }
}
