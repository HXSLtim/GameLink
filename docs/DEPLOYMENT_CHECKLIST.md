# GameLink 生产环境部署检查清单

**任务ID：** #36
**负责人：** DevOps-Engineer
**优先级：** P0
**预估工作量：** 5-7 天

---

## 部署前检查清单

### 1. 环境变量配置 ⏳

#### 1.1 数据库配置
- [ ] PostgreSQL 强密码生成
- [ ] 数据库连接字符串配置
- [ ] 连接池参数设置
- [ ] 主从复制配置（可选）

#### 1.2 缓存配置
- [ ] Redis 强密码生成
- [ ] Redis 连接配置
- [ ] 内存策略配置
- [ ] 持久化配置

#### 1.3 加密配置
- [ ] CRYPTO_SECRET_KEY 生成（32字节）
- [ ] CRYPTO_IV 生成（16字节）
- [ ] CRYPTO_USE_SIGNATURE 启用
- [ ] 前后端密钥一致性验证

#### 1.4 JWT 配置
- [ ] JWT_SECRET_KEY 生成（32+字符）
- [ ] JWT_TOKEN_TTL_HOURS 设置

#### 1.5 超级管理员配置
- [ ] SUPER_ADMIN_EMAIL 设置
- [ ] SUPER_ADMIN_PASSWORD 强密码生成
- [ ] SUPER_ADMIN_NAME 设置

#### 1.6 外部服务配置（可选）
- [ ] 微信支付参数配置
- [ ] 支付宝参数配置
- [ ] 短信服务参数配置
- [ ] 对象存储参数配置

---

### 2. 安全加固 ⏳

#### 2.1 加密启用
- [ ] 后端 CRYPTO_ENABLED=true
- [ ] 前端 VITE_CRYPTO_ENABLED=true
- [ ] 密钥配置正确性验证

#### 2.2 SSL/TLS 配置
- [ ] SSL 证书准备
- [ ] Nginx SSL 配置
- [ ] 强制 HTTPS 重定向
- [ ] HSTS 头配置

#### 2.3 防火墙规则
- [ ] 仅开放必要端口（80, 443, 22）
- [ ] 数据库端口仅内部访问
- [ ] Redis 端口仅内部访问
- [ ] 限制 SSH 访问 IP

#### 2.4 Redis 安全
- [ ] requirepass 配置
- [ ] 禁用危险命令（FLUSHDB, FLUSHALL）
- [ ] 绑定到本地回环地址

#### 2.5 容器安全
- [ ] 非 root 用户运行
- [ ] 资源限制配置
- [ ] 只读根文件系统（可选）
- [ ] 安全扫描通过

---

### 3. 部署配置 ⏳

#### 3.1 Docker Compose 配置
- [ ] docker-compose.prod.yml 更新
- [ ] 网络配置正确
- [ ] 卷挂载配置
- [ ] 健康检查配置

#### 3.2 Nginx 配置
- [ ] 反向代理配置
- [ ] 负载均衡配置（可选）
- [ ] 静态资源缓存配置
- [ ] WebSocket 代理配置
- [ ] Gzip 压缩启用

#### 3.3 日志配置
- [ ] 日志级别设置
- [ ] 日志轮转配置
- [ ] 日志收集配置（可选）
- [ ] 审计日志启用

#### 3.4 监控配置
- [ ] Prometheus 集成
- [ ] Grafana 仪表板
- [ ] 告警规则配置
- [ ] 健康检查端点

---

### 4. 部署测试 ⏳

#### 4.1 Staging 环境部署
- [ ] 基础设施启动
- [ ] 数据库迁移执行
- [ ] 后端服务部署
- [ ] 前端应用部署
- [ ] 健康检查通过

#### 4.2 功能验证测试
- [ ] 用户注册/登录
- [ ] 订单创建流程
- [ ] 支付流程测试
- [ ] WebSocket 连接测试
- [ ] 文件上传测试

