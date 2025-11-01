import React, { useState, useEffect } from 'react';
import { Modal, Input, Select, Button, Checkbox } from '../../components';
import type { Announcement, CreateAnnouncementRequest } from '../../types/settings';
import styles from './SettingsDashboard.module.less';

interface AnnouncementFormModalProps {
  visible: boolean;
  announcement: Announcement | null;
  onClose: () => void;
  onSuccess: (data: CreateAnnouncementRequest) => void;
}

/**
 * 公告表单模态框
 */
export const AnnouncementFormModal: React.FC<AnnouncementFormModalProps> = ({
  visible,
  announcement,
  onClose,
  onSuccess,
}) => {
  const [formData, setFormData] = useState<CreateAnnouncementRequest>({
    title: '',
    content: '',
    type: 'info',
    enabled: true,
  });
  const [error, setError] = useState<string | null>(null);

  const isEditMode = !!announcement;

  // 初始化表单数据
  useEffect(() => {
    if (announcement) {
      setFormData({
        title: announcement.title,
        content: announcement.content,
        type: announcement.type,
        startDate: announcement.startDate,
        endDate: announcement.endDate,
        enabled: announcement.enabled,
      });
    } else {
      setFormData({
        title: '',
        content: '',
        type: 'info',
        enabled: true,
      });
    }
  }, [announcement]);

  // 表单提交
  const handleSubmit = () => {
    // 表单验证
    if (!formData.title.trim() || !formData.content.trim()) {
      setError('标题和内容不能为空');
      return;
    }

    setError(null);
    onSuccess(formData);
  };

  return (
    <Modal
      visible={visible}
      title={isEditMode ? '编辑公告' : '创建公告'}
      onCancel={onClose}
      width={600}
      footer={
        <div className={styles.modalFooter}>
          <Button onClick={onClose}>取消</Button>
          <Button type="primary" onClick={handleSubmit}>
            {isEditMode ? '保存' : '创建'}
          </Button>
        </div>
      }
    >
      <div className={styles.formContent}>
        {error && <div className={styles.errorMessage}>{error}</div>}

        <div className={styles.formItem}>
          <label className={styles.formLabel}>
            公告标题 <span className={styles.required}>*</span>
          </label>
          <Input
            value={formData.title}
            onChange={(e) => setFormData({ ...formData, title: e.target.value })}
            placeholder="请输入公告标题"
          />
        </div>

        <div className={styles.formItem}>
          <label className={styles.formLabel}>
            公告内容 <span className={styles.required}>*</span>
          </label>
          <Input
            value={formData.content}
            onChange={(e) => setFormData({ ...formData, content: e.target.value })}
            placeholder="请输入公告内容"
          />
        </div>

        <div className={styles.formItem}>
          <label className={styles.formLabel}>
            公告类型 <span className={styles.required}>*</span>
          </label>
          <Select
            value={formData.type}
            onChange={(value) => setFormData({ ...formData, type: value as any })}
            options={[
              { label: '信息', value: 'info' },
              { label: '警告', value: 'warning' },
              { label: '错误', value: 'error' },
              { label: '成功', value: 'success' },
            ]}
          />
        </div>

        <div className={styles.formItem}>
          <label className={styles.formLabel}>开始日期</label>
          <Input
            type="date"
            value={formData.startDate || ''}
            onChange={(e) => setFormData({ ...formData, startDate: e.target.value })}
          />
        </div>

        <div className={styles.formItem}>
          <label className={styles.formLabel}>结束日期</label>
          <Input
            type="date"
            value={formData.endDate || ''}
            onChange={(e) => setFormData({ ...formData, endDate: e.target.value })}
          />
        </div>

        <div className={styles.formItem}>
          <Checkbox
            checked={formData.enabled}
            onChange={(e) => setFormData({ ...formData, enabled: e.target.checked })}
          >
            立即启用
          </Checkbox>
        </div>
      </div>
    </Modal>
  );
};



