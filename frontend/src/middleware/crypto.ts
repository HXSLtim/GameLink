import { AxiosRequestConfig, AxiosResponse, InternalAxiosRequestConfig } from 'axios';
import { CryptoUtil } from '../utils/crypto';

/**
 * 加密中间件配置
 */
interface CryptoMiddlewareConfig {
  // 是否启用加密
  enabled: boolean;
  // 需要加密的请求方法
  methods: string[];
  // 不需要加密的路径（白名单）
  excludePaths: string[];
  // 需要加密的敏感字段（部分加密模式）
  sensitiveFields?: string[];
  // 是否使用签名验证
  useSignature: boolean;
}

/**
 * 默认配置
 */
const DEFAULT_CONFIG: CryptoMiddlewareConfig = {
  enabled: import.meta.env.VITE_CRYPTO_ENABLED !== 'false',
  methods: ['POST', 'PUT', 'PATCH'],
  excludePaths: ['/api/v1/health', '/api/v1/ping'],
  sensitiveFields: ['password', 'phone', 'email', 'id_card'],
  useSignature: true,
};

/**
 * 加密中间件类
 */
export class CryptoMiddleware {
  private config: CryptoMiddlewareConfig;

  constructor(config?: Partial<CryptoMiddlewareConfig>) {
    this.config = { ...DEFAULT_CONFIG, ...config };
  }

  /**
   * 检查请求是否需要加密
   */
  private shouldEncrypt(config: InternalAxiosRequestConfig): boolean {
    if (!this.config.enabled) {
      return false;
    }

    // 检查请求方法
    const method = config.method?.toUpperCase() || 'GET';
    if (!this.config.methods.includes(method)) {
      return false;
    }

    // 检查路径白名单
    const url = config.url || '';
    if (this.config.excludePaths.some((path) => url.includes(path))) {
      return false;
    }

    return true;
  }

  /**
   * 请求拦截器：加密请求数据
   */
  requestInterceptor = (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
    if (!this.shouldEncrypt(config)) {
      return config;
    }

    try {
      const timestamp = Date.now();

      // 模式 1: 全量加密（推荐）
      if (config.data) {
        const encryptedData = CryptoUtil.encrypt(config.data);

        // 生成签名
        let signature = '';
        if (this.config.useSignature) {
          signature = CryptoUtil.generateSignature(config.data, timestamp);
        }

        // 包装加密数据
        config.data = {
          encrypted: true,
          payload: encryptedData,
          timestamp,
          signature,
        };

        console.log('🔒 请求数据已加密:', {
          url: config.url,
          method: config.method,
          timestamp,
        });
      }

      // 模式 2: 部分字段加密（可选）
      // if (config.data && this.config.sensitiveFields) {
      //   config.data = CryptoUtil.encryptFields(
      //     config.data,
      //     this.config.sensitiveFields,
      //   );
      // }

      return config;
    } catch (error) {
      console.error('❌ 请求加密失败:', error);
      throw error;
    }
  };

  /**
   * 响应拦截器：解密响应数据
   */
  responseInterceptor = (response: AxiosResponse): AxiosResponse => {
    if (!this.config.enabled) {
      return response;
    }

    try {
      const responseData = response.data;

      // 检查响应是否加密
      if (responseData && typeof responseData === 'object' && responseData.encrypted) {
        console.log('🔓 响应数据已加密，开始解密...', {
          url: response.config.url,
        });

        // 解密数据
        const decryptedData = CryptoUtil.decrypt(responseData.payload);

        // 验证签名（如果有）
        if (this.config.useSignature && responseData.signature && responseData.timestamp) {
          const expectedSignature = CryptoUtil.generateSignature(
            decryptedData,
            responseData.timestamp,
          );

          if (expectedSignature !== responseData.signature) {
            console.error('❌ 签名验证失败！数据可能被篡改');
            throw new Error('数据签名验证失败');
          }

          console.log('✅ 签名验证通过');
        }

        // 替换为解密后的数据
        response.data = decryptedData;

        console.log('✅ 响应数据解密成功');
      }

      return response;
    } catch (error) {
      console.error('❌ 响应解密失败:', error);
      throw error;
    }
  };

  /**
   * 获取配置
   */
  getConfig(): CryptoMiddlewareConfig {
    return { ...this.config };
  }

  /**
   * 更新配置
   */
  updateConfig(config: Partial<CryptoMiddlewareConfig>): void {
    this.config = { ...this.config, ...config };
  }

  /**
   * 启用加密
   */
  enable(): void {
    this.config.enabled = true;
  }

  /**
   * 禁用加密
   */
  disable(): void {
    this.config.enabled = false;
  }

  /**
   * 检查是否启用
   */
  isEnabled(): boolean {
    return this.config.enabled;
  }
}

/**
 * 创建加密中间件实例
 */
export const cryptoMiddleware = new CryptoMiddleware({
  enabled: import.meta.env.VITE_CRYPTO_ENABLED !== 'false',
  methods: ['POST', 'PUT', 'PATCH'],
  excludePaths: [
    '/api/v1/health',
    '/api/v1/ping',
    '/api/v1/auth/refresh', // Token 刷新不加密
  ],
  useSignature: true,
});

/**
 * 导出便捷方法
 */
export const { requestInterceptor, responseInterceptor } = cryptoMiddleware;

