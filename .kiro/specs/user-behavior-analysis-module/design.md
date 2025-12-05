# 用户行为分析模块设计文档

## 概述

用户行为分析模块是 GameLink 陪玩管理平台的核心数据分析系统,负责收集、处理、分析和可视化用户在平台上的各类行为数据。该模块采用前后端分离架构,前端使用 React + TypeScript + Ant Design 构建交互界面和数据可视化,后端使用 Go 语言提供高性能的数据处理和分析服务。

模块的核心目标是帮助管理员深入了解用户行为模式、消费习惯和活跃度,从而支持数据驱动的运营决策。系统支持实时数据监控、多维度数据分析、自定义事件追踪和智能报告生成等功能。

## 架构设计

### 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                        前端展示层                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │活跃度统计│  │消费分析  │  │留存分析  │  │路径分析  │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │用户画像  │  │转化漏斗  │  │游戏偏好  │  │报告生成  │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
└─────────────────────────────────────────────────────────────┘
                            ↕ HTTP/REST API
┌─────────────────────────────────────────────────────────────┐
│                        后端服务层                             │
│  ┌──────────────────┐  ┌──────────────────┐                │
│  │  行为分析服务     │  │  统计计算服务     │                │
│  │  - 事件收集       │  │  - DAU/MAU       │                │
│  │  - 数据聚合       │  │  - 留存率        │                │
│  │  - 实时监控       │  │  - 转化率        │                │
│  └──────────────────┘  └──────────────────┘                │
│  ┌──────────────────┐  ┌──────────────────┐                │
│  │  报告生成服务     │  │  数据导出服务     │                │
│  │  - 模板渲染       │  │  - Excel 导出    │                │
│  │  - PDF 生成       │  │  - CSV 导出      │                │
│  └──────────────────┘  └──────────────────┘                │
└─────────────────────────────────────────────────────────────┘
                            ↕ 数据访问
