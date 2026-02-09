# GameLink 页面组件使用计划

> 基于 `pages.json` 页面清单，为每个页面规划组件使用与抽取方案  
> 工作区：`app/src/`  
> 最后更新：2026-02-06  
> 状态：已覆盖全量页面，持续同步代码

---

## 一、组件使用总览

### 1.1 分层与职责

| 层级 | 职责 | 组件示例 |
|------|------|----------|
| Layout | 页面骨架与导航 | `BasePageLayout`, `PageShell`, `NavBar`, `CustomTabBar` |
| Base | 通用基础组件 | `PageState`, `LoadMore`, `Skeleton`, `GlEmpty` |
| Pattern | 页面级模块 | `SearchBar`, `TabsBar`, `FilterPanel`, `InfiniteList`, `ListItem` |
| Business | 业务卡片/区块 | `PlayerCard`, `OrderCard`, `MessageItem`, `ReviewCard` |

### 1.2 通用组合（建议）

- 列表页：`NavBar` + `SearchBar/TabsBar` + 业务卡片 + `LoadMore` + `PageState`
- 详情页：`NavBar` + 业务区块 + 操作区 + `PageState`
- 表单页：`NavBar` + 表单区块 + 提交按钮（统一 `GlButton`）
- 聊天页：`ChatNavBar` + `ChatMessageList` + `ChatInputBar`
- TabBar 页：`CustomTabBar` + 页面内容

---

## 二、全量页面清单（28）

| 页面 | 路径 | 类型 | 关键组件 | 状态 |
|------|------|------|----------|------|
| 首页 | `pages/index/index` | TabBar | `NavBar`, `SearchBar`, `OfflineBanner`, `HotGamesScroll`, `RecommendPlayersSection` | 已组件化 |
| 陪玩师列表 | `pages/player/list/index` | 列表页 | `SearchBar`, `FilterPanel`, `PlayerCard` | 已组件化 |
| 陪玩师详情 | `pages/player/detail/index` | 详情页 | `PlayerDetailHeader`, `PlayerReviewsSection` | 已组件化 |
| 消息列表 | `pages/message/list/index` | TabBar | `TabsBar`, `MessageItem` | 已组件化 |
| 聊天详情 | `pages/message/chat/index` | 聊天页 | `ChatMessageList`, `ChatInputBar` | 已组件化 |
| 个人中心 | `pages/profile/index/index` | TabBar | `ProfileHeader`, `MenuList`, `ThemeToggle` | 已组件化 |
| 编辑资料 | `pages/profile/edit/index` | 表单页 | `ProfileBasicSection`, `AvatarUploader` | 已组件化 |
| 订单列表 | `pages/order/list/index` | 列表页 | `TabsBar`, `OrderCard`, `InfiniteList` | 已组件化 |
| 订单详情 | `pages/order/detail/index` | 详情页 | `OrderInfoSection`, `OrderActionBar` | 已组件化 |
| 下单页 | `pages/order/create/index` | 表单页 | `ServiceSelector`, `SchedulePicker`, `SectionCard` | 已组件化 |
| 钱包 | `pages/wallet/index/index` | 列表页 | `WalletBalanceCard`, `WalletRecordsSection` | 已组件化 |
| 充值 | `pages/wallet/recharge/index` | 表单页 | `RechargeBalanceInfo`, `RechargeActionBar` | 已组件化 |
| 登录 | `pages/auth/login/index` | 表单页 | `LoginForm`, `AuthAgreementFooter` | 已组件化（已接入 GlInput） |
| 注册 | `pages/auth/register/index` | 表单页 | `RegisterForm`, `RoleSelector` | 已组件化（已接入 GlInput） |
| 游戏列表 | `pages/game/list/index` | 列表页 | `SearchBar`, `FilterPanel`, `GameCard` | 已组件化 |
| 收藏列表 | `pages/favorite/list/index` | 列表页 | `FavoriteListPanel`, `FavoriteEditBar` | 已组件化 |
| 设置 | `pages/settings/index/index` | 设置页 | `SettingsSection`, `ThemeToggle`, `GlSwitch` | ✅ 已组件化 |
| 评价列表 | `pages/review/list/index` | 列表页 | `ReviewCard`, `RatingStars` | 已组件化 |
| 工作台 | `pages/player/dashboard/index` | 仪表盘 | `DashboardStatusCard`, `StatsCard` | 已组件化 |
| 陪玩师订单 | `pages/player/orders/index` | 列表页 | `TabsBar`, `OrderCard` | 已组件化 |
| 收益中心 | `pages/player/earnings/index` | 详情页 | `EarningsSummaryCard`, `EarningsChart` | 已组件化 |
| 服务管理 | `pages/player/services/index` | 列表/表单 | `ServiceList`, `ServiceEditorPanel`, `GlInput` | ✅ 已组件化 |
| 陪玩认证 | `pages/player/certification/index` | 表单页 | `GameCertSection`, `IdCardUploader`, `FormItem(GlInput)` | ✅ 已组件化 |
| 公共频道 | `pages/channel/list/index` | 列表页 | `SearchBar`, `FilterPanel`, `ChannelCard` | 已组件化 |
| 支付结果 | `pages/payment/result/index` | 结果页 | `ResultCard`, `SectionCard` | 已组件化 |
| 协议 | `pages/agreement/index` | 文档页 | `AgreementContent` | 已组件化 |
| 帮助中心 | `pages/help/index` | 列表页 | `HelpCategoryGrid`, `HelpFaqSection` | 已组件化 |
| 在线客服 | `pages/service/index` | 聊天页 | `ServiceChatMessage`, `QuickQuestionBar` | 已组件化 |

