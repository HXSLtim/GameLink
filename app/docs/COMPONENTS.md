# GameLink App 组件职责文档（Uni-app + Vue 3）

> 文档状态：✅ 已对齐（2026-01-31）  
> 适用范围：`app/src`（用户端 + 陪玩师端）

本文件描述 **当前 Vue 3 SFC 组件结构**，替换旧版小程序 WXML/WXSS 体系说明。

---

## 组件架构概览

```
app/src/components/
├── gl/                       # 基础 UI 组件（全局基础能力）
│   ├── Avatar/index.vue
│   ├── Button/index.vue
│   ├── Card/index.vue
│   ├── Empty/index.vue
│   ├── Input/index.vue
│   ├── Switch/index.vue
│   ├── Tag/index.vue
│   └── index.ts
├── layout/                   # 页面/列表布局
│   ├── BasePageLayout/index.vue
│   ├── ListHeaderStack/index.vue
│   └── PageShell/index.vue
├── GlobalToast/index.vue      # 全局提示
├── PageState/index.vue        # 页面状态（加载/空/错误/离线/登录）
├── CustomTabBar/index.vue     # 自定义 TabBar（PC/移动）
├── NavBar/index.vue           # 页面导航栏
├── SectionCard/index.vue      # 统一 Section 卡片容器
├── StatusCard/index.vue       # 统一状态卡片容器
└── ...                        # 业务组件（按域命名）
```

> 业务组件按领域命名（如 Player* / Order* / Wallet*），目录化结构统一为 `ComponentName/index.vue`。

---

## 一、基础 UI 组件（gl/）

| 组件 | 职责 |
| --- | --- |
| `gl/Avatar` | 头像展示，支持状态/徽章 |
| `gl/Button` | 按钮（变体/尺寸/加载/禁用） |
| `gl/Card` | 卡片容器（标题/边框/阴影） |
| `gl/Empty` | 空状态 |
| `gl/Input` | 输入框 |
| `gl/Switch` | 开关 |
| `gl/Tag` | 标签 |

---

## 二、布局组件（layout/）

| 组件 | 职责 |
| --- | --- |
| `BasePageLayout` | 页面基础布局（Nav + Header Stack + CustomTabBar） |
| `ListHeaderStack` | 列表页头部堆栈（搜索/筛选/标签等） |
| `PageShell` | 页面壳（滚动、刷新、ConfirmDialog 注入） |

---

## 三、业务组件（按模块）

**认证/账户**
`AuthLogo` / `LoginForm` / `RegisterForm` / `RoleSelector` / `AuthAgreementFooter` 等

**聊天/消息**
`ChatNavBar` / `ChatMessageList` / `ChatMessageBubble` / `ChatInputBar` / `MessageItem` 等

**订单/评价**
`OrderCard` / `OrderStatusCard` / `OrderActionBar` / `ReviewCard` / `ReviewModal` / `OrderCard/*` 子组件 等

**陪玩师**
`PlayerCard` / `PlayerDetailHeader` / `PlayerServicesSection` / `PlayerReviewsSection` / `PlayerCard/*` 子组件 等

**钱包/支付**
`WalletBalanceCard` / `PaymentMethodSelector` / `RechargeBalanceInfo` / `WalletRecordsSection` 等

**通用**
`SearchBar` / `FilterPanel` / `TabsBar` / `Skeleton` / `InfiniteList` / `RecordListItem` / `StatusInfo` / `HeaderStatsRow` / `PageState` / `SectionCard` / `StatusCard` / `PriceTag` 等

---

## 四、旧组件替代关系（重要）

| 旧 gl-* 组件 | 当前替代 |
| --- | --- |
| `gl-icon` | `uv-icon` |
| `gl-loading` | `uv-loading-icon` / `Skeleton` |
| `gl-search` | `SearchBar` |
| `gl-navbar` | `NavBar` |
| `gl-page` | `layout/PageShell` |
| `gl-section` | `SectionHeader` |
| `gl-modal` | `ConfirmDialog` |
| `gl-toast` | `GlobalToast` + `useToast` |
| `gl-tabs` | `TabsBar` |
| `gl-form` | `FormItem` / `FormSection` |

### FormItem / FormSection 增强说明

- `FormItem`：支持 `required`（红色星号指示）、`disabled`、`type='input'`（内部使用 `GlInput`）
- `FormSection`：支持 `#extra` 插槽（透传到 GlCard 头部右侧）
- `StatsCard`：支持 `columns` 属性自定义网格列数，支持 `unit` 字段和 `#extra` 插槽

---

## 五、组件规范

- **目录化**：业务组件统一 `ComponentName/index.vue`
- **脚本**：统一 `<script setup lang="ts">`
- **样式**：统一 CSS 变量令牌（`--font-*` / `--spacing-*` / `--radius-*`）
- **响应式**：使用 `@include desktop` / `desktop-lg`

---

## 六、已完成合并/替换

- `GameCardLarge` → `GameCard`（统一命名）  
- `ChatBubble` → `ChatMessageBubble`（统一气泡组件）  
- `PlayerOrderCard` → `OrderCard`（通过 `viewMode` 区分视角）  

---

如需补充组件详细属性说明，请以组件源码为准。
