package admin

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	apierr "gamelink/internal/handler"
	"gamelink/internal/model"
	"gamelink/internal/repository"
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
	// Try common formats: RFC3339, "2006-01-02 15:04:05", "2006-01-02", unix seconds
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return &t, nil
	}
	layouts := []string{"2006-01-02 15:04:05", "2006-01-02"}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return &t, nil
		}
	}
	if sec, err := strconv.ParseInt(value, 10, 64); err == nil {
		t := time.Unix(sec, 0)
		return &t, nil
	}
	return nil, errors.New("invalid time format")
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

func ensureSlice[T any](items []T) []T {
	if items == nil {
		return make([]T, 0)
	}
	return items
}

// parsePagination parses page and page_size with defaults and writes error response when invalid.
func parsePagination(c *gin.Context) (int, int, bool) {
	page, err := queryIntDefault(c, "page", 1)
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidPage)
		return 0, 0, false
	}
	pageSize, err := queryIntDefault(c, "page_size", 20)
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidPageSize)
		return 0, 0, false
	}
	return page, pageSize, true
}

// buildOrderListOptions parses query parameters into OrderListOptions; on error responds and returns false.
func buildOrderListOptions(c *gin.Context) (repository.OrderListOptions, bool) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return repository.OrderListOptions{}, false
	}

	statusTokens := parseCSVParams(c.QueryArray("status"))
	statuses := make([]model.OrderStatus, 0, len(statusTokens))
	for _, token := range statusTokens {
		statuses = append(statuses, normalizeOrderStatus(token))
	}

	userID, err := queryUint64Ptr(c, "user_id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidUserID)
		return repository.OrderListOptions{}, false
	}
	playerID, err := queryUint64Ptr(c, "player_id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidPlayerID)
		return repository.OrderListOptions{}, false
	}
	gameID, err := queryUint64Ptr(c, "game_id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidGameID)
		return repository.OrderListOptions{}, false
	}
	dateFrom, err := queryTimePtr(c, "date_from")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidDateFrom)
		return repository.OrderListOptions{}, false
	}
	dateTo, err := queryTimePtr(c, "date_to")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidDateTo)
		return repository.OrderListOptions{}, false
	}

	return repository.OrderListOptions{
		Page:     page,
		PageSize: pageSize,
		Statuses: statuses,
		UserID:   userID,
		PlayerID: playerID,
		GameID:   gameID,
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Keyword:  strings.TrimSpace(c.Query("keyword")),
	}, true
}

// buildPaymentListOptions parses query parameters into PaymentListOptions; on error responds and returns false.
func buildPaymentListOptions(c *gin.Context) (repository.PaymentListOptions, bool) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return repository.PaymentListOptions{}, false
	}

	statusTokens := parseCSVParams(c.QueryArray("status"))
	statuses := make([]model.PaymentStatus, 0, len(statusTokens))
	for _, token := range statusTokens {
		statuses = append(statuses, model.PaymentStatus(strings.ToLower(token)))
	}

	methodTokens := parseCSVParams(c.QueryArray("method"))
	methods := make([]model.PaymentMethod, 0, len(methodTokens))
	for _, token := range methodTokens {
		methods = append(methods, model.PaymentMethod(strings.ToLower(token)))
	}

	userID, err := queryUint64Ptr(c, "user_id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidUserID)
		return repository.PaymentListOptions{}, false
	}
	orderID, err := queryUint64Ptr(c, "order_id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidOrderID)
		return repository.PaymentListOptions{}, false
	}
	dateFrom, err := queryTimePtr(c, "date_from")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidDateFrom)
		return repository.PaymentListOptions{}, false
	}
	dateTo, err := queryTimePtr(c, "date_to")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidDateTo)
		return repository.PaymentListOptions{}, false
	}

	return repository.PaymentListOptions{
		Page:     page,
		PageSize: pageSize,
		Statuses: statuses,
		Methods:  methods,
		UserID:   userID,
		OrderID:  orderID,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}, true
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

// buildUserListOptions parses query parameters for user listing.
func buildUserListOptions(c *gin.Context) (repository.UserListOptions, bool) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return repository.UserListOptions{}, false
	}

	roleTokens := parseCSVParams(c.QueryArray("role"))
	roles := make([]model.Role, 0, len(roleTokens))
	for _, t := range roleTokens {
		roles = append(roles, model.Role(strings.ToLower(t)))
	}

	statusTokens := parseCSVParams(c.QueryArray("status"))
	statuses := make([]model.UserStatus, 0, len(statusTokens))
	for _, t := range statusTokens {
		statuses = append(statuses, model.UserStatus(strings.ToLower(t)))
	}

	dateFrom, err := queryTimePtr(c, "date_from")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidDateFrom)
		return repository.UserListOptions{}, false
	}
	dateTo, err := queryTimePtr(c, "date_to")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidDateTo)
		return repository.UserListOptions{}, false
	}

	return repository.UserListOptions{
		Page:     page,
		PageSize: pageSize,
		Roles:    roles,
		Statuses: statuses,
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Keyword:  strings.TrimSpace(c.Query("keyword")),
	}, true
}
