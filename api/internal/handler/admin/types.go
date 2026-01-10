package admin

// BatchOperationResponse 批量操作响应
// 统一的批量操作响应格式，用于所有批量操作接口
type BatchOperationResponse struct {
	SuccessCount int                   `json:"success_count"`
	FailedCount  int                   `json:"failed_count"`
	TotalCount   int                   `json:"total_count"`
	FailedItems  []BatchOperationError `json:"failed_items,omitempty"`
	SuccessItems []uint64              `json:"success_items,omitempty"`
}

// BatchOperationError 单个操作错误详情
type BatchOperationError struct {
	ID      uint64 `json:"id"`
	Message string `json:"message"`
}
