import React, { useState, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, Table, Tag, Input, Select, Pagination, Badge } from '../../components';
import type { TableColumn } from '../../components/Table/Table';
import { mockOrders, mockStatistics } from '../../services/mockData';
import type { Order, OrderQueryParams } from '../../types/order.types';
import { OrderStatus, ReviewStatus, GameType, ServiceType } from '../../types/order.types';
import {
  formatOrderStatus,
  getOrderStatusColor,
  formatReviewStatus,
  getReviewStatusColor,
  formatGameType,
  formatServiceType,
  formatCurrency,
  formatDuration,
  formatDateTime,
} from '../../utils/formatters';
import styles from './OrderList.module.less';

export const OrderList: React.FC = () => {
  const navigate = useNavigate();

  // 查询参数
  const [queryParams, setQueryParams] = useState<OrderQueryParams>({
    page: 1,
    pageSize: 10,
    keyword: '',
    status: undefined,
    reviewStatus: undefined,
    gameType: undefined,
    serviceType: undefined,
  });

  // 筛选后的订单数据
  const filteredOrders = useMemo(() => {
    let result = [...mockOrders];

    // 关键词搜索
    if (queryParams.keyword) {
      const keyword = queryParams.keyword.toLowerCase();
      result = result.filter(
        (order) =>
          order.orderNo.toLowerCase().includes(keyword) ||
          order.user.username.toLowerCase().includes(keyword) ||
          order.player?.username.toLowerCase().includes(keyword),
      );
    }

    // 状态筛选
    if (queryParams.status) {
      result = result.filter((order) => order.status === queryParams.status);
    }

    // 审核状态筛选
    if (queryParams.reviewStatus) {
      result = result.filter((order) => order.reviewStatus === queryParams.reviewStatus);
    }

    // 游戏类型筛选
    if (queryParams.gameType) {
      result = result.filter((order) => order.gameType === queryParams.gameType);
    }

    // 服务类型筛选
    if (queryParams.serviceType) {
      result = result.filter((order) => order.serviceType === queryParams.serviceType);
    }

    return result;
  }, [queryParams]);

  // 分页数据
  const paginatedOrders = useMemo(() => {
    const start = (queryParams.page! - 1) * queryParams.pageSize!;
    const end = start + queryParams.pageSize!;
    return filteredOrders.slice(start, end);
  }, [filteredOrders, queryParams.page, queryParams.pageSize]);

  // 表格列定义
  const columns: TableColumn<Order>[] = [
    {
      key: 'orderNo',
      title: '订单号',
      dataIndex: 'orderNo',
      width: 150,
      render: (value, record) => (
        <div className={styles.orderNo} onClick={() => navigate(`/orders/${record.id}`)}>
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
          <div className={styles.userName}>{record.user.username}</div>
          {record.user.phone && <div className={styles.userPhone}>{record.user.phone}</div>}
        </div>
      ),
    },
    {
      key: 'player',
      title: '陪玩者',
      width: 120,
      render: (_, record) => (
        <div>
          {record.player ? (
            <>
              <div className={styles.playerName}>{record.player.username}</div>
              <div className={styles.playerRating}>
                ⭐ {record.player.rating} ({record.player.completedOrders}单)
              </div>
            </>
          ) : (
            <span className={styles.noPlayer}>待接单</span>
          )}
        </div>
      ),
    },
    {
      key: 'game',
      title: '游戏/服务',
      width: 150,
      render: (_, record) => (
        <div>
          <div>{formatGameType(record.gameType)}</div>
          <Tag color="info">{formatServiceType(record.serviceType)}</Tag>
        </div>
      ),
    },
    {
      key: 'price',
      title: '金额',
      dataIndex: 'price',
      width: 100,
      align: 'right',
      render: (value) => <span className={styles.price}>{formatCurrency(value)}</span>,
    },
    {
      key: 'duration',
      title: '时长',
      dataIndex: 'duration',
      width: 80,
      align: 'center',
      render: (value) => formatDuration(value),
    },
    {
      key: 'status',
      title: '订单状态',
      dataIndex: 'status',
      width: 120,
      render: (value) => <Tag color={getOrderStatusColor(value)}>{formatOrderStatus(value)}</Tag>,
    },
    {
      key: 'reviewStatus',
      title: '审核状态',
      width: 120,
      render: (_, record) =>
        record.reviewStatus ? (
          <Tag color={getReviewStatusColor(record.reviewStatus)}>
            {formatReviewStatus(record.reviewStatus)}
          </Tag>
        ) : (
          <span className={styles.noReview}>-</span>
        ),
    },
    {
      key: 'createdAt',
      title: '创建时间',
      dataIndex: 'createdAt',
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

  const handleReviewStatusChange = (value: string | number) => {
    setQueryParams({
      ...queryParams,
      reviewStatus: value as ReviewStatus,
      page: 1,
    });
  };

  const handleGameTypeChange = (value: string | number) => {
    setQueryParams({
      ...queryParams,
      gameType: value as GameType,
      page: 1,
    });
  };

  const handleServiceTypeChange = (value: string | number) => {
    setQueryParams({
      ...queryParams,
      serviceType: value as ServiceType,
      page: 1,
    });
  };

  // 处理分页
  const handlePageChange = (page: number) => {
    setQueryParams({ ...queryParams, page });
  };

  const handleSizeChange = (size: number) => {
    setQueryParams({ ...queryParams, pageSize: size, page: 1 });
  };

  // 重置筛选
  const handleReset = () => {
    setQueryParams({
      page: 1,
      pageSize: 10,
      keyword: '',
      status: undefined,
      reviewStatus: undefined,
      gameType: undefined,
      serviceType: undefined,
    });
  };

  return (
    <div className={styles.container}>
      {/* 统计卡片 */}
      <div className={styles.statsGrid}>
        <Card className={styles.statCard}>
          <div className={styles.statLabel}>总订单数</div>
          <div className={styles.statValue}>{mockStatistics.total}</div>
        </Card>
        <Card className={styles.statCard}>
          <div className={styles.statLabel}>待审核</div>
          <div className={styles.statValue}>
            <Badge count={mockStatistics.pendingReview} />
            {mockStatistics.pendingReview}
          </div>
        </Card>
        <Card className={styles.statCard}>
          <div className={styles.statLabel}>进行中</div>
          <div className={styles.statValue}>{mockStatistics.inProgress}</div>
        </Card>
        <Card className={styles.statCard}>
          <div className={styles.statLabel}>今日订单</div>
          <div className={styles.statValue}>{mockStatistics.todayOrders}</div>
        </Card>
        <Card className={styles.statCard}>
          <div className={styles.statLabel}>今日收入</div>
          <div className={styles.statValue}>{formatCurrency(mockStatistics.todayRevenue)}</div>
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

          <Select
            value={queryParams.reviewStatus}
            placeholder="审核状态"
            options={[
              { label: '全部审核状态', value: '' },
              ...Object.values(ReviewStatus).map((status) => ({
                label: formatReviewStatus(status),
                value: status,
              })),
            ]}
            onChange={handleReviewStatusChange}
            className={styles.filterSelect}
          />

          <Select
            value={queryParams.gameType}
            placeholder="游戏类型"
            options={[
              { label: '全部游戏', value: '' },
              ...Object.values(GameType).map((type) => ({
                label: formatGameType(type),
                value: type,
              })),
            ]}
            onChange={handleGameTypeChange}
            className={styles.filterSelect}
          />

          <Select
            value={queryParams.serviceType}
            placeholder="服务类型"
            options={[
              { label: '全部服务', value: '' },
              ...Object.values(ServiceType).map((type) => ({
                label: formatServiceType(type),
                value: type,
              })),
            ]}
            onChange={handleServiceTypeChange}
            className={styles.filterSelect}
          />

          <button className={styles.resetButton} onClick={handleReset}>
            重置
          </button>
        </div>
      </Card>

      {/* 订单表格 */}
      <Card className={styles.tableCard}>
        <Table columns={columns} dataSource={paginatedOrders} rowKey="id" scroll={{ x: 1400 }} />

        <Pagination
          current={queryParams.page!}
          total={filteredOrders.length}
          pageSize={queryParams.pageSize!}
          onChange={handlePageChange}
          onSizeChange={handleSizeChange}
          showSizeChanger
          className={styles.pagination}
        />
      </Card>
    </div>
  );
};
