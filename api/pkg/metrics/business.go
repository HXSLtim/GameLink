package metrics

import "github.com/prometheus/client_golang/prometheus"

var BusinessMetrics *BusinessMetricsCollector

type BusinessMetricsCollector struct {
	OrdersCreatedTotal     *prometheus.CounterVec
	OrdersCompletedTotal   *prometheus.CounterVec
	OrdersCanceledTotal    *prometheus.CounterVec
	OrdersRefundedTotal    *prometheus.CounterVec
	OrderDurationHours     *prometheus.HistogramVec
	PaymentsCreatedTotal   *prometheus.CounterVec
	PaymentsSucceededTotal *prometheus.CounterVec
	PaymentsFailedTotal    *prometheus.CounterVec
	PaymentsRefundedTotal  *prometheus.CounterVec
	PaymentAmountCents     *prometheus.HistogramVec
	UsersRegisteredTotal   *prometheus.CounterVec
	UsersLoggedInTotal     *prometheus.CounterVec
	UsersActive            *prometheus.GaugeVec
	PlayersRegisteredTotal *prometheus.CounterVec
	PlayersVerifiedTotal   *prometheus.CounterVec
	CommissionTotalCents   *prometheus.CounterVec
	CommissionRate         *prometheus.GaugeVec
	ErrorsTotal            *prometheus.CounterVec
}

// InitBusinessMetrics initializes business metrics, safe to call multiple times.
func InitBusinessMetrics(reg prometheus.Registerer) {
	if reg == nil {
		reg = prometheus.DefaultRegisterer
	}

	if BusinessMetrics == nil {
		BusinessMetrics = &BusinessMetricsCollector{
			OrdersCreatedTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{Name: "gamelink_orders_created_total", Help: "Total number of orders created"},
				[]string{"status", "game_type"},
			),
			OrdersCompletedTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{Name: "gamelink_orders_completed_total", Help: "Total number of orders completed"},
				[]string{"game_type", "player_tier"},
			),
			OrdersCanceledTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{Name: "gamelink_orders_canceled_total", Help: "Total number of orders canceled"},
				[]string{"reason", "canceled_by"},
			),
			OrdersRefundedTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{Name: "gamelink_orders_refunded_total", Help: "Total number of orders refunded"},
				[]string{"refund_type"},
			),
			OrderDurationHours: prometheus.NewHistogramVec(
				prometheus.HistogramOpts{Name: "gamelink_order_duration_hours", Help: "Order duration in hours", Buckets: prometheus.LinearBuckets(0.5, 0.5, 20)},
				[]string{"game_type"},
			),
			PaymentsCreatedTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{Name: "gamelink_payments_created_total", Help: "Total number of payments created"},
				[]string{"method", "currency"},
			),
			PaymentsSucceededTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{Name: "gamelink_payments_succeeded_total", Help: "Total number of successful payments"},
				[]string{"method", "currency"},
			),
			PaymentsFailedTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{Name: "gamelink_payments_failed_total", Help: "Total number of failed payments"},
				[]string{"method", "failure_reason"},
			),
			PaymentsRefundedTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{Name: "gamelink_payments_refunded_total", Help: "Total number of refunded payments"},
				[]string{"method"},
			),
			PaymentAmountCents: prometheus.NewHistogramVec(
				prometheus.HistogramOpts{Name: "gamelink_payment_amount_cents", Help: "Payment amount in cents", Buckets: prometheus.ExponentialBuckets(1000, 2, 15)},
				[]string{"method", "currency"},
			),
			UsersRegisteredTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{Name: "gamelink_users_registered_total", Help: "Total number of user registrations"},
				[]string{"role", "registration_method"},
			),
			UsersLoggedInTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{Name: "gamelink_users_logged_in_total", Help: "Total number of user logins"},
				[]string{"role"},
			),
			UsersActive: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{Name: "gamelink_users_active", Help: "Number of currently active users"},
				[]string{"role", "status"},
			),
			PlayersRegisteredTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{Name: "gamelink_players_registered_total", Help: "Total number of player registrations"},
				[]string{"verification_status"},
			),
			PlayersVerifiedTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{Name: "gamelink_players_verified_total", Help: "Total number of players verified"},
				[]string{"game_category"},
			),
			CommissionTotalCents: prometheus.NewCounterVec(
				prometheus.CounterOpts{Name: "gamelink_commission_total_cents", Help: "Total commission earned in cents"},
				[]string{"commission_type"},
			),
			CommissionRate: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{Name: "gamelink_commission_rate", Help: "Current commission rate"},
				[]string{"service_type"},
			),
			ErrorsTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{Name: "gamelink_errors_total", Help: "Total number of errors"},
				[]string{"error_type", "handler", "method"},
			),
		}
	}

	register := func(c prometheus.Collector) {
		if c == nil {
			return
		}
		if err := reg.Register(c); err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				panic(err)
			}
		}
	}

	register(BusinessMetrics.OrdersCreatedTotal)
	register(BusinessMetrics.OrdersCompletedTotal)
	register(BusinessMetrics.OrdersCanceledTotal)
	register(BusinessMetrics.OrdersRefundedTotal)
	register(BusinessMetrics.OrderDurationHours)
	register(BusinessMetrics.PaymentsCreatedTotal)
	register(BusinessMetrics.PaymentsSucceededTotal)
	register(BusinessMetrics.PaymentsFailedTotal)
	register(BusinessMetrics.PaymentsRefundedTotal)
	register(BusinessMetrics.PaymentAmountCents)
	register(BusinessMetrics.UsersRegisteredTotal)
	register(BusinessMetrics.UsersLoggedInTotal)
	register(BusinessMetrics.UsersActive)
	register(BusinessMetrics.PlayersRegisteredTotal)
	register(BusinessMetrics.PlayersVerifiedTotal)
	register(BusinessMetrics.CommissionTotalCents)
	register(BusinessMetrics.CommissionRate)
	register(BusinessMetrics.ErrorsTotal)
}

