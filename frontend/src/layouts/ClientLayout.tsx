import { Outlet, Link, useNavigate } from 'react-router-dom';
import { Button, Space, Dropdown, Avatar } from 'antd';
import { UserOutlined, LogoutOutlined, DashboardOutlined } from '@ant-design/icons';
import logo from '@/assets/logo.svg';

import { ThemeToggle } from '@/components';

const ClientLayout = () => {
    const navigate = useNavigate();
    const userStr = localStorage.getItem('user_info');
    const user = userStr ? JSON.parse(userStr) : null;

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
