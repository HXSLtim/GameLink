# GameLink 支付回调 URL 配置说明

**任务：** #34 支付渠道集成
**负责人：** DevOps-Engineer

---

## 回调 URL 配置规范

### 生产环境回调地址

**主域名配置：**
- 主域名：`gamelink.com` 或 `www.gamelink.com`
- API 域名：`api.gamelink.com`（推荐）或使用主域名路径

**微信支付回调地址：**
```
https://api.gamelink.com/api/v1/payments/wechat/notify
```

**支付宝回调地址：**
```
https://api.gamelink.com/api/v1/payments/alipay/notify
```

**Staging 环境回调地址：**
```
https://staging.gamelink.com/api/v1/payments/wechat/notify
https://staging.gamelink.com/api/v1/payments/alipay/notify
```

### 测试环境配置

**注意事项：**
1. 微信支付和支付宝要求回调地址必须是 HTTPS
2. 测试环境可以使用内网穿透工具（如 ngrok、frp）
3. 生产环境必须使用正式域名和 SSL 证书

---

## 环境变量配置

### .env.production 配置模板

```bash
# =============================================================================
# 微信支付配置
# =============================================================================
WECHAT_PAY_ENABLED=true

# 微信公众号/小程序 AppID
WECHAT_PAY_APP_ID=wx####################

# 微信商户号
WECHAT_PAY_MCH_ID=1####################

# API 密钥（在微信商户平台设置）
WECHAT_PAY_API_KEY=######################

# 商户证书路径（Docker 容器内路径）
WECHAT_PAY_CERT_PATH=/app/certs/apiclient_cert.p12
WECHAT_PAY_KEY_PATH=/app/certs/apiclient_key.pem

# 证书序列号
WECHAT_PAY_CERT_SERIAL_NO=######################

# 回调地址
WECHAT_PAY_NOTIFY_URL=https://api.gamelink.com/api/v1/payments/wechat/notify

# =============================================================================
# 支付宝配置
# =============================================================================
ALIPAY_ENABLED=true

# 支付宝应用 ID
ALIPAY_APP_ID=####################

# 支付宝网关地址
ALIPAY_GATEWAY=https://openapi.alipay.com/gateway.do

# 回调地址
ALIPAY_NOTIFY_URL=https://api.gamelink.com/api/v1/payments/alipay/notify

# =============================================================================
# 证书文件存储配置
# =============================================================================
# 证书文件将挂载到容器内的 /app/certs 目录
```

### .env.staging 配置模板

```bash
# Staging 环境配置

# 微信支付（测试环境）
WECHAT_PAY_ENABLED=false
# 测试环境暂不启用真实支付，使用 Mock 模式

# 支付宝（测试环境）
ALIPAY_ENABLED=false
# 测试环境暂不启用真实支付，使用 Mock 模式
```

---

## 证书文件存储

### Docker 卷挂载配置

**docker-compose.prod.yml 配置：**

```yaml
services:
  backend:
    volumes:
      # 挂载证书目录
      - ./certs/wechat:/app/certs/wechat:ro
      - ./certs/alipay:/app/certs/alipay:ro
```

### 目录结构

```
certs/
├── wechat/
│   ├── apiclient_cert.p12      # 微信商户证书（PKCS12格式）
│   ├── apiclient_key.pem       # 微信商户私钥
│   └── rootca.pem              # 微信支付根证书
└── alipay/
    ├── app_private_key.pem      # 支付宝应用私钥
    ├── app_public_key.pem       # 支付宝应用公钥
    └── alipay_rootCert.pem      # 支付宝根证书
```

### 安全注意事项

**证书文件权限：**
```bash
# 设置证书文件权限为只读
chmod 400 certs/wechat/apiclient_cert.p12
chmod 400 certs/alipay/app_private_key.pem

# 确保密钥文件不被意外提交到版本控制
echo "certs/" >> .gitignore
```

---

## 域名和 HTTPS 配置

### 生产环境域名要求

**域名备案：**
- ✅ 主域名：已备案（例如 `gamelink.com`）
- ✅ API 子域名：已备案（推荐 `api.gamelink.com`）
- ✅ 管理后台子域名：已备案（例如 `admin.gamelink.com`）

**域名规划：**
```
gamelink.com              # 主站（前端）
www.gamelink.com          # 主站（前端）
api.gamelink.com          # API 服务
admin.gamelink.com        # 管理后台
staging.gamelink.com      # 测试环境
```

### SSL/TLS 证书配置

**证书类型：**
- Let's Encrypt（免费，自动续期）
- 商业证书（付费，更高级别）

**获取 Let's Encrypt 证书：**
```bash
# 安装 certbot
sudo apt install certbot

# 获取证书（单域名）
sudo certbot certonly --standalone \
  -d api.gamelink.com \
  -d gamelink.com \
  -d www.gamelink.com \
  -d admin.gamelink.com

# 证书文件位置
/etc/letsencrypt/live/api.gamelink.com/fullchain.pem
/etc/letsencrypt/live/api.gamelink.com/privkey.pem
```

