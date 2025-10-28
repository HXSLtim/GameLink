# 加密中间件文档

## 📖 概述

本项目实现了完整的前后端通信加密方案，使用 AES-256-CBC 算法对敏感数据进行加密，并提供签名验证功能防止数据篡改。

## 🔒 加密方式

### 加密算法

- **对称加密**: AES-256-CBC
- **哈希算法**: SHA-256（签名）、MD5（可选）
- **模式**: CBC（密码块链接模式）
- **填充**: PKCS7

### 加密流程

#### 请求加密（前端 → 后端）

```typescript
// 原始请求数据
{
  username: "admin",
  password: "123456"
}

// 加密后的请求数据
{
  encrypted: true,
  payload: "U2FsdGVkX1+...", // Base64 加密字符串
  timestamp: 1703001234567,
  signature: "abc123..." // SHA-256 签名
}
```

#### 响应解密（后端 → 前端）

```typescript
// 加密的响应数据
{
  encrypted: true,
  payload: "U2FsdGVkX1+...",
  timestamp: 1703001234567,
  signature: "abc123..."
}

// 解密后的响应数据
{
  success: true,
  data: { ... },
  message: "success"
}
```

## 🚀 使用方法

### 1. 安装依赖

```bash
npm install crypto-js
npm install -D @types/crypto-js
```

### 2. 配置环境变量

创建 `.env.local` 文件：

```bash
# 启用加密
VITE_CRYPTO_ENABLED=true

# 加密密钥（32字节，必须与后端一致）
VITE_CRYPTO_SECRET_KEY=your-secret-key-here-32-bytes

# 加密向量（16字节，必须与后端一致）
VITE_CRYPTO_IV=your-iv-16-bytes
```

### 3. 自动加密

加密中间件已集成到 `apiClient`，无需额外配置：

```typescript
// 普通请求会自动加密
const result = await apiClient.post('/api/v1/auth/login', {
  username: 'admin',
  password: '123456',
});

// 响应会自动解密
console.log(result); // 解密后的数据
```

### 4. 手动加密/解密

```typescript
import { CryptoUtil } from 'utils/crypto';

// 加密
const encrypted = CryptoUtil.encrypt({ password: '123456' });

// 解密
const decrypted = CryptoUtil.decrypt(encrypted);

// 部分字段加密
const data = CryptoUtil.encryptFields({ username: 'admin', password: '123456' }, ['password']);

// 生成签名
const signature = CryptoUtil.generateSignature(data, Date.now());
```

## ⚙️ 配置选项

### 中间件配置

```typescript
import { cryptoMiddleware } from 'middleware/crypto';

// 查看当前配置
const config = cryptoMiddleware.getConfig();

// 更新配置
cryptoMiddleware.updateConfig({
  enabled: true,
  methods: ['POST', 'PUT', 'PATCH'],
  excludePaths: ['/api/v1/health'],
  useSignature: true,
});

// 启用/禁用加密
cryptoMiddleware.enable();
cryptoMiddleware.disable();
```

### 环境配置

| 变量                     | 说明               | 默认值                           |
| ------------------------ | ------------------ | -------------------------------- |
| `VITE_CRYPTO_ENABLED`    | 是否启用加密       | `true`                           |
| `VITE_CRYPTO_SECRET_KEY` | AES 密钥（32字节） | `GameLink2025SecretKey!@#123456` |
| `VITE_CRYPTO_IV`         | AES 向量（16字节） | `GameLink2025IV!!!`              |

## 🎯 加密策略

### 默认策略

1. **请求方法**: 仅加密 `POST`、`PUT`、`PATCH` 请求
2. **路径白名单**: 不加密健康检查等特殊接口
3. **签名验证**: 默认启用，防止数据篡改

### 两种加密模式

#### 模式 1：全量加密（推荐）

加密整个请求/响应体：

```typescript
// 优点：安全性高，实现简单
// 缺点：性能开销稍大

config.data = {
  encrypted: true,
  payload: encryptedData, // 整个 data 加密
  timestamp: Date.now(),
  signature: '...',
};
```

#### 模式 2：部分字段加密

仅加密敏感字段：

```typescript
// 优点：性能好，灵活
// 缺点：需要明确指定字段

const data = CryptoUtil.encryptFields(
  { username: 'admin', password: '123456' },
  ['password'], // 仅加密 password
);
```

## 🔐 安全建议

### 密钥管理

