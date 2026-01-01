/**
 * Issue Coupon Modal Component
 * Modal for issuing coupons to users
 */
import React, { useState, useEffect } from 'react';
import {
    Modal,
    Form,
    Select,
    message,
    Spin,
    Typography,
    Tag,
    Descriptions,
    Divider,
} from 'antd';
import {
    couponApi,
    type CouponTemplate,
    type CouponSource,
    getCouponTypeLabel,
    getCouponScopeLabel,
    getCouponSourceLabel,
    centsToYuan,
} from '@/api/coupon';
import { adminApi } from '@/api/admin';

const { Text } = Typography;

interface IssueModalProps {
    visible: boolean;
    templateId?: number;
    onCancel: () => void;
    onSuccess: () => void;
}

interface User {
    id: number;
    name: string;
    email: string;
    phone: string;
}

const IssueModal: React.FC<IssueModalProps> = ({
    visible,
    templateId,
    onCancel,
    onSuccess,
}) => {
    const [form] = Form.useForm();
    const [loading, setLoading] = useState(false);
    const [submitting, setSubmitting] = useState(false);
    const [templates, setTemplates] = useState<CouponTemplate[]>([]);
    const [users, setUsers] = useState<User[]>([]);
    const [selectedTemplate, setSelectedTemplate] = useState<CouponTemplate | null>(null);

    // Load templates and users
    useEffect(() => {
        if (visible) {
            loadData();
        }
    }, [visible]);

    useEffect(() => {
        if (templateId && templates.length > 0) {
            const template = templates.find((t) => t.id === templateId);
            if (template) {
                form.setFieldsValue({ templateId: template.id });
                setSelectedTemplate(template);
            }
        }
    }, [templateId, templates, form]);

    const loadData = async () => {
        setLoading(true);
        try {
            // Load active templates
            const templateRes = await couponApi.getTemplates({ isActive: true, page_size: 100 });
            if (templateRes.data?.success && templateRes.data?.data) {
                setTemplates(templateRes.data.data);
            }

            // Load users (for selection)
            const userRes = await adminApi.getUsers({ page_size: 100 });
            if (userRes.data?.success && userRes.data?.data) {
                setUsers(userRes.data.data);
            }
        } catch (err) {
            console.error('Failed to load data:', err);
            message.error('加载数据失败');
        } finally {
            setLoading(false);
        }
    };

    const handleTemplateChange = (templateId: number) => {
        const template = templates.find((t) => t.id === templateId);
        setSelectedTemplate(template || null);
    };

    const handleOk = async () => {
        try {
            const values = await form.validateFields();
            setSubmitting(true);

            // If templateId is provided, use it; otherwise use form value
            const finalTemplateId = templateId || values.templateId;

            if (values.targetType === 'single') {
                // Issue to single user
                const res = await couponApi.issueCoupon({
                    userId: values.userId,
                    templateId: finalTemplateId,
                    source: values.source,
                });
                if (res.data?.success) {
                    message.success('发放成功');
                    form.resetFields();
                    setSelectedTemplate(null);
                    onSuccess();
                } else {
                    message.error(res.data?.message || '发放失败');
                }
            } else if (values.targetType === 'batch') {
                // Batch issue (loop through users)
                const userIds = values.userIds || [];
                let successCount = 0;
                let failCount = 0;

                for (const userId of userIds) {
                    try {
                        const res = await couponApi.issueCoupon({
                            userId,
                            templateId: finalTemplateId,
                            source: values.source,
                        });
                        if (res.data?.success) {
                            successCount++;
                        } else {
                            failCount++;
                        }
                    } catch {
                        failCount++;
                    }
                }

                if (successCount > 0) {
                    message.success(`成功发放 ${successCount} 张优惠券${failCount > 0 ? `，失败 ${failCount} 张` : ''}`);
                    form.resetFields();
                    setSelectedTemplate(null);
                    onSuccess();
                } else {
                    message.error('发放失败');
                }
            }
        } catch (err) {
            console.error('Issue coupon failed:', err);
        } finally {
            setSubmitting(false);
        }
    };

    const renderTemplateInfo = () => {
        if (!selectedTemplate) return null;

        return (
            <div style={{ marginTop: 16 }}>
                <Divider>优惠券信息</Divider>
                <Descriptions column={2} size="small" bordered>
                    <Descriptions.Item label="名称" span={2}>
                        {selectedTemplate.name}
                    </Descriptions.Item>
                    <Descriptions.Item label="类型">
                        <Tag color="blue">{getCouponTypeLabel(selectedTemplate.type)}</Tag>
                    </Descriptions.Item>
                    <Descriptions.Item label="来源">
                        <Tag>{getCouponSourceLabel(selectedTemplate.source)}</Tag>
                    </Descriptions.Item>
                    <Descriptions.Item label="适用范围">
                        <Tag>{getCouponScopeLabel(selectedTemplate.scope)}</Tag>
                    </Descriptions.Item>
                    <Descriptions.Item label="最低消费">
                        ¥{centsToYuan(selectedTemplate.minAmountCents)}
                    </Descriptions.Item>
                    {selectedTemplate.type === 'deduct' && (
                        <>
                            <Descriptions.Item label="减免金额" span={2}>
                                <Text type="danger" strong>
                                    ¥{centsToYuan(selectedTemplate.deductAmountCents)}
                                </Text>
                            </Descriptions.Item>
                        </>
                    )}
                    {selectedTemplate.type === 'discount' && (
                        <>
                            <Descriptions.Item label="折扣率">
                                {(selectedTemplate.discountRate * 10).toFixed(1)} 折
                            </Descriptions.Item>
                            <Descriptions.Item label="最大优惠">
                                ¥{centsToYuan(selectedTemplate.maxDiscountCents)}
                            </Descriptions.Item>
                        </>
                    )}
                    {selectedTemplate.validityType === 'days' && (
                        <Descriptions.Item label="有效期" span={2}>
                            领取后 {selectedTemplate.validityDays} 天有效
                        </Descriptions.Item>
                    )}
                    {selectedTemplate.validityType === 'fixed' && selectedTemplate.fixedExpireAt && (
                        <Descriptions.Item label="固定过期时间" span={2}>
                            {selectedTemplate.fixedExpireAt}
                        </Descriptions.Item>
                    )}
                    <Descriptions.Item label="已领/总数" span={2}>
                        {selectedTemplate.claimedCount} / {selectedTemplate.totalCount}
                    </Descriptions.Item>
                </Descriptions>
            </div>
        );
    };

    return (
        <Modal
            title="发放优惠券"
            open={visible}
            onOk={handleOk}
            onCancel={onCancel}
            confirmLoading={submitting}
            width={600}
            okText="发放"
            cancelText="取消"
        >
            <Spin spinning={loading}>
                <Form
                    form={form}
                    layout="vertical"
                    initialValues={{
                        targetType: 'single',
                        source: 'manual' as CouponSource,
                    }}
                >
                    {!templateId && (
                        <Form.Item
                            name="templateId"
                            label="选择优惠券模板"
                            rules={[{ required: true, message: '请选择优惠券模板' }]}
                        >
                            <Select
                                placeholder="请选择优惠券模板"
                                showSearch
                                optionFilterProp="children"
                                onChange={handleTemplateChange}
                                options={templates.map((t) => ({
                                    label: `${t.name} (${getCouponTypeLabel(t.type)})`,
                                    value: t.id,
                                }))}
                            />
                        </Form.Item>
                    )}

                    <Form.Item
                        name="targetType"
                        label="发放方式"
                        rules={[{ required: true, message: '请选择发放方式' }]}
                    >
                        <Select
                            placeholder="请选择发放方式"
                            options={[
                                { label: '发放给单个用户', value: 'single' },
                                { label: '批量发放', value: 'batch' },
                            ]}
                        />
                    </Form.Item>

                    <Form.Item noStyle shouldUpdate={(prev, curr) => prev.targetType !== curr.targetType}>
                        {({ getFieldValue }) => {
                            const targetType = getFieldValue('targetType');
                            return targetType === 'single' ? (
                                <Form.Item
                                    name="userId"
                                    label="选择用户"
                                    rules={[{ required: true, message: '请选择用户' }]}
                                >
                                    <Select
                                        placeholder="请选择用户"
                                        showSearch
                                        optionFilterProp="children"
                                        filterOption={(input, option) =>
                                            (option?.label ?? '').toLowerCase().includes(input.toLowerCase())
                                        }
                                        options={users.map((u) => ({
                                            label: `${u.name} (${u.phone})`,
                                            value: u.id,
                                        }))}
                                    />
                                </Form.Item>
                            ) : (
                                <Form.Item
                                    name="userIds"
                                    label="选择用户（可多选）"
                                    rules={[{ required: true, message: '请选择用户' }]}
                                >
                                    <Select
                                        mode="multiple"
                                        placeholder="请选择用户"
                                        showSearch
                                        optionFilterProp="children"
                                        filterOption={(input, option) =>
                                            (option?.label ?? '').toLowerCase().includes(input.toLowerCase())
                                        }
                                        options={users.map((u) => ({
                                            label: `${u.name} (${u.phone})`,
                                            value: u.id,
                                        }))}
                                    />
                                </Form.Item>
                            );
                        }}
                    </Form.Item>

                    <Form.Item
                        name="source"
                        label="发放来源"
                        rules={[{ required: true, message: '请选择发放来源' }]}
                    >
                        <Select
                            placeholder="请选择发放来源"
                            options={[
                                { label: getCouponSourceLabel('manual'), value: 'manual' },
                                { label: getCouponSourceLabel('activity'), value: 'activity' },
                                { label: getCouponSourceLabel('vip'), value: 'vip' },
                                { label: getCouponSourceLabel('referral'), value: 'referral' },
                                { label: getCouponSourceLabel('team'), value: 'team' },
                            ]}
                        />
                    </Form.Item>
                </Form>

                {renderTemplateInfo()}
            </Spin>
        </Modal>
    );
};

export default IssueModal;
