package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gamelink/pkg/auth"
	"gamelink/internal/handler"
	"gamelink/internal/model"
	userrepo "gamelink/internal/repository/user"
	authservice "gamelink/internal/service/auth"
	"gamelink/pkg/testutil"
)

type apiResponse[T any] struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type loginResponse struct {
	Token     string     `json:"token"`
	ExpiresAt time.Time  `json:"expires_at"`
	User      model.User `json:"user"`
}

func setupAuthRouter(t *testing.T) (*gin.Engine, *authservice.AuthService, *gormCleanup) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.User{})

	userRepo := userrepo.NewUserRepository(db)
	// Use short TTL so refresh can be exercised in tests.
	jwtManager := auth.NewJWTManager("integration-secret", 5*time.Second)
	authSvc := authservice.NewAuthService(userRepo, jwtManager)

	router := gin.New()
	api := router.Group("/api/v1")
	handler.RegisterAuthRoutes(api, authSvc)

	return router, authSvc, &gormCleanup{t: t, db: db}
}

type gormCleanup struct {
	t  *testing.T
	db *gorm.DB
}

func (c *gormCleanup) Close() {
	if c != nil && c.db != nil {
		testutil.CleanDB(c.t, c.db)
	}
}

func performJSONRequest(router *gin.Engine, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}

	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// TestAuthIntegration_RegisterLoginRefreshMe covers happy-path auth flow with in-memory SQLite.
func TestAuthIntegration_RegisterLoginRefreshMe(t *testing.T) {
	router, _, cleanup := setupAuthRouter(t)
	defer cleanup.Close()

	registerPayload := map[string]interface{}{
		"name":     "Alice",
		"email":    "alice@example.com",
		"password": "secret123",
	}

	// Register
	regResp := performJSONRequest(router, http.MethodPost, "/api/v1/auth/register", registerPayload, "")
	if regResp.Code != http.StatusOK {
		t.Fatalf("register status = %d, body=%s", regResp.Code, regResp.Body.String())
	}

	var regParsed apiResponse[loginResponse]
	if err := json.Unmarshal(regResp.Body.Bytes(), &regParsed); err != nil {
		t.Fatalf("failed to parse register response: %v", err)
	}
	if !regParsed.Success || regParsed.Data.Token == "" {
		t.Fatalf("unexpected register response: %+v", regParsed)
	}

	token := regParsed.Data.Token

	// Login with same credentials
	loginPayload := map[string]interface{}{
		"username": registerPayload["email"],
		"password": registerPayload["password"],
	}
	loginResp := performJSONRequest(router, http.MethodPost, "/api/v1/auth/login", loginPayload, "")
	if loginResp.Code != http.StatusOK {
		t.Fatalf("login status = %d, body=%s", loginResp.Code, loginResp.Body.String())
	}

	var loginParsed apiResponse[loginResponse]
	if err := json.Unmarshal(loginResp.Body.Bytes(), &loginParsed); err != nil {
		t.Fatalf("failed to parse login response: %v", err)
	}
	if !loginParsed.Success || loginParsed.Data.Token == "" {
		t.Fatalf("unexpected login response: %+v", loginParsed)
	}

	// Refresh token
	refreshResp := performJSONRequest(router, http.MethodPost, "/api/v1/auth/refresh", nil, token)
	if refreshResp.Code != http.StatusOK {
		t.Fatalf("refresh status = %d, body=%s", refreshResp.Code, refreshResp.Body.String())
	}
	var refreshParsed apiResponse[struct {
		Token string `json:"token"`
	}]
	if err := json.Unmarshal(refreshResp.Body.Bytes(), &refreshParsed); err != nil {
		t.Fatalf("failed to parse refresh response: %v", err)
	}
	if !refreshParsed.Success || refreshParsed.Data.Token == "" {
		t.Fatalf("unexpected refresh response: %+v", refreshParsed)
	}
	refreshToken := refreshParsed.Data.Token

	// Me endpoint with refreshed token
	meResp := performJSONRequest(router, http.MethodGet, "/api/v1/auth/me", nil, refreshToken)
	if meResp.Code != http.StatusOK {
		t.Fatalf("me status = %d, body=%s", meResp.Code, meResp.Body.String())
	}
	var meParsed apiResponse[loginResponse]
	if err := json.Unmarshal(meResp.Body.Bytes(), &meParsed); err != nil {
		t.Fatalf("failed to parse me response: %v", err)
	}
	if !meParsed.Success {
		t.Fatalf("me response not success: %+v", meParsed)
	}
	if meParsed.Data.User.Email != registerPayload["email"] {
		t.Fatalf("me user email mismatch, got %s", meParsed.Data.User.Email)
	}
}
