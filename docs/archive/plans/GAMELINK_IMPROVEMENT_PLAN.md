# GameLink 系统改进规划 - 精确实施方案

**规划日期**: 2025年11月7日  
**项目阶段**: 未发布阶段 - 可进行大规模改进  
**当前状态**: 管理端完整，用户端和陪玩师端完全缺失

---

## 📋 目录

1. [数据模型改进方案](#1-数据模型改进方案)
2. [后端API新增方案](#2-后端api新增方案)
3. [前端页面实现方案](#3-前端页面实现方案)
4. [系统功能补充方案](#4-系统功能补充方案)
5. [实施时间表](#5-实施时间表)

---

## 1. 数据模型改进方案

### 1.1 需要新增的数据模型

#### 📄 文件: `backend/internal/model/dispute.go`
```go
package model

import "time"

// DisputeStatus 争议状态
type DisputeStatus string

const (
    DisputeStatusPending    DisputeStatus = "pending"     // 待处理
    DisputeStatusInProgress DisputeStatus = "in_progress" // 处理中
    DisputeStatusResolved   DisputeStatus = "resolved"    // 已解决
    DisputeStatusRejected   DisputeStatus = "rejected"    // 已驳回
)

// DisputeType 争议类型
type DisputeType string

const (
    DisputeTypeService DisputeType = "service" // 服务质量
    DisputeTypeRefund  DisputeType = "refund"  // 退款申请
    DisputeTypeOther   DisputeType = "other"   // 其他
)

// Dispute 争议/投诉记录
type Dispute struct {
    Base
    OrderID          uint64        `json:"orderId" gorm:"column:order_id;not null;index"`
    InitiatorID      uint64        `json:"initiatorId" gorm:"column:initiator_id;not null;index"`   // 发起人ID
    InitiatorType    string        `json:"initiatorType" gorm:"column:initiator_type;size:32"`     // 发起人类型 user/player
    RespondentID     uint64        `json:"respondentId" gorm:"column:respondent_id;not null;index"` // 被申诉人ID
    RespondentType   string        `json:"respondentType" gorm:"column:respondent_type;size:32"`   // 被申诉人类型
    Type             DisputeType   `json:"type" gorm:"size:32"`
    Status           DisputeStatus `json:"status" gorm:"size:32;index"`
    Title            string        `json:"title" gorm:"size:255"`
    Description      string        `json:"description" gorm:"type:text"`
    Evidence         string        `json:"evidence,omitempty" gorm:"type:json"`      // 证据（图片URL等）
    HandlerID        *uint64       `json:"handlerId,omitempty" gorm:"column:handler_id;index"` // 处理人ID
    HandlerNote      string        `json:"handlerNote,omitempty" gorm:"column:handler_note;type:text"` // 处理备注
    Resolution       string        `json:"resolution,omitempty" gorm:"type:text"`    // 处理结果
    RefundAmountCents int64        `json:"refundAmountCents,omitempty" gorm:"column:refund_amount_cents"` // 退款金额（分）
    ResolvedAt       *time.Time    `json:"resolvedAt,omitempty" gorm:"column:resolved_at"`
    
    // Relations
    Order      Order  `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:OrderID;references:ID"`
    Initiator  User   `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:InitiatorID;references:ID"`
    Respondent User   `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:RespondentID;references:ID"`
    Handler    *User  `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:HandlerID;references:ID"`
}
```

#### 📄 文件: `backend/internal/model/ticket.go`
```go
package model

import "time"

// TicketStatus 工单状态
type TicketStatus string

const (
    TicketStatusOpen       TicketStatus = "open"       // 待处理
    TicketStatusInProgress TicketStatus = "in_progress" // 处理中
    TicketStatusResolved   TicketStatus = "resolved"    // 已解决
    TicketStatusClosed     TicketStatus = "closed"      // 已关闭
)

// TicketPriority 工单优先级
type TicketPriority string

const (
    TicketPriorityLow      TicketPriority = "low"
    TicketPriorityMedium   TicketPriority = "medium"
    TicketPriorityHigh     TicketPriority = "high"
    TicketPriorityCritical TicketPriority = "critical"
)

// TicketCategory 工单类别
type TicketCategory string

const (
    TicketCategoryAccount  TicketCategory = "account"  // 账号问题
    TicketCategoryPayment  TicketCategory = "payment"  // 支付问题
    TicketCategoryService  TicketCategory = "service"  // 服务问题
    TicketCategoryTechnical TicketCategory = "technical" // 技术问题
    TicketCategoryOther    TicketCategory = "other"    // 其他
)

// Ticket 客服工单
type Ticket struct {
    Base
    TicketNo     string         `json:"ticketNo" gorm:"column:ticket_no;size:64;uniqueIndex"` // 工单号
    UserID       uint64         `json:"userId" gorm:"column:user_id;not null;index"`
    Category     TicketCategory `json:"category" gorm:"size:32"`
    Priority     TicketPriority `json:"priority" gorm:"size:32;default:'medium'"`
    Status       TicketStatus   `json:"status" gorm:"size:32;index;default:'open'"`
    Subject      string         `json:"subject" gorm:"size:255"`
    Description  string         `json:"description" gorm:"type:text"`
    AssignedToID *uint64        `json:"assignedToId,omitempty" gorm:"column:assigned_to_id;index"` // 分配给的客服ID
    ResolvedAt   *time.Time     `json:"resolvedAt,omitempty" gorm:"column:resolved_at"`
    ClosedAt     *time.Time     `json:"closedAt,omitempty" gorm:"column:closed_at"`
    
    // Relations
    User       User            `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:UserID;references:ID"`
    AssignedTo *User           `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:AssignedToID;references:ID"`
    Messages   []TicketMessage `json:"messages,omitempty" gorm:"foreignKey:TicketID"`
}

// TicketMessage 工单消息
type TicketMessage struct {
    Base
    TicketID  uint64 `json:"ticketId" gorm:"column:ticket_id;not null;index"`
    SenderID  uint64 `json:"senderId" gorm:"column:sender_id;not null"`
    Content   string `json:"content" gorm:"type:text"`
    IsStaff   bool   `json:"isStaff" gorm:"column:is_staff;default:false"` // 是否客服回复
    
    // Relations
    Ticket Ticket `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:TicketID;references:ID"`
    Sender User   `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:SenderID;references:ID"`
}
```

#### 📄 文件: `backend/internal/model/notification.go`
```go
package model

import "time"

// NotificationType 通知类型
type NotificationType string

const (
    NotificationTypeSystem  NotificationType = "system"  // 系统通知
    NotificationTypeOrder   NotificationType = "order"   // 订单通知
    NotificationTypePayment NotificationType = "payment" // 支付通知
    NotificationTypeReview  NotificationType = "review"  // 评价通知
    NotificationTypeDispute NotificationType = "dispute" // 争议通知
)

// Notification 站内通知
type Notification struct {
    Base
    UserID       uint64           `json:"userId" gorm:"column:user_id;not null;index"`
    Type         NotificationType `json:"type" gorm:"size:32"`
    Title        string           `json:"title" gorm:"size:255"`
    Content      string           `json:"content" gorm:"type:text"`
    RelatedID    *uint64          `json:"relatedId,omitempty" gorm:"column:related_id"` // 关联对象ID
    RelatedType  string           `json:"relatedType,omitempty" gorm:"column:related_type;size:32"` // 关联对象类型
    IsRead       bool             `json:"isRead" gorm:"column:is_read;default:false;index"`
    ReadAt       *time.Time       `json:"readAt,omitempty" gorm:"column:read_at"`
    
    // Relations
    User User `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID"`
}
```

#### 📄 文件: `backend/internal/model/chat.go`
```go
package model

// ChatMessage 聊天消息（订单内通信）
type ChatMessage struct {
    Base
    OrderID    uint64 `json:"orderId" gorm:"column:order_id;not null;index"`
    SenderID   uint64 `json:"senderId" gorm:"column:sender_id;not null"`
    ReceiverID uint64 `json:"receiverId" gorm:"column:receiver_id;not null"`
    Content    string `json:"content" gorm:"type:text"`
    IsRead     bool   `json:"isRead" gorm:"column:is_read;default:false"`
    MessageType string `json:"messageType,omitempty" gorm:"column:message_type;size:32;default:'text'"` // text/image/file
    
    // Relations
    Order    Order `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:OrderID;references:ID"`
    Sender   User  `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:SenderID;references:ID"`
    Receiver User  `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:ReceiverID;references:ID"`
}
```

#### 📄 文件: `backend/internal/model/favorite.go`
```go
package model

// Favorite 用户收藏
type Favorite struct {
    Base
    UserID   uint64 `json:"userId" gorm:"column:user_id;not null;index:idx_user_player,unique"`
    PlayerID uint64 `json:"playerId" gorm:"column:player_id;not null;index:idx_user_player,unique"`
    
    // Relations
    User   User   `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID"`
    Player Player `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:PlayerID;references:ID"`
}
```

#### 📄 文件: `backend/internal/model/tag.go`
```go
package model

// Tag 标签（用于陪玩师特长标记）
type Tag struct {
    Base
    Name        string `json:"name" gorm:"size:64;uniqueIndex"`
    DisplayName string `json:"displayName" gorm:"column:display_name;size:128"`
    Category    string `json:"category,omitempty" gorm:"size:64"` // 标签分类
    SortOrder   int    `json:"sortOrder" gorm:"column:sort_order;default:0"`
    IsActive    bool   `json:"isActive" gorm:"column:is_active;default:true"`
}

// PlayerTag 陪玩师标签关联
type PlayerTag struct {
    PlayerID uint64 `gorm:"column:player_id;primaryKey"`
    TagID    uint64 `gorm:"column:tag_id;primaryKey"`
    
    // Relations
    Player Player `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:PlayerID;references:ID"`
    Tag    Tag    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:TagID;references:ID"`
}
```

### 1.2 需要修改的现有模型

#### 📄 文件: `backend/internal/model/user.go` - 需要增强
```go
// 在User结构体中添加:
Favorites      []Favorite     `json:"-" gorm:"foreignKey:UserID"`
Notifications  []Notification `json:"-" gorm:"foreignKey:UserID"`
Tickets        []Ticket       `json:"-" gorm:"foreignKey:UserID"`
DisputesInitiated []Dispute   `json:"-" gorm:"foreignKey:InitiatorID"`
DisputesResponded []Dispute   `json:"-" gorm:"foreignKey:RespondentID"`

// 新增字段
Balance        int64          `json:"balance" gorm:"default:0;comment:账户余额（分）"`
FrozenBalance  int64          `json:"frozenBalance" gorm:"column:frozen_balance;default:0;comment:冻结余额（分）"`
RealName       string         `json:"realName,omitempty" gorm:"column:real_name;size:64;comment:实名"`
IDCard         string         `json:"idCard,omitempty" gorm:"column:id_card;size:32;comment:身份证号"`
IsVerified     bool           `json:"isVerified" gorm:"column:is_verified;default:false;comment:是否实名认证"`
VerifiedAt     *time.Time     `json:"verifiedAt,omitempty" gorm:"column:verified_at"`
```

#### 📄 文件: `backend/internal/model/player.go` - 需要增强
```go
// 在Player结构体中添加:
Tags           []Tag          `json:"tags,omitempty" gorm:"many2many:player_tags"`
Favorites      []Favorite     `json:"-" gorm:"foreignKey:PlayerID"`
OnlineStatus   string         `json:"onlineStatus" gorm:"column:online_status;size:32;default:'offline'"` // online/offline/busy
LastOnlineAt   *time.Time     `json:"lastOnlineAt,omitempty" gorm:"column:last_online_at"`
TotalOrders    uint32         `json:"totalOrders" gorm:"column:total_orders;default:0"`
CompletionRate float32        `json:"completionRate" gorm:"column:completion_rate;default:0"` // 完单率
ResponseTime   int            `json:"responseTime" gorm:"column:response_time;default:0;comment:平均响应时间（秒）"`
Specialty      string         `json:"specialty,omitempty" gorm:"type:text;comment:特长描述"`
```

#### 📄 文件: `backend/internal/model/order.go` - 需要增强
```go
// 在Order结构体中添加:
ChatMessages   []ChatMessage  `json:"-" gorm:"foreignKey:OrderID"`
Disputes       []Dispute      `json:"-" gorm:"foreignKey:OrderID"`
UserNote       string         `json:"userNote,omitempty" gorm:"column:user_notes;type:text;comment:用户备注"`
PlayerNote     string         `json:"playerNote,omitempty" gorm:"column:player_notes;type:text;comment:陪玩师备注"`
AcceptedAt     *time.Time     `json:"acceptedAt,omitempty" gorm:"column:accepted_at;comment:接单时间"`
RejectedReason string         `json:"rejectedReason,omitempty" gorm:"column:rejected_reason;type:text"`
```

---

## 2. 后端API新增方案

### 2.1 争议处理系统 API

#### 📄 文件: `backend/internal/handler/user/dispute.go`
```go
package user

import (
    "github.com/gin-gonic/gin"
    "gamelink/internal/service/dispute"
)

type DisputeHandler struct {
    service *dispute.Service
}

func NewDisputeHandler(service *dispute.Service) *DisputeHandler {
    return &DisputeHandler{service: service}
}

// CreateDispute 创建争议
// @Summary 创建争议
// @Tags User-Dispute
// @Accept json
// @Produce json
// @Param request body CreateDisputeRequest true "争议信息"
// @Success 200 {object} DisputeResponse
// @Router /user/disputes [post]
func (h *DisputeHandler) CreateDispute(c *gin.Context) {
    // 实现创建争议逻辑
}

// GetMyDisputes 获取我的争议列表
// @Summary 获取我的争议列表
// @Tags User-Dispute
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} PaginatedDisputeResponse
// @Router /user/disputes [get]
func (h *DisputeHandler) GetMyDisputes(c *gin.Context) {
    // 实现获取争议列表逻辑
}

// GetDisputeDetail 获取争议详情
// @Summary 获取争议详情
// @Tags User-Dispute
// @Produce json
// @Param id path int true "争议ID"
// @Success 200 {object} DisputeDetailResponse
// @Router /user/disputes/:id [get]
func (h *DisputeHandler) GetDisputeDetail(c *gin.Context) {
    // 实现获取争议详情逻辑
}

// WithdrawDispute 撤销争议
// @Summary 撤销争议
// @Tags User-Dispute
// @Param id path int true "争议ID"
// @Success 200 {object} SuccessResponse
// @Router /user/disputes/:id/withdraw [post]
func (h *DisputeHandler) WithdrawDispute(c *gin.Context) {
    // 实现撤销争议逻辑
}
```

#### 📄 文件: `backend/internal/service/dispute/service.go`
```go
package dispute

import (
    "context"
    "gamelink/internal/model"
    "gamelink/internal/repository/dispute"
)

type Service struct {
    repo         dispute.Repository
    orderRepo    order.Repository
    notificationService *notification.Service
}

func NewService(
    repo dispute.Repository,
    orderRepo order.Repository,
    notificationService *notification.Service,
) *Service {
    return &Service{
        repo:         repo,
        orderRepo:    orderRepo,
        notificationService: notificationService,
    }
}

// CreateDispute 创建争议
func (s *Service) CreateDispute(ctx context.Context, req *CreateDisputeRequest) (*model.Dispute, error) {
    // 1. 验证订单存在且属于当前用户
    // 2. 验证订单状态允许创建争议
    // 3. 创建争议记录
    // 4. 发送通知给相关方
    // 5. 记录操作日志
}

// ResolveDispute 解决争议（管理员）
func (s *Service) ResolveDispute(ctx context.Context, disputeID uint64, req *ResolveDisputeRequest) error {
    // 1. 获取争议信息
    // 2. 验证权限
    // 3. 更新争议状态
    // 4. 处理退款（如果有）
    // 5. 发送通知
}
```

#### 📄 文件: `backend/internal/repository/dispute/repository.go`
```go
package dispute

import (
    "context"
    "gorm.io/gorm"
    "gamelink/internal/model"
)

type Repository interface {
    Create(ctx context.Context, dispute *model.Dispute) error
    GetByID(ctx context.Context, id uint64) (*model.Dispute, error)
    List(ctx context.Context, filter *FilterParams) ([]model.Dispute, int64, error)
    Update(ctx context.Context, dispute *model.Dispute) error
    GetByOrderID(ctx context.Context, orderID uint64) ([]model.Dispute, error)
}

type repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
    return &repository{db: db}
}
```

### 2.2 客服工单系统 API

#### 📄 文件: `backend/internal/handler/user/ticket.go`
```go
package user

import (
    "github.com/gin-gonic/gin"
    "gamelink/internal/service/ticket"
)

type TicketHandler struct {
    service *ticket.Service
}

// CreateTicket 创建工单
// GetMyTickets 获取我的工单列表
// GetTicketDetail 获取工单详情
// ReplyTicket 回复工单
// CloseTicket 关闭工单
```

#### 📄 文件: `backend/internal/service/ticket/service.go`
#### 📄 文件: `backend/internal/repository/ticket/repository.go`

### 2.3 通知系统 API

#### 📄 文件: `backend/internal/handler/user/notification.go`
```go
package user

type NotificationHandler struct {
    service *notification.Service
}

// GetMyNotifications 获取我的通知列表
// MarkAsRead 标记为已读
// MarkAllAsRead 全部标记为已读
// DeleteNotification 删除通知
// GetUnreadCount 获取未读数量
```

#### 📄 文件: `backend/internal/service/notification/service.go`
#### 📄 文件: `backend/internal/repository/notification/repository.go`

### 2.4 聊天系统 API (WebSocket)

#### 📄 文件: `backend/internal/handler/websocket/chat.go`
```go
package websocket

import (
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

type ChatHandler struct {
    hub *ChatHub
}

// HandleWebSocket WebSocket连接处理
func (h *ChatHandler) HandleWebSocket(c *gin.Context) {
    // 1. 升级HTTP连接为WebSocket
    // 2. 验证用户身份
    // 3. 加入聊天室
    // 4. 处理消息收发
}

// SendMessage 发送消息
// GetChatHistory 获取聊天历史
// MarkMessagesAsRead 标记消息已读
```

#### 📄 文件: `backend/internal/service/chat/hub.go`
```go
package chat

// ChatHub WebSocket连接管理中心
type ChatHub struct {
    clients    map[uint64]*Client
    broadcast  chan *Message
    register   chan *Client
    unregister chan *Client
}

// Client 客户端连接
type Client struct {
    hub    *ChatHub
    conn   *websocket.Conn
    userID uint64
    send   chan []byte
}
```

#### 📄 文件: `backend/internal/repository/chat/repository.go`

### 2.5 收藏功能 API

#### 📄 文件: `backend/internal/handler/user/favorite.go`
```go
package user

type FavoriteHandler struct {
    service *favorite.Service
}

// AddFavorite 添加收藏
// RemoveFavorite 取消收藏
// GetMyFavorites 获取我的收藏列表
// CheckIsFavorite 检查是否已收藏
```

### 2.6 陪玩师在线状态 API

#### 📄 文件: `backend/internal/handler/player/online.go`
```go
package player

type OnlineHandler struct {
    service *player.Service
}

// UpdateOnlineStatus 更新在线状态
// GetOnlineStatus 获取在线状态
// Heartbeat 心跳保持
```

---

## 3. 前端页面实现方案

### 3.1 用户端页面 (7个核心页面)

#### 📂 目录: `frontend/src/pages/UserPortal/`

#### 📄 文件: `frontend/src/pages/UserPortal/Home/index.tsx`
```typescript
/**
 * 用户首页
 * 功能:
 * - 热门游戏展示
 * - 推荐陪玩师
 * - 最新活动
 * - 快速下单入口
 */
import React from 'react';
import { Link } from 'react-router-dom';
import { GameCard } from '@/components/GameCard';
import { PlayerCard } from '@/components/PlayerCard';

export const UserHomePage: React.FC = () => {
  return (
    <div className="user-home">
      {/* Banner轮播 */}
      <section className="hero-banner">
        <h1>找到你的游戏伙伴</h1>
        <Link to="/games">立即下单</Link>
      </section>
      
      {/* 热门游戏 */}
      <section className="hot-games">
        <h2>热门游戏</h2>
        <div className="game-grid">
          {/* 游戏卡片列表 */}
        </div>
      </section>
      
      {/* 推荐陪玩师 */}
      <section className="recommended-players">
        <h2>推荐陪玩师</h2>
        <div className="player-grid">
          {/* 陪玩师卡片列表 */}
        </div>
      </section>
    </div>
  );
};
```

#### 📄 文件: `frontend/src/pages/UserPortal/GameList/index.tsx`
```typescript
/**
 * 游戏列表页
 * 功能:
 * - 游戏分类筛选
 * - 游戏搜索
 * - 游戏列表展示
 * - 点击进入陪玩师列表
 */
import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { gameApi } from '@/services/api/game';
import { Game } from '@/types/game';

export const GameListPage: React.FC = () => {
  const [games, setGames] = useState<Game[]>([]);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  
  const loadGames = async () => {
    setLoading(true);
    try {
      const data = await gameApi.list();
      setGames(data.games);
    } catch (error) {
      console.error('Failed to load games:', error);
    } finally {
      setLoading(false);
    }
  };
  
  const handleGameClick = (gameId: number) => {
    navigate(`/players?gameId=${gameId}`);
  };
  
  return (
    <div className="game-list-page">
      <h1>选择游戏</h1>
      
      {/* 搜索和筛选 */}
      <div className="filters">
        <input type="text" placeholder="搜索游戏..." />
        <select>
          <option value="">全部分类</option>
          <option value="moba">MOBA</option>
          <option value="fps">FPS</option>
        </select>
      </div>
      
      {/* 游戏网格 */}
      <div className="game-grid">
        {games.map(game => (
          <GameCard 
            key={game.id} 
            game={game}
            onClick={() => handleGameClick(game.id)}
          />
        ))}
      </div>
    </div>
  );
};
```

#### 📄 文件: `frontend/src/pages/UserPortal/PlayerList/index.tsx`
```typescript
/**
 * 陪玩师列表页
 * 功能:
 * - 陪玩师筛选（价格、评分、在线状态等）
 * - 陪玩师排序
 * - 陪玩师卡片展示
 * - 点击查看详情
 * - 快速下单
 */
import React, { useState, useEffect } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import { playerApi } from '@/services/api/player';
import { Player } from '@/types/player';
import { PlayerCard } from '@/components/PlayerCard';

export const PlayerListPage: React.FC = () => {
  const [searchParams] = useSearchParams();
  const gameId = searchParams.get('gameId');
  const [players, setPlayers] = useState<Player[]>([]);
  const [filters, setFilters] = useState({
    minPrice: 0,
    maxPrice: 1000,
    minRating: 0,
    onlineOnly: false,
    sortBy: 'rating', // rating/price/orders
  });
  
  return (
    <div className="player-list-page">
      <h1>选择陪玩师</h1>
      
      {/* 筛选器 */}
      <aside className="filters-sidebar">
        <div className="filter-group">
          <h3>价格范围</h3>
          <input type="range" min="0" max="1000" />
        </div>
        
        <div className="filter-group">
          <h3>评分</h3>
          <label>
            <input type="checkbox" />
            4星及以上
          </label>
        </div>
        
        <div className="filter-group">
          <h3>在线状态</h3>
          <label>
            <input type="checkbox" />
            仅显示在线
          </label>
        </div>
      </aside>
      
      {/* 陪玩师列表 */}
      <main className="player-grid">
        {players.map(player => (
          <PlayerCard key={player.id} player={player} />
        ))}
      </main>
    </div>
  );
};
```

#### 📄 文件: `frontend/src/pages/UserPortal/PlayerDetail/index.tsx`
```typescript
/**
 * 陪玩师详情页
 * 功能:
 * - 陪玩师基本信息展示
 * - 服务项目列表
 * - 评价列表
 * - 收藏按钮
 * - 下单按钮
 */
import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { playerApi } from '@/services/api/player';
import { Player } from '@/types/player';
import { Rating } from '@/components/Rating';

export const PlayerDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [player, setPlayer] = useState<Player | null>(null);
  const [isFavorite, setIsFavorite] = useState(false);
  
  const handleOrder = (serviceId: number) => {
    navigate(`/order/create?playerId=${id}&serviceId=${serviceId}`);
  };
  
  const handleFavorite = async () => {
    // 添加/取消收藏
  };
  
  return (
    <div className="player-detail-page">
      {/* 陪玩师头部信息 */}
      <header className="player-header">
        <img src={player?.avatarUrl} alt={player?.nickname} />
        <div className="player-info">
          <h1>{player?.nickname}</h1>
          <div className="player-stats">
            <Rating value={player?.ratingAverage || 0} />
            <span>已接单 {player?.totalOrders} 次</span>
            <span className={`status ${player?.onlineStatus}`}>
              {player?.onlineStatus === 'online' ? '在线' : '离线'}
            </span>
          </div>
          <button onClick={handleFavorite}>
            {isFavorite ? '已收藏' : '收藏'}
          </button>
        </div>
      </header>
      
      {/* 服务项目 */}
      <section className="services">
        <h2>服务项目</h2>
        <div className="service-list">
          {player?.services?.map(service => (
            <div key={service.id} className="service-item">
              <h3>{service.title}</h3>
              <p>{service.description}</p>
              <span className="price">¥{service.priceYuan}/小时</span>
              <button onClick={() => handleOrder(service.id)}>
                立即下单
              </button>
            </div>
          ))}
        </div>
      </section>
      
      {/* 评价列表 */}
      <section className="reviews">
        <h2>用户评价</h2>
        {/* 评价列表组件 */}
      </section>
    </div>
  );
};
```

#### 📄 文件: `frontend/src/pages/UserPortal/OrderCreate/index.tsx`
```typescript
/**
 * 创建订单页
 * 功能:
 * - 服务信息确认
 * - 时长选择
 * - 特殊要求输入
 * - 价格计算
 * - 提交订单
 */
import React, { useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { orderApi } from '@/services/api/order';
import { Form } from '@/components/Form';

export const OrderCreatePage: React.FC = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    serviceItemId: searchParams.get('serviceId'),
    quantity: 1,
    scheduledStart: new Date(),
    userNotes: '',
  });
  
  const handleSubmit = async () => {
    try {
      const order = await orderApi.create(formData);
      // 跳转到支付页面
      navigate(`/payment?orderId=${order.id}`);
    } catch (error) {
      console.error('Failed to create order:', error);
    }
  };
  
  return (
    <div className="order-create-page">
      <h1>确认订单</h1>
      
      <Form onSubmit={handleSubmit}>
        {/* 服务信息展示 */}
        <section className="service-info">
          <h2>服务信息</h2>
          {/* 显示选择的服务详情 */}
        </section>
        
        {/* 时长选择 */}
        <div className="form-group">
          <label>服务时长</label>
          <select 
            value={formData.quantity}
            onChange={e => setFormData({...formData, quantity: parseInt(e.target.value)})}
          >
            <option value="1">1小时</option>
            <option value="2">2小时</option>
            <option value="3">3小时</option>
          </select>
        </div>
        
        {/* 预约时间 */}
        <div className="form-group">
          <label>预约时间</label>
          <input 
            type="datetime-local"
            value={formData.scheduledStart.toISOString().slice(0, 16)}
            onChange={e => setFormData({...formData, scheduledStart: new Date(e.target.value)})}
          />
        </div>
        
        {/* 特殊要求 */}
        <div className="form-group">
          <label>备注说明</label>
          <textarea 
            value={formData.userNotes}
            onChange={e => setFormData({...formData, userNotes: e.target.value})}
            placeholder="请输入您的特殊要求..."
          />
        </div>
        
        {/* 价格汇总 */}
        <div className="price-summary">
          <div className="row">
            <span>单价</span>
            <span>¥50/小时</span>
          </div>
          <div className="row">
            <span>时长</span>
            <span>{formData.quantity}小时</span>
          </div>
          <div className="row total">
            <span>总计</span>
            <span>¥{50 * formData.quantity}</span>
          </div>
        </div>
        
        <button type="submit">提交订单</button>
      </Form>
    </div>
  );
};
```

#### 📄 文件: `frontend/src/pages/UserPortal/MyOrders/index.tsx`
```typescript
/**
 * 我的订单页
 * 功能:
 * - 订单列表（全部/待支付/进行中/已完成/已取消）
 * - 订单状态筛选
 * - 订单操作（支付、取消、评价、申请退款）
 * - 查看订单详情
 */
