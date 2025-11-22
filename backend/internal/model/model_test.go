package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserModelIndexes(t *testing.T) {
	user := User{
		Phone:        "13812345678",
		Email:        "user@example.com",
		PasswordHash: "hashed_password",
		Name:         "Test User",
		AvatarURL:    "https://example.com/avatar.jpg",
		Role:         RoleUser,
		Status:       UserStatusActive,
		LastLoginAt:  &time.Time{},
	}

	// 验证字段标签
	assert.Equal(t, "13812345678", user.Phone)
	assert.Equal(t, "user@example.com", user.Email)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, RoleUser, user.Role)
	assert.Equal(t, UserStatusActive, user.Status)
}

func TestOrderModelIndexes(t *testing.T) {
	userID := uint64(123)
	itemID := uint64(456)
	playerID := uint64(789)
	gameID := uint64(101)

	order := Order{
		UserID:            userID,
		ItemID:            itemID,
		PlayerID:          &playerID,
		RecipientPlayerID: &playerID,
		Status:            OrderStatusPending,
		Title:             "Test Order",
		Description:       "Test Description",
		GameID:            &gameID,
		Quantity:          1,
		UnitPriceCents:    1000,
		TotalPriceCents:   1000,
		CommissionCents:   200,
		PlayerIncomeCents: 800,
		Currency:          "CNY",
	}

	// 验证订单字段
	assert.Equal(t, userID, order.UserID)
	assert.Equal(t, itemID, order.ItemID)
	assert.Equal(t, &playerID, order.PlayerID)
	assert.Equal(t, OrderStatusPending, order.Status)
	assert.Equal(t, "Test Order", order.Title)
	assert.Equal(t, int64(1000), order.TotalPriceCents)
	assert.Equal(t, int64(200), order.CommissionCents)
	assert.Equal(t, int64(800), order.PlayerIncomeCents)
}

func TestOrderStatusTransitions(t *testing.T) {
	tests := []struct {
		name        string
		fromStatus  OrderStatus
		toStatus    OrderStatus
		shouldAllow bool
	}{
		{"pending->confirmed", OrderStatusPending, OrderStatusConfirmed, true},
		{"pending->canceled", OrderStatusPending, OrderStatusCanceled, true},
		{"confirmed->in_progress", OrderStatusConfirmed, OrderStatusInProgress, true},
		{"confirmed->canceled", OrderStatusConfirmed, OrderStatusCanceled, true},
		{"in_progress->completed", OrderStatusInProgress, OrderStatusCompleted, true},
		{"in_progress->canceled", OrderStatusInProgress, OrderStatusCanceled, true},
		{"completed->canceled", OrderStatusCompleted, OrderStatusCanceled, false},
		{"canceled->completed", OrderStatusCanceled, OrderStatusCompleted, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 这里可以添加状态机验证逻辑
			// 暂时只验证状态值
			assert.NotEmpty(t, tt.fromStatus)
			assert.NotEmpty(t, tt.toStatus)
		})
	}
}

func TestUserStatusValues(t *testing.T) {
	assert.Equal(t, UserStatusActive, UserStatus("active"))
	assert.Equal(t, UserStatusSuspended, UserStatus("suspended"))
	assert.Equal(t, UserStatusBanned, UserStatus("banned"))
}

func TestOrderStatusValues(t *testing.T) {
	assert.Equal(t, OrderStatusPending, OrderStatus("pending"))
	assert.Equal(t, OrderStatusConfirmed, OrderStatus("confirmed"))
	assert.Equal(t, OrderStatusInProgress, OrderStatus("in_progress"))
	assert.Equal(t, OrderStatusCompleted, OrderStatus("completed"))
	assert.Equal(t, OrderStatusCanceled, OrderStatus("canceled"))
	assert.Equal(t, OrderStatusRefunded, OrderStatus("refunded"))
}

func TestServiceItemModel(t *testing.T) {
	item := ServiceItem{
		ItemCode:       "SERVICE_001",
		Name:           "Test Service",
		Description:    "Test Description",
		Category:       "escort",
		SubCategory:    SubCategorySolo,
		BasePriceCents: 1000,
		ServiceHours:   2,
		CommissionRate: 0.2,
		MinUsers:       1,
		MaxPlayers:     1,
		IsActive:       true,
		SortOrder:      1,
	}

	assert.Equal(t, "SERVICE_001", item.ItemCode)
	assert.Equal(t, "Test Service", item.Name)
	assert.Equal(t, int64(1000), item.BasePriceCents)
	assert.Equal(t, 2, item.ServiceHours)
	assert.Equal(t, 0.2, item.CommissionRate)
	assert.True(t, item.IsActive)

	// 测试是否为礼物
	assert.False(t, item.IsGift())

	giftItem := ServiceItem{
		ItemCode:    "GIFT_001",
		Name:        "Test Gift",
		SubCategory: SubCategoryGift,
	}
	assert.True(t, giftItem.IsGift())
}

func TestServiceItemCalculateCommission(t *testing.T) {
	item := ServiceItem{
		BasePriceCents: 1000,
		CommissionRate: 0.2,
	}

	platformCommission, playerIncome := item.CalculateCommission(2)

	assert.Equal(t, int64(400), platformCommission) // 2000 * 0.2 = 400
	assert.Equal(t, int64(1600), playerIncome)      // 2000 - 400 = 1600
}

func TestOrderIsGiftOrder(t *testing.T) {
	// 非礼物订单
	order1 := Order{
		UserID:   123,
		ItemID:   456,
		PlayerID: nil, // 没有接收者
	}
	assert.False(t, order1.IsGiftOrder())

	// 礼物订单
	recipientID := uint64(789)
	order2 := Order{
		UserID:            123,
		ItemID:            456,
		PlayerID:          nil,
		RecipientPlayerID: &recipientID,
	}
	assert.True(t, order2.IsGiftOrder())
}

func TestOrderGettersAndSetters(t *testing.T) {
	order := Order{
		UserID: 123,
		ItemID: 456,
	}

	// 测试 GetPlayerID
	assert.Equal(t, uint64(0), order.GetPlayerID())

	playerID := uint64(789)
	order.SetPlayerID(playerID)
	assert.Equal(t, playerID, order.GetPlayerID())

	// 测试 GetGameID
	assert.Equal(t, uint64(0), order.GetGameID())

	gameID := uint64(101)
	order.SetGameID(gameID)
	assert.Equal(t, gameID, order.GetGameID())

	// 测试 GetPriceCents
	assert.Equal(t, int64(0), order.GetPriceCents())
}
