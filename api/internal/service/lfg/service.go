package lfg

import (
	"context"
	"errors"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	gameroomservice "gamelink/internal/service/gameroom"
	"gamelink/pkg/apierr"
	"gamelink/pkg/cache"
)

// Errors specific to LFG domain.
var (
	ErrNotFound         = repository.ErrNotFound
	ErrAlreadyHasActive = apierr.BadRequest("已有进行中的匹配请求")
	ErrRequestExpired   = apierr.BadRequest("匹配请求已过期")
	ErrNotOwner         = apierr.Forbidden("只能操作自己的匹配请求")
	ErrCannotMatch      = apierr.BadRequest("无法匹配该请求")
)

// DefaultExpireDuration 默认过期时间
const DefaultExpireDuration = 30 * time.Minute

// CreateLFGRequest 创建LFG请求
type CreateLFGRequest struct {
	GameID          uint64               `json:"gameId" binding:"required"`
	RequestType     model.LFGRequestType `json:"requestType" binding:"required"`
	Title           string               `json:"title" binding:"max=64"`
	Description     string               `json:"description" binding:"max=256"`
	RequiredPlayers int                  `json:"requiredPlayers"`
	MinRank         string               `json:"minRank"`
	MaxPriceCents   int64                `json:"maxPriceCents"`
	ExpireMinutes   int                  `json:"expireMinutes"`
}

// Service LFG服务
type Service struct {
	lfgRepo     repository.LFGRequestRepository
	roomService *gameroomservice.Service
	cache       cache.Cache
}

// NewService 创建LFG服务
func NewService(
	lfgRepo repository.LFGRequestRepository,
	roomService *gameroomservice.Service,
	cache cache.Cache,
) *Service {
	return &Service{
		lfgRepo:     lfgRepo,
		roomService: roomService,
		cache:       cache,
	}
}

// CreateRequest 创建匹配请求
func (s *Service) CreateRequest(ctx context.Context, userID uint64, req *CreateLFGRequest) (*model.LFGRequest, error) {
	// 检查是否已有进行中的请求
	existing, err := s.lfgRepo.GetActiveByUserID(ctx, userID)
	if err == nil && existing != nil {
		return nil, ErrAlreadyHasActive
	}
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, apierr.InternalError("检查现有请求失败").WithDetails(err.Error())
	}

	// 设置默认值
	requiredPlayers := req.RequiredPlayers
	if requiredPlayers <= 0 {
		requiredPlayers = 1
	}

	expireMinutes := req.ExpireMinutes
	if expireMinutes <= 0 || expireMinutes > 120 {
		expireMinutes = 30
	}

	lfgRequest := &model.LFGRequest{
		UserID:          userID,
		GameID:          req.GameID,
		RequestType:     req.RequestType,
		Title:           req.Title,
		Description:     req.Description,
		RequiredPlayers: requiredPlayers,
		MinRank:         req.MinRank,
		MaxPriceCents:   req.MaxPriceCents,
		Status:          model.LFGPending,
		ExpiresAt:       time.Now().Add(time.Duration(expireMinutes) * time.Minute),
	}

	if err := s.lfgRepo.Create(ctx, lfgRequest); err != nil {
		return nil, apierr.InternalError("创建匹配请求失败").WithDetails(err.Error())
	}

	return lfgRequest, nil
}

// GetRequest 获取匹配请求
func (s *Service) GetRequest(ctx context.Context, id uint64) (*model.LFGRequest, error) {
	request, err := s.lfgRepo.GetWithRelations(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, apierr.InternalError("获取匹配请求失败").WithDetails(err.Error())
	}
	return request, nil
}

// GetActiveRequest 获取用户当前活跃的请求
func (s *Service) GetActiveRequest(ctx context.Context, userID uint64) (*model.LFGRequest, error) {
	request, err := s.lfgRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, apierr.InternalError("获取活跃请求失败").WithDetails(err.Error())
	}
	return request, nil
}

// CancelRequest 取消匹配请求
func (s *Service) CancelRequest(ctx context.Context, id, userID uint64) error {
	request, err := s.lfgRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return apierr.InternalError("获取匹配请求失败").WithDetails(err.Error())
	}

	// 检查权限
	if request.UserID != userID {
		return ErrNotOwner
	}

	// 检查状态
	if request.Status != model.LFGPending {
		return apierr.BadRequest("只能取消等待中的请求")
	}

	if err := s.lfgRepo.UpdateStatus(ctx, id, model.LFGCanceled); err != nil {
		return apierr.InternalError("取消请求失败").WithDetails(err.Error())
	}

	return nil
}

