/**
 * 验证工具函数
 */

/**
 * 验证手机号
 */
export function isValidPhone(phone: string): boolean {
  return /^1[3-9]\d{9}$/.test(phone)
}

/**
 * 验证邮箱
 */
export function isValidEmail(email: string): boolean {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
}

/**
 * 验证密码强度（至少6位，包含数字和字母）
 */
export function isValidPassword(password: string): boolean {
  return /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d@$!%*#?&]{6,20}$/.test(password)
}

/**
 * 验证身份证号
 */
export function isValidIdCard(idCard: string): boolean {
  // 18位身份证
  if (!/^\d{17}[\dXx]$/.test(idCard)) {
    return false
  }
  
  // 校验位验证
  const weights = [7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2]
  const checkCodes = ['1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2']
  
  let sum = 0
  for (let i = 0; i < 17; i++) {
    sum += parseInt(idCard[i]) * weights[i]
  }
  
  const checkCode = checkCodes[sum % 11]
  return idCard[17].toUpperCase() === checkCode
}

/**
 * 验证银行卡号（Luhn算法）
 */
export function isValidBankCard(cardNo: string): boolean {
  if (!/^\d{16,19}$/.test(cardNo)) {
    return false
  }
  
  let sum = 0
  const digits = cardNo.split('').reverse()
  
  for (let i = 0; i < digits.length; i++) {
    let digit = parseInt(digits[i])
    if (i % 2 === 1) {
      digit *= 2
      if (digit > 9) {
        digit -= 9
      }
    }
    sum += digit
  }
  
  return sum % 10 === 0
}

/**
 * 验证URL
 */
export function isValidUrl(url: string): boolean {
  try {
    new URL(url)
    return true
  } catch {
    return false
  }
}

/**
 * 验证是否为空（null、undefined、空字符串、空数组、空对象）
 */
export function isEmpty(value: unknown): boolean {
  if (value === null || value === undefined) return true
  if (typeof value === 'string') return value.trim() === ''
  if (Array.isArray(value)) return value.length === 0
  if (typeof value === 'object') return Object.keys(value).length === 0
  return false
}

/**
 * 验证是否为数字
 */
export function isNumber(value: unknown): value is number {
  return typeof value === 'number' && !isNaN(value)
}

/**
 * 验证金额（正数，最多两位小数）
 */
export function isValidAmount(amount: string | number): boolean {
  const num = typeof amount === 'string' ? parseFloat(amount) : amount
  if (isNaN(num) || num <= 0) return false
  
  const str = String(amount)
  const decimalIndex = str.indexOf('.')
  if (decimalIndex === -1) return true
  
  return str.length - decimalIndex - 1 <= 2
}

/**
 * 验证验证码（6位数字）
 */
export function isValidVerifyCode(code: string): boolean {
  return /^\d{6}$/.test(code)
}

/**
 * 验证昵称（2-12个字符，支持中英文数字）
 */
export function isValidNickname(nickname: string): boolean {
  const length = nickname.length
  // 中文算2个字符
  let charCount = 0
  for (const char of nickname) {
    charCount += /[\u4e00-\u9fa5]/.test(char) ? 2 : 1
  }
  return charCount >= 2 && charCount <= 24 && /^[\u4e00-\u9fa5a-zA-Z0-9_]+$/.test(nickname)
}
