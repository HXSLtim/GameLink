import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { DataTable, Button, Input, Select, Tag } from '../../components';
import type { FilterConfig } from '../../components/DataTable';
import type { TableColumn } from '../../components/Table/Table';
import { paymentApi } from '../../services/api/payment';
import type { Payment, PaymentListQuery, PaymentStatus, PaymentMethod } from '../../types/payment';
import { formatCurrency, formatDateTime } from '../../utils/formatters';
import styles from './PaymentList.module.less';

export const PaymentList: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [payments, setPayments] = useState<Payment[]>([]);
  const [total, setTotal] = useState(0);

  // 查询参数
  const [queryParams, setQueryParams] = useState<PaymentListQuery>({
    page: 1,
    page_size: 10,
    keyword: '',
    status: undefined,
    method: undefined,
  });

  // 加载支付列表
  const loadPayments = async () => {
    setLoading(true);

    try {
      const result = await paymentApi.getList({
        page: queryParams.page,
        page_size: queryParams.page_size,
        keyword: queryParams.keyword || undefined,
        status: queryParams.status,
        method: queryParams.method,
      });

      if (result && result.list) {
        setPayments(result.list);
        setTotal(result.total || 0);
      } else {
        setPayments([]);
        setTotal(0);
      }
    } catch (err) {
      console.error('加载支付列表失败:', err);
      setPayments([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  };

  // 搜索
  const handleSearch = () => {
    setQueryParams((prev) => ({ ...prev, page: 1 }));
  };

  // 重置
  const handleReset = () => {
    setQueryParams({
      page: 1,
      page_size: 10,
      keyword: '',
      status: undefined,
      method: undefined,
    });
  };

  // 分页变化
  const handlePageChange = (page: number) => {
    setQueryParams((prev) => ({ ...prev, page }));
  };

  // 加载数据
  useEffect(() => {
    loadPayments();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [queryParams.page, queryParams.page_size]);

  // 状态格式化
  const formatStatus = (status: PaymentStatus): string => {
    const statusMap: Record<PaymentStatus, string> = {
      pending: '待支付',
      paid: '已支付',
      failed: '支付失败',
      refunded: '已退款',
      cancelled: '已取消',
    };
    return statusMap[status] || status;
  };

  // 状态颜色
  const getStatusColor = (status: PaymentStatus): string => {
    const colorMap: Record<PaymentStatus, string> = {
      pending: 'orange',
      paid: 'green',
      failed: 'red',
      refunded: 'purple',
      cancelled: 'default',
    };
    return colorMap[status] || 'default';
  };

  // 支付方式格式化
  const formatMethod = (method: PaymentMethod): string => {
    const methodMap: Record<PaymentMethod, string> = {
      alipay: '支付宝',
      wechat: '微信支付',
      balance: '余额支付',
    };
    return methodMap[method] || method;
  };

  // 表格列定义
  const columns: TableColumn<Payment>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: '80px',
    },
    {
      title: '订单ID',
      dataIndex: 'order_id',
      key: 'order_id',
      width: '100px',
    },
    {
      title: '用户ID',
      dataIndex: 'user_id',
      key: 'user_id',
      width: '100px',
    },
    {
      title: '金额',
      key: 'amount',
      width: '120px',
      render: (_: unknown, record: Payment) => (
        <div className={styles.amount}>{formatCurrency(record.amount_cents)}</div>
      ),
    },
    {
      title: '支付方式',
      key: 'method',
      width: '120px',
      render: (_: unknown, record: Payment) => (
        <Tag>{formatMethod(record.method as PaymentMethod)}</Tag>
      ),
    },
    {
      title: '状态',
      key: 'status',
      width: '100px',
      render: (_: unknown, record: Payment) => (
        <Tag color={getStatusColor(record.status as PaymentStatus) as any}>
          {formatStatus(record.status as PaymentStatus)}
        </Tag>
      ),
    },
    {
      title: '交易号',
      dataIndex: 'transaction_id',
      key: 'transaction_id',
      width: '200px',
      render: (transaction_id: string) => (
        <div className={styles.transactionId}>{transaction_id || '-'}</div>
      ),
    },
    {
      title: '创建时间',
      key: 'createdAt',
      width: '160px',
      render: (_: unknown, record: Payment) => formatDateTime(record.created_at),
    },
    {
      title: '操作',
      key: 'actions',
      width: '120px',
      render: (_: unknown, record: Payment) => (
        <Button
          variant="text"
          onClick={() => navigate(`/payments/${record.id}`)}
          className={styles.actionButton}
        >
          查看详情
        </Button>
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
          onChange={(e) =>
            setQueryParams((prev) => ({ ...prev, keyword: e.target.value }))
          }
          placeholder="交易号/订单ID"
          onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
        />
      ),
    },
    {
      label: '支付状态',
      key: 'status',
      element: (
        <Select
          value={queryParams.status || ''}
          onChange={(value) =>
            setQueryParams((prev) => ({
              ...prev,
              status: value ? (value as PaymentStatus) : undefined,
            }))
          }
          options={[
            { label: '全部状态', value: '' },
            { label: '待支付', value: 'pending' },
            { label: '已支付', value: 'paid' },
            { label: '支付失败', value: 'failed' },
            { label: '已退款', value: 'refunded' },
            { label: '已取消', value: 'cancelled' },
          ]}
        />
      ),
    },
    {
      label: '支付方式',
      key: 'method',
      element: (
        <Select
          value={queryParams.method || ''}
          onChange={(value) =>
            setQueryParams((prev) => ({
              ...prev,
              method: value ? (value as PaymentMethod) : undefined,
            }))
          }
          options={[
            { label: '全部方式', value: '' },
            { label: '支付宝', value: 'alipay' },
            { label: '微信支付', value: 'wechat' },
            { label: '余额支付', value: 'balance' },
          ]}
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
      <Button variant="outlined" onClick={handleReset}>
        重置
      </Button>
    </>
  );

  return (
    <DataTable
      title="支付管理"
      filters={filters}
      filterActions={filterActions}
      columns={columns}
      dataSource={payments}
      loading={loading}
      rowKey="id"
      pagination={{
        current: queryParams.page || 1,
        pageSize: queryParams.page_size || 10,
        total,
        onChange: handlePageChange,
      }}
    />
  );
};
