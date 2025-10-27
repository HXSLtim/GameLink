import { Component, ReactNode } from 'react';
import { Button, Result, Typography } from '@arco-design/web-react';
import { IconRefresh, IconHome } from '@arco-design/web-react/icon';
import { errorHandler, AppError, ErrorSeverity } from '../../utils/errorHandler';
import styles from './ErrorBoundary.module.less';

interface ErrorBoundaryProps {
  /** Child components */
  children: ReactNode;
  /** Custom fallback UI */
  fallback?: (error: Error, resetError: () => void) => ReactNode;
  /** Callback when error occurs */
  onError?: (error: Error, errorInfo: React.ErrorInfo) => void;
}

interface ErrorBoundaryState {
  hasError: boolean;
  error: Error | null;
}

/**
 * Error Boundary component
 *
 * @component
 * @description Catches JavaScript errors anywhere in the child component tree
 */
export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: ErrorBoundaryProps) {
    super(props);
    this.state = {
      hasError: false,
      error: null,
    };
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return {
      hasError: true,
      error,
    };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo): void {
    // Log error to error handler
    const appError = new AppError(
      error.message,
      'REACT_ERROR_BOUNDARY',
      ErrorSeverity.CRITICAL,
      {
        componentStack: errorInfo.componentStack,
      },
    );

    errorHandler.handle(appError, false);

    // Call custom error handler if provided
    if (this.props.onError) {
      this.props.onError(error, errorInfo);
    }
  }

  handleReset = (): void => {
    this.setState({
      hasError: false,
      error: null,
    });
  };

  handleGoHome = (): void => {
    window.location.href = '/';
  };

  render(): ReactNode {
    const { hasError, error } = this.state;
    const { children, fallback } = this.props;

    if (hasError && error) {
      // Use custom fallback if provided
      if (fallback) {
        return fallback(error, this.handleReset);
      }

      // Default error UI
      return (
        <div className={styles.errorBoundary}>
          <Result
            status="error"
            title="哎呀，出错了！"
            subTitle={
              <div className={styles.errorDetails}>
                <Typography.Paragraph>
                  应用程序遇到了一个意外错误。我们已经记录了这个问题。
                </Typography.Paragraph>
                <Typography.Text type="secondary" className={styles.errorMessage}>
                  错误信息: {error.message}
                </Typography.Text>
              </div>
            }
            extra={
              <div className={styles.actions}>
                <Button
                  type="primary"
                  size="large"
                  icon={<IconRefresh />}
                  onClick={this.handleReset}
                >
                  重试
                </Button>
                <Button size="large" icon={<IconHome />} onClick={this.handleGoHome}>
                  返回首页
                </Button>
              </div>
            }
          />
        </div>
      );
    }

    return children;
  }
}



