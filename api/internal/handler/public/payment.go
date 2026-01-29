package public

import (
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/service/external"
	paymentservice "gamelink/internal/service/payment"
)

// PaymentHandler handles payment callbacks.
type PaymentHandler struct {
	svc            *paymentservice.PaymentService
	cfg            *external.Config
	wechatProvider *paymentservice.RealWeChatProvider
	alipayProvider *paymentservice.RealAlipayProvider
}

// NewPaymentHandler creates a new payment handler.
func NewPaymentHandler(svc *paymentservice.PaymentService, cfg *external.Config) *PaymentHandler {
	handler := &PaymentHandler{
		svc: svc,
		cfg: cfg,
	}
	if cfg != nil {
		handler.wechatProvider = paymentservice.NewWeChatProvider(cfg)
		if alipayProvider, err := paymentservice.NewAlipayProvider(cfg); err == nil {
			handler.alipayProvider = alipayProvider
		}
	}
	return handler
}

// RegisterRoutes registers payment callback routes.
func (h *PaymentHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/payments/wechat/notify", h.WeChatNotify)
	router.POST("/payments/alipay/notify", h.AlipayNotify)
}

// WeChatNotify handles WeChat Pay notifications.
func (h *PaymentHandler) WeChatNotify(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		writeWeChatResponse(c, "FAIL", "read body failed")
		return
	}

	params, err := paymentservice.ParseWeChatXML(body)
	if err != nil {
		writeWeChatResponse(c, "FAIL", "invalid xml")
		return
	}

	if params["return_code"] != "SUCCESS" {
		writeWeChatResponse(c, "FAIL", params["return_msg"])
		return
	}
	if h.wechatProvider == nil || h.cfg == nil || !h.cfg.WeChatPay.Enabled {
		writeWeChatResponse(c, "FAIL", "wechat disabled")
		return
	}
	if !h.wechatProvider.VerifySign(params) {
		writeWeChatResponse(c, "FAIL", "invalid sign")
		return
	}
	if params["result_code"] != "SUCCESS" {
		writeWeChatResponse(c, "SUCCESS", "OK")
		return
	}

	paymentID, err := strconv.ParseUint(strings.TrimSpace(params["out_trade_no"]), 10, 64)
	if err != nil {
		writeWeChatResponse(c, "FAIL", "invalid out_trade_no")
		return
	}

	data := map[string]interface{}{
		"payment_id": paymentID,
		"trade_no":   params["transaction_id"],
	}
	if totalFee := strings.TrimSpace(params["total_fee"]); totalFee != "" {
		if amount, parseErr := strconv.ParseInt(totalFee, 10, 64); parseErr == nil {
			data["amount_cents"] = amount
		}
	}

	if err := h.svc.HandlePaymentCallback(c.Request.Context(), "wechat", data); err != nil {
		writeWeChatResponse(c, "FAIL", "handle callback failed")
		return
	}
	writeWeChatResponse(c, "SUCCESS", "OK")
}

// AlipayNotify handles Alipay notifications.
func (h *PaymentHandler) AlipayNotify(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		c.String(http.StatusOK, "fail")
		return
	}
	if h.alipayProvider == nil || h.cfg == nil || !h.cfg.Alipay.Enabled {
		c.String(http.StatusOK, "fail")
		return
	}

	params := map[string]string{}
	for k, v := range c.Request.PostForm {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	sign := params["sign"]
	if sign == "" {
		c.String(http.StatusOK, "fail")
		return
	}
	delete(params, "sign")

	if !h.alipayProvider.VerifySign(params, sign) {
		c.String(http.StatusOK, "fail")
		return
	}

	tradeStatus := params["trade_status"]
	if tradeStatus != "TRADE_SUCCESS" && tradeStatus != "TRADE_FINISHED" {
		c.String(http.StatusOK, "success")
		return
	}

	paymentID, err := strconv.ParseUint(strings.TrimSpace(params["out_trade_no"]), 10, 64)
	if err != nil {
		c.String(http.StatusOK, "fail")
		return
	}

	data := map[string]interface{}{
		"payment_id": paymentID,
		"trade_no":   params["trade_no"],
	}
	if totalAmount := strings.TrimSpace(params["total_amount"]); totalAmount != "" {
		if amount, parseErr := strconv.ParseFloat(totalAmount, 64); parseErr == nil {
			data["amount_cents"] = int64(math.Round(amount * 100))
		}
	}

	if err := h.svc.HandlePaymentCallback(c.Request.Context(), "alipay", data); err != nil {
		c.String(http.StatusOK, "fail")
		return
	}
	c.String(http.StatusOK, "success")
}

func writeWeChatResponse(c *gin.Context, code, msg string) {
	c.Header("Content-Type", "application/xml")
	c.String(http.StatusOK, "<xml><return_code><![CDATA[%s]]></return_code><return_msg><![CDATA[%s]]></return_msg></xml>", code, msg)
}
