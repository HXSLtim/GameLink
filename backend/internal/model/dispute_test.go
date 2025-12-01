package model_test

import (
	"encoding/json"
	"testing"
	"time"

	"gamelink/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestOrderDisputeModel(t *testing.T) {
	now := time.Now()
	deadline := now.Add(30 * time.Minute)
	assignedUserID := uint64(50)
	resolvedUserID := uint64(60)
	evidenceURLs := model.EvidenceURLArray{"http://example.com/evidence1.jpg", "http://example.com/evidence2.png"}

	dispute := &model.OrderDispute{
		Base: model.Base{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		OrderID:            100,
		UserID:             200,
		Status:             model.DisputeStatusPending,
		Reason:             "服务质量问题",
		Description:        "陪玩师没有按照约定时间到达，服务态度也不好",
		EvidenceURLs:       evidenceURLs,
		AssignedToUserID:   &assignedUserID,
		AssignmentSource:   model.AssignmentSourceManual,
		AssignedAt:         &now,
		SLADeadline:        &deadline,
		SLABreached:        false,
		SLABreachedAt:      nil,
		Resolution:         model.ResolutionPending,
		ResolutionAmount:   5000, // 50元
		ResolutionNotes:    "经过调查，情况属实，部分退款",
		ResolvedAt:         &now,
		ResolvedByUserID:   &resolvedUserID,
		RolledBackAt:       nil,
		RolledBackByUserID: nil,
		RollbackReason:     "",
		TraceID:            "TRACE123456",
	}

	assert.Equal(t, uint64(1), dispute.ID)
	assert.Equal(t, uint64(100), dispute.OrderID)
	assert.Equal(t, uint64(200), dispute.UserID)
	assert.Equal(t, model.DisputeStatusPending, dispute.Status)
	assert.Equal(t, "服务质量问题", dispute.Reason)
	assert.Equal(t, "陪玩师没有按照约定时间到达，服务态度也不好", dispute.Description)
	assert.Equal(t, evidenceURLs, dispute.EvidenceURLs)
	assert.Equal(t, &assignedUserID, dispute.AssignedToUserID)
	assert.Equal(t, model.AssignmentSourceManual, dispute.AssignmentSource)
	assert.Equal(t, &now, dispute.AssignedAt)
	assert.Equal(t, &deadline, dispute.SLADeadline)
	assert.False(t, dispute.SLABreached)
	assert.Nil(t, dispute.SLABreachedAt)
	assert.Equal(t, model.ResolutionPending, dispute.Resolution)
	assert.Equal(t, int64(5000), dispute.ResolutionAmount)
	assert.Equal(t, "经过调查，情况属实，部分退款", dispute.ResolutionNotes)
	assert.Equal(t, &now, dispute.ResolvedAt)
	assert.Equal(t, &resolvedUserID, dispute.ResolvedByUserID)
	assert.Nil(t, dispute.RolledBackAt)
	assert.Nil(t, dispute.RolledBackByUserID)
	assert.Equal(t, "", dispute.RollbackReason)
	assert.Equal(t, "TRACE123456", dispute.TraceID)
}

func TestOrderDisputeJSONSerialization(t *testing.T) {
	now := time.Now()
	deadline := now.Add(30 * time.Minute)
	evidenceURLs := model.EvidenceURLArray{"http://example.com/evidence1.jpg"}

	dispute := &model.OrderDispute{
		Base: model.Base{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		OrderID:          100,
		UserID:           200,
		Status:           model.DisputeStatusPending,
		Reason:           "服务质量问题",
		Description:      "服务不符合预期",
		EvidenceURLs:     evidenceURLs,
		SLADeadline:      &deadline,
		Resolution:       model.ResolutionPending,
		ResolutionAmount: 0,
		ResolutionNotes:  "",
		TraceID:          "TRACE123",
	}

	// 序列化
	data, err := json.Marshal(dispute)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "服务质量问题")
	assert.Contains(t, string(data), "TRACE123")

	// 反序列化
	var decoded model.OrderDispute
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, dispute.ID, decoded.ID)
	assert.Equal(t, dispute.OrderID, decoded.OrderID)
	assert.Equal(t, dispute.UserID, decoded.UserID)
	assert.Equal(t, dispute.Status, decoded.Status)
	assert.Equal(t, dispute.Reason, decoded.Reason)
}

func TestOrderDisputeTableName(t *testing.T) {
	dispute := model.OrderDispute{}
	tableName := dispute.TableName()
	assert.Equal(t, "order_disputes", tableName)
}

