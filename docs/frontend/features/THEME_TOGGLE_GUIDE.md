# 主题切换功能指南

**实现日期**: 2025-10-28  
**设计风格**: Neo-brutalism 黑白主题切换  
**状态**: ✅ 已完成

---

## 🎨 功能概述

在全局导航栏（Header）中添加了主题切换按钮，支持浅色和深色两种模式，保持 Neo-brutalism 黑白极简设计风格。

---

## 🚀 快速使用

### 位置

主题切换按钮位于顶部导航栏右侧：

```
┌─────────────────────────────────────┐
│ ≡  GameLink          🌙  👤 admin  │
│                       ↑              │
│                  主题按钮             │
└─────────────────────────────────────┘
```

### 操作

1. **浅色模式** → 显示月亮图标 🌙
2. **深色模式** → 显示太阳图标 ☀️
3. **点击切换** → 在两种模式之间切换

---

## 🎨 视觉效果

### 浅色模式（Light Mode）

```
背景：纯白色 #ffffff
文字：纯黑色 #000000
边框：黑色 #000000
图标：月亮 🌙
```

**特点**：

- 高对比度
- 适合白天使用
- 清晰明了

### 深色模式（Dark Mode）

```
背景：深黑色 #0a0a0a
文字：纯白色 #ffffff
边框：白色 #ffffff
图标：太阳 ☀️
```

**特点**：

- 护眼舒适
- 适合夜间使用
- 节省电量（OLED 屏幕）

---

## 🔧 技术实现

### 1. 组件更新

#### Header 组件 (`src/components/Layout/Header.tsx`)

```tsx
import { useTheme } from 'contexts/ThemeContext';

export const Header = ({ user, onLogout, onToggleSidebar }) => {
  const { effective, setMode } = useTheme();

  const toggleTheme = () => {
    setMode(effective === 'light' ? 'dark' : 'light');
  };

  return (
    <header>
      {/* 主题切换按钮 */}
      <button onClick={toggleTheme}>{effective === 'light' ? <MoonIcon /> : <SunIcon />}</button>
    </header>
  );
};
```

### 2. CSS 变量系统

#### 浅色模式（默认）

```less
:root {
  --text-primary: #000000;
  --text-secondary: #666666;
  --bg-primary: #ffffff;
  --bg-secondary: #fafafa;
  --border-color: #000000;
}
```

#### 深色模式

```less
body.dark-theme {
  --text-primary: #ffffff;
  --text-secondary: #cccccc;
  --bg-primary: #0a0a0a;
  --bg-secondary: #1a1a1a;
  --border-color: #ffffff;
}
```

### 3. ThemeContext

**文件**: `src/contexts/ThemeContext.tsx`

**功能**:

- ✅ 管理主题状态（light / dark）
- ✅ 保存到 localStorage
- ✅ 自动应用主题类名
- ✅ 支持系统主题检测

**API**:

```typescript
const { mode, effective, setMode } = useTheme();

// mode: 'light' | 'dark' | 'system'
// effective: 'light' | 'dark'
// setMode: (mode) => void
```

---

## 🎯 按钮样式

### 设计特点

1. **尺寸**: 40px × 40px
2. **边框**: 2px 黑色/白色边框
3. **背景**: 透明
4. **图标**: 20px × 20px

### 交互动效

**悬停效果**:

- 阴影增强
- 上浮动画（-1px, -1px）
- 图标旋转 20°
- 水波纹扩散

**点击效果**:

- 立即切换主题
- 保存到 localStorage
- 全局样式更新

```less
.themeButton {
  width: 40px;
  height: 40px;
  border: 2px solid var(--color-black);

  &:hover {
    box-shadow: 2px 2px 0 var(--color-black);
    transform: translate(-1px, -1px);

    svg {
      transform: rotate(20deg);
    }
  }
}
```

---

## 💾 数据持久化

### LocalStorage

主题选择会自动保存：

```javascript
// 存储键
localStorage.setItem('gamelink_theme', 'light'); // 或 'dark'

// 查看方式
// 开发者工具 → Application → Local Storage
// 键: gamelink_theme
// 值: light 或 dark
```

### 自动恢复

刷新页面后，系统会自动读取并应用保存的主题：

```typescript
useEffect(() => {
  const savedTheme = localStorage.getItem('gamelink_theme');
  if (savedTheme) {
    applyTheme(savedTheme);
  }
}, []);
```

---

## 🎨 深色模式适配

### 所有组件自动适配

由于使用 CSS 变量系统，所有组件会自动适配深色模式：

| 组件        | 浅色模式 | 深色模式 |
| ----------- | -------- | -------- |
| **Header**  | 白底黑字 | 黑底白字 |
| **Sidebar** | 白底黑字 | 黑底白字 |
| **Card**    | 白底黑框 | 黑底白框 |
| **Button**  | 黑底白字 | 白底黑字 |
| **Input**   | 白底黑框 | 黑底白框 |

### 示例：Card 组件

```less
.card {
  background-color: var(--bg-primary); // 自动切换
  color: var(--text-primary); // 自动切换
  border: 2px solid var(--border-color); // 自动切换
}

// 浅色模式：白底 + 黑字 + 黑框
// 深色模式：黑底 + 白字 + 白框
```

---

## 🌓 图标说明

### 月亮图标（浅色模式显示）

```
🌙 MoonIcon
```

