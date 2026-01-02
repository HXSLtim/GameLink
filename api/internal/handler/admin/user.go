package admin

import (
	"net/mail"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	adminservice "gamelink/internal/service/admin"
	"gamelink/pkg/apierr"
)

// User 用户模型（类型别名）
type User = model.User

// UserHandler 处理后台用户管理接口
type UserHandler struct {
	svc *adminservice.AdminService
}

// NewUserHandler 创建Handler
func NewUserHandler(svc *adminservice.AdminService) *UserHandler {
	return &UserHandler{svc: svc}
}

// GetUserStats
// @Summary      获取用户统计数据
// @Description  获取用户总数、角色分布、状态分布等统计信息
// @Tags         Admin/Users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  model.SuccessResponse
// @Router       /admin/users/stats [get]
//
// GetUserStats returns user statistics.
func (h *UserHandler) GetUserStats(c *gin.Context) {
	stats, err := h.svc.GetUserStats(c.Request.Context())
	if err != nil {
		respondError(c, apierr.InternalError("获取用户统计失败").WithDetails(err.Error()))
		return
	}
	respondSuccess(c, stats)
}

// ListUsers
// @Summary      列出用户
// @Description  API endpoint// @Tags         Admin/Users
// @Security     BearerAuth
// @Param        page       query     int       false  "页码"
// @Param        pageSize   query     int       false  "每页数量"
// @Param        role           query    []string     false  "Role filter"// @Param        status         query    []string     false  "Status filter"// @Param        dateFrom       query    string       false  "Start date (YYYY-MM-DD)"// @Param        dateTo     query     string    false  "End date (YYYY-MM-DD)"
// @Param        keyword        query    string       false  "Parameter: keyword"// @Produce      json
// @Success      200  {object}  model.APIResponse[[]User]
// @Router       /admin/users [get]
//
// ListUsers returns a paginated list of users.
func (h *UserHandler) ListUsers(c *gin.Context) {
	opts, ok := buildUserListOptions(c)
	if !ok {
		return
	}
	users, pagination, err := h.svc.ListUsersWithOptions(c.Request.Context(), opts)
	if err != nil {
		respondError(c, apierr.InternalError("获取用户列表失败").WithDetails(err.Error()))
		return
	}
	respondList(c, users, pagination)
}

// GetUser
// @Summary      获取用户
// @Tags         Admin/Users
// @Security     BearerAuth
// @Param        id   path      int  true  "用户ID"
// @Produce      json
// @Success      200  {object}  model.APIResponse[User]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/users/{id} [get]
//
// GetUser returns a single user by id.
func (h *UserHandler) GetUser(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	user, err := h.svc.GetUser(c.Request.Context(), id)
	if err != nil {
		if apierr.IsNotFound(err) {
			respondError(c, err)
			return
		}
		respondError(c, apierr.InternalError("获取用户信息失败").WithDetails(err.Error()))
		return
	}
	respondSuccess(c, user)
}

// CreateUser
// @Summary      创建用户
// @Tags         Admin/Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  CreateUserPayload  true  "用户信息"
// @Success      201  {object}  model.APIResponse[User]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/users [post]
//
// CreateUser creates a new user.
func (h *UserHandler) CreateUser(c *gin.Context) {
	var payload CreateUserPayload
	if !ValidateAndRespond(c, &payload) {
		return
	}

	if payload.Email != "" && !isValidEmail(payload.Email) {
		respondBadRequest(c, apierr.ErrInvalidEmailFormat)
		return
	}
	if payload.Phone != "" && !isValidPhone(payload.Phone) {
		respondBadRequest(c, apierr.ErrInvalidPhoneFormat)
		return
	}

	user, err := h.svc.CreateUser(c.Request.Context(), adminservice.CreateUserInput{
		Phone:     payload.Phone,
		Email:     payload.Email,
		Password:  payload.Password,
		Name:      payload.Name,
		AvatarURL: payload.AvatarURL,
		Role:      model.Role(payload.Role),
		Status:    model.UserStatus(payload.Status),
	})
	if err != nil {
		if apierr.IsValidationError(err) {
			respondError(c, err)
			return
		}
		respondError(c, apierr.InternalError("创建用户失败").WithDetails(err.Error()))
		return
	}
	respondCreated(c, user)
}

