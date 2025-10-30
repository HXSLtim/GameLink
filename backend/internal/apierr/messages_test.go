package apierr

import "testing"

// TestErrorMessagesNotEmpty 确保所有错误消息常量都不为空。
func TestErrorMessagesNotEmpty(t *testing.T) {
	errorMessages := []struct {
		name  string
		value string
	}{
		{"ErrInvalidJSONPayload", ErrInvalidJSONPayload},
		{"ErrInvalidID", ErrInvalidID},
		{"ErrInvalidPage", ErrInvalidPage},
		{"ErrInvalidPageSize", ErrInvalidPageSize},
		{"ErrInvalidUserID", ErrInvalidUserID},
		{"ErrInvalidOrderID", ErrInvalidOrderID},
		{"ErrGameNotFound", ErrGameNotFound},
		{"ErrPlayerNotFound", ErrPlayerNotFound},
		{"ErrInvalidPlayerID", ErrInvalidPlayerID},
		{"ErrInvalidGameID", ErrInvalidGameID},
		{"ErrInvalidDateFrom", ErrInvalidDateFrom},
		{"ErrInvalidDateTo", ErrInvalidDateTo},
		{"ErrInvalidPaidAt", ErrInvalidPaidAt},
		{"ErrInvalidRefundedAt", ErrInvalidRefundedAt},
		{"ErrInvalidScheduledStart", ErrInvalidScheduledStart},
		{"ErrInvalidScheduledEnd", ErrInvalidScheduledEnd},
		{"ErrInvalidEmailFormat", ErrInvalidEmailFormat},
		{"ErrInvalidPhoneFormat", ErrInvalidPhoneFormat},
		{"ErrMissingFieldsOrShortPassword", ErrMissingFieldsOrShortPassword},
		{"ErrMissingRequiredFields", ErrMissingRequiredFields},
		{"ErrUserNotFound", ErrUserNotFound},
		{"ErrOrderNotFound", ErrOrderNotFound},
		{"ErrPaymentNotFound", ErrPaymentNotFound},
		{"ErrOrderInvalidTransition", ErrOrderInvalidTransition},
		{"ErrInvalidPaymentPayload", ErrInvalidPaymentPayload},
		{"ErrInvalidOrderPayload", ErrInvalidOrderPayload},
	}

	for _, tc := range errorMessages {
		t.Run(tc.name, func(t *testing.T) {
			if tc.value == "" {
				t.Errorf("%s should not be empty", tc.name)
			}
		})
	}
}

// TestErrorMessagesUnique 确保错误消息是唯一的（可选）。
func TestErrorMessagesUnique(t *testing.T) {
	messages := map[string]string{
		"ErrInvalidJSONPayload":           ErrInvalidJSONPayload,
		"ErrInvalidID":                    ErrInvalidID,
		"ErrInvalidPage":                  ErrInvalidPage,
		"ErrInvalidPageSize":              ErrInvalidPageSize,
		"ErrInvalidUserID":                ErrInvalidUserID,
		"ErrInvalidOrderID":               ErrInvalidOrderID,
		"ErrGameNotFound":                 ErrGameNotFound,
		"ErrPlayerNotFound":               ErrPlayerNotFound,
		"ErrInvalidPlayerID":              ErrInvalidPlayerID,
		"ErrInvalidGameID":                ErrInvalidGameID,
		"ErrInvalidDateFrom":              ErrInvalidDateFrom,
		"ErrInvalidDateTo":                ErrInvalidDateTo,
		"ErrInvalidPaidAt":                ErrInvalidPaidAt,
		"ErrInvalidRefundedAt":            ErrInvalidRefundedAt,
		"ErrInvalidScheduledStart":        ErrInvalidScheduledStart,
		"ErrInvalidScheduledEnd":          ErrInvalidScheduledEnd,
		"ErrInvalidEmailFormat":           ErrInvalidEmailFormat,
		"ErrInvalidPhoneFormat":           ErrInvalidPhoneFormat,
		"ErrMissingFieldsOrShortPassword": ErrMissingFieldsOrShortPassword,
		"ErrMissingRequiredFields":        ErrMissingRequiredFields,
		"ErrUserNotFound":                 ErrUserNotFound,
		"ErrOrderNotFound":                ErrOrderNotFound,
		"ErrPaymentNotFound":              ErrPaymentNotFound,
		"ErrOrderInvalidTransition":       ErrOrderInvalidTransition,
		"ErrInvalidPaymentPayload":        ErrInvalidPaymentPayload,
		"ErrInvalidOrderPayload":          ErrInvalidOrderPayload,
	}

	seen := make(map[string]string)
	for name, msg := range messages {
		if existing, found := seen[msg]; found {
			t.Logf("Warning: duplicate message '%s' found in %s and %s", msg, name, existing)
		}
		seen[msg] = name
	}
}
