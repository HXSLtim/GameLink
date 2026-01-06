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
    Popover,
    List,
    Typography,
    Empty,
    message,
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
    SettingFilled,
    TeamOutlined,
    SafetyCertificateOutlined,
    MenuOutlined,
    AppstoreOutlined,
    ShoppingCartOutlined,
    MonitorOutlined,
    LineChartOutlined,
    FundOutlined,
    FileTextOutlined,
    StarOutlined,
    UnorderedListOutlined,
    AuditOutlined,
    WarningOutlined,
    StopOutlined,
    BarChartOutlined,
    GiftOutlined,
    MessageOutlined,
    TagsOutlined,
    DollarOutlined,
    WalletOutlined,
    PercentageOutlined,
    BranchesOutlined,
    BankOutlined,
    TrophyOutlined,
    SwapOutlined,
    CrownOutlined,
    PayCircleOutlined,
    ShareAltOutlined,
    CalendarOutlined,
    TransactionOutlined,
} from '@ant-design/icons';
import { useAdmin } from '@/context/useAdmin';
import { useTheme } from '@/context/useTheme';
import { authApi } from '@/api/auth';
import { adminApi, type Menu as BackendMenuItem } from '@/api/admin';
import { userApi, type Notification, type ApiResponse, type NotificationListResponse } from '@/api/user';
import { ThemeToggle } from '@/components';
import { forceInit } from '@/services/init';
import styles from './index.module.css';

import { logger } from '@/utils/logger';
const { Header, Sider, Content } = Layout;
const { Text } = Typography;

type MenuItem = Required<MenuProps>['items'][number];

/**
 * 图标映射 - 将图标名称字符串映射到实际的图标组件
 */
const iconMap: Record<string, React.ComponentType> = {
    // 通用图标（完整名称）
    DashboardOutlined,
    SettingFilled,
    SettingOutlined,
    UserOutlined,
    TeamOutlined,
    SafetyCertificateOutlined,
    MenuOutlined,
    // 业务管理
    AppstoreOutlined,
    ShoppingCartOutlined,
    GiftOutlined,
    BranchesOutlined,
    // 财务管理
    DollarOutlined,
    WalletOutlined,
    PercentageOutlined,
    BankOutlined,
    TrophyOutlined,
    SwapOutlined,
    // 监控中心
    MonitorOutlined,
    LineChartOutlined,
    FundOutlined,
    BarChartOutlined,
    // 内容管理
    FileTextOutlined,
    MessageOutlined,
    TagsOutlined,
    // 评价管理
    StarOutlined,
    UnorderedListOutlined,
    AuditOutlined,
    WarningOutlined,
    StopOutlined,
    // 营销管理
    CrownOutlined,
    PayCircleOutlined,
    ShareAltOutlined,
    CalendarOutlined,
    TransactionOutlined,

    // 简写名称映射（数据库中存储的名称）
    'dashboard': DashboardOutlined,
    'setting': SettingOutlined,
    'user': UserOutlined,
    'team': TeamOutlined,
    'safety-certificate': SafetyCertificateOutlined,
    'menu': MenuOutlined,
    'appstore': AppstoreOutlined,
    'shopping-cart': ShoppingCartOutlined,
    'gift': GiftOutlined,
    'dollar': DollarOutlined,
    'wallet': WalletOutlined,
    'percentage': PercentageOutlined,
    'monitor': MonitorOutlined,
    'line-chart': LineChartOutlined,
    'fund': FundOutlined,
    'bar-chart': BarChartOutlined,
    'file-text': FileTextOutlined,
    'star': StarOutlined,
    'tags': TagsOutlined,
    'warning': WarningOutlined,
    'branches': BranchesOutlined,
    'bank': BankOutlined,
    'trophy': TrophyOutlined,
    'swap': SwapOutlined,
    'crown': CrownOutlined,
    'pay-circle': PayCircleOutlined,
    'share-alt': ShareAltOutlined,
    'calendar': CalendarOutlined,
    'transaction': TransactionOutlined,
};

