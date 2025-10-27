# 导入路径规范指南

## 统一规范

本项目使用 **绝对路径导入**（基于 `tsconfig.json` 的 `baseUrl` 配置）作为标准。

### ✅ 推荐：绝对路径导入

```typescript
// 从 src 根目录导入
import { Button } from 'components/Button';
import { useAuth } from 'contexts/AuthContext';
import { userService } from 'services/user';
import { API_BASE } from 'config';
import type { User } from 'types/user';
```

### ❌ 避免：相对路径导入

```typescript
// 不推荐：使用相对路径
import { Button } from '../../../components/Button';
import { useAuth } from '../../contexts/AuthContext';
```

## 配置说明

### TypeScript 配置 (`tsconfig.json`)

```json
{
  "compilerOptions": {
    "baseUrl": "./src",
    "paths": {
      "*": ["*"]
    }
  }
}
```

### Vite 配置 (`vite.config.ts`)

```typescript
import { defineConfig } from 'vite';
import path from 'path';

export default defineConfig({
  resolve: {
    alias: {
      // 确保 Vite 也能识别绝对路径
      '~': path.resolve(__dirname, './src'),
    },
  },
});
```

## 导入规范细则

### 1. 组件导入

```typescript
// ✅ 正确
import { Header } from 'components/Header';
import { Footer } from 'components/Footer';

// ❌ 错误
import { Header } from './components/Header';
import { Footer } from '../components/Footer';
```

### 2. 工具函数导入

```typescript
// ✅ 正确
import { formatDate } from 'utils/dateFormatter';
import { storage } from 'utils/storage';

// ❌ 错误
import { formatDate } from '../../utils/dateFormatter';
```

### 3. 类型导入

使用 `import type` 语法：

```typescript
// ✅ 正确
import type { User } from 'types/user';
import type { ApiResponse } from 'types/api';

// ❌ 错误
import { User } from 'types/user'; // 不使用 type 关键字
```

### 4. Hooks 导入

```typescript
// ✅ 正确
import { useAuth } from 'contexts/AuthContext';
import { useTable } from 'hooks/useTable';

// ❌ 错误
import { useAuth } from '../contexts/AuthContext';
```

### 5. 样式导入

```typescript
// ✅ 正确
import styles from './Button.module.less';

// 说明：样式文件使用相对路径，因为它们通常与组件在同一目录
```

### 6. 第三方库导入

```typescript
// ✅ 正确：第三方库保持原样
import React from 'react';
import { Button } from '@arco-design/web-react';
import axios from 'axios';
```

## 导入顺序

遵循以下顺序（由上至下）：

```typescript
// 1. React 相关
import React, { useState, useEffect } from 'react';

// 2. 第三方库
import { Button, Form } from '@arco-design/web-react';
import axios from 'axios';

// 3. 类型导入
import type { User } from 'types/user';
import type { ApiResponse } from 'types/api';

// 4. 项目内部导入（绝对路径）
import { Header } from 'components/Header';
import { useAuth } from 'contexts/AuthContext';
import { userService } from 'services/user';
import { formatDate } from 'utils/dateFormatter';

// 5. 相对路径导入（仅限样式和同目录文件）
import styles from './MyComponent.module.less';
import { helperFunction } from './helpers';
```

## 特殊情况

### 1. 同目录文件

同目录的辅助文件可以使用相对路径：

```typescript
// MyComponent/
// ├── MyComponent.tsx
// ├── helpers.ts
// └── types.ts

// 在 MyComponent.tsx 中
import { helperFunction } from './helpers'; // ✅ 允许
import type { LocalType } from './types'; // ✅ 允许
```

### 2. 测试文件

测试文件导入被测文件时，使用相对路径：

```typescript
// Button.test.tsx
import { Button } from './Button'; // ✅ 允许
```

### 3. index.ts 导出文件

```typescript
// components/Button/index.ts
export { Button } from './Button'; // ✅ 允许
export type { ButtonProps } from './Button';
```

## 迁移指南

如果现有代码使用相对路径，按以下步骤迁移：

### 步骤 1：识别需要改变的导入

```bash
# 查找相对路径导入
grep -r "from '\.\.\/" src/
```

### 步骤 2：替换相对路径

```typescript
// Before
import { Button } from '../../../components/Button';

// After
import { Button } from 'components/Button';
```

### 步骤 3：运行检查

```bash
npm run typecheck  # 检查类型
npm run lint       # 检查代码规范
npm run build      # 确保构建成功
```

## ESLint 规则（可选）

可以配置 ESLint 强制执行导入规范：

```javascript
// eslint.config.js
module.exports = {
  rules: {
    'no-restricted-imports': [
      'error',
      {
        patterns: [
          {
            group: ['../**/components/*', '../**/utils/*', '../**/hooks/*'],
            message: 'Use absolute imports instead of relative imports for shared modules.',
          },
        ],
      },
    ],
  },
};
```

## 常见问题

### Q: 为什么使用绝对路径？

**A:** 
- ✅ 更清晰：一目了然知道文件来自哪个目录
- ✅ 易维护：移动文件时不需要修改导入路径
- ✅ 避免混淆：`../../../utils` vs `utils`

### Q: 什么时候可以使用相对路径？

**A:**
- ✅ 样式文件 (`./Button.module.less`)
- ✅ 同目录辅助文件 (`./helpers.ts`)
- ✅ 测试文件导入被测文件
- ✅ index.ts 导出文件

### Q: IDE 自动导入不是绝对路径怎么办？

**A:** 配置 IDE：

**VSCode:**
```json
// settings.json
{
  "typescript.preferences.importModuleSpecifier": "non-relative",
  "javascript.preferences.importModuleSpecifier": "non-relative"
}
```

**WebStorm:**
Settings → Editor → Code Style → TypeScript → Imports → Use paths relative to tsconfig.json

## 检查清单

在提交代码前，确保：

- [ ] 所有跨目录导入使用绝对路径
- [ ] 类型导入使用 `import type`
- [ ] 导入顺序符合规范
- [ ] `npm run typecheck` 通过
- [ ] `npm run lint` 通过

## 参考资料

- [TypeScript Module Resolution](https://www.typescriptlang.org/docs/handbook/module-resolution.html)
- [Vite Resolve Alias](https://vitejs.dev/config/shared-options.html#resolve-alias)
- [ESLint no-restricted-imports](https://eslint.org/docs/latest/rules/no-restricted-imports)