┌─────────────────────────────────────────────────────────────┐
│                        数据存储层                             │
│  ┌──────────────────┐  ┌──────────────────┐                │
│  │  PostgreSQL      │  │  Redis 缓存       │                │
│  │  - 用户数据       │  │  - 实时统计       │                │
│  │  - 订单数据       │  │  - 会话数据       │                │
│  │  - 行为事件       │  │  - 在线用户       │                │
│  └──────────────────┘  └──────────────────┘                │
└─────────────────────────────────────────────────────────────┘
```

### 技术栈

**前端:**
- React 19.2.0 + TypeScript
- Ant Design 6.0.0 (UI 组件库)
- Recharts 3.5.0 (图表库)
- Axios 1.13.2 (HTTP 客户端)
- fast-check (属性测试库)

**后端:**
- Go 1.21+
- Gin (Web 框架)
- GORM (ORM)
- PostgreSQL (主数据库)
- Redis (缓存和实时数据)
- gopter (属性测试库)

## 组件和接口

### 前端组件

#### 1. BehaviorAnalysisDashboard (行为分析仪表盘)
主容器组件,整合所有分析模块。

**Props:**
```typescript
interface BehaviorAnalysisDashboardProps {
  dateRange?: [string, string]; // 时间范围
  refreshInterval?: number; // 自动刷新间隔(毫秒)
}
```

#### 2. ActiveUsersChart (活跃用户图表)
展示 DAU/MAU 趋势和实时在线用户。

**Props:**
```typescript
interface ActiveUsersChartProps {
  dateRange: [string, string];
  userType?: 'all' | 'user' | 'player'; // 用户类型筛选
  onDataLoad?: (data: ActiveUsersData) => void;
}
```

#### 3. UsageDurationAnalysis (使用时长分析)
分析用户会话时长和使用频率。

**Props:**
```typescript
interface UsageDurationAnalysisProps {
  filters: UserFilters; // 用户筛选条件
  chartType: 'bar' | 'line' | 'pie'; // 图表类型
}
```

#### 4. ConsumptionAnalysis (消费行为分析)
展示消费趋势、客单价和复购率。

**Props:**
```typescript
interface ConsumptionAnalysisProps {
  dateRange: [string, string];
  showTrend?: boolean; // 是否显示趋势线
  showDistribution?: boolean; // 是否显示分布图
}
```

#### 5. UserProfileMap (用户画像地图)
地域分布和年龄结构可视化。

**Props:**
```typescript
interface UserProfileMapProps {
  dimension: 'region' | 'age' | 'gender'; // 分析维度
  mapType: 'china' | 'world'; // 地图类型
}
```

#### 6. RetentionHeatmap (留存热力图)
同期群留存率分析。

**Props:**
```typescript
interface RetentionHeatmapProps {
  cohortType: 'daily' | 'weekly' | 'monthly'; // 同期群类型
  retentionDays: number[]; // 留存天数 [1, 7, 30]
}
```

#### 7. ConversionFunnel (转化漏斗)
自定义转化漏斗分析。

**Props:**
```typescript
interface ConversionFunnelProps {
  funnelConfig: FunnelStep[]; // 漏斗步骤配置
  compareMode?: boolean; // 是否对比模式
  compareDateRange?: [string, string]; // 对比时间范围
}
```

#### 8. UserPathFlow (用户路径流向)
桑基图展示用户行为路径。

**Props:**
```typescript
interface UserPathFlowProps {
  startPage?: string; // 起始页面筛选
  endPage?: string; // 目标页面筛选
  maxDepth: number; // 最大路径深度
}
```

#### 9. GamePreferenceRanking (游戏偏好排行)
游戏热度和收入排行。

**Props:**
```typescript
interface GamePreferenceRankingProps {
  sortBy: 'orders' | 'revenue' | 'users'; // 排序依据
  topN: number; // 显示前 N 名
  showTrend?: boolean; // 是否显示趋势
}
```

#### 10. EventTrackingManager (事件追踪管理)
自定义事件追踪配置和统计。

**Props:**
```typescript
interface EventTrackingManagerProps {
  onEventCreate: (event: EventConfig) => Promise<void>;
  onEventDelete: (eventId: number) => Promise<void>;
}
```

#### 11. ReportGenerator (报告生成器)
生成和导出分析报告。

**Props:**
```typescript
interface ReportGeneratorProps {
  reportType: 'daily' | 'weekly' | 'monthly' | 'custom';
  exportFormat: 'pdf' | 'excel';
  autoSchedule?: ScheduleConfig; // 自动生成配置
}
```

### 后端 API 接口

#### 活跃度统计 API

```go
// GET /api/v1/admin/users/behavior/stats
type ActiveStatsResponse struct {
    OnlineUsers    int     `json:"onlineUsers"`    // 当前在线用户数
    DAU            int     `json:"dau"`            // 日活跃用户
    MAU            int     `json:"mau"`            // 月活跃用户
    UserDAU        int     `json:"userDau"`        // 普通用户 DAU
    PlayerDAU      int     `json:"playerDau"`      // 陪玩师 DAU
    LastUpdateTime string  `json:"lastUpdateTime"` // 最后更新时间
}

// GET /api/v1/admin/users/behavior/trend
type ActiveTrendResponse struct {
    Dates      []string `json:"dates"`      // 日期列表
    DAUValues  []int    `json:"dauValues"`  // DAU 数值
    MAUValues  []int    `json:"mauValues"`  // MAU 数值
}
```

#### 使用时长分析 API

```go
// GET /api/v1/admin/users/behavior/duration
type UsageDurationResponse struct {
    AvgSessionDuration float64            `json:"avgSessionDuration"` // 平均会话时长(分钟)
    TotalUsageDuration float64            `json:"totalUsageDuration"` // 总使用时长(小时)
    AvgDailyLogins     float64            `json:"avgDailyLogins"`     // 日均登录次数
    AvgWeeklyDays      float64            `json:"avgWeeklyDays"`      // 周均登录天数
    DurationDistribution []DurationBucket `json:"durationDistribution"` // 时长分布
}

type DurationBucket struct {
    Range string `json:"range"` // 时长区间 "0-30min"
    Count int    `json:"count"` // 用户数量
}
```

#### 消费行为分析 API

```go
// GET /api/v1/admin/users/behavior/consumption
type ConsumptionResponse struct {
    TotalAmount      float64   `json:"totalAmount"`      // 总消费金额
    AvgOrderValue    float64   `json:"avgOrderValue"`    // 平均客单价
    PayingUsers      int       `json:"payingUsers"`      // 付费用户数
    RepurchaseRate   float64   `json:"repurchaseRate"`   // 复购率
    AvgRepurchaseDays float64  `json:"avgRepurchaseDays"` // 平均复购周期(天)
    TrendData        []TrendPoint `json:"trendData"`     // 趋势数据
}

