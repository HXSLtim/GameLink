package admin

import (
    "context"
    "testing"

    "gamelink/internal/cache"
    "gamelink/internal/model"
    "gamelink/internal/repository"
    "gamelink/internal/repository/common"
)

type reviewsPagerRepo struct{}
func (reviewsPagerRepo) List(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
    _ = ctx
    if opts.Page == 1 { return []model.Review{{Base: model.Base{ID: 1}}}, 201, nil }
    if opts.Page == 2 { return []model.Review{{Base: model.Base{ID: 2}}}, 201, nil }
    return []model.Review{}, 201, nil
}
func (reviewsPagerRepo) Get(context.Context, uint64) (*model.Review, error) { return nil, repository.ErrNotFound }
func (reviewsPagerRepo) Create(context.Context, *model.Review) error { return nil }
func (reviewsPagerRepo) Update(context.Context, *model.Review) error { return nil }
func (reviewsPagerRepo) Delete(context.Context, uint64) error { return nil }

type txReviewsPager struct{ repos common.Repos }
func (t *txReviewsPager) WithTx(ctx context.Context, fn func(r *common.Repos) error) error { return fn(&t.repos) }

func TestGetOrderReviews_AggregatesPages(t *testing.T) {
    svc := NewAdminService(nil, nil, nil, nil, nil, nil, cache.NewMemory())
    svc.SetTxManager(&txReviewsPager{repos: common.Repos{Reviews: reviewsPagerRepo{}}})
    items, err := svc.GetOrderReviews(context.Background(), 1)
    if err != nil { t.Fatalf("%v", err) }
    if len(items) != 2 { t.Fatalf("expected 2 reviews, got %d", len(items)) }
}

