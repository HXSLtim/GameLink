/**
 * @file user batch handler
 * @description 用户批量操作API接口
 */

package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"gamelink/internal/service/user"
	"github.com/gin-gonic/gin"
)

// RegisterBatchRoutes 注册批量操作路由
func RegisterBatchRoutes(rg *gin.RouterGroup, batchService *user.BatchOperationService) {
	batch := rg.Group("/users/batch")
	{
		batch.POST("/role", batchUpdateUserRoleHandler(batchService))
		batch.POST("/status", batchUpdateUserStatusHandler(batchService))
		batch.POST("/delete", batchDeleteUsersHandler(batchService))
		batch.POST("/points", batchAddPointsHandler(batchService))
		batch.POST("/notification", batchSendNotificationHandler(batchService))
	}
}

// BatchResponse 批量操作响应
type BatchResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	SuccessCount int    `json:"successCount"`
	FailedCount  int    `json:"failedCount"`
}

// BatchUpdateUserRoleRequest 批量更新用户角色请求
type BatchUpdateUserRoleRequest struct {
	UserIDs []uint64 `json:"userIds" binding:"required,min=1,max=1000"`
	Role    string   `json:"role" binding:"required,oneof=user player admin"`
}

// batchUpdateUserRoleHandler 批量更新用户角色
// @Summary 批量更新用户角色
// @Description 批量更新多个用户的角色（支持user/player/admin）
// @Tags 用户批量操作
// @Accept json
// @Produce json
// @Param body body BatchUpdateUserRoleRequest true "批量更新角色请求"
// @Success 200 {object} BatchResponse
// @Router /admin/users/batch/role [post]
func batchUpdateUserRoleHandler(s *user.BatchOperationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req BatchUpdateUserRoleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "参数验证失败",
				"error":   err.Error(),
			})
			return
		}


		successCount, failedCount, err := s.BatchUpdateUserRole(c.Request.Context(), &user.BatchUpdateUserRoleRequest{
			UserIDs: req.UserIDs,
			Role:    req.Role,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "批量更新角色失败",
				"error":   err.Error(),
			})
			return
		}

		message := "批量更新角色成功"
		if failedCount > 0 {
			message = fmt.Sprintf("批量更新角色完成，成功%d个，失败%d个", successCount, failedCount)
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": message,
			"successCount": successCount,
			"failedCount":  failedCount,
		})
	}
}

// BatchUpdateUserStatusRequest 批量更新用户状态请求
type BatchUpdateUserStatusRequest struct {
	UserIDs []uint64 `json:"userIds" binding:"required,min=1,max=1000"`
	Status  string   `json:"status" binding:"required,oneof=active banned suspended"`
	Reason  string   `json:"reason,omitempty" binding:"max=500"`
}

// batchUpdateUserStatusHandler 批量更新用户状态
// @Summary 批量更新用户状态
// @Description 批量更新多个用户的状态（active正常/banned封禁/suspended停用），支持填写原因
// @Tags 用户批量操作
// @Accept json
// @Produce json
// @Param body body BatchUpdateUserStatusRequest true "批量更新状态请求"
// @Success 200 {object} BatchResponse
// @Router /admin/users/batch/status [post]
func batchUpdateUserStatusHandler(s *user.BatchOperationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req BatchUpdateUserStatusRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "参数验证失败",
				"error":   err.Error(),
			})
			return
		}

		operatorID := getUserIDFromContext(c)

		successCount, failedCount, err := s.BatchUpdateUserStatus(
			c.Request.Context(),
			&user.BatchUpdateUserStatusRequest{
				UserIDs: req.UserIDs,
				Status:  req.Status,
				Reason:  req.Reason,
			},
			operatorID,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "批量更新状态失败",
				"error":   err.Error(),
			})
			return
		}

		message := "批量更新状态成功"
		if failedCount > 0 {
			message = fmt.Sprintf("批量更新状态完成，成功%d个，失败%d个", successCount, failedCount)
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": message,
			"successCount": successCount,
			"failedCount":  failedCount,
		})
	}
}

// BatchDeleteUsersRequest 批量删除用户请求
type BatchDeleteUsersRequest struct {
	UserIDs []uint64 `json:"userIds" binding:"required,min=1,max=1000"`
	Reason  string   `json:"reason,omitempty" binding:"max=500"`
}

