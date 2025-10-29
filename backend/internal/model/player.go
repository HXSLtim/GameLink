package model

// VerificationStatus indicates identity/skill verification flow.
type VerificationStatus string

// VerificationStatus values indicate the verification flow state.
const (
	VerificationPending  VerificationStatus = "pending"
	VerificationVerified VerificationStatus = "verified"
	VerificationRejected VerificationStatus = "rejected"
)

// Player is the pro/陪玩 profile bound to a User.
type Player struct {
	Base
	UserID             uint64             `json:"user_id" gorm:"uniqueIndex"`
	Nickname           string             `json:"nickname,omitempty" gorm:"size:64"`
	Bio                string             `json:"bio,omitempty" gorm:"type:text"`
	Rank               string             `json:"rank,omitempty" gorm:"size:32"`
	RatingAverage      float32            `json:"rating_average" gorm:"default:0;check:rating_average >= 0 AND rating_average <= 5"`
	RatingCount        uint32             `json:"rating_count" gorm:"default:0"`
	HourlyRateCents    int64              `json:"hourly_rate_cents"`
	MainGameID         uint64             `json:"main_game_id,omitempty" gorm:"index"`
	VerificationStatus VerificationStatus `json:"verification_status" gorm:"size:32;index"`
}
