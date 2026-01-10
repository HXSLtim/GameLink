/**
 * 状态容器组件
 * 统一管理加载、空状态、错误状态
 */
import React from 'react';
import LoadingState from '../LoadingState';
import EmptyState from '../EmptyState';

export interface StateContainerProps {
    /** 是否加载中 */
    loading?: boolean;
    /** 数据数组 */
    data?: unknown[] | unknown;
    /** 错误信息 */
    error?: string | null;
    /** 空状态类型 */
    emptyType?: 'no-data' | 'no-search' | 'no-permission';
    /** 空状态标题 */
    emptyTitle?: string;
    /** 空状态描述 */
    emptyDescription?: string;
    /** 空状态操作按钮 */
    emptyActionText?: string;
    /** 空状态操作回调 */
    onEmptyAction?: () => void;
    /** 加载状态配置 */
    loadingConfig?: {
        card?: boolean;
        title?: string;
        rows?: number;
    };
    /** 子内容 */
    children: React.ReactNode;
}

/**
 * StateContainer 组件
 * 优化: 使用 React.memo 避免不必要的重新渲染
 * 适用场景: 状态容器组件，仅在状态变化时需要重新渲染
 */
const StateContainer: React.FC<StateContainerProps> = React.memo(({
    loading = false,
    data,
    error,
    emptyType = 'no-data',
    emptyTitle,
    emptyDescription,
    emptyActionText,
    onEmptyAction,
    loadingConfig,
    children,
}) => {
    // 错误状态
    if (error) {
        return (
            <EmptyState
                type="error"
                description={error}
                actionText="重新加载"
                onAction={onEmptyAction}
            />
        );
    }

    // 加载状态
    if (loading) {
        return <LoadingState loading={loading} {...loadingConfig} />;
    }

    // 空状态判断
    const isEmpty = Array.isArray(data) ? data.length === 0 : !data;

    if (isEmpty) {
        return (
            <EmptyState
                type={emptyType}
                title={emptyTitle}
                description={emptyDescription}
                actionText={emptyActionText}
                onAction={onEmptyAction}
            />
        );
    }

    // 正常状态
    return <>{children}</>;
});

export default StateContainer;
