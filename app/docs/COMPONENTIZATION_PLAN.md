# GameLink Uni-app 组件化封装计划

> 技术栈：uni-app + Vue 3 + TypeScript + Vite + Pinia + uv-ui + SCSS  
> 工作区：`app/src/`  
> 最后更新：2026-01-31  
> 状态：持续推进（组件化 + UI/UX 令牌化）

---

## 一、目标与范围

- 覆盖 `app` 端全部页面与组件，进行“二次封装 + 高度组件化”改造
- H5 端支持 PC（≥1024px）与移动端响应式布局
- 复用现有 UI/UX 与主题体系，减少维护成本
- 产出仓库内可执行的计划文档与页面使用计划（Markdown）

---

## 二、现状基线

- **响应式布局已存在**：`ResponsiveLayout.vue` 与 `PCLayout.vue` 共存
- **页面清单来源**：`pages.json`，TabBar 入口与 easycom 已配置
- **主题与混入**：
  - `styles/variables.scss`：设计令牌与主题变量
  - `styles/mixins.scss`：断点与交互 mixins（mobile/tablet/desktop、press-effect、hover-lift 等）

---

## 三、资产盘点（组件 & 页面）

### 3.1 覆盖范围

| 范围 | 说明 |
|------|------|
| 组件 | `app/src/components/` 下约 **125+** 个 `.vue` 组件 |
| 页面 | `pages.json` 中 **28** 个页面 |
| 布局 | `app/src/layouts/` 下 **2** 个布局组件 |
| 样式 | `styles/variables.scss` + `styles/mixins.scss` + `styles/index.scss` |

### 3.2 组件分层与清单（已归类）

#### Base（gl 二次封装）

- GlButton, GlCard, GlTag, GlAvatar, GlEmpty, GlInput, GlSwitch

#### Base（通用基础组件）

- ErrorBoundary, GlobalToast, PageState
- LazyImage, ImageUploader, Skeleton, LoadMore, VirtualList
- PrivacyPopup, RatingStars, PriceTag, ThemeToggle
- FormItem, FormSection

#### Layout / 导航

- BasePageLayout, PageShell
- ResponsiveLayout（PC/历史）, PCLayout（PC/历史）
- NavBar, CustomTabBar
- SectionHeader, QuickActions, MenuList, ProfileHeader

#### Pattern（可复用模块）

- SearchBar, FilterPanel, TabsBar, StatusInfo, HeaderStatsRow
- RankSelector, GameSelector, ServiceSelector
- PaymentMethodSelector, AmountSelector, QuantitySelector
- SchedulePicker, CouponSelector
- QuickQuestionBar, InfiniteList, ListItem, RecordListItem, HotGamesScroll

#### 业务 / 模块组件（按域）

- **认证/账户**：AuthLogo, AuthDivider, AuthAgreementFooter, LoginForm, LoginOtherActions, LoginAccountPopup, RegisterForm, RoleSelector, ProfileBasicSection, ProfileContactSection, ProfileGamesSection
- **聊天/消息**：ChatNavBar, ChatConnectionStatus, ChatMessageList, ChatMessageBubble, ChatInputBar, ChatMorePanel, ChatMenuPanel, MessageItem
- **订单/评价**：OrderCard, OrderStatusCard, OrderActionBar, OrderInfoSection, OrderFeeSection, OrderSubmitBar, OrderReviewSection, ReviewModal, ReviewCard, OrderQuickEntry
- **陪玩师/认证**：PlayerCard, PlayerDetailHeader, PlayerGamesSection, PlayerServicesSection, PlayerReviewsSection, PlayerActionBar, FavoriteListPanel, RecommendPlayersSection, GameCertSection, GameCertItem, CertificationBasicSection, CertStatusCard, IdCardUploader, AvatarUploader, VoiceSampleCard
- **钱包/支付**：WalletBalanceCard, WalletRecordsSection, WalletVipTip, RechargeBalanceInfo, RechargeActionBar, AgreementCheckRow, TipsListCard, TransactionItem
- **客服/帮助/协议**：SupportAgentStatus, ServiceStatsBar, ServiceList, ServiceChatMessage, ServiceEditorPanel, PlayerServiceCard, HelpCategoryGrid, HelpFaqSection, HelpContactCard, FaqItem, AgreementContent
- **运营/统计**：DashboardStatusCard, StatsCard, ResultCard, EarningsSummaryCard, EarningsChart, EarningsItem, EarningsDetailSection, ChannelCard
- **游戏**：GameCard, GameSelector, GameCertItem, GameCertSection

#### 单文件组件例外（已完成迁移）

- ✅ `GlobalToast.vue` → `components/GlobalToast/index.vue`
- ✅ `PageState.vue` → `components/PageState/index.vue`

