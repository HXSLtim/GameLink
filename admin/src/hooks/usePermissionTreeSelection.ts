/**
 * 权限树选择状态管理 Hook
 * 实现父子节点联动选择逻辑
 * Requirements: 2.2 - 选中父节点自动选中子节点，取消父节点自动取消子节点
 */
import { useState, useCallback, useMemo } from 'react';
import type { PermissionTreeNode } from '@/types/permission';

/**
 * 节点映射类型
 */
interface NodeMap {
    [key: number]: {
        node: PermissionTreeNode;
        parentId: number | null;
        childIds: number[];
    };
}

/**
 * 构建节点映射
 */
const buildNodeMap = (
    nodes: PermissionTreeNode[],
    parentId: number | null = null
): NodeMap => {
    const map: NodeMap = {};

    const traverse = (items: PermissionTreeNode[], parent: number | null) => {
        items.forEach((item) => {
            const childIds = item.children?.map((c) => c.id) || [];
            map[item.id] = {
                node: item,
                parentId: parent,
                childIds,
            };

            if (item.children) {
                traverse(item.children, item.id);
            }
        });
    };

    traverse(nodes, parentId);
    return map;
};

/**
 * 获取所有后代节点ID
 */
const getAllDescendantIds = (nodeId: number, nodeMap: NodeMap): number[] => {
    const descendants: number[] = [];
    const nodeInfo = nodeMap[nodeId];

    if (!nodeInfo) return descendants;

    nodeInfo.childIds.forEach((childId) => {
        descendants.push(childId);
        descendants.push(...getAllDescendantIds(childId, nodeMap));
    });

    return descendants;
};

/**
 * 获取所有祖先节点ID
 */
const getAllAncestorIds = (nodeId: number, nodeMap: NodeMap): number[] => {
    const ancestors: number[] = [];
    let currentId: number | null = nodeMap[nodeId]?.parentId ?? null;

    while (currentId !== null) {
        ancestors.push(currentId);
        currentId = nodeMap[currentId]?.parentId ?? null;
    }

    return ancestors;
};

/**
 * 检查节点的所有子节点是否都被选中
 */
const areAllChildrenChecked = (
    nodeId: number,
    checkedSet: Set<number>,
    nodeMap: NodeMap
): boolean => {
    const nodeInfo = nodeMap[nodeId];
    if (!nodeInfo || nodeInfo.childIds.length === 0) return true;

    return nodeInfo.childIds.every(
        (childId) =>
            checkedSet.has(childId) && areAllChildrenChecked(childId, checkedSet, nodeMap)
    );
};

/**
 * 权限树选择状态管理 Hook
 * 
 * @param treeData 权限树数据
 * @param initialCheckedKeys 初始选中的权限ID列表
 * @returns 选择状态和操作方法
 * 
 * @example
 * ```tsx
 * const { checkedKeys, halfCheckedKeys, handleCheck, selectAll, deselectAll } = 
 *     usePermissionTreeSelection(treeData, initialKeys);
 * ```
 */
