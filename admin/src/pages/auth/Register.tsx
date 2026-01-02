import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { Form, Input, Button, Card, App, Tabs, theme } from 'antd';
import { UserOutlined, LockOutlined, MailOutlined, PhoneOutlined } from '@ant-design/icons';
import { authApi } from '@/api/auth';

import { logger } from '@/utils/logger';
const Register: React.FC = () => {
    const { message } = App.useApp();
    const { token } = theme.useToken();
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);

    const onFinish = async (values: { name: string; email: string; phone: string; password: string }) => {
        setLoading(true);
        try {
            await authApi.register(values);
            message.success('注册成功！请登录。');
            navigate('/login');
        } catch (error) {
            logger.error("Operation failed", error);
            message.error('注册失败，请重试。');
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
                        label: '账号注册',
                        children: (
                            <Form
                                name="register"
                                onFinish={onFinish}
                                layout="vertical"
                                scrollToFirstError
                            >
                                <Form.Item
                                    name="name"
                                    rules={[{ required: true, message: '请输入您的昵称！' }]}
                                >
                                    <Input prefix={<UserOutlined />} placeholder="昵称" size="large" />
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
                                    <Button type="primary" htmlType="submit" block size="large" loading={loading}>
                                        注册
                                    </Button>
                                </Form.Item>

                                <div style={{ textAlign: 'center' }}>
                                    已有账号？ <Link to="/login">立即登录</Link>
                                </div>
                            </Form>
                        )
                    }
                ]} />
            </Card>
        </div>
    );
};

export default Register;
