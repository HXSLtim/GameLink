# 环境变量配置指南

## 📋 环境变量文件说明

| 文件                     | 用途                 | Git 追踪  | 优先级    |
| ------------------------ | -------------------- | --------- | --------- |
| `.env`                   | 所有环境的默认值     | ✅ 提交   | 1（最低） |
| `.env.local`             | 本地覆盖（所有环境） | ❌ 不提交 | 4         |
| `.env.development`       | 开发环境             | ✅ 提交   | 2         |
| `.env.production`        | 生产环境             | ✅ 提交   | 2         |
| `.env.development.local` | 本地开发覆盖         | ❌ 不提交 | 3         |
| `.env.production.local`  | 本地生产覆盖         | ❌ 不提交 | 3         |

## ✅ 已创建的文件

### 1. `.env.example` - 配置模板

```bash
# 用途：提供配置示例，团队成员参考
# 状态：✅ 已创建
# 提交：是
```

### 2. `.env.development` - 开发环境

```bash
# 用途：npm run dev 时自动加载
# 状态：✅ 已创建
# 配置：
VITE_API_BASE_URL=http://localhost:8080
VITE_CRYPTO_ENABLED=false  # 开发环境禁用加密
VITE_CRYPTO_SECRET_KEY=GameLink2025SecretKey!@#123456
VITE_CRYPTO_IV=GameLink2025IV!!!
```

### 3. `.env.production` - 生产环境

```bash
# 用途：npm run build 时自动加载
# 状态：✅ 已创建
# 配置：
VITE_API_BASE_URL=https://api.gamelink.com
VITE_CRYPTO_ENABLED=true  # 生产环境启用加密
VITE_CRYPTO_SECRET_KEY=GameLink2025SecretKey!@#123456  # ⚠️ 需修改
VITE_CRYPTO_IV=GameLink2025IV!!!  # ⚠️ 需修改
```

### 4. `.env.local.example` - 本地配置示例

```bash
# 用途：本地测试时参考
# 状态：✅ 已创建
# 使用：复制为 .env.local
```

## 🔧 使用方法

### 开发环境（默认）

```bash
# 直接启动，自动加载 .env.development
npm run dev

# 访问环境变量
console.log(import.meta.env.VITE_API_BASE_URL)
// 输出: http://localhost:8080

console.log(import.meta.env.VITE_CRYPTO_ENABLED)
// 输出: "false"（开发环境禁用加密）
```

### 生产构建

```bash
# 构建时自动加载 .env.production
npm run build

# 预览生产构建
npm run preview
```

### 本地覆盖配置

如果你需要本地特殊配置（不影响其他开发者）：

```bash
# 1. 复制示例文件
cp .env.local.example .env.local

# 2. 修改 .env.local（此文件不会提交到 git）
VITE_API_BASE_URL=http://192.168.1.100:8080
VITE_CRYPTO_ENABLED=true

# 3. 启动开发服务器，会优先使用 .env.local 的配置
npm run dev
```

## 🔐 加密配置说明

### 开发环境（默认禁用）

```bash
# .env.development
VITE_CRYPTO_ENABLED=false  # 禁用加密，方便调试
```

**为什么禁用？**

- 方便查看网络请求
- 容易定位问题
- 后端可能还未实现解密中间件

### 生产环境（默认启用）

```bash
# .env.production
VITE_CRYPTO_ENABLED=true  # 启用加密，保护数据
```

### 临时启用/禁用

```bash
# 开发时临时启用加密测试
VITE_CRYPTO_ENABLED=true npm run dev

# 生产构建时临时禁用加密
VITE_CRYPTO_ENABLED=false npm run build
```

## 🔑 密钥配置

### 当前密钥（默认/测试）

```bash
VITE_CRYPTO_SECRET_KEY=GameLink2025SecretKey!@#123456  # 32字节
VITE_CRYPTO_IV=GameLink2025IV!!!  # 16字节
```

⚠️ **这些是测试密钥，生产环境必须修改！**

### 生成强密钥

#### 方法 1：使用 Node.js

```bash
# 生成 32 字节密钥（Base64）
node -e "console.log(require('crypto').randomBytes(32).toString('base64'))"
# 输出示例：J5kY8vZ2R6nP9xQ3mL7wC1bH4tS5dE8fG...

# 生成 16 字节向量（Base64）
node -e "console.log(require('crypto').randomBytes(16).toString('base64'))"
# 输出示例：aBcD123eFgH456iJ==
```

#### 方法 2：在线生成

访问：https://www.random.org/passwords/

- 密钥：长度 32，包含字母+数字+符号
- 向量：长度 16，包含字母+数字+符号

### 更新生产密钥

