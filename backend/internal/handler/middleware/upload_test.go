package middleware

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestValidateFile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("文件大小超过限制应该失败", func(t *testing.T) {
		config := UploadConfig{
			MaxFileSize: 1024, // 1KB
		}

		// 创建一个2KB的文件
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.jpg")
		part.Write(make([]byte, 2048)) // 2KB
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.ParseMultipartForm(10 << 20)

		files := req.MultipartForm.File["file"]
		if len(files) > 0 {
			err := ValidateFile(files[0], config)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "exceeds maximum allowed size")
		}
	})

	t.Run("不允许的文件扩展名应该失败", func(t *testing.T) {
		config := UploadConfig{
			MaxFileSize:       10 * 1024 * 1024,
			AllowedExtensions: []string{".jpg", ".png"},
		}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.exe")
		part.Write([]byte("fake exe content"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.ParseMultipartForm(10 << 20)

		files := req.MultipartForm.File["file"]
		if len(files) > 0 {
			err := ValidateFile(files[0], config)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "extension")
		}
	})

	t.Run("允许的文件应该通过验证", func(t *testing.T) {
		config := UploadConfig{
			MaxFileSize:       10 * 1024 * 1024,
			AllowedExtensions: []string{".jpg", ".png"},
			AllowedMimeTypes:  []string{"image/jpeg", "image/png"},
		}

		// 创建一个简单的JPEG文件头
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.jpg")
		// JPEG文件魔数
		jpegHeader := []byte{0xFF, 0xD8, 0xFF, 0xE0}
		part.Write(jpegHeader)
		part.Write(make([]byte, 100))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.ParseMultipartForm(10 << 20)

		files := req.MultipartForm.File["file"]
		if len(files) > 0 {
			err := ValidateFile(files[0], config)
			assert.NoError(t, err)
		}
	})
}

func TestSaveFile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建临时上传目录
	tempDir := "./test_uploads"
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	t.Run("成功保存文件", func(t *testing.T) {
		config := UploadConfig{
			MaxFileSize:       10 * 1024 * 1024,
			AllowedExtensions: []string{".txt"},
			// 允许text/plain和其他可能的文本类型
			AllowedMimeTypes:  []string{"text/plain", "text/plain; charset=utf-8", "application/octet-stream"},
			UploadPath:        tempDir,
			RandomizeFilename: true,
			PreserveExtension: true,
			CalculateHash:     true,
		}

		router := gin.New()
		router.POST("/upload", func(c *gin.Context) {
			file, err := c.FormFile("file")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
				return
			}
			result, err := SaveFile(c, file, config)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, result)
		})

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.txt")
		part.Write([]byte("test content"))
		writer.Close()

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Response body: %s", w.Body.String())
		
		// 验证返回的结果
		if w.Code == http.StatusOK {
			var result UploadResult
			json.Unmarshal(w.Body.Bytes(), &result)
			assert.NotEmpty(t, result.SavedName)
			assert.NotEmpty(t, result.Hash)
		}
	})

	t.Run("保存文件时保留扩展名", func(t *testing.T) {
		config := UploadConfig{
			MaxFileSize:       10 * 1024 * 1024,
			AllowedExtensions: []string{".jpg"},
			AllowedMimeTypes:  []string{"image/jpeg"},
			UploadPath:        tempDir,
			RandomizeFilename: true,
			PreserveExtension: true,
			CalculateHash:     false,
		}

		router := gin.New()
		var savedPath string
		router.POST("/upload", func(c *gin.Context) {
			file, _ := c.FormFile("file")
			result, err := SaveFile(c, file, config)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			savedPath = result.SavedName
			c.JSON(http.StatusOK, result)
		})

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "photo.jpg")
		// JPEG文件头
		jpegHeader := []byte{0xFF, 0xD8, 0xFF, 0xE0}
		part.Write(jpegHeader)
		part.Write(make([]byte, 100))
		writer.Close()

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, ".jpg", filepath.Ext(savedPath))
	})
}

func TestGetImageConfig(t *testing.T) {
	config := GetImageConfig()
	
	assert.Equal(t, int64(5*1024*1024), config.MaxFileSize)
	assert.Contains(t, config.AllowedMimeTypes, "image/jpeg")
	assert.Contains(t, config.AllowedExtensions, ".jpg")
	assert.True(t, config.RandomizeFilename)
	assert.True(t, config.PreserveExtension)
	assert.True(t, config.CalculateHash)
}

func TestGetVideoConfig(t *testing.T) {
	config := GetVideoConfig()
	
	assert.Equal(t, int64(100*1024*1024), config.MaxFileSize)
	assert.Contains(t, config.AllowedMimeTypes, "video/mp4")
	assert.Contains(t, config.AllowedExtensions, ".mp4")
}

func TestGetAudioConfig(t *testing.T) {
	config := GetAudioConfig()
	
	assert.Equal(t, int64(20*1024*1024), config.MaxFileSize)
	assert.Contains(t, config.AllowedMimeTypes, "audio/mpeg")
	assert.Contains(t, config.AllowedExtensions, ".mp3")
}

func TestGetDocumentConfig(t *testing.T) {
	config := GetDocumentConfig()
	
	assert.Equal(t, int64(10*1024*1024), config.MaxFileSize)
	assert.Contains(t, config.AllowedMimeTypes, "application/pdf")
	assert.Contains(t, config.AllowedExtensions, ".pdf")
}

func TestUploadResult(t *testing.T) {
	result := &UploadResult{
		OriginalName: "test.jpg",
		SavedName:    "uuid-123.jpg",
		FilePath:     "/uploads/uuid-123.jpg",
		FileSize:     1024,
		MimeType:     "image/jpeg",
		Extension:    ".jpg",
		Hash:         "abc123",
	}

	assert.Equal(t, "test.jpg", result.OriginalName)
	assert.Equal(t, "uuid-123.jpg", result.SavedName)
	assert.Equal(t, int64(1024), result.FileSize)
	assert.Equal(t, "image/jpeg", result.MimeType)
	assert.Equal(t, ".jpg", result.Extension)
	assert.Equal(t, "abc123", result.Hash)
}