import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { orderApi } from '@/services/api/order';
import { Order } from '@/types/order';
import { Tabs } from '@/components/Tabs';

export const MyOrdersPage: React.FC = () => {
  const navigate = useNavigate();
  const [orders, setOrders] = useState<Order[]>([]);
  const [activeTab, setActiveTab] = useState('all');
  
  const tabs = [
    { key: 'all', label: '全部' },
    { key: 'pending', label: '待支付' },
    { key: 'in_progress', label: '进行中' },
    { key: 'completed', label: '已完成' },
    { key: 'canceled', label: '已取消' },
  ];
  
  const handlePay = (orderId: number) => {
    navigate(`/payment?orderId=${orderId}`);
  };
  
  const handleCancel = async (orderId: number) => {
    // 取消订单
  };
  
  const handleReview = (orderId: number) => {
    navigate(`/orders/${orderId}/review`);
  };
  
  const handleDispute = (orderId: number) => {
    navigate(`/disputes/create?orderId=${orderId}`);
  };
  
  return (
    <div className="my-orders-page">
      <h1>我的订单</h1>
      
      <Tabs tabs={tabs} activeTab={activeTab} onChange={setActiveTab} />
      
      <div className="order-list">
        {orders.map(order => (
          <div key={order.id} className="order-card">
            <div className="order-header">
              <span className="order-no">订单号: {order.orderNo}</span>
              <span className={`status ${order.status}`}>
                {getStatusText(order.status)}
              </span>
            </div>
            
            <div className="order-body">
              <img src={order.player?.avatarUrl} alt="" />
              <div className="order-info">
                <h3>{order.player?.nickname}</h3>
                <p>{order.title}</p>
                <p>时长: {order.quantity}小时</p>
              </div>
              <div className="order-price">
                ¥{order.totalPriceCents / 100}
              </div>
            </div>
            
            <div className="order-actions">
              {order.status === 'pending' && (
                <>
                  <button onClick={() => handlePay(order.id)}>立即支付</button>
                  <button onClick={() => handleCancel(order.id)}>取消订单</button>
                </>
              )}
              {order.status === 'completed' && !order.review && (
                <button onClick={() => handleReview(order.id)}>评价</button>
              )}
              {order.status === 'in_progress' && (
                <button onClick={() => handleDispute(order.id)}>申请退款</button>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
```

#### 📄 文件: `frontend/src/pages/UserPortal/Profile/index.tsx`
```typescript
/**
 * 个人中心页
 * 功能:
 * - 用户基本信息展示
 * - 账户余额展示
 * - 我的收藏
 * - 我的工单
 * - 实名认证
 * - 账号设置
 */
import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { userApi } from '@/services/api/user';
import { User } from '@/types/user';

export const ProfilePage: React.FC = () => {
  const [user, setUser] = useState<User | null>(null);
  
  return (
    <div className="profile-page">
      <aside className="profile-sidebar">
        <div className="user-card">
          <img src={user?.avatarUrl} alt={user?.name} />
          <h2>{user?.name}</h2>
          <p>{user?.phone}</p>
        </div>
        
        <nav className="profile-nav">
          <Link to="/profile/info">个人信息</Link>
          <Link to="/profile/balance">账户余额</Link>
          <Link to="/profile/favorites">我的收藏</Link>
          <Link to="/profile/tickets">我的工单</Link>
          <Link to="/profile/verify">实名认证</Link>
          <Link to="/profile/settings">账号设置</Link>
        </nav>
      </aside>
      
      <main className="profile-content">
        {/* 根据路由显示不同的内容 */}
      </main>
    </div>
  );
};
```

### 3.2 陪玩师端页面 (7个核心页面)

#### 📂 目录: `frontend/src/pages/PlayerPortal/`

#### 📄 文件: `frontend/src/pages/PlayerPortal/Dashboard/index.tsx`
```typescript
/**
 * 陪玩师工作台
 * 功能:
 * - 今日数据统计（订单数、收益、评分）
 * - 待接单列表
 * - 进行中订单
 * - 快捷操作入口
 * - 在线状态切换
 */
import React, { useState, useEffect } from 'react';
import { playerApi } from '@/services/api/player';
import { Card } from '@/components/Card';

export const PlayerDashboard: React.FC = () => {
  const [stats, setStats] = useState({
    todayOrders: 0,
    todayEarnings: 0,
    rating: 0,
  });
  const [onlineStatus, setOnlineStatus] = useState('offline');
  
  const toggleOnlineStatus = async () => {
    const newStatus = onlineStatus === 'online' ? 'offline' : 'online';
    await playerApi.updateOnlineStatus(newStatus);
    setOnlineStatus(newStatus);
  };
  
  return (
    <div className="player-dashboard">
      <h1>工作台</h1>
      
      {/* 在线状态切换 */}
      <div className="online-toggle">
        <label className="switch">
          <input 
            type="checkbox" 
            checked={onlineStatus === 'online'}
            onChange={toggleOnlineStatus}
          />
          <span className="slider"></span>
        </label>
        <span>{onlineStatus === 'online' ? '在线接单中' : '离线'}</span>
      </div>
      
      {/* 今日数据 */}
      <div className="stats-grid">
        <Card>
          <h3>今日订单</h3>
          <p className="stat-value">{stats.todayOrders}</p>
        </Card>
        <Card>
          <h3>今日收益</h3>
          <p className="stat-value">¥{stats.todayEarnings / 100}</p>
        </Card>
        <Card>
          <h3>综合评分</h3>
          <p className="stat-value">{stats.rating}</p>
        </Card>
      </div>
      
      {/* 待接单列表 */}
      <section className="pending-orders">
        <h2>待接单</h2>
        {/* 订单卡片列表 */}
      </section>
      
      {/* 进行中订单 */}
      <section className="active-orders">
        <h2>进行中</h2>
        {/* 订单卡片列表 */}
      </section>
    </div>
  );
};
```

#### 📄 文件: `frontend/src/pages/PlayerPortal/Orders/index.tsx`
```typescript
/**
 * 陪玩师订单管理页
 * 功能:
 * - 订单列表（待接单/进行中/已完成/已拒绝）
 * - 接单/拒单操作
 * - 确认开始服务
 * - 确认完成服务
 * - 查看订单详情
 */
import React, { useState } from 'react';
import { orderApi } from '@/services/api/order';
import { Tabs } from '@/components/Tabs';

export const PlayerOrdersPage: React.FC = () => {
  const [activeTab, setActiveTab] = useState('pending');
  
  const tabs = [
    { key: 'pending', label: '待接单' },
    { key: 'accepted', label: '已接单' },
    { key: 'in_progress', label: '进行中' },
    { key: 'completed', label: '已完成' },
    { key: 'rejected', label: '已拒绝' },
  ];
  
  const handleAccept = async (orderId: number) => {
    await orderApi.acceptOrder(orderId);
    // 刷新列表
  };
  
  const handleReject = async (orderId: number, reason: string) => {
    await orderApi.rejectOrder(orderId, { reason });
    // 刷新列表
  };
  
  const handleStart = async (orderId: number) => {
    await orderApi.startService(orderId);
  };
  
  const handleComplete = async (orderId: number) => {
    await orderApi.completeService(orderId);
  };
  
  return (
    <div className="player-orders-page">
      <h1>订单管理</h1>
      <Tabs tabs={tabs} activeTab={activeTab} onChange={setActiveTab} />
      {/* 订单列表 */}
    </div>
  );
};
```

#### 📄 文件: `frontend/src/pages/PlayerPortal/Earnings/index.tsx`
```typescript
/**
 * 收益管理页
 * 功能:
 * - 收益统计（今日/本周/本月）
 * - 收益明细列表
 * - 提现申请
 * - 提现记录
 * - 收益趋势图表
 */
import React, { useState, useEffect } from 'react';
import { earningsApi } from '@/services/api/earnings';
import { Card } from '@/components/Card';
import { Chart } from '@/components/Chart';

export const PlayerEarningsPage: React.FC = () => {
  const [earnings, setEarnings] = useState({
    available: 0,  // 可提现金额
    pending: 0,    // 待结算金额
    total: 0,      // 累计收益
  });
  
  const [withdrawals, setWithdrawals] = useState([]);
  
  const handleWithdraw = () => {
    // 打开提现弹窗
  };
  
  return (
    <div className="player-earnings-page">
      <h1>收益管理</h1>
      
      {/* 收益概览 */}
      <div className="earnings-overview">
        <Card>
          <h3>可提现金额</h3>
          <p className="amount">¥{earnings.available / 100}</p>
          <button onClick={handleWithdraw}>立即提现</button>
        </Card>
        <Card>
          <h3>待结算金额</h3>
          <p className="amount">¥{earnings.pending / 100}</p>
        </Card>
        <Card>
          <h3>累计收益</h3>
          <p className="amount">¥{earnings.total / 100}</p>
        </Card>
      </div>
      
      {/* 收益趋势 */}
      <section className="earnings-chart">
        <h2>收益趋势</h2>
        <Chart type="line" data={[]} />
      </section>
      
      {/* 收益明细 */}
      <section className="earnings-details">
        <h2>收益明细</h2>
        {/* 明细表格 */}
      </section>
      
      {/* 提现记录 */}
      <section className="withdrawal-records">
        <h2>提现记录</h2>
        {/* 提现记录列表 */}
      </section>
    </div>
  );
};
```

#### 📄 文件: `frontend/src/pages/PlayerPortal/Services/index.tsx`
```typescript
/**
 * 服务管理页
 * 功能:
 * - 服务项目列表
 * - 添加服务项目
 * - 编辑服务项目
 * - 删除服务项目
 * - 设置服务价格和时长
 */
import React, { useState, useEffect } from 'react';
import { serviceItemApi } from '@/services/api/serviceItem';

export const PlayerServicesPage: React.FC = () => {
  const [services, setServices] = useState([]);
  const [showModal, setShowModal] = useState(false);
  
  return (
    <div className="player-services-page">
      <h1>服务管理</h1>
      <button onClick={() => setShowModal(true)}>添加服务</button>
      {/* 服务列表 */}
    </div>
  );
};
```

#### 📄 文件: `frontend/src/pages/PlayerPortal/Profile/index.tsx`
```typescript
/**
 * 陪玩师资料管理页
 * 功能:
 * - 基本信息编辑
 * - 头像上传
 * - 认证资料管理
 * - 特长标签设置
 * - 个人简介编辑
 */
```

#### 📄 文件: `frontend/src/pages/PlayerPortal/Reviews/index.tsx`
```typescript
/**
 * 评价管理页
 * 功能:
 * - 收到的评价列表
 * - 评价统计
 * - 回复评价
 */
```

#### 📄 文件: `frontend/src/pages/PlayerPortal/Schedule/index.tsx`
```typescript
/**
 * 时间管理页
 * 功能:
 * - 可接单时间设置
 * - 休息时间设置
 * - 日历视图
 */
```

### 3.3 通用组件新增

#### 📄 文件: `frontend/src/components/GameCard/index.tsx`
```typescript
/**
 * 游戏卡片组件
 * 用于展示游戏信息
 */
import React from 'react';
import { Game } from '@/types/game';
import styles from './GameCard.module.less';

interface GameCardProps {
  game: Game;
  onClick?: () => void;
}

export const GameCard: React.FC<GameCardProps> = ({ game, onClick }) => {
  return (
    <div className={styles.gameCard} onClick={onClick}>
      <img src={game.iconUrl} alt={game.name} />
      <h3>{game.name}</h3>
      <p>{game.playerCount} 陪玩师在线</p>
    </div>
  );
};
```

#### 📄 文件: `frontend/src/components/PlayerCard/index.tsx`
```typescript
/**
 * 陪玩师卡片组件
 * 用于展示陪玩师信息
 */
import React from 'react';
import { Player } from '@/types/player';
import { Rating } from '@/components/Rating';
import styles from './PlayerCard.module.less';

interface PlayerCardProps {
  player: Player;
  onClick?: () => void;
}

export const PlayerCard: React.FC<PlayerCardProps> = ({ player, onClick }) => {
  return (
    <div className={styles.playerCard} onClick={onClick}>
      <div className="player-avatar">
        <img src={player.avatarUrl} alt={player.nickname} />
        <span className={`status ${player.onlineStatus}`}></span>
      </div>
      <div className="player-info">
        <h3>{player.nickname}</h3>
        <Rating value={player.ratingAverage} />
        <p className="price">¥{player.hourlyRateCents / 100}/小时</p>
        <span className="orders">{player.totalOrders}单</span>
      </div>
    </div>
  );
};
```

#### 📄 文件: `frontend/src/components/OrderStatusBadge/index.tsx`
```typescript
/**
 * 订单状态徽章组件
 */
```

#### 📄 文件: `frontend/src/components/ChatWindow/index.tsx`
```typescript
/**
 * 聊天窗口组件
 * 用于订单内用户和陪玩师沟通
 */
```

### 3.4 前端服务层新增

#### 📄 文件: `frontend/src/services/api/dispute.ts`
```typescript
import { apiClient } from '../client';
import { Dispute } from '@/types/dispute';

export const disputeApi = {
  create: (data: CreateDisputeRequest) => 
    apiClient.post<Dispute>('/user/disputes', data),
  
  list: (params?: ListParams) => 
    apiClient.get<PaginatedResponse<Dispute>>('/user/disputes', { params }),
  
  getById: (id: number) => 
    apiClient.get<Dispute>(`/user/disputes/${id}`),
  
  withdraw: (id: number) => 
    apiClient.post(`/user/disputes/${id}/withdraw`),
};
```

#### 📄 文件: `frontend/src/services/api/ticket.ts`
```typescript
export const ticketApi = {
  create: (data: CreateTicketRequest) => 
    apiClient.post<Ticket>('/user/tickets', data),
  
  list: (params?: ListParams) => 
    apiClient.get<PaginatedResponse<Ticket>>('/user/tickets', { params }),
  
  getById: (id: number) => 
    apiClient.get<Ticket>(`/user/tickets/${id}`),
  
  reply: (id: number, content: string) => 
    apiClient.post(`/user/tickets/${id}/messages`, { content }),
  
  close: (id: number) => 
    apiClient.post(`/user/tickets/${id}/close`),
};
```

#### 📄 文件: `frontend/src/services/api/notification.ts`
```typescript
export const notificationApi = {
  list: (params?: ListParams) => 
    apiClient.get<PaginatedResponse<Notification>>('/user/notifications', { params }),
  
  getUnreadCount: () => 
    apiClient.get<{ count: number }>('/user/notifications/unread-count'),
  
  markAsRead: (id: number) => 
    apiClient.post(`/user/notifications/${id}/read`),
  
  markAllAsRead: () => 
    apiClient.post('/user/notifications/read-all'),
  
  delete: (id: number) => 
    apiClient.delete(`/user/notifications/${id}`),
};
```

#### 📄 文件: `frontend/src/services/api/favorite.ts`
```typescript
export const favoriteApi = {
  add: (playerId: number) => 
    apiClient.post('/user/favorites', { playerId }),
  
  remove: (playerId: number) => 
    apiClient.delete(`/user/favorites/${playerId}`),
  
  list: (params?: ListParams) => 
    apiClient.get<PaginatedResponse<Favorite>>('/user/favorites', { params }),
  
  check: (playerId: number) => 
    apiClient.get<{ isFavorite: boolean }>(`/user/favorites/check/${playerId}`),
};
```

#### 📄 文件: `frontend/src/services/websocket/chat.ts`
```typescript
/**
 * WebSocket聊天服务
 */
class ChatService {
  private ws: WebSocket | null = null;
  private listeners: Map<string, Function[]> = new Map();
  
  connect(orderId: number) {
    const token = localStorage.getItem('token');
    this.ws = new WebSocket(`ws://localhost:8080/ws/chat/${orderId}?token=${token}`);
    
    this.ws.onopen = () => {
      console.log('WebSocket connected');
      this.emit('connected');
    };
    
    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      this.emit('message', data);
    };
    
    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      this.emit('error', error);
    };
    
    this.ws.onclose = () => {
      console.log('WebSocket disconnected');
      this.emit('disconnected');
    };
  }
  
  sendMessage(content: string) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({ type: 'message', content }));
    }
  }
  
  on(event: string, callback: Function) {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, []);
    }
    this.listeners.get(event)!.push(callback);
  }
  
  off(event: string, callback: Function) {
    const callbacks = this.listeners.get(event);
    if (callbacks) {
      const index = callbacks.indexOf(callback);
      if (index > -1) {
        callbacks.splice(index, 1);
      }
    }
  }
  
  private emit(event: string, data?: any) {
    const callbacks = this.listeners.get(event);
    if (callbacks) {
      callbacks.forEach(callback => callback(data));
    }
  }
  
  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }
}

