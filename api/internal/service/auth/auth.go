package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service"
	"gamelink/pkg/apierr"
	"gamelink/pkg/auth"
)

var (
	ErrInvalidCredentials = service.ErrInvalidCredentials
	ErrUserDisabled       = service.ErrUserDisabled
	ErrNotFound           = service.ErrNotFound
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

// GetUser returns a user by id (for /auth/me endpoint).
func (s *AuthService) GetUser(ctx context.Context, id uint64) (*model.User, error) {
	return s.userRepo.Get(ctx, id)
}

// Me verifies Authorization header and returns current user.
func (s *AuthService) Me(ctx context.Context, authorizationHeader string) (*model.User, error) {
	if authorizationHeader == "" {
		return nil, apierr.Unauthorized("缺少认证头")
	}
	token, err := auth.ExtractTokenFromHeader(authorizationHeader)
	if err != nil {
		return nil, apierr.Unauthorized("提取Token失败")
	}
	claims, err := s.jwtManager.VerifyToken(token)
	if err != nil {
		return nil, apierr.Unauthorized("验证Token失败")
	}
	if auth.IsTokenExpired(claims) {
		return nil, apierr.Unauthorized("令牌已过期")
	}
	user, err := s.userRepo.Get(ctx, claims.UserID)
	if err != nil {
		return nil, apierr.Unauthorized("获取用户信息失败")
	}
	if user.Status != model.UserStatusActive {
		return nil, ErrUserDisabled
	}
	return user, nil
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
	Name     string     `json:"name"`     // 真实姓名（可选）
	Nickname string     `json:"nickname"` // 昵称（用于显示）
	Role     model.Role `json:"role"`
}

// PhoneCodeLoginRequest 手机验证码登录请求
type PhoneCodeLoginRequest struct {
	Phone    string `json:"phone"`
	Nickname string `json:"nickname,omitempty"`
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
		return nil, apierr.BadRequest("用户名和密码不能为空")
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
			return nil, ErrInvalidCredentials
		}
		return nil, apierr.BadRequest("查找用户失败")
	}

	// 检查用户状态
	if user.Status != model.UserStatusActive {
		return nil, ErrUserDisabled
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// 生成JWT Token
	token, err := s.jwtManager.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return nil, apierr.InternalError("生成Token失败").WithDetails(err.Error())
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
			return nil, apierr.Conflict("邮箱已被注册")
		} else if err != repository.ErrNotFound {
			return nil, apierr.BadRequest("检查邮箱失败")
		}
	}

	// 检查手机号是否已存在
	if req.Phone != "" {
		if _, err := s.userRepo.FindByPhone(ctx, req.Phone); err == nil {
			return nil, apierr.Conflict("手机号已被注册")
		} else if err != repository.ErrNotFound {
			return nil, apierr.BadRequest("检查手机号失败")
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apierr.InternalError("密码加密失败")
	}

	// 创建用户
	// 如果没有提供昵称，使用姓名作为昵称
	nickname := req.Nickname
	if nickname == "" {
		nickname = req.Name
	}

	user := &model.User{
		Phone:        req.Phone,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Name:         req.Name,
		Nickname:     nickname,
		Role:         req.Role,
		Status:       model.UserStatusActive, // 默认激活状态
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, apierr.InternalError("创建用户失败")
	}

	// 生成JWT Token
	token, err := s.jwtManager.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return nil, apierr.InternalError("生成Token失败").WithDetails(err.Error())
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
		return "", apierr.Unauthorized("验证Token失败")
	}

	// 检查用户是否仍然存在且处于激活状态
	user, err := s.userRepo.Get(ctx, claims.UserID)
	if err != nil {
		return "", apierr.Unauthorized("获取用户信息失败")
	}

	if user.Status != model.UserStatusActive {
		return "", ErrUserDisabled
	}

	// 刷新Token
	newToken, err := s.jwtManager.RefreshToken(claims)
	if err != nil {
		return "", apierr.Unauthorized("刷新Token失败")
	}

	return newToken, nil
}

