/**
 * 启用快速登录
 * ⚠️ 安全警告: 生产环境必须设置为 false
 */
export const ENABLE_QUICK_LOGIN = import.meta.env.VITE_DEBUG_USERS !== 'false' && import.meta.env.DEV;

export type DebugUser = {
    label: string;
    email: string;
    password: string;
    role: string;
    color: string;
};

/**
 * 解析调试用户配置
 * 从环境变量 JSON 字符串或使用默认值
 */
function parseDebugUsers(): DebugUser[] {
    if (import.meta.env.PROD) {
        // 生产环境禁用调试用户，避免运行时崩溃
        return [];
    }
    
    // 优先使用环境变量
    const envUsers = import.meta.env.VITE_DEBUG_USERS;
    if (envUsers) {
        try {
            return JSON.parse(envUsers) as DebugUser[];
        } catch (e) {
            console.error('Failed to parse VITE_DEBUG_USERS:', e);
            return [];
        }
    }
    
    // 默认仅用于本地开发
    return import.meta.env.DEV ? [
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
        }
    ] : [];
}

/**
 * 调试用户列表
 * ⚠️ 安全警告:
 * 1. 仅用于开发环境
 * 2. 生产构建会自动禁用
 * 3. 使用 .gitignore 防止提交
 * 4. 建议使用强密码
 */
export const DEBUG_USERS = parseDebugUsers();
