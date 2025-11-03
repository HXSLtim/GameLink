package admin

import (
    "database/sql"
    "os"
    "strings"

    "github.com/gin-gonic/gin"

    "gamelink/internal/cache"
    "gamelink/internal/config"
    mw "gamelink/internal/handler/middleware"
    "gamelink/internal/model"
)

// RegisterSystemRoutes 注册系统信息相关路由（管理端）
// 使用细粒度权限控制（method+path 级别）
func RegisterSystemRoutes(router gin.IRouter, cfg config.AppConfig, sqlDB *sql.DB, cacheClient cache.Cache, pm *mw.PermissionMiddleware) {
    h := NewSystemInfoHandler(cfg, sqlDB, cacheClient)

    group := router.Group("/admin")
    // 系统信息接口均需要认证 + 速率限制
    if os.Getenv("APP_ENV") == "production" {
        group.Use(pm.RequireAuth(), mw.RateLimitAdmin())
    } else {
        // 使用配置中的 admin_auth.mode
        switch strings.ToLower(cfg.AdminAuth.Mode) {
        case "jwt":
            group.Use(pm.RequireAuth(), mw.RateLimitAdmin())
        default:
            group.Use(mw.AdminAuth(), mw.RateLimitAdmin())
        }
    }

    // 系统信息接口 - 使用细粒度权限
    group.GET("/system/config", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/system/config"), h.Config)
    group.GET("/system/db", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/system/db"), h.DBStatus)
    group.GET("/system/cache", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/system/cache"), h.CacheStatus)
    group.GET("/system/resources", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/system/resources"), h.Resources)
    group.GET("/system/version", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/system/version"), h.Version)
}
