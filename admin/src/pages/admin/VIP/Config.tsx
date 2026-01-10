/**
 * VIP配置管理页面
 * 管理VIP系统的全局配置
 */
import React, { useState, useEffect } from 'react';
import {
    Card,
    Form,
    InputNumber,
    Button,
    Space,
    message,
    Typography,
    Row,
    Col,
    Statistic,
    Descriptions,
    Divider,
} from 'antd';
import {
    SaveOutlined,
    ReloadOutlined,
    SettingOutlined,
} from '@ant-design/icons';
import { vipApi, VIP_CONFIG_KEYS } from '@/api/vip';
import type { VIPConfig } from '@/api/vip';

import { logger } from '@/utils/logger';
const { Title, Text } = Typography;

interface ConfigFormData {
    unlock_by_consume: number;
    unlock_by_recharge: number;
    expire_days: number;
}

const VIPConfigPage: React.FC = () => {
    const [form] = Form.useForm<ConfigFormData>();
    const [loading, setLoading] = useState(false);
    const [saving, setSaving] = useState(false);
    const [configs, setConfigs] = useState<Record<string, VIPConfig>>({});

    const loadConfigs = async () => {
        setLoading(true);
        try {
            const response = await vipApi.getVIPConfigs();
            if (response.data.success) {
                const configMap: Record<string, VIPConfig> = {};
                response.data.data.forEach((config) => {
                    configMap[config.configKey] = config;
                });
                setConfigs(configMap);

                // 设置表单值
                form.setFieldsValue({
                    unlock_by_consume: parseInt(configMap[VIP_CONFIG_KEYS.UNLOCK_BY_CONSUME]?.configValue || '0'),
                    unlock_by_recharge: parseInt(configMap[VIP_CONFIG_KEYS.UNLOCK_BY_RECHARGE]?.configValue || '0'),
                    expire_days: parseInt(configMap[VIP_CONFIG_KEYS.EXPIRE_DAYS]?.configValue || '0'),
                });
            } else {
                message.error('加载配置失败');
            }
        } catch (error) {
            logger.error('Load VIP configs error:', error);
            message.error('加载配置失败');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        loadConfigs();
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    const handleSubmit = async (values: ConfigFormData) => {
        setSaving(true);
        try {
            // 保存解锁门槛
            await vipApi.saveVIPConfig({
                configKey: VIP_CONFIG_KEYS.UNLOCK_BY_CONSUME,
                configValue: String(values.unlock_by_consume),
                description: '累计消费解锁门槛（分）',
            });

            await vipApi.saveVIPConfig({
                configKey: VIP_CONFIG_KEYS.UNLOCK_BY_RECHARGE,
                configValue: String(values.unlock_by_recharge),
                description: '累计充值解锁门槛（分）',
            });

            await vipApi.saveVIPConfig({
                configKey: VIP_CONFIG_KEYS.EXPIRE_DAYS,
                configValue: String(values.expire_days),
                description: 'VIP过期天数（0=永久）',
            });

            message.success('保存成功');
            loadConfigs();
        } catch (error) {
            logger.error('Save VIP config error:', error);
            message.error('保存失败');
        } finally {
            setSaving(false);
        }
    };

    return (
        <div style={{ padding: 24 }}>
            <div style={{ marginBottom: 24, display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                <Title level={4} style={{ margin: 0 }}>
                    <SettingOutlined /> VIP系统配置
                </Title>
                <Button icon={<ReloadOutlined />} onClick={loadConfigs} loading={loading}>
                    刷新
                </Button>
            </div>

            <Card
                title="解锁门槛配置"
                extra={<Text type="secondary">设置用户解锁VIP功能的条件</Text>}
                style={{ marginBottom: 16 }}
            >
                <Row gutter={16}>
                    <Col span={8}>
                        <Statistic
                            title="当前消费门槛"
                            value={configs[VIP_CONFIG_KEYS.UNLOCK_BY_CONSUME]?.configValue || '0'}
                            suffix="分"
                            precision={0}
                            prefix=""
                        />
                    </Col>
                    <Col span={8}>
                        <Statistic
                            title="当前充值门槛"
                            value={configs[VIP_CONFIG_KEYS.UNLOCK_BY_RECHARGE]?.configValue || '0'}
                            suffix="分"
                            precision={0}
                            prefix=""
                        />
                    </Col>
                    <Col span={8}>
                        <Statistic
                            title="VIP有效期"
                            value={configs[VIP_CONFIG_KEYS.EXPIRE_DAYS]?.configValue || '0'}
                            suffix="天"
                            precision={0}
                            prefix=""
                        />
                        <Text type="secondary" style={{ fontSize: 12 }}>
                            {configs[VIP_CONFIG_KEYS.EXPIRE_DAYS]?.configValue === '0' ? '(永久有效)' : ''}
                        </Text>
                    </Col>
                </Row>
            </Card>

            <Card title="配置参数">
                <Form
                    form={form}
                    layout="vertical"
                    onFinish={handleSubmit}
                >
                    <Row gutter={16}>
                        <Col span={8}>
                            <Form.Item
                                name="unlock_by_consume"
                                label="累计消费解锁门槛"
                                extra="单位：分，100分 = 1元"
                                rules={[{ required: true, message: '请输入消费门槛' }]}
                            >
                                <InputNumber
                                    min={0}
                                    style={{ width: '100%' }}
                                    placeholder="0"
                                    formatter={(value) => `¥ ${(value || 0) / 100}`}
                                    parser={(value: string | undefined) => {
                                        const num = parseFloat(value?.replace(/[¥\s]/g, '') || '0');
                                        return Math.round(num * 100);
                                    }}
                                />
                            </Form.Item>
                        </Col>
                        <Col span={8}>
                            <Form.Item
                                name="unlock_by_recharge"
                                label="累计充值解锁门槛"
                                extra="单位：分，100分 = 1元"
                                rules={[{ required: true, message: '请输入充值门槛' }]}
                            >
                                <InputNumber
                                    min={0}
                                    style={{ width: '100%' }}
                                    placeholder="0"
                                    formatter={(value) => `¥ ${(value || 0) / 100}`}
                                    parser={(value: string | undefined) => {
                                        const num = parseFloat(value?.replace(/[¥\s]/g, '') || '0');
                                        return Math.round(num * 100);
                                    }}
                                />
                            </Form.Item>
                        </Col>
                        <Col span={8}>
                            <Form.Item
                                name="expire_days"
                                label="VIP过期天数"
                                extra="0 表示永久有效"
                                rules={[{ required: true, message: '请输入过期天数' }]}
                            >
                                <InputNumber
                                    min={0}
                                    style={{ width: '100%' }}
                                    placeholder="0"
                                />
                            </Form.Item>
                        </Col>
                    </Row>

                    <Form.Item>
                        <Space>
                            <Button type="primary" htmlType="submit" icon={<SaveOutlined />} loading={saving}>
                                保存配置
                            </Button>
                            <Button onClick={() => form.resetFields()}>
                                重置
                            </Button>
                        </Space>
                    </Form.Item>
                </Form>
            </Card>

            <Divider />

            <Card title="配置说明" type="inner">
                <Descriptions bordered column={1}>
                    <Descriptions.Item label="累计消费解锁门槛">
                        用户累计消费金额达到此门槛后，可解锁VIP功能。单位为分（100分 = 1元）。
                    </Descriptions.Item>
                    <Descriptions.Item label="累计充值解锁门槛">
                        用户累计充值金额达到此门槛后，可解锁VIP功能。单位为分（100分 = 1元）。
                    </Descriptions.Item>
                    <Descriptions.Item label="VIP过期天数">
                        VIP会员的有效期，单位为天。设置为0时表示VIP永久有效。
                    </Descriptions.Item>
                </Descriptions>
            </Card>
        </div>
    );
};

export default VIPConfigPage;
