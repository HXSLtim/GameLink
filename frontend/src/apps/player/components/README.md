# Player Components

## 📁 目录用途

该目录用于存放陪玩师端（Player）专用的业务组件。

## 📦 组件分类

### 个人中心组件
- `ProfileCard` - 个人资料卡片
- `WalletSummary` - 钱包概览
- `EarningsChart` - 收益图表
- `CommissionInfo` - 佣金信息展示
- `SkillTags` - 技能标签展示
- `AvailabilityStatus` - 接单状态设置
- `VerificationStatus` - 认证状态展示

### 订单管理组件
- `OrderCard` - 订单卡片（供陪玩师查看）
- `OrderAcceptButton` - 接单/拒单按钮
- `OrderProgress` - 订单进度展示
- `OrderChat` - 订单内聊天
- `OrderReviewForm` - 订单评价表单
- `OrderHistory` - 历史订单
- `QuickReject` - 快速拒绝模板

### 收益管理组件
- `EarningsList` - 收益明细列表
- `WithdrawForm` - 提现表单
- `BankCardList` - 银行卡列表
- `TransactionHistory` - 交易记录
- `CommissionCalculator` - 佣金计算器

### 接单设置组件
- `GameSelector` - 游戏选择器
- `PriceInput` - 价格输入组件
- `TimeSlotPicker` - 时间段选择
- `AutoAcceptToggle` - 自动接单开关
- `OrderSoundToggle` - 订单提醒音开关

### 技能展示组件
- `GameStats` - 游戏数据展示（胜率、段位等）
- `ReviewList` - 评价列表（陪玩师收到的）
- `RatingDisplay` - 评分展示（分数、星星）
- `ServiceGallery` - 服务展示相册
- `VideoIntro` - 视频介绍播放器

## 📋 命名规范

### 文件夹命名
```
ComponentName/
├── ComponentName.tsx       // 主组件
├── ComponentName.types.ts  // 类型定义
├── index.ts                // 导出
└── components/             // 子组件（如有）
    └── SubComponent.tsx
```

### 组件命名
- 使用 PascalCase（大驼峰）
- 以组件功能命名，例如：`OrderAcceptButton`, `EarningsChart`
- 避免使用通用名称，如：`Card`, `List`

## 🎯 开发规范

### TypeScript
- 所有组件使用 TypeScript
- 定义 Props 接口
- 使用明确的返回值类型

```typescript
// ✅ 推荐
interface OrderCardProps {
  order: Order;
  onAccept: () => void;
  onReject: () => void;
}

export const OrderCard: React.FC<OrderCardProps> = ({ order, onAccept, onReject }) => {
  return <div>...</div>;
};

// ❌ 避免
export const OrderCard = (props: any) => {
  return <div>...</div>;
};
```

### 组件设计原则
1. **移动端优先**: 陪玩师主要在手机上使用
2. **简洁明了**: 避免复杂交互
3. **快速操作**: 接单、拒单等操作要快
4. **实时更新**: 订单状态实时同步
5. **通知提醒**: 重要操作有明确反馈
6. **性能优化**: 使用 React.memo, useMemo, useCallback

### 样式规范
- 优先使用 Tailwind CSS 原子类
- 自定义组件样式放在 `src/styles/tailwind.css`
- 避免内联样式
- 使用响应式类名（md:, lg:）
- **移动端优先**：默认样式为移动端，使用 `md:` 适配桌面端

```tsx
// ✅ 推荐（移动端优先）
<div className="p-4 md:p-6">
  {/* 移动端 p-4, 桌面端 p-6 */}
</div>

// 按钮要大，方便手机点击
<button className="btn btn-primary min-h-12">
  接单
</button>

// ❌ 避免（按钮太小，移动端点击困难）
<button className="btn btn-primary text-sm">
  接单
</button>
```

### Visual Guidelines

#### 颜色规范
- **主色调**: primary（陪玩师品牌色）
- **接单按钮**: green（绿色，表示通过/接单）
- **拒绝按钮**: red（红色，表示拒绝/危险）
- **收益**: green（绿色，表示收入）

```tsx
// ✅ 推荐
<button className="bg-green-600 text-white hover:bg-green-700">
  接单
</button>
<button className="bg-red-600 text-white hover:bg-red-700">
  拒单
</button>
<span className="text-green-600 font-bold">
  +¥100
</span>
```

#### 字体和间距
- 字体大小：移动端最小 14px，按钮最小 16px
- 按钮高度：最低 44px（符合 iOS 规范）
- 间距：使用 Tailwind 标准间距（4, 6, 8）

## 🔗 相关资源

