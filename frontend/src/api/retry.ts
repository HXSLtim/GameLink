/**
 * HTTP 请求重试工具
 * 提供自动重试机制，处理临时性网络错误
 */

export interface RetryOptions {
  /**
   * 最大重试次数
   * @default 3
   */
  maxRetries?: number;

  /**
   * 初始延迟时间（毫秒）
   * @default 1000
   */
  initialDelay?: number;

  /**
   * 延迟倍增因子（指数退避）
   * @default 2
   */
  backoffFactor?: number;

  /**
   * 最大延迟时间（毫秒）
   * @default 10000
   */
  maxDelay?: number;

  /**
   * 判断错误是否应该重试
   * @param error 错误对象
   * @param attempt 当前尝试次数
   * @returns 是否应该重试
   */
  shouldRetry?: (error: unknown, attempt: number) => boolean;

  /**
   * 重试前的回调
   * @param error 错误对象
   * @param attempt 当前尝试次数
   * @param delay 延迟时间
   */
  onRetry?: (error: unknown, attempt: number, delay: number) => void;
}

const DEFAULT_RETRY_OPTIONS: Required<RetryOptions> = {
  maxRetries: 3,
  initialDelay: 1000,
  backoffFactor: 2,
  maxDelay: 10000,
  shouldRetry: (error: unknown) => {
    // 默认重试策略：仅重试网络错误
    if (error instanceof Error) {
      // 网络错误
      if (error.message.includes('NetworkError') || error.message.includes('Failed to fetch')) {
        return true;
      }
      // 超时错误
      if (error.message.includes('timeout')) {
        return true;
      }
    }

    // HTTP 状态码 5xx（服务器错误）应该重试
    if (typeof error === 'object' && error !== null && 'status' in error) {
      const status = (error as { status: number }).status;
      return status >= 500 && status < 600;
    }

    return false;
  },
  onRetry: () => {
    // 默认不做任何处理
  },
};

/**
 * 计算重试延迟（指数退避）
 * @param attempt 当前尝试次数
 * @param options 重试选项
 * @returns 延迟时间（毫秒）
 */
function calculateDelay(attempt: number, options: Required<RetryOptions>): number {
  const delay = options.initialDelay * Math.pow(options.backoffFactor, attempt - 1);
  return Math.min(delay, options.maxDelay);
}

/**
 * 等待指定时间
 * @param ms 毫秒数
 */
function wait(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

/**
 * 执行带重试的异步操作
 * @param fn 异步函数
 * @param options 重试选项
 * @returns Promise<T>
 */
export async function retryAsync<T>(
  fn: () => Promise<T>,
  options: RetryOptions = {},
): Promise<T> {
  const opts: Required<RetryOptions> = {
    ...DEFAULT_RETRY_OPTIONS,
    ...options,
  };

  let lastError: unknown;
  let attempt = 0;

  while (attempt <= opts.maxRetries) {
    attempt++;

    try {
      return await fn();
    } catch (error) {
      lastError = error;

      // 达到最大重试次数
      if (attempt > opts.maxRetries) {
        break;
      }

      // 判断是否应该重试
      if (!opts.shouldRetry(error, attempt)) {
        throw error;
      }

      // 计算延迟
      const delay = calculateDelay(attempt, opts);

      // 触发重试回调
      opts.onRetry(error, attempt, delay);

      // 等待后重试
      await wait(delay);
    }
  }

  // 所有重试都失败，抛出最后的错误
  throw lastError;
}

/**
 * 创建带重试的 fetch 包装函数
 * @param fetchFn 原始 fetch 函数
 * @param options 重试选项
 * @returns 包装后的 fetch 函数
 */
export function createRetryFetch(
  fetchFn: typeof fetch = fetch,
  options: RetryOptions = {},
): typeof fetch {
  return async (input: RequestInfo | URL, init?: RequestInit): Promise<Response> => {
    return retryAsync(
      () => fetchFn(input, init),
      {
        ...options,
        shouldRetry: (error: unknown, attempt: number) => {
          // 自定义重试逻辑
          if (options.shouldRetry) {
            return options.shouldRetry(error, attempt);
          }

          // 网络错误
          if (error instanceof TypeError && error.message.includes('Failed to fetch')) {
            return true;
          }

          // Response 对象
          if (error instanceof Response) {
            // 5xx 错误重试
            return error.status >= 500 && error.status < 600;
          }

          return false;
        },
        onRetry: (error: unknown, attempt: number, delay: number) => {
          console.warn(
            `[Retry] Attempt ${attempt}/${options.maxRetries || 3} failed, retrying in ${delay}ms...`,
            error,
          );
          options.onRetry?.(error, attempt, delay);
        },
      },
    );
  };
}

