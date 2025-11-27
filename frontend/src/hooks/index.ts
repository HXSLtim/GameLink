/**
 * Hooks 统一导出
 */

// 权限相关 Hooks
export {
    usePermission,
    usePermissions,
    useHasPermission,
    usePermissionChecker,
    type PermissionCheckResult,
} from './usePermission';

// 同步相关 Hooks
export {
    useSync,
    type SyncState,
    type UseSyncReturn,
} from './useSync';
