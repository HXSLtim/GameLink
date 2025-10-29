package admin

import (
	"errors"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	apierr "gamelink/internal/handler"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service"
)

// UserHandler 处理后台用户管理接口。
type UserHandler struct {
	svc *service.AdminService
}

// NewUserHandler 创建 Handler。
func NewUserHandler(svc *service.AdminService) *UserHandler {
	return &UserHandler{svc: svc}
}

// ListUsers
// @Summary      列出用户
// @Description  分页列出用户，支持按角色/状态/时间范围/关键字筛选
// @Tags         Admin/Users
// @Security     BearerAuth
// @Param        page       query     int       false  "页码"
// @Param        pageSize   query     int       false  "每页数量"
// @Param        role       query     []string  false  "角色过滤，可多值"
// @Param        status     query     []string  false  "状态过滤，可多值"
// @Param        dateFrom   query     string    false  "开始时间"
// @Param        dateTo     query     string    false  "结束时间"
// @Param        keyword    query     string    false  "关键字（匹配 name/email/phone）"
// @Produce      json
// @Success      200  {object}  map[string]any
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
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	users = ensureSlice(users)
	writeJSON(c, http.StatusOK, model.APIResponse[[]model.User]{
		Success:    true,
		Code:       http.StatusOK,
		Message:    "OK",
		Data:       users,
		Pagination: pagination,
	})
}

// GetUser
// @Summary      获取用户
// @Tags         Admin/Users
// @Security     BearerAuth
// @Param        id   path      int  true  "用户ID"
// @Produce      json
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/users/{id} [get]
//
// GetUser returns a single user by id.
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	user, err := h.svc.GetUser(c.Request.Context(), id)
	if errors.Is(err, service.ErrUserNotFound) {
		_ = c.Error(service.ErrUserNotFound)
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[*model.User]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    user,
	})
}

// CreateUser
// @Summary      创建用户
// @Tags         Admin/Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  CreateUserPayload  true  "用户信息"
// @Success      201  {object}  map[string]any
// @Failure      400  {object}  map[string]any
// @Router       /admin/users [post]
//
// CreateUser creates a new user.
func (h *UserHandler) CreateUser(c *gin.Context) {
	var payload CreateUserPayload
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}

	if payload.Email != "" && !isValidEmail(payload.Email) {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidEmailFormat)
		return
	}
	if payload.Phone != "" && !isValidPhone(payload.Phone) {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidPhoneFormat)
		return
	}

	user, err := h.svc.CreateUser(c.Request.Context(), service.CreateUserInput{
		Phone:     payload.Phone,
		Email:     payload.Email,
		Password:  payload.Password,
		Name:      payload.Name,
		AvatarURL: payload.AvatarURL,
		Role:      model.Role(payload.Role),
		Status:    model.UserStatus(payload.Status),
	})
	if errors.Is(err, service.ErrValidation) {
		_ = c.Error(service.ErrValidation)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusCreated, model.APIResponse[*model.User]{
		Success: true,
		Code:    http.StatusCreated,
		Message: "created",
		Data:    user,
	})
}

// UpdateUser
// @Summary      更新用户
// @Tags         Admin/Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                  true  "用户ID"
// @Param        request  body  UpdateUserPayload    true  "用户信息"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/users/{id} [put]
//
// UpdateUser updates user profile and optional password.
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}

	var payload UpdateUserPayload
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}

	if payload.Email != "" && !isValidEmail(payload.Email) {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidEmailFormat)
		return
	}
	if payload.Phone != "" && !isValidPhone(payload.Phone) {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidPhoneFormat)
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

	user, err := h.svc.UpdateUser(c.Request.Context(), id, service.UpdateUserInput{
		Phone:     payload.Phone,
		Email:     payload.Email,
		Name:      payload.Name,
		AvatarURL: payload.AvatarURL,
		Role:      model.Role(payload.Role),
		Status:    model.UserStatus(payload.Status),
		Password:  passwordPtr,
	})
	if errors.Is(err, service.ErrValidation) {
		_ = c.Error(service.ErrValidation)
		return
	}
	if errors.Is(err, service.ErrUserNotFound) {
		_ = c.Error(service.ErrUserNotFound)
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[*model.User]{
		Success: true,
		Code:    http.StatusOK,
		Message: "updated",
		Data:    user,
	})
}

