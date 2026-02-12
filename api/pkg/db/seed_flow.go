package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
)

// seedAdditionalFlowOrders creates extra orders that help cover end-to-end business flows:
// - Gift order (gift → payment → delivered)
// - Team order (team item → order_items + order_players)
// - Failed payment order (payment failed → order canceled/timeout log)
func seedAdditionalFlowOrders(
	tx *gorm.DB,
	now time.Time,
	users map[string]*model.User,
	players map[string]*model.Player,
	serviceItems map[string]*model.ServiceItem,
	orders map[string]*model.Order,
) error {
	mustUser := func(key string) (*model.User, error) {
		u := users[key]
		if u == nil {
			return nil, fmt.Errorf("seed flow missing user %s", key)
		}
		return u, nil
	}
	mustPlayer := func(key string) (*model.Player, error) {
		p := players[key]
		if p == nil {
			return nil, fmt.Errorf("seed flow missing player %s", key)
		}
		return p, nil
	}
	mustItem := func(code string) (*model.ServiceItem, error) {
		item := serviceItems[code]
		if item == nil {
			return nil, fmt.Errorf("seed flow missing service item %s", code)
		}
		return item, nil
	}

	// 1) Gift order (wallet payment)
	{
		user, err := mustUser("customerB")
		if err != nil {
			return err
		}
		recipient, err := mustPlayer("playerA")
		if err != nil {
			return err
		}
		item, err := mustItem("gift-rose")
		if err != nil {
			return err
		}

		title := "礼物订单-玫瑰（演示）"
		var order *model.Order
		var existing model.Order
		if err := tx.Where("title = ? AND user_id = ?", title, user.ID).First(&existing).Error; err == nil {
			orders["orderGift1"] = &existing
			order = &existing
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		} else {
			qty := 3
			unit := item.BasePriceCents
			total := unit * int64(qty)
			commission := int64(float64(total) * item.CommissionRate)
			playerIncome := total - commission

			deliveredAt := now.Add(-15 * time.Minute)
			created := &model.Order{
				OrderNo:           model.GenerateGiftOrderNo(),
				UserID:            user.ID,
				ItemID:            item.ID,
				RecipientPlayerID: &recipient.ID,
				Quantity:          qty,
				UnitPriceCents:    unit,
				TotalPriceCents:   total,
				CommissionCents:   commission,
				PlayerIncomeCents: playerIncome,
				Currency:          model.CurrencyCNY,
				Status:            model.OrderStatusCompleted,
				Title:             title,
				Description:       "演示数据：用户给陪玩师发送礼物（即时完成）",
				GiftMessage:       "辛苦啦，继续加油！",
				IsAnonymous:       true,
				DeliveredAt:       &deliveredAt,
				CompletedAt:       &deliveredAt,
				OrderConfig:       `{"seed":"demo","flow":"gift"}`,
			}
			if err := tx.Create(created).Error; err != nil {
				return err
			}
			orders["orderGift1"] = created
			order = created
		}

		paidAt := now.Add(-20 * time.Minute)
		if order != nil && order.CompletedAt != nil {
			paidAt = *order.CompletedAt
		}
		if order != nil {
			if err := seedPayment(tx, seedPaymentParams{
				OrderID:               order.ID,
				UserID:                user.ID,
				Method:                model.PaymentMethodWallet,
				AmountCents:           order.TotalPriceCents,
				Currency:              model.CurrencyCNY,
				Status:                model.PaymentStatusPaid,
				ProviderTradeNo:       fmt.Sprintf("WALLET-%s", order.OrderNo),
				ProviderRaw:           json.RawMessage(`{"channel":"wallet","seed":"demo"}`),
				PaidAt:                &paidAt,
				WalletAmountCents:     order.TotalPriceCents,
				ThirdPartyMethod:      "",
				ThirdPartyAmountCents: 0,
			}); err != nil {
				return err
			}
		}
	}

	// 2) Team order (creates order_items + order_players)
	{
		user, err := mustUser("customerA")
		if err != nil {
			return err
		}
		playerA, err := mustPlayer("playerA")
		if err != nil {
			return err
		}
		playerC, err := mustPlayer("playerC")
		if err != nil {
			return err
		}
		item, err := mustItem("escort-lol-team")
		if err != nil {
			return err
		}

		title := "双人车队订单（演示）"
		var order *model.Order
		var existing model.Order
		if err := tx.Where("title = ? AND user_id = ?", title, user.ID).First(&existing).Error; err == nil {
			orders["orderTeamInProgress1"] = &existing
			order = &existing
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		} else {
			required := 2
			unit := item.BasePriceCents
			total := unit * int64(required)
			commission := int64(float64(total) * item.CommissionRate)
			playerIncome := total - commission

			start := now.Add(-30 * time.Minute)
			end := now.Add(30 * time.Minute)
			created := &model.Order{
				OrderNo:           model.GenerateEscortOrderNo(),
				UserID:            user.ID,
				ItemID:            item.ID,
				GameID:            item.GameID,
				Quantity:          1,
				UnitPriceCents:    unit,
				TotalPriceCents:   total,
				CommissionCents:   commission,
				PlayerIncomeCents: playerIncome,
				Currency:          model.CurrencyCNY,
				Status:            model.OrderStatusInProgress,
				Title:             title,
				Description:       "演示数据：多人服务订单（2名陪玩师同时服务）",
				ScheduledStart:    &start,
				ScheduledEnd:      &end,
				StartedAt:         &start,
				RequiredPlayers:   required,
				CurrentPlayers:    required,
				OrderConfig:       `{"seed":"demo","flow":"team","requiredPlayers":2}`,
			}
			if err := tx.Create(created).Error; err != nil {
				return err
			}
			orders["orderTeamInProgress1"] = created
			order = created
		}

		paidAt := now.Add(-35 * time.Minute)
		if order != nil {
			if err := ensureOrderSlot(tx, order, item, 1, playerA); err != nil {
				return err
			}
			if err := ensureOrderSlot(tx, order, item, 2, playerC); err != nil {
				return err
			}

			walletPart := int64(10000)
			if walletPart > order.TotalPriceCents {
				walletPart = order.TotalPriceCents
			}
			thirdPart := order.TotalPriceCents - walletPart
			if err := seedPayment(tx, seedPaymentParams{
				OrderID:               order.ID,
				UserID:                user.ID,
				Method:                model.PaymentMethodCombined,
				AmountCents:           order.TotalPriceCents,
				Currency:              model.CurrencyCNY,
				Status:                model.PaymentStatusPaid,
				ProviderTradeNo:       fmt.Sprintf("WX-CMB-%s", order.OrderNo),
				ProviderRaw:           json.RawMessage(`{"channel":"combined","thirdParty":"wechat","seed":"demo"}`),
				PaidAt:                &paidAt,
				WalletAmountCents:     walletPart,
				ThirdPartyMethod:      model.PaymentMethodWeChat,
				ThirdPartyAmountCents: thirdPart,
			}); err != nil {
				return err
			}
		}
	}

	// 3) Payment failed order (used for timeout/failure flows)
	{
		user, err := mustUser("customerC")
		if err != nil {
			return err
		}
		item, err := mustItem("escort-default")
		if err != nil {
			return err
		}

		title := "支付失败订单（演示）"
		var order *model.Order
		var existing model.Order
		if err := tx.Where("title = ? AND user_id = ?", title, user.ID).First(&existing).Error; err == nil {
			orders["orderPaymentFailed1"] = &existing
			order = &existing
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		} else {
			start := now.Add(2 * time.Hour)
			end := now.Add(3 * time.Hour)
			created := &model.Order{
				OrderNo:         model.GenerateEscortOrderNo(),
				UserID:          user.ID,
				ItemID:          item.ID,
				GameID:          item.GameID,
				Quantity:        1,
				UnitPriceCents:  item.BasePriceCents,
				TotalPriceCents: item.BasePriceCents,
				Currency:        model.CurrencyCNY,
				Status:          model.OrderStatusCanceled,
				Title:           title,
				Description:     "演示数据：支付失败后自动取消",
				ScheduledStart:  &start,
				ScheduledEnd:    &end,
				CancelReason:    "支付失败自动取消",
				OrderConfig:     `{"seed":"demo","flow":"payment_failed"}`,
			}
			if err := tx.Create(created).Error; err != nil {
				return err
			}
			orders["orderPaymentFailed1"] = created
			order = created
		}

		if order != nil {
			if err := seedPayment(tx, seedPaymentParams{
				OrderID:         order.ID,
				UserID:          user.ID,
				Method:          model.PaymentMethodWeChat,
				AmountCents:     order.TotalPriceCents,
				Currency:        model.CurrencyCNY,
				Status:          model.PaymentStatusFailed,
				ProviderTradeNo: fmt.Sprintf("WX-FAIL-%s", order.OrderNo),
				ProviderRaw:     json.RawMessage(`{"result":"failed","code":"BANK_REJECT","seed":"demo"}`),
			}); err != nil {
				return err
			}

			// Create one timeout log sample for this canceled order (idempotent by order+type)
			var existingLog model.OrderTimeoutLog
			if err := tx.Where("order_id = ? AND timeout_type = ?", order.ID, model.OrderTimeoutTypePayment).First(&existingLog).Error; err == nil {
				// already exists
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			} else {
				timeoutAt := now.Add(-30 * time.Minute)
				logEntry := model.OrderTimeoutLog{
					OrderID:     order.ID,
					TimeoutType: model.OrderTimeoutTypePayment,
					TimeoutAt:   timeoutAt,
					Action:      model.OrderTimeoutActionCanceled,
					Remark:      "演示数据：支付失败/超时触发自动取消",
				}
				logEntry.ExtJSON = `{"seed":"demo"}`
				_ = tx.Create(&logEntry).Error
			}
		}
	}

	return nil
}

