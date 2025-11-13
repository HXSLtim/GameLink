package admin

import (
    "context"
    "testing"

    "gamelink/internal/model"
    "gamelink/internal/repository"
    "gamelink/internal/repository/common"
)

type txPlayers struct{ gotID uint64 }
func (p *txPlayers) List(ctx context.Context) ([]model.Player, error) { _ = ctx; return nil, nil }
func (p *txPlayers) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) { _ = ctx; _ = page; _ = pageSize; return nil, 0, nil }
func (p *txPlayers) Get(ctx context.Context, id uint64) (*model.Player, error) { _ = ctx; p.gotID = id; return &model.Player{UserID: 1, Nickname: "p"}, nil }
func (p *txPlayers) Create(ctx context.Context, pl *model.Player) error { _ = ctx; pl.ID = 2; return nil }
func (p *txPlayers) Update(ctx context.Context, _ *model.Player) error { _ = ctx; return nil }
func (p *txPlayers) Delete(ctx context.Context, _ uint64) error { _ = ctx; return nil }
func (p *txPlayers) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) { _ = ctx; return &model.Player{UserID: userID, Nickname: "p"}, nil }

type txUsers struct{}
func (u *txUsers) List(ctx context.Context) ([]model.User, error) { _ = ctx; return nil, nil }
func (u *txUsers) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) { _ = ctx; _ = page; _ = pageSize; return nil, 0, nil }
func (u *txUsers) ListWithFilters(ctx context.Context, _ repository.UserListOptions) ([]model.User, int64, error) { _ = ctx; return nil, 0, nil }
func (u *txUsers) Get(ctx context.Context, _ uint64) (*model.User, error) { _ = ctx; return &model.User{Name: "u"}, nil }
func (u *txUsers) GetByPhone(ctx context.Context, _ string) (*model.User, error) { _ = ctx; return nil, repository.ErrNotFound }
func (u *txUsers) FindByEmail(ctx context.Context, _ string) (*model.User, error) { _ = ctx; return nil, repository.ErrNotFound }
func (u *txUsers) FindByPhone(ctx context.Context, _ string) (*model.User, error) { _ = ctx; return nil, repository.ErrNotFound }
func (u *txUsers) Create(ctx context.Context, user *model.User) error { _ = ctx; user.ID = 1; return nil }
func (u *txUsers) Update(ctx context.Context, _ *model.User) error { _ = ctx; return nil }
func (u *txUsers) Delete(ctx context.Context, _ uint64) error { _ = ctx; return nil }

type txTags struct{ last []string }
func (t *txTags) GetTags(ctx context.Context, _ uint64) ([]string, error) { _ = ctx; return nil, nil }
func (t *txTags) ReplaceTags(ctx context.Context, _ uint64, tags []string) error { _ = ctx; t.last = tags; return nil }

type txOpLogs struct{}
func (o *txOpLogs) Append(ctx context.Context, log *model.OperationLog) error { _ = ctx; _ = log; return nil }
func (o *txOpLogs) ListByEntity(ctx context.Context, entityType string, entityID uint64, opts repository.OperationLogListOptions) ([]model.OperationLog, int64, error) {
    _ = ctx; _ = entityType; _ = entityID; _ = opts; return nil, 0, nil
}

type txTagsErr struct{}
func (t *txTagsErr) GetTags(ctx context.Context, _ uint64) ([]string, error) { _ = ctx; return nil, nil }
func (t *txTagsErr) ReplaceTags(ctx context.Context, _ uint64, _ []string) error { _ = ctx; return repository.ErrNotFound }

type fakeTx struct{ repos common.Repos }
func (f *fakeTx) WithTx(ctx context.Context, fn func(r *common.Repos) error) error { _ = ctx; return fn(&f.repos) }

// --- tests ---

func TestUpdatePlayerSkillTags_NoTx(t *testing.T) {
    svc := NewAdminService(nil, nil, nil, nil, nil, nil, nil)
    if err := svc.UpdatePlayerSkillTags(context.Background(), 1, []string{"a", "b"}); err == nil {
        t.Fatal("expected error when TxManager is not configured")
    }
}

func TestUpdatePlayerSkillTags_WithTx(t *testing.T) {
    svc := NewAdminService(nil, nil, nil, nil, nil, nil, nil)
    svc.SetTxManager(&fakeTx{repos: common.Repos{Players: &txPlayers{}, Tags: &txTags{}, OpLogs: &txOpLogs{}}})
    if err := svc.UpdatePlayerSkillTags(context.Background(), 7, []string{"speed", "team"}); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}

func TestUpdatePlayerSkillTags_TxReplaceError(t *testing.T) {
    svc := NewAdminService(nil, nil, nil, nil, nil, nil, nil)
    svc.SetTxManager(&fakeTx{repos: common.Repos{Players: &txPlayers{}, Tags: &txTagsErr{}, OpLogs: &txOpLogs{}}})
    if err := svc.UpdatePlayerSkillTags(context.Background(), 7, []string{"speed"}); err == nil {
        t.Fatal("expected error from ReplaceTags")
    }
}

func TestRegisterUserAndPlayer_Validation(t *testing.T) {
    svc := NewAdminService(nil, nil, nil, nil, nil, nil, nil)
    svc.SetTxManager(&fakeTx{repos: common.Repos{Users: &txUsers{}, Players: &txPlayers{}, OpLogs: &txOpLogs{}}})
    // missing VerificationStatus
    u := CreateUserInput{Phone:"13800138000", Email:"u@example.com", Password:"Abc123", Name:"U", Role:model.RoleUser, Status:model.UserStatusActive}
    p := CreatePlayerInput{Nickname:"p", Bio:"b", HourlyRateCents:1000, MainGameID:1, VerificationStatus:""}
    if _, _, err := svc.RegisterUserAndPlayer(context.Background(), u, p); err == nil {
        t.Fatal("expected validation error for empty VerificationStatus")
    }
}

func TestRegisterUserAndPlayer_Success(t *testing.T) {
    svc := NewAdminService(nil, nil, nil, nil, nil, nil, nil)
    svc.SetTxManager(&fakeTx{repos: common.Repos{Users: &txUsers{}, Players: &txPlayers{}, OpLogs: &txOpLogs{}}})
    u := CreateUserInput{Phone:"13800138000", Email:"u@example.com", Password:"Abc123", Name:"U", Role:model.RoleUser, Status:model.UserStatusActive}
    p := CreatePlayerInput{Nickname:"p", Bio:"b", HourlyRateCents:1000, MainGameID:1, VerificationStatus:"pending"}
    user, player, err := svc.RegisterUserAndPlayer(context.Background(), u, p)
    if err != nil { t.Fatalf("unexpected err: %v", err) }
    if user == nil || player == nil { t.Fatal("expected user and player") }
}
