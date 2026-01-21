/**
 * 充值选项表单组件
 * 用于创建和编辑充值档位
 */
import React, { useEffect } from 'react';
import {
    Modal,
    Form,
    Input,
    InputNumber,
    Switch,
    message,
    Row,
    Col,
    Divider,
    Typography,
    Space,
} from 'antd';
import { rechargeApi, type RechargeOption, type CreateRechargeOptionDto, type UpdateRechargeOptionDto } from '@/api/recharge';

import { logger } from '@/utils/logger';
const { TextArea } = Input;
const { Text } = Typography;

interface OptionFormProps {
    visible: boolean;
    option: RechargeOption | null;
    onCancel: () => void;
    onSuccess: () => void;
}

/**
 * 充值选项表单组件
 */
const OptionForm: React.FC<OptionFormProps> = ({ visible, option, onCancel, onSuccess }) => {
    const [form] = Form.useForm();
    const [loading, setLoading] = React.useState(false);
    const isEdit = !!option;

    useEffect(() => {
        if (visible && option) {
            form.setFieldsValue({
                name: option.name,
                amountCents: option.amountCents / 100,
                bonusCents: option.bonusCents / 100,
                originalCents: option.originalCents ? option.originalCents / 100 : undefined,
                discountPercent: option.discountPercent,
                description: option.description,
                tag: option.tag,
                iconUrl: option.iconUrl,
                sortOrder: option.sortOrder,
                isActive: option.isActive,
                isRecommended: option.isRecommended,
                couponTemplateId: option.couponTemplateId,
                couponCount: option.couponCount,
                minVipLevel: option.minVipLevel,
                perUserLimit: option.perUserLimit,
                totalLimit: option.totalLimit,
            });
        } else if (visible) {
            form.resetFields();
            form.setFieldsValue({
                amountCents: undefined,
                bonusCents: 0,
                sortOrder: 0,
                isActive: true,
                isRecommended: false,
                couponCount: 0,
                perUserLimit: 0,
                totalLimit: 0,
            });
        }
    }, [visible, option, form]);

    const handleSubmit = async () => {
        try {
            const values = await form.validateFields();
            setLoading(true);

            // 将元转换为分
            const data: CreateRechargeOptionDto | UpdateRechargeOptionDto = {
                name: values.name,
                amountCents: Math.round(values.amountCents * 100),
                bonusCents: values.bonusCents ? Math.round(values.bonusCents * 100) : 0,
                originalCents: values.originalCents ? Math.round(values.originalCents * 100) : undefined,
                discountPercent: values.discountPercent,
                description: values.description,
                tag: values.tag,
                iconUrl: values.iconUrl,
                sortOrder: values.sortOrder || 0,
                isActive: values.isActive !== undefined ? values.isActive : true,
                isRecommended: values.isRecommended || false,
                couponTemplateId: values.couponTemplateId,
                couponCount: values.couponCount || 0,
                minVipLevel: values.minVipLevel,
                perUserLimit: values.perUserLimit || 0,
                totalLimit: values.totalLimit || 0,
            };

            if (isEdit) {
                await rechargeApi.updateRechargeOption(option.id, data as UpdateRechargeOptionDto);
                message.success('更新成功');
            } else {
                await rechargeApi.createRechargeOption(data as CreateRechargeOptionDto);
                message.success('创建成功');
            }

            onSuccess();
            onCancel();
            form.resetFields();
        } catch (error) {
            logger.error('Save option error:', error);
            message.error('操作失败');
        } finally {
            setLoading(false);
        }
    };

    // 计算折扣百分比
    const handleAmountChange = () => {
        const amountCents = form.getFieldValue('amountCents');
        const originalCents = form.getFieldValue('originalCents');

        if (amountCents && originalCents && originalCents > amountCents) {
            const discount = (originalCents - amountCents) / originalCents;
            form.setFieldsValue({ discountPercent: parseFloat(discount.toFixed(2)) });
        } else if (!originalCents) {
            form.setFieldsValue({ discountPercent: undefined });
        }
    };

    return (
        <Modal
            title={isEdit ? '编辑充值选项' : '新增充值选项'}
            open={visible}
            onOk={handleSubmit}
            onCancel={onCancel}
            confirmLoading={loading}
            width={700}
            destroyOnHidden
        >
            <Form form={form} layout="vertical">
                {/* 基本信息 */}
                <Divider>基本信息</Divider>

                <Row gutter={16}>
                    <Col span={12}>
                        <Form.Item
                            name="name"
                            label="充值档位名称"
                            rules={[{ required: true, message: '请输入名称' }]}
                        >
                            <Input placeholder="如：基础档、尊享档" />
                        </Form.Item>
                    </Col>
                    <Col span={12}>
                        <Form.Item
                            name="tag"
                            label="标签"
                            extra="显示在卡片上的标签，如：热卖、推荐"
                        >
                            <Input placeholder="如：热卖" />
                        </Form.Item>
                    </Col>
                </Row>

                <Form.Item
                    name="description"
                    label="描述"
                    extra="充值档位的详细描述"
                >
                    <TextArea rows={2} placeholder="请输入描述" />
                </Form.Item>

                <Form.Item
                    name="iconUrl"
                    label="图标URL"
                    extra="充值档位图标的URL地址"
                >
                    <Input placeholder="https://example.com/icon.png" />
                </Form.Item>

                {/* 金额配置 */}
                <Divider>金额配置</Divider>

                <Row gutter={16}>
                    <Col span={8}>
                        <Form.Item
                            name="originalCents"
                            label="原价（元）"
                            extra="用于显示折扣效果"
                        >
                            <InputNumber
                                min={0}
                                precision={2}
                                style={{ width: '100%' }}
                                placeholder="原价"
                                onChange={handleAmountChange}
                            />
                        </Form.Item>
                    </Col>
                    <Col span={8}>
                        <Form.Item
                            name="amountCents"
                            label="售价（元）"
                            rules={[{ required: true, message: '请输入售价' }]}
                        >
                            <InputNumber
                                min={0}
                                precision={2}
                                style={{ width: '100%' }}
                                placeholder="售价"
                                onChange={handleAmountChange}
                            />
                        </Form.Item>
                    </Col>
                    <Col span={8}>
                        <Form.Item
                            name="discountPercent"
                            label="折扣"
                        >
                            <InputNumber
                                min={0}
                                max={1}
                                step={0.01}
                                precision={2}
                                style={{ width: '100%' }}
                                placeholder="0.85表示8.5折"
                                addonAfter="OFF"
                            />
                        </Form.Item>
                    </Col>
                </Row>

                <Row gutter={16}>
                    <Col span={12}>
                        <Form.Item
                            name="bonusCents"
                            label="赠送金额（元）"
                            extra="额外赠送的金额"
                        >
                            <InputNumber
                                min={0}
                                precision={2}
                                style={{ width: '100%' }}
                                placeholder="0"
                            />
                        </Form.Item>
                    </Col>
                    <Col span={12}>
                        <Form.Item
                            name="sortOrder"
                            label="排序顺序"
                            extra="数字越小排序越靠前"
                            rules={[{ required: true, message: '请输入排序' }]}
                        >
                            <InputNumber
                                min={0}
                                style={{ width: '100%' }}
                                placeholder="0"
                            />
                        </Form.Item>
                    </Col>
                </Row>

                {/* 优惠券配置 */}
                <Divider>优惠券配置</Divider>

                <Row gutter={16}>
                    <Col span={12}>
                        <Form.Item
                            name="couponTemplateId"
                            label="优惠券模板ID"
                            extra="充值后赠送的优惠券模板"
                        >
                            <InputNumber
                                min={0}
                                style={{ width: '100%' }}
                                placeholder="优惠券模板ID"
                            />
                        </Form.Item>
                    </Col>
                    <Col span={12}>
                        <Form.Item
                            name="couponCount"
                            label="赠送数量"
                            extra="充值后赠送的优惠券数量"
                        >
                            <InputNumber
                                min={0}
                                style={{ width: '100%' }}
                                placeholder="0"
                            />
                        </Form.Item>
                    </Col>
                </Row>

                {/* 限制配置 */}
                <Divider>限制配置</Divider>

                <Row gutter={16}>
                    <Col span={8}>
                        <Form.Item
                            name="minVipLevel"
                            label="最低VIP等级"
                            extra="限制购买的最低VIP等级"
                        >
                            <InputNumber
                                min={0}
                                style={{ width: '100%' }}
                                placeholder="不限制"
                            />
                        </Form.Item>
                    </Col>
                    <Col span={8}>
                        <Form.Item
                            name="perUserLimit"
                            label="每人限购次数"
                            extra="0表示不限购"
                        >
                            <InputNumber
                                min={0}
                                style={{ width: '100%' }}
                                placeholder="0"
                            />
                        </Form.Item>
                    </Col>
                    <Col span={8}>
                        <Form.Item
                            name="totalLimit"
                            label="总限量"
                            extra="0表示不限量"
                        >
                            <InputNumber
                                min={0}
                                style={{ width: '100%' }}
                                placeholder="0"
                            />
                        </Form.Item>
                    </Col>
                </Row>

                {/* 状态配置 */}
                <Divider>状态配置</Divider>

                <Row gutter={16}>
                    <Col span={12}>
                        <Form.Item
                            name="isActive"
                            label="启用状态"
                            valuePropName="checked"
                            extra="是否在用户端显示此档位"
                        >
                            <Switch checkedChildren="启用" unCheckedChildren="禁用" />
                        </Form.Item>
                    </Col>
                    <Col span={12}>
                        <Form.Item
                            name="isRecommended"
                            label="推荐状态"
                            valuePropName="checked"
                            extra="是否标记为推荐档位"
                        >
                            <Switch checkedChildren="推荐" unCheckedChildren="普通" />
                        </Form.Item>
                    </Col>
                </Row>

                {/* 说明信息 */}
                <div style={{ marginTop: 16, padding: 12, backgroundColor: '#f0f5ff', borderRadius: 4 }}>
                    <Space orientation="vertical" size={4}>
                        <Text style={{ fontSize: 12 }}>
                            • 用户充值时，将支付售价并获得到账金额（售价+赠送金额）
                        </Text>
                        <Text style={{ fontSize: 12 }}>
                            • 设置原价后可显示折扣效果，如原价100元，售价80元，显示8折
                        </Text>
                        <Text style={{ fontSize: 12 }}>
                            • 优惠券配置可选，配置后用户充值成功将自动获得对应优惠券
                        </Text>
                    </Space>
                </div>
            </Form>
        </Modal>
    );
};

export default OptionForm;
