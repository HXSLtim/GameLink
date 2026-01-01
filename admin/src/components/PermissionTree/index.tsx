/**
 * 权限树组件
 * 用于角色权限配置页面，支持虚拟滚动和父子节点联动选择
 * Requirements: 2.1, 2.2
 */
import React, { useState, useCallback, useMemo } from 'react';
import { Tree, Input, Empty, Spin, Button, Space, Tag, Typography } from 'antd';
import type { TreeProps } from 'antd';
import type { DataNode } from 'antd/es/tree';
import {
    SearchOutlined,
    CheckSquareOutlined,
    BorderOutlined,
    FolderOutlined,
    FileProtectOutlined,
    LockOutlined,
} from '@ant-design/icons';
import type { PermissionTreeNode } from '@/types/permission';
import styles from './index.module.less';

const { Text } = Typography;

/**
 * 权限树组件属性
 */
export interface PermissionTreeProps {
    /** 权限树数据 */
    treeData: PermissionTreeNode[];
    /** 已选中的权限ID列表 */
    checkedKeys?: number[];
    /** 选中状态变化回调 */
    onCheck?: (checkedKeys: number[], info: { checked: boolean; node: DataNode }) => void;
    /** 是否加载中 */
    loading?: boolean;
    /** 是否禁用 */
    disabled?: boolean;
    /** 是否显示搜索框 */
    showSearch?: boolean;
    /** 是否显示全选/反选按钮 */
    showSelectAll?: boolean;
    /** 树的高度（用于虚拟滚动） */
    height?: number;
    /** 是否使用虚拟滚动 */
    virtual?: boolean;
    /** 自定义类名 */
    className?: string;
    /** 是否为系统角色（显示特殊提示） */
    isSystemRole?: boolean;
}

/**
 * 将权限树节点转换为 Ant Design Tree 数据格式
 * Requirements: 2.4 - 高亮显示已分配的权限项
 */
const convertToTreeData = (
    nodes: PermissionTreeNode[],
    searchValue: string = '',
    checkedKeys: number[] = []
): DataNode[] => {
    return nodes.map((node) => {
        const title = node.description || node.code;
        const isMatched = searchValue
            ? title.toLowerCase().includes(searchValue.toLowerCase()) ||
              node.code.toLowerCase().includes(searchValue.toLowerCase())
            : true;

        const children = node.children
            ? convertToTreeData(node.children, searchValue, checkedKeys)
            : undefined;

        // 检查子节点是否有匹配项
        const hasMatchedChildren = children?.some(
            (child) => (child as DataNode & { isMatched?: boolean }).isMatched
        );

        // 检查当前节点是否被选中（用于高亮显示）
        const isChecked = checkedKeys.includes(node.id);

        return {
            key: node.id,
            title: (
                <span className={`${styles.treeNodeTitle} ${isChecked ? styles.checkedNode : ''}`}>
                    <span className={isMatched && searchValue ? styles.highlight : ''}>
                        {title}
                    </span>
                    <Text type="secondary" className={styles.permissionCode}>
                        {node.code}
                    </Text>
                    {node.isSystem && (
                        <Tag color="blue" className={styles.systemTag}>
                            <LockOutlined /> 系统
                        </Tag>
                    )}
                    {isChecked && (
                        <Tag color="green" className={styles.assignedTag}>
                            已分配
                        </Tag>
                    )}
                </span>
            ),
            icon: node.children?.length ? <FolderOutlined /> : <FileProtectOutlined />,
            children,
            isMatched: isMatched || hasMatchedChildren,
            disabled: false,
            selectable: false,
        } as DataNode & { isMatched?: boolean };
    });
};

/**
 * 获取所有节点的 key
 */
const getAllKeys = (nodes: PermissionTreeNode[]): number[] => {
    const keys: number[] = [];
    const traverse = (items: PermissionTreeNode[]) => {
        items.forEach((item) => {
            keys.push(item.id);
            if (item.children) {
                traverse(item.children);
            }
        });
    };
    traverse(nodes);
    return keys;
};

/**
 * 获取所有叶子节点的 key
 */
const getLeafKeys = (nodes: PermissionTreeNode[]): number[] => {
    const keys: number[] = [];
    const traverse = (items: PermissionTreeNode[]) => {
        items.forEach((item) => {
            if (!item.children || item.children.length === 0) {
                keys.push(item.id);
            } else {
                traverse(item.children);
            }
        });
    };
    traverse(nodes);
    return keys;
};

/**
 * 过滤树数据，只保留匹配的节点
 */
const filterTreeData = (
    nodes: DataNode[],
    searchValue: string
): DataNode[] => {
    if (!searchValue) return nodes;

    return nodes
        .filter((node) => (node as DataNode & { isMatched?: boolean }).isMatched)
        .map((node) => ({
            ...node,
            children: node.children
                ? filterTreeData(node.children, searchValue)
                : undefined,
        }));
};

/**
 * 权限树组件
 * 支持虚拟滚动、搜索、全选/反选功能
 * Requirements: 2.1 - 以树形结构展示所有权限，按模块和资源分组
 *
 * 优化: 使用 React.memo 避免不必要的重新渲染
 * 适用场景: 大型树组件，仅在关键 props 变化时重新渲染
 */
