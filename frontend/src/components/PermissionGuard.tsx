/**
 * 权限守卫组件
 * 提供按钮级权限控制功能
 * Requirements: 3.1, 3.2, 3.3, 3.4
 */
import React from 'react';
import { Button, Tooltip } from 'antd';
import type { ButtonProps } from 'antd';
import { useAdmin } from '@/context/AdminContext';

/**
 * 权限检查模式
 */
export type PermissionCheckMode = 'any' | 'all';

/**
 * PermissionGuard 组件属性
 * Requirements: 3.1, 3.2, 3.3, 3.4
 */
export interface PermissionGuardProps {
    /** 需要检查的权限码，支持单个或多个 */
    permission: string | string[];
    /** 权限检查模式：any-任一满足，all-全部满足 */
    mode?: PermissionCheckMode;
    /** 子元素 */
    children: React.ReactNode;
    /** 无权限时的回退内容 */
    fallback?: React.ReactNode;
    /** 加载中时显示的内容 */
    loading?: React.ReactNode;
    /** 是否使用禁用模式（无权限时显示禁用状态而非隐藏） */
    disabled?: boolean;
    /** 禁用模式下的提示文本 */
    tooltip?: string;
}

/**
 * PermissionGuard 组件
 * 根据用户权限决定子元素的显示/隐藏或启用/禁用
 *
 * @example
 * ```tsx
 * // 基本用法 - 无权限时隐藏
 * <PermissionGuard permission="admin.users.create">
 *     <Button type="primary">创建用户</Button>
 * </PermissionGuard>
 *
 * // 多权限检查（任一满足）
 * <PermissionGuard permission={['admin.users.update', 'admin.users.delete']} mode="any">
 *     <Button>操作</Button>
 * </PermissionGuard>
 *
 * // 多权限检查（全部满足）
 * <PermissionGuard permission={['admin.users.read', 'admin.users.update']} mode="all">
 *     <Button>编辑</Button>
 * </PermissionGuard>
 *
 * // 禁用模式 - 无权限时显示禁用按钮
 * <PermissionGuard permission="admin.users.delete" disabled tooltip="您没有删除权限">
 *     <Button danger>删除</Button>
 * </PermissionGuard>
 *
 * // 自定义加载状态
 * <PermissionGuard permission="admin.users.create" loading={<Spin size="small" />}>
 *     <Button>创建</Button>
 * </PermissionGuard>
 *
 * // 自定义回退内容
 * <PermissionGuard permission="admin.users.create" fallback={<span>无权限</span>}>
 *     <Button>创建</Button>
 * </PermissionGuard>
 * ```
 */
export const PermissionGuard: React.FC<PermissionGuardProps> = ({
    permission,
    mode = 'any',
    children,
    fallback = null,
    loading: loadingContent,
    disabled: disabledMode = false,
    tooltip = '您没有此操作的权限',
}) => {
    const { loading, hasPermission, isSuperAdmin } = useAdmin();

    // 处理加载状态 - 避免权限闪烁
    // Requirements: 3.3
    if (loading) {
        // 如果提供了自定义加载内容，显示它
        if (loadingContent !== undefined) {
            return <>{loadingContent}</>;
        }
        // 默认情况下，加载时不显示任何内容，避免闪烁
        return null;
    }

    // 检查权限
    const permissionList = Array.isArray(permission) ? permission : [permission];
    
    // 如果没有指定有效的权限要求，默认允许
    const hasValidPermission = permissionList.some(p => p && p.length > 0);
    if (!hasValidPermission) {
        return <>{children}</>;
    }

    // 超级管理员拥有所有权限
    // Requirements: 3.5
    if (isSuperAdmin) {
        return <>{children}</>;
    }

    // 使用 context 提供的 hasPermission 方法进行检查
    const hasRequiredPermission = hasPermission(permissionList, mode);

    // 有权限时直接显示子元素
    if (hasRequiredPermission) {
        return <>{children}</>;
    }

    // 无权限时的处理
    // Requirements: 3.2
    if (disabledMode) {
        // 禁用模式：克隆子元素并添加 disabled 属性
        return (
            <Tooltip title={tooltip}>
                {React.Children.map(children, (child) => {
                    if (React.isValidElement(child)) {
                        // 克隆元素并添加 disabled 属性
                        return React.cloneElement(child as React.ReactElement<{ disabled?: boolean }>, {
                            disabled: true,
                        });
                    }
                    return child;
                })}
            </Tooltip>
        );
    }

    // 默认：显示回退内容（默认为 null，即隐藏）
    return <>{fallback}</>;
};


/**
 * PermissionButton 组件属性
 * 带权限控制的按钮组件
 */
export interface PermissionButtonProps extends ButtonProps {
    /** 需要检查的权限码 */
    permission: string | string[];
    /** 权限检查模式 */
    mode?: PermissionCheckMode;
    /** 无权限时的提示文本 */
    tooltip?: string;
    /** 是否在无权限时隐藏（默认显示禁用状态） */
    hideOnNoPermission?: boolean;
}

/**
 * PermissionButton 组件
 * 带权限控制的按钮，无权限时显示禁用状态或隐藏
 *
 * @example
 * ```tsx
 * // 无权限时显示禁用按钮
 * <PermissionButton permission="admin.users.delete" tooltip="没有删除权限">
 *     删除
 * </PermissionButton>
 *
 * // 无权限时隐藏按钮
 * <PermissionButton permission="admin.users.delete" hideOnNoPermission>
 *     删除
 * </PermissionButton>
 * ```
 */
export const PermissionButton: React.FC<PermissionButtonProps> = ({
    permission,
    mode = 'any',
    tooltip = '您没有此操作的权限',
    hideOnNoPermission = false,
    children,
    disabled,
    ...props
}) => {
    const { loading, hasPermission, isSuperAdmin } = useAdmin();

    // 处理加载状态
    if (loading) {
        return (
            <Button disabled {...props}>
                {children}
            </Button>
        );
    }

    // 检查权限
    const permissionList = Array.isArray(permission) ? permission : [permission];
    const hasRequiredPermission = isSuperAdmin || hasPermission(permissionList, mode);

    if (!hasRequiredPermission) {
        // 无权限时隐藏
        if (hideOnNoPermission) {
            return null;
        }
        // 无权限时显示禁用按钮
        return (
            <Tooltip title={tooltip}>
                <Button disabled {...props}>
                    {children}
                </Button>
            </Tooltip>
        );
    }

    // 有权限时正常显示
    return (
        <Button disabled={disabled} {...props}>
            {children}
        </Button>
    );
};

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

export default PermissionGuard;
