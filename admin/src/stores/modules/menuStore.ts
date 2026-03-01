/**
 * Menu Store
 * Menu and navigation state management with permission filtering
 * Requirements: 8.1, 8.2, 8.4 - Menu permission filtering and caching
 */

import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { MenuItem } from '../types';
import { useAuthStore } from './authStore';
import { adminApi } from '@/api/admin';

import { logger } from '@/utils/logger';
interface MenuState {
  // State
  /** Raw menus from API (unfiltered) */
  rawMenus: MenuItem[];
  /** Filtered menus based on user permissions */
  menus: MenuItem[];
  /** Loading state */
  loading: boolean;
  /** Menu collapsed state */
  collapsed: boolean;
  /** Expanded menu keys */
  openKeys: string[];
  /** Selected menu keys */
  selectedKeys: string[];

  // Actions
  /** Fetch menus from API */
  fetchMenus: () => Promise<void>;
  /** Set collapsed state */
  setCollapsed: (collapsed: boolean) => void;
  /** Toggle collapsed state */
  toggleCollapsed: () => void;
  /** Set expanded keys */
  setOpenKeys: (keys: string[]) => void;
  /** Set selected keys */
  setSelectedKeys: (keys: string[]) => void;
  /** Filter menus by user permissions */
  filterMenusByPermission: (rawMenus: MenuItem[], permissions: string[]) => MenuItem[];

  // Selectors
  /** Get menu by key */
  getMenuByKey: (key: string) => MenuItem | undefined;
  /** Get breadcrumb items for a menu key */
  getBreadcrumbs: (key: string) => MenuItem[];
  /** Get all menu paths (for route generation) */
  getMenuPaths: () => Record<string, string>;
}

export const useMenuStore = create<MenuState>()(
  persist(
    (set, get) => ({
      // Initial State
      rawMenus: [],
      menus: [],
      loading: false,
      collapsed: false,
      openKeys: [],
      selectedKeys: [],

      // Actions
      fetchMenus: async () => {
        set({ loading: true });

        try {
          const response = await adminApi.getMyMenus();

          // Transform API response to MenuItem format
          const transformMenu = (menu: any): MenuItem => ({
            id: menu.id,
            key: menu.name, // Use 'name' as 'key' for routing
            label: menu.name,
            path: menu.path,
            icon: menu.icon,
            parentId: menu.parentId,
            sort: menu.order, // API uses 'order' instead of 'sort'
            permission: menu.permission,
            children: menu.children?.map(transformMenu),
          });

          const menus: MenuItem[] = response.data.data.map(transformMenu);

          set({ rawMenus: menus, loading: false });

          // Auto-filter menus based on current permissions
          const permissions = useAuthStore.getState().userInfo?.permissions || [];
          const filteredMenus = get().filterMenusByPermission(menus, permissions);

          set({ menus: filteredMenus });

        } catch (error: unknown) {
          logger.error('Failed to fetch menus:', error);
          set({ loading: false, menus: [] });
        }
      },

      setCollapsed: (collapsed: boolean) => {
        set({ collapsed });
      },

      toggleCollapsed: () => {
        set(state => ({ collapsed: !state.collapsed }));
      },

      setOpenKeys: (keys: string[]) => {
        set({ openKeys: keys });
      },

      setSelectedKeys: (keys: string[]) => {
        set({ selectedKeys: keys });
      },

      /**
       * Filter menus by user permissions
       * Supports nested menu structure and wildcard permissions
       */
      filterMenusByPermission: (rawMenus: MenuItem[], permissions: string[]) => {
        const filterMenu = (menus: MenuItem[]): MenuItem[] => {
          return menus
            .filter(menu => {
              // Super admin sees all menus
              if (permissions.includes('*')) {
                return true;
              }
              // Check menu permission
              if (!menu.permission) {
                return true;
              }
              return permissions.includes(menu.permission);
            })
            .map(menu => ({
              ...menu,
              // Recursively filter children
              children: menu.children ? filterMenu(menu.children) : undefined,
            }))
            .filter(menu => {
              // Remove parent menus with no visible children
              if (menu.children && menu.children.length === 0) {
                return false;
              }
              return true;
            })
            .sort((a, b) => a.sort - b.sort);
        };

        return filterMenu(rawMenus);
      },

      // Selectors
      getMenuByKey: (key: string) => {
        const findMenu = (menus: MenuItem[]): MenuItem | undefined => {
          for (const menu of menus) {
            if (menu.key === key) return menu;
            if (menu.children) {
              const found = findMenu(menu.children);
              if (found) return found;
            }
          }
          return undefined;
        };

        return findMenu(get().menus);
      },

      getBreadcrumbs: (key: string) => {
        const breadcrumbs: MenuItem[] = [];
        const menu = get().getMenuByKey(key);

        if (menu) {
          breadcrumbs.push(menu);
        }

        return breadcrumbs;
      },

      getMenuPaths: () => {
        const paths: Record<string, string> = {};

        const collectPaths = (menus: MenuItem[]) => {
          menus.forEach(menu => {
            if (menu.path) {
              paths[menu.key] = menu.path;
            }
            if (menu.children) {
              collectPaths(menu.children);
            }
          });
        };

        collectPaths(get().menus);
        return paths;
      },
    }),
    {
      name: 'menu-storage',
      // Cache raw menus and collapsed state (filtered menus will be recomputed)
      partialize: (state) => ({
        rawMenus: state.rawMenus,
        collapsed: state.collapsed,
      }),
    }
  )
);

/**
 * Listen to permission changes and auto-refresh menus
 * Requirements: 8.4 - Permission changes trigger menu updates
 */
// Use standard subscribe for zustand v5 with persist middleware
useAuthStore.subscribe(
  (state, prevState) => {
    const permissions = state.userInfo?.permissions;
    const prevPermissions = prevState.userInfo?.permissions;

    // Only refilter if permissions exist and actually changed
    if (permissions && permissions !== prevPermissions) {
      const { rawMenus, filterMenusByPermission } = useMenuStore.getState();
      const filteredMenus = filterMenusByPermission(rawMenus, permissions);
      useMenuStore.setState({ menus: filteredMenus });
    }
  }
);

export type { MenuState };
