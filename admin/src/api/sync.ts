/**
 * 路由和权限同步API
 * 用于在前端启动时自动同步菜单和权限到后端
 */
import apiClient from './client';
import type { ApiResponse } from './admin';
import type { MenuConfig, PermissionConfig } from '@/config/adminRoutes';

import { logger } from '@/utils/logger';
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
 * 批量同步请求
 */
export interface BatchSyncRequest {
    menus: MenuSyncItem[];
    permissions: PermissionSyncItem[];
    assignSuperAdminPermissions: boolean;
}

/**
 * 菜单同步项
 */
export interface MenuSyncItem {
    name: string;
    path: string;
    component: string;
    icon?: string;
    order: number;
    hidden?: boolean;
    visible?: boolean;
    permission?: string;
    redirect?: string;
    description?: string;
    children?: MenuSyncItem[];
}

/**
 * 权限同步项
 */
export interface PermissionSyncItem {
    code: string;
    name: string;
    method: string;
    path: string;
    group: string;
    description: string;
}

/**
 * 批量同步响应
 */
export interface BatchSyncResponse {
    success: boolean;
    menuSync?: SyncResult;
    permissionSync?: SyncResult;
    superAdminAssign?: { success: boolean; message: string };
    errors: string[];
}

/**
 * 同步API
 */
export const syncApi = {
    /**
     * 批量同步（一次性同步菜单、权限并分配超管权限）
     */
    batchSync: async (
        menus: MenuConfig[],
        permissions: PermissionConfig[],
        assignSuperAdmin: boolean = true
    ): Promise<BatchSyncResponse> => {
        // 转换菜单配置为同步格式
        const convertMenus = (items: MenuConfig[]): MenuSyncItem[] => {
            return items.map(item => ({
                name: item.name,
                path: item.path,
                component: item.component,
                icon: item.icon,
                order: item.order,
                hidden: item.hidden,
                visible: !item.hidden,
                permission: item.permission,
                redirect: item.redirect,
                description: item.description,
                children: item.children ? convertMenus(item.children) : undefined,
            }));
        };

        // 转换权限配置为同步格式
        const convertPermissions = (items: PermissionConfig[]): PermissionSyncItem[] => {
            return items.map(item => ({
                code: item.code,
                name: item.description,
                method: item.method,
                path: item.path,
                group: item.group,
                description: item.description,
            }));
        };

        const request: BatchSyncRequest = {
            menus: convertMenus(menus),
            permissions: convertPermissions(permissions),
            assignSuperAdminPermissions: assignSuperAdmin,
        };

        try {
            const response = (await apiClient.post<ApiResponse<BatchSyncResponse>>(
                '/admin/sync/batch',
                request
            )) as unknown as ApiResponse<BatchSyncResponse>;
            logger.info('[Sync] batchSync response:', response);
            return response.data;
        } catch (error) {
            return {
                success: false,
                errors: [`批量同步失败: ${error instanceof Error ? error.message : '未知错误'}`],
            };
        }
    },
};

export default syncApi;
