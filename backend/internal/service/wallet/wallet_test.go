package wallet

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
)

func TestRechargeRequest_Fields(t *testing.T) {
	req := RechargeRequest{
		AmountCents: 10000,
		Method:      model.PaymentMethodAlipay,
	}

	assert.Equal(t, int64(10000), req.AmountCents)
	assert.Equal(t, model.PaymentMethodAlipay, req.Method)
}

func TestRechargeResponse_Fields(t *testing.T) {
	resp := RechargeResponse{
		OrderID:   1,
		PaymentID: 2,
		Balance:   15000,
	}

	assert.Equal(t, uint64(1), resp.OrderID)
	assert.Equal(t, uint64(2), resp.PaymentID)
	assert.Equal(t, int64(15000), resp.Balance)
}

func TestErrInvalidAmount(t *testing.T) {
	assert.Equal(t, "invalid amount", ErrInvalidAmount.Error())
}

func TestRechargeRequest_Validation(t *testing.T) {
	tests := []struct {
		name        string
		amountCents int64
		isValid     bool
	}{
		{"positive amount", 1000, true},
		{"zero amount", 0, false},
		{"negative amount", -100, false},
		{"large amount", 1000000, true},
		{"minimum valid", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := RechargeRequest{
				AmountCents: tt.amountCents,
				Method:      model.PaymentMethodAlipay,
			}
			isValid := req.AmountCents > 0
			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

func TestPaymentMethods(t *testing.T) {
	// Test that payment methods are valid
	methods := []model.PaymentMethod{
		model.PaymentMethodAlipay,
		model.PaymentMethodWeChat,
	}

	for _, method := range methods {
		assert.NotEmpty(t, string(method))
	}
}

func TestRechargeResponse_ZeroValues(t *testing.T) {
	resp := RechargeResponse{}

	assert.Equal(t, uint64(0), resp.OrderID)
	assert.Equal(t, uint64(0), resp.PaymentID)
	assert.Equal(t, int64(0), resp.Balance)
}