// ListRequests 列出匹配请求
func (s *Service) ListRequests(ctx context.Context, opts repository.LFGRequestListOptions) ([]model.LFGRequest, int64, error) {
	requests, total, err := s.lfgRepo.List(ctx, opts)
	if err != nil {
		return nil, 0, apierr.InternalError("获取匹配请求列表失败").WithDetails(err.Error())
	}
	return requests, total, nil
}

// ListPendingRequests 列出等待中的请求
func (s *Service) ListPendingRequests(ctx context.Context, gameID *uint64, page, pageSize int) ([]model.LFGRequest, int64, error) {
	requests, total, err := s.lfgRepo.ListPending(ctx, gameID, page, pageSize)
	if err != nil {
		return nil, 0, apierr.InternalError("获取等待中请求列表失败").WithDetails(err.Error())
	}
	return requests, total, nil
}

// ListUserRequests 列出用户的请求
func (s *Service) ListUserRequests(ctx context.Context, userID uint64, status *model.LFGRequestStatus) ([]model.LFGRequest, error) {
	requests, err := s.lfgRepo.ListByUserID(ctx, userID, status)
	if err != nil {
		return nil, apierr.InternalError("获取用户请求列表失败").WithDetails(err.Error())
	}
	return requests, nil
}

// AcceptRequest 接受匹配请求（陪玩师接单）
func (s *Service) AcceptRequest(ctx context.Context, requestID, acceptorUserID uint64) (*model.ChatGroup, error) {
	request, err := s.lfgRepo.Get(ctx, requestID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, apierr.InternalError("获取匹配请求失败").WithDetails(err.Error())
	}

	// 检查状态
	if !request.CanMatch() {
		return nil, ErrCannotMatch
	}

	// 不能接受自己的请求
	if request.UserID == acceptorUserID {
		return nil, apierr.BadRequest("不能接受自己的请求")
	}

	// 创建游戏房间
	room, err := s.roomService.CreateRoom(ctx, request.UserID, &gameroomservice.CreateRoomRequest{
		Name:        request.Title,
		GroupType:   model.ChatGroupTypeLFG,
		GameID:      request.GameID,
		MaxMembers:  request.RequiredPlayers + 1,
		IsPrivate:   false,
		Description: request.Description,
	})
	if err != nil {
		return nil, apierr.InternalError("创建房间失败").WithDetails(err.Error())
	}

	// 接受者加入房间
	if err := s.roomService.JoinRoom(ctx, room.ID, acceptorUserID, ""); err != nil {
		return nil, apierr.InternalError("加入房间失败").WithDetails(err.Error())
	}

	// 更新请求状态
	if err := s.lfgRepo.UpdateMatched(ctx, requestID, room.ID); err != nil {
		return nil, apierr.InternalError("更新请求状态失败").WithDetails(err.Error())
	}

	// 更新房间关联
	room.RelatedLFGID = &requestID

	return room, nil
}

// FindMatches 查找匹配的请求
func (s *Service) FindMatches(ctx context.Context, requestID uint64, limit int) ([]model.LFGRequest, error) {
	request, err := s.lfgRepo.Get(ctx, requestID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, apierr.InternalError("获取匹配请求失败").WithDetails(err.Error())
	}

	if limit <= 0 {
		limit = 10
	}

	matches, err := s.lfgRepo.FindMatchingRequests(ctx, request, limit)
	if err != nil {
		return nil, apierr.InternalError("查找匹配失败").WithDetails(err.Error())
	}

	return matches, nil
}

// ExpireStaleRequests 过期处理（后台任务调用）
func (s *Service) ExpireStaleRequests(ctx context.Context, batchSize int) (int, error) {
	if batchSize <= 0 {
		batchSize = 100
	}

	expired, err := s.lfgRepo.ListExpired(ctx, batchSize)
	if err != nil {
		return 0, apierr.InternalError("获取过期请求失败").WithDetails(err.Error())
	}

	if len(expired) == 0 {
		return 0, nil
	}

	ids := make([]uint64, len(expired))
	for i, r := range expired {
		ids[i] = r.ID
	}

	if err := s.lfgRepo.BatchExpire(ctx, ids); err != nil {
		return 0, apierr.InternalError("批量过期失败").WithDetails(err.Error())
	}

	return len(ids), nil
}

// GetStats 获取统计信息
func (s *Service) GetStats(ctx context.Context) (map[model.LFGRequestStatus]int64, error) {
	return s.lfgRepo.CountByStatus(ctx)
}

// CountPending 获取等待中的请求数量
func (s *Service) CountPending(ctx context.Context, gameID *uint64) (int64, error) {
	return s.lfgRepo.CountPending(ctx, gameID)
}
