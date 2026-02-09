package db

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
)

// seedGameCategories 创建游戏分类种子数据
// 用于替代 Game.Category 字段（已废弃），实现游戏的分类管理
func seedGameCategories(tx *gorm.DB, games map[string]*model.Game) (map[string]*model.GameCategory, error) {
	seeds := []struct {
		Name        string
		Description string
		IconURL     string
		SortOrder   int
		GameKeys    []string // 关联的游戏
	}{
		{
			Name:        "MOBA",
			Description: "多人在线战术竞技游戏，如英雄联盟、DOTA2、王者荣耀",
			IconURL:     "https://example.com/categories/moba.png",
			SortOrder:   1,
			GameKeys:    []string{"lol", "dota2", "wzry"},
		},
		{
			Name:        "FPS",
			Description: "第一人称射击游戏，如CS:GO、Valorant、Apex等",
			IconURL:     "https://example.com/categories/fps.png",
			SortOrder:   2,
			GameKeys:    []string{"valorant", "csgo", "apex", "pubg", "overwatch"},
		},
		{
			Name:        "RPG",
			Description: "角色扮演游戏，如原神、魔兽世界等",
			IconURL:     "https://example.com/categories/rpg.png",
			SortOrder:   3,
			GameKeys:    []string{"genshin", "wow"},
		},
		{
			Name:        "体育竞技",
			Description: "体育模拟类游戏，如FIFA、NBA2K等",
			IconURL:     "https://example.com/categories/sports.png",
			SortOrder:   4,
			GameKeys:    []string{"fifa", "nba2k"},
		},
		{
			Name:        "休闲派对",
			Description: "轻松休闲的多人派对游戏",
			IconURL:     "https://example.com/categories/party.png",
			SortOrder:   5,
			GameKeys:    []string{"fallguys", "amongus"},
		},
		{
			Name:        "沙盒建造",
			Description: "自由建造的沙盒类游戏",
			IconURL:     "https://example.com/categories/sandbox.png",
			SortOrder:   6,
			GameKeys:    []string{"minecraft"},
		},
	}

	result := make(map[string]*model.GameCategory, len(seeds))
	for _, seed := range seeds {
		var category model.GameCategory
		if err := tx.Where("name = ?", seed.Name).First(&category).Error; err == nil {
			result[seed.Name] = &category
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		} else {
			category = model.GameCategory{
				Name:        seed.Name,
				Description: seed.Description,
				IconURL:     seed.IconURL,
				SortOrder:   seed.SortOrder,
				IsActive:    true,
			}
			category.ExtJSON = `{"seed":"demo"}`
			if err := tx.Create(&category).Error; err != nil {
				return nil, err
			}
			result[seed.Name] = &category
		}

		// 更新游戏的 CategoryID
		for _, gameKey := range seed.GameKeys {
			if game, ok := games[gameKey]; ok {
				if game.CategoryID == nil || *game.CategoryID != category.ID {
					if err := tx.Model(&model.Game{}).Where("id = ?", game.ID).Update("category_id", category.ID).Error; err != nil {
						log.Printf("warning: failed to update game %s category: %v", gameKey, err)
					}
				}
			}
		}
	}

	log.Printf("game categories ensured: %d\n", len(result))
	return result, nil
}

