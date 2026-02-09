package chat

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// NewChatGroupRepository creates a chat group repository implementation.
func NewChatGroupRepository(db *gorm.DB) repository.ChatGroupRepository {
	return &chatGroupRepository{db: db}

}

func (r *chatGroupRepository) ListDeactivatedBefore(ctx context.Context, cutoff time.Time, limit int) ([]model.ChatGroup, error) {
	if limit <= 0 || limit > 1000 {
		limit = 100
	}
	var groups []model.ChatGroup
	if err := r.db.WithContext(ctx).Model(&model.ChatGroup{}).
		Where("group_type = ? AND is_active = ? AND deactivated_at IS NOT NULL AND deactivated_at < ?", model.ChatGroupTypeOrder, false, cutoff).
		Order("deactivated_at ASC").
		Limit(limit).
		Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *chatGroupRepository) DeleteByIDs(ctx context.Context, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Unscoped().
		Where("id IN ?", ids).
		Delete(&model.ChatGroup{}).Error
}

type chatGroupRepository struct {
	db *gorm.DB
}

func (r *chatGroupRepository) Create(ctx context.Context, group *model.ChatGroup) error {
	return r.db.WithContext(ctx).Create(group).Error
}

func (r *chatGroupRepository) Get(ctx context.Context, id uint64) (*model.ChatGroup, error) {
	var group model.ChatGroup
	if err := r.db.WithContext(ctx).
		Preload("Members").
		First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *chatGroupRepository) GetByRelatedOrderID(ctx context.Context, orderID uint64) (*model.ChatGroup, error) {
	var group model.ChatGroup
	if err := r.db.WithContext(ctx).
		Preload("Members").
		Where("related_order_id = ?", orderID).
		First(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *chatGroupRepository) ListByUser(ctx context.Context, userID uint64, opts repository.ChatGroupListOptions) ([]model.ChatGroup, int64, error) {
	tx := r.db.WithContext(ctx).Model(&model.ChatGroup{}).
		Joins("JOIN chat_group_members AS m ON m.group_id = chat_groups.id AND m.user_id = ?", userID)

	if opts.GroupType != nil {
		tx = tx.Where("chat_groups.group_type = ?", *opts.GroupType)
	}
	if !opts.IncludeInactive {
		tx = tx.Where("chat_groups.is_active = ?", true)
	}
	if opts.Keyword != "" {
		like := fmt.Sprintf("%%%s%%", opts.Keyword)
		tx = tx.Where("chat_groups.group_name LIKE ? OR chat_groups.description LIKE ?", like, like)
	}
	if opts.RelatedOrderID != nil {
		tx = tx.Where("chat_groups.related_order_id = ?", *opts.RelatedOrderID)
	}

	page := opts.Page
	if page <= 0 {
		page = 1
	}
	pageSize := opts.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	totalTx := tx.Session(&gorm.Session{})
	var total int64
	if err := totalTx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var groups []model.ChatGroup
	if err := tx.
		Preload("Members").
		Order("chat_groups.updated_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&groups).Error; err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}

func (r *chatGroupRepository) ListMembers(ctx context.Context, groupID uint64, opts repository.ChatGroupMemberListOptions) ([]model.ChatGroupMember, int64, error) {
	tx := r.db.WithContext(ctx).Model(&model.ChatGroupMember{}).
		Where("group_id = ?", groupID)
	if opts.Role != "" {
		tx = tx.Where("role = ?", opts.Role)
	}
	if opts.Keyword != "" {
		like := fmt.Sprintf("%%%s%%", opts.Keyword)
		tx = tx.Where("nickname LIKE ?", like)
	}

	page := opts.Page
	if page <= 0 {
		page = 1
	}
	pageSize := opts.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	totalTx := tx.Session(&gorm.Session{})
	var total int64
	if err := totalTx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var members []model.ChatGroupMember
	if err := tx.
		Order("joined_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&members).Error; err != nil {
		return nil, 0, err
	}

	return members, total, nil
}

func (r *chatGroupRepository) ListPublicChannels(ctx context.Context, page, pageSize int) ([]model.ChatGroup, int64, error) {
	tx := r.db.WithContext(ctx).Model(&model.ChatGroup{}).
		Where("group_type = ?", model.ChatGroupTypePublic).
		Where("is_active = ?", true)

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	if err := tx.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var groups []model.ChatGroup
	if err := tx.
		Preload("Members").
		Order("updated_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&groups).Error; err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}

func (r *chatGroupRepository) Update(ctx context.Context, group *model.ChatGroup) error {
	return r.db.WithContext(ctx).Save(group).Error
}

func (r *chatGroupRepository) Deactivate(ctx context.Context, id uint64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.ChatGroup{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"is_active":      false,
			"deactivated_at": now,
		}).Error
}

// GetWithRelations 获取聊天组及其关联数据
func (r *chatGroupRepository) GetWithRelations(ctx context.Context, id uint64) (*model.ChatGroup, error) {
	var group model.ChatGroup
	if err := r.db.WithContext(ctx).
		Preload("Members").
		Preload("Game").
		First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

// GetByRelatedTeamID 根据关联战队ID获取聊天组
func (r *chatGroupRepository) GetByRelatedTeamID(ctx context.Context, teamID uint64) (*model.ChatGroup, error) {
	var group model.ChatGroup
	if err := r.db.WithContext(ctx).
		Preload("Members").
		Where("related_team_id = ?", teamID).
		First(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

// GetByRelatedLFGID 根据关联LFG请求ID获取聊天组
func (r *chatGroupRepository) GetByRelatedLFGID(ctx context.Context, lfgID uint64) (*model.ChatGroup, error) {
	var group model.ChatGroup
	if err := r.db.WithContext(ctx).
		Preload("Members").
		Where("related_lfg_id = ?", lfgID).
		First(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

// GetByVoiceRoomID 根据语音房间ID获取聊天组
func (r *chatGroupRepository) GetByVoiceRoomID(ctx context.Context, voiceRoomID string) (*model.ChatGroup, error) {
	var group model.ChatGroup
	if err := r.db.WithContext(ctx).
		Preload("Members").
		Where("voice_room_id = ?", voiceRoomID).
		First(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

// UpdateRoomStatus 更新房间状态
func (r *chatGroupRepository) UpdateRoomStatus(ctx context.Context, id uint64, status model.ChatGroupStatus) error {
	updates := map[string]any{
		"room_status": status,
	}
	// 根据状态更新相关时间字段
	now := time.Now()
	switch status {
	case model.ChatGroupStatusInGame:
		updates["started_at"] = now
	case model.ChatGroupStatusFinished, model.ChatGroupStatusCanceled:
		updates["finished_at"] = now
	}
	return r.db.WithContext(ctx).
		Model(&model.ChatGroup{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// ListGameRooms 列出游戏房间
func (r *chatGroupRepository) ListGameRooms(ctx context.Context, opts repository.GameRoomListOptions) ([]model.ChatGroup, int64, error) {
	tx := r.db.WithContext(ctx).Model(&model.ChatGroup{}).
		Where("group_type IN ?", []model.ChatGroupType{
			model.ChatGroupTypeTeam,
			model.ChatGroupTypeLFG,
			model.ChatGroupTypeCustom,
		})

	if opts.GameID != nil {
		tx = tx.Where("game_id = ?", *opts.GameID)
	}
	if opts.HostUserID != nil {
		tx = tx.Where("created_by = ?", *opts.HostUserID)
	}
	if opts.GroupType != nil {
		tx = tx.Where("group_type = ?", *opts.GroupType)
	}
	if opts.RoomStatus != nil {
		tx = tx.Where("room_status = ?", *opts.RoomStatus)
	}
	if len(opts.Statuses) > 0 {
		tx = tx.Where("room_status IN ?", opts.Statuses)
	}
	if opts.IsPrivate != nil {
		tx = tx.Where("is_private = ?", *opts.IsPrivate)
	}
	if opts.Keyword != "" {
		like := fmt.Sprintf("%%%s%%", opts.Keyword)
		tx = tx.Where("group_name LIKE ? OR description LIKE ?", like, like)
	}
	tx = tx.Where("is_active = ?", true)

	page := opts.Page
	if page <= 0 {
		page = 1
	}
	pageSize := opts.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	if err := tx.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var groups []model.ChatGroup
	if err := tx.
		Preload("Members").
		Preload("Game").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&groups).Error; err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}

// ListPublicRooms 列出公开房间
func (r *chatGroupRepository) ListPublicRooms(ctx context.Context, gameID *uint64, page, pageSize int) ([]model.ChatGroup, int64, error) {
	tx := r.db.WithContext(ctx).Model(&model.ChatGroup{}).
		Where("group_type IN ?", []model.ChatGroupType{
			model.ChatGroupTypeTeam,
			model.ChatGroupTypeLFG,
			model.ChatGroupTypeCustom,
		}).
		Where("is_private = ?", false).
		Where("is_active = ?", true).
		Where("room_status IN ?", []model.ChatGroupStatus{
			model.ChatGroupStatusWaiting,
			model.ChatGroupStatusReady,
		})

	if gameID != nil {
		tx = tx.Where("game_id = ?", *gameID)
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	if err := tx.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var groups []model.ChatGroup
	if err := tx.
		Preload("Members").
		Preload("Game").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&groups).Error; err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}

// ListByHostUserID 根据房主用户ID列出房间
func (r *chatGroupRepository) ListByHostUserID(ctx context.Context, hostUserID uint64, status *model.ChatGroupStatus) ([]model.ChatGroup, error) {
	tx := r.db.WithContext(ctx).Model(&model.ChatGroup{}).
		Where("created_by = ?", hostUserID).
		Where("group_type IN ?", []model.ChatGroupType{
			model.ChatGroupTypeTeam,
			model.ChatGroupTypeLFG,
			model.ChatGroupTypeCustom,
		}).
		Where("is_active = ?", true)

	if status != nil {
		tx = tx.Where("room_status = ?", *status)
	}

	var groups []model.ChatGroup
	if err := tx.
		Preload("Members").
		Order("created_at DESC").
		Find(&groups).Error; err != nil {
		return nil, err
	}

	return groups, nil
}

// IncrementMemberCount 增加成员数
func (r *chatGroupRepository) IncrementMemberCount(ctx context.Context, groupID uint64) error {
	return r.db.WithContext(ctx).
		Model(&model.ChatGroup{}).
		Where("id = ?", groupID).
		UpdateColumn("current_members", gorm.Expr("current_members + 1")).Error
}

// DecrementMemberCount 减少成员数
func (r *chatGroupRepository) DecrementMemberCount(ctx context.Context, groupID uint64) error {
	return r.db.WithContext(ctx).
		Model(&model.ChatGroup{}).
		Where("id = ? AND current_members > 0", groupID).
		UpdateColumn("current_members", gorm.Expr("current_members - 1")).Error
}

// CountByRoomStatus 按房间状态统计
func (r *chatGroupRepository) CountByRoomStatus(ctx context.Context) (map[model.ChatGroupStatus]int64, error) {
	type result struct {
		RoomStatus model.ChatGroupStatus
		Count      int64
	}
	var results []result
	if err := r.db.WithContext(ctx).
		Model(&model.ChatGroup{}).
		Where("group_type IN ?", []model.ChatGroupType{
			model.ChatGroupTypeTeam,
			model.ChatGroupTypeLFG,
			model.ChatGroupTypeCustom,
		}).
		Select("room_status, COUNT(*) as count").
		Group("room_status").
		Find(&results).Error; err != nil {
		return nil, err
	}

	counts := make(map[model.ChatGroupStatus]int64)
	for _, r := range results {
		counts[r.RoomStatus] = r.Count
	}
	return counts, nil
}

// ListAll lists all chat groups with admin-level filters (no user scope).
func (r *chatGroupRepository) ListAll(ctx context.Context, opts repository.AdminChatGroupListOptions) ([]model.ChatGroup, int64, error) {
	tx := r.db.WithContext(ctx).Model(&model.ChatGroup{})

	if opts.GroupType != nil {
		tx = tx.Where("group_type = ?", *opts.GroupType)
	}
	if opts.IsActive != nil {
		tx = tx.Where("is_active = ?", *opts.IsActive)
	}
	if opts.Keyword != "" {
		like := fmt.Sprintf("%%%s%%", opts.Keyword)
		tx = tx.Where("group_name LIKE ? OR description LIKE ?", like, like)
	}
	if opts.RelatedOrderID != nil {
		tx = tx.Where("related_order_id = ?", *opts.RelatedOrderID)
	}
	if opts.UserID != nil {
		tx = tx.Joins("JOIN chat_group_members AS m ON m.group_id = chat_groups.id AND m.user_id = ?", *opts.UserID)
	}
	if opts.DateFrom != nil {
		tx = tx.Where("chat_groups.created_at >= ?", *opts.DateFrom)
	}
	if opts.DateTo != nil {
		tx = tx.Where("chat_groups.created_at <= ?", *opts.DateTo)
	}

	page := opts.Page
	if page <= 0 {
		page = 1
	}
	pageSize := opts.PageSize
	if pageSize <= 0 || pageSize > 1000 {
		pageSize = 20
	}

	var total int64
	if err := tx.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var groups []model.ChatGroup
	if err := tx.
		Preload("Members").
		Order("chat_groups.updated_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&groups).Error; err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}

// CountAll returns the total number of chat groups.
func (r *chatGroupRepository) CountAll(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.ChatGroup{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Reactivate reactivates a deactivated chat group.
func (r *chatGroupRepository) Reactivate(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Model(&model.ChatGroup{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"is_active":      true,
			"deactivated_at": nil,
		}).Error
}

// CountActiveRooms 统计活跃房间数
func (r *chatGroupRepository) CountActiveRooms(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.ChatGroup{}).
		Where("group_type IN ?", []model.ChatGroupType{
			model.ChatGroupTypeTeam,
			model.ChatGroupTypeLFG,
			model.ChatGroupTypeCustom,
		}).
		Where("is_active = ?", true).
		Where("room_status IN ?", []model.ChatGroupStatus{
			model.ChatGroupStatusWaiting,
			model.ChatGroupStatusReady,
			model.ChatGroupStatusInGame,
		}).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
