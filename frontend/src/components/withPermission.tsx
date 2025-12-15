/**
 * withPermission 高阶组件
 * 为组件添加权限控制
 * 分离到单独文件以支持 React Fast Refresh
 */
import React from 'react';
import { PermissionGuard, type PermissionCheckMode } from './PermissionGuard';

/**
 * withPermission 高阶组件
 * 为组件添加权限控制
 *
 * @example
 * ```tsx
 * const ProtectedButton = withPermission(Button, 'admin.users.create');
 *
 * // 使用
 * <ProtectedButton type="primary">创建用户</ProtectedButton>
 * ```
 */
export const withPermission = <P extends object>(
    WrappedComponent: React.ComponentType<P>,
    permission: string | string[],
    options?: {
        mode?: PermissionCheckMode;
        fallback?: React.ComponentType<P> | null;
        loading?: React.ReactNode;
    }
) => {
    const { mode = 'any', fallback: FallbackComponent = null, loading: loadingContent } = options || {};

    const WithPermissionComponent = (props: P) => (
        <PermissionGuard
            permission={permission}
            mode={mode}
            fallback={FallbackComponent ? <FallbackComponent {...props} /> : null}
            loading={loadingContent}
        >
            <WrappedComponent {...props} />
        </PermissionGuard>
    );

    // 设置显示名称便于调试
    const wrappedName = WrappedComponent.displayName || WrappedComponent.name || 'Component';
    WithPermissionComponent.displayName = `withPermission(${wrappedName})`;

    return WithPermissionComponent;
};

export default withPermission;
