package trtc

import (
	"bytes"
	"compress/zlib"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
	"gamelink/pkg/cache"
)

// Errors specific to TRTC domain.
var (
	ErrNotFound       = repository.ErrNotFound
	ErrVoiceDisabled  = apierr.BadRequest("该房间未启用语音功能")
	ErrNotMember      = apierr.Forbidden("用户不是房间成员")
	ErrAlreadyInVoice = apierr.BadRequest("已在语音房间中")
	ErrNotInVoice     = apierr.BadRequest("不在语音房间中")
)

// Config TRTC配置
type Config struct {
	SDKAppID  uint64 `json:"sdkAppId"`
	SecretKey string `json:"secretKey"`
	ExpireSec int64  `json:"expireSec"` // UserSig过期时间（秒）
}

// Service TRTC语音服务
type Service struct {
	config  *Config
	groups  repository.ChatGroupRepository
	members repository.ChatMemberRepository
	cache   cache.Cache
}

// NewService 创建TRTC服务
func NewService(
	config *Config,
	groups repository.ChatGroupRepository,
	members repository.ChatMemberRepository,
	cache cache.Cache,
) *Service {
	if config.ExpireSec <= 0 {
		config.ExpireSec = 86400 * 7 // 默认7天
	}
	return &Service{
		config:  config,
		groups:  groups,
		members: members,
		cache:   cache,
	}
}

// UserSigResponse UserSig响应
type UserSigResponse struct {
	UserSig  string `json:"userSig"`
	SDKAppID uint64 `json:"sdkAppId"`
	UserID   string `json:"userId"`
	RoomID   string `json:"roomId"`
	ExpireAt int64  `json:"expireAt"`
}

// GetUserSig 获取用户签名（用于加入语音房间）
func (s *Service) GetUserSig(ctx context.Context, roomID uint64, userID uint64) (*UserSigResponse, error) {
	// 获取房间信息
	room, err := s.groups.Get(ctx, roomID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, apierr.InternalError("获取房间失败").WithDetails(err.Error())
	}

	// 检查语音是否启用
	if !room.VoiceEnabled {
		return nil, ErrVoiceDisabled
	}

	// 检查用户是否是房间成员
	member, err := s.members.Get(ctx, roomID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotMember
		}
		return nil, apierr.InternalError("获取成员信息失败").WithDetails(err.Error())
	}
	if !member.IsActive {
		return nil, ErrNotMember
	}

	// 生成UserSig
	userIDStr := fmt.Sprintf("%d", userID)
	expireAt := time.Now().Unix() + s.config.ExpireSec
	userSig, err := s.generateUserSig(userIDStr, s.config.ExpireSec)
	if err != nil {
		return nil, apierr.InternalError("生成UserSig失败").WithDetails(err.Error())
	}

	return &UserSigResponse{
		UserSig:  userSig,
		SDKAppID: s.config.SDKAppID,
		UserID:   userIDStr,
		RoomID:   room.VoiceRoomID,
		ExpireAt: expireAt,
	}, nil
}

// StartVoice 开启房间语音
func (s *Service) StartVoice(ctx context.Context, roomID, userID uint64) error {
	room, err := s.groups.Get(ctx, roomID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return apierr.InternalError("获取房间失败").WithDetails(err.Error())
	}

	// 检查权限（只有房主可以开启语音）
	if room.CreatedBy != userID {
		return apierr.Forbidden("只有房主可以开启语音")
	}

	// 如果已经启用，直接返回
	if room.VoiceEnabled {
		return nil
	}

	// 生成语音房间ID
	room.VoiceEnabled = true
	room.VoiceRoomID = fmt.Sprintf("voice_%d_%d", roomID, time.Now().UnixNano())
	room.VoiceProvider = "trtc"
	room.VoiceSDKAppID = s.config.SDKAppID
	now := time.Now()
	room.VoiceStartedAt = &now

	if err := s.groups.Update(ctx, room); err != nil {
		return apierr.InternalError("更新房间失败").WithDetails(err.Error())
	}

	return nil
}

