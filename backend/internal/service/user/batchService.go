/**
 * @file batch operation service
 * @description 用户批量操作业务逻辑层
 */

package user

import (
	"context"
	"fmt"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gorm.io/gorm"
)

// BatchOperationService 批量操作服务
type BatchOperationService struct {
	db       *gorm.DB
	userRepo repository.UserRepository
	tagRepo  repository.UserTagRepository
}

// NewBatchOperationService 创建批量操作服务
func NewBatchOperationService(
	db *gorm.DB,
	userRepo repository.UserRepository,
	tagRepo repository.UserTagRepository,
) *BatchOperationService {
	return &BatchOperationService{
		db:       db,
		userRepo: userRepo,
		tagRepo:  tagRepo,
	}
}

// BatchUpdateUserRoleRequest 批量更新用户角色请求
type BatchUpdateUserRoleRequest struct {
	UserIDs []uint64 `json:"userIds" binding:"required,min=1,max=1000"` // 用户ID列表（最多1000个）
	Role    string   `json:"role" binding:"required,oneof=user player admin"`
}

// BatchUpdateUserRole 批量更新用户角色
// @return successCount 成功数量
// @return failedCount 失败数量
// @return err 错误信息
func (s *BatchOperationService) BatchUpdateUserRole(ctx context.Context, req *BatchUpdateUserRoleRequest) (successCount, failedCount int, err error) {
	if len(req.UserIDs) == 0 {
		return 0, 0, fmt.Errorf("用户ID列表不能为空")
	}

	if len(req.UserIDs) > 1000 {
		return 0, 0, fmt.Errorf("一次批量操作最多支持1000个用户")
	}

	// 记录操作日志
	defer func() {
		operationLog := &model.OperationLog{
			EntityType:  string(model.OpEntityUser),
			Action:      string(model.OpActionUpdateRole),
			Reason:      fmt.Sprintf("批量更新角色为:%s", req.Role),
			MetadataJSON: []byte(fmt.Sprintf(`{"userIds": %v, "role": "%s"}`, req.UserIDs, req.Role)),
		}
		// 异步记录日志
		go s.recordOperation(ctx, operationLog)
	}()

	// 使用事务批量更新
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, userID := range req.UserIDs {
			if err := tx.Model(&model.User{}).Where("id = ?", userID).Update("role", req.Role).Error; err != nil {
				// 记录失败，但不回滚（继续处理其他用户）
				failedCount++
				continue
			}
			successCount++
		}
		return nil
	})

	if err != nil {
		return successCount, failedCount, fmt.Errorf("批量更新角色失败: %w", err)
	}

	return successCount, failedCount, nil
}

// BatchUpdateUserStatusRequest 批量更新用户状态请求
type BatchUpdateUserStatusRequest struct {
	UserIDs []uint64 `json:"userIds" binding:"required,min=1,max=1000"`
	Status  string   `json:"status" binding:"required,oneof=active banned suspended"`
	Reason  string   `json:"reason,omitempty" binding:"max=500"`
}

// BatchUpdateUserStatus 批量更新用户状态
func (s *BatchOperationService) BatchUpdateUserStatus(ctx context.Context, req *BatchUpdateUserStatusRequest, operatorID uint64) (successCount, failedCount int, err error) {
	if len(req.UserIDs) > 1000 {
		return 0, 0, fmt.Errorf("一次批量操作最多支持1000个用户")
	}

	// 记录操作日志
	defer func() {
		operationLog := &model.OperationLog{
			EntityType:  string(model.OpEntityUser),
			Action:      string(model.OpActionUpdateStatus),
			ActorUserID: &operatorID,
			Reason:      fmt.Sprintf("批量更新状态为:%s, 原因:%s", req.Status, req.Reason),
			MetadataJSON: []byte(fmt.Sprintf(`{"userIds": %v, "status": "%s"}`, req.UserIDs, req.Status)),
		}
		go s.recordOperation(ctx, operationLog)
	}()

	// 批量更新
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, userID := range req.UserIDs {
			updates := map[string]interface{}{
				"status": req.Status,
			}

			if err := tx.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
				failedCount++
				continue
			}
			successCount++
		}
		return nil
	})

	return successCount, failedCount, err
}

// BatchDeleteUsersRequest 批量删除用户请求
type BatchDeleteUsersRequest struct {
	UserIDs []uint64 `json:"userIds" binding:"required,min=1,max=1000"`
	Reason  string   `json:"reason,omitempty" binding:"max=500"`
}

