# 加密中间件使用示例

## 📖 快速开始

加密中间件已经集成到项目中，默认情况下会自动工作，无需额外配置。

## 🔧 基础配置

### 1. 环境变量配置

```bash
# .env.local 或 .env.production
VITE_CRYPTO_ENABLED=true
VITE_CRYPTO_SECRET_KEY=your-32-byte-secret-key-here
VITE_CRYPTO_IV=your-16-byte-iv
```

### 2. 开发环境禁用加密

```bash
# .env.development
VITE_CRYPTO_ENABLED=false
```

## 💻 使用示例

### 示例 1: 登录请求（自动加密）

```typescript
// src/pages/Login/Login.tsx
import { useAuth } from 'contexts/AuthContext';

export const Login = () => {
  const { login } = useAuth();

  const handleSubmit = async () => {
    try {
      // ✅ 请求会自动加密
      await login('admin', 'password123');

      // 实际发送的数据：
      // {
      //   encrypted: true,
      //   payload: "U2FsdGVkX1+...",
      //   timestamp: 1703001234567,
      //   signature: "abc123..."
      // }
    } catch (error) {
      console.error('登录失败', error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      {/* ... */}
    </form>
  );
};
```

### 示例 2: 创建订单（自动加密）

```typescript
// src/pages/Orders/CreateOrder.tsx
import { orderApi } from 'services/api/order';

export const CreateOrder = () => {
  const handleCreate = async () => {
    const orderData = {
      user_id: 1,
      game_id: 10,
      title: '王者荣耀陪玩',
      price_cents: 5000,
    };

    try {
      // ✅ 请求会自动加密
      const result = await orderApi.create(orderData);

      // ✅ 响应会自动解密
      console.log('创建成功:', result);
    } catch (error) {
      console.error('创建失败', error);
    }
  };

  return <button onClick={handleCreate}>创建订单</button>;
};
```

### 示例 3: 手动加密敏感数据

```typescript
// 在某些特殊场景下需要手动加密
import { CryptoUtil } from 'utils/crypto';

// 加密单个值
const encryptedPassword = CryptoUtil.encrypt('password123');
localStorage.setItem('pwd', encryptedPassword);

// 解密
const decryptedPassword = CryptoUtil.decrypt<string>(localStorage.getItem('pwd')!);

// 加密对象
const userData = {
  username: 'admin',
  password: 'secret',
  phone: '13800138000',
};

const encrypted = CryptoUtil.encrypt(userData);
// 发送到服务器或存储...

const decrypted = CryptoUtil.decrypt<typeof userData>(encrypted);
```

### 示例 4: 部分字段加密

```typescript
import { CryptoUtil } from 'utils/crypto';

// 仅加密敏感字段
const formData = {
  name: 'John Doe',
  email: 'john@example.com',
  password: 'secret123',
  phone: '13800138000',
};

// 只加密 password 和 phone
const encryptedData = CryptoUtil.encryptFields(formData, ['password', 'phone']);

console.log(encryptedData);
// {
//   name: 'John Doe',
//   email: 'john@example.com',
//   password: 'U2FsdGVkX1+...',  // 已加密
//   phone: 'U2FsdGVkX1+...'      // 已加密
// }
```

### 示例 5: 生成和验证签名

```typescript
import { CryptoUtil } from 'utils/crypto';

// 发送方：生成签名
const requestData = {
  orderId: 12345,
  amount: 100,
};

const timestamp = Date.now();
const signature = CryptoUtil.generateSignature(requestData, timestamp);

// 发送数据
const payload = {
  data: requestData,
  timestamp,
  signature,
};

// 接收方：验证签名
const receivedData = payload.data;
const receivedTimestamp = payload.timestamp;
const receivedSignature = payload.signature;

const expectedSignature = CryptoUtil.generateSignature(receivedData, receivedTimestamp);

if (expectedSignature === receivedSignature) {
  console.log('✅ 签名验证通过，数据未被篡改');
} else {
  console.error('❌ 签名验证失败，数据可能被篡改');
}
```

