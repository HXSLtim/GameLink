package upload

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service/external"
	"gamelink/internal/service/oss"
)

// Service handles upload workflows with optional OSS integration.
type Service struct {
	uploads repository.UploadRepository
	ossSvc  *oss.Service
	ossCfg  *external.Config
}

// Result represents a processed upload output.
type Result struct {
	FilePath string `json:"filePath"`
	FileURL  string `json:"fileUrl"`
	Hash     string `json:"hash"`
}

// NewService creates an upload service.
func NewService(uploadRepo repository.UploadRepository, cfg *external.Config) *Service {
	var ossSvc *oss.Service
	if cfg != nil {
		ossSvc = oss.NewService(cfg)
	}
	return &Service{
		uploads: uploadRepo,
		ossSvc:  ossSvc,
		ossCfg:  cfg,
	}
}

// UploadImage saves an image and optionally uploads it to OSS.
func (s *Service) UploadImage(c *gin.Context, userID uint64, file *multipart.FileHeader, uploadType model.UploadType, basePath string) (*Result, error) {
	cfg := middleware.GetImageConfig()
	cfg.UploadPath = basePath

	if err := os.MkdirAll(cfg.UploadPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("create upload path failed: %w", err)
	}

	localRes, err := middleware.SaveFile(c, file, cfg)
	if err != nil {
		return nil, err
	}

	filePath := localRes.FilePath
	fileURL := buildPublicURL(localRes.FilePath)

	// Upload to OSS if enabled
	if s.ossSvc != nil && s.ossCfg != nil && s.ossCfg.OSS.Enabled {
		key := buildOSSKey(uploadType, localRes.SavedName)
		f, err := os.Open(localRes.FilePath)
		if err != nil {
			return nil, fmt.Errorf("open local file failed: %w", err)
		}
		defer func() { _ = f.Close() }()

		url, err := s.ossSvc.Upload(c.Request.Context(), key, f)
		if err != nil {
			return nil, fmt.Errorf("oss upload failed: %w", err)
		}

		filePath = key
		fileURL = url
	}

	// Persist upload record (best-effort)
	if s.uploads != nil {
		record := &model.Upload{
			UserID:     userID,
			FileName:   localRes.OriginalName,
			FilePath:   filePath,
			FileURL:    fileURL,
			FileSize:   localRes.FileSize,
			MimeType:   localRes.MimeType,
			UploadType: uploadType,
			Status:     model.UploadStatusCompleted,
			Hash:       localRes.Hash,
		}
		_ = s.uploads.Create(c.Request.Context(), record)
	}

	return &Result{
		FilePath: filePath,
		FileURL:  fileURL,
		Hash:     localRes.Hash,
	}, nil
}

// UploadImages uploads multiple images and returns their URLs.
func (s *Service) UploadImages(c *gin.Context, userID uint64, files []*multipart.FileHeader, uploadType model.UploadType, basePath string) ([]Result, error) {
	results := make([]Result, 0, len(files))
	for _, file := range files {
		res, err := s.UploadImage(c, userID, file, uploadType, basePath)
		if err != nil {
			return nil, err
		}
		results = append(results, *res)
	}
	return results, nil
}

func buildOSSKey(uploadType model.UploadType, savedName string) string {
	date := time.Now().Format("2006/01/02")
	segments := []string{"uploads", "images", string(uploadType), date, savedName}
	return filepath.ToSlash(filepath.Join(segments...))
}

func buildPublicURL(filePath string) string {
	if filePath == "" {
		return ""
	}
	baseURL := strings.TrimSpace(os.Getenv("UPLOAD_PUBLIC_BASE_URL"))
	clean := filepath.ToSlash(filePath)
	clean = strings.TrimPrefix(clean, ".")
	if !strings.HasPrefix(clean, "/") {
		clean = "/" + clean
	}
	if baseURL == "" {
		return clean
	}
	return strings.TrimRight(baseURL, "/") + clean
}
