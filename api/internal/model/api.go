package model

// APIResponse is the standard envelope for API responses.
// Use APIResponse[YourData] to wrap data payloads consistently.
type APIResponse[T any] struct {
	Success    bool        `json:"success"`
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Data       T           `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Meta       any         `json:"meta,omitempty"`
	TraceID    string      `json:"traceId,omitempty"`
}

// Pagination describes list pagination metadata.
type Pagination struct {
	Page       int  `json:"page"`
	PageSize   int  `json:"page_size"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// ErrorResponse 错误响应（用于Swagger文档，避免泛型语法）
type ErrorResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	TraceID string `json:"traceId,omitempty"`
}

// SuccessResponse 通用成功响应（用于Swagger文档，避免泛型语法）
// 注意：这是非泛型版本，仅供Swagger文档使用
// 实际代码中请使用 resp.OK() 等响应函数
type SuccessResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
	TraceID string `json:"traceId,omitempty"`
}

// FeedView 动态视图（用于Swagger文档）
type FeedView struct {
	ID               uint64          `json:"id"`
	AuthorID         uint64          `json:"authorId"`
	Content          string          `json:"content"`
	Visibility       string          `json:"visibility"`
	ModerationStatus string          `json:"moderationStatus"`
	ModerationNote   string          `json:"moderationNote,omitempty"`
	CreatedAt        string          `json:"createdAt"`
	Images           []FeedImageView `json:"images"`
}

// FeedImageView 动态图片视图（用于Swagger文档）
type FeedImageView struct {
	URL       string `json:"url"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	SizeBytes int64  `json:"sizeBytes"`
	Order     int    `json:"order"`
}

// ListFeedsData 动态列表数据（用于Swagger文档）
type ListFeedsData struct {
	Items      []FeedView `json:"items"`
	NextCursor string     `json:"nextCursor,omitempty"`
}

// PlatformStatsResponse 平台统计响应（用于Swagger文档，避免泛型语法）
type PlatformStatsResponse struct {
	TotalCommission float64 `json:"total_commission"`
	TotalSettlement float64 `json:"total_settlement"`
	TotalOrders     int64   `json:"total_orders"`
	Month           string  `json:"month"`
}

// DashboardOverviewStats 仪表板总览统计（用于Swagger文档）
type DashboardOverviewStats struct {
	TotalUsers       int64 `json:"totalUsers"`
	TotalPlayers     int64 `json:"totalPlayers"`
	TotalOrders      int64 `json:"totalOrders"`
	TodayOrders      int64 `json:"todayOrders"`
	TodayRevenue     int64 `json:"todayRevenue"`
	PendingOrders    int64 `json:"pendingOrders"`
	ActiveServices   int64 `json:"activeServices"`
	PendingWithdraws int64 `json:"pendingWithdraws"`
}
