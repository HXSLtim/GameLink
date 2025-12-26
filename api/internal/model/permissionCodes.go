// Package model provides data models for the GameLink platform.
package model

// Permission code constants organized by module.
// Format: module.resource.action (three dot-separated segments)
// These constants are used for semantic permission identification
// and should be used instead of hardcoded strings throughout the codebase.

// ============================================================================
// Admin Module - 管理后台权限
// ============================================================================

// Admin Permission Module - 权限管理
const (
	// PermCodeAdminPermissionsRead 查看权限列表
	PermCodeAdminPermissionsRead = "admin.permissions.read"
	// PermCodeAdminPermissionsCreate 创建权限
	PermCodeAdminPermissionsCreate = "admin.permissions.create"
	// PermCodeAdminPermissionsUpdate 更新权限
	PermCodeAdminPermissionsUpdate = "admin.permissions.update"
	// PermCodeAdminPermissionsDelete 删除权限
	PermCodeAdminPermissionsDelete = "admin.permissions.delete"
)

// Admin Role Module - 角色管理
const (
	// PermCodeAdminRolesRead 查看角色列表
	PermCodeAdminRolesRead = "admin.roles.read"
	// PermCodeAdminRolesCreate 创建角色
	PermCodeAdminRolesCreate = "admin.roles.create"
	// PermCodeAdminRolesUpdate 更新角色
	PermCodeAdminRolesUpdate = "admin.roles.update"
	// PermCodeAdminRolesDelete 删除角色
	PermCodeAdminRolesDelete = "admin.roles.delete"
	// PermCodeAdminRolesAssign 分配角色权限
	PermCodeAdminRolesAssign = "admin.roles.assign"
)

// Admin User Module - 用户管理
const (
	// PermCodeAdminUsersRead 查看用户列表
	PermCodeAdminUsersRead = "admin.users.read"
	// PermCodeAdminUsersCreate 创建用户
	PermCodeAdminUsersCreate = "admin.users.create"
	// PermCodeAdminUsersUpdate 更新用户
	PermCodeAdminUsersUpdate = "admin.users.update"
	// PermCodeAdminUsersDelete 删除用户
	PermCodeAdminUsersDelete = "admin.users.delete"
	// PermCodeAdminUsersAssign 分配用户角色
	PermCodeAdminUsersAssign = "admin.users.assign"
	// PermCodeAdminUsersStatus 更新用户状态
	PermCodeAdminUsersStatus = "admin.users.status"
	// PermCodeAdminUsersPoints 管理用户积分
	PermCodeAdminUsersPoints = "admin.users.points"
)

// Admin Player Module - 陪玩师管理
const (
	// PermCodeAdminPlayersRead 查看陪玩师列表
	PermCodeAdminPlayersRead = "admin.players.read"
	// PermCodeAdminPlayersCreate 创建陪玩师
	PermCodeAdminPlayersCreate = "admin.players.create"
	// PermCodeAdminPlayersUpdate 更新陪玩师
	PermCodeAdminPlayersUpdate = "admin.players.update"
	// PermCodeAdminPlayersDelete 删除陪玩师
	PermCodeAdminPlayersDelete = "admin.players.delete"
	// PermCodeAdminPlayersVerify 审核陪玩师认证
	PermCodeAdminPlayersVerify = "admin.players.verify"
)

// Admin Game Module - 游戏管理
const (
	// PermCodeAdminGamesRead 查看游戏列表
	PermCodeAdminGamesRead = "admin.games.read"
	// PermCodeAdminGamesCreate 创建游戏
	PermCodeAdminGamesCreate = "admin.games.create"
	// PermCodeAdminGamesUpdate 更新游戏
	PermCodeAdminGamesUpdate = "admin.games.update"
	// PermCodeAdminGamesDelete 删除游戏
	PermCodeAdminGamesDelete = "admin.games.delete"
)

// Admin Order Module - 订单管理
const (
	// PermCodeAdminOrdersRead 查看订单列表
	PermCodeAdminOrdersRead = "admin.orders.read"
	// PermCodeAdminOrdersCreate 创建订单
	PermCodeAdminOrdersCreate = "admin.orders.create"
	// PermCodeAdminOrdersUpdate 更新订单
	PermCodeAdminOrdersUpdate = "admin.orders.update"
	// PermCodeAdminOrdersDelete 删除订单
	PermCodeAdminOrdersDelete = "admin.orders.delete"
	// PermCodeAdminOrdersAssign 指派订单
	PermCodeAdminOrdersAssign = "admin.orders.assign"
	// PermCodeAdminOrdersCancel 取消订单
	PermCodeAdminOrdersCancel = "admin.orders.cancel"
	// PermCodeAdminOrdersRefund 退款订单
	PermCodeAdminOrdersRefund = "admin.orders.refund"
	// PermCodeAdminOrdersConfirm 确认订单
	PermCodeAdminOrdersConfirm = "admin.orders.confirm"
	// PermCodeAdminOrdersStart 开始订单
	PermCodeAdminOrdersStart = "admin.orders.start"
	// PermCodeAdminOrdersComplete 完成订单
	PermCodeAdminOrdersComplete = "admin.orders.complete"
)

