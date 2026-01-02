import CryptoJS from 'crypto-js';

import { logger } from '@/utils/logger';
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
interface EncryptedRequest {
    encrypted: boolean;
    payload: string;
    timestamp: number;
    signature: string;
}

/**
 * 加密配置错误（当加密未正确配置时抛出）
 */
export class CryptoConfigError extends Error {
    constructor(message: string) {
        super(message);
        this.name = 'CryptoConfigError';
    }
}

/**
 * 获取加密配置
 * @throws {CryptoConfigError} 如果加密已启用但密钥未配置
 */
function getCryptoConfig(): CryptoConfig {
    const enabled = import.meta.env.VITE_CRYPTO_ENABLED === 'true';
    const secretKey = import.meta.env.VITE_CRYPTO_SECRET_KEY || '';
    const iv = import.meta.env.VITE_CRYPTO_IV || '';
    const useSignature = import.meta.env.VITE_CRYPTO_USE_SIGNATURE !== 'false';

    // 如果加密已启用但缺少密钥，抛出错误
    if (enabled) {
        if (!secretKey) {
            throw new CryptoConfigError(
                'VITE_CRYPTO_SECRET_KEY is not configured. ' +
                'Encryption is enabled but the secret key is missing. ' +
                'Please set VITE_CRYPTO_SECRET_KEY in your .env file.'
            );
        }
        if (!iv) {
            throw new CryptoConfigError(
                'VITE_CRYPTO_IV is not configured. ' +
                'Encryption is enabled but the IV is missing. ' +
                'Please set VITE_CRYPTO_IV in your .env file.'
            );
        }
    }

    return { secretKey, iv, enabled, useSignature };
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
 * 生成签名
 */
function generateSignature(plaintext: string, timestamp: number, secretKey: string): string {
    const message = plaintext + timestamp.toString() + secretKey;
    return CryptoJS.SHA256(message).toString(CryptoJS.enc.Hex);
}

/**
 * 加密请求数据
 * @throws {CryptoConfigError} 如果加密已启用但密钥未配置
 */
export function encryptRequest(data: unknown): EncryptedRequest | unknown {
    let config: CryptoConfig;

    try {
        config = getCryptoConfig();
    } catch (error) {
        if (error instanceof CryptoConfigError) {
            // 加密配置错误时，在控制台输出错误信息并返回原始数据
            logger.error('Encryption configuration error:', error.message);
            logger.warn('Request will be sent without encryption. Please fix your environment configuration.');
            return data;
        }
        throw error;
    }

    // 如果未启用加密，直接返回原始数据
    if (!config.enabled) {
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
        logger.error('Failed to encrypt request:', error);
        // 加密失败时返回原始数据
        return data;
    }
}

/**
 * 判断是否需要加密
 */
export function shouldEncrypt(method: string, url: string): boolean {
    try {
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
            '/api/v1/ping',
            '/api/v1/auth/refresh',
        ];

        return !excludePaths.some(path => url.includes(path));
    } catch (error) {
        // 如果配置检查失败，不加密
        logger.warn('Crypto configuration check failed, disabling encryption:', error);
        return false;
    }
}

/**
 * 检查加密是否已配置
 * @returns true 如果加密已正确配置
 * @throws {CryptoConfigError} 如果加密已启用但密钥未配置
 */
export function isCryptoConfigured(): boolean {
    try {
        const config = getCryptoConfig();
        return config.enabled && !!config.secretKey && !!config.iv;
    } catch (error) {
        if (error instanceof CryptoConfigError) {
            logger.warn('Crypto configuration check failed:', error.message);
            return false;
        }
        throw error;
    }
}
