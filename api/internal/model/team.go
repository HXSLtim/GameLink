package model

import "time"

// TeamStatus 团队状态
type TeamStatus string

const (
	TeamStatusActive   TeamStatus = "active"   // 正常
	TeamStatusBusy     TeamStatus = "busy"     // 忙碌（接单中）
	TeamStatusInactive TeamStatus = "inactive" // 停用
)

// TeamMemberRole 团队成员角色
type TeamMemberRole string

const (
	TeamMemberRoleLeader TeamMemberRole = "leader" // 队长
	TeamMemberRoleMember TeamMemberRole = "member" // 队员
)

// TeamMemberStatus 团队成员状态
type TeamMemberStatus string

const (
	TeamMemberStatusActive TeamMemberStatus = "active" // 正常
	TeamMemberStatusLeft   TeamMemberStatus = "left"   // 已离队
	TeamMemberStatusKicked TeamMemberStatus = "kicked" // 被踢出
)

// Team 陪玩师团队
type Team struct {
	Base
	Name        string     `json:"name" gorm:"size:64;not null"`                    // 团队名称
	Description string     `json:"description" gorm:"size:255"`                     // 团队简介
	AvatarURL   string     `json:"avatarUrl" gorm:"column:avatar_url;size:255"`     // 团队头像
	LeaderID    uint64     `json:"leaderId" gorm:"column:leader_id;not null;index"` // 队长ID（Player.ID）
	Status      TeamStatus `json:"status" gorm:"size:32;default:'active';index"`    // 状态

	// 配置
	MaxMembers  int `json:"maxMembers" gorm:"column:max_members;default:5"`   // 最大成员数
	MemberCount int `json:"memberCount" gorm:"column:member_count;default:1"` // 当前成员数

	// 收入分配配置
	IncomeShareType string  `json:"incomeShareType" gorm:"column:income_share_type;size:32;default:'equal'"` // equal=平分, custom=自定义
	LeaderBonusRate float64 `json:"leaderBonusRate" gorm:"column:leader_bonus_rate;default:0"`               // 队长额外抽成比例（预留）

	// 统计
	TotalOrderCount  int   `json:"totalOrderCount" gorm:"column:total_order_count;default:0"`   // 累计接单数
	TotalIncomeCents int64 `json:"totalIncomeCents" gorm:"column:total_income_cents;default:0"` // 累计收入（分）

	// 当前订单（忙碌状态时）
	CurrentOrderID *uint64 `json:"currentOrderId,omitempty" gorm:"column:current_order_id;index"` // 当前进行中的订单ID

	// Relations
	Leader  *Player      `json:"leader,omitempty" gorm:"foreignKey:LeaderID"`
	Members []TeamMember `json:"members,omitempty" gorm:"foreignKey:TeamID"`
}

// TableName 指定表名
func (Team) TableName() string {
	return "teams"
}

// IsBusy 是否忙碌
func (t *Team) IsBusy() bool {
	return t.Status == TeamStatusBusy || t.CurrentOrderID != nil
}

// CanAcceptOrder 是否可以接单
func (t *Team) CanAcceptOrder() bool {
	return t.Status == TeamStatusActive && t.CurrentOrderID == nil
}

// TeamMember 团队成员
type TeamMember struct {
	Base
	TeamID   uint64           `json:"teamId" gorm:"column:team_id;not null;index;uniqueIndex:idx_team_player"`     // 团队ID
	PlayerID uint64           `json:"playerId" gorm:"column:player_id;not null;index;uniqueIndex:idx_team_player"` // 陪玩师ID
	Role     TeamMemberRole   `json:"role" gorm:"size:32;default:'member'"`                                        // 角色
	Status   TeamMemberStatus `json:"status" gorm:"size:32;default:'active';index"`                                // 状态
	JoinedAt time.Time        `json:"joinedAt" gorm:"column:joined_at;not null"`                                   // 加入时间
	LeftAt   *time.Time       `json:"leftAt,omitempty" gorm:"column:left_at"`                                      // 离队时间

	// 顺位（用于队长继承）
	SortOrder int `json:"sortOrder" gorm:"column:sort_order;default:0"` // 排序（越小越靠前，队长继承顺序）

	// 统计
	OrderCount  int   `json:"orderCount" gorm:"column:order_count;default:0"`   // 在团队内接单数
	IncomeCents int64 `json:"incomeCents" gorm:"column:income_cents;default:0"` // 在团队内收入（分）

	// Relations
	Team   *Team   `json:"team,omitempty" gorm:"foreignKey:TeamID"`
	Player *Player `json:"player,omitempty" gorm:"foreignKey:PlayerID"`
}

// TableName 指定表名
func (TeamMember) TableName() string {
	return "team_members"
}

// IsLeader 是否是队长
func (tm *TeamMember) IsLeader() bool {
	return tm.Role == TeamMemberRoleLeader
}

// IsActive 是否在队
func (tm *TeamMember) IsActive() bool {
	return tm.Status == TeamMemberStatusActive
}

// TeamInviteStatus 团队邀请状态
type TeamInviteStatus string

const (
	TeamInviteStatusPending  TeamInviteStatus = "pending"  // 待处理
	TeamInviteStatusAccepted TeamInviteStatus = "accepted" // 已接受
	TeamInviteStatusRejected TeamInviteStatus = "rejected" // 已拒绝
	TeamInviteStatusExpired  TeamInviteStatus = "expired"  // 已过期
)

// TeamInvite 团队邀请（预留）
type TeamInvite struct {
	Base
	TeamID    uint64           `json:"teamId" gorm:"column:team_id;not null;index"`       // 团队ID
	PlayerID  uint64           `json:"playerId" gorm:"column:player_id;not null;index"`   // 被邀请的陪玩师ID
	InviterID uint64           `json:"inviterId" gorm:"column:inviter_id;not null;index"` // 邀请人ID
	Status    TeamInviteStatus `json:"status" gorm:"size:32;default:'pending';index"`     // 状态
	ExpireAt  time.Time        `json:"expireAt" gorm:"column:expire_at;not null;index"`   // 过期时间
	Message   string           `json:"message" gorm:"size:255"`                           // 邀请留言

	// Relations
	Team    *Team   `json:"team,omitempty" gorm:"foreignKey:TeamID"`
	Player  *Player `json:"player,omitempty" gorm:"foreignKey:PlayerID"`
	Inviter *Player `json:"inviter,omitempty" gorm:"foreignKey:InviterID"`
}

// TableName 指定表名
func (TeamInvite) TableName() string {
	return "team_invites"
}