含义：点击切换到**深色模式**

### 太阳图标（深色模式显示）

```
☀️ SunIcon
```

含义：点击切换到**浅色模式**

### 设计理念

- **当前是什么模式，就显示切换到另一个模式的图标**
- 浅色模式 → 显示月亮（切换到深色）
- 深色模式 → 显示太阳（切换到浅色）

---

## 📱 响应式支持

### 桌面端

```
┌────────────────────────────┐
│ ≡  GameLink   🌙  👤 admin │
└────────────────────────────┘
```

### 移动端

```
┌──────────────┐
│ ≡  GameLink  │
│    🌙  👤    │  ← 图标和用户名紧凑排列
└──────────────┘
```

---

## 🔍 使用场景

### 1. 日常使用

**白天**:

- 使用浅色模式
- 高对比度，清晰明了
- 适合明亮环境

**夜晚**:

- 使用深色模式
- 降低屏幕亮度
- 减少眼睛疲劳

### 2. 不同环境

**办公室**:

- 浅色模式
- 专业感强
- 适合演示

**家里/咖啡厅**:

- 深色模式
- 氛围柔和
- 省电护眼

---

## 🎯 最佳实践

### 用户体验

```typescript
// ✅ 好的实践
// 1. 记住用户选择
localStorage.setItem('gamelink_theme', theme);

// 2. 平滑过渡
transition: all 200ms ease-out;

// 3. 清晰的图标
浅色 → 🌙 月亮
深色 → ☀️ 太阳
```

### 开发建议

```typescript
// ✅ 使用 CSS 变量
color: var(--text-primary);
background: var(--bg-primary);

// ❌ 避免硬编码颜色
color: #000000;  // 不会自动切换
background: #ffffff;  // 不会自动切换
```

---

## 🐛 故障排除

### 问题1：主题切换不生效

**检查**:

1. 是否使用了 CSS 变量？
2. 是否正确导入了 ThemeProvider？
3. 浏览器是否支持 CSS 变量？

**解决方案**:

```tsx
// main.tsx
<ThemeProvider>
  <AuthProvider>
    <App />
  </AuthProvider>
</ThemeProvider>
```

### 问题2：刷新页面后主题重置

**原因**: localStorage 读取失败

**解决方案**:

- 检查浏览器是否禁用 localStorage
- 检查是否在隐私模式下

### 问题3：部分组件没有适配深色模式

**原因**: 使用了硬编码颜色

**解决方案**:

```less
// ❌ 错误
.component {
  color: #000000;
  background: #ffffff;
}

// ✅ 正确
.component {
  color: var(--text-primary);
  background: var(--bg-primary);
}
```

---

## 📊 对比表

| 特性         | 浅色模式      | 深色模式      |
| ------------ | ------------- | ------------- |
| **背景色**   | #ffffff（白） | #0a0a0a（黑） |
| **文字色**   | #000000（黑） | #ffffff（白） |
| **边框色**   | #000000（黑） | #ffffff（白） |
| **图标**     | 🌙 月亮       | ☀️ 太阳       |
| **适用场景** | 白天/明亮环境 | 夜晚/暗光环境 |
| **视觉特点** | 高对比/清晰   | 柔和/护眼     |
| **能耗**     | 标准          | 省电（OLED）  |

---

## 🔄 主题切换流程

```
1. 用户点击主题按钮
      ↓
2. toggleTheme() 调用
      ↓
3. setMode('light' | 'dark')
      ↓
4. ThemeContext 更新状态
      ↓
5. 添加/移除 .dark-theme 类
      ↓
6. CSS 变量自动切换
      ↓
7. 所有组件重新渲染（使用新颜色）
      ↓
8. 保存到 localStorage
      ↓
9. 图标更新（月亮 ↔ 太阳）
```

---

## 💡 扩展功能

### 可以添加的功能

1. **自动主题**

   ```typescript
   // 跟随系统主题
   setMode('system');
   ```

2. **定时切换**

   ```typescript
   // 白天自动浅色，晚上自动深色
   const hour = new Date().getHours();
   setMode(hour >= 6 && hour < 18 ? 'light' : 'dark');
   ```

3. **多种配色**

   ```typescript
   // 不只是黑白，还可以有其他配色
   setTheme('blue-dark');
   setTheme('green-light');
   ```

4. **主题动画**
   ```css
   /* 切换时的过渡动画 */
   * {
     transition:
       background-color 0.3s,
       color 0.3s;
   }
   ```

---

## 📚 相关文档

- [设计系统文档](./DESIGN_SYSTEM_V2.md)
- [导航系统文档](./NAVIGATION_SYSTEM.md)
- [CSS 变量系统](./src/styles/variables.less)

---

## ✅ 功能清单

- [x] Header 中添加主题切换按钮
- [x] 浅色/深色模式切换
- [x] 图标切换动画
- [x] CSS 变量系统
- [x] 深色模式样式
- [x] LocalStorage 持久化
- [x] 自动恢复主题
- [x] 所有组件适配
- [x] 响应式支持
- [x] 文档编写

---

**实现者**: GameLink Frontend Team  
**设计风格**: Neo-brutalism 黑白主题  
**最后更新**: 2025-10-28

---

<div align="center">

## 🌓 浅色与深色，随心切换！

**Light ⚪ ↔ ⚫ Dark**

⚫⚪ **极简 · 优雅 · 护眼**

</div>
