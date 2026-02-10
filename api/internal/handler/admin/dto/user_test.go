package dto

import (
	"testing"
	"time"

	"gamelink/internal/model"
)

func TestToUserResponse(t *testing.T) {
	now := time.Now()

	t.Run("nil user returns nil", func(t *testing.T) {
		resp := ToUserResponse(nil)
		if resp != nil {
			t.Error("expected nil response for nil user")
		}
	})

	t.Run("basic fields mapped correctly", func(t *testing.T) {
		user := &model.User{
			Phone:     "13800138000",
			Email:     "test@example.com",
			Name:      "Test User",
			Nickname:  "tester",
			AvatarURL: "https://example.com/avatar.png",
			Role:      model.RoleUser,
			Status:    model.UserStatusActive,
		}
		user.ID = 42
		user.CreatedAt = now
		user.UpdatedAt = now

		resp := ToUserResponse(user)

		if resp.ID != 42 {
			t.Errorf("expected ID 42, got %d", resp.ID)
		}
		if resp.Phone != "13800138000" {
			t.Errorf("expected phone 13800138000, got %s", resp.Phone)
		}
		if resp.Email != "test@example.com" {
			t.Errorf("expected email test@example.com, got %s", resp.Email)
		}
		if resp.Name != "Test User" {
			t.Errorf("expected name Test User, got %s", resp.Name)
		}
		if resp.Role != model.RoleUser {
			t.Errorf("expected role user, got %s", resp.Role)
		}
		if resp.Status != model.UserStatusActive {
			t.Errorf("expected status active, got %s", resp.Status)
		}
	})

	t.Run("vip level mapped when present", func(t *testing.T) {
		user := &model.User{
			Name:        "VIP User",
			Role:        model.RoleUser,
			Status:      model.UserStatusActive,
			VipUnlocked: true,
		}
		user.ID = 1
		user.VipLevel = &model.VipLevel{
			Slug:      "vip1",
			Title:     "VIP 1",
			SortOrder: 1,
		}
		user.VipLevel.ID = 10

		resp := ToUserResponse(user)

		if resp.VipLevel == nil {
			t.Fatal("expected VipLevel to be set")
		}
		if resp.VipLevel.Slug != "vip1" {
			t.Errorf("expected slug vip1, got %s", resp.VipLevel.Slug)
		}
		if resp.VipLevel.Title != "VIP 1" {
			t.Errorf("expected title VIP 1, got %s", resp.VipLevel.Title)
		}
		if !resp.VipUnlocked {
			t.Error("expected VipUnlocked to be true")
		}
	})

	t.Run("vip level nil when absent", func(t *testing.T) {
		user := &model.User{Name: "Normal", Role: model.RoleUser, Status: model.UserStatusActive}
		user.ID = 2

		resp := ToUserResponse(user)

		if resp.VipLevel != nil {
			t.Error("expected VipLevel to be nil")
		}
	})
}

func TestToUserResponseList(t *testing.T) {
	users := []model.User{
		{Phone: "111", Name: "A", Role: model.RoleUser, Status: model.UserStatusActive},
		{Phone: "222", Name: "B", Role: model.RolePlayer, Status: model.UserStatusActive},
	}
	users[0].ID = 1
	users[1].ID = 2

	responses := ToUserResponseList(users)

	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(responses))
	}
	if responses[0].Name != "A" {
		t.Errorf("expected name A, got %s", responses[0].Name)
	}
	if responses[1].Role != model.RolePlayer {
		t.Errorf("expected role player, got %s", responses[1].Role)
	}
}

func TestMaskSensitiveData(t *testing.T) {
	tests := []struct {
		name          string
		phone         string
		email         string
		expectedPhone string
		expectedEmail string
	}{
		{
			name:          "standard phone and email",
			phone:         "13800138000",
			email:         "testuser@example.com",
			expectedPhone: "138****8000",
			expectedEmail: "te****@example.com",
		},
		{
			name:          "short phone unchanged",
			phone:         "1380",
			email:         "a@b.com",
			expectedPhone: "1380",     // too short to mask
			expectedEmail: "a@b.com",  // username too short
		},
		{
			name:          "empty fields",
			phone:         "",
			email:         "",
			expectedPhone: "",
			expectedEmail: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &UserResponse{
				Phone: tt.phone,
				Email: tt.email,
			}
			resp.MaskSensitiveData()

			if resp.Phone != tt.expectedPhone {
				t.Errorf("phone: expected %q, got %q", tt.expectedPhone, resp.Phone)
			}
			if resp.Email != tt.expectedEmail {
				t.Errorf("email: expected %q, got %q", tt.expectedEmail, resp.Email)
			}
		})
	}
}

func TestToUserListResponseWithConfig(t *testing.T) {
	users := []model.User{
		{Phone: "13800138000", Name: "A", Role: model.RoleUser, Status: model.UserStatusActive},
		{Phone: "13900139000", Name: "B", Role: model.RoleUser, Status: model.UserStatusActive},
	}
	users[0].ID = 1
	users[1].ID = 2

	t.Run("without masking", func(t *testing.T) {
		resp := ToUserListResponseWithConfig(users, 15, 1, 10, MapperConfig{MaskSensitive: false})

		if resp.Total != 15 {
			t.Errorf("expected total 15, got %d", resp.Total)
		}
		if resp.TotalPages != 2 {
			t.Errorf("expected 2 pages, got %d", resp.TotalPages)
		}
		if resp.Items[0].Phone != "13800138000" {
			t.Errorf("expected unmasked phone, got %s", resp.Items[0].Phone)
		}
	})

	t.Run("with masking", func(t *testing.T) {
		resp := ToUserListResponseWithConfig(users, 15, 1, 10, MapperConfig{MaskSensitive: true})

		if resp.Items[0].Phone != "138****8000" {
			t.Errorf("expected masked phone 138****8000, got %s", resp.Items[0].Phone)
		}
	})

	t.Run("pagination math", func(t *testing.T) {
		resp := ToUserListResponseWithConfig(users, 21, 3, 10, MapperConfig{})

		if resp.TotalPages != 3 {
			t.Errorf("expected 3 pages for 21 items, got %d", resp.TotalPages)
		}
		if resp.Page != 3 {
			t.Errorf("expected page 3, got %d", resp.Page)
		}
	})
}

func TestToCreateUserInput(t *testing.T) {
	req := &CreateUserRequest{
		Phone:     "  13800138000  ",
		Email:     "  test@example.com  ",
		Password:  "MyPassword123",
		Name:      "  Test User  ",
		AvatarURL: "  https://example.com/avatar.png  ",
		Role:      model.RoleUser,
		Status:    model.UserStatusActive,
	}

	input := ToCreateUserInput(req)

	if input.Phone != "13800138000" {
		t.Errorf("expected trimmed phone, got %q", input.Phone)
	}
	if input.Email != "test@example.com" {
		t.Errorf("expected trimmed email, got %q", input.Email)
	}
	if input.Name != "Test User" {
		t.Errorf("expected trimmed name, got %q", input.Name)
	}
	if input.Password != "MyPassword123" {
		t.Errorf("password should not be trimmed, got %q", input.Password)
	}
}
