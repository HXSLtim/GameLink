import type { ReactNode } from 'react';
import { ErrorBoundary as ReactErrorBoundary } from 'react-error-boundary';
import type { FallbackProps } from 'react-error-boundary';
import { Button, Result, Typography, Space } from 'antd';
import { ReloadOutlined, BugOutlined, HomeOutlined } from '@ant-design/icons';

const { Paragraph, Text } = Typography;

/**
 * 错误回退 UI 组件
 */
function ErrorFallback({ error, resetErrorBoundary }: FallbackProps) {
  const isDev = import.meta.env.DEV;

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
                {error.name}: {error.message}
              </Text>
            </Paragraph>
            {error.stack && (
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
                  {error.stack}
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

  // TODO: 生产环境可以发送到错误监控服务 (如 Sentry)
  // if (import.meta.env.PROD) {
  //   sendToErrorService({ error, componentStack: info.componentStack });
  // }
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
        logError(error, info);
        onError?.(error, info);
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
