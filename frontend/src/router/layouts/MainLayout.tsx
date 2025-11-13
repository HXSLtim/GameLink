import { Outlet, useNavigate } from 'react-router-dom';
import { Layout, RouteCache } from 'components';
import { useAuth } from 'contexts/AuthContext';
import { useBreadcrumb } from '../../hooks/useBreadcrumb';
import { useRouteCache } from '../../hooks/useRouteCache';
import { FEATURE_FLAGS } from '../../config';

// Dashboard 图标
const DashboardIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <rect
      x="3"
      y="3"
      width="7"
      height="7"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <rect
      x="14"
      y="3"
      width="7"
      height="7"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <rect
      x="14"
      y="14"
      width="7"
      height="7"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <rect
      x="3"
      y="14"
      width="7"
      height="7"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

// 订单图标
const OrdersIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path d="M9 11L12 14L22 4" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    <path
      d="M21 12V19C21 19.5304 20.7893 20.0391 20.4142 20.4142C20.0391 20.7893 19.5304 21 19 21H5C4.46957 21 3.96086 20.7893 3.58579 20.4142C3.21071 20.0391 3 19.5304 3 19V5C3 4.46957 3.21071 3.96086 3.58579 3.58579C3.96086 3.21071 4.46957 3 5 3H16"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

const AssignmentsIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path
      d="M4 15.5C4 12.4624 6.46243 10 9.5 10C12.5376 10 15 12.4624 15 15.5V18H12"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path
      d="M6 10C6 6.68629 8.68629 4 12 4C15.3137 4 18 6.68629 18 10V12"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path d="M18 15L21 18L18 21" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
  </svg>
);

// 游戏图标
const GamesIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <rect x="2" y="7" width="20" height="14" rx="2" strokeWidth="2" />
    <path d="M16 3L8 3" strokeWidth="2" strokeLinecap="round" />
    <circle cx="8" cy="14" r="1" fill="currentColor" />
    <circle cx="12" cy="14" r="1" fill="currentColor" />
    <path d="M16 12V16" strokeWidth="2" strokeLinecap="round" />
    <path d="M14 14H18" strokeWidth="2" strokeLinecap="round" />
  </svg>
);

// 陪玩师图标
const PlayersIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <circle cx="9" cy="7" r="4" strokeWidth="2" />
    <path
      d="M2 21V19C2 17.3431 3.34315 16 5 16H13C14.6569 16 16 17.3431 16 19V21"
      strokeWidth="2"
    />
    <path d="M16 11L18 13L22 9" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
  </svg>
);

// 用户图标
const UsersIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path
      d="M17 21V19C17 17.9391 16.5786 16.9217 15.8284 16.1716C15.0783 15.4214 14.0609 15 13 15H5C3.93913 15 2.92172 15.4214 2.17157 16.1716C1.42143 16.9217 1 17.9391 1 19V21"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <circle cx="9" cy="7" r="4" strokeWidth="2" />
    <path
      d="M23 21V19C22.9993 18.1137 22.7044 17.2528 22.1614 16.5523C21.6184 15.8519 20.8581 15.3516 20 15.13"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path
      d="M16 3.13C16.8604 3.35031 17.623 3.85071 18.1676 4.55232C18.7122 5.25392 19.0078 6.11683 19.0078 7.005C19.0078 7.89318 18.7122 8.75608 18.1676 9.45769C17.623 10.1593 16.8604 10.6597 16 10.88"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

// 支付图标
const PaymentsIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <rect x="2" y="5" width="20" height="14" rx="2" strokeWidth="2" />
    <path d="M2 10H22" strokeWidth="2" />
  </svg>
);

// 组件演示图标
const ComponentsIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <rect x="3" y="3" width="7" height="7" strokeWidth="2" />
    <rect x="14" y="3" width="7" height="7" strokeWidth="2" />
    <rect x="3" y="14" width="7" height="7" strokeWidth="2" />
    <rect x="14" y="14" width="7" height="7" strokeWidth="2" />
  </svg>
);

// 社区图标
const CommunityIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path d="M5 20V9L12 4L19 9V20" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    <path d="M9 21V14H15V21" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    <circle cx="12" cy="10" r="1" fill="currentColor" />
  </svg>
);

