package model

import "time"

// ============================================================================
// 用户统计 (UserStatistics)
// ============================================================================

// UserStatistics 用户统计数据
// @Description 用户消费行为统计，用于标签判定和用户画像
type UserStatistics struct {
	Base
	UserID uint64 `json:"userId" gorm:"column:user_id;uniqueIndex;not null"`

	// 消费统计
	TotalSpentCents     int64 `json:"totalSpentCents" gorm:"column:total_spent_cents;default:0"`          // 累计消费金额（分）
	TotalOrderCount     int   `json:"totalOrderCount" gorm:"column:total_order_count;default:0"`          // 累计订单数
	CompletedOrderCount int   `json:"completedOrderCount" gorm:"column:completed_order_count;default:0"`  // 完成订单数
	CanceledOrderCount  int   `json:"canceledOrderCount" gorm:"column:canceled_order_count;default:0"`    // 取消订单数
	RefundOrderCount    int   `json:"refundOrderCount" gorm:"column:refund_order_count;default:0"`        // 退款订单数
	AvgOrderAmountCents int64 `json:"avgOrderAmountCents" gorm:"column:avg_order_amount_cents;default:0"` // 平均订单金额（分）

	// 争议统计
	DisputeCount     int `json:"disputeCount" gorm:"column:dispute_count;default:0"`          // 发起争议次数
	DisputeWonCount  int `json:"disputeWonCount" gorm:"column:dispute_won_count;default:0"`   // 争议胜诉次数
	DisputeLostCount int `json:"disputeLostCount" gorm:"column:dispute_lost_count;default:0"` // 争议败诉次数

	// 充值统计
	TotalRechargeCents     int64 `json:"totalRechargeCents" gorm:"column:total_recharge_cents;default:0"`          // 累计充值金额（分）
	RechargeCount          int   `json:"rechargeCount" gorm:"column:recharge_count;default:0"`                     // 充值次数
	AvgRechargeAmountCents int64 `json:"avgRechargeAmountCents" gorm:"column:avg_recharge_amount_cents;default:0"` // 平均充值金额（分）

	// 评价统计
	ReviewCount    int     `json:"reviewCount" gorm:"column:review_count;default:0"`        // 评价次数
	AvgReviewScore float32 `json:"avgReviewScore" gorm:"column:avg_review_score;default:0"` // 平均评分

	// 时间节点
	FirstOrderAt   *time.Time `json:"firstOrderAt,omitempty" gorm:"column:first_order_at"`     // 首次下单时间
	LastOrderAt    *time.Time `json:"lastOrderAt,omitempty" gorm:"column:last_order_at"`       // 最后下单时间
	LastRechargeAt *time.Time `json:"lastRechargeAt,omitempty" gorm:"column:last_recharge_at"` // 最后充值时间

	// 关联
	User User `json:"-" gorm:"foreignKey:UserID;references:ID"`
}

func (UserStatistics) TableName() string {
	return "user_statistics"
}

// GetCompletionRate 获取订单完成率
func (s *UserStatistics) GetCompletionRate() float32 {
	if s.TotalOrderCount == 0 {
		return 0
	}
	return float32(s.CompletedOrderCount) / float32(s.TotalOrderCount)
}

// GetCancelRate 获取订单取消率
func (s *UserStatistics) GetCancelRate() float32 {
	if s.TotalOrderCount == 0 {
		return 0
	}
	return float32(s.CanceledOrderCount) / float32(s.TotalOrderCount)
}

// GetDisputeWinRate 获取争议胜诉率
func (s *UserStatistics) GetDisputeWinRate() float32 {
	total := s.DisputeWonCount + s.DisputeLostCount
	if total == 0 {
		return 0
	}
	return float32(s.DisputeWonCount) / float32(total)
}

// ============================================================================
// 陪玩师统计 (PlayerStatistics)
// ============================================================================

