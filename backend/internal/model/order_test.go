package model_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gamelink/internal/model"
)

func TestOrderModel(t *testing.T) {
	now := time.Now()
	orderNo := "ORDER123456"
	userID := uint64(1)
	itemID := uint64(10)
	playerID := uint64(20)
	gameID := uint64(30)
	quantity := 2
	unitPriceCents := int64(5000) // 50元
	totalPriceCents := int64(10000) // 100元
	commissionCents := int64(2000) // 20元
	playerIncomeCents := int64(8000) // 80元

	order := &model.Order{
		Base: model.Base{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		OrderNo:           orderNo,
		UserID:            userID,
		ItemID:            itemID,
		PlayerID:          &playerID,
		RecipientPlayerID: nil,
		Quantity:          quantity,
		UnitPriceCents:    unitPriceCents,
		TotalPriceCents:   totalPriceCents,
		CommissionCents:   commissionCents,
		PlayerIncomeCents: playerIncomeCents,
		Currency:          model.CurrencyCNY,
		Status:            model.OrderStatusPending,
		Title:             "护航服务订单",
		Description:       "这是一个护航服务订单",
		GameID:            &gameID,
		ScheduledStart:    &now,
		ScheduledEnd:      &now,
		StartedAt:         &now,
		CompletedAt:       &now,
		GiftMessage:       "",
		IsAnonymous:       false,
		DeliveredAt:       nil,
		CancelReason:      "",
		RefundAmountCents: 0,
		RefundReason:      "",
		RefundedAt:        nil,
		OrderConfig:       "{}",
		UserNotes:         "用户备注",
		HasDispute:        false,
		DisputeID:         nil,
	}

	assert.Equal(t, orderNo, order.OrderNo)
	assert.Equal(t, userID, order.UserID)
	assert.Equal(t, itemID, order.ItemID)
	assert.Equal(t, &playerID, order.PlayerID)
	assert.Nil(t, order.RecipientPlayerID)
	assert.Equal(t, quantity, order.Quantity)
	assert.Equal(t, unitPriceCents, order.UnitPriceCents)
	assert.Equal(t, totalPriceCents, order.TotalPriceCents)
	assert.Equal(t, commissionCents, order.CommissionCents)
	assert.Equal(t, playerIncomeCents, order.PlayerIncomeCents)
	assert.Equal(t, model.CurrencyCNY, order.Currency)
	assert.Equal(t, model.OrderStatusPending, order.Status)
	assert.Equal(t, "护航服务订单", order.Title)
	assert.Equal(t, "这是一个护航服务订单", order.Description)
	assert.Equal(t, &gameID, order.GameID)
	assert.Equal(t, &now, order.ScheduledStart)
	assert.Equal(t, &now, order.ScheduledEnd)
	assert.Equal(t, &now, order.StartedAt)
	assert.Equal(t, &now, order.CompletedAt)
	assert.Equal(t, "", order.GiftMessage)
	assert.False(t, order.IsAnonymous)
	assert.Nil(t, order.DeliveredAt)
	assert.Equal(t, "", order.CancelReason)
	assert.Equal(t, int64(0), order.RefundAmountCents)
	assert.Equal(t, "", order.RefundReason)
	assert.Nil(t, order.RefundedAt)
	assert.Equal(t, "{}", order.OrderConfig)
	assert.Equal(t, "用户备注", order.UserNotes)
	assert.False(t, order.HasDispute)
	assert.Nil(t, order.DisputeID)
}

func TestOrderIsGiftOrder(t *testing.T) {
	// 普通订单
	order1 := &model.Order{
		RecipientPlayerID: nil,
	}
	assert.False(t, order1.IsGiftOrder())

	// 礼物订单
	recipientID := uint64(100)
	order2 := &model.Order{
		RecipientPlayerID: &recipientID,
	}
	assert.True(t, order2.IsGiftOrder())

	// 空的RecipientPlayerID
	recipientIDZero := uint64(0)
	order3 := &model.Order{
		RecipientPlayerID: &recipientIDZero,
	}
	assert.False(t, order3.IsGiftOrder()) // 因为值为0，所以不是礼物订单
}

func TestOrderGetPlayerID(t *testing.T) {
	// 有PlayerID的情况
	playerID := uint64(50)
	order1 := &model.Order{
		PlayerID: &playerID,
	}
	assert.Equal(t, uint64(50), order1.GetPlayerID())

	// PlayerID为nil的情况
	order2 := &model.Order{
		PlayerID: nil,
	}
	assert.Equal(t, uint64(0), order2.GetPlayerID())
}

func TestOrderGetGameID(t *testing.T) {
	// 有GameID的情况
	gameID := uint64(25)
	order1 := &model.Order{
		GameID: &gameID,
	}
	assert.Equal(t, uint64(25), order1.GetGameID())

	// GameID为nil的情况
	order2 := &model.Order{
		GameID: nil,
	}
	assert.Equal(t, uint64(0), order2.GetGameID())
}

func TestOrderGetPriceCents(t *testing.T) {
	order := &model.Order{
		TotalPriceCents: 15000,
	}
	assert.Equal(t, int64(15000), order.GetPriceCents())
}

func TestOrderSetPlayerID(t *testing.T) {
	order := &model.Order{}
	order.SetPlayerID(75)
	assert.NotNil(t, order.PlayerID)
	assert.Equal(t, uint64(75), *order.PlayerID)
}

func TestOrderSetGameID(t *testing.T) {
	order := &model.Order{}
	order.SetGameID(35)
	assert.NotNil(t, order.GameID)
	assert.Equal(t, uint64(35), *order.GameID)
}

func TestOrderJSONSerialization(t *testing.T) {
	now := time.Now()
	playerID := uint64(20)
	gameID := uint64(30)

	order := &model.Order{
		Base: model.Base{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		OrderNo:           "ORDER123456",
		UserID:            1,
		ItemID:            10,
		PlayerID:          &playerID,
		Quantity:          2,
		UnitPriceCents:    5000,
		TotalPriceCents:   10000,
		CommissionCents:   2000,
		PlayerIncomeCents: 8000,
		Currency:          model.CurrencyCNY,
		Status:            model.OrderStatusPending,
		Title:             "护航服务订单",
		Description:       "这是一个护航服务订单",
		GameID:            &gameID,
		ScheduledStart:    &now,
		ScheduledEnd:      &now,
		StartedAt:         &now,
		CompletedAt:       &now,
		UserNotes:         "用户备注",
		HasDispute:        false,
	}

	// 序列化
	data, err := json.Marshal(order)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "ORDER123456")
	assert.Contains(t, string(data), "护航服务订单")

	// 反序列化
	var decoded model.Order
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, order.OrderNo, decoded.OrderNo)
	assert.Equal(t, order.UserID, decoded.UserID)
	assert.Equal(t, order.ItemID, decoded.ItemID)
	assert.Equal(t, *order.PlayerID, *decoded.PlayerID)
	assert.Equal(t, order.TotalPriceCents, decoded.TotalPriceCents)
	assert.Equal(t, order.Status, decoded.Status)
}

