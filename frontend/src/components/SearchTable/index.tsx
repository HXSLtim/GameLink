/**
 * 搜索表格组件
 * 封装表格搜索、分页、批量操作等通用功能
 */
import React, { useState, useCallback, useEffect, ReactNode } from 'react';
import {
    Table,
    Card,
    Form,
    Row,
    Col,
    Button,
    Space,
    Input,
    Select,
    DatePicker,
    message,
    Popconfirm,
    Tooltip,
} from 'antd';
import type { TableProps, FormInstance } from 'antd';
import type { ColumnType } from 'antd/es/table';
import {
    SearchOutlined,
    ReloadOutlined,
    PlusOutlined,
    DeleteOutlined,
    DownOutlined,
    UpOutlined,
} from '@ant-design/icons';
import { PermissionGuard } from '@/components/PermissionGuard';
import styles from './index.module.css';

const { RangePicker } = DatePicker;

/**
 * 搜索字段配置
 */
export interface SearchField {
    /** 字段名 */
    name: string;
    /** 标签 */
    label: string;
    /** 类型 */
    type: 'input' | 'select' | 'dateRange' | 'date';
    /** 占位符 */
    placeholder?: string;
    /** 选项（select类型使用） */
    options?: { label: string; value: string | number }[];
    /** 栅格列数 */
    span?: number;
}

/**
 * 工具栏按钮配置
 */
export interface ToolbarButton {
    /** 按钮文本 */
    text: string;
    /** 按钮类型 */
    type?: 'primary' | 'default' | 'dashed' | 'link' | 'text';
    /** 图标 */
    icon?: ReactNode;
    /** 危险按钮 */
    danger?: boolean;
    /** 权限码 */
    permission?: string;
    /** 点击回调 */
    onClick?: () => void;
    /** 是否需要选中行 */
    needSelection?: boolean;
    /** 确认提示（批量删除等危险操作） */
    confirmText?: string;
}

/**
 * SearchTable组件属性
 */
export interface SearchTableProps<T> extends Omit<TableProps<T>, 'columns'> {
    /** 表格列配置 */
    columns: ColumnType<T>[];
    /** 搜索字段配置 */
    searchFields?: SearchField[];
    /** 工具栏按钮 */
    toolbarButtons?: ToolbarButton[];
    /** 是否显示新增按钮 */
    showCreate?: boolean;
    /** 新增按钮文本 */
    createText?: string;
    /** 新增权限码 */
    createPermission?: string;
    /** 新增回调 */
    onCreate?: () => void;
    /** 是否显示批量删除 */
    showBatchDelete?: boolean;
    /** 批量删除权限码 */
    batchDeletePermission?: string;
    /** 批量删除回调 */
    onBatchDelete?: (keys: React.Key[]) => Promise<void>;
    /** 搜索回调 */
    onSearch?: (values: Record<string, unknown>) => void;
    /** 刷新回调 */
    onRefresh?: () => void;
    /** 是否正在加载 */
    loading?: boolean;
    /** 表单实例 */
    form?: FormInstance;
    /** 卡片标题 */
    cardTitle?: ReactNode;
    /** 是否默认展开搜索 */
    defaultExpanded?: boolean;
}

/**
 * SearchTable组件
 */
