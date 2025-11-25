package model_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gamelink/internal/model"
)

func TestPlayerModel(t *testing.T) {
	player := &model.Player{
		Base: model.Base{
			ID: 1,
		},
		UserID:             100,
		Nickname:           "TestPlayer",
		Bio:                "This is a test player bio",
		Rank:               "Diamond",
		RatingAverage:      4.5,
		RatingCount:        50,
		HourlyRateCents:    10000, // 100元/小时
		MainGameID:         10,
		VerificationStatus: model.VerificationVerified,
	}

	assert.Equal(t, uint64(1), player.ID)
	assert.Equal(t, uint64(100), player.UserID)
	assert.Equal(t, "TestPlayer", player.Nickname)
	assert.Equal(t, "This is a test player bio", player.Bio)
	assert.Equal(t, "Diamond", player.Rank)
	assert.Equal(t, float32(4.5), player.RatingAverage)
	assert.Equal(t, uint32(50), player.RatingCount)
	assert.Equal(t, int64(10000), player.HourlyRateCents)
	assert.Equal(t, uint64(10), player.MainGameID)
	assert.Equal(t, model.VerificationVerified, player.VerificationStatus)
}

func TestPlayerJSONSerialization(t *testing.T) {
	player := &model.Player{
		Base: model.Base{
			ID: 1,
		},
		UserID:             100,
		Nickname:           "TestPlayer",
		Bio:                "This is a test player bio",
		Rank:               "Diamond",
		RatingAverage:      4.5,
		RatingCount:        50,
		HourlyRateCents:    10000,
		MainGameID:         10,
		VerificationStatus: model.VerificationVerified,
	}

	// 序列化
	data, err := json.Marshal(player)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "TestPlayer")
	assert.Contains(t, string(data), "Diamond")

	// 反序列化
	var decoded model.Player
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, player.ID, decoded.ID)
	assert.Equal(t, player.UserID, decoded.UserID)
	assert.Equal(t, player.Nickname, decoded.Nickname)
	assert.Equal(t, player.RatingAverage, decoded.RatingAverage)
	assert.Equal(t, player.VerificationStatus, decoded.VerificationStatus)
}

func TestPlayerConstants(t *testing.T) {
	// 测试验证状态常量
	assert.Equal(t, model.VerificationStatus("pending"), model.VerificationPending)
	assert.Equal(t, model.VerificationStatus("verified"), model.VerificationVerified)
	assert.Equal(t, model.VerificationStatus("rejected"), model.VerificationRejected)
}

func TestPlayerRatingValidation(t *testing.T) {
	// 测试有效的评分范围
	player1 := &model.Player{
		RatingAverage: 0.0,
	}
	assert.Equal(t, float32(0.0), player1.RatingAverage)

	player2 := &model.Player{
		RatingAverage: 5.0,
	}
	assert.Equal(t, float32(5.0), player2.RatingAverage)

	player3 := &model.Player{
		RatingAverage: 2.5,
	}
	assert.Equal(t, float32(2.5), player3.RatingAverage)

	// 测试超出范围的评分（虽然代码中没有强制验证，但测试边界值）
	player4 := &model.Player{
		RatingAverage: 6.0, // 超出5分制
	}
	assert.Equal(t, float32(6.0), player4.RatingAverage)

	player5 := &model.Player{
		RatingAverage: -1.0, // 负值
	}
	assert.Equal(t, float32(-1.0), player5.RatingAverage)
}

func TestPlayerZeroValues(t *testing.T) {
	player := &model.Player{
		UserID:          0,
		Nickname:        "",
		Bio:             "",
		Rank:            "",
		RatingAverage:   0.0,
		RatingCount:     0,
		HourlyRateCents: 0,
		MainGameID:      0,
	}

	assert.Equal(t, uint64(0), player.UserID)
	assert.Equal(t, "", player.Nickname)
	assert.Equal(t, "", player.Bio)
	assert.Equal(t, "", player.Rank)
	assert.Equal(t, float32(0.0), player.RatingAverage)
	assert.Equal(t, uint32(0), player.RatingCount)
	assert.Equal(t, int64(0), player.HourlyRateCents)
	assert.Equal(t, uint64(0), player.MainGameID)
}

