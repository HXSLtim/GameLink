package player

import (
	"context"
	"errors"

	"gamelink/internal/cache"
	"gamelink/internal/model"
	"gamelink/internal/repository"
)

var (
	// ErrNotFound 玩家不存在
	ErrNotFound = repository.ErrNotFound
	// ErrValidation 表示输入校验失败
	ErrValidation = errors.New("validation failed")
	// ErrPlayerNotVerified 陪玩师未通过审核
	ErrPlayerNotVerified = errors.New("player not verified")
	// ErrAlreadyPlayer 用户已是陪玩师
	ErrAlreadyPlayer = errors.New("user is already a player")
)

// PlayerService 陪玩师服务
//
// 功能：
// 1. 用户端查询陪玩师列表和详情
// 2. 陪玩师端管理个人资料
// 3. 陪玩师申请和审核
type PlayerService struct {
	players    repository.PlayerRepository
	users      repository.UserRepository
	games      repository.GameRepository
	orders     repository.OrderRepository
	reviews    repository.ReviewRepository
	playerTags repository.PlayerTagRepository
	cache      cache.Cache
}

// NewPlayerService 创建陪玩师服务
func NewPlayerService(
	players repository.PlayerRepository,
	users repository.UserRepository,
	games repository.GameRepository,
	orders repository.OrderRepository,
	reviews repository.ReviewRepository,
	playerTags repository.PlayerTagRepository,
	cache cache.Cache,
) *PlayerService {
	return &PlayerService{
		players:    players,
		users:      users,
		games:      games,
		orders:     orders,
		reviews:    reviews,
		playerTags: playerTags,
		cache:      cache,
	}
}

// PlayerCardDTO 陪玩师卡片信息（列表展示）
type PlayerCardDTO struct {
	ID              uint64  `json:"id"`
	UserID          uint64  `json:"userId"`
	Nickname        string  `json:"nickname"`
	AvatarURL       string  `json:"avatarUrl"`
	Bio             string  `json:"bio"`
	Rank            string  `json:"rank"`
	RatingAverage   float32 `json:"ratingAverage"`
	RatingCount     uint32  `json:"ratingCount"`
	HourlyRateCents int64   `json:"hourlyRateCents"`
	MainGame        string  `json:"mainGame"`   // 游戏名称
	IsOnline        bool    `json:"isOnline"`   // 在线状态
	OrderCount      int64   `json:"orderCount"` // 历史订单数
}

// PlayerDetailDTO 陪玩师详情信息
type PlayerDetailDTO struct {
	PlayerCardDTO
	Tags           []string `json:"tags"`           // 服务标签
	GoodRatio      float32  `json:"goodRatio"`      // 好评率
	AvgResponseMin int      `json:"avgResponseMin"` // 平均响应时间（分钟）
}

// PlayerStatsDTO 陪玩师统计数据
type PlayerStatsDTO struct {
	TotalOrders     int64   `json:"totalOrders"`
	CompletedOrders int64   `json:"completedOrders"`
	RepeatRate      float32 `json:"repeatRate"` // 复购率
}

// PlayerListRequest 陪玩师列表请求
type PlayerListRequest struct {
	GameID     *uint64  `form:"gameId"`     // 游戏筛选
	MinPrice   *int64   `form:"minPrice"`   // 最低价格（分）
	MaxPrice   *int64   `form:"maxPrice"`   // 最高价格（分）
	MinRating  *float32 `form:"minRating"`  // 最低评分
	OnlineOnly bool     `form:"onlineOnly"` // 仅在线
	SortBy     string   `form:"sortBy"`     // 排序：price/rating/orders
	Page       int      `form:"page"`
	PageSize   int      `form:"pageSize"`
}

// PlayerListResponse 陪玩师列表响应
type PlayerListResponse struct {
	Players []PlayerCardDTO `json:"players"`
	Total   int64           `json:"total"`
}

// PlayerDetailResponse 陪玩师详情响应
type PlayerDetailResponse struct {
	Player  PlayerDetailDTO `json:"player"`
	Reviews []ReviewDTO     `json:"reviews"` // 最新评价
	Stats   PlayerStatsDTO  `json:"stats"`
}

