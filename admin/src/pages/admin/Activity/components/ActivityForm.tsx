/**
 * Activity Form Component
 * Create/Edit activity form with validation
 */
import React, { useEffect } from 'react';
import {
    Modal,
    Form,
    Input,
    Select,
    InputNumber,
    Switch,
    DatePicker,
    Space,
    Divider,
    Typography,
    Row,
    Col,
    message,
} from 'antd';
import dayjs from 'dayjs';
import {
    type ActivityType,
    type ActivityStatus,
    type CreateActivityDto,
    type Activity,
    getActivityTypeLabel,
} from '@/api/activity';

const { TextArea } = Input;
const { Text } = Typography;

interface ActivityFormProps {
    visible: boolean;
    editing: boolean;
    initialValues?: Activity | null;
    loading: boolean;
    onCancel: () => void;
    onSubmit: (values: CreateActivityDto) => Promise<void>;
}

const ActivityForm: React.FC<ActivityFormProps> = ({
    visible,
    editing,
    initialValues,
    loading,
    onCancel,
    onSubmit,
}) => {
    const [form] = Form.useForm<CreateActivityDto>();
    const activityType = Form.useWatch('type', form);

    // Set initial form values
    useEffect(() => {
        if (visible && !editing) {
            form.resetFields();
            form.setFieldsValue({
                type: 'coupon' as ActivityType,
                status: 'draft' as ActivityStatus,
                totalLimit: 10000,
                dailyLimit: 1000,
                perUserLimit: 1,
                allowVipStack: false,
                sortOrder: 0,
                isVisible: false,
            });
        } else if (visible && editing && initialValues) {
            form.setFieldsValue({
                name: initialValues.name,
                description: initialValues.description,
                type: initialValues.type,
                status: initialValues.status,
                coverUrl: initialValues.coverUrl,
                bannerUrl: initialValues.bannerUrl,
                preheatAt: initialValues.preheatAt ? dayjs(initialValues.preheatAt) : undefined,
                startAt: dayjs(initialValues.startAt),
                endAt: dayjs(initialValues.endAt),
                totalLimit: initialValues.totalLimit,
                dailyLimit: initialValues.dailyLimit,
                perUserLimit: initialValues.perUserLimit,
                allowVipStack: initialValues.allowVipStack,
                rules: initialValues.rules,
                sortOrder: initialValues.sortOrder,
                isVisible: initialValues.isVisible,
            });
        }
    }, [visible, editing, initialValues, form]);

    const handleOk = async () => {
        try {
            const values = await form.validateFields();
            const data: CreateActivityDto = {
                ...values,
                preheatAt: values.preheatAt
                    ? (values.preheatAt as dayjs.Dayjs).format('YYYY-MM-DD HH:mm:ss')
                    : undefined,
                startAt: (values.startAt as dayjs.Dayjs).format('YYYY-MM-DD HH:mm:ss'),
                endAt: (values.endAt as dayjs.Dayjs).format('YYYY-MM-DD HH:mm:ss'),
            };
            await onSubmit(data);
        } catch (error) {
            console.error('Form validation failed:', error);
        }
    };

    return (
        <Modal
            title={editing ? '编辑活动' : '新建活动'}
            open={visible}
            onOk={handleOk}
            onCancel={onCancel}
            confirmLoading={loading}
            width={700}
            okText="保存"
            cancelText="取消"
        >
            <Form
                form={form}
                layout="vertical"
                autoComplete="off"
            >
                {/* Basic Information */}
                <Divider orientation="left">基本信息</Divider>

                <Row gutter={16}>
                    <Col span={12}>
                        <Form.Item
                            name="name"
                            label="活动名称"
                            rules={[{ required: true, message: '请输入活动名称' }]}
                        >
                            <Input placeholder="如：新人狂欢活动" maxLength={100} />
                        </Form.Item>
                    </Col>
                    <Col span={12}>
                        <Form.Item
                            name="type"
                            label="活动类型"
                            rules={[{ required: true, message: '请选择活动类型' }]}
                        >
                            <Select
                                placeholder="请选择类型"
                                options={[
                                    { label: getActivityTypeLabel('coupon'), value: 'coupon' },
                                    { label: getActivityTypeLabel('discount'), value: 'discount' },
                                    { label: getActivityTypeLabel('gift'), value: 'gift' },
                                ]}
                            />
                        </Form.Item>
                    </Col>
                </Row>

                <Form.Item
                    name="description"
                    label="活动描述"
                >
                    <TextArea
                        rows={3}
                        placeholder="请输入活动描述"
                        maxLength={500}
                    />
                </Form.Item>

                {/* Time Configuration */}
                <Divider orientation="left">时间配置</Divider>

                <Row gutter={16}>
                    <Col span={12}>
                        <Form.Item
                            name="preheatAt"
                            label="预热时间"
                            extra="活动开始前展示预告"
                        >
                            <DatePicker
                                showTime
                                style={{ width: '100%' }}
                                placeholder="选择预热时间"
                                disabledDate={(current) => current && current < dayjs().startOf('day')}
                            />
                        </Form.Item>
                    </Col>
                    <Col span={12}>
                        <Form.Item
                            name="startAt"
                            label="开始时间"
                            rules={[{ required: true, message: '请选择开始时间' }]}
                        >
                            <DatePicker
                                showTime
                                style={{ width: '100%' }}
                                placeholder="选择开始时间"
                                disabledDate={(current) => current && current < dayjs().startOf('day')}
                            />
                        </Form.Item>
                    </Col>
                </Row>

                <Form.Item
                    name="endAt"
                    label="结束时间"
                    rules={[{ required: true, message: '请选择结束时间' }]}
                >
                    <DatePicker
                        showTime
                        style={{ width: '100%' }}
                        placeholder="选择结束时间"
                        disabledDate={(current) => current && current < dayjs().startOf('day')}
                    />
                </Form.Item>

                {/* Participation Limits */}
                <Divider orientation="left">参与限制</Divider>

                <Row gutter={16}>
                    <Col span={8}>
                        <Form.Item
                            name="totalLimit"
                            label="总参与人数"
                            rules={[{ required: true, message: '请输入总参与人数限制' }]}
                            extra="0表示不限制"
                        >
                            <InputNumber
                                min={0}
                                max={10000000}
                                style={{ width: '100%' }}
                                placeholder="如：10000"
                            />
                        </Form.Item>
                    </Col>
                    <Col span={8}>
                        <Form.Item
                            name="dailyLimit"
                            label="每日参与人数"
                            rules={[{ required: true, message: '请输入每日参与人数限制' }]}
                            extra="0表示不限制"
                        >
                            <InputNumber
                                min={0}
                                max={1000000}
                                style={{ width: '100%' }}
                                placeholder="如：1000"
                            />
                        </Form.Item>
                    </Col>
                    <Col span={8}>
                        <Form.Item
                            name="perUserLimit"
                            label="每人参与次数"
                            rules={[{ required: true, message: '请输入每人参与次数限制' }]}
                            extra="0表示不限制"
                        >
                            <InputNumber
                                min={0}
                                max={100}
                                style={{ width: '100%' }}
                                placeholder="如：1"
                            />
                        </Form.Item>
                    </Col>
                </Row>

                {/* VIP Settings */}
                <Divider orientation="left">VIP设置</Divider>

                <Form.Item
                    name="allowVipStack"
                    label="允许VIP叠加"
                    valuePropName="checked"
                    extra="VIP用户是否可以叠加活动优惠"
                >
                    <Switch
                        checkedChildren="允许"
                        unCheckedChildren="不允许"
                    />
                </Form.Item>

                {/* Display Settings */}
                <Divider orientation="left">显示设置</Divider>

                <Row gutter={16}>
                    <Col span={12}>
                        <Form.Item
                            name="isVisible"
                            label="是否显示"
                            valuePropName="checked"
                            extra="活动是否在前端显示"
                        >
                            <Switch
                                checkedChildren="显示"
                                unCheckedChildren="隐藏"
                            />
                        </Form.Item>
                    </Col>
                    <Col span={12}>
                        <Form.Item
                            name="sortOrder"
                            label="排序顺序"
                            extra="数值越小越靠前"
                        >
                            <InputNumber
                                min={0}
                                max={10000}
                                style={{ width: '100%' }}
                                placeholder="如：0"
                            />
                        </Form.Item>
                    </Col>
                </Row>

                {/* Media URLs */}
                <Divider orientation="left">媒体资源</Divider>

                <Row gutter={16}>
                    <Col span={12}>
                        <Form.Item
                            name="coverUrl"
                            label="封面图片URL"
                        >
                            <Input placeholder="https://example.com/cover.jpg" />
                        </Form.Item>
                    </Col>
                    <Col span={12}>
                        <Form.Item
                            name="bannerUrl"
                            label="横幅图片URL"
                        >
                            <Input placeholder="https://example.com/banner.jpg" />
                        </Form.Item>
                    </Col>
                </Row>

                {/* Rules */}
                <Divider orientation="left">活动规则</Divider>

                <Form.Item
                    name="rules"
                    label="活动规则"
                >
                    <TextArea
                        rows={4}
                        placeholder="请输入活动规则说明"
                        maxLength={2000}
                    />
                </Form.Item>

                {/* Status */}
                <Divider orientation="left">状态</Divider>

                <Form.Item
                    name="status"
                    label="活动状态"
                    rules={[{ required: true, message: '请选择活动状态' }]}
                >
                    <Select
                        placeholder="请选择状态"
                        options={[
                            { label: '草稿', value: 'draft' },
                            { label: '预热中', value: 'preheat' },
                            { label: '进行中', value: 'active' },
                            { label: '已暂停', value: 'paused' },
                            { label: '已结束', value: 'ended' },
                            { label: '已取消', value: 'canceled' },
                        ]}
                    />
                </Form.Item>
            </Form>
        </Modal>
    );
};

export default ActivityForm;