export const PermissionTree: React.FC<PermissionTreeProps> = React.memo(({
    treeData,
    checkedKeys = [],
    onCheck,
    loading = false,
    disabled = false,
    showSearch = true,
    showSelectAll = true,
    height = 400,
    virtual = true,
    className,
    isSystemRole = false,
}) => {
    const [searchValue, setSearchValue] = useState('');
    const [expandedKeys, setExpandedKeys] = useState<React.Key[]>([]);
    const [autoExpandParent, setAutoExpandParent] = useState(true);

    // 转换树数据（包含已选中状态用于高亮显示）
    // Requirements: 2.4 - 高亮显示已分配的权限项
    const convertedTreeData = useMemo(
        () => convertToTreeData(treeData, searchValue, checkedKeys),
        [treeData, searchValue, checkedKeys]
    );

    // 过滤后的树数据
    const filteredTreeData = useMemo(
        () => filterTreeData(convertedTreeData, searchValue),
        [convertedTreeData, searchValue]
    );

    // 所有节点的 key
    const allKeys = useMemo(() => getAllKeys(treeData), [treeData]);

    // 所有叶子节点的 key
    const leafKeys = useMemo(() => getLeafKeys(treeData), [treeData]);

    // 当初始化时设置展开的 keys
    React.useLayoutEffect(() => {
        if (treeData.length > 0 && expandedKeys.length === 0) {
            setExpandedKeys(allKeys.map(String));
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [treeData.length]);

    // 搜索时自动展开匹配的节点
    React.useLayoutEffect(() => {
        if (searchValue) {
            setExpandedKeys(allKeys.map(String));
            setAutoExpandParent(true);
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [searchValue]);

    /**
     * 处理搜索
     */
    const handleSearch = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
        setSearchValue(e.target.value);
    }, []);

    /**
     * 处理展开/收起
     */
    const handleExpand: TreeProps['onExpand'] = useCallback((keys: React.Key[]) => {
        setExpandedKeys(keys);
        setAutoExpandParent(false);
    }, []);

    /**
     * 处理选中状态变化
     * Requirements: 2.2 - 支持批量选择（选中父节点自动选中所有子节点）
     */
    const handleCheck: TreeProps['onCheck'] = useCallback(
        (checked: React.Key[] | { checked: React.Key[]; halfChecked: React.Key[] }, info: { checked: boolean; node: DataNode }) => {
            if (disabled) return;

            const checkedArray = Array.isArray(checked) ? checked : checked.checked;
            const numericKeys = checkedArray.map((key) => Number(key));

            onCheck?.(numericKeys, {
                checked: info.checked,
                node: info.node,
            });
        },
        [disabled, onCheck]
    );

    /**
     * 全选
     */
    const handleSelectAll = useCallback(() => {
        if (disabled) return;
        onCheck?.(leafKeys, { checked: true, node: {} as DataNode });
    }, [disabled, leafKeys, onCheck]);

    /**
     * 反选
     */
    const handleDeselectAll = useCallback(() => {
        if (disabled) return;
        onCheck?.([], { checked: false, node: {} as DataNode });
    }, [disabled, onCheck]);

    /**
     * 计算选中状态
     */
    const selectionStatus = useMemo(() => {
        const checkedCount = checkedKeys.length;
        const totalCount = leafKeys.length;
        const isAllSelected = checkedCount === totalCount && totalCount > 0;
        const isPartialSelected = checkedCount > 0 && checkedCount < totalCount;

        return {
            checkedCount,
            totalCount,
            isAllSelected,
            isPartialSelected,
        };
    }, [checkedKeys, leafKeys]);

    if (loading) {
        return (
            <div className={styles.loadingContainer}>
                <Spin tip="加载权限数据中..." />
            </div>
        );
    }

    if (treeData.length === 0) {
        return <Empty description="暂无权限数据" />;
    }

    return (
        <div className={`${styles.permissionTree} ${className || ''}`}>
            {/* 系统角色提示 */}
            {isSystemRole && (
                <div className={styles.systemRoleAlert}>
                    <LockOutlined /> 超级管理员默认拥有所有权限，无需单独配置
                </div>
            )}

            {/* 搜索框 */}
            {showSearch && (
                <Input
                    placeholder="搜索权限..."
                    prefix={<SearchOutlined />}
                    value={searchValue}
                    onChange={handleSearch}
                    allowClear
                    className={styles.searchInput}
                    disabled={disabled}
                />
            )}

            {/* 全选/反选按钮 */}
            {showSelectAll && (
                <div className={styles.selectAllBar}>
                    <Space>
                        <Button
                            size="small"
                            icon={<CheckSquareOutlined />}
                            onClick={handleSelectAll}
                            disabled={disabled || selectionStatus.isAllSelected}
                        >
                            全选
                        </Button>
                        <Button
                            size="small"
                            icon={<BorderOutlined />}
                            onClick={handleDeselectAll}
                            disabled={disabled || selectionStatus.checkedCount === 0}
                        >
                            清空
                        </Button>
                    </Space>
                    <Text type="secondary" className={styles.selectionInfo}>
                        已选 {selectionStatus.checkedCount} / {selectionStatus.totalCount} 项
                    </Text>
                </div>
            )}

            {/* 权限树 */}
            <div className={styles.treeContainer}>
                <Tree
                    checkable
                    showIcon
                    showLine={{ showLeafIcon: false }}
                    checkedKeys={checkedKeys.map(String)}
                    expandedKeys={expandedKeys}
                    autoExpandParent={autoExpandParent}
                    onExpand={handleExpand}
                    onCheck={handleCheck}
                    treeData={filteredTreeData}
                    height={virtual ? height : undefined}
                    virtual={virtual}
                    disabled={disabled}
                    selectable={false}
                />
            </div>

            {/* 搜索无结果提示 */}
            {searchValue && filteredTreeData.length === 0 && (
                <Empty
                    description={`未找到包含 "${searchValue}" 的权限`}
                    className={styles.emptySearch}
                />
            )}
        </div>
    );
});

export default PermissionTree;
