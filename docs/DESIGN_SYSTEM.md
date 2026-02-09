# GameLink 设计系统规范

**版本**: 1.0.0
**更新日期**: 2026-02-09
**维护者**: Mobile-Lead, Frontend-Lead

---

## 设计理念

### 移动端 (App)
- **风格**: 活泼、友好、游戏风
- **参考**: KOOK、Discord 移动端
- **特点**: 大圆角、鲜艳色彩、丰富的动画效果

### 管理后台 (Admin)
- **风格**: 专业、简洁、克制
- **参考**: Ant Design、Material Design
- **特点**: 清晰的层级、克制使用品牌色、注重效率

### 统一元素
- **品牌色**: #7ACC35 (KOOK 绿)
- **字体**: 系统默认字体栈
- **基础间距**: 4px 基础单位
- **圆角**: 统一的圆角系统

---

## 1. 色彩系统

### 1.1 品牌色

```scss
// 主品牌色 - KOOK 绿
--color-primary: #7ACC35;
--color-primary-light: #87D149;
--color-primary-dark: #6DB72F;

// RGB (用于透明度)
--color-primary-rgb: 122, 204, 53;

// 渐变
--color-primary-gradient: linear-gradient(135deg, #7ACC35 0%, #6DB72F 100%);

// 光晕效果
--color-primary-glow: 0 0 20rpx rgba(122, 204, 53, 0.4);
```

### 1.2 应用场景

| 场景 | 移动端 | 管理后台 |
|------|--------|----------|
| 主按钮 | ✅ 品牌色渐变 | ✅ 品牌色纯色 |
| 强调元素 | ✅ 大量使用 | ⚠️ 克制使用 |
| 图标激活 | ✅ 品牌色 | ✅ 品牌色 |
| 链接 | ✅ 品牌色 | ✅ 品牌色 |
| 标签/徽章 | ✅ 根据状态 | ✅ 根据状态 |

### 1.3 语义色彩

```scss
// 成功
--color-success: #10B981;
--color-success-light: #34D399;

// 警告
--color-warning: #F59E0B;
--color-warning-light: #FBBF24;

// 错误
--color-error: #EF4444;
--color-error-light: #F87171;

// 信息
--color-info: #3B82F6;
--color-info-light: #60A5FA;
```

### 1.4 中性色彩

```scss
// 亮色模式
--color-bg: #F8FAFC;
--color-bg-card: #FFFFFF;
--color-bg-secondary: #F1F5F9;

--color-text: #0F172A;
--color-text-secondary: #64748B;
--color-text-placeholder: #94A3B8;

--color-border: #E2E8F0;

// 暗色模式
--color-bg-dark: #151618;
--color-bg-card-dark: #1C1D20;
--color-bg-secondary-dark: #1A1B1E;

--color-text-dark: #F1F5F9;
--color-text-secondary-dark: #94A3B8;

--color-border-dark: #2A2B30;
```

---

## 2. 间距系统

### 2.1 基础单位

```scss
// 移动端 (rpx)
$spacing-xs: 12rpx;   // 6px  - 最小间距
$spacing-sm: 20rpx;   // 10px - 小间距
$spacing-md: 28rpx;   // 14px - 中等间距
$spacing-lg: 36rpx;   // 18px - 大间距
$spacing-xl: 52rpx;   // 26px - 超大间距

// 管理后台 (px)
--spacing-xs: 6px;    // 最小间距
--spacing-sm: 8px;    // 小间距
--spacing-md: 12px;   // 中等间距
--spacing-lg: 16px;   // 大间距
--spacing-xl: 24px;   // 超大间距
--spacing-2xl: 32px;  // 特大间距
```

### 2.2 使用指南

| 场景 | 移动端 | 管理后台 |
|------|--------|----------|
| 卡片内边距 | md (28rpx) | md (12px) |
| 组件间距 | sm-md (20-28rpx) | md (12px) |
| section 间距 | lg-xl (36-52rpx) | lg-xl (16-24px) |
| 表单元素间距 | sm (20rpx) | sm (8px) |

---

## 3. 字体系统

### 3.1 字号层级

