/**
 * Encryption utilities for GameLink client
 * Implements AES-256-CBC encryption with SHA-256 signature
 *
 * Security: Production MUST enable encryption (VITE_CRYPTO_ENABLED=true)
 */

import CryptoJS from 'crypto-js';

/**
 * Encryption configuration
 */
interface CryptoConfig {
    secretKey: string;
    iv: string;
    enabled: boolean;
    useSignature: boolean;
}

/**
 * Encrypted request payload
 */
export interface EncryptedRequest {
    encrypted: boolean;
    payload: string;
    timestamp: number;
    signature: string;
}

/**
 * Encryption configuration error
 */
export class CryptoConfigError extends Error {
    constructor(message: string) {
        super(message);
        this.name = 'CryptoConfigError';
    }
}

/**
 * Get encryption configuration from environment
 * @throws {CryptoConfigError} If encryption enabled but keys missing
 */
function getCryptoConfig(): CryptoConfig {
    const enabled = import.meta.env.VITE_CRYPTO_ENABLED === 'true';
    const secretKey = import.meta.env.VITE_CRYPTO_SECRET_KEY || '';
    const iv = import.meta.env.VITE_CRYPTO_IV || '';
    const useSignature = import.meta.env.VITE_CRYPTO_USE_SIGNATURE !== 'false';

    // Production safety check
    if (import.meta.env.PROD && !enabled) {
        console.error('FATAL: Encryption must be enabled in production');
    }

    if (enabled) {
        if (!secretKey) {
            throw new CryptoConfigError(
                'VITE_CRYPTO_SECRET_KEY is not configured. ' +
                'Encryption is enabled but the secret key is missing.'
            );
        }
        if (!iv) {
            throw new CryptoConfigError(
                'VITE_CRYPTO_IV is not configured. ' +
                'Encryption is enabled but the IV is missing.'
            );
        }
    }

    return { secretKey, iv, enabled, useSignature };
}

/**
 * AES-256-CBC encryption
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
 * Generate SHA-256 signature
 */
function generateSignature(plaintext: string, timestamp: number, secretKey: string): string {
    const message = plaintext + timestamp.toString() + secretKey;
    return CryptoJS.SHA256(message).toString(CryptoJS.enc.Hex);
}

/**
 * Encrypt request data
 * @throws {CryptoConfigError} If encryption enabled but keys missing
 */
export function encryptRequest(data: unknown): EncryptedRequest | unknown {
    let config: CryptoConfig;

    try {
        config = getCryptoConfig();
    } catch (error) {
        if (error instanceof CryptoConfigError) {
            console.error('Encryption configuration error:', error.message);
            console.warn('Request will be sent without encryption.');
            return data;
        }
        throw error;
    }

    if (!config.enabled) {
        return data;
    }

    try {
        const plaintext = typeof data === 'string' ? data : JSON.stringify(data);
        const payload = encryptAES(plaintext, config.secretKey, config.iv);
        const timestamp = Date.now();
        const signature = config.useSignature
            ? generateSignature(plaintext, timestamp, config.secretKey)
            : '';

        return {
            encrypted: true,
            payload,
            timestamp,
            signature,
        };
    } catch (error) {
        console.error('Failed to encrypt request:', error);
        return data;
    }
}

/**
 * Determine if request should be encrypted
 */
export function shouldEncrypt(method: string, url: string): boolean {
    try {
        const config = getCryptoConfig();

        if (!config.enabled) {
            return false;
        }

        // Only encrypt POST, PUT, PATCH requests
        const encryptMethods = ['POST', 'PUT', 'PATCH'];
        if (!encryptMethods.includes(method.toUpperCase())) {
            return false;
        }

        // Exclude specific paths
        const excludePaths = [
            '/health',
            '/ping',
            '/auth/refresh',
        ];

        return !excludePaths.some(path => url.includes(path));
    } catch (error) {
        console.warn('Crypto configuration check failed, disabling encryption:', error);
        return false;
    }
}

/**
 * Check if encryption is properly configured
 */
export function isCryptoConfigured(): boolean {
    try {
        const config = getCryptoConfig();
        return config.enabled && !!config.secretKey && !!config.iv;
    } catch (error) {
        if (error instanceof CryptoConfigError) {
            console.warn('Crypto configuration check failed:', error.message);
            return false;
        }
        throw error;
    }
}
