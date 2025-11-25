package model_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gamelink/internal/model"
)

func TestGameModel(t *testing.T) {
	game := &model.Game{
		Base: model.Base{
			ID: 1,
		},
		Key:         "lol",
		Name:        "League of Legends",
		Category:    "moba",
		IconURL:     "https://example.com/lol-icon.png",
		Description: "英雄联盟是一款多人在线战术竞技游戏",
	}

	assert.Equal(t, uint64(1), game.ID)
	assert.Equal(t, "lol", game.Key)
	assert.Equal(t, "League of Legends", game.Name)
	assert.Equal(t, "moba", game.Category)
	assert.Equal(t, "https://example.com/lol-icon.png", game.IconURL)
	assert.Equal(t, "英雄联盟是一款多人在线战术竞技游戏", game.Description)
}

func TestGameJSONSerialization(t *testing.T) {
	game := &model.Game{
		Base: model.Base{
			ID: 1,
		},
		Key:         "dota2",
		Name:        "Dota 2",
		Category:    "moba",
		IconURL:     "https://example.com/dota2-icon.png",
		Description: "Dota 2是一款多人在线战术竞技游戏",
	}

	// 序列化
	data, err := json.Marshal(game)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "dota2")
	assert.Contains(t, string(data), "Dota 2")

	// 反序列化
	var decoded model.Game
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, game.ID, decoded.ID)
	assert.Equal(t, game.Key, decoded.Key)
	assert.Equal(t, game.Name, decoded.Name)
	assert.Equal(t, game.Category, decoded.Category)
	assert.Equal(t, game.IconURL, decoded.IconURL)
	assert.Equal(t, game.Description, decoded.Description)
}

func TestGameZeroValues(t *testing.T) {
	game := &model.Game{
		Key:         "",
		Name:        "",
		Category:    "",
		IconURL:     "",
		Description: "",
	}

	assert.Equal(t, "", game.Key)
	assert.Equal(t, "", game.Name)
	assert.Equal(t, "", game.Category)
	assert.Equal(t, "", game.IconURL)
	assert.Equal(t, "", game.Description)
}

func TestGameEdgeCases(t *testing.T) {
	// 测试长文本
	longKey := "very-long-game-key-with-multiple-words-for-testing-purposes"
	longName := "这是一个非常长的游戏名称，用于测试字符串长度的边界情况，可能包含很多描述性信息"
	longCategory := "这是一个非常长的游戏分类名称，可能包含多级分类信息"
	longIconURL := "https://example.com/very/long/icon/url/path/that/might/be/used/for/testing/purposes/and/should/handle/long/paths/without/issues.png"
	longDescription := "这是一个非常长的游戏描述，可以包含很多详细信息，比如游戏的具体玩法、特色、适用人群、系统要求、发布日期、开发商信息等等。这种长文本测试可以确保我们的模型能够处理各种长度的输入。"

	game1 := &model.Game{
		Key:         longKey,
		Name:        longName,
		Category:    longCategory,
		IconURL:     longIconURL,
		Description: longDescription,
	}

	assert.Equal(t, longKey, game1.Key)
	assert.Equal(t, longName, game1.Name)
	assert.Equal(t, longCategory, game1.Category)
	assert.Equal(t, longIconURL, game1.IconURL)
	assert.Equal(t, longDescription, game1.Description)

	// 测试特殊字符
	game2 := &model.Game{
		Key:         "game@123#test",
		Name:        "游戏@测试#123",
		Category:    "分类@#$%^&*()",
		IconURL:     "https://example.com/icon@special#chars.png",
		Description: "描述包含特殊字符：<>{}[]|\\\"quotes\" and 'apostrophes' and @#$%^&*()_+-=[]{}|;':\",./<>?😊🚀",
	}
	assert.Equal(t, "game@123#test", game2.Key)
	assert.Equal(t, "游戏@测试#123", game2.Name)
	assert.Equal(t, "分类@#$%^&*()", game2.Category)
	assert.Equal(t, "https://example.com/icon@special#chars.png", game2.IconURL)
	assert.Contains(t, game2.Description, "特殊字符")

	// 测试常见游戏分类
	categories := []string{"moba", "fps", "rpg", "rts", "sports", "racing", "fighting", "puzzle", "strategy", "simulation"}
	for _, category := range categories {
		game := &model.Game{
			Category: category,
		}
		assert.Equal(t, category, game.Category)
	}
}

func TestGameCommonGames(t *testing.T) {
	// 测试常见游戏
	games := []struct {
		key         string
		name        string
		category    string
		description string
	}{
		{"lol", "League of Legends", "moba", "英雄联盟是一款流行的多人在线战术竞技游戏"},
		{"dota2", "Dota 2", "moba", "Dota 2是一款经典的多人在线战术竞技游戏"},
		{"csgo", "Counter-Strike: Global Offensive", "fps", "CS:GO是一款热门的第一人称射击游戏"},
		{"valorant", "Valorant", "fps", "Valorant是一款战术射击游戏"},
		{"overwatch", "Overwatch", "fps", "守望先锋是一款团队射击游戏"},
		{"apex", "Apex Legends", "fps", "Apex英雄是一款大逃杀射击游戏"},
		{"pubg", "PUBG: Battlegrounds", "fps", "绝地求生是一款大逃杀游戏"},
		{"fortnite", "Fortnite", "fps", "堡垒之夜是一款建造射击游戏"},
		{"wow", "World of Warcraft", "rpg", "魔兽世界是一款经典的大型多人在线角色扮演游戏"},
		{"ffxiv", "Final Fantasy XIV", "rpg", "最终幻想14是一款热门的大型多人在线角色扮演游戏"},
	}

	for _, gameData := range games {
		game := &model.Game{
			Key:         gameData.key,
			Name:        gameData.name,
			Category:    gameData.category,
			Description: gameData.description,
		}
		assert.Equal(t, gameData.key, game.Key)
		assert.Equal(t, gameData.name, game.Name)
		assert.Equal(t, gameData.category, game.Category)
		assert.Equal(t, gameData.description, game.Description)
	}
}