// ReviewDTO 评价信息
type ReviewDTO struct {
	ID            uint64 `json:"id"`
	UserNickname  string `json:"userNickname"`
	UserAvatarURL string `json:"userAvatarUrl"`
	Rating        int    `json:"rating"`
	Comment       string `json:"comment"`
	CreatedAt     string `json:"createdAt"`
}

// ApplyPlayerRequest 申请成为陪玩师请求
type ApplyPlayerRequest struct {
	Nickname        string   `json:"nickname" binding:"required,max=64"`
	Bio             string   `json:"bio" binding:"max=500"`
	MainGameID      uint64   `json:"mainGameId" binding:"required"`
	Rank            string   `json:"rank" binding:"required,max=32"`
	HourlyRateCents int64    `json:"hourlyRateCents" binding:"required,min=1000"`
	Tags            []string `json:"tags"`
	ProofImages     []string `json:"proofImages"` // 段位证明图片
}

// ApplyPlayerResponse 申请陪玩师响应
type ApplyPlayerResponse struct {
	PlayerID           uint64                   `json:"playerId"`
	VerificationStatus model.VerificationStatus `json:"verificationStatus"`
}

// UpdatePlayerProfileRequest 更新陪玩师资料请求
type UpdatePlayerProfileRequest struct {
	Nickname        string   `json:"nickname" binding:"required,max=64"`
	Bio             string   `json:"bio" binding:"max=500"`
	Rank            string   `json:"rank"`
	HourlyRateCents int64    `json:"hourlyRateCents" binding:"min=1000"`
	Tags            []string `json:"tags"`
}

// SetPlayerStatusRequest 设置在线状态请求
type SetPlayerStatusRequest struct {
	Online bool `json:"online"`
}

