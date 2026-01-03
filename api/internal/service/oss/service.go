package oss

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"gamelink/internal/service/external"
)

// Provider defines OSS provider interface
type Provider interface {
	Upload(ctx context.Context, key string, reader io.Reader) (string, error)
	Delete(ctx context.Context, key string) error
	GetSignedURL(ctx context.Context, key string, expire time.Duration) (string, error)
}

// Service handles OSS operations
type Service struct {
	config *external.Config
}

// NewService creates OSS service
func NewService(cfg *external.Config) *Service {
	return &Service{config: cfg}
}

var (
	// ErrOSSDisabled OSS service is disabled
	ErrOSSDisabled = fmt.Errorf("oss service is disabled")
	// ErrUnsupportedProvider unsupported OSS provider
	ErrUnsupportedProvider = fmt.Errorf("unsupported oss provider")
)

// Upload uploads file to OSS
func (s *Service) Upload(ctx context.Context, key string, reader io.Reader) (string, error) {
	if !s.config.OSS.Enabled {
		return "", ErrOSSDisabled
	}

	provider, err := s.getProvider()
	if err != nil {
		return "", err
	}

	return provider.Upload(ctx, key, reader)
}

// Delete deletes file from OSS
func (s *Service) Delete(ctx context.Context, key string) error {
	if !s.config.OSS.Enabled {
		return ErrOSSDisabled
	}

	provider, err := s.getProvider()
	if err != nil {
		return err
	}

	return provider.Delete(ctx, key)
}

// GetSignedURL generates signed URL for file access
func (s *Service) GetSignedURL(ctx context.Context, key string, expire time.Duration) (string, error) {
	if !s.config.OSS.Enabled {
		return "", ErrOSSDisabled
	}

	provider, err := s.getProvider()
	if err != nil {
		return "", err
	}

	return provider.GetSignedURL(ctx, key, expire)
}

func (s *Service) getProvider() (Provider, error) {
	switch s.config.OSS.Provider {
	case "qcloud":
		return &TencentCOSProvider{
			SecretID:  s.config.OSS.AccessKey,
			SecretKey: s.config.OSS.SecretKey,
			Bucket:    s.config.OSS.Bucket,
			Region:    s.config.OSS.Region,
			Endpoint:  s.config.OSS.Endpoint,
		}, nil
	case "aliyun":
		return &AliyunOSSProvider{
			AccessKey: s.config.OSS.AccessKey,
			SecretKey: s.config.OSS.SecretKey,
			Bucket:    s.config.OSS.Bucket,
			Region:    s.config.OSS.Region,
			Endpoint:  s.config.OSS.Endpoint,
		}, nil
	case "minio":
		return &MinIOProvider{
			Endpoint:  s.config.OSS.Endpoint,
			AccessKey: s.config.OSS.AccessKey,
			SecretKey: s.config.OSS.SecretKey,
			Bucket:    s.config.OSS.Bucket,
		}, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedProvider, s.config.OSS.Provider)
	}
}

// TencentCOSProvider implements OSS provider for Tencent Cloud COS
type TencentCOSProvider struct {
	SecretID  string
	SecretKey string
	Bucket    string
	Region    string
	Endpoint  string
}

// Upload uploads file to Tencent COS
func (p *TencentCOSProvider) Upload(ctx context.Context, key string, reader io.Reader) (string, error) {
	// Build COS endpoint
	// Format: https://{bucket-name}-{appid}.cos.{region}.myqcloud.com/{key}
	bucketURL := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", p.Bucket, p.Region)
	fullURL := fmt.Sprintf("%s/%s", bucketURL, key)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "PUT", fullURL, reader)
	if err != nil {
		return "", err
	}

	// Add authorization header
	auth := p.generateAuth("PUT", key, "")
	req.Header.Set("Authorization", auth)
	req.Header.Set("Host", fmt.Sprintf("%s.cos.%s.myqcloud.com", p.Bucket, p.Region))

	// Execute request with timeout
	// Timeout choice: 60s for file upload operations (can be large files)
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Return file URL
	return fullURL, nil
}

// Delete deletes file from Tencent COS
func (p *TencentCOSProvider) Delete(ctx context.Context, key string) error {
	bucketURL := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", p.Bucket, p.Region)
	fullURL := fmt.Sprintf("%s/%s", bucketURL, key)

	req, err := http.NewRequestWithContext(ctx, "DELETE", fullURL, nil)
	if err != nil {
		return err
	}

	auth := p.generateAuth("DELETE", key, "")
	req.Header.Set("Authorization", auth)
	req.Header.Set("Host", fmt.Sprintf("%s.cos.%s.myqcloud.com", p.Bucket, p.Region))

	// Timeout choice: 10s for delete operations (should be fast)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete failed with status %d", resp.StatusCode)
	}

	return nil
}

