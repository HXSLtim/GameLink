// Package testutil provides utilities for testing.
package testutil

import (
	"time"

	"gamelink/internal/model"
)

// NewTestUser creates a test user.
func NewTestUser() *model.User {
	return &model.User{
		Base: model.Base{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: "hashed-password",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
	}
}

// NewTestPlayer creates a test player.
func NewTestPlayer() *model.Player {
	user := NewTestUser()
	return &model.Player{
		Base: model.Base{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserID:             user.ID,
		Nickname:           "Test Player",
		Bio:                "Test bio",
		HourlyRateCents:    1000,
		VerificationStatus: model.VerificationVerified,
	}
}

// NewTestOrder creates a test order.
func NewTestOrder() *model.Order {
	user := NewTestUser()
	player := NewTestPlayer()
	playerID := player.ID
	itemID := uint64(1)

	return &model.Order{
		Base: model.Base{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserID:          user.ID,
		ItemID:          itemID,
		Status:          model.OrderStatusPending,
		TotalPriceCents: 10000,
		UnitPriceCents:  10000,
		Quantity:        1,
		PlayerID:        &playerID,
	}
}

// NewTestPayment creates a test payment.
func NewTestPayment() *model.Payment {
	order := NewTestOrder()

	return &model.Payment{
		Base: model.Base{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		OrderID:     order.ID,
		UserID:      order.UserID,
		AmountCents: order.TotalPriceCents,
		Method:      model.PaymentMethodWeChat,
		Status:      model.PaymentStatusPending,
	}
}