// PlayerStatistics 陪玩师统计数据
// @Description 陪玩师业绩统计，用于标签判定和排名
type PlayerStatistics struct {
	Base
	PlayerID uint64 `json:"playerId" gorm:"column:player_id;uniqueIndex;not null"`

	// 收入统计
	TotalEarningsCents   int64 `json:"totalEarningsCents" gorm:"column:total_earnings_cents;default:0"`     // 累计收入（分）
	TotalCommissionCents int64 `json:"totalCommissionCents" gorm:"column:total_commission_cents;default:0"` // 累计被抽成（分）
	TotalWithdrawCents   int64 `json:"totalWithdrawCents" gorm:"column:total_withdraw_cents;default:0"`     // 累计提现（分）
	PendingWithdrawCents int64 `json:"pendingWithdrawCents" gorm:"column:pending_withdraw_cents;default:0"` // 待提现（分）

	// 订单统计
	TotalOrderCount     int `json:"totalOrderCount" gorm:"column:total_order_count;default:0"`         // 累计接单数
	CompletedOrderCount int `json:"completedOrderCount" gorm:"column:completed_order_count;default:0"` // 完成订单数
	CanceledOrderCount  int `json:"canceledOrderCount" gorm:"column:canceled_order_count;default:0"`   // 取消订单数
	RefundOrderCount    int `json:"refundOrderCount" gorm:"column:refund_order_count;default:0"`       // 退款订单数

	// 服务统计
	TotalServiceMinutes int   `json:"totalServiceMinutes" gorm:"column:total_service_minutes;default:0"`  // 累计服务时长（分钟）
	AvgResponseTimeSec  int   `json:"avgResponseTimeSec" gorm:"column:avg_response_time_sec;default:0"`   // 平均响应时间（秒）
	AvgOrderAmountCents int64 `json:"avgOrderAmountCents" gorm:"column:avg_order_amount_cents;default:0"` // 平均订单金额（分）

	// 客户统计
	TotalCustomerCount  int     `json:"totalCustomerCount" gorm:"column:total_customer_count;default:0"`   // 累计服务客户数
	RepeatCustomerCount int     `json:"repeatCustomerCount" gorm:"column:repeat_customer_count;default:0"` // 回头客数量
	RepeatOrderRate     float32 `json:"repeatOrderRate" gorm:"column:repeat_order_rate;default:0"`         // 复购率

	// 争议统计
	DisputeCount     int `json:"disputeCount" gorm:"column:dispute_count;default:0"`          // 被投诉次数
	DisputeWonCount  int `json:"disputeWonCount" gorm:"column:dispute_won_count;default:0"`   // 投诉胜诉次数
	DisputeLostCount int `json:"disputeLostCount" gorm:"column:dispute_lost_count;default:0"` // 投诉败诉次数

	// 礼物统计
	GiftReceivedCount       int   `json:"giftReceivedCount" gorm:"column:gift_received_count;default:0"`              // 收到礼物数
	GiftReceivedAmountCents int64 `json:"giftReceivedAmountCents" gorm:"column:gift_received_amount_cents;default:0"` // 收到礼物金额（分）

	// 时间节点
	FirstOrderAt *time.Time `json:"firstOrderAt,omitempty" gorm:"column:first_order_at"` // 首次接单时间
	LastOrderAt  *time.Time `json:"lastOrderAt,omitempty" gorm:"column:last_order_at"`   // 最后接单时间
	LastActiveAt *time.Time `json:"lastActiveAt,omitempty" gorm:"column:last_active_at"` // 最后活跃时间

	// 关联
	Player Player `json:"-" gorm:"foreignKey:PlayerID;references:ID"`
}

func (PlayerStatistics) TableName() string {
	return "player_statistics"
}

// GetCompletionRate 获取订单完成率
func (s *PlayerStatistics) GetCompletionRate() float32 {
	if s.TotalOrderCount == 0 {
		return 0
	}
	return float32(s.CompletedOrderCount) / float32(s.TotalOrderCount)
}

// GetCancelRate 获取订单取消率
func (s *PlayerStatistics) GetCancelRate() float32 {
	if s.TotalOrderCount == 0 {
		return 0
	}
	return float32(s.CanceledOrderCount) / float32(s.TotalOrderCount)
}

// GetTotalServiceHours 获取累计服务时长（小时）
func (s *PlayerStatistics) GetTotalServiceHours() float64 {
	return float64(s.TotalServiceMinutes) / 60.0
}

// ============================================================================
// 服务项目统计 (ServiceItemStatistics)
// ============================================================================

