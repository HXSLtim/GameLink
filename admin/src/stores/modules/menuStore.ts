/**
 * Menu Store
 * Menu and navigation state management with permission filtering
 * Requirements: 8.1, 8.2, 8.4 - Menu permission filtering and caching
 */

import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { MenuItem } from '../types';
import { useAuthStore } from './authStore';

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
          // TODO: Replace with actual API call
          // import { adminApi } from '@/api/admin';
          // const response = await adminApi.getMyMenus();

          // Simulate API call
          await new Promise(resolve => setTimeout(resolve, 300));

          // Mock menu data
          const mockMenus: MenuItem[] = [
            {
              id: 1,
              key: 'dashboard',
              label: '仪表盘',
              path: '/admin/dashboard',
              icon: 'DashboardOutlined',
              sort: 1,
              permission: 'admin.dashboard.read',
            },
            {
              id: 2,
              key: 'users',
              label: '用户管理',
              path: '/admin/users',
              icon: 'UserOutlined',
              sort: 2,
              permission: 'admin.users.read',
              children: [
                {
                  id: 21,
                  key: 'users-list',
                  label: '用户列表',
                  path: '/admin/users/list',
                  sort: 1,
                  permission: 'admin.users.read',
                },
                {
                  id: 22,
                  key: 'users-block',
                  label: '黑名单',
                  path: '/admin/users/block',
                  sort: 2,
                  permission: 'admin.users.block',
                },
              ],
            },
            {
              id: 3,
              key: 'players',
              label: '陪玩师管理',
              path: '/admin/players',
              icon: 'TeamOutlined',
              sort: 3,
              permission: 'admin.players.read',
            },
            {
              id: 4,
              key: 'orders',
              label: '订单管理',
              path: '/admin/orders',
              icon: 'ShoppingOutlined',
              sort: 4,
              permission: 'admin.orders.read',
            },
            {
              id: 5,
              key: 'chat',
              label: '聊天管理',
              path: '/admin/chat',
              icon: 'MessageOutlined',
              sort: 5,
              permission: 'admin.chat.read',
            },
            {
              id: 6,
              key: 'payment',
              label: '支付管理',
              path: '/admin/payment',
              icon: 'PayCircleOutlined',
              sort: 6,
              permission: 'admin.payment.read',
            },
            {
              id: 7,
              key: 'system',
              label: '系统管理',
              path: '/admin/system',
              icon: 'SettingOutlined',
              sort: 7,
              permission: 'admin.system.read',
              children: [
                {
                  id: 71,
                  key: 'system-admin',
                  label: '管理员管理',
                  path: '/admin/system/admin',
                  sort: 1,
                  permission: 'admin.admin.read',
                },
                {
                  id: 72,
                  key: 'system-role',
                  label: '角色管理',
                  path: '/admin/system/role',
                  sort: 2,
                  permission: 'admin.role.read',
                },
                {
                  id: 73,
                  key: 'system-permission',
                  label: '权限管理',
                  path: '/admin/system/permission',
                  sort: 3,
                  permission: 'admin.permission.read',
                },
                {
                  id: 74,
                  key: 'system-menu',
                  label: '菜单管理',
                  path: '/admin/system/menu',
                  sort: 4,
                  permission: 'admin.menu.read',
                },
              ],
            },
          ];

          set({ rawMenus: mockMenus, loading: false });

          // Auto-filter menus based on current permissions
          const permissions = useAuthStore.getState().userInfo?.permissions || [];
          const filteredMenus = get().filterMenusByPermission(mockMenus, permissions);

          set({ menus: filteredMenus });

        } catch (error: unknown) {
          console.error('Failed to fetch menus:', error);
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
