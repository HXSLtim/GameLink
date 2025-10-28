package repository

import (
    "context"

    "gamelink/internal/model"
)

// PlayerTagRepository 管理陪玩技能标签关联。
type PlayerTagRepository interface {
    // ReplaceTags 用新集合替换玩家的技能标签（幂等）。
    ReplaceTags(ctx context.Context, playerID uint64, tags []string) error
    List(ctx context.Context, playerID uint64) ([]model.PlayerSkillTag, error)
}

