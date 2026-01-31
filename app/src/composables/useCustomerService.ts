/**
 * 客服聊天 Hook
 */
import { ref, nextTick, onMounted } from 'vue'
import type { ServiceMessage } from '@/components/ServiceChatMessage/index.vue'
import type { QuickQuestion } from '@/components/QuickQuestionBar/index.vue'

export function useCustomerService() {
  // 状态
  const isOnline = ref(true)
  const sending = ref(false)
  const inputContent = ref('')
  const scrollToView = ref('')
  
  // 消息列表
  const messages = ref<ServiceMessage[]>([
    {
      id: 1,
      content: '您好！我是 GameLink 官方客服，很高兴为您服务。请问有什么可以帮助您的吗？',
      isMe: false,
      createdAt: new Date().toISOString(),
    },
  ])
  
  // 快捷问题
  const quickQuestions = ref<QuickQuestion[]>([
    {
      id: 1,
      icon: '💳',
      text: '退款问题',
      answer: '关于退款问题：\n1. 待支付订单可直接取消\n2. 已支付未开始的订单可申请全额退款\n3. 服务进行中的订单需协商处理\n4. 退款一般1-7个工作日到账\n\n如需进一步帮助，请提供您的订单号。',
    },
    {
      id: 2,
      icon: '📋',
      text: '订单问题',
      answer: '关于订单问题，我可以帮您：\n1. 查询订单状态\n2. 处理订单争议\n3. 修改订单信息\n\n请提供您的订单号，我来为您查询。',
    },
    {
      id: 3,
      icon: '👤',
      text: '账号问题',
      answer: '关于账号问题：\n1. 修改手机号：设置-账号安全-更换手机\n2. 忘记密码：登录页-忘记密码\n3. 账号被盗：请立即修改密码并联系我们\n\n请问您遇到的是哪种账号问题？',
    },
    {
      id: 4,
      icon: '⚠️',
      text: '投诉举报',
      answer: '感谢您的反馈。请提供以下信息：\n1. 被投诉的订单号或用户ID\n2. 具体问题描述\n3. 相关截图证据\n\n我们会在24小时内处理您的投诉。',
    },
  ])
  
  // 滚动到底部
  const scrollToBottom = () => {
    nextTick(() => {
      scrollToView.value = ''
      setTimeout(() => {
        scrollToView.value = 'chat-bottom'
      }, 50)
    })
  }
  
  // 模拟客服回复
  const simulateReply = (content: string) => {
    setTimeout(() => {
      messages.value.push({
        id: Date.now() + 1,
        content,
        isMe: false,
        createdAt: new Date().toISOString(),
      })
      scrollToBottom()
      sending.value = false
    }, 500 + Math.random() * 500)
  }
  
  // 选择快捷问题
  const selectQuestion = (item: QuickQuestion) => {
    messages.value.push({
      id: Date.now(),
      content: item.text,
      isMe: true,
      createdAt: new Date().toISOString(),
    })
    scrollToBottom()
    simulateReply(item.answer)
  }
  
  // 发送消息
  const sendMessage = () => {
    const content = inputContent.value.trim()
    if (!content || sending.value) return
    
    sending.value = true
    
    messages.value.push({
      id: Date.now(),
      content,
      isMe: true,
      createdAt: new Date().toISOString(),
    })
    
    inputContent.value = ''
    scrollToBottom()
    
    simulateReply('感谢您的咨询，客服正在为您查询中，请稍候...\n\n如需紧急帮助，您也可以拨打客服热线：400-XXX-XXXX')
  }
  
  // 导航
  const goBack = () => uni.navigateBack()
  
  onMounted(() => {
    scrollToBottom()
  })
  
  return {
    isOnline,
    sending,
    inputContent,
    scrollToView,
    messages,
    quickQuestions,
    selectQuestion,
    sendMessage,
    goBack,
  }
}
