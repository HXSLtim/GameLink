# GameLink Taro 小程序实施指南

> **版本**: v1.0.0
> **更新时间**: 2025-01-10
> **目标受众**: 开发者

---

## 📋 目录

1. [快速开始](#快速开始)
2. [项目结构](#项目结构)
3. [开发环境搭建](#开发环境搭建)
4. [组件使用指南](#组件使用指南)
5. [页面开发规范](#页面开发规范)
6. [状态管理](#状态管理)
7. [API 集成](#api-集成)
8. [最佳实践](#最佳实践)
9. [常见问题](#常见问题)
10. [部署流程](#部署流程)

---

## 快速开始

### 系统要求

- **Node.js**: >= 16.x
- **npm**: >= 7.x
- **Taro CLI**: >= 4.x
- **微信开发者工具**: 最新版

### 安装依赖

```bash
cd app
npm install
```

### 启动开发服务器

```bash
# 微信小程序
npm run dev:weapp

# H5
npm run dev:h5

# React Native
npm run dev:rn
```

### 项目初始化检查清单

- [ ] 安装 Node.js 和 npm
- [ ] 安装 Taro CLI: `npm install -g @tarojs/cli`
- [ ] 安装微信开发者工具
- [ ] 配置微信开发者工具（AppID、项目目录）
- [ ] 运行 `npm install` 安装项目依赖
- [ ] 运行 `npm run dev:weapp` 启动开发服务
- [ ] 在微信开发者工具中预览项目

---

## 项目结构

### 目录结构说明

```
app/
├── config/                    # Taro 配置文件
│   ├── index.ts              # 主配置
│   ├── dev.ts                # 开发环境配置
│   └── prod.ts               # 生产环境配置
├── src/
│   ├── app.config.ts         # 应用配置（页面路由、TabBar）
│   ├── app.scss              # 全局样式入口
│   ├── app.ts                # 应用入口
│   ├── components/           # 组件库
│   │   ├── basic/            # 基础组件
│   │   │   ├── Button/       # 按钮
│   │   │   ├── Input/        # 输入框
│   │   │   ├── Tag/          # 标签
│   │   │   └── Card/         # 卡片
│   │   ├── business/         # 业务组件
│   │   │   ├── PlayerCardNative/      # 陪玩师卡片
│   │   │   ├── OrderCardNative/       # 订单卡片
│   │   │   ├── PriceDisplay/          # 价格显示
│   │   │   ├── OnlineStatus/          # 在线状态
│   │   │   └── GameTag/              # 游戏标签
│   │   └── layout/           # 布局组件
│   │       ├── NavBar/       # 导航栏
│   │       ├── TabBar/       # 底部标签栏
│   │       ├── Empty/        # 空状态
│   │       └── Loading/      # 加载状态
│   ├── pages/                # 页面
│   │   ├── common/           # 通用页面
│   │   │   ├── login/        # 登录
│   │   │   ├── search/       # 搜索
│   │   │   └── settings/     # 设置
│   │   ├── user/             # 用户端页面
│   │   │   ├── home/         # 首页
│   │   │   ├── discover/     # 发现
│   │   │   ├── player-detail/# 陪玩师详情
│   │   │   ├── order-create/ # 创建订单
│   │   │   ├── orders/       # 订单列表
│   │   │   ├── wallet/       # 钱包
│   │   │   ├── vip/          # VIP中心
│   │   │   └── profile/      # 个人中心
│   │   └── player/           # 陪玩师端页面
│   │       ├── hall/         # 接单大厅
│   │       ├── orders/       # 订单管理
│   │       ├── income/       # 收入
│   │       └── profile/      # 个人中心
│   ├── stores/              # 状态管理 (Zustand)
│   ├── services/            # API 服务
│   ├── utils/               # 工具函数
│   ├── constants/           # 常量定义
│   ├── types/               # TypeScript 类型
│   └── styles/              # 全局样式
│       ├── taro-native.scss # 主设计系统文件
│       ├── variables.scss   # 设计变量
│       ├── mixins.scss      # 样式混入
│       └── global.scss      # 全局样式
├── package.json
├── tsconfig.json
└── project.config.json      # 微信小程序配置
```

### 文件命名规范

```
组件文件：PascalCase
├── Button/
│   ├── index.tsx           # 组件实现
│   └── index.scss          # 组件样式
├── PlayerCardNative/
│   ├── index.tsx
│   └── index.scss

页面文件：小写 + 连字符
├── user-home/
│   ├── index.tsx
│   ├── index.scss
│   └── index.config.ts     # 页面配置

样式文件：kebab-case 或与组件同名
├── index.scss
├── player-card.scss
```

---

## 开发环境搭建

### 1. 安装开发工具

```bash
# 安装 Taro CLI
npm install -g @tarojs/cli

# 验证安装
taro --version
```

### 2. 配置微信开发者工具

1. 打开微信开发者工具
2. 点击 "导入项目"
3. 选择 `app/dist` 目录
4. 填写 AppID（测试号或正式 AppID）
5. 点击 "导入"

### 3. 配置环境变量

创建 `app/.env.local`：

```bash
# API 基础地址
API_BASE_URL=http://localhost:8080/api/v1

# 微信 AppID
WECHAT_APP_ID=your_app_id
```

### 4. 配置 Taro

编辑 `app/config/index.ts`：

```typescript
export default defineConfig({
  projectName: 'gamelink-miniprogram',
  date: '2025-1-10',
  designWidth: 750,  // 设计稿宽度
  deviceRatio: {
    640: 2.34 / 2,
    750: 1,
    828: 1.81 / 2
  },
  sourceRoot: 'src',
  outputRoot: 'dist',
  framework: 'react',
  compiler: {
    type: 'webpack5',
  },
  // ...
})
```

---

## 组件使用指南

### 基础组件使用

#### 1. Button 按钮

```tsx
import { Button } from '@/components/basic/Button'

// 主要按钮
<Button type="primary" onClick={handleClick}>
  立即下单
</Button>

// 次要按钮
<Button type="secondary">取消</Button>

// 加载状态
<Button loading>提交中...</Button>

// 禁用状态
<Button disabled>不可用</Button>

// 不同尺寸
<Button size="sm">小按钮</Button>
<Button size="md">中按钮</Button>
<Button size="lg">大按钮</Button>
```

#### 2. Input 输入框

```tsx
import { Input } from '@/components/basic/Input'

// 基础输入
<Input
  placeholder="请输入昵称"
  onInput={(value) => console.log(value)}
/>

// 带标题
<Input
  title="昵称"
  required
  placeholder="请输入昵称"
/>

// 文本域
<Input type="textarea" placeholder="请输入简介" />

// 错误提示
<Input error="昵称不能为空" />
```

#### 3. Tag 标签

```tsx
import { Tag } from '@/components/basic/Tag'

// 基础标签
<Tag>默认</Tag>

// 类型标签
<Tag type="primary">VIP</Tag>
<Tag type="success">在线</Tag>
<Tag type="warning">忙碌</Tag>

// 可关闭
<Tag closable onClose={() => console.log('close')}>
  可关闭标签
</Tag>
```

### 业务组件使用

#### 1. PlayerCardNative 陪玩师卡片

```tsx
import { PlayerCardNative } from '@/components/business/PlayerCardNative'

<PlayerCardNative
  player={{
    id: 1,
    nickname: '技术流小哥哥',
    avatar: 'https://...',
    rating: 4.9,
    reviewCount: 120,
    onlineStatus: 'online',
    price: 50,
    tags: ['声音好', '技术流'],
    gameName: '王者荣耀',
  }}
  onClick={() => navigateToDetail(1)}
/>
```

#### 2. OrderCardNative 订单卡片

```tsx
import { OrderCardNative } from '@/components/business/OrderCardNative'

<OrderCardNative
  order={{
    id: 1,
    orderNo: 'ORD20250110001',
    gameName: '王者荣耀',
    playerName: '技术流小哥哥',
    playerAvatar: 'https://...',
    status: 'in_progress',
    priceCents: 5000,
    startTime: '2025-01-10 20:00',
    duration: 120,
  }}
  onClick={() => navigateToOrderDetail(1)}
/>
```

#### 3. PriceDisplay 价格显示

```tsx
import { PriceDisplay } from '@/components/business/PriceDisplay'

// 基础价格
<PriceDisplay priceCents={5000} />

// VIP 折扣
<PriceDisplay
  priceCents={5000}
  vipLevel={3}
  vipDiscount={0.9}
/>

// 显示划线价
<PriceDisplay
  priceCents={4500}
  originalPriceCents={5000}
/>
```

### 布局组件使用

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
    <Text onClick={handleComplete}>完成</Text>
  }
/>
```

#### 2. TabBar 底部标签栏

```tsx
import { TabBar } from '@/components/layout/TabBar'

<TabBar
  tabs={[
    { key: 'home', icon: 'home', text: '首页', path: '/pages/user/home/index' },
    { key: 'discover', icon: 'discover', text: '发现', path: '/pages/user/discover/index' },
    { key: 'orders', icon: 'orders', text: '订单', path: '/pages/user/orders/index' },
    { key: 'profile', icon: 'profile', text: '我的', path: '/pages/user/profile/index' },
  ]}
  activeKey="home"
  onChange={(key, path) => {
    Taro.switchTab({ url: `/${path}` })
  }}
/>
```

---

## 页面开发规范

### 页面结构模板

```tsx
import React, { useState, useEffect } from 'react'
import { View, ScrollView } from '@tarojs/components'
import Taro from '@tarojs/taro'
import { NavBar } from '@/components/layout/NavBar'
import './index.scss'

export default function MyPage() {
  // 1. 状态定义
  const [loading, setLoading] = useState(false)
  const [data, setData] = useState([])

  // 2. 生命周期
  useEffect(() => {
    loadData()
  }, [])

  // 3. 数据加载
  const loadData = async () => {
    setLoading(true)
    try {
      // TODO: 调用 API
      setData([])
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  // 4. 事件处理
  const handleClick = () => {
    // TODO: 处理点击事件
  }

  // 5. 渲染
  return (
    <View className="page my-page">
      <NavBar title="页面标题" showBack />

      <ScrollView className="my-page__scroll" scrollY>
        {/* 页面内容 */}
      </ScrollView>
    </View>
  )
}
```

### 页面配置

每个页面可以创建 `index.config.ts` 文件来配置页面：

```typescript
export default definePageConfig({
  navigationBarTitleText: '页面标题',
  navigationBarBackgroundColor: '#ffffff',
  navigationBarTextStyle: 'black',
  enablePullDownRefresh: true,
  onReachBottomDistance: 50,
})
```

### 样式编写规范

使用 SCSS 和设计系统变量：

```scss
@import '@/styles/taro-native.scss';

.my-page {
  min-height: 100vh;
  background-color: $bg-page;

  &__scroll {
    height: calc(100vh - #{$nav-bar-height});
  }

  &__section {
    padding: $spacing-md;
    background-color: $bg-card;
    margin-bottom: $spacing-sm;
  }

  &__title {
    font-size: $font-size-md;
    font-weight: $font-weight-semibold;
    color: $text-primary;
  }
}
```

---

## 状态管理

### Zustand Store 创建

```typescript
// src/stores/user.ts
import { create } from 'zustand'

interface User {
  id: number
  nickname: string
  avatar: string
  vipLevel: number
}

interface UserStore {
  user: User | null
  token: string | null
  setUser: (user: User) => void
  setToken: (token: string) => void
  logout: () => void
}

export const useUserStore = create<UserStore>((set) => ({
  user: null,
  token: null,
  setUser: (user) => set({ user }),
  setToken: (token) => set({ token }),
  logout: () => set({ user: null, token: null }),
}))
```

### Store 使用

```tsx
import { useUserStore } from '@/stores/user'

export default function MyPage() {
  const { user, logout } = useUserStore()

  return (
    <View>
      <Text>{user?.nickname}</Text>
      <Button onClick={logout}>退出登录</Button>
    </View>
  )
}
```

---

## API 集成

### API 服务封装

```typescript
// src/services/api.ts
import Taro from '@tarojs/taro'

const BASE_URL = process.env.API_BASE_URL || 'http://localhost:8080/api/v1'

export const request = async <T = any>(
  url: string,
  options: RequestInit = {}
): Promise<T> => {
  const token = Taro.getStorageSync('token')

  const response = await fetch(`${BASE_URL}${url}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      'Authorization': token ? `Bearer ${token}` : '',
      ...options.headers,
    },
  })

  if (!response.ok) {
    throw new Error(`HTTP Error: ${response.status}`)
  }

  return response.json()
}

export const api = {
  get: <T = any>(url: string) => request<T>(url, { method: 'GET' }),
  post: <T = any>(url: string, data: any) =>
    request<T>(url, { method: 'POST', body: JSON.stringify(data) }),
  put: <T = any>(url: string, data: any) =>
    request<T>(url, { method: 'PUT', body: JSON.stringify(data) }),
  delete: <T = any>(url: string) => request<T>(url, { method: 'DELETE' }),
}
```

### API 使用示例

```typescript
// src/services/player.ts
import { api } from './api'

export const playerService = {
  // 获取陪玩师列表
  getList: (params: any) => {
    return api.get('/players', params)
  },

  // 获取陪玩师详情
  getDetail: (id: number) => {
    return api.get(`/players/${id}`)
  },

  // 创建订单
  createOrder: (data: any) => {
    return api.post('/orders', data)
  },
}
```

---

## 最佳实践

### 1. 性能优化

#### 图片懒加载

```tsx
<Image lazyLoad mode="aspectFill" src={imageSrc} />
```

#### 列表虚拟化（长列表）

```tsx
// 使用 Taro 的 virtual list 组件
import { VirtualList } from '@tarojs/components'

<VirtualList
  items={data}
  itemSize={100}
  renderItem={(item) => <PlayerCardNative player={item} />}
/>
```

#### 防抖和节流

```typescript
import { debounce } from '@/utils'

const handleSearch = debounce((value) => {
  // 搜索逻辑
}, 300)
```

### 2. 错误处理

```tsx
const loadData = async () => {
  try {
    setLoading(true)
    const data = await api.get('/data')
    setData(data)
  } catch (error) {
    Taro.showToast({
      title: '加载失败',
      icon: 'error',
    })
  } finally {
    setLoading(false)
  }
}
```

### 3. 类型安全

```typescript
// 定义类型
interface Player {
  id: number
  nickname: string
  avatar: string
  // ...
}

// 使用类型
const [players, setPlayers] = useState<Player[]>([])

// 类型导出
export type { Player }
```

---

## 常见问题

### Q1: 如何处理微信小程序登录？

```typescript
const handleLogin = async () => {
  // 1. 获取微信登录码
  const { code } = await Taro.login()

  // 2. 发送到后端换取 token
  const { token } = await api.post('/auth/wechat-login', { code })

  // 3. 保存 token
  Taro.setStorageSync('token', token)

  // 4. 跳转到首页
  Taro.switchTab({ url: '/pages/user/home/index' })
}
```

### Q2: 如何实现下拉刷新和上拉加载？

```tsx
<ScrollView
  scrollY
  refresherEnabled
  refresherTriggered={refreshing}
  onRefresherRefresh={handleRefresh}
  onScrollToLower={handleLoadMore}
>
  {/* 列表内容 */}
</ScrollView>
```

### Q3: 如何处理图片上传？

```typescript
const handleUpload = async () => {
  const { tempFilePaths } = await Taro.chooseImage({
    count: 1,
  })

  const { url } = await Taro.uploadFile({
    url: `${BASE_URL}/upload`,
    filePath: tempFilePaths[0],
    name: 'file',
    header: {
      'Authorization': `Bearer ${token}`,
    },
  })

  console.log('上传成功:', url)
}
```

### Q4: 如何实现支付？

```typescript
const handlePay = async (orderId: number) => {
  // 1. 获取支付参数
  const payParams = await api.post(`/orders/${orderId}/payment`)

  // 2. 调起微信支付
  await Taro.requestPayment(payParams)

  // 3. 支付成功后跳转
  Taro.redirectTo({
    url: `/pages/user/order-detail/index?id=${orderId}`,
  })
}
```

---

## 部署流程

### 1. 构建生产版本

```bash
# 构建微信小程序
npm run build:weapp

# 构建会生成 dist 目录
```

### 2. 上传到微信平台

1. 打开微信开发者工具
2. 点击 "上传" 按钮
3. 填写版本号和项目备注
4. 点击 "上传"

### 3. 提交审核

1. 登录微信公众平台
2. 进入 "版本管理"
3. 点击 "提交审核"
4. 填写审核信息
5. 等待审核结果

### 4. 发布上线

审核通过后，点击 "发布" 按钮即可上线。

---

## 附录

### 相关文档

- [Taro 官方文档](https://taro-docs.jd.com/)
- [React 官方文档](https://react.dev/)
- [TypeScript 官方文档](https://www.typescriptlang.org/)
- [微信小程序文档](https://developers.weixin.qq.com/miniprogram/dev/framework/)
- [设计系统文档](./TARO_UI_DESIGN_SYSTEM.md)
- [原型设计文档](./MINIPROTOTYPE_DESIGN.md)

### 工具推荐

- **微信开发者工具**: 小程序开发和调试
- **VSCode**: 代码编辑器
- **ESLint**: 代码检查
- **Prettier**: 代码格式化
- **Sass**: CSS 预处理器

### 团队协作

- 使用 Git 进行版本控制
- 遵循 Conventional Commits 规范
- 代码审查（Code Review）
- 定期同步设计系统更新

---

**维护者**: GameLink 开发团队
**反馈邮箱**: dev@gamelink.com
**最后更新**: 2025-01-10