#### 4.3 性能基准测试
- [ ] API 响应时间测试
- [ ] 并发用户测试
- [ ] 数据库查询性能
- [ ] 内存和 CPU 使用率

#### 4.4 回滚演练
- [ ] 备份当前版本
- [ ] 部署新版本
- [ ] 验证新版本
- [ ] 执行回滚
- [ ] 验证回滚成功

---

### 5. 文档更新 ⏳

#### 5.1 配置文档
- [ ] DEPENDENCIES_AND_CONFIG.md 更新
- [ ] 环境变量说明更新
- [ ] 生产环境配置说明

#### 5.2 部署文档
- [ ] 部署操作手册创建
- [ ] 部署步骤说明
- [ ] 验证步骤说明
- [ ] 常见问题处理

#### 5.3 应急文档
- [ ] 应急处理手册创建
- [ ] 回滚操作步骤
- [ ] 故障排查流程
- [ ] 紧急联系方式

---

## 部署步骤

### 阶段 1：准备阶段（Day 1-2）

1. **生成密钥和密码**
   ```bash
   # 生成数据库密码
   openssl rand -base64 24

   # 生成 Redis 密码
   openssl rand -base64 24

   # 生成 JWT 密钥
   openssl rand -base64 32

   # 生成加密密钥（后端 base64）
   openssl rand -base64 32

   # 生成加密密钥（前端原始字节）
   openssl rand -base64 32 | base64 -d | xxd -p -c 32

   # 生成 IV（后端 base64）
   openssl rand -base64 16

   # 生成 IV（前端原始字节）
   openssl rand -base64 16 | base64 -d | xxd -p -c 16
   ```

2. **创建生产环境配置文件**
   ```bash
   cp .env.example .env.production
   # 编辑 .env.production，填入生成的密钥
   ```

3. **验证配置**
   ```bash
   # 检查配置完整性
   ./scripts/validate-config.sh
   ```

### 阶段 2：配置阶段（Day 2-3）

1. **更新 Docker Compose 配置**
2. **配置 Nginx 反向代理**
3. **配置 SSL/TLS**
4. **配置防火墙规则**

### 阶段 3：部署阶段（Day 3-4）

1. **部署到 Staging 环境**
   ```bash
   docker-compose -f docker-compose.staging.yml --env-file .env.staging up -d
   ```

2. **执行数据库迁移**
   ```bash
   docker exec gamelink-backend /app/server migrate
   ```

3. **验证部署**
   ```bash
   # 健康检查
   curl https://staging.gamelink.com/api/v1/healthz

   # 前端访问
   curl https://staging.gamelink.com/
   ```

### 阶段 4：测试阶段（Day 4-5）

1. **功能测试**
2. **性能测试**
3. **安全测试**
4. **回滚演练**

### 阶段 5：上线阶段（Day 5-6）

1. **生产环境部署**
   ```bash
   docker-compose -f docker-compose.prod.yml --env-file .env.production up -d
   ```

2. **监控验证**
3. **性能监控**
4. **日志检查**

### 阶段 6：文档阶段（Day 6-7）

1. **更新部署文档**
2. **创建运维手册**
3. **培训运维团队**

---

## 验收标准

- [ ] **环境变量**：所有必需的配置项已正确设置
- [ ] **Staging 部署**：Staging 环境成功部署并验证
- [ ] **安全加固**：所有安全措施已实施
- [ ] **监控告警**：监控和告警系统正常运行
- [ ] **文档更新**：所有文档已更新并审核

---

## 协作成员

**主导：** DevOps-Engineer
**协作：**
- Backend-Lead（后端配置验证）
- Database-Architect（数据库配置）
- Frontend-Lead（前端配置）
- Mobile-Lead（移动端配置）

---

## 进度跟踪

**当前状态：** 进行中
**开始日期：** 2026-02-09
**预计完成：** 2026-02-16

---

**最后更新：** 2026-02-09
**更新人：** DevOps-Engineer
