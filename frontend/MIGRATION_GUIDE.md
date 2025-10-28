# GameLink 设计系统迁移指南

从 Arco Design 到自定义 Neo-brutalism 黑白组件库

---

## 📋 迁移概述

**迁移时间**: 2025-10-28  
**原框架**: Arco Design  
**新框架**: 自定义组件库（Neo-brutalism 风格）  
**设计理念**: 纯黑白 + 极简 + 高对比度

---

## 🎯 为什么迁移？

### 设计需求

用户明确要求：

- ✅ 纯黑白配色（无任何彩色）
- ✅ Neo-brutalism 风格（直角、粗边框、实体阴影）
- ✅ 极简主义
- ✅ 适度动效

Arco Design 的问题：

- ❌ 彩色丰富（主题蓝、成功绿、警告橙等）
- ❌ 圆角设计（4px-8px 圆角）
- ❌ 柔和阴影（渐变阴影）
- ❌ 样式覆盖复杂

---

## 🔄 迁移对比

### 组件映射表

| Arco Design              | 自定义组件        | 说明                   |
| ------------------------ | ----------------- | ---------------------- |
| `@arco-design/web-react` | `components`      | 从 src/components 导入 |
| `<Button>`               | `<Button>`        | 保持相同 API           |
| `<Input>`                | `<Input>`         | 保持相同 API           |
| `<Input.Password>`       | `<PasswordInput>` | 独立组件               |
| `<Card>`                 | `<Card>`          | 简化 API               |
| `<Form>`                 | `<Form>`          | 简化表单验证           |
| `<Form.Item>`            | `<FormItem>`      | 独立组件               |

### API 变化

#### Button 组件

```tsx
// 之前 (Arco Design)
import { Button } from '@arco-design/web-react';

<Button type="primary" size="large" long>
  登录
</Button>;

// 之后 (自定义)
import { Button } from 'components';

<Button variant="primary" size="large" block>
  登录
</Button>;
```

**变化**：

- `type` → `variant`
- `long` → `block`

#### Input 组件

```tsx
// 之前
import { Input } from '@arco-design/web-react';
import { IconUser } from '@arco-design/web-react/icon';

<Input prefix={<IconUser />} allowClear />;

// 之后
import { Input } from 'components';

const UserIcon = () => <svg>...</svg>;

<Input prefix={<UserIcon />} allowClear />;
```

**变化**：

- 需要自定义图标（不再依赖图标库）
- API 保持一致

#### PasswordInput 组件

```tsx
// 之前
<Input.Password prefix={<IconLock />} />;

// 之后
import { PasswordInput } from 'components';

<PasswordInput prefix={<LockIcon />} />;
```

**变化**：

- 独立组件，不再是 Input 的子组件

#### Form 组件

```tsx
// 之前
import { Form } from '@arco-design/web-react';

<Form>
  <Form.Item field="username" rules={[{ required: true }]}>
    <Input />
  </Form.Item>
</Form>;

// 之后
import { Form, FormItem } from 'components';

<Form onSubmit={handleSubmit}>
  <FormItem error={errors.username}>
    <Input value={username} onChange={handleChange} />
  </FormItem>
</Form>;
```

**变化**：

- `Form.Item` → `FormItem`（独立组件）
- 移除内置表单验证（需手动处理）
- 通过 `error` prop 显示错误

---

## 📦 文件结构变化

### 新增文件

```
src/
├── styles/                    # 新增：全局样式
│   ├── variables.less         # CSS 变量定义
│   └── global.less            # 全局样式重置
├── components/                # 新增：自定义组件库
│   ├── Button/
│   │   ├── Button.tsx
│   │   ├── Button.module.less
│   │   └── index.ts
│   ├── Input/
│   │   ├── Input.tsx
│   │   ├── Input.module.less
│   │   └── index.ts
│   ├── Card/
│   ├── Form/
│   └── index.ts               # 统一导出
└── pages/
    └── Login/
        ├── Login.tsx          # 重构：使用新组件
        ├── Login.module.less  # 保持黑白风格
        └── README.md          # 登录页文档
```

### 修改文件

```
src/
├── main.tsx                   # 修改：导入全局样式
└── App.tsx                    # 修改：移除 ConfigProvider
```

### 删除文件

```
- node_modules/@arco-design/   # 卸载 Arco Design
```

---

## 🎨 样式系统变化

### CSS 变量

**之前** (Arco Design):

```less
--primary-6: #165dff; // Arco 蓝
--arcoblue-6: #165dff;
--border-radius-medium: 4px;
```

**之后** (自定义):

