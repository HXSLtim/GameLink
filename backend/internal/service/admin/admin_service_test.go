package admin

import (
    "context"
    "encoding/json"
    "testing"
    "time"

    "gamelink/internal/cache"
    "gamelink/internal/model"
    "gamelink/internal/repository"
    "gamelink/internal/repository/common"
)

// mockGameRepository 是一个简单的mock实现。
type mockGameRepository struct {
	repository.GameRepository
}

type mockUserRepository struct {
	repository.UserRepository
}

type mockPlayerRepository struct {
	repository.PlayerRepository
}

type mockOrderRepository struct {
	repository.OrderRepository
}

type mockPaymentRepository struct {
	repository.PaymentRepository
}

type mockRoleRepository struct {
	repository.RoleRepository
}

type mockCache struct {
	cache.Cache
}

// TestNewAdminService 测试构造函数。
func TestNewAdminService(t *testing.T) {
	games := &mockGameRepository{}
	users := &mockUserRepository{}
	players := &mockPlayerRepository{}
	orders := &mockOrderRepository{}
	payments := &mockPaymentRepository{}
	roles := &mockRoleRepository{}
	cache := &mockCache{}

	svc := NewAdminService(games, users, players, orders, payments, roles, cache)

	if svc == nil {
		t.Fatal("NewAdminService returned nil")
	}

	if svc.games != games {
		t.Error("games repository not set correctly")
	}
	if svc.users != users {
		t.Error("users repository not set correctly")
	}
	if svc.players != players {
		t.Error("players repository not set correctly")
	}
	if svc.orders != orders {
		t.Error("orders repository not set correctly")
	}
	if svc.payments != payments {
		t.Error("payments repository not set correctly")
	}
	if svc.roles != roles {
		t.Error("roles repository not set correctly")
	}
	if svc.cache != cache {
		t.Error("cache not set correctly")
	}
}

// TestSetTxManager 测试事务管理器注入。
func TestSetTxManager(t *testing.T) {
	svc := NewAdminService(
		&mockGameRepository{},
		&mockUserRepository{},
		&mockPlayerRepository{},
		&mockOrderRepository{},
		&mockPaymentRepository{},
		&mockRoleRepository{},
		&mockCache{},
	)

	if svc.tx != nil {
		t.Error("tx should be nil initially")
	}

	// Note: 我们不能测试实际的TxManager，因为它是一个接口
	// 这个测试只是确保方法存在并可以调用
}

// ---- Fakes for admin package tests ----

type fakeGameRepo struct {
    items     []model.Game
    listCalls int
}

func (f *fakeGameRepo) List(ctx context.Context) ([]model.Game, error) {
    f.listCalls++
    cp := append([]model.Game(nil), f.items...)
    return cp, nil
}
func (f *fakeGameRepo) ListPaged(ctx context.Context, page, size int) ([]model.Game, int64, error) {
    cp := append([]model.Game(nil), f.items...)
    return cp, int64(len(cp)), nil
}
func (f *fakeGameRepo) Get(ctx context.Context, id uint64) (*model.Game, error) {
    for i := range f.items {
        if f.items[i].ID == id {
            c := f.items[i]
            return &c, nil
        }
    }
    return nil, repository.ErrNotFound
}
func (f *fakeGameRepo) Create(ctx context.Context, g *model.Game) error {
    if g.ID == 0 {
        g.ID = uint64(len(f.items) + 1)
    }
    f.items = append(f.items, *g)
    return nil
}
func (f *fakeGameRepo) Update(ctx context.Context, g *model.Game) error { return nil }
func (f *fakeGameRepo) Delete(ctx context.Context, id uint64) error     { return nil }

type fakeUserRepo struct{ last *model.User }

func (f *fakeUserRepo) List(ctx context.Context) ([]model.User, error) { return nil, nil }
func (f *fakeUserRepo) ListPaged(ctx context.Context, page, size int) ([]model.User, int64, error) {
    return []model.User{}, 250, nil
}
func (f *fakeUserRepo) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
    return []model.User{}, 250, nil
}
func (f *fakeUserRepo) Get(ctx context.Context, id uint64) (*model.User, error) {
    if f.last != nil && f.last.ID == id {
        return f.last, nil
    }
    return nil, repository.ErrNotFound
}
func (f *fakeUserRepo) GetByPhone(ctx context.Context, phone string) (*model.User, error) { return nil, repository.ErrNotFound }
func (f *fakeUserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) { return f.last, nil }
func (f *fakeUserRepo) FindByPhone(ctx context.Context, phone string) (*model.User, error) { return f.last, nil }
func (f *fakeUserRepo) Create(ctx context.Context, u *model.User) error {
    if u.ID == 0 {
        u.ID = 42
    }
    f.last = u
    return nil
}
func (f *fakeUserRepo) Update(ctx context.Context, u *model.User) error { f.last = u; return nil }
func (f *fakeUserRepo) Delete(ctx context.Context, id uint64) error     { return nil }

