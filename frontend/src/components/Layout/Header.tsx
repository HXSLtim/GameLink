import { useState } from 'react';
import { useTheme } from 'contexts/ThemeContext';
import type { BreadcrumbItem } from '../Breadcrumb/Breadcrumb';
import styles from './Header.module.less';

export interface HeaderProps {
  /** 用户信息 */
  user?: {
    username: string;
    role: string;
  };
  /** 退出登录回调 */
  onLogout?: () => void;
  /** 切换侧边栏 */
  onToggleSidebar?: () => void;
  /** 面包屑数据 */
  breadcrumbs?: BreadcrumbItem[];
}

// 菜单图标
const MenuIcon = () => (
  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <line x1="3" y1="6" x2="21" y2="6" strokeWidth="2" strokeLinecap="round" />
    <line x1="3" y1="12" x2="21" y2="12" strokeWidth="2" strokeLinecap="round" />
    <line x1="3" y1="18" x2="21" y2="18" strokeWidth="2" strokeLinecap="round" />
  </svg>
);

// 用户图标
const UserIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path
      d="M20 21V19C20 17.9391 19.5786 16.9217 18.8284 16.1716C18.0783 15.4214 17.0609 15 16 15H8C6.93913 15 5.92172 15.4214 5.17157 16.1716C4.42143 16.9217 4 17.9391 4 19V21"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <circle cx="12" cy="7" r="4" strokeWidth="2" />
  </svg>
);

// 退出图标
const LogoutIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path
      d="M9 21H5C4.46957 21 3.96086 20.7893 3.58579 20.4142C3.21071 20.0391 3 19.5304 3 19V5C3 4.46957 3.21071 3.96086 3.58579 3.58579C3.96086 3.21071 4.46957 3 5 3H9"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <polyline
      points="16 17 21 12 16 7"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <line
      x1="21"
      y1="12"
      x2="9"
      y2="12"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

// 太阳图标（浅色模式）
const SunIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <circle cx="12" cy="12" r="5" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    <line
      x1="12"
      y1="1"
      x2="12"
      y2="3"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <line
      x1="12"
      y1="21"
      x2="12"
      y2="23"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <line
      x1="4.22"
      y1="4.22"
      x2="5.64"
      y2="5.64"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <line
      x1="18.36"
      y1="18.36"
      x2="19.78"
      y2="19.78"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <line
      x1="1"
      y1="12"
      x2="3"
      y2="12"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <line
      x1="21"
      y1="12"
      x2="23"
      y2="12"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <line
      x1="4.22"
      y1="19.78"
      x2="5.64"
      y2="18.36"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <line
      x1="18.36"
      y1="5.64"
      x2="19.78"
      y2="4.22"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

// 月亮图标（深色模式）
const MoonIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path
      d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

// 箭头图标
const ArrowIcon = () => (
  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <polyline
      points="9 18 15 12 9 6"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

export const Header: React.FC<HeaderProps> = ({ user, onLogout, onToggleSidebar, breadcrumbs }) => {
  const [showUserMenu, setShowUserMenu] = useState(false);
  const { effective, setMode } = useTheme();

  const handleLogout = () => {
    setShowUserMenu(false);
    onLogout?.();
  };

  const toggleTheme = (event: React.MouseEvent<HTMLButtonElement>) => {
    // 获取点击坐标
    const x = event.clientX;
    const y = event.clientY;

    // 在浅色和深色之间切换，传递点击坐标用于扩散动画
    setMode(effective === 'light' ? 'dark' : 'light', x, y);
  };

  return (
    <header className={styles.header}>
      <div className={styles.left}>
        <button className={styles.menuButton} onClick={onToggleSidebar} aria-label="切换菜单">
          <MenuIcon />
        </button>
        <h1 className={styles.logo}>GameLink</h1>
        
        {/* 面包屑 */}
        {breadcrumbs && breadcrumbs.length > 0 && (
          <>
            <span className={styles.separator}>|</span>
            <nav className={styles.breadcrumb} aria-label="面包屑导航">
              {breadcrumbs.map((item, index) => {
                const isLast = index === breadcrumbs.length - 1;
                return (
                  <span key={index} className={styles.breadcrumbItem}>
                    {item.path && !isLast ? (
                      <a href={item.path} className={styles.breadcrumbLink}>
                        {item.label}
                      </a>
                    ) : (
                      <span className={isLast ? styles.breadcrumbCurrent : styles.breadcrumbLabel}>
                        {item.label}
                      </span>
                    )}
                    {!isLast && (
                      <span className={styles.breadcrumbSeparator}>
                        <ArrowIcon />
                      </span>
                    )}
                  </span>
                );
              })}
            </nav>
          </>
        )}
      </div>

      <div className={styles.right}>
        {/* 主题切换按钮 */}
        <button
          className={styles.themeButton}
          onClick={toggleTheme}
          aria-label={effective === 'light' ? '切换到深色模式' : '切换到浅色模式'}
          title={effective === 'light' ? '切换到深色模式' : '切换到浅色模式'}
        >
          {effective === 'light' ? <MoonIcon /> : <SunIcon />}
        </button>

        {user && (
          <div className={styles.userSection}>
            <button
              className={styles.userButton}
              onClick={() => setShowUserMenu(!showUserMenu)}
              aria-label="用户菜单"
            >
              <UserIcon />
              <span className={styles.username}>{user.username}</span>
            </button>

            {showUserMenu && (
              <div className={styles.userMenu}>
                <div className={styles.userInfo}>
                  <div className={styles.userInfoName}>{user.username}</div>
                  <div className={styles.userInfoRole}>{user.role}</div>
                </div>
                <div className={styles.divider}></div>
                <button className={styles.logoutButton} onClick={handleLogout}>
                  <LogoutIcon />
                  <span>退出登录</span>
                </button>
              </div>
            )}
          </div>
        )}
      </div>
    </header>
  );
};