// ListPlayers 获取陪玩师列表（用户端）
func (s *PlayerService) ListPlayers(ctx context.Context, req PlayerListRequest) (*PlayerListResponse, error) {
	// 默认分页参数
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// 获取所有陪玩师
	players, total, err := s.players.ListPaged(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	playerCards := make([]PlayerCardDTO, 0, len(players))
	for _, p := range players {
		// 过滤条件
		if req.GameID != nil && p.MainGameID != *req.GameID {
			continue
		}
		if req.MinPrice != nil && p.HourlyRateCents < *req.MinPrice {
			continue
		}
		if req.MaxPrice != nil && p.HourlyRateCents > *req.MaxPrice {
			continue
		}
		if req.MinRating != nil && p.RatingAverage < *req.MinRating {
			continue
		}
		// 只显示已审核通过的陪玩师
		if p.VerificationStatus != model.VerificationVerified {
			continue
		}

		// 获取用户信息
		user, err := s.users.Get(ctx, p.UserID)
		if err != nil {
			continue
		}

		// 获取游戏信息
		var gameName string
		if p.MainGameID > 0 {
			game, err := s.games.Get(ctx, p.MainGameID)
			if err == nil {
				gameName = game.Name
			}
		}

		// 获取订单数
		orderCount, err := s.getPlayerOrderCount(ctx, p.ID)
		if err != nil {
			orderCount = 0
		}

		playerCards = append(playerCards, PlayerCardDTO{
			ID:              p.ID,
			UserID:          p.UserID,
			Nickname:        p.Nickname,
			AvatarURL:       user.AvatarURL,
			Bio:             p.Bio,
			Rank:            p.Rank,
			RatingAverage:   p.RatingAverage,
			RatingCount:     p.RatingCount,
			HourlyRateCents: p.HourlyRateCents,
			MainGame:        gameName,
			IsOnline:        false, // TODO: 从 Redis 获取在线状态
			OrderCount:      orderCount,
		})
	}

	return &PlayerListResponse{
		Players: playerCards,
		Total:   total,
	}, nil
}

// GetPlayerDetail 获取陪玩师详情（用户端）
func (s *PlayerService) GetPlayerDetail(ctx context.Context, id uint64) (*PlayerDetailResponse, error) {
	// 获取陪玩师信息
	player, err := s.players.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// 只显示已审核通过的陪玩师
	if player.VerificationStatus != model.VerificationVerified {
		return nil, ErrPlayerNotVerified
	}

	// 获取用户信息
	user, err := s.users.Get(ctx, player.UserID)
	if err != nil {
		return nil, err
	}

	// 获取游戏信息
	var gameName string
	if player.MainGameID > 0 {
		game, err := s.games.Get(ctx, player.MainGameID)
		if err == nil {
			gameName = game.Name
		}
	}

	// 获取标签
	tags, err := s.playerTags.GetTags(ctx, player.ID)
	if err != nil {
		tags = []string{}
	}

	// 获取订单统计
	stats, err := s.getPlayerStats(ctx, player.ID)
	if err != nil {
		stats = PlayerStatsDTO{}
	}

	// 获取评价列表
	reviews, err := s.getPlayerReviews(ctx, player.ID, 5)
	if err != nil {
		reviews = []ReviewDTO{}
	}

	// 计算好评率
	goodRatio := s.calculateGoodRatio(ctx, player.ID)

	// 获取订单数
	orderCount, _ := s.getPlayerOrderCount(ctx, player.ID)

	return &PlayerDetailResponse{
		Player: PlayerDetailDTO{
			PlayerCardDTO: PlayerCardDTO{
				ID:              player.ID,
				UserID:          player.UserID,
				Nickname:        player.Nickname,
				AvatarURL:       user.AvatarURL,
				Bio:             player.Bio,
				Rank:            player.Rank,
				RatingAverage:   player.RatingAverage,
				RatingCount:     player.RatingCount,
				HourlyRateCents: player.HourlyRateCents,
				MainGame:        gameName,
				IsOnline:        false, // TODO: 从 Redis 获取在线状态
				OrderCount:      orderCount,
			},
			Tags:           tags,
			GoodRatio:      goodRatio,
			AvgResponseMin: 30, // TODO: 计算平均响应时间
		},
		Reviews: reviews,
		Stats:   stats,
	}, nil
}

// ApplyAsPlayer 申请成为陪玩师
func (s *PlayerService) ApplyAsPlayer(ctx context.Context, userID uint64, req ApplyPlayerRequest) (*ApplyPlayerResponse, error) {
	// 检查用户是否存在
	user, err := s.users.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 检查用户是否已经是陪玩师
	players, _, err := s.players.ListPaged(ctx, 1, 1)
	if err != nil {
		return nil, err
	}
	for _, p := range players {
		if p.UserID == userID {
			return nil, ErrAlreadyPlayer
		}
	}

	// 验证游戏ID
	if _, err := s.games.Get(ctx, req.MainGameID); err != nil {
		return nil, errors.New("invalid game id")
	}

	// 创建陪玩师资料
	player := &model.Player{
		UserID:             userID,
		Nickname:           req.Nickname,
		Bio:                req.Bio,
		Rank:               req.Rank,
		HourlyRateCents:    req.HourlyRateCents,
		MainGameID:         req.MainGameID,
		VerificationStatus: model.VerificationPending, // 待审核
	}

	if err := s.players.Create(ctx, player); err != nil {
		return nil, err
	}

	// 保存标签
	if len(req.Tags) > 0 {
		if err := s.playerTags.ReplaceTags(ctx, player.ID, req.Tags); err != nil {
			// 标签保存失败不影响申请
		}
	}

	// 更新用户角色为 player（可选）
	user.Role = model.RolePlayer
	_ = s.users.Update(ctx, user)

	return &ApplyPlayerResponse{
		PlayerID:           player.ID,
		VerificationStatus: player.VerificationStatus,
	}, nil
}

// GetPlayerProfile 获取陪玩师自己的资料
func (s *PlayerService) GetPlayerProfile(ctx context.Context, userID uint64) (*PlayerDetailResponse, error) {
	// 查找该用户的陪玩师资料
	players, _, err := s.players.ListPaged(ctx, 1, 100)
	if err != nil {
		return nil, err
	}

	var player *model.Player
	for _, p := range players {
		if p.UserID == userID {
			player = &p
			break
		}
	}

	if player == nil {
		return nil, ErrNotFound
	}

	// 返回详情（不需要检查审核状态，陪玩师可以查看自己的资料）
	return s.GetPlayerDetail(ctx, player.ID)
}

// UpdatePlayerProfile 更新陪玩师资料
func (s *PlayerService) UpdatePlayerProfile(ctx context.Context, userID uint64, req UpdatePlayerProfileRequest) error {
	// 查找该用户的陪玩师资料
	players, _, err := s.players.ListPaged(ctx, 1, 100)
	if err != nil {
		return err
	}

	var player *model.Player
	for _, p := range players {
		if p.UserID == userID {
			player = &p
			break
		}
	}

	if player == nil {
		return ErrNotFound
	}

	// 更新资料
	player.Nickname = req.Nickname
	player.Bio = req.Bio
	if req.Rank != "" {
		player.Rank = req.Rank
	}
	if req.HourlyRateCents > 0 {
		player.HourlyRateCents = req.HourlyRateCents
	}

	if err := s.players.Update(ctx, player); err != nil {
		return err
	}

	// 更新标签
	if len(req.Tags) > 0 {
		if err := s.playerTags.ReplaceTags(ctx, player.ID, req.Tags); err != nil {
			// 标签更新失败不影响资料更新
		}
	}

	return nil
}

// SetPlayerOnlineStatus 设置陪玩师在线状态
func (s *PlayerService) SetPlayerOnlineStatus(ctx context.Context, userID uint64, online bool) error {
	// TODO: 使用 Redis 存储在线状态
	// 这里先返回 nil，表示操作成功
	return nil
}

// getPlayerOrderCount 获取陪玩师的订单数量
func (s *PlayerService) getPlayerOrderCount(ctx context.Context, playerID uint64) (int64, error) {
	playerIDPtr := &playerID
	orders, total, err := s.orders.List(ctx, repository.OrderListOptions{
		PlayerID: playerIDPtr,
		Page:     1,
		PageSize: 1,
	})
	if err != nil {
		return 0, err
	}
	_ = orders
	return total, nil
}

// getPlayerStats 获取陪玩师统计数据
func (s *PlayerService) getPlayerStats(ctx context.Context, playerID uint64) (PlayerStatsDTO, error) {
	playerIDPtr := &playerID

	// 总订单数
	_, total, err := s.orders.List(ctx, repository.OrderListOptions{
		PlayerID: playerIDPtr,
		Page:     1,
		PageSize: 1,
	})
	if err != nil {
		return PlayerStatsDTO{}, err
	}

	// 已完成订单数
	_, completed, err := s.orders.List(ctx, repository.OrderListOptions{
		PlayerID: playerIDPtr,
		Statuses: []model.OrderStatus{model.OrderStatusCompleted},
		Page:     1,
		PageSize: 1,
	})
	if err != nil {
		return PlayerStatsDTO{}, err
	}

	return PlayerStatsDTO{
		TotalOrders:     total,
		CompletedOrders: completed,
		RepeatRate:      0.0, // TODO: 计算复购率
	}, nil
}

// getPlayerReviews 获取陪玩师的评价列表
func (s *PlayerService) getPlayerReviews(ctx context.Context, playerID uint64, limit int) ([]ReviewDTO, error) {
	playerIDPtr := &playerID
	reviews, _, err := s.reviews.List(ctx, repository.ReviewListOptions{
		PlayerID: playerIDPtr,
		Page:     1,
		PageSize: limit,
	})
	if err != nil {
		return nil, err
	}

	reviewDTOs := make([]ReviewDTO, 0, len(reviews))
	for _, r := range reviews {
		user, err := s.users.Get(ctx, r.UserID)
		if err != nil {
			continue
		}

		reviewDTOs = append(reviewDTOs, ReviewDTO{
			ID:            r.ID,
			UserNickname:  user.Name,
			UserAvatarURL: user.AvatarURL,
			Rating:        int(r.Score),
			Comment:       r.Content,
			CreatedAt:     r.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return reviewDTOs, nil
}

// calculateGoodRatio 计算好评率
func (s *PlayerService) calculateGoodRatio(ctx context.Context, playerID uint64) float32 {
	playerIDPtr := &playerID
	reviews, _, err := s.reviews.List(ctx, repository.ReviewListOptions{
		PlayerID: playerIDPtr,
		Page:     1,
		PageSize: 1000,
	})
	if err != nil || len(reviews) == 0 {
		return 0.0
	}

	goodCount := 0
	for _, r := range reviews {
		if r.Score >= 4 { // 4分和5分算好评
			goodCount++
		}
	}

	return float32(goodCount) / float32(len(reviews))
}
