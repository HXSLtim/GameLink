// Package auth provides role switching functionality for mini-program users.
package auth

import (
	"context"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
	"gamelink/pkg/auth"
)

// RoleSwitchRequest 角色切换请求
type RoleSwitchRequest struct {
	TargetRole string `json:"targetRole" binding:"required,oneof=user player"` // 目标角色
}

// RoleSwitchResponse 角色切换响应
type RoleSwitchResponse struct {
	AccessToken string `json:"accessToken"`
	CurrentRole string `json:"currentRole"`
	IsPlayer    bool   `json:"isPlayer"`
}

// AvailableRole 可用角色
type AvailableRole struct {
	Role        string `json:"role"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Available   bool   `json:"available"`
}

// AvailableRolesResponse 可用角色列表响应
type AvailableRolesResponse struct {
	CurrentRole string          `json:"currentRole"`
	Roles       []AvailableRole `json:"roles"`
}

// RoleService 角色切换服务
type RoleService struct {
	users   repository.UserRepository
	players repository.PlayerRepository
}

// NewRoleService 创建角色切换服务
func NewRoleService(users repository.UserRepository, players repository.PlayerRepository) *RoleService {
	return &RoleService{
		users:   users,
		players: players,
	}
}

// GetAvailableRoles 获取用户可用角色列表
func (s *RoleService) GetAvailableRoles(ctx context.Context, userID uint64, currentRole string) (*AvailableRolesResponse, error) {
	// 检查用户是否存在
	user, err := s.users.Get(ctx, userID)
	if err != nil {
		return nil, apierr.NotFound("用户不存在")
	}

	if user.Status != model.UserStatusActive {
		return nil, apierr.Forbidden("用户已被禁用")
	}

	// 检查是否是陪玩师
	isPlayer := s.checkIsPlayer(ctx, userID)

	roles := []AvailableRole{
		{
			Role:        "user",
			Name:        "用户",
			Description: "浏览陪玩师、下单、评价",
			Available:   true,
		},
		{
			Role:        "player",
			Name:        "陪玩师",
			Description: "接单、管理服务、收益追踪",
			Available:   isPlayer,
		},
	}

	return &AvailableRolesResponse{
		CurrentRole: currentRole,
		Roles:       roles,
	}, nil
}

// SwitchRole 切换用户角色
func (s *RoleService) SwitchRole(ctx context.Context, userID uint64, targetRole string) (*RoleSwitchResponse, error) {
	// 检查用户是否存在
	user, err := s.users.Get(ctx, userID)
	if err != nil {
		return nil, apierr.NotFound("用户不存在")
	}

	if user.Status != model.UserStatusActive {
		return nil, apierr.Forbidden("用户已被禁用")
	}

	// 检查是否是陪玩师
	isPlayer := s.checkIsPlayer(ctx, userID)

	// 验证目标角色
	if targetRole == "player" && !isPlayer {
		return nil, apierr.Forbidden("您还不是认证陪玩师，无法切换到陪玩师角色")
	}

	if targetRole != "user" && targetRole != "player" {
		return nil, apierr.BadRequest("无效的目标角色")
	}

	// 生成新的 Token
	claims := auth.CustomClaims{
		UserID:      user.ID,
		Role:        string(user.Role),
		IsPlayer:    isPlayer,
		CurrentRole: targetRole,
	}

	accessToken, err := auth.GenerateToken(claims)
	if err != nil {
		return nil, apierr.InternalError("生成Token失败")
	}

	return &RoleSwitchResponse{
		AccessToken: accessToken,
		CurrentRole: targetRole,
		IsPlayer:    isPlayer,
	}, nil
}

// checkIsPlayer 检查用户是否是认证陪玩师
func (s *RoleService) checkIsPlayer(ctx context.Context, userID uint64) bool {
	if s.players == nil {
		return false
	}
	player, err := s.players.GetByUserID(ctx, userID)
	return err == nil && player != nil && player.VerificationStatus == model.VerificationVerified
}
