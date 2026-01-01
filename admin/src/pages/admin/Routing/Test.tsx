/**
 * Routing Rule Test Page
 * Test tool for routing rules
 */
import React, { useState } from 'react';
import {
    Card,
    Form,
    Input,
    InputNumber,
    Select,
    Button,
    Space,
    Result,
    Descriptions,
    Tag,
    Alert,
    Row,
    Col,
    Statistic,
    Spin,
} from 'antd';
import {
    ExperimentOutlined,
    CheckCircleOutlined,
    CloseCircleOutlined,
    ReloadOutlined,
} from '@ant-design/icons';
import { routingApi } from '@/api/routing';
import type { RoutingTestResponse, RoutingCondition } from '@/api/routing';
import type { TestFormValues } from './types';

// Game type options
const GAME_TYPE_OPTIONS = [
    { label: '王者荣耀', value: 'honor_of_kings' },
    { label: '和平精英', value: 'pubg_mobile' },
    { label: '英雄联盟', value: 'league_of_legends' },
    { label: '绝地求生', value: 'pubg' },
    { label: '原神', value: 'genshin_impact' },
    { label: '永劫无间', value: 'naraka' },
];

// Service type options
const SERVICE_TYPE_OPTIONS = [
    { label: '陪练', value: 'companion' },
    { label: '代练', value: 'boosting' },
    { label: '陪玩', value: 'play_with' },
    { label: '教学', value: 'coaching' },
];

// Region options
const REGION_OPTIONS = [
    { label: '华东', value: 'east_china' },
    { label: '华南', value: 'south_china' },
    { label: '华北', value: 'north_china' },
    { label: '西南', value: 'southwest_china' },
    { label: '西北', value: 'northwest_china' },
    { label: '东北', value: 'northeast_china' },
];

/**
 * Get condition display text
 */
const getConditionDisplay = (condition: RoutingCondition): string => {
    const fieldLabels: Record<string, string> = {
        game_type: '游戏类型',
        service_type: '服务类型',
        order_amount: '订单金额',
        region: '地区',
    };

    const operatorLabels: Record<string, string> = {
        eq: '等于',
        neq: '不等于',
        in: '包含于',
        not_in: '不包含于',
        gt: '大于',
        lt: '小于',
        between: '介于',
    };

    const fieldLabel = fieldLabels[condition.field] || condition.field;
    const operatorLabel = operatorLabels[condition.operator] || condition.operator;

    let valueDisplay: string;
    if (condition.operator === 'between' && Array.isArray(condition.value)) {
        valueDisplay = `${condition.value[0]} ~ ${condition.value[1]} 元`;
    } else if (Array.isArray(condition.value)) {
        valueDisplay = condition.value.join(', ');
    } else if (condition.field === 'order_amount') {
        valueDisplay = `${condition.value} 元`;
    } else {
        valueDisplay = String(condition.value);
    }

    return `${fieldLabel} ${operatorLabel} ${valueDisplay}`;
};

/**
 * Get condition field label
 */
const getFieldLabel = (field: string): string => {
    const labels: Record<string, string> = {
        game_type: '游戏类型',
        service_type: '服务类型',
        order_amount: '订单金额',
        region: '地区',
    };
    return labels[field] || field;
};

/**
 * Routing Test Page
 */
