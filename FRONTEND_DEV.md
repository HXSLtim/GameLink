# GameLink 陪玩管理平台前端规划

## 📋 项目概述

为 GameLink 陪玩管理平台规划完整的后台管理系统前端功能，基于 **React + TypeScript + Ant Design** 架构，构建功能完备、用户体验优秀的管理端界面。

---

## 🛠️ 技术栈

| 类别        | 技术方案           | 版本   |
| :---------- | :----------------- | :----- |
| 框架        | React + TypeScript | 19.2.0 |
| UI 库       | Ant Design         | 6.0.0  |
| 路由        | React Router       | v7.9.6 |
| HTTP 客户端 | Axios              | 1.13.2 |
| 图表库      | Recharts           | 3.5.0  |
| 构建工具    | Vite               | 7.2.4  |

---

## 📊 核心功能模块规划

### 1. 仪表盘模块 📈

#### 现有功能

- 主仪表盘：核心数据统计、订单状态分布、支付状态分布、收入趋势、用户增长、最新订单、热门陪玩师
- 数据可视化：饼图、折线图、统计卡片
- 时间筛选：近 7 天/30 天/90 天数据切换

#### 需要补充的功能

- 实时监控面板：系统运行状态、在线用户、订单处理队列
- 异常告警提示：系统异常、业务异常、安全告警
- 业务指标 KPI：转化率、复购率、用户留存率
- 运营数据概览：日活、月活、付费率、客单价

#### 页面列表

| 页面       | 状态      | 功能描述             | 优先级 |
| :--------- | :-------- | :------------------- | :----- |
| 主仪表盘   | ✅ 已存在 | 核心数据统计和可视化 | 高     |
| 实时监控   | ✅ 已存在   | 系统实时状态监控     | 中     |
| 运营分析   | ✅ 已存在   | 运营数据分析报告     | 中     |
| KPI 仪表板 | ✅ 已存在   | 关键业务指标展示     | 中     |

#### 关键 API 端点

```http
GET  /api/v1/admin/stats/dashboard        # 仪表板数据
GET  /api/v1/admin/stats/revenue-trend    # 收入趋势
GET  /api/v1/admin/stats/user-growth      # 用户增长
GET  /api/v1/admin/stats/realtime         # 实时监控数据
GET  /api/v1/admin/stats/kpi              # KPI 指标数据
```

#### 组件设计

- `StatCard`: 统计卡片，展示核心指标
- `RevenueChart`: 收入趋势图表
- `UserGrowthChart`: 用户增长图表
- `OrderStatusPie`: 订单状态饼图
- `RecentOrders`: 最新订单列表
- `TopPlayers`: 热门陪玩师排名

---

### 2. 用户管理模块 👥

#### 现有功能

- 用户 CRUD：创建、读取、更新、删除用户
- 角色管理：用户/陪玩师/管理员三种角色
- 状态管理：正常/封禁/停用状态切换
- 用户详情：基本信息、注册时间、最后登录
- 批量操作：批量删除用户
- 搜索过滤：关键词、角色、状态、时间范围

#### 需要补充的功能

- 用户行为分析：登录频率、使用习惯、消费行为
- 用户标签管理：给用户打标签，便于分类管理
- 用户等级体系：VIP 等级、积分系统
- 用户画像分析：年龄分布、地域分布、偏好分析
- 批量操作增强：批量修改角色、批量发送通知
- 用户登录历史：登录 IP、设备、地理位置
- 用户操作日志：关键操作记录

#### 页面列表

| 页面         | 状态      | 功能描述         | 优先级 | 后端API | 前端集成 | 数据来源 |
| :----------- | :-------- | :--------------- | :----- | :------ | :------- | :------- |
| 用户列表     | ✅ 已存在 | 用户管理和搜索   | 高     | ✅      | ✅       | ✅ 真实  |
| 用户详情     | ✅ 已存在 | 用户详细信息     | 高     | ✅      | ✅       | ✅ 真实  |
| 用户行为分析 | ✅ 已完成 | 用户行为数据统计 | 中     | ✅      | ✅       | ✅ 真实  |
| 用户标签管理 | ✅ 已完成 | 标签创建和管理   | 中     | ✅      | ✅       | ✅ 真实  |
| 用户等级管理 | ✅ 已完成 | 等级体系设置     | 低     | ✅      | ✅       | ✅ 真实  |
| 用户画像分析 | ✅ 已完成 | 用户群体分析     | 低     | ✅      | ✅       | ✅ 真实  |
| 用户登录历史 | ✅ 已完成 | 登录记录查询     | 中     | ✅      | ✅       | ⚠️ 待数据|

#### 关键表单字段

```typescript
interface UserForm {
  name: string; // 用户名
  email: string; // 邮箱
  phone: string; // 手机号
  password: string; // 密码
  avatarUrl?: string; // 头像 URL
  role: "user" | "player" | "admin"; // 角色
  status: "active" | "banned" | "suspended"; // 状态
  tags?: string[]; // 用户标签
  level?: number; // 用户等级
  vipExpiry?: string; // VIP 到期时间
}
```