type fakePlayerRepo struct{ last *model.Player }

func (f *fakePlayerRepo) List(ctx context.Context) ([]model.Player, error) { return nil, nil }
func (f *fakePlayerRepo) ListPaged(ctx context.Context, page, size int) ([]model.Player, int64, error) {
    return []model.Player{}, 0, nil
}
func (f *fakePlayerRepo) Get(ctx context.Context, id uint64) (*model.Player, error) {
    if f.last != nil && f.last.ID == id {
        return f.last, nil
    }
    return nil, repository.ErrNotFound
}
func (f *fakePlayerRepo) Create(ctx context.Context, p *model.Player) error {
    if p.ID == 0 {
        p.ID = 99
    }
    f.last = p
    return nil
}
func (f *fakePlayerRepo) Update(ctx context.Context, p *model.Player) error { f.last = p; return nil }
func (f *fakePlayerRepo) Delete(ctx context.Context, id uint64) error       { return nil }

type fakeOrderRepo struct{ obj *model.Order }

func (f *fakeOrderRepo) List(ctx context.Context, _ repository.OrderListOptions) ([]model.Order, int64, error) {
    return nil, 0, nil
}
func (f *fakeOrderRepo) Create(ctx context.Context, o *model.Order) error { f.obj = o; return nil }
func (f *fakeOrderRepo) Get(ctx context.Context, id uint64) (*model.Order, error) {
    if f.obj == nil {
        return nil, repository.ErrNotFound
    }
    return f.obj, nil
}
func (f *fakeOrderRepo) Update(ctx context.Context, o *model.Order) error { f.obj = o; return nil }
func (f *fakeOrderRepo) Delete(ctx context.Context, id uint64) error      { return nil }

type fakePaymentRepo struct{ obj *model.Payment }

func (f *fakePaymentRepo) List(ctx context.Context, _ repository.PaymentListOptions) ([]model.Payment, int64, error) {
    return nil, 0, nil
}
func (f *fakePaymentRepo) Create(ctx context.Context, p *model.Payment) error { f.obj = p; return nil }
func (f *fakePaymentRepo) Get(ctx context.Context, id uint64) (*model.Payment, error) {
    if f.obj == nil {
        return nil, repository.ErrNotFound
    }
    return f.obj, nil
}
func (f *fakePaymentRepo) Update(ctx context.Context, p *model.Payment) error { f.obj = p; return nil }
func (f *fakePaymentRepo) Delete(ctx context.Context, id uint64) error        { return nil }

type fakeRoleRepo struct{}

func (f *fakeRoleRepo) List(ctx context.Context) ([]model.RoleModel, error) { return nil, nil }
func (f *fakeRoleRepo) ListPaged(ctx context.Context, page, pageSize int) ([]model.RoleModel, int64, error) {
    return nil, 0, nil
}
func (f *fakeRoleRepo) ListWithPermissions(ctx context.Context) ([]model.RoleModel, error) {
    return nil, nil
}
func (f *fakeRoleRepo) Get(ctx context.Context, id uint64) (*model.RoleModel, error) {
    return nil, repository.ErrNotFound
}
func (f *fakeRoleRepo) GetWithPermissions(ctx context.Context, id uint64) (*model.RoleModel, error) {
    return nil, repository.ErrNotFound
}
func (f *fakeRoleRepo) GetBySlug(ctx context.Context, slug string) (*model.RoleModel, error) {
    return nil, repository.ErrNotFound
}
func (f *fakeRoleRepo) Create(ctx context.Context, role *model.RoleModel) error { return nil }
func (f *fakeRoleRepo) Update(ctx context.Context, role *model.RoleModel) error { return nil }
func (f *fakeRoleRepo) Delete(ctx context.Context, id uint64) error             { return nil }
func (f *fakeRoleRepo) AssignPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
    return nil
}
func (f *fakeRoleRepo) AddPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
    return nil
}
func (f *fakeRoleRepo) RemovePermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
    return nil
}
func (f *fakeRoleRepo) AssignToUser(ctx context.Context, userID uint64, roleIDs []uint64) error { return nil }
func (f *fakeRoleRepo) RemoveFromUser(ctx context.Context, userID uint64, roleIDs []uint64) error {
    return nil
}
func (f *fakeRoleRepo) ListByUserID(ctx context.Context, userID uint64) ([]model.RoleModel, error) { return nil, nil }
func (f *fakeRoleRepo) CheckUserHasRole(ctx context.Context, userID uint64, roleSlug string) (bool, error) {
    return false, nil
}

