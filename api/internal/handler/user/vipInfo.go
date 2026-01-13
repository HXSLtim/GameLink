package user

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/repository"
	vipservice "gamelink/internal/service/vip"
	"gamelink/pkg/apierr"
)

// VipInfoHandler 用户 VIP 信息处理器
type VipInfoHandler struct {
	vipSvc   *vipservice.Service
	userRepo repository.UserRepository
}

// NewVipInfoHandler 创建用户 VIP 信息处理器
func NewVipInfoHandler(vipSvc *vipservice.Service, userRepo repository.UserRepository) *VipInfoHandler {
	return &VipInfoHandler{
		vipSvc:   vipSvc,
		userRepo: userRepo,
	}
}

// VipLevelInfo VIP 等级信息
type VipLevelInfo struct {
	ID    uint64 `json:"id"`
	Slug  string `json:"slug"`
	Name  string `json:"name"`
	Level int    `json:"level"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
}

// VipBenefit VIP 权益
type VipBenefit struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

// UserVipInfoResponse 用户 VIP 信息响应
type UserVipInfoResponse struct {
	VipUnlocked             bool          `json:"vipUnlocked"`
	CurrentLevel            *VipLevelInfo `json:"currentLevel,omitempty"`
	CurrentExp              int64         `json:"currentExp"`
	NextLevelExp            int64         `json:"nextLevelExp"`
	ExpProgress             float64       `json:"expProgress"`
	VipUnlockedAt           *string       `json:"vipUnlockedAt,omitempty"`
	VipExpireAt             *string       `json:"vipExpireAt,omitempty"`
	Benefits                []VipBenefit  `json:"benefits"`
	MonthlyTicketsRemaining int           `json:"monthlyTicketsRemaining"`
	DiscountRate            int           `json:"discountRate"`
}

// GetUserVipInfo 获取用户 VIP 信息
// @Summary      获取用户 VIP 信息
// @Description  获取当前用户的 VIP 等级、经验值和权益信息
// @Tags         User - VIP
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  UserVipInfoResponse
// @Failure      401  {object}  apierr.APIError
// @Failure      500  {object}  apierr.APIError
// @Router       /user/vip/info [get]
func (h *VipInfoHandler) GetUserVipInfo(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	ctx := c.Request.Context()

	// 获取用户信息
	user, err := h.userRepo.Get(ctx, userID)
	if err != nil {
		resp.Error(c, apierr.InternalError("获取用户信息失败").WithDetails(err.Error()))
		return
	}

	response := UserVipInfoResponse{
		VipUnlocked: user.VipUnlocked,
		CurrentExp:  user.VipExp,
		Benefits:    []VipBenefit{},
	}

	// 如果 VIP 已解锁，获取等级信息
	if user.VipUnlocked && user.VipLevelID != nil {
		level, err := h.vipSvc.GetLevel(ctx, *user.VipLevelID)
		if err == nil && level != nil {
			response.CurrentLevel = &VipLevelInfo{
				ID:    level.ID,
				Slug:  level.Slug,
				Name:  level.Title,
				Level: level.SortOrder,
				Icon:  level.IconURL,
				Color: level.Color,
			}
			// 折扣率转换为百分比
			if level.OrderDiscount < 1.0 {
				response.DiscountRate = int((1.0 - level.OrderDiscount) * 100)
			}
			response.MonthlyTicketsRemaining = level.MonthlyCouponCount

			// 解析权益（Benefits 是 JSON 字符串）
			if level.Benefits != "" && level.Benefits != "{}" {
				response.Benefits = append(response.Benefits, VipBenefit{
					ID:   "b1",
					Name: "专属权益",
					Icon: "gift",
				})
			}
		}

		// 获取下一等级经验值
		nextLevel, err := h.vipSvc.GetLevelByExp(ctx, user.VipExp+1)
		if err == nil && nextLevel != nil && nextLevel.ID != *user.VipLevelID {
			response.NextLevelExp = nextLevel.ExpRequired
		} else {
			// 已是最高等级
			response.NextLevelExp = user.VipExp
		}

		// 计算进度
		if response.NextLevelExp > 0 && response.CurrentLevel != nil {
			currentMin := int64(0)
			if level != nil {
				currentMin = level.ExpRequired
			}
			totalNeeded := response.NextLevelExp - currentMin
			currentProgress := user.VipExp - currentMin
			if totalNeeded > 0 {
				response.ExpProgress = float64(currentProgress) / float64(totalNeeded)
			}
		}

		// 格式化时间
		if user.VipUnlockedAt != nil {
			t := user.VipUnlockedAt.Format("2006-01-02T15:04:05Z")
			response.VipUnlockedAt = &t
		}
		if user.VipExpireAt != nil {
			t := user.VipExpireAt.Format("2006-01-02T15:04:05Z")
			response.VipExpireAt = &t
		}
	} else {
		// 未解锁 VIP，获取解锁门槛
		consumeThreshold, rechargeThreshold, _ := h.vipSvc.GetUnlockThreshold(ctx)
		// 设置下一等级经验为解锁门槛
		if consumeThreshold > 0 {
			response.NextLevelExp = consumeThreshold
		} else if rechargeThreshold > 0 {
			response.NextLevelExp = rechargeThreshold
		}
	}

	resp.OK(c, response)
}

// RegisterVipInfoRoutes 注册 VIP 信息路由（追加到现有 VIP 路由）
func RegisterVipInfoRoutes(rg *gin.RouterGroup, vipSvc *vipservice.Service, userRepo repository.UserRepository, _ gin.HandlerFunc) {
	h := NewVipInfoHandler(vipSvc, userRepo)
	// 注意：这个路由会追加到 /user/vip 组下
	vipGroup := rg.Group("/vip")
	vipGroup.GET("/info", h.GetUserVipInfo)
}
