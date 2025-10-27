package model

// Package model defines shared domain entities and DTO structures for GameLink.
//
// The initial core domain includes:
//   - User: platform account with roles (user/player/admin)
//   - Player: pro/陪玩资料，关联 User，包含技能、认证与定价
//   - Game: 支持的游戏元数据
//   - Order: 订单，关联用户、打手与游戏，含状态机
//   - Payment: 支付记录（微信/支付宝），含支付状态
//   - Review: 评价与评分
//
// JSON 字段遵循 snake_case 命名，枚举使用小写字符串以便前后端一致。