// DeleteUser
// @Summary      删除用户
// @Tags         Admin/Users
// @Security     BearerAuth
// @Param        id   path  int  true  "用户ID"
// @Produce      json
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/users/{id} [delete]
//
// DeleteUser deletes a user by id.
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	if err := h.svc.DeleteUser(c.Request.Context(), id); errors.Is(err, service.ErrUserNotFound) {
		_ = c.Error(service.ErrUserNotFound)
		return
	} else if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	} else if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "deleted",
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
// @Param        dateFrom   query     string    false  "开始时间"
// @Param        dateTo     query     string    false  "结束时间"
// @Param        export       query  string false "导出格式" Enums(csv)
// @Param        fields       query  string false "导出列（逗号分隔）"
// @Param        header_lang  query  string false "列头语言" Enums(en,zh)
// @Success      200  {object}  map[string]any
// @Router       /admin/users/{id}/logs [get]
func (h *UserHandler) ListUserLogs(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidID)
		return
	}
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}
	var actorID *uint64
	if v, err := queryUint64Ptr(c, "actor_user_id"); err == nil {
		actorID = v
	}
	var dateFrom, dateTo *time.Time
	if v, err := queryTimePtr(c, "date_from"); err == nil {
		dateFrom = v
	} else if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidDateFrom)
		return
	}
	if v, err := queryTimePtr(c, "date_to"); err == nil {
		dateTo = v
	} else if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidDateTo)
		return
	}
	opts := repository.OperationLogListOptions{Page: page, PageSize: pageSize, Action: strings.TrimSpace(c.Query("action")), ActorUserID: actorID, DateFrom: dateFrom, DateTo: dateTo}
	items, p, err := h.svc.ListOperationLogs(c.Request.Context(), "user", id, opts)
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}
	if strings.EqualFold(strings.TrimSpace(c.Query("export")), "csv") {
		exportOperationLogsCSV(c, "user", id, items)
		return
	}
	items = ensureSlice(items)
	writeJSON(c, 200, model.APIResponse[[]model.OperationLog]{Success: true, Code: 200, Message: "OK", Data: items, Pagination: p})
}

// UpdateUserStatus
// @Summary      更新用户状态
// @Tags         Admin/Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int  true  "用户ID"
// @Param        request  body  map[string]string  true  "{status}"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/users/{id}/status [put]
func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	var payload struct {
		Status string `json:"status" binding:"required"`
	}
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}
	out, err := h.svc.UpdateUserStatus(c.Request.Context(), id, model.UserStatus(payload.Status))
	if errors.Is(err, service.ErrValidation) {
		_ = c.Error(service.ErrValidation)
		return
	}
	if errors.Is(err, service.ErrUserNotFound) {
		_ = c.Error(service.ErrUserNotFound)
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[*model.User]{Success: true, Code: http.StatusOK, Message: "updated", Data: out})
}

// UpdateUserRole
// @Summary      更新用户角色
// @Tags         Admin/Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int  true  "用户ID"
// @Param        request  body  map[string]string  true  "{role}"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/users/{id}/role [put]
func (h *UserHandler) UpdateUserRole(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	var payload struct {
		Role string `json:"role" binding:"required"`
	}
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}
	out, err := h.svc.UpdateUserRole(c.Request.Context(), id, model.Role(payload.Role))
	if errors.Is(err, service.ErrValidation) {
		_ = c.Error(service.ErrValidation)
		return
	}
	if errors.Is(err, service.ErrUserNotFound) {
		_ = c.Error(service.ErrUserNotFound)
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[*model.User]{Success: true, Code: http.StatusOK, Message: "updated", Data: out})
}

// ListUserOrders
// @Summary      获取用户的订单
// @Tags         Admin/Users
// @Security     BearerAuth
// @Produce      json
// @Param        id         path   int      true   "用户ID"
// @Param        page       query  int      false  "页码"
// @Param        pageSize   query     int       false  "每页数量"
// @Param        status     query  []string false  "订单状态"
// @Param        dateFrom   query     string    false  "开始时间"
// @Param        dateTo     query     string    false  "结束时间"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/users/{id}/orders [get]
func (h *UserHandler) ListUserOrders(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	// Ensure user exists
	if _, err := h.svc.GetUser(c.Request.Context(), id); errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	} else if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
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
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	orders = ensureSlice(orders)
	writeJSON(c, http.StatusOK, model.APIResponse[[]model.Order]{Success: true, Code: http.StatusOK, Message: "OK", Data: orders, Pagination: pagination})
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
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}
	if payload.Email != "" && !isValidEmail(payload.Email) {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidEmailFormat)
		return
	}
	if payload.Phone != "" && !isValidPhone(payload.Phone) {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidPhoneFormat)
		return
	}

	user, player, err := h.svc.RegisterUserAndPlayer(c.Request.Context(),
		service.CreateUserInput{
			Phone:     payload.Phone,
			Email:     payload.Email,
			Password:  payload.Password,
			Name:      payload.Name,
			AvatarURL: payload.AvatarURL,
			Role:      model.Role(payload.Role),
			Status:    model.UserStatus(payload.Status),
		},
		service.CreatePlayerInput{
			Nickname:           payload.Player.Nickname,
			Bio:                payload.Player.Bio,
			HourlyRateCents:    payload.Player.HourlyRateCents,
			MainGameID:         payload.Player.MainGameID,
			VerificationStatus: model.VerificationStatus(payload.Player.VerificationStatus),
		},
	)
	if errors.Is(err, service.ErrValidation) {
		_ = c.Error(service.ErrValidation)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusCreated, model.APIResponse[map[string]any]{
		Success: true,
		Code:    http.StatusCreated,
		Message: "created",
		Data:    map[string]any{"user": user, "player": player},
	})
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
