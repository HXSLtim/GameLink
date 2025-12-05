# StatCard Component

统计卡片组件 - 用于仪表盘展示关键指标

## 功能特性

✅ **已实现的功能**：

1. **统计卡片UI布局** - 清晰的卡片布局，包含标题、图标、数值和趋势信息
2. **加载动画效果** - 使用 Ant Design Skeleton 组件实现加载状态
3. **趋势指示器** - 显示上升/下降箭头和百分比变化
4. **点击跳转功能** - 支持点击卡片触发回调函数
5. **数字滚动动画** - 平滑的数字过渡效果（使用 useCountUp Hook）
6. **平滑过渡效果** - 使用缓动函数实现自然的动画效果

## 使用示例

### 基础用法

```tsx
import { StatCard } from '@/components';
import { UserOutlined } from '@ant-design/icons';

<StatCard
  title="总用户数"
  value={1000}
  icon={<UserOutlined />}
/>
```

### 带趋势指示器

```tsx
<StatCard
  title="总订单数"
  value={500}
  icon={<ShoppingCartOutlined />}
  iconBgColor="#52c41a"
  trend={15}  // 正数表示上升15%
  trendLabel="较昨日"
/>
```

### 带加载状态

```tsx
<StatCard
  title="总收入"
  value={50000}
  icon={<DollarOutlined />}
  loading={true}
  prefix="¥"
/>
```

### 带点击事件

```tsx
<StatCard
  title="活跃陪玩师"
  value={100}
  icon={<TeamOutlined />}
  onClick={() => navigate('/admin/players')}
/>
```

### 带提示信息

```tsx
<StatCard
  title="转化率"
  value={5.2}
  icon={<RiseOutlined />}
  tooltip="访问用户中完成下单的比例"
  suffix="%"
  precision={1}
/>
```

### 自定义动画

```tsx
<StatCard
  title="实时在线用户"
  value={1234}
  icon={<UserOutlined />}
  animated={true}
  animationDuration={2000}  // 2秒动画
/>
```

### 禁用动画

```tsx
<StatCard
  title="静态数据"
  value={999}
  icon={<InfoCircleOutlined />}
  animated={false}  // 禁用数字滚动动画
/>
```

### 完整示例

```tsx
import React from 'react';
import { Row, Col } from 'antd';
import { StatCard } from '@/components';
import {
  UserOutlined,
  ShoppingCartOutlined,
  DollarOutlined,
  TeamOutlined,
} from '@ant-design/icons';

const Dashboard: React.FC = () => {
  const [loading, setLoading] = useState(false);
  
  return (
    <Row gutter={[16, 16]}>
      <Col xs={24} sm={12} lg={6}>
        <StatCard
          title="总用户数"
          value={10000}
          icon={<UserOutlined />}
          iconBgColor="#1890ff"
          trend={12.5}
          trendLabel="较上周"
          loading={loading}
          onClick={() => console.log('Navigate to users')}
        />
      </Col>
      
      <Col xs={24} sm={12} lg={6}>
        <StatCard
          title="总订单数"
          value={5000}
          icon={<ShoppingCartOutlined />}
          iconBgColor="#52c41a"
          trend={-5.2}
          trendLabel="较上周"
          loading={loading}
        />
      </Col>
      
      <Col xs={24} sm={12} lg={6}>
        <StatCard
          title="总收入"
          value={250000}
          icon={<DollarOutlined />}
          iconBgColor="#faad14"
          prefix="¥"
          precision={2}
          trend={8.3}
          loading={loading}
        />
      </Col>
      
      <Col xs={24} sm={12} lg={6}>
        <StatCard
          title="活跃陪玩师"
          value={500}
          icon={<TeamOutlined />}
          iconBgColor="#722ed1"
          tooltip="最近30天内有订单的陪玩师数量"
          loading={loading}
        />
      </Col>
    </Row>
  );
};
```

## Props API

| 属性 | 说明 | 类型 | 默认值 | 必填 |
|------|------|------|--------|------|
| title | 卡片标题 | ReactNode | - | ✅ |
| value | 统计数值 | number \| string | - | ✅ |
| icon | 图标 | ReactNode | - | ❌ |
| iconBgColor | 图标背景色 | string | '#1890ff' | ❌ |
| trend | 趋势（正数上升，负数下降） | number | - | ❌ |
| trendLabel | 趋势描述文字 | string | '较昨日' | ❌ |
| tooltip | 提示信息 | string | - | ❌ |
| footer | 底部描述 | ReactNode | - | ❌ |
| loading | 是否加载中 | boolean | false | ❌ |
| onClick | 点击回调 | () => void | - | ❌ |
| animated | 是否启用数字动画 | boolean | true | ❌ |
| animationDuration | 动画持续时间（毫秒） | number | 1000 | ❌ |
| prefix | 数值前缀 | string | - | ❌ |
| suffix | 数值后缀 | string | - | ❌ |
| precision | 小数位数 | number | 0 | ❌ |

## 数字动画 Hook

StatCard 内部使用了 `useCountUp` Hook 来实现数字滚动动画。

### useCountUp API

```tsx
import { useCountUp } from '@/hooks/useCountUp';

const animatedValue = useCountUp(targetValue, {
  duration: 1000,        // 动画持续时间
  enabled: true,         // 是否启用动画
  easing: (t) => t,      // 缓动函数
});
```

### 缓动函数

默认使用 `easeOutExpo` 缓动函数，提供自然的减速效果：

```typescript
const easeOutExpo = (t: number): number => {
  return t === 1 ? 1 : 1 - Math.pow(2, -10 * t);
};
```

## 样式定制

组件使用 CSS Modules，可以通过覆盖以下类名来定制样式：

- `.card` - 卡片容器
- `.clickable` - 可点击状态
- `.header` - 头部区域
- `.titleWrapper` - 标题包装器
- `.title` - 标题文字
- `.iconWrapper` - 图标容器
- `.body` - 主体区域
- `.footer` - 底部区域
- `.trend` - 趋势指示器
- `.trendLabel` - 趋势标签

## 验证需求

该组件实现了以下需求：

- ✅ **需求 1.1**: 显示核心统计卡片
- ✅ **需求 1.2**: 显示加载动画
- ✅ **需求 1.3**: 显示趋势百分比和图标
- ✅ **需求 1.4**: 支持点击跳转
- ✅ **需求 1.5**: 数值动画效果

## 注意事项

1. **数字动画**：当 `value` 为字符串时，会尝试解析为数字进行动画。如果解析失败，将直接显示字符串。
2. **性能优化**：使用 `useMemo` 优化计算，避免不必要的重渲染。
3. **响应式设计**：建议配合 Ant Design 的 Grid 系统使用，确保在不同屏幕尺寸下正常显示。
4. **加载状态**：当 `loading={true}` 时，会显示骨架屏，此时动画会被禁用。

## 相关文件

- 组件实现: `frontend/src/components/StatCard/index.tsx`
- 样式文件: `frontend/src/components/StatCard/index.module.css`
- 动画Hook: `frontend/src/hooks/useCountUp.ts`
- 类型定义: `frontend/src/types/dashboard.ts`
- 常量定义: `frontend/src/constants/dashboard.ts`
