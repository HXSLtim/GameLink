package dto

import (
	"strings"

	"gamelink/internal/model"
	"gamelink/internal/service/admin"
)

// MapperConfig 转换器配置
type MapperConfig struct {
	MaskSensitive bool // 是否脱敏
}

// ==================== Request DTO → Service Input ====================

// ToCreateUserInput 将创建请求DTO转换为Service层输入
func ToCreateUserInput(req *CreateUserRequest) admin.CreateUserInput {
	return admin.CreateUserInput{
		Phone:     strings.TrimSpace(req.Phone),
		Email:     strings.TrimSpace(req.Email),
		Password:  req.Password,
		Name:      strings.TrimSpace(req.Name),
		AvatarURL: strings.TrimSpace(req.AvatarURL),
		Role:      req.Role,
		Status:    req.Status,
	}
}

// ToUpdateUserInput 将更新请求DTO转换为Service层输入
func ToUpdateUserInput(req *UpdateUserRequest) admin.UpdateUserInput {
	input := admin.UpdateUserInput{}

	if req.Phone != nil {
		trimmed := strings.TrimSpace(*req.Phone)
		input.Phone = trimmed
	}

	if req.Email != nil {
		trimmed := strings.TrimSpace(*req.Email)
		input.Email = trimmed
	}

	if req.Password != nil {
		input.Password = req.Password
	}

	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		input.Name = trimmed
	}

	if req.AvatarURL != nil {
		trimmed := strings.TrimSpace(*req.AvatarURL)
		input.AvatarURL = trimmed
	}

	if req.Role != nil {
		input.Role = *req.Role
	}

	if req.Status != nil {
		input.Status = *req.Status
	}

	return input
}

// ==================== Entity → Response DTO ====================

// ToUserResponseWithConfig 带配置的转换
func ToUserResponseWithConfig(user *model.User, config MapperConfig) *UserResponse {
	resp := ToUserResponse(user)
	if resp != nil && config.MaskSensitive {
		resp.MaskSensitiveData()
	}
	return resp
}

// ToUserListResponseWithConfig 列表转换（带分页）
func ToUserListResponseWithConfig(
	users []model.User,
	total int,
	page, pageSize int,
	config MapperConfig,
) *UserListResponse {
	items := ToUserResponseList(users)

	// 批量脱敏
	if config.MaskSensitive {
		for i := range items {
			items[i].MaskSensitiveData()
		}
	}

	totalPages := total / pageSize
	if total%pageSize > 0 {
		totalPages++
	}

	return &UserListResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
