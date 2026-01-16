/**
 * 个人中心页面
 * 显示当前登录用户信息，支持修改密码和个人信息
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Card,
    Form,
    Input,
    Button,
    Avatar,
    Descriptions,
    Divider,
    message,
    Space,
    Tag,
    Row,
    Col,
    Modal,
    theme,
    Spin,
    Upload,
} from 'antd';
import type { UploadProps } from 'antd';
import {
    UserOutlined,
    MailOutlined,
    PhoneOutlined,
    LockOutlined,
    EditOutlined,
    SafetyCertificateOutlined,
    UploadOutlined,
} from '@ant-design/icons';
import { PageContainer } from '@/components';
import { useUserInfo, useIsAuthenticated, useIsHydrated } from '@/stores/modules/authStore';
import { adminApi, type User } from '@/api/admin';
import dayjs from 'dayjs';
import { logger } from '@/utils/logger';

/**
 * 角色映射
 */
const roleMap: Record<string, { color: string; text: string }> = {
    user: { color: 'blue', text: '普通用户' },
    player: { color: 'purple', text: '陪玩师' },
    admin: { color: 'red', text: '管理员' },
};

/**
 * 状态映射
 */
const statusMap: Record<string, { color: string; text: string }> = {
    active: { color: 'success', text: '正常' },
    banned: { color: 'error', text: '已封禁' },
    suspended: { color: 'warning', text: '已停用' },
};

