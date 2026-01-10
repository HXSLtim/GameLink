/**
 * 统一的搜索/筛选组件
 * 标准化所有 CRUD 页面的筛选区域布局
 */
import React from 'react';
import type { ReactNode } from 'react';
import { Card, Row, Col, Space, Button, Segmented, Input, Select, DatePicker } from 'antd';
import { ReloadOutlined, SearchOutlined } from '@ant-design/icons';
import dayjs from 'dayjs';

const { RangePicker } = DatePicker;

// 定义筛选器类型
export interface FilterItem {
    type: 'input' | 'select' | 'rangePicker' | 'segmented';
    key: string;
    placeholder?: string;
    options?: Array<{ label: string; value: string | number }>;
    style?: React.CSSProperties;
    segmentedOptions?: Array<{ label: string; value: string }>;
    allowClear?: boolean;
    width?: number | string;
}

export interface SearchFiltersProps {
    /** 筛选器配置 */
    filters: FilterItem[];
    /** 筛选值变化回调 */
    onFilterChange: (key: string, value: unknown) => void;
    /** 筛选值对象 */
    filterValues: Record<string, unknown>;
    /** 操作按钮（右侧） */
    actions?: ReactNode;
    /** 是否使用 Card 包裹 */
    card?: boolean;
    /** 是否显示查询/重置按钮 */
    showQueryButtons?: boolean;
    /** 查询回调 */
    onQuery?: () => void;
    /** 重置回调 */
    onReset?: () => void;
    /** 加载状态 */
    loading?: boolean;
}

const SearchFilters: React.FC<SearchFiltersProps> = ({
    filters,
    onFilterChange,
    filterValues,
    actions,
    card = true,
    showQueryButtons = false,
    onQuery,
    onReset,
    loading = false,
}) => {
    const renderFilter = (filter: FilterItem) => {
        const value = filterValues[filter.key];

        switch (filter.type) {
            case 'input':
                return (
                    <Input
                        key={filter.key}
                        placeholder={filter.placeholder}
                        prefix={<SearchOutlined />}
                        allowClear={filter.allowClear ?? true}
                        style={{ width: filter.width || 200, ...filter.style }}
                        value={value as string}
                        onChange={(e) => onFilterChange(filter.key, e.target.value)}
                        onPressEnter={onQuery}
                    />
                );

            case 'select':
                return (
                    <Select
                        key={filter.key}
                        placeholder={filter.placeholder}
                        allowClear={filter.allowClear ?? true}
                        style={{ width: filter.width || 120, ...filter.style }}
                        value={value}
                        onChange={(v) => onFilterChange(filter.key, v)}
                        options={filter.options}
                    />
                );

            case 'rangePicker':
                return (
                    <RangePicker
                        key={filter.key}
                        style={filter.style}
                        value={
                            value
                                ? [
                                      (value as [Date | null, Date | null])[0]
                                          ? dayjs((value as [Date | null, Date | null])[0])
                                          : null,
                                      (value as [Date | null, Date | null])[1]
                                          ? dayjs((value as [Date | null, Date | null])[1])
                                          : null,
                                  ]
                                : null
                        }
                        onChange={(dates) => onFilterChange(filter.key, dates)}
                    />
                );

            case 'segmented':
                return (
                    <Segmented
                        key={filter.key}
                        options={filter.segmentedOptions || []}
                        value={value}
                        onChange={(v) => onFilterChange(filter.key, v)}
                    />
                );

            default:
                return null;
        }
    };

    const content = (
        <Row gutter={16} justify="space-between" align="middle">
            {/* 左侧：筛选器 */}
            <Col flex="auto">
                <Space wrap>
                    {filters.map(renderFilter)}
                    {showQueryButtons && (
                        <>
                            <Button type="primary" icon={<SearchOutlined />} onClick={onQuery} loading={loading}>
                                查询
                            </Button>
                            <Button icon={<ReloadOutlined />} onClick={onReset}>
                                重置
                            </Button>
                        </>
                    )}
                </Space>
            </Col>

            {/* 右侧：操作按钮 */}
            {actions && <Col>{actions}</Col>}
        </Row>
    );

    if (card) {
        return <Card style={{ marginBottom: 16 }}>{content}</Card>;
    }

    return <div style={{ marginBottom: 16 }}>{content}</div>;
};

export default SearchFilters;
