package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT工具类
//
// 功能：
// 1. 生成JWT Token
// 2. 解析和验证JWT Token
// 3. 提取用户信息

// Claims 定义JWT载荷结构
type Claims struct {
	UserID uint64 `json:"user_id"` // 用户ID
	Role   string `json:"role"`    // 用户角色
	jwt.RegisteredClaims
}

// JWTManager JWT管理器
type JWTManager struct {
	secretKey     string        // 签名密钥
	tokenDuration time.Duration // Token有效期
	maxRefresh    time.Duration // 允许刷新窗口（自签发起）
}

// NewJWTManager 创建JWT管理器
func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
		maxRefresh:    readMaxRefreshWindow(),
	}
}

func readMaxRefreshWindow() time.Duration {
	v := os.Getenv("JWT_MAX_REFRESH")
	if v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return 7 * 24 * time.Hour
}

// GenerateToken 生成JWT Token
//
// 参数：
// - userID: 用户ID
// - role: 用户角色
//
// 返回：
// - token: JWT字符串
// - err: 错误信息
func (manager *JWTManager) GenerateToken(userID uint64, role string) (string, error) {
	// 创建Claims
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(manager.tokenDuration)),
			// 设置签发时间
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// 设置生效时间（立即生效）
			NotBefore: jwt.NewNumericDate(time.Now()),
			// 设置签发者
			Issuer: "gamelink",
		},
	}

	// 创建Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名Token
	signedToken, err := token.SignedString([]byte(manager.secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// VerifyToken 验证JWT Token
//
// 参数：
// - tokenString: JWT字符串
//
// 返回：
// - claims: 用户信息
// - err: 错误信息
func (manager *JWTManager) VerifyToken(tokenString string) (*Claims, error) {
	// 解析Token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名算法")
		}
		return []byte(manager.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证Token是否有效
	if !token.Valid {
		return nil, errors.New("无效的Token")
	}

	// 提取Claims
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("无法解析Token Claims")
	}

	return claims, nil
}

// RefreshToken 刷新Token
//
// 当Token快要过期时，可以生成新的Token
func (manager *JWTManager) RefreshToken(claims *Claims) (string, error) {
	// 检查Token是否还有足够的时间
	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		return "", errors.New("Token还未到刷新时间")
	}

	// 限制刷新窗口：签发时间距今不得超过 maxRefresh
	if !claims.IssuedAt.Time.IsZero() && time.Since(claims.IssuedAt.Time) > manager.maxRefresh {
		return "", errors.New("Token已超过可刷新窗口")
	}

	// 生成新的Token
	return manager.GenerateToken(claims.UserID, claims.Role)
}

// ExtractTokenFromHeader 从HTTP头中提取Token
//
// Authorization Header格式：Bearer <token>
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("缺少Authorization头")
	}

	// 检查Bearer前缀
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.New("Authorization头格式错误，应为'Bearer <token>'")
	}

	// 提取Token
	token := authHeader[len(bearerPrefix):]
	if token == "" {
		return "", errors.New("Token为空")
	}

	return token, nil
}

// IsTokenExpired 检查Token是否过期
func IsTokenExpired(claims *Claims) bool {
	return time.Now().After(claims.ExpiresAt.Time)
}

// GetTokenRemainingTime 获取Token剩余有效时间
func GetTokenRemainingTime(claims *Claims) time.Duration {
	return time.Until(claims.ExpiresAt.Time)
}
