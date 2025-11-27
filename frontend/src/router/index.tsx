import { useRoutes } from 'react-router-dom';
import { routes } from './routes';
import RouteGuard from './Guard';
import type { RouteConfig } from './types';


import { useAdmin } from '@/context/AdminContext';
import { generateRoutesFromMenus } from '@/utils/dynamicRoutes';
import { useMemo } from 'react';

const renderRoutes = (routes: RouteConfig[]): any[] => {
    return routes.map((route) => {
        const element = route.element;

        // Wrap element in Guard if needed
        const guardedElement = (route.meta?.requiresAuth || route.meta?.roles) ? (
            <RouteGuard
                requiresAuth={route.meta?.requiresAuth}
                roles={route.meta?.roles}
            >
                {element}
            </RouteGuard>
        ) : element;

        return {
            path: route.path,
            element: guardedElement,
            children: route.children ? renderRoutes(route.children) : undefined,
            index: route.index
        };
    });
};

const AppRouter = () => {
    const { menus } = useAdmin();

    const finalRoutes = useMemo(() => {
        // Clone static routes
        const allRoutes = [...routes];

        // Find admin route
        const adminRouteIndex = allRoutes.findIndex(r => r.path === '/admin');
        if (adminRouteIndex !== -1 && menus.length > 0) {
            const dynamicRoutes = generateRoutesFromMenus(menus);

            // Merge dynamic routes into admin children
            const adminRoute = { ...allRoutes[adminRouteIndex] };
            adminRoute.children = [
                ...(adminRoute.children || []),
                ...dynamicRoutes
            ];
            allRoutes[adminRouteIndex] = adminRoute;
        }

        return allRoutes;
    }, [menus]);

    const element = useRoutes(renderRoutes(finalRoutes));
    return element;
};

export default AppRouter;