// UpdateUser
// @Summary      更新用户
// @Tags         Admin/Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                  true  "用户ID"
// @Param        request  body  UpdateUserPayload    true  "用户信息"
// @Success      200  {object}  model.APIResponse[User]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/users/{id} [put]
//
// UpdateUser updates user profile and optional password.
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var payload UpdateUserPayload
	if !ValidateAndRespond(c, &payload) {
		return
	}

	if payload.Email != "" && !isValidEmail(payload.Email) {
		respondBadRequest(c, apierr.ErrInvalidEmailFormat)
		return
	}
	if payload.Phone != "" && !isValidPhone(payload.Phone) {
		respondBadRequest(c, apierr.ErrInvalidPhoneFormat)
		return
	}

	var passwordPtr *string
	if payload.Password != nil {
		trimmed := strings.TrimSpace(*payload.Password)
		if trimmed != "" {
			passwordCopy := trimmed
			passwordPtr = &passwordCopy
		}
	}

	user, err := h.svc.UpdateUser(c.Request.Context(), id, adminservice.UpdateUserInput{
		Phone:     payload.Phone,
		Email:     payload.Email,
		Name:      payload.Name,
		AvatarURL: payload.AvatarURL,
		Role:      model.Role(payload.Role),
		Status:    model.UserStatus(payload.Status),
		Password:  passwordPtr,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	respondUpdated(c, user)
}

// DeleteUser
// @Summary      删除用户
// @Tags         Admin/Users
// @Security     BearerAuth
// @Param        id   path  int  true  "用户ID"
// @Produce      json
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/users/{id} [delete]
//
// DeleteUser deletes a user by id.
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	err := h.svc.DeleteUser(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondDeleted(c)
}

// BatchDeleteUsers
// @Summary      批量删除用户
// @Tags         Admin/Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  map[string][]int  true  "{ids: [1,2,3]}"
// @Success      200  {object}  model.SuccessResponse
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/users/batch-delete [post]
//
// BatchDeleteUsers deletes multiple users by ids.
func (h *UserHandler) BatchDeleteUsers(c *gin.Context) {
	var payload struct {
		IDs []uint64 `json:"ids" binding:"required,min=1"`
	}
	if !ValidateAndRespond(c, &payload) {
		return
	}

	// 删除每个用户
	var deletedCount int
	var failedCount int
	for _, id := range payload.IDs {
		err := h.svc.DeleteUser(c.Request.Context(), id)
		if err != nil {
			failedCount++
			continue
		}
		deletedCount++
	}

	if failedCount > 0 {
		respondSuccessWithMsg(c, "部分删除成功", map[string]int{
			"deleted": deletedCount,
			"failed":  failedCount,
		})
		return
	}
	respondSuccessWithMsg(c, "批量删除成功", map[string]int{
		"deleted": deletedCount,
	})
}

// ListUserLogs
// @Summary      获取用户操作日志
// @Tags         Admin/Users
// @Security     BearerAuth
// @Produce      json
// @Param        id           path   int  true  "用户ID"
// @Param        page         query  int  false "页码"
// @Param        pageSize   query     int       false  "每页数量"
// @Param        action       query  string false "动作过滤" Enums(create,update,delete)
// @Param        actor_user_id query int   false "操作者用户ID"
// @Param        dateFrom       query    string       false  "Start date (YYYY-MM-DD)"// @Param        dateTo     query     string    false  "End date (YYYY-MM-DD)"
// @Param        export       query  string false "导出格式" Enums(csv)
// @Param        fields         query    string       false  "Export fields (comma separated)"// @Param        header_lang  query  string false "列头语言" Enums(en,zh)
// @Success      200  {object}  model.APIResponse[[]model.OperationLog]
// @Router       /admin/users/{id}/logs [get]
func (h *UserHandler) ListUserLogs(c *gin.Context) {
	handleOperationLogList(c, "user", h.svc.ListOperationLogs)
}

// UpdateUserStatus
// @Summary      更新用户状
// @Tags         Admin/Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int  true  "用户ID"
// @Param        request  body  map[string]string  true  "{status}"
// @Success      200  {object}  model.APIResponse[User]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/users/{id}/status [put]
func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	var payload struct {
		Status string `json:"status" binding:"required"`
	}
	if !ValidateAndRespond(c, &payload) {
		return
	}
	out, err := h.svc.UpdateUserStatus(c.Request.Context(), id, model.UserStatus(payload.Status))
	if err != nil {
		respondError(c, err)
		return
	}
	respondUpdated(c, out)
}

// UpdateUserRole
// @Summary      更新用户角色
// @Tags         Admin/Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int  true  "用户ID"
// @Param        request  body  map[string]string  true  "{role}"
// @Success      200  {object}  model.APIResponse[User]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/users/{id}/role [put]
func (h *UserHandler) UpdateUserRole(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	var payload struct {
		Role string `json:"role" binding:"required"`
	}
	if !ValidateAndRespond(c, &payload) {
		return
	}
	out, err := h.svc.UpdateUserRole(c.Request.Context(), id, model.Role(payload.Role))
	if err != nil {
		respondError(c, err)
		return
	}
	respondUpdated(c, out)
}