```bash
# 1. 编辑 .env.production
vim .env.production

# 2. 修改密钥
VITE_CRYPTO_SECRET_KEY=你生成的32字节密钥
VITE_CRYPTO_IV=你生成的16字节向量

# 3. 通知后端同步密钥
# ⚠️ 前后端密钥必须完全一致！
```

## 🌍 多环境配置

### 测试环境

```bash
# 创建 .env.test
cat > .env.test << 'EOF'
VITE_API_BASE_URL=https://test-api.gamelink.com
VITE_CRYPTO_ENABLED=true
VITE_CRYPTO_SECRET_KEY=test-environment-key-32-bytes
VITE_CRYPTO_IV=test-env-iv-16b
EOF

# 使用测试环境构建
npm run build -- --mode test
```

### 预发布环境

```bash
# 创建 .env.staging
cat > .env.staging << 'EOF'
VITE_API_BASE_URL=https://staging-api.gamelink.com
VITE_CRYPTO_ENABLED=true
VITE_CRYPTO_SECRET_KEY=staging-environment-key-32bytes
VITE_CRYPTO_IV=staging-env-iv16
EOF

# 使用预发布环境构建
npm run build -- --mode staging
```

## 📊 环境变量优先级

从高到低：

1. `.env.*.local` （最高优先级，不提交）
2. `.env.local` （所有环境的本地覆盖，不提交）
3. `.env.production` / `.env.development` （环境特定）
4. `.env` （默认值，最低优先级）

**示例**：

```bash
# .env
VITE_API_BASE_URL=http://default.com

# .env.development
VITE_API_BASE_URL=http://localhost:8080

# .env.local
VITE_API_BASE_URL=http://192.168.1.100:8080

# npm run dev 时的实际值：
# http://192.168.1.100:8080  （.env.local 优先级最高）
```

## 🔍 调试环境变量

### 查看所有环境变量

```typescript
// src/main.tsx 或任意组件中
console.log('环境变量:', {
  mode: import.meta.env.MODE,
  dev: import.meta.env.DEV,
  prod: import.meta.env.PROD,
  apiUrl: import.meta.env.VITE_API_BASE_URL,
  cryptoEnabled: import.meta.env.VITE_CRYPTO_ENABLED,
});
```

### 检查加密状态

```typescript
import { cryptoMiddleware } from 'middleware/crypto';

console.log('加密中间件状态:', {
  enabled: cryptoMiddleware.isEnabled(),
  config: cryptoMiddleware.getConfig(),
});
```

## ⚠️ 安全注意事项

### ✅ 应该做的

- ✅ 使用 `VITE_` 前缀暴露给客户端
- ✅ 敏感配置文件加入 `.gitignore`
- ✅ 生产环境使用强密钥
- ✅ 定期轮换密钥
- ✅ 前后端密钥保持一致

### ❌ 不应该做的

- ❌ 将 `.env.local` 提交到 git
- ❌ 在代码中硬编码密钥
- ❌ 在生产环境使用默认密钥
- ❌ 将密钥直接写在前端代码中
- ❌ 在公共仓库暴露真实密钥

## 🚀 部署清单

### 部署前检查

- [ ] 修改了生产环境密钥
- [ ] 前后端密钥已同步
- [ ] `.env.production` 配置正确
- [ ] API 地址指向生产服务器
- [ ] 启用了加密（`VITE_CRYPTO_ENABLED=true`）
- [ ] 启用了 HTTPS

### 部署命令

```bash
# 1. 验证环境变量
cat .env.production

# 2. 构建生产版本
npm run build

# 3. 预览构建结果
npm run preview

# 4. 部署
# 将 dist 目录部署到服务器
```

## 📚 参考资料

- [Vite 环境变量文档](https://vitejs.dev/guide/env-and-mode.html)
- [加密中间件文档](./CRYPTO_MIDDLEWARE.md)
- [快速参考](./CRYPTO_README.md)

## 🆘 常见问题

### Q: 为什么我的环境变量不生效？

A: 检查以下几点：

1. 变量名是否以 `VITE_` 开头
2. 修改后是否重启了开发服务器
3. 是否有其他配置文件覆盖了该值

### Q: 如何在不同环境使用不同的密钥？

A: 在对应的环境文件中设置不同的密钥：

- `.env.development` - 开发密钥
- `.env.production` - 生产密钥
- `.env.test` - 测试密钥

### Q: 可以在运行时修改环境变量吗？

A: 不可以。环境变量在构建时确定，运行时无法修改。

### Q: 如何让团队成员快速配置？

A: 提供 `.env.example` 文件，新成员复制并修改即可：

```bash
cp .env.example .env.local
```

---

**总结**：

- ✅ 开发环境配置已完成（禁用加密）
- ✅ 生产环境配置已完成（启用加密）
- ⚠️ 生产环境需要修改默认密钥
- 📝 所有配置文件已准备就绪