func TestOrderDisputeConstants(t *testing.T) {
	// 测试争议状态常量
	assert.Equal(t, model.DisputeStatus("pending"), model.DisputeStatusPending)
	assert.Equal(t, model.DisputeStatus("assigned"), model.DisputeStatusAssigned)
	assert.Equal(t, model.DisputeStatus("mediating"), model.DisputeStatusMediating)
	assert.Equal(t, model.DisputeStatus("resolved"), model.DisputeStatusResolved)
	assert.Equal(t, model.DisputeStatus("rejected"), model.DisputeStatusRejected)
	assert.Equal(t, model.DisputeStatus("canceled"), model.DisputeStatusCanceled)

	// 测试解决方案常量
	assert.Equal(t, model.DisputeResolution("refund"), model.ResolutionRefund)
	assert.Equal(t, model.DisputeResolution("partial"), model.ResolutionPartial)
	assert.Equal(t, model.DisputeResolution("reassign"), model.ResolutionReassign)
	assert.Equal(t, model.DisputeResolution("reject"), model.ResolutionReject)
	assert.Equal(t, model.DisputeResolution("pending"), model.ResolutionPending)

	// 测试指派来源常量
	assert.Equal(t, model.AssignmentSource("system"), model.AssignmentSourceSystem)
	assert.Equal(t, model.AssignmentSource("manual"), model.AssignmentSourceManual)
	assert.Equal(t, model.AssignmentSource("team"), model.AssignmentSourceTeam)
}

func TestEvidenceURLArray(t *testing.T) {
	// 测试正常的URL数组
	evidenceURLs := model.EvidenceURLArray{
		"http://example.com/evidence1.jpg",
		"http://example.com/evidence2.png",
		"http://example.com/evidence3.gif",
	}

	assert.Len(t, evidenceURLs, 3)
	assert.Equal(t, "http://example.com/evidence1.jpg", evidenceURLs[0])
	assert.Equal(t, "http://example.com/evidence2.png", evidenceURLs[1])
	assert.Equal(t, "http://example.com/evidence3.gif", evidenceURLs[2])

	// 测试空数组
	emptyEvidence := model.EvidenceURLArray{}
	assert.Len(t, emptyEvidence, 0)

	// 测试单个URL
	singleEvidence := model.EvidenceURLArray{"http://example.com/single.jpg"}
	assert.Len(t, singleEvidence, 1)
	assert.Equal(t, "http://example.com/single.jpg", singleEvidence[0])
}

func TestEvidenceURLArrayJSONSerialization(t *testing.T) {
	evidenceURLs := model.EvidenceURLArray{
		"http://example.com/evidence1.jpg",
		"http://example.com/evidence2.png",
	}

	// 序列化
	data, err := json.Marshal(evidenceURLs)
	assert.NoError(t, err)
	assert.Equal(t, `["http://example.com/evidence1.jpg","http://example.com/evidence2.png"]`, string(data))

	// 反序列化
	var decoded model.EvidenceURLArray
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, evidenceURLs, decoded)
}

func TestEvidenceURLArraySQLInterfaces(t *testing.T) {
	// 测试 Value 方法
	evidenceURLs := model.EvidenceURLArray{
		"http://example.com/evidence1.jpg",
		"http://example.com/evidence2.png",
	}

	value, err := evidenceURLs.Value()
	assert.NoError(t, err)
	assert.NotNil(t, value)

	// 测试 Scan 方法
	jsonData := `["http://example.com/evidence1.jpg","http://example.com/evidence2.png"]`
	var scanned model.EvidenceURLArray
	err = scanned.Scan([]byte(jsonData))
	assert.NoError(t, err)
	assert.Equal(t, evidenceURLs, scanned)

	// 测试 Scan 方法的错误情况
	err = scanned.Scan("invalid data")
	assert.Error(t, err)

	// 测试 Scan 空值
	err = scanned.Scan(nil)
	assert.NoError(t, err)
}

func TestOrderDisputeIsOverSLA(t *testing.T) {
	now := time.Now()

	// 测试未超过SLA的情况
	futureDeadline := now.Add(1 * time.Hour)
	dispute1 := &model.OrderDispute{
		SLADeadline: &futureDeadline,
	}
	assert.False(t, dispute1.IsOverSLA())

	// 测试已超过SLA的情况
	pastDeadline := now.Add(-1 * time.Hour)
	dispute2 := &model.OrderDispute{
		SLADeadline: &pastDeadline,
	}
	assert.True(t, dispute2.IsOverSLA())

	// 测试没有SLA截止时间的情况
	dispute3 := &model.OrderDispute{
		SLADeadline: nil,
	}
	assert.False(t, dispute3.IsOverSLA())
}

