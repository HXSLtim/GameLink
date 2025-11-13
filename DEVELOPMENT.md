# 🛠️ GameLink 开发指南

本文档提供 GameLink 项目的完整开发环境搭建、编码规范和开发流程指南。

---

## 📋 目录

- [环境准备](#环境准备)
- [开发环境搭建](#开发环境搭建)
- [项目结构详解](#项目结构详解)
- [编码规范](#编码规范)
- [开发流程](#开发流程)
- [测试指南](#测试指南)
- [调试技巧](#调试技巧)
- [性能优化](#性能优化)
- [常见问题](#常见问题)

---

## 🔧 环境准备

### 系统要求
- **操作系统**: Windows 10+, macOS 10.15+, Ubuntu 18.04+
- **内存**: 8GB+ (推荐 16GB)
- **存储**: 20GB+ 可用空间

### 必需软件

#### 1. Go 语言环境
```bash
# 下载并安装 Go 1.25.3+
# Windows: 从 https://golang.org/dl/ 下载安装包
# macOS: brew install go
# Ubuntu: sudo apt install golang-go

# 验证安装
go version
# 应输出: go version go1.25.3 linux/amd64

# 设置环境变量
export GOPATH=$HOME/go
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
```

#### 2. Node.js 环境
```bash
# 使用 nvm 安装 Node.js 18+
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
nvm install 18
nvm use 18

# 验证安装
node --version  # v18.x.x
npm --version   # 9.x.x
```

#### 3. 数据库环境

**MySQL 8.0+**
```bash
# Windows: 下载 MySQL Installer
# macOS: brew install mysql
# Ubuntu: sudo apt install mysql-server

# 启动服务
brew services start mysql  # macOS
sudo systemctl start mysql # Ubuntu

# 创建数据库
mysql -u root -p
CREATE DATABASE gamelink CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'gamelink'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON gamelink.* TO 'gamelink'@'localhost';
FLUSH PRIVILEGES;
```

**Redis 6.0+**
```bash
# Windows: 下载 Redis for Windows
# macOS: brew install redis
# Ubuntu: sudo apt install redis-server

# 启动服务
brew services start redis  # macOS
sudo systemctl start redis # Ubuntu
```

#### 4. 开发工具

**推荐 IDE**: VS Code
```bash
# 安装 VS Code 扩展
code --install-extension golang.go
code --install-extension ms-vscode.vscode-typescript-next
code --install-extension bradlc.vscode-tailwindcss
code --install-extension ms-vscode.vscode-json
```

**其他工具**
```bash
# Go 工具链
go install -a github.com/cweill/gotests/gotests@latest
go install -a github.com/fatih/gomodifytags@latest
go install -a github.com/josharian/impl@latest
go install -a github.com/haya14busa/goplay/cmd/goplay@latest
go install -a github.com/go-delve/delve/cmd/dlv@latest
go install -a honnef.co/go/tools/cmd/staticcheck@latest
go install -a golang.org/x/tools/cmd/goimports@latest
go install -a golang.org/x/tools/cmd/godoc@latest
go install -a github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 前端工具
npm install -g nodemon
npm install -g concurrently
```

---

## 🏗️ 开发环境搭建

### 1. 克隆项目
```bash
git clone https://github.com/your-org/GameLink.git
cd GameLink
```

### 2. 后端环境配置

```bash
cd backend

# 安装依赖
go mod download
go mod tidy

# 复制配置文件
cp configs/config.example.yaml configs/config.yaml
# 编辑配置文件，填入数据库连接信息

# 运行数据库迁移
make migrate
# 或手动执行
go run scripts/migrate/main.go up

# 生成测试数据 (可选)
go run scripts/seed/main.go
```

**配置文件示例 (`configs/config.yaml`)**
```yaml
server:
  port: 8080
  mode: debug  # debug, release

database:
  host: localhost
  port: 3306
  user: gamelink
  password: your_password
  database: gamelink
  charset: utf8mb4

redis:
  host: localhost
  port: 6379
  password: ""
  database: 0

jwt:
  secret: your-jwt-secret-key
  expire_hours: 24

upload:
  max_size: 10485760  # 10MB
  allowed_types: ["jpg", "jpeg", "png", "gif"]
  path: "./uploads"

log:
  level: info
  file: "./logs/app.log"
  max_size: 100
  max_backups: 3
```

### 3. 前端环境配置

```bash
cd frontend

# 安装依赖
npm install

# 复制环境配置
cp .env.example .env.local
# 编辑配置文件
```

**环境配置文件 (`.env.local`)**
```env
# API 配置
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_WS_URL=ws://localhost:8080/ws

# 应用配置
VITE_APP_NAME=GameLink
VITE_APP_VERSION=2.1.0

# 第三方服务 (可选)
VITE_GOOGLE_CLIENT_ID=your_google_client_id
VITE_WECHAT_APP_ID=your_wechat_app_id
```

### 4. 启动开发服务

**方式一：分别启动**
```bash
# 终端 1 - 启动后端
cd backend
make run CMD=user-service

# 终端 2 - 启动前端
cd frontend
npm run dev
```

**方式二：同时启动**
```bash
# 在项目根目录
npm run dev
# 或
make dev
```

### 5. 验证环境
访问以下地址验证环境是否搭建成功：
- 前端应用: http://localhost:5173
- 后端健康检查: http://localhost:8080/health
- API 文档: http://localhost:8080/swagger/index.html

---

## 📁 项目结构详解

### 后端结构
```
backend/
├── cmd/                           # 应用入口
│   └── user-service/
│       ├── main.go               # 主程序入口
│       └── wire.go               # 依赖注入配置
├── internal/                     # 内部包 (不对外暴露)
│   ├── admin/                   # 管理端处理器
│   │   ├── handler/             # HTTP 处理器
│   │   ├── service/             # 业务逻辑
│   │   └── repository/          # 数据访问
│   ├── handler/                 # HTTP 处理器
│   │   ├── admin/               # 管理端接口
│   │   ├── user/                # 用户端接口
│   │   ├── player/              # 陪玩师接口
│   │   ├── middleware/          # 中间件
│   │   └── websocket/           # WebSocket 处理
│   ├── service/                 # 业务逻辑层
│   │   ├── user.go
│   │   ├── order.go
│   │   ├── payment.go
│   │   └── ...
│   ├── repository/              # 数据访问层
│   │   ├── interfaces.go        # 接口定义
│   │   ├── user_repo.go
│   │   ├── order_repo.go
│   │   └── ...
│   ├── model/                   # 数据模型
│   │   ├── user.go
│   │   ├── order.go
│   │   └── ...
│   ├── auth/                    # 认证授权
│   │   ├── jwt.go
│   │   ├── rbac.go
│   │   └── middleware.go
│   ├── cache/                   # 缓存层
│   │   ├── redis.go
│   │   └── memory.go
│   ├── config/                  # 配置管理
│   │   └── config.go
│   ├── utils/                   # 工具函数
│   │   ├── response.go
│   │   ├── validation.go
│   │   └── ...
│   └── websocket/               # WebSocket 处理
│       ├── hub.go
│       └── client.go
├── pkg/                         # 可复用的公共包
│   ├── logger/
│   ├── validator/
│   └── ...
├── configs/                     # 配置文件
│   ├── config.yaml
│   └── config.example.yaml
├── docs/                        # 文档
│   └── swagger/
├── scripts/                     # 脚本文件
│   ├── migrate/
│   ├── seed/
│   └── build/
├── tests/                       # 测试文件
│   ├── integration/
│   └── e2e/
├── uploads/                     # 文件上传目录
├── logs/                        # 日志文件
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### 前端结构
```
frontend/
├── public/                      # 静态资源
│   ├── index.html
│   ├── favicon.ico
│   └── ...
├── src/
│   ├── api/                     # API 调用层
│   │   ├── auth.ts
│   │   ├── user.ts
│   │   ├── order.ts
│   │   └── ...
│   ├── components/              # 可复用组件
│   │   ├── common/              # 通用组件
│   │   │   ├── Button/
│   │   │   ├── Modal/
│   │   │   ├── Table/
│   │   │   └── ...
│   │   ├── chat/                # 聊天组件
│   │   ├── order/               # 订单组件
│   │   └── user/                # 用户组件
│   ├── pages/                   # 页面组件
│   │   ├── admin/               # 管理端页面
│   │   │   ├── Dashboard/
│   │   │   ├── UserManagement/
│   │   │   └── ...
│   │   ├── user/                # 用户端页面
│   │   │   ├── Home/
│   │   │   ├── Profile/
│   │   │   └── ...
│   │   └── player/              # 陪玩师页面
│   │       ├── Dashboard/
│   │       └── ...
│   ├── layouts/                 # 布局组件
│   │   ├── AdminLayout/
│   │   ├── UserLayout/
│   │   └── PlayerLayout/
│   ├── hooks/                   # 自定义 Hooks
│   │   ├── useAuth.ts
│   │   ├── useWebSocket.ts
│   │   └── ...
│   ├── store/                   # 状态管理
│   │   ├── auth.ts
│   │   ├── order.ts
│   │   └── ...
│   ├── types/                   # TypeScript 类型
│   │   ├── api.ts
│   │   ├── user.ts
│   │   └── ...
│   ├── utils/                   # 工具函数
│   │   ├── request.ts
│   │   ├── storage.ts
│   │   └── ...
│   ├── styles/                  # 样式文件
│   │   ├── globals.less
│   │   ├── variables.less
│   │   └── ...
│   ├── App.tsx
│   ├── main.tsx
│   └── vite-env.d.ts
├── docs/                        # 前端文档
├── .env.example
├── .gitignore
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
└── README.md
```

---

## 📝 编码规范

### Go 编码规范

#### 1. 命名规范
```go
// 包名：小写，简短，有意义
package user

// 常量：大写，下划线分隔
const MAX_RETRY_COUNT = 3

// 变量：驼峰命名
var userService UserService

// 函数：驼峰命名，导出函数首字母大写
func CreateUser(user *User) error { }
func validateUser(user *User) bool { }

// 结构体：驼峰命名，导出结构体首字母大写
type UserService struct {
    repo UserRepository
}

// 接口：通常以 -er 结尾
type UserRepository interface {
    Create(user *User) error
    GetByID(id int64) (*User, error)
}
```

#### 2. 注释规范
```go
// UserService 用户服务层
// 提供用户相关的业务逻辑处理
type UserService struct {
    repo UserRepository
    cache Cache
}

// CreateUser 创建新用户
// 参数:
//   - user: 用户信息
// 返回值:
//   - error: 错误信息
func (s *UserService) CreateUser(user *User) error {
    // 参数验证
    if err := s.validateUser(user); err != nil {
        return fmt.Errorf("用户验证失败: %w", err)
    }

    // 业务逻辑处理
    return s.repo.Create(user)
}
```

#### 3. 错误处理
```go
// 使用 fmt.Errorf 包装错误
if err != nil {
    return fmt.Errorf("创建用户失败: %w", err)
}

// 使用自定义错误类型
var (
    ErrUserNotFound = errors.New("用户不存在")
    ErrInvalidInput = errors.New("输入参数无效")
)

// 错误处理最佳实践
user, err := s.repo.GetByID(id)
if err != nil {
    if errors.Is(err, repository.ErrNotFound) {
        return nil, ErrUserNotFound
    }
    return nil, fmt.Errorf("获取用户失败: %w", err)
}
```

#### 4. 并发安全
```go
// 使用互斥锁保护共享资源
type Counter struct {
    mu    sync.RWMutex
    value int64
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *Counter) Value() int64 {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.value
}
```

### TypeScript 编码规范

#### 1. 类型定义
```typescript
// 使用 interface 定义对象类型
interface User {
  id: number;
  username: string;
  email: string;
  createdAt: Date;
}

// 使用 type 定义联合类型或复杂类型
type UserRole = 'admin' | 'user' | 'player';
type ApiResponse<T> = {
  success: boolean;
  data: T;
  message?: string;
};

// 使用泛型
interface Repository<T> {
  create(data: T): Promise<T>;
  findById(id: number): Promise<T | null>;
  update(id: number, data: Partial<T>): Promise<T>;
}
```

#### 2. 函数定义
```typescript
// 箭头函数，明确指定参数和返回值类型
const fetchUser = async (id: number): Promise<User | null> => {
  try {
    const response = await api.get<ApiResponse<User>>(`/users/${id}`);
    return response.data.data;
  } catch (error) {
    console.error('获取用户失败:', error);
    return null;
  }
};

// 函数重载
function formatDate(date: Date): string;
function formatDate(date: string): string;
function formatDate(date: Date | string): string {
  const d = typeof date === 'string' ? new Date(date) : date;
  return d.toLocaleDateString();
}
```

#### 3. 组件定义
```typescript
// React 函数组件
interface UserCardProps {
  user: User;
  onEdit?: (user: User) => void;
  className?: string;
}

const UserCard: React.FC<UserCardProps> = ({
  user,
  onEdit,
  className
}) => {
  return (
    <div className={className}>
      <h3>{user.username}</h3>
      <p>{user.email}</p>
      {onEdit && (
        <button onClick={() => onEdit(user)}>
          编辑
        </button>
      )}
    </div>
  );
};

// 使用泛型组件
interface ListProps<T> {
  items: T[];
  renderItem: (item: T) => React.ReactNode;
  loading?: boolean;
}

function List<T>({ items, renderItem, loading }: ListProps<T>) {
  if (loading) return <div>加载中...</div>;

  return (
    <div>
      {items.map(renderItem)}
    </div>
  );
}
```

#### 4. 自定义 Hooks
```typescript
// 自定义 Hook 必须以 use 开头
interface UseApiResult<T> {
  data: T | null;
  loading: boolean;
  error: string | null;
  refetch: () => Promise<void>;
}

function useApi<T>(url: string): UseApiResult<T> {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchData = useCallback(async () => {
    try {
      setLoading(true);
      const response = await api.get<T>(url);
      setData(response.data);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : '未知错误');
    } finally {
      setLoading(false);
    }
  }, [url]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return { data, loading, error, refetch: fetchData };
}
```

---

## 🔄 开发流程

### 1. Git 工作流

```bash
# 1. 创建功能分支
git checkout -b feature/user-authentication

# 2. 开发并提交
git add .
git commit -m "feat: 添加用户认证功能

- 实现 JWT 登录
- 添加密码加密
- 完善错误处理

Closes #123"

# 3. 推送分支
git push origin feature/user-authentication

# 4. 创建 Pull Request
# 填写 PR 模板，等待代码审查

# 5. 合并后清理分支
git checkout main
git pull origin main
git branch -d feature/user-authentication
```

### 2. 提交规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**类型说明：**
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建或辅助工具变动

**示例：**
```bash
feat(auth): add JWT authentication
fix(order): resolve payment calculation error
docs(api): update user endpoints documentation
```

### 3. 代码审查清单

#### 后端审查要点
- [ ] 代码是否符合 Go 编码规范
- [ ] 是否有适当的错误处理
- [ ] 是否有必要的单元测试
- [ ] 是否有安全漏洞（SQL注入、XSS等）
- [ ] 是否有性能问题
- [ ] API 接口是否符合 RESTful 规范
- [ ] 数据库操作是否使用事务
- [ ] 是否有并发安全问题

#### 前端审查要点
- [ ] 代码是否符合 TypeScript 规范
- [ ] 组件是否可复用
- [ ] 是否有适当的错误边界
- [ ] 是否有性能优化（懒加载、防抖等）
- [ ] 是否有安全漏洞（XSS、CSRF等）
- [ ] 用户体验是否良好
- [ ] 是否有适当的测试用例
- [ ] 是否有无障碍访问支持

### 4. 测试流程

```bash
# 后端测试
cd backend

# 运行所有测试
make test

# 运行特定测试
go test ./internal/service/...

# 生成覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# 前端测试
cd frontend

# 运行单元测试
npm run test

# 运行 E2E 测试
npm run test:e2e

# 生成覆盖率报告
npm run test:coverage
```

---

## 🧪 测试指南

### 后端测试

#### 1. 单元测试
```go
// user_service_test.go
package service

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/suite"
)

// 使用 testify 进行测试
type UserServiceTestSuite struct {
    suite.Suite
    service *UserService
    repo    *MockUserRepository
}

func (suite *UserServiceTestSuite) SetupTest() {
    suite.repo = new(MockUserRepository)
    suite.service = NewUserService(suite.repo)
}

func (suite *UserServiceTestSuite) TestCreateUser_Success() {
    // Arrange
    user := &User{
        Username: "testuser",
        Email:    "test@example.com",
    }

    suite.repo.On("Create", user).Return(nil)

    // Act
    err := suite.service.CreateUser(user)

    // Assert
    assert.NoError(suite.T(), err)
    suite.repo.AssertExpectations(suite.T())
}

func TestUserServiceTestSuite(t *testing.T) {
    suite.Run(t, new(UserServiceTestSuite))
}
```

#### 2. 集成测试
```go
// user_integration_test.go
//go:build integration
// +build integration

package integration

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestUserHandler_CreateUser_Integration(t *testing.T) {
    // 设置测试数据库
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    // 创建测试服务器
    router := setupTestRouter(db)

    // 准备测试数据
    userReq := CreateUserRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }

    body, _ := json.Marshal(userReq)
    req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)

    var response CreateUserResponse
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.NotEmpty(t, response.Data.ID)
}
```

#### 3. Mock 测试
```go
// 使用 gomock 进行 mock
//go:generate mockgen -source=repository.go -destination=mocks/user_repository_mock.go

package service

import (
    "testing"
    "github.com/golang/mock/gomock"
    "github.com/stretchr/testify/assert"
)

func TestUserService_GetUser(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockUserRepository(ctrl)
    service := NewUserService(mockRepo)

    expectedUser := &User{
        ID:       1,
        Username: "testuser",
        Email:    "test@example.com",
    }

    mockRepo.EXPECT().
        GetByID(int64(1)).
        Return(expectedUser, nil).
        Times(1)

    user, err := service.GetUser(1)

    assert.NoError(t, err)
    assert.Equal(t, expectedUser, user)
}
```

### 前端测试

#### 1. 单元测试
```typescript
// userApi.test.ts
import { fetchUser, createUser } from '@/api/user';
import { api } from '@/utils/request';
import { User } from '@/types/user';

// Mock API 模块
jest.mock('@/utils/request');
const mockApi = api as jest.Mocked<typeof api>;

describe('User API', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('fetchUser', () => {
    it('should fetch user successfully', async () => {
      // Arrange
      const mockUser: User = {
        id: 1,
        username: 'testuser',
        email: 'test@example.com',
      };

      mockApi.get.mockResolvedValue({
        data: {
          success: true,
          data: mockUser,
        },
      });

      // Act
      const result = await fetchUser(1);

      // Assert
      expect(mockApi.get).toHaveBeenCalledWith('/users/1');
      expect(result).toEqual(mockUser);
    });

    it('should handle API error', async () => {
      // Arrange
      mockApi.get.mockRejectedValue(new Error('Network error'));

      // Act & Assert
      await expect(fetchUser(1)).rejects.toThrow('Network error');
    });
  });
});
```

#### 2. 组件测试
```typescript
// UserCard.test.tsx
import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { UserCard } from '@/components/UserCard';
import { User } from '@/types/user';

const mockUser: User = {
  id: 1,
  username: 'testuser',
  email: 'test@example.com',
};

describe('UserCard', () => {
  it('renders user information correctly', () => {
    render(<UserCard user={mockUser} />);

    expect(screen.getByText('testuser')).toBeInTheDocument();
    expect(screen.getByText('test@example.com')).toBeInTheDocument();
  });

  it('calls onEdit when edit button is clicked', () => {
    const mockOnEdit = jest.fn();
    render(<UserCard user={mockUser} onEdit={mockOnEdit} />);

    fireEvent.click(screen.getByText('编辑'));

    expect(mockOnEdit).toHaveBeenCalledWith(mockUser);
  });

  it('does not show edit button when onEdit is not provided', () => {
    render(<UserCard user={mockUser} />);

    expect(screen.queryByText('编辑')).not.toBeInTheDocument();
  });
});
```

#### 3. E2E 测试
```typescript
// e2e/user-auth.spec.ts
import { test, expect } from '@playwright/test';

test.describe('User Authentication', () => {
  test('should login successfully with valid credentials', async ({ page }) => {
    // 访问登录页面
    await page.goto('/login');

    // 填写登录表单
    await page.fill('[data-testid="username-input"]', 'testuser');
    await page.fill('[data-testid="password-input"]', 'password123');

    // 点击登录按钮
    await page.click('[data-testid="login-button"]');

    // 验证登录成功
    await expect(page).toHaveURL('/dashboard');
    await expect(page.locator('[data-testid="user-menu"]')).toContainText('testuser');
  });

  test('should show error message with invalid credentials', async ({ page }) => {
    await page.goto('/login');

    await page.fill('[data-testid="username-input"]', 'invaliduser');
    await page.fill('[data-testid="password-input"]', 'wrongpassword');
    await page.click('[data-testid="login-button"]');

    await expect(page.locator('[data-testid="error-message"]')).toBeVisible();
    await expect(page.locator('[data-testid="error-message"]')).toContainText('用户名或密码错误');
  });
});
```

---

## 🐛 调试技巧

### 后端调试

#### 1. 使用 Delve 调试器
```bash
# 安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 调试应用
dlv debug ./cmd/user-service

# 或调试测试
dlv test ./internal/service/

# Delve 命令
(dlv) break user_service.go:42  # 设置断点
(dlv) continue                  # 继续执行
(dlv) print user                # 打印变量
(dlv) locals                    # 显示局部变量
(dlv) stack                     # 显示调用栈
```

#### 2. VS Code 调试配置
```json
// .vscode/launch.json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch User Service",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/backend/cmd/user-service",
      "env": {
        "GIN_MODE": "debug"
      },
      "args": []
    },
    {
      "name": "Launch Tests",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${workspaceFolder}/backend",
      "env": {},
      "args": ["-test.run", "TestUserService"]
    }
  ]
}
```

#### 3. 日志调试
```go
// 使用结构化日志
import "github.com/sirupsen/logrus"

func (s *UserService) CreateUser(user *User) error {
    logger := logrus.WithFields(logrus.Fields{
        "user_id":   user.ID,
        "username":  user.Username,
        "email":     user.Email,
        "function":  "CreateUser",
    })

    logger.Info("开始创建用户")

    if err := s.validateUser(user); err != nil {
        logger.WithError(err).Error("用户验证失败")
        return err
    }

    if err := s.repo.Create(user); err != nil {
        logger.WithError(err).Error("数据库操作失败")
        return err
    }

    logger.Info("用户创建成功")
    return nil
}
```

### 前端调试

#### 1. VS Code 调试配置
```json
// .vscode/launch.json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Chrome",
      "type": "chrome",
      "request": "launch",
      "url": "http://localhost:5173",
      "webRoot": "${workspaceFolder}/frontend/src"
    }
  ]
}
```

#### 2. React DevTools
```bash
# 安装 React Developer Tools 浏览器扩展
# Chrome: https://chrome.google.com/webstore/detail/react-developer-tools/
# Firefox: https://addons.mozilla.org/en-US/firefox/addon/react-devtools/
```

#### 3. Redux DevTools (如使用 Redux)
```typescript
// store.ts
import { createStore } from 'redux';
import { rootReducer } from './reducers';

