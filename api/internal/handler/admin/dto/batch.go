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

// BatchUpdateUserRoleRequest 批量更新用户角色请求
type BatchUpdateUserRoleRequest struct {
	UserIDs []uint64 `json:"userIds" binding:"required,min=1,max=1000,dive,gt=0"`
	Role    string   `json:"role" binding:"required,oneof=user player admin"`
}

// BatchUpdateUserStatusRequest 批量更新用户状态请求
type BatchUpdateUserStatusRequest struct {
	UserIDs []uint64 `json:"userIds" binding:"required,min=1,max=1000,dive,gt=0"`
	Status  string   `json:"status" binding:"required,oneof=active suspended banned"`
	Reason  string   `json:"reason" binding:"omitempty,max=500"`
}

// BatchAddUsersPointsRequest 批量增加用户积分请求
type BatchAddUsersPointsRequest struct {
	UserIDs []uint64 `json:"userIds" binding:"required,min=1,max=1000,dive,gt=0"`
	Cents   int64    `json:"cents" binding:"required,min=1,max=1000000"`
	Reason  string   `json:"reason" binding:"required,max=200"`
}

// BatchNotifyUsersRequest 批量发送通知请求
type BatchNotifyUsersRequest struct {
	UserIDs  []uint64 `json:"userIds" binding:"required,min=1,max=1000,dive,gt=0"`
	Title    string   `json:"title" binding:"required,max=100"`
	Content  string   `json:"content" binding:"required,max=500"`
	Priority string   `json:"priority" binding:"omitempty,oneof=low normal high"`
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
