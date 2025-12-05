# 设计文档 - 仪表盘模块

## 概述

仪表盘模块是 GameLink 平台管理后台的核心入口，为管理员提供平台运营数据的全景视图。该模块采用响应式设计，支持实时数据更新，通过数据可视化技术将复杂的业务数据转化为直观易懂的图表和指标。模块设计遵循组件化原则，确保代码的可复用性和可维护性。

## 架构

### 系统架构图

```
┌─────────────────────────────────────────────────────────────┐
│                     前端展示层 (React)                       │
├─────────────────────────────────────────────────────────────┤
│  主仪表盘  │  实时监控  │  运营分析  │  KPI仪表板            │
└─────────────────────────────────────────────────────────────┘
                              │
                              ↓ HTTP/WebSocket
┌─────────────────────────────────────────────────────────────┐
│                      API 层 (Go Backend)                     │
├─────────────────────────────────────────────────────────────┤
│  统计API  │  实时数据API  │  KPI计算API  │  报表API         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                      数据聚合层                              │
├─────────────────────────────────────────────────────────────┤
│  订单聚合  │  用户聚合  │  收入聚合  │  实时数据聚合         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                      数据存储层                              │
├─────────────────────────────────────────────────────────────┤
│  PostgreSQL/SQLite  │  Redis (缓存/实时)  │  InfluxDB (时序)│
│  (业务数据)         │                     │                 │
└─────────────────────────────────────────────────────────────┘
```

### 技术栈

**前端**：
- React 19.2.0
- TypeScript
- Ant Design 6.0.0 (UI组件库)
- Recharts 3.5.0 (图表库)
- Socket.IO Client (WebSocket)

**后端**：
- Go 1.21+
- Gin Web Framework
- WebSocket
- Redis (缓存和实时数据)
- PostgreSQL (生产环境业务数据)
- SQLite (测试环境业务数据)

## 组件和接口

### 前端组件

#### 1. Dashboard 主组件
**职责**：仪表盘主页面容器，协调各子组件

**Props**：
```typescript
interface DashboardProps {
  timeRange: TimeRange;
  onTimeRangeChange: (range: TimeRange) => void;
}
```

**状态**：
```typescript
interface DashboardState {
  loading: boolean;
  stats: DashboardStats;
  error: Error | null;
}
```

#### 2. StatCard 组件
**职责**：展示单个统计指标

**Props**：
```typescript
interface StatCardProps {
  title: string;
  value: number | string;
  trend?: {
    value: number;
    isUp: boolean;
  };
  icon: ReactNode;
  loading?: boolean;
  onClick?: () => void;
}
```

#### 3. RevenueChart 组件
**职责**：收入趋势折线图

**Props**：
```typescript
interface RevenueChartProps {
  data: RevenueData[];
  timeRange: TimeRange;
  loading?: boolean;
}

interface RevenueData {
  date: string;
  revenue: number;
  netRevenue?: number;
}
```

#### 4. OrderStatusPie 组件
**职责**：订单状态分布饼图

**Props**：
```typescript
interface OrderStatusPieProps {
  data: OrderStatusData[];
  loading?: boolean;
}

interface OrderStatusData {
  status: OrderStatus;
  count: number;
  percentage: number;
}
```

#### 5. RealtimeMonitor 组件
**职责**：实时监控面板

**Props**：
```typescript
interface RealtimeMonitorProps {
  wsUrl: string;
  onAlert: (alert: Alert) => void;
}
```

**状态**：
```typescript
interface RealtimeMonitorState {
  connected: boolean;
  metrics: RealtimeMetrics;
}

interface RealtimeMetrics {
  onlineUsers: number;
  queueLength: number;
  cpuUsage: number;
  memoryUsage: number;
}
```

### 后端接口

