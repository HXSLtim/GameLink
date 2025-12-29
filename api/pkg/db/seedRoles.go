package db

import (
	"log"

	"gorm.io/gorm"

	"gamelink/internal/model"
)

// DefaultRolePermissions defines the default permission codes for each system role.
// superAdmin: All permissions (handled specially with "*")
// admin: Management permissions (users, players, orders, payments, content, reviews)
// finance: Financial permissions (withdraws, payments, refunds)
// customerService: Customer service permissions (disputes, orders, users)
// player: Player-specific permissions (view own data, manage services)
// user: Basic user permissions (view public data)
var DefaultRolePermissions = map[model.RoleSlug][]string{
	model.RoleSlugSuperAdmin: {
		// SuperAdmin gets all permissions via "*" wildcard
		model.SuperAdminPermissionCode,
	},
	// ============================================================
	// 财务角色 - Finance
	// ============================================================
	model.RoleSlugFinance: {
		// Withdraw Management (提现管理)
		model.PermCodeAdminWithdrawsRead,
		model.PermCodeAdminWithdrawsApprove,
		model.PermCodeAdminWithdrawsReject,
		model.PermCodeAdminWithdrawsExport,

		// Payment Management (支付管理)
		model.PermCodeAdminPaymentsRead,
		model.PermCodeAdminPaymentsRefund,
		model.PermCodeAdminPaymentsExport,

		// Order Refund (订单退款)
		model.PermCodeAdminOrdersRead,
		model.PermCodeAdminOrdersRefund,

		// Wallet Management (钱包管理)
		model.PermCodeWalletRead,
		model.PermCodeWalletTransactions,

		// Stats (统计)
		model.PermCodeAdminStatsRead,
		model.PermCodeAdminStatsExport,

		// Commission (佣金)
		model.PermCodeCommissionRead,
		model.PermCodeCommissionExport,

		// Operation Log
		model.PermCodeOperationLogList,
	},
	// ============================================================
	// 客服角色 - Customer Service
	// ============================================================
	model.RoleSlugCustomerService: {
		// User Management (用户管理 - 只读和状态更新)
		model.PermCodeAdminUsersRead,
		model.PermCodeAdminUsersStatus,

		// Player Management (陪玩师管理 - 只读)
		model.PermCodeAdminPlayersRead,

		// Order Management (订单管理 - 查看、确认、取消)
		model.PermCodeAdminOrdersRead,
		model.PermCodeAdminOrdersConfirm,
		model.PermCodeAdminOrdersCancel,
		model.PermCodeAdminOrdersAssign,

		// Payment Management (支付管理 - 只读)
		model.PermCodeAdminPaymentsRead,

		// Dispute Management (纠纷管理)
		model.PermCodeDisputeRead,
		model.PermCodeDisputeCreate,
		model.PermCodeDisputeResolve,

		// Review Management (评论管理 - 审核)
		model.PermCodeReviewList,
		model.PermCodeReviewGet,
		model.PermCodeReviewPending,
		model.PermCodeReviewApprove,
		model.PermCodeReviewReject,
		model.PermCodeReviewBatchApprove,
		model.PermCodeReviewBatchReject,

		// Content Management (内容管理 - 审核)
		model.PermCodeContentFeedList,
		model.PermCodeContentFeedGet,
		model.PermCodeContentFeedApprove,
		model.PermCodeContentFeedReject,
		model.PermCodeContentFeedBatchApprove,
		model.PermCodeContentFeedBatchReject,

		// Chat Monitoring (聊天监控)
		model.PermCodeContentChatList,

		// Content Report (内容举报)
		model.PermCodeContentReportList,
		model.PermCodeContentReportGet,
		model.PermCodeContentReportProcess,

		// Notification (通知)
		model.PermCodeNotificationRead,
		model.PermCodeNotificationCreate,
		model.PermCodeNotificationBatchSend,

		// Stats (统计)
		model.PermCodeAdminStatsRead,

		// Operation Log
		model.PermCodeOperationLogList,
	},
	model.RoleSlugAdmin: {
		// User Management
		model.PermCodeAdminUsersRead,
		model.PermCodeAdminUsersCreate,
		model.PermCodeAdminUsersUpdate,
		model.PermCodeAdminUsersDelete,
		model.PermCodeAdminUsersAssign,
		model.PermCodeAdminUsersStatus,
		model.PermCodeAdminUsersPoints,

		// Player Management
		model.PermCodeAdminPlayersRead,
		model.PermCodeAdminPlayersCreate,
		model.PermCodeAdminPlayersUpdate,
		model.PermCodeAdminPlayersDelete,
		model.PermCodeAdminPlayersVerify,

		// Game Management
		model.PermCodeAdminGamesRead,
		model.PermCodeAdminGamesCreate,
		model.PermCodeAdminGamesUpdate,
		model.PermCodeAdminGamesDelete,

		// Order Management
		model.PermCodeAdminOrdersRead,
		model.PermCodeAdminOrdersCreate,
		model.PermCodeAdminOrdersUpdate,
		model.PermCodeAdminOrdersAssign,
		model.PermCodeAdminOrdersCancel,
		model.PermCodeAdminOrdersRefund,
		model.PermCodeAdminOrdersConfirm,
		model.PermCodeAdminOrdersStart,
		model.PermCodeAdminOrdersComplete,

		// Payment Management
		model.PermCodeAdminPaymentsRead,
		model.PermCodeAdminPaymentsRefund,
		model.PermCodeAdminPaymentsExport,

		// Withdraw Management
		model.PermCodeAdminWithdrawsRead,
		model.PermCodeAdminWithdrawsApprove,
		model.PermCodeAdminWithdrawsReject,
		model.PermCodeAdminWithdrawsExport,

		// Audit Log (read only for admin)
		model.PermCodeAdminAuditRead,

		// Stats
		model.PermCodeAdminStatsRead,
		model.PermCodeAdminStatsExport,

		// Review Management
		model.PermCodeReviewList,
		model.PermCodeReviewGet,
		model.PermCodeReviewPending,
		model.PermCodeReviewLogs,
		model.PermCodeReviewPlayer,
		model.PermCodeReviewOrder,
		model.PermCodeReviewApprove,
		model.PermCodeReviewReject,
		model.PermCodeReviewBatchApprove,
		model.PermCodeReviewBatchReject,
		model.PermCodeReviewDelete,
		model.PermCodeReviewUpdate,
		model.PermCodeReviewStats,
		model.PermCodeReviewTrend,
		model.PermCodeReviewTopPlayers,
		model.PermCodeReviewGameStats,
		model.PermCodeReviewExport,
		model.PermCodeReviewDetectSensitive,

		// Review Report Management
		model.PermCodeReviewReportList,
		model.PermCodeReviewReportGet,
		model.PermCodeReviewReportHandle,

		// Review Reply Management
		model.PermCodeReviewReplyUpdate,
		model.PermCodeReviewReplyDelete,

		// Review Settings
		model.PermCodeReviewSettingsGet,
		model.PermCodeReviewSettingsUpdate,

		// Sensitive Word Management
		model.PermCodeSensitiveWordList,
		model.PermCodeSensitiveWordCreate,
		model.PermCodeSensitiveWordUpdate,
		model.PermCodeSensitiveWordDelete,

		// Content Management
		model.PermCodeContentFeedList,
		model.PermCodeContentFeedGet,
		model.PermCodeContentFeedApprove,
		model.PermCodeContentFeedReject,
		model.PermCodeContentFeedBatchApprove,
		model.PermCodeContentFeedBatchReject,
		model.PermCodeContentFeedDelete,

		// Chat Monitoring
		model.PermCodeContentChatList,
		model.PermCodeContentChatDelete,
		model.PermCodeContentChatMute,
		model.PermCodeContentChatUnmute,

		// Content Report Management
		model.PermCodeContentReportList,
		model.PermCodeContentReportGet,
		model.PermCodeContentReportProcess,

		// Content Stats
		model.PermCodeContentStats,

		// Content Category Management
		model.PermCodeContentCategoryList,
		model.PermCodeContentCategoryGet,
		model.PermCodeContentCategoryCreate,
		model.PermCodeContentCategoryUpdate,
		model.PermCodeContentCategoryDelete,

		// Commission Management
		model.PermCodeCommissionRead,
		model.PermCodeCommissionSettle,
		model.PermCodeCommissionExport,

		// Notification Management
		model.PermCodeNotificationRead,
		model.PermCodeNotificationCreate,
		model.PermCodeNotificationBatchSend,

		// Wallet Management
		model.PermCodeWalletRead,
		model.PermCodeWalletUpdate,
		model.PermCodeWalletTransactions,

		// Service Item Management
		model.PermCodeServiceItemRead,
		model.PermCodeServiceItemCreate,
		model.PermCodeServiceItemUpdate,
		model.PermCodeServiceItemDelete,

		// Dispute Management
		model.PermCodeDisputeRead,
		model.PermCodeDisputeCreate,
		model.PermCodeDisputeResolve,

		// Operation Log
		model.PermCodeOperationLogList,
	},
	model.RoleSlugPlayer: {
		// Players can view their own reviews
		model.PermCodeReviewList,
		model.PermCodeReviewGet,
		model.PermCodeReviewPlayer,
		model.PermCodeReviewOrder,
		model.PermCodeReviewStats,

		// Players can view service items
		model.PermCodeServiceItemRead,

		// Players can view their wallet
		model.PermCodeWalletRead,
		model.PermCodeWalletTransactions,

		// Players can view their commission
		model.PermCodeCommissionRead,

		// Players can view notifications
		model.PermCodeNotificationRead,

		// Players can view content categories
		model.PermCodeContentCategoryList,
		model.PermCodeContentCategoryGet,

		// Players can view their own disputes
		model.PermCodeDisputeRead,
	},
	model.RoleSlugUser: {
		// Users can view public reviews
		model.PermCodeReviewList,
		model.PermCodeReviewGet,
		model.PermCodeReviewOrder,

		// Users can create review reports
		model.PermCodeReviewReportCreate,

		// Users can view service items
		model.PermCodeServiceItemRead,

		// Users can view their wallet
		model.PermCodeWalletRead,
		model.PermCodeWalletTransactions,

		// Users can view notifications
		model.PermCodeNotificationRead,

		// Users can view content categories
		model.PermCodeContentCategoryList,
		model.PermCodeContentCategoryGet,

		// Users can view and create disputes
		model.PermCodeDisputeRead,
		model.PermCodeDisputeCreate,
	},
}

