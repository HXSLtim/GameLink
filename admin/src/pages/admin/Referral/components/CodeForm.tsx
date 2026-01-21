/**
 * Code Form Component
 * 邀请码表单组件
 *
 * Modal form for creating and editing referral codes.
 */
import React, { useEffect } from 'react';
import {
    Modal,
    Form,
    Select,
    InputNumber,
    DatePicker,
} from 'antd';
import dayjs from 'dayjs';
import { type ReferralCode } from '@/api/referral';
import type { CreateReferralCodeDto, UpdateReferralCodeDto } from '@/api/referral';

interface CodeFormProps {
    visible: boolean;
    code: ReferralCode | null;
    onSave: (values: CreateReferralCodeDto | UpdateReferralCodeDto) => void;
    onCancel: () => void;
}

// NOTE: UserSearchSelect component removed as it's currently unused.
// To implement user search, create a component that uses an API to search users by name/phone/email.

const CodeForm: React.FC<CodeFormProps> = ({
    visible,
    code,
    onSave,
    onCancel,
}) => {
    const [form] = Form.useForm();
    const isEdit = !!code;

    /**
     * Initialize form values when editing
     */
    useEffect(() => {
        if (code) {
            form.setFieldsValue({
                userId: code.ownerId,
                type: code.type,
                maxUse: code.maxUses,
                expireAt: dayjs(code.expiresAt),
                isActive: code.isActive,
            });
        } else {
            form.resetFields();
            // Set default values
            form.setFieldsValue({
                type: 'user',
                maxUse: 100,
                expireAt: dayjs().add(30, 'day'),
                isActive: true,
            });
        }
    }, [code, form, visible]);

    /**
     * Handle form submit
     */
    const handleOk = async () => {
        try {
            const values = await form.validateFields();

            if (isEdit) {
                // Update
                const updateData: UpdateReferralCodeDto = {
                    isActive: values.isActive,
                    maxUse: values.maxUse,
                    expireAt: values.expireAt.format('YYYY-MM-DD HH:mm:ss'),
                };
                onSave(updateData);
            } else {
                // Create
                const createData: CreateReferralCodeDto = {
                    userId: values.userId,
                    type: values.type,
                    maxUse: values.maxUse,
                    expireAt: values.expireAt.format('YYYY-MM-DD HH:mm:ss'),
                };
                onSave(createData);
            }
        } catch {
            // Form validation failed
        }
    };

    /**
     * Disable date before tomorrow for expire time
     */
    const disabledDate = (current: dayjs.Dayjs) => {
        return current && current < dayjs().startOf('day');
    };

    return (
        <Modal
            title={isEdit ? '编辑邀请码' : '创建邀请码'}
            open={visible}
            onOk={handleOk}
            onCancel={onCancel}
            width={500}
            destroyOnHidden
        >
            <Form
                form={form}
                layout="vertical"
                preserve={false}
            >
                {!isEdit && (
                    <Form.Item
                        name="userId"
                        label="拥有者"
                        rules={[{ required: true, message: '请选择拥有者' }]}
                    >
                        <InputNumber
                            placeholder="请输入用户ID"
                            min={1}
                            style={{ width: '100%' }}
                        />
                        {/* Note: Using simple input for user ID. You can replace with UserSearchSelect if needed */}
                    </Form.Item>
                )}

                <Form.Item
                    name="type"
                    label="类型"
                    rules={[{ required: true, message: '请选择类型' }]}
                >
                    <Select
                        placeholder="请选择邀请码类型"
                        disabled={isEdit}
                        options={[
                            { value: 'user', label: '用户推荐' },
                            { value: 'player', label: '陪玩师推荐' },
                        ]}
                    />
                </Form.Item>

                <Form.Item
                    name="maxUse"
                    label="最大使用次数"
                    rules={[
                        { required: true, message: '请输入最大使用次数' },
                        { type: 'number', min: 1, message: '至少为1次' },
                    ]}
                >
                    <InputNumber
                        placeholder="请输入最大使用次数"
                        min={1}
                        max={10000}
                        style={{ width: '100%' }}
                    />
                </Form.Item>

                <Form.Item
                    name="expireAt"
                    label="有效期至"
                    rules={[{ required: true, message: '请选择有效期' }]}
                >
                    <DatePicker
                        showTime
                        placeholder="请选择有效期"
                        style={{ width: '100%' }}
                        disabledDate={disabledDate}
                        format="YYYY-MM-DD HH:mm:ss"
                    />
                </Form.Item>

                {isEdit && (
                    <Form.Item
                        name="isActive"
                        label="状态"
                        rules={[{ required: true, message: '请选择状态' }]}
                    >
                        <Select
                            placeholder="请选择状态"
                            options={[
                                { value: true, label: '启用' },
                                { value: false, label: '禁用' },
                            ]}
                        />
                    </Form.Item>
                )}
            </Form>
        </Modal>
    );
};

export default CodeForm;
