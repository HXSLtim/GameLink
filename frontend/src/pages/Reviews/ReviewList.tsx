import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { DataTable, Button, Input, Select, Modal, SimpleRating } from '../../components';
import type { FilterConfig } from '../../components/DataTable';
import type { TableColumn } from '../../components/Table/Table';
import { reviewApi } from '../../services/api/review';
import type { Review, ReviewListQuery, UpdateReviewRequest } from '../../types/review';
import { formatDateTime, formatRelativeTime } from '../../utils/formatters';
import { ReviewFormModal } from './ReviewFormModal';
import styles from './ReviewList.module.less';

export const ReviewList: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [reviews, setReviews] = useState<Review[]>([]);
  const [total, setTotal] = useState(0);

  // 表单Modal状态
  const [formModalVisible, setFormModalVisible] = useState(false);
  const [editingReview, setEditingReview] = useState<Review | null>(null);

  // 删除确认Modal状态
  const [deleteModalVisible, setDeleteModalVisible] = useState(false);
  const [deletingReview, setDeletingReview] = useState<Review | null>(null);

  // 查询参数
  const [queryParams, setQueryParams] = useState<ReviewListQuery>({
    page: 1,
    pageSize: 10,
    keyword: '',
    minRating: undefined,
  });

  // 加载评价列表
  const loadReviews = async () => {
    setLoading(true);

    try {
      const result = await reviewApi.getList({
        page: queryParams.page,
        pageSize: queryParams.pageSize,
        keyword: queryParams.keyword || undefined,
        minRating: queryParams.minRating,
      });

      if (result && result.list) {
        setReviews(result.list);
        setTotal(result.total || 0);
      } else {
        setReviews([]);
        setTotal(0);
      }
    } catch (err) {
      console.error('加载评价列表失败:', err);
      setReviews([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  };

  // 搜索
  const handleSearch = async () => {
    setQueryParams((prev) => ({ ...prev, page: 1 }));
    await loadReviews();
  };

  // 重置
  const handleReset = async () => {
    const resetParams = {
      page: 1,
      pageSize: 10,
      keyword: '',
      minRating: undefined,
    };
    setQueryParams(resetParams);
    // 使用重置后的参数立即加载数据
    setLoading(true);
    try {
      const result = await reviewApi.getList({
        page: resetParams.page,
        pageSize: resetParams.pageSize,
      });
      if (result && result.list) {
        setReviews(result.list);
        setTotal(result.total || 0);
      } else {
        setReviews([]);
        setTotal(0);
      }
    } catch (err) {
      console.error('加载评价列表失败:', err);
      setReviews([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  };

  // 分页变化
  const handlePageChange = (page: number) => {
    setQueryParams((prev) => ({ ...prev, page }));
  };

  // 编辑评价
  const handleEdit = (review: Review) => {
    setEditingReview(review);
    setFormModalVisible(true);
  };

  // 提交表单
  const handleFormSubmit = async (data: UpdateReviewRequest) => {
    if (!editingReview) return;

    try {
      await reviewApi.update(editingReview.id, data);
      await loadReviews();
    } catch (err) {
      console.error('操作失败:', err);
      throw err;
    }
  };

  // 删除评价
  const handleDelete = (review: Review) => {
    setDeletingReview(review);
    setDeleteModalVisible(true);
  };

  // 确认删除
  const handleConfirmDelete = async () => {
    if (!deletingReview) return;

    try {
      await reviewApi.delete(deletingReview.id);
      setDeleteModalVisible(false);
      setDeletingReview(null);
      await loadReviews();
    } catch (err) {
      console.error('删除失败:', err);
    }
  };

  // 加载数据
  useEffect(() => {
    loadReviews();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [queryParams.page, queryParams.pageSize]);


  // 表格列定义
  const columns: TableColumn<Review>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: '80px',
    },
    {
      title: '订单',
      key: 'order',
      width: '150px',
      render: (_: unknown, record: Review) => (
        <div className={styles.orderInfo}>
          {record.order ? (
            <>
              <div className={styles.orderTitle}>{record.order.title || '订单'}</div>
              <div className={styles.orderId}>ID: {record.orderId}</div>
            </>
          ) : (
            <span>ID: {record.orderId}</span>
          )}
        </div>
      ),
    },
    {
      title: '评价人',
      key: 'reviewer',
      width: '150px',
      render: (_: unknown, record: Review) => (
        <div className={styles.reviewerInfo}>
          {record.reviewer ? (
            <>
              <div className={styles.reviewerName}>{record.reviewer.name}</div>
              <div className={styles.reviewerId}>ID: {record.reviewerId}</div>
            </>
          ) : (
            <span>ID: {record.reviewerId}</span>
          )}
        </div>
      ),
    },
    {
      title: '陪玩师',
      key: 'player',
      width: '150px',
      render: (_: unknown, record: Review) => (
        <div className={styles.playerInfo}>
          {record.player ? (
            <>
              <div className={styles.playerName}>{record.player.nickname || '未命名'}</div>
              <div className={styles.playerId}>ID: {record.playerId}</div>
            </>
          ) : (
            <span>ID: {record.playerId}</span>
          )}
        </div>
      ),
    },
    {
      title: '评分',
      key: 'rating',
      width: '150px',
      render: (_: unknown, record: Review) => (
        <SimpleRating value={record.rating} size={16} />
      ),
    },
    {
      title: '评价内容',
      dataIndex: 'comment',
      key: 'comment',
      render: (comment: string) => <div className={styles.comment}>{comment || '-'}</div>,
    },
    {
      title: '创建时间',
      key: 'createdAt',
      width: '160px',
      render: (_: unknown, record: Review) => (
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
      render: (_: unknown, record: Review) => (
        <div className={styles.actions}>
          <Button
            variant="text"
            onClick={() => navigate(`/reviews/${record.id}`)}
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
          placeholder="评价内容/订单ID"
          onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
        />
      ),
    },
    {
      label: '最低评分',
      key: 'minRating',
      element: (
        <Select
          value={queryParams.minRating?.toString() || ''}
          onChange={(value) =>
            setQueryParams((prev) => ({
              ...prev,
              minRating: value ? Number(value) : undefined,
            }))
          }
          options={[
            { label: '全部评分', value: '' },
            { label: '1星及以上', value: '1' },
            { label: '2星及以上', value: '2' },
            { label: '3星及以上', value: '3' },
            { label: '4星及以上', value: '4' },
            { label: '5星', value: '5' },
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
        title="评价管理"
        filters={filters}
        filterActions={filterActions}
        columns={columns}
        dataSource={reviews}
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
      <ReviewFormModal
        visible={formModalVisible}
        review={editingReview}
        onClose={() => {
          setFormModalVisible(false);
          setEditingReview(null);
        }}
        onSubmit={handleFormSubmit}
      />

      {/* 删除确认Modal */}
      <Modal
        visible={deleteModalVisible}
        title="确认删除"
        onClose={() => {
          setDeleteModalVisible(false);
          setDeletingReview(null);
        }}
        onOk={handleConfirmDelete}
        onCancel={() => {
          setDeleteModalVisible(false);
          setDeletingReview(null);
        }}
        okText="确定删除"
        cancelText="取消"
        width={400}
      >
        <p>确定要删除这条评价吗？此操作不可恢复。</p>
      </Modal>
    </>
  );
};
