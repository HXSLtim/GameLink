import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { DataTable, Button, Input, Select, Tag, Modal } from '../../components';
import type { FilterConfig } from '../../components/DataTable';
import type { TableColumn } from '../../components/Table/Table';
import { playerApi } from '../../services/api/user';
import type {
  Player,
  PlayerListQuery,
  CreatePlayerRequest,
  UpdatePlayerRequest,
} from '../../types/user';
import { formatCurrency, formatDateTime, formatRelativeTime } from '../../utils/formatters';
import { PlayerFormModal } from './PlayerFormModal';
import styles from './PlayerList.module.less';

export const PlayerList: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [players, setPlayers] = useState<Player[]>([]);
  const [total, setTotal] = useState(0);

  // 表单Modal状态
  const [formModalVisible, setFormModalVisible] = useState(false);
  const [editingPlayer, setEditingPlayer] = useState<Player | null>(null);

  // 删除确认Modal状态
  const [deleteModalVisible, setDeleteModalVisible] = useState(false);
  const [deletingPlayer, setDeletingPlayer] = useState<Player | null>(null);

  // 查询参数
  const [queryParams, setQueryParams] = useState<PlayerListQuery>({
    page: 1,
    pageSize: 10,
    keyword: '',
    isVerified: undefined,
  });

  // 加载陪玩师列表
  const loadPlayers = async () => {
    setLoading(true);

    try {
      const result = await playerApi.getList({
        page: queryParams.page,
        pageSize: queryParams.pageSize,
        keyword: queryParams.keyword || undefined,
        isVerified: queryParams.isVerified,
      });

      if (result && result.list) {
        setPlayers(result.list);
        setTotal(result.total || 0);
      } else {
        setPlayers([]);
        setTotal(0);
      }
    } catch (err) {
      console.error('加载陪玩师列表失败:', err);
      setPlayers([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  };

  // 搜索
  const handleSearch = async () => {
    setQueryParams((prev) => ({ ...prev, page: 1 }));
    await loadPlayers();
  };

  // 重置
  const handleReset = async () => {
    const resetParams = {
      page: 1,
      pageSize: 10,
      keyword: '',
      isVerified: undefined,
    };
    setQueryParams(resetParams);
    // 使用重置后的参数立即加载数据
    setLoading(true);
    try {
      const result = await playerApi.getList({
        page: resetParams.page,
        pageSize: resetParams.pageSize,
      });
      if (result && result.list) {
        setPlayers(result.list);
        setTotal(result.total || 0);
      } else {
        setPlayers([]);
        setTotal(0);
      }
    } catch (err) {
      console.error('加载陪玩师列表失败:', err);
      setPlayers([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  };

  // 分页变化
  const handlePageChange = (page: number) => {
    setQueryParams((prev) => ({ ...prev, page }));
  };

  // 新增陪玩师
  const handleCreate = () => {
    setEditingPlayer(null);
    setFormModalVisible(true);
  };

  // 编辑陪玩师
  const handleEdit = (player: Player) => {
    setEditingPlayer(player);
    setFormModalVisible(true);
  };

  // 提交表单
  const handleFormSubmit = async (data: CreatePlayerRequest | UpdatePlayerRequest) => {
    try {
      if (editingPlayer) {
        await playerApi.update(editingPlayer.id, data as UpdatePlayerRequest);
      } else {
        await playerApi.create(data as CreatePlayerRequest);
      }
      await loadPlayers();
    } catch (err) {
      console.error('操作失败:', err);
      throw err;
    }
  };

  // 删除陪玩师
  const handleDelete = (player: Player) => {
    setDeletingPlayer(player);
    setDeleteModalVisible(true);
  };

  // 确认删除
  const handleConfirmDelete = async () => {
    if (!deletingPlayer) return;

    try {
      await playerApi.delete(deletingPlayer.id);
      setDeleteModalVisible(false);
      setDeletingPlayer(null);
      await loadPlayers();
    } catch (err) {
      console.error('删除失败:', err);
    }
  };

  // 加载数据
  useEffect(() => {
    loadPlayers();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [queryParams.page, queryParams.pageSize]);

  // 表格列定义
  const columns: TableColumn<Player>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: '80px',
    },
    {
      title: '陪玩师信息',
      key: 'playerInfo',
      render: (_: unknown, record: Player) => (
        <div className={styles.playerInfo}>
          {record.user?.avatarUrl && (
            <img src={record.user.avatarUrl} alt={record.user.name} className={styles.avatar} />
          )}
          {!record.user?.avatarUrl && (
            <div className={styles.avatarPlaceholder}>{record.user?.name.charAt(0) || '?'}</div>
          )}
          <div className={styles.playerDetails}>
            <div className={styles.playerName}>{record.user?.name || '-'}</div>
            <div className={styles.playerContact}>
              {record.user?.phone || '-'} · {record.user?.email || '-'}
            </div>
          </div>
        </div>
      ),
    },
    {
      title: '主游戏',
      key: 'mainGame',
      width: '120px',
      render: (_: unknown, record: Player) => <div>{record.mainGame?.name || '-'}</div>,
    },
    {
      title: '段位',
      dataIndex: 'rank',
      key: 'rank',
      width: '100px',
      render: (rank: string) => rank || '-',
    },
    {
      title: '时薪',
      key: 'hourlyRate',
      width: '100px',
      render: (_: unknown, record: Player) => (
        <div className={styles.rate}>{formatCurrency(record.hourlyRateCents)}</div>
      ),
    },
    {
      title: '评分',
      dataIndex: 'rating',
      key: 'rating',
      width: '80px',
      render: (rating: number) => (
        <div className={styles.rating}>{rating ? `${rating.toFixed(1)} ⭐` : '-'}</div>
      ),
    },
    {
      title: '认证状态',
      key: 'verified',
      width: '100px',
      render: (_: unknown, record: Player) => (
        <Tag color={record.isVerified ? 'green' : 'orange'}>
          {record.isVerified ? '已认证' : '未认证'}
        </Tag>
      ),
    },
    {
      title: '接单状态',
      key: 'available',
      width: '100px',
      render: (_: unknown, record: Player) => (
        <Tag color={record.isAvailable ? 'blue' : 'default'}>
          {record.isAvailable ? '可接单' : '不可接单'}
        </Tag>
      ),
    },
    {
      title: '注册时间',
      key: 'createdAt',
      width: '160px',
      render: (_: unknown, record: Player) => (
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
      render: (_: unknown, record: Player) => (
        <div className={styles.actions}>
          <Button
            variant="text"
            onClick={() => navigate(`/players/${record.id}`)}
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
          placeholder="陪玩师姓名/手机号"
          onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
        />
      ),
    },
    {
      label: '认证状态',
      key: 'isVerified',
      element: (
        <Select
          value={
            queryParams.isVerified === undefined ? '' : queryParams.isVerified ? 'true' : 'false'
          }
          onChange={(value) =>
            setQueryParams((prev) => ({
              ...prev,
              isVerified: value === '' ? undefined : value === 'true',
            }))
          }
          options={[
            { label: '全部状态', value: '' },
            { label: '已认证', value: 'true' },
            { label: '未认证', value: 'false' },
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

  // 头部操作按钮
  const headerActions = (
    <Button variant="primary" onClick={handleCreate}>
      新增陪玩师
    </Button>
  );

  return (
    <>
      <DataTable
        title="陪玩师管理"
        headerActions={headerActions}
        filters={filters}
        filterActions={filterActions}
        columns={columns}
        dataSource={players}
        loading={loading}
        rowKey="id"
        pagination={{
          current: queryParams.page || 1,
          pageSize: queryParams.pageSize || 10,
          total,
          onChange: handlePageChange,
        }}
      />

      {/* 新增/编辑Modal */}
      <PlayerFormModal
        visible={formModalVisible}
        player={editingPlayer}
        onClose={() => {
          setFormModalVisible(false);
          setEditingPlayer(null);
        }}
        onSubmit={handleFormSubmit}
      />

      {/* 删除确认Modal */}
      <Modal
        visible={deleteModalVisible}
        title="确认删除"
        onClose={() => {
          setDeleteModalVisible(false);
          setDeletingPlayer(null);
        }}
        onOk={handleConfirmDelete}
        onCancel={() => {
          setDeleteModalVisible(false);
          setDeletingPlayer(null);
        }}
        okText="确定删除"
        cancelText="取消"
        width={400}
      >
        <p>确定要删除陪玩师 "{deletingPlayer?.user?.name}" 吗？此操作不可恢复。</p>
      </Modal>
    </>
  );
};
