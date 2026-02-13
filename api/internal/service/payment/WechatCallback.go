package payment

import (
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"sort"
	"strings"
)

// WeChatCallbackNotification 微信支付回调通知
type WeChatCallbackNotification struct {
	XMLName            xml.Name `xml:"xml"`
	ReturnCode         string   `xml:"return_code"`          // SUCCESS/FAIL
	ReturnMsg          string   `xml:"return_msg"`           // 返回信息
	AppID              string   `xml:"appid"`                // 微信分配的公众账号ID
	MchID              string   `xml:"mch_id"`               // 微信支付分配的商户号
	DeviceInfo         string   `xml:"device_info"`          // 微信支付分配的终端设备号
	NonceStr           string   `xml:"nonce_str"`            // 随机字符串
	Sign               string   `xml:"sign"`                 // 签名
	ResultCode         string   `xml:"result_code"`          // SUCCESS/FAIL
	ErrCode            string   `xml:"err_code"`             // 错误代码
	ErrCodeDes         string   `xml:"err_code_des"`         // 错误代码描述
	OpenID             string   `xml:"openid"`               // 用户在商户appid下的唯一标识
	IsSubscribe        string   `xml:"is_subscribe"`         // 是否关注公众账号
	TradeType          string   `xml:"trade_type"`           // 交易类型: JSAPI, NATIVE, APP
	BankType           string   `xml:"bank_type"`            // 银行类型
	TotalFee           string   `xml:"total_fee"`            // 订单总金额，单位为分
	SettlementTotalFee string   `xml:"settlement_total_fee"` // 应结订单金额=订单金额-非充值代金券金额，单位为分
	FeeType            string   `xml:"fee_type"`             // 货币类型
	CashFee            string   `xml:"cash_fee"`             // 现金支付金额
	CashFeeType        string   `xml:"cash_fee_type"`        // 现金支付货币类型
	CouponFee          string   `xml:"coupon_fee"`           // 代金券金额=订单金额-现金支付金额
	CouponCount        string   `xml:"coupon_count"`         // 代金券使用数量
	TransactionID      string   `xml:"transaction_id"`       // 微信支付订单号
	OutTradeNo         string   `xml:"out_trade_no"`         // 商户订单号
	Attach             string   `xml:"attach"`               // 附加数据，原样返回
	TimeEnd            string   `xml:"time_end"`             // 支付完成时间
}

// WeChatCallbackResponse 微信支付回调响应
type WeChatCallbackResponse struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
}

// NewWeChatCallbackResponse 创建成功回调响应
func NewWeChatCallbackResponse() *WeChatCallbackResponse {
	return &WeChatCallbackResponse{
		ReturnCode: "SUCCESS",
		ReturnMsg:  "OK",
	}
}

// NewWeChatCallbackFailResponse 创建失败回调响应
func NewWeChatCallbackFailResponse(msg string) *WeChatCallbackResponse {
	return &WeChatCallbackResponse{
		ReturnCode: "FAIL",
		ReturnMsg:  msg,
	}
}

// ToXML 转换为 XML 格式
func (r *WeChatCallbackResponse) ToXML() ([]byte, error) {
	return xml.Marshal(r)
}

// WeChatCallbackHandler 微信支付回调处理器
type WeChatCallbackHandler struct {
	appID       string
	mchID       string
	apiKey      string
	replayCache map[string]int64 // 防重放缓存
}

// NewWeChatCallbackHandler 创建微信支付回调处理器
func NewWeChatCallbackHandler(appID, mchID, apiKey string) *WeChatCallbackHandler {
	return &WeChatCallbackHandler{
		appID:       appID,
		mchID:       mchID,
		apiKey:      apiKey,
		replayCache: make(map[string]int64),
	}
}

// ParseCallback 解析回调通知
func (h *WeChatCallbackHandler) ParseCallback(body []byte) (*WeChatCallbackNotification, error) {
	var notification WeChatCallbackNotification
	if err := xml.Unmarshal(body, &notification); err != nil {
		return nil, fmt.Errorf("unmarshal callback: %w", err)
	}

	// 验证基本参数
	if notification.ReturnCode != "SUCCESS" {
		return nil, fmt.Errorf("return_code is not SUCCESS: %s", notification.ReturnCode)
	}

	if notification.ResultCode != "SUCCESS" {
		return nil, fmt.Errorf("result_code is not SUCCESS: %s - %s", notification.ErrCode, notification.ErrCodeDes)
	}

	if notification.OutTradeNo == "" {
		return nil, fmt.Errorf("out_trade_no is empty")
	}

	if notification.TransactionID == "" {
		return nil, fmt.Errorf("transaction_id is empty")
	}

	return &notification, nil
}

