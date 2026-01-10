/**
 * 权限容器组件
 * 根据用户权限控制子元素的显示
 */
import React from 'react';
import { Tooltip } from 'antd';
import { usePermission } from '@/hooks/usePermission';

export interface PermissionWrapperProps {
    /** 所需权限码，可以是单个或多个 */
    permission: string | string[];
    /** 多个权限时的检查模式：any-任一满足，all-全部满足 */
    mode?: 'any' | 'all';
    /** 无权限时的处理方式：hide-隐藏，disable-禁用样式 */
    fallback?: 'hide' | 'disable';
    /** 无权限时的提示文字（仅 fallback='disable' 时有效） */
    disabledTip?: string;
    /** 子元素 */
    children: React.ReactNode;
    /** 无权限时显示的替代内容 */
    fallbackContent?: React.ReactNode;
}

/**
 * 权限容器组件
 * 优化: 使用 React.memo 避免不必要的重新渲染
 * 适用场景: 权限容器组件，在权限和子元素相同时不需要重新渲染
 *
 * @example
 * ```tsx
 * // 基础用法 - 无权限时隐藏
 * <PermissionWrapper permission="admin.users.create">
 *     <Button type="primary">新建用户</Button>
 * </PermissionWrapper>
 *
 * // 无权限时显示禁用样式
 * <PermissionWrapper permission="admin.users.delete" fallback="disable">
 *     <Button danger>删除</Button>
 * </PermissionWrapper>
 *
 * // 无权限时显示替代内容
 * <PermissionWrapper
 *     permission="admin.reports.view"
 *     fallbackContent={<span>暂无权限查看</span>}
 * >
 *     <ReportChart />
 * </PermissionWrapper>
 * ```
 */
const PermissionWrapper: React.FC<PermissionWrapperProps> = React.memo(({
    permission,
    mode = 'any',
    fallback = 'hide',
    disabledTip = '您没有此操作的权限',
    children,
    fallbackContent,
}) => {
    const { hasPermission, loading } = usePermission(permission, mode);

    // 加载中时不显示
    if (loading) {
        return null;
    }

    // 有权限，正常渲染
    if (hasPermission) {
        return <>{children}</>;
    }

    // 无权限时的处理
    if (fallbackContent) {
        return <>{fallbackContent}</>;
    }

    if (fallback === 'hide') {
        return null;
    }

    // fallback === 'disable' - 添加禁用样式
    return (
        <Tooltip title={disabledTip}>
            <div style={{ opacity: 0.5, pointerEvents: 'none', cursor: 'not-allowed' }}>
                {children}
            </div>
        </Tooltip>
    );
});

export default PermissionWrapper;
