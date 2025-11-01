import React, { useState } from 'react';
import { Card, Button, Input, Checkbox, message } from '../../components';
import styles from './SettingsDashboard.module.less';

/**
 * 订单设置组件
 */
export const OrderSettings: React.FC = () => {
  const [formData, setFormData] = useState({
    autoConfirmMinutes: 30,
    autoCancelMinutes: 15,
    refundEnabled: true,
    refundDays: 7,
    minOrderAmount: 1000,
    maxOrderAmount: 1000000,
  });

  const [isSaving, setIsSaving] = useState(false);

  const handleChange = (field: string, value: any) => {
    setFormData({ ...formData, [field]: value });
  };

  const handleSave = async () => {
    try {
      setIsSaving(true);
      // TODO: 调用后端API保存设置
      await new Promise((resolve) => setTimeout(resolve, 1000));
      message.success('保存成功');
    } catch {
      message.error('保存失败');
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <Card className={styles.settingsCard}>
      <div className={styles.settingsContent}>
        <div className={styles.sectionTitle}>订单自动处理</div>

        <div className={styles.formItem}>
          <label className={styles.formLabel}>自动确认时间（分钟）</label>
          <Input
            type="number"
            value={formData.autoConfirmMinutes}
            onChange={(e) => handleChange('autoConfirmMinutes', Number(e.target.value))}
            placeholder="30"
          />
          <div className={styles.formHint}>
            订单创建后超过此时间未确认将自动确认，0 表示不自动确认
          </div>
        </div>

        <div className={styles.formItem}>
          <label className={styles.formLabel}>自动取消时间（分钟）</label>
          <Input
            type="number"
            value={formData.autoCancelMinutes}
            onChange={(e) => handleChange('autoCancelMinutes', Number(e.target.value))}
            placeholder="15"
          />
          <div className={styles.formHint}>
            待支付订单超过此时间未支付将自动取消，0 表示不自动取消
          </div>
        </div>

        <div className={styles.divider} />

        <div className={styles.sectionTitle}>退款设置</div>

        <div className={styles.formItem}>
          <Checkbox
            checked={formData.refundEnabled}
            onChange={(e) => handleChange('refundEnabled', e.target.checked)}
          >
            允许用户申请退款
          </Checkbox>
        </div>

        {formData.refundEnabled && (
          <div className={styles.formItem}>
            <label className={styles.formLabel}>退款期限（天）</label>
            <Input
              type="number"
              value={formData.refundDays}
              onChange={(e) => handleChange('refundDays', Number(e.target.value))}
              placeholder="7"
            />
            <div className={styles.formHint}>订单完成后多少天内可以申请退款</div>
          </div>
        )}

        <div className={styles.divider} />

        <div className={styles.sectionTitle}>订单金额限制</div>

        <div className={styles.formRow}>
          <div className={styles.formItem}>
            <label className={styles.formLabel}>最小订单金额（分）</label>
            <Input
              type="number"
              value={formData.minOrderAmount}
              onChange={(e) => handleChange('minOrderAmount', Number(e.target.value))}
              placeholder="1000"
            />
            <div className={styles.formHint}>约 {(formData.minOrderAmount / 100).toFixed(2)} 元</div>
          </div>

          <div className={styles.formItem}>
            <label className={styles.formLabel}>最大订单金额（分）</label>
            <Input
              type="number"
              value={formData.maxOrderAmount}
              onChange={(e) => handleChange('maxOrderAmount', Number(e.target.value))}
              placeholder="1000000"
            />
            <div className={styles.formHint}>
              约 {(formData.maxOrderAmount / 100).toFixed(2)} 元
            </div>
          </div>
        </div>

        <div className={styles.formActions}>
          <Button type="primary" onClick={handleSave} loading={isSaving}>
            保存设置
          </Button>
        </div>
      </div>
    </Card>
  );
};

