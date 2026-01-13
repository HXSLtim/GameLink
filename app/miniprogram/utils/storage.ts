/**
 * 本地存储工具
 */

const STORAGE_PREFIX = 'gl_'

export const StorageKeys = {
  TOKEN: `${STORAGE_PREFIX}token`,
  REFRESH_TOKEN: `${STORAGE_PREFIX}refresh_token`,
  USER_INFO: `${STORAGE_PREFIX}user_info`,
  IDENTITY: `${STORAGE_PREFIX}identity`, // 'user' | 'player'
  USER_MODE: `${STORAGE_PREFIX}user_mode`, // 'user' | 'player' - 角色模式
  THEME: `${STORAGE_PREFIX}theme`,
} as const

type StorageKey = typeof StorageKeys[keyof typeof StorageKeys]

/**
 * 同步获取存储
 */
export function getStorage<T>(key: StorageKey): T | null {
  try {
    return wx.getStorageSync(key) || null
  } catch {
    return null
  }
}

/**
 * 同步设置存储
 */
export function setStorage<T>(key: StorageKey, value: T): boolean {
  try {
    wx.setStorageSync(key, value)
    return true
  } catch {
    return false
  }
}

/**
 * 同步删除存储
 */
export function removeStorage(key: StorageKey): boolean {
  try {
    wx.removeStorageSync(key)
    return true
  } catch {
    return false
  }
}

/**
 * 清除所有 GameLink 存储
 */
export function clearAllStorage(): boolean {
  try {
    const keys = Object.values(StorageKeys)
    keys.forEach(key => wx.removeStorageSync(key))
    return true
  } catch {
    return false
  }
}

/**
 * 异步获取存储
 */
export function getStorageAsync<T>(key: StorageKey): Promise<T | null> {
  return new Promise((resolve) => {
    wx.getStorage({
      key,
      success: (res) => resolve(res.data as T),
      fail: () => resolve(null),
    })
  })
}

/**
 * 异步设置存储
 */
export function setStorageAsync<T>(key: StorageKey, value: T): Promise<boolean> {
  return new Promise((resolve) => {
    wx.setStorage({
      key,
      data: value,
      success: () => resolve(true),
      fail: () => resolve(false),
    })
  })
}
