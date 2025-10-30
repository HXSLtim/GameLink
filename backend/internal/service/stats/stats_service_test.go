package stats

import (
	"testing"

	"gamelink/internal/repository"
)

// mockStatsRepository 是一个简单的mock实现。
type mockStatsRepository struct {
	repository.StatsRepository
}

// TestNewStatsService 测试构造函数。
func TestNewStatsService(t *testing.T) {
	repo := &mockStatsRepository{}

	svc := NewStatsService(repo)

	if svc == nil {
		t.Fatal("NewStatsService returned nil")
	}

	if svc.repo != repo {
		t.Error("repo not set correctly")
	}
}
