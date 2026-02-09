# GameLink 监控仪表板配置指南

**版本：** 1.0.0
**创建日期：** 2026-02-09
**维护人：** DevOps-Engineer

---

## 1. 监控架构概述

### 推荐方案：Prometheus + Grafana

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Services  │────▶│ Prometheus  │────▶│   Grafana   │
│             │     │             │     │             │
│  - API      │     │  - Metrics  │     │ - Dashboards│
│  - Admin    │     │  - Storage  │     │ - Alerts    │
│  - App      │     │  - Alerts   │     │ - Reports   │
│  - Postgres │     │             │     │             │
│  - Redis    │     │             │     │             │
└─────────────┘     └─────────────┘     └─────────────┘
```

### 备选方案（轻量级）

1. **cAdvisor + Node Exporter** - 容器和主机监控
2. **uptime-kuma** - 简单的状态页面
3. **Netdata** - 一体化监控解决方案

---

## 2. 快速部署（Docker Compose）

### 2.1 创建监控配置文件

**文件：** `docker-compose.monitoring.yml`

```yaml
version: '3.8'

services:
  # Prometheus - 指标收集和存储
  prometheus:
    image: prom/prometheus:latest
    container_name: gamelink-prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./monitoring/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus-data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/usr/share/prometheus/console_libraries'
      - '--web.console.templates=/usr/share/prometheus/consoles'
    networks:
      - monitoring
    restart: unless-stopped

  # Grafana - 可视化仪表板
  grafana:
    image: grafana/grafana:latest
    container_name: gamelink-grafana
    ports:
      - "3001:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin123
      - GF_USERS_ALLOW_SIGN_UP=false
    volumes:
      - grafana-data:/var/lib/grafana
      - ./monitoring/grafana/provisioning:/etc/grafana/provisioning
      - ./monitoring/grafana/dashboards:/var/lib/grafana/dashboards
    networks:
      - monitoring
    restart: unless-stopped
    depends_on:
      - prometheus

  # cAdvisor - 容器监控
  cadvisor:
    image: gcr.io/cadvisor/cadvisor:latest
    container_name: gamelink-cadvisor
    ports:
      - "8081:8080"
    volumes:
      - /:/rootfs:ro
      - /var/run:/var/run:ro
      - /sys:/sys:ro
      - /var/lib/docker/:/var/lib/docker:ro
      - /dev/disk/:/dev/disk:ro
    networks:
      - monitoring
    restart: unless-stopped
    privileged: true

  # Node Exporter - 主机监控
  node-exporter:
    image: prom/node-exporter:latest
    container_name: gamelink-node-exporter
    ports:
      - "9100:9100"
    volumes:
      - /proc:/host/proc:ro
      - /sys:/host/sys:ro
      - /:/rootfs:ro
    command:
      - '--path.procfs=/host/proc'
      - '--path.sysfs=/host/sys'
      - '--collector.filesystem.mount-points-exclude=^/(sys|proc|dev|host|etc)($$|/)'
    networks:
      - monitoring
    restart: unless-stopped

networks:
  monitoring:
    driver: bridge

volumes:
  prometheus-data:
  grafana-data:
```

### 2.2 Prometheus 配置

**文件：** `monitoring/prometheus/prometheus.yml`

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s
  external_labels:
    cluster: 'gamelink'
    environment: 'development'

# 告警规则文件
rule_files:
  - 'alerts/*.yml'

# 抓取配置
scrape_configs:
  # Prometheus 自身监控
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  # 容器监控 (cAdvisor)
  - job_name: 'cadvisor'
    static_configs:
      - targets: ['cadvisor:8080']

  # 主机监控 (Node Exporter)
  - job_name: 'node'
    static_configs:
      - targets: ['node-exporter:9100']

  # GameLink API
  - job_name: 'gamelink-api'
    static_configs:
      - targets: ['host.docker.internal:8080']
    metrics_path: '/metrics'
    scrape_interval: 10s

  # PostgreSQL Exporter
  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres-exporter:9187']

  # Redis Exporter
  - job_name: 'redis'
    static_configs:
      - targets: ['redis-exporter:9121']

# 告警管理器配置
alerting:
  alertmanagers:
    - static_configs:
        - targets: ['alertmanager:9093']
```

