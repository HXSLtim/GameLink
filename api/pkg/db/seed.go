package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"gamelink/internal/model"
)

type seedUserInput struct {
	Key      string
	Email    string
	Phone    string
	Name     string
	Role     model.Role
	Password string
}

type seedPlayerSpec struct {
	Key                string
	UserKey            string
	Nickname           string
	Bio                string
	RatingAverage      float32
	RatingCount        uint32
	HourlyRateCents    int64
	MainGameKey        string
	VerificationStatus model.VerificationStatus
}

type seedOrderSpec struct {
	Key          string
	UserKey      string
	PlayerKey    string
	GameKey      string
	ItemCode     string
	Title        string
	Description  string
	Status       model.OrderStatus
	PriceCents   int64
	Currency     model.Currency
	StartOffset  time.Duration
	Duration     time.Duration
	CancelReason string
}

type seedPaymentSpec struct {
	OrderKey        string
	UserKey         string
	Method          model.PaymentMethod
	AmountCents     int64
	Currency        model.Currency
	Status          model.PaymentStatus
	ProviderTradeNo string
	ProviderRaw     string
	PaidAtOffset    *time.Duration
	RefundedOffset  *time.Duration
}

type seedReviewSpec struct {
	OrderKey        string
	UserKey         string
	PlayerKey       string
	Score           model.Rating
	Content         string
	Status          model.ReviewStatus
	Images          []string
	IsReported      bool
	RejectionReason string
}

