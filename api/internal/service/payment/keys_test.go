package payment

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
	"path/filepath"
	"testing"
)

// generateTestKeyPair 生成测试用的 RSA 密钥对
func generateTestKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	// 生成 2048 位 RSA 密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

// savePEMFile 保存密钥到 PEM 文件
func savePEMFile(path string, block *pem.Block) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return pem.Encode(file, block)
}

// TestLoadPrivateKey 测试加载私钥
func TestLoadPrivateKey(t *testing.T) {
	// 生成测试密钥对
	privateKey, _, err := generateTestKeyPair()
	if err != nil {
		t.Fatalf("failed to generate test key pair: %v", err)
	}

	// 创建临时目录
	tempDir := t.TempDir()

	// 测试 PKCS#1 格式
	pkcs1Path := filepath.Join(tempDir, "pkcs1_key.pem")
	pkcs1Bytes := x509.MarshalPKCS1PrivateKey(privateKey)
	pkcs1Block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: pkcs1Bytes,
	}
	if err := savePEMFile(pkcs1Path, pkcs1Block); err != nil {
		t.Fatalf("failed to save PKCS#1 key: %v", err)
	}

	loadedKey, err := loadPrivateKey(pkcs1Path)
	if err != nil {
		t.Fatalf("failed to load PKCS#1 key: %v", err)
	}

	if err := validatePrivateKey(loadedKey); err != nil {
		t.Fatalf("PKCS#1 key validation failed: %v", err)
	}

	// 验证密钥正确性
	if loadedKey.N.Cmp(privateKey.N) != 0 {
		t.Error("loaded PKCS#1 key modulus mismatch")
	}

	// 测试 PKCS#8 格式
	pkcs8Path := filepath.Join(tempDir, "pkcs8_key.pem")
	pkcs8Bytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("failed to marshal PKCS#8 key: %v", err)
	}
	pkcs8Block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs8Bytes,
	}
	if err := savePEMFile(pkcs8Path, pkcs8Block); err != nil {
		t.Fatalf("failed to save PKCS#8 key: %v", err)
	}

	loadedKey8, err := loadPrivateKey(pkcs8Path)
	if err != nil {
		t.Fatalf("failed to load PKCS#8 key: %v", err)
	}

	if err := validatePrivateKey(loadedKey8); err != nil {
		t.Fatalf("PKCS#8 key validation failed: %v", err)
	}

	if loadedKey8.N.Cmp(privateKey.N) != 0 {
		t.Error("loaded PKCS#8 key modulus mismatch")
	}
}

// TestLoadPublicKey 测试加载公钥
func TestLoadPublicKey(t *testing.T) {
	// 生成测试密钥对
	_, publicKey, err := generateTestKeyPair()
	if err != nil {
		t.Fatalf("failed to generate test key pair: %v", err)
	}

	// 创建临时目录
	tempDir := t.TempDir()

	// 测试 PKCS#1 格式
	pkcs1Path := filepath.Join(tempDir, "pkcs1_pub.pem")
	pkcs1Bytes := x509.MarshalPKCS1PublicKey(publicKey)
	pkcs1Block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pkcs1Bytes,
	}
	if err := savePEMFile(pkcs1Path, pkcs1Block); err != nil {
		t.Fatalf("failed to save PKCS#1 public key: %v", err)
	}

	loadedKey, err := loadPublicKey(pkcs1Path)
	if err != nil {
		t.Fatalf("failed to load PKCS#1 public key: %v", err)
	}

	if err := validatePublicKey(loadedKey); err != nil {
		t.Fatalf("PKCS#1 public key validation failed: %v", err)
	}

	if loadedKey.N.Cmp(publicKey.N) != 0 {
		t.Error("loaded PKCS#1 public key modulus mismatch")
	}

	// 测试 PKIX 格式
	pkixPath := filepath.Join(tempDir, "pkix_pub.pem")
	pkixBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		t.Fatalf("failed to marshal PKIX public key: %v", err)
	}
	pkixBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pkixBytes,
	}
	if err := savePEMFile(pkixPath, pkixBlock); err != nil {
		t.Fatalf("failed to save PKIX public key: %v", err)
	}

	loadedKeyIX, err := loadPublicKey(pkixPath)
	if err != nil {
		t.Fatalf("failed to load PKIX public key: %v", err)
	}

	if err := validatePublicKey(loadedKeyIX); err != nil {
		t.Fatalf("PKIX public key validation failed: %v", err)
	}

	if loadedKeyIX.N.Cmp(publicKey.N) != 0 {
		t.Error("loaded PKIX public key modulus mismatch")
	}
}

