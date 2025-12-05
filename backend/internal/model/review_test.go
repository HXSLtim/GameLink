package model_test

import (
	"encoding/json"
	"testing"

	"gamelink/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestReviewModel(t *testing.T) {
	review := &model.Review{
		Base: model.Base{
			ID: 1,
		},
		OrderID:         100,
		UserID:          200,
		PlayerID:        300,
		Score:           5,
		Content:         "非常好的服务，陪玩师很专业！",
		Status:          model.ReviewStatusApproved,
		IsReported:      false,
		Images:          model.StringArray{"https://example.com/image1.jpg"},
		RejectionReason: "",
	}

	assert.Equal(t, uint64(1), review.ID)
	assert.Equal(t, uint64(100), review.OrderID)
	assert.Equal(t, uint64(200), review.UserID)
	assert.Equal(t, uint64(300), review.PlayerID)
	assert.Equal(t, model.Rating(5), review.Score)
	assert.Equal(t, "非常好的服务，陪玩师很专业！", review.Content)
	assert.Equal(t, model.ReviewStatusApproved, review.Status)
	assert.False(t, review.IsReported)
	assert.Len(t, review.Images, 1)
	assert.Equal(t, "", review.RejectionReason)
}

func TestReviewJSONSerialization(t *testing.T) {
	review := &model.Review{
		Base: model.Base{
			ID: 1,
		},
		OrderID:  100,
		UserID:   200,
		PlayerID: 300,
		Score:    4,
		Content:  "Good service, very professional!",
	}

	// 序列化
	data, err := json.Marshal(review)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "Good service")
	assert.Contains(t, string(data), "4")

	// 反序列化
	var decoded model.Review
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, review.ID, decoded.ID)
	assert.Equal(t, review.OrderID, decoded.OrderID)
	assert.Equal(t, review.UserID, decoded.UserID)
	assert.Equal(t, review.PlayerID, decoded.PlayerID)
	assert.Equal(t, review.Score, decoded.Score)
	assert.Equal(t, review.Content, decoded.Content)
}

func TestReviewRatingValues(t *testing.T) {
	// 测试所有有效的评分值
	validRatings := []model.Rating{1, 2, 3, 4, 5}

	for _, rating := range validRatings {
		review := &model.Review{
			Score: rating,
		}
		assert.Equal(t, rating, review.Score)
	}
}

func TestReviewZeroValues(t *testing.T) {
	review := &model.Review{
		OrderID:  0,
		UserID:   0,
		PlayerID: 0,
		Score:    0,
		Content:  "",
	}

	assert.Equal(t, uint64(0), review.OrderID)
	assert.Equal(t, uint64(0), review.UserID)
	assert.Equal(t, uint64(0), review.PlayerID)
	assert.Equal(t, model.Rating(0), review.Score)
	assert.Equal(t, "", review.Content)
}

func TestReviewEmptyContent(t *testing.T) {
	review := &model.Review{
		OrderID:  100,
		UserID:   200,
		PlayerID: 300,
		Score:    5,
		Content:  "", // 空内容
	}

	assert.Equal(t, "", review.Content)
}

func TestReviewLongContent(t *testing.T) {
	longContent := "这是一段非常长的评价内容，可以包含很多详细信息。用户可能会在这里写下他们对服务的完整体验，包括优点、缺点、建议等等。这种长内容测试可以确保我们的模型能够处理各种长度的文本输入。"

	review := &model.Review{
		OrderID:  100,
		UserID:   200,
		PlayerID: 300,
		Score:    5,
		Content:  longContent,
	}

	assert.Equal(t, longContent, review.Content)
}

func TestReviewSpecialCharacters(t *testing.T) {
	review := &model.Review{
		OrderID:  100,
		UserID:   200,
		PlayerID: 300,
		Score:    4,
		Content:  "服务很好！😊 @user #推荐 5/5 ⭐⭐⭐⭐⭐",
	}

	assert.Equal(t, "服务很好！😊 @user #推荐 5/5 ⭐⭐⭐⭐⭐", review.Content)
}

