# 安全配置指南 (Security Configuration Guide)

本文档说明 GameLink 项目的安全配置要求，包括密钥生成、环境变量配置和最佳实践。

## 目录

- [概述](#概述)
- [P0 安全漏洞修复](#p0-安全漏洞修复)
- [密钥生成工具](#密钥生成工具)
- [环境变量配置](#环境变量配置)
- [配置验证规则](#配置验证规则)
- [最佳实践](#最佳实践)

---

## 概述

GameLink 使用 AES-256-CBC 加密保护敏感数据传输，使用 JWT 进行用户认证。所有密钥必须妥善保管，不得泄露或提交到版本控制系统。

### 安全密钥类型

| 密钥类型 | 长度要求 | 用途 | 环境变量 |
|---------|---------|------|---------|
| **加密密钥** | 32 字节 | AES-256-CBC 请求加密 | `CRYPTO_SECRET_KEY` |
| **初始化向量** | 16 字节 | AES 加密 IV | `CRYPTO_IV` |
| **JWT 密钥** | 32+ 字节 | JWT Token 签名 | `JWT_SECRET_KEY` |
| **超级管理员密码** | 8+ 字符 | 管理员登录凭证 | `SUPER_ADMIN_PASSWORD` |

---

## P0 安全漏洞修复

### 修复的安全问题

#### 1. 生产配置密钥为空 ✅

**问题**:
```yaml
# config.production.yaml (修复前)
crypto:
  secret_key: ""  # 空密钥
  iv: ""
```

**修复**:
- 添加启动时验证，加密启用时密钥不能为空
- 生产环境必须通过环境变量配置密钥
- 详见 [api/pkg/config/validate.go](../api/pkg/config/validate.go:72-77)

#### 2. 超级管理员弱密码 ✅

**问题**:
```yaml
# config.development.yaml (修复前)
super_admin:
  password: "123456"  # 弱密码
```

**修复**:
- 开发环境使用随机生成的强密码（24字节）
- 生产环境强制要求 8+ 字符，包含大小写字母、数字、特殊符号
- 详见 [api/pkg/config/validate.go](../api/pkg/config/validate.go:103-139)

#### 3. 开发配置弱密钥 ✅

**问题**:
```yaml
# config.development.yaml (修复前)
crypto:
  secret_key: "GameLink2025SecretKey!@#"  # 仅24字节
  iv: "GameLink2025IV!!!"  # 硬编码默认值
```

**修复**:
- 使用 OpenSSL 生成 32 字节随机密钥
- 使用 OpenSSL 生成 16 字节随机 IV
- 添加密钥长度验证（必须是 16/24/32 字节）
- 检测硬编码默认值并拒绝启动

### 修复涉及的文件

| 文件 | 变更类型 | 说明 |
|------|---------|------|
| `api/configs/config.development.yaml` | 更新 | 使用强密钥和强密码 |
| `api/configs/config.production.yaml` | 更新 | 清空默认值，强制环境变量配置 |
| `api/pkg/config/validate.go` | 增强 | 添加空密钥检查、加强密码强度验证 |
| `api/pkg/config/validate_test.go` | 新增 | 全面的配置验证测试 |
| `scripts/generate-secrets.ps1` | 新增 | Windows 密钥生成工具 |
| `scripts/generate-secrets.sh` | 新增 | Linux/Mac 密钥生成工具 |

---

## 密钥生成工具

项目提供了自动化密钥生成工具，支持 Windows 和 Linux/Mac。

### Windows PowerShell

```powershell
# 运行密钥生成工具
.\scripts\generate-secrets.ps1
```

**输出示例**:
```
========================================
  GameLink 安全密钥生成工具
========================================

生成的密钥如下（请妥善保存，不要泄露）：

1. 加密密钥 (CRYPTO_SECRET_KEY) - 32字节:
   H/oguKMv23lWlivgq8snNZmTzSUp6KSHZnEEo1c0Ook=

2. 初始化向量 (CRYPTO_IV) - 16字节:
   hTeObHJQ3nGDNs4H4O778A==

3. JWT 密钥 (JWT_SECRET_KEY) - 32字节:
   MiRSQJJKEW2euVXKpvxRzjS1C5TCFlXx4RXGUXSdWpJ

4. 超级管理员密码 (SUPER_ADMIN_PASSWORD) - 24字节:
   NNLeRYZN1IF3A/T80C7+Q6mU3xBZtdnu
```

### Linux/Mac Bash

```bash
# 添加执行权限（首次运行）
chmod +x scripts/generate-secrets.sh

# 运行密钥生成工具
./scripts/generate-secrets.sh
```

### 手动生成（使用 OpenSSL）

如果自动工具不可用，可以手动使用 OpenSSL 生成：

```bash
# 生成 32 字节加密密钥
openssl rand -base64 32

# 生成 16 字节初始化向量
openssl rand -base64 16

# 生成 32 字节 JWT 密钥
openssl rand -base64 32

# 生成 24 字节管理员密码
openssl rand -base64 24
```

---

## 环境变量配置

### 开发环境 (development)

开发环境可以使用配置文件中的默认值，但建议使用环境变量：

```bash
# Linux/Mac
export CRYPTO_SECRET_KEY="H/oguKMv23lWlivgq8snNZmTzSUp6KSHZnEEo1c0Ook="
export CRYPTO_IV="hTeObHJQ3nGDNs4H4O778A=="
export SUPER_ADMIN_PASSWORD="NNLeRYZN1IF3A/T80C7+Q6mU3xBZtdnu"

# Windows PowerShell
$env:CRYPTO_SECRET_KEY="H/oguKMv23lWlivgq8snNZmTzSUp6KSHZnEEo1c0Ook="
$env:CRYPTO_IV="hTeObHJQ3nGDNs4H4O778A=="
$env:SUPER_ADMIN_PASSWORD="NNLeRYZN1IF3A/T80C7+Q6mU3xBZtdnu"
```

### 生产环境 (production)

生产环境**必须**通过环境变量配置以下密钥：

```bash
# 加密配置（必须）
CRYPTO_ENABLED=true
CRYPTO_SECRET_KEY=<生成的32字节密钥>
CRYPTO_IV=<生成的16字节IV>

# JWT 配置（必须）
JWT_SECRET_KEY=<生成的32+字节密钥>

# 超级管理员配置（必须）
SUPER_ADMIN_EMAIL=admin@gamelink.com
SUPER_ADMIN_PASSWORD=<生成的8+字符强密码>
SUPER_ADMIN_NAME=Super Admin

# 数据库配置
POSTGRES_USER=gamelink
POSTGRES_PASSWORD=<安全密码>
POSTGRES_DB=gamelink

# Redis 配置
REDIS_PASSWORD=<安全密码>
```

#### Docker Compose 配置

创建 `.env` 文件（不要提交到 Git）：

```env
# GameLink 生产环境密钥
CRYPTO_ENABLED=true
CRYPTO_SECRET_KEY=H/oguKMv23lWlivgq8snNZmTzSUp6KSHZnEEo1c0Ook=
CRYPTO_IV=hTeObHJQ3nGDNs4H4O778A==
JWT_SECRET_KEY=MiRSQJJKEW2euVXKpvxRzjS1C5TCFlXx4RXGUXSdWpJ
SUPER_ADMIN_PASSWORD=NNLeRYZN1IF3A/T80C7+Q6mU3xBZtdnu
```

启动服务：
```bash
docker-compose up -d
```

---

## 配置验证规则

### 启动时验证

应用启动时会自动验证配置，不符合要求的将拒绝启动：

#### 加密配置验证

- ✅ 加密禁用时，不验证密钥
- ❌ 加密启用时，`CRYPTO_SECRET_KEY` 不能为空
- ❌ 加密启用时，`CRYPTO_IV` 不能为空
- ❌ 密钥长度必须是 16、24 或 32 字节
- ❌ 不能使用硬编码的默认值
- ❌ IV 长度至少 16 字节
- ❌ `CRYPTO_METHODS` 不能为空

#### JWT 配置验证

- ✅ 空密钥不报错（使用默认值）
- ❌ 密钥长度至少 16 字节
- ❌ 不能使用废弃的默认密钥

#### 超级管理员密码验证

**开发环境**:
- ✅ 任意 6+ 字符密码
- ❌ 少于 6 字符

**生产环境**:
- ❌ 密码不能为空
- ❌ 少于 8 字符
- ❌ 必须包含大写字母、小写字母、数字、特殊符号

### 测试验证

运行配置验证测试：

```bash
cd api
go test ./pkg/config/... -v
```

**测试覆盖**:
- 密钥为空检测
- 密钥长度验证
- 硬编码默认值检测
- 密码强度验证
- 生产环境强制要求

---

## 最佳实践

### 1. 密钥管理

- ✅ 使用密码管理器存储所有密钥（如 1Password、Bitwarden）
- ✅ 不同环境使用不同的密钥
- ✅ 定期轮换密钥（建议每 6 个月）
- ❌ 不要将密钥提交到 Git 仓库
- ❌ 不要在代码中硬编码密钥
- ❌ 不要在日志中输出密钥

### 2. 环境隔离

```
开发环境 (development)
  ├─ 使用配置文件默认值
  ├─ 加密可禁用 (CRYPTO_ENABLED=false)
  └─ 密码要求宽松

测试环境 (testing)
  ├─ 使用环境变量配置
  ├─ 加密启用
  └─ 使用与生产相同的安全标准

生产环境 (production)
  ├─ 强制使用环境变量
  ├─ 加密必须启用
  ├─ 强密码要求
  └─ 所有密钥必须配置
```

### 3. 密钥泄露应对

如果发现密钥泄露：

1. **立即行动**:
   - 重新生成所有密钥
   - 更新环境变量
   - 重启所有服务

2. **检查影响**:
   - 审查访问日志
   - 检查异常活动
   - 通知用户修改密码

3. **预防措施**:
   - 实施密钥轮换策略
   - 启用审计日志
   - 加强访问控制

### 4. Git 安全

确保 `.gitignore` 包含：

```gitignore
# 环境变量文件
.env
.env.local
.env.production

# 配置文件（如果包含密钥）
configs/config.production.yaml
```

检查是否已提交密钥：
```bash
# 搜索可能的密钥
git log --all --full-history --source -- "**/config.production.yaml"
git log --all --full-history --source -- ".env"
```

### 5. CI/CD 安全

- ✅ 使用 CI/CD 平台的环境变量或密钥管理服务
- ✅ 生产环境密钥不在 CI 日志中输出
- ✅ 使用加密的 secret 存储（GitHub Secrets、GitLab CI Variables）
- ❌ 不要在 CI 脚本中硬编码密钥

---

## 配置文件位置

| 环境 | 配置文件 | 优先级 |
|-----|---------|-------|
| Development | `api/configs/config.development.yaml` | 环境变量 > 配置文件 > 默认值 |
| Production | `api/configs/config.production.yaml` | 环境变量 > 配置文件 > 默认值 |

**配置加载流程**:
1. 加载 YAML 配置文件
2. 用环境变量覆盖
3. 验证配置完整性
4. 启动应用

---

## 相关文档

- [产品概述](../.kiro/steering/01-product.md) - 业务背景
- [技术栈](../.kiro/steering/02-tech-stack.md) - 环境变量说明
- [集成测试计划](INTEGRATION_TEST_PLAN.md) - 测试配置

---

## 快速参考

### 常用命令

```bash
# 生成密钥
.\scripts\generate-secrets.ps1          # Windows
./scripts/generate-secrets.sh           # Linux/Mac

# 手动生成密钥
openssl rand -base64 32                 # 32字节密钥
openssl rand -base64 16                 # 16字节IV

# 运行测试
cd api
go test ./pkg/config/... -v             # 配置验证测试
make test                               # 全部测试

# 检查配置
APP_ENV=production go run cmd/main.go   # 测试生产配置
```

### 故障排查

**问题**: `CRYPTO_SECRET_KEY is required when encryption is enabled`

**解决**:
```bash
# 检查环境变量
echo $CRYPTO_SECRET_KEY

# 重新生成密钥
openssl rand -base64 32

# 设置环境变量
export CRYPTO_SECRET_KEY=<生成的密钥>
```

**问题**: `SUPER_ADMIN_PASSWORD does not meet security requirements`

**解决**:
```bash
# 生成强密码（24字节）
openssl rand -base64 24

# 或手动创建（8+字符，包含大小写、数字、特殊符号）
# 例如: SecurePass123!@#
```

---

**文档版本**: 1.0.0
**最后更新**: 2025-01-01
**维护者**: GameLink Security Team
