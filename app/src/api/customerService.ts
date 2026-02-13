import { get, post, type RequestConfig } from './request'

export interface CustomerServiceAgent {
  userId: number
  nickname: string
  avatarUrl?: string
  status: string
}

export interface CustomerServiceSession {
  groupId: number
  agent: CustomerServiceAgent
  isOnline: boolean
}

export interface CustomerServiceMessage {
  id: number
  groupId: number
  senderId: number
  content: string
  messageType: 'text' | 'image' | 'file' | 'system' | 'voice' | 'emoji'
  isMe: boolean
  createdAt: string
}

export interface CustomerServiceMessageList {
  groupId: number
  messages: CustomerServiceMessage[]
  total: number
  page: number
  pageSize: number
  hasMore: boolean
}

export function getCustomerServiceSession(config?: Partial<RequestConfig>) {
  return get<CustomerServiceSession>('/user/customer-service/session', undefined, config)
}

export function getCustomerServiceMessages(
  params?: {
    page?: number
    pageSize?: number
    beforeId?: number
  },
  config?: Partial<RequestConfig>
) {
  return get<CustomerServiceMessageList>('/user/customer-service/messages', params, config)
}

export function sendCustomerServiceMessage(
  data: {
    content: string
  },
  config?: Partial<RequestConfig>
) {
  return post<CustomerServiceMessage>('/user/customer-service/messages', data, config)
}
