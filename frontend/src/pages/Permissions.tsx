import { useMemo } from 'react';
import type { TableColumnProps } from '@arco-design/web-react';
import { Table, Typography } from '@arco-design/web-react';
import type { Permission } from '../services/permission';
import { permissionService } from '../services/permission';
import { useTable } from '../hooks/useTable';
import styles from './Permissions.module.less';

/**
 * Permissions page component
 *
 * @component
 * @description Displays a paginated table of system permissions
 */
export const Permissions = () => {
  const { data, loading, pagination, handlePageChange } = useTable<Permission>({
    fetchData: (page, pageSize) => permissionService.list({ page, page_size: pageSize }),
    errorMessage: '获取权限列表',
  });

  const columns = useMemo<TableColumnProps<Permission>[]>(
    () => [
      { title: 'ID', dataIndex: 'id', width: 160 },
      { title: '名称', dataIndex: 'name' },
      { title: '描述', dataIndex: 'desc' },
    ],
    [],
  );

  return (
    <div className={styles.container}>
      <Typography.Title heading={5}>权限管理</Typography.Title>
      <Table
        rowKey="id"
        loading={loading}
        columns={columns}
        data={data}
        pagination={{
          current: pagination.page,
          pageSize: pagination.pageSize,
          total: pagination.total,
        }}
        onChange={handlePageChange}
        border
      />
    </div>
  );
};
