package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// BusinessMetrics stores all business-related metrics
	BusinessMetrics *BusinessMetricsCollector
)

// BusinessMetricsCollector holds all business metrics
// BusinessMetricsCollector holds all business metrics
type BusinessMetricsCollector struct {
	// Order metrics
	OrdersCreatedTotal *prometheus.CounterVec
	OrdersCompletedTotal *prometheus.CounterVec
	OrdersCancelledTotal *prometheus.CounterVec
	OrdersRefundedTotal *prometheus.CounterVec
	OrderDurationHours *prometheus.HistogramVec
	
	// Payment metrics
	PaymentsCreatedTotal *prometheus.CounterVec
	PaymentsSucceededTotal *prometheus.CounterVec
	PaymentsFailedTotal *prometheus.CounterVec
	PaymentsRefundedTotal *prometheus.CounterVec
	PaymentAmountCents *prometheus.HistogramVec
	
	// User metrics
	UsersRegisteredTotal *prometheus.CounterVec
	UsersLoggedInTotal *prometheus.CounterVec
	UsersActive *prometheus.GaugeVec
	
	// Player metrics
	PlayersRegisteredTotal *prometheus.CounterVec
	PlayersVerifiedTotal *prometheus.CounterVec
	
	// Commission metrics
	CommissionTotalCents *prometheus.CounterVec
	CommissionRate *prometheus.GaugeVec
	
	// Error metrics
	ErrorsTotal *prometheus.CounterVec
}

// InitBusinessMetrics initializes business metrics
func InitBusinessMetrics(reg prometheus.Registerer) {
	BusinessMetrics = &BusinessMetricsCollector{
		// Order metrics
		OrdersCreatedTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gamelink_orders_created_total",
				Help: "Total number of orders created",
			},
			[]string{"status", "game_type"},
		),
		OrdersCompletedTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gamelink_orders_completed_total",
				Help: "Total number of orders completed",
			},
			[]string{"game_type", "player_tier"},
		),
		OrdersCancelledTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gamelink_orders_cancelled_total",
				Help: "Total number of orders cancelled",
			},
			[]string{"reason", "cancelled_by"},
		),
		OrdersRefundedTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gamelink_orders_refunded_total",
				Help: "Total number of orders refunded",
			},
			[]string{"refund_type"},
		),
		OrderDurationHours: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "gamelink_order_duration_hours",
				Help:    "Order duration in hours",
				Buckets: prometheus.LinearBuckets(0.5, 0.5, 20), // 0.5h to 10h
			},
			[]string{"game_type"},
		),
		
		// Payment metrics
		PaymentsCreatedTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gamelink_payments_created_total",
				Help: "Total number of payments created",
			},
			[]string{"method", "currency"},
		),
		PaymentsSucceededTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gamelink_payments_succeeded_total",
				Help: "Total number of successful payments",
			},
			[]string{"method", "currency"},
		),
		PaymentsFailedTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gamelink_payments_failed_total",
				Help: "Total number of failed payments",
			},
			[]string{"method", "failure_reason"},
		),
		PaymentsRefundedTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gamelink_payments_refunded_total",
				Help: "Total number of refunded payments",
			},
			[]string{"method"},
		),
		PaymentAmountCents: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "gamelink_payment_amount_cents",
				Help:    "Payment amount in cents",
				Buckets: prometheus.ExponentialBuckets(1000, 2, 15), // 10元 to ~16万
			},
			[]string{"method", "currency"},
		),
		
		// User metrics
		UsersRegisteredTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gamelink_users_registered_total",
				Help: "Total number of user registrations",
			},
			[]string{"role", "registration_method"},
		),
		UsersLoggedInTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gamelink_users_logged_in_total",
				Help: "Total number of user logins",
			},
			[]string{"role"},
		),
		UsersActive: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "gamelink_users_active",
				Help: "Number of currently active users",
			},
			[]string{"role", "status"},
		),
		
		// Player metrics
		PlayersRegisteredTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gamelink_players_registered_total",
				Help: "Total number of player registrations",
			},
			[]string{"verification_status"},
		),
		PlayersVerifiedTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gamelink_players_verified_total",
				Help: "Total number of players verified",
			},
			[]string{"game_category"},
		),
		
		// Commission metrics
		CommissionTotalCents: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gamelink_commission_total_cents",
				Help: "Total commission earned in cents",
			},
			[]string{"commission_type"},
		),
		CommissionRate: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "gamelink_commission_rate",
				Help: "Current commission rate",
			},
			[]string{"service_type"},
		),
		
		// Error metrics
		ErrorsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gamelink_errors_total",
				Help: "Total number of errors",
			},
			[]string{"error_type", "handler", "method"},
		),
	}
	
	// Register all metrics
	reg.MustRegister(
		BusinessMetrics.OrdersCreatedTotal,
		BusinessMetrics.OrdersCompletedTotal,
		BusinessMetrics.OrdersCancelledTotal,
		BusinessMetrics.OrdersRefundedTotal,
		BusinessMetrics.OrderDurationHours,
		BusinessMetrics.PaymentsCreatedTotal,
		BusinessMetrics.PaymentsSucceededTotal,
		BusinessMetrics.PaymentsFailedTotal,
		BusinessMetrics.PaymentsRefundedTotal,
		BusinessMetrics.PaymentAmountCents,
		BusinessMetrics.UsersRegisteredTotal,
		BusinessMetrics.UsersLoggedInTotal,
		BusinessMetrics.UsersActive,
		BusinessMetrics.PlayersRegisteredTotal,
		BusinessMetrics.PlayersVerifiedTotal,
		BusinessMetrics.CommissionTotalCents,
		BusinessMetrics.CommissionRate,
		BusinessMetrics.ErrorsTotal,
	)
}

