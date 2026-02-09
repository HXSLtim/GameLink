# GameLink 加密配置指南

**任务ID：** #18
**更新日期：** 2026-02-09
**维护人：** DevOps-Engineer

---

## 概述

本文档说明如何正确配置 GameLink 的前后端加密功能，确保加密通讯正常工作。

## 关键注意事项

⚠️ **前后端密钥格式不同！**

- **后端使用 base64 编码的密钥**
- **前端使用原始字节字符串的密钥**

不要直接复制后端的密钥到前端！

---

## 快速配置指南

### 方法 1：使用自动化脚本（推荐）

运行密钥生成脚本，自动生成所有配置：

```bash
bash scripts/generate-production-keys.sh
```

脚本会自动创建：
- `.env.production` - 后端配置（使用 base64 编码密钥）
- `admin/.env.production` - 前端配置（使用原始字节密钥）
- `backups/keys/` - 密钥备份文件

### 方法 2：手动配置

#### 步骤 1：生成原始密钥

```bash
# 生成 32 字节原始密钥
openssl rand -base64 32 | base64 -d | xxd -p -c 32

# 输出示例：
# 1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b
```

#### 步骤 2：配置前端

将原始密钥（32个字符）直接填入前端配置：

```bash
# admin/.env.production
VITE_CRYPTO_SECRET_KEY=1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b
VITE_CRYPTO_IV=1a2b3c4d5e6f7a8b9c0d1e2f3a4b
```

#### 步骤 3：配置后端

将原始密钥进行 base64 编码后填入后端配置：

```bash
# 将原始密钥转换为 base64
echo -n "1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b" | base64

# 输出示例（44个字符）：
# MWEyYjNjNGQ1ZTZmN2E4YjljMGQxZTJmM2E0YjVjNmQ3ZThmOWEwYjFjMmQzZTRmNWE2YjdjOGQ5ZTBmMWEyYg==

# .env.production
CRYPTO_SECRET_KEY=MWEyYjNjNGQ1ZTZmN2E4YjljMGQxZTJmM2E0YjVjNmQ3ZThmOWEwYjFjMmQzZTRmNWE2YjdjOGQ5ZTBmMWEyYg==
CRYPTO_IV=MWEyYjNjNGQ1ZTZmN2E4YjljMGQxZTJmM2E0YjU=
```

---

## 配置检查清单

### 开发环境（加密关闭）

- [ ] 后端 `CRYPTO_ENABLED=false`
- [ ] 前端 `VITE_CRYPTO_ENABLED=false`
- [ ] 前后端无需配置密钥

### Staging 环境（加密启用）

- [ ] 后端 `CRYPTO_ENABLED=true`
- [ ] 前端 `VITE_CRYPTO_ENABLED=true`
- [ ] 后端 `CRYPTO_SECRET_KEY` 已配置（base64 格式，44字符）
- [ ] 前端 `VITE_CRYPTO_SECRET_KEY` 已配置（原始字节，32字符）
- [ ] 后端 `CRYPTO_IV` 已配置（base64 格式，24字符）
- [ ] 前端 `VITE_CRYPTO_IV` 已配置（原始字节，16字符）
- [ ] 后端 `CRYPTO_USE_SIGNATURE=true`
- [ ] 前端 `VITE_CRYPTO_USE_SIGNATURE=true`

### 生产环境（加密启用）

- [ ] 后端 `CRYPTO_ENABLED=true`
- [ ] 前端 `VITE_CRYPTO_ENABLED=true`
- [ ] 后端 `CRYPTO_SECRET_KEY` 已配置（base64 格式，44字符）
- [ ] 前端 `VITE_CRYPTO_SECRET_KEY` 已配置（原始字节，32字符）
- [ ] 后端 `CRYPTO_IV` 已配置（base64 格式，24字符）
- [ ] 前端 `VITE_CRYPTO_IV` 已配置（原始字节，16字符）
- [ ] 后端 `CRYPTO_USE_SIGNATURE=true`（推荐）
- [ ] 前端 `VITE_CRYPTO_USE_SIGNATURE=true`（推荐）
- [ ] 运行验证测试确保加密通讯正常

---

## 密钥格式对照表

| 配置项 | 后端格式 | 前端格式 | 长度 | 示例 |
|--------|---------|---------|------|------|
| **SECRET_KEY** | base64 编码 | 原始字节 | 44字符 / 32字符 | 后端: `MWEyYjNj...`<br>前端: `1a2b3c...` |
| **IV** | base64 编码 | 原始字节 | 24字符 / 16字符 | 后端: `MWEyYjNj...`<br>前端: `1a2b3c...` |
| **USE_SIGNATURE** | `true`/`false` | `true`/`false` | 布尔值 | 保持一致 |

