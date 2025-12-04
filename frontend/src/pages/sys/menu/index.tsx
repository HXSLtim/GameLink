import React, { useEffect, useState } from 'react';
import { Table, Button, Space, message, Tag, Popconfirm, Card } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, ReloadOutlined } from '@ant-design/icons';
import { adminApi } from '@/api/admin';
import type { Menu } from '@/api/admin';
import { getIcon } from '@/utils/iconMap';

const MenuList: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [data, setData] = useState<Menu[]>([]);
    const [pagination, setPagination] = useState({
        current: 1,
        pageSize: 10,
        total: 0,
    });

    const fetchData = async (page = pagination.current, pageSize = pagination.pageSize) => {
        setLoading(true);
        try {
            const res = await adminApi.getMenus({ page, page_size: pageSize }) as any;
            setData(res.data);
            if (res.pagination) {
                setPagination({
                    current: res.pagination.page,
                    pageSize: res.pagination.page_size,
                    total: res.pagination.total,
                });
            }
        } catch (error) {
            message.error('获取菜单列表失败');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, []);

    const handleDelete = async (id: number) => {
        try {
            await adminApi.deleteMenu(id);
            message.success('删除成功');
            fetchData();
        } catch (error) {
            message.error('删除失败');
        }
    };

    const handleTableChange = (newPagination: any) => {
        fetchData(newPagination.current, newPagination.pageSize);
    };

    const columns = [
        {
            title: '菜单名称',
            dataIndex: 'name',
            key: 'name',
        },
        {
            title: '图标',
            dataIndex: 'icon',
            key: 'icon',
            render: (icon: string) => icon ? getIcon(icon) : '-',
        },
        {
            title: '路由路径',
            dataIndex: 'path',
            key: 'path',
        },
        {
            title: '组件路径',
            dataIndex: 'component',
            key: 'component',
        },
        {
            title: '类型',
            dataIndex: 'type',
            key: 'type',
            render: (type: string) => {
                const safeType = type || 'unknown';
                let color = 'default';
                if (safeType === 'menu') color = 'blue';
                if (safeType === 'page') color = 'green';
                if (safeType === 'button') color = 'orange';
                const typeMap: Record<string, string> = { menu: '目录', page: '页面', button: '按钮' };
                return <Tag color={color}>{typeMap[safeType] || safeType}</Tag>;
            },
        },
        {
            title: '权限标识',
            dataIndex: 'permission',
            key: 'permission',
            render: (text: string) => text ? <Tag color="cyan">{text}</Tag> : '-',
        },
        {
            title: '排序',
            dataIndex: 'order',
            key: 'order',
        },
        {
            title: '状态',
            dataIndex: 'visible',
            key: 'visible',
            render: (visible: boolean) => (
                <Tag color={visible ? 'success' : 'default'}>
                    {visible ? '显示' : '隐藏'}
                </Tag>
            ),
        },
        {
            title: '操作',
            key: 'action',
            render: (_: any, record: Menu) => (
                <Space size="small">
                    <Button type="text" icon={<EditOutlined />} onClick={() => message.info('编辑功能暂未实现')} />
                    <Popconfirm title="确定要删除吗？" onConfirm={() => handleDelete(record.id)}>
                        <Button type="text" danger icon={<DeleteOutlined />} />
                    </Popconfirm>
                </Space>
            ),
        },
    ];

    return (
        <Card style={{ border: 'none' }}>
            <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
                <h2>菜单管理</h2>
                <Space>
                    <Button icon={<ReloadOutlined />} onClick={() => fetchData()}>刷新</Button>
                    <Button type="primary" icon={<PlusOutlined />} onClick={() => message.info('新增功能暂未实现')}>
                        新增菜单
                    </Button>
                </Space>
            </div>
            <Table
                columns={columns}
                dataSource={data}
                rowKey="id"
                loading={loading}
                pagination={pagination}
                onChange={handleTableChange}
                expandable={{
                    childrenColumnName: 'children', // Assuming backend returns nested 'children'
                }}
            />
        </Card>
    );
};

export default MenuList;
