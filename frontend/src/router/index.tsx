import { useRoutes } from 'react-router-dom';
import { routes } from './routes';
import RouteGuard from './Guard';
import { RouteConfig } from './types';
import { ReactNode } from 'react';

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
            children: route.children ? renderRoutes(route.children) : undefined
        };
    });
};

const AppRouter = () => {
    const element = useRoutes(renderRoutes(routes));
    return element;
};

export default AppRouter;