func TestOrderDisputeGetSLARemaining(t *testing.T) {
	now := time.Now()

	// 测试还有剩余时间的情况
	futureDeadline := now.Add(30 * time.Minute)
	dispute1 := &model.OrderDispute{
		SLADeadline: &futureDeadline,
	}
	remaining := dispute1.GetSLARemaining()
	assert.Greater(t, remaining, int64(0))
	assert.LessOrEqual(t, remaining, int64(1800)) // 30分钟 = 1800秒

	// 测试已过期的情况
	pastDeadline := now.Add(-1 * time.Hour)
	dispute2 := &model.OrderDispute{
		SLADeadline: &pastDeadline,
	}
	remaining2 := dispute2.GetSLARemaining()
	assert.Equal(t, int64(0), remaining2)

	// 测试没有SLA截止时间的情况
	dispute3 := &model.OrderDispute{
		SLADeadline: nil,
	}
	remaining3 := dispute3.GetSLARemaining()
	assert.Equal(t, int64(0), remaining3)
}

func TestCanInitiateDispute(t *testing.T) {
	now := time.Now()
	recentCompletion := now.Add(-1 * time.Hour)
	oldCompletion := now.Add(-25 * time.Hour) // 超过24小时

	// 测试订单进行中可以发起争议
	order1 := &model.Order{
		Status:      model.OrderStatusInProgress,
		CompletedAt: nil,
	}
	assert.True(t, model.CanInitiateDispute(order1))

	// 测试最近完成的订单可以发起争议
	order2 := &model.Order{
		Status:      model.OrderStatusCompleted,
		CompletedAt: &recentCompletion,
	}
	assert.True(t, model.CanInitiateDispute(order2))

	// 测试很久以前完成的订单不能发起争议
	order3 := &model.Order{
		Status:      model.OrderStatusCompleted,
		CompletedAt: &oldCompletion,
	}
	assert.False(t, model.CanInitiateDispute(order3))

	// 测试没有完成时间的已完成订单（边界情况）
	order4 := &model.Order{
		Status:      model.OrderStatusCompleted,
		CompletedAt: nil,
	}
	assert.False(t, model.CanInitiateDispute(order4))
}

func TestOrderDisputeZeroValues(t *testing.T) {
	dispute := &model.OrderDispute{
		OrderID:            0,
		UserID:             0,
		Status:             "",
		Reason:             "",
		Description:        "",
		EvidenceURLs:       nil,
		AssignedToUserID:   nil,
		AssignmentSource:   "",
		AssignedAt:         nil,
		SLADeadline:        nil,
		SLABreached:        false,
		SLABreachedAt:      nil,
		Resolution:         "",
		ResolutionAmount:   0,
		ResolutionNotes:    "",
		ResolvedAt:         nil,
		ResolvedByUserID:   nil,
		RolledBackAt:       nil,
		RolledBackByUserID: nil,
		RollbackReason:     "",
		TraceID:            "",
	}

	assert.Equal(t, uint64(0), dispute.OrderID)
	assert.Equal(t, uint64(0), dispute.UserID)
	assert.Equal(t, model.DisputeStatus(""), dispute.Status)
	assert.Equal(t, "", dispute.Reason)
	assert.Equal(t, "", dispute.Description)
	assert.Nil(t, dispute.EvidenceURLs)
	assert.Nil(t, dispute.AssignedToUserID)
	assert.Equal(t, model.AssignmentSource(""), dispute.AssignmentSource)
	assert.Nil(t, dispute.AssignedAt)
	assert.Nil(t, dispute.SLADeadline)
	assert.False(t, dispute.SLABreached)
	assert.Nil(t, dispute.SLABreachedAt)
	assert.Equal(t, model.DisputeResolution(""), dispute.Resolution)
	assert.Equal(t, int64(0), dispute.ResolutionAmount)
	assert.Equal(t, "", dispute.ResolutionNotes)
	assert.Nil(t, dispute.ResolvedAt)
	assert.Nil(t, dispute.ResolvedByUserID)
	assert.Nil(t, dispute.RolledBackAt)
	assert.Nil(t, dispute.RolledBackByUserID)
	assert.Equal(t, "", dispute.RollbackReason)
	assert.Equal(t, "", dispute.TraceID)
}