const store = createStore(
  rootReducer,
  window.__REDUX_DEVTOOLS_EXTENSION__ && window.__REDUX_DEVTOOLS_EXTENSION__()
);
```

---

## ⚡ 性能优化

### 后端优化

#### 1. 数据库优化
```go
// 使用索引
type User struct {
    ID       int64  `gorm:"primaryKey"`
    Username string `gorm:"uniqueIndex;size:50"`
    Email    string `gorm:"uniqueIndex;size:100"`
    CreatedAt time.Time `gorm:"index"`
}

// 预加载关联数据
func (r *UserRepository) GetWithOrders(id int64) (*User, error) {
    var user User
    err := r.db.
        Preload("Orders").
        Where("id = ?", id).
        First(&user).Error
    return &user, err
}

// 批量操作
func (r *UserRepository) CreateBatch(users []*User) error {
    return r.db.CreateInBatches(users, 100).Error
}
```

#### 2. 缓存策略
```go
// Redis 缓存
func (s *UserService) GetUser(id int64) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)

    // 先从缓存获取
    var user User
    if err := s.cache.Get(cacheKey, &user); err == nil {
        return &user, nil
    }

    // 缓存未命中，从数据库获取
    user, err := s.repo.GetByID(id)
    if err != nil {
        return nil, err
    }

    // 写入缓存，设置过期时间
    s.cache.Set(cacheKey, user, 10*time.Minute)

    return user, nil
}
```

#### 3. 并发处理
```go
// 使用 goroutine 池
type WorkerPool struct {
    workers    int
    jobQueue   chan Job
    workerPool chan chan Job
    quit       chan bool
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workers; i++ {
        worker := NewWorker(wp.workerPool)
        worker.Start()
    }

    go wp.dispatch()
}

