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
} from '@ant-design/icons';
import { useAdmin } from '@/context/AdminContext';
import { useTheme } from '@/context/ThemeContext';
import { authApi } from '@/api/auth';
import { adminApi, type Menu as BackendMenuItem } from '@/api/admin';
import { ThemeToggle } from '@/components';
import styles from './index.module.css';

const { Header, Sider, Content } = Layout;

type MenuItem = Required<MenuProps>['items'][number];

/**
 * 图标映射 - 将图标名称字符串映射到实际的图标组件
 */
const iconMap: Record<string, any> = {
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
    SettingOutlined,
    FileTextOutlined,
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

    // 菜单数据
    const [menuData, setMenuData] = useState<BackendMenuItem[]>([]);
    const [menuLoading, setMenuLoading] = useState(true);

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

    // 加载菜单数据
    useEffect(() => {
        const loadMenus = async () => {
            try {
                setMenuLoading(true);
                console.log('[AdminLayout] 开始加载菜单数据...');
                const response = await adminApi.getMenus({ parentId: undefined });
                console.log('[AdminLayout] 菜单API响应:', response);
                console.log('[AdminLayout] response 类型:', typeof response);

                // Axios 拦截器返回的是 response.data，所以 response 就是 ApiResponse
                // response.data 才是实际的菜单数组
                let menus: BackendMenuItem[] = [];
                if (response && Array.isArray(response.data)) {
                    // response.data 是菜单数组
                    menus = response.data;
                    console.log('[AdminLayout] 从 response.data 提取菜单，数量:', menus.length);
                } else if (Array.isArray(response)) {
                    // 如果后端直接返回数组（不应该发生）
                    menus = response;
                    console.log('[AdminLayout] 直接获取菜单数组，数量:', menus.length);
                } else {
                    console.log('[AdminLayout] 未识别的响应格式:', response);
                    menus = [];
                }

                console.log('[AdminLayout] 提取的菜单数据:', menus);

                // 检查 visible 字段（注意后端字段是 visible，不是 hidden）
                menus.forEach(menu => {
                    console.log(`[AdminLayout] 菜单: ${menu.name}, parentId: ${menu.parentId}, visible: ${menu.visible}`);
                });

                // 构建菜单树：将扁平数组转换为层级树结构
                const buildMenuTree = (menus: BackendMenuItem[], parentId: number | null = null): BackendMenuItem[] => {
                    return menus
                        .filter(menu => {
                            // 过滤：1) visible=true 2) parentId匹配
                            const isVisible = menu.visible !== false;
                            const isMatchParent = menu.parentId === parentId;
                            return isVisible && isMatchParent;
                        })
                        .sort((a, b) => a.order - b.order)  // 按order排序
                        .map(menu => {
                            // 递归查找子菜单
                            const children = buildMenuTree(menus, menu.id);
                            return {
                                ...menu,
                                children: children.length > 0 ? children : undefined
                            };
                        });
                };

                const menuTree = buildMenuTree(menus);
                console.log('[AdminLayout] 构建的菜单树:', menuTree);
                console.log('[AdminLayout] 根菜单数量:', menuTree.length);

                setMenuData(menuTree);
            } catch (error) {
                console.error('Failed to load menus:', error);
            } finally {
                setMenuLoading(false);
            }
        };

        loadMenus();
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
