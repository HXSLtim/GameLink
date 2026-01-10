// Package auth provides authentication services including WeChat mini-program login.
package auth

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
	"gamelink/pkg/auth"
)

// WeChatLoginRequest 微信登录请求
type WeChatLoginRequest struct {
	Code          string `json:"code" binding:"required"`  // 微信登录凭证
	EncryptedData string `json:"encryptedData,omitempty"`  // 加密数据（获取手机号）
	IV            string `json:"iv,omitempty"`             // 解密向量
	ReferralCode  string `json:"referralCode,omitempty"`   // 推荐码
}

// WeChatLoginResponse 微信登录响应
type WeChatLoginResponse struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresIn    int64     `json:"expiresIn"`
	User         *UserInfo `json:"user"`
}

// UserInfo 用户信息（登录响应用）
type UserInfo struct {
	ID          uint64 `json:"id"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Phone       string `json:"phone"`
	IsPlayer    bool   `json:"isPlayer"`
	CurrentRole string `json:"currentRole"`
}

// WeChatSession 微信 code2Session 响应
type WeChatSession struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid,omitempty"`
	ErrCode    int    `json:"errcode,omitempty"`
	ErrMsg     string `json:"errmsg,omitempty"`
}

// WeChatPhoneInfo 微信手机号信息
type WeChatPhoneInfo struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
}

var (
	ErrWeChatCodeInvalid   = errors.New("微信登录凭证无效")
	ErrWeChatSessionFailed = errors.New("获取微信会话失败")
	ErrDecryptFailed       = errors.New("解密数据失败")
)

// WeChatAuthService 微信认证服务
type WeChatAuthService struct {
	users   repository.UserRepository
	players repository.PlayerRepository
}

// NewWeChatAuthService 创建微信认证服务
func NewWeChatAuthService(users repository.UserRepository, players repository.PlayerRepository) *WeChatAuthService {
	return &WeChatAuthService{
		users:   users,
		players: players,
	}
}

// WeChatLogin 微信小程序登录
func (s *WeChatAuthService) WeChatLogin(ctx context.Context, req WeChatLoginRequest) (*WeChatLoginResponse, error) {
	// 1. 调用微信 code2Session 接口
	session, err := s.code2Session(req.Code)
	if err != nil {
		return nil, err
	}

	// 2. 解密手机号（如果提供了加密数据）
	var phone string
	if req.EncryptedData != "" && req.IV != "" {
		phoneInfo, decryptErr := s.decryptPhoneNumber(session.SessionKey, req.EncryptedData, req.IV)
		if decryptErr == nil {
			phone = phoneInfo.PhoneNumber
		}
		// 解密失败不阻塞登录
	}

	// 3. 查找或创建用户
	user, isNew, err := s.findOrCreateWeChatUser(ctx, session.OpenID, session.UnionID, phone)
	if err != nil {
		return nil, err
	}

	// 4. 处理推荐码（仅新用户）
	if isNew && req.ReferralCode != "" {
		_ = s.processReferral(ctx, user.ID, req.ReferralCode)
	}

	// 5. 检查用户是否是陪玩师
	isPlayer := s.checkIsPlayer(ctx, user.ID)

	// 6. 生成 JWT Token
	currentRole := "user"
	if isPlayer && user.Role == model.RolePlayer {
		currentRole = "player"
	}

	accessToken, err := s.generateWeChatToken(user, isPlayer, currentRole)
	if err != nil {
		return nil, apierr.InternalError("生成Token失败")
	}

	refreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return nil, apierr.InternalError("生成刷新Token失败")
	}

	return &WeChatLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    86400, // 24小时
		User: &UserInfo{
			ID:          user.ID,
			Nickname:    user.Name,
			Avatar:      user.AvatarURL,
			Phone:       user.Phone,
			IsPlayer:    isPlayer,
			CurrentRole: currentRole,
		},
	}, nil
}

// code2Session 调用微信 code2Session 接口
func (s *WeChatAuthService) code2Session(code string) (*WeChatSession, error) {
	appID := os.Getenv("WECHAT_APPID")
	secret := os.Getenv("WECHAT_SECRET")

	if appID == "" || secret == "" {
		return nil, apierr.InternalError("微信配置缺失")
	}

	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		appID, secret, code,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, ErrWeChatSessionFailed
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrWeChatSessionFailed
	}

	var session WeChatSession
	if err := json.Unmarshal(body, &session); err != nil {
		return nil, ErrWeChatSessionFailed
	}

	if session.ErrCode != 0 {
		return nil, fmt.Errorf("%w: %s", ErrWeChatCodeInvalid, session.ErrMsg)
	}

	return &session, nil
}

