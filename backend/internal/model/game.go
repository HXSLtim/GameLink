package model

// Game represents a supported game and its metadata.
type Game struct {
	Base
	Key         string `json:"key" gorm:"size:64;uniqueIndex"` // unique slug/key, e.g. "lol", "dota2"
	Name        string `json:"name" gorm:"size:128"`
	Category    string `json:"category,omitempty" gorm:"size:64"` // e.g. moba/fps
	IconURL     string `json:"icon_url,omitempty" gorm:"size:255"`
	Description string `json:"description,omitempty" gorm:"type:text"`
}
