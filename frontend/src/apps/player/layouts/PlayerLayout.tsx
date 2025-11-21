/**
 * 陪玩师端布局
 */
import { useState } from 'react';
import { Outlet, useNavigate } from 'react-router-dom';
import { Layout, Menu, Avatar, Dropdown, Button, Space, Typography } from '@arco-design/web-react';
import {
  IconDashboard,
  IconShoppingCart,
  IconDollarCircle,
  IconUser,
  IconPoweroff,
  IconMenu,
  IconMenuFold,
} from '@arco-design/web-react/icon';
import { useAuth } from '@/shared/hooks/useAuth';
import './PlayerLayout.less';

const { Header, Sider, Content } = Layout;
const { Text } = Typography;

/**
 * 陪玩师端布局组件
 */
export const PlayerLayout = () => {
  const navigate = useNavigate();
  const { user, logout } = useAuth();
  const [collapsed, setCollapsed] = useState(false);

  /**
   * 菜单项
   */
  const menuItems = [
    { key: '/player', icon: <IconDashboard />, title: '工作台' },
    { key: '/player/orders', icon: <IconShoppingCart />, title: '我的订单' },
    { key: '/player/earnings', icon: <IconDollarCircle />, title: '收益管理' },
    { key: '/player/profile', icon: <IconUser />, title: '个人资料' },
  ];

  /**
   * 用户下拉菜单
   */
  const userDroplist = (
    <Menu>
      <Menu.Item key="profile" onClick={() => navigate('/player/profile')}>
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
    <Layout className="player-layout">
      <Sider collapsed={collapsed} collapsible trigger={null} breakpoint="lg">
        <div className="layout-logo">
          {!collapsed && <Text bold>GameLink Player</Text>}
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
              <Avatar size={32}>{user?.username?.[0]?.toUpperCase() || 'P'}</Avatar>
              <Text>{user?.username || '陪玩师'}</Text>
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
