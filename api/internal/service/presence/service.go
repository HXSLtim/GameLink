package presence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/cache"
)

// Cache key patterns and TTL
const (
	cacheKeyPresence       = "presence:player:%d"
	cacheKeyOnlineCount    = "presence:online:count"
	cacheKeyOnlinePlayers  = "presence:online:players"
	cacheTTLPresence       = 5 * time.Minute
	cacheTTLOnlineCount    = 30 * time.Second
	cacheTTLOnlinePlayers  = 1 * time.Minute

	// Heartbeat timeout - if no heartbeat for this duration, mark as offline
	HeartbeatTimeout = 2 * time.Minute
)

// Service 陪玩师在线状态业务逻辑层
type Service struct {
	repo  repository.PlayerPresenceRepository
	cache cache.Cache
}

// NewPresenceService 创建在线状态服务
func NewPresenceService(repo repository.PlayerPresenceRepository) *Service {
	return &Service{repo: repo}
}

// NewPresenceServiceWithCache 创建带缓存的在线状态服务
func NewPresenceServiceWithCache(repo repository.PlayerPresenceRepository, c cache.Cache) *Service {
	return &Service{repo: repo, cache: c}
}

// SetCache 设置缓存实例
func (s *Service) SetCache(c cache.Cache) {
	s.cache = c
}

// invalidatePresenceCache 清除指定玩家的状态缓存
func (s *Service) invalidatePresenceCache(ctx context.Context, playerID uint64) {
	if s.cache == nil {
		return
	}
	_ = s.cache.Delete(ctx, fmt.Sprintf(cacheKeyPresence, playerID))
	_ = s.cache.Delete(ctx, cacheKeyOnlineCount)
	_ = s.cache.Delete(ctx, cacheKeyOnlinePlayers)
}

// ============================================================================
// 状态管理
// ============================================================================

// GetOrCreatePresence 获取或创建玩家状态记录
func (s *Service) GetOrCreatePresence(ctx context.Context, playerID uint64) (*model.PlayerPresence, error) {
	presence, err := s.repo.GetByPlayerID(ctx, playerID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			// 创建新的状态记录
			presence = &model.PlayerPresence{
				PlayerID:        playerID,
				Status:          model.PresenceOffline,
				LastHeartbeatAt: time.Now(),
			}
			if err := s.repo.Create(ctx, presence); err != nil {
				return nil, fmt.Errorf("create presence: %w", err)
			}
			return presence, nil
		}
		return nil, fmt.Errorf("get presence: %w", err)
	}
	return presence, nil
}

// GetPresence 获取玩家状态
func (s *Service) GetPresence(ctx context.Context, playerID uint64) (*model.PlayerPresence, error) {
	// 尝试从缓存获取
	if s.cache != nil {
		cacheKey := fmt.Sprintf(cacheKeyPresence, playerID)
		if val, ok, _ := s.cache.Get(ctx, cacheKey); ok {
			var presence model.PlayerPresence
			if err := json.Unmarshal([]byte(val), &presence); err == nil {
				return &presence, nil
			}
		}
	}

	presence, err := s.repo.GetByPlayerID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("get presence: %w", err)
	}

	// 写入缓存
	if s.cache != nil {
		cacheKey := fmt.Sprintf(cacheKeyPresence, playerID)
		if data, err := json.Marshal(presence); err == nil {
			_ = s.cache.Set(ctx, cacheKey, string(data), cacheTTLPresence)
		}
	}

	return presence, nil
}

// GetPresenceWithPlayer 获取玩家状态（包含玩家信息）
func (s *Service) GetPresenceWithPlayer(ctx context.Context, id uint64) (*model.PlayerPresence, error) {
	presence, err := s.repo.GetWithPlayer(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get presence with player: %w", err)
	}
	return presence, nil
}

// UpdatePresence 更新玩家状态
func (s *Service) UpdatePresence(ctx context.Context, playerID uint64, req *UpdatePresenceRequest) (*model.PlayerPresence, error) {
	presence, err := s.GetOrCreatePresence(ctx, playerID)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Status != "" {
		presence.Status = req.Status
	}
	if req.CurrentGameID != nil {
		presence.CurrentGameID = req.CurrentGameID
	}
	if req.CurrentGameName != "" {
		presence.CurrentGameName = req.CurrentGameName
	}
	if req.CustomStatus != "" {
		presence.CustomStatus = req.CustomStatus
	}
	if req.CurrentOrderID != nil {
		presence.CurrentOrderID = req.CurrentOrderID
	}
	if req.CurrentRoomID != nil {
		presence.CurrentRoomID = req.CurrentRoomID
	}
	if req.DeviceType != "" {
		presence.DeviceType = req.DeviceType
	}
	presence.LastHeartbeatAt = time.Now()

	if err := s.repo.Update(ctx, presence); err != nil {
		return nil, fmt.Errorf("update presence: %w", err)
	}

	// 清除缓存
	s.invalidatePresenceCache(ctx, playerID)

	return presence, nil
}

