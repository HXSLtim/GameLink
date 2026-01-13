package user

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	walletservice "gamelink/internal/service/wallet"
	"gamelink/pkg/apierr"
)

// RegisterWalletRoutes 注册钱包路由
func RegisterWalletRoutes(router gin.IRouter, svc *walletservice.WalletService, auth gin.HandlerFunc) {
	group := router.Group("/wallet")
	group.Use(auth)
	{
		group.GET("/balance", func(c *gin.Context) { getBalanceHandler(c, svc) })
		group.POST("/recharge", func(c *gin.Context) { rechargeHandler(c, svc) })
	}
}

type rechargeRequest struct {
	AmountCents int64               `json:"amountCents" binding:"required,min=1"`
	Method      model.PaymentMethod `json:"method" binding:"required,oneof=wechat alipay"`
}

// RechargeResponse 充值响应
type RechargeResponse struct {
	PaymentID   uint64 `json:"paymentId" example:"12345"`
	AmountCents int64  `json:"amountCents" example:"10000"`
	PayURL      string `json:"payUrl,omitempty" example:"https://pay.example.com/xxx"`
}

// WalletBalance 钱包余额
type WalletBalance struct {
	BalanceCents int64 `json:"balanceCents" example:"50000"`
	FrozenCents  int64 `json:"frozenCents" example:"1000"`
}

// rechargeHandler 充值
// @Summary      钱包充值
// @Description  用户钱包充值，支持微信/支付宝
// @Tags         User - Wallet
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      rechargeRequest  true  "充值请求"
// @Success      200      {object}  RechargeResponse
// @Failure      400      {object}  apierr.APIError
// @Failure      401      {object}  apierr.APIError
// @Failure      500      {object}  apierr.APIError
// @Router       /user/wallet/recharge [post]
func rechargeHandler(c *gin.Context, svc *walletservice.WalletService) {
	userID := getUserIDFromContext(c)
	var req rechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidJSONPayload).WithDetails(err.Error()))
		return
	}
	resp, err := svc.Recharge(c.Request.Context(), userID, walletservice.RechargeRequest{
		AmountCents: req.AmountCents,
		Method:      req.Method,
	})
	if err != nil {
		if err == walletservice.ErrInvalidAmount {
			respondAPIError(c, apierr.BadRequest(err.Error()))
			return
		}
		respondAPIError(c, apierr.InternalError("recharge failed").WithDetails(err.Error()))
		return
	}
	respondSuccess(c, "OK", resp)
}

// getBalanceHandler 获取钱包余额
// @Summary      获取钱包余额
// @Description  获取当前用户的钱包余额信息
// @Tags         User - Wallet
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  WalletBalance
// @Failure      401  {object}  apierr.APIError
// @Failure      500  {object}  apierr.APIError
// @Router       /user/wallet/balance [get]
func getBalanceHandler(c *gin.Context, svc *walletservice.WalletService) {
	userID := getUserIDFromContext(c)
	w, err := svc.GetBalance(c.Request.Context(), userID)
	if err != nil {
		respondAPIError(c, apierr.InternalError("get balance failed").WithDetails(err.Error()))
		return
	}
	respondSuccess(c, "OK", w)
}
