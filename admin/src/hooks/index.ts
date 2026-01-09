/**
 * Hooks 统一导出
 */

// 认证相关 Hooks (Super Dev 最佳实践)
export {
    useAuthCheck,
    useUserInfo,
    useIsAuthenticated,
    useIsHydrated,
    useAuthToken,
    useAuthLoading,
    useAuthError,
    useIsAdmin,
    type AuthCheckResult,
} from './useAuthCheck';

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

// 动画相关 Hooks
export {
    useCountUp,
    formatNumber,
} from './useCountUp';

// 主题感知消息 Hook
export { useAppMessage } from './useAppMessage';

// CRUD 操作 Hook
export {
    useCrud,
    type CrudId,
    type CrudQueryParams,
    type CrudPagination,
    type CrudMessages,
    type CrudApi,
    type UseCrudOptions,
    type UseCrudReturn,
} from './useCrud';

// 容器尺寸检测 Hook
export {
    useContainerReady,
    useResponsiveChartHeight,
} from './useContainerReady';

// localStorage Hook
export {
    useLocalStorage,
    useDashboardPreferences,
    type DashboardPreferences,
} from './useLocalStorage';

// 浏览器通知 Hook
export { useNotification } from './useNotification';
