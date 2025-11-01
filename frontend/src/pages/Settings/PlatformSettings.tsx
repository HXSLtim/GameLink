import React, { useState } from 'react';
import { Card, Button, Input, message } from '../../components';
import styles from './SettingsDashboard.module.less';

/**
 * 平台配置组件
 */
export const PlatformSettings: React.FC = () => {
  const [formData, setFormData] = useState({
    siteName: 'GameLink',
    siteUrl: 'https://gamelink.example.com',
    contactEmail: 'contact@gamelink.com',
    contactPhone: '400-123-4567',
    icp: '京ICP备12345678号',
    description: 'GameLink - 专业的游戏陪玩平台',
    keywords: '游戏陪玩,电竞陪练,游戏社交',
  });

  const [isSaving, setIsSaving] = useState(false);

  const handleChange = (field: string, value: string) => {
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
        <div className={styles.formItem}>
          <label className={styles.formLabel}>
            网站名称 <span className={styles.required}>*</span>
          </label>
          <Input
            value={formData.siteName}
            onChange={(e) => handleChange('siteName', e.target.value)}
            placeholder="请输入网站名称"
          />
        </div>

        <div className={styles.formItem}>
          <label className={styles.formLabel}>
            网站地址 <span className={styles.required}>*</span>
          </label>
          <Input
            value={formData.siteUrl}
            onChange={(e) => handleChange('siteUrl', e.target.value)}
            placeholder="https://example.com"
          />
        </div>

        <div className={styles.formItem}>
          <label className={styles.formLabel}>
            联系邮箱 <span className={styles.required}>*</span>
          </label>
          <Input
            value={formData.contactEmail}
            onChange={(e) => handleChange('contactEmail', e.target.value)}
            placeholder="contact@example.com"
          />
        </div>

        <div className={styles.formItem}>
          <label className={styles.formLabel}>联系电话</label>
          <Input
            value={formData.contactPhone}
            onChange={(e) => handleChange('contactPhone', e.target.value)}
            placeholder="400-123-4567"
          />
        </div>

        <div className={styles.formItem}>
          <label className={styles.formLabel}>ICP备案号</label>
          <Input
            value={formData.icp}
            onChange={(e) => handleChange('icp', e.target.value)}
            placeholder="京ICP备12345678号"
          />
        </div>

        <div className={styles.formItem}>
          <label className={styles.formLabel}>网站描述</label>
          <Input
            value={formData.description}
            onChange={(e) => handleChange('description', e.target.value)}
            placeholder="简要描述网站功能"
          />
        </div>

        <div className={styles.formItem}>
          <label className={styles.formLabel}>SEO关键词</label>
          <Input
            value={formData.keywords}
            onChange={(e) => handleChange('keywords', e.target.value)}
            placeholder="关键词1,关键词2,关键词3"
          />
          <div className={styles.formHint}>多个关键词用英文逗号分隔</div>
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

