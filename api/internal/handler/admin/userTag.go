/**
 * @file user tag handler
 * @description 用户标签管理API接口
 */

package admin

import (
	"net/http"
	"strconv"

	"gamelink/internal/model"
	"gamelink/internal/service/user"
	"gamelink/pkg/apierr"

	"github.com/gin-gonic/gin"
)

// RegisterTagRoutes 注册用户标签相关路由
func RegisterTagRoutes(rg *gin.RouterGroup, tagService *user.UserTagService) {
	// 标签管理
	tags := rg.Group("/user-tags")
	{
		tags.GET("", listTagsHandler(tagService))
		tags.POST("", createTagHandler(tagService))
		tags.GET("/:id", getTagHandler(tagService))
		tags.PUT("/:id", updateTagHandler(tagService))
		tags.DELETE("/:id", deleteTagHandler(tagService))
		tags.DELETE("/batch", batchDeleteTagsHandler(tagService))

		// 批量操作
		tags.POST("/batch/delete", BatchDeleteUserTags(tagService))
		tags.POST("/batch/assign", BatchAssignTagsToUsers(tagService))
		tags.POST("/batch/remove", BatchRemoveTagsFromUsers(tagService))

		// 批量查询
		tags.GET("/:id/users", getUsersByTagHandler(tagService))
	}

	// 用户标签操作
	users := rg.Group("/users")
	{
		users.GET("/:id/tags", getUserTagsHandler(tagService))
		users.POST("/:id/tags", addUserTagHandler(tagService))
		users.PUT("/:id/tags", batchSetUserTagsHandler(tagService))
		users.DELETE("/:id/tags/:tagId", removeUserTagHandler(tagService))
		users.POST("/user-tags/batch/assign", batchAssignTagsHandler(tagService))
	}
}

// TagResponse 标签响应结构
type TagResponse struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name        string `json:"name" binding:"required"`
	Color       string `json:"color" binding:"required"`
	Description string `json:"description"`
}

// UpdateTagRequest 更新标签请求
type UpdateTagRequest struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

// createTagHandler 创建标签
// @Summary 创建用户标签
// @Description 创建一个新的用户标签
// @Tags 用户标签管理
// @Accept json
// @Produce json
// @Param body body CreateTagRequest true "标签信息"
// @Success 200 {object} ApiResponse{data=TagResponse}
// @Router /admin/user-tags [post]
func createTagHandler(s *user.UserTagService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateTagRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "参数验证失败",
				"error":   err.Error(),
			})
			return
		}

		tag, err := s.CreateTag(c.Request.Context(), req.Name, req.Color, req.Description)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "创建标签失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "创建标签成功",
			"data":    toTagResponse(tag),
		})
	}
}

// listTagsHandler 获取标签列表
// @Summary 获取所有标签
// @Description 获取用户标签列表，按创建时间倒序
// @Tags 用户标签管理
// @Produce json
// @Success 200 {object} ApiResponse{data=[]TagResponse}
// @Router /admin/user-tags [get]
func listTagsHandler(s *user.UserTagService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tags, err := s.ListTags(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "获取标签列表失败",
				"error":   err.Error(),
			})
			return
		}

		responses := make([]*TagResponse, 0, len(tags))
		for _, tag := range tags {
			responses = append(responses, toTagResponse(&tag))
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "获取标签列表成功",
			"data":    responses,
		})
	}
}

// getTagHandler 获取标签详情
// @Summary 获取标签详情
// @Description 根据ID获取标签信息
// @Tags 用户标签管理
// @Produce json
// @Param id path int true "标签ID"
// @Success 200 {object} ApiResponse{data=TagResponse}
// @Router /admin/user-tags/{id} [get]
func getTagHandler(s *user.UserTagService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "无效的标签ID",
				"error":   err.Error(),
			})
			return
		}

		tag, err := s.GetTag(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "标签不存在",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "获取标签成功",
			"data":    toTagResponse(tag),
		})
	}
}