// Admin Payment Module - 支付管理
const (
	// PermCodeAdminPaymentsRead 查看支付记录
	PermCodeAdminPaymentsRead = "admin.payments.read"
	// PermCodeAdminPaymentsRefund 处理退款
	PermCodeAdminPaymentsRefund = "admin.payments.refund"
	// PermCodeAdminPaymentsExport 导出支付记录
	PermCodeAdminPaymentsExport = "admin.payments.export"
)

// Admin Withdraw Module - 提现管理
const (
	// PermCodeAdminWithdrawsRead 查看提现记录
	PermCodeAdminWithdrawsRead = "admin.withdraws.read"
	// PermCodeAdminWithdrawsApprove 审批提现
	PermCodeAdminWithdrawsApprove = "admin.withdraws.approve"
	// PermCodeAdminWithdrawsReject 拒绝提现
	PermCodeAdminWithdrawsReject = "admin.withdraws.reject"
	// PermCodeAdminWithdrawsExport 导出提现记录
	PermCodeAdminWithdrawsExport = "admin.withdraws.export"
)

// Admin Audit Module - 审计日志
const (
	// PermCodeAdminAuditRead 查看审计日志
	PermCodeAdminAuditRead = "admin.audit.read"
	// PermCodeAdminAuditExport 导出审计日志
	PermCodeAdminAuditExport = "admin.audit.export"
)

// Admin Stats Module - 统计数据
const (
	// PermCodeAdminStatsRead 查看统计数据
	PermCodeAdminStatsRead = "admin.stats.read"
	// PermCodeAdminStatsExport 导出统计数据
	PermCodeAdminStatsExport = "admin.stats.export"
)

// Admin Menu Module - 菜单管理
const (
	// PermCodeAdminMenusRead 查看菜单列表
	PermCodeAdminMenusRead = "admin.menus.read"
	// PermCodeAdminMenusCreate 创建菜单
	PermCodeAdminMenusCreate = "admin.menus.create"
	// PermCodeAdminMenusUpdate 更新菜单
	PermCodeAdminMenusUpdate = "admin.menus.update"
	// PermCodeAdminMenusDelete 删除菜单
	PermCodeAdminMenusDelete = "admin.menus.delete"
)

// ============================================================================
// Review Module - 评价管理权限
// ============================================================================

// Review Module - 评价管理
const (
	// PermCodeReviewList 查看评价列表
	PermCodeReviewList = "review.list"
	// PermCodeReviewGet 查看评价详情
	PermCodeReviewGet = "review.get"
	// PermCodeReviewPending 查看待审核评价
	PermCodeReviewPending = "review.pending"
	// PermCodeReviewLogs 查看评价操作日志
	PermCodeReviewLogs = "review.logs"
	// PermCodeReviewPlayer 查看陪玩师评价
	PermCodeReviewPlayer = "review.player"
	// PermCodeReviewOrder 查看订单评价
	PermCodeReviewOrder = "review.order"
	// PermCodeReviewApprove 批准评价
	PermCodeReviewApprove = "review.approve"
	// PermCodeReviewReject 拒绝评价
	PermCodeReviewReject = "review.reject"
	// PermCodeReviewBatchApprove 批量批准评价
	PermCodeReviewBatchApprove = "review.batch_approve"
	// PermCodeReviewBatchReject 批量拒绝评价
	PermCodeReviewBatchReject = "review.batch_reject"
	// PermCodeReviewDelete 删除评价
	PermCodeReviewDelete = "review.delete"
	// PermCodeReviewUpdate 更新评价
	PermCodeReviewUpdate = "review.update"
	// PermCodeReviewCreate 创建评价
	PermCodeReviewCreate = "review.create"
	// PermCodeReviewStats 查看评价统计
	PermCodeReviewStats = "review.stats"
	// PermCodeReviewTrend 查看评价趋势
	PermCodeReviewTrend = "review.trend"
	// PermCodeReviewTopPlayers 查看陪玩师排行榜
	PermCodeReviewTopPlayers = "review.top_players"
	// PermCodeReviewGameStats 查看游戏评价统计
	PermCodeReviewGameStats = "review.game_stats"
	// PermCodeReviewExport 导出评价统计
	PermCodeReviewExport = "review.export"
	// PermCodeReviewDetectSensitive 检测敏感词
	PermCodeReviewDetectSensitive = "review.detect_sensitive"
)

// Review Report Module - 评价举报管理
const (
	// PermCodeReviewReportList 查看评价举报列表
	PermCodeReviewReportList = "review_report.list"
	// PermCodeReviewReportGet 查看举报详情
	PermCodeReviewReportGet = "review_report.get"
	// PermCodeReviewReportCreate 创建评价举报
	PermCodeReviewReportCreate = "review_report.create"
	// PermCodeReviewReportHandle 处理评价举报
	PermCodeReviewReportHandle = "review_report.handle"
)

// Review Reply Module - 评价回复管理
const (
	// PermCodeReviewReplyUpdate 更新评价回复
	PermCodeReviewReplyUpdate = "review_reply.update"
	// PermCodeReviewReplyDelete 删除评价回复
	PermCodeReviewReplyDelete = "review_reply.delete"
)