export const usePermissionTreeSelection = (
    treeData: PermissionTreeNode[],
    initialCheckedKeys: number[] = []
) => {
    const [checkedKeys, setCheckedKeys] = useState<number[]>(initialCheckedKeys);

    // 构建节点映射
    const nodeMap = useMemo(() => buildNodeMap(treeData), [treeData]);

    // 获取所有节点ID
    const allNodeIds = useMemo(() => Object.keys(nodeMap).map(Number), [nodeMap]);

    // 获取所有叶子节点ID
    const leafNodeIds = useMemo(
        () => allNodeIds.filter((id) => nodeMap[id].childIds.length === 0),
        [allNodeIds, nodeMap]
    );

    // 计算半选状态的节点
    const halfCheckedKeys = useMemo(() => {
        const checkedSet = new Set(checkedKeys);
        const halfChecked: number[] = [];

        allNodeIds.forEach((nodeId) => {
            const nodeInfo = nodeMap[nodeId];
            if (nodeInfo.childIds.length === 0) return; // 叶子节点不需要半选状态

            const hasCheckedChild = nodeInfo.childIds.some(
                (childId) =>
                    checkedSet.has(childId) ||
                    getAllDescendantIds(childId, nodeMap).some((id) => checkedSet.has(id))
            );

            const allChildrenChecked = areAllChildrenChecked(nodeId, checkedSet, nodeMap);

            if (hasCheckedChild && !allChildrenChecked && !checkedSet.has(nodeId)) {
                halfChecked.push(nodeId);
            }
        });

        return halfChecked;
    }, [checkedKeys, allNodeIds, nodeMap]);

    /**
     * 处理节点选中/取消选中
     * Requirements: 2.2 - 选中父节点自动选中子节点，取消父节点自动取消子节点
     */
    const handleCheck = useCallback(
        (nodeId: number, checked: boolean) => {
            setCheckedKeys((prevKeys) => {
                const newCheckedSet = new Set(prevKeys);
                const descendants = getAllDescendantIds(nodeId, nodeMap);
                const ancestors = getAllAncestorIds(nodeId, nodeMap);

                if (checked) {
                    // 选中当前节点
                    newCheckedSet.add(nodeId);

                    // 选中所有后代节点
                    descendants.forEach((id) => newCheckedSet.add(id));

                    // 检查并更新祖先节点状态
                    ancestors.forEach((ancestorId) => {
                        if (areAllChildrenChecked(ancestorId, newCheckedSet, nodeMap)) {
                            newCheckedSet.add(ancestorId);
                        }
                    });
                } else {
                    // 取消选中当前节点
                    newCheckedSet.delete(nodeId);

                    // 取消选中所有后代节点
                    descendants.forEach((id) => newCheckedSet.delete(id));

                    // 取消选中所有祖先节点（因为子节点不完整了）
                    ancestors.forEach((id) => newCheckedSet.delete(id));
                }

                return Array.from(newCheckedSet);
            });
        },
        [nodeMap]
    );

    /**
     * 批量设置选中状态
     */
    const setChecked = useCallback((keys: number[]) => {
        setCheckedKeys(keys);
    }, []);

    /**
     * 全选
     */
    const selectAll = useCallback(() => {
        setCheckedKeys(allNodeIds);
    }, [allNodeIds]);

    /**
     * 清空选择
     */
    const deselectAll = useCallback(() => {
        setCheckedKeys([]);
    }, []);

    /**
     * 反选
     */
    const invertSelection = useCallback(() => {
        setCheckedKeys((prevKeys) => {
            const prevSet = new Set(prevKeys);
            return leafNodeIds.filter((id) => !prevSet.has(id));
        });
    }, [leafNodeIds]);

    /**
     * 检查节点是否被选中
     */
    const isChecked = useCallback(
        (nodeId: number) => checkedKeys.includes(nodeId),
        [checkedKeys]
    );

    /**
     * 检查节点是否为半选状态
     */
    const isHalfChecked = useCallback(
        (nodeId: number) => halfCheckedKeys.includes(nodeId),
        [halfCheckedKeys]
    );

    return {
        /** 选中的节点ID列表 */
        checkedKeys,
        /** 半选状态的节点ID列表 */
        halfCheckedKeys,
        /** 处理节点选中/取消选中 */
        handleCheck,
        /** 批量设置选中状态 */
        setChecked,
        /** 全选 */
        selectAll,
        /** 清空选择 */
        deselectAll,
        /** 反选 */
        invertSelection,
        /** 检查节点是否被选中 */
        isChecked,
        /** 检查节点是否为半选状态 */
        isHalfChecked,
        /** 所有节点ID */
        allNodeIds,
        /** 所有叶子节点ID */
        leafNodeIds,
    };
};

export default usePermissionTreeSelection;
