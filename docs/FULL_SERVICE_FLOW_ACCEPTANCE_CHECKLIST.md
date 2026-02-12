# GameLink 全服务流程验收清单

## 目标
- 一次性验证用户端、陪玩师端、客服/后台、数据一致性四条主链路。
- 给出可重复执行的脚本化结果，避免人工点点点。

## 前置条件
- Docker 依赖已启动：
  - `docker compose up -d`
- 后端已启动并可访问：
  - `cd api`
  - `go run cmd/main.go`
- 默认验收地址：
  - `http://127.0.0.1:8080/api/v1`

## 一键验收命令
- 根目录执行：
  - `powershell -ExecutionPolicy Bypass -File api/scripts/run_full_service_flow_acceptance.ps1`

## 脚本覆盖范围
- `Smoke`：前端关键 API 兼容检查（含 `/users/*` 旧路径）
- `OrderFlow`：下单/支付/完成/评价主链路守卫
- `DisputeFlow`：客服角色权限 + 争议处理链路
- `WithdrawFlow`：提现申请/审核/打款/流水校验
- `Integrity`：核心外键和业务一致性归零校验

## 判定标准（必须全部满足）
- 汇总表所有检查项均为 `PASS`。
- 最终输出包含：`[full-acceptance] 全部通过`。
- `Integrity` 阶段 `violations=0`。

## 常见失败定位
- `401/403`：账号密码错误或角色权限未同步（先检查 seed 和 RBAC）。
- `404`：后端未加载最新代码（重启 `go run cmd/main.go`）。
- `connection refused`：后端未启动或端口冲突。
- `violations>0`：先执行
  - `powershell -ExecutionPolicy Bypass -File api/scripts/run_data_integrity.ps1 -Fix`
  - 再重跑一键验收脚本。

## 默认验收账号
- 管理员：`admin@gamelink.com / Admin123456`
- 用户：`demo.user@gamelink.com / User@123456`
- 陪玩师：`pro.player@gamelink.com / Player@123456`
- 客服主管：`cs.leader@gamelink.com / CsLeader@123`
- 客服专员：`cs.agent@gamelink.com / CsAgent@123`
