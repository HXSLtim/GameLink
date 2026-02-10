package dto

import (
	"time"

	"gamelink/internal/model"
)

// CreateUserRequest 创建用户请求 DTO
type CreateUserRequest struct {
	Phone     string           `json:"phone" binding:"omitempty,e164"`           // E.164 格式
	Email     string           `json:"email" binding:"omitempty,email"`          // 邮箱格式
	Password  string           `json:"password" binding:"required,min=8,max=32"` // 密码必填，8-32字符
	Name      string           `json:"name" binding:"required,min=2,max=64"`     // 姓名必填
	Nickname  string           `json:"nickname" binding:"omitempty,max=64"`      // 昵称可选
	AvatarURL string           `json:"avatarUrl" binding:"omitempty,url"`        // 头像URL
	Role      model.Role       `json:"role" binding:"required,oneof=user player admin"`
	Status    model.UserStatus `json:"status" binding:"required,oneof=active suspended banned"`
}

// UpdateUserRequest 更新用户请求 DTO
type UpdateUserRequest struct {
	Phone     *string           `json:"phone" binding:"omitempty,e164"`
	Email     *string           `json:"email" binding:"omitempty,email"`
	Password  *string           `json:"password" binding:"omitempty,min=8,max=32"` // 可选更新密码
	Name      *string           `json:"name" binding:"omitempty,min=2,max=64"`
	Nickname  *string           `json:"nickname" binding:"omitempty,max=64"`
	AvatarURL *string           `json:"avatarUrl" binding:"omitempty,url"`
	Role      *model.Role       `json:"role" binding:"omitempty,oneof=user player admin"`
	Status    *model.UserStatus `json:"status" binding:"omitempty,oneof=active suspended banned"`
}

// UserResponse 用户响应 DTO（不暴露敏感字段）
type UserResponse struct {
	ID        uint64           `json:"id"`
	Phone     string           `json:"phone,omitempty"`      // 脱敏处理可在转换时做
	Email     string           `json:"email,omitempty"`      // 脱敏处理可在转换时做
	Name      string           `json:"name"`
	Nickname  string           `json:"nickname,omitempty"`
	AvatarURL string           `json:"avatarUrl,omitempty"`
	Role      model.Role       `json:"role"`
	Status    model.UserStatus `json:"status"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt time.Time        `json:"updatedAt"`
	
	// VIP 信息（可选）
	VipLevel    *VipLevelBrief `json:"vipLevel,omitempty"`
	VipUnlocked bool           `json:"vipUnlocked"`
	VipExpireAt *time.Time     `json:"vipExpireAt,omitempty"`
	
	// 最后登录
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
}

// VipLevelBrief VIP 等级简要信息
type VipLevelBrief struct {
	ID        uint64 `json:"id"`
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	SortOrder int    `json:"sortOrder"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Items      []UserResponse `json:"items"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"pageSize"`
	TotalPages int            `json:"totalPages"`
}

// UserStatsResponse 用户统计响应
type UserStatsResponse struct {
	TotalUsers      int64          `json:"totalUsers"`
	ActiveUsers     int64          `json:"activeUsers"`
	SuspendedUsers  int64          `json:"suspendedUsers"`
	BannedUsers     int64          `json:"bannedUsers"`
	RoleDistribution map[string]int64 `json:"roleDistribution"` // role -> count
}

// ==================== 转换函数 ====================

// ToUserResponse 将 model.User 转换为 UserResponse
func ToUserResponse(user *model.User) *UserResponse {
	if user == nil {
		return nil
	}
	
	resp := &UserResponse{
		ID:          user.ID,
		Phone:       user.Phone,
		Email:       user.Email,
		Name:        user.Name,
		Nickname:    user.Nickname,
		AvatarURL:   user.AvatarURL,
		Role:        user.Role,
		Status:      user.Status,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		VipUnlocked: user.VipUnlocked,
		VipExpireAt: user.VipExpireAt,
		LastLoginAt: user.LastLoginAt,
	}
	
	// 转换 VIP 等级信息
	if user.VipLevel != nil {
		resp.VipLevel = &VipLevelBrief{
			ID:        user.VipLevel.ID,
			Slug:      user.VipLevel.Slug,
			Title:     user.VipLevel.Title,
			SortOrder: user.VipLevel.SortOrder,
		}
	}
	
	return resp
}

// ToUserResponseList 批量转换
func ToUserResponseList(users []model.User) []UserResponse {
	responses := make([]UserResponse, 0, len(users))
	for i := range users {
		if resp := ToUserResponse(&users[i]); resp != nil {
			responses = append(responses, *resp)
		}
	}
	return responses
}

// MaskSensitiveData 脱敏处理（可选使用）
func (r *UserResponse) MaskSensitiveData() {
	// 手机号脱敏: 138****8000
	if len(r.Phone) >= 11 {
		r.Phone = r.Phone[:3] + "****" + r.Phone[len(r.Phone)-4:]
	}
	
	// 邮箱脱敏: te****@example.com
	if len(r.Email) > 0 {
		parts := splitEmail(r.Email)
		if len(parts) == 2 {
			username := parts[0]
			if len(username) > 2 {
				r.Email = username[:2] + "****@" + parts[1]
			}
		}
	}
}

// splitEmail 分割邮箱
func splitEmail(email string) []string {
	for i := 0; i < len(email); i++ {
		if email[i] == '@' {
			return []string{email[:i], email[i+1:]}
		}
	}
	return []string{email}
}
