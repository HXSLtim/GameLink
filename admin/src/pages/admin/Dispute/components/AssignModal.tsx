/**
 * AssignModal Component
 * Modal for assigning disputes to customer service representatives
 */
import React, { useState, useEffect } from 'react';
import { Modal, Form, Select, Alert, Space, Divider } from 'antd';
import type { Dispute } from '@/types/dispute';
import { adminApi } from '@/api/admin';

import { logger } from '@/utils/logger';
export interface AssignModalProps {
    /** Modal visibility */
    open: boolean;
    /** Dispute data */
    dispute: Dispute | null;
    /** Confirm callback with assignedServiceId and optional originalServiceId */
    onConfirm: (assignedServiceId: number, originalServiceId?: number) => void;
    /** Cancel callback */
    onCancel: () => void;
    /** Confirm loading state */
    confirmLoading?: boolean;
}

interface ServiceStaff {
    id: number;
    name: string;
    email?: string;
}

/**
 * AssignModal Component
 */
export const AssignModal: React.FC<AssignModalProps> = ({
    open,
    dispute,
    onConfirm,
    onCancel,
    confirmLoading = false,
}) => {
    const [form] = Form.useForm();
    const [loading, setLoading] = useState(false);
    const [staffList, setStaffList] = useState<ServiceStaff[]>([]);

    // Load service staff list when modal opens
    useEffect(() => {
        if (open) {
            form.resetFields();
            loadStaffList();
        }
    }, [open, form]);

    const loadStaffList = async () => {
        setLoading(true);
        try {
            // Fetch admin users (customer service staff)
            const response = await adminApi.getUsers({ role: ['admin'], page_size: 100 });
            if (response.data.success && response.data.data) {
                const staff = response.data.data.map((user: { id: number; name: string; email?: string }) => ({
                    id: user.id,
                    name: user.name,
                    email: user.email,
                }));
                setStaffList(staff);
            }
        } catch (error) {
            logger.error('Failed to load staff list:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleOk = async () => {
        try {
            const values = await form.validateFields();
            onConfirm(values.assignedServiceId as number, values.originalServiceId as number | undefined);
        } catch {
            // Validation failed
        }
    };

    return (
        <Modal
            title="分配纠纷"
            open={open}
            onOk={handleOk}
            onCancel={onCancel}
            confirmLoading={confirmLoading}
            width={550}
            okText="确认分配"
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
                                ({dispute.initiatorType === 'user' ? '用户' : '陪玩师'})
                            </div>
                            <div>
                                <strong>纠纷原因:</strong> {dispute.reason}
                            </div>
                        </Space>
                    </div>

                    <Form form={form} layout="vertical">
                        <Form.Item
                            name="assignedServiceId"
                            label="主客服"
                            rules={[{ required: true, message: '请选择主客服' }]}
                            extra="主要负责处理此纠纷的客服人员"
                        >
                            <Select
                                placeholder="请选择主客服"
                                loading={loading}
                                showSearch
                                optionFilterProp="children"
                            >
                                {staffList.map((staff) => (
                                    <Select.Option key={staff.id} value={staff.id}>
                                        {staff.name} {staff.email && `(${staff.email})`}
                                    </Select.Option>
                                ))}
                            </Select>
                        </Form.Item>

                        <Form.Item
                            name="originalServiceId"
                            label="原客服（可选）"
                            extra="如果订单原本有客服，可指定为原客服以保证公平"
                        >
                            <Select
                                placeholder="请选择原客服（可选）"
                                loading={loading}
                                allowClear
                                showSearch
                                optionFilterProp="children"
                            >
                                {staffList
                                    .filter((s) => s.id !== form.getFieldValue('assignedServiceId'))
                                    .map((staff) => (
                                        <Select.Option key={staff.id} value={staff.id}>
                                            {staff.name} {staff.email && `(${staff.email})`}
                                        </Select.Option>
                                    ))}
                            </Select>
                        </Form.Item>
                    </Form>

                    <Divider style={{ margin: '16px 0' }} />

                    {/* Dual-CS Info */}
                    <Alert
                        message="双客服机制说明"
                        description={
                            <div>
                                <p style={{ marginBottom: 8 }}>
                                    根据双客服机制，纠纷处理需要指定两名客服：
                                </p>
                                <ul style={{ margin: 0, paddingLeft: 20 }}>
                                    <li>
                                        <strong>主客服:</strong> 负责具体处理此纠纷
                                    </li>
                                    <li>
                                        <strong>原客服:</strong> 如果订单原本有客服，需指定为原客服以确保公平公正
                                    </li>
                                </ul>
                                <p style={{ marginBottom: 0, marginTop: 8 }}>
                                    <strong>SLA要求:</strong> 分配后，客服需在30分钟内响应并开始处理。
                                </p>
                            </div>
                        }
                        type="info"
                        showIcon
                    />
                </>
            )}
        </Modal>
    );
};

export default AssignModal;
