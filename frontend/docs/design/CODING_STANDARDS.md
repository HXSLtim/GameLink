# GameLink Frontend 代码编写规范

## 📋 目录

1. [项目概述](#项目概述)
2. [命名规范](#命名规范)
3. [文件和目录结构](#文件和目录结构)
4. [TypeScript 规范](#typescript-规范)
5. [React 组件规范](#react-组件规范)
6. [样式编写规范](#样式编写规范)
7. [API 和状态管理规范](#api-和状态管理规范)
8. [代码格式化规范](#代码格式化规范)
9. [注释和文档规范](#注释和文档规范)
10. [测试规范](#测试规范)
11. [Git 提交规范](#git-提交规范)
12. [最佳实践](#最佳实践)

---

## 项目概述

**技术栈：**

- React 18.3+ (Hooks)
- TypeScript 5.6+
- Vite 5.4+ (构建工具)
- Arco Design (UI 组件库)
- React Router 6.27+ (路由管理)
- Vitest (单元测试)

**开发工具：**

- ESLint 9+ (代码检查)
- Prettier (代码格式化)
- TypeScript Compiler (类型检查)

---

## 命名规范

### 1. 文件命名

#### React 组件文件

- **PascalCase**（大驼峰）命名
- 扩展名使用 `.tsx`

```
✅ 推荐
UserProfile.tsx
GameList.tsx
LoginForm.tsx

❌ 避免
userProfile.tsx
user-profile.tsx
user_profile.tsx
```

#### 工具函数/Hook 文件

- **camelCase**（小驼峰）命名
- 扩展名使用 `.ts` 或 `.tsx`

```
✅ 推荐
formatDate.ts
useAuth.tsx
apiClient.ts

❌ 避免
FormatDate.ts
use-auth.tsx
api_client.ts
```

#### 类型定义文件

- **camelCase** 或描述性命名
- 使用 `.types.ts` 或 `.d.ts` 后缀

```
✅ 推荐
user.types.ts
api.types.ts
global.d.ts
```

#### 样式文件

- 与组件同名
- 使用 `.module.less` 或 `.less`

```
✅ 推荐
UserProfile.module.less
global.less
```

### 2. 变量命名

#### 常规变量

- 使用 **camelCase**
- 语义化命名，避免缩写

```typescript
✅ 推荐
const userName = 'John';
const isLoading = false;
const gameList = [];

❌ 避免
const user_name = 'John';
const loading = false;  // 不明确
const gl = [];  // 缩写不清晰
```

#### 常量

- 使用 **UPPER_SNAKE_CASE**（全大写下划线）
- 放在文件顶部或单独的常量文件中

```typescript
✅ 推荐
const API_BASE_URL = 'https://api.example.com';
const MAX_RETRY_COUNT = 3;
const DEFAULT_PAGE_SIZE = 20;

❌ 避免
const apiBaseUrl = 'https://api.example.com';
const maxRetryCount = 3;
```

#### 布尔值

- 使用 `is`、`has`、`should` 等前缀

```typescript
✅ 推荐
const isVisible = true;
const hasPermission = false;
const shouldUpdate = true;

❌ 避免
const visible = true;
const permission = false;
```

### 3. 函数命名

#### 事件处理函数

- 使用 `handle` 前缀

```typescript
✅ 推荐
const handleClick = () => {};
const handleSubmit = () => {};
const handleInputChange = () => {};

❌ 避免
const onClick = () => {};
const submit = () => {};
```

#### 普通函数

- 使用动词开头
- 清晰描述功能

```typescript
✅ 推荐
function fetchUserData() {}
function validateEmail() {}
function calculateTotal() {}

❌ 避免
function user() {}
function email() {}
function total() {}
```

### 4. 类型/接口命名

#### 接口（Interface）

- 使用 **PascalCase**
- 描述性命名，不使用 `I` 前缀

```typescript
✅ 推荐
interface User {
  id: string;
  name: string;
}

interface ApiResponse<T> {
  data: T;
  message: string;
}

❌ 避免
interface IUser {}  // 不使用 I 前缀
interface user {}   // 小写
```

#### 类型别名（Type）

- 使用 **PascalCase**

```typescript
✅ 推荐
type UserId = string;
type UserRole = 'admin' | 'user' | 'guest';
type UserState = {
  isLogged: boolean;
  token: string;
};
```

#### 枚举（Enum）

- 枚举名使用 **PascalCase**
- 枚举值使用 **PascalCase** 或 **UPPER_SNAKE_CASE**

```typescript
✅ 推荐
enum UserRole {
  Admin = 'admin',
  User = 'user',
  Guest = 'guest',
}

enum HttpStatus {
  OK = 200,
  NOT_FOUND = 404,
  INTERNAL_ERROR = 500,
}
```

### 5. 组件命名

#### 组件名

- 使用 **PascalCase**
- 功能描述清晰

```typescript
✅ 推荐
export const UserProfile = () => {};
export const GameListItem = () => {};
export const LoginForm = () => {};

❌ 避免
export const userprofile = () => {};
export const Game_List_Item = () => {};
```

#### 自定义 Hook

- 必须以 `use` 开头
- 使用 **camelCase**

```typescript
✅ 推荐
export const useAuth = () => {};
export const useFetchData = () => {};
export const useLocalStorage = () => {};

❌ 避免
export const auth = () => {};
export const fetchData = () => {};
```

---

## 文件和目录结构

### 目录结构说明

```
src/
├── api/              # API 请求定义
│   ├── user.ts
│   └── game.ts
├── components/       # 公共组件
│   ├── Button/
│   │   ├── Button.tsx
│   │   ├── Button.module.less
│   │   ├── Button.test.tsx
│   │   └── index.ts
│   └── Layout/
├── contexts/         # React Context
│   └── AuthContext.tsx
├── hooks/            # 自定义 Hooks
│   ├── useAuth.tsx
│   └── useFetch.tsx
├── layouts/          # 页面布局组件
│   └── MainLayout.tsx
├── pages/            # 页面组件
│   ├── Home/
│   │   ├── Home.tsx
│   │   ├── Home.module.less
│   │   └── index.ts
│   └── Login/
├── services/         # 业务逻辑服务
│   └── authService.ts
├── styles/           # 全局样式
│   └── global.less
├── types/            # 类型定义
│   ├── user.types.ts
│   └── api.types.ts
├── utils/            # 工具函数
│   ├── formatDate.ts
│   └── validation.ts
├── config.ts         # 配置文件
└── main.tsx          # 入口文件
```

### 目录规范

#### 1. 组件目录结构

每个组件应该有独立的文件夹：

```
ComponentName/
├── ComponentName.tsx          # 组件主文件
├── ComponentName.module.less  # 组件样式
├── ComponentName.test.tsx     # 组件测试
├── index.ts                   # 导出文件
└── types.ts                   # 组件类型定义（可选）
```

**index.ts 示例：**

```typescript
export { ComponentName } from './ComponentName';
export type { ComponentNameProps } from './ComponentName';
```

#### 2. 页面目录结构

页面组件遵循类似规则：

```
PageName/
├── PageName.tsx
├── PageName.module.less
├── components/              # 页面专属组件
│   └── PageSection.tsx
└── index.ts
```

#### 3. API 目录结构

按业务模块划分：

```
api/
├── user.ts      # 用户相关 API
├── game.ts      # 游戏相关 API
├── auth.ts      # 认证相关 API
└── index.ts     # 统一导出
```

---

## TypeScript 规范

### 1. 类型定义

#### 优先使用 interface 定义对象类型

```typescript
✅ 推荐
interface User {
  id: string;
  name: string;
  email: string;
}

// 当需要联合类型或其他高级类型时使用 type
type UserRole = 'admin' | 'user' | 'guest';
type UserWithRole = User & { role: UserRole };
```

#### 避免使用 any

```typescript
✅ 推荐
interface ApiResponse<T> {
  data: T;
  status: number;
}

function fetchData<T>(url: string): Promise<ApiResponse<T>> {
  // ...
}

❌ 避免
function fetchData(url: string): Promise<any> {
  // ...
}
```

#### 使用联合类型和枚举

```typescript
✅ 推荐
type Status = 'pending' | 'success' | 'error';

enum UserRole {
  Admin = 'admin',
  User = 'user',
  Guest = 'guest',
}

❌ 避免
type Status = string;  // 太宽泛
```

### 2. 函数类型

#### 明确定义函数参数和返回值类型

```typescript
✅ 推荐
function calculateTotal(items: number[], tax: number): number {
  return items.reduce((sum, item) => sum + item, 0) * (1 + tax);
}

const handleSubmit = (event: React.FormEvent<HTMLFormElement>): void => {
  event.preventDefault();
  // ...
};

❌ 避免
function calculateTotal(items, tax) {
  return items.reduce((sum, item) => sum + item, 0) * (1 + tax);
}
```

#### 使用可选参数和默认值

```typescript
✅ 推荐
function fetchUsers(page: number = 1, pageSize: number = 20): Promise<User[]> {
  // ...
}

interface Config {
  timeout?: number;
  retry?: boolean;
}
```

### 3. 组件 Props 类型

#### 定义清晰的 Props 接口

```typescript
✅ 推荐
interface ButtonProps {
  text: string;
  onClick: () => void;
  disabled?: boolean;
  variant?: 'primary' | 'secondary';
}

export const Button: React.FC<ButtonProps> = ({
  text,
  onClick,
  disabled = false,
  variant = 'primary',
}) => {
  return (
    <button onClick={onClick} disabled={disabled} className={variant}>
      {text}
    </button>
  );
};

❌ 避免
export const Button = ({ text, onClick, disabled, variant }: any) => {
  // ...
};
```

#### 使用泛型组件

```typescript
✅ 推荐
interface ListProps<T> {
  items: T[];
  renderItem: (item: T) => React.ReactNode;
}

export const List = <T,>({ items, renderItem }: ListProps<T>) => {
  return (
    <ul>
      {items.map((item, index) => (
        <li key={index}>{renderItem(item)}</li>
      ))}
    </ul>
  );
};
```

### 4. 类型导入导出

#### 使用 type 关键字导入类型

```typescript
✅ 推荐
import type { User, UserRole } from './types/user.types';
import { fetchUser } from './api/user';

❌ 不推荐（但可接受）
import { User, UserRole, fetchUser } from './user';
```

#### 统一导出类型

```typescript
✅ 推荐
// types/index.ts
export type { User, UserProfile } from './user.types';
export type { Game, GameCategory } from './game.types';
export type { ApiResponse, ApiError } from './api.types';
```

---

## React 组件规范

### 1. 组件定义

#### 使用函数组件和 Hooks

```typescript
✅ 推荐
export const UserProfile: React.FC<UserProfileProps> = ({ userId }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    fetchUserData(userId);
  }, [userId]);

  return <div>{/* ... */}</div>;
};

❌ 避免
class UserProfile extends React.Component {
  // 不推荐使用类组件
}
```

#### 组件结构顺序

```typescript
export const MyComponent: React.FC<MyComponentProps> = (props) => {
  // 1. Props 解构
  const { title, onSubmit } = props;

  // 2. Hooks
  const [state, setState] = useState();
  const context = useContext(MyContext);
  const ref = useRef();

  // 3. 派生状态和计算值
  const computedValue = useMemo(() => {
    return expensiveCalculation(state);
  }, [state]);

  // 4. 事件处理函数
  const handleClick = () => {
    // ...
  };

  // 5. 副作用
  useEffect(() => {
    // ...
  }, []);

  // 6. 提前返回（条件渲染）
  if (isLoading) return <LoadingSpinner />;
  if (error) return <ErrorMessage error={error} />;

  // 7. 主渲染
  return (
    <div>
      {/* ... */}
    </div>
  );
};
```

### 2. Props 规范

#### Props 解构

```typescript
✅ 推荐
interface ButtonProps {
  text: string;
  onClick: () => void;
  disabled?: boolean;
}

export const Button: React.FC<ButtonProps> = ({ text, onClick, disabled = false }) => {
  return (
    <button onClick={onClick} disabled={disabled}>
      {text}
    </button>
  );
};

❌ 避免
export const Button: React.FC<ButtonProps> = (props) => {
  return (
    <button onClick={props.onClick} disabled={props.disabled}>
      {props.text}
    </button>
  );
};
```

#### Children Props

```typescript
✅ 推荐
interface CardProps {
  title: string;
  children: React.ReactNode;
}

export const Card: React.FC<CardProps> = ({ title, children }) => {
  return (
    <div className="card">
      <h2>{title}</h2>
      <div className="card-body">{children}</div>
    </div>
  );
};
```

### 3. 状态管理

#### 使用 useState

```typescript
✅ 推荐
const [count, setCount] = useState<number>(0);
const [user, setUser] = useState<User | null>(null);
const [filters, setFilters] = useState<Filters>({ status: 'all', page: 1 });

// 更新对象状态使用展开运算符
setFilters((prev) => ({ ...prev, page: prev.page + 1 }));

❌ 避免
const [state, setState] = useState({ count: 0, user: null, filters: {} });
// 避免将不相关的状态放在一起
```

#### 使用 useReducer（复杂状态）

```typescript
✅ 推荐
interface State {
  data: User[];
  isLoading: boolean;
  error: string | null;
}

type Action =
  | { type: 'FETCH_START' }
  | { type: 'FETCH_SUCCESS'; payload: User[] }
  | { type: 'FETCH_ERROR'; payload: string };

const reducer = (state: State, action: Action): State => {
  switch (action.type) {
    case 'FETCH_START':
      return { ...state, isLoading: true, error: null };
    case 'FETCH_SUCCESS':
      return { ...state, isLoading: false, data: action.payload };
    case 'FETCH_ERROR':
      return { ...state, isLoading: false, error: action.payload };
    default:
      return state;
  }
};

export const UserList = () => {
  const [state, dispatch] = useReducer(reducer, {
    data: [],
    isLoading: false,
    error: null,
  });

  // ...
};
```

### 4. 副作用（useEffect）

#### 明确依赖项

```typescript
✅ 推荐
useEffect(() => {
  fetchUserData(userId);
}, [userId]); // 明确依赖

// 仅执行一次
useEffect(() => {
  initializeApp();
}, []);

❌ 避免
useEffect(() => {
  fetchUserData(userId);
}); // 缺少依赖数组，每次渲染都执行
```

#### 清理副作用

```typescript
✅ 推荐
useEffect(() => {
  const timer = setTimeout(() => {
    console.log('Delayed action');
  }, 1000);

  return () => {
    clearTimeout(timer);
  };
}, []);

useEffect(() => {
  const subscription = eventBus.subscribe('event', handler);

  return () => {
    subscription.unsubscribe();
  };
}, []);
```

### 5. 自定义 Hooks

#### 提取可复用逻辑

```typescript
✅ 推荐
export const useFetch = <T>(url: string) => {
  const [data, setData] = useState<T | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      setIsLoading(true);
      try {
        const response = await fetch(url);
        const result = await response.json();
        setData(result);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unknown error');
      } finally {
        setIsLoading(false);
      }
    };

    fetchData();
  }, [url]);

  return { data, isLoading, error };
};

// 使用
const UserProfile = ({ userId }: { userId: string }) => {
  const { data: user, isLoading, error } = useFetch<User>(`/api/users/${userId}`);

  // ...
};
```

### 6. 性能优化

#### 使用 useMemo 和 useCallback

```typescript
✅ 推荐
export const ExpensiveComponent = ({ items, onItemClick }: Props) => {
  // 缓存计算结果
  const sortedItems = useMemo(() => {
    return [...items].sort((a, b) => a.name.localeCompare(b.name));
  }, [items]);

  // 缓存回调函数
  const handleItemClick = useCallback(
    (id: string) => {
      onItemClick(id);
    },
    [onItemClick]
  );

  return (
    <ul>
      {sortedItems.map((item) => (
        <ListItem key={item.id} item={item} onClick={handleItemClick} />
      ))}
    </ul>
  );
};

// 子组件使用 React.memo 避免不必要的重渲染
export const ListItem = React.memo<ListItemProps>(({ item, onClick }) => {
  return <li onClick={() => onClick(item.id)}>{item.name}</li>;
});
```

### 7. 条件渲染

#### 简洁的条件渲染

```typescript
✅ 推荐
// 使用 && 运算符
{isLoggedIn && <UserMenu />}

// 使用三元运算符
{isLoading ? <Spinner /> : <Content />}

// 提前返回
if (isLoading) return <Spinner />;
if (error) return <ErrorMessage error={error} />;
return <Content />;

❌ 避免
// 不要使用复杂的嵌套三元运算符
{isLoading ? <Spinner /> : error ? <ErrorMessage /> : data ? <Content /> : null}
```

---

## 样式编写规范

### 1. CSS Modules

#### 使用 CSS Modules

```typescript
✅ 推荐
// UserProfile.module.less
.container {
  padding: 20px;
  background: #fff;
}

.title {
  font-size: 24px;
  font-weight: bold;
}

.avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
}

// UserProfile.tsx
import styles from './UserProfile.module.less';

export const UserProfile = () => {
  return (
    <div className={styles.container}>
      <h1 className={styles.title}>Profile</h1>
      <img src="..." className={styles.avatar} />
    </div>
  );
};
```

### 2. 类名规范

#### BEM 命名约定（推荐但不强制）

```less
// 推荐
.card {
  // Block
  &__header {
    // Element
  }

  &__body {
    // Element
  }

  &--featured {
    // Modifier
  }
}

// 使用
<div className={styles.card}>
  <div className={styles.card__header}>Header</div>
  <div className={styles.card__body}>Body</div>
</div>
```

### 3. Arco Design 样式定制

#### 使用 Arco Design 主题变量

```less
// 使用 Arco Design 的设计令牌
@import '@arco-design/web-react/es/style/index.less';

.customButton {
  color: var(--color-primary);
  background: var(--color-bg-2);
  border-radius: var(--border-radius-small);
}
```

### 4. 响应式设计

#### 使用媒体查询

```less
.container {
  padding: 20px;

  @media (max-width: 768px) {
    padding: 10px;
  }

  @media (max-width: 480px) {
    padding: 5px;
  }
}
```

---

## API 和状态管理规范

### 1. API 请求

#### 统一的 API 客户端

```typescript
✅ 推荐
// api/client.ts
import { API_BASE_URL } from '../config';

interface ApiResponse<T> {
  data: T;
  message: string;
  status: number;
}

class ApiClient {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    const url = `${this.baseURL}${endpoint}`;
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    });

    if (!response.ok) {
      throw new Error(`API Error: ${response.statusText}`);
    }

    return response.json();
  }

  get<T>(endpoint: string) {
    return this.request<T>(endpoint, { method: 'GET' });
  }

  post<T>(endpoint: string, data: unknown) {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  put<T>(endpoint: string, data: unknown) {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  delete<T>(endpoint: string) {
    return this.request<T>(endpoint, { method: 'DELETE' });
  }
}

export const apiClient = new ApiClient(API_BASE_URL);
```

#### API 模块化

```typescript
✅ 推荐
// api/user.ts
import { apiClient } from './client';
import type { User, UserCreateDto, UserUpdateDto } from '../types/user.types';

export const userApi = {
  getUsers: () => apiClient.get<User[]>('/users'),

  getUserById: (id: string) => apiClient.get<User>(`/users/${id}`),

  createUser: (data: UserCreateDto) => apiClient.post<User>('/users', data),

  updateUser: (id: string, data: UserUpdateDto) =>
    apiClient.put<User>(`/users/${id}`, data),

  deleteUser: (id: string) => apiClient.delete(`/users/${id}`),
};
```

### 2. Context 使用

#### 创建 Context

```typescript
✅ 推荐
// contexts/AuthContext.tsx
import { createContext, useContext, useState, ReactNode } from 'react';
import type { User } from '../types/user.types';

interface AuthContextValue {
  user: User | null;
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);

  const login = async (email: string, password: string) => {
    // 登录逻辑
    const user = await authService.login(email, password);
    setUser(user);
  };

  const logout = () => {
    setUser(null);
  };

  const value: AuthContextValue = {
    user,
    isAuthenticated: user !== null,
    login,
    logout,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

// 自定义 Hook
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
};
```

---

## 代码格式化规范

### Prettier 配置

项目使用 Prettier 进行代码格式化，配置如下：

```json
{
  "singleQuote": true,
  "semi": true,
  "trailingComma": "all",
  "printWidth": 100
}
```

**说明：**

- `singleQuote`: 使用单引号
- `semi`: 使用分号
- `trailingComma`: 在多行结构中使用尾随逗号
- `printWidth`: 每行最大 100 字符

### ESLint 规则

主要规则：

- 未使用的变量会警告（以 `_` 开头的除外）
- 不允许使用 `any` 类型（已关闭，但不推荐使用）
- React Hooks 规则严格执行
- React 17+ 不需要导入 React

### 格式化命令

```bash
# 格式化所有文件
npm run format

# 检查代码规范
npm run lint

# 类型检查
npm run typecheck
```

---

## 注释和文档规范

### 1. 文件头部注释

```typescript
/**
 * 用户管理 API
 *
 * @module api/user
 * @description 提供用户 CRUD 操作的 API 接口
 */
```

### 2. 函数注释

````typescript
✅ 推荐
/**
 * 计算购物车总价
 *
 * @param items - 购物车商品列表
 * @param taxRate - 税率（0-1 之间）
 * @returns 包含税费的总价
 *
 * @example
 * ```typescript
 * const total = calculateCartTotal([{ price: 100 }, { price: 200 }], 0.1);
 * // 返回: 330
 * ```
 */
export function calculateCartTotal(
  items: CartItem[],
  taxRate: number
): number {
  const subtotal = items.reduce((sum, item) => sum + item.price, 0);
  return subtotal * (1 + taxRate);
}
````

### 3. 组件注释

````typescript
✅ 推荐
/**
 * 用户资料卡片组件
 *
 * @component
 * @example
 * ```tsx
 * <UserCard
 *   user={user}
 *   onEdit={(id) => console.log('Edit user:', id)}
 * />
 * ```
 */
export const UserCard: React.FC<UserCardProps> = ({ user, onEdit }) => {
  // ...
};
````

### 4. 行内注释

```typescript
✅ 推荐
// 获取用户权限列表
const permissions = await fetchUserPermissions(userId);

// TODO: 添加缓存机制以提高性能
const data = await fetchData();

// FIXME: 处理边界情况 - 空数组
const firstItem = items[0];

❌ 避免
// 这是一个变量
const x = 10;  // 不要写显而易见的注释
```

### 5. 类型注释

```typescript
✅ 推荐
/**
 * 用户信息
 */
interface User {
  /** 用户唯一标识 */
  id: string;

  /** 用户名 */
  name: string;

  /** 电子邮件地址 */
  email: string;

  /** 用户角色 */
  role: UserRole;

  /** 账户创建时间 */
  createdAt: Date;
}
```

---

## 测试规范

### 1. 测试文件命名

```
Component.test.tsx
utils.test.ts
useCustomHook.test.tsx
```

### 2. 组件测试

```typescript
✅ 推荐
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { Button } from './Button';

describe('Button', () => {
  it('should render with text', () => {
    render(<Button text="Click me" onClick={() => {}} />);
    expect(screen.getByText('Click me')).toBeInTheDocument();
  });

  it('should call onClick when clicked', () => {
    const handleClick = vi.fn();
    render(<Button text="Click me" onClick={handleClick} />);

    fireEvent.click(screen.getByText('Click me'));
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it('should be disabled when disabled prop is true', () => {
    render(<Button text="Click me" onClick={() => {}} disabled />);
    expect(screen.getByText('Click me')).toBeDisabled();
  });
});
```

### 3. Hook 测试

```typescript
import { renderHook, waitFor } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { useFetch } from './useFetch';

describe('useFetch', () => {
  it('should fetch data successfully', async () => {
    const { result } = renderHook(() => useFetch<User>('/api/user/1'));

    expect(result.current.isLoading).toBe(true);

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    expect(result.current.data).toBeDefined();
    expect(result.current.error).toBeNull();
  });
});
```

### 4. 测试覆盖率

运行测试并生成覆盖率报告：

```bash
npm run test -- --coverage
```

目标覆盖率：

- 语句覆盖率（Statements）：> 80%
- 分支覆盖率（Branches）：> 75%
- 函数覆盖率（Functions）：> 80%
- 行覆盖率（Lines）：> 80%

---

## Git 提交规范

### Commit Message 格式

使用 Conventional Commits 规范：

```
<type>(<scope>): <subject>

<body>

<footer>
```

#### Type 类型

- `feat`: 新功能
- `fix`: 修复 Bug
- `docs`: 文档更新
- `style`: 代码格式（不影响代码运行）
- `refactor`: 重构（既不是新功能也不是修复 Bug）
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动
- `ci`: CI/CD 相关

#### 示例

```bash
✅ 推荐
feat(auth): 添加用户登录功能

实现了基于 JWT 的用户登录认证系统
- 添加登录表单组件
- 实现 API 请求和响应处理
- 添加 token 存储和验证

Closes #123

fix(user): 修复用户列表分页问题

当页码超出范围时，重置为第一页

refactor(api): 优化 API 请求错误处理

统一错误处理逻辑，提高代码可维护性

❌ 避免
update code
fix bug
add feature
```

### 分支管理

```
main          - 主分支，生产环境代码
develop       - 开发分支
feature/*     - 功能分支（如 feature/user-login）
bugfix/*      - Bug 修复分支（如 bugfix/fix-pagination）
hotfix/*      - 紧急修复分支
release/*     - 发布分支
```

---

## 最佳实践

### 1. 代码复用

```typescript
✅ 推荐
// 提取可复用的自定义 Hook
export const useLocalStorage = <T>(key: string, initialValue: T) => {
  const [value, setValue] = useState<T>(() => {
    const stored = localStorage.getItem(key);
    return stored ? JSON.parse(stored) : initialValue;
  });

  useEffect(() => {
    localStorage.setItem(key, JSON.stringify(value));
  }, [key, value]);

  return [value, setValue] as const;
};
```

### 2. 错误处理

```typescript
✅ 推荐
// 统一的错误处理
try {
  const data = await fetchData();
  setData(data);
} catch (error) {
  if (error instanceof ApiError) {
    setError(error.message);
  } else if (error instanceof Error) {
    setError(error.message);
  } else {
    setError('An unknown error occurred');
  }
}

// 使用 Error Boundary
export class ErrorBoundary extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('Error caught by boundary:', error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return <ErrorFallback />;
    }

    return this.props.children;
  }
}
```

### 3. 性能优化

```typescript
✅ 推荐
// 懒加载路由
import { lazy, Suspense } from 'react';

const Home = lazy(() => import('./pages/Home'));
const UserProfile = lazy(() => import('./pages/UserProfile'));

export const App = () => {
  return (
    <Suspense fallback={<LoadingSpinner />}>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/profile" element={<UserProfile />} />
      </Routes>
    </Suspense>
  );
};

// 使用虚拟列表（长列表）
// 推荐使用 react-window 或 react-virtualized
```

### 4. 可访问性（A11y）

```typescript
✅ 推荐
// 添加适当的 ARIA 属性
<button
  aria-label="关闭对话框"
  onClick={handleClose}
>
  <CloseIcon />
</button>

<input
  type="text"
  id="username"
  aria-required="true"
  aria-describedby="username-help"
/>
<span id="username-help">请输入您的用户名</span>
```

### 5. 安全性

```typescript
✅ 推荐
// 防止 XSS 攻击 - React 默认会转义
<div>{userInput}</div>  // 安全

// 避免使用 dangerouslySetInnerHTML
❌ 避免
<div dangerouslySetInnerHTML={{ __html: userInput }} />

// 如果必须使用，先进行清理
import DOMPurify from 'dompurify';
<div dangerouslySetInnerHTML={{ __html: DOMPurify.sanitize(userInput) }} />
```

### 6. 环境变量

```typescript
✅ 推荐
// vite.config.ts 或 .env
// Vite 环境变量必须以 VITE_ 开头
VITE_API_BASE_URL=https://api.example.com
VITE_APP_NAME=GameLink

// 使用
const apiUrl = import.meta.env.VITE_API_BASE_URL;
```

### 7. 代码分割

```typescript
✅ 推荐
// 动态导入
const loadHeavyModule = async () => {
  const module = await import('./heavyModule');
  return module.default;
};

// 路由级代码分割
const routes = [
  {
    path: '/dashboard',
    component: lazy(() => import('./pages/Dashboard')),
  },
  {
    path: '/settings',
    component: lazy(() => import('./pages/Settings')),
  },
];
```

---

## 附录

### 推荐工具和插件

#### VS Code 插件

- ESLint
- Prettier - Code formatter
- TypeScript Vue Plugin (Volar)
- Auto Rename Tag
- Path Intellisense
- Import Cost

#### 浏览器扩展

- React Developer Tools
- Redux DevTools（如果使用 Redux）

### 参考资源

- [React 官方文档](https://react.dev/)
- [TypeScript 官方文档](https://www.typescriptlang.org/)
- [Arco Design 文档](https://arco.design/)
- [Vite 文档](https://vitejs.dev/)
- [React Testing Library](https://testing-library.com/react)

---

## 更新日志

- **2025-10-27**: 初始版本发布

---

**注意：** 本规范是团队共识的结晶，请所有团队成员遵守。如有疑问或建议，请及时与团队沟通。
