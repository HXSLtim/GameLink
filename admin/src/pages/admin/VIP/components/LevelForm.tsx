/**
 * VIP等级表单组件
 */
import React, { useEffect } from 'react';
import {
    Modal,
    Form,
    Input,
    InputNumber,
    Switch,
    ColorPicker,
    message,
    Row,
    Col,
    Space,
    Divider,
    Typography,
} from 'antd';
import { vipApi } from '@/api/vip';
import type { VIPLevel, CreateVIPLevelDto, UpdateVIPLevelDto } from '@/api/vip';
import { GAMELINK_PRIMARY } from '@/theme';

import { logger } from '@/utils/logger';
const { Text } = Typography;

interface LevelFormProps {
    visible: boolean;
    level: VIPLevel | null;
    onCancel: () => void;
    onSuccess: () => void;
}

/**
 * 解析 benefits 字符串为字符串数组
 */
const parseBenefits = (benefitsStr: string): string[] => {
    if (!benefitsStr) return [];
    try {
        const parsed = JSON.parse(benefitsStr);
        return Array.isArray(parsed) ? parsed : [];
    } catch {
        // 如果不是 JSON 数组，按换行符分割
        return benefitsStr.split('\n').filter(s => s.trim());
    }
};

/**
 * 将字符串数组转换为 benefits 字符串
 */
const stringifyBenefits = (benefits: string[]): string => {
    return JSON.stringify(benefits.filter(s => s.trim()));
};