// SystemRoleDefinitions defines the system roles with their metadata.
var SystemRoleDefinitions = []model.RoleModel{
	{
		Slug:        string(model.RoleSlugSuperAdmin),
		Name:        "超级管理员",
		Description: "拥有系统所有权限的超级管理员角色",
		IsSystem:    true,
		Priority:    1000,
		Level:       0,
	},
	{
		Slug:        string(model.RoleSlugAdmin),
		Name:        "管理员",
		Description: "平台管理员/店长，负责日常运营管理",
		IsSystem:    true,
		Priority:    500,
		Level:       0,
	},
	{
		Slug:        string(model.RoleSlugFinance),
		Name:        "财务",
		Description: "财务人员，负责提现审批、支付退款等财务操作",
		IsSystem:    true,
		Priority:    400,
		Level:       0,
	},
	{
		Slug:        string(model.RoleSlugCustomerService),
		Name:        "客服",
		Description: "客服人员，负责纠纷处理、订单协助、用户服务等",
		IsSystem:    true,
		Priority:    300,
		Level:       0,
	},
	{
		Slug:        string(model.RoleSlugPlayer),
		Name:        "陪玩师",
		Description: "平台陪玩师角色，提供陪玩服务",
		IsSystem:    true,
		Priority:    100,
		Level:       0,
	},
	{
		Slug:        string(model.RoleSlugUser),
		Name:        "普通用户",
		Description: "平台普通用户，可以下单和使用服务",
		IsSystem:    true,
		Priority:    10,
		Level:       0,
	},
}