export const chatService = new ChatService();
```

### 3.5 前端类型定义新增

#### 📄 文件: `frontend/src/types/dispute.ts`
```typescript
export interface Dispute {
  id: number;
  orderId: number;
  initiatorId: number;
  initiatorType: 'user' | 'player';
  respondentId: number;
  respondentType: 'user' | 'player';
  type: 'service' | 'refund' | 'other';
  status: 'pending' | 'in_progress' | 'resolved' | 'rejected';
  title: string;
  description: string;
  evidence?: string[];
  handlerId?: number;
  handlerNote?: string;
  resolution?: string;
  refundAmountCents?: number;
  resolvedAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateDisputeRequest {
  orderId: number;
  type: string;
  title: string;
  description: string;
  evidence?: string[];
}
```

#### 📄 文件: `frontend/src/types/ticket.ts`
```typescript
export interface Ticket {
  id: number;
  ticketNo: string;
  userId: number;
  category: 'account' | 'payment' | 'service' | 'technical' | 'other';
  priority: 'low' | 'medium' | 'high' | 'critical';
  status: 'open' | 'in_progress' | 'resolved' | 'closed';
  subject: string;
  description: string;
  assignedToId?: number;
  messages?: TicketMessage[];
  resolvedAt?: string;
  closedAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface TicketMessage {
  id: number;
  ticketId: number;
  senderId: number;
  content: string;
  isStaff: boolean;
  createdAt: string;
}
```

#### 📄 文件: `frontend/src/types/notification.ts`
```typescript
export interface Notification {
  id: number;
  userId: number;
  type: 'system' | 'order' | 'payment' | 'review' | 'dispute';
  title: string;
  content: string;
  relatedId?: number;
  relatedType?: string;
  isRead: boolean;
  readAt?: string;
  createdAt: string;
}
```

#### 📄 文件: `frontend/src/types/favorite.ts`
```typescript
export interface Favorite {
  id: number;
  userId: number;
  playerId: number;
  player?: Player;
  createdAt: string;
}
```

#### 📄 文件: `frontend/src/types/chat.ts`
```typescript
export interface ChatMessage {
  id: number;
  orderId: number;
  senderId: number;
  receiverId: number;
  content: string;
  messageType: 'text' | 'image' | 'file';
  isRead: boolean;
  createdAt: string;
}
```

---

## 4. 系统功能补充方案

### 4.1 支付系统改进

#### 📄 文件: `backend/internal/service/payment/alipay.go`
```go
package payment

import (
    "github.com/smartwalle/alipay/v3"
)

// AlipayService 支付宝支付服务
type AlipayService struct {
    client *alipay.Client
}

// CreatePayment 创建支付宝支付
func (s *AlipayService) CreatePayment(order *model.Order) (*PaymentResponse, error) {
    // 1. 生成支付参数
    // 2. 调用支付宝API
    // 3. 返回支付URL
}

// HandleCallback 处理支付宝回调
func (s *AlipayService) HandleCallback(params map[string]string) error {
    // 1. 验证签名
    // 2. 更新订单状态
    // 3. 更新支付记录
    // 4. 发送通知
}

// Refund 退款
func (s *AlipayService) Refund(payment *model.Payment, amount int64) error {
    // 1. 调用支付宝退款API
    // 2. 更新退款记录
}
```

#### 📄 文件: `backend/internal/service/payment/wechat.go`
```go
package payment

// WeChatPayService 微信支付服务
type WeChatPayService struct {
    // ...
}
```

### 4.2 文件上传系统

#### 📄 文件: `backend/internal/handler/upload/upload.go`
```go
package upload

import (
    "github.com/gin-gonic/gin"
)

type UploadHandler struct {
    service *upload.Service
}

// UploadImage 上传图片
// @Summary 上传图片
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "图片文件"
// @Success 200 {object} UploadResponse
// @Router /upload/image [post]
func (h *UploadHandler) UploadImage(c *gin.Context) {
    // 1. 获取上传文件
    // 2. 验证文件类型和大小
    // 3. 生成文件名
    // 4. 保存文件（本地/OSS）
    // 5. 返回文件URL
}

// UploadFile 上传文件
func (h *UploadHandler) UploadFile(c *gin.Context) {}
```

#### 📄 文件: `backend/internal/service/upload/service.go`
```go
package upload

type Service struct {
    storage Storage
}

type Storage interface {
    Save(file *multipart.FileHeader) (string, error)
    Delete(url string) error
}

// LocalStorage 本地存储
type LocalStorage struct {
    basePath string
}

// OSSStorage 阿里云OSS存储
type OSSStorage struct {
    client *oss.Client
    bucket string
}
```

### 4.3 实时通知系统

#### 📄 文件: `backend/internal/handler/websocket/notification.go`
```go
package websocket

// NotificationHub 通知推送中心
type NotificationHub struct {
    clients    map[uint64]*Client
    broadcast  chan *Notification
    register   chan *Client
    unregister chan *Client
}

func (h *NotificationHub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client.userID] = client
            
