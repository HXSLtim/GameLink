import React, { useState, useMemo } from 'react';
import styles from './Tree.module.less';

export interface TreeNode {
  key: string;
  title: React.ReactNode;
  children?: TreeNode[];
  disabled?: boolean;
}

export interface TreeProps {
  data: TreeNode[];
  defaultExpandedKeys?: string[];
  expandedKeys?: string[];
  onExpand?: (keys: string[]) => void;
  selectable?: boolean;
  selectedKeys?: string[];
  defaultSelectedKeys?: string[];
  onSelect?: (keys: string[]) => void;
  className?: string;
  style?: React.CSSProperties;
}

export const Tree: React.FC<TreeProps> = ({
  data,
  defaultExpandedKeys = [],
  expandedKeys,
  onExpand,
  selectable = true,
  selectedKeys,
  defaultSelectedKeys = [],
  onSelect,
  className = '',
  style,
}) => {
  const [internalExpanded, setInternalExpanded] = useState<string[]>(defaultExpandedKeys);
  const [internalSelected, setInternalSelected] = useState<string[]>(defaultSelectedKeys);

  const isExpanded = (key: string) => (expandedKeys || internalExpanded).includes(key);
  const isSelected = (key: string) => (selectedKeys || internalSelected).includes(key);

  const toggleExpand = (key: string) => {
    const current = expandedKeys || internalExpanded;
    const next = current.includes(key) ? current.filter((k) => k !== key) : [...current, key];
    if (expandedKeys) onExpand?.(next);
    else setInternalExpanded(next), onExpand?.(next);
  };

  const toggleSelect = (key: string) => {
    if (!selectable) return;
    const current = selectedKeys || internalSelected;
    const next = current.includes(key) ? current.filter((k) => k !== key) : [key];
    if (selectedKeys) onSelect?.(next);
    else setInternalSelected(next), onSelect?.(next);
  };

  const flatKeys = useMemo(() => new Set<string>(), []);

  const renderNodes = (nodes: TreeNode[], level = 1) => {
    return (
      <ul role={level === 1 ? 'tree' : 'group'} className={styles.treeGroup}>
        {nodes.map((node) => {
          flatKeys.add(node.key);
          const expanded = !!node.children && isExpanded(node.key);
          const selected = isSelected(node.key);

          return (
            <li
              key={node.key}
              role="treeitem"
              aria-expanded={!!node.children ? expanded : undefined}
              aria-selected={selectable ? selected : undefined}
              aria-level={level}
              className={`${styles.treeItem} ${selected ? styles.selected : ''} ${node.disabled ? styles.disabled : ''}`}
            >
              <div className={styles.itemRow}>
                {node.children && (
                  <button
                    type="button"
                    className={styles.expandBtn}
                    aria-label={expanded ? '折叠' : '展开'}
                    onClick={() => !node.disabled && toggleExpand(node.key)}
                    disabled={node.disabled}
                  >
                    {expanded ? '−' : '+'}
                  </button>
                )}
                <button
                  type="button"
                  className={styles.titleBtn}
                  onClick={() => !node.disabled && toggleSelect(node.key)}
                  disabled={node.disabled}
                >
                  {node.title}
                </button>
              </div>
              {node.children && expanded && renderNodes(node.children, level + 1)}
            </li>
          );
        })}
      </ul>
    );
  };

  return (
    <div className={[styles.tree, className].filter(Boolean).join(' ')} style={style}>
      {renderNodes(data)}
    </div>
  );
};