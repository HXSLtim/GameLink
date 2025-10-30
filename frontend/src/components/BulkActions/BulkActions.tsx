import React from 'react';
import { Button } from '../Button';
import styles from './BulkActions.module.less';

export interface BulkAction {
  key: string;
  label: string;
  variant?: 'primary' | 'outlined' | 'text';
  danger?: boolean;
  disabled?: boolean;
  onClick: () => void;
}

interface BulkActionsProps {
  selectedCount: number;
  totalCount: number;
  actions: BulkAction[];
  onClearSelection?: () => void;
}

export const BulkActions: React.FC<BulkActionsProps> = ({
  selectedCount,
  totalCount,
  actions,
  onClearSelection,
}) => {
  if (selectedCount === 0) {
    return null;
  }

  return (
    <div className={styles.container}>
      <div className={styles.info}>
        <span className={styles.selectedCount}>
          已选择 <strong>{selectedCount}</strong> 项
        </span>
        {totalCount > 0 && (
          <span className={styles.totalCount}>
            （共 {totalCount} 项）
          </span>
        )}
      </div>
      
      <div className={styles.actions}>
        {actions.map((action) => (
          <Button
            key={action.key}
            variant={action.variant || 'outlined'}
            onClick={action.onClick}
            disabled={action.disabled}
            className={action.danger ? styles.dangerButton : ''}
          >
            {action.label}
          </Button>
        ))}
        
        {onClearSelection && (
          <Button variant="text" onClick={onClearSelection}>
            取消选择
          </Button>
        )}
      </div>
    </div>
  );
};


