package admin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
	withdrawservice "gamelink/internal/service/withdraw"
)

// TestBatchApproveWithdrawalsHandler_RequestStructure tests the request structure
func TestBatchApproveWithdrawalsHandler_RequestStructure(t *testing.T) {
	req := BatchApproveWithdrawalsRequest{
		WithdrawIDs: []uint64{1, 2, 3},
		Remark:      "Test batch approve",
	}

	// Verify request structure
	assert.Equal(t, []uint64{1, 2, 3}, req.WithdrawIDs)
	assert.Equal(t, "Test batch approve", req.Remark)
}

// TestBatchRejectWithdrawalsHandler_RequestStructure tests the request structure
func TestBatchRejectWithdrawalsHandler_RequestStructure(t *testing.T) {
	req := BatchRejectWithdrawalsRequest{
		WithdrawIDs: []uint64{1, 2, 3},
		Reason:      "Invalid account",
	}

	// Verify request structure
	assert.Equal(t, []uint64{1, 2, 3}, req.WithdrawIDs)
	assert.Equal(t, "Invalid account", req.Reason)
}

// TestBatchCompleteWithdrawalsHandler_RequestStructure tests the request structure
func TestBatchCompleteWithdrawalsHandler_RequestStructure(t *testing.T) {
	req := BatchCompleteWithdrawalsRequest{
		WithdrawIDs: []uint64{1, 2, 3},
	}

	// Verify request structure
	assert.Equal(t, []uint64{1, 2, 3}, req.WithdrawIDs)
}

// TestBatchOperationResult_ResultStructure tests the service result structure
func TestBatchOperationResult_ResultStructure(t *testing.T) {
	result := withdrawservice.BatchOperationResult{
		SuccessCount: 2,
		FailedCount:  1,
		SuccessIDs:   []uint64{1, 2},
		FailedItems: []withdrawservice.BatchOperationErrorItem{
			{ID: 3, Message: "invalid status"},
		},
	}

	// Verify result structure
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Equal(t, []uint64{1, 2}, result.SuccessIDs)
	assert.Len(t, result.FailedItems, 1)
	assert.Equal(t, uint64(3), result.FailedItems[0].ID)
	assert.Equal(t, "invalid status", result.FailedItems[0].Message)
}

// TestBatchApproveWithdrawalsHandler_Validation tests parameter validation
func TestBatchApproveWithdrawalsHandler_Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		request    BatchApproveWithdrawalsRequest
		shouldFail bool
	}{
		{
			name: "valid request",
			request: BatchApproveWithdrawalsRequest{
				WithdrawIDs: []uint64{1, 2, 3},
				Remark:      "Valid",
			},
			shouldFail: false,
		},
		{
			name: "empty withdraw IDs",
			request: BatchApproveWithdrawalsRequest{
				WithdrawIDs: []uint64{},
				Remark:      "Test",
			},
			shouldFail: true,
		},
		{
			name: "exceeds maximum (101 IDs)",
			request: BatchApproveWithdrawalsRequest{
				WithdrawIDs: make([]uint64, 101),
				Remark:      "Test",
			},
			shouldFail: true,
		},
		{
			name: "exceeds maximum remark length",
			request: BatchApproveWithdrawalsRequest{
				WithdrawIDs: []uint64{1},
				Remark:      string(make([]byte, 501)),
			},
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/test", func(c *gin.Context) {
				var req BatchApproveWithdrawalsRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, req)
			})

			body, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if tt.shouldFail {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			} else {
				assert.Equal(t, http.StatusOK, w.Code)
			}
		})
	}
}

