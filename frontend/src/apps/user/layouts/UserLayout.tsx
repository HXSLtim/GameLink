/**
 * 用户端主布局
 */
import React from 'react';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { DiscordLayout } from '@/shared/components/DiscordLayout';
import { Avatar } from '@/shared/components/Avatar';
import './UserLayout.less';

const UserLayout: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();

  // TODO: 从状态管理或API获取用户信息
  const currentUser = {
    id: 1,
    name: '用户昵称',
    avatar: 'https://via.placeholder.com/40',
    status: 'online' as const,
  };

  const navigationItems = [
    {
      path: '/user/discover',
      icon: '🏠',
      label: '发现',
      description: '浏览陪玩师',
    },
    {
      path: '/user/players',
      icon: '🎮',
      label: '陪玩师',
      description: '找陪玩',
    },
    {
      path: '/user/orders',
      icon: '📋',
      label: '订单',
      description: '我的订单',
    },
    {
      path: '/user/messages',
      icon: '💬',
      label: '消息',
      description: '聊天',
    },
    {
      path: '/user/profile',
      icon: '👤',
      label: '我的',
      description: '个人中心',
    },
  ];

  const isActive = (path: string) => {
    return location.pathname.startsWith(path);
  };

  const handleNavigate = (path: string) => {
    navigate(path);
  };

  return (
    <DiscordLayout>
      {/* 左侧导航栏 */}
      <DiscordLayout.Sidebar>
        <div className="user-layout-sidebar">
          {/* Logo */}
          <div
            className="user-layout-sidebar__logo"
            onClick={() => handleNavigate('/user/discover')}
          >
            <span className="user-layout-sidebar__logo-icon">🎮</span>
          </div>

          {/* 导航按钮 */}
          <div className="user-layout-sidebar__nav">
            {navigationItems.map((item) => (
              <button
                key={item.path}
                className={`user-layout-sidebar__item ${
                  isActive(item.path) ? 'user-layout-sidebar__item--active' : ''
                }`}
                onClick={() => handleNavigate(item.path)}
                title={item.description}
              >
                <span className="user-layout-sidebar__item-icon">{item.icon}</span>
              </button>
            ))}
          </div>

          {/* 用户头像 */}
          <div className="user-layout-sidebar__user">
            <Avatar
              src={currentUser.avatar}
              alt={currentUser.name}
              size="medium"
              status={currentUser.status}
              showStatus
              onClick={() => handleNavigate('/user/profile')}
            />
          </div>
        </div>
      </DiscordLayout.Sidebar>

      {/* 主内容区域 - 由子路由填充 */}
      <Outlet />
    </DiscordLayout>
  );
};

export default UserLayout;
