# GameLink 黑白设计系统

**版本**: 2.0.0  
**更新日期**: 2025-10-28  
**设计风格**: Neo-brutalism 极简黑白风格  
**UI 框架**: 自定义组件库

---

## 📖 关于本设计系统

GameLink 管理系统采用 **Neo-brutalism（新野兽主义）** 设计风格，以**纯黑白配色**为核心，追求极简、直接、高对比度的视觉体验。

### 设计理念

> **极简即力量** - 去除一切不必要的装饰，用最纯粹的黑白演绎现代感

1. **极简主义** - 只保留必要元素，去除多余装饰
2. **高对比度** - 纯黑白搭配，清晰醒目
3. **直角设计** - 无圆角，直线条，展现力量感
4. **实体阴影** - 使用黑色实体阴影而非柔和阴影
5. **粗边框** - 2px 标准边框，强调边界感

---

## 🎨 色彩系统

### 核心颜色

```less
// 主色调 - 纯黑白
--color-white: #ffffff; // 纯白
--color-black: #000000; // 纯黑

// 灰度色阶（辅助色）
--color-gray-100: #f5f5f5; // 最浅灰
--color-gray-200: #e5e5e5;
--color-gray-300: #d4d4d4;
--color-gray-400: #a3a3a3;
--color-gray-500: #737373; // 中灰
--color-gray-600: #525252;
--color-gray-700: #404040;
--color-gray-800: #262626;
--color-gray-900: #171717; // 最深灰
```

### 语义化颜色

```less
// 文本颜色
--text-primary: #000000; // 主要文本
--text-secondary: #666666; // 次要文本
--text-tertiary: #999999; // 辅助文本
--text-disabled: #cccccc; // 禁用文本
--text-inverse: #ffffff; // 反色文本

// 背景颜色
--bg-primary: #ffffff; // 主背景
--bg-secondary: #fafafa; // 次要背景
--bg-tertiary: #f5f5f5; // 三级背景
--bg-inverse: #000000; // 反色背景
```

### 色彩使用原则

1. **主色为王** - 界面以白色为主，黑色为辅
2. **高对比度** - 文字与背景对比度 ≥ 7:1（WCAG AAA 级）
3. **限制灰度** - 仅在必要时使用灰度过渡
4. **无彩色** - 严禁使用任何彩色（红、绿、蓝等）

---

## 🧩 组件设计

### 组件列表

#### 1. Button（按钮）

**变体**：

- `primary` - 黑底白字（主要操作）
- `secondary` - 白底黑字（次要操作）
- `text` - 透明背景（文本按钮）

**尺寸**：

- `small` - 高度 32px
- `medium` - 高度 40px（默认）
- `large` - 高度 48px

**特性**：

- 2px 黑色边框
- 实体阴影（4px × 4px）
- 悬停时阴影增大 + 上浮动画
- 加载状态旋转动画
- 水波纹展开效果

```tsx
import { Button } from 'components';

<Button variant="primary" size="large">
  登录
</Button>;
```

#### 2. Input（输入框）

**特性**：

- 2px 黑色边框
- 聚焦时实体阴影
- 支持前缀/后缀图标
- 支持清空按钮
- 错误状态提示

```tsx
import { Input } from 'components';

<Input prefix={<UserIcon />} placeholder="用户名" allowClear />;
```

#### 3. PasswordInput（密码输入框）

**特性**：

- 继承 Input 所有特性
- 显示/隐藏密码切换
- 眼睛图标按钮

```tsx
import { PasswordInput } from 'components';

<PasswordInput prefix={<LockIcon />} placeholder="密码" />;
```

#### 4. Card（卡片）

**特性**：

- 2px 黑色边框
- 8px × 8px 实体阴影
- 可选标题和额外内容
- 悬停时上浮动画（可选）

```tsx
import { Card } from 'components';

<Card title="标题" hoverable>
  内容
</Card>;
```

#### 5. Form（表单）

**布局**：

- `vertical` - 垂直布局（默认）
- `horizontal` - 水平布局

**FormItem 特性**：

- 标签显示
- 必填标识（\*）
- 错误提示
- 帮助文本

```tsx
import { Form, FormItem } from 'components';

<Form onSubmit={handleSubmit}>
  <FormItem label="用户名" required error={errors.username}>
    <Input />
  </FormItem>
</Form>;
```