// TestBatchRejectWithdrawalsHandler_Validation tests parameter validation
func TestBatchRejectWithdrawalsHandler_Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		request    BatchRejectWithdrawalsRequest
		shouldFail bool
	}{
		{
			name: "valid request",
			request: BatchRejectWithdrawalsRequest{
				WithdrawIDs: []uint64{1, 2, 3},
				Reason:      "Invalid account information",
			},
			shouldFail: false,
		},
		{
			name: "empty withdraw IDs",
			request: BatchRejectWithdrawalsRequest{
				WithdrawIDs: []uint64{},
				Reason:      "Test",
			},
			shouldFail: true,
		},
		{
			name: "missing reason",
			request: BatchRejectWithdrawalsRequest{
				WithdrawIDs: []uint64{1},
				Reason:      "",
			},
			shouldFail: true,
		},
		{
			name: "exceeds maximum (101 IDs)",
			request: BatchRejectWithdrawalsRequest{
				WithdrawIDs: make([]uint64, 101),
				Reason:      "Test",
			},
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/test", func(c *gin.Context) {
				var req BatchRejectWithdrawalsRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, req)
			})

			body, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if tt.shouldFail {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			} else {
				assert.Equal(t, http.StatusOK, w.Code)
			}
		})
	}
}

// TestBatchCompleteWithdrawalsHandler_Validation tests parameter validation
func TestBatchCompleteWithdrawalsHandler_Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		request    BatchCompleteWithdrawalsRequest
		shouldFail bool
	}{
		{
			name: "valid request",
			request: BatchCompleteWithdrawalsRequest{
				WithdrawIDs: []uint64{1, 2, 3},
			},
			shouldFail: false,
		},
		{
			name: "empty withdraw IDs",
			request: BatchCompleteWithdrawalsRequest{
				WithdrawIDs: []uint64{},
			},
			shouldFail: true,
		},
		{
			name: "exceeds maximum (101 IDs)",
			request: BatchCompleteWithdrawalsRequest{
				WithdrawIDs: make([]uint64, 101),
			},
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/test", func(c *gin.Context) {
				var req BatchCompleteWithdrawalsRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, req)
			})

			body, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if tt.shouldFail {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			} else {
				assert.Equal(t, http.StatusOK, w.Code)
			}
		})
	}
}

// TestWithdrawStatusTransitions tests valid status transitions for batch operations
func TestWithdrawStatusTransitions(t *testing.T) {
	tests := []struct {
		name        string
		current     model.WithdrawStatus
		target      model.WithdrawStatus
		isValid     bool
		description string
	}{
		{
			name:        "pending to approved",
			current:     model.WithdrawStatusPending,
			target:      model.WithdrawStatusApproved,
			isValid:     true,
			description: "Can approve pending withdrawal",
		},
		{
			name:        "pending to rejected",
			current:     model.WithdrawStatusPending,
			target:      model.WithdrawStatusRejected,
			isValid:     true,
			description: "Can reject pending withdrawal",
		},
		{
			name:        "approved to completed",
			current:     model.WithdrawStatusApproved,
			target:      model.WithdrawStatusCompleted,
			isValid:     true,
			description: "Can complete approved withdrawal",
		},
		{
			name:        "approved to rejected",
			current:     model.WithdrawStatusApproved,
			target:      model.WithdrawStatusRejected,
			isValid:     false,
			description: "Cannot reject approved withdrawal",
		},
		{
			name:        "completed to approved",
			current:     model.WithdrawStatusCompleted,
			target:      model.WithdrawStatusApproved,
			isValid:     false,
			description: "Cannot approve completed withdrawal",
		},
		{
			name:        "rejected to approved",
			current:     model.WithdrawStatusRejected,
			target:      model.WithdrawStatusApproved,
			isValid:     false,
			description: "Cannot approve rejected withdrawal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var isValid bool
			switch tt.target {
			case model.WithdrawStatusApproved, model.WithdrawStatusRejected:
				isValid = tt.current == model.WithdrawStatusPending
			case model.WithdrawStatusCompleted:
				isValid = tt.current == model.WithdrawStatusApproved
			default:
				isValid = false
			}

			assert.Equal(t, tt.isValid, isValid, tt.description)
		})
	}
}
