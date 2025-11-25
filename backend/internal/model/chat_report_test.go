package model_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gamelink/internal/model"
)

func TestChatReportModel(t *testing.T) {
	now := time.Now()
	handledAt := now.Add(2 * time.Hour)
	handledBy := uint64(100)

	chatReport := &model.ChatReport{
		Base: model.Base{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		MessageID:  200,
		ReporterID: 300,
		Reason:     "不当言论",
		Evidence:   "用户发布了包含侮辱性语言的消息",
		Status:     "approved",
		HandledBy:  &handledBy,
		HandledAt:  &handledAt,
		Notes:      "经过审核，举报属实，已对消息进行处理",
	}

	assert.Equal(t, uint64(1), chatReport.ID)
	assert.Equal(t, uint64(200), chatReport.MessageID)
	assert.Equal(t, uint64(300), chatReport.ReporterID)
	assert.Equal(t, "不当言论", chatReport.Reason)
	assert.Equal(t, "用户发布了包含侮辱性语言的消息", chatReport.Evidence)
	assert.Equal(t, "approved", chatReport.Status)
	assert.Equal(t, &handledBy, chatReport.HandledBy)
	assert.Equal(t, &handledAt, chatReport.HandledAt)
	assert.Equal(t, "经过审核，举报属实，已对消息进行处理", chatReport.Notes)
}

func TestChatReportJSONSerialization(t *testing.T) {
	now := time.Now()

	chatReport := &model.ChatReport{
		Base: model.Base{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		MessageID:  200,
		ReporterID: 300,
		Reason:     "垃圾信息",
		Evidence:   "用户重复发送相同内容",
		Status:     "pending",
		HandledBy:  nil,
		HandledAt:  nil,
		Notes:      "",
	}

	// 序列化
	data, err := json.Marshal(chatReport)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "垃圾信息")
	assert.Contains(t, string(data), "pending")

	// 反序列化
	var decoded model.ChatReport
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, chatReport.ID, decoded.ID)
	assert.Equal(t, chatReport.MessageID, decoded.MessageID)
	assert.Equal(t, chatReport.ReporterID, decoded.ReporterID)
	assert.Equal(t, chatReport.Reason, decoded.Reason)
	assert.Equal(t, chatReport.Status, decoded.Status)
}

func TestChatReportTableName(t *testing.T) {
	chatReport := model.ChatReport{}
	tableName := chatReport.TableName()
	assert.Equal(t, "chat_reports", tableName)
}

func TestChatReportZeroValues(t *testing.T) {
	chatReport := &model.ChatReport{
		MessageID:  0,
		ReporterID: 0,
		Reason:     "",
		Evidence:   "",
		Status:     "",
		HandledBy:  nil,
		HandledAt:  nil,
		Notes:      "",
	}

	assert.Equal(t, uint64(0), chatReport.MessageID)
	assert.Equal(t, uint64(0), chatReport.ReporterID)
	assert.Equal(t, "", chatReport.Reason)
	assert.Equal(t, "", chatReport.Evidence)
	assert.Equal(t, "", chatReport.Status)
	assert.Nil(t, chatReport.HandledBy)
	assert.Nil(t, chatReport.HandledAt)
	assert.Equal(t, "", chatReport.Notes)
}

func TestChatReportEdgeCases(t *testing.T) {
	// 测试长文本
	longReason := "这是一个非常长的举报原因，可能包含很多详细信息，举报人会在这里描述他们发现的问题，包括具体的情况、时间、涉及的聊天消息内容等等。这种长文本测试可以确保我们的模型能够处理各种长度的输入。"
	longEvidence := "这是一个非常长的证据描述，可以包含更多详细信息，比如具体的聊天消息内容、截图描述、相关的上下文信息等等。这种长文本测试可以确保我们的模型能够处理各种长度的输入。"
	longNotes := "这是一个非常长的处理备注，审核人员会在这里记录他们的审核过程、分析结果、处理依据、最终决定等等详细信息。"

	chatReport1 := &model.ChatReport{
		Reason:     longReason,
		Evidence:   longEvidence,
		Notes:      longNotes,
	}
	assert.Equal(t, longReason, chatReport1.Reason)
	assert.Equal(t, longEvidence, chatReport1.Evidence)
	assert.Equal(t, longNotes, chatReport1.Notes)

	// 测试特殊字符
	chatReport2 := &model.ChatReport{
		Reason:     "举报@#$%^&*()_+-=[]{}|\\",
		Evidence:   "证据包含特殊字符：<>{}[]|\\\"quotes\" and 'apostrophes' and @#$%^&*()_+-=[]{}|;':\",./<>?",
		Notes:      "备注：\"引号\"和'单引号'和@#$%^&*()",
	}
	assert.Equal(t, "举报@#$%^&*()_+-=[]{}|\\", chatReport2.Reason)
	assert.Contains(t, chatReport2.Evidence, "特殊字符")
	assert.Equal(t, "备注：\"引号\"和'单引号'和@#$%^&*()", chatReport2.Notes)

	// 测试大数值
	chatReport3 := &model.ChatReport{
		MessageID:  ^uint64(0), // 最大uint64值
		ReporterID: ^uint64(0),
		HandledBy:  &[]uint64{^uint64(0)}[0],
	}
	assert.Equal(t, ^uint64(0), chatReport3.MessageID)
	assert.Equal(t, ^uint64(0), chatReport3.ReporterID)
	assert.Equal(t, ^uint64(0), *chatReport3.HandledBy)
}

func TestChatReportStatuses(t *testing.T) {
	// 测试不同状态
	statuses := []string{"pending", "approved", "rejected", "dismissed"}

	for _, status := range statuses {
		chatReport := &model.ChatReport{
			Status: status,
		}
		assert.Equal(t, status, chatReport.Status)
	}
}

