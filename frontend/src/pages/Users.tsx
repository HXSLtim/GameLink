import { useMemo } from 'react';
import type { TableColumnProps } from '@arco-design/web-react';
import { Table, Typography } from '@arco-design/web-react';
import { userService } from '../services/user';
import type { User } from '../types/user';
import { useTable } from '../hooks/useTable';
import styles from './Users.module.less';

/**
 * Users page component
 *
 * @component
 * @description Displays a paginated table of users with filtering and sorting capabilities
 */
export const Users = () => {
  const { data, loading, pagination, handlePageChange } = useTable<User>({
    fetchData: (page, pageSize) => userService.list({ page, page_size: pageSize }),
    errorMessage: '获取用户列表',
  });

  const columns = useMemo<TableColumnProps<User>[]>(
    () => [
      { title: 'ID', dataIndex: 'id', width: 160 },
      { title: '用户名', dataIndex: 'username' },
      { title: '角色', dataIndex: 'role', width: 120 },
      { title: '创建时间', dataIndex: 'created_at', width: 180 },
    ],
    [],
  );

  return (
    <div className={styles.container}>
      <Typography.Title heading={5}>用户管理</Typography.Title>
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
