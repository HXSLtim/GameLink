
import type { Menu } from '@/api/admin';
import type { RouteConfig } from '@/router/types';
import { getComponent } from '@/router/componentMap';

export const generateRoutesFromMenus = (menus: Menu[]): RouteConfig[] => {
    return menus.map(menu => {
        const Component = getComponent(menu.component);

        // Handle path: remove leading slash if present to make it relative to parent
        const path = menu.path.startsWith('/') ? menu.path.substring(1) : menu.path;

        const route: RouteConfig = {
            path: path,
            element: Component ? <Component /> : null,
            meta: {
                title: menu.name,
                requiresAuth: true,
                roles: ['ADMIN'], // Assuming all dynamic menus are for admins
                // permissions: [menu.permission] 
            },
            children: menu.children && menu.children.length > 0 ? generateRoutesFromMenus(menu.children) : undefined
        };
        return route;
    });
};