// StopVoice 关闭房间语音
func (s *Service) StopVoice(ctx context.Context, roomID, userID uint64) error {
	room, err := s.groups.Get(ctx, roomID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return apierr.InternalError("获取房间失败").WithDetails(err.Error())
	}

	// 检查权限
	if room.CreatedBy != userID {
		return apierr.Forbidden("只有房主可以关闭语音")
	}

	if !room.VoiceEnabled {
		return nil
	}

	// 计算语音时长
	now := time.Now()
	if room.VoiceStartedAt != nil {
		room.VoiceDuration = int(now.Sub(*room.VoiceStartedAt).Seconds())
	}

	room.VoiceEnabled = false
	room.VoiceEndedAt = &now

	if err := s.groups.Update(ctx, room); err != nil {
		return apierr.InternalError("更新房间失败").WithDetails(err.Error())
	}

	return nil
}

// GetVoiceStatus 获取语音状态
func (s *Service) GetVoiceStatus(ctx context.Context, roomID uint64) (*VoiceStatusResponse, error) {
	room, err := s.groups.Get(ctx, roomID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, apierr.InternalError("获取房间失败").WithDetails(err.Error())
	}

	return &VoiceStatusResponse{
		RoomID:       roomID,
		VoiceEnabled: room.VoiceEnabled,
		VoiceRoomID:  room.VoiceRoomID,
		Provider:     room.VoiceProvider,
		SDKAppID:     room.VoiceSDKAppID,
		StartedAt:    room.VoiceStartedAt,
		Duration:     room.VoiceDuration,
		MaxMembers:   room.VoiceMaxMembers,
	}, nil
}

// VoiceStatusResponse 语音状态响应
type VoiceStatusResponse struct {
	RoomID       uint64     `json:"roomId"`
	VoiceEnabled bool       `json:"voiceEnabled"`
	VoiceRoomID  string     `json:"voiceRoomId"`
	Provider     string     `json:"provider"`
	SDKAppID     uint64     `json:"sdkAppId"`
	StartedAt    *time.Time `json:"startedAt,omitempty"`
	Duration     int        `json:"duration"`
	MaxMembers   int        `json:"maxMembers"`
}

// generateUserSig 生成TRTC UserSig
// 参考腾讯云官方文档实现
func (s *Service) generateUserSig(userID string, expire int64) (string, error) {
	currTime := time.Now().Unix()

	// 构建签名内容
	sigDoc := map[string]interface{}{
		"TLS.ver":        "2.0",
		"TLS.identifier": userID,
		"TLS.sdkappid":   s.config.SDKAppID,
		"TLS.expire":     expire,
		"TLS.time":       currTime,
	}

	// 序列化
	sigDocJSON, err := json.Marshal(sigDoc)
	if err != nil {
		return "", err
	}

	// 计算签名
	sig := s.hmacSHA256(string(sigDocJSON))

	// 组合最终签名
	sigDoc["TLS.sig"] = sig
	finalJSON, err := json.Marshal(sigDoc)
	if err != nil {
		return "", err
	}

	// Base64编码
	return base64.StdEncoding.EncodeToString(s.compress(finalJSON)), nil
}

// hmacSHA256 计算HMAC-SHA256
func (s *Service) hmacSHA256(data string) string {
	h := hmac.New(sha256.New, []byte(s.config.SecretKey))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// compress 使用zlib压缩数据
func (s *Service) compress(data []byte) []byte {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	_, err := w.Write(data)
	if err != nil {
		// 压缩失败时返回原数据
		return data
	}
	if err := w.Close(); err != nil {
		return data
	}
	return buf.Bytes()
}

// RecordVoiceSession 记录语音会话（用于计费和统计）
func (s *Service) RecordVoiceSession(ctx context.Context, roomID uint64, duration int) error {
	room, err := s.groups.Get(ctx, roomID)
	if err != nil {
		return err
	}

	room.VoiceDuration += duration
	return s.groups.Update(ctx, room)
}
