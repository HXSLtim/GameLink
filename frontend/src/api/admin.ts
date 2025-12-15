import apiClient from './client';
import type { ApiResponse } from '../types/api';

// Re-export for backward compatibility
export type { ApiResponse, Pagination } from '../types/api';

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
    tags?: string[];
    level?: number;
    vipExpiry?: string;
    wallet?: {
        id: number;
        userId: number;
        balanceCents: number;
        frozenCents: number;
        createdAt: string;
        updatedAt: string;
    };
}

export interface Order {
    id: number;
    orderNumber: string;
    userId: number;
    playerId: number;
    gameId: number;
    title: string;
    description: string;
    totalPriceCents: number;
    currency: string;
    status: 'pending' | 'confirmed' | 'in_progress' | 'completed' | 'cancelled' | 'refunded';
    scheduledStart: string;
    scheduledEnd: string;
    completedAt: string;
    cancelReason: string;
    createdAt: string;
    updatedAt: string;
    user?: { id: number; name: string; avatarUrl?: string };
    player?: { id: number; nickname: string; user?: { avatarUrl?: string } };
    game?: { id: number; name: string };
}

export interface AuditLog {
    id: number;
    action: string;
    ip: string;
    location?: string;
    device?: string;
    createdAt: string;
    details?: string;
}

export interface LoginHistory {
    id: number;
    ip: string;
    location?: string;
    device?: string;
    loginAt: string;
    status: 'success' | 'failed';
}

export interface UserBehaviorStats {
    dau: number;
    avgOnlineTime: string;
    avgConsumption: number;
}

