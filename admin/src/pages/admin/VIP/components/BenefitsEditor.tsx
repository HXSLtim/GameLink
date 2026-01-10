/**
 * VIP权益编辑器组件
 * 用于快速编辑VIP等级的权益列表
 */
import React, { useState } from 'react';
import {
    Modal,
    Input,
    Button,
    Space,
    List,
    message,
    Typography,
    Divider,
    Tag,
} from 'antd';
import { PlusOutlined, DeleteOutlined, DragOutlined } from '@ant-design/icons';

const { TextArea } = Input;
const { Text } = Typography;

interface BenefitsEditorProps {
    visible: boolean;
    benefits: string; // JSON字符串数组
    onSave: (benefits: string) => void;
    onCancel: () => void;
}

/**
 * 解析 benefits 字符串为字符串数组
 */
const parseBenefits = (benefitsStr: string): string[] => {
    if (!benefitsStr) return [];
    try {
        const parsed = JSON.parse(benefitsStr);
        return Array.isArray(parsed) ? parsed : [];
    } catch {
        return benefitsStr.split('\n').filter(s => s.trim());
    }
};

/**
 * 将字符串数组转换为 benefits 字符串
 */
const stringifyBenefits = (benefits: string[]): string => {
    return JSON.stringify(benefits.filter(s => s.trim()));
};

const BenefitsEditor: React.FC<BenefitsEditorProps> = ({
    visible,
    benefits,
    onSave,
    onCancel,
}) => {
    // Initialize from props - component remounts when key changes (from parent)
    const [items, setItems] = useState<string[]>(() => parseBenefits(benefits));
    const [newItem, setNewItem] = useState('');
    const [viewMode, setViewMode] = useState<'list' | 'text'>('list');

    const handleAdd = () => {
        if (!newItem.trim()) {
            message.warning('请输入权益内容');
            return;
        }
        if (items.includes(newItem.trim())) {
            message.warning('该权益已存在');
            return;
        }
        setItems([...items, newItem.trim()]);
        setNewItem('');
    };

    const handleRemove = (index: number) => {
        setItems(items.filter((_, i) => i !== index));
    };

    const handleSave = () => {
        onSave(stringifyBenefits(items));
        onCancel();
    };

    const handleTextChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        const text = e.target.value;
        const lines = text.split('\n').filter(line => line.trim());
        setItems(lines);
    };

    return (
        <Modal
            title="编辑会员权益"
            open={visible}
            onOk={handleSave}
            onCancel={onCancel}
            width={600}
            okText="保存"
            cancelText="取消"
        >
            <div style={{ marginBottom: 16 }}>
                <Space>
                    <Button
                        type={viewMode === 'list' ? 'primary' : 'default'}
                        onClick={() => setViewMode('list')}
                    >
                        列表模式
                    </Button>
                    <Button
                        type={viewMode === 'text' ? 'primary' : 'default'}
                        onClick={() => setViewMode('text')}
                    >
                        文本模式
                    </Button>
                </Space>
            </div>

            {viewMode === 'list' ? (
                <>
                    <Space.Compact style={{ width: '100%', marginBottom: 16 }}>
                        <Input
                            placeholder="输入权益内容，如：专属客服、优先派单等"
                            value={newItem}
                            onChange={(e) => setNewItem(e.target.value)}
                            onPressEnter={handleAdd}
                            prefix={<DragOutlined />}
                        />
                        <Button
                            type="primary"
                            icon={<PlusOutlined />}
                            onClick={handleAdd}
                        >
                            添加
                        </Button>
                    </Space.Compact>

                    {items.length > 0 ? (
                        <List
                            bordered
                            dataSource={items}
                            renderItem={(item, index) => (
                                <List.Item
                                    actions={[
                                        <Button
                                            type="text"
                                            danger
                                            icon={<DeleteOutlined />}
                                            onClick={() => handleRemove(index)}
                                        >
                                            删除
                                        </Button>,
                                    ]}
                                >
                                    <Tag color="blue">#{index + 1}</Tag>
                                    <span style={{ marginLeft: 8 }}>{item}</span>
                                </List.Item>
                            )}
                        />
                    ) : (
                        <div
                            style={{
                                textAlign: 'center',
                                padding: '40px 0',
                                color: '#999',
                            }}
                        >
                            暂无权益，请添加
                        </div>
                    )}
                </>
            ) : (
                <>
                    <Text type="secondary" style={{ display: 'block', marginBottom: 8 }}>
                        每行一个权益，会自动分割
                    </Text>
                    <TextArea
                        rows={10}
                        value={items.join('\n')}
                        onChange={handleTextChange}
                        placeholder="每行输入一个权益内容"
                    />
                </>
            )}

            <Divider />

            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Text type="secondary">
                    共 {items.length} 项权益
                </Text>
            </div>
        </Modal>
    );
};

export default BenefitsEditor;
