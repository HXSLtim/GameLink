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
- 用户：
  - `demo.user@gamelink.com` / `User@123456`（普通用户）
  - `pro.player@gamelink.com` / `Player@123456`（陪玩师账号）
- 陪玩师档案：绑定上述陪玩师账号，含基础简介与评分
- 订单 & 支付：一条已完成的示例订单以及对应支付记录
- 评价：对该订单的五星评价

如需定制更多演示数据，可在 `internal/db/seed.go` 中扩展。
