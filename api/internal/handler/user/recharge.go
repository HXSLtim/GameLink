package user

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/repository"
	rechargeservice "gamelink/internal/service/recharge"
	"gamelink/pkg/apierr"
)

// RechargeHandler 用户充值处理器
type RechargeHandler struct {
	svc *rechargeservice.Service
}

// NewRechargeHandler 创建用户充值处理器
func NewRechargeHandler(svc *rechargeservice.Service) *RechargeHandler {
	return &RechargeHandler{svc: svc}
}

// ListOptions 获取充值档位列表
func (h *RechargeHandler) ListOptions(c *gin.Context) {
	// TODO: 从用户信息获取 VIP 等级
	var vipLevel *uint64

	options, err := h.svc.GetActiveOptions(c.Request.Context(), vipLevel)
	if err != nil {
		resp.Error(c, apierr.InternalError("获取充值档位失败").WithDetails(err.Error()))
		return
	}

	resp.OK(c, options)
}

// GetOption 获取档位详情
func (h *RechargeHandler) GetOption(c *gin.Context) {
	id, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	option, err := h.svc.GetOption(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			resp.Error(c, apierr.NotFound("档位不存在"))
			return
		}
		resp.Error(c, apierr.InternalError("获取档位失败").WithDetails(err.Error()))
		return
	}

	// 检查是否启用
	if !option.IsActive {
		resp.Error(c, apierr.NotFound("档位不存在"))
		return
	}

	resp.OK(c, option)
}

// CreateOrderRequest 创建充值订单请求
type CreateOrderRequest struct {
	OptionID       uint64 `json:"optionId" binding:"required"`
	PaymentChannel string `json:"paymentChannel" binding:"required,oneof=wechat alipay"`
	PaymentMethod  string `json:"paymentMethod"`
}

// CreateOrder 创建充值订单
func (h *RechargeHandler) CreateOrder(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("无效的请求参数").WithDetails(err.Error()))
		return
	}

	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	record, err := h.svc.CreateRechargeOrder(c.Request.Context(), userID, req.OptionID, req.PaymentChannel, req.PaymentMethod, clientIP, userAgent)
	if err != nil {
		resp.Error(c, apierr.BadRequest(err.Error()))
		return
	}

	// TODO: 调用支付渠道获取支付参数
	resp.Created(c, gin.H{
		"record":  record,
		"payInfo": nil, // 支付参数，需要对接支付渠道
	})
}

// GetOrder 获取充值订单详情
func (h *RechargeHandler) GetOrder(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	id, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	record, err := h.svc.GetRecord(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			resp.Error(c, apierr.NotFound("充值记录不存在"))
			return
		}
		resp.Error(c, apierr.InternalError("获取充值记录失败").WithDetails(err.Error()))
		return
	}

	// 验证所有权
	if record.UserID != userID {
		resp.Error(c, apierr.NotFound("充值记录不存在"))
		return
	}

	resp.OK(c, record)
}

// ListMyRecords 获取我的充值记录
func (h *RechargeHandler) ListMyRecords(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	limit := 20
	if v := c.Query("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	records, err := h.svc.GetUserRecords(c.Request.Context(), userID, limit)
	if err != nil {
		resp.Error(c, apierr.InternalError("获取充值记录失败").WithDetails(err.Error()))
		return
	}

	resp.OK(c, records)
}

// RegisterRechargeRoutes 注册用户充值路由
func RegisterRechargeRoutes(rg *gin.RouterGroup, svc *rechargeservice.Service, _ gin.HandlerFunc) {
	h := NewRechargeHandler(svc)

	rechargeGroup := rg.Group("/recharge")
	{
		rechargeGroup.GET("/options", h.ListOptions)
		rechargeGroup.GET("/options/:id", h.GetOption)
		rechargeGroup.POST("/orders", h.CreateOrder)
		rechargeGroup.GET("/orders/:id", h.GetOrder)
		rechargeGroup.GET("/records", h.ListMyRecords)
	}
}
