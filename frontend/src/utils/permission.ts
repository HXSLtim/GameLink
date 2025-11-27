/**
 * 权限工具函数
 * 提供非React环境下的权限检查能力
 */

/**
 * 权限存储管理器
 * 用于在非React组件中访问权限数据
 */
class PermissionStore {
    private permissions: string[] = [];
    private listeners: Set<(permissions: string[]) => void> = new Set();

    /**
     * 设置权限列表
     * @param permissions 权限码数组
     */
    setPermissions(permissions: string[]): void {
        this.permissions = permissions;
        this.notifyListeners();
    }

    /**
     * 获取权限列表
     * @returns 权限码数组
     */
    getPermissions(): string[] {
        return this.permissions;
    }

    /**
     * 清空权限
     */
    clearPermissions(): void {
        this.permissions = [];
        this.notifyListeners();
    }

    /**
     * 订阅权限变化
     * @param listener 监听函数
     * @returns 取消订阅函数
     */
    subscribe(listener: (permissions: string[]) => void): () => void {
        this.listeners.add(listener);
        return () => {
            this.listeners.delete(listener);
        };
    }

    /**
     * 通知所有监听者
     */
    private notifyListeners(): void {
        this.listeners.forEach(listener => listener(this.permissions));
    }

    /**
     * 检查是否拥有指定权限
     * @param permission 权限码或权限码数组
     * @param mode 检查模式：any-任一满足，all-全部满足
     * @returns 是否拥有权限
     */
    hasPermission(permission: string | string[], mode: 'any' | 'all' = 'any'): boolean {
        if (!this.permissions.length) {
            return false;
        }

        // 超级管理员拥有所有权限
        if (this.permissions.includes('*')) {
            return true;
        }

        const permissionList = Array.isArray(permission) ? permission : [permission];

        if (mode === 'all') {
            return permissionList.every(p => this.permissions.includes(p));
        } else {
            return permissionList.some(p => this.permissions.includes(p));
        }
    }
}

/** 权限存储单例 */
export const permissionStore = new PermissionStore();

/**
 * 检查是否拥有指定权限
 * 可在非React环境中使用
 *
 * @param permission 权限码或权限码数组
 * @param mode 检查模式：any-任一满足，all-全部满足
 * @returns 是否拥有权限
 *
 * @example
 * ```ts
 * // 在普通函数中检查权限
 * if (hasPermission('admin.games.delete')) {
 *     // 执行删除操作
 * }
 *
 * // 检查多个权限
 * if (hasPermission(['admin.games.create', 'admin.games.update'], 'any')) {
 *     // 有创建或编辑权限
 * }
 * ```
 */
export function hasPermission(permission: string | string[], mode: 'any' | 'all' = 'any'): boolean {
    return permissionStore.hasPermission(permission, mode);
}

/**
 * 检查是否拥有全部指定权限
 *
 * @param permissions 权限码数组
 * @returns 是否拥有全部权限
 */
export function hasAllPermissions(permissions: string[]): boolean {
    return permissionStore.hasPermission(permissions, 'all');
}

/**
 * 检查是否拥有任一指定权限
 *
 * @param permissions 权限码数组
 * @returns 是否拥有任一权限
 */
export function hasAnyPermission(permissions: string[]): boolean {
    return permissionStore.hasPermission(permissions, 'any');
}

/**
 * 根据权限过滤操作列表
 *
 * @param actions 操作列表
 * @param permissionKey 权限码字段名，默认为 'permission'
 * @returns 过滤后的操作列表
 *
 * @example
 * ```ts
 * const allActions = [
 *     { key: 'edit', label: '编辑', permission: 'admin.games.update' },
 *     { key: 'delete', label: '删除', permission: 'admin.games.delete' },
 *     { key: 'view', label: '查看', permission: 'admin.games.read' },
 * ];
 *
 * const allowedActions = filterActionsByPermission(allActions);
 * // 返回用户有权限的操作
 * ```
 */
export function filterActionsByPermission<T extends Record<string, unknown>>(
    actions: T[],
    permissionKey: keyof T = 'permission' as keyof T
): T[] {
    return actions.filter(action => {
        const permission = action[permissionKey];
        if (!permission) {
            return true; // 没有权限要求的操作默认允许
        }
        if (typeof permission === 'string') {
            return hasPermission(permission);
        }
        if (Array.isArray(permission)) {
            return hasPermission(permission as string[]);
        }
        return true;
    });
}

/**
 * 权限码解析工具
 */
export const PermissionParser = {
    /**
     * 从权限码中提取模块名
     * @param permissionCode 权限码，如 'admin.games.create'
     * @returns 模块名，如 'admin'
     */
    getModule(permissionCode: string): string {
        return permissionCode.split('.')[0] || '';
    },

    /**
     * 从权限码中提取资源名
     * @param permissionCode 权限码，如 'admin.games.create'
     * @returns 资源名，如 'games'
     */
    getResource(permissionCode: string): string {
        return permissionCode.split('.')[1] || '';
    },

    /**
     * 从权限码中提取操作名
     * @param permissionCode 权限码，如 'admin.games.create'
     * @returns 操作名，如 'create'
     */
    getAction(permissionCode: string): string {
        return permissionCode.split('.')[2] || '';
    },

    /**
     * 构建权限码
     * @param module 模块名
     * @param resource 资源名
     * @param action 操作名
     * @returns 权限码
     */
    build(module: string, resource: string, action: string): string {
        return `${module}.${resource}.${action}`;
    },

    /**
     * 检查权限码是否匹配模式
     * 支持通配符 *
     *
     * @param permissionCode 权限码
     * @param pattern 模式，如 'admin.games.*' 或 'admin.*.*'
     * @returns 是否匹配
     */
    matches(permissionCode: string, pattern: string): boolean {
        const codeParts = permissionCode.split('.');
        const patternParts = pattern.split('.');

        if (patternParts.length !== codeParts.length) {
            return false;
        }

        return patternParts.every((part, index) => {
            return part === '*' || part === codeParts[index];
        });
    }
};

export default {
    permissionStore,
    hasPermission,
    hasAllPermissions,
    hasAnyPermission,
    filterActionsByPermission,
    PermissionParser
};