```less
--color-black: #000000; // 纯黑
--color-white: #ffffff; // 纯白
--border-radius-none: 0; // 无圆角
--shadow-base: 8px 8px 0 #000; // 实体阴影
```

### 全局样式

```tsx
// 之前
import '@arco-design/web-react/dist/css/arco.css';

// 之后
import './styles/global.less';
```

---

## 🛠️ 开发工作流

### 安装依赖

```bash
# 卸载 Arco Design
npm uninstall @arco-design/web-react @arco-design/web-react/icon

# 无需额外安装（使用自定义组件）
```

### 启动开发

```bash
npm run dev
```

访问 http://localhost:5174/

### 构建生产

```bash
npm run build
```

---

## 🎯 迁移检查清单

### 代码迁移

- [x] 卸载 Arco Design 依赖
- [x] 创建全局样式系统
- [x] 创建基础组件（Button, Input, Card, Form）
- [x] 重构 Login 页面
- [x] 更新 main.tsx 导入
- [x] 更新 App.tsx 移除 ConfigProvider
- [x] 创建设计系统文档

### 样式迁移

- [x] 定义 CSS 变量
- [x] 创建全局样式重置
- [x] 实现黑白配色
- [x] 实现实体阴影
- [x] 实现直角设计
- [x] 实现动效系统

### 测试验证

- [x] 登录页面正常显示
- [x] 表单验证正常工作
- [x] 按钮交互正常
- [x] 输入框交互正常
- [x] 响应式布局正常
- [x] 无 Linter 错误

---

## 📝 组件使用示例

### 完整登录页面

```tsx
import { useState, FormEvent } from 'react';
import { Button, Input, PasswordInput, Form, FormItem } from 'components';

export const Login = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [errors, setErrors] = useState<any>({});
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();

    // 验证
    const newErrors: any = {};
    if (!username) newErrors.username = '请输入用户名';
    if (!password) newErrors.password = '请输入密码';

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    // 提交
    setLoading(true);
    try {
      await login(username, password);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Form onSubmit={handleSubmit}>
      <FormItem error={errors.username}>
        <Input
          prefix={<UserIcon />}
          placeholder="用户名"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          allowClear
        />
      </FormItem>

      <FormItem error={errors.password}>
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
  );
};
```

---

## 🚀 后续计划

### 短期（1周内）

- [ ] 创建 Message 提示组件
- [ ] 创建 Modal 弹窗组件
- [ ] 创建 Table 表格组件
- [ ] 创建 Checkbox / Radio 组件

### 中期（1个月内）

- [ ] 创建完整的表单验证系统
- [ ] 创建 Layout 布局组件
- [ ] 创建 Menu 菜单组件
- [ ] 创建 Dashboard 仪表盘页面

### 长期（3个月内）

- [ ] 完善设计系统文档
- [ ] 创建 Storybook 组件展示
- [ ] 编写单元测试
- [ ] 性能优化

---

## ❓ 常见问题

### Q: 为什么不使用 TailwindCSS？

A: TailwindCSS 虽然强大，但：

1. 会生成大量 utility classes
2. 不符合极简设计理念
3. 自定义组件更灵活，完全掌控

### Q: 为什么不保留 Arco Design？

A: Arco Design 的设计语言与 Neo-brutalism 冲突：

1. 彩色主题 vs 纯黑白
2. 圆角设计 vs 直角设计
3. 柔和阴影 vs 实体阴影

### Q: 如何添加新组件？

A: 遵循现有组件结构：

```bash
src/components/NewComponent/
├── NewComponent.tsx          # 组件实现
├── NewComponent.module.less  # 组件样式
├── index.ts                  # 导出
└── types.ts                  # 类型定义（可选）
```

然后在 `src/components/index.ts` 中导出。

### Q: 如何自定义主题？

A: 修改 `src/styles/variables.less` 中的 CSS 变量：

```less
:root {
  --color-black: #000000; // 改成其他颜色
  --shadow-base: 8px 8px 0 var(--color-black);
}
```

---

## 📚 参考资源

- [DESIGN_SYSTEM_V2.md](./DESIGN_SYSTEM_V2.md) - 新设计系统文档
- [Login README](./src/pages/Login/README.md) - 登录页面文档
- [Neo-brutalism](https://brutalistwebsites.com/) - 设计参考
- [WCAG 2.1](https://www.w3.org/WAI/WCAG21/quickref/) - 可访问性标准

---

**迁移完成时间**: 2025-10-28  
**负责人**: GameLink Frontend Team  
**审核状态**: ✅ 通过

---

<div align="center">

**从彩色到黑白 · 从复杂到极简**

🎨 → ⚫⚪

</div>
