package auth

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"gamelink/internal/auth"
	"gamelink/internal/model"
	"gamelink/internal/repository"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// mockUserRepository 模拟用户仓库
type mockUserRepository struct {
	users map[uint64]*model.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[uint64]*model.User),
	}
}

func (m *mockUserRepository) Create(ctx context.Context, user *model.User) error {
	user.ID = uint64(len(m.users) + 1)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) Get(ctx context.Context, id uint64) (*model.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, repository.ErrNotFound
	}
	return user, nil
}


func (m *mockUserRepository) Update(ctx context.Context, user *model.User) error {
	if _, exists := m.users[user.ID]; !exists {
		return repository.ErrNotFound
	}
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id uint64) error {
	delete(m.users, id)
	return nil
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *mockUserRepository) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	for _, user := range m.users {
		if user.Phone == phone {
			return user, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *mockUserRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	return m.FindByPhone(ctx, phone)
}

func (m *mockUserRepository) List(ctx context.Context) ([]model.User, error) {
	var result []model.User
	for _, user := range m.users {
		result = append(result, *user)
	}
	return result, nil
}

func (m *mockUserRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	allUsers, _ := m.List(ctx)
	total := int64(len(allUsers))
	
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= len(allUsers) {
		return []model.User{}, total, nil
	}
	if end > len(allUsers) {
		end = len(allUsers)
	}
	
	return allUsers[start:end], total, nil
}

func (m *mockUserRepository) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
	var result []model.User
	for _, user := range m.users {
		if len(opts.Roles) > 0 {
			found := false
			for _, role := range opts.Roles {
				if user.Role == role {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if len(opts.Statuses) > 0 {
			found := false
			for _, status := range opts.Statuses {
				if user.Status == status {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		result = append(result, *user)
	}
	
	total := int64(len(result))
	start := (opts.Page - 1) * opts.PageSize
	end := start + opts.PageSize
	if start >= len(result) {
		return []model.User{}, total, nil
	}
	if end > len(result) {
		end = len(result)
	}
	
	return result[start:end], total, nil
}

// BenchmarkUserRegistration 用户注册性能基准测试
func BenchmarkUserRegistration(b *testing.B) {
	// 创建模拟仓库
	userRepo := newMockUserRepository()
	
	// 创建JWT管理器
	jwtManager := auth.NewJWTManager("benchmark-secret-key", 24*time.Hour)
	
	// 创建认证服务
	service := NewAuthService(userRepo, jwtManager)
	
	ctx := context.Background()
	
	// 重置基准测试计时器
	b.ResetTimer()
	
	// 运行基准测试
	for i := 0; i < b.N; i++ {
		req := RegisterRequest{
			Email:     fmt.Sprintf("benchmark%d@example.com", i),
			Password:  "Password123!",
			Name:      fmt.Sprintf("BenchmarkUser%d", i),
			Phone:     fmt.Sprintf("1380013800%02d", i%100),
		}
		
		_, err := service.Register(ctx, req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUserLogin 用户登录性能基准测试
func BenchmarkUserLogin(b *testing.B) {
	// 创建模拟仓库
	userRepo := newMockUserRepository()
	
	// 创建JWT管理器
	jwtManager := auth.NewJWTManager("benchmark-secret-key", 24*time.Hour)
	
	// 创建认证服务
	service := NewAuthService(userRepo, jwtManager)
	
	ctx := context.Background()
	
	// 预创建测试用户
	email := "benchmark@example.com"
	password := "Password123!"
	
	registerReq := RegisterRequest{
		Email:    email,
		Password: password,
		Name:     "BenchmarkUser",
		Phone:    "13800138000",
		Role:     model.RoleUser,
	}
	
	_, err := service.Register(ctx, registerReq)
	require.NoError(b, err)
	
	// 重置基准测试计时器
	b.ResetTimer()
	
	// 运行基准测试
	for i := 0; i < b.N; i++ {
		loginReq := LoginRequest{
			Username: email,
			Password: password,
		}
		
		_, err := service.Login(ctx, loginReq)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkTokenVerification Token验证性能基准测试
func BenchmarkTokenVerification(b *testing.B) {
	// 创建模拟仓库
	userRepo := newMockUserRepository()
	
	// 创建JWT管理器
	jwtManager := auth.NewJWTManager("benchmark-secret-key", 24*time.Hour)
	
	// 创建认证服务
	service := NewAuthService(userRepo, jwtManager)
	
	ctx := context.Background()
	
	// 预创建测试用户并获取token
	email := "benchmark@example.com"
	password := "Password123!"
	
	registerReq := RegisterRequest{
		Email:    email,
		Password: password,
		Name:     "BenchmarkUser",
		Phone:    "13800138000",
		Role:     model.RoleUser,
	}
	
	registerResp, err := service.Register(ctx, registerReq)
	require.NoError(b, err)
	
	token := registerResp.Token
	
	// 重置基准测试计时器
	b.ResetTimer()
	
	// 运行基准测试
	for i := 0; i < b.N; i++ {
		user, err := service.Me(ctx, "Bearer "+token)
		if err != nil {
			b.Fatal(err)
		}
		if user == nil {
			b.Fatal("User should not be nil")
		}
	}
}

// BenchmarkPasswordHashing 密码哈希性能基准测试
func BenchmarkPasswordHashing(b *testing.B) {
	password := "Password123!"
	
	// 重置基准测试计时器
	b.ResetTimer()
	
	// 运行基准测试
	for i := 0; i < b.N; i++ {
		_, err := hashPassword(password)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkEmailValidation 邮箱验证性能基准测试
func BenchmarkEmailValidation(b *testing.B) {
	emails := []string{
		"valid@example.com",
		"user.name@domain.co.uk",
		"test+tag@example.org",
		"invalid-email",
		"@invalid.com",
		"test@",
	}
	
	// 重置基准测试计时器
	b.ResetTimer()
	
	// 运行基准测试
	for i := 0; i < b.N; i++ {
		email := emails[i%len(emails)]
		// 简单的邮箱格式验证
		_ = strings.Contains(email, "@")
	}
}



// 辅助函数：哈希密码
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}