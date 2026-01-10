import type { Menu } from '@/api/admin';
import type { RouteConfig } from '@/router/types';
import { getComponent } from '@/router/componentMap';

import { logger } from '@/utils/logger';
/**
 * 从菜单生成路由配置
 * @param menus 菜单列表（树形结构）
 * @param parentPath 父级路径（用于计算相对路径）
 */
export const generateRoutesFromMenus = (menus: Menu[], parentPath: string = '/admin'): RouteConfig[] => {
    return menus.map(menu => {
        const Component = getComponent(menu.component);
        const hasChildren = menu.children && menu.children.length > 0;

        // 计算相对路径
        let path = menu.path;
        
        // 移除父级路径前缀，得到相对路径
        if (path.startsWith(parentPath + '/')) {
            path = path.substring(parentPath.length + 1);
        } else if (path === parentPath) {
            // 如果路径与父路径相同，说明是 index route
            path = '';
        } else if (path.startsWith('/admin/')) {
            // 根级菜单，移除 /admin/ 前缀
            path = path.substring('/admin/'.length);
        } else if (path.startsWith('/admin')) {
            path = path.substring('/admin'.length) || '';
        }
        
        // 移除开头的斜杠（相对路径不需要）
        if (path.startsWith('/')) {
            path = path.substring(1);
        }

        // 是否是 index route（只有没有子菜单且路径为空时才是 index）
        const isIndex = path === '' && !hasChildren;

        // 递归处理子菜单，传入当前菜单的完整路径作为父路径
        const children = hasChildren
            ? generateRoutesFromMenus(menu.children!, menu.path) 
            : undefined;

        // 构建路由配置
        const route: RouteConfig = {
            path: isIndex ? undefined : path,
            element: Component ? <Component /> : null,
            meta: {
                title: menu.name,
                requiresAuth: true,
                roles: ['ADMIN'],
                permission: menu.permission || undefined,
            },
            children,
            index: isIndex
        };

        // 调试日志
        logger.info(`[Route] ${menu.name}: path="${path}", fullPath="${menu.path}", component="${menu.component}", hasChildren=${hasChildren}, isIndex=${isIndex}`);
        
        return route;
    });
};
