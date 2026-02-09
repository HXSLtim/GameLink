package oss

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"

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
	client    *cos.Client
}

// Upload uploads file to Tencent COS
func (p *TencentCOSProvider) Upload(ctx context.Context, key string, reader io.Reader) (string, error) {
	client, err := p.getClient()
	if err != nil {
		return "", err
	}

	_, err = client.Object.Put(ctx, key, reader, nil)
	if err != nil {
		return "", err
	}

	return p.objectURL(key)
}

// Delete deletes file from Tencent COS
func (p *TencentCOSProvider) Delete(ctx context.Context, key string) error {
	client, err := p.getClient()
	if err != nil {
		return err
	}

	_, err = client.Object.Delete(ctx, key)
	return err
}

// GetSignedURL generates signed URL for Tencent COS
func (p *TencentCOSProvider) GetSignedURL(ctx context.Context, key string, expire time.Duration) (string, error) {
	client, err := p.getClient()
	if err != nil {
		return "", err
	}

	signedURL, err := client.Object.GetPresignedURL(ctx, http.MethodGet, key, p.SecretID, p.SecretKey, expire, nil)
	if err != nil {
		return "", err
	}
	return signedURL.String(), nil
}

func (p *TencentCOSProvider) getClient() (*cos.Client, error) {
	if p.client != nil {
		return p.client, nil
	}

	bucketURL, err := p.buildBucketURL()
	if err != nil {
		return nil, err
	}

	baseURL := &cos.BaseURL{BucketURL: bucketURL}
	httpClient := &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  strings.TrimSpace(p.SecretID),
			SecretKey: strings.TrimSpace(p.SecretKey),
		},
		Timeout: 60 * time.Second,
	}
	p.client = cos.NewClient(baseURL, httpClient)
	return p.client, nil
}

func (p *TencentCOSProvider) buildBucketURL() (*url.URL, error) {
	endpoint := strings.TrimSpace(p.Endpoint)
	if endpoint != "" {
		if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
			endpoint = "https://" + endpoint
		}
		return url.Parse(endpoint)
	}
	if p.Bucket == "" || p.Region == "" {
		return nil, fmt.Errorf("missing cos bucket or region")
	}
	return url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", p.Bucket, p.Region))
}

func (p *TencentCOSProvider) objectURL(key string) (string, error) {
	base, err := p.buildBucketURL()
	if err != nil {
		return "", err
	}
	u := *base
	cleanKey := strings.TrimPrefix(key, "/")
	if cleanKey != "" {
		u.Path = path.Join(strings.TrimSuffix(u.Path, "/"), cleanKey)
	}
	return u.String(), nil
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
