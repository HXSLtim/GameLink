# 🛠️ GameLink 前端开发指南

**更新时间**: 2025-01-28
**文档类型**: 开发指南
**适用对象**: 前端开发人员、全栈开发人员

---

## 🎯 开发环境配置

### 前置要求

- **Node.js**: 18.0+
- **npm**: 8.0+ 或 **yarn**: 1.22+
- **Git**: 2.30+
- **现代浏览器**: Chrome 90+, Firefox 88+, Safari 14+

### 快速开始

```bash
# 1. 克隆项目
git clone https://github.com/your-org/gamelink.git
cd gamelink/frontend

# 2. 安装依赖
npm install

# 3. 启动开发服务器
npm run dev

# 4. 访问应用
# Local: http://localhost:5173
# Network: http://[your-ip]:5173
```

### 环境变量配置

创建 `.env.local` 文件：

```bash
# API 配置
VITE_API_BASE_URL=http://localhost:8080/api/v1

# 加密配置
VITE_CRYPTO_ENABLED=true
VITE_CRYPTO_SECRET_KEY=your-32-byte-secret-key-here-123456
VITE_CRYPTO_IV=your-iv-16-byte

# 开发配置
VITE_DEV_TOOLS=true
```

---

## 🏗️ 项目架构

### 技术栈

```
React 18.x         # 前端框架
TypeScript 5.x     # 类型系统
Vite 5.x           # 构建工具
Arco Design        # UI 组件库
Less               # CSS 预处理器
Vitest             # 测试框架
ESLint             # 代码检查
Prettier           # 代码格式化
```

### 目录结构

```
frontend/
├── public/                 # 静态资源
├── src/
│   ├── api/               # API 调用层
│   │   ├── client.ts      # HTTP 客户端
│   │   └── *.ts           # 各模块 API
│   ├── components/        # 可复用组件
│   │   ├── Button/        # 按钮组件
│   │   ├── Table/         # 表格组件
│   │   └── index.ts       # 组件导出
│   ├── contexts/          # React Context
│   │   ├── AuthContext.tsx
│   │   └── ThemeContext.tsx
│   ├── layouts/           # 布局组件
│   ├── pages/             # 页面组件
│   │   ├── Dashboard/
│   │   ├── Orders/
│   │   └── Users/
│   ├── services/          # 业务服务层
│   ├── types/             # TypeScript 类型
│   ├── utils/             # 工具函数
│   ├── styles/            # 样式文件
│   ├── i18n/              # 国际化
│   ├── router/            # 路由配置
│   ├── hooks/             # 自定义 Hooks
│   └── main.tsx           # 应用入口
├── docs/                  # 项目文档
├── tests/                 # 测试文件
├── package.json
├── tsconfig.json
├── vite.config.ts
└── .eslintrc.cjs
```

---

## 🔧 开发工作流

### 分支管理

```bash
# 主分支
main                    # 生产环境代码
develop                 # 开发环境代码

# 功能分支
feature/user-management # 新功能开发
feature/order-system    # 订单功能开发

# 修复分支
hotfix/login-bug        # 紧急修复
```

### 提交规范

```bash
# 提交格式
<type>(<scope>): <subject>

# 示例
feat(user): add user registration feature
fix(order): resolve order status update issue
docs(api): update payment API documentation
style(ui): improve button hover effects
refactor(auth): simplify login logic
test(components): add unit tests for Button
chore(deps): update dependencies
```

### 开发命令

```bash
# 开发环境
npm run dev              # 启动开发服务器
npm run dev -- --host    # 允许外部访问

# 构建
npm run build            # 构建生产版本
npm run build:analyze    # 分析构建产物

# 预览
npm run preview          # 预览构建结果

# 代码质量
npm run lint             # ESLint 检查
npm run lint:fix         # 自动修复 ESLint 问题
npm run format           # Prettier 格式化
npm run typecheck        # TypeScript 类型检查

# 测试
npm run test             # 运行测试（监听模式）
npm run test:run         # 运行测试（单次）
npm run test:coverage    # 生成覆盖率报告

# 依赖管理
npm outdated             # 检查过期依赖
npm update               # 更新依赖
npm audit                # 安全审计
```

