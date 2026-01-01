import React, { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Form, Input, Button, Card, App, Tabs, theme } from 'antd';
import { UserOutlined, LockOutlined, MailOutlined, PhoneOutlined } from '@ant-design/icons';
import { authApi } from '@/api/auth';
import { ENABLE_QUICK_LOGIN, DEBUG_USERS } from '@/config/debug';

const Auth: React.FC = () => {
    const { message } = App.useApp();
    const { token } = theme.useToken();
    const navigate = useNavigate();
    const location = useLocation();
    const [loginLoading, setLoginLoading] = useState(false);
    const [registerLoading, setRegisterLoading] = useState(false);
    
    // 根据路由决定默认 Tab
    const defaultTab = location.pathname === '/register' ? 'register' : 'login';

    const onLogin = async (values: { username: string; password: string }) => {
        setLoginLoading(true);
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
            const { token: authToken, user } = response.data;

            localStorage.setItem('token', authToken);
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
            setLoginLoading(false);
        }
    };

    const onRegister = async (values: { username: string; email: string; phone: string; password: string; confirm: string }) => {
        setRegisterLoading(true);
        try {
            await authApi.register({
                name: values.username,  // 后端期望 name 字段
                email: values.email,
                phone: values.phone,
                password: values.password,
                confirmPassword: values.confirm
            });
            message.success('注册成功！请登录。');
            // 切换到登录 Tab（通过改变 URL）
            navigate('/login');
        } catch (error) {
            console.error(error);
            message.error('注册失败，请重试。');
        } finally {
            setRegisterLoading(false);
        }
    };

    const handleTabChange = (key: string) => {
        navigate(key === 'register' ? '/register' : '/login', { replace: true });
    };

    const loginForm = (
        <Form
            name="login"
            onFinish={onLogin}
            layout="vertical"
            initialValues={{ username: 'admin@gameLink.com', password: '123456' }}
        >
            <Form.Item
                name="username"
                rules={[{ required: true, message: '请输入用户名！' }]}
            >
                <Input prefix={<UserOutlined />} placeholder="用户名/邮箱" size="large" />
            </Form.Item>

            <Form.Item
                name="password"
                rules={[{ required: true, message: '请输入密码！' }]}
            >
                <Input.Password prefix={<LockOutlined />} placeholder="密码" size="large" />
            </Form.Item>

            <Form.Item>
                <Button type="primary" htmlType="submit" block size="large" loading={loginLoading}>
                    登录
                </Button>
            </Form.Item>
        </Form>
    );

    const registerForm = (
        <Form
            name="register"
            onFinish={onRegister}
            layout="vertical"
            scrollToFirstError
        >
            <Form.Item
                name="username"
                rules={[{ required: true, message: '请输入用户名！' }]}
            >
                <Input prefix={<UserOutlined />} placeholder="用户名" size="large" />
            </Form.Item>

            <Form.Item
                name="email"
                rules={[
                    { type: 'email', message: '请输入有效的邮箱地址！' },
                    { required: true, message: '请输入邮箱！' },
                ]}
            >
                <Input prefix={<MailOutlined />} placeholder="邮箱" size="large" />
            </Form.Item>

            <Form.Item
                name="phone"
                rules={[{ required: true, message: '请输入手机号！' }]}
            >
                <Input prefix={<PhoneOutlined />} placeholder="手机号" size="large" />
            </Form.Item>

            <Form.Item
                name="password"
                rules={[
                    { required: true, message: '请输入密码！' },
                    { min: 6, message: '密码长度至少为6位！' },
                ]}
                hasFeedback
            >
                <Input.Password prefix={<LockOutlined />} placeholder="密码" size="large" />
            </Form.Item>

            <Form.Item
                name="confirm"
                dependencies={['password']}
                hasFeedback
                rules={[
                    { required: true, message: '请确认密码！' },
                    ({ getFieldValue }) => ({
                        validator(_, value) {
                            if (!value || getFieldValue('password') === value) {
                                return Promise.resolve();
                            }
                            return Promise.reject(new Error('两次输入的密码不一致！'));
                        },
                    }),
                ]}
            >
                <Input.Password prefix={<LockOutlined />} placeholder="确认密码" size="large" />
            </Form.Item>

            <Form.Item>
                <Button type="primary" htmlType="submit" block size="large" loading={registerLoading}>
                    注册
                </Button>
            </Form.Item>
        </Form>
    );

    return (
        <div style={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            minHeight: '100vh',
            background: token.colorBgLayout
        }}>
            <Card style={{ width: 400, boxShadow: '0 4px 12px rgba(0,0,0,0.1)' }}>
                <div style={{ textAlign: 'center', marginBottom: 24 }}>
                    <h1 style={{ margin: 0, color: token.colorTextHeading }}>GameLink</h1>
                    <p style={{ color: token.colorTextSecondary }}>游戏陪玩平台</p>
                </div>

                <Tabs 
                    activeKey={defaultTab}
                    onChange={handleTabChange}
                    centered
                    items={[
                        {
                            key: 'login',
                            label: '登录',
                            children: loginForm
                        },
                        {
                            key: 'register',
                            label: '注册',
                            children: registerForm
                        }
                    ]} 
                />

                {/* Quick Login for Development */}
                {ENABLE_QUICK_LOGIN && defaultTab === 'login' && (
                    <div style={{ marginTop: 16, borderTop: `1px solid ${token.colorBorderSecondary}`, paddingTop: 16 }}>
                        <div style={{ marginBottom: 8, color: token.colorTextSecondary, fontSize: 12, textAlign: 'center' }}>
                            开发环境快速登录
                        </div>
                        <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8, justifyContent: 'center' }}>
                            {DEBUG_USERS.map(user => (
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
            </Card>
        </div>
    );
};

export default Auth;
