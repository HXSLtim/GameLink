package resp

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/pkg/apierr"
)

// GetUserID extracts user_id from context. Returns 0 if not found.
func GetUserID(c *gin.Context) uint64 {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	userID, ok := userIDVal.(uint64)
	if !ok {
		return 0
	}
	return userID
}

// GetUserIDOrFail extracts user_id from context and responds with error if not found.
// Returns (id, ok). If ok is false, error response has been written.
func GetUserIDOrFail(c *gin.Context) (uint64, bool) {
	userID := GetUserID(c)
	if userID == 0 {
		Unauthorized(c, "missing user")
		return 0, false
	}
	return userID, true
}

// ParseUintParam parses a uint64 path parameter.
func ParseUintParam(c *gin.Context, name string) (uint64, error) {
	value := c.Param(name)
	return strconv.ParseUint(value, 10, 64)
}

// ParseIDOrFail parses an ID path parameter and responds with error if invalid.
// Returns (id, ok). If ok is false, error response has been written.
func ParseIDOrFail(c *gin.Context, paramName string) (uint64, bool) {
	idStr := c.Param(paramName)
	if idStr == "" {
		BadRequest(c, fmt.Sprintf("缺少%s参数", paramName))
		return 0, false
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		BadRequest(c, fmt.Sprintf("无效的%s格式", paramName))
		return 0, false
	}

	if id == 0 {
		BadRequest(c, fmt.Sprintf("%s不能为0", paramName))
		return 0, false
	}

	return id, true
}

// BindJSONOrFail binds JSON body and responds with error if invalid.
// Returns ok. If ok is false, error response has been written.
func BindJSONOrFail(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		Error(c, apierr.BadRequest(apierr.ErrInvalidJSONPayload).WithDetails(err.Error()))
		return false
	}
	return true
}

// BindQueryOrFail binds query parameters and responds with error if invalid.
// Returns ok. If ok is false, error response has been written.
func BindQueryOrFail(c *gin.Context, obj any) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		BadRequest(c, err.Error())
		return false
	}
	return true
}