func (wp *WorkerPool) dispatch() {
    for {
        select {
        case job := <-wp.jobQueue:
            go func() {
                workerChannel := <-wp.workerPool
                workerChannel <- job
            }()
        case <-wp.quit:
            return
        }
    }
}
```

### 前端优化

#### 1. 代码分割
```typescript
// 路由级别的代码分割
import { lazy, Suspense } from 'react';

const Dashboard = lazy(() => import('@/pages/Dashboard'));
const UserManagement = lazy(() => import('@/pages/UserManagement'));

function App() {
  return (
    <Router>
      <Suspense fallback={<div>Loading...</div>}>
        <Routes>
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/users" element={<UserManagement />} />
        </Routes>
      </Suspense>
    </Router>
  );
}
```

#### 2. 组件优化
```typescript
// 使用 React.memo 避免不必要的重渲染
const UserCard = React.memo<UserCardProps>(({ user, onEdit }) => {
  return (
    <div>
      <h3>{user.username}</h3>
      <button onClick={() => onEdit(user)}>编辑</button>
    </div>
  );
}, (prevProps, nextProps) => {
  // 自定义比较函数
  return prevProps.user.id === nextProps.user.id &&
         prevProps.user.username === nextProps.user.username;
});

// 使用 useMemo 缓存计算结果
const ExpensiveComponent: React.FC<{ items: Item[] }> = ({ items }) => {
  const expensiveValue = useMemo(() => {
    return items.reduce((sum, item) => sum + item.value, 0);
  }, [items]);

  return <div>Total: {expensiveValue}</div>;
};

