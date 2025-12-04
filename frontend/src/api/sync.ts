/**
 * 路由和权限同步API
 * 用于在前端启动时自动同步菜单和权限到后端
 */
import apiClient from './client';
import type { ApiResponse, Menu, CreateMenuDto } from './admin';
import type { MenuConfig, PermissionConfig } from '@/config/adminRoutes';

/**
 * 权限接口
 */
export interface Permission {
    id: number;
    code: string;
    name: string;
    method: string;
    path: string;
    group: string;
    description: string;
    createdAt: string;
    updatedAt: string;
}

/**
 * 创建权限DTO
 */
export interface CreatePermissionDto {
    code: string;
    name: string;
    method: string;
    path: string;
    group: string;
    description: string;
}

/**
 * 角色接口
 */
export interface Role {
    id: number;
    name: string;
    code: string;
    description: string;
    permissions: Permission[];
    createdAt: string;
    updatedAt: string;
}

/**
 * 同步结果
 */
export interface SyncResult {
    success: boolean;
    created: number;
    updated: number;
    skipped: number;
    errors: string[];
}

/**
 * 批量同步菜单请求
 */
export interface BatchSyncMenusRequest {
    menus: CreateMenuDto[];
    /** 是否清除未包含的菜单 */
    cleanOrphans?: boolean;
}

/**
 * 批量同步权限请求
 */
export interface BatchSyncPermissionsRequest {
    permissions: CreatePermissionDto[];
    /** 是否清除未包含的权限 */
    cleanOrphans?: boolean;
}

/**
 * 同步API
 */