type TrendPoint struct {
    Date   string  `json:"date"`   // 日期
    Amount float64 `json:"amount"` // 金额
}
```

#### 用户画像分析 API

```go
// GET /api/v1/admin/users/behavior/distribution
type UserDistributionResponse struct {
    RegionDistribution []RegionData `json:"regionDistribution"` // 地域分布
    AgeDistribution    []AgeData    `json:"ageDistribution"`    // 年龄分布
    GenderDistribution []GenderData `json:"genderDistribution"` // 性别分布
}

type RegionData struct {
    Province string  `json:"province"` // 省份
    City     string  `json:"city"`     // 城市
    Count    int     `json:"count"`    // 用户数
    Percent  float64 `json:"percent"`  // 占比
}

type AgeData struct {
    AgeRange string  `json:"ageRange"` // 年龄段 "18-25"
    Count    int     `json:"count"`    // 用户数
    Percent  float64 `json:"percent"`  // 占比
}
```

#### 留存分析 API

```go
// GET /api/v1/admin/users/behavior/retention
type RetentionResponse struct {
    Day1Retention  float64        `json:"day1Retention"`  // 次日留存率
    Day7Retention  float64        `json:"day7Retention"`  // 7日留存率
    Day30Retention float64        `json:"day30Retention"` // 30日留存率
    CohortData     []CohortRetention `json:"cohortData"`  // 同期群数据
    ChurnedUsers   int            `json:"churnedUsers"`   // 流失用户数
}

type CohortRetention struct {
    CohortDate string    `json:"cohortDate"` // 同期群日期
    UserCount  int       `json:"userCount"`  // 用户数
    Retention  []float64 `json:"retention"`  // 各天留存率
}
```

#### 转化漏斗 API

```go
// POST /api/v1/admin/users/behavior/funnel
type FunnelRequest struct {
    Steps     []FunnelStep `json:"steps"`     // 漏斗步骤
    DateRange [2]string    `json:"dateRange"` // 时间范围
}

type FunnelStep struct {
    Name      string `json:"name"`      // 步骤名称
    EventType string `json:"eventType"` // 事件类型
}

type FunnelResponse struct {
    Steps []FunnelStepResult `json:"steps"` // 步骤结果
}

type FunnelStepResult struct {
    Name           string  `json:"name"`           // 步骤名称
    UserCount      int     `json:"userCount"`      // 用户数
    ConversionRate float64 `json:"conversionRate"` // 转化率
}
```

#### 行为路径 API

```go
// GET /api/v1/admin/users/behavior/path
type UserPathResponse struct {
    Paths []PathSequence `json:"paths"` // 路径序列
}

type PathSequence struct {
    Sequence  []string `json:"sequence"`  // 页面序列
    UserCount int      `json:"userCount"` // 用户数
    Percent   float64  `json:"percent"`   // 占比
}
```

#### 游戏偏好 API

```go
// GET /api/v1/admin/users/behavior/game-preference
type GamePreferenceResponse struct {
    Rankings []GameRanking `json:"rankings"` // 游戏排行
}

type GameRanking struct {
    GameID      int     `json:"gameId"`      // 游戏 ID
    GameName    string  `json:"gameName"`    // 游戏名称
    OrderCount  int     `json:"orderCount"`  // 订单数
    Revenue     float64 `json:"revenue"`     // 收入
    UserCount   int     `json:"userCount"`   // 用户数
    TrendData   []int   `json:"trendData"`   // 趋势数据
}
```

#### 事件追踪 API

```go
// POST /api/v1/admin/users/behavior/events
type CreateEventRequest struct {
    Name       string            `json:"name"`       // 事件名称
    Type       string            `json:"type"`       // 事件类型
    Parameters map[string]string `json:"parameters"` // 事件参数
}

