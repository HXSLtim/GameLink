import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Form, Input, Button, Card, message, theme } from 'antd';
import { UserOutlined, LockOutlined, SafetyCertificateOutlined } from '@ant-design/icons';
import { authApi } from '@/api/auth';
import { ENABLE_QUICK_LOGIN, DEBUG_USERS } from '@/config/debug';

/**
 * 管理后台登录页面
 * 独立于用户端登录，专为管理员设计
 */
const AdminLogin: React.FC = () => {
    const { token } = theme.useToken();
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);

    const onLogin = async (values: { username: string; password: string }) => {
        setLoading(true);
        try {
            const res = await authApi.login({
                username: values.username,
                password: values.password
            });

            const response = res.data as { 
                success?: boolean; 
                data?: { token: string; user: { id: number; role: string; [key: string]: unknown } } 
            };
            
            if (!response.success || !response.data) {
                throw new Error('登录响应格式错误');
            }
            
            const { token: authToken, user } = response.data;
            const role = user.role.toUpperCase();

            // 验证是否为管理员角色
            if (role !== 'ADMIN') {
                message.error('您没有管理后台访问权限');
                return;
            }

            localStorage.setItem('token', authToken);
            localStorage.setItem('user_role', user.role);
            localStorage.setItem('user_info', JSON.stringify(user));

            message.success('登录成功');
            navigate('/admin');
        } catch (error) {
            console.error(error);
            message.error('登录失败，请检查用户名和密码');
        } finally {
            setLoading(false);
        }
    };

    // 过滤只显示管理员账号
    const adminUsers = DEBUG_USERS.filter(u => u.role === 'admin');

    return (
        <div style={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            minHeight: '100vh',
            background: `linear-gradient(135deg, ${token.colorPrimaryBg} 0%, ${token.colorBgLayout} 100%)`
        }}>
            <Card 
                style={{ 
                    width: 420, 
                    boxShadow: '0 8px 24px rgba(0,0,0,0.12)',
                    borderRadius: 12
                }}
            >
                <div style={{ textAlign: 'center', marginBottom: 32 }}>
                    <SafetyCertificateOutlined style={{ fontSize: 48, color: token.colorPrimary }} />
                    <h1 style={{ margin: '16px 0 8px', color: token.colorTextHeading }}>
                        GameLink 管理后台
                    </h1>
                    <p style={{ color: token.colorTextSecondary, margin: 0 }}>
                        管理员登录
                    </p>
                </div>

                <Form
                    name="admin_login"
                    onFinish={onLogin}
                    layout="vertical"
                    size="large"
                    initialValues={{ username: 'admin@gameLink.com', password: '' }}
                >
                    <Form.Item
                        name="username"
                        rules={[{ required: true, message: '请输入管理员账号！' }]}
                    >
                        <Input 
                            prefix={<UserOutlined />} 
                            placeholder="管理员账号/邮箱" 
                        />
                    </Form.Item>

                    <Form.Item
                        name="password"
                        rules={[{ required: true, message: '请输入密码！' }]}
                    >
                        <Input.Password 
                            prefix={<LockOutlined />} 
                            placeholder="密码" 
                        />
                    </Form.Item>

                    <Form.Item style={{ marginBottom: 16 }}>
                        <Button 
                            type="primary" 
                            htmlType="submit" 
                            block 
                            loading={loading}
                            style={{ height: 44 }}
                        >
                            登录管理后台
                        </Button>
                    </Form.Item>
                </Form>

                {/* 开发环境快速登录 */}
                {ENABLE_QUICK_LOGIN && adminUsers.length > 0 && (
                    <div style={{ 
                        borderTop: `1px solid ${token.colorBorderSecondary}`, 
                        paddingTop: 16,
                        marginTop: 8
                    }}>
                        <div style={{ 
                            marginBottom: 12, 
                            color: token.colorTextSecondary, 
                            fontSize: 12, 
                            textAlign: 'center' 
                        }}>
                            开发环境快速登录
                        </div>
                        <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8, justifyContent: 'center' }}>
                            {adminUsers.map(user => (
                                <Button
                                    key={user.email}
                                    size="small"
                                    onClick={() => onLogin({ username: user.email, password: user.password })}
                                    style={{
                                        borderColor: user.color,
                                        color: user.color,
                                        fontSize: 12
                                    }}
                                >
                                    {user.label}
                                </Button>
                            ))}
                        </div>
                    </div>
                )}

                <div style={{ 
                    textAlign: 'center', 
                    marginTop: 24, 
                    color: token.colorTextTertiary,
                    fontSize: 12
                }}>
                    © 2025 GameLink. All rights reserved.
                </div>
            </Card>
        </div>
    );
};

export default AdminLogin;