// 使用 useCallback 缓存函数
const ParentComponent: React.FC = () => {
  const [count, setCount] = useState(0);

  const handleClick = useCallback(() => {
    setCount(prev => prev + 1);
  }, []);

  return <ChildComponent onClick={handleClick} />;
};
```

#### 3. 请求优化
```typescript
// 请求防抖
import { debounce } from 'lodash';

const SearchInput: React.FC = () => {
  const [query, setQuery] = useState('');

  const debouncedSearch = useMemo(
    () => debounce(async (searchQuery: string) => {
      const results = await searchUsers(searchQuery);
      // 处理搜索结果
    }, 300),
    []
  );

  useEffect(() => {
    debouncedSearch(query);

    return () => {
      debouncedSearch.cancel();
    };
  }, [query, debouncedSearch]);

  return (
    <input
      value={query}
      onChange={(e) => setQuery(e.target.value)}
      placeholder="搜索用户..."
    />
  );
};

// 请求缓存和重试
import { useQuery } from '@tanstack/react-query';

const useUsers = () => {
  return useQuery({
    queryKey: ['users'],
    queryFn: fetchUsers,
    staleTime: 5 * 60 * 1000, // 5分钟
    retry: 3,
    retryDelay: attemptIndex => Math.min(1000 * 2 ** attemptIndex, 30000),
  });
};
```

---

## ❓ 常见问题

### 后端常见问题

#### 1. 数据库连接问题
```bash
# 检查 MySQL 服务状态
brew services list | grep mysql  # macOS
systemctl status mysql           # Linux

