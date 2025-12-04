/**
 * 运营分析页面
 * 展示用户活跃度、留存率、付费分析、转化漏斗等数据
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
  Row,
  Col,
  Card,
  Statistic,
  DatePicker,
  Select,
  Space,
  Table,
  Progress,
  Typography,
  Spin,
  theme,
} from 'antd';
import {
  UserOutlined,
  RiseOutlined,
  FallOutlined,
  DollarOutlined,
  FunnelPlotOutlined,
  TeamOutlined,
  PercentageOutlined,
} from '@ant-design/icons';
import {
  LineChart,
  Line,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  FunnelChart,
  Funnel,
  LabelList,
  Cell,
} from 'recharts';
import { PageContainer } from '@/components';
import { monitorApi } from '@/api/monitor';
import type {
  ActiveUsersData,
  RetentionData,
  PaymentData,
  ConversionFunnel,
} from '@/types/monitor';
import dayjs, { Dayjs } from 'dayjs';

const { RangePicker } = DatePicker;
const { Text } = Typography;

/**
 * 漏斗图颜色
 */
const FUNNEL_COLORS = ['#8884d8', '#83a6ed', '#8dd1e1', '#82ca9d', '#a4de6c', '#ffc658'];

/**
 * 运营分析页面组件
 */
const AnalyticsPage: React.FC = () => {
  const { token } = theme.useToken();

  // 状态
  const [dateRange, setDateRange] = useState<[Dayjs, Dayjs]>([
    dayjs().subtract(30, 'day'),
    dayjs(),
  ]);
  const [granularity, setGranularity] = useState<'day' | 'week' | 'month'>('day');
  const [loading, setLoading] = useState(true);
  const [activeUsers, setActiveUsers] = useState<ActiveUsersData | null>(null);
  const [retention, setRetention] = useState<RetentionData | null>(null);
  const [payment, setPayment] = useState<PaymentData | null>(null);
  const [funnel, setFunnel] = useState<ConversionFunnel | null>(null);

  /**
   * 加载数据
   */
  const loadData = useCallback(async () => {
    setLoading(true);
    try {
      const params = {
        start_date: dateRange[0].format('YYYY-MM-DD'),
        end_date: dateRange[1].format('YYYY-MM-DD'),
        granularity,
      };

      const [activeRes, retentionRes, paymentRes, funnelRes] = await Promise.all([
        monitorApi.getActiveUsers(params),
        monitorApi.getRetention(params),
        monitorApi.getPaymentAnalytics(params),
        monitorApi.getConversionFunnel(params),
      ]);

      if (activeRes.data?.success) setActiveUsers(activeRes.data.data);
      if (retentionRes.data?.success) setRetention(retentionRes.data.data);
      if (paymentRes.data?.success) setPayment(paymentRes.data.data);
      if (funnelRes.data?.success) setFunnel(funnelRes.data.data);
    } catch (error) {
      console.error('加载运营数据失败:', error);
    } finally {
      setLoading(false);
    }
  }, [dateRange, granularity]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  /**
   * 留存矩阵表格列配置
   */
  const retentionColumns = [
    {
      title: '注册日期',
      dataIndex: 'cohortDate',
      key: 'cohortDate',
      width: 120,
    },
    {
      title: '用户数',
      dataIndex: 'cohortSize',
      key: 'cohortSize',
      width: 80,
    },
    {
      title: '第1周',
      key: 'week1',
      width: 80,
      render: (_: unknown, record: { retention: number[] }) => (
        <span style={{ color: getRetentionColor(record.retention[0]) }}>
          {record.retention[0]?.toFixed(1)}%
        </span>
      ),
    },
    {
      title: '第2周',
      key: 'week2',
      width: 80,
      render: (_: unknown, record: { retention: number[] }) => (
        <span style={{ color: getRetentionColor(record.retention[1]) }}>
          {record.retention[1]?.toFixed(1) || '-'}%
        </span>
      ),
    },
    {
      title: '第3周',
      key: 'week3',
      width: 80,
      render: (_: unknown, record: { retention: number[] }) => (
        <span style={{ color: getRetentionColor(record.retention[2]) }}>
          {record.retention[2]?.toFixed(1) || '-'}%
        </span>
      ),
    },
    {
      title: '第4周',
      key: 'week4',
      width: 80,
      render: (_: unknown, record: { retention: number[] }) => (
        <span style={{ color: getRetentionColor(record.retention[3]) }}>
          {record.retention[3]?.toFixed(1) || '-'}%
        </span>
      ),
    },
  ];

  /**
   * 获取留存率颜色
   */
  const getRetentionColor = (rate: number | undefined): string => {
    if (rate === undefined) return token.colorTextDisabled;
    if (rate >= 40) return token.colorSuccess;
    if (rate >= 20) return token.colorWarning;
    return token.colorError;
  };

  return (
    <PageContainer
      title="运营分析"
      subTitle="用户活跃度、留存率、付费分析、转化漏斗"
      extra={
        <Space>
          <RangePicker
            value={dateRange}
            onChange={(dates) => {
              if (dates && dates[0] && dates[1]) {
                setDateRange([dates[0], dates[1]]);
              }
            }}
          />
          <Select
            value={granularity}
            onChange={setGranularity}
            style={{ width: 100 }}
            options={[
              { value: 'day', label: '按天' },
              { value: 'week', label: '按周' },
              { value: 'month', label: '按月' },
            ]}
          />
        </Space>
      }
    >
      <Spin spinning={loading}>
        {/* 活跃用户统计 */}
        <Row gutter={[16, 16]}>
          <Col xs={24} sm={12} lg={6}>
            <Card bordered={false}>
              <Statistic
                title={<Space><UserOutlined /> DAU (日活)</Space>}
                value={activeUsers?.dau || 0}
                valueStyle={{ color: token.colorPrimary }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card bordered={false}>
              <Statistic
                title={<Space><TeamOutlined /> WAU (周活)</Space>}
                value={activeUsers?.wau || 0}
                valueStyle={{ color: token.colorSuccess }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card bordered={false}>
              <Statistic
                title={<Space><TeamOutlined /> MAU (月活)</Space>}
                value={activeUsers?.mau || 0}
                valueStyle={{ color: token.colorWarning }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card bordered={false}>
              <Statistic
                title={<Space><PercentageOutlined /> DAU/MAU 比率</Space>}
                value={activeUsers?.dauMauRatio || 0}
                precision={1}
                suffix="%"
                valueStyle={{ color: token.colorInfo }}
              />
            </Card>
          </Col>
        </Row>

        {/* 活跃用户趋势 */}
        <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
          <Col xs={24} lg={16}>
            <Card title="活跃用户趋势" bordered={false}>
              <div style={{ height: 300 }}>
                <ResponsiveContainer width="100%" height="100%">
                  <LineChart data={activeUsers?.trend || []}>
                    <CartesianGrid strokeDasharray="3 3" stroke={token.colorSplit} />
                    <XAxis dataKey="date" stroke={token.colorTextSecondary} />
                    <YAxis stroke={token.colorTextSecondary} />
                    <Tooltip
                      contentStyle={{ backgroundColor: token.colorBgElevated, borderColor: token.colorBorder }}
                    />
                    <Line
                      type="monotone"
                      dataKey="value"
                      name="活跃用户"
                      stroke={token.colorPrimary}
                      strokeWidth={2}
                      dot={{ fill: token.colorPrimary }}
                    />
                  </LineChart>
                </ResponsiveContainer>
              </div>
            </Card>
          </Col>

          {/* 留存率卡片 */}
          <Col xs={24} lg={8}>
            <Card title="用户留存率" bordered={false}>
              <Space direction="vertical" style={{ width: '100%' }} size="large">
                <div>
                  <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 8 }}>
                    <Text>次日留存 (D1)</Text>
                    <Text strong style={{ color: getRetentionColor(retention?.day1) }}>
                      {retention?.day1?.toFixed(1) || 0}%
                    </Text>
                  </div>
                  <Progress
                    percent={retention?.day1 || 0}
                    showInfo={false}
                    strokeColor={getRetentionColor(retention?.day1)}
                    size="small"
                  />
                </div>
                <div>
                  <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 8 }}>
                    <Text>7日留存 (D7)</Text>
                    <Text strong style={{ color: getRetentionColor(retention?.day7) }}>
                      {retention?.day7?.toFixed(1) || 0}%
                    </Text>
                  </div>
                  <Progress
                    percent={retention?.day7 || 0}
                    showInfo={false}
                    strokeColor={getRetentionColor(retention?.day7)}
                    size="small"
                  />
                </div>
                <div>
                  <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 8 }}>
                    <Text>30日留存 (D30)</Text>
                    <Text strong style={{ color: getRetentionColor(retention?.day30) }}>
                      {retention?.day30?.toFixed(1) || 0}%
                    </Text>
                  </div>
                  <Progress
                    percent={retention?.day30 || 0}
                    showInfo={false}
                    strokeColor={getRetentionColor(retention?.day30)}
                    size="small"
                  />
                </div>
              </Space>
            </Card>
          </Col>
        </Row>

        {/* 付费分析 */}
        <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
          <Col xs={24} sm={12} lg={6}>
            <Card bordered={false}>
              <Statistic
                title={<Space><DollarOutlined /> 付费率</Space>}
                value={payment?.payingRate || 0}
                precision={2}
                suffix="%"
                prefix={payment?.payingRate && payment.payingRate > 5 ? <RiseOutlined /> : <FallOutlined />}
                valueStyle={{ color: payment?.payingRate && payment.payingRate > 5 ? token.colorSuccess : token.colorError }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card bordered={false}>
              <Statistic
                title="ARPU"
                value={payment?.arpu || 0}
                precision={2}
                prefix="¥"
                valueStyle={{ color: token.colorPrimary }}
              />
              <Text type="secondary" style={{ fontSize: 12 }}>
                平均每用户收入
              </Text>
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card bordered={false}>
              <Statistic
                title="ARPPU"
                value={payment?.arppu || 0}
                precision={2}
                prefix="¥"
                valueStyle={{ color: token.colorSuccess }}
              />
              <Text type="secondary" style={{ fontSize: 12 }}>
                平均每付费用户收入
              </Text>
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card bordered={false}>
              <Statistic
                title="平均客单价"
                value={payment?.avgOrderValue || 0}
                precision={2}
                prefix="¥"
                valueStyle={{ color: token.colorWarning }}
              />
            </Card>
          </Col>
        </Row>

        {/* 收入趋势和转化漏斗 */}
        <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
          <Col xs={24} lg={12}>
            <Card title="收入趋势" bordered={false}>
              <div style={{ height: 300 }}>
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={payment?.trend || []}>
                    <CartesianGrid strokeDasharray="3 3" stroke={token.colorSplit} />
                    <XAxis dataKey="date" stroke={token.colorTextSecondary} />
                    <YAxis stroke={token.colorTextSecondary} />
                    <Tooltip
                      contentStyle={{ backgroundColor: token.colorBgElevated, borderColor: token.colorBorder }}
                      formatter={(value: number) => [`¥${value.toFixed(2)}`, '收入']}
                    />
                    <Bar
                      dataKey="value"
                      name="收入"
                      fill={token.colorSuccess}
                      radius={[4, 4, 0, 0]}
                    />
                  </BarChart>
                </ResponsiveContainer>
              </div>
            </Card>
          </Col>

          <Col xs={24} lg={12}>
            <Card
              title={<Space><FunnelPlotOutlined /> 转化漏斗</Space>}
              bordered={false}
            >
              <div style={{ height: 300 }}>
                <ResponsiveContainer width="100%" height="100%">
                  <FunnelChart>
                    <Tooltip
                      contentStyle={{ backgroundColor: token.colorBgElevated, borderColor: token.colorBorder }}
                      formatter={(value: number, name: string, props: { payload?: { rate?: number } }) => [
                        `${value} (${props.payload?.rate?.toFixed(1) || 0}%)`,
                        name,
                      ]}
                    />
                    <Funnel
                      dataKey="value"
                      data={funnel?.steps || []}
                      isAnimationActive
                    >
                      {funnel?.steps?.map((_, index) => (
                        <Cell key={`cell-${index}`} fill={FUNNEL_COLORS[index % FUNNEL_COLORS.length]} />
                      ))}
                      <LabelList
                        position="right"
                        fill={token.colorText}
                        stroke="none"
                        dataKey="name"
                      />
                    </Funnel>
                  </FunnelChart>
                </ResponsiveContainer>
              </div>
            </Card>
          </Col>
        </Row>

        {/* 留存矩阵 */}
        <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
          <Col span={24}>
            <Card title="留存矩阵 (按周)" bordered={false}>
              <Table
                columns={retentionColumns}
                dataSource={retention?.matrix || []}
                rowKey="cohortDate"
                pagination={false}
                size="small"
                scroll={{ x: 600 }}
              />
            </Card>
          </Col>
        </Row>
      </Spin>
    </PageContainer>
  );
};

export default AnalyticsPage;
