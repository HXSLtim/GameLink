package payment

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// WeChatClient 微信支付客户端（简化版，用于演示核心逻辑）
type WeChatClient struct {
	appID     string
	mchID     string
	apiKey    string
	notifyURL string
}

// NewWeChatClient 创建微信支付客户端
func NewWeChatClient(appID, mchID, apiKey, notifyURL string) *WeChatClient {
	return &WeChatClient{
		appID:     appID,
		mchID:     mchID,
		apiKey:    apiKey,
		notifyURL: notifyURL,
	}
}

// GenerateSign 生成微信支付签名
func (c *WeChatClient) GenerateSign(params map[string]string) string {
	// 1. 过滤空值和 sign 字段
	filtered := make(map[string]string)
	for k, v := range params {
		if k != "sign" && v != "" {
			filtered[k] = v
		}
	}

	// 2. 按键名字典序排序
	var keys []string
	for k := range filtered {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 3. 拼接成字符串
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, filtered[k]))
	}
	signStr := strings.Join(parts, "&")

	// 4. 拼接 API 密钥
	signStr = signStr + "&key=" + c.apiKey

	// 5. MD5 加密并转大写
	h := md5.New()
	h.Write([]byte(signStr))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

// VerifySign 验证微信支付签名
func (c *WeChatClient) VerifySign(params map[string]string) bool {
	sign := params["sign"]
	if sign == "" {
		return false
	}

	calculatedSign := c.GenerateSign(params)
	return calculatedSign == sign
}

// BuildUnifiedOrderParams 构建统一下单参数
func (c *WeChatClient) BuildUnifiedOrderParams(orderID, body, clientIP, tradeType string, totalFee int64, openID string) map[string]string {
	params := map[string]string{
		"appid":            c.appID,
		"mch_id":           c.mchID,
		"nonce_str":        generateNonceStr(),
		"body":             body,
		"out_trade_no":     orderID,
		"total_fee":        fmt.Sprintf("%d", totalFee),
		"spbill_create_ip": clientIP,
		"notify_url":       c.notifyURL,
		"trade_type":       tradeType,
	}

	if openID != "" {
		params["openid"] = openID
	}

	// 生成签名
	params["sign"] = c.GenerateSign(params)

	return params
}

// BuildOrderQueryParams 构建订单查询参数
func (c *WeChatClient) BuildOrderQueryParams(outTradeNo, transactionID string) map[string]string {
	params := map[string]string{
		"appid":        c.appID,
		"mch_id":       c.mchID,
		"nonce_str":    generateNonceStr(),
		"out_trade_no": outTradeNo,
	}

	if transactionID != "" {
		params["transaction_id"] = transactionID
	}

	// 生成签名
	params["sign"] = c.GenerateSign(params)

	return params
}

// BuildRefundParams 构建退款参数
func (c *WeChatClient) BuildRefundParams(outTradeNo, outRefundNo string, totalFee, refundFee int64) map[string]string {
	params := map[string]string{
		"appid":         c.appID,
		"mch_id":        c.mchID,
		"nonce_str":     generateNonceStr(),
		"out_trade_no":  outTradeNo,
		"out_refund_no": outRefundNo,
		"total_fee":     fmt.Sprintf("%d", totalFee),
		"refund_fee":    fmt.Sprintf("%d", refundFee),
	}

	// 生成签名
	params["sign"] = c.GenerateSign(params)

	return params
}
