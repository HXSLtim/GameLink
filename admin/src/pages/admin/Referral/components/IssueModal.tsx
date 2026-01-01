/**
 * Issue Modal Component
 * 奖励发放/失败模态框组件
 *
 * Modal for issuing or marking referral rewards as failed.
 */
import React, { useState, useEffect } from 'react';
import {
    Modal,
    Form,
    Input,
    message,
    Descriptions,
    Avatar,
    Tag,
    Space,
} from 'antd';
import {
    DollarOutlined,
    CheckOutlined,
    CloseOutlined,
    GiftOutlined,
    UserOutlined,
} from '@ant-design/icons';
import dayjs from 'dayjs';
import type { ReferralReward } from '@/api/referral';
import {
    getRewardTypeLabel,
    getRewardStatusLabel,
    getRewardStatusColor,
    centsToYuan,
} from '@/api/referral';

interface IssueModalProps {
    visible: boolean;
    reward: ReferralReward | null;
    action: 'issue' | 'fail';
    onIssue: (reward: ReferralReward) => void;
    onFail: (reward: ReferralReward, reason: string) => void;
    onCancel: () => void;
}

const IssueModal: React.FC<IssueModalProps> = ({
    visible,
    reward,
    action,
    onIssue,
    onFail,
    onCancel,
}) => {
    const [form] = Form.useForm();
    const [loading, setLoading] = useState(false);
    const isIssue = action === 'issue';

    /**
     * Reset form when modal opens/closes
     */
    useEffect(() => {
        if (visible) {
            form.resetFields();
        }
    }, [visible, form]);

    /**
     * Handle confirm
     */
    const handleOk = async () => {
        if (!reward) return;

        if (isIssue) {
            // Issue reward - no form validation needed
            setLoading(true);
            try {
                await onIssue(reward);
            } finally {
                setLoading(false);
            }
        } else {
            // Fail reward - need reason
            try {
                const values = await form.validateFields();
                setLoading(true);
                await onFail(reward, values.reason);
            } catch {
                // Form validation failed
            } finally {
                setLoading(false);
            }
        }
    };

    /**
     * Modal title
     */
    const modalTitle = isIssue ? '发放奖励' : '标记奖励失败';

    return (
        <Modal
            title={
                <Space>
                    {isIssue ? <CheckOutlined /> : <CloseOutlined />}
                    {modalTitle}
                </Space>
            }
            open={visible}
            onOk={handleOk}
            onCancel={onCancel}
            confirmLoading={loading}
            okText={isIssue ? '确认发放' : '确认'}
            cancelText="取消"
            okButtonProps={{
                style: isIssue ? { backgroundColor: '#52c41a', borderColor: '#52c41a' } : { danger: true },
            }}
            width={500}
        >
            {reward && (
                <>
                    {/* Reward Summary */}
                    <div
                        style={{
                            textAlign: 'center',
                            padding: '16px 0',
                            marginBottom: 16,
                            background: '#fafafa',
                            borderRadius: 8,
                        }}
                    >
                        <div style={{ fontSize: 28, fontWeight: 'bold', color: '#faad14' }}>
                            <DollarOutlined style={{ marginRight: 8 }} />
                            ¥{centsToYuan(reward.amountCents)}
                        </div>
                        <div style={{ marginTop: 8 }}>
                            <Tag color={getRewardStatusColor(reward.status)}>
                                {getRewardStatusLabel(reward.status)}
                            </Tag>
                            <Tag color={reward.type === 'referrer' ? 'blue' : 'green'}>
                                {getRewardTypeLabel(reward.type)}
                            </Tag>
                        </div>
                    </div>

                    <Descriptions column={1} size="small" bordered>
                        <Descriptions.Item label="奖励ID">{reward.id}</Descriptions.Item>
                        <Descriptions.Item label="推荐关系ID">{reward.referralId}</Descriptions.Item>
                        <Descriptions.Item label="接收用户">
                            <Space>
                                <Avatar
                                    size="small"
                                    src={reward.user?.avatarUrl}
                                    icon={<UserOutlined />}
                                />
                                {reward.user?.name || `用户${reward.userId}`}
                            </Space>
                        </Descriptions.Item>
                        <Descriptions.Item label="创建时间">
                            {dayjs(reward.createdAt).format('YYYY-MM-DD HH:mm:ss')}
                        </Descriptions.Item>
                    </Descriptions>

                    {/* Issue confirmation message */}
                    {isIssue && (
                        <div
                            style={{
                                marginTop: 16,
                                padding: 12,
                                background: '#f6ffed',
                                border: '1px solid #b7eb8f',
                                borderRadius: 6,
                            }}
                        >
                            <Space>
                                <GiftOutlined style={{ color: '#52c41a' }} />
                                <span>
                                    确定向用户发放 <strong>¥{centsToYuan(reward.amountCents)}</strong> 奖励？
                                </span>
                            </Space>
                        </div>
                    )}

                    {/* Fail reason form */}
                    {!isIssue && (
                        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
                            <Form.Item
                                name="reason"
                                label="失败原因"
                                rules={[
                                    { required: true, message: '请输入失败原因' },
                                    { min: 5, message: '原因至少5个字符' },
                                    { max: 200, message: '原因不能超过200个字符' },
                                ]}
                            >
                                <Input.TextArea
                                    rows={3}
                                    placeholder="请输入奖励发放失败的原因（如：用户不存在、账户异常等）"
                                    maxLength={200}
                                    showCount
                                />
                            </Form.Item>
                        </Form>
                    )}
                </>
            )}
        </Modal>
    );
};

export default IssueModal;
