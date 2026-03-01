import type { ReactNode } from 'react';
import { ErrorBoundary as ReactErrorBoundary } from 'react-error-boundary';
import type { FallbackProps } from 'react-error-boundary';
import { Button, Result, Typography, Space } from 'antd';
import { ReloadOutlined, BugOutlined, HomeOutlined } from '@ant-design/icons';
import { captureException, addBreadcrumb } from '@/utils/monitoring';

const { Paragraph, Text } = Typography;

function toError(error: unknown): Error {
  if (error instanceof Error) {
    return error;
  }

  return new Error(String(error));
}

/**
 * 错误回退 UI 组件
 */
function ErrorFallback({ error, resetErrorBoundary }: FallbackProps) {
  const isDev = import.meta.env.DEV;
  const normalizedError = toError(error);

  const handleGoHome = () => {
    window.location.href = '/admin';
  };

  const handleRefresh = () => {
    window.location.reload();
  };

  return (
    <div
      style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        minHeight: '100vh',
        padding: 24,
        background: '#f5f5f5',
      }}
    >
      <Result
        status="error"
        title="页面出错了"
        subTitle="抱歉，页面遇到了一些问题。请尝试刷新页面或返回首页。"
        extra={
          <Space>
            <Button type="primary" icon={<ReloadOutlined />} onClick={resetErrorBoundary}>
              重试
            </Button>
            <Button icon={<ReloadOutlined />} onClick={handleRefresh}>
              刷新页面
            </Button>
            <Button icon={<HomeOutlined />} onClick={handleGoHome}>
              返回首页
            </Button>
          </Space>
        }
      >
        {isDev && (
          <div
            style={{
              background: '#fff1f0',
              border: '1px solid #ffa39e',
              borderRadius: 8,
              padding: 16,
              marginTop: 16,
              textAlign: 'left',
            }}
          >
            <Paragraph>
              <BugOutlined style={{ color: '#ff4d4f', marginRight: 8 }} />
              <Text strong style={{ color: '#ff4d4f' }}>
                开发模式错误详情:
              </Text>
            </Paragraph>
            <Paragraph>
              <Text type="danger" code>
                {normalizedError.name}: {normalizedError.message}
              </Text>
            </Paragraph>
            {normalizedError.stack && (
              <Paragraph>
                <pre
                  style={{
                    fontSize: 12,
                    overflow: 'auto',
                    maxHeight: 200,
                    background: '#fff',
                    padding: 8,
                    borderRadius: 4,
                    margin: 0,
                  }}
                >
                  {normalizedError.stack}
                </pre>
              </Paragraph>
            )}
          </div>
        )}
      </Result>
    </div>
  );
}

/**
 * 错误日志记录
 */
function logError(error: Error, info: { componentStack?: string | null }) {
  // 开发环境打印到控制台
  if (import.meta.env.DEV) {
    console.group('🔴 ErrorBoundary 捕获到错误');
    console.error('Error:', error);
    console.error('Component Stack:', info.componentStack);
    console.groupEnd();
  }

  // 生产环境发送到错误监控服务 (Sentry)
  if (import.meta.env.PROD) {
    // Add breadcrumb for context
    addBreadcrumb('error-boundary', 'ErrorBoundary caught an error', {
      errorMessage: error.message,
      errorName: error.name,
    });

    // Capture exception with component stack context
    captureException(error, {
      componentStack: info.componentStack ?? undefined,
      errorBoundary: true,
    });
  }
}

interface Props {
  children: ReactNode;
  /** 自定义回退 UI */
  fallback?: ReactNode;
  /** 错误发生时的回调 */
  onError?: (error: Error, info: { componentStack?: string | null }) => void;
  /** 重置时的回调 */
  onReset?: () => void;
}

/**
 * 应用错误边界组件
 *
 * 用于捕获子组件树中的 JavaScript 错误，防止整个应用崩溃白屏。
 *
 * @example
 * ```tsx
 * <ErrorBoundary>
 *   <App />
 * </ErrorBoundary>
 * ```
 */
export function ErrorBoundary({ children, fallback, onError, onReset }: Props) {
  return (
    <ReactErrorBoundary
      FallbackComponent={fallback ? () => <>{fallback}</> : ErrorFallback}
      onError={(error, info) => {
        const normalizedError = toError(error);
        logError(normalizedError, info);
        onError?.(normalizedError, info);
      }}
      onReset={onReset}
    >
      {children}
    </ReactErrorBoundary>
  );
}

/**
 * 页面级错误边界
 * 用于包装单个页面，错误不会影响其他页面
 */
export function PageErrorBoundary({ children }: { children: ReactNode }) {
  return (
    <ErrorBoundary
      onReset={() => {
        // 页面级重置：可以清理页面状态
      }}
    >
      {children}
    </ErrorBoundary>
  );
}

export default ErrorBoundary;
