package admin

import (
	"database/sql"
	"runtime"

	"github.com/gin-gonic/gin"

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
}

// NewSystemInfoHandler 创建系统信息Handler
func NewSystemInfoHandler(cfg config.AppConfig, sqlDB *sql.DB, cacheClient cache.Cache) *SystemInfoHandler {
	return &SystemInfoHandler{
		cfg:         cfg,
		sqlDB:       sqlDB,
		cacheClient: cacheClient,
	}
}

// RegisterSystemRoutes 注册系统信息路由
func RegisterSystemRoutes(router gin.IRouter, cfg config.AppConfig, sqlDB *sql.DB, cacheClient cache.Cache, pm *mw.PermissionMiddleware) {
	h := NewSystemInfoHandler(cfg, sqlDB, cacheClient)
	group := router.Group("/system")
	{
		group.GET("/config", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/system/config"), h.Config)
		group.GET("/db", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/system/db"), h.DBStatus)
		group.GET("/cache", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/system/cache"), h.CacheStatus)
		group.GET("/resources", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/system/resources"), h.Resources)
		group.GET("/version", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/system/version"), h.Version)
	}
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
