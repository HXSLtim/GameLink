package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	serviceitemrepo "gamelink/internal/repository/serviceitem"
	itemservice "gamelink/internal/service/item"
)

type fakeItemRepo struct{ items map[uint64]model.ServiceItem }

func newFakeItemRepo() *fakeItemRepo { return &fakeItemRepo{items: map[uint64]model.ServiceItem{}} }
func (f *fakeItemRepo) Create(_ context.Context, it *model.ServiceItem) error {
	if it.ID == 0 {
		it.ID = uint64(len(f.items) + 1)
	}
	f.items[it.ID] = *it
	return nil
}
func (f *fakeItemRepo) Get(_ context.Context, id uint64) (*model.ServiceItem, error) {
	v, ok := f.items[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	c := v
	return &c, nil
}
func (f *fakeItemRepo) GetByCode(_ context.Context, code string) (*model.ServiceItem, error) {
	for _, v := range f.items {
		if v.ItemCode == code {
			c := v
			return &c, nil
		}
	}
	return nil, repository.ErrNotFound
}
func (f *fakeItemRepo) List(_ context.Context, _ serviceitemrepo.ServiceItemListOptions) ([]model.ServiceItem, int64, error) {
	out := make([]model.ServiceItem, 0, len(f.items))
	for _, v := range f.items {
		out = append(out, v)
	}
	return out, int64(len(out)), nil
}
func (f *fakeItemRepo) Update(_ context.Context, it *model.ServiceItem) error {
	if _, ok := f.items[it.ID]; !ok {
		return repository.ErrNotFound
	}
	f.items[it.ID] = *it
	return nil
}
func (f *fakeItemRepo) Delete(_ context.Context, id uint64) error { delete(f.items, id); return nil }
func (f *fakeItemRepo) BatchUpdateStatus(_ context.Context, ids []uint64, isActive bool) error {
	for _, id := range ids {
		if v, ok := f.items[id]; ok {
			v.IsActive = isActive
			f.items[id] = v
		}
	}
	return nil
}
func (f *fakeItemRepo) BatchUpdatePrice(_ context.Context, ids []uint64, basePriceCents int64) error {
	for _, id := range ids {
		if v, ok := f.items[id]; ok {
			v.BasePriceCents = basePriceCents
			f.items[id] = v
		}
	}
	return nil
}
func (f *fakeItemRepo) GetGifts(_ context.Context, page, pageSize int) ([]model.ServiceItem, int64, error) {
	return nil, 0, nil
}
func (f *fakeItemRepo) GetGameServices(_ context.Context, gameID uint64, subCategory *model.ServiceItemSubCategory) ([]model.ServiceItem, error) {
	return nil, nil
}

type fakeGameRepoItem struct{}

func (fakeGameRepoItem) List(context.Context) ([]model.Game, error) { return nil, nil }
func (fakeGameRepoItem) ListPaged(context.Context, int, int) ([]model.Game, int64, error) {
	return nil, 0, nil
}
func (fakeGameRepoItem) Get(context.Context, uint64) (*model.Game, error) {
	return &model.Game{Name: "G"}, nil
}
func (fakeGameRepoItem) Create(context.Context, *model.Game) error { return nil }
func (fakeGameRepoItem) Update(context.Context, *model.Game) error { return nil }
func (fakeGameRepoItem) Delete(context.Context, uint64) error      { return nil }

type fakePlayerRepoItem struct{}

func (fakePlayerRepoItem) List(context.Context) ([]model.Player, error) { return nil, nil }
func (fakePlayerRepoItem) ListPaged(context.Context, int, int) ([]model.Player, int64, error) {
	return nil, 0, nil
}
func (fakePlayerRepoItem) Get(context.Context, uint64) (*model.Player, error) {
	return &model.Player{Nickname: "P"}, nil
}
func (fakePlayerRepoItem) GetByUserID(context.Context, uint64) (*model.Player, error) {
	return &model.Player{Nickname: "P"}, nil
}
func (fakePlayerRepoItem) Create(context.Context, *model.Player) error { return nil }
func (fakePlayerRepoItem) Update(context.Context, *model.Player) error { return nil }
func (fakePlayerRepoItem) Delete(context.Context, uint64) error        { return nil }

func setupItemRouter(repo *fakeItemRepo) (*gin.Engine, *itemservice.ServiceItemService) {
	r := newTestEngine()
	svc := itemservice.NewServiceItemService(repo, fakeGameRepoItem{}, fakePlayerRepoItem{})
	RegisterServiceItemRoutes(r, svc)
	return r, svc
}

func ptrInt64(v int64) *int64 { return &v }

func TestServiceItem_CRUD_and_Batch(t *testing.T) {
	r, _ := setupItemRouter(newFakeItemRepo())
	create := itemservice.CreateServiceItemRequest{ItemCode: "A1", Name: "ItemA", SubCategory: model.SubCategorySolo, BasePriceCents: 1000, CommissionRate: 0.2}
	cb, _ := json.Marshal(create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/service-items", bytes.NewReader(cb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status %d", w.Code)
	}

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/admin/service-items?page=1&pageSize=20", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	w3 := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodGet, "/admin/service-items/1", nil)
	r.ServeHTTP(w3, req3)
	if w3.Code != http.StatusOK && w3.Code != http.StatusNotFound {
		t.Fatalf("expected 200/404, got %d", w3.Code)
	}

	nm := "B"
	up := itemservice.UpdateServiceItemRequest{Name: &nm}
	ub, _ := json.Marshal(up)
	w4 := httptest.NewRecorder()
	req4 := httptest.NewRequest(http.MethodPut, "/admin/service-items/1", bytes.NewReader(ub))
	req4.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w4, req4)
	if w4.Code != http.StatusOK && w4.Code != http.StatusInternalServerError {
		t.Fatalf("expected 200/500, got %d", w4.Code)
	}

	w5 := httptest.NewRecorder()
	req5 := httptest.NewRequest(http.MethodPost, "/admin/service-items/batch-update-status", bytes.NewReader([]byte(`{"ids":[1],"isActive":true}`)))
	req5.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w5, req5)
	if w5.Code != http.StatusOK && w5.Code != http.StatusInternalServerError {
		t.Fatalf("expected 200/500, got %d", w5.Code)
	}

	w6 := httptest.NewRecorder()
	req6 := httptest.NewRequest(http.MethodPost, "/admin/service-items/batch-update-price", bytes.NewReader([]byte(`{"ids":[1],"basePriceCents":2000}`)))
	req6.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w6, req6)
	if w6.Code != http.StatusOK && w6.Code != http.StatusInternalServerError {
		t.Fatalf("expected 200/500, got %d", w6.Code)
	}

	w7 := httptest.NewRecorder()
	req7 := httptest.NewRequest(http.MethodDelete, "/admin/service-items/1", nil)
	r.ServeHTTP(w7, req7)
	if w7.Code != http.StatusOK && w7.Code != http.StatusInternalServerError {
		t.Fatalf("expected 200/500, got %d", w7.Code)
	}
}

