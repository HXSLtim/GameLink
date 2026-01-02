/**
 * 批量分配陪玩师到结算公司弹窗
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Modal,
    Form,
    Select,
    Input,
    DatePicker,
    message,
    Table,
    Space,
    Tag,
    Button,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { UserOutlined, DeleteOutlined } from '@ant-design/icons';
import { settlementApi } from '@/api/settlement';
import type { SettlementCompany } from '@/api/settlement';
import dayjs from 'dayjs';

import { logger } from '@/utils/logger';
export interface AssignModalProps {
    /** 是否可见 */
    open: boolean;
    /** 选中的陪玩师ID列表（用于批量分配） */
    playerIds?: number[];
    /** 确认回调 */
    onOk: () => void;
    /** 取消回调 */
    onCancel: () => void;
    /** 加载状态 */
    loading?: boolean;
}

/**
 * 陪玩师分配项
 */
interface PlayerAssignItem {
    playerId: number;
    nickname: string;
    phone?: string;
    currentCompany?: string;
}

/**
 * 批量分配陪玩师弹窗组件
 */
const AssignModal: React.FC<AssignModalProps> = ({
    open,
    playerIds = [],
    onOk,
    onCancel,
    loading: outerLoading,
}) => {
    const [form] = Form.useForm();
    const [loading, setLoading] = React.useState(false);
    const [companies, setCompanies] = useState<SettlementCompany[]>([]);
    const [players, setPlayers] = useState<PlayerAssignItem[]>([]);
    const [effectiveDate, setEffectiveDate] = useState<string>(dayjs().format('YYYY-MM-DD'));

    /**
     * 加载结算公司列表
     */
    const loadCompanies = async () => {
        try {
            const response = await settlementApi.getSettlementCompanies({
                status: 'active',
                pageSize: 1000,
            });
            if (response.data.success) {
                setCompanies(response.data.data || []);
            }
        } catch (error) {
            logger.error('Load companies error:', error);
        }
    };

    /**
     * 加载陪玩师信息
     */
    const loadPlayers = useCallback(async () => {
        if (playerIds.length === 0) return;

        try {
            // 获取每个陪玩师的当前归属公司
            const playerPromises = playerIds.map(async (playerId) => {
                try {
                    const response = await settlementApi.getPlayerCurrentAssignment(playerId);
                    if (response.data.success && response.data.data) {
                        const assignment = response.data.data;
                        return {
                            playerId,
                            nickname: assignment.player?.nickname || `ID:${playerId}`,
                            phone: assignment.player?.user?.phone,
                            currentCompany: assignment.settlementCompany?.name,
                        };
                    }
                    return {
                        playerId,
                        nickname: `ID:${playerId}`,
                        currentCompany: undefined,
                    };
                } catch {
                    return {
                        playerId,
                        nickname: `ID:${playerId}`,
                        currentCompany: undefined,
                    };
                }
            });

            const playerData = await Promise.all(playerPromises);
            setPlayers(playerData);
        } catch (error) {
            logger.error('Load players error:', error);
        }
    }, [playerIds]);

    useEffect(() => {
        if (open) {
            loadCompanies();
            loadPlayers();
            form.resetFields();
            setEffectiveDate(dayjs().format('YYYY-MM-DD'));
        }
    }, [open, playerIds, form, loadPlayers]);

    /**
     * 提交表单
     */
    const handleSubmit = async () => {
        try {
            const values = await form.validateFields();

            if (playerIds.length === 0) {
                message.warning('请选择要分配的陪玩师');
                return;
            }

            setLoading(true);
            const response = await settlementApi.batchAssignPlayersToCompany({
                playerIds,
                settlementCompanyId: values.settlementCompanyId,
                effectiveDate: values.effectiveDate.format('YYYY-MM-DD'),
                reason: values.reason || '批量分配',
            });

            if (response.data.success) {
                message.success(`成功分配 ${response.data.data.assignedCount} 名陪玩师`);
                onOk();
            }
        } catch (error) {
            logger.error('Assign players error:', error);
            message.error('分配陪玩师失败');
        } finally {
            setLoading(false);
        }
    };

    /**
     * 移除某个陪玩师
     */
    const handleRemovePlayer = (playerId: number) => {
        setPlayers(players.filter(p => p.playerId !== playerId));
    };

    /**
     * 表格列配置
     */
    const columns: ColumnsType<PlayerAssignItem> = [
        {
            title: '陪玩师ID',
            dataIndex: 'playerId',
            key: 'playerId',
            width: 100,
        },
        {
            title: '昵称',
            dataIndex: 'nickname',
            key: 'nickname',
            width: 150,
            render: (nickname: string, _record) => (
                <Space>
                    <UserOutlined />
                    <span>{nickname}</span>
                </Space>
            ),
        },
        {
            title: '手机号',
            dataIndex: 'phone',
            key: 'phone',
            width: 130,
            render: (phone?: string) => phone ? (
                <span>{phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')}</span>
            ) : '-',
        },
        {
            title: '当前归属',
            dataIndex: 'currentCompany',
            key: 'currentCompany',
            width: 150,
            render: (company?: string) => company ? (
                <Tag color="blue">{company}</Tag>
            ) : (
                <Tag color="default">未分配</Tag>
            ),
        },
        {
            title: '操作',
            key: 'action',
            width: 80,
            render: (_, record) => (
                <Button
                    type="link"
                    size="small"
                    danger
                    icon={<DeleteOutlined />}
                    onClick={() => handleRemovePlayer(record.playerId)}
                >
                    移除
                </Button>
            ),
        },
    ];

    // 公司选项（只显示启用状态的公司）
    const companyOptions = companies
        .filter(c => c.status === 'active')
        .map(c => ({
            label: `${c.name} (${c.playerCount}人)`,
            value: c.id,
        }));

    return (
        <Modal
            title={`批量分配陪玩师到结算公司 (${players.length}人)`}
            open={open}
            onOk={handleSubmit}
            onCancel={onCancel}
            confirmLoading={loading || outerLoading}
            width={700}
            okText="确认分配"
        >
            <Form
                form={form}
                layout="vertical"
                autoComplete="off"
                initialValues={{
                    effectiveDate: dayjs(),
                }}
            >
                <Form.Item
                    name="settlementCompanyId"
                    label="目标结算公司"
                    rules={[{ required: true, message: '请选择目标结算公司' }]}
                >
                    <Select
                        placeholder="请选择目标结算公司"
                        options={companyOptions}
                        showSearch
                        filterOption={(input, option) =>
                            (option?.label ?? '').toLowerCase().includes(input.toLowerCase())
                        }
                    />
                </Form.Item>

                <Form.Item
                    name="effectiveDate"
                    label="生效日期"
                    rules={[{ required: true, message: '请选择生效日期' }]}
                    tooltip="归属关系将从该日期开始生效"
                >
                    <DatePicker
                        style={{ width: '100%' }}
                        placeholder="请选择生效日期"
                        disabledDate={(current) => current && current < dayjs().startOf('day')}
                        onChange={(date) => {
                            if (date) {
                                setEffectiveDate(date.format('YYYY-MM-DD'));
                            }
                        }}
                    />
                </Form.Item>

                <Form.Item
                    name="reason"
                    label="分配原因"
                    tooltip="记录分配原因便于后续追溯"
                >
                    <Input.TextArea
                        rows={3}
                        placeholder="请输入分配原因（选填）"
                        maxLength={200}
                        showCount
                    />
                </Form.Item>

                <div style={{ marginBottom: 16 }}>
                    <div style={{ fontWeight: 500, marginBottom: 8 }}>
                        待分配陪玩师列表：
                    </div>
                    <Table
                        columns={columns}
                        dataSource={players}
                        rowKey="playerId"
                        size="small"
                        pagination={false}
                        scroll={{ y: 200 }}
                    />
                </div>

                {effectiveDate && (
                    <div style={{
                        padding: 12,
                        background: '#f0f0f0',
                        borderRadius: 4,
                        fontSize: 12,
                        color: '#666',
                    }}>
                        提示：陪玩师归属关系将从 {effectiveDate} 开始生效，生效后订单结算将分配到目标公司
                    </div>
                )}
            </Form>
        </Modal>
    );
};

export default AssignModal;
