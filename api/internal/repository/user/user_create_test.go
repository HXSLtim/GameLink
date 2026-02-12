package user

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"gamelink/internal/model"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}

	// User has optional relation to vip_levels; migrate both to avoid fk issues.
	if err := db.AutoMigrate(&model.VipLevel{}, &model.User{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}
	return db
}

func TestUserRepository_Create_EmailOnly_NoUniqueConflictOnEmptyPhone(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	u1 := &model.User{
		Email:        "email.only.1@gamelink.com",
		Name:         "EmailOnly1",
		Nickname:     "EmailOnly1",
		PasswordHash: "hash1",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
	}
	if err := repo.Create(ctx, u1); err != nil {
		t.Fatalf("create user1 failed: %v", err)
	}

	u2 := &model.User{
		Email:        "email.only.2@gamelink.com",
		Name:         "EmailOnly2",
		Nickname:     "EmailOnly2",
		PasswordHash: "hash2",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
	}
	if err := repo.Create(ctx, u2); err != nil {
		t.Fatalf("create user2 failed: %v", err)
	}
}

func TestUserRepository_Create_PhoneOnly_NoUniqueConflictOnEmptyEmail(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	u1 := &model.User{
		Phone:        "13900000001",
		Name:         "PhoneOnly1",
		Nickname:     "PhoneOnly1",
		PasswordHash: "hash1",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
	}
	if err := repo.Create(ctx, u1); err != nil {
		t.Fatalf("create user1 failed: %v", err)
	}

	u2 := &model.User{
		Phone:        "13900000002",
		Name:         "PhoneOnly2",
		Nickname:     "PhoneOnly2",
		PasswordHash: "hash2",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
	}
	if err := repo.Create(ctx, u2); err != nil {
		t.Fatalf("create user2 failed: %v", err)
	}
}

