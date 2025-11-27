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
     * 获取所有权限列表
     */
    getPermissions: () =>
        apiClient.get<ApiResponse<Permission[]>>('/admin/permissions'),

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
            const existingPermissions = response.data || [];
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
            }
        } catch (error) {
            result.success = false;
            result.errors.push(`获取权限列表失败: ${error instanceof Error ? error.message : '未知错误'}`);
        }

        return result;
    },

    /**
     * 获取所有角色列表
     */
    getRoles: () =>
        apiClient.get<ApiResponse<Role[]>>('/admin/roles'),

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
            const roles = rolesResponse.data || [];

            // 查找超级管理员角色
            const superAdminRole = roles.find(
                r => r.code === 'super_admin' || r.name === '超级管理员'
            );

            if (!superAdminRole) {
                return { success: false, message: '未找到超级管理员角色' };
            }

            // 获取所有权限
            const permissionsResponse = await syncApi.getPermissions();
            const permissions = permissionsResponse.data || [];

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
            const existingMenus = response.data || [];
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
