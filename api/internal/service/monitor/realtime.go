// Package monitor provides real-time system monitoring services.
package monitor

import (
	"context"
	"runtime"
	"sync"
	"time"

	"gamelink/internal/ws"

	"github.com/shirou/gopsutil/v3/cpu"

	"gorm.io/gorm"
)

// RealtimeService provides real-time system monitoring functionality.
type RealtimeService struct {
	hub              *ws.Hub
	db               *gorm.DB
	startTime        time.Time
	peakOnline       int
	mu               sync.RWMutex
	stopChan         chan struct{}
	requestCount     int64
	lastRequestCount int64
	lastRequestTime  time.Time
}

// NewRealtimeService creates a new realtime monitoring service.
func NewRealtimeService(hub *ws.Hub, db *gorm.DB) *RealtimeService {
	return &RealtimeService{
		hub:             hub,
		db:              db,
		startTime:       time.Now(),
		stopChan:        make(chan struct{}),
		lastRequestTime: time.Now(),
	}
}

// Start begins the monitoring loop.
func (s *RealtimeService) Start(ctx context.Context) {
	// System status broadcast every 5 seconds
	go s.systemStatusLoop(ctx, 5*time.Second)

	// Online users broadcast every 10 seconds
	go s.onlineUsersLoop(ctx, 10*time.Second)

	// Order queue broadcast every 15 seconds
	go s.orderQueueLoop(ctx, 15*time.Second)
}

// Stop stops the monitoring service.
func (s *RealtimeService) Stop() {
	close(s.stopChan)
}

// systemStatusLoop periodically broadcasts system status.
func (s *RealtimeService) systemStatusLoop(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopChan:
			return
		case <-ticker.C:
			status := s.GetSystemStatus()
			msg, err := ws.NewSystemStatusMessage(status)
			if err == nil {
				s.hub.Broadcast(msg)
			}
		}
	}
}

// onlineUsersLoop periodically broadcasts online user stats.
func (s *RealtimeService) onlineUsersLoop(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopChan:
			return
		case <-ticker.C:
			users := s.GetOnlineUsers()
			msg, err := ws.NewOnlineUsersMessage(users)
			if err == nil {
				s.hub.Broadcast(msg)
			}
		}
	}
}

// orderQueueLoop periodically broadcasts order queue status.
func (s *RealtimeService) orderQueueLoop(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopChan:
			return
		case <-ticker.C:
			queue := s.GetOrderQueue(ctx)
			msg, err := ws.NewOrderQueueMessage(queue)
			if err == nil {
				s.hub.Broadcast(msg)
			}
		}
	}
}

// GetSystemStatus returns current system status.
func (s *RealtimeService) GetSystemStatus() *ws.SystemStatus {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Calculate request rate
	s.mu.Lock()
	now := time.Now()
	elapsed := now.Sub(s.lastRequestTime).Seconds()
	requestsPerSec := float64(0)
	if elapsed > 0 {
		requestsPerSec = float64(s.requestCount-s.lastRequestCount) / elapsed
	}
	s.lastRequestCount = s.requestCount
	s.lastRequestTime = now
	s.mu.Unlock()

	// Get DB connection stats
	dbConn := ws.DBConnections{
		Active: 0,
		Idle:   0,
		Max:    50,
	}
	if s.db != nil {
		sqlDB, err := s.db.DB()
		if err == nil {
			stats := sqlDB.Stats()
			dbConn.Active = stats.InUse
			dbConn.Idle = stats.Idle
			dbConn.Max = stats.MaxOpenConnections
		}
	}

	// Calculate memory usage percentage
	memUsage := float64(memStats.Alloc) / float64(memStats.Sys) * 100

	// Get CPU usage - 使用 gopsutil
	cpuPercent, _ := cpu.Percent(0, false)
	var cpuUsage float64
	if len(cpuPercent) > 0 {
		cpuUsage = cpuPercent[0]
	}

	// Determine system status
	status := "healthy"
	if (memUsage > 90 || cpuUsage > 90) || dbConn.Active > dbConn.Max*8/10 {
		status = "critical"
	} else if (memUsage > 70 || cpuUsage > 70) || dbConn.Active > dbConn.Max*6/10 {
		status = "degraded"
	}

	return &ws.SystemStatus{
		CPUUsage:       cpuUsage,
		MemoryUsage:    memUsage,
		MemoryTotal:    memStats.Sys,
		MemoryUsed:     memStats.Alloc,
		Goroutines:     runtime.NumGoroutine(),
		DBConnections:  dbConn,
		Uptime:         int64(time.Since(s.startTime).Seconds()),
		RequestsPerSec: requestsPerSec,
		Status:         status,
	}
}

// GetOnlineUsers returns current online user statistics.
func (s *RealtimeService) GetOnlineUsers() *ws.OnlineUsers {
	current := s.hub.GetOnlineCount()
	byRole := s.hub.GetOnlineCountByRole()

	s.mu.Lock()
	if current > s.peakOnline {
		s.peakOnline = current
	}
	peak := s.peakOnline
	s.mu.Unlock()

	return &ws.OnlineUsers{
		Total:     current,
		Peak:      peak,
		ByRole:    byRole,
		UpdatedAt: time.Now(),
	}
}

// GetOrderQueue returns current order queue status.
func (s *RealtimeService) GetOrderQueue(ctx context.Context) *ws.OrderQueue {
	queue := &ws.OrderQueue{}

	if s.db == nil {
		return queue
	}

	// Count pending orders
	var pendingCount int64
	s.db.WithContext(ctx).Table("orders").
		Where("status = ?", "pending").
		Count(&pendingCount)
	queue.Pending = int(pendingCount)

	// Count processing orders
	var processingCount int64
	s.db.WithContext(ctx).Table("orders").
		Where("status = ?", "in_progress").
		Count(&processingCount)
	queue.Processing = int(processingCount)

	// Count completed orders in last hour (for speed calculation)
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	var completedCount int64
	s.db.WithContext(ctx).Table("orders").
		Where("status = ? AND updated_at >= ?", "completed", oneHourAgo).
		Count(&completedCount)

	queue.Completed = int(completedCount)
	queue.ProcessingSpeed = float64(completedCount) / 60.0 // per minute

	// Calculate average wait time (simplified)
	queue.AverageWaitTime = 0
	if queue.ProcessingSpeed > 0 {
		queue.AverageWaitTime = float64(queue.Pending) / queue.ProcessingSpeed * 60 // seconds
	}

	// Determine if there's a backlog
	queue.HasBacklog = queue.Pending > 100 || queue.AverageWaitTime > 3600 // > 1 hour wait

	return queue
}

// IncrementRequestCount increments the request counter.
func (s *RealtimeService) IncrementRequestCount() {
	s.mu.Lock()
	s.requestCount++
	s.mu.Unlock()
}

// BroadcastAlert sends an alert to all connected clients.
func (s *RealtimeService) BroadcastAlert(alert *ws.Alert) error {
	msg, err := ws.NewAlertMessage(alert)
	if err != nil {
		return err
	}
	s.hub.Broadcast(msg)
	return nil
}
