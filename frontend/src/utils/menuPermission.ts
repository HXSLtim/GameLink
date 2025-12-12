/**
 * 菜单权限过滤工具
 * 根据用户权限过滤菜单项
 * Requirements: 8.1, 8.2
 */

import type { Menu } from '@/api/admin';

/**
 * 检查用户是否拥有指定权限
 * @param userPermissions 用户权限列表
 * @param permission 需要检查的权限
 * @returns 是否拥有权限
 */
export const hasPermission = (userPermissions: string[], permission: string): boolean => {
    // 超级管理员拥有所有权限
    if (userPermissions.includes('*')) {
        return true;
    }
    // 空权限码表示不需要权限
    if (!permission) {
        return true;
    }
    return userPermissions.includes(permission);
};



/**
 * 根据用户权限过滤菜单
 * - 如果菜单项有权限要求，检查用户是否拥有该权限
 * - 如果菜单项是父菜单（有子菜单），检查是否有任何可访问的子菜单
 * - 无子菜单权限时隐藏父菜单
 * 
 * @param menus 原始菜单列表
 * @param userPermissions 用户权限列表
 * @returns 过滤后的菜单列表
 */
export const filterMenusByPermission = (
    menus: Menu[],
    userPermissions: string[]
): Menu[] => {
    // 超级管理员返回所有菜单
    if (userPermissions.includes('*')) {
        return menus;
    }

    return menus
        .filter(menu => {
            // 检查菜单是否可见
            if (menu.visible === false) {
                return false;
            }

            // 如果菜单有子菜单，检查是否有可访问的子菜单
            if (menu.children && menu.children.length > 0) {
                // 递归过滤子菜单
                const filteredChildren = filterMenusByPermission(menu.children, userPermissions);
                // 如果没有可访问的子菜单，隐藏父菜单
                return filteredChildren.length > 0;
            }

            // 叶子菜单：检查权限
            // 如果没有设置权限，默认可访问
            if (!menu.permission) {
                return true;
            }

            return hasPermission(userPermissions, menu.permission);
        })
        .map(menu => {
            // 递归过滤子菜单
            if (menu.children && menu.children.length > 0) {
                return {
                    ...menu,
                    children: filterMenusByPermission(menu.children, userPermissions),
                };
            }
            return menu;
        });
};

/**
 * 检查用户是否有访问指定路径的权限
 * @param menus 菜单列表
 * @param path 路径
 * @param userPermissions 用户权限列表
 * @returns 是否有权限访问
 */
export const hasRoutePermission = (
    menus: Menu[],
    path: string,
    userPermissions: string[]
): boolean => {
    // 超级管理员拥有所有权限
    if (userPermissions.includes('*')) {
        return true;
    }

    // 在菜单中查找对应路径
    const findMenuByPath = (menuList: Menu[], targetPath: string): Menu | null => {
        for (const menu of menuList) {
            if (menu.path === targetPath) {
                return menu;
            }
            if (menu.children && menu.children.length > 0) {
                const found = findMenuByPath(menu.children, targetPath);
                if (found) {
                    return found;
                }
            }
        }
        return null;
    };

    const menu = findMenuByPath(menus, path);
    
    // 如果路径不在菜单中，默认允许访问（可能是公共页面）
    if (!menu) {
        return true;
    }

    // 如果菜单没有设置权限，默认允许访问
    if (!menu.permission) {
        return true;
    }

    return hasPermission(userPermissions, menu.permission);
};

/**
 * 获取用户可访问的第一个菜单路径
 * 用于无权限时的重定向
 * @param menus 菜单列表
 * @param userPermissions 用户权限列表
 * @returns 第一个可访问的菜单路径，如果没有则返回 null
 */
export const getFirstAccessiblePath = (
    menus: Menu[],
    userPermissions: string[]
): string | null => {
    const filteredMenus = filterMenusByPermission(menus, userPermissions);
    
    if (filteredMenus.length === 0) {
        return null;
    }

    // 递归查找第一个叶子菜单
    const findFirstLeaf = (menuList: Menu[]): string | null => {
        for (const menu of menuList) {
            if (menu.children && menu.children.length > 0) {
                const leafPath = findFirstLeaf(menu.children);
                if (leafPath) {
                    return leafPath;
                }
            } else {
                return menu.path;
            }
        }
        return null;
    };

    return findFirstLeaf(filteredMenus);
};
