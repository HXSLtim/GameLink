/**
 * 管理后台仪表盘
 * 支持可折叠区块和用户偏好持久化
 * 性能优化：使用 useMemo 缓存图表数据和列定义
 */
import React, { useState, useEffect, useMemo, memo } from 'react';
import { Row, Col, Card, Table, Tag, Avatar, Typography, Space, App, Select, theme } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
  UserOutlined,
  ShoppingCartOutlined,
  TeamOutlined,
  DollarOutlined,
  ClockCircleOutlined,
  RiseOutlined,
} from '@ant-design/icons';
import { PageContainer, StatCard, CollapsibleSection } from '@/components';
import { useResponsiveChartHeight, useDashboardPreferences } from '@/hooks';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
} from 'recharts';
import { adminApi } from '@/api/admin';
import type { DashboardStats, TrendData, TopPlayer } from '@/api/admin';
import { logger } from '@/utils/logger';
import { chartColors } from '@/theme';
import type { GlobalToken } from 'antd';

const { Text } = Typography;
const { Option } = Select;

const COLORS = chartColors.palette;

const statusMap: Record<string, { color: string; text: string }> = {
  pending: { color: 'gold', text: '待确认' },
  confirmed: { color: 'blue', text: '已确认' },
  in_progress: { color: 'processing', text: '进行中' },
  completed: { color: 'success', text: '已完成' },
  canceled: { color: 'default', text: '已取消' },
  refunded: { color: 'error', text: '已退款' },
};

const paymentStatusMap: Record<string, string> = {
  pending: '待支付',
  paid: '已支付',
  refunded: '已退款',
  failed: '支付失败',
};

interface OrderRecord {
  id: number;
  orderNo: string;
  userId: number;
  totalPriceCents: number;
  status: string;
  createdAt: string;
}

interface ChartData {
  name: string;
  value: number;
  [key: string]: string | number;
}

// 抽取饼图组件并使用 memo 优化
interface PieChartCardProps {
  title: string;
  data: ChartData[];
  loading: boolean;
  chartHeight: number;
  token: GlobalToken;
  ariaLabel: string;
}

const PieChartCard = memo(({ title, data, loading, chartHeight, token, ariaLabel }: PieChartCardProps) => (
  <Card title={title} loading={loading} style={{ border: 'none' }}>
    <div
      style={{ height: chartHeight + 20, width: '100%', minWidth: 1, minHeight: 1, overflow: 'hidden' }}
      role="img"
      aria-label={ariaLabel}
    >
      {data.length > 0 && (
        <ResponsiveContainer width="99%" height={chartHeight} debounce={300}>
          <PieChart>
            <Pie
              data={data}
              cx="50%"
              cy="50%"
              labelLine={false}
              label={({ name, percent }) => `${name} ${((percent || 0) * 100).toFixed(0)}%`}
              outerRadius={100}
              fill="#8884d8"
              dataKey="value"
            >
              {data.map((_, index) => (
                <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
              ))}
            </Pie>
            <Tooltip
              contentStyle={{ backgroundColor: token.colorBgElevated, borderColor: token.colorBorder }}
              itemStyle={{ color: token.colorText }}
            />
            <Legend />
          </PieChart>
        </ResponsiveContainer>
      )}
    </div>
  </Card>
));
PieChartCard.displayName = 'PieChartCard';

// 抽取折线图组件并使用 memo 优化
interface LineChartCardProps {
  title: string;
  data: TrendData[];
  loading: boolean;
  chartHeight: number;
  token: GlobalToken;
  ariaLabel: string;
  dataName: string;
  lineColor: string;
}

const LineChartCard = memo(({ title, data, loading, chartHeight, token, ariaLabel, dataName, lineColor }: LineChartCardProps) => (
  <Card title={title} loading={loading} style={{ border: 'none' }}>
    <div
      style={{ height: chartHeight + 20, width: '100%', minWidth: 1, minHeight: 1, overflow: 'hidden' }}
      role="img"
      aria-label={ariaLabel}
    >
      <ResponsiveContainer width="99%" height={chartHeight} debounce={300}>
        <LineChart data={data}>
          <CartesianGrid strokeDasharray="3 3" stroke={token.colorSplit} />
          <XAxis dataKey="date" stroke={token.colorTextSecondary} />
          <YAxis stroke={token.colorTextSecondary} />
          <Tooltip
            contentStyle={{ backgroundColor: token.colorBgElevated, borderColor: token.colorBorder }}
            itemStyle={{ color: token.colorText }}
          />
          <Legend />
          <Line
            type="monotone"
            dataKey="value"
            name={dataName}
            stroke={lineColor}
            activeDot={{ r: 8 }}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  </Card>
));
LineChartCard.displayName = 'LineChartCard';

