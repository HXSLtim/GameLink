import { useRoutes, type RouteObject } from 'react-router-dom';
import { routes } from './routes';
import RouteGuard from './Guard';
import type { RouteConfig } from './types';


import { useAdmin } from '@/context/useAdmin';
import { generateRoutesFromMenus } from '@/utils/dynamicRoutes';
import { useMemo } from 'react';

const renderRoutes = (routes: RouteConfig[]): RouteObject[] => {
    return routes.map((route) => {
        const element = route.element;

        // Wrap element in Guard if needed
        const guardedElement = (route.meta?.requiresAuth || route.meta?.roles || route.meta?.permission) ? (
            <RouteGuard
                requiresAuth={route.meta?.requiresAuth}
                roles={route.meta?.roles}
                permission={route.meta?.permission}
            >
                {element}
            </RouteGuard>
        ) : element;

        // Handle index routes (cannot have path and index: true at the same time)
        if (route.index) {
            return {
                index: true,
                element: guardedElement,
            };
        }

        return {
            path: route.path,
            element: guardedElement,
            children: route.children ? renderRoutes(route.children) : undefined,
        };
    });
};

// 需要保留的静态子路由（不在菜单中但需要访问的页面）
const STATIC_CHILD_ROUTES = [
    'sys/role/:id/permissions',  // 角色权限配置页面
    'sys/user/:id/portrait',     // 用户画像页面
    'biz/service/create',        // 服务项目创建页面
    'biz/service/:id',           // 服务项目详情页面（编辑使用模态弹窗）
    'profile',                   // 个人中心页面（通过用户下拉菜单访问）
];

const AppRouter = () => {
    const { menus } = useAdmin();

    const finalRoutes = useMemo(() => {
        // Clone static routes (only keep non-admin routes like login, 404, etc.)
        const allRoutes = routes.filter(r => r.path !== '/admin');

        // Find the original admin route for layout and meta
        const originalAdminRoute = routes.find(r => r.path === '/admin');

        // Find or create admin route
        if (menus.length > 0) {
            const dynamicRoutes = generateRoutesFromMenus(menus);
            
            // 从静态路由中提取需要保留的子路由
            const staticChildRoutes = originalAdminRoute?.children?.filter(
                child => child.path && STATIC_CHILD_ROUTES.includes(child.path)
            ) || [];
            
            // 合并动态路由和静态子路由
            const mergedChildren = [...dynamicRoutes, ...staticChildRoutes];
            
            // Create admin route with merged children
            const adminRoute: RouteConfig = {
                path: '/admin',
                element: originalAdminRoute?.element,
                meta: originalAdminRoute?.meta,
                children: mergedChildren
            };
            
            allRoutes.push(adminRoute);
        } else {
            // If no menus loaded yet, keep original admin route structure
            if (originalAdminRoute) {
                allRoutes.push(originalAdminRoute);
            }
        }

        return allRoutes;
    }, [menus]);

    const element = useRoutes(renderRoutes(finalRoutes));
    return element;
};

export default AppRouter;
