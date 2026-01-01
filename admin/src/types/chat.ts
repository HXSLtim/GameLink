/**
 * Chat Module Types
 * 聊天模块类型定义
 *
 * Contains types for:
 * - Chat Conversation (聊天会话)
 * - Chat Message (聊天消息)
 * - Chat Broadcast (系统广播)
 * - Chat Statistics (聊天统计)
 * - Mute Management (禁言管理)
 */

// ============================================================================
// Chat Conversation Types (聊天会话)
// ============================================================================

/**
 * Chat Conversation Type
 * 聊天会话类型
 */
export type ChatConversationType =
    | 'user_order'   // 订单会话 (用户-陪玩师)
    | 'group_chat';  // 群聊会话 (多人)

/**
 * Chat Conversation Status
 * 聊天会话状态
 */
export type ChatConversationStatus =
    | 'active'       // 进行中
    | 'closed';      // 已关闭

/**
 * Chat Conversation
 * 聊天会话
 */
export interface ChatConversation {
    id: number;
    type: ChatConversationType;
    orderId?: number;                 // 关联订单ID
    orderNo?: string;                 // 订单号
    userId?: number;                  // 用户ID
    userName?: string;                // 用户名称
    userAvatar?: string;              // 用户头像
    playerId?: number;                // 陪玩师ID
    playerName?: string;              // 陪玩师名称
    playerAvatar?: string;            // 陪玩师头像
    messageCount: number;             // 消息数量
    lastMessageAt?: string;           // 最后消息时间
    lastMessageContent?: string;      // 最后消息内容
    status: ChatConversationStatus;
    createdAt: string;
    updatedAt?: string;
}

/**
 * Chat Conversation Detail
 * 聊天会话详情（含订单信息）
 */
export interface ChatConversationDetail extends ChatConversation {
    order?: {
        id: number;
        orderNo: string;
        status: string;
        gameName?: string;
        serviceItemName?: string;
    };
}

/**
 * Chat Conversation Query Params
 * 聊天会话查询参数
 */
export interface ChatConversationQueryParams {
    page?: number;
    pageSize?: number;
    type?: ChatConversationType;
    status?: ChatConversationStatus;
    userId?: number;
    playerId?: number;
    orderId?: number;
    keyword?: string;
    dateFrom?: string;
    dateTo?: string;
}

/**
 * Close Conversation Request
 * 关闭会话请求
 */
export interface CloseConversationRequest {
    reason?: string;                  // 关闭原因
    closedBy: number;                 // 操作人ID
}

// ============================================================================
// Chat Message Types (聊天消息)
// ============================================================================

/**
 * Message Sender Type
 * 消息发送者类型
 */
export type MessageSenderType =
    | 'user'        // 用户
    | 'player'      // 陪玩师
    | 'system';     // 系统

/**
 * Message Type
 * 消息类型
 */
export type MessageType =
    | 'text'        // 文本消息
    | 'system'      // 系统消息
    | 'image'       // 图片消息
    | 'file'        // 文件消息
    | 'voice'       // 语音消息
    | 'emoji';      // 表情消息

/**
 * Chat Message
 * 聊天消息
 */
export interface ChatMessage {
    id: number;
    conversationId: number;           // 会话ID
    senderId: number;                 // 发送者ID
    senderName: string;               // 发送者名称
    senderAvatar?: string;            // 发送者头像
    senderType: MessageSenderType;    // 发送者类型
    content: string;                  // 消息内容
    messageType: MessageType;         // 消息类型
    imageUrl?: string;                // 图片URL（图片消息）
    replyToId?: number;               // 回复的消息ID
    isDeleted: boolean;               // 是否已删除
    createdAt: string;
}

/**
 * Chat Message with Reply
 * 带回复信息的聊天消息
 */
export interface ChatMessageWithReply extends ChatMessage {
    replyToMessage?: ChatMessage;     // 被回复的消息
}

/**
 * Chat Message Query Params
 * 聊天消息查询参数
 */
export interface ChatMessageQueryParams {
    page?: number;
    pageSize?: number;
    conversationId?: number;          // 会话ID
    senderId?: number;                // 发送者ID
    senderType?: MessageSenderType;   // 发送者类型
    messageType?: MessageType;        // 消息类型
    keyword?: string;
    dateFrom?: string;
    dateTo?: string;
}

// ============================================================================
// Chat Broadcast Types (系统广播)
// ============================================================================

/**
 * Broadcast Target Type
 * 广播目标类型
 */
export type BroadcastTargetType =
    | 'all'         // 全部用户
    | 'users'       // 普通用户
    | 'players';    // 陪玩师

/**
 * Broadcast Priority
 * 广播优先级
 */
export type BroadcastPriority =
    | 'low'         // 低
    | 'normal'      // 普通
    | 'high';       // 高

/**
 * Send Broadcast Request
 * 发送广播请求
 */
export interface SendBroadcastRequest {
    title: string;                    // 广播标题
    content: string;                  // 广播内容
    targetType: BroadcastTargetType;  // 目标类型
    targetIds?: number[];             // 目标用户ID列表（可选）
    priority?: BroadcastPriority;     // 优先级
}

/**
 * Broadcast Record
 * 广播记录
 */
export interface BroadcastRecord {
    id: number;
    title: string;                    // 广播标题
    content: string;                  // 广播内容
    targetType: BroadcastTargetType;  // 目标类型
    targetIds?: number[];             // 目标用户ID列表
    priority: BroadcastPriority;      // 优先级
    sentBy: number;                   // 发送人ID
    sentByName?: string;              // 发送人名称
    reachCount: number;               // 触达人数
    readCount: number;                // 已读人数
    createdAt: string;
}

