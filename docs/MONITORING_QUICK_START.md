# 监控系统快速访问指南

**更新时间：** 2026-02-09
**维护者：** DevOps-Engineer

---

## 🎯 快速访问

### Grafana 仪表板（推荐）
- **URL：** http://localhost:3001
- **用户名：** `admin`
- **密码：** `admin123`
- **用途：** 查看所有监控仪表板和告警

### Prometheus 查询界面
- **URL：** http://localhost:9090
- **用途：** 原始指标查询、PromQL 查询、告警规则管理

### Alertmanager
- **URL：** http://localhost:9093
- **用途：** 查看和管理告警

---

## 📊 推荐查看顺序

### 1. 系统概览仪表板
**Grafana → Dashboards → System Overview**

查看内容：
- CPU 使用率趋势
- 内存使用情况
- 磁盘 I/O
- 网络流量
- 容器健康状态

### 2. 容器性能
**Grafana → Dashboards → Container Performance**

查看内容：
- 每个容器的资源使用
- 容器网络流量
- 容器文件系统使用

### 3. API 性能
**Prometheus → Graph**

输入查询：
```promql
# API 请求速率
rate(http_request_total[5m])

# API 响应时间
rate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m])

# API 错误率
rate(http_request_total{status=~"5.."}[5m])
```

### 4. 告警状态
**Alertmanager → Alerts**

查看内容：
- 当前活跃告警
- 告警历史
- 告警分组和路由

---

## 🔔 告警阈值参考

### CPU 告警
- ⚠️ Warning：> 50%
- 🚨 Critical：> 80%

### 内存告警
- ⚠️ Warning：> 70%
- 🚨 Critical：> 90%

### 磁盘告警
- ⚠️ Warning：> 70%
- 🚨 Critical：> 85%

### API 响应时间
- ⚠️ Warning：> 500ms
- 🚨 Critical：> 1000ms

### API 错误率
- ⚠️ Warning：> 1%
- 🚨 Critical：> 5%

---

## 🛠️ 常用 PromQL 查询

### 系统资源
```promql
# CPU 使用率
100 - (avg by (instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)

# 内存使用率
(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100

# 磁盘使用率
(1 - (node_filesystem_avail_bytes / node_filesystem_size_bytes)) * 100
```

### 容器监控
```promql
# 容器 CPU 使用
rate(container_cpu_usage_seconds_total{name!=""}[5m])

# 容器内存使用
container_memory_usage_bytes{name!=""}

# 容器网络接收
rate(container_network_receive_bytes_total{name!=""}[5m])
```

### API 监控
```promql
# QPS（每秒请求数）
rate(http_request_total[1m])

# 平均响应时间
rate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m])

# P95 响应时间
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))
```

---

## 🚨 故障排查流程

### 服务不可用
1. 检查 Grafana System Overview 仪表板
2. 查看 Prometheus Targets 状态
3. 检查容器状态：`docker ps -a`
4. 查看容器日志：`docker logs <container-name>`

### 性能下降
1. 查看 API 响应时间指标
2. 检查数据库连接数
3. 查看慢查询日志
4. 分析 CPU 和内存使用趋势

### 资源耗尽
1. 查看资源使用趋势图
2. 识别资源占用最高的容器
3. 检查是否有内存泄漏
4. 考虑扩容或优化

---

## 📝 日常检查清单

### 每日检查
- [ ] 查看 Grafana System Overview
- [ ] 检查告警通知
- [ ] 确认所有 targets 状态为 UP
- [ ] 查看容器健康状态

### 每周检查
- [ ] 审查告警规则效果
- [ ] 优化告警阈值
- [ ] 分析性能趋势
- [ ] 检查存储空间使用

### 每月检查
- [ ] 审查监控覆盖率
- [ ] 更新仪表板
- [ ] 清理过期数据
- [ ] 优化查询性能

---

## 🔗 相关链接

- **Prometheus 文档：** https://prometheus.io/docs/
- **Grafana 文档：** https://grafana.com/docs/
- **PromQL 简介：** https://prometheus.io/docs/prometheus/latest/querying/basics/

---

## 💡 提示

1. **收藏本页面**，方便快速访问
2. **使用 Grafana 的星标功能**，标记常用仪表板
3. **设置告警静默**，在维护时段避免告警轰炸
4. **定期检查告警规则**，确保阈值符合实际情况
5. **保存常用查询**，提高工作效率

---

**需要帮助？** 联系 DevOps-Engineer
