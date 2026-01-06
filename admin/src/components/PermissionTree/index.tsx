/**
 * 权限树组件
 * 用于角色权限配置页面，支持按分组懒加载和虚拟滚动
 * Requirements: 2.1, 2.2
 *
 * 优化特性：
 * - 真正的懒加载：初始只渲染分组，展开时才渲染该分组的权限节点
 * - 虚拟滚动：大数据量时只渲染可见区域
 * - 简化渲染：使用简单文本标题，避免复杂 JSX
 */
import React, {
    useState,
    useCallback,
    useMemo,
    useEffect,
    useRef,
} from 'react';
import { Tree, Input, Empty, Spin, Button, Space, Typography, Skeleton } from 'antd';
import type { TreeProps } from 'antd';
import type { DataNode } from 'antd/es/tree';
import {
    SearchOutlined,
    CheckSquareOutlined,
    BorderOutlined,
    LockOutlined,
} from '@ant-design/icons';
import type { PermissionTreeNode, Permission } from '@/types/permission';
import { permissionApi } from '@/api/permission';
import styles from './index.module.less';

const { Text } = Typography;

/**
 * 权限树组件属性
 */
export interface PermissionTreeProps {
    /** 权限树数据（传统模式，一次性加载所有数据） */
    treeData?: PermissionTreeNode[];
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
    /**
     * 是否启用按分组懒加载模式
     * 启用后会忽略 treeData，自动从 API 加载数据
     * @default false
     */
    lazyLoadByGroup?: boolean;
}

/**
 * 后端返回的分组数据结构（兼容两种格式）
 */
type ApiPermissionGroup = {
    group?: string;
    name?: string;
    permissions: Permission[];
};

/**
 * 获取分组名称（兼容 group 和 name 两种字段）
 */
const getGroupName = (group: ApiPermissionGroup): string => {
    return group.group || group.name || 'unknown';
};

/**
 * 将权限树节点转换为简单数据格式（传统模式）
 */
const convertToTreeData = (
    nodes: PermissionTreeNode[],
    searchValue: string = ''
): DataNode[] => {
    return nodes.map((node) => {
        const title = node.description || node.code;
        const _isMatched = searchValue
            ? title.toLowerCase().includes(searchValue.toLowerCase()) ||
              node.code.toLowerCase().includes(searchValue.toLowerCase())
            : true;

        const children = node.children ? convertToTreeData(node.children, searchValue) : undefined;

        return {
            key: node.id,
            title: `${title} [${node.code}]`,
            children,
            selectable: false,
        } as DataNode;
    });
};

/**
 * 获取所有节点的 key（传统模式）
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
 * 权限树组件
 */
