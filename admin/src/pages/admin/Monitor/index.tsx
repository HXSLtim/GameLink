/**
 * 实时监控页面
 * 显示系统状态、在线用户、订单队列、警告信息等实时数据
 */
import React, { useEffect, useCallback } from 'react';
import {
  Card,
  Row,
  Col,
  Statistic,
  Progress,
  Tag,
  List,
  Badge,
  Typography,
  Space,
  Button,
  Alert as AntAlert,
} from 'antd';
import {
  CloudServerOutlined,
  TeamOutlined,
  ShoppingOutlined,
  BellOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  CloseCircleOutlined,
  WifiOutlined,
} from '@ant-design/icons';
import { PageContainer } from '@/components';
import { useMonitorStore } from '@/stores/modules/monitorStore';
import { wsManager } from '@/utils/websocket';
import type {
  SystemStatus,
  OnlineUsers,
  OrderQueue,
  Alert,
} from '@/utils/websocket';
import { MessageType } from '@/utils/websocket';
import { logger } from '@/utils/logger';

const { Text } = Typography;

/**
 * 系统状态卡片
 */
const SystemStatusCard: React.FC = () => {
  const { systemStatus, wsConnected } = useMonitorStore();

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'healthy':
        return '#52c41a';
      case 'degraded':
        return '#faad14';
      case 'critical':
        return '#ff4d4f';
      default:
        return '#d9d9d9';
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'healthy':
        return '健康';
      case 'degraded':
        return '降级';
      case 'critical':
        return '严重';
      default:
        return '未知';
    }
  };

  if (!systemStatus) {
    return (
      <Card loading title="系统状态" extra={<WifiOutlined />}>
        <Text type="secondary">等待数据...</Text>
      </Card>
    );
  }

  const memoryPercent = (systemStatus.memoryUsed / systemStatus.memoryTotal) * 100;

  return (
    <Card
      title={
        <Space>
          <CloudServerOutlined />
          <span>系统状态</span>
        </Space>
      }
      extra={
        <Space>
          <Tag color={wsConnected ? 'success' : 'error'}>
            {wsConnected ? '已连接' : '未连接'}
          </Tag>
          <Tag color={getStatusColor(systemStatus.status)}>
            {getStatusText(systemStatus.status)}
          </Tag>
        </Space>
      }
    >
      <Row gutter={[16, 16]}>
        <Col span={6}>
          <Statistic title="CPU 使用率" value={systemStatus.cpuUsage} suffix="%" />
          <Progress
            percent={systemStatus.cpuUsage}
            status={systemStatus.cpuUsage > 80 ? 'exception' : 'active'}
            showInfo={false}
          />
        </Col>
        <Col span={6}>
          <Statistic
            title="内存使用"
            value={(systemStatus.memoryUsed / 1024 / 1024 / 1024).toFixed(2)}
            suffix={`/ ${(systemStatus.memoryTotal / 1024 / 1024 / 1024).toFixed(2)} GB`}
          />
          <Progress
            percent={memoryPercent}
            status={memoryPercent > 80 ? 'exception' : 'active'}
            showInfo={false}
          />
        </Col>
        <Col span={6}>
          <Statistic title="Goroutines" value={systemStatus.goroutines} />
        </Col>
        <Col span={6}>
          <Statistic
            title="数据库连接"
            value={systemStatus.dbConnections.active}
            suffix={`/ ${systemStatus.dbConnections.max}`}
          />
          <Text type="secondary" style={{ fontSize: 12 }}>
            空闲: {systemStatus.dbConnections.idle}
          </Text>
        </Col>
      </Row>
      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col span={12}>
          <Statistic
            title="运行时间"
            value={Math.floor(systemStatus.uptime / 3600)}
            suffix="小时"
          />
        </Col>
        <Col span={12}>
          <Statistic
            title="请求/秒"
            value={systemStatus.requestsPerSec}
            precision={2}
          />
        </Col>
      </Row>
    </Card>
  );
};

/**
 * 在线用户卡片
 */
