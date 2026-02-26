package item

import (
	"context"
	"errors"
	"testing"

	"gamelink/internal/model"
	"gamelink/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeServiceItemRepository struct {
	items         map[uint64]*model.ServiceItem
	orderRefs     map[uint64]int64
	orderRefErrs  map[uint64]error
	deleteErrs    map[uint64]error
	deletedItemID []uint64
}

func (r *fakeServiceItemRepository) Create(_ context.Context, item *model.ServiceItem) error {
	if r.items == nil {
		r.items = make(map[uint64]*model.ServiceItem)
	}
	r.items[item.ID] = item
	return nil
}

func (r *fakeServiceItemRepository) Get(_ context.Context, id uint64) (*model.ServiceItem, error) {
	item, ok := r.items[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return item, nil
}

func (r *fakeServiceItemRepository) GetByCode(_ context.Context, itemCode string) (*model.ServiceItem, error) {
	for _, item := range r.items {
		if item.ItemCode == itemCode {
			return item, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (r *fakeServiceItemRepository) List(_ context.Context, _ repository.ServiceItemListOptions) ([]model.ServiceItem, int64, error) {
	return nil, 0, nil
}

func (r *fakeServiceItemRepository) Update(_ context.Context, _ *model.ServiceItem) error {
	return nil
}

func (r *fakeServiceItemRepository) Delete(_ context.Context, id uint64) error {
	r.deletedItemID = append(r.deletedItemID, id)
	if err, ok := r.deleteErrs[id]; ok {
		return err
	}
	delete(r.items, id)
	return nil
}

func (r *fakeServiceItemRepository) BatchDelete(ctx context.Context, ids []uint64) (int64, error) {
	var successCount int64
	for _, id := range ids {
		if err := r.Delete(ctx, id); err != nil {
			return successCount, err
		}
		successCount++
	}
	return successCount, nil
}

func (r *fakeServiceItemRepository) BatchUpdateStatus(_ context.Context, _ []uint64, _ bool) error {
	return nil
}

func (r *fakeServiceItemRepository) BatchUpdatePrice(_ context.Context, _ []uint64, _ int64) error {
	return nil
}

func (r *fakeServiceItemRepository) BatchUpdateCommission(_ context.Context, _ []uint64, _ float64) error {
	return nil
}

func (r *fakeServiceItemRepository) GetGifts(_ context.Context, _, _ int) ([]model.ServiceItem, int64, error) {
	return nil, 0, nil
}

func (r *fakeServiceItemRepository) GetGameServices(_ context.Context, _ uint64, _ *model.ServiceItemSubCategory) ([]model.ServiceItem, error) {
	return nil, nil
}

// CountOrderReferences is intentionally defined for batch-delete reference checks.
func (r *fakeServiceItemRepository) CountOrderReferences(_ context.Context, itemID uint64) (int64, error) {
	if err, ok := r.orderRefErrs[itemID]; ok {
		return 0, err
	}
	if refs, ok := r.orderRefs[itemID]; ok {
		return refs, nil
	}
	return 0, nil
}

func TestBatchDeleteItems_OrderReferenced_ShouldContinueAndReportClearly(t *testing.T) {
	repo := &fakeServiceItemRepository{
		items: map[uint64]*model.ServiceItem{
			1: {ID: 1, Name: "item-1"},
			2: {ID: 2, Name: "item-2"},
			3: {ID: 3, Name: "item-3"},
		},
		orderRefs: map[uint64]int64{
			2: 2,
		},
	}
	svc := NewServiceItemService(repo, nil, nil)

	resp, err := svc.BatchDeleteItems(context.Background(), BatchDeleteItemsRequest{ItemIDs: []uint64{1, 2, 3, 999}})
	require.NoError(t, err)

	assert.Equal(t, 2, resp.SuccessCount)
	assert.Equal(t, 2, resp.FailedCount)
	assert.Equal(t, 4, resp.TotalCount)
	assert.ElementsMatch(t, []uint64{1, 3}, resp.SuccessItems)

	referenced := findBatchDeleteFailure(resp.FailedItems, 2)
	require.NotNil(t, referenced)
	assert.Contains(t, referenced.Message, "referenced")
	assert.Contains(t, referenced.Message, "order")

	notFound := findBatchDeleteFailure(resp.FailedItems, 999)
	require.NotNil(t, notFound)
	assert.Contains(t, notFound.Message, "not found")

	assert.NotContains(t, repo.deletedItemID, uint64(2), "referenced item should not be deleted")
}

func TestBatchDeleteItems_CheckReferenceError_ShouldFailOnlyCurrentItem(t *testing.T) {
	errCheck := errors.New("check failed")
	repo := &fakeServiceItemRepository{
		items: map[uint64]*model.ServiceItem{
			11: {ID: 11, Name: "item-11"},
			12: {ID: 12, Name: "item-12"},
		},
		orderRefErrs: map[uint64]error{
			11: errCheck,
		},
	}

	svc := NewServiceItemService(repo, nil, nil)

	resp, err := svc.BatchDeleteItems(context.Background(), BatchDeleteItemsRequest{ItemIDs: []uint64{11, 12}})
	require.NoError(t, err)

	assert.Equal(t, 1, resp.SuccessCount)
	assert.Equal(t, 1, resp.FailedCount)
	assert.ElementsMatch(t, []uint64{12}, resp.SuccessItems)

	failed := findBatchDeleteFailure(resp.FailedItems, 11)
	require.NotNil(t, failed)
	assert.Contains(t, failed.Message, "check")
}

func findBatchDeleteFailure(items []BatchOperationErrorItem, id uint64) *BatchOperationErrorItem {
	for i := range items {
		if items[i].ID == id {
			return &items[i]
		}
	}
	return nil
}