func ensureOrderSlot(tx *gorm.DB, order *model.Order, item *model.ServiceItem, slot int, player *model.Player) error {
	var existingItem model.OrderItem
	if err := tx.Where("order_id = ? AND slot = ?", order.ID, slot).First(&existingItem).Error; err == nil {
		// Ensure order_players also exists for this slot+player (idempotent)
		if existingItem.PlayerID == nil && player != nil {
			_ = tx.Model(&existingItem).Update("player_id", player.ID).Error
		}
		return ensureOrderPlayerLink(tx, order.ID, existingItem.ID, player.ID)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	orderItem := model.OrderItem{
		OrderID:        order.ID,
		ItemID:         item.ID,
		Slot:           slot,
		UnitPriceCents: item.BasePriceCents,
		Quantity:       1,
		TotalCents:     item.BasePriceCents,
		CommissionRate: item.CommissionRate,
		Status:         model.OrderItemStatusMatched,
		PlayerID:       &player.ID,
	}
	orderItem.ExtJSON = `{"seed":"demo"}`
	if err := tx.Create(&orderItem).Error; err != nil {
		return err
	}

	return ensureOrderPlayerLink(tx, order.ID, orderItem.ID, player.ID)
}

func ensureOrderPlayerLink(tx *gorm.DB, orderID, orderItemID, playerID uint64) error {
	var existing model.OrderPlayer
	if err := tx.Where("order_id = ? AND order_item_id = ? AND player_id = ?", orderID, orderItemID, playerID).First(&existing).Error; err == nil {
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	link := model.OrderPlayer{
		OrderID:     orderID,
		OrderItemID: orderItemID,
		PlayerID:    playerID,
		JoinedAt:    time.Now(),
		Status:      model.OrderPlayerStatusJoined,
	}
	link.ExtJSON = `{"seed":"demo"}`
	return tx.Create(&link).Error
}

func seedUserVipState(tx *gorm.DB, now time.Time, users map[string]*model.User, vipLevels map[string]*model.VipLevel) error {
	vip1 := vipLevels["vip1"]
	vip2 := vipLevels["vip2"]
	svip := vipLevels["svip"]

	type vipSpec struct {
		UserKey            string
		VipLevel           *model.VipLevel
		VipExp             int64
		TotalRechargeCents int64
		ExpireDays         int // 0=永久, 负数=已过期
	}

	specs := []vipSpec{
		// customerB: VIP2, 即将到期 (5天)
		{UserKey: "customerB", VipLevel: vip2, VipExp: 150000, TotalRechargeCents: 100000, ExpireDays: 5},
		// customerA: VIP1, 永久有效
		{UserKey: "customerA", VipLevel: vip1, VipExp: 50000, TotalRechargeCents: 30000, ExpireDays: 0},
		// customerH: SVIP, 长期有效
		{UserKey: "customerH", VipLevel: svip, VipExp: 350000, TotalRechargeCents: 200000, ExpireDays: 60},
		// customerC: VIP1, 已过期
		{UserKey: "customerC", VipLevel: vip1, VipExp: 80000, TotalRechargeCents: 50000, ExpireDays: -15},
	}

	for _, spec := range specs {
		user := users[spec.UserKey]
		if user == nil || spec.VipLevel == nil {
			continue
		}
		unlockedAt := now.Add(-30 * 24 * time.Hour)
		updates := map[string]interface{}{
			"vip_level_id":           spec.VipLevel.ID,
			"vip_unlocked":           true,
			"vip_exp":                spec.VipExp,
			"total_recharge_cents":   spec.TotalRechargeCents,
			"vip_unlocked_at":        &unlockedAt,
			"last_monthly_coupon_at": now.Add(-35 * 24 * time.Hour),
		}
		if spec.ExpireDays != 0 {
			expireAt := now.Add(time.Duration(spec.ExpireDays) * 24 * time.Hour)
			updates["vip_expire_at"] = &expireAt
		} else {
			updates["vip_expire_at"] = nil // 永久
		}
		if err := tx.Model(&model.User{}).Where("id = ?", user.ID).Updates(updates).Error; err != nil {
			return err
		}
	}
	return nil
}

// seedOrderGroupData 创建多时段（主订单+子订单）种子数据
// 覆盖 OrderGroup 模型及 Order.GroupID 关联
func seedOrderGroupData(
	tx *gorm.DB,
	now time.Time,
	users map[string]*model.User,
	players map[string]*model.Player,
	serviceItems map[string]*model.ServiceItem,
	orders map[string]*model.Order,
) error {
	customerA := users["customerA"]
	customerB := users["customerB"]
	playerA := players["playerA"]
	playerB := players["playerB"]
	item := serviceItems["escort-lol-solo"]
	if customerA == nil || customerB == nil || playerA == nil || playerB == nil || item == nil {
		return nil
	}
	if item.GameID == nil {
		return nil
	}

	ensureGroup := func(key string, group model.OrderGroup) (*model.OrderGroup, error) {
		var existing model.OrderGroup
		if err := tx.Where("group_no = ?", group.GroupNo).First(&existing).Error; err == nil {
			orders["_orderGroup_"+key] = nil // placeholder
			return &existing, nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if err := tx.Create(&group).Error; err != nil {
			return nil, err
		}
		return &group, nil
	}

	// ========== 1. 已完成的多时段订单（3小时，全部完成） ==========
	startCompleted := now.Add(-24 * time.Hour)
	endCompleted := startCompleted.Add(3 * time.Hour)
	completedGroup, err := ensureGroup("completed", model.OrderGroup{
		GroupNo:         "OG-DEMO-COMPLETED-001",
		UserID:          customerA.ID,
		GameID:          *item.GameID,
		ItemID:          item.ID,
		OriginalPlayer:  playerA.ID,
		TotalPriceCents: item.BasePriceCents * 3,
		TotalHours:      3,
		CompletedHours:  3,
		Status:          model.OrderGroupStatusCompleted,
		Title:           "英雄联盟3小时连续陪玩（已完成）",
		Description:     "演示数据：多时段订单完整完成场景",
		ScheduledStart:  &startCompleted,
		ScheduledEnd:    &endCompleted,
		Currency:        model.CurrencyCNY,
	})
	if err != nil {
		return err
	}

	// 为已完成group创建3个子订单
	for i := 0; i < 3; i++ {
		hourStart := startCompleted.Add(time.Duration(i) * time.Hour)
		hourEnd := hourStart.Add(1 * time.Hour)
		title := fmt.Sprintf("英雄联盟陪玩-第%d小时（演示）", i+1)
		order, err := seedOrder(tx, seedOrderParams{
			Title:          title,
			Description:    fmt.Sprintf("演示数据：多时段订单第%d小时子订单", i+1),
			UserID:         customerA.ID,
			PlayerID:       &playerA.ID,
			ItemID:         item.ID,
			GameID:         *item.GameID,
			Status:         model.OrderStatusCompleted,
			PriceCents:     item.BasePriceCents,
			Currency:       model.CurrencyCNY,
			ScheduledStart: &hourStart,
			ScheduledEnd:   &hourEnd,
			StartedAt:      &hourStart,
			CompletedAt:    &hourEnd,
		})
		if err != nil {
			return err
		}
		if err := seedPayment(tx, seedPaymentParams{
			OrderID:         order.ID,
			UserID:          customerA.ID,
			Method:          model.PaymentMethodWeChat,
			AmountCents:     item.BasePriceCents,
			Currency:        model.CurrencyCNY,
			Status:          model.PaymentStatusPaid,
			ProviderTradeNo: fmt.Sprintf("WX-DEMO-OG-COMP-%02d", i+1),
			ProviderRaw:     json.RawMessage(`{"seed":"order_group","group_no":"OG-DEMO-COMPLETED-001"}`),
			PaidAt:          ptrTime(hourStart.Add(-5 * time.Minute)),
		}); err != nil {
			return err
		}
		// Link to group (update GroupID)
		tx.Model(&model.Order{}).Where("title = ? AND user_id = ?", title, customerA.ID).
			Update("group_id", completedGroup.ID)
	}

	// ========== 2. 进行中的多时段订单（4小时，完成2小时） ==========
	startProgress := now.Add(-2 * time.Hour)
	endProgress := startProgress.Add(4 * time.Hour)
	progressGroup, err := ensureGroup("progress", model.OrderGroup{
		GroupNo:         "OG-DEMO-INPROGRESS-001",
		UserID:          customerB.ID,
		GameID:          *item.GameID,
		ItemID:          item.ID,
		OriginalPlayer:  playerA.ID,
		TotalPriceCents: item.BasePriceCents * 4,
		TotalHours:      4,
		CompletedHours:  2,
		Status:          model.OrderGroupStatusInProgress,
		Title:           "英雄联盟4小时陪练（进行中）",
		Description:     "演示数据：多时段订单部分完成场景",
		ScheduledStart:  &startProgress,
		ScheduledEnd:    &endProgress,
		Currency:        model.CurrencyCNY,
	})
	if err != nil {
		return err
	}

	// 已完成的2个子订单 + 1个进行中 + 1个待处理
	for i := 0; i < 4; i++ {
		hourStart := startProgress.Add(time.Duration(i) * time.Hour)
		hourEnd := hourStart.Add(1 * time.Hour)
		title := fmt.Sprintf("英雄联盟陪练-第%d小时（演示）", i+1)
		var status model.OrderStatus
		var startedAt, completedAt *time.Time
		switch {
		case i < 2:
			status = model.OrderStatusCompleted
			startedAt = &hourStart
			completedAt = &hourEnd
		case i == 2:
			status = model.OrderStatusInProgress
			startedAt = &hourStart
		default:
			status = model.OrderStatusPending
		}
		order, err := seedOrder(tx, seedOrderParams{
			Title:          title,
			Description:    fmt.Sprintf("演示数据：多时段订单第%d小时", i+1),
			UserID:         customerB.ID,
			PlayerID:       &playerA.ID,
			ItemID:         item.ID,
			GameID:         *item.GameID,
			Status:         status,
			PriceCents:     item.BasePriceCents,
			Currency:       model.CurrencyCNY,
			ScheduledStart: &hourStart,
			ScheduledEnd:   &hourEnd,
			StartedAt:      startedAt,
			CompletedAt:    completedAt,
		})
		if err != nil {
			return err
		}
		if status == model.OrderStatusCompleted {
			if err := seedPayment(tx, seedPaymentParams{
				OrderID:         order.ID,
				UserID:          customerB.ID,
				Method:          model.PaymentMethodWeChat,
				AmountCents:     item.BasePriceCents,
				Currency:        model.CurrencyCNY,
				Status:          model.PaymentStatusPaid,
				ProviderTradeNo: fmt.Sprintf("WX-DEMO-OG-PROG-%02d", i+1),
				ProviderRaw:     json.RawMessage(`{"seed":"order_group","group_no":"OG-DEMO-INPROGRESS-001"}`),
				PaidAt:          ptrTime(hourStart.Add(-5 * time.Minute)),
			}); err != nil {
				return err
			}
		}
		tx.Model(&model.Order{}).Where("title = ? AND user_id = ?", title, customerB.ID).
			Update("group_id", progressGroup.ID)
	}

	// ========== 3. 部分完成（含转单）的多时段订单 ==========
	startPartial := now.Add(-48 * time.Hour)
	endPartial := startPartial.Add(3 * time.Hour)
	partialGroup, err := ensureGroup("partial", model.OrderGroup{
		GroupNo:         "OG-DEMO-PARTIAL-001",
		UserID:          customerA.ID,
		GameID:          *item.GameID,
		ItemID:          item.ID,
		OriginalPlayer:  playerB.ID,
		TotalPriceCents: item.BasePriceCents * 3,
		TotalHours:      3,
		CompletedHours:  1,
		Status:          model.OrderGroupStatusPartial,
		Title:           "英雄联盟3小时陪玩（部分完成-含转单）",
		Description:     "演示数据：陪玩师第2小时后转单给其他陪玩师",
		ScheduledStart:  &startPartial,
		ScheduledEnd:    &endPartial,
		Currency:        model.CurrencyCNY,
	})
	if err != nil {
		return err
	}

	// 第1小时完成，第2小时取消（转单），第3小时由新陪玩师完成
	partialSubs := []struct {
		HourIndex int
		Status    model.OrderStatus
		PlayerKey string
		Note      string
	}{
		{0, model.OrderStatusCompleted, "playerB", "第1小时正常完成"},
		{1, model.OrderStatusCanceled, "playerB", "第2小时陪玩师有事取消"},
		{2, model.OrderStatusInProgress, "playerA", "转单给新陪玩师继续服务"},
	}
	for _, sub := range partialSubs {
		hourStart := startPartial.Add(time.Duration(sub.HourIndex) * time.Hour)
		hourEnd := hourStart.Add(1 * time.Hour)
		title := fmt.Sprintf("多时段转单-第%d小时（演示）", sub.HourIndex+1)
		player := players[sub.PlayerKey]
		if player == nil {
			continue
		}
		var startedAt, completedAt *time.Time
		cancelReason := ""
		if sub.Status == model.OrderStatusCompleted || sub.Status == model.OrderStatusInProgress {
			startedAt = &hourStart
		}
		if sub.Status == model.OrderStatusCompleted {
			completedAt = &hourEnd
		}
		if sub.Status == model.OrderStatusCanceled {
			cancelReason = "陪玩师临时有事，转单处理"
		}
		order, err := seedOrder(tx, seedOrderParams{
			Title:          title,
			Description:    "演示数据：" + sub.Note,
			UserID:         customerA.ID,
			PlayerID:       &player.ID,
			ItemID:         item.ID,
			GameID:         *item.GameID,
			Status:         sub.Status,
			PriceCents:     item.BasePriceCents,
			Currency:       model.CurrencyCNY,
			ScheduledStart: &hourStart,
			ScheduledEnd:   &hourEnd,
			StartedAt:      startedAt,
			CompletedAt:    completedAt,
			CancelReason:   cancelReason,
		})
		if err != nil {
			return err
		}
		if sub.Status == model.OrderStatusCompleted {
			if err := seedPayment(tx, seedPaymentParams{
				OrderID:         order.ID,
				UserID:          customerA.ID,
				Method:          model.PaymentMethodWeChat,
				AmountCents:     item.BasePriceCents,
				Currency:        model.CurrencyCNY,
				Status:          model.PaymentStatusPaid,
				ProviderTradeNo: fmt.Sprintf("WX-DEMO-OG-PART-%02d", sub.HourIndex+1),
				ProviderRaw:     json.RawMessage(`{"seed":"order_group","group_no":"OG-DEMO-PARTIAL-001"}`),
				PaidAt:          ptrTime(hourStart.Add(-5 * time.Minute)),
			}); err != nil {
				return err
			}
		}
		tx.Model(&model.Order{}).Where("title = ? AND user_id = ?", title, customerA.ID).
			Updates(map[string]interface{}{
				"group_id":   partialGroup.ID,
				"hour_index": sub.HourIndex,
			})
	}

	log.Println("order group seed data created successfully")
	return nil
}

func seedRefundAndTimeoutData(tx *gorm.DB, now time.Time, users map[string]*model.User) error {
	admin := users["adminA"]

	// Ensure timeout configs exist
	type cfgSeed struct {
		Key   string
		Value string
		Desc  string
	}
	cfgs := []cfgSeed{
		{Key: model.PaymentTimeoutMinutes, Value: "30", Desc: "演示：支付超时分钟数"},
		{Key: model.OrderAcceptTimeoutMinutes, Value: "30", Desc: "演示：接单超时分钟数"},
		{Key: model.AutoCancelEnabled, Value: "true", Desc: "演示：自动取消开关"},
		{Key: model.AutoRefundEnabled, Value: "true", Desc: "演示：自动退款开关"},
		{Key: model.AutoAssignServiceEnabled, Value: "true", Desc: "演示：自动分配客服开关"},
	}
	for _, cfg := range cfgs {
		var existing model.OrderTimeoutConfig
		if err := tx.Where("config_key = ?", cfg.Key).First(&existing).Error; err == nil {
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
		c := model.OrderTimeoutConfig{
			ConfigKey:   cfg.Key,
			ConfigValue: cfg.Value,
			Description: cfg.Desc,
		}
		c.ExtJSON = `{"seed":"demo"}`
		_ = tx.Create(&c).Error
	}

	// For each refunded payment, ensure a refund record exists and keep amounts consistent.
	var refundedPayments []model.Payment
	if err := tx.Where("status = ?", model.PaymentStatusRefunded).Find(&refundedPayments).Error; err != nil {
		return nil
	}

	for _, pay := range refundedPayments {
		if pay.RefundedAt == nil {
			refAt := now.Add(-1 * time.Hour)
			pay.RefundedAt = &refAt
		}
		refundAmount := pay.AmountCents

		if err := tx.Model(&model.Payment{}).Where("id = ?", pay.ID).Updates(map[string]interface{}{
			"refunded_amount_cents": refundAmount,
			"refunded_at":           pay.RefundedAt,
		}).Error; err != nil {
			return err
		}

		var existing model.RefundRecord
		if err := tx.Where("payment_id = ? AND order_id = ? AND amount_cents = ?", pay.ID, pay.OrderID, refundAmount).First(&existing).Error; err == nil {
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		rec := model.RefundRecord{
			PaymentID:       pay.ID,
			OrderID:         pay.OrderID,
			UserID:          pay.UserID,
			AmountCents:     refundAmount,
			Reason:          "演示数据：退款记录",
			Status:          model.RefundStatusProcessed,
			ProviderTradeNo: pay.ProviderTradeNo + "-REF",
			RefundedAt:      pay.RefundedAt,
			Note:            "seed/demo",
		}
		if admin != nil {
			rec.OperatorID = &admin.ID
		}
		rec.ExtJSON = `{"seed":"demo"}`
		if err := tx.Create(&rec).Error; err != nil {
			return err
		}

		// Also ensure order refund fields are filled.
		_ = tx.Model(&model.Order{}).Where("id = ?", pay.OrderID).Updates(map[string]interface{}{
			"refund_amount_cents": refundAmount,
			"refunded_at":         pay.RefundedAt,
			"refund_reason":       "演示数据：退款",
		}).Error

		// Create one timeout log sample per refunded order (idempotent by order+type)
		var existingLog model.OrderTimeoutLog
		if err := tx.Where("order_id = ? AND timeout_type = ?", pay.OrderID, model.OrderTimeoutTypePayment).First(&existingLog).Error; err == nil {
			continue
		}
		timeoutAt := pay.CreatedAt.Add(30 * time.Minute)
		logEntry := model.OrderTimeoutLog{
			OrderID:           pay.OrderID,
			TimeoutType:       model.OrderTimeoutTypePayment,
			TimeoutAt:         timeoutAt,
			Action:            model.OrderTimeoutActionRefunded,
			RefundAmountCents: refundAmount,
			RefundRecordID:    &rec.ID,
			Remark:            "演示数据：支付超时触发自动退款",
		}
		logEntry.ExtJSON = `{"seed":"demo"}`
		_ = tx.Create(&logEntry).Error
	}

	return nil
}

func seedOrderChatAndServiceAssignment(
	tx *gorm.DB,
	now time.Time,
	users map[string]*model.User,
	players map[string]*model.Player,
	orders map[string]*model.Order,
) error {
	order := orders["orderTeamInProgress1"]
	if order == nil {
		return nil
	}
	admin := users["adminA"]

	// Create chat group bound to order
	var group model.ChatGroup
	if err := tx.Where("group_type = ? AND related_order_id = ?", model.ChatGroupTypeOrder, order.ID).First(&group).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		creator := uint64(0)
		if p := players["playerA"]; p != nil {
			creator = p.UserID
		} else if admin != nil {
			creator = admin.ID
		}
		group = model.ChatGroup{
			GroupName:      fmt.Sprintf("订单服务群-%s（演示）", order.OrderNo),
			GroupType:      model.ChatGroupTypeOrder,
			RelatedOrderID: &order.ID,
			CreatedBy:      creator,
			MaxMembers:     50,
			IsActive:       true,
			Description:    "演示数据：订单服务群（含用户/陪玩师/客服）",
			Settings:       `{"seed":"demo"}`,
		}
		group.ExtJSON = `{"seed":"demo"}`
		if err := tx.Create(&group).Error; err != nil {
			return err
		}
	}

	// Members: order user + 2 players + admin
	addMember := func(userID uint64, role model.ChatMemberRole, nickname string) error {
		var existing model.ChatGroupMember
		if err := tx.Where("group_id = ? AND user_id = ?", group.ID, userID).First(&existing).Error; err == nil {
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		m := model.ChatGroupMember{
			GroupID:  group.ID,
			UserID:   userID,
			Role:     role,
			Nickname: nickname,
			JoinedAt: now.Add(-20 * time.Minute),
			IsActive: true,
		}
		m.ExtJSON = `{"seed":"demo"}`
		return tx.Create(&m).Error
	}

	_ = addMember(order.UserID, model.ChatMemberRoleMember, "下单用户")
	if p := players["playerA"]; p != nil {
		_ = addMember(p.UserID, model.ChatMemberRoleMember, "陪玩师A")
	}
	if p := players["playerC"]; p != nil {
		_ = addMember(p.UserID, model.ChatMemberRoleMember, "陪玩师C")
	}
	if admin != nil {
		_ = addMember(admin.ID, model.ChatMemberRoleAdmin, "客服")
	}

	// Seed a few messages (idempotent by group+content+sender)
	seedMsg := func(senderID uint64, content string, audit model.ChatMessageAuditStatus) {
		var existing model.ChatMessage
		if err := tx.Where("group_id = ? AND sender_id = ? AND content = ?", group.ID, senderID, content).First(&existing).Error; err == nil {
			return
		}
		msg := model.ChatMessage{
			GroupID:     group.ID,
			SenderID:    senderID,
			Content:     content,
			MessageType: model.ChatMessageTypeText,
			AuditStatus: audit,
			Metadata:    `{"seed":"demo"}`,
		}
		msg.ExtJSON = `{"seed":"demo"}`
		_ = tx.Create(&msg).Error
	}
	seedMsg(order.UserID, "大家好，我这单想练一下打野节奏。", model.ChatMessageAuditApproved)
	if admin != nil {
		seedMsg(admin.ID, "收到，我已加入群聊，如有问题随时@我。", model.ChatMessageAuditApproved)
	}

	// Ensure service assignment exists (idempotent by order_id)
	if admin == nil {
		return nil
	}
	var assignment model.OrderServiceAssignment
	if err := tx.Where("order_id = ?", order.ID).First(&assignment).Error; err == nil {
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	assignedAt := now.Add(-25 * time.Minute)
	assignment = model.OrderServiceAssignment{
		OrderID:       order.ID,
		ServiceUserID: admin.ID,
		ChatGroupID:   &group.ID,
		Status:        model.ServiceAssignmentStatusJoined,
		AssignedAt:    assignedAt,
		AssignType:    "auto",
		Remark:        "演示数据：自动分配客服并加入订单群聊",
	}
	assignment.ExtJSON = `{"seed":"demo"}`
	return tx.Create(&assignment).Error
}

// validateSeedAssociations performs a best-effort relational sanity check for demo data.
// It intentionally focuses on cross-table links that are easy to break during seed maintenance.
func validateSeedAssociations(tx *gorm.DB) error {
	isDemoSeed := func(extJSON string) bool {
		return strings.Contains(extJSON, `"seed":"demo"`)
	}

	// 1) payment.user_id must match order.user_id
	var payments []model.Payment
	if err := tx.Find(&payments).Error; err != nil {
		return fmt.Errorf("seed validation: query payments failed: %w", err)
	}
	for _, p := range payments {
		var order model.Order
		if err := tx.Select("id,user_id").Where("id = ?", p.OrderID).First(&order).Error; err != nil {
			return fmt.Errorf("seed validation: payment %d references missing order %d", p.ID, p.OrderID)
		}
		if order.UserID != p.UserID {
			return fmt.Errorf("seed validation: payment %d user_id %d != order.user_id %d", p.ID, p.UserID, order.UserID)
		}
		if p.Method == model.PaymentMethodCombined {
			if p.WalletAmountCents+p.ThirdPartyAmountCents != p.AmountCents {
				return fmt.Errorf("seed validation: combined payment %d amounts do not sum", p.ID)
			}
		}
		if p.Method == model.PaymentMethodWallet && p.WalletAmountCents != 0 && p.WalletAmountCents != p.AmountCents {
			return fmt.Errorf("seed validation: wallet payment %d wallet_amount_cents %d != amount_cents %d", p.ID, p.WalletAmountCents, p.AmountCents)
		}
	}

	// 2) review.user_id must match order.user_id; review.player_id should match order.player_id when present
	var reviews []model.Review
	if err := tx.Find(&reviews).Error; err != nil {
		return fmt.Errorf("seed validation: query reviews failed: %w", err)
	}
	for _, r := range reviews {
		var order model.Order
		if err := tx.Select("id,user_id,player_id").Where("id = ?", r.OrderID).First(&order).Error; err != nil {
			return fmt.Errorf("seed validation: review %d references missing order %d", r.ID, r.OrderID)
		}
		log.Printf("validate: review id=%d (user_id=%d, order_id=%d) <- order id=%d (user_id=%d)", r.ID, r.UserID, r.OrderID, order.ID, order.UserID)
		if order.UserID != r.UserID {
			return fmt.Errorf("seed validation: review %d user_id %d != order.user_id %d", r.ID, r.UserID, order.UserID)
		}
		if order.PlayerID != nil && *order.PlayerID != r.PlayerID {
			return fmt.Errorf("seed validation: review %d player_id %d != order.player_id %d", r.ID, r.PlayerID, *order.PlayerID)
		}
	}

	// 3) gift orders must use gift service item
	var giftOrders []model.Order
	if err := tx.Where("recipient_player_id IS NOT NULL").Find(&giftOrders).Error; err != nil {
		return fmt.Errorf("seed validation: query gift orders failed: %w", err)
	}
	for _, o := range giftOrders {
		var item model.ServiceItem
		if err := tx.Select("id,sub_category").Where("id = ?", o.ItemID).First(&item).Error; err != nil {
			return fmt.Errorf("seed validation: gift order %d missing service item %d", o.ID, o.ItemID)
		}
		if item.SubCategory != model.SubCategoryGift {
			return fmt.Errorf("seed validation: gift order %d item %d subCategory=%s", o.ID, item.ID, item.SubCategory)
		}
	}

	// 4) team orders should have order_items/order_players
	var teamOrders []model.Order
	if err := tx.Where("required_players > 1").Find(&teamOrders).Error; err != nil {
		return fmt.Errorf("seed validation: query team orders failed: %w", err)
	}
	for _, o := range teamOrders {
		var itemCount int64
		if err := tx.Model(&model.OrderItem{}).Where("order_id = ?", o.ID).Count(&itemCount).Error; err != nil {
			return fmt.Errorf("seed validation: query team order %d items failed: %w", o.ID, err)
		}
		if int(itemCount) < o.RequiredPlayers {
			return fmt.Errorf("seed validation: team order %d has %d order_items < required_players=%d", o.ID, itemCount, o.RequiredPlayers)
		}
		var playerCount int64
		if err := tx.Model(&model.OrderPlayer{}).Where("order_id = ?", o.ID).Count(&playerCount).Error; err != nil {
			return fmt.Errorf("seed validation: query team order %d players failed: %w", o.ID, err)
		}
		if int(playerCount) < o.RequiredPlayers {
			return fmt.Errorf("seed validation: team order %d has %d order_players < required_players=%d", o.ID, playerCount, o.RequiredPlayers)
		}
	}

	// 5) activity participation coupons should exist (if coupon_ids is JSON array)
	var parts []model.ActivityParticipation
	if err := tx.Find(&parts).Error; err != nil {
		return fmt.Errorf("seed validation: query activity participations failed: %w", err)
	}
	for _, p := range parts {
		if p.CouponIDs == "" {
			continue
		}
		var ids []uint64
		if err := json.Unmarshal([]byte(p.CouponIDs), &ids); err != nil {
			continue
		}
		for _, id := range ids {
			var c model.Coupon
			if err := tx.Select("id").Where("id = ?", id).First(&c).Error; err != nil {
				return fmt.Errorf("seed validation: activity participation %d references missing coupon %d", p.ID, id)
			}
		}
	}

	// 6) order_items should reference valid order/item/player/review
	var orderItems []model.OrderItem
	if err := tx.Find(&orderItems).Error; err != nil {
		return fmt.Errorf("seed validation: query order_items failed: %w", err)
	}
	for _, oi := range orderItems {
		var order model.Order
		if err := tx.Select("id").Where("id = ?", oi.OrderID).First(&order).Error; err != nil {
			return fmt.Errorf("seed validation: order_item %d references missing order %d", oi.ID, oi.OrderID)
		}
		var item model.ServiceItem
		if err := tx.Select("id").Where("id = ?", oi.ItemID).First(&item).Error; err != nil {
			return fmt.Errorf("seed validation: order_item %d references missing service item %d", oi.ID, oi.ItemID)
		}
		if oi.PlayerID != nil {
			var player model.Player
			if err := tx.Select("id").Where("id = ?", *oi.PlayerID).First(&player).Error; err != nil {
				return fmt.Errorf("seed validation: order_item %d references missing player %d", oi.ID, *oi.PlayerID)
			}
		}
		if oi.ReviewID != nil {
			var review model.Review
			if err := tx.Select("id").Where("id = ?", *oi.ReviewID).First(&review).Error; err != nil {
				return fmt.Errorf("seed validation: order_item %d references missing review %d", oi.ID, *oi.ReviewID)
			}
		}
	}

	// 7) order_players should reference valid order + order_item (same order) + player
	var orderPlayers []model.OrderPlayer
	if err := tx.Find(&orderPlayers).Error; err != nil {
		return fmt.Errorf("seed validation: query order_players failed: %w", err)
	}
	for _, op := range orderPlayers {
		var order model.Order
		if err := tx.Select("id").Where("id = ?", op.OrderID).First(&order).Error; err != nil {
			return fmt.Errorf("seed validation: order_player %d references missing order %d", op.ID, op.OrderID)
		}
		var orderItem model.OrderItem
		if err := tx.Select("id,order_id").Where("id = ?", op.OrderItemID).First(&orderItem).Error; err != nil {
			return fmt.Errorf("seed validation: order_player %d references missing order_item %d", op.ID, op.OrderItemID)
		}
		if orderItem.OrderID != op.OrderID {
			return fmt.Errorf("seed validation: order_player %d order_id %d != order_item.order_id %d", op.ID, op.OrderID, orderItem.OrderID)
		}
		var player model.Player
		if err := tx.Select("id").Where("id = ?", op.PlayerID).First(&player).Error; err != nil {
			return fmt.Errorf("seed validation: order_player %d references missing player %d", op.ID, op.PlayerID)
		}
	}

	// 8) coupons should reference valid template/user/orders when present
	var coupons []model.Coupon
	if err := tx.Find(&coupons).Error; err != nil {
		return fmt.Errorf("seed validation: query coupons failed: %w", err)
	}
	for _, c := range coupons {
		var tpl model.CouponTemplate
		if err := tx.Select("id").Where("id = ?", c.TemplateID).First(&tpl).Error; err != nil {
			return fmt.Errorf("seed validation: coupon %d references missing template %d", c.ID, c.TemplateID)
		}
		var user model.User
		if err := tx.Select("id").Where("id = ?", c.UserID).First(&user).Error; err != nil {
			return fmt.Errorf("seed validation: coupon %d references missing user %d", c.ID, c.UserID)
		}
		if c.LockedByOrderID != nil {
			var order model.Order
			if err := tx.Select("id,user_id").Where("id = ?", *c.LockedByOrderID).First(&order).Error; err != nil {
				return fmt.Errorf("seed validation: coupon %d locked_by_order_id missing order %d", c.ID, *c.LockedByOrderID)
			}
			if order.UserID != c.UserID {
				return fmt.Errorf("seed validation: coupon %d locked_by_order_id order.user_id %d != coupon.user_id %d", c.ID, order.UserID, c.UserID)
			}
		}
		if c.UsedOrderID != nil {
			var order model.Order
			if err := tx.Select("id,user_id").Where("id = ?", *c.UsedOrderID).First(&order).Error; err != nil {
				return fmt.Errorf("seed validation: coupon %d used_order_id missing order %d", c.ID, *c.UsedOrderID)
			}
			if order.UserID != c.UserID {
				return fmt.Errorf("seed validation: coupon %d used_order_id order.user_id %d != coupon.user_id %d", c.ID, order.UserID, c.UserID)
			}
		}
	}

	// 9) team demo data: ensure basic foreign keys
	var teams []model.Team
	if err := tx.Find(&teams).Error; err != nil {
		return fmt.Errorf("seed validation: query teams failed: %w", err)
	}
	for _, t := range teams {
		if !isDemoSeed(t.ExtJSON) {
			continue
		}
		var leader model.Player
		if err := tx.Select("id").Where("id = ?", t.LeaderID).First(&leader).Error; err != nil {
			return fmt.Errorf("seed validation: team %d references missing leader player %d", t.ID, t.LeaderID)
		}
		if t.CurrentOrderID != nil {
			var order model.Order
			if err := tx.Select("id").Where("id = ?", *t.CurrentOrderID).First(&order).Error; err != nil {
				return fmt.Errorf("seed validation: team %d references missing current order %d", t.ID, *t.CurrentOrderID)
			}
		}
	}

	var teamMembers []model.TeamMember
	if err := tx.Find(&teamMembers).Error; err != nil {
		return fmt.Errorf("seed validation: query team_members failed: %w", err)
	}
	for _, tm := range teamMembers {
		if !isDemoSeed(tm.ExtJSON) {
			continue
		}
		var team model.Team
		if err := tx.Select("id").Where("id = ?", tm.TeamID).First(&team).Error; err != nil {
			return fmt.Errorf("seed validation: team_member %d references missing team %d", tm.ID, tm.TeamID)
		}
		var player model.Player
		if err := tx.Select("id").Where("id = ?", tm.PlayerID).First(&player).Error; err != nil {
			return fmt.Errorf("seed validation: team_member %d references missing player %d", tm.ID, tm.PlayerID)
		}
	}

	var teamInvites []model.TeamInvite
	if err := tx.Find(&teamInvites).Error; err != nil {
		return fmt.Errorf("seed validation: query team_invites failed: %w", err)
	}
	for _, inv := range teamInvites {
		if !isDemoSeed(inv.ExtJSON) {
			continue
		}
		var team model.Team
		if err := tx.Select("id").Where("id = ?", inv.TeamID).First(&team).Error; err != nil {
			return fmt.Errorf("seed validation: team_invite %d references missing team %d", inv.ID, inv.TeamID)
		}
		var player model.Player
		if err := tx.Select("id").Where("id = ?", inv.PlayerID).First(&player).Error; err != nil {
			return fmt.Errorf("seed validation: team_invite %d references missing player %d", inv.ID, inv.PlayerID)
		}
		var inviter model.Player
		if err := tx.Select("id").Where("id = ?", inv.InviterID).First(&inviter).Error; err != nil {
			return fmt.Errorf("seed validation: team_invite %d references missing inviter %d", inv.ID, inv.InviterID)
		}
	}

	// 10) referral demo data: ensure user/code links exist
	var referralCodes []model.ReferralCode
	if err := tx.Find(&referralCodes).Error; err != nil {
		return fmt.Errorf("seed validation: query referral_codes failed: %w", err)
	}
	for _, rc := range referralCodes {
		if !isDemoSeed(rc.ExtJSON) {
			continue
		}
		var user model.User
		if err := tx.Select("id").Where("id = ?", rc.UserID).First(&user).Error; err != nil {
			return fmt.Errorf("seed validation: referral_code %d references missing user %d", rc.ID, rc.UserID)
		}
	}

	var referrals []model.Referral
	if err := tx.Find(&referrals).Error; err != nil {
		return fmt.Errorf("seed validation: query referrals failed: %w", err)
	}
	for _, r := range referrals {
		if !isDemoSeed(r.ExtJSON) {
			continue
		}
		var referrer model.User
		if err := tx.Select("id").Where("id = ?", r.ReferrerID).First(&referrer).Error; err != nil {
			return fmt.Errorf("seed validation: referral %d references missing referrer %d", r.ID, r.ReferrerID)
		}
		var referee model.User
		if err := tx.Select("id").Where("id = ?", r.RefereeID).First(&referee).Error; err != nil {
			return fmt.Errorf("seed validation: referral %d references missing referee %d", r.ID, r.RefereeID)
		}
		if r.CodeID != nil {
			var code model.ReferralCode
			if err := tx.Select("id").Where("id = ?", *r.CodeID).First(&code).Error; err != nil {
				return fmt.Errorf("seed validation: referral %d references missing referral_code %d", r.ID, *r.CodeID)
			}
		}
	}

	var referralRewards []model.ReferralReward
	if err := tx.Find(&referralRewards).Error; err != nil {
		return fmt.Errorf("seed validation: query referral_rewards failed: %w", err)
	}
	for _, rr := range referralRewards {
		if !isDemoSeed(rr.ExtJSON) {
			continue
		}
		var ref model.Referral
		if err := tx.Select("id").Where("id = ?", rr.ReferralID).First(&ref).Error; err != nil {
			return fmt.Errorf("seed validation: referral_reward %d references missing referral %d", rr.ID, rr.ReferralID)
		}
		var user model.User
		if err := tx.Select("id").Where("id = ?", rr.UserID).First(&user).Error; err != nil {
			return fmt.Errorf("seed validation: referral_reward %d references missing user %d", rr.ID, rr.UserID)
		}
		if rr.CouponID != nil {
			var c model.Coupon
			if err := tx.Select("id").Where("id = ?", *rr.CouponID).First(&c).Error; err != nil {
				return fmt.Errorf("seed validation: referral_reward %d references missing coupon %d", rr.ID, *rr.CouponID)
			}
		}
	}

	// 11) user block demo data: ensure blocker/blocked/canceler exist
	var blocks []model.UserBlock
	if err := tx.Find(&blocks).Error; err != nil {
		return fmt.Errorf("seed validation: query user_blocks failed: %w", err)
	}
	for _, b := range blocks {
		if !isDemoSeed(b.ExtJSON) {
			continue
		}
		if b.BlockerID == b.BlockedID {
			return fmt.Errorf("seed validation: user_block %d blocker_id == blocked_id (%d)", b.ID, b.BlockerID)
		}
		var blocker model.User
		if err := tx.Select("id").Where("id = ?", b.BlockerID).First(&blocker).Error; err != nil {
			return fmt.Errorf("seed validation: user_block %d references missing blocker user %d", b.ID, b.BlockerID)
		}
		var blocked model.User
		if err := tx.Select("id").Where("id = ?", b.BlockedID).First(&blocked).Error; err != nil {
			return fmt.Errorf("seed validation: user_block %d references missing blocked user %d", b.ID, b.BlockedID)
		}
		if b.CanceledBy != nil {
			var admin model.User
			if err := tx.Select("id").Where("id = ?", *b.CanceledBy).First(&admin).Error; err != nil {
				return fmt.Errorf("seed validation: user_block %d references missing canceled_by user %d", b.ID, *b.CanceledBy)
			}
		}
	}

	// 12) recharge demo data: ensure user/option/coupon_ids exist
	var rechargeRecords []model.RechargeRecord
	if err := tx.Find(&rechargeRecords).Error; err != nil {
		return fmt.Errorf("seed validation: query recharge_records failed: %w", err)
	}
	for _, r := range rechargeRecords {
		if !isDemoSeed(r.ExtJSON) {
			continue
		}
		var user model.User
		if err := tx.Select("id").Where("id = ?", r.UserID).First(&user).Error; err != nil {
			return fmt.Errorf("seed validation: recharge_record %d references missing user %d", r.ID, r.UserID)
		}
		if r.OptionID != nil {
			var opt model.RechargeOption
			if err := tx.Select("id").Where("id = ?", *r.OptionID).First(&opt).Error; err != nil {
				return fmt.Errorf("seed validation: recharge_record %d references missing recharge_option %d", r.ID, *r.OptionID)
			}
		}
		if r.CouponIDs == "" {
			continue
		}
		var ids []uint64
		if err := json.Unmarshal([]byte(r.CouponIDs), &ids); err != nil {
			continue
		}
		for _, id := range ids {
			var c model.Coupon
			if err := tx.Select("id").Where("id = ?", id).First(&c).Error; err != nil {
				return fmt.Errorf("seed validation: recharge_record %d references missing coupon %d", r.ID, id)
			}
		}
	}

	// 13) notifications demo data: ensure user/template/related exist
	var userNotifs []model.UserNotification
	if err := tx.Find(&userNotifs).Error; err != nil {
		return fmt.Errorf("seed validation: query user_notifications failed: %w", err)
	}
	for _, n := range userNotifs {
		if !isDemoSeed(n.ExtJSON) {
			continue
		}
		var user model.User
		if err := tx.Select("id").Where("id = ?", n.UserID).First(&user).Error; err != nil {
			return fmt.Errorf("seed validation: user_notification %d references missing user %d", n.ID, n.UserID)
		}
		if n.TemplateID != nil {
			var tpl model.NotificationTemplate
			if err := tx.Select("id").Where("id = ?", *n.TemplateID).First(&tpl).Error; err != nil {
				return fmt.Errorf("seed validation: user_notification %d references missing template %d", n.ID, *n.TemplateID)
			}
		}
		if n.RelatedID == nil || n.RelatedType == "" {
			continue
		}
		switch n.RelatedType {
		case "order":
			var o model.Order
			if err := tx.Select("id").Where("id = ?", *n.RelatedID).First(&o).Error; err != nil {
				return fmt.Errorf("seed validation: user_notification %d references missing order %d", n.ID, *n.RelatedID)
			}
		case "coupon":
			var c model.Coupon
			if err := tx.Select("id").Where("id = ?", *n.RelatedID).First(&c).Error; err != nil {
				return fmt.Errorf("seed validation: user_notification %d references missing coupon %d", n.ID, *n.RelatedID)
			}
		case "activity":
			var a model.Activity
			if err := tx.Select("id").Where("id = ?", *n.RelatedID).First(&a).Error; err != nil {
				return fmt.Errorf("seed validation: user_notification %d references missing activity %d", n.ID, *n.RelatedID)
			}
		case "vip":
			var v model.VipLevel
			if err := tx.Select("id").Where("id = ?", *n.RelatedID).First(&v).Error; err != nil {
				return fmt.Errorf("seed validation: user_notification %d references missing vip_level %d", n.ID, *n.RelatedID)
			}
		}
	}

	var notifSettings []model.UserNotificationSetting
	if err := tx.Find(&notifSettings).Error; err != nil {
		return fmt.Errorf("seed validation: query user_notification_settings failed: %w", err)
	}
	for _, s := range notifSettings {
		if !isDemoSeed(s.ExtJSON) {
			continue
		}
		var user model.User
		if err := tx.Select("id").Where("id = ?", s.UserID).First(&user).Error; err != nil {
			return fmt.Errorf("seed validation: user_notification_setting %d references missing user %d", s.ID, s.UserID)
		}
	}

	// 14) chat/service demo data: ensure group/member/message/assignment links exist
	var chatGroups []model.ChatGroup
	if err := tx.Find(&chatGroups).Error; err != nil {
		return fmt.Errorf("seed validation: query chat_groups failed: %w", err)
	}
	for _, g := range chatGroups {
		if !isDemoSeed(g.ExtJSON) {
			continue
		}
		var creator model.User
		if err := tx.Select("id").Where("id = ?", g.CreatedBy).First(&creator).Error; err != nil {
			return fmt.Errorf("seed validation: chat_group %d references missing created_by user %d", g.ID, g.CreatedBy)
		}
		if g.RelatedOrderID != nil {
			var order model.Order
			if err := tx.Select("id").Where("id = ?", *g.RelatedOrderID).First(&order).Error; err != nil {
				return fmt.Errorf("seed validation: chat_group %d references missing related_order_id %d", g.ID, *g.RelatedOrderID)
			}
		}
		if g.GroupType == model.ChatGroupTypeOrder && g.RelatedOrderID == nil {
			return fmt.Errorf("seed validation: chat_group %d group_type=order but related_order_id is NULL", g.ID)
		}
	}

	var chatMembers []model.ChatGroupMember
	if err := tx.Find(&chatMembers).Error; err != nil {
		return fmt.Errorf("seed validation: query chat_group_members failed: %w", err)
	}
	for _, m := range chatMembers {
		if !isDemoSeed(m.ExtJSON) {
			continue
		}
		var group model.ChatGroup
		if err := tx.Select("id").Where("id = ?", m.GroupID).First(&group).Error; err != nil {
			return fmt.Errorf("seed validation: chat_group_member %d references missing group %d", m.ID, m.GroupID)
		}
		var user model.User
		if err := tx.Select("id").Where("id = ?", m.UserID).First(&user).Error; err != nil {
			return fmt.Errorf("seed validation: chat_group_member %d references missing user %d", m.ID, m.UserID)
		}
	}

	var chatMessages []model.ChatMessage
	if err := tx.Find(&chatMessages).Error; err != nil {
		return fmt.Errorf("seed validation: query chat_messages failed: %w", err)
	}
	for _, m := range chatMessages {
		if !isDemoSeed(m.ExtJSON) {
			continue
		}
		var group model.ChatGroup
		if err := tx.Select("id").Where("id = ?", m.GroupID).First(&group).Error; err != nil {
			return fmt.Errorf("seed validation: chat_message %d references missing group %d", m.ID, m.GroupID)
		}
		var sender model.User
		if err := tx.Select("id").Where("id = ?", m.SenderID).First(&sender).Error; err != nil {
			return fmt.Errorf("seed validation: chat_message %d references missing sender user %d", m.ID, m.SenderID)
		}
	}

	var assignments []model.OrderServiceAssignment
	if err := tx.Find(&assignments).Error; err != nil {
		return fmt.Errorf("seed validation: query order_service_assignments failed: %w", err)
	}
	for _, a := range assignments {
		if !isDemoSeed(a.ExtJSON) {
			continue
		}
		var order model.Order
		if err := tx.Select("id").Where("id = ?", a.OrderID).First(&order).Error; err != nil {
			return fmt.Errorf("seed validation: order_service_assignment %d references missing order %d", a.ID, a.OrderID)
		}
		var serviceUser model.User
		if err := tx.Select("id").Where("id = ?", a.ServiceUserID).First(&serviceUser).Error; err != nil {
			return fmt.Errorf("seed validation: order_service_assignment %d references missing service user %d", a.ID, a.ServiceUserID)
		}
		if a.ChatGroupID != nil {
			var group model.ChatGroup
			if err := tx.Select("id,related_order_id,group_type").Where("id = ?", *a.ChatGroupID).First(&group).Error; err != nil {
				return fmt.Errorf("seed validation: order_service_assignment %d references missing chat_group %d", a.ID, *a.ChatGroupID)
			}
			if group.GroupType == model.ChatGroupTypeOrder && group.RelatedOrderID != nil && *group.RelatedOrderID != a.OrderID {
				return fmt.Errorf("seed validation: order_service_assignment %d chat_group.related_order_id %d != order_id %d", a.ID, *group.RelatedOrderID, a.OrderID)
			}
		}
	}

	log.Println("seed association validation passed")
	return nil
}

// seedBanners 种子数据：首页轮播图
func seedBanners(tx *gorm.DB) error {
	banners := []model.Banner{
		{
			Title:       "探索热门游戏",
			Description: "一键发现优质陪玩师，畅享游戏乐趣",
			ImageURL:    "/static/images/banner-jump.svg",
			Type:        model.BannerTypeLink,
			Link:        "/pages/game/list/index",
			ActionText:  "立即前往",
			SortOrder:   0,
			IsVisible:   true,
		},
		{
			Title:       "新赛季展示",
			Description: "查看最新活动海报与精彩内容",
			ImageURL:    "/static/images/banner-preview.svg",
			Type:        model.BannerTypePreview,
			ActionText:  "查看详情",
			SortOrder:   1,
			IsVisible:   true,
		},
	}

	for i := range banners {
		if err := tx.Create(&banners[i]).Error; err != nil {
			return fmt.Errorf("seed banner %d: %w", i, err)
		}
	}
	log.Printf("[startup] seed: created %d banners", len(banners))
	return nil
}
