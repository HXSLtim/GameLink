import React, { useState, useEffect } from 'react';
import { Modal, Input, Select, FormField } from '../../components';
import type { User, CreateUserRequest, UpdateUserRequest, UserRole, UserStatus } from '../../types/user';
import { USER_ROLE_OPTIONS, USER_STATUS_OPTIONS } from '../../utils/selectOptions';

interface UserFormModalProps {
  visible: boolean;
  onClose: () => void;
  onSave: (data: CreateUserRequest | UpdateUserRequest) => Promise<void>;
  initialData: User | null;
  isSubmitting: boolean;
}

export const UserFormModal: React.FC<UserFormModalProps> = ({
  visible,
  onClose,
  onSave,
  initialData,
  isSubmitting,
}) => {
  const [formData, setFormData] = useState<CreateUserRequest | UpdateUserRequest>({
    name: '',
    phone: '',
    email: '',
    password: '',
    role: 'user' as UserRole,
    status: 'active' as UserStatus,
  });

  useEffect(() => {
    if (initialData) {
      setFormData({
        name: initialData.name,
        phone: initialData.phone || '',
        email: initialData.email || '',
        role: initialData.role,
        status: initialData.status,
      });
    } else {
      setFormData({
        name: '',
        phone: '',
        email: '',
        password: '',
        role: 'user' as UserRole,
        status: 'active' as UserStatus,
      });
    }
  }, [initialData, visible]);

  const handleSubmit = async () => {
    if (!formData.name?.trim()) {
      alert('请输入用户名');
      return;
    }

    if (!initialData && !(formData as CreateUserRequest).password) {
      alert('请输入密码');
      return;
    }

    await onSave(formData);
  };

  return (
    <Modal
      visible={visible}
      title={initialData ? '编辑用户' : '新增用户'}
      onClose={onClose}
      onOk={handleSubmit}
      onCancel={onClose}
      okText={isSubmitting ? '提交中...' : '确定'}
      cancelText="取消"
      width={600}
    >
      <FormField label="用户名" required>
        <Input
          value={formData.name}
          onChange={(e) => setFormData({ ...formData, name: e.target.value })}
          placeholder="请输入用户名"
        />
      </FormField>

      <FormField label="手机号">
        <Input
          value={formData.phone || ''}
          onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
          placeholder="请输入手机号"
          type="tel"
        />
      </FormField>

      <FormField label="邮箱">
        <Input
          value={formData.email || ''}
          onChange={(e) => setFormData({ ...formData, email: e.target.value })}
          placeholder="请输入邮箱"
          type="email"
        />
      </FormField>

      {!initialData && (
        <FormField label="密码" required>
          <Input
            value={(formData as CreateUserRequest).password || ''}
            onChange={(e) => setFormData({ ...formData, password: e.target.value })}
            placeholder="请输入密码"
            type="password"
          />
        </FormField>
      )}

      <FormField label="角色" required>
        <Select
          value={formData.role}
          onChange={(value) => setFormData({ ...formData, role: value as any })}
          options={USER_ROLE_OPTIONS.filter((opt) => opt.value !== '')}
        />
      </FormField>

      <FormField label="状态" required>
        <Select
          value={formData.status}
          onChange={(value) => setFormData({ ...formData, status: value as any })}
          options={USER_STATUS_OPTIONS.filter((opt) => opt.value !== '')}
        />
      </FormField>
    </Modal>
  );
};
