package model_test

import (
	"encoding/json"
	"testing"
	"time"

	"gamelink/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestWithdrawModel(t *testing.T) {
	now := time.Now()
	processedAt := now.Add(1 * time.Hour)
	completedAt := now.Add(2 * time.Hour)
	processedBy := uint64(100)

	withdraw := &model.Withdraw{
		ID:           1,
		PlayerID:     200,
		UserID:       300,
		AmountCents:  10000, // 100元
		Method:       model.WithdrawMethodAlipay,
		AccountInfo:  "alipay_account_123456",
		Status:       model.WithdrawStatusCompleted,
		RejectReason: "",
		AdminRemark:  "审核通过，正常提现",
		ProcessedBy:  &processedBy,
		ProcessedAt:  &processedAt,
		CompletedAt:  &completedAt,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	assert.Equal(t, uint64(1), withdraw.ID)
	assert.Equal(t, uint64(200), withdraw.PlayerID)
	assert.Equal(t, uint64(300), withdraw.UserID)
	assert.Equal(t, int64(10000), withdraw.AmountCents)
	assert.Equal(t, model.WithdrawMethodAlipay, withdraw.Method)
	assert.Equal(t, "alipay_account_123456", withdraw.AccountInfo)
	assert.Equal(t, model.WithdrawStatusCompleted, withdraw.Status)
	assert.Equal(t, "", withdraw.RejectReason)
	assert.Equal(t, "审核通过，正常提现", withdraw.AdminRemark)
	assert.Equal(t, &processedBy, withdraw.ProcessedBy)
	assert.Equal(t, &processedAt, withdraw.ProcessedAt)
	assert.Equal(t, &completedAt, withdraw.CompletedAt)
	assert.Equal(t, now, withdraw.CreatedAt)
	assert.Equal(t, now, withdraw.UpdatedAt)
}

func TestWithdrawJSONSerialization(t *testing.T) {
	now := time.Now()

	withdraw := &model.Withdraw{
		ID:           1,
		PlayerID:     200,
		UserID:       300,
		AmountCents:  5000,
		Method:       model.WithdrawMethodWeChat,
		AccountInfo:  "wechat_account_test",
		Status:       model.WithdrawStatusPending,
		RejectReason: "",
		AdminRemark:  "",
		ProcessedBy:  nil,
		ProcessedAt:  nil,
		CompletedAt:  nil,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// 序列化
	data, err := json.Marshal(withdraw)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "wechat_account_test")
	assert.Contains(t, string(data), "pending")

	// 反序列化
	var decoded model.Withdraw
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, withdraw.ID, decoded.ID)
	assert.Equal(t, withdraw.PlayerID, decoded.PlayerID)
	assert.Equal(t, withdraw.UserID, decoded.UserID)
	assert.Equal(t, withdraw.AmountCents, decoded.AmountCents)
	assert.Equal(t, withdraw.Method, decoded.Method)
	assert.Equal(t, withdraw.Status, decoded.Status)
}

func TestWithdrawTableName(t *testing.T) {
	withdraw := model.Withdraw{}
	tableName := withdraw.TableName()
	assert.Equal(t, "withdraws", tableName)
}

func TestWithdrawConstants(t *testing.T) {
	// 测试提现状态常量
	assert.Equal(t, model.WithdrawStatus("pending"), model.WithdrawStatusPending)
	assert.Equal(t, model.WithdrawStatus("approved"), model.WithdrawStatusApproved)
	assert.Equal(t, model.WithdrawStatus("rejected"), model.WithdrawStatusRejected)
	assert.Equal(t, model.WithdrawStatus("completed"), model.WithdrawStatusCompleted)
	assert.Equal(t, model.WithdrawStatus("failed"), model.WithdrawStatusFailed)

	// 测试提现方式常量
	assert.Equal(t, model.WithdrawMethod("alipay"), model.WithdrawMethodAlipay)
	assert.Equal(t, model.WithdrawMethod("wechat"), model.WithdrawMethodWeChat)
	assert.Equal(t, model.WithdrawMethod("bank"), model.WithdrawMethodBank)
}

func TestWithdrawStatuses(t *testing.T) {
	// 测试所有提现状态
	statuses := []model.WithdrawStatus{
		model.WithdrawStatusPending,
		model.WithdrawStatusApproved,
		model.WithdrawStatusRejected,
		model.WithdrawStatusCompleted,
		model.WithdrawStatusFailed,
	}

	for _, status := range statuses {
		withdraw := &model.Withdraw{
			Status: status,
		}
		assert.Equal(t, status, withdraw.Status)
	}
}

