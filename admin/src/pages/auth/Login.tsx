import React, { useState, useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { Form, Input, Button, Card, App, Tabs, theme, Checkbox } from 'antd';
import { LockOutlined, UserOutlined } from '@ant-design/icons';
import { authApi } from '@/api/auth';
import { ENABLE_QUICK_LOGIN, DEBUG_USERS } from '@/config/debug';
import { useAdmin } from '@/context/useAdmin';

import { logger } from '@/utils/logger';
const REMEMBER_KEY = 'gamelink_remember_login';

const Login: React.FC = () => {
    const { message } = App.useApp();
    const { token } = theme.useToken();
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);
    const [form] = Form.useForm();
    const { refreshMenus } = useAdmin();

    // 加载记住的账号（不存储密码）
    useEffect(() => {
        const saved = localStorage.getItem(REMEMBER_KEY);
        if (saved) {
            try {
                const { username, remember } = JSON.parse(saved);
                form.setFieldsValue({ username, remember });
            } catch {
                localStorage.removeItem(REMEMBER_KEY);
            }
        }
    }, [form]);

    const onFinish = async (values: { username: string; password: string; remember?: boolean }) => {
        setLoading(true);
        try {
            const res = await authApi.login({
                username: values.username,
                password: values.password
            });

            // API 返回格式: { success, code, message, data: { token, user } }
            const response = res.data as { 
                success?: boolean; 
                code?: number;
                message?: string;
                data?: { token: string; user: { id: number; role: string; [key: string]: unknown } } 
            };
            
            if (!response.success || !response.data) {
                // 处理业务逻辑错误（success=false 但 HTTP 200）
                const errorMsg = response.message || '登录失败';
                message.error(errorMsg);
                return;
            }
            
            const { token: authToken, user } = response.data;

            // 保存或清除记住的账号（不存储密码）
            if (values.remember) {
                localStorage.setItem(REMEMBER_KEY, JSON.stringify({
                    username: values.username,
                    remember: true
                }));
            } else {
                localStorage.removeItem(REMEMBER_KEY);
            }

            localStorage.setItem('token', authToken);
            localStorage.setItem('user_role', user.role);
            // 保存用户信息，同时兼容两种字段名格式
            localStorage.setItem('user_info', JSON.stringify({
                ...user,
                username: user.name || user.username,
                avatar: user.avatarUrl || user.avatar,
            }));

            // 登录成功后刷新菜单和权限数据
            await refreshMenus();

            message.success('登录成功');

            const role = user.role.toUpperCase();
            if (role === 'ADMIN') {
                navigate('/admin');
            } else if (role === 'PLAYER') {
                navigate('/player');
            } else {
                navigate('/');
            }
        } catch (error: unknown) {
            logger.error('登录错误:', error);
            
            // 处理 Axios 错误响应
            if (error && typeof error === 'object' && 'response' in error) {
                const axiosError = error as { response?: { status?: number; data?: { message?: string; code?: number } } };
                const status = axiosError.response?.status;
                const errorData = axiosError.response?.data;
                
                if (status === 401) {
                    message.error(errorData?.message || '用户名或密码错误');
                } else if (status === 403) {
                    message.error('账号已被禁用，请联系管理员');
                } else if (status === 404) {
                    message.error('用户不存在');
                } else if (status === 429) {
                    message.error('登录尝试次数过多，请稍后再试');
                } else if (status && status >= 500) {
                    message.error('服务器错误，请稍后重试');
                } else if (errorData?.message) {
                    message.error(errorData.message);
                } else {
                    message.error('登录失败，请检查用户名和密码');
                }
            } else if (error && typeof error === 'object' && 'message' in error) {
                // 网络错误或其他错误
                const err = error as { message: string };
                if (err.message.includes('Network Error') || err.message.includes('timeout')) {
                    message.error('网络连接失败，请检查网络后重试');
                } else {
                    message.error(err.message || '登录失败');
                }
            } else {
                message.error('登录失败，请稍后重试');
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <div style={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            minHeight: '100vh',
            background: `linear-gradient(135deg, ${token.colorBgLayout} 0%, rgba(122, 204, 53, 0.05) 100%)`,
            padding: '24px',
            position: 'relative',
            overflow: 'hidden'
        }}>
            {/* 背景装饰元素 */}
            <div style={{
                position: 'absolute',
                width: '600px',
                height: '600px',
                borderRadius: '50%',
                background: 'radial-gradient(circle, rgba(122, 204, 53, 0.08) 0%, transparent 70%)',
                top: '-300px',
                right: '-300px',
                pointerEvents: 'none'
            }} />
            <div style={{
                position: 'absolute',
                width: '400px',
                height: '400px',
                borderRadius: '50%',
                background: 'radial-gradient(circle, rgba(122, 204, 53, 0.05) 0%, transparent 70%)',
                bottom: '-200px',
                left: '-200px',
                pointerEvents: 'none'
            }} />

            <Card
                style={{
                    width: '100%',
                    maxWidth: 420,
                    boxShadow: '0 8px 32px rgba(0, 0, 0, 0.08)',
                    borderRadius: token.borderRadiusLG,
                    border: `1px solid ${token.colorBorderSecondary}`,
                    position: 'relative',
                    zIndex: 1
                }}
                bodyStyle={{
                    padding: '40px 32px'
                }}
            >
                <div style={{ textAlign: 'center', marginBottom: 32 }}>
                    <h1 style={{
                        margin: '0 0 8px 0',
                        color: token.colorTextHeading,
                        fontSize: '28px',
                        fontWeight: 600,
                        letterSpacing: '-0.5px'
                    }}>
                        GameLink
                    </h1>
                    <p style={{
                        margin: 0,
                        color: token.colorTextSecondary,
                        fontSize: '14px'
                    }}>
                        游戏陪玩平台 · 管理后台
                    </p>
                </div>

                <Tabs
                    defaultActiveKey="account"
                    items={[
                        {
                            key: 'account',
                            label: '账号登录',
                            children: (
                                <Form
                                    name="login"
                                    form={form}
                                    onFinish={onFinish}
                                    layout="vertical"
                                    requiredMark={false}
                                >
                                    <Form.Item
                                        name="username"
                                        label={<span style={{ fontSize: '14px', fontWeight: 500 }}>用户名</span>}
                                        rules={[{ required: true, message: '请输入用户名！' }]}
                                    >
                                        <Input
                                            prefix={<UserOutlined style={{ color: token.colorTextSecondary }} />}
                                            placeholder="请输入用户名"
                                            size="large"
                                            style={{
                                                borderRadius: token.borderRadiusSM
                                            }}
                                        />
                                    </Form.Item>

                                    <Form.Item
                                        name="password"
                                        label={<span style={{ fontSize: '14px', fontWeight: 500 }}>密码</span>}
                                        rules={[{ required: true, message: '请输入密码！' }]}
                                    >
                                        <Input.Password
                                            prefix={<LockOutlined style={{ color: token.colorTextSecondary }} />}
                                            placeholder="请输入密码"
                                            size="large"
                                            style={{
                                                borderRadius: token.borderRadiusSM
                                            }}
                                        />
                                    </Form.Item>

                                    <Form.Item name="remember" valuePropName="checked">
                                        <Checkbox style={{ fontSize: '14px' }}>记住密码</Checkbox>
                                    </Form.Item>

                                    <Form.Item style={{ marginBottom: 16 }}>
                                        <Button
                                            type="primary"
                                            htmlType="submit"
                                            block
                                            size="large"
                                            loading={loading}
                                            style={{
                                                height: '44px',
                                                borderRadius: token.borderRadiusSM,
                                                fontSize: '16px',
                                                fontWeight: 500,
                                                boxShadow: loading ? 'none' : '0 2px 8px rgba(122, 204, 53, 0.25)'
                                            }}
                                        >
                                            登录
                                        </Button>
                                    </Form.Item>

                                    <div style={{
                                        textAlign: 'center',
                                        fontSize: '14px',
                                        color: token.colorTextSecondary
                                    }}>
                                        没有账号？{' '}
                                        <Link to="/register" style={{ fontWeight: 500 }}>
                                            立即注册
                                        </Link>
                                    </div>
                                </Form>
                            )
                        }
                    ]}
                />

                {/* Quick Login for Development */}
                {ENABLE_QUICK_LOGIN && (
                    <div style={{
                        marginTop: 32,
                        paddingTop: 24,
                        borderTop: `1px solid ${token.colorBorderSecondary}`
                    }}>
                        <div style={{
                            marginBottom: 16,
                            color: token.colorTextSecondary,
                            fontSize: '13px',
                            textAlign: 'center',
                            fontWeight: 500
                        }}>
                            开发环境快速登录
                        </div>
                        <div style={{
                            display: 'flex',
                            flexWrap: 'wrap',
                            gap: 12,
                            justifyContent: 'center'
                        }}>
                            {DEBUG_USERS.map(user => (
                                <Button
                                    key={user.email}
                                    size="large"
                                    onClick={() => onFinish({ username: user.email, password: user.password })}
                                    style={{
                                        borderColor: user.color,
                                        color: user.color,
                                        fontSize: '13px',
                                        borderRadius: token.borderRadiusSM,
                                        borderWidth: '1.5px'
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
