# 超级管理员初始化说明

应用启动时会在执行数据库迁移的同时调用 `ensureSuperAdmin`，用于自动补齐一条 `role=admin` 的管理员账号，确保后台入口可用。逻辑遵循以下规则：

| 环境变量 | 说明 | 默认值（非生产环境） |
|----------|------|----------------------|
| `SUPER_ADMIN_EMAIL` | 超管邮箱，用作唯一登录标识 | `superAdmin@GameLink.com` |
| `SUPER_ADMIN_PHONE` | 超管手机号，可选 | 空 |
| `SUPER_ADMIN_NAME` | 显示名称 | `Super Admin` |
| `SUPER_ADMIN_PASSWORD` | 登录密码（明文，启动时会自动加密） | `admin123` |

- 若邮箱和手机号均未设置且当前环境为生产（`APP_ENV=production`），迁移会中止并提示必须配置其中一个。
- 在生产环境中同样要求显式提供 `SUPER_ADMIN_PASSWORD`，否则启动会失败，避免使用默认弱口令。
- 当数据库中已存在同邮箱（或手机号）的管理员时不会重复创建，保证迁移幂等。
- 若同时启用了种子数据（`SEED_ENABLED=true`），超管之外还会插入演示用的普通用户/陪玩师账号，便于联调。
