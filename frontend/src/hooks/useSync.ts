/**
 * 同步功能Hook
 * 提供手动触发菜单和权限同步的功能
 */
import { useState, useCallback } from 'react';
import { message } from 'antd';
import {
    initApp,
    forceInit,
    syncPermissionsOnly,
    syncMenusOnly,
    assignSuperAdminOnly,
    type InitResult,
} from '@/services/init';
import type { SyncResult } from '@/api/sync';

/**
 * 同步状态
 */
export interface SyncState {
    loading: boolean;
    result: InitResult | SyncResult | null;
    error: string | null;
}

/**
 * 同步Hook返回值
 */
export interface UseSyncReturn {
    /** 同步状态 */
    state: SyncState;
    /** 完整同步（菜单+权限+超管权限） */
    syncAll: () => Promise<void>;
    /** 强制完整同步（忽略24小时缓存） */
    forceSyncAll: () => Promise<void>;
    /** 仅同步权限 */
    syncPermissions: () => Promise<void>;
    /** 仅同步菜单 */
    syncMenus: () => Promise<void>;
    /** 仅为超管分配权限 */
    assignSuperAdmin: () => Promise<void>;
    /** 重置状态 */
    reset: () => void;
}

/**
 * 同步功能Hook
 */
export const useSync = (): UseSyncReturn => {
    const [state, setState] = useState<SyncState>({
        loading: false,
        result: null,
        error: null,
    });

    /**
     * 完整同步
     */
    const syncAll = useCallback(async () => {
        setState({ loading: true, result: null, error: null });
        try {
            const result = await initApp({
                syncMenus: true,
                syncPermissions: true,
                assignSuperAdminPermissions: true,
                verbose: true,
            });
            setState({ loading: false, result, error: null });
            if (result.success) {
                message.success('同步成功');
            } else {
                message.warning(`同步完成，但有 ${result.errors.length} 个错误`);
            }
        } catch (error) {
            const errorMsg = error instanceof Error ? error.message : '未知错误';
            setState({ loading: false, result: null, error: errorMsg });
            message.error(`同步失败: ${errorMsg}`);
        }
    }, []);

    /**
     * 强制完整同步
     */
    const forceSyncAll = useCallback(async () => {
        setState({ loading: true, result: null, error: null });
        try {
            const result = await forceInit({
                syncMenus: true,
                syncPermissions: true,
                assignSuperAdminPermissions: true,
                verbose: true,
            });
            setState({ loading: false, result, error: null });
            if (result.success) {
                message.success('强制同步成功');
            } else {
                message.warning(`同步完成，但有 ${result.errors.length} 个错误`);
            }
        } catch (error) {
            const errorMsg = error instanceof Error ? error.message : '未知错误';
            setState({ loading: false, result: null, error: errorMsg });
            message.error(`强制同步失败: ${errorMsg}`);
        }
    }, []);

    /**
     * 仅同步权限
     */
    const syncPermissions = useCallback(async () => {
        setState({ loading: true, result: null, error: null });
        try {
            const result = await syncPermissionsOnly();
            setState({ loading: false, result, error: null });
            if (result.success) {
                message.success(`权限同步成功: 创建${result.created}个，更新${result.updated}个`);
            } else {
                message.warning(`同步完成，但有 ${result.errors.length} 个错误`);
            }
        } catch (error) {
            const errorMsg = error instanceof Error ? error.message : '未知错误';
            setState({ loading: false, result: null, error: errorMsg });
            message.error(`权限同步失败: ${errorMsg}`);
        }
    }, []);

    /**
     * 仅同步菜单
     */
    const syncMenus = useCallback(async () => {
        setState({ loading: true, result: null, error: null });
        try {
            const result = await syncMenusOnly();
            setState({ loading: false, result, error: null });
            if (result.success) {
                message.success(`菜单同步成功: 创建${result.created}个，更新${result.updated}个`);
            } else {
                message.warning(`同步完成，但有 ${result.errors.length} 个错误`);
            }
        } catch (error) {
            const errorMsg = error instanceof Error ? error.message : '未知错误';
            setState({ loading: false, result: null, error: errorMsg });
            message.error(`菜单同步失败: ${errorMsg}`);
        }
    }, []);

    /**
     * 仅为超管分配权限
     */
    const assignSuperAdmin = useCallback(async () => {
        setState({ loading: true, result: null, error: null });
        try {
            const result = await assignSuperAdminOnly();
            setState({ loading: false, result: null, error: result.success ? null : result.message });
            if (result.success) {
                message.success(result.message);
            } else {
                message.error(result.message);
            }
        } catch (error) {
            const errorMsg = error instanceof Error ? error.message : '未知错误';
            setState({ loading: false, result: null, error: errorMsg });
            message.error(`分配权限失败: ${errorMsg}`);
        }
    }, []);

    /**
     * 重置状态
     */
    const reset = useCallback(() => {
        setState({ loading: false, result: null, error: null });
    }, []);

    return {
        state,
        syncAll,
        forceSyncAll,
        syncPermissions,
        syncMenus,
        assignSuperAdmin,
        reset,
    };
};

export default useSync;
