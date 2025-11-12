package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"

	"gamelink/internal/repository"
)

// ChatRetentionScheduler purges chat data after retention period.
type ChatRetentionScheduler struct {
	groups        repository.ChatGroupRepository
	messages      repository.ChatMessageRepository
	cron          *cron.Cron
	RetentionDays int
}

func NewChatRetentionScheduler(groups repository.ChatGroupRepository, messages repository.ChatMessageRepository, retentionDays int) *ChatRetentionScheduler {
	if retentionDays <= 0 {
		retentionDays = 30
	}
	return &ChatRetentionScheduler{
		groups:        groups,
		messages:      messages,
		cron:          cron.New(),
		RetentionDays: retentionDays,
	}
}

// Start runs a daily purge at 03:15.
func (s *ChatRetentionScheduler) Start() {
	_, err := s.cron.AddFunc("15 3 * * *", s.purge)
	if err != nil {
		log.Printf("[ChatRetention] add job error: %v", err)
		return
	}
	s.cron.Start()
	log.Println("[ChatRetention] scheduler started - daily at 03:15")
}

func (s *ChatRetentionScheduler) Stop() { s.cron.Stop() }

// PurgeOnce allows manual purge for tests.
func (s *ChatRetentionScheduler) PurgeOnce() { s.purge() }

func (s *ChatRetentionScheduler) purge() {
	ctx := context.Background()
	cutoff := time.Now().AddDate(0, 0, -s.RetentionDays)
	const batch = 500
	groups, err := s.groups.ListDeactivatedBefore(ctx, cutoff, batch)
	if err != nil {
		log.Printf("[ChatRetention] list groups error: %v", err)
		return
	}
	if len(groups) == 0 {
		return
	}

	ids := make([]uint64, 0, len(groups))
	for i := range groups {
		ids = append(ids, groups[i].ID)
	}

	if err := s.messages.DeleteByGroupIDs(ctx, ids); err != nil {
		log.Printf("[ChatRetention] delete messages error: %v", err)
		return
	}
	if err := s.groups.DeleteByIDs(ctx, ids); err != nil {
		log.Printf("[ChatRetention] delete groups error: %v", err)
		return
	}
	log.Printf("[ChatRetention] purged %d groups older than %s", len(ids), cutoff.Format(time.RFC3339))
}