// Order helpers
func RecordOrderCreated(status, gameType string) {
	if BusinessMetrics != nil {
		BusinessMetrics.OrdersCreatedTotal.WithLabelValues(status, gameType).Inc()
	}
}

func RecordOrderCompleted(gameType, playerTier string, durationHours float64) {
	if BusinessMetrics != nil {
		BusinessMetrics.OrdersCompletedTotal.WithLabelValues(gameType, playerTier).Inc()
		BusinessMetrics.OrderDurationHours.WithLabelValues(gameType).Observe(durationHours)
	}
}

func RecordOrderCancelled(reason, cancelledBy string) {
	if BusinessMetrics != nil {
		BusinessMetrics.OrdersCancelledTotal.WithLabelValues(reason, cancelledBy).Inc()
	}
}

func RecordOrderRefunded(refundType string) {
	if BusinessMetrics != nil {
		BusinessMetrics.OrdersRefundedTotal.WithLabelValues(refundType).Inc()
	}
}

// Payment helpers
func RecordPaymentCreated(method, currency string, amountCents int64) {
	if BusinessMetrics != nil {
		BusinessMetrics.PaymentsCreatedTotal.WithLabelValues(method, currency).Inc()
		BusinessMetrics.PaymentAmountCents.WithLabelValues(method, currency).Observe(float64(amountCents))
	}
}

func RecordPaymentSucceeded(method, currency string) {
	if BusinessMetrics != nil {
		BusinessMetrics.PaymentsSucceededTotal.WithLabelValues(method, currency).Inc()
	}
}

func RecordPaymentFailed(method, failureReason string) {
	if BusinessMetrics != nil {
		BusinessMetrics.PaymentsFailedTotal.WithLabelValues(method, failureReason).Inc()
	}
}

func RecordPaymentRefunded(method string) {
	if BusinessMetrics != nil {
		BusinessMetrics.PaymentsRefundedTotal.WithLabelValues(method).Inc()
	}
}

// User helpers
func RecordUserRegistered(role, registrationMethod string) {
	if BusinessMetrics != nil {
		BusinessMetrics.UsersRegisteredTotal.WithLabelValues(role, registrationMethod).Inc()
	}
}

func RecordUserLogin(role string) {
	if BusinessMetrics != nil {
		BusinessMetrics.UsersLoggedInTotal.WithLabelValues(role).Inc()
	}
}

func SetActiveUsers(role, status string, count int) {
	if BusinessMetrics != nil {
		BusinessMetrics.UsersActive.WithLabelValues(role, status).Set(float64(count))
	}
}

// Player helpers
func RecordPlayerRegistered(verificationStatus string) {
	if BusinessMetrics != nil {
		BusinessMetrics.PlayersRegisteredTotal.WithLabelValues(verificationStatus).Inc()
	}
}

func RecordPlayerVerified(gameCategory string) {
	if BusinessMetrics != nil {
		BusinessMetrics.PlayersVerifiedTotal.WithLabelValues(gameCategory).Inc()
	}
}

// Commission helpers
func RecordCommission(amountCents int64, commissionType string) {
	if BusinessMetrics != nil {
		BusinessMetrics.CommissionTotalCents.WithLabelValues(commissionType).Add(float64(amountCents))
	}
}

func SetCommissionRate(rate float64, serviceType string) {
	if BusinessMetrics != nil {
		BusinessMetrics.CommissionRate.WithLabelValues(serviceType).Set(rate)
	}
}

// Error helpers
func RecordError(errorType, handler, method string) {
	if BusinessMetrics != nil {
		BusinessMetrics.ErrorsTotal.WithLabelValues(errorType, handler, method).Inc()
	}
}