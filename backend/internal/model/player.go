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
	UserID             uint64             `json:"userId" gorm:"column:user_id;uniqueIndex"`
	Nickname           string             `json:"nickname,omitempty" gorm:"size:64"`
	Bio                string             `json:"bio,omitempty" gorm:"type:text"`
	Rank               string             `json:"rank,omitempty" gorm:"size:32"`
	RatingAverage      float32            `json:"ratingAverage" gorm:"column:rating_average;default:0;check:rating_average >= 0 AND rating_average <= 5"`
	RatingCount        uint32             `json:"ratingCount" gorm:"column:rating_count;default:0"`
	HourlyRateCents    int64              `json:"hourlyRateCents" gorm:"column:hourly_rate_cents"`
	MainGameID         uint64             `json:"mainGameId,omitempty" gorm:"column:main_game_id;index"`
	VerificationStatus VerificationStatus `json:"verificationStatus" gorm:"column:verification_status;size:32;index"`
}
