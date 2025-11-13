import React, { useState } from 'react';
import {
  DataTable,
  Button,
  Input,
  Select,
  Tag,
  DeleteConfirmModal,
} from '../../components';
import type { TableColumn } from '../../components/Table/Table';
import { permissionApi } from '../../services/api/rbac';
import type { Permission, PermissionListQuery, HTTPMethod } from '../../types/rbac';
import { formatDateTime } from '../../utils/formatters';
import { useListPage } from '../../hooks/useListPage';
import { HTTP_METHOD_OPTIONS } from '../../utils/selectOptions';
import { getHTTPMethodColor } from '../../utils/statusHelpers';
import { PermissionFormModal } from './PermissionFormModal';
import styles from './PermissionManagement.module.less';

/**
 * 权限管理组件
 */
export const PermissionManagement: React.FC = () => {
  // 使用通用列表Hook
  const {
    loading,
    data: permissions,
    total,
    queryParams,
    setQueryParams,
    handleSearch,
    handleReset,
    handlePageChange,
    reload,
  } = useListPage<Permission, PermissionListQuery>({
    initialParams: {
      page: 1,
      pageSize: 10,
      keyword: '',
      method: undefined,
      group: undefined,
    },
    fetchData: async (params) => {
      return await permissionApi.getList({
        page: params.page,
        pageSize: params.pageSize,
        keyword: params.keyword || undefined,
        method: params.method,
        group: params.group,
      });
    },
  });

  // Modal状态
  const [showPermissionModal, setShowPermissionModal] = useState(false);
  const [editingPermission, setEditingPermission] = useState<Permission | null>(null);
  const [deletingPermission, setDeletingPermission] = useState<Permission | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  // 重置处理
  const handleResetClick = async () => {
    await handleReset({
      page: 1,
      pageSize: 10,
      keyword: '',
      method: undefined,
      group: undefined,
    });
  };

  // 创建权限
  const handleCreatePermission = () => {
    setEditingPermission(null);
    setShowPermissionModal(true);
  };

  // 编辑权限
  const handleEditPermission = (permission: Permission) => {
    setEditingPermission(permission);
    setShowPermissionModal(true);
  };

  // 删除权限
  const handleDeletePermission = (permission: Permission) => {
    setDeletingPermission(permission);
  };

  // 确认删除
  const handleConfirmDelete = async () => {
    if (!deletingPermission) return;

    try {
      setIsSubmitting(true);
      await permissionApi.delete(deletingPermission.id);
      setDeletingPermission(null);
      reload();
    } catch (error) {
      console.error('删除权限失败:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  // 表格列定义
  const columns: TableColumn<Permission>[] = [
    {
      key: 'id',
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      key: 'method',
      title: 'HTTP方法',
      dataIndex: 'method',
      width: 100,
      render: (method: string) => (
        <Tag color={getHTTPMethodColor(method as any) as any} className={styles.methodTag}>
          {method}
        </Tag>
      ),
    },
    {
      key: 'path',
      title: 'API路径',
      dataIndex: 'path',
      width: 300,
      render: (path: string) => <code className={styles.path}>{path}</code>,
    },
    {
      key: 'code',
      title: '权限代码',
      dataIndex: 'code',
      width: 200,
      render: (code: string) => <code className={styles.code}>{code}</code>,
    },
    {
      key: 'group',
      title: '分组',
      dataIndex: 'group',
      width: 150,
      render: (group: string) => (
        <span className={styles.groupTag}>{group}</span>
      ),
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
      key: 'createdAt',
      title: '创建时间',
      dataIndex: 'createdAt',
      width: 180,
      render: (date: string) => formatDateTime(date),
    },
    {
      key: 'actions',
      title: '操作',
      width: 150,
      fixed: 'right',
      render: (_: unknown, record: Permission) => (
        <div className={styles.actions}>
          <Button variant="text" onClick={() => handleEditPermission(record)}>
            编辑
          </Button>
          <Button variant="text" onClick={() => handleDeletePermission(record)}>
            删除
          </Button>
        </div>
      ),
    },
  ];

  // 头部操作按钮
  const headerActions = (
    <Button variant="primary" onClick={handleCreatePermission}>
      创建权限
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
          placeholder="搜索权限代码、路径或描述..."
          style={{ width: 200 }}
        />
      ),
    },
    {
      key: 'method',
      label: 'HTTP方法',
      element: (
        <Select
          value={queryParams.method || ''}
          onChange={(value) =>
            setQueryParams({
              ...queryParams,
              method: (value as HTTPMethod) || undefined,
            })
          }
          options={HTTP_METHOD_OPTIONS}
          placeholder="全部方法"
          style={{ width: 150 }}
        />
      ),
    },
    {
      key: 'group',
      label: '权限分组',
      element: (
        <Input
          value={queryParams.group || ''}
          onChange={(e) =>
            setQueryParams({
              ...queryParams,
              group: e.target.value || undefined,
            })
          }
          placeholder="输入分组"
          style={{ width: 200 }}
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
        title="权限管理"
        headerActions={headerActions}
        filters={filters}
        filterActions={filterActions}
        columns={columns}
        dataSource={permissions}
        loading={loading}
        rowKey="id"
        pagination={{
          current: queryParams.page || 1,
          pageSize: queryParams.pageSize || 10,
          total,
          onChange: handlePageChange,
        }}
      />

      {/* 权限表单模态框 */}
      {showPermissionModal && (
        <PermissionFormModal
          visible={showPermissionModal}
          permission={editingPermission}
          onClose={() => {
            setShowPermissionModal(false);
            setEditingPermission(null);
          }}
          onSuccess={() => {
            setShowPermissionModal(false);
            setEditingPermission(null);
            reload();
          }}
        />
      )}

      {/* 删除确认模态框 */}
      {deletingPermission && (
        <DeleteConfirmModal
          visible={!!deletingPermission}
          title="删除权限"
          content={
            <>
              确定要删除权限 <strong>{deletingPermission.code}</strong> 吗？
              <br />
              <span style={{ color: 'var(--text-secondary)' }}>
                警告：删除权限可能影响已关联的角色！
              </span>
            </>
          }
          onConfirm={handleConfirmDelete}
          onCancel={() => setDeletingPermission(null)}
          loading={isSubmitting}
        />
      )}
    </>
  );
};