const AdminLayout: React.FC = () => {
    const [collapsed, setCollapsed] = useState(false);
    const navigate = useNavigate();
    const location = useLocation();
    // 使用 AdminContext 中已过滤的菜单数据
    // Requirements: 8.1, 8.2, 8.4 - 菜单权限联动
    const { loading, menus: contextMenus, permissionVersion } = useAdmin();
    const { token } = theme.useToken();
    const screens = Grid.useBreakpoint();
    const { mode } = useTheme();
    const [messageApi, contextHolder] = message.useMessage();

    // 菜单数据 - 现在从 AdminContext 获取已过滤的菜单
    const [menuData, setMenuData] = useState<BackendMenuItem[]>([]);
    const [menuLoading, setMenuLoading] = useState(false); // 由 AdminContext.loading 掟一控制

    // 初始化状态
    const [initializing, setInitializing] = useState(false);

    // 数据是否已加载过（用于区分首次加载和数据为空）
    const [dataLoaded, setDataLoaded] = useState(false);

    // 用户信息
    const [userInfo, setUserInfo] = useState<{ username: string; avatar?: string; id?: number }>({
        username: 'Admin',
    });

    // 通知状态
    const [notifications, setNotifications] = useState<Notification[]>([]);
    const [unreadCount, setUnreadCount] = useState(0);
    const [loadingNotifications, setLoadingNotifications] = useState(false);
    const [notificationOpen, setNotificationOpen] = useState(false);

    const fetchNotifications = async () => {
        setLoadingNotifications(true);
        try {
            const response = await userApi.getNotifications({ page: 1, page_size: 5 });
            // Axios 返回 AxiosResponse，需要访问 response.data 获取 ApiResponse
            const res = response.data as ApiResponse<NotificationListResponse>;
            if (res.success && res.data) {
                setNotifications(res.data.items || []);
                if (res.data.unreadCount !== undefined) {
                    setUnreadCount(res.data.unreadCount);
                }
            }
        } catch (error) {
            logger.error('Failed to fetch notifications', error);
        } finally {
            setLoadingNotifications(false);
        }
    };

    useEffect(() => {
        fetchNotifications();
        // Poll every 60 seconds
        const interval = setInterval(fetchNotifications, 60000);
        return () => clearInterval(interval);
    }, []);

    const handleNotificationClick = async (item: Notification) => {
        if (!item.isRead) {
            try {
                await userApi.markAsRead(item.id);
                setUnreadCount(prev => Math.max(0, prev - 1));
                setNotifications(prev => prev.map(n => n.id === item.id ? { ...n, isRead: true } : n));
            } catch (error) {
                logger.error('Failed to mark as read', error);
            }
        }
    };

    const handleMarkAllRead = async () => {
        try {
            await userApi.markAllAsRead();
            setUnreadCount(0);
            setNotifications(prev => prev.map(n => ({ ...n, isRead: true })));
            messageApi.success('已全部标记为已读');
        } catch {
            messageApi.error('全部标记已读失败');
        }
    };

    const notificationContent = (
        <div style={{ width: 300 }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8, padding: '0 8px' }}>
                <Text strong>通知</Text>
                {unreadCount > 0 && (
                    <Button type="link" size="small" onClick={handleMarkAllRead}>
                        全部已读
                    </Button>
                )}
            </div>
            <List
                loading={loadingNotifications}
                dataSource={notifications}
                locale={{ emptyText: <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="暂无通知" /> }}
                renderItem={(item) => (
                    <List.Item
                        onClick={() => handleNotificationClick(item)}
                        style={{
                            cursor: 'pointer',
                            background: item.isRead ? 'transparent' : token.colorBgContainer, // Use token for bg
                            padding: '8px 12px',
                            borderRadius: '4px',
                            marginBottom: '4px',
                            transition: 'background 0.3s'
                        }}
                    >
                        <List.Item.Meta
                            title={
                                <Space>
                                    {!item.isRead && <Badge status="processing" />}
                                    <Text strong={!item.isRead}>{item.title}</Text>
                                </Space>
                            }
                            description={
                                <div>
                                    <div style={{ fontSize: '12px', color: token.colorTextSecondary }}>{item.message}</div>
                                    <div style={{ fontSize: '10px', marginTop: 4, color: token.colorTextTertiary }}>
                                        {new Date(item.createdAt).toLocaleString()}
                                    </div>
                                </div>
                            }
                        />
                    </List.Item>
                )}
                style={{ maxHeight: 400, overflowY: 'auto' }}
            />
        </div>
    );

    // 加载用户信息
    useEffect(() => {
        const loadUserInfo = async () => {
            // 先从 localStorage 获取缓存的用户信息
            const storedUser = localStorage.getItem('user_info');
            if (storedUser) {
                try {
                    const parsed = JSON.parse(storedUser);
                    // 兼容两种字段名格式：username/name 和 avatar/avatarUrl
                    setUserInfo({
                        username: parsed.username || parsed.name || 'Admin',
                        avatar: parsed.avatar || parsed.avatarUrl,
                        id: parsed.id,
                    });
                } catch {
                    // ignore
                }
            }

            // 从 API 获取最新用户信息
            try {
                const { authApi } = await import('@/api/auth');
                const response = await authApi.getMe();
                if (response.data?.success && response.data?.data) {
                    const userData = response.data.data.user || response.data.data;
                    const newUserInfo = {
                        username: userData.name || userData.username || 'Admin',
                        avatar: userData.avatarUrl || userData.avatar,
                        id: userData.id,
                    };
                    setUserInfo(newUserInfo);
                    // 更新 localStorage 缓存
                    localStorage.setItem('user_info', JSON.stringify({
                        ...userData,
                        username: newUserInfo.username,
                        avatar: newUserInfo.avatar,
                    }));
                }
            } catch (error) {
                logger.error('Failed to load user info', error);
            }
        };

        const token = localStorage.getItem('token');
        if (token) {
            loadUserInfo();
        }
    }, []);

    // 从 AdminContext 获取已过滤的菜单数据
    // Requirements: 8.1, 8.2, 8.4 - 菜单权限联动，权限变更后自动更新
    useEffect(() => {
        const loadMenus = async () => {
            try {
                setMenuLoading(true);

                // 优先使用 AdminContext 中已过滤的菜单
                if (contextMenus && contextMenus.length > 0) {
                    logger.info('[AdminLayout] 使用 AdminContext 中的菜单数据，数量:', contextMenus.length);

                    // 后端已经返回树形结构，直接使用
                    logger.info('[AdminLayout] 菜单数据（树形结构）:', contextMenus);
                    setMenuData(contextMenus);
                    setDataLoaded(true); // 标记数据已加载
                    setMenuLoading(false);
                    return;
                }

                // 如果 AdminContext 没有菜单数据，等待 AdminContext 加载完成
                // 不再直接调用 API，完全依赖 AdminContext
                if (loading) {
                    logger.info('[AdminLayout] AdminContext 正在加载中，等待...');
                    return;
                }

                // AdminContext 加载完成但没有菜单数据
                logger.info('[AdminLayout] AdminContext 加载完成但无菜单数据');
                setDataLoaded(true);
                setMenuData([]);
            } catch (error) {
                logger.error('Failed to load menus:', error);
                setDataLoaded(true);
            } finally {
                setMenuLoading(false);
            }
        };

        loadMenus();
    }, [contextMenus, loading, permissionVersion]); // 监听 loading 状态变化

    // 响应式处理：屏幕变窄时自动收起
    useEffect(() => {
        if (!screens.md) {
            setCollapsed(true);
        } else {
            setCollapsed(false);
        }
    }, [screens.md]);

    /**
     * 将后端菜单转换为 Ant Design 菜单格式
     */
    const menuItems = useMemo((): MenuItem[] => {
        if (!menuData || menuData.length === 0) return [];

        const transformMenu = (menus: BackendMenuItem[]): MenuItem[] => {
            return menus.map(menu => {
                const IconComponent = menu.icon ? iconMap[menu.icon] : null;
                const children = menu.children && menu.children.length > 0
                    ? transformMenu(menu.children)
                    : undefined;

                return {
                    key: menu.path,
                    icon: IconComponent ? React.createElement(IconComponent) : null,
                    label: menu.name,
                    children,
                } as MenuItem;
            });
        };

        return transformMenu(menuData);
    }, [menuData]);

    /**
     * 创建面包屑映射
     */
    const breadcrumbMap = useMemo((): Record<string, string> => {
        const map: Record<string, string> = {};

        const traverse = (menus: BackendMenuItem[]) => {
            menus.forEach(menu => {
                map[menu.path] = menu.name;
                if (menu.children && menu.children.length > 0) {
                    traverse(menu.children);
                }
            });
        };

        if (menuData && menuData.length > 0) {
            traverse(menuData);
        }

        return map;
    }, [menuData]);

    /**
     * 获取当前选中的菜单项
     */
    const selectedKeys = useMemo(() => {
        const path = location.pathname;
        const keys = Object.keys(breadcrumbMap).filter(key => path.startsWith(key));
        return keys.length > 0 ? [keys.sort((a, b) => b.length - a.length)[0]] : ['/admin'];
    }, [location.pathname, breadcrumbMap]);

    /**
     * 获取展开的菜单项
     */
    const openKeys = useMemo(() => {
        const path = location.pathname;

        // 根据菜单层级结构找出需要展开的父菜单
        const findParentKeys = (menus: BackendMenuItem[], currentPath: string): string[] => {
            const parentKeys: string[] = [];

            const checkMenu = (menu: BackendMenuItem, parentKey?: string) => {
                if (currentPath.startsWith(menu.path)) {
                    if (parentKey) parentKeys.push(parentKey);
                }
                if (menu.children) {
                    menu.children.forEach(child => checkMenu(child, menu.path));
                }
            };

            menus.forEach(menu => checkMenu(menu));
            return parentKeys;
        };

        return findParentKeys(menuData, path);
    }, [location.pathname, menuData]);

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
    }, [location.pathname, breadcrumbMap]);

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
        navigate('/admin/login');
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
        } else if (key === 'profile') {
            navigate('/admin/profile');
        }
    };

    // 首次加载时显示加载状态
    // 只有数据已加载且为空时才显示初始化按钮
    if (loading || menuLoading) {
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

    // 处理初始化
    const handleInitialize = async () => {
        setInitializing(true);
        try {
            const result = await forceInit();
            if (result.success) {
                messageApi.success('初始化成功，正在刷新...');
                // 重新加载菜单
                setTimeout(() => {
                    window.location.reload();
                }, 1000);
            } else {
                messageApi.error(`初始化失败: ${result.errors.join(', ')}`);
            }
        } catch (error) {
            messageApi.error('初始化失败');
            logger.error('Init error:', error);
        } finally {
            setInitializing(false);
        }
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

            {/* 只有数据已加载且为空时才显示初始化按钮 */}
            {dataLoaded && menuItems.length === 0 ? (
                <div style={{ padding: 16, textAlign: 'center' }}>
                    <Empty
                        image={Empty.PRESENTED_IMAGE_SIMPLE}
                        description="暂无菜单数据"
                    >
                        <Button
                            type="primary"
                            onClick={handleInitialize}
                            loading={initializing}
                        >
                            初始化系统
                        </Button>
                    </Empty>
                </div>
            ) : (
            /* 菜单 */
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
            )}
        </>
    );

    return (
        <Layout className={styles.layout} style={{ minHeight: '100vh' }}>
            {contextHolder}
            {/* 侧边栏 - 移动端使用 Drawer */}
            {!screens.md ? (
                <Drawer
                    placement="left"
                    onClose={() => setCollapsed(true)}
                    open={!collapsed}
                    styles={{ body: { padding: 0, width: 220 }, header: { display: 'none' } }}
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
                            {/* 通知 */}
                            <Popover
                                content={notificationContent}
                                trigger="click"
                                open={notificationOpen}
                                onOpenChange={setNotificationOpen}
                                placement="bottomRight"
                            >
                                <Badge count={unreadCount} size="small">
                                    <Button type="text" icon={<BellOutlined />} />
                                </Badge>
                            </Popover>

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
