package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
)

// This file adds additional demo seed coverage for marketing/notification/team/rank modules.

func seedServiceItems(tx *gorm.DB, games map[string]*model.Game) (map[string]*model.ServiceItem, error) {
	return seedServiceItemsInternal(tx, games)
}

func seedVipAndCouponTemplates(tx *gorm.DB, games map[string]*model.Game, serviceItems map[string]*model.ServiceItem) (map[string]*model.CouponTemplate, map[string]*model.VipLevel, error) {
	return seedVipAndCouponTemplatesInternal(tx, games, serviceItems)
}

func seedUserCoupons(tx *gorm.DB, templates map[string]*model.CouponTemplate, users map[string]*model.User, orders map[string]*model.Order) (map[string]*model.Coupon, error) {
	return seedUserCouponsInternal(tx, templates, users, orders)
}

func seedRechargeData(tx *gorm.DB, templates map[string]*model.CouponTemplate, users map[string]*model.User) error {
	return seedRechargeDataInternal(tx, templates, users)
}

func seedActivityData(tx *gorm.DB, templates map[string]*model.CouponTemplate, users map[string]*model.User) (map[string]*model.Activity, error) {
	return seedActivityDataInternal(tx, templates, users)
}

func seedTeamData(tx *gorm.DB, players map[string]*model.Player, orders map[string]*model.Order) error {
	return seedTeamDataInternal(tx, players, orders)
}

func seedReferralData(tx *gorm.DB, templates map[string]*model.CouponTemplate, users map[string]*model.User) error {
	return seedReferralDataInternal(tx, templates, users)
}

func seedUserBlockData(tx *gorm.DB, users map[string]*model.User) error {
	return seedUserBlockDataInternal(tx, users)
}

func seedGameRankAndCertificationData(tx *gorm.DB, games map[string]*model.Game, players map[string]*model.Player, users map[string]*model.User) error {
	return seedGameRankAndCertificationDataInternal(tx, games, players, users)
}

func seedNotificationData(tx *gorm.DB, users map[string]*model.User, orders map[string]*model.Order, coupons map[string]*model.Coupon, activities map[string]*model.Activity, vipLevels map[string]*model.VipLevel) error {
	return seedNotificationDataInternal(tx, users, orders, coupons, activities, vipLevels)
}

func couponFromTemplate(userID uint64, tpl *model.CouponTemplate, state model.CouponState, claimedAt *time.Time, expireAt time.Time) model.Coupon {
	c := model.Coupon{
		TemplateID: tpl.ID,
		UserID:     userID,
		State:      state,
		Name:       tpl.Name,
		Type:       tpl.Type,
		Source:     tpl.Source,
		MinAmountCents:    tpl.MinAmountCents,
		DeductAmountCents: tpl.DeductAmountCents,
		DiscountRate:      tpl.DiscountRate,
		MaxDiscountCents:  tpl.MaxDiscountCents,
		Scope:             tpl.Scope,
		GameIDs:           tpl.GameIDs,
		ItemIDs:           tpl.ItemIDs,
		ClaimedAt:         claimedAt,
		ExpireAt:          expireAt,
	}
	c.ExtJSON = `{}`
	return c
}

// ---- internal implementations (split for patch size) ----

