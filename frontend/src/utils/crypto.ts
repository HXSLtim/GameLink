import CryptoJS from 'crypto-js';

/**
 * 加密配置
 */
const CRYPTO_CONFIG = {
  // AES 加密密钥（生产环境应从环境变量读取）
  SECRET_KEY: import.meta.env.VITE_CRYPTO_SECRET_KEY || 'GameLink2025SecretKey!@#',
  // AES 加密向量（生产环境应从环境变量读取）
  IV: import.meta.env.VITE_CRYPTO_IV || 'GameLink2025IV!!!',
  // 是否启用加密（可通过环境变量控制）
  ENABLED: import.meta.env.VITE_CRYPTO_ENABLED !== 'false',
};

/**
 * AES 加密工具类
 */
export class CryptoUtil {
  private static key = CryptoJS.enc.Utf8.parse(CRYPTO_CONFIG.SECRET_KEY);
  private static iv = CryptoJS.enc.Utf8.parse(CRYPTO_CONFIG.IV);

  /**
   * 加密数据
   * @param data 原始数据（对象或字符串）
   * @returns 加密后的 Base64 字符串
   */
  static encrypt(data: unknown): string {
    if (!CRYPTO_CONFIG.ENABLED) {
      return typeof data === 'string' ? data : JSON.stringify(data);
    }

    try {
      // 将数据转换为字符串
      const dataStr = typeof data === 'string' ? data : JSON.stringify(data);

      // 使用 AES 加密
      const encrypted = CryptoJS.AES.encrypt(dataStr, this.key, {
        iv: this.iv,
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7,
      });

      return encrypted.toString();
    } catch (error) {
      console.error('❌ 加密失败:', error);
      throw new Error('数据加密失败');
    }
  }

  /**
   * 解密数据
   * @param encryptedData 加密的 Base64 字符串
   * @returns 解密后的原始数据
   */
  static decrypt<T = unknown>(encryptedData: string): T {
    if (!CRYPTO_CONFIG.ENABLED) {
      try {
        return JSON.parse(encryptedData) as T;
      } catch {
        return encryptedData as T;
      }
    }

    try {
      // 使用 AES 解密
      const decrypted = CryptoJS.AES.decrypt(encryptedData, this.key, {
        iv: this.iv,
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7,
      });

      const decryptedStr = decrypted.toString(CryptoJS.enc.Utf8);

      if (!decryptedStr) {
        throw new Error('解密结果为空');
      }

      // 尝试解析为 JSON
      try {
        return JSON.parse(decryptedStr) as T;
      } catch {
        return decryptedStr as T;
      }
    } catch (error) {
      console.error('❌ 解密失败:', error);
      throw new Error('数据解密失败');
    }
  }

  /**
   * 加密敏感字段（用于部分字段加密）
   * @param obj 包含敏感字段的对象
   * @param fields 需要加密的字段名数组
   * @returns 加密后的对象
   */
  static encryptFields<T extends Record<string, unknown>>(
    obj: T,
    fields: (keyof T)[],
  ): Record<string, unknown> {
    if (!CRYPTO_CONFIG.ENABLED) {
      return obj;
    }

    const result: Record<string, unknown> = { ...obj };

    fields.forEach((field) => {
      const value = obj[field];
      if (value !== undefined && value !== null) {
        result[field as string] = this.encrypt(value);
      }
    });

    return result;
  }

  /**
   * 解密敏感字段
   * @param obj 包含加密字段的对象
   * @param fields 需要解密的字段名数组
   * @returns 解密后的对象
   */
  static decryptFields<T extends Record<string, unknown>>(
    obj: T,
    fields: (keyof T)[],
  ): Record<string, unknown> {
    if (!CRYPTO_CONFIG.ENABLED) {
      return obj;
    }

    const result: Record<string, unknown> = { ...obj };

    fields.forEach((field) => {
      const value = obj[field];
      if (typeof value === 'string') {
        try {
          result[field as string] = this.decrypt(value);
        } catch (error) {
          console.error(`解密字段 ${String(field)} 失败:`, error);
        }
      }
    });

    return result;
  }

  /**
   * MD5 哈希（用于生成签名）
   * @param data 原始数据
   * @returns MD5 哈希值
   */
  static md5(data: string): string {
    return CryptoJS.MD5(data).toString();
  }

  /**
   * SHA256 哈希
   * @param data 原始数据
   * @returns SHA256 哈希值
   */
  static sha256(data: string): string {
    return CryptoJS.SHA256(data).toString();
  }

  /**
   * 生成请求签名（用于防篡改）
   * @param data 请求数据
   * @param timestamp 时间戳
   * @returns 签名字符串
   */
  static generateSignature(data: unknown, timestamp: number): string {
    const dataStr = typeof data === 'string' ? data : JSON.stringify(data);
    const signStr = `${dataStr}${timestamp}${CRYPTO_CONFIG.SECRET_KEY}`;
    return this.sha256(signStr);
  }

  /**
   * 检查是否启用加密
   */
  static isEnabled(): boolean {
    return CRYPTO_CONFIG.ENABLED;
  }
}

/**
 * 导出便捷方法
 */
export const { encrypt, decrypt, encryptFields, decryptFields, md5, sha256, generateSignature } =
  CryptoUtil;

/**
 * 导出加密状态
 */
export const isCryptoEnabled = CryptoUtil.isEnabled();
