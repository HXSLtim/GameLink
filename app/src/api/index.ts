/**
 * API 统一导出入口
 */

// 基础请求
export * from './request'

// 认证相关
export * from './auth'

// 聊天相关
export * from './chat'

// 通知相关
export * from './notification'

// 陪玩师相关 (管理端)
import * as playerApi from './player'
export { playerApi }
// 显式导出 player 中不冲突的类型
export type {
  PlayerService,
  CreateServiceData,
  PlayerSchedule,
  TodayStats,
  OverviewStats
} from './player'

// 用户设置
export * from './settings'

// 订单相关
export * from './order'

// 钱包相关
export * from './wallet'

// 游戏相关
export * from './game'

// 收藏相关
export * from './favorite'

// 公开陪玩师相关 (用户端)
import * as publicPlayerApi from './publicPlayer'
export { publicPlayerApi }
// 导出类型
export type {
  PlayerInfo,
  PlayerDetail,
  PlayerServiceItem,
  PlayerListParams
} from './publicPlayer'

// 用户信息相关
import * as userApi from './user'
export { userApi }

// 陪玩认证相关
import * as certificationApi from './certification'
export { certificationApi }