func applySeeds(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Clean up stale data to ensure consistency when seed data changes
		// Use raw SQL with IF EXISTS to handle tables that may not exist
		log.Println("Cleaning up stale seed data...")
		cleanupSQL := []string{
			// Child tables of orders (must be deleted first to avoid FK violations)
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'order_players') THEN DELETE FROM order_players; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'order_items') THEN DELETE FROM order_items; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'order_teams') THEN DELETE FROM order_teams; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'order_disputes') THEN DELETE FROM order_disputes; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'order_service_assignments') THEN DELETE FROM order_service_assignments; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'order_timeout_logs') THEN DELETE FROM order_timeout_logs; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'reviews') THEN DELETE FROM reviews; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'review_appeals') THEN DELETE FROM review_appeals; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'review_replies') THEN DELETE FROM review_replies; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'refund_records') THEN DELETE FROM refund_records; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'payments') THEN DELETE FROM payments; END IF; END $$;",
			// Parent tables (deleted after children)
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'orders') THEN DELETE FROM orders; END IF; END $$;",
			// Chat tables (order_service_assignments deleted before chat_groups)
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'chat_reports') THEN DELETE FROM chat_reports; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'chat_snapshots') THEN DELETE FROM chat_snapshots; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'chat_messages') THEN DELETE FROM chat_messages; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'chat_group_members') THEN DELETE FROM chat_group_members; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'chat_groups') THEN DELETE FROM chat_groups; END IF; END $$;",
			// Team tables (teams may reference orders via current_order_id)
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'team_invites') THEN DELETE FROM team_invites; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'team_members') THEN DELETE FROM team_members; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'teams') THEN DELETE FROM teams; END IF; END $$;",
			// Recharge and referral related tables (child tables first)
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'referral_rewards') THEN DELETE FROM referral_rewards; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'referrals') THEN DELETE FROM referrals; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'referral_codes') THEN DELETE FROM referral_codes; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'recharge_records') THEN DELETE FROM recharge_records; END IF; END $$;",
			// Wallet and settlement related tables
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'player_company_assignments') THEN DELETE FROM player_company_assignments; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'monthly_settlements') THEN DELETE FROM monthly_settlements; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'ranking_rewards') THEN DELETE FROM ranking_rewards; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'settlement_companies') THEN DELETE FROM settlement_companies; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'wallets') THEN DELETE FROM wallets; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'withdraws') THEN DELETE FROM withdraws; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'user_notifications') THEN DELETE FROM user_notifications; END IF; END $$;",
			// Marketing-related tables (child tables first)
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'recharge_options') THEN DELETE FROM recharge_options; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'activity_participations') THEN DELETE FROM activity_participations; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'activity_rewards') THEN DELETE FROM activity_rewards; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'activity_daily_stats') THEN DELETE FROM activity_daily_stats; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'user_coupons') THEN DELETE FROM user_coupons; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'coupons') THEN DELETE FROM coupons; END IF; END $$;",
			"DO $$ BEGIN IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'coupon_templates') THEN DELETE FROM coupon_templates; END IF; END $$;",
		}
		for _, sql := range cleanupSQL {
			if err := tx.Exec(sql).Error; err != nil {
				log.Printf("warning: failed to cleanup: %v", err)
			}
		}

		games, err := seedGames(tx)
		if err != nil {
			return err
		}

		now := time.Now()

		userInputs := []seedUserInput{
			{Key: "customerA", Email: "demo.user@gamelink.com", Phone: "13800138000", Name: "测试用户", Role: model.RoleUser, Password: "User@123456"},
			{Key: "proA", Email: "pro.player@gamelink.com", Phone: "13800138001", Name: "职业陪玩", Role: model.RolePlayer, Password: "Player@123456"},
			{Key: "customerB", Email: "vip.user@gamelink.com", Phone: "13800138002", Name: "高级会员", Role: model.RoleUser, Password: "Vip@123456"},
			{Key: "customerC", Email: "new.user@gamelink.com", Phone: "13800138003", Name: "体验用户", Role: model.RoleUser, Password: "User@123789"},
			{Key: "proB", Email: "streamer@gamelink.com", Phone: "13800138004", Name: "魔王主播", Role: model.RolePlayer, Password: "Player@654321"},
			{Key: "adminA", Email: "sysadmin@gamelink.com", Phone: "13800138100", Name: "系统管理员", Role: model.RoleAdmin, Password: "Admin@123456"},
			{Key: "customerD", Email: "casual.player@gamelink.com", Phone: "13800138006", Name: "休闲玩家", Role: model.RoleUser, Password: "User@123789"},
			{Key: "customerE", Email: "competitive.gamer@gamelink.com", Phone: "13800138007", Name: "竞技高手", Role: model.RoleUser, Password: "User@456789"},
			{Key: "proC", Email: "fps.master@gamelink.com", Phone: "13800138008", Name: "FPS大神", Role: model.RolePlayer, Password: "Player@987654"},
			{Key: "proD", Email: "rpg.explorer@gamelink.com", Phone: "13800138009", Name: "RPG探险家", Role: model.RolePlayer, Password: "Player@456123"},
			{Key: "customerF", Email: "weekend.gamer@gamelink.com", Phone: "13800138010", Name: "周末玩家", Role: model.RoleUser, Password: "User@789456"},
			{Key: "proE", Email: "sports.champion@gamelink.com", Phone: "13800138011", Name: "体育冠军", Role: model.RolePlayer, Password: "Player@789012"},
			{Key: "customerG", Email: "newbie.player@gamelink.com", Phone: "13800138012", Name: "新手玩家", Role: model.RoleUser, Password: "User@234567"},
			{Key: "proF", Email: "party.entertainer@gamelink.com", Phone: "13800138013", Name: "派对达人", Role: model.RolePlayer, Password: "Player@345678"},
			{Key: "customerH", Email: "business.professional@gamelink.com", Phone: "13800138014", Name: "商务人士", Role: model.RoleUser, Password: "User@567890"},
		}

		users := make(map[string]*model.User, len(userInputs))
		for _, input := range userInputs {
			user, err := seedUser(tx, input)
			if err != nil {
				return err
			}
			users[input.Key] = user
		}

		playerSpecs := []seedPlayerSpec{
			{
				Key:                "playerA",
				UserKey:            "proA",
				Nickname:           "峡谷守护者",
				Bio:                "全职陪玩，擅长打野位，帮助玩家快速上分。国服前100玩家。",
				RatingAverage:      4.9,
				RatingCount:        152,
				HourlyRateCents:    9900,
				MainGameKey:        "lol",
				VerificationStatus: model.VerificationVerified,
			},
			{
				Key:                "playerB",
				UserKey:            "proB",
				Nickname:           "王牌射手",
				Bio:                "FPS 资深选手，提供高强度陪练服务。参加过多次线下比赛。",
				RatingAverage:      4.7,
				RatingCount:        98,
				HourlyRateCents:    12900,
				MainGameKey:        "valorant",
				VerificationStatus: model.VerificationVerified,
			},
			{
				Key:                "playerC",
				UserKey:            "proC",
				Nickname:           "枪神降临",
				Bio:                "CS:GO职业选手，退役后专注陪玩教学。枪法精准，战术理解深入。",
				RatingAverage:      4.8,
				RatingCount:        127,
				HourlyRateCents:    15900,
				MainGameKey:        "csgo",
				VerificationStatus: model.VerificationVerified,
			},
			{
				Key:                "playerD",
				UserKey:            "proD",
				Nickname:           "异世界旅者",
				Bio:                "RPG游戏专家，熟悉各种MMORPG机制。带新手快速上手，老玩家攻克难关。",
				RatingAverage:      4.6,
				RatingCount:        89,
				HourlyRateCents:    11900,
				MainGameKey:        "wow",
				VerificationStatus: model.VerificationVerified,
			},
			{
				Key:                "playerE",
				UserKey:            "proE",
				Nickname:           "运动健将",
				Bio:                "FIFA/NBA2K专业玩家，体育游戏发烧友。战术教学，技巧提升。",
				RatingAverage:      4.5,
				RatingCount:        76,
				HourlyRateCents:    8900,
				MainGameKey:        "fifa",
				VerificationStatus: model.VerificationVerified,
			},
			{
				Key:                "playerF",
				UserKey:            "proF",
				Nickname:           "欢乐使者",
				Bio:                "派对游戏达人，擅长各种休闲游戏。轻松愉快的陪玩体验。",
				RatingAverage:      4.9,
				RatingCount:        203,
				HourlyRateCents:    7900,
				MainGameKey:        "fallguys",
				VerificationStatus: model.VerificationVerified,
			},
			{
				Key:                "playerG",
				UserKey:            "proA",
				Nickname:           "DOTA宗师",
				Bio:                "DOTA2老玩家，精通所有英雄。从新手教学到高端局指导全覆盖。",
				RatingAverage:      4.8,
				RatingCount:        143,
				HourlyRateCents:    13900,
				MainGameKey:        "dota2",
				VerificationStatus: model.VerificationVerified,
			},
			{
				Key:                "playerH",
				UserKey:            "proC",
				Nickname:           "大逃杀之王",
				Bio:                "PUBG/Apex资深玩家，枪法犀利，吃鸡率高。带你享受刺激的生存体验。",
				RatingAverage:      4.4,
				RatingCount:        67,
				HourlyRateCents:    10900,
				MainGameKey:        "pubg",
				VerificationStatus: model.VerificationPending,
			},
		}

		players := make(map[string]*model.Player, len(playerSpecs))
		for _, spec := range playerSpecs {
			user, ok := users[spec.UserKey]
			if !ok {
				return fmt.Errorf("seed player missing user %s", spec.UserKey)
			}
			game, ok := games[spec.MainGameKey]
			if !ok {
				return fmt.Errorf("seed player missing game %s", spec.MainGameKey)
			}
			player, err := seedPlayer(tx, seedPlayerParams{
				UserID:             user.ID,
				Nickname:           spec.Nickname,
				Bio:                spec.Bio,
				RatingAverage:      spec.RatingAverage,
				RatingCount:        spec.RatingCount,
				HourlyRateCents:    spec.HourlyRateCents,
				MainGameID:         game.ID,
				VerificationStatus: spec.VerificationStatus,
			})
			if err != nil {
				return err
			}
			players[spec.Key] = player
		}

		serviceItems, err := seedServiceItems(tx, games)
		if err != nil {
			return err
		}
		defaultServiceItem, ok := serviceItems["escort-default"]
		if !ok || defaultServiceItem == nil {
			return errors.New("seed service item escort-default missing")
		}

		hour := time.Hour

		orderSpecs := []seedOrderSpec{
			{
				Key:         "orderCompleted1",
				UserKey:     "customerA",
				PlayerKey:   "playerA",
				GameKey:     "lol",
				ItemCode:    "escort-lol-solo",
				Title:       "欢迎体验 GameLink 陪玩",
				Description: "我们为您匹配了经验丰富的高胜率陪玩，尽情享受游戏时光吧！",
				Status:      model.OrderStatusCompleted,
				PriceCents:  19900,
				Currency:    model.CurrencyCNY,
				StartOffset: -3 * hour,
				Duration:    1 * hour,
			},
			{
				Key:         "orderInProgress1",
				UserKey:     "customerB",
				PlayerKey:   "playerA",
				GameKey:     "dota2",
				ItemCode:    "escort-dota2-solo",
				Title:       "高端局连胜陪玩",
				Description: "DOTA2 冠军陪练，助你提升团队协作。",
				Status:      model.OrderStatusInProgress,
				PriceCents:  29900,
				Currency:    model.CurrencyCNY,
				StartOffset: -1 * hour,
				Duration:    2 * hour,
			},
			{
				Key:         "orderPending1",
				UserKey:     "customerC",
				PlayerKey:   "playerA",
				GameKey:     "lol",
				ItemCode:    "escort-lol-solo",
				Title:       "黄金段位冲刺",
				Description: "等待分配陪玩师，预计 30 分钟内开始。",
				Status:      model.OrderStatusPending,
				PriceCents:  15900,
				Currency:    model.CurrencyCNY,
				StartOffset: 1 * hour,
				Duration:    90 * time.Minute,
			},
			{
				Key:          "orderCanceled1",
				UserKey:      "customerB",
				PlayerKey:    "playerB",
				GameKey:      "valorant",
				ItemCode:     "escort-valorant-solo",
				Title:        "战术射击训练营",
				Description:  "因临时有事取消，等待重新安排。",
				Status:       model.OrderStatusCanceled,
				PriceCents:   12900,
				Currency:     model.CurrencyCNY,
				StartOffset:  -5 * hour,
				Duration:     2 * hour,
				CancelReason: "用户主动取消",
			},
			{
				Key:         "orderConfirmed1",
				UserKey:     "customerD",
				PlayerKey:   "playerC",
				GameKey:     "csgo",
				ItemCode:    "escort-default",
				Title:       "枪法强化训练",
				Description: "专业CS:GO选手带你提升枪法，从基础到进阶全覆盖。",
				Status:      model.OrderStatusConfirmed,
				PriceCents:  18900,
				Currency:    model.CurrencyCNY,
				StartOffset: 30 * time.Minute,
				Duration:    2 * hour,
			},
			{
				Key:         "orderPending2",
				UserKey:     "customerE",
				PlayerKey:   "playerB",
				GameKey:     "apex",
				ItemCode:    "escort-default",
				Title:       "大逃杀双人排位",
				Description: "寻找志同道合的队友一起冲击更高段位。",
				Status:      model.OrderStatusPending,
				PriceCents:  9900,
				Currency:    model.CurrencyCNY,
				StartOffset: 2 * hour,
				Duration:    3 * hour,
			},
			{
				Key:         "orderCompleted2",
				UserKey:     "customerF",
				PlayerKey:   "playerD",
				GameKey:     "wow",
				ItemCode:    "escort-wow-solo",
				Title:       "魔兽世界副本开荒",
				Description: "资深玩家带你体验经典副本，了解游戏机制和装备获取。",
				Status:      model.OrderStatusCompleted,
				PriceCents:  24900,
				Currency:    model.CurrencyCNY,
				StartOffset: -6 * hour,
				Duration:    4 * hour,
			},
			{
				Key:         "orderInProgress2",
				UserKey:     "customerG",
				PlayerKey:   "playerE",
				GameKey:     "fifa",
				ItemCode:    "escort-default",
				Title:       "FIFA在线友谊赛",
				Description: "体育游戏爱好者之间的友好比赛，享受竞技乐趣。",
				Status:      model.OrderStatusInProgress,
				PriceCents:  7900,
				Currency:    model.CurrencyCNY,
				StartOffset: -30 * time.Minute,
				Duration:    90 * time.Minute,
			},
			{
				Key:         "orderPending3",
				UserKey:     "customerH",
				PlayerKey:   "playerF",
				GameKey:     "fallguys",
				ItemCode:    "escort-default",
				Title:       "糖豆人欢乐时光",
				Description: "轻松愉快的派对游戏，适合休闲娱乐，放松心情。",
				Status:      model.OrderStatusPending,
				PriceCents:  5900,
				Currency:    model.CurrencyCNY,
				StartOffset: 4 * hour,
				Duration:    2 * hour,
			},
			{
				Key:          "orderRefunded1",
				UserKey:      "customerA",
				PlayerKey:    "playerB",
				GameKey:      "overwatch",
				ItemCode:     "escort-default",
				Title:        "守望先锋团队竞技",
				Description:  "因技术问题服务器维护，全额退款。",
				Status:       model.OrderStatusRefunded,
				PriceCents:   16900,
				Currency:     model.CurrencyCNY,
				StartOffset:  -2 * hour,
				Duration:     90 * time.Minute,
				CancelReason: "服务器维护退款",
			},
			{
				Key:         "orderConfirmed2",
				UserKey:     "customerB",
				PlayerKey:   "playerG",
				GameKey:     "dota2",
				ItemCode:    "escort-dota2-solo",
				Title:       "DOTA2新手教学",
				Description: "DOTA2老玩家带你熟悉游戏机制，学习基础操作和战术理解。",
				Status:      model.OrderStatusConfirmed,
				PriceCents:  21900,
				Currency:    model.CurrencyCNY,
				StartOffset: 1 * hour,
				Duration:    3 * hour,
			},
		}

		orders := make(map[string]*model.Order, len(orderSpecs))
		for _, spec := range orderSpecs {
			user, ok := users[spec.UserKey]
			if !ok {
				return fmt.Errorf("seed order missing user %s", spec.UserKey)
			}
			game, ok := games[spec.GameKey]
			if !ok {
				return fmt.Errorf("seed order missing game %s", spec.GameKey)
			}
			item := defaultServiceItem
			if spec.ItemCode != "" {
				if found, ok := serviceItems[spec.ItemCode]; ok && found != nil {
					item = found
				} else {
					log.Printf("warning: seed order %s references missing service item %q, fallback to escort-default", spec.Key, spec.ItemCode)
				}
			}
			var playerID *uint64
			if spec.PlayerKey != "" {
				player, ok := players[spec.PlayerKey]
				if !ok {
					return fmt.Errorf("seed order missing player %s", spec.PlayerKey)
				}
				playerID = &player.ID
			}
			var startPtr, endPtr *time.Time
			if spec.StartOffset != 0 || spec.Duration != 0 {
				startTime := now.Add(spec.StartOffset)
				startPtr = ptrTime(startTime)
				if spec.Duration != 0 {
					endPtr = ptrTime(startTime.Add(spec.Duration))
				}
			}
			var startedAt, completedAt *time.Time
			switch spec.Status {
			case model.OrderStatusInProgress, model.OrderStatusCompleted:
				startedAt = startPtr
			}
			if spec.Status == model.OrderStatusCompleted {
				completedAt = endPtr
			}
			order, err := seedOrder(tx, seedOrderParams{
				Title:          spec.Title,
				Description:    spec.Description,
				UserID:         user.ID,
				PlayerID:       playerID,
				ItemID:         item.ID,
				GameID:         game.ID,
				Status:         spec.Status,
				PriceCents:     spec.PriceCents,
				Currency:       spec.Currency,
				ScheduledStart: startPtr,
				ScheduledEnd:   endPtr,
				CancelReason:   spec.CancelReason,
				StartedAt:      startedAt,
				CompletedAt:    completedAt,
			})
			if err != nil {
				return err
			}
			orders[spec.Key] = order
		}

		paymentSpecs := []seedPaymentSpec{
			{
				OrderKey:        "orderCompleted1",
				UserKey:         "customerA",
				Method:          model.PaymentMethodWeChat,
				AmountCents:     19900,
				Currency:        model.CurrencyCNY,
				Status:          model.PaymentStatusPaid,
				ProviderTradeNo: "WX1234567890",
				ProviderRaw:     `{"result":"success"}`,
				PaidAtOffset:    ptrDuration(-2 * hour),
			},
			{
				OrderKey:        "orderInProgress1",
				UserKey:         "customerB",
				Method:          model.PaymentMethodAlipay,
				AmountCents:     29900,
				Currency:        model.CurrencyCNY,
				Status:          model.PaymentStatusPending,
				ProviderTradeNo: "ALI987654321",
				ProviderRaw:     `{"result":"processing"}`,
			},
			{
				OrderKey:        "orderCanceled1",
				UserKey:         "customerB",
				Method:          model.PaymentMethodWeChat,
				AmountCents:     12900,
				Currency:        model.CurrencyCNY,
				Status:          model.PaymentStatusRefunded,
				ProviderTradeNo: "WXREFUND001",
				ProviderRaw:     `{"result":"refunded"}`,
				PaidAtOffset:    ptrDuration(-5 * hour),
				RefundedOffset:  ptrDuration(-4 * hour),
			},
			{
				OrderKey:        "orderConfirmed1",
				UserKey:         "customerD",
				Method:          model.PaymentMethodWeChat,
				AmountCents:     18900,
				Currency:        model.CurrencyCNY,
				Status:          model.PaymentStatusPaid,
				ProviderTradeNo: "WXTRAIN123",
				ProviderRaw:     `{"result":"success"}`,
				PaidAtOffset:    ptrDuration(10 * time.Minute),
			},
			{
				OrderKey:        "orderCompleted2",
				UserKey:         "customerF",
				Method:          model.PaymentMethodAlipay,
				AmountCents:     24900,
				Currency:        model.CurrencyCNY,
				Status:          model.PaymentStatusPaid,
				ProviderTradeNo: "ALIWOWEXP456",
				ProviderRaw:     `{"result":"success"}`,
				PaidAtOffset:    ptrDuration(-7 * hour),
			},
			{
				OrderKey:        "orderInProgress2",
				UserKey:         "customerG",
				Method:          model.PaymentMethodWeChat,
				AmountCents:     7900,
				Currency:        model.CurrencyCNY,
				Status:          model.PaymentStatusPaid,
				ProviderTradeNo: "WXSPORTS789",
				ProviderRaw:     `{"result":"success"}`,
				PaidAtOffset:    ptrDuration(-45 * time.Minute),
			},
			{
				OrderKey:        "orderRefunded1",
				UserKey:         "customerA",
				Method:          model.PaymentMethodAlipay,
				AmountCents:     16900,
				Currency:        model.CurrencyCNY,
				Status:          model.PaymentStatusRefunded,
				ProviderTradeNo: "ALIREPAIR001",
				ProviderRaw:     `{"result":"refunded"}`,
				PaidAtOffset:    ptrDuration(-2 * hour),
				RefundedOffset:  ptrDuration(-90 * time.Minute),
			},
			{
				OrderKey:        "orderConfirmed2",
				UserKey:         "customerB",
				Method:          model.PaymentMethodWeChat,
				AmountCents:     21900,
				Currency:        model.CurrencyCNY,
				Status:          model.PaymentStatusPaid,
				ProviderTradeNo: "WXDOTATEACH012",
				ProviderRaw:     `{"result":"success"}`,
				PaidAtOffset:    ptrDuration(30 * time.Minute),
			},
		}

		for _, spec := range paymentSpecs {
			order, ok := orders[spec.OrderKey]
			if !ok {
				return fmt.Errorf("seed payment missing order %s", spec.OrderKey)
			}
			user, ok := users[spec.UserKey]
			if !ok {
				return fmt.Errorf("seed payment missing user %s", spec.UserKey)
			}
			paidAt := ptrTimeWithOffset(now, spec.PaidAtOffset)
			refundedAt := ptrTimeWithOffset(now, spec.RefundedOffset)
			if err := seedPayment(tx, seedPaymentParams{
				OrderID:         order.ID,
				UserID:          user.ID,
				Method:          spec.Method,
				AmountCents:     spec.AmountCents,
				Currency:        spec.Currency,
				Status:          spec.Status,
				ProviderTradeNo: spec.ProviderTradeNo,
				ProviderRaw:     json.RawMessage(spec.ProviderRaw),
				PaidAt:          paidAt,
				RefundedAt:      refundedAt,
			}); err != nil {
				return err
			}
		}

		reviewSpecs := []seedReviewSpec{
			// 已通过的好评
			{
				OrderKey:  "orderCompleted1",
				UserKey:   "customerA",
				PlayerKey: "playerA",
				Score:     model.MustRating(5),
				Content:   "很满意的陪玩体验，带我连胜！峡谷守护者技术确实强，打野节奏把控很好。推荐给大家！",
				Status:    model.ReviewStatusApproved,
				Images:    []string{"https://gamelink.oss.com/reviews/lol_win_1.jpg", "https://gamelink.oss.com/reviews/lol_win_2.jpg"},
			},
			{
				OrderKey:  "orderCompleted2",
				UserKey:   "customerF",
				PlayerKey: "playerD",
				Score:     model.MustRating(5),
				Content:   "异世界旅者对MMORPG的理解非常深入，带我了解了魔兽世界的核心玩法，收益很大！副本讲解清晰，装备搭配建议很实用。",
				Status:    model.ReviewStatusApproved,
				Images:    []string{"https://gamelink.oss.com/reviews/wow_raid_1.jpg"},
			},
			// 待审核的评价
			{
				OrderKey:  "orderInProgress1",
				UserKey:   "customerB",
				PlayerKey: "playerA",
				Score:     model.MustRating(4),
				Content:   "战术指导很专业，期待后续完成。DOTA2的复杂度很高，有专业指导确实不一样。",
				Status:    model.ReviewStatusPending,
			},
			{
				OrderKey:  "orderInProgress2",
				UserKey:   "customerG",
				PlayerKey: "playerE",
				Score:     model.MustRating(4),
				Content:   "运动健将的足球水平很高，学到了很多实用的技巧。FIFA游戏体验很好，传球和射门技巧讲解到位。",
				Status:    model.ReviewStatusPending,
				Images:    []string{"https://gamelink.oss.com/reviews/fifa_goal_1.jpg"},
			},
			{
				OrderKey:  "orderConfirmed2",
				UserKey:   "customerB",
				PlayerKey: "playerG",
				Score:     model.MustRating(5),
				Content:   "DOTA宗师的教学很系统，从基础到进阶都有涉及，受益匪浅。英雄池推荐和对线技巧都很实用！",
				Status:    model.ReviewStatusPending,
			},
			// 被拒绝的评价
			{
				OrderKey:        "orderRefunded1",
				UserKey:         "customerA",
				PlayerKey:       "playerB",
				Score:           model.MustRating(2),
				Content:         "服务态度一般，感觉不太专业。",
				Status:          model.ReviewStatusRejected,
				RejectionReason: "评价内容过于简短，缺乏具体描述",
			},
			// 被举报的评价
			{
				OrderKey:   "orderConfirmed1",
				UserKey:    "customerD",
				PlayerKey:  "playerC",
				Score:      model.MustRating(5),
				Content:    "枪神降临不愧是职业选手，枪法精准，教学耐心细致。CS:GO的水平确实提升了很多！强烈推荐！",
				Status:     model.ReviewStatusApproved,
				IsReported: true,
				Images:     []string{"https://gamelink.oss.com/reviews/csgo_ace_1.jpg", "https://gamelink.oss.com/reviews/csgo_ace_2.jpg"},
			},
			// 中等评价
			{
				OrderKey:  "orderPending3",
				UserKey:   "customerH",
				PlayerKey: "playerF",
				Score:     model.MustRating(3),
				Content:   "整体体验一般，陪玩师技术还可以，但沟通上有些问题，回复不够及时。希望能改进。",
				Status:    model.ReviewStatusApproved,
			},
			// 带图片的好评
			{
				OrderKey:  "orderPending2",
				UserKey:   "customerE",
				PlayerKey: "playerB",
				Score:     model.MustRating(5),
				Content:   "王牌射手带我打Apex双排上分，枪线和走位讲得很细，实战提升明显！体验很棒。",
				Status:    model.ReviewStatusPending,
				Images:    []string{"https://gamelink.oss.com/reviews/apex_rank_1.jpg", "https://gamelink.oss.com/reviews/apex_win_1.jpg"},
			},
		}

		for i, spec := range reviewSpecs {
			order, ok := orders[spec.OrderKey]
			if !ok {
				return fmt.Errorf("seed review missing order %s", spec.OrderKey)
			}
			user, ok := users[spec.UserKey]
			if !ok {
				return fmt.Errorf("seed review missing user %s", spec.UserKey)
			}
			player, ok := players[spec.PlayerKey]
			if !ok {
				return fmt.Errorf("seed review missing player %s", spec.PlayerKey)
			}
			log.Printf("seed review[%d]: OrderKey=%s (id=%d, user_id=%d), UserKey=%s (id=%d), PlayerKey=%s (id=%d)",
				i, spec.OrderKey, order.ID, order.UserID, spec.UserKey, user.ID, spec.PlayerKey, player.ID)
			if err := seedReview(tx, seedReviewParams{
				OrderID:         order.ID,
				UserID:          user.ID,
				PlayerID:        player.ID,
				Score:           spec.Score,
				Content:         spec.Content,
				Status:          spec.Status,
				Images:          spec.Images,
				IsReported:      spec.IsReported,
				RejectionReason: spec.RejectionReason,
			}); err != nil {
				return err
			}
		}

		// 菜单种子数据
		if err := seedUserManagementData(tx, users); err != nil {
			return err
		}
		if err := seedMenus(tx); err != nil {
			return err
		}

		// 监控模块种子数据
		if err := seedMonitorData(tx); err != nil {
			return err
		}

		// 评价管理权限种子数据
		if err := seedReviewPermissions(tx); err != nil {
			return err
		}

		// 内容管理种子数据
		if err := seedContentData(tx, users); err != nil {
			return err
		}

		// 系统权限种子数据 (必须在 seedDefaultRoles 之前执行)
		if err := seedSystemPermissions(tx); err != nil {
			return err
		}

		// 默认角色和权限种子数据
		if err := seedDefaultRoles(tx); err != nil {
			return err
		}

		// 提现种子数据
		if err := seedWithdrawData(tx, players); err != nil {
			return err
		}

		// 评价举报种子数据
		if err := seedReviewReportData(tx, users); err != nil {
			return err
		}

		// 敏感词种子数据
		if err := seedSensitiveWords(tx); err != nil {
			return err
		}

		// 钱包种子数据（陪玩师需要有余额才能提现）
		if err := seedWalletData(tx, users, players); err != nil {
			return err
		}

		// 佣金规则种子数据
		if err := seedCommissionRules(tx, games); err != nil {
			return err
		}

		// 佣金记录种子数据
		if err := seedCommissionRecords(tx, orders, players); err != nil {
			return err
		}

		// 收款主体和分流规则种子数据
		if err := seedCollectionEntities(tx); err != nil {
			return err
		}

		// 结算公司种子数据（传入players以创建分配关系）
		if err := seedSettlementCompanies(tx, players); err != nil {
			return err
		}

		// 排行榜抽成配置种子数据
		if err := seedRankingCommissionConfigs(tx); err != nil {
			return err
		}

		// 订单纠纷种子数据
		if err := seedOrderDisputes(tx, orders, users); err != nil {
			return err
		}

		couponTemplates, vipLevels, err := seedVipAndCouponTemplates(tx, games, serviceItems)
		if err != nil {
			return err
		}

		if err := seedUserVipState(tx, now, users, vipLevels); err != nil {
			return err
		}

		coupons, err := seedUserCoupons(tx, couponTemplates, users, orders)
		if err != nil {
			return err
		}

		activities, err := seedActivityData(tx, couponTemplates, users)
		if err != nil {
			return err
		}

		if err := seedRechargeData(tx, couponTemplates, users); err != nil {
			return err
		}

		if err := seedTeamData(tx, players, orders); err != nil {
			return err
		}

		if err := seedReferralData(tx, couponTemplates, users); err != nil {
			return err
		}

		if err := seedUserBlockData(tx, users); err != nil {
			return err
		}

		if err := seedGameRankAndCertificationData(tx, games, players, users); err != nil {
			return err
		}

		if err := seedNotificationData(tx, users, orders, coupons, activities, vipLevels); err != nil {
			return err
		}

		if err := seedAdditionalFlowOrders(tx, now, users, players, serviceItems, orders); err != nil {
			return err
		}

		if err := seedRefundAndTimeoutData(tx, now, users); err != nil {
			return err
		}

		if err := seedOrderChatAndServiceAssignment(tx, now, users, players, orders); err != nil {
			return err
		}

		if err := validateSeedAssociations(tx); err != nil {
			return err
		}

		log.Println("seed data ensured for demo environment")
		return nil
	})
}