---

## 三、页面使用计划（按模块）

### 3.1 首页与入口

#### `pages/index/index`
- 类型：TabBar 首页
- 组件：`NavBar`, `SearchBar`, `OfflineBanner`, `HotGamesScroll`, `RecommendPlayersSection`, `CustomTabBar`

#### `pages/channel/list/index`
- 类型：列表页
- 组件：`SearchBar`, `FilterPanel`, `ChannelCard`, `InfiniteList`, `PageState`

### 3.2 陪玩师

#### `pages/player/list/index`
- 类型：列表页 / TabBar
- 组件：`SearchBar`, `FilterPanel`, `PlayerCard`, `InfiniteList`, `PageState`, `CustomTabBar`

#### `pages/player/detail/index`
- 类型：详情页
- 组件：`PlayerDetailHeader`, `PlayerGamesSection`, `PlayerServicesSection`, `PlayerReviewsSection`, `PlayerActionBar`, `PageState`

#### `pages/player/certification/index`
- 类型：表单页
- 组件：`CertificationBasicSection`, `GameCertSection`, `GameCertItem`, `VoiceSampleCard`, `IdCardUploader`, `AvatarUploader`, `CertStatusCard`, `PageState`
- ✅ 输入类组件已统一为 `GlInput`（通过 `FormItem` 封装）

### 3.3 订单与评价

#### `pages/order/list/index`
- 类型：列表页
- 组件：`TabsBar`, `OrderCard`, `InfiniteList`, `PageState`

#### `pages/order/detail/index`
- 类型：详情页
- 组件：`OrderStatusCard`, `PlayerCard`, `OrderInfoSection`, `OrderFeeSection`, `OrderActionBar`, `OrderReviewSection`, `ReviewModal`, `PageState`

#### `pages/order/create/index`
- 类型：表单页
- 组件：`PlayerCard`, `ServiceSelector`, `GameSelector`, `SchedulePicker`, `QuantitySelector`, `CouponSelector`, `OrderSubmitBar`, `SectionCard`

#### `pages/review/list/index`
- 类型：列表页
- 组件：`NavBar`, `ReviewCard`, `RatingStars`, `LoadMore`, `PageState`

### 3.4 消息与聊天

#### `pages/message/list/index`
- 类型：列表页 / TabBar
- 组件：`TabsBar`, `MessageItem`, `InfiniteList`, `PageState`, `CustomTabBar`

#### `pages/message/chat/index`
- 类型：聊天页
- 组件：`ChatNavBar`, `ChatConnectionStatus`, `ChatMessageList`, `ChatMessageBubble`, `ChatInputBar`, `ChatMorePanel`, `ChatMenuPanel`

### 3.5 钱包与支付

#### `pages/wallet/index/index`
- 类型：详情页 + 列表页
- 组件：`WalletBalanceCard`, `WalletVipTip`, `WalletRecordsSection`, `TransactionItem`, `InfiniteList`, `PageState`

#### `pages/wallet/recharge/index`
- 类型：表单页
- 组件：`RechargeBalanceInfo`, `AmountSelector`, `PaymentMethodSelector`, `AgreementCheckRow`, `TipsListCard`, `RechargeActionBar`

#### `pages/payment/result/index`
- 类型：结果页
- 组件：`ResultCard`, `SectionCard`

### 3.6 认证与账户

#### `pages/auth/login/index`
- 类型：表单页
- 组件：`AuthLogo`, `LoginForm`, `AuthDivider`, `LoginOtherActions`, `LoginAccountPopup`, `AuthAgreementFooter`, `PrivacyPopup`
- ✅ 输入类组件已统一为 `GlInput`

#### `pages/auth/register/index`
- 类型：表单页
- 组件：`AuthLogo`, `RegisterForm`, `RoleSelector`, `AuthAgreementFooter`
- ✅ 输入类组件已统一为 `GlInput`

