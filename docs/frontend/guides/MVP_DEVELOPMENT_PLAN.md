# GameLink 陪玩管理平台 - MVP 开发计划

## 📋 项目概述

GameLink 是一个完整的陪玩管理平台，包含用户端、陪玩师端和管理后台三大模块。

## 🎯 当前状态

### ✅ 已完成功能

1. **基础架构**
   - ✅ React + TypeScript + Vite 项目搭建
   - ✅ Neo-brutalism 黑白设计系统
   - ✅ 自定义组件库（15+ 组件）
   - ✅ 主题切换系统（亮色/暗色）
   - ✅ 国际化支持（中英文）
   - ✅ 骨架屏加载动画

2. **认证系统**
   - ✅ Mock 登录功能
   - ✅ AuthContext 状态管理
   - ✅ 路由守卫
   - ✅ Token 持久化

3. **管理后台基础**
   - ✅ Dashboard 仪表盘
   - ✅ 订单列表页面（搜索、筛选、分页）
   - ✅ 导航和侧边栏
   - ✅ 面包屑导航

## 🚀 MVP 阶段开发计划

### 第一阶段：完善订单管理 (Week 1-2)

#### 1. 订单详情页面 ⏳ 进行中

**优先级**: 🔴 高

**功能点**:

- 订单基本信息展示
- 用户和陪玩者信息卡片
- 订单状态时间线
- 操作历史展示
- 审核状态和记录
- 操作按钮（审核、取消、退款等）

**文件清单**:

```
src/pages/Orders/
├── OrderDetail.tsx           # 订单详情页面
├── OrderDetail.module.less   # 页面样式
├── components/
│   ├── OrderInfo.tsx         # 订单信息卡片
│   ├── UserCard.tsx          # 用户/陪玩者卡片
│   ├── OrderTimeline.tsx     # 状态时间线
│   ├── OrderActions.tsx      # 操作按钮区
│   └── ReviewSection.tsx     # 审核区域
└── index.ts
```

#### 2. 订单审核工作流 ⏸️ 待开始

**优先级**: 🔴 高

**功能点**:

- 审核表单组件
- 审核意见填写
- 审核结果提交
- 状态实时更新
- 审核历史记录

**文件清单**:

```
src/components/ReviewModal/
├── ReviewModal.tsx
├── ReviewModal.module.less
└── index.ts
```

#### 3. 订单溯源功能 ⏸️ 待开始

**优先级**: 🟡 中

**功能点**:

- 操作日志列表
- 时间线可视化
- 操作者信息展示
- 状态变更记录
- 数据导出功能

---

### 第二阶段：基础数据管理 (Week 3-4)

#### 4. 游戏管理模块 ⏸️ 待开始

**优先级**: 🔴 高

**功能点**:

- 游戏列表展示
- 添加/编辑游戏
- 游戏图片上传
- 标签管理
- 启用/禁用状态

**文件清单**:

```
src/pages/Games/
├── GameList.tsx
├── GameList.module.less
├── GameForm.tsx              # 添加/编辑表单
├── components/
│   ├── GameCard.tsx
│   └── GameFormModal.tsx
└── index.ts
```

**数据结构**:

```typescript
interface Game {
  id: string;
  name: string;
  icon: string;
  category: string;
  tags: string[];
  status: 'active' | 'inactive';
  playerCount: number;
  createdAt: string;
  updatedAt: string;
}
```

#### 5. 陪玩师管理模块 ⏸️ 待开始

**优先级**: 🔴 高

**功能点**:

- 陪玩师列表
- 审核流程
- 技能标签管理
- 等级和评分
- 收益统计

**文件清单**:

```
src/pages/Players/
├── PlayerList.tsx
├── PlayerDetail.tsx
├── PlayerReview.tsx          # 审核页面
├── components/
│   ├── PlayerCard.tsx
│   ├── SkillTags.tsx
│   └── PlayerStats.tsx
└── index.ts
```

#### 6. 用户管理模块 ⏸️ 待开始

**优先级**: 🟡 中

**功能点**:

- 用户列表
- 用户详情
- 状态管理（正常/封禁）
- 订单历史
- 消费统计

**文件清单**:

```
src/pages/Users/
├── UserList.tsx
├── UserDetail.tsx
└── index.ts
```

---

### 第三阶段：财务和报表 (Week 5-6)

#### 7. 支付管理模块 ⏸️ 待开始

**优先级**: 🟡 中

**功能点**:

- 支付记录列表
- 支付状态查看
- 退款处理
- 对账功能
- 财务统计

**文件清单**:

```
src/pages/Payments/
├── PaymentList.tsx
├── PaymentDetail.tsx
├── RefundModal.tsx
└── index.ts
```

#### 8. 数据报表模块 ⏸️ 待开始

**优先级**: 🟡 中

**功能点**:

- 收入趋势图表
- 订单统计分析
- 用户增长曲线
- 陪玩师排行
- 数据导出

**文件清单**:

```
src/pages/Reports/
├── RevenueReport.tsx
├── OrderReport.tsx
├── UserReport.tsx
└── components/
    ├── LineChart.tsx
    ├── BarChart.tsx
    └── PieChart.tsx
```

---

### 第四阶段：系统配置 (Week 7-8)