#### 关键 API 端点

```http
# 用户基础管理
GET    /api/v1/admin/users                    # 用户列表
POST   /api/v1/admin/users                    # 创建用户
GET    /api/v1/admin/users/:id                # 用户详情
PUT    /api/v1/admin/users/:id                # 更新用户
DELETE /api/v1/admin/users/:id                # 删除用户
PUT    /api/v1/admin/users/:id/status         # 更新状态
PUT    /api/v1/admin/users/:id/role           # 更新角色
GET    /api/v1/admin/users/:id/orders         # 用户订单
GET    /api/v1/admin/users/:id/logs           # 用户操作日志
GET    /api/v1/admin/users/:id/login-history  # 用户登录历史 ✅

# 用户行为分析 ✅
GET    /api/v1/admin/users/behavior/stats         # 行为统计(DAU/时长/消费)
GET    /api/v1/admin/users/behavior/trend         # 活动趋势(7天)
GET    /api/v1/admin/users/behavior/distribution  # 用户分布(地域/年龄)

# 用户标签管理
GET    /api/v1/admin/tags                     # 标签列表
POST   /api/v1/admin/tags                     # 创建标签
PUT    /api/v1/admin/tags/:id                 # 更新标签
DELETE /api/v1/admin/tags/:id                 # 删除标签
POST   /api/v1/admin/users/:userId/tags/:tagId    # 分配标签
DELETE /api/v1/admin/users/:userId/tags/:tagId    # 移除标签
PUT    /api/v1/admin/users/:userId/tags       # 批量设置标签

# 批量操作
POST   /api/v1/admin/users/batch/role         # 批量更新角色
POST   /api/v1/admin/users/batch/status       # 批量更新状态
POST   /api/v1/admin/users/batch/delete       # 批量删除用户
POST   /api/v1/admin/users/batch/points       # 批量添加积分
POST   /api/v1/admin/users/batch/notification # 批量发送通知
```

---

### 3. 陪玩师管理模块 🎮

#### 现有功能

- 陪玩师列表：展示所有陪玩师信息
- 申请审核：待审核陪玩师的审核流程
- 状态管理：待审核/已通过/已拒绝/已封禁
- 详细信息：身份证照片、擅长游戏、个人介绍
- 在线状态：在线/离线/忙碌
- 封禁/解封：状态切换操作

#### 需要补充的功能

- 技能标签管理：游戏技能、服务技能标签
- 等级评定系统：根据订单量、评价等自动评级
- 绩效考核：订单完成率、用户满意度等指标
- 收益统计详细分析：收入趋势、游戏偏好收益
- 服务质量监控：投诉率、差评率跟踪
- 陪玩师培训管理：培训内容、考核记录
- 认证流程优化：多级认证、快速认证通道

#### 页面列表

| 页面         | 状态      | 功能描述           | 优先级 |
| :----------- | :-------- | :----------------- | :----- |
| 陪玩师列表   | ✅ 已存在 | 陪玩师管理和搜索   | 高     |
| 陪玩师详情   | ✅ 已存在 | 详细信息和认证资料 | 高     |
| 陪玩师审核   | ✅ 已存在 | 申请审核流程       | 高     |
| 技能标签管理 | 🆕 新增   | 技能标签创建和管理 | 中     |
| 绩效考核     | 🆕 新增   | 绩效指标和评价     | 中     |
| 培训管理     | 🆕 新增   | 培训内容和记录     | 低     |
| 认证管理     | 🆕 新增   | 认证流程和设置     | 低     |

#### 关键表单字段

```typescript
interface PlayerForm {
  name: string; // 陪玩师昵称
  gender: "male" | "female"; // 性别
  age: number; // 年龄
  games: string[]; // 擅长游戏
  introduction: string; // 个人介绍
  voiceSample?: string; // 语音样本
  idCardFront: string; // 身份证正面
  idCardBack: string; // 身份证反面
  skillTags: string[]; // 技能标签
  level: number; // 等级
  commissionRate: number; // 佣金比例
}
```

#### 关键 API 端点

```http
GET    /api/v1/admin/players                    # 陪玩师列表
POST   /api/v1/admin/players                    # 创建陪玩师
GET    /api/v1/admin/players/:id                # 陪玩师详情
PUT    /api/v1/admin/players/:id                # 更新陪玩师
PUT    /api/v1/admin/players/:id/verification   # 更新认证状态
PUT    /api/v1/admin/players/:id/games          # 更新擅长游戏
PUT    /api/v1/admin/players/:id/skill-tags     # 更新技能标签
```

---

### 4. 订单管理模块 📋

#### 现有功能

- 订单列表：完整订单信息展示
- 状态管理：待确认/已确认/进行中/已完成/已取消/已退款
- 支付跟踪：支付状态监控
- 订单详情：用户信息、陪玩师信息、订单进度
- 取消订单：订单取消功能
- 退款处理：退款申请和审核
- 订单时间线：订单状态变更历史

