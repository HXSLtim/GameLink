export const ENABLE_QUICK_LOGIN = import.meta.env.DEV; // Only enable in development

/**
 * 调试用户列表
 * ⚠️ 安全警告: 此配置仅用于开发环境
 * 生产构建时会通过环境变量禁用
 */
export const DEBUG_USERS = import.meta.env.DEV ? [
    {
        label: '超级管理员',
        email: 'admin@gameLink.com',
        password: '123456',
        role: 'ADMIN',
        color: '#f50'
    },
    {
        label: '测试用户',
        email: 'demo.user@gamelink.com',
        password: 'User@123456',
        role: 'USER',
        color: '#2db7f5'
    },
    {
        label: '高级会员',
        email: 'vip.user@gamelink.com',
        password: 'Vip@123456',
        role: 'USER',
        color: '#108ee9'
    },
    {
        label: '职业陪玩',
        email: 'pro.player@gamelink.com',
        password: 'Player@123456',
        role: 'PLAYER',
        color: '#722ed1'
    },
    {
        label: '魔王主播',
        email: 'streamer@gamelink.com',
        password: 'Player@654321',
        role: 'PLAYER',
        color: '#eb2f96'
    }
] : [];