const LevelForm: React.FC<LevelFormProps> = ({ visible, level, onCancel, onSuccess }) => {
    const [form] = Form.useForm();
    const [loading, setLoading] = React.useState(false);
    const [benefits, setBenefits] = React.useState<string[]>([]);
    const [newBenefit, setNewBenefit] = React.useState('');
    const isEdit = !!level;

    useEffect(() => {
        if (visible && level) {
            const parsedBenefits = parseBenefits(level.benefits);
            setBenefits(parsedBenefits);
            form.setFieldsValue({
                slug: level.slug,
                title: level.title,
                expRequired: level.expRequired,
                orderDiscount: level.orderDiscount,
                monthlyCouponTemplateId: level.monthlyCouponTemplateId,
                monthlyCouponCount: level.monthlyCouponCount,
                iconUrl: level.iconUrl,
                color: level.color || GAMELINK_PRIMARY.base,
                sortOrder: level.sortOrder,
                isDefault: level.isDefault,
                isActive: level.isActive,
            });
        } else if (visible) {
            setBenefits([]);
            form.resetFields();
            form.setFieldsValue({
                expRequired: 0,
                orderDiscount: 0,
                monthlyCouponCount: 0,
                color: GAMELINK_PRIMARY.base,
                sortOrder: 0,
                isDefault: false,
                isActive: true,
            });
        }
    }, [visible, level, form]);

    const handleSubmit = async () => {
        try {
            const values = await form.validateFields();
            setLoading(true);

            const data: CreateVIPLevelDto | UpdateVIPLevelDto = {
                slug: values.slug,
                title: values.title,
                expRequired: values.expRequired,
                orderDiscount: values.orderDiscount,
                monthlyCouponTemplateId: values.monthlyCouponTemplateId,
                monthlyCouponCount: values.monthlyCouponCount,
                iconUrl: values.iconUrl,
                color: typeof values.color === 'string' ? values.color : values.color?.toHexString(),
                benefits: stringifyBenefits(benefits),
                sortOrder: values.sortOrder,
                isDefault: values.isDefault,
                isActive: values.isActive,
            };

            if (isEdit) {
                await vipApi.updateVIPLevel(level.id, data as UpdateVIPLevelDto);
                message.success('更新成功');
            } else {
                await vipApi.createVIPLevel(data as CreateVIPLevelDto);
                message.success('创建成功');
            }

            onSuccess();
            form.resetFields();
            setBenefits([]);
        } catch (error) {
            logger.error('Save VIP level error:', error);
            message.error('操作失败');
        } finally {
            setLoading(false);
        }
    };

    const handleAddBenefit = () => {
        if (!newBenefit.trim()) {
            message.warning('请输入权益内容');
            return;
        }
        if (benefits.includes(newBenefit.trim())) {
            message.warning('该权益已存在');
            return;
        }
        setBenefits([...benefits, newBenefit.trim()]);
        setNewBenefit('');
    };

    const handleRemoveBenefit = (index: number) => {
        setBenefits(benefits.filter((_, i) => i !== index));
    };

    return (
        <Modal
            title={isEdit ? '编辑VIP等级' : '新增VIP等级'}
            open={visible}
            onOk={handleSubmit}
            onCancel={onCancel}
            confirmLoading={loading}
            width={700}
            destroyOnHidden
        >
            <Form form={form} layout="vertical">
                <Row gutter={16}>
                    <Col span={12}>
                        <Form.Item
                            name="slug"
                            label="等级标识"
                            rules={[{ required: true, message: '请输入等级标识' }]}
                            extra="唯一标识，如：bronze, silver, gold"
                        >
                            <Input placeholder="等级标识" disabled={isEdit} />
                        </Form.Item>
                    </Col>
                    <Col span={12}>
                        <Form.Item
                            name="title"
                            label="等级名称"
                            rules={[{ required: true, message: '请输入等级名称' }]}
                        >
                            <Input placeholder="如：青铜会员" />
                        </Form.Item>
                    </Col>
                </Row>

                <Row gutter={16}>
                    <Col span={12}>
                        <Form.Item
                            name="expRequired"
                            label="升级所需经验"
                            rules={[{ required: true, message: '请输入所需经验' }]}
                            extra="用户达到此经验值可升级到此等级"
                        >
                            <InputNumber min={0} style={{ width: '100%' }} placeholder="0" />
                        </Form.Item>
                    </Col>
                    <Col span={12}>
                        <Form.Item
                            name="orderDiscount"
                            label="订单折扣"
                            rules={[{ required: true, message: '请输入折扣' }]}
                            extra="0-1之间，0.1表示9折"
                        >
                            <InputNumber min={0} max={1} step={0.01} style={{ width: '100%' }} placeholder="0" />
                        </Form.Item>
                    </Col>
                </Row>

                <Row gutter={16}>
                    <Col span={12}>
                        <Form.Item
                            name="sortOrder"
                            label="排序顺序"
                            rules={[{ required: true, message: '请输入排序' }]}
                            extra="数字越小排序越靠前"
                        >
                            <InputNumber min={0} style={{ width: '100%' }} placeholder="0" />
                        </Form.Item>
                    </Col>
                    <Col span={6}>
                        <Form.Item
                            name="color"
                            label="主题颜色"
                            rules={[{ required: true, message: '请选择颜色' }]}
                        >
                            <ColorPicker showText />
                        </Form.Item>
                    </Col>
                    <Col span={6}>
                        <Form.Item
                            name="isActive"
                            label="启用状态"
                            valuePropName="checked"
                        >
                            <Switch checkedChildren="启用" unCheckedChildren="禁用" />
                        </Form.Item>
                    </Col>
                </Row>

                <Form.Item
                    name="iconUrl"
                    label="图标URL"
                    extra="VIP等级图标地址"
                >
                    <Input placeholder="https://example.com/icon.png" />
                </Form.Item>

                <Divider>月度优惠券配置</Divider>

                <Row gutter={16}>
                    <Col span={12}>
                        <Form.Item
                            name="monthlyCouponTemplateId"
                            label="优惠券模板ID"
                            extra="可选，每月发放的优惠券模板"
                        >
                            <InputNumber min={0} style={{ width: '100%' }} placeholder="优惠券模板ID" />
                        </Form.Item>
                    </Col>
                    <Col span={12}>
                        <Form.Item
                            name="monthlyCouponCount"
                            label="每月发放数量"
                            extra="每月发放的优惠券数量"
                        >
                            <InputNumber min={0} style={{ width: '100%' }} placeholder="0" />
                        </Form.Item>
                    </Col>
                </Row>

                <Divider>会员权益</Divider>

                <div style={{ marginBottom: 16 }}>
                    <Space.Compact style={{ width: '100%' }}>
                        <Input
                            placeholder="输入权益内容，如：专属客服、优先派单等"
                            value={newBenefit}
                            onChange={(e) => setNewBenefit(e.target.value)}
                            onPressEnter={handleAddBenefit}
                        />
                        <button
                            type="button"
                            onClick={handleAddBenefit}
                            style={{
                                padding: '0 16px',
                                backgroundColor: GAMELINK_PRIMARY.base,
                                color: '#fff',
                                border: 'none',
                                borderRadius: '0 6px 6px 0',
                                cursor: 'pointer',
                            }}
                        >
                            添加
                        </button>
                    </Space.Compact>
                </div>

                {benefits.length > 0 && (
                    <div style={{ marginBottom: 16 }}>
                        <Text strong>已添加的权益：</Text>
                        <div style={{ marginTop: 8 }}>
                            {benefits.map((benefit, index) => (
                                <div
                                    key={index}
                                    style={{
                                        display: 'flex',
                                        justifyContent: 'space-between',
                                        alignItems: 'center',
                                        padding: '8px 12px',
                                        marginBottom: '8px',
                                        backgroundColor: '#f5f5f5',
                                        borderRadius: '4px',
                                    }}
                                >
                                    <span>{benefit}</span>
                                    <button
                                        type="button"
                                        onClick={() => handleRemoveBenefit(index)}
                                        style={{
                                            padding: '4px 12px',
                                            backgroundColor: '#ff4d4f',
                                            color: '#fff',
                                            border: 'none',
                                            borderRadius: '4px',
                                            cursor: 'pointer',
                                        }}
                                    >
                                        删除
                                    </button>
                                </div>
                            ))}
                        </div>
                    </div>
                )}

                <Form.Item
                    name="isDefault"
                    label="设为默认等级"
                    valuePropName="checked"
                    extra="新用户注册时的默认VIP等级"
                >
                    <Switch checkedChildren="是" unCheckedChildren="否" />
                </Form.Item>
            </Form>
        </Modal>
    );
};

export default LevelForm;
