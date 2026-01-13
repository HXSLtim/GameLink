# GameLink Taro 原生组件设计系统

> **版本**: v1.0.0
> **更新时间**: 2025-01-10
> **状态**: 生产就绪

---

## 📋 目录

1. [概述](#概述)
2. [技术栈](#技术栈)
3. [设计原则](#设计原则)
4. [色彩系统](#色彩系统)
5. [字体系统](#字体系统)
6. [间距系统](#间距系统)
7. [组件库](#组件库)
8. [页面布局](#页面布局)
9. [导航系统](#导航系统)
10. [开发规范](#开发规范)

---

## 概述

### 设计系统介绍

GameLink Taro 设计系统是一套**完全基于 Taro 原生组件**的 UI/UX 解决方案，不依赖任何第三方 UI 库（如 NutUI），确保与微信小程序、H5、React Native 等平台的完全兼容。

### 核心特性

- ✅ **纯 Taro 原生组件** - 100% 兼容微信小程序
- ✅ **零依赖** - 不依赖任何 UI 库
- ✅ **完整的设计令牌** - 统一的设计变量系统
- ✅ **响应式设计** - 支持多种屏幕尺寸
- ✅ **无障碍支持** - 遵循 WCAG 2.1 标准
- ✅ **性能优化** - 轻量级，快速渲染

### 文件结构

```
app/src/
├── styles/
│   ├── taro-native.scss          # 主设计系统文件
│   ├── variables.scss            # 设计变量
│   ├── mixins.scss               # 样式混入
│   └── global.scss               # 全局样式
├── components/
│   ├── basic/                    # 基础组件
│   │   ├── Button/
│   │   ├── Input/
│   │   ├── Tag/
│   │   └── ...
│   ├── business/                 # 业务组件
│   │   ├── PlayerCard/
│   │   ├── OrderCard/
│   │   ├── PriceDisplay/
│   │   └── ...
│   └── layout/                   # 布局组件
│       ├── NavBar/
│       ├── TabBar/
│       └── ...
├── pages/                        # 页面
├── stores/                       # 状态管理
└── utils/                        # 工具函数
```

---

## 技术栈

### 核心技术

| 技术 | 版本 | 说明 |
|------|------|------|
| **Taro** | 3.x / 4.x | 多端开发框架 |
| **React** | 18.x | UI 框架 |
| **TypeScript** | 5.x | 类型系统 |
| **Sass** | 1.x | CSS 预处理器 |
| **Zustand** | 5.x | 状态管理 |

### Taro 原生组件映射

```
Taro 组件         →  微信小程序     →  H5          →  React Native
View             →  view          →  div         →  View
Text             →  text          →  span        →  Text
Image            →  image         →  img         →  Image
ScrollView       →  scroll-view   →  div         →  ScrollView
Swiper           →  swiper        →  div         →  ScrollView
Button           →  button        →  button      →  TouchableOpacity
Input            →  input         →  input       →  TextInput
TextArea         →  textarea      →  textarea    →  TextInput
Switch           →  switch        →  input       →  Switch
Slider           →  slider        →  input       →  Slider
Picker           →  picker        →  select      →  Picker
PickerView       →  picker-view   →  div         →  Picker
Navigator        →  navigator     →  a           →  View
```

---

## 设计原则

### 1. 简洁优先

- **最少步骤**：核心任务最多 3 步完成
- **减少干扰**：突出主要功能，隐藏次要信息
- **清晰层级**：使用字体大小、颜色、间距建立视觉层级

### 2. 一致性

- **视觉一致**：统一的设计语言和组件规范
- **交互一致**：相同的操作产生相同的反馈
- **内容一致**：统一的文案风格和术语

### 3. 性能优先

- **轻量级**：减少依赖，优化代码体积
- **快速渲染**：使用 Taro 原生组件，避免过度封装
- **流畅动画**：60fps 的动画体验

### 4. 可访问性

- **对比度**：文字与背景对比度 ≥ 4.5:1
- **触控目标**：最小点击区域 44x44pt
- **语义化**：使用正确的 HTML 语义

---

## 色彩系统

### 品牌色彩

```scss
// 主色调 - 紫色渐变
$primary-color: #667eea;        // 品牌主色
$primary-dark: #764ba2;         // 品牌辅助色
$primary-light: #8b9eef;        // 浅紫色

// 强调色 - 橙色（用于价格、CTA）
$accent-color: #f59e0b;
```

**应用场景**：
- 主按钮、导航栏选中状态
- 价格、重要数据高亮
- 链接、图标强调

### 功能色彩

```scss
$success-color: #10b981;        // 成功 - 绿色
$warning-color: #f59e0b;        // 警告 - 黄色
$error-color: #ef4444;          // 错误 - 红色
$info-color: #3b82f6;           // 信息 - 蓝色
```

**应用场景**：
- **成功**：订单完成、操作成功
- **警告**：余额不足、重要提示
- **错误**：网络错误、操作失败
- **信息**：系统通知、帮助提示

### 中性色彩

```scss
// 文字颜色
$text-primary: #1f2937;         // 主文字 - 深灰
$text-secondary: #6b7280;       // 副标题 - 中灰
$text-tertiary: #9ca3af;        // 提示文字 - 浅灰
$text-disabled: #d1d5db;        // 禁用文字

// 背景颜色
$bg-page: #f9fafb;              // 页面背景
$bg-card: #ffffff;              // 卡片背景
$bg-input: #f3f4f6;             // 输入框背景
$bg-active: #f3f4f6;            // 点击态背景

// 边框颜色
$border-color: #e5e7eb;         // 边框 - 极浅灰
$divider-color: #f3f4f6;        // 分割线
```

### 色彩使用指南

#### 文字颜色选择

| 场景 | 颜色变量 | 示例 |
|------|----------|------|
| 标题、正文 | `$text-primary` | 页面标题、卡片正文 |
| 次要信息 | `$text-secondary` | 时间戳、辅助说明 |
| 提示文字 | `$text-tertiary` | placeholder、禁用状态 |
| 品牌强调 | `$primary-color` | 链接、标签、图标 |

#### 背景颜色选择

| 场景 | 颜色变量 | 示例 |
|------|----------|------|
| 页面背景 | `$bg-page` | 整个页面背景 |
| 卡片背景 | `$bg-card` | 白色卡片、列表项 |
| 输入框 | `$bg-input` | 表单输入框背景 |

#### 功能色使用规范

```
成功状态：
✅ 订单完成、支付成功
✅ 操作成功提示
✅ 在线状态指示

警告状态：
⚠️ 余额不足、时间不足
⚠️ 重要操作确认
⚠️ 到期提醒

错误状态：
❌ 网络错误、操作失败
❌ 表单验证错误
❌ 权限不足
```

---

## 字体系统

### 字体家族

```scss
font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto',
             'Oxygen', 'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans',
             'Helvetica Neue', sans-serif;
```

**说明**：
- 优先使用系统字体，提升加载速度
- 自动适配 iOS、Android、Web 不同平台

### 字体大小

```scss
$font-size-xs: 20px;     // 10pt - 极小文字（标签、徽标）
$font-size-sm: 24px;     // 12pt - 小文字（辅助信息）
$font-size-base: 28px;   // 14pt - 基础文字（正文）
$font-size-md: 32px;     // 16pt - 中等文字（副标题）
$font-size-lg: 36px;     // 18pt - 大文字（标题）
$font-size-xl: 40px;     // 20pt - 大标题（页面标题）
$font-size-xxl: 48px;    // 24pt - 超大标题（特殊展示）
```

**使用场景**：

| 字体大小 | 应用场景 | 示例 |
|----------|----------|------|
| 20px (xs) | 标签、徽标、注释 | 标签文字、按钮小字 |
| 24px (sm) | 辅助信息 | 时间戳、价格单位 |
| 28px (base) | 正文内容 | 卡片内容、列表文字 |
| 32px (md) | 副标题 | 卡片标题、列表标题 |
| 36px (lg) | 标题 | 区块标题、弹窗标题 |
| 40px (xl) | 页面标题 | 页面主标题 |
| 48px (xxl) | 特殊展示 | 数字大屏、活动标题 |

### 字重

```scss
$font-weight-normal: 400;    // 常规 - 正文
$font-weight-medium: 500;    // 中等 - 强调文字
$font-weight-semibold: 600;  // 半粗 - 小标题
$font-weight-bold: 700;      // 粗体 - 大标题
```

**使用规范**：

| 字重 | 应用场景 | 示例 |
|------|----------|------|
| 400 (Regular) | 正文内容 | 列表项、卡片内容 |
| 500 (Medium) | 强调文字 | 价格、重要数据 |
| 600 (Semibold) | 小标题 | 卡片标题、按钮文字 |
| 700 (Bold) | 大标题 | 页面标题、弹窗标题 |

### 行高

```scss
$line-height-tight: 1.2;     // 标题行高
$line-height-normal: 1.5;    // 正文行高
$line-height-relaxed: 1.8;   // 宽松行高
```

**使用规范**：

| 行高 | 应用场景 | 示例 |
|------|----------|------|
| 1.2 | 标题文字 | 页面标题、卡片标题 |
| 1.5 | 正文内容 | 列表、卡片、表单 |
| 1.8 | 长文本 | 文章内容、详情描述 |

### 文本样式工具类

```tsx
// 使用示例
<Text className="u-text-primary u-text-base u-text-normal">
  正文内容
</Text>

<Text className="u-text-brand u-text-lg u-text-bold">
  品牌强调标题
</Text>

<Text className="u-text-secondary u-text-sm u-text-ellipsis">
  单行省略文字
</Text>
```

---

## 间距系统

### 间距刻度

基于 **4px 网格系统**的间距规范：

```scss
$spacing-xs: 8px;      // 极小间距
$spacing-sm: 12px;     // 小间距
$spacing-md: 16px;     // 中等间距
$spacing-lg: 24px;     // 大间距
$spacing-xl: 32px;     // 超大间距
$spacing-xxl: 48px;    // 极大间距
```

### 使用规范

| 间距 | 应用场景 | 示例 |
|------|----------|------|
| 8px (xs) | 紧密元素间距 | 标签与文字、图标与文字 |
| 12px (sm) | 小元素间距 | 列表项内边距 |
| 16px (md) | 标准间距 | 卡片内边距、表单元素间距 |
| 24px (lg) | 大间距 | 区块间距、页面边距 |
| 32px (xl) | 超大间距 | 页面区块间距 |
| 48px (xxl) | 页面级间距 | 大区块间距 |

### 间距工具类

```tsx
// Margin
<View className="u-mt-md">上边距 16px</View>
<View className="u-mb-lg">下边距 24px</View>
<View className="u-ml-sm">左边距 12px</View>
<View className="u-mr-sm">右边距 12px</View>

// Padding
<View className="u-p-md">内边距 16px</View>
<View className="u-p-lg">内边距 24px</View>
```

### 圆角系统

```scss
$radius-xs: 4px;       // 极小圆角 - 输入框、按钮
$radius-sm: 8px;       // 小圆角 - 卡片、标签
$radius-md: 12px;      // 中等圆角 - 弹窗
$radius-lg: 16px;      // 大圆角 - 大卡片
$radius-xl: 24px;      // 超大圆角 - 图片
$radius-full: 9999px;  // 完全圆角 - 头像、徽标
```

---

## 组件库

### 基础组件 (Basic Components)

#### 1. Button 按钮

```tsx
import { Button } from '@/components/basic/Button'

// 主要按钮
<Button type="primary" onClick={handleClick}>
  立即下单
</Button>

// 次要按钮
<Button type="secondary">取消</Button>

// 幽灵按钮
<Button type="ghost">查看详情</Button>

// 尺寸
<Button size="sm">小按钮</Button>
<Button size="md">中按钮</Button>
<Button size="lg">大按钮</Button>

// 块级按钮
<Button block>全宽按钮</Button>

// 禁用状态
<Button disabled>禁用按钮</Button>
```

**按钮类型**：

| 类型 | 样式 | 使用场景 |
|------|------|----------|
| `primary` | 紫色渐变背景 | 主要操作、提交表单 |
| `secondary` | 白色背景 + 紫色边框 | 次要操作、取消 |
| `ghost` | 透明背景 + 紫色文字 | 辅助操作、查看详情 |
| `success` | 绿色背景 | 成功操作 |
| `warning` | 黄色背景 | 警告操作 |
| `danger` | 红色背景 | 危险操作、删除 |

#### 2. Input 输入框

```tsx
import { Input } from '@/components/basic/Input'

// 基础输入框
<Input placeholder="请输入陪玩师昵称" />

// 带标题
<Input title="昵称" placeholder="请输入昵称" />

// 文本域
<Input.TextArea
  placeholder="请输入特殊要求"
  maxLength={200}
  showCount
/>

// 错误状态
<Input error="请输入有效的手机号" />

// 禁用状态
<Input disabled value="已禁用" />
```

#### 3. Tag 标签

```tsx
import { Tag } from '@/components/basic/Tag'

// 基础标签
<Tag>默认标签</Tag>

// 类型标签
<Tag type="primary">VIP</Tag>
<Tag type="success">在线</Tag>
<Tag type="warning">忙碌</Tag>
<Tag type="danger">离线</Tag>

// 圆角标签
<Tag circle>圆形标签</Tag>

// 可关闭标签
<Tag closable onClose={handleClose}>
  可关闭标签
</Tag>
```

#### 4. Card 卡片

```tsx
import { Card } from '@/components/basic/Card'

// 基础卡片
<Card title="卡片标题">
  <View>卡片内容</View>
</Card>

// 带副标题
<Card
  title="卡片标题"
  subtitle="卡片副标题"
>
  <View>卡片内容</View>
</Card>

// 带底部操作
<Card
  title="卡片标题"
  footer={
    <View className="u-flex-between">
      <Button size="sm">取消</Button>
      <Button size="sm" type="primary">确定</Button>
    </View>
  }
>
  <View>卡片内容</View>
</Card>
```

### 业务组件 (Business Components)

#### 1. PlayerCard 陪玩师卡片

```tsx
import { PlayerCard } from '@/components/business/PlayerCard'

<PlayerCard
  player={{
    id: 1,
    nickname: '技术流小哥哥',
    avatar: 'https://...',
    rating: 4.9,
    onlineStatus: 'online',
    price: 50,
    tags: ['声音好', '技术流']
  }}
  onClick={() => navigateToPlayerDetail(1)}
/>
```

#### 2. OrderCard 订单卡片

```tsx
import { OrderCard } from '@/components/business/OrderCard'

<OrderCard
  order={{
    id: 1,
    gameName: '王者荣耀',
    playerName: '技术流小哥哥',
    playerAvatar: 'https://...',
    status: 'in_progress',
    price: 100,
    startTime: '2025-01-10 20:00'
  }}
  onClick={() => navigateToOrderDetail(1)}
/>
```

#### 3. PriceDisplay 价格显示

```tsx
import { PriceDisplay } from '@/components/business/PriceDisplay'

// 基础价格
<PriceDisplay price={100} />

// 带 VIP 折扣
<PriceDisplay
  price={100}
  vipLevel={3}
  discount={0.9}
/>

// 显示划线价
<PriceDisplay
  price={80}
  originalPrice={100}
/>
```

#### 4. OnlineStatus 在线状态

```tsx
import { OnlineStatus } from '@/components/business/OnlineStatus'

<OnlineStatus status="online" />  // 🟢 在线
<OnlineStatus status="busy" />    // 🟡 忙碌
<OnlineStatus status="offline" /> // 🔴 离线
```

### 布局组件 (Layout Components)

#### 1. NavBar 导航栏

```tsx
import { NavBar } from '@/components/layout/NavBar'

// 基础导航栏
<NavBar title="页面标题" />

// 带返回按钮
<NavBar
  title="页面标题"
  showBack
  onBack={() => navigateBack()}
/>

// 自定义右侧操作
<NavBar
  title="页面标题"
  showBack
  renderRight={
    <View className="u-text-brand">完成</View>
  }
  onRightClick={handleComplete}
/>
```

#### 2. TabBar 底部标签栏

```tsx
import { TabBar } from '@/components/layout/TabBar'

// 用户端配置
<TabBar
  tabs={[
    { key: 'home', icon: 'home', text: '首页', path: '/pages/index/index' },
    { key: 'discover', icon: 'discover', text: '发现', path: '/pages/discover/index' },
    { key: 'orders', icon: 'orders', text: '订单', path: '/pages/orders/index' },
    { key: 'profile', icon: 'profile', text: '我的', path: '/pages/profile/index' }
  ]}
  activeKey="home"
/>

// 陪玩师端配置
<TabBar
  tabs={[
    { key: 'hall', icon: 'globe', text: '大厅', path: '/pages/hall/index' },
    { key: 'orders', icon: 'orders', text: '订单', path: '/pages/player-orders/index' },
    { key: 'income', icon: 'wallet', text: '收入', path: '/pages/income/index' },
    { key: 'profile', icon: 'profile', text: '我的', path: '/pages/profile/index' }
  ]}
  activeKey="hall"
/>
```

#### 3. Empty 空状态

```tsx
import { Empty } from '@/components/layout/Empty'

// 基础空状态
<Empty
  image="empty"
  text="暂无数据"
/>

// 带操作
<Empty
  image="search"
  text="未找到相关陪玩师"
  buttonText="去看看其他"
  onButtonClick={() => navigateToHome()}
/>
```

---

## 页面布局

### 标准页面结构

```tsx
import { View } from '@tarojs/components'
import { NavBar } from '@/components/layout/NavBar'
import './index.scss'

export default function Page() {
  return (
    <View className="page">
      {/* 导航栏 */}
      <NavBar title="页面标题" />

      {/* 主内容区 */}
      <View className="page__content">
        {/* 页面内容 */}
      </View>

      {/* 底部操作栏（可选） */}
      <View className="page__footer">
        <Button type="primary" block>
          提交
        </Button>
      </View>
    </View>
  )
}
```

### 常用布局模式

#### 1. 列表布局

```tsx
<View className="list">
  {items.map(item => (
    <View key={item.id} className="list__item list__item--clickable">
      <View className="list__content">
        <View className="list__title">{item.title}</View>
        <View className="list__subtitle">{item.subtitle}</View>
      </View>
      <View className="list__extra">></View>
    </View>
  ))}
</View>
```

#### 2. 卡片网格布局

```tsx
<View className="grid grid--2-col">
  {items.map(item => (
    <View key={item.id} className="grid__item">
      <Card>
        {/* 卡片内容 */}
      </Card>
    </View>
  ))}
</View>
```

```scss
// 样式
.grid {
  display: grid;
  gap: $spacing-md;
  padding: $spacing-md;

  &--2-col {
    grid-template-columns: repeat(2, 1fr);
  }

  &--3-col {
    grid-template-columns: repeat(3, 1fr);
  }
}
```

#### 3. 表单布局

```tsx
<View className="form">
  <View className="form__item">
    <View className="form__label">昵称</View>
    <Input placeholder="请输入昵称" />
  </View>

  <View className="form__item">
    <View className="form__label">简介</View>
    <Input.TextArea placeholder="请输入简介" />
  </View>

  <View className="form__actions">
    <Button type="primary" block>
      提交
    </Button>
  </View>
</View>
```

---

## 导航系统

### 路由配置

```typescript
// app.config.ts
export default defineAppConfig({
  pages: [
    // 用户端页面
    'pages/user/index/index',                // 首页
    'pages/user/discover/index',             // 发现
    'pages/user/player-detail/index',        // 陪玩师详情
    'pages/user/order-create/index',         // 创建订单
    'pages/user/orders/index',               // 订单列表
    'pages/user/order-detail/index',         // 订单详情
    'pages/user/chat/index',                 // 聊天
    'pages/user/wallet/index',               // 钱包
    'pages/user/vip/index',                  // VIP
    'pages/user/coupons/index',              // 优惠券
    'pages/user/recharge/index',             // 充值
    'pages/user/profile/index',              // 个人中心

    // 陪玩师端页面
    'pages/player/hall/index',               // 接单大厅
    'pages/player/orders/index',             // 订单管理
    'pages/player/income/index',             // 收入
    'pages/player/withdraw/index',           // 提现
    'pages/player/services/index',           // 我的服务
    'pages/player/stats/index',              // 数据统计

    // 通用页面
    'pages/common/login/index',              // 登录
    'pages/common/search/index',             // 搜索
    'pages/common/settings/index',           // 设置
  ],

  // 用户端 TabBar
  tabBar: {
    custom: true,
    list: [
      { pagePath: 'pages/user/index/index', text: '首页' },
      { pagePath: 'pages/user/discover/index', text: '发现' },
      { pagePath: 'pages/user/orders/index', text: '订单' },
      { pagePath: 'pages/user/profile/index', text: '我的' }
    ]
  }
})
```

### 页面跳转方法

```typescript
import Taro from '@tarojs/taro'

// 跳转到 tabBar 页面
Taro.switchTab({
  url: '/pages/user/index/index'
})

// 保留当前页面，跳转到新页面
Taro.navigateTo({
  url: '/pages/user/player-detail/index?id=123'
})

// 关闭当前页面，跳转到新页面
Taro.redirectTo({
  url: '/pages/common/login/index'
})

// 返回上一页
Taro.navigateBack({
  delta: 1
})
```

---

## 开发规范

### 组件命名规范

```
组件文件：PascalCase
├── PlayerCard.tsx       ✅
├── OrderCard.tsx        ✅
├── priceDisplay.tsx     ❌ 应为 PriceDisplay.tsx

组件导出：default export
export default function PlayerCard() { }

样式文件：kebab-case
├── PlayerCard.scss      ✅
├── order-card.scss      ✅
```

### 样式编写规范

```scss
// BEM 命名规范
.block {}              // 块
.block__element {}     // 元素
.block--modifier {}    // 修饰符

// 示例
.player-card {
  &__avatar {}
  &__info {}
  &__nickname {}
  &__rating {}
  &--online {}
  &--busy {}
}

// 嵌套不超过 3 层
.page {
  .section {
    .title {}          // ✅ 可接受
    .subtitle {
      .text {}         // ❌ 过深嵌套
    }
  }
}
```

### TypeScript 类型定义

```typescript
// types/player.ts
export interface Player {
  id: number
  nickname: string
  avatar: string
  rating: number
  onlineStatus: 'online' | 'busy' | 'offline'
  price: number
  tags: string[]
}

export interface PlayerCardProps {
  player: Player
  onClick?: () => void
}

// 使用类型
const PlayerCard: React.FC<PlayerCardProps> = ({ player, onClick }) => {
  // ...
}
```

### 性能优化建议

#### 1. 列表渲染优化

```tsx
// ❌ 不推荐：无 key
{items.map(item => (
  <PlayerCard player={item} />
))}

// ✅ 推荐：使用稳定 key
{items.map(item => (
  <PlayerCard key={item.id} player={item} />
))}
```

#### 2. 图片懒加载

```tsx
<Image
  src={player.avatar}
  lazyLoad          // 启用懒加载
  mode="aspectFill"
/>
```

#### 3. 避免不必要的渲染

```tsx
import { memo } from 'react'

// 使用 memo 包裹组件
export default memo(PlayerCard)
```

### 常见问题

#### Q: 如何处理安全区域（刘海屏）？

```scss
.footer-bar {
  padding-bottom: calc($spacing-md + env(safe-area-inset-bottom));
}
```

#### Q: 如何实现下拉刷新和上拉加载？

```tsx
import { ScrollView } from '@tarojs/components'

<ScrollView
  scrollY
  refresherEnabled
  refresherTriggered={refreshing}
  onRefresherRefresh={onRefresh}
  onScrollToLower={onLoadMore}
>
  {/* 列表内容 */}
</ScrollView>
```

#### Q: 如何实现底部固定按钮？

```tsx
<View className="page">
  <View className="page__content">
    {/* 内容 */}
  </View>

  <View className="footer-bar">
    <Button type="primary" block>
      提交
    </Button>
  </View>
</View>
```

```scss
.footer-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: $spacing-md;
  padding-bottom: calc($spacing-md + env(safe-area-inset-bottom));
  background-color: $bg-card;
  border-top: 1px solid $divider-color;
  z-index: $z-index-fixed;
}

.page__content {
  padding-bottom: 120px;  // 为底部按钮留出空间
}
```

---

## 附录

### 设计资源下载

- **Figma 设计文件**: `docs/design/GameLink-UI-Kit.fig`
- **图标库**: `docs/design/icons/`
- **插画库**: `docs/design/illustrations/`

### 相关文档

- [Taro 官方文档](https://taro-docs.jd.com/)
- [微信小程序设计指南](https://developers.weixin.qq.com/miniprogram/design/)
- [项目原型设计](./MINIPROTOTYPE_DESIGN.md)
- [配置数据汇总](./CONFIGURATION_DATA.md)

### 更新日志

| 版本 | 日期 | 更新内容 |
|------|------|----------|
| v1.0.0 | 2025-01-10 | 初始版本，完整设计系统 |

---

**维护者**: GameLink 开发团队
**反馈邮箱**: design@gamelink.com
**最后更新**: 2025-01-10
