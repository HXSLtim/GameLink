import React, { useState, useEffect } from 'react';
import { Card, Form, Button, Tabs, Divider, App, Space, Typography, Alert, Radio, InputNumber, Spin, Descriptions, Switch, Tooltip, Select, Row, Col, Badge } from 'antd';
import { SaveOutlined, SyncOutlined, StarOutlined, ReloadOutlined, UndoOutlined, CheckCircleOutlined, CloseCircleOutlined } from '@ant-design/icons';
import { motion } from 'framer-motion';
import { forceInit } from '@/services/init';
import { syncApi, type InitStatusResponse } from '@/api/sync';
import { reviewSettingsApi } from '@/api/review';
import type { ReviewDisplaySettings, UpdateSettingsFormData } from '@/types/review';
import { SORT_BY_TEXT } from '@/types/review';
import { usePermissions } from '@/hooks/usePermission';

import { logger } from '@/utils/logger';
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
    const { message, modal } = App.useApp();
    const [reviewForm] = Form.useForm<UpdateSettingsFormData>();
    const [initializing, setInitializing] = useState(false);
    const [initStatus, setInitStatus] = useState<InitStatusResponse | null>(null);
    const [loadingStatus, setLoadingStatus] = useState(true);

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
            const response = await reviewSettingsApi.getSettings();
            if (response.data.success) {
                setCurrentReviewSettings(response.data.data);
                reviewForm.setFieldsValue({
                    sortBy: response.data.data.sortBy,
                    minScore: response.data.data.minScore,
                    showAnonymous: response.data.data.showAnonymous,
                    pageSize: response.data.data.pageSize,
                    autoApprove: response.data.data.autoApprove,
                    autoApproveMinRating: response.data.data.autoApproveMinRating,
                });
            }
        } catch {
            // 如果API不存在，使用默认值
            reviewForm.setFieldsValue(DEFAULT_REVIEW_SETTINGS);
        } finally {
            setReviewLoading(false);
        }
    };

    // 加载初始化状态（从后端数据库）
    const fetchInitStatus = async () => {
        setLoadingStatus(true);
        try {
            const status = await syncApi.getInitStatus();
            setInitStatus(status);
            logger.info('[Settings] Init status:', status);
        } catch (error) {
            logger.error('[Settings] Failed to get init status:', error);
            setInitStatus({
                initialized: false,
                message: '无法获取初始化状态',
            });
        } finally {
            setLoadingStatus(false);
        }
    };

    useEffect(() => {
        fetchReviewSettings();
        fetchInitStatus();
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []); // 初始化时只执行一次，fetchReviewSettings 和 fetchInitStatus 在组件内部定义

    // 保存评价设置
    const handleSaveReviewSettings = async () => {
        try {
            const values = await reviewForm.validateFields();
            setReviewSaving(true);
            const response = await reviewSettingsApi.updateSettings(values);
            if (response.data.success) {
                message.success('评价设置已保存');
                setCurrentReviewSettings(response.data.data);
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
        modal.confirm({
            title: '确认同步',
            content: '此操作将同步菜单和权限配置到数据库，并为超级管理员分配所有权限。确定要继续吗？',
            okText: '确认同步',
            cancelText: '取消',
            onOk: async () => {
                setInitializing(true);
                const hide = message.loading('正在初始化系统...', 0);

                try {
                    const result = await forceInit({ verbose: true });

                    if (result.success) {
                        message.success({
                            content: `初始化成功！耗时 ${result.duration}ms`,
                            duration: 3,
                        });

                        // 重新加载初始化状态
                        await fetchInitStatus();

                        // 显示详细信息
                        if (result.menuSync) {
                            logger.info('[菜单同步]', result.menuSync);
                        }
                        if (result.permissionSync) {
                            logger.info('[权限同步]', result.permissionSync);
                        }
                        if (result.superAdminAssign) {
                            logger.info('[超管权限]', result.superAdminAssign);
                        }
                    } else {
                        message.error({
                            content: `初始化失败：${result.errors.join(', ')}`,
                            duration: 5,
                        });
                    }
                } catch (error) {
                    logger.error('初始化异常:', error);
                    message.error({
                        content: `初始化异常：${error instanceof Error ? error.message : '未知错误'}`,
                        duration: 5,
                    });
                } finally {
                    hide();
                    setInitializing(false);
                }
            },
        });
    };

    /**
     * 格式化时间显示
     */
    const formatLastInitTime = () => {
        if (!initStatus?.lastSyncAt) return '从未初始化';
        const date = new Date(initStatus.lastSyncAt);
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
                        {loadingStatus ? (
                            <div style={{ textAlign: 'center', padding: 20 }}>
                                <Spin />
                            </div>
                        ) : (
                            <Row gutter={[24, 24]} align="middle">
                                <Col xs={24} sm={8}>
                                    <div>
                                        <Text type="secondary" style={{ fontSize: 13 }}>初始化状态</Text>
                                        <div style={{ marginTop: 8, display: 'flex', alignItems: 'center', gap: 8 }}>
                                            {initStatus?.initialized ? (
                                                <>
                                                    <CheckCircleOutlined style={{ color: '#52c41a', fontSize: 20 }} />
                                                    <Text strong style={{ color: '#52c41a' }}>已初始化</Text>
                                                </>
                                            ) : (
                                                <>
                                                    <CloseCircleOutlined style={{ color: '#ff4d4f', fontSize: 20 }} />
                                                    <Text strong style={{ color: '#ff4d4f' }}>未初始化</Text>
                                                </>
                                            )}
                                        </div>
                                    </div>
                                </Col>
                                <Col xs={24} sm={8}>
                                    <div>
                                        <Text type="secondary" style={{ fontSize: 13 }}>上次同步时间</Text>
                                        <div style={{ fontSize: 14, fontWeight: 500, marginTop: 8 }}>
                                            {formatLastInitTime()}
                                        </div>
                                    </div>
                                </Col>
                                <Col xs={24} sm={8}>
                                    <div>
                                        <Text type="secondary" style={{ fontSize: 13 }}>数据统计</Text>
                                        <div style={{ marginTop: 8, display: 'flex', gap: 16 }}>
                                            <div>
                                                <Badge count={initStatus?.menuCount || 0} showZero style={{ backgroundColor: '#5865F2' }}>
                                                    <Text>菜单</Text>
                                                </Badge>
                                            </div>
                                            <div>
                                                <Badge count={initStatus?.permissionCount || 0} showZero style={{ backgroundColor: '#52c41a' }}>
                                                    <Text>权限</Text>
                                                </Badge>
                                            </div>
                                        </div>
                                    </div>
                                </Col>
                            </Row>
                        )}
                    </Card>

                    {/* 初始化按钮 */}
                    <Card style={{ marginBottom: 16, borderRadius: 12 }}>
                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: 16 }}>
                            <div>
                                <Text strong>手动同步菜单和权限</Text>
                                <div>
                                    <Text type="secondary" style={{ fontSize: 13 }}>
                                        当前端配置更新后，点击此按钮同步到后端数据库
                                    </Text>
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
                                    {initializing ? '正在同步...' : '立即同步'}
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
                                        立即同步
                                    </Button>
                                </Tooltip>
                            )}
                        </div>
                    </Card>

                    {/* 初始化说明 */}
                    <Card
                        size="small"
                        title={<span><Text strong>同步操作说明</Text></span>}
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
                                <li>同步过程通常需要 2-5 秒，请耐心等待</li>
                                <li>同步不会影响现有的用户数据和业务数据</li>
                                <li>同步不会影响现有角色的权限分配（权限ID不变）</li>
                                <li>建议在系统升级或配置变更后执行同步</li>
                                <li>超级管理员和普通管理员会自动获得所有新权限</li>
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
