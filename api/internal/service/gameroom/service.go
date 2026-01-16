package gameroom

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
	"gamelink/pkg/cache"
)

// Errors specific to game room domain.
var (
	ErrNotFound       = repository.ErrNotFound
	ErrNotMember      = apierr.Forbidden("用户不是房间成员")
	ErrNotHost        = apierr.Forbidden("只有房主可以执行此操作")
	ErrRoomFull       = apierr.BadRequest("房间已满")
	ErrRoomNotWaiting = apierr.BadRequest("房间不在等待状态")
	ErrRoomInactive   = apierr.BadRequest("房间已关闭")
	ErrWrongPassword  = apierr.BadRequest("房间密码错误")
	ErrAlreadyInRoom  = apierr.BadRequest("已在房间中")
	ErrNotAllReady    = apierr.BadRequest("还有成员未准备")
)

// CreateRoomRequest 创建房间请求
type CreateRoomRequest struct {
	Name        string               `json:"name" binding:"required,max=64"`
	GroupType   model.ChatGroupType  `json:"groupType" binding:"required"`
	GameID      uint64               `json:"gameId" binding:"required"`
	MaxMembers  int                  `json:"maxMembers"`
	IsPrivate   bool                 `json:"isPrivate"`
	Password    string               `json:"password"`
	Description string               `json:"description"`
	VoiceEnabled bool                `json:"voiceEnabled"`
}

// UpdateRoomRequest 更新房间请求
type UpdateRoomRequest struct {
	Name         *string `json:"name"`
	Description  *string `json:"description"`
	MaxMembers   *int    `json:"maxMembers"`
	IsPrivate    *bool   `json:"isPrivate"`
	Password     *string `json:"password"`
	VoiceEnabled *bool   `json:"voiceEnabled"`
}

// Service 游戏房间服务
type Service struct {
	groups  repository.ChatGroupRepository
	members repository.ChatMemberRepository
	cache   cache.Cache
}

// NewService 创建游戏房间服务
func NewService(
	groups repository.ChatGroupRepository,
	members repository.ChatMemberRepository,
	cache cache.Cache,
) *Service {
	return &Service{
		groups:  groups,
		members: members,
		cache:   cache,
	}
}

// CreateRoom 创建游戏房间
func (s *Service) CreateRoom(ctx context.Context, hostUserID uint64, req *CreateRoomRequest) (*model.ChatGroup, error) {
	// 验证房间类型
	if req.GroupType != model.ChatGroupTypeTeam &&
		req.GroupType != model.ChatGroupTypeLFG &&
		req.GroupType != model.ChatGroupTypeCustom {
		return nil, apierr.BadRequest("无效的房间类型")
	}

	// 设置默认值
	maxMembers := req.MaxMembers
	if maxMembers <= 0 || maxMembers > 20 {
		maxMembers = 5
	}

	// 创建房间
	room := &model.ChatGroup{
		GroupName:      req.Name,
		GroupType:      req.GroupType,
		CreatedBy:      hostUserID,
		MaxMembers:     maxMembers,
		IsActive:       true,
		Description:    req.Description,
		GameID:         &req.GameID,
		RoomStatus:     model.ChatGroupStatusWaiting,
		IsPrivate:      req.IsPrivate,
		Password:       req.Password,
		CurrentMembers: 1, // 房主
		VoiceEnabled:   req.VoiceEnabled,
	}

	// 如果启用语音，生成语音房间ID
	if req.VoiceEnabled {
		room.VoiceRoomID = generateVoiceRoomID()
		room.VoiceProvider = "trtc"
	}

	if err := s.groups.Create(ctx, room); err != nil {
		return nil, apierr.InternalError("创建房间失败").WithDetails(err.Error())
	}

	// 添加房主为成员
	member := &model.ChatGroupMember{
		GroupID:  room.ID,
		UserID:   hostUserID,
		Role:     model.ChatMemberRoleOwner,
		JoinedAt: time.Now(),
		IsActive: true,
	}
	if err := s.members.Add(ctx, member); err != nil {
		return nil, apierr.InternalError("添加房主失败").WithDetails(err.Error())
	}

	return room, nil
}