type seedPlayerParams struct {
	UserID             uint64
	Nickname           string
	Bio                string
	RatingAverage      float32
	RatingCount        uint32
	HourlyRateCents    int64
	MainGameID         uint64
	VerificationStatus model.VerificationStatus
}

type seedOrderParams struct {
	Title             string
	Description       string
	UserID            uint64
	PlayerID          *uint64
	ItemID            uint64
	GameID            uint64
	Status            model.OrderStatus
	PriceCents        int64
	Currency          model.Currency
	ScheduledStart    *time.Time
	ScheduledEnd      *time.Time
	CancelReason      string
	StartedAt         *time.Time
	CompletedAt       *time.Time
	RefundAmountCents *int64
	RefundReason      string
	RefundedAt        *time.Time
}

type seedPaymentParams struct {
	OrderID         uint64
	UserID          uint64
	Method          model.PaymentMethod
	AmountCents     int64
	Currency        model.Currency
	Status          model.PaymentStatus
	ProviderTradeNo string
	ProviderRaw     json.RawMessage
	PaidAt          *time.Time
	RefundedAt      *time.Time

	// Combined payment fields (optional)
	WalletAmountCents     int64
	ThirdPartyMethod      model.PaymentMethod
	ThirdPartyAmountCents int64
}

type seedReviewParams struct {
	OrderID         uint64
	UserID          uint64
	PlayerID        uint64
	Score           model.Rating
	Content         string
	Status          model.ReviewStatus
	Images          []string
	IsReported      bool
	RejectionReason string
}

