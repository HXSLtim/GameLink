import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { DataTable, Button, Input, Select, Tag, Modal } from '../../components';
import type { FilterConfig } from '../../components/DataTable';
import type { TableColumn } from '../../components/Table/Table';
import { gameApi } from '../../services/api/game';
import type {
  Game,
  GameListQuery,
  GameCategory,
  CreateGameRequest,
  UpdateGameRequest,
} from '../../types/game';
import { formatDateTime, formatRelativeTime } from '../../utils/formatters';
import { GameFormModal } from './GameFormModal';
import styles from './GameList.module.less';

export const GameList: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [games, setGames] = useState<Game[]>([]);
  const [total, setTotal] = useState(0);

  // 表单Modal状态
  const [formModalVisible, setFormModalVisible] = useState(false);
  const [editingGame, setEditingGame] = useState<Game | null>(null);

  // 删除确认Modal状态
  const [deleteModalVisible, setDeleteModalVisible] = useState(false);
  const [deletingGame, setDeletingGame] = useState<Game | null>(null);

  // 查询参数
  const [queryParams, setQueryParams] = useState<GameListQuery>({
    page: 1,
    pageSize: 10,
    keyword: '',
    category: undefined,
  });

  // 加载游戏列表
  const loadGames = async () => {
    setLoading(true);

    try {
      const result = await gameApi.getList({
        page: queryParams.page,
        pageSize: queryParams.pageSize,
        keyword: queryParams.keyword || undefined,
        category: queryParams.category,
      });

      if (result && result.list) {
        setGames(result.list);
        setTotal(result.total || 0);
      } else {
        setGames([]);
        setTotal(0);
      }
    } catch (err) {
      console.error('加载游戏列表失败:', err);
      setGames([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  };

  // 搜索
  const handleSearch = async () => {
    setQueryParams((prev) => ({ ...prev, page: 1 }));
    await loadGames();
  };

  // 重置
  const handleReset = async () => {
    const resetParams = {
      page: 1,
      pageSize: 10,
      keyword: '',
      category: undefined,
    };
    setQueryParams(resetParams);
    // 使用重置后的参数立即加载数据
    setLoading(true);
    try {
      const result = await gameApi.getList({
        page: resetParams.page,
        pageSize: resetParams.pageSize,
      });
      if (result && result.list) {
        setGames(result.list);
        setTotal(result.total || 0);
      } else {
        setGames([]);
        setTotal(0);
      }
    } catch (err) {
      console.error('加载游戏列表失败:', err);
      setGames([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  };

  // 分页变化
  const handlePageChange = (page: number) => {
    setQueryParams((prev) => ({ ...prev, page }));
  };

  // 新增游戏
  const handleCreate = () => {
    setEditingGame(null);
    setFormModalVisible(true);
  };

  // 编辑游戏
  const handleEdit = (game: Game) => {
    setEditingGame(game);
    setFormModalVisible(true);
  };

  // 提交表单
  const handleFormSubmit = async (data: CreateGameRequest | UpdateGameRequest) => {
    try {
      if (editingGame) {
        await gameApi.update(editingGame.id, data as UpdateGameRequest);
      } else {
        await gameApi.create(data as CreateGameRequest);
      }
      await loadGames();
    } catch (err) {
      console.error('操作失败:', err);
      throw err;
    }
  };

  // 删除游戏
  const handleDelete = (game: Game) => {
    setDeletingGame(game);
    setDeleteModalVisible(true);
  };

  // 确认删除
  const handleConfirmDelete = async () => {
    if (!deletingGame) return;

    try {
      await gameApi.delete(deletingGame.id);
      setDeleteModalVisible(false);
      setDeletingGame(null);
      await loadGames();
    } catch (err) {
      console.error('删除失败:', err);
    }
  };

  // 加载数据
  useEffect(() => {
    loadGames();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [queryParams.page, queryParams.pageSize]);

  // 分类格式化
  const formatCategory = (category: string): string => {
    const categoryMap: Record<string, string> = {
      moba: 'MOBA',
      fps: '射击',
      rpg: '角色扮演',
      strategy: '策略',
      sports: '体育',
      racing: '竞速',
      puzzle: '益智',
      other: '其他',
    };
    return categoryMap[category] || category;
  };

  // 分类颜色
  const getCategoryColor = (category: string): string => {
    const colorMap: Record<string, string> = {
      moba: 'blue',
      fps: 'red',
      rpg: 'purple',
      strategy: 'orange',
      sports: 'green',
      racing: 'cyan',
      puzzle: 'magenta',
      other: 'default',
    };
    return colorMap[category] || 'default';
  };

  // 表格列定义
  const columns: TableColumn<Game>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: '80px',
    },
    {
      title: '游戏信息',
      key: 'gameInfo',
      render: (_: unknown, record: Game) => (
        <div className={styles.gameInfo}>
          {record.iconUrl && (
            <img src={record.iconUrl} alt={record.name} className={styles.gameIcon} />
          )}
          {!record.iconUrl && (
            <div className={styles.iconPlaceholder}>{record.name.charAt(0)}</div>
          )}
          <div className={styles.gameDetails}>
            <div className={styles.gameName}>{record.name}</div>
            <div className={styles.gameKey}>KEY: {record.key}</div>
          </div>
        </div>
      ),
    },
    {
      title: '分类',
      key: 'category',
      width: '120px',
      render: (_: unknown, record: Game) => (
        <Tag color={getCategoryColor(record.category) as any}>
          {formatCategory(record.category)}
        </Tag>
      ),
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      render: (description: string) => (
        <div className={styles.description}>{description || '-'}</div>
      ),
    },
    {
      title: '创建时间',
      key: 'createdAt',
      width: '160px',
      render: (_: unknown, record: Game) => (
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
      render: (_: unknown, record: Game) => (
        <div className={styles.actions}>
          <Button
            variant="text"
            onClick={() => navigate(`/games/${record.id}`)}
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
          placeholder="游戏名称/KEY"
          onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
        />
      ),
    },
    {
      label: '分类',
      key: 'category',
      element: (
        <Select
          value={queryParams.category || ''}
          onChange={(value) =>
            setQueryParams((prev) => ({
              ...prev,
              category: value ? (value as GameCategory) : undefined,
            }))
          }
          options={[
            { label: '全部分类', value: '' },
            { label: 'MOBA', value: 'moba' },
            { label: '射击', value: 'fps' },
            { label: '角色扮演', value: 'rpg' },
            { label: '策略', value: 'strategy' },
            { label: '体育', value: 'sports' },
            { label: '竞速', value: 'racing' },
            { label: '益智', value: 'puzzle' },
            { label: '其他', value: 'other' },
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
      新增游戏
    </Button>
  );

  return (
    <>
      <DataTable
        title="游戏管理"
        headerActions={headerActions}
        filters={filters}
        filterActions={filterActions}
        columns={columns}
        dataSource={games}
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
      <GameFormModal
        visible={formModalVisible}
        game={editingGame}
        onClose={() => {
          setFormModalVisible(false);
          setEditingGame(null);
        }}
        onSubmit={handleFormSubmit}
      />

      {/* 删除确认Modal */}
      <Modal
        visible={deleteModalVisible}
        title="确认删除"
        onClose={() => {
          setDeleteModalVisible(false);
          setDeletingGame(null);
        }}
        onOk={handleConfirmDelete}
        onCancel={() => {
          setDeleteModalVisible(false);
          setDeletingGame(null);
        }}
        okText="确定删除"
        cancelText="取消"
        width={400}
      >
        <p>确定要删除游戏 "{deletingGame?.name}" 吗？此操作不可恢复。</p>
      </Modal>
    </>
  );
};
