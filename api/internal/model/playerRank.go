package model

import "time"

// ========== 游戏段位配置 ==========

// GameRank 游戏段位配置
// @Description 平台自定义的游戏段位，每个游戏可以有多个段位
type GameRank struct {
	Base
	// 游戏ID
	GameID uint64 `json:"gameId" gorm:"column:game_id;index;not null;index:idx_game_level,priority:1"`
	// 段位名称（如：青铜、白银、黄金、铂金、钻石、大师、王者）
	Name string `json:"name" gorm:"size:64;not null"`
	// 段位等级（数字，用于排序和比较，越大越高）
	Level int `json:"level" gorm:"default:0;index:idx_game_level,priority:2"`
	// 该段位定价（分）
	PriceCents int64 `json:"priceCents" gorm:"column:price_cents;default:0"`
	// 段位图标URL
	IconURL string `json:"iconUrl,omitempty" gorm:"column:icon_url;size:255"`
	// 段位颜色（前端展示用）
	Color string `json:"color,omitempty" gorm:"size:32"`
	// 段位描述
	Description string `json:"description,omitempty" gorm:"type:text"`
	// 排序（越小越靠前）
	SortOrder int `json:"sortOrder" gorm:"column:sort_order;default:0"`
	// 是否启用
	IsActive bool `json:"isActive" gorm:"column:is_active;default:true;index"`

	// 关联
	Game *Game `json:"game,omitempty" gorm:"foreignKey:GameID"`
}

// TableName 指定表名
func (GameRank) TableName() string {
	return "game_ranks"
}

// ========== 陪玩师段位认证 ==========

// PlayerRankStatus 陪玩师段位认证状态
type PlayerRankStatus string

const (
	PlayerRankStatusPending  PlayerRankStatus = "pending"  // 待审核
	PlayerRankStatusVerified PlayerRankStatus = "verified" // 已认证
	PlayerRankStatusRejected PlayerRankStatus = "rejected" // 已拒绝
	PlayerRankStatusRevoked  PlayerRankStatus = "revoked"  // 已撤销（降级/过期）
	PlayerRankStatusExpired  PlayerRankStatus = "expired"  // 已过期
)

// PlayerRankRecord 陪玩师段位认证记录
// @Description 陪玩师的游戏段位认证，一个陪玩师可以有多个游戏的多个段位
type PlayerRankRecord struct {
	Base
	// 陪玩师ID
	PlayerID uint64 `json:"playerId" gorm:"column:player_id;index;not null"`
	// 游戏ID
	GameID uint64 `json:"gameId" gorm:"column:game_id;index;not null"`
	// 段位ID
	RankID uint64 `json:"rankId" gorm:"column:rank_id;index;not null"`
	// 认证状态
	Status PlayerRankStatus `json:"status" gorm:"size:32;index;default:'pending'"`
	// 段位截图（JSON数组）
	ScreenshotURLs string `json:"screenshotUrls,omitempty" gorm:"column:screenshot_urls;type:text"`
	// 认证时间
	VerifiedAt *time.Time `json:"verifiedAt,omitempty" gorm:"column:verified_at"`
	// 审核人ID
	VerifiedBy *uint64 `json:"verifiedBy,omitempty" gorm:"column:verified_by"`
	// 拒绝/撤销原因
	RejectReason string `json:"rejectReason,omitempty" gorm:"column:reject_reason;size:500"`
	// 过期时间（预留，定期复审）
	ExpireAt *time.Time `json:"expireAt,omitempty" gorm:"column:expire_at"`
	// 备注
	Remark string `json:"remark,omitempty" gorm:"size:500"`

	// 关联
	Player   *Player   `json:"player,omitempty" gorm:"foreignKey:PlayerID"`
	Game     *Game     `json:"game,omitempty" gorm:"foreignKey:GameID"`
	Rank     *GameRank `json:"rank,omitempty" gorm:"foreignKey:RankID"`
	Verifier *User     `json:"verifier,omitempty" gorm:"foreignKey:VerifiedBy"`
}

// TableName 指定表名
func (PlayerRankRecord) TableName() string {
	return "player_rank_records"
}

// ========== 陪玩师实名认证 ==========

// CertificationStatus 实名认证状态
type CertificationStatus string

const (
	CertificationStatusPending  CertificationStatus = "pending"  // 待审核
	CertificationStatusVerified CertificationStatus = "verified" // 已认证
	CertificationStatusRejected CertificationStatus = "rejected" // 已拒绝
)

// PlayerCertification 陪玩师实名认证
// @Description 陪玩师的实名认证信息
type PlayerCertification struct {
	Base
	// 陪玩师ID（唯一）
	PlayerID uint64 `json:"playerId" gorm:"column:player_id;uniqueIndex;not null"`
	// 真实姓名
	RealName string `json:"realName" gorm:"column:real_name;size:64"`
	// 身份证号（加密存储）
	IDCardNo string `json:"-" gorm:"column:id_card_no;size:255"`
	// 身份证正面照片URL
	IDCardFrontURL string `json:"idCardFrontUrl,omitempty" gorm:"column:id_card_front_url;size:255"`
	// 身份证背面照片URL
	IDCardBackURL string `json:"idCardBackUrl,omitempty" gorm:"column:id_card_back_url;size:255"`
	// 认证状态
	Status CertificationStatus `json:"status" gorm:"size:32;index;default:'pending'"`
	// 认证时间
	VerifiedAt *time.Time `json:"verifiedAt,omitempty" gorm:"column:verified_at"`
	// 审核人ID
	VerifiedBy *uint64 `json:"verifiedBy,omitempty" gorm:"column:verified_by"`
	// 拒绝原因
	RejectReason string `json:"rejectReason,omitempty" gorm:"column:reject_reason;size:500"`

	// ========== 预留字段 ==========
	// 个人照片URL
	PhotoURL string `json:"photoUrl,omitempty" gorm:"column:photo_url;size:255"`
	// 语音介绍URL
	VoiceURL string `json:"voiceUrl,omitempty" gorm:"column:voice_url;size:255"`
	// 扩展字段（JSON格式，预留其他认证资料）
	ExtJSON string `json:"extJson,omitempty" gorm:"column:ext_json;type:text"`

	// 关联
	Player   *Player `json:"player,omitempty" gorm:"foreignKey:PlayerID"`
	Verifier *User   `json:"verifier,omitempty" gorm:"foreignKey:VerifiedBy"`
}

// TableName 指定表名
func (PlayerCertification) TableName() string {
	return "player_certifications"
}