        case client := <-h.unregister:
            delete(h.clients, client.userID)
            close(client.send)
            
        case notification := <-h.broadcast:
            if client, ok := h.clients[notification.UserID]; ok {
                select {
                case client.send <- notification.ToJSON():
                default:
                    close(client.send)
                    delete(h.clients, client.userID)
                }
            }
        }
    }
}
```

### 4.4 定时任务系统

#### 📄 文件: `backend/internal/scheduler/order_scheduler.go`
```go
package scheduler

import (
    "github.com/robfig/cron/v3"
)

type OrderScheduler struct {
    cron          *cron.Cron
    orderService  *order.Service
}

func (s *OrderScheduler) Start() {
    // 每5分钟检查订单超时
    s.cron.AddFunc("*/5 * * * *", s.checkOrderTimeout)
    
    // 每小时检查服务完成
    s.cron.AddFunc("0 * * * *", s.checkServiceCompletion)
    
    // 每天凌晨2点结算收益
    s.cron.AddFunc("0 2 * * *", s.settleEarnings)
    
    s.cron.Start()
}

func (s *OrderScheduler) checkOrderTimeout() {
    // 检查待支付订单超时
    // 自动取消超时订单
}

func (s *OrderScheduler) checkServiceCompletion() {
    // 检查服务是否按时完成
    // 发送提醒通知
}

