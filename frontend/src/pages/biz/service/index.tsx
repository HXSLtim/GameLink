import React, { useState, useEffect } from 'react';
import {
    Table, Button, Input, Select, Space, Tag, message,
    Popconfirm, Card, Row, Col, Tooltip
} from 'antd';
import {
    PlusOutlined, SearchOutlined, EditOutlined, DeleteOutlined,
    CheckCircleOutlined, StopOutlined, ReloadOutlined
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { adminApi } from '@/api/admin';
import type { ServiceItem } from '@/api/admin';
import { motion } from 'framer-motion';

const { Option } = Select;

const ServiceItemList: React.FC = () => {
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);
    const [data, setData] = useState<ServiceItem[]>([]);
    const [pagination, setPagination] = useState({ current: 1, pageSize: 20, total: 0 });
    const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

    // Filter States
    const [filters, setFilters] = useState({
        keyword: '',
        gameId: undefined as number | undefined,
        category: undefined as string | undefined,
        status: undefined as string | undefined,
    });

    // Data Loading
    const fetchData = async (page = 1) => {
        setLoading(true);
        try {
            const res = await adminApi.getServiceItems({
                ...filters,
                page,
                page_size: pagination.pageSize,
                game_id: filters.gameId,
            });
            // @ts-ignore
            const { data: items, meta } = res.data;
            setData(items);
            setPagination({ ...pagination, current: page, total: meta.total });
            setLoading(false);
        } catch (error) {
            console.error(error);
            message.error('Failed to fetch data');
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, [filters]);

    const handleTableChange = (newPagination: any) => {
        fetchData(newPagination.current);
    };

    const handleDelete = async (id: number) => {
        try {
            await adminApi.deleteServiceItem(id);
            message.success('Deleted successfully');
            fetchData(pagination.current);
        } catch (error) {
            message.error('Delete failed');
        }
    };

    const handleBatchStatus = async (status: 'active' | 'inactive') => {
        if (selectedRowKeys.length === 0) return;
        try {
            await adminApi.batchUpdateServiceItemStatus(selectedRowKeys as number[], status);
            message.success(`Batch ${status === 'active' ? 'enabled' : 'disabled'} successfully`);
            setSelectedRowKeys([]);
            fetchData(pagination.current);
        } catch (error) {
            message.error('Operation failed');
        }
    };

    const columns = [
        {
            title: 'ID',
            dataIndex: 'id',
            width: 80,
        },
        {
            title: '服务名称',
            dataIndex: 'name',
            render: (text: string, record: ServiceItem) => (
                <Space>
                    {record.icon && <img src={record.icon} alt="" style={{ width: 24, height: 24 }} />}
                    <span style={{ fontWeight: 500 }}>{text}</span>
                </Space>
            ),
        },
        {
            title: '所属游戏',
            dataIndex: 'gameName',
            width: 120,
        },
        {
            title: '分类',
            dataIndex: 'category',
            width: 100,
            render: (cat: string) => {
                const map: Record<string, string> = { rank: '上分', rush: '陪玩', teach: '教学', entertain: '娱乐' };
                return <Tag color="blue">{map[cat] || cat}</Tag>;
            }
        },
        {
            title: '价格',
            dataIndex: 'price',
            width: 100,
            render: (price: number) => <span style={{ color: '#faa61a', fontWeight: 'bold' }}>¥{price.toFixed(2)}</span>,
        },
        {
            title: '状态',
            dataIndex: 'status',
            width: 100,
            render: (status: string) => (
                <Tag color={status === 'active' ? 'success' : 'default'}>
                    {status === 'active' ? '已启用' : '已禁用'}
                </Tag>
            ),
        },
        {
            title: '排序',
            dataIndex: 'sortOrder',
            width: 80,
            sorter: true,
        },
        {
            title: '操作',
            key: 'action',
            width: 200,
            render: (_: any, record: ServiceItem) => (
                <Space size="small">
                    <Tooltip title="编辑">
                        <Button
                            type="text"
                            icon={<EditOutlined />}
                            onClick={() => navigate(`/admin/service-items/${record.id}/edit`)}
                        />
                    </Tooltip>
                    <Tooltip title={record.status === 'active' ? '禁用' : '启用'}>
                        <Button
                            type="text"
                            icon={record.status === 'active' ? <StopOutlined /> : <CheckCircleOutlined />}
                            style={{ color: record.status === 'active' ? '#ff4d4f' : '#52c41a' }}
                            onClick={() => {
                                // Toggle status logic
                                message.info('切换状态');
                            }}
                        />
                    </Tooltip>
                    <Popconfirm title="确定删除吗？" onConfirm={() => handleDelete(record.id)}>
                        <Tooltip title="删除">
                            <Button type="text" danger icon={<DeleteOutlined />} />
                        </Tooltip>
                    </Popconfirm>
                </Space>
            ),
        },
    ];

    return (
        <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.3 }}
        >
            <Card bordered={false}>
                <div style={{ marginBottom: 16 }}>
                    <Row gutter={[16, 16]} justify="space-between">
                        <Col>
                            <Space>
                                <Input
                                    placeholder="搜索服务名称"
                                    prefix={<SearchOutlined />}
                                    style={{ width: 200 }}
                                    value={filters.keyword}
                                    onChange={e => setFilters({ ...filters, keyword: e.target.value })}
                                />
                                <Select
                                    placeholder="选择游戏"
                                    style={{ width: 150 }}
                                    allowClear
                                    onChange={val => setFilters({ ...filters, gameId: val })}
                                >
                                    <Option value={1}>王者荣耀</Option>
                                    <Option value={2}>英雄联盟</Option>
                                </Select>
                                <Select
                                    placeholder="服务分类"
                                    style={{ width: 120 }}
                                    allowClear
                                    onChange={val => setFilters({ ...filters, category: val })}
                                >
                                    <Option value="rank">上分</Option>
                                    <Option value="rush">陪玩</Option>
                                    <Option value="teach">教学</Option>
                                    <Option value="entertain">娱乐</Option>
                                </Select>
                                <Select
                                    placeholder="状态"
                                    style={{ width: 100 }}
                                    allowClear
                                    onChange={val => setFilters({ ...filters, status: val })}
                                >
                                    <Option value="active">启用</Option>
                                    <Option value="inactive">禁用</Option>
                                </Select>
                                <Button icon={<SearchOutlined />} type="primary" onClick={() => fetchData(1)}>查询</Button>
                                <Button icon={<ReloadOutlined />} onClick={() => setFilters({ keyword: '', gameId: undefined, category: undefined, status: undefined })}>重置</Button>
                            </Space>
                        </Col>
                        <Col>
                            <Space>
                                <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/admin/service-items/create')}>
                                    新建服务
                                </Button>
                                {selectedRowKeys.length > 0 && (
                                    <>
                                        <Button onClick={() => handleBatchStatus('active')}>批量启用</Button>
                                        <Button danger onClick={() => handleBatchStatus('inactive')}>批量禁用</Button>
                                    </>
                                )}
                            </Space>
                        </Col>
                    </Row>
                </div>

                <Table
                    rowSelection={{
                        selectedRowKeys,
                        onChange: setSelectedRowKeys,
                    }}
                    columns={columns}
                    dataSource={data}
                    rowKey="id"
                    pagination={pagination}
                    loading={loading}
                    onChange={handleTableChange}
                />
            </Card>
        </motion.div>
    );
};

export default ServiceItemList;