type fakeTxManager struct{ repos *common.Repos }

func (m *fakeTxManager) WithTx(ctx context.Context, fn func(r *common.Repos) error) error { return fn(m.repos) }

// ---- Tests covering cache, validation, state machine, pagination, tx ----

func TestAdmin_ListGames_UsesCacheAndInvalidatesOnWrite(t *testing.T) {
    gRepo := &fakeGameRepo{items: []model.Game{{Base: model.Base{ID: 1}, Key: "lol", Name: "League"}}}
    s := NewAdminService(gRepo, &fakeUserRepo{}, &fakePlayerRepo{}, &fakeOrderRepo{}, &fakePaymentRepo{}, &fakeRoleRepo{}, cache.NewMemory())

    ctx := context.Background()

    // first call hits repo
    if _, err := s.ListGames(ctx); err != nil {
        t.Fatalf("ListGames err: %v", err)
    }
    if gRepo.listCalls != 1 {
        t.Fatalf("expected 1 repo call, got %d", gRepo.listCalls)
    }
    // second call served from cache
    if _, err := s.ListGames(ctx); err != nil {
        t.Fatalf("ListGames err: %v", err)
    }
    if gRepo.listCalls != 1 {
        t.Fatalf("expected cached result, repo calls=%d", gRepo.listCalls)
    }
    // write invalidates cache
    if _, err := s.CreateGame(ctx, CreateGameInput{Key: "dota2", Name: "DOTA2"}); err != nil {
        t.Fatalf("CreateGame err: %v", err)
    }
    if _, err := s.ListGames(ctx); err != nil {
        t.Fatalf("ListGames err: %v", err)
    }
    if gRepo.listCalls != 2 {
        t.Fatalf("expected cache invalidation; repo calls=%d", gRepo.listCalls)
    }
}

func TestAdmin_CreateGame_Validation(t *testing.T) {
    s := NewAdminService(&fakeGameRepo{}, &fakeUserRepo{}, &fakePlayerRepo{}, &fakeOrderRepo{}, &fakePaymentRepo{}, &fakeRoleRepo{}, cache.NewMemory())
    if _, err := s.CreateGame(context.Background(), CreateGameInput{Key: "", Name: ""}); err == nil {
        t.Fatalf("expected validation error for empty key/name")
    }
}

func TestAdmin_UpdateOrder_Validation(t *testing.T) {
    now := time.Now()
    order := &model.Order{Base: model.Base{ID: 1}, Status: model.OrderStatusPending}
    oRepo := &fakeOrderRepo{obj: order}
    s := NewAdminService(&fakeGameRepo{}, &fakeUserRepo{}, &fakePlayerRepo{}, oRepo, &fakePaymentRepo{}, &fakeRoleRepo{}, cache.NewMemory())

    // invalid status
    if _, err := s.UpdateOrder(context.Background(), 1, UpdateOrderInput{Status: "bad", PriceCents: 1, Currency: model.CurrencyCNY}); err == nil {
        t.Fatalf("expected validation error for bad status")
    }
    // invalid currency
    if _, err := s.UpdateOrder(context.Background(), 1, UpdateOrderInput{Status: model.OrderStatusConfirmed, PriceCents: 1, Currency: "XYZ"}); err == nil {
        t.Fatalf("expected validation error for bad currency")
    }
    // negative price
    if _, err := s.UpdateOrder(context.Background(), 1, UpdateOrderInput{Status: model.OrderStatusConfirmed, PriceCents: -1, Currency: model.CurrencyCNY}); err == nil {
        t.Fatalf("expected validation error for negative price")
    }
    // valid update
    in := UpdateOrderInput{Status: model.OrderStatusCanceled, PriceCents: 100, Currency: model.CurrencyUSD, ScheduledStart: &now}
    out, err := s.UpdateOrder(context.Background(), 1, in)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if out.Status != model.OrderStatusCanceled || out.PriceCents != 100 || out.Currency != model.CurrencyUSD {
        t.Fatalf("unexpected order after update: %#v", out)
    }
}

