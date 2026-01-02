import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Form, Input, Button, Card, App, theme, Checkbox } from 'antd';
import { UserOutlined, LockOutlined, SafetyCertificateOutlined } from '@ant-design/icons';
import { authApi } from '@/api/auth';
import { ENABLE_QUICK_LOGIN, DEBUG_USERS } from '@/config/debug';

import { logger } from '@/utils/logger';
const REMEMBER_KEY = 'gamelink_admin_remember';

/**
 * 管理后台登录页面
 * 独立于用户端登录，专为管理员设计
 */
const AdminLogin: React.FC = () => {
    const { token } = theme.useToken();
    const { message } = App.useApp(); // 使用 App.useApp() 获取 message 实例
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);
    const [form] = Form.useForm();

    // 加载记住的账号密码
    useEffect(() => {
        const saved = localStorage.getItem(REMEMBER_KEY);
        if (saved) {
            try {
                const { username, password, remember } = JSON.parse(saved);
                form.setFieldsValue({ username, password, remember });
            } catch {
                localStorage.removeItem(REMEMBER_KEY);
            }
        }
    }, [form]);

    const onLogin = async (values: { username: string; password: string; remember?: boolean }) => {
        setLoading(true);
        try {
            const res = await authApi.login({
                username: values.username,
                password: values.password
            });

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
            const role = user.role.toUpperCase();

            // 验证是否为管理员角色
            if (role !== 'ADMIN') {
                message.error('您没有管理后台访问权限');
                return;
            }

            // 保存或清除记住的账号密码
            if (values.remember) {
                localStorage.setItem(REMEMBER_KEY, JSON.stringify({
                    username: values.username,
                    password: values.password,
                    remember: true
                }));
            } else {
                localStorage.removeItem(REMEMBER_KEY);
            }

            localStorage.setItem('token', authToken);
            localStorage.setItem('user_role', user.role);
            localStorage.setItem('user_info', JSON.stringify(user));

            message.success('登录成功');
            navigate('/admin');
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
                    form={form}
                    onFinish={onLogin}
                    layout="vertical"
                    size="large"
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

                    <Form.Item name="remember" valuePropName="checked">
                        <Checkbox>记住密码</Checkbox>
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