#### 需要补充的功能

- 订单分配算法：智能匹配陪玩师
- 争议处理流程：订单纠纷处理机制
- 订单评价管理：评价审核和回复
- 订单统计分析：转化率、客单价等分析
- 批量订单操作：批量确认、批量取消等
- 订单模板管理：常用订单模板
- 自动接单规则：陪玩师自动接单设置

#### 页面列表

| 页面         | 状态      | 功能描述           | 优先级 |
| :----------- | :-------- | :----------------- | :----- |
| 订单列表     | ✅ 已存在 | 订单管理和搜索     | 高     |
| 订单详情     | ✅ 已存在 | 订单详细信息和进度 | 高     |
| 争议处理     | 🆕 新增   | 订单纠纷处理       | 高     |
| 订单统计     | 🆕 新增   | 订单数据分析       | 中     |
| 自动接单设置 | 🆕 新增   | 自动接单规则配置   | 中     |
| 订单模板     | 🆕 新增   | 订单模板管理       | 低     |

#### 关键 API 端点

```http
GET    /api/v1/admin/orders              # 订单列表
POST   /api/v1/admin/orders              # 创建订单
GET    /api/v1/admin/orders/:id          # 订单详情
PUT    /api/v1/admin/orders/:id          # 更新订单
POST   /api/v1/admin/orders/:id/cancel   # 取消订单
POST   /api/v1/admin/orders/:id/refund   # 订单退款
POST   /api/v1/admin/orders/:id/assign   # 指派陪玩师
POST   /api/v1/admin/orders/:id/confirm  # 确认订单
GET    /api/v1/admin/orders/:id/timeline # 订单时间线
```

---

### 5. 游戏管理模块 🎯

#### 现有功能

- 游戏 CRUD：创建、读取、更新、删除游戏
- 分类管理：MOBA/射击/RPG/卡牌/休闲/其他
- 平台支持：PC/iOS/Android/主机
- 状态管理：上架/下架控制
- 排序管理：游戏显示顺序
- 图标上传：游戏图标管理

#### 需要补充的功能

- 游戏标签系统：游戏特色标签
- 游戏热度排行：根据订单量统计热度
- 游戏数据分析：用户偏好、收益分析
- 游戏推荐算法：个性化游戏推荐
- 游戏版本管理：游戏版本更新记录
- 游戏活动管理：游戏相关活动
- 游戏社区管理：游戏讨论区管理

#### 页面列表

| 页面         | 状态      | 功能描述       | 优先级 |
| :----------- | :-------- | :------------- | :----- |
| 游戏列表     | ✅ 已存在 | 游戏管理和搜索 | 高     |
| 游戏编辑     | ✅ 已存在 | 游戏信息编辑   | 高     |
| 游戏标签管理 | 🆕 新增   | 游戏标签系统   | 中     |
| 游戏数据分析 | 🆕 新增   | 游戏数据统计   | 中     |
| 游戏活动管理 | 🆕 新增   | 游戏活动配置   | 低     |
| 游戏推荐设置 | 🆕 新增   | 推荐算法配置   | 低     |

#### 关键 API 端点

```http
GET    /api/v1/admin/games           # 游戏列表
POST   /api/v1/admin/games           # 创建游戏
GET    /api/v1/admin/games/:id       # 游戏详情
PUT    /api/v1/admin/games/:id       # 更新游戏
DELETE /api/v1/admin/games/:id       # 删除游戏
GET    /api/v1/admin/games/:id/logs  # 游戏操作日志
```

---

### 6. 服务项目管理模块 🎁

#### 现有功能

- 服务项目 CRUD：创建、读取、更新、删除服务
- 分类管理：上分/陪玩/教学/娱乐
- 价格管理：服务价格设置
- 状态管理：启用/禁用控制
- 批量操作：批量状态更新
- 搜索过滤：游戏、分类、状态筛选

#### 需要补充的功能

- 服务包管理：多个服务打包销售
- 价格策略设置：时段价格、会员价格等
- 服务评价模板：标准化评价模板
- 服务标签系统：服务特色标签
- 服务推荐算法：个性化服务推荐
- 服务使用统计：服务热度分析
- 服务质量监控：服务满意度跟踪

#### 页面列表

| 页面          | 状态      | 功能描述     | 优先级 |
| :------------ | :-------- | :----------- | :----- |
| 服务列表      | ✅ 已存在 | 服务项目管理 | 高     |
| 创建/编辑服务 | ✅ 已存在 | 服务表单编辑 | 高     |
| 服务详情      | ✅ 已存在 | 服务详情展示 | 高     |
| 服务包管理    | 🆕 新增   | 服务打包销售 | 中     |
| 价格策略      | 🆕 新增   | 动态价格配置 | 中     |
| 服务统计      | 🆕 新增   | 服务热度分析 | 中     |

#### 关键表单字段

