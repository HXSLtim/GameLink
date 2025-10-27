/**
 * 安全的 localStorage 工具类
 * 提供异常处理和类型安全的存储操作
 */

/**
 * 安全地从 localStorage 获取值
 * @param key 存储键名
 * @param defaultValue 默认值
 * @returns 存储的值或默认值
 */
export function getItem<T = string>(key: string, defaultValue: T | null = null): T | null {
  try {
    const item = localStorage.getItem(key);
    if (item === null) return defaultValue;
    
    // 尝试解析 JSON
    try {
      return JSON.parse(item) as T;
    } catch {
      // 如果不是 JSON，返回原始字符串
      return item as T;
    }
  } catch (error) {
    console.error(`[Storage] Failed to get item "${key}":`, error);
    return defaultValue;
  }
}

/**
 * 安全地设置 localStorage 值
 * @param key 存储键名
 * @param value 要存储的值
 * @returns 是否成功
 */
export function setItem<T>(key: string, value: T): boolean {
  try {
    const serialized = typeof value === 'string' ? value : JSON.stringify(value);
    localStorage.setItem(key, serialized);
    return true;
  } catch (error) {
    console.error(`[Storage] Failed to set item "${key}":`, error);
    // 可能是存储空间已满或隐私模式
    if (error instanceof Error) {
      if (error.name === 'QuotaExceededError') {
        console.warn('[Storage] localStorage quota exceeded');
      }
    }
    return false;
  }
}

/**
 * 安全地移除 localStorage 值
 * @param key 存储键名
 * @returns 是否成功
 */
export function removeItem(key: string): boolean {
  try {
    localStorage.removeItem(key);
    return true;
  } catch (error) {
    console.error(`[Storage] Failed to remove item "${key}":`, error);
    return false;
  }
}

/**
 * 安全地清空 localStorage
 * @returns 是否成功
 */
export function clear(): boolean {
  try {
    localStorage.clear();
    return true;
  } catch (error) {
    console.error('[Storage] Failed to clear localStorage:', error);
    return false;
  }
}

/**
 * 检查 localStorage 是否可用
 * @returns 是否可用
 */
export function isAvailable(): boolean {
  try {
    const testKey = '__storage_test__';
    localStorage.setItem(testKey, 'test');
    localStorage.removeItem(testKey);
    return true;
  } catch {
    return false;
  }
}

/**
 * 获取 localStorage 使用情况（仅在支持的浏览器中）
 * @returns 使用情况对象或 null
 */
export function getStorageInfo(): { used: number; quota: number } | null {
  try {
    if ('storage' in navigator && 'estimate' in navigator.storage) {
      // 注意：这是异步的，这里提供同步版本的占位
      return null;
    }
    return null;
  } catch {
    return null;
  }
}

/**
 * 安全的 localStorage 工具对象
 */
export const storage = {
  getItem,
  setItem,
  removeItem,
  clear,
  isAvailable,
  getStorageInfo,
};
