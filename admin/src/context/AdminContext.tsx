/**
 * 管理员上下文
 * 提供菜单、权限数据和权限检查方法
 * Requirements: 8.1, 8.2, 8.4 - 菜单权限联动
 */
import React, { createContext, useState, useEffect, useCallback, useMemo } from 'react';
import { adminApi } from '@/api/admin';
import type { Menu } from '@/api/admin';
import { permissionStore } from '@/utils/permission';
import { filterMenusByPermission } from '@/utils/menuPermission';

import { PERMISSION_CHANGE_EVENT, triggerPermissionChange } from './permissionEvents';

/**
 * 管理员上下文类型接口
 */
interface AdminContextType {
    /** 原始菜单列表（未过滤） */
    rawMenus: Menu[];
    /** 根据权限过滤后的菜单列表 */
    menus: Menu[];
    /** 权限码数组 */
    permissions: string[];
    /** 是否正在加载 */
    loading: boolean;
    /** 刷新菜单和权限 */
    refreshMenus: () => Promise<void>;
    /** 检查是否拥有指定权限 */
    hasPermission: (permission: string | string[], mode?: 'any' | 'all') => boolean;
    /** 检查是否拥有全部指定权限 */
    hasAllPermissions: (permissions: string[]) => boolean;
    /** 检查是否拥有任一指定权限 */
    hasAnyPermission: (permissions: string[]) => boolean;
    /** 是否为超级管理员 */
    isSuperAdmin: boolean;
    /** 权限版本号，用于触发权限变更后的更新 */
    permissionVersion: number;
    /** 触发权限变更通知（通知其他组件和标签页） */
    notifyPermissionChange: () => void;
}

const AdminContext = createContext<AdminContextType>({
    rawMenus: [],
    menus: [],
    permissions: [],
    loading: false,
    refreshMenus: async () => { },
    hasPermission: () => false,
    hasAllPermissions: () => false,
    hasAnyPermission: () => false,
    isSuperAdmin: false,
    permissionVersion: 0,
    notifyPermissionChange: () => { },
});

// Note: useAdmin is available from './useAdmin' for Fast Refresh compatibility

/**
 * 管理员上下文提供者
 *
 * @description
 * 管理用户的菜单和权限数据，提供权限检查方法
 */
export const AdminProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [rawMenus, setRawMenus] = useState<Menu[]>([]);
    const [permissions, setPermissions] = useState<string[]>([]);
    const [loading, setLoading] = useState(false);
    const [permissionVersion, setPermissionVersion] = useState(0);

    /**
     * 刷新菜单和权限数据
     * Requirements: 8.4 - 权限变更后菜单更新
     */
    const refreshMenus = useCallback(async () => {
        const token = localStorage.getItem('token');
        if (!token) {
            setPermissions([]);
            setRawMenus([]);
            permissionStore.clearPermissions();
            return;
        }

        setLoading(true);
        try {
            // 并行获取权限和当前用户可访问的菜单
            const [permRes, menuRes] = await Promise.all([
                adminApi.getMyPermissions().catch(() => ({ data: { success: false, data: [] } })),
                adminApi.getMyMenus().catch(() => ({ data: { success: false, data: [] } }))
            ]);

            console.log('[AdminContext] permRes:', permRes);
            console.log('[AdminContext] menuRes:', menuRes);

            // 从响应中提取数据
            // apiClient 返回完整的 AxiosResponse，所以需要访问 response.data
            // permRes.data 格式: { success, code, message, data: string[] }
            // menuRes.data 格式: { success, code, message, data: Menu[] }
            const apiPermRes = permRes as unknown as { data?: { success?: boolean; data?: string[] } };
            const apiMenuRes = menuRes as unknown as { data?: { success?: boolean; data?: Menu[] } };
            
            const permData = apiPermRes.data?.data || [];
            const menuData = apiMenuRes.data?.data || [];

            // 更新状态
            setPermissions(permData);
            setRawMenus(menuData);

            // 同步到权限存储（供非React环境使用）
            permissionStore.setPermissions(permData);
            
            // 增加权限版本号，触发依赖组件更新
            setPermissionVersion(v => v + 1);
        } catch (error) {
            console.error('Failed to fetch admin info', error);
            setPermissions([]);
            setRawMenus([]);
            permissionStore.clearPermissions();
        } finally {
            setLoading(false);
        }
    }, []);

    /**
     * 检查是否拥有指定权限
     */
    const hasPermission = useCallback(
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

    /**
     * 检查是否拥有全部指定权限
     */
    const hasAllPermissions = useCallback(
        (perms: string[]): boolean => hasPermission(perms, 'all'),
        [hasPermission]
    );

    /**
     * 检查是否拥有任一指定权限
     */
    const hasAnyPermission = useCallback(
        (perms: string[]): boolean => hasPermission(perms, 'any'),
        [hasPermission]
    );

    /**
     * 是否为超级管理员
     */
    const isSuperAdmin = useMemo(
        () => permissions.includes('*'),
        [permissions]
    );

    /**
     * 根据权限过滤后的菜单
     * Requirements: 8.1, 8.2 - 菜单权限过滤
     */
    const menus = useMemo(
        () => filterMenusByPermission(rawMenus, permissions),
        [rawMenus, permissions]
    );

    /**
     * 触发权限变更通知
     * 用于在权限分配后通知其他组件和标签页刷新
     * Requirements: 8.4 - 权限变更后菜单更新
     */
    const notifyPermissionChange = useCallback(() => {
        triggerPermissionChange();
        // 同时刷新当前上下文
        refreshMenus();
    }, [refreshMenus]);

    // 初始化时加载权限
    useEffect(() => {
        refreshMenus();
    }, [refreshMenus]);

    // 监听登录/登出事件和权限变更事件
    // Requirements: 8.4 - 权限变更后菜单更新
    useEffect(() => {
        const handleStorageChange = (e: StorageEvent) => {
            if (e.key === 'token') {
                refreshMenus();
            }
            // 监听其他标签页的权限变更
            if (e.key === 'permission_change_timestamp') {
                console.log('[AdminContext] 检测到其他标签页权限变更，刷新权限...');
                refreshMenus();
            }
        };

        // 监听当前标签页的权限变更事件
        const handlePermissionChange = () => {
            console.log('[AdminContext] 检测到权限变更事件，刷新权限...');
            refreshMenus();
        };

        window.addEventListener('storage', handleStorageChange);
        window.addEventListener(PERMISSION_CHANGE_EVENT, handlePermissionChange);
        
        return () => {
            window.removeEventListener('storage', handleStorageChange);
            window.removeEventListener(PERMISSION_CHANGE_EVENT, handlePermissionChange);
        };
    }, [refreshMenus]);

    const contextValue = useMemo(
        () => ({
            rawMenus,
            menus,
            permissions,
            loading,
            refreshMenus,
            hasPermission,
            hasAllPermissions,
            hasAnyPermission,
            isSuperAdmin,
            permissionVersion,
            notifyPermissionChange,
        }),
        [rawMenus, menus, permissions, loading, refreshMenus, hasPermission, hasAllPermissions, hasAnyPermission, isSuperAdmin, permissionVersion, notifyPermissionChange]
    );

    return (
        <AdminContext.Provider value={contextValue}>
            {children}
        </AdminContext.Provider>
    );
};

export default AdminContext;
