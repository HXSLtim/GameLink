import React, { useEffect, useRef } from 'react';
import { createPortal } from 'react-dom';
import styles from './Modal.module.less';

export interface ModalProps {
  visible: boolean;
  title?: string;
  children: React.ReactNode;
  onClose?: () => void;
  onOk?: () => void;
  onCancel?: () => void;
  okText?: string;
  cancelText?: string;
  width?: number | string;
  footer?: React.ReactNode | null;
  closable?: boolean;
  maskClosable?: boolean;
  className?: string;
}

export const Modal: React.FC<ModalProps> = ({
  visible,
  title,
  children,
  onClose,
  onOk,
  onCancel,
  okText = '确定',
  cancelText = '取消',
  width = 520,
  footer,
  closable = true,
  maskClosable = true,
  className,
}) => {
  const modalRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleEsc = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && visible && onClose) {
        onClose();
      }
    };

    if (visible) {
      document.addEventListener('keydown', handleEsc);
      document.body.style.overflow = 'hidden';
    }

    return () => {
      document.removeEventListener('keydown', handleEsc);
      document.body.style.overflow = '';
    };
  }, [visible, onClose]);

  const handleMaskClick = (e: React.MouseEvent) => {
    if (maskClosable && e.target === e.currentTarget && onClose) {
      onClose();
    }
  };

  const handleCancel = () => {
    onCancel?.();
    onClose?.();
  };

  const handleOk = () => {
    onOk?.();
  };

  if (!visible) return null;

  const modalContent = (
    <div className={styles.modalWrapper} onClick={handleMaskClick}>
      <div className={styles.modalMask} />
      <div className={styles.modalContainer}>
        <div ref={modalRef} className={`${styles.modal} ${className || ''}`} style={{ width }}>
          {closable && (
            <button className={styles.closeButton} onClick={onClose} aria-label="关闭">
              <CloseIcon />
            </button>
          )}

          {title && (
            <div className={styles.modalHeader}>
              <h3 className={styles.modalTitle}>{title}</h3>
            </div>
          )}

          <div className={styles.modalBody}>{children}</div>

          {footer !== null && (
            <div className={styles.modalFooter}>
              {footer || (
                <>
                  <button className={styles.cancelButton} onClick={handleCancel}>
                    {cancelText}
                  </button>
                  <button className={styles.okButton} onClick={handleOk}>
                    {okText}
                  </button>
                </>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );

  return createPortal(modalContent, document.body);
};

const CloseIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path strokeLinecap="square" strokeLinejoin="miter" strokeWidth="2" d="M6 6l12 12M6 18L18 6" />
  </svg>
);
