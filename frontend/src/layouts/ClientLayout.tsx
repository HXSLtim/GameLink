import { Outlet, Link, useNavigate } from 'react-router-dom';
import { Button, Space, Dropdown, Avatar, Badge, Popover, List, Typography, Empty, message } from 'antd';
import { UserOutlined, LogoutOutlined, DashboardOutlined, BellOutlined } from '@ant-design/icons';
import logo from '@/assets/logo.svg';
import { userApi, type Notification, type ApiResponse } from '@/api/user';
import { ThemeToggle } from '@/components';
import { useState, useEffect } from 'react';

const { Text } = Typography;

const ClientLayout = () => {
    const navigate = useNavigate();
    const userStr = localStorage.getItem('user_info');
    const user = userStr ? JSON.parse(userStr) : null;

    const [notifications, setNotifications] = useState<Notification[]>([]);
    const [unreadCount, setUnreadCount] = useState(0);
    const [loadingNotifications, setLoadingNotifications] = useState(false);
    const [notificationOpen, setNotificationOpen] = useState(false);

    const fetchNotifications = async () => {
        if (!user) return;
        setLoadingNotifications(true);
        try {
            const res = await userApi.getNotifications({ page: 1, page_size: 5 }) as unknown as ApiResponse<Notification[]>;
            if (res.success) {
                setNotifications(res.data || []);
            }
            const countRes = await userApi.getUnreadCount() as unknown as ApiResponse<{ count: number }>;
            if (countRes.success) {
                setUnreadCount(countRes.data.count);
            }
        } catch (error) {
            console.error('Failed to fetch notifications', error);
        } finally {
            setLoadingNotifications(false);
        }
    };

    useEffect(() => {
        if (user) {
            fetchNotifications();
            // Poll every 60 seconds
            const interval = setInterval(fetchNotifications, 60000);
            return () => clearInterval(interval);
        }
    }, [user?.id]); // Depend on user ID to refetch if user changes

    const handleNotificationClick = async (item: Notification) => {
        if (!item.isRead) {
            try {
                await userApi.markAsRead(item.id);
                setUnreadCount(prev => Math.max(0, prev - 1));
                setNotifications(prev => prev.map(n => n.id === item.id ? { ...n, isRead: true } : n));
            } catch (error) {
                console.error('Failed to mark as read', error);
            }
        }
        // Navigate or show detail if needed
    };

    const handleMarkAllRead = async () => {
        try {
            await userApi.markAllAsRead();
            setUnreadCount(0);
            setNotifications(prev => prev.map(n => ({ ...n, isRead: true })));
            message.success('All marked as read');
        } catch (error) {
            message.error('Failed to mark all as read');
        }
    };

    const notificationContent = (
        <div style={{ width: 300 }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8, padding: '0 8px' }}>
                <Text strong>Notifications</Text>
                {unreadCount > 0 && (
                    <Button type="link" size="small" onClick={handleMarkAllRead}>
                        Mark all as read
                    </Button>
                )}
            </div>
            <List
                loading={loadingNotifications}
                dataSource={notifications}
                locale={{ emptyText: <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="No notifications" /> }}
                renderItem={(item) => (
                    <List.Item
                        onClick={() => handleNotificationClick(item)}
                        style={{
                            cursor: 'pointer',
                            background: item.isRead ? 'transparent' : 'var(--bg-secondary)',
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
                                    <Text strong={!item.isRead} style={{ color: 'var(--text-normal)' }}>{item.title}</Text>
                                </Space>
                            }
                            description={
                                <div>
                                    <div style={{ color: 'var(--text-muted)', fontSize: '12px' }}>{item.content}</div>
                                    <div style={{ color: 'var(--text-muted)', fontSize: '10px', marginTop: 4 }}>
                                        {new Date(item.createdAt).toLocaleString()}
                                    </div>
                                </div>
                            }
                        />
                    </List.Item>
                )}
                style={{ maxHeight: 400, overflowY: 'auto' }}
            />
            <div style={{ textAlign: 'center', marginTop: 8, borderTop: '1px solid var(--background-modifier-accent)', paddingTop: 8 }}>
                <Button type="link" onClick={() => navigate('/notifications')}>View All</Button>
            </div>
        </div>
    );

    const handleLogout = () => {
        localStorage.removeItem('token');
        localStorage.removeItem('user_role');
        localStorage.removeItem('user_info');
        navigate('/login');
    };

    const userMenu = [
        {
            key: 'dashboard',
            label: 'Dashboard',
            icon: <DashboardOutlined />,
            onClick: () => navigate(user?.role === 'ADMIN' ? '/admin' : '/companion'),
        },
        {
            key: 'logout',
            label: 'Logout',
            icon: <LogoutOutlined />,
            onClick: handleLogout,
        },
    ];

    return (
        <div style={{ minHeight: '100vh', display: 'flex', flexDirection: 'column', background: 'var(--bg-primary)', color: 'var(--text-normal)' }}>
            <header style={{
                padding: '0 40px',
                height: '80px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                background: 'var(--bg-header, var(--bg-secondary))', // Fallback or specific variable
                backdropFilter: 'blur(10px)',
                position: 'fixed',
                top: 0,
                left: 0,
                right: 0,
                zIndex: 1000,
                borderBottom: '1px solid var(--background-modifier-active)'
            }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '12px', cursor: 'pointer' }} onClick={() => navigate('/')}>
                    <img src={logo} alt="GameLink" style={{ height: '40px' }} />
                    <span style={{ fontSize: '24px', fontWeight: 'bold', background: 'linear-gradient(90deg, #5865F2, #eb2f96)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>
                        GameLink
                    </span>
                </div>

                <nav style={{ display: 'flex', gap: '40px' }}>
                    <Link to="/" style={{ color: 'var(--text-normal)', fontSize: '16px', fontWeight: 500 }}>Home</Link>
                    <Link to="/companions" style={{ color: 'var(--text-muted)', fontSize: '16px', fontWeight: 500 }}>Companions</Link>
                    <Link to="/games" style={{ color: 'var(--text-muted)', fontSize: '16px', fontWeight: 500 }}>Games</Link>
                    <Link to="/about" style={{ color: 'var(--text-muted)', fontSize: '16px', fontWeight: 500 }}>About</Link>
                </nav>

                <div style={{ display: 'flex', alignItems: 'center', gap: '20px' }}>
                    <ThemeToggle />
                    {user && (
                        <Popover
                            content={notificationContent}
                            trigger="click"
                            open={notificationOpen}
                            onOpenChange={setNotificationOpen}
                            placement="bottomRight"
                            overlayInnerStyle={{ padding: 0 }}
                        >
                            <Badge count={unreadCount} size="small" offset={[-5, 5]}>
                                <Button
                                    type="text"
                                    icon={<BellOutlined style={{ fontSize: '20px', color: 'var(--text-normal)' }} />}
                                    style={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}
                                />
                            </Badge>
                        </Popover>
                    )}
                    {user ? (
                        <Dropdown menu={{ items: userMenu }} placement="bottomRight">
                            <Space style={{ cursor: 'pointer' }}>
                                <Avatar src={user.avatar} icon={<UserOutlined />} style={{ backgroundColor: '#5865F2' }} />
                                <span style={{ fontWeight: 500 }}>{user.name || user.username}</span>
                            </Space>
                        </Dropdown>
                    ) : (
                        <Space>
                            <Link to="/login">
                                <Button type="text" style={{ color: 'var(--text-normal)' }}>Login</Button>
                            </Link>
                            <Link to="/register">
                                <Button type="primary" style={{ backgroundColor: '#5865F2', borderColor: '#5865F2' }}>
                                    Sign Up
                                </Button>
                            </Link>
                        </Space>
                    )}
                </div>
            </header>

            <main style={{ flex: 1, marginTop: '80px' }}>
                <Outlet />
            </main>

            <footer style={{ padding: '40px', background: 'var(--bg-tertiary)', textAlign: 'center', color: 'var(--text-muted)' }}>
                <div style={{ marginBottom: '20px' }}>
                    <span style={{ fontSize: '20px', fontWeight: 'bold', color: 'var(--text-normal)' }}>GameLink</span>
                </div>
                <p>&copy; 2024 GameLink. All rights reserved.</p>
            </footer>
        </div>
    );
};

export default ClientLayout;
