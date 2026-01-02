/**
 * 应用初始化服务
 * 在应用启动时自动同步路由和权限
 */
import { syncApi, type SyncResult } from '@/api/sync';
import { ADMIN_MENUS, ADMIN_PERMISSIONS } from '@/config/adminRoutes';
import { authApi } from '@/api/auth';

import { logger } from '@/utils/logger';
/**
 * 初始化配置
 */
export interface InitConfig {
    /** 是否同步菜单 */
    syncMenus?: boolean;
    /** 是否同步权限 */
    syncPermissions?: boolean;
    /** 是否为超管分配所有权限 */
    assignSuperAdminPermissions?: boolean;
    /** 开发模式下是否跳过同步 */
    skipInDev?: boolean;
    /** 是否显示控制台日志 */
    verbose?: boolean;
}

/**
 * 初始化结果
 */
export interface InitResult {
    success: boolean;
    menuSync?: SyncResult;
    permissionSync?: SyncResult;
    superAdminAssign?: { success: boolean; message: string };
    errors: string[];
    duration: number;
}

/**
 * 默认配置
 */
const defaultConfig: InitConfig = {
    syncMenus: true,
    syncPermissions: true,
    assignSuperAdminPermissions: true,
    skipInDev: false,
    verbose: import.meta.env.DEV,
};

/**
 * 用户数据接口（匹配后端实际返回格式）
 */
interface UserData {
    user?: {
        role?: string;
    };
}

/**
 * 检查是否有管理员权限（已登录且是管理员）
 */
const hasAdminAccess = async (): Promise<boolean> => {
    const token = localStorage.getItem('token');
    if (!token) {
        logger.info('[Init] No token found');
        return false;
    }

    try {
        const axiosResponse = await authApi.getMe();

        logger.info('[Init] getMe response:', axiosResponse);

        // axios 响应拦截器返回完整响应对象，实际数据在 response.data 中
        const api = (axiosResponse as { data: unknown }).data as {
            success?: boolean;
            code?: number;
            message?: string;
            data?: UserData;
        };

        if (!api) {
            logger.info('[Init] No ApiResponse payload');
            return false;
        }

        const isSuccess = api.success === true || api.code === 200;
        if (!isSuccess) {
            logger.info('[Init] ApiResponse not successful:', api);
            return false;
        }

        if (!api.data) {
            logger.info('[Init] ApiResponse.data missing');
            return false;
        }

        const userData = api.data as UserData;
        const role = userData.user?.role || '';
        logger.info('[Init] User role:', role);
        const isAdmin = ['admin', 'superAdmin', 'ADMIN', 'CS', 'FINANCE'].includes(role);
        logger.info('[Init] Is admin:', isAdmin);
        return isAdmin;
    } catch (error) {
        logger.info('[Init] API call failed:', error);
        return false;
    }
};

/**
 * 日志输出
 */
const log = (verbose: boolean, ...args: unknown[]) => {
    if (verbose) {
        logger.info('[Init]', ...args);
    }
};

/**
 * 初始化应用 - 批量同步菜单和权限到后端
 */
export const initApp = async (config: InitConfig = {}): Promise<InitResult> => {
    const startTime = Date.now();
    const cfg = { ...defaultConfig, ...config };
    const result: InitResult = {
        success: true,
        errors: [],
        duration: 0,
    };

    log(cfg.verbose!, '开始应用初始化...');

    if (cfg.skipInDev && import.meta.env.DEV) {
        log(cfg.verbose!, '开发模式跳过同步');
        result.duration = Date.now() - startTime;
        return result;
    }

    const isAdmin = await hasAdminAccess();
    if (!isAdmin) {
        log(cfg.verbose!, '非管理员用户，跳过同步');
        result.duration = Date.now() - startTime;
        return result;
    }

    try {
        log(cfg.verbose!, '开始批量同步...');
        const batchResult = await syncApi.batchSync(
            cfg.syncMenus ? ADMIN_MENUS : [],
            cfg.syncPermissions ? ADMIN_PERMISSIONS : [],
            cfg.assignSuperAdminPermissions ?? true
        );
        log(cfg.verbose!, '批量同步完成:', batchResult);

        if (batchResult.permissionSync) {
            result.permissionSync = batchResult.permissionSync;
            if (!batchResult.permissionSync.success) {
                result.errors.push(...batchResult.permissionSync.errors);
            }
        }

        if (batchResult.menuSync) {
            result.menuSync = batchResult.menuSync;
            if (!batchResult.menuSync.success) {
                result.errors.push(...batchResult.menuSync.errors);
            }
        }

        if (batchResult.superAdminAssign) {
            result.superAdminAssign = batchResult.superAdminAssign;
            if (!result.superAdminAssign.success) {
                log(cfg.verbose!, '超管权限分配警告:', result.superAdminAssign.message);
            } else {
                log(cfg.verbose!, result.superAdminAssign.message);
            }
        }
    } catch (error) {
        result.success = false;
        result.errors.push(`初始化失败: ${error instanceof Error ? error.message : '未知错误'}`);
    }

    result.duration = Date.now() - startTime;
    log(cfg.verbose!, `初始化完成，耗时 ${result.duration}ms`);

    return result;
};

const INIT_STORAGE_KEY = 'app_init_timestamp';
const INIT_INTERVAL = 24 * 60 * 60 * 1000; // 24小时

const shouldReInit = (): boolean => {
    const lastInit = localStorage.getItem(INIT_STORAGE_KEY);
    if (!lastInit) return true;
    const lastTime = parseInt(lastInit, 10);
    return Date.now() - lastTime > INIT_INTERVAL;
};

const markInitComplete = () => {
    localStorage.setItem(INIT_STORAGE_KEY, Date.now().toString());
};

/**
 * 智能初始化（带缓存）- 每24小时只执行一次完整同步
 */
export const smartInit = async (config: InitConfig = {}): Promise<InitResult | null> => {
    if (!shouldReInit()) {
        const cfg = { ...defaultConfig, ...config };
        log(cfg.verbose!, '24小时内已初始化，跳过');
        return null;
    }

    const result = await initApp(config);
    if (result.success) {
        markInitComplete();
    }
    return result;
};

/**
 * 强制重新初始化
 */
export const forceInit = async (config: InitConfig = {}): Promise<InitResult> => {
    localStorage.removeItem(INIT_STORAGE_KEY);
    return initApp(config);
};

export default {
    initApp,
    smartInit,
    forceInit,
};
