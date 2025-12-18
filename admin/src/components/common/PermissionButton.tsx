/**
 * 权限按钮组件
 * 根据用户权限自动控制按钮的显示/禁用状态
 */
import React from 'react';
import { Button, Tooltip } from 'antd';
import type { ButtonProps } from 'antd';
import { usePermission } from '@/hooks/usePermission';

export interface PermissionButtonProps extends ButtonProps {
    /** 所需权限码，可以是单个或多个 */
    permission: string | string[];
    /** 多个权限时的检查模式：any-任一满足，all-全部满足 */
    mode?: 'any' | 'all';
    /** 无权限时的处理方式：hide-隐藏，disable-禁用 */
    fallback?: 'hide' | 'disable';
    /** 无权限时的提示文字 */
    disabledTip?: string;
    /** 子元素 */
    children?: React.ReactNode;
}

/**
 * 权限按钮组件
 * 
 * @example
 * ```tsx
 * // 基础用法 - 无权限时隐藏
 * <PermissionButton permission="admin.users.create" type="primary">
 *     新建用户
 * </PermissionButton>
 * 
 * // 无权限时禁用并显示提示
 * <PermissionButton 
 *     permission="admin.users.delete" 
 *     fallback="disable"
 *     disabledTip="您没有删除权限"
 *     danger
 * >
 *     删除
 * </PermissionButton>
 * 
 * // 多个权限（任一满足）
 * <PermissionButton permission={['admin.users.update', 'admin.users.create']}>
 *     编辑
 * </PermissionButton>
 * 
 * // 多个权限（全部满足）
 * <PermissionButton permission={['admin.users.update', 'admin.users.delete']} mode="all">
 *     批量操作
 * </PermissionButton>
 * ```
 */
const PermissionButton: React.FC<PermissionButtonProps> = ({
    permission,
    mode = 'any',
    fallback = 'hide',
    disabledTip = '您没有此操作的权限',
    children,
    disabled,
    ...buttonProps
}) => {
    const { hasPermission, loading } = usePermission(permission, mode);

    // 加载中时显示加载状态
    if (loading) {
        return (
            <Button {...buttonProps} loading disabled>
                {children}
            </Button>
        );
    }

    // 无权限时的处理
    if (!hasPermission) {
        if (fallback === 'hide') {
            return null;
        }
        // fallback === 'disable'
        return (
            <Tooltip title={disabledTip}>
                <Button {...buttonProps} disabled>
                    {children}
                </Button>
            </Tooltip>
        );
    }

    // 有权限，正常渲染
    return (
        <Button {...buttonProps} disabled={disabled}>
            {children}
        </Button>
    );
};

export default PermissionButton;