// updateTagHandler 更新标签
// @Summary 更新标签信息
// @Description 根据ID更新标签
// @Tags 用户标签管理
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Param body body UpdateTagRequest true "标签信息"
// @Success 200 {object} ApiResponse{data=TagResponse}
// @Router /admin/user-tags/{id} [put]
func updateTagHandler(s *user.UserTagService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "无效的标签ID",
				"error":   err.Error(),
			})
			return
		}

		var req UpdateTagRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "参数验证失败",
				"error":   err.Error(),
			})
			return
		}

		if err := s.UpdateTag(c.Request.Context(), id, req.Name, req.Color, req.Description); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "更新标签失败",
				"error":   err.Error(),
			})
			return
		}

		// 获取更新后的标签
		tag, err := s.GetTag(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "获取更新后的标签失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "更新标签成功",
			"data":    toTagResponse(tag),
		})
	}
}

// deleteTagHandler 删除标签
// @Summary 删除标签
// @Description 根据ID删除标签（同时会移除所有用户的该标签）
// @Tags 用户标签管理
// @Param id path int true "标签ID"
// @Success 200 {object} ApiResponse
// @Router /admin/user-tags/{id} [delete]
func deleteTagHandler(s *user.UserTagService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "无效的标签ID",
				"error":   err.Error(),
			})
			return
		}

		if err := s.DeleteTag(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "删除标签失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "删除标签成功",
		})
	}
}

// AddUserTagRequest 为用户添加标签请求
type AddUserTagRequest struct {
	TagID uint64 `json:"tagId" binding:"required"`
}

// addUserTagHandler 为用户添加标签
// @Summary 为用户添加标签
// @Description 为指定用户添加一个标签
// @Tags 用户标签管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param body body AddUserTagRequest true "标签ID"
// @Success 200 {object} ApiResponse
// @Router /admin/users/{id}/tags [post]
func addUserTagHandler(s *user.UserTagService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "无效的用户ID",
				"error":   err.Error(),
			})
			return
		}

		var req AddUserTagRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "参数验证失败",
				"error":   err.Error(),
			})
			return
		}

		if err := s.AddTagToUser(c.Request.Context(), userID, req.TagID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "添加标签失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "添加标签成功",
		})
	}
}

// BatchSetUserTagsRequest 批量设置用户标签请求
type BatchSetUserTagsRequest struct {
	TagIDs []uint64 `json:"tagIds" binding:"required"`
}

// batchSetUserTagsHandler 批量设置用户标签（覆盖）
// @Summary 批量设置用户标签
// @Description 批量设置用户的标签（会覆盖原有标签）
// @Tags 用户标签管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param body body BatchSetUserTagsRequest true "标签ID列表"
// @Success 200 {object} ApiResponse
// @Router /admin/users/{id}/tags [put]
func batchSetUserTagsHandler(s *user.UserTagService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "无效的用户ID",
				"error":   err.Error(),
			})
			return
		}

		var req BatchSetUserTagsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "参数验证失败",
				"error":   err.Error(),
			})
			return
		}

		if err := s.BatchSetUserTags(c.Request.Context(), userID, req.TagIDs); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "批量设置标签失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "批量设置标签成功",
		})
	}
}

// getUserTagsHandler 获取用户标签列表
// @Summary 获取用户标签
// @Description 获取指定用户的所有标签
// @Tags 用户标签管理
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} ApiResponse{data=[]TagResponse}
// @Router /admin/users/{id}/tags [get]
func getUserTagsHandler(s *user.UserTagService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "无效的用户ID",
				"error":   err.Error(),
			})
			return
		}

		tags, err := s.GetUserTags(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "获取用户标签失败",
				"error":   err.Error(),
			})
			return
		}

		responses := make([]*TagResponse, 0, len(tags))
		for _, tag := range tags {
			responses = append(responses, toTagResponse(&tag))
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "获取用户标签成功",
			"data":    responses,
		})
	}
}

// removeUserTagHandler 移除用户标签
// @Summary 移除用户标签
// @Description 从用户移除一个标签
// @Tags 用户标签管理
// @Param id path int true "用户ID"
// @Param tagId path int true "标签ID"
// @Success 200 {object} ApiResponse
// @Router /admin/users/{id}/tags/{tagId} [delete]
func removeUserTagHandler(s *user.UserTagService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "无效的用户ID",
				"error":   err.Error(),
			})
			return
		}

		tagID, err := strconv.ParseUint(c.Param("tagId"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "无效的标签ID",
				"error":   err.Error(),
			})
			return
		}

		if err := s.RemoveTagFromUser(c.Request.Context(), userID, tagID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "移除标签失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "移除标签成功",
		})
	}
}

