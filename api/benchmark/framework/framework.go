package benchmark

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/user"
	svcauth "gamelink/internal/service/auth"
	svcpayment "gamelink/internal/service/payment"
	authpkg "gamelink/pkg/auth"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// BenchmarkConfig holds configuration for benchmarks
type BenchmarkConfig struct {
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	EnableProfiling bool
	CPUProfile      string
	MemProfile      string
}

// DefaultConfig returns default benchmark configuration
func DefaultConfig() *BenchmarkConfig {
	return &BenchmarkConfig{
		DBHost:          getEnv("BENCH_DB_HOST", "localhost"),
		DBPort:          getEnv("BENCH_DB_PORT", "5432"),
		DBUser:          getEnv("BENCH_DB_USER", "gamelink"),
		DBPassword:      getEnv("BENCH_DB_PASSWORD", "gamelink"),
		DBName:          getEnv("BENCH_DB_NAME", "gamelink_bench"),
		EnableProfiling: false,
	}
}

// BenchmarkSuite holds all services and repositories for benchmarking
type BenchmarkSuite struct {
	DB             *gorm.DB
	AuthService    *svcauth.AuthService
	PaymentService *svcpayment.PaymentService
	UserRepo       repository.UserRepository
	PlayerRepo     repository.PlayerRepository
	GameRepo       repository.GameRepository
	Config         *BenchmarkConfig
	CleanupFunc    func()
}

// NewBenchmarkSuite creates a new benchmark suite
func NewBenchmarkSuite(t testing.TB, config *BenchmarkConfig) *BenchmarkSuite {
	if config == nil {
		config = DefaultConfig()
	}

	// Setup database
	db, cleanup := setupBenchmarkDB(t, config)

	// Create repositories
	userRepo := user.NewUserRepository(db)

	// Create JWT manager for auth
	jwtManager := authpkg.NewJWTManager("test-secret-key-for-benchmark", 24*time.Hour)

	// Create services
	authService := svcauth.NewAuthService(userRepo, jwtManager)

	suite := &BenchmarkSuite{
		DB:          db,
		AuthService: authService,
		UserRepo:    userRepo,
		Config:      config,
		CleanupFunc: cleanup,
	}

	// Setup cleanup on test completion
	t.Cleanup(func() {
		suite.CleanupFunc()
	})

	return suite
}

// setupBenchmarkDB creates a database connection for benchmarking
func setupBenchmarkDB(t testing.TB, config *BenchmarkConfig) (*gorm.DB, func()) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Silent logging for benchmarks
	})
	if err != nil {
		t.Fatalf("Failed to connect to benchmark database: %v", err)
	}

	// Test connection
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get database instance: %v", err)
	}

	// Configure connection pool for high-performance benchmarking
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("Failed to ping benchmark database: %v", err)
	}

	log.Printf("Connected to benchmark database: %s", config.DBName)

	cleanup := func() {
		log.Printf("Closing benchmark database connection")
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}

	return db, cleanup
}

// CreateBenchmarkUser creates a test user for benchmarking
func (s *BenchmarkSuite) CreateBenchmarkUser(t testing.TB, phone string) *model.User {
	user := &model.User{
		Phone:        phone,
		Name:         "Benchmark User",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "$2a$10$benchmark.hash.for.testing",
	}
	if err := s.DB.Create(user).Error; err != nil {
		t.Fatalf("Failed to create benchmark user: %v", err)
	}
	return user
}

// CreateBenchmarkPlayer creates a test player for benchmarking
func (s *BenchmarkSuite) CreateBenchmarkPlayer(t testing.TB, userID uint64, nickname string) *model.Player {
	player := &model.Player{
		UserID:          userID,
		Nickname:        nickname,
		Rank:            "bronze",
		RatingAverage:   5.0,
		RatingCount:     0,
		HourlyRateCents: 2000, // 20 CNY
		OnlineStatus:    model.PlayerOnlineStatusOnline,
		AcceptingOrders: true,
	}
	if err := s.DB.Create(player).Error; err != nil {
		t.Fatalf("Failed to create benchmark player: %v", err)
	}
	return player
}

// CreateBenchmarkGame creates a test game for benchmarking
func (s *BenchmarkSuite) CreateBenchmarkGame(t testing.TB, name string) *model.Game {
	game := &model.Game{
		Key:         "bench-" + name,
		Name:        name,
		Description: "Benchmark test game",
		IsActive:    true,
		IconURL:     "https://example.com/icon.png",
	}
	if err := s.DB.Create(game).Error; err != nil {
		t.Fatalf("Failed to create benchmark game: %v", err)
	}
	return game
}

// CleanBenchmarkData removes all benchmark data
func (s *BenchmarkSuite) CleanBenchmarkData(t testing.TB) {
	// Clean in order of dependencies
	tables := []string{
		"order_players", "payments", "reviews", "disputes",
		"orders", "players", "users", "games",
	}
	for _, table := range tables {
		if err := s.DB.Exec("DELETE FROM " + table).Error; err != nil {
			log.Printf("Warning: Failed to clean table %s: %v", table, err)
		}
	}
}