```typescript
interface ServiceItemForm {
  name: string; // 服务名称
  gameId: number; // 所属游戏
  description: string; // 服务描述
  price: number; // 价格
  duration: number; // 时长（小时）
  category: string; // 分类
  tags: string[]; // 标签
  icon?: string; // 图标
  status: "active" | "inactive"; // 状态
  sortOrder: number; // 排序
  packageId?: number; // 服务包 ID
  discountRules?: any[]; // 折扣规则
}
```

#### 关键 API 端点

```http
GET    /api/v1/admin/service-items                    # 服务列表
POST   /api/v1/admin/service-items                    # 创建服务
GET    /api/v1/admin/service-items/:id                # 服务详情
PUT    /api/v1/admin/service-items/:id                # 更新服务
DELETE /api/v1/admin/service-items/:id                # 删除服务
POST   /api/v1/admin/service-items/batch-update-status # 批量更新状态
POST   /api/v1/admin/service-items/batch-update-price  # 批量更新价格
```

---

### 7. 支付与财务模块 💰

#### 需要实现的功能

- 支付记录管理：所有支付流水记录
- 退款处理：退款申请审核和处理
- 财务对账：资金流水对账
- 资金流水：详细的资金变动记录
- 财务报表：收支统计、利润分析
- 税务管理：税务相关配置和报表
- 第三方支付配置：支付宝、微信等配置

#### 页面列表

| 页面         | 状态    | 功能描述       | 优先级 |
| :----------- | :------ | :------------- | :----- |
| 支付记录列表 | 🆕 新增 | 支付流水管理   | 高     |
| 退款管理     | 🆕 新增 | 退款申请和处理 | 高     |
| 财务对账     | 🆕 新增 | 资金对账功能   | 中     |
| 资金流水     | 🆕 新增 | 详细流水记录   | 中     |
| 财务报表     | 🆕 新增 | 财务数据统计   | 中     |
| 支付配置     | 🆕 新增 | 支付方式配置   | 低     |

#### 关键表单字段

```typescript
interface PaymentForm {
  orderId: number; // 订单 ID
  userId: number; // 用户 ID
  method: string; // 支付方式
  amount: number; // 金额
  currency: string; // 货币
  status: "pending" | "paid" | "refunded" | "failed"; // 状态
  providerTradeNo: string; // 第三方交易号
  refundedAmount?: number; // 已退款金额
}
```

#### 关键 API 端点

```http
GET    /api/v1/admin/payments           # 支付记录列表
POST   /api/v1/admin/payments           # 创建支付记录
GET    /api/v1/admin/payments/:id       # 支付详情
PUT    /api/v1/admin/payments/:id       # 更新支付
POST   /api/v1/admin/payments/:id/refund   # 退款处理
POST   /api/v1/admin/payments/:id/capture  # 确认入账
GET    /api/v1/admin/payments/:id/logs     # 支付操作日志
```

---

### 8. 佣金与结算模块 💵

#### 需要实现的功能

- 抽成规则设置：平台抽成比例配置
- 结算周期管理：日结/周结/月结设置
- 自动结算配置：自动结算触发条件
- 结算记录查询：历史结算记录
- 佣金统计报表：佣金收入分析
- 提现审核流程：陪玩师提现审核
- 打款管理：批量打款操作

#### 页面列表

| 页面         | 状态    | 功能描述       | 优先级 |
| :----------- | :------ | :------------- | :----- |
| 佣金规则设置 | 🆕 新增 | 抽成比例配置   | 高     |
| 结算管理     | 🆕 新增 | 结算周期和记录 | 高     |
| 提现审核     | 🆕 新增 | 提现申请处理   | 高     |
| 佣金统计     | 🆕 新增 | 佣金收入分析   | 中     |
| 结算记录     | 🆕 新增 | 历史结算查询   | 中     |
| 打款管理     | 🆕 新增 | 批量打款操作   | 低     |

#### 关键表单字段

```typescript
interface CommissionRuleForm {
  name: string; // 规则名称
  gameId?: number; // 适用游戏
  minAmount: number; // 最低金额
  maxAmount?: number; // 最高金额
  rate: number; // 抽成比例
  description?: string; // 规则描述
  isActive: boolean; // 是否启用
}

interface SettlementForm {
  playerId: number; // 陪玩师 ID
  startDate: string; // 结算开始日期
  endDate: string; // 结算结束日期
  totalAmount: number; // 总结算金额
  commission: number; // 平台佣金
  netAmount: number; // 实际结算金额
  status: "pending" | "processed" | "paid"; // 状态
}
```

#### 关键 API 端点

```http
POST   /api/v1/admin/commission/rules                  # 创建抽成规则
PUT    /api/v1/admin/commission/rules/:id              # 更新规则
GET    /api/v1/admin/commission/rules                  # 规则列表
GET    /api/v1/admin/commission/stats                  # 佣金统计
POST   /api/v1/admin/commission/settlements/trigger    # 触发结算
GET    /api/v1/admin/withdraws                         # 提现列表
POST   /api/v1/admin/withdraws/:id/approve             # 审核通过
POST   /api/v1/admin/withdraws/:id/reject              # 拒绝提现
POST   /api/v1/admin/withdraws/:id/complete            # 完成打款
```