func TestWithdrawMethods(t *testing.T) {
	// 测试所有提现方式
	methods := []model.WithdrawMethod{
		model.WithdrawMethodAlipay,
		model.WithdrawMethodWeChat,
		model.WithdrawMethodBank,
	}

	for _, method := range methods {
		withdraw := &model.Withdraw{
			Method: method,
		}
		assert.Equal(t, method, withdraw.Method)
	}
}

func TestWithdrawZeroValues(t *testing.T) {
	withdraw := &model.Withdraw{
		ID:           0,
		PlayerID:     0,
		UserID:       0,
		AmountCents:  0,
		Method:       "",
		AccountInfo:  "",
		Status:       "",
		RejectReason: "",
		AdminRemark:  "",
		ProcessedBy:  nil,
		ProcessedAt:  nil,
		CompletedAt:  nil,
	}

	assert.Equal(t, uint64(0), withdraw.ID)
	assert.Equal(t, uint64(0), withdraw.PlayerID)
	assert.Equal(t, uint64(0), withdraw.UserID)
	assert.Equal(t, int64(0), withdraw.AmountCents)
	assert.Equal(t, model.WithdrawMethod(""), withdraw.Method)
	assert.Equal(t, "", withdraw.AccountInfo)
	assert.Equal(t, model.WithdrawStatus(""), withdraw.Status)
	assert.Equal(t, "", withdraw.RejectReason)
	assert.Equal(t, "", withdraw.AdminRemark)
	assert.Nil(t, withdraw.ProcessedBy)
	assert.Nil(t, withdraw.ProcessedAt)
	assert.Nil(t, withdraw.CompletedAt)
}

func TestWithdrawEdgeCases(t *testing.T) {
	// 测试大金额
	withdraw1 := &model.Withdraw{
		AmountCents: ^int64(0), // 最大int64值
	}
	assert.Equal(t, ^int64(0), withdraw1.AmountCents)

	// 测试长文本
	longAccountInfo := "这是一个非常长的账户信息，可能包含银行账号、支付宝账号、微信账号等各种支付方式的详细信息，用于测试字符串长度的边界情况。"
	longRejectReason := "这是一个非常长的拒绝原因，管理员会在这里详细说明为什么拒绝这笔提现申请，可能包括风险评估、账户异常、信息不完整等各种原因。"
	longAdminRemark := "这是一个非常长的管理员备注，客服会在这里记录他们的审核过程、分析结果、处理依据、最终决定等等详细信息。"

	withdraw2 := &model.Withdraw{
		AccountInfo:  longAccountInfo,
		RejectReason: longRejectReason,
		AdminRemark:  longAdminRemark,
	}
	assert.Equal(t, longAccountInfo, withdraw2.AccountInfo)
	assert.Equal(t, longRejectReason, withdraw2.RejectReason)
	assert.Equal(t, longAdminRemark, withdraw2.AdminRemark)

	// 测试特殊字符
	withdraw3 := &model.Withdraw{
		AccountInfo:  "account@123#test_with_special-chars",
		RejectReason: "拒绝原因@#$%^&*()_+-=[]{}|;':\",./<>?",
		AdminRemark:  "管理员备注：\"引号\"和'单引号'和@#$%^&*()",
	}
	assert.Equal(t, "account@123#test_with_special-chars", withdraw3.AccountInfo)
	assert.Equal(t, "拒绝原因@#$%^&*()_+-=[]{}|;':\",./<>?", withdraw3.RejectReason)
	assert.Equal(t, "管理员备注：\"引号\"和'单引号'和@#$%^&*()", withdraw3.AdminRemark)

	// 测试零值金额
	withdraw4 := &model.Withdraw{
		AmountCents: 0,
	}
	assert.Equal(t, int64(0), withdraw4.AmountCents)

	// 测试负数金额（虽然业务上不合理）
	withdraw5 := &model.Withdraw{
		AmountCents: -1000,
	}
	assert.Equal(t, int64(-1000), withdraw5.AmountCents)
}

