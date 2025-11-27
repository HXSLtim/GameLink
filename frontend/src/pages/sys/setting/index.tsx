import React from 'react';
import { Card, Form, Input, Switch, Button, Tabs, Select, ColorPicker, Divider, message, List } from 'antd';
import { SaveOutlined, SettingOutlined, SafetyOutlined, BellOutlined, BgColorsOutlined } from '@ant-design/icons';
import { motion } from 'framer-motion';

const Settings: React.FC = () => {
    const [form] = Form.useForm();

    const onFinish = (values: any) => {
        console.log('Success:', values);
        message.success('设置已保存');
    };

    const items = [
        {
            key: '1',
            label: (
                <span>
                    <SettingOutlined />
                    基本设置
                </span>
            ),
            children: (
                <Form
                    form={form}
                    layout="vertical"
                    initialValues={{
                        siteName: 'GameLink 游戏陪玩平台',
                        siteDescription: '专业的游戏陪玩与社交平台',
                        maintenanceMode: false,
                        language: 'zh_CN',
                    }}
                    onFinish={onFinish}
                >
                    <Form.Item label="平台名称" name="siteName" rules={[{ required: true, message: '请输入平台名称' }]}>
                        <Input placeholder="请输入平台名称" />
                    </Form.Item>
                    <Form.Item label="平台描述" name="siteDescription">
                        <Input.TextArea rows={4} placeholder="请输入平台描述" />
                    </Form.Item>
                    <Form.Item label="默认语言" name="language">
                        <Select>
                            <Select.Option value="zh_CN">简体中文</Select.Option>
                            <Select.Option value="en_US">English</Select.Option>
                        </Select>
                    </Form.Item>
                    <Form.Item label="维护模式" name="maintenanceMode" valuePropName="checked" help="开启后，除管理员外用户将无法访问">
                        <Switch checkedChildren="开启" unCheckedChildren="关闭" />
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" htmlType="submit" icon={<SaveOutlined />} style={{ backgroundColor: '#5865F2' }}>
                            保存更改
                        </Button>
                    </Form.Item>
                </Form>
            ),
        },
        {
            key: '2',
            label: (
                <span>
                    <BgColorsOutlined />
                    主题设置
                </span>
            ),
            children: (
                <Form layout="vertical" initialValues={{ primaryColor: '#5865F2', darkMode: true }}>
                    <Form.Item label="主色调" name="primaryColor">
                        <ColorPicker showText />
                    </Form.Item>
                    <Form.Item label="深色模式" name="darkMode" valuePropName="checked">
                        <Switch checkedChildren="开启" unCheckedChildren="关闭" />
                    </Form.Item>
                    <Divider />
                    <Form.Item label="侧边栏风格">
                        <Select defaultValue="dark">
                            <Select.Option value="dark">深色 (Discord 风格)</Select.Option>
                            <Select.Option value="light">浅色</Select.Option>
                        </Select>
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" icon={<SaveOutlined />} style={{ backgroundColor: '#5865F2' }}>
                            应用主题
                        </Button>
                    </Form.Item>
                </Form>
            ),
        },
        {
            key: '3',
            label: (
                <span>
                    <SafetyOutlined />
                    安全设置
                </span>
            ),
            children: (
                <Form layout="vertical">
                    <Form.Item label="注册限制">
                        <Switch checkedChildren="开放注册" unCheckedChildren="停止注册" defaultChecked />
                    </Form.Item>
                    <Form.Item label="密码强度要求">
                        <Select defaultValue="medium">
                            <Select.Option value="low">低 (仅长度限制)</Select.Option>
                            <Select.Option value="medium">中 (需包含字母和数字)</Select.Option>
                            <Select.Option value="high">高 (需包含特殊字符)</Select.Option>
                        </Select>
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" style={{ backgroundColor: '#5865F2' }}>
                            更新安全策略
                        </Button>
                    </Form.Item>
                </Form>
            ),
        },
        {
            key: '4',
            label: (
                <span>
                    <BellOutlined />
                    通知设置
                </span>
            ),
            children: (
                <List
                    dataSource={[
                        { title: '新订单通知', desc: '当有新订单创建时发送通知' },
                        { title: '用户注册通知', desc: '当有新用户注册时发送通知' },
                        { title: '系统异常报警', desc: '系统发生严重错误时发送邮件报警' },
                    ]}
                    renderItem={(item: { title: string; desc: string }) => (
                        <List.Item extra={<Switch defaultChecked />}>
                            <List.Item.Meta
                                title={<span style={{ color: '#fff' }}>{item.title}</span>}
                                description={<span style={{ color: 'rgba(255,255,255,0.45)' }}>{item.desc}</span>}
                            />
                        </List.Item>
                    )}
                />
            ),
        },
    ];

    return (
        <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.3 }}
        >
            <Card bordered={false} title="系统设置" bodyStyle={{ paddingTop: 0 }}>
                <Tabs defaultActiveKey="1" items={items} />
            </Card>
        </motion.div>
    );
};

export default Settings;
