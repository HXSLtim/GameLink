# 仪表盘模块增强实施计划

## 一、项目概述

### 需求目标
为 GameLink 管理后台新增3个仪表盘相关页面：
1. **实时监控** - 系统实时状态监控（WebSocket推送）
2. **运营分析** - 运营数据分析报告
3. **KPI仪表板** - 关键业务指标展示

### 技术选型确认
- 实时方案：WebSocket 实时推送
- KPI周期：支持 日/周/月 快捷选择 + 自定义日期范围
- 权限控制：需要（不同角色看到不同内容）

---

## 二、页面功能详细设计

### 2.1 实时监控页面 (`/admin/monitor/realtime`)

#### 功能模块
| 模块 | 功能 | 数据来源 |
|------|------|----------|
| 系统状态 | CPU、内存、Go协程数、数据库连接 | WebSocket 实时推送 |
| 在线用户 | 当前在线人数、峰值、分布 | WebSocket + Redis |
| 订单处理队列 | 待处理订单数、处理速度、积压告警 | WebSocket |
| 异常告警 | 系统异常、业务异常、安全告警 | WebSocket + 告警服务 |

#### UI布局
```
┌──────────────────────────────────────────────────────┐
│  系统运行状态                           🟢 正常运行    │
├──────────────┬──────────────┬──────────────┬─────────┤
│  CPU使用率   │  内存使用率   │  Go协程数    │ DB连接池 │
│    45%       │    62%       │   1,234     │  15/50  │
├──────────────┴──────────────┴──────────────┴─────────┤
│  实时监控图表（折线图 - 最近10分钟）                    │
│  [CPU] [内存] [请求数]                                │
├──────────────────────────────────────────────────────┤
│  在线用户统计                                         │
│  ┌─────────┬─────────┬─────────┬─────────┐          │
│  │ 当前在线 │ 今日峰值 │ 普通用户 │ 陪玩师  │          │
│  │   523   │   892   │   412   │   111   │          │
│  └─────────┴─────────┴─────────┴─────────┘          │
├──────────────────────────────────────────────────────┤
│  订单处理队列                                         │
│  待处理: 23 | 处理中: 8 | 处理速度: 45/min           │
│  [进度条显示队列状态]                                  │
├──────────────────────────────────────────────────────┤
│  异常告警列表                          [全部标记已读]  │
│  🔴 [高] 数据库连接超时 - 2024-01-15 14:32:01        │
│  🟡 [中] API响应超时 /api/v1/orders - 14:30:22       │
│  🟢 [低] 用户登录异常尝试 - 14:28:15                  │
└──────────────────────────────────────────────────────┘
```

### 2.2 运营分析页面 (`/admin/monitor/analytics`)

#### 功能模块
| 模块 | 指标 | 可视化 |
|------|------|--------|
| 日活/月活 | DAU/MAU/WAU、活跃趋势 | 折线图 + 数值卡片 |
| 用户留存 | 次日留存、7日留存、30日留存 | 留存矩阵热力图 |
| 付费分析 | 付费率、ARPU、ARPPU | 柱状图 + 环形图 |
| 订单分析 | 客单价、订单转化漏斗 | 漏斗图 + 趋势图 |
| 陪玩师分析 | 活跃度、接单率、好评率 | 排行榜 + 散点图 |

