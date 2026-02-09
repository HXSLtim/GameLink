/**
 * 工具函数统一导出
 */

export * from './format'
export * from './validate'
export * from './storage'

/**
 * 防抖函数
 */
export function debounce<T extends (...args: unknown[]) => unknown>(
  fn: T,
  delay: number
): (...args: Parameters<T>) => void {
  let timer: ReturnType<typeof setTimeout> | null = null
  
  return function (this: unknown, ...args: Parameters<T>) {
    if (timer) {
      clearTimeout(timer)
    }
    timer = setTimeout(() => {
      fn.apply(this, args)
      timer = null
    }, delay)
  }
}

/**
 * 节流函数
 */
export function throttle<T extends (...args: unknown[]) => unknown>(
  fn: T,
  delay: number
): (...args: Parameters<T>) => void {
  let lastTime = 0
  
  return function (this: unknown, ...args: Parameters<T>) {
    const now = Date.now()
    if (now - lastTime >= delay) {
      fn.apply(this, args)
      lastTime = now
    }
  }
}

/**
 * 延迟执行
 */
export function sleep(ms: number): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms))
}

/**
 * 生成 UUID
 */
export function uuid(): string {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0
    const v = c === 'x' ? r : (r & 0x3) | 0x8
    return v.toString(16)
  })
}

/**
 * 深拷贝
 */
export function deepClone<T>(obj: T): T {
  if (obj === null || typeof obj !== 'object') {
    return obj
  }
  
  if (Array.isArray(obj)) {
    return obj.map(item => deepClone(item)) as unknown as T
  }
  
  const cloned = {} as T
  for (const key in obj) {
    if (Object.prototype.hasOwnProperty.call(obj, key)) {
      cloned[key] = deepClone(obj[key])
    }
  }
  
  return cloned
}

/**
 * 获取平台信息
 */
export function getPlatform(): 'h5' | 'mp-weixin' | 'app' | 'unknown' {
  // #ifdef H5
  return 'h5'
  // #endif
  
  // #ifdef MP-WEIXIN
  return 'mp-weixin'
  // #endif
  
  // #ifdef APP-PLUS
  return 'app'
  // #endif
  
  // #ifndef H5 || MP-WEIXIN || APP-PLUS
  return 'unknown'
  // #endif
}

/**
 * 复制到剪贴板
 */
export function copyToClipboard(text: string): Promise<void> {
  return new Promise((resolve, reject) => {
    uni.setClipboardData({
      data: text,
      success: () => {
        uni.showToast({ title: '已复制', icon: 'success' })
        resolve()
      },
      fail: reject,
    })
  })
}

/**
 * 拨打电话
 */
export function makePhoneCall(phoneNumber: string): void {
  uni.makePhoneCall({
    phoneNumber,
    fail: () => {
      uni.showToast({ title: '拨打失败', icon: 'none' })
    },
  })
}

/**
 * 预览图片
 */
export function previewImage(urls: string[], current: number | string = 0): void {
  if (!urls || urls.length === 0) return
  const resolvedCurrent = typeof current === 'number' ? (urls[current] || urls[0]) : current
  uni.previewImage({
    urls,
    current: resolvedCurrent,
  })
}

/**
 * 保存图片到相册
 */
export function saveImageToPhotosAlbum(filePath: string): Promise<void> {
  return new Promise((resolve, reject) => {
    uni.saveImageToPhotosAlbum({
      filePath,
      success: () => {
        uni.showToast({ title: '保存成功', icon: 'success' })
        resolve()
      },
      fail: (err) => {
        if (err.errMsg?.includes('auth deny')) {
          uni.showToast({ title: '请授权相册权限', icon: 'none' })
        } else {
          uni.showToast({ title: '保存失败', icon: 'none' })
        }
        reject(err)
      },
    })
  })
}
