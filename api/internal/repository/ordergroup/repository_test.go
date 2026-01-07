package ordergroup

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
)

// TestOrderGroup_UpdateStatusFromSubOrders 测试根据子订单更新主订单状态
func TestOrderGroup_UpdateStatusFromSubOrders(t *testing.T) {
	tests := []struct {
		name           string
		subOrders      []model.Order
		expectedStatus model.OrderGroupStatus
		expectedHours  int
	}{
		{
			name:           "空子订单列表",
			subOrders:      []model.Order{},
			expectedStatus: model.OrderGroupStatusPending, // 不变
			expectedHours:  0,
		},
		{
			name: "全部待处理",
			subOrders: []model.Order{
				{Status: model.OrderStatusPending},
				{Status: model.OrderStatusPending},
				{Status: model.OrderStatusConfirmed},
			},
			expectedStatus: model.OrderGroupStatusPending,
			expectedHours:  0,
		},
		{
			name: "全部完成",
			subOrders: []model.Order{
				{Status: model.OrderStatusCompleted},
				{Status: model.OrderStatusCompleted},
				{Status: model.OrderStatusCompleted},
			},
			expectedStatus: model.OrderGroupStatusCompleted,
			expectedHours:  3,
		},
		{
			name: "全部取消",
			subOrders: []model.Order{
				{Status: model.OrderStatusCanceled},
				{Status: model.OrderStatusCanceled},
			},
			expectedStatus: model.OrderGroupStatusCanceled,
			expectedHours:  0,
		},
		{
			name: "进行中",
			subOrders: []model.Order{
				{Status: model.OrderStatusCompleted},
				{Status: model.OrderStatusInProgress},
				{Status: model.OrderStatusPending},
			},
			expectedStatus: model.OrderGroupStatusInProgress,
			expectedHours:  1,
		},
		{
			name: "部分完成（有待处理）",
			subOrders: []model.Order{
				{Status: model.OrderStatusCompleted},
				{Status: model.OrderStatusPending},
				{Status: model.OrderStatusPending},
			},
			expectedStatus: model.OrderGroupStatusPartial,
			expectedHours:  1,
		},
		{
			name: "混合状态 - 完成和退款",
			subOrders: []model.Order{
				{Status: model.OrderStatusCompleted},
				{Status: model.OrderStatusCompleted},
				{Status: model.OrderStatusRefunded},
			},
			expectedStatus: model.OrderGroupStatusPending, // 根据实际逻辑，不满足任何条件时保持 pending
			expectedHours:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &model.OrderGroup{
				Status:         model.OrderGroupStatusPending,
				CompletedHours: 0,
			}

			group.UpdateStatusFromSubOrders(tt.subOrders)

			assert.Equal(t, tt.expectedHours, group.CompletedHours)
			// 注意：空列表时状态不变
			if len(tt.subOrders) > 0 {
				assert.Equal(t, tt.expectedStatus, group.Status)
			}
		})
	}
}

// TestGenerateGroupOrderNo 测试主订单号生成
func TestGenerateGroupOrderNo(t *testing.T) {
	orderNo := model.GenerateGroupOrderNo()

	assert.NotEmpty(t, orderNo)
	assert.True(t, len(orderNo) > 1)
	assert.Equal(t, byte('G'), orderNo[0]) // G 开头
}

// TestListOptions_Defaults 测试列表选项默认值
func TestListOptions_Defaults(t *testing.T) {
	opts := ListOptions{}

	assert.Nil(t, opts.UserID)
	assert.Nil(t, opts.Status)
	assert.Nil(t, opts.GameID)
	assert.Equal(t, 0, opts.Page)
	assert.Equal(t, 0, opts.PageSize)
}

// TestListOptions_WithFilters 测试列表选项过滤
func TestListOptions_WithFilters(t *testing.T) {
	userID := uint64(1)
	status := model.OrderGroupStatusInProgress
	gameID := uint64(10)

	opts := ListOptions{
		UserID:   &userID,
		Status:   &status,
		GameID:   &gameID,
		Page:     2,
		PageSize: 20,
	}

	assert.Equal(t, uint64(1), *opts.UserID)
	assert.Equal(t, model.OrderGroupStatusInProgress, *opts.Status)
	assert.Equal(t, uint64(10), *opts.GameID)
	assert.Equal(t, 2, opts.Page)
	assert.Equal(t, 20, opts.PageSize)
}
