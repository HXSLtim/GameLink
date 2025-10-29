import React, { useState, useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { DataTable, Button, Input, Select, Tag, Modal } from '../../components';
import type { FilterConfig } from '../../components/DataTable';
import type { TableColumn } from '../../components/Table/Table';
import { orderApi } from '../../services/api/order';
import type { Order, OrderListQuery, OrderStatus, UpdateOrderRequest } from '../../types/order';
import { formatCurrency, formatDateTime, formatRelativeTime } from '../../utils/formatters';
import { OrderFormModal } from './OrderFormModal';
import styles from './OrderList.module.less';

export const OrderList: React.FC = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const [loading, setLoading] = useState(false);
  const [orders, setOrders] = useState<Order[]>([]);
  const [total, setTotal] = useState(0);

  // 表单Modal状态
  const [formModalVisible, setFormModalVisible] = useState(false);
  const [editingOrder, setEditingOrder] = useState<Order | null>(null);

  // 删除确认Modal状态
  const [deleteModalVisible, setDeleteModalVisible] = useState(false);
  const [deletingOrder, setDeletingOrder] = useState<Order | null>(null);

  // 从URL读取初始状态
  const getInitialParams = (): OrderListQuery => {
    const statusFromUrl = searchParams.get('status') as OrderStatus | null;
    return {
      page: 1,
      pageSize: 10,
      keyword: '',
      status: statusFromUrl || undefined,
    };
  };

  // 查询参数
  const [queryParams, setQueryParams] = useState<OrderListQuery>(getInitialParams());

  // 加载订单列表
  const loadOrders = async () => {
    setLoading(true);

    try {
      const result = await orderApi.getList({
        page: queryParams.page,
        pageSize: queryParams.pageSize,
        keyword: queryParams.keyword || undefined,
        status: queryParams.status,
      });

      if (result && result.list) {
        setOrders(result.list);
        setTotal(result.total || 0);
      } else {
        setOrders([]);
        setTotal(0);
      }
    } catch (err) {
      console.error('加载订单列表失败:', err);
      setOrders([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  };

  // 搜索
  const handleSearch = async () => {
    setQueryParams((prev) => ({ ...prev, page: 1 }));
    await loadOrders();
  };

  // 重置
  const handleReset = async () => {
    const resetParams = {
      page: 1,
      pageSize: 10,
      keyword: '',
      status: undefined,
    };
    setQueryParams(resetParams);
    // 使用重置后的参数立即加载数据
    setLoading(true);
    try {
      const result = await orderApi.getList({
        page: resetParams.page,
        pageSize: resetParams.pageSize,
      });
      if (result && result.list) {
        setOrders(result.list);
        setTotal(result.total || 0);
      } else {
        setOrders([]);
        setTotal(0);
      }
    } catch (err) {
      console.error('加载订单列表失败:', err);
      setOrders([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  };

  // 分页变化
  const handlePageChange = (page: number) => {
    setQueryParams((prev) => ({ ...prev, page }));
  };

  // 编辑订单
  const handleEdit = (order: Order) => {
    setEditingOrder(order);
    setFormModalVisible(true);
  };

  // 提交表单
  const handleFormSubmit = async (data: UpdateOrderRequest) => {
    if (!editingOrder) return;

    try {
      await orderApi.update(editingOrder.id, data);
      await loadOrders();
    } catch (err) {
      console.error('操作失败:', err);
      throw err;
    }
  };

  // 删除订单
  const handleDelete = (order: Order) => {
    setDeletingOrder(order);
    setDeleteModalVisible(true);
  };

  // 确认删除
  const handleConfirmDelete = async () => {
    if (!deletingOrder) return;

    try {
      await orderApi.delete(deletingOrder.id);
      setDeleteModalVisible(false);
      setDeletingOrder(null);
      await loadOrders();
    } catch (err) {
      console.error('删除失败:', err);
    }
  };

  // 监听URL参数变化
  useEffect(() => {
    const statusFromUrl = searchParams.get('status') as OrderStatus | null;
    setQueryParams((prev) => ({
      ...prev,
      status: statusFromUrl || undefined,
      page: 1, // 重置到第一页
    }));
  }, [searchParams]);

  // 加载数据
  useEffect(() => {
    loadOrders();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [queryParams.page, queryParams.pageSize, queryParams.status]);

  // 状态格式化
  const formatStatus = (status: OrderStatus): string => {
    const statusMap: Record<OrderStatus, string> = {
      pending: '待处理',
      confirmed: '已确认',
      in_progress: '进行中',
      completed: '已完成',
      canceled: '已取消',
      refunded: '已退款',
    };
    return statusMap[status] || status;
  };

  // 状态颜色
  const getStatusColor = (status: OrderStatus): string => {
    const colorMap: Record<OrderStatus, string> = {
      pending: 'orange',
      confirmed: 'blue',
      in_progress: 'cyan',
      completed: 'green',
      canceled: 'default',
      refunded: 'purple',
    };
    return colorMap[status] || 'default';
  };

  // 表格列定义
  const columns: TableColumn<Order>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: '80px',
    },
    {
      title: '订单标题',
      key: 'title',
      render: (_: unknown, record: Order) => (
        <div className={styles.orderInfo}>
          <div className={styles.orderTitle}>{record.title}</div>
          {record.description && <div className={styles.orderDesc}>{record.description}</div>}
        </div>
      ),
    },
    {
      title: '用户',
      key: 'user',
      width: '150px',
      render: (_: unknown, record: Order) => (
        <div className={styles.userInfo}>
          {record.user ? (
            <>
              <div className={styles.userName}>{record.user.name}</div>
              <div className={styles.userId}>ID: {record.userId}</div>
            </>
          ) : (
            <span>ID: {record.userId}</span>
          )}
        </div>
      ),
    },
    {
      title: '陪玩师',
      key: 'player',
      width: '150px',
      render: (_: unknown, record: Order) => (
        <div className={styles.playerInfo}>
          {record.player ? (
            <>
              <div className={styles.playerName}>{record.player.nickname || '未命名'}</div>
              <div className={styles.playerId}>ID: {record.playerId}</div>
            </>
          ) : record.playerId ? (
            <span>ID: {record.playerId}</span>
          ) : (
            <span className={styles.noData}>未分配</span>
          )}
        </div>
      ),
    },
    {
      title: '游戏',
      key: 'game',
      width: '120px',
      render: (_: unknown, record: Order) => (
        <div>{record.game?.name || `ID: ${record.gameId}`}</div>
      ),
    },
    {
      title: '金额',
      key: 'price',
      width: '100px',
      render: (_: unknown, record: Order) => (
        <div className={styles.price}>{formatCurrency(record.priceCents)}</div>
      ),
    },
    {
      title: '状态',
      key: 'status',
      width: '100px',
      render: (_: unknown, record: Order) => (
        <Tag color={getStatusColor(record.status as OrderStatus) as any}>
          {formatStatus(record.status as OrderStatus)}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      key: 'createdAt',
      width: '160px',
      render: (_: unknown, record: Order) => (
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
      render: (_: unknown, record: Order) => (
        <div className={styles.actions}>
          <Button
            variant="text"
            onClick={() => navigate(`/orders/${record.id}`)}
            className={styles.actionButton}
          >
            详情
          </Button>
          <Button variant="text" onClick={() => handleEdit(record)} className={styles.actionButton}>
            编辑
          </Button>
          <Button
            variant="text"
            onClick={() => handleDelete(record)}
            className={styles.deleteButton}
          >
            删除
          </Button>
        </div>
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
          placeholder="订单标题/ID"
          onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
        />
      ),
    },
    {
      label: '订单状态',
      key: 'status',
      element: (
        <Select
          value={queryParams.status || ''}
          onChange={(value) =>
            setQueryParams((prev) => ({
              ...prev,
              status: value ? (value as OrderStatus) : undefined,
            }))
          }
          options={[
            { label: '全部状态', value: '' },
            { label: '待处理', value: 'pending' },
            { label: '已确认', value: 'confirmed' },
            { label: '进行中', value: 'in_progress' },
            { label: '已完成', value: 'completed' },
            { label: '已取消', value: 'canceled' },
            { label: '已退款', value: 'refunded' },
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
    <>
      <DataTable
        title="订单管理"
        filters={filters}
        filterActions={filterActions}
        columns={columns}
        dataSource={orders}
        loading={loading}
        rowKey="id"
        pagination={{
          current: queryParams.page || 1,
          pageSize: queryParams.pageSize || 10,
          total,
          onChange: handlePageChange,
        }}
      />

      {/* 编辑Modal */}
      <OrderFormModal
        visible={formModalVisible}
        order={editingOrder}
        onClose={() => {
          setFormModalVisible(false);
          setEditingOrder(null);
        }}
        onSubmit={handleFormSubmit}
      />

      {/* 删除确认Modal */}
      <Modal
        visible={deleteModalVisible}
        title="确认删除"
        onClose={() => {
          setDeleteModalVisible(false);
          setDeletingOrder(null);
        }}
        onOk={handleConfirmDelete}
        onCancel={() => {
          setDeleteModalVisible(false);
          setDeletingOrder(null);
        }}
        okText="确定删除"
        cancelText="取消"
        width={400}
      >
        <p>确定要删除订单 "{deletingOrder?.title}" 吗？此操作不可恢复。</p>
      </Modal>
    </>
  );
};
