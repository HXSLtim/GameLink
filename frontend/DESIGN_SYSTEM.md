# GameLink Frontend 设计系统

**版本**: 1.0.0  
**更新日期**: 2025-10-27  
**设计语言**: 现代化、专业、游戏化

---

## 🎨 设计理念

GameLink 管理系统采用现代化的设计语言，结合游戏行业特点，打造专业、高效、美观的管理界面。

### 核心设计原则

1. **简洁明了** - 清晰的信息层级，减少认知负担
2. **视觉愉悦** - 渐变色彩、流畅动画、精致细节
3. **响应式设计** - 适配各种设备和屏幕尺寸
4. **一致性** - 统一的交互模式和视觉语言
5. **可访问性** - 良好的对比度和可读性

---

## 🌈 色彩系统

### 主色调（Primary Colors）

```less
// 主渐变色
--gradient-primary: linear-gradient(135deg, #667eea 0%, #764ba2 100%);

// 主色
--color-primary: #667eea;
--color-primary-light: #8b9ff5;
--color-primary-dark: #5568d3;

// 次要色
--color-secondary: #764ba2;
--color-secondary-light: #9470b8;
--color-secondary-dark: #5e3c82;
```

### 功能色（Functional Colors）

```less
// 成功
--color-success: #00d084;
--color-success-light: #33da9f;
--color-success-dark: #00b56f;

// 警告
--color-warning: #ff9800;
--color-warning-light: #ffa726;
--color-warning-dark: #f57c00;

// 错误
--color-error: #f5222d;
--color-error-light: #ff4d4f;
--color-error-dark: #cf1322;

// 信息
--color-info: #1890ff;
--color-info-light: #40a9ff;
--color-info-dark: #096dd9;
```

### 中性色（Neutral Colors）

```less
// 文本颜色
--text-primary: #1f2937;
--text-secondary: #6b7280;
--text-tertiary: #9ca3af;
--text-placeholder: #d1d5db;

// 背景颜色
--bg-white: #ffffff;
--bg-gray-50: #f9fafb;
--bg-gray-100: #f3f4f6;
--bg-gray-200: #e5e7eb;

// 边框颜色
--border-light: #f0f0f0;
--border-normal: #d9d9d9;
--border-dark: #bfbfbf;
```

### 渐变装饰色

```less
// 装饰渐变球
--gradient-orb-1: radial-gradient(circle, rgba(255, 107, 107, 0.8) 0%, transparent 70%);
--gradient-orb-2: radial-gradient(circle, rgba(78, 205, 196, 0.8) 0%, transparent 70%);
--gradient-orb-3: radial-gradient(circle, rgba(255, 195, 113, 0.8) 0%, transparent 70%);
```

---

## 📐 布局系统

### 间距规范（Spacing Scale）

```less
--spacing-1: 4px;   // 0.25rem
--spacing-2: 8px;   // 0.5rem
--spacing-3: 12px;  // 0.75rem
--spacing-4: 16px;  // 1rem
--spacing-5: 20px;  // 1.25rem
--spacing-6: 24px;  // 1.5rem
--spacing-8: 32px;  // 2rem
--spacing-10: 40px; // 2.5rem
--spacing-12: 48px; // 3rem
--spacing-16: 64px; // 4rem
```

### 圆角规范（Border Radius）

```less
--radius-small: 4px;
--radius-medium: 8px;
--radius-large: 12px;
--radius-xlarge: 16px;
--radius-2xlarge: 20px;
--radius-round: 50%;
```

### 阴影规范（Shadows）

```less
// 卡片阴影
--shadow-sm: 0 2px 8px rgba(0, 0, 0, 0.08);
--shadow-md: 0 4px 16px rgba(0, 0, 0, 0.1);
--shadow-lg: 0 8px 24px rgba(0, 0, 0, 0.12);
--shadow-xl: 0 20px 60px rgba(0, 0, 0, 0.2);

// 悬浮阴影
--shadow-hover: 0 12px 30px rgba(102, 126, 234, 0.4);

// 按钮阴影
--shadow-button: 0 8px 20px rgba(102, 126, 234, 0.3);
```

### 容器宽度

```less
--container-xs: 480px;
--container-sm: 640px;
--container-md: 768px;
--container-lg: 1024px;
--container-xl: 1280px;
--container-2xl: 1536px;
```

---

## ✏️ 字体系统