// 评价图标
const ReviewsIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path
      d="M21 11.5C21 16.75 16.5 21 12 21C11.5 21 10.96 20.93 10.5 20.85C9.5 21.5 8 22 6.5 22C6.5 22 6.78 20.5 6.5 19.5C4.5 18 3 15.5 3 11.5C3 6.25 7.5 2 12 2C16.5 2 21 6.25 21 11.5Z"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path d="M9 11H9.01" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    <path d="M12 11H12.01" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    <path d="M15 11H15.01" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
  </svg>
);

// 报表图标
const ReportsIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path d="M3 3V18C3 19.1046 3.89543 20 5 20H21" strokeWidth="2" strokeLinecap="round" />
    <path d="M7 14L12 9L16 13L21 8" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
  </svg>
);

// 权限图标
const PermissionsIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <rect x="3" y="11" width="18" height="10" rx="2" strokeWidth="2" />
    <path
      d="M7 11V7C7 5.67392 7.52678 4.40215 8.46447 3.46447C9.40215 2.52678 10.6739 2 12 2C13.3261 2 14.5979 2.52678 15.5355 3.46447C16.4732 4.40215 17 5.67392 17 7V11"
      strokeWidth="2"
    />
    <circle cx="12" cy="16" r="1" fill="currentColor" />
  </svg>
);

// 设置图标
const SettingsIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <circle cx="12" cy="12" r="3" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    <path d="M12 1V3" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    <path d="M12 21V23" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    <path d="M4.22 4.22L5.64 5.64" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    <path
      d="M18.36 18.36L19.78 19.78"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path d="M1 12H3" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    <path d="M21 12H23" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    <path d="M4.22 19.78L5.64 18.36" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    <path d="M18.36 5.64L19.78 4.22" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
  </svg>
);

/**
 * 主布局组件
 */
export const MainLayout = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  // 菜单配置
  const menuItems = [
    {
      key: 'dashboard',
      label: '仪表盘',
      icon: <DashboardIcon />,
      path: '/dashboard',
    },
    {
      key: 'orders',
      label: '订单管理',
      icon: <OrdersIcon />,
      path: '/orders',
    },
    {
      key: 'assignments',
      label: '客服指派',
      icon: <AssignmentsIcon />,
      path: '/assignments',
    },
    {
      key: 'games',
      label: '游戏管理',
      icon: <GamesIcon />,
      path: '/games',
    },
    {
      key: 'players',
      label: '陪玩师管理',
      icon: <PlayersIcon />,
      path: '/players',
    },
    {
      key: 'users',
      label: '用户管理',
      icon: <UsersIcon />,
      path: '/users',
    },
    {
      key: 'payments',
      label: '支付管理',
      icon: <PaymentsIcon />,
      path: '/payments',
    },
    {
      key: 'reviews',
      label: '评价管理',
      icon: <ReviewsIcon />,
      path: '/reviews',
    },
    {
      key: 'community',
      label: '社区动态',
      icon: <CommunityIcon />,
      path: '/community',
    },
    {
      key: 'reports',
      label: '数据报表',
      icon: <ReportsIcon />,
      path: '/reports',
    },
    {
      key: 'permissions',
      label: '权限管理',
      icon: <PermissionsIcon />,
      path: '/permissions',
    },
    {
      key: 'settings',
      label: '系统设置',
      icon: <SettingsIcon />,
      path: '/settings',
    },
    ...(FEATURE_FLAGS.showcase.enableComponentsRoute
      ? [
          {
            key: 'components',
            label: '组件演示',
            icon: <ComponentsIcon />,
            path: '/components',
          },
        ]
      : []),
  ];

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  // 获取面包屑
  const breadcrumbs = useBreadcrumb();

  // 路由缓存配置
  const cacheControl = useRouteCache({
    enabled: true,
    maxCache: 10,
    cacheRoutes: [
      '/users',
      '/orders',
      '/assignments',
      '/games',
      '/players',
      '/payments',
      '/reviews',
    ],
    excludeRoutes: ['/login', '/register'],
  });

  return (
    <Layout
      headerProps={{
        user: user
          ? {
              username: user.username || user.name,
              role: (user.role as string) || 'user',
            }
          : undefined,
        onLogout: handleLogout,
        breadcrumbs: breadcrumbs.length > 0 ? breadcrumbs : undefined,
      }}
      sidebarProps={{
        menuItems,
      }}
      showSidebar={true}
    >
      <RouteCache
        enabled={cacheControl.config.enabled}
        maxCache={cacheControl.config.maxCache}
        cacheRoutes={cacheControl.config.cacheRoutes}
        excludeRoutes={cacheControl.config.excludeRoutes}
      >
        <Outlet />
      </RouteCache>
    </Layout>
  );
};
