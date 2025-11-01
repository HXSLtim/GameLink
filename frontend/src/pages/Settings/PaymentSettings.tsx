import React, { useState } from 'react';
import { Card, Button, Input, Checkbox, message } from '../../components';
import styles from './SettingsDashboard.module.less';

/**
 * 支付设置组件
 */
export const PaymentSettings: React.FC = () => {
  const [formData, setFormData] = useState({
    alipayEnabled: true,
    alipayAppId: '',
    wechatEnabled: true,
    wechatAppId: '',
    wechatMchId: '',
    balanceEnabled: true,
    minRechargeAmount: 100,
    maxRechargeAmount: 100000,
    serviceFeeRate: 5,
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
        <div className={styles.sectionTitle}>支付方式</div>

        <div className={styles.formItem}>
          <Checkbox
            checked={formData.alipayEnabled}
            onChange={(e) => handleChange('alipayEnabled', e.target.checked)}
          >
            启用支付宝支付
          </Checkbox>
        </div>

        {formData.alipayEnabled && (
          <div className={styles.formItem}>
            <label className={styles.formLabel}>支付宝 AppID</label>
            <Input
              value={formData.alipayAppId}
              onChange={(e) => handleChange('alipayAppId', e.target.value)}
              placeholder="请输入支付宝 AppID"
            />
          </div>
        )}

        <div className={styles.formItem}>
          <Checkbox
            checked={formData.wechatEnabled}
            onChange={(e) => handleChange('wechatEnabled', e.target.checked)}
          >
            启用微信支付
          </Checkbox>
        </div>

        {formData.wechatEnabled && (
          <>
            <div className={styles.formItem}>
              <label className={styles.formLabel}>微信 AppID</label>
              <Input
                value={formData.wechatAppId}
                onChange={(e) => handleChange('wechatAppId', e.target.value)}
                placeholder="请输入微信 AppID"
              />
            </div>

            <div className={styles.formItem}>
              <label className={styles.formLabel}>微信商户号</label>
              <Input
                value={formData.wechatMchId}
                onChange={(e) => handleChange('wechatMchId', e.target.value)}
                placeholder="请输入微信商户号"
              />
            </div>
          </>
        )}

        <div className={styles.formItem}>
          <Checkbox
            checked={formData.balanceEnabled}
            onChange={(e) => handleChange('balanceEnabled', e.target.checked)}
          >
            启用余额支付
          </Checkbox>
        </div>

        <div className={styles.divider} />

        <div className={styles.sectionTitle}>支付金额限制</div>

        <div className={styles.formRow}>
          <div className={styles.formItem}>
            <label className={styles.formLabel}>最小充值金额（元）</label>
            <Input
              type="number"
              value={formData.minRechargeAmount}
              onChange={(e) => handleChange('minRechargeAmount', Number(e.target.value))}
              placeholder="100"
            />
          </div>

          <div className={styles.formItem}>
            <label className={styles.formLabel}>最大充值金额（元）</label>
            <Input
              type="number"
              value={formData.maxRechargeAmount}
              onChange={(e) => handleChange('maxRechargeAmount', Number(e.target.value))}
              placeholder="100000"
            />
          </div>
        </div>

        <div className={styles.divider} />

        <div className={styles.sectionTitle}>平台费率</div>

        <div className={styles.formItem}>
          <label className={styles.formLabel}>服务费率（%）</label>
          <Input
            type="number"
            value={formData.serviceFeeRate}
            onChange={(e) => handleChange('serviceFeeRate', Number(e.target.value))}
            placeholder="5"
            step="0.1"
            min="0"
            max="100"
          />
          <div className={styles.formHint}>
            平台从每笔交易中收取的服务费百分比，当前设置为 {formData.serviceFeeRate}%
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

