// Package external 提供第三方服务集成
//
// 包含以下第三方服务：
// 1. 身份证实名验证 - IdentityVerifier
// 2. 支付服务 - 微信支付/支付宝 (见 payment package)
// 3. 短信服务 - SMS (见 sms package)
// 4. 对象存储 - OSS (见 oss package)
//
// TODO: 生产环境部署前需要完成以下接口对接：
// - [ ] 身份证实名验证接口 (阿里云/腾讯云身份证二要素验证)
// - [ ] 微信支付回调验签
// - [ ] 支付宝支付回调验签
// - [ ] 短信发送接口 (阿里云/腾讯云 SMS)
package external

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

// ============================================================================
// 第三方服务模式配置
// ============================================================================

// ThirdPartyMode 第三方服务模式
type ThirdPartyMode string

const (
	// ModeMock 模拟模式 - 所有验证自动通过，用于开发测试
	ModeMock ThirdPartyMode = "mock"
	// ModeReal 真实模式 - 调用真实第三方API
	ModeReal ThirdPartyMode = "real"
)

var (
	ErrIdentityVerificationNotImplemented = errors.New("identity verification integration is not implemented in production")
	ErrPaymentCallbackNotImplemented      = errors.New("payment callback verification integration is not implemented in production")
)

func isProductionEnv() bool {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	return env == "production" || env == "prod"
}

// ThirdPartyConfig 第三方服务配置
type ThirdPartyConfig struct {
	// Mode 运行模式: mock(模拟) / real(真实)
	// 开发环境默认 mock，生产环境应设为 real
	Mode ThirdPartyMode `yaml:"mode" json:"mode"`

	// Identity 身份证验证配置
	Identity IdentityConfig `yaml:"identity" json:"identity"`
}

// IdentityConfig 身份证验证配置
type IdentityConfig struct {
	// Provider 服务提供商: aliyun / tencent / mock
	Provider string `yaml:"provider" json:"provider"`
	// AppID 应用ID
	AppID string `yaml:"app_id" json:"app_id"`
	// AppSecret 应用密钥
	AppSecret string `yaml:"app_secret" json:"app_secret"`
	// Endpoint API端点
	Endpoint string `yaml:"endpoint" json:"endpoint"`
}

// ============================================================================
// 身份证实名验证接口
// ============================================================================

// IdentityVerifyResult 身份证验证结果
type IdentityVerifyResult struct {
	// Verified 是否验证通过
	Verified bool `json:"verified"`
	// Score 置信度分数 (0-100)
	Score int `json:"score"`
	// Message 验证消息
	Message string `json:"message"`
	// RequestID 请求ID (用于追踪)
	RequestID string `json:"request_id"`
	// VerifiedAt 验证时间
	VerifiedAt time.Time `json:"verified_at"`
}

// IdentityVerifier 身份证实名验证接口
//
// TODO: 生产环境需要实现以下提供商之一：
// - 阿里云身份证二要素验证: https://market.aliyun.com/products/57000002/cmapi00035097.html
// - 腾讯云身份证二要素验证: https://cloud.tencent.com/document/product/1007/33188
type IdentityVerifier interface {
	// VerifyIdentity 验证身份证姓名与号码是否匹配
	// realName: 真实姓名
	// idCardNo: 身份证号码
	VerifyIdentity(ctx context.Context, realName, idCardNo string) (*IdentityVerifyResult, error)
}

// ============================================================================
// Mock 实现 - 开发测试用
// ============================================================================

// MockIdentityVerifier 模拟身份证验证器 (万能通过)
//
// 功能：
// - 基础格式校验（身份证号格式）
// - 所有格式正确的验证都返回通过
// - 用于开发和测试环境
type MockIdentityVerifier struct{}

// NewMockIdentityVerifier 创建模拟身份证验证器
func NewMockIdentityVerifier() *MockIdentityVerifier {
	return &MockIdentityVerifier{}
}

