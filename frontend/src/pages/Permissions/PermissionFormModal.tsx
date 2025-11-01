import React, { useState, useEffect } from 'react';
import { Modal, Input, Select, Button } from '../../components';
import { permissionApi } from '../../services/api/rbac';
import type { Permission, CreatePermissionRequest, UpdatePermissionRequest } from '../../types/rbac';
import { HTTPMethod } from '../../types/rbac';
import { HTTP_METHOD_OPTIONS } from '../../utils/selectOptions';
import styles from './PermissionFormModal.module.less';

interface PermissionFormModalProps {
  visible: boolean;
  permission: Permission | null;
  onClose: () => void;
  onSuccess: () => void;
}

/**
 * 权限表单模态框
 */
export const PermissionFormModal: React.FC<PermissionFormModalProps> = ({
  visible,
  permission,
  onClose,
  onSuccess,
}) => {
  const [formData, setFormData] = useState<{
    method: HTTPMethod;
    path: string;
    code: string;
    group: string;
    description: string;
  }>({
    method: HTTPMethod.GET,
    path: '',
    code: '',
    group: '',
    description: '',
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const isEditMode = !!permission;

  // 初始化表单数据
  useEffect(() => {
    if (permission) {
      setFormData({
        method: permission.method,
        path: permission.path,
        code: permission.code,
        group: permission.group,
        description: permission.description || '',
      });
    } else {
      setFormData({
        method: HTTPMethod.GET,
        path: '',
        code: '',
        group: '',
        description: '',
      });
    }
  }, [permission]);

  // 表单提交
  const handleSubmit = async () => {
    // 表单验证
    if (!formData.method || !formData.path.trim() || !formData.code.trim() || !formData.group.trim()) {
      setError('HTTP方法、API路径、权限代码和分组不能为空');
      return;
    }

    try {
      setIsSubmitting(true);
      setError(null);

      if (isEditMode) {
        // 更新权限
        const updateData: UpdatePermissionRequest = {
          code: formData.code,
          group: formData.group,
          description: formData.description || undefined,
        };
        await permissionApi.update(permission.id, updateData);
      } else {
        // 创建权限
        const createData: CreatePermissionRequest = {
          method: formData.method,
          path: formData.path,
          code: formData.code,
          group: formData.group,
          description: formData.description || undefined,
        };
        await permissionApi.create(createData);
      }

      onSuccess();
    } catch (err: any) {
      console.error('保存权限失败:', err);
      setError(err.message || '保存失败，请重试');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Modal
      visible={visible}
      title={isEditMode ? '编辑权限' : '创建权限'}
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
            HTTP方法 <span className={styles.required}>*</span>
          </label>
          <Select
            value={formData.method}
            onChange={(value) => setFormData({ ...formData, method: value as HTTPMethod })}
            options={HTTP_METHOD_OPTIONS.filter(opt => opt.value !== '')}
            disabled={isEditMode}
          />
          {isEditMode && (
            <div className={styles.formHint}>HTTP方法不可修改</div>
          )}
        </div>

        <div className={styles.formItem}>
          <label className={styles.formLabel}>
            API路径 <span className={styles.required}>*</span>
          </label>
          <Input
            value={formData.path}
            onChange={(e) => setFormData({ ...formData, path: e.target.value })}
            placeholder="如：/api/v1/admin/users"
            disabled={isEditMode}
          />
          {isEditMode && (
            <div className={styles.formHint}>API路径不可修改</div>
          )}
        </div>

        <div className={styles.formItem}>
          <label className={styles.formLabel}>
            权限代码 <span className={styles.required}>*</span>
          </label>
          <Input
            value={formData.code}
            onChange={(e) => setFormData({ ...formData, code: e.target.value })}
            placeholder="如：admin.users.read"
          />
        </div>

        <div className={styles.formItem}>
          <label className={styles.formLabel}>
            权限分组 <span className={styles.required}>*</span>
          </label>
          <Input
            value={formData.group}
            onChange={(e) => setFormData({ ...formData, group: e.target.value })}
            placeholder="如：/admin/users"
          />
        </div>

        <div className={styles.formItem}>
          <label className={styles.formLabel}>权限描述</label>
          <Input
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            placeholder="描述该权限的用途"
          />
        </div>
      </div>
    </Modal>
  );
};

