package service

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"gamelink/internal/auth"
	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// AuthService 认证服务
//
// 功能：
// 1. 用户登录验证
// 2. Token生成和验证
// 3. 用户注册
type AuthService struct {
	userRepo   repository.UserRepository
	jwtManager *auth.JWTManager
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo repository.UserRepository, jwtManager *auth.JWTManager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username"` // 用户名（可以是邮箱或手机号）
	Password string `json:"password"` // 密码
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string     `json:"token"`      // JWT Token
	ExpiresAt time.Time  `json:"expires_at"` // 过期时间
	User      model.User `json:"user"`       // 用户信息
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Phone    string     `json:"phone"`
	Email    string     `json:"email"`
	Password string     `json:"password"`
	Name     string     `json:"name"`
	Role     model.Role `json:"role"`
}

// Login 用户登录
//
// 参数：
// - ctx: 上下文
// - req: 登录请求
//
// 返回：
// - LoginResponse: 包含Token和用户信息
// - error: 错误信息
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// 验证输入
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("用户名和密码不能为空")
	}

	// 查找用户（通过邮箱或手机号）
	var user *model.User
	var err error

	// 先尝试通过邮箱查找
	if isValidEmail(req.Username) {
		user, err = s.userRepo.FindByEmail(ctx, req.Username)
	} else {
		// 通过手机号查找
		user, err = s.userRepo.FindByPhone(ctx, req.Username)
	}

	if err != nil {
		if err == repository.ErrNotFound {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 检查用户状态
	if user.Status != model.UserStatusActive {
		return nil, errors.New("用户账户已被禁用")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("密码错误")
	}

	// 生成JWT Token
	token, err := s.jwtManager.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	// 忽略更新时间错误，不影响登录流程
	_ = s.userRepo.Update(ctx, user)

	return &LoginResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 与JWT Token有效期一致
		User:      *user,
	}, nil
}

// Register 用户注册
//
// 参数：
// - ctx: 上下文
// - req: 注册请求
//
// 返回：
// - LoginResponse: 包含Token和用户信息
// - error: 错误信息
func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*LoginResponse, error) {
	// 验证输入
	if err := validateRegisterInput(req); err != nil {
		return nil, err
	}

	// 检查邮箱是否已存在
	if req.Email != "" {
		if _, err := s.userRepo.FindByEmail(ctx, req.Email); err == nil {
			return nil, errors.New("邮箱已被注册")
		} else if err != repository.ErrNotFound {
			return nil, err
		}
	}

	// 检查手机号是否已存在
	if req.Phone != "" {
		if _, err := s.userRepo.FindByPhone(ctx, req.Phone); err == nil {
			return nil, errors.New("手机号已被注册")
		} else if err != repository.ErrNotFound {
			return nil, err
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &model.User{
		Phone:        req.Phone,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Name:         req.Name,
		Role:         req.Role,
		Status:       model.UserStatusActive, // 默认激活状态
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// 生成JWT Token
	token, err := s.jwtManager.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		User:      *user,
	}, nil
}

// RefreshToken 刷新Token
func (s *AuthService) RefreshToken(ctx context.Context, tokenString string) (string, error) {
	// 验证当前Token
	claims, err := s.jwtManager.VerifyToken(tokenString)
	if err != nil {
		return "", err
	}

	// 检查用户是否仍然存在且处于激活状态
	user, err := s.userRepo.Get(ctx, claims.UserID)
	if err != nil {
		return "", err
	}

	if user.Status != model.UserStatusActive {
		return "", errors.New("用户账户已被禁用")
	}

	// 刷新Token
	newToken, err := s.jwtManager.RefreshToken(claims)
	if err != nil {
		return "", err
	}

	return newToken, nil
}

// validateRegisterInput 验证注册输入
func validateRegisterInput(req RegisterRequest) error {
	if req.Name == "" {
		return errors.New("姓名不能为空")
	}
	if req.Email == "" && req.Phone == "" {
		return errors.New("邮箱和手机号必须至少提供一个")
	}
	if req.Password == "" {
		return errors.New("密码不能为空")
	}
	if len(req.Password) < 6 {
		return errors.New("密码长度不能少于6位")
	}
	if req.Role == "" {
		req.Role = model.RoleUser // 默认角色
	}
	return nil
}

// isValidEmail 检查是否是有效的邮箱格式
func isValidEmail(email string) bool {
	// 简单的邮箱验证
	return len(email) > 3 && email[len(email)-4:] == ".com" ||
		len(email) > 3 && email[len(email)-3:] == ".cn" ||
		len(email) > 6 && email[len(email)-7:] == ".org.cn"
}