// GetSignedURL generates signed URL for Tencent COS
func (p *TencentCOSProvider) GetSignedURL(ctx context.Context, key string, expire time.Duration) (string, error) {
	bucketURL := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", p.Bucket, p.Region)
	fullURL := fmt.Sprintf("%s/%s", bucketURL, key)

	// Add signed parameters
	expiredTime := time.Now().Add(expire).Unix()

	// Generate signature
	keyTime := fmt.Sprintf("%d;%d", time.Now().Unix(), expiredTime)

	// COS signature generation will be implemented for production
	// Reference: https://cloud.tencent.com/document/product/436/7778
	// Requires: HMAC-SHA1, q-sign-algorithm, q-ak, q-sign-time headers

	u, _ := url.Parse(fullURL)
	q := u.Query()
	q.Set("sign", p.generateSign(key, keyTime))
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// generateAuth generates Tencent COS authorization header
func (p *TencentCOSProvider) generateAuth(method, key, signTime string) string {
	// COS authorization will be implemented for production
	// Reference: https://cloud.tencent.com/document/product/436/7778
	requestTime := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	return fmt.Sprintf("q-sign-algorithm=sha1&q-ak=%s&q-sign-time=%s&q-key-time=%s&q-header-list=&q-url-param-list=&q-signature=",
		p.SecretID, requestTime, requestTime)
}

// generateSign generates signature for Tencent COS
func (p *TencentCOSProvider) generateSign(key, signTime string) string {
	// Simple implementation
	h := hmac.New(sha1.New, []byte(p.SecretKey))
	h.Write([]byte(key + signTime))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// AliyunOSSProvider implements OSS provider for Aliyun OSS
type AliyunOSSProvider struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
	Endpoint  string
}

// Upload uploads file to Aliyun OSS
func (p *AliyunOSSProvider) Upload(ctx context.Context, key string, reader io.Reader) (string, error) {
	// Aliyun OSS upload will be implemented for production
	// Reference: https://help.aliyun.com/document_detail/32016.html
	// Requires: SDK for Go (github.com/aliyun/aliyun-oss-go-sdk)
	fmt.Printf("[Aliyun OSS] Upload: bucket=%s, key=%s\n", p.Bucket, key)
	return "", fmt.Errorf("not implemented")
}

// Delete deletes file from Aliyun OSS
func (p *AliyunOSSProvider) Delete(ctx context.Context, key string) error {
	fmt.Printf("[Aliyun OSS] Delete: bucket=%s, key=%s\n", p.Bucket, key)
	return fmt.Errorf("not implemented")
}

// GetSignedURL generates signed URL for Aliyun OSS
func (p *AliyunOSSProvider) GetSignedURL(ctx context.Context, key string, expire time.Duration) (string, error) {
	fmt.Printf("[Aliyun OSS] GetSignedURL: bucket=%s, key=%s, expire=%v\n", p.Bucket, key, expire)
	return "", fmt.Errorf("not implemented")
}

// MinIOProvider implements OSS provider for MinIO
type MinIOProvider struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
}

// Upload uploads file to MinIO
func (p *MinIOProvider) Upload(ctx context.Context, key string, reader io.Reader) (string, error) {
	// MinIO upload will be implemented for production
	// Requires: minio-go SDK (github.com/minio/minio-go/v7)
	fmt.Printf("[MinIO] Upload: endpoint=%s, bucket=%s, key=%s\n", p.Endpoint, p.Bucket, key)
	return "", fmt.Errorf("not implemented")
}

// Delete deletes file from MinIO
func (p *MinIOProvider) Delete(ctx context.Context, key string) error {
	fmt.Printf("[MinIO] Delete: endpoint=%s, bucket=%s, key=%s\n", p.Endpoint, p.Bucket, key)
	return fmt.Errorf("not implemented")
}

// GetSignedURL generates signed URL for MinIO
func (p *MinIOProvider) GetSignedURL(ctx context.Context, key string, expire time.Duration) (string, error) {
	fmt.Printf("[MinIO] GetSignedURL: endpoint=%s, bucket=%s, key=%s, expire=%v\n", p.Endpoint, p.Bucket, key, expire)
	return "", fmt.Errorf("not implemented")
}