// getUsersByTagHandler 获取拥有指定标签的用户列表
// @Summary 获取标签下的用户
// @Description 获取拥有指定标签的用户列表（分页）
// @Tags 用户标签管理
// @Produce json
// @Param id path int true "标签ID"
// @Param page query int false "页码（默认1）"
// @Param page_size query int false "每页数量（默认10）"
// @Success 200 {object} ApiResponse{data=[]model.User,pagination=object}
// @Router /admin/user-tags/{id}/users [get]
func getUsersByTagHandler(s *user.UserTagService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tagID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "无效的标签ID",
				"error":   err.Error(),
			})
			return
		}

		page, pageSize := getPaginationParams(c)

		users, total, err := s.GetUsersByTag(c.Request.Context(), tagID, page, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "获取用户列表失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "获取用户列表成功",
			"data":    users,
			"pagination": gin.H{
				"page":      page,
				"page_size": pageSize,
				"total":     total,
			},
		})
	}
}

// toTagResponse 转换Model到响应
func toTagResponse(tag *model.UserTag) *TagResponse {
	return &TagResponse{
		ID:          tag.ID,
		Name:        tag.Name,
		Color:       tag.Color,
		Description: tag.Description,
		CreatedAt:   tag.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   tag.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// getPaginationParams 获取分页参数
func getPaginationParams(c *gin.Context) (page, pageSize int) {
	page = 1
	pageSize = 10

	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = p
	}
	if ps, err := strconv.Atoi(c.Query("page_size")); err == nil && ps > 0 {
		pageSize = ps
	}

	return
}

// BatchDeleteTagsRequest 批量删除标签请求
type BatchDeleteTagsRequest struct {
	IDs []uint64 `json:"ids" binding:"required,min=1"`
}

// batchDeleteTagsHandler 批量删除标签
// @Summary 批量删除用户标签
// @Description 批量删除用户标签（同时会移除所有用户的这些标签）
// @Tags 用户标签管理
// @Accept json
// @Produce json
// @Param body body BatchDeleteTagsRequest true "标签ID列表"
// @Success 200 {object} ApiResponse{data=map[string]interface{}}
// @Router /admin/user-tags/batch [delete]
func batchDeleteTagsHandler(s *user.UserTagService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req BatchDeleteTagsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "参数验证失败",
				"error":   err.Error(),
			})
			return
		}

		result, err := s.BatchDeleteTags(c.Request.Context(), req.IDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "批量删除标签失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "批量删除标签成功",
			"data":    result,
		})
	}
}

// BatchDeleteUserTagsRequest 批量删除用户标签请求
type BatchDeleteUserTagsRequest struct {
	TagIDs []uint64 `json:"tag_ids" binding:"required,min=1,max=100"`
}

// BatchDeleteUserTags 批量删除用户标签
// @Summary      批量删除用户标签
// @Description  批量删除用户标签（同时会移除所有用户的这些标签）
// @Tags         用户标签管理
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchDeleteUserTagsRequest  true  "标签ID列表"
// @Success      200  {object}  model.APIResponse[BatchOperationResponse]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/user-tags/batch/delete [post]
func BatchDeleteUserTags(s *user.UserTagService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req BatchDeleteUserTagsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
			return
		}

		if len(req.TagIDs) == 0 {
			respondAPIError(c, apierr.BadRequest("tag_ids is required"))
			return
		}
		if len(req.TagIDs) > 100 {
			respondAPIError(c, apierr.BadRequest("maximum 100 tags can be deleted at once"))
			return
		}

		result, err := s.BatchDeleteTags(c.Request.Context(), req.TagIDs)
		if err != nil {
			respondAPIError(c, apierr.InternalError("batch delete tags failed").WithDetails(err.Error()))
			return
		}

		respondSuccess(c, result)
	}
}

// BatchAssignTagsToUsersRequest 批量分配标签给用户请求
type BatchAssignTagsToUsersRequest struct {
	UserIDs []uint64 `json:"user_ids" binding:"required,min=1,max=100"`
	TagIDs  []uint64 `json:"tag_ids" binding:"required,min=1,max=100"`
}

