package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/**
 * Test_InitBusinessMetrics_InitializesAllMetrics 测试初始化所有业务指标
 */
func Test_InitBusinessMetrics_InitializesAllMetrics(t *testing.T) {
	// Arrange
	reg := prometheus.NewRegistry()

	// Act
	InitBusinessMetrics(reg)

	// Assert
	assert.NotNil(t, BusinessMetrics, "BusinessMetrics应该被初始化")
	assert.NotNil(t, BusinessMetrics.OrdersCreatedTotal)
	assert.NotNil(t, BusinessMetrics.OrdersCompletedTotal)
	assert.NotNil(t, BusinessMetrics.OrdersCancelledTotal)
	assert.NotNil(t, BusinessMetrics.OrdersRefundedTotal)
	assert.NotNil(t, BusinessMetrics.OrderDurationHours)
	assert.NotNil(t, BusinessMetrics.PaymentsCreatedTotal)
	assert.NotNil(t, BusinessMetrics.PaymentsSucceededTotal)
	assert.NotNil(t, BusinessMetrics.PaymentsFailedTotal)
	assert.NotNil(t, BusinessMetrics.PaymentsRefundedTotal)
	assert.NotNil(t, BusinessMetrics.PaymentAmountCents)
	assert.NotNil(t, BusinessMetrics.UsersRegisteredTotal)
	assert.NotNil(t, BusinessMetrics.UsersLoggedInTotal)
	assert.NotNil(t, BusinessMetrics.UsersActive)
	assert.NotNil(t, BusinessMetrics.PlayersRegisteredTotal)
	assert.NotNil(t, BusinessMetrics.PlayersVerifiedTotal)
	assert.NotNil(t, BusinessMetrics.CommissionTotalCents)
	assert.NotNil(t, BusinessMetrics.CommissionRate)
	assert.NotNil(t, BusinessMetrics.ErrorsTotal)
}

/**
 * Test_InitBusinessMetrics_RegistersAllMetrics 测试所有指标都注册到registry
 */
func Test_InitBusinessMetrics_RegistersAllMetrics(t *testing.T) {
	// Arrange
	reg := prometheus.NewRegistry()

	// Act
	InitBusinessMetrics(reg)

	// 使用一些指标以便Gather能看到它们
	RecordOrderCreated("pending", "LOL")
	RecordPaymentCreated("alipay", "CNY", 10000)
	RecordUserRegistered("user", "email")

	// Assert
	metrics, err := reg.Gather()
	require.NoError(t, err)

	metricNames := make(map[string]bool)
	for _, m := range metrics {
		metricNames[m.GetName()] = true
	}

	// 验证关键指标已注册
	assert.True(t, metricNames["gamelink_orders_created_total"])
	assert.True(t, metricNames["gamelink_payments_created_total"])
	assert.True(t, metricNames["gamelink_payment_amount_cents"])
	assert.True(t, metricNames["gamelink_users_registered_total"])
}

/**
 * Test_RecordOrderCreated_RecordsMetric 测试记录订单创建
 */
func Test_RecordOrderCreated_RecordsMetric(t *testing.T) {
	// Arrange
	reg := prometheus.NewRegistry()
	InitBusinessMetrics(reg)

	// Act
	RecordOrderCreated("pending", "LOL")
	RecordOrderCreated("pending", "DOTA2")
	RecordOrderCreated("confirmed", "LOL")

	// Assert
	metrics, err := reg.Gather()
	require.NoError(t, err)

	for _, m := range metrics {
		if m.GetName() == "gamelink_orders_created_total" {
			assert.Equal(t, 3, len(m.GetMetric()), "应该有3个不同的标签组合")
		}
	}
}

/**
 * Test_RecordOrderCreated_WithNilBusinessMetrics 测试nil BusinessMetrics不会panic
 */
func Test_RecordOrderCreated_WithNilBusinessMetrics(t *testing.T) {
	// Arrange
	oldMetrics := BusinessMetrics
	BusinessMetrics = nil
	defer func() { BusinessMetrics = oldMetrics }()

	// Act & Assert - 不应该panic
	assert.NotPanics(t, func() {
		RecordOrderCreated("pending", "LOL")
	})
}

/**
 * Test_RecordOrderCompleted_RecordsMetric 测试记录订单完成
 */