const OnlineUsersCard: React.FC = () => {
  const { onlineUsers } = useMonitorStore();

  if (!onlineUsers) {
    return (
      <Card loading title={<Space><TeamOutlined /><span>在线用户</span></Space>}>
        <Text type="secondary">等待数据...</Text>
      </Card>
    );
  }

  return (
    <Card
      title={
        <Space>
          <TeamOutlined />
          <span>在线用户</span>
        </Space>
      }
    >
      <Row gutter={16}>
        <Col span={8}>
          <Statistic title="当前在线" value={onlineUsers.total} />
        </Col>
        <Col span={8}>
          <Statistic title="峰值" value={onlineUsers.peak} />
        </Col>
        <Col span={8}>
          <Statistic title="更新时间" value={new Date(onlineUsers.updatedAt).toLocaleTimeString()} />
        </Col>
      </Row>
      {onlineUsers.byRole && Object.keys(onlineUsers.byRole).length > 0 && (
        <div style={{ marginTop: 16 }}>
          <Text strong>按角色分布：</Text>
          <div style={{ marginTop: 8 }}>
            {Object.entries(onlineUsers.byRole).map(([role, count]) => (
              <Tag key={role} style={{ margin: '4px' }}>
                {role}: {count}
              </Tag>
            ))}
          </div>
        </div>
      )}
    </Card>
  );
};

/**
 * 订单队列卡片
 */
const OrderQueueCard: React.FC = () => {
  const { orderQueue } = useMonitorStore();

  if (!orderQueue) {
    return (
      <Card loading title={<Space><ShoppingOutlined /><span>订单队列</span></Space>}>
        <Text type="secondary">等待数据...</Text>
      </Card>
    );
  }

  return (
    <Card
      title={
        <Space>
          <ShoppingOutlined />
          <span>订单队列</span>
        </Space>
      }
    >
      <Row gutter={16}>
        <Col span={6}>
          <Statistic title="待处理" value={orderQueue.pending} />
        </Col>
        <Col span={6}>
          <Statistic title="处理中" value={orderQueue.processing} />
        </Col>
        <Col span={6}>
          <Statistic title="已完成" value={orderQueue.completed} />
        </Col>
        <Col span={6}>
          <Statistic
            title="队列状态"
            value={orderQueue.hasBacklog ? '积压' : '正常'}
            valueStyle={{ color: orderQueue.hasBacklog ? '#ff4d4f' : '#52c41a' }}
          />
        </Col>
      </Row>
      <Row gutter={16} style={{ marginTop: 16 }}>
        <Col span={12}>
          <Statistic
            title="处理速度"
            value={orderQueue.processingSpeed}
            precision={2}
            suffix="单/分钟"
          />
        </Col>
        <Col span={12}>
          <Statistic
            title="平均等待"
            value={Math.floor(orderQueue.averageWaitTime)}
            suffix="秒"
          />
        </Col>
      </Row>
    </Card>
  );
};

/**
 * 警告列表卡片
 */
const AlertsCard: React.FC = () => {
  const { alerts, markAlertAsRead, clearAlert, clearAllAlerts } = useMonitorStore();

  const getAlertIcon = (level: string) => {
    switch (level) {
      case 'high':
        return <CloseCircleOutlined style={{ color: '#ff4d4f' }} />;
      case 'medium':
        return <ExclamationCircleOutlined style={{ color: '#faad14' }} />;
      case 'low':
        return <CheckCircleOutlined style={{ color: '#52c41a' }} />;
      default:
        return <BellOutlined />;
    }
  };

  const getAlertTag = (level: string) => {
    const colorMap = {
      high: 'error',
      medium: 'warning',
      low: 'success',
    };
    return (
      <Tag color={colorMap[level as keyof typeof colorMap] || 'default'}>
        {level.toUpperCase()}
      </Tag>
    );
  };

  return (
    <Card
      title={
        <Space>
          <Badge count={alerts.filter((a) => !a.isRead).length} offset={[10, 0]}>
            <BellOutlined />
          </Badge>
          <span>警告通知</span>
        </Space>
      }
      extra={
        <Space>
          <Button size="small" onClick={clearAllAlerts} disabled={alerts.length === 0}>
            清空全部
          </Button>
        </Space>
      }
    >
      {alerts.length === 0 ? (
        <AntAlert
          message="暂无警告"
          description="系统运行正常，没有新的警告通知"
          type="success"
          showIcon
        />
      ) : (
        <List
          dataSource={alerts}
          renderItem={(alert) => (
            <List.Item
              style={{
                backgroundColor: alert.isRead ? 'transparent' : '#f6f6f6',
                padding: '12px',
                borderRadius: '4px',
                marginBottom: '8px',
              }}
              actions={[
                !alert.isRead && (
                  <Button
                    size="small"
                    type="link"
                    onClick={() => markAlertAsRead(alert.id)}
                  >
                    标记已读
                  </Button>
                ),
                <Button
                  size="small"
                  type="link"
                  danger
                  onClick={() => clearAlert(alert.id)}
                >
                  删除
                </Button>,
              ]}
            >
              <List.Item.Meta
                avatar={getAlertIcon(alert.level)}
                title={
                  <Space>
                    <Text strong={!alert.isRead}>{alert.title}</Text>
                    {getAlertTag(alert.level)}
                    <Tag>{alert.type}</Tag>
                  </Space>
                }
                description={
                  <div>
                    <Text type="secondary">{alert.message}</Text>
                    <br />
                    <Text type="secondary" style={{ fontSize: 12 }}>
                      {alert.source} · {new Date(alert.createdAt).toLocaleString()}
                    </Text>
                  </div>
                }
              />
            </List.Item>
          )}
        />
      )}
    </Card>
  );
};