---

### 9. 评价管理模块 ⭐

#### 需要实现的功能

- 评价审核：用户评价内容审核
- 举报处理：评价举报处理流程
- 评价回复：管理员回复用户评价
- 评价统计分析：评价数据统计
- 敏感词过滤：评价内容敏感词检测
- 评价展示设置：评价显示规则配置

#### 页面列表

| 页面       | 状态    | 功能描述     | 优先级 |
| :--------- | :------ | :----------- | :----- |
| 评价列表   | 🆕 新增 | 评价内容管理 | 中     |
| 评价审核   | 🆕 新增 | 评价内容审核 | 中     |
| 举报管理   | 🆕 新增 | 评价举报处理 | 中     |
| 评价统计   | 🆕 新增 | 评价数据分析 | 低     |
| 敏感词管理 | 🆕 新增 | 敏感词库管理 | 中     |
| 评价设置   | 🆕 新增 | 评价规则配置 | 低     |

#### 关键表单字段

```typescript
interface ReviewForm {
  orderId: number; // 订单 ID
  reviewerId: number; // 评价者 ID
  playerId: number; // 陪玩师 ID
  score: number; // 评分（1-5）
  content: string; // 评价内容
  images?: string[]; // 评价图片
  status: "pending" | "approved" | "rejected"; // 审核状态
  isReported: boolean; // 是否被举报
}
```

#### 关键 API 端点

```http
GET    /api/v1/admin/reviews           # 评价列表
POST   /api/v1/admin/reviews           # 创建评价
GET    /api/v1/admin/reviews/:id       # 评价详情
PUT    /api/v1/admin/reviews/:id       # 更新评价
DELETE /api/v1/admin/reviews/:id       # 删除评价
GET    /api/v1/admin/reviews/:id/logs  # 评价操作日志
POST   /api/v1/admin/reviews/:id/reply # 回复评价
```

---

### 10. RBAC 权限系统模块 🔐

#### 现有功能（已完善）

- 权限 CRUD：创建、读取、更新、删除权限
- 角色管理：角色创建和权限分配
- 菜单管理：动态菜单配置
- 权限分组：权限按模块分组
- 细粒度控制：方法+路径级别的权限控制
- 权限同步：前后端权限同步机制

#### 需要补充的功能

- 权限模板管理：常用权限组合模板
- 权限变更审计：权限修改历史记录
- 角色继承机制：角色层级继承关系
- 权限测试工具：权限验证测试功能
- 权限性能优化：权限查询性能优化

#### 页面列表

| 页面     | 状态      | 功能描述       | 优先级 |
| :------- | :-------- | :------------- | :----- |
| 权限列表 | ✅ 已存在 | 权限管理和分组 | 高     |
| 角色管理 | ✅ 已存在 | 角色和权限分配 | 高     |
| 菜单管理 | ✅ 已存在 | 动态菜单配置   | 高     |
| 权限模板 | 🆕 新增   | 权限模板管理   | 中     |
| 操作审计 | 🆕 新增   | 权限变更审计   | 中     |
| 权限测试 | 🆕 新增   | 权限测试工具   | 低     |

#### 关键表单字段

```typescript
interface PermissionForm {
  name: string; // 权限名称
  method: string; // HTTP 方法
  path: string; // API 路径
  group: string; // 所属分组
  description?: string; // 描述
}

interface RoleForm {
  name: string; // 角色名称
  slug: string; // 角色标识
  description?: string; // 描述
  permissions: number[]; // 权限 ID 列表
  isSystem?: boolean; // 是否系统角色
}

interface MenuForm {
  name: string; // 菜单名称
  path?: string; // 路由路径
  icon?: string; // 图标
  parentId?: number; // 父菜单 ID
  sortOrder: number; // 排序
  permission?: string; // 所需权限
}
```

#### 关键 API 端点

```http
GET    /api/v1/admin/permissions           # 权限列表
POST   /api/v1/admin/permissions           # 创建权限
GET    /api/v1/admin/permissions/groups    # 权限分组
GET    /api/v1/admin/permissions/:id       # 权限详情
PUT    /api/v1/admin/permissions/:id       # 更新权限
GET    /api/v1/admin/roles                 # 角色列表
POST   /api/v1/admin/roles                 # 创建角色
PUT    /api/v1/admin/roles/:id             # 更新角色
PUT    /api/v1/admin/roles/:id/permissions # 分配权限
POST   /api/v1/admin/roles/assign-user     # 分配角色给用户
GET    /api/v1/admin/menus                 # 菜单列表
POST   /api/v1/admin/menus                 # 创建菜单
PUT    /api/v1/admin/menus/:id             # 更新菜单
DELETE /api/v1/admin/menus/:id             # 删除菜单
```

