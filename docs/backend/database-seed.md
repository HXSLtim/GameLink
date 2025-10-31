# 数据种子开关

为方便前端在本地或测试环境快速看到示例数据，后端提供可配置的种子数据生成：

- 启动时在迁移完成后自动写入基础游戏、普通用户、陪玩师、示例订单/支付/评价等记录。
- 所有插入都具备幂等性，多次运行不会重复创建。
- 默认仅在开发配置 (`configs/config.development.yaml`) 中启用，可通过配置或环境变量控制。

## 配置方式

`configs/config.<env>.yaml` 中新增节点：

```yaml
seed:
  enabled: true   # 开启自动种子
```

或通过环境变量：

```bash
export SEED_ENABLED=true
```

生产环境建议保持关闭（默认 `false`），以免污染真实数据。

## 种子内容

- 游戏：示例游戏（英雄联盟、DOTA 2、无畏契约）
- 用户（均为演示用途，密码仅限本地环境）：
  1. `demo.user@gamelink.com` / `User@123456`（普通用户）
  2. `pro.player@gamelink.com` / `Player@123456`（陪玩师账号）
  3. `vip.user@gamelink.com` / `Vip@123456`（高级会员用户）
  4. `new.user@gamelink.com` / `User@123789`（体验用户）
  5. `streamer@gamelink.com` / `Player@654321`（陪玩主播账号）
- 陪玩师档案：为陪玩师账号生成两条档案，涵盖 MOBA/FPS 不同领域
- 订单：4 笔不同状态的订单（已完成、进行中、待开始、已取消）
- 支付：针对订单生成对应支付记录（包含已支付、待支付、已退款等状态）
- 评价：示例评价 2 条，对应不同订单状态

如需定制更多演示数据，可在 `internal/db/seed.go` 中扩展。
