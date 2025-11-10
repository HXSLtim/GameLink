package admin

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"gamelink/internal/cache"
	"gamelink/internal/model"
)

// 深度业务逻辑测试

// 跳过需要audit log的复杂测试，专注于可测试的业务逻辑

func TestPasswordValidation(t *testing.T) {
	tests := []struct {
		name     string
		password string
		valid    bool
	}{
		{"有效密码-字母数字", "pass123", true},
		{"有效密码-大小写数字", "Pass123", true},
		{"无效-太短", "p1", false},
		{"无效-只有字母", "password", false},
		{"无效-只有数字", "123456", false},
		{"无效-空字符串", "", false},
		{"有效-长密码", "password123456", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validPassword(tt.password)
			assert.Equal(t, tt.valid, result)
		})
	}
}

func TestHashPassword(t *testing.T) {
	t.Run("成功加密密码", func(t *testing.T) {
		password := "testpass123"
		hashed, err := hashPassword(password)
		
		assert.NoError(t, err)
		assert.NotEmpty(t, hashed)
		assert.NotEqual(t, password, hashed)
		
		// 验证可以解密
		err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
		assert.NoError(t, err)
	})
	
	t.Run("空密码失败", func(t *testing.T) {
		hashed, err := hashPassword("")
		
		assert.Error(t, err)
		assert.Empty(t, hashed)
		assert.Equal(t, ErrValidation, err)
	})
	
	t.Run("空格密码失败", func(t *testing.T) {
		hashed, err := hashPassword("   ")
		
		assert.Error(t, err)
		assert.Empty(t, hashed)
		assert.Equal(t, ErrValidation, err)
	})
}

// TestAssignOrder 跳过，需要audit log支持

func TestConfirmOrder(t *testing.T) {
	ctx := context.Background()
	
	t.Run("成功确认订单", func(t *testing.T) {
		orderRepo := &fakeOrderRepo{
			obj: &model.Order{
				Base:            model.Base{ID: 1},
				Status:          model.OrderStatusPending,
				TotalPriceCents: 10000,
				Currency:        model.CurrencyCNY,
			},
		}
		
		svc := NewAdminService(
			&fakeGameRepo{},
			&fakeUserRepo{},
			&fakePlayerRepo{},
			orderRepo,
			&fakePaymentRepo{},
			&fakeRoleRepo{},
			cache.NewMemory(),
		)
		
		order, err := svc.ConfirmOrder(ctx, 1, "确认订单")
		
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, model.OrderStatusConfirmed, order.Status)
	})
}

func TestStartOrder(t *testing.T) {
	ctx := context.Background()
	
	t.Run("成功开始订单", func(t *testing.T) {
		orderRepo := &fakeOrderRepo{
			obj: &model.Order{
				Base:            model.Base{ID: 1},
				Status:          model.OrderStatusConfirmed,
				TotalPriceCents: 10000,
				Currency:        model.CurrencyCNY,
			},
		}
		
		svc := NewAdminService(
			&fakeGameRepo{},
			&fakeUserRepo{},
			&fakePlayerRepo{},
			orderRepo,
			&fakePaymentRepo{},
			&fakeRoleRepo{},
			cache.NewMemory(),
		)
		
		order, err := svc.StartOrder(ctx, 1, "开始服务")
		
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, model.OrderStatusInProgress, order.Status)
		assert.NotNil(t, order.StartedAt)
	})
}

func TestCreateOrder_Validation(t *testing.T) {
	ctx := context.Background()
	
	tests := []struct {
		name    string
		input   CreateOrderInput
		wantErr bool
	}{
		{
			name: "有效订单",
			input: CreateOrderInput{
				UserID:          1,
				GameID:          1,
				TotalPriceCents: 10000,
				Currency:        model.CurrencyCNY,
				Title:           "测试订单",
			},
			wantErr: false,
		},
		{
			name: "UserID为0",
			input: CreateOrderInput{
				UserID:          0,
				GameID:          1,
				TotalPriceCents: 10000,
				Currency:        model.CurrencyCNY,
			},
			wantErr: true,
		},
		{
			name: "GameID为0",
			input: CreateOrderInput{
				UserID:          1,
				GameID:          0,
				TotalPriceCents: 10000,
				Currency:        model.CurrencyCNY,
			},
			wantErr: true,
		},
		{
			name: "价格为负数",
			input: CreateOrderInput{
				UserID:          1,
				GameID:          1,
				TotalPriceCents: -100,
				Currency:        model.CurrencyCNY,
			},
			wantErr: true,
		},
		{
			name: "无效货币",
			input: CreateOrderInput{
				UserID:          1,
				GameID:          1,
				TotalPriceCents: 10000,
				Currency:        "INVALID",
			},
			wantErr: true,
		},
		{
			name: "结束时间早于开始时间",
			input: CreateOrderInput{
				UserID:          1,
				GameID:          1,
				TotalPriceCents: 10000,
				Currency:        model.CurrencyCNY,
				ScheduledStart:  timePtr(time.Now().Add(2 * time.Hour)),
				ScheduledEnd:    timePtr(time.Now().Add(1 * time.Hour)),
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAdminService(
				&fakeGameRepo{},
				&fakeUserRepo{},
				&fakePlayerRepo{},
				&fakeOrderRepo{},
				&fakePaymentRepo{},
				&fakeRoleRepo{},
				cache.NewMemory(),
			)
			
			order, err := svc.CreateOrder(ctx, tt.input)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, order)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, order)
			}
		})
	}
}

// Helper functions
func timePtr(t time.Time) *time.Time {
	return &t
}

// Helper functions removed - not needed for current tests
