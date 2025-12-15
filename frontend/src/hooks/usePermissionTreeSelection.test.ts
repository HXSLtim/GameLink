/**
 * usePermissionTreeSelection Hook Property-Based Tests
 *
 * **Feature: rbac-button-level-permission, Property 4: 角色权限树形选择一致性**
 * **Validates: Requirements 2.2**
 *
 * Property: For any permission tree node, selecting a parent node should result
 * in all its child nodes being selected; deselecting a parent should deselect all children.
 */
import { describe, it } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import * as fc from 'fast-check';
import type { PermissionTreeNode } from '@/types/permission';

// Import the hook
import { usePermissionTreeSelection } from './usePermissionTreeSelection';

/**
 * Generate a valid PermissionTreeNode with unique IDs
 */
const generatePermissionNode = (
  id: number,
  depth: number,
  maxDepth: number,
  idCounter: { value: number }
): PermissionTreeNode => {
  const node: PermissionTreeNode = {
    id,
    method: 'GET',
    path: `/api/test/${id}`,
    code: `test.resource.action${id}`,
    group: 'test',
    description: `Test permission ${id}`,
    sortOrder: 0,
    isSystem: false,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  };

  // Generate children if not at max depth
  if (depth < maxDepth) {
    const numChildren = Math.floor(Math.random() * 3) + 1; // 1-3 children
    node.children = [];
    for (let i = 0; i < numChildren; i++) {
      idCounter.value++;
      node.children.push(
        generatePermissionNode(idCounter.value, depth + 1, maxDepth, idCounter)
      );
    }
  }

  return node;
};

/**
 * Arbitrary generator for permission tree
 */
const permissionTreeArb = fc
  .record({
    numRoots: fc.integer({ min: 1, max: 3 }),
    maxDepth: fc.integer({ min: 1, max: 3 }),
  })
  .map(({ numRoots, maxDepth }) => {
    const idCounter = { value: 0 };
    const roots: PermissionTreeNode[] = [];

    for (let i = 0; i < numRoots; i++) {
      idCounter.value++;
      roots.push(generatePermissionNode(idCounter.value, 0, maxDepth, idCounter));
    }

    return roots;
  });

/**
 * Get all node IDs from a tree
 */
const getAllNodeIds = (nodes: PermissionTreeNode[]): number[] => {
  const ids: number[] = [];
  const traverse = (items: PermissionTreeNode[]) => {
    items.forEach((item) => {
      ids.push(item.id);
      if (item.children) {
        traverse(item.children);
      }
    });
  };
  traverse(nodes);
  return ids;
};

/**
 * Get all descendant IDs of a node
 */
const getDescendantIds = (
  nodeId: number,
  nodes: PermissionTreeNode[]
): number[] => {
  const descendants: number[] = [];

  const findNode = (items: PermissionTreeNode[]): PermissionTreeNode | null => {
    for (const item of items) {
      if (item.id === nodeId) return item;
      if (item.children) {
        const found = findNode(item.children);
        if (found) return found;
      }
    }
    return null;
  };

  const collectDescendants = (node: PermissionTreeNode) => {
    if (node.children) {
      node.children.forEach((child) => {
        descendants.push(child.id);
        collectDescendants(child);
      });
    }
  };

  const node = findNode(nodes);
  if (node) {
    collectDescendants(node);
  }

  return descendants;
};

/**
 * Find a node with children in the tree
 */
const findNodeWithChildren = (
  nodes: PermissionTreeNode[]
): PermissionTreeNode | null => {
  for (const node of nodes) {
    if (node.children && node.children.length > 0) {
      return node;
    }
    if (node.children) {
      const found = findNodeWithChildren(node.children);
      if (found) return found;
    }
  }
  return null;
};

/**
 * Get all nodes with children
 */
const getAllNodesWithChildren = (
  nodes: PermissionTreeNode[]
): PermissionTreeNode[] => {
  const result: PermissionTreeNode[] = [];

  const traverse = (items: PermissionTreeNode[]) => {
    items.forEach((item) => {
      if (item.children && item.children.length > 0) {
        result.push(item);
        traverse(item.children);
      }
    });
  };

  traverse(nodes);
  return result;
};

