package admin

import (
    "bytes"
    "context"
    "testing"

    "gamelink/internal/cache"
    "gamelink/internal/logging"
    "gamelink/internal/model"
    "gamelink/internal/repository"
    "gamelink/internal/repository/common"
)

type usersRepo struct{ u *model.User }
func (r *usersRepo) List(ctx context.Context) ([]model.User, error) { _ = ctx; return nil, nil }
func (r *usersRepo) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) { _ = ctx; _ = page; _ = pageSize; return nil, 0, nil }
func (r *usersRepo) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) { _ = ctx; _ = opts; return nil, 0, nil }
func (r *usersRepo) Get(ctx context.Context, _ uint64) (*model.User, error) { _ = ctx; return r.u, nil }
func (r *usersRepo) GetByPhone(ctx context.Context, _ string) (*model.User, error) { _ = ctx; return nil, repository.ErrNotFound }
func (r *usersRepo) FindByEmail(ctx context.Context, _ string) (*model.User, error) { _ = ctx; return nil, repository.ErrNotFound }
func (r *usersRepo) FindByPhone(ctx context.Context, _ string) (*model.User, error) { _ = ctx; return nil, repository.ErrNotFound }
func (r *usersRepo) Create(ctx context.Context, _ *model.User) error { _ = ctx; return nil }
func (r *usersRepo) Update(ctx context.Context, u *model.User) error { _ = ctx; r.u = u; return nil }
func (r *usersRepo) Delete(ctx context.Context, _ uint64) error { _ = ctx; return nil }

type rolesRepo struct{ gotAssign []uint64 }
func (r *rolesRepo) List(ctx context.Context) ([]model.RoleModel, error) { _ = ctx; return nil, nil }
func (r *rolesRepo) ListPaged(ctx context.Context, page, pageSize int) ([]model.RoleModel, int64, error) { _ = ctx; _ = page; _ = pageSize; return nil, 0, nil }
func (r *rolesRepo) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string, isSystem *bool) ([]model.RoleModel, int64, error) { _ = ctx; _ = page; _ = pageSize; _ = keyword; _ = isSystem; return nil, 0, nil }
func (r *rolesRepo) ListWithPermissions(ctx context.Context) ([]model.RoleModel, error) { _ = ctx; return nil, nil }
func (r *rolesRepo) Get(ctx context.Context, _ uint64) (*model.RoleModel, error) { _ = ctx; return nil, repository.ErrNotFound }
func (r *rolesRepo) GetWithPermissions(ctx context.Context, _ uint64) (*model.RoleModel, error) { _ = ctx; return nil, repository.ErrNotFound }
func (r *rolesRepo) GetBySlug(ctx context.Context, _ string) (*model.RoleModel, error) { _ = ctx; return &model.RoleModel{Base: model.Base{ID: 10}}, nil }
func (r *rolesRepo) Create(ctx context.Context, _ *model.RoleModel) error { _ = ctx; return nil }
func (r *rolesRepo) Update(ctx context.Context, _ *model.RoleModel) error { _ = ctx; return nil }
func (r *rolesRepo) Delete(ctx context.Context, _ uint64) error { _ = ctx; return nil }
func (r *rolesRepo) AssignPermissions(ctx context.Context, _ uint64, _ []uint64) error { _ = ctx; return nil }
func (r *rolesRepo) AddPermissions(ctx context.Context, _ uint64, _ []uint64) error { _ = ctx; return nil }
func (r *rolesRepo) RemovePermissions(ctx context.Context, _ uint64, _ []uint64) error { _ = ctx; return nil }
func (r *rolesRepo) AssignToUser(ctx context.Context, _ uint64, roleIDs []uint64) error { _ = ctx; r.gotAssign = roleIDs; return nil }
func (r *rolesRepo) RemoveFromUser(ctx context.Context, _ uint64, _ []uint64) error { _ = ctx; return nil }
func (r *rolesRepo) ListByUserID(ctx context.Context, _ uint64) ([]model.RoleModel, error) { _ = ctx; return nil, nil }
func (r *rolesRepo) CheckUserHasRole(ctx context.Context, _ uint64, _ string) (bool, error) { _ = ctx; return false, nil }

type opLogsPager struct{ calls int }
func (o *opLogsPager) Append(ctx context.Context, _ *model.OperationLog) error { _ = ctx; return nil }
func (o *opLogsPager) ListByEntity(ctx context.Context, entityType string, entityID uint64, opts repository.OperationLogListOptions) ([]model.OperationLog, int64, error) {
    _ = ctx; _ = entityType; _ = entityID; _ = opts
    o.calls++
    if opts.Page == 1 { return []model.OperationLog{{Base: model.Base{ID:1}}}, 201, nil }
    return []model.OperationLog{}, 201, nil
}