// POST /api/v1/admin/users/behavior/events/track
type TrackEventRequest struct {
    EventName  string                 `json:"eventName"`  // 事件名称
    UserID     int                    `json:"userId"`     // 用户 ID
    Timestamp  string                 `json:"timestamp"`  // 时间戳
    Parameters map[string]interface{} `json:"parameters"` // 事件参数
}

// GET /api/v1/admin/users/behavior/events/stats
type EventStatsResponse struct {
    Events []EventStat `json:"events"` // 事件统计
}

type EventStat struct {
    EventName   string `json:"eventName"`   // 事件名称
    TriggerCount int   `json:"triggerCount"` // 触发次数
    UniqueUsers int    `json:"uniqueUsers"`  // 独立用户数
}
```

#### 报告生成 API

```go
// POST /api/v1/admin/users/behavior/reports
type GenerateReportRequest struct {
    ReportType string    `json:"reportType"` // 报告类型
    DateRange  [2]string `json:"dateRange"`  // 时间范围
    Format     string    `json:"format"`     // 导出格式
}

type GenerateReportResponse struct {
    ReportID   int    `json:"reportId"`   // 报告 ID
    DownloadURL string `json:"downloadUrl"` // 下载链接
}

// POST /api/v1/admin/users/behavior/reports/schedule
type ScheduleReportRequest struct {
    ReportType string `json:"reportType"` // 报告类型
    Frequency  string `json:"frequency"`  // 频率 daily/weekly/monthly
    Recipients []string `json:"recipients"` // 接收人邮箱
}
```

## 数据模型

### 用户行为事件表 (user_behavior_events)

```sql
CREATE TABLE user_behavior_events (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    event_name VARCHAR(100) NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    event_data JSONB,
    page_url VARCHAR(500),
    session_id VARCHAR(100),
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_event_name (event_name),
    INDEX idx_created_at (created_at),
    INDEX idx_session_id (session_id)
);
```

### 用户会话表 (user_sessions)

```sql
CREATE TABLE user_sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    session_id VARCHAR(100) UNIQUE NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    duration_seconds INTEGER,
    page_views INTEGER DEFAULT 0,
    ip_address VARCHAR(45),
    device_type VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_session_id (session_id),
    INDEX idx_start_time (start_time)
);
```

### 用户活跃度统计表 (user_activity_stats)

```sql
CREATE TABLE user_activity_stats (
    id BIGSERIAL PRIMARY KEY,
    stat_date DATE NOT NULL UNIQUE,
    dau INTEGER NOT NULL DEFAULT 0,
    mau INTEGER NOT NULL DEFAULT 0,
    user_dau INTEGER NOT NULL DEFAULT 0,
    player_dau INTEGER NOT NULL DEFAULT 0,
    new_users INTEGER NOT NULL DEFAULT 0,
    active_sessions INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_stat_date (stat_date)
);
```

### 用户留存数据表 (user_retention_data)

```sql
CREATE TABLE user_retention_data (
    id BIGSERIAL PRIMARY KEY,
    cohort_date DATE NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    day_1_active BOOLEAN DEFAULT FALSE,
    day_7_active BOOLEAN DEFAULT FALSE,
    day_30_active BOOLEAN DEFAULT FALSE,
    last_active_date DATE,
    is_churned BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_cohort_date (cohort_date),
    INDEX idx_user_id (user_id),
    UNIQUE (cohort_date, user_id)
);
```

### 转化漏斗配置表 (conversion_funnels)

```sql
CREATE TABLE conversion_funnels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    steps JSONB NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### 事件追踪配置表 (event_tracking_configs)

```sql
CREATE TABLE event_tracking_configs (
    id SERIAL PRIMARY KEY,
    event_name VARCHAR(100) UNIQUE NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    parameters JSONB,
    is_active BOOLEAN DEFAULT TRUE,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_event_name (event_name)
);
```

### TypeScript 类型定义

