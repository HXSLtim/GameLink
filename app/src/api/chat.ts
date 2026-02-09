/**
 * 聊天相关 API
 * Discord 风格频道系统
 */

import { get, post, type RequestConfig } from './request'
import type { ChatMemberRole, ChatMessageStatus, ChatMessageType, ChatType } from '@/types/message'

// 群组类型
export type ChatGroupType = ChatType

// 聊天群组
export interface ChatGroup {
  id: number
  type: ChatGroupType
  name: string
  avatar?: string
  memberCount: number
  lastMessage?: ChatMessage
  unreadCount: number
  orderId?: number
  createdAt: string
  updatedAt: string
}

// 群组成员
export interface ChatMember {
  id: number
  userId: number
  nickname: string
  avatar?: string
  role: ChatMemberRole
  joinedAt: string
}

// 群组详情
export interface ChatGroupDetail extends ChatGroup {
  members: ChatMember[]
}

// 聊天消息
export interface ChatMessage {
  id: string
  groupId: number
  senderId: number
  senderName: string
  senderAvatar?: string
  type: ChatMessageType
  content: string
  duration?: number
  orderId?: number
  status?: ChatMessageStatus
  createdAt: string
}

// 公共频道
export interface PublicChannel {
  id: number
  name: string
  description?: string
  avatar?: string
  icon?: string
  memberCount: number
  maxMembers?: number
  isActive: boolean
  isJoined: boolean
  gameId?: number
  gameName?: string
}

// 公共频道查询参数
export interface PublicChannelParams {
  gameId?: number
  keyword?: string
  page?: number
  page_size?: number
}

/**
 * 获取聊天群组列表
 */
export function getChatGroups(
  params?: { page?: number; page_size?: number },
  config?: Partial<RequestConfig>
) {
  return get<ChatGroup[]>('/users/chat/groups', params, config)
}

/**
 * 创建聊天群组（私聊或订单群组）
 */
export function createChatGroup(
  data: {
  targetUserId?: number
  groupType: ChatGroupType
  orderId?: number
},
  config?: Partial<RequestConfig>
) {
  return post<ChatGroup>('/users/chat/groups', data, config)
}

/**
 * 获取群组详情
 */
export function getChatGroupDetail(groupId: number, config?: Partial<RequestConfig>) {
  return get<ChatGroupDetail>(`/users/chat/groups/${groupId}`, undefined, config)
}

/**
 * 加入公共频道
 */
export function joinChatGroup(groupId: number, config?: Partial<RequestConfig>) {
  return post<void>(`/users/chat/groups/${groupId}/join`, undefined, config)
}

/**
 * 离开群组/频道
 */
export function leaveChatGroup(groupId: number, config?: Partial<RequestConfig>) {
  return post<void>(`/users/chat/groups/${groupId}/leave`, undefined, config)
}

/**
 * 标记消息已读
 */
export function markMessagesRead(
  groupId: number,
  messageId: string,
  config?: Partial<RequestConfig>
) {
  return post<void>(`/users/chat/groups/${groupId}/read`, { messageId }, config)
}

/**
 * 获取群组消息列表
 */
export function getChatMessages(
  groupId: number,
  params?: { 
    page?: number
    page_size?: number
    beforeId?: string 
  },
  config?: Partial<RequestConfig>
) {
  return get<ChatMessage[]>(`/users/chat/groups/${groupId}/messages`, params, config)
}

/**
 * 发送消息
 */
export function sendChatMessage(
  groupId: number,
  data: {
    type: 'text' | 'image' | 'voice'
    content: string
    duration?: number
  },
  config?: Partial<RequestConfig>
) {
  return post<ChatMessage>(`/users/chat/groups/${groupId}/messages`, data, config)
}

/**
 * 举报消息
 */
export function reportMessage(messageId: string, reason: string, config?: Partial<RequestConfig>) {
  return post<void>(`/users/chat/messages/${messageId}/report`, { reason }, config)
}

/**
 * 获取公共频道列表
 */
export function getPublicChannels(params?: PublicChannelParams, config?: Partial<RequestConfig>) {
  return get<PublicChannel[]>('/public/chat/public-channels', params, config)
}

export default {
  getChatGroups,
  createChatGroup,
  getChatGroupDetail,
  joinChatGroup,
  leaveChatGroup,
  markMessagesRead,
  getChatMessages,
  sendChatMessage,
  reportMessage,
  getPublicChannels,
}
