/**
 * withPermission HOC
 * 高阶组件，为任意组件添加权限控制
 */
import React, { type ReactNode } from 'react';
import { PermissionGuard } from './index';

/**
 * withPermission HOC
 *
 * @description
 * 高阶组件，为任意组件添加权限控制
 *
 * @example
 * ```tsx
 * const ProtectedButton = withPermission(Button, 'admin.games.create');
 *
 * // 使用
 * <ProtectedButton type="primary">创建</ProtectedButton>
 * ```
 */
export function withPermission<P extends object>(
  WrappedComponent: React.ComponentType<P>,
  permission: string | string[],
  options: {
    mode?: 'any' | 'all';
    fallback?: ReactNode;
    disabled?: boolean;
  } = {}
): React.FC<P> {
  const { mode = 'any', fallback = null, disabled = false } = options;

  const WithPermissionComponent: React.FC<P> = (props) => {
    return (
      <PermissionGuard
        permission={permission}
        mode={mode}
        fallback={fallback}
        disabled={disabled}
      >
        <WrappedComponent {...props} />
      </PermissionGuard>
    );
  };

  WithPermissionComponent.displayName = `withPermission(${WrappedComponent.displayName || WrappedComponent.name || 'Component'})`;

  return WithPermissionComponent;
}
