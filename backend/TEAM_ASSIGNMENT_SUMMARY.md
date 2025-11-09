# 团队分配快速参考

## 👥 3人团队分配 (推荐)

### 成员A: Handler层 (17-19小时)
**任务**: 所有Handler层测试
- Admin Handler核心: 9个文件 (13.5h)
- User/Player Handler: 5个文件 (3.5h)
- **预计提升**: +13%

### 成员B: Service+Repository (7小时)
**任务**: Service层和Repository层
- Service层: 4个文件 (5.5h)
- Repository层: 3个文件 (1.5h)
- **预计提升**: +11%

### 成员C: 小模块 (4.5-5.5小时)
**任务**: 小模块和工具类
- 7个小模块测试 (3.5h)
- 边界条件测试 (1-2h)
- **预计提升**: +5%

---

## 📅 3天时间线

**第1天**: 35.5% → 50% (+14.5%)
- 成员B: Service层 (2h)
- 成员C: 小模块 (2h)
- 成员A: Admin Handler开始 (4h)

**第2天**: 50% → 63% (+13%)
- 成员A: Admin Handler继续 (6h)
- 成员B: User/Player Handler (2h)

**第3天**: 63% → 70%+ (+7%)
- 成员A: Admin Handler完成 (2h)
- 成员B: Repository层 (1.5h)
- 全体: 代码审查和优化 (4h)

---

## 📋 详细任务清单

### 成员A任务清单
1. `handler/admin/game_test.go` (1.5h)
2. `handler/admin/user_test.go` (2h)
3. `handler/admin/player_test.go` (1.5h)
4. `handler/admin/order_test.go` (3h)
5. `handler/admin/payment_test.go` (1.5h)
6. `handler/admin/review_test.go` (1h)
7. `handler/admin/role_test.go` (1h)
8. `handler/admin/permission_test.go` (1h)
9. `handler/admin/helpers_test.go` (1h)
10. `handler/user/order_test.go` (1h)
11. `handler/user/payment_test.go` (0.5h)
12. `handler/user/review_test.go` (0.5h)
13. `handler/user/player_test.go` (0.5h)
14. `handler/player/*_test.go` (1h)

### 成员B任务清单
1. `service/admin/admin_test.go` 增强 (2h)
2. `service/role/role_test.go` 增强 (1.5h)
3. `service/player/player_test.go` 增强 (1h)
4. `service/order/order_test.go` 增强 (1h)
5. `repository/commission/repository_test.go` 增强 (0.5h)
6. `repository/serviceitem/repository_test.go` 增强 (0.5h)
7. `repository/permission/repository_test.go` 增强 (0.5h)

### 成员C任务清单
1. `cache/redis_test.go` (0.5h)
2. `auth/jwt_test.go` 增强 (0.5h)
3. `db/db_test.go` (0.5h)
4. `db/seed_test.go` (0.5h)
5. `logging/logger_test.go` 增强 (0.5h)
6. `metrics/metrics_test.go` 增强 (0.5h)
7. `config/env_test.go` 增强 (0.5h)
8. 边界条件测试 (1-2h)

---

## ✅ 每日检查点

**第1天结束**:
- [ ] Service层 70%+
- [ ] 小模块 50%+
- [ ] 总体覆盖率 50%

**第2天结束**:
- [ ] Admin Handler核心完成
- [ ] 总体覆盖率 63%

**第3天结束**:
- [ ] 所有任务完成
- [ ] 代码审查通过
- [ ] 总体覆盖率 70%+

---

## 📚 参考文档

- **详细任务**: `REMAINING_WORK_FILE_LEVEL.md`
- **完整计划**: `TEAM_ASSIGNMENT_PLAN.md`
- **测试规范**: `.cursor/rules/backend-testing.mdc`

---

**快速开始**: 查看 `TEAM_ASSIGNMENT_PLAN.md` 获取详细指南