// BatchDeleteUsers 批量删除用户（软删除）
func (s *BatchOperationService) BatchDeleteUsers(ctx context.Context, req *BatchDeleteUsersRequest, operatorID uint64) (successCount, failedCount int, err error) {
	if len(req.UserIDs) > 1000 {
		return 0, 0, fmt.Errorf("一次批量操作最多支持1000个用户")
	}

	// 记录操作日志
	defer func() {
		operationLog := &model.OperationLog{
			EntityType:  string(model.OpEntityUser),
			Action:      string(model.OpActionDelete),
			ActorUserID: &operatorID,
			Reason:      fmt.Sprintf("批量删除用户,原因:%s", req.Reason),
			MetadataJSON: []byte(fmt.Sprintf(`{"userIds": %v}`, req.UserIDs)),
		}
		go s.recordOperation(ctx, operationLog)
	}()

	// 批量软删除
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, userID := range req.UserIDs {
			if err := tx.Delete(&model.User{}, userID).Error; err != nil {
				failedCount++
				continue
			}
			successCount++
		}
		return nil
	})

	return successCount, failedCount, err
}

// BatchAddPointsRequest 批量增加积分请求
type BatchAddPointsRequest struct {
	UserIDs []uint64 `json:"userIds" binding:"required,min=1,max=1000"`
	Points  int64    `json:"points" binding:"required,min=1,max=10000"` // 最多10000积分
	Reason  string   `json:"reason" binding:"required,max=200"`
	Type    string   `json:"type" binding:"required,oneof=admin activity compensation"`
}

// BatchAddPoints 批量增加用户积分
func (s *BatchOperationService) BatchAddPoints(ctx context.Context, req *BatchAddPointsRequest, operatorID uint64) (successCount, failedCount int, err error) {
	if len(req.UserIDs) > 1000 {
		return 0, 0, fmt.Errorf("一次批量操作最多支持1000个用户")
	}

	// 记录操作日志
	defer func() {
		operationLog := &model.OperationLog{
			EntityType:  string(model.OpEntityUser),
			Action:      "batch_add_points",
			ActorUserID: &operatorID,
			Reason:      fmt.Sprintf("批量增加积分:%d,原因:%s", req.Points, req.Reason),
			MetadataJSON: []byte(fmt.Sprintf(`{"userIds": %v, "points": %d, "type": "%s"}`, req.UserIDs, req.Points, req.Type)),
		}
		go s.recordOperation(ctx, operationLog)
	}()

	// 批量增加积分
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, userID := range req.UserIDs {
			// 这里应该调用积分服务，简化处理直接在用户表更新
			if err := tx.Model(&model.User{}).
				Where("id = ?", userID).
				UpdateColumn("points", gorm.Expr("points + ?", req.Points)).Error; err != nil {
				failedCount++
				continue
			}
			successCount++
		}
		return nil
	})

	return successCount, failedCount, err
}

// BatchSendNotificationRequest 批量发送通知请求
type BatchSendNotificationRequest struct {
	UserIDs []uint64 `json:"userIds" binding:"required,min=1,max=1000"`
	Title   string   `json:"title" binding:"required,max=100"`
	Content string   `json:"content" binding:"required,max=500"`
	Type    string   `json:"type" binding:"required,oneof=system marketing personal"`
}

// BatchSendNotification 批量发送通知
func (s *BatchOperationService) BatchSendNotification(ctx context.Context, req *BatchSendNotificationRequest, operatorID uint64) error {
	if len(req.UserIDs) > 1000 {
		return fmt.Errorf("一次批量操作最多支持1000个用户")
	}

	// 记录操作日志
	operationLog := &model.OperationLog{
		EntityType:  string(model.OpEntityUser),
		Action:      "batch_send_notification",
		ActorUserID: &operatorID,
		Reason:      fmt.Sprintf("批量发送通知,标题:%s", req.Title),
		MetadataJSON: []byte(fmt.Sprintf(`{"userIds": %v, "title": "%s", "type": "%s"}`, req.UserIDs, req.Title, req.Type)),
	}
	go s.recordOperation(ctx, operationLog)

	// 这里实际应该调用通知服务，简化处理：直接打印日志
	// TODO: 集成实际的通知服务（消息队列或通知模块）
	for _, userID := range req.UserIDs {
		fmt.Printf("[Notification] Sending to user %d: %s - %s\n", userID, req.Title, req.Type)
	}

	return nil
}

// recordOperation 异步记录操作日志
func (s *BatchOperationService) recordOperation(ctx context.Context, log *model.OperationLog) {
	// 这是一个占位方法，实际需要调用操作日志服务
	// 简化实现：直接打印日志
	fmt.Printf("[OperationLog] Entity:%s, Action:%s, Reason:%s\n",
		log.EntityType, log.Action, log.Reason)
}
