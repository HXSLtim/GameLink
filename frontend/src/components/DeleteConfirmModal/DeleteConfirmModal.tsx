import React from 'react';
import { Modal } from '../Modal';

export interface DeleteConfirmModalProps {
  visible: boolean;
  title?: string;
  content: string | React.ReactNode;
  onConfirm: () => void | Promise<void>;
  onCancel: () => void;
  loading?: boolean;
}

/**
 * 通用删除确认Modal组件
 */
export const DeleteConfirmModal: React.FC<DeleteConfirmModalProps> = ({
  visible,
  title = '确认删除',
  content,
  onConfirm,
  onCancel,
  loading = false,
}) => {
  return (
    <Modal
      visible={visible}
      title={title}
      onClose={onCancel}
      onOk={onConfirm}
      onCancel={onCancel}
      okText={loading ? '删除中...' : '确定删除'}
      cancelText="取消"
      width={400}
    >
      {typeof content === 'string' ? <p>{content}</p> : content}
    </Modal>
  );
};
