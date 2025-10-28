import { ReactNode } from 'react';
import { Link, useLocation } from 'react-router-dom';
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

export const Sidebar: React.FC<SidebarProps> = ({ menuItems, collapsed = false }) => {
  const location = useLocation();

  const isActive = (path: string) => {
    return location.pathname === path;
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