const Dashboard: React.FC = () => {
  const { message } = App.useApp();
  const { token } = theme.useToken();
  const [loading, setLoading] = useState(true);
  const [days, setDays] = useState(7);
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [revenueTrend, setRevenueTrend] = useState<TrendData[]>([]);
  const [userGrowth, setUserGrowth] = useState<TrendData[]>([]);
  const [recentOrders, setRecentOrders] = useState<OrderRecord[]>([]);
  const [topPlayers, setTopPlayers] = useState<TopPlayer[]>([]);

  // 响应式图表高度
  const chartHeight = useResponsiveChartHeight({
    mobile: 200,
    tablet: 250,
    desktop: 280,
  });

  // 用户偏好（可折叠区块状态持久化）
  const { isSectionCollapsed, toggleSection } = useDashboardPreferences();

  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      try {
        const [dashboardRes, revenueRes, userGrowthRes, ordersRes, topPlayersRes] = await Promise.all([
          adminApi.getDashboardStats(),
          adminApi.getRevenueTrend({ days }),
          adminApi.getUserGrowth({ days }),
          adminApi.getOrders({ page: 1, page_size: 5 }),
          adminApi.getTopPlayers({ limit: 5 }),
        ]);

        setStats(dashboardRes.data?.data || dashboardRes.data || {});
        setRevenueTrend(revenueRes.data?.data || revenueRes.data || []);
        setUserGrowth(userGrowthRes.data?.data || userGrowthRes.data || []);

        const ordersData = ordersRes.data?.data || ordersRes.data || [];
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        setRecentOrders(Array.isArray(ordersData) ? ordersData : (ordersData as any).list || []);

        setTopPlayers(topPlayersRes.data?.data || topPlayersRes.data || []);
      } catch (error) {
        logger.error('Failed to load dashboard data', error);
        message.error('加载仪表盘数据失败');
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, [days, message]);

  // 使用 useMemo 缓存订单列定义
  const orderColumns: ColumnsType<OrderRecord> = useMemo(() => [
    {
      title: '订单号',
      dataIndex: 'orderNo',
      key: 'orderNo',
      render: (text) => <Text copyable={{ text }}>{text}</Text>,
    },
    {
      title: '用户',
      dataIndex: 'userId',
      key: 'userId',
      render: (userId) => (
        <Space>
          <Avatar style={{ backgroundColor: token.colorPrimary }} icon={<UserOutlined />} />
          <Text>用户 #{userId}</Text>
        </Space>
      ),
    },
    {
      title: '金额',
      dataIndex: 'totalPriceCents',
      key: 'totalPriceCents',
      render: (price) => <Text strong>¥{(price / 100).toFixed(2)}</Text>,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status) => (
        <Tag color={statusMap[status]?.color || 'default'}>{statusMap[status]?.text || status}</Tag>
      ),
    },
    {
      title: '时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (text) => (
        <Text type="secondary" style={{ fontSize: 12 }}>
          {new Date(text).toLocaleString()}
        </Text>
      ),
    },
  ], [token.colorPrimary]);

  // 使用 useMemo 缓存热门陪玩列定义
  const topPlayersColumns: ColumnsType<TopPlayer> = useMemo(() => [
    {
      title: '排名',
      key: 'rank',
      width: 60,
      render: (_: unknown, __: TopPlayer, index: number) => <Tag color="gold">#{index + 1}</Tag>,
    },
    {
      title: '陪玩',
      key: 'player',
      render: (_: unknown, record: TopPlayer, index: number) => (
        <Space>
          <Avatar style={{ backgroundColor: `hsl(${index * 60}, 70%, 50%)` }} icon={<UserOutlined />} />
          <Text>{record.nickname}</Text>
        </Space>
      ),
    },
    {
      title: '评分',
      key: 'rating',
      render: (_: unknown, record: TopPlayer) => (
        <Space>
          <Text strong>{record.ratingAverage}</Text>
          <Text type="secondary">({record.ratingCount})</Text>
        </Space>
      ),
    },
  ], []);

  // 使用 useMemo 缓存图表数据计算
  const orderStatusData = useMemo(() => 
    stats?.ordersByStatus
      ? Object.entries(stats.ordersByStatus).map(([name, value]) => ({
          name: statusMap[name]?.text || name,
          value,
        }))
      : [],
    [stats?.ordersByStatus]
  );

  const paymentStatusData = useMemo(() =>
    stats?.paymentsByStatus
      ? Object.entries(stats.paymentsByStatus).map(([name, value]) => ({
          name: paymentStatusMap[name] || name,
          value,
        }))
      : [],
    [stats?.paymentsByStatus]
  );

  // 计算图表 aria-label
  const orderChartAriaLabel = useMemo(() => 
    `订单状态分布饼图，共${orderStatusData.reduce((sum, d) => sum + d.value, 0)}个订单`,
    [orderStatusData]
  );

  const paymentChartAriaLabel = useMemo(() =>
    `支付状态分布饼图，共${paymentStatusData.reduce((sum, d) => sum + d.value, 0)}笔支付`,
    [paymentStatusData]
  );

  return (
    <PageContainer
      title="仪表盘"
      subTitle="系统性能概览"
      extra={
        <Select value={days} onChange={setDays} style={{ width: 120 }}>
          <Option value={7}>近7天</Option>
          <Option value={30}>近30天</Option>
          <Option value={90}>近90天</Option>
        </Select>
      }
    >
      {/* Stats Cards - 始终显示 */}
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={6}>
          <StatCard
            title="总订单数"
            value={stats?.totalOrders || 0}
            icon={<ShoppingCartOutlined />}
            iconBgColor={token.colorSuccess}
            loading={loading}
            suffix=""
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StatCard
            title="交易总额"
            value={(stats?.totalPaidAmountCents || 0) / 100}
            icon={<DollarOutlined />}
            iconBgColor={token.colorWarning}
            loading={loading}
            prefix="¥"
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StatCard
            title="总用户数"
            value={stats?.totalUsers || 0}
            icon={<UserOutlined />}
            iconBgColor={token.colorPrimary}
            loading={loading}
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StatCard
            title="总陪玩数"
            value={stats?.totalPlayers || 0}
            icon={<TeamOutlined />}
            iconBgColor={chartColors.players}
            loading={loading}
          />
        </Col>
      </Row>

      {/* 状态分布图表 - 可折叠 */}
      <CollapsibleSection
        id="status-charts"
        title="状态分布"
        collapsed={isSectionCollapsed('status-charts')}
        onCollapsedChange={() => toggleSection('status-charts')}
      >
        <Row gutter={[16, 16]}>
          <Col xs={24} lg={12}>
            <PieChartCard
              title="订单状态分布"
              data={orderStatusData}
              loading={loading}
              chartHeight={chartHeight}
              token={token}
              ariaLabel={orderChartAriaLabel}
            />
          </Col>
          <Col xs={24} lg={12}>
            <PieChartCard
              title="支付状态分布"
              data={paymentStatusData}
              loading={loading}
              chartHeight={chartHeight}
              token={token}
              ariaLabel={paymentChartAriaLabel}
            />
          </Col>
        </Row>
      </CollapsibleSection>

      {/* 趋势图表 - 可折叠 */}
      <CollapsibleSection
        id="trend-charts"
        title="趋势分析"
        collapsed={isSectionCollapsed('trend-charts')}
        onCollapsedChange={() => toggleSection('trend-charts')}
      >
        <Row gutter={[16, 16]}>
          <Col xs={24} lg={12}>
            <LineChartCard
              title={`收入趋势 (近${days}天)`}
              data={revenueTrend}
              loading={loading}
              chartHeight={chartHeight}
              token={token}
              ariaLabel={`收入趋势折线图，展示近${days}天收入变化`}
              dataName="收入"
              lineColor={chartColors.revenue}
            />
          </Col>
          <Col xs={24} lg={12}>
            <LineChartCard
              title={`用户增长 (近${days}天)`}
              data={userGrowth}
              loading={loading}
              chartHeight={chartHeight}
              token={token}
              ariaLabel={`用户增长折线图，展示近${days}天新增用户变化`}
              dataName="新增用户"
              lineColor={chartColors.users}
            />
          </Col>
        </Row>
      </CollapsibleSection>

      {/* 数据列表 - 可折叠 */}
      <CollapsibleSection
        id="data-tables"
        title="数据详情"
        collapsed={isSectionCollapsed('data-tables')}
        onCollapsedChange={() => toggleSection('data-tables')}
      >
        <Row gutter={[16, 16]} style={{ display: 'flex', flexWrap: 'wrap', alignItems: 'stretch' }}>
          <Col xs={24} xl={16} style={{ display: 'flex', flexDirection: 'column' }}>
            <Card
              title={
                <Space>
                  <ClockCircleOutlined />
                  <span>最新订单</span>
                </Space>
              }
              extra={<a href="/admin/biz/order">查看全部</a>}
              loading={loading}
              style={{ flex: 1, width: '100%', border: 'none' }}
            >
              <Table columns={orderColumns} dataSource={recentOrders} rowKey="id" pagination={false} size="small" />
            </Card>
          </Col>
          <Col xs={24} xl={8} style={{ display: 'flex', flexDirection: 'column' }}>
            <Card
              title={
                <Space>
                  <RiseOutlined />
                  <span>热门陪玩</span>
                </Space>
              }
              loading={loading}
              style={{ flex: 1, width: '100%', border: 'none' }}
            >
              <Table
                columns={topPlayersColumns}
                dataSource={topPlayers}
                rowKey="playerId"
                pagination={false}
                size="small"
              />
            </Card>
          </Col>
        </Row>
      </CollapsibleSection>
    </PageContainer>
  );
};

export default Dashboard;