```scss
// 移动端 (rpx)
$font-xs: 24rpx;   // 12px - 辅助信息
$font-sm: 28rpx;   // 14px - 次要内容
$font-md: 32rpx;   // 16px - 正文
$font-base: 36rpx; // 18px - 基础
$font-lg: 40rpx;   // 20px - 小标题
$font-xl: 48rpx;   // 24px - 标题
$font-2xl: 56rpx;  // 28px - 大标题
$font-3xl: 64rpx;  // 32px - 特大标题

// 管理后台 (px)
--font-xs: 12px;   // 辅助信息
--font-sm: 13px;   // 次要内容
--font-base: 14px; // 正文
--font-md: 16px;   // 强调正文
--font-lg: 18px;   // 小标题
--font-xl: 20px;   // 标题
```

### 3.2 字重

```scss
--font-weight-normal: 400;
--font-weight-medium: 500;
--font-weight-semibold: 600;
--font-weight-bold: 700;
```

### 3.3 使用指南

| 场景 | 字号 | 字重 |
|------|------|------|
| 页面标题 | xl-2xl | Semibold/Bold |
| 卡片标题 | md-lg | Semibold |
| 正文 | md/base | Normal |
| 次要文字 | sm | Normal |
| 辅助信息 | xs | Normal |

---

## 4. 圆角系统

### 4.1 圆角定义

```scss
// 移动端 (rpx)
$radius-sm: 8rpx;   // 小圆角 - 标签、徽章
$radius-md: 12rpx;  // 中圆角 - 按钮、输入框
$radius-lg: 24rpx;  // 大圆角 - 卡片
$radius-xl: 32rpx;  // 超大圆角 - 模态框
$radius-full: 9999rpx; // 完整圆形 - 头像、按钮

// 管理后台 (px)
--radius-sm: 4px;   // 小圆角
--radius-md: 6px;   // 中圆角
--radius-lg: 8px;   // 大圆角
--radius-xl: 12px;  // 超大圆角
--radius-full: 9999px; // 完整圆形
```

### 4.2 使用指南

| 组件 | 移动端 | 管理后台 |
|------|--------|----------|
| 按钮 | md/full | md |
| 输入框 | md | md |
| 卡片 | lg | lg |
| 标签 | sm/full | sm |
| 头像 | full | full |
| 模态框 | xl | lg |

---

## 5. 阴影系统

### 5.1 阴影定义

```scss
// 移动端 (rpx)
$shadow-sm: 0 2rpx 6rpx rgba(0, 0, 0, 0.06);
$shadow-md: 0 4rpx 12rpx rgba(0, 0, 0, 0.08);
$shadow-lg: 0 8rpx 24rpx rgba(0, 0, 0, 0.12);
$shadow-xl: 0 16rpx 48rpx rgba(0, 0, 0, 0.2);

// 管理后台 (px)
--shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05);
--shadow-md: 0 2px 8px rgba(0, 0, 0, 0.08);
--shadow-lg: 0 4px 16px rgba(0, 0, 0, 0.12);
--shadow-xl: 0 8px 32px rgba(0, 0, 0, 0.16);
```

### 5.2 使用指南

| 场景 | 移动端 | 管理后台 |
|------|--------|----------|
| 卡片默认 | md | md |
| 卡片悬浮 | lg | lg |
| 模态框 | xl | xl |
| 下拉菜单 | md | lg |

---

## 6. 动画系统

### 6.1 缓动函数

```scss
--transition-fast: 150ms cubic-bezier(0.4, 0, 0.2, 1);
--transition-normal: 250ms cubic-bezier(0.4, 0, 0.2, 1);
--transition-slow: 350ms cubic-bezier(0.4, 0, 0.2, 1);
--transition-spring: 400ms cubic-bezier(0.34, 1.56, 0.64, 1);
```

### 6.2 使用指南

| 交互 | 时长 | 缓动 |
|------|------|------|
| 按钮点击 | fast | ease-out |
| 卡片悬浮 | normal | ease-out |
| 页面切换 | slow | ease-in-out |
| 弹性效果 | spring | spring |

---

## 7. 组件规范

### 7.1 按钮

| 类型 | 移动端 | 管理后台 |
|------|--------|----------|
| 主按钮 | 品牌色渐变 + 大圆角 | 品牌色纯色 + 中圆角 |
| 次要按钮 | 描边 + 中圆角 | 描边 + 小圆角 |
| 文本按钮 | 品牌色文字 | 品牌色文字 |
| 危险按钮 | 错误色 + 中圆角 | 错误色 + 小圆角 |

