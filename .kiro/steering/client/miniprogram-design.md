# 小程序 UI/UX 设计规范

> GameLink 微信小程序端设计指南 - Discord Dark Theme

## 技术栈

- **框架**: 原生微信小程序
- **语言**: TypeScript
- **样式**: Less
- **渲染**: Skyline + glass-easel

## 设计风格

**Discord Dark Theme + Gaming Accent**
- 深色背景，减少视觉疲劳
- 双身份（用户/陪玩师）通过主题色区分
- 简洁的卡片式布局，清晰的信息层级

---

## 色彩系统

### 背景层级

| 用途 | 变量名 | 色值 | 说明 |
|------|--------|------|------|
| 主背景 | `@bg-primary` | `#313338` | 页面背景 |
| 次级背景 | `@bg-secondary` | `#2B2D31` | 卡片背景 |
| 最深背景 | `@bg-tertiary` | `#1E1F22` | TabBar/输入框 |
| 悬停态 | `@bg-modifier-hover` | `#35373C` | 交互反馈 |
| 激活态 | `@bg-modifier-active` | `#3F4147` | 按下状态 |

### 文字颜色

| 用途 | 变量名 | 色值 |
|------|--------|------|
| 标题/强调 | `@text-header` | `#F2F3F5` |
| 正文 | `@text-normal` | `#dbdee1` |
| 次要/占位符 | `@text-muted` | `#949BA4` |
| 链接 | `@text-link` | `#00AFF4` |

### 主题色

| 模式 | 变量名 | 色值 | 用途 |
|------|--------|------|------|
| 用户模式 | `@user-primary` | `#5865F2` | Discord Blurple |
| 用户悬停 | `@user-primary-hover` | `#4752C4` | 按钮悬停 |
| 陪玩师模式 | `@player-primary` | `#F0B132` | Gold |
| 陪玩师悬停 | `@player-primary-hover` | `#E29D28` | 按钮悬停 |

### 状态色

| 用途 | 变量名 | 色值 |
|------|--------|------|
| 成功/在线 | `@color-success` | `#23A559` |
| 警告 | `@color-warning` | `#F0B132` |
| 错误/离线 | `@color-error` | `#DA373C` |
| 信息 | `@color-info` | `#5865F2` |

---

## 字体规范

```less
@font-size-xs: 20rpx;      // 10px - 辅助说明
@font-size-sm: 24rpx;      // 12px - 次要文字
@font-size-base: 28rpx;    // 14px - 正文
@font-size-md: 32rpx;      // 16px - 小标题
@font-size-lg: 36rpx;      // 18px - 标题
@font-size-xl: 44rpx;      // 22px - 大标题
@font-size-2xl: 48rpx;     // 24px - 页面标题

@font-weight-normal: 400;
@font-weight-medium: 500;
@font-weight-semibold: 600;
@font-weight-bold: 700;
@font-weight-extrabold: 800;
```

## 间距规范

```less
@spacing-xs: 8rpx;         // 4px
@spacing-sm: 16rpx;        // 8px
@spacing-md: 24rpx;        // 12px
@spacing-lg: 32rpx;        // 16px
@spacing-xl: 48rpx;        // 24px
@spacing-2xl: 64rpx;       // 32px

@page-padding: 40rpx;      // 20px - 页面边距
@card-padding: 32rpx;      // 16px - 卡片内边距
@card-gap: 32rpx;          // 16px - 卡片间距
```

## 圆角规范

```less
@radius-xs: 4rpx;          // 2px
@radius-sm: 8rpx;          // 4px - 按钮/输入框
@radius-md: 16rpx;         // 8px - 小卡片
@radius-lg: 24rpx;         // 12px - 大卡片
@radius-xl: 32rpx;         // 16px
@radius-2xl: 48rpx;        // 24px
@radius-full: 9999rpx;     // 圆形
```

---

## 页面结构

