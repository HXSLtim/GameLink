import { ReactNode } from 'react';
import { Link, useLocation, matchPath } from 'react-router-dom';
import styles from './Sidebar.module.less';

export interface MenuItem {
  key: string;
  label: string;
  icon?: ReactNode;
  path: string;
  children?: MenuItem[];
}

export interface SidebarProps {
  /** 菜单项 */
  menuItems: MenuItem[];
  /** 是否收起 */
  collapsed?: boolean;
}

/**
 * 侧边栏导航组件
 * 支持精确匹配和子路由匹配
 */
export const Sidebar: React.FC<SidebarProps> = ({ menuItems, collapsed = false }) => {
  const location = useLocation();

  /**
   * 判断菜单项是否激活
   * 支持：
   * 1. 精确路径匹配
   * 2. 子路由匹配（如 /orders/:id 激活 /orders）
   * 3. 忽略查询参数
   */
  const isActive = (path: string): boolean => {
    // 精确匹配（忽略查询参数和hash）
    if (location.pathname === path) {
      return true;
    }

    // 子路由匹配
    // 例如：/orders/:id 应该激活 /orders 菜单项
    if (path !== '/' && path !== '/dashboard') {
      const match = matchPath(
        {
          path: `${path}/*`,
          caseSensitive: false,
          end: false,
        },
        location.pathname,
      );

      if (match) {
        return true;
      }
    }

    return false;
  };

  return (
    <aside className={`${styles.sidebar} ${collapsed ? styles.collapsed : ''}`}>
      <nav className={styles.nav}>
        {menuItems.map((item) => (
          <Link
            key={item.key}
            to={item.path}
            className={`${styles.menuItem} ${isActive(item.path) ? styles.active : ''}`}
          >
            {item.icon && <span className={styles.icon}>{item.icon}</span>}
            {!collapsed && <span className={styles.label}>{item.label}</span>}
          </Link>
        ))}
      </nav>
    </aside>
  );
};