#### UI布局
```
┌──────────────────────────────────────────────────────┐
│  运营数据分析              [日期范围选择器] [导出报告]  │
├──────────────────────────────────────────────────────┤
│  用户活跃概览                                         │
│  ┌─────────┬─────────┬─────────┬─────────┐          │
│  │  DAU    │  WAU    │  MAU    │ DAU/MAU │          │
│  │  1,234  │  5,678  │ 12,345  │  10.0%  │          │
│  │  ↑12%   │  ↑8%    │  ↑15%   │  ↓2%    │          │
│  └─────────┴─────────┴─────────┴─────────┘          │
├─────────────────────────┬────────────────────────────┤
│  用户活跃趋势            │  用户留存矩阵              │
│  [折线图:DAU/WAU/MAU]   │  [热力图:留存率]           │
├─────────────────────────┴────────────────────────────┤
│  付费数据分析                                         │
│  ┌─────────┬─────────┬─────────┬─────────┐          │
│  │  付费率  │  ARPU   │  ARPPU  │  客单价  │          │
│  │  5.2%   │  ¥12.5  │  ¥240   │  ¥89    │          │
│  └─────────┴─────────┴─────────┴─────────┘          │
├─────────────────────────┬────────────────────────────┤
│  订单转化漏斗            │  收入趋势                  │
│  [漏斗图]               │  [面积图]                  │
└─────────────────────────┴────────────────────────────┘
```

### 2.3 KPI仪表板页面 (`/admin/monitor/kpi`)

#### 功能模块
| KPI分类 | 指标 | 目标值 |
|---------|------|--------|
| 增长指标 | 新增用户、新增陪玩师、订单增长率 | 可配置 |
| 转化指标 | 注册转化率、下单转化率、复购率 | 可配置 |
| 质量指标 | 好评率、投诉率、退款率 | 可配置 |
| 财务指标 | GMV、收入、客单价 | 可配置 |

#### UI布局
```
┌──────────────────────────────────────────────────────┐
│  KPI仪表板     [今日|本周|本月|自定义] [对比:环比|同比] │
├──────────────────────────────────────────────────────┤
│  核心KPI概览                                          │
│  ┌────────────────┬────────────────┬────────────────┐│
│  │   GMV完成率    │  订单完成率     │  用户增长率    ││
│  │     85%       │     92%        │     78%       ││
│  │  [仪表盘图]    │  [仪表盘图]     │  [仪表盘图]    ││
│  │  目标:100万   │  目标:1000单    │  目标:500人   ││
│  └────────────────┴────────────────┴────────────────┘│
├──────────────────────────────────────────────────────┤
│  增长指标                                             │
│  ┌─────────┬─────────┬─────────┬─────────┐          │
│  │ 新增用户 │新增陪玩师│ 订单数   │ 复购用户 │          │
│  │  +234   │  +12    │  +567   │  +89    │          │
│  │  ↑15%   │  ↓5%    │  ↑23%   │  ↑8%    │          │
│  └─────────┴─────────┴─────────┴─────────┘          │
├─────────────────────────┬────────────────────────────┤
│  KPI趋势对比            │  目标完成进度              │
│  [多系列折线图]          │  [环形进度图列表]          │
├─────────────────────────┴────────────────────────────┤
│  KPI明细表格                                          │
│  [可展开的分类KPI表格，支持钻取查看详情]                │
└──────────────────────────────────────────────────────┘
```

---

## 三、技术实施方案

### 3.1 前端实现

#### 新增文件结构
```
frontend/src/
├── pages/admin/Monitor/
│   ├── Realtime/
│   │   └── index.tsx           # 实时监控页面
│   ├── Analytics/
│   │   └── index.tsx           # 运营分析页面
│   └── KPI/
│       └── index.tsx           # KPI仪表板页面
├── api/
│   └── monitor.ts              # 监控相关API
├── hooks/
│   └── useWebSocket.ts         # WebSocket Hook
├── components/Monitor/
│   ├── SystemStatusCard.tsx    # 系统状态卡片
│   ├── OnlineUsersCard.tsx     # 在线用户卡片
│   ├── AlertList.tsx           # 告警列表
│   ├── GaugeChart.tsx          # 仪表盘图表
│   ├── RetentionMatrix.tsx     # 留存矩阵
│   └── FunnelChart.tsx         # 漏斗图
└── types/
    └── monitor.ts              # 监控相关类型定义
```

