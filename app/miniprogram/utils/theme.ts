/**
 * 主题切换工具
 */

import { getStorage, setStorage, StorageKeys } from './storage'

export type Identity = 'user' | 'player'

/**
 * 获取当前身份
 */
export function getIdentity(): Identity {
  return getStorage<Identity>(StorageKeys.IDENTITY) || 'user'
}

/**
 * 设置身份
 */
export function setIdentity(identity: Identity): boolean {
  return setStorage(StorageKeys.IDENTITY, identity)
}

/**
 * 切换身份
 */
export function toggleIdentity(): Identity {
  const current = getIdentity()
  const next: Identity = current === 'user' ? 'player' : 'user'
  setIdentity(next)
  return next
}

/**
 * 获取主题类名
 */
export function getThemeClass(): string {
  const identity = getIdentity()
  return identity === 'player' ? 'theme-player' : 'theme-user'
}

/**
 * 是否为陪玩师模式
 */
export function isPlayerMode(): boolean {
  return getIdentity() === 'player'
}

/**
 * 是否为用户模式
 */
export function isUserMode(): boolean {
  return getIdentity() === 'user'
}
