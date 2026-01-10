import { ReactNode } from 'react';
import { useNavigate, useLocation, Outlet } from 'react-router-dom';
import { Tooltip, Dropdown, Avatar } from 'antd';
import {
    HomeOutlined,
    TeamOutlined,
    ShoppingOutlined,
    SettingOutlined,
    PlusOutlined
} from '@ant-design/icons';
import { useAuth } from '@/context/AuthContext';
import ThemeToggle from '@/components/ThemeToggle';
import './style.css';

interface ServerIconProps {
    id: string;
    name: string;
    icon?: ReactNode;
    isActive?: boolean;
    hasNotification?: boolean;
    onClick?: () => void;
}

const ServerIcon = ({ name, icon, isActive, hasNotification, onClick }: ServerIconProps) => {
    return (
        <Tooltip title={name} placement="right">
            <div
                className={`server-icon ${isActive ? 'active' : ''} ${hasNotification ? 'has-notification' : ''}`}
                onClick={onClick}
            >
                <div className="server-icon-inner">
                    {icon || name.charAt(0).toUpperCase()}
                </div>
                <div className="server-pill" />
            </div>
        </Tooltip>
    );
};

const Separator = () => <div className="server-separator" />;

interface ChannelItemProps {
    name: string;
    icon?: ReactNode;
    isActive?: boolean;
    path: string;
    onClick?: () => void;
}

const ChannelItem = ({ name, icon, isActive, onClick }: ChannelItemProps) => {
    return (
        <div
            className={`channel-item ${isActive ? 'active' : ''}`}
            onClick={onClick}
        >
            <span className="channel-icon">{icon || '#'}</span>
            <span className="channel-name">{name}</span>
        </div>
    );
};

export default function DiscordLayout() {
    const navigate = useNavigate();
    const location = useLocation();
    const { user, isAuthenticated, logout } = useAuth();

    // Server navigation items
    const servers = [
        { id: 'home', name: '首页', icon: <HomeOutlined />, path: '/' },
    ];

    // Channel navigation items
    const channels = [
        {
            category: '陪玩服务',
            items: [
                { id: 'players', name: '找陪玩', icon: <TeamOutlined />, path: '/players' },
                { id: 'orders', name: '我的订单', icon: <ShoppingOutlined />, path: '/orders' },
            ]
        },
        {
            category: '个人中心',
            items: [
                { id: 'profile', name: '个人资料', icon: <SettingOutlined />, path: '/profile' },
            ]
        }
    ];

    const handleServerClick = (path: string) => {
        navigate(path);
    };

    const handleChannelClick = (path: string) => {
        navigate(path);
    };

    const handleLogout = () => {
        logout();
        navigate('/login');
    };

    const userMenuItems = [
        { key: 'profile', label: '个人中心', onClick: () => navigate('/profile') },
        { type: 'divider' as const },
        { key: 'logout', label: '退出登录', danger: true, onClick: handleLogout },
    ];

    const isPathActive = (path: string) => {
        if (path === '/') {
            return location.pathname === '/';
        }
        return location.pathname.startsWith(path);
    };

    return (
        <div className="discord-layout">
            {/* Server Sidebar */}
            <nav className="server-sidebar">
                <div className="server-list">
                    {servers.map((server) => (
                        <ServerIcon
                            key={server.id}
                            id={server.id}
                            name={server.name}
                            icon={server.icon}
                            isActive={isPathActive(server.path)}
                            onClick={() => handleServerClick(server.path)}
                        />
                    ))}

                    <Separator />

                    <ServerIcon
                        id="add"
                        name="添加服务器"
                        icon={<PlusOutlined />}
                        onClick={() => { }}
                    />
                </div>

                <ThemeToggle />
            </nav>

            {/* Channel Sidebar */}
            <aside className="channel-sidebar">
                <header className="channel-header">
                    <h2 className="server-name">GameLink</h2>
                </header>

                <div className="channel-list-wrapper">
                    <div className="channel-list">
                        {channels.map((category, idx) => (
                            <div key={idx} className="channel-category">
                                <div className="category-header">
                                    <span className="category-name">{category.category}</span>
                                </div>
                                {category.items.map((channel) => (
                                    <ChannelItem
                                        key={channel.id}
                                        name={channel.name}
                                        icon={channel.icon}
                                        isActive={isPathActive(channel.path)}
                                        path={channel.path}
                                        onClick={() => handleChannelClick(channel.path)}
                                    />
                                ))}
                            </div>
                        ))}
                    </div>
                </div>

                {/* User Panel */}
                <div className="user-panel">
                    <Dropdown
                        menu={{ items: userMenuItems }}
                        trigger={['click']}
                        placement="topLeft"
                    >
                        <div className="user-panel-inner">
                            <Avatar
                                size={32}
                                src={user?.avatar}
                                className="user-avatar"
                            >
                                {user?.nickname?.charAt(0) || 'U'}
                            </Avatar>
                            <div className="user-info">
                                <div className="user-name">{isAuthenticated ? (user?.nickname || '用户') : '未登录'}</div>
                                <div className="user-status">
                                    {isAuthenticated ? '在线' : '点击登录'}
                                </div>
                            </div>
                            {!isAuthenticated && (
                                <button
                                    className="login-btn"
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        navigate('/login');
                                    }}
                                >
                                    登录
                                </button>
                            )}
                        </div>
                    </Dropdown>
                </div>
            </aside>

            {/* Main Content Area */}
            <main className="main-content">
                <header className="content-header">
                    <div className="header-title">
                        <span className="header-icon">#</span>
                        <span className="header-name">
                            {channels.flatMap(c => c.items).find(ch => isPathActive(ch.path))?.name || '首页'}
                        </span>
                    </div>
                </header>

                <div className="content-body">
                    <Outlet />
                </div>
            </main>
        </div>
    );
}
