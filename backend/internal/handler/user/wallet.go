package user

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/apierr"
	"gamelink/internal/model"
	"gamelink/internal/service/wallet"
)

// RegisterWalletRoutes 注册钱包路由
func RegisterWalletRoutes(router gin.IRouter, svc *wallet.Service, auth gin.HandlerFunc) {
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

func rechargeHandler(c *gin.Context, svc *wallet.Service) {
	userID := getUserIDFromContext(c)
	var req rechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidJSONPayload).WithDetails(err.Error()))
		return
	}
	resp, err := svc.Recharge(c.Request.Context(), userID, wallet.RechargeRequest{
		AmountCents: req.AmountCents,
		Method:      req.Method,
	})
	if err != nil {
		if err == wallet.ErrInvalidAmount {
			respondAPIError(c, apierr.BadRequest(err.Error()))
			return
		}
		respondAPIError(c, apierr.InternalError("recharge failed").WithDetails(err.Error()))
		return
	}
	respondSuccess(c, "OK", resp)
}

func getBalanceHandler(c *gin.Context, svc *wallet.Service) {
	userID := getUserIDFromContext(c)
	w, err := svc.GetBalance(c.Request.Context(), userID)
	if err != nil {
		respondAPIError(c, apierr.InternalError("get balance failed").WithDetails(err.Error()))
		return
	}
	respondSuccess(c, "OK", w)
}