export interface UserDistribution {
    byRegion: Array<{ name: string; value: number }>;
    byAge: Array<{ name: string; value: number }>;
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

export interface UserTag {
    id: number;
    name: string;
    color: string;
    description?: string;
    userCount?: number;
    createdAt: string;
    updatedAt: string;
}

export interface CreateTagDto {
    name: string;
    color: string;
    description?: string;
}

export type UpdateTagDto = Partial<CreateTagDto>;

export interface BatchRoleDto {
    userIds: number[];
    role: string;
}

export interface BatchStatusDto {
    userIds: number[];
    status: string;
}

export interface BatchPointsDto {
    target: 'users' | 'role' | 'all';
    userIds?: number[];
    roles?: string[];
    cents: number;
    reason: string;
    type: string;
}

export interface BatchNotificationDto {
    target: 'users' | 'role' | 'all';
    userIds?: number[];
    roles?: string[];
    title: string;
    content: string;
    type: 'system' | 'marketing' | 'personal' | 'activity';
}

export interface UserOrderParams {
    page?: number;
    page_size?: number;
    status?: string;
}

export interface UserLogParams {
    page?: number;
    page_size?: number;
    type?: string;
}

export interface Game {
    id: string;
    name: string;
    icon: string;
    category: string;
}

export interface CreateGameDto {
    name: string;
    icon?: string;
    category?: string;
    status?: 'active' | 'inactive';
}

export type UpdateGameDto = Partial<CreateGameDto>;

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

export type UpdateMenuDto = Partial<CreateMenuDto>;

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

export interface CreatePermissionDto {
    name: string;
    code: string;
    group: string;
    method?: string;
    path?: string;
    description?: string;
}

export type UpdatePermissionDto = Partial<CreatePermissionDto>;

export interface Role {
    id: number;
    name: string;
    slug: string;
    description: string;
    isSystem: boolean;
    permissions: Permission[] | null;
    users: User[] | null;
    createdAt: string;
    updatedAt: string;
    deletedAt: string | null;
}

export interface CreateRoleDto {
    name: string;
    description?: string;
    permissionIds?: number[];
}

export type UpdateRoleDto = Partial<CreateRoleDto>;

export interface OrderQueryParams {
    page?: number;
    page_size?: number;
    status?: string;
    userId?: number;
    orderNumber?: string;
    dateFrom?: string;
    dateTo?: string;
}

// Player (陪玩师) Interfaces
export interface Player {
    id: number;
    userId: number;
    nickname: string;
    bio: string;
    rank: string;
    hourlyRateCents: number;
    mainGameId: number;
    verificationStatus: 'pending' | 'verified' | 'rejected';
    ratingAverage: number;
    ratingCount: number;
    skillTags: string[];
    createdAt: string;
    updatedAt: string;
    user?: User;
    mainGame?: Game;
}

export interface CreatePlayerDto {
    userId: number;
    nickname?: string;
    bio?: string;
    rank?: string;
    hourlyRateCents?: number;
    mainGameId?: number;
    verificationStatus: string;
}

export interface UpdatePlayerDto {
    nickname?: string;
    bio?: string;
    rank?: string;
    hourlyRateCents?: number;
    mainGameId?: number;
    verificationStatus: string;
}

export interface PlayerQueryParams {
    page?: number;
    page_size?: number;
    status?: string;
    keyword?: string;
}

export interface WithdrawQueryParams {
    page?: number;
    page_size?: number;
    status?: 'pending' | 'approved' | 'rejected' | 'completed';
    playerId?: number;
    dateFrom?: string;
    dateTo?: string;
}

export interface Withdraw {
    id: number;
    playerId: number;
    amountCents: number;
    status: 'pending' | 'approved' | 'rejected' | 'completed';
    bankName?: string;
    bankAccount?: string;
    accountName?: string;
    remark?: string;
    rejectReason?: string;
    adminRemark?: string;
    processedBy?: number;
    processedAt?: string;
    completedAt?: string;
    createdAt: string;
    updatedAt: string;
    player?: Player;
}

export interface ApproveWithdrawDto {
    remark?: string;
}

export interface RejectWithdrawDto {
    reason: string;
}

// Commission Interfaces
export interface CommissionRule {
    id: number;
    name: string;
    ratePercent: number;
    minAmountCents?: number;
    maxAmountCents?: number;
    gameId?: number;
    categoryId?: number;
    isDefault: boolean;
    status: 'active' | 'inactive';
    createdAt: string;
    updatedAt: string;
}

export interface CreateCommissionRuleDto {
    name: string;
    ratePercent: number;
    minAmountCents?: number;
    maxAmountCents?: number;
    gameId?: number;
    categoryId?: number;
    isDefault?: boolean;
}

export interface UpdateCommissionRuleDto extends Partial<CreateCommissionRuleDto> {
    status?: 'active' | 'inactive';
}

export interface PlatformStats {
    month: string;
    totalRevenueCents: number;
    totalCommissionCents: number;
    totalOrderCount: number;
    completedOrderCount: number;
}

export const adminApi = {
    // User Management
    getUsers: (params?: UserQueryParams) => apiClient.get<ApiResponse<User[]>>('/admin/users', { params }),
    getUser: (id: number) => apiClient.get<ApiResponse<User>>(`/admin/users/${id}`),
    createUser: (data: CreateUserDto) => apiClient.post<ApiResponse<User>>('/admin/users', data),
    updateUser: (id: number, data: UpdateUserDto) => apiClient.put<ApiResponse<User>>(`/admin/users/${id}`, data),
    deleteUser: (id: number) => apiClient.delete<ApiResponse<void>>(`/admin/users/${id}`),
    updateUserStatus: (id: number, status: string) => apiClient.put<ApiResponse<User>>(`/admin/users/${id}/status`, { status }),
    updateUserRole: (id: number, role: string) => apiClient.put<ApiResponse<User>>(`/admin/users/${id}/role`, { role }),
    getUserOrders: (id: number, params?: UserOrderParams) => apiClient.get<ApiResponse<Order[]>>(`/admin/users/${id}/orders`, { params }),
    getUserLogs: (id: number, params?: UserLogParams) => apiClient.get<ApiResponse<AuditLog[]>>(`/admin/users/${id}/logs`, { params }),
    getUserLoginHistory: (id: number, params?: { page?: number; page_size?: number }) => apiClient.get<ApiResponse<LoginHistory[]>>(`/admin/users/${id}/login-history`, { params }),

    // User Behavior Analysis
    getUserBehaviorStats: () =>
        apiClient.get<ApiResponse<UserBehaviorStats>>('/admin/users/behavior/stats'),
    getUserActivityTrend: (params?: { days?: number }) =>
        apiClient.get<ApiResponse<TrendData[]>>('/admin/users/behavior/trend', { params }),
    getUserDistribution: () =>
        apiClient.get<ApiResponse<UserDistribution>>('/admin/users/behavior/distribution'),

    // User Tags
    getTags: () => apiClient.get<ApiResponse<UserTag[]>>('/admin/tags'),
    createTag: (data: CreateTagDto) => apiClient.post<ApiResponse<UserTag>>('/admin/tags', data),
    updateTag: (id: number, data: UpdateTagDto) => apiClient.put<ApiResponse<UserTag>>(`/admin/tags/${id}`, data),
    deleteTag: (id: number) => apiClient.delete<ApiResponse<void>>(`/admin/tags/${id}`),
    assignUserTag: (userId: number, tagId: number) => apiClient.post<ApiResponse<void>>(`/admin/users/${userId}/tags/${tagId}`),
    removeUserTag: (userId: number, tagId: number) => apiClient.delete<ApiResponse<void>>(`/admin/users/${userId}/tags/${tagId}`),
    batchSetUserTags: (userId: number, tagIds: number[]) => apiClient.put<ApiResponse<void>>(`/admin/users/${userId}/tags`, { tagIds }),
    getUserTags: (userId: number) => apiClient.get<ApiResponse<UserTag[]>>(`/admin/users/${userId}/tags`),

    // Batch Operations
    batchUpdateUserRole: (data: BatchRoleDto) => apiClient.post<ApiResponse<void>>('/admin/users/batch/role', data),
    batchUpdateUserStatus: (data: BatchStatusDto) => apiClient.post<ApiResponse<void>>('/admin/users/batch/status', data),
    batchDeleteUsers: (userIds: number[]) => apiClient.post<ApiResponse<void>>('/admin/users/batch/delete', { userIds }),
    batchAddUserPoints: (data: BatchPointsDto) => apiClient.post<ApiResponse<void>>('/admin/users/batch/points', data),
    batchSendNotification: (data: BatchNotificationDto) => apiClient.post<ApiResponse<void>>('/admin/users/batch/notification', data),

    // Game Management
    getGames: (params?: { status?: string; page_size?: number }) => apiClient.get('/admin/games', { params }),
    createGame: (data: CreateGameDto) => apiClient.post('/admin/games', data),
    updateGame: (id: string, data: UpdateGameDto) => apiClient.put(`/admin/games/${id}`, data),
    deleteGame: (id: string) => apiClient.delete(`/admin/games/${id}`),
    // Game Batch Operations
    batchDeleteGames: (gameIds: string[]) => apiClient.post<ApiResponse<void>>('/admin/games/batch/delete', { gameIds }),

    // Service Item Management
    getServiceItems: (params?: ServiceItemQueryParams) => apiClient.get<ApiResponse<ServiceItem[]>>('/admin/service-items', { params }),
    getServiceItem: (id: number) => apiClient.get<ApiResponse<ServiceItem>>(`/admin/service-items/${id}`),
    createServiceItem: (data: CreateServiceItemDto) => apiClient.post<ApiResponse<ServiceItem>>('/admin/service-items', data),
    updateServiceItem: (id: number, data: UpdateServiceItemDto) => apiClient.put<ApiResponse<ServiceItem>>(`/admin/service-items/${id}`, data),
    deleteServiceItem: (id: number) => apiClient.delete<ApiResponse<void>>(`/admin/service-items/${id}`),
    batchUpdateServiceItemStatus: (ids: number[], status: 'active' | 'inactive') => apiClient.put<ApiResponse<void>>('/admin/service-items/batch-status', { ids, status }),

    // Menu Management
    getMenus: (params?: { parentId?: number; page?: number; page_size?: number }) => apiClient.get<ApiResponse<Menu[]>>('/admin/menus', { params }),
    getMyMenus: () => apiClient.get<ApiResponse<Menu[]>>('/admin/menus/me'),
    getMenu: (id: number) => apiClient.get<ApiResponse<Menu>>(`/admin/menus/${id}`),
    createMenu: (data: CreateMenuDto) => apiClient.post<ApiResponse<Menu>>('/admin/menus', data),
    updateMenu: (id: number, data: UpdateMenuDto) => apiClient.put<ApiResponse<Menu>>(`/admin/menus/${id}`, data),
    deleteMenu: (id: number) => apiClient.delete<ApiResponse<void>>(`/admin/menus/${id}`),

    // Permissions
    getPermissions: (params?: { page?: number; page_size?: number; group?: string }) => apiClient.get<ApiResponse<{ items: Permission[], totalCount: number, page: number, pageSize: number }>>('/admin/permissions', { params }),
    getPermissionGroups: () => apiClient.get<ApiResponse<string[]>>('/admin/permissions/groups'),
    getPermission: (id: number) => apiClient.get<ApiResponse<Permission>>(`/admin/permissions/${id}`),
    createPermission: (data: CreatePermissionDto) => apiClient.post<ApiResponse<Permission>>('/admin/permissions', data),
    updatePermission: (id: number, data: UpdatePermissionDto) => apiClient.put<ApiResponse<Permission>>(`/admin/permissions/${id}`, data),
    deletePermission: (id: number) => apiClient.delete<ApiResponse<void>>(`/admin/permissions/${id}`),
    getMyPermissions: () => apiClient.get<ApiResponse<string[]>>('/admin/permissions/me'),

    // Role Management
    getRoles: (params?: { page?: number; page_size?: number }) => apiClient.get<ApiResponse<{ items: Role[], totalCount: number, page: number, pageSize: number }>>('/admin/roles', { params }),
    getRole: (id: number) => apiClient.get(`/admin/roles/${id}`),
    createRole: (data: CreateRoleDto) => apiClient.post('/admin/roles', data),
    updateRole: (id: number, data: UpdateRoleDto) => apiClient.put(`/admin/roles/${id}`, data),
    deleteRole: (id: number) => apiClient.delete(`/admin/roles/${id}`),
    assignRolePermissions: (id: number, permissionIds: number[]) => apiClient.put(`/admin/roles/${id}/permissions`, { permissionIds }),
    assignRoleUser: (roleId: number, userId: number) => apiClient.post('/admin/roles/assign-user', { roleId, userId }),
    getUserRoles: (userId: number) => apiClient.get(`/admin/users/${userId}/roles`),
    getUserStats: () => apiClient.get<ApiResponse<UserStats>>('/admin/users/stats'),

    // Player Management
    getPlayers: (params?: PlayerQueryParams) => apiClient.get<ApiResponse<Player[]>>('/admin/players', { params }),
    getPlayer: (id: number) => apiClient.get<ApiResponse<Player>>(`/admin/players/${id}`),
    createPlayer: (data: CreatePlayerDto) => apiClient.post<ApiResponse<Player>>('/admin/players', data),
    updatePlayer: (id: number, data: UpdatePlayerDto) => apiClient.put<ApiResponse<Player>>(`/admin/players/${id}`, data),
    deletePlayer: (id: number) => apiClient.delete<ApiResponse<void>>(`/admin/players/${id}`),
    updatePlayerVerification: (id: number, status: string) => apiClient.put<ApiResponse<Player>>(`/admin/players/${id}/verification`, { verification_status: status }),
    updatePlayerSkillTags: (id: number, tags: string[]) => apiClient.put<ApiResponse<void>>(`/admin/players/${id}/skill-tags`, { tags }),
    // Player Batch Operations
    batchUpdatePlayerStatus: (data: { playerIds: number[]; status: string }) => apiClient.post<ApiResponse<void>>('/admin/players/batch/status', data),
    batchDeletePlayers: (playerIds: number[]) => apiClient.post<ApiResponse<void>>('/admin/players/batch/delete', { playerIds }),

    // Order Management
    getOrders: (params?: OrderQueryParams) => apiClient.get<ApiResponse<Order[]>>('/admin/orders', { params }),
    getOrder: (id: number) => apiClient.get<ApiResponse<Order>>(`/admin/orders/${id}`),
    cancelOrder: (id: number, note?: string) => apiClient.post<ApiResponse<Order>>(`/admin/orders/${id}/cancel`, { note }),
    refundOrder: (id: number, data: { reason: string; amount_cents: number; note?: string }) => apiClient.post<ApiResponse<Order>>(`/admin/orders/${id}/refund`, data),
    // Order Batch Operations
    batchCancelOrders: (orderIds: number[], reason?: string) => apiClient.post<ApiResponse<void>>('/admin/orders/batch/cancel', { orderIds, reason }),
    batchCompleteOrders: (orderIds: number[]) => apiClient.post<ApiResponse<void>>('/admin/orders/batch/complete', { orderIds }),

    // Withdraw Management
    getWithdraws: (params?: WithdrawQueryParams) => apiClient.get<ApiResponse<{ withdraws: Withdraw[]; total: number }>>('/admin/withdraws', { params }),
    getWithdraw: (id: number) => apiClient.get<ApiResponse<Withdraw>>(`/admin/withdraws/${id}`),
    approveWithdraw: (id: number, data?: ApproveWithdrawDto) => apiClient.post<ApiResponse<void>>(`/admin/withdraws/${id}/approve`, data),
    rejectWithdraw: (id: number, data: RejectWithdrawDto) => apiClient.post<ApiResponse<void>>(`/admin/withdraws/${id}/reject`, data),
    completeWithdraw: (id: number) => apiClient.post<ApiResponse<void>>(`/admin/withdraws/${id}/complete`),

    // Commission Management
    createCommissionRule: (data: CreateCommissionRuleDto) => apiClient.post<ApiResponse<CommissionRule>>('/admin/commission/rules', data),
    updateCommissionRule: (id: number, data: UpdateCommissionRuleDto) => apiClient.put<ApiResponse<void>>(`/admin/commission/rules/${id}`, data),
    triggerSettlement: (month?: string) => apiClient.post<ApiResponse<void>>('/admin/commission/settlements/trigger', null, { params: { month } }),
    getPlatformStats: (month?: string) => apiClient.get<ApiResponse<PlatformStats>>('/admin/commission/stats', { params: { month } }),

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

export interface UserStats {
    total: number;
    byRole: {
        user: number;
        player: number;
        admin: number;
    };
    byStatus: {
        active: number;
        banned: number;
        suspended: number;
    };
    recentRegistrations: number; // 最近7天注册数
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
