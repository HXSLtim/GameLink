import React, { useState, useEffect } from 'react';
import { Modal, Select } from '../../components';
import type { Review, UpdateReviewRequest } from '../../types/review';
import styles from './ReviewFormModal.module.less';

interface ReviewFormModalProps {
  visible: boolean;
  review: Review | null;
  onClose: () => void;
  onSubmit: (data: UpdateReviewRequest) => Promise<void>;
}

export const ReviewFormModal: React.FC<ReviewFormModalProps> = ({
  visible,
  review,
  onClose,
  onSubmit,
}) => {
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState<UpdateReviewRequest>({
    rating: 5,
    comment: '',
  });

  useEffect(() => {
    if (review) {
      setFormData({
        rating: review.rating,
        comment: review.comment || '',
      });
    }
  }, [review]);

  const handleSubmit = async () => {
    setLoading(true);
    try {
      await onSubmit(formData);
      onClose();
    } catch (err) {
      console.error('提交失败:', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal
      visible={visible}
      title="编辑评价"
      onClose={onClose}
      onOk={handleSubmit}
      onCancel={onClose}
      okText={loading ? '提交中...' : '确定'}
      cancelText="取消"
      width={600}
    >
      <div className={styles.form}>
        <div className={styles.formItem}>
          <label className={styles.label}>
            评分 <span className={styles.required}>*</span>
          </label>
          <Select
            value={formData.rating?.toString() || '5'}
            onChange={(value) => setFormData({ ...formData, rating: parseInt(value as string) })}
            options={[
              { label: '⭐ 1星 - 非常差', value: '1' },
              { label: '⭐⭐ 2星 - 较差', value: '2' },
              { label: '⭐⭐⭐ 3星 - 一般', value: '3' },
              { label: '⭐⭐⭐⭐ 4星 - 满意', value: '4' },
              { label: '⭐⭐⭐⭐⭐ 5星 - 非常满意', value: '5' },
            ]}
          />
        </div>

        <div className={styles.formItem}>
          <label className={styles.label}>评价内容</label>
          <textarea
            className={styles.textarea}
            value={formData.comment || ''}
            onChange={(e) => setFormData({ ...formData, comment: e.target.value })}
            placeholder="请输入评价内容"
            rows={6}
          />
        </div>
      </div>
    </Modal>
  );
};
