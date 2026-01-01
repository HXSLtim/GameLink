/**
 * Coupon Template Form Component
 * Create/Edit coupon template form with validation
 */
import React, { useEffect } from 'react';
import {
    Modal,
    Form,
    Input,
    Select,
    InputNumber,
    Switch,
    DatePicker,
    Divider,
    Row,
    Col,
} from 'antd';
import dayjs from 'dayjs';
import type { Dayjs } from 'dayjs';
import {
    type CouponType,
    type CouponScope,
    type CouponSource,
    type ValidityType,
    type CreateTemplateDto,
    getCouponTypeLabel,
    getCouponScopeLabel,
    getCouponSourceLabel,
    yuanToCents,
} from '@/api/coupon';
import { MONEY, LAYOUT, TABLE, MODAL, BUSINESS } from '@/constants/common';

const { TextArea } = Input;

interface TemplateFormProps {
    visible: boolean;
    editing: boolean;
    initialValues?: CreateTemplateDto;
    loading: boolean;
    onCancel: () => void;
    onSubmit: (values: CreateTemplateDto) => Promise<void>;
}

// Form-specific type with dayjs for DatePicker
interface TemplateFormData extends Omit<CreateTemplateDto, 'fixedExpireAt'> {
    fixedExpireAt?: Dayjs;
}

