/**
 * 游戏管理页面
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Tag,
    Space,
    Button,
    Modal,
    Form,
    Input,
    Select,
    message,
    Popconfirm,
    Switch,
    Upload,
    Image,
    InputNumber,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import type { UploadFile } from 'antd';
import {
    EditOutlined,
    DeleteOutlined,
    PlusOutlined,
} from '@ant-design/icons';
import { PageContainer, SearchTable } from '@/components';
import type { SearchField } from '@/components';
import { GAME_PERMISSIONS } from '@/constants/permissions';
import { PermissionGuard } from '@/components/PermissionGuard';

/**
 * 游戏数据接口
 */
interface Game {
    id: number;
    name: string;
    icon: string;
    category: string;
    platform: string[];
    status: 'active' | 'inactive';
    playerCount: number;
    orderCount: number;
    sortOrder: number;
    createdAt: string;
    updatedAt: string;
}

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
 * 平台选项
 */
const platformOptions = [
    { label: 'PC', value: 'pc' },
    { label: 'iOS', value: 'ios' },
    { label: 'Android', value: 'android' },
    { label: '主机', value: 'console' },
];

/**
 * 游戏管理页面
 */
const GamePage: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [games, setGames] = useState<Game[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);

    // 弹窗状态
    const [editModalVisible, setEditModalVisible] = useState(false);
    const [currentGame, setCurrentGame] = useState<Game | null>(null);
    const [form] = Form.useForm();
    const [fileList, setFileList] = useState<UploadFile[]>([]);

    /**
     * 加载游戏数据
     */
    const loadData = useCallback(async () => {
        setLoading(true);
        await new Promise(resolve => setTimeout(resolve, 500));

        const mockGames: Game[] = [
            { id: 1, name: '王者荣耀', icon: 'https://via.placeholder.com/48', category: 'moba', platform: ['ios', 'android'], status: 'active', playerCount: 256, orderCount: 1280, sortOrder: 1, createdAt: '2024-01-01 00:00:00', updatedAt: '2024-11-20 10:30:00' },
            { id: 2, name: '英雄联盟', icon: 'https://via.placeholder.com/48', category: 'moba', platform: ['pc'], status: 'active', playerCount: 189, orderCount: 920, sortOrder: 2, createdAt: '2024-01-01 00:00:00', updatedAt: '2024-11-18 15:20:00' },
            { id: 3, name: '和平精英', icon: 'https://via.placeholder.com/48', category: 'fps', platform: ['ios', 'android'], status: 'active', playerCount: 145, orderCount: 680, sortOrder: 3, createdAt: '2024-01-05 10:00:00', updatedAt: '2024-11-15 09:00:00' },
            { id: 4, name: '原神', icon: 'https://via.placeholder.com/48', category: 'rpg', platform: ['pc', 'ios', 'android'], status: 'active', playerCount: 98, orderCount: 450, sortOrder: 4, createdAt: '2024-02-01 14:30:00', updatedAt: '2024-11-10 11:45:00' },
            { id: 5, name: '永劫无间', icon: 'https://via.placeholder.com/48', category: 'fps', platform: ['pc'], status: 'active', playerCount: 67, orderCount: 280, sortOrder: 5, createdAt: '2024-03-15 09:00:00', updatedAt: '2024-11-05 16:30:00' },
            { id: 6, name: '炉石传说', icon: 'https://via.placeholder.com/48', category: 'card', platform: ['pc', 'ios', 'android'], status: 'inactive', playerCount: 23, orderCount: 85, sortOrder: 10, createdAt: '2024-04-20 11:00:00', updatedAt: '2024-10-01 10:00:00' },
        ];

        setGames(mockGames);
        setTotal(mockGames.length);
        setLoading(false);
    }, []);

    useEffect(() => {
        // eslint-disable-next-line react-hooks/set-state-in-effect
        loadData();
    }, [loadData]);

    /**
     * 编辑游戏
     */
    const handleEdit = (game: Game) => {
        setCurrentGame(game);
        form.setFieldsValue({
            ...game,
        });
        if (game.icon) {
            setFileList([{ uid: '-1', name: 'icon', status: 'done', url: game.icon }]);
        } else {
            setFileList([]);
        }
        setEditModalVisible(true);
    };

    /**
     * 保存编辑
     */
    const handleSaveEdit = async () => {
        try {
            const values = await form.validateFields();
            console.log('Save:', values);
            message.success('保存成功');
            setEditModalVisible(false);
            loadData();
        } catch {
            // 验证失败
        }
    };

    /**
     * 切换状态
     */
    const handleToggleStatus = async (game: Game) => {
        message.success(`${game.status === 'active' ? '下架' : '上架'}成功`);
        loadData();
    };

    /**
     * 删除游戏
     */
    const handleDelete = async (game: Game) => {
        message.success(`删除游戏 ${game.name} 成功`);
        loadData();
    };

    /**
     * 搜索字段配置
     */
    const searchFields: SearchField[] = [
        { name: 'keyword', label: '游戏名称', type: 'input', placeholder: '请输入游戏名称' },
        { name: 'category', label: '分类', type: 'select', options: categoryOptions },
        {
            name: 'status',
            label: '状态',
            type: 'select',
            options: [
                { label: '已上架', value: 'active' },
                { label: '已下架', value: 'inactive' },
            ],
        },
    ];

    /**
     * 表格列配置
     */
    const columns: ColumnsType<Game> = [
        {
            title: '排序',
            dataIndex: 'sortOrder',
            key: 'sortOrder',
            width: 80,
        },
        {
            title: '游戏',
            key: 'game',
            width: 200,
            render: (_, record) => (
                <Space>
                    <Image
                        src={record.icon}
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
            title: '分类',
            dataIndex: 'category',
            key: 'category',
            width: 100,
            render: category => {
                const option = categoryOptions.find(o => o.value === category);
                return <Tag>{option?.label || category}</Tag>;
            },
        },
        {
            title: '平台',
            dataIndex: 'platform',
            key: 'platform',
            width: 150,
            render: platforms => (
                <Space size={4} wrap>
                    {platforms.map((p: string) => {
                        const option = platformOptions.find(o => o.value === p);
                        return <Tag key={p} color="blue">{option?.label || p}</Tag>;
                    })}
                </Space>
            ),
        },
        {
            title: '陪玩师',
            dataIndex: 'playerCount',
            key: 'playerCount',
            width: 100,
            render: count => <span>{count} 人</span>,
        },
        {
            title: '订单数',
            dataIndex: 'orderCount',
            key: 'orderCount',
            width: 100,
            render: count => <span>{count} 单</span>,
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: (status, record) => (
                <PermissionGuard permission={GAME_PERMISSIONS.UPDATE} fallback={
                    <Tag color={status === 'active' ? 'success' : 'default'}>
                        {status === 'active' ? '已上架' : '已下架'}
                    </Tag>
                }>
                    <Switch
                        checked={status === 'active'}
                        checkedChildren="上架"
                        unCheckedChildren="下架"
                        onChange={() => handleToggleStatus(record)}
                    />
                </PermissionGuard>
            ),
        },
        {
            title: '更新时间',
            dataIndex: 'updatedAt',
            key: 'updatedAt',
            width: 180,
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
                        <Popconfirm
                            title="确定要删除该游戏吗？"
                            onConfirm={() => handleDelete(record)}
                        >
                            <Button type="link" size="small" danger icon={<DeleteOutlined />}>
                                删除
                            </Button>
                        </Popconfirm>
                    </PermissionGuard>
                </Space>
            ),
        },
    ];

    return (
        <PageContainer title="游戏管理" subTitle="管理平台支持的游戏">
            <SearchTable
                columns={columns}
                dataSource={games}
                rowKey="id"
                searchFields={searchFields}
                onSearch={() => loadData()}
                onRefresh={loadData}
                loading={loading}
                showCreate={true}
                createText="新增游戏"
                createPermission={GAME_PERMISSIONS.CREATE}
                onCreate={() => {
                    setCurrentGame(null);
                    form.resetFields();
                    setFileList([]);
                    setEditModalVisible(true);
                }}
                pagination={{
                    current,
                    pageSize,
                    total,
                    showSizeChanger: true,
                    showTotal: total => `共 ${total} 条`,
                    onChange: (page, size) => {
                        setCurrent(page);
                        setPageSize(size);
                    },
                }}
                scroll={{ x: 1200 }}
            />

            {/* 编辑弹窗 */}
            <Modal
                title={currentGame ? '编辑游戏' : '新增游戏'}
                open={editModalVisible}
                onOk={handleSaveEdit}
                onCancel={() => setEditModalVisible(false)}
                width={600}
            >
                <Form form={form} layout="vertical">
                    <Form.Item
                        name="name"
                        label="游戏名称"
                        rules={[{ required: true, message: '请输入游戏名称' }]}
                    >
                        <Input placeholder="请输入游戏名称" />
                    </Form.Item>

                    <Form.Item name="icon" label="游戏图标">
                        <Upload
                            listType="picture-card"
                            fileList={fileList}
                            onChange={({ fileList }) => setFileList(fileList)}
                            maxCount={1}
                            beforeUpload={() => false}
                        >
                            {fileList.length === 0 && (
                                <div>
                                    <PlusOutlined />
                                    <div style={{ marginTop: 8 }}>上传</div>
                                </div>
                            )}
                        </Upload>
                    </Form.Item>

                    <Form.Item
                        name="category"
                        label="分类"
                        rules={[{ required: true, message: '请选择分类' }]}
                    >
                        <Select placeholder="请选择分类" options={categoryOptions} />
                    </Form.Item>

                    <Form.Item
                        name="platform"
                        label="支持平台"
                        rules={[{ required: true, message: '请选择支持平台' }]}
                    >
                        <Select
                            mode="multiple"
                            placeholder="请选择支持平台"
                            options={platformOptions}
                        />
                    </Form.Item>

                    <Form.Item name="sortOrder" label="排序" initialValue={99}>
                        <InputNumber min={1} max={999} style={{ width: '100%' }} />
                    </Form.Item>

                    <Form.Item
                        name="status"
                        label="状态"
                        initialValue="active"
                        rules={[{ required: true }]}
                    >
                        <Select
                            options={[
                                { label: '已上架', value: 'active' },
                                { label: '已下架', value: 'inactive' },
                            ]}
                        />
                    </Form.Item>
                </Form>
            </Modal>
        </PageContainer>
    );
};

export default GamePage;
