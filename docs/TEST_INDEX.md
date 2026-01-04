# GameLink 管理后台全量测试文档索引

**文档版本**: v1.0  
**创建日期**: 2024-12-18  
**测试环境**: Docker 生产环境 (docker-compose.prod.yml)

> 自动化测试（后端集成测试 / 管理后台 E2E）环境搭建请参考：[TEST_ENVIRONMENT_SETUP.md](TEST_ENVIRONMENT_SETUP.md)。

---

## 测试文档清单

按依赖关系排序，建议按顺序执行测试：

### 第一阶段：基础配置模块（无依赖）

| 序号 | 模块名称 | 文档路径 | 按钮数量 | 优先级 |
|------|---------|---------|---------|--------|
| 01 | 用户标签管理 | [TEST_MODULE_01_UserTag.md](TEST_MODULE_01_UserTag.md) | 6 | P0 |
| 02 | 分流规则管理 | [TEST_MODULE_02_RoutingRule.md](TEST_MODULE_02_RoutingRule.md) | 7 | P0 |
| 03 | 结算公司管理 | [TEST_MODULE_03_SettlementCompany.md](TEST_MODULE_03_SettlementCompany.md) | 6 | P0 |
| 04 | 排行榜抽成配置 | [TEST_MODULE_04_RankingCommission.md](TEST_MODULE_04_RankingCommission.md) | 6 | P1 |
| 05 | 提现分流统计 | [TEST_MODULE_05_WithdrawRouting.md](TEST_MODULE_05_WithdrawRouting.md) | 4 | P1 |
| 06 | 用户行为分析 | [TEST_MODULE_06_UserBehavior.md](TEST_MODULE_06_UserBehavior.md) | 5 | P1 |

### 第二阶段：权限与用户模块

| 序号 | 模块名称 | 文档路径 | 按钮数量 | 优先级 |
|------|---------|---------|---------|--------|
| 21 | 权限管理 | [TEST_MODULE_21_Permission.md](TEST_MODULE_21_Permission.md) | 6 | P0 |
| 08 | 角色管理 | [TEST_MODULE_08_Role.md](TEST_MODULE_08_Role.md) | 8 | P0 |
| 09 | 用户管理 | [TEST_MODULE_09_User.md](TEST_MODULE_09_User.md) | 10 | P0 |

### 第三阶段：业务主体模块

| 序号 | 模块名称 | 文档路径 | 按钮数量 | 优先级 |
|------|---------|---------|---------|--------|
| 10 | 游戏管理 | [TEST_MODULE_10_Game.md](TEST_MODULE_10_Game.md) | 8 | P0 |
| 11 | 陪玩师管理 | [TEST_MODULE_11_Player.md](TEST_MODULE_11_Player.md) | 12 | P0 |
| 12 | 服务项管理 | [TEST_MODULE_12_Service.md](TEST_MODULE_12_Service.md) | 8 | P0 |

### 第四阶段：核心业务模块

| 序号 | 模块名称 | 文档路径 | 按钮数量 | 优先级 |
|------|---------|---------|---------|--------|
| 13 | 订单管理 | [TEST_MODULE_13_Order.md](TEST_MODULE_13_Order.md) | 12 | P0 |
| 07 | 纠纷管理 | [TEST_MODULE_07_Dispute.md](TEST_MODULE_07_Dispute.md) | 8 | P0 |
| 14 | 评价管理 | [TEST_MODULE_14_Review.md](TEST_MODULE_14_Review.md) | 8 | P0 |

### 第五阶段：财务模块

| 序号 | 模块名称 | 文档路径 | 按钮数量 | 优先级 |
|------|---------|---------|---------|--------|
| 15 | 提现管理 | [TEST_MODULE_15_Withdraw.md](TEST_MODULE_15_Withdraw.md) | 10 | P0 |
| 16 | 佣金管理 | [TEST_MODULE_16_Commission.md](TEST_MODULE_16_Commission.md) | 8 | P0 |

### 第六阶段：运营支撑模块

| 序号 | 模块名称 | 文档路径 | 按钮数量 | 优先级 |
|------|---------|---------|---------|--------|
| 20 | 内容管理 | [TEST_MODULE_20_Content.md](TEST_MODULE_20_Content.md) | 10 | P1 |
| 19 | 通知管理 | [TEST_MODULE_19_Notifications.md](TEST_MODULE_19_Notifications.md) | 8 | P1 |
| 18 | 告警管理 | [TEST_MODULE_18_Alert.md](TEST_MODULE_18_Alert.md) | 8 | P1 |
| 17 | 仪表盘 | [TEST_MODULE_17_Dashboard.md](TEST_MODULE_17_Dashboard.md) | 6 | P0 |
| 22 | 系统监控 | [TEST_MODULE_22_Monitor.md](TEST_MODULE_22_Monitor.md) | 15 | P1 |

---

## 测试统计

| 统计项 | 数量 |
|--------|------|
| 总模块数 | 22 |
| 总按钮数 | ~169 |
| P0 优先级模块 | 14 |
| P1 优先级模块 | 8 |

---

## 测试环境信息

### Docker 环境

```powershell
# 启动环境
docker compose -f docker-compose.prod.yml up -d

# 检查状态
docker compose -f docker-compose.prod.yml ps
```

### 访问地址

| 服务 | 地址 |
|------|------|
| 前端 | http://localhost:80 |
| 后端 API | http://localhost:8081/api/v1 |
| 健康检查 | http://localhost:8081/api/v1/healthz |

### 测试账号

| 角色 | 配置来源 | 说明 |
|------|----------|------|
| 超级管理员 | `.env` 中的 `SUPER_ADMIN_EMAIL` / `SUPER_ADMIN_PASSWORD` | `docker-compose.prod.yml` 会注入到后端容器 |

### 数据库连接

```powershell
docker exec -it gamelink-postgres psql -U gamelink -d gamelink
```

### Redis 连接

```powershell
docker exec -it gamelink-redis redis-cli -a <REDIS_PASSWORD>
```

---

## 测试执行流程

1. **环境准备**: 确保 Docker 环境正常运行
2. **按阶段执行**: 按上述阶段顺序执行测试
3. **证据收集**: 每个按钮收集 5 张截图 + 日志
4. **异常测试**: 每个按钮至少 3 个异常场景
5. **签字确认**: 测试人和组长签字
6. **归档**: 测试报告和证据打包归档

---

## 测试报告提交要求

每个模块测试完成后，需提交：

1. ✅ 填写完整的测试文档（所有 ☐ 改为 ☑️ 或标注失败原因）
2. ✅ 截图文件夹（命名格式：`MODULE_XX_按钮ID_步骤.png`）
3. ✅ 日志文件（`MODULE_XX_logs.tar.gz`）
4. ✅ 测试人签字
5. ✅ 组长审核签字

---

## 联系方式

如有问题，请联系：
- 技术负责人: [姓名]
- 测试组长: [姓名]

---

**文档维护**: 如有模块变更，请同步更新对应测试文档
