/**
 * 操作日志展示组件
 */
import React, { useMemo } from 'react';
import { Timeline, Tag, Typography, Space, Empty, Spin, Tooltip, Avatar } from 'antd';
import {
  UserOutlined,
  EditOutlined,
  DeleteOutlined,
  PlusOutlined,
  EyeOutlined,
  SettingOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  SyncOutlined,
  ClockCircleOutlined,
} from '@ant-design/icons';
import styles from './index.module.css';

const { Text, Paragraph } = Typography;

type AuditActionType =
  | 'create'
  | 'update'
  | 'delete'
  | 'view'
  | 'export'
  | 'import'
  | 'approve'
  | 'reject'
  | 'login'
  | 'logout'
  | 'other';

interface AuditLogItem {
  id: number | string;
  action: AuditActionType;
  actionLabel?: string;
  module: string;
  description: string;
  operatorId?: number;
  operatorName: string;
  operatorAvatar?: string;
  ipAddress?: string;
  userAgent?: string;
  changes?: {
    field: string;
    oldValue?: unknown;
    newValue?: unknown;
  }[];
  metadata?: Record<string, unknown>;
  createdAt: string;
}

interface AuditLogProps {
  /**
   * 日志数据
   */
  logs: AuditLogItem[];
  /**
   * 是否加载中
   */
  loading?: boolean;
  /**
   * 是否显示操作人头像
   * @default true
   */
  showAvatar?: boolean;
  /**
   * 是否显示变更详情
   * @default true
   */
  showChanges?: boolean;
  /**
   * 是否显示 IP 地址
   * @default false
   */
  showIP?: boolean;
  /**
   * 最大显示条数
   */
  maxItems?: number;
  /**
   * 空数据文本
   */
  emptyText?: string;
  /**
   * 点击日志项回调
   */
  onItemClick?: (log: AuditLogItem) => void;
  /**
   * 自定义渲染日志项
   */
  renderItem?: (log: AuditLogItem) => React.ReactNode;
  /**
   * 自定义类名
   */
  className?: string;
}

// 操作类型配置
const ACTION_CONFIG: Record<
  AuditActionType,
  { color: string; icon: React.ReactNode; label: string }
> = {
  create: { color: 'green', icon: <PlusOutlined />, label: '创建' },
  update: { color: 'blue', icon: <EditOutlined />, label: '更新' },
  delete: { color: 'red', icon: <DeleteOutlined />, label: '删除' },
  view: { color: 'default', icon: <EyeOutlined />, label: '查看' },
  export: { color: 'purple', icon: <SyncOutlined />, label: '导出' },
  import: { color: 'cyan', icon: <SyncOutlined />, label: '导入' },
  approve: { color: 'green', icon: <CheckCircleOutlined />, label: '审批通过' },
  reject: { color: 'red', icon: <CloseCircleOutlined />, label: '审批拒绝' },
  login: { color: 'blue', icon: <UserOutlined />, label: '登录' },
  logout: { color: 'default', icon: <UserOutlined />, label: '登出' },
  other: { color: 'default', icon: <SettingOutlined />, label: '其他' },
};

// 格式化时间
const formatTime = (dateString: string): string => {
  const date = new Date(dateString);
  const now = new Date();
  const diff = now.getTime() - date.getTime();

  // 小于1分钟
  if (diff < 60000) {
    return '刚刚';
  }
  // 小于1小时
  if (diff < 3600000) {
    return `${Math.floor(diff / 60000)} 分钟前`;
  }
  // 小于24小时
  if (diff < 86400000) {
    return `${Math.floor(diff / 3600000)} 小时前`;
  }
  // 小于7天
  if (diff < 604800000) {
    return `${Math.floor(diff / 86400000)} 天前`;
  }
  // 其他
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');

  if (year === now.getFullYear()) {
    return `${month}-${day} ${hours}:${minutes}`;
  }
  return `${year}-${month}-${day} ${hours}:${minutes}`;
};