// TestLoadPrivateKeyFromBytes 测试从字节加载私钥
func TestLoadPrivateKeyFromBytes(t *testing.T) {
	// 生成测试密钥对
	privateKey, _, err := generateTestKeyPair()
	if err != nil {
		t.Fatalf("failed to generate test key pair: %v", err)
	}

	// 测试 PKCS#1 格式
	pkcs1Bytes := x509.MarshalPKCS1PrivateKey(privateKey)
	pkcs1PEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: pkcs1Bytes,
	})

	loadedKey, err := loadPrivateKeyFromBytes(pkcs1PEM)
	if err != nil {
		t.Fatalf("failed to load PKCS#1 key from bytes: %v", err)
	}

	if err := validatePrivateKey(loadedKey); err != nil {
		t.Fatalf("PKCS#1 key validation failed: %v", err)
	}

	// 测试 PKCS#8 格式
	pkcs8Bytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("failed to marshal PKCS#8 key: %v", err)
	}
	pkcs8PEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs8Bytes,
	})

	loadedKey8, err := loadPrivateKeyFromBytes(pkcs8PEM)
	if err != nil {
		t.Fatalf("failed to load PKCS#8 key from bytes: %v", err)
	}

	if err := validatePrivateKey(loadedKey8); err != nil {
		t.Fatalf("PKCS#8 key validation failed: %v", err)
	}
}

// TestLoadPublicKeyFromBytes 测试从字节加载公钥
func TestLoadPublicKeyFromBytes(t *testing.T) {
	// 生成测试密钥对
	_, publicKey, err := generateTestKeyPair()
	if err != nil {
		t.Fatalf("failed to generate test key pair: %v", err)
	}

	// 测试 PKCS#1 格式
	pkcs1Bytes := x509.MarshalPKCS1PublicKey(publicKey)
	pkcs1PEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pkcs1Bytes,
	})

	loadedKey, err := loadPublicKeyFromBytes(pkcs1PEM)
	if err != nil {
		t.Fatalf("failed to load PKCS#1 public key from bytes: %v", err)
	}

	if err := validatePublicKey(loadedKey); err != nil {
		t.Fatalf("PKCS#1 public key validation failed: %v", err)
	}

	// 测试 PKIX 格式
	pkixBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		t.Fatalf("failed to marshal PKIX public key: %v", err)
	}
	pkixPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pkixBytes,
	})

	loadedKeyIX, err := loadPublicKeyFromBytes(pkixPEM)
	if err != nil {
		t.Fatalf("failed to load PKIX public key from bytes: %v", err)
	}

	if err := validatePublicKey(loadedKeyIX); err != nil {
		t.Fatalf("PKIX public key validation failed: %v", err)
	}
}

// TestLoadPrivateKeyErrors 测试私钥加载错误处理
func TestLoadPrivateKeyErrors(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr error
	}{
		{
			name:    "empty path",
			path:    "",
			wantErr: ErrKeyNotFound,
		},
		{
			name:    "non-existent file",
			path:    "/non/existent/path/key.pem",
			wantErr: ErrKeyNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := loadPrivateKey(tt.path)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}

	// 测试无效的 PEM 数据
	t.Run("invalid PEM data", func(t *testing.T) {
		tempDir := t.TempDir()
		invalidPath := filepath.Join(tempDir, "invalid.pem")

		if err := os.WriteFile(invalidPath, []byte("not valid PEM data"), 0644); err != nil {
			t.Fatalf("failed to write invalid PEM file: %v", err)
		}

		_, err := loadPrivateKey(invalidPath)
		if err == nil {
			t.Fatal("expected error for invalid PEM, got nil")
		}
		if err != ErrInvalidPEMData {
			t.Errorf("expected ErrInvalidPEMData, got %v", err)
		}
	})
}

// TestLoadPublicKeyErrors 测试公钥加载错误处理
func TestLoadPublicKeyErrors(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr error
	}{
		{
			name:    "empty path",
			path:    "",
			wantErr: ErrKeyNotFound,
		},
		{
			name:    "non-existent file",
			path:    "/non/existent/path/key.pem",
			wantErr: ErrKeyNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := loadPublicKey(tt.path)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

// TestValidatePrivateKey 测试私钥验证
func TestValidatePrivateKey(t *testing.T) {
	tests := []struct {
		name string
		key  *rsa.PrivateKey
		err  string
	}{
		{
			name: "nil key",
			key:  nil,
			err:  "private key is nil",
		},
		{
			name: "weak key (1024 bits)",
			key: &rsa.PrivateKey{
				PublicKey: rsa.PublicKey{
					N: big.NewInt(1).Lsh(big.NewInt(1), 1024),
					E: 65537,
				},
				D: big.NewInt(1),
			},
			err: "too weak",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePrivateKey(tt.key)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if tt.err != "" && err.Error() == "" {
				t.Errorf("expected error containing %q, got %v", tt.err, err)
			}
		})
	}
}

// TestValidatePublicKey 测试公钥验证
func TestValidatePublicKey(t *testing.T) {
	tests := []struct {
		name string
		key  *rsa.PublicKey
		err  string
	}{
		{
			name: "nil key",
			key:  nil,
			err:  "public key is nil",
		},
		{
			name: "weak key (1024 bits)",
			key: &rsa.PublicKey{
				N: big.NewInt(1).Lsh(big.NewInt(1), 1024),
				E: 65537,
			},
			err: "too weak",
		},
		{
			name: "invalid exponent",
			key: &rsa.PublicKey{
				N: big.NewInt(1).Lsh(big.NewInt(1), 2048),
				E: 0,
			},
			err: "invalid exponent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePublicKey(tt.key)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if tt.err != "" && err.Error() == "" {
				t.Errorf("expected error containing %q, got %v", tt.err, err)
			}
		})
	}
}
