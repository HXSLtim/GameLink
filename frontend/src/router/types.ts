import type { ReactNode } from 'react';

export type Role = 'USER' | 'COMPANION' | 'ADMIN' | 'CS' | 'FINANCE';

export interface RouteConfig {
    path?: string;
    index?: boolean;
    element?: ReactNode;
    children?: RouteConfig[];
    meta?: {
        title: string;
        roles?: Role[]; // Allowed roles. If undefined, accessible by everyone (or authenticated users depending on context)
        requiresAuth?: boolean;
        permission?: string; // Required permission code for this route
        icon?: string;
        hideInMenu?: boolean;
    };
}
