package auth

import "time"

const (
	// DefaultTokenDuration 是默认的 JWT Token 有效期
	DefaultTokenDuration = 24 * time.Hour

	// TokenAutoRefreshWindow 是在 Token 过期前多少时间内自动刷新
	// 当剩余时间小于此值时，中间件会自动刷新 Token 并在响应头中返回新 Token
	TokenAutoRefreshWindow = 15 * time.Minute

	// TokenRefreshRecommendationWindow 是在 Token 过期前多少时间提醒前端刷新
	// 当剩余时间小于此值但大于自动刷新窗口时，会在响应头中设置刷新建议
	TokenRefreshRecommendationWindow = 1 * time.Hour

	// TokenMinRefreshThreshold 是 Token 刷新的最小阈值
	// Token 必须至少在过期前此时间内才能被刷新
	TokenMinRefreshThreshold = 30 * time.Second

	// DefaultMaxRefreshWindow 是默认的 Token 刷新窗口
	// 从 Token 签发时间开始计算，超过此时间的 Token 将无法再刷新
	DefaultMaxRefreshWindow = 7 * 24 * time.Hour // 7 days
)