// Review Settings Module - 评价展示设置
const (
	// PermCodeReviewSettingsGet 查看评价展示设置
	PermCodeReviewSettingsGet = "review_settings.get"
	// PermCodeReviewSettingsUpdate 更新评价展示设置
	PermCodeReviewSettingsUpdate = "review_settings.update"
)

// Sensitive Word Module - 敏感词管理
const (
	// PermCodeSensitiveWordList 查看敏感词列表
	PermCodeSensitiveWordList = "sensitive_word.list"
	// PermCodeSensitiveWordCreate 添加敏感词
	PermCodeSensitiveWordCreate = "sensitive_word.create"
	// PermCodeSensitiveWordUpdate 更新敏感词
	PermCodeSensitiveWordUpdate = "sensitive_word.update"
	// PermCodeSensitiveWordDelete 删除敏感词
	PermCodeSensitiveWordDelete = "sensitive_word.delete"
)

// Operation Log Module - 操作日志
const (
	// PermCodeOperationLogList 查看操作日志
	PermCodeOperationLogList = "operation_log.list"
	// PermCodeOperationLogExport 导出操作日志
	PermCodeOperationLogExport = "operation_log.export"
)

// ============================================================================
// Content Module - 内容管理权限
// ============================================================================

// Content Feed Module - 动态管理
const (
	// PermCodeContentFeedList 查看动态列表
	PermCodeContentFeedList = "content.feed.list"
	// PermCodeContentFeedGet 查看动态详情
	PermCodeContentFeedGet = "content.feed.get"
	// PermCodeContentFeedApprove 批准动态
	PermCodeContentFeedApprove = "content.feed.approve"
	// PermCodeContentFeedReject 拒绝动态
	PermCodeContentFeedReject = "content.feed.reject"
	// PermCodeContentFeedBatchApprove 批量批准动态
	PermCodeContentFeedBatchApprove = "content.feed.batch_approve"
	// PermCodeContentFeedBatchReject 批量拒绝动态
	PermCodeContentFeedBatchReject = "content.feed.batch_reject"
	// PermCodeContentFeedDelete 删除动态
	PermCodeContentFeedDelete = "content.feed.delete"
)

// Content Chat Module - 聊天监控
const (
	// PermCodeContentChatList 查看聊天消息
	PermCodeContentChatList = "content.chat.list"
	// PermCodeContentChatDelete 删除聊天消息
	PermCodeContentChatDelete = "content.chat.delete"
	// PermCodeContentChatMute 禁言用户
	PermCodeContentChatMute = "content.chat.mute"
	// PermCodeContentChatUnmute 解除禁言
	PermCodeContentChatUnmute = "content.chat.unmute"
)

// Content Report Module - 举报管理
const (
	// PermCodeContentReportList 查看举报列表
	PermCodeContentReportList = "content.report.list"
	// PermCodeContentReportGet 查看举报详情
	PermCodeContentReportGet = "content.report.get"
	// PermCodeContentReportProcess 处理举报
	PermCodeContentReportProcess = "content.report.process"
)

// Content Stats Module - 内容统计
const (
	// PermCodeContentStats 查看内容统计
	PermCodeContentStats = "content.stats"
)

// Content Category Module - 内容分类管理
const (
	// PermCodeContentCategoryList 查看内容分类
	PermCodeContentCategoryList = "content.category.list"
	// PermCodeContentCategoryGet 查看分类详情
	PermCodeContentCategoryGet = "content.category.get"
	// PermCodeContentCategoryCreate 创建内容分类
	PermCodeContentCategoryCreate = "content.category.create"
	// PermCodeContentCategoryUpdate 更新内容分类
	PermCodeContentCategoryUpdate = "content.category.update"
	// PermCodeContentCategoryDelete 删除内容分类
	PermCodeContentCategoryDelete = "content.category.delete"
)

// ============================================================================
// Commission Module - 佣金管理权限
// ============================================================================

// Commission Module - 佣金管理
const (
	// PermCodeCommissionRead 查看佣金记录
	PermCodeCommissionRead = "commission.read"
	// PermCodeCommissionSettle 结算佣金
	PermCodeCommissionSettle = "commission.settle"
	// PermCodeCommissionExport 导出佣金记录
	PermCodeCommissionExport = "commission.export"
)

// ============================================================================
// Notification Module - 通知管理权限
// ============================================================================

// Notification Module - 通知管理
const (
	// PermCodeNotificationRead 查看通知
	PermCodeNotificationRead = "notification.read"
	// PermCodeNotificationCreate 创建通知
	PermCodeNotificationCreate = "notification.create"
	// PermCodeNotificationBatchSend 批量发送通知
	PermCodeNotificationBatchSend = "notification.batch_send"
)

// ============================================================================
// Wallet Module - 钱包管理权限
// ============================================================================

// Wallet Module - 钱包管理
const (
	// PermCodeWalletRead 查看钱包
	PermCodeWalletRead = "wallet.read"
	// PermCodeWalletUpdate 更新钱包
	PermCodeWalletUpdate = "wallet.update"
	// PermCodeWalletTransactions 查看交易记录
	PermCodeWalletTransactions = "wallet.transactions"
)

// ============================================================================
// Service Item Module - 服务项目管理权限
// ============================================================================

