package chat

import (
    "testing"
    "context"
    "time"

    "gamelink/internal/cache"
    "gamelink/internal/model"
    "gamelink/internal/repository"
)

type gRepo struct{ g model.ChatGroup }
func (r gRepo) Create(ctx context.Context, group *model.ChatGroup) error { return nil }
func (r gRepo) Get(ctx context.Context, id uint64) (*model.ChatGroup, error) { return &r.g, nil }
func (r gRepo) GetByRelatedOrderID(ctx context.Context, orderID uint64) (*model.ChatGroup, error) { return &r.g, nil }
func (r gRepo) ListByUser(ctx context.Context, userID uint64, opts repository.ChatGroupListOptions) ([]model.ChatGroup, int64, error) { return []model.ChatGroup{r.g}, 1, nil }
func (r gRepo) ListMembers(ctx context.Context, groupID uint64, opts repository.ChatGroupMemberListOptions) ([]model.ChatGroupMember, int64, error) { return nil, 0, nil }
func (r gRepo) Update(ctx context.Context, group *model.ChatGroup) error { return nil }
func (r gRepo) Deactivate(ctx context.Context, id uint64) error { return nil }
func (r gRepo) ListDeactivatedBefore(ctx context.Context, cutoff time.Time, limit int) ([]model.ChatGroup, error) { return nil, nil }
func (r gRepo) DeleteByIDs(ctx context.Context, ids []uint64) error { return nil }

type mRepo struct{ active bool; lastRead *uint64; lastMod *model.ChatMessageAuditStatus; lastModMid *uint64; lastModReason string }
func (r mRepo) Add(ctx context.Context, member *model.ChatGroupMember) error { return nil }
func (r mRepo) AddBatch(ctx context.Context, members []*model.ChatGroupMember) error { return nil }
func (r mRepo) Update(ctx context.Context, member *model.ChatGroupMember) error { return nil }
func (r mRepo) Remove(ctx context.Context, groupID, userID uint64) error { return nil }
func (r mRepo) Get(ctx context.Context, groupID, userID uint64) (*model.ChatGroupMember, error) {
    if !r.active { return nil, repository.ErrNotFound }
    return &model.ChatGroupMember{GroupID: groupID, UserID: userID, IsActive: true}, nil
}

type msgRepo struct{ created *model.ChatMessage }
func (r msgRepo) Create(ctx context.Context, message *model.ChatMessage) error { r.created = message; return nil }
func (r msgRepo) CreateBatch(ctx context.Context, messages []*model.ChatMessage) error { return nil }
func (r msgRepo) ListByGroup(ctx context.Context, opts repository.ChatMessageListOptions) ([]model.ChatMessage, int64, error) { return nil, 0, nil }
func (r msgRepo) Get(ctx context.Context, id uint64) (*model.ChatMessage, error) { return nil, repository.ErrNotFound }
func (r msgRepo) MarkDeleted(ctx context.Context, id uint64, deletedBy uint64) error { return nil }
func (r msgRepo) ListForModeration(ctx context.Context, opts repository.ChatMessageModerationListOptions) ([]model.ChatMessage, int64, error) { return nil, 0, nil }
func (r msgRepo) UpdateAuditStatus(ctx context.Context, id uint64, status model.ChatMessageAuditStatus, moderatorID *uint64, reason string) error { return nil }
func (r msgRepo) DeleteByGroupIDs(ctx context.Context, groupIDs []uint64) error { return nil }

type repRepo struct{ created bool }
func (r *repRepo) Create(ctx context.Context, report *model.ChatReport) error { r.created = true; return nil }
func (r *repRepo) Get(ctx context.Context, id uint64) (*model.ChatReport, error) { return nil, repository.ErrNotFound }
func (r *repRepo) Update(ctx context.Context, report *model.ChatReport) error { return nil }
func (r *repRepo) List(ctx context.Context, opts repository.ChatReportListOptions) ([]model.ChatReport, int64, error) { return nil, 0, nil }

func TestEnsureMembership_NotMember(t *testing.T) {
    s := NewChatService(gRepo{g: model.ChatGroup{Base: model.Base{ID: 1}, IsActive: true}}, mRepo{active: false}, msgRepo{}, &repRepo{}, cache.NewMemory())
    if _, err := s.EnsureMembership(context.Background(), 1, 1); err == nil { t.Fatalf("expected error") }
}

func TestSendMessage_PublicThrottle(t *testing.T) {
    grp := gRepo{g: model.ChatGroup{Base: model.Base{ID: 1}, IsActive: true, GroupType: model.ChatGroupTypePublic}}
    s := NewChatService(grp, mRepo{active: true}, msgRepo{}, &repRepo{}, cache.NewMemory())
    _, err := s.SendMessage(context.Background(), SendMessageInput{GroupID:1, SenderID:1, Content:"a"})
    if err != nil { t.Fatalf("%v", err) }
    _, err = s.SendMessage(context.Background(), SendMessageInput{GroupID:1, SenderID:1, Content:"b"})
    if err == nil { t.Fatalf("expected error") }
}

