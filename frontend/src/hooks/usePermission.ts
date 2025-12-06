/**
 * 权限检查Hook
 * 提供按钮级权限控制功能
 * Requirements: 3.1, 3.4
 */
import { useMemo, useCallback } from 'react';
import { useAdmin } from '@/context/AdminContext';

/**
 * 权限检查结果接口
 */
export interface PermissionCheckResult {
    /** 是否拥有权限 */
    hasPermission: boolean;
    /** 权限检查是否正在加载 */
    loading: boolean;
}

/**
 * 批量权限检查结果类型
 */
export type BatchPermissionResult<T extends Record<string, string>> = {
    [K in keyof T]: boolean;
} & { loading: boolean };

/**
 * 权限检查函数类型
 */
export type PermissionChecker = (
    permission: string | string[],
    mode?: 'any' | 'all'
) => boolean;

/**
 * 将权限数组转换为稳定的字符串key，用于依赖比较
 * Requirements: 3.1 - 性能优化
 */
function getPermissionKey(permission: string | string[]): string {
    if (Array.isArray(permission)) {
        return permission.slice().sort().join(',');
    }
    return permission;
}


/**
 * 核心权限检查逻辑
 * 抽取为纯函数，便于复用和测试
 */
function checkPermission(
    userPermissions: string[],
    requiredPermission: string | string[],
    mode: 'any' | 'all' = 'any'
): boolean {
    const permissionList = Array.isArray(requiredPermission) 
        ? requiredPermission 
        : [requiredPermission];
    
    // 如果没有指定有效的权限要求，默认允许
    const hasValidPermission = permissionList.some(p => p && p.length > 0);
    if (!hasValidPermission) {
        return true;
    }

    if (!userPermissions.length) {
        return false;
    }

    // 超级管理员拥有所有权限
    if (userPermissions.includes('*')) {
        return true;
    }

    if (mode === 'all') {
        // 全部满足模式
        return permissionList.every(p => userPermissions.includes(p));
    } else {
        // 任一满足模式
        return permissionList.some(p => userPermissions.includes(p));
    }
}

/**
 * usePermission Hook
 * 用于检查当前用户是否拥有指定权限
 * 
 * 性能优化：
 * - 使用 useMemo 缓存计算结果
 * - 使用稳定的权限key避免不必要的重渲染
 * - 减少依赖项，只在权限真正变化时重新计算
 *
 * Requirements: 3.1
 *
 * @example
 * ```tsx
 * // 单个权限检查
 * const { hasPermission, loading } = usePermission('admin.games.create');
 *
 * // 多个权限检查（任一满足）
 * const { hasPermission } = usePermission(['admin.games.create', 'admin.games.update']);
 *
 * // 多个权限检查（全部满足）
 * const { hasPermission } = usePermission(['admin.games.create', 'admin.games.update'], 'all');
 * ```
 */
