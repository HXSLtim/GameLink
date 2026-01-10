package model

// UserBehavior 用户行为事件
// 记录页面及操作行为与上下文
type UserBehavior struct {
	Base
	UserID     uint64 `json:"userId" gorm:"column:user_id;index;not null;comment:用户ID"`
	Action     string `json:"action" gorm:"size:64;index;not null;comment:行为动作"`
	TargetType string `json:"targetType" gorm:"column:target_type;size:32;comment:目标类型"`
	TargetID   uint64 `json:"targetId" gorm:"column:target_id;comment:目标ID"`
	Duration   int    `json:"duration,omitempty" gorm:"comment:持续时长(毫秒)"`
	PagePath   string `json:"pagePath" gorm:"column:page_path;size:255;comment:页面路径"`
	SessionID  string `json:"sessionId" gorm:"column:session_id;size:128;comment:会话ID"`
	Metadata   string `json:"metadata,omitempty" gorm:"type:text;comment:扩展数据JSON"`
}

// TableName 指定表名
func (UserBehavior) TableName() string { return "user_behaviors" }