func TestOrderConstants(t *testing.T) {
	// 测试订单状态常量
	assert.Equal(t, model.OrderStatus("pending"), model.OrderStatusPending)
	assert.Equal(t, model.OrderStatus("confirmed"), model.OrderStatusConfirmed)
	assert.Equal(t, model.OrderStatus("in_progress"), model.OrderStatusInProgress)
	assert.Equal(t, model.OrderStatus("completed"), model.OrderStatusCompleted)
	assert.Equal(t, model.OrderStatus("canceled"), model.OrderStatusCanceled)
	assert.Equal(t, model.OrderStatus("refunded"), model.OrderStatusRefunded)
}

func TestOrderEdgeCases(t *testing.T) {
	// 测试零值
	order := &model.Order{
		OrderNo:           "",
		UserID:            0,
		ItemID:            0,
		Quantity:          0,
		UnitPriceCents:    0,
		TotalPriceCents:   0,
		CommissionCents:   0,
		PlayerIncomeCents: 0,
		Title:             "",
		Description:       "",
		GiftMessage:       "",
		CancelReason:      "",
		RefundReason:      "",
		OrderConfig:       "",
		UserNotes:         "",
	}

	assert.Equal(t, "", order.OrderNo)
	assert.Equal(t, uint64(0), order.UserID)
	assert.Equal(t, uint64(0), order.ItemID)
	assert.Equal(t, 0, order.Quantity)
	assert.Equal(t, int64(0), order.UnitPriceCents)
	assert.Equal(t, int64(0), order.TotalPriceCents)
	assert.Equal(t, int64(0), order.CommissionCents)
	assert.Equal(t, int64(0), order.PlayerIncomeCents)

	// 测试最大整数值
	order2 := &model.Order{
		UnitPriceCents:    ^int64(0),
		TotalPriceCents:   ^int64(0),
		CommissionCents:   ^int64(0),
		PlayerIncomeCents: ^int64(0),
	}

	assert.Equal(t, ^int64(0), order2.UnitPriceCents)
	assert.Equal(t, ^int64(0), order2.TotalPriceCents)
	assert.Equal(t, ^int64(0), order2.CommissionCents)
	assert.Equal(t, ^int64(0), order2.PlayerIncomeCents)
}

func TestOrderRelations(t *testing.T) {
	order := &model.Order{
		User: model.User{
			Name: "test user",
		},
		Player: &model.Player{
			Nickname: "test player",
		},
	}

	assert.Equal(t, "test user", order.User.Name)
	assert.NotNil(t, order.Player)
	assert.Equal(t, "test player", order.Player.Nickname)
}

func TestOrderGiftOrderWithZeroRecipientID(t *testing.T) {
	recipientID := uint64(0)
	order := &model.Order{
		RecipientPlayerID: &recipientID,
	}

	// 即使RecipientPlayerID为0，因为值不大于0，所以不是礼物订单
	assert.False(t, order.IsGiftOrder())
}

func TestOrderMethodsWithNilPointers(t *testing.T) {
	order := &model.Order{
		PlayerID: nil,
		GameID:   nil,
	}

	assert.Equal(t, uint64(0), order.GetPlayerID())
	assert.Equal(t, uint64(0), order.GetGameID())
	assert.False(t, order.IsGiftOrder())
}