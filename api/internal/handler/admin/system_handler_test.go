package admin

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/service/integration"
	"gamelink/pkg/cache"
	"gamelink/pkg/config"
)

// TestSystemHandler_InitStatus_NotInitialized tests the init status endpoint when system is not initialized
func TestSystemHandler_InitStatus_NotInitialized(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	db := integration.SetupTestDB(t)

	// Create handler
	cfg := config.AppConfig{
		Database: config.DatabaseConfig{
			Type: "postgres",
		},
		Cache: config.CacheConfig{
			Type: "memory",
		},
	}
	sqlDB, _ := db.DB()
	cacheClient := cache.NewMemory()

	handler := NewSystemInfoHandler(cfg, sqlDB, cacheClient, db)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/system/init-status", handler.InitStatus)

	// Make request
	req, _ := http.NewRequest("GET", "/system/init-status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"initialized":false`)
}

// TestSystemHandler_InitStatus_Initialized tests the init status endpoint when system is initialized
func TestSystemHandler_InitStatus_Initialized(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	db := integration.SetupTestDB(t)

	// Create test permission
	perm := &model.Permission{
		Method:      model.HTTPMethodGET,
		Path:        "/api/v1/admin/test",
		Code:        "admin.test",
		Group:       "admin",
		Description: "Test permission",
		IsSystem:    false,
	}
	require.NoError(t, db.Create(perm).Error)

	// Create test menu
	menu := &model.Menu{
		Name:     "Test Menu",
		Path:     "/test",
		Icon:     "test",
		Order: 1,
		Hidden:   false,
	}
	require.NoError(t, db.Create(menu).Error)

	// Create handler
	cfg := config.AppConfig{
		Database: config.DatabaseConfig{
			Type: "postgres",
		},
		Cache: config.CacheConfig{
			Type: "memory",
		},
	}
	sqlDB, _ := db.DB()
	cacheClient := cache.NewMemory()

	handler := NewSystemInfoHandler(cfg, sqlDB, cacheClient, db)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/system/init-status", handler.InitStatus)

	// Make request
	req, _ := http.NewRequest("GET", "/system/init-status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"initialized":true`)
	assert.Contains(t, w.Body.String(), `"menuCount":1`)
	assert.Contains(t, w.Body.String(), `"permissionCount":1`)
}

// TestSystemHandler_Config tests the system config endpoint
func TestSystemHandler_Config(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	db := integration.SetupTestDB(t)

	cfg := config.AppConfig{
		Database: config.DatabaseConfig{
			Type: "postgres",
		},
		Cache: config.CacheConfig{
			Type: "redis",
		},
		AdminAuth: config.AdminAuthConfig{
			Mode: "jwt",
		},
	}
	sqlDB, _ := db.DB()
	cacheClient := cache.NewMemory()

	handler := NewSystemInfoHandler(cfg, sqlDB, cacheClient, db)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/system/config", handler.Config)

	// Make request
	req, _ := http.NewRequest("GET", "/system/config", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"databaseType":"postgres"`)
	assert.Contains(t, w.Body.String(), `"cacheType":"redis"`)
	assert.Contains(t, w.Body.String(), `"adminAuthMode":"jwt"`)
}

// TestSystemHandler_DBStatus tests the database status endpoint
func TestSystemHandler_DBStatus(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	db := integration.SetupTestDB(t)

	cfg := config.AppConfig{}
	sqlDB, _ := db.DB()
	cacheClient := cache.NewMemory()

	handler := NewSystemInfoHandler(cfg, sqlDB, cacheClient, db)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/system/db", handler.DBStatus)

	// Make request
	req, _ := http.NewRequest("GET", "/system/db", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"openConnections"`)
	assert.Contains(t, w.Body.String(), `"inUse"`)
	assert.Contains(t, w.Body.String(), `"idle"`)
	assert.Contains(t, w.Body.String(), `"maxOpenConns"`)
}

// TestSystemHandler_CacheStatus tests the cache status endpoint
func TestSystemHandler_CacheStatus(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	db := integration.SetupTestDB(t)

	cfg := config.AppConfig{
		Cache: config.CacheConfig{
			Type: "memory",
		},
	}
	sqlDB, _ := db.DB()
	cacheClient := cache.NewMemory()

	handler := NewSystemInfoHandler(cfg, sqlDB, cacheClient, db)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/system/cache", handler.CacheStatus)

	// Make request
	req, _ := http.NewRequest("GET", "/system/cache", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"connected":true`)
	assert.Contains(t, w.Body.String(), `"type":"memory"`)
}

// TestSystemHandler_Resources tests the system resources endpoint
func TestSystemHandler_Resources(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	db := integration.SetupTestDB(t)

	cfg := config.AppConfig{}
	sqlDB, _ := db.DB()
	cacheClient := cache.NewMemory()

	handler := NewSystemInfoHandler(cfg, sqlDB, cacheClient, db)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/system/resources", handler.Resources)

	// Make request
	req, _ := http.NewRequest("GET", "/system/resources", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"goroutines"`)
	assert.Contains(t, w.Body.String(), `"allocMB"`)
	assert.Contains(t, w.Body.String(), `"totalAllocMB"`)
	assert.Contains(t, w.Body.String(), `"sysMB"`)
	assert.Contains(t, w.Body.String(), `"numGC"`)
	assert.Contains(t, w.Body.String(), `"cpuCores"`)
}

// TestSystemHandler_Version tests the system version endpoint
func TestSystemHandler_Version(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	db := integration.SetupTestDB(t)

	cfg := config.AppConfig{}
	sqlDB, _ := db.DB()
	cacheClient := cache.NewMemory()

	handler := NewSystemInfoHandler(cfg, sqlDB, cacheClient, db)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/system/version", handler.Version)

	// Make request
	req, _ := http.NewRequest("GET", "/system/version", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"version":"1.0.0"`)
	assert.Contains(t, w.Body.String(), `"goVersion"`)
	assert.Contains(t, w.Body.String(), `"buildTime":"2024-11-06"`)
}

// TestSystemHandler_InitStatus_WithSystemState tests init status with system state record
func TestSystemHandler_InitStatus_WithSystemState(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	db := integration.SetupTestDB(t)

	// Create test permission and menu
	perm := &model.Permission{
		Method:      model.HTTPMethodGET,
		Path:        "/api/v1/admin/test",
		Code:        "admin.test",
		Group:       "admin",
		Description: "Test permission",
		IsSystem:    false,
	}
	require.NoError(t, db.Create(perm).Error)

	menu := &model.Menu{
		Name:     "Test Menu",
		Path:     "/test",
		Icon:     "test",
		Order: 1,
		Hidden:   false,
	}
	require.NoError(t, db.Create(menu).Error)

	// Create system state record
	state := model.SystemState{
		Key:       model.SystemStateKeyAdminInit,
		Version:   "1.0.0",
		SyncedBy:  1,
		LastSyncAt: time.Now(),
	}
	require.NoError(t, db.Create(&state).Error)

	// Create handler
	cfg := config.AppConfig{}
	sqlDB, _ := db.DB()
	cacheClient := cache.NewMemory()

	handler := NewSystemInfoHandler(cfg, sqlDB, cacheClient, db)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/system/init-status", handler.InitStatus)

	// Make request
	req, _ := http.NewRequest("GET", "/system/init-status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"initialized":true`)
	assert.Contains(t, w.Body.String(), `"version":"1.0.0"`)
	assert.Contains(t, w.Body.String(), `"syncedBy":1`)
}

// TestNewSystemInfoHandler tests the handler constructor
func TestNewSystemInfoHandler(t *testing.T) {
	cfg := config.AppConfig{}
	var sqlDB *sql.DB
	cacheClient := cache.NewMemory()
	var db *gorm.DB

	handler := NewSystemInfoHandler(cfg, sqlDB, cacheClient, db)

	assert.NotNil(t, handler)
	assert.Equal(t, cfg, handler.cfg)
	assert.Equal(t, sqlDB, handler.sqlDB)
	assert.Equal(t, cacheClient, handler.cacheClient)
	assert.Equal(t, db, handler.db)
}
