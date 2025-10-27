import { useMemo } from 'react';
import type { TableColumnProps, BadgeProps } from '@arco-design/web-react';
import { Table, Badge, Typography } from '@arco-design/web-react';
import { orderService } from '../services/order';
import type { Order } from '../types/order';
import { useTable } from '../hooks/useTable';
import styles from './Orders.module.less';

type BadgeStatus = BadgeProps['status'];

const STATUS_BADGE: Record<string, BadgeStatus> = {
  paid: 'success',
  pending: 'processing',
  cancelled: 'warning',
  refunded: 'warning',
  failed: 'error',
};

/**
 * Orders page component
 *
 * @component
 * @description Displays a paginated table of orders with status badges
 */
export const Orders = () => {
  const { data, loading, pagination, handlePageChange } = useTable<Order>({
    fetchData: (page, pageSize) => orderService.list({ page, page_size: pageSize }),
    errorMessage: '获取订单列表',
  });

  const columns = useMemo<TableColumnProps<Order>[]>(
    () => [
      { title: '订单ID', dataIndex: 'id', width: 200 },
      { title: '用户ID', dataIndex: 'user_id', width: 160 },
      {
        title: '金额',
        dataIndex: 'amount',
        width: 120,
        render: (value: number, record: Order) => `${value} ${record.currency || ''}`,
      },
      {
        title: '状态',
        dataIndex: 'status',
        width: 120,
        render: (status: string) => (
          <Badge status={STATUS_BADGE[status] || 'processing'} text={status} />
        ),
      },
      { title: '创建时间', dataIndex: 'created_at', width: 180 },
    ],
    [],
  );

  return (
    <div className={styles.container}>
      <Typography.Title heading={5}>订单管理</Typography.Title>
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
