package auth

import (
	"testing"
	"time"
)

func TestJWTManagerRoundTrip(t *testing.T) {
	manager := NewJWTManager("secret-key", time.Minute)
	token, err := manager.GenerateToken(42, "admin")
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}

	claims, err := manager.VerifyToken(token)
	if err != nil {
		t.Fatalf("VerifyToken error: %v", err)
	}
	if claims.UserID != 42 || claims.Role != "admin" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestExtractTokenFromHeader(t *testing.T) {
	token := "abc.def"
	got, err := ExtractTokenFromHeader("Bearer " + token)
	if err != nil {
		t.Fatalf("ExtractTokenFromHeader error: %v", err)
	}
	if got != token {
		t.Fatalf("expected %s, got %s", token, got)
	}
	if _, err := ExtractTokenFromHeader("Invalid header"); err == nil {
		t.Fatal("expected error for invalid header")
	}
}
