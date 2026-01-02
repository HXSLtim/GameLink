/**
 * 团队创建/编辑表单
 */
import React, { useEffect } from 'react';
import {
    Modal,
    Form,
    Input,
    InputNumber,
    Select,
    message,
} from 'antd';
import { teamApi, type Team, type TeamCreateRequest, type TeamUpdateRequest } from '@/api/team';

import { logger } from '@/utils/logger';
export interface TeamFormProps {
    /** 是否可见 */
    open: boolean;
    /** 编辑的团队数据（新建时为空） */
    team?: Team | null;
    /** 确认回调 */
    onOk: () => void;
    /** 取消回调 */
    onCancel: () => void;
    /** 加载状态 */
    loading?: boolean;
}

/**
 * 团队表单组件
 */
const TeamForm: React.FC<TeamFormProps> = ({
    open,
    team,
    onOk,
    onCancel,
    loading: outerLoading,
}) => {
    const [form] = Form.useForm();
    const [loading, setLoading] = React.useState(false);
    const isEdit = !!team;

    useEffect(() => {
        if (open) {
            if (team) {
                // 编辑模式：填充表单
                form.setFieldsValue({
                    name: team.name,
                    description: team.description,
                    avatarUrl: team.avatarUrl,
                    maxMembers: team.maxMembers,
                    incomeShareType: team.incomeShareType,
                    leaderBonusRate: team.leaderBonusRate,
                });
            } else {
                // 新建模式：重置表单
                form.resetFields();
                form.setFieldsValue({
                    maxMembers: 10,
                    incomeShareType: 'equal',
                    leaderBonusRate: 0,
                });
            }
        }
    }, [open, team, form]);

    /**
     * 提交表单
     */
    const handleSubmit = async () => {
        try {
            const values = await form.validateFields();

            setLoading(true);
            if (isEdit && team) {
                // 更新
                const updateData: TeamUpdateRequest = {
                    name: values.name,
                    description: values.description,
                    avatarUrl: values.avatarUrl,
                    maxMembers: values.maxMembers,
                    incomeShareType: values.incomeShareType,
                    leaderBonusRate: values.leaderBonusRate,
                };
                await teamApi.updateTeam(team.id, updateData);
                message.success('更新团队成功');
            } else {
                // 新建
                const createData: TeamCreateRequest = {
                    name: values.name,
                    description: values.description,
                    avatarUrl: values.avatarUrl,
                    leaderPlayerId: values.leaderPlayerId,
                    maxMembers: values.maxMembers,
                    incomeShareType: values.incomeShareType,
                    leaderBonusRate: values.leaderBonusRate,
                };
                await teamApi.createTeam(createData);
                message.success('创建团队成功');
            }

            onOk();
        } catch (error) {
            logger.error('Submit team form error:', error);
            message.error(isEdit ? '更新团队失败' : '创建团队失败');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Modal
            title={isEdit ? '编辑团队' : '创建团队'}
            open={open}
            onOk={handleSubmit}
            onCancel={onCancel}
            confirmLoading={loading || outerLoading}
            width={600}
            destroyOnClose
        >
            <Form
                form={form}
                layout="vertical"
                autoComplete="off"
            >
                <Form.Item
                    name="name"
                    label="团队名称"
                    rules={[
                        { required: true, message: '请输入团队名称' },
                        { min: 2, max: 50, message: '团队名称长度为2-50个字符' },
                    ]}
                >
                    <Input placeholder="请输入团队名称" maxLength={50} showCount />
                </Form.Item>

                <Form.Item
                    name="description"
                    label="团队描述"
                    rules={[
                        { max: 500, message: '描述最多500个字符' },
                    ]}
                >
                    <Input.TextArea
                        placeholder="请输入团队描述"
                        rows={4}
                        maxLength={500}
                        showCount
                    />
                </Form.Item>

                <Form.Item
                    name="avatarUrl"
                    label="团队头像URL"
                    rules={[
                        { type: 'url', message: '请输入有效的URL' },
                    ]}
                >
                    <Input placeholder="请输入团队头像URL" />
                </Form.Item>

                {!isEdit && (
                    <Form.Item
                        name="leaderPlayerId"
                        label="队长陪玩师ID"
                        rules={[
                            { required: true, message: '请输入队长陪玩师ID' },
                            { type: 'number', min: 1, message: '请输入有效的陪玩师ID' },
                        ]}
                    >
                        <InputNumber
                            placeholder="请输入队长陪玩师ID"
                            min={1}
                            style={{ width: '100%' }}
                        />
                    </Form.Item>
                )}

                <Form.Item
                    name="maxMembers"
                    label="最大成员数"
                    rules={[
                        { required: true, message: '请输入最大成员数' },
                        { type: 'number', min: 2, max: 50, message: '成员数为2-50' },
                    ]}
                >
                    <InputNumber
                        placeholder="请输入最大成员数"
                        min={2}
                        max={50}
                        style={{ width: '100%' }}
                    />
                </Form.Item>

                <Form.Item
                    name="incomeShareType"
                    label="收益分配方式"
                    rules={[{ required: true, message: '请选择收益分配方式' }]}
                >
                    <Select
                        placeholder="请选择收益分配方式"
                        options={[
                            { label: '平均分配', value: 'equal' },
                            { label: '自定义分配', value: 'custom' },
                        ]}
                    />
                </Form.Item>

                <Form.Item
                    name="leaderBonusRate"
                    label="队长额外加成比例（%）"
                    rules={[
                        { required: true, message: '请输入队长加成比例' },
                        { type: 'number', min: 0, max: 100, message: '加成比例为0-100' },
                    ]}
                    tooltip="队长可以获得的额外收益百分比，0表示无加成"
                >
                    <InputNumber
                        placeholder="请输入队长加成比例"
                        min={0}
                        max={100}
                        precision={2}
                        style={{ width: '100%' }}
                        addonAfter="%"
                    />
                </Form.Item>
            </Form>
        </Modal>
    );
};

export default TeamForm;
