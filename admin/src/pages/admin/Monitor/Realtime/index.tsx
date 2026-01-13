/**
 * 实时监控页面
 * 展示系统运行状态、在线用户、订单队列、异常告警等实时数据
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
  Row,
  Col,
  Card,
  Statistic,
  Progress,
  Badge,
  List,
  Tag,
  Button,
  Space,
  Typography,
  Alert as AntAlert,
  theme,
  message,
} from 'antd';
import {
  CloudServerOutlined,
  TeamOutlined,
  ShoppingCartOutlined,
  SyncOutlined,
  ApiOutlined,
  DatabaseOutlined,
  ClockCircleOutlined,
  BellOutlined,
} from '@ant-design/icons';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip as RechartsTooltip, ResponsiveContainer } from 'recharts';
import { PageContainer } from '@/components';
import { useWebSocket } from '@/hooks/useWebSocket';
import { getMonitorWSUrl, monitorApi } from '@/api/monitor';
import type {
  SystemStatus,
  OnlineUsers,
  OrderQueue,
  Alert,
  WSMessage,
} from '@/types/monitor';
import dayjs from 'dayjs';

import { logger } from '@/utils/logger';
const { Text } = Typography;

/**
 * 系统状态颜色映射
 */
const statusColorMap = {
  healthy: '#52c41a',
  degraded: '#faad14',
  critical: '#ff4d4f',
};

/**
 * 告警级别配置
 */
const alertLevelConfig = {
  high: { color: 'error', icon: '🔴', text: '高' },
  medium: { color: 'warning', icon: '🟡', text: '中' },
  low: { color: 'processing', icon: '🟢', text: '低' },
};

/**
 * 格式化运行时间
 */
const formatUptime = (seconds: number): string => {
  const days = Math.floor(seconds / 86400);
  const hours = Math.floor((seconds % 86400) / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);

  if (days > 0) return `${days}天 ${hours}小时`;
  if (hours > 0) return `${hours}小时 ${minutes}分钟`;
  return `${minutes}分钟`;
};

/**
 * 格式化字节数
 * 注意：后端返回的 memoryUsed 和 memoryTotal 单位是字节
 */
const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  const index = Math.min(i, sizes.length - 1);
  return parseFloat((bytes / Math.pow(k, index)).toFixed(2)) + ' ' + sizes[index];
};

/**
 * 实时监控页面组件
 */
