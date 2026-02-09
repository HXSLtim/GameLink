/**
 * Pinia Store 统一导出
 */

import { createPinia } from 'pinia'

// 创建 Pinia 实例
const pinia = createPinia()

export default pinia

// 导出所有 Store
export { useUserStore, normalizeUserInfo } from './user'
export type { UserInfo } from './user'
export { useAppStore } from './app'
export { usePlayerStore } from './player'
