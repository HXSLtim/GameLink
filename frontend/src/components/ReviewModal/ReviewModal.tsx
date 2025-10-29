import React, { useState } from 'react';
import { Modal, Button, Form, FormItem } from '../index';
import { Select } from '../Select/Select';
import styles from './ReviewModal.module.less';

export interface ReviewFormData {
  approved: boolean; // true=通过, false=拒绝
  reason: string;
}

export interface ReviewModalProps {
  visible: boolean;
  orderNo: string;
  onClose: () => void;
  onSubmit: (data: ReviewFormData) => Promise<void>;
}

export const ReviewModal: React.FC<ReviewModalProps> = ({
  visible,
  orderNo,
  onClose,
  onSubmit,
}) => {
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState<ReviewFormData>({
    approved: true,
    reason: '',
  });
  const [errors, setErrors] = useState<Partial<Record<keyof ReviewFormData, string>>>({});

  const validateForm = (): boolean => {
    const newErrors: Partial<Record<keyof ReviewFormData, string>> = {};

    if (!formData.approved && !formData.reason) {
      newErrors.reason = '拒绝时必须填写原因';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async () => {
    if (!validateForm()) {
      return;
    }

    setLoading(true);
    try {
      await onSubmit(formData);
      // 重置表单
      setFormData({
        approved: true,
        reason: '',
      });
      setErrors({});
      onClose();
    } catch (error) {
      console.error('审核失败:', error);
      // 这里可以显示错误提示
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    // 重置表单
    setFormData({
      approved: true,
      reason: '',
    });
    setErrors({});
    onClose();
  };

  return (
    <Modal
      visible={visible}
      title="订单审核"
      onClose={handleCancel}
      width={600}
      footer={
        <div className={styles.footer}>
          <Button variant="outlined" onClick={handleCancel} disabled={loading}>
            取消
          </Button>
          <Button variant="primary" onClick={handleSubmit} disabled={loading}>
            {loading ? '提交中...' : '提交审核'}
          </Button>
        </div>
      }
    >
      <div className={styles.modalContent}>
        <div className={styles.orderInfo}>
          <span className={styles.label}>订单号:</span>
          <span className={styles.value}>{orderNo}</span>
        </div>

        <Form>
          <FormItem label="审核结果" required>
            <Select
              value={formData.approved ? 'approved' : 'rejected'}
              options={[
                { label: '✅ 审核通过', value: 'approved' },
                { label: '❌ 审核拒绝', value: 'rejected' },
              ]}
              onChange={(value) => {
                setFormData({ ...formData, approved: value === 'approved' });
              }}
            />
          </FormItem>

          {!formData.approved && (
            <FormItem label="拒绝原因" required error={errors.reason}>
              <Select
                value={formData.reason}
                options={[
                  { label: '服务未完成', value: '服务未完成' },
                  { label: '服务质量不达标', value: '服务质量不达标' },
                  { label: '用户投诉', value: '用户投诉' },
                  { label: '违反平台规则', value: '违反平台规则' },
                  { label: '证据不足', value: '证据不足' },
                  { label: '其他原因', value: '其他原因' },
                ]}
                onChange={(value) => {
                  setFormData({ ...formData, reason: value as string });
                  setErrors({ ...errors, reason: undefined });
                }}
                placeholder="请选择拒绝原因"
              />
            </FormItem>
          )}

          <FormItem label="备注说明">
            <textarea
              className={styles.textarea}
              value={formData.reason}
              onChange={(e) => setFormData({ ...formData, reason: e.target.value })}
              placeholder={
                formData.approved ? '请填写审核通过的备注（选填）' : '请详细说明拒绝的具体原因'
              }
              rows={4}
            />
          </FormItem>
        </Form>

        <div className={styles.tips}>
          <div className={styles.tipsTitle}>📌 审核提示</div>
          <ul className={styles.tipsList}>
            {formData.approved ? (
              <>
                <li>审核通过后，订单将自动完成</li>
                <li>陪玩师将收到相应的报酬</li>
                <li>用户可以对本次服务进行评价</li>
              </>
            ) : (
              <>
                <li>审核拒绝后，订单将被标记为审核失败</li>
                <li>需要与陪玩师和用户进行沟通协调</li>
                <li>请务必填写清晰的拒绝原因</li>
                <li>可能需要进行退款处理</li>
              </>
            )}
          </ul>
        </div>
      </div>
    </Modal>
  );
};
