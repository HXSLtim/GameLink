/**
 * 评价管理 - 统计分析页面
 * 需求: 6.1, 6.2, 6.3, 6.4, 6.5
 */
import React, { useState, useEffect } from 'react';
import {
  Card,
  Row,
  Col,
  Statistic,
  Table,
  Button,
  Space,
  message,
  Spin,
  Select,
  Rate,
  Avatar,
  Typography,
  Progress,
} from 'antd';
import {
  StarOutlined,
  FileTextOutlined,
  CheckCircleOutlined,
  ClockCircleOutlined,
  DownloadOutlined,
  ReloadOutlined,
  UserOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { reviewStatsApi } from '@/api/review';
import type { ReviewStats, ReviewTrend, TopPlayer, GameStats } from '@/types/review';

const { Text } = Typography;

const ReviewStatsPage: React.FC = () => {
  // 状态
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState<ReviewStats | null>(null);
  const [trend, setTrend] = useState<ReviewTrend[]>([]);
  const [topPlayers, setTopPlayers] = useState<TopPlayer[]>([]);
  const [gameStats, setGameStats] = useState<GameStats[]>([]);
  const [rankType, setRankType] = useState<'count' | 'rating'>('count');

  // 加载所有数据
  const fetchAllData = async () => {
    setLoading(true);
    try {
      const [statsRes, trendRes, playersRes, gamesRes] = await Promise.all([
        reviewStatsApi.getStats() as unknown as { success: boolean; data: ReviewStats },
        reviewStatsApi.getTrend({ days: 30 }) as unknown as { success: boolean; data: ReviewTrend[] },
        reviewStatsApi.getTopPlayers({ limit: 10, sortBy: rankType }) as unknown as { success: boolean; data: TopPlayer[] },
        reviewStatsApi.getGameStats() as unknown as { success: boolean; data: GameStats[] },
      ]);

      if (statsRes.success) setStats(statsRes.data);
      if (trendRes.success) setTrend(trendRes.data || []);
      if (playersRes.success) setTopPlayers(playersRes.data || []);
      if (gamesRes.success) setGameStats(gamesRes.data || []);
    } catch {
      message.error('获取统计数据失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchAllData();
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [rankType]);

  // 导出数据
  const handleExport = async (type: 'overview' | 'trend' | 'players' | 'games') => {
    try {
      const response = await reviewStatsApi.exportStats(type);
      // 创建下载链接
      const blob = new Blob([response.data as BlobPart], { type: 'text/csv' });
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `review_${type}_${new Date().toISOString().split('T')[0]}.csv`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      message.success('导出成功');
    } catch {
      message.error('导出失败');
    }
  };

  // 陪玩师排行表格列
  const playerColumns: ColumnsType<TopPlayer> = [
    {
      title: '排名',
      dataIndex: 'rank',
      key: 'rank',
      width: 80,
      render: (rank: number) => (
        <span style={{
          color: rank <= 3 ? '#faad14' : undefined,
          fontWeight: rank <= 3 ? 'bold' : undefined,
        }}>
          {rank}
        </span>
      ),
    },
    {
      title: '陪玩师',
      dataIndex: 'playerName',
      key: 'playerName',
      render: (name: string, record: TopPlayer) => (
        <Space>
          <Avatar src={record.avatarUrl || undefined} icon={<UserOutlined />} />
          <span>{name}</span>
        </Space>
      ),
    },
    {
      title: '评价数',
      dataIndex: 'reviewCount',
      key: 'reviewCount',
      width: 100,
    },
    {
      title: '平均评分',
      dataIndex: 'averageRating',
      key: 'averageRating',
      width: 150,
      render: (rating: number) => (
        <Space>
          <Rate disabled value={rating} style={{ fontSize: 14 }} />
          <Text>{rating.toFixed(1)}</Text>
        </Space>
      ),
    },
  ];

  // 游戏统计表格列
  const gameColumns: ColumnsType<GameStats> = [
    {
      title: '游戏',
      dataIndex: 'gameName',
      key: 'gameName',
      render: (name: string, record: GameStats) => (
        <Space>
          {record.gameIcon && <Avatar src={record.gameIcon} shape="square" />}
          <span>{name}</span>
        </Space>
      ),
    },
    {
      title: '评价数',
      dataIndex: 'reviewCount',
      key: 'reviewCount',
      width: 100,
    },
    {
      title: '平均评分',
      dataIndex: 'averageRating',
      key: 'averageRating',
      width: 150,
      render: (rating: number) => (
        <Space>
          <Rate disabled value={rating} style={{ fontSize: 14 }} />
          <Text>{rating.toFixed(1)}</Text>
        </Space>
      ),
    },
  ];

  if (loading && !stats) {
    return (
      <div style={{ textAlign: 'center', padding: 100 }}>
        <Spin size="large" />
      </div>
    );
  }

  return (
    <div>
      {/* 操作按钮 */}
      <Card style={{ marginBottom: 16 }}>
        <Space>
          <Button icon={<ReloadOutlined />} onClick={fetchAllData} loading={loading}>
            刷新数据
          </Button>
          <Button icon={<DownloadOutlined />} onClick={() => handleExport('overview')}>
            导出概览
          </Button>
          <Button icon={<DownloadOutlined />} onClick={() => handleExport('trend')}>
            导出趋势
          </Button>
          <Button icon={<DownloadOutlined />} onClick={() => handleExport('players')}>
            导出排行
          </Button>
          <Button icon={<DownloadOutlined />} onClick={() => handleExport('games')}>
            导出游戏统计
          </Button>
        </Space>
      </Card>

      {/* 统计概览 */}
      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={6}>
          <Card>
            <Statistic
              title="总评价数"
              value={stats?.totalCount || 0}
              prefix={<FileTextOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="平均评分"
              value={stats?.averageRating || 0}
              precision={2}
              prefix={<StarOutlined />}
              suffix="分"
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="待审核"
              value={stats?.pendingCount || 0}
              prefix={<ClockCircleOutlined />}
              styles={{ content: { color: '#faad14' } }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="已通过"
              value={stats?.approvedCount || 0}
              prefix={<CheckCircleOutlined />}
              styles={{ content: { color: '#52c41a' } }}
            />
          </Card>
        </Col>
      </Row>

      {/* 评分分布 */}
      <Card title="评分分布" style={{ marginBottom: 16 }}>
        <Row gutter={16}>
          {Array.isArray(stats?.ratingDistribution) && stats.ratingDistribution.length > 0 ? (
            stats.ratingDistribution.map(item => (
              <Col span={4} key={item.rating}>
                <div style={{ textAlign: 'center', marginBottom: 8 }}>
                  <Rate disabled value={item.rating} style={{ fontSize: 14 }} />
                </div>
                <Progress
                  percent={item.percentage}
                  format={() => `${item.count}条`}
                  strokeColor={item.rating >= 4 ? '#52c41a' : item.rating >= 3 ? '#faad14' : '#ff4d4f'}
                />
              </Col>
            ))
          ) : (
            <Col span={24}>
              <Text type="secondary">暂无数据</Text>
            </Col>
          )}
        </Row>
      </Card>

      {/* 评价趋势 */}
      <Card title="最近30天评价趋势" style={{ marginBottom: 16 }}>
        {Array.isArray(trend) && trend.length > 0 ? (
          <div style={{ overflowX: 'auto' }}>
            <Table
              dataSource={trend}
              rowKey="date"
              pagination={false}
              size="small"
              columns={[
                { title: '日期', dataIndex: 'date', key: 'date' },
                { title: '评价数', dataIndex: 'count', key: 'count' },
                {
                  title: '平均评分',
                  dataIndex: 'averageRating',
                  key: 'averageRating',
                  render: (v: number) => v?.toFixed(2) || '-',
                },
              ]}
              scroll={{ x: 400 }}
            />
          </div>
        ) : (
          <Text type="secondary">暂无趋势数据</Text>
        )}
      </Card>

      <Row gutter={16}>
        {/* 陪玩师排行 */}
        <Col span={12}>
          <Card
            title="陪玩师排行"
            extra={
              <Select
                value={rankType}
                onChange={setRankType}
                style={{ width: 120 }}
                options={[
                  { value: 'count', label: '按评价数' },
                  { value: 'rating', label: '按评分' },
                ]}
              />
            }
          >
            <Table
              columns={playerColumns}
              dataSource={Array.isArray(topPlayers) ? topPlayers : []}
              rowKey="playerId"
              pagination={false}
              size="small"
            />
          </Card>
        </Col>

        {/* 游戏统计 */}
        <Col span={12}>
          <Card title="游戏评价统计">
            <Table
              columns={gameColumns}
              dataSource={Array.isArray(gameStats) ? gameStats : []}
              rowKey="gameId"
              pagination={false}
              size="small"
            />
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default ReviewStatsPage;
