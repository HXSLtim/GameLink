import React, { useState, useEffect } from 'react';
import { Card, Form, Input, Switch, Button, Tabs, Select, ColorPicker, Divider, message, List, Space, Typography, Alert, Radio, InputNumber, Spin, Descriptions } from 'antd';
import { SaveOutlined, SettingOutlined, SafetyOutlined, BellOutlined, BgColorsOutlined, SyncOutlined, StarOutlined, ReloadOutlined, UndoOutlined } from '@ant-design/icons';
import { motion } from 'framer-motion';
import { forceInit } from '@/services/init';
import { reviewSettingsApi } from '@/api/review';
import type { ReviewDisplaySettings, UpdateSettingsFormData } from '@/types/review';
import { SORT_BY_TEXT } from '@/types/review';

const { Text, Paragraph } = Typography;

// 评价设置默认值
const DEFAULT_REVIEW_SETTINGS: UpdateSettingsFormData = {
    sortBy: 'time',
    minScore: 1,
    showAnonymous: true,
    pageSize: 10,
};

const Settings: React.FC = () => {
    const [form] = Form.useForm();
    const [reviewForm] = Form.useForm<UpdateSettingsFormData>();
    const [initializing, setInitializing] = useState(false);
    const [lastInitTime, setLastInitTime] = useState<string | null>(
        localStorage.getItem('app_init_timestamp')
    );

    // 评价设置状态
    const [reviewLoading, setReviewLoading] = useState(true);
    const [reviewSaving, setReviewSaving] = useState(false);
    const [currentReviewSettings, setCurrentReviewSettings] = useState<ReviewDisplaySettings | null>(null);

    // 加载评价设置
    const fetchReviewSettings = async () => {
        setReviewLoading(true);
        try {
            const response = await reviewSettingsApi.getSettings() as unknown as {
                success: boolean;
                data: ReviewDisplaySettings;
            };
            if (response.success) {
                setCurrentReviewSettings(response.data);
                reviewForm.setFieldsValue({
                    sortBy: response.data.sortBy,
                    minScore: response.data.minScore,
                    showAnonymous: response.data.showAnonymous,
                    pageSize: response.data.pageSize,
                });
            }
        } catch {
            // 如果API不存在，使用默认值
            reviewForm.setFieldsValue(DEFAULT_REVIEW_SETTINGS);
        } finally {
            setReviewLoading(false);
        }
    };

    useEffect(() => {
        fetchReviewSettings();
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    // 保存评价设置
    const handleSaveReviewSettings = async () => {
        try {
            const values = await reviewForm.validateFields();
            setReviewSaving(true);
            const response = await reviewSettingsApi.updateSettings(values) as unknown as {
                success: boolean;
                data: ReviewDisplaySettings;
            };
            if (response.success) {
                message.success('评价设置已保存');
                setCurrentReviewSettings(response.data);
            }
        } catch {
            message.error('保存评价设置失败');
        } finally {
            setReviewSaving(false);
        }
    };

    // 重置评价设置为默认值
    const handleResetReviewSettings = () => {
        reviewForm.setFieldsValue(DEFAULT_REVIEW_SETTINGS);
        message.info('已重置为默认值，请点击保存生效');
    };

    // 恢复评价设置
    const handleRevertReviewSettings = () => {
        if (currentReviewSettings) {
            reviewForm.setFieldsValue({
                sortBy: currentReviewSettings.sortBy,
                minScore: currentReviewSettings.minScore,
                showAnonymous: currentReviewSettings.showAnonymous,
                pageSize: currentReviewSettings.pageSize,
            });
            message.info('已恢复为当前保存的设置');
        }
    };

    const onFinish = (values: Record<string, unknown>) => {
        console.log('Success:', values);
        message.success('设置已保存');
    };

    /**
     * 执行系统初始化
     * 同步菜单、权限并为超级管理员分配所有权限
     */
    const handleInit = async () => {
        setInitializing(true);
        const hide = message.loading('正在初始化系统...', 0);

        try {
            const result = await forceInit({ verbose: true });

            if (result.success) {
                message.success({
                    content: `初始化成功！耗时 ${result.duration}ms`,
                    duration: 3,
                });

                // 更新最后初始化时间
                setLastInitTime(Date.now().toString());

                // 显示详细信息
                if (result.menuSync) {
                    console.log('[菜单同步]', result.menuSync);
                }
                if (result.permissionSync) {
                    console.log('[权限同步]', result.permissionSync);
                }
                if (result.superAdminAssign) {
                    console.log('[超管权限]', result.superAdminAssign);
                }
            } else {
                message.error({
                    content: `初始化失败：${result.errors.join(', ')}`,
                    duration: 5,
                });
            }
        } catch (error) {
            console.error('初始化异常:', error);
            message.error({
                content: `初始化异常：${error instanceof Error ? error.message : '未知错误'}`,
                duration: 5,
            });
        } finally {
            hide();
            setInitializing(false);
        }
    };

    /**
     * 格式化时间显示
     */
    const formatLastInitTime = () => {
        if (!lastInitTime) return '从未初始化';
        const date = new Date(parseInt(lastInitTime, 10));
        return date.toLocaleString('zh-CN', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit',
        });
    };

    const items = [
        {
            key: '1',
            label: (
                <span>
                    <SettingOutlined />
                    基本设置
                </span>
            ),
            children: (
                <Form
                    form={form}
                    layout="vertical"
                    initialValues={{
                        siteName: 'GameLink 游戏陪玩平台',
                        siteDescription: '专业的游戏陪玩与社交平台',
                        maintenanceMode: false,
                        language: 'zh_CN',
                    }}
                    onFinish={onFinish}
                >
                    <Form.Item label="平台名称" name="siteName" rules={[{ required: true, message: '请输入平台名称' }]}>
                        <Input placeholder="请输入平台名称" />
                    </Form.Item>
                    <Form.Item label="平台描述" name="siteDescription">
                        <Input.TextArea rows={4} placeholder="请输入平台描述" />
                    </Form.Item>
                    <Form.Item label="默认语言" name="language">
                        <Select>
                            <Select.Option value="zh_CN">简体中文</Select.Option>
                            <Select.Option value="en_US">English</Select.Option>
                        </Select>
                    </Form.Item>
                    <Form.Item label="维护模式" name="maintenanceMode" valuePropName="checked" help="开启后，除管理员外用户将无法访问">
                        <Switch checkedChildren="开启" unCheckedChildren="关闭" />
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" htmlType="submit" icon={<SaveOutlined />} style={{ backgroundColor: '#5865F2' }}>
                            保存更改
                        </Button>
                    </Form.Item>
                </Form>
            ),
        },
        {
            key: '2',
            label: (
                <span>
                    <BgColorsOutlined />
                    主题设置
                </span>
            ),
            children: (
                <Form layout="vertical" initialValues={{ primaryColor: '#5865F2', darkMode: true }}>
                    <Form.Item label="主色调" name="primaryColor">
                        <ColorPicker showText />
                    </Form.Item>
                    <Form.Item label="深色模式" name="darkMode" valuePropName="checked">
                        <Switch checkedChildren="开启" unCheckedChildren="关闭" />
                    </Form.Item>
                    <Divider />
                    <Form.Item label="侧边栏风格">
                        <Select defaultValue="dark">
                            <Select.Option value="dark">深色 (Discord 风格)</Select.Option>
                            <Select.Option value="light">浅色</Select.Option>
                        </Select>
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" icon={<SaveOutlined />} style={{ backgroundColor: '#5865F2' }}>
                            应用主题
                        </Button>
                    </Form.Item>
                </Form>
            ),
        },
        {
            key: '3',
            label: (
                <span>
                    <SafetyOutlined />
                    安全设置
                </span>
            ),
            children: (
                <Form layout="vertical">
                    <Form.Item label="注册限制">
                        <Switch checkedChildren="开放注册" unCheckedChildren="停止注册" defaultChecked />
                    </Form.Item>
                    <Form.Item label="密码强度要求">
                        <Select defaultValue="medium">
                            <Select.Option value="low">低 (仅长度限制)</Select.Option>
                            <Select.Option value="medium">中 (需包含字母和数字)</Select.Option>
                            <Select.Option value="high">高 (需包含特殊字符)</Select.Option>
                        </Select>
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" style={{ backgroundColor: '#5865F2' }}>
                            更新安全策略
                        </Button>
                    </Form.Item>
                </Form>
            ),
        },
        {
            key: '4',
            label: (
                <span>
                    <BellOutlined />
                    通知设置
                </span>
            ),
            children: (
                <List
                    dataSource={[
                        { title: '新订单通知', desc: '当有新订单创建时发送通知' },
                        { title: '用户注册通知', desc: '当有新用户注册时发送通知' },
                        { title: '系统异常报警', desc: '系统发生严重错误时发送邮件报警' },
                    ]}
                    renderItem={(item: { title: string; desc: string }) => (
                        <List.Item extra={<Switch defaultChecked />}>
                            <List.Item.Meta
                                title={<span style={{ color: '#fff' }}>{item.title}</span>}
                                description={<span style={{ color: 'rgba(255,255,255,0.45)' }}>{item.desc}</span>}
                            />
                        </List.Item>
                    )}
                />
            ),
        },
        {
            key: '5',
            label: (
                <span>
                    <StarOutlined />
                    评价设置
                </span>
            ),
            children: reviewLoading ? (
                <div style={{ textAlign: 'center', padding: 50 }}>
                    <Spin size="large" />
                </div>
            ) : (
                <div>
                    {/* 当前设置预览 */}
                    <Card size="small" title="当前设置" style={{ marginBottom: 16 }}>
                        <Descriptions column={2} size="small">
                            <Descriptions.Item label="排序方式">
                                {currentReviewSettings ? SORT_BY_TEXT[currentReviewSettings.sortBy] : '-'}
                            </Descriptions.Item>
                            <Descriptions.Item label="最低评分">
                                {currentReviewSettings?.minScore || '-'} 分
                            </Descriptions.Item>
                            <Descriptions.Item label="显示匿名评价">
                                {currentReviewSettings?.showAnonymous ? '是' : '否'}
                            </Descriptions.Item>
                            <Descriptions.Item label="每页显示数量">
                                {currentReviewSettings?.pageSize || '-'} 条
                            </Descriptions.Item>
                        </Descriptions>
                    </Card>

                    {/* 设置表单 */}
                    <Form
                        form={reviewForm}
                        layout="vertical"
                        initialValues={DEFAULT_REVIEW_SETTINGS}
                        style={{ maxWidth: 600 }}
                    >
                        <Form.Item
                            name="sortBy"
                            label="评价排序规则"
                            tooltip="设置前端评价列表的默认排序方式"
                        >
                            <Radio.Group>
                                <Radio.Button value="time">按时间</Radio.Button>
                                <Radio.Button value="score">按评分</Radio.Button>
                                <Radio.Button value="likes">按点赞数</Radio.Button>
                            </Radio.Group>
                        </Form.Item>

                        <Form.Item
                            name="minScore"
                            label="最低评分阈值"
                            tooltip="低于此评分的评价将不在前端展示"
                            rules={[
                                { required: true, message: '请输入最低评分' },
                                { type: 'number', min: 1, max: 5, message: '评分范围为1-5' },
                            ]}
                        >
                            <InputNumber min={1} max={5} style={{ width: 120 }} suffix="分" />
                        </Form.Item>

                        <Form.Item
                            name="showAnonymous"
                            label="显示匿名评价"
                            tooltip="是否在前端展示匿名评价"
                            valuePropName="checked"
                        >
                            <Switch checkedChildren="显示" unCheckedChildren="隐藏" />
                        </Form.Item>

                        <Form.Item
                            name="pageSize"
                            label="每页显示数量"
                            tooltip="前端评价列表每页显示的评价数量"
                            rules={[
                                { required: true, message: '请输入每页数量' },
                                { type: 'number', min: 5, max: 50, message: '数量范围为5-50' },
                            ]}
                        >
                            <InputNumber min={5} max={50} style={{ width: 120 }} suffix="条" />
                        </Form.Item>

                        <Divider />

                        <Form.Item>
                            <Space>
                                <Button
                                    type="primary"
                                    icon={<SaveOutlined />}
                                    onClick={handleSaveReviewSettings}
                                    loading={reviewSaving}
                                    style={{ backgroundColor: '#5865F2' }}
                                >
                                    保存设置
                                </Button>
                                <Button icon={<UndoOutlined />} onClick={handleRevertReviewSettings}>
                                    恢复当前
                                </Button>
                                <Button icon={<ReloadOutlined />} onClick={handleResetReviewSettings}>
                                    重置为默认
                                </Button>
                            </Space>
                        </Form.Item>
                    </Form>

                    <Text type="secondary">
                        说明：这些设置将影响用户端评价的展示方式。修改后需要点击"保存设置"才能生效。
                    </Text>
                </div>
            ),
        },
        {
            key: '6',
            label: (
                <span>
                    <SyncOutlined />
                    系统初始化
                </span>
            ),
            children: (
                <Space direction="vertical" size="large" style={{ width: '100%' }}>
                    <Alert
                        message="系统初始化说明"
                        description={
                            <div>
                                <Paragraph style={{ marginBottom: 8 }}>
                                    系统初始化会执行以下操作：
                                </Paragraph>
                                <ul style={{ marginBottom: 0, paddingLeft: 20 }}>
                                    <li>同步前端菜单配置到后端数据库</li>
                                    <li>同步前端权限配置到后端数据库</li>
                                    <li>为超级管理员自动分配所有权限</li>
                                </ul>
                            </div>
                        }
                        type="info"
                        showIcon
                    />

                    <Card size="small" title="初始化状态">
                        <Space direction="vertical" style={{ width: '100%' }}>
                            <div>
                                <Text type="secondary">上次初始化时间：</Text>
                                <Text strong style={{ marginLeft: 8 }}>
                                    {formatLastInitTime()}
                                </Text>
                            </div>
                            <Button
                                type="primary"
                                size="large"
                                icon={<SyncOutlined spin={initializing} />}
                                onClick={handleInit}
                                loading={initializing}
                                style={{ backgroundColor: '#5865F2', width: '100%' }}
                            >
                                {initializing ? '正在初始化...' : '立即初始化'}
                            </Button>
                        </Space>
                    </Card>

                    <Alert
                        message="注意事项"
                        description={
                            <ul style={{ marginBottom: 0, paddingLeft: 20 }}>
                                <li>初始化过程可能需要几秒钟，请耐心等待</li>
                                <li>初始化不会影响现有的用户数据和业务数据</li>
                                <li>建议在系统升级或配置变更后执行初始化</li>
                                <li>普通管理员执行初始化不会影响超级管理员权限</li>
                            </ul>
                        }
                        type="warning"
                        showIcon
                    />
                </Space>
            ),
        },
    ];

    return (
        <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.3 }}
        >
            <Card variant="borderless" title="系统设置" styles={{ body: { paddingTop: 0 } }}>
                <Tabs defaultActiveKey="1" items={items} />
            </Card>
        </motion.div>
    );
};

export default Settings;