// decryptPhoneNumber 解密微信手机号
func (s *WeChatAuthService) decryptPhoneNumber(sessionKey, encryptedData, iv string) (*WeChatPhoneInfo, error) {
	// Base64 解码
	keyBytes, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, ErrDecryptFailed
	}

	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, ErrDecryptFailed
	}

	dataBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, ErrDecryptFailed
	}

	// AES-128-CBC 解密
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, ErrDecryptFailed
	}

	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(dataBytes, dataBytes)

	// 去除 PKCS7 填充
	dataBytes = pkcs7Unpad(dataBytes)

	// 解析 JSON
	var result struct {
		PhoneNumber     string `json:"phoneNumber"`
		PurePhoneNumber string `json:"purePhoneNumber"`
		CountryCode     string `json:"countryCode"`
	}
	if err := json.Unmarshal(dataBytes, &result); err != nil {
		return nil, ErrDecryptFailed
	}

	return &WeChatPhoneInfo{
		PhoneNumber:     result.PhoneNumber,
		PurePhoneNumber: result.PurePhoneNumber,
		CountryCode:     result.CountryCode,
	}, nil
}

// pkcs7Unpad 去除 PKCS7 填充
func pkcs7Unpad(data []byte) []byte {
	length := len(data)
	if length == 0 {
		return data
	}
	padding := int(data[length-1])
	if padding > length {
		return data
	}
	return data[:length-padding]
}

// findOrCreateWeChatUser 查找或创建微信用户
func (s *WeChatAuthService) findOrCreateWeChatUser(ctx context.Context, openID, unionID, phone string) (*model.User, bool, error) {
	// 先通过 OpenID 查找
	user, err := s.users.GetByWeChatOpenID(ctx, openID)
	if err == nil && user != nil {
		// 更新手机号（如果有新的）
		if phone != "" && user.Phone == "" {
			user.Phone = phone
			_ = s.users.Update(ctx, user)
		}
		return user, false, nil
	}

	// 如果有 UnionID，尝试通过 UnionID 查找
	if unionID != "" {
		user, err = s.users.GetByWeChatUnionID(ctx, unionID)
		if err == nil && user != nil {
			// 更新 OpenID
			user.WeChatOpenID = openID
			if phone != "" && user.Phone == "" {
				user.Phone = phone
			}
			_ = s.users.Update(ctx, user)
			return user, false, nil
		}
	}

	// 创建新用户
	now := time.Now()
	newUser := &model.User{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Name:          fmt.Sprintf("用户%d", now.UnixNano()%1000000),
		Phone:         phone,
		WeChatOpenID:  openID,
		WeChatUnionID: unionID,
		Role:          model.RoleUser,
		Status:        model.UserStatusActive,
	}

	if err := s.users.Create(ctx, newUser); err != nil {
		return nil, false, apierr.InternalError("创建用户失败")
	}

	return newUser, true, nil
}

// checkIsPlayer 检查用户是否是陪玩师
func (s *WeChatAuthService) checkIsPlayer(ctx context.Context, userID uint64) bool {
	if s.players == nil {
		return false
	}
	player, err := s.players.GetByUserID(ctx, userID)
	return err == nil && player != nil && player.VerificationStatus == model.VerificationVerified
}

// processReferral 处理推荐码
func (s *WeChatAuthService) processReferral(_ context.Context, _ uint64, _ string) error {
	// TODO: 调用推荐服务处理推荐关系
	return nil
}

// generateWeChatToken 生成微信登录 Token
func (s *WeChatAuthService) generateWeChatToken(user *model.User, isPlayer bool, currentRole string) (string, error) {
	claims := auth.CustomClaims{
		UserID:      user.ID,
		Role:        string(user.Role),
		IsPlayer:    isPlayer,
		CurrentRole: currentRole,
	}
	return auth.GenerateToken(claims)
}

// generateRefreshToken 生成刷新 Token
func (s *WeChatAuthService) generateRefreshToken(userID uint64) (string, error) {
	claims := auth.CustomClaims{
		UserID:    userID,
		IsRefresh: true,
	}
	return auth.GenerateRefreshToken(claims)
}

// RefreshAccessToken 使用刷新 Token 获取新的访问 Token
func (s *WeChatAuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (*WeChatLoginResponse, error) {
	// 验证刷新 Token
	claims, err := auth.VerifyCustomToken(refreshToken)
	if err != nil {
		return nil, apierr.Unauthorized("刷新Token无效")
	}

	if !claims.IsRefresh {
		return nil, apierr.Unauthorized("不是有效的刷新Token")
	}

	// 获取用户信息
	user, err := s.users.Get(ctx, claims.UserID)
	if err != nil {
		return nil, apierr.Unauthorized("用户不存在")
	}

	if user.Status != model.UserStatusActive {
		return nil, apierr.Forbidden("用户已被禁用")
	}

	// 检查是否是陪玩师
	isPlayer := s.checkIsPlayer(ctx, user.ID)
	currentRole := "user"
	if isPlayer && user.Role == model.RolePlayer {
		currentRole = "player"
	}

	// 生成新的访问 Token
	accessToken, err := s.generateWeChatToken(user, isPlayer, currentRole)
	if err != nil {
		return nil, apierr.InternalError("生成Token失败")
	}

	// 生成新的刷新 Token
	newRefreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return nil, apierr.InternalError("生成刷新Token失败")
	}

	return &WeChatLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    86400,
		User: &UserInfo{
			ID:          user.ID,
			Nickname:    user.Name,
			Avatar:      user.AvatarURL,
			Phone:       user.Phone,
			IsPlayer:    isPlayer,
			CurrentRole: currentRole,
		},
	}, nil
}
