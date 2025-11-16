package model

// APIResponse is the standard envelope for API responses.
// Use APIResponse[YourData] to wrap data payloads consistently.
type APIResponse[T any] struct {
	Success    bool        `json:"success"`
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Data       T           `json:"data,omitempty"`
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

// SuccessResponse 成功响应（用于Swagger文档，避免泛型语法）
type SuccessResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