### 7.2 卡片

| 属性 | 移动端 | 管理后台 |
|------|--------|----------|
| 内边距 | md (28rpx) | md (12px) |
| 圆角 | lg (24rpx) | lg (8px) |
| 阴影 | md | md |
| 边框 | 无 | 细边框 |

### 7.3 输入框

| 属性 | 移动端 | 管理后台 |
|------|--------|----------|
| 高度 | 80-88rpx | 32-36px |
| 内边距 | md (28rpx) | sm-md (8-12px) |
| 圆角 | md (12rpx) | md (6px) |
| 边框 | 细边框 | 细边框 |

---

## 8. 响应式断点

### 8.1 移动端

```scss
// 小程序/UniApp 默认响应式
// 无需额外断点
```

### 8.2 管理后台

```scss
// 平板
@media (min-width: 768px) { }

// 桌面
@media (min-width: 1024px) { }

// 大屏
@media (min-width: 1440px) { }
```

---

## 9. 暗色模式

### 9.1 切换原则

- **移动端**: 自动跟随系统
- **管理后台**: 用户手动切换

### 9.2 色彩映射

| 元素 | 亮色模式 | 暗色模式 |
|------|----------|----------|
| 主背景 | #F8FAFC | #151618 |
| 卡片背景 | #FFFFFF | #1C1D20 |
| 次级背景 | #F1F5F9 | #1A1B1E |
| 主文字 | #0F172A | #F1F5F9 |
| 次级文字 | #64748B | #94A3B8 |
| 边框 | #E2E8F0 | #2A2B30 |

---

## 10. 品牌元素使用

### 10.1 Logo

- **移动端**: 简化版 Logo
- **管理后台**: 完整版 Logo

### 10.2 品牌色使用原则

✅ **推荐使用**:
- 主要操作按钮
- 激活状态
- 重要信息强调
- 数据可视化

⚠️ **克制使用**:
- 大面积背景（管理后台）
- 文字颜色（除非强调）
- 装饰性元素

❌ **避免使用**:
- 错误/警告状态
- 大量重复使用
- 与其他主色混用

---

## 11. 特殊元素

### 11.1 游戏元素（仅移动端）

```scss
// 等级/稀有度配色
$rank-bronze: #CD7F32;
$rank-silver: #C0C0C0;
$rank-gold: #FFD700;
$rank-platinum: #E5E4E2;
$rank-diamond: #B9F2FF;
```

### 11.2 VIP 系统

```scss
// VIP 等级渐变
$vip-1: linear-gradient(135deg, #FFD700 0%, #FFA500 100%);
$vip-2: linear-gradient(135deg, #FF6B6B 0%, #EE5A24 100%);
$vip-3: linear-gradient(135deg, #A78BFA 0%, #7C3AED 100%);
$vip-4: linear-gradient(135deg, #22D3EE 0%, #0EA5E9 100%);
```

---

## 12. 可访问性

### 12.1 对比度要求

- **正文文字**: 最小 4.5:1
- **大文字 (18px+)**: 最小 3:1
- **UI 组件**: 最小 3:1

### 12.2 触摸目标

- **移动端**: 最小 44×44px (88×88rpx)
- **管理后台**: 最小 32×32px

---

## 13. 文件组织

### 13.1 移动端

```
app/src/styles/
├── variables.scss       # 设计令牌
├── mixins.scss          # 通用 mixins
├── index.scss           # 全局样式
└── themes/
    ├── light.scss       # 亮色主题
    └── dark.scss        # 暗色主题
```

### 13.2 管理后台

```
admin/src/styles/
├── variables.css        # 设计令牌
├── global.css           # 全局样式
└── themes/
    ├── light.css        # 亮色主题
    └── dark.css         # 暗色主题
```

---

## 14. 版本历史

| 版本 | 日期 | 变更 |
|------|------|------|
| 1.0.0 | 2026-02-09 | 初始版本，统一移动端和管理后台设计规范 |

---

## 15. 维护者

- **Mobile-Lead**: 移动端设计系统
- **Frontend-Lead**: 管理后台设计系统

---

**最后更新**: 2026-02-09
