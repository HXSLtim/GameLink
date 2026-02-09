package public

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
)

// PublicChatHandler handles public chat channels.
type PublicChatHandler struct {
	groups repository.ChatGroupRepository
}

// NewPublicChatHandler creates a public chat handler.
func NewPublicChatHandler(groups repository.ChatGroupRepository) *PublicChatHandler {
	return &PublicChatHandler{groups: groups}
}

// RegisterRoutes registers public chat routes.
func (h *PublicChatHandler) RegisterRoutes(rg *gin.RouterGroup) {
	chat := rg.Group("/chat")
	chat.GET("/public-channels", h.ListPublicChannels)
}

// ListPublicChannels returns public channels list.
// @Summary 获取公共频道列表
// @Description 获取公共频道列表（大厅、游戏频道等）
// @Tags 公共-聊天
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} resp.PagedResponse
// @Failure 500 {object} apierr.APIError
// @Router /public/chat/public-channels [get]
func (h *PublicChatHandler) ListPublicChannels(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	groups, total, err := h.groups.ListPublicChannels(c.Request.Context(), page, pageSize)
	if err != nil {
		resp.Error(c, apierr.InternalError("获取公共频道失败"))
		return
	}
	resp.List(c, groups, &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	})
}
