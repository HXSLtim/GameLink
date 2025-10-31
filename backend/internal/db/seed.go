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
	OrderKey  string
	UserKey   string
	PlayerKey string
	Score     model.Rating
	Content   string
}

func applySeeds(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
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
			{Key: "adminA", Email: "admin@gamelink.com", Phone: "13800138005", Name: "系统管理员", Role: model.RoleAdmin, Password: "Admin@123456"},
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

		hour := time.Hour

		orderSpecs := []seedOrderSpec{
			{
				Key:         "orderCompleted1",
				UserKey:     "customerA",
				PlayerKey:   "playerA",
				GameKey:     "lol",
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
			{
				OrderKey:  "orderCompleted1",
				UserKey:   "customerA",
				PlayerKey: "playerA",
				Score:     model.MustRating(5),
				Content:   "很满意的陪玩体验，带我连胜！峡谷守护者技术确实强，打野节奏把控很好。",
			},
			{
				OrderKey:  "orderInProgress1",
				UserKey:   "customerB",
				PlayerKey: "playerA",
				Score:     model.MustRating(4),
				Content:   "战术指导很专业，期待后续完成。DOTA2的复杂度很高，有专业指导确实不一样。",
			},
			{
				OrderKey:  "orderCompleted2",
				UserKey:   "customerF",
				PlayerKey: "playerD",
				Score:     model.MustRating(5),
				Content:   "异世界旅者对MMORPG的理解非常深入，带我了解了魔兽世界的核心玩法，收益很大！",
			},
			{
				OrderKey:  "orderInProgress2",
				UserKey:   "customerG",
				PlayerKey: "playerE",
				Score:     model.MustRating(4),
				Content:   "运动健将的足球水平很高，学到了很多实用的技巧。FIFA游戏体验很好。",
			},
			{
				OrderKey:  "orderCompleted1",
				UserKey:   "customerD",
				PlayerKey: "playerC",
				Score:     model.MustRating(5),
				Content:   "枪神降临不愧是职业选手，枪法精准，教学耐心细致。CS:GO的水平确实提升了很多！",
			},
			{
				OrderKey:  "orderConfirmed2",
				UserKey:   "customerB",
				PlayerKey: "playerG",
				Score:     model.MustRating(4),
				Content:   "DOTA宗师的教学很系统，从基础到进阶都有涉及，受益匪浅。",
			},
			{
				OrderKey:  "orderRefunded1",
				UserKey:   "customerA",
				PlayerKey: "playerB",
				Score:     model.MustRating(3),
				Content:   "虽然因为服务器维护退款了，但之前的服务还不错。期待下次能正常完成。",
			},
		}

		for _, spec := range reviewSpecs {
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
			if err := seedReview(tx, seedReviewParams{
				OrderID:  order.ID,
				UserID:   user.ID,
				PlayerID: player.ID,
				Score:    spec.Score,
				Content:  spec.Content,
			}); err != nil {
				return err
			}
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
}

type seedReviewParams struct {
	OrderID  uint64
	UserID   uint64
	PlayerID uint64
	Score    model.Rating
	Content  string
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
		if err := tx.Where("key = ?", game.Key).First(&existing).Error; err == nil {
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
	lookup := tx.Model(&model.User{})
	if input.Email != "" {
		lookup = lookup.Where("email = ?", input.Email)
	} else {
		lookup = lookup.Where("phone = ?", input.Phone)
	}
	var existing model.User
	if err := lookup.First(&existing).Error; err == nil {
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

func seedOrder(tx *gorm.DB, input seedOrderParams) (*model.Order, error) {
	var existing model.Order
	if err := tx.Where("title = ? AND user_id = ?", input.Title, input.UserID).First(&existing).Error; err == nil {
		return &existing, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	order := &model.Order{
		UserID:         input.UserID,
		GameID:         input.GameID,
		Title:          input.Title,
		Description:    input.Description,
		Status:         input.Status,
		PriceCents:     input.PriceCents,
		Currency:       input.Currency,
		ScheduledStart: input.ScheduledStart,
		ScheduledEnd:   input.ScheduledEnd,
		CancelReason:   strings.TrimSpace(input.CancelReason),
	}
	if input.PlayerID != nil {
		order.PlayerID = *input.PlayerID
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
	}
	if err := tx.Create(payment).Error; err != nil {
		return err
	}
	return nil
}

func seedReview(tx *gorm.DB, input seedReviewParams) error {
	var existing model.Review
	if err := tx.Where("order_id = ?", input.OrderID).First(&existing).Error; err == nil {
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	review := &model.Review{
		OrderID:  input.OrderID,
		UserID:   input.UserID,
		PlayerID: input.PlayerID,
		Score:    input.Score,
		Content:  input.Content,
	}
	if err := tx.Create(review).Error; err != nil {
		return err
	}
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
