import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  DataTable,
  Button,
  Input,
  Select,
  Tag,
  ActionButtons,
  DeleteConfirmModal,
} from '../../components';
import type { FilterConfig } from '../../components/DataTable';
import type { TableColumn } from '../../components/Table/Table';
import { userApi } from '../../services/api/user';
import type { User, UserListQuery, CreateUserRequest, UpdateUserRequest } from '../../types/user';
import { formatDateTime, formatRelativeTime } from '../../utils/formatters';
import { useListPage } from '../../hooks/useListPage';
import {
  formatUserRole,
  getUserRoleColor,
  formatUserStatus,
  getUserStatusColor,
} from '../../utils/statusHelpers';
import { USER_ROLE_OPTIONS, USER_STATUS_OPTIONS } from '../../utils/selectOptions';
import { UserFormModal } from './UserFormModal';
import styles from './UserList.module.less';

export const UserList: React.FC = () => {
  const navigate = useNavigate();

  // 使用通用列表Hook
  const {
    loading,
    data: users,
    total,
    queryParams,
    setQueryParams,
    handleSearch,
    handleReset,
    handlePageChange,
    reload,
  } = useListPage<User, UserListQuery>({
    initialParams: {
      page: 1,
      pageSize: 10,
      keyword: '',
      role: undefined,
      status: undefined,
    },
    fetchData: async (params) => {
      return await userApi.getList({
        page: params.page,
        pageSize: params.pageSize,
        keyword: params.keyword || undefined,
        role: params.role,
        status: params.status,
      });
    },
    urlParamKeys: ['role', 'status'], // 从URL读取角色和状态参数
  });

  // Modal状态
  const [showUserModal, setShowUserModal] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [deletingUser, setDeletingUser] = useState<User | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  // 重置处理
  const handleResetClick = async () => {
    await handleReset({
      page: 1,
      pageSize: 10,
      keyword: '',
      role: undefined,
      status: undefined,
    });
  };

  // 创建用户
  const handleCreateUser = () => {
    setEditingUser(null);
    setShowUserModal(true);
  };

  // 编辑用户
  const handleEditUser = (user: User) => {
    setEditingUser(user);
    setShowUserModal(true);
  };

  // 删除用户
  const handleDeleteUser = (user: User) => {
    setDeletingUser(user);
  };

  // 确认删除
  const handleConfirmDelete = async () => {
    if (!deletingUser) return;

    setIsSubmitting(true);
    try {
      await userApi.delete(deletingUser.id);
      setDeletingUser(null);
      await reload();
    } catch (err) {
      console.error('删除失败:', err);
    } finally {
      setIsSubmitting(false);
    }
  };

  // 保存用户
  const handleSaveUser = async (data: CreateUserRequest | UpdateUserRequest) => {
    setIsSubmitting(true);
    try {
      if (editingUser) {
        await userApi.update(editingUser.id, data as UpdateUserRequest);
      } else {
        await userApi.create(data as CreateUserRequest);
      }
      setShowUserModal(false);
      setEditingUser(null);
      await reload();
    } catch (err) {
      console.error('操作失败:', err);
      throw err;
    } finally {
      setIsSubmitting(false);
    }
  };

  // 表格列定义
  const columns: TableColumn<User>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: '80px',
    },
    {
      title: '用户信息',
      key: 'userInfo',
      render: (_: unknown, record: User) => (
        <div className={styles.userInfo}>
          {record.avatarUrl && (
            <img src={record.avatarUrl} alt={record.name} className={styles.avatar} />
          )}
          {!record.avatarUrl && (
            <div className={styles.avatarPlaceholder}>{record.name.charAt(0)}</div>
          )}
          <div className={styles.userDetails}>
            <div className={styles.userName}>{record.name}</div>
            <div className={styles.userContact}>
              {record.phone && <span>📱 {record.phone}</span>}
              {record.email && <span>✉️ {record.email}</span>}
            </div>
          </div>
        </div>
      ),
    },
    {
      title: '角色',
      key: 'role',
      width: '100px',
      render: (_: unknown, record: User) => (
        <Tag color={getUserRoleColor(record.role) as any}>{formatUserRole(record.role)}</Tag>
      ),
    },
    {
      title: '状态',
      key: 'status',
      width: '100px',
      render: (_: unknown, record: User) => (
        <Tag color={getUserStatusColor(record.status) as any}>
          {formatUserStatus(record.status)}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      key: 'createdAt',
      width: '160px',
      render: (_: unknown, record: User) => (
        <div className={styles.timeInfo}>
          <div>{formatRelativeTime(record.createdAt)}</div>
          <div className={styles.timeDetail}>{formatDateTime(record.createdAt)}</div>
        </div>
      ),
    },
    {
      title: '操作',
      key: 'actions',
      width: '200px',
      render: (_: unknown, record: User) => (
        <ActionButtons
          onView={() => navigate(`/users/${record.id}`)}
          onEdit={() => handleEditUser(record)}
          onDelete={() => handleDeleteUser(record)}
        />
      ),
    },
  ];

  // 筛选配置
  const filters: FilterConfig[] = [
    {
      label: '搜索',
      key: 'keyword',
      element: (
        <Input
          value={queryParams.keyword}
          onChange={(e) => setQueryParams((prev) => ({ ...prev, keyword: e.target.value }))}
          placeholder="用户名/手机/邮箱"
          onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
        />
      ),
    },
    {
      label: '角色',
      key: 'role',
      element: (
        <Select
          value={queryParams.role || ''}
          onChange={(value) =>
            setQueryParams((prev) => ({ ...prev, role: value ? (value as any) : undefined }))
          }
          options={USER_ROLE_OPTIONS}
        />
      ),
    },
    {
      label: '状态',
      key: 'status',
      element: (
        <Select
          value={queryParams.status || ''}
          onChange={(value) =>
            setQueryParams((prev) => ({ ...prev, status: value ? (value as any) : undefined }))
          }
          options={USER_STATUS_OPTIONS}
        />
      ),
    },
  ];

  // 筛选操作按钮
  const filterActions = (
    <>
      <Button variant="primary" onClick={handleSearch}>
        搜索
      </Button>
      <Button variant="outlined" onClick={handleResetClick}>
        重置
      </Button>
    </>
  );

  // 头部操作按钮
  const headerActions = (
    <Button variant="primary" onClick={handleCreateUser}>
      新增用户
    </Button>
  );

  return (
    <>
      <DataTable
        title="用户管理"
        headerActions={headerActions}
        filters={filters}
        filterActions={filterActions}
        columns={columns}
        dataSource={users}
        loading={loading}
        rowKey="id"
        pagination={{
          current: queryParams.page || 1,
          pageSize: queryParams.pageSize || 10,
          total,
          onChange: handlePageChange,
        }}
      />

      {/* 表单Modal */}
      <UserFormModal
        visible={showUserModal}
        onClose={() => {
          setShowUserModal(false);
          setEditingUser(null);
        }}
        onSave={handleSaveUser}
        initialData={editingUser}
        isSubmitting={isSubmitting}
      />

      {/* 删除确认Modal */}
      <DeleteConfirmModal
        visible={!!deletingUser}
        content={`确定要删除用户 "${deletingUser?.name}" 吗？此操作不可恢复。`}
        onConfirm={handleConfirmDelete}
        onCancel={() => setDeletingUser(null)}
        loading={isSubmitting}
      />
    </>
  );
};