type txPager struct{ repos common.Repos }
func (t *txPager) WithTx(ctx context.Context, fn func(r *common.Repos) error) error { _ = ctx; return fn(&t.repos) }

func TestUpdateUserStatus_And_Role(t *testing.T) {
    u := &usersRepo{u: &model.User{Base: model.Base{ID:1}, Name:"u", Role:model.RoleUser, Status:model.UserStatusActive}}
    r := &rolesRepo{}
    svc := NewAdminService(nil, u, nil, nil, nil, r, cache.NewMemory())
    out1, err1 := svc.UpdateUserStatus(context.Background(), 1, model.UserStatusSuspended)
    if err1 != nil || out1.Status != model.UserStatusSuspended { t.Fatalf("status update failed: %v", err1) }
    out2, err2 := svc.UpdateUserRole(context.Background(), 1, model.RoleAdmin)
    if err2 != nil || out2.Role != model.RoleAdmin { t.Fatalf("role update failed: %v", err2) }
    if len(r.gotAssign) != 1 || r.gotAssign[0] != 10 { t.Fatalf("expected assign role id 10, got %v", r.gotAssign) }
}

func TestCollectOperationLogs_AggregatesPages(t *testing.T) {
    svc := NewAdminService(nil, nil, nil, nil, nil, nil, cache.NewMemory())
    svc.SetTxManager(&txPager{repos: common.Repos{OpLogs: &opLogsPager{}}})
    items, err := svc.collectOperationLogs(context.Background(), "order", 1)
    if err != nil || len(items) < 1 { t.Fatalf("unexpected: err=%v items=%d", err, len(items)) }
}

type recordOpLogs struct{ last *model.OperationLog }
func (r *recordOpLogs) Append(ctx context.Context, log *model.OperationLog) error { _ = ctx; r.last = log; return nil }
func (r *recordOpLogs) ListByEntity(context.Context, string, uint64, repository.OperationLogListOptions) ([]model.OperationLog, int64, error) { return nil, 0, nil }

type paymentsRepoPending struct{ p model.Payment }
func (p paymentsRepoPending) Create(context.Context, *model.Payment) error { return nil }
func (p paymentsRepoPending) List(context.Context, repository.PaymentListOptions) ([]model.Payment, int64, error) { return nil, 0, nil }
func (p paymentsRepoPending) Get(context.Context, uint64) (*model.Payment, error) { return &p.p, nil }
func (p paymentsRepoPending) Update(context.Context, *model.Payment) error { return nil }
func (p paymentsRepoPending) Delete(context.Context, uint64) error { return nil }

type txRecorder struct{ repos common.Repos }
func (t *txRecorder) WithTx(ctx context.Context, fn func(r *common.Repos) error) error { return fn(&t.repos) }

func TestCapturePayment_AppendsActorAndTradeNo(t *testing.T) {
    rec := &recordOpLogs{}
    svc := NewAdminService(nil, nil, nil, nil, paymentsRepoPending{p: model.Payment{Base: model.Base{ID: 101}, Status: model.PaymentStatusPending}}, nil, cache.NewMemory())
    svc.SetTxManager(&txRecorder{repos: common.Repos{OpLogs: rec}})
    ctx := logging.WithRequestID(context.Background(), "req-1")
    ctx = logging.WithActorUserID(ctx, 999)
    _, err := svc.CapturePayment(ctx, 101, UpdatePaymentCapture("trade-abc"))
    if err != nil { t.Fatalf("%v", err) }
    if rec.last == nil || rec.last.ActorUserID == nil { t.Fatal("expected actor user id in op log") }
    if !bytes.Contains(rec.last.MetadataJSON, []byte("trade_no")) { t.Fatal("expected trade_no in metadata") }
}

func UpdatePaymentCapture(trade string) CapturePaymentInput { return CapturePaymentInput{ProviderTradeNo: trade} }