## ⚙️ 高级配置

### 自定义加密配置

```typescript
// src/api/customClient.ts
import axios from 'axios';
import { CryptoMiddleware } from 'middleware/crypto';

// 创建自定义加密中间件
const customCrypto = new CryptoMiddleware({
  enabled: true,
  methods: ['POST', 'PUT'], // 只加密 POST 和 PUT
  excludePaths: [
    '/api/v1/health',
    '/api/v1/public/*', // 公开接口不加密
  ],
  sensitiveFields: ['password', 'phone', 'id_card'],
  useSignature: true,
});

// 创建自定义 API 客户端
const customClient = axios.create({
  baseURL: 'https://api.example.com',
});

// 应用加密中间件
customClient.interceptors.request.use(customCrypto.requestInterceptor, (error) =>
  Promise.reject(error),
);

customClient.interceptors.response.use(customCrypto.responseInterceptor, (error) =>
  Promise.reject(error),
);
```

### 动态控制加密

```typescript
import { cryptoMiddleware } from 'middleware/crypto';

// 临时禁用加密
cryptoMiddleware.disable();
await someApiCall();
cryptoMiddleware.enable();

// 或使用开关
if (isDevelopment) {
  cryptoMiddleware.disable();
}

// 更新配置
cryptoMiddleware.updateConfig({
  methods: ['POST'], // 只加密 POST 请求
  excludePaths: ['/api/v1/upload'], // 上传接口不加密
});

// 检查状态
if (cryptoMiddleware.isEnabled()) {
  console.log('加密已启用');
}
```

## 🧪 测试

### 单元测试示例

```typescript
// src/middleware/crypto.test.ts
import { describe, it, expect, beforeEach } from 'vitest';
import { CryptoMiddleware } from './crypto';
import { InternalAxiosRequestConfig } from 'axios';

describe('CryptoMiddleware', () => {
  let middleware: CryptoMiddleware;

  beforeEach(() => {
    middleware = new CryptoMiddleware({
      enabled: true,
      methods: ['POST'],
      excludePaths: ['/health'],
    });
  });

  it('should encrypt POST request data', () => {
    const config: InternalAxiosRequestConfig = {
      method: 'POST',
      url: '/api/users',
      data: { username: 'test', password: 'secret' },
    } as InternalAxiosRequestConfig;

    const result = middleware.requestInterceptor(config);

    expect(result.data.encrypted).toBe(true);
    expect(result.data.payload).toBeTruthy();
    expect(result.data.timestamp).toBeTruthy();
  });

  it('should not encrypt GET requests', () => {
    const config: InternalAxiosRequestConfig = {
      method: 'GET',
      url: '/api/users',
    } as InternalAxiosRequestConfig;

    const result = middleware.requestInterceptor(config);

    expect(result.data).toBeUndefined();
  });

  it('should not encrypt excluded paths', () => {
    const config: InternalAxiosRequestConfig = {
      method: 'POST',
      url: '/health',
      data: { status: 'ok' },
    } as InternalAxiosRequestConfig;

    const result = middleware.requestInterceptor(config);

    expect(result.data.encrypted).toBeUndefined();
    expect(result.data.status).toBe('ok');
  });
});
```

## 🔍 调试技巧

### 查看加密/解密日志

加密中间件会自动输出详细日志：

```typescript
// 请求加密
🔒 请求数据已加密: {
  url: '/api/v1/auth/login',
  method: 'POST',
  timestamp: 1703001234567
}

// 响应解密
🔓 响应数据已加密，开始解密...
✅ 签名验证通过
✅ 响应数据解密成功
```

### 在 Chrome DevTools 中查看

打开 Network 面板，查看实际发送的请求：

