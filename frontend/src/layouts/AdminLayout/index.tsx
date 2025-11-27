/**
 * Admin布局组件
 * 使用Ant Design Layout构建管理后台框架
 */
import React, { useState, useEffect, useMemo } from 'react';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import {
    Layout,
    Menu,
    Avatar,
    Dropdown,
    Breadcrumb,
    theme,
    Button,
    Space,
    Badge,
    Spin,
    Grid,
    Drawer,
    ConfigProvider,
} from 'antd';
import type { MenuProps } from 'antd';
import {
    MenuFoldOutlined,
    MenuUnfoldOutlined,
    UserOutlined,
    LogoutOutlined,
    SettingOutlined,
    BellOutlined,
    DashboardOutlined,
    TeamOutlined,
    AppstoreOutlined,
    ShoppingCartOutlined,
    SafetyCertificateOutlined,
    SettingFilled,
    MenuOutlined,
} from '@ant-design/icons';
import { useAdmin } from '@/context/AdminContext';
import { useTheme } from '@/context/ThemeContext';
import { authApi } from '@/api/auth';
import { ThemeToggle } from '@/components';
import styles from './index.module.css';

const { Header, Sider, Content } = Layout;

type MenuItem = Required<MenuProps>['items'][number];

/**
 * 菜单配置
 */
const menuItems: MenuItem[] = [
    {
        key: '/admin',
        icon: <DashboardOutlined />,
        label: '仪表盘',
    },
    {
        key: '/admin/sys',
        icon: <SettingFilled />,
        label: '系统管理',
        children: [
            { key: '/admin/sys/user', icon: <TeamOutlined />, label: '用户管理' },
            { key: '/admin/sys/role', icon: <SafetyCertificateOutlined />, label: '角色管理' },
            { key: '/admin/sys/permission', icon: <SafetyCertificateOutlined />, label: '权限管理' },
            { key: '/admin/sys/menu', icon: <MenuOutlined />, label: '菜单管理' },
        ],
    },
    {
        key: '/admin/biz',
        icon: <AppstoreOutlined />,
        label: '业务管理',
        children: [
            { key: '/admin/biz/game', icon: <AppstoreOutlined />, label: '游戏管理' },
            { key: '/admin/biz/player', icon: <TeamOutlined />, label: '陪玩师管理' },
            { key: '/admin/biz/order', icon: <ShoppingCartOutlined />, label: '订单管理' },
            { key: '/admin/biz/service', icon: <AppstoreOutlined />, label: '服务项目' },
        ],
    },
    {
        key: '/admin/settings',
        icon: <SettingOutlined />,
        label: '系统设置',
    },
];

/**
 * 面包屑映射
 */
const breadcrumbMap: Record<string, string> = {
    '/admin': '仪表盘',
    '/admin/sys': '系统管理',
    '/admin/sys/user': '用户管理',
    '/admin/sys/role': '角色管理',
    '/admin/sys/permission': '权限管理',
    '/admin/sys/menu': '菜单管理',
    '/admin/biz': '业务管理',
    '/admin/biz/game': '游戏管理',
    '/admin/biz/player': '陪玩师管理',
    '/admin/biz/order': '订单管理',
    '/admin/biz/service': '服务项目',
    '/admin/settings': '系统设置',
};

/**
 * AdminLayout组件
 */