// VerifyCallback 验证回调签名
func (h *WeChatCallbackHandler) VerifyCallback(notification *WeChatCallbackNotification) bool {
	// 构建签名参数
	params := map[string]string{
		"return_code": notification.ReturnCode,
		"return_msg":  notification.ReturnMsg,
	}

	if notification.AppID != "" {
		params["appid"] = notification.AppID
	}
	if notification.MchID != "" {
		params["mch_id"] = notification.MchID
	}
	if notification.DeviceInfo != "" {
		params["device_info"] = notification.DeviceInfo
	}
	if notification.NonceStr != "" {
		params["nonce_str"] = notification.NonceStr
	}
	if notification.ResultCode != "" {
		params["result_code"] = notification.ResultCode
	}
	if notification.ErrCode != "" {
		params["err_code"] = notification.ErrCode
	}
	if notification.ErrCodeDes != "" {
		params["err_code_des"] = notification.ErrCodeDes
	}
	if notification.OpenID != "" {
		params["openid"] = notification.OpenID
	}
	if notification.IsSubscribe != "" {
		params["is_subscribe"] = notification.IsSubscribe
	}
	if notification.TradeType != "" {
		params["trade_type"] = notification.TradeType
	}
	if notification.BankType != "" {
		params["bank_type"] = notification.BankType
	}
	if notification.TotalFee != "" {
		params["total_fee"] = notification.TotalFee
	}
	if notification.FeeType != "" {
		params["fee_type"] = notification.FeeType
	}
	if notification.CashFee != "" {
		params["cash_fee"] = notification.CashFee
	}
	if notification.CashFeeType != "" {
		params["cash_fee_type"] = notification.CashFeeType
	}
	if notification.CouponFee != "" {
		params["coupon_fee"] = notification.CouponFee
	}
	if notification.CouponCount != "" {
		params["coupon_count"] = notification.CouponCount
	}
	if notification.TransactionID != "" {
		params["transaction_id"] = notification.TransactionID
	}
	if notification.OutTradeNo != "" {
		params["out_trade_no"] = notification.OutTradeNo
	}
	if notification.Attach != "" {
		params["attach"] = notification.Attach
	}
	if notification.TimeEnd != "" {
		params["time_end"] = notification.TimeEnd
	}

	// 生成签名
	expectedSign := h.generateSign(params)

	// 验证签名
	return expectedSign == notification.Sign
}

// generateSign 生成签名（MD5 + 字典排序）
func (h *WeChatCallbackHandler) generateSign(params map[string]string) string {
	// 1. 过滤空值和 sign 字段
	filteredParams := make(map[string]string)
	for k, v := range params {
		if v != "" && k != "sign" {
			filteredParams[k] = v
		}
	}

	// 2. 按字典序排序参数
	keys := make([]string, 0, len(filteredParams))
	for k := range filteredParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 3. 拼接参数字符串
	var buf strings.Builder
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(filteredParams[k])
	}

	// 4. 追加 API 密钥
	buf.WriteString("&key=")
	buf.WriteString(h.apiKey)

	// 5. MD5 哈希并转大写
	hash := md5.Sum([]byte(buf.String()))
	return strings.ToUpper(fmt.Sprintf("%x", hash))
}

// CheckReplay 检查是否为重放攻击
func (h *WeChatCallbackHandler) CheckReplay(notification *WeChatCallbackNotification) bool {
	// 使用 TransactionID 作为唯一标识
	key := notification.TransactionID

	// 检查缓存
	if _, exists := h.replayCache[key]; exists {
		return true // 检测到重放
	}

	// 记录到缓存
	// 注意：实际应用中应该使用 Redis 等外部缓存，并设置过期时间
	h.replayCache[key] = 1

	return false
}

// ValidateCallback 完整的回调验证流程
func (h *WeChatCallbackHandler) ValidateCallback(body []byte) (*WeChatCallbackNotification, error) {
	// 1. 解析回调
	notification, err := h.ParseCallback(body)
	if err != nil {
		return nil, fmt.Errorf("parse callback failed: %w", err)
	}

	// 2. 验证签名
	if !h.VerifyCallback(notification) {
		return nil, fmt.Errorf("signature verification failed")
	}

	// 3. 检查重放攻击
	if h.CheckReplay(notification) {
		return nil, fmt.Errorf("replay attack detected")
	}

	return notification, nil
}

// GetCallbackUniqueKey 获取回调的唯一标识（用于幂等性处理）
func (h *WeChatCallbackHandler) GetCallbackUniqueKey(notification *WeChatCallbackNotification) string {
	return notification.TransactionID
}

// GetOrderID 获取商户订单号
func (h *WeChatCallbackHandler) GetOrderID(notification *WeChatCallbackNotification) string {
	return notification.OutTradeNo
}

// GetTransactionID 获取微信支付订单号
func (h *WeChatCallbackHandler) GetTransactionID(notification *WeChatCallbackNotification) string {
	return notification.TransactionID
}

// GetTotalFee 获取订单总金额（分）
func (h *WeChatCallbackHandler) GetTotalFee(notification *WeChatCallbackNotification) int64 {
	var fee int64
	fmt.Sscanf(notification.TotalFee, "%d", &fee)
	return fee
}

// GetPaymentTime 获取支付完成时间
func (h *WeChatCallbackHandler) GetPaymentTime(notification *WeChatCallbackNotification) string {
	return notification.TimeEnd
}

// GetAttach 获取附加数据
func (h *WeChatCallbackHandler) GetAttach(notification *WeChatCallbackNotification) string {
	return notification.Attach
}
