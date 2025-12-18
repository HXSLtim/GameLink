import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { Form, Input, Button, Card, message, Tabs, theme } from 'antd';
import { LockOutlined, UserOutlined } from '@ant-design/icons';
import { authApi } from '@/api/auth';
import { ENABLE_QUICK_LOGIN, DEBUG_USERS } from '@/config/debug';

const Login: React.FC = () => {
    const { token } = theme.useToken();
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);

    const onFinish = async (values: { username: string; password: string }) => {
        // ... (keep existing logic)
        setLoading(true);
        try {
            const res = await authApi.login({
                username: values.username,
                password: values.password
            });

            // API 返回格式: { success, code, message, data: { token, user } }
            const response = res.data as { success?: boolean; data?: { token: string; user: { id: number; role: string; [key: string]: unknown } } };
            if (!response.success || !response.data) {
                throw new Error('登录响应格式错误');
            }
            const { token, user } = response.data;

            localStorage.setItem('token', token);
            localStorage.setItem('user_role', user.role);
            localStorage.setItem('user_info', JSON.stringify(user));

            message.success('登录成功');

            const role = user.role.toUpperCase();
            if (role === 'ADMIN') {
                navigate('/admin');
            } else if (role === 'COMPANION') {
                navigate('/companion');
            } else {
                navigate('/');
            }
        } catch (error) {
            console.error(error);
            message.error('登录失败，请检查用户名和密码');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div style={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            height: '100vh',
            background: token.colorBgLayout
        }}>
            <Card style={{ width: 400, boxShadow: '0 4px 12px rgba(0,0,0,0.1)' }}>
                <div style={{ textAlign: 'center', marginBottom: 24 }}>
                    <h1 style={{ margin: 0, color: token.colorTextHeading }}>GameLink</h1>
                    <p style={{ color: token.colorTextSecondary }}>游戏陪玩平台</p>
                </div>

                <Tabs defaultActiveKey="account" items={[
                    {
                        key: 'account',
                        label: '账号登录',
                        children: (
                            <Form
                                name="login"
                                onFinish={onFinish}
                                layout="vertical"
                                initialValues={{ username: 'admin@gameLink.com', password: '123456' }}
                            >
                                <Form.Item
                                    name="username"
                                    rules={[{ required: true, message: '请输入用户名！' }]}
                                >
                                    <Input prefix={<UserOutlined />} placeholder="用户名" size="large" />
                                </Form.Item>

                                <Form.Item
                                    name="password"
                                    rules={[{ required: true, message: '请输入密码！' }]}
                                >
                                    <Input.Password prefix={<LockOutlined />} placeholder="密码" size="large" />
                                </Form.Item>

                                <Form.Item>
                                    <Button type="primary" htmlType="submit" block size="large" loading={loading}>
                                        登录
                                    </Button>
                                </Form.Item>

                                <div style={{ textAlign: 'center' }}>
                                    没有账号？ <Link to="/register">立即注册</Link>
                                </div>
                            </Form>
                        )
                    }
                ]} />

                {/* Quick Login for Development */}
                {ENABLE_QUICK_LOGIN && (
                    <div style={{ marginTop: 24, borderTop: `1px solid ${token.colorBorderSecondary}`, paddingTop: 16 }}>
                        <div style={{ marginBottom: 8, color: token.colorTextSecondary, fontSize: 12, textAlign: 'center' }}>
                            开发环境快速登录
                        </div>
                        <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8, justifyContent: 'center' }}>
                            {DEBUG_USERS.map(user => (
                                <Button
                                    key={user.email}
                                    size="small"
                                    onClick={() => onFinish({ username: user.email, password: user.password })}
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
            </Card>
        </div>
    );
};

export default Login;
