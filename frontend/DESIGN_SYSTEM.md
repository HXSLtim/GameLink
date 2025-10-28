# GameLink Frontend 设计系统

**版本**: 2.0.0  
**更新日期**: 2025-10-27  
**设计语言**: 基于 Arco Design，融合 GameLink 品牌特色  
**UI 框架**: [Arco Design](https://arco.design/) by 字节跳动

---

## 📖 关于本设计系统

GameLink 管理系统基于 **Arco Design** 企业级设计语言构建，遵循清晰、一致、开放与韵律感的设计价值观，同时融入游戏行业特色，打造专业、高效、美观的管理界面。

### 设计价值观

> **务实与浪漫并重** - 既强调清晰效率与品牌一致性，又追求设计韵律与开放包容

1. **清晰明了** - 清晰的信息层级，减少认知负担
2. **视觉愉悦** - 现代化配色、流畅动画、精致细节
3. **一致性** - 统一的交互模式和视觉语言
4. **可访问性** - 良好的对比度和可读性
5. **响应式** - 适配各种设备和屏幕尺寸

---

## 🎨 色彩系统

### Arco Design 官方色卡

基于 Arco Design 色彩体系，支持亮/暗模式无缝切换。

#### 主品牌色（Primary）

```less
// Arco Blue - 主品牌色
--arcoblue-1: #e8f3ff;
--arcoblue-2: #bedaff;
--arcoblue-3: #94bfff;
--arcoblue-4: #6aa1ff;
--arcoblue-5: #4080ff;
--arcoblue-6: #165dff; // 主色 - 用于主按钮、链接、高亮
--arcoblue-7: #0e42d2;
--arcoblue-8: #072ca6;
--arcoblue-9: #031a79;
--arcoblue-10: #01114d;
```

**使用场景**：

- 主按钮（Primary Button）
- 链接文字（Link）
- 选中状态（Selected State）
- 进度指示（Progress）

#### 功能色（Functional Colors）

**成功色（Success）**

```less
--green-1: #e8ffea;
--green-2: #aff0b5;
--green-3: #7be188;
--green-4: #4cd263;
--green-5: #23c343;
--green-6: #00b42a; // 主成功色
--green-7: #009a29;
--green-8: #008026;
--green-9: #006622;
--green-10: #004d1c;
```

**警告色（Warning）**

```less
--orange-1: #fff7e8;
--orange-2: #ffe4ba;
--orange-3: #ffcf8b;
--orange-4: #ffb65d;
--orange-5: #ff9a2e;
--orange-6: #ff7d00; // 主警告色
--orange-7: #d25f00;
--orange-8: #a64500;
--orange-9: #792e00;
--orange-10: #4d1b00;
```

**错误色（Error）**

```less
--red-1: #ffece8;
--red-2: #fdcdc5;
--red-3: #fbaca3;
--red-4: #f98981;
--red-5: #f76560;
--red-6: #f53f3f; // 主错误色
--red-7: #cb272d;
--red-8: #a1151e;
--red-9: #770813;
--red-10: #4d000a;
```

**信息色（Info）**

```less
--cyan-1: #e8fffb;
--cyan-2: #b7f4ec;
--cyan-3: #87e9d9;
--cyan-4: #58d9c3;
--cyan-5: #2cc9aa;
--cyan-6: #14c9c9; // 主信息色
--cyan-7: #0da5aa;
--cyan-8: #07828b;
--cyan-9: #03616c;
--cyan-10: #00424d;
```

#### 中性色（Neutral Colors）

**亮色模式**

```less
// 文本颜色
--color-text-1: #1d2129; // 标题、重要文本
--color-text-2: #4e5969; // 次要文本
--color-text-3: #86909c; // 占位符、说明文字
--color-text-4: #c9cdd4; // 禁用状态

// 背景颜色
--color-bg-1: #ffffff; // 主背景
--color-bg-2: #f7f8fa; // 次要背景
--color-bg-3: #f2f3f5; // 三级背景
--color-bg-4: #e5e6eb; // 分割线背景
--color-bg-5: #c9cdd4; // 边框

// 填充颜色
--color-fill-1: #f7f8fa;
--color-fill-2: #f2f3f5;
--color-fill-3: #e5e6eb;
--color-fill-4: #c9cdd4;
```

**暗色模式**

```less
// 文本颜色
--color-text-1: #e5e6eb;
--color-text-2: #c9cdd4;
--color-text-3: #86909c;
--color-text-4: #4e5969;

// 背景颜色
--color-bg-1: #17171a;
--color-bg-2: #232324;
--color-bg-3: #2a2a2b;
--color-bg-4: #313132;
--color-bg-5: #3c3c3d;
```

### GameLink 品牌渐变色

为突出游戏行业特色，在 Arco 基础上增加品牌渐变：

```less
// GameLink 主渐变（用于登录页、品牌展示）
--gradient-gamelink-primary: linear-gradient(135deg, #667eea 0%, #764ba2 100%);

// 装饰渐变
--gradient-gamelink-accent: linear-gradient(120deg, #f093fb 0%, #f5576c 100%);
--gradient-gamelink-success: linear-gradient(120deg, #84fab0 0%, #8fd3f4 100%);
--gradient-gamelink-warning: linear-gradient(120deg, #ffc837 0%, #ff8008 100%);

// 装饰渐变球（用于背景装饰）
--gradient-orb-1: radial-gradient(circle, rgba(255, 107, 107, 0.8) 0%, transparent 70%);
--gradient-orb-2: radial-gradient(circle, rgba(78, 205, 196, 0.8) 0%, transparent 70%);
--gradient-orb-3: radial-gradient(circle, rgba(255, 195, 113, 0.8) 0%, transparent 70%);
```

### 色彩使用原则

1. **主色为主** - 界面主要交互使用 Arco Blue (`#165DFF`)
2. **功能色辅助** - 成功、警告、错误等状态使用对应功能色
3. **中性色打底** - 大面积使用中性色，保持界面简洁
4. **品牌色点缀** - 在登录页、品牌展示等场景使用 GameLink 渐变色
5. **对比度保证** - 确保文字与背景对比度 ≥ 4.5:1

---

## 🧩 组件设计方案

基于 Arco Design **60+ 高质量组件**，覆盖中后台常见场景。

### 组件分类

#### 1. 通用组件（General）

| 组件       | 用途 | 示例                        |
| ---------- | ---- | --------------------------- |
| Button     | 按钮 | 主按钮、次要按钮、文本按钮  |
| Icon       | 图标 | @arco-design/web-react/icon |
| Typography | 排版 | 标题、段落、文本            |

#### 2. 布局组件（Layout）

| 组件    | 用途     | 示例                           |
| ------- | -------- | ------------------------------ |
| Layout  | 页面布局 | Header、Sider、Content、Footer |
| Grid    | 栅格布局 | Row、Col (24栅格)              |
| Space   | 间距     | 自动间距布局                   |
| Divider | 分割线   | 水平/垂直分割                  |

#### 3. 导航组件（Navigation）

| 组件       | 用途     | 示例               |
| ---------- | -------- | ------------------ |
| Menu       | 导航菜单 | 侧边菜单、顶部菜单 |
| Breadcrumb | 面包屑   | 页面路径导航       |
| Tabs       | 标签页   | 内容切换           |
| Steps      | 步骤条   | 流程指示           |
| Pagination | 分页     | 数据分页           |

#### 4. 数据录入（Data Entry）

| 组件        | 用途     | 示例          |
| ----------- | -------- | ------------- |
| Form        | 表单     | 数据录入表单  |
| Input       | 输入框   | 文本输入      |
| InputNumber | 数字输入 | 数字输入      |
| Select      | 选择器   | 下拉选择      |
| DatePicker  | 日期选择 | 日期/时间选择 |
| Checkbox    | 复选框   | 多选          |
| Radio       | 单选框   | 单选          |
| Switch      | 开关     | 开关切换      |
| Upload      | 上传     | 文件上传      |

#### 5. 数据展示（Data Display）

| 组件         | 用途     | 示例                             |
| ------------ | -------- | -------------------------------- |
| Table        | 表格     | 数据表格（支持排序、筛选、分页） |
| List         | 列表     | 数据列表                         |
| Card         | 卡片     | 内容容器                         |
| Tree         | 树形控件 | 层级数据展示                     |
| Descriptions | 描述列表 | 详情展示                         |
| Statistic    | 统计数值 | 数据统计展示                     |
| Badge        | 徽标     | 状态标识                         |
| Tag          | 标签     | 标签分类                         |
| Avatar       | 头像     | 用户头像                         |

#### 6. 反馈组件（Feedback）

| 组件         | 用途       | 示例         |
| ------------ | ---------- | ------------ |
| Message      | 消息提示   | 全局提示信息 |
| Notification | 通知       | 系统通知     |
| Modal        | 对话框     | 弹窗对话     |
| Drawer       | 抽屉       | 侧边抽屉面板 |
| Popconfirm   | 气泡确认框 | 操作确认     |
| Progress     | 进度条     | 进度展示     |
| Spin         | 加载中     | 加载状态     |
| Skeleton     | 骨架屏     | 加载占位     |
| Result       | 结果页     | 操作结果反馈 |

### 组件设计原则

#### 1. 系统一致性

- 统一的交互逻辑（如：所有输入框获焦行为一致）
- 统一的视觉风格（如：统一的圆角、间距、阴影）
- 统一的文案表达（如：确认/取消 vs 确定/关闭）

#### 2. 防止错误

- 危险操作二次确认（如：删除前弹出确认框）
- 表单验证提示（如：实时验证并显示错误信息）
- 禁用状态明确（如：禁用按钮灰色且不可点击）

#### 3. 及时反馈

- 操作响应（如：点击按钮有 loading 状态）
- 状态变化（如：数据更新后显示成功提示）
- 进度展示（如：文件上传显示进度条）

#### 4. 高效操作

- 快捷键支持（如：Enter 提交、Esc 关闭）
- 批量操作（如：表格多选批量删除）
- 智能默认（如：表单默认值、记住选择）

---

## 📐 布局系统

### Arco Design 栅格系统

采用 **24 栅格系统**，支持响应式布局。

```tsx
import { Grid } from '@arco-design/web-react';
const { Row, Col } = Grid;

// 基础栅格
<Row>
  <Col span={12}>左侧（50%）</Col>
  <Col span={12}>右侧（50%）</Col>
</Row>

// 响应式栅格
<Row>
  <Col xs={24} sm={12} md={8} lg={6} xl={4}>
    响应式列
  </Col>
</Row>
```

### 间距规范（Spacing Scale）

基于 **4px** 基准单位：

```less
--spacing-mini: 4px; // 0.25rem
--spacing-small: 8px; // 0.5rem
--spacing-medium: 12px; // 0.75rem
--spacing-default: 16px; // 1rem (基准)
--spacing-large: 20px; // 1.25rem
--spacing-xlarge: 24px; // 1.5rem
--spacing-2xlarge: 32px; // 2rem
--spacing-3xlarge: 40px; // 2.5rem
--spacing-4xlarge: 48px; // 3rem
```

**使用建议**：

- 组件内边距：16px (default)
- 卡片内边距：20px-24px (large/xlarge)
- 页面边距：24px-32px (xlarge/2xlarge)
- 区块间距：24px-48px (xlarge/4xlarge)

### 圆角规范（Border Radius）

```less
--border-radius-none: 0;
--border-radius-small: 2px; // 小按钮、标签
--border-radius-medium: 4px; // 标准圆角
--border-radius-large: 8px; // 卡片、输入框
--border-radius-circle: 50%; // 圆形头像
```

### 阴影规范（Box Shadow）

```less
// 卡片阴影
--shadow-1: 0 1px 2px rgba(0, 0, 0, 0.04);
--shadow-2: 0 2px 8px rgba(0, 0, 0, 0.08);
--shadow-3: 0 4px 16px rgba(0, 0, 0, 0.1);
--shadow-4: 0 8px 24px rgba(0, 0, 0, 0.12);

// 悬浮阴影
--shadow-hover: 0 4px 12px rgba(22, 93, 255, 0.15);

// 按钮阴影
--shadow-button: 0 2px 4px rgba(22, 93, 255, 0.2);
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
// 系统字体栈（Arco Design 推荐）
--font-family-base:
  -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, 'PingFang SC',
  'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;

// 等宽字体
--font-family-code: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', 'Courier New', monospace;

// 数字字体
--font-family-number: 'Helvetica Neue', Arial, sans-serif;
```

### 字体大小

```less
--font-size-caption: 12px; // 辅助说明
--font-size-body: 14px; // 正文 (基准)
--font-size-title: 16px; // 小标题
--font-size-h6: 16px; // H6
--font-size-h5: 18px; // H5
--font-size-h4: 20px; // H4
--font-size-h3: 24px; // H3
--font-size-h2: 28px; // H2
--font-size-h1: 32px; // H1
--font-size-display: 36px; // 展示文字
```

### 字重（Font Weight）

```less
--font-weight-light: 300;
--font-weight-regular: 400; // 正文
--font-weight-medium: 500; // 强调
--font-weight-semibold: 600; // 标题
--font-weight-bold: 700; // 加粗
```

### 行高（Line Height）

```less
--line-height-dense: 1.25; // 密集
--line-height-default: 1.5; // 标准
--line-height-relaxed: 1.75; // 宽松
--line-height-loose: 2; // 松散
```

---

## 🎨 主题定制

### Arco 主题配置

使用 [Arco Design 风格配置平台 2.0](https://arco.design/themes) 自定义主题。

#### 1. 在线配置

访问 https://arco.design/themes，在线调整：

- 主色调
- 圆角大小
- 阴影强度
- 字体大小
- 组件样式

#### 2. 导出主题包

下载 JSON 配置文件，应用到项目：

```typescript
// vite.config.ts
import { vitePluginForArco } from '@arco-plugins/vite-react';

export default defineConfig({
  plugins: [
    vitePluginForArco({
      theme: '@arco-themes/gamelink', // 自定义主题包
    }),
  ],
});
```

#### 3. CSS 变量覆盖

```less
// 覆盖 Arco Design 主色
:root {
  --primary-6: #667eea; // GameLink 主色
}
```

### 亮/暗模式切换

```tsx
import { ConfigProvider } from '@arco-design/web-react';

function App() {
  const [theme, setTheme] = useState('light');

  return (
    <ConfigProvider
      componentConfig={{
        Card: { bordered: false },
        List: { bordered: false },
      }}
      theme={theme}
    >
      <YourApp />
    </ConfigProvider>
  );
}
```

---

## 🎬 动画规范

### 动画时长（Duration）

```less
--motion-duration-immediately: 0ms; // 立即
--motion-duration-fast: 100ms; // 快速（微交互）
--motion-duration-moderate: 200ms; // 中速（标准动画）
--motion-duration-slow: 300ms; // 慢速（复杂动画）
--motion-duration-slower: 400ms; // 更慢（页面级）
```

### 缓动函数（Easing）

```less
--motion-ease-linear: cubic-bezier(0, 0, 1, 1);
--motion-ease-in: cubic-bezier(0.4, 0, 1, 1);
--motion-ease-out: cubic-bezier(0, 0, 0.2, 1);
--motion-ease-in-out: cubic-bezier(0.4, 0, 0.2, 1);
--motion-ease-out-back: cubic-bezier(0.12, 0.4, 0.29, 1.46);
--motion-ease-in-back: cubic-bezier(0.71, -0.46, 0.88, 0.6);
--motion-ease-in-out-back: cubic-bezier(0.71, -0.46, 0.29, 1.46);
```

### 动画使用原则

1. **性能优先** - 使用 `transform` 和 `opacity` 而非 `top`/`left`
2. **适度使用** - 避免过度动画导致眩晕
3. **尊重用户** - 支持 `prefers-reduced-motion`
4. **目的明确** - 动画应帮助用户理解界面变化

```less
// 尊重用户动画偏好
@media (prefers-reduced-motion: reduce) {
  * {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}
```

---

## 📱 响应式设计

### 断点定义（Breakpoints）

与 Arco Design Grid 系统保持一致：

```less
@screen-xs: 0px; // < 576px  手机竖屏
@screen-sm: 576px; // ≥ 576px  手机横屏
@screen-md: 768px; // ≥ 768px  平板竖屏
@screen-lg: 992px; // ≥ 992px  平板横屏
@screen-xl: 1200px; // ≥ 1200px 桌面显示器
@screen-xxl: 1600px; // ≥ 1600px 大屏显示器
```

### 响应式栅格

```tsx
<Row gutter={24}>
  {/* 手机 100%，平板 50%，桌面 25% */}
  <Col xs={24} md={12} lg={6}>
    <Card>响应式卡片</Card>
  </Col>
</Row>
```

### 移动优先策略

```less
// 移动优先
.component {
  padding: 16px;
  font-size: 14px;

  @media (min-width: @screen-md) {
    padding: 24px;
    font-size: 16px;
  }

  @media (min-width: @screen-lg) {
    padding: 32px;
  }
}
```

---

## 🎯 组件实例：登录页面

### 设计规范

基于 Arco Design 组件，融合 GameLink 品牌特色。

#### 视觉层级

```
LoginPage
├── Background (品牌渐变背景)
│   ├── GradientOrb × 3 (装饰元素)
├── Container
│   ├── Header
│   │   ├── Logo (品牌标识)
│   │   ├── Title (系统名称)
│   │   └── Subtitle (欢迎语)
│   ├── Card (Arco Card)
│   │   └── Form (Arco Form)
│   │       ├── Input (用户名)
│   │       ├── InputPassword (密码)
│   │       ├── Checkbox (记住我)
│   │       ├── Link (忘记密码)
│   │       ├── Button (登录按钮)
│   │       └── Divider + Tips
│   └── Footer (版权信息)
```

#### 组件配置

```tsx
// 使用 Arco Design 组件
import { Card, Form, Input, Button, Checkbox } from '@arco-design/web-react';

<Card
  bordered={false}
  style={{
    width: 400,
    background: 'rgba(255, 255, 255, 0.95)',
    backdropFilter: 'blur(10px)',
    borderRadius: 20,
    boxShadow: '0 20px 60px rgba(0, 0, 0, 0.2)',
  }}
>
  <Form size="large" labelCol={{ span: 0 }} wrapperCol={{ span: 24 }}>
    <Form.Item
      field="username"
      rules={[
        { required: true, message: '请输入用户名' },
        { minLength: 3, message: '用户名至少3个字符' },
      ]}
    >
      <Input prefix={<IconUser />} placeholder="请输入用户名" allowClear />
    </Form.Item>

    <Form.Item
      field="password"
      rules={[
        { required: true, message: '请输入密码' },
        { minLength: 6, message: '密码至少6个字符' },
      ]}
    >
      <Input.Password prefix={<IconLock />} placeholder="请输入密码" allowClear />
    </Form.Item>

    <Form.Item>
      <Button
        type="primary"
        long
        htmlType="submit"
        style={{
          height: 44,
          background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
          border: 'none',
        }}
      >
        立即登录
      </Button>
    </Form.Item>
  </Form>
</Card>;
```

#### 动画效果

```less
// 页面淡入
@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

// 卡片上滑
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

.loginPage {
  animation: fadeIn 0.6s ease-out;

  .card {
    animation: slideUp 0.6s ease-out;
  }
}
```

---

## ♿ 可访问性（Accessibility）

### WCAG 2.1 AAA 级标准

#### 颜色对比度

- **普通文本**: 最小对比度 **7:1**
- **大文本** (18px+ 或 粗体14px+): 最小对比度 **4.5:1**
- **UI 组件**: 最小对比度 **3:1**

#### 键盘导航

所有交互元素支持键盘操作：

| 按键          | 功能              |
| ------------- | ----------------- |
| Tab           | 焦点前进          |
| Shift + Tab   | 焦点后退          |
| Enter / Space | 激活按钮/链接     |
| Esc           | 关闭对话框/抽屉   |
| Arrow Keys    | 在菜单/列表中导航 |

#### ARIA 标签

```tsx
// 按钮
<Button aria-label="提交表单" aria-busy={loading}>
  提交
</Button>

// 输入框
<Input
  aria-label="用户名"
  aria-required="true"
  aria-invalid={hasError}
  aria-describedby="username-error"
/>

// 错误提示
<div id="username-error" role="alert">
  用户名格式不正确
</div>

// 加载状态
<Spin aria-label="加载中" />

// 页面区域
<main role="main" aria-label="主内容区">
  <nav aria-label="主导航">...</nav>
  <section aria-labelledby="page-title">...</section>
</main>
```

#### 焦点指示器

```less
// 焦点样式
*:focus-visible {
  outline: 2px solid var(--arcoblue-6);
  outline-offset: 2px;
  border-radius: 4px;
}

// 移除默认焦点样式（仅在有自定义样式时）
*:focus:not(:focus-visible) {
  outline: none;
}
```

---

## 🔧 Figma 到代码工作流

### 1. Figma 资源准备

您已下载的 Arco Design Figma 资源包包含：

- ✅ 完整组件库
- ✅ 色彩系统
- ✅ 图标库
- ✅ 模板页面

### 2. 设计阶段

在 Figma 中使用 Arco 组件设计页面：

1. **复用组件** - 使用 Arco 组件库中的组件
2. **保持规范** - 遵循间距、颜色、字体规范
3. **标注说明** - 添加交互说明、状态说明
4. **导出资源** - 导出图标、图片等资源

### 3. 开发阶段

从 Figma 设计到 React 代码：

#### 方法一：使用 Arco Figma 插件

Arco Design 提供 Figma 插件，可自动生成代码：

1. 安装插件：Figma 中搜索 "Arco Design"
2. 选中设计稿
3. 点击"生成代码"
4. 复制 React/Vue 代码

#### 方法二：手动实现

根据设计稿，使用 Arco Design 组件：

```tsx
// Figma 设计: 用户列表卡片
// 实现代码:
import { Card, Table, Button, Space } from '@arco-design/web-react';

<Card title="用户列表" extra={<Button type="primary">添加用户</Button>}>
  <Table columns={columns} data={data} pagination={{ pageSize: 10 }} />
</Card>;
```

### 4. 设计交付清单

设计师交付给开发者：

- [ ] Figma 设计稿链接
- [ ] 组件状态说明（hover、active、disabled等）
- [ ] 交互逻辑说明
- [ ] 响应式断点说明
- [ ] 图标/图片资源（SVG/PNG）
- [ ] 特殊字体文件（如有）
- [ ] 动画效果说明

### 5. 开发检查清单

开发者实现时检查：

- [ ] 使用 Arco Design 组件
- [ ] 遵循设计系统色彩规范
- [ ] 应用正确的间距和布局
- [ ] 实现响应式设计
- [ ] 添加键盘导航支持
- [ ] 添加 ARIA 标签
- [ ] 测试亮/暗模式
- [ ] 验证设计还原度

---

## 📚 资源和工具

### 官方资源

- 🎨 [Arco Design 官网](https://arco.design/)
- 📖 [组件文档](https://arco.design/react/docs/start)
- 🎭 [Figma 组件库](https://www.figma.com/community/file/1068364551746333840)
- 🛠️ [风格配置平台](https://arco.design/themes)
- 💎 [Arco Pro 模板](https://github.com/arco-design/arco-design-pro)

### 开发工具

- **VSCode 插件**: Arco Design Snippets
- **Chrome 插件**: React Developer Tools
- **设计工具**: Figma + Arco 插件

### 参考指南

- [Material Design Guidelines](https://material.io/design)
- [Apple HIG](https://developer.apple.com/design/)
- [WCAG 2.1](https://www.w3.org/WAI/WCAG21/quickref/)
- [MDN Web 文档](https://developer.mozilla.org/)

---

## 🎯 待设计页面规划

### 短期（1-2周）

1. **仪表盘（Dashboard）**
   - 数据统计卡片（Statistic）
   - 图表展示（可集成 ECharts）
   - 快捷操作（Button Group）

2. **用户管理（Users）**
   - 用户列表表格（Table）
   - 筛选和搜索（Form + Input）
   - 用户详情抽屉（Drawer）

### 中期（1个月）

3. **订单管理（Orders）**
   - 订单状态流程（Steps）
   - 订单详情（Descriptions）
   - 数据导出（Button）

4. **权限管理（Permissions）**
   - 权限树（Tree）
   - 角色配置（Form）
   - 权限矩阵（Table）

### 长期（3个月）

5. **系统设置（Settings）**
   - 个人信息（Form）
   - 主题切换（Switch）
   - 通知配置（Checkbox Group）

---

## 📝 更新日志

### v2.0.0 (2025-10-27)

- ✅ 整合 Arco Design 官方色卡系统
- ✅ 完善 60+ 组件使用规范
- ✅ 添加 Figma 到代码工作流
- ✅ 增强可访问性规范（ARIA、键盘导航）
- ✅ 补充主题定制指南
- ✅ 优化响应式设计规范

### v1.0.0 (2025-10-27)

- 初始版本
- 基础色彩系统
- 登录页面设计规范

---

**维护者**: GameLink Frontend Team  
**设计系统**: Arco Design + GameLink Brand  
**最后更新**: 2025-10-27

---

<div align="center">

**基于 [Arco Design](https://arco.design/) 构建**

Made with ❤️ by 字节跳动 + GameLink Team

</div>
