/**
 * Routing Form Component
 * Form for creating/editing routing rules
 */
import React, { useEffect, useState } from 'react';
import {
    Modal,
    Form,
    Input,
    InputNumber,
    Select,
    Switch,
    Space,
    Divider,
    message,
} from 'antd';
import type { RoutingRule, RoutingCondition } from '@/api/routing';
import { routingApi } from '@/api/routing';
import { settlementApi, type SettlementCompany } from '@/api/settlement';
import ConditionBuilder from './ConditionBuilder';
import type {
    ConditionFormItem,
    RoutingRuleFormValues,
} from '../types';

interface RoutingFormProps {
    visible: boolean;
    rule: RoutingRule | null;
    onCancel: () => void;
    onOk: () => void;
}

/**
 * Convert form conditions to API conditions
 */
const toApiConditions = (formConditions: ConditionFormItem[]) => {
    return formConditions.map(({ id: _id, ...rest }) => rest);
};

/**
 * Convert API conditions to form conditions
 */
const fromApiConditions = (apiConditions: RoutingCondition[]): ConditionFormItem[] => {
    return apiConditions.map((cond, index) => ({
        id: `api-${index}`,
        field: cond.field,
        operator: cond.operator,
        value: cond.value,
    }));
};

/**
 * Routing Form Component
 */
const RoutingForm: React.FC<RoutingFormProps> = ({
    visible,
    rule,
    onCancel,
    onOk,
}) => {
    const [form] = Form.useForm<RoutingRuleFormValues>();
    const [submitting, setSubmitting] = useState(false);
    const [loadingEntities, setLoadingEntities] = useState(false);
    const [collectionEntities, setCollectionEntities] = useState<SettlementCompany[]>([]);

    // Load collection entities (settlement companies)
    useEffect(() => {
        if (visible) {
            loadCollectionEntities();
            if (rule) {
                form.setFieldsValue({
                    name: rule.name,
                    description: rule.description,
                    priority: rule.priority,
                    targetEntityId: rule.targetEntityId,
                    conditions: fromApiConditions(rule.conditions || []),
                    status: rule.status,
                });
            } else {
                form.resetFields();
                // Set default values
                form.setFieldsValue({
                    priority: 100,
                    conditions: [{ id: Date.now().toString(), field: 'game_type', operator: 'eq', value: '' }],
                });
            }
        }
    }, [visible, rule, form]);

    const loadCollectionEntities = async () => {
        setLoadingEntities(true);
        try {
            const response = await settlementApi.getSettlementCompanies({
                pageSize: 1000,
                status: 'active'
            });
            if (response.data?.data) {
                setCollectionEntities(response.data.data);
            }
        } catch {
            message.error('加载收款主体失败');
        } finally {
            setLoadingEntities(false);
        }
    };

    const handleSubmit = async () => {
        try {
            const values = await form.validateFields();
            setSubmitting(true);

            const payload = {
                name: values.name,
                description: values.description || '',
                priority: values.priority,
                targetEntityId: values.targetEntityId,
                conditions: toApiConditions(values.conditions),
            };

            if (rule) {
                await routingApi.updateRoutingRule(rule.id, payload);
                message.success('更新成功');
            } else {
                await routingApi.createRoutingRule(payload);
                message.success('创建成功');
            }

            onOk();
        } catch (error: unknown) {
            // Validation error or API error
            if (error && typeof error === 'object' && 'errorFields' in error) {
                // Form validation error - do nothing, Form will display errors
            } else {
                message.error(rule ? '更新失败' : '创建失败');
            }
        } finally {
            setSubmitting(false);
        }
    };

    return (
        <Modal
            title={rule ? '编辑路由规则' : '创建路由规则'}
            open={visible}
            onOk={handleSubmit}
            onCancel={onCancel}
            confirmLoading={submitting}
            width={700}
            destroyOnHidden
        >
            <Form
                form={form}
                layout="vertical"
                preserve={false}
            >
                <Form.Item
                    name="name"
                    label="规则名称"
                    rules={[
                        { required: true, message: '请输入规则名称' },
                        { max: 100, message: '名称不能超过100个字符' },
                    ]}
                >
                    <Input placeholder="例如: 高金额订单路由到公司A" />
                </Form.Item>

                <Form.Item
                    name="description"
                    label="规则描述"
                    rules={[{ max: 500, message: '描述不能超过500个字符' }]}
                >
                    <Input.TextArea
                        placeholder="描述此规则的用途和适用场景"
                        rows={2}
                        maxLength={500}
                        showCount
                    />
                </Form.Item>

                <Form.Item
                    name="priority"
                    label="优先级"
                    rules={[
                        { required: true, message: '请输入优先级' },
                        { type: 'number', min: 1, max: 1000, message: '优先级范围: 1-1000' },
                    ]}
                    tooltip="数字越小优先级越高，系统按优先级顺序匹配规则"
                >
                    <InputNumber
                        placeholder="1-1000，数字越小优先级越高"
                        min={1}
                        max={1000}
                        style={{ width: '100%' }}
                    />
                </Form.Item>

                <Form.Item
                    name="targetEntityId"
                    label="目标收款主体"
                    rules={[{ required: true, message: '请选择收款主体' }]}
                    tooltip="订单匹配此规则时，将路由到此主体收款"
                >
                    <Select
                        placeholder="请选择收款主体"
                        loading={loadingEntities}
                        showSearch
                        optionFilterProp="label"
                    >
                        {collectionEntities.map(entity => (
                            <Select.Option
                                key={entity.id}
                                value={entity.id}
                                label={entity.name}
                            >
                                <Space>
                                    <span>{entity.name}</span>
                                    <span style={{ color: '#999', fontSize: 12 }}>
                                        ({entity.type === 'company' ? '企业' : '个人'})
                                    </span>
                                </Space>
                            </Select.Option>
                        ))}
                    </Select>
                </Form.Item>

                <Divider>匹配条件</Divider>

                <Form.Item
                    name="conditions"
                    label="路由条件"
                    rules={[{ required: true, message: '请至少添加一个条件' }]}
                >
                    <ConditionBuilder />
                </Form.Item>

                {rule && (
                    <>
                        <Divider />
                        <Form.Item
                            name="status"
                            label="规则状态"
                            valuePropName="checked"
                            initialValue={rule.status === 'active'}
                        >
                            <Switch
                                checkedChildren="启用"
                                unCheckedChildren="禁用"
                            />
                        </Form.Item>
                    </>
                )}

                <div style={{
                    padding: 12,
                    background: '#f0f5ff',
                    borderRadius: 4,
                    marginTop: 8
                }}>
                    <Space orientation="vertical" size={4}>
                        <strong>规则说明:</strong>
                        <div style={{ fontSize: 12, color: '#666' }}>
                            1. 优先级数字越小，规则越优先匹配
                        </div>
                        <div style={{ fontSize: 12, color: '#666' }}>
                            2. 所有条件之间为 AND 关系，需同时满足
                        </div>
                        <div style={{ fontSize: 12, color: '#666' }}>
                            3. 若无规则匹配，将使用默认收款主体
                        </div>
                    </Space>
                </div>
            </Form>
        </Modal>
    );
};

export default RoutingForm;
