import React from 'react';
import { Button, Tooltip } from 'antd';
import type { ButtonProps } from 'antd';
import { useAdmin } from '@/context/AdminContext';

export interface PermissionGuardProps {
    permission: string;
    children: React.ReactNode;
    fallback?: React.ReactNode;
}

export const PermissionGuard: React.FC<PermissionGuardProps> = ({
    permission,
    children,
    fallback = null
}) => {
    const { permissions } = useAdmin();

    // Check if user has the required permission
    // You might want to handle '*' (superuser) logic here if applicable
    const hasPermission = permissions.includes(permission) || permissions.includes('*');

    if (!hasPermission) {
        return <>{fallback}</>;
    }

    return <>{children}</>;
};

export interface PermissionButtonProps extends ButtonProps {
    permission: string;
    tooltip?: string;
}

export const PermissionButton: React.FC<PermissionButtonProps> = ({
    permission,
    tooltip,
    children,
    disabled,
    ...props
}) => {
    const { permissions } = useAdmin();
    const hasPermission = permissions.includes(permission) || permissions.includes('*');

    if (!hasPermission) {
        if (tooltip) {
            return (
                <Tooltip title={tooltip}>
                    <Button disabled {...props}>
                        {children}
                    </Button>
                </Tooltip>
            );
        }
        return null;
    }

    return (
        <Button {...props}>
            {children}
        </Button>
    );
};

export const withPermission = <P extends object>(
    WrappedComponent: React.ComponentType<P>,
    permission: string,
    FallbackComponent: React.ComponentType<P> | null = null
) => {
    return (props: P) => (
        <PermissionGuard permission={permission} fallback={FallbackComponent ? <FallbackComponent {...props} /> : null}>
            <WrappedComponent {...props} />
        </PermissionGuard>
    );
};