const TemplateForm: React.FC<TemplateFormProps> = ({
    visible,
    editing,
    initialValues,
    loading,
    onCancel,
    onSubmit,
}) => {
    const [form] = Form.useForm<TemplateFormData>();
    const couponType = Form.useWatch('type', form);
    const couponScope = Form.useWatch('scope', form);
    const validityType = Form.useWatch('validityType', form);

    // Set initial form values
    useEffect(() => {
        if (visible && !editing) {
            form.resetFields();
            form.setFieldsValue({
                type: 'deduct' as CouponType,
                scope: 'all' as CouponScope,
                source: 'manual' as CouponSource,
                validityType: 'days' as ValidityType,
                validityDays: MONEY.DEFAULT_VALIDITY_DAYS,
                perUserLimit: MONEY.DEFAULT_PER_USER_LIMIT,
                totalCount: MONEY.DEFAULT_TOTAL_COUNT,
                minAmountCents: yuanToCents(MONEY.DEFAULT_MIN_AMOUNT),
                deductAmountCents: yuanToCents(MONEY.DEFAULT_DEDUCT_AMOUNT),
                discountRate: BUSINESS.DISCOUNT_RATE_STEP,
                maxDiscountCents: yuanToCents(MONEY.DEFAULT_MAX_DISCOUNT),
                isActive: true,
            });
        } else if (visible && editing && initialValues) {
            form.setFieldsValue({
                ...initialValues,
            });
            // Fix: Convert dayjs to string for form field values
            form.setFieldValue('fixedExpireAt', initialValues.fixedExpireAt
                ? dayjs(initialValues.fixedExpireAt)
                : undefined);
        }
    }, [visible, editing, initialValues, form]);

    const handleOk = async () => {
        try {
            const values = await form.validateFields();
            const data: CreateTemplateDto = {
                ...values,
                deductAmountCents: values.deductAmountCents
                    ? yuanToCents(values.deductAmountCents as number)
                    : 0,
                minAmountCents: values.minAmountCents
                    ? yuanToCents(values.minAmountCents as number)
                    : 0,
                maxDiscountCents: values.maxDiscountCents
                    ? yuanToCents(values.maxDiscountCents as number)
                    : 0,
                fixedExpireAt: values.fixedExpireAt
                    ? dayjs(values.fixedExpireAt as Dayjs).format('YYYY-MM-DD HH:mm:ss')
                    : undefined,
            };
            await onSubmit(data);
        } catch (error) {
            console.error('Form validation failed:', error);
        }
    };

    return (
        <Modal
            title={editing ? '编辑优惠券模板' : '新建优惠券模板'}
            open={visible}
            onOk={handleOk}
            onCancel={onCancel}
            confirmLoading={loading}
            width={MODAL.WIDTH.XLARGE}
            okText="保存"
            cancelText="取消"
        >
            <Form
                form={form}
                layout="vertical"
                autoComplete="off"
            >
                {/* Basic Information */}
                <Divider>基本信息</Divider>

                <Row gutter={LAYOUT.GUTTER}>
                    <Col span={TABLE.COL_SPAN.HALF}>
                        <Form.Item
                            name="name"
                            label="优惠券名称"
                            rules={[{ required: true, message: '请输入优惠券名称' }]}
                        >
                            <Input placeholder="如：新人满减券" maxLength={MONEY.COUPON_NAME_MAX_LENGTH} />
                        </Form.Item>
                    </Col>
                    <Col span={TABLE.COL_SPAN.HALF}>
                        <Form.Item
                            name="type"
                            label="优惠券类型"
                            rules={[{ required: true, message: '请选择优惠券类型' }]}
                        >
                            <Select
                                placeholder="请选择类型"
                                options={[
                                    { label: getCouponTypeLabel('deduct'), value: 'deduct' },
                                    { label: getCouponTypeLabel('discount'), value: 'discount' },
                                ]}
                            />
                        </Form.Item>
                    </Col>
                </Row>

                <Row gutter={LAYOUT.GUTTER}>
                    <Col span={TABLE.COL_SPAN.HALF}>
                        <Form.Item
                            name="source"
                            label="发放来源"
                            rules={[{ required: true, message: '请选择发放来源' }]}
                        >
                            <Select
                                placeholder="请选择来源"
                                options={[
                                    { label: getCouponSourceLabel('new_user'), value: 'new_user' },
                                    { label: getCouponSourceLabel('link'), value: 'link' },
                                    { label: getCouponSourceLabel('vip'), value: 'vip' },
                                    { label: getCouponSourceLabel('activity'), value: 'activity' },
                                    { label: getCouponSourceLabel('manual'), value: 'manual' },
                                    { label: getCouponSourceLabel('referral'), value: 'referral' },
                                    { label: getCouponSourceLabel('team'), value: 'team' },
                                ]}
                            />
                        </Form.Item>
                    </Col>
                    <Col span={TABLE.COL_SPAN.HALF}>
                        <Form.Item
                            name="scope"
                            label="适用范围"
                            rules={[{ required: true, message: '请选择适用范围' }]}
                        >
                            <Select
                                placeholder="请选择范围"
                                options={[
                                    { label: getCouponScopeLabel('all'), value: 'all' },
                                    { label: getCouponScopeLabel('game'), value: 'game' },
                                    { label: getCouponScopeLabel('item'), value: 'item' },
                                ]}
                            />
                        </Form.Item>
                    </Col>
                </Row>

                <Form.Item
                    name="description"
                    label="优惠券描述"
                >
                    <TextArea
                        rows={BUSINESS.FORM_ROWS.DEFAULT}
                        placeholder="请输入优惠券使用说明"
                        maxLength={MONEY.COUPON_DESC_MAX_LENGTH}
                    />
                </Form.Item>

                {/* Discount Configuration */}
                <Divider>优惠配置</Divider>

                {couponType === 'deduct' && (
                    <Row gutter={LAYOUT.GUTTER}>
                        <Col span={TABLE.COL_SPAN.HALF}>
                            <Form.Item
                                name="deductAmountCents"
                                label="减免金额（元）"
                                rules={[{ required: true, message: '请输入减免金额' }]}
                                normalize={(value) => (value !== undefined ? Number(value) : undefined)}
                            >
                                <InputNumber
                                    min={0}
                                    precision={BUSINESS.PRECISION.AMOUNT}
                                    style={{ width: '100%' }}
                                    placeholder="如：10"
                                />
                            </Form.Item>
                        </Col>
                        <Col span={TABLE.COL_SPAN.HALF}>
                            <Form.Item
                                name="minAmountCents"
                                label="最低消费（元）"
                                rules={[{ required: true, message: '请输入最低消费金额' }]}
                                normalize={(value) => (value !== undefined ? Number(value) : undefined)}
                            >
                                <InputNumber
                                    min={0}
                                    precision={BUSINESS.PRECISION.AMOUNT}
                                    style={{ width: '100%' }}
                                    placeholder="如：100"
                                />
                            </Form.Item>
                        </Col>
                    </Row>
                )}

                {couponType === 'discount' && (
                    <Row gutter={LAYOUT.GUTTER}>
                        <Col span={TABLE.COL_SPAN.HALF}>
                            <Form.Item
                                name="discountRate"
                                label="折扣率（0-1）"
                                rules={[{ required: true, message: '请输入折扣率' }]}
                                normalize={(value) => (value !== undefined ? Number(value) : undefined)}
                                extra="如：0.9 表示 9 折"
                            >
                                <InputNumber
                                    min={MONEY.MIN_DISCOUNT_RATE}
                                    max={MONEY.MAX_DISCOUNT_RATE}
                                    step={MONEY.DISCOUNT_RATE_STEP}
                                    precision={BUSINESS.PRECISION.RATE}
                                    style={{ width: '100%' }}
                                    placeholder="如：0.9"
                                />
                            </Form.Item>
                        </Col>
                        <Col span={TABLE.COL_SPAN.HALF}>
                            <Form.Item
                                name="maxDiscountCents"
                                label="最大优惠金额（元）"
                                normalize={(value) => (value !== undefined ? Number(value) : undefined)}
                            >
                                <InputNumber
                                    min={0}
                                    precision={BUSINESS.PRECISION.AMOUNT}
                                    style={{ width: '100%' }}
                                    placeholder="如：50（0表示不限制）"
                                />
                            </Form.Item>
                        </Col>
                        <Col span={TABLE.COL_SPAN.HALF}>
                            <Form.Item
                                name="minAmountCents"
                                label="最低消费（元）"
                                rules={[{ required: true, message: '请输入最低消费金额' }]}
                                normalize={(value) => (value !== undefined ? Number(value) : undefined)}
                            >
                                <InputNumber
                                    min={0}
                                    precision={BUSINESS.PRECISION.AMOUNT}
                                    style={{ width: '100%' }}
                                    placeholder="如：100"
                                />
                            </Form.Item>
                        </Col>
                    </Row>
                )}

                {couponScope !== 'all' && (
                    <Form.Item
                        name="gameIds"
                        label={couponScope === 'game' ? '适用游戏ID（JSON数组）' : '适用服务ID（JSON数组）'}
                        rules={[{ required: true, message: `请输入${couponScope === 'game' ? '游戏' : '服务'}ID` }]}
                        extra='格式：[1, 2, 3]'
                    >
                        <Input placeholder='[1, 2, 3]' />
                    </Form.Item>
                )}

                {/* Validity Configuration */}
                <Divider>有效期配置</Divider>

                <Form.Item
                    name="validityType"
                    label="有效期类型"
                    rules={[{ required: true, message: '请选择有效期类型' }]}
                >
                    <Select
                        placeholder="请选择"
                        options={[
                            { label: '相对天数（领取后N天有效）', value: 'days' },
                            { label: '固定时间（指定过期日期）', value: 'fixed' },
                        ]}
                    />
                </Form.Item>

                {validityType === 'days' && (
                    <Form.Item
                        name="validityDays"
                        label="有效天数"
                        rules={[{ required: true, message: '请输入有效天数' }]}
                    >
                        <InputNumber
                            min={1}
                            max={MONEY.MAX_VALIDITY_DAYS}
                            style={{ width: '100%' }}
                            placeholder="如：30"
                            addonAfter="天"
                        />
                    </Form.Item>
                )}

                {validityType === 'fixed' && (
                    <Form.Item
                        name="fixedExpireAt"
                        label="固定过期时间"
                        rules={[{ required: true, message: '请选择过期时间' }]}
                    >
                        <DatePicker
                            showTime
                            style={{ width: '100%' }}
                            placeholder="请选择过期时间"
                            disabledDate={(current) => current && current < dayjs().startOf('day')}
                        />
                    </Form.Item>
                )}

                {/* Claim Configuration */}
                <Divider>发放配置</Divider>

                <Row gutter={LAYOUT.GUTTER}>
                    <Col span={TABLE.COL_SPAN.HALF}>
                        <Form.Item
                            name="totalCount"
                            label="发放总数"
                            rules={[{ required: true, message: '请输入发放总数' }]}
                        >
                            <InputNumber
                                min={1}
                                max={MONEY.MAX_TOTAL_COUNT}
                                style={{ width: '100%' }}
                                placeholder="如：100"
                            />
                        </Form.Item>
                    </Col>
                    <Col span={TABLE.COL_SPAN.HALF}>
                        <Form.Item
                            name="perUserLimit"
                            label="每人限领"
                            rules={[{ required: true, message: '请输入每人限领数量' }]}
                            extra="0表示不限制"
                        >
                            <InputNumber
                                min={0}
                                max={MONEY.MAX_PER_USER_LIMIT}
                                style={{ width: '100%' }}
                                placeholder="如：1"
                            />
                        </Form.Item>
                    </Col>
                </Row>

                <Form.Item
                    name="claimLink"
                    label="领取链接（可选）"
                    extra="用户通过此链接领取优惠券"
                >
                    <Input placeholder="如：https://gamelink.com/coupon/newuser" />
                </Form.Item>

                {/* Status */}
                <Divider>状态</Divider>

                <Form.Item
                    name="isActive"
                    label="启用状态"
                    valuePropName="checked"
                >
                    <Switch
                        checkedChildren="启用"
                        unCheckedChildren="禁用"
                    />
                </Form.Item>
            </Form>
        </Modal>
    );
};

export default TemplateForm;
