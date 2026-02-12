package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
)

func seedReviewExtensions(tx *gorm.DB, users map[string]*model.User, players map[string]*model.Player, orders map[string]*model.Order) error {
	type replySpec struct {
		OrderKey string
		PlayerKey string
		Content string
		Status model.ReviewReplyStatus
		Offset time.Duration
	}
	replies := []replySpec{
		{OrderKey: "orderCompleted1", PlayerKey: "playerA", Content: "感谢认可，下次继续为您服务。", Status: model.ReviewReplyStatusApproved, Offset: -18 * time.Hour},
		{OrderKey: "orderCompleted2", PlayerKey: "playerD", Content: "感谢反馈，后续我会继续优化沟通节奏。", Status: model.ReviewReplyStatusApproved, Offset: -12 * time.Hour},
		{OrderKey: "orderCompleted3", PlayerKey: "playerE", Content: "感谢支持，欢迎随时再约。", Status: model.ReviewReplyStatusPending, Offset: -6 * time.Hour},
	}
	for _, spec := range replies {
		player := players[spec.PlayerKey]
		if player == nil {
			continue
		}
		var review model.Review
		found := false
		if order := orders[spec.OrderKey]; order != nil {
			if err := tx.Where("order_id = ?", order.ID).First(&review).Error; err == nil {
				found = true
			}
		}
		if !found {
			if err := tx.Where("player_id = ?", player.ID).Order("id asc").First(&review).Error; err == nil {
				found = true
			}
		}
		if !found {
			if err := tx.Order("id asc").First(&review).Error; err != nil {
				continue
			}
		}
		var existing model.ReviewReply
		err := tx.Where("review_id = ? AND player_id = ?", review.ID, player.ID).First(&existing).Error
		if err == nil {
			updates := map[string]interface{}{
				"content":      spec.Content,
				"status":       spec.Status,
				"reply_count":  1,
				"moderated_at": time.Now().Add(spec.Offset),
			}
			if err := tx.Model(&model.ReviewReply{}).Where("id = ?", existing.ID).Updates(updates).Error; err != nil {
				return err
			}
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		reply := &model.ReviewReply{
			ReviewID:    review.ID,
			PlayerID:    player.ID,
			AuthorID:    player.UserID,
			Content:     spec.Content,
			ReplyCount:  1,
			Status:      spec.Status,
			ModeratedAt: ptrTime(time.Now().Add(spec.Offset)),
		}
		if err := tx.Create(reply).Error; err != nil {
			return err
		}
	}

	type appealSpec struct {
		OrderKey string
		PlayerKey string
		Reason string
		Status model.AppealStatus
		Remark string
	}
	appeals := []appealSpec{
		{OrderKey: "orderRefunded1", PlayerKey: "playerB", Reason: "该评价与实际服务过程不符，申请复核。", Status: model.AppealStatusPending, Remark: ""},
		{OrderKey: "orderDispute1", PlayerKey: "playerC", Reason: "用户评价包含情绪化指责，请求重新判定。", Status: model.AppealStatusApproved, Remark: "申诉成立，评分已重新评估"},
		{OrderKey: "orderCanceled1", PlayerKey: "playerB", Reason: "订单取消非陪玩师责任，评价应无效。", Status: model.AppealStatusRejected, Remark: "证据不足，维持原判"},
	}
	var handlerID *uint64
	if admin := users["adminA"]; admin != nil {
		handlerID = &admin.ID
	}
	for _, spec := range appeals {
		player := players[spec.PlayerKey]
		if player == nil {
			continue
		}
		var review model.Review
		found := false
		if order := orders[spec.OrderKey]; order != nil {
			if err := tx.Where("order_id = ?", order.ID).First(&review).Error; err == nil {
				found = true
			}
		}
		if !found {
			if err := tx.Where("player_id = ? AND score <= ?", player.ID, 3).Order("id asc").First(&review).Error; err == nil {
				found = true
			}
		}
		if !found {
			if err := tx.Where("player_id = ?", player.ID).Order("id asc").First(&review).Error; err == nil {
				found = true
			}
		}
		if !found {
			continue
		}
		var existing model.ReviewAppeal
		err := tx.Where("review_id = ? AND player_id = ?", review.ID, player.ID).First(&existing).Error
		evidence := `["https://cdn.gamelink.com/evidence/review-appeal-1.png"]`
		if err == nil {
			updates := map[string]interface{}{
				"reason":         spec.Reason,
				"status":         spec.Status,
				"evidence_urls":  evidence,
				"handle_remark":  spec.Remark,
				"handled_by":     handlerID,
				"handled_at":     ptrTime(time.Now().Add(-4 * time.Hour)),
			}
			if spec.Status == model.AppealStatusPending {
				updates["handled_by"] = nil
				updates["handled_at"] = nil
			}
			if err := tx.Model(&model.ReviewAppeal{}).Where("id = ?", existing.ID).Updates(updates).Error; err != nil {
				return err
			}
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		appeal := &model.ReviewAppeal{
			ReviewID:      review.ID,
			PlayerID:      player.ID,
			Reason:        spec.Reason,
			EvidenceURLs:  evidence,
			Status:        spec.Status,
			HandleRemark:  spec.Remark,
		}
		if spec.Status != model.AppealStatusPending {
			appeal.HandledBy = handlerID
			appeal.HandledAt = ptrTime(time.Now().Add(-4 * time.Hour))
		}
		if err := tx.Create(appeal).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedChatSnapshots(tx *gorm.DB, orders map[string]*model.Order) error {
	var disputes []model.OrderDispute
	if err := tx.Order("id asc").Limit(3).Find(&disputes).Error; err != nil {
		return err
	}
	for i := range disputes {
		d := disputes[i]
		if d.ID == 0 || d.OrderID == 0 {
			continue
		}
		if d.ChatSnapshotID != nil && *d.ChatSnapshotID > 0 {
			continue
		}
		msgs := []map[string]interface{}{
			{"from": "user", "text": "我这边感觉服务节奏不太对", "ts": time.Now().Add(-55 * time.Minute).Format(time.RFC3339)},
			{"from": "player", "text": "收到，我马上调整打法", "ts": time.Now().Add(-52 * time.Minute).Format(time.RFC3339)},
			{"from": "user", "text": "沟通有点少，希望多报点位", "ts": time.Now().Add(-50 * time.Minute).Format(time.RFC3339)},
			{"from": "player", "text": "好的，后续我会实时报点", "ts": time.Now().Add(-48 * time.Minute).Format(time.RFC3339)},
			{"from": "service", "text": "请双方先继续完成本局，问题我们记录处理", "ts": time.Now().Add(-45 * time.Minute).Format(time.RFC3339)},
		}
		b, _ := json.Marshal(msgs)
		snapshot := &model.ChatSnapshot{
			DisputeID:   d.ID,
			OrderID:     d.OrderID,
			ChatGroupID: d.OrderID,
			Messages:    string(b),
			SnapshotAt:  time.Now().Add(-40 * time.Minute),
		}
		if err := tx.Create(snapshot).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.OrderDispute{}).Where("id = ?", d.ID).Update("chat_snapshot_id", snapshot.ID).Error; err != nil {
			return err
		}
	}
	_ = orders
	return nil
}

func seedUserSettings(tx *gorm.DB, users map[string]*model.User) error {
	type settingSpec struct {
		UserKey string
		Theme string
		Language string
		Notifications map[string]interface{}
		Privacy map[string]interface{}
	}
	specs := []settingSpec{
		{UserKey: "customerA", Theme: "auto", Language: "zh-CN", Notifications: map[string]interface{}{"order": true, "promo": true, "system": true}, Privacy: map[string]interface{}{"hideOnline": false}},
		{UserKey: "customerB", Theme: "light", Language: "zh-CN", Notifications: map[string]interface{}{"order": true, "promo": false, "system": true}, Privacy: map[string]interface{}{"hideOnline": true}},
		{UserKey: "customerC", Theme: "dark", Language: "en-US", Notifications: map[string]interface{}{"order": true, "promo": true, "system": false}, Privacy: map[string]interface{}{"hideOnline": false}},
		{UserKey: "proA", Theme: "auto", Language: "zh-CN", Notifications: map[string]interface{}{"order": true, "promo": false, "system": true}, Privacy: map[string]interface{}{"hideOnline": false}},
		{UserKey: "proB", Theme: "dark", Language: "zh-CN", Notifications: map[string]interface{}{"order": true, "promo": false, "system": true}, Privacy: map[string]interface{}{"hideOnline": false}},
		{UserKey: "adminA", Theme: "light", Language: "zh-CN", Notifications: map[string]interface{}{"order": true, "promo": false, "system": true}, Privacy: map[string]interface{}{"hideOnline": true}},
	}
	for _, spec := range specs {
		u := users[spec.UserKey]
		if u == nil {
			continue
		}
		nb, _ := json.Marshal(spec.Notifications)
		pb, _ := json.Marshal(spec.Privacy)
		var existing model.UserSettings
		err := tx.Where("user_id = ?", u.ID).First(&existing).Error
		if err == nil {
			if err := tx.Model(&model.UserSettings{}).Where("id = ?", existing.ID).Updates(map[string]interface{}{
				"theme": spec.Theme, "language": spec.Language, "notifications": string(nb), "privacy": string(pb),
			}).Error; err != nil {
				return err
			}
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		row := &model.UserSettings{
			UserID: u.ID, Theme: spec.Theme, Language: spec.Language, Notifications: string(nb), Privacy: string(pb),
		}
		if err := tx.Create(row).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedPaymentRoutingExtensions(tx *gorm.DB, users map[string]*model.User, orders map[string]*model.Order) error {
	var creatorID uint64
	if admin := users["adminA"]; admin != nil {
		creatorID = admin.ID
	}
	var entities []model.CollectionEntity
	if err := tx.Order("id asc").Find(&entities).Error; err != nil {
		return err
	}
	if len(entities) == 0 {
		return nil
	}
	defaultEntity := entities[0]

	type channelSpec struct {
		EntityID uint64
		Channel model.PaymentMethod
		MerchantNo string
		Enabled bool
		Priority int
		Remark string
	}
	channelSpecs := []channelSpec{
		{EntityID: defaultEntity.ID, Channel: model.PaymentMethodWeChat, MerchantNo: "wx_mch_10001", Enabled: true, Priority: 1, Remark: "微信主商户"},
		{EntityID: defaultEntity.ID, Channel: model.PaymentMethodAlipay, MerchantNo: "ali_mch_20001", Enabled: true, Priority: 1, Remark: "支付宝主商户"},
		{EntityID: defaultEntity.ID, Channel: model.PaymentMethodWallet, MerchantNo: "wallet_mch_30001", Enabled: false, Priority: 2, Remark: "钱包支付通道（停用）"},
	}
	for _, spec := range channelSpecs {
		var existing model.PaymentChannelConfig
		err := tx.Where("collection_entity_id = ? AND channel = ?", spec.EntityID, spec.Channel).First(&existing).Error
		if err == nil {
			if err := tx.Model(&model.PaymentChannelConfig{}).Where("id = ?", existing.ID).Updates(map[string]interface{}{
				"merchant_no": spec.MerchantNo,
				"enabled":     spec.Enabled,
				"priority":    spec.Priority,
				"remark":      spec.Remark,
				"callback_url": "https://api.gamelink.com/payment/callback",
			}).Error; err != nil {
				return err
			}
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		row := &model.PaymentChannelConfig{
			CollectionEntityID: spec.EntityID,
			Channel: spec.Channel,
			MerchantNo: spec.MerchantNo,
			MerchantKey: "seed-demo-key",
			CallbackURL: "https://api.gamelink.com/payment/callback",
			Enabled: spec.Enabled,
			Priority: spec.Priority,
			Remark: spec.Remark,
		}
		if err := tx.Create(row).Error; err != nil {
			return err
		}
	}

	makeCond := func(field model.ConditionField, op model.ConditionOperator, value interface{}) json.RawMessage {
		b, _ := json.Marshal([]model.RoutingCondition{{Field: field, Operator: op, Value: mustJSON(value)}})
		return b
	}
	ruleSpecs := []struct {
		Name string
		Priority int
		Conditions json.RawMessage
		TargetEntityID uint64
		Status model.RuleStatus
		Description string
	}{
		{Name: "大额订单分流", Priority: 10, Conditions: makeCond(model.ConditionFieldOrderAmount, model.ConditionOperatorGreaterThan, 50000), TargetEntityID: defaultEntity.ID, Status: model.RuleStatusActive, Description: "金额大于500元优先走主主体"},
		{Name: "MOBA 游戏分流", Priority: 20, Conditions: makeCond(model.ConditionFieldGameType, model.ConditionOperatorIn, []string{"moba"}), TargetEntityID: defaultEntity.ID, Status: model.RuleStatusActive, Description: "MOBA订单专项分流"},
		{Name: "默认兜底规则", Priority: 999, Conditions: json.RawMessage(`[]`), TargetEntityID: defaultEntity.ID, Status: model.RuleStatusInactive, Description: "历史兜底规则，现停用"},
	}
	for _, spec := range ruleSpecs {
		var existing model.RoutingRule
		err := tx.Where("name = ?", spec.Name).First(&existing).Error
		if err == nil {
			if err := tx.Model(&model.RoutingRule{}).Where("id = ?", existing.ID).Updates(map[string]interface{}{
				"priority": spec.Priority, "conditions": spec.Conditions, "target_entity_id": spec.TargetEntityID, "status": spec.Status, "description": spec.Description, "updated_by": creatorID,
			}).Error; err != nil {
				return err
			}
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			row := &model.RoutingRule{
				Name: spec.Name, Priority: spec.Priority, Conditions: spec.Conditions, TargetEntityID: spec.TargetEntityID,
				Status: spec.Status, Description: spec.Description, CreatedBy: creatorID, UpdatedBy: &creatorID,
			}
			if err := tx.Create(row).Error; err != nil {
				return err
			}
			existing = *row
		} else {
			return err
		}

		// history: each rule at least one change record
		var hCount int64
		if err := tx.Model(&model.RoutingRuleHistory{}).Where("routing_rule_id = ?", existing.ID).Count(&hCount).Error; err != nil {
			return err
		}
		if hCount == 0 {
			h := &model.RoutingRuleHistory{
				RoutingRuleID: existing.ID,
				FieldName: "status",
				OldValue: "draft",
				NewValue: string(existing.Status),
				ChangedBy: creatorID,
			}
			if err := tx.Create(h).Error; err != nil {
				return err
			}
		}
	}

	// Routing logs: attach to existing payments
	var payments []model.Payment
	if err := tx.Order("id asc").Limit(8).Find(&payments).Error; err != nil {
		return err
	}
	if len(payments) > 0 {
		var activeRule model.RoutingRule
		_ = tx.Where("status = ?", model.RuleStatusActive).Order("priority asc").First(&activeRule).Error
		for i, p := range payments {
			var existing model.RoutingLog
			if err := tx.Where("payment_id = ?", p.ID).First(&existing).Error; err == nil {
				continue
			}
			matchedID := activeRule.ID
			isFallback := i%3 == 0
			detail := fmt.Sprintf(`{"seed":"demo","rule":"%d","seq":%d}`, matchedID, i+1)
			logRow := &model.RoutingLog{
				PaymentID: p.ID,
				OrderID: p.OrderID,
				MatchedRuleID: &matchedID,
				CollectionEntityID: defaultEntity.ID,
				MerchantNo: "wx_mch_10001",
				IsDefault: false,
				IsFallback: isFallback,
				MatchDetails: detail,
			}
			if err := tx.Create(logRow).Error; err != nil {
				return err
			}
		}
	}

	_ = orders
	return nil
}

func mustJSON(v interface{}) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}
