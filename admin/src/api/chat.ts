/**
 * 聊天管理 API
 * 包含聊天会话管理、消息管理、系统广播、聊天统计等功能
 */
import apiClient from './client';
import type { ApiResponse } from '@/types/api';

// ==================== 类型定义 ====================

/** 聊天会话类型 */
export type ChatConversationType = 'user_order' | 'group_chat';

/** 聊天会话状态 */
export type ChatConversationStatus = 'active' | 'closed';

/** 消息发送者类型 */
export type MessageSenderType = 'user' | 'player' | 'system';

/** 消息类型 */
export type MessageType = 'text' | 'system' | 'image' | 'file' | 'voice' | 'emoji';

/** 聊天会话 */
export interface ChatConversation {
  id: number;
  type: ChatConversationType;
  orderId?: number;
  orderNo?: string;
  userId?: number;
  userName?: string;
  userAvatar?: string;
  playerId?: number;
  playerName?: string;
  playerAvatar?: string;
  messageCount: number;
  lastMessageAt?: string;
  lastMessageContent?: string;
  status: ChatConversationStatus;
  createdAt: string;
  updatedAt?: string;
}

/** 聊天消息 */
export interface ChatMessage {
  id: number;
  conversationId: number;
  senderId: number;
  senderName: string;
  senderAvatar?: string;
  senderType: MessageSenderType;
  content: string;
  messageType: MessageType;
  imageUrl?: string;
  replyToId?: number;
  isDeleted: boolean;
  createdAt: string;
}

/** 聊天会话查询参数 */
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

/** 聊天消息查询参数 */
export interface ChatMessageQueryParams {
  page?: number;
  pageSize?: number;
  conversationId?: number;
  senderId?: number;
  senderType?: MessageSenderType;
  messageType?: MessageType;
  keyword?: string;
  dateFrom?: string;
  dateTo?: string;
}

/** 关闭会话请求 */
export interface CloseConversationRequest {
  reason?: string;
  closedBy: number;
}

/** 系统广播请求 */
export interface SendBroadcastRequest {
  title: string;
  content: string;
  targetType: 'all' | 'users' | 'players';
  targetIds?: number[];
  priority?: 'low' | 'normal' | 'high';
}

/** 聊天统计概览 */
export interface ChatStatsOverview {
  totalConversations: number;
  activeConversations: number;
  totalMessages: number;
  todayMessages: number;
  totalUsers: number;
  onlineUsers: number;
}

/** 聊天趋势数据 */
export interface ChatTrendData {
  date: string;
  conversations: number;
  messages: number;
  activeUsers: number;
}

/** 聊天统计响应 */
export interface ChatStatsResponse {
  overview: ChatStatsOverview;
  trends: ChatTrendData[];
}

/** 禁言用户请求 */
export interface MuteUserRequest {
  conversationId: number;
  userId: number;
  duration: number; // 分钟
  reason?: string;
}

/** 解除禁言请求 */
export interface UnmuteUserRequest {
  conversationId: number;
  userId: number;
}

// ==================== 会话管理 API ====================

/**
 * 聊天会话管理 API
 */
export const chatConversationApi = {
  /**
   * 获取聊天会话列表
   * @param params 查询参数
   */
  getConversations: (params?: ChatConversationQueryParams) =>
    apiClient.get<ApiResponse<{ items: ChatConversation[]; total: number }>>('/admin/chat/conversations', { params }),

  /**
   * 获取聊天会话详情
   * @param id 会话ID
   */
  getConversation: (id: number) =>
    apiClient.get<ApiResponse<ChatConversation>>(`/admin/chat/conversations/${id}`),

  /**
   * 关闭聊天会话
   * @param id 会话ID
   * @param data 关闭原因
   */
  closeConversation: (id: number, data?: CloseConversationRequest) =>
    apiClient.post<ApiResponse<void>>(`/admin/chat/conversations/${id}/close`, data),

  /**
   * 重新打开聊天会话
   * @param id 会话ID
   */
  reopenConversation: (id: number) =>
    apiClient.post<ApiResponse<void>>(`/admin/chat/conversations/${id}/reopen`),

  /**
   * 批量关闭会话
   * @param conversationIds 会话ID列表
   * @param reason 关闭原因
   */
  batchCloseConversations: (conversationIds: number[], reason?: string) =>
    apiClient.post<ApiResponse<void>>('/admin/chat/conversations/batch-close', { conversationIds, reason }),
};

// ==================== 消息管理 API ====================

/**
 * 聊天消息管理 API
 */
export const chatMessageApi = {
  /**
   * 获取会话消息列表
   * @param conversationId 会话ID
   * @param params 查询参数
   */
  getMessages: (conversationId: number, params?: ChatMessageQueryParams) =>
    apiClient.get<ApiResponse<{ items: ChatMessage[]; total: number }>>(
      `/admin/chat/conversations/${conversationId}/messages`,
      { params }
    ),

  /**
   * 获取所有消息列表（跨会话）
   * @param params 查询参数
   */
  getAllMessages: (params?: ChatMessageQueryParams) =>
    apiClient.get<ApiResponse<{ items: ChatMessage[]; total: number }>>('/admin/chat/messages', { params }),

  /**
   * 删除消息
   * @param id 消息ID
   * @param reason 删除原因
   */
  deleteMessage: (id: number, reason?: string) =>
    apiClient.delete<ApiResponse<void>>(`/admin/chat/messages/${id}`, { data: { reason } }),

  /**
   * 批量删除消息
   * @param messageIds 消息ID列表
   * @param reason 删除原因
   */
  batchDeleteMessages: (messageIds: number[], reason?: string) =>
    apiClient.post<ApiResponse<void>>('/admin/chat/messages/batch-delete', { messageIds, reason }),
};

