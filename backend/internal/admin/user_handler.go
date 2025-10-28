package admin

import (
    "errors"
    "net/http"
    "net/mail"
    "regexp"
    "strings"

    "github.com/gin-gonic/gin"

    apierr "gamelink/internal/handler"
    "gamelink/internal/model"
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
// @Param        page_size  query     int       false  "每页数量"
// @Param        role       query     []string  false  "角色过滤，可多值"
// @Param        status     query     []string  false  "状态过滤，可多值"
// @Param        date_from  query     string    false  "开始时间"
// @Param        date_to    query     string    false  "结束时间"
// @Param        keyword    query     string    false  "关键字（匹配 name/email/phone）"
// @Produce      json
// @Success      200  {object}  map[string]any
// @Router       /admin/users [get]
//
// ListUsers returns a paginated list of users.
func (h *UserHandler) ListUsers(c *gin.Context) {
    opts, ok := buildUserListOptions(c)
    if !ok { return }
    users, pagination, err := h.svc.ListUsersWithOptions(c.Request.Context(), opts)
    if err != nil {
        writeJSONError(c, http.StatusInternalServerError, err.Error())
        return
    }
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
    if err := h.svc.DeleteUser(c.Request.Context(), id); errors.Is(err, service.ErrNotFound) {
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
    var payload struct{ Status string `json:"status" binding:"required"` }
    if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
        writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
        return
    }
    user, err := h.svc.GetUser(c.Request.Context(), id)
    if errors.Is(err, service.ErrNotFound) {
        _ = c.Error(service.ErrNotFound)
        return
    }
    if err != nil {
        writeJSONError(c, http.StatusInternalServerError, err.Error())
        return
    }
    out, err := h.svc.UpdateUser(c.Request.Context(), id, service.UpdateUserInput{
        Phone:     user.Phone,
        Email:     user.Email,
        Name:      user.Name,
        AvatarURL: user.AvatarURL,
        Role:      user.Role,
        Status:    model.UserStatus(payload.Status),
    })
    if errors.Is(err, service.ErrValidation) { _ = c.Error(service.ErrValidation); return }
    if err != nil { writeJSONError(c, http.StatusInternalServerError, err.Error()); return }
    writeJSON(c, http.StatusOK, model.APIResponse[*model.User]{ Success: true, Code: http.StatusOK, Message: "updated", Data: out })
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
    var payload struct{ Role string `json:"role" binding:"required"` }
    if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
        writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
        return
    }
    user, err := h.svc.GetUser(c.Request.Context(), id)
    if errors.Is(err, service.ErrNotFound) { _ = c.Error(service.ErrNotFound); return }
    if err != nil { writeJSONError(c, http.StatusInternalServerError, err.Error()); return }
    out, err := h.svc.UpdateUser(c.Request.Context(), id, service.UpdateUserInput{
        Phone:     user.Phone,
        Email:     user.Email,
        Name:      user.Name,
        AvatarURL: user.AvatarURL,
        Role:      model.Role(payload.Role),
        Status:    user.Status,
    })
    if errors.Is(err, service.ErrValidation) { _ = c.Error(service.ErrValidation); return }
    if err != nil { writeJSONError(c, http.StatusInternalServerError, err.Error()); return }
    writeJSON(c, http.StatusOK, model.APIResponse[*model.User]{ Success: true, Code: http.StatusOK, Message: "updated", Data: out })
}

// ListUserOrders
// @Summary      获取用户的订单
// @Tags         Admin/Users
// @Security     BearerAuth
// @Produce      json
// @Param        id         path   int      true   "用户ID"
// @Param        page       query  int      false  "页码"
// @Param        page_size  query  int      false  "每页数量"
// @Param        status     query  []string false  "订单状态"
// @Param        date_from  query  string  false  "开始时间"
// @Param        date_to    query  string  false  "结束时间"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/users/{id}/orders [get]
func (h *UserHandler) ListUserOrders(c *gin.Context) {
    id, err := parseUintParam(c, "id")
    if err != nil { writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID); return }
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
    if !ok { return }
    // Force user scope
    uid := id
    opts.UserID = &uid
    orders, pagination, err := h.svc.ListOrders(c.Request.Context(), opts)
    if err != nil {
        writeJSONError(c, http.StatusInternalServerError, err.Error())
        return
    }
    writeJSON(c, http.StatusOK, model.APIResponse[[]model.Order]{ Success: true, Code: http.StatusOK, Message: "OK", Data: orders, Pagination: pagination })
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
