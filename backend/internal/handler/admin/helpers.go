package admin

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	repoiface "gamelink/internal/repository/interfaces"
	adminservice "gamelink/internal/service/admin"
	"gamelink/pkg/apierr"
)

// ============================================================================
// 统一响应函数
// ============================================================================

// respondSuccess 统一成功响应
func respondSuccess[T any](c *gin.Context, data T) {
	respondSuccessWithMsg(c, "OK", data)
}

// respondSuccessWithMsg 统一成功响应（带自定义消息）
func respondSuccessWithMsg[T any](c *gin.Context, message string, data T) {
	writeJSON(c, http.StatusOK, model.APIResponse[T]{
		Success: true,
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

// respondCreated 统一创建成功响应
func respondCreated[T any](c *gin.Context, data T) {
	writeJSON(c, http.StatusCreated, model.APIResponse[T]{
		Success: true,
		Code:    http.StatusCreated,
		Message: "created",
		Data:    data,
	})
}

// respondList 统一列表响应（带分页）
func respondList[T any](c *gin.Context, data []T, pagination *model.Pagination) {
	data = ensureSlice(data)
	writeJSON(c, http.StatusOK, model.APIResponse[[]T]{
		Success:    true,
		Code:       http.StatusOK,
		Message:    "OK",
		Data:       data,
		Pagination: pagination,
	})
}

// respondDeleted 统一删除成功响应
func respondDeleted(c *gin.Context) {
	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "deleted",
	})
}

// respondUpdated 统一更新成功响应
func respondUpdated[T any](c *gin.Context, data T) {
	writeJSON(c, http.StatusOK, model.APIResponse[T]{
		Success: true,
		Code:    http.StatusOK,
		Message: "updated",
		Data:    data,
	})
}

// ============================================================================
// 统一错误处理函数
// ============================================================================

// respondError 统一错误响应（推荐使用）
// 自动处理 apierr.APIError、adminservice.ErrNotFound 等常见错误类型
func respondError(c *gin.Context, err error) {
	// 处理 apierr.APIError
	if apiErr, ok := err.(*apierr.APIError); ok {
		writeJSON(c, apiErr.Code, model.APIResponse[any]{
			Success: false,
			Code:    apiErr.Code,
			Message: apiErr.Message,
			TraceID: apiErr.RequestID,
		})
		return
	}

	// 处理 adminservice.ErrNotFound
	if errors.Is(err, adminservice.ErrNotFound) || errors.Is(err, repository.ErrNotFound) {
		writeJSON(c, http.StatusNotFound, model.APIResponse[any]{
			Success: false,
			Code:    http.StatusNotFound,
			Message: "resource not found",
		})
		return
	}

	// 处理验证错误
	if apierr.IsValidationError(err) {
		writeJSON(c, http.StatusBadRequest, model.APIResponse[any]{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	// 默认内部错误
	writeJSON(c, http.StatusInternalServerError, model.APIResponse[any]{
		Success: false,
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	})
}

// respondBadRequest 统一400错误响应
func respondBadRequest(c *gin.Context, message string) {
	writeJSON(c, http.StatusBadRequest, model.APIResponse[any]{
		Success: false,
		Code:    http.StatusBadRequest,
		Message: message,
	})
}

// respondNotFound 统一404错误响应
func respondNotFound(c *gin.Context, message string) {
	if message == "" {
		message = "resource not found"
	}
	writeJSON(c, http.StatusNotFound, model.APIResponse[any]{
		Success: false,
		Code:    http.StatusNotFound,
		Message: message,
	})
}

// respondUnauthorized 统一401错误响应
func respondUnauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "unauthorized"
	}
	writeJSON(c, http.StatusUnauthorized, model.APIResponse[any]{
		Success: false,
		Code:    http.StatusUnauthorized,
		Message: message,
	})
}

// respondForbidden 统一403错误响应
func respondForbidden(c *gin.Context, message string) {
	if message == "" {
		message = "forbidden"
	}
	writeJSON(c, http.StatusForbidden, model.APIResponse[any]{
		Success: false,
		Code:    http.StatusForbidden,
		Message: message,
	})
}

// respondInternalError 统一500错误响应
func respondInternalError(c *gin.Context, message string) {
	writeJSON(c, http.StatusInternalServerError, model.APIResponse[any]{
		Success: false,
		Code:    http.StatusInternalServerError,
		Message: message,
	})
}

// ============================================================================
// 统一参数解析函数
// ============================================================================

// ParseIDAndRespond 解析ID参数并响应错误（如果需要）
// 返回 (id, ok)，ok=false 时已写入错误响应，调用方应直接 return
func ParseIDAndRespond(c *gin.Context, paramName string) (uint64, bool) {
	id, err := handler.ParseIDParam(c, paramName)
	if err != nil {
		respondError(c, err)
		return 0, false
	}
	return id, true
}

// ValidateAndRespond 验证请求体并响应错误（如果需要）
// 返回 ok，ok=false 时已写入错误响应，调用方应直接 return
func ValidateAndRespond(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		respondError(c, apierr.BadRequest(apierr.ErrInvalidJSONPayload).WithDetails(err.Error()))
		return false
	}
	return true
}

// QueryUint64PtrAndRespond 解析可选的 uint64 查询参数
// 返回 (value, ok)，ok=false 时已写入错误响应
func QueryUint64PtrAndRespond(c *gin.Context, key string, errMsg string) (*uint64, bool) {
	v, err := queryUint64Ptr(c, key)
	if err != nil {
		respondBadRequest(c, errMsg)
		return nil, false
	}
	return v, true
}

// QueryTimePtrAndRespond 解析可选的时间查询参数
// 返回 (value, ok)，ok=false 时已写入错误响应
func QueryTimePtrAndRespond(c *gin.Context, key string, errMsg string) (*time.Time, bool) {
	v, err := queryTimePtr(c, key)
	if err != nil {
		respondBadRequest(c, errMsg)
		return nil, false
	}
	return v, true
}

// ============================================================================
// 兼容旧代码的函数（逐步废弃）
// ============================================================================

// parseUintParam 解析路径参数为 uint64（建议使用 ParseIDAndRespond）
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
	result := make([]string, 0)
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

// writeJSON 底层 JSON 响应函数
func writeJSON[T any](c *gin.Context, status int, payload model.APIResponse[T]) {
	// 从上下文中获取TraceID
	if payload.TraceID == "" {
		if rid, ok := c.Get("request_id"); ok {
			if ridStr, ok := rid.(string); ok {
				payload.TraceID = ridStr
			}
		}
	}
	c.JSON(status, payload)
}

// writeJSONError 底层错误响应函数（建议使用 respondError 或 respondBadRequest）
func writeJSONError(c *gin.Context, status int, message string) {
	writeJSON(c, status, model.APIResponse[any]{
		Success: false,
		Code:    status,
		Message: message,
	})
}

// respondAPIError 使用apierr包的错误响应（建议使用 respondError）
func respondAPIError(c *gin.Context, err error) {
	respondError(c, err)
}

func ensureSlice[T any](items []T) []T {
	if items == nil {
		return make([]T, 0)
	}
	return items
}

// getAdminUserID 从上下文中读取管理端用户 ID，并在缺失或类型不匹配时写入统一错误响应
// 返回值中的 bool 表示是否读取成功，调用方应在返回 false 时直接中断处理
func getAdminUserID(c *gin.Context) (uint64, bool) {
	val, ok := c.Get("user_id")
	if !ok {
		writeJSONError(c, 401, "missing admin user")
		return 0, false
	}

	id, ok := val.(uint64)
	if !ok {
		writeJSONError(c, 500, "invalid admin user id type")
		return 0, false
	}

	return id, true
}

// parsePagination 解析分页参数，支持 page/page_size 以及向后兼容的 pageSize
// 出错时会写入统一的 JSON 错误响应并返回 false
func parsePagination(c *gin.Context) (int, int, bool) {
	page, err := queryIntDefault(c, "page", 1)
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidPage)
		return 0, 0, false
	}

	// pageSize 同时兼容 page_size �?pageSize，优先读取蛇形命名以保持与大部分接口一�?
	pageSizeStr := strings.TrimSpace(c.Query("page_size"))
	if pageSizeStr == "" {
		pageSizeStr = strings.TrimSpace(c.Query("pageSize"))
	}
	pageSize := 20
	if pageSizeStr != "" {
		parsed, convErr := strconv.Atoi(pageSizeStr)
		if convErr != nil {
			writeJSONError(c, 400, apierr.ErrInvalidPageSize)
			return 0, 0, false
		}
		pageSize = parsed
	}

	return page, pageSize, true
}