func Test_RecordOrderCompleted_RecordsMetric(t *testing.T) {
	// Arrange
	reg := prometheus.NewRegistry()
	InitBusinessMetrics(reg)

	// Act
	RecordOrderCompleted("LOL", "gold", 2.5)
	RecordOrderCompleted("DOTA2", "silver", 3.0)

	// Assert
	metrics, err := reg.Gather()
	require.NoError(t, err)

	var foundCounter, foundHistogram bool
	for _, m := range metrics {
		if m.GetName() == "gamelink_orders_completed_total" {
			foundCounter = true
			assert.Equal(t, 2, len(m.GetMetric()))
		}
		if m.GetName() == "gamelink_order_duration_hours" {
			foundHistogram = true
		}
	}
	assert.True(t, foundCounter, "应该找到订单完成计数器")
	assert.True(t, foundHistogram, "应该找到订单时长直方图")
}

/**
 * Test_RecordOrderCancelled_RecordsMetric 测试记录订单取消
 */
func Test_RecordOrderCancelled_RecordsMetric(t *testing.T) {
	// Arrange
	reg := prometheus.NewRegistry()
	InitBusinessMetrics(reg)

	// Act
	RecordOrderCancelled("user_request", "user")
	RecordOrderCancelled("timeout", "system")
	RecordOrderCancelled("player_unavailable", "player")

	// Assert
	metrics, err := reg.Gather()
	require.NoError(t, err)

	for _, m := range metrics {
		if m.GetName() == "gamelink_orders_cancelled_total" {
			assert.Equal(t, 3, len(m.GetMetric()))
		}
	}
}

/**
 * Test_RecordOrderRefunded_RecordsMetric 测试记录订单退款
 */
func Test_RecordOrderRefunded_RecordsMetric(t *testing.T) {
	// Arrange
	reg := prometheus.NewRegistry()
	InitBusinessMetrics(reg)

	// Act
	RecordOrderRefunded("full")
	RecordOrderRefunded("partial")

	// Assert
	metrics, err := reg.Gather()
	require.NoError(t, err)

	for _, m := range metrics {
		if m.GetName() == "gamelink_orders_refunded_total" {
			assert.Equal(t, 2, len(m.GetMetric()))
		}
	}
}

/**
 * Test_RecordPaymentCreated_RecordsMetricAndAmount 测试记录支付创建和金额
 */
func Test_RecordPaymentCreated_RecordsMetricAndAmount(t *testing.T) {
	// Arrange
	reg := prometheus.NewRegistry()
	InitBusinessMetrics(reg)

	// Act
	RecordPaymentCreated("alipay", "CNY", 10000)
	RecordPaymentCreated("wechat", "CNY", 20000)
	RecordPaymentCreated("stripe", "USD", 5000)

	// Assert
	metrics, err := reg.Gather()
	require.NoError(t, err)

	var foundCounter, foundHistogram bool
	for _, m := range metrics {
		if m.GetName() == "gamelink_payments_created_total" {
			foundCounter = true
			assert.Equal(t, 3, len(m.GetMetric()))
		}
		if m.GetName() == "gamelink_payment_amount_cents" {
			foundHistogram = true
			assert.Equal(t, 3, len(m.GetMetric()))
		}
	}
	assert.True(t, foundCounter)
	assert.True(t, foundHistogram)
}

/**
 * Test_RecordPaymentSucceeded_RecordsMetric 测试记录支付成功
 */
func Test_RecordPaymentSucceeded_RecordsMetric(t *testing.T) {
	// Arrange
	reg := prometheus.NewRegistry()
	InitBusinessMetrics(reg)

	// Act
	RecordPaymentSucceeded("alipay", "CNY")
	RecordPaymentSucceeded("wechat", "CNY")

	// Assert
	metrics, err := reg.Gather()
	require.NoError(t, err)

	for _, m := range metrics {
		if m.GetName() == "gamelink_payments_succeeded_total" {
			assert.Equal(t, 2, len(m.GetMetric()))
		}
	}
}

/**
 * Test_RecordPaymentFailed_RecordsMetric 测试记录支付失败
 */
func Test_RecordPaymentFailed_RecordsMetric(t *testing.T) {
	// Arrange
	reg := prometheus.NewRegistry()
	InitBusinessMetrics(reg)

	// Act
	RecordPaymentFailed("alipay", "insufficient_balance")
	RecordPaymentFailed("stripe", "card_declined")

	// Assert
	metrics, err := reg.Gather()
	require.NoError(t, err)

	for _, m := range metrics {
		if m.GetName() == "gamelink_payments_failed_total" {
			assert.Equal(t, 2, len(m.GetMetric()))
		}
	}
}

/**
 * Test_RecordPaymentRefunded_RecordsMetric 测试记录支付退款
 */
