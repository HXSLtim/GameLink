import apiClient from './client';

export interface Pagination {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
    has_next: boolean;
    has_prev: boolean;
}

export interface ApiResponse<T> {
    success: boolean;
    code: number;
    message: string;
    data: T;
    pagination?: Pagination;
}

export interface User {
    id: number;
    name: string;
    email: string;
    phone: string;
    avatarUrl?: string;
    role: 'user' | 'player' | 'admin';
    status: 'active' | 'banned' | 'suspended';
    lastLoginAt?: string;
    createdAt: string;
    updatedAt?: string;
}

export interface CreateUserDto {
    name: string;
    email: string;
    phone: string;
    password: string;
    avatarUrl?: string;
    role: 'user' | 'player' | 'admin';
    status: 'active' | 'banned' | 'suspended';
}

export interface UpdateUserDto {
    name: string;
    email: string;
    phone: string;
    avatarUrl?: string;
    role: 'user' | 'player' | 'admin';
    status: 'active' | 'banned' | 'suspended';
    password?: string;
}

export interface UserQueryParams {
    page?: number;
    page_size?: number;
    keyword?: string;
    role?: string[];
    status?: string[];
    date_from?: string;
    date_to?: string;
}

export interface Game {
    id: string;
    name: string;
    icon: string;
    category: string;
}

// Service Item Interfaces
export interface ServiceItem {
    id: number;
    name: string;
    gameId: number;
    gameName?: string;
    description: string;
    price: number;
    duration: number;
    category: string;
    tags: string[];
    icon: string;
    status: 'active' | 'inactive';
    sortOrder: number;
    createdAt: string;
    updatedAt: string;
    createdBy: number;
    updatedBy: number;
}

export interface CreateServiceItemDto {
    name: string;
    gameId: number;
    description: string;
    price: number;
    duration: number;
    category: string;
    tags?: string[];
    icon?: string;
    sortOrder?: number;
}

export interface UpdateServiceItemDto extends Partial<CreateServiceItemDto> {
    status?: 'active' | 'inactive';
}

export interface ServiceItemQueryParams {
    page?: number;
    page_size?: number;
    game_id?: number;
    category?: string;
    status?: string;
    keyword?: string;
    sort_by?: string;
    sort_order?: 'asc' | 'desc';
}

export interface Menu {
    id: number;
    name: string;
    path: string;
    component: string;
    parentId: number | null;
    order: number;
    visible: boolean;
    type: 'menu' | 'page' | 'button';
    permission: string;
    icon?: string;
    redirect?: string;
    description?: string;
    children?: Menu[];
}

export interface CreateMenuDto {
    name: string;
    path: string;
    component: string;
    parentId?: number | null;
    order?: number;
    visible?: boolean;
    type?: 'menu' | 'page' | 'button';
    permission?: string;
    icon?: string;
    redirect?: string;
    description?: string;
    children?: Menu[];
}

export interface UpdateMenuDto extends Partial<CreateMenuDto> { }

export interface Permission {
    id: number;
    name?: string;
    code: string;
    group: string;
    method: string;
    path: string;
    description?: string;
    createdAt: string;
    updatedAt: string;
}