// buildOrderListOptions parses query parameters into OrderListOptions; on error responds and returns false.
func buildOrderListOptions(c *gin.Context) (repoiface.OrderListOptions, bool) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return repoiface.OrderListOptions{}, false
	}

	statusTokens := parseCSVParams(c.QueryArray("status"))
	statuses := make([]model.OrderStatus, 0, len(statusTokens))
	for _, token := range statusTokens {
		statuses = append(statuses, normalizeOrderStatus(token))
	}

	userID, err := queryUint64Ptr(c, "user_id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidUserID)
		return repoiface.OrderListOptions{}, false
	}
	playerID, err := queryUint64Ptr(c, "player_id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidPlayerID)
		return repoiface.OrderListOptions{}, false
	}
	gameID, err := queryUint64Ptr(c, "game_id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidGameID)
		return repoiface.OrderListOptions{}, false
	}
	dateFrom, err := queryTimePtr(c, "date_from")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidDateFrom)
		return repoiface.OrderListOptions{}, false
	}
	dateTo, err := queryTimePtr(c, "date_to")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidDateTo)
		return repoiface.OrderListOptions{}, false
	}

	return repoiface.OrderListOptions{
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

// ============================================================================
// CSV 导出函数
// ============================================================================

// exportOperationLogsCSV writes operation logs as CSV attachment.
// 支持字段选择、多语言表头、时区转换等功能
func exportOperationLogsCSV(c *gin.Context, entity string, entityID uint64, items []model.OperationLog) {
	// default columns
	allowed := []string{"id", "entity_type", "entity_id", "actor_user_id", "action", "reason", "metadata", "created_at"}
	// parse fields
	rawFields := strings.TrimSpace(c.Query("fields"))
	fields := allowed
	if rawFields != "" {
		req := parseCSVParams([]string{rawFields})
		// validate and keep order
		pick := make([]string, 0, len(req))
		for _, f := range req {
			for _, a := range allowed {
				if f == a {
					pick = append(pick, f)
					break
				}
			}
		}
		if len(pick) > 0 {
			fields = pick
		}
	}

	// header i18n
	lang := strings.ToLower(strings.TrimSpace(c.Query("header_lang")))
	headerMapEn := map[string]string{
		"id": "id", "entity_type": "entity_type", "entity_id": "entity_id", "actor_user_id": "actor_user_id",
		"action": "action", "reason": "reason", "metadata": "metadata", "created_at": "created_at",
	}
	headerMapZh := map[string]string{
		"id": "编号", "entity_type": "实体", "entity_id": "实体ID", "actor_user_id": "操作人ID",
		"action": "动作", "reason": "原因", "metadata": "元数据", "created_at": "创建时间",
	}
	var header []string
	for _, f := range fields {
		if lang == "zh" {
			header = append(header, headerMapZh[f])
		} else {
			header = append(header, headerMapEn[f])
		}
	}

	filename := entity + "_" + strconv.FormatUint(entityID, 10) + "_logs.csv"
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	// excel-friendly BOM when requested or zh header
	bom := strings.EqualFold(strings.TrimSpace(c.Query("bom")), "true") || lang == "zh"
	if bom {
		_, _ = c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})
	}
	w := csv.NewWriter(c.Writer)
	_ = w.Write(header)
	// timezone
	tz := strings.TrimSpace(c.Query("tz"))
	var loc *time.Location
	if tz != "" {
		if l, err := time.LoadLocation(tz); err == nil {
			loc = l
		}
	}
	for _, it := range items {
		row := make([]string, 0, len(fields))
		for _, f := range fields {
			switch f {
			case "id":
				row = append(row, strconv.FormatUint(it.ID, 10))
			case "entity_type":
				row = append(row, it.EntityType)
			case "entity_id":
				row = append(row, strconv.FormatUint(it.EntityID, 10))
			case "actor_user_id":
				if it.ActorUserID != nil {
					row = append(row, strconv.FormatUint(*it.ActorUserID, 10))
				} else {
					row = append(row, "")
				}
			case "action":
				row = append(row, it.Action)
			case "reason":
				row = append(row, it.Reason)
			case "metadata":
				row = append(row, fmt.Sprintf("%q", string(it.MetadataJSON)))
			case "created_at":
				t := it.CreatedAt
				if loc != nil {
					t = t.In(loc)
				}
				row = append(row, t.Format(time.RFC3339))
			default:
				row = append(row, "")
			}
		}
		_ = w.Write(row)
	}
	w.Flush()
}