// ResetBenchmarkSequence resets database sequences
func (s *BenchmarkSuite) ResetBenchmarkSequence(t testing.TB) {
	// Reset sequences to avoid ID exhaustion during benchmarks
	sequences := []string{
		"users_id_seq", "players_id_seq", "games_id_seq",
		"orders_id_seq", "payments_id_seq",
	}
	for _, seq := range sequences {
		if err := s.DB.Exec(fmt.Sprintf("SELECT setval('%s', 1, false)", seq)).Error; err != nil {
			log.Printf("Warning: Failed to reset sequence %s: %v", seq, err)
		}
	}
}

// BenchmarkTimer measures execution time for operations
type BenchmarkTimer struct {
	start time.Time
	name  string
}

// StartTimer starts a new benchmark timer
func StartTimer(name string) *BenchmarkTimer {
	return &BenchmarkTimer{
		start: time.Now(),
		name:  name,
	}
}

// End stops the timer and reports the duration
func (bt *BenchmarkTimer) End() time.Duration {
	duration := time.Since(bt.start)
	log.Printf("%s took %v", bt.name, duration)
	return duration
}

// ReportMetrics reports benchmark metrics
func ReportMetrics(b *testing.B, opName string, durations []time.Duration) {
	if len(durations) == 0 {
		return
	}

	var total time.Duration
	var min, max time.Duration = durations[0], durations[0]

	for _, d := range durations {
		total += d
		if d < min {
			min = d
		}
		if d > max {
			max = d
		}
	}

	avg := total / time.Duration(len(durations))

	b.ReportMetric(float64(avg.Nanoseconds()), "ns/op")
	b.ReportMetric(float64(min.Nanoseconds()), "ns/min")
	b.ReportMetric(float64(max.Nanoseconds()), "ns/max")

	log.Printf("%s - Avg: %v, Min: %v, Max: %v", opName, avg, min, max)
}

// getEnv retrieves environment variable with default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SetupProfiling enables CPU and memory profiling
func SetupProfiling(config *BenchmarkConfig) (func(), error) {
	if !config.EnableProfiling {
		return func() {}, nil
	}

	var cleanupFuncs []func()

	if config.CPUProfile != "" {
		f, err := os.Create(config.CPUProfile)
		if err != nil {
			return nil, fmt.Errorf("could not create CPU profile: %w", err)
		}
		// Note: In actual usage, you would use runtime/pprof here
		cleanupFuncs = append(cleanupFuncs, func() {
			f.Close()
		})
	}

	if config.MemProfile != "" {
		f, err := os.Create(config.MemProfile)
		if err != nil {
			return nil, fmt.Errorf("could not create memory profile: %w", err)
		}
		// Note: In actual usage, you would use runtime/pprof here
		cleanupFuncs = append(cleanupFuncs, func() {
			f.Close()
		})
	}

	return func() {
		for _, cleanup := range cleanupFuncs {
			cleanup()
		}
	}, nil
}

// CheckDBConnection checks if database connection is alive
func CheckDBConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx)
}

// BenchmarkStats holds statistics from benchmark runs
type BenchmarkStats struct {
	Operations    int
	TotalDuration time.Duration
	AvgDuration   time.Duration
	MinDuration   time.Duration
	MaxDuration   time.Duration
	Throughput    float64 // operations per second
	ErrorCount    int
	DBConnections int
	CacheHits     int
	CacheMisses   int
}

// CalculateStats calculates benchmark statistics
func CalculateStats(operations int, totalDuration time.Duration, durations []time.Duration) *BenchmarkStats {
	stats := &BenchmarkStats{
		Operations:    operations,
		TotalDuration: totalDuration,
		AvgDuration:   totalDuration / time.Duration(operations),
		Throughput:    float64(operations) / totalDuration.Seconds(),
	}

	if len(durations) > 0 {
		stats.MinDuration = durations[0]
		stats.MaxDuration = durations[0]
		for _, d := range durations {
			if d < stats.MinDuration {
				stats.MinDuration = d
			}
			if d > stats.MaxDuration {
				stats.MaxDuration = d
			}
		}
	}

	return stats
}

// PrintStats prints benchmark statistics
func PrintStats(stats *BenchmarkStats, opName string) {
	fmt.Printf("\n=== %s Benchmark Results ===\n", opName)
	fmt.Printf("Operations:      %d\n", stats.Operations)
	fmt.Printf("Total Duration:  %v\n", stats.TotalDuration)
	fmt.Printf("Avg Duration:    %v\n", stats.AvgDuration)
	fmt.Printf("Min Duration:    %v\n", stats.MinDuration)
	fmt.Printf("Max Duration:    %v\n", stats.MaxDuration)
	fmt.Printf("Throughput:      %.2f ops/sec\n", stats.Throughput)
	if stats.ErrorCount > 0 {
		fmt.Printf("Errors:          %d\n", stats.ErrorCount)
	}
	fmt.Printf("===============================\n\n")
}
