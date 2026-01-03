/**
 * 游戏管理页面
 */
import React, { useState, useCallback } from 'react';
import {
    Tag,
    Space,
    Button,
    Modal,
    Form,
    Input,
    Select,
    Radio,
    Image,
    message,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    EditOutlined,
    DeleteOutlined,
    DownloadOutlined,
} from '@ant-design/icons';
import { exportToCSV, gameExportColumns } from '@/utils/export';
import { PageContainer, SearchTable, type ToolbarButton } from '@/components';
import type { SearchField } from '@/components';
import { GAME_PERMISSIONS } from '@/constants/permissions';
import { PermissionGuard } from '@/components/PermissionGuard';
import { adminApi, type Game, type CreateGameDto, type UpdateGameDto } from '@/api/admin';
import { useCrud } from '@/hooks';
import dayjs from 'dayjs';

/**
 * 分类选项
 */
const categoryOptions = [
    { label: 'MOBA', value: 'moba' },
    { label: '射击', value: 'fps' },
    { label: 'RPG', value: 'rpg' },
    { label: '卡牌', value: 'card' },
    { label: '休闲', value: 'casual' },
    { label: '其他', value: 'other' },
];

/**
 * 游戏管理页面
 */
const GamePage: React.FC = () => {
    // 弹窗状态
    const [editModalVisible, setEditModalVisible] = useState(false);
    const [currentGame, setCurrentGame] = useState<Game | null>(null);
    const [form] = Form.useForm();
    const [submitting, setSubmitting] = useState(false);

    // 批量操作状态
    const [batchDeleteVisible, setBatchDeleteVisible] = useState(false);
    const [selectedGameIds, setSelectedGameIds] = useState<number[]>([]);
    const [batchTarget, setBatchTarget] = useState<'selected' | 'category' | 'all'>('selected');
    const [batchForm] = Form.useForm();

    /**
     * 使用 CRUD Hook 管理游戏数据
     */
    const {
        data: games,
        loading,
        pagination,
        fetchAll,
        create,
        update,
        remove,
        setSearchParams,
        queryParams,
    } = useCrud<Game, CreateGameDto, UpdateGameDto>({
        api: {
            getAll: adminApi.getGames as any,
            create: adminApi.createGame as any,
            update: adminApi.updateGame as any,
            remove: adminApi.deleteGame as any,
        },
        messages: {
            fetchError: '加载游戏列表失败',
            createSuccess: '创建游戏成功',
            updateSuccess: '更新游戏成功',
            deleteSuccess: '删除游戏成功',
        },
        initialPagination: {
            pageSize: 10,
        },
    });

    /**
     * 搜索
     */
    const handleSearch = useCallback((values: Record<string, unknown>) => {
        setSearchParams(values);
    }, [setSearchParams]);

    /**
     * 编辑游戏
     */
    const handleEdit = useCallback((game: Game) => {
        setCurrentGame(game);
        form.setFieldsValue({
            key: game.key,
            name: game.name,
            category: game.category,
            description: game.description,
            icon_url: game.iconUrl,
        });
        setEditModalVisible(true);
    }, [form]);

    /**
     * 新增游戏
     */
    const handleCreate = useCallback(() => {
        setCurrentGame(null);
        form.resetFields();
        setEditModalVisible(true);
    }, [form]);

    /**
     * 保存编辑
     */
    const handleSaveEdit = useCallback(async () => {
        try {
            const values = await form.validateFields();
            setSubmitting(true);

            const data = {
                key: values.key,
                name: values.name,
                category: values.category || '',
                description: values.description || '',
                icon_url: values.icon_url || '',
            };

            if (currentGame) {
                await update(String(currentGame.id), data as UpdateGameDto);
            } else {
                await create(data as CreateGameDto);
            }

            setEditModalVisible(false);
        } catch (err) {
            // Form validation error or API error (handled by hook)
            if (err && typeof err === 'object' && 'errorFields' in err) {
                // Form validation error from Ant Design
                console.error('Form validation error:', err);
            } else {
                console.error('Save error:', err);
            }
        } finally {
            setSubmitting(false);
        }
    }, [currentGame, form, update, create]);

    /**
     * 删除游戏
     */
    const handleDelete = useCallback(async (game: Game) => {
        await remove(game.id, {
            confirmMessage: `确定要删除游戏 "${game.name}" 吗？`,
        });
    }, [remove]);

    /**
     * 批量删除
     */
    const handleBatchDelete = useCallback((keys: React.Key[]) => {
        setSelectedGameIds(keys ? keys.map(k => Number(k)) : []);
        batchForm.resetFields();
        batchForm.setFieldsValue({
            target: (keys && keys.length > 0) ? 'selected' : 'all',
        });
        setBatchTarget((keys && keys.length > 0) ? 'selected' : 'all');
        setBatchDeleteVisible(true);
    }, []);

    const submitBatchDelete = useCallback(async () => {
        try {
            const values = await batchForm.validateFields();
            let gameIds: number[] = [];

            if (values.target === 'selected') {
                gameIds = selectedGameIds;
            } else if (values.target === 'category') {
                // 按分类筛选
                const response = await adminApi.getGames({ page_size: 1000 });
                if (response.data.success && response.data.data) {
                    gameIds = response.data.data
                        .filter((g: Game) => g.category === values.filterCategory)
                        .map((g: Game) => g.id);
                }
            } else {
                // 全部
                const response = await adminApi.getGames({ page_size: 1000 });
                if (response.data.success && response.data.data) {
                    gameIds = response.data.data.map((g: Game) => g.id);
                }
            }

            if (gameIds.length === 0) {
                Modal.warning({
                    title: '提示',
                    content: '没有符合条件的游戏',
                });
                return;
            }

            const res = await adminApi.batchDeleteGames(gameIds);

            if (res.data.success) {
                Modal.success({
                    title: '成功',
                    content: `批量删除 ${gameIds.length} 个游戏成功`,
                });
                setBatchDeleteVisible(false);
                fetchAll();
            }
        } catch (error) {
            console.error('Batch delete error:', error);
            Modal.error({
                title: '失败',
                content: '操作失败',
            });
        }
    }, [batchForm, selectedGameIds, fetchAll]);

    /**
     * 搜索字段配置
     */
    const searchFields: SearchField[] = [
        { name: 'keyword', label: '游戏名称', type: 'input', placeholder: '请输入游戏名称' },
        { name: 'category', label: '分类', type: 'select', options: categoryOptions },
    ];

    /**
     * 表格列配置
     */
    const columns: ColumnsType<Game> = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 80,
        },
        {
            title: '游戏',
            key: 'game',
            width: 200,
            render: (_, record) => (
                <Space>
                    <Image
                        src={record.iconUrl || 'https://via.placeholder.com/40'}
                        width={40}
                        height={40}
                        style={{ borderRadius: 8 }}
                        fallback="https://via.placeholder.com/40"
                    />
                    <span style={{ fontWeight: 500 }}>{record.name}</span>
                </Space>
            ),
        },
        {
            title: 'Key',
            dataIndex: 'key',
            key: 'key',
            width: 120,
        },
        {
            title: '分类',
            dataIndex: 'category',
            key: 'category',
            width: 100,
            render: category => {
                const option = categoryOptions.find(o => o.value === category);
                return <Tag>{option?.label || category || '-'}</Tag>;
            },
        },
        {
            title: '描述',
            dataIndex: 'description',
            key: 'description',
            width: 200,
            ellipsis: true,
            render: desc => desc || '-',
        },
        {
            title: '创建时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 180,
            render: date => date ? dayjs(date).format('YYYY-MM-DD HH:mm:ss') : '-',
        },
        {
            title: '更新时间',
            dataIndex: 'updatedAt',
            key: 'updatedAt',
            width: 180,
            render: date => date ? dayjs(date).format('YYYY-MM-DD HH:mm:ss') : '-',
        },
        {
            title: '操作',
            key: 'action',
            width: 150,
            fixed: 'right',
            render: (_, record) => (
                <Space size="small">
                    <PermissionGuard permission={GAME_PERMISSIONS.UPDATE}>
                        <Button
                            type="link"
                            size="small"
                            icon={<EditOutlined />}
                            onClick={() => handleEdit(record)}
                        >
                            编辑
                        </Button>
                    </PermissionGuard>
                    <PermissionGuard permission={GAME_PERMISSIONS.DELETE}>
                        <Button
                            type="link"
                            size="small"
                            danger
                            icon={<DeleteOutlined />}
                            onClick={() => {
                                Modal.confirm({
                                    title: '确认删除',
                                    content: `确定要删除游戏 "${record.name}" 吗？`,
                                    onOk: () => handleDelete(record),
                                });
                            }}
                        >
                            删除
                        </Button>
                    </PermissionGuard>
                </Space>
            ),
        },
    ];

    /**
     * 导出游戏数据
     */
    const handleExport = useCallback(async () => {
        try {
            message.loading({ content: '正在导出...', key: 'export' });
            const response = await adminApi.getGames({ ...queryParams, page_size: 10000 });
            if (response.data.success && response.data.data) {
                exportToCSV(response.data.data as unknown as Record<string, unknown>[], gameExportColumns, 'games');
                message.success({ content: '导出成功', key: 'export' });
            } else {
                message.error({ content: '导出失败', key: 'export' });
            }
        } catch {
            message.error({ content: '导出失败', key: 'export' });
        }
    }, [queryParams]);

    /**
     * 工具栏按钮
     */
    const toolbarButtons: ToolbarButton[] = [
        {
            text: '批量删除',
            icon: <DeleteOutlined />,
            needSelection: false,
            danger: true,
            onClick: (keys) => handleBatchDelete(keys || []),
            permission: GAME_PERMISSIONS.DELETE,
        },
        {
            text: '导出数据',
            icon: <DownloadOutlined />,
            needSelection: false,
            onClick: () => handleExport(),
            permission: GAME_PERMISSIONS.LIST,
        },
    ];

    return (
        <PageContainer title="游戏管理" subTitle="管理平台支持的游戏">
            <SearchTable
                columns={columns}
                dataSource={games}
                rowKey="id"
                searchFields={searchFields}
                onSearch={handleSearch}
                onRefresh={() => fetchAll()}
                loading={loading}
                showCreate={true}
                createText="新增游戏"
                createPermission={GAME_PERMISSIONS.CREATE}
                onCreate={handleCreate}
                toolbarButtons={toolbarButtons}
                pagination={pagination}
                scroll={{ x: 1200 }}
            />

            {/* 编辑弹窗 */}
            <Modal
                title={currentGame ? '编辑游戏' : '新增游戏'}
                open={editModalVisible}
                onOk={handleSaveEdit}
                onCancel={() => setEditModalVisible(false)}
                confirmLoading={submitting}
                width={600}
            >
                <Form form={form} layout="vertical">
                    <Form.Item
                        name="key"
                        label="游戏Key"
                        rules={[{ required: true, message: '请输入游戏Key' }]}
                    >
                        <Input placeholder="请输入游戏Key（唯一标识）" disabled={!!currentGame} />
                    </Form.Item>

                    <Form.Item
                        name="name"
                        label="游戏名称"
                        rules={[{ required: true, message: '请输入游戏名称' }]}
                    >
                        <Input placeholder="请输入游戏名称" />
                    </Form.Item>

                    <Form.Item name="icon_url" label="游戏图标URL">
                        <Input placeholder="请输入图标URL" />
                    </Form.Item>

                    <Form.Item name="category" label="分类">
                        <Select placeholder="请选择分类" options={categoryOptions} allowClear />
                    </Form.Item>

                    <Form.Item name="description" label="描述">
                        <Input.TextArea rows={3} placeholder="请输入游戏描述" />
                    </Form.Item>
                </Form>
            </Modal>

            {/* 批量删除弹窗 */}
            <Modal
                title="批量删除游戏"
                open={batchDeleteVisible}
                onOk={submitBatchDelete}
                onCancel={() => setBatchDeleteVisible(false)}
                okText="确认删除"
                okButtonProps={{ danger: true }}
            >
                <Form form={batchForm} layout="vertical">
                    <Form.Item name="target" label="目标对象" rules={[{ required: true }]}>
                        <Radio.Group onChange={(e) => setBatchTarget(e.target.value)}>
                            <Radio value="selected" disabled={selectedGameIds.length === 0}>
                                选中的游戏 {selectedGameIds.length > 0 ? `(${selectedGameIds.length})` : ''}
                            </Radio>
                            <Radio value="category">按分类筛选</Radio>
                            <Radio value="all">全部游戏</Radio>
                        </Radio.Group>
                    </Form.Item>

                    {batchTarget === 'category' && (
                        <Form.Item name="filterCategory" label="筛选分类" rules={[{ required: true, message: '请选择筛选分类' }]}>
                            <Select placeholder="请选择要筛选的分类" options={categoryOptions} />
                        </Form.Item>
                    )}

                    <div style={{ color: '#ff4d4f', marginTop: 16 }}>
                        ⚠️ 警告：此操作不可恢复，请谨慎操作！
                    </div>
                </Form>
            </Modal>
        </PageContainer>
    );
};

export default GamePage;