// ServiceItemStatistics 服务项目统计数据
type ServiceItemStatistics struct {
	Base
	ServiceItemID uint64 `json:"serviceItemId" gorm:"column:service_item_id;uniqueIndex;not null"`

	// 销售统计
	TotalSalesCount       int   `json:"totalSalesCount" gorm:"column:total_sales_count;default:0"`              // 累计销量
	TotalSalesAmountCents int64 `json:"totalSalesAmountCents" gorm:"column:total_sales_amount_cents;default:0"` // 累计销售额（分）
	TotalCommissionCents  int64 `json:"totalCommissionCents" gorm:"column:total_commission_cents;default:0"`    // 累计抽成（分）

	// 评价统计
	ReviewCount   int     `json:"reviewCount" gorm:"column:review_count;default:0"`      // 评价数
	AvgRating     float32 `json:"avgRating" gorm:"column:avg_rating;default:0"`          // 平均评分
	FiveStarCount int     `json:"fiveStarCount" gorm:"column:five_star_count;default:0"` // 五星评价数

	// 退款统计
	RefundCount       int     `json:"refundCount" gorm:"column:refund_count;default:0"`              // 退款次数
	RefundAmountCents int64   `json:"refundAmountCents" gorm:"column:refund_amount_cents;default:0"` // 退款金额（分）
	RefundRate        float32 `json:"refundRate" gorm:"column:refund_rate;default:0"`                // 退款率

	// 时间节点
	FirstSoldAt *time.Time `json:"firstSoldAt,omitempty" gorm:"column:first_sold_at"` // 首次售出时间
	LastSoldAt  *time.Time `json:"lastSoldAt,omitempty" gorm:"column:last_sold_at"`   // 最后售出时间

	// 关联
	ServiceItem ServiceItem `json:"-" gorm:"foreignKey:ServiceItemID;references:ID"`
}

func (ServiceItemStatistics) TableName() string {
	return "service_item_statistics"
}

// GetFiveStarRate 获取五星好评率
func (s *ServiceItemStatistics) GetFiveStarRate() float32 {
	if s.ReviewCount == 0 {
		return 0
	}
	return float32(s.FiveStarCount) / float32(s.ReviewCount)
}

// ============================================================================
// 游戏统计 (GameStatistics)
// ============================================================================

// GameStatistics 游戏统计数据
type GameStatistics struct {
	Base
	GameID uint64 `json:"gameId" gorm:"column:game_id;uniqueIndex;not null"`

	// 订单统计
	TotalOrderCount      int   `json:"totalOrderCount" gorm:"column:total_order_count;default:0"`           // 累计订单数
	TotalGMVCents        int64 `json:"totalGmvCents" gorm:"column:total_gmv_cents;default:0"`               // 累计GMV（分）
	TotalCommissionCents int64 `json:"totalCommissionCents" gorm:"column:total_commission_cents;default:0"` // 累计抽成（分）
	AvgOrderAmountCents  int64 `json:"avgOrderAmountCents" gorm:"column:avg_order_amount_cents;default:0"`  // 平均客单价（分）

	// 陪玩师统计
	TotalPlayerCount  int `json:"totalPlayerCount" gorm:"column:total_player_count;default:0"`   // 总陪玩师数
	ActivePlayerCount int `json:"activePlayerCount" gorm:"column:active_player_count;default:0"` // 活跃陪玩师数（7天内有接单）

	// 用户统计
	TotalUserCount  int `json:"totalUserCount" gorm:"column:total_user_count;default:0"`   // 总用户数
	ActiveUserCount int `json:"activeUserCount" gorm:"column:active_user_count;default:0"` // 活跃用户数（7天内有下单）

	// 评价统计
	AvgRating float32 `json:"avgRating" gorm:"column:avg_rating;default:0"` // 平均评分

	// 关联
	Game Game `json:"-" gorm:"foreignKey:GameID;references:ID"`
}

func (GameStatistics) TableName() string {
	return "game_statistics"
}

// ============================================================================
// 平台统计 (PlatformStatistics) - 全局汇总
// ============================================================================

