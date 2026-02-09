export type ChatMessageType = 'text' | 'image' | 'voice' | 'order' | 'system'
export type ChatMessageStatus = 'sending' | 'sent' | 'read' | 'failed'

export interface ChatMessageData {
  id: string
  type: ChatMessageType
  content: string
  senderId: number
  senderName?: string
  senderAvatar?: string
  createdAt: string
  status?: ChatMessageStatus
  duration?: number
  orderId?: number
}

export type ChatShowTimeFn = (message: ChatMessageData, index: number) => boolean

export interface MessageData {
  id: number
  conversationId: number
  avatar: string
  name: string
  lastMessage: string
  lastMessageType?: ChatMessageType
  lastTime: number
  unread: number
  type: 'chat' | 'system' | 'order'
}

export type ChatType = 'private' | 'order' | 'public'
export type WsStatus = 'connecting' | 'connected' | 'disconnected'
export type ChatMemberRole = 'owner' | 'admin' | 'member'

export interface ChatInfo {
  id: number
  type: ChatType
  name: string
  targetId?: number
  orderId?: number
  isOnline?: boolean
  memberCount?: number
}

export interface WsMessage {
  type: 'chat_message' | 'order_status' | 'order_new' | 'notification' | 'pong'
  timestamp: number
  data: unknown
}

export type WsEventHandler = (message: WsMessage) => void