// GetRoom 获取房间详情
func (s *Service) GetRoom(ctx context.Context, roomID uint64) (*model.ChatGroup, error) {
	room, err := s.groups.GetWithRelations(ctx, roomID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, apierr.InternalError("获取房间失败").WithDetails(err.Error())
	}
	return room, nil
}

// ListRooms 列出游戏房间
func (s *Service) ListRooms(ctx context.Context, opts repository.GameRoomListOptions) ([]model.ChatGroup, int64, error) {
	rooms, total, err := s.groups.ListGameRooms(ctx, opts)
	if err != nil {
		return nil, 0, apierr.InternalError("获取房间列表失败").WithDetails(err.Error())
	}
	return rooms, total, nil
}

// ListPublicRooms 列出公开房间
func (s *Service) ListPublicRooms(ctx context.Context, gameID *uint64, page, pageSize int) ([]model.ChatGroup, int64, error) {
	rooms, total, err := s.groups.ListPublicRooms(ctx, gameID, page, pageSize)
	if err != nil {
		return nil, 0, apierr.InternalError("获取公开房间列表失败").WithDetails(err.Error())
	}
	return rooms, total, nil
}

// JoinRoom 加入房间
func (s *Service) JoinRoom(ctx context.Context, roomID, userID uint64, password string) error {
	room, err := s.groups.Get(ctx, roomID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return apierr.InternalError("获取房间失败").WithDetails(err.Error())
	}

	// 检查房间状态
	if !room.IsActive {
		return ErrRoomInactive
	}
	if room.RoomStatus != model.ChatGroupStatusWaiting {
		return ErrRoomNotWaiting
	}

	// 检查密码
	if room.IsPrivate && room.Password != "" && room.Password != password {
		return ErrWrongPassword
	}

	// 检查是否已满
	if room.CurrentMembers >= room.MaxMembers {
		return ErrRoomFull
	}

	// 检查是否已在房间
	existingMember, err := s.members.Get(ctx, roomID, userID)
	if err == nil && existingMember.IsActive {
		return ErrAlreadyInRoom
	}

	// 添加成员
	if existingMember != nil {
		// 重新激活
		existingMember.IsActive = true
		existingMember.JoinedAt = time.Now()
		if err := s.members.Update(ctx, existingMember); err != nil {
			return apierr.InternalError("更新成员失败").WithDetails(err.Error())
		}
	} else {
		member := &model.ChatGroupMember{
			GroupID:  roomID,
			UserID:   userID,
			Role:     model.ChatMemberRoleMember,
			JoinedAt: time.Now(),
			IsActive: true,
		}
		if err := s.members.Add(ctx, member); err != nil {
			return apierr.InternalError("添加成员失败").WithDetails(err.Error())
		}
	}

	// 增加成员数
	if err := s.groups.IncrementMemberCount(ctx, roomID); err != nil {
		return apierr.InternalError("更新成员数失败").WithDetails(err.Error())
	}

	return nil
}

// LeaveRoom 离开房间
func (s *Service) LeaveRoom(ctx context.Context, roomID, userID uint64) error {
	room, err := s.groups.Get(ctx, roomID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return apierr.InternalError("获取房间失败").WithDetails(err.Error())
	}

	member, err := s.members.Get(ctx, roomID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotMember
		}
		return apierr.InternalError("获取成员失败").WithDetails(err.Error())
	}

	if !member.IsActive {
		return ErrNotMember
	}

	// 如果是房主离开，需要转移房主或关闭房间
	if room.CreatedBy == userID {
		// 简单处理：关闭房间
		return s.CloseRoom(ctx, roomID, userID)
	}

	// 标记成员离开
	member.IsActive = false
	if err := s.members.Update(ctx, member); err != nil {
		return apierr.InternalError("更新成员失败").WithDetails(err.Error())
	}

	// 减少成员数
	if err := s.groups.DecrementMemberCount(ctx, roomID); err != nil {
		return apierr.InternalError("更新成员数失败").WithDetails(err.Error())
	}

	return nil
}