func TestGameWithBaseFields(t *testing.T) {
	game := &model.Game{
		Base: model.Base{
			ID: 123,
		},
		Key:  "test-game",
		Name: "测试游戏",
	}

	assert.Equal(t, uint64(123), game.ID)
	assert.Equal(t, "test-game", game.Key)
	assert.Equal(t, "测试游戏", game.Name)
}

func TestGameJSONFields(t *testing.T) {
	game := &model.Game{
		Base: model.Base{
			ID: 1,
		},
		Key:         "lol",
		Name:        "League of Legends",
		Category:    "moba",
		IconURL:     "https://example.com/lol-icon.png",
		Description: "英雄联盟游戏",
	}

	data, err := json.Marshal(game)
	assert.NoError(t, err)

	// 验证JSON结构
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// 检查必需的字段
	assert.Contains(t, result, "id")
	assert.Contains(t, result, "key")
	assert.Contains(t, result, "name")
	assert.Contains(t, result, "category")
	assert.Contains(t, result, "iconUrl")
	assert.Contains(t, result, "description")

	// 验证值
	assert.Equal(t, float64(1), result["id"])
	assert.Equal(t, "lol", result["key"])
	assert.Equal(t, "League of Legends", result["name"])
	assert.Equal(t, "moba", result["category"])
	assert.Equal(t, "https://example.com/lol-icon.png", result["iconUrl"])
	assert.Equal(t, "英雄联盟游戏", result["description"])
}

func TestGameEmptyFields(t *testing.T) {
	game := &model.Game{}

	assert.Equal(t, uint64(0), game.ID)
	assert.Equal(t, "", game.Key)
	assert.Equal(t, "", game.Name)
	assert.Equal(t, "", game.Category)
	assert.Equal(t, "", game.IconURL)
	assert.Equal(t, "", game.Description)
}

func TestGameMultilingualContent(t *testing.T) {
	// 测试多语言内容
	games := []struct {
		key         string
		name        string
		description string
		lang        string
	}{
		{"honorofkings", "Honor of Kings", "A popular mobile MOBA game", "English"},
		{"王者荣耀", "王者荣耀", "一款流行的手机MOBA游戏", "Chinese"},
		{"honor-des-rois", "Honor of Kings", "Un jeu MOBA mobile populaire", "French"},
		{"honor-de-reyes", "Honor of Kings", "Un popular juego MOBA móvil", "Spanish"},
		{"honra-dos-reis", "Honor of Kings", "Um popular jogo MOBA móvel", "Portuguese"},
	}

	for _, gameData := range games {
		game := &model.Game{
			Key:         gameData.key,
			Name:        gameData.name,
			Description: gameData.description,
		}
		assert.Equal(t, gameData.key, game.Key, "Failed for language: %s", gameData.lang)
		assert.Equal(t, gameData.name, game.Name, "Failed for language: %s", gameData.lang)
		assert.Equal(t, gameData.description, game.Description, "Failed for language: %s", gameData.lang)
	}
}

func TestGameSpecialCharacters(t *testing.T) {
	// 测试包含特殊字符的内容
	game := &model.Game{
		Key:         "game-with-dashes_and_underscores.123",
		Name:        "游戏-名称_测试.123 | Special Edition",
		Category:    "action-rpg_fantasy",
		IconURL:     "https://example.com/game-icon@special#chars-v2.0.png",
		Description: "游戏描述包含特殊字符：<>{}[]|\\\"quotes\" and 'apostrophes' and @#$%^&*()_+-=[]{}|;':\",./<>?😊🚀",
	}

	assert.Equal(t, "game-with-dashes_and_underscores.123", game.Key)
	assert.Equal(t, "游戏-名称_测试.123 | Special Edition", game.Name)
	assert.Equal(t, "action-rpg_fantasy", game.Category)
	assert.Equal(t, "https://example.com/game-icon@special#chars-v2.0.png", game.IconURL)
	assert.Contains(t, game.Description, "特殊字符")
}

func TestGameOptionalFields(t *testing.T) {
	// 测试可选字段为空的情况
	game1 := &model.Game{
		Key:         "test-game",
		Name:        "测试游戏",
		Category:    "", // 空分类
		IconURL:     "", // 空图标URL
		Description: "", // 空描述
	}

	assert.Equal(t, "", game1.Category)
	assert.Equal(t, "", game1.IconURL)
	assert.Equal(t, "", game1.Description)

	// 测试只有必需字段的情况
	game2 := &model.Game{
		Key:  "minimal-game",
		Name: "Minimal Game",
	}

	assert.Equal(t, "minimal-game", game2.Key)
	assert.Equal(t, "Minimal Game", game2.Name)
	assert.Equal(t, "", game2.Category)
	assert.Equal(t, "", game2.IconURL)
	assert.Equal(t, "", game2.Description)
}