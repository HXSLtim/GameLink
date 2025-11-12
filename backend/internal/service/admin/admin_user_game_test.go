package admin

import (
    "context"
    "testing"

    "gamelink/internal/cache"
    "gamelink/internal/model"
    "gamelink/internal/repository"
)

type gamesStub struct{ getErr error; listItems []model.Game; listTotal int64 }
func (g *gamesStub) List(ctx context.Context) ([]model.Game, error) { _=ctx; return []model.Game{{Name:"g"}}, nil }
func (g *gamesStub) ListPaged(ctx context.Context, page, pageSize int) ([]model.Game, int64, error) { _=ctx; _=page; _=pageSize; return g.listItems, g.listTotal, nil }
func (g *gamesStub) Get(ctx context.Context, id uint64) (*model.Game, error) { _=ctx; _=id; if g.getErr!=nil { return nil, g.getErr }; return &model.Game{Base: model.Base{ID:1}, Name:"g"}, nil }
func (g *gamesStub) Create(ctx context.Context, game *model.Game) error { _=ctx; game.ID=1; return nil }
func (g *gamesStub) Update(ctx context.Context, game *model.Game) error { _=ctx; _=game; return nil }
func (g *gamesStub) Delete(ctx context.Context, id uint64) error { _=ctx; _=id; return nil }

type usersStub struct{ user *model.User; delErr error; listTotal int64 }
func (u *usersStub) List(ctx context.Context) ([]model.User, error) { _=ctx; return []model.User{{Name:"u"}}, nil }
func (u *usersStub) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) { _=ctx; _=page; _=pageSize; return []model.User{{Name:"u"}}, u.listTotal, nil }
func (u *usersStub) ListWithFilters(ctx context.Context, _ repository.UserListOptions) ([]model.User, int64, error) { _=ctx; return []model.User{{Name:"u"}}, u.listTotal, nil }
func (u *usersStub) Get(ctx context.Context, _ uint64) (*model.User, error) { _=ctx; return u.user, nil }
func (u *usersStub) GetByPhone(context.Context, string) (*model.User, error) { return nil, repository.ErrNotFound }
func (u *usersStub) FindByEmail(context.Context, string) (*model.User, error) { return nil, repository.ErrNotFound }
func (u *usersStub) FindByPhone(context.Context, string) (*model.User, error) { return nil, repository.ErrNotFound }
func (u *usersStub) Create(ctx context.Context, user *model.User) error { _=ctx; user.ID=1; return nil }
func (u *usersStub) Update(ctx context.Context, user *model.User) error { _=ctx; u.user = user; return nil }
func (u *usersStub) Delete(ctx context.Context, _ uint64) error { _=ctx; return u.delErr }

type rolesStub struct{}
func (rolesStub) List(ctx context.Context) ([]model.RoleModel, error) { _=ctx; return nil, nil }
func (rolesStub) ListPaged(ctx context.Context, page, pageSize int) ([]model.RoleModel, int64, error) { _=ctx; _=page; _=pageSize; return nil, 0, nil }
func (rolesStub) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string, isSystem *bool) ([]model.RoleModel, int64, error) { _=ctx; _=page; _=pageSize; _=keyword; _=isSystem; return nil, 0, nil }
func (rolesStub) ListWithPermissions(ctx context.Context) ([]model.RoleModel, error) { _=ctx; return nil, nil }
func (rolesStub) Get(ctx context.Context, id uint64) (*model.RoleModel, error) { _=ctx; _=id; return nil, repository.ErrNotFound }
func (rolesStub) GetWithPermissions(ctx context.Context, id uint64) (*model.RoleModel, error) { _=ctx; _=id; return nil, repository.ErrNotFound }
func (rolesStub) GetBySlug(ctx context.Context, slug string) (*model.RoleModel, error) { _=ctx; _=slug; return nil, repository.ErrNotFound }
func (rolesStub) Create(ctx context.Context, role *model.RoleModel) error { _=ctx; _=role; return nil }
func (rolesStub) Update(ctx context.Context, role *model.RoleModel) error { _=ctx; _=role; return nil }
func (rolesStub) Delete(ctx context.Context, id uint64) error { _=ctx; _=id; return nil }
func (rolesStub) AssignPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error { _=ctx; _=roleID; _=permissionIDs; return nil }
func (rolesStub) AddPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error { _=ctx; _=roleID; _=permissionIDs; return nil }
func (rolesStub) RemovePermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error { _=ctx; _=roleID; _=permissionIDs; return nil }
func (rolesStub) AssignToUser(ctx context.Context, userID uint64, roleIDs []uint64) error { _=ctx; _=userID; _=roleIDs; return nil }
func (rolesStub) RemoveFromUser(ctx context.Context, userID uint64, roleIDs []uint64) error { _=ctx; _=userID; _=roleIDs; return nil }
func (rolesStub) ListByUserID(ctx context.Context, userID uint64) ([]model.RoleModel, error) { _=ctx; _=userID; return nil, nil }
func (rolesStub) CheckUserHasRole(ctx context.Context, userID uint64, slug string) (bool, error) { _=ctx; _=userID; _=slug; return false, nil }

