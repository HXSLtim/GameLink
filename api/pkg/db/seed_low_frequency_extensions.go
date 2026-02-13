package db

import (
	"errors"

	"gorm.io/gorm"

	"gamelink/internal/model"
)

func seedPlayerGames(tx *gorm.DB, players map[string]*model.Player, games map[string]*model.Game) error {
	specs := map[string][]string{
		"playerA": {"lol", "dota2", "wzry"},
		"playerB": {"valorant", "apex", "pubg"},
		"playerC": {"csgo", "valorant"},
		"playerD": {"wow", "genshin", "minecraft"},
		"playerE": {"fifa", "nba2k"},
		"playerF": {"fallguys", "amongus"},
		"playerG": {"dota2", "lol"},
		"playerH": {"pubg", "apex"},
	}

	for playerKey, gameKeys := range specs {
		player := players[playerKey]
		if player == nil {
			continue
		}

		mainGameID := player.MainGameID
		if mainGameID == 0 && len(gameKeys) > 0 {
			if mainGame := games[gameKeys[0]]; mainGame != nil {
				mainGameID = mainGame.ID
			}
		}

		for _, gameKey := range gameKeys {
			game := games[gameKey]
			if game == nil {
				continue
			}

			var existing model.PlayerGame
			err := tx.Where("player_id = ? AND game_id = ?", player.ID, game.ID).First(&existing).Error
			isMain := game.ID == mainGameID
			if err == nil {
				if err := tx.Model(&model.PlayerGame{}).Where("id = ?", existing.ID).
					Update("is_main", isMain).Error; err != nil {
					return err
				}
				continue
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			row := &model.PlayerGame{
				PlayerID: player.ID,
				GameID:   game.ID,
				IsMain:   isMain,
			}
			if err := tx.Create(row).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func seedPlayerSkillTags(tx *gorm.DB, players map[string]*model.Player) error {
	specs := map[string][]string{
		"playerA": {"打野", "指挥", "上分"},
		"playerB": {"狙击", "战术沟通", "突破"},
		"playerC": {"枪法", "道具配合", "残局"},
		"playerD": {"副本带练", "机制讲解"},
		"playerE": {"体育竞技", "阵型教学"},
		"playerF": {"娱乐氛围", "休闲陪玩"},
		"playerG": {"中单", "运营", "教学"},
		"playerH": {"生存", "报点"},
	}

	for playerKey, tags := range specs {
		player := players[playerKey]
		if player == nil {
			continue
		}

		for _, tag := range tags {
			var existing model.PlayerSkillTag
			err := tx.Where("player_id = ? AND tag = ?", player.ID, tag).First(&existing).Error
			if err == nil {
				continue
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if err := tx.Create(&model.PlayerSkillTag{PlayerID: player.ID, Tag: tag}).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func seedReviewDisplaySettings(tx *gorm.DB) error {
	defaultSettings := model.DefaultReviewDisplaySettings()

	var existing model.ReviewDisplaySettings
	err := tx.Where("id = ?", defaultSettings.ID).First(&existing).Error
	if err == nil {
		updates := map[string]interface{}{
			"sort_by":                 defaultSettings.SortBy,
			"min_score":               defaultSettings.MinScore,
			"show_anonymous":          defaultSettings.ShowAnonymous,
			"page_size":               defaultSettings.PageSize,
			"auto_approve":            defaultSettings.AutoApprove,
			"auto_approve_min_rating": defaultSettings.AutoApproveMinRating,
		}
		return tx.Model(&model.ReviewDisplaySettings{}).Where("id = ?", defaultSettings.ID).Updates(updates).Error
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return tx.Create(defaultSettings).Error
}

func seedRankingRewards(tx *gorm.DB) error {
	specs := []model.RankingReward{
		{
			RankingType: model.RankingTypeIncome,
			Period:      "monthly",
			RankStart:   1,
			RankEnd:     3,
			RewardType:  "fixed",
			RewardValue: 100000,
			Description: "月度收入榜 TOP3 固定奖金",
			IsActive:    true,
		},
		{
			RankingType: model.RankingTypeIncome,
			Period:      "monthly",
			RankStart:   4,
			RankEnd:     10,
			RewardType:  "fixed",
			RewardValue: 30000,
			Description: "月度收入榜 4-10 名固定奖金",
			IsActive:    true,
		},
		{
			RankingType: model.RankingTypeOrderCount,
			Period:      "monthly",
			RankStart:   1,
			RankEnd:     5,
			RewardType:  "percentage",
			RewardValue: 5,
			Description: "月度订单榜 TOP5 手续费减免 5%",
			IsActive:    true,
		},
	}

	for _, spec := range specs {
		var existing model.RankingReward
		err := tx.Where("ranking_type = ? AND period = ? AND rank_start = ? AND rank_end = ?",
			spec.RankingType, spec.Period, spec.RankStart, spec.RankEnd).First(&existing).Error
		if err == nil {
			updates := map[string]interface{}{
				"reward_type":  spec.RewardType,
				"reward_value": spec.RewardValue,
				"description":  spec.Description,
				"is_active":    spec.IsActive,
			}
			if err := tx.Model(&model.RankingReward{}).Where("id = ?", existing.ID).Updates(updates).Error; err != nil {
				return err
			}
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		row := spec
		if err := tx.Create(&row).Error; err != nil {
			return err
		}
	}

	return nil
}
