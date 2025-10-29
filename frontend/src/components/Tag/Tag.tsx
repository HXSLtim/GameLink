import React from 'react';
import styles from './Tag.module.less';

export type TagColor =
  | 'default'
  | 'success'
  | 'warning'
  | 'error'
  | 'info'
  | 'pending'
  | 'processing'
  | 'green'
  | 'blue'
  | 'red'
  | 'orange'
  | 'purple'
  | 'cyan'
  | 'magenta'
  | 'yellow'
  | 'lime'
  | 'gold'
  | 'volcano'
  | 'geekblue';

export interface TagProps {
  children: React.ReactNode;
  color?: TagColor;
  bordered?: boolean;
  closable?: boolean;
  onClose?: () => void;
  className?: string;
  style?: React.CSSProperties;
}

export const Tag: React.FC<TagProps> = ({
  children,
  color = 'default',
  bordered = true,
  closable = false,
  onClose,
  className,
  style,
}) => {
  const handleClose = (e: React.MouseEvent) => {
    e.stopPropagation();
    onClose?.();
  };

  return (
    <span
      className={`${styles.tag} ${styles[color]} ${bordered ? styles.bordered : ''} ${className || ''}`}
      style={style}
    >
      <span className={styles.text}>{children}</span>
      {closable && (
        <button className={styles.closeButton} onClick={handleClose} aria-label="关闭">
          <CloseIcon />
        </button>
      )}
    </span>
  );
};

const CloseIcon = () => (
  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path strokeLinecap="square" strokeLinejoin="miter" strokeWidth="3" d="M6 6l12 12M6 18L18 6" />
  </svg>
);
