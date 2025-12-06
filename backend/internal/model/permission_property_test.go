package model

import (
	"reflect"
	"strings"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// TestPermissionCodeFormatValidation tests Property 1: Permission Code Format Validation
// For any permission code string, if it does not match the pattern {module}.{resource}.{action}
// (three dot-separated segments), the system should reject it with a validation error.
// **Feature: rbac-button-level-permission, Property 1: 权限码格式验证**
// **Validates: Requirements 1.3**
func TestPermissionCodeFormatValidation(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property 1: Valid permission codes (module.resource.action) should pass validation
	properties.Property("valid permission codes should pass validation", prop.ForAll(
		func(module, resource, action string) bool {
			code := module + "." + resource + "." + action
			p := &Permission{Code: code}
			return p.ValidateCode()
		},
		genValidSegment(),
		genValidSegment(),
		genValidSegment(),
	))

	// Property 2: Empty permission code should fail validation
	properties.Property("empty permission code should fail validation", prop.ForAll(
		func(_ int) bool {
			p := &Permission{Code: ""}
			return !p.ValidateCode()
		},
		gen.Int(),
	))

	// Property 3: Permission codes with fewer than 3 segments should fail validation
	properties.Property("permission codes with fewer than 3 segments should fail validation", prop.ForAll(
		func(segment1, segment2 string) bool {
			// Test single segment
			p1 := &Permission{Code: segment1}
			if p1.ValidateCode() {
				return false
			}

			// Test two segments
			p2 := &Permission{Code: segment1 + "." + segment2}
			if p2.ValidateCode() {
				return false
			}

			return true
		},
		genValidSegment(),
		genValidSegment(),
	))

	// Property 4: Permission codes with more than 3 segments should fail validation
	properties.Property("permission codes with more than 3 segments should fail validation", prop.ForAll(
		func(s1, s2, s3, s4 string) bool {
			code := s1 + "." + s2 + "." + s3 + "." + s4
			p := &Permission{Code: code}
			return !p.ValidateCode()
		},
		genValidSegment(),
		genValidSegment(),
		genValidSegment(),
		genValidSegment(),
	))

	// Property 5: Permission codes with segments starting with numbers should fail validation
	properties.Property("permission codes with segments starting with numbers should fail validation", prop.ForAll(
		func(num int, validSeg1, validSeg2 string) bool {
			// Segment starting with number
			invalidSegment := string(rune('0'+num%10)) + validSeg1

			// Test in module position
			p1 := &Permission{Code: invalidSegment + "." + validSeg1 + "." + validSeg2}
			if p1.ValidateCode() {
				return false
			}

			// Test in resource position
			p2 := &Permission{Code: validSeg1 + "." + invalidSegment + "." + validSeg2}
			if p2.ValidateCode() {
				return false
			}

			// Test in action position
			p3 := &Permission{Code: validSeg1 + "." + validSeg2 + "." + invalidSegment}
			if p3.ValidateCode() {
				return false
			}

			return true
		},
		gen.IntRange(0, 9),
		genValidSegment(),
		genValidSegment(),
	))

	// Property 6: Permission codes with uppercase letters should fail validation
	properties.Property("permission codes with uppercase letters should fail validation", prop.ForAll(
		func(validSeg1, validSeg2, validSeg3 string) bool {
			// Add uppercase letter to first segment
			upperSegment := strings.ToUpper(validSeg1[:1]) + validSeg1[1:]
			p := &Permission{Code: upperSegment + "." + validSeg2 + "." + validSeg3}
			return !p.ValidateCode()
		},
		genValidSegment(),
		genValidSegment(),
		genValidSegment(),
	))

	// Property 7: Permission codes with special characters should fail validation
	properties.Property("permission codes with special characters should fail validation", prop.ForAll(
		func(validSeg1, validSeg2, validSeg3 string, specialChar rune) bool {
			// Insert special character into segment
			invalidSegment := validSeg1 + string(specialChar)
			p := &Permission{Code: invalidSegment + "." + validSeg2 + "." + validSeg3}
			return !p.ValidateCode()
		},
		genValidSegment(),
		genValidSegment(),
		genValidSegment(),
		gen.OneConstOf('@', '#', '$', '%', '^', '&', '*', '!', '-', '_', ' '),
	))

	// Property 8: Permission codes with empty segments should fail validation
	properties.Property("permission codes with empty segments should fail validation", prop.ForAll(
		func(validSeg1, validSeg2 string) bool {
			// Empty module
			p1 := &Permission{Code: "." + validSeg1 + "." + validSeg2}
			if p1.ValidateCode() {
				return false
			}

			// Empty resource
			p2 := &Permission{Code: validSeg1 + ".." + validSeg2}
			if p2.ValidateCode() {
				return false
			}

			// Empty action
			p3 := &Permission{Code: validSeg1 + "." + validSeg2 + "."}
			if p3.ValidateCode() {
				return false
			}

			return true
		},
		genValidSegment(),
		genValidSegment(),
	))

	// Property 9: GetModule, GetResource, GetAction should correctly extract parts from valid codes
	properties.Property("GetModule, GetResource, GetAction should correctly extract parts", prop.ForAll(
		func(module, resource, action string) bool {
			code := module + "." + resource + "." + action
			p := &Permission{Code: code}

			return p.GetModule() == module &&
				p.GetResource() == resource &&
				p.GetAction() == action
		},
		genValidSegment(),
		genValidSegment(),
		genValidSegment(),
	))

	// Property 10: Whitespace in permission codes should fail validation
	properties.Property("permission codes with whitespace should fail validation", prop.ForAll(
		func(validSeg1, validSeg2, validSeg3 string) bool {
			// Leading whitespace
			p1 := &Permission{Code: " " + validSeg1 + "." + validSeg2 + "." + validSeg3}
			if p1.ValidateCode() {
				return false
			}

			// Trailing whitespace
			p2 := &Permission{Code: validSeg1 + "." + validSeg2 + "." + validSeg3 + " "}
			if p2.ValidateCode() {
				return false
			}

			// Whitespace between segments
			p3 := &Permission{Code: validSeg1 + " ." + validSeg2 + "." + validSeg3}
			if p3.ValidateCode() {
				return false
			}

			return true
		},
		genValidSegment(),
		genValidSegment(),
		genValidSegment(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// genValidSegment generates a valid segment for permission codes.
// A valid segment starts with a lowercase letter and contains only lowercase letters and digits.
func genValidSegment() gopter.Gen {
	// Generate lowercase letters for the first character (a-z)
	firstChar := gen.IntRange(0, 25).Map(func(i int) rune {
		return rune('a' + i)
	})

	// Generate length for rest of string (0-9 additional characters)
	restLength := gen.IntRange(0, 9)

	// Generate each character as lowercase letter or digit
	charGen := gen.IntRange(0, 35).Map(func(i int) rune {
		if i < 26 {
			return rune('a' + i)
		}
		return rune('0' + (i - 26))
	})

	return gopter.CombineGens(firstChar, restLength).FlatMap(func(vals interface{}) gopter.Gen {
		v := vals.([]interface{})
		first := v[0].(rune)
		length := v[1].(int)

		if length == 0 {
			return gen.Const(string(first))
		}

		return gen.SliceOfN(length, charGen).Map(func(rest []rune) string {
			return string(first) + string(rest)
		})
	}, reflect.TypeOf(""))
}

// TestPermissionCodeExamples tests specific examples of valid and invalid permission codes.
// This complements the property tests with concrete examples.
func TestPermissionCodeExamples(t *testing.T) {
	validCodes := []string{
		"admin.users.create",
		"admin.users.read",
		"admin.users.update",
		"admin.users.delete",
		"admin.orders.read",
		"admin.roles.assign",
		"admin.audit.export",
		"player.profile.update",
		"user.orders.create",
		"a.b.c",
		"module1.resource2.action3",
	}

	invalidCodes := []string{
		"",                         // Empty
		"admin",                    // Single segment
		"admin.users",              // Two segments
		"admin.users.create.extra", // Four segments
		"Admin.users.create",       // Uppercase in module
		"admin.Users.create",       // Uppercase in resource
		"admin.users.Create",       // Uppercase in action
		"1admin.users.create",      // Starts with number
		"admin.1users.create",      // Resource starts with number
		"admin.users.1create",      // Action starts with number
		"admin-users.create.read",  // Hyphen in segment
		"admin_users.create.read",  // Underscore in segment
		"admin users.create.read",  // Space in segment
		".users.create",            // Empty module
		"admin..create",            // Empty resource
		"admin.users.",             // Empty action
		"admin.users.create ",      // Trailing space
		" admin.users.create",      // Leading space
		"admin.users .create",      // Space before dot
		"admin. users.create",      // Space after dot
	}

	for _, code := range validCodes {
		p := &Permission{Code: code}
		if !p.ValidateCode() {
			t.Errorf("Expected valid code %q to pass validation, but it failed", code)
		}
	}

	for _, code := range invalidCodes {
		p := &Permission{Code: code}
		if p.ValidateCode() {
			t.Errorf("Expected invalid code %q to fail validation, but it passed", code)
		}
	}
}
