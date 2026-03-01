package dto

import "gamelink/internal/model"

// BatchRoleUpdateRequest 批量修改用户角色请求
type BatchRoleUpdateRequest struct {
	UserIDs []uint64   `json:"userIds" binding:"required,min=1,max=100"`
	Role    model.Role `json:"role" binding:"required,oneof=user player admin"`
}

// BatchDeleteUsersRequest 批量删除用户请求
type BatchDeleteUsersRequest struct {
	IDs []uint64 `json:"ids" binding:"required,min=1,max=100"`
}

// BatchOperationResponse 批量操作响应
type BatchOperationResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Data    BatchOperationResult `json:"data"`
}

// BatchOperationResult 批量操作结果
type BatchOperationResult struct {
	Updated int          `json:"updated"`
	Failed  int          `json:"failed"`
	Errors  []BatchError `json:"errors"`
}

// BatchError 批量操作错误
type BatchError struct {
	UserID uint64 `json:"userId"`
	Reason string `json:"reason"`
}
