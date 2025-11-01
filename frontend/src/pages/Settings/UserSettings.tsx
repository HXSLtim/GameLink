import React, { useState } from 'react';
import { Card, Button, Input, Select, Checkbox, message } from '../../components';
import styles from './SettingsDashboard.module.less';

/**
 * 用户设置组件
 */
export const UserSettings: React.FC = () => {
  const [formData, setFormData] = useState({
    registrationEnabled: true,
    emailVerificationRequired: false,
    phoneVerificationRequired: true,
    minPasswordLength: 6,
    defaultUserRole: 'user',
    sessionTimeout: 1440,
    playerVerificationRequired: true,
    playerMinAge: 18,
    playerMinHourlyRate: 1000,
    playerMaxHourlyRate: 50000,
    playerCommissionRate: 10,
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
        <div className={styles.sectionTitle}>用户注册</div>

        <div className={styles.formItem}>
          <Checkbox
            checked={formData.registrationEnabled}
            onChange={(e) => handleChange('registrationEnabled', e.target.checked)}
          >
            允许新用户注册
          </Checkbox>
        </div>

        <div className={styles.formItem}>
          <Checkbox
            checked={formData.emailVerificationRequired}
            onChange={(e) => handleChange('emailVerificationRequired', e.target.checked)}
          >
            要求邮箱验证
          </Checkbox>
        </div>

        <div className={styles.formItem}>
          <Checkbox
            checked={formData.phoneVerificationRequired}
            onChange={(e) => handleChange('phoneVerificationRequired', e.target.checked)}
          >
            要求手机验证
          </Checkbox>
        </div>

        <div className={styles.formItem}>
          <label className={styles.formLabel}>最小密码长度</label>
          <Input
            type="number"
            value={formData.minPasswordLength}
            onChange={(e) => handleChange('minPasswordLength', Number(e.target.value))}
            placeholder="6"
            min="6"
            max="20"
          />
        </div>

        <div className={styles.formItem}>
          <label className={styles.formLabel}>默认用户角色</label>
          <Select
            value={formData.defaultUserRole}
            onChange={(value) => handleChange('defaultUserRole', value)}
            options={[
              { label: '普通用户', value: 'user' },
              { label: '陪玩师', value: 'player' },
            ]}
          />
        </div>

        <div className={styles.divider} />

        <div className={styles.sectionTitle}>会话管理</div>

        <div className={styles.formItem}>
          <label className={styles.formLabel}>会话超时时间（分钟）</label>
          <Input
            type="number"
            value={formData.sessionTimeout}
            onChange={(e) => handleChange('sessionTimeout', Number(e.target.value))}
            placeholder="1440"
          />
          <div className={styles.formHint}>
            用户登录后超过此时间未活动将自动退出，默认 1440 分钟（24小时）
          </div>
        </div>

        <div className={styles.divider} />

        <div className={styles.sectionTitle}>陪玩师设置</div>

        <div className={styles.formItem}>
          <Checkbox
            checked={formData.playerVerificationRequired}
            onChange={(e) => handleChange('playerVerificationRequired', e.target.checked)}
          >
            陪玩师需要实名认证
          </Checkbox>
        </div>

        <div className={styles.formItem}>
          <label className={styles.formLabel}>最低年龄要求</label>
          <Input
            type="number"
            value={formData.playerMinAge}
            onChange={(e) => handleChange('playerMinAge', Number(e.target.value))}
            placeholder="18"
            min="16"
            max="60"
          />
        </div>

        <div className={styles.formRow}>
          <div className={styles.formItem}>
            <label className={styles.formLabel}>最低时薪（分）</label>
            <Input
              type="number"
              value={formData.playerMinHourlyRate}
              onChange={(e) => handleChange('playerMinHourlyRate', Number(e.target.value))}
              placeholder="1000"
            />
            <div className={styles.formHint}>
              约 {(formData.playerMinHourlyRate / 100).toFixed(2)} 元/小时
            </div>
          </div>

          <div className={styles.formItem}>
            <label className={styles.formLabel}>最高时薪（分）</label>
            <Input
              type="number"
              value={formData.playerMaxHourlyRate}
              onChange={(e) => handleChange('playerMaxHourlyRate', Number(e.target.value))}
              placeholder="50000"
            />
            <div className={styles.formHint}>
              约 {(formData.playerMaxHourlyRate / 100).toFixed(2)} 元/小时
            </div>
          </div>
        </div>

        <div className={styles.formItem}>
          <label className={styles.formLabel}>平台抽成比例（%）</label>
          <Input
            type="number"
            value={formData.playerCommissionRate}
            onChange={(e) => handleChange('playerCommissionRate', Number(e.target.value))}
            placeholder="10"
            step="0.1"
            min="0"
            max="50"
          />
          <div className={styles.formHint}>
            平台从陪玩师收入中抽取的比例，当前设置为 {formData.playerCommissionRate}%
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