function SearchTable<T extends object>({
    columns,
    searchFields = [],
    toolbarButtons = [],
    showCreate = true,
    createText = '新增',
    createPermission,
    onCreate,
    showBatchDelete = false,
    batchDeletePermission,
    onBatchDelete,
    onSearch,
    onRefresh,
    loading = false,
    form: externalForm,
    cardTitle,
    defaultExpanded = true,
    rowSelection,
    ...tableProps
}: SearchTableProps<T>) {
    const [internalForm] = Form.useForm();
    const form = externalForm || internalForm;
    const [expanded, setExpanded] = useState(defaultExpanded);
    const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

    // 搜索
    const handleSearch = useCallback(() => {
        const values = form.getFieldsValue();
        onSearch?.(values);
    }, [form, onSearch]);

    // 重置
    const handleReset = useCallback(() => {
        form.resetFields();
        handleSearch();
    }, [form, handleSearch]);

    // 刷新
    const handleRefresh = useCallback(() => {
        onRefresh?.();
    }, [onRefresh]);

    // 批量删除
    const handleBatchDelete = useCallback(async () => {
        if (selectedRowKeys.length === 0) {
            message.warning('请选择要删除的数据');
            return;
        }
        try {
            await onBatchDelete?.(selectedRowKeys);
            setSelectedRowKeys([]);
            message.success('删除成功');
        } catch {
            message.error('删除失败');
        }
    }, [selectedRowKeys, onBatchDelete]);

    // 行选择配置
    const mergedRowSelection = rowSelection !== undefined ? rowSelection : (showBatchDelete ? {
        selectedRowKeys,
        onChange: (keys: React.Key[]) => setSelectedRowKeys(keys),
    } : undefined);

    // 搜索表单回车提交
    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                handleSearch();
            }
        };
        document.addEventListener('keydown', handleKeyDown);
        return () => document.removeEventListener('keydown', handleKeyDown);
    }, [handleSearch]);

    /**
     * 渲染搜索字段
     */
    const renderSearchField = (field: SearchField) => {
        switch (field.type) {
            case 'select':
                return (
                    <Select
                        placeholder={field.placeholder || `请选择${field.label}`}
                        allowClear
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

    /**
     * 渲染工具栏按钮
     */
    const renderToolbarButton = (btn: ToolbarButton, index: number) => {
        const disabled = btn.needSelection && selectedRowKeys.length === 0;

        const button = btn.confirmText ? (
            <Popconfirm
                title={btn.confirmText}
                onConfirm={btn.onClick}
                disabled={disabled}
            >
                <Button
                    key={index}
                    type={btn.type}
                    icon={btn.icon}
                    danger={btn.danger}
                    disabled={disabled}
                >
                    {btn.text}
                    {btn.needSelection && selectedRowKeys.length > 0 && ` (${selectedRowKeys.length})`}
                </Button>
            </Popconfirm>
        ) : (
            <Button
                key={index}
                type={btn.type}
                icon={btn.icon}
                danger={btn.danger}
                onClick={btn.onClick}
                disabled={disabled}
            >
                {btn.text}
                {btn.needSelection && selectedRowKeys.length > 0 && ` (${selectedRowKeys.length})`}
            </Button>
        );

        if (btn.permission) {
            return (
                <PermissionGuard key={index} permission={btn.permission}>
                    {button}
                </PermissionGuard>
            );
        }

        return button;
    };

    return (
        <div className={styles.container}>
            {/* 搜索区域 */}
            {searchFields.length > 0 && (
                <Card className={styles.searchCard} size="small">
                    <Form form={form} layout="horizontal">
                        <Row gutter={16}>
                            {searchFields.slice(0, expanded ? undefined : 3).map(field => (
                                <Col key={field.name} span={field.span || 6}>
                                    <Form.Item
                                        name={field.name}
                                        label={field.label}
                                        className={styles.formItem}
                                    >
                                        {renderSearchField(field)}
                                    </Form.Item>
                                </Col>
                            ))}
                            <Col span={6} className={styles.searchActions}>
                                <Space>
                                    <Button
                                        type="primary"
                                        icon={<SearchOutlined />}
                                        onClick={handleSearch}
                                        loading={loading}
                                    >
                                        搜索
                                    </Button>
                                    <Button onClick={handleReset}>重置</Button>
                                    {searchFields.length > 3 && (
                                        <Button
                                            type="link"
                                            onClick={() => setExpanded(!expanded)}
                                            icon={expanded ? <UpOutlined /> : <DownOutlined />}
                                        >
                                            {expanded ? '收起' : '展开'}
                                        </Button>
                                    )}
                                </Space>
                            </Col>
                        </Row>
                    </Form>
                </Card>
            )}

            {/* 表格区域 */}
            <Card
                title={cardTitle}
                className={styles.tableCard}
                extra={
                    <Space>
                        {/* 自定义工具栏按钮 */}
                        {toolbarButtons.map(renderToolbarButton)}

                        {/* 批量删除 */}
                        {showBatchDelete && batchDeletePermission && (
                            <PermissionGuard permission={batchDeletePermission}>
                                <Popconfirm
                                    title={`确定要删除选中的 ${selectedRowKeys.length} 条数据吗？`}
                                    onConfirm={handleBatchDelete}
                                    disabled={selectedRowKeys.length === 0}
                                >
                                    <Button
                                        icon={<DeleteOutlined />}
                                        danger
                                        disabled={selectedRowKeys.length === 0}
                                    >
                                        批量删除 {selectedRowKeys.length > 0 && `(${selectedRowKeys.length})`}
                                    </Button>
                                </Popconfirm>
                            </PermissionGuard>
                        )}

                        {/* 新增按钮 */}
                        {showCreate && (
                            <PermissionGuard permission={createPermission || ''}>
                                <Button
                                    type="primary"
                                    icon={<PlusOutlined />}
                                    onClick={onCreate}
                                >
                                    {createText}
                                </Button>
                            </PermissionGuard>
                        )}

                        {/* 刷新按钮 */}
                        <Tooltip title="刷新">
                            <Button
                                icon={<ReloadOutlined spin={loading} />}
                                onClick={handleRefresh}
                            />
                        </Tooltip>
                    </Space>
                }
            >
                <Table<T>
                    columns={columns}
                    loading={loading}
                    rowSelection={mergedRowSelection}
                    {...tableProps}
                />
            </Card>
        </div>
    );
}

export default SearchTable;