// VerifyIdentity 模拟验证身份证
// 规则：
// - 身份证号格式正确即通过
// - 姓名长度 >= 2 即通过
func (v *MockIdentityVerifier) VerifyIdentity(ctx context.Context, realName, idCardNo string) (*IdentityVerifyResult, error) {
	requestID := fmt.Sprintf("mock_%d", time.Now().UnixNano())
	result := &IdentityVerifyResult{
		RequestID:  requestID,
		VerifiedAt: time.Now(),
	}

	// 基础校验
	if len(realName) < 2 {
		result.Verified = false
		result.Score = 0
		result.Message = "姓名长度不符合要求"
		return result, nil
	}

	// 身份证号格式校验 (18位)
	idCardRegex := regexp.MustCompile(`^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$`)
	if !idCardRegex.MatchString(idCardNo) {
		result.Verified = false
		result.Score = 0
		result.Message = "身份证号格式不正确"
		return result, nil
	}

	// 模拟模式：格式正确即通过
	result.Verified = true
	result.Score = 100
	result.Message = "[MOCK] 身份证验证通过（模拟模式）"

	return result, nil
}

// FailClosedIdentityVerifier 生产环境下的失败保护实现
type FailClosedIdentityVerifier struct{}

// VerifyIdentity 在生产环境下未接入实名服务时显式失败
func (v *FailClosedIdentityVerifier) VerifyIdentity(ctx context.Context, realName, idCardNo string) (*IdentityVerifyResult, error) {
	return nil, ErrIdentityVerificationNotImplemented
}

// ============================================================================
// 真实接口实现 - TODO: 生产环境实现
// ============================================================================

// AliyunIdentityVerifier 阿里云身份证验证器
//
// TODO: 生产环境实现
// 参考文档: https://market.aliyun.com/products/57000002/cmapi00035097.html
// 需要：
// 1. 注册阿里云账号并开通身份证二要素验证服务
// 2. 获取 AppCode/AppKey/AppSecret
// 3. 实现 API 调用和签名
type AliyunIdentityVerifier struct {
	config IdentityConfig
}

// NewAliyunIdentityVerifier 创建阿里云身份证验证器
func NewAliyunIdentityVerifier(config IdentityConfig) *AliyunIdentityVerifier {
	return &AliyunIdentityVerifier{config: config}
}

// VerifyIdentity 调用阿里云身份证验证API
func (v *AliyunIdentityVerifier) VerifyIdentity(ctx context.Context, realName, idCardNo string) (*IdentityVerifyResult, error) {
	if isProductionEnv() {
		return nil, fmt.Errorf("%w: aliyun", ErrIdentityVerificationNotImplemented)
	}

	// TODO: 实现阿里云身份证二要素验证 API 调用
	//
	// API 调用示例：
	// POST https://idcert.market.alicloudapi.com/idcard
	// Headers:
	//   Authorization: APPCODE {appcode}
	// Body:
	//   idCard={身份证号}&name={姓名}
	//
	// 响应示例：
	// {
	//   "error_code": 0,
	//   "reason": "成功",
	//   "result": {
	//     "realname": "张三",
	//     "idcard": "110101199003076019",
	//     "res": 1  // 1=一致, 2=不一致
	//   }
	// }

	// 当前返回 Mock 结果
	return NewMockIdentityVerifier().VerifyIdentity(ctx, realName, idCardNo)
}

// TencentIdentityVerifier 腾讯云身份证验证器
//
// TODO: 生产环境实现
// 参考文档: https://cloud.tencent.com/document/product/1007/33188
type TencentIdentityVerifier struct {
	config IdentityConfig
}

// NewTencentIdentityVerifier 创建腾讯云身份证验证器
func NewTencentIdentityVerifier(config IdentityConfig) *TencentIdentityVerifier {
	return &TencentIdentityVerifier{config: config}
}