// batchDeleteUsersHandler 批量删除用户
// @Summary 批量删除用户
// @Description 批量软删除多个用户（慎用），支持填写删除原因
// @Tags 用户批量操作
// @Accept json
// @Produce json
// @Param body body BatchDeleteUsersRequest true "批量删除请求"
// @Success 200 {object} BatchResponse
// @Router /admin/users/batch/delete [post]
func batchDeleteUsersHandler(s *user.BatchOperationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req BatchDeleteUsersRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "参数验证失败",
				"error":   err.Error(),
			})
			return
		}

		operatorID := getUserIDFromContext(c)

		successCount, failedCount, err := s.BatchDeleteUsers(
			c.Request.Context(),
			&user.BatchDeleteUsersRequest{
				UserIDs: req.UserIDs,
				Reason:  req.Reason,
			},
			operatorID,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "批量删除失败",
				"error":   err.Error(),
			})
			return
		}

		message := "批量删除成功"
		if failedCount > 0 {
			message = fmt.Sprintf("批量删除完成，成功%d个，失败%d个", successCount, failedCount)
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": message,
			"successCount": successCount,
			"failedCount":  failedCount,
		})
	}
}

// BatchAddPointsRequest 批量增加积分请求
type BatchAddPointsRequest struct {
	UserIDs []uint64 `json:"userIds" binding:"required,min=1,max=1000"`
	Points  int64    `json:"points" binding:"required,min=1,max=10000"`
	Reason  string   `json:"reason" binding:"required,max=200"`
	Type    string   `json:"type" binding:"required,oneof=admin activity compensation"`
}

// batchAddPointsHandler 批量增加积分
// @Summary 批量增加用户积分
// @Description 为多个用户批量增加积分（用于奖励、补偿等场景）
// @Tags 用户批量操作
// @Accept json
// @Produce json
// @Param body body BatchAddPointsRequest true "批量增加积分请求"
// @Success 200 {object} BatchResponse
// @Router /admin/users/batch/points [post]
func batchAddPointsHandler(s *user.BatchOperationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req BatchAddPointsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "参数验证失败",
				"error":   err.Error(),
			})
			return
		}

		operatorID := getUserIDFromContext(c)

		successCount, failedCount, err := s.BatchAddPoints(
			c.Request.Context(),
			&user.BatchAddPointsRequest{
				UserIDs: req.UserIDs,
				Points:  req.Points,
				Reason:  req.Reason,
				Type:    req.Type,
			},
			operatorID,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "批量增加积分失败",
				"error":   err.Error(),
			})
			return
		}

		message := fmt.Sprintf("批量增加积分成功，共增加%d积分", req.Points*int64(successCount))
		if failedCount > 0 {
			message = fmt.Sprintf("批量增加积分完成，成功%d个，失败%d个", successCount, failedCount)
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": message,
			"successCount": successCount,
			"failedCount":  failedCount,
		})
	}
}

// BatchSendNotificationRequest 批量发送通知请求
type BatchSendNotificationRequest struct {
	UserIDs []uint64 `json:"userIds" binding:"required,min=1,max=1000"`
	Title   string   `json:"title" binding:"required,max=100"`
	Content string   `json:"content" binding:"required,max=500"`
	Type    string   `json:"type" binding:"required,oneof=system marketing personal"`
}

// batchSendNotificationHandler 批量发送通知
// @Summary 批量发送站内通知
// @Description 向多个用户批量发送通知消息
// @Tags 用户批量操作
// @Accept json
// @Produce json
// @Param body body BatchSendNotificationRequest true "批量发送通知请求"
// @Success 200 {object} ApiResponse
// @Router /admin/users/batch/notification [post]
func batchSendNotificationHandler(s *user.BatchOperationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req BatchSendNotificationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "参数验证失败",
				"error":   err.Error(),
			})
			return
		}

		operatorID := getUserIDFromContext(c)

		if err := s.BatchSendNotification(
			c.Request.Context(),
			&user.BatchSendNotificationRequest{
				UserIDs: req.UserIDs,
				Title:   req.Title,
				Content: req.Content,
				Type:    req.Type,
			},
			operatorID,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "批量发送通知失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "批量发送通知成功",
		})
	}
}

// getUserIDFromContext 从上下文获取用户ID
func getUserIDFromContext(c *gin.Context) uint64 {
	// 从JWT token或session获取用户ID
	// 这里简化处理，实际应该从中间件获取
	userID, _ := strconv.ParseUint(c.GetString("userId"), 10, 64)
	return userID
}
