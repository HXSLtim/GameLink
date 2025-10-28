package model

import "encoding/json"

// OperationLog 记录后台关键业务操作日志，用于审计与可视化。
type OperationLog struct {
    Base
    EntityType   string          `json:"entity_type" gorm:"size:32;index"` // order | payment
    EntityID     uint64          `json:"entity_id" gorm:"index"`
    ActorUserID  *uint64         `json:"actor_user_id" gorm:"index"`
    Action       string          `json:"action" gorm:"size:64;index"`
    Reason       string          `json:"reason,omitempty" gorm:"type:text"`
    MetadataJSON json.RawMessage `json:"metadata,omitempty" gorm:"type:json"`
}

