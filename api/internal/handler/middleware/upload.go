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
	// StrictMimeExtMatch 严格校验MIME类型与扩展名一致性（防止后缀伪装）
	StrictMimeExtMatch bool
}

// mimeToExtensions MIME类型到合法扩展名的映射（用于防止后缀伪装攻击）
var mimeToExtensions = map[string][]string{
	// 图片
	"image/jpeg": {".jpg", ".jpeg"},
	"image/png":  {".png"},
	"image/gif":  {".gif"},
	"image/webp": {".webp"},
	// 视频
	"video/mp4":       {".mp4"},
	"video/mpeg":      {".mpeg", ".mpg"},
	"video/quicktime": {".mov"},
	"video/webm":      {".webm"},
	// 音频
	"audio/mpeg": {".mp3"},
	"audio/wav":  {".wav"},
	"audio/ogg":  {".ogg"},
	"audio/webm": {".webm"},
	"audio/aac":  {".aac"},
	// 文档
	"application/pdf":    {".pdf"},
	"application/msword": {".doc"},
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": {".docx"},
}

// dangerousExtensions 危险扩展名黑名单（即使MIME类型正确也拒绝）
var dangerousExtensions = map[string]bool{
	".php": true, ".php3": true, ".php4": true, ".php5": true, ".phtml": true,
	".asp": true, ".aspx": true,
	".jsp": true, ".jspx": true,
	".exe": true, ".dll": true, ".bat": true, ".cmd": true, ".com": true,
	".sh": true, ".bash": true,
	".py": true, ".pl": true, ".rb": true,
	".js": true, ".mjs": true,
	".html": true, ".htm": true, ".xhtml": true,
	".svg": true, // SVG可包含脚本
	".xml": true,
}

// DefaultUploadConfig 默认上传配置
var DefaultUploadConfig = UploadConfig{
	MaxFileSize:        10 * 1024 * 1024, // 10MB
	AllowedMimeTypes:   []string{"image/jpeg", "image/png", "image/gif", "image/webp"},
	AllowedExtensions:  []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
	UploadPath:         "./uploads",
	RandomizeFilename:  true,
	PreserveExtension:  true,
	CalculateHash:      true,
	StrictMimeExtMatch: true, // 默认启用严格校验
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
	// 1. 检查文件大小
	if file.Size > config.MaxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", config.MaxFileSize)
	}

	// 2. 获取并规范化扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))

	// 3. 检查危险扩展名黑名单（优先级最高）
	if isDangerousExtension(file.Filename) {
		return fmt.Errorf("file extension %s is not allowed for security reasons", ext)
	}

	// 4. 检查文件扩展名白名单
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

	// 5. 检测真实MIME类型（基于文件内容magic bytes）
	if len(config.AllowedMimeTypes) > 0 {
		src, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer func() { _ = src.Close() }()

		// 读取文件头512字节用于MIME类型检测
		buffer := make([]byte, 512)
		n, err := src.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read file: %w", err)
		}

		// 检测真实MIME类型
		mimeType := http.DetectContentType(buffer[:n])
		mimeType = strings.Split(mimeType, ";")[0] // 移除参数

		// 检查MIME类型白名单
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

		// 6. 严格校验：MIME类型与扩展名必须一致（防止后缀伪装攻击）
		if config.StrictMimeExtMatch {
			if !isMimeExtensionMatch(mimeType, ext) {
				return fmt.Errorf("file extension %s does not match content type %s (possible file spoofing)", ext, mimeType)
			}
		}
	}

	return nil
}

// isDangerousExtension 检查是否包含危险扩展名（支持双扩展名检测）
func isDangerousExtension(filename string) bool {
	// 转小写
	lower := strings.ToLower(filename)

	// 检查最终扩展名
	ext := filepath.Ext(lower)
	if dangerousExtensions[ext] {
		return true
	}

	// 检查双扩展名攻击（如 image.jpg.php）
	// 移除最后一个扩展名后再检查
	withoutExt := strings.TrimSuffix(lower, ext)
	secondExt := filepath.Ext(withoutExt)
	if dangerousExtensions[secondExt] {
		return true
	}

	// 检查文件名中是否包含危险扩展名（如 image.php.jpg）
	for dangerous := range dangerousExtensions {
		if strings.Contains(lower, dangerous+".") {
			return true
		}
	}

	return false
}