#### 1. DashboardService 接口
```go
type DashboardService interface {
    // 获取仪表盘统计数据
    GetDashboardStats(ctx context.Context, req *DashboardStatsRequest) (*DashboardStats, error)
    
    // 获取收入趋势
    GetRevenueTrend(ctx context.Context, req *TrendRequest) (*RevenueTrend, error)
    
    // 获取用户增长趋势
    GetUserGrowth(ctx context.Context, req *TrendRequest) (*UserGrowth, error)
    
    // 获取最新订单
    GetRecentOrders(ctx context.Context, limit int) ([]*Order, error)
    
    // 获取热门陪玩师
    GetTopPlayers(ctx context.Context, limit int) ([]*PlayerRanking, error)
}
```

#### 2. RealtimeService 接口
```go
type RealtimeService interface {
    // 获取实时监控数据
    GetRealtimeMetrics(ctx context.Context) (*RealtimeMetrics, error)
    
    // 订阅实时数据更新
    SubscribeMetrics(ctx context.Context) (<-chan *RealtimeMetrics, error)
    
    // 获取告警列表
    GetAlerts(ctx context.Context, req *AlertRequest) ([]*Alert, error)
}
```

#### 3. KPIService 接口
```go
type KPIService interface {
    // 计算KPI指标
    CalculateKPI(ctx context.Context, req *KPIRequest) (*KPIMetrics, error)
    
    // 获取运营数据
    GetOperationalData(ctx context.Context, req *OperationalRequest) (*OperationalData, error)
}
```

## 数据模型

### DashboardStats (仪表盘统计)
```typescript
interface DashboardStats {
  totalUsers: number;              // 总用户数
  totalOrders: number;             // 总订单数
  totalRevenue: number;            // 总收入
  activePlayers: number;           // 活跃陪玩师数
  userGrowth: TrendData;           // 用户增长趋势
  revenueTrend: TrendData;         // 收入趋势
  orderStatusDistribution: OrderStatusData[];  // 订单状态分布
  recentOrders: Order[];           // 最新订单
  topPlayers: PlayerRanking[];     // 热门陪玩师
}

interface TrendData {
  current: number;
  previous: number;
  changePercentage: number;
  isUp: boolean;
}
```

### RealtimeMetrics (实时监控指标)
```typescript
interface RealtimeMetrics {
  timestamp: Date;                 // 时间戳
  onlineUsers: number;             // 在线用户数
  queueLength: number;             // 订单处理队列长度
  cpuUsage: number;                // CPU使用率 (0-100)
  memoryUsage: number;             // 内存使用率 (0-100)
  diskUsage: number;               // 磁盘使用率 (0-100)
  networkIn: number;               // 网络入流量 (MB/s)
  networkOut: number;              // 网络出流量 (MB/s)
}
```

### KPIMetrics (KPI指标)
```typescript
interface KPIMetrics {
  conversionRate: number;          // 转化率 (%)
  repeatPurchaseRate: number;      // 复购率 (%)
  userRetentionRate: number;       // 用户留存率 (%)
  averageOrderValue: number;       // 平均客单价
  dau: number;                     // 日活用户数
  mau: number;                     // 月活用户数
  paymentRate: number;             // 付费率 (%)
  arpu: number;                    // 平均每用户收入
}
```

### Alert (告警)
```typescript
interface Alert {
  id: number;                      // 告警ID
  type: AlertType;                 // 告警类型
  level: AlertLevel;               // 告警级别
  title: string;                   // 告警标题
  message: string;                 // 告警内容
  timestamp: Date;                 // 发生时间
  isRead: boolean;                 // 是否已读
  isResolved: boolean;             // 是否已解决
}

enum AlertType {
  SYSTEM = 'system',               // 系统异常
  BUSINESS = 'business',           // 业务异常
  SECURITY = 'security'            // 安全告警
}

enum AlertLevel {
  INFO = 'info',
  WARNING = 'warning',
  ERROR = 'error',
  CRITICAL = 'critical'
}
```

## 正确性属性

*属性是一个特征或行为，应该在系统的所有有效执行中保持为真——本质上是关于系统应该做什么的正式陈述。属性作为人类可读规范和机器可验证正确性保证之间的桥梁。*

