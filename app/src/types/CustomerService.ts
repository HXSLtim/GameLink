export interface ServiceMessage {
  id: number
  content: string
  isMe: boolean
  createdAt: string
}

export interface QuickQuestion {
  id: number
  icon: string
  text: string
  answer: string
}
