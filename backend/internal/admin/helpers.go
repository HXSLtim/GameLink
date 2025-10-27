package admin

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
)

func parseUintParam(c *gin.Context, key string) (uint64, error) {
	return strconv.ParseUint(c.Param(key), 10, 64)
}

func queryIntDefault(c *gin.Context, key string, defaults int) (int, error) {
	value := strings.TrimSpace(c.Query(key))
	if value == "" {
		return defaults, nil
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func queryUint64Ptr(c *gin.Context, key string) (*uint64, error) {
	value := strings.TrimSpace(c.Query(key))
	if value == "" {
		return nil, nil
	}
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func queryTimePtr(c *gin.Context, key string) (*time.Time, error) {
	value := strings.TrimSpace(c.Query(key))
	if value == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func parseCSVParams(values []string) []string {
	var result []string
	for _, raw := range values {
		parts := strings.Split(raw, ",")
		for _, part := range parts {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				result = append(result, trimmed)
			}
		}
	}
	return result
}

func writeJSON[T any](c *gin.Context, status int, payload model.APIResponse[T]) {
	c.JSON(status, payload)
}

func writeJSONError(c *gin.Context, status int, message string) {
	writeJSON(c, status, model.APIResponse[any]{
		Success: false,
		Code:    status,
		Message: message,
	})
}

// normalizeOrderStatus maps legacy spellings to canonical values.
// Accepts "cancelled" (legacy) and returns "canceled".
func normalizeOrderStatus(s string) model.OrderStatus { //nolint:misspell // accepts legacy 'cancelled'
	v := strings.TrimSpace(strings.ToLower(s))
	switch v {
	case "cancelled": // legacy spelling
		return model.OrderStatusCanceled
	default:
		return model.OrderStatus(v)
	}
}
