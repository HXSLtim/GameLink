package player

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
)

func TestRespondJSON_Player(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	resp := model.APIResponse[string]{
		Success: true,
		Code:    http.StatusOK,
		Message: "Test message",
		Data:    "test data",
	}

	respondJSON(c, http.StatusOK, resp)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	if w.Body.String() == "" {
		t.Fatal("expected non-empty response body")
	}
}

func TestRespondError_Player(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	respondError(c, http.StatusBadRequest, "test error")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	if w.Body.String() == "" {
		t.Fatal("expected non-empty response body")
	}
}

func TestRespondError_Player_InternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	respondError(c, http.StatusInternalServerError, "internal error")

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}