func TestAdmin_UpdatePayment_Validation(t *testing.T) {
    p := &model.Payment{Base: model.Base{ID: 1}, Status: model.PaymentStatusPending}
    pRepo := &fakePaymentRepo{obj: p}
    s := NewAdminService(&fakeGameRepo{}, &fakeUserRepo{}, &fakePlayerRepo{}, &fakeOrderRepo{}, pRepo, &fakeRoleRepo{}, cache.NewMemory())

    // invalid status
    if _, err := s.UpdatePayment(context.Background(), 1, UpdatePaymentInput{Status: "oops"}); err == nil {
        t.Fatalf("expected validation error for bad status")
    }
    // valid update
    raw := json.RawMessage(`{"from":"test"}`)
    out, err := s.UpdatePayment(context.Background(), 1, UpdatePaymentInput{Status: model.PaymentStatusPaid, ProviderRaw: raw})
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if out.Status != model.PaymentStatusPaid {
        t.Fatalf("unexpected status: %v", out.Status)
    }
    if string(out.ProviderRaw) != string(raw) {
        t.Fatalf("unexpected raw: %s", string(out.ProviderRaw))
    }
}

func TestAdmin_OrderStateMachine(t *testing.T) {
    o := &model.Order{Base: model.Base{ID: 1}, Status: model.OrderStatusPending}
    oRepo := &fakeOrderRepo{obj: o}
    s := NewAdminService(&fakeGameRepo{}, &fakeUserRepo{}, &fakePlayerRepo{}, oRepo, &fakePaymentRepo{}, &fakeRoleRepo{}, cache.NewMemory())

    // pending -> confirmed ok
    if _, err := s.UpdateOrder(context.Background(), 1, UpdateOrderInput{Status: model.OrderStatusConfirmed, PriceCents: 1, Currency: model.CurrencyCNY}); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    // confirmed -> pending not allowed
    if _, err := s.UpdateOrder(context.Background(), 1, UpdateOrderInput{Status: model.OrderStatusPending, PriceCents: 1, Currency: model.CurrencyCNY}); err == nil {
        t.Fatalf("expected invalid transition error")
    }
}

func TestAdmin_PaymentStateMachine(t *testing.T) {
    p := &model.Payment{Base: model.Base{ID: 1}, Status: model.PaymentStatusPending}
    pRepo := &fakePaymentRepo{obj: p}
    s := NewAdminService(&fakeGameRepo{}, &fakeUserRepo{}, &fakePlayerRepo{}, &fakeOrderRepo{}, pRepo, &fakeRoleRepo{}, cache.NewMemory())

    // pending -> paid ok
    if _, err := s.UpdatePayment(context.Background(), 1, UpdatePaymentInput{Status: model.PaymentStatusPaid}); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    // paid -> failed not allowed
    if _, err := s.UpdatePayment(context.Background(), 1, UpdatePaymentInput{Status: model.PaymentStatusFailed}); err == nil {
        t.Fatalf("expected invalid transition error")
    }
}

func TestAdmin_Pagination_Normalization(t *testing.T) {
    uRepo := &fakeUserRepo{}
    s := NewAdminService(&fakeGameRepo{}, uRepo, &fakePlayerRepo{}, &fakeOrderRepo{}, &fakePaymentRepo{}, &fakeRoleRepo{}, cache.NewMemory())
    // page=0 and size=1000 should be normalized to 1 and 100
    _, p, err := s.ListUsersPaged(context.Background(), 0, 1000)
    if err != nil {
        t.Fatalf("ListUsersPaged err: %v", err)
    }
    if p.TotalPages != 3 || p.HasNext != true || p.HasPrev != false {
        t.Fatalf("unexpected pagination: %#v", p)
    }
    // last page
    _, p2, err := s.ListUsersPaged(context.Background(), 3, 100)
    if err != nil {
        t.Fatalf("ListUsersPaged err: %v", err)
    }
    if p2.TotalPages != 3 || p2.HasNext != false || p2.HasPrev != true {
        t.Fatalf("unexpected pagination last page: %#v", p2)
    }
}

