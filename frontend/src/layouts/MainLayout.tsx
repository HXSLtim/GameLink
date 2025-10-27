import { useMemo } from 'react';
import { Layout, Menu, Breadcrumb, Typography, Space } from '@arco-design/web-react';
import { IconApps, IconHome, IconSettings, IconUser, IconLock } from '@arco-design/web-react/icon';
import type { UIMatch } from 'react-router-dom';
import { Link, Outlet, useLocation, useMatches } from 'react-router-dom';
import { ThemeSwitcher } from '../components/ThemeSwitcher';
import { Footer } from '../components/Footer';
import styles from './MainLayout.module.less';

const { Header, Sider, Content } = Layout;

interface RouteHandle {
  crumb?: string;
}

/**
 * Breadcrumb navigation component
 *
 * @component
 * @description Displays breadcrumb navigation based on current route
 */
const Breadcrumbs = () => {
  const matches = useMatches() as UIMatch<unknown, RouteHandle>[];

  const crumbItems = useMemo(() => matches.filter((match) => match.handle?.crumb), [matches]);

  return (
    <Breadcrumb className={styles.breadcrumb}>
      <Breadcrumb.Item>
        <Link to="/">
          <IconHome /> 首页
        </Link>
      </Breadcrumb.Item>
      {crumbItems.map((match) => (
        <Breadcrumb.Item key={match.id}>{match.handle?.crumb}</Breadcrumb.Item>
      ))}
    </Breadcrumb>
  );
};

/**
 * Main layout component
 *
 * @component
 * @description Provides the main application layout with header, sidebar navigation, and content area
 */
export const MainLayout = () => {
  const location = useLocation();

  const selectedKeys = useMemo(() => [location.pathname], [location.pathname]);

  return (
    <Layout className={styles.layout}>
      <Header>
        <div className={styles.headerContent}>
          <div className={styles.logo}>
            <IconApps className={styles.logoIcon} />
            <Typography.Title className={styles.title} heading={5}>
              GameLink 管理端
            </Typography.Title>
          </div>
          <Space size={12}>
            <ThemeSwitcher />
          </Space>
        </div>
      </Header>
      <Layout>
        <Sider breakpoint="lg" width={220}>
          <Menu selectedKeys={selectedKeys} className={styles.menu}>
            <Menu.Item key="/">
              <Link to="/">
                <IconHome /> 总览
              </Link>
            </Menu.Item>
            <Menu.Item key="/users">
              <Link to="/users">
                <IconUser /> 用户
              </Link>
            </Menu.Item>
            <Menu.Item key="/orders">
              <Link to="/orders">
                <IconApps /> 订单
              </Link>
            </Menu.Item>
            <Menu.Item key="/permissions">
              <Link to="/permissions">
                <IconLock /> 权限
              </Link>
            </Menu.Item>
            <Menu.Item key="/settings">
              <Link to="/settings">
                <IconSettings /> 设置
              </Link>
            </Menu.Item>
            <Menu.Item key="/login">
              <Link to="/login">登录</Link>
            </Menu.Item>
          </Menu>
        </Sider>
        <Layout>
          <Breadcrumbs />
          <Content className={styles.content}>
            <Outlet />
          </Content>
          <Footer />
        </Layout>
      </Layout>
    </Layout>
  );
};