```
miniprogram/
├── pages/
│   ├── index/              # 首页（Discover - 陪玩师列表）
│   ├── category/           # 游戏分类（Browse Games）
│   ├── message/            # 消息列表（Direct Messages）
│   ├── profile/            # 个人中心（Profile）
│   ├── order/
│   │   ├── create/         # 创建订单
│   │   ├── detail/         # 订单详情
│   │   └── list/           # 订单列表
│   ├── player/             # 陪玩师详情
│   └── player-center/      # 陪玩师工作台（Player模式）
│       ├── dashboard/      # 工作台首页
│       ├── orders/         # 接单管理
│       └── earnings/       # 收益统计
├── components/             # 组件库
│   ├── gl-page/            # 页面容器（状态栏/TabBar占位）
│   ├── gl-button/          # 按钮
│   ├── gl-search/          # 搜索框
│   ├── gl-loading/         # 加载状态
│   ├── gl-empty/           # 空状态
│   ├── gl-section/         # 区块标题
│   ├── gl-card/            # 卡片
│   ├── gl-avatar/          # 头像
│   ├── gl-tag/             # 标签
│   ├── gl-navbar/          # 导航栏
│   ├── player-card/        # 陪玩师卡片
│   ├── game-card/          # 游戏分类卡片
│   ├── message-item/       # 消息项
│   └── order-card/         # 订单卡片
├── custom-tab-bar/         # 自定义 TabBar（支持 SVG 图标）
├── styles/
│   ├── variables.less      # Discord 风格变量定义
│   ├── theme.less          # 主题类和工具类
│   └── mixins.less         # 通用 mixins
├── assets/
│   └── icons/              # SVG 图标
└── utils/
    ├── request.ts          # 请求封装
    ├── auth.ts             # 认证工具
    ├── storage.ts          # 存储工具（含 USER_MODE 键）
    └── theme.ts            # 主题切换
```

## TabBar 配置

### 用户模式 (4 Tab)
1. **Discover** - 陪玩师列表（首页）
2. **Games** - 游戏分类
3. **Messages** - 消息列表
4. **Profile** - 个人中心

### 陪玩师模式 (5 Tab)
1. **Dashboard** - 工作台首页
2. **Orders** - 接单管理
3. **Messages** - 消息列表
4. **Earnings** - 收益统计
5. **Profile** - 个人中心

---

## 组件规范

### gl-page 页面容器

统一处理状态栏占位、页面背景、TabBar 占位。

```wxml
<gl-page isTabBar="{{true}}">
  <!-- 页面内容 -->
</gl-page>
```

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| showStatusBar | Boolean | true | 是否显示状态栏占位 |
| isTabBar | Boolean | false | 是否是 TabBar 页面 |
| bgColor | String | '' | 自定义背景色 |
| customHeader | Boolean | false | 是否有自定义头部（如 Banner） |

### gl-button 按钮

```wxml
<gl-button type="primary" size="medium" block>Button</gl-button>
```

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| type/variant | String | 'primary' | 按钮类型 |
| size | String | 'medium' | 尺寸: small/medium/large |
| block | Boolean | false | 是否占满宽度 |
| disabled | Boolean | false | 是否禁用 |
| loading | Boolean | false | 是否加载中 |

| 类型 | 背景色 | 文字色 | 用途 |
|------|--------|--------|------|
| primary | `@user-primary` | `#FFFFFF` | 主要操作 |
| secondary | `@btn-secondary-bg` | `@text-normal` | 次要操作 |
| outline | transparent | `@text-normal` | 边框按钮 |
| ghost | transparent | `@text-muted` | 文字按钮 |
| gold | `@player-primary` | `#000000` | 陪玩师专用 |

### gl-search 搜索框

```wxml
<gl-search placeholder="Search..." bind:search="onSearch" />
```

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| placeholder | String | 'Search...' | 占位文字 |
| value | String | '' | 输入值 |
| showIcon | Boolean | true | 是否显示搜索图标 |

| 事件 | 说明 |
|------|------|
| input | 输入时触发 |
| search | 确认搜索时触发 |
| clear | 清空时触发 |

### gl-loading 加载组件

