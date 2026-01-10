import type { ReactNode } from 'react';

/**
 * 用户角色类型
 * 与后端 model.Role 保持一致：user, player, admin
 * 注意：这里使用大写是为了路由守卫中的兼容性处理
 */
export type Role = 'USER' | 'PLAYER' | 'ADMIN';

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