**Nginx SSL 配置（已包含在 nginx-production.conf）：**
```nginx
ssl_certificate /etc/nginx/ssl/fullchain.pem;
ssl_certificate_key /etc/nginx/ssl/privkey.pem;
ssl_protocols TLSv1.2 TLSv1.3;
```

---

## Nginx 反向代理配置

### 支付回调路由配置

**已在 `deploy/nginx-production.conf` 中包含：**

```nginx
# API 代理（包含支付回调）
location /api/ {
    proxy_pass http://backend;
    proxy_http_version 1.1;

    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;

    # 超时设置（支付回调可能需要更长时间）
    proxy_connect_timeout 60s;
    proxy_send_timeout 60s;
    proxy_read_timeout 60s;
}
```

### 支付平台 IP 白名单

**微信支付服务器 IP：**
- 最新 IP 段：在微信商户平台查询
- 需要在防火墙中白名单这些 IP

**支付宝服务器 IP：**
- 最新 IP 段：在支付宝开放平台查询
- 需要在防火墙中白名单这些 IP

---

## 测试环境配置

### 内网穿透方案

**使用 ngrok（临时测试）：**

```bash
# 安装 ngrok
# 下载：https://ngrok.com/download

# 启动 ngrok（映射本地 8080 端口）
ngrok http 8080

# 获得 HTTPS 回调地址
# 例如：https://abc123.ngrok.io/api/v1/payments/wechat/notify
```

**使用 frp（推荐）：**

```ini
# frps.ini（服务器端）
[common]
bind_port = 7000
vhost_http_port = 8080

[https]
listen_port = 443
tls_cert_file = /path/to/cert.pem
tls_key_file = /path/to/key.pem
```

```ini
# frpc.ini（客户端）
[common]
server_addr = your-server.com:7000

[web]
type = http
local_ip = 127.0.0.1
local_port = 8080
custom_domains = staging.gamelink.com
```

---

## 支付平台配置要求

### 微信支付要求

**必须满足：**
1. ✅ 回调地址必须是 HTTPS
2. ✅ 域名必须备案（中国大陆）
3. ✅ SSL/TLS 证书必须有效
4. ✅ 回调地址必须在商户平台配置

**建议配置：**
- 回调超时时间：5 秒
- 重试次数：3 次
- 重试间隔：1 秒

### 支付宝要求

**必须满足：**
1. ✅ 回调地址必须是 HTTPS（生产环境）
2. ✅ 应用公钥必须上传到支付宝平台
3. ✅ 支付宝公钥必须下载并配置

**建议配置：**
- 回调超时时间：30 秒
- 同步跳转地址：前端支付结果页面
- 异步通知地址：后端回调接口

---

## 回调接口实现

### 后端接口路由（Backend-Lead 负责）

**微信支付回调接口：**
- 路径：`/api/v1/payments/wechat/notify`
- 方法：POST
- Content-Type：application/xml

**支付宝回调接口：**
- 路径：`/api/v1/payments/alipay/notify`
- 方法：POST
- Content-Type：application/x-www-form-urlencoded

### 验证步骤

**微信支付：**
1. 接收回调数据
2. 验证签名（使用商户证书）
3. 验证订单金额
4. 更新订单状态
5. 返回成功响应（XML 格式）

**支付宝：**
1. 接收回调数据
2. 验证签名（使用支付宝公钥）
3. 验证订单金额
4. 更新订单状态
5. 返回成功响应（文本 "success"）

---

## 当前配置状态

### ✅ 已配置

- [x] Nginx 反向代理配置（`deploy/nginx-production.conf`）
- [x] SSL/TLS 配置模板
- [x] 回调 URL 路由配置
- [x] 环境变量配置模板

### ⏳ 待配置

- [ ] 实际域名注册和备案
- [ ] SSL 证书申请和配置
- [ ] 微信支付商户号申请
- [ ] 支付宝应用申请
- [ ] 证书文件生成和存储
- [ ] Docker 卷挂载配置

---

## 下一步行动

### 立即需要 Backend-Lead 确认

1. **回调接口路由：**
   - 确认路由路径是否符合上述规范
   - 确认请求格式（XML/Form-Data）
   - 确认响应格式

2. **环境变量：**
   - 确认所有必需的环境变量
   - 确认证书文件路径
   - 确认回调 URL 格式

3. **证书准备：**
   - 准备微信商户证书（P12 格式）
   - 准备支付宝 RSA 密钥对
   - 生成证书存储目录

### 需要提供的实际信息

1. **生产环境域名：**
   - 主域名：`_____________`
   - API 域名：`_____________`
   - 是否已备案：`[ ]`

2. **SSL 证书：**
   - 证书类型：`[ ] Let's Encrypt / [ ] 商业证书`
   - 证书状态：`[ ] 已申请 / [ ] 待申请`

3. **支付商户号：**
   - 微信商户号：`_____________`
   - 支付宝应用ID：`_____________`

---

**最后更新：** 2026-02-09
**更新人：** DevOps-Engineer
**文档状态：** 待 Backend-Lead 确认
