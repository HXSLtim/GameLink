/**
 * ResolveModal Component
 * Modal for resolving disputes with resolution type and remarks
 */
import React, { useEffect } from 'react';
import { Modal, Form, Input, Radio, Space } from 'antd';
import { DISPUTE_RESOLUTION_LABELS } from '@/types/dispute';
import type { Dispute, DisputeResolution } from '@/types/dispute';

export interface ResolveModalProps {
    /** Modal visibility */
    open: boolean;
    /** Dispute data */
    dispute: Dispute | null;
    /** Confirm callback */
    onConfirm: (resolution: string, resolveRemark: string) => void;
    /** Cancel callback */
    onCancel: () => void;
    /** Confirm loading state */
    confirmLoading?: boolean;
}

const resolutionOptions: { value: DisputeResolution; label: string }[] = [
    { value: 'refund', label: DISPUTE_RESOLUTION_LABELS.refund },
    { value: 'partial', label: DISPUTE_RESOLUTION_LABELS.partial },
    { value: 'reassign', label: DISPUTE_RESOLUTION_LABELS.reassign },
    { value: 'reject', label: DISPUTE_RESOLUTION_LABELS.reject },
];

/**
 * ResolveModal Component
 */
export const ResolveModal: React.FC<ResolveModalProps> = ({
    open,
    dispute,
    onConfirm,
    onCancel,
    confirmLoading = false,
}) => {
    const [form] = Form.useForm();

    useEffect(() => {
        if (open) {
            form.resetFields();
        }
    }, [open, form]);

    const handleOk = async () => {
        try {
            const values = await form.validateFields();
            onConfirm(values.resolution, values.resolveRemark);
        } catch {
            // Validation failed
        }
    };

    return (
        <Modal
            title="解决纠纷"
            open={open}
            onOk={handleOk}
            onCancel={onCancel}
            confirmLoading={confirmLoading}
            width={550}
            okText="确认解决"
            cancelText="取消"
        >
            {dispute && (
                <>
                    {/* Dispute Info Summary */}
                    <div
                        style={{
                            marginBottom: 16,
                            padding: 12,
                            backgroundColor: '#fafafa',
                            borderRadius: 4,
                        }}
                    >
                        <Space direction="vertical" size={4}>
                            <div>
                                <strong>纠纷ID:</strong> {dispute.id}
                            </div>
                            <div>
                                <strong>订单号:</strong> {dispute.orderNo || `ID: ${dispute.orderId}`}
                            </div>
                            <div>
                                <strong>发起人:</strong> {dispute.initiatorName || `ID: ${dispute.initiatorId}`}
                            </div>
                            <div>
                                <strong>纠纷原因:</strong> {dispute.reason}
                            </div>
                        </Space>
                    </div>

                    <Form form={form} layout="vertical">
                        <Form.Item
                            name="resolution"
                            label="解决方案"
                            rules={[{ required: true, message: '请选择解决方案' }]}
                        >
                            <Radio.Group>
                                <Space direction="vertical">
                                    {resolutionOptions.map((option) => (
                                        <Radio key={option.value} value={option.value}>
                                            {option.label}
                                        </Radio>
                                    ))}
                                </Space>
                            </Radio.Group>
                        </Form.Item>

                        <Form.Item
                            name="resolveRemark"
                            label="处理备注"
                            rules={[{ required: true, message: '请输入处理备注' }]}
                            extra="请详细说明处理依据和结果，用于后续审计"
                        >
                            <Input.TextArea
                                rows={4}
                                placeholder="请输入处理备注，例如：经核实，陪玩师确实未完成约定服务时长，同意全额退款。"
                                maxLength={500}
                                showCount
                            />
                        </Form.Item>
                    </Form>

                    {/* Resolution Info */}
                    <div
                        style={{
                            padding: 12,
                            backgroundColor: '#fffbe6',
                            borderRadius: 4,
                            fontSize: 12,
                            color: '#d48806',
                        }}
                    >
                        <div>
                            <strong>解决方案说明:</strong>
                        </div>
                        <ul style={{ margin: '8px 0', paddingLeft: 20 }}>
                            <li>
                                <strong>全额退款:</strong> 订单金额全额退款给用户
                            </li>
                            <li>
                                <strong>部分退款:</strong> 按一定比例退款（需在备注中说明退款金额）
                            </li>
                            <li>
                                <strong>重新指派:</strong> 重新为用户分配陪玩师
                            </li>
                            <li>
                                <strong>驳回:</strong> 驳回纠纷申请，维持原订单状态
                            </li>
                        </ul>
                    </div>
                </>
            )}
        </Modal>
    );
};

export default ResolveModal;
