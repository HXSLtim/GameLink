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

		// 确保有一个默认护航服务项，供订单外键引用
		defaultGame := games["lol"]
		serviceItem, err := ensureServiceItem(tx, "escort-default", "默认护航服务", defaultGame.ID)
		if err != nil {
			return err
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
				ItemID:         serviceItem.ID,
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
				OrderKey:   "orderCompleted1",
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
				OrderKey:  "orderConfirmed1",
				UserKey:   "customerC",
				PlayerKey: "playerF",
				Score:     model.MustRating(3),
				Content:   "整体体验一般，陪玩师技术还可以，但沟通上有些问题，回复不够及时。希望能改进。",
				Status:    model.ReviewStatusApproved,
			},
			// 带图片的好评
			{
				OrderKey:  "orderPending1",
				UserKey:   "customerE",
				PlayerKey: "playerH",
				Score:     model.MustRating(5),
				Content:   "卡牌大师对炉石传说的理解太深了！帮我组了一套超强卡组，连胜上传说！感谢！",
				Status:    model.ReviewStatusPending,
				Images:    []string{"https://gamelink.oss.com/reviews/hs_legend_1.jpg", "https://gamelink.oss.com/reviews/hs_deck_1.jpg", "https://gamelink.oss.com/reviews/hs_win_1.jpg"},
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
func seedMenus(tx *gorm.DB) error {
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

// seedReviewPermissions 创建评价管理权限种子数据
func seedReviewPermissions(tx *gorm.DB) error {
	// 不再检查是否已有权限，改为逐条 upsert，确保新增的权限能被添加

	// 定义评价管理权限
	// 注意：code 字段有唯一索引，所以每个权限必须有唯一的 code
	permissions := []model.Permission{
		// 评价查看权限 (review:view)
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/reviews",
			Code:        "review.list",
			Group:       "/admin/reviews",
			Description: "查看评价列表",
		},
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/reviews/:id",
			Code:        "review.get",
			Group:       "/admin/reviews",
			Description: "查看评价详情",
		},
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/reviews/pending",
			Code:        "review.pending",
			Group:       "/admin/reviews",
			Description: "查看待审核评价",
		},
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/reviews/:id/logs",
			Code:        "review.logs",
			Group:       "/admin/reviews",
			Description: "查看评价操作日志",
		},
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/players/:id/reviews",
			Code:        "review.player",
			Group:       "/admin/reviews",
			Description: "查看陪玩师评价",
		},
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/orders/:id/reviews",
			Code:        "review.order",
			Group:       "/admin/reviews",
			Description: "查看订单评价",
		},

		// 评价审核权限 (review:approve)
		{
			Method:      model.HTTPMethodPUT,
			Path:        "/api/v1/admin/reviews/:id/approve",
			Code:        "review.approve",
			Group:       "/admin/reviews",
			Description: "批准评价",
		},
		{
			Method:      model.HTTPMethodPUT,
			Path:        "/api/v1/admin/reviews/:id/reject",
			Code:        "review.reject",
			Group:       "/admin/reviews",
			Description: "拒绝评价",
		},
		{
			Method:      model.HTTPMethodPUT,
			Path:        "/api/v1/admin/reviews/batch-approve",
			Code:        "review.batch_approve",
			Group:       "/admin/reviews",
			Description: "批量批准评价",
		},
		{
			Method:      model.HTTPMethodPUT,
			Path:        "/api/v1/admin/reviews/batch-reject",
			Code:        "review.batch_reject",
			Group:       "/admin/reviews",
			Description: "批量拒绝评价",
		},

		// 评价删除权限 (review:delete)
		{
			Method:      model.HTTPMethodDELETE,
			Path:        "/api/v1/admin/reviews/:id",
			Code:        "review.delete",
			Group:       "/admin/reviews",
			Description: "删除评价",
		},
		{
			Method:      model.HTTPMethodPUT,
			Path:        "/api/v1/admin/reviews/:id",
			Code:        "review.update",
			Group:       "/admin/reviews",
			Description: "更新评价",
		},
		{
			Method:      model.HTTPMethodPOST,
			Path:        "/api/v1/admin/reviews",
			Code:        "review.create",
			Group:       "/admin/reviews",
			Description: "创建评价",
		},

		// 评价管理权限 (review:manage) - 敏感词管理
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/sensitive-words",
			Code:        "sensitive_word.list",
			Group:       "/admin/reviews",
			Description: "查看敏感词列表",
		},
		{
			Method:      model.HTTPMethodPOST,
			Path:        "/api/v1/admin/sensitive-words",
			Code:        "sensitive_word.create",
			Group:       "/admin/reviews",
			Description: "添加敏感词",
		},
		{
			Method:      model.HTTPMethodPUT,
			Path:        "/api/v1/admin/sensitive-words/:id",
			Code:        "sensitive_word.update",
			Group:       "/admin/reviews",
			Description: "更新敏感词",
		},
		{
			Method:      model.HTTPMethodDELETE,
			Path:        "/api/v1/admin/sensitive-words/:id",
			Code:        "sensitive_word.delete",
			Group:       "/admin/reviews",
			Description: "删除敏感词",
		},
		{
			Method:      model.HTTPMethodPOST,
			Path:        "/api/v1/admin/reviews/detect-sensitive",
			Code:        "review.detect_sensitive",
			Group:       "/admin/reviews",
			Description: "检测敏感词",
		},

		// 评价举报管理权限
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/review-reports",
			Code:        "review_report.list",
			Group:       "/admin/reviews",
			Description: "查看评价举报列表",
		},
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/review-reports/:id",
			Code:        "review_report.get",
			Group:       "/admin/reviews",
			Description: "查看举报详情",
		},
		{
			Method:      model.HTTPMethodPOST,
			Path:        "/api/v1/admin/reviews/:id/reports",
			Code:        "review_report.create",
			Group:       "/admin/reviews",
			Description: "创建评价举报",
		},
		{
			Method:      model.HTTPMethodPUT,
			Path:        "/api/v1/admin/review-reports/:id/handle",
			Code:        "review_report.handle",
			Group:       "/admin/reviews",
			Description: "处理评价举报",
		},

		// 评价回复管理权限
		{
			Method:      model.HTTPMethodPUT,
			Path:        "/api/v1/admin/review-replies/:id",
			Code:        "review_reply.update",
			Group:       "/admin/reviews",
			Description: "更新评价回复",
		},
		{
			Method:      model.HTTPMethodDELETE,
			Path:        "/api/v1/admin/review-replies/:id",
			Code:        "review_reply.delete",
			Group:       "/admin/reviews",
			Description: "删除评价回复",
		},

		// 评价统计权限
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/reviews/stats",
			Code:        "review.stats",
			Group:       "/admin/reviews",
			Description: "查看评价统计",
		},
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/reviews/trend",
			Code:        "review.trend",
			Group:       "/admin/reviews",
			Description: "查看评价趋势",
		},
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/reviews/top-players",
			Code:        "review.top_players",
			Group:       "/admin/reviews",
			Description: "查看陪玩师排行榜",
		},
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/reviews/game-stats",
			Code:        "review.game_stats",
			Group:       "/admin/reviews",
			Description: "查看游戏评价统计",
		},
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/reviews/export",
			Code:        "review.export",
			Group:       "/admin/reviews",
			Description: "导出评价统计",
		},

		// 评价展示设置权限
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/review-settings",
			Code:        "review_settings.get",
			Group:       "/admin/reviews",
			Description: "查看评价展示设置",
		},
		{
			Method:      model.HTTPMethodPUT,
			Path:        "/api/v1/admin/review-settings",
			Code:        "review_settings.update",
			Group:       "/admin/reviews",
			Description: "更新评价展示设置",
		},

		// 操作日志权限
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/operation-logs",
			Code:        "operation_log.list",
			Group:       "/admin/reviews",
			Description: "查看操作日志",
		},
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/operation-logs/export",
			Code:        "operation_log.export",
			Group:       "/admin/reviews",
			Description: "导出操作日志",
		},

		// ========== 内容管理权限 ==========
		// 动态审核权限 (content:view)
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/content/feeds",
			Code:        "content.feed.list",
			Group:       "/admin/content",
			Description: "查看动态列表",
		},
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/content/feeds/:id",
			Code:        "content.feed.get",
			Group:       "/admin/content",
			Description: "查看动态详情",
		},
		// 动态审核权限 (content:moderate)
		{
			Method:      model.HTTPMethodPUT,
			Path:        "/api/v1/admin/content/feeds/:id/approve",
			Code:        "content.feed.approve",
			Group:       "/admin/content",
			Description: "批准动态",
		},
		{
			Method:      model.HTTPMethodPUT,
			Path:        "/api/v1/admin/content/feeds/:id/reject",
			Code:        "content.feed.reject",
			Group:       "/admin/content",
			Description: "拒绝动态",
		},
		{
			Method:      model.HTTPMethodPOST,
			Path:        "/api/v1/admin/content/feeds/batch-approve",
			Code:        "content.feed.batch_approve",
			Group:       "/admin/content",
			Description: "批量批准动态",
		},
		{
			Method:      model.HTTPMethodPOST,
			Path:        "/api/v1/admin/content/feeds/batch-reject",
			Code:        "content.feed.batch_reject",
			Group:       "/admin/content",
			Description: "批量拒绝动态",
		},
		// 动态删除权限 (content:delete)
		{
			Method:      model.HTTPMethodDELETE,
			Path:        "/api/v1/admin/content/feeds/:id",
			Code:        "content.feed.delete",
			Group:       "/admin/content",
			Description: "删除动态",
		},

		// 聊天监控权限
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/content/chat/messages",
			Code:        "content.chat.list",
			Group:       "/admin/content",
			Description: "查看聊天消息",
		},
		{
			Method:      model.HTTPMethodDELETE,
			Path:        "/api/v1/admin/content/chat/messages/:id",
			Code:        "content.chat.delete",
			Group:       "/admin/content",
			Description: "删除聊天消息",
		},
		{
			Method:      model.HTTPMethodPOST,
			Path:        "/api/v1/admin/content/chat/mute",
			Code:        "content.chat.mute",
			Group:       "/admin/content",
			Description: "禁言用户",
		},
		{
			Method:      model.HTTPMethodPOST,
			Path:        "/api/v1/admin/content/chat/unmute",
			Code:        "content.chat.unmute",
			Group:       "/admin/content",
			Description: "解除禁言",
		},

		// 举报管理权限
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/content/reports",
			Code:        "content.report.list",
			Group:       "/admin/content",
			Description: "查看举报列表",
		},
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/content/reports/:id",
			Code:        "content.report.get",
			Group:       "/admin/content",
			Description: "查看举报详情",
		},
		{
			Method:      model.HTTPMethodPOST,
			Path:        "/api/v1/admin/content/reports/:id/process",
			Code:        "content.report.process",
			Group:       "/admin/content",
			Description: "处理举报",
		},

		// 内容统计权限
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/content/stats",
			Code:        "content.stats",
			Group:       "/admin/content",
			Description: "查看内容统计",
		},

		// 内容分类管理权限
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/content/categories",
			Code:        "content.category.list",
			Group:       "/admin/content",
			Description: "查看内容分类",
		},
		{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/content/categories/:id",
			Code:        "content.category.get",
			Group:       "/admin/content",
			Description: "查看分类详情",
		},
		{
			Method:      model.HTTPMethodPOST,
			Path:        "/api/v1/admin/content/categories",
			Code:        "content.category.create",
			Group:       "/admin/content",
			Description: "创建内容分类",
		},
		{
			Method:      model.HTTPMethodPUT,
			Path:        "/api/v1/admin/content/categories/:id",
			Code:        "content.category.update",
			Group:       "/admin/content",
			Description: "更新内容分类",
		},
		{
			Method:      model.HTTPMethodDELETE,
			Path:        "/api/v1/admin/content/categories/:id",
			Code:        "content.category.delete",
			Group:       "/admin/content",
			Description: "删除内容分类",
		},
	}

	// 批量创建权限
	for _, perm := range permissions {
		if err := tx.Create(&perm).Error; err != nil {
			// 如果权限已存在（唯一索引冲突），跳过
			if strings.Contains(err.Error(), "UNIQUE constraint failed") ||
				strings.Contains(err.Error(), "duplicate key") {
				log.Printf("permission already exists: %s %s, skipping\n", perm.Method, perm.Path)
				continue
			}
			return fmt.Errorf("failed to create permission %s %s: %w", perm.Method, perm.Path, err)
		}
	}

	log.Println("review permissions seed data created successfully")
	return nil
}
