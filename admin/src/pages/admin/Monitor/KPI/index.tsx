/**
 * KPI 仪表板页面
 * 展示业务 KPI 指标概览、趋势图表和目标管理
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
  Progress,
  Typography,
  Spin,
  theme,
  Button,
  Table,
  Modal,
  Form,
  InputNumber,
  Tooltip,
  Tag,
  App,
} from 'antd';
import {
  RiseOutlined,
  FallOutlined,
  MinusOutlined,
  DollarOutlined,
  ShoppingCartOutlined,
  UserAddOutlined,
  TeamOutlined,
  AimOutlined,
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  SyncOutlined,
} from '@ant-design/icons';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip as RechartsTooltip,
  ResponsiveContainer,
  ReferenceLine,
} from 'recharts';
import { PageContainer } from '@/components';
import { monitorApi } from '@/api/monitor';
import type { ApiResponse } from '@/api/admin';
import type { KPIQueryParams } from '@/api/monitor';
import type { KPIOverview, KPIMetric, KPITrendPoint, KPITarget } from '@/types/monitor';
import dayjs, { Dayjs } from 'dayjs';

const { RangePicker } = DatePicker;
const { Text } = Typography;

/**
 * 周期选项
 */
const PERIOD_OPTIONS = [
  { value: 'today', label: '今日' },
  { value: 'week', label: '本周' },
  { value: 'month', label: '本月' },
  { value: 'custom', label: '自定义' },
];

/**
 * 对比选项
 */
const COMPARE_OPTIONS = [
  { value: 'mom', label: '环比' },
  { value: 'yoy', label: '同比' },
];

/**
 * 指标配置
 */
const METRIC_CONFIG: Record<string, { label: string; icon: React.ReactNode; unit?: string; precision?: number }> = {
  gmv: { label: 'GMV', icon: <DollarOutlined />, unit: '¥', precision: 2 },
  orders: { label: '订单数', icon: <ShoppingCartOutlined />, precision: 0 },
  newUsers: { label: '新用户', icon: <UserAddOutlined />, precision: 0 },
  newPlayers: { label: '新陪玩师', icon: <TeamOutlined />, precision: 0 },
  dau: { label: '日活用户', icon: <TeamOutlined />, precision: 0 },
  retention: { label: '留存率', icon: <AimOutlined />, unit: '%', precision: 1 },
  repurchaseRate: { label: '复购率', icon: <SyncOutlined />, unit: '%', precision: 1 },
};

/**
 * KPI 仪表板页面组件
 */