### 2.3 告警规则配置

**文件：** `monitoring/prometheus/alerts/gamelink.yml`

```yaml
groups:
  - name: gamelink_alerts
    interval: 30s
    rules:
      # 服务可用性告警
      - alert: ServiceDown
        expr: up == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "服务 {{ $labels.instance }} 宕机"
          description: "{{ $labels.job }} 服务已经宕机超过 1 分钟"

      # 高 CPU 使用率告警
      - alert: HighCPUUsage
        expr: rate(process_cpu_seconds_total[5m]) * 100 > 80
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "CPU 使用率过高"
          description: "{{ $labels.instance }} CPU 使用率超过 80%"

      # 高内存使用率告警
      - alert: HighMemoryUsage
        expr: (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100 > 85
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "内存使用率过高"
          description: "内存使用率超过 85%"

      # 磁盘空间告警
      - alert: LowDiskSpace
        expr: (node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes{mountpoint="/"}) * 100 < 15
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "磁盘空间不足"
          description: "根分区可用空间低于 15%"

      # API 响应时间告警
      - alert: SlowAPIResponse
        expr: histogram_quantile(0.95, http_request_duration_seconds_bucket) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "API 响应时间过长"
          description: "95% 的 API 请求响应时间超过 1 秒"

      # API 错误率告警
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "API 错误率过高"
          description: "API 错误率超过 5%"
```

---

## 3. Grafana 仪表板配置

### 3.1 数据源配置

**文件：** `monitoring/grafana/provisioning/datasources/prometheus.yml`

```yaml
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
    editable: true
```

### 3.2 仪表板配置

**文件：** `monitoring/grafana/provisioning/dashboards/dashboard.yml`

```yaml
apiVersion: 1

providers:
  - name: 'Default'
    orgId: 1
    folder: ''
    type: file
    disableDeletion: false
    updateIntervalSeconds: 10
    allowUiUpdates: true
    options:
      path: /var/lib/grafana/dashboards
```

### 3.3 关键指标仪表板

创建以下仪表板 JSON 文件：

1. **系统概览** (`System Overview.json`)
   - CPU、内存、磁盘、网络使用率
   - 容器状态
   - 服务健康状态

2. **API 性能** (`API Performance.json`)
   - 请求速率
   - 响应时间 (P50, P95, P99)
   - 错误率
   - 端点性能排行

3. **数据库监控** (`Database Monitoring.json`)
   - 连接数
   - 查询性能
   - 缓存命中率
   - 复制延迟

---

## 4. 部署步骤

### 4.1 启动监控栈

```bash
# 创建必要的目录
mkdir -p monitoring/{prometheus,grafana/{provisioning,dashboards}}

# 复制配置文件
# (将上述配置文件复制到相应目录)

# 启动监控服务
docker-compose -f docker-compose.monitoring.yml up -d

# 检查服务状态
docker-compose -f docker-compose.monitoring.yml ps
```

### 4.2 访问监控界面

- **Prometheus:** http://localhost:9090
- **Grafana:** http://localhost:3001
  - 用户名: `admin`
  - 密码: `admin123`
- **cAdvisor:** http://localhost:8081

### 4.3 验证监控数据

1. 打开 Prometheus UI
2. 检查 Targets 页面，确认所有目标都在抓取
3. 执行查询验证数据：`up`
4. 打开 Grafana，导入预配置的仪表板

---

## 5. 关键监控指标

### 5.1 系统级指标

| 指标 | 类型 | 预警阈值 | 危险阈值 |
|------|------|---------|---------|
| CPU 使用率 | Gauge | > 50% | > 80% |
| 内存使用率 | Gauge | > 70% | > 90% |
| 磁盘使用率 | Gauge | > 70% | > 85% |
| 网络 I/O | Counter | - | - |

### 5.2 容器级指标

