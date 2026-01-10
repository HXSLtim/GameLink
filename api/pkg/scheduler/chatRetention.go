package scheduler

import (
	"context"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"

	"gamelink/internal/repository"
	"gamelink/pkg/cache"
)

// ChatRetentionScheduler purges chat data after retention period.
type ChatRetentionScheduler struct {
	groups        repository.ChatGroupRepository
	messages      repository.ChatMessageRepository
	lock          cache.DistributedLock  // Redis distributed lock
	logger        *slog.Logger
	cron          *cron.Cron
	RetentionDays int
}

func NewChatRetentionScheduler(groups repository.ChatGroupRepository, messages repository.ChatMessageRepository, lock cache.DistributedLock, retentionDays int) *ChatRetentionScheduler {
	if retentionDays <= 0 {
		retentionDays = 30
	}
	return &ChatRetentionScheduler{
		groups:        groups,
		messages:      messages,
		lock:          lock,
		logger:        slog.Default(),
		cron:          cron.New(),
		RetentionDays: retentionDays,
	}
}

// Start runs a daily purge at 03:15.
func (s *ChatRetentionScheduler) Start() {
	_, err := s.cron.AddFunc("15 3 * * *", s.purgeWithLock)
	if err != nil {
		s.logger.Error("Failed to add chat retention job", "error", err)
		return
	}
	s.cron.Start()
	s.logger.Info("Chat retention scheduler started - daily at 03:15 with distributed lock")
}

func (s *ChatRetentionScheduler) Stop() {
	s.cron.Stop()
	s.logger.Info("Chat retention scheduler stopped")
}

// PurgeOnce allows manual purge for tests.
func (s *ChatRetentionScheduler) PurgeOnce() {
	s.purgeWithLock()
}

// purgeWithLock executes purge with distributed lock
func (s *ChatRetentionScheduler) purgeWithLock() {
	ctx := context.Background()
	lockKey := "scheduler:chat:purge"

	// Try to acquire distributed lock with 30 minutes TTL
	// Chat purge should complete within 30 minutes
	locked, err := s.lock.TryLock(ctx, lockKey, 30*time.Minute, 1, time.Second)
	if err != nil {
		s.logger.Error("Failed to acquire chat purge lock", "error", err, "key", lockKey)
		return
	}

	if !locked {
		s.logger.Info("Another instance is running chat purge, skipping", "key", lockKey)
		return
	}

	s.logger.Info("Acquired chat purge lock, starting purge", "key", lockKey)

	// Ensure lock is released
	defer func() {
		if unlockErr := s.lock.Unlock(ctx, lockKey); unlockErr != nil {
			s.logger.Error("Failed to release chat purge lock", "error", unlockErr, "key", lockKey)
		} else {
			s.logger.Info("Released chat purge lock", "key", lockKey)
		}
	}()

	// Execute the actual purge
	s.purge(ctx)
}

func (s *ChatRetentionScheduler) purge(ctx context.Context) {
	cutoff := time.Now().AddDate(0, 0, -s.RetentionDays)
	const batch = 500
	groups, err := s.groups.ListDeactivatedBefore(ctx, cutoff, batch)
	if err != nil {
		s.logger.Error("Failed to list groups for purge", "error", err, "cutoff", cutoff)
		return
	}
	if len(groups) == 0 {
		s.logger.Debug("No groups to purge", "cutoff", cutoff)
		return
	}

	ids := make([]uint64, 0, len(groups))
	for i := range groups {
		ids = append(ids, groups[i].ID)
	}

	if err := s.messages.DeleteByGroupIDs(ctx, ids); err != nil {
		s.logger.Error("Failed to delete messages", "error", err, "count", len(ids))
		return
	}
	if err := s.groups.DeleteByIDs(ctx, ids); err != nil {
		s.logger.Error("Failed to delete groups", "error", err, "count", len(ids))
		return
	}
	s.logger.Info("Purged chat data", "groups", len(ids), "cutoff", cutoff.Format(time.RFC3339))
}
