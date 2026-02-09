# 🎨 Figma 到代码实现指南

**适用项目**: GameLink Frontend  
**设计系统**: Arco Design  
**更新日期**: 2025-10-27

---

## 📋 目录

1. [准备工作](#准备工作)
2. [设计阶段](#设计阶段)
3. [设计交付](#设计交付)
4. [开发实现](#开发实现)
5. [质量检查](#质量检查)
6. [常见问题](#常见问题)

---

## 🛠️ 准备工作

### 设计师准备

#### 1. 安装 Arco Design Figma 资源

✅ **您已完成**: 下载 Arco Design Figma 资源包

资源包包含：

- 完整组件库（60+ 组件）
- 色彩系统（主色、功能色、中性色）
- 图标库（200+ 图标）
- 字体规范
- 间距系统
- 模板页面

#### 2. Figma 插件推荐

| 插件名称        | 用途           | 必装    |
| --------------- | -------------- | ------- |
| Arco Design     | 组件代码生成   | ✅ 是   |
| Unsplash        | 免费图片素材   | ⭐ 推荐 |
| Iconify         | 图标库         | ⭐ 推荐 |
| Contrast        | 颜色对比度检查 | ⭐ 推荐 |
| Table Generator | 表格生成       | ⚪ 可选 |

### 开发者准备

#### 1. 环境检查

```bash
# 确保已安装 Arco Design
npm list @arco-design/web-react

# 如未安装
npm install @arco-design/web-react @arco-design/web-react/icon
```

#### 2. 熟悉组件库

- 阅读 [Arco Design 文档](https://arco.design/react/docs/start)
- 查看 [组件演示](https://arco.design/react/components/button)
- 了解项目 [DESIGN_SYSTEM.md](./DESIGN_SYSTEM.md)

---

## 🎨 设计阶段

### 1. 创建设计文件

#### 文件组织结构

```
GameLink 设计文件/
├── 📁 0-Resources（资源库）
│   ├── Arco Design Components
│   ├── Icons
│   ├── Images
│   └── Brand Assets
├── 📁 1-Pages（页面设计）
│   ├── Dashboard（仪表盘）
│   ├── Users（用户管理）
│   ├── Orders（订单管理）
│   └── Settings（设置）
└── 📁 2-Components（自定义组件）
    ├── DataCard
    ├── UserAvatar
    └── Tag
```

### 2. 使用 Arco 组件

#### ✅ 推荐做法

```
1. 从 Arco 组件库拖拽组件到画布
2. 使用 Arco 预设的颜色变量
3. 遵循 Arco 的间距系统（4px 倍数）
4. 使用 Arco 图标库
5. 保持组件实例而非分离（Detach）
```

#### ❌ 避免做法

```
1. 自己绘制标准组件（如 Button、Input）
2. 使用非 Arco 色彩系统的颜色
3. 随意设置间距（不符合4px规则）
4. 过度定制导致无法用代码实现
5. 分离组件实例后大幅修改
```

### 3. 设计规范遵循

#### 颜色使用

```
✅ 使用 Arco 色彩变量
  - Primary/arcoblue-6: #165DFF
  - Success/green-6: #00B42A
  - Warning/orange-6: #FF7D00
  - Error/red-6: #F53F3F

❌ 避免随意取色
  - #1234AB (非系统颜色)
  - RGB(100, 200, 50) (非系统颜色)
```

#### 间距使用

```
✅ 使用 4px 倍数
  - 4px, 8px, 12px, 16px, 24px, 32px...

❌ 避免随意数值
  - 13px, 19px, 27px...
```

#### 文字规范

```
✅ 使用系统字体大小
  - 12px (caption)
  - 14px (body)
  - 16px (title)
  - 20px, 24px (heading)

❌ 避免奇怪的字号
  - 13px, 15px, 17px...
```

### 4. 设计页面示例

#### 用户列表页面设计

**Frame**: 1440 × 900 (Desktop)

```
结构:
┌─────────────────────────────────────┐
│ Header (Breadcrumb)                 │ 60px
├─────────────────────────────────────┤
│ ┌─────────────────────────────────┐ │
│ │ Card: 用户列表                   │ │
│ │ ┌────────────────┬──────────┐  │ │
│ │ │ Search & Filter │ Add User │  │ │ 56px
│ │ └────────────────┴──────────┘  │ │
│ │ ┌─────────────────────────────┐ │ │
│ │ │ Table                        │ │ │
│ │ │  - 姓名  邮箱  角色  操作     │ │ │
│ │ │  ...                         │ │ │
│ │ └─────────────────────────────┘ │ │
│ │ Pagination                       │ │ 48px
│ └─────────────────────────────────┘ │
└─────────────────────────────────────┘

间距:
- Page padding: 24px
- Card padding: 20px
- Element gap: 16px
```

**使用的 Arco 组件**:

- Breadcrumb (面包屑)
- Card (卡片容器)
- Input (搜索框)
- Select (筛选下拉)
- Button (操作按钮)
- Table (数据表格)
- Pagination (分页器)

---

## 📤 设计交付

### 1. 设计稿准备

#### Figma 文件整理

```
在交付前确保:
✅ 所有图层命名清晰（如: btn-primary, input-username）
✅ 组件整理成组（Group）或 Frame
✅ 删除无用图层和隐藏元素
✅ 标注特殊交互状态（hover、active、disabled）
✅ 添加 Comments 说明复杂交互
✅ 检查设计稿与 Arco 规范一致性
```

#### 导出资源

**需要导出的资源**:

| 资源类型 | 格式    | 命名规范              | 示例                           |
| -------- | ------- | --------------------- | ------------------------------ |
| 图标     | SVG     | `icon-{name}.svg`     | `icon-user.svg`                |
| 插图     | SVG/PNG | `illustration-{name}` | `illustration-empty-state.png` |
| 头像     | PNG/JPG | `avatar-{name}`       | `avatar-default.png`           |
| 背景     | PNG/JPG | `bg-{name}`           | `bg-login.jpg`                 |
| Logo     | SVG     | `logo-{variant}`      | `logo-light.svg`               |

**导出设置**:

```
SVG:
- 包含 "id"
- 简化描边
- 内联样式

PNG:
- 1x, 2x, 3x（可选）
- 优化压缩
```

### 2. 设计交付清单

给开发者的完整交付物：

```
📁 设计交付/
├── 📄 Figma 链接
│   └── 查看权限已开启
├── 📁 导出资源/
│   ├── icons/
│   ├── images/
│   └── fonts/（如有自定义字体）
├── 📄 交互说明文档
│   ├── 页面流程图
│   ├── 状态说明（hover、active、disabled等）
│   └── 动画效果说明
└── 📄 响应式断点说明
    ├── Mobile: 375px
    ├── Tablet: 768px
    └── Desktop: 1440px
```

### 3. 使用 Arco Figma 插件生成代码

#### 步骤

1. **选中设计元素**
   - 在 Figma 中选中要生成代码的组件/Frame

2. **打开 Arco 插件**
   - Plugins → Arco Design → Generate Code

3. **选择框架**
   - 选择 "React"
   - 选择 "TypeScript"

4. **复制代码**
   - 点击 "Copy Code"
   - 交给开发者

#### 生成的代码示例

```tsx
// Figma 插件自动生成
import { Card, Button, Input, Table } from '@arco-design/web-react';

function UserList() {
  return (
    <Card title="用户列表" extra={<Button type="primary">添加用户</Button>}>
      <Input.Search placeholder="搜索用户..." style={{ marginBottom: 16 }} />
      <Table
        columns={[
          { title: '姓名', dataIndex: 'name' },
          { title: '邮箱', dataIndex: 'email' },
          { title: '角色', dataIndex: 'role' },
        ]}
        data={[]}
        pagination={{ pageSize: 10 }}
      />
    </Card>
  );
}
```

---

## 💻 开发实现

### 1. 创建组件文件

```bash
# 创建组件目录
mkdir -p src/pages/Users

# 创建文件
touch src/pages/Users/Users.tsx
touch src/pages/Users/Users.module.less
touch src/pages/Users/index.ts
```

### 2. 实现基础结构

```tsx
// src/pages/Users/Users.tsx
import { useState } from 'react';
import { Card, Button, Input, Table, Space } from '@arco-design/web-react';
import { IconPlus, IconSearch } from '@arco-design/web-react/icon';
import type { User } from 'types/user';
import styles from './Users.module.less';

export const Users: React.FC = () => {
  const [searchText, setSearchText] = useState('');

  const columns = [
    {
      title: '姓名',
      dataIndex: 'name',
      sorter: (a: User, b: User) => a.name.localeCompare(b.name),
    },
    {
      title: '邮箱',
      dataIndex: 'email',
    },
    {
      title: '角色',
      dataIndex: 'role',
      filters: [
        { text: '管理员', value: 'admin' },
        { text: '用户', value: 'user' },
      ],
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: unknown, record: User) => (
        <Space>
          <Button type="text" size="small">
            编辑
          </Button>
          <Button type="text" size="small" status="danger">
            删除
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div className={styles.usersPage}>
      <Card
        title="用户列表"
        extra={
          <Button type="primary" icon={<IconPlus />}>
            添加用户
          </Button>
        }
      >
        <Input.Search
          prefix={<IconSearch />}
          placeholder="搜索用户..."
          value={searchText}
          onChange={setSearchText}
          style={{ marginBottom: 16, maxWidth: 400 }}
          allowClear
        />

        <Table
          columns={columns}
          data={[]}
          pagination={{
            pageSize: 10,
            showTotal: true,
            showJumper: true,
          }}
          stripe
          borderCell
        />
      </Card>
    </div>
  );
};
```

### 3. 添加样式

```less
// src/pages/Users/Users.module.less
.usersPage {
  padding: 24px;

  :global {
    .arco-card {
      border-radius: 8px;
      box-shadow: var(--shadow-2);
    }

    .arco-table {
      font-size: 14px;
    }
  }
}

// 响应式
@media (max-width: 768px) {
  .usersPage {
    padding: 16px;
  }
}
```

### 4. 对照设计稿检查

使用浏览器开发者工具对比：

```
✅ 检查项:
  □ 颜色是否一致
  □ 间距是否一致
  □ 字体大小是否一致
  □ 圆角是否一致
  □ 阴影是否一致
  □ 布局是否一致
  □ 响应式断点是否正确
```

### 5. 实现交互状态

```tsx
// 添加 hover、active、disabled 状态

// Hover 效果（通常 Arco 组件已内置）
<Button
  type="primary"
  // Arco 自动处理 hover 样式
>
  按钮
</Button>

// Disabled 状态
<Button disabled>
  禁用按钮
</Button>

// Loading 状态
<Button loading={isSubmitting}>
  {isSubmitting ? '提交中...' : '提交'}
</Button>
```

---

## ✅ 质量检查

### 1. 设计还原度检查

#### 视觉对比

```bash
# 1. 打开 Figma 设计稿
# 2. 打开本地开发服务器
npm run dev

# 3. 使用浏览器插件对比
# - PerfectPixel
# - PixelParallel
```

#### 检查清单

```
视觉还原:
  □ 颜色 100% 匹配
  □ 间距误差 < 2px
  □ 字体大小一致
  □ 圆角一致
  □ 阴影效果一致

交互还原:
  □ Hover 状态正确
  □ Focus 状态正确
  □ Active 状态正确
  □ Disabled 状态正确
  □ Loading 状态正确

响应式:
  □ Mobile 布局正确
  □ Tablet 布局正确
  □ Desktop 布局正确
  □ 断点切换流畅
```

### 2. 代码质量检查

```bash
# TypeScript 类型检查
npm run typecheck

# ESLint 代码检查
npm run lint

# 单元测试
npm run test

# 构建检查
npm run build
```

### 3. 可访问性检查

```bash
# 使用 axe DevTools 浏览器插件
# 或手动检查:

✅ 可访问性检查:
  □ 所有交互元素可键盘访问
  □ 焦点指示器清晰可见
  □ ARIA 标签完整
  □ 颜色对比度 ≥ 4.5:1
  □ 图片有 alt 文本
  □ 表单有 label
```

### 4. 性能检查

```bash
# Lighthouse 检查
# Chrome DevTools → Lighthouse → Generate Report

目标分数:
  Performance: ≥ 90
  Accessibility: 100
  Best Practices: ≥ 95
  SEO: ≥ 90
```

---

## ❓ 常见问题

### Q1: Figma 设计稿中的组件在代码里找不到？

**A:** 检查以下几点：

1. 确认使用的是 Arco Design 组件
2. 查看 [Arco 组件文档](https://arco.design/react/components/overview)
3. 可能是自定义组件，需要手动实现

### Q2: 生成的代码无法运行？

**A:** Figma 插件生成的是模板代码，需要：

1. 补充完整的 imports
2. 添加状态管理（useState、useEffect）
3. 连接 API 数据
4. 添加事件处理函数

### Q3: 设计稿颜色与实际效果不一致？

**A:** 可能原因：

1. 显示器色彩配置不同
2. 浏览器渲染差异
3. 未使用 Arco 色彩变量

解决方案：

```tsx
// 使用 Arco CSS 变量
style={{ color: 'rgb(var(--arcoblue-6))' }}

// 或在 LESS 中
color: var(--color-primary-6);
```

### Q4: 响应式布局不生效？

**A:** 检查：

1. 是否使用了 Arco Grid 组件
2. 是否设置了正确的断点（xs、sm、md、lg、xl）
3. 是否添加了媒体查询

```tsx
// 正确的响应式写法
<Row>
  <Col xs={24} md={12} lg={8}>
    <Card>响应式卡片</Card>
  </Col>
</Row>
```

### Q5: Figma 插件生成的代码太简单？

**A:** Figma 插件只能生成基础结构，需要手动补充：

- 状态管理
- 数据请求
- 表单验证
- 错误处理
- 动画效果

这是正常的，插件只是加速初始开发。

---

## 📚 参考资源

### 官方文档

- [Arco Design React](https://arco.design/react/docs/start)
- [Arco Figma 资源](https://www.figma.com/community/file/1068364551746333840)
- [Arco 主题配置](https://arco.design/themes)

### 工具推荐

- **Figma 插件**: Arco Design, Contrast, Iconify
- **Chrome 插件**: React DevTools, PerfectPixel, axe DevTools
- **VSCode 插件**: Arco Design Snippets

### 学习资源

- [Figma 官方教程](https://www.figma.com/resources/learn-design/)
- [React 官方文档](https://react.dev/)
- [TypeScript 手册](https://www.typescriptlang.org/docs/)

---

## 📝 更新日志

### v1.0.0 (2025-10-27)

- 初始版本
- Figma 到代码完整流程
- 质量检查清单
- 常见问题解答

---

**维护者**: GameLink Frontend Team  
**最后更新**: 2025-10-27

---

<div align="center">

**从设计到代码，让开发更高效** 🚀

[返回设计系统](./DESIGN_SYSTEM.md)

</div>