```typescript
// 活跃用户数据
interface ActiveUsersData {
  onlineUsers: number;
  dau: number;
  mau: number;
  userDau: number;
  playerDau: number;
  lastUpdateTime: string;
}

// 使用时长数据
interface UsageDurationData {
  avgSessionDuration: number;
  totalUsageDuration: number;
  avgDailyLogins: number;
  avgWeeklyDays: number;
  durationDistribution: DurationBucket[];
}

interface DurationBucket {
  range: string;
  count: number;
}

// 消费行为数据
interface ConsumptionData {
  totalAmount: number;
  avgOrderValue: number;
  payingUsers: number;
  repurchaseRate: number;
  avgRepurchaseDays: number;
  trendData: TrendPoint[];
}

interface TrendPoint {
  date: string;
  amount: number;
}

// 用户分布数据
interface UserDistribution {
  regionDistribution: RegionData[];
  ageDistribution: AgeData[];
  genderDistribution: GenderData[];
}

interface RegionData {
  province: string;
  city: string;
  count: number;
  percent: number;
}

interface AgeData {
  ageRange: string;
  count: number;
  percent: number;
}

// 留存数据
interface RetentionData {
  day1Retention: number;
  day7Retention: number;
  day30Retention: number;
  cohortData: CohortRetention[];
  churnedUsers: number;
}

interface CohortRetention {
  cohortDate: string;
  userCount: number;
  retention: number[];
}

// 转化漏斗
interface FunnelConfig {
  steps: FunnelStep[];
  dateRange: [string, string];
}

interface FunnelStep {
  name: string;
  eventType: string;
}

interface FunnelResult {
  steps: FunnelStepResult[];
}

interface FunnelStepResult {
  name: string;
  userCount: number;
  conversionRate: number;
}

// 用户路径
interface UserPath {
  paths: PathSequence[];
}

interface PathSequence {
  sequence: string[];
  userCount: number;
  percent: number;
}

// 游戏偏好
interface GamePreference {
  rankings: GameRanking[];
}

interface GameRanking {
  gameId: number;
  gameName: string;
  orderCount: number;
  revenue: number;
  userCount: number;
  trendData: number[];
}

// 事件追踪
interface EventConfig {
  name: string;
  type: string;
  parameters: Record<string, string>;
}

interface EventStat {
  eventName: string;
  triggerCount: number;
  uniqueUsers: number;
}

// 报告配置
interface ReportConfig {
  reportType: 'daily' | 'weekly' | 'monthly' | 'custom';
  dateRange: [string, string];
  format: 'pdf' | 'excel';
}

interface ScheduleConfig {
  frequency: 'daily' | 'weekly' | 'monthly';
  recipients: string[];
}

// 用户筛选条件
interface UserFilters {
  role?: 'user' | 'player' | 'admin';
  registrationDateRange?: [string, string];
  region?: string;
  ageRange?: string;
}
```

## 正确性属性

*属性是指在系统所有有效执行中都应该成立的特征或行为——本质上是关于系统应该做什么的形式化陈述。属性是人类可读规范和机器可验证正确性保证之间的桥梁。*


### 属性反思

在编写正确性属性之前,我需要识别逻辑上冗余的属性并进行合并:

**冗余分析:**
1. 属性 2.2(使用频率计算)和属性 2.5(会话记录)可以合并,因为会话记录是计算使用频率的基础
2. 属性 3.4(支付记录)和属性 3.1(消费统计显示)、3.3(消费分布)、3.5(复购计算)存在依赖关系,但它们测试不同层面,保留
3. 属性 4.5(用户信息记录)是属性 4.2(地域统计)和属性 4.4(画像分析)的数据来源,但测试不同功能,保留
4. 属性 7.4(页面访问记录)是属性 7.1(路径统计)和属性 7.2(路径分析)的基础,但测试不同层面,保留
5. 属性 8.4(订单记录)是属性 8.1(游戏排行)、8.2(游戏统计)、8.3(游戏趋势)的数据来源,但测试不同功能,保留

**合并决策:**
- 将属性 2.2 和 2.5 合并为一个综合属性:会话记录完整性和频率计算正确性
- 其他属性保持独立,因为它们测试不同的功能层面

### 正确性属性列表

**属性 1: 时间范围数据完整性**
*对于任意*时间范围查询,返回的 DAU 和 MAU 数据应该包含完整的日期序列,且数值非负
**验证需求: 1.2**

**属性 2: 用户类型活跃度分离**
*对于任意*活跃用户数据集,按用户类型(普通用户/陪玩师)分类统计的活跃度之和应该等于总活跃度
**验证需求: 1.3**

**属性 3: 自动刷新时间间隔**
*对于任意*启用自动刷新的页面,数据刷新函数应该在每 5 分钟(300000 毫秒)被调用一次
**验证需求: 1.5**

