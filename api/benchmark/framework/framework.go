package benchmark

import (
	"fmt"
	"log"
	"os"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/user"
	svcauth "gamelink/internal/service/auth"
	authpkg "gamelink/pkg/auth"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type BenchmarkConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func DefaultConfig() *BenchmarkConfig {
	return &BenchmarkConfig{
		DBHost:     getEnv("BENCH_DB_HOST", "localhost"),
		DBPort:     getEnv("BENCH_DB_PORT", "5432"),
		DBUser:     getEnv("BENCH_DB_USER", "gamelink"),
		DBPassword: getEnv("BENCH_DB_PASSWORD", "gamelink"),
		DBName:     getEnv("BENCH_DB_NAME", "gamelink_bench"),
	}
}

type BenchmarkSuite struct {
	DB          *gorm.DB
	AuthService *svcauth.AuthService
	UserRepo    repository.UserRepository
	Config      *BenchmarkConfig
	Cleanup     func()
}

func NewBenchmarkSuite(_ testingTB, config *BenchmarkConfig) *BenchmarkSuite {
	if config == nil {
		config = DefaultConfig()
	}

	db, cleanup := setupBenchmarkDB(config)

	userRepo := user.NewUserRepository(db)

	jwtManager := authpkg.NewJWTManager("test-secret-key-for-benchmark", 24*time.Hour)

	authService := svcauth.NewAuthService(userRepo, jwtManager)

	suite := &BenchmarkSuite{
		DB:          db,
		AuthService: authService,
		UserRepo:    userRepo,
		Config:      config,
		Cleanup:     cleanup,
	}

	return suite
}

func setupBenchmarkDB(config *BenchmarkConfig) (*gorm.DB, func()) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to benchmark database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping benchmark database: %v", err)
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

func (s *BenchmarkSuite) CreateBenchmarkUser(_ testingTB, phone string) *model.User {
	testUser := &model.User{
		Phone:        phone,
		Name:         "Benchmark User",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "$2a$10$benchmark.hash.for.testing",
	}
	if err := s.DB.Create(testUser).Error; err != nil {
		log.Fatalf("Failed to create benchmark user: %v", err)
	}
	return testUser
}

func (s *BenchmarkSuite) CreateBenchmarkPlayer(_ testingTB, userID uint64, nickname string) *model.Player {
	player := &model.Player{
		UserID:          userID,
		Nickname:        nickname,
		Rank:            "bronze",
		RatingAverage:   5.0,
		RatingCount:     0,
		HourlyRateCents: 2000,
		OnlineStatus:    model.PlayerOnlineStatusOnline,
		AcceptingOrders: true,
	}
	if err := s.DB.Create(player).Error; err != nil {
		log.Fatalf("Failed to create benchmark player: %v", err)
	}
	return player
}

func (s *BenchmarkSuite) CreateBenchmarkGame(_ testingTB, name string) *model.Game {
	game := &model.Game{
		Key:         "bench-" + name,
		Name:        name,
		Description: "Benchmark test game",
		IsActive:    true,
		IconURL:     "https://example.com/icon.png",
	}
	if err := s.DB.Create(game).Error; err != nil {
		log.Fatalf("Failed to create benchmark game: %v", err)
	}
	return game
}

func (s *BenchmarkSuite) CleanBenchmarkData(_ testingTB) {
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

func (s *BenchmarkSuite) ResetBenchmarkSequence(_ testingTB) {
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

type BenchmarkTimer struct {
	start time.Time
	name  string
}

func NewBenchmarkTimer(name string) *BenchmarkTimer {
	return &BenchmarkTimer{
		start: time.Now(),
		name:  name,
	}
}

func (bt *BenchmarkTimer) End() time.Duration {
	duration := time.Since(bt.start)
	log.Printf("%s took %v", bt.name, duration)
	return duration
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

type testingTB interface{}

func init() {
	_ = NewBenchmarkSuite
	_ = NewBenchmarkTimer
}