export const syncApi = {
    /**
     * 获取所有权限列表（同步专用接口，不受限流限制）
     */
    getPermissions: () =>
        apiClient.get<ApiResponse<Permission[]>>('/admin/sync/permissions'),

    /**
     * 创建权限
     */
    createPermission: (data: CreatePermissionDto) =>
        apiClient.post<ApiResponse<Permission>>('/admin/permissions', data),

    /**
     * 更新权限
     */
    updatePermission: (id: number, data: Partial<CreatePermissionDto>) =>
        apiClient.put<ApiResponse<Permission>>(`/admin/permissions/${id}`, data),

    /**
     * 批量同步权限（创建或更新）
     */
    syncPermissions: async (permissions: PermissionConfig[]): Promise<SyncResult> => {
        const result: SyncResult = {
            success: true,
            created: 0,
            updated: 0,
            skipped: 0,
            errors: [],
        };

        try {
            // 获取现有权限
            const response = await syncApi.getPermissions();
            console.log('[Sync] getPermissions response:', response);
            console.log('[Sync] response.data:', response.data);
            console.log('[Sync] response.data?.items:', response.data?.items);

            // 后端返回的分页结构: { items: [...], page: 1, pageSize: 10, totalCount: 50 }
            let existingPermissions: Permission[] = [];
            if (response.data && typeof response.data === 'object' && 'items' in response.data) {
                existingPermissions = (response.data as any).items || [];
            } else if (Array.isArray(response.data)) {
                existingPermissions = response.data;
            } else {
                console.error('[Sync] Unexpected data structure:', response.data);
                existingPermissions = [];
            }

            console.log('[Sync] existingPermissions:', existingPermissions);
            const existingMap = new Map(existingPermissions.map(p => [p.code, p]));

            // 逐个处理权限
            for (const perm of permissions) {
                try {
                    const existing = existingMap.get(perm.code);
                    const permData: CreatePermissionDto = {
                        code: perm.code,
                        name: perm.description,
                        method: perm.method,
                        path: perm.path,
                        group: perm.group,
                        description: perm.description,
                    };

                    if (existing) {
                        // 检查是否需要更新
                        if (
                            existing.method !== perm.method ||
                            existing.path !== perm.path ||
                            existing.description !== perm.description
                        ) {
                            await syncApi.updatePermission(existing.id, permData);
                            result.updated++;
                        } else {
                            result.skipped++;
                        }
                    } else {
                        await syncApi.createPermission(permData);
                        result.created++;
                    }
                } catch (error) {
                    result.errors.push(`权限 ${perm.code}: ${error instanceof Error ? error.message : '未知错误'}`);
                }

                // 添加延迟避免限流
                await new Promise(resolve => setTimeout(resolve, 100));
            }
        } catch (error) {
            result.success = false;
            result.errors.push(`获取权限列表失败: ${error instanceof Error ? error.message : '未知错误'}`);
        }

        return result;
    },

    /**
     * 获取所有角色列表（同步专用接口，不受限流限制）
     */
    getRoles: () =>
        apiClient.get<ApiResponse<Role[]>>('/admin/sync/roles'),

    /**
     * 获取角色详情
     */
    getRole: (id: number) =>
        apiClient.get<ApiResponse<Role>>(`/admin/roles/${id}`),

    /**
     * 分配角色权限
     */
    assignRolePermissions: (roleId: number, permissionIds: number[]) =>
        apiClient.put<ApiResponse<void>>(`/admin/roles/${roleId}/permissions`, { permissionIds }),

    /**
     * 为超级管理员分配所有权限
     */
    assignAllPermissionsToSuperAdmin: async (): Promise<{ success: boolean; message: string }> => {
        try {
            // 获取所有角色
            const rolesResponse = await syncApi.getRoles();
            console.log('[Sync] getRoles response:', rolesResponse);

            // 后端返回的分页结构: { items: [...], page: 1, pageSize: 10, totalCount: 50 }
            const roles: Role[] = rolesResponse.data?.items || [];
            console.log('[Sync] roles:', roles);

            // 查找超级管理员角色
            const superAdminRole = roles.find(
                r => r.code === 'super_admin' || r.name === '超级管理员'
            );

            if (!superAdminRole) {
                return { success: false, message: '未找到超级管理员角色' };
            }

            // 获取所有权限
            const permissionsResponse = await syncApi.getPermissions();
            console.log('[Sync] getPermissions response:', permissionsResponse);

            // 后端返回的分页结构: { items: [...], page: 1, pageSize: 10, totalCount: 50 }
            const permissions: Permission[] = permissionsResponse.data?.items || [];
            console.log('[Sync] permissions:', permissions);

            if (permissions.length === 0) {
                return { success: false, message: '没有可分配的权限' };
            }

            // 分配所有权限给超级管理员
            const permissionIds = permissions.map(p => p.id);
            await syncApi.assignRolePermissions(superAdminRole.id, permissionIds);

            return {
                success: true,
                message: `成功为超级管理员分配 ${permissions.length} 个权限`
            };
        } catch (error) {
            return {
                success: false,
                message: `分配权限失败: ${error instanceof Error ? error.message : '未知错误'}`
            };
        }
    },

    /**
     * 同步菜单（递归创建）
     */
    syncMenus: async (menus: MenuConfig[], parentId: number | null = null): Promise<SyncResult> => {
        const result: SyncResult = {
            success: true,
            created: 0,
            updated: 0,
            skipped: 0,
            errors: [],
        };

        try {
            // 获取现有菜单
            const response = await apiClient.get<ApiResponse<Menu[]>>('/admin/menus', {
                params: { parentId }
            });
            console.log('[Sync] getMenus response:', response);

            // 菜单接口在无分页参数时直接返回数组，有分页时返回 { data: [...], pagination: {...} }
            let existingMenus: Menu[] = [];
            if (Array.isArray(response.data)) {
                existingMenus = response.data;
            } else if (response.data && typeof response.data === 'object' && 'data' in response.data) {
                // 有分页的情况
                existingMenus = (response.data as any).data || [];
            } else {
                existingMenus = response.data || [];
            }

            console.log('[Sync] existingMenus:', existingMenus);
            const existingMap = new Map(existingMenus.map(m => [m.path, m]));

            for (const menu of menus) {
                try {
                    const existing = existingMap.get(menu.path);
                    const menuData: CreateMenuDto = {
                        name: menu.name,
                        path: menu.path,
                        component: menu.component,
                        parentId: parentId,
                        order: menu.order,
                        hidden: menu.hidden || false,
                        permission: menu.permission,
                        icon: menu.icon,
                        redirect: menu.redirect,
                        description: menu.description,
                    };

                    let currentMenuId: number;

                    if (existing) {
                        // 检查是否需要更新
                        if (
                            existing.name !== menu.name ||
                            existing.component !== menu.component ||
                            existing.order !== menu.order ||
                            existing.icon !== menu.icon
                        ) {
                            await apiClient.put(`/admin/menus/${existing.id}`, menuData);
                            result.updated++;
                        } else {
                            result.skipped++;
                        }
                        currentMenuId = existing.id;
                    } else {
                        const createResponse = await apiClient.post<ApiResponse<Menu>>('/admin/menus', menuData);
                        result.created++;
                        currentMenuId = createResponse.data.id;
                    }

                    // 递归处理子菜单
                    if (menu.children && menu.children.length > 0) {
                        const childResult = await syncApi.syncMenus(menu.children, currentMenuId);
                        result.created += childResult.created;
                        result.updated += childResult.updated;
                        result.skipped += childResult.skipped;
                        result.errors.push(...childResult.errors);
                    }
                } catch (error) {
                    result.errors.push(`菜单 ${menu.path}: ${error instanceof Error ? error.message : '未知错误'}`);
                }
            }
        } catch (error) {
            result.success = false;
            result.errors.push(`获取菜单列表失败: ${error instanceof Error ? error.message : '未知错误'}`);
        }

        return result;
    },
};

export default syncApi;
