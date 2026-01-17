/**
 * API Response Types for GameLink Client
 * Complete type definitions for all API modules
 */

// ============ Common Types ============

export interface PaginatedResponse<T> {
    items: T[];
    total: number;
    page: number;
    pageSize: number;
    totalPages: number;
}

export interface PaginationParams {
    page: number;
    pageSize: number;
}

// ============ Auth Types ============

export interface LoginRequest {
    username: string;
    password: string;
}

export interface LoginResponse {
    token: string;
    refreshToken: string;
    user: {
        id: string;
        username: string;
        avatar: string;
        email?: string;
        nickname?: string;
    };
    role?: string;
    permissions?: string[];
}

export interface RegisterRequest {
    phone?: string;
    email?: string;
    password: string;
    name: string;
}

export interface RegisterResponse extends LoginResponse {}

export interface RefreshResponse {
    token: string;
    refreshToken?: string;
    user?: {
        id: string;
        username: string;
        avatar: string;
        email?: string;
        nickname?: string;
    };
}

export interface MeResponse {
    user?: {
        id: string;
        username: string;
        avatar: string;
        email?: string;
        nickname?: string;
        role?: string;
    };
    id?: string;
    username?: string;
    avatar?: string;
    email?: string;
    nickname?: string;
    role?: string;
}

// ============ User Types ============

export interface User {
    id: number;
    username: string;
    name: string;
    avatar: string;
    email?: string;
    phone?: string;
    role: 'user' | 'player' | 'admin';
    status: 'active' | 'suspended' | 'banned';
    createdAt: string;
    updatedAt: string;
}

export interface UpdateUserRequest {
    name?: string;
    avatar?: string;
    email?: string;
    phone?: string;
}

export interface UserPreferences {
    language: string;
    theme: 'light' | 'dark' | 'auto';
    notifications: {
        email: boolean;
        push: boolean;
        sms: boolean;
    };
}

// ============ Player Types ============

export interface Player {
    id: number;
    userId: number;
    username: string;
    name: string;
    avatar: string;
    rating: number;
    reviewCount: number;
    completedOrders: number;
    online: boolean;
    introduction?: string;
    games: Array<{
        gameId: number;
        gameName: string;
        rankId: number;
        rankName: string;
    }>;
    createdAt: string;
}

export interface PlayerProfile extends Player {
    totalEarnings: number;
    totalOrders: number;
    acceptanceRate: number;
    completionRate: number;
    averageResponseTime: number;
}

export interface UpdatePlayerProfileRequest {
    name?: string;
    avatar?: string;
    introduction?: string;
    games?: Array<{
        gameId: number;
        rankId: number;
    }>;
}

export interface PlayerSearchParams extends PaginationParams {
    gameId?: number;
    rankId?: number;
    minRating?: number;
    maxRating?: number;
    online?: boolean;
    search?: string;
}

// ============ Order Types ============

export interface Order {
    id: number;
    userId: number;
    playerId: number;
    serviceItemId: number;
    type: 'solo' | 'team' | 'gift';
    status: 'pending' | 'accepted' | 'in_progress' | 'completed' | 'cancelled';
    amount: number;
    duration: number;
    startTime?: string;
    endTime?: string;
    createdAt: string;
    updatedAt: string;
}

export interface CreateOrderRequest {
    playerId: number;
    serviceItemId: number;
    duration: number;
    couponId?: number;
    message?: string;
}

export interface UpdateOrderRequest {
    status?: string;
    endTime?: string;
}

export interface OrderListParams extends PaginationParams {
    status?: string;
    type?: string;
    startDate?: string;
    endDate?: string;
}

// ============ Payment Types ============

export interface Payment {
    id: number;
    orderId: number;
    userId: number;
    amount: number;
    method: string;
    status: 'pending' | 'paid' | 'failed' | 'refunded';
    paidAt?: string;
    createdAt: string;
}

export interface CreatePaymentRequest {
    orderId: number;
    method: string;
    amount: number;
}

export interface PaymentListParams extends PaginationParams {
    status?: string;
    startDate?: string;
    endDate?: string;
}

// ============ Chat Types ============

export interface ChatRoom {
    id: number;
    userId: number;
    otherUserId: number;
    otherUser: {
        id: number;
        username: string;
        avatar: string;
        online: boolean;
    };
    lastMessage?: ChatMessage;
    unreadCount: number;
    createdAt: string;
    updatedAt: string;
}

export interface ChatMessage {
    id: number;
    roomId: number;
    senderId: number;
    content: string;
    type: 'text' | 'image' | 'system';
    read: boolean;
    createdAt: string;
}

export interface SendMessageRequest {
    content: string;
    type?: 'text' | 'image';
}

// ============ Dispute Types ============