---

## 🎨 组件开发规范

### 组件结构

```typescript
// Button/Button.tsx
import React from 'react';
import styles from './Button.module.less';

export interface ButtonProps {
  children: React.ReactNode;
  variant?: 'primary' | 'secondary' | 'outlined';
  size?: 'small' | 'medium' | 'large';
  disabled?: boolean;
  onClick?: () => void;
}

export const Button: React.FC<ButtonProps> = ({
  children,
  variant = 'primary',
  size = 'medium',
  disabled = false,
  onClick
}) => {
  return (
    <button
      className={`${styles.button} ${styles[variant]} ${styles[size]}`}
      disabled={disabled}
      onClick={onClick}
    >
      {children}
    </button>
  );
};
```

### 样式规范

```less
// Button/Button.module.less
.button {
  // 基础样式
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: none;
  cursor: pointer;
  transition: all var(--duration-base) var(--ease-out);

  // 变体样式
  &.primary {
    background-color: var(--color-primary);
    color: var(--color-white);
  }

  &.outlined {
    background-color: transparent;
    border: var(--border-width-base) solid var(--color-primary);
    color: var(--color-primary);
  }

  // 尺寸样式
  &.small {
    height: var(--button-height-sm);
    padding: 0 var(--spacing-sm);
  }

  // 状态样式
  &:hover:not(:disabled) {
    opacity: 0.8;
  }

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
}
```

### 类型定义

```typescript
// types/button.ts
export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'outlined' | 'text';
  size?: 'small' | 'medium' | 'large';
  loading?: boolean;
  icon?: ReactNode;
  block?: boolean;
}
```

---

## 📡 API 开发

### API 客户端配置

```typescript
// api/client.ts
import axios, { AxiosInstance, AxiosRequestConfig } from 'axios';

class ApiClient {
  private instance: AxiosInstance;

  constructor() {
    this.instance = axios.create({
      baseURL: import.meta.env.VITE_API_BASE_URL,
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    this.setupInterceptors();
  }

  private setupInterceptors() {
    // 请求拦截器
    this.instance.interceptors.request.use(
      (config) => {
        // 添加认证 token
        const token = localStorage.getItem('auth_token');
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => Promise.reject(error),
    );

    // 响应拦截器
    this.instance.interceptors.response.use(
      (response) => response.data,
      (error) => {
        // 统一错误处理
        if (error.response?.status === 401) {
          // 处理认证失败
          this.handleAuthError();
        }
        return Promise.reject(error);
      },
    );
  }

  async get<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    return this.instance.get(url, config);
  }

  async post<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    return this.instance.post(url, data, config);
  }

  // ... 其他 HTTP 方法
}

export const apiClient = new ApiClient();
```

### API 服务示例

```typescript
// services/user.ts
import { apiClient } from '../api/client';
import type { User, CreateUserData, UpdateUserData } from '../types/user';

export const userService = {
  // 获取用户列表
  async getList(params: {
    page?: number;
    page_size?: number;
    keyword?: string;
  }): Promise<{ list: User[]; total: number }> {
    return apiClient.get('/admin/users', { params });
  },

  // 获取用户详情
  async getDetail(id: number): Promise<User> {
    return apiClient.get(`/admin/users/${id}`);
  },

  // 创建用户
  async create(data: CreateUserData): Promise<User> {
    return apiClient.post('/admin/users', data);
  },

  // 更新用户
  async update(id: number, data: UpdateUserData): Promise<User> {
    return apiClient.put(`/admin/users/${id}`, data);
  },

  // 删除用户
  async delete(id: number): Promise<void> {
    return apiClient.delete(`/admin/users/${id}`);
  },
};
```

### 错误处理

```typescript
// utils/errorHandler.ts
export class ApiError extends Error {
  constructor(
    public code: number,
    message: string,
    public details?: any,
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

export const handleApiError = (error: any): ApiError => {
  if (error.response) {
    const { status, data } = error.response;
    return new ApiError(status, data.message || '请求失败', data);
  } else if (error.request) {
    return new ApiError(0, '网络请求失败');
  } else {
    return new ApiError(-1, error.message || '未知错误');
  }
};
```