### 字体家族

```less
--font-family-base: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 
                    'Helvetica Neue', Arial, sans-serif;
--font-family-code: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', 
                    'Courier New', monospace;
```

### 字体大小

```less
--font-size-xs: 12px;   // 0.75rem
--font-size-sm: 13px;   // 0.8125rem
--font-size-base: 14px; // 0.875rem
--font-size-lg: 16px;   // 1rem
--font-size-xl: 18px;   // 1.125rem
--font-size-2xl: 20px;  // 1.25rem
--font-size-3xl: 24px;  // 1.5rem
--font-size-4xl: 30px;  // 1.875rem
--font-size-5xl: 36px;  // 2.25rem
```

### 字重（Font Weight）

```less
--font-weight-light: 300;
--font-weight-normal: 400;
--font-weight-medium: 500;
--font-weight-semibold: 600;
--font-weight-bold: 700;
```

### 行高（Line Height）

```less
--line-height-tight: 1.25;
--line-height-normal: 1.5;
--line-height-relaxed: 1.75;
--line-height-loose: 2;
```

---

## 🎭 组件设计规范

### 1. 登录页面（Login Page）

#### 设计特点
- ✨ 全屏渐变背景
- 🎨 动态装饰球体
- 💳 毛玻璃卡片效果
- 🎯 居中对齐布局
- ⚡ 流畅的动画效果

#### 组件层级

```
LoginPage
├── Background (背景层)
│   ├── GradientOrb1 (装饰球体1)
│   ├── GradientOrb2 (装饰球体2)
│   └── GradientOrb3 (装饰球体3)
├── Container (主容器)
│   ├── Header (头部)
│   │   ├── LogoContainer (Logo容器)
│   │   ├── Title (标题)
│   │   └── Subtitle (副标题)
│   ├── Card (登录卡片)
│   │   └── Form (表单)
│   │       ├── UsernameInput (用户名输入)
│   │       ├── PasswordInput (密码输入)
│   │       ├── FormOptions (记住我 + 忘记密码)
│   │       ├── LoginButton (登录按钮)
│   │       └── DevInfo (开发环境提示)
│   └── Footer (页脚)
```

#### 视觉规范

**背景渐变**
```less
background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
```

**卡片样式**
```less
background: rgba(255, 255, 255, 0.95);
backdrop-filter: blur(10px);
border-radius: 20px;
box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
```

**Logo 容器**
```less
width: 80px;
height: 80px;
background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
border-radius: 20px;
box-shadow: 0 10px 40px rgba(102, 126, 234, 0.4);
```

**登录按钮**
```less
height: 44px;
font-size: 16px;
font-weight: 600;
border-radius: 10px;
background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
box-shadow: 0 8px 20px rgba(102, 126, 234, 0.3);
```

#### 动画效果

**页面进入动画**
```less
@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
duration: 0.6s
easing: ease-out
```

**Logo 弹跳动画**
```less
@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}
duration: 2s
easing: ease-in-out
iteration: infinite
```

**装饰球浮动动画**
```less
@keyframes float {
  0%, 100% { transform: translate(0, 0) scale(1); }
  33% { transform: translate(30px, -30px) scale(1.1); }
  66% { transform: translate(-20px, 20px) scale(0.9); }
}
duration: 20s
easing: ease-in-out
iteration: infinite
```

**卡片悬浮效果**
```less
transition: transform 0.3s ease, box-shadow 0.3s ease;
&:hover {
  transform: translateY(-5px);
  box-shadow: 0 25px 70px rgba(0, 0, 0, 0.25);
}
```

#### 响应式设计

**平板设备 (≤768px)**
```less
- Logo: 64px × 64px
- Card padding: 24px
- 装饰球体: 300px × 300px
```

**移动设备 (≤480px)**
```less
- Card padding: 20px
- 表单选项垂直排列
- 最大宽度: 100%
```

#### 交互状态

**输入框状态**
- Normal: 默认边框
- Focus: 主色边框 + 阴影
- Error: 红色边框 + 错误提示
- Disabled: 灰色背景 + 禁用光标

**按钮状态**
- Normal: 渐变背景 + 阴影
- Hover: 向上移动 2px + 加深阴影
- Active: 恢复原位
- Loading: 显示加载图标 + 禁用点击
- Disabled: 灰色背景 + 禁用光标

---

## 🎬 动画规范