// Service Item Module - 服务项目管理
const (
	// PermCodeServiceItemRead 查看服务项目
	PermCodeServiceItemRead = "service_item.read"
	// PermCodeServiceItemCreate 创建服务项目
	PermCodeServiceItemCreate = "service_item.create"
	// PermCodeServiceItemUpdate 更新服务项目
	PermCodeServiceItemUpdate = "service_item.update"
	// PermCodeServiceItemDelete 删除服务项目
	PermCodeServiceItemDelete = "service_item.delete"
)

// ============================================================================
// Dispute Module - 纠纷管理权限
// ============================================================================

// Dispute Module - 纠纷管理
const (
	// PermCodeDisputeRead 查看纠纷
	PermCodeDisputeRead = "dispute.read"
	// PermCodeDisputeCreate 创建纠纷
	PermCodeDisputeCreate = "dispute.create"
	// PermCodeDisputeResolve 解决纠纷
	PermCodeDisputeResolve = "dispute.resolve"
)

// ============================================================================
// System Module - 系统管理权限
// ============================================================================

// System Module - 系统管理
const (
	// PermCodeSystemConfigRead 查看系统配置
	PermCodeSystemConfigRead = "system.config.read"
	// PermCodeSystemConfigUpdate 更新系统配置
	PermCodeSystemConfigUpdate = "system.config.update"
	// PermCodeSystemMonitorRead 查看系统监控
	PermCodeSystemMonitorRead = "system.monitor.read"
)

// ============================================================================
// Permission Groups - 权限分组定义
// ============================================================================

// PermissionGroupInfo represents a permission group with its metadata.
type PermissionGroupInfo struct {
	// Group is the group identifier (e.g., "/admin/permissions")
	Group string
	// Name is the display name of the group
	Name string
	// Description is a brief description of the group
	Description string
	// Module is the module this group belongs to
	Module string
}

// Permission groups organized by module
var (
	// PermGroupAdminPermissions 权限管理分组
	PermGroupAdminPermissions = PermissionGroupInfo{
		Group:       "/admin/permissions",
		Name:        "权限管理",
		Description: "管理系统权限定义",
		Module:      "admin",
	}

	// PermGroupAdminRoles 角色管理分组
	PermGroupAdminRoles = PermissionGroupInfo{
		Group:       "/admin/roles",
		Name:        "角色管理",
		Description: "管理系统角色和权限分配",
		Module:      "admin",
	}

	// PermGroupAdminUsers 用户管理分组
	PermGroupAdminUsers = PermissionGroupInfo{
		Group:       "/admin/users",
		Name:        "用户管理",
		Description: "管理平台用户",
		Module:      "admin",
	}

	// PermGroupAdminPlayers 陪玩师管理分组
	PermGroupAdminPlayers = PermissionGroupInfo{
		Group:       "/admin/players",
		Name:        "陪玩师管理",
		Description: "管理陪玩师资料和认证",
		Module:      "admin",
	}

	// PermGroupAdminGames 游戏管理分组
	PermGroupAdminGames = PermissionGroupInfo{
		Group:       "/admin/games",
		Name:        "游戏管理",
		Description: "管理游戏列表",
		Module:      "admin",
	}

	// PermGroupAdminOrders 订单管理分组
	PermGroupAdminOrders = PermissionGroupInfo{
		Group:       "/admin/orders",
		Name:        "订单管理",
		Description: "管理订单和订单状态",
		Module:      "admin",
	}

	// PermGroupAdminPayments 支付管理分组
	PermGroupAdminPayments = PermissionGroupInfo{
		Group:       "/admin/payments",
		Name:        "支付管理",
		Description: "管理支付和退款",
		Module:      "admin",
	}

	// PermGroupAdminWithdraws 提现管理分组
	PermGroupAdminWithdraws = PermissionGroupInfo{
		Group:       "/admin/withdraws",
		Name:        "提现管理",
		Description: "管理提现申请和审批",
		Module:      "admin",
	}

	// PermGroupAdminAudit 审计日志分组
	PermGroupAdminAudit = PermissionGroupInfo{
		Group:       "/admin/audit",
		Name:        "审计日志",
		Description: "查看和导出审计日志",
		Module:      "admin",
	}

	// PermGroupAdminStats 统计数据分组
	PermGroupAdminStats = PermissionGroupInfo{
		Group:       "/admin/stats",
		Name:        "统计数据",
		Description: "查看平台统计数据",
		Module:      "admin",
	}

	// PermGroupAdminMenus 菜单管理分组
	PermGroupAdminMenus = PermissionGroupInfo{
		Group:       "/admin/menus",
		Name:        "菜单管理",
		Description: "管理后台菜单",
		Module:      "admin",
	}

	// PermGroupReviews 评价管理分组
	PermGroupReviews = PermissionGroupInfo{
		Group:       "/admin/reviews",
		Name:        "评价管理",
		Description: "管理用户评价和审核",
		Module:      "review",
	}

	// PermGroupContent 内容管理分组
	PermGroupContent = PermissionGroupInfo{
		Group:       "/admin/content",
		Name:        "内容管理",
		Description: "管理动态、聊天和举报",
		Module:      "content",
	}

	// PermGroupCommission 佣金管理分组
	PermGroupCommission = PermissionGroupInfo{
		Group:       "/admin/commission",
		Name:        "佣金管理",
		Description: "管理佣金结算",
		Module:      "commission",
	}

	// PermGroupNotification 通知管理分组
	PermGroupNotification = PermissionGroupInfo{
		Group:       "/admin/notification",
		Name:        "通知管理",
		Description: "管理系统通知",
		Module:      "notification",
	}

	// PermGroupWallet 钱包管理分组
	PermGroupWallet = PermissionGroupInfo{
		Group:       "/admin/wallet",
		Name:        "钱包管理",
		Description: "管理用户钱包",
		Module:      "wallet",
	}

	// PermGroupServiceItem 服务项目分组
	PermGroupServiceItem = PermissionGroupInfo{
		Group:       "/admin/service-items",
		Name:        "服务项目管理",
		Description: "管理服务项目",
		Module:      "service_item",
	}

	// PermGroupDispute 纠纷管理分组
	PermGroupDispute = PermissionGroupInfo{
		Group:       "/admin/disputes",
		Name:        "纠纷管理",
		Description: "管理订单纠纷",
		Module:      "dispute",
	}

	// PermGroupSystem 系统管理分组
	PermGroupSystem = PermissionGroupInfo{
		Group:       "/admin/system",
		Name:        "系统管理",
		Description: "系统配置和监控",
		Module:      "system",
	}
)

