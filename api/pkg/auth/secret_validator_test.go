package auth

import (
	"strings"
	"testing"
)

func TestValidateJWTSecret(t *testing.T) {
	tests := []struct {
		name      string
		secret    string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid strong secret",
			secret:    "Abc123!@#Xyz789$%^Def456&*()Ghi012+=-Jkl345~`Mno678[]{}Pqr901<>?",
			wantError: false,
		},
		{
			name:      "too short - less than 64 chars",
			secret:    "Abc123!@#Xyz789$%^Def456&*()",
			wantError: true,
			errorMsg:  "at least 64 characters",
		},
		{
			name:      "missing uppercase",
			secret:    "abc123!@#xyz789$%^def456&*()ghi012+=-jkl345~`mno678[]{}pqr901<>?",
			wantError: true,
			errorMsg:  "uppercase letter",
		},
		{
			name:      "missing lowercase",
			secret:    "ABC123!@#XYZ789$%^DEF456&*()GHI012+=-JKL345~`MNO678[]{}PQR901<>?",
			wantError: true,
			errorMsg:  "lowercase letter",
		},
		{
			name:      "missing digit",
			secret:    "Abc!@#Xyz$%^Def&*()Ghi+=-Jkl~`Mno[]{}Pqr<>?Stu|\\Vwx;:'Yza\",./?Bcd",
			wantError: true,
			errorMsg:  "digit",
		},
		{
			name:      "missing special character",
			secret:    "Abc123XyzDef456GhiJkl789MnoPqr012StuVwx345YzaAbc678DefGhi901Jkl234",
			wantError: true,
			errorMsg:  "special character",
		},
		{
			name:      "contains weak pattern - 'secret'",
			secret:    "MySecretKey123!@#Xyz789$%^Def456&*()Ghi012+=-Jkl345~`Mno678[]{}Pqr",
			wantError: true,
			errorMsg:  "weak patterns",
		},
		{
			name:      "contains weak pattern - 'password'",
			secret:    "MyPassword123!@#Xyz789$%^Def456&*()Ghi012+=-Jkl345~`Mno678[]{}Pqr",
			wantError: true,
			errorMsg:  "weak patterns",
		},
		{
			name:      "repeated characters",
			secret:    "Aaaaaaa123!@#Xyz789$%^Def456&*()Ghi012+=-Jkl345~`Mno678[]{}Pqr901",
			wantError: true,
			errorMsg:  "repeated characters",
		},
		{
			name:      "insufficient entropy",
			secret:    "AAAA1111!!!!BBBBaaaa2222@@@@CCCCbbbb3333####DDDDcccc4444$$$$EEEEdddd",
			wantError: true,
			errorMsg:  "insufficient entropy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateJWTSecret(tt.secret)
			if tt.wantError {
				if err == nil {
					t.Errorf("ValidateJWTSecret() expected error but got nil")
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("ValidateJWTSecret() error = %v, want error containing %v", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateJWTSecret() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestIsSecretWeak(t *testing.T) {
	tests := []struct {
		name   string
		secret string
		want   bool
	}{
		{
			name:   "strong secret",
			secret: "Abc123!@#Xyz789$%^Def456&*()Ghi012+=-Jkl345~`Mno678[]{}Pqr901<>?",
			want:   false,
		},
		{
			name:   "too short",
			secret: "Abc123!@#",
			want:   true,
		},
		{
			name:   "all same character",
			secret: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			want:   true,
		},
		{
			name:   "starts with 'secret'",
			secret: "secret123456789012345678901234567890",
			want:   true,
		},
		{
			name:   "starts with 'password'",
			secret: "password123456789012345678901234567890",
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSecretWeak(tt.secret); got != tt.want {
				t.Errorf("IsSecretWeak() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateSecretHash(t *testing.T) {
	secret := "test-secret-key-123"
	hash1 := GenerateSecretHash(secret)
	hash2 := GenerateSecretHash(secret)

	// Same secret should produce same hash
	if hash1 != hash2 {
		t.Errorf("GenerateSecretHash() produced different hashes for same secret")
	}

	// Hash should be 64 characters (SHA-256 hex)
	if len(hash1) != 64 {
		t.Errorf("GenerateSecretHash() hash length = %d, want 64", len(hash1))
	}

	// Different secret should produce different hash
	differentSecret := "different-secret-key-456"
	hash3 := GenerateSecretHash(differentSecret)
	if hash1 == hash3 {
		t.Errorf("GenerateSecretHash() produced same hash for different secrets")
	}
}

func TestValidateJWTSecretOrPanic(t *testing.T) {
	// Test that weak secret causes panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("ValidateJWTSecretOrPanic() did not panic for weak secret")
		}
	}()

	ValidateJWTSecretOrPanic("weak")
}

func TestHasRepeatedChars(t *testing.T) {
	tests := []struct {
		name      string
		s         string
		maxRepeat int
		want      bool
	}{
		{"no repeats", "abcdef", 3, false},
		{"3 repeats below threshold", "aaa", 4, false},
		{"5 repeats at threshold", "aaaaa", 5, true},
		{"6 repeats above threshold", "aaaaaa", 5, true},
		{"repeats in middle", "abcccccdef", 5, true},
		{"multiple separate repeats", "aaabbbccc", 3, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasRepeatedChars(tt.s, tt.maxRepeat); got != tt.want {
				t.Errorf("hasRepeatedChars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasSufficientEntropy(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"high entropy", "Abc123!@#Xyz789$%^", true},
		{"low entropy - repeated chars", "AAAA1111!!!!BBBB", false},
		{"medium entropy", "abcdefghijklmnop", true},
		{"very low entropy", "aaaaaaaaaaaaaaaa", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasSufficientEntropy(tt.s); got != tt.want {
				t.Errorf("hasSufficientEntropy() = %v, want %v", got, tt.want)
			}
		})
	}
}
