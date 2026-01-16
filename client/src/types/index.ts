/**
 * Common type definitions for GameLink Client
 */

// API Response wrapper
export interface ApiResponse<T> {
    success: boolean;
    code: number;
    message: string;
    data: T;
}

// Pagination
export interface Pagination {
    page: number;
    pageSize: number;
    total: number;
    totalPages: number;
}

export interface PaginatedResponse<T> {
    items: T[];
    pagination: Pagination;
}

// User types (backend uses 'name' not 'nickname')
export interface User {
    id: number;
    phone: string;
    name: string; // Backend User.Name
    avatar: string;
    gender: 'male' | 'female' | 'unknown';
    birthday?: string;
    bio?: string;
    status: 'active' | 'banned' | 'inactive';
    isPlayer: boolean;
    createdAt: string;
    updatedAt: string;
}

// Player types
export interface Player {
    id: number;
    userId: number;
    nickname: string;
    avatar: string;
    gender: 'male' | 'female' | 'unknown';
    bio?: string;
    status: 'pending' | 'approved' | 'rejected' | 'suspended';
    rating: number;
    orderCount: number;
    services: PlayerService[];
    tags: string[];
    isOnline: boolean;
    createdAt: string;
}

export interface PlayerService {
    id: number;
    gameId: number;
    gameName: string;
    gameIcon: string;
    serviceItemId: number;
    serviceItemName: string;
    rankId: number;
    rankName: string;
    priceCents: number;
    status: 'active' | 'inactive';
}

// Order types (matching backend OrderStatus)
export type OrderStatus =
    | 'pending'      // Initial state
    | 'confirmed'    // Paid, waiting for player
    | 'in_progress'  // Player joined, service in progress
    | 'completed'    // Service completed
    | 'canceled'     // Order cancelled
    | 'refunded'     // Order refunded
    | 'disputed';    // Dispute raised

export type OrderType = 'solo' | 'team' | 'gift';

export interface Order {
    id: number;
    orderNo: string;
    userId: number;
    type: OrderType;
    status: OrderStatus;
    totalCents: number;
    paidCents: number;
    quantity: number;
    unit: string;
    gameName: string;
    serviceItemName: string;
    players: OrderPlayer[];
    createdAt: string;
    paidAt?: string;
    completedAt?: string;
}

export interface OrderPlayer {
    id: number;
    playerId: number;
    playerName: string;
    playerAvatar: string;
    status: 'pending' | 'accepted' | 'rejected' | 'completed';
}

// Wallet types
export interface Wallet {
    id: number;
    userId: number;
    balanceCents: number;
    frozenCents: number;
    totalRechargeCents: number;
    totalSpentCents: number;
}

export interface WalletTransaction {
    id: number;
    walletId: number;
    type: 'recharge' | 'payment' | 'refund' | 'withdraw' | 'bonus';
    amountCents: number;
    balanceAfterCents: number;
    description: string;
    createdAt: string;
}

// VIP types
export interface VipLevel {
    id: number;
    name: string;
    level: number;
    icon: string;
    requiredPoints: number;
    discountPercent: number;
    benefits: string[];
}

export interface UserVip {
    userId: number;
    levelId: number;
    level: VipLevel;
    points: number;
    expireAt?: string;
}

// Chat types
export interface ChatGroup {
    id: number;
    name: string;
    type: 'public' | 'order';
    avatar?: string;
    memberCount: number;
    lastMessage?: ChatMessage;
    unreadCount: number;
}

export interface ChatMessage {
    id: number;
    groupId: number;
    senderId: number;
    senderName: string;
    senderAvatar: string;
    content: string;
    type: 'text' | 'image' | 'system';
    createdAt: string;
}

// Notification types
export interface Notification {
    id: number;
    userId: number;
    type: 'system' | 'order' | 'chat' | 'wallet' | 'promotion';
    title: string;
    content: string;
    isRead: boolean;
    data?: Record<string, unknown>;
    createdAt: string;
}

// Game types
export interface Game {
    id: number;
    name: string;
    icon: string;
    banner?: string;
    description?: string;
    status: 'active' | 'inactive';
    sortOrder: number;
}

export interface ServiceItem {
    id: number;
    gameId: number;
    name: string;
    description?: string;
    unit: string;
    basePriceCents: number;
    status: 'active' | 'inactive';
}