func TestAdmin_RegisterUserAndPlayer_TxAndCacheInvalidation(t *testing.T) {
    mem := cache.NewMemory()
    g := &fakeGameRepo{}
    s := NewAdminService(g, &fakeUserRepo{}, &fakePlayerRepo{}, &fakeOrderRepo{}, &fakePaymentRepo{}, &fakeRoleRepo{}, mem)

    // when tx is nil
    if _, _, err := s.RegisterUserAndPlayer(context.Background(), CreateUserInput{Name: "alice", Password: "password1", Role: model.RoleUser, Status: model.UserStatusActive}, CreatePlayerInput{Nickname: "p", VerificationStatus: model.VerificationVerified}); err == nil {
        t.Fatalf("expected tx not configured error")
    }

    // prepare tx repos
    txU := &fakeUserRepo{}
    txP := &fakePlayerRepo{}
    s.SetTxManager(&fakeTxManager{repos: &common.Repos{Users: txU, Players: txP}})

    // seed caches and ensure invalidation happens later
    _ = mem.Set(context.Background(), cacheKeyUsers, "x", 10*time.Minute)
    _ = mem.Set(context.Background(), cacheKeyPlayers, "y", 10*time.Minute)

    u, p, err := s.RegisterUserAndPlayer(context.Background(), CreateUserInput{Name: "alice", Password: "password1", Role: model.RoleUser, Status: model.UserStatusActive}, CreatePlayerInput{Nickname: "pro", VerificationStatus: model.VerificationVerified})
    if err != nil {
        t.Fatalf("RegisterUserAndPlayer err: %v", err)
    }
    if u == nil || p == nil || u.ID == 0 || p.UserID != u.ID {
        t.Fatalf("unexpected created entities: user=%#v player=%#v", u, p)
    }
    if _, ok, _ := mem.Get(context.Background(), cacheKeyUsers); ok {
        t.Fatalf("users cache should be invalidated")
    }
    if _, ok, _ := mem.Get(context.Background(), cacheKeyPlayers); ok {
        t.Fatalf("players cache should be invalidated")
    }
}

// ===== 新增测试：用户管理 =====

func TestCreateUser_Success(t *testing.T) {
	users := &fakeUserRepo{}
	svc := NewAdminService(
		&fakeGameRepo{},
		users,
		&fakePlayerRepo{},
		&fakeOrderRepo{},
		&fakePaymentRepo{},
		&fakeRoleRepo{},
		cache.NewMemory(),
	)

	input := CreateUserInput{
		Phone:    "13812345678",
		Email:    "test@example.com",
		Password: "pass123",
		Name:     "Test User",
		Role:     model.RoleUser,
		Status:   model.UserStatusActive,
	}

	user, err := svc.CreateUser(context.Background(), input)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if user.ID == 0 {
		t.Error("Expected user ID to be set")
	}
	if user.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got '%s'", user.Name)
	}
	if user.PasswordHash == "" {
		t.Error("Expected password hash to be set")
	}
}

func TestCreateUser_InvalidInput(t *testing.T) {
	svc := NewAdminService(
		&fakeGameRepo{},
		&fakeUserRepo{},
		&fakePlayerRepo{},
		&fakeOrderRepo{},
		&fakePaymentRepo{},
		&fakeRoleRepo{},
		cache.NewMemory(),
	)

	tests := []struct {
		name  string
		input CreateUserInput
	}{
		{
			name: "空名称",
			input: CreateUserInput{
				Name:     "",
				Role:     model.RoleUser,
				Status:   model.UserStatusActive,
				Password: "pass123",
			},
		},
		{
			name: "空角色",
			input: CreateUserInput{
				Name:     "Test",
				Role:     "",
				Status:   model.UserStatusActive,
				Password: "pass123",
			},
		},
		{
			name: "空状态",
			input: CreateUserInput{
				Name:     "Test",
				Role:     model.RoleUser,
				Status:   "",
				Password: "pass123",
			},
		},
		{
			name: "无效密码",
			input: CreateUserInput{
				Name:     "Test",
				Role:     model.RoleUser,
				Status:   model.UserStatusActive,
				Password: "123", // too short, no letter
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateUser(context.Background(), tt.input)
			if err != ErrValidation {
				t.Errorf("Expected ErrValidation, got %v", err)
			}
		})
	}
}