#### 可合并/变体候选（提升复用度）

- `StatsCard` 与 `ResultCard` → 统一为统计/结果卡片变体

### 3.3 页面清单（来自 `pages.json`）

**用户端/公共页面（23）**  
`index`, `player/list`, `player/detail`, `message/list`, `message/chat`, `profile/index`, `profile/edit`, `order/list`, `order/detail`, `order/create`, `wallet/index`, `wallet/recharge`, `auth/login`, `auth/register`, `game/list`, `favorite/list`, `settings/index`, `review/list`, `channel/list`, `payment/result`, `agreement`, `help`, `service`

**陪玩师端页面（5）**  
`player/dashboard`, `player/orders`, `player/earnings`, `player/services`, `player/certification`

---

## 四、响应式策略

### 4.1 断点与 mixins

来自 `styles/mixins.scss`：

- Mobile：`< 768px`
- Tablet：`768px - 1023px`
- Desktop：`>= 1024px`
- Desktop LG：`>= 1440px`

### 4.2 布局容器策略

- **主入口**：页面统一以 `BasePageLayout`（`PageShell`）作为页面容器
- **兼容策略**：`ResponsiveLayout` / `PCLayout` 作为 PC 端三栏布局或历史页面保留
- **H5 / 小程序差异**：H5 使用 `window.innerWidth`，小程序用 `uni.getSystemInfoSync`
- **统一策略**：
  - 页面统一 slot 约定：`nav` / `search` / `banner` / `tabs` / `header-extra`
  - PC 端样式与断点保持与 `mixins.scss` 一致（≥ 1024px）

### 4.3 组件级响应式规范

- 组件统一 `size` / `variant`，在 PC 端默认使用更舒展的密度
- 列表类组件支持 `compact`/`expanded` 变体
- 交互效果统一使用 `press-effect`、`hover-lift` mixins

---

## 五、组件封装规范

### 5.1 分层职责与边界

| 层级 | 职责 | 示例 | 规则 |
|------|------|------|------|
| Base | 对 uv-ui 二次封装 + 通用基础 UI | GlButton, GlCard, EmptyState | 无业务耦合、样式令牌化 |
| Pattern | 页面级可复用模块 | SearchBar, TabsBar, FilterPanel | 复用度高、输入输出明确 |
| Business | 业务卡片/区块 | PlayerCard, OrderCard, MessageItem | 绑定业务模型与展示规则 |

### 5.2 目录结构

当前目录为“扁平 + gl 子目录”，计划逐步演进到分层目录（可渐进实施）：

```
components/
├── gl/                   # Base（uv-ui 二次封装）
├── layout/               # Layout（可选物理分组）
├── pattern/              # Pattern（可选物理分组）
└── business/             # Business（可选物理分组）
```

### 5.3 命名与 easycom

- Base 使用 `Gl` 前缀（如 `GlButton`）
- 其它组件使用 PascalCase
- easycom 已配置：
  - `^gl-(.*)` 与 `^Gl([A-Z].*)` 指向 `components/gl`

### 5.4 Props / Events / Slots

```ts
interface Props {
  title: string
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  disabled: false,
})

const emit = defineEmits<{
  click: [event: Event]
  change: [value: string]
}>()
```

### 5.5 样式规范

- 统一使用设计令牌（`styles/variables.scss`）
- 避免硬编码，优先使用 `var(--color-*)`、`var(--spacing-*)`、`var(--radius-*)`
- 交互统一使用 mixins：`press-effect` / `hover-lift` / `text-ellipsis`

---

## 六、实施步骤（与 Todo 对齐）

1. 组件资产盘点与分层归类（✅ 已完成）
2. 响应式断点与布局容器策略（✅ 已完成）
3. Base/Pattern/Business 组件封装规范（✅ 已完成）
4. 页面组件使用计划（✅ 已完成，见 `PAGE_COMPONENT_USAGE_PLAN.md`）
5. 文档交付与交叉引用（✅ 已完成，本文件 + 页面计划）

---

## 七、风险与注意事项

| 风险 | 影响 | 缓解措施 |
|------|------|---------|
| 一次性全量重构 | 高 | 采用增量替换，保留兼容路径 |
| easycom 命名冲突 | 中 | 统一 `Gl` 前缀，避免与 `uv-` 冲突 |
| H5 与小程序差异 | 中 | 保留条件编译，避免小程序受影响 |
| 单文件组件混杂 | 低 | 按阶段迁移到目录结构 |

---

## 八、相关文档

- [组件职责文档](./COMPONENTS.md)
- [页面组件使用计划](./PAGE_COMPONENT_USAGE_PLAN.md)
- [前端重构计划](./REFACTOR_PLAN.md)
- [主题变量与 mixins](../src/styles/variables.scss)
