package payment

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

var (
	// ErrInvalidKeyFormat 表示密钥格式无效
	ErrInvalidKeyFormat = errors.New("invalid key format")
	// ErrInvalidPEMData 表示 PEM 数据无效
	ErrInvalidPEMData = errors.New("invalid PEM data")
	// ErrKeyNotFound 表示密钥文件不存在
	ErrKeyNotFound = errors.New("key file not found")
)

// loadPrivateKey 从 PEM 文件加载 RSA 私钥
// 支持两种格式：
// 1. PKCS#1 (RSA PRIVATE KEY)
// 2. PKCS#8 (PRIVATE KEY)
func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	if path == "" {
		return nil, ErrKeyNotFound
	}

	// 读取文件内容
	pemData, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", ErrKeyNotFound, path)
		}
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	// 解码 PEM 块
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, ErrInvalidPEMData
	}

	// 尝试解析为 PKCS#1 格式
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return key, nil
	}

	// 尝试解析为 PKCS#8 格式
	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse private key (tried PKCS#1 and PKCS#8)", ErrInvalidKeyFormat)
	}

	rsaKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("%w: not an RSA private key", ErrInvalidKeyFormat)
	}

	return rsaKey, nil
}

// loadPublicKey 从 PEM 文件加载 RSA 公钥
// 支持两种格式：
// 1. PKCS#1 (RSA PUBLIC KEY)
// 2. PKIX (PUBLIC KEY)
func loadPublicKey(path string) (*rsa.PublicKey, error) {
	if path == "" {
		return nil, ErrKeyNotFound
	}

	// 读取文件内容
	pemData, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", ErrKeyNotFound, path)
		}
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	// 解码 PEM 块
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, ErrInvalidPEMData
	}

	// 尝试解析为 PKCS#1 格式
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err == nil {
		return pub, nil
	}

	// 尝试解析为 PKIX 格式
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse public key (tried PKCS#1 and PKIX)", ErrInvalidKeyFormat)
	}

	rsaKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("%w: not an RSA public key", ErrInvalidKeyFormat)
	}

	return rsaKey, nil
}

// loadPrivateKeyFromBytes 从字节切片加载 RSA 私钥
// 用于从环境变量或配置直接加载密钥内容
func loadPrivateKeyFromBytes(pemData []byte) (*rsa.PrivateKey, error) {
	if len(pemData) == 0 {
		return nil, errors.New("empty private key data")
	}

	// 解码 PEM 块
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, ErrInvalidPEMData
	}

	// 尝试解析为 PKCS#1 格式
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return key, nil
	}

	// 尝试解析为 PKCS#8 格式
	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse private key", ErrInvalidKeyFormat)
	}

	rsaKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("%w: not an RSA private key", ErrInvalidKeyFormat)
	}

	return rsaKey, nil
}

// loadPublicKeyFromBytes 从字节切片加载 RSA 公钥
// 用于从环境变量或配置直接加载密钥内容
func loadPublicKeyFromBytes(pemData []byte) (*rsa.PublicKey, error) {
	if len(pemData) == 0 {
		return nil, errors.New("empty public key data")
	}

	// 解码 PEM 块
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, ErrInvalidPEMData
	}

	// 尝试解析为 PKCS#1 格式
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err == nil {
		return pub, nil
	}

	// 尝试解析为 PKIX 格式
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse public key", ErrInvalidKeyFormat)
	}

	rsaKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("%w: not an RSA public key", ErrInvalidKeyFormat)
	}

	return rsaKey, nil
}

// validatePrivateKey 验证私钥的有效性
func validatePrivateKey(key *rsa.PrivateKey) error {
	if key == nil {
		return errors.New("private key is nil")
	}

	// 检查密钥大小（建议至少 2048 位）
	if key.N.BitLen() < 2048 {
		return fmt.Errorf("private key too weak: %d bits (minimum 2048)", key.N.BitLen())
	}

	// 验证公钥
	if key.PublicKey.N == nil || key.D.Sign() == 0 {
		return errors.New("invalid private key: missing critical parameters")
	}

	return nil
}

// validatePublicKey 验证公钥的有效性
func validatePublicKey(key *rsa.PublicKey) error {
	if key == nil {
		return errors.New("public key is nil")
	}

	// 检查密钥大小（建议至少 2048 位）
	if key.N.BitLen() < 2048 {
		return fmt.Errorf("public key too weak: %d bits (minimum 2048)", key.N.BitLen())
	}

	// 验证 E 值（通常为 65537）
	if key.E <= 0 {
		return errors.New("invalid public key: invalid exponent")
	}

	return nil
}