func TestUpdateUser_Success(t *testing.T) {
	users := &fakeUserRepo{
		last: &model.User{
			Base:   model.Base{ID: 42},
			Name:   "Old Name",
			Role:   model.RoleUser,
			Status: model.UserStatusActive,
			Phone:  "13800000000",
			Email:  "old@example.com",
		},
	}

	svc := NewAdminService(
		&fakeGameRepo{},
		users,
		&fakePlayerRepo{},
		&fakeOrderRepo{},
		&fakePaymentRepo{},
		&fakeRoleRepo{},
		cache.NewMemory(),
	)

	input := UpdateUserInput{
		Name:   "New Name",
		Role:   model.RoleAdmin,
		Status: model.UserStatusSuspended,
	}

	user, err := svc.UpdateUser(context.Background(), 42, input)
	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}

	if user.Name != "New Name" {
		t.Errorf("Expected name 'New Name', got '%s'", user.Name)
	}
	if user.Role != model.RoleAdmin {
		t.Errorf("Expected role Admin, got %s", user.Role)
	}
}

func TestValidPassword(t *testing.T) {
	tests := []struct {
		password string
		valid    bool
	}{
		{"pass123", true},    // 有字母和数字，长度>=6
		{"abc123", true},     // 有字母和数字，长度>=6
		{"Test99", true},     // 有字母和数字，长度>=6
		{"123456", false},    // 只有数字
		{"abcdef", false},    // 只有字母
		{"12345", false},     // 长度<6
		{"abc", false},       // 长度<6且只有字母
		{"", false},          // 空字符串
		{"Pass1", false},     // 长度不足6
		{"Password1", true},  // 有字母和数字，长度>=6
	}

	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			result := validPassword(tt.password)
			if result != tt.valid {
				t.Errorf("validPassword(%q) = %v, want %v", tt.password, result, tt.valid)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	t.Run("成功生成密码哈希", func(t *testing.T) {
		hash, err := hashPassword("test123")
		if err != nil {
			t.Fatalf("hashPassword failed: %v", err)
		}
		if hash == "" {
			t.Error("Expected non-empty hash")
		}
		if hash == "test123" {
			t.Error("Hash should not equal plain password")
		}
	})

	t.Run("空密码返回错误", func(t *testing.T) {
		_, err := hashPassword("")
		if err != ErrValidation {
			t.Errorf("Expected ErrValidation for empty password, got %v", err)
		}
	})

	t.Run("空格密码返回错误", func(t *testing.T) {
		_, err := hashPassword("   ")
		if err != ErrValidation {
			t.Errorf("Expected ErrValidation for whitespace password, got %v", err)
		}
	})
}

func TestValidateGameInput(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		gameName  string
		expectErr bool
	}{
		{"有效输入", "lol", "League of Legends", false},
		{"Key为空", "", "LOL", true},
		{"Name为空", "lol", "", true},
		{"Key为空格", "   ", "LOL", true},
		{"Name为空格", "lol", "   ", true},
		{"都为空", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateGameInput(tt.key, tt.gameName)
			if tt.expectErr && err != ErrValidation {
				t.Errorf("Expected ErrValidation, got %v", err)
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}

func TestValidateUserInput(t *testing.T) {
	tests := []struct {
		name      string
		userName  string
		role      model.Role
		status    model.UserStatus
		password  string
		expectErr bool
	}{
		{"有效输入", "John", model.RoleUser, model.UserStatusActive, "pass123", false},
		{"空名称", "", model.RoleUser, model.UserStatusActive, "pass123", true},
		{"空角色", "John", "", model.UserStatusActive, "pass123", true},
		{"空状态", "John", model.RoleUser, "", "pass123", true},
		{"无效密码", "John", model.RoleUser, model.UserStatusActive, "123", true},
		{"空密码可接受", "John", model.RoleUser, model.UserStatusActive, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUserInput(tt.userName, tt.role, tt.status, tt.password)
			if tt.expectErr && err != ErrValidation {
				t.Errorf("Expected ErrValidation, got %v", err)
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}

func TestValidatePlayerInput(t *testing.T) {
	tests := []struct {
		name         string
		userID       uint64
		verification model.VerificationStatus
		expectErr    bool
	}{
		{"有效输入", 42, model.VerificationVerified, false},
		{"UserID为0", 0, model.VerificationVerified, true},
		{"空验证状态", 42, "", true},
		{"都无效", 0, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePlayerInput(tt.userID, tt.verification)
			if tt.expectErr && err != ErrValidation {
				t.Errorf("Expected ErrValidation, got %v", err)
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}

func TestBuildPagination(t *testing.T) {
	tests := []struct {
		name            string
		page, pageSize  int
		total           int64
		expectedPages   int
		expectedHasNext bool
		expectedHasPrev bool
	}{
		{"第1页", 1, 20, 100, 5, true, false},
		{"中间页", 3, 20, 100, 5, true, true},
		{"最后一页", 5, 20, 100, 5, false, true},
		{"单页", 1, 100, 50, 1, false, false},
		{"空结果", 1, 20, 0, 0, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := buildPagination(tt.page, tt.pageSize, tt.total)
			if p.TotalPages != tt.expectedPages {
				t.Errorf("Expected %d total pages, got %d", tt.expectedPages, p.TotalPages)
			}
			if p.HasNext != tt.expectedHasNext {
				t.Errorf("Expected HasNext=%v, got %v", tt.expectedHasNext, p.HasNext)
			}
			if p.HasPrev != tt.expectedHasPrev {
				t.Errorf("Expected HasPrev=%v, got %v", tt.expectedHasPrev, p.HasPrev)
			}
		})
	}
}

func TestGetUser_Success(t *testing.T) {
	users := &fakeUserRepo{
		last: &model.User{Base: model.Base{ID: 42}, Name: "Test User"},
	}

	svc := NewAdminService(
		&fakeGameRepo{},
		users,
		&fakePlayerRepo{},
		&fakeOrderRepo{},
		&fakePaymentRepo{},
		&fakeRoleRepo{},
		cache.NewMemory(),
	)

	user, err := svc.GetUser(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}

	if user.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got '%s'", user.Name)
	}
}

func TestDeleteUser_Success(t *testing.T) {
	svc := NewAdminService(
		&fakeGameRepo{},
		&fakeUserRepo{},
		&fakePlayerRepo{},
		&fakeOrderRepo{},
		&fakePaymentRepo{},
		&fakeRoleRepo{},
		cache.NewMemory(),
	)

	err := svc.DeleteUser(context.Background(), 42)
	if err != nil {
		t.Errorf("DeleteUser failed: %v", err)
	}
}

func TestCreateGame_Success(t *testing.T) {
	games := &fakeGameRepo{}
	svc := NewAdminService(
		games,
		&fakeUserRepo{},
		&fakePlayerRepo{},
		&fakeOrderRepo{},
		&fakePaymentRepo{},
		&fakeRoleRepo{},
		cache.NewMemory(),
	)

	input := CreateGameInput{
		Key:         "lol",
		Name:        "League of Legends",
		Category:    "MOBA",
		IconURL:     "https://example.com/icon.png",
		Description: "5v5 team game",
	}

	game, err := svc.CreateGame(context.Background(), input)
	if err != nil {
		t.Fatalf("CreateGame failed: %v", err)
	}

	if game.Key != "lol" {
		t.Errorf("Expected key 'lol', got '%s'", game.Key)
	}
}

func TestUpdateGame_Success(t *testing.T) {
	games := &fakeGameRepo{
		items: []model.Game{
			{Base: model.Base{ID: 1}, Key: "lol", Name: "LOL"},
		},
	}

	svc := NewAdminService(
		games,
		&fakeUserRepo{},
		&fakePlayerRepo{},
		&fakeOrderRepo{},
		&fakePaymentRepo{},
		&fakeRoleRepo{},
		cache.NewMemory(),
	)

	input := UpdateGameInput{
		Key:  "lol",
		Name: "League of Legends Updated",
	}

	game, err := svc.UpdateGame(context.Background(), 1, input)
	if err != nil {
		t.Fatalf("UpdateGame failed: %v", err)
	}

	if game.Name != "League of Legends Updated" {
		t.Errorf("Expected name updated, got '%s'", game.Name)
	}
}

func TestGetGame_Success(t *testing.T) {
	games := &fakeGameRepo{
		items: []model.Game{
			{Base: model.Base{ID: 1}, Key: "lol", Name: "LOL"},
		},
	}

	svc := NewAdminService(
		games,
		&fakeUserRepo{},
		&fakePlayerRepo{},
		&fakeOrderRepo{},
		&fakePaymentRepo{},
		&fakeRoleRepo{},
		cache.NewMemory(),
	)

	game, err := svc.GetGame(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetGame failed: %v", err)
	}

	if game.Key != "lol" {
		t.Errorf("Expected key 'lol', got '%s'", game.Key)
	}
}

func TestDeleteGame_Success(t *testing.T) {
	svc := NewAdminService(
		&fakeGameRepo{},
		&fakeUserRepo{},
		&fakePlayerRepo{},
		&fakeOrderRepo{},
		&fakePaymentRepo{},
		&fakeRoleRepo{},
		cache.NewMemory(),
	)

	err := svc.DeleteGame(context.Background(), 1)
	if err != nil {
		t.Errorf("DeleteGame failed: %v", err)
	}
}

func TestCreatePlayer_Success(t *testing.T) {
	players := &fakePlayerRepo{}
	svc := NewAdminService(
		&fakeGameRepo{},
		&fakeUserRepo{},
		players,
		&fakeOrderRepo{},
		&fakePaymentRepo{},
		&fakeRoleRepo{},
		cache.NewMemory(),
	)

	input := CreatePlayerInput{
		UserID:             42,
		Nickname:           "ProGamer",
		Bio:                "Expert gamer",
		HourlyRateCents:    5000,
		MainGameID:         1,
		VerificationStatus: model.VerificationVerified,
	}

	player, err := svc.CreatePlayer(context.Background(), input)
	if err != nil {
		t.Fatalf("CreatePlayer failed: %v", err)
	}

	if player.Nickname != "ProGamer" {
		t.Errorf("Expected nickname 'ProGamer', got '%s'", player.Nickname)
	}
}

func TestCreatePlayer_InvalidInput(t *testing.T) {
	svc := NewAdminService(
		&fakeGameRepo{},
		&fakeUserRepo{},
		&fakePlayerRepo{},
		&fakeOrderRepo{},
		&fakePaymentRepo{},
		&fakeRoleRepo{},
		cache.NewMemory(),
	)

	tests := []struct {
		name  string
		input CreatePlayerInput
	}{
		{
			name: "UserID为0",
			input: CreatePlayerInput{
				UserID:             0,
				VerificationStatus: model.VerificationVerified,
			},
		},
		{
			name: "空验证状态",
			input: CreatePlayerInput{
				UserID:             42,
				VerificationStatus: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreatePlayer(context.Background(), tt.input)
			if err != ErrValidation {
				t.Errorf("Expected ErrValidation, got %v", err)
			}
		})
	}
}

func TestUpdatePlayer_Success(t *testing.T) {
	players := &fakePlayerRepo{
		last: &model.Player{
			Base:               model.Base{ID: 99},
			UserID:             42,
			Nickname:           "OldNick",
			VerificationStatus: model.VerificationPending,
		},
	}

	svc := NewAdminService(
		&fakeGameRepo{},
		&fakeUserRepo{},
		players,
		&fakeOrderRepo{},
		&fakePaymentRepo{},
		&fakeRoleRepo{},
		cache.NewMemory(),
	)

	input := UpdatePlayerInput{
		Nickname:           "NewNick",
		Bio:                "Updated bio",
		HourlyRateCents:    6000,
		MainGameID:         2,
		VerificationStatus: model.VerificationVerified,
	}

	player, err := svc.UpdatePlayer(context.Background(), 99, input)
	if err != nil {
		t.Fatalf("UpdatePlayer failed: %v", err)
	}

	if player.Nickname != "NewNick" {
		t.Errorf("Expected nickname 'NewNick', got '%s'", player.Nickname)
	}
}

func TestGetPlayer_Success(t *testing.T) {
	players := &fakePlayerRepo{
		last: &model.Player{
			Base:     model.Base{ID: 99},
			Nickname: "TestPlayer",
		},
	}

	svc := NewAdminService(
		&fakeGameRepo{},
		&fakeUserRepo{},
		players,
		&fakeOrderRepo{},
		&fakePaymentRepo{},
		&fakeRoleRepo{},
		cache.NewMemory(),
	)

	player, err := svc.GetPlayer(context.Background(), 99)
	if err != nil {
		t.Fatalf("GetPlayer failed: %v", err)
	}

	if player.Nickname != "TestPlayer" {
		t.Errorf("Expected nickname 'TestPlayer', got '%s'", player.Nickname)
	}
}

func TestDeletePlayer_Success(t *testing.T) {
	svc := NewAdminService(
		&fakeGameRepo{},
		&fakeUserRepo{},
		&fakePlayerRepo{},
		&fakeOrderRepo{},
		&fakePaymentRepo{},
		&fakeRoleRepo{},
		cache.NewMemory(),
	)

	err := svc.DeletePlayer(context.Background(), 99)
	if err != nil {
		t.Errorf("DeletePlayer failed: %v", err)
	}
}