package admin

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	mw "gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	adminservice "gamelink/internal/service/admin"
	permissionservice "gamelink/internal/service/admin" //nolint:ST1019
	roleservice "gamelink/internal/service/admin"       //nolint:ST1019
	statsservice "gamelink/internal/service/admin"      //nolint:ST1019
	reviewservice "gamelink/internal/service/review"
	"gamelink/internal/service/sensitiveword"
	"gamelink/pkg/config"
)

// RegisterRoutesOption 路由注册选项
type RegisterRoutesOption func(*ReviewHandler)

// WithSensitiveWordService 设置敏感词服务
func WithSensitiveWordService(svc *sensitiveword.SensitiveWordService) RegisterRoutesOption {
	return func(h *ReviewHandler) {
		h.SetSensitiveWordService(svc)
	}
}

// RegisterRoutes 注册后台管理相关路由
// 使用细粒度权限控制（method+path 级别）
func RegisterRoutes(router gin.IRouter, svc *adminservice.AdminService, statsSvc *statsservice.StatsService, pm *mw.PermissionMiddleware, opts ...RegisterRoutesOption) {
	// 先注册同步专用路由（不受限流限制）
	RegisterSyncRoutes(router, svc, pm)
	gameHandler := NewGameHandler(svc)
	userHandler := NewUserHandler(svc)
	userBehaviorHandler := NewUserBehaviorHandler(statsSvc)
	playerHandler := NewPlayerHandler(svc)
	orderHandler := NewOrderHandler(svc)
	paymentHandler := NewPaymentHandler(svc)
	reviewHandler := NewReviewHandler(svc)

	// 应用可选配置
	for _, opt := range opts {
		opt(reviewHandler)
	}

	group := router
	// 所有管理接口均需要认+ 速率限制
	cfg := config.Load()
	if os.Getenv("APP_ENV") == "production" {
		group.Use(pm.RequireAuth(), mw.RateLimitAdmin())
	} else {
		// 使用配置中的 admin_auth.mode
		switch strings.ToLower(cfg.AdminAuth.Mode) {
		case "jwt":
			group.Use(pm.RequireAuth(), mw.RateLimitAdmin())
		default:
			// 开发模式：保留旧的 AdminAuth（Bearer Token			group.Use(mw.AdminAuth(), mw.RateLimitAdmin())
		}
	}
	{
		// 游戏管理 - 使用细粒度权		// @Summary      列出游戏
		// @Tags         Admin/Games
		// @Security     BearerAuth
		// @Param        page       query  int  false  "页码"
		// @Param        pageSize   query     int       false  "每页数量"
		// @Produce      json
		// @Success      200  {object}  model.APIResponse[[]model.Game]
		// @Router       /admin/games [get]
		group.GET("/games", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/games"), gameHandler.ListGames)
		// @Summary      创建游戏
		// @Tags         Admin/Games
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  GamePayload  true  "游戏信息"
		// @Success      201  {object}  model.APIResponse[model.Game]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/games [post]
		group.POST("/games", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/games"), gameHandler.CreateGame)
		// @Summary      获取游戏
		// @Tags         Admin/Games
		// @Security     BearerAuth
		// @Param        id   path  int  true  "游戏ID"
		// @Produce      json
		// @Success      200  {object}  model.APIResponse[model.Game]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/games/{id} [get]
		group.GET("/games/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/games/:id"), gameHandler.GetGame)
		// @Summary      更新游戏
		// @Tags         Admin/Games
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int         true  "游戏ID"
		// @Param        request  body  GamePayload true  "游戏信息"
		// @Success      200  {object}  model.APIResponse[model.Game]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/games/{id} [put]
		group.PUT("/games/:id", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/games/:id"), gameHandler.UpdateGame)
		// @Summary      删除游戏
		// @Tags         Admin/Games
		// @Security     BearerAuth
		// @Param        id   path  int  true  "游戏ID"
		// @Produce      json
		// @Success      200  {object}  model.SuccessResponse
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/games/{id} [delete]
		group.DELETE("/games/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/games/:id"), gameHandler.DeleteGame)
		// @Summary      批量删除游戏
		// @Tags         Admin/Games
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  BatchDeleteGamesRequest  true  "游戏ID列表"
		// @Success      200  {object}  model.APIResponse[map[string]int64]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/games/batch/delete [post]
		group.POST("/games/batch/delete", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/games/batch/delete"), gameHandler.BatchDeleteGames)
		// @Summary      批量更新游戏状态
		// @Tags         Admin/Games
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  BatchUpdateGameStatusRequest  true  "游戏ID列表和状态"
		// @Success      200  {object}  model.APIResponse[adminservice.BatchOperationResponse]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/games/batch/status [post]
		group.POST("/games/batch/status", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/games/batch/status"), gameHandler.BatchUpdateGamesStatus)
		// @Summary      批量更新游戏分类
		// @Tags         Admin/Games
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  BatchUpdateGameCategoryRequest  true  "游戏ID列表和分类"
		// @Success      200  {object}  model.APIResponse[adminservice.BatchOperationResponse]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/games/batch/category [post]
		group.POST("/games/batch/category", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/games/batch/category"), gameHandler.BatchUpdateGamesCategory)
		// @Summary      获取游戏操作日志
		// @Tags         Admin/Games
		// @Security     BearerAuth
		// @Produce      json
		// @Param        id           path   int  true  "游戏ID"
		// @Param        page         query  int  false "页码"
		// @Param        pageSize   query     int       false  "每页数量"
		// @Param        action       query  string false "动作过滤" Enums(create,update,delete)
		// @Param        actor_user_id query int   false "操作者用户ID"
		// @Param        dateFrom   query     string    false  "Start date (YYYY-MM-DD)"
		// @Param        dateTo     query     string    false  "End date (YYYY-MM-DD)"
		// @Param        export       query  string false "导出格式" Enums(csv)
		// @Param        fields         query    string       false  "Export fields (comma separated)"
		// @Param        header_lang  query  string false "列头语言" Enums(en,zh)
		// @Success      200  {object}  model.APIResponse[[]model.OperationLog]
		// @Router       /admin/games/{id}/logs [get]
		group.GET("/games/:id/logs", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/games/:id/logs"), gameHandler.ListGameLogs)

		// 用户管理 - 使用细粒度权		// @Summary      列出用户
		// @Tags         Admin/Users
		// @Security     BearerAuth
		// @Param        page       query     int       false  "页码"
		// @Param        pageSize   query     int       false  "每页数量"
		// @Param        role           query    []string     false  "Role filter"
		// @Param        status         query    []string     false  "Status filter"
		// @Param        dateFrom       query    string       false  "Start date (YYYY-MM-DD)"
		// @Param        dateTo     query     string    false  "End date (YYYY-MM-DD)"
		// @Param        keyword        query    string       false  "Parameter: keyword"
		// @Produce      json
		// @Success      200  {object}  model.APIResponse[[]model.User]
		// @Router       /admin/users [get]
		group.GET("/users", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/users"), userHandler.ListUsers)
		// @Summary      获取用户统计数据
		// @Tags         Admin/Users
		// @Security     BearerAuth
		// @Produce      json
		// @Success      200  {object}  model.APIResponse[adminservice.UserStatsResponse]
		// @Router       /admin/users/stats [get]
		group.GET("/users/stats", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/users/stats"), userHandler.GetUserStats)
		// @Summary      创建用户
		// @Tags         Admin/Users
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  CreateUserPayload  true  "用户信息"
		// @Success      201  {object}  model.APIResponse[model.User]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/users [post]
		group.POST("/users", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/users"), userHandler.CreateUser)
		// @Summary      创建用户和陪玩师
		// @Tags         Admin/Users
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  CreateUserWithPlayerPayload  true  "用户信息和陪玩师信息"
		// @Success      201  {object}  model.SuccessResponse
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/users/with-player [post]
		group.POST("/users/with-player", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/users/with-player"), userHandler.CreateUserWithPlayer)
		// @Summary      批量删除用户
		// @Tags         Admin/Users
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  map[string][]int  true  "{ids: [1,2,3]}"
		// @Success      200  {object}  model.SuccessResponse
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/users/batch-delete [post]
		group.POST("/users/batch-delete", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/users/batch-delete"), userHandler.BatchDeleteUsers)
		// @Summary      获取用户
		// @Tags         Admin/Users
		// @Security     BearerAuth
		// @Param        id   path      int  true  "用户ID"
		// @Produce      json
		// @Success      200  {object}  model.APIResponse[model.User]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/users/{id} [get]
		group.GET("/users/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/users/:id"), userHandler.GetUser)
		// @Summary      更新用户
		// @Tags         Admin/Users
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int                  true  "用户ID"
		// @Param        request  body  UpdateUserPayload    true  "用户信息"
		// @Success      200  {object}  model.APIResponse[model.User]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/users/{id} [put]
		group.PUT("/users/:id", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/users/:id"), userHandler.UpdateUser)
		// @Summary      删除用户
		// @Tags         Admin/Users
		// @Security     BearerAuth
		// @Param        id   path  int  true  "用户ID"
		// @Produce      json
		// @Success      200  {object}  model.SuccessResponse
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/users/{id} [delete]
		group.DELETE("/users/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/users/:id"), userHandler.DeleteUser)
		// @Summary      更新用户状		// @Tags         Admin/Users
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int  true  "用户ID"
		// @Param        request  body  map[string]string  true  "{status}"
		// @Success      200  {object}  model.APIResponse[model.User]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/users/{id}/status [put]
		group.PUT("/users/:id/status", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/users/:id/status"), userHandler.UpdateUserStatus)
		// @Summary      更新用户角色
		// @Tags         Admin/Users
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int  true  "用户ID"
		// @Param        request  body  map[string]string  true  "{role}"
		// @Success      200  {object}  model.APIResponse[model.User]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/users/{id}/role [put]
		group.PUT("/users/:id/role", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/users/:id/role"), userHandler.UpdateUserRole)
		// @Summary      获取用户的订		// @Tags         Admin/Users
		// @Security     BearerAuth
		// @Produce      json
		// @Param        id         path   int      true   "用户ID"
		// @Param        page       query  int      false  "页码"
		// @Param        pageSize   query     int       false  "每页数量"
		// @Param        status         query    []string     false  "Status filter"
		// @Param        dateFrom       query    string       false  "Start date (YYYY-MM-DD)"
		// @Param        dateTo     query     string    false  "End date (YYYY-MM-DD)"
		// @Success      200  {object}  model.APIResponse[[]model.Order]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/users/{id}/orders [get]
		group.GET("/users/:id/orders", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/users/:id/orders"), userHandler.ListUserOrders)
		// @Summary      获取用户操作日志
		// @Tags         Admin/Users
		// @Security     BearerAuth
		// @Produce      json
		// @Param        id           path   int  true  "用户ID"
		// @Param        page         query  int  false "页码"
		// @Param        pageSize   query     int       false  "每页数量"
		// @Param        action       query  string false "动作过滤" Enums(create,update,delete)
		// @Param        actor_user_id query int   false "操作者用户ID"
		// @Param        dateFrom       query    string       false  "Start date (YYYY-MM-DD)"
		// @Param        dateTo     query     string    false  "End date (YYYY-MM-DD)"
		// @Param        export       query  string false "导出格式" Enums(csv)
		// @Param        fields         query    string       false  "Export fields (comma separated)"
		// @Param        header_lang  query  string false "列头语言" Enums(en,zh)
		// @Success      200  {object}  model.APIResponse[[]model.OperationLog]
		// @Router       /admin/users/{id}/logs [get]
		group.GET("/users/:id/logs", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/users/:id/logs"), userHandler.ListUserLogs)
		// @Summary      获取用户登录历史
		// @Tags         Admin/Users
		// @Security     BearerAuth
		// @Produce      json
		// @Param        id          path   int      true   "用户ID"
		// @Param        page        query  int      false  "页码"
		// @Param        page_size   query  int      false  "每页数量"
		// @Success      200  {object}  model.APIResponse[[]model.UserLoginHistory]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/users/{id}/login-history [get]
		group.GET("/users/:id/login-history", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/users/:id/login-history"), userHandler.ListUserLoginHistory)

		// 用户行为分析 - 使用细粒度权限
		// @Summary      获取用户行为统计
		// @Description  获取DAU、平均在线时长、人均消费等用户行为统计数据
		// @Tags         Admin/UserBehavior
		// @Security     BearerAuth
		// @Produce      json
		// @Success      200  {object}  model.APIResponse[adminservice.UserBehaviorStatsResponse]
		// @Router       /admin/users/behavior/stats [get]
		group.GET("/users/behavior/stats", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/users/behavior/stats"), userBehaviorHandler.GetBehaviorStats)
		// @Summary      获取用户活动趋势
		// @Description  获取最近N天的用户活动趋势数据
		// @Tags         Admin/UserBehavior
		// @Security     BearerAuth
		// @Produce      json
		// @Param        days  query  int  false  "统计天数（默认7天）"
		// @Success      200  {object}  model.SuccessResponse
		// @Router       /admin/users/behavior/trend [get]
		group.GET("/users/behavior/trend", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/users/behavior/trend"), userBehaviorHandler.GetActivityTrend)
		// @Summary      获取用户分布
		// @Description  获取用户地域分布、年龄分布等统计数据
		// @Tags         Admin/UserBehavior
		// @Security     BearerAuth
		// @Produce      json
		// @Success      200  {object}  model.APIResponse[adminservice.UserDistributionResponse]
		// @Router       /admin/users/behavior/distribution [get]
		group.GET("/users/behavior/distribution", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/users/behavior/distribution"), userBehaviorHandler.GetUserDistribution)

		// 陪玩师管- 使用细粒度权		// @Summary      列出玩家资料
		// @Tags         Admin/Players
		// @Security     BearerAuth
		// @Param        page       query  int  false  "页码"
		// @Param        pageSize   query     int       false  "每页数量"
		// @Produce      json
		// @Success      200  {object}  model.APIResponse[[]model.Player]
		// @Router       /admin/players [get]
		group.GET("/players", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/players"), playerHandler.ListPlayers)
		// @Summary      新建玩家资料
		// @Tags         Admin/Players
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  CreatePlayerPayload  true  "玩家信息"
		// @Success      201  {object}  model.APIResponse[model.Player]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/players [post]
		group.POST("/players", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/players"), playerHandler.CreatePlayer)
		// @Summary      获取玩家资料
		// @Tags         Admin/Players
		// @Security     BearerAuth
		// @Param        id   path  int  true  "玩家ID"
		// @Produce      json
		// @Success      200  {object}  model.APIResponse[model.Player]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/players/{id} [get]
		group.GET("/players/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/players/:id"), playerHandler.GetPlayer)
		// @Summary      更新玩家资料
		// @Tags         Admin/Players
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int                   true  "玩家ID"
		// @Param        request  body  UpdatePlayerPayload   true  "玩家信息"
		// @Success      200  {object}  model.APIResponse[model.Player]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/players/{id} [put]
		group.PUT("/players/:id", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/players/:id"), playerHandler.UpdatePlayer)
		// @Summary      删除玩家资料
		// @Tags         Admin/Players
		// @Security     BearerAuth
		// @Param        id   path  int  true  "玩家ID"
		// @Produce      json
		// @Success      200  {object}  model.SuccessResponse
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/players/{id} [delete]
		group.DELETE("/players/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/players/:id"), playerHandler.DeletePlayer)
		// @Summary      更新玩家认证状		// @Tags         Admin/Players
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int  true  "玩家ID"
		// @Param        request  body  map[string]string  true  "{verification_status}"
		// @Success      200  {object}  model.APIResponse[model.Player]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/players/{id}/verification [put]
		group.PUT("/players/:id/verification", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/players/:id/verification"), playerHandler.UpdatePlayerVerification)
		// @Summary      更新玩家主游		// @Tags         Admin/Players
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int  true  "玩家ID"
		// @Param        request  body  map[string]uint64  true  "{main_game_id}"
		// @Success      200  {object}  model.APIResponse[model.Player]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/players/{id}/games [put]
		group.PUT("/players/:id/games", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/players/:id/games"), playerHandler.UpdatePlayerGames)
		// @Summary      更新玩家技能标		// @Tags         Admin/Players
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int            true  "玩家ID"
		// @Param        request  body  SkillTagsBody  true  "标签集合"
		// @Success      200  {object}  model.SuccessResponse
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/players/{id}/skill-tags [put]
		group.PUT("/players/:id/skill-tags", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/players/:id/skill-tags"), playerHandler.UpdatePlayerSkillTags)
		// @Summary      获取玩家操作日志
		// @Tags         Admin/Players
		// @Security     BearerAuth
		// @Produce      json
		// @Param        id           path   int  true  "玩家ID"
		// @Param        page         query  int  false "页码"
		// @Param        pageSize   query     int       false  "每页数量"
		// @Param        action       query  string false "动作过滤" Enums(create,update,delete)
		// @Param        actor_user_id query int   false "操作者用户ID"
		// @Param        dateFrom       query    string       false  "Start date (YYYY-MM-DD)"
		// @Param        dateTo     query     string    false  "End date (YYYY-MM-DD)"
		// @Param        export       query  string false "导出格式" Enums(csv)
		// @Param        fields         query    string       false  "Export fields (comma separated)"
		// @Param        header_lang  query  string false "列头语言" Enums(en,zh)
		// @Success      200  {object}  model.APIResponse[[]model.OperationLog]
		// @Router       /admin/players/{id}/logs [get]
		group.GET("/players/:id/logs", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/players/:id/logs"), playerHandler.ListPlayerLogs)
		// @Summary      批量更新陪玩师状态
		// @Tags         Admin/Players
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  BatchUpdateStatusPayload  true  "批量更新请求"
		// @Success      200  {object}  model.APIResponse[map[string]int64]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/players/batch/status [put]
		group.PUT("/players/batch/status", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/players/batch/status"), playerHandler.BatchUpdatePlayerStatus)
		// @Summary      批量删除陪玩师
		// @Tags         Admin/Players
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  BatchDeletePlayersPayload  true  "批量删除请求"
		// @Success      200  {object}  model.APIResponse[map[string]int64]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/players/batch/delete [post]
		group.POST("/players/batch/delete", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/players/batch/delete"), playerHandler.BatchDeletePlayers)
		// @Summary      批量更新陪玩师认证状态
		// @Description  批量审核陪玩师认证，支持批量通过/拒绝
		// @Tags         Admin/Players
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  BatchUpdateVerificationStatusRequest  true  "批量更新认证状态请求"
		// @Success      200  {object}  model.APIResponse[BatchOperationResponse]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/players/batch/verification [post]
		// TODO: Implement BatchUpdateVerificationStatus method in PlayerHandler
		// group.POST("/players/batch/verification", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/players/batch/verification"), playerHandler.BatchUpdateVerificationStatus)
		// @Summary      批量撤销陪玩师认证
		// @Description  批量撤销已通过认证的陪玩师，将状态改为待审核或已拒绝
		// @Tags         Admin/Players
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  BatchRevokeCertificationRequest  true  "批量撤销认证请求"
		// @Success      200  {object}  model.APIResponse[BatchOperationResponse]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/players/batch/revoke-certification [post]
		// TODO: Implement BatchRevokeCertification method in PlayerHandler
		// group.POST("/players/batch/revoke-certification", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/players/batch/revoke-certification"), playerHandler.BatchRevokeCertification)

		// 订单管理 - 使用细粒度权		// @Summary      列出订单
		// @Tags         Admin/Orders
		// @Security     BearerAuth
		// @Param        page        query  int     false  "页码"
		// @Param        pageSize   query     int       false  "每页数量"
		// @Param        status         query    []string     false  "Status filter"
		// @Param        userId     query     int       false  "用户ID"
		// @Param        player_id   query  int     false  "玩家ID"
		// @Param        gameId     query     int       false  "游戏ID"
		// @Param        dateFrom       query    string       false  "Start date (YYYY-MM-DD)"
		// @Param        dateTo     query     string    false  "End date (YYYY-MM-DD)"
		// @Produce      json
		// @Success      200  {object}  model.APIResponse[[]model.Order]
		// @Router       /admin/orders [get]
		group.GET("/orders", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/orders"), orderHandler.ListOrders)
		// @Summary      创建订单
		// @Tags         Admin/Orders
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  CreateOrderPayload  true  "订单信息"
		// @Success      201  {object}  model.APIResponse[model.Order]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/orders [post]
		group.POST("/orders", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/orders"), orderHandler.CreateOrder)
		// @Summary      获取订单
		// @Tags         Admin/Orders
		// @Security     BearerAuth
		// @Param        id   path  int  true  "订单ID"
		// @Produce      json
		// @Success      200  {object}  model.APIResponse[model.Order]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/orders/{id} [get]
		group.GET("/orders/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/orders/:id"), orderHandler.GetOrder)
		// @Summary      更新订单
		// @Tags         Admin/Orders
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int                true  "订单ID"
		// @Param        request  body  UpdateOrderPayload true  "订单信息"
		// @Success      200  {object}  model.APIResponse[model.Order]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/orders/{id} [put]
		group.PUT("/orders/:id", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/orders/:id"), orderHandler.UpdateOrder)
		// @Summary      删除订单
		// @Tags         Admin/Orders
		// @Security     BearerAuth
		// @Param        id   path  int  true  "订单ID"
		// @Produce      json
		// @Success      200  {object}  model.SuccessResponse
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/orders/{id} [delete]
		group.DELETE("/orders/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/orders/:id"), orderHandler.DeleteOrder)
		// @Summary      评价订单
		// @Tags         Admin/Orders
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int                true  "订单ID"
		// @Param        request  body  ReviewOrderPayload true  "评价信息"
		// @Success      200  {object}  model.APIResponse[model.Review]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/orders/{id}/review [post]
		group.POST("/orders/:id/review", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/orders/:id/review"), orderHandler.ReviewOrder)
		// @Summary      取消订单
		// @Tags         Admin/Orders
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int                true  "订单ID"
		// @Param        request  body  CancelOrderPayload true  "取消信息"
		// @Success      200  {object}  model.APIResponse[model.Order]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/orders/{id}/cancel [post]
		group.POST("/orders/:id/cancel", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/orders/:id/cancel"), orderHandler.CancelOrder)
		// @Summary      指派订单的陪玩师
		// @Tags         Admin/Orders
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int                 true  "订单ID"
		// @Param        request  body  AssignOrderPayload  true  "指派信息"
		// @Success      200  {object}  model.APIResponse[model.Order]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/orders/{id}/assign [post]
		group.POST("/orders/:id/assign", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/orders/:id/assign"), orderHandler.AssignOrder)
		// @Summary      确认订单
		// @Description  将订单状态从 pending 置为 confirmed
		// @Tags         Admin/Orders
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int               true  "订单ID"
		// @Param        request  body  orderNotePayload  false "备注（可选）"
		// @Success      200  {object}  model.APIResponse[model.Order]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/orders/{id}/confirm [post]
		group.POST("/orders/:id/confirm", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/orders/:id/confirm"), orderHandler.ConfirmOrder)
		// @Summary      开始订		// @Description  将订单状态从 confirmed 置为 in_progress
		// @Tags         Admin/Orders
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int               true  "订单ID"
		// @Param        request  body  orderNotePayload  false "备注（可选）"
		// @Success      200  {object}  model.APIResponse[model.Order]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/orders/{id}/start [post]
		group.POST("/orders/:id/start", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/orders/:id/start"), orderHandler.StartOrder)
		// @Summary      完成订单
		// @Description  将订单状态从 in_progress 置为 completed，并记录完成时间
		// @Tags         Admin/Orders
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int               true  "订单ID"
		// @Param        request  body  orderNotePayload  false "备注（可选）"
		// @Success      200  {object}  model.APIResponse[model.Order]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/orders/{id}/complete [post]
		group.POST("/orders/:id/complete", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/orders/:id/complete"), orderHandler.CompleteOrder)
		// @Summary      退款订		// @Tags         Admin/Orders
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int                  true  "订单ID"
		// @Param        request        body     orderRefundPayload true   "Request body"
		// @Success      200  {object}  model.APIResponse[model.Order]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/orders/{id}/refund [post]
		group.POST("/orders/:id/refund", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/orders/:id/refund"), orderHandler.RefundOrder)
		// @Summary      获取订单操作日志
		// @Tags         Admin/Orders
		// @Security     BearerAuth
		// @Produce      json
		// @Param        id         path   int  true  "订单ID"
		// @Param        page       query  int  false "页码"
		// @Param        pageSize   query     int       false  "每页数量"
		// @Param        action     query  string false "动作过滤" Enums(create,assign_player,update_status,cancel,delete)
		// @Param        actor_user_id query int false "操作者用户ID"
		// @Param        dateFrom       query    string       false  "Start date (YYYY-MM-DD)"
		// @Param        dateTo     query     string    false  "End date (YYYY-MM-DD)"
		// @Param        export     query  string false "导出格式" Enums(csv)
		// @Param        fields     query  string false "Export columns (comma separated)"
		// @Param        header_lang query string false "列头语言" Enums(en,zh)
		// @Success      200  {object}  model.APIResponse[[]model.OperationLog]
		// @Router       /admin/orders/{id}/logs [get]
		group.GET("/orders/:id/logs", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/orders/:id/logs"), orderHandler.ListOrderLogs)
		// @Summary      获取订单时间		// @Tags         Admin/Orders
		// @Security     BearerAuth
		// @Produce      json
		// @Param        id   path  int  true  "订单ID"
		// @Success      200  {object}  model.APIResponse[[]service.OrderTimelineItem]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/orders/{id}/timeline [get]
		group.GET("/orders/:id/timeline", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/orders/:id/timeline"), orderHandler.GetOrderTimeline)
		// @Summary      获取订单支付记录
		// @Tags         Admin/Orders
		// @Security     BearerAuth
		// @Produce      json
		// @Param        id   path  int  true  "订单ID"
		// @Success      200  {object}  model.APIResponse[[]model.Payment]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/orders/{id}/payments [get]
		group.GET("/orders/:id/payments", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/orders/:id/payments"), orderHandler.ListOrderPayments)
		// @Summary      获取订单退款记		// @Tags         Admin/Orders
		// @Security     BearerAuth
		// @Produce      json
		// @Param        id   path  int  true  "订单ID"
		// @Success      200  {object}  model.APIResponse[[]service.OrderRefundItem]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/orders/{id}/refunds [get]
		group.GET("/orders/:id/refunds", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/orders/:id/refunds"), orderHandler.ListOrderRefunds)
		// @Summary      获取订单评价列表
		// @Tags         Admin/Orders
		// @Security     BearerAuth
		// @Produce      json
		// @Param        id   path  int  true  "订单ID"
		// @Success      200  {object}  model.APIResponse[[]model.Review]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/orders/{id}/reviews [get]
		group.GET("/orders/:id/reviews", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/orders/:id/reviews"), orderHandler.ListOrderReviews)

		// 批量订单操作 - Batch Order Operations
		// @Summary      批量取消订单
		// @Description  取消多个待处理或已确认的订单
		// @Tags         Admin/Orders
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  BatchCancelOrdersRequest  true  "订单ID列表和取消原因"
		// @Success      200  {object}  model.APIResponse[BatchOperationResponse]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/orders/batch/cancel [post]
		group.POST("/orders/batch/cancel", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/orders/batch/cancel"), orderHandler.BatchCancelOrders)
		// @Summary      批量确认订单
		// @Description  确认多个待处理的订单
		// @Tags         Admin/Orders
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  BatchConfirmOrdersRequest  true  "订单ID列表"
		// @Success      200  {object}  model.APIResponse[BatchOperationResponse]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/orders/batch/confirm [post]
		group.POST("/orders/batch/confirm", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/orders/batch/confirm"), orderHandler.BatchConfirmOrders)
		// @Summary      批量完成订单
		// @Description  将多个进行中的订单标记为已完成
		// @Tags         Admin/Orders
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  BatchCompleteOrdersRequest  true  "订单ID列表"
		// @Success      200  {object}  model.APIResponse[BatchOperationResponse]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/orders/batch/complete [post]
		group.POST("/orders/batch/complete", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/orders/batch/complete"), orderHandler.BatchCompleteOrders)
		// @Summary      批量退款订单
		// @Description  对多个订单执行退款操作
		// @Tags         Admin/Orders
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  BatchRefundOrdersRequest  true  "订单ID列表和退款信息"
		// @Success      200  {object}  model.APIResponse[BatchOperationResponse]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/orders/batch/refund [post]
		group.POST("/orders/batch/refund", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/orders/batch/refund"), orderHandler.BatchRefundOrders)
		// @Summary      批量删除订单
		// @Description  软删除多个订单
		// @Tags         Admin/Orders
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  BatchDeleteOrdersRequest  true  "订单ID列表"
		// @Success      200  {object}  model.APIResponse[BatchOperationResponse]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/orders/batch/delete [post]
		group.POST("/orders/batch/delete", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/orders/batch/delete"), orderHandler.BatchDeleteOrders)
		// @Summary      批量更新订单状态
		// @Description  更新多个订单的状态
		// @Tags         Admin/Orders
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  BatchUpdateOrderStatusRequest  true  "订单ID列表和新状态"
		// @Success      200  {object}  model.APIResponse[BatchOperationResponse]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/orders/batch/status [put]
		group.PUT("/orders/batch/status", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/orders/batch/status"), orderHandler.BatchUpdateOrderStatus)
		// @Summary      批量指派订单
		// @Description  将多个订单指派给指定的陪玩师
		// @Tags         Admin/Orders
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  BatchAssignOrdersRequest  true  "订单ID列表和陪玩师ID"
		// @Success      200  {object}  model.APIResponse[BatchOperationResponse]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/orders/batch/assign [post]
		group.POST("/orders/batch/assign", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/orders/batch/assign"), orderHandler.BatchAssignOrders)

		// 支付管理 - 使用细粒度权		// @Summary      列出支付
		// @Tags         Admin/Payments
		// @Security     BearerAuth
		// @Param        page        query  int       false  "页码"
		// @Param        pageSize   query     int       false  "每页数量"
		// @Param        status         query    []string     false  "Status filter"
		// @Param        method      query  []string  false  "支付方式"
		// @Param        userId     query     int       false  "用户ID"
		// @Param        orderId     query     int       false  "订单ID"
		// @Param        dateFrom       query    string       false  "Start date (YYYY-MM-DD)"
		// @Param        dateTo     query     string    false  "End date (YYYY-MM-DD)"
		// @Produce      json
		// @Success      200  {object}  model.APIResponse[[]model.Payment]
		// @Router       /admin/payments [get]
		group.GET("/payments", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/payments"), paymentHandler.ListPayments)
		// @Summary      创建支付记录
		// @Tags         Admin/Payments
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  CreatePaymentPayload  true  "支付信息"
		// @Success      201  {object}  model.APIResponse[model.Payment]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/payments [post]
		group.POST("/payments", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/payments"), paymentHandler.CreatePayment)
		// @Summary      获取支付
		// @Tags         Admin/Payments
		// @Security     BearerAuth
		// @Param        id   path  int  true  "支付ID"
		// @Produce      json
		// @Success      200  {object}  model.APIResponse[model.Payment]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/payments/{id} [get]
		group.GET("/payments/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/payments/:id"), paymentHandler.GetPayment)
		// @Summary      更新支付
		// @Tags         Admin/Payments
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int                  true  "支付ID"
		// @Param        request  body  UpdatePaymentPayload true  "支付信息"
		// @Success      200  {object}  model.APIResponse[model.Payment]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/payments/{id} [put]
		group.PUT("/payments/:id", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/payments/:id"), paymentHandler.UpdatePayment)
		// @Summary      删除支付
		// @Tags         Admin/Payments
		// @Security     BearerAuth
		// @Param        id   path  int  true  "支付ID"
		// @Produce      json
		// @Success      200  {object}  model.SuccessResponse
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/payments/{id} [delete]
		group.DELETE("/payments/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/payments/:id"), paymentHandler.DeletePayment)
		// @Summary      退款支		// @Tags         Admin/Payments
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int                  true  "支付ID"
		// @Param        request  body  RefundPaymentPayload true  "退款信息"
		// @Success      200  {object}  model.APIResponse[model.Payment]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/payments/{id}/refund [post]
		group.POST("/payments/:id/refund", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/payments/:id/refund"), paymentHandler.RefundPayment)
		// @Summary      确认支付入账
		// @Tags         Admin/Payments
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int                     true  "支付ID"
		// @Param        request  body  CapturePaymentPayload   true  "入账信息"
		// @Success      200  {object}  model.APIResponse[model.Payment]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/payments/{id}/capture [post]
		group.POST("/payments/:id/capture", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/payments/:id/capture"), paymentHandler.CapturePayment)
		// @Summary      获取支付操作日志
		// @Tags         Admin/Payments
		// @Security     BearerAuth
		// @Produce      json
		// @Param        id         path   int  true  "支付ID"
		// @Param        page       query  int  false "页码"
		// @Param        pageSize   query     int       false  "每页数量"
		// @Param        action     query  string false "动作过滤" Enums(create,capture,update_status,refund,delete)
		// @Param        actor_user_id query int false "操作者用户ID"
		// @Param        dateFrom       query    string       false  "Start date (YYYY-MM-DD)"
		// @Param        dateTo     query     string    false  "End date (YYYY-MM-DD)"
		// @Param        export     query  string false "导出格式" Enums(csv)
		// @Param        fields     query  string false "Export columns (comma separated)"
		// @Param        header_lang query string false "列头语言" Enums(en,zh)
		// @Success      200  {object}  model.APIResponse[[]model.OperationLog]
		// @Router       /admin/payments/{id}/logs [get]
		group.GET("/payments/:id/logs", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/payments/:id/logs"), paymentHandler.ListPaymentLogs)
		// @Summary      获取退款历史
		// @Description  获取支付记录的所有退款操作历史
		// @Tags         Admin/Payments
		// @Security     BearerAuth
		// @Produce      json
		// @Param        id   path  int  true  "支付ID"
		// @Success      200  {object}  model.APIResponse[[]model.OperationLog]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/payments/{id}/refunds [get]
		group.GET("/payments/:id/refunds", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/payments/:id/refunds"), paymentHandler.GetRefundHistory)

		// 批量支付操作路由 - TODO: Implement batch payment handlers
		// batchGroup := group.Group("/payments/batch")
		// {
		//     // Route definitions commented out pending handler implementation
		// }

		// 评价管理 - 使用细粒度权		// @Summary      评价列表
		// @Tags         Admin/Reviews
		// @Security     BearerAuth
		// @Produce      json
		// @Param        page       query  int  false  "页码"
		// @Param        pageSize   query     int       false  "每页数量"
		// @Param        orderId     query     int       false  "订单ID"
		// @Param        userId     query     int       false  "用户ID"
		// @Param        playerId   query     int       false  "陪玩师ID"
		// @Param        dateFrom       query    string       false  "Start date (YYYY-MM-DD)"
		// @Param        dateTo     query     string    false  "End date (YYYY-MM-DD)"
		// @Success      200  {object}  model.APIResponse[[]model.Review]
		// @Router       /admin/reviews [get]
		group.GET("/reviews", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/reviews"), reviewHandler.ListReviews)
		// @Summary      获取待审核评价列表
		// @Tags         Admin/Reviews
		// @Security     BearerAuth
		// @Produce      json
		// @Param        page       query  int  false  "页码"
		// @Param        pageSize   query  int  false  "每页数量"
		// @Success      200  {object}  model.APIResponse[[]model.Review]
		// @Router       /admin/reviews/pending [get]
		group.GET("/reviews/pending", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/reviews/pending"), reviewHandler.ListPendingReviews)
		// @Summary      创建评价
		// @Tags         Admin/Reviews
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  CreateReviewPayload  true  "评价"
		// @Success      201  {object}  model.APIResponse[model.Review]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/reviews [post]
		group.POST("/reviews", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/reviews"), reviewHandler.CreateReview)
		// @Summary      获取评价
		// @Tags         Admin/Reviews
		// @Security     BearerAuth
		// @Produce      json
		// @Param        id   path  int  true  "评价ID"
		// @Success      200  {object}  model.APIResponse[model.Review]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/reviews/{id} [get]
		group.GET("/reviews/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/reviews/:id"), reviewHandler.GetReview)
		// @Summary      更新评价
		// @Tags         Admin/Reviews
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int                    true  "评价ID"
		// @Param        request  body  UpdateReviewPayload    true  "评价"
		// @Success      200  {object}  model.APIResponse[model.Review]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/reviews/{id} [put]
		group.PUT("/reviews/:id", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/reviews/:id"), reviewHandler.UpdateReview)
		// @Summary      删除评价
		// @Tags         Admin/Reviews
		// @Security     BearerAuth
		// @Produce      json
		// @Param        id   path  int  true  "评价ID"
		// @Success      200  {object}  model.SuccessResponse
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/reviews/{id} [delete]
		group.DELETE("/reviews/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/reviews/:id"), reviewHandler.DeleteReview)
		// @Summary      批准评价
		// @Tags         Admin/Reviews
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id   path  int  true  "评价ID"
		// @Success      200  {object}  model.APIResponse[any]
		// @Failure      400  {object}  model.ErrorResponse
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/reviews/{id}/approve [put]
		group.PUT("/reviews/:id/approve", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/reviews/:id/approve"), reviewHandler.ApproveReview)
		// @Summary      拒绝评价
		// @Tags         Admin/Reviews
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int                   true  "评价ID"
		// @Param        request  body  RejectReviewPayload   true  "拒绝信息"
		// @Success      200  {object}  model.APIResponse[any]
		// @Failure      400  {object}  model.ErrorResponse
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/reviews/{id}/reject [put]
		group.PUT("/reviews/:id/reject", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/reviews/:id/reject"), reviewHandler.RejectReview)
		// @Summary      批量批准评价
		// @Tags         Admin/Reviews
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  BatchApprovePayload  true  "批量批准信息"
		// @Success      200  {object}  model.APIResponse[any]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/reviews/batch-approve [put]
		group.PUT("/reviews/batch-approve", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/reviews/batch-approve"), reviewHandler.BatchApproveReviews)
		// @Summary      批量拒绝评价
		// @Tags         Admin/Reviews
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  BatchRejectPayload  true  "批量拒绝信息"
		// @Success      200  {object}  model.APIResponse[any]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/reviews/batch-reject [put]
		group.PUT("/reviews/batch-reject", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/reviews/batch-reject"), reviewHandler.BatchRejectReviews)
		// @Summary      批准所有不含敏感词的待审核评价
		// @Tags         Admin/Reviews
		// @Security     BearerAuth
		// @Produce      json
		// @Success      200  {object}  model.APIResponse[ApproveAllNonSensitiveResponse]
		// @Failure      500  {object}  model.ErrorResponse
		// @Router       /admin/reviews/approve-all-non-sensitive [put]
		group.PUT("/reviews/approve-all-non-sensitive", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/reviews/batch-approve"), reviewHandler.ApproveAllNonSensitiveReviews)
		// @Summary      获取陪玩师的评价
		// @Tags         Admin/Players
		// @Security     BearerAuth
		// @Produce      json
		// @Param        id         path   int  true  "陪玩师ID"
		// @Param        page       query  int  false  "页码"
		// @Param        pageSize   query     int       false  "每页数量"
		// @Success      200  {object}  model.APIResponse[[]model.Review]
		// @Router       /admin/players/{id}/reviews [get]
		group.GET("/players/:id/reviews", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/players/:id/reviews"), reviewHandler.ListPlayerReviews)
		// @Summary      获取评价操作日志
		// @Tags         Admin/Reviews
		// @Security     BearerAuth
		// @Produce      json
		// @Param        id           path   int  true  "评价ID"
		// @Param        page         query  int  false "页码"
		// @Param        pageSize   query     int       false  "每页数量"
		// @Param        action       query  string false "动作过滤" Enums(create,update,delete)
		// @Param        actor_user_id query int   false "操作者用户ID"
		// @Param        dateFrom       query    string       false  "Start date (YYYY-MM-DD)"
		// @Param        dateTo     query     string    false  "End date (YYYY-MM-DD)"
		// @Param        export       query  string false "导出格式" Enums(csv)
		// @Param        fields         query    string       false  "Export fields (comma separated)"
		// @Param        header_lang  query  string false "列头语言" Enums(en,zh)
		// @Success      200  {object}  model.APIResponse[[]model.OperationLog]
		// @Router       /admin/reviews/{id}/logs [get]
		group.GET("/reviews/:id/logs", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/reviews/:id/logs"), reviewHandler.ListReviewLogs)

		// 评价举报管理
		// @Summary      创建评价举报
		// @Tags         Admin/Reviews
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int                      true  "评价ID"
		// @Param        request  body  CreateReviewReportPayload  true  "举报信息"
		// @Success      201  {object}  model.APIResponse[CreateReviewReportResponse]
		// @Failure      400  {object}  model.ErrorResponse
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/reviews/{id}/reports [post]
		group.POST("/reviews/:id/reports", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/reviews/:id/reports"), reviewHandler.CreateReviewReport)

		// @Summary      列出评价举报
		// @Tags         Admin/Reviews
		// @Security     BearerAuth
		// @Produce      json
		// @Param        page         query  int     false  "页码"
		// @Param        pageSize     query  int     false  "每页数量"
		// @Param        review_id    query  int     false  "评价ID"
		// @Param        reporter_id  query  int     false  "举报人ID"
		// @Param        status       query  string  false  "状态" Enums(pending,approved,rejected)
		// @Param        date_from    query  string  false  "开始日期 (YYYY-MM-DD)"
		// @Param        date_to      query  string  false  "结束日期 (YYYY-MM-DD)"
		// @Success      200  {object}  model.APIResponse[[]ReviewReportDTO]
		// @Router       /admin/review-reports [get]
		group.GET("/review-reports", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/review-reports"), reviewHandler.ListReviewReports)

		// @Summary      获取举报详情
		// @Tags         Admin/Reviews
		// @Security     BearerAuth
		// @Produce      json
		// @Param        id   path  int  true  "举报ID"
		// @Success      200  {object}  model.APIResponse[ReviewReportDTO]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/review-reports/{id} [get]
		group.GET("/review-reports/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/review-reports/:id"), reviewHandler.GetReviewReport)

		// @Summary      处理评价举报
		// @Tags         Admin/Reviews
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int                       true  "举报ID"
		// @Param        request  body  HandleReviewReportPayload  true  "处理信息"
		// @Success      200  {object}  model.APIResponse[HandleReviewReportResponse]
		// @Failure      400  {object}  model.ErrorResponse
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/review-reports/{id}/handle [put]
		group.PUT("/review-reports/:id/handle", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/review-reports/:id/handle"), reviewHandler.HandleReviewReport)

		// 评价回复管理
		// @Summary      更新评价回复
		// @Tags         Admin/Reviews
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int                  true  "回复ID"
		// @Param        request  body  UpdateReplyPayload   true  "回复内容"
		// @Success      200  {object}  model.SuccessResponse
		// @Failure      400  {object}  model.ErrorResponse
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/review-replies/{id} [put]
		group.PUT("/review-replies/:id", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/review-replies/:id"), reviewHandler.UpdateReply)

		// @Summary      删除评价回复
		// @Tags         Admin/Reviews
		// @Security     BearerAuth
		// @Produce      json
		// @Param        id   path  int  true  "回复ID"
		// @Success      200  {object}  model.SuccessResponse
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/review-replies/{id} [delete]
		group.DELETE("/review-replies/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/review-replies/:id"), reviewHandler.DeleteReply)

		// 操作日志管理
		// @Summary      搜索操作日志
		// @Tags         Admin/OperationLogs
		// @Security     BearerAuth
		// @Produce      json
		// @Param        page          query  int     false  "页码"
		// @Param        pageSize      query  int     false  "每页数量"
		// @Param        entity_type   query  string  false  "实体类型" Enums(review,order,payment,user,player)
		// @Param        entity_id     query  int     false  "实体ID"
		// @Param        action        query  string  false  "动作过滤"
		// @Param        actor_user_id query  int     false  "操作者用户ID"
		// @Param        date_from     query  string  false  "开始日期 (YYYY-MM-DD)"
		// @Param        date_to       query  string  false  "结束日期 (YYYY-MM-DD)"
		// @Param        export        query  string  false  "导出格式" Enums(csv)
		// @Success      200  {object}  model.APIResponse[[]model.OperationLog]
		// @Router       /admin/operation-logs [get]
		group.GET("/operation-logs", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/operation-logs"), reviewHandler.SearchOperationLogs)

		// @Summary      导出操作日志
		// @Tags         Admin/OperationLogs
		// @Security     BearerAuth
		// @Produce      text/csv
		// @Param        entity_type   query  string  false  "实体类型" Enums(review,order,payment,user,player)
		// @Param        entity_id     query  int     false  "实体ID"
		// @Param        action        query  string  false  "动作过滤"
		// @Param        actor_user_id query  int     false  "操作者用户ID"
		// @Param        date_from     query  string  false  "开始日期 (YYYY-MM-DD)"
		// @Param        date_to       query  string  false  "结束日期 (YYYY-MM-DD)"
		// @Success      200  {file}  file
		// @Router       /admin/operation-logs/export [get]
		group.GET("/operation-logs/export", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/operation-logs/export"), reviewHandler.ExportOperationLogs)
	}
}

// RegisterReviewSettingsRoutes 注册评价展示设置相关路由
func RegisterReviewSettingsRoutes(router gin.IRouter, settingsSvc *reviewservice.SettingsService, pm *mw.PermissionMiddleware) {
	h := NewReviewSettingsHandler(settingsSvc)
	group := router
	// 设置接口均需要认证 + 速率限制
	cfg := config.Load()
	if os.Getenv("APP_ENV") == "production" {
		group.Use(pm.RequireAuth(), mw.RateLimitAdmin())
	} else {
		switch strings.ToLower(cfg.AdminAuth.Mode) {
		case "jwt":
			group.Use(pm.RequireAuth(), mw.RateLimitAdmin())
		default:
			group.Use(mw.AdminAuth(), mw.RateLimitAdmin())
		}
	}

	// 评价展示设置 - 使用细粒度权限
	// @Summary      获取评价展示设置
	// @Description  获取当前的评价展示规则配置
	// @Tags         Admin/ReviewSettings
	// @Security     BearerAuth
	// @Produce      json
	// @Success      200  {object}  model.APIResponse[model.ReviewDisplaySettings]
	// @Failure      500  {object}  model.ErrorResponse
	// @Router       /admin/review-settings [get]
	group.GET("/review-settings", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/review-settings"), h.GetReviewSettings)

	// @Summary      更新评价展示设置
	// @Description  更新评价展示规则配置，支持部分更新
	// @Tags         Admin/ReviewSettings
	// @Security     BearerAuth
	// @Accept       json
	// @Produce      json
	// @Param        request  body  UpdateReviewSettingsPayload  true  "设置信息"
	// @Success      200  {object}  model.APIResponse[model.ReviewDisplaySettings]
	// @Failure      400  {object}  model.ErrorResponse
	// @Failure      500  {object}  model.ErrorResponse
	// @Router       /admin/review-settings [put]
	group.PUT("/review-settings", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/review-settings"), h.UpdateReviewSettings)
}

// RegisterStatsRoutes 注册统计相关路由
// 使用细粒度权限控制（method+path 级别）
func RegisterStatsRoutes(router gin.IRouter, stats *statsservice.StatsService, pm *mw.PermissionMiddleware) {
	h := NewStatsHandler(stats)
	group := router
	// 统计接口均需要认+ 速率限制
	cfg := config.Load()
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
	// 统计接口 - 使用细粒度权	// @Summary      仪表板数	// @Description  获取平台统计数据总览
	// @Tags         Admin - Stats
	// @Accept       json
	// @Produce      json
	// @Success      200            {object}  model.APIResponse[stats.DashboardData]
	// @Router       /admin/stats/dashboard [get]
	group.GET("/stats/dashboard", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/stats/dashboard"), h.Dashboard)
	// @Summary      收入趋势
	// @Description  获取指定天数的收入趋	// @Tags         Admin - Stats
	// @Accept       json
	// @Produce      json
	// @Param        days           query     int     false  "天数" default(7)
	// @Success      200            {object}  model.APIResponse[[]stats.DateValue]
	// @Router       /admin/stats/revenue-trend [get]
	group.GET("/stats/revenue-trend", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/stats/revenue-trend"), h.RevenueTrend)
	// @Summary      用户增长趋势
	// @Description  获取指定天数的用户增长趋	// @Tags         Admin - Stats
	// @Accept       json
	// @Produce      json
	// @Param        days           query     int     false  "天数" default(7)
	// @Success      200            {object}  model.APIResponse[[]stats.DateValue]
	// @Router       /admin/stats/user-growth [get]
	group.GET("/stats/user-growth", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/stats/user-growth"), h.UserGrowth)
	// @Summary      订单状态汇	// @Description  获取各状态订单数量统	// @Tags         Admin - Stats
	// @Accept       json
	// @Produce      json
	// @Success      200            {object}  model.APIResponse[map[string]int64]
	// @Router       /admin/stats/orders [get]
	group.GET("/stats/orders", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/stats/orders"), h.OrdersSummary)
	// @Summary      顶级陪玩	// @Description  获取收入最高的陪玩师列	// @Tags         Admin - Stats
	// @Accept       json
	// @Produce      json
	// @Param        limit          query     int     false  "数量限制" default(10)
	// @Success      200            {object}  model.APIResponse[[]stats.PlayerValue]
	// @Router       /admin/stats/top-players [get]
	group.GET("/stats/top-players", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/stats/top-players"), h.TopPlayers)
	// @Summary      审计概览
	// @Description  获取审计日志统计概览
	// @Tags         Admin - Stats
	// @Accept       json
	// @Produce      json
	// @Param        from           query     string  false  "开始日期"
	// @Param        to             query     string  false  "结束日期"
	// @Success      200            {object}  model.SuccessResponse
	// @Router       /admin/stats/audit/overview [get]
	group.GET("/stats/audit/overview", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/stats/audit/overview"), h.AuditOverview)
	// @Summary      审计趋势
	// @Description  获取审计日志时间趋势
	// @Tags         Admin - Stats
	// @Accept       json
	// @Produce      json
	// @Param        from           query     string  false  "开始日期"
	// @Param        to             query     string  false  "结束日期"
	// @Param        entity         query     string  false  "实体类型"
	// @Param        action         query     string  false  "操作类型"
	// @Success      200            {object}  model.APIResponse[[]stats.DateValue]
	// @Router       /admin/stats/audit/trend [get]
	group.GET("/stats/audit/trend", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/stats/audit/trend"), h.AuditTrend)
}

// RegisterSyncRoutes 注册同步专用路由（不受限流限制）
// 注意：此函数已废弃，同步路由现在通过 RegisterSyncRoutesWithServices 注册
func RegisterSyncRoutes(router gin.IRouter, svc *adminservice.AdminService, pm *mw.PermissionMiddleware) {
	// 此函数保留为空，实际的同步路由在 router.go 的 registerRBACRoutes 中注册
}

// RegisterSyncRoutesWithServices 注册同步专用路由（使用已初始化的服务）
// 只保留批量同步接口，前端通过一次请求完成所有同步
func RegisterSyncRoutesWithServices(router gin.IRouter, roleSvc *roleservice.RoleService, permService *permissionservice.PermissionService, menuSvc *adminservice.MenuService, pm *mw.PermissionMiddleware) {
	batchSyncHandler := NewBatchSyncHandler(menuSvc, permService, roleSvc)

	// 创建同步专用路由组 - 只认证，不限流
	syncGroup := router.Group("/sync")
	syncGroup.Use(pm.RequireAuth())
	{
		// 批量同步接口 - 一次性同步菜单、权限并分配超管权限
		syncGroup.POST("/batch", batchSyncHandler.BatchSync)
	}
}

// RegisterSensitiveWordRoutes 注册敏感词管理路由
// 这是一个独立的函数，需要在主路由注册时调用
func RegisterSensitiveWordRoutes(router gin.IRouter, handler *SensitiveWordHandler, pm *mw.PermissionMiddleware) {
	group := router
	{
		// @Summary      列出敏感词
		// @Tags         Admin/SensitiveWords
		// @Security     BearerAuth
		// @Produce      json
		// @Param        page       query  int     false  "页码"
		// @Param        pageSize   query  int     false  "每页数量"
		// @Param        keyword    query  string  false  "关键词搜索"
		// @Param        category   query  string  false  "分类" Enums(political,pornographic,violent,advertising,other)
		// @Param        severity   query  string  false  "严重程度" Enums(low,medium,high)
		// @Success      200  {object}  model.APIResponse[sensitiveword.ListSensitiveWordsResponse]
		// @Router       /admin/sensitive-words [get]
		group.GET("/sensitive-words", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/sensitive-words"), handler.ListSensitiveWords)

		// @Summary      添加敏感词
		// @Tags         Admin/SensitiveWords
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  sensitiveword.AddSensitiveWordRequest  true  "敏感词信息"
		// @Success      201  {object}  model.APIResponse[sensitiveword.SensitiveWordDTO]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/sensitive-words [post]
		group.POST("/sensitive-words", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/sensitive-words"), handler.AddSensitiveWord)

		// @Summary      更新敏感词
		// @Tags         Admin/SensitiveWords
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        id       path  int                                      true  "敏感词ID"
		// @Param        request  body  sensitiveword.UpdateSensitiveWordRequest  true  "敏感词信息"
		// @Success      200  {object}  model.APIResponse[any]
		// @Failure      400  {object}  model.ErrorResponse
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/sensitive-words/{id} [put]
		group.PUT("/sensitive-words/:id", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/sensitive-words/:id"), handler.UpdateSensitiveWord)

		// @Summary      删除敏感词
		// @Tags         Admin/SensitiveWords
		// @Security     BearerAuth
		// @Produce      json
		// @Param        id   path  int  true  "敏感词ID"
		// @Success      200  {object}  model.APIResponse[any]
		// @Failure      404  {object}  model.ErrorResponse
		// @Router       /admin/sensitive-words/{id} [delete]
		group.DELETE("/sensitive-words/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/sensitive-words/:id"), handler.DeleteSensitiveWord)

		// @Summary      检测敏感词
		// @Tags         Admin/Reviews
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  sensitiveword.DetectSensitiveWordsRequest  true  "检测内容"
		// @Success      200  {object}  model.APIResponse[sensitiveword.DetectSensitiveWordsResponse]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/reviews/detect-sensitive [post]
		group.POST("/reviews/detect-sensitive", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/reviews/detect-sensitive"), handler.DetectSensitiveWords)

		// 批量操作路由
		// @Summary      批量添加敏感词
		// @Description  批量添加敏感词，最多100条
		// @Tags         Admin/SensitiveWords
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  BatchAddSensitiveWordsRequest  true  "批量添加请求"
		// @Success      200  {object}  model.APIResponse[sensitiveword.BatchOperationResult]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/sensitive-words/batch/add [post]
		group.POST("/sensitive-words/batch/add", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/sensitive-words/batch/add"), handler.BatchAddSensitiveWords)

		// @Summary      批量删除敏感词
		// @Description  批量删除敏感词，最多100条
		// @Tags         Admin/SensitiveWords
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  BatchDeleteSensitiveWordsRequest  true  "批量删除请求"
		// @Success      200  {object}  model.APIResponse[sensitiveword.BatchOperationResult]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/sensitive-words/batch/delete [post]
		group.POST("/sensitive-words/batch/delete", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/sensitive-words/batch/delete"), handler.BatchDeleteSensitiveWords)

		// @Summary      批量更新敏感词状态
		// @Description  批量启用/禁用敏感词，最多100条
		// @Tags         Admin/SensitiveWords
		// @Security     BearerAuth
		// @Accept       json
		// @Produce      json
		// @Param        request  body  BatchUpdateSensitiveWordStatusRequest  true  "批量更新状态请求"
		// @Success      200  {object}  model.APIResponse[sensitiveword.BatchOperationResult]
		// @Failure      400  {object}  model.ErrorResponse
		// @Router       /admin/sensitive-words/batch/status [put]
		group.PUT("/sensitive-words/batch/status", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/sensitive-words/batch/status"), handler.BatchUpdateSensitiveWordStatus)
	}
}

// RegisterReviewStatsRoutes 注册评价统计路由
// 使用细粒度权限控制（method+path 级别）
func RegisterReviewStatsRoutes(router gin.IRouter, svc *reviewservice.ReviewStatsService, pm *mw.PermissionMiddleware) {
	h := NewReviewStatsHandler(svc)
	group := router
	{
		// @Summary      获取评价统计概览
		// @Description  获取总评价数、平均评分、各评分段分布
		// @Tags         Admin/Reviews
		// @Security     BearerAuth
		// @Produce      json
		// @Success      200  {object}  model.APIResponse[reviewservice.GetReviewStatsResponse]
		// @Failure      500  {object}  model.ErrorResponse
		// @Router       /admin/reviews/stats [get]
		group.GET("/reviews/stats", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/reviews/stats"), h.GetReviewStats)

		// @Summary      获取评价趋势
		// @Description  获取最近N天的评价数量趋势
		// @Tags         Admin/Reviews
		// @Security     BearerAuth
		// @Produce      json
		// @Param        days  query  int  false  "统计天数（默认30天）"
		// @Success      200  {object}  model.APIResponse[reviewservice.GetReviewTrendResponse]
		// @Failure      500  {object}  model.ErrorResponse
		// @Router       /admin/reviews/trend [get]
		group.GET("/reviews/trend", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/reviews/trend"), h.GetReviewTrend)

		// @Summary      获取陪玩师排行榜
		// @Description  获取评价最多或评分最高的陪玩师排行榜
		// @Tags         Admin/Reviews
		// @Security     BearerAuth
		// @Produce      json
		// @Param        limit    query  int     false  "数量限制（默认10）"
		// @Param        sort_by  query  string  false  "排序方式：count（评价数量）或 rating（评分）"
		// @Success      200  {object}  model.APIResponse[reviewservice.GetTopPlayersResponse]
		// @Failure      500  {object}  model.ErrorResponse
		// @Router       /admin/reviews/top-players [get]
		group.GET("/reviews/top-players", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/reviews/top-players"), h.GetTopPlayers)

		// @Summary      获取游戏统计
		// @Description  获取各游戏的评价数量和平均评分
		// @Tags         Admin/Reviews
		// @Security     BearerAuth
		// @Produce      json
		// @Success      200  {object}  model.APIResponse[reviewservice.GetGameStatsResponse]
		// @Failure      500  {object}  model.ErrorResponse
		// @Router       /admin/reviews/game-stats [get]
		group.GET("/reviews/game-stats", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/reviews/game-stats"), h.GetGameStats)

		// @Summary      导出评价统计数据
		// @Description  导出评价统计数据为CSV格式
		// @Tags         Admin/Reviews
		// @Security     BearerAuth
		// @Produce      text/csv
		// @Param        type  query  string  false  "导出类型：overview（概览）、trend（趋势）、players（陪玩师排行）、games（游戏统计）"
		// @Param        days  query  int     false  "趋势统计天数（默认30天）"
		// @Param        limit query  int     false  "排行榜数量限制（默认10）"
		// @Success      200  {file}  file  "CSV文件"
		// @Failure      500  {object}  model.ErrorResponse
		// @Router       /admin/reviews/export [get]
		group.GET("/reviews/export", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/reviews/export"), h.ExportReviewStats)
	}
}

// RegisterContentRoutes 注册内容管理相关路由
func RegisterContentRoutes(
	group gin.IRouter,
	contentHandler *ContentHandler,
	categoryHandler *ContentCategoryHandler,
	pm *mw.PermissionMiddleware,
) {
	// 动态审核
	group.GET("/content/feeds", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/content/feeds"), contentHandler.ListFeeds)
	group.GET("/content/feeds/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/content/feeds/:id"), contentHandler.GetFeed)
	group.PUT("/content/feeds/:id/approve", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/content/feeds/:id/approve"), contentHandler.ApproveFeed)
	group.PUT("/content/feeds/:id/reject", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/content/feeds/:id/reject"), contentHandler.RejectFeed)
	group.DELETE("/content/feeds/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/content/feeds/:id"), contentHandler.DeleteFeed)
	group.POST("/content/feeds/batch-approve", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/content/feeds/batch-approve"), contentHandler.BatchApproveFeed)
	group.POST("/content/feeds/batch-reject", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/content/feeds/batch-reject"), contentHandler.BatchRejectFeed)

	// 聊天监控
	group.GET("/content/chat/messages", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/content/chat/messages"), contentHandler.ListChatMessages)
	group.DELETE("/content/chat/messages/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/content/chat/messages/:id"), contentHandler.DeleteChatMessage)
	group.POST("/content/chat/mute", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/content/chat/mute"), contentHandler.MuteUser)
	group.POST("/content/chat/unmute", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/content/chat/unmute"), contentHandler.UnmuteUser)

	// 举报管理
	group.GET("/content/reports", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/content/reports"), contentHandler.ListFeedReports)
	group.GET("/content/reports/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/content/reports/:id"), contentHandler.GetFeedReport)
	group.POST("/content/reports/:id/process", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/content/reports/:id/process"), contentHandler.ProcessFeedReport)

	// 内容统计
	group.GET("/content/stats", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/content/stats"), contentHandler.GetContentStats)
	group.GET("/content/stats/export", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/content/stats/export"), contentHandler.ExportContentStats)

	// 内容分类管理
	group.GET("/content/categories", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/content/categories"), categoryHandler.List)
	group.GET("/content/categories/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/content/categories/:id"), categoryHandler.Get)
	group.POST("/content/categories", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/content/categories"), categoryHandler.Create)
	group.PUT("/content/categories/:id", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/content/categories/:id"), categoryHandler.Update)
	group.DELETE("/content/categories/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/content/categories/:id"), categoryHandler.Delete)
}