// ============================================================================
// Mute Management Types (禁言管理)
// ============================================================================

/**
 * Mute User Request
 * 禁言用户请求
 */
export interface MuteUserRequest {
    conversationId: number;           // 会话ID
    userId: number;                   // 用户ID
    duration: number;                 // 禁言时长（分钟）
    reason?: string;                  // 禁言原因
}

/**
 * Unmute User Request
 * 解除禁言请求
 */
export interface UnmuteUserRequest {
    conversationId: number;           // 会话ID
    userId: number;                   // 用户ID
}

/**
 * Mute Record
 * 禁言记录
 */
export interface MuteRecord {
    id: number;
    conversationId: number;           // 会话ID
    userId: number;                   // 被禁言用户ID
    userName?: string;                // 用户名称
    mutedBy: number;                  // 操作人ID
    mutedByName?: string;             // 操作人名称
    reason?: string;                  // 禁言原因
    mutedAt: string;                  // 禁言时间
    expiresAt: string;                // 到期时间
    unmutedAt?: string;               // 解除时间
    isActive: boolean;                // 是否生效中
}

// ============================================================================
// Chat Statistics Types (聊天统计)
// ============================================================================

/**
 * Chat Statistics Overview
 * 聊天统计概览
 */
export interface ChatStatsOverview {
    totalConversations: number;       // 总会话数
    activeConversations: number;      // 活跃会话数
    totalMessages: number;            // 总消息数
    todayMessages: number;            // 今日消息数
    totalUsers: number;               // 总用户数
    onlineUsers: number;              // 在线用户数
}

/**
 * Chat Trend Data
 * 聊天趋势数据
 */
export interface ChatTrendData {
    date: string;                     // 日期
    conversations: number;            // 会话数
    messages: number;                 // 消息数
    activeUsers: number;              // 活跃用户数
}

/**
 * Chat Statistics Response
 * 聊天统计响应
 */
export interface ChatStatsResponse {
    overview: ChatStatsOverview;
    trends: ChatTrendData[];
}

/**
 * Conversation Statistics
 * 会话统计
 */
export interface ConversationStats {
    total: number;                    // 总会话数
    active: number;                   // 活跃会话数
    closed: number;                   // 已关闭会话数
    trend: ChatTrendData[];           // 趋势数据
}

/**
 * Message Statistics
 * 消息统计
 */
export interface MessageStats {
    total: number;                    // 总消息数
    today: number;                    // 今日消息数
    byType: Record<string, number>;   // 按类型统计
    trend: ChatTrendData[];           // 趋势数据
}

/**
 * User Activity Statistics
 * 用户活跃统计
 */
export interface UserActivityStats {
    totalUsers: number;               // 总用户数
    activeUsers: number;              // 活跃用户数
    onlineUsers: number;              // 在线用户数
    trend: ChatTrendData[];           // 趋势数据
}

// ============================================================================
// Display Constants (UI展示用常量)
// ============================================================================

/**
 * Conversation Type Display Labels
 * 会话类型显示标签
 */
export const CONVERSATION_TYPE_LABELS: Record<ChatConversationType, string> = {
    user_order: '订单会话',
    group_chat: '群聊会话',
};

/**
 * Conversation Status Display Labels
 * 会话状态显示标签
 */
export const CONVERSATION_STATUS_LABELS: Record<ChatConversationStatus, string> = {
    active: '进行中',
    closed: '已关闭',
};

/**
 * Conversation Status Colors
 * 会话状态颜色（用于Ant Design Tag）
 */
export const CONVERSATION_STATUS_COLORS: Record<ChatConversationStatus, string> = {
    active: 'green',
    closed: 'default',
};

/**
 * Sender Type Display Labels
 * 发送者类型显示标签
 */
export const SENDER_TYPE_LABELS: Record<MessageSenderType, string> = {
    user: '用户',
    player: '陪玩师',
    system: '系统',
};

/**
 * Sender Type Colors
 * 发送者类型颜色
 */
export const SENDER_TYPE_COLORS: Record<MessageSenderType, string> = {
    user: 'blue',
    player: 'purple',
    system: 'default',
};

/**
 * Message Type Display Labels
 * 消息类型显示标签
 */
export const MESSAGE_TYPE_LABELS: Record<MessageType, string> = {
    text: '文本',
    system: '系统消息',
    image: '图片',
    file: '文件',
    voice: '语音',
    emoji: '表情',
};

/**
 * Message Type Icons
 * 消息类型图标
 */
export const MESSAGE_TYPE_ICONS: Record<MessageType, string> = {
    text: 'message',
    system: 'notification',
    image: 'picture',
    file: 'file',
    voice: 'audio',
    emoji: 'smile',
};

/**
 * Broadcast Target Type Display Labels
 * 广播目标类型显示标签
 */
export const BROADCAST_TARGET_LABELS: Record<BroadcastTargetType, string> = {
    all: '全部用户',
    users: '普通用户',
    players: '陪玩师',
};

/**
 * Broadcast Priority Display Labels
 * 广播优先级显示标签
 */
export const BROADCAST_PRIORITY_LABELS: Record<BroadcastPriority, string> = {
    low: '低',
    normal: '普通',
    high: '高',
};

/**
 * Broadcast Priority Colors
 * 广播优先级颜色
 */
export const BROADCAST_PRIORITY_COLORS: Record<BroadcastPriority, string> = {
    low: 'default',
    normal: 'blue',
    high: 'red',
};
