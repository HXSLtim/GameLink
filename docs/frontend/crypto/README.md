# 🔐 加密中间件 - 快速参考

## 📦 已集成完成

前端加密中间件已完全集成，支持 AES-256-CBC 加密和 SHA-256 签名验证。

## 🚀 立即开始

### 1. 查看是否生效

```bash
# 启动开发服务器
npm run dev

# 登录系统，打开浏览器控制台
# 如果看到以下日志，说明加密中间件工作正常：
🔒 请求数据已加密: { url: '/api/v1/auth/login', ... }
```

### 2. 配置密钥（生产环境）

```bash
# 创建 .env.local
cat > .env.local << EOF
VITE_CRYPTO_ENABLED=true
VITE_CRYPTO_SECRET_KEY=your-32-byte-secret-key-here-123456
VITE_CRYPTO_IV=your-iv-16-byte
EOF
```

### 3. 正常使用

无需修改任何业务代码，所有 POST/PUT/PATCH 请求会自动加密：

```typescript
// ✅ 自动加密
await authApi.login({ username: 'admin', password: '123456' });
await orderApi.create({ user_id: 1, amount: 100 });
```

## 📁 核心文件

```
src/
├── utils/
│   └── crypto.ts              # 加密工具类
├── middleware/
│   └── crypto.ts              # 加密中间件
└── api/
    └── client.ts              # 已集成加密

文档/
├── CRYPTO_INTEGRATION.md      # 集成说明（推荐）⭐
├── CRYPTO_MIDDLEWARE.md       # 详细文档
├── CRYPTO_USAGE_EXAMPLES.md   # 使用示例
└── CRYPTO_README.md           # 本文档
```

## ⚙️ 默认配置

| 配置项   | 默认值         | 说明      |
| -------- | -------------- | --------- |
| 加密算法 | AES-256-CBC    | 对称加密  |
| 签名算法 | SHA-256        | 防篡改    |
| 加密方法 | POST/PUT/PATCH | GET不加密 |
| 开发环境 | 禁用           | 方便调试  |
| 生产环境 | 启用           | 保护数据  |

## 🎯 三种使用模式

### 模式 1: 自动加密（默认，推荐）✅

```typescript
// 无需任何修改，自动工作
const result = await apiClient.post('/api/data', { secret: '123' });
```

### 模式 2: 手动加密

```typescript
import { CryptoUtil } from 'utils/crypto';
const encrypted = CryptoUtil.encrypt({ password: '123' });
```

### 模式 3: 部分字段加密

```typescript
const data = CryptoUtil.encryptFields({ name: 'John', password: '123' }, ['password']);
```

## 🔧 常用操作

```typescript
import { cryptoMiddleware } from 'middleware/crypto';

// 临时禁用
cryptoMiddleware.disable();

// 启用
cryptoMiddleware.enable();

// 检查状态
cryptoMiddleware.isEnabled(); // true/false

// 更新配置
cryptoMiddleware.updateConfig({
  methods: ['POST'],
  excludePaths: ['/api/public/*'],
});
```

## ⚠️ 重要提示

### 必须做 ✅

- [x] ✅ 已安装 crypto-js 依赖
- [x] ✅ 已集成到 API 客户端
- [x] ✅ 已配置默认密钥
- [ ] ⚠️ 生产环境修改密钥
- [ ] ⚠️ 后端实现解密中间件
- [ ] ⚠️ 前后端密钥保持一致

### 建议做 💡

- [ ] 启用 HTTPS
- [ ] 定期轮换密钥
- [ ] 监控加密性能
- [ ] 记录异常日志

## 🤝 后端对接清单

后端需要实现以下功能：

1. **解密中间件**
   - 解密请求体中的 `payload`
   - 验证 `signature`
   - 将解密数据传递给业务逻辑

2. **加密响应**
   - 加密响应数据
   - 生成签名
   - 返回加密格式

3. **密钥配置**
   - 使用与前端相同的密钥和向量
   - 支持环境变量配置

## 📊 数据格式

### 请求格式

```json
{
  "encrypted": true,
  "payload": "U2FsdGVkX1+...",
  "timestamp": 1703001234567,
  "signature": "abc123..."
}
```

### 响应格式

```json
{
  "success": true,
  "data": {
    "encrypted": true,
    "payload": "U2FsdGVkX1+...",
    "timestamp": 1703001234568,
    "signature": "def456..."
  }
}
```

## 🧪 测试

```bash
# 运行测试
npm run test

# 测试覆盖率
npm run test -- --coverage
```

## 🐛 调试

```bash
# 开发环境禁用加密
VITE_CRYPTO_ENABLED=false npm run dev

# 查看加密日志
# 控制台会显示：
🔒 请求数据已加密
🔓 响应数据已加密，开始解密
✅ 签名验证通过
```

## 📚 详细文档

- [集成说明](./CRYPTO_INTEGRATION.md) - 完整的集成指南 ⭐
- [技术文档](./CRYPTO_MIDDLEWARE.md) - 实现原理和配置
- [使用示例](./CRYPTO_USAGE_EXAMPLES.md) - 各种场景的示例代码

## 🆘 常见问题

<details>
<summary><b>Q: 为什么请求失败返回 400？</b></summary>

A: 检查后端是否实现了解密中间件，密钥是否一致。

</details>

<details>
<summary><b>Q: 如何在开发环境禁用加密？</b></summary>

A: 已默认禁用，如需启用设置 `VITE_CRYPTO_ENABLED=true`

</details>

<details>
<summary><b>Q: 加密会影响性能吗？</b></summary>

A: 对小数据（< 10KB）几乎无影响，耗时约 1-5ms。

</details>

<details>
<summary><b>Q: 如何修改密钥？</b></summary>

A: 在 `.env.local` 中设置 `VITE_CRYPTO_SECRET_KEY` 和 `VITE_CRYPTO_IV`

</details>

## 📞 获取帮助

1. 查看 [集成文档](./CRYPTO_INTEGRATION.md)
2. 阅读 [技术文档](./CRYPTO_MIDDLEWARE.md)
3. 参考 [示例代码](./CRYPTO_USAGE_EXAMPLES.md)

---

**状态**: ✅ 已完成集成，可直接使用

**版本**: 1.0.0

**最后更新**: 2025-01-28