---

## 🎨 主题系统

### 主题配置

```typescript
// contexts/ThemeContext.tsx
import React, { createContext, useContext, useEffect, useState } from 'react';

interface ThemeContextType {
  theme: 'light' | 'dark';
  toggleTheme: () => void;
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

export const ThemeProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [theme, setTheme] = useState<'light' | 'dark'>(() => {
    return localStorage.getItem('theme') as 'light' | 'dark' || 'light';
  });

  useEffect(() => {
    document.body.classList.toggle('dark-theme', theme === 'dark');
    localStorage.setItem('theme', theme);
  }, [theme]);

  const toggleTheme = () => {
    setTheme(prev => prev === 'light' ? 'dark' : 'light');
  };

  return (
    <ThemeContext.Provider value={{ theme, toggleTheme }}>
      {children}
    </ThemeContext.Provider>
  );
};

export const useTheme = () => {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error('useTheme must be used within ThemeProvider');
  }
  return context;
};
```

### CSS 变量系统

```less
// styles/variables.less
:root {
  // 主色调
  --color-primary: #165dff;
  --color-success: #00b42a;
  --color-warning: #ff7d00;
  --color-error: #f53f3f;

  // 中性色
  --color-white: #ffffff;
  --color-black: #000000;
  --color-gray-100: #f7f8fa;
  --color-gray-900: #1f2329;

  // 文本颜色（浅色模式）
  --text-primary: var(--color-black);
  --text-secondary: #4e5969;
  --text-tertiary: #86909c;

  // 背景颜色（浅色模式）
  --bg-primary: var(--color-white);
  --bg-secondary: var(--color-gray-100);
  --bg-tertiary: #f2f3f5;
}

body.dark-theme {
  // 文本颜色（深色模式）
  --text-primary: var(--color-white);
  --text-secondary: #c9cdd4;
  --text-tertiary: #86909c;

  // 背景颜色（深色模式）
  --bg-primary: #17171a;
  --bg-secondary: #252529;
  --bg-tertiary: #2e2e33;
}
```

---

## 🌍 国际化

### i18n 配置

```typescript
// i18n/index.ts
import { createI18n } from 'vue-i18n';
import zhCN from './locales/zh-CN.json';
import enUS from './locales/en-US.json';

export const i18n = createI18n({
  legacy: false,
  locale: localStorage.getItem('locale') || 'zh-CN',
  fallbackLocale: 'zh-CN',
  messages: {
    'zh-CN': zhCN,
    'en-US': enUS,
  },
});
```

### 语言文件

```json
// i18n/locales/zh-CN.json
{
  "common": {
    "confirm": "确认",
    "cancel": "取消",
    "save": "保存",
    "delete": "删除",
    "edit": "编辑",
    "search": "搜索"
  },
  "user": {
    "title": "用户管理",
    "name": "姓名",
    "email": "邮箱",
    "phone": "电话",
    "status": "状态",
    "create": "新建用户",
    "edit": "编辑用户"
  }
}
```

### 使用示例

```typescript
// 组件中使用
import { useTranslation } from 'react-i18next';

const UserList = () => {
  const { t } = useTranslation();

  return (
    <div>
      <h1>{t('user.title')}</h1>
      <button>{t('common.save')}</button>
    </div>
  );
};
```

---

## 🧪 测试

### 单元测试

```typescript
// components/Button/Button.test.tsx
import { render, screen, fireEvent } from '@testing-library/react';
import { Button } from './Button';

describe('Button', () => {
  it('renders correctly', () => {
    render(<Button>Click me</Button>);
    expect(screen.getByRole('button')).toBeInTheDocument();
    expect(screen.getByText('Click me')).toBeInTheDocument();
  });

  it('handles click events', () => {
    const handleClick = jest.fn();
    render(<Button onClick={handleClick}>Click me</Button>);

    fireEvent.click(screen.getByRole('button'));
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it('applies variant styles correctly', () => {
    render(<Button variant="outlined">Outlined</Button>);
    expect(screen.getByRole('button')).toHaveClass('outlined');
  });

  it('is disabled when disabled prop is true', () => {
    render(<Button disabled>Disabled</Button>);
    expect(screen.getByRole('button')).toBeDisabled();
  });
});
```