// LoginOrRegisterByPhone verifies a phone login flow after code verification.
// If user does not exist, it will create a default active user account.
func (s *AuthService) LoginOrRegisterByPhone(ctx context.Context, req PhoneCodeLoginRequest) (*LoginResponse, error) {
	phone := strings.TrimSpace(req.Phone)
	if !isValidPhone(phone) {
		return nil, apierr.BadRequest("手机号格式错误")
	}

	user, err := s.userRepo.FindByPhone(ctx, phone)
	if err != nil {
		if err != repository.ErrNotFound {
			return nil, apierr.InternalError("查找用户失败").WithDetails(err.Error())
		}

		displayName := strings.TrimSpace(req.Nickname)
		if displayName == "" {
			displayName = "手机用户" + phone[len(phone)-4:]
		}

		passwordHash, hashErr := bcrypt.GenerateFromPassword([]byte(generateRandomPassword()), bcrypt.DefaultCost)
		if hashErr != nil {
			return nil, apierr.InternalError("创建账号失败").WithDetails(hashErr.Error())
		}

		newUser := &model.User{
			Phone:        phone,
			PasswordHash: string(passwordHash),
			Name:         displayName,
			Nickname:     displayName,
			Role:         model.RoleUser,
			Status:       model.UserStatusActive,
		}

		if createErr := s.userRepo.Create(ctx, newUser); createErr != nil {
			// 并发场景下手机号已被其他请求创建，回读即可
			if existed, findErr := s.userRepo.FindByPhone(ctx, phone); findErr == nil {
				user = existed
			} else {
				return nil, apierr.InternalError("创建账号失败").WithDetails(createErr.Error())
			}
		} else {
			user = newUser
		}
	}

	if user.Status != model.UserStatusActive {
		return nil, ErrUserDisabled
	}

	token, err := s.jwtManager.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return nil, apierr.InternalError("生成Token失败").WithDetails(err.Error())
	}

	now := time.Now()
	user.LastLoginAt = &now
	_ = s.userRepo.Update(ctx, user)

	return &LoginResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		User:      *user,
	}, nil
}

// validateRegisterInput 验证注册输入
func validateRegisterInput(req RegisterRequest) error {
	if req.Name == "" {
		return apierr.BadRequest("姓名不能为空")
	}
	if req.Email == "" && req.Phone == "" {
		return apierr.BadRequest("邮箱或手机号不能为空")
	}
	// 验证邮箱格式
	if req.Email != "" && !isValidEmail(req.Email) {
		return apierr.BadRequest("邮箱格式错误")
	}
	if req.Password == "" {
		return apierr.BadRequest("密码不能为空")
	}
	if len(req.Password) < 6 {
		return apierr.BadRequest("密码长度至少为6位")
	}
	if req.Role == "" {
		req.Role = model.RoleUser // 默认角色
	}
	return nil
}

// emailRegex 邮箱正则表达式
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

// isValidEmail 检查是否是有效的邮箱格式
func isValidEmail(email string) bool {
	if email == "" || len(email) > 128 {
		return false
	}

	// 基本格式验证
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	// 正则表达式验证
	if !emailRegex.MatchString(email) {
		return false
	}

	// 检查常见临时邮箱域名（可选）
	disposableDomains := []string{"tempmail.com", "10minutemail.com", "guerrillamail.com", "mailinator.com"}
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		domain := strings.ToLower(parts[1])
		for _, disposable := range disposableDomains {
			if domain == disposable || strings.HasSuffix(domain, "."+disposable) {
				return false // 拒绝临时邮箱
			}
		}
	}

	return true
}

func isValidPhone(phone string) bool {
	return phoneRegex.MatchString(strings.TrimSpace(phone))
}

func generateRandomPassword() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		// fallback
		return "PhoneLogin@123456"
	}
	return hex.EncodeToString(buf)
}