// UpdateRoom 更新房间信息
func (s *Service) UpdateRoom(ctx context.Context, roomID, userID uint64, req *UpdateRoomRequest) (*model.ChatGroup, error) {
	room, err := s.groups.Get(ctx, roomID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, apierr.InternalError("获取房间失败").WithDetails(err.Error())
	}

	// 检查权限
	if room.CreatedBy != userID {
		return nil, ErrNotHost
	}

	// 更新字段
	if req.Name != nil {
		room.GroupName = *req.Name
	}
	if req.Description != nil {
		room.Description = *req.Description
	}
	if req.MaxMembers != nil && *req.MaxMembers > 0 && *req.MaxMembers <= 20 {
		room.MaxMembers = *req.MaxMembers
	}
	if req.IsPrivate != nil {
		room.IsPrivate = *req.IsPrivate
	}
	if req.Password != nil {
		room.Password = *req.Password
	}
	if req.VoiceEnabled != nil {
		room.VoiceEnabled = *req.VoiceEnabled
		if *req.VoiceEnabled && room.VoiceRoomID == "" {
			room.VoiceRoomID = generateVoiceRoomID()
			room.VoiceProvider = "trtc"
		}
	}

	if err := s.groups.Update(ctx, room); err != nil {
		return nil, apierr.InternalError("更新房间失败").WithDetails(err.Error())
	}

	return room, nil
}

// CloseRoom 关闭房间
func (s *Service) CloseRoom(ctx context.Context, roomID, userID uint64) error {
	room, err := s.groups.Get(ctx, roomID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return apierr.InternalError("获取房间失败").WithDetails(err.Error())
	}

	// 检查权限
	if room.CreatedBy != userID {
		return ErrNotHost
	}

	// 更新状态
	if err := s.groups.UpdateRoomStatus(ctx, roomID, model.ChatGroupStatusCanceled); err != nil {
		return apierr.InternalError("关闭房间失败").WithDetails(err.Error())
	}

	// 停用房间
	if err := s.groups.Deactivate(ctx, roomID); err != nil {
		return apierr.InternalError("停用房间失败").WithDetails(err.Error())
	}

	return nil
}

// StartGame 开始游戏
func (s *Service) StartGame(ctx context.Context, roomID, userID uint64) error {
	room, err := s.groups.Get(ctx, roomID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return apierr.InternalError("获取房间失败").WithDetails(err.Error())
	}

	// 检查权限
	if room.CreatedBy != userID {
		return ErrNotHost
	}

	// 检查状态
	if room.RoomStatus != model.ChatGroupStatusWaiting && room.RoomStatus != model.ChatGroupStatusReady {
		return apierr.BadRequest("房间状态不允许开始游戏")
	}

	// 更新状态为游戏中
	if err := s.groups.UpdateRoomStatus(ctx, roomID, model.ChatGroupStatusInGame); err != nil {
		return apierr.InternalError("更新房间状态失败").WithDetails(err.Error())
	}

	return nil
}

// FinishGame 结束游戏
func (s *Service) FinishGame(ctx context.Context, roomID, userID uint64) error {
	room, err := s.groups.Get(ctx, roomID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return apierr.InternalError("获取房间失败").WithDetails(err.Error())
	}

	// 检查权限
	if room.CreatedBy != userID {
		return ErrNotHost
	}

	// 检查状态
	if room.RoomStatus != model.ChatGroupStatusInGame {
		return apierr.BadRequest("房间不在游戏中")
	}

	// 更新状态为已结束
	if err := s.groups.UpdateRoomStatus(ctx, roomID, model.ChatGroupStatusFinished); err != nil {
		return apierr.InternalError("更新房间状态失败").WithDetails(err.Error())
	}

	return nil
}

// KickMember 踢出成员
func (s *Service) KickMember(ctx context.Context, roomID, hostUserID, targetUserID uint64) error {
	room, err := s.groups.Get(ctx, roomID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return apierr.InternalError("获取房间失败").WithDetails(err.Error())
	}

	// 检查权限
	if room.CreatedBy != hostUserID {
		return ErrNotHost
	}

	// 不能踢自己
	if hostUserID == targetUserID {
		return apierr.BadRequest("不能踢出自己")
	}

	member, err := s.members.Get(ctx, roomID, targetUserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotMember
		}
		return apierr.InternalError("获取成员失败").WithDetails(err.Error())
	}

	if !member.IsActive {
		return ErrNotMember
	}

	// 标记成员离开
	member.IsActive = false
	if err := s.members.Update(ctx, member); err != nil {
		return apierr.InternalError("更新成员失败").WithDetails(err.Error())
	}

	// 减少成员数
	if err := s.groups.DecrementMemberCount(ctx, roomID); err != nil {
		return apierr.InternalError("更新成员数失败").WithDetails(err.Error())
	}

	return nil
}

