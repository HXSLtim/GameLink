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
	db               *gorm.DB
	userRepo         repository.UserRepository
	tagRepo          repository.UserTagRepository
	notificationRepo repository.NotificationRepository
}

// NewBatchOperationService 创建批量操作服务
func NewBatchOperationService(
	db *gorm.DB,
	userRepo repository.UserRepository,
	tagRepo repository.UserTagRepository,
	notificationRepo repository.NotificationRepository,
) *BatchOperationService {
	return &BatchOperationService{
		db:               db,
		userRepo:         userRepo,
		tagRepo:          tagRepo,
		notificationRepo: notificationRepo,
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
			EntityType:   string(model.OpEntityUser),
			Action:       string(model.OpActionUpdateRole),
			Reason:       fmt.Sprintf("批量更新角色为:%s", req.Role),
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
			EntityType:   string(model.OpEntityUser),
			Action:       string(model.OpActionUpdateStatus),
			ActorUserID:  &operatorID,
			Reason:       fmt.Sprintf("批量更新状态为:%s, 原因:%s", req.Status, req.Reason),
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
			EntityType:   string(model.OpEntityUser),
			Action:       string(model.OpActionDelete),
			ActorUserID:  &operatorID,
			Reason:       fmt.Sprintf("批量删除用户,原因:%s", req.Reason),
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
	// Target指定目标类型：users（指定用户列表）、role（按角色）、all（全体用户）
	Target string `json:"target" binding:"required,oneof=users role all"`

	// 当Target=users时使用，最多1000个用户
	UserIDs []uint64 `json:"userIds,omitempty"`

	// 当Target=role时使用，可指定多个角色
	Roles []string `json:"roles,omitempty"`

	Cents  int64  `json:"cents" binding:"required,min=1,max=1000000"` // 积分金额（分），最多10000元=1000000分
	Reason string `json:"reason" binding:"required,max=200"`
	Type   string `json:"type" binding:"required,oneof=admin activity compensation"`
}

// BatchAddPoints 批量增加用户积分
func (s *BatchOperationService) BatchAddPoints(ctx context.Context, req *BatchAddPointsRequest, operatorID uint64) (successCount, failedCount int, err error) {
	// 根据Target类型获取用户ID列表
	var userIDs []uint64
	var targetDesc string

	switch req.Target {
	case "users":
		// 指定用户列表模式
		if len(req.UserIDs) == 0 {
			return 0, 0, fmt.Errorf("target为users时，userIds不能为空")
		}
		if len(req.UserIDs) > 1000 {
			return 0, 0, fmt.Errorf("一次批量操作最多支持1000个用户")
		}
		userIDs = req.UserIDs
		targetDesc = fmt.Sprintf("指定用户（%d个）", len(userIDs))

	case "role":
		// 按角色筛选模式
		if len(req.Roles) == 0 {
			return 0, 0, fmt.Errorf("target为role时，roles不能为空")
		}
		// 转换角色字符串为model.Role类型
		roles := make([]model.Role, 0, len(req.Roles))
		for _, roleStr := range req.Roles {
			roles = append(roles, model.Role(roleStr))
		}
		// 查询符合角色条件的用户
		var users []model.User
		if err := s.db.WithContext(ctx).Where("role IN ?", roles).Find(&users).Error; err != nil {
			return 0, 0, fmt.Errorf("查询角色用户失败: %w", err)
		}
		if len(users) > 1000 {
			return 0, 0, fmt.Errorf("符合条件的用户超过1000个，请缩小范围或分批操作")
		}
		for _, user := range users {
			userIDs = append(userIDs, user.ID)
		}
		targetDesc = fmt.Sprintf("角色%v（%d个用户）", req.Roles, len(userIDs))

	case "all":
		// 全体用户模式
		var users []model.User
		if err := s.db.WithContext(ctx).Select("id").Find(&users).Error; err != nil {
			return 0, 0, fmt.Errorf("查询全体用户失败: %w", err)
		}
		if len(users) > 1000 {
			return 0, 0, fmt.Errorf("用户总数超过1000个，请使用分批操作")
		}
		for _, user := range users {
			userIDs = append(userIDs, user.ID)
		}
		targetDesc = fmt.Sprintf("全体用户（%d个）", len(userIDs))

	default:
		return 0, 0, fmt.Errorf("无效的target类型: %s", req.Target)
	}

	if len(userIDs) == 0 {
		return 0, 0, fmt.Errorf("没有找到符合条件的用户")
	}

	// 记录操作日志
	defer func() {
		operationLog := &model.OperationLog{
			EntityType:   string(model.OpEntityUser),
			Action:       "batch_add_points",
			ActorUserID:  &operatorID,
			Reason:       fmt.Sprintf("批量增加积分:%d分,目标:%s,原因:%s", req.Cents, targetDesc, req.Reason),
			MetadataJSON: []byte(fmt.Sprintf(`{"target": "%s", "userCount": %d, "cents": %d, "type": "%s"}`, req.Target, len(userIDs), req.Cents, req.Type)),
		}
		go s.recordOperation(ctx, operationLog)
	}()

	// 批量增加积分（积分就是余额，直接增加钱包余额）
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, userID := range userIDs {
			// 先确保用户有钱包记录（如果不存在则创建）
			var wallet model.Wallet
			result := tx.Where("user_id = ?", userID).FirstOrCreate(&wallet, model.Wallet{
				UserID:       userID,
				BalanceCents: 0,
				FrozenCents:  0,
			})
			if result.Error != nil {
				failedCount++
				continue
			}

			// 增加余额
			if err := tx.Model(&wallet).
				UpdateColumn("balance_cents", gorm.Expr("balance_cents + ?", req.Cents)).Error; err != nil {
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
	// Target指定目标类型：users（指定用户列表）、role（按角色）、all（全体用户）
	Target string `json:"target" binding:"required,oneof=users role all"`

	// 当Target=users时使用，最多1000个用户
	UserIDs []uint64 `json:"userIds,omitempty"`

	// 当Target=role时使用，可指定多个角色
	Roles []string `json:"roles,omitempty"`

	Title   string `json:"title" binding:"required,max=100"`
	Content string `json:"content" binding:"required,max=500"`
	Type    string `json:"type" binding:"required,oneof=system marketing personal activity"`
}

// BatchSendNotification 批量发送通知
func (s *BatchOperationService) BatchSendNotification(ctx context.Context, req *BatchSendNotificationRequest, operatorID uint64) error {
	// 根据Target类型获取用户ID列表
	var userIDs []uint64
	var targetDesc string

	switch req.Target {
	case "users":
		// 指定用户列表模式
		if len(req.UserIDs) == 0 {
			return fmt.Errorf("target为users时，userIds不能为空")
		}
		if len(req.UserIDs) > 1000 {
			return fmt.Errorf("一次批量操作最多支持1000个用户")
		}
		userIDs = req.UserIDs
		targetDesc = fmt.Sprintf("指定用户（%d个）", len(userIDs))

	case "role":
		// 按角色筛选模式
		if len(req.Roles) == 0 {
			return fmt.Errorf("target为role时，roles不能为空")
		}
		// 转换角色字符串为model.Role类型
		roles := make([]model.Role, 0, len(req.Roles))
		for _, roleStr := range req.Roles {
			roles = append(roles, model.Role(roleStr))
		}
		// 查询符合角色条件的用户
		var users []model.User
		if err := s.db.WithContext(ctx).Where("role IN ?", roles).Find(&users).Error; err != nil {
			return fmt.Errorf("查询角色用户失败: %w", err)
		}
		if len(users) > 1000 {
			return fmt.Errorf("符合条件的用户超过1000个，请缩小范围或分批操作")
		}
		for _, user := range users {
			userIDs = append(userIDs, user.ID)
		}
		targetDesc = fmt.Sprintf("角色%v（%d个用户）", req.Roles, len(userIDs))

	case "all":
		// 全体用户模式
		var users []model.User
		if err := s.db.WithContext(ctx).Select("id").Find(&users).Error; err != nil {
			return fmt.Errorf("查询全体用户失败: %w", err)
		}
		if len(users) > 1000 {
			return fmt.Errorf("用户总数超过1000个，请使用分批操作")
		}
		for _, user := range users {
			userIDs = append(userIDs, user.ID)
		}
		targetDesc = fmt.Sprintf("全体用户（%d个）", len(userIDs))

	default:
		return fmt.Errorf("无效的target类型: %s", req.Target)
	}

	if len(userIDs) == 0 {
		return fmt.Errorf("没有找到符合条件的用户")
	}

	// 记录操作日志
	operationLog := &model.OperationLog{
		EntityType:   string(model.OpEntityUser),
		Action:       "batch_send_notification",
		ActorUserID:  &operatorID,
		Reason:       fmt.Sprintf("批量发送通知,目标:%s,标题:%s", targetDesc, req.Title),
		MetadataJSON: []byte(fmt.Sprintf(`{"target": "%s", "userCount": %d, "title": "%s", "type": "%s"}`, req.Target, len(userIDs), req.Title, req.Type)),
	}
	go s.recordOperation(ctx, operationLog)

	// 根据通知类型确定优先级
	var priority model.NotificationPriority
	switch req.Type {
	case "system":
		priority = model.NotificationPriorityHigh
	case "marketing":
		priority = model.NotificationPriorityNormal
	case "personal":
		priority = model.NotificationPriorityLow
	default:
		priority = model.NotificationPriorityNormal
	}

	// 批量创建通知
	for _, userID := range userIDs {
		note := &model.NotificationEvent{
			UserID:   userID,
			Title:    req.Title,
			Message:  req.Content,
			Channel:  "web",
			Priority: priority,
		}
		if err := s.notificationRepo.Create(ctx, note); err != nil {
			fmt.Printf("[Notification] Failed to create for user %d: %v\n", userID, err)
			// 单个通知失败不中断整个批处理
			continue
		}
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