func TestReviewJSONFields(t *testing.T) {
	review := &model.Review{
		Base: model.Base{
			ID: 123,
		},
		OrderID:  456,
		UserID:   789,
		PlayerID: 101112,
		Score:    5,
		Content:  "Excellent service!",
	}

	data, err := json.Marshal(review)
	assert.NoError(t, err)

	// 验证JSON结构
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// 检查必需的字段
	assert.Contains(t, result, "id")
	assert.Contains(t, result, "orderId")
	assert.Contains(t, result, "reviewerId") // 注意：JSON标签是reviewerId
	assert.Contains(t, result, "playerId")
	assert.Contains(t, result, "rating")
	assert.Contains(t, result, "comment")

	// 验证值
	assert.Equal(t, float64(123), result["id"])
	assert.Equal(t, float64(456), result["orderId"])
	assert.Equal(t, float64(789), result["reviewerId"])
	assert.Equal(t, float64(101112), result["playerId"])
	assert.Equal(t, float64(5), result["rating"])
	assert.Equal(t, "Excellent service!", result["comment"])
}

func TestReviewEdgeCases(t *testing.T) {
	// 测试超出正常范围的评分（虽然业务上不应该出现）
	review1 := &model.Review{
		Score: model.Rating(0), // 0分
	}
	assert.Equal(t, model.Rating(0), review1.Score)

	review2 := &model.Review{
		Score: model.Rating(6), // 6分（超出5分制）
	}
	assert.Equal(t, model.Rating(6), review2.Score)

	review3 := &model.Review{
		Score: model.Rating(10), // 10分
	}
	assert.Equal(t, model.Rating(10), review3.Score)

	// 测试包含换行符的内容
	review4 := &model.Review{
		Content: "第一行\n第二行\n第三行",
	}
	assert.Equal(t, "第一行\n第二行\n第三行", review4.Content)

	// 测试包含引号的内容
	review5 := &model.Review{
		Content: `用户说："服务非常好"，推荐给大家！`,
	}
	assert.Equal(t, `用户说："服务非常好"，推荐给大家！`, review5.Content)
}

func TestReviewEmptyFields(t *testing.T) {
	review := &model.Review{}

	assert.Equal(t, uint64(0), review.ID)
	assert.Equal(t, uint64(0), review.OrderID)
	assert.Equal(t, uint64(0), review.UserID)
	assert.Equal(t, uint64(0), review.PlayerID)
	assert.Equal(t, model.Rating(0), review.Score)
	assert.Equal(t, "", review.Content)
}

func TestReviewWithBaseFields(t *testing.T) {
	review := &model.Review{
		Base: model.Base{
			ID: 999,
		},
		OrderID:  111,
		PlayerID: 222,
		Score:    5,
	}

	assert.Equal(t, uint64(999), review.ID)
	assert.Equal(t, uint64(111), review.OrderID)
	assert.Equal(t, uint64(222), review.PlayerID)
	assert.Equal(t, model.Rating(5), review.Score)
}

func TestReviewMultilingualContent(t *testing.T) {
	// 测试多语言内容
	reviews := []struct {
		content string
		lang    string
	}{
		{"Excellent service! Very professional.", "English"},
		{"服务非常好！非常专业。", "Chinese"},
		{"素晴らしいサービス！とてもプロフェッショナルです。", "Japanese"},
		{"Отличный сервис! Очень профессионально.", "Russian"},
		{"¡Excelente servicio! Muy profesional.", "Spanish"},
	}

	for _, rev := range reviews {
		review := &model.Review{
			OrderID: 100,
			UserID:  200,
			Score:   5,
			Content: rev.content,
		}
		assert.Equal(t, rev.content, review.Content, "Failed for language: %s", rev.lang)
	}
}