const RoutingTestPage: React.FC = () => {
    const [form] = Form.useForm<TestFormValues>();
    const [testing, setTesting] = useState(false);
    const [testResult, setTestResult] = useState<RoutingTestResponse | null>(null);

    const handleTest = async () => {
        try {
            const values = await form.validateFields();
            setTesting(true);
            setTestResult(null);

            const response = await routingApi.testRouting({
                gameType: values.gameType,
                serviceType: values.serviceType,
                amountCents: values.amountCents ? values.amountCents * 100 : undefined,
                region: values.region,
            });

            if (response.data?.data) {
                setTestResult(response.data.data);
            }
        } catch {
            // Form validation error or API error
        } finally {
            setTesting(false);
        }
    };

    const handleReset = () => {
        form.resetFields();
        setTestResult(null);
    };

    const getGameTypeLabel = (value: string) => {
        return GAME_TYPE_OPTIONS.find(o => o.value === value)?.label || value;
    };

    const getServiceTypeLabel = (value: string) => {
        return SERVICE_TYPE_OPTIONS.find(o => o.value === value)?.label || value;
    };

    const getRegionLabel = (value: string) => {
        return REGION_OPTIONS.find(o => o.value === value)?.label || value;
    };

    return (
        <div style={{ padding: 24 }}>
            <Card
                title={
                    <Space>
                        <ExperimentOutlined />
                        <span>路由规则测试工具</span>
                    </Space>
                }
                extra={
                    <Button
                        icon={<ReloadOutlined />}
                        onClick={handleReset}
                        disabled={testing}
                    >
                        重置
                    </Button>
                }
            >
                <Row gutter={24}>
                    <Col span={testResult ? 12 : 24}>
                        <Form
                            form={form}
                            layout="vertical"
                            onFinish={handleTest}
                        >
                            <Form.Item
                                name="gameType"
                                label="游戏类型"
                                tooltip="模拟订单的游戏类型"
                            >
                                <Select
                                    placeholder="选择游戏类型"
                                    options={GAME_TYPE_OPTIONS}
                                    allowClear
                                />
                            </Form.Item>

                            <Form.Item
                                name="serviceType"
                                label="服务类型"
                                tooltip="模拟订单的服务类型"
                            >
                                <Select
                                    placeholder="选择服务类型"
                                    options={SERVICE_TYPE_OPTIONS}
                                    allowClear
                                />
                            </Form.Item>

                            <Form.Item
                                name="amountCents"
                                label="订单金额（元）"
                                tooltip="模拟订单的金额"
                            >
                                <InputNumber
                                    placeholder="输入订单金额"
                                    min={0}
                                    precision={2}
                                    style={{ width: '100%' }}
                                />
                            </Form.Item>

                            <Form.Item
                                name="region"
                                label="地区"
                                tooltip="模拟订单所属地区"
                            >
                                <Select
                                    placeholder="选择地区"
                                    options={REGION_OPTIONS}
                                    allowClear
                                />
                            </Form.Item>

                            <Form.Item>
                                <Space>
                                    <Button
                                        type="primary"
                                        htmlType="submit"
                                        icon={<ExperimentOutlined />}
                                        loading={testing}
                                    >
                                        测试路由
                                    </Button>
                                    <Button onClick={handleReset} disabled={testing}>
                                        重置
                                    </Button>
                                </Space>
                            </Form.Item>
                        </Form>
                    </Col>

                    <Col span={testResult ? 12 : 24}>
                        {testing && (
                            <div style={{ textAlign: 'center', padding: 40 }}>
                                <Spin tip="正在测试路由规则..." />
                            </div>
                        )}

                        {testResult && !testing && (
                            <div>
                                <Result
                                    icon={<CheckCircleOutlined style={{ color: '#52c41a' }} />}
                                    title="路由测试完成"
                                    subTitle="系统已匹配到合适的路由规则"
                                    style={{ padding: '20px 0' }}
                                />

                                <Card
                                    type="inner"
                                    title="测试结果"
                                    size="small"
                                    styles={{
                                        body: { padding: 16 }
                                    }}
                                >
                                    <Descriptions column={1} size="small" bordered>
                                        <Descriptions.Item label="匹配规则">
                                            {testResult.matchedRuleId ? (
                                                <Tag color="success">{testResult.matchedRuleName}</Tag>
                                            ) : (
                                                <Tag color="default">默认规则</Tag>
                                            )}
                                        </Descriptions.Item>
                                        <Descriptions.Item label="收款主体">
                                            <Tag color="blue">{testResult.entityName}</Tag>
                                        </Descriptions.Item>
                                        <Descriptions.Item label="商户号">
                                            <code>{testResult.merchantNo}</code>
                                        </Descriptions.Item>
                                        <Descriptions.Item label="是否默认">
                                            {testResult.isDefault ? (
                                                <Tag color="orange">默认主体</Tag>
                                            ) : (
                                                <Tag color="green">规则路由</Tag>
                                            )}
                                        </Descriptions.Item>
                                    </Descriptions>

                                    {testResult.matchDetails && testResult.matchDetails.length > 0 && (
                                        <div style={{ marginTop: 16 }}>
                                            <Alert
                                                message="匹配条件"
                                                type="info"
                                                showIcon
                                                style={{ marginBottom: 12 }}
                                            />
                                            <Space direction="vertical" size={8} style={{ width: '100%' }}>
                                                {testResult.matchDetails.map((cond, index) => (
                                                    <Tag key={index} color="cyan" style={{
                                                        padding: '4px 12px',
                                                        fontSize: 13
                                                    }}>
                                                        {getConditionDisplay(cond)}
                                                    </Tag>
                                                ))}
                                            </Space>
                                        </div>
                                    )}
                                </Card>

                                <Card
                                    type="inner"
                                    title="输入参数"
                                    size="small"
                                    styles={{
                                        body: { padding: 16 }
                                    }}
                                    style={{ marginTop: 12 }}
                                >
                                    <Row gutter={16}>
                                        <Col span={12}>
                                            <Statistic
                                                title="游戏类型"
                                                value={form.getFieldValue('gameType') ?
                                                    getGameTypeLabel(form.getFieldValue('gameType')) : '-'}
                                                valueStyle={{ fontSize: 14 }}
                                            />
                                        </Col>
                                        <Col span={12}>
                                            <Statistic
                                                title="服务类型"
                                                value={form.getFieldValue('serviceType') ?
                                                    getServiceTypeLabel(form.getFieldValue('serviceType')) : '-'}
                                                valueStyle={{ fontSize: 14 }}
                                            />
                                        </Col>
                                        <Col span={12}>
                                            <Statistic
                                                title="订单金额"
                                                value={form.getFieldValue('amountCents') || '-'}
                                                suffix="元"
                                                precision={2}
                                                valueStyle={{ fontSize: 14 }}
                                            />
                                        </Col>
                                        <Col span={12}>
                                            <Statistic
                                                title="地区"
                                                value={form.getFieldValue('region') ?
                                                    getRegionLabel(form.getFieldValue('region')) : '-'}
                                                valueStyle={{ fontSize: 14 }}
                                            />
                                        </Col>
                                    </Row>
                                </Card>
                            </div>
                        )}

                        {!testing && !testResult && (
                            <Alert
                                message="请输入测试参数并点击测试按钮"
                                type="info"
                                showIcon
                                style={{ marginTop: 24 }}
                            />
                        )}
                    </Col>
                </Row>
            </Card>

            <Card
                title="使用说明"
                style={{ marginTop: 16 }}
                size="small"
            >
                <Space direction="vertical" size={8}>
                    <div>1. 选择模拟订单的游戏类型、服务类型、金额和地区</div>
                    <div>2. 点击"测试路由"按钮进行规则匹配</div>
                    <div>3. 系统将按优先级顺序检查路由规则，返回匹配的主体</div>
                    <div>4. 若无规则匹配，将返回默认收款主体</div>
                    <div>5. 可以使用此工具验证路由规则配置是否正确</div>
                </Space>
            </Card>
        </div>
    );
};

export default RoutingTestPage;