#### WebSocket Hook 设计
```typescript
// useWebSocket.ts
interface UseWebSocketOptions {
  url: string;
  onMessage?: (data: any) => void;
  onError?: (error: Event) => void;
  reconnectInterval?: number;
  maxRetries?: number;
}

function useWebSocket(options: UseWebSocketOptions) {
  // 自动连接、断线重连、消息处理
  return {
    connected: boolean,
    send: (data: any) => void,
    disconnect: () => void
  };
}
```

### 3.2 后端实现

#### 新增文件结构
```
backend/
├── internal/handler/admin/
│   └── monitor.go              # 监控API处理器
├── internal/service/
│   └── monitor/
│       ├── realtime.go         # 实时监控服务
│       ├── analytics.go        # 运营分析服务
│       └── kpi.go              # KPI服务
├── internal/repository/
│   └── monitor/
│       └── repository.go       # 监控数据仓储
├── internal/ws/
│   ├── hub.go                  # WebSocket Hub
│   ├── client.go               # WebSocket Client
│   └── handler.go              # WebSocket Handler
└── internal/model/
    └── alert.go                # 告警模型
```

#### WebSocket 架构
```
┌─────────────────────────────────────────────┐
│                 WebSocket Hub               │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐      │
│  │ Client1 │ │ Client2 │ │ Client3 │      │
│  └─────────┘ └─────────┘ └─────────┘      │
└─────────────────────────────────────────────┘
         ↑                    ↑
         │                    │
┌─────────────────┐  ┌─────────────────┐
│  System Monitor │  │  Alert Service  │
│  (定时采集)      │  │  (事件触发)     │
└─────────────────┘  └─────────────────┘
```

### 3.3 API 设计

#### REST API 端点
```
# 实时监控
GET  /api/v1/admin/monitor/system-status    # 系统状态快照
GET  /api/v1/admin/monitor/online-users     # 在线用户统计
GET  /api/v1/admin/monitor/order-queue      # 订单队列状态
GET  /api/v1/admin/monitor/alerts           # 告警列表
PUT  /api/v1/admin/monitor/alerts/:id/read  # 标记告警已读

# 运营分析
GET  /api/v1/admin/analytics/active-users   # 活跃用户数据
GET  /api/v1/admin/analytics/retention      # 用户留存数据
GET  /api/v1/admin/analytics/payment        # 付费分析数据
GET  /api/v1/admin/analytics/conversion     # 转化漏斗数据

# KPI
GET  /api/v1/admin/kpi/overview             # KPI概览
GET  /api/v1/admin/kpi/trend                # KPI趋势
GET  /api/v1/admin/kpi/targets              # KPI目标配置
PUT  /api/v1/admin/kpi/targets              # 更新KPI目标
```

#### WebSocket 端点
```
WS   /api/v1/admin/ws/monitor               # 实时监控WebSocket
```

#### WebSocket 消息格式
```typescript
// 服务端推送消息
interface WSMessage {
  type: 'system_status' | 'online_users' | 'order_queue' | 'alert';
  timestamp: string;
  data: SystemStatus | OnlineUsers | OrderQueue | Alert;
}

// 系统状态
interface SystemStatus {
  cpuUsage: number;
  memoryUsage: number;
  goroutines: number;
  dbConnections: { active: number; idle: number; max: number };
  uptime: number;
}

// 在线用户
interface OnlineUsers {
  total: number;
  peak: number;
  byRole: { user: number; player: number; admin: number };
}

// 告警
interface Alert {
  id: string;
  level: 'high' | 'medium' | 'low';
  type: 'system' | 'business' | 'security';
  message: string;
  createdAt: string;
  isRead: boolean;
}
```

### 3.4 数据模型设计

