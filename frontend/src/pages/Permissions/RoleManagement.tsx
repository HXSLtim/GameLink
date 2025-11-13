import React, { useState } from 'react';
import {
  DataTable,
  Button,
  Input,
  Select,
  DeleteConfirmModal,
} from '../../components';
import type { TableColumn } from '../../components/Table/Table';
import { roleApi } from '../../services/api/rbac';
import type { Role, RoleListQuery } from '../../types/rbac';
import { formatDateTime } from '../../utils/formatters';
import { useListPage } from '../../hooks/useListPage';
import { BOOLEAN_OPTIONS } from '../../utils/selectOptions';
import { RoleFormModal } from './RoleFormModal';
import { PermissionAssignModal } from './PermissionAssignModal';
import styles from './RoleManagement.module.less';

/**
 * 角色管理组件
 */
export const RoleManagement: React.FC = () => {
  // 使用通用列表Hook
  const {
    loading,
    data: roles,
    total,
    queryParams,
    setQueryParams,
    handleSearch,
    handleReset,
    handlePageChange,
    reload,
  } = useListPage<Role, RoleListQuery>({
    initialParams: {
      page: 1,
      pageSize: 10,
      keyword: '',
      isSystem: undefined,
    },
    fetchData: async (params) => {
      return await roleApi.getList({
        page: params.page,
        pageSize: params.pageSize,
        keyword: params.keyword || undefined,
        isSystem: params.isSystem,
      });
    },
  });

  // Modal状态
  const [showRoleModal, setShowRoleModal] = useState(false);
  const [editingRole, setEditingRole] = useState<Role | null>(null);
  const [deletingRole, setDeletingRole] = useState<Role | null>(null);
  const [assigningPermissionRole, setAssigningPermissionRole] = useState<Role | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  // 重置处理
  const handleResetClick = async () => {
    await handleReset({
      page: 1,
      pageSize: 10,
      keyword: '',
      isSystem: undefined,
    });
  };

  // 创建角色
  const handleCreateRole = () => {
    setEditingRole(null);
    setShowRoleModal(true);
  };

  // 编辑角色
  const handleEditRole = (role: Role) => {
    setEditingRole(role);
    setShowRoleModal(true);
  };

  // 删除角色
  const handleDeleteRole = (role: Role) => {
    setDeletingRole(role);
  };

  // 分配权限
  const handleAssignPermission = (role: Role) => {
    setAssigningPermissionRole(role);
  };

  // 确认删除
  const handleConfirmDelete = async () => {
    if (!deletingRole) return;

    try {
      setIsSubmitting(true);
      await roleApi.delete(deletingRole.id);
      setDeletingRole(null);
      reload();
    } catch (error) {
      console.error('删除角色失败:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  // 表格列定义
  const columns: TableColumn<Role>[] = [
    {
      key: 'id',
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      key: 'slug',
      title: '角色标识',
      dataIndex: 'slug',
      width: 150,
      render: (slug: string) => <code>{slug}</code>,
    },
    {
      key: 'name',
      title: '角色名称',
      dataIndex: 'name',
      width: 150,
    },
    {
      key: 'description',
      title: '描述',
      dataIndex: 'description',
      render: (description: string) => (
        <span className={styles.description}>{description || '-'}</span>
      ),
    },
    {
      key: 'isSystem',
      title: '系统角色',
      dataIndex: 'isSystem',
      width: 100,
      render: (isSystem: boolean) => (
        <span className={isSystem ? styles.systemBadge : styles.normalBadge}>
          {isSystem ? '系统' : '自定义'}
        </span>
      ),
    },
    {
      key: 'permissionCount',
      title: '权限数',
      width: 100,
      render: (_: unknown, record: Role) => record.permissions?.length || 0,
    },
    {
      key: 'createdAt',
      title: '创建时间',
      dataIndex: 'createdAt',
      width: 180,
      render: (date: string) => formatDateTime(date),
    },
    {
      key: 'actions',
      title: '操作',
      width: 200,
      fixed: 'right',
      render: (_: unknown, record: Role) => (
        <div className={styles.actions}>
          <Button variant="text" onClick={() => handleAssignPermission(record)}>
            权限
          </Button>
          <Button 
            variant="text" 
            onClick={() => handleEditRole(record)}
            disabled={record.isSystem}
          >
            编辑
          </Button>
          <Button 
            variant="text" 
            onClick={() => handleDeleteRole(record)}
            disabled={record.isSystem}
            style={{ color: record.isSystem ? undefined : 'var(--text-primary)' }}
          >
            删除
          </Button>
        </div>
      ),
    },
  ];

  // 头部操作按钮
  const headerActions = (
    <Button variant="primary" onClick={handleCreateRole}>
      创建角色
    </Button>
  );

  // 筛选器
  const filters = [
    {
      key: 'keyword',
      label: '关键词',
      element: (
        <Input
          value={queryParams.keyword || ''}
          onChange={(e) => setQueryParams({ ...queryParams, keyword: e.target.value })}
          placeholder="搜索角色名称或标识..."
          style={{ width: 200 }}
        />
      ),
    },
    {
      key: 'isSystem',
      label: '角色类型',
      element: (
        <Select
          value={
            queryParams.isSystem === undefined ? '' : queryParams.isSystem ? 'true' : 'false'
          }
          onChange={(value) =>
            setQueryParams({
              ...queryParams,
              isSystem: value === '' ? undefined : value === 'true',
            })
          }
          options={BOOLEAN_OPTIONS}
          placeholder="全部类型"
          style={{ width: 150 }}
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
      <Button variant="secondary" onClick={handleResetClick}>重置</Button>
    </>
  );

  return (
    <>
      <DataTable
        title="角色管理"
        headerActions={headerActions}
        filters={filters}
        filterActions={filterActions}
        columns={columns}
        dataSource={roles}
        loading={loading}
        rowKey="id"
        pagination={{
          current: queryParams.page || 1,
          pageSize: queryParams.pageSize || 10,
          total,
          onChange: handlePageChange,
        }}
      />

      {/* 角色表单模态框 */}
      {showRoleModal && (
        <RoleFormModal
          visible={showRoleModal}
          role={editingRole}
          onClose={() => {
            setShowRoleModal(false);
            setEditingRole(null);
          }}
          onSuccess={() => {
            setShowRoleModal(false);
            setEditingRole(null);
            reload();
          }}
        />
      )}

      {/* 权限分配模态框 */}
      {assigningPermissionRole && (
        <PermissionAssignModal
          visible={!!assigningPermissionRole}
          role={assigningPermissionRole}
          onClose={() => setAssigningPermissionRole(null)}
          onSuccess={() => {
            setAssigningPermissionRole(null);
            reload();
          }}
        />
      )}

      {/* 删除确认模态框 */}
      {deletingRole && !deletingRole.isSystem && (
        <DeleteConfirmModal
          visible={!!deletingRole}
          title="删除角色"
          content={
            <>
              确定要删除角色 <strong>{deletingRole.name}</strong> 吗？
              <br />
              <span style={{ color: 'var(--text-secondary)' }}>
                删除后将无法恢复。
              </span>
            </>
          }
          onConfirm={handleConfirmDelete}
          onCancel={() => setDeletingRole(null)}
          loading={isSubmitting}
        />
      )}
    </>
  );
};