// VerifyIdentity 调用腾讯云身份证验证API
func (v *TencentIdentityVerifier) VerifyIdentity(ctx context.Context, realName, idCardNo string) (*IdentityVerifyResult, error) {
	if isProductionEnv() {
		return nil, fmt.Errorf("%w: tencent", ErrIdentityVerificationNotImplemented)
	}

	// TODO: 实现腾讯云身份核验 API 调用
	//
	// 需要使用腾讯云 SDK: tencentcloud-sdk-go
	// 服务: faceid.tencentcloudapi.com
	// 接口: IdCardVerification

	// 当前返回 Mock 结果
	return NewMockIdentityVerifier().VerifyIdentity(ctx, realName, idCardNo)
}

// ============================================================================
// 工厂函数
// ============================================================================

// NewIdentityVerifier 根据配置创建身份证验证器
func NewIdentityVerifier(config *ThirdPartyConfig) IdentityVerifier {
	if config == nil || config.Mode == ModeMock {
		if isProductionEnv() {
			return &FailClosedIdentityVerifier{}
		}
		return NewMockIdentityVerifier()
	}

	switch config.Identity.Provider {
	case "aliyun":
		return NewAliyunIdentityVerifier(config.Identity)
	case "tencent":
		return NewTencentIdentityVerifier(config.Identity)
	default:
		// 非生产默认使用 Mock；生产环境必须显式失败，避免误用
		if isProductionEnv() {
			return &FailClosedIdentityVerifier{}
		}
		return NewMockIdentityVerifier()
	}
}

// ============================================================================
// 支付回调验签 - TODO
// ============================================================================

// PaymentCallbackVerifier 支付回调验签接口
//
// TODO: 生产环境需要实现：
// - 微信支付回调验签: https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_7
// - 支付宝回调验签: https://opendocs.alipay.com/open/270/105902
type PaymentCallbackVerifier interface {
	// VerifyWeChatCallback 验证微信支付回调签名
	VerifyWeChatCallback(body []byte, signature string) (bool, error)
	// VerifyAlipayCallback 验证支付宝回调签名
	VerifyAlipayCallback(params map[string]string) (bool, error)
}

// FailClosedPaymentCallbackVerifier 生产环境下的失败保护实现
type FailClosedPaymentCallbackVerifier struct{}

// VerifyWeChatCallback 在生产环境下未接入验签时显式失败
func (v *FailClosedPaymentCallbackVerifier) VerifyWeChatCallback(body []byte, signature string) (bool, error) {
	return false, ErrPaymentCallbackNotImplemented
}

// VerifyAlipayCallback 在生产环境下未接入验签时显式失败
func (v *FailClosedPaymentCallbackVerifier) VerifyAlipayCallback(params map[string]string) (bool, error) {
	return false, ErrPaymentCallbackNotImplemented
}

// MockPaymentCallbackVerifier 模拟支付回调验签器 (万能通过)
type MockPaymentCallbackVerifier struct{}

// VerifyWeChatCallback 模拟验证微信回调
func (v *MockPaymentCallbackVerifier) VerifyWeChatCallback(body []byte, signature string) (bool, error) {
	if isProductionEnv() {
		return false, fmt.Errorf("%w: wechat", ErrPaymentCallbackNotImplemented)
	}

	// TODO: 生产环境实现真实验签
	// 参考: https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_3
	return true, nil
}

// VerifyAlipayCallback 模拟验证支付宝回调
func (v *MockPaymentCallbackVerifier) VerifyAlipayCallback(params map[string]string) (bool, error) {
	if isProductionEnv() {
		return false, fmt.Errorf("%w: alipay", ErrPaymentCallbackNotImplemented)
	}

	// TODO: 生产环境实现真实验签
	// 参考: https://opendocs.alipay.com/open/270/105902
	return true, nil
}

// NewPaymentCallbackVerifier 创建支付回调验签器
func NewPaymentCallbackVerifier(config *ThirdPartyConfig) PaymentCallbackVerifier {
	if isProductionEnv() {
		return &FailClosedPaymentCallbackVerifier{}
	}

	// TODO: 根据配置返回真实验签器
	return &MockPaymentCallbackVerifier{}
}