---

### 11. 内容管理模块 📝

#### 需要实现的功能

- 动态内容审核：用户动态内容审核
- 聊天内容监控：实时聊天内容监控
- 敏感词管理：敏感词库维护
- 内容举报处理：用户举报内容处理
- 内容分类管理：内容类型分类
- 内容发布审核：发布前内容审核
- 用户生成内容管理：UGC 内容统一管理

#### 页面列表

| 页面       | 状态    | 功能描述         | 优先级 |
| :--------- | :------ | :--------------- | :----- |
| 动态审核   | 🆕 新增 | 用户动态内容审核 | 中     |
| 聊天监控   | 🆕 新增 | 实时聊天内容监控 | 中     |
| 敏感词管理 | 🆕 新增 | 敏感词库管理     | 中     |
| 举报处理   | 🆕 新增 | 内容举报处理     | 中     |
| 内容分类   | 🆕 新增 | 内容类型管理     | 低     |
| 发布审核   | 🆕 新增 | 内容发布审核     | 低     |

#### 关键表单字段

```typescript
interface FeedForm {
  userId: number; // 发布者 ID
  content: string; // 内容
  images?: string[]; // 图片
  visibility: "public" | "private" | "friends"; // 可见性
  status: "pending" | "approved" | "rejected"; // 审核状态
}
```

#### 关键 API 端点

```http
GET    /api/v1/user/feeds                 # 动态列表
POST   /api/v1/admin/feeds/:id/approve    # 批准动态
POST   /api/v1/admin/feeds/:id/reject     # 拒绝动态
GET    /api/v1/user/chat/groups/:id/messages  # 聊天消息
POST   /api/v1/admin/chat/messages/:id/report # 举报消息
```

---

### 12. 系统管理模块 ⚙️

#### 现有功能（部分存在）

- 系统配置：基础系统参数配置
- 审计日志：操作日志记录

#### 需要补充的功能

- 系统监控：服务器性能、数据库状态监控
- 数据备份恢复：自动备份和恢复功能
- 系统维护工具：系统维护相关工具
- 性能监控分析：系统性能分析报表
- 错误日志管理：系统错误日志查看
- 操作日志审计：详细的操作审计追踪

#### 页面列表

| 页面     | 状态      | 功能描述     | 优先级 |
| :------- | :-------- | :----------- | :----- |
| 系统配置 | ✅ 已存在 | 系统参数配置 | 高     |
| 审计日志 | ✅ 已存在 | 操作日志审计 | 高     |
| 系统监控 | 🆕 新增   | 系统状态监控 | 中     |
| 数据备份 | 🆕 新增   | 备份恢复管理 | 中     |
| 性能分析 | 🆕 新增   | 性能监控分析 | 中     |
| 错误日志 | 🆕 新增   | 错误日志管理 | 中     |
| 维护工具 | 🆕 新增   | 系统维护工具 | 低     |

#### 关键 API 端点

```http
GET    /api/v1/system/config          # 系统配置
GET    /api/v1/system/db              # 数据库状态
GET    /api/v1/system/cache           # 缓存状态
GET    /api/v1/system/resources       # 系统资源
GET    /api/v1/system/version         # 系统版本
GET    /api/v1/admin/stats/audit/overview  # 审计概览
GET    /api/v1/admin/stats/audit/trend     # 审计趋势
```

---

## 📋 实施优先级规划

### 🚀 阶段一：核心功能完善（Week 1-2）

**优先级：高**

- 增强仪表盘实时监控功能
- 完善支付与财务模块前端界面
- 实现评价管理模块
- 补充内容管理模块

### 🏗️ 阶段二：新增核心模块（Week 3-4）

**优先级：中**

- 开发佣金与结算模块
- 实现系统监控功能
- 增强用户行为分析
- 完善陪玩师绩效考核

### 🎨 阶段三：优化与增强（Week 5-6）

**优先级：低**

- 智能推荐算法
- 高级数据分析
- 自动化运营工具
- 移动端适配优化

---

## 🎯 关键组件设计模式

### 通用组件架构

