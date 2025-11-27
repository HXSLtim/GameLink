export const ENABLE_QUICK_LOGIN = import.meta.env.DEV; // Only enable in development

export const DEBUG_USERS = [
    {
        label: 'Super Admin',
        email: 'admin@gameLink.com',
        password: '123456',
        role: 'ADMIN',
        color: '#f50'
    },
    {
        label: 'Demo User',
        email: 'demo.user@gamelink.com',
        password: 'User@123456',
        role: 'USER',
        color: '#2db7f5'
    },
    {
        label: 'VIP User',
        email: 'vip.user@gamelink.com',
        password: 'Vip@123456',
        role: 'USER',
        color: '#108ee9'
    },
    {
        label: 'New User',
        email: 'new.user@gamelink.com',
        password: 'User@123789',
        role: 'USER',
        color: '#87d068'
    },
    {
        label: 'Pro Player',
        email: 'pro.player@gamelink.com',
        password: 'Player@123456',
        role: 'COMPANION',
        color: '#722ed1'
    },
    {
        label: 'Streamer',
        email: 'streamer@gamelink.com',
        password: 'Player@654321',
        role: 'COMPANION',
        color: '#eb2f96'
    }
];
