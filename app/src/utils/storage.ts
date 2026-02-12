/**
 * 本地存储工具
 */

const PREFIX = 'gamelink_'

/**
 * 存储数据
 */
export function setStorage<T>(key: string, value: T): void {
  try {
    const data = JSON.stringify(value)
    uni.setStorageSync(PREFIX + key, data)
  } catch (error) {
    console.error('Storage set error:', error)
  }
}

/**
 * 获取数据
 */
export function getStorage<T>(key: string, defaultValue?: T): T | undefined {
  try {
    const data = uni.getStorageSync(PREFIX + key)
    if (!data) return defaultValue
    return JSON.parse(data) as T
  } catch (error) {
    console.error('Storage get error:', error)
    return defaultValue
  }
}

/**
 * 删除数据
 */
export function removeStorage(key: string): void {
  try {
    uni.removeStorageSync(PREFIX + key)
  } catch (error) {
    console.error('Storage remove error:', error)
  }
}

/**
 * 清空所有数据
 */
export function clearStorage(): void {
  try {
    uni.clearStorageSync()
  } catch (error) {
    console.error('Storage clear error:', error)
  }
}

/**
 * 获取存储信息
 */
export function getStorageInfo(): UniNamespace.GetStorageInfoSuccess | null {
  try {
    return uni.getStorageInfoSync()
  } catch (error) {
    console.error('Storage info error:', error)
    return null
  }
}

/**
 * 设置带过期时间的数据
 */
export function setStorageWithExpiry<T>(key: string, value: T, ttlMs: number): void {
  const item = {
    value,
    expiry: Date.now() + ttlMs,
  }
  setStorage(key, item)
}

/**
 * 获取带过期时间的数据
 */
export function getStorageWithExpiry<T>(key: string): T | null {
  const item = getStorage<{ value: T; expiry: number }>(key)
  if (!item) return null
  
  if (Date.now() > item.expiry) {
    removeStorage(key)
    return null
  }
  
  return item.value
}