func TestUpdatePayment_AuditNoRequestID(t *testing.T) {
    rec := &recordOpLogs{}
    svc := NewAdminService(nil, nil, nil, nil, paymentsRepoGet{p: model.Payment{Base: model.Base{ID: 202}, Status: model.PaymentStatusPending}}, nil, cache.NewMemory())
    svc.SetTxManager(&txRecorder{repos: common.Repos{OpLogs: rec}})
    _, err := svc.UpdatePayment(context.Background(), 202, UpdatePaymentInput{Status: model.PaymentStatusFailed})
    if err != nil { t.Fatalf("%v", err) }
    if rec.last == nil || rec.last.ActorUserID != nil { t.Fatal("expected op log without actor user id") }
    if !bytes.Contains(rec.last.MetadataJSON, []byte("status")) { t.Fatal("expected status in metadata") }
}

type reviewsRepoTxErr struct{}
func (reviewsRepoTxErr) List(context.Context, repository.ReviewListOptions) ([]model.Review, int64, error) { return nil, 0, nil }
func (reviewsRepoTxErr) Get(context.Context, uint64) (*model.Review, error) { return &model.Review{Base: model.Base{ID:1}}, nil }
func (reviewsRepoTxErr) Create(context.Context, *model.Review) error { return nil }
func (reviewsRepoTxErr) Update(context.Context, *model.Review) error { return nil }
func (reviewsRepoTxErr) Delete(context.Context, uint64) error { return nil }

type txRejectReviews struct{}
func (txRejectReviews) WithTx(context.Context, func(r *common.Repos) error) error { return repository.ErrNotFound }

func TestAdminService_CreateReview_TxRollback(t *testing.T) {
    svc := NewAdminService(nil, nil, nil, nil, nil, nil, cache.NewMemory())
    svc.SetTxManager(txRejectReviews{})
    _, err := svc.CreateReview(context.Background(), model.Review{OrderID:1, UserID:1, PlayerID:1, Score:5})
    if err == nil { t.Fatal("expected tx error") }
}

type paymentsRepo struct{}
func (p paymentsRepo) Create(context.Context, *model.Payment) error { return nil }
func (p paymentsRepo) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
    _ = ctx
    if opts.Page == 1 { return []model.Payment{{Base: model.Base{ID:1}}}, 201, nil }
    if opts.Page == 2 { return []model.Payment{{Base: model.Base{ID:2}}}, 201, nil }
    return []model.Payment{}, 201, nil
}
func (p paymentsRepo) Get(context.Context, uint64) (*model.Payment, error) { return nil, repository.ErrNotFound }
func (p paymentsRepo) Update(context.Context, *model.Payment) error { return nil }
func (p paymentsRepo) Delete(context.Context, uint64) error { return nil }

func TestListPaymentsByOrder_AggregatesPages(t *testing.T) {
    svc := NewAdminService(nil, nil, nil, nil, paymentsRepo{}, nil, cache.NewMemory())
    items, err := svc.listPaymentsByOrder(context.Background(), 1)
    if err != nil || len(items) != 2 { t.Fatalf("unexpected: err=%v items=%d", err, len(items)) }
}

type usersRepoCount struct{ calls int }
func (u *usersRepoCount) List(context.Context) ([]model.User, error) { return nil, nil }
func (u *usersRepoCount) ListPaged(context.Context, int, int) ([]model.User, int64, error) { return nil, 0, nil }
func (u *usersRepoCount) ListWithFilters(context.Context, repository.UserListOptions) ([]model.User, int64, error) { return nil, 0, nil }
func (u *usersRepoCount) Get(context.Context, uint64) (*model.User, error) { u.calls++; return nil, repository.ErrNotFound }
func (u *usersRepoCount) GetByPhone(context.Context, string) (*model.User, error) { return nil, repository.ErrNotFound }
func (u *usersRepoCount) FindByEmail(context.Context, string) (*model.User, error) { return nil, repository.ErrNotFound }
func (u *usersRepoCount) FindByPhone(context.Context, string) (*model.User, error) { return nil, repository.ErrNotFound }
func (u *usersRepoCount) Create(context.Context, *model.User) error { return nil }
func (u *usersRepoCount) Update(context.Context, *model.User) error { return nil }
func (u *usersRepoCount) Delete(context.Context, uint64) error { return nil }

func TestResolveUser_CachesNotFound(t *testing.T) {
    u := &usersRepoCount{}
    svc := NewAdminService(nil, u, nil, nil, nil, nil, cache.NewMemory())
    m := map[uint64]*model.User{}
    r1 := svc.resolveUser(context.Background(), m, 1)
    r2 := svc.resolveUser(context.Background(), m, 1)
    if r1 != nil || r2 != nil || u.calls != 1 { t.Fatalf("expected single repo call and nil user, got calls=%d", u.calls) }
}