const AdminLayout: React.FC = () => {
    const [collapsed, setCollapsed] = useState(false);
    const navigate = useNavigate();
    const location = useLocation();
    const { loading } = useAdmin();
    const { token } = theme.useToken();
    const screens = Grid.useBreakpoint();
    const { mode } = useTheme();

    // 用户信息
    const [userInfo, setUserInfo] = useState<{ username: string; avatar?: string }>({
        username: 'Admin',
    });

    useEffect(() => {
        const storedUser = localStorage.getItem('user_info');
        if (storedUser) {
            try {
                const parsed = JSON.parse(storedUser);
                setUserInfo({ username: parsed.username || 'Admin', avatar: parsed.avatar });
            } catch {
                // ignore
            }
        }
    }, []);

    // 响应式处理：屏幕变窄时自动收起
    useEffect(() => {
        if (!screens.md) {
            setCollapsed(true);
        } else {
            setCollapsed(false);
        }
    }, [screens.md]);

    /**
     * 获取当前选中的菜单项
     */
    const selectedKeys = useMemo(() => {
        const path = location.pathname;
        // 精确匹配或找到最长匹配的父路径
        const keys = Object.keys(breadcrumbMap).filter(key => path.startsWith(key));
        return keys.length > 0 ? [keys.sort((a, b) => b.length - a.length)[0]] : ['/admin'];
    }, [location.pathname]);

    /**
     * 获取展开的菜单项
     */
    const openKeys = useMemo(() => {
        const path = location.pathname;
        const keys: string[] = [];
        if (path.startsWith('/admin/sys')) keys.push('/admin/sys');
        if (path.startsWith('/admin/biz')) keys.push('/admin/biz');
        return keys;
    }, [location.pathname]);

    /**
     * 生成面包屑
     */
    const breadcrumbItems = useMemo(() => {
        const pathSnippets = location.pathname.split('/').filter(i => i);
        const items = [{ title: '首页', href: '/admin' }];

        let currentPath = '';
        pathSnippets.forEach((snippet) => {
            currentPath += `/${snippet}`;
            if (breadcrumbMap[currentPath]) {
                items.push({
                    title: breadcrumbMap[currentPath],
                    href: currentPath,
                });
            }
        });

        return items;
    }, [location.pathname]);

    /**
     * 菜单点击
     */
    const handleMenuClick: MenuProps['onClick'] = ({ key }) => {
        navigate(key);
        // 移动端点击菜单后自动收起
        if (!screens.md) {
            setCollapsed(true);
        }
    };

    /**
     * 退出登录
     */
    const handleLogout = async () => {
        try {
            await authApi.logout();
        } catch {
            // ignore
        }
        localStorage.removeItem('token');
        localStorage.removeItem('user_info');
        localStorage.removeItem('user_role');
        navigate('/login');
    };

    /**
     * 用户下拉菜单
     */
    const userMenuItems: MenuProps['items'] = [
        {
            key: 'profile',
            icon: <UserOutlined />,
            label: '个人中心',
        },
        {
            key: 'settings',
            icon: <SettingOutlined />,
            label: '账号设置',
        },
        { type: 'divider' },
        {
            key: 'logout',
            icon: <LogoutOutlined />,
            label: '退出登录',
            danger: true,
        },
    ];

    const handleUserMenuClick: MenuProps['onClick'] = ({ key }) => {
        if (key === 'logout') {
            handleLogout();
        } else if (key === 'settings') {
            navigate('/admin/settings');
        }
    };

    if (loading) {
        return (
            <div className={styles.loading} style={{ background: token.colorBgLayout }}>
                <Spin size="large" tip="加载中...">
                    <div style={{ padding: 50 }} />
                </Spin>
            </div>
        );
    }

    // 计算左侧边距
    const getMarginLeft = () => {
        if (!screens.md) return 0;
        return collapsed ? 80 : 220;
    };

    const renderMenu = () => (
        <>
            {/* Logo */}
            <div className={styles.logo} style={{ borderBottom: `1px solid ${token.colorSplit}` }}>
                <img src="/logo.svg" alt="GameLink" className={styles.logoIcon} />
                {(!collapsed || !screens.md) && (
                    <span className={styles.logoText} style={{ color: token.colorPrimary }}>
                        GameLink
                    </span>
                )}
            </div>

            {/* 菜单 */}
            <ConfigProvider
                theme={{
                    components: {
                        Menu: {
                            itemBorderRadius: 6,
                            itemMarginInline: 8,
                        },
                    },
                }}
            >
                <Menu
                    mode="inline"
                    theme={mode}
                    selectedKeys={selectedKeys}
                    defaultOpenKeys={openKeys}
                    items={menuItems}
                    onClick={handleMenuClick}
                    className={styles.menu}
                    style={{ background: 'transparent', borderRight: 0 }}
                />
            </ConfigProvider>
        </>
    );

    return (
        <Layout className={styles.layout} style={{ minHeight: '100vh' }}>
            {/* 侧边栏 - 移动端使用 Drawer */}
            {!screens.md ? (
                <Drawer
                    placement="left"
                    onClose={() => setCollapsed(true)}
                    open={!collapsed}
                    width={220}
                    styles={{ body: { padding: 0 }, header: { display: 'none' } }}
                    closable={false}
                >
                    <div style={{ display: 'flex', flexDirection: 'column', height: '100%', background: token.colorBgContainer }}>
                        {renderMenu()}
                    </div>
                </Drawer>
            ) : (
                /* 侧边栏 - 桌面端使用 Sider */
                <Sider
                    trigger={null}
                    collapsible
                    collapsed={collapsed}
                    width={220}
                    className={styles.sider}
                    theme={mode}
                    style={{
                        background: token.colorBgContainer,
                        boxShadow: '2px 0 8px rgba(0, 0, 0, 0.05)',
                        zIndex: 100
                    }}
                >
                    {renderMenu()}
                </Sider>
            )}

            <Layout>
                {/* 顶部栏 */}
                <Header
                    className={styles.header}
                    style={{
                        background: token.colorBgContainer,
                        marginLeft: getMarginLeft(),
                        width: !screens.md ? '100%' : `calc(100% - ${getMarginLeft()}px)`,
                        transition: 'all 0.2s',
                    }}
                >
                    <div className={styles.headerLeft}>
                        <Button
                            type="text"
                            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
                            onClick={() => setCollapsed(!collapsed)}
                            className={styles.trigger}
                        />
                        <Breadcrumb
                            items={breadcrumbItems.map(item => ({
                                title: item.href ? (
                                    <a onClick={() => navigate(item.href)}>{item.title}</a>
                                ) : (
                                    item.title
                                ),
                            }))}
                            className={styles.breadcrumb}
                        />
                    </div>

                    <div className={styles.headerRight}>
                        <Space size={16}>
                            {/* 主题切换 */}
                            <ThemeToggle />

                            {/* 通知 */}
                            <Badge count={5} size="small">
                                <Button type="text" icon={<BellOutlined />} />
                            </Badge>

                            {/* 用户信息 */}
                            <Dropdown
                                menu={{ items: userMenuItems, onClick: handleUserMenuClick }}
                                placement="bottomRight"
                            >
                                <Space className={styles.userInfo}>
                                    <Avatar
                                        size="small"
                                        icon={<UserOutlined />}
                                        src={userInfo.avatar || undefined}
                                    />
                                    <span className={styles.username} style={{ display: screens.md ? 'inline' : 'none' }}>
                                        {userInfo.username}
                                    </span>
                                </Space>
                            </Dropdown>
                        </Space>
                    </div>
                </Header>

                {/* 内容区 */}
                <Content
                    className={styles.content}
                    style={{
                        marginLeft: getMarginLeft() + (screens.md ? 24 : 0), // 桌面端保持原有间距逻辑
                        marginRight: screens.md ? 24 : 0,
                        marginTop: 24,
                        marginBottom: 24,
                        padding: screens.md ? 24 : 16, // 移动端减少内边距
                        transition: 'all 0.2s',
                        background: token.colorBgContainer,
                    }}
                >
                    <Outlet />
                </Content>
            </Layout>
        </Layout>
    );
};

export default AdminLayout;
