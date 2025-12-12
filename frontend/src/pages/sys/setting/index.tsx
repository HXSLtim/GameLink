import React, { useState, useEffect } from 'react';
import { Card, Form, Button, Tabs, Divider, message, Space, Typography, Alert, Radio, InputNumber, Spin, Descriptions, Switch, Tooltip, Select } from 'antd';
import { SaveOutlined, SyncOutlined, StarOutlined, ReloadOutlined, UndoOutlined } from '@ant-design/icons';
import { motion } from 'framer-motion';
import { forceInit } from '@/services/init';
import { reviewSettingsApi } from '@/api/review';
import type { ReviewDisplaySettings, UpdateSettingsFormData } from '@/types/review';
import { SORT_BY_TEXT } from '@/types/review';
import { usePermissions } from '@/hooks/usePermission';

const { Text, Paragraph } = Typography;

// 评价设置默认值
const DEFAULT_REVIEW_SETTINGS: UpdateSettingsFormData = {
    sortBy: 'time',
    minScore: 1,
    showAnonymous: true,
    pageSize: 10,
    autoApprove: false,
    autoApproveMinRating: 4,
};

const Settings: React.FC = () => {
    const [reviewForm] = Form.useForm<UpdateSettingsFormData>();
    const [initializing, setInitializing] = useState(false);
    const [lastInitTime, setLastInitTime] = useState<string | null>(
        localStorage.getItem('app_init_timestamp')
    );

    // 权限检查
    const permissions = usePermissions({
        canUpdateReviewSettings: 'admin.review-settings.update',
        canInit: 'admin.system.init',
    });

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
                    autoApprove: response.data.autoApprove,
                    autoApproveMinRating: response.data.autoApproveMinRating,
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
                autoApprove: currentReviewSettings.autoApprove,
                autoApproveMinRating: currentReviewSettings.autoApproveMinRating,
            });
            message.info('已恢复为当前保存的设置');
        }
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
                            <Descriptions.Item label="自动批准">
                                {currentReviewSettings?.autoApprove ? '开启' : '关闭'}
                            </Descriptions.Item>
                            <Descriptions.Item label="自动批准最低评分">
                                {currentReviewSettings?.autoApproveMinRating || '-'} 星
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

                        <Divider>审核设置</Divider>

                        <Form.Item
                            name="autoApprove"
                            label="自动批准评价"
                            tooltip="开启后，符合条件的评价将自动批准，无需人工审核"
                            valuePropName="checked"
                        >
                            <Switch checkedChildren="开启" unCheckedChildren="关闭" />
                        </Form.Item>

                        <Form.Item
                            noStyle
                            shouldUpdate={(prevValues, currentValues) => prevValues.autoApprove !== currentValues.autoApprove}
                        >
                            {({ getFieldValue }) => 
                                getFieldValue('autoApprove') ? (
                                    <Form.Item
                                        name="autoApproveMinRating"
                                        label="自动批准最低评分"
                                        tooltip="只有评分大于等于此值的评价才会自动批准"
                                        rules={[
                                            { required: true, message: '请选择最低评分' },
                                        ]}
                                    >
                                        <Select style={{ width: 120 }}>
                                            <Select.Option value={1}>1 星及以上</Select.Option>
                                            <Select.Option value={2}>2 星及以上</Select.Option>
                                            <Select.Option value={3}>3 星及以上</Select.Option>
                                            <Select.Option value={4}>4 星及以上</Select.Option>
                                            <Select.Option value={5}>仅 5 星</Select.Option>
                                        </Select>
                                    </Form.Item>
                                ) : null
                            }
                        </Form.Item>

                        <Divider />

                        <Form.Item>
                            <Space>
                                {permissions.canUpdateReviewSettings ? (
                                    <Button
                                        type="primary"
                                        icon={<SaveOutlined />}
                                        onClick={handleSaveReviewSettings}
                                        loading={reviewSaving}
                                        style={{ backgroundColor: '#5865F2' }}
                                    >
                                        保存设置
                                    </Button>
                                ) : (
                                    <Tooltip title="无修改权限">
                                        <Button type="primary" icon={<SaveOutlined />} disabled>
                                            保存设置
                                        </Button>
                                    </Tooltip>
                                )}
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
            key: '2',
            label: (
                <span>
                    <SyncOutlined />
                    系统初始化
                </span>
            ),
            children: (
                <div style={{ maxWidth: 800 }}>
                    {/* 初始化状态卡片 */}
                    <Card 
                        style={{ marginBottom: 24, borderRadius: 12 }}
                        styles={{ body: { padding: 24 } }}
                    >
                        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', flexWrap: 'wrap', gap: 16 }}>
                            <div>
                                <Text type="secondary" style={{ fontSize: 13 }}>上次初始化时间</Text>
                                <div style={{ fontSize: 18, fontWeight: 500, marginTop: 4 }}>
                                    {formatLastInitTime()}
                                </div>
                            </div>
                            {permissions.canInit ? (
                                <Button
                                    type="primary"
                                    size="large"
                                    icon={<SyncOutlined spin={initializing} />}
                                    onClick={handleInit}
                                    loading={initializing}
                                    style={{ 
                                        backgroundColor: '#5865F2', 
                                        height: 48,
                                        paddingLeft: 32,
                                        paddingRight: 32,
                                        borderRadius: 8,
                                        fontSize: 15,
                                    }}
                                >
                                    {initializing ? '正在初始化...' : '立即初始化'}
                                </Button>
                            ) : (
                                <Tooltip title="无初始化权限">
                                    <Button
                                        type="primary"
                                        size="large"
                                        icon={<SyncOutlined />}
                                        disabled
                                        style={{ 
                                            height: 48,
                                            paddingLeft: 32,
                                            paddingRight: 32,
                                            borderRadius: 8,
                                            fontSize: 15,
                                        }}
                                    >
                                        立即初始化
                                    </Button>
                                </Tooltip>
                            )}
                        </div>
                    </Card>

                    {/* 初始化说明 */}
                    <Card 
                        size="small" 
                        title={<span><Text strong>初始化操作说明</Text></span>}
                        style={{ marginBottom: 16, borderRadius: 12 }}
                    >
                        <Space direction="vertical" size={12} style={{ width: '100%' }}>
                            <div style={{ display: 'flex', alignItems: 'flex-start', gap: 12 }}>
                                <div style={{ 
                                    width: 24, height: 24, borderRadius: '50%', 
                                    backgroundColor: 'rgba(88, 101, 242, 0.1)', 
                                    display: 'flex', alignItems: 'center', justifyContent: 'center',
                                    flexShrink: 0, marginTop: 2
                                }}>
                                    <Text style={{ color: '#5865F2', fontSize: 12, fontWeight: 600 }}>1</Text>
                                </div>
                                <div>
                                    <Text strong>同步菜单配置</Text>
                                    <Paragraph type="secondary" style={{ margin: 0, fontSize: 13 }}>
                                        将前端定义的菜单结构同步到后端数据库
                                    </Paragraph>
                                </div>
                            </div>
                            <div style={{ display: 'flex', alignItems: 'flex-start', gap: 12 }}>
                                <div style={{ 
                                    width: 24, height: 24, borderRadius: '50%', 
                                    backgroundColor: 'rgba(88, 101, 242, 0.1)', 
                                    display: 'flex', alignItems: 'center', justifyContent: 'center',
                                    flexShrink: 0, marginTop: 2
                                }}>
                                    <Text style={{ color: '#5865F2', fontSize: 12, fontWeight: 600 }}>2</Text>
                                </div>
                                <div>
                                    <Text strong>同步权限配置</Text>
                                    <Paragraph type="secondary" style={{ margin: 0, fontSize: 13 }}>
                                        将前端定义的权限码同步到后端数据库
                                    </Paragraph>
                                </div>
                            </div>
                            <div style={{ display: 'flex', alignItems: 'flex-start', gap: 12 }}>
                                <div style={{ 
                                    width: 24, height: 24, borderRadius: '50%', 
                                    backgroundColor: 'rgba(88, 101, 242, 0.1)', 
                                    display: 'flex', alignItems: 'center', justifyContent: 'center',
                                    flexShrink: 0, marginTop: 2
                                }}>
                                    <Text style={{ color: '#5865F2', fontSize: 12, fontWeight: 600 }}>3</Text>
                                </div>
                                <div>
                                    <Text strong>分配超管权限</Text>
                                    <Paragraph type="secondary" style={{ margin: 0, fontSize: 13 }}>
                                        为超级管理员自动分配所有系统权限
                                    </Paragraph>
                                </div>
                            </div>
                        </Space>
                    </Card>

                    {/* 注意事项 */}
                    <Alert
                        message="注意事项"
                        description={
                            <ul style={{ marginBottom: 0, paddingLeft: 20, lineHeight: 2 }}>
                                <li>初始化过程通常需要 2-5 秒，请耐心等待</li>
                                <li>初始化不会影响现有的用户数据和业务数据</li>
                                <li>建议在系统升级或配置变更后执行初始化</li>
                            </ul>
                        }
                        type="info"
                        showIcon
                        style={{ borderRadius: 8 }}
                    />
                </div>
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
