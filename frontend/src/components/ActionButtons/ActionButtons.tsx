import React from 'react';
import { Button } from '../Button';
import styles from './ActionButtons.module.less';

export interface ActionButtonsProps {
  onView?: () => void;
  onEdit?: () => void;
  onDelete?: () => void;
  viewText?: string;
  editText?: string;
  deleteText?: string;
  className?: string;
}

/**
 * 通用操作按钮组件
 * 包含查看、编辑、删除按钮
 */
export const ActionButtons: React.FC<ActionButtonsProps> = ({
  onView,
  onEdit,
  onDelete,
  viewText = '详情',
  editText = '编辑',
  deleteText = '删除',
  className = '',
}) => {
  return (
    <div className={`${styles.actions} ${className}`}>
      {onView && (
        <Button variant="text" onClick={onView} className={styles.actionButton}>
          {viewText}
        </Button>
      )}
      {onEdit && (
        <Button variant="text" onClick={onEdit} className={styles.actionButton}>
          {editText}
        </Button>
      )}
      {onDelete && (
        <Button variant="text" onClick={onDelete} className={styles.deleteButton}>
          {deleteText}
        </Button>
      )}
    </div>
  );
};