// seedDefaultRoles creates the default system roles and assigns their permissions.
func seedDefaultRoles(tx *gorm.DB) error {
	// 1. Create or update system roles
	roles := make(map[string]*model.RoleModel)
	for _, roleDef := range SystemRoleDefinitions {
		role := roleDef // Create a copy to avoid pointer issues
		var existing model.RoleModel
		if err := tx.Where("slug = ?", role.Slug).First(&existing).Error; err == nil {
			// Role exists, update it
			existing.Name = role.Name
			existing.Description = role.Description
			existing.IsSystem = role.IsSystem
			existing.Priority = role.Priority
			if err := tx.Save(&existing).Error; err != nil {
				return err
			}
			roles[role.Slug] = &existing
		} else if err == gorm.ErrRecordNotFound {
			// Role doesn't exist, create it
			if err := tx.Create(&role).Error; err != nil {
				return err
			}
			roles[role.Slug] = &role
		} else {
			return err
		}
	}

	// 2. Assign permissions to each role
	for roleSlug, permCodes := range DefaultRolePermissions {
		role, ok := roles[string(roleSlug)]
		if !ok {
			log.Printf("Warning: role %s not found, skipping permission assignment\n", roleSlug)
			continue
		}

		// Skip superAdmin - they get all permissions via "*" wildcard
		if roleSlug == model.RoleSlugSuperAdmin {
			log.Printf("SuperAdmin role uses wildcard permission, skipping explicit assignment\n")
			continue
		}

		// Get permission IDs for the codes
		var permissionIDs []uint64
		for _, code := range permCodes {
			var perm model.Permission
			if err := tx.Where("code = ?", code).First(&perm).Error; err == nil {
				permissionIDs = append(permissionIDs, perm.ID)
			} else if err != gorm.ErrRecordNotFound {
				return err
			} else {
				log.Printf("Warning: permission code %s not found, skipping\n", code)
			}
		}

		if len(permissionIDs) == 0 {
			log.Printf("No permissions found for role %s\n", roleSlug)
			continue
		}

		// Clear existing permissions and assign new ones
		if err := tx.Exec("DELETE FROM role_permissions WHERE role_id = ?", role.ID).Error; err != nil {
			return err
		}

		// Insert new permissions
		for _, permID := range permissionIDs {
			rp := model.RolePermission{
				RoleID:       role.ID,
				PermissionID: permID,
			}
			if err := tx.Create(&rp).Error; err != nil {
				// Skip duplicates
				if !isDuplicateKeyError(err) {
					return err
				}
			}
		}

		log.Printf("Assigned %d permissions to role %s\n", len(permissionIDs), roleSlug)
	}

	log.Println("Default roles and permissions seeded successfully")
	return nil
}

// isDuplicateKeyError checks if the error is a duplicate key error.
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return contains(errStr, "UNIQUE constraint failed") ||
		contains(errStr, "duplicate key") ||
		contains(errStr, "Duplicate entry")
}

// contains checks if s contains substr (simple implementation to avoid importing strings).
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
