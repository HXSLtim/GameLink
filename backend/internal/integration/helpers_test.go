package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gamelink/internal/model"
)

type apiResp[T any] struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func doJSON(router *gin.Engine, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
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

func uintToStr(v uint64) string {
	return strconv.FormatUint(v, 10)
}

func fakeAuthMiddleware(userID uint64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	}
}

func ctx() context.Context { return context.Background() }

func migratePaymentModels(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(
		&model.Order{},
		&model.Payment{},
	); err != nil {
		t.Fatalf("migrate payment models: %v", err)
	}
}

func migrateServiceItemModels(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(
		&model.Game{},
		&model.Player{},
		&model.ServiceItem{},
	); err != nil {
		t.Fatalf("migrate service item models: %v", err)
	}
}

func migrateWithdrawModels(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(
		&model.User{},
		&model.Player{},
		&model.Withdraw{},
	); err != nil {
		t.Fatalf("migrate withdraw models: %v", err)
	}
}
