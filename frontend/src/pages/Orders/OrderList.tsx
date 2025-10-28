import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, Table, Tag, Input, Select, Pagination, Badge } from '../../components';
import type { TableColumn } from '../../components/Table/Table';
import { orderApi, OrderInfo } from '../../services/api/order';
import type { OrderListQuery } from '../../types/order';
import { OrderStatus } from '../../types/order';
import {
  formatOrderStatus,
  getOrderStatusColor,
  formatCurrency,
  formatDateTime,
} from '../../utils/formatters';
import styles from './OrderList.module.less';


export const OrderList: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [orders, setOrders] = useState<OrderInfo[]>([]);
  const [total, setTotal] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [statistics, setStatistics] = useState({
    total: 0,
    pending: 0,
    in_progress: 0,
    today_orders: 0,
    today_revenue: 0,
  });

  // 查询参数
  const [queryParams, setQueryParams] = useState<OrderListQuery>({
    page: 1,
    page_size: 10,
    keyword: '',
    status: undefined,
  });

  // 加载订单数据
  const loadOrders = async () => {
    setLoading(true);
    setError(null);

    try {
      const result = await orderApi.getList({
        page: queryParams.page,
        page_size: queryParams.page_size,
        keyword: queryParams.keyword || undefined,
        status: queryParams.status,
      });

      setOrders(result.list);
      setTotal(result.total);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '加载订单列表失败';
      setError(errorMessage);
      console.error('加载订单列表失败:', err);
    } finally {
      setLoading(false);
    }
  };

  // 加载统计数据
  const loadStatistics = async () => {
    try {
      const stats = await orderApi.getStatistics();
      setStatistics(stats);
    } catch (err) {
      console.error('加载统计数据失败:', err);
    }
  };

  // 查询参数变化时重新加载
  useEffect(() => {
    loadOrders();
    loadStatistics();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [queryParams.page, queryParams.page_size, queryParams.status]);

  // 表格列定义
  const columns: TableColumn<OrderInfo>[] = [
    {
      key: 'id',
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      key: 'title',
      title: '订单标题',
      dataIndex: 'title',
      width: 200,
      render: (value, record) => (
        <div className={styles.orderTitle} onClick={() => navigate(`/orders/${record.id}`)}>
          {value}
        </div>
      ),
    },
    {
      key: 'user',
      title: '用户',
      width: 120,
      render: (_, record) => (
        <div>
          {record.user ? (
            <div className={styles.userName}>{record.user.name}</div>
          ) : (
            <span>-</span>
          )}
        </div>
      ),
    },
    {
      key: 'player',
      title: '陪玩师',
      width: 120,
      render: (_, record) => (
        <div>
          {record.player ? (
            <div className={styles.playerName}>{record.player.nickname || record.player.id}</div>
          ) : (
            <span className={styles.noPlayer}>待接单</span>
          )}
        </div>
      ),
    },
    {
      key: 'game',
      title: '游戏',
      width: 120,
      render: (_, record) => (
        <div>
          {record.game ? record.game.name : '-'}
        </div>
      ),
    },
    {
      key: 'price_cents',
      title: '金额',
      dataIndex: 'price_cents',
      width: 100,
      align: 'right',
      render: (value) => <span className={styles.price}>{formatCurrency(value)}</span>,
    },
    {
      key: 'status',
      title: '订单状态',
      dataIndex: 'status',
      width: 120,
      render: (value) => <Tag color={getOrderStatusColor(value)}>{formatOrderStatus(value)}</Tag>,
    },
    {
      key: 'created_at',
      title: '创建时间',
      dataIndex: 'created_at',
      width: 160,
      render: (value) => formatDateTime(value),
    },
    {
      key: 'actions',
      title: '操作',
      width: 100,
      fixed: 'right',
      render: (_, record) => (
        <div className={styles.actions}>
          <button className={styles.actionButton} onClick={() => navigate(`/orders/${record.id}`)}>
            查看
          </button>
        </div>
      ),
    },
  ];

  // 处理搜索
  const handleSearch = (value: string) => {
    setQueryParams({ ...queryParams, keyword: value, page: 1 });
  };

  // 处理筛选
  const handleStatusChange = (value: string | number) => {
    setQueryParams({
      ...queryParams,
      status: value as OrderStatus,
      page: 1,
    });
  };


  // 处理分页
  const handlePageChange = (page: number) => {
    setQueryParams({ ...queryParams, page });
  };

  const handleSizeChange = (size: number) => {
    setQueryParams({ ...queryParams, page_size: size, page: 1 });
  };

  // 重置筛选
  const handleReset = () => {
    setQueryParams({
      page: 1,
      page_size: 10,
      keyword: '',
      status: undefined,
    });
  };

  return (
    <div className={styles.container}>
      {/* 统计卡片 */}
      <div className={styles.statsGrid}>
        <Card className={styles.statCard}>
          <div className={styles.statLabel}>总订单数</div>
          <div className={styles.statValue}>{statistics.total}</div>
        </Card>
        <Card className={styles.statCard}>
          <div className={styles.statLabel}>待审核</div>
          <div className={styles.statValue}>
            <Badge count={statistics.pending} />
            {statistics.pending}
          </div>
        </Card>
        <Card className={styles.statCard}>
          <div className={styles.statLabel}>进行中</div>
          <div className={styles.statValue}>{statistics.in_progress}</div>
        </Card>
        <Card className={styles.statCard}>
          <div className={styles.statLabel}>今日订单</div>
          <div className={styles.statValue}>{statistics.today_orders}</div>
        </Card>
        <Card className={styles.statCard}>
          <div className={styles.statLabel}>今日收入</div>
          <div className={styles.statValue}>{formatCurrency(statistics.today_revenue)}</div>
        </Card>
      </div>

      {/* 筛选区域 */}
      <Card className={styles.filterCard}>
        <div className={styles.filterRow}>
          <Input
            placeholder="搜索订单号、用户、陪玩者"
            value={queryParams.keyword}
            onChange={(e) => handleSearch(e.target.value)}
            className={styles.searchInput}
          />

          <Select
            value={queryParams.status}
            placeholder="订单状态"
            options={[
              { label: '全部状态', value: '' },
              ...Object.values(OrderStatus).map((status) => ({
                label: formatOrderStatus(status),
                value: status,
              })),
            ]}
            onChange={handleStatusChange}
            className={styles.filterSelect}
          />


          <button className={styles.resetButton} onClick={handleReset}>
            重置
          </button>
        </div>
      </Card>

      {/* 订单表格 */}
      <Card className={styles.tableCard}>
        {loading ? (
          <p>加载中...</p>
        ) : error ? (
          <p>加载失败: {error}</p>
        ) : (
          <>
            <Table columns={columns} dataSource={orders} rowKey="id" scroll={{ overflowX: 'auto', minWidth: '1400px' } as React.CSSProperties} />

            <Pagination
              current={queryParams.page!}
              total={total}
              pageSize={queryParams.page_size!}
              onChange={handlePageChange}
              onSizeChange={handleSizeChange}
              showSizeChanger
              className={styles.pagination}
            />
          </>
        )}
      </Card>
    </div>
  );
};