#### `pages/profile/index/index`
- 类型：TabBar 页面
- 组件：`ProfileHeader`, `QuickActions`, `MenuList`, `ThemeToggle`, `CustomTabBar`

#### `pages/profile/edit/index`
- 类型：表单页
- 组件：`ProfileBasicSection`, `ProfileContactSection`, `ProfileGamesSection`, `IntroTextCard`, `AvatarUploader`

### 3.7 设置 / 帮助 / 协议 / 客服

#### `pages/settings/index/index`
- 类型：设置页
- 组件：`SettingsSection`, `ThemeToggle`, `MenuList`
- ✅ 开关类组件已统一为 `GlSwitch`

#### `pages/help/index`
- 类型：帮助中心
- 组件：`HelpCategoryGrid`, `HelpFaqSection`, `HelpContactCard`, `FaqItem`

#### `pages/agreement/index`
- 类型：文档页
- 组件：`AgreementContent`

#### `pages/service/index`
- 类型：客服页
- 组件：`SupportAgentStatus`, `QuickQuestionBar`, `ServiceChatMessage`

### 3.8 陪玩师端

#### `pages/player/dashboard/index`
- 类型：仪表盘
- 组件：`DashboardStatusCard`, `StatsCard`, `QuickActions`, `MenuList`, `PageState`

#### `pages/player/orders/index`
- 类型：列表页
- 组件：`TabsBar`, `OrderCard`, `InfiniteList`, `PageState`

#### `pages/player/earnings/index`
- 类型：详情页 + 列表页
- 组件：`EarningsSummaryCard`, `FilterPanel`, `EarningsChart`, `EarningsDetailSection`, `EarningsItem`, `LoadMore`, `PageState`

#### `pages/player/services/index`
- 类型：列表页 + 表单页
- 组件：`ServiceStatsBar`, `ServiceList`, `ServiceEditorPanel`, `PlayerServiceCard`, `GameSelector`, `RankSelector`, `PriceTag`, `PageState`
- ✅ 输入类组件已统一为 `GlInput`

---

## 四、已完成的统一与优化

1. ✅ **输入类 Base 组件**：`GlInput` 已覆盖所有表单页（RegisterForm、ServiceEditorPanel、FormItem、IntroTextCard 等）
2. ✅ **开关类 Base 组件**：`GlSwitch` 已在 SettingsSection 中统一使用
3. ✅ **布局合并**：统一使用 `BasePageLayout/PageShell`，`ResponsiveLayout` 仅用于 PC 专用布局
4. ✅ **组件目录结构**：已完成 `GlobalToast` / `PageState` 目录化迁移
5. ✅ **FormItem required 属性**：已补齐 `required` 属性，支持红色星号指示
6. ✅ **FormSection 插槽透传**：已修复 `#extra` 插槽透传到 GlCard
7. ✅ **StatsCard 灵活化**：支持 `columns` 属性自定义列数、`unit` 字段、`#extra` 插槽
8. ℹ️ **StatsCard / ResultCard**：保持独立（职责不同：统计网格 vs 状态结果展示）
9. ✅ **主题色一致**：修复 ChatConnectionStatus、RatingStars、PlayerDetailHeader 硬编码 → CSS 变量
10. ✅ **Emoji → uv-icon**：全量替换 emoji 图标为 uv-icon（StatusCard presets、StatusInfo、ErrorBoundary、LazyImage、HelpContactCard、PlayerActionBar、ChatMorePanel、ChatInputBar、PlayerMeta、RatingStars、ReviewModal、RoleSelector、QuickActions、QuickQuestionBar、HelpCategoryGrid、PaymentMethodSelector、EarningsItem、ResponsiveLayout、useRegister、usePlayerDashboard、useHelp、useWallet、useRecharge、useCustomerService）
11. ✅ **硬编码颜色清理**：TransactionItem、OfflineBanner、GlobalToast、CustomTabBar、ResponsiveLayout、login 页微信按钮 → CSS 变量；新增 `--color-wechat` 设计令牌
12. ✅ **console.log 清理**：移除 useWebSocket、App.vue 中的调试日志（保留 console.error/console.warn）
13. ✅ **硬编码尺寸清理**：EarningsItem、ChatNavBar 等组件的像素值 → CSS 变量（`--spacing-*`、`--font-*`、`--radius-*`）

---

## 五、相关文档

- [组件化封装计划](./COMPONENTIZATION_PLAN.md)
- [组件职责文档](./COMPONENTS.md)
- [前端重构计划](./REFACTOR_PLAN.md)
- [Emoji/图标使用检查报告](./EMOJI_ICON_AUDIT.md)（使用 emoji 的组件与统一为 uv-icon 的建议）
