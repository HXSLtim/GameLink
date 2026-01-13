# GameLink 小程序色卡

> **版本**: v1.0.0
> **更新时间**: 2025-01-10
> **用途**: 开发参考、设计规范

---

## 📋 目录

1. [品牌色彩](#品牌色彩)
2. [功能色彩](#功能色彩)
3. [中性色彩](#中性色彩)
4. [渐变色彩](#渐变色彩)
5. [色彩组合](#色彩组合)
6. [CSS 变量](#css-变量)
7. [使用规范](#使用规范)

---

## 品牌色彩

### 主色调 - 紫色系

```
┌─────────────────────────────────────────────────────────────┐
│  Primary Color - 紫色                                      │
├─────────────────────────────────────────────────────────────┤
│  ■ $primary-color      #667eea  RGB(102, 126, 234)        │
│  ■ $primary-dark       #764ba2  RGB(118, 75, 162)         │
│  ■ $primary-light      #8b9eef  RGB(139, 158, 239)        │
│  ■ $primary-lighter    #a5b3f5  RGB(165, 179, 245)       │
│  ■ $primary-lightest   #e8ecfe  RGB(232, 236, 254)       │
└─────────────────────────────────────────────────────────────┘
```

#### 色值预览

| 颜色名称 | 色值 | RGB | 预览 |
|----------|------|-----|------|
| Primary | #667eea | 102, 126, 234 | ![Primary](https://via.placeholder.com/50/667eea/667eea) |
| Primary Dark | #764ba2 | 118, 75, 162 | ![Primary Dark](https://via.placeholder.com/50/764ba2/764ba2) |
| Primary Light | #8b9eef | 139, 158, 239 | ![Primary Light](https://via.placeholder.com/50/8b9eef/8b9eef) |
| Primary Lighter | #a5b3f5 | 165, 179, 245 | ![Primary Lighter](https://via.placeholder.com/50/a5b3f5/a5b3f5) |
| Primary Lightest | #e8ecfe | 232, 236, 254 | ![Primary Lightest](https://via.placeholder.com/50/e8ecfe/e8ecfe) |

#### 使用场景

| 颜色 | 使用场景 | 代码示例 |
|------|----------|----------|
| Primary | 主按钮、导航栏选中、重要操作 | `<Button type="primary">` |
| Primary Dark | 主按钮悬停、深色背景 | `.btn-primary:hover` |
| Primary Light | 次要强调、标签背景 | `<Tag type="primary">` |
| Primary Lighter | 浅色背景、分隔线 | `.bg-primary-lighter` |
| Primary Lightest | 极浅背景、页面底色 | `.page-bg` |

#### SCSS 变量

```scss
// 主色调
$primary-color: #667eea;
$primary-dark: #764ba2;
$primary-light: #8b9eef;
$primary-lighter: #a5b3f5;
$primary-lightest: #e8ecfe;

// 主色调渐变
$primary-gradient: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
```

---

### 强调色 - 橙色系

```
┌─────────────────────────────────────────────────────────────┐
│  Accent Color - 橙色                                       │
├─────────────────────────────────────────────────────────────┤
│  ■ $accent-color       #f59e0b  RGB(245, 158, 11)         │
│  ■ $accent-dark        #d97706  RGB(217, 119, 6)          │
│  ■ $accent-light       #fbbf24  RGB(251, 191, 36)         │
│  ■ $accent-lighter     #fcd34d  RGB(252, 211, 77)         │
│  ■ $accent-lightest    #fef3c7  RGB(254, 243, 199)        │
└─────────────────────────────────────────────────────────────┘
```

#### 色值预览

| 颜色名称 | 色值 | RGB | 预览 |
|----------|------|-----|------|
| Accent | #f59e0b | 245, 158, 11 | ![Accent](https://via.placeholder.com/50/f59e0b/f59e0b) |
| Accent Dark | #d97706 | 217, 119, 6 | ![Accent Dark](https://via.placeholder.com/50/d97706/d97706) |
| Accent Light | #fbbf24 | 251, 191, 36 | ![Accent Light](https://via.placeholder.com/50/fbbf24/fbbf24) |
| Accent Lighter | #fcd34d | 252, 211, 77 | ![Accent Lighter](https://via.placeholder.com/50/fcd34d/fcd34d) |
| Accent Lightest | #fef3c7 | 254, 243, 199 | ![Accent Lightest](https://via.placeholder.com/50/fef3c7/fef3c7) |

#### 使用场景

| 颜色 | 使用场景 | 示例 |
|------|----------|------|
| Accent | 价格、重要数字、CTA 按钮 | ¥50.00 |
| Accent Dark | 价格强调、悬停状态 | VIP 价格 |
| Accent Light | 标签背景、徽章 | 新标签 |
| Accent Lighter | 浅色背景 | 优惠区域 |
| Accent Lightest | 极浅背景 | 活动卡片 |

#### SCSS 变量

```scss
// 强调色
$accent-color: #f59e0b;
$accent-dark: #d97706;
$accent-light: #fbbf24;
$accent-lighter: #fcd34d;
$accent-lightest: #fef3c7;

// 强调色渐变
$accent-gradient: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
```

---

## 功能色彩

### 成功色 - 绿色系

```
┌─────────────────────────────────────────────────────────────┐
│  Success Color - 绿色                                     │
├─────────────────────────────────────────────────────────────┤
│  ■ $success-color     #10b981  RGB(16, 185, 129)         │
│  ■ $success-dark      #059669  RGB(5, 150, 105)          │
│  ■ $success-light     #34d399  RGB(52, 211, 153)         │
│  ■ $success-lighter   #6ee7b7  RGB(110, 231, 183)        │
│  ■ $success-lightest  #d1fae5  RGB(209, 250, 229)        │
└─────────────────────────────────────────────────────────────┘
```

#### 使用场景

- ✅ 订单完成、支付成功
- ✅ 操作成功提示
- ✅ 在线状态指示
- ✅ 验证通过状态

### 警告色 - 黄色系

```
┌─────────────────────────────────────────────────────────────┐
│  Warning Color - 黄色                                     │
├─────────────────────────────────────────────────────────────┤
│  ■ $warning-color     #f59e0b  RGB(245, 158, 11)         │
│  ■ $warning-dark      #d97706  RGB(217, 119, 6)          │
│  ■ $warning-light     #fbbf24  RGB(251, 191, 36)         │
│  ■ $warning-lighter   #fcd34d  RGB(252, 211, 77)         │
│  ■ $warning-lightest  #fef3c7  RGB(254, 243, 199)        │
└─────────────────────────────────────────────────────────────┘
```

#### 使用场景

- ⚠️ 余额不足、时间不足
- ⚠️ 重要操作确认
- ⚠️ 到期提醒
- ⚠️ 忙碌状态

### 错误色 - 红色系

```
┌─────────────────────────────────────────────────────────────┐
│  Error Color - 红色                                       │
├─────────────────────────────────────────────────────────────┤
│  ■ $error-color      #ef4444  RGB(239, 68, 68)           │
│  ■ $error-dark       #dc2626  RGB(220, 38, 38)           │
│  ■ $error-light      #f87171  RGB(248, 113, 113)         │
│  ■ $error-lighter    #fca5a5  RGB(252, 165, 165)         │
│  ■ $error-lightest   #fee2e2  RGB(254, 226, 226)         │
└─────────────────────────────────────────────────────────────┘
```

#### 使用场景

- ❌ 网络错误、操作失败
- ❌ 表单验证错误
- ❌ 权限不足
- ❌ 离线状态

### 信息色 - 蓝色系

```
┌─────────────────────────────────────────────────────────────┐
│  Info Color - 蓝色                                        │
├─────────────────────────────────────────────────────────────┤
│  ■ $info-color       #3b82f6  RGB(59, 130, 246)          │
│  ■ $info-dark        #2563eb  RGB(37, 99, 235)           │
│  ■ $info-light       #60a5fa  RGB(96, 165, 250)          │
│  ■ $info-lighter     #93c5fd  RGB(147, 197, 253)         │
│  ■ $info-lightest    #dbeafe  RGB(219, 234, 254)         │
└─────────────────────────────────────────────────────────────┘
```

#### 使用场景

- ℹ️ 系统通知
- ℹ️ 帮助提示
- ℹ️ 信息气泡
- ℹ️ 链接文字

---

## 中性色彩

### 文字颜色

```
┌─────────────────────────────────────────────────────────────┐
│  Text Colors                                               │
├─────────────────────────────────────────────────────────────┤
│  ■ $text-primary      #1f2937  RGB(31, 41, 55)           │
│  ■ $text-secondary    #6b7280  RGB(107, 114, 128)        │
│  ■ $text-tertiary     #9ca3af  RGB(156, 163, 175)        │
│  ■ $text-quaternary   #d1d5db  RGB(209, 213, 219)        │
│  ■ $text-disabled     #e5e7eb  RGB(229, 231, 235)        │
│  ■ $text-inverse      #ffffff  RGB(255, 255, 255)        │
└─────────────────────────────────────────────────────────────┘
```

#### 文字颜色使用规范

| 颜色 | 使用场景 | 示例 | 对比度 |
|------|----------|------|--------|
| text-primary | 标题、正文 | 页面标题、卡片正文 | 13.4:1 ✅ |
| text-secondary | 副标题、辅助信息 | 时间戳、辅助说明 | 6.2:1 ✅ |
| text-tertiary | 提示文字 | placeholder、禁用状态 | 3.9:1 ⚠️ |
| text-quaternary | 弱化文字 | 极次要信息 | 2.8:1 ❌ |
| text-disabled | 禁用文字 | 禁用状态按钮 | 2.1:1 ❌ |
| text-inverse | 反白文字（深色背景） | 按钮文字、标签 | 12.6:1 ✅ |

### 背景颜色

```
┌─────────────────────────────────────────────────────────────┐
│  Background Colors                                         │
├─────────────────────────────────────────────────────────────┤
│  ■ $bg-page          #f9fafb  RGB(249, 250, 251)         │
│  ■ $bg-card          #ffffff  RGB(255, 255, 255)         │
│  ■ $bg-input         #f3f4f6  RGB(243, 244, 246)         │
│  ■ $bg-active        #f3f4f6  RGB(243, 244, 246)         │
│  ■ $bg-hover         #f9fafb  RGB(249, 250, 251)         │
│  ■ $bg-disabled      #f3f4f6  RGB(243, 244, 246)         │
│  ■ $bg-overlay       rgba(0,0,0,0.5)                      │
└─────────────────────────────────────────────────────────────┘
```

#### 背景颜色使用规范

| 颜色 | 使用场景 | 示例 |
|------|----------|------|
| bg-page | 页面整体背景 | 整个页面背景色 |
| bg-card | 卡片、弹窗背景 | 白色卡片背景 |
| bg-input | 输入框背景 | 表单输入框 |
| bg-active | 点击态背景 | 按钮按下状态 |
| bg-hover | 悬停背景 | 列表项悬停 |
| bg-disabled | 禁用背景 | 禁用状态背景 |
| bg-overlay | 遮罩背景 | 弹窗遮罩 |

### 边框颜色

```
┌─────────────────────────────────────────────────────────────┐
│  Border Colors                                             │
├─────────────────────────────────────────────────────────────┤
│  ■ $border-color      #e5e7eb  RGB(229, 231, 235)        │
│  ■ $divider-color     #f3f4f6  RGB(243, 244, 246)         │
│  ■ $border-primary    #667eea  RGB(102, 126, 234)         │
│  ■ $border-success    #10b981  RGB(16, 185, 129)         │
│  ■ $border-warning    #f59e0b  RGB(245, 158, 11)          │
│  ■ $border-error      #ef4444  RGB(239, 68, 68)           │
└─────────────────────────────────────────────────────────────┘
```

---

## 渐变色彩

### 品牌渐变

```scss
// 主色渐变（紫色）
$gradient-primary: linear-gradient(135deg, #667eea 0%, #764ba2 100%);

// 强调色渐变（橙色）
$gradient-accent: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);

// 成功渐变（绿色）
$gradient-success: linear-gradient(135deg, #10b981 0%, #059669 100%);

// 错误渐变（红色）
$gradient-error: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
```

### 背景渐变

```scss
// 页面背景渐变
$gradient-bg: linear-gradient(180deg, #f9fafb 0%, #ffffff 100%);

// 卡片背景渐变
$gradient-card: linear-gradient(135deg, #ffffff 0%, #f9fafb 100%);

// 按钮背景渐变
$gradient-btn: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
```

### 文字渐变

```scss
// VIP 标题渐变
$gradient-text-vip: linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%);

// 品牌文字渐变
$gradient-text-brand: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
```

#### 渐变使用示例

```tsx
// 渐变背景
<View style={{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }}>
  <Text className="text-white">渐变背景</Text>
</View>

// 渐变文字
<Text style={{
  background: 'linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%)',
  WebkitBackgroundClip: 'text',
  WebkitTextFillColor: 'transparent'
}}>
  VIP 会员
</Text>
```

---

## 色彩组合

### 常用色彩组合

#### 1. 品牌组合（紫色 + 白色）

```
主按钮: Primary Gradient
背景色: White
文字色: Text Primary
边框色: Border Color
```

**应用场景**: 主要操作按钮、导航栏

#### 2. 强调组合（橙色 + 白色）

```
价格: Accent Color
背景色: Accent Lightest
文字色: Text Primary
边框色: Accent Lighter
```

**应用场景**: VIP 价格、优惠标签

#### 3. 成功组合（绿色 + 白色）

```
状态: Success Color
背景色: Success Lightest
文字色: Text Primary
图标: Success Color
```

**应用场景**: 在线状态、订单完成

#### 4. 错误组合（红色 + 白色）

```
状态: Error Color
背景色: Error Lightest
文字色: Error Color
图标: Error Color
```

**应用场景**: 离线状态、错误提示

### 对比度标准

| 组合 | 前景色 | 背景色 | 对比度 | 等级 |
|------|--------|--------|--------|------|
| ✅ Primary on White | #667eea | #ffffff | 3.1:1 | AA Large |
| ✅ Text on White | #1f2937 | #ffffff | 13.4:1 | AAA |
| ✅ White on Primary | #ffffff | #667eea | 12.6:1 | AAA |
| ✅ Accent on White | #f59e0b | #ffffff | 2.8:1 | AA Large |
| ⚠️ Tertiary on White | #9ca3af | #ffffff | 3.9:1 | AA Large |
| ❌ Disabled on White | #e5e7eb | #ffffff | 1.2:1 | Fail |

**说明**:
- **AAA**: 对比度 ≥ 7:1，最佳可读性
- **AA**: 对比度 ≥ 4.5:1，正常文字可读性
- **AA Large**: 对比度 ≥ 3:1，大号文字（18pt+）可读性

---

## CSS 变量

### 完整 CSS 变量定义

```css
:root {
  /* 品牌色彩 - 紫色 */
  --color-primary: #667eea;
  --color-primary-dark: #764ba2;
  --color-primary-light: #8b9eef;
  --color-primary-lighter: #a5b3f5;
  --color-primary-lightest: #e8ecfe;

  /* 强调色 - 橙色 */
  --color-accent: #f59e0b;
  --color-accent-dark: #d97706;
  --color-accent-light: #fbbf24;
  --color-accent-lighter: #fcd34d;
  --color-accent-lightest: #fef3c7;

  /* 功能色彩 */
  --color-success: #10b981;
  --color-success-dark: #059669;
  --color-success-light: #34d399;
  --color-success-lighter: #6ee7b7;
  --color-success-lightest: #d1fae5;

  --color-warning: #f59e0b;
  --color-warning-dark: #d97706;
  --color-warning-light: #fbbf24;
  --color-warning-lighter: #fcd34d;
  --color-warning-lightest: #fef3c7;

  --color-error: #ef4444;
  --color-error-dark: #dc2626;
  --color-error-light: #f87171;
  --color-error-lighter: #fca5a5;
  --color-error-lightest: #fee2e2;

  --color-info: #3b82f6;
  --color-info-dark: #2563eb;
  --color-info-light: #60a5fa;
  --color-info-lighter: #93c5fd;
  --color-info-lightest: #dbeafe;

  /* 中性色彩 - 文字 */
  --color-text-primary: #1f2937;
  --color-text-secondary: #6b7280;
  --color-text-tertiary: #9ca3af;
  --color-text-quaternary: #d1d5db;
  --color-text-disabled: #e5e7eb;
  --color-text-inverse: #ffffff;

  /* 中性色彩 - 背景 */
  --color-bg-page: #f9fafb;
  --color-bg-card: #ffffff;
  --color-bg-input: #f3f4f6;
  --color-bg-active: #f3f4f6;
  --color-bg-hover: #f9fafb;
  --color-bg-disabled: #f3f4f6;
  --color-bg-overlay: rgba(0, 0, 0, 0.5);

  /* 中性色彩 - 边框 */
  --color-border: #e5e7eb;
  --color-divider: #f3f4f6;
  --color-border-primary: #667eea;
  --color-border-success: #10b981;
  --color-border-warning: #f59e0b;
  --color-border-error: #ef4444;

  /* 渐变 */
  --gradient-primary: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  --gradient-accent: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  --gradient-success: linear-gradient(135deg, #10b981 0%, #059669 100%);
  --gradient-error: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);

  /* 阴影 */
  --shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
  --shadow-xl: 0 20px 25px -5px rgba(0, 0, 0, 0.1);

  /* 圆角 */
  --radius-xs: 4px;
  --radius-sm: 8px;
  --radius-md: 12px;
  --radius-lg: 16px;
  --radius-xl: 24px;
  --radius-full: 9999px;
}
```

### Taro 项目中使用

```tsx
// app.scss
:root {
  --color-primary: #667eea;
  --color-text-primary: #1f2937;
  // ... 其他变量
}

// 组件中使用
<View style={{
  backgroundColor: 'var(--color-primary)',
  color: 'var(--color-text-inverse)'
}}>
  主要按钮
</View>
```

---

## 使用规范

### 1. 按钮色彩规范

#### 主要按钮

```tsx
<Button type="primary">主要操作</Button>
// 背景: Primary Gradient
// 文字: White
// 悬停: Primary Dark
```

#### 次要按钮

```tsx
<Button type="secondary">次要操作</Button>
// 背景: White
// 边框: Primary
// 文字: Primary
```

#### 危险按钮

```tsx
<Button type="danger">删除</Button>
// 背景: Error Gradient
// 文字: White
```

### 2. 状态色彩规范

#### 在线状态

```tsx
<OnlineStatus status="online" />
// 在线: Success Color
// 忙碌: Warning Color
// 离线: Text Tertiary
```

#### 订单状态

```tsx
<OrderStatus status="completed" />
// 已完成: Success
// 进行中: Primary
// 待支付: Warning
// 已取消: Text Tertiary
```

### 3. 价格显示规范

#### 普通价格

```tsx
<PriceDisplay price={100} />
// 文字颜色: Text Primary
// 字号: 32px (lg)
// 字重: 600 (semibold)
```

#### VIP 折扣价格

```tsx
<PriceDisplay
  price={80}
  originalPrice={100}
  vipLevel={3}
/>
// 折后价: Accent Color
// 原价: Text Tertiary + 删除线
// VIP 标签: Accent Gradient
```

### 4. 标签色彩规范

#### 游戏标签

```tsx
<Tag type="game">王者荣耀</Tag>
// 背景: Primary Lightest
// 文字: Primary
// 边框: Primary Lighter
```

#### VIP 标签

```tsx
<Tag type="vip">VIP3</Tag>
// 背景: Accent Gradient
// 文字: White
// 圆角: Circle
```

#### 状态标签

```tsx
<Tag type="success">在线</Tag>
<Tag type="warning">忙碌</Tag>
<Tag type="error">离线</Tag>
```

---

## 快速参考

### 色彩速查表

| 用途 | 颜色 | 色值 | 变量名 |
|------|------|------|--------|
| 主按钮 | 紫色 | #667eea | `--color-primary` |
| 价格 | 橙色 | #f59e0b | `--color-accent` |
| 成功 | 绿色 | #10b981 | `--color-success` |
| 警告 | 黄色 | #f59e0b | `--color-warning` |
| 错误 | 红色 | #ef4444 | `--color-error` |
| 主文字 | 深灰 | #1f2937 | `--color-text-primary` |
| 次要文字 | 中灰 | #6b7280 | `--color-text-secondary` |
| 页面背景 | 浅灰 | #f9fafb | `--color-bg-page` |
| 卡片背景 | 白色 | #ffffff | `--color-bg-card` |
| 边框 | 极浅灰 | #e5e7eb | `--color-border` |

### SCSS 完整变量

```scss
// 品牌色彩
$primary-color: #667eea;
$primary-dark: #764ba2;
$primary-light: #8b9eef;
$primary-lighter: #a5b3f5;
$primary-lightest: #e8ecfe;

// 强调色
$accent-color: #f59e0b;
$accent-dark: #d97706;
$accent-light: #fbbf24;
$accent-lighter: #fcd34d;
$accent-lightest: #fef3c7;

// 功能色彩
$success-color: #10b981;
$warning-color: #f59e0b;
$error-color: #ef4444;
$info-color: #3b82f6;

// 文字色彩
$text-primary: #1f2937;
$text-secondary: #6b7280;
$text-tertiary: #9ca3af;
$text-disabled: #d1d5db;
$text-inverse: #ffffff;

// 背景色彩
$bg-page: #f9fafb;
$bg-card: #ffffff;
$bg-input: #f3f4f6;
$bg-active: #f3f4f6;

// 边框色彩
$border-color: #e5e7eb;
$divider-color: #f3f4f6;

// 渐变
$gradient-primary: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
$gradient-accent: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
```

---

## 附录

### 相关文档

- [Taro UI 设计系统](./TARO_UI_DESIGN_SYSTEM.md)
- [原型设计文档](./MINIPROTOTYPE_DESIGN.md)
- [配置数据汇总](./CONFIGURATION_DATA.md)

### 色彩工具

- 在线对比度检查: https://contrast-ratio.com/
- 色彩对比生成: https://coolors.co/
- 渐变生成器: https://cssgradient.io/

### 联系方式

如有设计问题，请联系:
- **设计团队**: design@gamelink.com
- **开发团队**: dev@gamelink.com

---

**维护者**: GameLink 设计团队
**最后更新**: 2025-01-10
**版本**: v1.0.0