**属性 4: 会话记录完整性和使用频率计算**
*对于任意*用户的会话记录集合,每条记录都应该包含开始时间、结束时间和持续时长,且基于这些记录计算的日均登录次数和周均登录天数应该准确反映实际登录模式
**验证需求: 2.2, 2.5**

**属性 5: 用户筛选结果正确性**
*对于任意*用户数据集和筛选条件(角色/注册时间/地域),筛选结果中的所有用户都应该满足所有指定的筛选条件
**验证需求: 2.3**

**属性 6: 消费分布统计准确性**
*对于任意*消费记录集合,按金额区间统计的用户数量之和应该等于总付费用户数,且每个用户只出现在一个区间中
**验证需求: 3.3**

**属性 7: 支付事件记录完整性**
*对于任意*支付完成事件,系统记录应该包含订单金额、支付方式和支付时间三个必要字段,且金额应该为正数
**验证需求: 3.4**

**属性 8: 复购率计算正确性**
*对于任意*用户订单历史,如果用户有多次购买记录,则应该被计入复购用户;复购率应该等于复购用户数除以总付费用户数
**验证需求: 3.5**

**属性 9: 地域分布统计一致性**
*对于任意*用户地域数据集,所有地区的用户数量之和应该等于总用户数,且所有地区的占比之和应该等于 100%
**验证需求: 4.2**

**属性 10: 用户画像多维度完整性**
*对于任意*用户画像分析结果,应该同时包含性别、年龄和地域三个维度的数据,且每个维度的数据都应该非空
**验证需求: 4.4**

**属性 11: 用户信息记录完整性**
*对于任意*用户注册或资料更新事件,系统应该记录地理位置和年龄信息,且这些信息应该可以被后续的画像分析查询到
**验证需求: 4.5**

**属性 12: 同期群留存率计算正确性**
*对于任意*用户注册同期群,留存率应该按注册日期分组计算,且后续天数的留存率应该小于或等于前一天的留存率(单调递减)
**验证需求: 5.2**

**属性 13: 流失用户识别准确性**
*对于任意*用户登录记录,如果用户最后登录时间距今超过 30 天,则应该被标记为流失用户;反之则不应该被标记
**验证需求: 5.4**

**属性 14: 转化漏斗计算正确性**
*对于任意*漏斗配置和用户行为数据,每个步骤的用户数应该小于或等于前一步骤的用户数,且转化率应该等于当前步骤用户数除以前一步骤用户数
**验证需求: 6.2**

**属性 15: 漏斗步骤记录完整性**
*对于任意*用户完成的漏斗步骤,系统应该记录步骤完成时间和用户标识,且这些记录应该能被漏斗分析查询到
**验证需求: 6.4**

**属性 16: 漏斗对比时间段独立性**
*对于任意*两个不同的时间段,漏斗对比分析应该分别计算每个时间段的转化数据,且两个时间段的数据应该相互独立
**验证需求: 6.5**

**属性 17: 行为路径统计准确性**
*对于任意*用户访问记录集合,统计出的最常见路径序列应该按用户数量降序排列,且所有路径的用户数之和应该小于或等于总用户数
**验证需求: 7.1**

**属性 18: 路径占比计算正确性**
*对于任意*路径数据集,每条路径的占比应该等于该路径的用户数除以总用户数,且所有路径的占比之和应该小于或等于 100%
**验证需求: 7.2**

**属性 19: 页面访问记录完整性**
*对于任意*页面访问事件,系统应该记录页面 URL、访问时间和停留时长,且这些记录应该能被路径分析查询到
**验证需求: 7.4**

**属性 20: 路径筛选结果正确性**
*对于任意*路径数据集和筛选条件(起始页面/目标页面),筛选结果中的所有路径都应该满足指定的起始或目标页面条件
**验证需求: 7.5**

**属性 21: 游戏排行榜排序正确性**
*对于任意*游戏订单数据集,排行榜应该按照指定的排序依据(订单数/收入/用户数)降序排列,且排名靠前的游戏的对应指标值应该大于或等于排名靠后的游戏
**验证需求: 8.1**

