package admin

import (
	"errors"
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

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

// ListUsers returns a paginated list of users.
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, err := queryIntDefault(c, "page", 1)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid page")
		return
	}
	pageSize, err := queryIntDefault(c, "page_size", 20)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid page_size")
		return
	}

	users, pagination, err := h.svc.ListUsersPaged(c.Request.Context(), page, pageSize)
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

// GetUser returns a single user by id.
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	user, err := h.svc.GetUser(c.Request.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		writeJSONError(c, http.StatusNotFound, "user not found")
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

// CreateUser creates a new user.
func (h *UserHandler) CreateUser(c *gin.Context) {
	var payload CreateUserPayload
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid json payload")
		return
	}

	if payload.Email != "" && !isValidEmail(payload.Email) {
		writeJSONError(c, http.StatusBadRequest, "invalid email format")
		return
	}
	if payload.Phone != "" && !isValidPhone(payload.Phone) {
		writeJSONError(c, http.StatusBadRequest, "invalid phone format")
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
		writeJSONError(c, http.StatusBadRequest, "缺少必填字段或密码过短")
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

// UpdateUser updates user profile and optional password.
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var payload UpdateUserPayload
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid json payload")
		return
	}

	if payload.Email != "" && !isValidEmail(payload.Email) {
		writeJSONError(c, http.StatusBadRequest, "invalid email format")
		return
	}
	if payload.Phone != "" && !isValidPhone(payload.Phone) {
		writeJSONError(c, http.StatusBadRequest, "invalid phone format")
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
		writeJSONError(c, http.StatusBadRequest, "缺少必填字段或密码过短")
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		writeJSONError(c, http.StatusNotFound, "user not found")
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

// DeleteUser deletes a user by id.
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteUser(c.Request.Context(), id); errors.Is(err, service.ErrNotFound) {
		writeJSONError(c, http.StatusNotFound, "user not found")
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
