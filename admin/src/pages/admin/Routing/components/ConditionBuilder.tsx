/**
 * Condition Builder Component
 * Visual condition builder for routing rules
 */
import React, { useState, useCallback } from 'react';
import {
    Card,
    Form,
    Select,
    Input,
    InputNumber,
    Button,
    Space,
    Tag,
    Row,
    Col,
    Divider,
    Typography,
} from 'antd';
import {
    PlusOutlined,
    DeleteOutlined,
    InfoCircleOutlined,
} from '@ant-design/icons';
import type { ConditionField, ConditionOperator } from '@/api/routing';
import type {
    ConditionFormItem,
    FieldOption,
    OperatorOption,
} from '../types';

const { Text } = Typography;

/**
 * Field configurations for condition builder
 */
const FIELD_OPTIONS: FieldOption[] = [
    {
        value: 'game_type',
        label: '游戏类型',
        type: 'select',
        operators: ['eq', 'neq', 'in', 'not_in'],
        options: [
            { label: '王者荣耀', value: 'honor_of_kings' },
            { label: '和平精英', value: 'pubg_mobile' },
            { label: '英雄联盟', value: 'league_of_legends' },
            { label: '绝地求生', value: 'pubg' },
            { label: '原神', value: 'genshin_impact' },
            { label: '永劫无间', value: 'naraka' },
            { label: '其他', value: 'other' },
        ],
    },
    {
        value: 'service_type',
        label: '服务类型',
        type: 'select',
        operators: ['eq', 'neq', 'in', 'not_in'],
        options: [
            { label: '陪练', value: 'companion' },
            { label: '代练', value: 'boosting' },
            { label: '陪玩', value: 'play_with' },
            { label: '教学', value: 'coaching' },
        ],
    },
    {
        value: 'order_amount',
        label: '订单金额',
        type: 'number',
        operators: ['eq', 'neq', 'gt', 'lt', 'between'],
    },
    {
        value: 'region',
        label: '地区',
        type: 'select',
        operators: ['eq', 'neq', 'in', 'not_in'],
        options: [
            { label: '华东', value: 'east_china' },
            { label: '华南', value: 'south_china' },
            { label: '华北', value: 'north_china' },
            { label: '西南', value: 'southwest_china' },
            { label: '西北', value: 'northwest_china' },
            { label: '东北', value: 'northeast_china' },
        ],
    },
];

/**
 * Operator configurations
 */
const OPERATOR_OPTIONS: OperatorOption[] = [
    { value: 'eq', label: '等于', requiresMultiple: false },
    { value: 'neq', label: '不等于', requiresMultiple: false },
    { value: 'in', label: '包含于', requiresMultiple: true },
    { value: 'not_in', label: '不包含于', requiresMultiple: true },
    { value: 'gt', label: '大于', requiresMultiple: false },
    { value: 'lt', label: '小于', requiresMultiple: false },
    { value: 'between', label: '介于', requiresMultiple: true },
];

interface ConditionBuilderProps {
    value?: ConditionFormItem[];
    onChange?: (conditions: ConditionFormItem[]) => void;
    disabled?: boolean;
}

/**
 * Condition Builder Component
 */