func TestOrderDisputeEdgeCases(t *testing.T) {
	// 测试大金额退款
	dispute1 := &model.OrderDispute{
		ResolutionAmount: ^int64(0), // 最大int64值
	}
	assert.Equal(t, ^int64(0), dispute1.ResolutionAmount)

	// 测试长文本
	longReason := "这是一个非常长的争议原因，可能包含很多详细信息，用户会在这里描述他们遇到的问题，包括具体的情况、时间、地点、涉及的人员等等。"
	longDescription := "这是一个非常长的争议描述，可以包含更多详细信息，比如用户的具体体验、他们期望的解决方案、相关的背景信息等等。这种长文本测试可以确保我们的模型能够处理各种长度的输入。"
	longResolutionNotes := "这是一个非常长的处理备注，客服会在这里记录他们的调查过程、分析结果、处理依据、最终决定等等。"
	longRollbackReason := "这是一个非常长的回退原因，管理员会在这里详细说明为什么要回退之前的处理决定。"

	dispute2 := &model.OrderDispute{
		Reason:          longReason,
		Description:     longDescription,
		ResolutionNotes: longResolutionNotes,
		RollbackReason:  longRollbackReason,
	}
	assert.Equal(t, longReason, dispute2.Reason)
	assert.Equal(t, longDescription, dispute2.Description)
	assert.Equal(t, longResolutionNotes, dispute2.ResolutionNotes)
	assert.Equal(t, longRollbackReason, dispute2.RollbackReason)

	// 测试特殊字符
	dispute3 := &model.OrderDispute{
		Reason:          "服务质量问题@#$%^&*()",
		Description:     "描述包含特殊字符：<>{}[]|\\",
		ResolutionNotes: "处理备注：\"引号\"和'单引号'",
		TraceID:         "TRACE_123#456@789",
	}
	assert.Equal(t, "服务质量问题@#$%^&*()", dispute3.Reason)
	assert.Equal(t, "描述包含特殊字符：<>{}[]|\\", dispute3.Description)
	assert.Equal(t, "处理备注：\"引号\"和'单引号'", dispute3.ResolutionNotes)
	assert.Equal(t, "TRACE_123#456@789", dispute3.TraceID)

	// 测试所有争议状态
	statuses := []model.DisputeStatus{
		model.DisputeStatusPending,
		model.DisputeStatusAssigned,
		model.DisputeStatusMediating,
		model.DisputeStatusResolved,
		model.DisputeStatusRejected,
		model.DisputeStatusCanceled,
	}

	for _, status := range statuses {
		dispute := &model.OrderDispute{
			Status: status,
		}
		assert.Equal(t, status, dispute.Status)
	}

	// 测试所有解决方案
	resolutions := []model.DisputeResolution{
		model.ResolutionRefund,
		model.ResolutionPartial,
		model.ResolutionReassign,
		model.ResolutionReject,
		model.ResolutionPending,
	}

	for _, resolution := range resolutions {
		dispute := &model.OrderDispute{
			Resolution: resolution,
		}
		assert.Equal(t, resolution, dispute.Resolution)
	}

	// 测试所有指派来源
	sources := []model.AssignmentSource{
		model.AssignmentSourceSystem,
		model.AssignmentSourceManual,
		model.AssignmentSourceTeam,
	}

	for _, source := range sources {
		dispute := &model.OrderDispute{
			AssignmentSource: source,
		}
		assert.Equal(t, source, dispute.AssignmentSource)
	}
}

func TestOrderDisputeRelations(t *testing.T) {
	dispute := &model.OrderDispute{
		Order: model.Order{
			OrderNo: "ORDER123",
		},
		User: model.User{
			Name: "test user",
		},
		AssignedToUser: &model.User{
			Name: "assigned user",
		},
		ResolvedByUser: &model.User{
			Name: "resolved user",
		},
		RolledBackByUser: &model.User{
			Name: "rollback user",
		},
	}

	assert.Equal(t, "ORDER123", dispute.Order.OrderNo)
	assert.Equal(t, "test user", dispute.User.Name)
	assert.NotNil(t, dispute.AssignedToUser)
	assert.Equal(t, "assigned user", dispute.AssignedToUser.Name)
	assert.NotNil(t, dispute.ResolvedByUser)
	assert.Equal(t, "resolved user", dispute.ResolvedByUser.Name)
	assert.NotNil(t, dispute.RolledBackByUser)
	assert.Equal(t, "rollback user", dispute.RolledBackByUser.Name)
}