func seedGames(tx *gorm.DB) (map[string]*model.Game, error) {
	seeds := []model.Game{
		{Key: "lol", Name: "英雄联盟", Category: "moba", Description: "召唤师峡谷 5v5 对战"},
		{Key: "dota2", Name: "DOTA 2", Category: "moba", Description: "经典即时战略竞技"},
		{Key: "valorant", Name: "无畏契约", Category: "fps", Description: "英雄战术射击"},
		{Key: "csgo", Name: "反恐精英：全球攻势", Category: "fps", Description: "经典第一人称射击"},
		{Key: "apex", Name: "Apex英雄", Category: "fps", Description: "大逃杀类射击游戏"},
		{Key: "pubg", Name: "绝地求生", Category: "fps", Description: "百人竞技生存游戏"},
		{Key: "overwatch", Name: "守望先锋", Category: "fps", Description: "团队英雄射击游戏"},
		{Key: "fifa", Name: "FIFA足球", Category: "sports", Description: "足球模拟游戏"},
		{Key: "nba2k", Name: "NBA2K篮球", Category: "sports", Description: "篮球模拟游戏"},
		{Key: "wzry", Name: "王者荣耀", Category: "moba", Description: "移动端MOBA游戏"},
		{Key: "genshin", Name: "原神", Category: "rpg", Description: "开放世界冒险游戏"},
		{Key: "wow", Name: "魔兽世界", Category: "rpg", Description: "大型多人在线角色扮演游戏"},
		{Key: "minecraft", Name: "我的世界", Category: "sandbox", Description: "沙盒建造游戏"},
		{Key: "amongus", Name: "Among Us", Category: "social", Description: "社交推理游戏"},
		{Key: "fallguys", Name: "糖豆人", Category: "party", Description: "多人派对游戏"},
	}
	result := make(map[string]*model.Game, len(seeds))
	for i := range seeds {
		game := &seeds[i]
		var existing model.Game
		// 使用 Unscoped 包含软删除的记录，避免唯一索引冲突
		if err := tx.Unscoped().Where("key = ?", game.Key).First(&existing).Error; err == nil {
			// 如果记录被软删除，恢复它
			if existing.DeletedAt.Valid {
				if err := tx.Unscoped().Model(&existing).Update("deleted_at", nil).Error; err != nil {
					return nil, err
				}
			}
			ex := existing
			result[game.Key] = &ex
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if err := tx.Create(game).Error; err != nil {
			return nil, err
		}
		result[game.Key] = game
	}
	return result, nil
}

func seedUser(tx *gorm.DB, input seedUserInput) (*model.User, error) {
	if input.Email == "" && input.Phone == "" {
		return nil, errors.New("seed user requires email or phone")
	}
	// 使用 Unscoped 包含软删除的记录，避免唯一索引冲突
	lookup := tx.Unscoped().Model(&model.User{})
	if input.Email != "" {
		lookup = lookup.Where("email = ?", input.Email)
	} else {
		lookup = lookup.Where("phone = ?", input.Phone)
	}
	var existing model.User
	if err := lookup.First(&existing).Error; err == nil {
		// 如果记录被软删除，恢复它
		if existing.DeletedAt.Valid {
			if err := tx.Unscoped().Model(&existing).Update("deleted_at", nil).Error; err != nil {
				return nil, err
			}
		}
		return &existing, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Email:        input.Email,
		Phone:        input.Phone,
		Name:         input.Name,
		Role:         input.Role,
		Status:       model.UserStatusActive,
		PasswordHash: string(hashed),
	}
	if err := tx.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func seedPlayer(tx *gorm.DB, input seedPlayerParams) (*model.Player, error) {
	var existing model.Player
	if err := tx.Where("user_id = ?", input.UserID).First(&existing).Error; err == nil {
		return &existing, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	player := &model.Player{
		UserID:             input.UserID,
		Nickname:           input.Nickname,
		Bio:                input.Bio,
		RatingAverage:      input.RatingAverage,
		RatingCount:        input.RatingCount,
		HourlyRateCents:    input.HourlyRateCents,
		MainGameID:         input.MainGameID,
		VerificationStatus: input.VerificationStatus,
	}
	if err := tx.Create(player).Error; err != nil {
		return nil, err
	}
	return player, nil
}

func ensureServiceItem(tx *gorm.DB, code, name string, gameID uint64) (*model.ServiceItem, error) {
	var item model.ServiceItem
	if err := tx.Where("item_code = ?", code).First(&item).Error; err == nil {
		return &item, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	item = model.ServiceItem{
		ItemCode:       code,
		Name:           name,
		Description:    "系统默认护航服务项",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		GameID:         &gameID,
		BasePriceCents: 9900,
		ServiceHours:   1,
		CommissionRate: 0.20,
		IsActive:       true,
		Tags:           "[]", // 空 JSON 数组，PostgreSQL JSON 类型不接受空字符串
	}
	if err := tx.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func seedOrder(tx *gorm.DB, input seedOrderParams) (*model.Order, error) {
	var existing model.Order
	if err := tx.Where("title = ? AND user_id = ?", input.Title, input.UserID).First(&existing).Error; err == nil {
		return &existing, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	order := &model.Order{
		OrderNo:         model.GenerateEscortOrderNo(),
		UserID:          input.UserID,
		ItemID:          input.ItemID,
		GameID:          &input.GameID,
		Quantity:        1,
		UnitPriceCents:  input.PriceCents,
		TotalPriceCents: input.PriceCents,
		Currency:        input.Currency,
		Status:          input.Status,
		Title:           input.Title,
		Description:     input.Description,
		ScheduledStart:  input.ScheduledStart,
		ScheduledEnd:    input.ScheduledEnd,
		CancelReason:    strings.TrimSpace(input.CancelReason),
		OrderConfig:     "{}", // 空 JSON 对象，PostgreSQL JSON 类型不接受空字符串
	}
	if input.PlayerID != nil {
		order.PlayerID = input.PlayerID
	}
	if input.StartedAt != nil {
		order.StartedAt = input.StartedAt
	}
	if input.CompletedAt != nil {
		order.CompletedAt = input.CompletedAt
	}
	if input.RefundAmountCents != nil {
		order.RefundAmountCents = *input.RefundAmountCents
	}
	if input.RefundReason != "" {
		order.RefundReason = strings.TrimSpace(input.RefundReason)
	}
	if input.RefundedAt != nil {
		order.RefundedAt = input.RefundedAt
	}
	if err := tx.Create(order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

func seedPayment(tx *gorm.DB, input seedPaymentParams) error {
	var existing model.Payment
	if err := tx.Where("order_id = ?", input.OrderID).First(&existing).Error; err == nil {
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	payment := &model.Payment{
		OrderID:         input.OrderID,
		UserID:          input.UserID,
		Method:          input.Method,
		AmountCents:     input.AmountCents,
		Currency:        input.Currency,
		Status:          input.Status,
		ProviderTradeNo: input.ProviderTradeNo,
		ProviderRaw:     input.ProviderRaw,
		PaidAt:          input.PaidAt,
		RefundedAt:      input.RefundedAt,

		WalletAmountCents:     input.WalletAmountCents,
		ThirdPartyMethod:      input.ThirdPartyMethod,
		ThirdPartyAmountCents: input.ThirdPartyAmountCents,
	}
	if err := tx.Create(payment).Error; err != nil {
		return err
	}
	return nil
}

func seedReview(tx *gorm.DB, input seedReviewParams) error {
	var existing model.Review
	if err := tx.Where("order_id = ?", input.OrderID).First(&existing).Error; err == nil {
		// Delete and recreate to ensure data consistency (user_id or player_id may have changed)
		log.Printf("seedReview: deleting stale review id=%d (order_id=%d, old user_id=%d) to recreate with correct data",
			existing.ID, existing.OrderID, existing.UserID)
		if err := tx.Delete(&existing).Error; err != nil {
			return err
		}
		// Continue to create new review below
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	status := input.Status
	if status == "" {
		status = model.ReviewStatusApproved
	}
	review := &model.Review{
		OrderID:         input.OrderID,
		UserID:          input.UserID,
		PlayerID:        input.PlayerID,
		Score:           input.Score,
		Content:         input.Content,
		Status:          status,
		Images:          input.Images,
		IsReported:      input.IsReported,
		RejectionReason: input.RejectionReason,
	}
	if err := tx.Create(review).Error; err != nil {
		return err
	}
	log.Printf("seedReview: created review id=%d (order_id=%d, user_id=%d)", review.ID, review.OrderID, review.UserID)
	return nil
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

func ptrDuration(d time.Duration) *time.Duration {
	return &d
}

func ptrTimeWithOffset(base time.Time, offset *time.Duration) *time.Time {
	if offset == nil {
		return nil
	}
	return ptrTime(base.Add(*offset))
}

// seedMenus 创建后台管理菜单
// 注意：菜单数据现在由前端初始化服务同步，不再在后端种子数据中创建
func seedMenus(_ *gorm.DB) error {
	// 菜单数据由前端 init.ts 服务同步到后端
	// 当管理员首次登录时，前端会自动将 ADMIN_MENUS 配置同步到数据库
	log.Println("menu seed skipped - menus are synced from frontend")
	return nil
}

// seedMonitorData 创建监控模块示例数据
func seedMonitorData(tx *gorm.DB) error {
	// 创建示例告警
	var alertCount int64
	if err := tx.Model(&model.Alert{}).Count(&alertCount).Error; err != nil {
		return err
	}
	if alertCount == 0 {
		alerts := []model.Alert{
			{
				Level:   model.AlertLevelHigh,
				Type:    model.AlertTypeSystem,
				Title:   "CPU 使用率过高",
				Message: "服务器 CPU 使用率超过 90%，请检查系统负载",
				Source:  "system-monitor",
				IsRead:  false,
			},
			{
				Level:   model.AlertLevelMedium,
				Type:    model.AlertTypeBusiness,
				Title:   "订单处理延迟",
				Message: "订单处理队列积压，平均等待时间超过 5 分钟",
				Source:  "order-service",
				IsRead:  false,
			},
			{
				Level:   model.AlertLevelLow,
				Type:    model.AlertTypeSecurity,
				Title:   "异常登录尝试",
				Message: "检测到来自异常 IP 的多次登录尝试",
				Source:  "security-monitor",
				IsRead:  true,
			},
			{
				Level:   model.AlertLevelMedium,
				Type:    model.AlertTypeSystem,
				Title:   "数据库连接池接近上限",
				Message: "数据库连接池使用率达到 85%",
				Source:  "db-monitor",
				IsRead:  false,
			},
			{
				Level:   model.AlertLevelLow,
				Type:    model.AlertTypeBusiness,
				Title:   "新陪玩师注册待审核",
				Message: "有 5 位新陪玩师等待审核认证",
				Source:  "player-service",
				IsRead:  true,
			},
		}
		for _, alert := range alerts {
			if err := tx.Create(&alert).Error; err != nil {
				return err
			}
		}
		log.Println("alert seed data created")
	}

	// 创建 KPI 目标
	var kpiCount int64
	if err := tx.Model(&model.KPITarget{}).Count(&kpiCount).Error; err != nil {
		return err
	}
	if kpiCount == 0 {
		now := time.Now()
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		monthEnd := monthStart.AddDate(0, 1, -1)

		kpiTargets := []model.KPITarget{
			{
				PeriodType:  "monthly",
				MetricName:  "gmv",
				TargetValue: 1000000, // 100万GMV
				StartDate:   monthStart,
				EndDate:     monthEnd,
				CreatedBy:   1,
			},
			{
				PeriodType:  "monthly",
				MetricName:  "orders",
				TargetValue: 5000, // 5000订单
				StartDate:   monthStart,
				EndDate:     monthEnd,
				CreatedBy:   1,
			},
			{
				PeriodType:  "monthly",
				MetricName:  "new_users",
				TargetValue: 1000, // 1000新用户
				StartDate:   monthStart,
				EndDate:     monthEnd,
				CreatedBy:   1,
			},
			{
				PeriodType:  "monthly",
				MetricName:  "new_players",
				TargetValue: 100, // 100新陪玩师
				StartDate:   monthStart,
				EndDate:     monthEnd,
				CreatedBy:   1,
			},
			{
				PeriodType:  "monthly",
				MetricName:  "dau",
				TargetValue: 500, // 500日活
				StartDate:   monthStart,
				EndDate:     monthEnd,
				CreatedBy:   1,
			},
			{
				PeriodType:  "monthly",
				MetricName:  "retention",
				TargetValue: 40, // 40%留存率
				StartDate:   monthStart,
				EndDate:     monthEnd,
				CreatedBy:   1,
			},
			{
				PeriodType:  "monthly",
				MetricName:  "repurchase",
				TargetValue: 30, // 30%复购率
				StartDate:   monthStart,
				EndDate:     monthEnd,
				CreatedBy:   1,
			},
		}
		for _, target := range kpiTargets {
			if err := tx.Create(&target).Error; err != nil {
				return err
			}
		}
		log.Println("KPI target seed data created")
	}

	return nil
}

// seedUserManagementData 创建用户管理模块种子数据
func seedUserManagementData(tx *gorm.DB, users map[string]*model.User) error {
	// 检查表是否存在
	if !tx.Migrator().HasTable(&model.UserTag{}) {
		// 如果表不存在，先创建表
		if err := tx.AutoMigrate(&model.UserTag{}, &model.UserTagRelation{}, &model.UserLoginHistory{}, &model.UserBehavior{}); err != nil {
			return fmt.Errorf("failed to migrate user management tables: %w", err)
		}
		log.Println("created user management tables")
	}

	// 检查是否已有用户标签数据
	var tagCount int64
	if err := tx.Model(&model.UserTag{}).Count(&tagCount).Error; err != nil {
		return err
	}

	if tagCount > 0 {
		log.Println("user management seed data already exists, skipping")
		return nil
	}

	// 创建用户标签
	tags := []struct {
		name        string
		color       string
		description string
		key         string // 用于后续关联
	}{
		{"VIP用户", "#FFD700", "付费高级用户", "vip"},
		{"活跃用户", "#4CAF50", "近期有登录记录的用户", "active"},
		{"新用户", "#2196F3", "注册30天内的用户", "new"},
		{"高消费用户", "#FF5722", "累计消费超过1000元的用户", "highspend"},
		{"陪玩师", "#9C27B0", "认证的游戏陪玩师", "player"},
		{"潜力用户", "#00BCD4", "有消费意向但未下单的用户", "potential"},
	}

	tagModels := make(map[string]*model.UserTag)
	for _, t := range tags {
		tag := &model.UserTag{
			Name:        t.name,
			Color:       t.color,
			Description: t.description,
		}
		if err := tx.Create(tag).Error; err != nil {
			return err
		}
		tagModels[t.key] = tag
		log.Printf("created user tag: %s\n", t.name)
	}

	// 为用户分配标签
	tagAssignments := []struct {
		userKey string
		tagKeys []string
	}{
		// VIP用户
		{"customerA", []string{"vip", "active"}},
		{"customerB", []string{"vip", "active", "highspend"}},
		{"customerH", []string{"vip", "active"}},

		// 活跃用户
		{"proA", []string{"active", "player"}},
		{"proB", []string{"active", "player"}},
		{"proC", []string{"active", "player"}},
		{"customerD", []string{"active"}},
		{"customerF", []string{"active"}},

		// 新用户
		{"customerG", []string{"new", "potential"}},
		{"customerE", []string{"new"}},

		// 高消费用户
		{"customerC", []string{"highspend"}},
	}

	for _, assignment := range tagAssignments {
		user, ok := users[assignment.userKey]
		if !ok {
			log.Printf("warning: user %s not found, skipping tag assignment\n", assignment.userKey)
			continue
		}

		for _, tagKey := range assignment.tagKeys {
			tag, ok := tagModels[tagKey]
			if !ok {
				log.Printf("warning: tag %s not found\n", tagKey)
				continue
			}

			relation := &model.UserTagRelation{
				UserID: user.ID,
				TagID:  tag.ID,
			}
			if err := tx.Create(relation).Error; err != nil {
				return err
			}
		}
		log.Printf("assigned tags to user %s\n", user.Name)
	}

	// 创建用户登录历史
	loginHistory := []model.UserLoginHistory{
		{UserID: users["customerA"].ID, IPAddress: "192.168.1.100", UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36", Location: "北京市", DeviceType: "desktop"},
		{UserID: users["customerA"].ID, IPAddress: "192.168.1.101", UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X)", Location: "上海市", DeviceType: "mobile"},
		{UserID: users["proA"].ID, IPAddress: "192.168.1.102", UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36", Location: "广州市", DeviceType: "desktop"},
		{UserID: users["proB"].ID, IPAddress: "192.168.1.103", UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)", Location: "深圳市", DeviceType: "desktop"},
		{UserID: users["customerB"].ID, IPAddress: "192.168.1.104", UserAgent: "Mozilla/5.0 (iPad; CPU OS 14_6 like Mac OS X)", Location: "成都市", DeviceType: "tablet"},
	}

	for _, history := range loginHistory {
		if err := tx.Create(&history).Error; err != nil {
			return err
		}
	}
	log.Println("created user login history")

	// 创建用户行为数据
	behaviors := []model.UserBehavior{
		{UserID: users["customerA"].ID, Action: "view_order", TargetType: "order", TargetID: 1, Metadata: `{"page": "order_detail"}`},
		{UserID: users["customerA"].ID, Action: "create_order", TargetType: "order", TargetID: 2, Metadata: `{"game": "dota2"}`},
		{UserID: users["customerB"].ID, Action: "view_player", TargetType: "player", TargetID: 1, Metadata: `{"player_id": 1}`},
		{UserID: users["customerB"].ID, Action: "payment", TargetType: "payment", TargetID: 1, Metadata: `{"method": "wechat", "amount": 29900}`},
		{UserID: users["proA"].ID, Action: "update_profile", TargetType: "player", TargetID: 1, Metadata: `{"field": "bio"}`},
		{UserID: users["proB"].ID, Action: "accept_order", TargetType: "order", TargetID: 1, Metadata: `{"status": "confirmed"}`},
		{UserID: users["customerC"].ID, Action: "view_game", TargetType: "game", TargetID: 1, Metadata: `{"game": "lol"}`},
		{UserID: users["customerC"].ID, Action: "search", TargetType: "search", TargetID: 0, Metadata: `{"keyword": "fps", "results": 15}`},
		{UserID: users["adminA"].ID, Action: "view_dashboard", TargetType: "admin", TargetID: 0, Metadata: `{"page": "dashboard"}`},
		{UserID: users["adminA"].ID, Action: "manage_user", TargetType: "user", TargetID: 2, Metadata: `{"action": "update_status"}`},
	}

	for _, behavior := range behaviors {
		if err := tx.Create(&behavior).Error; err != nil {
			return err
		}
	}
	log.Println("created user behavior data")

	log.Println("user management seed data created successfully")
	return nil
}

// seedReviewPermissions 创建所有业务模块权限种子数据
// 按业务模块组织，便于维护和扩展
func seedReviewPermissions(tx *gorm.DB) error {
	var allPermissions []model.Permission

	// 1. 评价管理模块
	allPermissions = append(allPermissions, getReviewPermissions()...)

	// 2. 内容管理模块
	allPermissions = append(allPermissions, getContentPermissions()...)

	// 3. 监控模块
	allPermissions = append(allPermissions, getMonitorPermissions()...)

	// 4. 运营分析模块
	allPermissions = append(allPermissions, getAnalyticsPermissions()...)

	// 5. KPI 仪表板模块
	allPermissions = append(allPermissions, getKPIPermissions()...)

	// 6. 财务管理模块
	allPermissions = append(allPermissions, getFinancePermissions()...)

	// 7. 系统管理模块
	allPermissions = append(allPermissions, getSystemPermissions()...)

	// 批量创建权限
	for _, perm := range allPermissions {
		if err := upsertPermission(tx, &perm); err != nil {
			return err
		}
	}

	log.Println("all business permissions seed data created successfully")
	return nil
}

// getReviewPermissions 评价管理模块权限
func getReviewPermissions() []model.Permission {
	return []model.Permission{
		// 评价查看
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/reviews", Code: "admin.reviews.list", Group: "评价管理", Description: "查看评价列表"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/reviews/:id", Code: "admin.reviews.read", Group: "评价管理", Description: "查看评价详情"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/reviews/pending", Code: "admin.reviews.pending.list", Group: "评价管理", Description: "查看待审核评价"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/reviews/:id/logs", Code: "admin.reviews.logs.list", Group: "评价管理", Description: "查看评价操作日志"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/players/:id/reviews", Code: "admin.reviews.player.list", Group: "评价管理", Description: "查看陪玩师评价"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/orders/:id/reviews", Code: "admin.reviews.order.list", Group: "评价管理", Description: "查看订单评价"},
		// 评价审核
		{Method: model.HTTPMethodPUT, Path: "/api/v1/admin/reviews/:id/approve", Code: "admin.reviews.approve.update", Group: "评价管理", Description: "批准评价"},
		{Method: model.HTTPMethodPUT, Path: "/api/v1/admin/reviews/:id/reject", Code: "admin.reviews.reject.update", Group: "评价管理", Description: "拒绝评价"},
		{Method: model.HTTPMethodPUT, Path: "/api/v1/admin/reviews/batch-approve", Code: "admin.reviews.batch-approve.update", Group: "评价管理", Description: "批量批准评价"},
		{Method: model.HTTPMethodPUT, Path: "/api/v1/admin/reviews/batch-reject", Code: "admin.reviews.batch-reject.update", Group: "评价管理", Description: "批量拒绝评价"},
		// 评价操作
		{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/reviews", Code: "admin.reviews.create", Group: "评价管理", Description: "创建评价"},
		{Method: model.HTTPMethodPUT, Path: "/api/v1/admin/reviews/:id", Code: "admin.reviews.update", Group: "评价管理", Description: "更新评价"},
		{Method: model.HTTPMethodDELETE, Path: "/api/v1/admin/reviews/:id", Code: "admin.reviews.delete", Group: "评价管理", Description: "删除评价"},
		// 评价统计
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/reviews/stats", Code: "admin.reviews.stats.list", Group: "评价管理", Description: "查看评价统计"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/reviews/trend", Code: "admin.reviews.trend.list", Group: "评价管理", Description: "查看评价趋势"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/reviews/top-players", Code: "admin.reviews.top-players.list", Group: "评价管理", Description: "查看陪玩师排行榜"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/reviews/game-stats", Code: "admin.reviews.game-stats.list", Group: "评价管理", Description: "查看游戏评价统计"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/reviews/export", Code: "admin.reviews.export", Group: "评价管理", Description: "导出评价统计"},
		// 评价设置
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/review-settings", Code: "admin.review-settings.list", Group: "评价管理", Description: "查看评价展示设置"},
		{Method: model.HTTPMethodPUT, Path: "/api/v1/admin/review-settings", Code: "admin.review-settings.update", Group: "评价管理", Description: "更新评价展示设置"},
		// 敏感词管理
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/sensitive-words", Code: "admin.sensitive-words.list", Group: "敏感词管理", Description: "查看敏感词列表"},
		{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/sensitive-words", Code: "admin.sensitive-words.create", Group: "敏感词管理", Description: "添加敏感词"},
		{Method: model.HTTPMethodPUT, Path: "/api/v1/admin/sensitive-words/:id", Code: "admin.sensitive-words.update", Group: "敏感词管理", Description: "更新敏感词"},
		{Method: model.HTTPMethodDELETE, Path: "/api/v1/admin/sensitive-words/:id", Code: "admin.sensitive-words.delete", Group: "敏感词管理", Description: "删除敏感词"},
		{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/reviews/detect-sensitive", Code: "review.detect_sensitive", Group: "敏感词管理", Description: "检测敏感词"},
		// 评价举报
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/review-reports", Code: "admin.review-reports.list", Group: "评价管理", Description: "查看评价举报列表"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/review-reports/:id", Code: "admin.review-reports.read", Group: "评价管理", Description: "查看举报详情"},
		{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/reviews/:id/reports", Code: "admin.review-reports.create", Group: "评价管理", Description: "创建评价举报"},
		{Method: model.HTTPMethodPUT, Path: "/api/v1/admin/review-reports/:id/handle", Code: "admin.review-reports.handle.update", Group: "评价管理", Description: "处理评价举报"},
		// 评价回复
		{Method: model.HTTPMethodPUT, Path: "/api/v1/admin/review-replies/:id", Code: "admin.review-replies.update", Group: "评价管理", Description: "更新评价回复"},
		{Method: model.HTTPMethodDELETE, Path: "/api/v1/admin/review-replies/:id", Code: "admin.review-replies.delete", Group: "评价管理", Description: "删除评价回复"},
	}
}

// getContentPermissions 内容管理模块权限
func getContentPermissions() []model.Permission {
	return []model.Permission{
		// 动态管理
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/content/feeds", Code: "content.feed.list", Group: "内容管理", Description: "查看动态列表"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/content/feeds/:id", Code: "content.feed.get", Group: "内容管理", Description: "查看动态详情"},
		{Method: model.HTTPMethodPUT, Path: "/api/v1/admin/content/feeds/:id/approve", Code: "content.feed.approve", Group: "内容管理", Description: "批准动态"},
		{Method: model.HTTPMethodPUT, Path: "/api/v1/admin/content/feeds/:id/reject", Code: "content.feed.reject", Group: "内容管理", Description: "拒绝动态"},
		{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/content/feeds/batch-approve", Code: "content.feed.batch_approve", Group: "内容管理", Description: "批量批准动态"},
		{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/content/feeds/batch-reject", Code: "content.feed.batch_reject", Group: "内容管理", Description: "批量拒绝动态"},
		{Method: model.HTTPMethodDELETE, Path: "/api/v1/admin/content/feeds/:id", Code: "content.feed.delete", Group: "内容管理", Description: "删除动态"},
		// 聊天监控
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/content/chat/messages", Code: "content.chat.list", Group: "内容管理", Description: "查看聊天消息"},
		{Method: model.HTTPMethodDELETE, Path: "/api/v1/admin/content/chat/messages/:id", Code: "content.chat.delete", Group: "内容管理", Description: "删除聊天消息"},
		{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/content/chat/mute", Code: "content.chat.mute", Group: "内容管理", Description: "禁言用户"},
		{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/content/chat/unmute", Code: "content.chat.unmute", Group: "内容管理", Description: "解除禁言"},
		// 举报管理
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/content/reports", Code: "content.report.list", Group: "内容管理", Description: "查看举报列表"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/content/reports/:id", Code: "content.report.get", Group: "内容管理", Description: "查看举报详情"},
		{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/content/reports/:id/process", Code: "content.report.process", Group: "内容管理", Description: "处理举报"},
		// 内容统计
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/content/stats", Code: "content.stats", Group: "内容管理", Description: "查看内容统计"},
		// 内容分类
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/content/categories", Code: "content.category.list", Group: "内容管理", Description: "查看内容分类"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/content/categories/:id", Code: "content.category.get", Group: "内容管理", Description: "查看分类详情"},
		{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/content/categories", Code: "content.category.create", Group: "内容管理", Description: "创建内容分类"},
		{Method: model.HTTPMethodPUT, Path: "/api/v1/admin/content/categories/:id", Code: "content.category.update", Group: "内容管理", Description: "更新内容分类"},
		{Method: model.HTTPMethodDELETE, Path: "/api/v1/admin/content/categories/:id", Code: "content.category.delete", Group: "内容管理", Description: "删除内容分类"},
	}
}

// getMonitorPermissions 监控模块权限
func getMonitorPermissions() []model.Permission {
	return []model.Permission{
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/monitor/system-status", Code: "monitor.system_status", Group: "系统监控", Description: "查看系统状态"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/monitor/online-users", Code: "monitor.online_users", Group: "系统监控", Description: "查看在线用户"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/monitor/order-queue", Code: "monitor.order_queue", Group: "系统监控", Description: "查看订单队列"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/monitor/alerts", Code: "monitor.alerts", Group: "系统监控", Description: "查看告警列表"},
		{Method: model.HTTPMethodPUT, Path: "/api/v1/admin/monitor/alerts/:id/read", Code: "monitor.alert_read", Group: "系统监控", Description: "标记告警已读"},
		{Method: model.HTTPMethodPUT, Path: "/api/v1/admin/monitor/alerts/batch-read", Code: "monitor.alert_batch_read", Group: "系统监控", Description: "批量标记告警已读"},
	}
}

// getAnalyticsPermissions 运营分析模块权限
func getAnalyticsPermissions() []model.Permission {
	return []model.Permission{
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/analytics/active-users", Code: "analytics.active_users", Group: "运营分析", Description: "查看活跃用户"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/analytics/retention", Code: "analytics.retention", Group: "运营分析", Description: "查看留存分析"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/analytics/payment", Code: "analytics.payment", Group: "运营分析", Description: "查看支付分析"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/analytics/conversion", Code: "analytics.conversion", Group: "运营分析", Description: "查看转化漏斗"},
	}
}

// getKPIPermissions KPI 仪表板模块权限
func getKPIPermissions() []model.Permission {
	return []model.Permission{
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/kpi/overview", Code: "kpi.overview", Group: "KPI管理", Description: "查看KPI概览"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/kpi/trend", Code: "kpi.trend", Group: "KPI管理", Description: "查看KPI趋势"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/kpi/targets", Code: "kpi.targets.list", Group: "KPI管理", Description: "查看KPI目标"},
		{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/kpi/targets", Code: "kpi.targets.create", Group: "KPI管理", Description: "创建KPI目标"},
		{Method: model.HTTPMethodPUT, Path: "/api/v1/admin/kpi/targets/:id", Code: "kpi.targets.update", Group: "KPI管理", Description: "更新KPI目标"},
		{Method: model.HTTPMethodDELETE, Path: "/api/v1/admin/kpi/targets/:id", Code: "kpi.targets.delete", Group: "KPI管理", Description: "删除KPI目标"},
	}
}

// getFinancePermissions 财务管理模块权限
func getFinancePermissions() []model.Permission {
	return []model.Permission{
		// 佣金管理
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/commission/rules", Code: "commission.rules.list", Group: "财务管理", Description: "查看佣金规则"},
		{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/commission/rules", Code: "commission.rules.create", Group: "财务管理", Description: "创建佣金规则"},
		{Method: model.HTTPMethodPUT, Path: "/api/v1/admin/commission/rules/:id", Code: "commission.rules.update", Group: "财务管理", Description: "更新佣金规则"},
		{Method: model.HTTPMethodDELETE, Path: "/api/v1/admin/commission/rules/:id", Code: "commission.rules.delete", Group: "财务管理", Description: "删除佣金规则"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/commission/records", Code: "commission.records.list", Group: "财务管理", Description: "查看佣金记录"},
		{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/commission/settlement/trigger", Code: "commission.settlement.trigger", Group: "财务管理", Description: "触发结算"},
		// 提现管理
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/withdraws", Code: "withdraw.list", Group: "财务管理", Description: "查看提现列表"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/withdraws/:id", Code: "withdraw.get", Group: "财务管理", Description: "查看提现详情"},
		{Method: model.HTTPMethodPUT, Path: "/api/v1/admin/withdraws/:id/approve", Code: "withdraw.approve", Group: "财务管理", Description: "批准提现"},
		{Method: model.HTTPMethodPUT, Path: "/api/v1/admin/withdraws/:id/reject", Code: "withdraw.reject", Group: "财务管理", Description: "拒绝提现"},
		// 收款主体
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/collection-entities", Code: "collection_entity.list", Group: "财务管理", Description: "查看收款主体"},
		{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/collection-entities", Code: "collection_entity.create", Group: "财务管理", Description: "创建收款主体"},
		{Method: model.HTTPMethodPUT, Path: "/api/v1/admin/collection-entities/:id", Code: "collection_entity.update", Group: "财务管理", Description: "更新收款主体"},
		{Method: model.HTTPMethodDELETE, Path: "/api/v1/admin/collection-entities/:id", Code: "collection_entity.delete", Group: "财务管理", Description: "删除收款主体"},
		// 分流规则
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/routing-rules", Code: "routing_rule.list", Group: "财务管理", Description: "查看分流规则"},
		{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/routing-rules", Code: "routing_rule.create", Group: "财务管理", Description: "创建分流规则"},
		{Method: model.HTTPMethodPUT, Path: "/api/v1/admin/routing-rules/:id", Code: "routing_rule.update", Group: "财务管理", Description: "更新分流规则"},
		{Method: model.HTTPMethodDELETE, Path: "/api/v1/admin/routing-rules/:id", Code: "routing_rule.delete", Group: "财务管理", Description: "删除分流规则"},
		// 排名抽成
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/ranking-commission", Code: "ranking_commission.list", Group: "财务管理", Description: "查看排名抽成"},
		{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/ranking-commission", Code: "ranking_commission.create", Group: "财务管理", Description: "创建排名抽成"},
		{Method: model.HTTPMethodPUT, Path: "/api/v1/admin/ranking-commission/:id", Code: "ranking_commission.update", Group: "财务管理", Description: "更新排名抽成"},
		{Method: model.HTTPMethodDELETE, Path: "/api/v1/admin/ranking-commission/:id", Code: "ranking_commission.delete", Group: "财务管理", Description: "删除排名抽成"},
	}
}

// getSystemPermissions 系统管理模块权限
func getSystemPermissions() []model.Permission {
	return []model.Permission{
		// 系统信息
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/system/info", Code: "system.info", Group: "系统管理", Description: "查看系统信息"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/system/health", Code: "system.health", Group: "系统管理", Description: "查看系统健康"},
		// 操作日志
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/operation-logs", Code: "operation_log.list", Group: "系统管理", Description: "查看操作日志"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/operation-logs/export", Code: "operation_log.export", Group: "系统管理", Description: "导出操作日志"},
		// 用户标签
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/user-tags", Code: "user_tag.list", Group: "系统管理", Description: "查看用户标签"},
		{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/user-tags", Code: "user_tag.create", Group: "系统管理", Description: "创建用户标签"},
		{Method: model.HTTPMethodPUT, Path: "/api/v1/admin/user-tags/:id", Code: "user_tag.update", Group: "系统管理", Description: "更新用户标签"},
		{Method: model.HTTPMethodDELETE, Path: "/api/v1/admin/user-tags/:id", Code: "user_tag.delete", Group: "系统管理", Description: "删除用户标签"},
		// 批量操作
		{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/users/batch/status", Code: "user_batch.status", Group: "系统管理", Description: "批量更新用户状态"},
		{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/users/batch/tags", Code: "user_batch.tags", Group: "系统管理", Description: "批量更新用户标签"},
		{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/users/batch/export", Code: "user_batch.export", Group: "系统管理", Description: "批量导出用户"},
	}
}

// upsertPermission 安全地插入或更新权限
// 检查 code 和 method+path 两个唯一约束，避免 PostgreSQL 事务中断
func upsertPermission(tx *gorm.DB, perm *model.Permission) error {
	// 1. 先检查 code 是否存在
	var existingByCode model.Permission
	if err := tx.Where("code = ?", perm.Code).First(&existingByCode).Error; err == nil {
		// code 已存在，跳过
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check permission by code %s: %w", perm.Code, err)
	}

	// 2. 再检查 method+path 是否存在
	var existingByPath model.Permission
	if err := tx.Where("method = ? AND path = ?", perm.Method, perm.Path).First(&existingByPath).Error; err == nil {
		// method+path 已存在但 code 不同，更新 code
		existingByPath.Code = perm.Code
		existingByPath.Description = perm.Description
		existingByPath.Group = perm.Group
		if err := tx.Save(&existingByPath).Error; err != nil {
			return fmt.Errorf("failed to update permission %s %s: %w", perm.Method, perm.Path, err)
		}
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check permission by path %s %s: %w", perm.Method, perm.Path, err)
	}

	// 3. 都不存在，创建新权限
	if err := tx.Create(perm).Error; err != nil {
		return fmt.Errorf("failed to create permission %s %s: %w", perm.Method, perm.Path, err)
	}
	return nil
}

// seedWithdrawData 创建提现种子数据
func seedWithdrawData(tx *gorm.DB, players map[string]*model.Player) error {
	now := time.Now()
	hour := time.Hour

	withdrawSpecs := []struct {
		PlayerKey    string
		AmountCents  int64
		Method       model.WithdrawMethod
		AccountInfo  string
		Status       model.WithdrawStatus
		RejectReason string
		CreatedAt    time.Time
	}{
		// 待审核
		{PlayerKey: "playerA", AmountCents: 50000, Method: model.WithdrawMethodAlipay, AccountInfo: "138****0001", Status: model.WithdrawStatusPending, CreatedAt: now.Add(-2 * hour)},
		{PlayerKey: "playerB", AmountCents: 80000, Method: model.WithdrawMethodWeChat, AccountInfo: "wx_user_001", Status: model.WithdrawStatusPending, CreatedAt: now.Add(-1 * hour)},
		// 已批准
		{PlayerKey: "playerC", AmountCents: 120000, Method: model.WithdrawMethodBank, AccountInfo: "6222****1234", Status: model.WithdrawStatusApproved, CreatedAt: now.Add(-24 * hour)},
		{PlayerKey: "playerD", AmountCents: 65000, Method: model.WithdrawMethodAlipay, AccountInfo: "138****0009", Status: model.WithdrawStatusApproved, CreatedAt: now.Add(-12 * hour)},
		// 已完成
		{PlayerKey: "playerE", AmountCents: 200000, Method: model.WithdrawMethodBank, AccountInfo: "6228****5678", Status: model.WithdrawStatusCompleted, CreatedAt: now.Add(-72 * hour)},
		{PlayerKey: "playerF", AmountCents: 45000, Method: model.WithdrawMethodWeChat, AccountInfo: "wx_user_002", Status: model.WithdrawStatusCompleted, CreatedAt: now.Add(-48 * hour)},
		{PlayerKey: "playerA", AmountCents: 30000, Method: model.WithdrawMethodAlipay, AccountInfo: "138****0001", Status: model.WithdrawStatusCompleted, CreatedAt: now.Add(-96 * hour)},
		// 已拒绝
		{PlayerKey: "playerB", AmountCents: 500000, Method: model.WithdrawMethodBank, AccountInfo: "6222****9999", Status: model.WithdrawStatusRejected, RejectReason: "提现金额超过单日限额", CreatedAt: now.Add(-36 * hour)},
	}

	for _, spec := range withdrawSpecs {
		player, ok := players[spec.PlayerKey]
		if !ok {
			continue
		}

		// 检查是否已存在
		var existing model.Withdraw
		if err := tx.Where("player_id = ? AND amount_cents = ? AND created_at = ?", player.ID, spec.AmountCents, spec.CreatedAt).First(&existing).Error; err == nil {
			continue
		}

		withdraw := &model.Withdraw{
			PlayerID:    player.ID,
			UserID:      player.UserID,
			AmountCents: spec.AmountCents,
			Method:      spec.Method,
			AccountInfo: spec.AccountInfo,
			Status:      spec.Status,
			CreatedAt:   spec.CreatedAt,
		}

		if spec.RejectReason != "" {
			withdraw.RejectReason = spec.RejectReason
			processedAt := spec.CreatedAt.Add(2 * hour)
			withdraw.ProcessedAt = &processedAt
		}

		if spec.Status == model.WithdrawStatusApproved {
			processedAt := spec.CreatedAt.Add(4 * hour)
			withdraw.ProcessedAt = &processedAt
		}

		if spec.Status == model.WithdrawStatusCompleted {
			processedAt := spec.CreatedAt.Add(4 * hour)
			completedAt := spec.CreatedAt.Add(24 * hour)
			withdraw.ProcessedAt = &processedAt
			withdraw.CompletedAt = &completedAt
			withdraw.ActualAmountCents = spec.AmountCents
		}

		if err := tx.Create(withdraw).Error; err != nil {
			return fmt.Errorf("failed to create withdraw: %w", err)
		}
	}

	log.Println("withdraw seed data created successfully")
	return nil
}

// seedReviewReportData 创建评价举报种子数据
func seedReviewReportData(tx *gorm.DB, users map[string]*model.User) error {
	// 先获取一些评价
	var reviews []model.Review
	if err := tx.Limit(5).Find(&reviews).Error; err != nil {
		return nil // 没有评价则跳过
	}

	if len(reviews) == 0 {
		return nil
	}

	now := time.Now()
	reportSpecs := []struct {
		ReviewIdx   int
		ReporterKey string
		Reason      string
		Status      model.ReviewReportStatus
	}{
		{ReviewIdx: 0, ReporterKey: "customerB", Reason: "评价内容包含不实信息", Status: model.ReviewReportStatusPending},
		{ReviewIdx: 1, ReporterKey: "customerC", Reason: "恶意差评，与实际服务不符", Status: model.ReviewReportStatusPending},
		{ReviewIdx: 2, ReporterKey: "customerD", Reason: "评价包含广告内容", Status: model.ReviewReportStatusApproved},
	}

	for _, spec := range reportSpecs {
		if spec.ReviewIdx >= len(reviews) {
			continue
		}

		reporter, ok := users[spec.ReporterKey]
		if !ok {
			continue
		}

		review := reviews[spec.ReviewIdx]

		// 检查是否已存在
		var existing model.ReviewReport
		if err := tx.Where("review_id = ? AND reporter_id = ?", review.ID, reporter.ID).First(&existing).Error; err == nil {
			continue
		}

		report := &model.ReviewReport{
			ReviewID:   review.ID,
			ReporterID: reporter.ID,
			Reason:     spec.Reason,
			Status:     spec.Status,
		}

		if spec.Status != model.ReviewReportStatusPending {
			handledAt := now.Add(-12 * time.Hour)
			report.HandledAt = &handledAt
			report.HandlingNote = "已处理"
		}

		if err := tx.Create(report).Error; err != nil {
			// 忽略错误，可能是表不存在
			log.Printf("skip creating review report: %v", err)
			continue
		}
	}

	log.Println("review report seed data created successfully")
	return nil
}

// seedSensitiveWords 创建敏感词种子数据
func seedSensitiveWords(tx *gorm.DB) error {
	words := []struct {
		Word     string
		Category model.SensitiveWordCategory
	}{
		// 广告类
		{Word: "加微信", Category: model.SensitiveWordCategoryAd},
		{Word: "加QQ", Category: model.SensitiveWordCategoryAd},
		{Word: "私聊", Category: model.SensitiveWordCategoryAd},
		{Word: "代练", Category: model.SensitiveWordCategoryAd},
		// 违规类
		{Word: "骗子", Category: model.SensitiveWordCategoryOther},
		{Word: "垃圾", Category: model.SensitiveWordCategoryOther},
		// 其他
		{Word: "退款", Category: model.SensitiveWordCategoryOther},
		{Word: "投诉", Category: model.SensitiveWordCategoryOther},
	}

	for _, w := range words {
		var existing model.SensitiveWord
		if err := tx.Where("word = ?", w.Word).First(&existing).Error; err == nil {
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			// 表可能不存在，跳过
			log.Printf("skip sensitive word check: %v", err)
			return nil
		}

		sw := &model.SensitiveWord{
			Word:     w.Word,
			Category: w.Category,
		}

		if err := tx.Create(sw).Error; err != nil {
			log.Printf("skip creating sensitive word: %v", err)
			continue
		}
	}

	log.Println("sensitive words seed data created successfully")
	return nil
}

// seedWalletData 创建钱包种子数据
func seedWalletData(tx *gorm.DB, users map[string]*model.User, _ map[string]*model.Player) error {
	// 为陪玩师创建钱包余额
	walletSpecs := []struct {
		UserKey      string
		BalanceCents int64
		FrozenCents  int64
	}{
		// 陪玩师钱包（有收入）
		{UserKey: "proA", BalanceCents: 358000, FrozenCents: 50000}, // ¥3580 可用，¥500 冻结
		{UserKey: "proB", BalanceCents: 256000, FrozenCents: 80000}, // ¥2560 可用，¥800 冻结
		{UserKey: "proC", BalanceCents: 189000, FrozenCents: 0},     // ¥1890 可用
		{UserKey: "proD", BalanceCents: 125000, FrozenCents: 30000}, // ¥1250 可用，¥300 冻结
		{UserKey: "proE", BalanceCents: 98000, FrozenCents: 0},      // ¥980 可用
		{UserKey: "proF", BalanceCents: 67000, FrozenCents: 15000},  // ¥670 可用，¥150 冻结
		// 普通用户钱包（充值余额）
		{UserKey: "customerA", BalanceCents: 50000, FrozenCents: 0},  // ¥500 可用
		{UserKey: "customerB", BalanceCents: 120000, FrozenCents: 0}, // ¥1200 可用（VIP用户）
		{UserKey: "customerC", BalanceCents: 20000, FrozenCents: 0},  // ¥200 可用
		{UserKey: "customerH", BalanceCents: 80000, FrozenCents: 0},  // ¥800 可用（商务人士）
	}

	for _, spec := range walletSpecs {
		user, ok := users[spec.UserKey]
		if !ok {
			continue
		}

		// 检查是否已存在
		var existing model.Wallet
		if err := tx.Where("user_id = ?", user.ID).First(&existing).Error; err == nil {
			// 已存在，更新余额
			if err := tx.Model(&existing).Updates(map[string]interface{}{
				"balance_cents": spec.BalanceCents,
				"frozen_cents":  spec.FrozenCents,
			}).Error; err != nil {
				return err
			}
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		wallet := &model.Wallet{
			UserID:       user.ID,
			BalanceCents: spec.BalanceCents,
			FrozenCents:  spec.FrozenCents,
		}
		if err := tx.Create(wallet).Error; err != nil {
			return err
		}
	}

	log.Println("wallet seed data created successfully")
	return nil
}

// seedCommissionRules 创建佣金规则种子数据
func seedCommissionRules(tx *gorm.DB, games map[string]*model.Game) error {
	// 检查是否已有规则（除了默认规则）
	var ruleCount int64
	if err := tx.Model(&model.CommissionRule{}).Where("type != ?", "default").Count(&ruleCount).Error; err != nil {
		return err
	}
	if ruleCount > 0 {
		log.Println("commission rules already exist, skipping")
		return nil
	}

	ruleSpecs := []struct {
		Name        string
		Description string
		Type        string
		Rate        int
		GameKey     string
		ServiceType string
	}{
		// 游戏特定规则
		{
			Name:        "英雄联盟专属抽成",
			Description: "英雄联盟游戏订单享受较低抽成",
			Type:        "special",
			Rate:        18,
			GameKey:     "lol",
		},
		{
			Name:        "DOTA2专属抽成",
			Description: "DOTA2游戏订单抽成规则",
			Type:        "special",
			Rate:        18,
			GameKey:     "dota2",
		},
		{
			Name:        "FPS游戏抽成",
			Description: "FPS类游戏（CS:GO、Valorant等）抽成规则",
			Type:        "special",
			Rate:        22,
			GameKey:     "csgo",
		},
		// 服务类型规则
		{
			Name:        "陪练服务抽成",
			Description: "专业陪练服务抽成比例",
			Type:        "special",
			Rate:        15,
			ServiceType: "training",
		},
		{
			Name:        "娱乐陪玩抽成",
			Description: "休闲娱乐陪玩服务抽成",
			Type:        "special",
			Rate:        25,
			ServiceType: "entertainment",
		},
	}

	for _, spec := range ruleSpecs {
		rule := &model.CommissionRule{
			Name:        spec.Name,
			Description: spec.Description,
			Type:        model.CommissionRuleType(spec.Type),
			Rate:        spec.Rate,
			IsActive:    true,
		}

		if spec.GameKey != "" {
			if game, ok := games[spec.GameKey]; ok {
				rule.GameID = &game.ID
			}
		}

		if spec.ServiceType != "" {
			rule.ServiceType = &spec.ServiceType
		}

		if err := tx.Create(rule).Error; err != nil {
			return err
		}
	}

	log.Println("commission rules seed data created successfully")
	return nil
}

// seedCommissionRecords 创建佣金记录种子数据
func seedCommissionRecords(tx *gorm.DB, orders map[string]*model.Order, players map[string]*model.Player) error {
	now := time.Now()
	currentMonth := now.Format("2006-01")
	lastMonth := now.AddDate(0, -1, 0).Format("2006-01")

	// 检查是否已有记录
	var recordCount int64
	if err := tx.Model(&model.CommissionRecord{}).Count(&recordCount).Error; err != nil {
		return err
	}
	if recordCount > 0 {
		log.Println("commission records already exist, skipping")
		return nil
	}

	// 为已完成的订单创建佣金记录
	completedOrders := []string{"orderCompleted1", "orderCompleted2"}
	for _, orderKey := range completedOrders {
		order, ok := orders[orderKey]
		if !ok || order.PlayerID == nil {
			continue
		}

		commissionRate := 20 // 默认20%
		commissionCents := order.TotalPriceCents * int64(commissionRate) / 100
		playerIncome := order.TotalPriceCents - commissionCents

		record := &model.CommissionRecord{
			OrderID:           order.ID,
			PlayerID:          *order.PlayerID,
			TotalAmountCents:  order.TotalPriceCents,
			CommissionRate:    commissionRate,
			CommissionCents:   commissionCents,
			PlayerIncomeCents: playerIncome,
			SettlementStatus:  "settled",
			SettlementMonth:   lastMonth,
			SettledAt:         ptrTime(now.AddDate(0, 0, -5)),
		}

		if err := tx.Create(record).Error; err != nil {
			return err
		}
	}

	// 为进行中的订单创建待结算记录
	pendingOrders := []string{"orderInProgress1", "orderInProgress2", "orderConfirmed1", "orderConfirmed2"}
	for _, orderKey := range pendingOrders {
		order, ok := orders[orderKey]
		if !ok || order.PlayerID == nil {
			continue
		}

		commissionRate := 20
		commissionCents := order.TotalPriceCents * int64(commissionRate) / 100
		playerIncome := order.TotalPriceCents - commissionCents

		record := &model.CommissionRecord{
			OrderID:           order.ID,
			PlayerID:          *order.PlayerID,
			TotalAmountCents:  order.TotalPriceCents,
			CommissionRate:    commissionRate,
			CommissionCents:   commissionCents,
			PlayerIncomeCents: playerIncome,
			SettlementStatus:  "pending",
			SettlementMonth:   currentMonth,
		}

		if err := tx.Create(record).Error; err != nil {
			return err
		}
	}

	// 创建月度结算记录
	settlementSpecs := []struct {
		PlayerKey            string
		Month                string
		OrderCount           int64
		TotalAmountCents     int64
		TotalCommissionCents int64
		TotalIncomeCents     int64
		Status               string
		IncomeRank           int
	}{
		{PlayerKey: "playerA", Month: lastMonth, OrderCount: 45, TotalAmountCents: 450000, TotalCommissionCents: 90000, TotalIncomeCents: 360000, Status: "paid", IncomeRank: 1},
		{PlayerKey: "playerB", Month: lastMonth, OrderCount: 38, TotalAmountCents: 380000, TotalCommissionCents: 76000, TotalIncomeCents: 304000, Status: "paid", IncomeRank: 2},
		{PlayerKey: "playerC", Month: lastMonth, OrderCount: 32, TotalAmountCents: 320000, TotalCommissionCents: 64000, TotalIncomeCents: 256000, Status: "paid", IncomeRank: 3},
		{PlayerKey: "playerD", Month: lastMonth, OrderCount: 28, TotalAmountCents: 280000, TotalCommissionCents: 56000, TotalIncomeCents: 224000, Status: "confirmed", IncomeRank: 4},
		{PlayerKey: "playerE", Month: lastMonth, OrderCount: 22, TotalAmountCents: 220000, TotalCommissionCents: 44000, TotalIncomeCents: 176000, Status: "confirmed", IncomeRank: 5},
		{PlayerKey: "playerF", Month: lastMonth, OrderCount: 18, TotalAmountCents: 180000, TotalCommissionCents: 36000, TotalIncomeCents: 144000, Status: "pending", IncomeRank: 6},
		// 当月数据
		{PlayerKey: "playerA", Month: currentMonth, OrderCount: 12, TotalAmountCents: 120000, TotalCommissionCents: 24000, TotalIncomeCents: 96000, Status: "pending", IncomeRank: 1},
		{PlayerKey: "playerB", Month: currentMonth, OrderCount: 10, TotalAmountCents: 100000, TotalCommissionCents: 20000, TotalIncomeCents: 80000, Status: "pending", IncomeRank: 2},
	}

	for _, spec := range settlementSpecs {
		player, ok := players[spec.PlayerKey]
		if !ok {
			continue
		}

		// 检查是否已存在
		var existing model.MonthlySettlement
		if err := tx.Where("player_id = ? AND settlement_month = ?", player.ID, spec.Month).First(&existing).Error; err == nil {
			continue
		}

		settlement := &model.MonthlySettlement{
			PlayerID:             player.ID,
			SettlementMonth:      spec.Month,
			TotalOrderCount:      spec.OrderCount,
			TotalAmountCents:     spec.TotalAmountCents,
			TotalCommissionCents: spec.TotalCommissionCents,
			TotalIncomeCents:     spec.TotalIncomeCents,
			BonusCents:           0,
			FinalIncomeCents:     spec.TotalIncomeCents,
			Status:               model.MonthlySettlementStatus(spec.Status),
			IncomeRank:           &spec.IncomeRank,
		}

		if spec.Status == "paid" {
			settlement.SettledAt = ptrTime(now.AddDate(0, 0, -3))
		}

		if err := tx.Create(settlement).Error; err != nil {
			return err
		}
	}

	log.Println("commission records seed data created successfully")
	return nil
}

// seedCollectionEntities 创建收款主体和分流规则种子数据
func seedCollectionEntities(tx *gorm.DB) error {
	// 检查是否已有收款主体
	var entityCount int64
	if err := tx.Model(&model.CollectionEntity{}).Count(&entityCount).Error; err != nil {
		return err
	}
	if entityCount > 0 {
		log.Println("collection entities already exist, skipping")
		return nil
	}

	// 创建收款主体
	entities := []model.CollectionEntity{
		{
			Name:                 "游戏联盟科技有限公司",
			CreditCode:           "91110108MA01ABCD1X",
			TaxRegistrationNo:    "110108MA01ABCD1X",
			Status:               model.EntityStatusActive,
			IsDefault:            true,
			TotalCollectionCents: 1580000,
			TransactionCount:     156,
			CreatedBy:            1,
		},
		{
			Name:                 "星耀互娱网络科技公司",
			CreditCode:           "91310115MA1KEFGH2Y",
			TaxRegistrationNo:    "310115MA1KEFGH2Y",
			Status:               model.EntityStatusActive,
			IsDefault:            false,
			TotalCollectionCents: 890000,
			TransactionCount:     89,
			CreatedBy:            1,
		},
		{
			Name:                 "电竞梦想文化传媒公司",
			CreditCode:           "91440300MA5DIJKL3Z",
			TaxRegistrationNo:    "440300MA5DIJKL3Z",
			Status:               model.EntityStatusInactive,
			IsDefault:            false,
			TotalCollectionCents: 320000,
			TransactionCount:     32,
			CreatedBy:            1,
		},
	}

	for i := range entities {
		if err := tx.Create(&entities[i]).Error; err != nil {
			return fmt.Errorf("failed to create collection entity: %w", err)
		}
	}

	// 创建分流规则
	rules := []struct {
		Name           string
		Priority       int
		TargetEntityID uint64
		Conditions     []model.RoutingCondition
		Description    string
	}{
		{
			Name:           "大额订单分流",
			Priority:       1,
			TargetEntityID: 1,
			Conditions: []model.RoutingCondition{
				{Field: model.ConditionFieldOrderAmount, Operator: model.ConditionOperatorGreaterThan, Value: json.RawMessage(`500`)},
			},
			Description: "订单金额超过500元走主收款主体",
		},
		{
			Name:           "LOL游戏分流",
			Priority:       2,
			TargetEntityID: 2,
			Conditions: []model.RoutingCondition{
				{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: json.RawMessage(`"lol"`)},
			},
			Description: "英雄联盟游戏订单走星耀互娱",
		},
		{
			Name:           "华东地区分流",
			Priority:       3,
			TargetEntityID: 2,
			Conditions: []model.RoutingCondition{
				{Field: model.ConditionFieldRegion, Operator: model.ConditionOperatorIn, Value: json.RawMessage(`["上海","江苏","浙江"]`)},
			},
			Description: "华东地区订单走星耀互娱",
		},
	}

	for _, spec := range rules {
		rule := &model.RoutingRule{
			Name:           spec.Name,
			Priority:       spec.Priority,
			TargetEntityID: spec.TargetEntityID,
			Status:         model.RuleStatusActive,
			Description:    spec.Description,
			CreatedBy:      1,
		}
		if err := rule.SetConditions(spec.Conditions); err != nil {
			return fmt.Errorf("failed to set conditions: %w", err)
		}
		if err := tx.Create(rule).Error; err != nil {
			return fmt.Errorf("failed to create routing rule: %w", err)
		}
	}

	log.Println("collection entities and routing rules seed data created successfully")
	return nil
}

// seedSettlementCompanies 创建结算公司种子数据并分配陪玩师
// players 参数用于创建真实的陪玩师-结算公司分配关系
func seedSettlementCompanies(tx *gorm.DB, players map[string]*model.Player) error {
	// 结算公司定义
	companySpecs := []struct {
		Name             string
		CreditCode       string
		BankName         string
		BankAccount      string
		ContactName      string
		ContactPhone     string
		Status           model.CompanyStatus
		TotalPayoutCents int64
	}{
		{
			Name:             "游戏联盟结算中心",
			CreditCode:       "91110108MA01WXYZ1A",
			BankName:         "中国工商银行北京分行",
			BankAccount:      "6222021234567890123",
			ContactName:      "张经理",
			ContactPhone:     "13912345678",
			Status:           model.CompanyStatusActive,
			TotalPayoutCents: 2580000,
		},
		{
			Name:             "星耀支付结算公司",
			CreditCode:       "91310115MA1KMNOP2B",
			BankName:         "招商银行上海分行",
			BankAccount:      "6225881234567890456",
			ContactName:      "李总监",
			ContactPhone:     "13898765432",
			Status:           model.CompanyStatusActive,
			TotalPayoutCents: 1890000,
		},
		{
			Name:             "电竞梦想财务公司",
			CreditCode:       "91440300MA5DQRST3C",
			BankName:         "中国建设银行深圳分行",
			BankAccount:      "6227001234567890789",
			ContactName:      "王会计",
			ContactPhone:     "13765432109",
			Status:           model.CompanyStatusInactive,
			TotalPayoutCents: 560000,
		},
	}

	// 创建或获取结算公司
	companies := make([]*model.SettlementCompany, len(companySpecs))
	for i, spec := range companySpecs {
		var company model.SettlementCompany
		// 通过 CreditCode 查找是否已存在
		err := tx.Where("credit_code = ?", spec.CreditCode).First(&company).Error
		if err == nil {
			// 已存在，使用现有记录
			companies[i] = &company
			log.Printf("settlement company %s already exists (ID:%d)", spec.Name, company.ID)
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			// 不存在，创建新记录
			company = model.SettlementCompany{
				Name:             spec.Name,
				CreditCode:       spec.CreditCode,
				BankName:         spec.BankName,
				BankAccount:      spec.BankAccount,
				ContactName:      spec.ContactName,
				ContactPhone:     spec.ContactPhone,
				Status:           spec.Status,
				TotalPayoutCents: spec.TotalPayoutCents,
				PlayerCount:      0, // 将根据实际分配更新
				CreatedBy:        1,
			}
			if err := tx.Create(&company).Error; err != nil {
				return fmt.Errorf("failed to create settlement company %s: %w", spec.Name, err)
			}
			companies[i] = &company
			log.Printf("created settlement company %s (ID:%d)", spec.Name, company.ID)
		} else {
			return fmt.Errorf("failed to query settlement company: %w", err)
		}
	}

	log.Println("settlement companies seed data ensured successfully")

	// 如果没有陪玩师数据，跳过分配
	if len(players) == 0 {
		log.Println("no players provided, skipping player-company assignments")
		return nil
	}

	// 检查是否已有分配记录
	var assignmentCount int64
	if err := tx.Model(&model.PlayerCompanyAssignment{}).Count(&assignmentCount).Error; err != nil {
		return err
	}
	if assignmentCount > 0 {
		log.Println("player company assignments already exist, skipping")
		return nil
	}

	// 定义陪玩师到结算公司的分配关系
	// 只分配已认证的陪玩师（VerificationVerified）
	// 公司0（游戏联盟结算中心）: playerA, playerB, playerC
	// 公司1（星耀支付结算公司）: playerD, playerE, playerF
	// 公司2（电竞梦想财务公司）: 不分配（已禁用）
	assignmentSpecs := []struct {
		PlayerKey  string
		CompanyIdx int // 0-based index in companies slice
	}{
		{"playerA", 0}, // 峡谷守护者 -> 游戏联盟结算中心
		{"playerB", 0}, // 王牌射手 -> 游戏联盟结算中心
		{"playerC", 0}, // 枪神降临 -> 游戏联盟结算中心
		{"playerD", 1}, // 异世界旅者 -> 星耀支付结算公司
		{"playerE", 1}, // 运动健将 -> 星耀支付结算公司
		{"playerF", 1}, // 欢乐使者 -> 星耀支付结算公司
		// playerG 和 playerH 不分配（playerG 是重复用户，playerH 是待审核状态）
	}

	now := time.Now()
	companyPlayerCounts := make(map[uint64]int) // 统计每个公司分配的陪玩师数量

	for _, spec := range assignmentSpecs {
		player, ok := players[spec.PlayerKey]
		if !ok {
			log.Printf("player %s not found, skipping assignment", spec.PlayerKey)
			continue
		}

		// 只分配已认证的陪玩师
		if player.VerificationStatus != model.VerificationVerified {
			log.Printf("player %s not verified, skipping assignment", spec.PlayerKey)
			continue
		}

		company := companies[spec.CompanyIdx]

		// 创建分配记录
		assignment := &model.PlayerCompanyAssignment{
			PlayerID:            player.ID,
			SettlementCompanyID: company.ID,
			EffectiveDate:       now.AddDate(0, -1, 0), // 一个月前生效
			Reason:              "初始分配 - 系统种子数据",
			AssignedBy:          1, // 管理员ID
			IsCurrent:           true,
		}

		if err := tx.Create(assignment).Error; err != nil {
			return fmt.Errorf("failed to create player company assignment for player %d: %w", player.ID, err)
		}

		companyPlayerCounts[company.ID]++
		log.Printf("assigned player %s (ID:%d) to company %s (ID:%d)", spec.PlayerKey, player.ID, company.Name, company.ID)
	}

	// 更新结算公司的 PlayerCount 字段
	for companyID, count := range companyPlayerCounts {
		if err := tx.Model(&model.SettlementCompany{}).Where("id = ?", companyID).Update("player_count", count).Error; err != nil {
			return fmt.Errorf("failed to update player count for company %d: %w", companyID, err)
		}
		log.Printf("updated company %d player_count to %d", companyID, count)
	}

	log.Printf("player company assignments created successfully, total: %d assignments", len(companyPlayerCounts))
	return nil
}

// seedRankingCommissionConfigs 创建排行榜抽成配置种子数据
func seedRankingCommissionConfigs(tx *gorm.DB) error {
	// 检查是否已有配置
	var configCount int64
	if err := tx.Model(&model.RankingCommissionConfig{}).Count(&configCount).Error; err != nil {
		return err
	}
	if configCount > 0 {
		log.Println("ranking commission configs already exist, skipping")
		return nil
	}

	now := time.Now()
	currentMonth := now.Format("2006-01")
	nextMonth := now.AddDate(0, 1, 0).Format("2006-01")

	configs := []struct {
		Name        string
		RankingType model.RankingType
		Period      string
		Month       string
		Rules       []model.RankingCommissionRule
		Description string
		IsActive    bool
	}{
		{
			Name:        "月度收入排行抽成",
			RankingType: model.RankingTypeIncome,
			Period:      "monthly",
			Month:       currentMonth,
			Rules: []model.RankingCommissionRule{
				{RankStart: 1, RankEnd: 3, CommissionRate: 10},
				{RankStart: 4, RankEnd: 10, CommissionRate: 12},
				{RankStart: 11, RankEnd: 50, CommissionRate: 15},
				{RankStart: 51, RankEnd: 100, CommissionRate: 18},
			},
			Description: "本月收入排行榜抽成规则：TOP3享10%低抽成",
			IsActive:    true,
		},
		{
			Name:        "月度订单量排行抽成",
			RankingType: model.RankingTypeOrderCount,
			Period:      "monthly",
			Month:       currentMonth,
			Rules: []model.RankingCommissionRule{
				{RankStart: 1, RankEnd: 5, CommissionRate: 8},
				{RankStart: 6, RankEnd: 20, CommissionRate: 12},
				{RankStart: 21, RankEnd: 100, CommissionRate: 16},
			},
			Description: "本月订单量排行榜抽成规则：TOP5享8%超低抽成",
			IsActive:    true,
		},
		{
			Name:        "下月收入排行抽成（预设）",
			RankingType: model.RankingTypeIncome,
			Period:      "monthly",
			Month:       nextMonth,
			Rules: []model.RankingCommissionRule{
				{RankStart: 1, RankEnd: 3, CommissionRate: 10},
				{RankStart: 4, RankEnd: 10, CommissionRate: 12},
				{RankStart: 11, RankEnd: 50, CommissionRate: 15},
				{RankStart: 51, RankEnd: 100, CommissionRate: 18},
			},
			Description: "下月收入排行榜抽成规则（预设）",
			IsActive:    false,
		},
	}

	for _, spec := range configs {
		rulesJSON, err := json.Marshal(spec.Rules)
		if err != nil {
			return fmt.Errorf("failed to marshal rules: %w", err)
		}

		config := &model.RankingCommissionConfig{
			Name:        spec.Name,
			RankingType: spec.RankingType,
			Period:      spec.Period,
			Month:       spec.Month,
			RulesJSON:   string(rulesJSON),
			Description: spec.Description,
			IsActive:    spec.IsActive,
		}

		if err := tx.Create(config).Error; err != nil {
			return fmt.Errorf("failed to create ranking commission config: %w", err)
		}
	}

	log.Println("ranking commission configs seed data created successfully")
	return nil
}

// seedOrderDisputes 创建订单纠纷种子数据
func seedOrderDisputes(tx *gorm.DB, orders map[string]*model.Order, users map[string]*model.User) error {
	// 检查是否已有纠纷
	var disputeCount int64
	if err := tx.Model(&model.OrderDispute{}).Count(&disputeCount).Error; err != nil {
		return err
	}
	if disputeCount > 0 {
		log.Println("order disputes already exist, skipping")
		return nil
	}

	now := time.Now()
	hour := time.Hour

	disputeSpecs := []struct {
		OrderKey         string
		UserKey          string
		Status           model.DisputeStatus
		Reason           string
		Description      string
		Resolution       model.DisputeResolution
		ResolutionAmount int64
		ResolutionNotes  string
		AssignedToKey    string
		CreatedOffset    time.Duration
	}{
		// 待处理纠纷
		{
			OrderKey:      "orderInProgress1",
			UserKey:       "customerA",
			Status:        model.DisputeStatusPending,
			Reason:        "陪玩师迟到",
			Description:   "约定时间是晚上8点，但陪玩师9点才上线，严重影响游戏体验",
			Resolution:    model.ResolutionPending,
			CreatedOffset: -2 * hour,
		},
		// 已指派纠纷
		{
			OrderKey:      "orderInProgress2",
			UserKey:       "customerB",
			Status:        model.DisputeStatusAssigned,
			Reason:        "服务态度差",
			Description:   "陪玩师在游戏中频繁挂机，态度敷衍",
			Resolution:    model.ResolutionPending,
			AssignedToKey: "adminA",
			CreatedOffset: -4 * hour,
		},
		// 调解中纠纷
		{
			OrderKey:      "orderConfirmed1",
			UserKey:       "customerC",
			Status:        model.DisputeStatusMediating,
			Reason:        "技术水平不符",
			Description:   "陪玩师声称是王者段位，实际游戏表现只有黄金水平",
			Resolution:    model.ResolutionPending,
			AssignedToKey: "adminA",
			CreatedOffset: -8 * hour,
		},
		// 已解决纠纷（全额退款）
		{
			OrderKey:         "orderCompleted1",
			UserKey:          "customerD",
			Status:           model.DisputeStatusResolved,
			Reason:           "陪玩师中途离开",
			Description:      "游戏进行到一半，陪玩师突然下线不再回来",
			Resolution:       model.ResolutionRefund,
			ResolutionAmount: 9900,
			ResolutionNotes:  "经核实，陪玩师确实中途离开，全额退款",
			AssignedToKey:    "adminA",
			CreatedOffset:    -24 * hour,
		},
		// 已解决纠纷（部分退款）
		{
			OrderKey:         "orderCompleted2",
			UserKey:          "customerE",
			Status:           model.DisputeStatusResolved,
			Reason:           "服务时长不足",
			Description:      "购买了2小时服务，实际只陪玩了1.5小时",
			Resolution:       model.ResolutionPartial,
			ResolutionAmount: 5000,
			ResolutionNotes:  "按实际服务时长计算，退还0.5小时费用",
			AssignedToKey:    "adminA",
			CreatedOffset:    -48 * hour,
		},
		// 已驳回纠纷
		{
			OrderKey:        "orderCanceled1",
			UserKey:         "customerF",
			Status:          model.DisputeStatusRejected,
			Reason:          "要求退款",
			Description:     "游戏输了，要求全额退款",
			Resolution:      model.ResolutionReject,
			ResolutionNotes: "游戏胜负不在服务保障范围内，驳回退款请求",
			AssignedToKey:   "adminA",
			CreatedOffset:   -72 * hour,
		},
	}

	for _, spec := range disputeSpecs {
		order, ok := orders[spec.OrderKey]
		if !ok {
			continue
		}
		user, ok := users[spec.UserKey]
		if !ok {
			continue
		}

		createdAt := now.Add(spec.CreatedOffset)
		slaDeadline := createdAt.Add(30 * time.Minute)

		dispute := &model.OrderDispute{
			OrderID:       order.ID,
			InitiatorID:   user.ID,
			InitiatorType: model.DisputeInitiatorUser,
			Type:          model.DisputeTypeServiceQuality,
			Status:        spec.Status,
			Reason:        spec.Reason,
			EvidenceText:  spec.Description,
			Resolution:    spec.Resolution,
			ResolveRemark: spec.ResolutionNotes,
			SLADeadline:   &slaDeadline,
			EvidenceURLs:  model.EvidenceURLArray{"https://example.com/evidence1.jpg", "https://example.com/evidence2.jpg"},
		}
		dispute.CreatedAt = createdAt

		// 设置指派信息
		if spec.AssignedToKey != "" {
			if assignedTo, ok := users[spec.AssignedToKey]; ok {
				dispute.AssignedServiceID = &assignedTo.ID
			}
		}

		// 设置解决信息
		if spec.Status == model.DisputeStatusResolved || spec.Status == model.DisputeStatusRejected {
			if admin, ok := users["adminA"]; ok {
				dispute.ResolvedBy = &admin.ID
				resolvedAt := createdAt.Add(2 * hour)
				dispute.ResolvedAt = &resolvedAt
			}
		}

		// 检查SLA是否超时
		if now.After(slaDeadline) && spec.Status == model.DisputeStatusPending {
			dispute.SLABreached = true
			dispute.SLABreachedAt = &slaDeadline
		}

		if err := tx.Create(dispute).Error; err != nil {
			return fmt.Errorf("failed to create order dispute: %w", err)
		}
	}

	log.Println("order disputes seed data created successfully")
	return nil
}