// SetStatus 设置玩家状态
func (s *Service) SetStatus(ctx context.Context, playerID uint64, status model.PlayerPresenceStatus) error {
	if err := s.repo.UpdateStatus(ctx, playerID, status); err != nil {
		return fmt.Errorf("set status: %w", err)
	}

	// 清除缓存
	s.invalidatePresenceCache(ctx, playerID)

	return nil
}

// GoOnline 上线
func (s *Service) GoOnline(ctx context.Context, playerID uint64, deviceType string) (*model.PlayerPresence, error) {
	return s.UpdatePresence(ctx, playerID, &UpdatePresenceRequest{
		Status:     model.PresenceOnline,
		DeviceType: deviceType,
	})
}

// GoOffline 下线
func (s *Service) GoOffline(ctx context.Context, playerID uint64) error {
	return s.SetStatus(ctx, playerID, model.PresenceOffline)
}

// SetAccepting 设置为接单中
func (s *Service) SetAccepting(ctx context.Context, playerID uint64) error {
	return s.SetStatus(ctx, playerID, model.PresenceAccepting)
}

// SetInGame 设置为游戏中
func (s *Service) SetInGame(ctx context.Context, playerID uint64, gameID *uint64, gameName string) (*model.PlayerPresence, error) {
	return s.UpdatePresence(ctx, playerID, &UpdatePresenceRequest{
		Status:          model.PresenceInGame,
		CurrentGameID:   gameID,
		CurrentGameName: gameName,
	})
}

// SetMatching 设置为匹配中
func (s *Service) SetMatching(ctx context.Context, playerID uint64) error {
	return s.SetStatus(ctx, playerID, model.PresenceMatching)
}

// SetResting 设置为休息中
func (s *Service) SetResting(ctx context.Context, playerID uint64) error {
	return s.SetStatus(ctx, playerID, model.PresenceResting)
}

// SetInvisible 设置为隐身
func (s *Service) SetInvisible(ctx context.Context, playerID uint64) error {
	return s.SetStatus(ctx, playerID, model.PresenceInvisible)
}

// SetCustomStatus 设置自定义状态文字
func (s *Service) SetCustomStatus(ctx context.Context, playerID uint64, customStatus string) (*model.PlayerPresence, error) {
	return s.UpdatePresence(ctx, playerID, &UpdatePresenceRequest{
		CustomStatus: customStatus,
	})
}

// ============================================================================
// 心跳管理
// ============================================================================

// Heartbeat 心跳保活
func (s *Service) Heartbeat(ctx context.Context, playerID uint64) error {
	if err := s.repo.UpdateHeartbeat(ctx, playerID); err != nil {
		return fmt.Errorf("heartbeat: %w", err)
	}

	// 更新缓存中的心跳时间
	if s.cache != nil {
		cacheKey := fmt.Sprintf(cacheKeyPresence, playerID)
		_ = s.cache.Delete(ctx, cacheKey)
	}

	return nil
}

// CheckAndMarkStaleOffline 检查并标记超时的玩家为离线
func (s *Service) CheckAndMarkStaleOffline(ctx context.Context) (int, error) {
	threshold := time.Now().Add(-HeartbeatTimeout)
	stalePresences, err := s.repo.ListStalePresences(ctx, threshold)
	if err != nil {
		return 0, fmt.Errorf("list stale presences: %w", err)
	}

	if len(stalePresences) == 0 {
		return 0, nil
	}

	playerIDs := make([]uint64, len(stalePresences))
	for i, p := range stalePresences {
		playerIDs[i] = p.PlayerID
	}

	if err := s.repo.BatchUpdateOffline(ctx, playerIDs); err != nil {
		return 0, fmt.Errorf("batch update offline: %w", err)
	}

	// 清除缓存
	for _, playerID := range playerIDs {
		s.invalidatePresenceCache(ctx, playerID)
	}

	return len(playerIDs), nil
}

// ============================================================================
// 查询方法
// ============================================================================