# 检查连接配置
mysql -u gamelink -p -h localhost gamelink

# 常见错误解决
# Error 1045: Access denied - 检查用户名密码
# Error 2003: Can't connect - 检查服务是否启动
```

#### 2. Go modules 问题
```bash
# 清理模块缓存
go clean -modcache

# 重新下载依赖
go mod download
go mod tidy

# 代理设置 (如果需要)
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=sum.golang.google.cn
```

#### 3. 端口占用问题
```bash
# 查看端口占用
lsof -i :8080  # macOS/Linux
netstat -ano | findstr :8080  # Windows

# 杀死进程
kill -9 <PID>  # macOS/Linux
taskkill /PID <PID> /F  # Windows
```

### 前端常见问题

#### 1. Node.js 版本问题
```bash
# 使用 nvm 管理 Node.js 版本
nvm list
nvm use 18
nvm install 18
```

#### 2. 依赖安装问题
```bash
# 清理缓存
npm cache clean --force

# 删除 node_modules 重新安装
rm -rf node_modules package-lock.json
npm install
```

#### 3. 端口冲突
```bash
# 杀死占用端口的进程
npm run dev -- --port 3001  # 指定其他端口

# 或修改 vite.config.ts
export default defineConfig({
  server: {
    port: 3001,
  },
});
```

### 开发工具问题

#### 1. VS Code 扩展问题
```bash
# 重新加载窗口
Ctrl+Shift+P -> "Developer: Reload Window"

# 禁用冲突扩展
Ctrl+Shift+X -> 搜索扩展 -> 禁用
```

#### 2. Git 问题
```bash
# 查看配置
git config --list

# 设置用户信息
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"

# 解决合并冲突
git mergetool
```

---

## 📞 获取帮助

### 技术支持
- 📋 **Issues**: [GitHub Issues](https://github.com/your-org/GameLink/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/your-org/GameLink/discussions)
- 📧 **邮件**: dev-team@gamelink.com

### 学习资源
- [Go 官方文档](https://golang.org/doc/)
- [React 官方文档](https://react.dev/)
- [TypeScript 手册](https://www.typescriptlang.org/docs/)
- [项目内部文档](../docs/)

---

*最后更新: 2025-11-13*