#### 9. 权限管理模块 ⏸️ 待开始

**优先级**: 🟡 中

**功能点**:

- 角色列表
- 权限分配
- 管理员账户
- 操作日志
- 安全审计

**文件清单**:

```
src/pages/Permissions/
├── RoleList.tsx
├── PermissionMatrix.tsx
├── AdminList.tsx
└── OperationLog.tsx
```

#### 10. 系统设置模块 ⏸️ 待开始

**优先级**: 🟢 低

**功能点**:

- 平台参数配置
- 佣金比例设置
- 客服配置
- 系统维护
- 通知设置

**文件清单**:

```
src/pages/Settings/
├── GeneralSettings.tsx
├── CommissionSettings.tsx
├── NotificationSettings.tsx
└── MaintenanceSettings.tsx
```

---

## 📊 开发进度追踪

### 当前冲刺（Sprint 1）

**时间**: Week 1-2  
**目标**: 完善订单管理核心功能

- [x] 订单列表页面
- [ ] 订单详情页面（进行中）
- [ ] 订单审核工作流
- [ ] 订单溯源功能

### 里程碑

| 里程碑           | 目标日期 | 状态      | 完成度 |
| ---------------- | -------- | --------- | ------ |
| M1: 订单管理完整 | Week 2   | 🟡 进行中 | 60%    |
| M2: 数据管理完成 | Week 4   | ⏸️ 未开始 | 0%     |
| M3: 财务报表就绪 | Week 6   | ⏸️ 未开始 | 0%     |
| M4: MVP 上线     | Week 8   | ⏸️ 未开始 | 0%     |

---

## 🛠️ 技术栈

### 前端核心

- **框架**: React 18.3+
- **语言**: TypeScript 5.6+
- **构建**: Vite 5.4+
- **路由**: React Router 6.27+
- **样式**: Less + CSS Modules

### 组件库

- **自定义组件**: 15+ Neo-brutalism 风格组件
- **图表**: 待定（考虑 Recharts 或 Chart.js）
- **表单**: 自定义 Form 组件
- **表格**: 自定义 Table 组件

### 状态管理

- **Context API**: 认证、主题、国际化
- **Local State**: 页面级状态

### 工具库

- **日期处理**: 考虑引入 dayjs
- **请求**: 自定义 http 客户端
- **工具函数**: 自定义 utils

---

## 📁 项目结构

```
src/
├── components/           # 共享组件（15+）
│   ├── Button/
│   ├── Input/
│   ├── Table/
│   ├── Modal/
│   ├── Skeleton/
│   └── ...
├── pages/               # 页面组件
│   ├── Dashboard/       ✅ 已完成
│   ├── Orders/          🟡 进行中
│   │   ├── OrderList    ✅
│   │   └── OrderDetail  ⏸️
│   ├── Games/           ⏸️ 待开始
│   ├── Players/         ⏸️ 待开始
│   ├── Users/           ⏸️ 待开始
│   ├── Payments/        ⏸️ 待开始
│   ├── Reports/         ⏸️ 待开始
│   ├── Permissions/     ⏸️ 待开始
│   └── Settings/        ⏸️ 待开始
├── contexts/            # 全局状态
│   ├── AuthContext      ✅
│   ├── ThemeContext     ✅
│   └── I18nContext      ✅
├── hooks/               # 自定义 Hooks
├── services/            # API 服务
├── types/               # 类型定义
├── utils/               # 工具函数
└── i18n/               # 国际化
```

---

## 🎨 设计规范

### 颜色系统

- **主色**: 纯黑 (#000000) / 纯白 (#FFFFFF)
- **灰度**: 10 级灰度系统
- **主题**: 支持亮色/暗色切换

### 组件规范

- **无圆角**: border-radius: 0
- **实体阴影**: 8px x 8px black
- **粗边框**: 2px solid
- **高对比度**: 黑白分明

### 响应式

- **桌面**: 1920x1080 主要适配
- **平板**: 768x1024 适配
- **移动**: 375x667 基准

---

## 📝 开发规范

### 代码质量

- ✅ TypeScript 严格模式
- ✅ ESLint + Prettier
- ✅ 100% 类型安全
- ✅ 组件测试（Vitest）

### Git 工作流

```
main (生产)
  ├── develop (开发)
  │   ├── feature/order-detail
  │   ├── feature/game-management
  │   └── ...
```

### 提交规范

```
feat: 新功能
fix: 修复bug
docs: 文档更新
style: 代码格式
refactor: 重构
test: 测试
chore: 构建/工具
```

---

## 🚀 下一步行动

### 立即开始

1. **创建订单详情页面** 🔥
   - 设计页面布局
   - 创建子组件
   - 集成 Mock 数据
   - 添加交互逻辑

2. **实现审核工作流**
   - 设计审核流程
   - 创建审核表单
   - 状态流转逻辑

3. **完善溯源功能**
   - 操作日志展示
   - 时间线组件

### 本周目标

- [ ] 完成订单详情页面 (3天)
- [ ] 实现审核工作流 (2天)
- [ ] 测试和优化 (1天)

---

**更新日期**: 2025-01-05  
**版本**: v1.0.0-mvp  
**负责人**: 开发团队