### 集成测试

```typescript
// services/user.test.ts
import { userService } from '../services/user';
import { apiClient } from '../api/client';

jest.mock('../api/client');
const mockedApiClient = apiClient as jest.Mocked<typeof apiClient>;

describe('UserService', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should fetch user list', async () => {
    const mockData = {
      list: [{ id: 1, name: 'John Doe' }],
      total: 1,
    };

    mockedApiClient.get.mockResolvedValue(mockData);

    const result = await userService.getList({ page: 1 });

    expect(result).toEqual(mockData);
    expect(mockedApiClient.get).toHaveBeenCalledWith('/admin/users', {
      params: { page: 1 },
    });
  });
});
```

---

## 🔐 安全最佳实践

### 输入验证

```typescript
// utils/validation.ts
export const validateEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
};

export const validatePhone = (phone: string): boolean => {
  const phoneRegex = /^1[3-9]\d{9}$/;
  return phoneRegex.test(phone);
};

export const sanitizeInput = (input: string): string => {
  return input.trim().replace(/[<>]/g, '');
};
```

### XSS 防护

```typescript
// utils/sanitize.ts
import DOMPurify from 'dompurify';

export const sanitizeHtml = (html: string): string => {
  return DOMPurify.sanitize(html);
};

export const escapeHtml = (text: string): string => {
  const map: { [key: string]: string } = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#39;',
  };

  return text.replace(/[&<>"']/g, (char) => map[char]);
};
```

### 敏感信息处理

```typescript
// utils/security.ts
export const maskSensitiveData = (data: string, visibleChars: number = 4): string => {
  if (!data || data.length <= visibleChars) {
    return data;
  }

  const start = data.substring(0, visibleChars);
  const end = data.substring(data.length - visibleChars);
  const mask = '*'.repeat(data.length - visibleChars * 2);

  return `${start}${mask}${end}`;
};

// 使用示例
const phone = '13812345678';
const maskedPhone = maskSensitiveData(phone, 3); // 138*****678
```

---

## 🚀 性能优化

### 代码分割

```typescript
// router/index.ts
import { lazy, Suspense } from 'react';
import { createBrowserRouter } from 'react-router-dom';

// 懒加载组件
const Dashboard = lazy(() => import('../pages/Dashboard'));
const UserList = lazy(() => import('../pages/Users/UserList'));
const OrderList = lazy(() => import('../pages/Orders/OrderList'));

export const router = createBrowserRouter([
  {
    path: '/',
    element: <Layout />,
    children: [
      {
        path: 'dashboard',
        element: (
          <Suspense fallback={<div>Loading...</div>}>
            <Dashboard />
          </Suspense>
        ),
      },
      // ... 其他路由
    ],
  },
]);
```

### 组件优化

```typescript
// 使用 React.memo 优化组件
export const UserCard = React.memo<UserCardProps>(({ user, onUpdate }) => {
  return (
    <div className="user-card">
      <h3>{user.name}</h3>
      <p>{user.email}</p>
      <button onClick={() => onUpdate(user.id)}>
        编辑
      </button>
    </div>
  );
});

// 使用 useMemo 优化计算
const ExpensiveComponent = ({ data }: { data: any[] }) => {
  const processedData = useMemo(() => {
    return data.map(item => expensiveCalculation(item));
  }, [data]);

  return <div>{/* 渲染处理后的数据 */}</div>;
};

// 使用 useCallback 优化函数
const ParentComponent = () => {
  const [count, setCount] = useState(0);

  const handleUpdate = useCallback((id: number) => {
    // 处理更新逻辑
  }, []);

  return (
    <div>
      <UserCard user={userData} onUpdate={handleUpdate} />
      <button onClick={() => setCount(count + 1)}>
        Count: {count}
      </button>
    </div>
  );
};
```