describe('usePermissionTreeSelection Property Tests', () => {
  /**
   * **Feature: rbac-button-level-permission, Property 4: 角色权限树形选择一致性**
   * **Validates: Requirements 2.2**
   *
   * Property: Selecting a parent node should select all its descendants
   */
  it('Property 4a: selecting a parent node should select all its child nodes', () => {
    fc.assert(
      fc.property(permissionTreeArb, (treeData) => {
        // Find a node with children
        const parentNode = findNodeWithChildren(treeData);
        if (!parentNode) {
          // Skip if no parent nodes exist (tree is flat)
          return true;
        }

        const { result } = renderHook(() =>
          usePermissionTreeSelection(treeData, [])
        );

        // Select the parent node
        act(() => {
          result.current.handleCheck(parentNode.id, true);
        });

        // Get all descendant IDs
        const descendantIds = getDescendantIds(parentNode.id, treeData);

        // Verify parent is selected
        const parentSelected = result.current.checkedKeys.includes(parentNode.id);

        // Verify all descendants are selected
        const allDescendantsSelected = descendantIds.every((id) =>
          result.current.checkedKeys.includes(id)
        );

        return parentSelected && allDescendantsSelected;
      }),
      { numRuns: 100 }
    );
  });

  /**
   * **Feature: rbac-button-level-permission, Property 4: 角色权限树形选择一致性**
   * **Validates: Requirements 2.2**
   *
   * Property: Deselecting a parent node should deselect all its descendants
   */
  it('Property 4b: deselecting a parent node should deselect all its child nodes', () => {
    fc.assert(
      fc.property(permissionTreeArb, (treeData) => {
        // Find a node with children
        const parentNode = findNodeWithChildren(treeData);
        if (!parentNode) {
          return true;
        }

        // Start with all nodes selected
        const allIds = getAllNodeIds(treeData);
        const { result } = renderHook(() =>
          usePermissionTreeSelection(treeData, allIds)
        );

        // Deselect the parent node
        act(() => {
          result.current.handleCheck(parentNode.id, false);
        });

        // Get all descendant IDs
        const descendantIds = getDescendantIds(parentNode.id, treeData);

        // Verify parent is deselected
        const parentDeselected = !result.current.checkedKeys.includes(parentNode.id);

        // Verify all descendants are deselected
        const allDescendantsDeselected = descendantIds.every(
          (id) => !result.current.checkedKeys.includes(id)
        );

        return parentDeselected && allDescendantsDeselected;
      }),
      { numRuns: 100 }
    );
  });

  /**
   * Property: Selecting all children should auto-select the parent
   */
  it('Property 4c: selecting all children should auto-select the parent', () => {
    fc.assert(
      fc.property(permissionTreeArb, (treeData) => {
        const parentNode = findNodeWithChildren(treeData);
        if (!parentNode || !parentNode.children) {
          return true;
        }

        const { result } = renderHook(() =>
          usePermissionTreeSelection(treeData, [])
        );

        // Select all children one by one
        act(() => {
          parentNode.children!.forEach((child) => {
            // Select child and all its descendants
            result.current.handleCheck(child.id, true);
          });
        });

        // Parent should be auto-selected when all children are selected
        return result.current.checkedKeys.includes(parentNode.id);
      }),
      { numRuns: 100 }
    );
  });

  /**
   * Property: Deselecting any child should deselect the parent
   */
  it('Property 4d: deselecting any child should deselect the parent', () => {
    fc.assert(
      fc.property(permissionTreeArb, (treeData) => {
        const parentNode = findNodeWithChildren(treeData);
        if (!parentNode || !parentNode.children || parentNode.children.length === 0) {
          return true;
        }

        // Start with all nodes selected
        const allIds = getAllNodeIds(treeData);
        const { result } = renderHook(() =>
          usePermissionTreeSelection(treeData, allIds)
        );

        // Deselect one child
        const childToDeselect = parentNode.children[0];
        act(() => {
          result.current.handleCheck(childToDeselect.id, false);
        });

        // Parent should be deselected
        return !result.current.checkedKeys.includes(parentNode.id);
      }),
      { numRuns: 100 }
    );
  });

  /**
   * Property: selectAll should select all nodes
   */
  it('Property 4e: selectAll should select all nodes in the tree', () => {
    fc.assert(
      fc.property(permissionTreeArb, (treeData) => {
        const { result } = renderHook(() =>
          usePermissionTreeSelection(treeData, [])
        );

        act(() => {
          result.current.selectAll();
        });

        const allIds = getAllNodeIds(treeData);
        return allIds.every((id) => result.current.checkedKeys.includes(id));
      }),
      { numRuns: 100 }
    );
  });

  /**
   * Property: deselectAll should deselect all nodes
   */
  it('Property 4f: deselectAll should deselect all nodes in the tree', () => {
    fc.assert(
      fc.property(permissionTreeArb, (treeData) => {
        // Start with all selected
        const allIds = getAllNodeIds(treeData);
        const { result } = renderHook(() =>
          usePermissionTreeSelection(treeData, allIds)
        );

        act(() => {
          result.current.deselectAll();
        });

        return result.current.checkedKeys.length === 0;
      }),
      { numRuns: 100 }
    );
  });

  /**
   * Property: Half-checked state should be correct
   */
  it('Property 4g: parent should be half-checked when some but not all children are selected', () => {
    fc.assert(
      fc.property(permissionTreeArb, (treeData) => {
        const parentNode = findNodeWithChildren(treeData);
        if (!parentNode || !parentNode.children || parentNode.children.length < 2) {
          return true;
        }

        const { result } = renderHook(() =>
          usePermissionTreeSelection(treeData, [])
        );

        // Select only the first child
        act(() => {
          result.current.handleCheck(parentNode.children![0].id, true);
        });

        // Parent should be in half-checked state (not fully checked)
        const parentNotFullyChecked = !result.current.checkedKeys.includes(parentNode.id);
        const parentHalfChecked = result.current.halfCheckedKeys.includes(parentNode.id);

        return parentNotFullyChecked && parentHalfChecked;
      }),
      { numRuns: 100 }
    );
  });

  /**
   * Property: Selection should be idempotent when node and all descendants are already selected
   */
  it('Property 4h: selecting an already fully selected subtree should not change state', () => {
    fc.assert(
      fc.property(permissionTreeArb, (treeData) => {
        const allIds = getAllNodeIds(treeData);
        if (allIds.length === 0) return true;

        const nodeId = allIds[0];
        
        // Start with all nodes selected (fully selected state)
        const { result } = renderHook(() =>
          usePermissionTreeSelection(treeData, allIds)
        );

        const beforeKeys = [...result.current.checkedKeys].sort((a, b) => a - b);

        act(() => {
          result.current.handleCheck(nodeId, true);
        });

        const afterKeys = [...result.current.checkedKeys].sort((a, b) => a - b);

        // Keys should be the same (idempotent when already fully selected)
        return (
          beforeKeys.length === afterKeys.length &&
          beforeKeys.every((k, i) => afterKeys[i] === k)
        );
      }),
      { numRuns: 100 }
    );
  });

  /**
   * Property: Deselection should be idempotent
   */
  it('Property 4i: deselecting an already deselected node should not change state', () => {
    fc.assert(
      fc.property(permissionTreeArb, (treeData) => {
        const allIds = getAllNodeIds(treeData);
        if (allIds.length === 0) return true;

        const nodeId = allIds[0];
        const { result } = renderHook(() =>
          usePermissionTreeSelection(treeData, [])
        );

        const beforeKeys = [...result.current.checkedKeys];

        act(() => {
          result.current.handleCheck(nodeId, false);
        });

        const afterKeys = result.current.checkedKeys;

        return (
          beforeKeys.length === afterKeys.length &&
          beforeKeys.every((k) => afterKeys.includes(k))
        );
      }),
      { numRuns: 100 }
    );
  });

  /**
   * Property: Tree structure consistency - all ancestors should be checked
   * when all their descendants are checked
   */
  it('Property 4j: all ancestors should be checked when all descendants are checked', () => {
    fc.assert(
      fc.property(permissionTreeArb, (treeData) => {
        const nodesWithChildren = getAllNodesWithChildren(treeData);
        if (nodesWithChildren.length === 0) return true;

        // Select all nodes
        const allIds = getAllNodeIds(treeData);
        const { result } = renderHook(() =>
          usePermissionTreeSelection(treeData, allIds)
        );

        // All parent nodes should be checked
        return nodesWithChildren.every((node) =>
          result.current.checkedKeys.includes(node.id)
        );
      }),
      { numRuns: 100 }
    );
  });
});

