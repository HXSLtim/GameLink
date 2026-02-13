/**
 * 客服聊天 Hook
 */
import { ref, nextTick, onMounted, onUnmounted } from 'vue'
import { useUserStore } from '@/store/user'
import {
  getCustomerServiceSession,
  getCustomerServiceMessages,
  sendCustomerServiceMessage,
} from '@/api/customerService'
import type { ServiceMessage, QuickQuestion } from '@/types/CustomerService'

export function useCustomerService() {
  const userStore = useUserStore()

  // 状态
  const isOnline = ref(false)
  const sending = ref(false)
  const inputContent = ref('')
  const scrollToView = ref('')
  const sessionGroupId = ref(0)
  
  // 消息列表
  const messages = ref<ServiceMessage[]>([])
  let refreshTimer: number | null = null
  
  // 快捷问题（icon 为 uv-icon 名称）
  const quickQuestions = ref<QuickQuestion[]>([
    {
      id: 1,
      icon: 'wallet',
      text: '退款问题',
      answer: '关于退款问题：\n1. 待支付订单可直接取消\n2. 已支付未开始的订单可申请全额退款\n3. 服务进行中的订单需协商处理\n4. 退款一般1-7个工作日到账\n\n如需进一步帮助，请提供您的订单号。',
    },
    {
      id: 2,
      icon: 'list',
      text: '订单问题',
      answer: '关于订单问题，我可以帮您：\n1. 查询订单状态\n2. 处理订单争议\n3. 修改订单信息\n\n请提供您的订单号，我来为您查询。',
    },
    {
      id: 3,
      icon: 'account',
      text: '账号问题',
      answer: '关于账号问题：\n1. 修改手机号：设置-账号安全-更换手机\n2. 忘记密码：登录页-忘记密码\n3. 账号被盗：请立即修改密码并联系我们\n\n请问您遇到的是哪种账号问题？',
    },
    {
      id: 4,
      icon: 'info-circle',
      text: '投诉举报',
      answer: '感谢您的反馈。请提供以下信息：\n1. 被投诉的订单号或用户ID\n2. 具体问题描述\n3. 相关截图证据\n\n我们会在24小时内处理您的投诉。',
    },
  ])

  const ensureGreeting = () => {
    if (messages.value.length > 0) return
    messages.value = [
      {
        id: 1,
        content: '您好！我是 GameLink 官方客服，很高兴为您服务。请问有什么可以帮助您的吗？',
        isMe: false,
        createdAt: new Date().toISOString(),
      },
    ]
  }
  
  // 滚动到底部
  const scrollToBottom = () => {
    nextTick(() => {
      scrollToView.value = ''
      setTimeout(() => {
        scrollToView.value = 'chat-bottom'
      }, 50)
    })
  }

  const loadSession = async () => {
    if (!userStore.isLoggedIn) {
      ensureGreeting()
      return
    }

    try {
      const res = await getCustomerServiceSession({ showError: false })
      const session = res.data
      if (session) {
        sessionGroupId.value = session.groupId
        isOnline.value = !!session.isOnline
      }
    } catch (error) {
      ensureGreeting()
    }
  }

  const loadMessages = async (silent = true) => {
    if (!userStore.isLoggedIn) {
      ensureGreeting()
      return
    }

    try {
      const res = await getCustomerServiceMessages(
        {
          page: 1,
          pageSize: 100,
        },
        { showError: !silent }
      )

      const currentUserId = userStore.userInfo?.id || 0
      const list = res.data?.messages || []
      sessionGroupId.value = res.data?.groupId || sessionGroupId.value

      messages.value = list.map((item) => ({
        id: item.id,
        content: item.content,
        isMe: item.senderId === currentUserId,
        createdAt: item.createdAt,
      }))

      ensureGreeting()
      scrollToBottom()
    } catch (error) {
      if (!silent) {
        uni.showToast({ title: '加载客服消息失败', icon: 'none' })
      }
      ensureGreeting()
    }
  }

  const startRefresh = () => {
    if (refreshTimer) return
    refreshTimer = setInterval(() => {
      loadMessages(true)
    }, 8000) as unknown as number
  }

  const stopRefresh = () => {
    if (!refreshTimer) return
    clearInterval(refreshTimer)
    refreshTimer = null
  }
  
  // 选择快捷问题
  const selectQuestion = (item: QuickQuestion) => {
    inputContent.value = item.text
    sendMessage()
  }
  
  // 发送消息
  const sendMessage = async () => {
    const content = inputContent.value.trim()
    if (!content || sending.value) return

    if (!userStore.isLoggedIn) {
      uni.showToast({ title: '请先登录', icon: 'none' })
      return
    }

    sending.value = true

    const tempId = Date.now()
    const tempMessage: ServiceMessage = {
      id: tempId,
      content,
      isMe: true,
      createdAt: new Date().toISOString(),
    }
    messages.value.push(tempMessage)
    inputContent.value = ''
    scrollToBottom()

    try {
      const res = await sendCustomerServiceMessage({ content }, { showError: false })
      const saved = res.data
      if (saved) {
        const idx = messages.value.findIndex(item => item.id === tempId)
        if (idx >= 0) {
          messages.value[idx] = {
            id: saved.id,
            content: saved.content,
            isMe: true,
            createdAt: saved.createdAt,
          }
        }
        sessionGroupId.value = saved.groupId || sessionGroupId.value
      }

      await loadMessages(true)
    } catch (error) {
      uni.showToast({ title: '发送失败，请稍后重试', icon: 'none' })
    } finally {
      sending.value = false
    }
  }
  
  // 导航
  const goBack = () => uni.navigateBack()
  
  onMounted(async () => {
    await loadSession()
    await loadMessages(true)
    scrollToBottom()
    startRefresh()
  })

  onUnmounted(() => {
    stopRefresh()
  })
  
  return {
    isOnline,
    sending,
    inputContent,
    scrollToView,
    sessionGroupId,
    messages,
    quickQuestions,
    loadMessages,
    selectQuestion,
    sendMessage,
    goBack,
  }
}
