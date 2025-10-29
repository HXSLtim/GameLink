import React, { useState, useEffect } from 'react';
import { Modal, Input, Select } from '../../components';
import type { Order, UpdateOrderRequest, OrderStatus } from '../../types/order';
import styles from './OrderFormModal.module.less';

interface OrderFormModalProps {
  visible: boolean;
  order: Order | null;
  onClose: () => void;
  onSubmit: (data: UpdateOrderRequest) => Promise<void>;
}

export const OrderFormModal: React.FC<OrderFormModalProps> = ({
  visible,
  order,
  onClose,
  onSubmit,
}) => {
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState<UpdateOrderRequest>({
    currency: 'CNY',
    priceCents: 0,
    status: 'pending' as OrderStatus,
    scheduledStart: '',
    scheduledEnd: '',
    cancelReason: '',
  });

  useEffect(() => {
    if (order) {
      setFormData({
        currency: order.currency || 'CNY',
        priceCents: order.priceCents,
        status: order.status as OrderStatus,
        scheduledStart: order.scheduledStart || '',
        scheduledEnd: order.scheduledEnd || '',
        cancelReason: order.cancelReason || '',
      });
    }
  }, [order]);

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
      title="编辑订单"
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
            订单状态 <span className={styles.required}>*</span>
          </label>
          <Select
            value={formData.status}
            onChange={(value) => setFormData({ ...formData, status: value as OrderStatus })}
            options={[
              { label: '待处理', value: 'pending' },
              { label: '已确认', value: 'confirmed' },
              { label: '进行中', value: 'in_progress' },
              { label: '已完成', value: 'completed' },
              { label: '已取消', value: 'canceled' },
              { label: '已退款', value: 'refunded' },
            ]}
          />
        </div>

        <div className={styles.formItem}>
          <label className={styles.label}>
            金额（分） <span className={styles.required}>*</span>
          </label>
          <Input
            value={formData.priceCents.toString()}
            onChange={(e) =>
              setFormData({ ...formData, priceCents: parseInt(e.target.value) || 0 })
            }
            placeholder="请输入金额（单位：分）"
            type="number"
          />
          <div className={styles.hint}>例如：10000 表示 100.00 元</div>
        </div>

        <div className={styles.formItem}>
          <label className={styles.label}>货币</label>
          <Select
            value={formData.currency}
            onChange={(value) => setFormData({ ...formData, currency: value as string })}
            options={[
              { label: '人民币（CNY）', value: 'CNY' },
              { label: '美元（USD）', value: 'USD' },
            ]}
          />
        </div>

        <div className={styles.formItem}>
          <label className={styles.label}>预约开始时间</label>
          <Input
            value={formData.scheduledStart || ''}
            onChange={(e) => setFormData({ ...formData, scheduledStart: e.target.value })}
            placeholder="YYYY-MM-DD HH:mm:ss"
            type="datetime-local"
          />
        </div>

        <div className={styles.formItem}>
          <label className={styles.label}>预约结束时间</label>
          <Input
            value={formData.scheduledEnd || ''}
            onChange={(e) => setFormData({ ...formData, scheduledEnd: e.target.value })}
            placeholder="YYYY-MM-DD HH:mm:ss"
            type="datetime-local"
          />
        </div>

        {formData.status === 'canceled' && (
          <div className={styles.formItem}>
            <label className={styles.label}>取消原因</label>
            <textarea
              className={styles.textarea}
              value={formData.cancelReason || ''}
              onChange={(e) => setFormData({ ...formData, cancelReason: e.target.value })}
              placeholder="请输入取消原因"
              rows={3}
            />
          </div>
        )}
      </div>
    </Modal>
  );
};
