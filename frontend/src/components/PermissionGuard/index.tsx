/**
 * 权限守卫组件
 * 用于控制按钮、菜单等UI元素的显示/隐藏
 */
import React, { type ReactNode, type ReactElement } from 'react';
import { usePermission } from '@/hooks/usePermission';

/**
 * PermissionGuard组件属性接口
 */
export interface PermissionGuardProps {
    /** 需要的权限码，可以是单个或多个 */
    permission: string | string[];
    /** 权限检查模式：any-任一满足，all-全部满足 */
    mode?: 'any' | 'all';
    /** 有权限时显示的内容 */
    children: ReactNode;
    /** 无权限时显示的内容（可选） */
    fallback?: ReactNode;
    /** 加载中时显示的内容（可选） */
    loading?: ReactNode;
    /** 是否在无权限时禁用而非隐藏（适用于按钮） */
    disabled?: boolean;
}

/**
 * PermissionGuard 权限守卫组件
 *
 * @description
 * 包裹需要权限控制的UI元素，根据用户权限决定是否渲染
 *
 * @example
 * ```tsx
 * // 基础用法 - 有权限则显示
 * <PermissionGuard permission="admin.games.create">
 *     <Button type="primary">创建游戏</Button>
 * </PermissionGuard>
 *
 * // 多个权限（任一满足）
 * <PermissionGuard permission={['admin.games.create', 'admin.games.update']}>
 *     <Button>操作</Button>
 * </PermissionGuard>
 *
 * // 多个权限（全部满足）
 * <PermissionGuard permission={['admin.games.read', 'admin.games.update']} mode="all">
 *     <Button>编辑</Button>
 * </PermissionGuard>
 *
 * // 带有fallback
 * <PermissionGuard
 *     permission="admin.games.delete"
 *     fallback={<span>无权限</span>}
 * >
 *     <Button danger>删除</Button>
 * </PermissionGuard>
 *
 * // 禁用模式（无权限时禁用按钮而非隐藏）
 * <PermissionGuard permission="admin.games.create" disabled>
 *     <Button type="primary">创建</Button>
 * </PermissionGuard>
 * ```
 */
export const PermissionGuard: React.FC<PermissionGuardProps> = ({
    permission,
    mode = 'any',
    children,
    fallback = null,
    loading: loadingContent = null,
    disabled = false,
}) => {
    const { hasPermission, loading } = usePermission(permission, mode);

    // 加载中
    if (loading && loadingContent) {
        return <>{loadingContent}</>;
    }

    // 有权限，直接渲染子组件
    if (hasPermission) {
        return <>{children}</>;
    }

    // 无权限，禁用模式
    if (disabled && React.isValidElement(children)) {
        return React.cloneElement(children as ReactElement<{ disabled?: boolean }>, {
            disabled: true,
        });
    }

    // 无权限，显示fallback或不渲染
    return <>{fallback}</>;
};

/**
 * 权限按钮组件属性接口
 */
export interface PermissionButtonProps extends Omit<PermissionGuardProps, 'disabled'> {
    /** 无权限时是否禁用按钮 */
    disableOnNoPermission?: boolean;
}

/**
 * PermissionButton 权限按钮包装器
 *
 * @description
 * 专门为按钮设计的权限包装器，支持禁用模式
 *
 * @example
 * ```tsx
 * <PermissionButton
 *     permission="admin.games.create"
 *     disableOnNoPermission
 * >
 *     <Button type="primary">创建游戏</Button>
 * </PermissionButton>
 * ```
 */
export const PermissionButton: React.FC<PermissionButtonProps> = ({
    disableOnNoPermission = false,
    ...props
}) => {
    return <PermissionGuard {...props} disabled={disableOnNoPermission} />;
};

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

export default PermissionGuard;