func TestReviewHTMLContent(t *testing.T) {
	// 测试包含HTML标签的内容
	review := &model.Review{
		OrderID:  100,
		UserID:   200,
		PlayerID: 300,
		Score:    4,
		Content:  "<b>很棒的服务</b><br><i>推荐给大家</i><br><a href='#'>查看更多</a>",
	}

	assert.Equal(t, "<b>很棒的服务</b><br><i>推荐给大家</i><br><a href='#'>查看更多</a>", review.Content)
}

func TestReviewStatus(t *testing.T) {
	tests := []struct {
		name   string
		status model.ReviewStatus
		valid  bool
	}{
		{"Pending status", model.ReviewStatusPending, true},
		{"Approved status", model.ReviewStatusApproved, true},
		{"Rejected status", model.ReviewStatusRejected, true},
		{"Deleted status", model.ReviewStatusDeleted, true},
		{"Invalid status", model.ReviewStatus("invalid"), false},
		{"Empty status", model.ReviewStatus(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.valid, tt.status.Valid())
		})
	}
}

func TestReviewWithStatus(t *testing.T) {
	review := &model.Review{
		Base: model.Base{
			ID: 1,
		},
		OrderID:  100,
		UserID:   200,
		PlayerID: 300,
		Score:    5,
		Content:  "Great service!",
		Status:   model.ReviewStatusPending,
	}

	assert.Equal(t, model.ReviewStatusPending, review.Status)

	// Test status transitions
	review.Status = model.ReviewStatusApproved
	assert.Equal(t, model.ReviewStatusApproved, review.Status)

	review.Status = model.ReviewStatusRejected
	assert.Equal(t, model.ReviewStatusRejected, review.Status)

	review.Status = model.ReviewStatusDeleted
	assert.Equal(t, model.ReviewStatusDeleted, review.Status)
}

func TestReviewIsReported(t *testing.T) {
	review := &model.Review{
		Base: model.Base{
			ID: 1,
		},
		OrderID:    100,
		UserID:     200,
		PlayerID:   300,
		Score:      5,
		Content:    "Good service",
		IsReported: false,
	}

	assert.False(t, review.IsReported)

	review.IsReported = true
	assert.True(t, review.IsReported)
}

func TestReviewImages(t *testing.T) {
	images := model.StringArray{
		"https://example.com/image1.jpg",
		"https://example.com/image2.jpg",
		"https://example.com/image3.jpg",
	}

	review := &model.Review{
		Base: model.Base{
			ID: 1,
		},
		OrderID:  100,
		UserID:   200,
		PlayerID: 300,
		Score:    5,
		Content:  "Great service with photos!",
		Images:   images,
	}

	assert.Len(t, review.Images, 3)
	assert.Equal(t, "https://example.com/image1.jpg", review.Images[0])
	assert.Equal(t, "https://example.com/image2.jpg", review.Images[1])
	assert.Equal(t, "https://example.com/image3.jpg", review.Images[2])
}

func TestReviewEmptyImages(t *testing.T) {
	review := &model.Review{
		Base: model.Base{
			ID: 1,
		},
		OrderID:  100,
		UserID:   200,
		PlayerID: 300,
		Score:    5,
		Content:  "Good service without photos",
		Images:   model.StringArray{},
	}

	assert.Len(t, review.Images, 0)
	assert.NotNil(t, review.Images)
}

func TestReviewRejectionReason(t *testing.T) {
	review := &model.Review{
		Base: model.Base{
			ID: 1,
		},
		OrderID:         100,
		UserID:          200,
		PlayerID:        300,
		Score:           3,
		Content:         "Contains inappropriate content",
		Status:          model.ReviewStatusRejected,
		RejectionReason: "评价内容包含敏感词",
	}

	assert.Equal(t, model.ReviewStatusRejected, review.Status)
	assert.Equal(t, "评价内容包含敏感词", review.RejectionReason)
}

