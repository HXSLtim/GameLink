package cache

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// DistributedLock provides distributed locking capabilities
type DistributedLock interface {
	// Lock attempts to acquire a lock
	Lock(ctx context.Context, key string, ttl time.Duration) (bool, error)
	// Unlock releases a lock
	Unlock(ctx context.Context, key string) error
	// TryLock attempts to acquire a lock with retry
	TryLock(ctx context.Context, key string, ttl time.Duration, retry int, interval time.Duration) (bool, error)
}

// NewDistributedLock creates a new distributed lock instance
func NewDistributedLock(cache Cache) DistributedLock {
	return &distributedLockImpl{
		cache: cache,
	}
}

type distributedLockImpl struct {
	cache Cache
}

// Lock attempts to acquire a lock by setting a key with NX semantics
func (l *distributedLockImpl) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	lockKey := fmt.Sprintf("lock:%s", key)
	// Try to set the key with NX semantics (only if not exists)
	// For Redis, this would be SETNX command
	// For memory cache, we simulate this by checking existence first

	_, exists, err := l.cache.Get(ctx, lockKey)
	if err != nil {
		return false, err
	}

	if exists {
		return false, nil // Lock already held
	}

	// Set the lock
	err = l.cache.Set(ctx, lockKey, "1", ttl)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Unlock releases a lock
func (l *distributedLockImpl) Unlock(ctx context.Context, key string) error {
	lockKey := fmt.Sprintf("lock:%s", key)
	return l.cache.Delete(ctx, lockKey)
}

// TryLock attempts to acquire a lock with retry
func (l *distributedLockImpl) TryLock(ctx context.Context, key string, ttl time.Duration, retry int, interval time.Duration) (bool, error) {
	for i := 0; i < retry; i++ {
		locked, err := l.Lock(ctx, key, ttl)
		if err != nil {
			return false, err
		}

		if locked {
			return true, nil
		}

		if i < retry-1 {
			// Add jitter to avoid thundering herd
			jitter := time.Duration(rand.Int63n(int64(interval)))
			time.Sleep(interval + jitter)
		}
	}

	return false, nil
}
