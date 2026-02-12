package db

import (
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"gamelink/internal/model"
)

func seedPlayerServices(tx *gorm.DB, players map[string]*model.Player, games map[string]*model.Game) error {
	type serviceSpec struct {
		PlayerKey    string
		GameKey      string
		RankLevel    int
		Description  string
		IsActive     bool
	}

	specs := []serviceSpec{
		{PlayerKey: "playerA", GameKey: "lol", RankLevel: 7, Description: "LOL 王者段位护航，擅长打野与节奏带动。", IsActive: true},
		{PlayerKey: "playerA", GameKey: "dota2", RankLevel: 6, Description: "DOTA2 大师局双排，强调团队沟通与地图控制。", IsActive: true},
		{PlayerKey: "playerB", GameKey: "valorant", RankLevel: 6, Description: "无畏契约高分局枪法与战术配合训练。", IsActive: true},
		{PlayerKey: "playerC", GameKey: "csgo", RankLevel: 5, Description: "CS:GO 进阶枪法与道具点位教学。", IsActive: true},
		{PlayerKey: "playerD", GameKey: "wow", RankLevel: 4, Description: "魔兽世界副本带练与职业天赋指导。", IsActive: true},
		{PlayerKey: "playerE", GameKey: "fifa", RankLevel: 3, Description: "FIFA 对战教学，阵型与防守反击专项训练。", IsActive: true},
		{PlayerKey: "playerF", GameKey: "fallguys", RankLevel: 2, Description: "糖豆人轻松娱乐局，氛围陪玩。", IsActive: true},
		{PlayerKey: "playerG", GameKey: "dota2", RankLevel: 7, Description: "DOTA2 冠军冲分局，高强度上分。", IsActive: true},
		// 待审核陪玩师示例：服务存在但暂未上架
		{PlayerKey: "playerH", GameKey: "pubg", RankLevel: 3, Description: "PUBG 生存技巧教学，待审核期间下架。", IsActive: false},
	}

	for _, spec := range specs {
		player := players[spec.PlayerKey]
		game := games[spec.GameKey]
		if player == nil || game == nil {
			continue
		}

		var rank model.GameRank
		if err := tx.Where("game_id = ? AND level = ?", game.ID, spec.RankLevel).First(&rank).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			return err
		}

		var existing model.PlayerService
		err := tx.Where("player_id = ? AND game_id = ? AND rank_id = ?", player.ID, game.ID, rank.ID).First(&existing).Error
		if err == nil {
			updates := map[string]interface{}{
				"description": spec.Description,
				"is_active":   spec.IsActive,
			}
			if err := tx.Model(&model.PlayerService{}).Where("id = ?", existing.ID).Updates(updates).Error; err != nil {
				return err
			}
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		row := &model.PlayerService{
			PlayerID:    player.ID,
			GameID:      game.ID,
			RankID:      rank.ID,
			Description: spec.Description,
			IsActive:    spec.IsActive,
		}
		if err := tx.Create(row).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedPlayerSchedules(tx *gorm.DB, players map[string]*model.Player) error {
	type scheduleSpec struct {
		PlayerKey        string
		WeeklySchedule   map[string]interface{}
		AutoOffline      bool
		MaxOrdersPerDay  int
	}

	fullDay := []interface{}{[]string{"00:00", "23:59"}}
	weekdayNight := []interface{}{[]string{"19:00", "23:00"}}
	weekendDay := []interface{}{[]string{"14:00", "22:00"}}

	specs := []scheduleSpec{
		{
			PlayerKey: "playerA",
			WeeklySchedule: map[string]interface{}{
				"timezone": "Asia/Shanghai",
				"mon":      fullDay, "tue": fullDay, "wed": fullDay, "thu": fullDay, "fri": fullDay, "sat": fullDay, "sun": fullDay,
			},
			AutoOffline:     true,
			MaxOrdersPerDay: 10,
		},
		{
			PlayerKey: "playerB",
			WeeklySchedule: map[string]interface{}{
				"timezone": "Asia/Shanghai",
				"mon":      weekdayNight, "tue": weekdayNight, "wed": weekdayNight, "thu": weekdayNight, "fri": weekdayNight,
				"sat": fullDay, "sun": fullDay,
			},
			AutoOffline:     true,
			MaxOrdersPerDay: 5,
		},
		{
			PlayerKey: "playerE",
			WeeklySchedule: map[string]interface{}{
				"timezone": "Asia/Shanghai",
				"mon":      []interface{}{}, "tue": []interface{}{}, "wed": []interface{}{}, "thu": []interface{}{}, "fri": []interface{}{},
				"sat": weekendDay, "sun": weekendDay,
			},
			AutoOffline:     true,
			MaxOrdersPerDay: 3,
		},
		{
			PlayerKey: "playerF",
			WeeklySchedule: map[string]interface{}{
				"timezone": "Asia/Shanghai",
				"mon":      fullDay, "tue": fullDay, "wed": fullDay, "thu": fullDay, "fri": fullDay, "sat": fullDay, "sun": fullDay,
			},
			AutoOffline:     false,
			MaxOrdersPerDay: 0,
		},
	}

	// 其余已认证陪玩师给默认晚间排班
	defaultSchedule := map[string]interface{}{
		"timezone": "Asia/Shanghai",
		"mon":      weekdayNight, "tue": weekdayNight, "wed": weekdayNight, "thu": weekdayNight, "fri": weekdayNight,
		"sat": []interface{}{}, "sun": []interface{}{},
	}
	indexed := map[string]scheduleSpec{}
	for _, s := range specs {
		indexed[s.PlayerKey] = s
	}
	for key, p := range players {
		if p == nil || p.VerificationStatus != model.VerificationVerified {
			continue
		}
		if _, ok := indexed[key]; ok {
			continue
		}
		indexed[key] = scheduleSpec{
			PlayerKey:       key,
			WeeklySchedule:  defaultSchedule,
			AutoOffline:     true,
			MaxOrdersPerDay: 5,
		}
	}

	for _, spec := range indexed {
		player := players[spec.PlayerKey]
		if player == nil || player.VerificationStatus != model.VerificationVerified {
			continue
		}
		b, err := json.Marshal(spec.WeeklySchedule)
		if err != nil {
			return fmt.Errorf("marshal weekly schedule for %s failed: %w", spec.PlayerKey, err)
		}

		var existing model.PlayerSchedule
		err = tx.Where("player_id = ?", player.ID).First(&existing).Error
		if err == nil {
			updates := map[string]interface{}{
				"weekly_schedule":   string(b),
				"auto_offline":      spec.AutoOffline,
				"max_orders_per_day": spec.MaxOrdersPerDay,
			}
			if err := tx.Model(&model.PlayerSchedule{}).Where("id = ?", existing.ID).Updates(updates).Error; err != nil {
				return err
			}
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		row := &model.PlayerSchedule{
			PlayerID:        player.ID,
			WeeklySchedule:  string(b),
			AutoOffline:     spec.AutoOffline,
			MaxOrdersPerDay: spec.MaxOrdersPerDay,
		}
		if err := tx.Create(row).Error; err != nil {
			return err
		}
	}

	return nil
}