const RealtimeMonitor: React.FC = () => {
  const { token } = theme.useToken();

  // 状态
  const [systemStatus, setSystemStatus] = useState<SystemStatus | null>(null);
  const [onlineUsers, setOnlineUsers] = useState<OnlineUsers | null>(null);
  const [orderQueue, setOrderQueue] = useState<OrderQueue | null>(null);
  const [alerts, setAlerts] = useState<Alert[]>([]);
  const [metricsHistory, setMetricsHistory] = useState<{ time: string; memory: number; goroutines: number }[]>([]);
  const [loading, setLoading] = useState(true);
  const [markingRead, setMarkingRead] = useState(false);

  /**
   * WebSocket 消息处理
   */
  const handleWSMessage = useCallback((message: WSMessage) => {
    switch (message.type) {
      case 'system_status':
        const status = message.data as SystemStatus;
        // 调试：打印原始内存值和计算过程
        const usedInKB = status.memoryUsed;
        const totalInKB = status.memoryTotal;
        logger.info('内存原始值 (KB):', {
          memoryUsed: status.memoryUsed,
          memoryTotal: status.memoryTotal,
          usedInKB,
          totalInKB,
          // 假设单位是 KB
          usedGB_fromKB: (usedInKB / 1024 / 1024).toFixed(2),
          totalGB_fromKB: (totalInKB / 1024 / 1024).toFixed(2),
          // 假设单位是字节
          usedGB_fromBytes: (usedInKB / 1024 / 1024 / 1024).toFixed(2),
          totalGB_fromBytes: (totalInKB / 1024 / 1024 / 1024).toFixed(2),
        });
        setSystemStatus(status);
        // 更新历史数据
        setMetricsHistory(prev => {
          const newPoint = {
            time: dayjs().format('HH:mm:ss'),
            memory: (status.memoryUsed || 0) / 1024 / 1024, // 转换为 MB
            goroutines: status.goroutines,
          };
          const updated = [...prev, newPoint];
          return updated.slice(-20); // 保留最近20个点
        });
        break;
      case 'online_users':
        setOnlineUsers(message.data as OnlineUsers);
        break;
      case 'order_queue':
        setOrderQueue(message.data as OrderQueue);
        break;
      case 'alert':
        setAlerts(prev => [message.data as Alert, ...prev].slice(0, 50));
        break;
    }
  }, []);

  // WebSocket 连接
  const { connected } = useWebSocket({
    url: getMonitorWSUrl(),
    onMessage: handleWSMessage,
    autoConnect: true,
    maxRetries: 10,
  });

  /**
   * 初始加载数据
   */
  const loadInitialData = useCallback(async () => {
    setLoading(true);
    try {
      const [statusRes, usersRes, queueRes, alertsRes] = await Promise.all([
        monitorApi.getSystemStatus(),
        monitorApi.getOnlineUsers(),
        monitorApi.getOrderQueue(),
        monitorApi.getAlerts({ page_size: 20 }),
      ]);

      if (statusRes.data?.success) setSystemStatus(statusRes.data.data);
      if (usersRes.data?.success) setOnlineUsers(usersRes.data.data);
      if (queueRes.data?.success) setOrderQueue(queueRes.data.data);
      if (alertsRes.data?.success) setAlerts(alertsRes.data.data || []);
    } catch (error) {
      logger.error('加载监控数据失败:', error);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadInitialData();
  }, [loadInitialData]);

  /**
   * 标记告警已读
   */
  const handleMarkAlertRead = async (alertId: string) => {
    try {
      await monitorApi.markAlertRead(alertId);
      setAlerts(prev => prev.map(a => (a.id === alertId ? { ...a, isRead: true } : a)));
      message.success('标记已读成功');
    } catch (error) {
      logger.error('标记告警已读失败:', error);
      message.error('标记已读失败，请稍后重试');
    }
  };

  /**
   * 标记全部已读
   */
  const handleMarkAllRead = async () => {
    const unreadIds = alerts.filter(a => !a.isRead).map(a => a.id);
    if (unreadIds.length === 0) {
      message.info('没有未读告警');
      return;
    }

    setMarkingRead(true);
    try {
      logger.info('准备标记已读的告警IDs:', unreadIds);
      const response = await monitorApi.markAlertsRead(unreadIds);
      logger.info('标记已读响应:', response);
      setAlerts(prev => prev.map(a => ({ ...a, isRead: true })));
      message.success(`成功标记 ${unreadIds.length} 条告警为已读`);
    } catch (error) {
      logger.error('批量标记已读失败:', error);
      message.error('批量标记已读失败，请稍后重试');
    } finally {
      setMarkingRead(false);
    }
  };

  const unreadCount = alerts.filter(a => !a.isRead).length;

  return (
    <PageContainer
      title="实时监控"
      subTitle={
        <Space>
          <Badge
            status={connected ? 'success' : 'error'}
            text={connected ? '已连接' : '未连接'}
          />
          {systemStatus && (
            <Tag color={statusColorMap[systemStatus.status]}>
              {systemStatus.status === 'healthy' ? '系统正常' :
                systemStatus.status === 'degraded' ? '系统降级' : '系统异常'}
            </Tag>
          )}
        </Space>
      }
      extra={
        <Button
          icon={<SyncOutlined spin={loading} />}
          onClick={loadInitialData}
        >
          刷新
        </Button>
      }
    >
      {/* 系统状态卡片 */}
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={6}>
          <Card variant="borderless" loading={loading}>
            <Statistic
              title={<Space><ApiOutlined /> CPU 使用率</Space>}
              value={systemStatus?.cpuUsage || 0}
              precision={1}
              suffix="%"
              styles={{ content: { color: (systemStatus?.cpuUsage || 0) > 80 ? token.colorError : token.colorPrimary } }}
            />
            <Progress
              percent={systemStatus?.cpuUsage || 0}
              showInfo={false}
              strokeColor={(systemStatus?.cpuUsage || 0) > 80 ? token.colorError : token.colorPrimary}
              size="small"
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card variant="borderless" loading={loading}>
            <Statistic
              title={<Space><CloudServerOutlined /> 内存使用</Space>}
              value={`${formatBytes(systemStatus?.memoryUsed || 0)} / ${formatBytes(systemStatus?.memoryTotal || 0)}`}
              styles={{ content: { color: (systemStatus?.memoryUsage || 0) > 80 ? token.colorError : token.colorWarning } }}
            />
            <Progress
              percent={systemStatus?.memoryUsage || 0}
              showInfo={false}
              strokeColor={(systemStatus?.memoryUsage || 0) > 80 ? token.colorError : token.colorWarning}
              size="small"
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card variant="borderless" loading={loading}>
            <Statistic
              title={<Space><ApiOutlined /> Go协程数</Space>}
              value={systemStatus?.goroutines || 0}
              styles={{ content: { color: token.colorPrimary } }}
            />
            <Text type="secondary" style={{ fontSize: 12 }}>
              请求速率: {systemStatus?.requestsPerSec?.toFixed(1) || 0}/s
            </Text>
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card variant="borderless" loading={loading}>
            <Statistic
              title={<Space><DatabaseOutlined /> 数据库连接</Space>}
              value={systemStatus?.dbConnections?.active || 0}
              suffix={`/ ${systemStatus?.dbConnections?.max || 50}`}
              styles={{ content: { color: token.colorSuccess } }}
            />
            <Text type="secondary" style={{ fontSize: 12 }}>
              空闲: {systemStatus?.dbConnections?.idle || 0} | 运行时间: {formatUptime(systemStatus?.uptime || 0)}
            </Text>
          </Card>
        </Col>
      </Row>

      {/* 监控图表 */}
      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col xs={24} lg={16}>
          <Card title="实时监控趋势" variant="borderless" loading={loading} style={{ height: '100%' }}>
            {!loading && (
              <div style={{ height: 250, width: '100%', minWidth: 0, overflow: 'hidden' }}>
                <ResponsiveContainer width="100%" height={250}>
                  <LineChart data={metricsHistory}>
                    <CartesianGrid strokeDasharray="3 3" stroke={token.colorSplit} />
                    <XAxis dataKey="time" stroke={token.colorTextSecondary} />
                    <YAxis yAxisId="left" stroke={token.colorTextSecondary} />
                    <YAxis yAxisId="right" orientation="right" stroke={token.colorTextSecondary} />
                    <RechartsTooltip
                      contentStyle={{ backgroundColor: token.colorBgElevated, borderColor: token.colorBorder }}
                    />
                    <Line
                      yAxisId="left"
                      type="monotone"
                      dataKey="memory"
                      name="内存使用(MB)"
                      stroke={token.colorWarning}
                      dot={false}
                    />
                    <Line
                      yAxisId="right"
                      type="monotone"
                      dataKey="goroutines"
                      name="协程数"
                      stroke={token.colorPrimary}
                      dot={false}
                    />
                  </LineChart>
                </ResponsiveContainer>
              </div>
            )}
          </Card>
        </Col>

        {/* 在线用户统计 */}
        <Col xs={24} lg={8}>
          <Card title={<Space><TeamOutlined /> 在线用户</Space>} variant="borderless" loading={loading} style={{ height: '100%' }}>
            <Row gutter={16}>
              <Col span={12}>
                <Statistic
                  title="当前在线"
                  value={onlineUsers?.total || 0}
                  styles={{ content: { color: token.colorSuccess } }}
                />
              </Col>
              <Col span={12}>
                <Statistic
                  title="今日峰值"
                  value={onlineUsers?.peak || 0}
                  styles={{ content: { color: token.colorPrimary } }}
                />
              </Col>
            </Row>
            <div style={{ marginTop: 16 }}>
              <Space orientation="vertical" style={{ width: '100%' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                  <Text>普通用户</Text>
                  <Text strong>{onlineUsers?.byRole?.user || 0}</Text>
                </div>
                <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                  <Text>陪玩师</Text>
                  <Text strong>{onlineUsers?.byRole?.player || 0}</Text>
                </div>
                <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                  <Text>管理员</Text>
                  <Text strong>{onlineUsers?.byRole?.admin || 0}</Text>
                </div>
              </Space>
            </div>
          </Card>
        </Col>
      </Row>

      {/* 订单队列和告警 */}
      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        {/* 订单处理队列 */}
        <Col xs={24} lg={8}>
          <Card title={<Space><ShoppingCartOutlined /> 订单处理队列</Space>} variant="borderless" loading={loading} style={{ height: '100%' }}>
            <Space orientation="vertical" style={{ width: '100%' }} size="large">
              <div>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 8 }}>
                  <Text>待处理</Text>
                  <Text strong style={{ color: orderQueue?.hasBacklog ? token.colorError : token.colorText }}>
                    {orderQueue?.pending || 0}
                  </Text>
                </div>
                <Progress
                  percent={Math.min((orderQueue?.pending || 0) / 100 * 100, 100)}
                  showInfo={false}
                  strokeColor={orderQueue?.hasBacklog ? token.colorError : token.colorPrimary}
                  size="small"
                />
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <Text>处理中</Text>
                <Text strong>{orderQueue?.processing || 0}</Text>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <Text>处理速度</Text>
                <Text strong>{orderQueue?.processingSpeed?.toFixed(1) || 0} 单/分钟</Text>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <Text>预计等待</Text>
                <Text strong>
                  {orderQueue?.averageWaitTime
                    ? `${Math.round(orderQueue.averageWaitTime / 60)} 分钟`
                    : '-'}
                </Text>
              </div>
              {orderQueue?.hasBacklog && (
                <AntAlert
                  title="订单积压告警"
                  description="当前待处理订单过多，请关注处理进度"
                  type="warning"
                  showIcon
                />
              )}
            </Space>
          </Card>
        </Col>

        {/* 异常告警列表 */}
        <Col xs={24} lg={16}>
          <Card
            title={
              <Space>
                <BellOutlined />
                异常告警
                {unreadCount > 0 && <Badge count={unreadCount} />}
              </Space>
            }
            variant="borderless"
            extra={
              <Button
                size="small"
                onClick={handleMarkAllRead}
                disabled={unreadCount === 0 || markingRead}
                loading={markingRead}
              >
                全部标记已读
              </Button>
            }
            style={{ height: '100%' }}
          >
            <List
              size="small"
              dataSource={alerts.slice(0, 10)}
              locale={{ emptyText: '暂无告警' }}
              renderItem={(alert) => (
                <List.Item
                  style={{
                    backgroundColor: alert.isRead ? 'transparent' : token.colorBgTextHover,
                    padding: '8px 12px',
                    borderRadius: 4,
                  }}
                  actions={[
                    !alert.isRead && (
                      <Button
                        type="link"
                        size="small"
                        onClick={() => handleMarkAlertRead(alert.id)}
                      >
                        标记已读
                      </Button>
                    ),
                  ]}
                >
                  <List.Item.Meta
                    avatar={
                      <Tag color={alertLevelConfig[alert.level].color}>
                        {alertLevelConfig[alert.level].text}
                      </Tag>
                    }
                    title={
                      <Space>
                        <Text strong={!alert.isRead}>{alert.title}</Text>
                        <Tag>{alert.type === 'system' ? '系统' : alert.type === 'business' ? '业务' : '安全'}</Tag>
                      </Space>
                    }
                    description={
                      <Space>
                        <Text type="secondary" style={{ fontSize: 12 }}>{alert.message}</Text>
                        <Text type="secondary" style={{ fontSize: 12 }}>
                          <ClockCircleOutlined /> {dayjs(alert.createdAt).format('HH:mm:ss')}
                        </Text>
                      </Space>
                    }
                  />
                </List.Item>
              )}
            />
          </Card>
        </Col>
      </Row>
    </PageContainer>
  );
};

export default RealtimeMonitor;