// AllPermissionGroups returns all permission groups.
func AllPermissionGroups() []PermissionGroupInfo {
	return []PermissionGroupInfo{
		PermGroupAdminPermissions,
		PermGroupAdminRoles,
		PermGroupAdminUsers,
		PermGroupAdminPlayers,
		PermGroupAdminGames,
		PermGroupAdminOrders,
		PermGroupAdminPayments,
		PermGroupAdminWithdraws,
		PermGroupAdminAudit,
		PermGroupAdminStats,
		PermGroupAdminMenus,
		PermGroupReviews,
		PermGroupContent,
		PermGroupCommission,
		PermGroupNotification,
		PermGroupWallet,
		PermGroupServiceItem,
		PermGroupDispute,
		PermGroupSystem,
	}
}

// PermissionDefinition represents a complete permission definition.
type PermissionDefinition struct {
	Code        string
	Method      HTTPMethod
	Path        string
	Group       string
	Description string
	IsSystem    bool
}

// GetAllPermissionDefinitions returns all permission definitions.
// This is used for seeding permissions in the database.
func GetAllPermissionDefinitions() []PermissionDefinition {
	return []PermissionDefinition{
		// ========== Admin Permission Management ==========
		{Code: PermCodeAdminPermissionsRead, Method: HTTPMethodGET, Path: "/api/v1/admin/permissions", Group: PermGroupAdminPermissions.Group, Description: "查看权限列表", IsSystem: true},
		{Code: PermCodeAdminPermissionsRead, Method: HTTPMethodGET, Path: "/api/v1/admin/permissions/:id", Group: PermGroupAdminPermissions.Group, Description: "查看权限详情", IsSystem: true},
		{Code: PermCodeAdminPermissionsRead, Method: HTTPMethodGET, Path: "/api/v1/admin/permissions/tree", Group: PermGroupAdminPermissions.Group, Description: "查看权限树", IsSystem: true},
		{Code: PermCodeAdminPermissionsRead, Method: HTTPMethodGET, Path: "/api/v1/admin/permissions/groups", Group: PermGroupAdminPermissions.Group, Description: "查看权限分组", IsSystem: true},
		{Code: PermCodeAdminPermissionsCreate, Method: HTTPMethodPOST, Path: "/api/v1/admin/permissions", Group: PermGroupAdminPermissions.Group, Description: "创建权限", IsSystem: true},
		{Code: PermCodeAdminPermissionsUpdate, Method: HTTPMethodPUT, Path: "/api/v1/admin/permissions/:id", Group: PermGroupAdminPermissions.Group, Description: "更新权限", IsSystem: true},
		{Code: PermCodeAdminPermissionsUpdate, Method: HTTPMethodPATCH, Path: "/api/v1/admin/permissions/:id", Group: PermGroupAdminPermissions.Group, Description: "部分更新权限", IsSystem: true},
		{Code: PermCodeAdminPermissionsDelete, Method: HTTPMethodDELETE, Path: "/api/v1/admin/permissions/:id", Group: PermGroupAdminPermissions.Group, Description: "删除权限", IsSystem: true},

		// ========== Admin Role Management ==========
		{Code: PermCodeAdminRolesRead, Method: HTTPMethodGET, Path: "/api/v1/admin/roles", Group: PermGroupAdminRoles.Group, Description: "查看角色列表", IsSystem: true},
		{Code: PermCodeAdminRolesRead, Method: HTTPMethodGET, Path: "/api/v1/admin/roles/:id", Group: PermGroupAdminRoles.Group, Description: "查看角色详情", IsSystem: true},
		{Code: PermCodeAdminRolesRead, Method: HTTPMethodGET, Path: "/api/v1/admin/roles/:id/permissions", Group: PermGroupAdminRoles.Group, Description: "查看角色权限", IsSystem: true},
		{Code: PermCodeAdminRolesCreate, Method: HTTPMethodPOST, Path: "/api/v1/admin/roles", Group: PermGroupAdminRoles.Group, Description: "创建角色", IsSystem: true},
		{Code: PermCodeAdminRolesUpdate, Method: HTTPMethodPUT, Path: "/api/v1/admin/roles/:id", Group: PermGroupAdminRoles.Group, Description: "更新角色", IsSystem: true},
		{Code: PermCodeAdminRolesUpdate, Method: HTTPMethodPATCH, Path: "/api/v1/admin/roles/:id", Group: PermGroupAdminRoles.Group, Description: "部分更新角色", IsSystem: true},
		{Code: PermCodeAdminRolesDelete, Method: HTTPMethodDELETE, Path: "/api/v1/admin/roles/:id", Group: PermGroupAdminRoles.Group, Description: "删除角色", IsSystem: true},
		{Code: PermCodeAdminRolesAssign, Method: HTTPMethodPUT, Path: "/api/v1/admin/roles/:id/permissions/batch", Group: PermGroupAdminRoles.Group, Description: "批量分配角色权限", IsSystem: true},
		{Code: PermCodeAdminRolesAssign, Method: HTTPMethodPOST, Path: "/api/v1/admin/roles/:id/permissions/:pid", Group: PermGroupAdminRoles.Group, Description: "添加角色权限", IsSystem: true},
		{Code: PermCodeAdminRolesAssign, Method: HTTPMethodDELETE, Path: "/api/v1/admin/roles/:id/permissions/:pid", Group: PermGroupAdminRoles.Group, Description: "移除角色权限", IsSystem: true},

		// ========== Admin User Management ==========
		{Code: PermCodeAdminUsersRead, Method: HTTPMethodGET, Path: "/api/v1/admin/users", Group: PermGroupAdminUsers.Group, Description: "查看用户列表", IsSystem: true},
		{Code: PermCodeAdminUsersRead, Method: HTTPMethodGET, Path: "/api/v1/admin/users/stats", Group: PermGroupAdminUsers.Group, Description: "查看用户统计", IsSystem: true},
		{Code: PermCodeAdminUsersRead, Method: HTTPMethodGET, Path: "/api/v1/admin/users/:id", Group: PermGroupAdminUsers.Group, Description: "查看用户详情", IsSystem: true},
		{Code: PermCodeAdminUsersRead, Method: HTTPMethodGET, Path: "/api/v1/admin/users/:id/orders", Group: PermGroupAdminUsers.Group, Description: "查看用户订单", IsSystem: true},
		{Code: PermCodeAdminUsersRead, Method: HTTPMethodGET, Path: "/api/v1/admin/users/:id/logs", Group: PermGroupAdminUsers.Group, Description: "查看用户操作日志", IsSystem: true},
		{Code: PermCodeAdminUsersRead, Method: HTTPMethodGET, Path: "/api/v1/admin/users/:id/login-history", Group: PermGroupAdminUsers.Group, Description: "查看用户登录历史", IsSystem: true},
		{Code: PermCodeAdminUsersRead, Method: HTTPMethodGET, Path: "/api/v1/admin/users/:id/roles", Group: PermGroupAdminUsers.Group, Description: "查看用户角色", IsSystem: true},
		{Code: PermCodeAdminUsersRead, Method: HTTPMethodGET, Path: "/api/v1/admin/users/:id/permissions", Group: PermGroupAdminUsers.Group, Description: "查看用户权限", IsSystem: true},
		{Code: PermCodeAdminUsersCreate, Method: HTTPMethodPOST, Path: "/api/v1/admin/users", Group: PermGroupAdminUsers.Group, Description: "创建用户", IsSystem: true},
		{Code: PermCodeAdminUsersCreate, Method: HTTPMethodPOST, Path: "/api/v1/admin/users/with-player", Group: PermGroupAdminUsers.Group, Description: "创建用户和陪玩师", IsSystem: true},
		{Code: PermCodeAdminUsersUpdate, Method: HTTPMethodPUT, Path: "/api/v1/admin/users/:id", Group: PermGroupAdminUsers.Group, Description: "更新用户", IsSystem: true},
		{Code: PermCodeAdminUsersDelete, Method: HTTPMethodDELETE, Path: "/api/v1/admin/users/:id", Group: PermGroupAdminUsers.Group, Description: "删除用户", IsSystem: true},
		{Code: PermCodeAdminUsersDelete, Method: HTTPMethodPOST, Path: "/api/v1/admin/users/batch-delete", Group: PermGroupAdminUsers.Group, Description: "批量删除用户", IsSystem: true},
		{Code: PermCodeAdminUsersStatus, Method: HTTPMethodPUT, Path: "/api/v1/admin/users/:id/status", Group: PermGroupAdminUsers.Group, Description: "更新用户状态", IsSystem: true},
		{Code: PermCodeAdminUsersAssign, Method: HTTPMethodPUT, Path: "/api/v1/admin/users/:id/role", Group: PermGroupAdminUsers.Group, Description: "更新用户角色", IsSystem: true},
		{Code: PermCodeAdminUsersAssign, Method: HTTPMethodPUT, Path: "/api/v1/admin/users/:id/roles", Group: PermGroupAdminUsers.Group, Description: "分配用户角色", IsSystem: true},
		{Code: PermCodeAdminUsersPoints, Method: HTTPMethodPOST, Path: "/api/v1/admin/users/batch-add-points", Group: PermGroupAdminUsers.Group, Description: "批量添加积分", IsSystem: true},

		// ========== Admin Player Management ==========
		{Code: PermCodeAdminPlayersRead, Method: HTTPMethodGET, Path: "/api/v1/admin/players", Group: PermGroupAdminPlayers.Group, Description: "查看陪玩师列表", IsSystem: true},
		{Code: PermCodeAdminPlayersRead, Method: HTTPMethodGET, Path: "/api/v1/admin/players/:id", Group: PermGroupAdminPlayers.Group, Description: "查看陪玩师详情", IsSystem: true},
		{Code: PermCodeAdminPlayersRead, Method: HTTPMethodGET, Path: "/api/v1/admin/players/:id/logs", Group: PermGroupAdminPlayers.Group, Description: "查看陪玩师操作日志", IsSystem: true},
		{Code: PermCodeAdminPlayersCreate, Method: HTTPMethodPOST, Path: "/api/v1/admin/players", Group: PermGroupAdminPlayers.Group, Description: "创建陪玩师", IsSystem: true},
		{Code: PermCodeAdminPlayersUpdate, Method: HTTPMethodPUT, Path: "/api/v1/admin/players/:id", Group: PermGroupAdminPlayers.Group, Description: "更新陪玩师", IsSystem: true},
		{Code: PermCodeAdminPlayersDelete, Method: HTTPMethodDELETE, Path: "/api/v1/admin/players/:id", Group: PermGroupAdminPlayers.Group, Description: "删除陪玩师", IsSystem: true},
		{Code: PermCodeAdminPlayersVerify, Method: HTTPMethodPUT, Path: "/api/v1/admin/players/:id/verification", Group: PermGroupAdminPlayers.Group, Description: "审核陪玩师认证", IsSystem: true},
		{Code: PermCodeAdminPlayersUpdate, Method: HTTPMethodPUT, Path: "/api/v1/admin/players/:id/games", Group: PermGroupAdminPlayers.Group, Description: "更新陪玩师游戏", IsSystem: true},
		{Code: PermCodeAdminPlayersUpdate, Method: HTTPMethodPUT, Path: "/api/v1/admin/players/:id/skill-tags", Group: PermGroupAdminPlayers.Group, Description: "更新陪玩师技能标签", IsSystem: true},

		// ========== Admin Game Management ==========
		{Code: PermCodeAdminGamesRead, Method: HTTPMethodGET, Path: "/api/v1/admin/games", Group: PermGroupAdminGames.Group, Description: "查看游戏列表", IsSystem: true},
		{Code: PermCodeAdminGamesRead, Method: HTTPMethodGET, Path: "/api/v1/admin/games/:id", Group: PermGroupAdminGames.Group, Description: "查看游戏详情", IsSystem: true},
		{Code: PermCodeAdminGamesRead, Method: HTTPMethodGET, Path: "/api/v1/admin/games/:id/logs", Group: PermGroupAdminGames.Group, Description: "查看游戏操作日志", IsSystem: true},
		{Code: PermCodeAdminGamesCreate, Method: HTTPMethodPOST, Path: "/api/v1/admin/games", Group: PermGroupAdminGames.Group, Description: "创建游戏", IsSystem: true},
		{Code: PermCodeAdminGamesUpdate, Method: HTTPMethodPUT, Path: "/api/v1/admin/games/:id", Group: PermGroupAdminGames.Group, Description: "更新游戏", IsSystem: true},
		{Code: PermCodeAdminGamesDelete, Method: HTTPMethodDELETE, Path: "/api/v1/admin/games/:id", Group: PermGroupAdminGames.Group, Description: "删除游戏", IsSystem: true},

		// ========== Admin Order Management ==========
		{Code: PermCodeAdminOrdersRead, Method: HTTPMethodGET, Path: "/api/v1/admin/orders", Group: PermGroupAdminOrders.Group, Description: "查看订单列表", IsSystem: true},
		{Code: PermCodeAdminOrdersRead, Method: HTTPMethodGET, Path: "/api/v1/admin/orders/:id", Group: PermGroupAdminOrders.Group, Description: "查看订单详情", IsSystem: true},
		{Code: PermCodeAdminOrdersRead, Method: HTTPMethodGET, Path: "/api/v1/admin/orders/:id/logs", Group: PermGroupAdminOrders.Group, Description: "查看订单操作日志", IsSystem: true},
		{Code: PermCodeAdminOrdersRead, Method: HTTPMethodGET, Path: "/api/v1/admin/orders/:id/timeline", Group: PermGroupAdminOrders.Group, Description: "查看订单时间线", IsSystem: true},
		{Code: PermCodeAdminOrdersRead, Method: HTTPMethodGET, Path: "/api/v1/admin/orders/:id/payments", Group: PermGroupAdminOrders.Group, Description: "查看订单支付记录", IsSystem: true},
		{Code: PermCodeAdminOrdersRead, Method: HTTPMethodGET, Path: "/api/v1/admin/orders/:id/refunds", Group: PermGroupAdminOrders.Group, Description: "查看订单退款记录", IsSystem: true},
		{Code: PermCodeAdminOrdersRead, Method: HTTPMethodGET, Path: "/api/v1/admin/orders/:id/reviews", Group: PermGroupAdminOrders.Group, Description: "查看订单评价", IsSystem: true},
		{Code: PermCodeAdminOrdersCreate, Method: HTTPMethodPOST, Path: "/api/v1/admin/orders", Group: PermGroupAdminOrders.Group, Description: "创建订单", IsSystem: true},
		{Code: PermCodeAdminOrdersUpdate, Method: HTTPMethodPUT, Path: "/api/v1/admin/orders/:id", Group: PermGroupAdminOrders.Group, Description: "更新订单", IsSystem: true},
		{Code: PermCodeAdminOrdersDelete, Method: HTTPMethodDELETE, Path: "/api/v1/admin/orders/:id", Group: PermGroupAdminOrders.Group, Description: "删除订单", IsSystem: true},
		{Code: PermCodeAdminOrdersAssign, Method: HTTPMethodPOST, Path: "/api/v1/admin/orders/:id/assign", Group: PermGroupAdminOrders.Group, Description: "指派订单", IsSystem: true},
		{Code: PermCodeAdminOrdersCancel, Method: HTTPMethodPOST, Path: "/api/v1/admin/orders/:id/cancel", Group: PermGroupAdminOrders.Group, Description: "取消订单", IsSystem: true},
		{Code: PermCodeAdminOrdersRefund, Method: HTTPMethodPOST, Path: "/api/v1/admin/orders/:id/refund", Group: PermGroupAdminOrders.Group, Description: "退款订单", IsSystem: true},
		{Code: PermCodeAdminOrdersConfirm, Method: HTTPMethodPOST, Path: "/api/v1/admin/orders/:id/confirm", Group: PermGroupAdminOrders.Group, Description: "确认订单", IsSystem: true},
		{Code: PermCodeAdminOrdersStart, Method: HTTPMethodPOST, Path: "/api/v1/admin/orders/:id/start", Group: PermGroupAdminOrders.Group, Description: "开始订单", IsSystem: true},
		{Code: PermCodeAdminOrdersComplete, Method: HTTPMethodPOST, Path: "/api/v1/admin/orders/:id/complete", Group: PermGroupAdminOrders.Group, Description: "完成订单", IsSystem: true},
		{Code: PermCodeAdminOrdersCreate, Method: HTTPMethodPOST, Path: "/api/v1/admin/orders/:id/review", Group: PermGroupAdminOrders.Group, Description: "评价订单", IsSystem: true},

		// ========== Admin Audit Log ==========
		{Code: PermCodeAdminAuditRead, Method: HTTPMethodGET, Path: "/api/v1/admin/audit/permissions", Group: PermGroupAdminAudit.Group, Description: "查看权限审计日志", IsSystem: true},
		{Code: PermCodeAdminAuditExport, Method: HTTPMethodGET, Path: "/api/v1/admin/audit/permissions/export", Group: PermGroupAdminAudit.Group, Description: "导出权限审计日志", IsSystem: true},

		// ========== Admin Stats ==========
		{Code: PermCodeAdminStatsRead, Method: HTTPMethodGET, Path: "/api/v1/admin/stats", Group: PermGroupAdminStats.Group, Description: "查看统计数据", IsSystem: true},
		{Code: PermCodeAdminStatsRead, Method: HTTPMethodGET, Path: "/api/v1/admin/users/behavior/stats", Group: PermGroupAdminStats.Group, Description: "查看用户行为统计", IsSystem: true},
		{Code: PermCodeAdminStatsRead, Method: HTTPMethodGET, Path: "/api/v1/admin/users/behavior/trend", Group: PermGroupAdminStats.Group, Description: "查看用户活动趋势", IsSystem: true},
		{Code: PermCodeAdminStatsRead, Method: HTTPMethodGET, Path: "/api/v1/admin/users/behavior/distribution", Group: PermGroupAdminStats.Group, Description: "查看用户分布", IsSystem: true},
	}
}

