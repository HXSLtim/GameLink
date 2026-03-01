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
import type { ApiResponse } from '@/api/admin';
import type {
  ActiveUsersData,
  RetentionData,
  PaymentData,
  ConversionFunnel,
} from '@/types/monitor';
import dayjs, { Dayjs } from 'dayjs';

import { logger } from '@/utils/logger';
const { RangePicker } = DatePicker;
const { Text } = Typography;

/**
 * 漏斗图颜色
 */
const FUNNEL_COLORS = ['#8884d8', '#83a6ed', '#8dd1e1', '#82ca9d', '#a4de6c', '#ffc658'];

/**
 * 图表高度
 */
const CHART_HEIGHT = 320;

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
        monitorApi.getActiveUsers(params) as unknown as Promise<ApiResponse<ActiveUsersData>>,
        monitorApi.getRetention(params) as unknown as Promise<ApiResponse<RetentionData>>,
        monitorApi.getPaymentAnalytics(params) as unknown as Promise<ApiResponse<PaymentData>>,
        monitorApi.getConversionFunnel(params) as unknown as Promise<ApiResponse<ConversionFunnel>>,
      ]);

      if (activeRes.success) setActiveUsers(activeRes.data);
      if (retentionRes.success) setRetention(retentionRes.data);
      if (paymentRes.success) setPayment(paymentRes.data);
      if (funnelRes.success) setFunnel(funnelRes.data);
    } catch (error) {
      logger.error('加载运营数据失败:', error);
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
        <span style={{ color: getRetentionColor(record.retention?.[0]) }}>
          {record.retention?.[0]?.toFixed(1)}%
        </span>
      ),
    },
    {
      title: '第2周',
      key: 'week2',
      width: 80,
      render: (_: unknown, record: { retention: number[] }) => (
        <span style={{ color: getRetentionColor(record.retention?.[1]) }}>
          {record.retention?.[1]?.toFixed(1) || '-'}%
        </span>
      ),
    },
    {
      title: '第3周',
      key: 'week3',
      width: 80,
      render: (_: unknown, record: { retention: number[] }) => (
        <span style={{ color: getRetentionColor(record.retention?.[2]) }}>
          {record.retention?.[2]?.toFixed(1) || '-'}%
        </span>
      ),
    },
    {
      title: '第4周',
      key: 'week4',
      width: 80,
      render: (_: unknown, record: { retention: number[] }) => (
        <span style={{ color: getRetentionColor(record.retention?.[3]) }}>
          {record.retention?.[3]?.toFixed(1) || '-'}%
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
            <Card variant="borderless" style={{ height: '100%' }}>
              <Statistic
                title={<Space><UserOutlined /> DAU (日活)</Space>}
                value={activeUsers?.dau || 0}
                styles={{ content: { color: token.colorPrimary } }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card variant="borderless" style={{ height: '100%' }}>
              <Statistic
                title={<Space><TeamOutlined /> WAU (周活)</Space>}
                value={activeUsers?.wau || 0}
                styles={{ content: { color: token.colorSuccess } }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card variant="borderless" style={{ height: '100%' }}>
              <Statistic
                title={<Space><TeamOutlined /> MAU (月活)</Space>}
                value={activeUsers?.mau || 0}
                styles={{ content: { color: token.colorWarning } }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card variant="borderless" style={{ height: '100%' }}>
              <Statistic
                title={<Space><PercentageOutlined /> DAU/MAU 比率</Space>}
                value={activeUsers?.dauMauRatio || 0}
                precision={1}
                suffix="%"
                styles={{ content: { color: token.colorInfo } }}
              />
            </Card>
          </Col>
        </Row>

        {/* 活跃用户趋势 */}
        <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
          <Col xs={24} lg={16}>
            <Card title="活跃用户趋势" variant="borderless" style={{ height: '100%' }}>
              {!loading && (
                <div style={{ height: CHART_HEIGHT, width: '100%', minWidth: 0, overflow: 'hidden' }}>
                  <ResponsiveContainer width="100%" height={CHART_HEIGHT}>
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
              )}
            </Card>
          </Col>

          {/* 留存率卡片 */}
          <Col xs={24} lg={8}>
            <Card title="用户留存率" variant="borderless" style={{ height: '100%' }}>
              <Space orientation="vertical" style={{ width: '100%' }} size="large">
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
            <Card variant="borderless" style={{ height: '100%' }}>
              <Statistic
                title={<Space><DollarOutlined /> 付费率</Space>}
                value={payment?.payingRate || 0}
                precision={2}
                suffix="%"
                prefix={payment?.payingRate && payment.payingRate > 5 ? <RiseOutlined /> : <FallOutlined />}
                styles={{ content: { color: payment?.payingRate && payment.payingRate > 5 ? token.colorSuccess : token.colorError } }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card variant="borderless" style={{ height: '100%' }}>
              <Statistic
                title="ARPU"
                value={payment?.arpu || 0}
                precision={2}
                prefix="¥"
                styles={{ content: { color: token.colorPrimary } }}
              />
              <Text type="secondary" style={{ fontSize: 12 }}>
                平均每用户收入
              </Text>
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card variant="borderless" style={{ height: '100%' }}>
              <Statistic
                title="ARPPU"
                value={payment?.arppu || 0}
                precision={2}
                prefix="¥"
                styles={{ content: { color: token.colorSuccess } }}
              />
              <Text type="secondary" style={{ fontSize: 12 }}>
                平均每付费用户收入
              </Text>
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card variant="borderless" style={{ height: '100%' }}>
              <Statistic
                title="平均客单价"
                value={payment?.avgOrderValue || 0}
                precision={2}
                prefix="¥"
                styles={{ content: { color: token.colorWarning } }}
              />
            </Card>
          </Col>
        </Row>

        {/* 收入趋势和转化漏斗 */}
        <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
          <Col xs={24} lg={12}>
            <Card title="收入趋势" variant="borderless" style={{ height: '100%' }}>
              <div style={{ height: CHART_HEIGHT, width: '100%', minWidth: 0 }}>
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={payment?.trend || []}>
                    <CartesianGrid strokeDasharray="3 3" stroke={token.colorSplit} />
                    <XAxis dataKey="date" stroke={token.colorTextSecondary} />
                    <YAxis stroke={token.colorTextSecondary} />
                    <Tooltip
                      contentStyle={{ backgroundColor: token.colorBgElevated, borderColor: token.colorBorder }}
                      formatter={(value) => [`¥${Number(value ?? 0).toFixed(2)}`, '收入']}
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
              variant="borderless"
              style={{ height: '100%' }}
            >
              <div style={{ height: CHART_HEIGHT, width: '100%', minWidth: 0 }}>
                <ResponsiveContainer width="100%" height="100%">
                  <FunnelChart>
                    <Tooltip
                      contentStyle={{ backgroundColor: token.colorBgElevated, borderColor: token.colorBorder }}
                      formatter={(value, name, props) => [
                        `${Number(value ?? 0)} (${Number(props?.payload?.rate ?? 0).toFixed(1)}%)`,
                        String(name),
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
            <Card title="留存矩阵 (按周)" variant="borderless">
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