type opLogsErr struct{}
func (o *opLogsErr) Append(context.Context, *model.OperationLog) error { return nil }
func (o *opLogsErr) ListByEntity(context.Context, string, uint64, repository.OperationLogListOptions) ([]model.OperationLog, int64, error) {
    return nil, 0, repository.ErrNotFound
}

type txErr struct{ repos common.Repos }
func (t *txErr) WithTx(ctx context.Context, fn func(r *common.Repos) error) error { return fn(&t.repos) }

func TestListOperationLogs_ErrorPath(t *testing.T) {
    svc := NewAdminService(nil, nil, nil, nil, nil, nil, cache.NewMemory())
    svc.SetTxManager(&txErr{repos: common.Repos{OpLogs: &opLogsErr{}}})
    _, _, err := svc.ListOperationLogs(context.Background(), "order", 1, repository.OperationLogListOptions{Page:1, PageSize:10})
    if err == nil { t.Fatal("expected error from OpLogs.ListByEntity") }
}

func TestListOperationLogs_NoTx(t *testing.T) {
    svc := NewAdminService(nil, nil, nil, nil, nil, nil, cache.NewMemory())
    _, _, err := svc.ListOperationLogs(context.Background(), "order", 1, repository.OperationLogListOptions{Page:1, PageSize:10})
    if err == nil { t.Fatal("expected transaction manager not configured") }
}

type paymentsRepoGet struct{ p model.Payment }
func (p paymentsRepoGet) Create(context.Context, *model.Payment) error { return nil }
func (p paymentsRepoGet) List(context.Context, repository.PaymentListOptions) ([]model.Payment, int64, error) { return nil, 0, nil }
func (p paymentsRepoGet) Get(context.Context, uint64) (*model.Payment, error) { return &p.p, nil }
func (p paymentsRepoGet) Update(context.Context, *model.Payment) error { return nil }
func (p paymentsRepoGet) Delete(context.Context, uint64) error { return nil }

func TestUpdatePayment_InvalidStatusAndTransition(t *testing.T) {
    svc := NewAdminService(nil, nil, nil, nil, paymentsRepoGet{p: model.Payment{Status: model.PaymentStatusPending}}, nil, cache.NewMemory())
    _, err := svc.UpdatePayment(context.Background(), 1, UpdatePaymentInput{Status: model.PaymentStatus("bad")})
    if err == nil { t.Fatal("expected validation error for invalid status") }

    svc2 := NewAdminService(nil, nil, nil, nil, paymentsRepoGet{p: model.Payment{Status: model.PaymentStatusFailed}}, nil, cache.NewMemory())
    _, err2 := svc2.UpdatePayment(context.Background(), 1, UpdatePaymentInput{Status: model.PaymentStatusPaid})
    if err2 == nil { t.Fatal("expected validation error for invalid transition") }
}

type playersRepoAssign struct{}
func (playersRepoAssign) List(context.Context) ([]model.Player, error) { return nil, nil }
func (playersRepoAssign) ListPaged(context.Context, int, int) ([]model.Player, int64, error) { return nil, 0, nil }
func (playersRepoAssign) Get(context.Context, uint64) (*model.Player, error) { return &model.Player{Nickname:"p"}, nil }
func (playersRepoAssign) Create(context.Context, *model.Player) error { return nil }
func (playersRepoAssign) Update(context.Context, *model.Player) error { return nil }
func (playersRepoAssign) Delete(context.Context, uint64) error { return nil }

type ordersRepoAssign struct{ status model.OrderStatus }
func (o ordersRepoAssign) Create(context.Context, *model.Order) error { return nil }
func (o ordersRepoAssign) List(context.Context, repository.OrderListOptions) ([]model.Order, int64, error) { return nil, 0, nil }
func (o ordersRepoAssign) Get(context.Context, uint64) (*model.Order, error) { return &model.Order{Status: o.status}, nil }
func (o ordersRepoAssign) Update(context.Context, *model.Order) error { return nil }
func (o ordersRepoAssign) Delete(context.Context, uint64) error { return nil }

func TestAssignOrder_Validations(t *testing.T) {
    svc := NewAdminService(nil, nil, playersRepoAssign{}, ordersRepoAssign{status: model.OrderStatusCompleted}, nil, nil, cache.NewMemory())
    if _, err := svc.AssignOrder(context.Background(), 1, 0); err == nil { t.Fatal("expected validation error for zero playerID") }
    if _, err := svc.AssignOrder(context.Background(), 1, 2); err == nil { t.Fatal("expected validation error for completed order") }
}