```tsx
const CrudPage: React.FC = () => {
  // 状态管理
  const [data, setData] = useState<DataItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });

  // 权限检查
  const permissions = usePermissions({
    canCreate: MODULE_PERMISSIONS.CREATE,
    canEdit: MODULE_PERMISSIONS.UPDATE,
    canDelete: MODULE_PERMISSIONS.DELETE,
  });

  // 数据加载
  const loadData = useCallback(async (params?: any) => {
    setLoading(true);
    try {
      const response = await api.getList(params);
      setData(response.data.data || []);
      setPagination({
        current: response.data.pagination?.page || 1,
        pageSize: response.data.pagination?.pageSize || 10,
        total: response.data.pagination?.total || 0,
      });
    } finally {
      setLoading(false);
    }
  }, []);

  // 创建和编辑
  const [formVisible, setFormVisible] = useState(false);
  const [editingRecord, setEditingRecord] = useState<DataItem | null>(null);

  const handleCreate = () => {
    setEditingRecord(null);
    setFormVisible(true);
  };

  const handleEdit = (record: DataItem) => {
    setEditingRecord(record);
    setFormVisible(true);
  };

  const handleDelete = async (id: number) => {
    await api.delete(id);
    message.success("删除成功");
    loadData();
  };

  return (
    <PageContainer title="标题">
      <SearchTable
        columns={columns}
        dataSource={data}
        loading={loading}
        pagination={pagination}
        searchFields={searchFields}
        onSearch={loadData}
        onCreate={permissions.canCreate ? handleCreate : undefined}
        onEdit={permissions.canEdit ? handleEdit : undefined}
        onDelete={permissions.canDelete ? handleDelete : undefined}
      />
      {formVisible && (
        <CrudForm
          visible={formVisible}
          onCancel={() => setFormVisible(false)}
          onSubmit={async (values) => {
            if (editingRecord) {
              await api.update(editingRecord.id, values);
              message.success("更新成功");
            } else {
              await api.create(values);
              message.success("创建成功");
            }
            setFormVisible(false);
            loadData();
          }}
          initialValues={editingRecord || undefined}
        />
      )}
    </PageContainer>
  );
};
```

---

## 📁 关键文件路径

### 核心配置文件

```
frontend/src/constants/permissions.ts    # 权限常量定义
frontend/src/api/admin.ts                # 管理端 API 接口
frontend/src/router/routes.tsx           # 路由配置
frontend/src/config/adminRoutes.ts       # 菜单和权限配置
frontend/src/types/admin.ts              # TypeScript 类型定义
```

### 组件文件

```
frontend/src/components/SearchTable/     # 搜索表格组件
frontend/src/components/PageContainer/   # 页面容器组件
frontend/src/components/StatCard/        # 统计卡片组件
frontend/src/components/PermissionGuard/ # 权限守卫组件
```

### 页面文件

```
frontend/src/pages/admin/                # 核心管理页面
frontend/src/pages/sys/                  # 系统管理页面
frontend/src/pages/biz/                  # 业务管理页面
```

### 后端文件

```
backend/internal/handler/admin/router.go # 后端路由注册
backend/internal/router/router.go        # 总路由器配置
backend/internal/router/adminRoutes.go   # 权限分配逻辑
```

---

## 🛡️ 权限设计原则

1. **最小权限原则**：每个角色只赋予完成工作所必需的最小权限
2. **权限分级**：
   - 超级管理员：拥有所有权限
   - 普通管理员：根据职责分配权限
   - 模块管理员：只管理特定模块
3. **权限自检**：敏感操作需要二次确认
4. **操作审计**：所有重要操作记录日志
5. **动态权限**：支持运行时动态调整权限

---

## 🎨 UI/UX 设计规范

### 一致性

- 所有页面遵循统一的设计语言

### 响应式

- 支持多种屏幕尺寸

### 可访问性

- 符合 WCAG 2.1 标准

### 性能优化

- 代码分割（Code Splitting）
- 懒加载（Lazy Loading）
- 虚拟滚动（Virtual Scrolling）

### 用户体验

- 加载状态提示
- 操作反馈（成功/失败提示）
- 表单验证
- 操作确认对话框

---

## ✅ 成功标准

### 功能完整性

- 所有规划的功能模块均可正常使用

### 用户体验

- 页面加载时间 < 2 秒
- 操作响应流畅

### 代码质量

- TypeScript 严格模式
- 无 any 类型

### 测试覆盖率

- 核心功能测试覆盖率 > 80%

### 文档完整性

- 每个模块都有详细的使用文档

---

## 📅 详细实施计划

### 第一阶段：核心功能完善（Week 1-2）

#### 第 1 周：仪表盘增强和基础设施

- **Day 1-2**: 实时监控面板开发
  - 实现 WebSocket 连接获取实时数据
  - 创建系统运行状态监控组件
  - 实现在线用户实时监控
- **Day 3-4**: KPI 指标系统
  - 转化率、复购率、留存率计算逻辑
  - KPI 数据可视化组件开发
  - 时间筛选功能实现
- **Day 5**: 异常告警系统
  - 告警规则配置界面
  - 实时告警通知组件

#### 第 2 周：支付与财务模块

- **Day 1-2**: 支付记录管理
  - 支付列表页面（带搜索、筛选）
  - 支付详情页面
  - 退款处理功能
- **Day 3-4**: 财务对账功能
  - 对账报表生成
  - 资金流水查询
  - 财务报表导出（Excel/PDF）
- **Day 5**: 支付配置管理
  - 第三方支付配置界面

### 第二阶段：新增核心模块（Week 3-4）

#### 第 3 周：佣金结算和用户管理增强

- **Day 1-2**: 佣金规则设置
  - 抽成比例配置表单
  - 规则优先级管理
  - 规则测试工具
