/**
 * 消息列表专用 Hook
 * 封装消息列表的数据加载和标签切换逻辑
 */
import { ref, reactive } from 'vue'
import { getChatGroups, type ChatGroup } from '@/api/chat'
import { getNotifications, getUnreadCount, type Notification } from '@/api/notification'
import { useUserStore } from '@/store/user'
import type { PageStateType } from '@/components/PageState.vue'
import type { MessageData } from '@/components/MessageItem/index.vue'
import type { TabItem } from '@/components/TabsBar/index.vue'

export function useMessageList() {
  const userStore = useUserStore()
  
  // 标签配置
  const tabs = reactive<TabItem[]>([
    { key: 'chat', label: '聊天', badge: 0 },
    { key: 'system', label: '通知', badge: 0 },
  ])
  
  // 状态
  const currentTab = ref('chat')
  const pageState = ref<PageStateType>('loading')
  const errorMessage = ref('')
  const messages = ref<MessageData[]>([])
  
  // 切换标签
  const switchTab = (key: string) => {
    currentTab.value = key
    loadMessages()
  }
  
  // 加载消息
  const loadMessages = async () => {
    if (!userStore.isLoggedIn) {
      pageState.value = 'login'
      return
    }
    
    pageState.value = 'loading'
    errorMessage.value = ''
    
    try {
      if (currentTab.value === 'chat') {
        await loadChatMessages()
      } else {
        await loadSystemNotifications()
      }
      
      // 设置页面状态
      pageState.value = messages.value.length === 0 ? 'empty' : 'content'
      
      // 加载未读数量
      await loadUnreadCounts()
    } catch (error: any) {
      console.error('加载消息失败:', error)
      pageState.value = 'error'
      errorMessage.value = error?.message || '网络请求失败，请检查网络后重试'
    }
  }
  
  // 加载聊天消息
  const loadChatMessages = async () => {
    const res = await getChatGroups(undefined, { showError: false })
    const groups = (res.data as any)?.groups || res.data || []
    
    messages.value = groups.map((g: ChatGroup): MessageData => ({
      id: g.id,
      conversationId: g.id,
      avatar: g.avatar || `https://picsum.photos/100?random=${g.id}`,
      name: g.name,
      lastMessage: g.lastMessage?.content || '',
      lastTime: new Date(g.updatedAt).getTime(),
      unread: g.unreadCount,
      type: 'chat',
    }))
    
    tabs[0].badge = messages.value.reduce((sum, m) => sum + m.unread, 0)
  }
  
  // 加载系统通知
  const loadSystemNotifications = async () => {
    const res = await getNotifications({}, { showError: false })
    const notifications = (res.data as any)?.items || res.data || []
    
    messages.value = notifications.map((n: Notification): MessageData => ({
      id: n.id,
      conversationId: 0,
      avatar: '/static/icons/system.png',
      name: n.title,
      lastMessage: n.content,
      lastTime: new Date(n.createdAt).getTime(),
      unread: n.isRead ? 0 : 1,
      type: 'system',
    }))
  }
  
  // 加载未读数量
  const loadUnreadCounts = async () => {
    try {
      const unreadRes = await getUnreadCount({ showError: false })
      if (unreadRes.data) {
        tabs[0].badge = unreadRes.data.chat
        tabs[1].badge = unreadRes.data.system
      }
    } catch (e) {
      // 未读数加载失败不影响主流程
    }
  }
  
  // 跳转聊天
  const goToChat = (item: MessageData) => {
    if (item.type === 'chat') {
      uni.navigateTo({ url: `/pages/message/chat/index?groupId=${item.conversationId}` })
    } else if (item.type === 'order') {
      uni.navigateTo({ url: '/pages/order/list/index' })
    }
  }
  
  // 跳转陪玩师列表
  const goToPlayerList = () => {
    uni.switchTab({ url: '/pages/player/list/index' })
  }
  
  return {
    // 数据
    messages,
    pageState,
    errorMessage,
    
    // 标签
    tabs,
    currentTab,
    
    // 方法
    switchTab,
    loadMessages,
    goToChat,
    goToPlayerList,
  }
}