const ProfilePage: React.FC = () => {
    const { token } = theme.useToken();

    // Super Dev 最佳实践: 使用选择器精确订阅，避免不必要的重渲染
    const authUser = useUserInfo();
    const isAuthenticated = useIsAuthenticated();
    const isHydrated = useIsHydrated();

    const [loading, setLoading] = useState(false);
    const [userInfo, setUserInfo] = useState<User | null>(null);
    const [passwordModalVisible, setPasswordModalVisible] = useState(false);
    const [editModalVisible, setEditModalVisible] = useState(false);
    const [passwordLoading, setPasswordLoading] = useState(false);
    const [editLoading, setEditLoading] = useState(false);
    const [passwordForm] = Form.useForm();
    const [editForm] = Form.useForm();

    // 加载用户信息 - 定义为 useCallback 以便在多处调用
    const loadUserInfo = useCallback(async () => {
        if (!authUser) {
            return;
        }
        try {
            setLoading(true);
            // 从 authUser 构建 User 对象
            const user: User = {
                id: authUser.id,
                name: authUser.name || '未知',
                email: authUser.email || '',
                phone: authUser.phone || '',
                avatarUrl: authUser.avatar, // UserInfo.avatar → User.avatarUrl
                role: (authUser.role as 'user' | 'player' | 'admin') || 'admin',
                status: 'active', // 管理员默认状态
                createdAt: authUser.createdAt || new Date().toISOString(),
                updatedAt: authUser.updatedAt,
            };
            setUserInfo(user);
            logger.info('[Profile] User info loaded:', user);
        } catch (error) {
            logger.error('Failed to load user info', error);
            message.error('获取用户信息失败');
        } finally {
            setLoading(false);
        }
    }, [authUser]);

    // 初始加载用户信息
    useEffect(() => {
        if (!isHydrated || !authUser) {
            return;
        }
        loadUserInfo();
    }, [isHydrated, authUser, loadUserInfo]);

    // 使用水合状态而非 setTimeout，确保 Zustand persist 完成后再渲染
    if (!isHydrated) {
        return (
            <PageContainer title="个人中心">
                <div style={{ textAlign: 'center', padding: 100 }}>
                    <Spin size="large" tip="加载中..." />
                </div>
            </PageContainer>
        );
    }

    // 添加调试日志
    logger.info('[Profile] Component render', { authUser, isAuthenticated });

    // 修改密码
    const handleChangePassword = async () => {
        try {
            const values = await passwordForm.validateFields();
            
            if (values.newPassword !== values.confirmPassword) {
                message.error('两次输入的密码不一致');
                return;
            }

            setPasswordLoading(true);
            
            // 调用修改密码 API（通过更新用户信息实现）
            if (userInfo) {
                const response = await adminApi.updateUser(userInfo.id, {
                    name: userInfo.name,
                    email: userInfo.email,
                    phone: userInfo.phone || '',
                    role: userInfo.role,
                    status: userInfo.status,
                    password: values.newPassword,
                });

                if (response.data.success) {
                    message.success('密码修改成功');
                    setPasswordModalVisible(false);
                    passwordForm.resetFields();
                } else {
                    message.error(response.data.message || '密码修改失败');
                }
            }
        } catch (error) {
            logger.error('Failed to change password', error);
            message.error('密码修改失败');
        } finally {
            setPasswordLoading(false);
        }
    };

    // 修改个人信息
    const handleEditProfile = async () => {
        try {
            const values = await editForm.validateFields();
            setEditLoading(true);

            if (userInfo) {
                const response = await adminApi.updateUser(userInfo.id, {
                    name: values.name,
                    email: values.email,
                    phone: values.phone || '',
                    avatarUrl: values.avatarUrl,
                    role: userInfo.role,
                    status: userInfo.status,
                });

                if (response.data.success) {
                    message.success('个人信息更新成功');
                    setEditModalVisible(false);
                    // 更新本地存储
                    const storedUser = localStorage.getItem('user_info');
                    if (storedUser) {
                        const parsed = JSON.parse(storedUser);
                        parsed.username = values.name;
                        parsed.avatar = values.avatarUrl;
                        localStorage.setItem('user_info', JSON.stringify(parsed));
                    }
                    loadUserInfo();
                } else {
                    message.error(response.data.message || '更新失败');
                }
            }
        } catch (error) {
            logger.error('Failed to update profile', error);
            message.error('更新个人信息失败');
        } finally {
            setEditLoading(false);
        }
    };

    // 打开编辑弹窗
    const openEditModal = () => {
        if (userInfo) {
            editForm.setFieldsValue({
                name: userInfo.name,
                email: userInfo.email,
                phone: userInfo.phone,
                avatarUrl: userInfo.avatarUrl,
            });
        }
        setEditModalVisible(true);
    };

    // 头像上传配置
    const uploadProps: UploadProps = {
        name: 'file',
        action: '/api/v1/upload/image',
        headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
        showUploadList: false,
        onChange(info) {
            if (info.file.status === 'done') {
                const url = info.file.response?.data?.url;
                if (url) {
                    editForm.setFieldValue('avatarUrl', url);
                    message.success('头像上传成功');
                }
            } else if (info.file.status === 'error') {
                message.error('头像上传失败');
            }
        },
    };

    // 未登录状态
    if (!authUser || !isAuthenticated) {
        return (
            <PageContainer title="个人中心">
                <div style={{ textAlign: 'center', padding: 100 }}>
                    <p>请先登录</p>
                </div>
            </PageContainer>
        );
    }

    if (loading) {
        return (
            <PageContainer title="个人中心">
                <div style={{ textAlign: 'center', padding: 100 }}>
                    <Spin size="large" />
                </div>
            </PageContainer>
        );
    }

    return (
        <PageContainer title="个人中心" subTitle="查看和管理您的账号信息">
            <Row gutter={[24, 24]}>
                {/* 左侧：用户信息卡片 */}
                <Col xs={24} lg={8}>
                    <Card>
                        <div style={{ textAlign: 'center', marginBottom: 24 }}>
                            <Avatar
                                size={100}
                                src={userInfo?.avatarUrl}
                                icon={<UserOutlined />}
                                style={{ backgroundColor: token.colorPrimary }}
                            />
                            <h2 style={{ marginTop: 16, marginBottom: 8 }}>{userInfo?.name}</h2>
                            <Space>
                                <Tag color={roleMap[userInfo?.role || 'user']?.color}>
                                    {roleMap[userInfo?.role || 'user']?.text}
                                </Tag>
                                <Tag color={statusMap[userInfo?.status || 'active']?.color}>
                                    {statusMap[userInfo?.status || 'active']?.text}
                                </Tag>
                            </Space>
                        </div>

                        <Divider />

                        <Space direction="vertical" style={{ width: '100%' }} size="middle">
                            <Button
                                type="primary"
                                icon={<EditOutlined />}
                                block
                                onClick={openEditModal}
                            >
                                编辑资料
                            </Button>
                            <Button
                                icon={<LockOutlined />}
                                block
                                onClick={() => setPasswordModalVisible(true)}
                            >
                                修改密码
                            </Button>
                        </Space>
                    </Card>
                </Col>

                {/* 右侧：详细信息 */}
                <Col xs={24} lg={16}>
                    <Card title="账号信息">
                        <Descriptions column={{ xs: 1, sm: 2 }} labelStyle={{ fontWeight: 500 }}>
                            <Descriptions.Item label={<><UserOutlined /> 用户名</>}>
                                {userInfo?.name || '-'}
                            </Descriptions.Item>
                            <Descriptions.Item label={<><MailOutlined /> 邮箱</>}>
                                {userInfo?.email || '-'}
                            </Descriptions.Item>
                            <Descriptions.Item label={<><PhoneOutlined /> 手机号</>}>
                                {userInfo?.phone || '-'}
                            </Descriptions.Item>
                            <Descriptions.Item label={<><SafetyCertificateOutlined /> 用户ID</>}>
                                {userInfo?.id || '-'}
                            </Descriptions.Item>
                            <Descriptions.Item label="注册时间">
                                {userInfo?.createdAt ? dayjs(userInfo.createdAt).format('YYYY-MM-DD HH:mm') : '-'}
                            </Descriptions.Item>
                            <Descriptions.Item label="最后登录">
                                {userInfo?.lastLoginAt ? dayjs(userInfo.lastLoginAt).format('YYYY-MM-DD HH:mm') : '-'}
                            </Descriptions.Item>
                        </Descriptions>
                    </Card>

                    <Card title="安全设置" style={{ marginTop: 24 }}>
                        <Descriptions column={1} labelStyle={{ fontWeight: 500 }}>
                            <Descriptions.Item label="登录密码">
                                <Space>
                                    <span>******</span>
                                    <Button
                                        type="link"
                                        size="small"
                                        onClick={() => setPasswordModalVisible(true)}
                                    >
                                        修改
                                    </Button>
                                </Space>
                            </Descriptions.Item>
                        </Descriptions>
                    </Card>
                </Col>
            </Row>

            {/* 修改密码弹窗 */}
            <Modal
                title="修改密码"
                open={passwordModalVisible}
                onOk={handleChangePassword}
                onCancel={() => {
                    setPasswordModalVisible(false);
                    passwordForm.resetFields();
                }}
                confirmLoading={passwordLoading}
                okText="确认修改"
                cancelText="取消"
            >
                <Form form={passwordForm} layout="vertical" style={{ marginTop: 16 }}>
                    <Form.Item
                        name="newPassword"
                        label="新密码"
                        rules={[
                            { required: true, message: '请输入新密码' },
                            { min: 8, message: '密码至少8位' },
                            {
                                pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/,
                                message: '密码需包含大小写字母和数字',
                            },
                        ]}
                    >
                        <Input.Password placeholder="请输入新密码" prefix={<LockOutlined />} />
                    </Form.Item>
                    <Form.Item
                        name="confirmPassword"
                        label="确认密码"
                        dependencies={['newPassword']}
                        rules={[
                            { required: true, message: '请确认新密码' },
                            ({ getFieldValue }) => ({
                                validator(_, value) {
                                    if (!value || getFieldValue('newPassword') === value) {
                                        return Promise.resolve();
                                    }
                                    return Promise.reject(new Error('两次输入的密码不一致'));
                                },
                            }),
                        ]}
                    >
                        <Input.Password placeholder="请再次输入新密码" prefix={<LockOutlined />} />
                    </Form.Item>
                </Form>
            </Modal>

            {/* 编辑资料弹窗 */}
            <Modal
                title="编辑资料"
                open={editModalVisible}
                onOk={handleEditProfile}
                onCancel={() => {
                    setEditModalVisible(false);
                    editForm.resetFields();
                }}
                confirmLoading={editLoading}
                okText="保存"
                cancelText="取消"
            >
                <Form form={editForm} layout="vertical" style={{ marginTop: 16 }}>
                    <Form.Item label="头像" name="avatarUrl">
                        <Space direction="vertical" align="center" style={{ width: '100%' }}>
                            <Avatar
                                size={80}
                                src={editForm.getFieldValue('avatarUrl')}
                                icon={<UserOutlined />}
                            />
                            <Upload {...uploadProps}>
                                <Button icon={<UploadOutlined />}>上传头像</Button>
                            </Upload>
                            <Input
                                placeholder="或输入头像URL"
                                style={{ marginTop: 8 }}
                                onChange={(e) => editForm.setFieldValue('avatarUrl', e.target.value)}
                            />
                        </Space>
                    </Form.Item>
                    <Form.Item
                        name="name"
                        label="用户名"
                        rules={[{ required: true, message: '请输入用户名' }]}
                    >
                        <Input placeholder="请输入用户名" prefix={<UserOutlined />} />
                    </Form.Item>
                    <Form.Item
                        name="email"
                        label="邮箱"
                        rules={[
                            { required: true, message: '请输入邮箱' },
                            { type: 'email', message: '请输入有效的邮箱地址' },
                        ]}
                    >
                        <Input placeholder="请输入邮箱" prefix={<MailOutlined />} />
                    </Form.Item>
                    <Form.Item name="phone" label="手机号">
                        <Input placeholder="请输入手机号" prefix={<PhoneOutlined />} />
                    </Form.Item>
                </Form>
            </Modal>
        </PageContainer>
    );
};

export default ProfilePage;
