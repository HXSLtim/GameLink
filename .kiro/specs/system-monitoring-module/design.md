# 设计文档 - 系统监控模块

## 概述

系统监控模块负责实时监控 GameLink 平台的系统资源、数据库状态、缓存状态、API性能等关键指标，通过告警机制及时发现和处理系统异常，确保平台的稳定运行。

## 架构

```
┌─────────────────────────────────────────────────────────────┐
│                        前端层 (React)                        │
├─────────────────────────────────────────────────────────────┤
│  资源监控  │  数据库监控  │  缓存监控  │  API监控  │  日志查看  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ↓ HTTP/WebSocket
┌─────────────────────────────────────────────────────────────┐
│                      监控服务层 (Go)                         │
├─────────────────────────────────────────────────────────────┤
│  指标采集  │  告警检测  │  日志聚合  │  健康检查  │  报表生成  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                      数据存储层                              │
├─────────────────────────────────────────────────────────────┤
│  InfluxDB (时序数据)  │  Elasticsearch (日志)  │  Redis (缓存) │
└─────────────────────────────────────────────────────────────┘
```

## 数据模型

### SystemMetrics (系统指标)
```typescript
interface SystemMetrics {
  timestamp: Date;
  cpuUsage: number;          // CPU使用率 (0-100)
  memoryUsage: number;       // 内存使用率 (0-100)
  diskUsage: number;         // 磁盘使用率 (0-100)
  networkIn: number;         // 网络入流量 (MB/s)
  networkOut: number;        // 网络出流量 (MB/s)
}
```

### AlertRule (告警规则)
```typescript
interface AlertRule {
  id: number;
  name: string;
  metric: string;
  operator: ComparisonOperator;
  threshold: number;
  level: AlertLevel;
  isActive: boolean;
  createdAt: Date;
}

enum ComparisonOperator {
  GT = '>',
  LT = '<',
  GTE = '>=',
  LTE = '<=',
  EQ = '=='
}

enum AlertLevel {
  INFO = 'info',
  WARNING = 'warning',
  ERROR = 'error',
  CRITICAL = 'critical'
}
```

## 正确性属性

### 属性 1：监控指标范围约束
*对于任何*监控指标（CPU、内存、磁盘使用率），数值必须在0到100之间
**验证：需求 1.1**

### 属性 2：告警阈值合理性
*对于任何*告警规则，阈值必须在监控指标的有效范围内
**验证：需求 6.3**

### 属性 3：数据更新频率
*对于任何*实时监控数据，两次连续更新之间的时间间隔必须在4.5秒到5.5秒之间
**验证：需求 1.2**

## 测试策略

使用 Vitest + Testing Library 进行单元测试，fast-check 进行属性测试，Playwright 进行E2E测试。
