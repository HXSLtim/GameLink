package auth

import (
    "os"
    "testing"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

func TestJWT_GenerateVerifyAndRefresh(t *testing.T) {
    os.Setenv("JWT_MAX_REFRESH", "1h")
    m := NewJWTManager("secret", 1*time.Second)
    tok, err := m.GenerateToken(123, "super_admin")
    if err != nil { t.Fatalf("generate: %v", err) }
    claims, err := m.VerifyToken(tok)
    if err != nil { t.Fatalf("verify: %v", err) }
    if claims.UserID != 123 || claims.Role != "super_admin" { t.Fatalf("claims mismatch") }
    if IsTokenExpired(claims) { t.Fatalf("unexpected expired") }
    time.Sleep(1200 * time.Millisecond)
    if !IsTokenExpired(claims) { t.Fatalf("expected expired") }
    // refresh when within window
    claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(0))
    claims.IssuedAt = jwt.NewNumericDate(time.Now())
    _, err = m.RefreshToken(claims)
    if err != nil { t.Fatalf("refresh: %v", err) }
}

func TestExtractTokenFromHeader(t *testing.T) {
    if _, err := ExtractTokenFromHeader(""); err == nil { t.Fatalf("expected error") }
    if _, err := ExtractTokenFromHeader("Bad abc"); err == nil { t.Fatalf("expected error") }
    tok := "Bearer xyz"
    out, err := ExtractTokenFromHeader(tok)
    if err != nil || out != "xyz" { t.Fatalf("extract: %v", err) }
}