### 资源优化

```typescript
// 图片懒加载
import { useState, useRef, useEffect } from 'react';

const LazyImage = ({ src, alt, ...props }: any) => {
  const [isLoaded, setIsLoaded] = useState(false);
  const [isInView, setIsInView] = useState(false);
  const imgRef = useRef<HTMLImageElement>(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsInView(true);
          observer.disconnect();
        }
      },
      { threshold: 0.1 }
    );

    if (imgRef.current) {
      observer.observe(imgRef.current);
    }

    return () => observer.disconnect();
  }, []);

  return (
    <div ref={imgRef} {...props}>
      {isInView && (
        <img
          src={src}
          alt={alt}
          onLoad={() => setIsLoaded(true)}
          style={{ opacity: isLoaded ? 1 : 0 }}
        />
      )}
    </div>
  );
};
```

---

## 📦 构建和部署

### 构建配置

```typescript
// vite.config.ts
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { resolve } from 'path';

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      '@components': resolve(__dirname, 'src/components'),
      '@pages': resolve(__dirname, 'src/pages'),
      '@utils': resolve(__dirname, 'src/utils'),
    },
  },
  build: {
    target: 'es2015',
    outDir: 'dist',
    assetsDir: 'assets',
    sourcemap: false,
    minify: 'terser',
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom'],
          ui: ['@arco-design/web-react'],
          utils: ['lodash-es', 'dayjs'],
        },
      },
    },
    chunkSizeWarningLimit: 1000,
  },
  server: {
    port: 5173,
    host: true,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
      },
    },
  },
});
```

### Docker 配置

```dockerfile
# Dockerfile
FROM node:18-alpine as builder

WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### CI/CD 配置

```yaml
# .github/workflows/deploy.yml
name: Deploy to Production

on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
          cache: 'npm'

      - run: npm ci
      - run: npm run lint
      - run: npm run test:coverage
      - run: npm run build

  deploy:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Deploy to server
        run: |
          # 部署脚本
          scp -r dist/* user@server:/var/www/html/
```

---

## 🐛 调试技巧

### 浏览器调试

```typescript
// 添加调试日志
const debug = (message: string, data?: any) => {
  if (import.meta.env.DEV) {
    console.log(`[DEBUG] ${message}`, data);
  }
};

// 性能监控
const measurePerformance = (name: string, fn: () => void) => {
  const start = performance.now();
  fn();
  const end = performance.now();
  console.log(`${name} took ${end - start} milliseconds`);
};
```

### React DevTools

- 安装 React DevTools 浏览器扩展
- 使用 Components 面板检查组件层次
- 使用 Profiler 面板分析性能

### 网络调试

```typescript
// 添加请求拦截器用于调试
if (import.meta.env.DEV) {
  apiClient.interceptors.request.use((config) => {
    console.log(`[API Request] ${config.method?.toUpperCase()} ${config.url}`, config.data);
    return config;
  });

  apiClient.interceptors.response.use((response) => {
    console.log(`[API Response] ${response.config.url}`, response.data);
    return response;
  });
}
```

---

## 📚 学习资源

### 官方文档

- [React 官方文档](https://react.dev/)
- [TypeScript 官方文档](https://www.typescriptlang.org/)
- [Vite 官方文档](https://vitejs.dev/)
- [Arco Design 文档](https://arco.design/)

### 推荐博客

- [React 博客](https://react.dev/blog)
- [TypeScript Weekly](https://www.typescriptlang.org/blog)
- [前端周刊](https://github.com/FrontendMagazine/frontend-weekly)

### 在线课程

- React 完整指南
- TypeScript 深入浅出
- 前端性能优化实战

---

## 🔗 相关链接

- [技术文档总览](./TECHNICAL_DOCUMENTATION.md)
- [用户使用指南](./USER_DOCUMENTATION.md)
- [项目 README](../README.md)
- [API 接口文档](./api/)

---

**文档维护**: GameLink 前端团队
**最后更新**: 2025-01-28
**版本**: v1.0.0

如有开发相关问题，请联系前端团队或提交 Issue。