func (s *OrderScheduler) settleEarnings() {
    // 结算陪玩师收益
    // 更新可提现金额
}
```

### 4.5 监控和日志系统

#### 📄 文件: `backend/internal/middleware/prometheus.go`
```go
package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus"
)

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
        },
        []string{"method", "path"},
    )
)

func PrometheusMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(c.Request.Method, c.Request.URL.Path))
        defer timer.ObserveDuration()
        
        c.Next()
        
        httpRequestsTotal.WithLabelValues(
            c.Request.Method,
            c.Request.URL.Path,
            strconv.Itoa(c.Writer.Status()),
        ).Inc()
    }
}
```

---

## 5. 实施时间表

### 第一周: 数据模型和核心API (2024.11.11 - 2024.11.17)

#### Day 1-2: 数据模型实现
- [ ] 创建争议模型 (dispute.go)
- [ ] 创建工单模型 (ticket.go)
- [ ] 创建通知模型 (notification.go)
- [ ] 创建聊天模型 (chat.go)
- [ ] 创建收藏模型 (favorite.go)
- [ ] 创建标签模型 (tag.go)
- [ ] 修改现有模型 (user.go, player.go, order.go)
- [ ] 运行数据库迁移
- [ ] 编写单元测试

#### Day 3-4: Repository层实现
- [ ] 实现DisputeRepository
- [ ] 实现TicketRepository
- [ ] 实现NotificationRepository
- [ ] 实现ChatRepository
- [ ] 实现FavoriteRepository
- [ ] 实现TagRepository
- [ ] 编写Repository测试

#### Day 5-7: Service层实现
- [ ] 实现DisputeService
- [ ] 实现TicketService
- [ ] 实现NotificationService
- [ ] 实现ChatService
- [ ] 实现FavoriteService
- [ ] 支付服务改进（支付宝/微信真实集成）
- [ ] 编写Service测试

### 第二周: 后端API完成 (2024.11.18 - 2024.11.24)

#### Day 8-10: Handler层实现
- [ ] 实现争议处理Handler
- [ ] 实现工单Handler
- [ ] 实现通知Handler
- [ ] 实现收藏Handler
- [ ] 实现文件上传Handler
- [ ] 实现WebSocket Chat Handler
- [ ] 编写Handler测试

#### Day 11-12: 支付系统改进
- [ ] 实现支付宝支付集成
- [ ] 实现微信支付集成
- [ ] 实现支付回调处理
- [ ] 实现退款流程
- [ ] 编写支付测试

#### Day 13-14: 系统功能补充
- [ ] 实现定时任务系统
- [ ] 实现文件上传服务
- [ ] 实现监控中间件
- [ ] API文档更新
- [ ] 集成测试

### 第三周: 用户端前端开发 (2024.11.25 - 2024.12.1)

#### Day 15-16: 用户端基础页面
- [ ] 用户首页 (Home)
- [ ] 游戏列表页 (GameList)
- [ ] 陪玩师列表页 (PlayerList)
- [ ] API集成和测试

#### Day 17-18: 用户端订单页面
- [ ] 陪玩师详情页 (PlayerDetail)
- [ ] 创建订单页 (OrderCreate)
- [ ] 支付页面
- [ ] API集成和测试

#### Day 19-20: 用户端个人中心
- [ ] 我的订单页 (MyOrders)
- [ ] 个人中心页 (Profile)
- [ ] 收藏页
- [ ] 工单页
- [ ] API集成和测试

#### Day 21: 用户端优化
- [ ] 响应式适配
- [ ] 交互优化
- [ ] 错误处理
- [ ] 测试和修复

### 第四周: 陪玩师端前端开发 (2024.12.2 - 2024.12.8)

#### Day 22-23: 陪玩师端核心页面
- [ ] 陪玩师工作台 (Dashboard)
- [ ] 订单管理页 (Orders)
- [ ] API集成和测试

#### Day 24-25: 陪玩师端收益管理
- [ ] 收益管理页 (Earnings)
- [ ] 提现管理
- [ ] 收益统计图表
- [ ] API集成和测试

#### Day 26-27: 陪玩师端资料管理
- [ ] 服务管理页 (Services)
- [ ] 资料管理页 (Profile)
- [ ] 评价管理页 (Reviews)
- [ ] 时间管理页 (Schedule)
- [ ] API集成和测试

#### Day 28: 陪玩师端优化
- [ ] 响应式适配
- [ ] 交互优化
- [ ] 测试和修复

### 第五周: 通用功能和组件 (2024.12.9 - 2024.12.15)

#### Day 29-30: 通用组件开发
- [ ] GameCard组件
- [ ] PlayerCard组件
- [ ] OrderStatusBadge组件
- [ ] ChatWindow组件
- [ ] 组件测试

#### Day 31-32: WebSocket集成
- [ ] 聊天功能实现
- [ ] 实时通知功能
- [ ] 在线状态管理
- [ ] 测试和优化

#### Day 33-34: 争议和工单系统
- [ ] 争议创建页面
- [ ] 争议详情页面
- [ ] 工单创建页面
- [ ] 工单详情页面
- [ ] 测试和优化

#### Day 35: 服务层完善
- [ ] 完善API服务层
- [ ] 添加类型定义
- [ ] 错误处理优化
- [ ] 测试覆盖

### 第六周: 测试和优化 (2024.12.16 - 2024.12.22)

#### Day 36-37: 后端测试
- [ ] 单元测试补充（目标80%覆盖率）
- [ ] 集成测试
- [ ] API测试
- [ ] 性能测试

#### Day 38-39: 前端测试
- [ ] 组件测试补充
- [ ] 页面测试
- [ ] E2E测试
- [ ] 性能优化

#### Day 40-41: 系统集成测试
- [ ] 完整业务流程测试
- [ ] 支付流程测试
- [ ] 争议处理流程测试
- [ ] 工单流程测试

#### Day 42: 文档和部署准备
- [ ] API文档完善
- [ ] 用户手册编写
- [ ] 部署文档编写
- [ ] 代码审查和优化

---

## 6. 关键技术决策

### 6.1 数据库设计决策

#### 软删除策略
- 所有核心业务表使用软删除 (DeletedAt)
- 保留历史数据用于数据分析和审计

#### 金额存储
- 统一使用int64存储分为单位的金额
- 避免浮点数精度问题

#### 索引策略
- 外键字段添加索引
- 状态字段添加索引
- 时间字段根据查询需求添加索引
- 复合索引用于高频组合查询

### 6.2 支付系统决策

#### 支付方式
- 支持支付宝和微信支付
- 预留其他支付方式扩展接口

#### 支付安全
- 所有回调必须验证签名
- 支付金额双重验证
- 支付状态机严格控制

#### 退款策略
- 支持部分退款和全额退款
- 退款需要管理员审核
- 退款记录完整保留

### 6.3 实时通信决策

#### WebSocket使用场景
- 订单内聊天
- 实时通知推送
- 在线状态更新

#### 消息存储
- 聊天消息存储在数据库
- 消息有效期设置（可选）
- 支持消息撤回

### 6.4 文件存储决策

#### 存储方式
- 开发环境：本地存储
- 生产环境：阿里云OSS/腾讯云COS

#### 文件类型
- 图片：头像、证据、服务图片
- 文件：身份证、认证资料

#### 安全控制
- 文件大小限制
- 文件类型白名单
- 文件扫描（病毒、敏感内容）

### 6.5 性能优化策略

#### 缓存策略
- Redis缓存热点数据
- 游戏列表缓存
- 陪玩师列表缓存
- 用户Session缓存

#### 数据库优化
- 使用连接池
- 慢查询监控
- 索引优化
- 读写分离（后期）

#### 前端优化
- 路由懒加载
- 图片懒加载
- 组件缓存
- API请求防抖

---

## 7. 风险评估和应对

### 7.1 技术风险

#### 风险: 支付集成复杂度高
- **影响**: 可能延期
- **概率**: 中
- **应对**: 
  - 提前阅读支付接口文档
  - 准备测试环境
  - 预留充足测试时间

#### 风险: WebSocket稳定性
- **影响**: 实时功能不稳定
- **概率**: 中
- **应对**:
  - 实现自动重连机制
  - 添加心跳检测
  - 降级方案（轮询）

#### 风险: 数据迁移问题
- **影响**: 数据丢失
- **概率**: 低
- **应对**:
  - 数据备份
  - 迁移脚本测试
  - 回滚方案

### 7.2 业务风险

#### 风险: 争议处理流程不完善
- **影响**: 用户体验差
- **概率**: 中
- **应对**:
  - 详细的流程设计
  - 多方测试
  - 快速迭代

#### 风险: 支付安全问题
- **影响**: 资金损失
- **概率**: 低
- **应对**:
  - 严格的安全审查
  - 多层验证
  - 监控告警

### 7.3 时间风险

#### 风险: 开发进度延期
- **影响**: 上线延迟
- **概率**: 中
- **应对**:
  - 功能优先级排序
  - MVP优先
  - 并行开发

---

## 8. 质量保证计划

### 8.1 代码质量

- [ ] 代码审查流程
- [ ] 代码规范检查 (golangci-lint, ESLint)
- [ ] 单元测试覆盖率 >= 80%
- [ ] 集成测试覆盖核心流程

### 8.2 功能测试

- [ ] 用户端完整流程测试
- [ ] 陪玩师端完整流程测试
- [ ] 支付流程测试
- [ ] 争议处理流程测试
- [ ] 工单流程测试

### 8.3 性能测试

- [ ] 接口响应时间 < 200ms
- [ ] 并发用户测试 (1000+)
- [ ] 数据库查询优化
- [ ] 前端加载性能优化

### 8.4 安全测试

- [ ] SQL注入测试
- [ ] XSS攻击测试
- [ ] CSRF攻击测试
- [ ] 支付安全测试
- [ ] 权限控制测试

---

## 9. 部署和发布计划

### 9.1 环境准备

#### 开发环境
- [ ] 本地开发环境配置
- [ ] Docker容器化

#### 测试环境
- [ ] 测试服务器部署
- [ ] 数据库配置
- [ ] Redis配置
- [ ] OSS配置

#### 生产环境
- [ ] 生产服务器配置
- [ ] 数据库主从配置
- [ ] Redis集群配置
- [ ] CDN配置
- [ ] 监控系统配置

### 9.2 发布策略

#### 灰度发布
- [ ] 10% 用户灰度测试
- [ ] 监控关键指标
- [ ] 逐步扩大范围

#### 回滚方案
- [ ] 数据库备份
- [ ] 代码版本管理
- [ ] 快速回滚脚本

---

## 10. 后续优化计划

### 10.1 短期优化 (1-3个月)

- [ ] 用户行为分析系统
- [ ] 推荐算法优化
- [ ] 移动端APP开发
- [ ] 消息推送系统（短信/邮件）

### 10.2 中期优化 (3-6个月)

- [ ] 微服务架构改造
- [ ] 数据中台建设
- [ ] AI智能客服
- [ ] 营销活动系统

### 10.3 长期规划 (6-12个月)

- [ ] 国际化支持
- [ ] 多业务线扩展
- [ ] 云原生部署
- [ ] 大数据分析平台

---

## 附录

### A. 技术栈清单

#### 后端
- Go 1.25+
- Gin框架
- GORM
- Redis
- PostgreSQL/SQLite
- WebSocket
- Prometheus

#### 前端
- React 18
- TypeScript 5.6
- Vite 5.4
- Less
- Axios
- WebSocket

#### 工具
- Docker
- Git
- Swagger
- Postman

### B. 参考文档

- [Go编码规范](./docs/api/go-coding-standards.md)
- [后端项目指南](./backend/PROJECT_GUIDELINES.md)
- [前端开发规范](./frontend/README.md)
- [API文档](./docs/api/)

---

**文档版本**: v1.0  
**最后更新**: 2025年11月7日  
**维护人**: 开发团队  
**审核状态**: ✅ 已审核