// isMimeExtensionMatch 检查MIME类型与扩展名是否匹配
func isMimeExtensionMatch(mimeType, ext string) bool {
	allowedExts, ok := mimeToExtensions[mimeType]
	if !ok {
		// 未知MIME类型，默认不匹配
		return false
	}

	ext = strings.ToLower(ext)
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

// getExtensionForMime 根据MIME类型获取推荐的扩展名
func getExtensionForMime(mimeType string) string {
	if exts, ok := mimeToExtensions[mimeType]; ok && len(exts) > 0 {
		return exts[0]
	}
	return ""
}

// SaveFile 保存上传的文件
func SaveFile(c *gin.Context, file *multipart.FileHeader, config UploadConfig) (*UploadResult, error) {
	// 验证文件
	if err := ValidateFile(file, config); err != nil {
		return nil, err
	}

	// 检测真实MIME类型
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	buffer := make([]byte, 512)
	n, err := src.Read(buffer)
	if err != nil && err != io.EOF {
		_ = src.Close()
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	mimeType := http.DetectContentType(buffer[:n])
	mimeType = strings.Split(mimeType, ";")[0]
	_ = src.Close()

	// 生成安全的文件名
	var savedName string
	var finalExt string

	if config.PreserveExtension {
		// 使用基于MIME类型的安全扩展名（而非用户提供的扩展名）
		finalExt = getExtensionForMime(mimeType)
		if finalExt == "" {
			// 回退到原始扩展名（已通过白名单校验）
			finalExt = strings.ToLower(filepath.Ext(file.Filename))
		}
	}

	if config.RandomizeFilename {
		// 使用UUID生成随机文件名
		savedName = uuid.New().String() + finalExt
	} else {
		// 使用时间戳 + 安全化的文件名
		timestamp := time.Now().Format("20060102150405")
		// 只保留字母数字和下划线，移除原始扩展名
		baseName := sanitizeFilename(strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename)))
		savedName = fmt.Sprintf("%s_%s%s", timestamp, baseName, finalExt)
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
		defer func() { _ = src.Close() }()

		h := sha256.New()
		if _, err := io.Copy(h, src); err != nil {
			return nil, fmt.Errorf("failed to calculate hash: %w", err)
		}
		hash = hex.EncodeToString(h.Sum(nil))
	}

	return &UploadResult{
		OriginalName: file.Filename,
		SavedName:    savedName,
		FilePath:     savePath,
		FileSize:     file.Size,
		MimeType:     mimeType,
		Extension:    finalExt,
		Hash:         hash,
	}, nil
}

// sanitizeFilename 清理文件名，只保留安全字符
func sanitizeFilename(name string) string {
	var result strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			result.WriteRune(r)
		}
	}
	s := result.String()
	if s == "" {
		return "file"
	}
	// 限制长度
	if len(s) > 50 {
		s = s[:50]
	}
	return s
}

// GetImageConfig 获取图片上传配置
func GetImageConfig() UploadConfig {
	return UploadConfig{
		MaxFileSize: 5 * 1024 * 1024, // 5MB
		AllowedMimeTypes: []string{
			"image/jpeg",
			"image/png",
			"image/gif",
			"image/webp",
		},
		AllowedExtensions:  []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
		UploadPath:         "./uploads/images",
		RandomizeFilename:  true,
		PreserveExtension:  true,
		CalculateHash:      true,
		StrictMimeExtMatch: true,
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
		AllowedExtensions:  []string{".mp4", ".mpeg", ".mov", ".webm"},
		UploadPath:         "./uploads/videos",
		RandomizeFilename:  true,
		PreserveExtension:  true,
		CalculateHash:      true,
		StrictMimeExtMatch: true,
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
		AllowedExtensions:  []string{".mp3", ".wav", ".ogg", ".webm", ".aac"},
		UploadPath:         "./uploads/audio",
		RandomizeFilename:  true,
		PreserveExtension:  true,
		CalculateHash:      true,
		StrictMimeExtMatch: true,
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
		AllowedExtensions:  []string{".pdf", ".doc", ".docx"},
		UploadPath:         "./uploads/documents",
		RandomizeFilename:  true,
		PreserveExtension:  true,
		CalculateHash:      true,
		StrictMimeExtMatch: true,
	}
}
