/**
 * 格式化工具函数
 */

/**
 * 格式化金额（分转元）
 */
export function formatMoney(cents: number, showSign = false): string {
  const yuan = cents / 100
  const formatted = yuan.toFixed(2)
  if (showSign && cents > 0) {
    return `+${formatted}`
  }
  return formatted
}

/**
 * 格式化金额（元）
 */
export function formatYuan(amount: number): string {
  if (!Number.isFinite(amount)) {
    return '0.00'
  }
  return amount.toFixed(2)
}

/**
 * 格式化金额（带货币符号）
 */
export function formatPrice(cents: number): string {
  return `¥${formatMoney(cents)}`
}

/**
 * 格式化数字（大数简化）
 */
export function formatNumber(num: number): string {
  if (num >= 10000) {
    return `${(num / 10000).toFixed(1)}w`
  }
  if (num >= 1000) {
    return `${(num / 1000).toFixed(1)}k`
  }
  return String(num)
}

/**
 * 格式化人数/数量（中文单位）
 */
export function formatCount(count?: number): string {
  if (!count) return '0'
  if (count >= 10000) {
    return `${(count / 10000).toFixed(1)}万`
  }
  if (count >= 1000) {
    return `${(count / 1000).toFixed(1)}k`
  }
  return String(count)
}

/**
 * 格式化时间戳为相对时间
 */
export function formatRelativeTime(timestamp: number | string | Date): string {
  const date = typeof timestamp === 'object' ? timestamp : new Date(timestamp)
  const now = Date.now()
  const diff = now - date.getTime()
  
  const minute = 60 * 1000
  const hour = 60 * minute
  const day = 24 * hour
  const week = 7 * day
  const month = 30 * day
  
  if (diff < minute) {
    return '刚刚'
  }
  if (diff < hour) {
    return `${Math.floor(diff / minute)}分钟前`
  }
  if (diff < day) {
    return `${Math.floor(diff / hour)}小时前`
  }
  if (diff < week) {
    return `${Math.floor(diff / day)}天前`
  }
  if (diff < month) {
    return `${Math.floor(diff / week)}周前`
  }
  
  // 超过一个月显示日期
  const m = date.getMonth() + 1
  const d = date.getDate()
  return `${m}月${d}日`
}

/**
 * 格式化日期时间
 */
export function formatDateTime(
  timestamp: number | string | Date, 
  format = 'YYYY-MM-DD HH:mm'
): string {
  const date = typeof timestamp === 'object' ? timestamp : new Date(timestamp)
  
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hour = String(date.getHours()).padStart(2, '0')
  const minute = String(date.getMinutes()).padStart(2, '0')
  const second = String(date.getSeconds()).padStart(2, '0')
  
  return format
    .replace('YYYY', String(year))
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hour)
    .replace('mm', minute)
    .replace('ss', second)
}

/**
 * 安全格式化日期时间（支持空值回退）
 */
export function formatDateTimeSafe(
  timestamp?: number | string | Date,
  format = 'YYYY-MM-DD HH:mm',
  fallback = '-'
): string {
  if (timestamp === undefined || timestamp === null || timestamp === '') {
    return fallback
  }
  return formatDateTime(timestamp, format)
}

/**
 * 格式化日期（仅日期）
 */
export function formatDate(timestamp: number | string | Date): string {
  return formatDateTime(timestamp, 'YYYY-MM-DD')
}

/**
 * 格式化时间（仅时间）
 */
export function formatTime(timestamp: number | string | Date): string {
  return formatDateTime(timestamp, 'HH:mm')
}

/**
 * 格式化日期（月/日）
 */
export function formatMonthDay(timestamp: number | string | Date): string {
  const date = typeof timestamp === 'object' ? timestamp : new Date(timestamp)
  const m = date.getMonth() + 1
  const d = date.getDate()
  return `${m}/${d}`
}

/**
 * 格式化日期（中文月日）
 */
export function formatDateChinese(timestamp: number | string | Date): string {
  return `${formatMonthDay(timestamp).replace('/', '月')}日`
}

/**
 * 格式化日期时间（月/日 + 时间）
 */
export function formatMonthDayTime(timestamp: number | string | Date): string {
  const date = typeof timestamp === 'object' ? timestamp : new Date(timestamp)
  const hour = String(date.getHours()).padStart(2, '0')
  const minute = String(date.getMinutes()).padStart(2, '0')
  return `${formatMonthDay(date)} ${hour}:${minute}`
}

/**
 * 相对时间（短格式，超过一周显示月/日）
 */
export function formatRelativeTimeShort(timestamp: number | string | Date): string {
  const date = typeof timestamp === 'object' ? timestamp : new Date(timestamp)
  const now = Date.now()
  const diff = now - date.getTime()
  
  const minute = 60 * 1000
  const hour = 60 * minute
  const day = 24 * hour
  
  if (diff < minute) {
    return '刚刚'
  }
  if (diff < hour) {
    return `${Math.floor(diff / minute)}分钟前`
  }
  if (diff < day) {
    return `${Math.floor(diff / hour)}小时前`
  }
  if (diff < 7 * day) {
    return `${Math.floor(diff / day)}天前`
  }
  
  return formatMonthDay(date)
}

/**
 * 格式化时长（秒转可读格式）
 */
export function formatDuration(seconds: number): string {
  if (seconds < 60) {
    return `${seconds}秒`
  }
  if (seconds < 3600) {
    const minutes = Math.floor(seconds / 60)
    const secs = seconds % 60
    return secs > 0 ? `${minutes}分${secs}秒` : `${minutes}分钟`
  }
  
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  return minutes > 0 ? `${hours}小时${minutes}分` : `${hours}小时`
}

/**
 * 格式化文件大小
 */
export function formatFileSize(bytes: number): string {
  if (bytes < 1024) {
    return `${bytes}B`
  }
  if (bytes < 1024 * 1024) {
    return `${(bytes / 1024).toFixed(1)}KB`
  }
  if (bytes < 1024 * 1024 * 1024) {
    return `${(bytes / 1024 / 1024).toFixed(1)}MB`
  }
  return `${(bytes / 1024 / 1024 / 1024).toFixed(2)}GB`
}

/**
 * 手机号脱敏
 */
export function maskPhone(phone: string): string {
  if (!phone || phone.length !== 11) return phone
  return `${phone.slice(0, 3)}****${phone.slice(7)}`
}

/**
 * 姓名脱敏
 */
export function maskName(name: string): string {
  if (!name) return ''
  if (name.length <= 1) return '*'
  if (name.length === 2) return `${name[0]}*`
  return `${name[0]}${'*'.repeat(name.length - 2)}${name[name.length - 1]}`
}