func TestReviewJSONSerializationWithNewFields(t *testing.T) {
	review := &model.Review{
		Base: model.Base{
			ID: 1,
		},
		OrderID:  100,
		UserID:   200,
		PlayerID: 300,
		Score:    4,
		Content:  "Good service!",
		Status:   model.ReviewStatusApproved,
		IsReported: false,
		Images: model.StringArray{
			"https://example.com/image1.jpg",
			"https://example.com/image2.jpg",
		},
		RejectionReason: "",
	}

	// 序列化
	data, err := json.Marshal(review)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "Good service")
	assert.Contains(t, string(data), "approved")
	assert.Contains(t, string(data), "image1.jpg")

	// 反序列化
	var decoded model.Review
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, review.ID, decoded.ID)
	assert.Equal(t, review.OrderID, decoded.OrderID)
	assert.Equal(t, review.Status, decoded.Status)
	assert.Equal(t, review.IsReported, decoded.IsReported)
	assert.Len(t, decoded.Images, 2)
	assert.Equal(t, review.Images[0], decoded.Images[0])
}

func TestStringArrayScanAndValue(t *testing.T) {
	// Test Scan with valid JSON bytes
	var sa model.StringArray
	err := sa.Scan([]byte(`["image1.jpg","image2.jpg"]`))
	assert.NoError(t, err)
	assert.Len(t, sa, 2)
	assert.Equal(t, "image1.jpg", sa[0])
	assert.Equal(t, "image2.jpg", sa[1])

	// Test Scan with string
	var sa6 model.StringArray
	err = sa6.Scan(`["image3.jpg","image4.jpg"]`)
	assert.NoError(t, err)
	assert.Len(t, sa6, 2)
	assert.Equal(t, "image3.jpg", sa6[0])
	assert.Equal(t, "image4.jpg", sa6[1])

	// Test Scan with nil
	var sa2 model.StringArray
	err = sa2.Scan(nil)
	assert.NoError(t, err)
	assert.Len(t, sa2, 0)

	// Test Scan with empty bytes
	var sa3 model.StringArray
	err = sa3.Scan([]byte{})
	assert.NoError(t, err)
	assert.Len(t, sa3, 0)

	// Test Value with data
	sa4 := model.StringArray{"image1.jpg", "image2.jpg"}
	value, err := sa4.Value()
	assert.NoError(t, err)
	assert.NotNil(t, value)

	// Test Value with empty array
	sa5 := model.StringArray{}
	value2, err := sa5.Value()
	assert.NoError(t, err)
	assert.Equal(t, "[]", value2)
}

func TestReviewCompleteWorkflow(t *testing.T) {
	// Create a pending review
	review := &model.Review{
		Base: model.Base{
			ID: 1,
		},
		OrderID:    100,
		UserID:     200,
		PlayerID:   300,
		Score:      5,
		Content:    "Excellent service!",
		Status:     model.ReviewStatusPending,
		IsReported: false,
		Images: model.StringArray{
			"https://example.com/proof1.jpg",
		},
	}

	assert.Equal(t, model.ReviewStatusPending, review.Status)
	assert.False(t, review.IsReported)

	// Approve the review
	review.Status = model.ReviewStatusApproved
	assert.Equal(t, model.ReviewStatusApproved, review.Status)

	// Report the review
	review.IsReported = true
	assert.True(t, review.IsReported)

	// Delete the review
	review.Status = model.ReviewStatusDeleted
	assert.Equal(t, model.ReviewStatusDeleted, review.Status)
}

func TestReviewRejectionWorkflow(t *testing.T) {
	// Create a review with inappropriate content
	review := &model.Review{
		Base: model.Base{
			ID: 1,
		},
		OrderID:  100,
		UserID:   200,
		PlayerID: 300,
		Score:    2,
		Content:  "Contains bad words",
		Status:   model.ReviewStatusPending,
	}

	assert.Equal(t, model.ReviewStatusPending, review.Status)
	assert.Equal(t, "", review.RejectionReason)

	// Reject the review with reason
	review.Status = model.ReviewStatusRejected
	review.RejectionReason = "评价内容包含敏感词汇"

	assert.Equal(t, model.ReviewStatusRejected, review.Status)
	assert.Equal(t, "评价内容包含敏感词汇", review.RejectionReason)
}