#### 新增数据库表
```sql
-- 告警记录表
CREATE TABLE alerts (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    level VARCHAR(10) NOT NULL,           -- high, medium, low
    type VARCHAR(20) NOT NULL,            -- system, business, security
    title VARCHAR(200) NOT NULL,
    message TEXT,
    source VARCHAR(100),                  -- 告警来源
    is_read BOOLEAN DEFAULT FALSE,
    read_by BIGINT,                       -- 标记已读的用户ID
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_level_created (level, created_at),
    INDEX idx_is_read (is_read)
);

-- KPI目标配置表
CREATE TABLE kpi_targets (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    period_type VARCHAR(10) NOT NULL,     -- daily, weekly, monthly
    metric_name VARCHAR(50) NOT NULL,     -- gmv, orders, users, etc.
    target_value DECIMAL(15,2) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_period_metric (period_type, metric_name, start_date)
);

-- 用户活跃统计表（每日聚合）
CREATE TABLE user_activity_daily (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    stat_date DATE NOT NULL,
    dau INT DEFAULT 0,                    -- 日活跃用户
    new_users INT DEFAULT 0,              -- 新增用户
    returning_users INT DEFAULT 0,        -- 回访用户
    paying_users INT DEFAULT 0,           -- 付费用户
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_stat_date (stat_date)
);
```

---

## 四、权限设计

### 4.1 新增权限码
```go
// 监控模块权限
MONITOR_REALTIME_VIEW    // 查看实时监控
MONITOR_ANALYTICS_VIEW   // 查看运营分析
MONITOR_KPI_VIEW         // 查看KPI仪表板
MONITOR_KPI_EDIT         // 编辑KPI目标
MONITOR_ALERT_MANAGE     // 管理告警
```

### 4.2 角色权限分配
| 角色 | 实时监控 | 运营分析 | KPI查看 | KPI编辑 | 告警管理 |
|------|---------|---------|--------|--------|---------|
| 超级管理员 | ✅ | ✅ | ✅ | ✅ | ✅ |
| 运营 | ❌ | ✅ | ✅ | ❌ | ❌ |
| 客服 | ❌ | ❌ | ❌ | ❌ | ❌ |
| 财务 | ❌ | ✅ | ✅ | ❌ | ❌ |

---

## 五、实施步骤

### Phase 1: 基础设施搭建（2个任务）
- [ ] 1.1 创建 WebSocket 基础架构（Hub、Client、Handler）
- [ ] 1.2 创建前端 WebSocket Hook 和类型定义

### Phase 2: 实时监控页面（4个任务）
- [ ] 2.1 后端：实现系统状态采集服务
- [ ] 2.2 后端：实现告警服务和存储
- [ ] 2.3 后端：实现 WebSocket 消息推送
- [ ] 2.4 前端：实现实时监控页面组件

### Phase 3: 运营分析页面（3个任务）
- [ ] 3.1 后端：实现活跃用户、留存分析 API
- [ ] 3.2 后端：实现付费分析、转化漏斗 API
- [ ] 3.3 前端：实现运营分析页面组件

### Phase 4: KPI仪表板页面（3个任务）
- [ ] 4.1 后端：实现 KPI 计算和目标管理 API
- [ ] 4.2 前端：实现 KPI 仪表板页面组件
- [ ] 4.3 前端：实现 KPI 目标配置功能

### Phase 5: 收尾工作（2个任务）
- [ ] 5.1 配置路由和权限
- [ ] 5.2 测试和文档

---

## 六、风险与注意事项

1. **WebSocket 连接管理**：需要处理断线重连、心跳检测
2. **系统指标采集**：需要考虑采集频率对性能的影响
3. **数据一致性**：统计数据的聚合计算需要考虑时区问题
4. **权限验证**：WebSocket 连接也需要进行 JWT 验证

---

## 七、预期成果

完成后，管理后台将新增：
- `/admin/monitor/realtime` - 实时监控页面
- `/admin/monitor/analytics` - 运营分析页面
- `/admin/monitor/kpi` - KPI仪表板页面

提供完整的系统监控、运营数据分析和业务指标追踪能力。
