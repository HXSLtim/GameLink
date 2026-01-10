/**
 * Reward Form Component
 * Create/Edit activity reward form
 */
import React, { useEffect } from 'react';
import {
    Modal,
    Form,
    Select,
    InputNumber,
    Row,
    Col,
} from 'antd';
import {
    type CreateRewardDto,
    type ActivityReward,
} from '@/api/activity';
import { couponApi } from '@/api/coupon';

import { logger } from '@/utils/logger';
interface RewardFormProps {
    visible: boolean;
    editing: boolean;
    activityId: number;
    initialValues?: ActivityReward | null;
    loading: boolean;
    onCancel: () => void;
    onSubmit: (values: CreateRewardDto) => Promise<void>;
}

const RewardForm: React.FC<RewardFormProps> = ({
    visible,
    editing,
    activityId,
    initialValues,
    loading,
    onCancel,
    onSubmit,
}) => {
    const [form] = Form.useForm<CreateRewardDto>();
    const [couponTemplates, setCouponTemplates] = React.useState<Array<{id: number; name: string}>>([]);
    const [loadingTemplates, setLoadingTemplates] = React.useState(false);

    // Load coupon templates for selection
    useEffect(() => {
        const loadCouponTemplates = async () => {
            setLoadingTemplates(true);
            try {
                const res = await couponApi.getTemplates({ isActive: true, page_size: 1000 });
                if (res.data?.success && res.data?.data) {
                    setCouponTemplates(res.data.data.map(t => ({ id: t.id, name: t.name })));
                }
            } catch (err) {
                logger.error('Failed to load coupon templates:', err);
            } finally {
                setLoadingTemplates(false);
            }
        };

        if (visible) {
            loadCouponTemplates();
        }
    }, [visible]);

    // Set initial form values
    useEffect(() => {
        if (visible && !editing) {
            form.resetFields();
            form.setFieldsValue({
                activityId,
                couponCount: 1,
                probability: 100,
                totalStock: 1000,
                sortOrder: 0,
            });
        } else if (visible && editing && initialValues) {
            form.setFieldsValue({
                activityId: initialValues.activityId,
                couponTemplateId: initialValues.couponTemplateId,
                couponCount: initialValues.couponCount,
                probability: initialValues.probability,
                totalStock: initialValues.totalStock,
                sortOrder: initialValues.sortOrder,
            });
        }
    }, [visible, editing, initialValues, activityId, form]);

    const handleOk = async () => {
        try {
            const values = await form.validateFields();
            const data: CreateRewardDto = {
                ...values,
                activityId,
            };
            await onSubmit(data);
        } catch (error) {
            logger.error('Form validation failed:', error);
        }
    };

    return (
        <Modal
            title={editing ? '编辑奖励' : '添加奖励'}
            open={visible}
            onOk={handleOk}
            onCancel={onCancel}
            confirmLoading={loading || loadingTemplates}
            width={600}
            okText="保存"
            cancelText="取消"
        >
            <Form
                form={form}
                layout="vertical"
                autoComplete="off"
            >
                <Form.Item
                    name="couponTemplateId"
                    label="优惠券模板"
                    rules={[{ required: true, message: '请选择优惠券模板' }]}
                >
                    <Select
                        placeholder="请选择优惠券模板"
                        loading={loadingTemplates}
                        showSearch
                        optionFilterProp="children"
                        options={couponTemplates.map(t => ({
                            label: t.name,
                            value: t.id,
                        }))}
                    />
                </Form.Item>

                <Row gutter={16}>
                    <Col span={12}>
                        <Form.Item
                            name="couponCount"
                            label="发放数量"
                            rules={[{ required: true, message: '请输入发放数量' }]}
                            extra="每个中奖用户获得的优惠券数量"
                        >
                            <InputNumber
                                min={1}
                                max={100}
                                style={{ width: '100%' }}
                                placeholder="如：1"
                            />
                        </Form.Item>
                    </Col>
                    <Col span={12}>
                        <Form.Item
                            name="probability"
                            label="中奖概率（%）"
                            rules={[
                                { required: true, message: '请输入中奖概率' },
                                { type: 'number', min: 0, max: 100, message: '概率范围为0-100' },
                            ]}
                            extra="所有奖励概率之和应为100%"
                        >
                            <InputNumber
                                min={0}
                                max={100}
                                precision={2}
                                style={{ width: '100%' }}
                                placeholder="如：10"
                                addonAfter="%"
                            />
                        </Form.Item>
                    </Col>
                </Row>

                <Row gutter={16}>
                    <Col span={12}>
                        <Form.Item
                            name="totalStock"
                            label="库存数量"
                            rules={[{ required: true, message: '请输入库存数量' }]}
                            extra="0表示不限制"
                        >
                            <InputNumber
                                min={0}
                                max={1000000}
                                style={{ width: '100%' }}
                                placeholder="如：1000"
                            />
                        </Form.Item>
                    </Col>
                    <Col span={12}>
                        <Form.Item
                            name="sortOrder"
                            label="排序顺序"
                            extra="数值越小越靠前"
                        >
                            <InputNumber
                                min={0}
                                max={10000}
                                style={{ width: '100%' }}
                                placeholder="如：0"
                            />
                        </Form.Item>
                    </Col>
                </Row>
            </Form>
        </Modal>
    );
};

export default RewardForm;
