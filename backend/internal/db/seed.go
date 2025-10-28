package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"gamelink/internal/model"
)

func applySeeds(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		games, err := seedGames(tx)
		if err != nil {
			return err
		}

		lolGame, ok := games["lol"]
		if !ok {
			return fmt.Errorf("seed game 'lol' not found")
		}

		customer, err := seedUser(tx, seedUserInput{
			Email:    "demo.user@gamelink.com",
			Phone:    "13800138000",
			Name:     "测试用户",
			Role:     model.RoleUser,
			Password: "User@123456",
		})
		if err != nil {
			return err
		}

		proUser, err := seedUser(tx, seedUserInput{
			Email:    "pro.player@gamelink.com",
			Phone:    "13800138001",
			Name:     "职业陪玩",
			Role:     model.RolePlayer,
			Password: "Player@123456",
		})
		if err != nil {
			return err
		}

		proPlayer, err := seedPlayer(tx, proUser.ID, lolGame.ID)
		if err != nil {
			return err
		}

		order, err := seedOrder(tx, &model.Order{
			UserID:      customer.ID,
			PlayerID:    proPlayer.ID,
			GameID:      lolGame.ID,
			Title:       "欢迎体验 GameLink 陪玩",
			Description: "我们为您匹配了经验丰富的高胜率陪玩，尽情享受游戏时光吧！",
			Status:      model.OrderStatusCompleted,
			PriceCents:  19900,
			Currency:    model.CurrencyCNY,
		})
		if err != nil {
			return err
		}

		if err := seedPayment(tx, &model.Payment{
			OrderID:         order.ID,
			UserID:          customer.ID,
			Method:          model.PaymentMethodWeChat,
			AmountCents:     19900,
			Currency:        model.CurrencyCNY,
			Status:          model.PaymentStatusPaid,
			ProviderTradeNo: "WX1234567890",
			ProviderRaw:     json.RawMessage(`{"result":"success"}`),
			PaidAt:          ptrTime(time.Now().Add(-2 * time.Hour)),
		}); err != nil {
			return err
		}

		if err := seedReview(tx, &model.Review{
			OrderID:  order.ID,
			UserID:   customer.ID,
			PlayerID: proPlayer.ID,
			Score:    5,
			Content:  "很满意的陪玩体验，带我连胜！",
		}); err != nil {
			return err
		}

		log.Println("seed data ensured for demo environment")
		return nil
	})
}

type seedUserInput struct {
	Email    string
	Phone    string
	Name     string
	Role     model.Role
	Password string
}

func seedGames(tx *gorm.DB) (map[string]*model.Game, error) {
	seeds := []model.Game{
		{Key: "lol", Name: "英雄联盟", Category: "moba", Description: "召唤师峡谷 5v5 对战"},
		{Key: "dota2", Name: "DOTA 2", Category: "moba", Description: "经典即时战略竞技"},
		{Key: "valorant", Name: "无畏契约", Category: "fps", Description: "英雄战术射击"},
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

func seedPlayer(tx *gorm.DB, userID uint64, mainGameID uint64) (*model.Player, error) {
	var existing model.Player
	if err := tx.Where("user_id = ?", userID).First(&existing).Error; err == nil {
		return &existing, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	player := &model.Player{
		UserID:             userID,
		Nickname:           "峡谷守护者",
		Bio:                "全职陪玩，擅长打野位，帮助玩家快速上分。",
		RatingAverage:      4.9,
		RatingCount:        152,
		HourlyRateCents:    9900,
		MainGameID:         mainGameID,
		VerificationStatus: model.VerificationVerified,
	}
	if err := tx.Create(player).Error; err != nil {
		return nil, err
	}
	return player, nil
}

func seedOrder(tx *gorm.DB, order *model.Order) (*model.Order, error) {
	var existing model.Order
	if err := tx.Where("title = ? AND user_id = ?", order.Title, order.UserID).First(&existing).Error; err == nil {
		return &existing, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	now := time.Now()
	order.ScheduledStart = ptrTime(now.Add(-3 * time.Hour))
	order.ScheduledEnd = ptrTime(now.Add(-2 * time.Hour))
	if err := tx.Create(order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

func seedPayment(tx *gorm.DB, payment *model.Payment) error {
	var existing model.Payment
	if err := tx.Where("order_id = ?", payment.OrderID).First(&existing).Error; err == nil {
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err := tx.Create(payment).Error; err != nil {
		return err
	}
	return nil
}

func seedReview(tx *gorm.DB, review *model.Review) error {
	var existing model.Review
	if err := tx.Where("order_id = ?", review.OrderID).First(&existing).Error; err == nil {
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err := tx.Create(review).Error; err != nil {
		return err
	}
	return nil
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