export const PermissionTree: React.FC<PermissionTreeProps> = React.memo(
    ({
        treeData = [],
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
        lazyLoadByGroup = false,
    }) => {
        const [searchValue, setSearchValue] = useState('');
        const [expandedKeys, setExpandedKeys] = useState<React.Key[]>([]);
        const [autoExpandParent, setAutoExpandParent] = useState(true);

        // 懒加载相关状态
        const [groupedData, setGroupedData] = useState<ApiPermissionGroup[]>([]);
        const [loadedGroups, setLoadedGroups] = useState<Set<string>>(new Set());
        const [isInitialLoading, setIsInitialLoading] = useState(lazyLoadByGroup);

        // 搜索防抖
        const searchTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
        const [debouncedSearchValue, setDebouncedSearchValue] = useState('');

        // 加载分组数据
        useEffect(() => {
            if (!lazyLoadByGroup) return;

            const loadGroups = async () => {
                setIsInitialLoading(true);
                try {
                    const res = await permissionApi.getTreeByGroup();
                    if (res.data.success && res.data.data) {
                        setGroupedData(res.data.data);
                        setLoadedGroups(new Set());
                    }
                } catch (error) {
                    console.error('Failed to load permission groups:', error);
                } finally {
                    setIsInitialLoading(false);
                }
            };

            loadGroups();
        }, [lazyLoadByGroup]);

        // 计算所有权限 ID（用于全选）
        const allPermissionIds = useMemo(() => {
            if (lazyLoadByGroup) {
                return groupedData.flatMap((g) => g.permissions.map((p) => p.id));
            }
            return getAllKeys(treeData);
        }, [lazyLoadByGroup, groupedData, treeData]);

        /**
         * 构建树数据 - 真正的懒加载
         * 只有展开的分组才会渲染其权限节点
         */
        const treeDataMemo = useMemo(() => {
            if (!lazyLoadByGroup) {
                return convertToTreeData(treeData, debouncedSearchValue);
            }

            // 懒加载模式：只渲染分组 + 已展开分组的权限
            return groupedData
                .map((group) => {
                    const groupName = getGroupName(group);
                    const groupKey = `group-${groupName}`;
                    const displayName = groupName.split('/').pop() || groupName;
                    const isExpanded = loadedGroups.has(groupName);

                    // 过滤匹配搜索的权限
                    const filteredPermissions = debouncedSearchValue
                        ? group.permissions.filter(
                              (p) =>
                                  p.description?.toLowerCase().includes(debouncedSearchValue.toLowerCase()) ||
                                  p.code.toLowerCase().includes(debouncedSearchValue.toLowerCase())
                          )
                        : group.permissions;

                    // 如果搜索时该分组没有匹配项，跳过
                    if (debouncedSearchValue && filteredPermissions.length === 0) {
                        return null;
                    }

                    // 构建子节点
                    let children: DataNode[];
                    if (isExpanded) {
                        // 已展开：渲染实际权限节点
                        children = filteredPermissions.map((p) => ({
                            key: p.id,
                            title: `${p.description || p.code} [${p.code}]`,
                            isLeaf: true,
                        }));
                    } else {
                        // 未展开：只显示一个占位节点
                        children = [
                            {
                                key: `${groupKey}-placeholder`,
                                title: `点击展开加载 ${group.permissions.length} 个权限...`,
                                isLeaf: true,
                                disabled: true,
                                checkable: false,
                            },
                        ];
                    }

                    return {
                        key: groupKey,
                        title: `${displayName} (${group.permissions.length})`,
                        children,
                        checkable: false,
                        selectable: false,
                    } as DataNode;
                })
                .filter(Boolean) as DataNode[];
        }, [lazyLoadByGroup, groupedData, loadedGroups, debouncedSearchValue, treeData]);

        // 搜索时自动展开所有分组
        useEffect(() => {
            if (debouncedSearchValue && lazyLoadByGroup) {
                const allGroupNames = new Set(groupedData.map((g) => getGroupName(g)));
                setLoadedGroups(allGroupNames);
                setExpandedKeys(groupedData.map((g) => `group-${getGroupName(g)}`));
                setAutoExpandParent(true);
            }
        }, [debouncedSearchValue, lazyLoadByGroup, groupedData]);

        /**
         * 处理搜索（带防抖）
         */
        const handleSearch = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
            const value = e.target.value;
            setSearchValue(value);

            if (searchTimeoutRef.current) {
                clearTimeout(searchTimeoutRef.current);
            }

            searchTimeoutRef.current = setTimeout(() => {
                setDebouncedSearchValue(value);
            }, 300);
        }, []);

        /**
         * 处理展开/收起 - 懒加载核心逻辑
         */
        const handleExpand: TreeProps['onExpand'] = useCallback(
            (keys: React.Key[]) => {
                setExpandedKeys(keys);
                setAutoExpandParent(false);

                if (lazyLoadByGroup) {
                    // 找出新展开的分组
                    const newExpandedGroups: string[] = [];
                    keys.forEach((key) => {
                        if (typeof key === 'string' && key.startsWith('group-')) {
                            const groupName = key.replace('group-', '');
                            if (!loadedGroups.has(groupName)) {
                                newExpandedGroups.push(groupName);
                            }
                        }
                    });

                    // 标记为已加载（触发重新渲染，显示实际权限节点）
                    if (newExpandedGroups.length > 0) {
                        setLoadedGroups((prev) => {
                            const next = new Set(prev);
                            newExpandedGroups.forEach((g) => next.add(g));
                            return next;
                        });
                    }
                }
            },
            [lazyLoadByGroup, loadedGroups]
        );

        /**
         * 处理选中状态变化
         */
        const handleCheck: TreeProps['onCheck'] = useCallback(
            (
                checked: React.Key[] | { checked: React.Key[]; halfChecked: React.Key[] },
                info: { checked: boolean; node: DataNode }
            ) => {
                if (disabled) return;

                const checkedArray = Array.isArray(checked) ? checked : checked.checked;
                // 过滤掉分组 key 和占位符 key，只保留权限 ID
                const numericKeys = checkedArray
                    .filter(
                        (key) =>
                            typeof key === 'number' ||
                            (typeof key === 'string' &&
                                !key.toString().startsWith('group-') &&
                                !key.toString().includes('-placeholder'))
                    )
                    .map((key) => Number(key));

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
            onCheck?.(allPermissionIds, { checked: true, node: {} as DataNode });
        }, [disabled, allPermissionIds, onCheck]);

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
            const totalCount = allPermissionIds.length;
            const isAllSelected = checkedCount === totalCount && totalCount > 0;

            return {
                checkedCount,
                totalCount,
                isAllSelected,
            };
        }, [checkedKeys.length, allPermissionIds.length]);

        // 显示加载状态
        const isLoading = loading || isInitialLoading;

        if (isLoading) {
            return (
                <div className={styles.loadingContainer}>
                    <div className={styles.skeletonTree}>
                        <Skeleton.Input active size="small" style={{ width: '100%', marginBottom: 8 }} />
                        <div style={{ paddingLeft: 24 }}>
                            <Skeleton.Input active size="small" style={{ width: '90%', marginBottom: 8 }} />
                            <Skeleton.Input active size="small" style={{ width: '85%', marginBottom: 8 }} />
                        </div>
                        <Skeleton.Input active size="small" style={{ width: '100%', marginBottom: 8 }} />
                    </div>
                    <div style={{ textAlign: 'center', marginTop: 16 }}>
                        <Spin tip="加载权限数据中..." size="default" />
                    </div>
                </div>
            );
        }

        const hasData = lazyLoadByGroup ? groupedData.length > 0 : treeData.length > 0;

        if (!hasData) {
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
                        showIcon={false}
                        showLine={false}
                        checkedKeys={checkedKeys.map(String)}
                        expandedKeys={expandedKeys}
                        autoExpandParent={autoExpandParent}
                        onExpand={handleExpand}
                        onCheck={handleCheck}
                        treeData={treeDataMemo}
                        height={virtual ? height : undefined}
                        virtual={virtual}
                        disabled={disabled}
                        selectable={false}
                    />
                </div>

                {/* 搜索无结果提示 */}
                {debouncedSearchValue && treeDataMemo.length === 0 && (
                    <Empty
                        description={`未找到包含 "${debouncedSearchValue}" 的权限`}
                        className={styles.emptySearch}
                    />
                )}
            </div>
        );
    }
);

export default PermissionTree;