/**
 * 实时监控主页面
 */
const MonitorPage: React.FC = () => {
  const { setSystemStatus, setOnlineUsers, setOrderQueue, addAlert, setWsConnected } =
    useMonitorStore();

  /**
   * 注册 WebSocket 消息处理器
   */
  const registerHandlers = useCallback(() => {
    // 系统状态消息
    const unsubSystemStatus = wsManager.on(MessageType.SystemStatus, (message) => {
      logger.debug('[Monitor] System status received:', message.data);
      if (message.data) {
        setSystemStatus(message.data as unknown as SystemStatus);
      }
    });

    // 在线用户消息
    const unsubOnlineUsers = wsManager.on(MessageType.OnlineUsers, (message) => {
      logger.debug('[Monitor] Online users received:', message.data);
      if (message.data) {
        setOnlineUsers(message.data as unknown as OnlineUsers);
      }
    });

    // 订单队列消息
    const unsubOrderQueue = wsManager.on(MessageType.OrderQueue, (message) => {
      logger.debug('[Monitor] Order queue received:', message.data);
      if (message.data) {
        setOrderQueue(message.data as unknown as OrderQueue);
      }
    });

    // 警告消息
    const unsubAlert = wsManager.on(MessageType.Alert, (message) => {
      logger.info('[Monitor] Alert received:', message.data);
      if (message.data) {
        addAlert(message.data as unknown as Alert);
      }
    });

    // 设置连接状态监听
    wsManager.setEventListeners({
      onOpen: () => {
        logger.info('[Monitor] WebSocket connected');
        setWsConnected(true);
      },
      onClose: () => {
        logger.warn('[Monitor] WebSocket disconnected');
        setWsConnected(false);
      },
      onError: () => {
        logger.error('[Monitor] WebSocket error');
        setWsConnected(false);
      },
    });

    // 返回清理函数
    return () => {
      unsubSystemStatus();
      unsubOnlineUsers();
      unsubOrderQueue();
      unsubAlert();
    };
  }, [setSystemStatus, setOnlineUsers, setOrderQueue, addAlert, setWsConnected]);

  useEffect(() => {
    const cleanup = registerHandlers();
    return cleanup;
  }, [registerHandlers]);

  return (
    <PageContainer title="实时监控" subTitle="系统运行状态实时数据">
      <Row gutter={[16, 16]}>
        {/* 系统状态 */}
        <Col span={24}>
          <SystemStatusCard />
        </Col>

        {/* 在线用户 & 订单队列 */}
        <Col span={12}>
          <OnlineUsersCard />
        </Col>
        <Col span={12}>
          <OrderQueueCard />
        </Col>

        {/* 警告列表 */}
        <Col span={24}>
          <AlertsCard />
        </Col>
      </Row>
    </PageContainer>
  );
};

export default MonitorPage;