// Helper functions (all are nil-safe)

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

func RecordOrderCanceled(reason, canceledBy string) {
	if BusinessMetrics != nil {
		BusinessMetrics.OrdersCanceledTotal.WithLabelValues(reason, canceledBy).Inc()
	}
}

func RecordOrderRefunded(refundType string) {
	if BusinessMetrics != nil {
		BusinessMetrics.OrdersRefundedTotal.WithLabelValues(refundType).Inc()
	}
}

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

func RecordUserRegistered(role, registrationMethod string) {
	if BusinessMetrics != nil {
		BusinessMetrics.UsersRegisteredTotal.WithLabelValues(role, registrationMethod).Inc()
	}
}

func RecordUserLoggedIn(role string) {
	if BusinessMetrics != nil {
		BusinessMetrics.UsersLoggedInTotal.WithLabelValues(role).Inc()
	}
}

// Alias for backward compatibility
func RecordUserLogin(role string) { RecordUserLoggedIn(role) }

func SetUsersActive(role, status string, value float64) {
	if BusinessMetrics != nil {
		BusinessMetrics.UsersActive.WithLabelValues(role, status).Set(value)
	}
}

// Alias for backward compatibility
func SetActiveUsers(role, status string, value float64) { SetUsersActive(role, status, value) }

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

func RecordCommission(amountCents int64, commissionType string) {
	if BusinessMetrics != nil {
		BusinessMetrics.CommissionTotalCents.WithLabelValues(commissionType).Add(float64(amountCents))
	}
}

func RecordCommissionRate(rate float64, serviceType string) {
	if BusinessMetrics != nil {
		BusinessMetrics.CommissionRate.WithLabelValues(serviceType).Set(rate)
	}
}

// Alias for backward compatibility
func SetCommissionRate(rate float64, serviceType string) { RecordCommissionRate(rate, serviceType) }

func RecordError(errorType, handler, method string) {
	if BusinessMetrics != nil {
		BusinessMetrics.ErrorsTotal.WithLabelValues(errorType, handler, method).Inc()
	}
}
