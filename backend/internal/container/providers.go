package container

import (
	"context"
	"database/sql"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"

	"gamelink/internal/cache"
	"gamelink/internal/config"
	"gamelink/internal/db"
	"gamelink/internal/lifecycle"
	"gamelink/internal/repository/common"
	gamerepo "gamelink/internal/repository/game"
	orderrepo "gamelink/internal/repository/implementations"
	paymentrepo "gamelink/internal/repository/payment"
	playerrepo "gamelink/internal/repository/player"
	rolerepo "gamelink/internal/repository/role"
	serviceitemrepo "gamelink/internal/repository/serviceitem"
	userrepo "gamelink/internal/repository/user"
	"gamelink/internal/router"
	adminservice "gamelink/internal/service/admin"
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
	return db.Open(cfg)
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
	client, err := cache.New(cfg.Cache)
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
		playerrepo.NewPlayerRepository(orm),
		orderrepo.NewOrderRepository(orm),
		paymentrepo.NewPaymentRepository(orm),
		rolerepo.NewRoleRepository(orm),
		serviceitemrepo.NewServiceItemRepository(orm),
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

// ProvideEngine builds the gin engine by invoking Router.Setup().
func ProvideEngine(r *router.Router) *gin.Engine {
	return r.Setup()
}
