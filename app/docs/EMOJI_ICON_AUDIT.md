# App 端 Emoji / 图标使用检查报告

本文档列出 `app/src` 中所有**使用 emoji 作为图标或装饰**的组件与数据源，并给出统一为 `uv-icon` 或设计令牌的建议，便于多端一致性与可访问性。

**状态**：已按本报告完成替换（emoji → uv-icon），仅保留：PlayerDetailHeader 的 ♂♀（标准 Unicode 符号）。RatingStars / ReviewModal 的星星已改为 uv-icon（star-fill / star）。

---

## 一、清单总览

| 类别 | 文件/位置 | 用途 | 建议 |
|------|------------|------|------|
| **预设数据** | `StatusCard/presets.ts` | 订单/认证/结果状态 icon | 可改为 uv-icon 名 + StatusInfo 支持图标类型 |
| **UI 组件** | `StatusInfo` | 展示 `icon` 字符串（多为 emoji） | 支持 `iconType: 'emoji' \| 'uv-icon'`，按类型渲染 |
| **UI 组件** | `ResultCard` | 依赖 presets 的 icon（emoji） | 随 presets 改为图标名后自动统一 |
| **UI 组件** | `ErrorBoundary` | 错误类型对应 emoji | 建议改为 uv-icon（见下表） |
| **UI 组件** | `LazyImage` | 占位/错误态 emoji | 建议改为 uv-icon |
| **UI 组件** | `HelpContactCard` | 联系图标 💬 | 建议改为 uv-icon `chat` |
| **UI 组件** | `PlayerActionBar` | 收藏 ❤️/🤍、私信 💬 | 建议改为 uv-icon |
| **UI 组件** | `PlayerCard/PlayerMeta` | 评分 ⭐ | 可与 RatingStars 统一为 ★/☆ 或图标 |
| **UI 组件** | `ChatMorePanel` | 图片 🖼️、警告 ⚠️ | 建议改为 uv-icon |
| **UI 组件** | `ChatInputBar` | 键盘 ⌨️、麦克风 🎤 | 建议改为 uv-icon |
| **UI 组件** | `PlayerDetailHeader` | 性别 ♂/♀ | 可保留或改为小图标 |
| **UI 组件** | `RatingStars` / `ReviewModal` | ★ / ☆ | 已与主题色统一，可保留或换图标集 |
| **布局/数据** | `ResponsiveLayout` | PC 主导航 mainNavItems icon（emoji） | 建议改为 uv-icon 名 + 左侧 nav 按名渲染 |
| **Composables** | `useRegister` | 角色选择 👤/🎮 | 若 UI 用 uv-icon，改为图标名 |
| **Composables** | `usePlayerDashboard` | 订单/收益/服务等入口 icon | 同上 |
| **Composables** | `useHelp` | 帮助分类 icon | 同上 |
| **Composables** | `useWallet` | 钱包入口 icon | 同上 |
| **Composables** | `useRecharge` | 支付方式 icon | 同上 |
| **Composables** | `useCustomerService` | 客服入口 icon | 同上 |
| **其他** | `SupportAgentStatus` | GlAvatar 兜底文字 👩‍💼 | 可改为首字或图标 |

---

## 二、按文件详细列表

### 1. 组件内联 Emoji（模板中直接写死）

| 文件 | 位置/含义 | Emoji | 建议 uv-icon 替代 |
|------|------------|--------|-------------------|
| `ErrorBoundary/index.vue` | 网络错误 | 📡 | `wifi-off` 或 `error-circle` |
| `ErrorBoundary/index.vue` | 服务器错误 | 🔧 | `setting` 或 `error-circle` |
| `ErrorBoundary/index.vue` | 数据加载失败 | 📋 | `list` 或 `file-text` |
| `ErrorBoundary/index.vue` | 权限不足 | 🔒 | `lock` |
| `ErrorBoundary/index.vue` | 404 | 🔍 | `search` 或 `file-delete` |
| `ErrorBoundary/index.vue` | 通用错误 | 😵 | `error-circle` |
| `LazyImage/index.vue` | 占位 | 🖼️ | `image` 或 `photo` |
| `LazyImage/index.vue` | 加载失败 | ⚠️ | `error-circle` 或 `info-circle` |
| `HelpContactCard/index.vue` | 联系图标 | 💬 | `chat` |
| `PlayerActionBar/index.vue` | 收藏/未收藏 | ❤️ / 🤍 | `heart-fill` / `heart` |
| `PlayerActionBar/index.vue` | 私信 | 💬 | `chat` |
| `PlayerCard/PlayerMeta.vue` | 评分 | ⭐ | 保留或 `star-fill` |
| `ChatMorePanel/index.vue` | 图片、警告 | 🖼️、⚠️ | `image`、`info-circle` |
| `ChatInputBar/index.vue` | 键盘/麦克风切换 | ⌨️ / 🎤 | `keyboard` / `mic`（需确认 uv 是否有） |
| `PlayerDetailHeader/index.vue` | 性别 | ♂ / ♀ | 可保留或 `male`/`female` 图标 |
| `RatingStars/index.vue` | 星标 | ★ / ☆ | 已用主题色，可保留 |
| `ReviewModal/index.vue` | 星标 | ★ | 同上 |
| `SupportAgentStatus/index.vue` | 客服头像兜底 | 👩‍💼 | 可改为首字或 `account` |