// GetRoomMembers 获取房间成员列表
func (s *Service) GetRoomMembers(ctx context.Context, roomID uint64) ([]model.ChatGroupMember, error) {
	members, _, err := s.groups.ListMembers(ctx, roomID, repository.ChatGroupMemberListOptions{
		Page:     1,
		PageSize: 100,
	})
	if err != nil {
		return nil, apierr.InternalError("获取成员列表失败").WithDetails(err.Error())
	}

	// 只返回活跃成员
	activeMembers := make([]model.ChatGroupMember, 0)
	for _, m := range members {
		if m.IsActive {
			activeMembers = append(activeMembers, m)
		}
	}

	return activeMembers, nil
}

// GetUserRooms 获取用户的房间列表
func (s *Service) GetUserRooms(ctx context.Context, userID uint64, page, pageSize int) ([]model.ChatGroup, int64, error) {
	rooms, total, err := s.groups.ListByUser(ctx, userID, repository.ChatGroupListOptions{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, 0, apierr.InternalError("获取用户房间列表失败").WithDetails(err.Error())
	}

	// 过滤出游戏房间类型
	gameRooms := make([]model.ChatGroup, 0)
	for _, r := range rooms {
		if r.GroupType == model.ChatGroupTypeTeam ||
			r.GroupType == model.ChatGroupTypeLFG ||
			r.GroupType == model.ChatGroupTypeCustom {
			gameRooms = append(gameRooms, r)
		}
	}

	return gameRooms, total, nil
}

// GetRoomStats 获取房间统计
func (s *Service) GetRoomStats(ctx context.Context) (map[model.ChatGroupStatus]int64, error) {
	return s.groups.CountByRoomStatus(ctx)
}

// ToggleReady 切换成员准备状态
func (s *Service) ToggleReady(ctx context.Context, roomID, userID uint64) (bool, error) {
	room, err := s.groups.Get(ctx, roomID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return false, ErrNotFound
		}
		return false, apierr.InternalError("获取房间失败").WithDetails(err.Error())
	}

	// 检查房间状态
	if !room.IsActive {
		return false, ErrRoomInactive
	}
	if room.RoomStatus != model.ChatGroupStatusWaiting {
		return false, ErrRoomNotWaiting
	}

	// 获取成员
	member, err := s.members.Get(ctx, roomID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return false, ErrNotMember
		}
		return false, apierr.InternalError("获取成员失败").WithDetails(err.Error())
	}

	if !member.IsActive {
		return false, ErrNotMember
	}

	// 切换准备状态
	member.IsReady = !member.IsReady
	if err := s.members.Update(ctx, member); err != nil {
		return false, apierr.InternalError("更新成员状态失败").WithDetails(err.Error())
	}

	// 检查是否所有成员都准备好了，如果是则更新房间状态
	if member.IsReady {
		allReady, err := s.checkAllMembersReady(ctx, roomID)
		if err == nil && allReady {
			_ = s.groups.UpdateRoomStatus(ctx, roomID, model.ChatGroupStatusReady)
		}
	} else {
		// 有人取消准备，房间状态回到等待
		if room.RoomStatus == model.ChatGroupStatusReady {
			_ = s.groups.UpdateRoomStatus(ctx, roomID, model.ChatGroupStatusWaiting)
		}
	}

	return member.IsReady, nil
}

// checkAllMembersReady 检查所有成员是否都准备好了
func (s *Service) checkAllMembersReady(ctx context.Context, roomID uint64) (bool, error) {
	members, err := s.GetRoomMembers(ctx, roomID)
	if err != nil {
		return false, err
	}

	if len(members) < 2 {
		return false, nil // 至少需要2人才能开始
	}

	for _, m := range members {
		// 房主不需要准备
		if m.Role == model.ChatMemberRoleOwner {
			continue
		}
		if !m.IsReady {
			return false, nil
		}
	}
	return true, nil
}

// generateVoiceRoomID 生成语音房间ID
func generateVoiceRoomID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return fmt.Sprintf("voice_%s_%d", hex.EncodeToString(bytes), time.Now().UnixNano())
}
