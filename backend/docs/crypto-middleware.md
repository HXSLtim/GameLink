# 加密中间件（后端）

本文档说明 GameLink 后端对前端 AES 加密中间件的对接方式、配置项与调试要点。

## 功能概览

- 支持与前端一致的 AES-CBC 对称加密协议（默认秘钥与 IV 与前端示例同步）。
- 仅对指定 HTTP 方法（默认：POST/PUT/PATCH）执行解密，其余请求直接透传。
- 可配置请求路径白名单，命中白名单时跳过解密逻辑。
- 支持基于原始明文 + 时间戳 + 秘钥的 SHA256 签名校验，防止数据被篡改。
- 解密成功后将明文请求体写回 `gin.Context`，路由和业务层逻辑保持原有写法。

## 配置项

通过 `configs/config.<env>.yaml` 或环境变量设置以下参数：

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `crypto.enabled` | `false`（开发） / `true`（生产样例） | `CRYPTO_ENABLED` | 是否启用加密中间件 |
| `crypto.secret_key` | `GameLink2025SecretKey!@#` | `CRYPTO_SECRET_KEY` | AES 秘钥，长度需为 16/24/32 字节 |
| `crypto.iv` | `GameLink2025IV!!!` | `CRYPTO_IV` | AES 向量，至少 16 字节，将截取前 16 字节参与计算 |
| `crypto.methods` | `["POST","PUT","PATCH"]` | `CRYPTO_METHODS` | 需要解密的 HTTP 方法，逗号分隔 |
| `crypto.exclude_paths` | `["/api/v1/health","/api/v1/ping","/api/v1/auth/refresh"]` | `CRYPTO_EXCLUDE_PATHS` | 白名单路径，命中后跳过解密（子串匹配） |
| `crypto.use_signature` | `true` | `CRYPTO_USE_SIGNATURE` | 是否校验签名 |

> ⚠️ 生产环境请务必通过环境变量覆盖 `secret_key` 与 `iv`，并定期轮换。

## 请求报文格式

当前端命中需要加密的请求时，报文结构如下：

```json
{
  "encrypted": true,
  "payload": "<Base64 编码的 AES 密文>",
  "timestamp": 1703001234567,
  "signature": "<SHA256 签名>"
}
```

后端中间件会：

1. Base64 解码 `payload`，使用 CBC 模式 + PKCS7 去填充得到原始明文；
2. （可选）计算 `sha256(明文 + timestamp + secret_key)` 与 `signature` 比对；
3. 将明文回写至 `gin.Context`，供后续路由处理。

若解析失败、签名校验不通过或缺失必要字段，中间件将直接返回统一错误响应：

```json
{
  "success": false,
  "code": 400,
  "message": "请求数据解密失败"
}
```

## 签名策略

- 签名由前端在加密前计算，格式为 `sha256(JSON.stringify(data) + timestamp + secret_key)`。
- 后端默认开启签名校验，可按需通过配置关闭（例如联调或灰度阶段）。
- 请求需携带 `timestamp` 字段（毫秒时间戳），用于生成签名与追踪日志。

## 白名单与调试

- 白名单采用子串匹配，例如 `/api/v1/ping` 可覆盖 `/api/v1/ping/status`。
- 路径被白名单命中后，请求体不会被解密，业务层会收到原始报文。
- 可通过在日志或断点中读取 `ctx.Get("crypto.encrypted")` 来判断请求是否经由加密流程。

## 启用步骤

1. 设置密钥与向量：
   ```bash
   export CRYPTO_SECRET_KEY="$(openssl rand -base64 24)"   # 生成 32 字节可用密钥
   export CRYPTO_IV="$(openssl rand -base64 12)"           # 取前 16 字节作为 IV
   ```
2. 打开中间件：
   ```bash
   export CRYPTO_ENABLED=true
   ```
3. 重启 `user-service`，观察启动日志确认：
   ```
   crypto middleware enabled, methods=[POST PUT PATCH] exclude=[...] use_signature=true
   ```
4. 使用前端集成或 Postman 的加密脚本发起测试请求，确认能成功解密并进入业务逻辑。

## 常见问题

| 问题 | 排查建议 |
|------|----------|
| 返回 400 且提示“请求数据解密失败” | 检查是否使用了正确的密钥和 IV、payload 是否为合法 Base64 字符串 |
| 提示“请求签名验证失败” | 检查前后端签名算法是否一致、timestamp 是否为发送时刻生成、密钥是否一致 |
| 正常请求被误解密 | 确认 `crypto.methods` 是否包含该 HTTP 方法，或将路径加入 `exclude_paths` |

如需同时对响应加密，可在服务层组装响应后统一封装，复用 `Crypto` 模块中的签名方法，确保与前端保持一致的协议。