func TestReportMessage_Creates(t *testing.T) {
    r := &repRepo{}
    s := NewChatService(gRepo{g: model.ChatGroup{Base: model.Base{ID: 1}, IsActive: true}}, mRepo{active: true}, msgRepo{}, r, cache.NewMemory())
    if err := s.ReportMessage(context.Background(), 1, 2, "x", "y"); err != nil { t.Fatalf("%v", err) }
    if !r.created { t.Fatalf("not created") }
}

func TestListMessages_PublicApprovedFilter(t *testing.T) {
    grp := gRepo{g: model.ChatGroup{Base: model.Base{ID: 1}, IsActive: true, GroupType: model.ChatGroupTypePublic}}
    s := NewChatService(grp, mRepo{active: true}, msgRepo{}, &repRepo{}, cache.NewMemory())
    if _, _, err := s.ListMessages(context.Background(), 1, 1, ListMessagesOptions{Page:1, PageSize:10}); err != nil { t.Fatalf("%v", err) }
}

func TestApproveRejectMessage(t *testing.T) {
    s := NewChatService(gRepo{g: model.ChatGroup{Base: model.Base{ID: 1}, IsActive: true}}, mRepo{active: true}, msgRepo{}, &repRepo{}, cache.NewMemory())
    mid := uint64(9)
    if err := s.ApproveMessage(context.Background(), 1, mid); err != nil { t.Fatalf("%v", err) }
    if err := s.RejectMessage(context.Background(), 1, mid, "r"); err != nil { t.Fatalf("%v", err) }
}

// stateful member repo for join/leave/mark-read
type memMembers struct{ store map[uint64]map[uint64]*model.ChatGroupMember }
func newMemMembers() *memMembers { return &memMembers{store: map[uint64]map[uint64]*model.ChatGroupMember{}} }
func (m *memMembers) Add(ctx context.Context, member *model.ChatGroupMember) error {
    if m.store[member.GroupID] == nil { m.store[member.GroupID] = map[uint64]*model.ChatGroupMember{} }
    m.store[member.GroupID][member.UserID] = member
    return nil
}
func (m *memMembers) AddBatch(ctx context.Context, members []*model.ChatGroupMember) error { for _, mem := range members { _ = m.Add(ctx, mem) } ; return nil }
func (m *memMembers) Update(ctx context.Context, member *model.ChatGroupMember) error { if m.store[member.GroupID] == nil { m.store[member.GroupID] = map[uint64]*model.ChatGroupMember{} } ; m.store[member.GroupID][member.UserID] = member ; return nil }
func (m *memMembers) Remove(ctx context.Context, groupID, userID uint64) error { if m.store[groupID]!=nil { delete(m.store[groupID], userID) } ; return nil }
func (m *memMembers) Get(ctx context.Context, groupID, userID uint64) (*model.ChatGroupMember, error) { if m.store[groupID]==nil { return nil, repository.ErrNotFound } ; if mem := m.store[groupID][userID]; mem != nil { return mem, nil } ; return nil, repository.ErrNotFound }

func TestJoinLeaveMarkReadFlows(t *testing.T) {
    grp := gRepo{g: model.ChatGroup{Base: model.Base{ID: 7}, IsActive: true}}
    mem := newMemMembers()
    msg := msgRepo{}
    s := NewChatService(grp, mem, msg, &repRepo{}, cache.NewMemory())
    // first join creates member
    if err := s.JoinGroup(context.Background(), 7, 100, "nick"); err != nil { t.Fatalf("%v", err) }
    // join again updates nickname
    if err := s.JoinGroup(context.Background(), 7, 100, "nick2"); err != nil { t.Fatalf("%v", err) }
    // mark read
    if err := s.MarkRead(context.Background(), 7, 100, 5); err != nil { t.Fatalf("%v", err) }
    // leave
    if err := s.LeaveGroup(context.Background(), 7, 100); err != nil { t.Fatalf("%v", err) }
    // mark read after leave should still succeed (MarkRead does not enforce active)
    if err := s.MarkRead(context.Background(), 7, 100, 6); err != nil { t.Fatalf("%v", err) }
}

func TestSendMessage_TooLargeAndInactive(t *testing.T) {
    grp := gRepo{g: model.ChatGroup{Base: model.Base{ID: 8}, IsActive: false}}
    mem := newMemMembers()
    _ = mem.Add(context.Background(), &model.ChatGroupMember{GroupID:8, UserID:200, IsActive:true})
    s := NewChatService(grp, mem, msgRepo{}, &repRepo{}, cache.NewMemory())
    // empty content and no image
    if _, err := s.SendMessage(context.Background(), SendMessageInput{GroupID:8, SenderID:200, Content:""}); err == nil { t.Fatalf("expected error") }
    // too long content
    long := make([]rune, 2001); for i:=range long { long[i] = 'a' }
    if _, err := s.SendMessage(context.Background(), SendMessageInput{GroupID:8, SenderID:200, Content:string(long)}); err == nil { t.Fatalf("expected error") }
    // inactive group
    if _, err := s.SendMessage(context.Background(), SendMessageInput{GroupID:8, SenderID:200, Content:"ok"}); err == nil { t.Fatalf("expected error") }
}