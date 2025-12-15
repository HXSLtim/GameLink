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

const AppRouter = () => {
    const { menus } = useAdmin();

    const finalRoutes = useMemo(() => {
        // Clone static routes (only keep non-admin routes like login, 404, etc.)
        const allRoutes = routes.filter(r => r.path !== '/admin');

        // Find or create admin route
        if (menus.length > 0) {
            const dynamicRoutes = generateRoutesFromMenus(menus);
            
            // Find the original admin route for layout and meta
            const originalAdminRoute = routes.find(r => r.path === '/admin');
            
            // Create admin route with dynamic children only
            const adminRoute: RouteConfig = {
                path: '/admin',
                element: originalAdminRoute?.element,
                meta: originalAdminRoute?.meta,
                children: dynamicRoutes
            };
            
            allRoutes.push(adminRoute);
        } else {
            // If no menus loaded yet, keep original admin route structure
            const originalAdminRoute = routes.find(r => r.path === '/admin');
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
