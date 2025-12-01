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
            message.error('Failed to fetch menus');
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
            message.success('Deleted successfully');
            fetchData();
        } catch (error) {
            message.error('Delete failed');
        }
    };

    const handleTableChange = (newPagination: any) => {
        fetchData(newPagination.current, newPagination.pageSize);
    };

    const columns = [
        {
            title: 'Name',
            dataIndex: 'name',
            key: 'name',
        },
        {
            title: 'Icon',
            dataIndex: 'icon',
            key: 'icon',
            render: (icon: string) => icon ? getIcon(icon) : '-',
        },
        {
            title: 'Path',
            dataIndex: 'path',
            key: 'path',
        },
        {
            title: 'Component',
            dataIndex: 'component',
            key: 'component',
        },
        {
            title: 'Type',
            dataIndex: 'type',
            key: 'type',
            render: (type: string) => {
                const safeType = type || 'unknown';
                let color = 'default';
                if (safeType === 'menu') color = 'blue';
                if (safeType === 'page') color = 'green';
                if (safeType === 'button') color = 'orange';
                return <Tag color={color}>{safeType.toUpperCase()}</Tag>;
            },
        },
        {
            title: 'Permission',
            dataIndex: 'permission',
            key: 'permission',
            render: (text: string) => text ? <Tag color="cyan">{text}</Tag> : '-',
        },
        {
            title: 'Sort',
            dataIndex: 'order',
            key: 'order',
        },
        {
            title: 'Status',
            dataIndex: 'visible',
            key: 'visible',
            render: (visible: boolean) => (
                <Tag color={visible ? 'success' : 'default'}>
                    {visible ? 'Visible' : 'Hidden'}
                </Tag>
            ),
        },
        {
            title: 'Actions',
            key: 'action',
            render: (_: any, record: Menu) => (
                <Space size="small">
                    <Button type="text" icon={<EditOutlined />} onClick={() => message.info('Edit not implemented yet')} />
                    <Popconfirm title="Are you sure?" onConfirm={() => handleDelete(record.id)}>
                        <Button type="text" danger icon={<DeleteOutlined />} />
                    </Popconfirm>
                </Space>
            ),
        },
    ];

    return (
        <Card style={{ border: 'none' }}>
            <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
                <h2>Menu Management</h2>
                <Space>
                    <Button icon={<ReloadOutlined />} onClick={() => fetchData()}>Refresh</Button>
                    <Button type="primary" icon={<PlusOutlined />} onClick={() => message.info('Create not implemented yet')}>
                        Add Menu
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