// GetPresencesByPlayerIDs 批量获取玩家状态
func (s *Service) GetPresencesByPlayerIDs(ctx context.Context, playerIDs []uint64) ([]model.PlayerPresence, error) {
	if len(playerIDs) == 0 {
		return []model.PlayerPresence{}, nil
	}

	presences, err := s.repo.ListByPlayerIDs(ctx, playerIDs)
	if err != nil {
		return nil, fmt.Errorf("list presences by player ids: %w", err)
	}
	return presences, nil
}

// ListOnline 获取在线玩家列表
func (s *Service) ListOnline(ctx context.Context, opts repository.PlayerPresenceListOptions) ([]model.PlayerPresence, int64, error) {
	presences, total, err := s.repo.ListOnline(ctx, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("list online: %w", err)
	}
	return presences, total, nil
}

// ListByStatus 按状态获取玩家列表
func (s *Service) ListByStatus(ctx context.Context, status model.PlayerPresenceStatus, page, pageSize int) ([]model.PlayerPresence, int64, error) {
	presences, total, err := s.repo.ListByStatus(ctx, status, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list by status: %w", err)
	}
	return presences, total, nil
}

// CountOnline 获取在线玩家数量
func (s *Service) CountOnline(ctx context.Context) (int64, error) {
	// 尝试从缓存获取
	if s.cache != nil {
		if val, ok, _ := s.cache.Get(ctx, cacheKeyOnlineCount); ok {
			var count int64
			if err := json.Unmarshal([]byte(val), &count); err == nil {
				return count, nil
			}
		}
	}

	count, err := s.repo.CountOnline(ctx)
	if err != nil {
		return 0, fmt.Errorf("count online: %w", err)
	}

	// 写入缓存
	if s.cache != nil {
		if data, err := json.Marshal(count); err == nil {
			_ = s.cache.Set(ctx, cacheKeyOnlineCount, string(data), cacheTTLOnlineCount)
		}
	}

	return count, nil
}

// CountByStatus 按状态统计玩家数量
func (s *Service) CountByStatus(ctx context.Context) (map[model.PlayerPresenceStatus]int64, error) {
	counts, err := s.repo.CountByStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("count by status: %w", err)
	}
	return counts, nil
}

// ============================================================================
// 房间/订单关联
// ============================================================================

// SetCurrentRoom 设置当前房间
func (s *Service) SetCurrentRoom(ctx context.Context, playerID uint64, roomID *uint64) (*model.PlayerPresence, error) {
	return s.UpdatePresence(ctx, playerID, &UpdatePresenceRequest{
		CurrentRoomID: roomID,
	})
}

// SetCurrentOrder 设置当前订单
func (s *Service) SetCurrentOrder(ctx context.Context, playerID uint64, orderID *uint64) (*model.PlayerPresence, error) {
	return s.UpdatePresence(ctx, playerID, &UpdatePresenceRequest{
		CurrentOrderID: orderID,
	})
}

// ClearCurrentRoom 清除当前房间
func (s *Service) ClearCurrentRoom(ctx context.Context, playerID uint64) (*model.PlayerPresence, error) {
	presence, err := s.GetOrCreatePresence(ctx, playerID)
	if err != nil {
		return nil, err
	}

	presence.CurrentRoomID = nil
	if err := s.repo.Update(ctx, presence); err != nil {
		return nil, fmt.Errorf("clear current room: %w", err)
	}

	s.invalidatePresenceCache(ctx, playerID)
	return presence, nil
}

// ClearCurrentOrder 清除当前订单
func (s *Service) ClearCurrentOrder(ctx context.Context, playerID uint64) (*model.PlayerPresence, error) {
	presence, err := s.GetOrCreatePresence(ctx, playerID)
	if err != nil {
		return nil, err
	}

	presence.CurrentOrderID = nil
	if err := s.repo.Update(ctx, presence); err != nil {
		return nil, fmt.Errorf("clear current order: %w", err)
	}

	s.invalidatePresenceCache(ctx, playerID)
	return presence, nil
}

// ============================================================================
// DTO
// ============================================================================

// UpdatePresenceRequest 更新状态请求
type UpdatePresenceRequest struct {
	Status          model.PlayerPresenceStatus `json:"status,omitempty"`
	CurrentGameID   *uint64                    `json:"currentGameId,omitempty"`
	CurrentGameName string                     `json:"currentGameName,omitempty"`
	CustomStatus    string                     `json:"customStatus,omitempty"`
	CurrentOrderID  *uint64                    `json:"currentOrderId,omitempty"`
	CurrentRoomID   *uint64                    `json:"currentRoomId,omitempty"`
	DeviceType      string                     `json:"deviceType,omitempty"`
}