---

## 📐 布局系统

### 间距规范

基于 **4px** 基准单位：

```less
--spacing-xs: 4px; // 0.25rem
--spacing-sm: 8px; // 0.5rem
--spacing-md: 12px; // 0.75rem
--spacing-base: 16px; // 1rem (基准)
--spacing-lg: 20px; // 1.25rem
--spacing-xl: 24px; // 1.5rem
--spacing-2xl: 32px; // 2rem
--spacing-3xl: 40px; // 2.5rem
--spacing-4xl: 48px; // 3rem
```

**使用建议**：

- 组件内边距：16px
- 表单间距：24px
- 卡片内边距：24px
- 页面边距：32px

### 边框规范

```less
// 边框宽度
--border-width-thin: 1px;
--border-width-base: 2px; // 标准宽度
--border-width-thick: 3px;
--border-width-heavy: 4px;

// 圆角 - Neo-brutalism 使用直角
--border-radius-none: 0; // 标准（无圆角）
--border-radius-sm: 2px; // 可选
--border-radius-base: 4px; // 可选
```

### 阴影规范

**实体阴影** - Neo-brutalism 特色：

```less
--shadow-xs: 2px 2px 0 #000000;
--shadow-sm: 4px 4px 0 #000000;
--shadow-base: 8px 8px 0 #000000; // 标准
--shadow-lg: 12px 12px 0 #000000;
--shadow-xl: 16px 16px 0 #000000;
```

**使用原则**：

- 卡片：8px × 8px
- 按钮：4px × 4px
- 输入框聚焦：4px × 4px
- 悬停效果：增大阴影 + 位移

---

## ✏️ 字体系统

### 字体家族

```less
--font-family-base:
  -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, 'PingFang SC',
  'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;

--font-family-mono: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, monospace;
```

### 字体大小

```less
--font-size-xs: 12px;
--font-size-sm: 14px;
--font-size-base: 16px; // 基准
--font-size-lg: 18px;
--font-size-xl: 20px;
--font-size-2xl: 24px;
--font-size-3xl: 28px;
--font-size-4xl: 32px;
--font-size-5xl: 36px;
--font-size-6xl: 48px;
```

### 字重

```less
--font-weight-light: 300;
--font-weight-normal: 400; // 正文
--font-weight-medium: 500;
--font-weight-semibold: 600; // 按钮、标题
--font-weight-bold: 700; // 强调
--font-weight-black: 900;
```

### 行高

```less
--line-height-tight: 1.25; // 标题
--line-height-normal: 1.5; // 正文（默认）
--line-height-relaxed: 1.75;
--line-height-loose: 2;
```

---

## 🎬 动画规范

### 动画时长

```less
--duration-instant: 0ms;
--duration-fast: 100ms; // 微交互
--duration-base: 200ms; // 标准动画（默认）
--duration-medium: 300ms;
--duration-slow: 400ms;
--duration-slower: 600ms; // 复杂动画
```

### 缓动函数

```less
--ease-linear: cubic-bezier(0, 0, 1, 1);
--ease-in: cubic-bezier(0.4, 0, 1, 1);
--ease-out: cubic-bezier(0, 0, 0.2, 1); // 常用
--ease-in-out: cubic-bezier(0.4, 0, 0.2, 1); // 常用
--ease-bounce: cubic-bezier(0.68, -0.55, 0.265, 1.55);
```

### 动画类型

#### 1. 页面进入动画

```less
@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(40px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
```

#### 2. 交互动画

- **悬停效果**: 阴影增大 + 元素上浮（-2px, -2px）
- **点击效果**: 阴影缩小 + 元素回位
- **加载效果**: 旋转动画（0.6s linear infinite）

#### 3. 水波纹效果

按钮点击时的圆形扩散效果：

```less
&::before {
  content: '';
  position: absolute;
  width: 0;
  height: 0;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  transition:
    width 0.6s,
    height 0.6s;
}

&:hover::before {
  width: 300px;
  height: 300px;
}
```

### 动画使用原则

1. **性能优先** - 仅使用 `transform` 和 `opacity`
2. **适度使用** - 避免过度动画
3. **尊重用户** - 支持 `prefers-reduced-motion`
4. **目的明确** - 动画应帮助理解