---

## 验证工具

### 1. 验证配置一致性

```bash
bash scripts/verify-crypto-keys.sh
```

检查内容：
- 密钥长度
- IV长度
- 前后端配置一致性
- 启用状态一致性

### 2. 部署前检查

```bash
bash scripts/pre-deployment-check.sh
```

### 3. 测试加密通讯

**启动后端（启用加密）：**
```bash
# .env.production
CRYPTO_ENABLED=true
CRYPTO_SECRET_KEY=<base64密钥>
CRYPTO_IV=<base64 IV>
CRYPTO_USE_SIGNATURE=true

docker-compose -f docker-compose.prod.yml --env-file .env.production up backend
```

**启动前端（启用加密）：**
```bash
# admin/.env.production
VITE_CRYPTO_ENABLED=true
VITE_CRYPTO_SECRET_KEY=<原始密钥>
VITE_CRYPTO_IV=<原始IV>
VITE_CRYPTO_USE_SIGNATURE=true

cd admin
npm run build
npm run preview
```

**测试登录请求：**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@gamelink.com",
    "password": "Admin123456"
  }'
```

**预期结果：**
- 请求成功返回 JWT token
- 后端日志显示解密成功
- 前端显示登录成功

---

## 常见问题

### Q1: 后端报错 "crypto: key length is not supported"

**原因：** 后端密钥格式错误

**解决方案：**
1. 检查 `CRYPTO_SECRET_KEY` 是否为 base64 编码（44字符）
2. 检查 `CRYPTO_IV` 是否为 base64 编码（24字符）
3. 使用脚本重新生成：`bash scripts/generate-production-keys.sh`

### Q2: 前端报错 "CryptoConfigError: Invalid key length"

**原因：** 前端密钥格式错误

**解决方案：**
1. 检查 `VITE_CRYPTO_SECRET_KEY` 是否为原始字节（32字符）
2. 检查 `VITE_CRYPTO_IV` 是否为原始字节（16字符）
3. 不要直接复制后端的 base64 密钥

### Q3: 签名验证失败

**原因：** 前后端签名配置不一致

**解决方案：**
1. 检查 `CRYPTO_USE_SIGNATURE` 和 `VITE_CRYPTO_USE_SIGNATURE` 是否一致
2. 检查密钥是否完全匹配
3. 检查时间戳是否同步

### Q4: 开发环境需要加密吗？

**建议：** 不需要

**原因：**
- 开发环境通常在受信任网络中
- 加密会增加调试难度
- 生产环境启用即可

### Q5: 如何切换密钥？

**方法：**
1. 生成新的密钥对
2. 同时更新前后端配置
3. 重启前后端服务
4. 验证加密通讯

**注意：** 密钥切换会导致所有旧会话失效，建议在低峰期操作。

---

## 安全建议

### 密钥管理

1. **永远不要将密钥提交到版本控制**
   ```bash
   # .gitignore
   .env.production
   admin/.env.production
   backups/keys/
   ```

2. **定期轮换密钥**
   - 建议每 3-6 个月轮换一次
   - 使用密钥版本管理（未来实现）

3. **安全存储密钥**
   - 生产环境：使用 Kubernetes Secrets / AWS Secrets Manager
   - 测试环境：使用环境变量文件（限制访问权限）
   - 密钥文件权限：`chmod 400`

### 加密最佳实践

1. **生产环境必须启用加密**
   - `CRYPTO_ENABLED=true`
   - `CRYPTO_USE_SIGNATURE=true`

2. **使用强密钥**
   - 使用 openssl 或专用工具生成
   - 不要使用弱密钥（如 "12345678"）

3. **定期审计**
   - 运行验证工具检查配置
   - 检查后端日志是否有解密失败记录
   - 监控异常请求模式

---

## 参考资料

- **加密验证报告：** `docs/ENCRYPTION_VERIFICATION_REPORT.md`
- **依赖配置文档：** `docs/DEPENDENCIES_AND_CONFIG.md`
- **安全加固指南：** `docs/SECURITY_HARDENING.md`
- **部署检查清单：** `docs/DEPLOYMENT_CHECKLIST.md`

---

## 更新历史

| 日期 | 版本 | 更新内容 | 更新人 |
|------|------|---------|--------|
| 2026-02-09 | 1.0.0 | 初始版本 | DevOps-Engineer |
