import React, { useState, useEffect } from 'react';
import { Modal, Input, Button } from '../../components';
import { roleApi } from '../../services/api/rbac';
import type { Role, CreateRoleRequest, UpdateRoleRequest } from '../../types/rbac';
import styles from './RoleFormModal.module.less';

interface RoleFormModalProps {
  visible: boolean;
  role: Role | null;
  onClose: () => void;
  onSuccess: () => void;
}

/**
 * 角色表单模态框
 */
export const RoleFormModal: React.FC<RoleFormModalProps> = ({
  visible,
  role,
  onClose,
  onSuccess,
}) => {
  const [formData, setFormData] = useState({
    slug: '',
    name: '',
    description: '',
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const isEditMode = !!role;

  // 初始化表单数据
  useEffect(() => {
    if (role) {
      setFormData({
        slug: role.slug,
        name: role.name,
        description: role.description || '',
      });
    } else {
      setFormData({
        slug: '',
        name: '',
        description: '',
      });
    }
  }, [role]);

  // 表单提交
  const handleSubmit = async () => {
    // 表单验证
    if (!formData.slug.trim() || !formData.name.trim()) {
      setError('角色标识和名称不能为空');
      return;
    }

    try {
      setIsSubmitting(true);
      setError(null);

      if (isEditMode) {
        // 更新角色
        const updateData: UpdateRoleRequest = {
          name: formData.name,
          description: formData.description || undefined,
        };
        await roleApi.update(role.id, updateData);
      } else {
        // 创建角色
        const createData: CreateRoleRequest = {
          slug: formData.slug,
          name: formData.name,
          description: formData.description || undefined,
        };
        await roleApi.create(createData);
      }

      onSuccess();
    } catch (err: any) {
      console.error('保存角色失败:', err);
      setError(err.message || '保存失败，请重试');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Modal
      visible={visible}
      title={isEditMode ? '编辑角色' : '创建角色'}
      onCancel={onClose}
      footer={
        <div className={styles.modalFooter}>
          <Button variant="secondary" onClick={onClose} disabled={isSubmitting}>
            取消
          </Button>
          <Button variant="primary" onClick={handleSubmit} loading={isSubmitting}>
            {isEditMode ? '保存' : '创建'}
          </Button>
        </div>
      }
    >
      <div className={styles.formContent}>
        {error && <div className={styles.errorMessage}>{error}</div>}

        <div className={styles.formItem}>
          <label className={styles.formLabel}>
            角色标识 <span className={styles.required}>*</span>
          </label>
          <Input
            value={formData.slug}
            onChange={(e) => setFormData({ ...formData, slug: e.target.value })}
            placeholder="如：custom_role"
            disabled={isEditMode}
          />
          {isEditMode && (
            <div className={styles.formHint}>系统角色的标识不可修改</div>
          )}
        </div>

        <div className={styles.formItem}>
          <label className={styles.formLabel}>
            角色名称 <span className={styles.required}>*</span>
          </label>
          <Input
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            placeholder="如：自定义角色"
          />
        </div>

        <div className={styles.formItem}>
          <label className={styles.formLabel}>角色描述</label>
          <Input
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            placeholder="描述该角色的用途"
          />
        </div>
      </div>
    </Modal>
  );
};

