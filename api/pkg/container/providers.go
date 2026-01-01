package container

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"

	menurepo "gamelink/internal/repository/admin"
	permissionrepo "gamelink/internal/repository/admin"
	rolerepo "gamelink/internal/repository/admin"
	"gamelink/internal/repository/common"
	gamerepo "gamelink/internal/repository/game"
	gamecategoryrepo "gamelink/internal/repository/gamecategory"
	orderrepo "gamelink/internal/repository/implementations"
	paymentrepo "gamelink/internal/repository/order"
	serviceitemrepo "gamelink/internal/repository/serviceitem"
	statsrepo "gamelink/internal/repository/stats"
	userrepo "gamelink/internal/repository/user"
	walletrepo "gamelink/internal/repository/wallet"
	"gamelink/internal/router"
	adminservice "gamelink/internal/service/admin"
	"gamelink/pkg/cache"
	"gamelink/pkg/config"
	"gamelink/pkg/db"
	"gamelink/pkg/lifecycle"
	"gamelink/pkg/metrics"
)

// ProviderSet is consumed by wire to build the dependency graph.
var ProviderSet = wire.NewSet(
	ProvideConfig,
	lifecycle.NewRegistry,
	ProvideDB,
	ProvideSQLDB,
	ProvideCache,
	ProvideAdminService,
	ProvideRouter,
	ProvideEngine,
	lifecycle.NewManager,
	wire.Struct(new(Application), "Engine", "Config", "Lifecycle"),
)

// ProvideConfig loads and validates the application configuration.
func ProvideConfig() (config.AppConfig, error) {
	cfg := config.Load()
	if err := config.Validate(os.Getenv("APP_ENV"), cfg); err != nil {
		return config.AppConfig{}, err
	}
	return cfg, nil
}

// ProvideDB opens the primary gorm DB connection.
func ProvideDB(cfg config.AppConfig) (*gorm.DB, error) {
	// Create config provider adapter
	configProvider := &configProviderAdapter{cfg: cfg}
	metricsProvider := &metricsProviderAdapter{}
	return db.Open(configProvider, metricsProvider)
}

// ProvideSQLDB exposes the underlying *sql.DB and registers its shutdown hook.
func ProvideSQLDB(reg *lifecycle.Registry, orm *gorm.DB) (*sql.DB, error) {
	sqlDB, err := orm.DB()
	if err != nil {
		return nil, err
	}

	reg.RegisterHook("database", nil, func(context.Context) error {
		return sqlDB.Close()
	})

	return sqlDB, nil
}

// ProvideCache initializes the cache backend and registers its shutdown hook.
func ProvideCache(cfg config.AppConfig, reg *lifecycle.Registry) (cache.Cache, error) {
	// Create cache config provider adapter
	cacheConfigProvider := &cacheConfigProviderAdapter{cfg: cfg}
	client, err := cache.New(cacheConfigProvider)
	if err != nil {
		return nil, err
	}

	reg.RegisterHook("cache", nil, func(ctx context.Context) error {
		return client.Close(ctx)
	})

	return client, nil
}

// ProvideAdminService wires the admin service and its transactional manager.
func ProvideAdminService(orm *gorm.DB, cacheClient cache.Cache) *adminservice.AdminService {
	svc := adminservice.NewAdminService(
		gamerepo.NewGameRepository(orm),
		userrepo.NewUserRepository(orm),
		userrepo.NewPlayerRepository(orm),
		orderrepo.NewOrderRepository(orm),
		paymentrepo.NewPaymentRepository(orm),
		rolerepo.NewRoleRepository(orm),
		serviceitemrepo.NewServiceItemRepository(orm),
		permissionrepo.NewPermissionRepository(orm),
		menurepo.NewMenuRepository(orm),
		statsrepo.NewStatsRepository(orm),
		walletrepo.NewWalletRepository(orm),
		gamecategoryrepo.NewGameCategoryRepository(orm),
		cacheClient,
	)
	svc.SetTxManager(common.NewUnitOfWork(orm))
	return svc
}

// ProvideRouter constructs the HTTP router abstraction.
func ProvideRouter(
	cfg config.AppConfig,
	orm *gorm.DB,
	sqlDB *sql.DB,
	cacheClient cache.Cache,
	adminSvc *adminservice.AdminService,
	reg *lifecycle.Registry,
) *router.Router {
	return router.NewRouter(cfg, orm, sqlDB, cacheClient, adminSvc, reg)
}

// Adapter types for pkg/db interfaces
type configProviderAdapter struct {
	cfg config.AppConfig
}

func (a *configProviderAdapter) GetDatabaseDSN() string {
	return a.cfg.Database.DSN
}

func (a *configProviderAdapter) IsSeedEnabled() bool {
	return a.cfg.Seed.Enabled
}

func (a *configProviderAdapter) GetMaxOpenConns() int {
	return 1 // Default for SQLite, can be made configurable
}

func (a *configProviderAdapter) GetDatabaseType() string {
	return a.cfg.Database.Type
}

type metricsProviderAdapter struct{}

func (a *metricsProviderAdapter) InstrumentGorm(db interface{}) error {
	if gormDB, ok := db.(*gorm.DB); ok {
		return metrics.InstrumentGorm(gormDB)
	}
	return fmt.Errorf("expected *gorm.DB, got %T", db)
}

// Cache config provider adapter
type cacheConfigProviderAdapter struct {
	cfg config.AppConfig
}

func (a *cacheConfigProviderAdapter) GetCacheType() string {
	return a.cfg.Cache.Type
}

func (a *cacheConfigProviderAdapter) GetRedisConfig() cache.RedisConfig {
	return cache.RedisConfig{
		Addr:     a.cfg.Cache.Redis.Addr,
		Password: a.cfg.Cache.Redis.Password,
		DB:       a.cfg.Cache.Redis.DB,
	}
}

// ProvideEngine builds the gin engine by invoking Router.Setup().
func ProvideEngine(r *router.Router) *gin.Engine {
	return r.Setup()
}