### 属性 1：统计数据完整性
*对于任何*仪表盘统计数据请求，返回的数据必须包含所有必填字段（总用户数、总订单数、总收入、活跃陪玩师数），且数值必须为非负数
**验证：需求 1.1**

### 属性 2：趋势数据一致性
*对于任何*趋势数据，当前值与前一周期值的差值计算出的变化百分比必须与返回的changePercentage字段一致
**验证：需求 1.3**

### 属性 3：饼图数据总和
*对于任何*订单状态分布饼图数据，所有状态的百分比之和必须等于100%（误差范围±0.1%）
**验证：需求 2.1**

### 属性 4：时间范围数据点数量
*对于任何*折线图数据，当选择7天时间范围时数据点数量必须为7，选择30天时必须为30，选择90天时必须为90
**验证：需求 3.2**

### 属性 5：最新订单排序
*对于任何*最新订单列表，订单必须按创建时间降序排列，第一条记录的创建时间必须大于或等于最后一条记录的创建时间
**验证：需求 5.1**

### 属性 6：热门陪玩师排序
*对于任何*热门陪玩师列表，陪玩师必须按订单数降序排列，排名序号必须从1开始连续递增
**验证：需求 6.2**

### 属性 7：实时数据更新频率
*对于任何*实时监控数据，两次连续更新之间的时间间隔必须在0.9秒到1.1秒之间
**验证：需求 7.5**

### 属性 8：监控指标范围
*对于任何*监控指标（CPU、内存、磁盘使用率），数值必须在0到100之间
**验证：需求 7.2**

### 属性 9：KPI指标有效性
*对于任何*KPI指标（转化率、复购率、留存率、付费率），数值必须在0到100之间
**验证：需求 9.2**

### 属性 10：WebSocket重连机制
*对于任何*WebSocket连接断开事件，系统必须在5秒内尝试重新连接，最多重试3次
**验证：需求 7.4**

## 错误处理

### 错误分类

#### 1. 数据加载错误
- **网络错误**: API请求失败
  - 处理：显示错误提示，提供重试按钮
  - 示例：网络连接超时、服务器无响应
  
- **数据格式错误**: 返回数据格式不符合预期
  - 处理：记录错误日志，显示默认值或空状态
  - 示例：字段缺失、类型不匹配

#### 2. WebSocket连接错误
- **连接失败**: 无法建立WebSocket连接
  - 处理：显示连接失败提示，自动重试
  - 示例：服务器拒绝连接、网络不可达
  
- **连接断开**: WebSocket连接意外断开
  - 处理：显示断开提示，自动重连
  - 示例：网络波动、服务器重启

#### 3. 图表渲染错误
- **数据为空**: 图表数据为空数组
  - 处理：显示"暂无数据"占位符
  
- **数据异常**: 数据包含异常值
  - 处理：过滤异常值或标记异常点

### 错误处理策略

#### 前端错误处理
```typescript
// 统一错误处理Hook
const useDashboardData = (timeRange: TimeRange) => {
  const [data, setData] = useState<DashboardStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);
  
  const loadData = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await api.getDashboardStats(timeRange);
      setData(response.data);
    } catch (err) {
      setError(err as Error);
      message.error('加载数据失败，请稍后重试');
    } finally {
      setLoading(false);
    }
  };
  
  useEffect(() => {
    loadData();
  }, [timeRange]);
  
  return { data, loading, error, retry: loadData };
};
```

