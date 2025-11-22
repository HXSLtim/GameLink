# 环境变量配置指南

## 📋 概述

本项目使用 Vite 的环境变量系统，支持按环境加载不同的配置。配置文件遵循以下规则：

- 开发环境使用 `.env.development`
- 生产环境使用 `.env.production`
- 所有变量必须以 `VITE_` 开头才能在客户端代码中访问
- 本地覆盖配置可以使用 `.env.local`（不会被 Git 追踪）

## 🗂️ 配置文件说明

### `.env.development` - 开发环境配置

在 `npm run dev` 时自动加载，包含开发环境专用的配置：
- 本地 API 地址 (http://localhost:8080)
- Mock 登录凭据
- 演示功能开关（默认开启）

### `.env.production` - 生产环境配置

在 `npm run build` 时自动加载，包含生产环境专用的配置：
- 正式 API 地址
- 生产环境加密密钥（占位符，需修改）
- 演示功能开关（默认关闭）

### `.env.example` - 配置模板

包含所有可用的环境变量说明，可复制创建新的环境配置文件。

### `.env.local` - 本地覆盖配置（可选）

用于覆盖特定环境的配置，适合存储个人开发环境的敏感信息：
- 不会被 Git 追踪（在 .gitignore 中）
- 优先级高于 `.env.development` 和 `.env.production`
- 适合存储个人测试账号、本地特殊配置等

## ⚙️ 环境变量列表

### API 配置（必需）

| 变量名 | 说明 | 开发环境 | 生产环境 |
|--------|------|----------|----------|
| `VITE_API_BASE_URL` | 后端 API 基础地址 | `http://localhost:8080` | `https://api.gamelink.com` |

### Mock 配置（仅开发环境）

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `VITE_DEV_MOCK_USERNAME` | Mock 登录用户名 | `admin` |
| `VITE_DEV_MOCK_PASSWORD` | Mock 登录密码 | `admin123` |
| `VITE_DEV_MOCK_TOKEN` | Mock Token | `dev-token` |

### 加密配置（安全相关）

| 变量名 | 说明 | 开发环境 | 生产环境 |
|--------|------|----------|----------|
| `VITE_CRYPTO_ENABLED` | 是否启用数据加密 | `false` | `true` |
| `VITE_CRYPTO_SECRET_KEY` | AES 密钥（32字符） | `GameLink2025SecretKey!@#123456` | **必须修改！** |
| `VITE_CRYPTO_IV` | AES 初始化向量（16字符） | `GameLink2025IV!!!` | **必须修改！** |

**⚠️ 安全警告：**
- 生产环境的密钥必须通过以下命令生成强随机值：
  ```bash
  # 生成 32 字节密钥
  node -e "console.log(require('crypto').randomBytes(32).toString('base64'))"

  # 生成 16 字节 IV
  node -e "console.log(require('crypto').randomBytes(16).toString('base64'))"
  ```

### 功能开关

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `VITE_ENABLE_MOCK` | 是否启用 Mock 数据 | `false` |
| `VITE_ENABLE_ANALYZER` | 是否启用 Bundle 分析器 | `false` |

### 演示功能配置

| 变量名 | 说明 | 开发环境 | 生产环境 |
|--------|------|----------|----------|
| `VITE_SHOWCASE_ENABLE_ROUTE` | 是否启用组件演示路由 | `true` | `false` |
| `VITE_SHOWCASE_ENABLE_COMPONENTS` | 是否展示组件示例 | `true` | `false` |
| `VITE_SHOWCASE_EXPAND_EXAMPLES` | 是否展开示例代码 | `true` | `false` |
| `VITE_CACHE_DEMO_ENABLE_ROUTE` | 是否启用缓存演示路由 | `true` | `false` |

### 第三方服务配置（可选）

| 变量名 | 说明 | 示例 |
|--------|------|------|
| `VITE_SENTRY_DSN` | Sentry DSN（错误监控） | `https://xxx@sentry.io/xxx` |
| `VITE_GA_ID` | Google Analytics ID | `G-XXXXXXXXXX` |

### 性能监控（可选）

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `VITE_ENABLE_PERFORMANCE_MONITOR` | 是否启用性能监控 | `false` |

### 日志配置（可选）

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `VITE_LOG_LEVEL` | 日志级别 | `info` |
| `VITE_ENABLE_LOG_UPLOAD` | 是否上传日志到服务器 | `false` |

### 部署配置（可选）

| 变量名 | 说明 | 示例 |
|--------|------|------|
| `VITE_DEPLOY_ENV` | 部署环境 | `development` / `staging` / `production` |
| `VITE_APP_VERSION` | 应用版本号 | `0.1.0` |
| `VITE_BUILD_TIME` | 构建时间（自动注入） | `2025-01-01T00:00:00Z` |

## 💻 在代码中使用环境变量

### TypeScript/React 代码

```typescript
// 直接通过 import.meta.env 访问
const apiUrl = import.meta.env.VITE_API_BASE_URL;
const isCryptoEnabled = import.meta.env.VITE_CRYPTO_ENABLED === 'true';

// 在组件中使用
function LoginComponent() {
  const mockUsername = import.meta.env.VITE_DEV_MOCK_USERNAME;
  const mockPassword = import.meta.env.VITE_DEV_MOCK_PASSWORD;

  return (
    <div>
      <p>API地址: {apiUrl}</p>
      <p>加密: {isCryptoEnabled ? '开启' : '关闭'}</p>
    </div>
  );
}
```

### Vite 配置中使用

```typescript
// vite.config.ts
export default defineConfig(({ mode }) => {
  // 加载环境变量
  const env = loadEnv(mode, process.cwd());

  return {
    // 使用环境变量
    base: env.VITE_APP_BASE_URL || '/',
    plugins: [
      // ...
    ],
  };
});
```

## 🚀 使用场景示例

### 场景 1：切换后端 API 地址

在开发过程中需要连接不同的后端环境：

```bash
# .env.development
VITE_API_BASE_URL=http://localhost:8080          # 本地开发环境
# VITE_API_BASE_URL=http://192.168.1.100:8080    # 局域网测试环境
# VITE_API_BASE_URL=https://dev-api.gamelink.com # 开发服务器
```

### 场景 2：本地覆盖敏感配置

创建 `.env.local` 文件，存储个人测试账号：

```bash
# .env.local（不会被 Git 追踪）
VITE_DEV_MOCK_USERNAME=mytestaccount
VITE_DEV_MOCK_PASSWORD=mypassword123
VITE_DEV_MOCK_TOKEN=my-custom-token
```

### 场景 3：启用 Bundle 分析

分析生产构建的包大小：

```bash
# 方法 1：修改 .env.production
VITE_ENABLE_ANALYZER=true

# 方法 2：使用命令行
VITE_ENABLE_ANALYZER=true npm run build

# 方法 3：使用 npm script
npm run build:analyze
```

### 场景 4：集成第三方服务

在生产环境启用错误监控：

```bash
# .env.production
VITE_SENTRY_DSN=https://xxx@sentry.io/xxx
VITE_GA_ID=G-XXXXXXXXXX
VITE_ENABLE_PERFORMANCE_MONITOR=true
```

## 🔧 管理命令

### 开发环境

```bash
# 启动开发服务器（自动加载 .env.development）
npm run dev
```

### 生产构建

```bash
# 生产构建（自动加载 .env.production）
npm run build

# 带 Bundle 分析的生产构建
npm run build:analyze
```

### 验证配置

```bash
# 检查环境变量是否加载
npm run dev

# 在浏览器控制台输入
console.log(import.meta.env)
```

## 📝 最佳实践

### ✅ 推荐做法

1. **不要提交敏感信息**
   - 将 `.env.local` 添加到 `.gitignore`
   - 使用占位符值提交 `.env.development` 和 `.env.production`

2. **使用描述性的变量名**
   ```bash
   # ✅ 推荐
   VITE_API_BASE_URL=http://localhost:8080

   # ❌ 避免
   VITE_URL=http://localhost:8080
   ```

3. **为布尔值使用字符串**
   ```bash
   # ✅ 推荐
   VITE_FEATURE_ENABLED=true

   # ❌ 避免
   VITE_FEATURE_ENABLED=1
   ```

4. **添加详细注释**
   ```bash
   # API 基础地址
   # 开发环境: http://localhost:8080
   # 生产环境: https://api.gamelink.com
   VITE_API_BASE_URL=http://localhost:8080
   ```

5. **分组管理**
   按功能分组，使用注释分隔：
   ```bash
   # ========================================
   # API 配置
   # ========================================
   VITE_API_BASE_URL=...

   # ========================================
   # 功能开关
   # ========================================
   VITE_ENABLE_FEATURE=...
   ```

### ⚠️ 安全提醒

1. **生产环境必须修改加密密钥**
   - 使用强随机生成的密钥
   - 不要重复使用开发环境的密钥

2. **不要在客户端暴露敏感信息**
   - 所有 `VITE_` 前缀的变量都会在构建时注入到客户端代码
   - 不要在环境变量中存储：API Secret、数据库密码、私钥等

3. **使用环境变量管理敏感配置**
   - 在 CI/CD 流程中注入敏感配置
   - 不要硬编码到代码中

4. **定期轮换密钥**
   - 特别是生产环境的加密密钥
   - 及时撤销泄露的密钥

## 📚 相关文档

- [Vite 环境变量文档](https://vite.dev/guide/env-and-mode.html)
- [npm scripts 配置](../package.json)
- [项目部署指南](./DEPLOYMENT.md)

## ❓ 常见问题

### Q: 修改环境变量后需要重启开发服务器吗？

**A:** 需要。环境变量在启动时加载，修改后必须重启才能生效。

### Q: 为什么我的环境变量不生效？

**A:** 检查以下几点：
1. 变量名是否以 `VITE_` 开头
2. 文件命名是否正确（.env.development / .env.production）
3. 是否重启了开发服务器
4. 检查控制台是否有错误信息

### Q: 如何在运行时切换环境？

**A:** Vite 在构建时确定环境，运行时无法切换。如需切换环境，需重新构建应用。

### Q: 如何添加新的环境变量？

**A:** 按以下步骤：
1. 更新 `.env.example`（添加说明）
2. 更新 `.env.development`（开发环境值）
3. 更新 `.env.production`（生产环境占位符）
4. 在代码中使用 `import.meta.env.VITE_YOUR_VAR`

## 🔍 调试技巧

### 查看所有环境变量

```typescript
// 在浏览器控制台
console.log('所有环境变量:', import.meta.env);

// 检查特定变量
console.log('API地址:', import.meta.env.VITE_API_BASE_URL);
console.log('加密开关:', import.meta.env.VITE_CRYPTO_ENABLED);
```

### 环境变量加载顺序

Vite 按以下顺序加载环境变量（后面的优先级更高）：

1. `.env` - 默认环境
2. `.env.local` - 本地覆盖（不被 Git 追踪）
3. `.env.[mode]` - 特定模式（development/production）
4. `.env.[mode].local` - 特定模式的本地覆盖

这意味着：
- `.env.development.local` 会覆盖 `.env.development` 中的变量
- 可以在 `.env.local` 中临时修改配置，而不修改提交到 Git 的文件

---

**最后更新：** 2025-11-22
**维护者：** GameLink 开发团队
