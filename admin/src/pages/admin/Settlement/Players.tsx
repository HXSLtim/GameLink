/**
 * 陪玩师归属管理页面
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Tag,
    Space,
    Button,
    Modal,
    message,
    Card,
    Statistic,
    Select,
    Form,
    Input,
    Drawer,
    List,
    theme,
    Avatar,
    Typography,
    Tooltip,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    UserOutlined,
    BankOutlined,
    SwapOutlined,
    HistoryOutlined,
    TeamOutlined,
} from '@ant-design/icons';
import { PageContainer, SearchTable, type ToolbarButton } from '@/components';
import type { SearchField } from '@/components';
import { settlementApi } from '@/api/settlement';
import type {
    PlayerCompanyAssignment,
    SettlementCompany,
    CompanyType,
} from '@/api/settlement';
import AssignModal from './components/AssignModal';
import dayjs from 'dayjs';

const { Text } = Typography;

/**
 * 公司类型映射
 */
const companyTypeMap: Record<CompanyType, { text: string; color: string }> = {
    individual: { text: '个人', color: 'blue' },
    company: { text: '企业', color: 'green' },
};

/**
 * 陪玩师归属管理页面
 */
const SettlementPlayersPage: React.FC = () => {
    const { token } = theme.useToken();
    const [loading, setLoading] = useState(false);
    const [players, setPlayers] = useState<PlayerCompanyAssignment[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const [searchParams, setSearchParams] = useState<Record<string, unknown>>({});

    // 弹窗状态
    const [assignModalVisible, setAssignModalVisible] = useState(false);
    const [historyDrawerVisible, setHistoryDrawerVisible] = useState(false);
    const [transferModalVisible, setTransferModalVisible] = useState(false);
    const [selectedPlayerIds, setSelectedPlayerIds] = useState<number[]>([]);
    const [currentPlayer, setCurrentPlayer] = useState<PlayerCompanyAssignment | null>(null);
    const [playerHistory, setPlayerHistory] = useState<PlayerCompanyAssignment[]>([]);
    const [companies, setCompanies] = useState<SettlementCompany[]>([]);

    const [assignForm] = Form.useForm();

    /**
     * 加载陪玩师归属数据
     */
    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            // 由于API返回的是单个公司的陪玩师列表，这里需要通过其他方式获取
            // 暂时使用空数据，实际需要后端提供专门的陪玩师归属列表API
            const response = await settlementApi.getSettlementCompanies({
                pageSize: 100,
            });
            if (response.data.success) {
                const companyData = response.data.data || [];
                setCompanies(companyData);

                // 展开所有公司的陪玩师
                const allPlayers: PlayerCompanyAssignment[] = [];
                for (const company of companyData) {
                    const playersResponse = await settlementApi.getCompanyPlayers(company.id, {
                        pageSize: 100,
                    });
                    if (playersResponse.data.success) {
                        allPlayers.push(...(playersResponse.data.data || []));
                    }
                }
                setPlayers(allPlayers);
                setTotal(allPlayers.length);
            }
        } catch (error) {
            console.error('Load players error:', error);
            message.error('加载陪玩师归属列表失败');
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        loadData();
    }, [loadData]);

    /**
     * 搜索
     */
    const handleSearch = (values: Record<string, unknown>) => {
        setSearchParams(values);
        setCurrent(1);
    };

    /**
     * 批量分配陪玩师
     */
    const handleBatchAssign = (keys: React.Key[]) => {
        setSelectedPlayerIds(keys ? keys.map(k => Number(k)) : []);
        setAssignModalVisible(true);
    };

    /**
     * 单个陪玩师转移归属
     */
    const handleTransfer = (player: PlayerCompanyAssignment) => {
        setCurrentPlayer(player);
        assignForm.resetFields();
        setTransferModalVisible(true);
    };

    const submitTransfer = async () => {
        if (!currentPlayer) return;
        try {
            const values = await assignForm.validateFields();
            await settlementApi.assignPlayerToCompany(currentPlayer.playerId, {
                settlementCompanyId: values.settlementCompanyId,
                effectiveDate: values.effectiveDate.format('YYYY-MM-DD'),
                reason: values.reason || '管理员调整',
            });
            message.success('转移归属成功');
            setTransferModalVisible(false);
            loadData();
        } catch (error) {
            console.error('Transfer error:', error);
            message.error('转移归属失败');
        }
    };

    /**
     * 查看陪玩师归属历史
     */
    const handleViewHistory = async (player: PlayerCompanyAssignment) => {
        try {
            const response = await settlementApi.getPlayerAssignmentHistory(player.playerId, {
                page: 1,
                pageSize: 20,
            });
            if (response.data.success) {
                setPlayerHistory(response.data.data.assignments || []);
                setCurrentPlayer(player);
                setHistoryDrawerVisible(true);
            }
        } catch (error) {
            console.error('Load history error:', error);
            message.error('加载历史记录失败');
        }
    };

    /**
     * 搜索字段配置
     */
    const searchFields: SearchField[] = [
        { name: 'keyword', label: '关键词', type: 'input', placeholder: '陪玩师昵称/ID' },
        {
            name: 'settlementCompanyId',
            label: '归属公司',
            type: 'select',
            options: companies
                .filter(c => c.status === 'active')
                .map(c => ({ label: c.name, value: c.id })),
        },
    ];

    /**
     * 表格列配置
     */
    const columns: ColumnsType<PlayerCompanyAssignment> = [
        {
            title: '陪玩师ID',
            dataIndex: 'playerId',
            key: 'playerId',
            width: 100,
        },
        {
            title: '陪玩师',
            key: 'player',
            width: 200,
            render: (_, record) => (
                <Space>
                    <Avatar size={32} icon={<UserOutlined />} />
                    <div>
                        <div style={{ fontWeight: 500 }}>
                            {record.player?.nickname || `ID:${record.playerId}`}
                        </div>
                        {record.player?.user?.phone && (
                            <Text type="secondary" style={{ fontSize: 12 }}>
                                {record.player.user.phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')}
                            </Text>
                        )}
                    </div>
                </Space>
            ),
        },
        {
            title: '归属公司',
            key: 'company',
            width: 200,
            render: (_, record) => {
                if (record.settlementCompany) {
                    return (
                        <Space>
                            <BankOutlined />
                            <span>{record.settlementCompany.name}</span>
                            <Tag color={companyTypeMap[record.settlementCompany.type].color} style={{ fontSize: 11 }}>
                                {companyTypeMap[record.settlementCompany.type].text}
                            </Tag>
                        </Space>
                    );
                }
                return <Tag color="default">未分配</Tag>;
            },
        },
        {
            title: '生效日期',
            dataIndex: 'effectiveDate',
            key: 'effectiveDate',
            width: 120,
            render: (date: string) => dayjs(date).format('YYYY-MM-DD'),
        },
        {
            title: '分配原因',
            dataIndex: 'reason',
            key: 'reason',
            width: 150,
            ellipsis: { showTitle: false },
            render: (reason: string) => (
                <Tooltip title={reason}>
                    <span>{reason || '-'}</span>
                </Tooltip>
            ),
        },
        {
            title: '分配时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 180,
            render: (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss'),
        },
        {
            title: '操作',
            key: 'action',
            width: 200,
            fixed: 'right',
            render: (_, record) => (
                <Space size="small">
                    <Button
                        type="link"
                        size="small"
                        icon={<SwapOutlined />}
                        onClick={() => handleTransfer(record)}
                    >
                        转移
                    </Button>
                    <Button
                        type="link"
                        size="small"
                        icon={<HistoryOutlined />}
                        onClick={() => handleViewHistory(record)}
                    >
                        历史
                    </Button>
                </Space>
            ),
        },
    ];

    /**
     * 工具栏按钮
     */
    const toolbarButtons: ToolbarButton[] = [
        {
            text: '批量分配',
            icon: <TeamOutlined />,
            needSelection: true,
            onClick: (keys) => handleBatchAssign(keys || []),
        },
    ];

    // 统计数据
    const stats = {
        total: total,
        assigned: players.filter(p => p.settlementCompany).length,
        unassigned: players.filter(p => !p.settlementCompany).length,
    };

    // 公司选项（只显示启用状态的公司）
    const companyOptions = companies
        .filter(c => c.status === 'active')
        .map(c => ({
            label: c.name,
            value: c.id,
        }));

    return (
        <PageContainer
            title="陪玩师归属管理"
            subTitle="管理陪玩师与结算公司的归属关系"
            extra={
                <Space size="large">
                    <Statistic title="总陪玩师" value={stats.total} prefix={<UserOutlined />} />
                    <Statistic title="已分配" value={stats.assigned} valueStyle={{ color: token.colorSuccess }} />
                    <Statistic title="未分配" value={stats.unassigned} valueStyle={{ color: token.colorWarning }} />
                </Space>
            }
        >
            <SearchTable
                columns={columns}
                dataSource={players}
                rowKey="id"
                searchFields={searchFields}
                onSearch={handleSearch}
                onRefresh={() => loadData()}
                loading={loading}
                showCreate={false}
                toolbarButtons={toolbarButtons}
                pagination={{
                    current,
                    pageSize,
                    total,
                    showSizeChanger: true,
                    showQuickJumper: true,
                    showTotal: t => `共 ${t} 条`,
                    onChange: (page, size) => {
                        setCurrent(page);
                        setPageSize(size);
                    },
                }}
                scroll={{ x: 1400 }}
            />

            {/* 批量分配弹窗 */}
            <AssignModal
                open={assignModalVisible}
                playerIds={selectedPlayerIds}
                onOk={() => {
                    setAssignModalVisible(false);
                    setSelectedPlayerIds([]);
                    loadData();
                }}
                onCancel={() => {
                    setAssignModalVisible(false);
                    setSelectedPlayerIds([]);
                }}
            />

            {/* 单个转移弹窗 */}
            <Modal
                title="转移陪玩师归属"
                open={transferModalVisible}
                onOk={submitTransfer}
                onCancel={() => setTransferModalVisible(false)}
                width={500}
            >
                {currentPlayer && (
                    <Form form={assignForm} layout="vertical">
                        <Card size="small" style={{ marginBottom: 16 }}>
                            <Space>
                                <Avatar size={40} icon={<UserOutlined />} />
                                <div>
                                    <div style={{ fontWeight: 500 }}>
                                        {currentPlayer.player?.nickname || `ID:${currentPlayer.playerId}`}
                                    </div>
                                    <Text type="secondary" style={{ fontSize: 12 }}>
                                        当前归属: {currentPlayer.settlementCompany?.name || '未分配'}
                                    </Text>
                                </div>
                            </Space>
                        </Card>

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
                            initialValue={dayjs()}
                        >
                            <Select
                                style={{ width: '100%' }}
                                placeholder="请选择生效日期"
                            >
                                <Select.Option value={dayjs().format('YYYY-MM-DD')}>立即生效</Select.Option>
                                <Select.Option value={dayjs().add(1, 'day').format('YYYY-MM-DD')}>
                                    次日生效
                                </Select.Option>
                                <Select.Option value={dayjs().add(1, 'month').format('YYYY-MM-DD')}>
                                    次月生效
                                </Select.Option>
                            </Select>
                        </Form.Item>

                        <Form.Item
                            name="reason"
                            label="转移原因"
                        >
                            <Input.TextArea
                                rows={3}
                                placeholder="请输入转移原因（选填）"
                                maxLength={200}
                                showCount
                            />
                        </Form.Item>
                    </Form>
                )}
            </Modal>

            {/* 历史记录抽屉 */}
            <Drawer
                title="陪玩师归属历史"
                open={historyDrawerVisible}
                onClose={() => setHistoryDrawerVisible(false)}
                width={600}
            >
                {currentPlayer && (
                    <>
                        <Card size="small" style={{ marginBottom: 16 }}>
                            <Space>
                                <Avatar size={48} icon={<UserOutlined />} />
                                <div>
                                    <div style={{ fontSize: 16, fontWeight: 500 }}>
                                        {currentPlayer.player?.nickname || `ID:${currentPlayer.playerId}`}
                                    </div>
                                    <Text type="secondary">
                                        陪玩师ID: {currentPlayer.playerId}
                                    </Text>
                                </div>
                            </Space>
                        </Card>

                        <div style={{ marginBottom: 16, fontWeight: 500 }}>归属变更记录：</div>

                        <List
                            dataSource={playerHistory}
                            renderItem={(item) => (
                                <List.Item>
                                    <List.Item.Meta
                                        avatar={<Avatar icon={<BankOutlined />} />}
                                        title={
                                            <Space>
                                                <span>{item.settlementCompany?.name || '未分配'}</span>
                                                {item.settlementCompany && (
                                                    <Tag color={companyTypeMap[item.settlementCompany.type].color}>
                                                        {companyTypeMap[item.settlementCompany.type].text}
                                                    </Tag>
                                                )}
                                            </Space>
                                        }
                                        description={
                                            <div>
                                                <div>生效日期: {dayjs(item.effectiveDate).format('YYYY-MM-DD')}</div>
                                                <div>分配原因: {item.reason || '-'}</div>
                                                <div>
                                                    分配时间: {dayjs(item.createdAt).format('YYYY-MM-DD HH:mm:ss')}
                                                </div>
                                                {item.assignedByAdmin && (
                                                    <div>操作人: {item.assignedByAdmin.name}</div>
                                                )}
                                            </div>
                                        }
                                    />
                                </List.Item>
                            )}
                        />
                    </>
                )}
            </Drawer>
        </PageContainer>
    );
};

export default SettlementPlayersPage;
