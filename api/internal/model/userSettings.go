package model

// UserSettings stores user preferences.
type UserSettings struct {
	Base
	UserID        uint64 `json:"userId" gorm:"column:user_id;uniqueIndex;not null"`
	Theme         string `json:"theme" gorm:"size:16;default:'auto'"`
	Language      string `json:"language" gorm:"size:32;default:'zh-CN'"`
	Notifications string `json:"notifications" gorm:"column:notifications;type:json;default:'{}'"`
	Privacy       string `json:"privacy" gorm:"column:privacy;type:json;default:'{}'"`
}

// TableName specifies table name.
func (UserSettings) TableName() string {
	return "user_settings"
}
