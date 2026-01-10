import React, { useEffect, useState, useCallback } from 'react';
import { Card, Descriptions, Tag, Button, Space, Spin, Image, Modal, Form, Input, InputNumber, Select, Switch } from 'antd';
import { ArrowLeftOutlined, EditOutlined } from '@ant-design/icons';
import { useNavigate, useParams } from 'react-router-dom';
import { adminApi } from '@/api/admin';
import type { ServiceItem } from '@/api/admin';
import { motion } from 'framer-motion';
import { logger } from '@/utils/logger';

const ServiceItemDetail: React.FC = () => {
    const navigate = useNavigate();
    const { id } = useParams<{ id: string }>();
    const [loading, setLoading] = useState(false);
    const [data, setData] = useState<ServiceItem | null>(null);
    const [editModalVisible, setEditModalVisible] = useState(false);
    const [form] = Form.useForm();

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

    const handleOpenEdit = useCallback(() => {
        if (!data) return;
        form.setFieldsValue({
            name: data.name,
            itemCode: data.itemCode,
            subCategory: data.subCategory,
            basePriceCents: data.basePriceCents ? data.basePriceCents / 100 : undefined,
            commissionRate: data.commissionRate,
            sortOrder: data.sortOrder,
            iconUrl: data.iconUrl,
            description: data.description,
            isActive: data.isActive ?? true,
        });
        setEditModalVisible(true);
    }, [data, form]);

    const handleSaveEdit = useCallback(async () => {
        if (!data) return;
        try {
            const values = await form.validateFields();
            const payload = {
                ...values,
                basePriceCents: values.basePriceCents ? Math.round(values.basePriceCents * 100) : 0,
            };

            const res = await adminApi.updateServiceItem(data.id, payload);
            if (res.data.success) {
                // Refresh data
                const updatedRes = await adminApi.getServiceItem(data.id);
                // @ts-expect-error API response type mismatch
                setData(updatedRes.data);
                setEditModalVisible(false);
                Modal.success({
                    title: '成功',
                    content: '更新成功',
                });
            } else {
                Modal.error({
                    title: '失败',
                    content: res.data.message || '更新失败',
                });
            }
        } catch (error) {
            logger.error('Save failed:', error);
            Modal.error({
                title: '失败',
                content: '保存失败',
            });
        }
    }, [data, form]);

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
                    <Button type="primary" icon={<EditOutlined />} onClick={handleOpenEdit}>
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

            {/* 编辑弹窗 */}
            <Modal
                title="编辑服务"
                open={editModalVisible}
                onOk={handleSaveEdit}
                onCancel={() => setEditModalVisible(false)}
                width={600}
                okText="保存"
                cancelText="取消"
            >
                <Form form={form} layout="vertical">
                    <Form.Item
                        name="name"
                        label="服务名称"
                        rules={[{ required: true, message: '请输入服务名称' }]}
                    >
                        <Input placeholder="请输入服务名称" />
                    </Form.Item>
                    <Form.Item
                        name="itemCode"
                        label="服务编码"
                        rules={[{ required: true, message: '请输入服务编码' }]}
                    >
                        <Input placeholder="请输入服务编码" />
                    </Form.Item>
                    <Form.Item
                        name="subCategory"
                        label="服务分类"
                        rules={[{ required: true, message: '请选择服务分类' }]}
                    >
                        <Select placeholder="请选择服务分类">
                            <Select.Option value="solo">单人护航</Select.Option>
                            <Select.Option value="team">团队护航</Select.Option>
                            <Select.Option value="gift">礼物</Select.Option>
                        </Select>
                    </Form.Item>
                    <Form.Item
                        name="basePriceCents"
                        label="基础价格（元）"
                        rules={[{ required: true, message: '请输入基础价格' }]}
                    >
                        <InputNumber
                            min={0}
                            precision={2}
                            placeholder="请输入基础价格"
                            style={{ width: '100%' }}
                            prefix="¥"
                        />
                    </Form.Item>
                    <Form.Item
                        name="commissionRate"
                        label="佣金比例（%）"
                        rules={[{ required: true, message: '请输入佣金比例' }]}
                    >
                        <InputNumber
                            min={0}
                            max={100}
                            precision={2}
                            placeholder="请输入佣金比例"
                            style={{ width: '100%' }}
                        />
                    </Form.Item>
                    <Form.Item
                        name="sortOrder"
                        label="排序"
                    >
                        <InputNumber
                            min={0}
                            placeholder="请输入排序值"
                            style={{ width: '100%' }}
                        />
                    </Form.Item>
                    <Form.Item
                        name="iconUrl"
                        label="图标URL"
                    >
                        <Input placeholder="请输入图标URL" />
                    </Form.Item>
                    <Form.Item
                        name="description"
                        label="描述"
                    >
                        <Input.TextArea rows={3} placeholder="请输入服务描述" />
                    </Form.Item>
                    <Form.Item
                        name="isActive"
                        label="状态"
                        valuePropName="checked"
                    >
                        <Switch checkedChildren="启用" unCheckedChildren="禁用" />
                    </Form.Item>
                </Form>
            </Modal>
        </motion.div>
    );
};

export default ServiceItemDetail;
