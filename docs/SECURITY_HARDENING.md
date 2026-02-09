# GameLink 安全加固指南

**任务ID：** #36
**负责人：** DevOps-Engineer
**优先级：** P0

---

## 1. 基础设施安全

### 1.1 网络安全

**防火墙配置：**
```bash
# 仅开放必要端口
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable

# 限制 SSH 访问（可选）
sudo ufw allow from 203.0.113.0/24 to any port 22
```

**禁用 ping：**
```bash
# 临时禁用
sudo echo "1" > /proc/sys/net/ipv4/icmp_echo_ignore_all

# 永久禁用（编辑 /etc/sysctl.conf）
net.ipv4.icmp_echo_ignore_all = 1
```

### 1.2 SSH 安全

**配置文件：** `/etc/ssh/sshd_config`

```bash
# 禁用 root 登录
PermitRootLogin no

# 禁用密码登录，仅允许密钥
PasswordAuthentication no
PubkeyAuthentication yes

# 限制登录用户
AllowUsers gamelink

# 禁用空密码
PermitEmptyPasswords no

# 更改默认端口（可选）
Port 2222
```

**重启 SSH：**
```bash
sudo systemctl restart sshd
```

### 1.3 系统加固

**自动更新：**
```bash
# Ubuntu/Debian
sudo apt install unattended-upgrades
sudo dpkg-reconfigure -plow unattended-upgrades

# CentOS/RHEL
sudo yum install yum-cron
sudo systemctl enable yum-cron
sudo systemctl start yum-cron
```

**失败登录审计：**
```bash
# 安装 fail2ban
sudo apt install fail2ban

# 启动服务
sudo systemctl enable fail2ban
sudo systemctl start fail2ban
```

---

## 2. 容器安全

### 2.1 Docker 安全

**非 root 用户运行：**
```dockerfile
# 已在 Dockerfile 中配置
USER gamelink
```

**只读根文件系统：**
```yaml
# docker-compose.yml
services:
  backend:
    read_only: true
    tmpfs:
      - /tmp
      - /app/logs
      - /app/var
```

**资源限制：**
```yaml
deploy:
  resources:
    limits:
      cpus: '1'
      memory: 512M
    reservations:
      memory: 256M
```

**安全扫描：**
```bash
# 使用 Trivy 扫描镜像
trivy image gamelink-backend:latest

# 扫描并显示漏洞
trivy image --severity HIGH,CRITICAL gamelink-backend:latest
```

### 2.2 网络隔离

**创建独立网络：**
```yaml
networks:
  gamelink-network:
    driver: bridge
    internal: false  # 设为 true 可完全隔离
```

**容器间通信限制：**
```yaml
services:
  backend:
    networks:
      - gamelink-network
    # 不连接默认网络
```

---

## 3. 应用安全

### 3.1 加密配置

**启用加密（必须）：**
```bash
# .env.production
CRYPTO_ENABLED=true
CRYPTO_SECRET_KEY=<32字节密钥>
CRYPTO_IV=<16字节IV>
CRYPTO_USE_SIGNATURE=true
```

**验证加密一致性：**
```bash
bash scripts/verify-crypto-keys.sh
```

### 3.2 JWT 安全

**强密钥：**
```bash
# 生成 32+ 字符密钥
JWT_SECRET_KEY=$(openssl rand -base64 32)
```

**合理过期时间：**
```bash
JWT_TOKEN_TTL_HOURS=24  # 24小时
```

**刷新机制：**
- 实现 refresh token
- 短期 access token（15-30分钟）
- 长期 refresh token（7-30天）

### 3.3 CORS 配置

**严格 CORS 策略：**
```go
// 仅允许特定源
AllowOrigins: []string{
    "https://gamelink.com",
    "https://admin.gamelink.com",
}

// 仅允许必要的方法
AllowMethods: []string{"GET", "POST", "PUT", "DELETE"}

// 仅允许必要的头
AllowHeaders: []string{
    "Origin",
    "Content-Type",
    "Authorization",
    "X-Signature",
}
```

---

## 4. 数据库安全

### 4.1 PostgreSQL 安全

**强密码：**
```bash
# 生成强密码
POSTGRES_PASSWORD=$(openssl rand -base64 24)
```

**本地连接加密：**
```bash
# pg_hba.conf
host    all             all             127.0.0.1/32            scram-sha-256
host    all             all             172.18.0.0/16           scram-sha-256
```

**SSL 连接（可选）：**
```bash
# postgresql.conf
ssl = on
ssl_cert_file = '/var/lib/postgresql/server.crt'
ssl_key_file = '/var/lib/postgresql/server.key'
```