// BatchAssignTagsToUsers 批量分配标签给用户
// @Summary      批量分配标签给用户
// @Description  批量为多个用户分配多个标签（笛卡尔积：每个用户都会获得所有指定的标签）
// @Tags         用户标签管理
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchAssignTagsToUsersRequest  true  "用户ID列表和标签ID列表"
// @Success      200  {object}  model.APIResponse[BatchOperationResponse]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/user-tags/batch/assign [post]
func BatchAssignTagsToUsers(s *user.UserTagService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req BatchAssignTagsToUsersRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
			return
		}

		if len(req.UserIDs) == 0 {
			respondAPIError(c, apierr.BadRequest("user_ids is required"))
			return
		}
		if len(req.UserIDs) > 100 {
			respondAPIError(c, apierr.BadRequest("maximum 100 users per batch"))
			return
		}
		if len(req.TagIDs) == 0 {
			respondAPIError(c, apierr.BadRequest("tag_ids is required"))
			return
		}
		if len(req.TagIDs) > 100 {
			respondAPIError(c, apierr.BadRequest("maximum 100 tags per batch"))
			return
		}

		result, err := s.BatchAssignTagsToUsers(c.Request.Context(), req.UserIDs, req.TagIDs)
		if err != nil {
			respondAPIError(c, apierr.InternalError("batch assign tags failed").WithDetails(err.Error()))
			return
		}

		respondSuccess(c, result)
	}
}

// BatchRemoveTagsFromUsersRequest 批量移除用户标签请求
type BatchRemoveTagsFromUsersRequest struct {
	UserIDs []uint64 `json:"user_ids" binding:"required,min=1,max=100"`
	TagIDs  []uint64 `json:"tag_ids" binding:"required,min=1,max=100"`
}

// BatchRemoveTagsFromUsers 批量移除用户标签
// @Summary      批量移除用户标签
// @Description  批量移除多个用户的多个标签（笛卡尔积：移除所有指定用户的所有指定标签）
// @Tags         用户标签管理
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchRemoveTagsFromUsersRequest  true  "用户ID列表和标签ID列表"
// @Success      200  {object}  model.APIResponse[BatchOperationResponse]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/user-tags/batch/remove [post]
func BatchRemoveTagsFromUsers(s *user.UserTagService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req BatchRemoveTagsFromUsersRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
			return
		}

		if len(req.UserIDs) == 0 {
			respondAPIError(c, apierr.BadRequest("user_ids is required"))
			return
		}
		if len(req.UserIDs) > 100 {
			respondAPIError(c, apierr.BadRequest("maximum 100 users per batch"))
			return
		}
		if len(req.TagIDs) == 0 {
			respondAPIError(c, apierr.BadRequest("tag_ids is required"))
			return
		}
		if len(req.TagIDs) > 100 {
			respondAPIError(c, apierr.BadRequest("maximum 100 tags per batch"))
			return
		}

		result, err := s.BatchRemoveTagsFromUsers(c.Request.Context(), req.UserIDs, req.TagIDs)
		if err != nil {
			respondAPIError(c, apierr.InternalError("batch remove tags failed").WithDetails(err.Error()))
			return
		}

		respondSuccess(c, result)
	}
}

// BatchAssignTagsRequest 批量分配标签请求
type BatchAssignTagsRequest struct {
	UserTagPairs []UserTagPair `json:"userTagPairs" binding:"required,min=1"`
}

// UserTagPair 用户标签对
type UserTagPair struct {
	UserID uint64 `json:"userId" binding:"required"`
	TagID  uint64 `json:"tagId" binding:"required"`
}

// batchAssignTagsHandler 批量为用户分配标签（旧版本兼容）
// @Summary 批量为用户分配标签
// @Description 批量为多个用户分配标签
// @Tags 用户标签管理
// @Accept json
// @Produce json
// @Param body body BatchAssignTagsRequest true "用户标签对列表"
// @Success 200 {object} ApiResponse{data=map[string]interface{}}
// @Router /admin/users/user-tags/batch/assign [post]
func batchAssignTagsHandler(s *user.UserTagService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req BatchAssignTagsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "参数验证失败",
				"error":   err.Error(),
			})
			return
		}

		// Convert to service input
		pairs := make([]user.UserTagPair, len(req.UserTagPairs))
		for i, p := range req.UserTagPairs {
			pairs[i] = user.UserTagPair{
				UserID: p.UserID,
				TagID:  p.TagID,
			}
		}

		result, err := s.BatchAssignTags(c.Request.Context(), pairs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "批量分配标签失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "批量分配标签成功",
			"data":    result,
		})
	}
}

// ApiResponse 标准API响应
type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
