package model

import (
	"time"
)

// UploadType 上传类型
type UploadType string

const (
	// UploadTypeAvatar 用户头像
	UploadTypeAvatar UploadType = "avatar"
	// UploadTypeCertification 认证材料
	UploadTypeCertification UploadType = "certification"
	// UploadTypeIDCard 身份证
	UploadTypeIDCard UploadType = "id_card"
	// UploadTypeGameScreenshot 游戏截图
	UploadTypeGameScreenshot UploadType = "game_screenshot"
	// UploadTypeReviewImage 评价图片
	UploadTypeReviewImage UploadType = "review_image"
	// UploadTypeChatImage 聊天图片
	UploadTypeChatImage UploadType = "chat_image"
	// UploadTypeOther 其他
	UploadTypeOther UploadType = "other"
)

// UploadStatus 上传状态
type UploadStatus string

const (
	// UploadStatusPending 待处理
	UploadStatusPending UploadStatus = "pending"
	// UploadStatusProcessing 处理中
	UploadStatusProcessing UploadStatus = "processing"
	// UploadStatusCompleted 已完成
	UploadStatusCompleted UploadStatus = "completed"
	// UploadStatusFailed 失败
	UploadStatusFailed UploadStatus = "failed"
	// UploadStatusDeleted 已删除
	UploadStatusDeleted UploadStatus = "deleted"
)

// Upload 文件上传记录
type Upload struct {
	ID         uint64       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint64       `gorm:"not null;index:idx_uploads_user" json:"userId"`
	FileName   string       `gorm:"size:255;not null" json:"fileName"`
	FilePath   string       `gorm:"size:500;not null" json:"filePath"`
	FileURL    string       `gorm:"size:500" json:"fileUrl"`
	FileSize   int64        `gorm:"not null" json:"fileSize"`
	MimeType   string       `gorm:"size:100;not null" json:"mimeType"`
	UploadType UploadType   `gorm:"size:50;not null;index:idx_uploads_type" json:"uploadType"`
	Status     UploadStatus `gorm:"size:20;default:'pending';index:idx_uploads_status" json:"status"`
	Hash       string       `gorm:"size:64;index:idx_uploads_hash" json:"hash"` // 文件哈希，用于去重
	Width      int          `gorm:"default:0" json:"width,omitempty"`            // 图片宽度
	Height     int          `gorm:"default:0" json:"height,omitempty"`           // 图片高度
	ErrorMsg   string       `gorm:"size:500" json:"errorMsg,omitempty"`          // 错误信息
	CreatedAt  time.Time    `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time    `gorm:"autoUpdateTime" json:"updatedAt"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (Upload) TableName() string {
	return "uploads"
}

// IsImage 判断是否为图片
func (u *Upload) IsImage() bool {
	switch u.MimeType {
	case "image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp", "image/bmp":
		return true
	default:
		return false
	}
}

// IsVideo 判断是否为视频
func (u *Upload) IsVideo() bool {
	switch u.MimeType {
	case "video/mp4", "video/mpeg", "video/quicktime", "video/x-msvideo", "video/webm":
		return true
	default:
		return false
	}
}

// IsAudio 判断是否为音频
func (u *Upload) IsAudio() bool {
	switch u.MimeType {
	case "audio/mpeg", "audio/wav", "audio/ogg", "audio/webm", "audio/aac":
		return true
	default:
		return false
	}
}

// GetSizeInMB 获取文件大小（MB）
func (u *Upload) GetSizeInMB() float64 {
	return float64(u.FileSize) / 1024 / 1024
}
