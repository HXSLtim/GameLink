package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadConfig 文件上传配置
type UploadConfig struct {
	// MaxFileSize 最大文件大小（字节）
	MaxFileSize int64
	// AllowedMimeTypes 允许的MIME类型白名单
	AllowedMimeTypes []string
	// AllowedExtensions 允许的文件扩展名白名单
	AllowedExtensions []string
	// UploadPath 上传路径
	UploadPath string
	// RandomizeFilename 是否随机化文件名
	RandomizeFilename bool
	// PreserveExtension 保留原始扩展名
	PreserveExtension bool
	// CalculateHash 是否计算文件哈希
	CalculateHash bool
}

// DefaultUploadConfig 默认上传配置
var DefaultUploadConfig = UploadConfig{
	MaxFileSize:       10 * 1024 * 1024, // 10MB
	AllowedMimeTypes:  []string{"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp"},
	AllowedExtensions: []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
	UploadPath:        "./uploads",
	RandomizeFilename: true,
	PreserveExtension: true,
	CalculateHash:     true,
}

// UploadResult 上传结果
type UploadResult struct {
	OriginalName string
	SavedName    string
	FilePath     string
	FileSize     int64
	MimeType     string
	Extension    string
	Hash         string
}

// FileUpload 返回文件上传中间件
func FileUpload(config ...UploadConfig) gin.HandlerFunc {
	cfg := DefaultUploadConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	// 设置默认值
	if cfg.MaxFileSize == 0 {
		cfg.MaxFileSize = DefaultUploadConfig.MaxFileSize
	}
	if len(cfg.AllowedMimeTypes) == 0 {
		cfg.AllowedMimeTypes = DefaultUploadConfig.AllowedMimeTypes
	}
	if cfg.UploadPath == "" {
		cfg.UploadPath = DefaultUploadConfig.UploadPath
	}

	return func(c *gin.Context) {
		c.Next()
	}
}

// ValidateFile 验证上传的文件
func ValidateFile(file *multipart.FileHeader, config UploadConfig) error {
	// 检查文件大小
	if file.Size > config.MaxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", config.MaxFileSize)
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if len(config.AllowedExtensions) > 0 {
		allowed := false
		for _, allowedExt := range config.AllowedExtensions {
			if ext == strings.ToLower(allowedExt) {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("file extension %s is not allowed", ext)
		}
	}

	// 检查MIME类型
	if len(config.AllowedMimeTypes) > 0 {
		// 打开文件检测MIME类型
		src, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer func() {
			_ = src.Close()
		}()

		// 读取文件头512字节用于MIME类型检测
		buffer := make([]byte, 512)
		_, err = src.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read file: %w", err)
		}

		// 检测MIME类型
		mimeType := http.DetectContentType(buffer)
		// 移除可能的参数（如 "image/jpeg; charset=utf-8"）
		mimeType = strings.Split(mimeType, ";")[0]

		allowed := false
		for _, allowedMime := range config.AllowedMimeTypes {
			if mimeType == allowedMime {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("file type %s is not allowed", mimeType)
		}
	}

	return nil
}

// SaveFile 保存上传的文件
func SaveFile(c *gin.Context, file *multipart.FileHeader, config UploadConfig) (*UploadResult, error) {
	// 验证文件
	if err := ValidateFile(file, config); err != nil {
		return nil, err
	}

	// 生成文件名
	var savedName string
	ext := filepath.Ext(file.Filename)

	if config.RandomizeFilename {
		// 使用UUID生成随机文件名
		savedName = uuid.New().String()
		if config.PreserveExtension {
			savedName += ext
		}
	} else {
		// 使用时间戳 + 原始文件名
		timestamp := time.Now().Format("20060102150405")
		savedName = fmt.Sprintf("%s_%s", timestamp, file.Filename)
	}

	// 构建保存路径
	savePath := filepath.Join(config.UploadPath, savedName)

	// 保存文件
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// 计算文件哈希
	var hash string
	if config.CalculateHash {
		src, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file for hashing: %w", err)
		}
		defer func() {
			_ = src.Close()
		}()

		h := sha256.New()
		if _, err := io.Copy(h, src); err != nil {
			return nil, fmt.Errorf("failed to calculate hash: %w", err)
		}
		hash = hex.EncodeToString(h.Sum(nil))
	}

	// 检测MIME类型
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		_ = src.Close()
	}()

	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	mimeType := http.DetectContentType(buffer)
	mimeType = strings.Split(mimeType, ";")[0]

	return &UploadResult{
		OriginalName: file.Filename,
		SavedName:    savedName,
		FilePath:     savePath,
		FileSize:     file.Size,
		MimeType:     mimeType,
		Extension:    ext,
		Hash:         hash,
	}, nil
}

// GetImageConfig 获取图片上传配置
func GetImageConfig() UploadConfig {
	return UploadConfig{
		MaxFileSize: 5 * 1024 * 1024, // 5MB
		AllowedMimeTypes: []string{
			"image/jpeg",
			"image/jpg",
			"image/png",
			"image/gif",
			"image/webp",
		},
		AllowedExtensions: []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
		UploadPath:        "./uploads/images",
		RandomizeFilename: true,
		PreserveExtension: true,
		CalculateHash:     true,
	}
}

// GetVideoConfig 获取视频上传配置
func GetVideoConfig() UploadConfig {
	return UploadConfig{
		MaxFileSize: 100 * 1024 * 1024, // 100MB
		AllowedMimeTypes: []string{
			"video/mp4",
			"video/mpeg",
			"video/quicktime",
			"video/webm",
		},
		AllowedExtensions: []string{".mp4", ".mpeg", ".mov", ".webm"},
		UploadPath:        "./uploads/videos",
		RandomizeFilename: true,
		PreserveExtension: true,
		CalculateHash:     true,
	}
}

// GetAudioConfig 获取音频上传配置
func GetAudioConfig() UploadConfig {
	return UploadConfig{
		MaxFileSize: 20 * 1024 * 1024, // 20MB
		AllowedMimeTypes: []string{
			"audio/mpeg",
			"audio/wav",
			"audio/ogg",
			"audio/webm",
			"audio/aac",
		},
		AllowedExtensions: []string{".mp3", ".wav", ".ogg", ".webm", ".aac"},
		UploadPath:        "./uploads/audio",
		RandomizeFilename: true,
		PreserveExtension: true,
		CalculateHash:     true,
	}
}

// GetDocumentConfig 获取文档上传配置
func GetDocumentConfig() UploadConfig {
	return UploadConfig{
		MaxFileSize: 10 * 1024 * 1024, // 10MB
		AllowedMimeTypes: []string{
			"application/pdf",
			"application/msword",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		},
		AllowedExtensions: []string{".pdf", ".doc", ".docx"},
		UploadPath:        "./uploads/documents",
		RandomizeFilename: true,
		PreserveExtension: true,
		CalculateHash:     true,
	}
}
