package handler

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"

    "gamelink/internal/auth"
    "gamelink/internal/model"
    "gamelink/internal/repository"
    "gamelink/internal/service"
)

type fakeUserRepoAuth struct{ u *model.User }

func (f *fakeUserRepoAuth) List(context.Context) ([]model.User, error) { return nil, nil }
func (f *fakeUserRepoAuth) ListPaged(context.Context, int, int) ([]model.User, int64, error) {
    return nil, 0, nil
}
func (f *fakeUserRepoAuth) ListWithFilters(context.Context, repository.UserListOptions) ([]model.User, int64, error) {
    return nil, 0, nil
}
func (f *fakeUserRepoAuth) Get(context.Context, uint64) (*model.User, error) { return f.u, nil }
func (f *fakeUserRepoAuth) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return f.u, nil
}
func (f *fakeUserRepoAuth) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	return f.u, nil
}
func (f *fakeUserRepoAuth) Create(context.Context, *model.User) error { return nil }
func (f *fakeUserRepoAuth) Update(context.Context, *model.User) error { return nil }
func (f *fakeUserRepoAuth) Delete(context.Context, uint64) error      { return nil }

func TestAuth_LoginAndRefresh(t *testing.T) {
	gin.SetMode(gin.TestMode)

	pwd, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	user := &model.User{Base: model.Base{ID: 42}, Email: "a@b.com", PasswordHash: string(pwd), Name: "alice", Role: model.RoleUser, Status: model.UserStatusActive}
	repo := &fakeUserRepoAuth{u: user}

	mgr := auth.NewJWTManager("test-secret", 2*time.Second)
	svc := service.NewAuthService(repo, mgr)

	r := gin.New()
	RegisterAuthRoutes(r, svc)

	// login
	body := map[string]string{"username": "a@b.com", "password": "secret123"}
	buf, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login failed: code=%d body=%s", w.Code, w.Body.String())
	}
	var lr struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &lr); err != nil || lr.Data.Token == "" {
		t.Fatalf("invalid login response: %v %s", err, w.Body.String())
	}

	// refresh
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req2.Header.Set("Authorization", "Bearer "+lr.Data.Token)
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("refresh failed: code=%d body=%s", w2.Code, w2.Body.String())
	}
}
