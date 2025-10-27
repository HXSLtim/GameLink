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