export interface Dispute {
    id: number;
    orderId: number;
    userId: number;
    playerId: number;
    reason: string;
    description: string;
    status: 'pending' | 'processing' | 'resolved' | 'rejected';
    resolution?: string;
    createdAt: string;
    updatedAt: string;
}

export interface CreateDisputeRequest {
    orderId: number;
    reason: string;
    description: string;
    images?: string[];
}

export interface DisputeListParams extends PaginationParams {
    status?: string;
    orderId?: number;
}

// ============ Wallet Types ============

export interface Wallet {
    id: number;
    userId: number;
    balanceCents: number;
    frozenCents: number;
    totalIncome: number;
    totalWithdrawn: number;
    updatedAt: string;
}

export interface Transaction {
    id: number;
    walletId: number;
    type: 'recharge' | 'payment' | 'refund' | 'withdraw' | 'income';
    amount: number;
    balance: number;
    description: string;
    createdAt: string;
}

export interface WithdrawRequest {
    amount: number;
    method: string;
    account: string;
}

export interface TransactionListParams extends PaginationParams {
    type?: string;
    startDate?: string;
    endDate?: string;
}

// Error response
export interface ApiError {
    response?: {
        data?: {
            message?: string;
        };
    };
    message?: string;
}

// Helper to extract error message
export function getErrorMessage(err: unknown): string {
    if (err && typeof err === 'object') {
        const apiErr = err as ApiError;
        return apiErr.response?.data?.message || apiErr.message || 'An error occurred';
    }
    return 'An error occurred';
}

// ============ Game Types ============

export interface Game {
    id: number;
    name: string;
    icon: string;
    category: string;
    description?: string;
    playerCount: number;
    createdAt: string;
}

export interface GameRank {
    id: number;
    gameId: number;
    name: string;
    level: number;
    icon?: string;
}

// ============ Service Item Types ============

export interface ServiceItem {
    id: number;
    gameId: number;
    playerId?: number;
    name: string;
    description: string;
    price: number;
    duration: number;
    commissionRate: number;
    status: 'active' | 'inactive';
    createdAt: string;
}

// ============ Review Types ============

export interface Review {
    id: number;
    orderId: number;
    userId: number;
    playerId: number;
    rating: number;
    content: string;
    images?: string[];
    reply?: string;
    repliedAt?: string;
    createdAt: string;
}

export interface CreateReviewRequest {
    orderId: number;
    rating: number;
    content: string;
    images?: string[];
}

export interface ReviewListParams extends PaginationParams {
    playerId?: number;
    userId?: number;
    rating?: number;
}

// ============ Notification Types ============

export interface Notification {
    id: number;
    userId: number;
    type: 'order' | 'chat' | 'system' | 'promotion';
    title: string;
    content: string;
    data?: Record<string, any>;
    read: boolean;
    createdAt: string;
}

export interface NotificationListParams extends PaginationParams {
    type?: string;
    read?: boolean;
}

// ============ VIP Types ============

export interface VipLevel {
    id: number;
    level: number;
    name: string;
    icon: string;
    requiredExp: number;
    benefits: string[];
    discount: number;
    monthlyCouponAmount: number;
    createdAt: string;
}

export interface VipBenefit {
    id: number;
    levelId: number;
    type: string;
    name: string;
    description: string;
    value: string;
}

export interface UserVipInfo {
    userId: number;
    levelId: number;
    level: VipLevel;
    exp: number;
    unlocked: boolean;
    unlockedAt?: string;
    expireAt?: string;
    nextLevel?: VipLevel;
    progressToNext: number;
}

// ============ Coupon Types ============

export interface Coupon {
    id: number;
    code: string;
    name: string;
    type: 'discount' | 'fixed' | 'gift';
    value: number;
    minAmount: number;
    maxDiscount?: number;
    startTime: string;
    endTime: string;
    totalCount: number;
    usedCount: number;
    status: 'active' | 'inactive' | 'expired';
    createdAt: string;
}

export interface CouponListParams extends PaginationParams {
    type?: string;
    status?: string;
}

// ============ Recharge Types ============

export interface RechargePackage {
    id: number;
    amount: number;
    bonus: number;
    totalAmount: number;
    popular: boolean;
    description?: string;
}

export interface CreateRechargeRequest {
    packageId: number;
    paymentMethod: string;
}

export interface RechargeRecord {
    id: number;
    userId: number;
    amount: number;
    bonus: number;
    totalAmount: number;
    paymentId: number;
    status: 'pending' | 'completed' | 'failed';
    completedAt?: string;
    createdAt: string;
}


// ============ Activity Types ============

export interface Activity {
    id: number;
    name: string;
    description: string;
    type: string;
    startTime: string;
    endTime: string;
    status: 'upcoming' | 'ongoing' | 'ended';
    rewards?: string[];
    createdAt: string;
}

export interface ActivityListParams extends PaginationParams {
    type?: string;
    status?: string;
}