const ConditionBuilder: React.FC<ConditionBuilderProps> = ({
    value = [],
    onChange,
    disabled = false,
}) => {
    const [conditions, setConditions] = useState<ConditionFormItem[]>(
        value.length > 0 ? value : [{ id: Date.now().toString(), field: 'game_type', operator: 'eq', value: '' }]
    );

    const updateConditions = useCallback((newConditions: ConditionFormItem[]) => {
        setConditions(newConditions);
        onChange?.(newConditions);
    }, [onChange]);

    const addCondition = useCallback(() => {
        const newCondition: ConditionFormItem = {
            id: Date.now().toString(),
            field: 'game_type',
            operator: 'eq',
            value: '',
        };
        updateConditions([...conditions, newCondition]);
    }, [conditions, updateConditions]);

    const removeCondition = useCallback((id: string) => {
        if (conditions.length > 1) {
            updateConditions(conditions.filter(c => c.id !== id));
        }
    }, [conditions, updateConditions]);

    const updateCondition = useCallback((id: string, updates: Partial<ConditionFormItem>) => {
        const newConditions = conditions.map(c =>
            c.id === id ? { ...c, ...updates } : c
        );
        updateConditions(newConditions);
    }, [conditions, updateConditions]);

    const getFieldConfig = (field: ConditionField): FieldOption | undefined => {
        return FIELD_OPTIONS.find(f => f.value === field);
    };

    const getOperatorConfig = (operator: ConditionOperator): OperatorOption | undefined => {
        return OPERATOR_OPTIONS.find(o => o.value === operator);
    };

    const renderValueInput = (condition: ConditionFormItem) => {
        const fieldConfig = getFieldConfig(condition.field);
        const operatorConfig = getOperatorConfig(condition.operator);

        if (!fieldConfig) return null;

        const isMultiple = operatorConfig?.requiresMultiple || false;
        const isNumber = fieldConfig.type === 'number';

        // For 'between' operator, render two inputs
        if (condition.operator === 'between') {
            return (
                <Input.Group compact>
                    <InputNumber
                        placeholder="最小值"
                        value={Array.isArray(condition.value) ? condition.value[0] : undefined}
                        onChange={(val) => updateCondition(condition.id, {
                            value: [val || 0, Array.isArray(condition.value) ? condition.value[1] : 0]
                        })}
                        disabled={disabled}
                        style={{ width: '50%' }}
                        min={0}
                    />
                    <InputNumber
                        placeholder="最大值"
                        value={Array.isArray(condition.value) ? condition.value[1] : undefined}
                        onChange={(val) => updateCondition(condition.id, {
                            value: [Array.isArray(condition.value) ? condition.value[0] : 0, (val || 0) as number]
                        })}
                        disabled={disabled}
                        style={{ width: '50%' }}
                        min={0}
                    />
                </Input.Group>
            );
        }

        // For select with multiple values
        if (fieldConfig.type === 'select' && isMultiple) {
            return (
                <Select
                    mode="multiple"
                    placeholder="选择多个值"
                    value={Array.isArray(condition.value) ? condition.value : []}
                    onChange={(val) => updateCondition(condition.id, { value: val })}
                    disabled={disabled}
                    options={fieldConfig.options}
                    style={{ width: '100%' }}
                />
            );
        }

        // For select with single value
        if (fieldConfig.type === 'select') {
            return (
                <Select
                    placeholder="选择值"
                    value={condition.value as string}
                    onChange={(val) => updateCondition(condition.id, { value: val })}
                    disabled={disabled}
                    options={fieldConfig.options}
                    style={{ width: '100%' }}
                />
            );
        }

        // For number input
        if (isNumber) {
            return (
                <InputNumber
                    placeholder="输入数值"
                    value={typeof condition.value === 'number' ? condition.value : undefined}
                    onChange={(val) => updateCondition(condition.id, { value: val || 0 })}
                    disabled={disabled}
                    style={{ width: '100%' }}
                    min={0}
                />
            );
        }

        // Default text input
        return (
            <Input
                placeholder="输入值"
                value={condition.value as string}
                onChange={(e) => updateCondition(condition.id, { value: e.target.value })}
                disabled={disabled}
            />
        );
    };

    const renderConditionPreview = (condition: ConditionFormItem) => {
        const fieldConfig = getFieldConfig(condition.field);
        const operatorConfig = getOperatorConfig(condition.operator);

        const fieldLabel = fieldConfig?.label || condition.field;
        const operatorLabel = operatorConfig?.label || condition.operator;

        let displayValue: string;
        if (condition.operator === 'between' && Array.isArray(condition.value)) {
            displayValue = `${condition.value[0]} ~ ${condition.value[1]}`;
        } else if (Array.isArray(condition.value)) {
            displayValue = condition.value.join(', ');
        } else {
            displayValue = String(condition.value || '-');
        }

        return (
            <Tag color="blue" style={{ margin: 2 }}>
                {fieldLabel} {operatorLabel} {displayValue}
            </Tag>
        );
    };

    return (
        <div>
            <div style={{ marginBottom: 12 }}>
                <Space>
                    <Text strong>条件列表</Text>
                    <Text type="secondary">（满足所有条件时匹配）</Text>
                </Space>
            </div>

            {conditions.map((condition, index) => {
                const fieldConfig = getFieldConfig(condition.field);
                const availableOperators = fieldConfig?.operators || ['eq'];
                const operatorConfig = getOperatorConfig(condition.operator);

                return (
                    <Card
                        key={condition.id}
                        size="small"
                        style={{ marginBottom: 8 }}
                        styles={{
                            body: { padding: 12 }
                        }}
                    >
                        <Row gutter={16} align="middle">
                            <Col span={1}>
                                <Text type="secondary">{index + 1}</Text>
                            </Col>
                            <Col span={5}>
                                <Select
                                    placeholder="选择字段"
                                    value={condition.field}
                                    onChange={(val) => updateCondition(condition.id, {
                                        field: val,
                                        operator: 'eq',
                                        value: ''
                                    })}
                                    disabled={disabled}
                                    options={FIELD_OPTIONS.map(f => ({
                                        value: f.value,
                                        label: f.label
                                    }))}
                                    style={{ width: '100%' }}
                                />
                            </Col>
                            <Col span={5}>
                                <Select
                                    placeholder="操作符"
                                    value={condition.operator}
                                    onChange={(val) => updateCondition(condition.id, {
                                        operator: val,
                                        value: operatorConfig?.requiresMultiple ? [] : ''
                                    })}
                                    disabled={disabled}
                                    options={availableOperators.map(op => {
                                        const config = getOperatorConfig(op);
                                        return { value: op, label: config?.label || op };
                                    })}
                                    style={{ width: '100%' }}
                                />
                            </Col>
                            <Col span={10}>
                                {renderValueInput(condition)}
                            </Col>
                            <Col span={2}>
                                <Button
                                    type="text"
                                    danger
                                    icon={<DeleteOutlined />}
                                    onClick={() => removeCondition(condition.id)}
                                    disabled={disabled || conditions.length === 1}
                                />
                            </Col>
                            <Col span={1}>
                                <InfoCircleOutlined style={{ color: '#999' }} />
                            </Col>
                        </Row>
                        <Divider style={{ margin: '8px 0' }} />
                        <Space size={4} wrap>
                            <Text type="secondary" style={{ fontSize: 12 }}>
                                条件预览:
                            </Text>
                            {renderConditionPreview(condition)}
                        </Space>
                    </Card>
                );
            })}

            <Button
                type="dashed"
                onClick={addCondition}
                disabled={disabled}
                block
                icon={<PlusOutlined />}
                style={{ marginTop: 8 }}
            >
                添加条件
            </Button>

            <div style={{ marginTop: 12, padding: 8, background: '#f5f5f5', borderRadius: 4 }}>
                <Text type="secondary" style={{ fontSize: 12 }}>
                    说明: 所有条件之间为 AND 关系，即订单需同时满足所有条件才会匹配此规则
                </Text>
            </div>
        </div>
    );
};

export default ConditionBuilder;
