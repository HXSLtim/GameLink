package event

import (
	"context"
	"log/slog"

	"gorm.io/gorm"

	"gamelink/internal/model"
)

// EventPublisher 事件发布器，用于在业务模块间解耦
type EventPublisher struct {
	db     *gorm.DB
	logger *slog.Logger
}

// NewEventPublisher 创建事件发布器
func NewEventPublisher(db *gorm.DB) *EventPublisher {
	return &EventPublisher{
		db:     db,
		logger: slog.Default(),
	}
}

// PublishGiftReceived 发布礼物收到事件
func (p *EventPublisher) PublishGiftReceived(ctx context.Context, playerUserID uint64, senderName string, giftName string, quantity int, orderID uint64) {
	title := "收到新礼物"
	message := senderName + " 赠送了您 " + giftName
	if quantity > 1 {
		message += " x" + formatInt(quantity)
	}

	p.createNotification(ctx, playerUserID, "gift", title, message, &orderID)
}

// PublishOrderAccepted 发布订单被接单事件
func (p *EventPublisher) PublishOrderAccepted(ctx context.Context, userID uint64, playerName string, orderID uint64) {
	title := "订单已被接单"
	message := "陪玩师 " + playerName + " 已接受您的订单"
	p.createNotification(ctx, userID, "order", title, message, &orderID)
}

// PublishOrderCompleted 发布订单完成事件
func (p *EventPublisher) PublishOrderCompleted(ctx context.Context, userID uint64, orderID uint64) {
	title := "订单已完成"
	message := "您的订单已完成，欢迎评价"
	p.createNotification(ctx, userID, "order", title, message, &orderID)
}

// PublishIncomeSettled 发布收入结算事件
func (p *EventPublisher) PublishIncomeSettled(ctx context.Context, playerUserID uint64, amountCents int64) {
	title := "收入已到账"
	message := "您有 ¥" + formatCents(amountCents) + " 收入已解冻，可申请提现"
	p.createNotification(ctx, playerUserID, "income", title, message, nil)
}

// PublishWithdrawResult 发布提现结果事件
func (p *EventPublisher) PublishWithdrawResult(ctx context.Context, playerUserID uint64, amountCents int64, success bool, reason string) {
	var title, message string
	if success {
		title = "提现成功"
		message = "您的提现 ¥" + formatCents(amountCents) + " 已到账"
	} else {
		title = "提现失败"
		message = "您的提现申请被拒绝：" + reason
	}
	p.createNotification(ctx, playerUserID, "withdraw", title, message, nil)
}

// PublishRankVerified 发布段位认证结果事件
func (p *EventPublisher) PublishRankVerified(ctx context.Context, playerUserID uint64, gameName string, rankName string, approved bool, reason string) {
	var title, message string
	if approved {
		title = "段位认证通过"
		message = "您的 " + gameName + " " + rankName + " 段位认证已通过"
	} else {
		title = "段位认证未通过"
		message = "您的 " + gameName + " 段位认证未通过：" + reason
	}
	p.createNotification(ctx, playerUserID, "rank", title, message, nil)
}

// PublishCertificationVerified 发布实名认证结果事件
func (p *EventPublisher) PublishCertificationVerified(ctx context.Context, playerUserID uint64, approved bool, reason string) {
	var title, message string
	if approved {
		title = "实名认证通过"
		message = "您的实名认证已通过"
	} else {
		title = "实名认证未通过"
		message = "您的实名认证未通过：" + reason
	}
	p.createNotification(ctx, playerUserID, "certification", title, message, nil)
}

// PublishPlayerVerified 发布陪玩师入驻审核结果事件
func (p *EventPublisher) PublishPlayerVerified(ctx context.Context, userID uint64, approved bool, reason string) {
	var title, message string
	if approved {
		title = "入驻申请通过"
		message = "恭喜您，陪玩师入驻申请已通过，现在可以开始接单了"
	} else {
		title = "入驻申请未通过"
		message = "您的陪玩师入驻申请未通过：" + reason
	}
	p.createNotification(ctx, userID, "player", title, message, nil)
}

// PublishDisputeResolved 发布争议处理结果事件
func (p *EventPublisher) PublishDisputeResolved(ctx context.Context, userID uint64, disputeID uint64, result string, refundAmount int64) {
	var title, message string
	if result == "refund" {
		title = "争议已处理 - 退款"
		message = "您的争议已处理，退款 ¥" + formatCents(refundAmount) + " 将原路退回"
	} else {
		title = "争议已处理"
		message = "您的争议已处理，详情请查看订单"
	}
	p.createNotification(ctx, userID, "dispute", title, message, &disputeID)
}

// PublishLFGMatched 发布 LFG 匹配成功事件
func (p *EventPublisher) PublishLFGMatched(ctx context.Context, userID uint64, playerName string, requestID uint64) {
	title := "匹配成功"
	message := "陪玩师 " + playerName + " 已响应您的匹配请求，请前往聊天室确认"
	p.createNotification(ctx, userID, "lfg", title, message, &requestID)
}

// PublishLFGExpired 发布 LFG 请求过期事件
func (p *EventPublisher) PublishLFGExpired(ctx context.Context, userID uint64, requestID uint64) {
	title := "匹配请求已过期"
	message := "您的匹配请求已过期，可以重新发起"
	p.createNotification(ctx, userID, "lfg", title, message, &requestID)
}

// PublishReferralReward 发布推荐奖励事件
func (p *EventPublisher) PublishReferralReward(ctx context.Context, userID uint64, rewardType string, amountCents int64) {
	title := "推荐奖励已到账"
	var message string
	if rewardType == "cash" {
		message = "恭喜您获得推荐奖励 ¥" + formatCents(amountCents) + "，已存入您的钱包"
	} else {
		message = "恭喜您获得推荐奖励，优惠券已发放到您的账户"
	}
	p.createNotification(ctx, userID, "referral", title, message, nil)
}

// createNotification 创建通知记录
func (p *EventPublisher) createNotification(ctx context.Context, userID uint64, refType, title, message string, refID *uint64) {
	notification := &model.NotificationEvent{
		UserID:        userID,
		Title:         title,
		Message:       message,
		Priority:      model.NotificationPriorityNormal,
		Channel:       "in_app",
		ReferenceType: refType,
		ReferenceID:   refID,
	}

	if err := p.db.WithContext(ctx).Create(notification).Error; err != nil {
		p.logger.Error("Failed to create notification",
			"userID", userID,
			"type", refType,
			"error", err)
	}
}

// formatInt 格式化整数
func formatInt(n int) string {
	if n < 10 {
		return string(rune('0' + n))
	}
	return formatInt(n/10) + string(rune('0'+n%10))
}

// formatCents 格式化分为元
func formatCents(cents int64) string {
	yuan := cents / 100
	fen := cents % 100
	if fen == 0 {
		return formatInt64(yuan)
	}
	return formatInt64(yuan) + "." + padZero(int(fen))
}

func formatInt64(n int64) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + formatInt64(-n)
	}
	if n < 10 {
		return string(rune('0' + n))
	}
	return formatInt64(n/10) + string(rune('0'+n%10))
}

func padZero(n int) string {
	if n < 10 {
		return "0" + string(rune('0'+n))
	}
	return string(rune('0'+n/10)) + string(rune('0'+n%10))
}
