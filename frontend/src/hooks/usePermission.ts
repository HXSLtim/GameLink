/**
 * 权限检查Hook
 * 提供按钮级权限控制功能
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
 * usePermission Hook
 * 用于检查当前用户是否拥有指定权限
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

    const hasPermission = useMemo(() => {
        // 如果没有指定权限要求，默认允许
        const permissionList = Array.isArray(permission) ? permission : [permission];
        const hasValidPermission = permissionList.some(p => p && p.length > 0);
        if (!hasValidPermission) {
            return true;
        }

        if (loading || !permissions.length) {
            return false;
        }

        // 超级管理员拥有所有权限
        if (permissions.includes('*')) {
            return true;
        }

        if (mode === 'all') {
            // 全部满足模式
            return permissionList.every(p => permissions.includes(p));
        } else {
            // 任一满足模式
            return permissionList.some(p => permissions.includes(p));
        }
    }, [permissions, permission, mode, loading]);

    return { hasPermission, loading };
}

/**
 * usePermissions Hook
 * 批量检查多个权限，返回每个权限的检查结果
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
 * ```
 */
export function usePermissions<T extends Record<string, string>>(
    permissionMap: T
): Record<keyof T, boolean> & { loading: boolean } {
    const { permissions, loading } = useAdmin();

    const result = useMemo(() => {
        const checkResult = {} as Record<keyof T, boolean>;

        // 超级管理员拥有所有权限
        const isSuperAdmin = permissions.includes('*');

        for (const key in permissionMap) {
            if (Object.prototype.hasOwnProperty.call(permissionMap, key)) {
                checkResult[key] = isSuperAdmin || permissions.includes(permissionMap[key]);
            }
        }

        return { ...checkResult, loading };
    }, [permissions, permissionMap, loading]);

    return result;
}

/**
 * useHasPermission Hook
 * 简化版本，直接返回boolean
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
 * ```
 */
export function usePermissionChecker(): (permission: string | string[], mode?: 'any' | 'all') => boolean {
    const { permissions } = useAdmin();

    return useCallback(
        (permission: string | string[], mode: 'any' | 'all' = 'any'): boolean => {
            if (!permissions.length) {
                return false;
            }

            // 超级管理员拥有所有权限
            if (permissions.includes('*')) {
                return true;
            }

            const permissionList = Array.isArray(permission) ? permission : [permission];

            if (mode === 'all') {
                return permissionList.every(p => permissions.includes(p));
            } else {
                return permissionList.some(p => permissions.includes(p));
            }
        },
        [permissions]
    );
}

export default usePermission;
