package public

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/model"
	"gamelink/pkg/apierr"

	"gorm.io/gorm"
)

// BannerHandler 公共 Banner 处理器
type BannerHandler struct {
	db *gorm.DB
}

// NewBannerHandler 创建公共 Banner 处理器
func NewBannerHandler(db *gorm.DB) *BannerHandler {
	return &BannerHandler{db: db}
}

// PublicBannerItem 公开的 banner 信息
type PublicBannerItem struct {
	ID          uint64           `json:"id"`
	Title       string           `json:"title,omitempty"`
	Description string           `json:"description,omitempty"`
	ImageURL    string           `json:"imageUrl"`
	Type        model.BannerType `json:"type"`
	Link        string           `json:"link,omitempty"`
	ActionText  string           `json:"actionText,omitempty"`
}

// BannerListResponse banner 列表响应
type BannerListResponse struct {
	Banners []PublicBannerItem `json:"banners"`
}

// ListBanners 获取首页 banner 列表（公开）
// @Summary 获取首页 banner 列表
// @Description 获取当前可见的 banner 列表，无需登录
// @Tags 公共-Banner
// @Produce json
// @Success 200 {object} BannerListResponse
// @Router /public/banners [get]
func (h *BannerHandler) ListBanners(c *gin.Context) {
	var banners []model.Banner
	err := h.db.
		Where("is_visible = ? AND deleted_at IS NULL", true).
		Order("sort_order ASC, id ASC").
		Find(&banners).Error
	if err != nil {
		resp.Error(c, apierr.InternalError("获取 banner 列表失败"))
		return
	}

	result := make([]PublicBannerItem, 0, len(banners))
	for _, b := range banners {
		item := PublicBannerItem{
			ID:          b.ID,
			Title:       b.Title,
			Description: b.Description,
			ImageURL:    b.ImageURL,
			Type:        b.Type,
			Link:        b.Link,
			ActionText:  b.ActionText,
		}
		result = append(result, item)
	}

	resp.OK(c, BannerListResponse{Banners: result})
}

// RegisterRoutes 注册公共 banner 路由
func (h *BannerHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/banners", h.ListBanners)
}