func TestWithdrawTimeFields(t *testing.T) {
	now := time.Now()
	processedAt := now.Add(1 * time.Hour)
	completedAt := now.Add(2 * time.Hour)

	// 测试所有时间字段都有值
	withdraw1 := &model.Withdraw{
		ProcessedAt: &processedAt,
		CompletedAt: &completedAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	assert.NotNil(t, withdraw1.ProcessedAt)
	assert.NotNil(t, withdraw1.CompletedAt)
	assert.Equal(t, processedAt, *withdraw1.ProcessedAt)
	assert.Equal(t, completedAt, *withdraw1.CompletedAt)
	assert.Equal(t, now, withdraw1.CreatedAt)
	assert.Equal(t, now, withdraw1.UpdatedAt)

	// 测试只有创建和更新时间
	withdraw2 := &model.Withdraw{
		ProcessedAt: nil,
		CompletedAt: nil,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	assert.Nil(t, withdraw2.ProcessedAt)
	assert.Nil(t, withdraw2.CompletedAt)
	assert.Equal(t, now, withdraw2.CreatedAt)
	assert.Equal(t, now, withdraw2.UpdatedAt)
}

func TestWithdrawWithProcessor(t *testing.T) {
	processorID := uint64(100)

	// 测试有处理人的情况
	withdraw1 := &model.Withdraw{
		ProcessedBy: &processorID,
	}
	assert.NotNil(t, withdraw1.ProcessedBy)
	assert.Equal(t, uint64(100), *withdraw1.ProcessedBy)

	// 测试没有处理人的情况
	withdraw2 := &model.Withdraw{
		ProcessedBy: nil,
	}
	assert.Nil(t, withdraw2.ProcessedBy)
}

func TestWithdrawJSONFields(t *testing.T) {
	now := time.Now()
	processorID := uint64(50)
	processedAt := now.Add(30 * time.Minute)

	withdraw := &model.Withdraw{
		ID:          1,
		PlayerID:    100,
		UserID:      200,
		AmountCents: 15000,
		Method:      model.WithdrawMethodBank,
		AccountInfo: "bank_account_123456",
		Status:      model.WithdrawStatusApproved,
		AdminRemark: "审核通过",
		ProcessedBy: &processorID,
		ProcessedAt: &processedAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	data, err := json.Marshal(withdraw)
	assert.NoError(t, err)

	// 验证JSON结构
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// 检查必需的字段
	assert.Contains(t, result, "id")
	assert.Contains(t, result, "playerId")
	assert.Contains(t, result, "userId")
	assert.Contains(t, result, "amountCents")
	assert.Contains(t, result, "method")
	assert.Contains(t, result, "accountInfo")
	assert.Contains(t, result, "status")
	assert.Contains(t, result, "adminRemark")
	assert.Contains(t, result, "processedBy")
	assert.Contains(t, result, "processedAt")
	assert.Contains(t, result, "createdAt")
	assert.Contains(t, result, "updatedAt")

	// 验证值
	assert.Equal(t, float64(1), result["id"])
	assert.Equal(t, float64(100), result["playerId"])
	assert.Equal(t, float64(200), result["userId"])
	assert.Equal(t, float64(15000), result["amountCents"])
	assert.Equal(t, "bank", result["method"])
	assert.Equal(t, "bank_account_123456", result["accountInfo"])
	assert.Equal(t, "approved", result["status"])
	assert.Equal(t, "审核通过", result["adminRemark"])
	assert.Equal(t, float64(50), result["processedBy"])
}

func TestWithdrawStatusTransitions(t *testing.T) {
	// 测试不同状态组合的提现记录
	withdraws := []struct {
		status       model.WithdrawStatus
		rejectReason string
		adminRemark  string
		desc         string
	}{
		{model.WithdrawStatusPending, "", "", "待审核"},
		{model.WithdrawStatusApproved, "", "审核通过，准备打款", "已批准"},
		{model.WithdrawStatusRejected, "账户信息不完整", "审核不通过", "已拒绝"},
		{model.WithdrawStatusCompleted, "", "提现成功完成", "已完成"},
		{model.WithdrawStatusFailed, "银行系统异常", "处理失败，已退回", "处理失败"},
	}

	for _, w := range withdraws {
		withdraw := &model.Withdraw{
			Status:       w.status,
			RejectReason: w.rejectReason,
			AdminRemark:  w.adminRemark,
		}
		assert.Equal(t, w.status, withdraw.Status)
		assert.Equal(t, w.rejectReason, withdraw.RejectReason)
		assert.Equal(t, w.adminRemark, withdraw.AdminRemark)
	}
}

func TestWithdrawAccountInfoFormats(t *testing.T) {
	// 测试不同格式的账户信息
	accountInfos := []struct {
		method      model.WithdrawMethod
		accountInfo string
		desc        string
	}{
		{model.WithdrawMethodAlipay, "alipay_user_123456", "支付宝账号"},
		{model.WithdrawMethodWeChat, "wechat_user_abcdef", "微信账号"},
		{model.WithdrawMethodBank, "6222001234567890123", "银行卡号"},
		{model.WithdrawMethodBank, "中国银行-1234567890-张三", "银行账户信息"},
		{model.WithdrawMethodAlipay, "13800138000", "支付宝手机号"},
		{model.WithdrawMethodWeChat, "test_openid_123456789", "微信OpenID"},
	}

	for _, info := range accountInfos {
		withdraw := &model.Withdraw{
			Method:      info.method,
			AccountInfo: info.accountInfo,
		}
		assert.Equal(t, info.method, withdraw.Method)
		assert.Equal(t, info.accountInfo, withdraw.AccountInfo)
	}
}