// --- tests ---

func TestCreateGame_ValidationError(t *testing.T) {
    svc := NewAdminService(&gamesStub{}, nil, nil, nil, nil, &rolesStub{}, cache.NewMemory())
    _, err := svc.CreateGame(context.Background(), CreateGameInput{Key:"", Name:""})
    if err == nil { t.Fatal("expected validation error") }
}

func TestUpdateGame_NotFound(t *testing.T) {
    g := &gamesStub{getErr: repository.ErrNotFound}
    svc := NewAdminService(g, nil, nil, nil, nil, &rolesStub{}, cache.NewMemory())
    _, err := svc.UpdateGame(context.Background(), 99, UpdateGameInput{Key:"k", Name:"n"})
    if err == nil { t.Fatal("expected error") }
}

func TestListGamesPaged_NormalizeAndPagination(t *testing.T) {
    g := &gamesStub{listItems: []model.Game{{Name:"g"}}, listTotal: 1}
    svc := NewAdminService(g, nil, nil, nil, nil, &rolesStub{}, cache.NewMemory())
    items, p, err := svc.ListGamesPaged(context.Background(), -1, -1)
    if err != nil || p == nil { t.Fatalf("unexpected: err=%v p=%v", err, p) }
    if len(items) != 1 || p.Total != 1 { t.Fatalf("expected one item and total=1, got %d total=%d", len(items), p.Total) }
}

func TestUpdateUser_PasswordWhitespaceValidationError(t *testing.T) {
    u := &usersStub{user: &model.User{Base: model.Base{ID:1}, Name:"n", Role:model.RoleUser, Status:model.UserStatusActive, PasswordHash:"old"}}
    svc := NewAdminService(nil, u, nil, nil, nil, &rolesStub{}, cache.NewMemory())
    pw := "   "
    _, err := svc.UpdateUser(context.Background(), 1, UpdateUserInput{Name:"n", Role:model.RoleUser, Status:model.UserStatusActive, Password:&pw})
    if err == nil { t.Fatal("expected validation error for whitespace password") }
}

func TestUpdateUser_ValidPasswordChangesHash(t *testing.T) {
    u := &usersStub{user: &model.User{Base: model.Base{ID:1}, Name:"n", Role:model.RoleUser, Status:model.UserStatusActive, PasswordHash:"old"}}
    svc := NewAdminService(nil, u, nil, nil, nil, &rolesStub{}, cache.NewMemory())
    pw := "Abc123"
    out, err := svc.UpdateUser(context.Background(), 1, UpdateUserInput{Name:"n", Role:model.RoleUser, Status:model.UserStatusActive, Password:&pw})
    if err != nil { t.Fatalf("unexpected %v", err) }
    if out.PasswordHash == "old" { t.Fatal("expected password hash updated") }
}

func TestDeleteUser_MapUserError(t *testing.T) {
    u := &usersStub{delErr: repository.ErrNotFound}
    svc := NewAdminService(nil, u, nil, nil, nil, &rolesStub{}, cache.NewMemory())
    if err := svc.DeleteUser(context.Background(), 123); err != ErrUserNotFound { t.Fatalf("expected ErrUserNotFound, got %v", err) }
}

func TestListUsersWithOptions_Normalize(t *testing.T) {
    u := &usersStub{listTotal: 2}
    svc := NewAdminService(nil, u, nil, nil, nil, &rolesStub{}, cache.NewMemory())
    items, p, err := svc.ListUsersWithOptions(context.Background(), repository.UserListOptions{Page:-1, PageSize:-1})
    if err != nil || p == nil { t.Fatalf("unexpected: err=%v p=%v", err, p) }
    if len(items) < 1 || p.Total != 2 { t.Fatalf("expected total=2, got %d", p.Total) }
}