```json
// Request Payload
{
  "encrypted": true,
  "payload": "U2FsdGVkX1+abcdef123456...",
  "timestamp": 1703001234567,
  "signature": "abc123def456..."
}

// Response
{
  "encrypted": true,
  "payload": "U2FsdGVkX1+zyxwvu987654...",
  "timestamp": 1703001234568,
  "signature": "def789ghi012..."
}
```

### 禁用加密进行调试

```typescript
// 方法 1: 环境变量
VITE_CRYPTO_ENABLED=false npm run dev

// 方法 2: 代码
import { cryptoMiddleware } from 'middleware/crypto';
cryptoMiddleware.disable();
```

## 🚨 错误处理

### 常见错误和解决方案

```typescript
// 1. 解密失败
try {
  const decrypted = CryptoUtil.decrypt(encryptedData);
} catch (error) {
  console.error('解密失败，可能原因：');
  console.error('- 密钥不匹配');
  console.error('- 数据已损坏');
  console.error('- 数据不是加密格式');
}

// 2. 签名验证失败
try {
  // 响应拦截器会自动验证签名
  const response = await apiClient.post('/api/data', data);
} catch (error) {
  if (error.message.includes('签名验证失败')) {
    console.error('数据可能被篡改！');
    // 通知用户或上报
  }
}

// 3. 加密超时
// 对于大数据，可能需要增加超时时间
const largeClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000, // 30 秒
});
```

## 📊 性能优化

### 1. 选择性加密

```typescript
// 只加密真正敏感的接口
cryptoMiddleware.updateConfig({
  methods: ['POST'],
  excludePaths: [
    '/api/v1/health',
    '/metrics',
    '/api/v1/logs', // 日志接口不加密
  ],
});
```

### 2. 缓存加密结果

```typescript
const encryptionCache = new Map<string, string>();

function getCachedEncryption(data: string): string {
  const key = CryptoUtil.md5(data);

  if (encryptionCache.has(key)) {
    return encryptionCache.get(key)!;
  }

  const encrypted = CryptoUtil.encrypt(data);
  encryptionCache.set(key, encrypted);

  return encrypted;
}
```

### 3. 异步加密大数据

```typescript
// 对于大数据，可以使用 Web Worker
async function encryptLargeData(data: unknown): Promise<string> {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve(CryptoUtil.encrypt(data));
    }, 0);
  });
}
```

## 🔐 安全最佳实践

### 1. 密钥管理

```typescript
// ❌ 错误：硬编码密钥
const SECRET_KEY = 'my-secret-key';

// ✅ 正确：使用环境变量
const SECRET_KEY = import.meta.env.VITE_CRYPTO_SECRET_KEY;

// ✅ 更好：从安全的密钥管理服务获取
const SECRET_KEY = await fetchSecretKey();
```

### 2. HTTPS + 加密

```typescript
// 加密不能替代 HTTPS！
// 应该同时使用：
// - HTTPS：防止中间人攻击
// - 数据加密：保护敏感数据
// - 签名验证：防止数据篡改
```

### 3. 定期轮换密钥

```typescript
// 实现密钥版本控制
interface EncryptedPayload {
  encrypted: true;
  payload: string;
  timestamp: number;
  signature: string;
  keyVersion: number; // 密钥版本
}

// 根据版本使用不同的密钥
function getKeyByVersion(version: number): string {
  const keys = {
    1: 'old-key',
    2: 'current-key',
    3: 'new-key',
  };
  return keys[version] || keys[2];
}
```

## 📚 相关资源

- [加密中间件文档](./CRYPTO_MIDDLEWARE.md)
- [加密工具源码](./src/utils/crypto.ts)
- [中间件源码](./src/middleware/crypto.ts)
- [CryptoJS 文档](https://cryptojs.gitbook.io/docs/)

## 💡 提示

1. 开发环境建议禁用加密，方便调试
2. 生产环境必须启用加密和 HTTPS
3. 定期更新密钥，建议每季度一次
4. 监控加密性能，避免影响用户体验
5. 记录所有加密/解密错误，及时发现问题
