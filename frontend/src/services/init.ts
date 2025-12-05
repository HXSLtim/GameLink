/**
 * 应用初始化服务
 * 在应用启动时自动同步路由和权限
 */
import { syncApi, type SyncResult } from '@/api/sync';
import { ADMIN_MENUS, ADMIN_PERMISSIONS } from '@/config/adminRoutes';

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
 * 通过调用后端API验证，避免客户端JWT解析安全漏洞
 */
const hasAdminAccess = async (): Promise<boolean> => {
    const token = localStorage.getItem('token');
    if (!token) {
        console.log('[Init] No token found');
        return false;
    }

    try {
        // ✅ 安全修复: 通过后端API验证用户角色，不在客户端解析JWT
        // 避免了JWT伪造导致的权限绕过漏洞
        const { authApi } = await import('@/api/auth');
        const response = await authApi.getMe();

        console.log('[Init] getMe response:', response);

        // AxiosResponse<ApiResponse<LoginResponse>>
        const api = response?.data;
        if (!api) {
            console.log('[Init] No ApiResponse payload');
            return false;
        }
        if (api.code !== 200) {
            console.log('[Init] ApiResponse.code is not 200:', api.code);
            return false;
        }
        if (!api.data) {
            console.log('[Init] ApiResponse.data missing');
            return false;
        }

        // ApiResponse<LoginResponse> -> data: LoginResponse
        const userData = api.data as unknown as UserData;
        const role = userData.user?.role || '';
        console.log('[Init] User role:', role);
        const isAdmin = ['admin', 'superAdmin', 'ADMIN', 'CS', 'FINANCE'].includes(role);
        console.log('[Init] Is admin:', isAdmin);
        return isAdmin;
    } catch (error) {
        // API调用失败（网络错误、认证失败等）
        console.log('[Init] API call failed:', error);
        return false;
    }
};

/**
 * 日志输出
 */
const log = (verbose: boolean, ...args: unknown[]) => {
    if (verbose) {
        console.log('[Init]', ...args);
    }
};

/**
 * 初始化应用
 * 同步菜单和权限到后端，为超管分配所有权限
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

    // 开发模式跳过检查
    if (cfg.skipInDev && import.meta.env.DEV) {
        log(cfg.verbose!, '开发模式跳过同步');
        result.duration = Date.now() - startTime;
        return result;
    }

    // 检查管理员权限
    const isAdmin = await hasAdminAccess();
    if (!isAdmin) {
        log(cfg.verbose!, '非管理员用户，跳过同步');
        result.duration = Date.now() - startTime;
        return result;
    }

    try {
        // 1. 同步权限
        if (cfg.syncPermissions) {
            log(cfg.verbose!, '同步权限...');
            result.permissionSync = await syncApi.syncPermissions(ADMIN_PERMISSIONS);
            if (!result.permissionSync.success) {
                result.success = false;
                result.errors.push(...result.permissionSync.errors);
            }
            log(cfg.verbose!, '权限同步完成:', result.permissionSync);
        }

        // 2. 同步菜单
        if (cfg.syncMenus) {
            log(cfg.verbose!, '同步菜单...');
            result.menuSync = await syncApi.syncMenus(ADMIN_MENUS);
            if (!result.menuSync.success) {
                result.success = false;
                result.errors.push(...result.menuSync.errors);
            }
            log(cfg.verbose!, '菜单同步完成:', result.menuSync);
        }

        // 3. 为超管分配所有权限
        if (cfg.assignSuperAdminPermissions) {
            log(cfg.verbose!, '为超管分配权限...');
            result.superAdminAssign = await syncApi.assignAllPermissionsToSuperAdmin();
            if (!result.superAdminAssign.success) {
                // 超管分配失败不算整体失败，只记录警告
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

/**
 * 仅同步权限
 */
export const syncPermissionsOnly = async (): Promise<SyncResult> => {
    return syncApi.syncPermissions(ADMIN_PERMISSIONS);
};

/**
 * 仅同步菜单
 */
export const syncMenusOnly = async (): Promise<SyncResult> => {
    return syncApi.syncMenus(ADMIN_MENUS);
};

/**
 * 仅分配超管权限
 */
export const assignSuperAdminOnly = async (): Promise<{ success: boolean; message: string }> => {
    return syncApi.assignAllPermissionsToSuperAdmin();
};

/**
 * 初始化存储键
 */
const INIT_STORAGE_KEY = 'app_init_timestamp';
const INIT_INTERVAL = 24 * 60 * 60 * 1000; // 24小时

/**
 * 检查是否需要重新初始化
 */
const shouldReInit = (): boolean => {
    const lastInit = localStorage.getItem(INIT_STORAGE_KEY);
    if (!lastInit) return true;

    const lastTime = parseInt(lastInit, 10);
    return Date.now() - lastTime > INIT_INTERVAL;
};

/**
 * 标记初始化完成
 */
const markInitComplete = () => {
    localStorage.setItem(INIT_STORAGE_KEY, Date.now().toString());
};

/**
 * 智能初始化（带缓存）
 * 每24小时只执行一次完整同步
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
    syncPermissionsOnly,
    syncMenusOnly,
    assignSuperAdminOnly,
};
