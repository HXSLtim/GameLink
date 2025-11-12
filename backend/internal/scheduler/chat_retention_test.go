package scheduler

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

type fakeGroupRepo struct {
	listedAt    time.Time
	limit       int
	returned    []model.ChatGroup
	deletedIDs  []uint64
}

func (r *fakeGroupRepo) Create(ctx context.Context, group *model.ChatGroup) error { return nil }
func (r *fakeGroupRepo) Get(ctx context.Context, id uint64) (*model.ChatGroup, error) { return nil, repository.ErrNotFound }
func (r *fakeGroupRepo) GetByRelatedOrderID(ctx context.Context, orderID uint64) (*model.ChatGroup, error) { return nil, repository.ErrNotFound }
func (r *fakeGroupRepo) ListByUser(ctx context.Context, userID uint64, opts repository.ChatGroupListOptions) ([]model.ChatGroup, int64, error) { return nil, 0, nil }
func (r *fakeGroupRepo) ListMembers(ctx context.Context, groupID uint64, opts repository.ChatGroupMemberListOptions) ([]model.ChatGroupMember, int64, error) { return nil, 0, nil }
func (r *fakeGroupRepo) Update(ctx context.Context, group *model.ChatGroup) error { return nil }
func (r *fakeGroupRepo) Deactivate(ctx context.Context, id uint64) error { return nil }
func (r *fakeGroupRepo) ListDeactivatedBefore(ctx context.Context, cutoff time.Time, limit int) ([]model.ChatGroup, error) {
	r.listedAt = cutoff
	r.limit = limit
	return r.returned, nil
}
func (r *fakeGroupRepo) DeleteByIDs(ctx context.Context, ids []uint64) error {
	r.deletedIDs = append(r.deletedIDs, ids...)
	return nil
}

type fakeMessageRepo struct{ deletedGroupIDs []uint64 }

func (m *fakeMessageRepo) Create(ctx context.Context, message *model.ChatMessage) error { return nil }
func (m *fakeMessageRepo) CreateBatch(ctx context.Context, messages []*model.ChatMessage) error { return nil }
func (m *fakeMessageRepo) ListByGroup(ctx context.Context, opts repository.ChatMessageListOptions) ([]model.ChatMessage, int64, error) { return nil, 0, nil }
func (m *fakeMessageRepo) Get(ctx context.Context, id uint64) (*model.ChatMessage, error) { return nil, repository.ErrNotFound }
func (m *fakeMessageRepo) MarkDeleted(ctx context.Context, id uint64, deletedBy uint64) error { return nil }
func (m *fakeMessageRepo) ListForModeration(ctx context.Context, opts repository.ChatMessageModerationListOptions) ([]model.ChatMessage, int64, error) { return nil, 0, nil }
func (m *fakeMessageRepo) UpdateAuditStatus(ctx context.Context, id uint64, status model.ChatMessageAuditStatus, moderatorID *uint64, reason string) error { return nil }
func (m *fakeMessageRepo) DeleteByGroupIDs(ctx context.Context, groupIDs []uint64) error {
	m.deletedGroupIDs = append(m.deletedGroupIDs, groupIDs...)
	return nil
}

func TestChatRetentionScheduler_PurgeOnce(t *testing.T) {
	g := &fakeGroupRepo{ returned: []model.ChatGroup{
		{ Base: model.Base{ID: 1}, GroupType: model.ChatGroupTypeOrder },
		{ Base: model.Base{ID: 2}, GroupType: model.ChatGroupTypeOrder },
	}}
	m := &fakeMessageRepo{}
	s := NewChatRetentionScheduler(g, m, 30)
	// run once
	s.PurgeOnce()
	// verify deleted message groups & groups
	if len(m.deletedGroupIDs) != 2 || m.deletedGroupIDs[0] != 1 || m.deletedGroupIDs[1] != 2 {
		t.Fatalf("unexpected msg delete calls: %+v", m.deletedGroupIDs)
	}
	if len(g.deletedIDs) != 2 || g.deletedIDs[0] != 1 || g.deletedIDs[1] != 2 {
		t.Fatalf("unexpected group deletes: %+v", g.deletedIDs)
	}
}

func TestChatRetentionScheduler_StartStop(t *testing.T) {
	g := &fakeGroupRepo{}
	m := &fakeMessageRepo{}
	s := NewChatRetentionScheduler(g, m, 30)
	s.Start()
	s.Stop()
}