// 格式化变更值
const formatValue = (value: unknown): string => {
  if (value === null || value === undefined) return '-';
  if (typeof value === 'boolean') return value ? '是' : '否';
  if (typeof value === 'object') return JSON.stringify(value);
  return String(value);
};

export const AuditLog: React.FC<AuditLogProps> = ({
  logs,
  loading = false,
  showAvatar = true,
  showChanges = true,
  showIP = false,
  maxItems,
  emptyText = '暂无操作日志',
  onItemClick,
  renderItem,
  className,
}) => {
  const displayLogs = useMemo(() => {
    if (maxItems && logs.length > maxItems) {
      return logs.slice(0, maxItems);
    }
    return logs;
  }, [logs, maxItems]);

  if (loading) {
    return (
      <div className={styles.loading}>
        <Spin />
      </div>
    );
  }

  if (!logs || logs.length === 0) {
    return <Empty description={emptyText} />;
  }

  const renderLogItem = (log: AuditLogItem) => {
    if (renderItem) {
      return renderItem(log);
    }

    const actionConfig = ACTION_CONFIG[log.action] || ACTION_CONFIG.other;

    return (
      <div
        className={`${styles.logItem} ${onItemClick ? styles.clickable : ''}`}
        onClick={() => onItemClick?.(log)}
      >
        <div className={styles.logHeader}>
          <Space>
            {showAvatar && (
              <Avatar
                size="small"
                src={log.operatorAvatar}
                icon={<UserOutlined />}
              />
            )}
            <Text strong>{log.operatorName}</Text>
            <Tag color={actionConfig.color}>
              {actionConfig.icon} {log.actionLabel || actionConfig.label}
            </Tag>
            <Text type="secondary">{log.module}</Text>
          </Space>
        </div>

        <Paragraph className={styles.logDescription}>
          {log.description}
        </Paragraph>

        {/* 变更详情 */}
        {showChanges && log.changes && log.changes.length > 0 && (
          <div className={styles.changes}>
            {log.changes.map((change, index) => (
              <div key={index} className={styles.changeItem}>
                <Text type="secondary">{change.field}:</Text>
                <Text delete type="secondary" className={styles.oldValue}>
                  {formatValue(change.oldValue)}
                </Text>
                <Text>→</Text>
                <Text className={styles.newValue}>
                  {formatValue(change.newValue)}
                </Text>
              </div>
            ))}
          </div>
        )}

        <div className={styles.logFooter}>
          <Space>
            <Tooltip title={new Date(log.createdAt).toLocaleString()}>
              <Text type="secondary">
                <ClockCircleOutlined /> {formatTime(log.createdAt)}
              </Text>
            </Tooltip>
            {showIP && log.ipAddress && (
              <Text type="secondary">IP: {log.ipAddress}</Text>
            )}
          </Space>
        </div>
      </div>
    );
  };

  return (
    <div className={`${styles.auditLog} ${className || ''}`}>
      <Timeline
        items={displayLogs.map((log) => {
          const actionConfig = ACTION_CONFIG[log.action] || ACTION_CONFIG.other;
          return {
            key: log.id,
            color: actionConfig.color,
            children: renderLogItem(log),
          };
        })}
      />
    </div>
  );
};

/**
 * 简化版 - 单行日志展示
 */
interface AuditLogLineProps {
  log: AuditLogItem;
  showTime?: boolean;
  className?: string;
}

export const AuditLogLine: React.FC<AuditLogLineProps> = ({
  log,
  showTime = true,
  className,
}) => {
  const actionConfig = ACTION_CONFIG[log.action] || ACTION_CONFIG.other;

  return (
    <div className={`${styles.logLine} ${className || ''}`}>
      <Space>
        <Text strong>{log.operatorName}</Text>
        <Tag color={actionConfig.color} style={{ marginRight: 0 }}>
          {log.actionLabel || actionConfig.label}
        </Tag>
        <Text>{log.description}</Text>
        {showTime && (
          <Text type="secondary">{formatTime(log.createdAt)}</Text>
        )}
      </Space>
    </div>
  );
};

export default AuditLog;
