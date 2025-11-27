import React, { useState } from 'react';
import { Outlet, useLocation, useNavigate, Link } from 'react-router-dom';
import { ConfigProvider, Layout, Menu, theme, Avatar, Tooltip, Breadcrumb, Button, Drawer, Radio, ColorPicker, Divider, Space } from 'antd';
import {
    DashboardOutlined,
    UserOutlined,
    TrophyOutlined,
    ShoppingOutlined,
    FileTextOutlined,
    SettingOutlined,
    HomeFilled,
    SafetyCertificateFilled,
    MenuFoldOutlined,
    MenuUnfoldOutlined,
    BgColorsOutlined,
    LayoutOutlined
} from '@ant-design/icons';
import { motion, AnimatePresence } from 'framer-motion';
import zhCN from 'antd/locale/zh_CN';
import { useAdmin } from '@/context/AdminContext';
import { getIcon } from '@/utils/iconMap';

const { Sider, Content, Header } = Layout;

const AdminLayout: React.FC = () => {
    const location = useLocation();
    const navigate = useNavigate();

    // Theme & Layout State
    const [darkMode, setDarkMode] = useState(true);
    const [primaryColor, setPrimaryColor] = useState('#5865F2');
    const [collapsed, setCollapsed] = useState(false);
    const [settingsVisible, setSettingsVisible] = useState(false);
    const [contentWidth, setContentWidth] = useState<'fluid' | 'fixed'>('fluid');

    // Breadcrumb Mapping
    const breadcrumbNameMap: Record<string, string> = {
        '/admin': '仪表盘',
        '/admin/users': '用户管理',
        '/admin/games': '游戏管理',
        '/admin/orders': '订单管理',
        '/admin/audit': '审计日志',
        '/admin/settings': '系统设置',
    };

    const pathSnippets = location.pathname.split('/').filter(i => i);
    const breadcrumbItems = [
        { title: <Link to="/admin">首页</Link> },
        ...pathSnippets.map((_, index) => {
            const url = `/${pathSnippets.slice(0, index + 1).join('/')}`;
            return {
                title: breadcrumbNameMap[url] || url,
            };
        }),
    ];

    // Dynamic Theme Algorithm
    const getThemeAlgorithm = () => {
        return darkMode ? theme.darkAlgorithm : theme.defaultAlgorithm;
    };

    const customTheme = {
        algorithm: getThemeAlgorithm(),
        token: {
            colorPrimary: primaryColor,
            borderRadius: 8,
            ...(darkMode ? {
                colorBgBase: '#36393f',
                colorBgContainer: '#2f3136',
                colorBgElevated: '#18191c',
                colorText: '#dcddde',
                colorTextSecondary: '#b9bbbe',
            } : {})
        },
        components: {
            Menu: {
                itemBg: 'transparent',
            },
            Layout: {
                ...(darkMode ? {
                    bodyBg: '#36393f',
                    siderBg: '#2f3136',
                    headerBg: '#36393f',
                } : {
                    bodyBg: '#f0f2f5',
                    siderBg: '#fff',
                    headerBg: '#fff',
                })
            }
        }
    };



    const { menus } = useAdmin();

    const staticMenuItems = [
        { key: '/admin/sys/dashboard', icon: <DashboardOutlined />, label: '仪表盘' },
        { key: '/admin/sys/user', icon: <UserOutlined />, label: '用户管理' },
        { key: '/admin/biz/game', icon: <TrophyOutlined />, label: '游戏管理' },
        { key: '/admin/biz/service', icon: <ShoppingOutlined />, label: '服务项目' },
        { key: '/admin/biz/order', icon: <FileTextOutlined />, label: '订单管理' },
        { key: '/admin/sys/log', icon: <SafetyCertificateFilled />, label: '审计日志' },
        { key: '/admin/sys/setting', icon: <SettingOutlined />, label: '系统设置' },
    ];

    const generateMenuItems = (items: any[]): any[] => {
        return items.map(item => {
            if (item.hidden) return null;
            const path = item.path.startsWith('/') ? item.path.substring(1) : item.path;
            const fullPath = `/admin/${path}`;

            return {
                key: fullPath,
                icon: getIcon(item.icon) || <FileTextOutlined />, // Fallback icon
                label: item.name,
                children: item.children && item.children.length > 0 ? generateMenuItems(item.children) : undefined
            };
        }).filter(Boolean);
    };

    const menuItems = menus.length > 0 ? generateMenuItems(menus) : staticMenuItems;

    return (
        <ConfigProvider theme={customTheme} locale={zhCN}>
            <Layout style={{ height: '100vh', overflow: 'hidden' }}>
                {/* Primary Server Sidebar (Discord Style - Always Dark) */}
                <Sider
                    width={72}
                    style={{
                        backgroundColor: '#202225',
                        borderRight: 'none',
                        display: 'flex',
                        flexDirection: 'column',
                        alignItems: 'center',
                        paddingTop: 12,
                        zIndex: 100
                    }}
                >
                    <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 8, width: '100%' }}>
                        <Tooltip title="GameLink 首页" placement="right">
                            <div
                                style={{
                                    width: 48, height: 48,
                                    display: 'flex', alignItems: 'center', justifyContent: 'center',
                                    cursor: 'pointer',
                                    transition: 'all 0.3s ease'
                                }}
                            >
                                <img src="/logo.svg" alt="GameLink" style={{ width: '100%', height: '100%', borderRadius: 16 }} />
                            </div>
                        </Tooltip>

                        <div style={{ width: 32, height: 2, background: '#4f545c', borderRadius: 1 }} />

                        <Tooltip title="管理控制台" placement="right">
                            <div style={{
                                width: 48, height: 48, borderRadius: 16,
                                background: primaryColor, color: '#fff',
                                display: 'flex', alignItems: 'center', justifyContent: 'center',
                                fontSize: 24, cursor: 'pointer'
                            }}>
                                <SafetyCertificateFilled />
                            </div>
                        </Tooltip>

                        <Tooltip title="客户端首页" placement="right">
                            <div
                                onClick={() => navigate('/')}
                                style={{
                                    width: 48, height: 48, borderRadius: 24,
                                    background: '#36393f', color: '#3ba55c',
                                    display: 'flex', alignItems: 'center', justifyContent: 'center',
                                    fontSize: 24, cursor: 'pointer',
                                    transition: 'all 0.2s'
                                }}
                                onMouseEnter={(e) => e.currentTarget.style.borderRadius = '16px'}
                                onMouseLeave={(e) => e.currentTarget.style.borderRadius = '24px'}
                            >
                                <HomeFilled />
                            </div>
                        </Tooltip>
                    </div>
                </Sider>

                {/* Secondary Navigation Sidebar */}
                <Sider
                    width={240}
                    collapsible
                    collapsed={collapsed}
                    trigger={null}
                    style={{
                        backgroundColor: darkMode ? '#2f3136' : '#fff',
                        borderRight: darkMode ? 'none' : '1px solid #f0f0f0'
                    }}
                >
                    <div style={{
                        height: 48, padding: '0 16px', display: 'flex', alignItems: 'center',
                        boxShadow: '0 1px 0 rgba(4,4,5,0.05)', fontWeight: 600,
                        color: darkMode ? '#fff' : '#000',
                        justifyContent: collapsed ? 'center' : 'flex-start'
                    }}>
                        {!collapsed && 'GameLink 管理后台'}
                        {collapsed && <img src="/logo.svg" alt="GL" style={{ width: 32, height: 32, borderRadius: 8 }} />}
                    </div>
                    <div style={{ padding: '16px 8px' }}>
                        <Menu
                            mode="inline"
                            selectedKeys={[location.pathname]}
                            items={menuItems}
                            onClick={({ key }) => navigate(key)}
                            style={{ borderRight: 'none', background: 'transparent' }}
                        />
                    </div>

                    {/* User Profile Footer */}
                    <div style={{
                        position: 'absolute', bottom: 0, width: '100%',
                        height: 52, background: darkMode ? '#292b2f' : '#f9f9f9',
                        display: 'flex', alignItems: 'center', padding: '0 8px',
                        borderTop: darkMode ? 'none' : '1px solid #f0f0f0'
                    }}>
                        <Avatar style={{ backgroundColor: primaryColor }} size="small">A</Avatar>
                        {!collapsed && (
                            <div style={{ marginLeft: 8 }}>
                                <div style={{ fontSize: 12, fontWeight: 600, color: darkMode ? '#fff' : '#000' }}>管理员</div>
                                <div style={{ fontSize: 10, color: darkMode ? '#b9bbbe' : '#666' }}>#0001</div>
                            </div>
                        )}
                    </div>
                </Sider>

                {/* Main Content */}
                <Layout style={{ backgroundColor: darkMode ? '#36393f' : '#f0f2f5' }}>
                    <Header style={{
                        padding: '0 16px',
                        background: darkMode ? '#36393f' : '#fff',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'space-between',
                        boxShadow: '0 1px 2px rgba(0,0,0,0.03)',
                        height: 48,
                        zIndex: 99
                    }}>
                        <div style={{ display: 'flex', alignItems: 'center' }}>
                            <Button
                                type="text"
                                icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
                                onClick={() => setCollapsed(!collapsed)}
                                style={{ fontSize: '16px', width: 48, height: 48, color: darkMode ? '#fff' : '#000' }}
                            />
                            <Breadcrumb items={breadcrumbItems} style={{ marginLeft: 16 }} />
                        </div>

                        <Space>
                            <Tooltip title="主题设置">
                                <Button
                                    type="text"
                                    icon={<BgColorsOutlined />}
                                    onClick={() => setSettingsVisible(true)}
                                    style={{ color: darkMode ? '#fff' : '#000' }}
                                />
                            </Tooltip>
                        </Space>
                    </Header>

                    <Content style={{
                        padding: 24,
                        overflowY: 'auto',
                        margin: contentWidth === 'fixed' ? '0 auto' : 0,
                        width: contentWidth === 'fixed' ? '1200px' : '100%',
                        maxWidth: '100%'
                    }}>
                        <AnimatePresence mode="wait">
                            <motion.div
                                key={location.pathname}
                                initial={{ opacity: 0, y: 20 }}
                                animate={{ opacity: 1, y: 0 }}
                                exit={{ opacity: 0, y: -20 }}
                                transition={{ duration: 0.2 }}
                            >
                                <Outlet />
                            </motion.div>
                        </AnimatePresence>
                    </Content>
                </Layout>

                {/* Theme & Layout Settings Drawer */}
                <Drawer
                    title="布局设置"
                    placement="right"
                    onClose={() => setSettingsVisible(false)}
                    open={settingsVisible}
                    width={300}
                >
                    <Divider><BgColorsOutlined /> 主题风格</Divider>
                    <div style={{ marginBottom: 24, textAlign: 'center' }}>
                        <Radio.Group value={darkMode} onChange={e => setDarkMode(e.target.value)}>
                            <Radio.Button value={false}>亮色</Radio.Button>
                            <Radio.Button value={true}>暗色</Radio.Button>
                        </Radio.Group>
                    </div>

                    <Divider><BgColorsOutlined /> 主题色</Divider>
                    <div style={{ marginBottom: 24, textAlign: 'center' }}>
                        <ColorPicker
                            value={primaryColor}
                            onChange={(color) => setPrimaryColor(color.toHexString())}
                            showText
                        />
                    </div>

                    <Divider><LayoutOutlined /> 内容区域</Divider>
                    <div style={{ marginBottom: 24, textAlign: 'center' }}>
                        <Radio.Group value={contentWidth} onChange={e => setContentWidth(e.target.value)}>
                            <Radio.Button value="fluid">流式</Radio.Button>
                            <Radio.Button value="fixed">定宽</Radio.Button>
                        </Radio.Group>
                    </div>
                </Drawer>
            </Layout>
        </ConfigProvider>
    );
};

export default AdminLayout;
