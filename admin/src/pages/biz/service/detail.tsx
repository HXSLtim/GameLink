import React, { useEffect, useState } from 'react';
import { Card, Descriptions, Tag, Button, Space, Spin, Image } from 'antd';
import { ArrowLeftOutlined, EditOutlined } from '@ant-design/icons';
import { useNavigate, useParams } from 'react-router-dom';
import { adminApi } from '@/api/admin';
import type { ServiceItem } from '@/api/admin';
import { motion } from 'framer-motion';

const ServiceItemDetail: React.FC = () => {
    const navigate = useNavigate();
    const { id } = useParams<{ id: string }>();
    const [loading, setLoading] = useState(false);
    const [data, setData] = useState<ServiceItem | null>(null);

    useEffect(() => {
        let isMounted = true;
        const loadData = async () => {
            try {
                const res = await adminApi.getServiceItem(Number(id));
                // @ts-expect-error API response type mismatch
                if (isMounted) setData(res.data);
            } catch {
                // message.error('Failed to load data');
            } finally {
                if (isMounted) setLoading(false);
            }
        };
        setLoading(true);
        loadData();
        return () => { isMounted = false; };
    }, [id]);

    if (loading || !data) return <Spin />;

    return (
        <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.3 }}
        >
            <Card
                title={
                    <Space>
                        <Button icon={<ArrowLeftOutlined />} type="text" onClick={() => navigate(-1)} />
                        服务详情
                    </Space>
                }
                extra={
                    <Button type="primary" icon={<EditOutlined />} onClick={() => navigate(`/admin/biz/service/${id}/edit`)}>
                        编辑
                    </Button>
                }
                variant="borderless"
            >
                <Descriptions bordered column={{ xxl: 4, xl: 3, lg: 3, md: 3, sm: 2, xs: 1 }}>
                    <Descriptions.Item label="服务名称">{data.name}</Descriptions.Item>
                    <Descriptions.Item label="服务编码">{data.itemCode}</Descriptions.Item>
                    <Descriptions.Item label="分类">
                        <Tag color="blue">{data.subCategory === 'solo' ? '单人护航' : data.subCategory === 'team' ? '团队护航' : '礼物'}</Tag>
                    </Descriptions.Item>
                    <Descriptions.Item label="价格">
                        <span style={{ color: '#faa61a', fontWeight: 'bold' }}>¥{(data.basePriceCents / 100).toFixed(2)}</span>
                    </Descriptions.Item>
                    <Descriptions.Item label="服务时长">{data.serviceHours} 小时</Descriptions.Item>
                    <Descriptions.Item label="状态">
                        <Tag color={data.isActive ? 'success' : 'default'}>
                            {data.isActive ? '已启用' : '已禁用'}
                        </Tag>
                    </Descriptions.Item>
                    <Descriptions.Item label="标签">
                        {data.tags ? (typeof data.tags === 'string' ? JSON.parse(data.tags) : data.tags).map((tag: string) => <Tag key={tag}>{tag}</Tag>) : '-'}
                    </Descriptions.Item>
                    <Descriptions.Item label="排序权重">{data.sortOrder}</Descriptions.Item>
                    <Descriptions.Item label="创建时间">{data.createdAt}</Descriptions.Item>
                    <Descriptions.Item label="更新时间">{data.updatedAt}</Descriptions.Item>
                    <Descriptions.Item label="描述" span={3}>
                        {data.description}
                    </Descriptions.Item>
                    <Descriptions.Item label="图标" span={3}>
                        {data.iconUrl ? <Image width={100} src={data.iconUrl} /> : '-'}
                    </Descriptions.Item>
                </Descriptions>
            </Card>
        </motion.div>
    );
};

export default ServiceItemDetail;