### 4.2 备份加密

**加密备份：**
```bash
# 备份并加密
pg_dump -U gamelink -d gamelink | gzip | \
  openssl enc -aes-256-cbc -salt -out backup_$(date +%Y%m%d).sql.gz.enc
```

**安全存储：**
```bash
# 使用密钥管理服务
# - AWS Secrets Manager
# - Azure Key Vault
# - HashiCorp Vault
```

---

## 5. Redis 安全

### 5.1 访问控制

**设置密码：**
```bash
# redis.conf
requirepass <强密码>
```

**禁用危险命令：**
```bash
# redis.conf
rename-command FLUSHDB ""
rename-command FLUSHALL ""
rename-command CONFIG ""
rename-command DEBUG ""
```

**绑定到本地：**
```bash
# redis.conf
bind 127.0.0.1 ::1
```

### 5.2 持久化安全

**AOF 持久化：**
```bash
appendonly yes
appendfsync everysec
```

**备份 RDB 文件：**
```bash
# 定期备份 dump.rdb
cp /var/lib/redis/dump.rdb /backups/redis_$(date +%Y%m%d).rdb
```

---

## 6. Web 服务器安全

### 6.1 Nginx 安全

**安全头配置：**
```nginx
# HSTS
add_header Strict-Transport-Security "max-age=31536000; includeSubDomains; preload" always;

# 防止点击劫持
add_header X-Frame-Options "SAMEORIGIN" always;

# 防止 MIME 类型嗅探
add_header X-Content-Type-Options "nosniff" always;

# XSS 保护
add_header X-XSS-Protection "1; mode=block" always;

# Referrer 策略
add_header Referrer-Policy "no-referrer-when-downgrade" always;
```

**限制请求大小：**
```nginx
client_max_body_size 10M;
```

**限速：**
```nginx
limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
limit_req zone=api burst=20;
```

### 6.2 SSL/TLS 配置

**Let's Encrypt 证书：**
```bash
# 安装 certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d gamelink.com -d www.gamelink.com

# 自动续期
sudo certbot renew --dry-run
```

**现代 SSL 配置：**
```nginx
ssl_protocols TLSv1.2 TLSv1.3;
ssl_ciphers 'ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256';
ssl_prefer_server_ciphers off;
```

---

## 7. 监控和日志

### 7.1 日志管理

**集中日志：**
```bash
# 使用 ELK Stack 或 Loki
# 1. Filebeat 收集日志
# 2. Logstash 处理日志
# 3. Elasticsearch 存储日志
# 4. Kibana 可视化日志
```

**日志轮转：**
```bash
# /etc/logrotate.d/gamelink
/var/log/gamelink/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 0640 gamelink gamelink
}
```

### 7.2 监控

**Prometheus + Grafana：**
```yaml
# docker-compose.yml
services:
  prometheus:
    image: prom/prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    ports:
      - "9090:9090"

  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
```

**告警规则：**
```yaml
# alerts.yml
groups:
  - name: gamelink
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        annotations:
          summary: "高错误率"
```

---

## 8. 安全检查清单

### 部署前检查

- [ ] 所有密码已更新为强密码
- [ ] SSH 仅允许密钥登录
- [ ] 防火墙已配置
- [ ] SSL/TLS 证书已配置
- [ ] 加密已启用（CRYPTO_ENABLED=true）
- [ ] JWT 密钥已设置
- [ ] 数据库备份已配置
- [ ] 日志收集已配置
- [ ] 监控已配置
- [ ] 容器资源限制已设置

### 运行时检查

- [ ] Fail2ban 已启用
- [ ] 自动更新已配置
- [ ] 备份任务已安排
- [ ] 告警规则已测试
- [ ] 灾难恢复已演练

---

## 9. 应急响应

### 数据泄露事件

**立即行动：**
1. 隔离受影响系统
2. 更换所有密钥和密码
3. 检查日志确定泄露范围
4. 通知受影响用户
5. 向监管机构报告（如需要）

### DDoS 攻击

**缓解措施：**
1. 启用 Cloudflare 或类似服务
2. 配置 Nginx 限速
3. 启用缓存
4. 扩容服务器
5. 联系 ISP

---

## 10. 合规性

### GDPR 合规

- [ ] 数据加密（传输和存储）
- [ ] 用户访问/删除权
- [ ] 数据处理协议
- [ ] 数据 breach 通知

### PCI DSS 合规（如处理支付）

- [ ] 网络隔离
- [ ] 数据加密
- [ ] 访问控制
- [ ] 定期审计
- [ ] 漏洞扫描

---

**最后更新：** 2026-02-09
**更新人：** DevOps-Engineer