1. **生产环境**：使用强随机密钥，不要使用默认值
2. **密钥轮换**：定期更换密钥（建议每季度）
3. **密钥存储**：使用环境变量，不要硬编码
4. **前后端同步**：确保前后端密钥完全一致

### 生成强密钥

```bash
# 生成 32 字节密钥
node -e "console.log(require('crypto').randomBytes(32).toString('base64'))"

# 生成 16 字节向量
node -e "console.log(require('crypto').randomBytes(16).toString('base64'))"
```

### 部署检查清单

- [ ] 修改默认密钥和向量
- [ ] 前后端密钥保持一致
- [ ] 启用 HTTPS（防止中间人攻击）
- [ ] 启用签名验证
- [ ] 设置合理的请求白名单
- [ ] 监控加密/解密性能
- [ ] 定期审查加密日志

## 🧪 测试

### 单元测试

```typescript
import { CryptoUtil } from 'utils/crypto';

describe('CryptoUtil', () => {
  it('should encrypt and decrypt correctly', () => {
    const original = { test: 'data' };
    const encrypted = CryptoUtil.encrypt(original);
    const decrypted = CryptoUtil.decrypt(encrypted);

    expect(decrypted).toEqual(original);
  });

  it('should generate valid signature', () => {
    const data = { test: 'data' };
    const timestamp = Date.now();
    const signature = CryptoUtil.generateSignature(data, timestamp);

    expect(signature).toBeTruthy();
    expect(signature.length).toBe(64); // SHA-256 = 64 hex chars
  });
});
```

### 集成测试

```typescript
import { apiClient } from 'api/client';

describe('Crypto Middleware', () => {
  it('should encrypt request and decrypt response', async () => {
    const response = await apiClient.post('/api/v1/auth/login', {
      username: 'test',
      password: 'test123',
    });

    expect(response).toBeDefined();
    // 响应应该是解密后的数据
  });
});
```

## 🐛 调试

### 启用日志

加密中间件会自动输出日志：

```
🔒 请求数据已加密: { url: '/api/v1/auth/login', method: 'POST', timestamp: ... }
🔓 响应数据已加密，开始解密...
✅ 签名验证通过
✅ 响应数据解密成功
```

### 临时禁用加密

```typescript
// 方法 1: 环境变量
VITE_CRYPTO_ENABLED = false;

// 方法 2: 代码
cryptoMiddleware.disable();
```

## 📊 性能影响

| 操作          | 耗时   | 影响   |
| ------------- | ------ | ------ |
| 加密 1KB 数据 | ~1ms   | 可忽略 |
| 解密 1KB 数据 | ~1ms   | 可忽略 |
| 签名生成      | ~0.5ms | 可忽略 |
| 签名验证      | ~0.5ms | 可忽略 |

**建议**：

- 小数据（< 10KB）：全量加密
- 大数据（> 10KB）：部分字段加密或压缩后加密

## 🔄 后端对接

后端需要实现对应的解密中间件，示例（Go）：

```go
// 解密请求
func DecryptMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        var req EncryptedRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            return
        }

        if req.Encrypted {
            // 解密
            decrypted := Decrypt(req.Payload, secretKey, iv)
            // 验证签名
            if !VerifySignature(decrypted, req.Timestamp, req.Signature) {
                c.JSON(400, gin.H{"error": "invalid signature"})
                return
            }
            // 替换请求体
            c.Set("decrypted_data", decrypted)
        }

        c.Next()
    }
}
```

## 📚 相关文件

- 加密工具：`src/utils/crypto.ts`
- 加密中间件：`src/middleware/crypto.ts`
- API 客户端：`src/api/client.ts`
- 环境配置：`.env.example`, `.env.development`

## 🆘 常见问题

### Q: 加密后请求失败？

A: 检查前后端密钥是否一致，后端是否实现了解密中间件。

### Q: 性能影响大吗？

A: 对于小数据（< 10KB）几乎无影响，大数据建议使用部分加密。

### Q: 可以只加密敏感字段吗？

A: 可以，修改中间件配置使用模式 2（部分字段加密）。

### Q: 如何在开发环境禁用加密？

A: 设置 `VITE_CRYPTO_ENABLED=false` 或调用 `cryptoMiddleware.disable()`。

## 📝 更新日志

- **v1.0.0** (2025-01-28)
  - 初始版本
  - 支持 AES-256-CBC 加密
  - 支持 SHA-256 签名验证
  - 支持全量/部分加密模式