#### WebSocket错误处理
```typescript
const useRealtimeData = (wsUrl: string) => {
  const [connected, setConnected] = useState(false);
  const [metrics, setMetrics] = useState<RealtimeMetrics | null>(null);
  const socketRef = useRef<Socket | null>(null);
  const retryCountRef = useRef(0);
  const maxRetries = 3;
  
  const connect = () => {
    socketRef.current = io(wsUrl, {
      reconnection: true,
      reconnectionDelay: 5000,
      reconnectionAttempts: maxRetries
    });
    
    socketRef.current.on('connect', () => {
      setConnected(true);
      retryCountRef.current = 0;
      message.success('实时监控已连接');
    });
    
    socketRef.current.on('disconnect', () => {
      setConnected(false);
      message.warning('实时监控连接断开');
    });
    
    socketRef.current.on('metrics', (data: RealtimeMetrics) => {
      setMetrics(data);
    });
    
    socketRef.current.on('connect_error', () => {
      retryCountRef.current++;
      if (retryCountRef.current >= maxRetries) {
        message.error('实时监控连接失败，请刷新页面重试');
      }
    });
  };
  
  useEffect(() => {
    connect();
    return () => {
      socketRef.current?.disconnect();
    };
  }, [wsUrl]);
  
  return { connected, metrics };
};
```

## 测试策略

### 单元测试

#### 前端组件测试
```typescript
describe('StatCard Component', () => {
  it('should render stat value correctly', () => {
    render(<StatCard title="总用户数" value={1000} />);
    expect(screen.getByText('1000')).toBeInTheDocument();
  });
  
  it('should display trend indicator when trend prop is provided', () => {
    render(
      <StatCard 
        title="总用户数" 
        value={1000} 
        trend={{ value: 10, isUp: true }}
      />
    );
    expect(screen.getByText('↑ 10%')).toBeInTheDocument();
  });
  
  it('should call onClick when card is clicked', () => {
    const onClick = vi.fn();
    render(<StatCard title="总用户数" value={1000} onClick={onClick} />);
    fireEvent.click(screen.getByRole('button'));
    expect(onClick).toHaveBeenCalled();
  });
});
```

### 属性测试

```typescript
// 属性 3：饼图数据总和
describe('Property: Pie Chart Data Sum', () => {
  it('sum of all percentages should equal 100%', () => {
    fc.assert(
      fc.property(
        fc.array(
          fc.record({
            status: fc.constantFrom('pending', 'confirmed', 'in_progress', 'completed', 'cancelled', 'refunded'),
            count: fc.integer({ min: 0, max: 1000 }),
            percentage: fc.float({ min: 0, max: 100 })
          }),
          { minLength: 1, maxLength: 6 }
        ),
        (data) => {
          const totalPercentage = data.reduce((sum, item) => sum + item.percentage, 0);
          // 允许0.1%的误差
          expect(Math.abs(totalPercentage - 100)).toBeLessThanOrEqual(0.1);
        }
      ),
      { numRuns: 100 }
    );
  });
});
```

### 集成测试

```typescript
describe('Dashboard Integration', () => {
  it('should load and display all dashboard data', async () => {
    // Mock API responses
    mockAPI.getDashboardStats.mockResolvedValue({
      data: {
        totalUsers: 1000,
        totalOrders: 500,
        totalRevenue: 50000,
        activePlayers: 100
      }
    });
    
    render(<Dashboard timeRange="7d" />);
    
    // 等待数据加载
    await waitFor(() => {
      expect(screen.getByText('1000')).toBeInTheDocument();
      expect(screen.getByText('500')).toBeInTheDocument();
      expect(screen.getByText('¥50,000')).toBeInTheDocument();
      expect(screen.getByText('100')).toBeInTheDocument();
    });
  });
});
```

### E2E测试

```typescript
test('dashboard complete workflow', async ({ page }) => {
  // 登录
  await page.goto('/login');
  await page.fill('[name="username"]', 'admin');
  await page.fill('[name="password"]', 'password');
  await page.click('button[type="submit"]');
  
  // 等待跳转到仪表盘
  await page.waitForURL('/admin/dashboard');
  
  // 验证统计卡片显示
  await expect(page.locator('.stat-card')).toHaveCount(4);
  
  // 切换时间范围
  await page.click('[data-testid="time-range-selector"]');
  await page.click('[data-value="30d"]');
  
  // 验证图表更新
  await expect(page.locator('.revenue-chart')).toBeVisible();
  
  // 点击统计卡片跳转
  await page.click('.stat-card:first-child');
  await page.waitForURL('/admin/users');
});
```
