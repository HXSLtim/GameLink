package admin

import (
	"database/sql"
	"fmt"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	mw "gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	"gamelink/pkg/cache"
	"gamelink/pkg/config"
)

// SystemInfoHandler 系统信息Handler
type SystemInfoHandler struct {
	cfg         config.AppConfig
	sqlDB       *sql.DB
	cacheClient cache.Cache
	db          *gorm.DB
}

// NewSystemInfoHandler 创建系统信息Handler
func NewSystemInfoHandler(cfg config.AppConfig, sqlDB *sql.DB, cacheClient cache.Cache, db *gorm.DB) *SystemInfoHandler {
	return &SystemInfoHandler{
		cfg:         cfg,
		sqlDB:       sqlDB,
		cacheClient: cacheClient,
		db:          db,
	}
}

// RegisterSystemRoutes 注册系统信息路由
func RegisterSystemRoutes(router gin.IRouter, cfg config.AppConfig, sqlDB *sql.DB, cacheClient cache.Cache, db *gorm.DB, pm *mw.PermissionMiddleware) {
	h := NewSystemInfoHandler(cfg, sqlDB, cacheClient, db)
	group := router.Group("/system")
	{
		group.GET("/config", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/system/config"), h.Config)
		group.GET("/db", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/system/db"), h.DBStatus)
		group.GET("/cache", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/system/cache"), h.CacheStatus)
		group.GET("/resources", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/system/resources"), h.Resources)
		group.GET("/version", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/system/version"), h.Version)
		group.GET("/init-status", h.InitStatus) // No permission required - used for initialization check
	}
}

// InitStatusResponse represents the initialization status response
type InitStatusResponse struct {
	Initialized    bool      `json:"initialized"`
	LastSyncAt     *time.Time `json:"lastSyncAt,omitempty"`
	MenuCount      int       `json:"menuCount,omitempty"`
	PermissionCount int      `json:"permissionCount,omitempty"`
	Version        string    `json:"version,omitempty"`
	SyncedBy       uint64    `json:"syncedBy,omitempty"`
	Message        string    `json:"message,omitempty"`
}

// InitStatus 获取系统初始化状态
// @Summary      获取系统初始化状态
// @Description  检查系统是否已初始化（菜单和权限同步）
// @Tags         Admin - System
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Success      200            {object}  model.SuccessResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/system/init-status [get]
func (h *SystemInfoHandler) InitStatus(c *gin.Context) {
	response := InitStatusResponse{
		Initialized: false,
		Message:     "System not initialized",
	}

	// 判断准则：检查数据库内是否有权限和菜单记录
	var permCount int64
	var menuCount int64

	// 统计所有权限数量（不再过滤，直接统计）
	if err := h.db.Model(&model.Permission{}).Count(&permCount).Error; err != nil {
		respondError(c, fmt.Errorf("failed to check permissions: %w", err))
		return
	}

	// 统计菜单数量
	if err := h.db.Model(&model.Menu{}).Count(&menuCount).Error; err != nil {
		respondError(c, fmt.Errorf("failed to check menus: %w", err))
		return
	}

	// 真正的判断标准：有业务权限和菜单记录
	if permCount > 0 && menuCount > 0 {
		response.Initialized = true
		response.MenuCount = int(menuCount)
		response.PermissionCount = int(permCount)
		response.Message = "System initialized"

		// 尝试从 system_states 获取额外的同步信息（如果存在）
		var state model.SystemState
		if err := h.db.Where("key = ?", model.SystemStateKeyAdminInit).First(&state).Error; err == nil {
			response.LastSyncAt = &state.LastSyncAt
			response.Version = state.Version
			response.SyncedBy = state.SyncedBy
		} else {
			// 如果没有同步记录，尝试从最近更新的权限或菜单获取时间
			var latestPerm model.Permission
			var latestMenu model.Menu
			var latestTime time.Time

			if err := h.db.Order("updated_at DESC").First(&latestPerm).Error; err == nil {
				latestTime = latestPerm.UpdatedAt
			}
			if err := h.db.Order("updated_at DESC").First(&latestMenu).Error; err == nil {
				if latestMenu.UpdatedAt.After(latestTime) {
					latestTime = latestMenu.UpdatedAt
				}
			}
			if !latestTime.IsZero() {
				response.LastSyncAt = &latestTime
			}
		}
	} else {
		response.Message = fmt.Sprintf("System not initialized - permissions: %d, menus: %d", permCount, menuCount)
	}

	respondSuccess(c, response)
}

// Config 获取系统配置
// @Summary      系统配置
// @Description  获取系统配置信息
// @Tags         Admin - System
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Success      200            {object}  model.SuccessResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/system/config [get]
func (h *SystemInfoHandler) Config(c *gin.Context) {
	respondSuccess(c, gin.H{
		"databaseType":  h.cfg.Database.Type,
		"cacheType":     h.cfg.Cache.Type,
		"adminAuthMode": h.cfg.AdminAuth.Mode,
	})
}

// DBStatus 获取数据库状态
// @Summary      数据库状态
// @Description  获取数据库连接状态
// @Tags         Admin - System
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Success      200            {object}  model.SuccessResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/system/db [get]
func (h *SystemInfoHandler) DBStatus(c *gin.Context) {
	stats := h.sqlDB.Stats()
	respondSuccess(c, gin.H{
		"openConnections": stats.OpenConnections,
		"inUse":           stats.InUse,
		"idle":            stats.Idle,
		"maxOpenConns":    stats.MaxOpenConnections,
	})
}

// CacheStatus 获取缓存状态
// @Summary      缓存状态
// @Description  获取缓存连接状态
// @Tags         Admin - System
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Success      200            {object}  model.SuccessResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/system/cache [get]
func (h *SystemInfoHandler) CacheStatus(c *gin.Context) {
	testKey := "system:health:check"
	_, _, err := h.cacheClient.Get(c.Request.Context(), testKey)
	respondSuccess(c, gin.H{
		"connected": err == nil,
		"type":      h.cfg.Cache.Type,
	})
}

// Resources 获取系统资源信息
// @Summary      系统资源
// @Description  获取系统资源使用情况
// @Tags         Admin - System
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Success      200            {object}  model.SuccessResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/system/resources [get]
func (h *SystemInfoHandler) Resources(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	respondSuccess(c, gin.H{
		"goroutines":   runtime.NumGoroutine(),
		"allocMB":      m.Alloc / 1024 / 1024,
		"totalAllocMB": m.TotalAlloc / 1024 / 1024,
		"sysMB":        m.Sys / 1024 / 1024,
		"numGC":        m.NumGC,
		"cpuCores":     runtime.NumCPU(),
	})
}

// Version 获取系统版本
// @Summary      系统版本
// @Description  获取系统版本信息
// @Tags         Admin - System
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Success      200            {object}  model.SuccessResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/system/version [get]
func (h *SystemInfoHandler) Version(c *gin.Context) {
	respondSuccess(c, gin.H{
		"version":   "1.0.0",
		"goVersion": runtime.Version(),
		"buildTime": "2024-11-06",
	})
}