func TestChatReportWithHandler(t *testing.T) {
	handlerID := uint64(150)
	handledAt := time.Now()

	// 测试有处理人的情况
	chatReport1 := &model.ChatReport{
		HandledBy: &handlerID,
		HandledAt: &handledAt,
	}
	assert.NotNil(t, chatReport1.HandledBy)
	assert.NotNil(t, chatReport1.HandledAt)
	assert.Equal(t, uint64(150), *chatReport1.HandledBy)
	assert.Equal(t, handledAt, *chatReport1.HandledAt)

	// 测试没有处理人的情况
	chatReport2 := &model.ChatReport{
		HandledBy: nil,
		HandledAt: nil,
	}
	assert.Nil(t, chatReport2.HandledBy)
	assert.Nil(t, chatReport2.HandledAt)
}

func TestChatReportTimeFields(t *testing.T) {
	now := time.Now()
	handledAt := now.Add(1 * time.Hour)

	// 测试所有时间字段
	chatReport := &model.ChatReport{
		Base: model.Base{
			CreatedAt: now,
			UpdatedAt: now,
		},
		HandledAt: &handledAt,
	}

	assert.Equal(t, now, chatReport.CreatedAt)
	assert.Equal(t, now, chatReport.UpdatedAt)
	assert.Equal(t, &handledAt, chatReport.HandledAt)
}

func TestChatReportJSONFields(t *testing.T) {
	handlerID := uint64(50)
	handledAt := time.Now()

	chatReport := &model.ChatReport{
		Base: model.Base{
			ID: 1,
		},
		MessageID:  100,
		ReporterID: 200,
		Reason:     "垃圾信息",
		Evidence:   "用户重复发送广告内容",
		Status:     "approved",
		HandledBy:  &handlerID,
		HandledAt:  &handledAt,
		Notes:      "举报属实，已处理",
	}

	data, err := json.Marshal(chatReport)
	assert.NoError(t, err)

	// 验证JSON结构
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// 检查必需的字段
	assert.Contains(t, result, "id")
	assert.Contains(t, result, "messageId")
	assert.Contains(t, result, "reporterId")
	assert.Contains(t, result, "reason")
	assert.Contains(t, result, "evidence")
	assert.Contains(t, result, "status")
	assert.Contains(t, result, "handledBy")
	assert.Contains(t, result, "handledAt")
	assert.Contains(t, result, "notes")

	// 验证值
	assert.Equal(t, float64(1), result["id"])
	assert.Equal(t, float64(100), result["messageId"])
	assert.Equal(t, float64(200), result["reporterId"])
	assert.Equal(t, "垃圾信息", result["reason"])
	assert.Equal(t, "用户重复发送广告内容", result["evidence"])
	assert.Equal(t, "approved", result["status"])
	assert.Equal(t, float64(50), result["handledBy"])
}

func TestChatReportMultilingualContent(t *testing.T) {
	// 测试多语言内容
	reports := []struct {
		reason   string
		evidence string
		notes    string
		lang     string
	}{
		{"Inappropriate language", "User posted messages with offensive content", "Report verified, message has been removed", "English"},
		{"不当言论", "用户发布了包含冒犯性内容的消息", "举报属实，消息已被移除", "Chinese"},
		{"Langage inapproprié", "L'utilisateur a publié des messages avec du contenu offensant", "Signalement vérifié, le message a été supprimé", "French"},
		{"Lenguaje inapropiado", "El usuario publicó mensajes con contenido ofensivo", "Reporte verificado, el mensaje ha sido eliminado", "Spanish"},
		{"Linguaggio inappropriato", "L'utente ha pubblicato messaggi con contenuti offensivi", "Segnalazione verificata, il messaggio è stato rimosso", "Italian"},
	}

	for _, report := range reports {
		chatReport := &model.ChatReport{
			Reason:   report.reason,
			Evidence: report.evidence,
			Notes:    report.notes,
		}
		assert.Equal(t, report.reason, chatReport.Reason, "Failed for language: %s", report.lang)
		assert.Equal(t, report.evidence, chatReport.Evidence, "Failed for language: %s", report.lang)
		assert.Equal(t, report.notes, chatReport.Notes, "Failed for language: %s", report.lang)
	}
}

func TestChatReportSpecialCharacters(t *testing.T) {
	// 测试包含特殊字符的内容
	chatReport := &model.ChatReport{
		Reason:     "举报@#$%^&*()_+-=[]{}|\\",
		Evidence:   "证据包含特殊字符：<>{}[]|\\\"quotes\" and 'apostrophes' and @#$%^&*()_+-=[]{}|;':\",./<>?",
		Notes:      "备注：\"引号\"和'单引号'和@#$%^&*()😊🚀",
	}

	assert.Equal(t, "举报@#$%^&*()_+-=[]{}|\\", chatReport.Reason)
	assert.Contains(t, chatReport.Evidence, "特殊字符")
	assert.Equal(t, "备注：\"引号\"和'单引号'和@#$%^&*()😊🚀", chatReport.Notes)
}

func TestChatReportEmptyFields(t *testing.T) {
	chatReport := &model.ChatReport{}

	assert.Equal(t, uint64(0), chatReport.ID)
	assert.Equal(t, uint64(0), chatReport.MessageID)
	assert.Equal(t, uint64(0), chatReport.ReporterID)
	assert.Equal(t, "", chatReport.Reason)
	assert.Equal(t, "", chatReport.Evidence)
	assert.Equal(t, "", chatReport.Status)
	assert.Nil(t, chatReport.HandledBy)
	assert.Nil(t, chatReport.HandledAt)
	assert.Equal(t, "", chatReport.Notes)
}