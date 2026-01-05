// ============================================
// Taro App Types
// ============================================

export interface UserInfo {
  id: number;
  name: string;
  phone?: string;
  email?: string;
  avatar?: string;
  role: string;
  status?: 'active' | 'inactive' | 'blocked';
  permissions?: string[];
  createdAt: string;
  updatedAt?: string;
}

export interface LoginRequest {
  phone: string;
  code: string;
}

export interface LoginResponse {
  token: string;
  user: UserInfo;
}

// User types
export interface User {
  id: number;
  name: string;
  phone: string;
  status: 'active' | 'inactive' | 'blocked';
  createdAt: string;
}

// Order types
export type OrderStatus =
  | 'pending'       // Created, waiting for payment
  | 'paid'          // Paid, waiting for players to accept
  | 'in_progress'   // Service in progress
  | 'completed'     // Service completed
  | 'canceled'      // Canceled
  | 'refunded';     // Refunded

export interface Order {
  id: number;
  orderNo: string;
  userId: number;
  playerIds: number[];
  playerNames?: string[];       // Player names for display
  playerAvatars?: string[];     // Player avatars for display
  status: OrderStatus;
  amount: number;               // Total amount in cents (分)
  amountYuan?: number;          // Total amount in yuan (元) - computed
  duration: number;             // Duration in hours
  createdAt: string;
  updatedAt?: string;
  completedAt?: string;
  gameId?: number;
  gameName?: string;
  itemId?: number;              // Service item ID
  itemName?: string;            // Service item name
  scheduledStart?: string;      // Scheduled start time
  hasDispute?: boolean;         // Whether order has dispute
  canCancel?: boolean;          // Whether can cancel
  canReview?: boolean;          // Whether can review
  canDispute?: boolean;         // Whether can file dispute
}

export interface CreateOrderRequest {
  itemId: number;                // Service item ID
  playerId?: number;             // Specific player ID (optional)
  quantity?: number;             // Quantity (hours), default 1
  scheduledStart?: string;       // Scheduled start time (ISO 8601)
  gameIds?: number[];            // Game IDs
  remark?: string;               // Order remark
}

export interface OrderDraft {
  itemId: number;
  itemName: string;
  itemPrice: number;             // Price in cents
  playerId?: number;
  playerName?: string;
  playerAvatar?: string;
  quantity: number;
  totalPrice: number;            // Total in cents
  scheduledStart?: string;
  gameIds?: number[];
  remark?: string;
}

export interface OrderListParams {
  status?: OrderStatus | 'all';  // Filter by status
  page?: number;                 // Page number, default 1
  pageSize?: number;             // Page size, default 10
}

export interface PaymentRequest {
  orderId: number;
  method: 'wechat' | 'alipay' | 'wallet';
  amount?: number;               // Amount in cents (for partial payment)
}

export interface ReviewRequest {
  orderId: number;
  playerId: number;
  rating: number;                // 1-5
  content?: string;
  tags?: string[];               // Review tags
}

export interface DisputeRequest {
  orderId: number;
  type: string;                  // Dispute type code
  reason: string;                // Detailed reason
  evidenceUrls?: string[];       // Evidence image URLs (max 5)
  evidenceText?: string;         // Text evidence
}

// Player types
export interface Player {
  id: number;
  userId: number;
  nickname: string;
  avatar: string;
  rank: number;
  level: number;
  pricePerHour: number;
  status: 'available' | 'busy' | 'offline';
  tags: string[];
  // Detail fields
  bio?: string;
  games?: string[];
  rating?: number;
  reviewCount?: number;
  orderCount?: number;
  isOnline?: boolean;
  lastOnlineAt?: string;
  // Computed fields
  isFavorite?: boolean;
}

export interface PlayerDetail extends Player {
  bio: string;
  games: string[];
  rating: number;
  reviewCount: number;
  orderCount: number;
  isOnline: boolean;
  lastOnlineAt: string;
  images?: PlayerImage[];
  reviews?: PlayerReview[];
}

export interface PlayerImage {
  id: number;
  playerId: number;
  url: string;
  description?: string;
  createdAt: string;
}

export interface PlayerReview {
  id: number;
  orderId: number;
  userId: number;
  userName: string;
  userAvatar?: string;
  rating: number;
  content: string;
  createdAt: string;
}

export interface PlayerFilters {
  // Price range (in yuan)
  minPrice?: number;
  maxPrice?: number;
  // Game filter
  games?: string[];
  // Tags filter
  tags?: string[];
  // Rank filter (1-5)
  rank?: number;
  // Status filter
  status?: Player['status'];
  // Search keyword
  keyword?: string;
}

export interface PlayerListResponse {
  players: Player[];
  total: number;
  page: number;
  pageSize: number;
}

// ============================================
// Chat Types
// ============================================

export type ChatMessageType = 'text' | 'image' | 'file' | 'system' | 'voice' | 'emoji';

export type ChatMessageAuditStatus = 'pending' | 'approved' | 'rejected' | 'deleted';

export type ChatGroupType = 'public' | 'order';

export type ChatMemberRole = 'owner' | 'admin' | 'member';

export interface ChatMessage {
  id: number;
  groupId: number;
  senderId: number;
  senderName?: string;
  senderAvatar?: string;
  content: string;
  messageType: ChatMessageType;
  imageUrl?: string;
  replyToId?: number;
  metadata?: Record<string, unknown>;
  isDeleted: boolean;
  auditStatus: ChatMessageAuditStatus;
  createdAt: string;
  readAt?: string;
  // UI state (not from API)
  sending?: boolean;
  sendError?: boolean;
  retryCount?: number;
}

export interface ChatGroupMember {
  id: number;
  groupId: number;
  userId: number;
  nickname?: string;
  role: ChatMemberRole;
  joinedAt: string;
  lastReadAt?: string;
  lastReadMessageId?: number;
  isMuted: boolean;
  isActive: boolean;
}

export interface ChatGroup {
  id: number;
  groupName: string;
  groupType: ChatGroupType;
  relatedOrderId?: number;
  createdBy: number;
  maxMembers: number;
  isActive: boolean;
  avatarUrl?: string;
  description?: string;
  messageRetentionDays: number;
  members: ChatGroupMember[];
  lastMessage?: ChatMessage;
  unreadCount: number;
  createdAt: string;
  // UI state
  isOnline?: boolean;
}

export interface ChatRoom {
  id: number;
  groupId: number;
  orderId?: number;
  name: string;
  avatar?: string;
  type: ChatGroupType;
  participants: ChatGroupMember[];
  lastMessage?: ChatMessage;
  unreadCount: number;
  isActive: boolean;
  isMuted: boolean;
  createdAt: string;
}

// WebSocket message types
export interface WSMessage {
  type: string;
  timestamp: string;
  data?: unknown;
}

export type WSMessageType =
  | 'ping'
  | 'pong'
  | 'subscribe'
  | 'chat_message'
  | 'chat_typing'
  | 'chat_read'
  | 'system_status'
  | 'online_users'
  | 'order_queue'
  | 'alert';

export interface ChatTypingData {
  groupId: number;
  userId: number;
  userName: string;
}

export interface ChatReadData {
  groupId: number;
  userId: number;
  messageId: number;
}

export interface ChatMessageData {
  message: ChatMessage;
  groupId: number;
}

// Upload state
export interface UploadProgress {
  fileId: string;
  fileName: string;
  progress: number;
  status: 'uploading' | 'success' | 'error';
  url?: string;
  error?: string;
}
