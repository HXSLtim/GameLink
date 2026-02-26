# GameLink 开发者入职指南

> **版本**: v2.0
> **最后更新**: 2026-02-11
> **目标**: 帮助新开发者快速融入项目

---

## 目录

1. [第一天：环境搭建](#第一天环境搭建)
2. [第一周：代码入门](#第一周代码入门)
3. [第一个月：独立开发](#第一个月独立开发)
4. [开发工具配置](#开发工具配置)
5. [代码规范](#代码规范)
6. [常用命令](#常用命令)
7. [团队协作](#团队协作)
8. [常见问题](#常见问题)

---

## 第一天：环境搭建

### 上午（2-3小时）

#### 1.1 账号与权限配置

- [ ] 申请公司邮箱账号
- [ ] 申请 GitHub/GitLab 访问权限
- [ ] 申请开发服务器 SSH 密钥
- [ ] 申请数据库只读权限（仅查看）
- [ ] 加入团队沟通群（钉钉/飞书/Slack）

#### 1.2 必需软件安装

**必装软件**:
```bash
# 1. Git
git --version  # >= 2.30

# 2. Go
go version  # >= 1.24

# 3. Node.js
node --version  # >= 20
npm --version

# 4. Docker
docker --version
docker-compose --version

# 5. IDE（二选一）
# Visual Studio Code
code --version

# GoLand
# 下载并安装
```

**推荐工具**:
- Postman（API 测试）
- DBeaver（数据库客户端）
-浏览器开发者工具（Web 调试）

### 下午（3-4小时）

#### 1.3 克隆项目并启动

```bash
# 1. 克隆项目
git clone https://github.com/your-org/gamelink.git
cd gamelink

# 2. 查看项目结构
tree -L 2 -I 'node_modules|dist'

# 3. 配置环境变量
cp .env.example .env
# 编辑 .env 文件（默认配置即可启动）

# 4. 启动基础服务
docker compose up -d postgres redis

# 5. 启动后端
cd api
go mod download
go run cmd/main.go
# 验证: curl http://localhost:8080/health

# 6. 启动管理后台
cd ../admin
npm install
npm run dev
# 访问: http://localhost:5173
# 登录: admin@gamelink.com / Admin123456

# 7. 启动用户端 Web
cd ../app
npm install
npm run dev
# 访问: http://localhost:5175
```

#### 1.4 第一天检查清单

- [ ] 后端 API 可以访问
- [ ] 管理后台可以登录
- [ ] 用户端 Web 可以打开
- [ ] 可以查看 Swagger 文档
- [ ] 可以访问数据库

---

## 第一周：代码入门

### Day 1-2: 理解项目结构

**学习目标**: 熟悉项目目录和各模块职责

```
GameLink/
├── api/                # Go 后端服务
│   ├── cmd/            # 应用入口
│   ├── internal/       # 内部代码
│   │   ├── handler/   # HTTP 处理器（控制器）
│   │   ├── service/   # 业务逻辑层
│   │   ├── repository/# 数据访问层
│   │   ├── model/     # 数据模型
│   │   └── router/    # 路由注册
│   └── pkg/           # 公共包
├── admin/             # React 管理后台
│   └── src/
│       ├── pages/     # 页面组件
│       ├── components/# 通用组件
│       ├── api/       # API 封装
│       └── store/     # 状态管理
└── app/              # 用户端 Web
    └── src/
        ├── features/   # 页面功能
        ├── components/ # 组件
        └── services/   # API 调用
```

**阅读顺序**:
1. `README.md` - 项目概述
2. `PROJECT_OVERVIEW.md` - 详细介绍
3. `docs/ARCHITECTURE_DIAGRAMS.md` - 架构图
4. `api/internal/handler/` - 选择一个模块阅读

### Day 3-4: 运行测试和调试

```bash
# 后端测试
cd api
go test ./... -v
go test ./internal/service/order -v -cover

# 管理后台测试
cd admin
npm run test

# 调试配置
# VS Code: .vscode/launch.json
# GoLand: Run -> Edit Configurations
```

**练习任务**:
1. 修改一个接口的返回值
2. 添加一个测试用例
3. 使用断点调试代码

### Day 5: 第一次代码提交

```bash
# 1. 创建功能分支
git checkout -b feature/hello-world

# 2. 进行简单修改
# 例如：修改首页欢迎文字

# 3. 提交代码
git add .
git commit -m "feat: update welcome message"

# 4. 推送到远程
git push origin feature/hello-world

# 5. 创建 Pull Request
# 在 GitHub 上操作
```

**第一周检查清单**:
- [ ] 理解项目目录结构
- [ ] 能够运行所有测试
- [ ] 完成第一次代码提交
- [ ] 了解 Git 工作流程

---

## 第一个月：独立开发

### Week 1: 熟悉核心业务

**学习内容**:
1. 数据库表关系（`docs/ARCHITECTURE_DIAGRAMS.md`）
2. 核心业务流程（下单、支付、完成）
3. API 接口规范（`docs/API_ALIGNMENT.md`）

**练习任务**:
1. 实现一个简单的查询接口
2. 添加一个新的管理后台页面

### Week 2-3: 完成一个小功能

**推荐任务**（从简单到复杂）:
1. 添加一个筛选条件
2. 实现批量操作功能
3. 添加一个新的数据导出功能

**开发流程**:
```
1. 需求确认（与产品经理沟通）
2. 技术设计（简单方案即可）
3. 开发实现
4. 自测
5. 提交代码审查
6. 修改反馈意见
7. 合并到主分支
```

### Week 4: 参与代码审查

**审查要点**:
1. 代码是否遵循规范
2. 是否有安全问题
3. 是否有性能问题
4. 测试是否充分
5. 文档是否完整

**第一个月检查清单**:
- [ ] 独立完成至少一个功能
- [ ] 参与 5+ 次代码审查
- [ ] 修复至少一个 Bug
- [ ] 理解整个业务流程

---

## 开发工具配置

### VS Code

**推荐扩展**:
```json
{
  "recommendations": [
    "golang.go",
    "bradlc.vscode-tailwindcss",
    "dbaeumer.vscode-eslint",
    "esbenp.prettier-vscode",
    "ms-vscode.live-server",
    "humao.rest-client"
  ]
}
```

**工作区设置** (`.vscode/settings.json`):
```json
{
  "go.formatTool": "goimports",
  "go.lintTool": "golangci-lint",
  "go.testFlags": ["-v"],
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.fixAll.eslint": true
  }
}
```

### Git 配置

```bash
# 全局配置
git config --global user.name "Your Name"
git config --global user.email "your.email@company.com"

# Git 别名
git config --global alias.st status
git config --global alias.co checkout
git config --global alias.br branch
git config --global alias.ci commit

# GitHub 认证（推荐使用 SSH）
ssh-keygen -t ed25519 -C "your.email@company.com"
# 将公钥添加到 GitHub
```

### 数据库客户端配置

**DBeaver**（推荐）:
```
连接名称: GameLink Dev
驱动: PostgreSQL
主机: localhost
端口: 5432
数据库: gamelink_dev
用户名: gamelink
密码: [见 .env 文件]
```

---

## 代码规范

### Go 代码规范

**命名**:
```go
// ✅ 好的命名
type UserService struct {}
func (s *UserService) GetUserByID(ctx context.Context, id int) (*User, error) {}
const MAX_RETRY_COUNT = 3
var userCache = make(map[int]*User)

// ❌ 不好的命名
type US struct {}
func (s *US) Get(ctx context.Context, i int) (*User, error) {}
const max = 3
var cache = make(map[int]*User)
```

**错误处理**:
```go
// ✅ 好的错误处理
user, err := s.repo.Get(ctx, userID)
if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, ErrUserNotFound
    }
    return nil, fmt.Errorf("failed to get user: %w", err)
}

// ❌ 不好的错误处理
user, err := s.repo.Get(ctx, userID)
if err != nil {
    return nil, err
}
```

**注释规范**:
```go
// UserService 处理用户相关业务逻辑
type UserService struct {
    repo   UserRepository
    logger *zap.Logger
}

// GetUserByID 根据用户ID获取用户信息
// 如果用户不存在，返回 ErrUserNotFound
func (s *UserService) GetUserByID(ctx context.Context, id int) (*User, error) {
    // 实现细节...
}
```

### React 代码规范

**组件定义**:
```tsx
// ✅ 函数组件 + Hooks
interface UserListProps {
  users: User[];
  onEdit: (user: User) => void;
}

export const UserList: React.FC<UserListProps> = ({ users, onEdit }) => {
  return (
    <div>
      {users.map(user => (
        <UserCard key={user.id} user={user} onEdit={onEdit} />
      ))}
    </div>
  );
};

// ❌ 类组件（不推荐用于新代码）
class UserList extends React.Component {
  // ...
}
```

**状态管理**:
```tsx
// ✅ 使用 Zustand
interface UserStore {
  users: User[];
  fetchUsers: () => Promise<void>;
}

export const useUserStore = create<UserStore>((set) => ({
  users: [],
  fetchUsers: async () => {
    const users = await api.getUsers();
    set({ users });
  },
}));
```

### Vue 代码规范

**组件定义**:
```vue
<script setup lang="ts">
// ✅ Composition API
import { ref, computed } from 'vue';

interface Props {
  userId: number;
}
const props = defineProps<Props>();

const user = ref<User | null>(null);
const isLoading = ref(false);

const displayName = computed(() => user.value?.name ?? 'Unknown');
</script>
```

---

## 常用命令

### 后端开发

```bash
# 运行
go run cmd/main.go

# 构建
go build -o gamelink-api cmd/main.go

# 测试
go test ./... -v
go test ./internal/service/order -v -cover

# 代码检查
go vet ./...
golangci-lint run

# 格式化
go fmt ./...
goimports -w .

# 查看 Go 版本
go version
go env
```

### 前端开发（管理后台）

```bash
# 安装依赖
npm install

# 开发
npm run dev

# 构建
npm run build

# 预览构建结果
npm run preview

# 测试
npm run test

# 代码检查
npm run lint

# 格式化
npm run format
```

### 前端开发（用户端 Web）

```bash
# 启动开发
npm run dev

# 构建
npm run build

# 预览
npm run preview
```

### Docker 命令

```bash
# 启动服务
docker compose up -d

# 查看日志
docker compose logs -f api

# 停止服务
docker compose down

# 重启服务
docker compose restart api

# 进入容器
docker compose exec api bash
docker compose exec postgres psql -U gamelink -d gamelink_dev
```

### Git 命令

```bash
# 查看状态
git status

# 查看分支
git branch -a

# 切换分支
git checkout -b feature/new-feature

# 提交
git add .
git commit -m "feat: add new feature"

# 推送
git push origin feature/new-feature

# 拉取最新
git pull origin main

# 合并分支
git merge feature/new-feature
```

---

## 团队协作

### 分支策略

```
main      # 生产环境分支
  ↑
dev       # 开发环境分支
  ↑
feature/* # 功能分支
bugfix/*  # 修复分支
```

**工作流程**:
```
1. 从 dev 创建功能分支
   git checkout dev
   git pull origin dev
   git checkout -b feature/your-feature

2. 开发并提交
   git add .
   git commit -m "feat: your feature description"

3. 推送到远程
   git push origin feature/your-feature

4. 创建 Pull Request
   在 GitHub 上创建 PR 到 dev

5. 代码审查通过后合并
   Squash and Merge
```

### 代码审查清单

**功能审查**:
- [ ] 功能是否符合需求
- [ ] 边界情况是否处理
- [ ] 错误处理是否完善

**代码质量**:
- [ ] 代码是否易读
- [ ] 是否有重复代码
- [ ] 命名是否清晰

**安全审查**:
- [ ] SQL 注入风险
- [ ] XSS 风险
- [ ] 敏感信息泄露

**测试审查**:
- [ ] 是否有单元测试
- [ ] 测试覆盖率是否足够
- [ ] 是否有集成测试

### 沟通渠道

| 沟通类型 | 工具 | 说明 |
|---------|------|------|
| 日常沟通 | 钉钉/飞书 | 快速问题讨论 |
| 代码审查 | GitHub PR | 正式代码审查 |
| 文档协作 | Wiki/Notion | 文档编写 |
| 会议 | 腾讯会议/Zoom | 技术讨论 |

---

## 常见问题

### Q1: 如何查看数据库表结构？

```bash
# 方法1: 使用 DBeaver 图形界面
# 方法2: 使用 psql 命令
\c gamelink_dev
\d users  # 查看表结构
\d+ users # 查看详细信息

# 方法3: 查看 model 文件
cat api/internal/model/user.go
```

### Q2: 如何调试 API？

```bash
# 方法1: 使用 curl
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer your_token"

# 方法2: 使用 Postman
# 导入 API_COLLECTION.json

# 方法3: 使用 Swagger
# 访问 http://localhost:8080/swagger/index.html
```

### Q3: 代码提交被拒绝了怎么办？

```bash
# 1. 查看审查意见
# 在 GitHub PR 页面查看评论

# 2. 修改代码
# 按照意见修改

# 3. 更新分支
git add .
git commit -m "fix: address review comments"
git push origin feature/your-feature

# 4. PR 会自动更新
```

### Q4: 遇到无法解决的问题？

```
1. 查看项目文档
   - README.md
   - ARCHITECTURE_DIAGRAMS.md
   - TROUBLESHOOTING.md

2. 搜索问题
   - GitHub Issues
   - 项目 Wiki

3. 寻求帮助
   - 在团队群里提问
   - @技术负责人

4. 记录问题
   - 创建新的 GitHub Issue
   - 或更新 TROUBLESHOOTING.md
```

---

## 检查清单

### 第一周结束
- [ ] 开发环境正常运行
- [ ] 理解项目结构
- [ ] 完成第一次代码提交
- [ ] 知道如何寻求帮助

### 第一个月结束
- [ ] 独立完成一个功能
- [ ] 参与 5+ 次代码审查
- [ ] 修复至少一个 Bug
- [ ] 理解核心业务流程

### 第三个月结束
- [ ] 能独立设计功能方案
- [ ] 能指导新人
- [ ] 能进行代码审查
- [ ] 能优化代码性能

---

**祝你在 GameLink 工作愉快！** 🚀

如有任何问题，随时联系团队。

**文档维护**: 产品经理 + 技术负责人
**更新频率**: 每月更新
