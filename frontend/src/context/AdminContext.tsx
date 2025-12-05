/**
 * 管理员上下文
 * 提供菜单、权限数据和权限检查方法
 */
import React, { createContext, useContext, useState, useEffect, useCallback, useMemo } from 'react';
import { adminApi } from '@/api/admin';
import type { Menu } from '@/api/admin';
import { permissionStore } from '@/utils/permission';

/**
 * 管理员上下文类型接口
 */
interface AdminContextType {
    /** 可访问的菜单列表 */
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
}

const AdminContext = createContext<AdminContextType>({
    menus: [],
    permissions: [],
    loading: false,
    refreshMenus: async () => { },
    hasPermission: () => false,
    hasAllPermissions: () => false,
    hasAnyPermission: () => false,
    isSuperAdmin: false,
});

/**
 * 获取管理员上下文Hook
 *
 * @example
 * ```tsx
 * const { permissions, hasPermission, isSuperAdmin } = useAdmin();
 *
 * if (hasPermission('admin.games.create')) {
 *     // 显示创建按钮
 * }
 * ```
 */
export const useAdmin = () => useContext(AdminContext);

/**
 * 管理员上下文提供者
 *
 * @description
 * 管理用户的菜单和权限数据，提供权限检查方法
 */
export const AdminProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [menus, setMenus] = useState<Menu[]>([]);
    const [permissions, setPermissions] = useState<string[]>([]);
    const [loading, setLoading] = useState(false);

    /**
     * 刷新菜单和权限数据
     */
    const refreshMenus = useCallback(async () => {
        const token = localStorage.getItem('token');
        if (!token) {
            setPermissions([]);
            setMenus([]);
            permissionStore.clearPermissions();
            return;
        }

        setLoading(true);
        try {
            // 并行获取权限和当前用户可访问的菜单
            const [permRes, menuRes] = await Promise.all([
                adminApi.getMyPermissions().catch(() => ({ data: [] })),
                adminApi.getMyMenus().catch(() => ({ data: [] }))
            ]);

            console.log('[AdminContext] permRes:', permRes);
            console.log('[AdminContext] menuRes:', menuRes);

            // 从响应中提取数据
            // 注意: apiClient 响应拦截器已返回 response.data，所以这里直接是 ApiResponse
            // permRes 格式: { success, code, message, data: string[] }
            // menuRes 格式: { success, code, message, data: Menu[] }
            const apiPermRes = permRes as unknown as { data?: string[] };
            const apiMenuRes = menuRes as unknown as { data?: Menu[] };
            
            const permData = apiPermRes.data || [];
            const menuData = apiMenuRes.data || [];

            // 更新状态
            setPermissions(permData);
            setMenus(menuData);

            // 同步到权限存储（供非React环境使用）
            permissionStore.setPermissions(permData);
        } catch (error) {
            console.error('Failed to fetch admin info', error);
            setPermissions([]);
            setMenus([]);
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

    // 初始化时加载权限
    useEffect(() => {
        refreshMenus();
    }, [refreshMenus]);

    // 监听登录/登出事件
    useEffect(() => {
        const handleStorageChange = (e: StorageEvent) => {
            if (e.key === 'token') {
                refreshMenus();
            }
        };

        window.addEventListener('storage', handleStorageChange);
        return () => window.removeEventListener('storage', handleStorageChange);
    }, [refreshMenus]);

    const contextValue = useMemo(
        () => ({
            menus,
            permissions,
            loading,
            refreshMenus,
            hasPermission,
            hasAllPermissions,
            hasAnyPermission,
            isSuperAdmin,
        }),
        [menus, permissions, loading, refreshMenus, hasPermission, hasAllPermissions, hasAnyPermission, isSuperAdmin]
    );

    return (
        <AdminContext.Provider value={contextValue}>
            {children}
        </AdminContext.Provider>
    );
};

export default AdminContext;