func seedServiceItemsInternal(tx *gorm.DB, games map[string]*model.Game) (map[string]*model.ServiceItem, error) {
	defaultGame := games["lol"]
	if defaultGame == nil {
		return nil, errors.New("seed service items missing game lol")
	}

	items := make(map[string]*model.ServiceItem)

	defaultItem, err := ensureServiceItem(tx, "escort-default", "默认护航服务", defaultGame.ID)
	if err != nil {
		return nil, err
	}
	items[defaultItem.ItemCode] = defaultItem

	type itemSeed struct {
		Code            string
		Name            string
		GameKey         string
		SubCategory     model.ServiceItemSubCategory
		BasePriceCents  int64
		ServiceHours    int
		RequiredPlayers int
		MaxPlayers      int
		IsActive        bool
		VipPriceCents   *int64
		SortOrder       int
		Tag             string
	}

	vipPrice := int64(8900)
	seeds := []itemSeed{
		{Code: "escort-lol-solo", Name: "英雄联盟单人护航", GameKey: "lol", SubCategory: model.SubCategorySolo, BasePriceCents: 9900, ServiceHours: 1, RequiredPlayers: 1, MaxPlayers: 1, IsActive: true, VipPriceCents: &vipPrice, SortOrder: 10, Tag: "热门"},
		{Code: "escort-dota2-solo", Name: "DOTA2单人护航", GameKey: "dota2", SubCategory: model.SubCategorySolo, BasePriceCents: 10900, ServiceHours: 1, RequiredPlayers: 1, MaxPlayers: 1, IsActive: true, SortOrder: 20, Tag: "推荐"},
		{Code: "escort-valorant-solo", Name: "无畏契约单人陪练", GameKey: "valorant", SubCategory: model.SubCategorySolo, BasePriceCents: 12900, ServiceHours: 1, RequiredPlayers: 1, MaxPlayers: 1, IsActive: true, SortOrder: 30, Tag: "上分"},
		{Code: "escort-wow-solo", Name: "魔兽世界副本教学", GameKey: "wow", SubCategory: model.SubCategorySolo, BasePriceCents: 15900, ServiceHours: 2, RequiredPlayers: 1, MaxPlayers: 1, IsActive: true, SortOrder: 40, Tag: "新手"},
		{Code: "escort-lol-team", Name: "英雄联盟双人车队", GameKey: "lol", SubCategory: model.SubCategoryTeam, BasePriceCents: 19900, ServiceHours: 1, RequiredPlayers: 2, MaxPlayers: 2, IsActive: true, SortOrder: 50, Tag: "组队"},
		{Code: "gift-rose", Name: "礼物：玫瑰", GameKey: "", SubCategory: model.SubCategoryGift, BasePriceCents: 990, ServiceHours: 0, RequiredPlayers: 1, MaxPlayers: 1, IsActive: true, SortOrder: 60, Tag: "礼物"},
		{Code: "gift-sponsor", Name: "礼物：打赏", GameKey: "", SubCategory: model.SubCategoryGift, BasePriceCents: 5000, ServiceHours: 0, RequiredPlayers: 1, MaxPlayers: 1, IsActive: false, SortOrder: 70, Tag: "下架"},
	}

	now := time.Now()
	for _, seed := range seeds {
		var existing model.ServiceItem
		if err := tx.Where("item_code = ?", seed.Code).First(&existing).Error; err == nil {
			ex := existing
			items[seed.Code] = &ex
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		var gameID *uint64
		if seed.GameKey != "" {
			game := games[seed.GameKey]
			if game == nil {
				return nil, fmt.Errorf("seed service item %s missing game %s", seed.Code, seed.GameKey)
			}
			gameID = &game.ID
		}

		item := model.ServiceItem{
			ItemCode:       seed.Code,
			Name:           seed.Name,
			Description:    "演示数据：用于本地/测试环境展示不同服务类型与配置项",
			Category:       "escort",
			SubCategory:    seed.SubCategory,
			GameID:         gameID,
			BasePriceCents: seed.BasePriceCents,
			ServiceHours:   seed.ServiceHours,
			CommissionRate: 0.20,
			MinUsers:       1,
			MaxPlayers:     seed.MaxPlayers,
			Tags:           "[]",
			IconURL:        "https://example.com/icon.png",
			IsActive:       seed.IsActive,
			SortOrder:      seed.SortOrder,
			CreatedAt:      now,
			UpdatedAt:      now,

			RequiredPlayers: seed.RequiredPlayers,
			VipPriceCents:   seed.VipPriceCents,

			UsageLimitType:  model.UsageLimitNone,
			UsageLimitCount: 0,
			MaxPerOrder:     0,
		}
		item.Tags = fmt.Sprintf(`["seed:demo","tag:%s"]`, seed.Tag)

		if err := tx.Create(&item).Error; err != nil {
			return nil, err
		}
		items[seed.Code] = &item
	}

	log.Printf("service items seeded: %d\n", len(items))
	return items, nil
}

func seedVipAndCouponTemplatesInternal(tx *gorm.DB, games map[string]*model.Game, serviceItems map[string]*model.ServiceItem) (map[string]*model.CouponTemplate, map[string]*model.VipLevel, error) {
	templates := make(map[string]*model.CouponTemplate)

	ensureTemplate := func(key string, tpl model.CouponTemplate) (*model.CouponTemplate, error) {
		if tpl.ClaimLink == "" {
			return nil, errors.New("coupon template requires claim link for seeding")
		}

		var existing model.CouponTemplate
		if err := tx.Where("claim_link = ?", tpl.ClaimLink).First(&existing).Error; err == nil {
			ex := existing
			templates[key] = &ex
			return &ex, nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		tpl.ExtJSON = `{"seed":"demo"}`
		if tpl.GameIDs == "" {
			tpl.GameIDs = "[]"
		}
		if tpl.ItemIDs == "" {
			tpl.ItemIDs = "[]"
		}
		if err := tx.Create(&tpl).Error; err != nil {
			return nil, err
		}
		templates[key] = &tpl
		return &tpl, nil
	}

	_, err := ensureTemplate("new_user_deduct", model.CouponTemplate{
		Name:              "新用户立减券（演示）",
		Type:              model.CouponTypeDeduct,
		Source:            model.CouponSourceNewUser,
		Description:       "演示数据：注册后可领取",
		MinAmountCents:    0,
		DeductAmountCents: 1000,
		Scope:             model.CouponScopeAll,
		ValidityType:      "days",
		ValidityDays:      30,
		PerUserLimit:      1,
		ClaimLink:         "seed-new-user-deduct",
		IsActive:          true,
	})
	if err != nil {
		return nil, nil, err
	}

	_, err = ensureTemplate("vip_monthly_discount", model.CouponTemplate{
		Name:             "VIP月度折扣券（演示）",
		Type:             model.CouponTypeDiscount,
		Source:           model.CouponSourceVip,
		Description:      "演示数据：VIP每月发放",
		MinAmountCents:   0,
		DiscountRate:     0.90,
		MaxDiscountCents: 3000,
		Scope:            model.CouponScopeAll,
		ValidityType:     "days",
		ValidityDays:     30,
		PerUserLimit:     99,
		ClaimLink:        "seed-vip-monthly-discount",
		IsActive:         true,
	})
	if err != nil {
		return nil, nil, err
	}

	_, err = ensureTemplate("recharge_deduct", model.CouponTemplate{
		Name:              "充值赠送满减券（演示）",
		Type:              model.CouponTypeDeduct,
		Source:            model.CouponSourceRecharge,
		Description:       "演示数据：充值档位赠送",
		MinAmountCents:    5000,
		DeductAmountCents: 500,
		Scope:             model.CouponScopeAll,
		ValidityType:      "days",
		ValidityDays:      15,
		PerUserLimit:      99,
		ClaimLink:         "seed-recharge-deduct",
		IsActive:          true,
	})
	if err != nil {
		return nil, nil, err
	}

	_, err = ensureTemplate("activity_deduct", model.CouponTemplate{
		Name:              "活动奖励满减券（演示）",
		Type:              model.CouponTypeDeduct,
		Source:            model.CouponSourceActivity,
		Description:       "演示数据：活动参与奖励",
		MinAmountCents:    10000,
		DeductAmountCents: 1500,
		Scope:             model.CouponScopeAll,
		ValidityType:      "days",
		ValidityDays:      7,
		PerUserLimit:      99,
		ClaimLink:         "seed-activity-deduct",
		IsActive:          true,
	})
	if err != nil {
		return nil, nil, err
	}

	// 游戏限定券（覆盖 CouponScopeGame）
	if lol := games["lol"]; lol != nil {
		_, err = ensureTemplate("lol_only_discount", model.CouponTemplate{
			Name:             "英雄联盟专属折扣券（演示）",
			Type:             model.CouponTypeDiscount,
			Source:           model.CouponSourceManual,
			Description:      "演示数据：仅限英雄联盟下单使用",
			MinAmountCents:   0,
			DiscountRate:     0.85,
			MaxDiscountCents: 5000,
			Scope:            model.CouponScopeGame,
			GameIDs:          fmt.Sprintf("[%d]", lol.ID),
			ValidityType:     "days",
			ValidityDays:     30,
			PerUserLimit:     1,
			ClaimLink:        "seed-game-lol-only",
			IsActive:         true,
		})
		if err != nil {
			return nil, nil, err
		}
	}

	// 服务项限定券（覆盖 CouponScopeItem）
	if teamItem := serviceItems["escort-lol-team"]; teamItem != nil {
		_, err = ensureTemplate("team_item_deduct", model.CouponTemplate{
			Name:              "双人车队满减券（演示）",
			Type:              model.CouponTypeDeduct,
			Source:            model.CouponSourceManual,
			Description:       "演示数据：仅限双人车队服务使用",
			MinAmountCents:    10000,
			DeductAmountCents: 1200,
			Scope:             model.CouponScopeItem,
			ItemIDs:           fmt.Sprintf("[%d]", teamItem.ID),
			ValidityType:      "days",
			ValidityDays:      10,
			PerUserLimit:      1,
			ClaimLink:         "seed-item-team-only",
			IsActive:          true,
		})
		if err != nil {
			return nil, nil, err
		}
	}

	// VIP 全局配置（确保存在）
	vipConfigs := []model.VipConfig{
		{ConfigKey: model.VipConfigUnlockByConsume, ConfigValue: "0", Description: "演示：消费达到门槛解锁（分）"},
		{ConfigKey: model.VipConfigUnlockByRecharge, ConfigValue: "0", Description: "演示：充值达到门槛解锁（分）"},
		{ConfigKey: model.VipConfigExpireDays, ConfigValue: "0", Description: "演示：VIP有效期（天），0=永久"},
	}
	for _, cfg := range vipConfigs {
		var existing model.VipConfig
		if err := tx.Where("config_key = ?", cfg.ConfigKey).First(&existing).Error; err == nil {
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, err
		}
		cfg.ExtJSON = `{"seed":"demo"}`
		if err := tx.Create(&cfg).Error; err != nil {
			return nil, nil, err
		}
	}

	levels := make(map[string]*model.VipLevel)
	vipMonthlyTpl := templates["vip_monthly_discount"]
	if vipMonthlyTpl == nil {
		return nil, nil, errors.New("vip monthly coupon template missing")
	}

	levelSeeds := []model.VipLevel{
		{
			Slug:                    "vip1",
			Title:                   "VIP 1",
			ExpRequired:             0,
			OrderDiscount:           0.98,
			MonthlyCouponTemplateID: &vipMonthlyTpl.ID,
			MonthlyCouponCount:      1,
			IconURL:                 "https://example.com/vip/vip1.png",
			Color:                   "#FFD700",
			Benefits:                `{"badge":"VIP1","monthlyCoupon":true}`,
			SortOrder:               10,
			IsDefault:               true,
			IsActive:                true,
		},
		{
			Slug:                    "vip2",
			Title:                   "VIP 2",
			ExpRequired:             100000,
			OrderDiscount:           0.95,
			MonthlyCouponTemplateID: &vipMonthlyTpl.ID,
			MonthlyCouponCount:      2,
			IconURL:                 "https://example.com/vip/vip2.png",
			Color:                   "#FF8C00",
			Benefits:                `{"badge":"VIP2","monthlyCoupon":true}`,
			SortOrder:               20,
			IsDefault:               false,
			IsActive:                true,
		},
		{
			Slug:                    "svip",
			Title:                   "SVIP",
			ExpRequired:             300000,
			OrderDiscount:           0.90,
			MonthlyCouponTemplateID: &vipMonthlyTpl.ID,
			MonthlyCouponCount:      3,
			IconURL:                 "https://example.com/vip/svip.png",
			Color:                   "#C0C0C0",
			Benefits:                `{"badge":"SVIP","monthlyCoupon":true}`,
			SortOrder:               30,
			IsDefault:               false,
			IsActive:                true,
		},
	}

	for i := range levelSeeds {
		seed := &levelSeeds[i]
		var existing model.VipLevel
		if err := tx.Where("slug = ?", seed.Slug).First(&existing).Error; err == nil {
			ex := existing
			levels[seed.Slug] = &ex
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, err
		}
		seed.ExtJSON = `{"seed":"demo"}`
		if err := tx.Create(seed).Error; err != nil {
			return nil, nil, err
		}
		levels[seed.Slug] = seed
	}

	log.Printf("vip levels ensured: %d, coupon templates ensured: %d\n", len(levels), len(templates))
	return templates, levels, nil
}

func seedUserCouponsInternal(tx *gorm.DB, templates map[string]*model.CouponTemplate, users map[string]*model.User, orders map[string]*model.Order) (map[string]*model.Coupon, error) {
	result := make(map[string]*model.Coupon)

	ensureCoupon := func(key string, c model.Coupon) (*model.Coupon, error) {
		var existing model.Coupon
		q := tx.Where(
			"template_id = ? AND user_id = ? AND state = ? AND name = ? AND expire_at = ?",
			c.TemplateID, c.UserID, c.State, c.Name, c.ExpireAt,
		)
		if c.LockedByOrderID != nil {
			q = q.Where("locked_by_order_id = ?", *c.LockedByOrderID)
		}
		if c.UsedOrderID != nil {
			q = q.Where("used_order_id = ?", *c.UsedOrderID)
		}
		if err := q.First(&existing).Error; err == nil {
			ex := existing
			result[key] = &ex
			return &ex, nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		c.ExtJSON = `{"seed":"demo"}`
		if err := tx.Create(&c).Error; err != nil {
			return nil, err
		}
		result[key] = &c
		return &c, nil
	}

	mustUser := func(k string) (*model.User, error) {
		u := users[k]
		if u == nil {
			return nil, fmt.Errorf("seed coupon missing user %s", k)
		}
		return u, nil
	}
	mustTpl := func(k string) (*model.CouponTemplate, error) {
		t := templates[k]
		if t == nil {
			return nil, fmt.Errorf("seed coupon missing coupon template %s", k)
		}
		return t, nil
	}

	farFuture := time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC)
	farPast := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	claimedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	lockedAt := time.Date(2025, 1, 2, 10, 0, 0, 0, time.UTC)
	usedAt := time.Date(2025, 1, 3, 10, 0, 0, 0, time.UTC)

	customerG, err := mustUser("customerG")
	if err != nil {
		return nil, err
	}
	newUserTpl, err := mustTpl("new_user_deduct")
	if err != nil {
		return nil, err
	}
	_, err = ensureCoupon("customerG_new_user_available", couponFromTemplate(customerG.ID, newUserTpl, model.CouponStateAvailable, &claimedAt, farFuture))
	if err != nil {
		return nil, err
	}

	customerA, err := mustUser("customerA")
	if err != nil {
		return nil, err
	}
	vipTpl, err := mustTpl("vip_monthly_discount")
	if err != nil {
		return nil, err
	}
	_, err = ensureCoupon("customerA_vip_available", couponFromTemplate(customerA.ID, vipTpl, model.CouponStateAvailable, &claimedAt, farFuture))
	if err != nil {
		return nil, err
	}

	customerB, err := mustUser("customerB")
	if err != nil {
		return nil, err
	}
	orderConfirmed := orders["orderConfirmed2"]
	if orderConfirmed != nil {
		locked := couponFromTemplate(customerB.ID, vipTpl, model.CouponStateLocked, &claimedAt, farFuture)
		locked.LockedByOrderID = &orderConfirmed.ID
		locked.LockedAt = &lockedAt
		_, err = ensureCoupon("customerB_vip_locked", locked)
		if err != nil {
			return nil, err
		}
	}

	customerC, err := mustUser("customerC")
	if err != nil {
		return nil, err
	}
	rechargeTpl, err := mustTpl("recharge_deduct")
	if err != nil {
		return nil, err
	}
	orderUsed := orders["orderPending1"]
	if orderUsed != nil {
		used := couponFromTemplate(customerC.ID, rechargeTpl, model.CouponStateUsed, &claimedAt, farFuture)
		used.UsedOrderID = &orderUsed.ID
		used.UsedAt = &usedAt
		used.DiscountCents = 500
		_, err = ensureCoupon("customerC_recharge_used", used)
		if err != nil {
			return nil, err
		}
	}

	customerD, err := mustUser("customerD")
	if err != nil {
		return nil, err
	}
	expired := couponFromTemplate(customerD.ID, newUserTpl, model.CouponStateExpired, &claimedAt, farPast)
	_, err = ensureCoupon("customerD_new_user_expired", expired)
	if err != nil {
		return nil, err
	}

	log.Printf("user coupons ensured: %d\n", len(result))
	return result, nil
}

func seedRechargeDataInternal(tx *gorm.DB, templates map[string]*model.CouponTemplate, users map[string]*model.User) error {
	var optionCount int64
	if err := tx.Model(&model.RechargeOption{}).Count(&optionCount).Error; err != nil {
		return err
	}

	rechargeTpl := templates["recharge_deduct"]
	var rechargeTplID *uint64
	if rechargeTpl != nil {
		rechargeTplID = &rechargeTpl.ID
	}

	if optionCount == 0 {
		original := int64(10000)
		discount := 90
		options := []model.RechargeOption{
			{
				Name:             "50元（演示）",
				AmountCents:      5000,
				BonusCents:       200,
				TotalCents:       5200,
				Description:      "演示数据：小额体验档位",
				Tag:              "新手",
				IconURL:          "https://example.com/recharge/50.png",
				SortOrder:        10,
				IsActive:         true,
				IsRecommended:    false,
				CouponTemplateID: rechargeTplID,
				CouponCount:      1,
			},
			{
				Name:             "100元（演示）",
				AmountCents:      10000,
				BonusCents:       800,
				TotalCents:       10800,
				OriginalCents:    &original,
				DiscountPercent:  &discount,
				Description:      "演示数据：推荐档位",
				Tag:              "推荐",
				IconURL:          "https://example.com/recharge/100.png",
				SortOrder:        20,
				IsActive:         true,
				IsRecommended:    true,
				CouponTemplateID: rechargeTplID,
				CouponCount:      2,
			},
			{
				Name:             "300元（演示）",
				AmountCents:      30000,
				BonusCents:       4500,
				TotalCents:       34500,
				Description:      "演示数据：大额优惠档位",
				Tag:              "热门",
				IconURL:          "https://example.com/recharge/300.png",
				SortOrder:        30,
				IsActive:         true,
				IsRecommended:    false,
				CouponTemplateID: rechargeTplID,
				CouponCount:      3,
			},
		}
		for i := range options {
			options[i].ExtJSON = `{"seed":"demo"}`
			if err := tx.Create(&options[i]).Error; err != nil {
				return err
			}
		}
		log.Println("recharge options seed data created")
	}

	u := users["customerB"]
	if u == nil {
		return nil
	}

	// 幂等：用 order_no 作为键
	var existing model.RechargeRecord
	if err := tx.Where("order_no = ?", "RCG-DEMO-0001").First(&existing).Error; err == nil {
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	var option model.RechargeOption
	if err := tx.Where("name = ?", "100元（演示）").First(&option).Error; err != nil {
		return err
	}

	now := time.Now()
	expireAt := now.Add(30 * time.Minute)
	paidAt := now.Add(-10 * time.Minute)

	record := model.RechargeRecord{
		UserID:           u.ID,
		OptionID:         &option.ID,
		AmountCents:      option.AmountCents,
		BonusCents:       option.BonusCents,
		TotalCents:       option.TotalCents,
		Status:           model.RechargeStatusPaid,
		OrderNo:          "RCG-DEMO-0001",
		MerchantOrderNo:  "MCH-RCG-DEMO-0001",
		ProviderTradeNo:  "PROVIDER-RCG-DEMO-0001",
		MerchantID:       "mch-demo-001",
		CollectionEntity: "demo_entity_001",
		PaymentChannel:   "wechat",
		PaymentMethod:    "wechat_h5",
		PaidAt:           &paidAt,
		ExpireAt:         &expireAt,
		CouponIssued:     false,
		CouponIDs:        "[]",
		ClientIP:         "127.0.0.1",
		UserAgent:        "seed/1.0",
		DeviceInfo:       `{"os":"windows","channel":"demo"}`,
		Remark:           "演示数据：已支付充值记录",
	}
	record.ExtJSON = `{"seed":"demo"}`
	if err := tx.Create(&record).Error; err != nil {
		return err
	}

	// 如果配置了充值赠送券模板，则创建并回写 CouponIDs
	if rechargeTpl != nil {
		issuedCoupons := make([]uint64, 0, option.CouponCount)
		for i := 0; i < option.CouponCount; i++ {
			coupon := couponFromTemplate(u.ID, rechargeTpl, model.CouponStateAvailable, &now, time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC))
			coupon.Source = model.CouponSourceRecharge
			if err := tx.Create(&coupon).Error; err != nil {
				return err
			}
			issuedCoupons = append(issuedCoupons, coupon.ID)
		}
		rawIDs, _ := json.Marshal(issuedCoupons)
		if err := tx.Model(&record).Updates(map[string]interface{}{
			"coupon_issued": true,
			"coupon_ids":    string(rawIDs),
		}).Error; err != nil {
			return err
		}
	}

	log.Println("recharge records seed data created")
	return nil
}

func seedActivityDataInternal(tx *gorm.DB, templates map[string]*model.CouponTemplate, users map[string]*model.User) (map[string]*model.Activity, error) {
	result := make(map[string]*model.Activity)

	activityTpl := templates["activity_deduct"]
	if activityTpl == nil {
		return result, nil
	}

	ensureActivity := func(key string, a model.Activity) (*model.Activity, error) {
		var existing model.Activity
		if err := tx.Where("name = ?", a.Name).First(&existing).Error; err == nil {
			ex := existing
			result[key] = &ex
			return &ex, nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		a.ExtJSON = `{"seed":"demo"}`
		if err := tx.Create(&a).Error; err != nil {
			return nil, err
		}
		result[key] = &a
		return &a, nil
	}

	now := time.Now()
	start := now.Add(-24 * time.Hour)
	end := now.Add(7 * 24 * time.Hour)
	preheat := now.Add(-48 * time.Hour)

	active, err := ensureActivity("activity_active", model.Activity{
		Name:         "新手福利活动（演示）",
		Description:  "演示数据：参与活动可领取优惠券",
		Type:         model.ActivityTypeCoupon,
		Status:       model.ActivityStatusActive,
		CoverURL:     "https://example.com/activity/cover.png",
		BannerURL:    "https://example.com/activity/banner.png",
		PreheatAt:    &preheat,
		StartAt:      start,
		EndAt:        end,
		PerUserLimit: 1,

		AllowVipStack: false,
		Rules:         "演示规则：每人限领1次；优惠券有效期7天。",
		SortOrder:     10,
		IsVisible:     true,
	})
	if err != nil {
		return nil, err
	}

	endedStart := now.Add(-30 * 24 * time.Hour)
	endedEnd := now.Add(-15 * 24 * time.Hour)
	_, err = ensureActivity("activity_ended", model.Activity{
		Name:         "周末加码活动（演示-已结束）",
		Description:  "演示数据：已结束活动，用于筛选与状态展示",
		Type:         model.ActivityTypeCoupon,
		Status:       model.ActivityStatusEnded,
		StartAt:      endedStart,
		EndAt:        endedEnd,
		PerUserLimit: 1,
		Rules:        "演示规则：活动已结束。",
		SortOrder:    20,
		IsVisible:    true,
	})
	if err != nil {
		return nil, err
	}

	// 确保奖励配置存在
	var reward model.ActivityReward
	if err := tx.Where("activity_id = ? AND coupon_template_id = ?", active.ID, activityTpl.ID).First(&reward).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		reward = model.ActivityReward{
			ActivityID:       active.ID,
			CouponTemplateID: activityTpl.ID,
			CouponCount:      1,
			Probability:      100,
			TotalStock:       0,
			RemainingStock:   0,
			SortOrder:        10,
		}
		reward.ExtJSON = `{"seed":"demo"}`
		if err := tx.Create(&reward).Error; err != nil {
			return nil, err
		}
	}

	// 创建参与记录（幂等：按 activity+user）
	customerA := users["customerA"]
	if customerA != nil {
		var existing model.ActivityParticipation
		if err := tx.Where("activity_id = ? AND user_id = ?", active.ID, customerA.ID).First(&existing).Error; err == nil {
			return result, nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		coupon := couponFromTemplate(customerA.ID, activityTpl, model.CouponStateAvailable, &now, time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC))
		coupon.Source = model.CouponSourceActivity
		if err := tx.Create(&coupon).Error; err != nil {
			return nil, err
		}
		rawIDs, _ := json.Marshal([]uint64{coupon.ID})

		p := model.ActivityParticipation{
			ActivityID: active.ID,
			UserID:     customerA.ID,
			RewardID:   reward.ID,
			CouponIDs:  string(rawIDs),
			ClaimedAt:  now.Add(-2 * time.Hour),
			ClientIP:   "127.0.0.1",
		}
		p.ExtJSON = `{"seed":"demo"}`
		if err := tx.Create(&p).Error; err != nil {
			return nil, err
		}
	}

	log.Printf("activities ensured: %d\n", len(result))
	return result, nil
}

func seedTeamDataInternal(tx *gorm.DB, players map[string]*model.Player, orders map[string]*model.Order) error {
	playerA := players["playerA"]
	playerB := players["playerB"]
	playerC := players["playerC"]
	if playerA == nil || playerB == nil || playerC == nil {
		return nil
	}

	ensureTeam := func(name string, t model.Team) (*model.Team, error) {
		var existing model.Team
		if err := tx.Where("name = ?", name).First(&existing).Error; err == nil {
			// Update CurrentOrderID if it has changed (order IDs may change after cleanup)
			if (existing.CurrentOrderID == nil && t.CurrentOrderID != nil) ||
				(existing.CurrentOrderID != nil && t.CurrentOrderID == nil) ||
				(existing.CurrentOrderID != nil && t.CurrentOrderID != nil && *existing.CurrentOrderID != *t.CurrentOrderID) {
				if err := tx.Model(&existing).Update("current_order_id", t.CurrentOrderID).Error; err != nil {
					return nil, err
				}
				existing.CurrentOrderID = t.CurrentOrderID
			}
			return &existing, nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		t.ExtJSON = `{"seed":"demo"}`
		if err := tx.Create(&t).Error; err != nil {
			return nil, err
		}
		return &t, nil
	}

	activeTeam, err := ensureTeam("峡谷车队（演示）", model.Team{
		Name:            "峡谷车队（演示）",
		Description:     "演示数据：用于团队管理模块测试",
		AvatarURL:       "https://example.com/team/avatar1.png",
		LeaderID:        playerA.ID,
		Status:          model.TeamStatusActive,
		MaxMembers:      5,
		MemberCount:     3,
		IncomeShareType: "equal",
		TotalOrderCount: 12,
		TotalIncomeCents: 258000,
	})
	if err != nil {
		return err
	}

	var currentOrderID *uint64
	if o := orders["orderInProgress1"]; o != nil {
		currentOrderID = &o.ID
	}
	busyTeam, err := ensureTeam("冠军冲分队（演示）", model.Team{
		Name:            "冠军冲分队（演示）",
		Description:     "演示数据：忙碌状态团队（CurrentOrderID 非空）",
		AvatarURL:       "https://example.com/team/avatar2.png",
		LeaderID:        playerB.ID,
		Status:          model.TeamStatusBusy,
		MaxMembers:      5,
		MemberCount:     2,
		IncomeShareType: "equal",
		CurrentOrderID:  currentOrderID,
		TotalOrderCount: 5,
		TotalIncomeCents: 98000,
	})
	if err != nil {
		return err
	}

	now := time.Now()

	ensureMember := func(teamID, playerID uint64, role model.TeamMemberRole, sortOrder int) error {
		var existing model.TeamMember
		if err := tx.Where("team_id = ? AND player_id = ?", teamID, playerID).First(&existing).Error; err == nil {
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		member := model.TeamMember{
			TeamID:    teamID,
			PlayerID:  playerID,
			Role:      role,
			Status:    model.TeamMemberStatusActive,
			JoinedAt:  now.Add(-7 * 24 * time.Hour),
			SortOrder: sortOrder,
		}
		member.ExtJSON = `{"seed":"demo"}`
		return tx.Create(&member).Error
	}

	_ = ensureMember(activeTeam.ID, playerA.ID, model.TeamMemberRoleLeader, 0)
	_ = ensureMember(activeTeam.ID, playerC.ID, model.TeamMemberRoleMember, 10)
	_ = ensureMember(activeTeam.ID, playerB.ID, model.TeamMemberRoleMember, 20)

	_ = ensureMember(busyTeam.ID, playerB.ID, model.TeamMemberRoleLeader, 0)
	_ = ensureMember(busyTeam.ID, playerA.ID, model.TeamMemberRoleMember, 10)

	// 邀请记录（演示：pending）
	inviteExpire := now.Add(3 * 24 * time.Hour)
	var existingInvite model.TeamInvite
	if err := tx.Where("team_id = ? AND player_id = ? AND status = ?", activeTeam.ID, playerB.ID, model.TeamInviteStatusPending).First(&existingInvite).Error; err == nil {
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	invite := model.TeamInvite{
		TeamID:    activeTeam.ID,
		PlayerID:  playerB.ID,
		InviterID: playerA.ID,
		Status:    model.TeamInviteStatusPending,
		ExpireAt:  inviteExpire,
		Message:   "演示数据：邀请加入团队",
	}
	invite.ExtJSON = `{"seed":"demo"}`
	_ = tx.Create(&invite).Error

	log.Println("team seed data ensured")
	return nil
}

func seedReferralDataInternal(tx *gorm.DB, templates map[string]*model.CouponTemplate, users map[string]*model.User) error {
	// 全局配置（确保存在）
	configs := []model.ReferralConfig{
		{ConfigKey: model.ReferralConfigEnabled, ConfigValue: "true", Description: "演示：启用推荐系统"},
		{ConfigKey: model.ReferralConfigExpireDays, ConfigValue: "30", Description: "演示：邀请码过期天数"},
		{ConfigKey: model.ReferralConfigMaxLevel, ConfigValue: "2", Description: "演示：最大层级"},
		{ConfigKey: model.ReferralConfigUserRewardType, ConfigValue: string(model.RewardTypeCash), Description: "演示：用户邀请奖励类型"},
		{ConfigKey: model.ReferralConfigUserRewardAmount, ConfigValue: "500", Description: "演示：用户邀请奖励金额（分）"},
	}
	for _, cfg := range configs {
		var existing model.ReferralConfig
		if err := tx.Where("config_key = ?", cfg.ConfigKey).First(&existing).Error; err == nil {
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		cfg.ExtJSON = `{"seed":"demo"}`
		if err := tx.Create(&cfg).Error; err != nil {
			return err
		}
	}

	referrer := users["customerA"]
	referee := users["customerG"]
	if referrer == nil || referee == nil {
		return nil
	}

	// 邀请码（uniqueIndex by code）
	code := model.ReferralCode{
		Code:     "REF-CUST-A-DEMO",
		UserID:   referrer.ID,
		Type:     model.ReferralTypeUserToUser,
		IsActive: true,
		MaxUse:   0,
	}
	code.ExtJSON = `{"seed":"demo"}`
	var existingCode model.ReferralCode
	if err := tx.Where("code = ?", code.Code).First(&existingCode).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err := tx.Create(&code).Error; err != nil {
			return err
		}
		existingCode = code
	}

	// 推荐记录（按 referrer+referee+type 做幂等）
	var existing model.Referral
	if err := tx.Where("referrer_id = ? AND referee_id = ? AND type = ? AND level = ?", referrer.ID, referee.ID, model.ReferralTypeUserToUser, 1).First(&existing).Error; err == nil {
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	now := time.Now()
	completedAt := now.Add(-24 * time.Hour)
	r := model.Referral{
		ReferrerID:        referrer.ID,
		RefereeID:         referee.ID,
		CodeID:            &existingCode.ID,
		Type:              model.ReferralTypeUserToUser,
		Level:             1,
		Status:            model.ReferralStatusRewarded,
		CompletedAt:       &completedAt,
		RewardType:        model.RewardTypeCash,
		RewardAmountCents: 500,
		RewardedAt:        &now,
		RewardNote:        "演示数据：邀请奖励已发放",
		RefereeCondition:  "first_order",
	}
	r.ExtJSON = `{"seed":"demo"}`
	if err := tx.Create(&r).Error; err != nil {
		return err
	}

	// 奖励明细（演示：现金 + 可选优惠券）
	reward := model.ReferralReward{
		ReferralID:  r.ID,
		UserID:      referrer.ID,
		Type:        model.RewardTypeCash,
		AmountCents: 500,
		Status:      model.ReferralRewardStatusIssued,
		IssuedAt:    &now,
	}
	reward.ExtJSON = `{"seed":"demo"}`
	_ = tx.Create(&reward).Error

	manualTpl := templates["lol_only_discount"]
	if manualTpl != nil {
		coupon := couponFromTemplate(referrer.ID, manualTpl, model.CouponStateAvailable, &now, time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC))
		if err := tx.Create(&coupon).Error; err == nil {
			rewardCoupon := model.ReferralReward{
				ReferralID:  r.ID,
				UserID:      referrer.ID,
				Type:        model.RewardTypeCoupon,
				AmountCents: 0,
				CouponID:    &coupon.ID,
				Status:      model.ReferralRewardStatusIssued,
				IssuedAt:    &now,
			}
			rewardCoupon.ExtJSON = `{"seed":"demo"}`
			_ = tx.Create(&rewardCoupon).Error
		}
	}

	log.Println("referral seed data ensured")
	return nil
}

func seedUserBlockDataInternal(tx *gorm.DB, users map[string]*model.User) error {
	var cnt int64
	if err := tx.Model(&model.UserBlock{}).Count(&cnt).Error; err != nil {
		return err
	}
	if cnt > 0 {
		return nil
	}

	blocker := users["customerA"]
	blocked := users["proB"]
	admin := users["adminA"]
	if blocker == nil || blocked == nil {
		return nil
	}

	now := time.Now()

	active := model.UserBlock{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypePlayer,
		Reason:      "演示数据：服务体验不佳",
		Status:      model.BlockStatusActive,
		BlockedAt:   now.Add(-72 * time.Hour),
	}
	active.ExtJSON = `{"seed":"demo"}`
	if err := tx.Create(&active).Error; err != nil {
		return err
	}

	if admin != nil {
		canceledAt := now.Add(-24 * time.Hour)
		adminCanceled := model.UserBlock{
			BlockerID:   blocked.ID,
			BlockerType: model.BlockUserTypePlayer,
			BlockedID:   blocker.ID,
			BlockedType: model.BlockUserTypeUser,
			Reason:      "演示数据：误操作拉黑",
			Status:      model.BlockStatusAdminCanceled,
			BlockedAt:   now.Add(-96 * time.Hour),
			CanceledAt:  &canceledAt,
			CanceledBy:  &admin.ID,
			AdminRemark: "演示数据：管理员强制解除",
		}
		adminCanceled.ExtJSON = `{"seed":"demo"}`
		_ = tx.Create(&adminCanceled).Error
	}

	log.Println("user block seed data created")
	return nil
}

func seedGameRankAndCertificationDataInternal(tx *gorm.DB, games map[string]*model.Game, players map[string]*model.Player, users map[string]*model.User) error {
	lol := games["lol"]
	valorant := games["valorant"]
	if lol == nil || valorant == nil {
		return nil
	}

	ensureRank := func(gameID uint64, level int, name string, price int64, color string) (*model.GameRank, error) {
		var existing model.GameRank
		if err := tx.Where("game_id = ? AND level = ?", gameID, level).First(&existing).Error; err == nil {
			return &existing, nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		rank := model.GameRank{
			GameID:       gameID,
			Name:         name,
			Level:        level,
			PriceCents:   price,
			IconURL:      "https://example.com/rank/icon.png",
			Color:        color,
			Description:  "演示数据：用于段位与定价展示",
			SortOrder:    level,
			IsActive:     true,
		}
		rank.ExtJSON = `{"seed":"demo"}`
		if err := tx.Create(&rank).Error; err != nil {
			return nil, err
		}
		return &rank, nil
	}

	lolGold, err := ensureRank(lol.ID, 3, "黄金", 9900, "#FFD700")
	if err != nil {
		return err
	}
	lolDiamond, err := ensureRank(lol.ID, 5, "钻石", 15900, "#4B9CD3")
	if err != nil {
		return err
	}
	valSilver, err := ensureRank(valorant.ID, 2, "白银", 10900, "#C0C0C0")
	if err != nil {
		return err
	}

	admin := users["adminA"]
	playerA := players["playerA"]
	playerB := players["playerB"]
	playerC := players["playerC"]
	if admin == nil || playerA == nil || playerB == nil || playerC == nil {
		return nil
	}

	ensureRankRecord := func(p *model.Player, g *model.Game, rank *model.GameRank, status model.PlayerRankStatus, rejectReason string) error {
		var existing model.PlayerRankRecord
		if err := tx.Where("player_id = ? AND rank_id = ?", p.ID, rank.ID).First(&existing).Error; err == nil {
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		now := time.Now()
		rec := model.PlayerRankRecord{
			PlayerID:       p.ID,
			GameID:         g.ID,
			RankID:         rank.ID,
			Status:         status,
			ScreenshotURLs: `["https://example.com/rank/shot1.png"]`,
			VerifiedAt:     &now,
			VerifiedBy:     &admin.ID,
			RejectReason:   rejectReason,
			Remark:         "演示数据：用于陪玩师段位审核列表",
		}
		if status == model.PlayerRankStatusPending || status == model.PlayerRankStatusRejected {
			rec.VerifiedAt = nil
			rec.VerifiedBy = nil
		}
		rec.ExtJSON = `{"seed":"demo"}`
		return tx.Create(&rec).Error
	}

	_ = ensureRankRecord(playerA, lol, lolDiamond, model.PlayerRankStatusVerified, "")
	_ = ensureRankRecord(playerB, valorant, valSilver, model.PlayerRankStatusPending, "")
	_ = ensureRankRecord(playerC, lol, lolGold, model.PlayerRankStatusRejected, "截图不清晰")

	ensureCertification := func(p *model.Player, status model.CertificationStatus, rejectReason string) error {
		var existing model.PlayerCertification
		if err := tx.Where("player_id = ?", p.ID).First(&existing).Error; err == nil {
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		now := time.Now()
		cert := model.PlayerCertification{
			PlayerID:       p.ID,
			RealName:       "张三",
			IDCardNo:       "ENCRYPTED_DEMO",
			IDCardFrontURL: "https://example.com/id/front.png",
			IDCardBackURL:  "https://example.com/id/back.png",
			Status:         status,
			VerifiedAt:     &now,
			VerifiedBy:     &admin.ID,
			RejectReason:   rejectReason,
			PhotoURL:       "https://example.com/photo.png",
			VoiceURL:       "https://example.com/voice.mp3",
			ExtJSON:        `{"seed":"demo"}`,
		}
		if status == model.CertificationStatusPending || status == model.CertificationStatusRejected {
			cert.VerifiedAt = nil
			cert.VerifiedBy = nil
		}
		return tx.Create(&cert).Error
	}

	_ = ensureCertification(playerA, model.CertificationStatusVerified, "")
	_ = ensureCertification(playerB, model.CertificationStatusPending, "")
	_ = ensureCertification(playerC, model.CertificationStatusRejected, "证件信息不一致")

	log.Println("game ranks and certifications seed data ensured")
	return nil
}

func seedNotificationDataInternal(tx *gorm.DB, users map[string]*model.User, orders map[string]*model.Order, coupons map[string]*model.Coupon, activities map[string]*model.Activity, vipLevels map[string]*model.VipLevel) error {
	ensureTemplate := func(tpl model.NotificationTemplate) (*model.NotificationTemplate, error) {
		var existing model.NotificationTemplate
		if err := tx.Where("code = ?", tpl.Code).First(&existing).Error; err == nil {
			return &existing, nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		tpl.ExtJSON = `{"seed":"demo"}`
		if tpl.Channels == "" {
			tpl.Channels = `["in_app"]`
		}
		if err := tx.Create(&tpl).Error; err != nil {
			return nil, err
		}
		return &tpl, nil
	}

	orderTpl, err := ensureTemplate(model.NotificationTemplate{
		Code:     "tpl_order_status_demo",
		Name:     "订单状态通知（演示）",
		Type:     model.NotificationTypeOrderStatus,
		Title:    "订单状态更新",
		Content:  "您的订单状态已更新：{{status}}（演示模板）",
		Channels: `["in_app","email"]`,
		IsActive: true,
		IsSystem: true,
	})
	if err != nil {
		return err
	}

	couponTpl, err := ensureTemplate(model.NotificationTemplate{
		Code:     "tpl_coupon_expire_demo",
		Name:     "优惠券过期提醒（演示）",
		Type:     model.NotificationTypeCouponExpire,
		Title:    "优惠券即将过期",
		Content:  "您有一张优惠券即将过期，请尽快使用（演示模板）",
		Channels: `["in_app"]`,
		IsActive: true,
		IsSystem: true,
	})
	if err != nil {
		return err
	}

	activityTpl, err := ensureTemplate(model.NotificationTemplate{
		Code:     "tpl_activity_start_demo",
		Name:     "活动开始提醒（演示）",
		Type:     model.NotificationTypeActivityStart,
		Title:    "活动开始啦",
		Content:  "活动已开始，快来参与领取奖励（演示模板）",
		Channels: `["in_app","push"]`,
		IsActive: true,
		IsSystem: true,
	})
	if err != nil {
		return err
	}

	_, err = ensureTemplate(model.NotificationTemplate{
		Code:     "tpl_system_announcement_demo",
		Name:     "系统公告（演示）",
		Type:     model.NotificationTypeSystem,
		Title:    "系统公告",
		Content:  "欢迎体验 GameLink（演示公告）",
		Channels: `["in_app"]`,
		IsActive: true,
		IsSystem: true,
	})
	if err != nil {
		return err
	}

	// 配置：按 config_key 幂等
	configs := []model.NotificationConfig{
		{ConfigKey: model.NotificationConfigVipExpireDays, ConfigValue: "[7,3,1]", Description: "演示：VIP到期提醒天数"},
		{ConfigKey: model.NotificationConfigCouponExpireDays, ConfigValue: "[7,3,1]", Description: "演示：优惠券到期提醒天数"},
		{ConfigKey: model.NotificationConfigPushProvider, ConfigValue: "demo", Description: "演示：推送服务商"},
	}
	for _, cfg := range configs {
		var existing model.NotificationConfig
		if err := tx.Where("config_key = ?", cfg.ConfigKey).First(&existing).Error; err == nil {
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		cfg.ExtJSON = `{"seed":"demo"}`
		_ = tx.Create(&cfg).Error
	}

	// 用户设置（按 user_id 幂等）
	if customerA := users["customerA"]; customerA != nil {
		var existing model.UserNotificationSetting
		if err := tx.Where("user_id = ?", customerA.ID).First(&existing).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			setting := model.UserNotificationSetting{
				UserID:              customerA.ID,
				OrderStatusEnabled:  true,
				VipExpireEnabled:    true,
				CouponExpireEnabled: true,
				ActivityEnabled:     true,
				SystemEnabled:       true,
				InAppEnabled:        true,
				PushEnabled:         true,
				EmailEnabled:        true,
				DoNotDisturbEnabled: true,
				DoNotDisturbStart:   "23:00",
				DoNotDisturbEnd:     "08:00",
			}
			setting.ExtJSON = `{"seed":"demo"}`
			_ = tx.Create(&setting).Error
		}
	}

	// 通知记录：按 (user_id, type, related_type, related_id, title) 做幂等
	ensureUserNotif := func(n model.UserNotification) error {
		var existing model.UserNotification
		q := tx.Where("user_id = ? AND type = ? AND channel = ? AND title = ?", n.UserID, n.Type, n.Channel, n.Title).
			Where("related_type = ?", n.RelatedType)
		if n.RelatedID != nil {
			q = q.Where("related_id = ?", *n.RelatedID)
		} else {
			q = q.Where("related_id IS NULL")
		}
		if err := q.First(&existing).Error; err == nil {
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		n.ExtJSON = `{"seed":"demo"}`
		return tx.Create(&n).Error
	}

	now := time.Now()

	if customerA := users["customerA"]; customerA != nil {
		if o := orders["orderCompleted1"]; o != nil {
			_ = ensureUserNotif(model.UserNotification{
				UserID:      customerA.ID,
				TemplateID:  &orderTpl.ID,
				Type:        model.NotificationTypeOrderStatus,
				Channel:     model.NotificationChannelInApp,
				Title:       "订单已完成（演示）",
				Content:     fmt.Sprintf("订单 %s 已完成，感谢使用（演示数据）", o.OrderNo),
				Status:      model.NotificationStatusRead,
				ReadAt:      &now,
				SentAt:      &now,
				RelatedType: "order",
				RelatedID:   &o.ID,
			})
		}

		if c := coupons["customerD_new_user_expired"]; c != nil {
			sentAt := now.Add(-2 * time.Hour)
			_ = ensureUserNotif(model.UserNotification{
				UserID:      customerA.ID,
				TemplateID:  &couponTpl.ID,
				Type:        model.NotificationTypeCouponExpire,
				Channel:     model.NotificationChannelInApp,
				Title:       "优惠券过期提醒（演示）",
				Content:     "您有一张优惠券即将过期，请尽快使用。",
				Status:      model.NotificationStatusSent,
				SentAt:      &sentAt,
				RelatedType: "coupon",
				RelatedID:   &c.ID,
			})
		}

		if a := activities["activity_active"]; a != nil {
			_ = ensureUserNotif(model.UserNotification{
				UserID:      customerA.ID,
				TemplateID:  &activityTpl.ID,
				Type:        model.NotificationTypeActivityStart,
				Channel:     model.NotificationChannelInApp,
				Title:       "活动开始提醒（演示）",
				Content:     fmt.Sprintf("活动「%s」正在进行中，快来参与！", a.Name),
				Status:      model.NotificationStatusSent,
				SentAt:      &now,
				RelatedType: "activity",
				RelatedID:   &a.ID,
			})
		}
	}

	// VIP 到期提醒（演示，不依赖用户绑定关系）
	if customerB := users["customerB"]; customerB != nil {
		if vip := vipLevels["vip1"]; vip != nil {
			_ = ensureUserNotif(model.UserNotification{
				UserID:      customerB.ID,
				Type:        model.NotificationTypeVipExpire,
				Channel:     model.NotificationChannelInApp,
				Title:       "VIP到期提醒（演示）",
				Content:     fmt.Sprintf("您的 %s 即将到期（演示数据）", vip.Title),
				Status:      model.NotificationStatusPending,
				RelatedType: "vip",
				RelatedID:   &vip.ID,
			})
		}
	}

	// 为管理员创建系统通知（演示管理员 adminA）
	if adminA := users["adminA"]; adminA != nil {
		sentAt := now.Add(-1 * time.Hour)
		_ = ensureUserNotif(model.UserNotification{
			UserID:      adminA.ID,
			Type:        model.NotificationTypeSystem,
			Channel:     model.NotificationChannelInApp,
			Title:       "系统公告",
			Content:     "欢迎使用 GameLink 管理后台！",
			Status:      model.NotificationStatusSent,
			SentAt:      &sentAt,
			RelatedType: "system",
		})

		sentAt2 := now.Add(-30 * time.Minute)
		_ = ensureUserNotif(model.UserNotification{
			UserID:      adminA.ID,
			Type:        model.NotificationTypeSystem,
			Channel:     model.NotificationChannelInApp,
			Title:       "新用户注册提醒",
			Content:     "今日有 5 位新用户注册，请及时审核。",
			Status:      model.NotificationStatusSent,
			SentAt:      &sentAt2,
			RelatedType: "system",
		})

		_ = ensureUserNotif(model.UserNotification{
			UserID:      adminA.ID,
			Type:        model.NotificationTypeSystem,
			Channel:     model.NotificationChannelInApp,
			Title:       "陪玩师入驻申请",
			Content:     "有 3 位陪玩师提交了入驻申请，等待审核。",
			Status:      model.NotificationStatusPending,
			RelatedType: "system",
		})
	}

	// 为超级管理员创建系统通知（从数据库查询超管账户）
	var superAdmin model.User
	if err := tx.Where("role = ? AND email LIKE ?", model.RoleAdmin, "%admin%").
		Order("id ASC").First(&superAdmin).Error; err == nil {
		// 检查是否已有通知
		var existingCount int64
		tx.Model(&model.UserNotification{}).Where("user_id = ?", superAdmin.ID).Count(&existingCount)
		if existingCount == 0 {
			sentAt := now.Add(-2 * time.Hour)
			_ = ensureUserNotif(model.UserNotification{
				UserID:      superAdmin.ID,
				Type:        model.NotificationTypeSystem,
				Channel:     model.NotificationChannelInApp,
				Title:       "欢迎使用 GameLink",
				Content:     "您已成功登录管理后台，开始管理您的平台吧！",
				Status:      model.NotificationStatusSent,
				SentAt:      &sentAt,
				RelatedType: "system",
			})

			sentAt2 := now.Add(-1 * time.Hour)
			_ = ensureUserNotif(model.UserNotification{
				UserID:      superAdmin.ID,
				Type:        model.NotificationTypeSystem,
				Channel:     model.NotificationChannelInApp,
				Title:       "系统初始化完成",
				Content:     "演示数据已加载完成，您可以开始体验各项功能。",
				Status:      model.NotificationStatusSent,
				SentAt:      &sentAt2,
				RelatedType: "system",
			})

			_ = ensureUserNotif(model.UserNotification{
				UserID:      superAdmin.ID,
				Type:        model.NotificationTypeSystem,
				Channel:     model.NotificationChannelInApp,
				Title:       "待处理事项",
				Content:     "您有 2 个陪玩师入驻申请待审核，1 个提现申请待处理。",
				Status:      model.NotificationStatusPending,
				RelatedType: "system",
			})
		}
	}

	// ========== 为 notification_events 表创建数据（API 实际查询的表）==========
	// NotificationEvent 是 API /notifications 端点实际使用的模型
	ensureNotificationEvent := func(n model.NotificationEvent) error {
		var existing model.NotificationEvent
		q := tx.Where("user_id = ? AND title = ? AND channel = ?", n.UserID, n.Title, n.Channel)
		if n.ReferenceID != nil {
			q = q.Where("reference_id = ?", *n.ReferenceID)
		} else {
			q = q.Where("reference_id IS NULL")
		}
		if err := q.First(&existing).Error; err == nil {
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		n.ExtJSON = `{"seed":"demo"}`
		return tx.Create(&n).Error
	}

	// 为超级管理员创建 NotificationEvent（API 实际查询的表）
	if superAdmin.ID > 0 {
		_ = ensureNotificationEvent(model.NotificationEvent{
			UserID:        superAdmin.ID,
			Title:         "欢迎使用 GameLink",
			Message:       "您已成功登录管理后台，开始管理您的平台吧！",
			Channel:       "web",
			Priority:      model.NotificationPriorityNormal,
			ReferenceType: "system",
		})

		_ = ensureNotificationEvent(model.NotificationEvent{
			UserID:        superAdmin.ID,
			Title:         "系统初始化完成",
			Message:       "演示数据已加载完成，您可以开始体验各项功能。",
			Channel:       "web",
			Priority:      model.NotificationPriorityNormal,
			ReferenceType: "system",
		})

		_ = ensureNotificationEvent(model.NotificationEvent{
			UserID:        superAdmin.ID,
			Title:         "待处理事项",
			Message:       "您有 2 个陪玩师入驻申请待审核，1 个提现申请待处理。",
			Channel:       "web",
			Priority:      model.NotificationPriorityHigh,
			ReferenceType: "system",
		})
	}

	// 为演示管理员 adminA 创建 NotificationEvent
	if adminA := users["adminA"]; adminA != nil {
		_ = ensureNotificationEvent(model.NotificationEvent{
			UserID:        adminA.ID,
			Title:         "系统公告",
			Message:       "欢迎使用 GameLink 管理后台！",
			Channel:       "web",
			Priority:      model.NotificationPriorityNormal,
			ReferenceType: "system",
		})

		_ = ensureNotificationEvent(model.NotificationEvent{
			UserID:        adminA.ID,
			Title:         "新用户注册提醒",
			Message:       "今日有 5 位新用户注册，请及时审核。",
			Channel:       "web",
			Priority:      model.NotificationPriorityNormal,
			ReferenceType: "user",
		})

		_ = ensureNotificationEvent(model.NotificationEvent{
			UserID:        adminA.ID,
			Title:         "陪玩师入驻申请",
			Message:       "有 3 位陪玩师提交了入驻申请，等待审核。",
			Channel:       "web",
			Priority:      model.NotificationPriorityHigh,
			ReferenceType: "player",
		})
	}

	// 为普通用户创建 NotificationEvent（测试用户端通知）
	if customerA := users["customerA"]; customerA != nil {
		_ = ensureNotificationEvent(model.NotificationEvent{
			UserID:        customerA.ID,
			Title:         "订单已完成",
			Message:       "您的订单已完成，感谢使用 GameLink！",
			Channel:       "web",
			Priority:      model.NotificationPriorityNormal,
			ReferenceType: "order",
		})

		_ = ensureNotificationEvent(model.NotificationEvent{
			UserID:        customerA.ID,
			Title:         "优惠券即将过期",
			Message:       "您有一张优惠券将于 3 天后过期，请尽快使用。",
			Channel:       "web",
			Priority:      model.NotificationPriorityNormal,
			ReferenceType: "coupon",
		})
	}

	// 定时任务（演示）
	var schedule model.NotificationSchedule
	if err := tx.Where("name = ?", "每日优惠券到期提醒（演示）").First(&schedule).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		s := model.NotificationSchedule{
			Name:       "每日优惠券到期提醒（演示）",
			Type:       model.NotificationTypeCouponExpire,
			TemplateID: couponTpl.ID,
			ScheduleAt: now.Add(24 * time.Hour),
			Status:     model.NotificationScheduleStatusPending,
			TargetType: "all",
			TargetIDs:  "[]",
		}
		s.ExtJSON = `{"seed":"demo"}`
		_ = tx.Create(&s).Error
	}

	log.Println("notification seed data ensured")
	return nil
}
