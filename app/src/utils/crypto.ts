import CryptoJS from 'crypto-js';

/**
 * 加密配置
 */
interface CryptoConfig {
  secretKey: string;
  iv: string;
  enabled: boolean;
  useSignature: boolean;
}

/**
 * 加密请求体
 */
export interface EncryptedRequest {
  encrypted: boolean;
  payload: string;
  timestamp: number;
  signature: string;
}

/**
 * 获取加密配置
 * Taro 使用 process.env 访问环境变量
 */
function getCryptoConfig(): CryptoConfig {
  return {
    secretKey: process.env.TARO_APP_CRYPTO_SECRET_KEY || '',
    iv: process.env.TARO_APP_CRYPTO_IV || '',
    enabled: process.env.TARO_APP_CRYPTO_ENABLED === 'true',
    useSignature: process.env.TARO_APP_CRYPTO_USE_SIGNATURE === 'true',
  };
}

/**
 * AES-256-CBC 加密
 */
function encryptAES(plaintext: string, secretKey: string, iv: string): string {
  const key = CryptoJS.enc.Utf8.parse(secretKey);
  const ivParsed = CryptoJS.enc.Utf8.parse(iv);

  const encrypted = CryptoJS.AES.encrypt(plaintext, key, {
    iv: ivParsed,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7,
  });

  return encrypted.toString();
}

/**
 * AES-256-CBC 解密
 */
function decryptAES(ciphertext: string, secretKey: string, iv: string): string {
  const key = CryptoJS.enc.Utf8.parse(secretKey);
  const ivParsed = CryptoJS.enc.Utf8.parse(iv);

  const decrypted = CryptoJS.AES.decrypt(ciphertext, key, {
    iv: ivParsed,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7,
  });

  return decrypted.toString(CryptoJS.enc.Utf8);
}

/**
 * 生成签名
 */
function generateSignature(plaintext: string, timestamp: number, secretKey: string): string {
  const message = plaintext + timestamp.toString() + secretKey;
  return CryptoJS.SHA256(message).toString(CryptoJS.enc.Hex);
}

/**
 * 加密请求数据
 */
export function encryptRequest(data: unknown): EncryptedRequest | unknown {
  const config = getCryptoConfig();

  // 如果未启用加密，直接返回原始数据
  if (!config.enabled || !config.secretKey || !config.iv) {
    return data;
  }

  try {
    // 将数据转换为 JSON 字符串
    const plaintext = typeof data === 'string' ? data : JSON.stringify(data);

    // 加密数据
    const payload = encryptAES(plaintext, config.secretKey, config.iv);

    // 生成时间戳
    const timestamp = Date.now();

    // 生成签名
    const signature = config.useSignature
      ? generateSignature(plaintext, timestamp, config.secretKey)
      : '';

    // 返回加密请求体
    return {
      encrypted: true,
      payload,
      timestamp,
      signature,
    };
  } catch (error) {
    console.error('Failed to encrypt request:', error);
    // 加密失败时返回原始数据
    return data;
  }
}

/**
 * 解密响应数据
 */
export function decryptResponse<T = unknown>(data: unknown): T {
  const config = getCryptoConfig();

  // 如果未启用加密或数据不是加密格式，直接返回
  if (!config.enabled || !config.secretKey || !config.iv) {
    return data as T;
  }

  // 检查是否是加密响应
  if (
    typeof data === 'object' &&
    data !== null &&
    'encrypted' in data &&
    (data as { encrypted: boolean }).encrypted === true &&
    'payload' in data
  ) {
    try {
      const encryptedData = data as { payload: string };
      const decrypted = decryptAES(encryptedData.payload, config.secretKey, config.iv);
      return JSON.parse(decrypted) as T;
    } catch (error) {
      console.error('Failed to decrypt response:', error);
      return data as T;
    }
  }

  return data as T;
}

/**
 * 判断是否需要加密
 */
export function shouldEncrypt(method: string, url: string): boolean {
  const config = getCryptoConfig();

  if (!config.enabled) {
    return false;
  }

  // 只加密 POST, PUT, PATCH 请求
  const encryptMethods = ['POST', 'PUT', 'PATCH'];
  if (!encryptMethods.includes(method.toUpperCase())) {
    return false;
  }

  // 排除特定路径
  const excludePaths = [
    '/api/v1/health',
    '/api/v1/healthz',
    '/api/v1/ping',
    '/api/v1/auth/refresh',
  ];

  return !excludePaths.some((path) => url.includes(path));
}

/**
 * 检查加密是否已配置
 */
export function isCryptoConfigured(): boolean {
  const config = getCryptoConfig();
  return config.enabled && !!config.secretKey && !!config.iv;
}