func Test_RecordPaymentRefunded_RecordsMetric(t *testing.T) {
	// Arrange
	reg := prometheus.NewRegistry()
	InitBusinessMetrics(reg)

	// Act
	RecordPaymentRefunded("alipay")
	RecordPaymentRefunded("wechat")

	// Assert
	metrics, err := reg.Gather()
	require.NoError(t, err)

	for _, m := range metrics {
		if m.GetName() == "gamelink_payments_refunded_total" {
			assert.Equal(t, 2, len(m.GetMetric()))
		}
	}
}

/**
 * Test_RecordUserRegistered_RecordsMetric 测试记录用户注册
 */
func Test_RecordUserRegistered_RecordsMetric(t *testing.T) {
	// Arrange
	reg := prometheus.NewRegistry()
	InitBusinessMetrics(reg)

	// Act
	RecordUserRegistered("user", "email")
	RecordUserRegistered("user", "phone")
	RecordUserRegistered("player", "email")

	// Assert
	metrics, err := reg.Gather()
	require.NoError(t, err)

	for _, m := range metrics {
		if m.GetName() == "gamelink_users_registered_total" {
			assert.Equal(t, 3, len(m.GetMetric()))
		}
	}
}

/**
 * Test_RecordUserLogin_RecordsMetric 测试记录用户登录
 */
func Test_RecordUserLogin_RecordsMetric(t *testing.T) {
	// Arrange
	reg := prometheus.NewRegistry()
	InitBusinessMetrics(reg)

	// Act
	RecordUserLogin("user")
	RecordUserLogin("player")
	RecordUserLogin("admin")

	// Assert
	metrics, err := reg.Gather()
	require.NoError(t, err)

	for _, m := range metrics {
		if m.GetName() == "gamelink_users_logged_in_total" {
			assert.Equal(t, 3, len(m.GetMetric()))
		}
	}
}

/**
 * Test_SetActiveUsers_SetsGauge 测试设置活跃用户数（Gauge）
 */
func Test_SetActiveUsers_SetsGauge(t *testing.T) {
	// Arrange
	reg := prometheus.NewRegistry()
	InitBusinessMetrics(reg)

	// Act
	SetActiveUsers("user", "online", 100)
	SetActiveUsers("player", "online", 50)
	SetActiveUsers("user", "idle", 20)

	// Assert
	metrics, err := reg.Gather()
	require.NoError(t, err)

	for _, m := range metrics {
		if m.GetName() == "gamelink_users_active" {
			assert.Equal(t, 3, len(m.GetMetric()))
			// 验证gauge值被正确设置
			for _, metric := range m.GetMetric() {
				labels := metric.GetLabel()
				for _, label := range labels {
					if label.GetName() == "role" && label.GetValue() == "user" {
						// 可能是online=100或idle=20
						value := metric.GetGauge().GetValue()
						assert.True(t, value == 100 || value == 20)
					}
				}
			}
		}
	}
}

/**
 * Test_RecordPlayerRegistered_RecordsMetric 测试记录打手注册
 */
func Test_RecordPlayerRegistered_RecordsMetric(t *testing.T) {
	// Arrange
	reg := prometheus.NewRegistry()
	InitBusinessMetrics(reg)

	// Act
	RecordPlayerRegistered("pending")
	RecordPlayerRegistered("verified")

	// Assert
	metrics, err := reg.Gather()
	require.NoError(t, err)

	for _, m := range metrics {
		if m.GetName() == "gamelink_players_registered_total" {
			assert.Equal(t, 2, len(m.GetMetric()))
		}
	}
}

/**
 * Test_RecordPlayerVerified_RecordsMetric 测试记录打手认证
 */
func Test_RecordPlayerVerified_RecordsMetric(t *testing.T) {
	// Arrange
	reg := prometheus.NewRegistry()
	InitBusinessMetrics(reg)

	// Act
	RecordPlayerVerified("LOL")
	RecordPlayerVerified("DOTA2")
	RecordPlayerVerified("CSGO")

	// Assert
	metrics, err := reg.Gather()
	require.NoError(t, err)

	for _, m := range metrics {
		if m.GetName() == "gamelink_players_verified_total" {
			assert.Equal(t, 3, len(m.GetMetric()))
		}
	}
}

/**
 * Test_RecordCommission_RecordsMetric 测试记录佣金
 */
func Test_RecordCommission_RecordsMetric(t *testing.T) {
	// Arrange
	reg := prometheus.NewRegistry()
	InitBusinessMetrics(reg)

	// Act
	RecordCommission(1000, "platform")
	RecordCommission(500, "referral")
	RecordCommission(2000, "platform")

	// Assert
	metrics, err := reg.Gather()
	require.NoError(t, err)

	for _, m := range metrics {
		if m.GetName() == "gamelink_commission_total_cents" {
			assert.Equal(t, 2, len(m.GetMetric()), "应该有2种佣金类型")
			// 验证累加正确
			for _, metric := range m.GetMetric() {
				labels := metric.GetLabel()
				for _, label := range labels {
					if label.GetName() == "commission_type" && label.GetValue() == "platform" {
						// platform应该是1000+2000=3000
						assert.Equal(t, float64(3000), metric.GetCounter().GetValue())
					}
				}
			}
		}
	}
}