**属性 22: 游戏统计数据准确性**
*对于任意*游戏的订单记录集合,统计的订单数量应该等于该游戏的订单记录数,收入金额应该等于所有订单金额之和
**验证需求: 8.2**

**属性 23: 游戏热度趋势单调性**
*对于任意*游戏的时间序列订单数据,热度趋势曲线应该准确反映每个时间点的订单数量变化,且数据点的时间顺序应该正确
**验证需求: 8.3**

**属性 24: 订单游戏关联记录完整性**
*对于任意*订单创建事件,系统应该记录订单关联的游戏和服务项目,且这些关联信息应该能被游戏偏好分析查询到
**验证需求: 8.4**

**属性 25: 游戏偏好分组正确性**
*对于任意*用户游戏偏好数据,按游戏偏好分组的用户应该根据其订单中最常玩的游戏进行分类,且每个用户只应该出现在一个分组中
**验证需求: 8.5**

**属性 26: 事件追踪记录完整性**
*对于任意*用户触发的追踪事件,系统应该记录事件名称、时间戳和相关参数,且记录的时间戳应该在事件触发时间的合理范围内
**验证需求: 9.2**

**属性 27: 事件统计数据准确性**
*对于任意*事件记录集合,统计的触发次数应该等于该事件的记录总数,独立用户数应该等于触发该事件的不重复用户数
**验证需求: 9.3**

**属性 28: 事件追踪删除后行为**
*对于任意*被删除的事件追踪配置,系统应该停止收集新的事件数据,但历史记录应该保持可查询状态
**验证需求: 9.5**

**属性 29: 报告内容完整性**
*对于任意*报告生成请求,生成的报告内容应该包含关键指标、图表和分析结论三个必要部分,且内容应该与选择的报告类型和时间范围相匹配
**验证需求: 10.2**

## 错误处理

### 前端错误处理

1. **API 请求失败**
   - 显示友好的错误提示信息
   - 提供重试机制
   - 记录错误日志到监控系统

2. **数据加载超时**
   - 设置合理的超时时间(30 秒)
   - 显示加载超时提示
   - 允许用户手动刷新

3. **数据格式错误**
   - 验证 API 返回数据的格式
   - 对异常数据进行容错处理
   - 显示数据异常提示

4. **图表渲染失败**
   - 捕获图表库的渲染错误
   - 显示图表加载失败提示
   - 提供降级展示方案(表格形式)

5. **权限不足**
   - 检查用户权限
   - 显示权限不足提示
   - 引导用户联系管理员

### 后端错误处理

1. **数据库查询失败**
   ```go
   if err := db.Query(&result); err != nil {
       log.Error("Database query failed", "error", err)
       return nil, errors.New("数据查询失败,请稍后重试")
   }
   ```

2. **数据计算异常**
   ```go
   if totalUsers == 0 {
       return 0, errors.New("无有效用户数据")
   }
   retentionRate := float64(retainedUsers) / float64(totalUsers)
   ```

3. **缓存失效**
   ```go
   data, err := cache.Get(key)
   if err != nil {
       // 缓存失效,从数据库加载
       data, err = db.Query()
       if err == nil {
           cache.Set(key, data, expiration)
       }
   }
   ```

4. **并发访问控制**
   ```go
   mutex.Lock()
   defer mutex.Unlock()
   // 执行需要同步的操作
   ```

5. **参数验证**
   ```go
   if dateRange[0] > dateRange[1] {
       return nil, errors.New("开始日期不能晚于结束日期")
   }
   ```

## 测试策略

### 单元测试

**前端单元测试 (Vitest + Testing Library):**

1. **组件渲染测试**
   - 测试组件在不同 props 下的渲染结果
   - 验证关键 UI 元素的存在性
   - 测试条件渲染逻辑

2. **数据处理函数测试**
   - 测试数据转换函数的正确性
   - 测试边界条件和异常输入
   - 验证计算逻辑的准确性

3. **API 客户端测试**
   - 使用 Mock 测试 API 调用
   - 验证请求参数的正确性
   - 测试错误处理逻辑

**后端单元测试 (Go testing + testify):**

1. **业务逻辑测试**
   - 测试统计计算函数
   - 测试数据聚合逻辑
   - 验证边界条件处理

