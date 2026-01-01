/**
 * 退款模态框组件
 * 用于处理充值记录退款
 */
import React, { useState, useEffect } from 'react';
import {
    Modal,
    Form,
    Input,
    message,
    Descriptions,
    Tag,
    Typography,
    Space,
    Alert,
} from 'antd';
import { rechargeApi, type RechargeRecord } from '@/api/recharge';
import { RECHARGE_PERMISSIONS } from '@/constants/permissions';
import { PermissionGuard } from '@/components/PermissionGuard';
import dayjs from 'dayjs';

const { TextArea } = Input;
const { Text } = Typography;

interface RefundModalProps {
    visible: boolean;
    record: RechargeRecord;
    onCancel: () => void;
    onSuccess: () => void;
}

/**
 * 退款模态框组件
 */
const RefundModal: React.FC<RefundModalProps> = ({ visible, record, onCancel, onSuccess }) => {
    const [form] = Form.useForm();
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        if (visible) {
            form.resetFields();
        }
    }, [visible, form]);

    const handleSubmit = async () => {
        try {
            const values = await form.validateFields();
            setLoading(true);

            await rechargeApi.refundRechargeRecord(record.id, { reason: values.reason });

            message.success('退款成功');
            form.resetFields();
            onSuccess();
        } catch (error) {
            console.error('Refund error:', error);
            message.error('退款失败');
        } finally {
            setLoading(false);
        }
    };

    const statusColorMap: Record<string, string> = {
        pending: 'orange',
        paid: 'success',
        failed: 'error',
        refunded: 'default',
        expired: 'default',
    };

    const statusTextMap: Record<string, string> = {
        pending: '待支付',
        paid: '已支付',
        failed: '支付失败',
        refunded: '已退款',
        expired: '已过期',
    };

    return (
        <Modal
            title="充值退款"
            open={visible}
            onOk={handleSubmit}
            onCancel={onCancel}
            confirmLoading={loading}
            okText="确认退款"
            okButtonProps={{ danger: true }}
            width={700}
            destroyOnClose
        >
            {/* 充值信息 */}
            <Descriptions column={2} bordered size="small" style={{ marginBottom: 16 }}>
                <Descriptions.Item label="订单号" span={2}>
                    <Text copyable>{record.orderNo}</Text>
                </Descriptions.Item>
                <Descriptions.Item label="用户">
                    {record.user?.name || `用户${record.userId}`}
                </Descriptions.Item>
                <Descriptions.Item label="状态">
                    <Tag color={statusColorMap[record.status]}>
                        {statusTextMap[record.status]}
                    </Tag>
                </Descriptions.Item>
                <Descriptions.Item label="充值档位">
                    {record.option?.name || '-'}
                </Descriptions.Item>
                <Descriptions.Item label="充值金额">
                    ¥{(record.amountCents / 100).toFixed(2)}
                </Descriptions.Item>
                <Descriptions.Item label="赠送金额">
                    {record.bonusCents > 0 ? `¥${(record.bonusCents / 100).toFixed(2)}` : '-'}
                </Descriptions.Item>
                <Descriptions.Item label="到账金额">
                    <Text style={{ color: '#52c41a', fontWeight: 500 }}>
                        ¥{(record.totalCents / 100).toFixed(2)}
                    </Text>
                </Descriptions.Item>
                <Descriptions.Item label="支付方式">
                    {record.paymentChannel || '-'}
                </Descriptions.Item>
                <Descriptions.Item label="支付时间" span={2}>
                    {record.paidAt ? dayjs(record.paidAt).format('YYYY-MM-DD HH:mm:ss') : '-'}
                </Descriptions.Item>
                <Descriptions.Item label="创建时间" span={2}>
                    {dayjs(record.createdAt).format('YYYY-MM-DD HH:mm:ss')}
                </Descriptions.Item>
            </Descriptions>

            {/* 警告提示 */}
            <Alert
                message="退款警告"
                description={
                    <Space direction="vertical" size={4}>
                        <Text>• 退款将从用户余额中扣除已到账的充值金额（包含赠送金额）</Text>
                        <Text>• 退款后用户获得的优惠券将被收回</Text>
                        <Text>• 此操作不可撤销，请谨慎操作</Text>
                    </Space>
                }
                type="warning"
                showIcon
                style={{ marginBottom: 16 }}
            />

            {/* 退款表单 */}
            <PermissionGuard permission={RECHARGE_PERMISSIONS.REFUND}>
                <Form form={form} layout="vertical">
                    <Form.Item
                        name="reason"
                        label="退款原因"
                        rules={[{ required: true, message: '请输入退款原因' }]}
                        extra="请详细说明退款原因，此信息将记录在退款记录中"
                    >
                        <TextArea
                            rows={4}
                            placeholder="请输入退款原因，如：用户误操作、系统错误等"
                            maxLength={500}
                            showCount
                        />
                    </Form.Item>
                </Form>
            </PermissionGuard>
        </Modal>
    );
};

export default RefundModal;
