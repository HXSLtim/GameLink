import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { DataTable, Button, Input, Select, Tag } from '../../components';
import type { FilterConfig } from '../../components/DataTable';
import type { TableColumn } from '../../components/Table/Table';
import { playerApi } from '../../services/api/user';
import type { Player, PlayerListQuery } from '../../types/user';
import { formatCurrency, formatDateTime, formatRelativeTime } from '../../utils/formatters';
import styles from './PlayerList.module.less';

export const PlayerList: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [players, setPlayers] = useState<Player[]>([]);
  const [total, setTotal] = useState(0);

  // 查询参数
  const [queryParams, setQueryParams] = useState<PlayerListQuery>({
    page: 1,
    page_size: 10,
    keyword: '',
    is_verified: undefined,
  });

  // 加载陪玩师列表
  const loadPlayers = async () => {
    setLoading(true);

    try {
      const result = await playerApi.getList({
        page: queryParams.page,
        page_size: queryParams.page_size,
        keyword: queryParams.keyword || undefined,
        is_verified: queryParams.is_verified,
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
  const handleSearch = () => {
    setQueryParams((prev) => ({ ...prev, page: 1 }));
  };

  // 重置
  const handleReset = () => {
    setQueryParams({
      page: 1,
      page_size: 10,
      keyword: '',
      is_verified: undefined,
    });
  };

  // 分页变化
  const handlePageChange = (page: number) => {
    setQueryParams((prev) => ({ ...prev, page }));
  };

  // 加载数据
  useEffect(() => {
    loadPlayers();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [queryParams.page, queryParams.page_size]);

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
          {record.user?.avatar_url && (
            <img 
              src={record.user.avatar_url} 
              alt={record.user.name} 
              className={styles.avatar} 
            />
          )}
          {!record.user?.avatar_url && (
            <div className={styles.avatarPlaceholder}>
              {record.user?.name.charAt(0) || '?'}
            </div>
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
      render: (_: unknown, record: Player) => (
        <div>{record.main_game?.name || '-'}</div>
      ),
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
        <div className={styles.rate}>
          {formatCurrency(record.hourly_rate_cents)}
        </div>
      ),
    },
    {
      title: '评分',
      dataIndex: 'rating',
      key: 'rating',
      width: '80px',
      render: (rating: number) => (
        <div className={styles.rating}>
          {rating ? `${rating.toFixed(1)} ⭐` : '-'}
        </div>
      ),
    },
    {
      title: '认证状态',
      key: 'verified',
      width: '100px',
      render: (_: unknown, record: Player) => (
        <Tag color={record.is_verified ? 'green' : 'orange'}>
          {record.is_verified ? '已认证' : '未认证'}
        </Tag>
      ),
    },
    {
      title: '接单状态',
      key: 'available',
      width: '100px',
      render: (_: unknown, record: Player) => (
        <Tag color={record.is_available ? 'blue' : 'default'}>
          {record.is_available ? '可接单' : '不可接单'}
        </Tag>
      ),
    },
    {
      title: '注册时间',
      key: 'createdAt',
      width: '160px',
      render: (_: unknown, record: Player) => (
        <div className={styles.timeInfo}>
          <div>{formatRelativeTime(record.created_at)}</div>
          <div className={styles.timeDetail}>{formatDateTime(record.created_at)}</div>
        </div>
      ),
    },
    {
      title: '操作',
      key: 'actions',
      width: '120px',
      render: (_: unknown, record: Player) => (
        <Button
          variant="text"
          onClick={() => navigate(`/players/${record.id}`)}
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
          placeholder="陪玩师姓名/手机号"
          onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
        />
      ),
    },
    {
      label: '认证状态',
      key: 'is_verified',
      element: (
        <Select
          value={
            queryParams.is_verified === undefined
              ? ''
              : queryParams.is_verified
                ? 'true'
                : 'false'
          }
          onChange={(value) =>
            setQueryParams((prev) => ({
              ...prev,
              is_verified:
                value === '' ? undefined : value === 'true',
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

  return (
    <DataTable
      title="陪玩师管理"
      filters={filters}
      filterActions={filterActions}
      columns={columns}
      dataSource={players}
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
