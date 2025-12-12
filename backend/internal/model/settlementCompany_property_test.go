package model_test

import (
	"strings"
	"testing"
	"time"

	"gamelink/internal/model"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: payment-finance-module, Property 11: 统一社会信用代码格式验证**
// **Validates: Requirements 11.2, 15.2**
//
// Property 11: Unified Social Credit Code Format Validation
// *For any* settlement company or collection entity creation/update operation,
// the unified social credit code must conform to the 18-character standard format
// and be unique within the system.

// validCreditCodeChars contains all valid characters for credit code positions
// Excludes I, O, S, V, Z to avoid confusion with numbers
const validCreditCodeChars = "0123456789ABCDEFGHJKLMNPQRTUWXY"

// genValidCreditCodeChar generates a valid character for credit code
func genValidCreditCodeChar() gopter.Gen {
	return gen.IntRange(0, len(validCreditCodeChars)-1).Map(func(i int) byte {
		return validCreditCodeChars[i]
	})
}

// genValidCreditCode generates a valid 18-character credit code format
// Note: This generates format-valid codes but may not pass checksum validation
func genValidCreditCode() gopter.Gen {
	return gopter.CombineGens(
		genValidCreditCodeChar(),                  // Position 1: Registration management department code
		genValidCreditCodeChar(),                  // Position 2: Organization category code
		gen.IntRange(100000, 999999),              // Positions 3-8: Administrative division code (6 digits)
		gen.SliceOfN(9, genValidCreditCodeChar()), // Positions 9-17: Subject identifier (9 chars)
		genValidCreditCodeChar(),                  // Position 18: Check code
	).Map(func(values []interface{}) string {
		pos1 := values[0].(byte)
		pos2 := values[1].(byte)
		divisionCode := values[2].(int)
		subjectChars := values[3].([]byte)
		checkCode := values[4].(byte)

		var sb strings.Builder
		sb.WriteByte(pos1)
		sb.WriteByte(pos2)
		sb.WriteString(strings.Repeat("0", 6-len(string(rune(divisionCode)))))
		sb.WriteString(string(rune('0' + divisionCode/100000)))
		sb.WriteString(string(rune('0' + (divisionCode/10000)%10)))
		sb.WriteString(string(rune('0' + (divisionCode/1000)%10)))
		sb.WriteString(string(rune('0' + (divisionCode/100)%10)))
		sb.WriteString(string(rune('0' + (divisionCode/10)%10)))
		sb.WriteString(string(rune('0' + divisionCode%10)))
		for _, c := range subjectChars {
			sb.WriteByte(c)
		}
		sb.WriteByte(checkCode)
		return sb.String()
	})
}

// genValidCreditCodeSimple generates a simpler valid credit code format
func genValidCreditCodeSimple() gopter.Gen {
	return gen.SliceOfN(18, genValidCreditCodeChar()).Map(func(chars []byte) string {
		// Ensure positions 3-8 are digits
		for i := 2; i < 8; i++ {
			chars[i] = byte('0' + (int(chars[i]) % 10))
		}
		return string(chars)
	})
}

// genInvalidCreditCode generates an invalid credit code
func genInvalidCreditCode() gopter.Gen {
	return gen.OneGenOf(
		// Wrong length (too short)
		gen.SliceOfN(17, genValidCreditCodeChar()).Map(func(chars []byte) string {
			return string(chars)
		}),
		// Wrong length (too long)
		gen.SliceOfN(19, genValidCreditCodeChar()).Map(func(chars []byte) string {
			return string(chars)
		}),
		// Contains invalid characters (I, O, S, V, Z)
		gen.SliceOfN(18, genValidCreditCodeChar()).Map(func(chars []byte) string {
			// Replace one character with an invalid one
			invalidChars := []byte{'I', 'O', 'S', 'V', 'Z'}
			chars[5] = invalidChars[int(chars[0])%len(invalidChars)]
			return string(chars)
		}),
		// Empty string
		gen.Const(""),
		// Contains lowercase letters
		gen.SliceOfN(18, genValidCreditCodeChar()).Map(func(chars []byte) string {
			chars[10] = 'a' // lowercase
			return string(chars)
		}),
		// Contains special characters
		gen.SliceOfN(18, genValidCreditCodeChar()).Map(func(chars []byte) string {
			chars[10] = '@'
			return string(chars)
		}),
	)
}

// TestProperty11_CreditCodeFormatValidation tests credit code format validation
func TestProperty11_CreditCodeFormatValidation(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property 11.1: Valid format credit codes should pass basic validation
	properties.Property("valid format credit code should pass validation", prop.ForAll(
		func(code string) bool {
			return model.ValidateCreditCode(code)
		},
		genValidCreditCodeSimple(),
	))

	// Property 11.2: Credit code must be exactly 18 characters
	properties.Property("credit code must be exactly 18 characters", prop.ForAll(
		func(code string) bool {
			if len(code) != 18 {
				return !model.ValidateCreditCode(code)
			}
			return true // Skip codes that happen to be 18 chars
		},
		gen.AnyString(),
	))

	// Property 11.3: Credit codes with invalid characters should fail
	properties.Property("credit code with invalid characters should fail", prop.ForAll(
		func(code string) bool {
			// Check if code contains invalid characters
			hasInvalid := false
			for _, c := range code {
				if !strings.ContainsRune(validCreditCodeChars, c) {
					hasInvalid = true
					break
				}
			}
			if hasInvalid {
				return !model.ValidateCreditCode(code)
			}
			return true
		},
		genInvalidCreditCode(),
	))

	// Property 11.4: Empty credit code should fail validation
	properties.Property("empty credit code should fail validation", prop.ForAll(
		func(_ int) bool {
			return !model.ValidateCreditCode("")
		},
		gen.Int(),
	))

	// Property 11.5: Credit code with spaces should fail validation
	properties.Property("credit code with spaces should fail validation", prop.ForAll(
		func(code string) bool {
			codeWithSpaces := code[:9] + " " + code[10:]
			return !model.ValidateCreditCode(codeWithSpaces)
		},
		genValidCreditCodeSimple(),
	))

	// Property 11.6: Credit code validation is case-sensitive (lowercase should fail)
	properties.Property("lowercase credit code should fail validation", prop.ForAll(
		func(code string) bool {
			lowerCode := strings.ToLower(code)
			// If the code contains letters, lowercase version should fail
			if lowerCode != code {
				return !model.ValidateCreditCode(lowerCode)
			}
			return true
		},
		genValidCreditCodeSimple(),
	))

	// Property 11.7: Credit code length validation
	properties.Property("credit code length must be 18", prop.ForAll(
		func(length int) bool {
			if length == 18 {
				return true // Skip valid length
			}
			// Generate a code of wrong length
			code := strings.Repeat("A", length)
			return !model.ValidateCreditCode(code)
		},
		gen.IntRange(0, 30),
	))

	// Property 11.8: Valid credit codes should have consistent validation results
	properties.Property("validation should be deterministic", prop.ForAll(
		func(code string) bool {
			result1 := model.ValidateCreditCode(code)
			result2 := model.ValidateCreditCode(code)
			return result1 == result2
		},
		gen.AnyString(),
	))

	// Property 11.9: Credit code with only digits in positions 3-8 should pass format check
	properties.Property("positions 3-8 must be digits", prop.ForAll(
		func(code string) bool {
			if len(code) != 18 {
				return true // Skip invalid length
			}
			// Check if positions 3-8 (index 2-7) are digits
			for i := 2; i < 8; i++ {
				if code[i] < '0' || code[i] > '9' {
					return !model.ValidateCreditCode(code)
				}
			}
			return true
		},
		gen.AnyString(),
	))

	properties.TestingRun(t)
}

// TestCreditCodeKnownValues tests credit code validation with known valid/invalid codes
func TestCreditCodeKnownValues(t *testing.T) {
	// Known valid credit code formats (format valid, checksum may vary)
	validFormats := []string{
		"91110000100000000A", // Example format
		"11110000123456789X",
		"91440300MA5EQRQP0K", // Real format example
	}

	for _, code := range validFormats {
		if !model.ValidateCreditCode(code) {
			t.Errorf("Expected valid credit code format: %s", code)
		}
	}

	// Known invalid credit codes
	invalidCodes := []string{
		"",                    // Empty
		"12345678901234567",   // 17 chars
		"1234567890123456789", // 19 chars
		"9111000010000000OA",  // Contains 'O' (invalid)
		"9111000010000000IA",  // Contains 'I' (invalid)
		"9111000010000000SA",  // Contains 'S' (invalid)
		"9111000010000000VA",  // Contains 'V' (invalid)
		"9111000010000000ZA",  // Contains 'Z' (invalid)
		"91110000100000000a",  // Lowercase
		"9111 000100000000A",  // Contains space
	}

	for _, code := range invalidCodes {
		if model.ValidateCreditCode(code) {
			t.Errorf("Expected invalid credit code: %s", code)
		}
	}
}

// TestSettlementCompanyStatus tests company status methods
func TestSettlementCompanyStatus(t *testing.T) {
	company := &model.SettlementCompany{
		Status: model.CompanyStatusActive,
	}

	if !company.IsActive() {
		t.Error("Expected company to be active")
	}

	company.Status = model.CompanyStatusInactive
	if company.IsActive() {
		t.Error("Expected company to be inactive")
	}
}

// TestPlayerCompanyAssignmentEffective tests assignment effectiveness
func TestPlayerCompanyAssignmentEffective(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	// Property: Assignment with future effective date should not be effective
	properties.Property("future assignment should not be effective", prop.ForAll(
		func(daysInFuture int) bool {
			futureDate := time.Now().AddDate(0, 0, daysInFuture)
			assignment := &model.PlayerCompanyAssignment{
				EffectiveDate: futureDate,
			}
			return !assignment.IsEffective()
		},
		gen.IntRange(1, 365),
	))

	properties.TestingRun(t)
}
