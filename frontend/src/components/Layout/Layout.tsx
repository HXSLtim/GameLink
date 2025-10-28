import { ReactNode, useState } from 'react';
import { Header, HeaderProps } from './Header';
import { Sidebar, SidebarProps } from './Sidebar';
import styles from './Layout.module.less';

export interface LayoutProps {
  /** 页面内容 */
  children: ReactNode;
  /** Header 配置 */
  headerProps?: HeaderProps;
  /** Sidebar 配置 */
  sidebarProps?: SidebarProps;
  /** 是否显示侧边栏 */
  showSidebar?: boolean;
}

export const Layout: React.FC<LayoutProps> = ({
  children,
  headerProps,
  sidebarProps,
  showSidebar = true,
}) => {
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);

  const handleToggleSidebar = () => {
    setSidebarCollapsed(!sidebarCollapsed);
  };

  return (
    <div className={styles.layout}>
      <Header {...headerProps} onToggleSidebar={handleToggleSidebar} />
      <div className={styles.container}>
        {showSidebar && sidebarProps && <Sidebar {...sidebarProps} collapsed={sidebarCollapsed} />}
        <main className={styles.main}>{children}</main>
      </div>
    </div>
  );
};