### 动画时长（Duration）

```less
--duration-fast: 150ms;    // 快速交互
--duration-normal: 300ms;  // 标准动画
--duration-slow: 500ms;    // 缓慢过渡
--duration-slower: 800ms;  // 页面级动画
```

### 缓动函数（Easing）

```less
--ease-in: cubic-bezier(0.4, 0, 1, 1);
--ease-out: cubic-bezier(0, 0, 0.2, 1);
--ease-in-out: cubic-bezier(0.4, 0, 0.2, 1);
--ease-bounce: cubic-bezier(0.68, -0.55, 0.265, 1.55);
```

### 动画使用原则

1. **微交互优先** - 优先为用户交互添加反馈动画
2. **性能优化** - 使用 `transform` 和 `opacity` 而非 `top`/`left`
3. **适度使用** - 避免过度动画导致眩晕
4. **可关闭** - 尊重用户的 `prefers-reduced-motion` 设置

---

## 📱 响应式断点

### 断点定义

```less
// 超小屏幕（手机竖屏）
@screen-xs: 480px;

// 小屏幕（手机横屏）
@screen-sm: 640px;

// 中等屏幕（平板竖屏）
@screen-md: 768px;

// 大屏幕（平板横屏、笔记本）
@screen-lg: 1024px;

// 超大屏幕（桌面显示器）
@screen-xl: 1280px;

// 超超大屏幕（大显示器）
@screen-2xl: 1536px;
```

### 使用示例

```less
// 移动优先
.component {
  padding: 16px;
  
  @media (min-width: @screen-md) {
    padding: 24px;
  }
  
  @media (min-width: @screen-lg) {
    padding: 32px;
  }
}
```

---

## ♿ 可访问性规范

### 颜色对比度

- 普通文本: 最小对比度 4.5:1
- 大文本 (18px+): 最小对比度 3:1
- UI 组件: 最小对比度 3:1

### 焦点指示器

```less
&:focus-visible {
  outline: 2px solid var(--color-primary);
  outline-offset: 2px;
}
```

### 键盘导航

- 所有交互元素支持 Tab 键导航
- 支持 Enter/Space 触发按钮
- 支持 Escape 关闭模态框

### ARIA 标签

```tsx
<button aria-label="登录" aria-busy={loading}>
  {loading ? '登录中...' : '登录'}
</button>
```

---

## 🔧 设计工具和资源

### 使用的设计系统
- **Arco Design**: 基础组件库
- **CSS Modules**: 样式隔离
- **LESS**: CSS 预处理器

### 图标系统
- **@arco-design/web-react/icon**: 官方图标库
- 自定义图标: SVG 格式，24×24px 基准

### 字体资源
- **系统字体栈**: 保证各平台最佳显示效果
- **中文字体**: 苹方、微软雅黑
- **英文字体**: San Francisco、Segoe UI、Roboto

---

## 📋 设计检查清单

### 新组件设计检查

- [ ] 符合色彩系统规范
- [ ] 遵循间距和布局规范
- [ ] 使用统一的圆角和阴影
- [ ] 实现响应式布局
- [ ] 添加合适的动画效果
- [ ] 支持深色模式（可选）
- [ ] 满足可访问性要求
- [ ] 有 Loading 和 Error 状态
- [ ] 支持键盘操作
- [ ] 测试移动端体验

---

## 🎯 后续设计计划

### 待设计页面

1. **仪表盘（Dashboard）**
   - 数据可视化卡片
   - 实时统计图表
   - 快捷操作入口

2. **用户管理（Users）**
   - 用户列表表格
   - 筛选和搜索
   - 用户详情抽屉

3. **订单管理（Orders）**
   - 订单状态流程
   - 订单详情页
   - 数据导出功能

4. **权限管理（Permissions）**
   - 角色权限矩阵
   - 权限树形选择器
   - 分配权限界面

5. **设置页面（Settings）**
   - 系统配置表单
   - 个人信息设置
   - 主题切换

---

## 📚 参考资源

- [Arco Design 官方文档](https://arco.design/)
- [Material Design Guidelines](https://material.io/design)
- [Apple Human Interface Guidelines](https://developer.apple.com/design/)
- [Web Content Accessibility Guidelines (WCAG)](https://www.w3.org/WAI/WCAG21/quickref/)

---

**维护者**: GameLink Team  
**最后更新**: 2025-10-27