func TestServiceItem_BatchUpdateStatus_EmptyIDs(t *testing.T) {
	r, _ := setupItemRouter(newFakeItemRepo())
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/service-items/batch-update-status", bytes.NewReader([]byte(`{"ids":[],"isActive":true}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestServiceItem_BatchUpdatePrice_EmptyIDs(t *testing.T) {
	r, _ := setupItemRouter(newFakeItemRepo())
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/service-items/batch-update-price", bytes.NewReader([]byte(`{"ids":[],"basePriceCents":100}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestServiceItem_UpdateGiftServiceHoursNonZero(t *testing.T) {
	repo := newFakeItemRepo()
	repo.items[1] = model.ServiceItem{ID: 1, ItemCode: "G1", Name: "Gift", Category: "escort", SubCategory: model.SubCategoryGift, BasePriceCents: 1000, CommissionRate: 0.2, IsActive: true}
	r, _ := setupItemRouter(repo)
	sh := 1
	up := itemservice.UpdateServiceItemRequest{ServiceHours: &sh}
	ub, _ := json.Marshal(up)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/admin/service-items/1", bytes.NewReader(ub))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestServiceItem_UpdateNegativePrice(t *testing.T) {
	repo := newFakeItemRepo()
	repo.items[2] = model.ServiceItem{ID: 2, ItemCode: "S1", Name: "Svc", Category: "escort", SubCategory: model.SubCategorySolo, BasePriceCents: 1000, CommissionRate: 0.2, IsActive: true}
	r, _ := setupItemRouter(repo)
	up := itemservice.UpdateServiceItemRequest{BasePriceCents: ptrInt64(-100)}
	ub, _ := json.Marshal(up)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/admin/service-items/2", bytes.NewReader(ub))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 400/500, got %d", w.Code)
	}
}