// ============ Team Types ============

export interface Team {
    id: number;
    leaderId: number;
    name: string;
    description?: string;
    gameId: number;
    maxMembers: number;
    currentMembers: number;
    status: 'recruiting' | 'full' | 'disbanded';
    createdAt: string;
}

export interface CreateTeamRequest {
    name: string;
    description?: string;
    gameId: number;
    maxMembers: number;
}

export interface TeamMember {
    userId: number;
    username: string;
    avatar: string;
    role: 'leader' | 'member';
    joinedAt: string;
}

export interface TeamOrder {
    id: number;
    teamId: number;
    serviceItemId: number;
    status: string;
    amount: number;
    createdAt: string;
}

// ============ Referral Types ============

export interface ReferralInfo {
    code: string;
    totalReferrals: number;
    successfulReferrals: number;
    totalRewards: number;
    pendingRewards: number;
}

export interface ReferralRecord {
    id: number;
    referrerId: number;
    referredUserId: number;
    referredUser: {
        username: string;
        avatar: string;
    };
    reward: number;
    status: 'pending' | 'completed';
    createdAt: string;
}

// ============ Feature Types ============

export interface BlockedUser {
    id: number;
    userId: number;
    blockedUserId: number;
    blockedUser: {
        username: string;
        avatar: string;
    };
    reason?: string;
    createdAt: string;
}

export interface CommissionRule {
    id: number;
    playerId?: number;
    serviceItemId?: number;
    rate: number;
    type: 'player' | 'item' | 'ranking';
    createdAt: string;
}

export interface CommissionRecord {
    id: number;
    orderId: number;
    playerId: number;
    amount: number;
    rate: number;
    createdAt: string;
}

export interface Certification {
    id: number;
    playerId: number;
    type: 'real_name' | 'skill';
    status: 'pending' | 'approved' | 'rejected';
    createdAt: string;
}

export interface RealNameVerification {
    id: number;
    playerId: number;
    realName: string;
    idCard: string;
    status: 'pending' | 'approved' | 'rejected';
    reason?: string;
    submittedAt: string;
    reviewedAt?: string;
}

export interface SkillAuthentication {
    id: number;
    playerId: number;
    gameId: number;
    rankId: number;
    proofImages: string[];
    status: 'pending' | 'approved' | 'rejected';
    reason?: string;
    submittedAt: string;
    reviewedAt?: string;
}

export interface RankingPlayer {
    playerId: number;
    username: string;
    avatar: string;
    rank: number;
    score: number;
    change: number;
}

export interface RankingConfig {
    period: 'daily' | 'weekly' | 'monthly';
    scoreFormula: string;
    rewards: Array<{
        rank: number;
        reward: string;
    }>;
}

// ============ Special Feature Types ============

export interface Room {
    id: number;
    ownerId: number;
    gameId: number;
    name: string;
    description?: string;
    maxMembers: number;
    currentMembers: number;
    isPrivate: boolean;
    status: 'open' | 'in_game' | 'closed';
    createdAt: string;
}

export interface CreateRoomRequest {
    gameId: number;
    name: string;
    description?: string;
    maxMembers: number;
    isPrivate?: boolean;
    password?: string;
}

export interface RoomListParams extends PaginationParams {
    gameId?: number;
    status?: string;
    isPrivate?: boolean;
}

export interface RoomMember {
    userId: number;
    username: string;
    avatar: string;
    role: 'owner' | 'member';
    joinedAt: string;
}

export interface LFGPost {
    id: number;
    userId: number;
    gameId: number;
    title: string;
    description: string;
    requiredPlayers: number;
    currentPlayers: number;
    status: 'open' | 'full' | 'closed';
    createdAt: string;
}

export interface CreateLFGRequest {
    gameId: number;
    title: string;
    description: string;
    requiredPlayers: number;
}

export interface LFGListParams extends PaginationParams {
    gameId?: number;
    status?: string;
}

export interface VoiceChannel {
    id: number;
    roomId: number;
    name: string;
    maxParticipants: number;
    currentParticipants: number;
    status: 'active' | 'inactive';
}

export interface VoiceSession {
    id: string;
    channelId: number;
    userId: number;
    joinedAt: string;
    muted: boolean;
}

export interface UserPresence {
    userId: number;
    status: 'online' | 'away' | 'busy' | 'offline';
    customStatus?: string;
    lastSeen: string;
}

export interface Gift {
    id: number;
    name: string;
    description: string;
    icon: string;
    price: number;
    category: string;
}

export interface CreateGiftOrderRequest {
    giftId: number;
    recipientId: number;
    message?: string;
}

export interface GiftOrder {
    id: number;
    giftId: number;
    senderId: number;
    recipientId: number;
    amount: number;
    message?: string;
    status: 'completed';
    createdAt: string;
}