- **Day 3-4**: 提现审核流程
  - 提现申请列表
  - 审核通过/拒绝功能
  - 批量打款操作
- **Day 5**: 结算管理
  - 结算周期配置
  - 自动结算触发

#### 第 4 周：系统监控和玩家管理

- **Day 1-2**: 系统监控面板
  - 服务器性能监控
  - 数据库状态监控
  - 缓存状态监控
- **Day 3-4**: 陪玩师绩效考核
  - KPI 指标设置
  - 绩效报表生成
  - 等级自动评定
- **Day 5**: 技能标签管理
  - 标签创建和管理
  - 标签分类体系

### 第三阶段：优化与增强（Week 5-6）

#### 第 5 周：内容管理和用户体验

- **Day 1-2**: 动态内容审核
  - 动态列表和审核
  - 批量审核功能
- **Day 3-4**: 敏感词管理
  - 敏感词库维护
  - 自动检测规则
- **Day 5**: 评价管理
  - 评价审核流程
  - 评价回复功能

#### 第 6 周：性能优化和测试

- **Day 1-2**: 代码性能优化
  - 虚拟滚动实现
  - 懒加载优化
  - 缓存策略
- **Day 3-4**: 测试用例完善
  - 单元测试编写
  - 集成测试编写
  - E2E 测试编写
- **Day 5**: Bug 修复和文档
  - 修复测试中发现的问题
  - 完善开发文档

---

## 🧪 测试计划

### 测试策略

1. **单元测试**（目标覆盖率：80%+）

   - 工具：Vitest + Testing Library
   - 覆盖范围：API 客户端函数、工具函数、自定义 Hooks、纯组件

2. **集成测试**

   - 覆盖范围：组件与 API 的交互、权限控制逻辑、表单验证流程、数据流管理

3. **E2E 测试**
   - 工具：Playwright
   - 覆盖范围：用户登录流程、订单创建流程、支付流程、管理操作全流程

### 测试优先级

| 优先级 | 测试内容                                                         |
| :----- | :--------------------------------------------------------------- |
| **高** | 认证和权限系统、订单创建和支付流程、资金相关操作、核心 CRUD 操作 |
| **中** | 数据分析和报表、搜索和筛选功能、批量操作、用户交互体验           |
| **低** | 静态展示页面、非核心功能、视觉回归测试                           |

### 测试环境

1. **开发环境**：使用 Mock 数据，快速反馈，集成到 CI/CD 流程
2. **测试环境**：连接测试 API，真实数据测试，完整功能验证
3. **预生产环境**：生产数据镜像，性能测试，安全测试

---

## 📊 成功指标

### 功能指标

- ✅ 所有规划的功能模块完整实现
- ✅ 核心业务流程通过率：100%
- ✅ 权限控制正确率：100%

### 性能指标

- ⏱️ 页面首次加载时间 < 2 秒
- ⏱️ API 响应时间 < 500ms
- ⏱️ 列表渲染性能：1000 条数据无卡顿
- 📦 构建包大小 < 5MB（gzip）

### 质量指标

- 🧪 测试覆盖率 > 80%
- 🐛 Bug 密度 < 0.5/KLOC
- 📄 代码注释覆盖率 > 30%
- ♿ 可访问性评分 > 90（Lighthouse）

### 用户体验指标

- 🔍 用户操作成功率 > 95%
- 👁️ 界面一致性：100%
- 📱 移动端适配：覆盖主流设备
- 🎨 设计还原度 > 98%

---

## ⚠️ 风险评估与应对

| 风险等级   | 风险描述                                           | 应对策略                                             | 缓解措施                       |
| :--------- | :------------------------------------------------- | :--------------------------------------------------- | :----------------------------- |
| **高风险** | **API 接口变更**<br>后端接口不稳定导致前端频繁修改 | 建立接口契约，使用 Mock 数据先行开发                 | 增加集成测试覆盖               |
| **高风险** | **性能瓶颈**<br>大数据量下页面卡顿                 | 实施虚拟滚动、分页加载、数据缓存                     | 性能监控和优化工具集成         |
| **中风险** | **权限系统复杂**<br>权限配置错误导致安全隐患       | 权限测试工具，双重验证机制                           | 详细的权限变更审计日志         |
| **中风险** | **第三方库兼容性**<br>依赖库版本冲突               | 锁定版本号，定期更新测试                             | 使用 yarn resolutions 管理依赖 |
| **低风险** | **浏览器兼容性**<br>旧版浏览器不支持新特性         | 设置浏览器支持列表（Chrome、Firefox、Safari 最新版） | 使用 Polyfill 和 Babel 转译    |

---

## 📈 后续优化方向

1. **移动端适配**：开发移动端 H5 管理界面
2. **PWA 支持**：添加渐进式 Web 应用能力
3. **微前端架构**：将大型应用拆分为微前端
4. **低代码方案**：部分功能支持可视化配置
5. **AI 智能助手**：集成 AI 辅助管理功能

---