- [全局组件目录](../../shared/components/) - 通用组件（跨应用）
- [陪玩师页面目录](../pages/) - 页面级组件
- [Tailwind CSS 使用指南](../../../docs/STYLE_AND_STATE_MANAGEMENT.md)
- [移动端设计规范](https://developers.google.com/speed/docs/insights/rules)

## 🌟 特色组件示例

### 接单卡片（OrderAcceptCard）

移动端专用的订单卡片，包含接单/拒单按钮。

```typescript
// OrderAcceptCard.tsx
import React from 'react';
import type { Order } from '@/api';

interface OrderAcceptCardProps {
  order: Order;
  onAccept: (orderId: number) => void;
  onReject: (orderId: number) => void;
  countdown?: number; // 接单倒计时（秒）
}

export const OrderAcceptCard: React.FC<OrderAcceptCardProps> = ({
  order,
  onAccept,
  onReject,
  countdown,
}) => {
  return (
    <div className="card p-4 mx-2 mb-4">
      {/* 顶部：游戏和模式 */}
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center">
          <div className="w-10 h-10 bg-primary-100 rounded-lg flex items-center justify-center mr-3">
            <span className="text-primary-800 font-bold">{order.game?.[0]}</span>
          </div>
          <div>
            <h3 className="font-bold text-gray-900">{order.game?.name}</h3>
            <p className="text-sm text-gray-600">{order.serviceType}</p>
          </div>
        </div>
        <div className="text-right">
          <p className="text-lg font-bold text-green-600">¥{order.amount}</p>
          {countdown !== undefined && (
            <p className="text-sm text-red-600">{countdown}秒后过期</p>
          )}
        </div>
      </div>

      {/* 中部：客户需求 */}
      <div className="mb-4">
        <div className="flex items-center mb-2">
          <div className="w-8 h-8 bg-gray-200 rounded-full mr-2" />
          <div>
            <p className="font-medium text-sm">{order.user?.username}</p>
            <p className="text-xs text-gray-500">Lv.{order.user?.level}</p>
          </div>
        </div>
        <p className="text-sm text-gray-700 bg-gray-50 p-2 rounded">
          {order.notes || '无特殊要求'}
        </p>
      </div>

      {/* 底部：操作按钮 */}
      <div className="grid grid-cols-2 gap-3">
        <button
          onClick={() => onReject(order.id)}
          className="btn bg-gray-200 text-gray-800 hover:bg-gray-300 py-3"
        >
          拒单
        </button>
        <button
          onClick={() => onAccept(order.id)}
          className="btn bg-green-500 text-white hover:bg-green-600 py-3 font-bold"
        >
          接单
        </button>
      </div>
    </div>
  );
};
```

### 收益概览（EarningsSummary）

展示陪玩师的收益数据，包括今日、本周、本月。

```typescript
// EarningsSummary.tsx
import React from 'react';

interface EarningsSummaryProps {
  today: number;
  thisWeek: number;
  thisMonth: number;
}

export const EarningsSummary: React.FC<EarningsSummaryProps> = ({
  today,
  thisWeek,
  thisMonth,
}) => {
  return (
    <div className="card p-4 mb-4">
      <h2 className="text-lg font-bold mb-3">收益概览</h2>

      <div className="grid grid-cols-3 gap-3">
        <div className="text-center">
          <p className="text-xs text-gray-500 mb-1">今日</p>
          <p className="text-lg font-bold text-green-600">¥{today}</p>
        </div>
        <div className="text-center">
          <p className="text-xs text-gray-500 mb-1">本周</p>
          <p className="text-lg font-bold text-green-600">¥{thisWeek}</p>
        </div>
        <div className="text-center">
          <p className="text-xs text-gray-500 mb-1">本月</p>
          <p className="text-lg font-bold text-green-600">¥{thisMonth}</p>
        </div>
      </div>

      <button className="btn btn-primary w-full mt-4">
        查看明细
      </button>
    </div>
  );
};
```

## 📱 移动端适配策略

### 1. 断点设置

```css
/* tailwind.config.ts */
module.exports = {
  theme: {
    screens: {
      'sm': '640px',   // 平板横屏
      'md': '768px',   // 平板
      'lg': '1024px',  // 小型笔记本
      'xl': '1280px',  // 桌面
    },
  },
}
```

### 2. 常见移动端问题

#### 字体太小
```tsx
// ✅ 推荐 - 移动端最小 14px
text-sm   // 14px
text-base // 16px (默认)

// ❌ 避免 - 移动端不要使用 text-xs (12px)
text-xs   // 12px - 太小
```

#### 按钮太小
```tsx
// ✅ 推荐
min-h-12  // 最小高度 48px
py-3      // 充足的垂直间距

// ❌ 避免
py-1      // 太小，点击困难
```

#### 输入框太小
```tsx
// ✅ 推荐
input {
  @apply py-3 px-4 text-base min-h-12;
}

// ❌ 避免
input {
  @apply py-1 px-2 text-sm; // 太小
}
```

#### 点击区域太小
```tsx
// ✅ 推荐 - 使用 touch 类
<button className="touch-manipulation">
  按钮
</button>

// 或者在 CSS 中
@layer base {
  button {
    @apply min-h-12 touch-manipulation;
  }
}
```

### 3. 移动端交互优化

```typescript
// 防止双击缩放
<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">

// 禁止长按选择文本（移动端）
* {
  -webkit-touch-callout: none;
  -webkit-user-select: none;
  user-select: none;
}

// 输入框允许选择
input, textarea {
  -webkit-user-select: auto;
  user-select: auto;
}
```

## 📝 注意事项

1. **优先使用全局组件** - 先检查 [shared/components](../../shared/components/) 是否已有可用组件
2. **避免重复造轮子** - 与 user 端组件保持统一（如果功能相似）
3. **保持组件轻量** - 移动端要考虑性能
4. **提供加载状态** - 网络请求要有 loading 状态
5. **处理边界情况** - 空状态、无数据状态
6. **添加注释** - 复杂的业务逻辑需要说明
7. **统一数据格式** - 与后端接口格式保持一致

---

**最后更新**: 2025-11-22
**维护者**: GameLink 前端团队
**特点**: 🎯 移动端优先 | 📱 陪玩师专属 | 🎨 简洁易用