```wxml
<gl-loading text="Loading..." />
```

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| text | String | 'Loading...' | 加载文字 |
| showText | Boolean | true | 是否显示文字 |
| size | String | 'medium' | 尺寸: small/medium/large |
| fullscreen | Boolean | false | 是否全屏居中 |

### gl-empty 空状态组件

```wxml
<gl-empty 
  icon="/assets/icons/message.svg"
  title="No messages"
  description="Start a conversation"
  showAction
  actionText="Get Started"
  bind:action="onAction"
/>
```

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| icon | String | '' | 图标路径 |
| title | String | 'No data' | 标题 |
| description | String | '' | 描述文字 |
| showAction | Boolean | false | 是否显示操作按钮 |
| actionText | String | 'Action' | 按钮文字 |

### gl-section 区块标题

Discord 风格的大写标题。

```wxml
<gl-section title="RECOMMENDED" showMore bind:more="onMore">
  <!-- 内容 -->
</gl-section>
```

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| title | String | '' | 标题文字 |
| showMore | Boolean | false | 是否显示更多按钮 |
| moreText | String | 'See All' | 更多按钮文字 |

### game-card 游戏分类卡片

```wxml
<game-card game="{{item}}" bind:tap="onTap" />
```

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| game | Object | {} | 游戏数据 {id, name, icon, playerCount} |
| showCount | Boolean | true | 是否显示在线人数 |

### message-item 消息项

```wxml
<message-item message="{{item}}" bind:tap="onTap" />
```

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| message | Object | {} | 消息数据 {id, name, avatar, lastMessage, time, unread, online} |

---

### 卡片样式

```less
.card-base() {
  background-color: @bg-secondary;
  border-radius: @radius-lg;
  border: 2rpx solid @border-subtle;
}

.card-hover() {
  background-color: @bg-modifier-hover;
}
```

### 在线状态点

```less
.status-dot(@size: 16rpx, @color: @color-success) {
  width: @size;
  height: @size;
  border-radius: 50%;
  background-color: @color;
  border: 4rpx solid @bg-secondary;
}
```

---

## 交互规范

### 动画时长
- 快速反馈: 150ms (`@duration-fast`)
- 常规过渡: 200ms (`@duration-normal`)
- 复杂动画: 300ms (`@duration-slow`)
- 缓动函数: `@ease-out` (进入) / `@ease-in` (退出)

### 点击反馈
- 卡片: `background-color: @bg-modifier-hover`
- 按钮: `transform: scale(0.95)`

### 页面标题样式
```less
.page-title {
  font-size: @font-size-sm;
  font-weight: @font-weight-bold;
  color: @text-muted;
  letter-spacing: 1rpx;
  text-transform: uppercase;
}
```

---

## 命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 页面目录 | kebab-case | `player-center` |
| 组件目录 | kebab-case | `gl-button` |
| TS 文件 | 与目录同名 | `index.ts` |
| Less 文件 | 与目录同名 | `index.less` |
| 类名 | kebab-case + BEM | `.player-card__name` |
| 变量 | camelCase | `userInfo` |
| 常量 | UPPER_SNAKE | `API_BASE_URL` |
| Less 变量 | kebab-case with @ | `@bg-primary` |

---

## 图标资源

位置: `miniprogram/assets/icons/`

| 图标 | 文件名 | 用途 |
|------|--------|------|
| 首页 | `home.svg` | TabBar |
| 分类 | `category.svg` | TabBar |
| 消息 | `message.svg` | TabBar |
| 用户 | `user.svg` | TabBar/Profile |
| 搜索 | `search.svg` | 搜索框 |
| 麦克风 | `mic.svg` | 语音时长 |
| 铃铛 | `bell.svg` | 通知 |
| 钱包 | `wallet.svg` | 收益 |
| 订单 | `order.svg` | 订单管理 |
| 星星 | `star.svg` / `star-outline.svg` | 评分 |
| 编辑 | `edit.svg` | 编辑操作 |
| 箭头 | `arrow-right.svg` | 导航 |
| 关闭 | `close.svg` | 关闭/清除 |