| 指标 | 说明 |
|------|------|
| `container_cpu_usage_seconds_total` | CPU 使用时间 |
| `container_memory_usage_bytes` | 内存使用量 |
| `container_network_*` | 网络 I/O |
| `container_fs_*` | 文件系统使用 |

### 5.3 应用级指标

| 指标 | 类型 | 预警阈值 |
|------|------|---------|
| `http_requests_total` | Counter | - |
| `http_request_duration_seconds` | Histogram | P95 > 1s |
| `http_requests_total{status="500"}` | Counter | 错误率 > 5% |
| `api_response_time` | Gauge | > 500ms |

---

## 6. 告警通知配置

### 6.1 Email 通知

**修改 Prometheus 配置：**

```yaml
alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - alertmanager:9093

# Alertmanager 配置文件 (alertmanager.yml)
global:
  resolve_timeout: 5m
  smtp_smarthost: 'smtp.gmail.com:587'
  smtp_from: 'alertmanager@gamelink.com'
  smtp_auth_username: 'your-email@gmail.com'
  smtp_auth_password: 'your-password'

route:
  group_by: ['alertname', 'cluster', 'service']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 12h
  receiver: 'default'
  routes:
    - match:
        severity: critical
      receiver: 'critical-alerts'

receivers:
  - name: 'default'
    email_configs:
      - to: 'team@gamelink.com'
        headers:
          Subject: '[Alert] {{ .GroupLabels.alertname }}'

  - name: 'critical-alerts'
    email_configs:
      - to: 'oncall@gamelink.com'
        headers:
          Subject: '[CRITICAL] {{ .GroupLabels.alertname }}'
```

### 6.2 钉钉通知

使用 Webhook 配置钉钉机器人：

```yaml
receivers:
  - name: 'dingtalk'
    webhook_configs:
      - url: 'https://oapi.dingtalk.com/robot/send?access_token=TOKEN'
        send_resolved: true
```

---

## 7. 维护和优化

### 7.1 数据保留策略

**Prometheus 配置：**

```yaml
# prometheus.yml
global:
  scrape_interval: 15s

# 数据保留时间
storage:
  tsdb:
    retention.time: 30d
    retention.size: 10GB
```

### 7.2 性能优化

1. **降低抓取频率** - 对非关键指标使用 30s 或 60s
2. **使用 recording rules** - 预计算复杂查询
3. **分区数据** - 按环境或服务分区
4. **清理旧数据** - 定期删除不需要的指标

### 7.3 备份和恢复

```bash
# 备份 Prometheus 数据
docker exec gamelink-prometheus tar czf /backup/prometheus-$(date +%Y%m%d).tar.gz /prometheus

# 备份 Grafana 数据
docker exec gamelink-grafana tar czf /backup/grafana-$(date +%Y%m%d).tar.gz /var/lib/grafana

# 恢复
docker exec gamelink-prometheus tar xzf /backup/prometheus-20260209.tar.gz -C /
```

---

## 8. 故障排查

### 8.1 Prometheus 无法启动

```bash
# 检查配置文件
docker exec gamelink-prometheus promtool check config /etc/prometheus/prometheus.yml

# 查看日志
docker logs gamelink-prometheus

# 重新启动
docker-compose -f docker-compose.monitoring.yml restart prometheus
```

### 8.2 Grafana 无法连接 Prometheus

1. 确认 Prometheus 正在运行
2. 检查数据源配置
3. 验证网络连接
4. 查看 Grafana 日志

### 8.3 指标数据缺失

1. 检查 Exporter 是否运行
2. 验证抓取目标状态
3. 查看 Prometheus 日志
4. 确认网络连接

---

## 9. 下一步

- [ ] 部署监控栈
- [ ] 配置告警规则
- [ ] 设置告警通知
- [ ] 创建自定义仪表板
- [ ] 测试告警流程
- [ ] 编写运维文档

---

**相关文档：**
- **服务状态报告：** `docs/SERVICE_STATUS_REPORT.md`
- **监控告警指南：** `docs/MONITORING_ALERT_GUIDE.md`
- **健康检查脚本：** `scripts/health-check.sh` / `scripts/health-check.bat`