export function usePermission(
    permission: string | string[],
    mode: 'any' | 'all' = 'any'
): PermissionCheckResult {
    const { permissions, loading } = useAdmin();
    
    // 使用稳定的权限key作为依赖，避免数组引用变化导致的重渲染
    const permissionKey = useMemo(
        () => getPermissionKey(permission),
        [permission]
    );

    const hasPermission = useMemo(() => {
        if (loading) {
            return false;
        }
        return checkPermission(permissions, permission, mode);
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [permissions, permissionKey, mode, loading]);

    return useMemo(
        () => ({ hasPermission, loading }),
        [hasPermission, loading]
    );
}


/**
 * usePermissions Hook
 * 批量检查多个权限，返回每个权限的检查结果
 * 
 * 性能优化：
 * - 使用 JSON.stringify 生成稳定的 key 进行深度比较
 * - 单次遍历计算所有权限结果
 * - 使用 useMemo 缓存结果，避免不必要的重渲染
 *
 * Requirements: 3.4
 *
 * @example
 * ```tsx
 * const permissionMap = usePermissions({
 *     canCreate: 'admin.games.create',
 *     canEdit: 'admin.games.update',
 *     canDelete: 'admin.games.delete'
 * });
 *
 * if (permissionMap.canCreate) {
 *     // 显示创建按钮
 * }
 * 
 * // 检查加载状态
 * if (permissionMap.loading) {
 *     return <Spin />;
 * }
 * ```
 */
export function usePermissions<T extends Record<string, string>>(
    permissionMap: T
): BatchPermissionResult<T> {
    const { permissions, loading } = useAdmin();
    
    // 计算当前 permissionMap 的稳定 key，用于依赖比较
    const permissionMapKey = useMemo(() => {
        try {
            // 对 keys 排序以确保稳定性
            const sortedEntries = Object.entries(permissionMap).sort(([a], [b]) => a.localeCompare(b));
            return JSON.stringify(sortedEntries);
        } catch {
            return '';
        }
    }, [permissionMap]);

    const result = useMemo(() => {
        const checkResult = {} as Record<keyof T, boolean>;

        // 超级管理员拥有所有权限
        const isSuperAdmin = permissions.includes('*');

        for (const key in permissionMap) {
            if (Object.prototype.hasOwnProperty.call(permissionMap, key)) {
                if (loading) {
                    checkResult[key] = false;
                } else {
                    checkResult[key] = isSuperAdmin || permissions.includes(permissionMap[key]);
                }
            }
        }

        return { ...checkResult, loading } as BatchPermissionResult<T>;
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [permissions, permissionMapKey, loading]);

    return result;
}

/**
 * useHasPermission Hook
 * 简化版本，直接返回boolean
 * 
 * Requirements: 3.1
 *
 * @example
 * ```tsx
 * const canCreate = useHasPermission('admin.games.create');
 *
 * return canCreate && <Button>创建</Button>;
 * ```
 */
export function useHasPermission(permission: string): boolean {
    const { hasPermission } = usePermission(permission);
    return hasPermission;
}


/**
 * usePermissionChecker Hook
 * 返回一个检查函数，可以动态检查权限
 * 
 * 性能优化：
 * - 使用 useCallback 确保函数引用稳定
 * - 只依赖 permissions 数组，减少不必要的重新创建
 * 
 * Requirements: 3.4
 *
 * @example
 * ```tsx
 * const checkPermission = usePermissionChecker();
 *
 * const handleClick = () => {
 *     if (checkPermission('admin.games.delete')) {
 *         // 执行删除
 *     } else {
 *         message.error('没有删除权限');
 *     }
 * };
 * 
 * // 支持批量检查
 * const canOperate = checkPermission(['admin.games.update', 'admin.games.delete'], 'any');
 * ```
 */
export function usePermissionChecker(): PermissionChecker {
    const { permissions } = useAdmin();

    return useCallback(
        (permission: string | string[], mode: 'any' | 'all' = 'any'): boolean => {
            return checkPermission(permissions, permission, mode);
        },
        [permissions]
    );
}

/**
 * usePermissionCheckerWithLoading Hook
 * 返回检查函数和加载状态，适用于需要处理加载状态的场景
 * 
 * Requirements: 3.4
 *
 * @example
 * ```tsx
 * const { check, loading } = usePermissionCheckerWithLoading();
 *
 * if (loading) {
 *     return <Spin />;
 * }
 * 
 * const handleClick = () => {
 *     if (check('admin.games.delete')) {
 *         // 执行删除
 *     }
 * };
 * ```
 */
export function usePermissionCheckerWithLoading(): {
    check: PermissionChecker;
    loading: boolean;
} {
    const { permissions, loading } = useAdmin();

    const check = useCallback(
        (permission: string | string[], mode: 'any' | 'all' = 'any'): boolean => {
            if (loading) {
                return false;
            }
            return checkPermission(permissions, permission, mode);
        },
        [permissions, loading]
    );

    return useMemo(
        () => ({ check, loading }),
        [check, loading]
    );
}

/**
 * 导出核心检查函数，供非Hook场景使用（如工具函数）
 */
export { checkPermission };

export default usePermission;
