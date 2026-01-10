# 数据种子开关

为方便前端在本地或测试环境快速看到示例数据，后端提供可配置的种子数据生成：

- 启动时在迁移完成后自动写入基础游戏、普通用户、陪玩师、示例订单/支付/评价等记录。
- 所有插入都具备幂等性，多次运行不会重复创建。
- 默认仅在开发配置 (`api/configs/config.development.yaml`) 中启用，可通过配置或环境变量控制。
- 启动时会做一次“关联性自检”（跨表引用一致性检查），若发现关键关联断裂会直接返回错误，便于尽早发现 seed 维护问题。

## 配置方式

`api/configs/config.<env>.yaml` 中新增节点：

```yaml
seed:
  enabled: true   # 开启自动种子
```

或通过环境变量：

```bash
export SEED_ENABLED=true
```

生产环境建议保持关闭（默认 `false`），以免污染真实数据。

## 验证方式（快速自查）

启动后端时观察日志，应该能看到类似 `seed data ensured for demo environment` 的输出。

也可以直接跑一次关联性自检测试（会在内存 SQLite 中执行迁移 + 种子 + 关联校验）：

```bash
cd api
go test ./pkg/db -run TestSeedAssociations -count=1
```

也可以通过数据库快速检查关键表是否已有演示数据（示例 SQL）：

```sql
SELECT COUNT(*) FROM vip_levels;
SELECT COUNT(*) FROM coupon_templates;
SELECT COUNT(*) FROM coupons;
SELECT COUNT(*) FROM recharge_options;
SELECT COUNT(*) FROM recharge_records;
SELECT COUNT(*) FROM activities;
SELECT COUNT(*) FROM teams;
SELECT COUNT(*) FROM referrals;
SELECT COUNT(*) FROM user_blocks;
SELECT COUNT(*) FROM game_ranks;
SELECT COUNT(*) FROM player_rank_records;
SELECT COUNT(*) FROM player_certifications;
SELECT COUNT(*) FROM user_notifications;
```

## 种子内容

- 游戏：多品类示例游戏（MOBA/FPS/RPG/体育/派对等）
- 用户（均为演示用途，密码仅限本地环境；可用于前后端联调与管理后台登录）：
  1. `demo.user@gamelink.com` / `User@123456`（普通用户）
  2. `pro.player@gamelink.com` / `Player@123456`（陪玩师账号）
  3. `vip.user@gamelink.com` / `Vip@123456`（高级会员用户）
  4. `new.user@gamelink.com` / `User@123789`（体验用户）
  5. `streamer@gamelink.com` / `Player@654321`（陪玩主播账号）
- 陪玩师档案：覆盖不同主玩游戏、认证状态、评分与价格档位
- 服务项：覆盖 `solo/team/gift` 多种子类别（含上下架与 VIP 专属价示例）
- 订单：覆盖 `pending/confirmed/in_progress/completed/canceled/refunded` 等状态
- 支付：覆盖 `paid/pending/failed/refunded`，并包含 `wallet/combined` 等支付方式示例
- 评价与举报：覆盖通过/待审/驳回、举报处理等场景
- 用户管理：用户标签、登录历史、行为记录等
- 内容管理：分类、动态、动态举报、订单群聊/消息、客服分配等
- 监控与告警：告警、KPI 目标等示例数据
- 财务相关：钱包余额、佣金规则/记录、提现记录、结算公司与分流规则等
- 营销相关：VIP 等级/配置、优惠券模板与用户优惠券、充值档位/记录、活动与参与记录、推荐与奖励记录
- 其它：用户拉黑、段位配置、陪玩师段位认证、陪玩师实名认证、通知模板/通知记录/定时任务等

如需定制更多演示数据，可在 `api/pkg/db/seed.go` 与 `api/pkg/db/seed_demo_extensions.go` 中扩展。
