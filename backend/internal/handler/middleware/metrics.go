package middleware

import (
	"strconv"
	"time"

	"gamelink/internal/metrics"
	"github.com/gin-gonic/gin"
)

// MetricsMiddleware returns a gin middleware that records HTTP metrics
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		// Process request
		c.Next()

		// Record metrics after request
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method

		// Record HTTP request metrics
		if metrics.HTTPRequestsTotal != nil {
			metrics.HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
		}
		if metrics.HTTPRequestDuration != nil {
			metrics.HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)
		}

		// Record error metrics for 4xx and 5xx status codes
		if c.Writer.Status() >= 400 && metrics.BusinessMetrics != nil {
			errorType := "client_error"
			if c.Writer.Status() >= 500 {
				errorType = "server_error"
			}
			metrics.BusinessMetrics.ErrorsTotal.WithLabelValues(
				errorType,
				path,
				method,
			).Inc()
		}
	}
}

// RecordOrderMetrics records order-related business metrics
func RecordOrderMetrics(status, gameType string, durationHours float64) {
	if metrics.BusinessMetrics == nil {
		return
	}

	switch status {
	case "created":
		metrics.BusinessMetrics.OrdersCreatedTotal.WithLabelValues(status, gameType).Inc()
	case "completed":
		metrics.BusinessMetrics.OrdersCompletedTotal.WithLabelValues(gameType, "standard").Inc()
		if durationHours > 0 {
			metrics.BusinessMetrics.OrderDurationHours.WithLabelValues(gameType).Observe(durationHours)
		}
	case "cancelled":
		metrics.BusinessMetrics.OrdersCancelledTotal.WithLabelValues("user_cancelled", "user").Inc()
	case "refunded":
		metrics.BusinessMetrics.OrdersRefundedTotal.WithLabelValues("full").Inc()
	}
}

// RecordPaymentMetrics records payment-related business metrics
func RecordPaymentMetrics(method, currency string, amountCents int64, status string) {
	if metrics.BusinessMetrics == nil {
		return
	}

	switch status {
	case "created":
		metrics.BusinessMetrics.PaymentsCreatedTotal.WithLabelValues(method, currency).Inc()
		metrics.BusinessMetrics.PaymentAmountCents.WithLabelValues(method, currency).Observe(float64(amountCents))
	case "succeeded":
		metrics.BusinessMetrics.PaymentsSucceededTotal.WithLabelValues(method, currency).Inc()
	case "failed":
		metrics.BusinessMetrics.PaymentsFailedTotal.WithLabelValues(method, "unknown").Inc()
	case "refunded":
		metrics.BusinessMetrics.PaymentsRefundedTotal.WithLabelValues(method).Inc()
	}
}

// RecordUserMetrics records user-related business metrics
func RecordUserMetrics(role, method, action string) {
	if metrics.BusinessMetrics == nil {
		return
	}

	switch action {
	case "registered":
		metrics.BusinessMetrics.UsersRegisteredTotal.WithLabelValues(role, method).Inc()
	case "login":
		metrics.BusinessMetrics.UsersLoggedInTotal.WithLabelValues(role).Inc()
	}
}

// RecordPlayerMetrics records player-related business metrics
func RecordPlayerMetrics(verificationStatus, action string) {
	if metrics.BusinessMetrics == nil {
		return
	}

	switch action {
	case "registered":
		metrics.BusinessMetrics.PlayersRegisteredTotal.WithLabelValues(verificationStatus).Inc()
	case "verified":
		metrics.BusinessMetrics.PlayersVerifiedTotal.WithLabelValues("game_category").Inc()
	}
}

// RecordCommissionMetrics records commission-related business metrics
func RecordCommissionMetrics(amountCents int64, commissionType string) {
	if metrics.BusinessMetrics == nil {
		return
	}

	metrics.BusinessMetrics.CommissionTotalCents.WithLabelValues(commissionType).Add(float64(amountCents))
}

// RecordErrorMetrics records error-related business metrics
func RecordErrorMetrics(errorType, handler, method string) {
	if metrics.BusinessMetrics == nil {
		return
	}

	metrics.BusinessMetrics.ErrorsTotal.WithLabelValues(errorType, handler, method).Inc()
}
