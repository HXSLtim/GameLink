/**
 * 佣金管理页面
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Card,
    Row,
    Col,
    Statistic,
    Button,
    Modal,
    Form,
    Input,
    InputNumber,
    Select,
    message,
    DatePicker,
    Table,
    Space,
    Typography,
    Divider,
    Alert,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    DollarOutlined,
    PlusOutlined,
    SettingOutlined,
    SyncOutlined,
} from '@ant-design/icons';
import { PageContainer } from '@/components';
import { COMMISSION_PERMISSIONS } from '@/constants/permissions';
import { PermissionGuard } from '@/components/PermissionGuard';
import { adminApi, type PlatformStats, type CreateCommissionRuleDto } from '@/api/admin';
import dayjs from 'dayjs';

const { Title, Text } = Typography;

interface CommissionRule {
    id: number;
    name: string;
    ratePercent: number;
    minAmountCents?: number;
    maxAmountCents?: number;
    gameId?: number;
    isDefault: boolean;
    status: 'active' | 'inactive';
    createdAt: string;
}

/**
 * 佣金管理页面
 */
const CommissionPage: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [stats, setStats] = useState<PlatformStats | null>(null);
    const [selectedMonth, setSelectedMonth] = useState(dayjs().format('YYYY-MM'));
    // TODO: Load rules from API when backend supports listing
    const rules: CommissionRule[] = [];

    // 弹窗状态
    const [ruleModalVisible, setRuleModalVisible] = useState(false);
    const [currentRule, setCurrentRule] = useState<CommissionRule | null>(null);
    const [form] = Form.useForm();
    const [submitting, setSubmitting] = useState(false);
    const [triggeringSettlement, setTriggeringSettlement] = useState(false);

    /**
     * 加载平台统计
     */
    const loadStats = useCallback(async () => {
        setLoading(true);
        try {
            const response = await adminApi.getPlatformStats(selectedMonth);
            if (response.data.success) {
                setStats(response.data.data);
            }
        } catch (error) {
            console.error('Load stats error:', error);
            message.error('加载统计数据失败');
        } finally {
            setLoading(false);
        }
    }, [selectedMonth]);

    useEffect(() => {
        loadStats();
    }, [loadStats]);

    /**
     * 触发结算
     */
    const handleTriggerSettlement = async () => {
        Modal.confirm({
            title: '确认触发结算',
            content: `确定要触发 ${selectedMonth} 的月度结算吗？此操作将计算所有陪玩师的收益。`,
            onOk: async () => {
                try {
                    setTriggeringSettlement(true);
                    await adminApi.triggerSettlement(selectedMonth);
                    message.success('结算任务已触发');
                    loadStats();
                } catch (error) {
                    console.error('Trigger settlement error:', error);
                    message.error('触发结算失败');
                } finally {
                    setTriggeringSettlement(false);
                }
            },
        });
    };

    /**
     * 创建佣金规则
     */
    const handleCreateRule = () => {
        setCurrentRule(null);
        form.resetFields();
        form.setFieldsValue({ ratePercent: 20, isDefault: false });
        setRuleModalVisible(true);
    };

    /**
     * 编辑佣金规则
     */
    const handleEditRule = (rule: CommissionRule) => {
        setCurrentRule(rule);
        form.setFieldsValue({
            name: rule.name,
            ratePercent: rule.ratePercent,
            minAmountCents: rule.minAmountCents ? rule.minAmountCents / 100 : undefined,
            maxAmountCents: rule.maxAmountCents ? rule.maxAmountCents / 100 : undefined,
            isDefault: rule.isDefault,
        });
        setRuleModalVisible(true);
    };

    /**
     * 保存佣金规则
     */
    const handleSaveRule = async () => {
        try {
            const values = await form.validateFields();
            setSubmitting(true);

            const data: CreateCommissionRuleDto = {
                name: values.name,
                ratePercent: values.ratePercent,
                minAmountCents: values.minAmountCents ? Math.round(values.minAmountCents * 100) : undefined,
                maxAmountCents: values.maxAmountCents ? Math.round(values.maxAmountCents * 100) : undefined,
                isDefault: values.isDefault,
            };

            if (currentRule) {
                await adminApi.updateCommissionRule(currentRule.id, data);
                message.success('更新成功');
            } else {
                await adminApi.createCommissionRule(data);
                message.success('创建成功');
            }
            setRuleModalVisible(false);
            // 刷新规则列表（如果有API的话）
        } catch (error) {
            console.error('Save rule error:', error);
            message.error('保存失败');
        } finally {
            setSubmitting(false);
        }
    };

    /**
     * 规则表格列
     */
    const ruleColumns: ColumnsType<CommissionRule> = [
        { title: '规则名称', dataIndex: 'name', key: 'name' },
        {
            title: '抽成比例',
            dataIndex: 'ratePercent',
            key: 'ratePercent',
            render: (rate: number) => `${rate}%`,
        },
        {
            title: '最低金额',
            dataIndex: 'minAmountCents',
            key: 'minAmountCents',
            render: (cents?: number) => cents ? `¥${(cents / 100).toFixed(2)}` : '-',
        },
        {
            title: '最高金额',
            dataIndex: 'maxAmountCents',
            key: 'maxAmountCents',
            render: (cents?: number) => cents ? `¥${(cents / 100).toFixed(2)}` : '-',
        },
        {
            title: '默认规则',
            dataIndex: 'isDefault',
            key: 'isDefault',
            render: (isDefault: boolean) => isDefault ? '是' : '否',
        },
        {
            title: '操作',
            key: 'action',
            fixed: 'right',
            width: 100,
            render: (_, record) => (
                <PermissionGuard permission={COMMISSION_PERMISSIONS.UPDATE}>
                    <Button type="link" size="small" onClick={() => handleEditRule(record)}>
                        编辑
                    </Button>
                </PermissionGuard>
            ),
        },
    ];

    // 默认规则数据（实际应从API获取）
    const defaultRules: CommissionRule[] = [
        {
            id: 1,
            name: '默认抽成规则',
            ratePercent: 20,
            isDefault: true,
            status: 'active',
            createdAt: '2024-01-01',
        },
    ];

    return (
        <PageContainer title="佣金管理" subTitle="平台佣金设置与结算管理">
            {/* 月份选择 */}
            <Card style={{ marginBottom: 16 }}>
                <Space>
                    <Text>选择月份：</Text>
                    <DatePicker
                        picker="month"
                        value={dayjs(selectedMonth)}
                        onChange={(date) => setSelectedMonth(date?.format('YYYY-MM') || dayjs().format('YYYY-MM'))}
                        allowClear={false}
                    />
                    <Button icon={<SyncOutlined />} onClick={loadStats} loading={loading}>
                        刷新
                    </Button>
                    <PermissionGuard permission={COMMISSION_PERMISSIONS.UPDATE}>
                        <Button
                            type="primary"
                            icon={<SettingOutlined />}
                            onClick={handleTriggerSettlement}
                            loading={triggeringSettlement}
                        >
                            触发月度结算
                        </Button>
                    </PermissionGuard>
                </Space>
            </Card>

            {/* 统计卡片 */}
            <Row gutter={16} style={{ marginBottom: 16 }}>
                <Col span={6}>
                    <Card>
                        <Statistic
                            title="平台总收入"
                            value={stats?.totalRevenueCents ? stats.totalRevenueCents / 100 : 0}
                            precision={2}
                            prefix="¥"
                            valueStyle={{ color: '#3f8600' }}
                        />
                    </Card>
                </Col>
                <Col span={6}>
                    <Card>
                        <Statistic
                            title="平台佣金"
                            value={stats?.totalCommissionCents ? stats.totalCommissionCents / 100 : 0}
                            precision={2}
                            prefix="¥"
                            valueStyle={{ color: '#1890ff' }}
                        />
                    </Card>
                </Col>
                <Col span={6}>
                    <Card>
                        <Statistic
                            title="总订单数"
                            value={stats?.totalOrderCount || 0}
                        />
                    </Card>
                </Col>
                <Col span={6}>
                    <Card>
                        <Statistic
                            title="完成订单数"
                            value={stats?.completedOrderCount || 0}
                            valueStyle={{ color: '#52c41a' }}
                        />
                    </Card>
                </Col>
            </Row>

            {/* 佣金规则 */}
            <Card
                title={
                    <Space>
                        <DollarOutlined />
                        <span>佣金规则</span>
                    </Space>
                }
                extra={
                    <PermissionGuard permission={COMMISSION_PERMISSIONS.CREATE}>
                        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreateRule}>
                            新增规则
                        </Button>
                    </PermissionGuard>
                }
            >
                <Alert
                    message="佣金说明"
                    description="平台从每笔完成的订单中按比例抽取佣金。默认抽成比例为20%，可根据游戏类型或陪玩师等级设置不同的抽成规则。"
                    type="info"
                    showIcon
                    style={{ marginBottom: 16 }}
                />
                <Table
                    columns={ruleColumns}
                    dataSource={rules.length > 0 ? rules : defaultRules}
                    rowKey="id"
                    pagination={false}
                    scroll={{ x: 1000 }}
                />
            </Card>

            <Divider />

            {/* 结算说明 */}
            <Card title="结算说明">
                <Row gutter={24}>
                    <Col span={12}>
                        <Title level={5}>结算周期</Title>
                        <ul>
                            <li>结算周期为每月一次</li>
                            <li>每月1日自动触发上月结算</li>
                            <li>也可手动触发指定月份的结算</li>
                        </ul>
                    </Col>
                    <Col span={12}>
                        <Title level={5}>结算规则</Title>
                        <ul>
                            <li>只结算已完成的订单</li>
                            <li>退款订单不计入结算</li>
                            <li>结算金额 = 订单金额 × (1 - 抽成比例)</li>
                        </ul>
                    </Col>
                </Row>
            </Card>

            {/* 规则编辑弹窗 */}
            <Modal
                title={currentRule ? '编辑佣金规则' : '新增佣金规则'}
                open={ruleModalVisible}
                onOk={handleSaveRule}
                onCancel={() => setRuleModalVisible(false)}
                confirmLoading={submitting}
                width={500}
            >
                <Form form={form} layout="vertical">
                    <Form.Item
                        name="name"
                        label="规则名称"
                        rules={[{ required: true, message: '请输入规则名称' }]}
                    >
                        <Input placeholder="请输入规则名称" />
                    </Form.Item>

                    <Form.Item
                        name="ratePercent"
                        label="抽成比例 (%)"
                        rules={[{ required: true, message: '请输入抽成比例' }]}
                    >
                        <InputNumber
                            min={0}
                            max={100}
                            precision={1}
                            style={{ width: '100%' }}
                            placeholder="请输入抽成比例"
                            addonAfter="%"
                        />
                    </Form.Item>

                    <Form.Item name="minAmountCents" label="最低订单金额 (元)">
                        <InputNumber
                            min={0}
                            precision={2}
                            style={{ width: '100%' }}
                            placeholder="不限制请留空"
                            addonBefore="¥"
                        />
                    </Form.Item>

                    <Form.Item name="maxAmountCents" label="最高订单金额 (元)">
                        <InputNumber
                            min={0}
                            precision={2}
                            style={{ width: '100%' }}
                            placeholder="不限制请留空"
                            addonBefore="¥"
                        />
                    </Form.Item>

                    <Form.Item name="isDefault" label="是否默认规则">
                        <Select>
                            <Select.Option value={true}>是</Select.Option>
                            <Select.Option value={false}>否</Select.Option>
                        </Select>
                    </Form.Item>
                </Form>
            </Modal>
        </PageContainer>
    );
};

export default CommissionPage;