/**
 * Test_SetCommissionRate_SetsGauge 测试设置佣金率
 */
func Test_SetCommissionRate_SetsGauge(t *testing.T) {
	// Arrange
	reg := prometheus.NewRegistry()
	InitBusinessMetrics(reg)

	// Act
	SetCommissionRate(0.15, "basic")
	SetCommissionRate(0.20, "premium")

	// Assert
	metrics, err := reg.Gather()
	require.NoError(t, err)

	for _, m := range metrics {
		if m.GetName() == "gamelink_commission_rate" {
			assert.Equal(t, 2, len(m.GetMetric()))
		}
	}
}

/**
 * Test_RecordError_RecordsMetric 测试记录错误
 */
func Test_RecordError_RecordsMetric(t *testing.T) {
	// Arrange
	reg := prometheus.NewRegistry()
	InitBusinessMetrics(reg)

	// Act
	RecordError("validation", "UserHandler", "POST")
	RecordError("database", "OrderHandler", "GET")
	RecordError("validation", "UserHandler", "PUT")

	// Assert
	metrics, err := reg.Gather()
	require.NoError(t, err)

	for _, m := range metrics {
		if m.GetName() == "gamelink_errors_total" {
			assert.Equal(t, 3, len(m.GetMetric()))
		}
	}
}

/**
 * Test_AllHelperFunctions_WithNilBusinessMetrics 测试所有辅助函数在nil时不panic
 */
func Test_AllHelperFunctions_WithNilBusinessMetrics(t *testing.T) {
	// Arrange
	oldMetrics := BusinessMetrics
	BusinessMetrics = nil
	defer func() { BusinessMetrics = oldMetrics }()

	// Act & Assert - 所有函数都不应该panic
	assert.NotPanics(t, func() {
		RecordOrderCreated("pending", "LOL")
		RecordOrderCompleted("LOL", "gold", 2.5)
		RecordOrderCancelled("timeout", "system")
		RecordOrderRefunded("full")
		RecordPaymentCreated("alipay", "CNY", 10000)
		RecordPaymentSucceeded("alipay", "CNY")
		RecordPaymentFailed("alipay", "error")
		RecordPaymentRefunded("alipay")
		RecordUserRegistered("user", "email")
		RecordUserLogin("user")
		SetActiveUsers("user", "online", 100)
		RecordPlayerRegistered("pending")
		RecordPlayerVerified("LOL")
		RecordCommission(1000, "platform")
		SetCommissionRate(0.15, "basic")
		RecordError("validation", "Handler", "POST")
	})
}

/**
 * Test_BusinessMetrics_CompleteScenario 测试完整的业务场景
 */
func Test_BusinessMetrics_CompleteScenario(t *testing.T) {
	// Arrange
	reg := prometheus.NewRegistry()
	InitBusinessMetrics(reg)

	// Act - 模拟一个完整的订单流程
	// 1. 用户注册和登录
	RecordUserRegistered("user", "email")
	RecordUserLogin("user")
	SetActiveUsers("user", "online", 1)

	// 2. 打手注册和认证
	RecordPlayerRegistered("pending")
	RecordPlayerVerified("LOL")

	// 3. 创建订单
	RecordOrderCreated("pending", "LOL")

	// 4. 创建支付
	RecordPaymentCreated("alipay", "CNY", 10000)
	RecordPaymentSucceeded("alipay", "CNY")

	// 5. 完成订单
	RecordOrderCompleted("LOL", "gold", 2.5)
	RecordCommission(1000, "platform")

	// 6. 记录一些错误
	RecordError("validation", "OrderHandler", "POST")

	// Assert - 验证所有指标都被正确记录
	metrics, err := reg.Gather()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(metrics), 10, "应该有至少10个不同的指标")

	// 验证关键指标存在
	metricNames := make(map[string]bool)
	for _, m := range metrics {
		metricNames[m.GetName()] = true
	}

	assert.True(t, metricNames["gamelink_users_registered_total"])
	assert.True(t, metricNames["gamelink_orders_created_total"])
	assert.True(t, metricNames["gamelink_payments_created_total"])
	assert.True(t, metricNames["gamelink_orders_completed_total"])
	assert.True(t, metricNames["gamelink_commission_total_cents"])
}