// ListUserOrders
// @Summary      获取用户的订
// @Tags         Admin/Users
// @Security     BearerAuth
// @Produce      json
// @Param        id         path   int      true   "用户ID"
// @Param        page       query  int      false  "页码"
// @Param        pageSize   query     int       false  "每页数量"
// @Param        status         query    []string     false  "Status filter"// @Param        dateFrom       query    string       false  "Start date (YYYY-MM-DD)"// @Param        dateTo     query     string    false  "End date (YYYY-MM-DD)"
// @Success      200  {object}  model.APIResponse[[]model.Order]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/users/{id}/orders [get]
func (h *UserHandler) ListUserOrders(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	// Ensure user exists
	_, err := h.svc.GetUser(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	// Build options from query, then enforce user scope
	opts, ok := buildOrderListOptions(c)
	if !ok {
		return
	}
	// Force user scope
	uid := id
	opts.UserID = &uid
	orders, pagination, err := h.svc.ListOrders(c.Request.Context(), opts)
	if err != nil {
		respondError(c, apierr.InternalError("获取用户订单失败").WithDetails(err.Error()))
		return
	}
	respondList(c, orders, pagination)
}

// CreateUserPayload defines the request body for creating a user.
type CreateUserPayload struct {
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Password  string `json:"password" binding:"required,min=6"`
	Name      string `json:"name" binding:"required"`
	AvatarURL string `json:"avatar_url"`
	Role      string `json:"role" binding:"required"`
	Status    string `json:"status" binding:"required"`
}

// UpdateUserPayload defines the request body for updating a user.
type UpdateUserPayload struct {
	Phone     string  `json:"phone"`
	Email     string  `json:"email"`
	Name      string  `json:"name" binding:"required"`
	AvatarURL string  `json:"avatar_url"`
	Role      string  `json:"role" binding:"required"`
	Status    string  `json:"status" binding:"required"`
	Password  *string `json:"password"`
}

// CreateUserWithPlayerPayload defines a combined create payload.
type CreateUserWithPlayerPayload struct {
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Password  string `json:"password" binding:"required,min=6"`
	Name      string `json:"name" binding:"required"`
	AvatarURL string `json:"avatar_url"`
	Role      string `json:"role" binding:"required"`
	Status    string `json:"status" binding:"required"`

	Player struct {
		Nickname           string `json:"nickname"`
		Bio                string `json:"bio"`
		HourlyRateCents    int64  `json:"hourly_rate_cents"`
		MainGameID         uint64 `json:"main_game_id"`
		VerificationStatus string `json:"verification_status" binding:"required"`
	} `json:"player" binding:"required"`
}

// CreateUserWithPlayer creates a user and a player profile transactionally.
func (h *UserHandler) CreateUserWithPlayer(c *gin.Context) {
	var payload CreateUserWithPlayerPayload
	if !ValidateAndRespond(c, &payload) {
		return
	}
	if payload.Email != "" && !isValidEmail(payload.Email) {
		respondBadRequest(c, apierr.ErrInvalidEmailFormat)
		return
	}
	if payload.Phone != "" && !isValidPhone(payload.Phone) {
		respondBadRequest(c, apierr.ErrInvalidPhoneFormat)
		return
	}

	user, player, err := h.svc.RegisterUserAndPlayer(c.Request.Context(),
		adminservice.CreateUserInput{
			Phone:     payload.Phone,
			Email:     payload.Email,
			Password:  payload.Password,
			Name:      payload.Name,
			AvatarURL: payload.AvatarURL,
			Role:      model.Role(payload.Role),
			Status:    model.UserStatus(payload.Status),
		},
		adminservice.CreatePlayerInput{
			Nickname:           payload.Player.Nickname,
			Bio:                payload.Player.Bio,
			HourlyRateCents:    payload.Player.HourlyRateCents,
			MainGameID:         payload.Player.MainGameID,
			VerificationStatus: model.VerificationStatus(payload.Player.VerificationStatus),
		},
	)
	if err != nil {
		respondError(c, err)
		return
	}
	respondCreated(c, map[string]any{"user": user, "player": player})
}

var phoneRegexp = regexp.MustCompile(`^1[3-9]\d{9}$`)

func isValidPhone(p string) bool {
	p = strings.TrimSpace(p)
	return phoneRegexp.MatchString(p)
}

func isValidEmail(e string) bool {
	e = strings.TrimSpace(e)
	if e == "" {
		return false
	}
	if _, err := mail.ParseAddress(e); err != nil {
		return false
	}
	return true
}

// ListUserLoginHistory
// @Summary      获取用户登录历史
// @Description  获取指定用户的登录历史记录（分页）
// @Tags         Admin/Users
// @Security     BearerAuth
// @Produce      json
// @Param        id          path   int      true   "用户ID"
// @Param        page        query  int      false  "页码"
// @Param        page_size   query  int      false  "每页数量"
// @Success      200  {object}  model.APIResponse[[]model.UserLoginHistory]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/users/{id}/login-history [get]
func (h *UserHandler) ListUserLoginHistory(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	// 确认用户存在
	_, err := h.svc.GetUser(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	// User login history will be implemented when UserLoginHistoryRepository is integrated
	// For now, return empty data
	histories := []model.UserLoginHistory{}
	total := int64(0)

	// 构建分页响应
	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}
	pagination := &model.Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
	respondList(c, histories, pagination)
}