const KPIDashboard: React.FC = () => {
  const { message, modal } = App.useApp();
  const { token } = theme.useToken();
  const [form] = Form.useForm();

  // 状态
  const [period, setPeriod] = useState<'today' | 'week' | 'month' | 'custom'>('month');
  const [dateRange, setDateRange] = useState<[Dayjs, Dayjs] | null>(null);
  const [compare, setCompare] = useState<'mom' | 'yoy'>('mom');
  const [loading, setLoading] = useState(true);
  const [overview, setOverview] = useState<KPIOverview | null>(null);
  const [selectedMetric, setSelectedMetric] = useState<string>('gmv');
  const [trendData, setTrendData] = useState<KPITrendPoint[]>([]);
  const [trendLoading, setTrendLoading] = useState(false);
  const [targets, setTargets] = useState<KPITarget[]>([]);
  const [targetModalVisible, setTargetModalVisible] = useState(false);
  const [editingTarget, setEditingTarget] = useState<KPITarget | null>(null);
  const [targetSubmitting, setTargetSubmitting] = useState(false);

  /**
   * 构建查询参数
   */
  const buildQueryParams = useCallback((): KPIQueryParams => {
    const params: KPIQueryParams = {
      period,
      compare,
    };
    if (period === 'custom' && dateRange) {
      params.start_date = dateRange[0].format('YYYY-MM-DD');
      params.end_date = dateRange[1].format('YYYY-MM-DD');
    }
    return params;
  }, [period, dateRange, compare]);

  /**
   * 加载 KPI 概览数据
   */
  const loadOverview = useCallback(async () => {
    setLoading(true);
    try {
      const params = buildQueryParams();
      const res = await monitorApi.getKPIOverview(params) as unknown as ApiResponse<KPIOverview>;
      if (res.success) {
        setOverview(res.data);
      }
    } catch (error) {
      console.error('加载 KPI 概览失败:', error);
      message.error('加载 KPI 概览失败');
    } finally {
      setLoading(false);
    }
  }, [buildQueryParams]);

  /**
   * 加载趋势数据
   */
  const loadTrend = useCallback(async () => {
    setTrendLoading(true);
    try {
      const params = buildQueryParams();
      const res = await monitorApi.getKPITrend(selectedMetric, params) as unknown as ApiResponse<KPITrendPoint[]>;
      if (res.success) {
        setTrendData(res.data || []);
      }
    } catch (error) {
      console.error('加载趋势数据失败:', error);
    } finally {
      setTrendLoading(false);
    }
  }, [selectedMetric, buildQueryParams]);

  /**
   * 加载目标列表
   */
  const loadTargets = useCallback(async () => {
    try {
      const res = await monitorApi.getKPITargets() as unknown as ApiResponse<KPITarget[]>;
      if (res.success) {
        setTargets(res.data || []);
      }
    } catch (error) {
      console.error('加载目标列表失败:', error);
    }
  }, []);

  useEffect(() => {
    loadOverview();
    loadTargets();
  }, [loadOverview, loadTargets]);

  useEffect(() => {
    loadTrend();
  }, [loadTrend]);

  /**
   * 获取趋势图标
   */
  const getTrendIcon = (trend: string) => {
    switch (trend) {
      case 'up':
        return <RiseOutlined style={{ color: token.colorSuccess }} />;
      case 'down':
        return <FallOutlined style={{ color: token.colorError }} />;
      default:
        return <MinusOutlined style={{ color: token.colorTextSecondary }} />;
    }
  };

  /**
   * 获取变化率颜色
   */
  const getChangeColor = (rate: number) => {
    if (rate > 0) return token.colorSuccess;
    if (rate < 0) return token.colorError;
    return token.colorTextSecondary;
  };

  /**
   * 渲染 KPI 卡片
   */
  const renderKPICard = (key: string, metric: KPIMetric) => {
    const config = METRIC_CONFIG[key];
    if (!config) return null;

    const isSelected = selectedMetric === key;

    return (
      <Col xs={24} sm={12} lg={6} key={key}>
        <Card
          variant="borderless"
          hoverable
          style={{
            cursor: 'pointer',
            borderColor: isSelected ? token.colorPrimary : 'transparent',
            borderWidth: 2,
            borderStyle: 'solid',
          }}
          onClick={() => setSelectedMetric(key)}
        >
          <Space orientation="vertical" style={{ width: '100%' }} size="small">
            <Space>
              {config.icon}
              <Text type="secondary">{config.label}</Text>
              {getTrendIcon(metric.trend)}
            </Space>
            <Statistic
              value={metric.currentValue}
              precision={config.precision}
              prefix={config.unit === '¥' ? config.unit : undefined}
              suffix={config.unit === '%' ? config.unit : undefined}
              styles={{ content: { fontSize: 24 } }}
            />
            <Space style={{ width: '100%', justifyContent: 'space-between' }}>
              <Tooltip title={`目标: ${config.unit === '¥' ? '¥' : ''}${metric.targetValue.toLocaleString()}${config.unit === '%' ? '%' : ''}`}>
                <Text type="secondary" style={{ fontSize: 12 }}>
                  完成率 {metric.completionRate.toFixed(1)}%
                </Text>
              </Tooltip>
              <Text style={{ color: getChangeColor(metric.changeRate), fontSize: 12 }}>
                {metric.changeRate > 0 ? '+' : ''}{metric.changeRate.toFixed(1)}%
              </Text>
            </Space>
            <Progress
              percent={Math.min(metric.completionRate, 100)}
              showInfo={false}
              strokeColor={metric.completionRate >= 100 ? token.colorSuccess : token.colorPrimary}
              size="small"
            />
          </Space>
        </Card>
      </Col>
    );
  };

  /**
   * 打开目标编辑弹窗
   */
  const openTargetModal = (target?: KPITarget) => {
    setEditingTarget(target || null);
    if (target) {
      form.setFieldsValue({
        periodType: target.periodType,
        metricName: target.metricName,
        targetValue: target.targetValue,
        dateRange: [dayjs(target.startDate), dayjs(target.endDate)],
      });
    } else {
      form.resetFields();
    }
    setTargetModalVisible(true);
  };

  /**
   * 提交目标表单
   */
  const handleTargetSubmit = async () => {
    try {
      const values = await form.validateFields();
      setTargetSubmitting(true);

      const data = {
        periodType: values.periodType,
        metricName: values.metricName,
        targetValue: values.targetValue,
        startDate: values.dateRange[0].format('YYYY-MM-DD'),
        endDate: values.dateRange[1].format('YYYY-MM-DD'),
      };

      if (editingTarget) {
        await monitorApi.updateKPITarget(editingTarget.id, data);
        message.success('目标更新成功');
      } else {
        await monitorApi.createKPITarget(data);
        message.success('目标创建成功');
      }

      setTargetModalVisible(false);
      loadTargets();
      loadOverview();
    } catch (error) {
      console.error('保存目标失败:', error);
      message.error('保存目标失败');
    } finally {
      setTargetSubmitting(false);
    }
  };

  /**
   * 删除目标
   */
  const handleDeleteTarget = async (id: number) => {
    modal.confirm({
      title: '确认删除',
      content: '确定要删除这个 KPI 目标吗？',
      onOk: async () => {
        try {
          await monitorApi.deleteKPITarget(id);
          message.success('删除成功');
          loadTargets();
          loadOverview();
        } catch (error) {
          console.error('删除目标失败:', error);
          message.error('删除目标失败');
        }
      },
    });
  };

  /**
   * 目标表格列配置
   */
  const targetColumns = [
    {
      title: '指标',
      dataIndex: 'metricName',
      key: 'metricName',
      render: (name: string) => METRIC_CONFIG[name]?.label || name,
    },
    {
      title: '周期类型',
      dataIndex: 'periodType',
      key: 'periodType',
      render: (type: string) => {
        const map: Record<string, string> = {
          daily: '日',
          weekly: '周',
          monthly: '月',
        };
        return <Tag>{map[type] || type}</Tag>;
      },
    },
    {
      title: '目标值',
      dataIndex: 'targetValue',
      key: 'targetValue',
      render: (value: number, record: KPITarget) => {
        const config = METRIC_CONFIG[record.metricName];
        if (config?.unit === '¥') return `¥${value.toLocaleString()}`;
        if (config?.unit === '%') return `${value}%`;
        return value.toLocaleString();
      },
    },
    {
      title: '有效期',
      key: 'period',
      render: (_: unknown, record: KPITarget) => (
        <Text type="secondary">
          {dayjs(record.startDate).format('YYYY-MM-DD')} ~ {dayjs(record.endDate).format('YYYY-MM-DD')}
        </Text>
      ),
    },
    {
      title: '操作',
      key: 'actions',
      width: 120,
      render: (_: unknown, record: KPITarget) => (
        <Space>
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => openTargetModal(record)} />
          <Button type="link" size="small" danger icon={<DeleteOutlined />} onClick={() => handleDeleteTarget(record.id)} />
        </Space>
      ),
    },
  ];

  return (
    <PageContainer
      title="KPI 仪表板"
      subTitle="业务核心指标概览与目标管理"
      extra={
        <Space>
          <Select
            value={period}
            onChange={(val) => {
              setPeriod(val);
              if (val !== 'custom') {
                setDateRange(null);
              }
            }}
            style={{ width: 100 }}
            options={PERIOD_OPTIONS}
          />
          {period === 'custom' && (
            <RangePicker
              value={dateRange}
              onChange={(dates) => {
                if (dates && dates[0] && dates[1]) {
                  setDateRange([dates[0], dates[1]]);
                }
              }}
            />
          )}
          <Select
            value={compare}
            onChange={setCompare}
            style={{ width: 80 }}
            options={COMPARE_OPTIONS}
          />
        </Space>
      }
    >
      <Spin spinning={loading}>
        {/* KPI 概览卡片 */}
        <Row gutter={[16, 16]}>
          {overview && Object.entries(overview).map(([key, metric]) => renderKPICard(key, metric))}
        </Row>

        {/* 趋势图表 */}
        <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
          <Col span={24}>
            <Card
              title={
                <Space>
                  <AimOutlined />
                  <span>{METRIC_CONFIG[selectedMetric]?.label || selectedMetric} 趋势</span>
                </Space>
              }
              variant="borderless"
            >
              <Spin spinning={trendLoading}>
                {!trendLoading && (
                  <div style={{ height: 350, width: '100%', minWidth: 0, overflow: 'hidden' }}>
                    <ResponsiveContainer width="100%" height={350}>
                      <LineChart data={trendData}>
                        <CartesianGrid strokeDasharray="3 3" stroke={token.colorSplit} />
                        <XAxis dataKey="date" stroke={token.colorTextSecondary} />
                        <YAxis stroke={token.colorTextSecondary} />
                        <RechartsTooltip
                          contentStyle={{
                            backgroundColor: token.colorBgElevated,
                            borderColor: token.colorBorder,
                          }}
                          formatter={(value: number, name: string) => {
                            const config = METRIC_CONFIG[selectedMetric];
                            let formatted = value.toLocaleString();
                            if (config?.unit === '¥') formatted = `¥${formatted}`;
                            if (config?.unit === '%') formatted = `${formatted}%`;
                            return [formatted, name === 'value' ? '实际值' : '目标值'];
                          }}
                        />
                        <Line
                          type="monotone"
                          dataKey="value"
                          name="实际值"
                          stroke={token.colorPrimary}
                          strokeWidth={2}
                          dot={{ fill: token.colorPrimary }}
                        />
                        {trendData.some(d => d.target !== undefined) && (
                          <Line
                            type="monotone"
                            dataKey="target"
                            name="目标值"
                            stroke={token.colorWarning}
                            strokeWidth={2}
                            strokeDasharray="5 5"
                            dot={false}
                          />
                        )}
                        {overview && (overview as unknown as Record<string, KPIMetric>)[selectedMetric]?.targetValue && (
                          <ReferenceLine
                            y={(overview as unknown as Record<string, KPIMetric>)[selectedMetric]?.targetValue}
                            stroke={token.colorSuccess}
                            strokeDasharray="3 3"
                            label={{ value: '目标', fill: token.colorSuccess, position: 'right' }}
                          />
                        )}
                      </LineChart>
                    </ResponsiveContainer>
                  </div>
                )}
              </Spin>
            </Card>
          </Col>
        </Row>

        {/* KPI 目标管理 */}
        <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
          <Col span={24}>
            <Card
              title={
                <Space>
                  <AimOutlined />
                  <span>KPI 目标管理</span>
                </Space>
              }
              variant="borderless"
              extra={
                <Button type="primary" icon={<PlusOutlined />} onClick={() => openTargetModal()}>
                  新建目标
                </Button>
              }
            >
              <Table
                columns={targetColumns}
                dataSource={targets}
                rowKey="id"
                pagination={{ pageSize: 5 }}
                size="small"
              />
            </Card>
          </Col>
        </Row>
      </Spin>

      {/* 目标编辑弹窗 */}
      <Modal
        title={editingTarget ? '编辑 KPI 目标' : '新建 KPI 目标'}
        open={targetModalVisible}
        onOk={handleTargetSubmit}
        onCancel={() => setTargetModalVisible(false)}
        confirmLoading={targetSubmitting}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="metricName"
            label="指标"
            rules={[{ required: true, message: '请选择指标' }]}
          >
            <Select
              placeholder="选择指标"
              options={Object.entries(METRIC_CONFIG).map(([key, config]) => ({
                value: key,
                label: config.label,
              }))}
            />
          </Form.Item>
          <Form.Item
            name="periodType"
            label="周期类型"
            rules={[{ required: true, message: '请选择周期类型' }]}
          >
            <Select
              placeholder="选择周期类型"
              options={[
                { value: 'daily', label: '日' },
                { value: 'weekly', label: '周' },
                { value: 'monthly', label: '月' },
              ]}
            />
          </Form.Item>
          <Form.Item
            name="targetValue"
            label="目标值"
            rules={[{ required: true, message: '请输入目标值' }]}
          >
            <InputNumber style={{ width: '100%' }} min={0} placeholder="输入目标值" />
          </Form.Item>
          <Form.Item
            name="dateRange"
            label="有效期"
            rules={[{ required: true, message: '请选择有效期' }]}
          >
            <RangePicker style={{ width: '100%' }} />
          </Form.Item>
        </Form>
      </Modal>
    </PageContainer>
  );
};

export default KPIDashboard;
