import React, { useState } from 'react';
import { Form, Row, Col, Input, Select, DatePicker, Button, Space } from 'antd';
import type { FormInstance } from 'antd';
import { SearchOutlined, DownOutlined, UpOutlined } from '@ant-design/icons';
import styles from './index.module.css';

const { RangePicker } = DatePicker;

export interface SearchField {
    name: string;
    label: string;
    type: 'input' | 'select' | 'dateRange' | 'date';
    placeholder?: string;
    options?: { label: string; value: string | number }[];
    span?: number;
    mode?: 'multiple' | 'tags';
}

export interface SearchFormProps {
    fields: SearchField[];
    form: FormInstance;
    onSearch: () => void;
    onReset: () => void;
    loading?: boolean;
    defaultExpanded?: boolean;
}

export const SearchForm: React.FC<SearchFormProps> = ({
    fields,
    form,
    onSearch,
    onReset,
    loading,
    defaultExpanded = true,
}) => {
    const [expanded, setExpanded] = useState(defaultExpanded);

    const renderSearchField = (field: SearchField) => {
        switch (field.type) {
            case 'select':
                return (
                    <Select
                        placeholder={field.placeholder || `请选择${field.label}`}
                        allowClear
                        mode={field.mode}
                        options={field.options}
                        style={{ width: '100%' }}
                    />
                );
            case 'dateRange':
                return (
                    <RangePicker
                        placeholder={['开始日期', '结束日期']}
                        style={{ width: '100%' }}
                    />
                );
            case 'date':
                return (
                    <DatePicker
                        placeholder={field.placeholder || `请选择${field.label}`}
                        style={{ width: '100%' }}
                    />
                );
            default:
                return (
                    <Input
                        placeholder={field.placeholder || `请输入${field.label}`}
                        allowClear
                    />
                );
        }
    };

    if (fields.length === 0) {
        return null;
    }

    return (
        <Form form={form} layout="horizontal">
            <Row gutter={[24, 16]}>
                {fields.slice(0, expanded ? undefined : 3).map(field => (
                    <Col
                        key={field.name}
                        xs={24}
                        sm={12}
                        md={8}
                        lg={6}
                        xl={6}
                    >
                        <Form.Item
                            name={field.name}
                            label={field.label}
                            className={styles.formItem}
                        >
                            {renderSearchField(field)}
                        </Form.Item>
                    </Col>
                ))}
                <Col flex="auto" style={{ textAlign: 'right' }}>
                    <Space>
                        <Button
                            type="primary"
                            icon={<SearchOutlined />}
                            onClick={onSearch}
                            loading={loading}
                        >
                            搜索
                        </Button>
                        <Button onClick={onReset}>重置</Button>
                        {fields.length > 3 && (
                            <Button
                                type="link"
                                onClick={() => setExpanded(!expanded)}
                                style={{ padding: 0 }}
                            >
                                {expanded ? '收起' : '展开'}
                                {expanded ? <UpOutlined /> : <DownOutlined />}
                            </Button>
                        )}
                    </Space>
                </Col>
            </Row>
        </Form >
    );
};
