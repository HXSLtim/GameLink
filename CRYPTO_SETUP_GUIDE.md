# 🔐 GameLink 加密配置指南

## 概述

GameLink 生产环境使用 AES-256-CBC 加密算法对敏感的 API 请求进行加密，确保数据传输安全。

## 加密机制

### 后端加密中间件

- **算法**: AES-256-CBC
- **密钥长度**: 256 位（32 字节）
- **IV 长度**: 128 位（16 字节）
- **填充方式**: PKCS7
- **签名算法**: SHA-256

### 加密流程

1. **前端**:
   - 将请求数据序列化为 JSON 字符串
   - 使用 AES-256-CBC 加密数据
   - 生成时间戳
   - 使用 SHA-256 生成签名（可选）
   - 发送加密请求体

2. **后端**:
   - 接收加密请求
   - 验证签名（如果启用）
   - 解密请求体
   - 处理业务逻辑

### 加密请求格式

```json
{
  "encrypted": true,
  "payload": "base64_encoded_encrypted_data",
  "timestamp": 1702345678901,
  "signature": "sha256_hex_signature"
}
```

## 配置步骤

### 1. 后端配置

后端配置文件 `backend/configs/config.production.yaml`:

```yaml
crypto:
  enabled: true
  secret_key: ""  # 通过环境变量 CRYPTO_SECRET_KEY 提供
  iv: ""          # 通过环境变量 CRYPTO_IV 提供
  methods:
    - POST
    - PUT
    - PATCH
  exclude_paths:
    - "/api/v1/health"
    - "/api/v1/ping"
    - "/api/v1/auth/refresh"
  use_signature: true
```

环境变量（`.env` 文件）:

```bash
CRYPTO_SECRET_KEY=your_32_character_secret_key_here
CRYPTO_IV=your_16_char_iv
```

### 2. 前端配置

#### 安装依赖

```bash
cd frontend
npm install crypto-js
npm install --save-dev @types/crypto-js
```

#### 环境变量配置

前端配置文件 `frontend/.env.production`:

```bash
VITE_CRYPTO_ENABLED=true
VITE_CRYPTO_SECRET_KEY=your_32_character_secret_key_here
VITE_CRYPTO_IV=your_16_char_iv
VITE_CRYPTO_USE_SIGNATURE=true
```

⚠️ **重要**: 前端的密钥必须与后端完全一致！

### 3. 同步密钥

使用提供的脚本自动同步密钥：

```powershell
.\scripts\sync-crypto-keys.ps1
```

该脚本会：
- 从后端 `.env` 读取 `CRYPTO_SECRET_KEY` 和 `CRYPTO_IV`
- 自动写入前端 `.env.production` 文件

## 部署流程

### 完整部署步骤

```powershell
# 1. 确保后端 .env 文件已配置加密密钥
Get-Content .env | Select-String "CRYPTO"

# 2. 同步密钥到前端
.\scripts\sync-crypto-keys.ps1

# 3. 安装前端依赖
cd frontend
npm install crypto-js @types/crypto-js

# 4. 构建前端
npm run build

# 5. 返回项目根目录
cd ..

# 6. 重新构建 Docker 镜像
docker-compose -f docker-compose.prod.yml build

# 7. 重启服务
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d

# 8. 验证服务状态
docker ps --filter "name=gamelink"
docker logs gamelink-backend --tail=50
```

## 验证加密

### 检查后端日志

```powershell
docker logs gamelink-backend --tail=100 | Select-String "crypto"
```

应该看到：
```
crypto middleware enabled, methods=[POST PUT PATCH] exclude=[...] use_signature=true
```

### 测试加密请求

使用浏览器开发者工具查看网络请求：

1. 打开前端应用
2. 执行一个 POST 请求（如登录）
3. 在 Network 标签中查看请求体

**未加密的请求**:
```json
{
  "email": "admin@gamelink.com",
  "password": "password123"
}
```

**加密后的请求**:
```json
{
  "encrypted": true,
  "payload": "U2FsdGVkX1+...",
  "timestamp": 1702345678901,
  "signature": "a1b2c3d4..."
}
```

## 故障排查

### 问题 1: 前端加密失败

**症状**: 请求返回 400 错误，提示"请求数据解密失败"

**解决方案**:
1. 检查前端环境变量是否正确配置
2. 确认密钥与后端完全一致
3. 检查 crypto-js 是否正确安装

```powershell
# 检查前端环境变量
Get-Content frontend/.env.production | Select-String "CRYPTO"

# 重新安装依赖
cd frontend
npm install crypto-js @types/crypto-js
```

### 问题 2: 签名验证失败

**症状**: 请求返回 400 错误，提示"请求签名验证失败"

**解决方案**:
1. 检查时间戳是否正确生成
2. 确认签名算法实现正确
3. 验证密钥一致性

### 问题 3: 加密未生效

**症状**: 请求体仍然是明文

**解决方案**:
1. 检查 `VITE_CRYPTO_ENABLED` 是否为 `true`
2. 确认前端已重新构建
3. 清除浏览器缓存

```powershell
# 重新构建前端
cd frontend
npm run build

# 重启前端容器
docker-compose -f docker-compose.prod.yml restart frontend
```

## 安全建议

### 密钥管理

1. ✅ **使用强密钥**: 密钥应该是随机生成的，至少 32 字节
2. ✅ **定期轮换**: 建议每 3-6 个月更换一次密钥
3. ✅ **环境隔离**: 开发、测试、生产环境使用不同的密钥
4. ⚠️ **不要提交到 Git**: 确保 `.env` 文件在 `.gitignore` 中

### 生成安全密钥

```powershell
# 生成 32 字节密钥（用于 CRYPTO_SECRET_KEY）
$bytes = New-Object byte[] 32
[Security.Cryptography.RandomNumberGenerator]::Create().GetBytes($bytes)
[Convert]::ToBase64String($bytes).Substring(0, 32)

# 生成 16 字节 IV（用于 CRYPTO_IV）
$bytes = New-Object byte[] 16
[Security.Cryptography.RandomNumberGenerator]::Create().GetBytes($bytes)
[Convert]::ToBase64String($bytes).Substring(0, 16)
```

### HTTPS 配置

⚠️ **重要**: 加密中间件不能替代 HTTPS！

在生产环境中，必须同时配置：
1. ✅ 应用层加密（本文档描述的 AES 加密）
2. ✅ 传输层加密（HTTPS/TLS）

## 性能影响

### 加密开销

- **CPU 使用**: 增加约 5-10%
- **延迟**: 每个请求增加约 1-3ms
- **内存**: 可忽略不计

### 优化建议

1. 只加密敏感请求（POST, PUT, PATCH）
2. 排除不需要加密的路径（如健康检查）
3. 使用 CDN 缓存静态资源

## 开发环境

开发环境建议**禁用加密**以简化调试：

```bash
# frontend/.env.development
VITE_CRYPTO_ENABLED=false
```

```yaml
# backend/configs/config.development.yaml
crypto:
  enabled: false
```

## 相关文档

- [Docker 生产部署](DOCKER_PRODUCTION_DEPLOYED.md)
- [生产环境快速启动](PRODUCTION_QUICK_START.md)
- [后端加密中间件源码](backend/internal/handler/middleware/crypto.go)
- [前端加密工具源码](frontend/src/utils/crypto.ts)

---

**更新时间**: 2025-12-13  
**版本**: 1.0.0