// GetPermissionCodesByModule returns all permission codes for a specific module.
func GetPermissionCodesByModule(module string) []string {
	allDefs := GetAllPermissionDefinitions()
	codeSet := make(map[string]bool)
	var codes []string

	for _, def := range allDefs {
		// Extract module from code (first segment)
		parts := splitPermissionCode(def.Code)
		if len(parts) > 0 && parts[0] == module {
			if !codeSet[def.Code] {
				codeSet[def.Code] = true
				codes = append(codes, def.Code)
			}
		}
	}

	return codes
}

// splitPermissionCode splits a permission code into its parts.
func splitPermissionCode(code string) []string {
	if code == "" {
		return nil
	}
	var parts []string
	start := 0
	for i, c := range code {
		if c == '.' {
			parts = append(parts, code[start:i])
			start = i + 1
		}
	}
	if start < len(code) {
		parts = append(parts, code[start:])
	}
	return parts
}

// IsValidPermissionCode checks if a permission code is valid.
func IsValidPermissionCode(code string) bool {
	return PermissionCodePattern.MatchString(code)
}

// SuperAdminPermissionCode is the special permission code that grants all permissions.
const SuperAdminPermissionCode = "*"

// IsSuperAdminPermission checks if the permission code is the super admin wildcard.
func IsSuperAdminPermission(code string) bool {
	return code == SuperAdminPermissionCode
}
