import { describe, it, expect, vi, beforeEach } from 'vitest';
import { AppError, ErrorSeverity, errorHandler, handleApiError } from './errorHandler';

describe('AppError', () => {
  it('should create an AppError with default values', () => {
    const error = new AppError('Test error');

    expect(error.message).toBe('Test error');
    expect(error.code).toBe('UNKNOWN_ERROR');
    expect(error.severity).toBe(ErrorSeverity.ERROR);
    expect(error.timestamp).toBeInstanceOf(Date);
  });

  it('should create an AppError with custom values', () => {
    const context = { userId: '123' };
    const error = new AppError('Custom error', 'CUSTOM_CODE', ErrorSeverity.WARNING, context);

    expect(error.message).toBe('Custom error');
    expect(error.code).toBe('CUSTOM_CODE');
    expect(error.severity).toBe(ErrorSeverity.WARNING);
    expect(error.context).toEqual(context);
  });

  it('should have correct name', () => {
    const error = new AppError('Test');
    expect(error.name).toBe('AppError');
  });

  it('should be instance of Error', () => {
    const error = new AppError('Test');
    expect(error).toBeInstanceOf(Error);
    expect(error).toBeInstanceOf(AppError);
  });
});

describe('errorHandler', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should handle Error instances', () => {
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    const error = new Error('Test error');

    errorHandler.handle(error, false);

    expect(consoleSpy).toHaveBeenCalled();
    consoleSpy.mockRestore();
  });

  it('should handle AppError instances', () => {
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    const error = new AppError('Test error', 'TEST_ERROR', ErrorSeverity.WARNING);

    errorHandler.handle(error, false);

    expect(consoleSpy).toHaveBeenCalled();
    consoleSpy.mockRestore();
  });

  it('should handle async errors successfully', async () => {
    const successPromise = Promise.resolve('success');

    const [data, error] = await errorHandler.handleAsync(successPromise);

    expect(data).toBe('success');
    expect(error).toBeNull();
  });

  it('should handle async errors with failure', async () => {
    const failurePromise = Promise.reject(new Error('Failed'));
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

    const [data, error] = await errorHandler.handleAsync(failurePromise);

    expect(data).toBeNull();
    expect(error).toBeInstanceOf(Error);
    expect(error?.message).toBe('Failed');

    consoleSpy.mockRestore();
  });

  it('should handle async errors with custom message', async () => {
    const failurePromise = Promise.reject('Unknown error');
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

    const [data, error] = await errorHandler.handleAsync(failurePromise, 'Custom error');

    expect(data).toBeNull();
    expect(error).toBeInstanceOf(Error);
    // 当error是字符串时，normalizeError会直接使用字符串作为错误消息，而不是defaultMessage
    // 只有当error不是Error实例也不是字符串时，才会使用defaultMessage
    expect(error?.message).toBe('Unknown error');

    consoleSpy.mockRestore();
  });

  it('should use custom message when error is not string or Error', async () => {
    const failurePromise = Promise.reject({ unknown: 'object' });
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

    const [data, error] = await errorHandler.handleAsync(failurePromise, 'Custom error');

    expect(data).toBeNull();
    expect(error).toBeInstanceOf(Error);
    // 当error是对象时，会使用defaultMessage
    expect(error?.message).toBe('Custom error');

    consoleSpy.mockRestore();
  });
});

describe('handleApiError', () => {
  it('should handle Error instances', () => {
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    const error = new Error('API failed');

    handleApiError(error, '获取数据');

    expect(consoleSpy).toHaveBeenCalled();
    consoleSpy.mockRestore();
  });

  it('should handle string errors', () => {
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

    handleApiError('String error', '操作');

    expect(consoleSpy).toHaveBeenCalled();
    consoleSpy.mockRestore();
  });

  it('should handle unknown errors', () => {
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

    handleApiError({ unknown: 'object' }, '未知操作');

    expect(consoleSpy).toHaveBeenCalled();
    consoleSpy.mockRestore();
  });
});
