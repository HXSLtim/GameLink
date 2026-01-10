import React, { useEffect, useState } from 'react';
import {
    Form, Input, InputNumber, Select, Button, Card,
    Radio, Upload, message, Space, Spin
} from 'antd';
import {
    ArrowLeftOutlined, SaveOutlined,
    PlusOutlined
} from '@ant-design/icons';
import { useNavigate, useParams } from 'react-router-dom';
import { adminApi } from '@/api/admin';
import { motion } from 'framer-motion';

const { Option } = Select;
const { TextArea } = Input;

const ServiceItemForm: React.FC = () => {
    const navigate = useNavigate();
    const { id } = useParams<{ id: string }>();
    const isEdit = !!id;
    const [form] = Form.useForm();
    const [loading, setLoading] = useState(false);
    const [submitting, setSubmitting] = useState(false);

    useEffect(() => {
        if (isEdit) {
            let isMounted = true;
            const loadData = async () => {
                try {
                    const res = await adminApi.getServiceItem(Number(id));
                    // eslint-disable-next-line @typescript-eslint/no-explicit-any
                    if (isMounted) form.setFieldsValue((res as any).data);
                } catch {
                    message.error('Failed to load data');
                } finally {
                    if (isMounted) setLoading(false);
                }
            };
            setLoading(true);
            loadData();
            return () => { isMounted = false; };
        }
    }, [isEdit, form, id]);

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const onFinish = async (values: any) => {
        setSubmitting(true);
        try {
            if (isEdit) {
                await adminApi.updateServiceItem(Number(id), values);
                message.success('Updated successfully');
            } else {
                await adminApi.createServiceItem(values);
                message.success('Created successfully');
            }
            navigate('/admin/biz/service');
        } catch {
            message.error('Operation failed');
        } finally {
            setSubmitting(false);
        }
    };

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
                        {isEdit ? '编辑服务项目' : '新建服务项目'}
                    </Space>
                }
                variant="borderless"
            >
                <Spin spinning={loading}>
                    <Form
                        form={form}
                        layout="vertical"
                        onFinish={onFinish}
                        initialValues={{
                            duration: 60,
                            sortOrder: 0,
                            status: 'active',
                        }}
                        style={{ maxWidth: 800, margin: '0 auto' }}
                    >
                        <Form.Item
                            label="服务名称"
                            name="name"
                            rules={[
                                { required: true, message: '请输入服务名称' },
                                { min: 2, max: 50, message: '长度在 2 到 50 个字符' }
                            ]}
                        >
                            <Input placeholder="请输入服务名称" />
                        </Form.Item>

                        <Form.Item
                            label="关联游戏"
                            name="gameId"
                            rules={[{ required: true, message: '请选择关联游戏' }]}
                        >
                            <Select placeholder="请选择游戏">
                                <Option value={1}>王者荣耀</Option>
                                <Option value={2}>英雄联盟</Option>
                            </Select>
                        </Form.Item>

                        <Form.Item
                            label="服务分类"
                            name="category"
                            rules={[{ required: true, message: '请选择服务分类' }]}
                        >
                            <Select placeholder="请选择分类">
                                <Option value="rank">上分 (Rank)</Option>
                                <Option value="rush">陪玩 (Rush)</Option>
                                <Option value="teach">教学 (Teach)</Option>
                                <Option value="entertain">娱乐 (Entertain)</Option>
                            </Select>
                        </Form.Item>

                        <Space size="large" style={{ display: 'flex' }}>
                            <Form.Item
                                label="基础价格 (元/小时)"
                                name="price"
                                rules={[{ required: true, message: '请输入价格' }]}
                            >
                                <InputNumber min={0.01} precision={2} style={{ width: 200 }} prefix="¥" />
                            </Form.Item>

                            <Form.Item
                                label="服务时长 (分钟)"
                                name="duration"
                                rules={[
                                    { required: true, message: '请输入时长' },
                                    { type: 'number', min: 30, max: 480, message: '时长在 30-480 分钟之间' }
                                ]}
                            >
                                <InputNumber step={30} style={{ width: 200 }} suffix="分钟" />
                            </Form.Item>
                        </Space>

                        <Form.Item
                            label="服务描述"
                            name="description"
                            rules={[
                                { required: true, message: '请输入服务描述' },
                                { min: 10, max: 500, message: '长度在 10 到 500 个字符' }
                            ]}
                        >
                            <TextArea rows={4} placeholder="请输入服务描述" showCount maxLength={500} />
                        </Form.Item>

                        <Form.Item label="服务标签" name="tags">
                            <Select mode="tags" placeholder="输入标签后回车" maxTagCount={5} />
                        </Form.Item>

                        <Form.Item label="服务图标" name="icon">
                            <Upload listType="picture-card">
                                <div>
                                    <PlusOutlined />
                                    <div style={{ marginTop: 8 }}>上传</div>
                                </div>
                            </Upload>
                        </Form.Item>

                        <Space size="large">
                            <Form.Item label="排序权重" name="sortOrder">
                                <InputNumber min={0} style={{ width: 150 }} />
                            </Form.Item>

                            <Form.Item label="状态" name="status">
                                <Radio.Group>
                                    <Radio value="active">启用</Radio>
                                    <Radio value="inactive">禁用</Radio>
                                </Radio.Group>
                            </Form.Item>
                        </Space>

                        <Form.Item style={{ marginTop: 24 }}>
                            <Space>
                                <Button type="primary" htmlType="submit" loading={submitting} icon={<SaveOutlined />}>
                                    保存提交
                                </Button>
                                <Button onClick={() => navigate(-1)}>取消</Button>
                            </Space>
                        </Form.Item>
                    </Form>
                </Spin>
            </Card>
        </motion.div>
    );
};

export default ServiceItemForm;
