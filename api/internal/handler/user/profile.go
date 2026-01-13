package user

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
)

// ProfileHandler 用户资料处理器
type ProfileHandler struct {
	userRepo repository.UserRepository
}

// NewProfileHandler 创建用户资料处理器
func NewProfileHandler(userRepo repository.UserRepository) *ProfileHandler {
	return &ProfileHandler{userRepo: userRepo}
}

// UpdateProfileRequest 更新资料请求
type UpdateProfileRequest struct {
	Name      string `json:"name" binding:"omitempty,min=2,max=32"`
	AvatarURL string `json:"avatarUrl" binding:"omitempty,url,max=255"`
}

// UpdateProfileResponse 更新资料响应
type UpdateProfileResponse struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatarUrl"`
	UpdatedAt string `json:"updatedAt"`
}

// UserProfileResponse 用户资料响应
type UserProfileResponse struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	AvatarURL string `json:"avatarUrl"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// UpdateProfile 更新用户资料
// @Summary      更新用户资料
// @Description  更新当前用户的个人资料（昵称、头像）
// @Tags         User - Profile
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      UpdateProfileRequest  true  "更新信息"
// @Success      200      {object}  UpdateProfileResponse
// @Failure      400      {object}  apierr.APIError
// @Failure      401      {object}  apierr.APIError
// @Failure      500      {object}  apierr.APIError
// @Router       /user/profile [put]
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("无效的请求参数").WithDetails(err.Error()))
		return
	}

	// 获取当前用户
	user, err := h.userRepo.Get(c.Request.Context(), userID)
	if err != nil {
		resp.Error(c, apierr.InternalError("获取用户信息失败").WithDetails(err.Error()))
		return
	}

	// 更新字段
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	// 保存更新
	if err := h.userRepo.Update(c.Request.Context(), user); err != nil {
		resp.Error(c, apierr.InternalError("更新用户资料失败").WithDetails(err.Error()))
		return
	}

	resp.OK(c, UpdateProfileResponse{
		ID:        user.ID,
		Name:      user.Name,
		AvatarURL: user.AvatarURL,
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// GetProfile 获取用户资料
// @Summary      获取用户资料
// @Description  获取当前用户的个人资料
// @Tags         User - Profile
// @Security     BearerAuth
// @Produce      json
// @Success      200      {object}  UserProfileResponse
// @Failure      401      {object}  apierr.APIError
// @Failure      500      {object}  apierr.APIError
// @Router       /user/profile [get]
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	user, err := h.userRepo.Get(c.Request.Context(), userID)
	if err != nil {
		resp.Error(c, apierr.InternalError("获取用户信息失败").WithDetails(err.Error()))
		return
	}

	resp.OK(c, UserProfileResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		AvatarURL: user.AvatarURL,
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// RegisterProfileRoutes 注册用户资料路由
func RegisterProfileRoutes(rg *gin.RouterGroup, userRepo repository.UserRepository, _ gin.HandlerFunc) {
	h := NewProfileHandler(userRepo)

	rg.GET("/profile", h.GetProfile)
	rg.PUT("/profile", h.UpdateProfile)
}