// PlatformStatistics 平台统计数据（按日期汇总）
type PlatformStatistics struct {
	Base
	StatDate time.Time `json:"statDate" gorm:"column:stat_date;uniqueIndex;not null"` // 统计日期

	// 订单统计
	DailyOrderCount        int   `json:"dailyOrderCount" gorm:"column:daily_order_count;default:0"`                // 日订单数
	DailyCompletedCount    int   `json:"dailyCompletedCount" gorm:"column:daily_completed_count;default:0"`        // 日完成订单数
	DailyCanceledCount     int   `json:"dailyCanceledCount" gorm:"column:daily_canceled_count;default:0"`          // 日取消订单数
	DailyGMVCents          int64 `json:"dailyGmvCents" gorm:"column:daily_gmv_cents;default:0"`                    // 日GMV（分）
	DailyCommissionCents   int64 `json:"dailyCommissionCents" gorm:"column:daily_commission_cents;default:0"`      // 日抽成（分）
	DailyRefundAmountCents int64 `json:"dailyRefundAmountCents" gorm:"column:daily_refund_amount_cents;default:0"` // 日退款金额（分）

	// 用户统计
	DailyNewUserCount    int `json:"dailyNewUserCount" gorm:"column:daily_new_user_count;default:0"`       // 日新增用户数
	DailyActiveUserCount int `json:"dailyActiveUserCount" gorm:"column:daily_active_user_count;default:0"` // 日活跃用户数
	DailyPayingUserCount int `json:"dailyPayingUserCount" gorm:"column:daily_paying_user_count;default:0"` // 日付费用户数

	// 陪玩师统计
	DailyNewPlayerCount    int `json:"dailyNewPlayerCount" gorm:"column:daily_new_player_count;default:0"`       // 日新增陪玩师数
	DailyActivePlayerCount int `json:"dailyActivePlayerCount" gorm:"column:daily_active_player_count;default:0"` // 日活跃陪玩师数

	// 充值提现统计
	DailyRechargeCents int64 `json:"dailyRechargeCents" gorm:"column:daily_recharge_cents;default:0"` // 日充值金额（分）
	DailyWithdrawCents int64 `json:"dailyWithdrawCents" gorm:"column:daily_withdraw_cents;default:0"` // 日提现金额（分）
	DailyRechargeCount int   `json:"dailyRechargeCount" gorm:"column:daily_recharge_count;default:0"` // 日充值笔数
	DailyWithdrawCount int   `json:"dailyWithdrawCount" gorm:"column:daily_withdraw_count;default:0"` // 日提现笔数

	// 争议统计
	DailyDisputeCount   int `json:"dailyDisputeCount" gorm:"column:daily_dispute_count;default:0"`      // 日争议数
	DailyResolvedCount  int `json:"dailyResolvedCount" gorm:"column:daily_resolved_count;default:0"`    // 日解决争议数
	DailySLABreachCount int `json:"dailySlaBreachCount" gorm:"column:daily_sla_breach_count;default:0"` // 日SLA超时数
}

func (PlatformStatistics) TableName() string {
	return "platform_statistics"
}

// GetCompletionRate 获取日订单完成率
func (s *PlatformStatistics) GetCompletionRate() float32 {
	if s.DailyOrderCount == 0 {
		return 0
	}
	return float32(s.DailyCompletedCount) / float32(s.DailyOrderCount)
}

// ============================================================================
// 标签阈值配置 (TagThreshold)
// ============================================================================

// TagTargetType 标签目标类型
type TagTargetType string

const (
	TagTargetUser        TagTargetType = "user"
	TagTargetPlayer      TagTargetType = "player"
	TagTargetServiceItem TagTargetType = "service_item"
	TagTargetGame        TagTargetType = "game"
)

// TagOperator 比较运算符
type TagOperator string

const (
	TagOperatorGT  TagOperator = "gt"  // 大于
	TagOperatorGTE TagOperator = "gte" // 大于等于
	TagOperatorLT  TagOperator = "lt"  // 小于
	TagOperatorLTE TagOperator = "lte" // 小于等于
	TagOperatorEQ  TagOperator = "eq"  // 等于
)

// TagThreshold 标签阈值配置
// @Description 定义标签的自动判定规则
type TagThreshold struct {
	Base
	TagID       uint64        `json:"tagId" gorm:"column:tag_id;not null;index"`                // 关联的标签ID
	TargetType  TagTargetType `json:"targetType" gorm:"column:target_type;size:32;index"`       // 目标类型
	MetricField string        `json:"metricField" gorm:"column:metric_field;size:64;not null"`  // 指标字段名
	Operator    TagOperator   `json:"operator" gorm:"column:operator;size:10;not null"`         // 比较运算符
	Threshold   float64       `json:"threshold" gorm:"column:threshold;not null"`               // 阈值
	Priority    int           `json:"priority" gorm:"column:priority;default:0"`                // 优先级（多条件时）
	IsActive    bool          `json:"isActive" gorm:"column:is_active;default:true"`            // 是否启用
	Description string        `json:"description,omitempty" gorm:"column:description;size:255"` // 描述
}

func (TagThreshold) TableName() string {
	return "tag_thresholds"
}

// Evaluate 评估指标值是否满足阈值条件
func (t *TagThreshold) Evaluate(value float64) bool {
	switch t.Operator {
	case TagOperatorGT:
		return value > t.Threshold
	case TagOperatorGTE:
		return value >= t.Threshold
	case TagOperatorLT:
		return value < t.Threshold
	case TagOperatorLTE:
		return value <= t.Threshold
	case TagOperatorEQ:
		return value == t.Threshold
	default:
		return false
	}
}
