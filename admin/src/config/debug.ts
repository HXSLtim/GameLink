export const ENABLE_QUICK_LOGIN = import.meta.env.DEV; // Only enable in development

/**
 * 调试用户列表
 * 角色与后端 model.Role 保持一致：user, player, admin
 * 这里使用大写是为了与路由守卫中的 Role 类型兼容
 */
export const DEBUG_USERS = [
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
];
