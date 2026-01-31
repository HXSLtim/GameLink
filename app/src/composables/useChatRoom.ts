/**
 * 聊天室专用 Hook
 */
import { ref, computed, nextTick } from 'vue'
import { useUserStore } from '@/store/user'
import { 
  createChatGroup,
  getChatGroupDetail, 
  getChatMessages, 
  sendChatMessage, 
  markMessagesRead 
} from '@/api/chat'
import type { ChatMessageData } from '@/components/ChatMessageBubble/index.vue'

type ChatType = 'private' | 'order' | 'public'
type WsStatus = 'connecting' | 'connected' | 'disconnected'

interface ChatInfo {
  id: number
  type: ChatType
  name: string
  targetId?: number
  orderId?: number
  isOnline?: boolean
  memberCount?: number
}

export function useChatRoom() {
  const userStore = useUserStore()
  
  // 状态
  const loading = ref(true)
  const loadingHistory = ref(false)
  const hasMoreHistory = ref(true)
  const showMore = ref(false)
  const showMenu = ref(false)
  const recording = ref(false)
  const playingVoiceId = ref<string | null>(null)
  
  // WebSocket
  const wsStatus = ref<WsStatus>('connecting')
  let ws: UniApp.SocketTask | null = null
  let reconnectTimer: number | null = null
  let heartbeatTimer: number | null = null
  
  // 数据
  const chatId = ref(0)
  const chatInfo = ref<ChatInfo>({} as ChatInfo)
  const messages = ref<ChatMessageData[]>([])
  const inputText = ref('')
  const scrollToId = ref('')
  
  // 用户信息
  const currentUserId = computed(() => userStore.userInfo?.id || 0)
  const currentUserName = computed(() => userStore.userInfo?.nickname || '我')
  const currentUserAvatar = computed(() => userStore.userInfo?.avatar || '')
  
  // 初始化聊天
  const initChat = async (options: Record<string, any>) => {
    const resolved = await resolveChatId(options)
    if (!resolved) return false
    
    await loadChatInfo(chatId.value)
    await loadMessages()
    connectWebSocket()
    
    return true
  }
  
  // 解析聊天 ID
  const resolveChatId = async (options: Record<string, any>) => {
    const groupId = options?.groupId || options?.id
    if (groupId) {
      chatId.value = parseInt(groupId as string)
      return true
    }
    
    if (options?.playerId) {
      const targetUserId = parseInt(options.playerId as string)
      try {
        const res = await createChatGroup({
          targetUserId,
          groupType: 'private',
        }, { showError: false })
        if (res?.data?.id) {
          chatId.value = res.data.id
          return true
        }
      } catch (error) {
        console.error('创建聊天失败', error)
      }
    }
    
    uni.showToast({ title: '无效的会话', icon: 'none' })
    return false
  }
  
  // 加载聊天信息
  const loadChatInfo = async (id: number) => {
    loading.value = true
    try {
      const res = await getChatGroupDetail(id, { showError: false })
      if (res.data) {
        const data = res.data as any
        chatInfo.value = {
          id: data.id,
          type: data.groupType || 'private',
          name: data.name || data.targetNickname || '聊天',
          targetId: data.targetUserId,
          orderId: data.orderId,
          isOnline: data.targetIsOnline,
          memberCount: data.memberCount,
        }
      }
    } catch (error) {
      console.error('加载聊天信息失败', error)
    } finally {
      loading.value = false
    }
  }
  
  // 加载消息
  const loadMessages = async (beforeId?: string) => {
    if (beforeId) {
      loadingHistory.value = true
    }
    
    try {
      const res = await getChatMessages(chatId.value, {
        beforeId,
        limit: 20,
      }, { showError: false })
      
      const items = ((res.data as any)?.items || []).map((m: any): ChatMessageData => ({
        id: String(m.id),
        type: m.messageType || 'text',
        content: m.content,
        senderId: m.senderId,
        senderName: m.senderNickname,
        senderAvatar: m.senderAvatar,
        createdAt: m.createdAt,
        status: m.status || 'sent',
        duration: m.voiceDuration,
        orderId: m.orderId,
      }))
      
      if (beforeId) {
        messages.value = [...items, ...messages.value]
      } else {
        messages.value = items
        scrollToBottom()
      }
      
      hasMoreHistory.value = items.length >= 20
      
      // 标记已读
      if (items.length > 0) {
        markMessagesRead(chatId.value).catch(() => {})
      }
    } catch (error) {
      console.error('加载消息失败', error)
    } finally {
      loadingHistory.value = false
    }
  }
  
  // 加载更多历史
  const loadMoreHistory = () => {
    if (loadingHistory.value || !hasMoreHistory.value) return
    const firstMsg = messages.value[0]
    if (firstMsg) {
      loadMessages(firstMsg.id)
    }
  }
  
  // 滚动到底部
  const scrollToBottom = () => {
    nextTick(() => {
      scrollToId.value = 'msg-bottom'
    })
  }
  
  // 发送文本消息
  const sendTextMessage = async () => {
    const text = inputText.value.trim()
    if (!text) return
    
    const tempId = `temp-${Date.now()}`
    const newMessage: ChatMessageData = {
      id: tempId,
      type: 'text',
      content: text,
      senderId: currentUserId.value,
      senderName: currentUserName.value,
      senderAvatar: currentUserAvatar.value,
      createdAt: new Date().toISOString(),
      status: 'sending',
    }
    
    messages.value.push(newMessage)
    inputText.value = ''
    scrollToBottom()
    
    try {
      const res = await sendChatMessage(chatId.value, {
        messageType: 'text',
        content: text,
      })
      
      // 更新消息 ID 和状态
      const idx = messages.value.findIndex(m => m.id === tempId)
      if (idx !== -1 && res.data) {
        messages.value[idx].id = String((res.data as any).id)
        messages.value[idx].status = 'sent'
      }
    } catch (error) {
      const idx = messages.value.findIndex(m => m.id === tempId)
      if (idx !== -1) {
        messages.value[idx].status = 'failed'
      }
    }
  }
  
  // 发送图片
  const sendImage = async (url: string) => {
    const tempId = `temp-${Date.now()}`
    const newMessage: ChatMessageData = {
      id: tempId,
      type: 'image',
      content: url,
      senderId: currentUserId.value,
      senderName: currentUserName.value,
      senderAvatar: currentUserAvatar.value,
      createdAt: new Date().toISOString(),
      status: 'sending',
    }
    
    messages.value.push(newMessage)
    scrollToBottom()
    
    try {
      const res = await sendChatMessage(chatId.value, {
        messageType: 'image',
        content: url,
      })
      
      const idx = messages.value.findIndex(m => m.id === tempId)
      if (idx !== -1 && res.data) {
        messages.value[idx].id = String((res.data as any).id)
        messages.value[idx].status = 'sent'
      }
    } catch (error) {
      const idx = messages.value.findIndex(m => m.id === tempId)
      if (idx !== -1) {
        messages.value[idx].status = 'failed'
      }
    }
  }
  
  // 选择图片
  const chooseImage = () => {
    uni.chooseImage({
      count: 1,
      sizeType: ['compressed'],
      success: (res) => {
        sendImage(res.tempFilePaths[0])
      }
    })
    showMore.value = false
  }
  
  // 拍照
  const takePhoto = () => {
    uni.chooseImage({
      count: 1,
      sourceType: ['camera'],
      success: (res) => {
        sendImage(res.tempFilePaths[0])
      }
    })
    showMore.value = false
  }
  
  // 重发消息
  const resendMessage = (message: ChatMessageData) => {
    if (message.type === 'text') {
      inputText.value = message.content
      messages.value = messages.value.filter(m => m.id !== message.id)
      sendTextMessage()
    }
  }
  
  // 预览图片
  const previewImage = (url: string) => {
    uni.previewImage({
      urls: messages.value.filter(m => m.type === 'image').map(m => m.content),
      current: url,
    })
  }
  
  // 播放语音
  const playVoice = (message: ChatMessageData) => {
    if (playingVoiceId.value === message.id) {
      playingVoiceId.value = null
      // 停止播放
    } else {
      playingVoiceId.value = message.id
      // 播放语音
    }
  }
  
  // 时间分割判断
  const shouldShowTime = (message: ChatMessageData, index: number) => {
    if (index === 0) return true
    const prev = messages.value[index - 1]
    const prevTime = new Date(prev.createdAt).getTime()
    const currTime = new Date(message.createdAt).getTime()
    return currTime - prevTime > 5 * 60 * 1000 // 5分钟
  }
  
  // WebSocket 连接
  const connectWebSocket = () => {
    if (!userStore.token) return
    
    wsStatus.value = 'connecting'
    
    const wsUrl = `wss://api.gamelink.com/ws/chat/${chatId.value}?token=${userStore.token}`
    
    ws = uni.connectSocket({
      url: wsUrl,
      complete: () => {}
    })
    
    ws.onOpen(() => {
      wsStatus.value = 'connected'
      startHeartbeat()
    })
    
    ws.onMessage((res) => {
      try {
        const data = JSON.parse(res.data as string)
        handleWsMessage(data)
      } catch (error) {
        console.error('解析消息失败', error)
      }
    })
    
    ws.onClose(() => {
      wsStatus.value = 'disconnected'
      stopHeartbeat()
      scheduleReconnect()
    })
    
    ws.onError(() => {
      wsStatus.value = 'disconnected'
    })
  }
  
  const handleWsMessage = (data: any) => {
    if (data.type === 'message') {
      const msg: ChatMessageData = {
        id: String(data.id),
        type: data.messageType,
        content: data.content,
        senderId: data.senderId,
        senderName: data.senderNickname,
        senderAvatar: data.senderAvatar,
        createdAt: data.createdAt,
        status: 'sent',
      }
      
      if (msg.senderId !== currentUserId.value) {
        messages.value.push(msg)
        scrollToBottom()
        markMessagesRead(chatId.value).catch(() => {})
      }
    } else if (data.type === 'read') {
      // 更新已读状态
      messages.value.forEach(m => {
        if (m.senderId === currentUserId.value && m.status === 'sent') {
          m.status = 'read'
        }
      })
    } else if (data.type === 'online') {
      chatInfo.value.isOnline = data.online
    }
  }
  
  const startHeartbeat = () => {
    heartbeatTimer = setInterval(() => {
      if (ws && wsStatus.value === 'connected') {
        ws.send({ data: JSON.stringify({ type: 'ping' }) })
      }
    }, 30000) as unknown as number
  }
  
  const stopHeartbeat = () => {
    if (heartbeatTimer) {
      clearInterval(heartbeatTimer)
      heartbeatTimer = null
    }
  }
  
  const scheduleReconnect = () => {
    if (reconnectTimer) return
    reconnectTimer = setTimeout(() => {
      reconnectTimer = null
      if (wsStatus.value === 'disconnected') {
        connectWebSocket()
      }
    }, 3000) as unknown as number
  }
  
  const reconnectWs = () => {
    wsStatus.value = 'connecting'
    connectWebSocket()
  }
  
  const disconnectWs = () => {
    if (ws) {
      ws.close({})
      ws = null
    }
    stopHeartbeat()
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
  }
  
  // 导航和操作
  const goBack = () => uni.navigateBack()
  
  const viewProfile = (userId: number) => {
    uni.navigateTo({ url: `/pages/player/detail/index?id=${userId}` })
  }
  
  const viewOrder = (orderId?: number) => {
    if (orderId) {
      uni.navigateTo({ url: `/pages/order/detail/index?id=${orderId}` })
    } else if (chatInfo.value.orderId) {
      uni.navigateTo({ url: `/pages/order/detail/index?id=${chatInfo.value.orderId}` })
    }
    showMore.value = false
  }
  
  const clearHistory = () => {
    uni.showModal({
      title: '清空聊天记录',
      content: '确定要清空聊天记录吗？',
      success: (res) => {
        if (res.confirm) {
          messages.value = []
          showMenu.value = false
        }
      }
    })
  }
  
  const blockUser = () => {
    uni.showModal({
      title: '拉黑用户',
      content: '拉黑后将无法接收对方消息，确定继续？',
      success: (res) => {
        if (res.confirm) {
          uni.showToast({ title: '已拉黑', icon: 'success' })
          showMenu.value = false
          setTimeout(() => uni.navigateBack(), 500)
        }
      }
    })
  }
  
  const reportChat = () => {
    uni.navigateTo({ url: `/pages/report/index?type=chat&targetId=${chatId.value}` })
    showMore.value = false
  }
  
  return {
    // 状态
    loading,
    loadingHistory,
    hasMoreHistory,
    showMore,
    showMenu,
    recording,
    playingVoiceId,
    wsStatus,
    
    // 数据
    chatId,
    chatInfo,
    messages,
    inputText,
    scrollToId,
    currentUserId,
    currentUserName,
    currentUserAvatar,
    
    // 方法
    initChat,
    loadMoreHistory,
    sendTextMessage,
    chooseImage,
    takePhoto,
    resendMessage,
    previewImage,
    playVoice,
    shouldShowTime,
    reconnectWs,
    disconnectWs,
    goBack,
    viewProfile,
    viewOrder,
    clearHistory,
    blockUser,
    reportChat,
  }
}