// ==================== 系统广播 API ====================

/**
 * 系统广播管理 API
 */
export const chatBroadcastApi = {
  /**
   * 发送系统广播
   * @param data 广播内容
   */
  sendBroadcast: (data: SendBroadcastRequest) =>
    apiClient.post<ApiResponse<void>>('/admin/chat/broadcasts', data),

  /**
   * 获取广播历史
   * @param params 分页参数
   */
  getBroadcastHistory: (params?: { page?: number; pageSize?: number }) =>
    apiClient.get<ApiResponse<{ items: BroadcastRecord[]; total: number }>>('/admin/chat/broadcasts', { params }),

  /**
   * 获取广播详情
   * @param id 广播ID
   */
  getBroadcast: (id: number) =>
    apiClient.get<ApiResponse<BroadcastRecord>>(`/admin/chat/broadcasts/${id}`),
};

/** 广播记录 */
export interface BroadcastRecord {
  id: number;
  title: string;
  content: string;
  targetType: 'all' | 'users' | 'players';
  targetIds?: number[];
  priority: 'low' | 'normal' | 'high';
  sentBy: number;
  sentByName?: string;
  reachCount: number;
  readCount: number;
  createdAt: string;
}

// ==================== 用户管理 API ====================

/**
 * 聊天用户管理 API
 */
export const chatUserApi = {
  /**
   * 禁言用户
   * @param data 禁言请求
   */
  muteUser: (data: MuteUserRequest) =>
    apiClient.post<ApiResponse<void>>('/admin/chat/users/mute', data),

  /**
   * 解除禁言
   * @param data 解除禁言请求
   */
  unmuteUser: (data: UnmuteUserRequest) =>
    apiClient.post<ApiResponse<void>>('/admin/chat/users/unmute', data),

  /**
   * 获取用户禁言记录
   * @param userId 用户ID
   */
  getUserMuteRecords: (userId: number) =>
    apiClient.get<ApiResponse<MuteRecord[]>>(`/admin/chat/users/${userId}/mute-records`),
};

/** 禁言记录 */
export interface MuteRecord {
  id: number;
  conversationId: number;
  userId: number;
  userName?: string;
  mutedBy: number;
  mutedByName?: string;
  reason?: string;
  mutedAt: string;
  expiresAt: string;
  unmutedAt?: string;
  isActive: boolean;
}

// ==================== 聊天统计 API ====================

/**
 * 聊天统计 API
 */
export const chatStatsApi = {
  /**
   * 获取聊天统计概览
   * @param days 统计天数（默认7天）
   */
  getStats: (days?: number) =>
    apiClient.get<ApiResponse<ChatStatsResponse>>('/admin/chat/stats', { params: { days } }),

  /**
   * 获取会话统计
   * @param days 统计天数
   */
  getConversationStats: (days?: number) =>
    apiClient.get<ApiResponse<{ total: number; active: number; closed: number; trend: ChatTrendData[] }>>(
      '/admin/chat/stats/conversations',
      { params: { days } }
    ),

  /**
   * 获取消息统计
   * @param days 统计天数
   */
  getMessageStats: (days?: number) =>
    apiClient.get<ApiResponse<{ total: number; today: number; byType: Record<string, number>; trend: ChatTrendData[] }>>(
      '/admin/chat/stats/messages',
      { params: { days } }
    ),

  /**
   * 获取用户活跃统计
   * @param days 统计天数
   */
  getUserActivityStats: (days?: number) =>
    apiClient.get<ApiResponse<{ totalUsers: number; activeUsers: number; onlineUsers: number; trend: ChatTrendData[] }>>(
      '/admin/chat/stats/user-activity',
      { params: { days } }
    ),

  /**
   * 导出统计数据
   * @param days 统计天数
   */
  exportStats: (days?: number) =>
    apiClient.get('/admin/chat/stats/export', {
      params: { days },
      responseType: 'blob',
    }),
};

// ==================== 常量映射 ====================

/** 会话类型显示文本 */
export const CONVERSATION_TYPE_TEXT: Record<ChatConversationType, string> = {
  user_order: '订单会话',
  group_chat: '群聊会话',
};

/** 会话状态显示文本 */
export const CONVERSATION_STATUS_TEXT: Record<ChatConversationStatus, string> = {
  active: '进行中',
  closed: '已关闭',
};

/** 会话状态颜色 */
export const CONVERSATION_STATUS_COLOR: Record<ChatConversationStatus, string> = {
  active: 'green',
  closed: 'default',
};

/** 发送者类型显示文本 */
export const SENDER_TYPE_TEXT: Record<MessageSenderType, string> = {
  user: '用户',
  player: '陪玩师',
  system: '系统',
};

/** 发送者类型颜色 */
export const SENDER_TYPE_COLOR: Record<MessageSenderType, string> = {
  user: 'blue',
  player: 'purple',
  system: 'default',
};

/** 消息类型显示文本 */
export const MESSAGE_TYPE_TEXT: Record<MessageType, string> = {
  text: '文本',
  system: '系统消息',
  image: '图片',
  file: '文件',
  voice: '语音',
  emoji: '表情',
};

/** 广播目标类型显示文本 */
export const BROADCAST_TARGET_TEXT: Record<'all' | 'users' | 'players', string> = {
  all: '全部用户',
  users: '普通用户',
  players: '陪玩师',
};

/** 广播优先级显示文本 */
export const BROADCAST_PRIORITY_TEXT: Record<'low' | 'normal' | 'high', string> = {
  low: '低',
  normal: '普通',
  high: '高',
};

/** 广播优先级颜色 */
export const BROADCAST_PRIORITY_COLOR: Record<'low' | 'normal' | 'high', string> = {
  low: 'default',
  normal: 'blue',
  high: 'red',
};