2. **数据访问层测试**
   - 使用测试数据库
   - 测试 CRUD 操作
   - 验证查询结果的正确性

3. **API 处理器测试**
   - 测试请求参数验证
   - 测试响应数据格式
   - 验证错误处理逻辑

### 属性测试

**前端属性测试 (fast-check):**

使用 fast-check 库进行属性测试,每个测试至少运行 100 次迭代。

1. **数据转换属性测试**
   ```typescript
   // 测试时间范围数据完整性
   fc.assert(
     fc.property(
       fc.date(), fc.date(),
       (startDate, endDate) => {
         const range = [startDate, endDate].sort();
         const result = fetchActiveUsers(range);
         return result.dates.length > 0 && 
                result.dauValues.every(v => v >= 0);
       }
     ),
     { numRuns: 100 }
   );
   ```

2. **筛选逻辑属性测试**
   ```typescript
   // 测试用户筛选结果正确性
   fc.assert(
     fc.property(
       fc.array(fc.record({
         role: fc.constantFrom('user', 'player', 'admin'),
         region: fc.string(),
         age: fc.integer(18, 80)
       })),
       fc.record({
         role: fc.option(fc.constantFrom('user', 'player', 'admin')),
         region: fc.option(fc.string())
       }),
       (users, filters) => {
         const filtered = filterUsers(users, filters);
         return filtered.every(user => 
           (!filters.role || user.role === filters.role) &&
           (!filters.region || user.region === filters.region)
         );
       }
     ),
     { numRuns: 100 }
   );
   ```

**后端属性测试 (gopter):**

使用 gopter 库进行属性测试,每个测试至少运行 100 次迭代。

1. **统计计算属性测试**
   ```go
   // 测试地域分布统计一致性
   properties := gopter.NewProperties(nil)
   properties.Property("region distribution sum equals total", 
     prop.ForAll(
       func(users []User) bool {
         distribution := CalculateRegionDistribution(users)
         totalCount := 0
         totalPercent := 0.0
         for _, region := range distribution {
           totalCount += region.Count
           totalPercent += region.Percent
         }
         return totalCount == len(users) && 
                math.Abs(totalPercent - 100.0) < 0.01
       },
       gen.SliceOf(genUser()),
     ),
   )
   properties.TestingRun(t, gopter.ConsoleReporter(false))
   ```

2. **留存率计算属性测试**
   ```go
   // 测试同期群留存率单调递减
   properties.Property("retention rate monotonic decrease",
     prop.ForAll(
       func(cohort CohortData) bool {
         retention := CalculateRetention(cohort)
         for i := 1; i < len(retention); i++ {
           if retention[i] > retention[i-1] {
             return false
           }
         }
         return true
       },
       genCohortData(),
     ),
   )
   ```

### 集成测试

1. **前后端集成测试**
   - 测试完整的数据流程
   - 验证 API 集成的正确性
   - 测试实时数据更新机制

2. **数据库集成测试**
   - 测试复杂查询的性能
   - 验证事务处理的正确性
   - 测试数据一致性

3. **缓存集成测试**
   - 测试缓存命中和失效
   - 验证缓存数据的一致性
   - 测试缓存更新机制

### 端到端测试

使用 Playwright 进行关键用户流程的端到端测试:

1. **查看活跃度统计流程**
   - 访问行为分析页面
   - 选择时间范围
   - 验证数据展示

2. **生成分析报告流程**
   - 配置报告参数
   - 生成报告
   - 下载报告文件

3. **创建转化漏斗流程**
   - 配置漏斗步骤
   - 查看漏斗分析
   - 对比不同时期

### 测试覆盖率目标

- 单元测试覆盖率: > 80%
- 属性测试: 覆盖所有 29 个正确性属性
- 集成测试: 覆盖所有 API 端点
- 端到端测试: 覆盖 3 个核心用户流程

### 测试数据管理

1. **测试数据生成**
   - 使用工厂模式生成测试数据
   - 提供多种数据场景的生成器
   - 支持随机数据和固定数据

2. **测试数据隔离**
   - 每个测试使用独立的数据集
   - 测试后自动清理数据
   - 避免测试间的相互影响

3. **测试数据版本控制**
   - 维护测试数据的版本
   - 支持数据回滚和恢复
   - 记录数据变更历史
