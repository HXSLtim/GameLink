import { useAuthStore } from '@/stores';

interface PermissionGateProps {
    permission: string;
    children: React.ReactNode;
    fallback?: React.ReactNode;
}

export function PermissionGate({ permission, children, fallback = null }: PermissionGateProps) {
    const { permissions } = useAuthStore();

    const hasPermission = permissions?.includes(permission);

    if (!hasPermission) {
        return <>{fallback}</>;
    }

    return <>{children}</>;
}