// seedLFGRequests 创建快速匹配请求种子数据
// 覆盖 LFGRequest 模型的各种状态
func seedLFGRequests(tx *gorm.DB, users map[string]*model.User, games map[string]*model.Game) error {
	now := time.Now()

	// 从数据库获取已有的公开聊天群组（用于匹配成功的请求）
	var publicGroups []model.ChatGroup
	tx.Where("group_type = ?", "public").Find(&publicGroups)
	chatGroupMap := make(map[string]*model.ChatGroup)
	for i := range publicGroups {
		chatGroupMap[publicGroups[i].GroupName] = &publicGroups[i]
	}

	seeds := []struct {
		UserKey         string
		GameKey         string
		RequestType     model.LFGRequestType
		Title           string
		Description     string
		RequiredPlayers int
		MinRank         string
		MaxPriceCents   int64
		Status          model.LFGRequestStatus
		ExpiresOffset   time.Duration
		MatchedGroupKey string // 匹配成功后关联的聊天群组（按名称）
	}{
		// 待匹配的请求
		{
			UserKey:         "customerA",
			GameKey:         "lol",
			RequestType:     model.LFGFindPlayer,
			Title:           "求黄金陪玩上分",
			Description:     "需要一个打野位大神带我上分，黄金段位，希望能连胜几把",
			RequiredPlayers: 1,
			MinRank:         "黄金",
			MaxPriceCents:   15000,
			Status:          model.LFGPending,
			ExpiresOffset:   2 * time.Hour,
		},
		{
			UserKey:         "customerB",
			GameKey:         "valorant",
			RequestType:     model.LFGFindTeam,
			Title:           "找队友一起打排位",
			Description:     "铂金段位，找几个稳定队友一起冲分，最好有语音",
			RequiredPlayers: 4,
			MinRank:         "铂金",
			Status:          model.LFGPending,
			ExpiresOffset:   3 * time.Hour,
		},
		{
			UserKey:         "customerC",
			GameKey:         "dota2",
			RequestType:     model.LFGFindPlayer,
			Title:           "新手求带",
			Description:     "刚开始玩DOTA2，希望有大佬愿意教教我基础操作",
			RequiredPlayers: 1,
			MaxPriceCents:   10000,
			Status:          model.LFGPending,
			ExpiresOffset:   1 * time.Hour,
		},
		// 已匹配的请求
		{
			UserKey:         "customerD",
			GameKey:         "lol",
			RequestType:     model.LFGFindPlayer,
			Title:           "钻石陪练",
			Description:     "找一个钻石以上的陪玩师帮我练习中单",
			RequiredPlayers: 1,
			MinRank:         "钻石",
			MaxPriceCents:   20000,
			Status:          model.LFGMatched,
			ExpiresOffset:   -30 * time.Minute, // 已过期但已匹配
			MatchedGroupKey: "英雄联盟大厅",          // 使用实际群组名
		},
		// 已过期的请求
		{
			UserKey:         "customerE",
			GameKey:         "csgo",
			RequestType:     model.LFGFindTeam,
			Title:           "找队友打匹配",
			Description:     "想找几个人一起玩CS:GO，不要太菜就行",
			RequiredPlayers: 4,
			Status:          model.LFGExpired,
			ExpiresOffset:   -2 * time.Hour,
		},
		// 已取消的请求
		{
			UserKey:         "customerF",
			GameKey:         "apex",
			RequestType:     model.LFGFindPlayer,
			Title:           "双排大逃杀",
			Description:     "找一个队友一起吃鸡",
			RequiredPlayers: 1,
			Status:          model.LFGCanceled,
			ExpiresOffset:   -1 * time.Hour,
		},
	}

	for _, seed := range seeds {
		user, ok := users[seed.UserKey]
		if !ok {
			continue
		}
		game, ok := games[seed.GameKey]
		if !ok {
			continue
		}

		// 检查是否已存在（按 user_id + title 判断）
		var existing model.LFGRequest
		if err := tx.Where("user_id = ? AND title = ?", user.ID, seed.Title).First(&existing).Error; err == nil {
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		expiresAt := now.Add(seed.ExpiresOffset)
		request := model.LFGRequest{
			UserID:          user.ID,
			GameID:          game.ID,
			RequestType:     seed.RequestType,
			Title:           seed.Title,
			Description:     seed.Description,
			RequiredPlayers: seed.RequiredPlayers,
			MinRank:         seed.MinRank,
			MaxPriceCents:   seed.MaxPriceCents,
			Status:          seed.Status,
			ExpiresAt:       expiresAt,
		}
		request.ExtJSON = `{"seed":"demo"}`

		// 如果已匹配，设置匹配信息
		if seed.Status == model.LFGMatched && seed.MatchedGroupKey != "" {
			if group, ok := chatGroupMap[seed.MatchedGroupKey]; ok {
				request.MatchedRoomID = &group.ID
				matchedAt := now.Add(-15 * time.Minute)
				request.MatchedAt = &matchedAt
			}
		}

		if err := tx.Create(&request).Error; err != nil {
			return err
		}
	}

	log.Println("LFG requests seed data created successfully")
	return nil
}

// seedFavorites 创建用户收藏种子数据
// 覆盖 Favorite 模型
func seedFavorites(tx *gorm.DB, users map[string]*model.User, players map[string]*model.Player) error {
	seeds := []struct {
		UserKey   string
		PlayerKey string
	}{
		// customerA 收藏了多个陪玩师
		{"customerA", "playerA"},
		{"customerA", "playerB"},
		{"customerA", "playerC"},
		// customerB 收藏了 playerA
		{"customerB", "playerA"},
		{"customerB", "playerD"},
		// customerC 收藏了 playerB
		{"customerC", "playerB"},
		{"customerC", "playerE"},
		// customerD 收藏了 playerA 和 playerF
		{"customerD", "playerA"},
		{"customerD", "playerF"},
		// customerE 收藏了 playerC
		{"customerE", "playerC"},
		// customerF 收藏了多个陪玩师
		{"customerF", "playerA"},
		{"customerF", "playerB"},
		{"customerF", "playerD"},
		{"customerF", "playerE"},
		// customerG（新用户）收藏了 playerA
		{"customerG", "playerA"},
		// customerH 收藏了 playerF
		{"customerH", "playerF"},
	}

	for _, seed := range seeds {
		user, ok := users[seed.UserKey]
		if !ok {
			continue
		}
		player, ok := players[seed.PlayerKey]
		if !ok {
			continue
		}

		// 检查是否已存在（唯一索引）
		var existing model.Favorite
		if err := tx.Where("user_id = ? AND player_id = ?", user.ID, player.ID).First(&existing).Error; err == nil {
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		favorite := model.Favorite{
			UserID:   user.ID,
			PlayerID: player.ID,
		}
		favorite.ExtJSON = `{"seed":"demo"}`

		if err := tx.Create(&favorite).Error; err != nil {
			log.Printf("warning: failed to create favorite %s->%s: %v", seed.UserKey, seed.PlayerKey, err)
			continue
		}
	}

	log.Println("favorites seed data created successfully")
	return nil
}

// seedPlayerPresences 创建陪玩师在线状态种子数据
// 覆盖 PlayerPresence 模型
func seedPlayerPresences(tx *gorm.DB, players map[string]*model.Player, games map[string]*model.Game, orders map[string]*model.Order) error {
	now := time.Now()

	seeds := []struct {
		PlayerKey      string
		Status         model.PlayerPresenceStatus
		CurrentGameKey string
		CustomStatus   string
		DeviceType     string
		OrderKey       string // 关联进行中的订单
		HeartbeatDelta time.Duration
	}{
		// 在线空闲
		{
			PlayerKey:      "playerA",
			Status:         model.PresenceOnline,
			CurrentGameKey: "",
			CustomStatus:   "随时开车，速来！",
			DeviceType:     "web",
			HeartbeatDelta: -1 * time.Minute,
		},
		// 接单中
		{
			PlayerKey:      "playerB",
			Status:         model.PresenceAccepting,
			CurrentGameKey: "valorant",
			CustomStatus:   "FPS专业陪练中",
			DeviceType:     "desktop",
			HeartbeatDelta: -30 * time.Second,
		},
		// 游戏中
		{
			PlayerKey:      "playerC",
			Status:         model.PresenceInGame,
			CurrentGameKey: "csgo",
			CustomStatus:   "",
			DeviceType:     "desktop",
			OrderKey:       "orderInProgress1",
			HeartbeatDelta: -2 * time.Minute,
		},
		// 匹配中
		{
			PlayerKey:      "playerD",
			Status:         model.PresenceMatching,
			CurrentGameKey: "wow",
			CustomStatus:   "副本匹配中...",
			DeviceType:     "web",
			HeartbeatDelta: -45 * time.Second,
		},
		// 休息中
		{
			PlayerKey:      "playerE",
			Status:         model.PresenceResting,
			CurrentGameKey: "",
			CustomStatus:   "午休中，14点后接单",
			DeviceType:     "mobile",
			HeartbeatDelta: -10 * time.Minute,
		},
		// 离线
		{
			PlayerKey:      "playerF",
			Status:         model.PresenceOffline,
			CurrentGameKey: "",
			CustomStatus:   "",
			DeviceType:     "",
			HeartbeatDelta: -2 * time.Hour,
		},
		// 隐身
		{
			PlayerKey:      "playerG",
			Status:         model.PresenceInvisible,
			CurrentGameKey: "dota2",
			CustomStatus:   "私密模式",
			DeviceType:     "web",
			HeartbeatDelta: -5 * time.Minute,
		},
	}

	for _, seed := range seeds {
		player, ok := players[seed.PlayerKey]
		if !ok {
			continue
		}

		// 检查是否已存在（按 player_id 唯一索引）
		var existing model.PlayerPresence
		if err := tx.Where("player_id = ?", player.ID).First(&existing).Error; err == nil {
			// 更新现有记录
			updates := map[string]interface{}{
				"status":            seed.Status,
				"custom_status":     seed.CustomStatus,
				"device_type":       seed.DeviceType,
				"last_heartbeat_at": now.Add(seed.HeartbeatDelta),
			}
			if seed.CurrentGameKey != "" {
				if game, ok := games[seed.CurrentGameKey]; ok {
					updates["current_game_id"] = game.ID
					updates["current_game_name"] = game.Name
				}
			} else {
				updates["current_game_id"] = nil
				updates["current_game_name"] = ""
			}
			if seed.OrderKey != "" {
				if order, ok := orders[seed.OrderKey]; ok {
					updates["current_order_id"] = order.ID
				}
			}
			if err := tx.Model(&existing).Updates(updates).Error; err != nil {
				log.Printf("warning: failed to update presence for %s: %v", seed.PlayerKey, err)
			}
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		presence := model.PlayerPresence{
			PlayerID:        player.ID,
			Status:          seed.Status,
			CustomStatus:    seed.CustomStatus,
			DeviceType:      seed.DeviceType,
			LastHeartbeatAt: now.Add(seed.HeartbeatDelta),
		}
		presence.ExtJSON = `{"seed":"demo"}`

		if seed.CurrentGameKey != "" {
			if game, ok := games[seed.CurrentGameKey]; ok {
				presence.CurrentGameID = &game.ID
				presence.CurrentGameName = game.Name
			}
		}

		if seed.OrderKey != "" {
			if order, ok := orders[seed.OrderKey]; ok {
				presence.CurrentOrderID = &order.ID
			}
		}

		if err := tx.Create(&presence).Error; err != nil {
			return err
		}
	}

	log.Println("player presences seed data created successfully")
	return nil
}
