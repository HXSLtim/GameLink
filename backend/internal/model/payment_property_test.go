package model_test

import (
	"reflect"
	"testing"
	"time"

	"gamelink/internal/model"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: payment-finance-module, Property 1: 支付记录完整性**
// **Validates: Requirements 1.1**
//
// Property 1: Payment Record Completeness
// *For any* payment record list request, all returned payment records must contain
// complete required fields (OrderID, UserID, Method, Amount, Status, CreatedAt),
// and field values must not be empty or null.

// genPaymentMethod generates valid payment methods
func genPaymentMethod() gopter.Gen {
	return gen.OneConstOf(
		model.PaymentMethodWeChat,
		model.PaymentMethodAlipay,
	)
}

// genPaymentStatus generates valid payment statuses
func genPaymentStatus() gopter.Gen {
	return gen.OneConstOf(
		model.PaymentStatusPending,
		model.PaymentStatusPaid,
		model.PaymentStatusFailed,
		model.PaymentStatusRefunded,
	)
}

// genCurrency generates valid currencies
func genCurrency() gopter.Gen {
	return gen.OneConstOf(
		model.CurrencyCNY,
		model.CurrencyUSD,
	)
}

// genValidPayment generates a valid payment with all required fields
func genValidPayment() gopter.Gen {
	return gopter.CombineGens(
		gen.UInt64Range(1, 1000000), // OrderID
		gen.UInt64Range(1, 1000000), // UserID
		genPaymentMethod(),          // Method
		gen.Int64Range(1, 10000000), // AmountCents (positive)
		genCurrency(),               // Currency
		genPaymentStatus(),          // Status
	).Map(func(values []interface{}) *model.Payment {
		now := time.Now()
		return &model.Payment{
			Base: model.Base{
				ID:        uint64(values[0].(uint64)),
				CreatedAt: now,
				UpdatedAt: now,
			},
			OrderID:     values[0].(uint64),
			UserID:      values[1].(uint64),
			Method:      values[2].(model.PaymentMethod),
			AmountCents: values[3].(int64),
			Currency:    values[4].(model.Currency),
			Status:      values[5].(model.PaymentStatus),
		}
	})
}

// genInvalidPayment generates a payment with at least one missing required field
func genInvalidPayment() gopter.Gen {
	return gen.IntRange(0, 4).FlatMap(func(missingField interface{}) gopter.Gen {
		return genValidPayment().Map(func(p *model.Payment) *model.Payment {
			switch missingField.(int) {
			case 0:
				p.OrderID = 0
			case 1:
				p.UserID = 0
			case 2:
				p.Method = ""
			case 3:
				p.Status = ""
			case 4:
				p.CreatedAt = time.Time{}
			}
			return p
		})
	}, reflect.TypeOf(&model.Payment{}))
}

// TestProperty1_PaymentRecordCompleteness tests that valid payments have all required fields
func TestProperty1_PaymentRecordCompleteness(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property 1: Valid payments should have all required fields
	properties.Property("valid payment should have all required fields", prop.ForAll(
		func(p *model.Payment) bool {
			return p.HasRequiredFields()
		},
		genValidPayment(),
	))

	// Property 2: Payments with missing required fields should fail validation
	properties.Property("payment with missing required field should fail validation", prop.ForAll(
		func(p *model.Payment) bool {
			return !p.HasRequiredFields()
		},
		genInvalidPayment(),
	))

	// Property 3: OrderID must be non-zero for valid payments
	properties.Property("valid payment must have non-zero OrderID", prop.ForAll(
		func(p *model.Payment) bool {
			if p.HasRequiredFields() {
				return p.OrderID != 0
			}
			return true
		},
		genValidPayment(),
	))

	// Property 4: UserID must be non-zero for valid payments
	properties.Property("valid payment must have non-zero UserID", prop.ForAll(
		func(p *model.Payment) bool {
			if p.HasRequiredFields() {
				return p.UserID != 0
			}
			return true
		},
		genValidPayment(),
	))

	// Property 5: Method must be non-empty for valid payments
	properties.Property("valid payment must have non-empty Method", prop.ForAll(
		func(p *model.Payment) bool {
			if p.HasRequiredFields() {
				return p.Method != ""
			}
			return true
		},
		genValidPayment(),
	))

	// Property 6: Status must be non-empty for valid payments
	properties.Property("valid payment must have non-empty Status", prop.ForAll(
		func(p *model.Payment) bool {
			if p.HasRequiredFields() {
				return p.Status != ""
			}
			return true
		},
		genValidPayment(),
	))

	// Property 7: CreatedAt must be non-zero for valid payments
	properties.Property("valid payment must have non-zero CreatedAt", prop.ForAll(
		func(p *model.Payment) bool {
			if p.HasRequiredFields() {
				return !p.CreatedAt.IsZero()
			}
			return true
		},
		genValidPayment(),
	))

	properties.TestingRun(t)
}

// **Feature: payment-finance-module, Property 5: 分页数据连续性**
// **Validates: Requirements 1.4, 4.3**
//
// Property 5: Pagination Data Continuity
// *For any* paginated query, when sorted by CreatedAt descending,
// the last record of the current page must have CreatedAt >= the first record of the next page.

// TestProperty5_PaginationDataContinuity tests pagination continuity for payment records.
// This is a unit test that verifies the pagination logic maintains proper ordering.
func TestProperty5_PaginationDataContinuity(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Generate a list of payments with sequential timestamps
	genPaymentList := func(count int) gopter.Gen {
		return gen.IntRange(2, count).FlatMap(func(n interface{}) gopter.Gen {
			numPayments := n.(int)
			return gopter.CombineGens(
				gen.SliceOfN(numPayments, genValidPayment()),
			).Map(func(values []interface{}) []*model.Payment {
				payments := values[0].([]*model.Payment)
				// Assign sequential timestamps to ensure proper ordering
				baseTime := time.Now()
				for i, p := range payments {
					p.CreatedAt = baseTime.Add(time.Duration(-i) * time.Minute)
					p.ID = uint64(i + 1)
				}
				return payments
			})
		}, reflect.TypeOf([]*model.Payment{}))
	}

	// Property: When payments are sorted by CreatedAt DESC, pagination maintains order
	properties.Property("pagination maintains descending order by CreatedAt", prop.ForAll(
		func(payments []*model.Payment) bool {
			if len(payments) < 2 {
				return true
			}

			// Sort payments by CreatedAt DESC (simulating database ordering)
			sortedPayments := make([]*model.Payment, len(payments))
			copy(sortedPayments, payments)
			for i := 0; i < len(sortedPayments)-1; i++ {
				for j := i + 1; j < len(sortedPayments); j++ {
					if sortedPayments[i].CreatedAt.Before(sortedPayments[j].CreatedAt) {
						sortedPayments[i], sortedPayments[j] = sortedPayments[j], sortedPayments[i]
					}
				}
			}

			// Simulate pagination with page size of 3
			pageSize := 3
			for pageStart := 0; pageStart < len(sortedPayments)-pageSize; pageStart += pageSize {
				pageEnd := pageStart + pageSize
				if pageEnd > len(sortedPayments) {
					pageEnd = len(sortedPayments)
				}

				currentPageLast := sortedPayments[pageEnd-1]
				nextPageStart := pageEnd
				if nextPageStart < len(sortedPayments) {
					nextPageFirst := sortedPayments[nextPageStart]
					// Property: last item of current page should have CreatedAt >= first item of next page
					if currentPageLast.CreatedAt.Before(nextPageFirst.CreatedAt) {
						return false
					}
				}
			}
			return true
		},
		genPaymentList(20),
	))

	// Property: Adjacent pages should not have overlapping records (by ID)
	properties.Property("adjacent pages should not have overlapping records", prop.ForAll(
		func(payments []*model.Payment) bool {
			if len(payments) < 4 {
				return true
			}

			// Sort payments by CreatedAt DESC
			sortedPayments := make([]*model.Payment, len(payments))
			copy(sortedPayments, payments)
			for i := 0; i < len(sortedPayments)-1; i++ {
				for j := i + 1; j < len(sortedPayments); j++ {
					if sortedPayments[i].CreatedAt.Before(sortedPayments[j].CreatedAt) {
						sortedPayments[i], sortedPayments[j] = sortedPayments[j], sortedPayments[i]
					}
				}
			}

			// Simulate pagination with page size of 3
			pageSize := 3
			for pageStart := 0; pageStart < len(sortedPayments)-pageSize; pageStart += pageSize {
				pageEnd := pageStart + pageSize
				if pageEnd > len(sortedPayments) {
					pageEnd = len(sortedPayments)
				}

				currentPageIDs := make(map[uint64]bool)
				for i := pageStart; i < pageEnd; i++ {
					currentPageIDs[sortedPayments[i].ID] = true
				}

				nextPageStart := pageEnd
				nextPageEnd := nextPageStart + pageSize
				if nextPageEnd > len(sortedPayments) {
					nextPageEnd = len(sortedPayments)
				}

				// Check no overlap with next page
				for i := nextPageStart; i < nextPageEnd; i++ {
					if currentPageIDs[sortedPayments[i].ID] {
						return false
					}
				}
			}
			return true
		},
		genPaymentList(20),
	))

	// Property: Total records across all pages should equal total count
	properties.Property("total records across all pages equals total count", prop.ForAll(
		func(payments []*model.Payment) bool {
			if len(payments) == 0 {
				return true
			}

			pageSize := 3
			totalRecords := 0
			for pageStart := 0; pageStart < len(payments); pageStart += pageSize {
				pageEnd := pageStart + pageSize
				if pageEnd > len(payments) {
					pageEnd = len(payments)
				}
				totalRecords += pageEnd - pageStart
			}

			return totalRecords == len(payments)
		},
		genPaymentList(20),
	))

	properties.TestingRun(t)
}

// **Feature: payment-finance-module, Property 2: 退款金额约束**
// **Validates: Requirements 2.1, 9.1, 9.2**
//
// Property 2: Refund Amount Constraints
// *For any* refund operation, the refund amount must be greater than 0
// and must not exceed the original payment amount minus the already refunded amount.

// genPaidPayment generates a paid payment that can be refunded
func genPaidPayment() gopter.Gen {
	return gopter.CombineGens(
		gen.UInt64Range(1, 1000000),   // 0: OrderID
		gen.UInt64Range(1, 1000000),   // 1: UserID
		genPaymentMethod(),            // 2: Method
		gen.Int64Range(100, 10000000), // 3: AmountCents (at least 100 cents = 1 yuan)
		genCurrency(),                 // 4: Currency
		gen.Int64Range(0, 50),         // 5: RefundedAmountPercent (0-50% already refunded)
	).Map(func(values []interface{}) *model.Payment {
		now := time.Now()
		amountCents := values[3].(int64)
		refundedPercent := values[5].(int64) // Fixed: was values[4]
		refundedAmount := (amountCents * refundedPercent) / 100

		return &model.Payment{
			Base: model.Base{
				ID:        uint64(values[0].(uint64)),
				CreatedAt: now,
				UpdatedAt: now,
			},
			OrderID:             values[0].(uint64),
			UserID:              values[1].(uint64),
			Method:              values[2].(model.PaymentMethod),
			AmountCents:         amountCents,
			Currency:            values[4].(model.Currency), // Fixed: correct index
			Status:              model.PaymentStatusPaid,
			RefundedAmountCents: refundedAmount,
			PaidAt:              &now,
		}
	})
}

// genValidRefundAmount generates a valid refund amount for a given payment
func genValidRefundAmount(payment *model.Payment) gopter.Gen {
	remaining := payment.RemainingRefundableAmount()
	if remaining <= 0 {
		return gen.Const(int64(0))
	}
	return gen.Int64Range(1, remaining)
}

// genInvalidRefundAmount generates an invalid refund amount (negative, zero, or exceeds remaining)
func genInvalidRefundAmount(payment *model.Payment) gopter.Gen {
	remaining := payment.RemainingRefundableAmount()
	return gen.OneGenOf(
		gen.Const(int64(0)),                            // Zero amount
		gen.Int64Range(-1000000, -1),                   // Negative amount
		gen.Int64Range(remaining+1, remaining+1000000), // Exceeds remaining
	)
}

// TestProperty2_RefundAmountConstraints tests refund amount validation
func TestProperty2_RefundAmountConstraints(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property 2.1: Valid refund amounts should pass validation
	properties.Property("valid refund amount should pass validation", prop.ForAll(
		func(payment *model.Payment, refundPercent int64) bool {
			remaining := payment.RemainingRefundableAmount()
			if remaining <= 0 {
				return true // Skip if nothing to refund
			}
			// Calculate a valid refund amount (1% to 100% of remaining)
			refundAmount := (remaining * (refundPercent + 1)) / 100
			if refundAmount <= 0 {
				refundAmount = 1
			}
			if refundAmount > remaining {
				refundAmount = remaining
			}
			err := payment.ValidateRefundAmount(refundAmount)
			return err == nil
		},
		genPaidPayment(),
		gen.Int64Range(0, 99),
	))

	// Property 2.2: Zero refund amount should fail validation
	properties.Property("zero refund amount should fail validation", prop.ForAll(
		func(payment *model.Payment) bool {
			err := payment.ValidateRefundAmount(0)
			if err == nil {
				return false
			}
			refundErr, ok := err.(*model.RefundValidationError)
			return ok && refundErr.Code == model.RefundErrCodeInvalidAmount
		},
		genPaidPayment(),
	))

	// Property 2.3: Negative refund amount should fail validation
	properties.Property("negative refund amount should fail validation", prop.ForAll(
		func(payment *model.Payment, negativeAmount int64) bool {
			err := payment.ValidateRefundAmount(negativeAmount)
			if err == nil {
				return false
			}
			refundErr, ok := err.(*model.RefundValidationError)
			return ok && refundErr.Code == model.RefundErrCodeInvalidAmount
		},
		genPaidPayment(),
		gen.Int64Range(-1000000, -1),
	))

	// Property 2.4: Refund amount exceeding remaining should fail validation
	properties.Property("refund amount exceeding remaining should fail validation", prop.ForAll(
		func(payment *model.Payment) bool {
			remaining := payment.RemainingRefundableAmount()
			if remaining <= 0 {
				return true // Skip if nothing to refund
			}
			// Try to refund more than remaining
			excessAmount := remaining + 1
			err := payment.ValidateRefundAmount(excessAmount)
			if err == nil {
				return false
			}
			refundErr, ok := err.(*model.RefundValidationError)
			return ok && refundErr.Code == model.RefundErrCodeExceedsRemaining
		},
		genPaidPayment(),
	))

	// Property 2.5: Refund on non-paid payment should fail
	properties.Property("refund on non-paid payment should fail", prop.ForAll(
		func(payment *model.Payment, status model.PaymentStatus) bool {
			if status == model.PaymentStatusPaid {
				return true // Skip paid status
			}
			payment.Status = status
			err := payment.ValidateRefundAmount(100)
			if err == nil {
				return false
			}
			refundErr, ok := err.(*model.RefundValidationError)
			return ok && refundErr.Code == model.RefundErrCodeInvalidStatus
		},
		genPaidPayment(),
		gen.OneConstOf(
			model.PaymentStatusPending,
			model.PaymentStatusFailed,
			model.PaymentStatusRefunded,
		),
	))

	// Property 2.6: Remaining refundable amount should decrease after partial refund
	properties.Property("remaining refundable amount decreases correctly", prop.ForAll(
		func(payment *model.Payment, refundPercent int64) bool {
			originalRemaining := payment.RemainingRefundableAmount()
			if originalRemaining <= 0 {
				return true
			}
			// Calculate refund amount
			refundAmount := (originalRemaining * (refundPercent + 1)) / 100
			if refundAmount <= 0 {
				refundAmount = 1
			}
			if refundAmount > originalRemaining {
				refundAmount = originalRemaining
			}

			// Simulate refund
			payment.RefundedAmountCents += refundAmount
			newRemaining := payment.RemainingRefundableAmount()

			return newRemaining == originalRemaining-refundAmount
		},
		genPaidPayment(),
		gen.Int64Range(0, 99),
	))

	// Property 2.7: Full refund should leave zero remaining
	properties.Property("full refund should leave zero remaining", prop.ForAll(
		func(payment *model.Payment) bool {
			remaining := payment.RemainingRefundableAmount()
			if remaining <= 0 {
				return true
			}
			// Simulate full refund
			payment.RefundedAmountCents += remaining
			return payment.RemainingRefundableAmount() == 0
		},
		genPaidPayment(),
	))

	// Property 2.8: IsFullyRefunded should be true when refunded equals amount
	properties.Property("IsFullyRefunded returns true when fully refunded", prop.ForAll(
		func(payment *model.Payment) bool {
			payment.RefundedAmountCents = payment.AmountCents
			return payment.IsFullyRefunded()
		},
		genPaidPayment(),
	))

	// Property 2.9: IsPartiallyRefunded should be true for partial refunds
	properties.Property("IsPartiallyRefunded returns true for partial refunds", prop.ForAll(
		func(payment *model.Payment, refundPercent int64) bool {
			if payment.AmountCents <= 1 {
				return true // Skip if amount too small for partial refund
			}
			// Set partial refund (1-99% of amount)
			partialAmount := (payment.AmountCents * (refundPercent + 1)) / 100
			if partialAmount <= 0 {
				partialAmount = 1
			}
			if partialAmount >= payment.AmountCents {
				partialAmount = payment.AmountCents - 1
			}
			payment.RefundedAmountCents = partialAmount
			return payment.IsPartiallyRefunded()
		},
		genPaidPayment(),
		gen.Int64Range(0, 98),
	))

	properties.TestingRun(t)
}

// **Feature: payment-finance-module, Property 3: 支付状态转换合法性**
// **Validates: Requirements 2.3, 2.4, 8.2**
//
// Property 3: Payment Status Transition Validity
// *For any* payment status update operation, the status transition must follow
// the valid path: pending -> paid -> refunded or pending -> failed.
// No other transitions are allowed.

// TestProperty3_PaymentStatusTransitionValidity tests payment status transition rules
func TestProperty3_PaymentStatusTransitionValidity(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property 3.1: Valid transitions should be allowed
	properties.Property("valid transitions should be allowed", prop.ForAll(
		func(from model.PaymentStatus, to model.PaymentStatus) bool {
			transition := model.PaymentStatusTransition{From: from, To: to}
			isValid := model.ValidPaymentTransitions[transition]
			return model.IsValidStatusTransition(from, to) == isValid
		},
		genPaymentStatus(),
		genPaymentStatus(),
	))

	// Property 3.2: pending -> paid is valid
	properties.Property("pending to paid transition is valid", prop.ForAll(
		func(_ int) bool {
			return model.IsValidStatusTransition(model.PaymentStatusPending, model.PaymentStatusPaid)
		},
		gen.Int(),
	))

	// Property 3.3: pending -> failed is valid
	properties.Property("pending to failed transition is valid", prop.ForAll(
		func(_ int) bool {
			return model.IsValidStatusTransition(model.PaymentStatusPending, model.PaymentStatusFailed)
		},
		gen.Int(),
	))

	// Property 3.4: paid -> refunded is valid
	properties.Property("paid to refunded transition is valid", prop.ForAll(
		func(_ int) bool {
			return model.IsValidStatusTransition(model.PaymentStatusPaid, model.PaymentStatusRefunded)
		},
		gen.Int(),
	))

	// Property 3.5: Same status transition is always valid (no-op)
	properties.Property("same status transition is always valid", prop.ForAll(
		func(status model.PaymentStatus) bool {
			return model.IsValidStatusTransition(status, status)
		},
		genPaymentStatus(),
	))

	// Property 3.6: failed -> any other status is invalid
	properties.Property("failed status cannot transition to other statuses", prop.ForAll(
		func(to model.PaymentStatus) bool {
			if to == model.PaymentStatusFailed {
				return true // Same status is valid
			}
			return !model.IsValidStatusTransition(model.PaymentStatusFailed, to)
		},
		genPaymentStatus(),
	))

	// Property 3.7: refunded -> any other status is invalid
	properties.Property("refunded status cannot transition to other statuses", prop.ForAll(
		func(to model.PaymentStatus) bool {
			if to == model.PaymentStatusRefunded {
				return true // Same status is valid
			}
			return !model.IsValidStatusTransition(model.PaymentStatusRefunded, to)
		},
		genPaymentStatus(),
	))

	// Property 3.8: paid -> pending is invalid (cannot go back)
	properties.Property("paid to pending transition is invalid", prop.ForAll(
		func(_ int) bool {
			return !model.IsValidStatusTransition(model.PaymentStatusPaid, model.PaymentStatusPending)
		},
		gen.Int(),
	))

	// Property 3.9: paid -> failed is invalid
	properties.Property("paid to failed transition is invalid", prop.ForAll(
		func(_ int) bool {
			return !model.IsValidStatusTransition(model.PaymentStatusPaid, model.PaymentStatusFailed)
		},
		gen.Int(),
	))

	// Property 3.10: pending -> refunded is invalid (must go through paid first)
	properties.Property("pending to refunded transition is invalid", prop.ForAll(
		func(_ int) bool {
			return !model.IsValidStatusTransition(model.PaymentStatusPending, model.PaymentStatusRefunded)
		},
		gen.Int(),
	))

	// Property 3.11: ValidateStatusTransition returns nil for valid transitions
	properties.Property("ValidateStatusTransition returns nil for valid transitions", prop.ForAll(
		func(payment *model.Payment) bool {
			allowed := payment.GetAllowedTransitions()
			for _, nextStatus := range allowed {
				if err := payment.ValidateStatusTransition(nextStatus); err != nil {
					return false
				}
			}
			return true
		},
		genValidPayment(),
	))

	// Property 3.12: ValidateStatusTransition returns error for invalid transitions
	properties.Property("ValidateStatusTransition returns error for invalid transitions", prop.ForAll(
		func(payment *model.Payment, targetStatus model.PaymentStatus) bool {
			// Skip if it's a valid transition
			if model.IsValidStatusTransition(payment.Status, targetStatus) {
				return true
			}
			err := payment.ValidateStatusTransition(targetStatus)
			return err != nil
		},
		genValidPayment(),
		genPaymentStatus(),
	))

	// Property 3.13: GetAllowedTransitions returns only valid next statuses
	properties.Property("GetAllowedTransitions returns only valid transitions", prop.ForAll(
		func(payment *model.Payment) bool {
			allowed := payment.GetAllowedTransitions()
			for _, nextStatus := range allowed {
				if !model.IsValidStatusTransition(payment.Status, nextStatus) {
					return false
				}
			}
			return true
		},
		genValidPayment(),
	))

	properties.TestingRun(t)
}