export const adminApi = {
    // User Management
    getUsers: (params?: UserQueryParams) => apiClient.get<ApiResponse<User[]>>('/admin/users', { params }),
    getUser: (id: number) => apiClient.get<ApiResponse<User>>(`/admin/users/${id}`),
    createUser: (data: CreateUserDto) => apiClient.post<ApiResponse<User>>('/admin/users', data),
    updateUser: (id: number, data: UpdateUserDto) => apiClient.put<ApiResponse<User>>(`/admin/users/${id}`, data),
    deleteUser: (id: number) => apiClient.delete<ApiResponse<void>>(`/admin/users/${id}`),
    batchDeleteUsers: (ids: number[]) => apiClient.post<ApiResponse<void>>('/admin/users/batch-delete', { ids }),
    updateUserStatus: (id: number, status: string) => apiClient.put<ApiResponse<User>>(`/admin/users/${id}/status`, { status }),
    updateUserRole: (id: number, role: string) => apiClient.put<ApiResponse<User>>(`/admin/users/${id}/role`, { role }),
    getUserOrders: (id: number, params?: any) => apiClient.get<ApiResponse<any[]>>(`/admin/users/${id}/orders`, { params }),
    getUserLogs: (id: number, params?: any) => apiClient.get<ApiResponse<any[]>>(`/admin/users/${id}/logs`, { params }),

    // Game Management
    getGames: (params?: { status?: string; page_size?: number }) => apiClient.get('/admin/games', { params }),
    createGame: (data: any) => apiClient.post('/admin/games', data),
    updateGame: (id: string, data: any) => apiClient.put(`/admin/games/${id}`, data),
    deleteGame: (id: string) => apiClient.delete(`/admin/games/${id}`),

    // Service Item Management
    getServiceItems: (params?: ServiceItemQueryParams) => apiClient.get<ApiResponse<ServiceItem[]>>('/admin/service-items', { params }),
    getServiceItem: (id: number) => apiClient.get<ApiResponse<ServiceItem>>(`/admin/service-items/${id}`),
    createServiceItem: (data: CreateServiceItemDto) => apiClient.post<ApiResponse<ServiceItem>>('/admin/service-items', data),
    updateServiceItem: (id: number, data: UpdateServiceItemDto) => apiClient.put<ApiResponse<ServiceItem>>(`/admin/service-items/${id}`, data),
    deleteServiceItem: (id: number) => apiClient.delete<ApiResponse<void>>(`/admin/service-items/${id}`),
    batchUpdateServiceItemStatus: (ids: number[], status: 'active' | 'inactive') => apiClient.put<ApiResponse<void>>('/admin/service-items/batch-status', { ids, status }),

    // Menu Management
    getMenus: (params?: { parentId?: number; page?: number; page_size?: number }) => apiClient.get<ApiResponse<Menu[]>>('/admin/menus', { params }),
    getMenu: (id: number) => apiClient.get<ApiResponse<Menu>>(`/admin/menus/${id}`),
    createMenu: (data: CreateMenuDto) => apiClient.post<ApiResponse<Menu>>('/admin/menus', data),
    updateMenu: (id: number, data: UpdateMenuDto) => apiClient.put<ApiResponse<Menu>>(`/admin/menus/${id}`, data),
    deleteMenu: (id: number) => apiClient.delete<ApiResponse<void>>(`/admin/menus/${id}`),

    // Permissions
    getPermissions: (params?: any) => apiClient.get<ApiResponse<{ items: Permission[], totalCount: number, page: number, pageSize: number }>>('/admin/permissions', { params }),
    getPermissionGroups: () => apiClient.get<ApiResponse<string[]>>('/admin/permissions/groups'),
    getPermission: (id: number) => apiClient.get<ApiResponse<Permission>>(`/admin/permissions/${id}`),
    createPermission: (data: any) => apiClient.post<ApiResponse<Permission>>('/admin/permissions', data),
    updatePermission: (id: number, data: any) => apiClient.put<ApiResponse<Permission>>(`/admin/permissions/${id}`, data),
    deletePermission: (id: number) => apiClient.delete<ApiResponse<void>>(`/admin/permissions/${id}`),
    getMyPermissions: () => apiClient.get<ApiResponse<string[]>>('/admin/permissions/me'),

    // Role Management
    getRoles: (params?: any) => apiClient.get('/admin/roles', { params }),
    getRole: (id: number) => apiClient.get(`/admin/roles/${id}`),
    createRole: (data: any) => apiClient.post('/admin/roles', data),
    updateRole: (id: number, data: any) => apiClient.put(`/admin/roles/${id}`, data),
    deleteRole: (id: number) => apiClient.delete(`/admin/roles/${id}`),
    assignRolePermissions: (id: number, permissionIds: number[]) => apiClient.put(`/admin/roles/${id}/permissions`, { permissionIds }),
    assignRoleUser: (roleId: number, userId: number) => apiClient.post('/admin/roles/assign-user', { roleId, userId }),
    getUserRoles: (userId: number) => apiClient.get(`/admin/users/${userId}/roles`),

    // Order Management
    getOrders: (params?: any) => apiClient.get('/admin/orders', { params }),

    // Withdraw Management
    getWithdraws: (params?: any) => apiClient.get('/admin/withdraws', { params }),

    // Dashboard & Stats
    getDashboardStats: () => apiClient.get<ApiResponse<DashboardStats>>('/admin/stats/dashboard'),
    getRevenueTrend: (params?: { days?: number }) => apiClient.get<ApiResponse<TrendData[]>>('/admin/stats/revenue-trend', { params }),
    getUserGrowth: (params?: { days?: number }) => apiClient.get<ApiResponse<TrendData[]>>('/admin/stats/user-growth', { params }),
    getOrderStats: (params?: { days?: number }) => apiClient.get<ApiResponse<OrderStats>>('/admin/stats/orders', { params }),
    getTopPlayers: (params?: { limit?: number }) => apiClient.get<ApiResponse<TopPlayer[]>>('/admin/stats/top-players', { params }),
    getAuditStats: () => apiClient.get<ApiResponse<AuditStats>>('/admin/stats/audit/overview'),

    // Legacy (keep if needed, or remove)
    getStats: () => apiClient.get('/admin/stats'),
};

export interface DashboardStats {
    totalUsers: number;
    totalPlayers: number;
    totalGames: number;
    totalOrders: number;
    ordersByStatus: Record<string, number>;
    paymentsByStatus: Record<string, number>;
    totalPaidAmountCents: number;
}

export interface TrendData {
    date: string;
    value: number;
    type?: string;
}

export interface OrderStats {
    total: number;
    completed: number;
    cancelled: number;
    refunded: number;
}

export interface TopPlayer {
    playerId: number;
    nickname: string;
    ratingAverage: number;
    ratingCount: number;
}

export interface AuditStats {
    pending: number;
    approved: number;
    rejected: number;
}
