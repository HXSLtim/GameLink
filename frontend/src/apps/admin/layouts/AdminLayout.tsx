/**
 * 管理后台布局
 */
import { useState } from 'react';
import { Outlet, useNavigate } from 'react-router-dom';
import { Layout, Menu, Avatar, Dropdown, Button, Space, Typography } from '@arco-design/web-react';
import {
  IconDashboard,
  IconUser,
  IconShoppingCart,
  IconDollarCircle,
  IconSettings,
  IconPoweroff,
  IconMenu,
  IconMenuFold,
} from '@arco-design/web-react/icon';
import { useAuth } from '@/shared/hooks/useAuth';
import './AdminLayout.less';

const { Header, Sider, Content } = Layout;
const { Text } = Typography;

/**
 * 管理后台布局组件
 */
export const AdminLayout = () => {
  const navigate = useNavigate();
  const { user, logout } = useAuth();
  const [collapsed, setCollapsed] = useState(false);

  /**
   * 菜单项
   */
  const menuItems = [
    { key: '/admin', icon: <IconDashboard />, title: '工作台' },
    { key: '/admin/users', icon: <IconUser />, title: '用户管理' },
    { key: '/admin/orders', icon: <IconShoppingCart />, title: '订单管理' },
    { key: '/admin/finance', icon: <IconDollarCircle />, title: '财务管理' },
    { key: '/admin/settings', icon: <IconSettings />, title: '系统设置' },
  ];

  /**
   * 用户下拉菜单
   */
  const userDroplist = (
    <Menu>
      <Menu.Item key="profile" onClick={() => navigate('/admin/profile')}>
        <IconUser style={{ marginRight: 8 }} />
        个人中心
      </Menu.Item>
      <Menu.Divider />
      <Menu.Item key="logout" onClick={logout}>
        <IconPoweroff style={{ marginRight: 8 }} />
        退出登录
      </Menu.Item>
    </Menu>
  );

  return (
    <Layout className="admin-layout">
      <Sider collapsed={collapsed} collapsible trigger={null} breakpoint="lg">
        <div className="layout-logo">
          {!collapsed && <Text bold>GameLink Admin</Text>}
        </div>
        <Menu onClickMenuItem={(key) => navigate(key)}>
          {menuItems.map((item) => (
            <Menu.Item key={item.key}>
              {item.icon}
              {item.title}
            </Menu.Item>
          ))}
        </Menu>
      </Sider>

      <Layout>
        <Header className="layout-header">
          <Button
            shape="circle"
            icon={collapsed ? <IconMenu /> : <IconMenuFold />}
            onClick={() => setCollapsed(!collapsed)}
          />
          <Dropdown droplist={userDroplist} position="br">
            <Space style={{ cursor: 'pointer' }}>
              <Avatar size={32}>{user?.username?.[0]?.toUpperCase() || 'A'}</Avatar>
              <Text>{user?.username || '管理员'}</Text>
            </Space>
          </Dropdown>
        </Header>

        <Content className="layout-content">
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
};