### 2. 预设/配置中的 Emoji（数据驱动）

| 文件 | 用途 | Emoji 列表 | 建议 |
|------|------|------------|------|
| `StatusCard/presets.ts` | 订单状态 | ⏳✅🎮🎉🔄💰❌⚠️ | 统一为 uv-icon 名，如 `clock`、`checkmark-circle`、`gamepad` 等 |
| `StatusCard/presets.ts` | 认证状态 | 📝⏳✅❌ | 同上 |
| `StatusCard/presets.ts` | 结果状态 | ✅⏳❌⚠️ | 同上 |
| `ResponsiveLayout.vue` | PC 主导航 mainNavItems | 🏠🎮💬📨📋💰 | 改为 icon 名，左侧 nav 用 uv-icon 渲染 |
| `useRegister.ts` | 注册角色选项 | 👤、🎮 | 若列表用图标展示，改为 icon 名 |
| `usePlayerDashboard.ts` | 陪玩工作台入口 | 📋💰🎮📅 | 同上 |
| `useHelp.ts` | 帮助分类 | 📋💳🔐📖🎮❓ | 同上 |
| `useWallet.ts` | 钱包入口 | 🎫📊🎁❓ | 同上 |
| `useRecharge.ts` | 支付方式 | 💚💙🍎 | 同上 |
| `useCustomerService.ts` | 客服入口 | 💳📋👤⚠️ | 同上 |
| `EarningsItem/index.vue` | 收益类型 icons 映射 | 💰💳🎁 | 改为 uv-icon 名，组件内用 uv-icon 渲染 |

---

## 三、与现有 uv-icon 用法的一致性

- **已使用 uv-icon 的组件**：ChatConnectionStatus、FormItem、PaymentMethodSelector、PlayerCard、OrderQuickEntry、ChatNavBar、GlAvatar、PageState、RegisterForm、ServiceEditorPanel、TransactionItem、WalletBalanceCard、CouponSelector、HelpContactCard（箭头）、CustomTabBar（PC 侧）、EarningsDetailSection、ReviewCard、SettingsSection、MenuList、ChannelCard、LoginAccountPopup、ChatMessageBubble、GameCertItem、GameSelector、AvatarUploader、PlayerServicesSection、IdCardUploader、ServiceSelector、SearchBar、GlobalToast 等。
- **移动端 TabBar**：CustomTabBar 已使用本地 SVG（`/static/icons/*.svg`），未使用 emoji。
- **不一致点**：同一产品内部分为 emoji、部分为 uv-icon，不同系统/字体下 emoji 展示不一，且不利于无障碍（屏幕阅读器对 emoji 的解读不统一）。建议逐步将「状态、结果、错误、导航入口」等处的 emoji 改为 uv-icon 或本地图标。

---

## 四、建议实施顺序（可选）

1. **StatusCard presets + StatusInfo**  
   - presets 中 `icon` 改为 uv-icon 的 `name`（如 `checkmark-circle-fill`、`clock`、`close-circle` 等）。  
   - StatusInfo 增加 `iconType: 'emoji' | 'uv-icon'`（或根据 `icon` 是否像 emoji 自动判断），`iconType === 'uv-icon'` 时用 `<uv-icon :name="icon" />` 渲染。  
   - ResultCard 无需改组件，只依赖 presets 即可统一。

2. **ErrorBoundary / LazyImage**  
   - 按上表将各类型 emoji 改为对应 uv-icon，样式用现有 CSS 变量（如 `--color-error`）。

3. **ResponsiveLayout 主导航**  
   - mainNavItems 的 `icon` 改为 uv-icon 名，模板中左侧 nav 用 `<uv-icon :name="item.icon" />` 替代 `<text>{{ item.icon }}</text>`。

4. **业务组件**  
   - HelpContactCard、PlayerActionBar、ChatMorePanel、ChatInputBar、PlayerMeta、SupportAgentStatus 等，按上表逐个替换为 uv-icon，并保持颜色使用主题变量。

5. **Composables 与列表类页面**  
   - useRegister、usePlayerDashboard、useHelp、useWallet、useRecharge、useCustomerService、EarningsItem 等，将配置中的 `icon` 改为 uv-icon 名，对应列表/卡片用 `<uv-icon>` 渲染。

---

## 五、备注

- **RatingStars / ReviewModal 的 ★☆**：已与主题色（如 `--color-gold`、`--color-text-disabled`）绑定，若设计无强制要求，可保留；若需与图标体系完全统一，再改为 uv 的 star 系列图标。
- **PlayerDetailHeader 的 ♂♀**：占位小，可保留字符或换成小尺寸 uv-icon（若 uv 提供性别图标）。
- 替换前请在工程内确认 uv-ui 实际提供的 icon 名称（如 `checkmark-circle-fill`、`clock`、`lock`、`chat` 等），避免写错导致不显示。

以上为 app 端 emoji/图标使用检查结果与统一建议。