func TestPlayerVerificationStatuses(t *testing.T) {
	// 测试所有验证状态
	statuses := []model.VerificationStatus{
		model.VerificationPending,
		model.VerificationVerified,
		model.VerificationRejected,
	}

	for _, status := range statuses {
		player := &model.Player{
			VerificationStatus: status,
		}
		assert.Equal(t, status, player.VerificationStatus)
	}
}

func TestPlayerEdgeCases(t *testing.T) {
	// 测试长文本
	longBio := "This is a very long bio that could potentially contain a lot of information about the player, their gaming experience, skills, achievements, and other relevant details that users might want to know."
	player1 := &model.Player{
		Bio: longBio,
	}
	assert.Equal(t, longBio, player1.Bio)

	// 测试特殊字符
	player2 := &model.Player{
		Nickname: "玩家_123#测试",
		Bio:      "Bio with special chars: @#$%^&*()",
		Rank:     "Rank#1",
	}
	assert.Equal(t, "玩家_123#测试", player2.Nickname)
	assert.Equal(t, "Bio with special chars: @#$%^&*()", player2.Bio)
	assert.Equal(t, "Rank#1", player2.Rank)

	// 测试高评分数量
	player3 := &model.Player{
		RatingCount: ^uint32(0), // 最大uint32值
	}
	assert.Equal(t, ^uint32(0), player3.RatingCount)

	// 测试高费率
	player4 := &model.Player{
		HourlyRateCents: ^int64(0), // 最大int64值
	}
	assert.Equal(t, ^int64(0), player4.HourlyRateCents)
}

func TestPlayerFloatPrecision(t *testing.T) {
	// 测试浮点数精度
	player1 := &model.Player{
		RatingAverage: 4.333333,
	}
	assert.Equal(t, float32(4.333333), player1.RatingAverage)

	player2 := &model.Player{
		RatingAverage: 4.7,
	}
	assert.Equal(t, float32(4.7), player2.RatingAverage)
}

func TestPlayerEmptyFields(t *testing.T) {
	player := &model.Player{}

	assert.Equal(t, uint64(0), player.ID)
	assert.Equal(t, uint64(0), player.UserID)
	assert.Equal(t, "", player.Nickname)
	assert.Equal(t, "", player.Bio)
	assert.Equal(t, "", player.Rank)
	assert.Equal(t, float32(0.0), player.RatingAverage)
	assert.Equal(t, uint32(0), player.RatingCount)
	assert.Equal(t, int64(0), player.HourlyRateCents)
	assert.Equal(t, uint64(0), player.MainGameID)
	assert.Equal(t, model.VerificationStatus(""), player.VerificationStatus)
}

func TestPlayerWithBaseFields(t *testing.T) {
	player := &model.Player{
		Base: model.Base{
			ID: 12345,
		},
		Nickname: "TestPlayer",
	}

	assert.Equal(t, uint64(12345), player.ID)
	assert.Equal(t, "TestPlayer", player.Nickname)
}

func TestPlayerJSONFields(t *testing.T) {
	player := &model.Player{
		UserID:          100,
		Nickname:        "TestPlayer",
		Bio:             "Test bio",
		Rank:            "Diamond",
		RatingAverage:   4.5,
		RatingCount:     25,
		HourlyRateCents: 15000,
		MainGameID:      5,
	}

	data, err := json.Marshal(player)
	assert.NoError(t, err)

	// 验证JSON结构
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// 检查必需的字段
	assert.Contains(t, result, "userId")
	assert.Contains(t, result, "nickname")
	assert.Contains(t, result, "bio")
	assert.Contains(t, result, "rank")
	assert.Contains(t, result, "ratingAverage")
	assert.Contains(t, result, "ratingCount")
	assert.Contains(t, result, "hourlyRateCents")
	assert.Contains(t, result, "mainGameId")

	// 验证值
	assert.Equal(t, float64(100), result["userId"])
	assert.Equal(t, "TestPlayer", result["nickname"])
	assert.Equal(t, float64(4.5), result["ratingAverage"])
	assert.Equal(t, float64(25), result["ratingCount"])
	assert.Equal(t, float64(15000), result["hourlyRateCents"])
	assert.Equal(t, float64(5), result["mainGameId"])
}