---

## 📱 响应式设计

### 断点定义

```less
--breakpoint-xs: 0px;
--breakpoint-sm: 576px; // 手机横屏
--breakpoint-md: 768px; // 平板竖屏
--breakpoint-lg: 992px; // 平板横屏
--breakpoint-xl: 1200px; // 桌面显示器
--breakpoint-2xl: 1600px; // 大屏
```

### 移动端适配

**登录页面示例**：

- 桌面端：固定宽度 420px
- 移动端：90% 宽度，最大 420px
- 阴影尺寸自适应减小

```less
@media (max-width: 768px) {
  .loginCard {
    width: 90%;
    box-shadow: var(--shadow-sm);
  }
}
```

---

## ♿ 可访问性

### 对比度标准

遵循 **WCAG 2.1 AAA 级**标准：

- 普通文本：≥ 7:1
- 大文本：≥ 4.5:1
- UI 组件：≥ 3:1

黑白配色天然满足高对比度要求。

### 键盘导航

所有交互元素支持：

| 按键          | 功能     |
| ------------- | -------- |
| Tab           | 焦点前进 |
| Shift + Tab   | 焦点后退 |
| Enter / Space | 激活按钮 |
| Esc           | 关闭弹窗 |

### 焦点样式

```less
*:focus-visible {
  outline: 2px solid var(--color-black);
  outline-offset: 2px;
}
```

### 动画禁用

```less
@media (prefers-reduced-motion: reduce) {
  * {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
  }
}
```

---

## 🎯 组件使用示例

### 登录页面

```tsx
import { useState } from 'react';
import { Button, Input, PasswordInput, Form, FormItem } from 'components';

export const Login = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    // ... 登录逻辑
  };

  return (
    <div className="login-container">
      <div className="login-card">
        <h1>GameLink</h1>
        <p>欢迎回来</p>

        <Form onSubmit={handleSubmit}>
          <FormItem>
            <Input
              prefix={<UserIcon />}
              placeholder="用户名"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              allowClear
            />
          </FormItem>

          <FormItem>
            <PasswordInput
              prefix={<LockIcon />}
              placeholder="密码"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </FormItem>

          <FormItem>
            <Button type="submit" variant="primary" size="large" block loading={loading}>
              登录
            </Button>
          </FormItem>
        </Form>
      </div>
    </div>
  );
};
```

---

## 📚 资源和工具

### 设计参考

- **Dribbble**: 搜索 "Neo-brutalism" / "Black and White Design"
- **Behance**: 搜索 "Brutalism Web Design"
- **Awwwards**: 极简主义网站案例

### 开发工具

- **VSCode**: 代码编辑器
- **Chrome DevTools**: 调试工具
- **React Developer Tools**: React 调试

### 灵感网站

- [Brutalist Websites](https://brutalistwebsites.com/)
- [Hoverstat.es](https://www.hoverstat.es/)
- [Minimal Gallery](https://minimal.gallery/)

---

## 🎯 设计原则总结

### DO ✅

1. ✅ 使用纯黑白配色
2. ✅ 使用直角设计（无圆角）
3. ✅ 使用 2px 粗边框
4. ✅ 使用实体阴影
5. ✅ 保持高对比度
6. ✅ 保持极简风格
7. ✅ 使用适度动画

### DON'T ❌

1. ❌ 使用任何彩色
2. ❌ 使用圆角（除非特殊需要）
3. ❌ 使用柔和阴影
4. ❌ 使用渐变（除非装饰用）
5. ❌ 过度装饰
6. ❌ 低对比度配色
7. ❌ 过度动画

---

## 📝 更新日志

### v2.0.0 (2025-10-28)

- ✅ 从 Arco Design 迁移到自定义组件库
- ✅ 采用 Neo-brutalism 设计风格
- ✅ 纯黑白配色系统
- ✅ 创建基础组件（Button, Input, Card, Form）
- ✅ 实体阴影系统
- ✅ 直角设计规范
- ✅ 重构登录页面

---

**维护者**: GameLink Frontend Team  
**设计风格**: Neo-brutalism 黑白极简  
**最后更新**: 2025-10-28

---

<div align="center">

**纯粹黑白 · 极简力量**

Made with ⚫⚪ by GameLink Team

</div>
